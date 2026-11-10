package workspace

import (
	"fmt"
	"os"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

const backupCommittedDocumentPrefix = "backup-committed-dryrun-"

// BackupCommittedWork describes an inert Git patch. Its header carries the
// source identity even when the source worktree is no longer available.
// Bytes are transient; native TaskContent CAS owns their durable preservation.
type BackupCommittedWork struct {
	RunID         string      `json:"run_id"`
	RequestSHA256 string      `json:"request_sha256"`
	TaskID        TaskID      `json:"task_id"`
	Goal          TaskGoalRef `json:"goal"`
	GitCommonDir  string      `json:"git_common_dir"`
	SourceRef     string      `json:"source_ref"`
	BaseOID       string      `json:"base_oid"`
	BaseTree      string      `json:"base_tree"`
	SourceOID     string      `json:"source_oid"`
	SourceTree    string      `json:"source_tree"`
	Locator       string      `json:"locator"`
	SHA256        string      `json:"sha256"`
	SizeBytes     int         `json:"size_bytes"`
	Bytes         []byte      `json:"-"`
}

type BackupCommittedGit interface {
	ObserveBackupCommitted(path, common, ref, base, baseTree, source, sourceTree string) (*BackupCommittedWork, error)
}

func (git *systemTaskIntegrationGit) ObserveBackupCommitted(path, common, ref, base, baseTree, source, sourceTree string) (*BackupCommittedWork, error) {
	fail := func(message string) (*BackupCommittedWork, error) {
		return nil, fmt.Errorf("backup_goal_committed_source_invalid: %s", message)
	}
	if !validOIDText(base) || !validOIDText(source) || !validOIDText(baseTree) || !validOIDText(sourceTree) {
		return fail("invalid commit/tree identity")
	}
	observed, err := git.ObserveIntegrationWorktree(path, ref)
	if err != nil {
		return nil, err
	}
	if observed.GitCommonDir != common || observed.Ref != ref || !observed.Clean || len(observed.InProgress) != 0 || observed.OID != source || observed.Tree != sourceTree {
		return fail("source identity or clean status changed")
	}
	tree, err := integrationLine(git.read(path, "rev-parse", "--verify", base+"^{tree}"), "observe original base tree")
	if err != nil || tree != baseTree {
		return fail("original base/tree unavailable or changed")
	}
	ancestor, err := git.CheckAncestor(RepoRecord{Locator: path, GitCommonDir: common}, base, source)
	if err != nil {
		return nil, err
	}
	if !ancestor {
		return fail("source is not a descendant of the bound original base")
	}
	patch := git.read(path, "-c", "core.quotePath=true", "-c", "diff.suppressBlankEmpty=false", "diff", "--binary", "--full-index", "--no-ext-diff", "--no-textconv", "--no-renames", "--no-color", "--diff-algorithm=myers", "--no-indent-heuristic", "--unified=3", "--inter-hunk-context=0", "--no-relative", "--ignore-submodules=none", "--submodule=short", "-O"+os.DevNull, "--output-indicator-new=+", "--output-indicator-old=-", "--output-indicator-context= ", "--src-prefix=a/", "--dst-prefix=b/", base, source, "--")
	if patch.Err != nil {
		return nil, fmt.Errorf("read committed dry-run patch: %s", safeGitDetail(patch))
	}
	if len(patch.Stdout) == 0 {
		return fail("no committed changes from original base")
	}
	if err = taskDocumentBytesValid(patch.Stdout); err != nil {
		return nil, fmt.Errorf("backup_goal_committed_input_limit: %w", err)
	}
	return &BackupCommittedWork{GitCommonDir: common, SourceRef: ref, BaseOID: base, BaseTree: baseTree, SourceOID: source, SourceTree: sourceTree, Bytes: patch.Stdout}, nil
}

// Bind adds a deterministic, Git-readable provenance header and validates the
// existing document limit before any retirement or publication can happen.
func (w *BackupCommittedWork) Bind(root string) error {
	metadata := *w
	metadata.Locator, metadata.SHA256, metadata.SizeBytes = "", "", 0
	raw, err := canonicaljson.Marshal(contentRefValue(metadata))
	if err != nil {
		return err
	}
	header := "# Preserved committed source; dry-run input only.\n# Assess, optionally apply/adapt, test, commit and deliver under the fresh mandate.\n# No old acceptance, permission, tests, review, QA, result, qualification or delivery is inherited.\n# Source: " + string(raw) + "\n\n"
	w.Bytes = append([]byte(header), w.Bytes...)
	if err = taskDocumentBytesValid(w.Bytes); err != nil {
		return fmt.Errorf("backup_goal_committed_input_limit: %w", err)
	}
	w.SHA256, w.SizeBytes = digestTaskBytes(w.Bytes), len(w.Bytes)
	w.Locator = taskContentPath(root, "objects", w.SHA256)
	return nil
}

// DocumentID scopes each retained draft to its exact bytes, including source identity.
func (w *BackupCommittedWork) DocumentID() string {
	return backupCommittedDocumentPrefix + strings.TrimPrefix(w.SHA256, "sha256:")
}

func (w *BackupCommittedWork) document() canonicaljson.Value {
	return contentRefValue(map[string]any{"id": w.DocumentID(), "source": map[string]any{"kind": "file", "locator": w.Locator, "sha256": w.SHA256, "size_bytes": w.SizeBytes, "media_type": "text/plain", "git_provenance": nil}})
}

func validateBackupCommitted(w *BackupCommittedWork, root string) error {
	if w == nil {
		return nil
	}
	if err := taskDocumentBytesValid(w.Bytes); err != nil {
		return err
	}
	if digestTaskBytes(w.Bytes) != w.SHA256 || len(w.Bytes) != w.SizeBytes || w.Locator != taskContentPath(root, "objects", w.SHA256) || !strings.HasPrefix(string(w.Bytes), "# Preserved committed source;") {
		return WorkInvalidArguments("committed dry-run bytes or descriptor differ")
	}
	return nil
}
