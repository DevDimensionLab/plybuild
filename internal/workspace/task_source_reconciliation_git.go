package workspace

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// This optional adapter keeps source Git effects and observations injectable.
// Missing adapters fail closed; they never fall back to an OS/provider probe.
type TaskSourceReconciliationGit interface {
	ReconcileSource(WorkspaceTaskPreparePlan, TaskSourceReconciliation, string) (TaskSourceReconciliation, error)
}

func (git *systemTaskIntegrationGit) sourceRun(path string, op TaskSourceReconciliation, args ...string) GitCommandOutcome {
	prefix := []string{"-c", "core.hooksPath=" + os.DevNull, "-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false", "-c", "submodule.recurse=false", "-c", "merge.autoStash=false", "-c", "rerere.enabled=false", "-c", "commit.gpgSign=false", "-c", "merge.gpgSign=false", "-c", "gc.auto=0", "-c", "user.name=Ply source reconciliation", "-c", "user.email=ply-source@localhost", "-C", path}
	return git.run(append(prefix, args...), []string{"GIT_REFLOG_ACTION=ply source " + op.ID, "GIT_EDITOR=false", "GIT_MERGE_AUTOEDIT=no"})
}
func (git *systemTaskIntegrationGit) sourceConflictFingerprint(path string, op TaskSourceReconciliation) (string, error) {
	parts := [][]byte{}
	for _, args := range [][]string{{"diff", "--binary", "--no-ext-diff", "--no-textconv"}, {"diff", "--cached", "--binary", "--no-ext-diff", "--no-textconv"}, {"status", "--porcelain=v1", "-z", "--untracked-files=all"}} {
		out := git.sourceRun(path, op, args...)
		if out.Err != nil {
			return "", out.Err
		}
		parts = append(parts, out.Stdout)
	}
	// Include untracked content: an unchanged path list does not prove unchanged edits.
	files := git.sourceRun(path, op, "ls-files", "--others", "--exclude-standard", "-z")
	if files.Err != nil {
		return "", files.Err
	}
	for _, name := range bytes.Split(files.Stdout, []byte{0}) {
		if len(name) == 0 {
			continue
		}
		full := filepath.Join(path, string(name))
		st, e := os.Lstat(full)
		if e != nil || !st.Mode().IsRegular() {
			return "", fmt.Errorf("cannot fingerprint unexpected source file")
		}
		raw, e := os.ReadFile(full)
		if e != nil {
			return "", e
		}
		parts = append(parts, []byte(digestTaskBytes(raw)))
	}
	return queueDigest(parts), nil
}
func (git *systemTaskIntegrationGit) sourceMergeHead(path string, op TaskSourceReconciliation) string {
	out := git.sourceRun(path, op, "rev-parse", "--verify", "MERGE_HEAD")
	if out.Err != nil {
		return ""
	}
	return strings.TrimSpace(string(out.Stdout))
}
func (git *systemTaskIntegrationGit) sourceConflicts(path string, op TaskSourceReconciliation) ([]string, error) {
	out := git.sourceRun(path, op, "diff", "--name-only", "--diff-filter=U", "-z")
	if out.Err != nil {
		return nil, out.Err
	}
	names := []string{}
	for _, p := range bytes.Split(out.Stdout, []byte{0}) {
		if len(p) > 0 {
			names = append(names, string(p))
		}
	}
	sort.Strings(names)
	return names, nil
}

