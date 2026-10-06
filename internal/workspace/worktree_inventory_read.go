package workspace

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// InventoryFactFreshness separates a Git inventory entry from observations of
// its checkout. A missing/prunable directory can still have a known entry.
type InventoryFactFreshness struct {
	Inventory     string `json:"inventory"`
	Identity      string `json:"identity"`
	Clean         string `json:"clean"`
	BaseRelation  string `json:"base_relation"`
	LastCommitUTC string `json:"last_commit_utc"`
}
type InventoryWorktreeFacts struct {
	Entry          GitWorktreeInventoryEntry
	Clean          *bool
	Ahead, Behind  *int
	MergedIntoBase *bool
	LastCommitUTC  *string
	Freshness      InventoryFactFreshness
	Reasons        []string
}
type InventoryRepositoryFacts struct {
	Worktrees []InventoryWorktreeFacts
	Freshness string
	Reasons   []string
}
type WorktreeInventoryGit interface {
	ReadWorktreeInventory(RepoRecord, map[string]string) (InventoryRepositoryFacts, error)
}

// ReadRegisteredWorktreeInventory is intentionally separate from readiness and
// integration observers. It tolerates incomplete checkouts, never fetches, and
// compares only the explicitly supplied local base OIDs.
func ReadRegisteredWorktreeInventory(d Dependencies, repo RepoRecord, bases map[string]string) (InventoryRepositoryFacts, error) {
	reader, ok := d.IntegrationGit.(WorktreeInventoryGit)
	if !ok {
		return InventoryRepositoryFacts{}, workError(ErrorWorkGitObservation, "read-only inventory observer is unavailable", nil)
	}
	return reader.ReadWorktreeInventory(repo, bases)
}

func (g *systemTaskIntegrationGit) inventoryRead(path string, args ...string) GitCommandOutcome {
	prefix := []string{"-c", "core.hooksPath=" + os.DevNull, "-c", "gc.auto=0", "-c", "maintenance.auto=false"}
	prefix = append(prefix, integrationGitPrefix(path)...)
	return g.run(append(prefix, args...), []string{"GIT_OPTIONAL_LOCKS=0"})
}
func (g *systemTaskIntegrationGit) ReadWorktreeInventory(repo RepoRecord, bases map[string]string) (InventoryRepositoryFacts, error) {
	out := InventoryRepositoryFacts{Worktrees: []InventoryWorktreeFacts{}, Freshness: "fresh", Reasons: []string{}}
	common, err := integrationLine(g.inventoryRead(repo.Locator, "rev-parse", "--path-format=absolute", "--git-common-dir"), "observe inventory repository")
	if err != nil {
		return out, err
	}
	actual, err := canonicalDirectory(g.files, "", common, "inventory Git directory")
	if err != nil || actual != repo.GitCommonDir {
		return out, workError(ErrorWorkIdentityConflict, "registered repository Git identity differs", err)
	}
	before := g.inventoryRead(repo.Locator, "worktree", "list", "--porcelain", "-z")
	if before.Err != nil {
		return out, gitObservationError("read worktree inventory", repo.Locator, before)
	}
	entries, err := parseReadInventory(before.Stdout)
	if err != nil {
		return out, err
	}
	for _, entry := range entries {
		out.Worktrees = append(out.Worktrees, g.inventoryWorktree(repo, entry, bases[entry.Locator]))
	}
	after := g.inventoryRead(repo.Locator, "worktree", "list", "--porcelain", "-z")
	if after.Err != nil || !bytes.Equal(before.Stdout, after.Stdout) {
		out.Freshness = "stale"
		out.Reasons = append(out.Reasons, "worktree_inventory_changed_during_read")
		for i := range out.Worktrees {
			out.Worktrees[i].Freshness.Inventory = "stale"
			out.Worktrees[i].Freshness.Identity = "stale"
			out.Worktrees[i].Freshness.Clean = "stale"
			out.Worktrees[i].Freshness.BaseRelation = "stale"
			out.Worktrees[i].Reasons = append(out.Worktrees[i].Reasons, "worktree_inventory_changed_during_read")
		}
	}
	return out, nil
}

