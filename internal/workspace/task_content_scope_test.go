package workspace

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
	r, err := f.dependencies.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	for i := range r.TaskContentPublications {
		p := &r.TaskContentPublications[i]
		if p.TaskID != task || p.OutcomeRef.Kind != "spec" {
			continue
		}
		request, err := f.dependencies.TaskContent.Read(root, "requests", p.RequestSHA256)
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
		requestDigest, err := f.dependencies.TaskContent.Publish(root, "requests", request)
		if err != nil {
			t.Fatal(err)
		}
		manifest, err := readContentManifest(f.dependencies.TaskContent, root, p.OutcomeRef.ManifestSHA256)
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
		digest, err := f.dependencies.TaskContent.Publish(root, "manifests", raw)
		if err != nil {
			t.Fatal(err)
		}
		for j := range r.TaskSpecRevisions {
			if r.TaskSpecRevisions[j].ManifestSHA256 == p.OutcomeRef.ManifestSHA256 {
				r.TaskSpecRevisions[j].ManifestSHA256 = digest
			}
		}
		p.RequestSHA256, p.IntentSHA256, p.OutcomeRef.ManifestSHA256 = requestDigest, requestDigest, digest
		raw, err = encodeWorkItemRegistry(r)
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(workItemsPath(root), raw, 0600); err != nil {
			t.Fatal(err)
		}
		return digest
	}
	t.Fatal("fixture publication missing")
	return ""
}

func TestTaskScopedContentRequiresItsArtifactsAndPreservesOutsideRecords(t *testing.T) {
	for _, failure := range []string{"missing_manifest", "missing_request", "tampered_document", "outside_write"} {
		t.Run(failure, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			if _, err := CreateTask(f.dependencies, TaskCreateInput{TaskID: "sibling", Title: "Sibling", Description: "Independent delivery", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
				t.Fatal(err)
			}
			goalQueueFixture(t, f, []TaskGoalQueueEntry{})
			root := filepath.Dir(f.wrapper)
			d := WithTaskContentScope(f.dependencies, "task")
			r, err := d.WorkItems.Snapshot(root)
			if err != nil {
				t.Fatal(err)
			}
			before, err := os.ReadFile(workItemsPath(root))
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
				err = os.Remove(taskContentPath(root, "manifests", p.OutcomeRef.ManifestSHA256))
			case "missing_request":
				err = os.Remove(taskContentPath(root, "requests", p.RequestSHA256))
			case "tampered_document":
				m, e := readContentManifest(d.TaskContent, root, p.OutcomeRef.ManifestSHA256)
				if e != nil {
					t.Fatal(e)
				}
				doc := contentFields(contentArray(contentFields(m), "documents")[0])
				path := taskContentPath(root, "objects", contentString(doc, "sha256"))
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
				err = d.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
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
				_, err = ShowTask(d, "task")
			}
			var contentErr *TaskContentError
			if !errors.As(err, &contentErr) || contentErr.Code != want {
				t.Fatalf("required failure %s: got %v, want %s", failure, err, want)
			}
			if failure != "outside_write" && (contentErr.Artifact == "" || contentErr.Operation == "" || contentErr.ReaderCapability == "" || contentErr.NextAction == "") {
				t.Fatalf("diagnostic lacks artifact/operation/capability/action: %#v", contentErr)
			}
			after, e := os.ReadFile(workItemsPath(root))
			if e != nil || string(before) != string(after) {
				t.Fatal("blocked operation changed registry", e)
			}
		})
	}
}

