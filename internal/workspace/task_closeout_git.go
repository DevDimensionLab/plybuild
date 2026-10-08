package workspace

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// This adapter deliberately has no fetch, prune, unlock, force or remote effect.
func closeoutGit(p TaskCloseoutPlan, path string, write bool, args ...string) GitCommandOutcome {
	prefix := []string{"-c", "core.hooksPath=" + os.DevNull, "-c", "gc.auto=0", "-c", "maintenance.auto=false"}
	prefix = append(prefix, integrationGitPrefix(path)...)
	env := []string{}
	if !write {
		env = append(env, "GIT_OPTIONAL_LOCKS=0")
	}
	return runSafeIntegrationGit(append(prefix, args...), env)
}

func closeoutInventory(p TaskCloseoutPlan) ([]GitWorktreeInventoryEntry, error) {
	if err := contentPhysical(p.RepositoryLocator); err != nil {
		return nil, err
	}
	if err := contentPhysical(p.Source.GitCommonDir); err != nil {
		return nil, err
	}
	common, err := integrationLine(closeoutGit(p, p.RepositoryLocator, false, "rev-parse", "--path-format=absolute", "--git-common-dir"), "check closeout repository")
	if err != nil || common != p.Source.GitCommonDir {
		return nil, closeoutError("cleanup_blocked", "Repository common directory changed.", err)
	}
	o := closeoutGit(p, p.RepositoryLocator, false, "worktree", "list", "--porcelain", "-z")
	if o.Err != nil {
		return nil, fmt.Errorf("read closeout inventory: %s", safeGitDetail(o))
	}
	return parseReadInventory(o.Stdout)
}

