// Package buildinfo exposes the identity of the executable without reading a workspace.
package buildinfo

import (
	"crypto/sha256"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
)

// Release builds set these with -ldflags. Development builds remain identifiable.
var Version = "dev"
var SourceCommit string

type Identity struct {
	Kind            string `json:"kind"`
	SchemaVersion   int    `json:"schema_version"`
	Version         string `json:"version"`
	SourceCommit    string `json:"source_commit"`
	Dirty           bool   `json:"dirty"`
	OS              string `json:"os"`
	Arch            string `json:"arch"`
	SHA256          string `json:"sha256"`
	RegistryFormats []int  `json:"registry_formats"`
}

func Current() (Identity, error) {
	i := Identity{Kind: "PlyExecutable@1", SchemaVersion: 1, Version: Version, SourceCommit: SourceCommit, OS: runtime.GOOS, Arch: runtime.GOARCH, RegistryFormats: []int{1, 2, 3, 4, 5}}
	if b, ok := debug.ReadBuildInfo(); ok {
		for _, s := range b.Settings {
			if s.Key == "vcs.revision" && i.SourceCommit == "" {
				i.SourceCommit = s.Value
			}
			if s.Key == "vcs.modified" {
				i.Dirty = s.Value == "true"
			}
		}
	}
	p, err := os.Executable()
	if err != nil {
		return i, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return i, err
	}
	i.SHA256 = fmt.Sprintf("sha256:%x", sha256.Sum256(b))
	return i, nil
}
