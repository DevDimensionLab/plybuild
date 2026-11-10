package workspace

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"os"
	"path/filepath"
	"testing"
)

// historicalDelivery is a delivery member exactly as Goals stored it before DP01.
func goalFixture(t *testing.T, f workItemJourneyFixture, task TaskID, key string, historicalDelivery ...map[string]any) TaskGoalRef {
	t.Helper()
	ensureFixtureTaskProblem(t, f, task)
	in := contentFixtureSpec(t, f, key, nil)
	raw, e := os.ReadFile(in.File)
	if e != nil {
		t.Fatal(e)
	}
	v, e := canonicaljson.DecodeStrict(raw)
	if e != nil {
		t.Fatal(e)
	}
	m := contentFields(v)
	r := mustQueueRegistry(t, f)
	m["task_id"] = string(task)
	m["problem"] = contentRefValue(taskContentState(r, task).ProblemHead)
	m["expected_previous"] = contentRefValue(taskContentState(r, task).SpecHead)
	for _, k := range []string{"parts", "supporting", "removed_requirement_ids", "phases", "implementation_basis", "dependencies"} {
		delete(m, k)
	}
	first := contentString(contentFields(contentArray(m, "documents")[0]), "id")
	m["kind"] = "WorkspaceTaskSpecDraft@2"
	m["schema_version"] = int64(2)
	m["contract_kind"] = "goal"
	m["objective"] = "Make the chosen behavior observable without choosing the detailed implementation."
	m["requirements"] = []canonicaljson.Value{contentObject(map[string]canonicaljson.Value{"id": "behavior", "acceptance": "The meaningful acceptance test observes the chosen behavior."})}
	m["design"] = []canonicaljson.Value{contentObject(map[string]canonicaljson.Value{"document_id": first, "section": nil})}
	m["constraints"] = []canonicaljson.Value{"Preserve existing behavior outside this goal."}
	m["executor"] = contentObject(map[string]canonicaljson.Value{"provider": "codex", "model": "fixture-model", "effort": "high"})
	if len(historicalDelivery) != 0 {
		m["delivery"] = contentRefValue(historicalDelivery[0])
	}
	raw, e = canonicaljson.Marshal(contentObject(m))
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(in.File, raw, 0600); e != nil {
		t.Fatal(e)
	}
	in.TaskID = task
	out, e := RecordTaskSpec(f.dependencies, in)
	if e != nil {
		t.Fatal(e)
	}
	return TaskGoalRef{SpecID: *out.OutcomeRef.SpecID, Spec: TaskRevisionRef{Revision: *out.OutcomeRef.Revision, ManifestSHA256: out.OutcomeRef.ManifestSHA256}}
}

func goalQueueFixture(t *testing.T, f workItemJourneyFixture, entries []TaskGoalQueueEntry) WorkspaceTaskQueueReadback {
	t.Helper()
	root := filepath.Dir(f.wrapper)
	r := mustQueueRegistry(t, f)
	q := foldQueue(r, queueID(root, QueueTarget{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}))
	draft := WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: fmt.Sprintf("goals/queue/%d", q.Revision), ProjectID: "ply", RepoID: "ply", EpicID: "epic", ExpectedRevision: q.Revision, Entries: entries, Recorder: TaskGoalQueueRecorder{ActorClaim: "fixture planner", ControlSurface: "synthetic test", RecordedAtUTC: "2026-10-05T00:00:00Z"}, RegistryUpgrade: queueUpgrade(r)}
	raw, e := contentCanonical(draft)
	if e != nil {
		t.Fatal(e)
	}
	file := filepath.Join(root, "goal-queue.json")
	if e = os.WriteFile(file, raw, 0600); e != nil {
		t.Fatal(e)
	}
	out, e := SetTaskQueue(f.dependencies, file)
	if e != nil {
		t.Fatal(e)
	}
	return out
}
