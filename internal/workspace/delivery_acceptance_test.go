package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func automaticAgreementFixture() DeliveryAgreement {
	a := fixtureAgreement(DeliveryLocalEpic, "refs/heads/epic")
	a.SchemaVersion = 3
	a.TargetWorktree = "/workspace/epic"
	a.Acceptance = &DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "fixture orchestrator", TimeoutSeconds: 60}
	return a
}

func automaticResultFixture() TaskResultRecord {
	oid := strings.Repeat("a", 40)
	return TaskResultRecord{ID: "trs_11111111111111111111111111111111", TaskID: "task", ResultOID: oid, ResultTree: strings.Repeat("b", 40), TerminalResultSHA256: "sha256:" + strings.Repeat("c", 64), TechnicalGate: "passed", ReportedOutcome: "complete", VerifierResults: []TaskVerifierResultRecord{{VerifierID: "acceptance", Outcome: "passed", Argv: []string{"/bin/sh", "/workspace/acceptance.sh"}, CWD: "/workspace/task", BoundOIDOrSHA256: oid, Exit: 0}}}
}

func TestAutomaticAcceptanceIsExplicitAndLocalEpicOnly(t *testing.T) {
	a := automaticAgreementFixture()
	if err := ValidateDeliveryAgreement(a); err != nil || !a.AutomaticAcceptance() {
		t.Fatalf("valid automatic policy: %v", err)
	}
	for _, mutate := range []func(*DeliveryAgreement){
		func(v *DeliveryAgreement) { v.SchemaVersion = 1 },
		func(v *DeliveryAgreement) { v.SchemaVersion = 2; v.IntegrationOwner = IntegrationOwnerHuman },
		func(v *DeliveryAgreement) { v.Acceptance = nil },
		func(v *DeliveryAgreement) { v.Mode = DeliveryLocalBranch },
		func(v *DeliveryAgreement) {
			v.Mode = DeliveryPullRequest
			v.TargetWorktree = ""
			v.Remote = "origin"
			v.GitHubRepository = "owner/repo"
		},
		func(v *DeliveryAgreement) { v.TargetRef = "refs/heads/main" },
		func(v *DeliveryAgreement) { v.TargetRef = "refs/heads/master" },
		func(v *DeliveryAgreement) { v.IntegrationOwner = "agent" },
		func(v *DeliveryAgreement) { v.Acceptance.ResponsibleActor = "" },
		func(v *DeliveryAgreement) { v.Acceptance.Mode = "human" },
		func(v *DeliveryAgreement) { v.Acceptance.TimeoutSeconds = 0 },
		func(v *DeliveryAgreement) { v.Acceptance.TimeoutSeconds = 86401 },
	} {
		bad := automaticAgreementFixture()
		mutate(&bad)
		if ValidateDeliveryAgreement(bad) == nil {
			t.Fatalf("invalid automatic agreement accepted: %+v", bad)
		}
	}
	for _, version := range []int{1, 2} {
		old := a
		old.SchemaVersion, old.Acceptance = version, nil
		if version == 2 {
			old.IntegrationOwner = IntegrationOwnerHuman
		}
		if _, present := contentFields(contentRefValue(old))["acceptance"]; present {
			t.Fatal("historical agreement acquired an acceptance field")
		}
		decision, err := EvaluateDeliveryAcceptance(&old, automaticResultFixture(), nil)
		if err != nil || decision.Outcome != "blocked" || decision.Mode != "human" {
			t.Fatalf("historical contract became automatic: %+v %v", decision, err)
		}
	}
	a.IntegrationOwner = IntegrationOwnerHuman
	if err := ValidateDeliveryAgreement(a); err != nil || !a.HumanOwnedIntegration() || len(a.AllowedEffects()) != 0 || a.StopAfter() != "qualified_candidate_before_human_integration" {
		t.Fatalf("automatic verification expanded human integration ownership: %+v %v", a, err)
	}
}

