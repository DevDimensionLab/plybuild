package weekrun

import (
	"encoding/xml"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestWeekPlistBindsFrozenControlAndLoginEnvironment(t *testing.T) {
	f := newFixture(t)
	f.b.Environment = map[string]string{"HOME": f.c.Workspace, "PATH": "/custom/codex:/usr/bin:/bin", "CODEX_HOME": "/synthetic/account<&>"}
	data := Plist(f.b)
	d := xml.NewDecoder(strings.NewReader(string(data)))
	for {
		_, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal("invalid XML", err)
		}
	}
	for _, want := range []string{f.b.Control.Path, "_supervise", f.b.ID, "<key>RunAtLoad</key><true/>", "<key>SuccessfulExit</key><false/>", "/synthetic/account&lt;&amp;&gt;"} {
		if !strings.Contains(string(data), want) {
			t.Fatalf("plist misses %q", want)
		}
	}
	if strings.Contains(string(data), "--model") || strings.Contains(string(data), "API_KEY") {
		t.Fatal("service invented provider options")
	}
	if runtime.GOOS == "darwin" {
		path := filepath.Join(f.c.Workspace, "service.plist")
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		if out, err := exec.Command("/usr/bin/plutil", "-lint", path).CombinedOutput(); err != nil {
			t.Fatalf("actual macOS plist parser: %s: %v", out, err)
		}
	}
}
func TestWeekServiceLifecycleUsesOnlyItsNamedFrozenPlist(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("launchd is macOS only")
	}
	f := newFixture(t, selection("first"))
	f.tick(t)
	f.tick(t)
	// Redirect only the synthetic binding's service destination; no real user
	// LaunchAgents, launchctl service, login files or provider process is used.
	f.b.Environment = map[string]string{"HOME": f.c.Workspace, "PATH": "/usr/bin:/bin"}
	if err := publish(filepath.Join(f.dir, "binding.json"), f.b, true); err != nil {
		t.Fatal(err)
	}
	loaded := false
	calls := []string{}
	f.m.Launchctl = func(args ...string) (string, error) {
		calls = append(calls, args[0])
		switch args[0] {
		case "print":
			if !loaded {
				return "", errors.New("not loaded")
			}
			return "state = running\n pid = 731\n last exit code = 0\n", nil
		case "bootstrap":
			if args[2] != servicePath(f.b) {
				t.Fatal(args)
			}
			loaded = true
		case "bootout":
			if !reflect.DeepEqual(f.n.controls, []string{"stop"}) {
				t.Fatal("unloaded before native worker stop was acknowledged", f.n.controls)
			}
			loaded = false
		case "kickstart":
			for _, a := range args {
				if a == "-k" {
					t.Fatal("unsafe replacement")
				}
			}
		default:
			t.Fatal(args)
		}
		return "", nil
	}
	f.m.Launch = func(b Binding) error { f.tick(t); return nil }
	st, err := f.m.Service(f.c.Workspace, f.s.ID, "install")
	if err != nil || st.State != "running" || st.PID != 731 {
		t.Fatal(st, err)
	}
	if st.Label != "com.ply.week."+f.s.ID || st.ControlSHA256 != f.b.Control.SHA256 {
		t.Fatal(st)
	}
	st, err = f.m.Service(f.c.Workspace, f.s.ID, "stop")
	if err != nil || st.State != "installed" || loaded {
		t.Fatal(st, err)
	}
	if f.s.Phase != "stopping" || f.s.Entries[0].RunID != "run-first" {
		t.Fatal("stop did not preserve native Run", f.s)
	}
	st, err = f.m.Service(f.c.Workspace, f.s.ID, "start")
	if err != nil || st.State != "running" {
		t.Fatal(st, err)
	}
	// Stop repeats idempotently while the same native request is already visible.
	st, err = f.m.Service(f.c.Workspace, f.s.ID, "uninstall")
	if err != nil || st.State != "absent" || loaded {
		t.Fatal(st, err)
	}
	if _, err = os.Lstat(servicePath(f.b)); !os.IsNotExist(err) {
		t.Fatal("managed plist remained")
	}
	if _, err = os.Lstat(filepath.Join(f.dir, "binding.json")); err != nil {
		t.Fatal("uninstall removed native/operational evidence", err)
	}
	if len(calls) == 0 {
		t.Fatal("no service effects")
	}
}
func TestWeekServiceRefusesDifferentPlist(t *testing.T) {
	f := newFixture(t)
	f.b.Environment = map[string]string{"HOME": f.c.Workspace}
	path := servicePath(f.b)
	if err := directory(filepath.Dir(path)); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("unrelated service"), 0600); err != nil {
		t.Fatal(err)
	}
	calls := 0
	_, err := ServiceReadback(f.b, func(...string) (string, error) { calls++; return "", nil })
	if err == nil || calls != 0 {
		t.Fatal("mismatched plist was accepted or observed through launchctl", err, calls)
	}
}
