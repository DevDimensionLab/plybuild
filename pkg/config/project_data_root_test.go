package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProjectDataRootSurvivesConfigUpdate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "local-config.yaml")
	if err := os.WriteFile(path, []byte("projects:\n  fixture:\n    data_root: /physical/project-data\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg := OpenLocalConfig(dir)
	value, err := cfg.Config()
	if err != nil {
		t.Fatal(err)
	}
	if err = cfg.UpdateLocalConfig(value); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(raw), "data_root: /physical/project-data") {
		t.Fatalf("project data root lost during config update: %s, %v", raw, err)
	}
}