func TestAutomaticAcceptanceRequiresExactVerificationReviewAndHumanOverrides(t *testing.T) {
	a := automaticAgreementFixture()
	r := automaticResultFixture()
	decision, err := EvaluateDeliveryAcceptance(&a, r, nil)
	if err != nil || decision.Outcome != "pass" || decision.Mode != "automatic" || decision.TaskResultID != r.ID || decision.EvidenceSHA256 != r.TerminalResultSHA256 || decision.HumanQARecordID != "" || decision.PolicySHA256 == "" {
		t.Fatalf("automatic result lost its exact evidence: %+v %v", decision, err)
	}
	for name, mutate := range map[string]func(*TaskResultRecord){
		"no verifier":         func(v *TaskResultRecord) { v.VerifierResults = nil },
		"failed test":         func(v *TaskResultRecord) { v.VerifierResults[0].Outcome = "failed"; v.VerifierResults[0].Exit = 1 },
		"different candidate": func(v *TaskResultRecord) { v.ResultOID = strings.Repeat("d", 40) },
		"missing invocation":  func(v *TaskResultRecord) { v.VerifierResults[0].Argv = nil },
		"missing cwd":         func(v *TaskResultRecord) { v.VerifierResults[0].CWD = "" },
		"unqualified":         func(v *TaskResultRecord) { v.TechnicalGate = "blocked" },
		"missing terminal":    func(v *TaskResultRecord) { v.TerminalResultSHA256 = "" },
		"review finding":      func(v *TaskResultRecord) { v.Review.OpenActionableFindings = []TaskReviewEntry{{ID: "bug"}} },
	} {
		t.Run(name, func(t *testing.T) {
			changed := automaticResultFixture()
			mutate(&changed)
			got, err := EvaluateDeliveryAcceptance(&a, changed, nil)
			if err != nil || got.Outcome == "pass" {
				t.Fatalf("invalid candidate accepted: %+v %v", got, err)
			}
		})
	}
	qa := TaskHumanQARecord{ID: "hqa_11111111111111111111111111111111", TaskID: r.TaskID, TaskResultID: r.ID, ResultOID: r.ResultOID, ResultTree: r.ResultTree, Outcome: "pass", Actor: TaskHumanActorRecord{CompletedAtUTC: "2026-10-09T10:00:00Z"}}
	a.Acceptance.RequireHumanQA = true
	if got, _ := EvaluateDeliveryAcceptance(&a, r, nil); got.Outcome != "blocked" {
		t.Fatal("required human gate was omitted")
	}
	if got, err := EvaluateDeliveryAcceptance(&a, r, &qa); err != nil || got.Outcome != "pass" {
		t.Fatalf("explicit human gate failed: %+v %v", got, err)
	}
	a.Acceptance.RequireHumanQA = false
	for _, outcome := range []string{"fail", "blocked"} {
		later := qa
		later.ID, later.Outcome, later.Actor.CompletedAtUTC = "hqa_22222222222222222222222222222222", outcome, "2026-10-09T10:01:00Z"
		latest, err := LatestTaskHumanQA([]TaskHumanQARecord{later, qa}, r.ID, r.ResultOID, r.ResultTree)
		if err != nil {
			t.Fatal(err)
		}
		got, err := EvaluateDeliveryAcceptance(&a, r, latest)
		if err != nil || got.Outcome != outcome || got.HumanQARecordID != later.ID {
			t.Fatalf("later human answer ignored: %+v %v", got, err)
		}
	}
}

