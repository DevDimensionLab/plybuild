package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func publicationReferences(t *testing.T, root string) *yaml.Node {
	t.Helper()
	raw, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	node, err := parseRegistryYAML(raw)
	if err != nil {
		t.Fatal(err)
	}
	refs := registryNodeField(registryNodeField(node, "registry"), "task_content_publications")
	if refs == nil || len(refs.Content) == 0 {
		t.Fatal("publication history missing")
	}
	for _, ref := range refs.Content {
		if ref.Kind != yaml.ScalarNode || !digestPattern.MatchString(ref.Value) {
			t.Fatal("publication history is still inline; want immutable digest references")
		}
	}
	return refs
}

func TestPublicationEvidenceNativeWritesPreserveHistory(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	refs := publicationReferences(t, root)
	before := evidenceFileSnapshot(t, root)
	r, err := f.dependencies.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	prior := r.TaskContentPublications
	input, err := ParseTaskCreateInput("another-task", "Another Task", "Synthetic measurement", "epic", "ply", "ply")
	if err != nil {
		t.Fatal(err)
	}
	if out, err := CreateTask(f.dependencies, input); err != nil || !out.Created {
		t.Fatalf("create: %v", err)
	}
	for path, info := range before {
		if filepath.Base(path) == WorkItemsFile {
			continue
		}
		after, err := os.Lstat(path)
		if err != nil || !os.SameFile(info, after) || !sameRegistryFileState(info, after) {
			t.Fatalf("historical payload rewritten: %s %v", path, err)
		}
	}
	afterRefs := publicationReferences(t, root)
	if len(refs.Content) != len(afterRefs.Content) {
		t.Fatal("historical reference count changed")
	}
	for i := range refs.Content {
		if refs.Content[i].Value != afterRefs.Content[i].Value {
			t.Fatal("historical references changed")
		}
	}
	if false {
		t.Fatal("historical references changed")
	}
	in := contentFixtureProblem(t, f, "fixture/revised-problem")
	if _, err := RecordTaskProblem(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	stable := evidenceFileSnapshot(t, root)
	if out, err := RecordTaskProblem(f.dependencies, in); err != nil || out.Classification != "existing" {
		t.Fatalf("retry: %#v %v", out, err)
	}
	if n, _ := evidencePublishedDelta(t, stable, evidenceFileSnapshot(t, root)); n != 0 {
		t.Fatal("retry republished payloads")
	}
	r, err = f.dependencies.WorkItems.Snapshot(root)
	if err != nil || len(r.TaskContentPublications) != len(prior)+1 {
		t.Fatalf("read history: %v", err)
	}
	if _, err := ShowTaskProblem(f.dependencies, TaskContentQuery{TaskID: "task"}); err != nil {
		t.Fatal(err)
	}
	publicationReferences(t, root)
}

func TestPublicationEvidenceFailClosed(t *testing.T) {
	for _, fault := range []string{"missing", "hash", "mode", "symlink", "wrong-type", "invalid-record", "request-mismatch", "inline-in-format10"} {
		t.Run(fault, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			root := filepath.Dir(f.wrapper)
			s := f.dependencies.WorkItems
			input := contentFixtureProblem(t, f, "fixture/check-corruption")
			err := s.WithLock(root, func(session WorkItemStoreSession) error {
				r, err := session.Snapshot()
				if err != nil {
					return err
				}
				refs := publicationReferences(t, root)
				path := taskContentPath(root, integrationEvidenceKind, refs.Content[0].Value)
				raw, err := os.ReadFile(workItemsPath(root))
				if err != nil {
					return err
				}
				node, err := parseRegistryYAML(raw)
				if err != nil {
					return err
				}
				refs = registryNodeField(registryNodeField(node, "registry"), "task_content_publications")
				switch fault {
				case "missing":
					err = os.Remove(path)
				case "mode":
					err = os.Chmod(path, 0600)
				case "hash":
					if err = os.Chmod(path, 0600); err == nil {
						err = os.WriteFile(path, []byte("corrupt: true\n"), 0400)
					}
				case "symlink":
					if err = os.Rename(path, path+"-real"); err == nil {
						err = os.Symlink(path+"-real", path)
					}
				default:
					payload, e := os.ReadFile(path)
					if e != nil {
						return e
					}
					record, e := parseRegistryYAML(payload)
					if e != nil {
						return e
					}
					switch fault {
					case "inline-in-format10":
						refs.Content[0] = record
					case "wrong-type":
						payload = []byte("id: wrong-record\n")
					case "invalid-record":
						registryNodeField(record, "operation").Value = "unknown"
					case "request-mismatch":
						registryNodeField(record, "request_sha256").Value = digestTaskBytes([]byte("wrong request"))
					}
					if fault == "invalid-record" || fault == "request-mismatch" {
						payload, err = marshalRegistryNode(record)
					}
					if err != nil {
						return err
					}
					if fault != "inline-in-format10" {
						refs.Content[0].Value, err = (&TaskContentStorage{}).Publish(root, integrationEvidenceKind, payload)
					}
					if err == nil {
						raw, err = marshalRegistryNode(node)
					}
					if err == nil {
						err = os.WriteFile(workItemsPath(root), raw, 0600)
					}
				}
				if err != nil {
					return err
				}
				if fault == "request-mismatch" {
					return nil
				} // Request chains are validated by native content operations, not registry-only reads.
				if _, err := session.Snapshot(); err == nil {
					return fmt.Errorf("%s accepted by confirmed snapshot", fault)
				}
				if err := session.Publish(r); err == nil {
					return fmt.Errorf("%s repaired by publication", fault)
				}
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
			if fault != "request-mismatch" {
				if _, err := s.Snapshot(root); err == nil {
					t.Fatal("fresh reader accepted corrupted publication")
				}
			}
			if read, err := ShowTaskProblem(f.dependencies, TaskContentQuery{TaskID: "task"}); err == nil && contentString(contentFields(read.Value), "integrity") == "valid" {
				t.Fatal("native content reader accepted corruption")
			}
			if _, err := RecordTaskProblem(f.dependencies, input); err == nil {
				t.Fatal("native content writer accepted corruption")
			}
		})
	}
}

func TestPublicationEvidenceInterruptedNativePublication(t *testing.T) {
	for _, stage := range []string{"integration-evidence-object-write", "integration-evidence-object-rename", "integration-evidence-object-directory-sync", "replace", "directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			root := filepath.Dir(f.wrapper)
			input := contentFixtureProblem(t, f, "fixture/recovery")
			before := evidenceFileSnapshot(t, root)
			s := newSystemWorkItemStore()
			s.faults = &workItemStoreFaults{fail: func(op string) error {
				if op == stage {
					return errors.New("injected " + stage)
				}
				return nil
			}}
			f.dependencies.WorkItems = s
			if _, err := RecordTaskProblem(f.dependencies, input); err == nil {
				t.Fatal("injected fault reported success")
			}
			s.faults = nil
			r, err := s.Snapshot(root)
			if err != nil {
				t.Fatal(err)
			}
			want := 1
			if stage == "directory-sync" {
				want = 2
			}
			if len(r.TaskProblemRevisions) != want {
				t.Fatalf("partial publication: %d", len(r.TaskProblemRevisions))
			}
			if _, err := RecordTaskProblem(f.dependencies, input); err != nil {
				t.Fatal(err)
			}
			r, err = s.Snapshot(root)
			if err != nil || len(r.TaskProblemRevisions) != 2 || len(r.TaskContentPublications) != 2 {
				t.Fatalf("recovery lost exact outcome: %v", err)
			}
			stable := evidenceFileSnapshot(t, root)
			if out, err := RecordTaskProblem(f.dependencies, input); err != nil || out.Classification != "existing" {
				t.Fatalf("replay: %v", err)
			}
			if n, _ := evidencePublishedDelta(t, stable, evidenceFileSnapshot(t, root)); n != 0 {
				t.Fatal("replay rewrote payloads")
			}
			for path, info := range before {
				if filepath.Base(path) == WorkItemsFile {
					continue
				}
				after, err := os.Lstat(path)
				if err != nil || !os.SameFile(info, after) || !sameRegistryFileState(info, after) {
					t.Fatalf("recovery rewrote history: %s", path)
				}
			}
		})
	}
}

