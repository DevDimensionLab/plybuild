package workspace

import (
	"bytes"
	"reflect"
	"testing"
)

func TestProjectRegistryRoundTripWithoutDataPaths(t *testing.T) {
	registry := projectRegistry{FormatVersion: FormatVersion,
		Projects: []ProjectRecord{{ID: "ply", Name: "Ply", Wrapper: "/w/ply", RepoIDs: []RepoID{"ply"}}},
		Repos:    []RepoRecord{{ID: "ply", Locator: "/w/ply/main", GitCommonDir: "/w/ply/main/.git"}}}
	want := "format_version: 1\nprojects:\n  - id: ply\n    name: Ply\n    wrapper: /w/ply\n    repo_ids:\n      - ply\nrepos:\n  - id: ply\n    locator: /w/ply/main\n    git_common_dir: /w/ply/main/.git\n"
	raw, err := encodeProjectRegistry(registry)
	if err != nil || string(raw) != want {
		t.Fatalf("encoded %q, %v", raw, err)
	}
	decoded, err := decodeProjectRegistry(raw)
	if err != nil || !reflect.DeepEqual(decoded, registry) {
		t.Fatalf("decoded %+v, %v", decoded, err)
	}
	again, _ := encodeProjectRegistry(decoded)
	if !bytes.Equal(again, raw) {
		t.Fatalf("round trip changed bytes: %q", again)
	}
	clone := cloneProjects([]ProjectRecord{{ID: "ply", DataPaths: map[string]string{"backlog": "b"}}})
	if clone[0].DataPaths["backlog"] != "b" {
		t.Fatalf("clone %+v", clone)
	}
}
