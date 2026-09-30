package workspace

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func TestTaskContentConcurrentMigrationAcceptsOneExactRegistryAcknowledgement(t *testing.T) {
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
			before, err := os.ReadFile(workItemsPath(root))
			if err != nil {
				t.Fatal(err)
			}
			start := make(chan struct{})
			type outcome struct {
				result TaskContentMutationResult
				err    error
			}
			done := make(chan outcome, 2)
			for _, key := range []string{"fixture/migrate-a", "fixture/migrate-b"} {
				in := contentFixtureProblem(t, f, key)
				raw, err := os.ReadFile(in.File)
				if err != nil {
					t.Fatal(err)
				}
				value, err := canonicaljson.DecodeStrict(raw)
				if err != nil {
					t.Fatal(err)
				}
				m := contentFields(value)
				m["registry_upgrade"] = contentRefValue(TaskRegistryUpgrade{FromVersion: version, RegistrySHA256: digestTaskBytes(before)})
				raw, err = canonicaljson.Marshal(contentObject(m))
				if err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(in.File, raw, 0600); err != nil {
					t.Fatal(err)
				}
				go func(in TaskContentInput) {
					<-start
					result, err := RecordTaskProblem(f.dependencies, in)
					done <- outcome{result, err}
				}(in)
			}
			close(start)
			wins := 0
			for i := 0; i < 2; i++ {
				out := <-done
				if out.err == nil {
					wins++
					continue
				}
				var contentErr *TaskContentError
				if !errors.As(out.err, &contentErr) || contentErr.Code != "task_content_upgrade_conflict" || out.result.Attempted.ObjectPublicationStarted {
					t.Fatalf("losing migration did not stop before publication: %#v %v", out.result, out.err)
				}
			}
			r, err := f.dependencies.WorkItems.Snapshot(root)
			if err != nil || wins != 1 || r.FormatVersion != 3 || len(r.TaskProblemRevisions) != 1 || len(r.TaskContentPublications) != 1 {
				t.Fatalf("migration did not have one exact winner: wins=%d error=%v", wins, err)
			}
			backup, err := f.dependencies.TaskContent.Read(root, "backups", digestTaskBytes(before))
			if err != nil || !bytes.Equal(before, backup) {
				t.Fatal("competing migration changed the raw backup")
			}
			legacy, err := decodeWorkItemRegistry(before)
			if err != nil || !contentTypedEqual(r.Tasks, legacy.Tasks) || !contentTypedEqual(r.Epics, legacy.Epics) || !contentTypedEqual(r.WorktreeOperations, legacy.WorktreeOperations) {
				t.Fatal("competing migration changed legacy records")
			}
		})
	}
}

func TestTaskContentConcurrentCASPublishesExactlyOneRevision(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	a := contentFixtureProblem(t, f, "fixture/a")
	b := contentFixtureProblem(t, f, "fixture/b")
	start := make(chan struct{})
	done := make(chan error, 2)
	for _, in := range []TaskContentInput{a, b} {
		go func(in TaskContentInput) { <-start; _, err := RecordTaskProblem(f.dependencies, in); done <- err }(in)
	}
	close(start)
	wins := 0
	for i := 0; i < 2; i++ {
		if <-done == nil {
			wins++
		}
	}
	r, err := f.dependencies.WorkItems.Snapshot(filepath.Dir(f.wrapper))
	if err != nil || wins != 1 || len(r.TaskProblemRevisions) != 2 {
		t.Fatalf("CAS wins=%d registry=%#v error=%v", wins, r, err)
	}
}

func TestTaskSpecGuardUsesExistingWriterLocksInBothOrders(t *testing.T) {
	for _, writerFirst := range []bool{false, true} {
		t.Run(map[bool]string{true: "problem-before-start", false: "start-before-problem"}[writerFirst], func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			root := filepath.Dir(f.wrapper)
			path := filepath.Join(f.wrapper, "task")
			if _, err := CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: "task", Branch: "task", Path: path, ExpectedParentOID: f.oid}); err != nil {
				t.Fatal(err)
			}
			basis := fixtureSelectTaskSpec(t, f)
			in := contentFixtureProblem(t, f, "fixture/p2")
			entered, release := make(chan struct{}), make(chan struct{})
			first, second := make(chan error, 1), make(chan error, 1)
			guard := func() error {
				return WithTaskSpecSnapshot(f.dependencies, root, func(s *TaskSpecSession) error {
					_, err := s.ValidateTarget(path, "refs/heads/task", s.Registry.Tasks[0].GitCommonDir, &basis)
					return err
				})
			}
			if writerFirst {
				f.dependencies.TaskContent = &TaskContentStorage{fault: func(point string) error {
					if point == "before-registry-replace" {
						close(entered)
						<-release
					}
					return nil
				}}
				go func() { _, err := RecordTaskProblem(f.dependencies, in); first <- err }()
				<-entered
				go func() { second <- guard() }()
			} else {
				go func() {
					first <- WithTaskSpecSnapshot(f.dependencies, root, func(s *TaskSpecSession) error {
						close(entered)
						<-release
						_, err := s.ValidateTarget(path, "refs/heads/task", s.Registry.Tasks[0].GitCommonDir, &basis)
						return err
					})
				}()
				<-entered
				go func() { _, err := RecordTaskProblem(f.dependencies, in); second <- err }()
			}
			select {
			case err := <-second:
				t.Fatalf("operation bypassed existing advisory lock: %v", err)
			case <-time.After(40 * time.Millisecond):
			}
			close(release)
			if err := <-first; err != nil {
				t.Fatal(err)
			}
			err := <-second
			if writerFirst && err == nil {
				t.Fatal("old selection was accepted after P2")
			}
			if !writerFirst && err != nil {
				t.Fatal(err)
			}
		})
	}
}