func newAutomaticIntegrationFixture(t *testing.T, keepHuman bool) (integrationJourneyFixture, TaskIntegrationInput) {
	t.Helper()
	f := newIntegrationJourneyFixtureHistory(t, false, false, automaticAgreementFixture())
	var verifier TaskVerifierResultRecord
	if err := f.dependencies.WorkItems.WithLock(f.root, func(session WorkItemStoreSession) error {
		r, err := session.Snapshot()
		if err != nil {
			return err
		}
		tr := findTaskResult(r, f.resultID)
		verifier = TaskVerifierResultRecord{VerifierID: "acceptance", Outcome: "passed", Argv: []string{"/bin/sh", "/private/acceptance.sh"}, CWD: f.taskPath, Exit: 0, BoundOIDOrSHA256: tr.ResultOID}
		for i := range r.TaskResults {
			if r.TaskResults[i].ID == f.resultID {
				r.TaskResults[i].VerifierResults = []TaskVerifierResultRecord{verifier}
			}
		}
		if !keepHuman {
			r.HumanQARecords = []TaskHumanQARecord{}
		}
		return session.Publish(r)
	}); err != nil {
		t.Fatal(err)
	}
	f.evidence.evidence.VerifierResults = []TaskVerifierResultRecord{verifier}
	f.evidence.evidence.ExpectedVerifierIDs = []string{verifier.VerifierID}
	auth := f.evidence.evidence.DeliveryAuthorization
	in := deliveryIntegrationInput(f)
	in.HumanQARecordID = ""
	in.DeliveryOwner = &DeliveryIntegrationOwner{RunID: auth.RunID, RequestSHA256: auth.RequestSHA256, ActorClaim: "fixture orchestrator", PreparationID: "fixture-preparation"}
	return f, in
}

func TestAutomaticNativeIntegrationRecoversAfterGitWithoutHumanQA(t *testing.T) {
	f, in := newAutomaticIntegrationFixture(t, false)
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	check, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil || check.Readback.Classification != "ready" || check.Readback.HumanQA != nil || check.Readback.Acceptance == nil || check.Readback.Acceptance.Outcome != "pass" {
		t.Fatalf("automatic candidate was not ready without human QA: %+v %v", check.Readback, err)
	}
	in.Apply, in.Confirmation = true, integrationDigest(t, check.Readback)
	publishes := 0
	f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(operation string) error {
		if operation == "temp-open" {
			publishes++
			if publishes == 3 {
				return errInjected
			}
		}
		return nil
	}}}
	first, err := ApplyTaskIntegration(f.dependencies, in)
	if err == nil || git.count() != 1 || first.Readback.Classification != "exact_effect" {
		t.Fatalf("missing actual post-Git interruption: %+v %v merges=%d", first.Readback, err, git.count())
	}
	f.dependencies.WorkItems = newSystemWorkItemStore()
	for i := 0; i < 2; i++ {
		recovered, err := CheckTaskIntegration(f.dependencies, in)
		if err == nil && recovered.Readback.RecoveryStatus != "complete" {
			recovered, err = ApplyTaskIntegration(f.dependencies, in)
		}
		if err != nil || recovered.Readback.RecoveryStatus != "complete" || recovered.Readback.ParentOID != f.resultOID || recovered.Readback.HumanQA != nil || git.count() != 1 {
			t.Fatalf("recovery repeated Git or lost automatic evidence: recovery=%s parent=%s human=%v error=%v merges=%d", recovered.Readback.RecoveryStatus, recovered.Readback.ParentOID, recovered.Readback.HumanQA, err, git.count())
		}
	}
	r, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil || len(r.HumanQARecords) != 0 || len(r.IntegrationAuthorities) != 1 || len(r.IntegrationAttempts) != 1 || len(r.IntegrationResults) != 1 {
		t.Fatalf("automatic recovery wrote unexpected history: %+v %v", r, err)
	}
	if r.IntegrationAuthorities[0].Mode != "delivery_owner_after_automatic_pass" || r.IntegrationAuthorities[0].Plan.HumanQAReady || r.IntegrationAuthorities[0].Plan.Acceptance == nil {
		t.Fatal("automatic acceptance was presented as human QA")
	}
	facts := RecordedTaskProgressFacts(r)["task"]
	if facts.Ambiguous || facts.IntegrationResult == nil || facts.HumanQA != nil {
		t.Fatalf("native Task progress lost machine integration: %+v", facts)
	}
	show, err := ShowTask(f.dependencies, "task")
	if err != nil || show.Integration == nil || show.Integration.Classification != "exact_effect" || show.Integration.HumanQA != nil || show.Integration.Acceptance == nil {
		t.Fatalf("Task status lost automatic recovery: %+v %v", show.Integration, err)
	}
}

