package taskwork

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

// Preserve actual canonical publication/request bytes with a newer semantic
// field. The current reader deliberately does not understand that field.
func scopedFuturePublication(t *testing.T, f workItemJourneyFixture, task TaskID) string {
	t.Helper()
	root := filepath.Dir(f.wrapper)
	r, err := f.dependencies.WorkItems.Snapshot(testWork(t, root))
	if err != nil {
		t.Fatal(err)
	}
	for i := range r.TaskContentPublications {
		p := &r.TaskContentPublications[i]
		if p.TaskID != task || p.OutcomeRef.Kind != "spec" {
			continue
		}
		request, err := f.dependencies.TaskContent.Read(testWork(t, root), "requests", p.RequestSHA256)
		if err != nil {
			t.Fatal(err)
		}
		v, err := canonicaljson.DecodeStrict(request)
		if err != nil {
			t.Fatal(err)
		}
		fields := contentFields(v)
		fields["future_review_policy"] = contentObject(map[string]canonicaljson.Value{"required": true})
		request, err = canonicaljson.Marshal(contentObject(fields))
		if err != nil {
			t.Fatal(err)
		}
		requestDigest, err := f.dependencies.TaskContent.Publish(testWork(t, root), "requests", request)
		if err != nil {
			t.Fatal(err)
		}
		manifest, err := readContentManifest(f.dependencies.TaskContent, testWork(t, root), p.OutcomeRef.ManifestSHA256)
		if err != nil {
			t.Fatal(err)
		}
		fields = contentFields(manifest)
		fields["future_review_policy"] = contentObject(map[string]canonicaljson.Value{"required": true})
		fields["source_draft_sha256"] = requestDigest
		raw, err := canonicaljson.Marshal(contentObject(fields))
		if err != nil {
			t.Fatal(err)
		}
		digest, err := f.dependencies.TaskContent.Publish(testWork(t, root), "manifests", raw)
		if err != nil {
			t.Fatal(err)
		}
		for j := range r.TaskSpecRevisions {
			if r.TaskSpecRevisions[j].ManifestSHA256 == p.OutcomeRef.ManifestSHA256 {
				r.TaskSpecRevisions[j].ManifestSHA256 = digest
			}
		}
		p.RequestSHA256, p.IntentSHA256, p.OutcomeRef.ManifestSHA256 = requestDigest, requestDigest, digest
		raw = storedRegistryBytes(t, root, r)
		if err = os.WriteFile(testWork(t, root).Registry(), raw, 0600); err != nil {
			t.Fatal(err)
		}
		return digest
	}
	t.Fatal("fixture publication missing")
	return ""
}

func TestTaskScopedContentRequiresItsArtifactsAndPreservesOutsideRecords(t *testing.T) {
	t.Parallel()
	for _, failure := range []string{"missing_manifest", "missing_request", "tampered_document", "outside_write"} {
		t.Run(failure, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			if _, err := CreateTask(f.dependencies, TaskCreateInput{TaskID: "sibling", Title: "Sibling", Description: "Independent delivery", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
				t.Fatal(err)
			}
			goalQueueFixture(t, f, []TaskGoalQueueEntry{})
			next := contentFixtureProblem(t, f, "task/scoped-revision")
			root := filepath.Dir(f.wrapper)
			d := WithTaskContentScope(f.dependencies, "task")
			r, err := d.WorkItems.Snapshot(testWork(t, root))
			if err != nil {
				t.Fatal(err)
			}
			before, err := os.ReadFile(testWork(t, root).Registry())
			if err != nil {
				t.Fatal(err)
			}
			p := taskPublication(r, "task-create/task")
			if p == nil {
				t.Fatal("initial publication missing")
			}
			want := "task_content_missing"
			switch failure {
			case "missing_manifest":
				err = os.Remove(taskContentPath(testWork(t, root), "manifests", p.OutcomeRef.ManifestSHA256))
			case "missing_request":
				err = os.Remove(taskContentPath(testWork(t, root), "requests", p.RequestSHA256))
			case "tampered_document":
				m, e := readContentManifest(d.TaskContent, testWork(t, root), p.OutcomeRef.ManifestSHA256)
				if e != nil {
					t.Fatal(e)
				}
				doc := contentFields(contentArray(contentFields(m), "documents")[0])
				path := taskContentPath(testWork(t, root), "objects", contentString(doc, "sha256"))
				if err = os.Chmod(path, 0600); err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(path, []byte("tampered"), 0400); err != nil {
					t.Fatal(err)
				}
				err = os.Chmod(path, 0400)
				want = "task_content_integrity_conflict"
			case "outside_write":
				want = "task_content_scope_conflict"
			}
			if err != nil {
				t.Fatal(err)
			}
			if failure == "outside_write" {
				err = d.WorkItems.WithLock(testWork(t, root), func(s WorkItemStoreSession) error {
					r, e := s.Snapshot()
					if e != nil {
						return e
					}
					for i := range r.Tasks {
						if r.Tasks[i].ID == "sibling" {
							r.Tasks[i].Description = "Out-of-scope replacement"
						}
					}
					return s.Publish(r)
				})
			} else {
				_, err = RecordTaskProblem(d, next)
			}
			var contentErr *TaskContentError
			if !errors.As(err, &contentErr) || contentErr.Code != want {
				t.Fatalf("required failure %s: got %v, want %s", failure, err, want)
			}
			if failure != "outside_write" && (contentErr.Artifact == "" || contentErr.Operation == "" || contentErr.ReaderCapability == "" || contentErr.NextAction == "") {
				t.Fatalf("diagnostic lacks artifact/operation/capability/action: %#v", contentErr)
			}
			after, e := os.ReadFile(testWork(t, root).Registry())
			if e != nil || string(before) != string(after) {
				t.Fatal("blocked operation changed registry", e)
			}
		})
	}
}

