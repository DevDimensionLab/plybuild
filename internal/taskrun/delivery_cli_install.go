package taskrun

import (
	"context"
	gobuild "debug/buildinfo"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/buildinfo"
	"github.com/devdimensionlab/plybuild/internal/clirelease"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/sirupsen/logrus"
)

// CLIInstallation is coordinator evidence, separate from the closed owner Run.
// The retained native snapshot is also the supported first-install controller.
type CLIInstallation struct {
	Kind           string                  `json:"kind"`
	RunID          string                  `json:"run_id"`
	Coordinator    string                  `json:"coordinator"`
	Destination    string                  `json:"destination"`
	Controller     FileBinding             `json:"controller"`
	Source         clirelease.QualifiedCLI `json:"source"`
	Receipt        *clirelease.Receipt     `json:"receipt,omitempty"`
	ObservedSHA256 string                  `json:"observed_sha256,omitempty"`
	Smoke          *CLIInstallSmoke        `json:"smoke,omitempty"`
	State          string                  `json:"state"`
	NextAction     string                  `json:"next_action"`
}
type CLIInstallSmoke struct {
	Started bool   `json:"started"`
	Exit    *int   `json:"exit"`
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
	Error   string `json:"error,omitempty"`
}

func qualifiedCLIInstallation(d Dependencies, root, id, destination, coordinator string) (CLIInstallation, []byte, error) {
	out := CLIInstallation{Kind: "PlyQualifiedCLIInstallation@1", RunID: id, Destination: destination, Coordinator: coordinator, State: "ready", NextAction: "Run install-cli with the same destination and coordinator; retain this controller outside the Task."}
	if strings.TrimSpace(coordinator) != coordinator || coordinator == "" {
		return out, nil, workflowError(2, "--coordinator must name the actual authorized caller")
	}
	if !filepath.IsAbs(destination) || filepath.Clean(destination) != destination {
		return out, nil, workflowError(2, "--destination requires an explicit physical absolute Ply CLI path")
	}
	parent, err := filepath.EvalSymlinks(filepath.Dir(destination))
	if err != nil || parent != filepath.Dir(destination) {
		return out, nil, workflowError(4, "destination parent must be physical")
	}
	s, err := workflowRead(root, id)
	if err != nil {
		return out, nil, err
	}
	a := deliveryEffectiveAgreement(s)
	if a == nil || a.ProjectID != "ply" || a.RepoID != "ply" || a.EpicID == "" || a.Mode != workspace.DeliveryLocalEpic || !a.AutomaticAcceptance() || s.Result.Delivery == nil || s.Result.Delivery.Phase != "completed" || len(s.Result.Delivery.Candidates) == 0 {
		return out, nil, workflowError(4, "install-cli requires completed native self-hosted Ply local Epic delivery")
	}
	cwd, err := d.CWD()
	if err != nil || cwd != a.TargetWorktree {
		return out, nil, workflowError(4, "coordinator must run from the registered delivery target, outside the owner Task")
	}
	if _, err = deliveryAuthorityState(d, s); err != nil {
		return out, nil, err
	}
	c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
	integrated, handled, err := ObserveAutomaticLocalDelivery(d, root, id, c.TaskResult.ID)
	if err != nil {
		return out, nil, err
	}
	if !handled || !integrated.Completed || integrated.Closeout == nil || integrated.Closeout.State != "complete" || !integrated.Closeout.LifecycleCompleted {
		return out, nil, workflowError(4, "native delivery and closeout must be observed before coordinator installation")
	}
	registry, err := d.Workspace.WorkItems.Snapshot(root)
	if err != nil {
		return out, nil, err
	}
	found := false
	for _, r := range registry.TaskResults {
		if equal(r, c.TaskResult) {
			found = true
		}
	}
	if !found {
		return out, nil, workflowError(4, "qualified result is absent from native registry")
	}
	qa, err := LatestDeliveryAcceptanceHumanQA(registry.HumanQARecords, c.TaskResult, a)
	if err != nil {
		return out, nil, err
	}
	if err = requireDeliveryPass(DeliveryAuthority{Agreement: *a, Candidate: DeliveryCandidate{TaskResult: c.TaskResult, HumanQA: qa}}); err != nil {
		return out, nil, err
	}
	raw, err := workflowBound(c.Verification, 4<<20)
	if err != nil {
		return out, nil, err
	}
	var v deliveryVerificationReceipt
	if err = decode(raw, 4<<20, &v); err != nil {
		return out, nil, err
	}
	if v.RunID != id || v.RequestSHA256 != s.Result.RequestSHA256 || v.AttemptID != c.Key || v.CandidateOID != c.OID || v.CandidateTree != c.Tree {
		return out, nil, workflowError(4, "verification differs from delivered source")
	}
	if err = validateAutomaticReceipt(s, v, false); err != nil {
		return out, nil, err
	}
	// Qualification's passed review is native evidence. Re-read it after source
	// removal without requiring any live owner, old executable or working inputs.
	reviewRaw, err := workflowBound(v.Review, 1<<20)
	if err != nil {
		return out, nil, err
	}
	var review struct {
		Decision               string            `json:"decision"`
		CandidateOID           string            `json:"candidate_oid"`
		CandidateTree          string            `json:"candidate_tree"`
		ReviewerSessionID      *string           `json:"reviewer_session_id"`
		OpenActionableFindings []json.RawMessage `json:"open_actionable_findings"`
	}
	if err = json.Unmarshal(reviewRaw, &review); err != nil {
		return out, nil, err
	}
	if review.Decision != "passed" || review.CandidateOID != c.OID || review.CandidateTree != c.Tree || review.ReviewerSessionID == nil || *review.ReviewerSessionID == deliveryOwnerSession(s) || len(review.OpenActionableFindings) > 0 {
		return out, nil, workflowError(4, "independent passed exact-candidate review is required")
	}
	out.Controller = FileBinding{Locator: v.Automatic.CandidateBinarySnapshot.Locator, SHA256: v.Automatic.CandidateBinarySnapshot.SHA256}
	binary, err := workflowBound(out.Controller, workflowhandoff.AutomaticCandidateBinaryLimit)
	if err != nil {
		return out, nil, err
	}
	info, err := gobuild.ReadFile(out.Controller.Locator)
	if err != nil || info.Path != "github.com/devdimensionlab/plybuild/cmd/ply" && info.Path != "github.com/devdimensionlab/plybuild" {
		return out, nil, workflowError(4, "qualified candidate is not the self-hosted Ply CLI")
	}
	out.Source = clirelease.QualifiedCLI{RunID: id, Project: string(a.ProjectID), Repo: string(a.RepoID), Epic: string(a.EpicID), OID: c.OID, Tree: c.Tree, BinarySHA256: out.Controller.SHA256, VerificationSHA256: c.Verification.SHA256}
	return out, binary, nil
}

