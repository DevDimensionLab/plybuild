package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"time"
)

type TaskVerificationInput struct {
	TaskID      TaskID
	ReceiptPath string
}

type taskVerificationReceipt struct {
	Kind          string                  `json:"kind"`
	SchemaVersion int                     `json:"schema_version"`
	RunID         string                  `json:"run_id"`
	AttemptID     string                  `json:"attempt_id"`
	CandidateOID  string                  `json:"candidate_oid"`
	CandidateTree string                  `json:"candidate_tree"`
	Argv          []string                `json:"argv"`
	CWD           string                  `json:"cwd"`
	Exit          *int                    `json:"exit"`
	Stdout        taskVerificationBinding `json:"stdout"`
	Stderr        taskVerificationBinding `json:"stderr"`
	StartedAt     string                  `json:"started_at"`
	FinishedAt    string                  `json:"finished_at"`
	Error         string                  `json:"error"`
	Automatic     *struct {
		ExecutedArgv []string `json:"executed_argv"`
	} `json:"automatic"`
}

func (r taskVerificationReceipt) executedArgv() ([]string, error) {
	var argv []string
	switch {
	case r.Kind == "PlyDeliveryVerification@1" && r.SchemaVersion == 1 && r.Automatic == nil:
		argv = r.Argv
	case r.Kind == "PlyDeliveryVerification@2" && r.SchemaVersion == 2 && r.Automatic != nil:
		argv = r.Automatic.ExecutedArgv
	default:
		return nil, fmt.Errorf("unsupported verification receipt")
	}
	if len(argv) == 0 || argv[0] == "" {
		return nil, fmt.Errorf("verification receipt lacks the executed command")
	}
	return argv, nil
}

type taskVerificationBinding struct {
	Locator string `json:"locator"`
	SHA256  string `json:"sha256"`
}

func verificationArtifact(d Dependencies, path, id, role string) (TaskArtifactRecord, []byte, error) {
	if !filepath.IsAbs(path) {
		return TaskArtifactRecord{}, nil, fmt.Errorf("verification artifact must have an absolute path")
	}
	info, err := d.Files.Lstat(path)
	if err != nil {
		return TaskArtifactRecord{}, nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > 64<<20 {
		return TaskArtifactRecord{}, nil, fmt.Errorf("verification artifact must be a bounded regular file")
	}
	raw, err := d.Files.ReadFile(path)
	if err != nil {
		return TaskArtifactRecord{}, nil, err
	}
	sum := sha256.Sum256(raw)
	return TaskArtifactRecord{ArtifactID: id, Role: role, Locator: path, SHA256: "sha256:" + hex.EncodeToString(sum[:]), SizeBytes: int64(len(raw))}, raw, nil
}

func RecordTaskVerification(d Dependencies, in TaskVerificationInput) (TaskResultRecord, error) {
	d = WithTaskContentScope(d, in.TaskID)
	var out TaskResultRecord
	if _, err := ParseTaskID(string(in.TaskID)); err != nil {
		return out, err
	}
	if err := requireTaskDeliveryDependencies(d, false); err != nil {
		return out, err
	}
	if d.WorkGit == nil {
		return out, fmt.Errorf("verification requires Git observations")
	}
	receipt, raw, err := verificationArtifact(d, in.ReceiptPath, "verification", "command_receipt")
	if err != nil {
		return out, err
	}
	var observed taskVerificationReceipt
	if err = json.Unmarshal(raw, &observed); err != nil {
		return out, err
	}
	argv, err := observed.executedArgv()
	if err != nil {
		return out, err
	}
	started, e1 := time.Parse(time.RFC3339Nano, observed.StartedAt)
	finished, e2 := time.Parse(time.RFC3339Nano, observed.FinishedAt)
	if observed.RunID == "" || observed.AttemptID == "" || observed.Exit == nil || observed.Error != "" || e1 != nil || e2 != nil || finished.Before(started) {
		return out, fmt.Errorf("verification receipt lacks an observed command completion")
	}
	stdout, _, err := verificationArtifact(d, observed.Stdout.Locator, "stdout", "command_stdout")
	if err != nil {
		return out, err
	}
	stderr, _, err := verificationArtifact(d, observed.Stderr.Locator, "stderr", "command_stderr")
	if err != nil {
		return out, err
	}
	if stdout.SHA256 != observed.Stdout.SHA256 || stderr.SHA256 != observed.Stderr.SHA256 {
		return out, fmt.Errorf("verification output differs from the command receipt")
	}
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return out, err
	}
	err = d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			registry, err := session.Snapshot()
			if err != nil {
				return err
			}
			publication := "verification/" + observed.RunID + "/" + observed.AttemptID
			for _, prior := range registry.TaskResults {
				if prior.PublicationKey == publication {
					if prior.TaskID != in.TaskID || prior.DraftSHA256 != receipt.SHA256 {
						return fmt.Errorf("verification operation already records different bytes")
					}
					if err := ValidateTaskVerificationReceipt(d, prior); err != nil {
						return err
					}
					out = prior
					return nil
				}
			}
			task, _ := findTask(registry, in.TaskID)
			if task == nil || task.Worktree == nil || task.WorktreeState != WorkItemReady {
				return fmt.Errorf("verification requires an existing Task worktree")
			}
			_, repo, err := projectAndRepo(projects, task.ProjectID, task.RepoID)
			if err != nil {
				return err
			}
			source, err := d.WorkGit.ObserveWorktree(task.Worktree.Locator)
			if err != nil {
				return err
			}
			if observed.CWD != task.Worktree.Locator || source.Ref != task.Worktree.Ref || source.GitCommonDir != task.GitCommonDir || source.GitCommonDir != repo.GitCommonDir || source.OID != observed.CandidateOID || source.Tree != observed.CandidateTree || !source.Clean || !source.InventoryMatch {
				return fmt.Errorf("verification differs from the clean exact registered Task source")
			}
			id, err := d.TaskLifecycleIDs.NewTaskResultID()
			if err != nil {
				return err
			}
			outcome := "failed"
			if *observed.Exit == 0 {
				outcome = "passed"
			}
			stdoutID, stderrID := stdout.ArtifactID, stderr.ArtifactID
			transition := "none"
			if registry.FormatVersion == 1 {
				transition = "format_1_to_2"
				upgradeRegistryToV2(&registry)
			}
			out = TaskResultRecord{ID: id, PublicationKey: publication, DraftSHA256: receipt.SHA256, StoreTransition: transition, TaskID: task.ID, TaskWorktreeID: task.Worktree.ID, ProjectID: task.ProjectID, RepoID: task.RepoID, GitCommonDir: source.GitCommonDir, SourceLocator: task.Worktree.Locator, SourceRef: source.Ref, ResultOID: source.OID, ResultTree: source.Tree, RunID: observed.RunID, ReportedOutcome: "observed", VerificationReceipt: &receipt, VerifierResults: []TaskVerifierResultRecord{{VerifierID: observed.AttemptID, Outcome: outcome, Argv: argv, CWD: observed.CWD, Exit: *observed.Exit, BoundOIDOrSHA256: source.OID, StdoutArtifactID: &stdoutID, StderrArtifactID: &stderrID}}, Artifacts: []TaskArtifactRecord{stderr, stdout}, Recorder: TaskRecorderRecord{ActorClaim: observed.RunID, ControlSurface: "explicit verification", RecordedAtUTC: observed.FinishedAt}}
			registry.TaskResults = append(registry.TaskResults, out)
			sortWorkRegistry(&registry)
			_, err = publishWorkItemRegistryRecover(session, registry)
			return err
		})
	})
	return out, err
}

