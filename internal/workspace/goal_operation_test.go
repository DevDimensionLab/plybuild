package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoalOperationIgnoresUnrelatedHistory(t *testing.T) {
	for _, mode := range []string{"missing", "corrupt", "next-missing", "next-corrupt"} {
		t.Run(mode, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			in, _ := ParseTaskCreateInput("history", "history", "Unrelated history", "epic", "ply", "ply")
			if _, err := CreateTask(f.dependencies, in); err != nil {
				t.Fatal(err)
			}
			old := goalFixture(t, f, "history", "historical-goal")
			goal := goalFixture(t, f, "task", "selected-goal")
			goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
			root := filepath.Dir(f.wrapper)
			path := taskContentPath(root, "manifests", old.Spec.ManifestSHA256)
			if strings.HasSuffix(mode, "missing") {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := os.Chmod(path, 0600); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte("corrupt"), 0400); err != nil {
					t.Fatal(err)
				}
			}
			input := TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: goal.SpecID}
			if strings.HasPrefix(mode, "next-") {
				input.SpecID = ""
				input.Next = true
			}
			p, err := PreviewTaskGoalExecution(f.dependencies, input)
			if err != nil {
				t.Fatalf("unrelated history blocked preview: %v", err)
			}
			prepared, err := PrepareTaskGoalExecution(f.dependencies, p.Input, p.Confirmation, goalHumanFixture())
			if err != nil {
				t.Fatalf("unrelated history blocked preparation: %v", err)
			}
			if prepared.Preparation.State != "prepared" {
				t.Fatalf("not prepared: %#v", prepared)
			}
			if _, err := f.dependencies.WorkItems.Snapshot(root); err == nil {
				t.Fatal("full history validation accepted unavailable history")
			}
			r, err := f.dependencies.WorkItems.SnapshotRegistrations(root)
			if err != nil {
				t.Fatal(err)
			}
			if task, _ := findTask(r, "history"); task == nil {
				t.Fatal("mutation lost unrelated Task")
			}
		})
	}
}

func TestGoalOperationRelevantChangesBlockBeforeEffect(t *testing.T) {
	for _, change := range []string{"content", "structure", "queue"} {
		t.Run(change, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			goal := goalFixture(t, f, "task", "selected-goal")
			goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
			p, err := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: goal.SpecID})
			if err != nil {
				t.Fatal(err)
			}
			root := filepath.Dir(f.wrapper)
			switch change {
			case "content":
				if err = os.Remove(taskContentPath(root, "manifests", goal.Spec.ManifestSHA256)); err != nil {
					t.Fatal(err)
				}
			case "structure":
				b, e := os.ReadFile(workItemsPath(root))
				if e != nil {
					t.Fatal(e)
				}
				if err = os.WriteFile(workItemsPath(root), append(b, []byte("unknown_authority: true\n")...), 0600); err != nil {
					t.Fatal(err)
				}
			case "queue":
				goalQueueFixture(t, f, []TaskGoalQueueEntry{})
			}
			if _, err = PrepareTaskGoalExecution(f.dependencies, p.Input, p.Confirmation, goalHumanFixture()); err == nil {
				t.Fatal("relevant change accepted")
			}
			if _, err = os.Stat(p.WorktreePath); !os.IsNotExist(err) {
				t.Fatal("effect started after relevant drift")
			}
		})
	}
}

func TestStartupDependencyDigestKeepsSharedFacts(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	p, _ := goalPreparedFixture(t, f)
	r := mustQueueRegistry(t, f)
	before := TaskStartupDependencySHA256(r, "task")
	sibling := TaskContentPublication{TaskID: "other", PublicationKey: "unrelated"}
	r.TaskContentPublications = append(r.TaskContentPublications, sibling)
	if got := TaskStartupDependencySHA256(r, "task"); got != before {
		t.Fatal("unrelated external publication changed dependency digest")
	}
	for _, change := range []string{"selected-content", "ownership", "queue", "base"} {
		t.Run(change, func(t *testing.T) {
			next := r
			switch change {
			case "selected-content":
				next.TaskContentPublications = append(append([]TaskContentPublication{}, r.TaskContentPublications...), TaskContentPublication{TaskID: "task", PublicationKey: "changed"})
			case "ownership":
				next.Tasks = append(append([]TaskRecord{}, r.Tasks...), TaskRecord{ID: "conflicting-owner"})
			case "queue":
				next.TaskQueueEvents = append(append([]TaskQueueEvent{}, r.TaskQueueEvents...), TaskQueueEvent{QueueID: p.QueueID, Revision: 100})
			case "base":
				next.EpicBaseVersions = append(append([]EpicBaseVersion{}, r.EpicBaseVersions...), EpicBaseVersion{})
			}
			if TaskStartupDependencySHA256(next, "task") == before {
				t.Fatal("relevant shared facts omitted")
			}
		})
	}
}

func TestGoalOperationUnboundSessionCannotPublish(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goalQueueFixture(t, f, []TaskGoalQueueEntry{})
	root := filepath.Dir(f.wrapper)
	before, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	d, _ := beginGoalOperationContent(f.dependencies)
	err = d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
		r, err := session.Snapshot()
		if err != nil {
			return err
		}
		return session.Publish(r)
	})
	if err == nil || !strings.Contains(err.Error(), "dependencies are not bound") {
		t.Fatalf("unbound publication accepted: %v", err)
	}
	after, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("unbound operation changed shared registry")
	}
}
