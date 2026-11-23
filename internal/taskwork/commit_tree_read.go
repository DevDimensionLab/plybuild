package taskwork

import (
	"bytes"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Reasons a commit tree could not be read.
const (
	CommitTreeNoCommit      = "no_commit"
	CommitTreeCommitUnknown = "commit_unknown"
	CommitTreeGitFailed     = "git_failed"
)

// CommitTreeEntry is one blob, symlink or gitlink of a commit, with its Git
// mode and slash-separated path.
type CommitTreeEntry struct {
	Mode, Path string
}

// CommitTree is the full tree of one commit, or the Reason it was not read.
type CommitTree struct {
	OID     string
	Entries []CommitTreeEntry
	Reason  string
}

type CommitTreeGit interface {
	ReadCommitTree(repo workspace.RepoRecord, commit string) CommitTree
}

// ReadRegisteredCommitTree reads the tree of commit (HEAD when empty) in the
// repo's registered checkout. It reads objects only, never the work tree.
func ReadRegisteredCommitTree(d Dependencies, repo workspace.RepoRecord, commit string) CommitTree {
	reader, ok := d.IntegrationGit.(CommitTreeGit)
	if !ok {
		return CommitTree{Reason: CommitTreeGitFailed}
	}
	return reader.ReadCommitTree(repo, commit)
}

func (g *systemTaskIntegrationGit) ReadCommitTree(repo workspace.RepoRecord, commit string) CommitTree {
	missing := CommitTreeCommitUnknown
	if commit == "" {
		commit, missing = "HEAD", CommitTreeNoCommit
	}
	resolved := g.inventoryRead(repo.Locator, "rev-parse", "--verify", "--quiet", "--end-of-options", commit+"^{commit}")
	if resolved.Err != nil {
		// --verify --quiet exits 1 without output for a name that is no commit.
		if resolved.Exit == 1 && len(bytes.TrimSpace(resolved.Stderr)) == 0 {
			return CommitTree{Reason: missing}
		}
		return CommitTree{Reason: CommitTreeGitFailed}
	}
	oid := strings.TrimSpace(string(resolved.Stdout))
	listed := g.inventoryRead(repo.Locator, "ls-tree", "-r", "-z", "--full-tree", oid)
	if listed.Err != nil {
		return CommitTree{Reason: CommitTreeGitFailed}
	}
	out := CommitTree{OID: oid, Entries: []CommitTreeEntry{}}
	for _, record := range bytes.Split(listed.Stdout, []byte{0}) {
		meta, path, ok := bytes.Cut(record, []byte{'\t'})
		if !ok {
			continue
		}
		mode, _, _ := bytes.Cut(meta, []byte{' '})
		out.Entries = append(out.Entries, CommitTreeEntry{Mode: string(mode), Path: string(path)})
	}
	return out
}
