package taskrun

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type testTerminal struct {
	calls        []string
	restoreError error
}

func (t *testTerminal) Check() error { t.calls = append(t.calls, "check"); return nil }
func (t *testTerminal) Capture() (func() error, error) {
	t.calls = append(t.calls, "capture")
	return func() error { t.calls = append(t.calls, "restore"); return t.restoreError }, nil
}
func (t *testTerminal) Foreground(int) error { t.calls = append(t.calls, "foreground"); return nil }

type syntheticChild struct {
	cmd *exec.Cmd
	p   Process
}

func (c *syntheticChild) Process() Process { return c.p }
func (c *syntheticChild) Forward(os.Signal) error {
	return errors.New("synthetic child has no authorized signal group")
}
func (c *syntheticChild) Wait() (Process, error) {
	e := c.cmd.Wait()
	p := c.p
	p.State = "exited"
	p.ExitCode = ptr(c.cmd.ProcessState.ExitCode())
	p.Quiescence = ptr(true)
	return p, e
}
func TestR11InjectedTerminalRestorationWithSyntheticProcess(t *testing.T) {
	binary, e := os.Executable()
	if e != nil {
		t.Fatal(e)
	}
	binary, e = filepath.EvalSymlinks(binary)
	if e != nil {
		t.Fatal(e)
	}
	for _, restoreFails := range []bool{false, true} {
		terminal := &testTerminal{}
		if restoreFails {
			terminal.restoreError = errors.New("synthetic restore failure")
		}
		r := &InteractiveRunner{Terminal: terminal, Spawn: func(s LaunchSpec) (Child, error) {
			cmd := exec.Command(binary, "-test.run=^TestSyntheticProcessHelper$")
			cmd.Env = append(os.Environ(), "PLY_SYNTHETIC_PROCESS=1")
			if e := cmd.Start(); e != nil {
				return nil, e
			}
			pid := cmd.Process.Pid
			return &syntheticChild{cmd, Process{State: "running", PID: &pid, ProcessGroup: &pid, StartIdentity: ptr("synthetic-birth-only"), Quiescence: ptr(false)}}, nil
		}}
		called := false
		p, e := r.Run(LaunchSpec{Executable: Executable{binary, hashFileTest(t, binary)}}, func(Process) error { called = true; return nil })
		if !called || terminal.calls[len(terminal.calls)-1] != "restore" {
			t.Fatal("terminal lifecycle incomplete")
		}
		if restoreFails {
			if e == nil || p.State != "unknown" || p.Quiescence != nil {
				t.Fatalf("restoration failure became green: %+v", p)
			}
		} else if p.State != "exited" || p.ExitCode == nil || *p.ExitCode != 17 {
			t.Fatalf("synthetic exit lost: %+v %v", p, e)
		}
	}
}
func TestSyntheticProcessHelper(t *testing.T) {
	if os.Getenv("PLY_SYNTHETIC_PROCESS") == "1" {
		os.Exit(17)
	}
}
