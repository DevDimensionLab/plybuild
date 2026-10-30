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

func TestParallelSerialPrepareNextKeepsCurrent(t *testing.T) {
	for _, version := range []int{0, 1, 2} {
		name := "format4"
		if version > 0 {
			name = "legacy-" + string(rune('0'+version))
		}
		t.Run(name, func(t *testing.T) {
			f := newWorkItemJourneyFixtureVersion(t, version > 0)
			if version > 0 {
				root := filepath.Dir(f.wrapper)
				if version == 2 {
					if e := f.dependencies.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
						r, e := session.Snapshot()
						if e != nil {
							return e
						}
						upgradeRegistryToV2(&r)
						return session.Publish(r)
					}); e != nil {
						t.Fatal(e)
					}
				}
				before, e := os.ReadFile(workItemsPath(root))
				if e != nil {
					t.Fatal(e)
				}
				in := contentFixtureProblem(t, f, "serial-next-legacy")
				if _, e = RecordTaskProblem(f.dependencies, in); e == nil {
					t.Fatal("legacy import lost explicit upgrade acknowledgement")
				}
				raw, e := os.ReadFile(in.File)
				if e != nil {
					t.Fatal(e)
				}
				v, e := canonicaljson.DecodeStrict(raw)
				if e != nil {
					t.Fatal(e)
				}
				m := contentFields(v)
				m["registry_upgrade"] = contentRefValue(TaskRegistryUpgrade{FromVersion: version, RegistrySHA256: digestTaskBytes(before)})
				raw, e = canonicaljson.Marshal(contentObject(m))
				if e != nil {
					t.Fatal(e)
				}
				if e = os.WriteFile(in.File, raw, 0600); e != nil {
					t.Fatal(e)
				}
				if _, e = RecordTaskProblem(f.dependencies, in); e != nil {
					t.Fatal(e)
				}
				backup, e := f.dependencies.TaskContent.Read(root, "backups", digestTaskBytes(before))
				if e != nil || string(backup) != string(before) {
					t.Fatal("legacy raw backup changed")
				}
			}
			selection := queueSelect(t, f, "task", "serial-next", nil)
			for _, id := range []TaskID{"ply-consumer-agent-status", "ply-registered-cross-repo-plan"} {
				if _, e := CreateTask(f.dependencies, TaskCreateInput{TaskID: id, Title: string(id), Description: "Protected pending reference", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
					t.Fatal(e)
				}
			}
			queueSetFixture(t, f, "serial-next", []QueueEntry{{TaskID: "task", Selection: &selection}, {TaskID: "ply-consumer-agent-status"}, {TaskID: "ply-registered-cross-repo-plan"}})
			first := queuePrepareFixture(t, f, true)
			before := queueStateBytes(t, f)
			for _, apply := range []bool{false, true} {
				out, e := PrepareTask(f.dependencies, TaskPrepareInput{Target: queueTargetFixture(), Selector: QueueSelector{Kind: "next"}, Apply: apply, Confirmation: first.Preparation.PlanSHA256})
				if e != nil || out.Preparation == nil || out.Preparation.ID != first.Preparation.ID || queueStateBytes(t, f) != before {
					t.Fatalf("serial --next retry apply=%t changed current/pending: %+v %v", apply, out, e)
				}
			}
			q, e := ListTaskQueue(f.dependencies, queueTargetFixture(), false)
			if e != nil || q.Registry.FormatVersion != 4 || len(q.Pending) != 2 || q.Pending[0].TaskID != "ply-consumer-agent-status" || q.Pending[1].TaskID != "ply-registered-cross-repo-plan" {
				t.Fatalf("serial pending/format changed: %+v %v", q, e)
			}
		})
	}
}

func TestParallelLegacyPrepareRefusesKnownParallelConfirmation(t *testing.T) {
	f, p, _ := sourceReconcileFixture(t, false)
	task := p.Preparation.Plan.TaskID
	before := queueStateBytes(t, f)
	for _, apply := range []bool{false, true} {
		_, e := PrepareTask(f.dependencies, TaskPrepareInput{Target: queueTargetFixture(), Selector: QueueSelector{Kind: "task", TaskID: &task}, Apply: apply, Confirmation: p.Preparation.Preparation.PlanSHA256})
		if e == nil {
			t.Errorf("legacy prepare adopted known parallel confirmation apply=%t", apply)
		}
		if queueStateBytes(t, f) != before {
			t.Fatal("exclusive read/recovery changed parallel preparation")
		}
	}
}
