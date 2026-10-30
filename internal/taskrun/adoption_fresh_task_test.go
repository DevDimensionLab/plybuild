package taskrun

import (
	"fmt"
	"os"

	"path/filepath"
	"reflect"

	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Reject only the execution Spec's registry replacement, after the native
// adoption intent and immutable content writes. No registry bytes are patched.
type unstartedAdoptionStore struct {
	workspace.WorkItemStore
	specs    int
	rejected bool
}
type unstartedAdoptionSession struct {
	workspace.WorkItemStoreSession
	owner *unstartedAdoptionStore
}

func (s *unstartedAdoptionStore) WithLock(root string, fn func(workspace.WorkItemStoreSession) error) error {
	return s.WorkItemStore.WithLock(root, func(inner workspace.WorkItemStoreSession) error { return fn(unstartedAdoptionSession{inner, s}) })
}
func (s unstartedAdoptionSession) Publish(r workspace.WorkItemRegistry) error {
	if len(r.TaskSpecRevisions) > s.owner.specs {
		s.owner.rejected = true
		return fmt.Errorf("injected before execution content registry replacement")
	}
	return s.WorkItemStoreSession.Publish(r)
}

type freshAdoptionDraft struct {
	preservedPatch                                                   string
	source, sourceOID, oldHistory, patch, patchHash, factsHash, base string
	pending                                                          string
}

func fixtureFreshHuman() workspace.QueueHumanDecision {
	return workspace.QueueHumanDecision{ActorClaim: "modeled coordinator", DecidedAtUTC: "2026-10-05T00:00:00Z", Source: "explicit_human_instruction", Statement: "Isolated fixture only; no actual runtime acceptance or human QA."}
}
func (draft freshAdoptionDraft) assertPreserved(t *testing.T, w workspace.Dependencies) {
	t.Helper()
	old, err := w.WorkItems.Snapshot(filepath.Dir(draft.patch))
	if err != nil {
		t.Fatal(err)
	}
	if retainedAdoptionHistory(old) != draft.oldHistory {
		t.Fatal("old Task, Goal, intent or history changed")
	}
	if gitOutput(t, draft.source, "status", "--porcelain") != "?? external-dirty.txt" {
		t.Fatal("old tracked or dirty source changed")
	}
	if gitOutput(t, draft.source, "rev-parse", "HEAD") != draft.sourceOID {
		t.Fatal("old source ref changed")
	}
	dirty, err := os.ReadFile(filepath.Join(draft.source, "external-dirty.txt"))
	if err != nil || string(dirty) != "retain external writer bytes\n" {
		t.Fatal("external dirty bytes lost")
	}
	q, err := workspace.ListTaskQueue(w, workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, false)
	if err != nil {
		t.Fatal(err)
	}
	retained := q.Pending[:0]
	for _, entry := range q.Pending {
		if entry.TaskID != "fresh" {
			retained = append(retained, entry)
		}
	}
	if digest(retained) != draft.pending {
		t.Fatal("old/foreign pending queue changed")
	}
}

// Real native services and Git with injected writer/provider/Herdr observations;
// no actual old Task, provider, OS inventory or permission policy is exercised.

// Snapshot durable Task-owned records, excluding shared Epic observations and
// the fresh Task. Queue order is checked separately through native readback.
func retainedAdoptionHistory(r workspace.WorkItemRegistry) string {
	rows := map[string][]any{}
	value := reflect.ValueOf(r)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Kind() != reflect.Slice {
			continue
		}
		for j := 0; j < field.Len(); j++ {
			item := field.Index(j)
			if item.Kind() != reflect.Struct {
				continue
			}
			owner := item.FieldByName("TaskID")
			if value.Type().Field(i).Name == "Tasks" {
				owner = item.FieldByName("ID")
			}
			if owner.IsValid() && owner.Kind() == reflect.String && owner.String() != "fresh" {
				name := value.Type().Field(i).Name
				rows[name] = append(rows[name], item.Interface())
			}
		}
	}
	return digest(rows)
}