func TestAutomaticNativeIntegrationRechecksLatestHumanVetoAndAuthority(t *testing.T) {
	for _, defect := range []string{"later fail", "missing native permission", "human-owned integration"} {
		t.Run(defect, func(t *testing.T) {
			f, in := newAutomaticIntegrationFixture(t, true)
			git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
			f.dependencies.IntegrationGit = git
			check, err := CheckTaskIntegration(f.dependencies, in)
			if err != nil {
				t.Fatal(err)
			}
			in.Apply, in.Confirmation = true, integrationDigest(t, check.Readback)
			switch defect {
			case "later fail":
				err = f.dependencies.WorkItems.WithLock(f.root, func(session WorkItemStoreSession) error {
					r, e := session.Snapshot()
					if e != nil {
						return e
					}
					q := r.HumanQARecords[0]
					q.ID, q.PublicationKey, q.Outcome = "hqa_22222222222222222222222222222222", "later/fail", "fail"
					q.Actor.CompletedAtUTC = "2026-09-28T12:03:00Z"
					r.HumanQARecords = append(r.HumanQARecords, q)
					sortWorkRegistry(&r)
					return session.Publish(r)
				})
			case "missing native permission":
				f.evidence.evidence.DeliveryAuthorization = nil
			case "human-owned integration":
				in.DeliveryAuthorization.HumanIntegrationRequired = true
			}
			if err != nil {
				t.Fatal(err)
			}
			if _, err = ApplyTaskIntegration(f.dependencies, in); err == nil || git.count() != 0 {
				t.Fatalf("vetoed automatic integration changed Git: %v merges=%d", err, git.count())
			}
		})
	}
}

func TestAutomaticNativeRetryPreservesFailuresAndRequiresFreshConfirmation(t *testing.T) {
	f, in := newAutomaticIntegrationFixture(t, false)
	base := f.dependencies.IntegrationGit
	git := &noEffectIntegrationGit{TaskIntegrationGit: base}
	f.dependencies.IntegrationGit = git
	var previous IntegrationResultID
	for i := 0; i < 2; i++ {
		check, err := CheckTaskIntegration(f.dependencies, in)
		if err != nil || check.Readback.Classification != "ready" {
			t.Fatalf("retry %d did not recheck candidate: %s %v", i, check.Readback.Classification, err)
		}
		in.Apply, in.Confirmation = true, integrationDigest(t, check.Readback)
		out, err := ApplyTaskIntegration(f.dependencies, in)
		if err != nil || out.Readback.Classification != "no_effect" || git.merges != i+1 {
			t.Fatalf("retry %d did not preserve no-effect attempt: %s %v merges=%d", i, out.Readback.Classification, err, git.merges)
		}
		// The preceding confirmation cannot authorize another attempt after its
		// result is known. A new plan must name that exact no-effect result.
		if _, err = ApplyTaskIntegration(f.dependencies, in); err == nil || git.merges != i+1 {
			t.Fatalf("old confirmation repeated a failed attempt: %v merges=%d", err, git.merges)
		}
		r, err := f.dependencies.WorkItems.Snapshot(f.root)
		if err != nil || len(r.IntegrationResults) != i+1 || len(r.IntegrationAttempts) != i+1 || len(r.IntegrationAuthorities) != i+1 {
			t.Fatalf("retry %d lost its durable history: %v", i, err)
		}
		var latest *IntegrationResult
		for j := range r.IntegrationResults {
			result := &r.IntegrationResults[j]
			if result.RetryAfterResultID == nil && previous == "" || result.RetryAfterResultID != nil && *result.RetryAfterResultID == previous {
				latest = result
			}
		}
		if latest == nil {
			t.Fatal("retry did not name its exact prior result")
		}
		previous = latest.ID
	}
	f.dependencies.IntegrationGit = base
	check, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	in.Confirmation = integrationDigest(t, check.Readback)
	if out, err := ApplyTaskIntegration(f.dependencies, in); err != nil || out.Readback.RecoveryStatus != "complete" {
		t.Fatalf("resolved retry did not integrate: %s %v", out.Readback.Classification, err)
	}
	r, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil || len(r.IntegrationResults) != 3 || len(r.IntegrationAttempts) != 3 || len(r.HumanQARecords) != 0 {
		t.Fatalf("resolved retry lost prior failures or fabricated QA: %v", err)
	}
}

