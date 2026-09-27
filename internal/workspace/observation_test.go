package workspace

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestObserveContainingAndRoot(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	dependencies := SystemDependencies()
	if _, err := os.Stat(root); err != nil {
		t.Fatal(err)
	}
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	if _, err := Init(dependencies); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(nested); err != nil {
		t.Fatal(err)
	}
	marker, err := os.ReadFile(filepath.Join(root, MarkerDirectory, MarkerFile))
	if err != nil {
		t.Fatal(err)
	}
	wantDigest := fmt.Sprintf("sha256:%x", sha256.Sum256(marker))
	for name, observe := range map[string]func() (WorkspaceObservation, error){
		"containing": func() (WorkspaceObservation, error) { return ObserveContaining(dependencies) },
		"root":       func() (WorkspaceObservation, error) { return ObserveRoot(dependencies, root) },
	} {
		t.Run(name, func(t *testing.T) {
			got, err := observe()
			if err != nil {
				t.Fatal(err)
			}
			if got.Root != root || got.MarkerFormatVersion != 1 || got.MarkerSHA256 != wantDigest {
				t.Fatalf("observation = %#v", got)
			}
		})
	}
}

func TestObserveRootRejectsNonRootAndMarkerDrift(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	dependencies := SystemDependencies()
	old, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(old) })
	_ = os.Chdir(root)
	if _, err := Init(dependencies); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(root, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := ObserveRoot(dependencies, nested); err == nil {
		t.Fatal("nested directory accepted as workspace root")
	}
	if err := os.WriteFile(filepath.Join(root, MarkerDirectory, MarkerFile), []byte("format_version: 1\nroot: /wrong\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ObserveRoot(dependencies, root); err == nil {
		t.Fatal("drifted marker accepted")
	}
}
