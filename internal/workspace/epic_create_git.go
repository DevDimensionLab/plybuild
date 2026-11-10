package workspace

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

func epicGitArgs(args []string) []string {
	return append([]string{"-c", "core.hooksPath=" + os.DevNull, "-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false", "-c", "submodule.recurse=false", "-c", "protocol.allow=never", "-c", "core.attributesFile=" + os.DevNull}, args...)
}

// Neutralize every locally configured driver before status inspects working
// attributes. Global/system config and inherited Git environment are already
// excluded by the native safe runner. Configuration reads execute no filters.
func epicSafeGitArgs(args []string) ([]string, GitCommandOutcome) {
	path := ""
	for i := 0; i+1 < len(args); i++ {
		if args[i] == "-C" {
			path = args[i+1]
		}
	}
	prefix := []string{}
	if path != "" {
		config := runSafeIntegrationGit(epicGitArgs([]string{"-C", path, "config", "--null", "--name-only", "--get-regexp", `^filter\..*\.(clean|smudge|process|required)$`}), []string{"GIT_ALLOW_PROTOCOL="})
		if config.Err != nil && config.Exit != 1 {
			return nil, config
		}
		drivers := map[string]bool{}
		for _, raw := range bytes.Split(config.Stdout, []byte{0}) {
			key := string(raw)
			if key == "" {
				continue
			}
			dot := strings.LastIndexByte(key, '.')
			if dot <= len("filter.") {
				return nil, GitCommandOutcome{Exit: 1, Err: fmt.Errorf("malformed filter configuration key")}
			}
			drivers[key[:dot]] = true
		}
		for driver := range drivers {
			for _, suffix := range []string{"clean=", "smudge=", "process=", "required=false"} {
				prefix = append(prefix, "-c", driver+"."+suffix)
			}
		}
	}
	return epicGitArgs(append(prefix, args...)), GitCommandOutcome{}
}
func runEpicCreateGit(args, env []string) GitCommandOutcome {
	argv, failed := epicSafeGitArgs(args)
	if failed.Err != nil {
		return failed
	}
	return runSafeIntegrationGit(argv, append(env, "GIT_ALLOW_PROTOCOL="))
}
func (g *systemWorkItemGit) checkEpicParentOperation(path string) error {
	for _, name := range []string{"MERGE_HEAD", "CHERRY_PICK_HEAD", "REVERT_HEAD", "BISECT_LOG", "rebase-merge", "rebase-apply", "sequencer", "index.lock"} {
		p, e := g.gitLine(path, "rev-parse", "--path-format=absolute", "--git-path", name)
		if e != nil {
			return e
		}
		if _, e = g.files.Lstat(p); e == nil {
			return fmt.Errorf("parent has active Git operation %s; complete it before creating an Epic", name)
		} else if !errors.Is(e, os.ErrNotExist) {
			return e
		}
	}
	return nil
}
func (g *systemWorkItemGit) checkEpicMaterialization(repo RepoRecord, oid string) error {
	objects := g.read(repo.Locator, "rev-list", "--objects", "--missing=print", oid)
	if objects.Err != nil {
		return gitObservationError("inspect complete parent objects", repo.Locator, objects)
	}
	for _, line := range bytes.Split(objects.Stdout, []byte{'\n'}) {
		if bytes.HasPrefix(line, []byte{'?'}) {
			return fmt.Errorf("parent objects are incomplete; materialize them independently before Epic creation")
		}
	}
	tree := g.read(repo.Locator, "ls-tree", "-r", "-z", oid)
	if tree.Err != nil {
		return gitObservationError("inspect parent tree", repo.Locator, tree)
	}
	paths := []byte{}
	for _, entry := range bytes.Split(tree.Stdout, []byte{0}) {
		if len(entry) == 0 {
			continue
		}
		header, path, ok := bytes.Cut(entry, []byte{'\t'})
		if !ok {
			return fmt.Errorf("malformed parent tree")
		}
		if !bytes.HasPrefix(header, []byte("100644 blob ")) && !bytes.HasPrefix(header, []byte("100755 blob ")) {
			return fmt.Errorf("Epic create v1 requires regular files; symlinks and submodules need independent materialization")
		}
		paths = append(paths, path...)
		paths = append(paths, 0)
	}
	if len(paths) > 0 {
		args, failed := epicSafeGitArgs([]string{"-C", repo.Locator, "check-attr", "--source=" + oid, "-z", "--stdin", "filter"})
		if failed.Err != nil {
			return failed.Err
		}
		attr := runSafeIntegrationGitInput(args, []string{"GIT_ALLOW_PROTOCOL=", "GIT_OPTIONAL_LOCKS=0"}, paths)
		if attr.Err != nil {
			return gitObservationError("inspect whole-tree attributes", repo.Locator, attr)
		}
		if e := validateFilterAttributeFrames(paths, attr.Stdout); e != nil {
			return e
		}
	}
	// Sparse checkout and worktree-specific config can change the chosen tree.
	for _, key := range []string{"core.sparseCheckout", "extensions.worktreeConfig"} {
		c := g.read(repo.Locator, "config", "--bool", "--get", key)
		if c.Exit != 1 && (c.Err != nil || strings.TrimSpace(string(c.Stdout)) != "false") {
			return fmt.Errorf("Epic create requires %s disabled; inspect repository configuration", key)
		}
	}
	return nil
}
