package weekrun

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestPlyUpgradeExistingWeek(t *testing.T) {
	f := newMR05Week(t)
	before, err := os.ReadFile(filepath.Join(f.dir, "binding.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(f.b.Control.Path); err != nil {
		t.Fatal(err)
	}
	p, err := f.m.Controls().Preview(f.c.Workspace, f.b.ID, "pause", "tc0e6e1fe-upgrade-pause")
	if err != nil || p.Request == nil {
		t.Fatalf("historical Ply blocked week control: %+v %v", p, err)
	}
	if _, err = f.m.Controls().Apply(*p.Request); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(filepath.Join(f.dir, "binding.json"))
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("week history rewritten: %v", err)
	}
}

func TestPlyUpgradeLegacyWeekServiceReadback(t *testing.T) {
	f := newFixture(t)
	f.b.Environment = map[string]string{"HOME": f.c.Workspace, "PATH": "/usr/bin:/bin"}
	legacy := bytes.Replace(Plist(f.b), []byte("<string>/usr/bin/env</string>\n<string>ply</string>"), []byte("<string>"+f.b.Control.Path+"</string>"), 1)
	path := servicePath(f.b)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, legacy, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(f.b.Control.Path); err != nil {
		t.Fatal(err)
	}
	if _, err := ServiceReadback(f.b, func(...string) (string, error) { return "state = running\n pid = 42\n", nil }); err != nil {
		t.Fatalf("historical runtime path invalidated service history: %v", err)
	}
	after, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(after, legacy) {
		t.Fatalf("service history rewritten: %v", err)
	}
}
