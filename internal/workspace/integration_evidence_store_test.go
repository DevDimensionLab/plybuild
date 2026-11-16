package workspace

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Build real, validated evidence in a private fixture, then seed a fresh store.
// This is test seeding, not a production migration path.
func evidenceHistoryFixture(t *testing.T) WorkItemRegistry {
	t.Helper()
	f := newIntegrationJourneyFixture(t)
	f.dependencies.TaskLifecycleIDs = cryptoTaskLifecycleIDSource{}
	git := &noEffectIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	in := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
	for i := 0; i < 8; i++ {
		in.Apply, in.Confirmation = false, ""
		check, err := CheckTaskIntegration(f.dependencies, in)
		if err != nil {
			t.Fatal(err)
		}
		in.Apply, in.Confirmation = true, integrationDigest(t, check.Readback)
		out, err := ApplyTaskIntegration(f.dependencies, in)
		if err != nil || out.Readback.Classification != "no_effect" {
			t.Fatalf("retry %d: %s %v", i, out.Readback.Classification, err)
		}
		r, err := f.dependencies.WorkItems.Snapshot(testWork(t, f.root))
		if err != nil {
			t.Fatal(err)
		}
		for j := range r.IntegrationResults {
			v := &r.IntegrationResults[j]
			if (in.RetryAfterResultID == nil && v.RetryAfterResultID == nil) || (in.RetryAfterResultID != nil && v.RetryAfterResultID != nil && *v.RetryAfterResultID == *in.RetryAfterResultID) {
				id := v.ID
				in.RetryAfterResultID = &id
				break
			}
		}
	}
	r, err := f.dependencies.WorkItems.Snapshot(testWork(t, f.root))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func publishEvidenceFixture(t *testing.T, s *systemWorkItemStore, root string, r WorkItemRegistry) {
	t.Helper()
	if err := s.WithLock(testWork(t, root), func(session WorkItemStoreSession) error { return session.Publish(r) }); err != nil {
		t.Fatal(err)
	}
}

func TestIntegrationEvidenceTaskCreationDoesNotRewriteHistory(t *testing.T) {
	r := evidenceHistoryFixture(t)
	root := createProjectWorkspace(t)
	s := newSystemWorkItemStore()
	publishEvidenceFixture(t, s, root, r)
	before := evidenceFileSnapshot(t, root)
	inline, err := encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	// A normal new unbound Task, without any integration operation.
	r.Tasks = append(r.Tasks, TaskRecord{ID: "another-task", Title: "Another Task", Description: "Synthetic measurement", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply", GitCommonDir: r.Tasks[0].GitCommonDir, WorktreeState: WorkItemUnbound})
	sortWorkRegistry(&r)
	publishEvidenceFixture(t, s, root, r)
	actualBytes, actualFiles := evidencePublishedDelta(t, before, evidenceFileSnapshot(t, root))
	if actualFiles != 1 {
		t.Fatal("unexpected published files")
	}
	for path, info := range before {
		if filepath.Base(path) == WorkItemsFile {
			continue
		}
		after, err := os.Lstat(path)
		if err != nil || !os.SameFile(info, after) || !sameRegistryFileState(info, after) {
			t.Fatalf("historical object rewritten: %s %v", path, err)
		}
	}
	disk, err := os.ReadFile(testWork(t, root).Registry())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("MEASURE published Task payload: %d bytes/%d files; authorities/intents/attempts/results=8/8/8/8; unchanged historical objects=%d", actualBytes, actualFiles, 32+len(r.TaskContentPublications))
	if bytes.Contains(disk, []byte("plan_sha256:")) || len(disk) >= len(inline)/2 {
		t.Fatal("ordinary Task publication serialized historical integration payloads again")
	}
	objects, err := filepath.Glob(filepath.Join(testWork(t, root).Dir, "task-content", "v1", "integration-evidence", "sha256", "*"))
	if err != nil || len(objects) != 32+len(r.TaskContentPublications) {
		t.Fatalf("want %d immutable records, got %d: %v", 32+len(r.TaskContentPublications), len(objects), err)
	}
	restored, err := s.Snapshot(testWork(t, root))
	if err != nil || len(restored.IntegrationResults) != 8 {
		t.Fatalf("history unavailable: %v", err)
	}
}

// Count payloads actually published by inode/state changes, including new object
// files. Ignore the lock and marker; directories are not payloads. This detects
// a rewrite even when an object is replaced with byte-identical content.
func evidenceFileSnapshot(t *testing.T, root string) map[string]os.FileInfo {
	t.Helper()
	out := map[string]os.FileInfo{}
	err := filepath.WalkDir(testWork(t, root).Dir, func(path string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if e.IsDir() || filepath.Base(path) == workItemsLock || filepath.Base(path) == MarkerFile {
			return nil
		}
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		out[path] = info
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}
func evidencePublishedDelta(t *testing.T, before, after map[string]os.FileInfo) (int64, int) {
	t.Helper()
	var size int64
	files := 0
	for path, info := range after {
		old, found := before[path]
		if !found || !os.SameFile(old, info) || !sameRegistryFileState(old, info) {
			size += info.Size()
			files++
		}
	}
	return size, files
}

func TestIntegrationEvidenceFailClosed(t *testing.T) {
	r := evidenceHistoryFixture(t)
	for _, fault := range []string{"missing", "hash", "symlink", "wrong-reference", "valid-hash-invalid-record", "unknown-envelope-field"} {
		t.Run(fault, func(t *testing.T) {
			root := createProjectWorkspace(t)
			s := newSystemWorkItemStore()
			publishEvidenceFixture(t, s, root, r)
			err := s.WithLock(testWork(t, root), func(session WorkItemStoreSession) error {
				if _, err := session.Snapshot(); err != nil {
					return err
				}
				raw, err := os.ReadFile(testWork(t, root).Registry())
				if err != nil {
					return err
				}
				node, err := parseRegistryYAML(raw)
				if err != nil {
					return err
				}
				logical := registryNodeField(node, "registry")
				refs := registryNodeField(logical, "integration_authorities")
				path := taskContentPath(testWork(t, root), integrationEvidenceKind, refs.Content[0].Value)
				switch fault {
				case "missing":
					err = os.Remove(path)
				case "hash":
					if err = os.Chmod(path, 0600); err == nil {
						err = os.WriteFile(path, []byte("corrupt: true\n"), 0400)
					}
				case "symlink":
					target := filepath.Join(root, "fake-evidence")
					data, e := os.ReadFile(path)
					if e != nil {
						return e
					}
					if err = os.WriteFile(target, data, 0400); err == nil {
						err = os.Remove(path)
					}
					if err == nil {
						err = os.Symlink(target, path)
					}
				case "wrong-reference":
					refs.Content[0].Value = registryNodeField(logical, "integration_results").Content[0].Value
				case "valid-hash-invalid-record":
					data, e := os.ReadFile(path)
					if e != nil {
						return e
					}
					record, e := parseRegistryYAML(data)
					if e != nil {
						return e
					}
					registryNodeField(record, "plan_sha256").Value = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
					data, e = marshalRegistryNode(record)
					if e != nil {
						return e
					}
					refs.Content[0].Value, err = evidenceStorage(s).Publish(testWork(t, root), integrationEvidenceKind, data)
				case "unknown-envelope-field":
					node.Content[2].Value = "unknown"
				}
				if err != nil {
					return err
				}
				if fault == "wrong-reference" || fault == "valid-hash-invalid-record" || fault == "unknown-envelope-field" {
					raw, err = marshalRegistryNode(node)
					if err != nil {
						return err
					}
					if err := os.WriteFile(testWork(t, root).Registry(), raw, 0600); err != nil {
						return err
					}
				}
				if _, err := session.Snapshot(); err == nil {
					return fmt.Errorf("corruption accepted inside confirmed session")
				}
				if err := session.Publish(r); err == nil {
					return fmt.Errorf("corruption repaired or bypassed by publication")
				}
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
			if _, err := s.SnapshotRegistrations(testWork(t, root)); err == nil {
				t.Fatal("registration reader accepted corruption")
			}
		})
	}
}

func TestIntegrationEvidenceInterruptedPublicationAndLostReply(t *testing.T) {
	r := evidenceHistoryFixture(t)
	for _, stage := range []string{"integration-evidence-object-write", "integration-evidence-object-rename", "integration-evidence-object-directory-sync", "replace", "directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			root := createProjectWorkspace(t)
			s := newSystemWorkItemStore()
			publishEvidenceFixture(t, s, root, emptyWorkItemRegistry())
			s.faults = &workItemStoreFaults{fail: func(op string) error {
				if op == stage {
					return errInjected
				}
				return nil
			}}
			err := s.WithLock(testWork(t, root), func(session WorkItemStoreSession) error { return session.Publish(r) })
			if err == nil {
				t.Fatal("fault did not fire")
			}
			s.faults = nil
			observed, err := s.Snapshot(testWork(t, root))
			if err != nil {
				t.Fatal(err)
			}
			want := 0
			if stage == "directory-sync" {
				want = 8
			}
			if len(observed.IntegrationResults) != want {
				t.Fatal("partial state exposed")
			}
			before := evidenceFileSnapshot(t, root)
			publishEvidenceFixture(t, s, root, r)
			after := evidenceFileSnapshot(t, root)
			for path, info := range before {
				if filepath.Base(path) == WorkItemsFile {
					continue
				}
				now := after[path]
				if now == nil || !os.SameFile(info, now) || !sameRegistryFileState(info, now) {
					t.Fatalf("retry rewrote an orphan object: %s", path)
				}
			}
			err = s.WithLock(testWork(t, root), func(session WorkItemStoreSession) error {
				got, err := publishWorkItemRegistryRecover(session, r)
				if err != nil {
					return err
				}
				a, _ := encodeWorkItemRegistry(got)
				b, _ := encodeWorkItemRegistry(r)
				if !bytes.Equal(a, b) {
					return fmt.Errorf("lost reply retry changed history")
				}
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

// The store has one physical format. An inline logical registry on disk is
// rejected without writes, and the logical decoder rejects the envelope.
func TestIntegrationEvidenceSingleStoredFormat(t *testing.T) {
	r := evidenceHistoryFixture(t)
	inlineRoot := createProjectWorkspace(t)
	inline, err := encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	writeWorkItems(t, inlineRoot, string(inline))
	s := newSystemWorkItemStore()
	if _, err := s.Snapshot(testWork(t, inlineRoot)); RegistryError(err) == nil || RegistryError(err).Code != "registry_invalid_content" {
		t.Fatalf("inline registry accepted: %v", err)
	}
	if err := s.WithLock(testWork(t, inlineRoot), func(session WorkItemStoreSession) error { return session.Publish(r) }); err == nil {
		t.Fatal("inline registry overwritten")
	}
	if after, err := os.ReadFile(testWork(t, inlineRoot).Registry()); err != nil || !bytes.Equal(after, inline) {
		t.Fatal("inline registry changed", err)
	}
	root := createProjectWorkspace(t)
	publishEvidenceFixture(t, s, root, r)
	raw, _ := os.ReadFile(testWork(t, root).Registry())
	if _, err := decodeWorkItemRegistry(raw); RegistryError(err) == nil {
		t.Fatal("logical decoder accepted the physical envelope", err)
	}
	got, err := s.Snapshot(testWork(t, root))
	if err != nil {
		t.Fatal(err)
	}
	a, _ := encodeWorkItemRegistry(got)
	b, _ := encodeWorkItemRegistry(r)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("historical logical evidence changed")
	}
}

func TestIntegrationEvidenceReaderFreshnessAndNoGitReplay(t *testing.T) {
	for _, fault := range []string{"missing", "object-mode", "directory-link"} {
		t.Run(fault, func(t *testing.T) {
			f := newIntegrationJourneyFixture(t)
			git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit, lost: true}
			f.dependencies.IntegrationGit = git
			in := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
			check, err := CheckTaskIntegration(f.dependencies, in)
			if err != nil {
				t.Fatal(err)
			}
			in.Apply, in.Confirmation = true, integrationDigest(t, check.Readback)
			if _, err := ApplyTaskIntegration(f.dependencies, in); err != nil {
				t.Fatal(err)
			}
			basis, err := ReadViewBasis(f.dependencies)
			if err != nil {
				t.Fatal(err)
			}
			journal, err := ReadTaskJournalBasis(f.dependencies, "task")
			if err != nil {
				t.Fatal(err)
			}
			sources := IntegrationEvidenceSources(testWork(t, f.root), journal.Registry)
			if len(sources) != 4+len(journal.Registry.TaskContentPublications) {
				t.Fatalf("missing evidence sources: %d", len(sources))
			}
			for _, src := range sources {
				if basis.Sources[src.Locator] != src.SHA256 {
					t.Fatal("read view omits evidence freshness")
				}
			}
			switch fault {
			case "missing":
				err = os.Remove(sources[0].Locator)
			case "object-mode":
				err = os.Chmod(sources[0].Locator, 0600)
			case "directory-link":
				dir := filepath.Dir(sources[0].Locator)
				err = os.Rename(dir, dir+"-retained")
				if err == nil {
					err = os.Symlink(dir+"-retained", dir)
				}
			}
			if err != nil {
				t.Fatal(err)
			}
			if status, _ := CheckReadViewBasis(f.dependencies, basis); status != "unknown" {
				t.Fatal("read view lost its evidence but remains fresh")
			}
			if _, err := ReadTaskJournalBasis(f.dependencies, "task"); err == nil {
				t.Fatal("journal accepted missing evidence")
			}
			if _, err := ApplyTaskIntegration(f.dependencies, in); err == nil {
				t.Fatal("replay accepted missing authority evidence")
			}
			if git.count() != 1 {
				t.Fatal("corruption started another Git effect")
			}
		})
	}
}

func TestIntegrationEvidencePreservesUnchangedObjectReferences(t *testing.T) {
	r := evidenceHistoryFixture(t)
	root := createProjectWorkspace(t)
	s := newSystemWorkItemStore()
	publishEvidenceFixture(t, s, root, r)
	raw, err := os.ReadFile(testWork(t, root).Registry())
	if err != nil {
		t.Fatal(err)
	}
	node, err := parseRegistryYAML(raw)
	if err != nil {
		t.Fatal(err)
	}
	refs := registryNodeField(registryNodeField(node, "registry"), "integration_authorities")
	payload, err := evidenceStorage(s).Read(testWork(t, root), integrationEvidenceKind, refs.Content[0].Value)
	if err != nil {
		t.Fatal(err)
	}
	// A valid pre-existing object may have another YAML spelling. Its exact
	// digest remains authoritative; an unrelated Task write must not normalize it.
	digest, err := evidenceStorage(s).Publish(testWork(t, root), integrationEvidenceKind, append(payload, []byte("\n# preserved historical spelling\n")...))
	if err != nil {
		t.Fatal(err)
	}
	refs.Content[0].Value = digest
	raw, err = marshalRegistryNode(node)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(testWork(t, root).Registry(), raw, 0600); err != nil {
		t.Fatal(err)
	}
	r, err = s.Snapshot(testWork(t, root))
	if err != nil {
		t.Fatal(err)
	}
	r.Tasks[0].Title = "Unrelated Task edit"
	before := evidenceFileSnapshot(t, root)
	publishEvidenceFixture(t, s, root, r)
	changed, files := evidencePublishedDelta(t, before, evidenceFileSnapshot(t, root))
	if files != 1 || changed == 0 {
		t.Fatal("unexpected evidence publication")
	}
	r, err = s.Snapshot(testWork(t, root))
	if err != nil || r.integrationEvidenceRefs[0] != digest {
		t.Fatal("unchanged historical reference was normalized", err)
	}
}

func TestIntegrationEvidenceUnknownAndPartialNeverRepeatGitEffect(t *testing.T) {
	for _, outcome := range []string{"unknown", "partial"} {
		t.Run(outcome, func(t *testing.T) {
			f := newIntegrationJourneyFixture(t)
			if outcome == "unknown" {
				f.dependencies.IntegrationGit = &sourceFailureAfterCommandGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
			} else {
				f.dependencies.IntegrationGit = &partialAfterCommandGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
			}
			git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
			f.dependencies.IntegrationGit = git
			in := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
			check, err := CheckTaskIntegration(f.dependencies, in)
			if err != nil {
				t.Fatal(err)
			}
			in.Apply, in.Confirmation = true, integrationDigest(t, check.Readback)
			for i := 0; i < 2; i++ {
				out, err := ApplyTaskIntegration(f.dependencies, in)
				if err != nil || out.Readback.Classification != outcome || git.count() != 1 {
					t.Fatalf("apply/replay: %s %v effects=%d", out.Readback.Classification, err, git.count())
				}
			}
			r, err := f.dependencies.WorkItems.Snapshot(testWork(t, f.root))
			if err != nil || len(r.IntegrationResults) != 1 {
				t.Fatal("uncertain history lost", err)
			}
			prior := r.IntegrationResults[0].ID
			in.Apply, in.Confirmation, in.RetryAfterResultID = false, "", &prior
			if check, err := CheckTaskIntegration(f.dependencies, in); err == nil && check.Readback.Classification == "ready" {
				t.Fatal("uncertain result authorized retry")
			}
			if git.count() != 1 {
				t.Fatal("unknown/partial retry caused another Git effect")
			}
		})
	}
}

func TestIntegrationEvidenceConcurrentTasksPreserveHistory(t *testing.T) {
	r := evidenceHistoryFixture(t)
	root := createProjectWorkspace(t)
	s := newSystemWorkItemStore()
	publishEvidenceFixture(t, s, root, r)
	before := evidenceFileSnapshot(t, root)
	start, done := make(chan struct{}), make(chan error, 2)
	for _, id := range []TaskID{"concurrent-a", "concurrent-b"} {
		go func(id TaskID) {
			<-start
			done <- s.WithLock(testWork(t, root), func(session WorkItemStoreSession) error {
				current, err := session.Snapshot()
				if err != nil {
					return err
				}
				current.Tasks = append(current.Tasks, TaskRecord{ID: id, Title: string(id), Description: "Concurrent isolated fixture", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply", GitCommonDir: r.Tasks[0].GitCommonDir, WorktreeState: WorkItemUnbound})
				sortWorkRegistry(&current)
				return session.Publish(current)
			})
		}(id)
	}
	close(start)
	for i := 0; i < 2; i++ {
		if err := <-done; err != nil {
			t.Fatal(err)
		}
	}
	got, err := s.Snapshot(testWork(t, root))
	if err != nil || len(got.Tasks) != len(r.Tasks)+2 || len(got.IntegrationResults) != 8 {
		t.Fatal("concurrent Task/evidence update lost", err)
	}
	after := evidenceFileSnapshot(t, root)
	for path, info := range before {
		if filepath.Base(path) == WorkItemsFile {
			continue
		}
		now := after[path]
		if now == nil || !os.SameFile(info, now) || !sameRegistryFileState(info, now) {
			t.Fatal("concurrent Task creation rewrote evidence", path)
		}
	}
}
