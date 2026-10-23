package workspace

import (
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Runtime-free tests: real isolated registry/Git; no OS/provider inspection.
func TestParallelRequiresExplicitFormatActivation(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	g := goalFixture(t, f, "task", "goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &g}})
	before := queueStateBytes(t, f)
	var in TaskGoalExecuteInput
	if err := json.Unmarshal([]byte(`{"target":{"project_id":"ply","repo_id":"ply","epic_id":"epic"},"spec_id":"solution","parallel":true}`), &in); err != nil {
		t.Fatal(err)
	}
	in.SpecID = g.SpecID
	_, err := PreviewTaskGoalExecution(f.dependencies, in)
	if err == nil || !strings.Contains(err.Error(), "parallel_activation_required") {
		t.Fatalf("parallel must require explicit format activation: %v", err)
	}
	if before != queueStateBytes(t, f) {
		t.Fatal("preview changed registry or Git")
	}
}

func TestParallelFormatEightDirectPredecessorPreservesSerialHistory(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	g := goalFixture(t, f, "task", "goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &g}})
	r, err := f.dependencies.WorkItems.Snapshot(filepath.Dir(f.wrapper))
	if err != nil {
		t.Fatal(err)
	}
	raw, err := encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	boundary, err := yaml.Marshal(struct {
		Boundary RegistryFormatBoundary `yaml:"parallel_boundary"`
	}{RegistryFormatBoundary{4, r.RawSHA256, "sha256:" + strings.Repeat("a", 64), "v0.0.0-test", f.oid, "sha256:" + strings.Repeat("b", 64)}})
	if err != nil {
		t.Fatal(err)
	}
	// Native encoding is ordered and strict. Preserve the predecessor bytes.
	raw = []byte(strings.Replace(string(raw), "format_version: 4", "format_version: 8", 1))
	raw = append(raw, boundary...)
	raw = append(raw, []byte("source_reconciliations: []\n")...)
	decoded, err := decodeWorkItemRegistry(raw)
	if err != nil {
		t.Fatalf("direct 4->8 reader: %v", err)
	}
	if !contentTypedEqual(r.TaskQueueEvents, decoded.TaskQueueEvents) || decoded.FormatBoundary != nil || decoded.AdoptionBoundary != nil || decoded.EpicCreateBoundary != nil {
		t.Fatal("reader fabricated intermediate history")
	}
}

func TestParallelParticipationNamedClosurePreservesOtherMember(t *testing.T) {
	q := emptyQueueFold()
	for i, e := range []struct{ kind, id string }{{"reserve_parallel", "a"}, {"reserve_parallel", "b"}, {"advance", "a"}} {
		applyQueueEvent(&q, TaskQueueEvent{Revision: i + 1, Kind: e.kind, Request: map[string]any{"preparation_id": e.id}})
	}
	// Inspect the native fold's public serialization so this regression executes
	// against the old implementation rather than failing to compile on new fields.
	b, err := json.Marshal(q)
	if err != nil {
		t.Fatal(err)
	}
	var v struct{ Active []string }
	if err = json.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	if len(v.Active) != 1 || v.Active[0] != "b" || q.Terminal["a"] != "advanced" {
		t.Fatalf("named closure lost B: %s", b)
	}
}

func parallelGoalFixture(t *testing.T, f workItemJourneyFixture, task TaskID, key string, agreement ...DeliveryAgreement) TaskGoalRef {
	t.Helper()
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
	m["spec_id"] = key
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
	if len(agreement) != 0 {
		m["delivery"] = contentRefValue(agreement[0])
	} else {
		m["delivery"] = contentRefValue(DeliveryAgreement{SchemaVersion: 1, Mode: DeliveryLocalEpic, ProjectID: "ply", RepoID: "ply", EpicID: "epic", TargetRef: "refs/heads/epic", TargetWorktree: filepath.Join(f.wrapper, "epic")})
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

// Only test seeding: explicit format boundary with simulated release identity.
// Live activation is exercised separately through its qualified native seam.
func parallelFixture(t *testing.T) workItemJourneyFixture {
	t.Helper()
	f := newWorkItemJourneyFixture(t)
	queueSetFixture(t, f, "seed-format-four", []QueueEntry{})
	r := mustQueueRegistry(t, f)
	r.ParallelBoundary = &RegistryFormatBoundary{4, r.RawSHA256, "sha256:" + strings.Repeat("a", 64), "v0.0.0-fixture", f.oid, "sha256:" + strings.Repeat("b", 64)}
	r.FormatVersion = ParallelRegistryFormat
	r.SourceReconciliations = []TaskSourceReconciliation{}
	raw, err := encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(workItemsPath(filepath.Dir(f.wrapper)), raw, 0600); err != nil {
		t.Fatal(err)
	}
	return f
}

func TestParallelNamedGoalPreparationsAndCAS(t *testing.T) {
	f := parallelFixture(t)
	if _, err := CreateTask(f.dependencies, TaskCreateInput{TaskID: "task-b", Title: "B", Description: "Independent B", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
		t.Fatal(err)
	}
	a := parallelGoalFixture(t, f, "task", "parallel-a")
	b := parallelGoalFixture(t, f, "task-b", "parallel-b")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &a}, {TaskID: "task-b", Goal: &b}})
	ai := TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: a.SpecID, Parallel: true}
	bi := TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: b.SpecID, Parallel: true}
	ap, err := PreviewTaskGoalExecution(f.dependencies, ai)
	if err != nil {
		t.Fatal(err)
	}
	bp, err := PreviewTaskGoalExecution(f.dependencies, bi)
	if err != nil {
		t.Fatal(err)
	}
	if ap.SchemaVersion != 2 {
		t.Fatalf("parallel plan did not freeze versioned authority: %#v", ap)
	}
	ar, err := PrepareTaskGoalExecution(f.dependencies, ai, ap.Confirmation, goalHumanFixture())
	if err != nil {
		t.Fatal(err)
	}
	before := queueStateBytes(t, f)
	if _, err = PrepareTaskGoalExecution(f.dependencies, bi, bp.Confirmation, goalHumanFixture()); err == nil {
		t.Fatal("unused stale B preview admitted")
	}
	if before != queueStateBytes(t, f) {
		t.Fatal("stale preview caused effects")
	}
	bp, err = PreviewTaskGoalExecution(f.dependencies, bi)
	if err != nil {
		t.Fatal(err)
	}
	br, err := PrepareTaskGoalExecution(f.dependencies, bi, bp.Confirmation, goalHumanFixture())
	if err != nil {
		t.Fatal(err)
	}
	if ar.Preparation.Preparation.ID == br.Preparation.Preparation.ID || ar.Preparation.Plan.WorktreePath == br.Preparation.Plan.WorktreePath || ar.Preparation.Plan.Branch == br.Preparation.Plan.Branch {
		t.Fatal("owners share source identity")
	}
	r := mustQueueRegistry(t, f)
	q := foldQueue(r, ap.QueueID)
	if len(q.Active) != 2 || q.Current != nil || len(q.Pending) != 0 {
		t.Fatalf("membership: %#v", q)
	}
	ws, err := ObserveRoot(f.dependencies, filepath.Dir(f.wrapper))
	if err != nil {
		t.Fatal(err)
	}
	projects, repos, err := f.dependencies.Projects.Snapshot(filepath.Dir(f.wrapper))
	if err != nil {
		t.Fatal(err)
	}
	registered := ReadRegisteredTaskQueues(f.dependencies, ws, ProjectSnapshot{Projects: projects, Repos: repos}, r)
	raw, err := json.Marshal(registered)
	if err != nil {
		t.Fatal(err)
	}
	var public struct {
		Queues []struct {
			Active []struct {
				TaskID TaskID `json:"task_id"`
			}
		}
	}
	if err = json.Unmarshal(raw, &public); err != nil {
		t.Fatal(err)
	}
	if len(public.Queues) != 1 || len(public.Queues[0].Active) != 2 {
		t.Fatalf("registered projection hid parallel members: %s", raw)
	}

	again, err := PrepareTaskGoalExecution(f.dependencies, ai, ap.Confirmation, goalHumanFixture())
	if err != nil || again.Preparation.Preparation.ID != ar.Preparation.Preparation.ID {
		t.Fatalf("A retry lost identity: %v", err)
	}
	if len(mustQueueRegistry(t, f).TaskPreparations) != 2 {
		t.Fatal("retry duplicated preparation")
	}
}
