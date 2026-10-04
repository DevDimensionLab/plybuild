package workflownotification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

const syntheticSecret = "https://hooks.slack.com/services/TSYNTHETIC/BSYNTHETIC/fixtureOnly"

type fixture struct {
	root    string
	request Request
	route   Route
	input   Input
	deps    Dependencies
	posts   int
}

func writeJSON(t *testing.T, path string, v any) {
	t.Helper()
	b, e := json.Marshal(v)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(path, b, 0600); e != nil {
		t.Fatal(e)
	}
}
func newFixture(t *testing.T) *fixture {
	t.Helper()
	root, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	f := &fixture{root: root}
	work := filepath.Join(root, "worktree")
	if e = os.Mkdir(work, 0700); e != nil {
		t.Fatal(e)
	}
	startKind, reportKind := "FixtureStart@1", "FixtureReport@1"
	f.request = Request{Kind: "ply.workflow.notification-request", SchemaVersion: 1, Route: "fixture", Source: Source{Kind: "external", Activity: "notification/tests", Run: "run-1", Worktree: work}, Gate: &Gate{ID: "result-control", Revision: 1, OpenedAt: "2026-10-03T19:00:00Z", Reason: "result_control", State: "waiting_for_human"}, Public: Public{TaskTitle: "Test delivery", NextAction: "Start the planning launcher for result control."}}
	f.request.Sender.ActorClaim = "fixture-only"
	for name, kind := range map[string]*string{"handoff": nil, "start": &startKind, "report": &reportKind} {
		p := filepath.Join(root, name)
		var b []byte
		if kind == nil {
			b = []byte("Fixture handoff\nPRIVATE_SOURCE_SENTINEL\n")
		} else {
			b, _ = json.Marshal(map[string]any{"kind": *kind, "activity": f.request.Source.Activity, "run": f.request.Source.Run, "private": "PRIVATE_SOURCE_SENTINEL"})
		}
		if e = os.WriteFile(p, b, 0600); e != nil {
			t.Fatal(e)
		}
		loc := Locator{p, Digest(b), kind}
		switch name {
		case "handoff":
			f.request.Source.Handoff = loc
		case "start":
			f.request.Source.Start = &loc
		case "report":
			f.request.Source.Report = &loc
		}
	}
	f.route = Route{"ply.workflow.notification-route", 1, "fixture", "slack_incoming_webhook", "#fixture", "FIXTURE_WEBHOOK", filepath.Join(root, "state")}
	f.input = Input{Operation: "send", File: filepath.Join(root, "request.json"), Route: filepath.Join(root, "route.json")}
	f.save(t)
	f.deps = Dependencies{Now: func() time.Time { return time.Date(2026, 10, 3, 19, 0, 0, 0, time.UTC) }, LookupEnv: func(string) (string, bool) { return syntheticSecret, true }, Transport: func(context.Context, string, []byte) Observation {
		f.posts++
		return Observation{Status: 200, Body: []byte("ok"), Dispatch: "started"}
	}}
	return f
}
func (f *fixture) save(t *testing.T) {
	writeJSON(t, f.input.File, f.request)
	writeJSON(t, f.input.Route, f.route)
}
func (f *fixture) preview(t *testing.T) Preview {
	t.Helper()
	v, e := Execute(f.deps, f.input)
	if e != nil {
		t.Fatalf("preview: %v (%+v)", e, v)
	}
	return v.(Preview)
}
func (f *fixture) apply(t *testing.T, p Preview) (any, error) {
	t.Helper()
	in := f.input
	in.Apply = true
	in.Confirm = *p.Confirmation
	return Execute(f.deps, in)
}
func exitOf(e error) int {
	if e == nil {
		return 0
	}
	var v *Error
	if errors.As(e, &v) {
		return v.Exit
	}
	return 1
}
func stateBytes(t *testing.T, f *fixture) []byte {
	t.Helper()
	b, e := os.ReadFile(filepath.Join(f.route.StateRoot, "state.json"))
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func TestReportReadyJourneyPreservesLocalReturn(t *testing.T) {
	f := newFixture(t)
	before, e := os.ReadFile(f.request.Source.Report.Path)
	if e != nil {
		t.Fatal(e)
	}
	p := f.preview(t)
	if _, e = os.Stat(f.route.StateRoot); !os.IsNotExist(e) {
		t.Fatal("preview wrote state")
	}
	v, e := f.apply(t, p)
	if e != nil {
		t.Fatal(e)
	}
	r := v.(Result)
	if r.State != "transport_acknowledged" || r.Knowledge != "reported" || r.NextAction != f.request.Public.NextAction || f.posts != 1 {
		t.Fatalf("bad result %+v", r)
	}
	stored := stateBytes(t, f)
	for _, s := range []string{syntheticSecret, "PRIVATE_SOURCE_SENTINEL"} {
		if bytes.Contains(stored, []byte(s)) {
			t.Fatal("private source or credential leaked")
		}
	}
	for _, operation := range []string{"send", "show"} {
		in := f.input
		if operation == "show" {
			in.Operation = "show"
			in.ID = p.ID
		}
		if _, e = Execute(f.deps, in); e != nil {
			t.Fatal(e)
		}
	}
	if !bytes.Equal(stored, stateBytes(t, f)) {
		t.Fatal("read changed state")
	}
	after, _ := os.ReadFile(f.request.Source.Report.Path)
	if !bytes.Equal(before, after) {
		t.Fatal("local return changed")
	}
	v, e = f.apply(t, p)
	if e != nil || v.(Result).State != "transport_acknowledged" || f.posts != 1 {
		t.Fatal("duplicate dispatched")
	}
}
func TestStrictInputAndPathBoundaries(t *testing.T) {
	for _, variant := range []string{"duplicate", "nested-unknown", "invalid-utf8", "surrogate", "float", "null-revision", "bool-revision", "oversize", "source-duplicate", "source-symlink", "parent-symlink", "state-symlink", "state-permissions", "missing-source", "missing-state-parent", "credential-in-actor", "webhook-in-channel", "mention", "source-worktree-symlink"} {
		t.Run(variant, func(t *testing.T) {
			f := newFixture(t)
			switch variant {
			case "duplicate":
				b, _ := os.ReadFile(f.input.File)
				b = bytes.Replace(b, []byte(`"schema_version":1`), []byte(`"schema_version":1,"schema_version":1`), 1)
				os.WriteFile(f.input.File, b, 0600)
			case "nested-unknown":
				b, _ := os.ReadFile(f.input.File)
				b = bytes.Replace(b, []byte(`"actor_claim":`), []byte(`"unexpected":true,"actor_claim":`), 1)
				os.WriteFile(f.input.File, b, 0600)
			case "invalid-utf8":
				b, _ := os.ReadFile(f.input.File)
				b = append(b, 0xff)
				os.WriteFile(f.input.File, b, 0600)
			case "surrogate":
				b, _ := os.ReadFile(f.input.File)
				b = bytes.Replace(b, []byte(`"Test delivery"`), []byte(`"\ud800"`), 1)
				os.WriteFile(f.input.File, b, 0600)
			case "float", "bool-revision", "null-revision":
				b, _ := os.ReadFile(f.input.File)
				value := map[string]string{"float": "1.0", "bool-revision": "true", "null-revision": "null"}[variant]
				b = bytes.Replace(b, []byte(`"revision":1`), []byte(`"revision":`+value), 1)
				os.WriteFile(f.input.File, b, 0600)
			case "oversize":
				b, _ := os.ReadFile(f.input.File)
				b = append(b, bytes.Repeat([]byte(" "), 64<<10)...)
				os.WriteFile(f.input.File, b, 0600)
			case "source-duplicate":
				p := f.request.Source.Report.Path
				b, _ := os.ReadFile(p)
				b = bytes.Replace(b, []byte(`"run":"run-1"`), []byte(`"run":"run-1","run":"run-1"`), 1)
				os.WriteFile(p, b, 0600)
				f.request.Source.Report.SHA256 = Digest(b)
				f.save(t)
			case "source-symlink":
				p := f.request.Source.Report.Path
				os.Rename(p, p+"-real")
				os.Symlink(p+"-real", p)
			case "parent-symlink":
				p := filepath.Join(f.root, "alias")
				os.Symlink(f.root, p)
				f.request.Source.Report.Path = filepath.Join(p, "report")
				f.save(t)
			case "state-symlink":
				os.Mkdir(filepath.Join(f.root, "actual-state"), 0700)
				os.Symlink(filepath.Join(f.root, "actual-state"), f.route.StateRoot)
			case "state-permissions":
				os.Mkdir(f.route.StateRoot, 0755)
				os.Chmod(f.route.StateRoot, 0755)
			case "missing-source":
				os.Remove(f.request.Source.Report.Path)
			case "missing-state-parent":
				f.route.StateRoot = filepath.Join(f.root, "absent", "state")
				f.save(t)
			case "credential-in-actor":
				f.request.Sender.ActorClaim = syntheticSecret
				f.save(t)
			case "webhook-in-channel":
				f.route.ChannelLabel = syntheticSecret
				f.save(t)
			case "mention":
				f.request.Public.TaskTitle = "Notify @HeRe"
				f.save(t)
			case "source-worktree-symlink":
				p := filepath.Join(f.root, "alias-worktree")
				os.Symlink(f.request.Source.Worktree, p)
				f.request.Source.Worktree = p
				f.save(t)
			}
			if _, e := Execute(f.deps, f.input); e == nil {
				t.Fatal("invalid input accepted")
			}
			if f.posts != 0 {
				t.Fatal("validation dispatched")
			}
		})
	}
}
func TestCredentialURLBoundary(t *testing.T) {
	for _, s := range []string{"http://hooks.slack.com/services/A/B/C", "https://hooks.slack.com:443/services/A/B/C", "https://user@hooks.slack.com/services/A/B/C", "https://hooks.slack.com/services/A/B/C?x=1", "https://hooks.slack.com/services/A/B/C#x", "https://hooks.slack.com/services/A/B/%43", "https://other.test/services/A/B/C", "https://hooks.slack.com/services/A/B/C?", "https://hooks.slack.com/services/A/B/C#"} {
		if credential(s) {
			t.Errorf("invalid credential accepted")
		}
	}
}
func TestSourceDriftAfterReservationReturnsBoundResult(t *testing.T) {
	f := newFixture(t)
	p := f.preview(t)
	f.apply(t, p)
	os.WriteFile(f.request.Source.Report.Path, []byte("changed"), 0600)
	v, e := f.apply(t, p)
	if exitOf(e) != 4 {
		t.Fatal(e)
	}
	r, ok := v.(Result)
	if !ok || r.State != "transport_acknowledged" || r.Freshness != "changed" || f.posts != 1 {
		t.Fatalf("misreported previous attempt: %+v", v)
	}
}
func TestMissingCorruptStateAndUnknownID(t *testing.T) {
	f := newFixture(t)
	in := f.input
	in.Operation = "show"
	in.ID = "ntf_" + strings.Repeat("0", 64)
	if _, e := Execute(f.deps, in); exitOf(e) != 4 {
		t.Fatal(e)
	}
	p := f.preview(t)
	f.apply(t, p)
	path := filepath.Join(f.route.StateRoot, "state.json")
	os.WriteFile(path, []byte("{broken"), 0600)
	if _, e := f.apply(t, p); exitOf(e) != 1 {
		t.Fatalf("corrupt state: %v", e)
	}
	if f.posts != 1 {
		t.Fatal("corrupt state resent")
	}
}
func TestReservationAndReceiptFilesystemFailure(t *testing.T) {
	for _, point := range []string{"before_reservation", "receipt_commit"} {
		t.Run(point, func(t *testing.T) {
			f := newFixture(t)
			p := f.preview(t)
			f.deps.Fault = func(stage string) error {
				if stage == point {
					return os.Chmod(f.route.StateRoot, 0500)
				}
				return nil
			}
			v, e := f.apply(t, p)
			os.Chmod(f.route.StateRoot, 0700)
			if e == nil {
				t.Fatal("write failure reported success")
			}
			if point == "before_reservation" && f.posts != 0 {
				t.Fatal("failed reservation dispatched")
			}
			if point == "receipt_commit" {
				r := v.(Result)
				if r.State != "unknown" || r.Persistence != "reservation_only" || f.posts != 1 {
					t.Fatalf("wrong uncertain outcome %+v", r)
				}
				in := f.input
				in.Operation = "show"
				in.ID = p.ID
				v, e = Execute(f.deps, in)
				if e != nil || v.(Result).State != "unknown" {
					t.Fatal("restart did not preserve uncertainty")
				}
			}
		})
	}
}
func TestNotificationProcessHelper(t *testing.T) {
	config := os.Getenv("PLY_NOTIFICATION_PROCESS_FIXTURE")
	if config == "" {
		return
	}
	var in Input
	if json.Unmarshal([]byte(config), &in) != nil {
		os.Exit(80)
	}
	d := SystemDependencies()
	d.LookupEnv = func(string) (string, bool) { return syntheticSecret, true }
	d.Now = func() time.Time { return time.Date(2026, 10, 3, 19, 0, 0, 0, time.UTC) }
	d.Transport = func(context.Context, string, []byte) Observation {
		f, e := os.OpenFile(os.Getenv("PLY_NOTIFICATION_POST_LOG"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if e != nil {
			os.Exit(81)
		}
		f.Write([]byte("post\n"))
		f.Sync()
		f.Close()
		return Observation{Status: 200, Body: []byte("ok"), Dispatch: "started"}
	}
	d.Fault = func(point string) error {
		if point == os.Getenv("PLY_NOTIFICATION_CRASH_POINT") {
			os.Exit(83)
		}
		return nil
	}
	_, e := Execute(d, in)
	if e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
	}
	os.Exit(exitOf(e))
}
func process(t *testing.T, f *fixture, p Preview, crash string) *exec.Cmd {
	t.Helper()
	in := f.input
	in.Apply = true
	in.Confirm = *p.Confirmation
	b, _ := json.Marshal(in)
	c := exec.Command(os.Args[0], "-test.run=^TestNotificationProcessHelper$")
	c.Env = append(os.Environ(), "PLY_NOTIFICATION_PROCESS_FIXTURE="+string(b), "PLY_NOTIFICATION_CRASH_POINT="+crash, "PLY_NOTIFICATION_POST_LOG="+filepath.Join(f.root, "posts"))
	return c
}
func TestConcurrentProcessesReserveOneAttempt(t *testing.T) {
	f := newFixture(t)
	p := f.preview(t)
	var wg sync.WaitGroup
	codes := make(chan int, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			output, e := process(t, f, p, "").CombinedOutput()
			if len(output) > 0 {
				t.Log(string(output))
			}
			if e == nil {
				codes <- 0
			} else {
				codes <- e.(*exec.ExitError).ExitCode()
			}
		}()
	}
	wg.Wait()
	close(codes)
	success := false
	for code := range codes {
		if code != 0 && code != 4 {
			t.Fatalf("exit %d", code)
		}
		success = success || code == 0
	}
	if !success {
		t.Fatal("no successful process")
	}
	b, _ := os.ReadFile(filepath.Join(f.root, "posts"))
	if string(b) != "post\n" {
		t.Fatal("concurrent dispatch count", string(b))
	}
}
func TestActualProcessDeathBoundaries(t *testing.T) {
	for _, stage := range []string{"before_reservation", "after_reservation", "after_transport"} {
		t.Run(stage, func(t *testing.T) {
			f := newFixture(t)
			p := f.preview(t)
			e := process(t, f, p, stage).Run()
			if e == nil || e.(*exec.ExitError).ExitCode() != 83 {
				t.Fatal(e)
			}
			in := f.input
			in.Operation = "show"
			in.ID = p.ID
			v, e := Execute(f.deps, in)
			if stage == "before_reservation" {
				if exitOf(e) != 4 {
					t.Fatal(e)
				}
				return
			}
			if e != nil || v.(Result).State != "unknown" {
				t.Fatal("crash appeared retryable", v, e)
			}
			v, e = f.apply(t, p)
			if e != nil || v.(Result).State != "unknown" || f.posts != 0 {
				t.Fatal("crash resend", v, e)
			}
			log, _ := os.ReadFile(filepath.Join(f.root, "posts"))
			want := ""
			if stage == "after_transport" {
				want = "post\n"
			}
			if string(log) != want {
				t.Fatal("wrong crash dispatch count")
			}
		})
	}
}

func TestCorruptPriorStateNeverClaimsNotAttempted(t *testing.T) {
	f := newFixture(t)
	p := f.preview(t)
	f.apply(t, p)
	os.WriteFile(filepath.Join(f.route.StateRoot, "state.json"), []byte("{broken"), 0600)
	value, e := f.apply(t, p)
	if e == nil {
		t.Fatal("corrupt state accepted")
	}
	b, _ := json.Marshal(value)
	var out map[string]any
	json.Unmarshal(b, &out)
	if out["state"] == "not_attempted" || out["kind"] != "ply.workflow.notification" {
		t.Fatal("a prior attempt may exist; false absence claim", string(b))
	}
}
func TestSourceChangeAfterReservationNeverDispatches(t *testing.T) {
	f := newFixture(t)
	p := f.preview(t)
	f.deps.Fault = func(point string) error {
		if point == "after_reservation" {
			return os.WriteFile(f.request.Source.Report.Path, []byte("changed before dispatch"), 0600)
		}
		return nil
	}
	value, e := f.apply(t, p)
	if e == nil || f.posts != 0 {
		t.Fatal("changed source dispatched", value, e)
	}
	if value.(Result).State != "unknown" {
		t.Fatal("reservation must remain uncertain")
	}
}
func TestRecordCorruptionFailsClosed(t *testing.T) {
	f := newFixture(t)
	p := f.preview(t)
	f.apply(t, p)
	path := filepath.Join(f.route.StateRoot, "state.json")
	b := stateBytes(t, f)
	b = bytes.Replace(b, []byte(`"state":"transport_acknowledged"`), []byte(`"state":"rejected"`), 1)
	os.WriteFile(path, b, 0600)
	in := f.input
	in.Operation = "retry"
	in.ID = p.ID
	if _, e := Execute(f.deps, in); e == nil {
		t.Fatal("corrupt acknowledgement became retryable")
	}
}

func TestDeletedStateDoesNotEraseRootBinding(t *testing.T) {
	f := newFixture(t)
	p := f.preview(t)
	f.apply(t, p)
	os.Remove(filepath.Join(f.route.StateRoot, "state.json"))
	if _, e := f.apply(t, p); e == nil || f.posts != 1 {
		t.Fatal("missing state silently allowed resend")
	}
}
func TestStateRootReplacementAfterReservationNeverDispatches(t *testing.T) {
	f := newFixture(t)
	p := f.preview(t)
	f.deps.Fault = func(stage string) error {
		if stage == "after_reservation" {
			if e := os.Rename(f.route.StateRoot, f.route.StateRoot+"-original"); e != nil {
				return e
			}
			return os.Mkdir(f.route.StateRoot, 0700)
		}
		return nil
	}
	value, e := f.apply(t, p)
	if e == nil || f.posts != 0 || value.(Result).State != "unknown" {
		t.Fatal("root swap dispatched")
	}
}
func TestDocumentationFixturePassesRealValidator(t *testing.T) {
	f := newFixture(t)
	base := filepath.Join("..", "..", "test", "fixtures", "notification-docs")
	docs, e := os.ReadFile(filepath.Join("..", "..", "docs", "workflow-notification.md"))
	if e != nil {
		t.Fatal(e)
	}
	for _, name := range []string{"handoff.md", "start-receipt.json", "report.json", "route.json", "request.json"} {
		b, e := os.ReadFile(filepath.Join(base, name))
		if e != nil {
			t.Fatal(e)
		}
		if name == "route.json" || name == "request.json" {
			if !bytes.Contains(docs, bytes.TrimSpace(b)) {
				t.Fatal("docs JSON differs from tested fixture", name)
			}
		}
		b = bytes.ReplaceAll(b, []byte("/private/tmp/ply-notification-example"), []byte(f.root))
		if e = os.WriteFile(filepath.Join(f.root, name), b, 0600); e != nil {
			t.Fatal(e)
		}
	}
	f.input.Route = filepath.Join(f.root, "route.json")
	f.input.File = filepath.Join(f.root, "request.json")
	p := f.preview(t)
	if _, e = f.apply(t, p); e != nil {
		t.Fatal(e)
	}
}

func TestExternalMetadataAcceptsStandardJSONNumbers(t *testing.T) {
	f := newFixture(t)
	b, _ := os.ReadFile(f.request.Source.Report.Path)
	b = append(b[:len(b)-1], []byte(`,"elapsed_seconds":1.25,"measurements":[1e30,-0.1,9223372036854775808],"flag":true}`)...)
	os.WriteFile(f.request.Source.Report.Path, b, 0600)
	f.request.Source.Report.SHA256 = Digest(b)
	f.save(t)
	p := f.preview(t)
	if _, e := f.apply(t, p); e != nil {
		t.Fatal(e)
	}
}
