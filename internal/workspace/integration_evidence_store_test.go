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
		r, err := f.dependencies.WorkItems.Snapshot(f.root)
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
	r, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func publishEvidenceFixture(t *testing.T, s *systemWorkItemStore, root string, r WorkItemRegistry) {
	t.Helper()
	if err := s.WithLock(root, func(session WorkItemStoreSession) error { return session.Publish(r) }); err != nil {
		t.Fatal(err)
	}
}

func TestIntegrationEvidenceTaskCreationDoesNotRewriteHistory(t *testing.T) {
	r := evidenceHistoryFixture(t)
	root := createProjectWorkspace(t)
	s := newSystemWorkItemStore()
	publishEvidenceFixture(t, s, root, r)
	before := evidenceFileSnapshot(t, root)
	legacyRoot := createProjectWorkspace(t)
	baseline, err := encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(workItemsPath(legacyRoot), baseline, 0600); err != nil {
		t.Fatal(err)
	}
	legacyBefore := evidenceFileSnapshot(t, legacyRoot)
	// A normal new unbound Task, without any integration operation.
	r.Tasks = append(r.Tasks, TaskRecord{ID: "another-task", Title: "Another Task", Description: "Synthetic measurement", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply", GitCommonDir: r.Tasks[0].GitCommonDir, WorktreeState: WorkItemUnbound})
	sortWorkRegistry(&r)
	legacy, err := encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	publishEvidenceFixture(t, s, root, r)
	publishEvidenceFixture(t, s, legacyRoot, r)
	actualBytes, actualFiles := evidencePublishedDelta(t, before, evidenceFileSnapshot(t, root))
	baselineBytes, baselineFiles := evidencePublishedDelta(t, legacyBefore, evidenceFileSnapshot(t, legacyRoot))
	if actualFiles != 1 || baselineFiles != 1 || baselineBytes != int64(len(legacy)) {
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
	disk, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("MEASURE published Task payload: before=%d bytes/%d files; after=%d bytes/%d files; reduction=%.2f%%; authorities/intents/attempts/results=8/8/8/8; unchanged historical objects=%d", baselineBytes, baselineFiles, actualBytes, actualFiles, 100*(1-float64(actualBytes)/float64(baselineBytes)), 32+len(r.TaskContentPublications))
	var stock int64
	for _, info := range evidenceFileSnapshot(t, root) {
		stock += info.Size()
	}
	t.Logf("STOCK new store=%d bytes (registry plus retained evidence); stock is not publication I/O", stock)
	if bytes.Contains(disk, []byte("plan_sha256:")) || len(disk) >= len(legacy)/2 {
		t.Fatal("ordinary Task publication serialized historical integration payloads again")
	}
	objects, err := filepath.Glob(filepath.Join(root, MarkerDirectory, "task-content", "v1", "integration-evidence", "sha256", "*"))
	if err != nil || len(objects) != 32+len(r.TaskContentPublications) {
		t.Fatalf("want %d immutable records, got %d: %v", 32+len(r.TaskContentPublications), len(objects), err)
	}
	restored, err := s.Snapshot(root)
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
	err := filepath.WalkDir(filepath.Join(root, MarkerDirectory), func(path string, e fs.DirEntry, err error) error {
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
			err := s.WithLock(root, func(session WorkItemStoreSession) error {
				if _, err := session.Snapshot(); err != nil {
					return err
				}
				raw, err := os.ReadFile(workItemsPath(root))
				if err != nil {
					return err
				}
				node, err := parseRegistryYAML(raw)
				if err != nil {
					return err
				}
				logical := registryNodeField(node, "registry")
				refs := registryNodeField(logical, "integration_authorities")
				path := taskContentPath(root, integrationEvidenceKind, refs.Content[0].Value)
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
					refs.Content[0].Value, err = evidenceStorage(s).Publish(root, integrationEvidenceKind, data)
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
					if err := os.WriteFile(workItemsPath(root), raw, 0600); err != nil {
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
			if _, err := s.SnapshotRegistrations(root); err == nil {
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
			err := s.WithLock(root, func(session WorkItemStoreSession) error { return session.Publish(r) })
			if err == nil {
				t.Fatal("fault did not fire")
			}
			s.faults = nil
			observed, err := s.Snapshot(root)
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
			err = s.WithLock(root, func(session WorkItemStoreSession) error {
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

func TestIntegrationEvidenceLegacyBoundary(t *testing.T) {
	r := evidenceHistoryFixture(t)
	root := createProjectWorkspace(t)
	raw, err := encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(workItemsPath(root), raw, 0600); err != nil {
		t.Fatal(err)
	}
	s := newSystemWorkItemStore()
	publishEvidenceFixture(t, s, root, r)
	after, err := os.ReadFile(workItemsPath(root))
	if err != nil || !bytes.Equal(raw, after) {
		t.Fatal("legacy registry silently migrated", err)
	}
	if _, err := os.Lstat(filepath.Join(root, MarkerDirectory, "task-content")); !os.IsNotExist(err) {
		t.Fatal("legacy writer introduced external evidence")
	}
	fresh := createProjectWorkspace(t)
	publishEvidenceFixture(t, s, fresh, r)
	raw, _ = os.ReadFile(workItemsPath(fresh))
	if _, err := decodeWorkItemRegistry(raw); RegistryError(err) == nil || RegistryError(err).Code != "registry_unsupported_format" {
		t.Fatal("old reader must reject new physical format", err)
	}
	got, err := s.Snapshot(fresh)
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
			sources := IntegrationEvidenceSources(f.root, journal.Registry)
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
	raw, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	node, err := parseRegistryYAML(raw)
	if err != nil {
		t.Fatal(err)
	}
	refs := registryNodeField(registryNodeField(node, "registry"), "integration_authorities")
	payload, err := evidenceStorage(s).Read(root, integrationEvidenceKind, refs.Content[0].Value)
	if err != nil {
		t.Fatal(err)
	}
	// A valid pre-existing object may have another YAML spelling. Its exact
	// digest remains authoritative; an unrelated Task write must not normalize it.
	digest, err := evidenceStorage(s).Publish(root, integrationEvidenceKind, append(payload, []byte("\n# preserved historical spelling\n")...))
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
	r, err = s.Snapshot(root)
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
	r, err = s.Snapshot(root)
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
			r, err := f.dependencies.WorkItems.Snapshot(f.root)
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
			done <- s.WithLock(root, func(session WorkItemStoreSession) error {
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
	got, err := s.Snapshot(root)
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
