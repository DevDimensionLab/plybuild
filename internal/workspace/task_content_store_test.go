package workspace

import (
	"errors"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestTaskContentRejectsManifestWhoseBytesDisagreeWithPreservedRequest(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	first, err := RecordTaskSpec(f.dependencies, contentFixtureSpec(t, f, "fixture/spec", nil))
	if err != nil {
		t.Fatal(err)
	}
	r, err := f.dependencies.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	m, err := readContentManifest(f.dependencies.TaskContent, root, first.OutcomeRef.ManifestSHA256)
	if err != nil {
		t.Fatal(err)
	}
	fields := contentFields(m)
	fields["title"] = "Unrequested replacement title"
	b, err := canonicaljson.Marshal(contentObject(fields))
	if err != nil {
		t.Fatal(err)
	}
	digest, err := f.dependencies.TaskContent.Publish(root, "manifests", b)
	if err != nil {
		t.Fatal(err)
	}
	r.TaskSpecRevisions[0].ManifestSHA256 = digest
	for i := range r.TaskContentPublications {
		if r.TaskContentPublications[i].OutcomeRef.Kind == "spec" {
			r.TaskContentPublications[i].OutcomeRef.ManifestSHA256 = digest
		}
	}
	b, err = encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(workItemsPath(root), b, 0644); err != nil {
		t.Fatal(err)
	}
	_, err = RecordTaskSpec(f.dependencies, contentFixtureSpec(t, f, "fixture/spec-2", &TaskRevisionRef{1, digest}))
	if err == nil {
		t.Fatal("accepted a manifest inconsistent with its immutable request")
	}
}

func TestTaskContentPublicationFaultsPreserveOutcomeAndExactRetry(t *testing.T) {
	for _, point := range []string{"object-temp-open", "object-directory-sync", "before-registry-replace", "postread"} {
		t.Run(point, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			root := filepath.Dir(f.wrapper)
			in := contentFixtureProblem(t, f, "fixture/problem-2")
			f.dependencies.TaskContent = &TaskContentStorage{fault: func(p string) error {
				if p == point {
					return errors.New("fixture publication fault")
				}
				return nil
			}}
			out, err := RecordTaskProblem(f.dependencies, in)
			if err == nil {
				t.Fatal("fault reported success")
			}
			if point == "postread" {
				if out.Classification != "unknown" || out.Attempted.RegistryReplace != "completed" || out.PostState.ProblemHead.Revision != 2 {
					t.Fatalf("lost possible publication: %#v", out)
				}
			} else if out.Attempted.RegistryReplace != "not_attempted" {
				t.Fatalf("unexpected registry attempt: %#v", out)
			}
			f.dependencies.TaskContent = &TaskContentStorage{}
			retry, err := RecordTaskProblem(f.dependencies, in)
			if err != nil {
				t.Fatal(err)
			}
			if retry.OutcomeRef == nil || *retry.OutcomeRef.Revision != 2 {
				t.Fatalf("retry allocated another revision: %#v", retry)
			}
			r, err := f.dependencies.WorkItems.Snapshot(root)
			if err != nil || len(r.TaskProblemRevisions) != 2 {
				t.Fatalf("registry after recovery: %#v %v", r, err)
			}
		})
	}
}

func TestTaskContentSourceHashAndSymlinkRejectedBeforePublication(t *testing.T) {
	for _, mode := range []string{"hash", "symlink", "short_read", "git_provenance"} {
		t.Run(mode, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			root := filepath.Dir(f.wrapper)
			in := contentFixtureSpec(t, f, "fixture/spec", nil)
			doc := filepath.Join(root, "solution.md")
			if mode == "hash" {
				if err := os.WriteFile(doc, []byte("changed"), 0600); err != nil {
					t.Fatal(err)
				}
			} else if mode == "symlink" {
				if err := os.Rename(doc, doc+".real"); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(doc+".real", doc); err != nil {
					t.Fatal(err)
				}
			} else if mode == "short_read" {
				reads := 0
				f.dependencies.TaskContent = &TaskContentStorage{fault: func(point string) error {
					if point == "source-read" {
						reads++
						if reads == 2 {
							return io.ErrUnexpectedEOF
						}
					}
					return nil
				}}
			} else {
				r, _ := f.dependencies.WorkItems.Snapshot(root)
				parent := r.Epics[0].RepoBindings[0].Worktree
				raw, _ := os.ReadFile(in.File)
				v, _ := canonicaljson.DecodeStrict(raw)
				m := contentFields(v)
				docs := contentArray(m, "documents")
				dm := contentFields(docs[0])
				src := contentFields(dm["source"])
				src["git_provenance"] = contentObject(map[string]canonicaljson.Value{"git_common_dir": r.Tasks[0].GitCommonDir, "ref": parent.Ref, "oid": parent.OID, "tree": parent.Tree, "repo_path": "not-the-solution.md", "blob": parent.OID})
				dm["source"] = contentObject(src)
				m["documents"] = []canonicaljson.Value{contentObject(dm)}
				raw, _ = canonicaljson.Marshal(contentObject(m))
				if err := os.WriteFile(in.File, raw, 0600); err != nil {
					t.Fatal(err)
				}
			}
			before, _ := os.ReadFile(workItemsPath(root))
			out, err := RecordTaskSpec(f.dependencies, in)
			after, _ := os.ReadFile(workItemsPath(root))
			if err == nil || out.Attempted.ObjectPublicationStarted || string(before) != string(after) {
				t.Fatalf("source drift published: %#v %v", out, err)
			}
		})
	}
}

