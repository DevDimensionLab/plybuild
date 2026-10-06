package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
)

const workItemLifecycleStoreKind = "WorkspaceWorkItemLifecycleStore@1"

type workItemLifecycleStore struct {
	Kind          string                   `json:"kind"`
	SchemaVersion int                      `json:"schema_version"`
	Events        []WorkItemLifecycleEvent `json:"events"`
}

func readWorkItemLifecycle(root string) (WorkItemLifecycleSnapshot, error) {
	out := WorkItemLifecycleSnapshot{Locator: lifecycleLocator(root), Events: []WorkItemLifecycleEvent{}, states: map[string]LifecycleState{}}
	bytes, exists, err := readWorkspaceMetadata(root, WorkItemLifecycleFile)
	if err != nil {
		return out, err
	}
	if !exists {
		return out, nil
	}
	var store workItemLifecycleStore
	if err := decodeWorkspaceMetadata(bytes, &store); err != nil {
		return out, metadataStoreError(out.Locator, err)
	}
	if store.Kind != workItemLifecycleStoreKind || store.SchemaVersion != 1 || store.Events == nil {
		return out, metadataStoreError(out.Locator, errors.New("unsupported or incomplete lifecycle store; expected version 1 with an events array"))
	}
	if err := validateLifecycleEvents(store.Events); err != nil {
		return out, metadataStoreError(out.Locator, err)
	}
	out.Events, out.Revision, out.Exists, out.SHA256 = store.Events, len(store.Events), true, digestTaskBytes(bytes)
	for _, event := range store.Events {
		out.states[event.SubjectKind+"/"+event.SubjectID] = event.To
	}
	return out, nil
}

func encodeWorkItemLifecycle(events []WorkItemLifecycleEvent) ([]byte, error) {
	if err := validateLifecycleEvents(events); err != nil {
		return nil, err
	}
	bytes, err := json.MarshalIndent(workItemLifecycleStore{workItemLifecycleStoreKind, 1, events}, "", "  ")
	return append(bytes, '\n'), err
}

func validateLifecycleEvents(events []WorkItemLifecycleEvent) error {
	states := map[string]LifecycleState{}
	for i, event := range events {
		if err := validateLifecycleSubject(event.SubjectKind, event.SubjectID); err != nil {
			return err
		}
		key := event.SubjectKind + "/" + event.SubjectID
		previous, ok := states[key]
		if !ok {
			previous = LifecycleActive
		}
		if event.Sequence != i+1 || event.ID != lifecycleEventID(event) || event.From != previous || event.From == event.To || event.To.Validate() != nil || !validTaskText(event.ActorClaim, 1, 256) || !validTaskUTC(event.RecordedAtUTC) {
			return fmt.Errorf("invalid lifecycle event at sequence %d", i+1)
		}
		states[key] = event.To
	}
	return nil
}
