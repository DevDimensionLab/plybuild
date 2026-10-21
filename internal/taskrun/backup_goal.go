package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/pelletier/go-toml/v2"
)

type BackupGoalInput struct {
	TaskID   workspace.TaskID `json:"task_id"`
	SpecID   string           `json:"spec_id"`
	Provider string           `json:"provider"`
	Model    string           `json:"model"`
	Effort   string           `json:"effort"`
}

// WorkflowBackupGoal is registration only. Preview never acquires a writer
// lock or touches the source. Apply holds the existing Run store lock while
// validating retirement and publishing through native workspace seams.
func WorkflowBackupGoal(d Dependencies, root, id string, in BackupGoalInput, apply bool, confirmation string) (out workspace.BackupGoalRegistrationResult, err error) {
	if err = containing(d, root); err != nil {
		return out, err
	}
	action := func() error {
		assignment, executorSource, e := backupExecutor(in)
		if e != nil {
			return e
		}
		in.Model, in.Effort = *assignment.Model, *assignment.Effort

		s, e := workflowRead(root, id)
		if e != nil {
			return e
		}
		if !deliveryRun(s.Request) || headlessRun(s.Request) || s.Request.Delivery == nil || s.Result.Delivery == nil {
			return workflowError(4, "backup_goal_unsupported: requires a native interactive DeliveryGoal Run")
		}
		retired, e := deliveryRetirement(s)
		if e != nil {
			return e
		}
		if retired == nil {
			return workflowError(4, "backup_goal_source_not_retired: use a confirmed natively retired owner; this operation cannot retire or stop it")
		}
		rounds, e := deliveryReviewRounds(s)
		if e != nil {
			return e
		}
		if len(rounds) != 0 {
			return workflowError(4, "backup_goal_review_unsupported: review quota is later work")
		}
		dl := s.Result.Delivery
		if dl.Phase == "completed" || dl.OwnershipRelease != nil || dl.HumanIntegration != nil || dl.Attempt != nil && dl.Attempt.State != "recorded" {
			return workflowError(4, "backup_goal_source_unresolved: delivered, released or unresolved attempt is unsupported")
		}
		if len(dl.Candidates) != 0 {
			return workflowError(4, "backup_goal_post_work_unsupported: candidate/review transfer requires later work")
		}
		stops, e := deliveryUnresolvedStops(s)
		if e != nil {
			return e
		}
		for _, stop := range stops {
			if stop.Input.EffectStatus == DeliveryEffectUnknown || stop.Input.EffectStatus == DeliveryEffectPartial {
				ok, e := deliveryQuotaAllowsRetirement(s, stop)
				if e != nil {
					return e
				}
				if !ok {
					return workflowError(4, "backup_goal_source_unresolved: unknown or partial source effects")
				}
			}
		}
		x, e := deliveryTarget(d, s)
		if e != nil {
			return e
		}
		if x.OID != s.Observed.Target.OID || x.Tree != s.Observed.Target.Tree {
			return workflowError(4, "backup_goal_post_work_unsupported: source changed; no patch transfer is supported")
		}
		raw, e := workflowBound(s.Request.Delivery.Goal, 4<<20)
		if e != nil {
			return e
		}
		var goal struct {
			Kind          string                       `json:"kind"`
			SchemaVersion int                          `json:"schema_version"`
			ContractKind  string                       `json:"contract_kind"`
			TaskID        workspace.TaskID             `json:"task_id"`
			SpecID        string                       `json:"spec_id"`
			Revision      int                          `json:"revision"`
			Delivery      *workspace.DeliveryAgreement `json:"delivery"`
		}
		if json.Unmarshal(raw, &goal) != nil || goal.Kind != "WorkspaceTaskSpecRevision@2" || goal.SchemaVersion != 2 || goal.ContractKind != "goal" {
			return workflowError(4, "backup_goal_unsupported: requires Goal Spec@2; legacy/report contracts remain unsupported")
		}
		a := s.Request.Delivery.Agreement
		if a == nil {
			return workflowError(4, "backup_goal_unsupported: missing delivery agreement")
		}
		unbound := *a
		unbound.SourceRef = ""
		if !equal(&unbound, goal.Delivery) || a.Mode != workspace.DeliveryLocalEpic || !a.AutomaticAcceptance() {
			return workflowError(4, "backup_goal_unsupported: requires automatic local Epic delivery")
		}
		source := workspace.TaskGoalRef{SpecID: goal.SpecID, Spec: workspace.TaskRevisionRef{Revision: goal.Revision, ManifestSHA256: s.Request.Delivery.Goal.SHA256}}
		provenance := fmt.Sprintf("Copied from retired Run %s request %s, Task %s, Goal %s@%d %s, source ref %s. Retirement is historical prerequisite, not new runtime authority.", id, s.Result.RequestSHA256, goal.TaskID, goal.SpecID, goal.Revision, source.Spec.ManifestSHA256, s.Observed.Target.Ref)
		out, e = workspace.RegisterBackupGoal(d.Workspace, workspace.BackupGoalRegistration{SourcePreparation: s.Request.PreparationID, SourceTask: goal.TaskID, SourceGoal: source, Target: workspace.QueueTargetInput{ProjectID: a.ProjectID, RepoID: a.RepoID, EpicID: a.EpicID}, TaskID: in.TaskID, SpecID: in.SpecID, Provider: in.Provider, Model: in.Model, Effort: in.Effort, ExecutorSource: executorSource, Provenance: provenance, RecordedAtUTC: retired.RecordedAtUTC, PublicationKey: "backup-goal/" + id}, apply, confirmation)
		return e
	}
	if apply {
		err = withStore(root, action)
	} else {
		err = action()
	}
	return out, err
}

func backupExecutor(in BackupGoalInput) (workspace.TaskExecutorAssignment, string, error) {
	var settings map[string]any
	source := "native Claude defaults"
	if in.Provider == "codex" {
		home := os.Getenv("CODEX_HOME")
		if home == "" {
			userHome, e := os.UserHomeDir()
			if e != nil {
				return workspace.TaskExecutorAssignment{}, "", e
			}
			home = filepath.Join(userHome, ".codex")
		}
		path := filepath.Join(home, "config.toml")
		raw, e := os.ReadFile(path)
		if e != nil {
			return workspace.TaskExecutorAssignment{}, "", fmt.Errorf("read existing Codex settings: %w", e)
		}
		if len(raw) > 4<<20 {
			return workspace.TaskExecutorAssignment{}, "", fmt.Errorf("Codex settings exceed 4 MiB")
		}
		if e = toml.Unmarshal(raw, &settings); e != nil {
			return workspace.TaskExecutorAssignment{}, "", e
		}
		source = "native Codex settings " + path + " " + hash(raw)
	}
	model, effort := workspace.TaskExecutorDefaults(in.Provider, settings)
	if in.Model != "" {
		model = in.Model
		source += "; explicit model override"
	}
	if in.Effort != "" {
		effort = in.Effort
		source += "; explicit effort override"
	}
	if e := workspace.ValidateTaskExecutorChoice(in.Provider, model, effort); e != nil {
		return workspace.TaskExecutorAssignment{}, "", e
	}
	return workspace.TaskExecutorAssignment{Provider: &in.Provider, Model: &model, Effort: &effort}, source, nil
}
