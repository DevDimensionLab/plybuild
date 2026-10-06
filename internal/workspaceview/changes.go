package workspaceview

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const ChangesKind = "WorkspaceStatusReadback@1"

// RegisterChange is an invalidation token, not validation of the underlying
// documents or evidence that Git and running processes have been observed.
type RegisterChange struct {
	ID             string   `json:"id"`
	RegistrySHA256 *string  `json:"registry_sha256"`
	Freshness      string   `json:"freshness"`
	FileCount      int      `json:"file_count"`
	Reasons        []string `json:"reasons"`
}

type ChangesResult struct {
	Kind          string `json:"kind"`
	SchemaVersion int    `json:"schema_version"`
	Workspace     struct {
		Root string `json:"root"`
	} `json:"workspace"`
	Scope          string           `json:"scope"`
	RegistrySHA256 *string          `json:"registry_sha256"`
	Freshness      string           `json:"freshness"`
	Registers      []RegisterChange `json:"registers"`
}

type changeOwner struct {
	id           string
	files, trees []string
}

// ReadChanges fingerprints only owners of registered read-model facts. It does
// not scan immutable content objects, terminal transcripts, locks or Git state.
// Two complete passes detect movement without acquiring a reader lock. They do
// not promise a transactionally atomic snapshot across the independent stores.
func ReadChanges(d workspace.Dependencies) (ChangesResult, error) {
	return readChanges(d, nil)
}

func readChanges(d workspace.Dependencies, betweenPasses func()) (ChangesResult, error) {
	w, err := workspace.ObserveContaining(d)
	if err != nil {
		return ChangesResult{}, err
	}
	owners := []changeOwner{
		{"projects", []string{"projects.yaml", "project-metadata.json"}, nil},
		{"work_items", []string{"work-items.yaml", "work-item-lifecycle.json"}, nil},
		{"journal", nil, []string{"task-process/v1"}},
		{"runs", nil, []string{"task-runs/v1/runs", "workflow-runs/requests", "workflow-runs/runs"}},
	}
	first := make([]RegisterChange, len(owners))
	for i, owner := range owners {
		first[i] = fingerprintOwner(w.Root, owner)
	}
	if betweenPasses != nil {
		betweenPasses()
	}
	result := ChangesResult{Kind: ChangesKind, SchemaVersion: 1, Scope: "registered_ply_facts", Freshness: "fresh", Registers: []RegisterChange{}}
	result.Workspace.Root = w.Root
	for i, owner := range owners {
		current := fingerprintOwner(w.Root, owner)
		previous := first[i]
		if previous.Freshness == "unknown" || current.Freshness == "unknown" {
			current.RegistrySHA256 = nil
			current.Freshness = "unknown"
			if len(current.Reasons) == 0 {
				current.Reasons = previous.Reasons
			}
			result.Freshness = "unknown"
		} else if *previous.RegistrySHA256 != *current.RegistrySHA256 {
			current.RegistrySHA256 = nil
			current.Freshness = "stale"
			current.Reasons = append(current.Reasons, "registered facts changed during observation")
			if result.Freshness != "unknown" {
				result.Freshness = "stale"
			}
		}
		result.Registers = append(result.Registers, current)
	}
	if result.Freshness == "fresh" {
		// Workspace marker participates so tokens cannot migrate between roots.
		b, _ := json.Marshal(struct {
			Marker    string
			Registers []RegisterChange
		}{w.MarkerSHA256, result.Registers})
		digest := fmt.Sprintf("sha256:%x", sha256.Sum256(b))
		result.RegistrySHA256 = &digest
	}
	return result, nil
}