func TestPublicationEvidenceFormat9Boundary(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	s := newSystemWorkItemStore()
	r, err := s.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	// Construct the prior physical format in a private fixture, never live data.
	old := r
	old.publicationEvidence = false
	old.integrationEvidenceRefs = nil // This fixture has no integration records.
	raw, err := s.encodeStoredRegistry(root, &r, old)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(workItemsPath(root), raw, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f, "fixture/format9")); err != nil {
		t.Fatal(err)
	}
	raw, err = os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	node, err := parseRegistryYAML(raw)
	if err != nil {
		t.Fatal(err)
	}
	if registryNodeField(node, "format_version").Value != "9" {
		t.Fatal("existing format9 store migrated")
	}
	pubs := registryNodeField(registryNodeField(node, "registry"), "task_content_publications")
	if len(pubs.Content) != 2 || pubs.Content[0].Kind != yaml.MappingNode {
		t.Fatal("format9 publications must remain inline")
	}
	r, err = s.Snapshot(root)
	if err != nil || r.publicationEvidence {
		t.Fatalf("format9 read: %v", err)
	}
	if !reflect.DeepEqual(*taskPublication(r, old.TaskContentPublications[0].PublicationKey), old.TaskContentPublications[0]) {
		t.Fatal("format9 history changed")
	}
}

