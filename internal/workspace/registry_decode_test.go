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
		"{epics: [], tasks: [], worktree_operations: []}",
		"{epics: wrong, unknown: true, tasks: [], worktree_operations: []}",
		"{unknown: true, epics: wrong, tasks: [], worktree_operations: []}",
		"epics: [{id: [], extra: true, title: {bad: value}}]\ntasks: []\nworktree_operations: []\n",
		"epics: [{extra: true, id: [], title: {bad: value}}]\ntasks: []\nworktree_operations: []\n",
		"epics: null\ntasks: []\nworktree_operations: []\n",
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
	valid := "epics: []\ntasks: []\nworktree_operations: []\n"
	for name, input := range map[string]string{
		"sequence key":           valid + "? !!str [bad]\n: true\n",
		"mapping key":            valid + "? !!str {bad: value}\n: true\n",
		"key before unknown":     valid + "? !!str [bad]\n: {ignored: true}\nunknown: true\n",
		"unknown before key":     valid + "unknown: true\n? !!str [bad]\n: true\n",
		"type before key":        strings.Replace(valid, "epics: []", "epics: wrong", 1) + "? !!str [bad]\n: true\n",
		"nested key before type": "epics:\n- ? !!str [bad]\n  : true\n  id: []\ntasks: []\nworktree_operations: []\n",
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
				physical := physicalRegistry([]byte(input))
				if err := writeRegistryDecodeFixture(t, root, physical); err != nil {
					t.Fatal(err)
				}
				// The store decodes the logical mapping re-encoded from the envelope.
				node, err := parseRegistryYAML(physical)
				if err != nil {
					t.Fatal(err)
				}
				logical, err := marshalRegistryNode(registryNodeField(node, "registry"))
				if err != nil {
					t.Fatal(err)
				}
				oracle := yaml.NewDecoder(bytes.NewReader(logical))
				oracle.KnownFields(true)
				var stored WorkItemRegistry
				_, err = newSystemWorkItemStore().Snapshot(testWork(t, root))
				got, want := RegistryError(err), RegistryError(classifyRegistryDecode(oracle.Decode(&stored), RegistryFormat))
				if got == nil || want == nil || got.Code != want.Code || got.Detail != want.Detail || got.FormatVersion != RegistryFormat || got.Artifact != testWork(t, root).Registry() {
					t.Fatalf("consumed diagnostic differs:\nwant %+v\ngot %+v (%v)", want, got, err)
				}
			})
		})
	}
}

func TestRegistryConsumedSnapshotRejectsInvalidSyntaxAndShape(t *testing.T) {
	valid, err := encodeWorkItemRegistry(emptyWorkItemRegistry())
	if err != nil {
		t.Fatal(err)
	}
	wrap := func(b []byte) []byte { return physicalRegistry(b) }
	for _, tc := range []struct {
		name  string
		input []byte
		code  string
	}{
		{"alias", wrap(bytes.Replace(bytes.Replace(valid, []byte("epics: []"), []byte("epics: &a []"), 1), []byte("tasks: []"), []byte("tasks: *a"), 1)), "registry_invalid_content"},
		{"merge", wrap(append(append([]byte{}, valid...), []byte("<<: {}\n")...)), "registry_invalid_content"},
		{"null", wrap(bytes.Replace(valid, []byte("tasks: []"), []byte("tasks: null"), 1)), "registry_invalid_content"},
		{"missing", wrap(bytes.Replace(valid, []byte("tasks: []\n"), nil, 1)), "registry_invalid_content"},
		{"order", wrap(bytes.Replace(valid, []byte("epics: []\ntasks: []\n"), []byte("tasks: []\nepics: []\n"), 1)), "registry_invalid_content"},
		{"second document", append(wrap(valid), []byte("---\n{}\n")...), "registry_malformed_yaml"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if err := writeRegistryDecodeFixture(t, root, tc.input); err != nil {
				t.Fatal(err)
			}
			_, err := newSystemWorkItemStore().Snapshot(testWork(t, root))
			d := RegistryError(err)
			if d == nil || d.Code != tc.code || d.Artifact != testWork(t, root).Registry() {
				t.Fatalf("diagnostic lost code/artifact: %+v (%v)", d, err)
			}
		})
	}
}

func writeRegistryDecodeFixture(t *testing.T, root string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(testWork(t, root).Registry()), 0755); err != nil {
		return err
	}
	return os.WriteFile(testWork(t, root).Registry(), b, 0600)
}