func (g *systemTaskIntegrationGit) inventoryWorktree(repo RepoRecord, entry GitWorktreeInventoryEntry, base string) InventoryWorktreeFacts {
	out := InventoryWorktreeFacts{Entry: entry, Freshness: InventoryFactFreshness{"fresh", "unknown", "unknown", "unknown", "unknown"}, Reasons: []string{}}
	physical, err := canonicalDirectory(g.files, "", entry.Locator, "inventory worktree")
	if err != nil || physical != entry.Locator || entry.Bare {
		out.Reasons = append(out.Reasons, "worktree_checkout_unavailable")
	} else {
		identity := g.inventoryRead(entry.Locator, "rev-parse", "--path-format=absolute", "--git-common-dir", "--show-toplevel", "HEAD")
		parts := strings.Split(strings.TrimSuffix(string(identity.Stdout), "\n"), "\n")
		if identity.Err != nil || len(parts) != 3 || parts[0] != repo.GitCommonDir || parts[1] != entry.Locator || parts[2] != entry.OID {
			out.Reasons = append(out.Reasons, "worktree_identity_differs")
		} else {
			out.Freshness.Identity = "fresh"
			status := g.inventoryRead(entry.Locator, "status", "--porcelain=v1", "-z", "--untracked-files=normal")
			if status.Err != nil {
				out.Reasons = append(out.Reasons, "worktree_status_unavailable")
			} else {
				clean := len(status.Stdout) == 0
				out.Clean = &clean
				out.Freshness.Clean = "fresh"
			}
		}
	}
	if validOIDText(entry.OID) {
		commit := g.inventoryRead(repo.Locator, "show", "-s", "--format=%cI", entry.OID, "--")
		if commit.Err == nil {
			if t, e := time.Parse(time.RFC3339, strings.TrimSpace(string(commit.Stdout))); e == nil {
				value := t.UTC().Format(time.RFC3339)
				out.LastCommitUTC = &value
				out.Freshness.LastCommitUTC = "fresh"
			}
		}
		if out.LastCommitUTC == nil {
			out.Reasons = append(out.Reasons, "commit_time_unavailable")
		}
	}
	if base != "" && validOIDText(base) && validOIDText(entry.OID) {
		relation := g.inventoryRead(repo.Locator, "rev-list", "--left-right", "--count", "--end-of-options", base+"..."+entry.OID, "--")
		parts := strings.Fields(string(relation.Stdout))
		if relation.Err == nil && len(parts) == 2 {
			behind, e1 := strconv.Atoi(parts[0])
			ahead, e2 := strconv.Atoi(parts[1])
			if e1 == nil && e2 == nil && ahead >= 0 && behind >= 0 {
				merged := ahead == 0
				out.Ahead = &ahead
				out.Behind = &behind
				out.MergedIntoBase = &merged
				out.Freshness.BaseRelation = "fresh"
			}
		}
		if out.Ahead == nil {
			out.Reasons = append(out.Reasons, "base_relation_unavailable")
		}
	}
	return out
}

// Unlike the effect-oriented parser, this preserves a missing/prunable path;
// it makes no claim that the reported checkout exists or is a physical root.
func parseReadInventory(b []byte) ([]GitWorktreeInventoryEntry, error) {
	out := []GitWorktreeInventoryEntry{}
	var current *GitWorktreeInventoryEntry
	seen := map[string]bool{}
	finish := func() error {
		if current == nil {
			return nil
		}
		if !filepath.IsAbs(current.Locator) || filepath.Clean(current.Locator) != current.Locator || (!current.Bare && !validOIDText(current.OID)) || (!current.Bare && !current.Detached && !validFullBranchRef(current.Ref)) {
			return fmt.Errorf("invalid Git inventory record")
		}
		out = append(out, *current)
		current = nil
		return nil
	}
	for _, field := range bytes.Split(b, []byte{0}) {
		text := string(field)
		if text == "" {
			if err := finish(); err != nil {
				return nil, err
			}
			continue
		}
		key, value, _ := strings.Cut(text, " ")
		if key == "worktree" {
			if err := finish(); err != nil {
				return nil, err
			}
			current = &GitWorktreeInventoryEntry{Locator: value}
			seen = map[string]bool{"worktree": true}
			continue
		}
		if current == nil || seen[key] {
			return nil, fmt.Errorf("invalid Git inventory field %q", key)
		}
		seen[key] = true
		switch key {
		case "HEAD":
			current.OID = value
		case "branch":
			current.Ref = value
		case "locked":
			current.Locked = true
		case "prunable":
			current.Prunable = true
		case "bare":
			current.Bare = true
		case "detached":
			current.Detached = true
		default:
			return nil, fmt.Errorf("unknown Git inventory field %q", key)
		}
	}
	if err := finish(); err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Locator < out[j].Locator })
	return out, nil
}
