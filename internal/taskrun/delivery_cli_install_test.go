package taskrun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/buildinfo"
)

// Consumes the built candidate's public CLI and real native verify/integrate/
// closeout boundaries; the disposable fixture is not an ordinary installation.
func TestQualifiedCLIInstallationPublicConsumer(t *testing.T) {
	qualifiedCLIInstallationConsumer(t, automaticCLIBinary(t))
}

func qualifiedCLIInstallationConsumer(t *testing.T, binary string) {
	t.Helper()
	f, o := cleanupCLIFixture(t, binary)
	destination := filepath.Join(f.R.WorkspaceRoot, "installed-ply")
	old := []byte("#!/bin/sh\necho old\n")
	if err := os.WriteFile(destination, old, 0755); err != nil {
		t.Fatal(err)
	}
	var lastStderr bytes.Buffer
	call := 0
	invoke := func(controller, dest string, extra ...string) (CLIInstallation, []byte, error) {
		args := []string{"workflow", "execute", "install-cli", o.RunID, "--destination", dest, "--coordinator", "synthetic authorized coordinator", "--format", "json"}
		args = append(args, extra...)
		c := exec.Command(controller, args...)
		c.Dir = f.Parent
		lastStderr.Reset()
		c.Stderr = &lastStderr
		out, err := c.Output()
		call++
		if root := os.Getenv("PLY_AL16_CLI_EVIDENCE"); root != "" {
			dir := filepath.Join(root, t.Name())
			if e := os.MkdirAll(dir, 0700); e != nil {
				t.Fatal(e)
			}
			prefix := filepath.Join(dir, fmt.Sprintf("%02d", call))
			exit := -1
			if c.ProcessState != nil {
				exit = c.ProcessState.ExitCode()
			}
			argv, _ := json.Marshal(c.Args)
			for suffix, data := range map[string][]byte{".stdout": out, ".stderr": lastStderr.Bytes(), ".exit": []byte(fmt.Sprintln(exit)), ".argv.json": argv} {
				if e := os.WriteFile(prefix+suffix, data, 0600); e != nil {
					t.Fatal(e)
				}
			}
		}
		if exit, ok := err.(*exec.ExitError); ok {
			t.Logf("CLI actual exit=%d stderr=%s", exit.ExitCode(), lastStderr.String())
		}
		var r CLIInstallation
		if decodeErr := json.Unmarshal(out, &r); decodeErr != nil {
			t.Fatalf("stdout is not clean JSON: %v %s (process %v)", decodeErr, out, err)
		}
		return r, out, err
	}
	if _, raw, err := invoke(binary, destination); err == nil {
		t.Fatalf("unclosed owner installed CLI: %s", raw)
	}
	_, reflog := cleanupOwnerIntegrates(t, binary, f, o)
	cleanupOwnerLeaves(t, f)
	reviewRoundIn(t, "w-fixture:root")
	coordinator := f.D
	coordinator.CWD = func() (string, error) { return f.Parent, nil }
	coordinator.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: os.Getpid(), CWD: f.Parent}}, nil
	}
	closed, err := WorkflowFollow(coordinator, f.R.WorkspaceRoot, o.RunID, 30*time.Second)
	if err != nil {
		t.Fatalf("modeled fixture closeout: %v", err)
	}
	assertCleanupRemoved(t, f, closed, reflog)
	ready, raw, err := invoke(binary, destination, "--check")
	if err != nil || ready.Controller.Locator == "" {
		t.Fatalf("check: %v %s", err, raw)
	}
	// Bootstrap-free first installation uses the native retained snapshot after
	// Task deletion while the ordinary target still holds the old CLI.
	controller := ready.Controller.Locator
	before, err := os.ReadFile(destination)
	if err != nil || string(before) != string(old) {
		t.Fatal("check changed target")
	}
	installed, raw, err := invoke(controller, destination)
	if err != nil || installed.State != "verified" || installed.Smoke == nil || installed.Smoke.Exit == nil || *installed.Smoke.Exit != 0 {
		t.Fatalf("install: %v %s", err, raw)
	}
	t.Logf("actual installed smoke exit=%d stdout=%s install stderr=%s", *installed.Smoke.Exit, installed.Smoke.Stdout, lastStderr.String())
	var identity buildinfo.Identity
	if err := json.Unmarshal([]byte(installed.Smoke.Stdout), &identity); err != nil {
		t.Fatal(err)
	}
	if identity.SourceCommit == "" {
		for _, fact := range []string{"embedded source identity is unknown", "qualification_oid", installed.Source.OID, "qualification_tree", installed.Source.Tree, installed.Source.BinarySHA256, "installed executable", "continues"} {
			if !strings.Contains(lastStderr.String(), fact) {
				t.Fatalf("missing sourced warning %q: %s", fact, lastStderr.String())
			}
		}
	} else if strings.Contains(lastStderr.String(), "embedded source identity is unknown") {
		t.Fatalf("known identity warned as unknown: %s", lastStderr.String())
	}
	candidateBytes, err := os.ReadFile(controller)
	if err != nil {
		t.Fatal(err)
	}
	installedBytes, err := os.ReadFile(destination)
	if err != nil || !bytes.Equal(candidateBytes, installedBytes) {
		t.Fatal("qualified bytes changed", err)
	}
	backup := installed.Receipt.Backup
	repeated, raw, err := invoke(controller, destination)
	if err != nil || repeated.Receipt.Backup != backup || !equal(repeated.Receipt, installed.Receipt) {
		t.Fatalf("lost reply recovery %v %s", err, raw)
	}
	if strings.Contains(lastStderr.String(), "embedded source identity is unknown") {
		t.Fatalf("repeat warning storm: %s", lastStderr.String())
	}
	if _, raw, err := invoke(controller, destination+"-changed"); err == nil {
		t.Fatalf("changed destination reused operation %s", raw)
	}
	b, err := os.ReadFile(backup)
	if err != nil || string(b) != string(old) {
		t.Fatal("backup differs", err)
	}
	if err = os.WriteFile(destination, []byte("separately chosen later result"), 0755); err != nil {
		t.Fatal(err)
	}
	unknown, raw, err := invoke(controller, destination, "--observe")
	if err == nil || unknown.State != "effect_unknown" {
		t.Fatalf("unknown effect false pass %v %s", err, raw)
	}
	if _, raw, err = invoke(controller, destination, "--rollback"); err == nil {
		t.Fatalf("rollback overwrote later result %s", raw)
	}
	// Restore exact installed bytes to exercise guarded rollback, not a new install.
	b, err = os.ReadFile(controller)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(destination, b, 0755); err != nil {
		t.Fatal(err)
	}
	restored, raw, err := invoke(controller, destination, "--rollback")
	if err != nil || restored.State != "rolled_back" {
		t.Fatalf("rollback %v %s", err, raw)
	}
	b, err = os.ReadFile(destination)
	if err != nil || string(b) != string(old) {
		t.Fatal("rollback did not restore old executable")
	}
	if _, raw, err = invoke(controller, destination); err == nil || !strings.Contains(string(raw), "rolled_back") {
		t.Fatalf("rolled back operation retried %v %s", err, raw)
	}
}