func TestTaskScopedConcurrentPublicationsPreserveBothWritersAndUnknownSibling(t *testing.T) {
	t.Parallel()
	f := newWorkItemJourneyFixture(t)
	for _, id := range []TaskID{"sibling", "future"} {
		if _, err := CreateTask(f.dependencies, TaskCreateInput{TaskID: id, Title: string(id), Description: "Independent work", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
			t.Fatal(err)
		}
	}
	goalFixture(t, f, "future", "future/goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{})
	root := filepath.Dir(f.wrapper)
	r, err := f.dependencies.WorkItems.Snapshot(testWork(t, root))
	if err != nil {
		t.Fatal(err)
	}
	inputs := []TaskContentInput{}
	for _, id := range []TaskID{"task", "sibling"} {
		in := contentFixtureProblem(t, f, string(id)+"/new-"+string(id))
		raw, err := os.ReadFile(in.File)
		if err != nil {
			t.Fatal(err)
		}
		v, err := canonicaljson.DecodeStrict(raw)
		if err != nil {
			t.Fatal(err)
		}
		m := contentFields(v)
		m["task_id"], m["expected_previous"] = string(id), contentRefValue(taskContentState(r, id).ProblemHead)
		raw, err = canonicaljson.Marshal(contentObject(m))
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(in.File, raw, 0600); err != nil {
			t.Fatal(err)
		}
		in.TaskID = id
		inputs = append(inputs, in)
	}
	future := scopedFuturePublication(t, f, "future")
	start := make(chan struct{})
	done := make(chan error, len(inputs))
	for _, in := range inputs {
		go func(in TaskContentInput) { <-start; _, err := RecordTaskProblem(f.dependencies, in); done <- err }(in)
	}
	close(start)
	for range inputs {
		if err := <-done; err != nil {
			t.Fatal(err)
		}
	}
	after, err := f.dependencies.WorkItems.SnapshotRegistrations(testWork(t, root))
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []TaskID{"task", "sibling"} {
		previous := 0
		if head := taskContentState(r, id).ProblemHead; head != nil {
			previous = head.Revision
		}
		if taskContentState(after, id).ProblemHead.Revision != previous+1 {
			t.Fatalf("lost %s publication", id)
		}
	}
	if got := taskContentState(after, "future").SpecHead; got == nil || got.ManifestSHA256 != future {
		t.Fatal("rewrote unknown sibling")
	}
	if _, err = f.dependencies.WorkItems.Snapshot(testWork(t, root)); err != nil {
		t.Fatalf("registry observation required unrelated document semantics: %v", err)
	}
	if _, err = readRegisteredTaskManifest(f.dependencies.TaskContent, testWork(t, root), after, future); err == nil || !strings.Contains(err.Error(), "future_review_policy") {
		t.Fatalf("explicit document read accepted unsupported semantics: %v", err)
	}
}

func TestQueueReadbackDoesNotInterpretAnotherEpicsContent(t *testing.T) {
	t.Parallel()
	f := newWorkItemJourneyFixture(t)
	goalFixture(t, f, "task", "task/goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{})
	r := mustQueueRegistry(t, f)
	parent := r.Epics[0].RepoBindings[0].Worktree
	other := EpicAdoptInput{EpicID: "other", Title: "Other", ProjectID: "ply", RepoID: "ply", Worktree: parent.Locator, Ref: parent.Ref, ExpectedOID: parent.OID}
	if _, err := AdoptEpic(f.dependencies, other); err != nil {
		t.Fatal(err)
	}
	scopedFuturePublication(t, f, "task")
	before, err := os.ReadFile(testWork(t, filepath.Dir(f.wrapper)).Registry())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ListTaskQueue(f.dependencies, QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "other"}, false); err != nil {
		t.Fatalf("unrelated historical content blocked queue: %v", err)
	}
	if _, err := ListTaskQueue(f.dependencies, queueTargetFixture(), false); err == nil {
		t.Fatal("selected target's unknown content was accepted")
	}
	after, err := os.ReadFile(testWork(t, filepath.Dir(f.wrapper)).Registry())
	if err != nil || string(before) != string(after) {
		t.Fatal("queue read changed registry", err)
	}
}
