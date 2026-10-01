package taskrun

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestR2TTYStopBeforeWrite(t *testing.T) {
	d, r, f := fixture(t)
	d.Runner.(*fakeRunner).checkErr = failure("task_run_tty_required", 3, "synthetic no terminal")
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = Start(d, f, *p.(Preview).Confirmation); e == nil {
		t.Fatal("no TTY accepted")
	}
	if _, e = os.Stat(storeRoot(r.WorkspaceRoot)); !os.IsNotExist(e) {
		t.Fatal("TTY preflight wrote files")
	}
}
func TestR4ReservationCrashesNeverRelaunch(t *testing.T) {
	for _, point := range []string{"after_request_reservation", "after_request", "after_reservation", "after_launch_intent"} {
		t.Run(point, func(t *testing.T) {
			d, r, f := fixture(t)
			d.Fault = func(p string) error {
				if p == point {
					return errors.New("synthetic crash")
				}
				return nil
			}
			p, e := PreviewStart(d, f)
			if e != nil {
				t.Fatal(e)
			}
			if _, e = Start(d, f, *p.(Preview).Confirmation); e == nil {
				t.Fatal("fault did not fire")
			}
			d.Fault = nil
			out, e := Start(d, f, "retry")
			if e != nil {
				t.Fatal(e)
			}
			if out.Launch.State != "unknown" || d.Runner.(*fakeRunner).calls != 0 {
				t.Fatalf("partial reservation relaunched: %+v", out)
			}
			other := r
			other.RequestKey = "another/run"
			file := writeAny(t, r.WorkspaceRoot, "other.json", other)
			if _, e = PreviewStart(d, file); e == nil {
				t.Fatal("another run bypassed unknown target reservation")
			}
		})
	}
}
func TestR4KnownSpawnFailure(t *testing.T) {
	d, r, f := fixture(t)
	runner := d.Runner.(*fakeRunner)
	runner.process = Process{State: "not_started", Quiescence: ptr(true)}
	runner.errorAfter = errors.New("synthetic spawn failure")
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	out, e := Start(d, f, *p.(Preview).Confirmation)
	var code *Error
	if !errors.As(e, &code) || code.Exit != 4 || out.Launch.State != "failed" || out.Launch.Attempts != 1 || out.Acceptance.State != "missing" {
		t.Fatalf("failure facts: %v %+v", e, out)
	}
	if _, e = Start(d, f, "retry"); e != nil {
		t.Fatal(e)
	}
	if runner.calls != 1 {
		t.Fatal("spawn retried")
	}
	_ = r
}
func TestR5CallbackBindingsAndImmutableAcceptance(t *testing.T) {
	d, r, f := fixture(t)
	runner := d.Runner.(*fakeRunner)
	runner.inside = func(s LaunchSpec) error {
		good := acceptance(t, d, r)
		for _, field := range []string{"runtime", "model", "policy", "session", "cwd", "scope"} {
			bad := good
			old := d.CWD
			switch field {
			case "runtime":
				bad.RuntimeClaim.RuntimeID = ptr("other")
			case "model":
				bad.RuntimeClaim.ModelID = ptr("other")
			case "policy":
				bad.RuntimeClaim.EffectivePolicySHA256 = ptr(hash([]byte("other")))
			case "session":
				bad.SessionID = "other"
			case "cwd":
				d.CWD = func() (string, error) { return r.WorkspaceRoot, nil }
			case "scope":
				bad.Sandbox = json.RawMessage(`{"read_roots":[],"write_roots":[],"temp_root":"/","matches_contract":true}`)
			}
			_, e := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "bad-"+field+".json", bad))
			d.CWD = old
			if e == nil {
				return fmt.Errorf("accepted %s drift", field)
			}
		}
		cp := writeAny(t, r.WorkspaceRoot, "good-claim.json", good)
		if _, e := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), cp); e != nil {
			return e
		}
		if _, e := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), cp); e != nil {
			return e
		}
		good.RuntimeClaim.NativeSessionID = ptr("changed-native-session")
		if _, e := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "changed-claim.json", good)); e == nil {
			return fmt.Errorf("accepted changed immutable claim")
		}
		return nil
	}
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	out, _ := Start(d, f, *p.(Preview).Confirmation)
	if runner.lastErr != nil {
		t.Fatal(runner.lastErr)
	}
	if out.Acceptance.State != "started" || out.Collection.State == "qualified" {
		t.Fatalf("wrong facts: %+v", out)
	}
}
func completeInside(t *testing.T, d *Dependencies, r Request, change func(*Report)) {
	t.Helper()
	d.Runner.(*fakeRunner).inside = func(s LaunchSpec) error {
		a := acceptance(t, *d, r)
		if _, e := Accept(*d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "claim.json", a)); e != nil {
			return e
		}
		report := semanticReport(t, *d, r)
		if change != nil {
			change(&report)
		}
		_, e := SubmitReport(*d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "return.json", report))
		return e
	}
}
func TestR7ProcessEvidenceCannotBeInferredFromReport(t *testing.T) {
	for _, state := range []string{"missing-report", "live-group", "signal", "identity-unknown"} {
		t.Run(state, func(t *testing.T) {
			d, r, f := fixture(t)
			runner := d.Runner.(*fakeRunner)
			if state != "missing-report" {
				completeInside(t, &d, r, nil)
				p := Process{State: "exited", PID: ptr(731), ProcessGroup: ptr(731), StartIdentity: ptr("synthetic-boot:birth-123"), ExitCode: ptr(0), Quiescence: ptr(true)}
				switch state {
				case "live-group":
					p.Quiescence = ptr(false)
				case "signal":
					p.Signal = ptr("terminated")
					p.ExitCode = nil
				case "identity-unknown":
					p.State = "unknown"
					p.Quiescence = nil
				}
				runner.process = p
			}
			preview, e := PreviewStart(d, f)
			if e != nil {
				t.Fatal(e)
			}
			out, e := Start(d, f, *preview.(Preview).Confirmation)
			if runner.lastErr != nil {
				t.Fatal(runner.lastErr)
			}
			if e == nil || out.Collection.State == "qualified" {
				t.Fatalf("false green: %v %+v", e, out)
			}
		})
	}
}
func TestR10BudgetBoundaryAndUnknown(t *testing.T) {
	for _, change := range []string{"rounds", "seconds", "environment", "unknown"} {
		t.Run(change, func(t *testing.T) {
			d, r, f := fixture(t)
			completeInside(t, &d, r, func(report *Report) {
				switch change {
				case "rounds":
					report.BudgetUsage.CorrectionRounds = 4
				case "seconds":
					report.BudgetUsage.ActiveSeconds = 5401
				case "environment":
					report.BudgetUsage.EnvironmentMeasures = 3
				case "unknown":
					report.BudgetUsage = nil
				}
			})
			p, e := PreviewStart(d, f)
			if e != nil {
				t.Fatal(e)
			}
			out, e := Start(d, f, *p.(Preview).Confirmation)
			if e == nil || out.Collection.State == "qualified" {
				t.Fatalf("budget falsely qualified: %+v", out)
			}
			if d.Runner.(*fakeRunner).lastErr != nil {
				t.Fatal(d.Runner.(*fakeRunner).lastErr)
			}
		})
	}
}
func TestR9StrictDocumentsAndSymlinks(t *testing.T) {
	d, r, f := fixture(t)
	raw, _ := os.ReadFile(f)
	for _, b := range [][]byte{append(raw, []byte(" {}")...), []byte(strings.Replace(string(raw), `"schema_version":1`, `"schema_version":1,"schema_version":1`, 1)), []byte(strings.Replace(string(raw), `"schema_version":1`, `"schema_version":1.0`, 1)), []byte(strings.Replace(string(raw), `"schema_version":1`, `"schema_version":1,"unknown":true`, 1)), []byte(strings.Replace(string(raw), `"max_active_seconds":5400`, `"max_active_seconds":null`, 1))} {
		if _, e := parseRequest(b); e == nil {
			t.Fatal("invalid JSON accepted")
		}
	}
	link := filepath.Join(r.WorkspaceRoot, "request-link.json")
	os.Symlink(f, link)
	if _, e := ReadRequest(link); e == nil {
		t.Fatal("symlink accepted")
	}
	b := r
	b.Runtime.Executable.SHA256 = hash([]byte("drift"))
	if _, e := PreviewStart(d, writeAny(t, r.WorkspaceRoot, "drift.json", b)); e == nil {
		t.Fatal("executable drift accepted")
	}
}
func TestR9BlockedReportDoesNotFabricateVerifier(t *testing.T) {
	d, r, f := fixture(t)
	completeInside(t, &d, r, func(p *Report) {
		p.Outcome = "blocked"
		p.StopReasons = json.RawMessage(`[{"type":"verifier_unrecoverable","detail":"Synthetic verifier was not run.","primary":true}]`)
		p.VerifierResults = json.RawMessage(`[]`)
		p.TechnicalAssessment.Gate = "unknown"
	})
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	out, _ := Start(d, f, *p.(Preview).Confirmation)
	if d.Runner.(*fakeRunner).lastErr != nil {
		t.Fatal(d.Runner.(*fakeRunner).lastErr)
	}
	if out.Delivery.State != "not_representable" || out.Collection.State == "qualified" {
		t.Fatalf("blocked report: %+v", out)
	}
}
func TestR3CompetingProcesses(t *testing.T) {
	d, r, f := fixture(t)
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	binary, e := os.Executable()
	if e != nil {
		t.Fatal(e)
	}
	count := filepath.Join(r.WorkspaceRoot, "synthetic-spawn-count")
	children := []*exec.Cmd{}
	for i := 0; i < 2; i++ {
		c := exec.Command(binary, "-test.run=^TestReservationChildHelper$")
		c.Dir = r.WorkspaceRoot
		c.Env = append(os.Environ(), "PLY_SYNTHETIC_RESERVATION_CHILD=1", "PLY_SYNTHETIC_REQUEST="+f, "PLY_SYNTHETIC_CONFIRM="+*p.(Preview).Confirmation, "PLY_SYNTHETIC_COUNT="+count)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if e = c.Start(); e != nil {
			t.Fatal(e)
		}
		children = append(children, c)
	}
	for _, c := range children {
		if e = c.Wait(); e != nil {
			t.Fatal(e)
		}
	}
	b, e := os.ReadFile(count)
	if e != nil || string(b) != "spawn\n" {
		t.Fatalf("spawn count %q: %v", b, e)
	}
}
func TestReservationChildHelper(t *testing.T) {
	if os.Getenv("PLY_SYNTHETIC_RESERVATION_CHILD") != "1" {
		return
	}
	d := SystemDependencies(workspace.SystemDependencies())
	d.Runner = &fakeRunner{inside: func(LaunchSpec) error {
		f, e := os.OpenFile(os.Getenv("PLY_SYNTHETIC_COUNT"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if e != nil {
			return e
		}
		defer f.Close()
		_, e = f.WriteString("spawn\n")
		return e
	}}
	_, e := Start(d, os.Getenv("PLY_SYNTHETIC_REQUEST"), os.Getenv("PLY_SYNTHETIC_CONFIRM"))
	if e != nil {
		var err *Error
		if !errors.As(e, &err) || err.Exit != 5 && err.Exit != 3 {
			t.Fatal(e)
		}
	}
}

type failingPublication struct{ systemFileSystem }

func (failingPublication) WriteOnce(string, []byte) error {
	return errors.New("synthetic disk failure")
}
func TestR4FilesystemFailureCannotSpawn(t *testing.T) {
	d, r, f := fixture(t)
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	d.Files = failingPublication{}
	if _, e = Start(d, f, *p.(Preview).Confirmation); e == nil {
		t.Fatal("disk failure ignored")
	}
	if d.Runner.(*fakeRunner).calls != 0 {
		t.Fatal("spawned without durable intent")
	}
	_ = r
}

func TestR9ImpossibleProcessFactsRejected(t *testing.T) {
	for _, p := range []Process{{State: "exited", Quiescence: ptr(true)}, {State: "running", PID: ptr(12), ProcessGroup: ptr(12)}, {State: "invented"}, {State: "exited", PID: ptr(-1), ExitCode: ptr(0)}} {
		if validateProcess(p) == nil {
			t.Fatalf("accepted impossible process: %+v", p)
		}
	}
}
func TestR1PolicyAndTargetDriftBlockBeforeReservation(t *testing.T) {
	for _, kind := range []string{"policy", "epic", "target"} {
		t.Run(kind, func(t *testing.T) {
			d, r, f := fixture(t)
			p, e := PreviewStart(d, f)
			if e != nil {
				t.Fatal(e)
			}
			plan := p.(Preview).Preparation.Plan
			switch kind {
			case "policy":
				os.WriteFile(r.Runtime.PermissionBinding.Evidence[0].Locator, []byte("changed synthetic policy"), 0600)
			case "epic":
				os.WriteFile(filepath.Join(plan.Target.ParentLocator, "drift.txt"), []byte("drift"), 0644)
			case "target":
				os.WriteFile(filepath.Join(plan.WorktreePath, "drift.txt"), []byte("drift"), 0644)
			}
			if _, e = Start(d, f, *p.(Preview).Confirmation); e == nil {
				t.Fatal("drift accepted")
			}
			if _, e = os.Stat(storeRoot(r.WorkspaceRoot)); !os.IsNotExist(e) {
				t.Fatal("drift created reservation")
			}
		})
	}
}