func TestQualifiedCLIInstallationInterruptionAndKnownSmokeFailure(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := cleanupCLIFixture(t, binary)
	_, _ = cleanupOwnerIntegrates(t, binary, f, o)
	workflowTestModel(t, f.workflowFixture, map[string]any{"agent_status": "idle"})
	if _, raw, err := cleanupCLIAt(binary, f.Parent, "execute", "follow", o.RunID, "--timeout", "30"); err != nil {
		t.Fatalf("closeout %v %s", err, raw)
	}
	dest := filepath.Join(f.R.WorkspaceRoot, "installed-ply")
	if err := os.WriteFile(dest, []byte("old binary"), 0755); err != nil {
		t.Fatal(err)
	}
	d := f.D
	d.CWD = func() (string, error) { return f.Parent, nil }
	ready, err := InstallQualifiedCLI(d, f.R.WorkspaceRoot, o.RunID, dest, "test coordinator", "check")
	if err != nil {
		t.Fatal(err)
	}
	d.Executable = func() (string, error) { return ready.Controller.Locator, nil }
	d.Fault = func(stage string) error {
		if stage == "cli_install_before_effect" {
			return fmt.Errorf("injected before effect")
		}
		return nil
	}
	if _, err = InstallQualifiedCLI(d, f.R.WorkspaceRoot, o.RunID, dest, "test coordinator", "install"); err == nil {
		t.Fatal("missing before-effect failure")
	}
	before, _ := os.ReadFile(dest)
	if string(before) != "old binary" {
		t.Fatal("before-effect interruption replaced destination")
	}
	d.Fault = func(stage string) error {
		if stage == "cli_install_after_effect" {
			return fmt.Errorf("injected lost response after replace")
		}
		return nil
	}
	interrupted, err := InstallQualifiedCLI(d, f.R.WorkspaceRoot, o.RunID, dest, "test coordinator", "install")
	if err == nil || interrupted.Receipt == nil || interrupted.State != "effect_unknown" {
		t.Fatalf("missing effect receipt: %+v %v", interrupted, err)
	}
	inode, err := os.Stat(dest)
	if err != nil {
		t.Fatal(err)
	}
	// The actual installed process fails to execute. Preserve that known failure,
	// then prove repeated installation cannot repair permissions and hide it.
	d.Fault = func(stage string) error {
		if stage == "cli_install_before_smoke" {
			return os.Chmod(dest, 0600)
		}
		return nil
	}
	failed, err := InstallQualifiedCLI(d, f.R.WorkspaceRoot, o.RunID, dest, "test coordinator", "install")
	if err == nil || failed.State != "smoke_failed" || failed.Smoke == nil || failed.Smoke.Error == "" {
		t.Fatalf("false smoke pass %+v %v", failed, err)
	}
	t.Logf("actual process failure exit=%v error=%s stderr=%s", failed.Smoke.Exit, failed.Smoke.Error, failed.Smoke.Stderr)
	d.Fault = nil
	observed, observeErr := InstallQualifiedCLI(d, f.R.WorkspaceRoot, o.RunID, dest, "test coordinator", "observe")
	if observeErr == nil || observed.ObservedSHA256 != ready.Source.BinarySHA256 || observed.Smoke == nil || observed.Smoke.Error != failed.Smoke.Error {
		t.Fatalf("observe hid failed smoke or destination: %+v %v", observed, observeErr)
	}
	if err := os.WriteFile(dest, []byte("unrelated changed result"), 0600); err != nil {
		t.Fatal(err)
	}
	unknown, unknownErr := InstallQualifiedCLI(d, f.R.WorkspaceRoot, o.RunID, dest, "test coordinator", "observe")
	if unknownErr == nil || unknown.State != "effect_unknown" || unknown.ObservedSHA256 == ready.Source.BinarySHA256 || unknown.Smoke == nil {
		t.Fatalf("observe hid changed destination behind smoke failure: %+v %v", unknown, unknownErr)
	}
	retained, err := os.ReadFile(ready.Controller.Locator)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dest, retained, 0600); err != nil {
		t.Fatal(err)
	}
	repeated, err := InstallQualifiedCLI(d, f.R.WorkspaceRoot, o.RunID, dest, "test coordinator", "install")
	after, statErr := os.Stat(dest)
	if err == nil || repeated.State != "smoke_failed" || statErr != nil || !os.SameFile(inode, after) || after.Mode().Perm() != 0600 {
		t.Fatalf("known failure was hidden by retry: %+v %v", repeated, err)
	}
	restored, err := InstallQualifiedCLI(d, f.R.WorkspaceRoot, o.RunID, dest, "test coordinator", "rollback")
	if err != nil || restored.State != "rolled_back" {
		t.Fatalf("guarded failure rollback %+v %v", restored, err)
	}
}

func TestQualifiedCLIInstalledSmokePreservesNonzeroExit(t *testing.T) {
	target := filepath.Join(t.TempDir(), "failed-cli")
	if err := os.WriteFile(target, []byte("#!/bin/sh\necho actual-smoke-failure >&2\nexit 7\n"), 0700); err != nil {
		t.Fatal(err)
	}
	smoke := probeInstalledCLI(target, "unused")
	if smoke.Exit == nil || *smoke.Exit != 7 || !strings.Contains(smoke.Stderr, "actual-smoke-failure") || smoke.Error == "" {
		t.Fatalf("lost actual process failure: %+v", smoke)
	}
	t.Logf("actual smoke exit=%d stderr=%s error=%s", *smoke.Exit, smoke.Stderr, smoke.Error)
}
