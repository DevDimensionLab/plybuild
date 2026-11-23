package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/machine"
)

func runMachine(t *testing.T, read func() (machine.Status, error), args ...string) (string, string, error) {
	t.Helper()
	return runMachineWith(t, read, func() ([]machine.Process, error) { return nil, errors.New("not called") }, args...)
}

func runMachineWith(t *testing.T, read func() (machine.Status, error), processes func() ([]machine.Process, error), args ...string) (string, string, error) {
	t.Helper()
	c := newMachineCommand(read, processes)
	c.SilenceUsage, c.SilenceErrors = true, true // as ExecuteE sets on the root
	var out, errOut bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&errOut)
	c.SetArgs(append([]string{"status"}, args...))
	err := c.Execute()
	return out.String(), errOut.String(), err
}

func TestMachineStatusFormats(t *testing.T) {
	pressure := "normal"
	s := machine.Status{Kind: machine.StatusKind, SchemaVersion: 1, Cores: 16,
		Load: machine.Windows{M1: 4.91, M5: 9.34, M15: 9.97}, LoadPercent: machine.Windows{M1: 30.7, M5: 58.4, M15: 62.3}, MemoryPressure: &pressure}
	read := func() (machine.Status, error) { return s, nil }
	out, _, err := runMachine(t, read, "--format", "json")
	if err != nil || !strings.Contains(out, `"kind":"MachineStatusReadback@1"`) {
		t.Fatalf("json: %v %q", err, out)
	}
	out, _, err = runMachine(t, read)
	if err != nil || out != "Load 31% (1m) 58% (5m) 62% (15m) of 16 cores; memory unknown, swap unknown; memory pressure normal\n" {
		t.Fatalf("text: %v %q", err, out)
	}
	if strings.Contains(out, "processes") {
		t.Fatalf("default call shows processes: %q", out)
	}
	out, _, err = runMachineWith(t, read, func() ([]machine.Process, error) {
		return []machine.Process{{PID: 7, Name: "node", CPUPercent: 99.5, RSSBytes: 1 << 20}}, nil
	}, "--processes", "--format", "json")
	if err != nil || !strings.Contains(out, `"processes":[{"pid":7,"name":"node","cpu_percent":99.5,"rss_bytes":1048576}],"processes_reason":null`) {
		t.Fatalf("processes json: %v %q", err, out)
	}
	out, _, err = runMachineWith(t, read, func() ([]machine.Process, error) { return nil, errors.New("boom") }, "--processes")
	if err != nil || !strings.HasSuffix(out, "\n  processes unknown (read_failed)\n") {
		t.Fatalf("processes failure: %v %q", err, out)
	}
	if _, _, err = runMachine(t, read, "--format", "yaml"); err == nil {
		t.Fatal("accepted an unknown format")
	}
}

func TestMachineStatusFailureWritesNothingToStdout(t *testing.T) {
	read := func() (machine.Status, error) { return machine.Status{}, errors.New("read_failed: boom") }
	out, _, err := runMachine(t, read, "--format", "json")
	if err == nil || out != "" || !strings.HasPrefix(err.Error(), "read_failed") {
		t.Fatalf("got %v %q", err, out)
	}
}

// Acceptance 5 and 6: the real command reads macOS without a workspace or profile.
func TestMachineStatusReadsThisMachineWithoutWorkspace(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("machine status reads macOS only")
	}
	t.Setenv("HOME", t.TempDir())
	t.Chdir(t.TempDir())
	var out bytes.Buffer
	RootCmd.SetOut(&out)
	RootCmd.SetArgs([]string{"machine", "status", "--format", "json"})
	t.Cleanup(func() { RootCmd.SetOut(nil); RootCmd.SetArgs(nil) })
	if err := ExecuteE(); err != nil {
		t.Fatal(err)
	}
	var s machine.Status
	if err := json.Unmarshal(out.Bytes(), &s); err != nil {
		t.Fatalf("%v: %s", err, out.String())
	}
	raw, err := exec.Command("sysctl", "-n", "hw.logicalcpu").Output()
	if err != nil {
		t.Fatal(err)
	}
	cores, _ := strconv.Atoi(strings.TrimSpace(string(raw)))
	if s.Cores != cores || s.Load.M1 < 0 || s.Load.M5 < 0 || s.Load.M15 < 0 || s.MemoryPressure == nil {
		t.Fatalf("got %+v, want %d cores", s, cores)
	}
	raw, err = exec.Command("sysctl", "-n", "hw.memsize").Output()
	if err != nil {
		t.Fatal(err)
	}
	memsize, _ := strconv.ParseUint(strings.TrimSpace(string(raw)), 10, 64)
	if s.MemoryTotalBytes != memsize || s.MemoryUsedPercent == nil || *s.MemoryUsedPercent <= 0 || *s.MemoryUsedPercent > 100 || s.SwapUsedBytes == nil || s.SwapTotalBytes == nil || *s.SwapUsedBytes > *s.SwapTotalBytes {
		t.Fatalf("memory or swap: %+v, want %d bytes", s, memsize)
	}
}