func TestAutomaticNativeRetryNeverRepeatsUncertainEffect(t *testing.T) {
	f, in := newAutomaticIntegrationFixture(t, false)
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = &repositoryFailureAfterCommandGit{TaskIntegrationGit: git}
	check, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	in.Apply, in.Confirmation = true, integrationDigest(t, check.Readback)
	out, err := ApplyTaskIntegration(f.dependencies, in)
	if err != nil || out.Readback.Classification != "unknown" || git.count() != 1 {
		t.Fatalf("did not preserve actual uncertainty: %s %v merges=%d", out.Readback.Classification, err, git.count())
	}
	f.dependencies.IntegrationGit = git
	check, err = CheckTaskIntegration(f.dependencies, in)
	if err != nil || check.Readback.Classification != "unknown" {
		t.Fatalf("unknown effect became a new retry: %s %v", check.Readback.Classification, err)
	}
	if out, err = ApplyTaskIntegration(f.dependencies, in); err != nil || out.Readback.Classification != "unknown" || git.count() != 1 {
		t.Fatalf("unknown effect was attempted again: %s %v merges=%d", out.Readback.Classification, err, git.count())
	}
	r, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil || len(r.IntegrationResults) != 1 || len(r.IntegrationAttempts) != 1 || len(r.IntegrationAuthorities) != 1 {
		t.Fatalf("unknown effect gained a competing attempt: %v", err)
	}
}

func TestAutomaticNativeGatePreservesHumanVetoAcrossUnchangedResultIDs(t *testing.T) {
	f, in := newAutomaticIntegrationFixture(t, true)
	preserveSameByteHumanJudgment(t, f, "fail")
	checked, err := CheckTaskIntegration(f.dependencies, in)
	if err == nil && checked.Readback.Classification == "ready" {
		t.Fatal("new TaskResult ID erased an actual human fail for unchanged candidate bytes")
	}
}

func TestAutomaticNativeGateNeverRebindsPriorPositiveHumanRecord(t *testing.T) {
	f, in := newAutomaticIntegrationFixture(t, true)
	preserveSameByteHumanJudgment(t, f, "pass")
	in.HumanQARecordID = f.qaID
	checked, err := CheckTaskIntegration(f.dependencies, in)
	if err == nil && checked.Readback.Classification == "ready" {
		t.Fatal("new TaskResult ID selected a prior result's positive human record")
	}
}

func preserveSameByteHumanJudgment(t *testing.T, f integrationJourneyFixture, outcome string) {
	t.Helper()
	if err := f.dependencies.WorkItems.WithLock(f.root, func(s WorkItemStoreSession) error {
		r, err := s.Snapshot()
		if err != nil {
			return err
		}
		previous := *findTaskResult(r, f.resultID)
		previous.ID, previous.PublicationKey = "trs_22222222222222222222222222222222", "previous-result"
		previous.TerminalResultID = "res_22222222222222222222222222222222"
		r.TaskResults = append(r.TaskResults, previous)
		link := *taskResultSpecLink(r, f.resultID)
		link.TaskResultID = previous.ID
		r.TaskResultSpecBindings = append(r.TaskResultSpecBindings, link)
		for i := range r.HumanQARecords {
			r.HumanQARecords[i].TaskResultID, r.HumanQARecords[i].Outcome = previous.ID, outcome
		}
		sortWorkRegistry(&r)
		return s.Publish(r)
	}); err != nil {
		t.Fatal(err)
	}
}

