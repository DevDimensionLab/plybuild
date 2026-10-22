// Package clirelease validates and installs a prebuilt CLI independently of Ply stores.
package clirelease

import (
	"bytes"
	"context"
	"crypto/sha256"
	"debug/macho"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"time"

	"github.com/devdimensionlab/plybuild/internal/buildinfo"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

type Manifest struct {
	Qualification   *QualifiedCLI `json:"qualification,omitempty"`
	Kind            string        `json:"kind"`
	SchemaVersion   int           `json:"schema_version"`
	Version         string        `json:"version"`
	SourceCommit    string        `json:"source_commit"`
	OS              string        `json:"os"`
	Arch            string        `json:"arch"`
	SHA256          string        `json:"sha256"`
	SizeBytes       int64         `json:"size_bytes"`
	RegistryFormats []int         `json:"registry_formats"`
	InstallerSHA256 string        `json:"installer_sha256"`
}

var digestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
var commitPattern = regexp.MustCompile(`^[0-9a-f]{40}([0-9a-f]{24})?$`)
var versionPattern = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+[A-Za-z0-9.+-]*$`)

func Digest(b []byte) string { return fmt.Sprintf("sha256:%x", sha256.Sum256(b)) }

func DecodeStrict(b []byte, out any) error {
	if _, err := canonicaljson.DecodeStrict(b); err != nil {
		return err
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if err := d.Decode(out); err != nil {
		return err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return fmt.Errorf("expected exactly one JSON document")
	}
	return nil
}

func (m Manifest) Validate() error {
	if m.Qualification != nil || m.Kind != "PlyRelease@1" || m.SchemaVersion != 1 || !versionPattern.MatchString(m.Version) || !commitPattern.MatchString(m.SourceCommit) || m.OS != "darwin" || m.Arch != "arm64" || !digestPattern.MatchString(m.SHA256) || !digestPattern.MatchString(m.InstallerSHA256) || m.SizeBytes <= 0 || m.SizeBytes > 128<<20 {
		return fmt.Errorf("release_invalid_manifest: expected an identified macOS arm64 release")
	}
	if !reflect.DeepEqual(m.RegistryFormats, []int{1, 2, 3, 4, 5}) && !reflect.DeepEqual(m.RegistryFormats, []int{1, 2, 3, 4, 5, 6}) {
		return fmt.Errorf("release_invalid_manifest: unsupported registry capability declaration")
	}
	return nil
}

// ReadRegular rejects symlinks and bounds memory use. It does not touch a workspace.
func ReadRegular(path string, limit int64) ([]byte, error) {
	st, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !st.Mode().IsRegular() || st.Size() > limit {
		return nil, fmt.Errorf("not a bounded regular file: %s", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(st, opened) {
		return nil, fmt.Errorf("file changed before read: %s", path)
	}
	b, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil {
		return nil, err
	}
	after, err := os.Lstat(path)
	if err != nil || !os.SameFile(opened, after) || len(b) > int(limit) || int64(len(b)) != opened.Size() || after.ModTime() != opened.ModTime() {
		return nil, fmt.Errorf("file changed during read: %s", path)
	}
	return b, nil
}

func ValidateExecutable(path string, b []byte, m Manifest) error {
	if m.OS != runtime.GOOS || m.Arch != runtime.GOARCH {
		return fmt.Errorf("release_platform_mismatch: release is %s/%s, host is %s/%s", m.OS, m.Arch, runtime.GOOS, runtime.GOARCH)
	}
	if int64(len(b)) != m.SizeBytes || Digest(b) != m.SHA256 {
		return fmt.Errorf("release_digest_mismatch: candidate bytes differ from the manifest")
	}
	f, err := macho.NewFile(bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("release_platform_mismatch: not a Mach-O executable: %w", err)
	}
	defer f.Close()
	if f.Cpu != macho.CpuArm64 || f.Type != macho.TypeExec {
		return fmt.Errorf("release_platform_mismatch: expected an arm64 executable")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, path, "version", "--format", "json")
	c.Env = []string{"PATH=/usr/bin:/bin", "LC_ALL=C", "LANG=C"}
	c.Dir = filepath.Dir(path)
	out, err := c.Output()
	if err != nil {
		return fmt.Errorf("release_identity_failed: %w", err)
	}
	var i buildinfo.Identity
	if err := DecodeStrict(out, &i); err != nil {
		return err
	}
	if i.Kind != "PlyExecutable@1" || i.SchemaVersion != 1 || i.Dirty || i.Version != m.Version || i.SourceCommit != m.SourceCommit || i.OS != m.OS || i.Arch != m.Arch || i.SHA256 != m.SHA256 || !reflect.DeepEqual(i.RegistryFormats, m.RegistryFormats) {
		return fmt.Errorf("release_identity_mismatch: executable identity differs from the manifest")
	}
	return nil
}

// Load verifies a package's caller-pinned manifest and delivered bytes. It does
// not invoke the old installed CLI, use Go/Git, or inspect workspace state.
func Load(directory, expected string) (Manifest, []byte, error) {
	var m Manifest
	b, err := ReadRegular(filepath.Join(directory, "manifest.json"), 64<<10)
	if err != nil {
		return m, nil, err
	}
	if !digestPattern.MatchString(expected) || Digest(b) != expected {
		return m, nil, fmt.Errorf("release_digest_mismatch: pin the exact manifest SHA-256")
	}
	if err = DecodeStrict(b, &m); err != nil {
		return m, nil, err
	}
	if err = m.Validate(); err != nil {
		return m, nil, err
	}
	b, err = ReadRegular(filepath.Join(directory, "ply"), 128<<20)
	if err != nil {
		return m, nil, err
	}
	if err = ValidateExecutable(filepath.Join(directory, "ply"), b, m); err != nil {
		return m, nil, err
	}
	installer, err := ReadRegular(filepath.Join(directory, "ply-install"), 128<<20)
	if err != nil {
		return m, nil, err
	}
	if Digest(installer) != m.InstallerSHA256 {
		return m, nil, fmt.Errorf("release_digest_mismatch: installer bytes differ")
	}
	return m, b, nil
}
