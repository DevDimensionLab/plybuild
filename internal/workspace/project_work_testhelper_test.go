package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

// testWork is the work of Project id, or of the only Project, of the
// workspace at root. Tests reach registry and task-content paths only through
// it. A root without definitions or Projects is a bare store directory: its
// work is <root>/work of Project "test".
func testWork(t testing.TB, root string, id ...ProjectID) ProjectWork {
	t.Helper()
	if _, err := os.Lstat(filepath.Join(root, DefinitionsFile)); os.IsNotExist(err) {
		return ProjectWork{Workspace: root, Project: "test", Dir: filepath.Join(root, "work")}
	}
	works, err := DefinedWork(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(works) == 0 {
		return ProjectWork{Workspace: root, Project: "test", Dir: filepath.Join(root, "work")}
	}
	for _, work := range works {
		if len(id) == 1 && work.Project == id[0] {
			return work
		}
	}
	if len(id) == 0 && len(works) == 1 {
		return works[0]
	}
	t.Fatalf("workspace %s has no single Project %v among %v", root, id, works)
	return ProjectWork{}
}
