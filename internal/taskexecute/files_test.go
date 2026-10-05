package taskexecute

import (
	"os"
	"path/filepath"
	"testing"
)

func physicalTemp(t *testing.T) string {
	t.Helper()
	p, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestInstalledReplacementDoesNotChangeControl(t *testing.T) {
	root := physicalTemp(t)
	installed := filepath.Join(root, "installed")
	old := []byte("#!/bin/sh\necho old\n")
	if err := os.WriteFile(installed, old, 0700); err != nil {
		t.Fatal(err)
	}
	binding, err := bindExecutable(installed)
	if err != nil {
		t.Fatal(err)
	}
	control, err := PreserveControl(binding, filepath.Join(root, "run"))
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(installed, []byte("#!/bin/sh\necho new\n"), 0700); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(control.Path)
	if err != nil || string(b) != string(old) {
		t.Fatalf("control changed with install: %q, %v", b, err)
	}
	if _, err = PreserveControl(binding, filepath.Join(root, "run")); err == nil {
		t.Fatal("changed source was accepted")
	}
	if err = writeOnce(filepath.Join(root, "run", "intent.json"), []byte("first"), 0600); err != nil {
		t.Fatal(err)
	}
	if err = writeOnce(filepath.Join(root, "run", "intent.json"), []byte("second"), 0600); err == nil {
		t.Fatal("different retry replaced intent")
	}
}

func TestControlRejectsSymlinkAncestor(t *testing.T) {
	root := physicalTemp(t)
	actual := filepath.Join(root, "actual")
	if err := os.Mkdir(actual, 0700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "link")
	if err := os.Symlink(actual, link); err != nil {
		t.Fatal(err)
	}
	if err := writeOnce(filepath.Join(link, "missing", "intent.json"), []byte("data"), 0600); err == nil {
		t.Fatal("symlink ancestor accepted")
	}
	if _, err := os.Stat(filepath.Join(actual, "missing")); !os.IsNotExist(err) {
		t.Fatal("unsafe path was written")
	}
}