func TestPublicationEvidencePreservesExactHistoricalBytes(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	s := newSystemWorkItemStore()
	raw, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	node, err := parseRegistryYAML(raw)
	if err != nil {
		t.Fatal(err)
	}
	refs := registryNodeField(registryNodeField(node, "registry"), "task_content_publications")
	payload, err := evidenceStorage(s).Read(root, integrationEvidenceKind, refs.Content[0].Value)
	if err != nil {
		t.Fatal(err)
	}
	digest, err := evidenceStorage(s).Publish(root, integrationEvidenceKind, append(payload, []byte("\n# exact historical bytes\n")...))
	if err != nil {
		t.Fatal(err)
	}
	refs.Content[0].Value = digest
	raw, err = marshalRegistryNode(node)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(workItemsPath(root), raw, 0600); err != nil {
		t.Fatal(err)
	}
	prior := evidenceFileSnapshot(t, root)
	// This publication key sorts before the original Task-create publication.
	if _, err := RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f, "aaa/new-problem")); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, ref := range publicationReferences(t, root).Content {
		found = found || ref.Value == digest
	}
	if !found {
		t.Fatal("historical reference was normalized or matched by array position")
	}
	for path, info := range prior {
		if filepath.Base(path) == WorkItemsFile {
			continue
		}
		after, err := os.Lstat(path)
		if err != nil || !os.SameFile(info, after) || !sameRegistryFileState(info, after) {
			t.Fatalf("historical payload rewritten: %s", path)
		}
	}
	basis, err := ReadViewBasis(f.dependencies)
	if err != nil {
		t.Fatal(err)
	}
	path := taskContentPath(root, integrationEvidenceKind, digest)
	if basis.Sources[path] != digest {
		t.Fatal("read view omitted publication source")
	}
	if err := os.Chmod(path, 0600); err != nil {
		t.Fatal(err)
	}
	if status, _ := CheckReadViewBasis(f.dependencies, basis); status != "unknown" {
		t.Fatal("read view accepted publication mode corruption")
	}
}
