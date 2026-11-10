package workflownotification

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAgentEventJourneys(t *testing.T) {
	for _, c := range []struct{ kind, status, phase string }{{"agent_finished", "done", "after_start"}, {"agent_stopped", "stopped", "before_start"}, {"agent_stopped", "stopped", "after_start"}, {"feedback_required", "needs_answer", "after_start"}} {
		t.Run(c.kind+"/"+c.phase, func(t *testing.T) {
			f, _ := compactFixture(t, c.kind, c.status, "claude")
			f.request.Event.Phase = c.phase
			f.save(t)
			before, _ := os.ReadFile(f.input.File)
			p := f.preview(t)
			if !strings.Contains(p.Payload.Text, f.request.Public.Summary) || strings.Contains(p.Payload.Text, f.root) {
				t.Fatal("bad public payload", p.Payload)
			}
			out, e := f.apply(t, p)
			if e != nil {
				t.Fatal(e)
			}
			r := out.(Result)
			if r.Event == nil || r.State != "transport_acknowledged" || r.Knowledge != "reported" || f.posts != 1 {
				t.Fatal(r)
			}
			b, _ := json.Marshal(r)
			for _, field := range []string{`"gate"`, `"handoff"`, `"start_receipt"`, `"report"`, `"record"`} {
				if bytes.Contains(b, []byte(field)) {
					t.Fatal("fabricated", field)
				}
			}
			after, _ := os.ReadFile(f.input.File)
			if !bytes.Equal(before, after) {
				t.Fatal("changed request")
			}
			out, e = f.apply(t, p)
			if e != nil || out.(Result).State != "transport_acknowledged" || f.posts != 1 {
				t.Fatal("replayed POST")
			}
			in := Input{Operation: "show", Route: f.input.Route, ID: p.ID}
			d := f.deps
			d.LookupEnv = func(string) (string, bool) { return "", false }
			out, e = Execute(d, in)
			if e != nil || out.(Result).Freshness != "current" {
				t.Fatal("credential-free readback", out, e)
			}
		})
	}
}
func TestAgentEventValidationBeforeReservation(t *testing.T) {
	for _, variant := range []string{"before-start-finished", "native", "bad-time", "bad-type", "empty-summary", "mention", "bad-actor", "feedback-coordinator", "feedback-no-question", "event-with-gate", "request-duplicate", "unknown-public", "fraction-schema", "worktree-in-state"} {
		t.Run(variant, func(t *testing.T) {
			f, _ := compactFixture(t, "feedback_required", "needs_answer", "claude")
			switch variant {
			case "before-start-finished":
				f.request.Event.Type = "agent_finished"
				f.request.Event.Phase = "before_start"
				f.request.Public.Status = "done"
			case "native":
				f.request.Source.Kind = "native"
			case "bad-time":
				f.request.Event.OccurredAt = "2026-10-04T14:00:00.1Z"
			case "bad-type":
				f.request.Event.Type = "report_ready"
			case "empty-summary":
				f.request.Public.Summary = ""
			case "mention":
				f.request.Public.Summary = "Notify @Here?"
			case "bad-actor":
				f.request.Public.NextActor = "agent"
			case "feedback-coordinator":
				f.request.Public.NextActor = "coordinator"
			case "feedback-no-question":
				f.request.Public.Summary = "Please clarify"
			case "worktree-in-state":
				f.request.Source.Worktree = filepath.Join(f.route.StateRoot, "worktree")
				os.MkdirAll(f.request.Source.Worktree, 0700)
			}
			f.save(t)
			b, _ := os.ReadFile(f.input.File)
			switch variant {
			case "event-with-gate":
				b = append([]byte(`{"gate":null,`), b[1:]...)
			case "request-duplicate":
				b = append([]byte(`{"schema_version":4,`), b[1:]...)
			case "unknown-public":
				b = bytes.Replace(b, []byte(`"public":{`), []byte(`"public":{"extra":true,`), 1)
			case "fraction-schema":
				b = bytes.Replace(b, []byte(`"schema_version":4`), []byte(`"schema_version":4.0`), 1)
			}
			os.WriteFile(f.input.File, b, 0600)
			_, err := Execute(f.deps, f.input)
			wantExit := 2
			if variant == "worktree-in-state" {
				wantExit = 1
			} // Existing nonempty state is rejected as corrupt before source validation.
			if exitOf(err) != wantExit || f.posts != 0 {
				t.Fatal("not rejected before reservation", err)
			}
			if _, err = os.Stat(filepath.Join(f.route.StateRoot, "state.json")); !os.IsNotExist(err) {
				t.Fatal("reserved invalid event")
			}
		})
	}
}
func TestAgentUnknownConflictAndFreshness(t *testing.T) {
	f, _ := compactFixture(t, "feedback_required", "needs_answer", "claude")
	f.deps.Transport = func(context.Context, string, []byte) Observation {
		f.posts++
		return Observation{Failed: true, Dispatch: "started"}
	}
	p := f.preview(t)
	out, e := f.apply(t, p)
	if exitOf(e) != 5 || out.(Result).State != "unknown" {
		t.Fatal(out, e)
	}
	out, e = f.apply(t, p)
	if e != nil || out.(Result).State != "unknown" || f.posts != 1 {
		t.Fatal("unknown replay")
	}
	_, e = Execute(f.deps, Input{Operation: "retry", Route: f.input.Route, ID: p.ID})
	if exitOf(e) != 4 {
		t.Fatal("unknown retry")
	}
	f.request.Public.Summary = "Hvilket annet mål?"
	f.save(t)
	_, e = Execute(f.deps, f.input)
	if exitOf(e) != 4 || f.posts != 1 {
		t.Fatal("changed bytes dispatched", e)
	}
	out, e = Execute(f.deps, Input{Operation: "show", Route: f.input.Route, ID: p.ID})
	if e != nil || out.(Result).Freshness != "changed" || out.(Result).State != "unknown" {
		t.Fatal(out, e)
	}
}
func TestAgentEventRequestDriftAfterPreview(t *testing.T) {
	f := newFixture(t)
	p := f.preview(t)
	b, _ := os.ReadFile(f.input.File)
	os.WriteFile(f.input.File, append(b, '\n'), 0600)
	_, e := f.apply(t, p)
	if exitOf(e) != 4 || f.posts != 0 {
		t.Fatal("request drift dispatched", e)
	}
}
func TestCorruptStateRetainsEventWireContract(t *testing.T) {
	f := newFixture(t)
	p := f.preview(t)
	f.apply(t, p)
	os.WriteFile(filepath.Join(f.route.StateRoot, "state.json"), []byte("broken"), 0600)
	out, e := Execute(f.deps, f.input)
	if e == nil {
		t.Fatal("corrupt state accepted")
	}
	b, _ := json.Marshal(out)
	if !bytes.Contains(b, []byte(`"event":{`)) || bytes.Contains(b, []byte(`"gate"`)) {
		t.Fatalf("unknown result changed wire contract: %s", b)
	}
}
func TestMalformedStoredIdentityDoesNotPanic(t *testing.T) {
	for _, variant := range []string{"nil-record", "malformed-bytes", "missing-event", "historical-identity-mismatch", "invalid-key", "identity-mismatch"} {
		t.Run(variant, func(t *testing.T) {
			f := newFixture(t)
			p := f.preview(t)
			if _, err := f.apply(t, p); err != nil {
				t.Fatal(err)
			}
			dir, err := stateDirectory(f.route.StateRoot, false)
			if err != nil {
				t.Fatal(err)
			}
			defer dir.Close()
			db, err := readDatabase(dir)
			if err != nil {
				t.Fatal(err)
			}
			r := db.Records[p.ID]
			switch variant {
			case "nil-record":
				db.Records[p.ID] = nil
			case "malformed-bytes":
				r.RequestBytes = []byte("{broken")
			case "missing-event":
				r.Request.Event = nil
				r.RequestBytes, err = json.Marshal(r.Request)
			case "historical-identity-mismatch":
				r.Request.Source.Run, r.Request.SchemaVersion = "other-run", 3
				r.RequestBytes, err = json.Marshal(r.Request)
			case "invalid-key":
				delete(db.Records, p.ID)
				db.Records[""] = r
			case "identity-mismatch":
				r.Request.Source.Run = "other-run"
				r.RequestBytes, err = json.Marshal(r.Request)
			}
			if err != nil {
				t.Fatal(err)
			}
			if err = writeDatabase(dir, db); err != nil {
				t.Fatal(err)
			}
			if _, err = Execute(f.deps, Input{Operation: "show", ID: p.ID, Route: f.input.Route}); err == nil || f.posts != 1 {
				t.Fatal("malformed state accepted or dispatched", err)
			}
		})
	}
}