func TestTaskScopedConcurrentPublicationsPreserveBothWritersAndUnknownSibling(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	for _, id := range []TaskID{"sibling", "future"} {
		if _, err := CreateTask(f.dependencies, TaskCreateInput{TaskID: id, Title: string(id), Description: "Independent work", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
			t.Fatal(err)
		}
	}
	goalFixture(t, f, "future", "future/goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{})
	root := filepath.Dir(f.wrapper)
	r, err := f.dependencies.WorkItems.Snapshot(root)
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
	after, err := f.dependencies.WorkItems.SnapshotRegistrations(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []TaskID{"task", "sibling"} {
		if taskContentState(after, id).ProblemHead.Revision != 2 {
			t.Fatalf("lost %s publication", id)
		}
	}
	if got := taskContentState(after, "future").SpecHead; got == nil || got.ManifestSHA256 != future {
		t.Fatal("rewrote unknown sibling")
	}
	if _, err = f.dependencies.WorkItems.Snapshot(root); err == nil || !strings.Contains(err.Error(), "future_review_policy") {
		t.Fatalf("global audit must still reject unsupported semantics: %v", err)
	}
}

func TestTaskScopedContentPreservesUnrelatedNewerPublications(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	if _, err := CreateTask(f.dependencies, TaskCreateInput{TaskID: "sibling", Title: "Sibling", Description: "Independent delivery", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
		t.Fatal(err)
	}
	goalFixture(t, f, "sibling", "sibling/goal", fixtureAgreement(DeliveryPullRequest, "refs/heads/main"))
	goalQueueFixture(t, f, []TaskGoalQueueEntry{})
	next := contentFixtureProblem(t, f, "active/revision")
	digest := scopedFuturePublication(t, f, "sibling")
	root := filepath.Dir(f.wrapper)
	before, err := f.dependencies.WorkItems.SnapshotRegistrations(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = ShowTask(f.dependencies, "task"); err != nil {
		t.Fatalf("unrelated newer publication blocked active Task readback: %v", err)
	}
	listed, err := ListTaskSpecs(f.dependencies, TaskContentQuery{TaskID: "task"})
	if err != nil || contentString(contentFields(listed.Value), "freshness") != "fresh" {
		t.Fatalf("unrelated publication made Spec list stale: %#v %v", listed, err)
	}
	if _, err = RecordTaskProblem(f.dependencies, next); err != nil {
		t.Fatalf("unrelated newer publication blocked active Task write: %v", err)
	}
	after, err := f.dependencies.WorkItems.SnapshotRegistrations(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range before.TaskContentPublications {
		if p.TaskID == "sibling" && !contentTypedEqual(&p, taskPublication(after, p.PublicationKey)) {
			t.Fatal("active Task publication changed sibling publication")
		}
	}
	if _, err = f.dependencies.TaskContent.Read(root, "manifests", digest); err != nil {
		t.Fatal("sibling bytes changed", err)
	}
	_, err = ShowTask(f.dependencies, "sibling")
	var incompatible *TaskContentError
	if !errors.As(err, &incompatible) || incompatible.Code != "task_content_reader_incompatible" {
		t.Fatalf("required unknown semantics must be incompatible, got %v", err)
	}
}

func TestTaskStoredKnownDeliveryVersionMalformationIsInvalidContent(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	a := fixtureAgreement(DeliveryPullRequest, "refs/heads/main")
	a.SchemaVersion, a.IntegrationOwner = 2, IntegrationOwnerHuman
	goal := goalFixture(t, f, "task", "goal/v2-delivery", a)
	root := filepath.Dir(f.wrapper)
	raw, err := f.dependencies.TaskContent.Read(root, "manifests", goal.Spec.ManifestSHA256)
	if err != nil {
		t.Fatal(err)
	}
	v, err := canonicaljson.DecodeStrict(raw)
	if err != nil {
		t.Fatal(err)
	}
	m := contentFields(v)
	delivery := contentFields(m["delivery"])
	delivery["target_ref"] = "invalid"
	m["delivery"] = contentObject(delivery)
	raw, err = canonicaljson.Marshal(contentObject(m))
	if err != nil {
		t.Fatal(err)
	}
	_, err = decodeStoredTaskContent(raw, "spec_record", true, true)
	var contentErr *TaskContentError
	if !errors.As(err, &contentErr) || contentErr.Code != "task_content_invalid_input" {
		t.Fatalf("known delivery v2 malformation misclassified: %v", err)
	}
	if !strings.Contains(taskContentReaderCapability, "delivery@2") {
		t.Fatal("reader capability omitted supported delivery v2")
	}
}