func fingerprintOwner(root string, owner changeOwner) RegisterChange {
	result := RegisterChange{ID: owner.id, Freshness: "fresh", Reasons: []string{}}
	manifest := [][2]string{}
	fail := func(err error) RegisterChange {
		result.RegistrySHA256 = nil
		result.Freshness = "unknown"
		result.Reasons = append(result.Reasons, err.Error())
		return result
	}
	// Root-relative opens also prevent a concurrently replaced parent directory
	// from redirecting a reader outside the observed workspace.
	bounded, err := os.OpenRoot(root)
	if err != nil {
		return fail(fmt.Errorf("open workspace: %w", err))
	}
	defer func() { _ = bounded.Close() }()
	readFile := func(relative string, info fs.FileInfo) error {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("%s is not a regular file", relative)
		}
		f, err := bounded.Open(filepath.Join(".ply", filepath.FromSlash(relative)))
		if err != nil {
			return fmt.Errorf("read %s: %w", relative, err)
		}
		defer func() { _ = f.Close() }()
		opened, err := f.Stat()
		if err != nil {
			return fmt.Errorf("inspect %s: %w", relative, err)
		}
		if !os.SameFile(info, opened) {
			return fmt.Errorf("%s changed during open", relative)
		}
		const limit = 16 << 20
		if opened.Size() > limit {
			return fmt.Errorf("%s exceeds the 16 MiB fingerprint limit", relative)
		}
		b, err := io.ReadAll(io.LimitReader(f, limit+1))
		if err != nil {
			return fmt.Errorf("read %s: %w", relative, err)
		}
		if len(b) > limit {
			return fmt.Errorf("%s exceeds the 16 MiB fingerprint limit", relative)
		}
		manifest = append(manifest, [2]string{relative, fmt.Sprintf("file:%d:sha256:%x", opened.Mode(), sha256.Sum256(b))})
		result.FileCount++
		return nil
	}
	for _, relative := range owner.files {
		info, err := inspectChangePath(bounded, relative)
		if errors.Is(err, fs.ErrNotExist) {
			manifest = append(manifest, [2]string{relative, "absent"})
			continue
		}
		if err != nil {
			return fail(err)
		}
		if err := readFile(relative, info); err != nil {
			return fail(err)
		}
	}
	var walk func(string) error
	walk = func(relative string) error {
		directory, err := bounded.Open(filepath.Join(".ply", filepath.FromSlash(relative)))
		if err != nil {
			return fmt.Errorf("open %s: %w", relative, err)
		}
		entries, err := directory.ReadDir(-1)
		closeErr := directory.Close()
		if err != nil {
			return fmt.Errorf("list %s: %w", relative, err)
		}
		if closeErr != nil {
			return fmt.Errorf("close %s: %w", relative, closeErr)
		}
		// File.ReadDir uses directory order, so sort before hashing the manifest.
		sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
		for _, entry := range entries {
			name := entry.Name()
			// Match the source stores' publication/lock exclusions. Arbitrary
			// dotfiles or .lock files can invalidate a strict event sequence.
			if strings.HasPrefix(name, ".publish-") || owner.id == "journal" && name == "append.lock" {
				continue
			}
			child := relative + "/" + name
			if entry.Type()&fs.ModeSymlink != 0 {
				return fmt.Errorf("%s is a symbolic link", child)
			}
			if entry.IsDir() {
				info, err := entry.Info()
				if err != nil {
					return fmt.Errorf("inspect %s: %w", child, err)
				}
				manifest = append(manifest, [2]string{child, fmt.Sprintf("directory:%d", info.Mode())})
				if err := walk(child); err != nil {
					return err
				}
				continue
			}
			if !strings.HasSuffix(name, ".json") {
				// Preserve unexpected source entries as structure, without reading
				// arbitrary output streams. They change run/journal read validity.
				if owner.id == "journal" || relative == "task-runs/v1/runs" || relative == "workflow-runs/requests" || strings.HasSuffix(relative, "/events") {
					info, err := entry.Info()
					if err != nil {
						return fmt.Errorf("inspect %s: %w", child, err)
					}
					manifest = append(manifest, [2]string{child, fmt.Sprintf("entry:%d", info.Mode())})
					result.FileCount++
				}
				continue
			}
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("inspect %s: %w", child, err)
			}
			if err := readFile(child, info); err != nil {
				return err
			}
		}
		return nil
	}
	for _, relative := range owner.trees {
		info, err := inspectChangePath(bounded, relative)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return fail(err)
		}
		if !info.IsDir() {
			return fail(fmt.Errorf("%s is not a directory", relative))
		}
		manifest = append(manifest, [2]string{relative, fmt.Sprintf("directory:%d", info.Mode())})
		if err := walk(relative); err != nil {
			return fail(err)
		}
	}
	b, _ := json.Marshal(manifest)
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256(b))
	result.RegistrySHA256 = &digest
	return result
}

// Check every component: Lstat of only the final file would follow a linked
// parent directory. The fixed relative inputs never contain .. or absolute paths.
func inspectChangePath(root *os.Root, relative string) (fs.FileInfo, error) {
	p := ""
	var info fs.FileInfo
	for _, part := range strings.Split(".ply/"+relative, "/") {
		p = filepath.Join(p, part)
		var err error
		info, err = root.Lstat(p)
		if err != nil {
			return nil, err
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			return nil, fmt.Errorf("%s contains a symbolic link", relative)
		}
	}
	return info, nil
}
