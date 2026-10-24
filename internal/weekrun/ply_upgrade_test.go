package weekrun

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
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

// Reviewer-only fixture: fake launchctl launches exactly its serialized argv.
// It never calls real launchd or operates on the assigned product Run.
func TestPlyUpgradeLegacyWeekServiceStart(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("named launchd service requires macOS")
	}
	home, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	f := newFixture(t)
	legacy := weekPlist(f.b, []string{f.b.Control.Path})
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
	replacement := filepath.Join(home, "current-ply")
	if err := os.WriteFile(replacement, []byte("#!/bin/sh\nprintf current > \"$HOME/current-controller\"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	f.m.Executable = func() (string, error) { return replacement, nil }
	f.m.Launch = func(b Binding) error { return launchDetached(b, replacement) }
	f.m.Launchctl = func(args ...string) (string, error) {
		switch args[0] {
		case "print":
			return "state = waiting\n", nil
		case "kickstart", "bootstrap":
			cmd := exec.Command(f.b.Control.Path, "workflow", "week", "_supervise", f.b.ID, "--workspace", f.c.Workspace)
			cmd.Dir = f.c.ReturnWorktree
			out, err := cmd.CombinedOutput()
			return string(out), fmt.Errorf("synthetic launchctl used unchanged legacy ProgramArguments: %w", err)
		}
		return "", fmt.Errorf("unexpected launchctl: %v", args)
	}
	if _, err := f.m.Service(f.c.Workspace, f.b.ID, "start"); err != nil {
		t.Fatalf("ordinary upgraded CLI still starts removed historical controller: %v", err)
	}
	marker := filepath.Join(home, "current-controller")
	until := time.Now().Add(time.Second)
	for {
		if data, err := os.ReadFile(marker); err == nil && string(data) == "current" {
			break
		}
		if time.Now().After(until) {
			t.Fatal("current controller did not execute")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if after, err := os.ReadFile(path); err != nil || !bytes.Equal(after, legacy) {
		t.Fatalf("service start rewrote legacy plist: %v", err)
	}
}
