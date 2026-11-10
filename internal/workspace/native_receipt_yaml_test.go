package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNativeReceiptVersionRequiresYAMLInteger(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(f.wrapper, "task"), f.oid)
	if _, err := CreateTaskWorktree(f.dependencies, input); err != nil {
		t.Fatal(err)
	}
	r, err := f.dependencies.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	logical, err := encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	physical, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	for _, reader := range []struct {
		name  string
		bytes []byte
		read  func([]byte) (WorkItemRegistry, error)
	}{
		{"logical", logical, decodeWorkItemRegistry},
		{"store", physical, func(b []byte) (WorkItemRegistry, error) {
			return newSystemWorkItemStore().decodeStoredRegistry(root, b)
		}},
	} {
		t.Run(reader.name, func(t *testing.T) {
			const marker = "native_receipt_version: 1"
			if strings.Count(string(reader.bytes), marker) != 1 {
				t.Fatal("fixture must contain one actual writer marker")
			}
			for _, value := range []string{"1", "1.5", "1.9", "1.0", "1e0", "!!float 1", "null", "false", `"1"`, "0", "2", "-1"} {
				t.Run(value, func(t *testing.T) {
					b := []byte(strings.Replace(string(reader.bytes), marker, "native_receipt_version: "+value, 1))
					got, err := reader.read(b)
					if value == "1" {
						if err != nil || got.Tasks[0].Worktree.NativeReceiptVersion != nativeReceiptVersion {
							t.Fatalf("valid integer marker rejected: %v", err)
						}
					} else if err == nil {
						t.Fatalf("invalid raw YAML marker %s accepted as version %d", value, got.Tasks[0].Worktree.NativeReceiptVersion)
					}
				})
			}
		})
	}
}
