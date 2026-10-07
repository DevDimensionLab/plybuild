package taskexecute

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

// This stand-in only runs in launcherFixture's disposable workspace. It models
// a shell left by a failed pre-Task startup and the single replacement owner.
const recoveryLauncherHerdr = `#!/usr/bin/env python3
import json, sys
from pathlib import Path
root = Path(__file__).parent
path = root / "recovery-model.json"
model = json.loads(path.read_text()) if path.exists() else {}
args = sys.argv[1:]
cmd = args[:2]
pane = {"workspace_id":model.get("workspace", "w-fixture"), "tab_id":model.get("tab", "tab-fixture"), "pane_id":model.get("pane", "w-fixture:p1"), "terminal_id":model.get("terminal", "terminal-fixture"), "agent_status":"unknown", "cwd":model.get("cwd", ""), "foreground_cwd":model.get("cwd", "")}
if cmd == ["tab", "create"]:
    model["tabs"] = model.get("tabs", 0) + 1
    if model.get("closed"):
        model.update(closed=False, workspace=args[args.index("--workspace") + 1], tab="replacement-tab", pane="replacement-pane", terminal="replacement-terminal")
        pane.update(workspace_id=model["workspace"], tab_id=model["tab"], pane_id=model["pane"], terminal_id=model["terminal"])
    model["cwd"] = args[args.index("--cwd") + 1]
    pane.update(cwd=model["cwd"], foreground_cwd=model["cwd"])
    result = {"root_pane":pane}
elif cmd == ["pane", "get"] and model.get("closed"):
    print(json.dumps({"error":{"code":"pane_not_found","message":"synthetic closed Task tab"}}), file=sys.stderr)
    raise SystemExit(1)
elif cmd == ["pane", "list"]:
    result = {"panes": [] if model.get("closed") else [pane]}
elif cmd == ["workspace", "get"]:
    result = {"workspace":{"workspace_id":args[2]}}
elif cmd == ["pane", "get"]:
    result = {"pane":pane}
elif cmd == ["pane", "process-info"]:
    result = {"process_info":{"pane_id":"w-fixture:p1", "shell_pid":531, "foreground_process_group_id":531, "foreground_processes":[{"pid":531,"name":"zsh","argv0":"zsh","argv":["-zsh"],"cwd":model["cwd"]}]}}
elif cmd == ["agent", "get"] and not model.get("starts"):
    print(json.dumps({"error":{"code":"agent_not_found","message":"synthetic exited provider"}}), file=sys.stderr)
    raise SystemExit(1)
elif cmd in (["agent", "start"], ["agent", "get"], ["agent", "prompt"]):
    if cmd == ["agent", "start"]:
        model["name"] = args[2]
        model["starts"] = model.get("starts",0) + 1
    if cmd == ["agent", "prompt"]:
        model["prompts"] = model.get("prompts",0) + 1
    result = {"agent":{**pane,"name":model["name"],"agent":"codex","agent_status":"working" if model.get("prompts") else "idle","agent_session":{"agent":"codex","kind":"id","value":"replacement-fixture-session"}}}
else:
    raise SystemExit("unexpected synthetic command: " + repr(args))
path.write_text(json.dumps(model))
print(json.dumps({"result":result}))
`

func recoveryFiles(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			out[path] = digestBytes(b)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestRecoveryLauncherPreservesControlAndReusesOneAttempt(t *testing.T) {
	d, in, root, _ := launcherFixtureForProvider(t, "codex")
	d.StartupProcessTree = func(shellPID int) ([]taskrun.StartupProcess, error) {
		return []taskrun.StartupProcess{{PID: shellPID, ParentPID: 1, ProcessGroupID: shellPID}}, nil
	}
	originalProvider := filepath.Join(root, "bin", "codex-original")
	if err := os.Rename(in.Runtime.ProviderExecutable, originalProvider); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(originalProvider, in.Runtime.ProviderExecutable); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(in.Runtime.HerdrExecutable, []byte(recoveryLauncherHerdr), 0700); err != nil {
		t.Fatal(err)
	}
	d.Fault = func(point string) error {
		if point == "workflow_before_agent_send" {
			return errors.New("synthetic startup interrupted before Task input")
		}
		return nil
	}
	started, err := Execute(d, in)
	if err == nil || started.Run == nil || started.Run.RunID == "" {
		t.Fatalf("fixture did not preserve startup: %+v %v", started, err)
	}
	d.Fault = nil
	originalRequest, err := os.ReadFile(started.RequestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(originalProvider); err != nil {
		t.Fatal(err)
	}
	newProvider := filepath.Join(root, "bin", "codex-updated")
	if err = os.WriteFile(newProvider, []byte("#!/bin/sh\n# Updated isolated provider\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(in.Runtime.ProviderExecutable); err != nil {
		t.Fatal(err)
	}
	if err = os.Symlink(newProvider, in.Runtime.ProviderExecutable); err != nil {
		t.Fatal(err)
	}
	recovery := RecoveryInput{RunID: started.Run.RunID, Check: true, Timeout: 3 * time.Second, ProviderExecutable: newProvider}
	before := recoveryFiles(t, started.Run.Paths.RunRoot)
	registry, err := d.Workspace.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	checked, err := RecoverStartup(d, recovery)
	if err != nil || checked.State != "ready" || checked.Preview == nil {
		t.Fatalf("recovery check failed: %+v %v", checked, err)
	}
	if after := recoveryFiles(t, started.Run.Paths.RunRoot); !reflect.DeepEqual(before, after) {
		t.Fatal("recovery check changed the preserved run")
	}
	if _, err = os.Lstat(checked.Preview.ControlExecutable.Path); !os.IsNotExist(err) {
		t.Fatal("recovery check preserved a control binary")
	}
	recovery.Check = false
	applied, err := RecoverStartup(d, recovery)
	if err != nil || applied.State != "started" || applied.Run == nil || applied.Run.RunID != started.Run.RunID || applied.Run.StartupRecovery == nil {
		t.Fatalf("replacement startup failed: %+v %v", applied, err)
	}
	control, err := bindExecutable(checked.Preview.ControlExecutable.Path)
	if err != nil || control != checked.Preview.ControlExecutable {
		t.Fatalf("control was not preserved: %+v %v", control, err)
	}
	request, err := os.ReadFile(started.RequestPath)
	if err != nil || string(request) != string(originalRequest) {
		t.Fatal("recovery changed the original goal request")
	}
	if after, err := d.Workspace.WorkItems.Snapshot(root); err != nil || after.RawSHA256 != registry.RawSHA256 {
		t.Fatal("recovery changed Task preparation or queue state")
	}
	// A later removed provider must not make retry select another owner or
	// require a newly working PATH merely to inspect the already reserved start.
	if err = os.Remove(newProvider); err != nil {
		t.Fatal(err)
	}
	again, err := RecoverStartup(d, recovery)
	if err != nil || again.State != "existing" || again.Run.RunID != applied.Run.RunID {
		t.Fatalf("repeat did not inspect the same replacement: %+v %v", again, err)
	}
	var model struct{ Starts, Prompts int }
	raw, err := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json"))
	if err != nil || json.Unmarshal(raw, &model) != nil || model.Starts != 1 || model.Prompts != 1 {
		t.Fatalf("replacement effects = %s (%v)", strings.TrimSpace(string(raw)), err)
	}
}