func TestAutomaticNativeCloseoutRetainsSourceAndRecoversLifecycleReceipt(t *testing.T) {
	f, in := newAutomaticIntegrationFixture(t, false)
	bytes := []byte(`{"fixture":"explicit automatic owner release and retained native evidence"}`)
	sha := digestTaskBytes(bytes)
	for _, name := range []string{"handoff", "start", "terminal", "owner-release"} {
		if err := os.WriteFile(filepath.Join(f.root, name), bytes, 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := f.dependencies.WorkItems.WithLock(f.root, func(s WorkItemStoreSession) error {
		r, err := s.Snapshot()
		if err != nil {
			return err
		}
		tr := findTaskResult(r, f.resultID)
		tr.HandoffSHA256, tr.StartReceiptSHA256, tr.TerminalResultSHA256 = sha, sha, sha
		f.evidence.evidence.HandoffSHA256, f.evidence.evidence.StartReceiptSHA256, f.evidence.evidence.TerminalResultSHA256 = sha, sha, sha
		for i := range r.TaskResultSpecBindings {
			r.TaskResultSpecBindings[i].HandoffSHA256, r.TaskResultSpecBindings[i].StartReceiptSHA256 = sha, sha
		}
		return s.Publish(r)
	}); err != nil {
		t.Fatal(err)
	}
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	check, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	in.Apply, in.Confirmation = true, integrationDigest(t, check.Readback)
	if _, err = ApplyTaskIntegration(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	queue, err := ListTaskQueue(f.dependencies, queueTargetFixture(), false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = CloseTaskQueue(f.dependencies, queueTargetFixture(), queue.Current.PreparationID, queue.Revision, "fixture automatic return observed", "advance"); err != nil {
		t.Fatal(err)
	}
	base, err := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *base.Confirmation}); err != nil {
		t.Fatal(err)
	}
	owner := TaskCloseoutOwnership{TaskID: "task", TaskResultID: f.resultID, ResultOID: f.resultOID, ResultTree: f.resultTree, EvidenceLocator: filepath.Join(f.root, "owner-release"), EvidenceSHA256: sha}
	preview, err := PreviewTaskReconcile(f.dependencies, "task", owner, true)
	if err != nil || preview.Ready || !strings.Contains(strings.Join(preview.Reasons, ";"), "owner_active") {
		t.Fatalf("active owner was not preserved: ready=%t reasons=%v err=%v", preview.Ready, preview.Reasons, err)
	}
	owner.Released = true
	preview, err = PreviewTaskReconcile(f.dependencies, "task", owner, true)
	if err != nil || !preview.Ready || !preview.IntegrationObserved {
		t.Fatalf("released automatic closeout not ready: ready=%t reasons=%v err=%v", preview.Ready, preview.Reasons, err)
	}
	closeIn := TaskCloseoutInput{OperationID: "automatic-closeout-fixture", TaskID: "task", TaskResultID: f.resultID, ExpectedResultOID: f.resultOID, ExpectedResultTree: f.resultTree, TargetRef: "refs/heads/epic", Ownership: owner, Keep: true, RequireIntegration: true}
	_, err = applyTaskCloseoutGuarded(f.dependencies, closeIn, preview.PlanSHA256, func() error { return nil }, func(stage string) error {
		if stage == "after-lifecycle" {
			return errInjected
		}
		return nil
	})
	if err == nil {
		t.Fatal("lifecycle receipt interruption was not reached")
	}
	state, err := readWorkItemLifecycle(f.root)
	if err != nil || state.Task("task") != LifecycleCompleted {
		t.Fatalf("actual lifecycle completion was lost: %+v %v", state, err)
	}
	receipt, err := ApplyTaskCloseoutWithOwnershipGuard(f.dependencies, closeIn, preview.PlanSHA256, func() error { return nil })
	if err != nil || receipt.State != "complete" || receipt.ResourceState != "kept" || !receipt.LifecycleCompleted || receipt.WorktreeRemoved || receipt.BranchRemoved || receipt.Retention == nil || git.count() != 1 {
		t.Fatalf("closeout did not resume retained lifecycle: %+v %v merges=%d", receipt, err, git.count())
	}
	if _, err := os.Stat(f.taskPath); err != nil {
		t.Fatalf("automatic closeout removed source: %v", err)
	}
	if got := runLocalGit(t, f.taskPath, "rev-parse", "HEAD"); got != f.resultOID {
		t.Fatal("automatic closeout changed source")
	}
	r, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil || len(r.HumanQARecords) != 0 || len(r.IntegrationAttempts) != 1 {
		t.Fatalf("automatic closeout fabricated acceptance or repeated integration: %v", err)
	}
}