func ValidateTaskVerificationReceipt(d Dependencies, record TaskResultRecord) error {
	if record.VerificationReceipt == nil {
		return fmt.Errorf("verification receipt missing")
	}
	bound := record.VerificationReceipt
	actual, raw, err := verificationArtifact(d, bound.Locator, bound.ArtifactID, bound.Role)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(actual, *bound) || record.DraftSHA256 != actual.SHA256 {
		return fmt.Errorf("verification receipt bytes changed")
	}
	var receipt taskVerificationReceipt
	if err := json.Unmarshal(raw, &receipt); err != nil {
		return err
	}
	argv, err := receipt.executedArgv()
	if err != nil {
		return err
	}
	if record.ReportedOutcome != "observed" || record.TechnicalGate != "" || record.HandoffID != "" || len(record.Artifacts) != 2 {
		return fmt.Errorf("direct verification record has process claims or missing output")
	}
	if receipt.Exit == nil || receipt.Error != "" || receipt.RunID != record.RunID || receipt.CandidateOID != record.ResultOID || receipt.CandidateTree != record.ResultTree || receipt.CWD != record.SourceLocator || len(record.VerifierResults) != 1 {
		return fmt.Errorf("verification receipt differs from the recorded candidate")
	}
	result := record.VerifierResults[0]
	outcome := "failed"
	if *receipt.Exit == 0 {
		outcome = "passed"
	}
	if result.Outcome != outcome || result.VerifierID != receipt.AttemptID || result.StdoutArtifactID == nil || *result.StdoutArtifactID != "stdout" || result.StderrArtifactID == nil || *result.StderrArtifactID != "stderr" || result.Exit != *receipt.Exit || !reflect.DeepEqual(result.Argv, argv) || result.CWD != receipt.CWD || result.BoundOIDOrSHA256 != receipt.CandidateOID {
		return fmt.Errorf("verification result differs from the command receipt")
	}
	for _, artifact := range record.Artifacts {
		wanted := receipt.Stdout
		role := "command_stdout"
		if artifact.ArtifactID == "stderr" {
			wanted, role = receipt.Stderr, "command_stderr"
		} else if artifact.ArtifactID != "stdout" {
			return fmt.Errorf("unexpected verification artifact")
		}
		if artifact.Locator != wanted.Locator || artifact.SHA256 != wanted.SHA256 || artifact.Role != role {
			return fmt.Errorf("verification output binding differs from receipt")
		}
		actual, _, err := verificationArtifact(d, artifact.Locator, artifact.ArtifactID, artifact.Role)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(actual, artifact) {
			return fmt.Errorf("verification output bytes changed")
		}
	}
	return nil
}
