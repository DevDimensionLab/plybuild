package clirelease

import (
	"strings"
	"testing"
)

func TestManifestDeclaresOnlyTheSingleRegistryFormat(t *testing.T) {
	digest := "sha256:" + strings.Repeat("a", 64)
	for _, tc := range []struct {
		formats []int
		valid   bool
	}{
		{[]int{11}, true},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, false},
		{[]int{10, 11}, false},
		{[]int{}, false},
	} {
		m := Manifest{Kind: "PlyRelease@1", SchemaVersion: 1, Version: "v1.2.3", SourceCommit: strings.Repeat("c", 40), OS: "darwin", Arch: "arm64", SHA256: digest, SizeBytes: 1, RegistryFormats: tc.formats, InstallerSHA256: digest}
		if err := m.Validate(); (err == nil) != tc.valid {
			t.Fatalf("registry_formats %v: %v", tc.formats, err)
		}
	}
}
