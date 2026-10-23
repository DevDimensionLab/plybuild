package workspace

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParallelHistoricalHumanVetoIsNotReattributed(t *testing.T) {
	a := automaticAgreementFixture()
	r := automaticResultFixture()
	for _, outcome := range []string{"fail", "blocked"} {
		qa := TaskHumanQARecord{ID: "hqa_11111111111111111111111111111111", TaskID: r.TaskID, TaskResultID: "trs_22222222222222222222222222222222", ResultOID: strings.Repeat("d", 40), ResultTree: strings.Repeat("e", 40), Outcome: outcome, Actor: TaskHumanActorRecord{CompletedAtUTC: "2026-10-23T10:00:00Z"}}
		got, err := EvaluateDeliveryAcceptance(&a, r, &qa)
		if err != nil || got.Outcome != outcome || got.HumanQARecordID != qa.ID {
			t.Fatalf("historical human veto must survive changed parallel source without rebinding QA: %+v %v", got, err)
		}
		selected, err := LatestParallelTaskHumanQA([]TaskHumanQARecord{qa}, r)
		if err != nil || selected == nil || selected.ID != qa.ID {
			t.Fatal("reconciliation lost unresolved Task veto")
		}
		pass := qa
		pass.ID = "hqa_33333333333333333333333333333333"
		pass.TaskResultID = r.ID
		pass.ResultOID = r.ResultOID
		pass.ResultTree = r.ResultTree
		pass.Outcome = "pass"
		pass.Actor.CompletedAtUTC = "2026-10-23T11:00:00Z"
		selected, err = LatestParallelTaskHumanQA([]TaskHumanQARecord{qa, pass}, r)
		if err != nil || selected == nil || selected.ID != pass.ID {
			t.Fatal("later exact human judgment could not resolve veto")
		}
		if qa.Outcome != outcome || qa.TaskResultID == r.ID {
			t.Fatal("human record was reattributed")
		}
	}
}

func TestParallelEpicBasePlanNamesAllActiveMembers(t *testing.T) {
	f := parallelFixture(t)
	g := parallelGoalFixture(t, f, "task", "base-plan")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &g}})
	in := TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: g.SpecID, Parallel: true}
	p, err := PreviewTaskGoalExecution(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	prepared, err := PrepareTaskGoalExecution(f.dependencies, in, p.Confirmation, goalHumanFixture())
	if err != nil {
		t.Fatal(err)
	}
	parent := prepared.Preparation.Plan.Target.ParentLocator
	runLocalGit(t, parent, "commit", "--allow-empty", "-m", "external parent advance")
	out, err := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(out.Plan)
	if err != nil {
		t.Fatal(err)
	}
	var got struct {
		SchemaVersion int      `json:"schema_version"`
		Active        []string `json:"active_preparation_ids"`
	}
	if err = json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got.SchemaVersion != 2 || len(got.Active) != 1 || got.Active[0] != prepared.Preparation.Preparation.ID {
		t.Fatalf("format8 base plan omitted authoritative membership: %s", raw)
	}
	_, err = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *out.Confirmation})
	if err != nil {
		t.Fatal(err)
	}
}

func TestParallelEmptyActiveReadbackIsExplicitArray(t *testing.T) {
	f := parallelFixture(t)
	q, e := ListTaskQueueVersion(f.dependencies, queueTargetFixture(), false, 2)
	if e != nil {
		t.Fatal(e)
	}
	raw, e := MarshalTaskQueue(q)
	if e != nil {
		t.Fatal(e)
	}
	var m map[string]json.RawMessage
	if e = json.Unmarshal(raw, &m); e != nil {
		t.Fatal(e)
	}
	if string(m["active"]) != "[]" {
		t.Fatalf("schema2 empty membership is not explicit: %s", raw)
	}
	q, e = ListTaskQueue(f.dependencies, queueTargetFixture(), false)
	if e != nil {
		t.Fatal(e)
	}
	raw, e = MarshalTaskQueue(q)
	if e != nil {
		t.Fatal(e)
	}
	if strings.Contains(string(raw), `"active"`) {
		t.Fatal("legacy readback acquired parallel meaning")
	}
}
