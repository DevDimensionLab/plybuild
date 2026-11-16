package workspace

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// physicalRegistry wraps logical registry bytes in the single stored format.
// Fixtures without integration evidence or publications need no objects.
func physicalRegistry(logical []byte) []byte {
	out := []byte(fmt.Sprintf("format_version: %d\nregistry:\n", RegistryFormat))
	for _, line := range strings.SplitAfter(string(logical), "\n") {
		if line != "" {
			out = append(out, "  "+line...)
		}
	}
	return out
}

func TestRegistryDiagnosticsUseDeclaredFormatAndValidation(t *testing.T) {
	logical, err := encodeWorkItemRegistry(emptyWorkItemRegistry())
	if err != nil {
		t.Fatal(err)
	}
	valid := string(physicalRegistry(logical))
	for name, tc := range map[string]struct{ input, code string }{
		"future":             {strings.Replace(valid, "version: 11", "version: 99", 1), "registry_unsupported_format"},
		"previous envelope":  {strings.Replace(valid, "version: 11", "version: 10", 1), "registry_unsupported_format"},
		"inline registry":    {"format_version: 4\n" + string(logical), "registry_unsupported_format"},
		"authorization":      {valid + "  acceptance_override: true\n", "registry_unknown_field"},
		"syntax":             {valid + "broken: [", "registry_malformed_yaml"},
		"missing version":    {strings.TrimPrefix(valid, "format_version: 11\n"), "registry_invalid_content"},
		"wrong version type": {strings.Replace(valid, "version: 11", "version: '11'", 1), "registry_invalid_content"},
		"duplicate":          {valid + "format_version: 11\n", "registry_invalid_content"},
		"multiple documents": {valid + "---\nformat_version: 11\n", "registry_malformed_yaml"},
		"no envelope":        {"format_version: 11\n" + string(logical), "registry_invalid_content"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := newSystemWorkItemStore().decodeStoredRegistry(testWork(t, t.TempDir()), []byte(tc.input))
			d := RegistryError(err)
			if d == nil || d.Code != tc.code || d.NextAction == "" || len(d.SupportedFormats) != 1 || d.SupportedFormats[0] != RegistryFormat {
				t.Fatalf("diagnostic: %+v %v", d, err)
			}
		})
	}
	f := newWorkItemJourneyFixture(t)
	r, _ := f.dependencies.WorkItems.Snapshot(testWork(t, filepath.Dir(f.wrapper)))
	b, _ := encodeWorkItemRegistry(r)
	b = bytes.Replace(b, []byte("parent_epic_id: epic"), []byte("parent_epic_id: missing"), 1)
	_, err = decodeWorkItemRegistry(b)
	if d := RegistryError(err); d == nil || d.Code != "registry_invalid_content" || d.FormatVersion != RegistryFormat {
		t.Fatalf("invalid reference: %v", err)
	}
}

// A fresh workspace has no registry file. Its first write is the single format,
// and adopt, Task create and Epic create work without any upgrade step.
func TestFreshRegistrySupportsAdoptTaskAndEpicCreate(t *testing.T) {
	f, in := epicCreateFixture(t)
	root := filepath.Dir(f.wrapper)
	raw, err := os.ReadFile(testWork(t, root).Registry())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(raw, []byte("format_version: 11\nregistry:\n")) {
		t.Fatalf("adopt and Task create wrote another format:\n%s", raw)
	}
	r, err := f.dependencies.WorkItems.Snapshot(testWork(t, root))
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Epics) != 1 || len(r.EpicBaseVersions) != 1 || len(r.Tasks) != 1 || len(r.TaskContentPublications) != 1 || r.EpicCreateOperations == nil {
		t.Fatalf("fresh registry = %#v", r)
	}
	created, err := CreateEpic(f.dependencies, in)
	if err != nil || created.Outcome != "created" {
		t.Fatalf("Epic create on a fresh registry: %#v %v", created, err)
	}
	if _, err = CreateTask(f.dependencies, TaskCreateInput{TaskID: "later", Title: "Later", Description: "On the created Epic", ParentEpicID: in.EpicID, ProjectID: "ply", RepoID: "ply"}); err != nil {
		t.Fatal(err)
	}
	if raw, err = os.ReadFile(testWork(t, root).Registry()); err != nil || !bytes.HasPrefix(raw, []byte("format_version: 11\nregistry:\n")) {
		t.Fatalf("registry format changed: %v\n%s", err, raw)
	}
}

func TestOtherRegistryFormatIsRejectedWithoutWrites(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	raw, err := os.ReadFile(testWork(t, root).Registry())
	if err != nil {
		t.Fatal(err)
	}
	for _, version := range []string{"10", "12"} {
		t.Run(version, func(t *testing.T) {
			other := bytes.Replace(raw, []byte("format_version: 11\n"), []byte("format_version: "+version+"\n"), 1)
			if err := os.WriteFile(testWork(t, root).Registry(), other, 0o600); err != nil {
				t.Fatal(err)
			}
			_, err := CreateTask(f.dependencies, TaskCreateInput{TaskID: "rejected", Title: "Rejected", Description: "No write", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"})
			if d := RegistryError(err); d == nil || d.Code != "registry_unsupported_format" || d.FormatVersion == RegistryFormat {
				t.Fatalf("other format accepted: %v", err)
			}
			after, err := os.ReadFile(testWork(t, root).Registry())
			if err != nil || !bytes.Equal(after, other) {
				t.Fatalf("rejection changed the registry: %v", err)
			}
		})
	}
}

// storedRegistryBytes encodes a fixture registry in the stored format,
// publishing its evidence objects exactly as the store does.
func storedRegistryBytes(t *testing.T, root string, r WorkItemRegistry) []byte {
	t.Helper()
	s := newSystemWorkItemStore()
	old, err := s.Snapshot(testWork(t, root))
	if err != nil {
		t.Fatal(err)
	}
	raw, err := s.encodeStoredRegistry(testWork(t, root), &r, old)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}
