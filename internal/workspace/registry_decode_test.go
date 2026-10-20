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
