package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTaskContentReadbackUsesNormativeIntegrityReasons(t *testing.T) {
	for _, tc := range []struct {
		failure error
		reason  string
	}{{errors.New("fixture reader unavailable"), "task_content_observation_unknown"}, {contentError("task_content_integrity_conflict", "fixture mismatch", nil), "task_content_integrity_conflict"}} {
		t.Run(tc.reason, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			f.dependencies.TaskContent = &TaskContentStorage{fault: func(point string) error {
				if point == "managed-read" {
					return tc.failure
				}
				return nil
			}}
			out, err := ShowTaskProblem(f.dependencies, TaskContentQuery{TaskID: "task"})
			if err != nil {
				t.Fatal(err)
			}
			b, _ := MarshalTaskContentReadback(out)
			if !strings.Contains(string(b), tc.reason) {
				t.Fatalf("missing normative reason: %s", b)
			}
		})
	}
}

func TestTaskContentShowPreservesResourceProblemAndHistoricalResults(t *testing.T) {
	f := newTaskDeliveryFixture(t, false)
	recorded, err := RecordTaskResult(f.dependencies, TaskResultRecordInput{TaskID: "task", File: f.draftPath})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f.workItemJourneyFixture, "fixture/p2")); err != nil {
		t.Fatal(err)
	}
	show, err := ShowTask(f.dependencies, "task")
	if err != nil {
		t.Fatal(err)
	}
	m := contentFields(show.Content)
	if contentString(m, "kind") != "WorkspaceTaskReadback@3" || m["resource"] == nil || m["integration"] == nil || m["ready_for_spec_handoff"] != false {
		t.Fatalf("lost resource: %#v", m)
	}
	results := contentArray(m, "results")
	if len(results) != 1 || contentString(contentFields(results[0]), "relevance") != "stale" || contentString(contentFields(results[0]), "basis_status") != "bound" {
		t.Fatalf("history rewritten: %#v", results)
	}
	text := TaskContentText(show.Content)
	if !strings.Contains(text, string(recorded.Record.ID)) || !strings.Contains(text, "revision 1") || !strings.Contains(text, "stale") {
		t.Fatalf("human readback hides result basis: %s", text)
	}
}

func TestTaskContentReadDoesNotCreateMissingLocksOrStore(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	before, _ := os.ReadFile(workItemsPath(root))
	lock := filepath.Join(root, MarkerDirectory, workItemsLock)
	if err := os.Remove(lock); err != nil {
		t.Fatal(err)
	}
	if _, err := ShowTaskProblem(f.dependencies, TaskContentQuery{TaskID: "task"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(lock); !os.IsNotExist(err) {
		t.Fatalf("read created a lock: %v", err)
	}
	after, _ := os.ReadFile(workItemsPath(root))
	if string(before) != string(after) {
		t.Fatal("read rewrote store")
	}
	if err := WithTaskSpecSnapshot(f.dependencies, root, func(*TaskSpecSession) error { return nil }); err == nil {
		t.Fatal("guard accepted missing work-item lock")
	}
}

func TestTaskContentMissingProblemDocumentClearsSummaryAndRequiresInspection(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	r, err := f.dependencies.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	head := taskContentState(r, "task").ProblemHead
	manifest, err := readRegisteredTaskManifest(f.dependencies.TaskContent, root, r, head.ManifestSHA256)
	if err != nil {
		t.Fatal(err)
	}
	doc := contentFields(contentArray(contentFields(manifest), "documents")[0])
	if err := os.Remove(taskContentPath(root, "objects", contentString(doc, "sha256"))); err != nil {
		t.Fatal(err)
	}
	show, err := ShowTask(f.dependencies, "task")
	if err != nil {
		t.Fatal(err)
	}
	m := contentFields(show.Content)
	problem := contentFields(m["problem"])
	if problem["head"] == nil || problem["title"] != nil || problem["summary"] != nil || problem["integrity"] != "missing" || contentString(contentFields(m["next_action"]), "kind") != "inspect_task" {
		t.Fatalf("unreadable problem was presented as usable content: %#v", m)
	}
}
