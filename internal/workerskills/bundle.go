// Package workerskills provides pinned, offline Addy guidance adapted to the
// native Ply worker roles. It never reads a live source checkout or user config.
package workerskills

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"
)

const (
	SourceRepository = "https://github.com/addyosmani/agent-skills"
	SourceCommit     = "1401c8b8030e023baeebb31781a6653fe8e93026"
	ManifestSHA256   = "29466c59aa8bc8c3a6eed0abc3765d009cd0d1095df4785f6fb948831e703561"

	RolePlanning   = "planning"
	RoleDeveloper  = "developer"
	RoleCorrection = "correction"
	RoleReviewer   = "reviewer"
)

//go:embed upstream adapters manifest.json
var embedded embed.FS

// ManifestData identifies the exact upstream source, Ply adapters, and ordered
// role selections included in the executable. Manifest returns a fresh value.
type ManifestData struct {
	Format           int                 `json:"format"`
	SourceRepository string              `json:"source_repository"`
	SourceCommit     string              `json:"source_commit"`
	License          string              `json:"license"`
	LicensePath      string              `json:"license_path"`
	Files            []File              `json:"files"`
	Roles            map[string]RoleData `json:"roles"`
}

// File hashes the exact embedded bytes. SourcePath is the upstream path at
// SourceCommit; its absence identifies a Ply-owned adapter.
type File struct {
	Path       string `json:"path"`
	SourcePath string `json:"source_path,omitempty"`
	SHA256     string `json:"sha256"`
	Bytes      int    `json:"bytes"`
}

type RoleData struct {
	Adapters []string `json:"adapters"`
	Files    []string `json:"files"`
}

// Manifest verifies and returns the embedded provenance manifest.
func Manifest() (ManifestData, error) {
	m, _, err := verifiedManifest(embedded)
	return m, err
}

// Verify rejects changed, missing, or unmanifested embedded material, including
// changes to the role selections or adapter text. This is a consistency check
// against the reviewed pin, not a substitute for qualifying the executable.
func Verify() error {
	_, err := Manifest()
	return err
}

// Bundle returns deterministic instructions for exactly one role. Prepend the
// complete result to the controller's prompt before hashing and preserving that
// prompt in the native attempt. Unknown roles and invalid assets fail closed.
func Bundle(role string) (string, error) {
	return bundle(embedded, role)
}

func bundle(files fs.FS, role string) (string, error) {
	m, manifest, err := verifiedManifest(files)
	if err != nil {
		return "", err
	}
	selection, ok := m.Roles[role]
	if !ok {
		return "", fmt.Errorf("worker skills: unknown role %q", role)
	}
	var out strings.Builder
	fmt.Fprintf(&out, "# Pinned Ply worker role package\n\nRole: %s\nSource: %s\nSource commit: %s\nManifest SHA-256: %s\n\n", role, m.SourceRepository, m.SourceCommit, ManifestSHA256)
	for _, name := range selection.Adapters {
		if err := appendSection(&out, files, name); err != nil {
			return "", err
		}
	}
	out.WriteString("## Embedded provenance manifest\n\n```json\n")
	out.Write(manifest)
	out.WriteString("```\n\n## Upstream material\n\nThe following exact source bytes are interpreted under the Ply adapter above.\n")
	for _, name := range selection.Files {
		if err := appendSection(&out, files, name); err != nil {
			return "", err
		}
	}
	out.WriteString("\nEnd of embedded upstream material. The Ply adapter and the following controller task instruction govern this attempt.\n\n")
	return out.String(), nil
}

func appendSection(out *strings.Builder, files fs.FS, name string) error {
	raw, err := fs.ReadFile(files, name)
	if err != nil {
		return fmt.Errorf("worker skills: read %s: %w", name, err)
	}
	fmt.Fprintf(out, "\n--- BEGIN EMBEDDED FILE: %s ---\n", name)
	out.Write(raw)
	fmt.Fprintf(out, "\n--- END EMBEDDED FILE: %s ---\n", name)
	return nil
}

func verifiedManifest(files fs.FS) (ManifestData, []byte, error) {
	var m ManifestData
	raw, err := fs.ReadFile(files, "manifest.json")
	if err != nil {
		return m, nil, fmt.Errorf("worker skills: read manifest: %w", err)
	}
	if digest(raw) != ManifestSHA256 {
		return m, nil, fmt.Errorf("worker skills: manifest SHA-256 differs from the pinned manifest")
	}
	if err := json.Unmarshal(raw, &m); err != nil {
		return m, nil, fmt.Errorf("worker skills: decode manifest: %w", err)
	}
	if m.Format != 1 || m.SourceRepository != SourceRepository || m.SourceCommit != SourceCommit || m.License != "MIT" {
		return m, nil, fmt.Errorf("worker skills: unexpected source provenance")
	}
	listed := map[string]bool{"manifest.json": true}
	for _, file := range m.Files {
		if !fs.ValidPath(file.Path) || listed[file.Path] {
			return m, nil, fmt.Errorf("worker skills: invalid or duplicate manifest path %q", file.Path)
		}
		content, err := fs.ReadFile(files, file.Path)
		if err != nil {
			return m, nil, fmt.Errorf("worker skills: read %s: %w", file.Path, err)
		}
		if len(content) != file.Bytes || digest(content) != file.SHA256 {
			return m, nil, fmt.Errorf("worker skills: integrity mismatch for %s", file.Path)
		}
		listed[file.Path] = true
	}
	// Fail on accidentally embedded additions, such as upstream AGENTS.md, even
	// when they would not otherwise be selected by a role.
	err = fs.WalkDir(files, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && !listed[name] {
			return fmt.Errorf("worker skills: unmanifested embedded file %s", name)
		}
		return nil
	})
	if err != nil {
		return m, nil, err
	}
	for role, selection := range m.Roles {
		for _, name := range append(append([]string{}, selection.Adapters...), selection.Files...) {
			if !listed[name] {
				return m, nil, fmt.Errorf("worker skills: role %s selects unmanifested file %s", role, name)
			}
		}
	}
	return m, raw, nil
}

func digest(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
