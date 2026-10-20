package workspace

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// The previous KnownFields decoder remains the diagnostic oracle. Include flow
// mappings so line sorting cannot accidentally substitute for traversal order.
func TestRegistryNodeDecodeMatchesStrictDiagnostics(t *testing.T) {
	for _, input := range []string{
		"{format_version: 1, epics: [], tasks: [], worktree_operations: []}",
		"{format_version: 1, epics: wrong, unknown: true, tasks: [], worktree_operations: []}",
		"{format_version: 1, unknown: true, epics: wrong, tasks: [], worktree_operations: []}",
		"format_version: 1\nepics: [{id: [], extra: true, title: {bad: value}}]\ntasks: []\nworktree_operations: []\n",
		"format_version: 1\nepics: [{extra: true, id: [], title: {bad: value}}]\ntasks: []\nworktree_operations: []\n",
		"format_version: 1\nepics: null\ntasks: []\nworktree_operations: []\n",
	} {
		t.Run(input, func(t *testing.T) {
			decoder := yaml.NewDecoder(strings.NewReader(input))
			decoder.KnownFields(true)
			var old, got WorkItemRegistry
			wantErr := decoder.Decode(&old)
			root, err := parseRegistryYAML([]byte(input))
			if err != nil {
				t.Fatal(err)
			}
			gotErr := decodeRegistryNode(root, &got)
			if (wantErr == nil) != (gotErr == nil) || (wantErr != nil && wantErr.Error() != gotErr.Error()) {
				t.Fatalf("diagnostics differ:\nwant %v\ngot %v", wantErr, gotErr)
			}
			if wantErr == nil && !reflect.DeepEqual(old, got) {
				t.Fatal("typed values differ")
			}
		})
	}
}

func TestRegistryTaggedCollectionKeysPreserveStrictDiagnostics(t *testing.T) {
	valid := "format_version: 1\nepics: []\ntasks: []\nworktree_operations: []\n"
	for name, input := range map[string]string{
		"sequence key":           valid + "? !!str [bad]\n: true\n",
		"mapping key":            valid + "? !!str {bad: value}\n: true\n",
		"key before unknown":     valid + "? !!str [bad]\n: {ignored: true}\nunknown: true\n",
		"unknown before key":     valid + "unknown: true\n? !!str [bad]\n: true\n",
		"type before key":        strings.Replace(valid, "epics: []", "epics: wrong", 1) + "? !!str [bad]\n: true\n",
		"nested key before type": "format_version: 1\nepics:\n- ? !!str [bad]\n  : true\n  id: []\ntasks: []\nworktree_operations: []\n",
	} {
		t.Run(name, func(t *testing.T) {
			decoder := yaml.NewDecoder(strings.NewReader(input))
			decoder.KnownFields(true)
			var old WorkItemRegistry
			wantErr := decoder.Decode(&old)
			if wantErr == nil {
				t.Fatal("strict decoder must reject a collection key")
			}
			t.Run("oracle", func(t *testing.T) {
				root, err := parseRegistryYAML([]byte(input))
				if err != nil {
					t.Fatal(err)
				}
				var got WorkItemRegistry
				if err := decodeRegistryNode(root, &got); err == nil || err.Error() != wantErr.Error() {
					t.Fatalf("diagnostics differ:\nwant %v\ngot %v", wantErr, err)
				}
			})
			t.Run("consumed snapshot", func(t *testing.T) {
				root := t.TempDir()
				if err := writeRegistryDecodeFixture(root, []byte(input)); err != nil {
					t.Fatal(err)
				}
				_, err := newSystemWorkItemStore().Snapshot(root)
				got, want := RegistryError(err), RegistryError(classifyRegistryDecode(wantErr, 1))
				if got == nil || got.Code != want.Code || got.Detail != want.Detail || got.FormatVersion != 1 || got.Artifact != workItemsPath(root) {
					t.Fatalf("consumed diagnostic differs:\nwant %+v\ngot %+v (%v)", want, got, err)
				}
			})
		})
	}
}

func TestRegistryConsumedSnapshotRejectsInvalidSyntaxAndShape(t *testing.T) {
	valid := []byte("format_version: 1\nepics: []\ntasks: []\nworktree_operations: []\n")
	for _, tc := range []struct {
		name  string
		input []byte
		code  string
	}{
		{"alias", []byte("format_version: 1\nepics: &a []\ntasks: *a\nworktree_operations: []\n"), "registry_invalid_content"},
		{"merge", append(append([]byte{}, valid...), []byte("<<: {}\n")...), "registry_invalid_content"},
		{"null", bytes.Replace(valid, []byte("tasks: []"), []byte("tasks: null"), 1), "registry_invalid_content"},
		{"missing", bytes.Replace(valid, []byte("tasks: []\n"), nil, 1), "registry_invalid_content"},
		{"order", []byte("format_version: 1\ntasks: []\nepics: []\nworktree_operations: []\n"), "registry_invalid_content"},
		{"second document", append(append([]byte{}, valid...), []byte("---\n{}\n")...), "registry_malformed_yaml"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if err := writeRegistryDecodeFixture(root, tc.input); err != nil {
				t.Fatal(err)
			}
			_, err := newSystemWorkItemStore().Snapshot(root)
			d := RegistryError(err)
			if d == nil || d.Code != tc.code || d.Artifact != workItemsPath(root) {
				t.Fatalf("diagnostic lost code/artifact: %+v (%v)", d, err)
			}
		})
	}
}

func writeRegistryDecodeFixture(root string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(workItemsPath(root)), 0755); err != nil {
		return err
	}
	return os.WriteFile(workItemsPath(root), b, 0600)
}
