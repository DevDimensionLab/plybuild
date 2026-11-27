package workspaceview

import (
	"os"
	"path/filepath"
	"testing"
)

// The project list carries the revision of the projects register, the same
// value workspace status reports, so a reader caches against one truth.
func TestProjectListCarriesTheProjectsRegisterRevision(t *testing.T) {
	d, root := changesFixture(t)
	defineChangeProject(t, root)
	revision := func() string {
		t.Helper()
		listed, err := ListProjects(d.Dependencies)
		if err != nil {
			t.Fatal(err)
		}
		status, err := ReadChanges(d.Dependencies)
		if err != nil {
			t.Fatal(err)
		}
		if listed.RegistrySHA256 == nil || status.Registers[0].ID != "projects" || status.Registers[0].RegistrySHA256 == nil {
			t.Fatalf("no revision: list %v, status %+v", listed.RegistrySHA256, status.Registers[0])
		}
		if *listed.RegistrySHA256 != *status.Registers[0].RegistrySHA256 {
			t.Fatalf("list revision %s; status has %s", *listed.RegistrySHA256, *status.Registers[0].RegistrySHA256)
		}
		return *listed.RegistrySHA256
	}
	before := revision()
	definitions := filepath.Join(root, "ply-workspace.json")
	raw, err := os.ReadFile(definitions)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(definitions, append(raw, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	if revision() == before {
		t.Fatal("changed definitions kept the revision")
	}
}