// InstallQualifiedCLI does not grant authority to the owner and never updates
// its acceptance, mandate, or closed state. The explicit caller owns this
// separate operation. check is read-only; observe only probes known bytes.
func InstallQualifiedCLI(d Dependencies, root, id, destination, coordinator, action string) (CLIInstallation, error) {
	out, binary, err := qualifiedCLIInstallation(d, root, id, destination, coordinator)
	if err != nil {
		out.State = "blocked"
		out.NextAction = "Inspect the actual validation error; no installation was started."
		return out, err
	}
	if action == "check" {
		return out, nil
	}
	sourceState, err := workflowRead(root, id)
	if err != nil {
		return out, err
	}
	executable, err := d.Executable()
	if err != nil {
		return out, err
	}
	controller, err := filepath.EvalSymlinks(executable)
	if err != nil {
		return out, err
	}
	if controller == destination || strings.HasPrefix(controller, sourceState.Observed.Target.WorktreeLocator+string(filepath.Separator)) {
		return out, workflowError(4, "preserve controller outside destination and removed Task")
	}
	if left, e := os.Stat(controller); e == nil {
		if right, e := os.Stat(destination); e == nil && os.SameFile(left, right) {
			return out, workflowError(4, "destination aliases the running controller")
		}
	}
	if action != "install" && action != "observe" && action != "rollback" {
		return out, workflowError(2, "invalid install action")
	}
	// One immutable operation per delivered Run. Changed destination or caller
	// cannot reuse an interrupted operation or overwrite its previous backup.
	path := filepath.Join(workflowRoot(root), "cli-installations", id)
	err = withStore(root, func() error {
		fresh, _, err := qualifiedCLIInstallation(d, root, id, destination, coordinator)
		if err != nil {
			return err
		}
		if !equal(fresh, out) {
			return workflowError(4, "installation source changed before effect")
		}
		if action == "install" {
			if err := d.writeValue(filepath.Join(path, "request.json"), out); err != nil {
				return err
			}
		} else {
			var prior CLIInstallation
			if err := readPrivateValue(filepath.Join(path, "request.json"), privateRunState, 1<<20, &prior); err != nil {
				return err
			}
			if !equal(prior, out) {
				return workflowError(4, "installation request differs")
			}
		}
		if action == "install" {
			var prior CLIInstallSmoke
			if err := readPrivateValue(filepath.Join(path, "smoke.json"), privateRunState, 1<<20, &prior); err == nil && (prior.Error != "" || prior.Exit == nil || *prior.Exit != 0) {
				out.Smoke = &prior
				receipt, actual, observeErr := (clirelease.Installer{}).ObserveQualified(out.Source, destination)
				out.Receipt = &receipt
				out.ObservedSHA256 = actual
				if observeErr != nil {
					out.State = "effect_unknown"
					out.NextAction = "Inspect the receipt and preserved smoke failure; no installation repeated."
					return observeErr
				}
				if receipt.State == "rolled_back" {
					out.State = "rolled_back"
					out.NextAction = "Previous binary restored; the failed smoke evidence is preserved."
					return workflowError(4, "installation was rolled back after smoke failure")
				}
				if actual != out.Source.BinarySHA256 {
					out.State = "effect_unknown"
					out.NextAction = "Inspect changed destination and preserved smoke failure; no installation repeated."
					return workflowError(4, "destination changed after smoke failure")
				}
				out.State = cliSmokeFailureState(prior)
				out.NextAction = "Inspect the preserved failure and use guarded --rollback; do not repeat installation."
				return workflowError(4, "preserved installed CLI smoke failed")
			} else if err != nil && !os.IsNotExist(err) {
				return err
			}
		}
		out.State = "effect_unknown"
		out.NextAction = "Inspect the same operation with --observe before resuming; this invocation has not started replacement."
		if err := d.fault("cli_install_before_effect"); err != nil {
			return err
		}
		installer := clirelease.Installer{}
		var receipt clirelease.Receipt
		var effectErr error
		if action == "rollback" {
			receipt, effectErr = installer.RollbackQualified(out.Source, destination)
		} else if action == "observe" {
			receipt, out.ObservedSHA256, effectErr = installer.ObserveQualified(out.Source, destination)
		} else {
			receipt, effectErr = installer.InstallQualified(out.Source, binary, destination)
		}
		out.Receipt = &receipt
		out.State = "effect_unknown"
		out.NextAction = "Observe the destination and preserved receipt before resuming the same operation."
		if effectErr != nil {
			out.State = "effect_unknown"
			out.NextAction = "Inspect the preserved receipt and destination with --observe. Resume only the same inputs or use guarded --rollback."
			return effectErr
		}
		if err := d.fault("cli_install_after_effect"); err != nil {
			return err
		}
		if receipt.State == "rolled_back" {
			actual, _, readErr := installedCLIBytes(destination)
			out.ObservedSHA256 = actual
			if readErr != nil || actual != receipt.BeforeSHA256 {
				out.State = "effect_unknown"
				return workflowError(4, "destination changed after recorded rollback")
			}
			var prior CLIInstallSmoke
			if err := readPrivateValue(filepath.Join(path, "smoke.json"), privateRunState, 1<<20, &prior); err == nil {
				out.Smoke = &prior
			} else if !os.IsNotExist(err) {
				return err
			}
			out.State = "rolled_back"
			out.NextAction = "Previous installed bytes restored; preserve the failed installation evidence."
			return nil
		}
		if action == "observe" {
			var prior CLIInstallSmoke
			if err := readPrivateValue(filepath.Join(path, "smoke.json"), privateRunState, 1<<20, &prior); err == nil {
				out.Smoke = &prior
			} else if !os.IsNotExist(err) {
				return err
			}
		}
		if action == "observe" && out.ObservedSHA256 != out.Source.BinarySHA256 {
			out.State = "effect_unknown"
			out.NextAction = "Destination does not contain the qualified candidate. Inspect before resuming or rolling back."
			return workflowError(4, "installed effect is not the qualified candidate")
		}
		if receipt.State != "installed" {
			out.State = "pending"
			out.NextAction = "Resume install-cli with identical inputs to reconcile the pending replacement."
			return workflowError(4, "installation receipt is pending")
		}
		var smoke CLIInstallSmoke
		smokePath := filepath.Join(path, "smoke.json")
		if err := readPrivateValue(smokePath, privateRunState, 1<<20, &smoke); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			if action == "observe" {
				out.State = "smoke_not_recorded"
				out.NextAction = "Resume the same installation to run its bounded installed CLI smoke; no replacement is needed."
				return workflowError(4, "installed hash observed but no smoke result is recorded")
			}
			if err := d.fault("cli_install_before_smoke"); err != nil {
				return err
			}
			smoke = probeInstalledCLI(destination, out.Source.BinarySHA256)
			if err = d.writeValue(smokePath, smoke); err != nil {
				return err
			}
			// Log only when recording the first successful installed smoke. Check,
			// observe and identical retries reuse it without repeating this warning.
			var embedded buildinfo.Identity
			if smoke.Error == "" && smoke.Exit != nil && *smoke.Exit == 0 && json.Unmarshal([]byte(smoke.Stdout), &embedded) == nil && embedded.SourceCommit == "" {
				logrus.WithFields(logrus.Fields{
					"run_id": id, "embedded_identity_source": "installed executable version --format json",
					"qualification_source": out.Source.VerificationSHA256,
					"qualification_oid":    out.Source.OID, "qualification_tree": out.Source.Tree,
					"binary_sha256": out.Source.BinarySHA256,
					"handling":      "installation continues with exact qualified bytes; embedded commit remains unknown; use make build for future candidates",
				}).Warn("Qualified CLI embedded source identity is unknown")
			}
		}
		out.Smoke = &smoke
		actual, _, err := installedCLIBytes(destination)
		out.ObservedSHA256 = actual
		if err != nil || actual != out.Source.BinarySHA256 {
			out.State = "effect_unknown"
			return workflowError(4, "installed bytes changed since smoke; no installation repeated")
		}
		if smoke.Error != "" || smoke.Exit == nil || *smoke.Exit != 0 {
			out.State = cliSmokeFailureState(smoke)
			out.NextAction = "Inspect preserved smoke stderr and exit; use --rollback to restore the guarded backup."
			return workflowError(4, "installed CLI smoke failed")
		}
		out.State = "verified"
		out.NextAction = "Installation and installed CLI smoke are observed. No action required."
		return nil
	})
	if err != nil {
		logrus.WithError(err).Warn("qualified CLI installation requires observation")
	}
	return out, err
}
func installedCLIBytes(path string) (string, []byte, error) {
	b, err := clirelease.ReadRegular(path, workflowhandoff.AutomaticCandidateBinaryLimit)
	return hash(b), b, err
}
func probeInstalledCLI(path, expected string) CLIInstallSmoke {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, "version", "--format", "json")
	cmd.Dir = filepath.Dir(path)
	cmd.Env = []string{"PATH=/usr/bin:/bin", "LC_ALL=C"}
	var stdout, stderr workflowLimitedOutput
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	r := CLIInstallSmoke{Started: cmd.Process != nil, Stdout: string(stdout.data), Stderr: string(stderr.data)}
	if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
		r.Exit = ptr(cmd.ProcessState.ExitCode())
	}
	if err != nil {
		r.Error = err.Error()
		return r
	}
	var identity buildinfo.Identity
	if err = json.Unmarshal(stdout.data, &identity); err != nil || stdout.overflow || stderr.overflow || identity.Kind != "PlyExecutable@1" || identity.SHA256 != expected {
		r.Error = fmt.Sprint("installed CLI identity mismatch: ", err)
	}
	return r
}

func cliSmokeFailureState(s CLIInstallSmoke) string {
	if s.Started && s.Exit == nil {
		return "smoke_unknown"
	}
	return "smoke_failed"
}
