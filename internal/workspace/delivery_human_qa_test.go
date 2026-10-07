package workspace

import "testing"

func TestLatestTaskHumanQAUsesActualTimeAndRejectsAmbiguousLatestAnswers(t *testing.T) {
	qa := func(id, at, outcome string) TaskHumanQARecord {
		return TaskHumanQARecord{ID: HumanQARecordID(id), TaskResultID: "result", ResultOID: "oid", ResultTree: "tree", Outcome: outcome, Actor: TaskHumanActorRecord{CompletedAtUTC: at}}
	}
	older := qa("ffff", "2026-10-01T12:00:00Z", "pass")
	newer := qa("0000", "2026-10-01T12:00:00.1Z", "fail")
	equalPass := qa("ffff", newer.Actor.CompletedAtUTC, "pass")
	equalFail := qa("aaaa", newer.Actor.CompletedAtUTC, "fail")
	later := qa("0001", "2026-10-01T12:00:01Z", "pass")
	unrelated := qa("other", "2026-10-01T13:00:00Z", "fail")
	unrelated.ResultOID = "other"
	for _, tc := range []struct {
		name    string
		records []TaskHumanQARecord
		want    HumanQARecordID
		wantErr bool
	}{
		{"fractional time beats ID order", []TaskHumanQARecord{newer, older}, newer.ID, false},
		{"opposite input order", []TaskHumanQARecord{older, newer}, newer.ID, false},
		{"equal time conflict", []TaskHumanQARecord{newer, equalPass}, "", true},
		{"equal time conflict reversed", []TaskHumanQARecord{equalPass, newer}, "", true},
		{"same outcome stable ID", []TaskHumanQARecord{newer, equalFail}, equalFail.ID, false},
		{"same outcome reversed", []TaskHumanQARecord{equalFail, newer}, equalFail.ID, false},
		{"old ambiguity superseded", []TaskHumanQARecord{newer, equalPass, later}, later.ID, false},
		{"different candidate ignored", []TaskHumanQARecord{later, unrelated}, later.ID, false},
		{"no matching answer", []TaskHumanQARecord{unrelated}, "", false},
		{"invalid time", []TaskHumanQARecord{qa("bad", "invalid", "pass")}, "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := LatestTaskHumanQA(tc.records, "result", "oid", "tree")
			if (err != nil) != tc.wantErr {
				t.Fatalf("selection error = %v; want error %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			if tc.want == "" {
				if got != nil {
					t.Fatalf("unexpected answer: %#v", got)
				}
			} else if got == nil || got.ID != tc.want {
				t.Fatalf("answer = %#v; want %s", got, tc.want)
			}
		})
	}
}

func TestDeliveryIntegrationRejectsOlderPassAfterLaterNativeHumanFailure(t *testing.T) {
	for _, mode := range []string{DeliveryLocalEpic, DeliveryLocalBranch} {
		for _, outcome := range []string{"fail", "blocked"} {
			t.Run(mode+"/"+outcome, func(t *testing.T) {
				ref := "refs/heads/epic"
				if mode == DeliveryLocalBranch {
					ref = "refs/heads/main"
				}
				f := newIntegrationJourneyFixtureHistory(t, false, false, fixtureAgreement(mode, ref))
				input := deliveryIntegrationInput(f)
				before, err := CheckTaskIntegration(f.dependencies, input)
				if err != nil || before.Readback.Classification != "ready" {
					t.Fatalf("initial passed candidate not ready: %#v %v", before.Readback, err)
				}
				input.Apply, input.Confirmation = true, integrationDigest(t, before.Readback)
				registry := mustQueueRegistry(t, f.workItemJourneyFixture)
				qaFile, _ := writeHumanQADraft(t, taskDeliveryFixture{workItemJourneyFixture: f.workItemJourneyFixture, root: f.root}, registry.TaskResults[0], outcome, "report")
				newQA, err := RecordTaskHumanQA(f.dependencies, TaskHumanQARecordInput{TaskID: "task", File: qaFile})
				if err != nil {
					t.Fatal(err)
				}
				if newQA.Record.ID >= f.qaID {
					t.Fatal("regression requires later native answer to have a lower random ID")
				}
				checked, err := CheckTaskIntegration(f.dependencies, input)
				if err == nil && checked.Readback.Classification == "ready" {
					t.Error("direct integration check accepted older pass after later native human failure")
				}
				applied, err := ApplyTaskIntegration(f.dependencies, input)
				if err == nil && applied.Readback.Classification == "exact_effect" {
					t.Error("saved integration confirmation reused older human pass")
				}
				if got := runLocalGit(t, input.DeliveryAuthorization.Agreement.TargetWorktree, "rev-parse", "HEAD"); got != f.oid {
					t.Fatal("later human failure did not stop the local target effect")
				}
			})
		}
	}
}

func TestDeliveryInterruptedIntentCheckRechecksLatestHumanAnswer(t *testing.T) {
	f := newIntegrationJourneyFixtureHistory(t, false, false, fixtureAgreement(DeliveryLocalBranch, "refs/heads/main"))
	f.dependencies.TaskLifecycleIDs = &failFirstAttemptIDs{TaskLifecycleIDSource: f.dependencies.TaskLifecycleIDs}
	input := deliveryIntegrationInput(f)
	before, err := CheckTaskIntegration(f.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, before.Readback)
	if _, err = ApplyTaskIntegration(f.dependencies, input); err == nil {
		t.Fatal("injected interruption before native attempt was not observed")
	}
	registry := mustQueueRegistry(t, f.workItemJourneyFixture)
	if len(registry.IntegrationAuthorities) != 1 || len(registry.IntegrationAttempts) != 0 {
		t.Fatal("fixture must stop with preserved intent before the Git effect")
	}
	qaFile, _ := writeHumanQADraft(t, taskDeliveryFixture{workItemJourneyFixture: f.workItemJourneyFixture, root: f.root}, registry.TaskResults[0], "fail", "report")
	if _, err = RecordTaskHumanQA(f.dependencies, TaskHumanQARecordInput{TaskID: "task", File: qaFile}); err != nil {
		t.Fatal(err)
	}
	checked, err := CheckTaskIntegration(f.dependencies, input)
	if err == nil && checked.Readback.Classification == "ready" {
		t.Fatal("persisted unattempted intent still reports ready after a newer human failure")
	}
	if got := runLocalGit(t, input.DeliveryAuthorization.Agreement.TargetWorktree, "rev-parse", "HEAD"); got != f.oid {
		t.Fatal("read-only check moved target")
	}
}
