package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
)

const projectMetadataStoreKind = "WorkspaceProjectMetadataStore@1"

type projectMetadataStore struct {
	Kind          string                 `json:"kind"`
	SchemaVersion int                    `json:"schema_version"`
	Events        []ProjectMetadataEvent `json:"events"`
}

func emptyProjectMetadata(root string) ProjectMetadataSnapshot {
	return ProjectMetadataSnapshot{Locator: filepath.Join(root, MarkerDirectory, ProjectMetadataFile), Events: []ProjectMetadataEvent{}, companions: map[ProjectID]map[string]ProjectCompanion{}, wrappers: map[RepoID]string{}}
}
func readProjectMetadata(root string) (ProjectMetadataSnapshot, error) {
	out := emptyProjectMetadata(root)
	bytes, exists, err := readWorkspaceMetadata(root, ProjectMetadataFile)
	if err != nil {
		return out, err
	}
	if !exists {
		return out, nil
	}
	var store projectMetadataStore
	if err := decodeWorkspaceMetadata(bytes, &store); err != nil {
		return out, metadataStoreError(out.Locator, err)
	}
	if store.Kind != projectMetadataStoreKind || store.SchemaVersion != 1 || store.Events == nil {
		return out, metadataStoreError(out.Locator, errors.New("unsupported or incomplete project metadata store; expected version 1 with an events array"))
	}
	if err := replayProjectMetadataEvents(&out, store.Events); err != nil {
		return out, metadataStoreError(out.Locator, err)
	}
	out.Exists, out.SHA256 = true, digestTaskBytes(bytes)
	return out, nil
}
func encodeProjectMetadata(events []ProjectMetadataEvent) ([]byte, error) {
	state := emptyProjectMetadata("")
	if err := replayProjectMetadataEvents(&state, events); err != nil {
		return nil, err
	}
	bytes, err := json.MarshalIndent(projectMetadataStore{projectMetadataStoreKind, 1, events}, "", "  ")
	return append(bytes, '\n'), err
}
func replayProjectMetadataEvents(state *ProjectMetadataSnapshot, events []ProjectMetadataEvent) error {
	for i, event := range events {
		if event.Sequence != i+1 || event.ID != projectMetadataEventID(event) || !validTaskText(event.ActorClaim, 1, 256) || !validTaskUTC(event.RecordedAtUTC) {
			return fmt.Errorf("invalid project metadata event at sequence %d", i+1)
		}
		changed, err := applyProjectMetadataEvent(state, event)
		if err != nil {
			return err
		}
		if !changed {
			return fmt.Errorf("project metadata event %d records no change", i+1)
		}
	}
	state.Events, state.Revision = events, len(events)
	return nil
}
func applyProjectMetadataEvent(state *ProjectMetadataSnapshot, event ProjectMetadataEvent) (bool, error) {
	if _, err := ParseProjectID(string(event.ProjectID)); err != nil {
		return false, err
	}
	switch event.Operation {
	case ProjectCompanionAdd, ProjectCompanionRemove:
		if event.RepoID != nil || event.Role == nil || event.Locator == nil || !validStoredPath(*event.Locator) {
			return false, errors.New("invalid companion metadata event")
		}
		role := *event.Role
		if role != CompanionPlanning && role != CompanionDocs && role != CompanionOther {
			return false, errors.New("invalid companion role")
		}
		if state.companions[event.ProjectID] == nil {
			state.companions[event.ProjectID] = map[string]ProjectCompanion{}
		}
		items := state.companions[event.ProjectID]
		key := string(role) + "\x00" + *event.Locator
		_, exists := items[key]
		if event.Operation == ProjectCompanionAdd {
			if exists {
				return false, nil
			}
			items[key] = ProjectCompanion{role, *event.Locator}
			return true, nil
		}
		if !exists {
			return false, nil
		}
		delete(items, key)
		return true, nil
	case ProjectRepoWrapperSet, ProjectRepoWrapperClear:
		if event.RepoID == nil || event.Role != nil {
			return false, errors.New("invalid repository wrapper event")
		}
		if err := validateIdentifier("repository", string(*event.RepoID)); err != nil {
			return false, err
		}
		if event.Operation == ProjectRepoWrapperClear {
			if event.Locator != nil {
				return false, errors.New("wrapper clear must have null locator")
			}
			if _, exists := state.wrappers[*event.RepoID]; !exists {
				return false, nil
			}
			delete(state.wrappers, *event.RepoID)
			return true, nil
		}
		if event.Locator == nil || !validStoredPath(*event.Locator) {
			return false, errors.New("invalid repository wrapper path")
		}
		if state.wrappers[*event.RepoID] == *event.Locator {
			return false, nil
		}
		state.wrappers[*event.RepoID] = *event.Locator
		return true, nil
	default:
		return false, errors.New("unknown project metadata operation")
	}
}