func closeoutObserveSource(p TaskCloseoutPlan, allowRemoved, allowBranchRemoved bool) error {
	entries, err := closeoutInventory(p)
	if err != nil {
		return err
	}
	var entry *GitWorktreeInventoryEntry
	for i := range entries {
		e := &entries[i]
		if e.Locator == p.Source.Locator {
			entry = e
		}
		if e.Ref == p.Source.Ref && e.Locator != p.Source.Locator {
			return fmt.Errorf("source branch is checked out elsewhere")
		}
	}
	_, statErr := os.Lstat(p.Source.Locator)
	if os.IsNotExist(statErr) {
		if !allowRemoved || entry != nil {
			return fmt.Errorf("source_missing: absence has no completed worktree removal observation")
		}
	} else {
		if statErr != nil {
			return statErr
		}
		if err := contentPhysical(p.Source.Locator); err != nil {
			return err
		}
		if entry == nil || entry.Ref != p.Source.Ref || entry.OID != p.ResultOID || entry.Locked || entry.Prunable || entry.Bare || entry.Detached {
			return fmt.Errorf("source inventory identity changed, locked or unavailable")
		}
		// Linked worktrees have a regular .git file and a private administrative
		// directory. The repository's main checkout must never qualify.
		info, err := os.Lstat(filepath.Join(p.Source.Locator, ".git"))
		if err != nil || !info.Mode().IsRegular() {
			return fmt.Errorf("source is not a linked worktree")
		}
		identity := closeoutGit(p, p.Source.Locator, false, "rev-parse", "--path-format=absolute", "--git-common-dir", "--git-dir", "--show-toplevel", "HEAD", "HEAD^{tree}")
		fields := strings.Split(strings.TrimSuffix(string(identity.Stdout), "\n"), "\n")
		if identity.Err != nil || len(fields) != 5 || fields[0] != p.Source.GitCommonDir || fields[1] == fields[0] || !pathWithin(filepath.Join(fields[0], "worktrees"), fields[1]) || fields[2] != p.Source.Locator || fields[3] != p.ResultOID || fields[4] != p.ResultTree {
			return fmt.Errorf("source path, common directory or candidate changed")
		}
		if err := contentPhysical(fields[1]); err != nil {
			return err
		}
		branch, err := integrationLine(closeoutGit(p, p.Source.Locator, false, "symbolic-ref", "--quiet", "HEAD"), "check source branch")
		if err != nil || branch != p.Source.Ref {
			return fmt.Errorf("source branch changed")
		}
		status := closeoutGit(p, p.Source.Locator, false, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--ignored=matching", "--ignore-submodules=none")
		if status.Err != nil || len(status.Stdout) != 0 {
			return fmt.Errorf("source contains tracked, untracked, ignored or submodule changes")
		}
		index := closeoutGit(p, p.Source.Locator, false, "ls-files", "--stage", "-z")
		if index.Err != nil {
			return fmt.Errorf("source index could not be inspected")
		}
		for _, line := range bytes.Split(index.Stdout, []byte{0}) {
			if bytes.HasPrefix(line, []byte("160000 ")) {
				return fmt.Errorf("submodule cleanup requires separate resolution")
			}
		}
		// Status and even no-force worktree removal trust index flags that can
		// hide changed tracked bytes. Require ordinary cached entries; lowercase
		// tags mark assume-unchanged and S marks skip-worktree (including sparse).
		flags := closeoutGit(p, p.Source.Locator, false, "ls-files", "-v", "-z")
		if flags.Err != nil {
			return fmt.Errorf("source index flags could not be inspected")
		}
		for _, entry := range bytes.Split(flags.Stdout, []byte{0}) {
			if len(entry) > 0 && (len(entry) < 3 || entry[0] != 'H' || entry[1] != ' ') {
				return fmt.Errorf("source index flags may hide tracked changes; clear assume-unchanged or skip-worktree and inspect the files before cleanup")
			}
		}
		for _, name := range []string{"MERGE_HEAD", "CHERRY_PICK_HEAD", "REVERT_HEAD", "rebase-merge", "rebase-apply", "BISECT_LOG", "index.lock", "locked"} {
			if _, err := os.Lstat(filepath.Join(fields[1], name)); err == nil || !os.IsNotExist(err) {
				return fmt.Errorf("source has an in-progress Git operation or lock: %s", name)
			}
		}
		if err := filepath.WalkDir(p.Source.Locator, func(path string, de fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if path == p.Source.Locator {
				return nil
			}
			if de.Name() == ".git" && filepath.Dir(path) != p.Source.Locator {
				return fmt.Errorf("source contains a nested repository")
			}
			if de.IsDir() {
				if info, err := os.Lstat(filepath.Join(path, "HEAD")); err == nil && info.Mode().IsRegular() {
					if info, err := os.Lstat(filepath.Join(path, "objects")); err == nil && info.IsDir() {
						return fmt.Errorf("source contains a nested bare repository")
					}
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}
	ref := closeoutGit(p, p.RepositoryLocator, false, "rev-parse", "--verify", "--quiet", p.Source.Ref)
	if ref.Err != nil {
		if allowBranchRemoved && ref.Exit == 1 {
			return nil
		}
		return fmt.Errorf("source branch missing or unreadable")
	}
	if strings.TrimSpace(string(ref.Stdout)) != p.ResultOID {
		return fmt.Errorf("source ref/OID changed")
	}
	if symbolic := closeoutGit(p, p.RepositoryLocator, false, "symbolic-ref", "--quiet", p.Source.Ref); symbolic.Exit != 1 {
		return fmt.Errorf("source branch must remain a direct ref, not a symbolic ref")
	}
	return nil
}

func closeoutCheckAncestor(p TaskCloseoutPlan) error {
	out := closeoutGit(p, p.RepositoryLocator, false, "merge-base", "--is-ancestor", p.ResultOID, p.TargetRef)
	if out.Err != nil {
		return fmt.Errorf("observed local integration is no longer present at the target ref")
	}
	return nil
}

func closeoutArchiveCandidate(p TaskCloseoutPlan) error {
	if err := closeoutCheckCandidate(p); err != nil {
		return err
	}
	current := closeoutGit(p, p.RepositoryLocator, false, "rev-parse", "--verify", "--quiet", p.ArchiveRef)
	if current.Err == nil {
		if strings.TrimSpace(string(current.Stdout)) != p.ResultOID {
			return fmt.Errorf("candidate archive ref conflicts")
		}
		if symbolic := closeoutGit(p, p.RepositoryLocator, false, "symbolic-ref", "--quiet", p.ArchiveRef); symbolic.Exit != 1 {
			return fmt.Errorf("candidate archive must be a direct ref")
		}
		return nil
	}
	if current.Exit != 1 {
		return fmt.Errorf("cannot observe candidate archive ref")
	}
	out := closeoutGit(p, p.RepositoryLocator, true, "update-ref", "--no-deref", p.ArchiveRef, p.ResultOID, strings.Repeat("0", len(p.ResultOID)))
	if out.Err != nil {
		return fmt.Errorf("archive candidate: %s", safeGitDetail(out))
	}
	return nil
}

func closeoutCheckCandidate(p TaskCloseoutPlan) error {
	oid, err := integrationLine(closeoutGit(p, p.RepositoryLocator, false, "rev-parse", "--verify", p.ResultOID+"^{commit}"), "check retained candidate")
	if err != nil || oid != p.ResultOID {
		return fmt.Errorf("candidate commit unavailable")
	}
	tree, err := integrationLine(closeoutGit(p, p.RepositoryLocator, false, "rev-parse", "--verify", p.ResultOID+"^{tree}"), "check retained candidate tree")
	if err != nil || tree != p.ResultTree {
		return fmt.Errorf("candidate tree differs")
	}
	return nil
}

func closeoutRemoveWorktree(p TaskCloseoutPlan) error {
	// A durable reservation is required by the caller before invoking this.
	if err := closeoutObserveSource(p, true, false); err != nil {
		return err
	}
	if _, err := os.Lstat(p.Source.Locator); os.IsNotExist(err) {
		return nil
	}
	out := closeoutGit(p, p.RepositoryLocator, true, "worktree", "remove", "--", p.Source.Locator)
	if out.Err != nil {
		return closeoutError("cleanup_blocked", "Git did not remove the source worktree.", fmt.Errorf("%s", safeGitDetail(out)))
	}
	if err := closeoutObserveSource(p, true, false); err != nil {
		return err
	}
	if _, err := os.Lstat(p.Source.Locator); !os.IsNotExist(err) {
		return fmt.Errorf("worktree removal has not been observed")
	}
	return nil
}

func closeoutRemoveBranch(p TaskCloseoutPlan) error {
	if err := closeoutObserveSource(p, true, true); err != nil {
		return err
	}
	if _, err := os.Lstat(p.Source.Locator); !os.IsNotExist(err) {
		return fmt.Errorf("source worktree remains; branch removal is blocked")
	}
	ref := closeoutGit(p, p.RepositoryLocator, false, "rev-parse", "--verify", "--quiet", p.Source.Ref)
	if ref.Err != nil && ref.Exit == 1 {
		return nil
	}
	// update-ref's expected old OID provides a compare-and-swap deletion, including
	// squash candidates that git branch -d cannot establish by ancestry.
	out := closeoutGit(p, p.RepositoryLocator, true, "update-ref", "--no-deref", "-d", p.Source.Ref, p.ResultOID)
	if out.Err != nil {
		return closeoutError("cleanup_blocked", "The exact local branch could not be removed.", fmt.Errorf("%s", safeGitDetail(out)))
	}
	ref = closeoutGit(p, p.RepositoryLocator, false, "rev-parse", "--verify", "--quiet", p.Source.Ref)
	if ref.Err == nil || ref.Exit != 1 {
		return fmt.Errorf("branch removal has not been observed")
	}
	return nil
}