func TestTaskContentRegistrySyncAndUnlockFailuresRemainUnknown(t *testing.T) {
	for _, point := range []string{"directory-sync", "unlock", "lock-close"} {
		t.Run(point, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			root := filepath.Dir(f.wrapper)
			in := contentFixtureProblem(t, f, "fixture/p2")
			f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(p string) error {
				if p == point {
					return errInjected
				}
				return nil
			}}}
			out, err := RecordTaskProblem(f.dependencies, in)
			if err == nil || out.Classification != "unknown" || out.Durability != "unknown" || out.PostState.ProblemHead == nil || out.PostState.ProblemHead.Revision != 2 {
				t.Fatalf("false success or lost poststate: %#v %v", out, err)
			}
			f.dependencies.WorkItems = newSystemWorkItemStore()
			before, _ := os.ReadFile(workItemsPath(root))
			retry, err := RecordTaskProblem(f.dependencies, in)
			after, _ := os.ReadFile(workItemsPath(root))
			if err != nil || retry.Classification != "existing" || string(before) != string(after) {
				t.Fatalf("recovery rewrote: %#v %v", retry, err)
			}
		})
	}
}

func TestTaskCreateRetryDoesNotRequireOptionalDocumentContent(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	r, err := f.dependencies.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	p, err := readContentManifest(f.dependencies.TaskContent, root, r.TaskProblemRevisions[0].ManifestSHA256)
	if err != nil {
		t.Fatal(err)
	}
	doc := contentFields(contentArray(contentFields(p), "documents")[0])
	if err = os.Remove(contentString(doc, "locator")); err != nil {
		t.Fatal(err)
	}
	in, _ := ParseTaskCreateInput("task", "Task", "Description", "epic", "ply", "ply")
	out, err := CreateTask(f.dependencies, in)
	if err != nil || out.Created || out.Content != nil {
		t.Fatalf("registration retry required optional document content: %#v %v", out, err)
	}
	if err = verifyManifestDocuments(f.dependencies.TaskContent, root, p); err == nil {
		t.Fatal("explicit document verification hid missing immutable content")
	}
}

func TestTaskContentGeneratedProblemIsBoundToOriginalSummary(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	r, err := f.dependencies.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	m, err := readContentManifest(f.dependencies.TaskContent, root, r.TaskProblemRevisions[0].ManifestSHA256)
	if err != nil {
		t.Fatal(err)
	}
	fields := contentFields(m)
	docs := contentArray(fields, "documents")
	doc := contentFields(docs[0])
	raw := []byte("# Unrequested replacement problem\n")
	digest, err := f.dependencies.TaskContent.Publish(root, "objects", raw)
	if err != nil {
		t.Fatal(err)
	}
	doc["locator"], doc["sha256"], doc["size_bytes"] = taskContentPath(root, "objects", digest), digest, int64(len(raw))
	fields["documents"] = []canonicaljson.Value{contentObject(doc)}
	raw, _ = canonicaljson.Marshal(contentObject(fields))
	digest, err = f.dependencies.TaskContent.Publish(root, "manifests", raw)
	if err != nil {
		t.Fatal(err)
	}
	r.TaskProblemRevisions[0].ManifestSHA256 = digest
	r.TaskContentPublications[0].OutcomeRef.ManifestSHA256 = digest
	raw, err = encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(workItemsPath(root), raw, 0644); err != nil {
		t.Fatal(err)
	}
	shown, err := ShowTaskProblem(f.dependencies, TaskContentQuery{TaskID: "task"})
	if err != nil {
		t.Fatal(err)
	}
	if contentString(contentFields(shown.Value), "integrity") != "conflict" {
		t.Fatal("generated document differs from original Task summary but was accepted")
	}
}
