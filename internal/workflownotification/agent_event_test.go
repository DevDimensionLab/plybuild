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

func agentFixture(t *testing.T, kind string) *fixture {
	t.Helper()
	f := newFixture(t)
	f.request.SchemaVersion = 2
	f.request.Gate = nil
	f.request.Public.Summary = "The local return is preserved."
	f.request.Public.NextActor = "coordinator"
	eventKind := agentEventKind
	f.request.Event = &Event{ID: "event-1", Type: kind, Phase: "after_start", OccurredAt: "2026-10-04T14:00:00Z", Record: Locator{Path: filepath.Join(f.root, "event.json"), Kind: &eventKind}}
	if kind == "feedback_required" {
		f.request.Public.Summary = "Which target should this task use?"
		f.request.Public.NextActor = "user"
	}
	saveAgent(t, f)
	return f
}
func saveAgent(t *testing.T, f *fixture) {
	t.Helper()
	e := f.request.Event
	writeJSON(t, e.Record.Path, eventRecord{agentEventKind, f.request.Source.Activity, f.request.Source.Run, e.ID, e.Type, e.Phase, e.OccurredAt, f.request.Public})
	b, err := os.ReadFile(e.Record.Path)
	if err != nil {
		t.Fatal(err)
	}
	e.Record.SHA256 = Digest(b)
	f.save(t)
}
func TestAgentEventJourneys(t *testing.T) {
	for _, kind := range []string{"agent_finished", "agent_stopped", "feedback_required"} {
		t.Run(kind, func(t *testing.T) {
			f := agentFixture(t, kind)
			if kind == "agent_stopped" {
				f.request.Source.Start = nil
				f.request.Source.Report = nil
				f.request.Event.Phase = "before_start"
				saveAgent(t, f)
			}
			before, _ := os.ReadFile(f.request.Event.Record.Path)
			p := f.preview(t)
			if !strings.Contains(p.Payload.Text, f.request.Public.Summary) || strings.Contains(p.Payload.Text, f.root) {
				t.Fatal("bad public payload", p.Payload)
			}
			if f.request.Public.NextActor == "coordinator" && !strings.Contains(p.Payload.Text, "Coordinator continues") {
				t.Fatal("wrong next actor")
			}
			out, e := f.apply(t, p)
			if e != nil {
				t.Fatal(e)
			}
			r := out.(Result)
			if r.Event == nil || r.Gate != nil || r.State != "transport_acknowledged" || r.Knowledge != "reported" || f.posts != 1 {
				t.Fatal(r)
			}
			b, _ := json.Marshal(r)
			if bytes.Contains(b, []byte(`"gate"`)) {
				t.Fatal("fabricated gate")
			}
			if kind == "agent_stopped" && (bytes.Contains(b, []byte(`"start_receipt"`)) || bytes.Contains(b, []byte(`"report"`))) {
				t.Fatal("fabricated sources")
			}
			after, _ := os.ReadFile(f.request.Event.Record.Path)
			if !bytes.Equal(before, after) {
				t.Fatal("changed event")
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
	for _, variant := range []string{"finished-no-report", "before-start-finished", "before-start-receipt", "native", "bad-time", "bad-type", "empty-summary", "mention", "bad-actor", "feedback-coordinator", "feedback-no-question", "source-run", "event-run", "event-public", "event-unknown", "event-duplicate", "event-null", "event-oversize", "event-symlink", "event-with-gate", "null-receipt", "request-duplicate", "unknown-public", "fraction-schema", "record-kind", "event-in-state"} {
		t.Run(variant, func(t *testing.T) {
			f := agentFixture(t, "feedback_required")
			switch variant {
			case "finished-no-report":
				f.request.Event.Type = "agent_finished"
				f.request.Source.Report = nil
			case "before-start-finished":
				f.request.Event.Type = "agent_finished"
				f.request.Event.Phase = "before_start"
				f.request.Source.Start = nil
			case "before-start-receipt":
				f.request.Event.Type = "agent_stopped"
				f.request.Event.Phase = "before_start"
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
			case "source-run":
				f.request.Source.Run = "different"
			case "record-kind":
				v := "Other@1"
				f.request.Event.Record.Kind = &v
			case "event-in-state":
				os.Mkdir(f.route.StateRoot, 0700)
				f.request.Event.Record.Path = filepath.Join(f.route.StateRoot, "event.json")
			}
			saveAgent(t, f)
			if strings.HasPrefix(variant, "event-") && variant != "event-in-state" && variant != "event-with-gate" {
				path := f.request.Event.Record.Path
				b, _ := os.ReadFile(path)
				switch variant {
				case "event-run":
					b = bytes.Replace(b, []byte(`"run-1"`), []byte(`"other"`), 1)
				case "event-public":
					b = bytes.Replace(b, []byte("Which target"), []byte("Which branch"), 1)
				case "event-unknown":
					b = append([]byte(`{"unknown":true,`), b[1:]...)
				case "event-duplicate":
					b = append([]byte(`{"kind":"duplicate",`), b[1:]...)
				case "event-null":
					b = []byte("null")
				case "event-oversize":
					b = bytes.Repeat([]byte(" "), (1<<20)+1)
				case "event-symlink":
					os.Rename(path, path+".real")
					os.Symlink(path+".real", path)
				}
				if variant != "event-symlink" {
					os.WriteFile(path, b, 0600)
					f.request.Event.Record.SHA256 = Digest(b)
					f.save(t)
				}
			}
			b, _ := os.ReadFile(f.input.File)
			switch variant {
			case "event-with-gate":
				b = append([]byte(`{"gate":null,`), b[1:]...)
			case "null-receipt":
				b = bytes.Replace(b, []byte(`"start_receipt":{`), []byte(`"start_receipt":null,"extra":{`), 1)
			case "request-duplicate":
				b = append([]byte(`{"schema_version":2,`), b[1:]...)
			case "unknown-public":
				b = bytes.Replace(b, []byte(`"public":{`), []byte(`"public":{"extra":true,`), 1)
			case "fraction-schema":
				b = bytes.Replace(b, []byte(`"schema_version":2`), []byte(`"schema_version":2.0`), 1)
			}
			os.WriteFile(f.input.File, b, 0600)
			_, err := Execute(f.deps, f.input)
			wantExit := 2
			if variant == "event-in-state" {
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
	f := agentFixture(t, "feedback_required")
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
	f.request.Public.Summary = "Which changed target should this task use?"
	saveAgent(t, f)
	_, e = Execute(f.deps, f.input)
	if exitOf(e) != 4 || f.posts != 1 {
		t.Fatal("changed bytes dispatched", e)
	}
	out, e = Execute(f.deps, Input{Operation: "show", Route: f.input.Route, ID: p.ID})
	if e != nil || out.(Result).Freshness != "changed" || out.(Result).State != "unknown" {
		t.Fatal(out, e)
	}
}
func TestAgentEventSourceDriftAfterPreview(t *testing.T) {
	f := agentFixture(t, "agent_finished")
	p := f.preview(t)
	b, _ := os.ReadFile(f.request.Event.Record.Path)
	os.WriteFile(f.request.Event.Record.Path, append(b, '\n'), 0600)
	_, e := f.apply(t, p)
	if exitOf(e) != 4 || f.posts != 0 {
		t.Fatal("event drift dispatched", e)
	}
}
func TestAgentIdentityNamespaceAndConcurrentProcesses(t *testing.T) {
	f := newFixture(t)
	v1 := identity(f.request)
	g := agentFixture(t, "agent_finished")
	g.request.Event.ID = f.request.Gate.ID
	saveAgent(t, g)
	if identity(g.request) == v1 {
		t.Fatal("v1 namespace collision")
	}
	p := g.preview(t)
	a, b := process(t, g, p, ""), process(t, g, p, "")
	if err := a.Start(); err != nil {
		t.Fatal(err)
	}
	if err := b.Start(); err != nil {
		t.Fatal(err)
	}
	a.Wait()
	b.Wait()
	posts, _ := os.ReadFile(filepath.Join(g.root, "posts"))
	if string(posts) != "post\n" {
		t.Fatal("concurrent event POST count", string(posts))
	}
}

func TestCorruptV1StateRetainsGateWireContract(t *testing.T) {
	f := newFixture(t)
	p := f.preview(t)
	f.apply(t, p)
	os.WriteFile(filepath.Join(f.route.StateRoot, "state.json"), []byte("broken"), 0600)
	out, e := Execute(f.deps, f.input)
	if e == nil {
		t.Fatal("corrupt state accepted")
	}
	b, _ := json.Marshal(out)
	if !bytes.Contains(b, []byte(`"gate":{`)) || bytes.Contains(b, []byte(`"event"`)) {
		t.Fatalf("v1 unknown result changed wire contract: %s", b)
	}
}
func TestMalformedStoredIdentityDoesNotPanic(t *testing.T) {
	f := newFixture(t)
	p := f.preview(t)
	f.apply(t, p)
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
	delete(db.Records, p.ID)
	r.Request.Gate = nil
	db.Records[""] = r
	if err = writeDatabase(dir, db); err != nil {
		t.Fatal(err)
	}
	if _, err = Execute(f.deps, Input{Operation: "show", ID: p.ID, Route: f.input.Route}); err == nil {
		t.Fatal("malformed identity accepted")
	}
}

func TestAfterStartStopAndFeedbackNeedNoFabricatedSources(t *testing.T) {
	for _, kind := range []string{"agent_stopped", "feedback_required"} {
		t.Run(kind, func(t *testing.T) {
			f := agentFixture(t, kind)
			f.request.Source.Start = nil
			f.request.Source.Report = nil
			saveAgent(t, f)
			p := f.preview(t)
			out, err := f.apply(t, p)
			if err != nil || out.(Result).State != "transport_acknowledged" || f.posts != 1 {
				t.Fatal(out, err)
			}
			data, _ := json.Marshal(out)
			if bytes.Contains(data, []byte(`"start_receipt"`)) || bytes.Contains(data, []byte(`"report"`)) {
				t.Fatal("fabricated source")
			}
		})
	}
}
