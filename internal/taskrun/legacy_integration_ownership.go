package taskrun

import (
	"fmt"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type LegacyTaskOwnerObservation struct {
	RunID        string                       `json:"run_id"`
	State        string                       `json:"state"`
	Evidence     []FileBinding                `json:"evidence"`
	ProviderExit *DeliveryStartupExitEvidence `json:"provider_exit,omitempty"`
}

type LegacyTaskOwnershipObservation struct {
	Kind          string                       `json:"kind"`
	SchemaVersion int                          `json:"schema_version"`
	TaskID        workspace.TaskID             `json:"task_id"`
	TaskResultID  workspace.TaskResultID       `json:"task_result_id"`
	ResultOID     string                       `json:"result_oid"`
	ResultTree    string                       `json:"result_tree"`
	Ready         bool                         `json:"ready"`
	Reason        string                       `json:"reason"`
	Owners        []LegacyTaskOwnerObservation `json:"owners"`
}

// CheckLegacyTaskOwnership detects known native writers for a human's explicit
// legacy ownership release. Absence is reported literally: it is not evidence
// that an unrecorded process exited, and is not by itself release authority.
func CheckLegacyTaskOwnership(d Dependencies, root string, result workspace.TaskResultRecord) (LegacyTaskOwnershipObservation, error) {
	out := LegacyTaskOwnershipObservation{Kind: "PlyLegacyTaskOwnershipObservation@1", SchemaVersion: 1, TaskID: result.TaskID, TaskResultID: result.ID, ResultOID: result.ResultOID, ResultTree: result.ResultTree, Reason: "owner_unknown", Owners: []LegacyTaskOwnerObservation{}}
	registry, e := d.Workspace.WorkItems.Snapshot(root)
	if e != nil {
		return out, e
	}
	found := false
	for _, r := range registry.TaskResults {
		if equal(r, result) {
			found = true
		}
	}
	if !found {
		return out, workflowError(4, "owner_unknown: selected legacy TaskResult is not the exact registered record")
	}
	inventory, e := ReadInventory(d, root)
	if e != nil {
		return out, e
	}
	if len(inventory.Reasons) > 0 {
		return out, workflowError(4, "owner_unknown: native owner inventory is incomplete; inspect its recorded reasons")
	}
	for _, row := range inventory.Runs {
		if row.TaskID == nil {
			return out, workflowError(4, "owner_unknown: a preserved run cannot be assigned to its Task")
		}
		if *row.TaskID != string(result.TaskID) {
			continue
		}
		if row.Freshness != "fresh" {
			return out, workflowError(4, fmt.Sprintf("owner_unknown: preserved Task owner evidence is incomplete: %+v", row.Reasons))
		}
		observation := LegacyTaskOwnerObservation{RunID: row.RunID, Evidence: append([]FileBinding{}, row.Sources...)}
		if workflowRunIDPattern.MatchString(row.RunID) {
			s, e := workflowRead(root, row.RunID)
			if e != nil {
				return out, e
			}
			if s.Observed.Target.WorktreeLocator != result.SourceLocator || s.Observed.Target.Ref != result.SourceRef || s.Observed.Target.GitCommonDir != result.GitCommonDir {
				return out, workflowError(4, "owner_unknown: native Task owner source identity differs")
			}
			if deliveryRun(s.Request) && s.Result.Delivery != nil {
				if attempt := s.Result.Delivery.Attempt; attempt != nil && attempt.State == "attempted" {
					return out, workflowError(4, "owner_active: a native delivery effect is unresolved")
				}
				if len(s.Result.Delivery.Candidates) > 0 {
					c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
					release, e := deliveryRelease(s, c)
					if e != nil {
						return out, e
					}
					if release != nil {
						if c.TaskResult.ID != result.ID {
							return out, workflowError(4, "owner_unknown: release binds another candidate")
						}
						observation.State = "explicitly_released"
						observation.Evidence = append(observation.Evidence, *release)
						out.Owners = append(out.Owners, observation)
						continue
					}
				}
			}
			exit, e := deliveryRecoveryExit(d, s)
			if e != nil {
				return out, workflowError(4, "owner_active: native Task owner has no explicit release or positively observed process exit")
			}
			observation.State, observation.ProviderExit = "provider_exited", &exit
		} else {
			j, e := readJournal(root, row.RunID)
			if e != nil {
				return out, e
			}
			if e = validateRunEvidence(d, j); e != nil {
				return out, e
			}
			if !quiescent(j.Result.Process) || j.Result.Process.State != "exited" {
				return out, workflowError(4, "owner_active: native Task process has not recorded a quiescent exit")
			}
			observation.State = "recorded_quiescent_exit"
			for i, ev := range j.Events {
				if ev.Type == "process_exited" {
					observation.Evidence = append(observation.Evidence, FileBinding{filepath.Join(runPaths(j.Request).RunRoot, "events", fmt.Sprintf("%06d.json", ev.Sequence)), j.Hashes[i]})
				}
			}
		}
		out.Owners = append(out.Owners, observation)
	}
	out.Ready, out.Reason = true, "native_owners_released_or_exited"
	if len(out.Owners) == 0 {
		out.Reason = "no_native_owner_recorded"
	}
	return out, nil
}
