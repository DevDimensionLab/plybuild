package workspace

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestWorkItemSessionConfirmedGeneration(t *testing.T) {
	root := createProjectWorkspace(t)
	reads := 0
	store := &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(op string) error {
		if op == "registry-read" {
			reads++
		}
		return nil
	}}}
	initial := minimalEpicRegistry(root)
	if err := store.WithLock(root, func(s WorkItemStoreSession) error { return s.Publish(initial) }); err != nil {
		t.Fatal(err)
	}
	reads = 0
	err := store.WithLock(root, func(s WorkItemStoreSession) error {
		r, err := s.Snapshot()
		if err != nil {
			return err
		}
		r.Epics[0].Title = "Published"
		if err = s.Publish(r); err != nil {
			return err
		}
		r.Epics[0].Title = "Caller mutation"
		got, err := s.Snapshot()
		if err != nil {
			return err
		}
		if got.Epics[0].Title != "Published" {
			t.Fatal("publication input aliases confirmed generation")
		}
		bytes, err := os.ReadFile(workItemsPath(root))
		if err != nil {
			return err
		}
		if got.RawSHA256 != digestTaskBytes(bytes) {
			t.Fatal("confirmed hash does not identify persisted bytes")
		}
		got.Epics[0].Title = "Snapshot mutation"
		again, err := s.Snapshot()
		if err != nil {
			return err
		}
		if again.Epics[0].Title != "Published" {
			t.Fatal("snapshot aliases confirmed generation")
		}
		if reads != 1 {
			t.Fatalf("physical registry reads=%d, want one per unchanged lock period", reads)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if err = store.WithLock(root, func(s WorkItemStoreSession) error { _, err := s.Snapshot(); return err }); err != nil {
		t.Fatal(err)
	}
	if reads != 2 {
		t.Fatalf("relock reads=%d, want 2", reads)
	}
}

func TestWorkItemSessionPublicationRecovery(t *testing.T) {
	for _, fault := range []string{"temp-open", "write", "chmod", "file-sync", "close", "replace", "directory-sync"} {
		t.Run(fault, func(t *testing.T) {
			root := createProjectWorkspace(t)
			store := newSystemWorkItemStore()
			if err := store.WithLock(root, func(s WorkItemStoreSession) error { return s.Publish(minimalEpicRegistry(root)) }); err != nil {
				t.Fatal(err)
			}
			reads := 0
			armed := false
			store.faults = &workItemStoreFaults{fail: func(op string) error {
				if op == "registry-read" {
					reads++
				}
				if armed && op == fault {
					return errInjected
				}
				return nil
			}}
			if err := store.WithLock(root, func(s WorkItemStoreSession) error {
				r, err := s.Snapshot()
				if err != nil {
					return err
				}
				r.Epics[0].Title = "New"
				armed = true
				if err = s.Publish(r); err == nil {
					t.Fatal("fault succeeded")
				}
				armed = false
				got, err := s.Snapshot()
				if err != nil {
					return err
				}
				want := "Epic"
				if fault == "directory-sync" {
					want = "New"
				} else {
					want = minimalEpicRegistry(root).Epics[0].Title
				}
				if got.Epics[0].Title != want {
					t.Fatalf("recovery got %q want %q", got.Epics[0].Title, want)
				}
				if reads != 2 {
					t.Fatalf("recovery reads=%d want 2", reads)
				}
				physical, err := store.Snapshot(root)
				if err != nil {
					return err
				}
				if !reflect.DeepEqual(got, physical) {
					t.Fatal("recovery differs from disk")
				}
				return nil
			}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestWorkItemSessionDetectsRegistryDrift(t *testing.T) {
	root := createProjectWorkspace(t)
	store := newSystemWorkItemStore()
	if err := store.WithLock(root, func(s WorkItemStoreSession) error {
		if err := s.Publish(minimalEpicRegistry(root)); err != nil {
			return err
		}
		if _, err := s.Snapshot(); err != nil {
			return err
		}
		// Simulate an out-of-protocol replacement, including a same-size generation.
		r := minimalEpicRegistry(root)
		r.Epics[0].Title = "Replacement"
		b, err := encodeWorkItemRegistry(r)
		if err != nil {
			return err
		}
		temp := filepath.Join(root, MarkerDirectory, "replacement")
		if err = os.WriteFile(temp, b, 0600); err != nil {
			return err
		}
		if err = os.Rename(temp, workItemsPath(root)); err != nil {
			return err
		}
		got, err := s.Snapshot()
		if err != nil {
			return err
		}
		if got.Epics[0].Title != "Replacement" {
			t.Fatal("missed physical replacement")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestWorkItemSessionCurrentScopeAndExternalFreshness(t *testing.T) {
	for _, boundary := range []string{"snapshot", "publish"} {
		t.Run(boundary, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			goalQueueFixture(t, f, []TaskGoalQueueEntry{})
			root := filepath.Dir(f.wrapper)
			scope := &taskContentScope{}
			d := goalOperationContentScope(f.dependencies, scope)
			if err := d.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
				r, err := s.Snapshot()
				if err != nil {
					return err
				}
				if err = s.Publish(r); err == nil {
					t.Fatal("unbound publication succeeded")
				}
				// Observe again structurally, then bind the content scope in this lock.
				r, err = s.Snapshot()
				if err != nil {
					return err
				}
				scope.task = "task"
				if err = s.Publish(r); err != nil {
					return err
				}
				p := taskPublication(r, "task-create/task")
				if p == nil {
					t.Fatal("publication missing")
				}
				if err = os.Remove(taskContentPath(root, "requests", p.RequestSHA256)); err != nil {
					return err
				}
				if boundary == "snapshot" {
					_, err = s.Snapshot()
				} else {
					err = s.Publish(r)
				}
				if err == nil {
					t.Fatal("retained registry hid removed external content")
				}
				return nil
			}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestWorkItemSessionLateScopeBinding(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goalQueueFixture(t, f, []TaskGoalQueueEntry{})
	root := filepath.Dir(f.wrapper)
	scope := &taskContentScope{}
	d := goalOperationContentScope(f.dependencies, scope)
	if err := d.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
		r, err := s.Snapshot()
		if err != nil {
			return err
		}
		p := taskPublication(r, "task-create/task")
		if p == nil {
			t.Fatal("publication missing")
		}
		if err = os.Remove(taskContentPath(root, "requests", p.RequestSHA256)); err != nil {
			return err
		}
		scope.task = "task"
		if _, err = s.Snapshot(); err == nil {
			t.Fatal("snapshot reused initial unbound content check")
		}
		if err = s.Publish(r); err == nil {
			t.Fatal("publish reused initial unbound content check")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestWorkItemSessionNestedSnapshotOwnership(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	g := goalFixture(t, f, "task", "owned/goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &g}})
	root := filepath.Dir(f.wrapper)
	if err := f.dependencies.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
		r, err := s.Snapshot()
		if err != nil {
			return err
		}
		before := cloneWorkItemRegistry(r)
		r.TaskQueueEvents[0].Request["publication_key"] = "caller mutation"
		r.TaskQueueEvents[0].Request["entries"].([]any)[0] = nil // Nested map/slice ownership.
		got, err := s.Snapshot()
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(before, got) {
			t.Fatal("nested snapshot mutation leaked")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

type generationGuard struct{ committed WorkItemRegistry }

func (g *generationGuard) Check(_ string, _ WorkItemRegistry) error { return nil }
func (g *generationGuard) Committed(_ string, r WorkItemRegistry)   { g.committed = r }

func TestWorkItemSessionGuardCommittedHash(t *testing.T) {
	root := createProjectWorkspace(t)
	store := newSystemWorkItemStore()
	g := &generationGuard{}
	guarded := guardedTaskStartStore{store, g}
	if err := guarded.WithLock(root, func(s WorkItemStoreSession) error { return s.Publish(minimalEpicRegistry(root)) }); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	if g.committed.RawSHA256 != digestTaskBytes(b) {
		t.Fatal("guard committed hash differs from published bytes")
	}
}

func TestWorkItemSessionConfirmedMatchesPersistedShape(t *testing.T) {
	for version := 1; version <= 4; version++ {
		t.Run(string(rune('0'+version)), func(t *testing.T) {
			root := createProjectWorkspace(t)
			store := newSystemWorkItemStore()
			r := emptyWorkItemRegistry()
			r.FormatVersion = version
			if version < 3 {
				r.TaskSpecPolicies = nil
				r.TaskProblemRevisions = nil
				r.TaskSpecRevisions = nil
				r.TaskSpecAssessments = nil
				r.TaskSolutionSelections = nil
				r.TaskResultSpecBindings = nil
				r.TaskContentPublications = nil
			}
			if version == 1 {
				r.TaskResults = nil
				r.HumanQARecords = nil
				r.IntegrationAuthorities = nil
				r.IntegrationIntents = nil
				r.IntegrationAttempts = nil
				r.IntegrationResults = nil
			}
			if version == 4 {
				r.EpicBaseVersions = []EpicBaseVersion{}
				r.EpicBaseUpdates = []EpicBaseUpdate{}
				r.TaskQueueEvents = []TaskQueueEvent{}
				r.TaskPreparations = []TaskPreparation{}
			}
			seed, err := encodeWorkItemRegistry(r)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(workItemsPath(root), seed, 0600); err != nil {
				t.Fatal(err)
			}
			if err := store.WithLock(root, func(s WorkItemStoreSession) error {
				if err := s.Publish(r); err != nil {
					return err
				}
				got, err := s.Snapshot()
				if err != nil {
					return err
				}
				disk, err := store.Snapshot(root)
				if err != nil {
					return err
				}
				if !reflect.DeepEqual(got, disk) {
					t.Fatalf("confirmed differs from persisted format %d", version)
				}
				return nil
			}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestWorkItemSessionDoesNotRetainAfterUnlock(t *testing.T) {
	root := createProjectWorkspace(t)
	store := newSystemWorkItemStore()
	var escaped WorkItemStoreSession
	if err := store.WithLock(root, func(s WorkItemStoreSession) error { escaped = s; return s.Publish(minimalEpicRegistry(root)) }); err != nil {
		t.Fatal(err)
	}
	concrete := escaped.(*systemWorkItemStoreSession)
	if concrete.confirmed != nil {
		t.Fatal("generation retained after unlock")
	}
	if err := escaped.Publish(minimalEpicRegistry(root)); err == nil {
		t.Fatal("publication without lock succeeded")
	}
}

func TestWorkItemSessionDoesNotConfirmReplacementDuringObservation(t *testing.T) {
	for _, during := range []string{"read-closure", "post-rename"} {
		t.Run(during, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			goalQueueFixture(t, f, []TaskGoalQueueEntry{})
			root := filepath.Dir(f.wrapper)
			r, err := f.dependencies.WorkItems.Snapshot(root)
			if err != nil {
				t.Fatal(err)
			}
			r.Epics[0].Title = "Replacement"
			replacement, err := encodeWorkItemRegistry(r)
			if err != nil {
				t.Fatal(err)
			}
			replaced := false
			armed := false
			store := &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(op string) error {
				match := during == "read-closure" && op == "managed-read" || during == "post-rename" && op == "directory-sync"
				if armed && match && !replaced {
					replaced = true
					temp := filepath.Join(root, MarkerDirectory, "replacement")
					if err := os.WriteFile(temp, replacement, 0600); err != nil {
						return err
					}
					return os.Rename(temp, workItemsPath(root))
				}
				return nil
			}}}
			if err := store.WithLock(root, func(s WorkItemStoreSession) error {
				if during == "post-rename" {
					r, err = s.Snapshot()
					if err != nil {
						return err
					}
				}
				armed = true
				if during == "read-closure" {
					_, err = s.Snapshot()
				} else {
					err = s.Publish(r)
				}
				if err != nil {
					return err
				}
				if !replaced {
					t.Fatal("replacement hook did not run")
				}
				if s.(*systemWorkItemStoreSession).confirmed != nil {
					t.Fatal("paired stale registry with replacement metadata")
				}
				got, err := s.Snapshot()
				if err != nil {
					return err
				}
				if got.Epics[0].Title != "Replacement" {
					t.Fatal("reused stale registry")
				}
				return nil
			}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestWorkItemSessionPublishedArrayOwnership(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	g := goalFixture(t, f, "task", "array/goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &g}})
	root := filepath.Dir(f.wrapper)
	if err := f.dependencies.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
		r, err := s.Snapshot()
		if err != nil {
			return err
		}
		entry := r.TaskQueueEvents[0].Request["entries"].([]any)[0]
		r.TaskQueueEvents[0].Request["entries"] = [1]any{entry}
		// A Go array of reference-bearing interface values is accepted by the
		// existing YAML/JSON queue schema and publication transition checks.
		if err = s.Publish(r); err != nil {
			return err
		}
		before, err := s.Snapshot()
		if err != nil {
			return err
		}
		want, err := encodeWorkItemRegistry(before)
		if err != nil {
			return err
		}
		entry.(map[string]any)["task_id"] = "caller mutation"
		got, err := s.Snapshot()
		if err != nil {
			return err
		}
		actual, err := encodeWorkItemRegistry(got)
		if err != nil {
			return err
		}
		if string(want) != string(actual) {
			t.Fatal("published array references alias confirmed state")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestWorkItemSessionDetectsInPlaceDriftWithRestoredMtime(t *testing.T) {
	root := createProjectWorkspace(t)
	store := newSystemWorkItemStore()
	if err := store.WithLock(root, func(s WorkItemStoreSession) error {
		r := minimalEpicRegistry(root)
		r.Epics[0].Title = "Original"
		if err := s.Publish(r); err != nil {
			return err
		}
		info, err := os.Stat(workItemsPath(root))
		if err != nil {
			return err
		}
		r.Epics[0].Title = "Replaced"
		b, err := encodeWorkItemRegistry(r)
		if err != nil {
			return err
		}
		if err = os.WriteFile(workItemsPath(root), b, 0644); err != nil {
			return err
		}
		if err = os.Chtimes(workItemsPath(root), info.ModTime(), info.ModTime()); err != nil {
			return err
		}
		got, err := s.Snapshot()
		if err != nil {
			return err
		}
		if got.Epics[0].Title != "Replaced" {
			t.Fatal("missed same-size rewrite with restored modification time")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestWorkItemSessionPreservesQueueTimestampScalars(t *testing.T) {
	for _, boundary := range []string{"snapshot", "publication"} {
		t.Run(boundary, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			queueSetFixture(t, f, "timestamp/queue", []QueueEntry{})
			root := filepath.Dir(f.wrapper)
			path := workItemsPath(root)
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			// YAML accepts an unquoted timestamp as time.Time in a request interface.
			// Its canonical JSON value and request digest remain the same UTC string.
			quoted := []byte(`decided_at_utc: "2026-09-30T00:00:00Z"`)
			if bytes.Count(raw, quoted) != 1 {
				t.Fatal("expected one queue decision timestamp")
			}
			raw = bytes.Replace(raw, quoted, []byte("decided_at_utc: 2026-09-30T00:00:00Z"), 1)
			if err = os.WriteFile(path, raw, 0600); err != nil {
				t.Fatal(err)
			}
			if err = f.dependencies.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
				first, err := s.Snapshot()
				if err != nil {
					return err
				}
				request := first.TaskQueueEvents[0].Request
				decision := request["human_decision"].(map[string]any)
				instant, ok := decision["decided_at_utc"].(time.Time)
				if !ok || instant.IsZero() {
					t.Fatal("fixture did not consume an accepted timestamp scalar")
				}
				want := first.TaskQueueEvents[0].RequestSHA256
				if queueDigest(request) != want {
					t.Fatal("fixture changed semantic request digest")
				}
				if boundary == "publication" {
					if err = s.Publish(first); err != nil {
						return err
					}
				}
				// Changing the caller's nested map and its timestamp value cannot leak.
				decision["decided_at_utc"] = instant.Add(time.Hour)
				got, err := s.Snapshot()
				if err != nil {
					return err
				}
				if queueDigest(got.TaskQueueEvents[0].Request) != want {
					t.Fatal("confirmed snapshot changed accepted timestamp scalar or retained caller aliases")
				}
				if err = s.Publish(got); err != nil {
					return err
				}
				disk, err := f.dependencies.WorkItems.Snapshot(root)
				if err != nil {
					return err
				}
				confirmed, err := s.Snapshot()
				if err != nil {
					return err
				}
				if queueDigest(confirmed.TaskQueueEvents[0].Request) != queueDigest(disk.TaskQueueEvents[0].Request) || confirmed.RawSHA256 != disk.RawSHA256 {
					t.Fatal("confirmed queue timestamp or hash differs from persisted request")
				}
				return nil
			}); err != nil {
				t.Fatal(err)
			}
		})
	}
}
