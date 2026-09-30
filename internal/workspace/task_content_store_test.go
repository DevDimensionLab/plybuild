package workspace

import (
	"bytes"
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

func TestTaskContentMigrationPreservesRawBackupAndLegacyRecords(t *testing.T) {
	for _, version := range []int{1, 2} {
		t.Run(string(rune('0'+version)), func(t *testing.T) {
			f := newWorkItemJourneyFixtureVersion(t, true)
			root := filepath.Dir(f.wrapper)
			if version == 2 {
				if err := f.dependencies.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
					r, err := s.Snapshot()
					if err != nil {
						return err
					}
					upgradeRegistryToV2(&r)
					return s.Publish(r)
				}); err != nil {
					t.Fatal(err)
				}
			}
			before, _ := os.ReadFile(workItemsPath(root))
			in := contentFixtureProblem(t, f, "fixture/import")
			out, err := RecordTaskProblem(f.dependencies, in)
			if err == nil || out.Attempted.ObjectPublicationStarted {
				t.Fatal("migration did not require acknowledgement")
			}
			raw, _ := os.ReadFile(in.File)
			v, _ := canonicaljson.DecodeStrict(raw)
			m := contentFields(v)
			m["registry_upgrade"] = contentObject(map[string]canonicaljson.Value{"from_version": int64(version), "registry_sha256": digestTaskBytes(before)})
			raw, _ = canonicaljson.Marshal(contentObject(m))
			if err = os.WriteFile(in.File, raw, 0600); err != nil {
				t.Fatal(err)
			}
			if _, err = RecordTaskProblem(f.dependencies, in); err != nil {
				t.Fatal(err)
			}
			r, err := f.dependencies.WorkItems.Snapshot(root)
			if err != nil || r.FormatVersion != 3 || r.Tasks[0].Description != "Description" || len(r.Epics) != 1 || len(r.TaskProblemRevisions) != 1 {
				t.Fatalf("migration changed legacy records: %#v %v", r, err)
			}
			backup, err := f.dependencies.TaskContent.Read(root, "backups", digestTaskBytes(before))
			if err != nil || string(before) != string(backup) {
				t.Fatalf("raw backup differs: %v", err)
			}
		})
	}
}

func TestTaskCreateRetryRechecksPublishedContent(t *testing.T) {
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
	if err == nil || out.Content == nil || out.Content.Classification != "unknown" {
		t.Fatalf("create retry hid missing immutable content: %#v %v", out, err)
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

func TestTaskContentLegacyImportRequiresExactAcknowledgementAndPreservesSummary(t *testing.T) {
	for _, version := range []int{1, 2} {
		t.Run(string(rune('0'+version)), func(t *testing.T) {
			f := newWorkItemJourneyFixtureVersion(t, true)
			root := filepath.Dir(f.wrapper)
			if version == 2 {
				if err := f.dependencies.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
					r, e := s.Snapshot()
					if e != nil {
						return e
					}
					upgradeRegistryToV2(&r)
					return s.Publish(r)
				}); err != nil {
					t.Fatal(err)
				}
			}
			before, err := os.ReadFile(workItemsPath(root))
			if err != nil {
				t.Fatal(err)
			}
			r, err := f.dependencies.WorkItems.Snapshot(root)
			if err != nil {
				t.Fatal(err)
			}
			task := r.Tasks[0]
			summary, _ := contentCanonical(map[string]any{"task_id": task.ID, "title": task.Title, "description": task.Description})
			fields := map[string]canonicaljson.Value{"expected_previous": nil, "origin": contentObject(map[string]canonicaljson.Value{"kind": "legacy_summary_import", "legacy_summary_sha256": digestTaskBytes(summary)}), "title": task.Title, "summary": task.Description, "problem_document_id": "problem", "documents": []canonicaljson.Value{}, "sources": []canonicaljson.Value{}, "claims": []canonicaljson.Value{}, "deadline": nil, "change_reason": "Explicit fixture legacy import."}
			in := contentFixtureDraft(t, root, "WorkspaceTaskProblemDraft@1", "fixture/legacy-import", fields)
			raw, _ := os.ReadFile(in.File)
			value, _ := canonicaljson.DecodeStrict(raw)
			m := contentFields(value)
			for _, wrong := range []string{"digest", "version"} {
				v := version
				digest := digestTaskBytes(before)
				if wrong == "digest" {
					digest = digestTaskBytes([]byte("different registry"))
				} else {
					v = 3 - version
				}
				m["registry_upgrade"] = contentObject(map[string]canonicaljson.Value{"from_version": int64(v), "registry_sha256": digest})
				b, _ := canonicaljson.Marshal(contentObject(m))
				if err = os.WriteFile(in.File, b, 0600); err != nil {
					t.Fatal(err)
				}
				out, e := RecordTaskProblem(f.dependencies, in)
				if e == nil || out.Attempted.ObjectPublicationStarted {
					t.Fatalf("bad %s acknowledgement published: %#v %v", wrong, out, e)
				}
				after, _ := os.ReadFile(workItemsPath(root))
				if !bytes.Equal(before, after) {
					t.Fatal("bad acknowledgement changed registry")
				}
			}
			m["registry_upgrade"] = contentObject(map[string]canonicaljson.Value{"from_version": int64(version), "registry_sha256": digestTaskBytes(before)})
			b, _ := canonicaljson.Marshal(contentObject(m))
			if err = os.WriteFile(in.File, b, 0600); err != nil {
				t.Fatal(err)
			}
			out, err := RecordTaskProblem(f.dependencies, in)
			if err != nil {
				t.Fatal(err)
			}
			manifest, err := readContentManifest(f.dependencies.TaskContent, root, out.OutcomeRef.ManifestSHA256)
			if err != nil {
				t.Fatal(err)
			}
			mm := contentFields(manifest)
			if contentString(contentFields(mm["origin"]), "kind") != "legacy_summary_import" || contentInt(mm, "revision") != 1 {
				t.Fatalf("fabricated import history: %#v", mm)
			}
			doc := contentFields(contentArray(mm, "documents")[0])
			snapshot, err := os.ReadFile(contentString(doc, "locator"))
			if err != nil || string(snapshot) != "# "+task.Title+"\n\n"+task.Description+"\n" {
				t.Fatalf("legacy summary bytes changed: %q %v", snapshot, err)
			}
			backup, err := f.dependencies.TaskContent.Read(root, "backups", digestTaskBytes(before))
			if err != nil || !bytes.Equal(backup, before) {
				t.Fatal("migration backup differs")
			}
			current, err := f.dependencies.WorkItems.Snapshot(root)
			if err != nil || !contentTypedEqual(current.Tasks, r.Tasks) || !contentTypedEqual(current.Epics, r.Epics) || len(current.TaskSolutionSelections) != 0 || len(current.TaskSpecAssessments) != 0 {
				t.Fatalf("import changed legacy facts or invented approval: %v", err)
			}
			retry, err := RecordTaskProblem(f.dependencies, in)
			if err != nil || retry.Classification != "existing" || retry.OutcomeRef.ManifestSHA256 != out.OutcomeRef.ManifestSHA256 {
				t.Fatalf("migration retry did not precede old ack: %#v %v", retry, err)
			}
		})
	}
}