func (git *systemTaskIntegrationGit) ReconcileSource(plan WorkspaceTaskPreparePlan, op TaskSourceReconciliation, action string) (TaskSourceReconciliation, error) {
	if action != "" && action != "continue" && action != "abort" {
		return op, WorkInvalidArguments("reconciliation action must be continue or abort")
	}
	path, ref := plan.WorktreePath, "refs/heads/"+plan.Branch
	x, err := git.ObserveIntegrationWorktree(path, ref)
	if err != nil {
		return op, err
	}
	if x.Ref != ref || x.GitCommonDir != plan.Target.GitCommonDir {
		return op, queueError("source_reconciliation_identity", "source identity changed")
	}
	marker, err := integrationLine(git.sourceRun(path, op, "rev-parse", "--path-format=absolute", "--git-path", "ply-source-"+op.ID), "source operation marker")
	if err != nil {
		return op, err
	}
	markerBytes := []byte(op.PlanSHA256 + "\n")
	existing, e := os.ReadFile(marker)
	if e != nil && !os.IsNotExist(e) {
		return op, e
	}
	if e == nil && !bytes.Equal(existing, markerBytes) {
		return op, queueError("reconciliation_effect_unknown", "source operation marker differs")
	}
	mergeHead := git.sourceMergeHead(path, op)
	// A changed ref is adopted only with this exact marker, merge parents and
	// native commit message. An unrelated/ambiguous commit remains unknown.
	if x.OID != op.Plan.BeforeOID {
		parents := git.sourceRun(path, op, "show", "-s", "--format=%P", x.OID)
		message := git.sourceRun(path, op, "show", "-s", "--format=%B", x.OID)
		expected := op.Plan.BeforeOID + " " + op.Plan.Parent.OID
		if e == nil && x.Clean && mergeHead == "" && parents.Err == nil && message.Err == nil && strings.TrimSpace(string(parents.Stdout)) == expected && strings.TrimSpace(string(message.Stdout)) == "Ply source reconciliation "+op.ID {
			op.State = "committed"
			op.AfterOID = x.OID
			op.AfterTree = x.Tree
			op.ConflictPaths = []string{}
			return op, nil
		}
		op.State = "effect_unknown"
		return op, nil
	}
	if mergeHead != "" {
		if e != nil || mergeHead != op.Plan.Parent.OID {
			op.State = "effect_unknown"
			return op, nil
		}
		conflicts, e2 := git.sourceConflicts(path, op)
		if e2 != nil {
			return op, e2
		}
		if action == "abort" {
			fingerprint, e2 := git.sourceConflictFingerprint(path, op)
			if e2 != nil {
				return op, e2
			}
			if op.State != "source_conflict" || fingerprint != op.ConflictSHA256 {
				return op, queueError("source_abort_edits_present", "abort refuses changed resolutions or unexpected edits; preserve them before choosing a supported action")
			}
			result := git.sourceRun(path, op, "merge", "--abort")
			after, e2 := git.ObserveIntegrationWorktree(path, ref)
			if e2 == nil && result.Exit == 0 && after.Clean && after.OID == op.Plan.BeforeOID && after.Tree == op.Plan.BeforeTree && git.sourceMergeHead(path, op) == "" {
				op.State = "aborted"
				op.AfterOID = after.OID
				op.AfterTree = after.Tree
				return op, nil
			}
			op.State = "effect_unknown"
			return op, nil
		}
		if action == "continue" {
			if op.State != "source_conflict" || len(conflicts) != 0 {
				return op, queueError("source_conflict", "resolve and stage every conflict before native continue")
			}
			unstaged := git.sourceRun(path, op, "diff", "--quiet", "--no-ext-diff", "--no-textconv")
			if unstaged.Exit != 0 {
				return op, queueError("source_resolution_unstaged", "stage the complete resolution before native continue")
			}
			result := git.sourceRun(path, op, "commit", "--no-gpg-sign", "-m", "Ply source reconciliation "+op.ID)
			if result.Exit != 0 {
				return op, fmt.Errorf("native source resolution commit failed: %s", strings.TrimSpace(string(result.Stderr)))
			}
			return git.ReconcileSource(plan, op, "")
		}
		op.State = "source_conflict"
		op.ConflictPaths = conflicts
		op.ConflictSHA256, err = git.sourceConflictFingerprint(path, op)
		return op, err
	}
	if action != "" {
		return op, queueError("source_reconciliation_not_pending", "continue/abort requires this known native pending merge")
	}
	if !x.Clean || len(x.InProgress) > 0 || x.Tree != op.Plan.BeforeTree {
		op.State = "effect_unknown"
		return op, nil
	}
	if op.State == "source_conflict" {
		op.State = "effect_unknown"
		return op, nil
	}
	drivers := git.sourceRun(path, op, "config", "--local", "--name-only", "--get-regexp", `^merge\..*\.driver$`)
	if drivers.Exit != 1 {
		return op, queueError("source_merge_driver_unsupported", "custom or unreadable merge drivers require a separately supported adapter")
	}
	if e != nil {
		file, e2 := os.OpenFile(marker, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if e2 != nil {
			return op, e2
		}
		_, e2 = file.Write(markerBytes)
		if e2 == nil {
			e2 = file.Sync()
		}
		closeErr := file.Close()
		if e2 != nil {
			return op, e2
		}
		if closeErr != nil {
			return op, closeErr
		}
	}
	contained, e2 := git.CheckAncestor(RepoRecord{Locator: path}, op.Plan.Parent.OID, op.Plan.BeforeOID)
	if e2 != nil {
		return op, e2
	}
	if contained {
		op.State = "committed"
		op.AfterOID = x.OID
		op.AfterTree = x.Tree
		return op, nil
	}
	git.sourceRun(path, op, "merge", "--no-ff", "--no-edit", "--no-stat", "--no-autostash", "--no-gpg-sign", "-m", "Ply source reconciliation "+op.ID, op.Plan.Parent.OID)
	// Observe exact effects even when the command lost its reply or returned nonzero.
	after, e2 := git.ObserveIntegrationWorktree(path, ref)
	if e2 != nil {
		return op, e2
	}
	if after.OID == op.Plan.BeforeOID && git.sourceMergeHead(path, op) == "" {
		op.State = "effect_unknown"
		return op, nil
	}
	return git.ReconcileSource(plan, op, "")
}
