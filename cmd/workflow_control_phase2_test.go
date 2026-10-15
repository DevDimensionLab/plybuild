package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/weekrun"
)

func (f *mr05CLI) put(path string, v any) {
	f.t.Helper()
	if e := os.MkdirAll(filepath.Dir(path), 0700); e != nil {
		f.t.Fatal(e)
	}
	b, e := taskrun.Canonical(v)
	if e != nil {
		f.t.Fatal(e)
	}
	if e = os.WriteFile(path, b, 0600); e != nil {
		f.t.Fatal(e)
	}
}

// Retained fixtures exercise the real process CLI without launching a scheduler.
func (f *mr05CLI) week(phase string) (string, string) {
	c := weekrun.Contract{Kind: "ply.workflow.week-contract", SchemaVersion: 2, Name: "cli-week", Workspace: f.root, ReturnWorktree: f.root, ReturnRef: "refs/heads/epic", ProjectID: "test", RepoID: "test", EpicID: "epic", Selected: []weekrun.Selection{}, Model: "gpt-6-astra", Effort: "max", DurationSeconds: 3600, MaxAttempts: 1, AgentTimeoutSeconds: 300, MaxProviderRetries: 1, RuntimeControl: &weekrun.RuntimeControl{SchemaVersion: 1, Mode: taskrun.BoundControlMode}}
	id := "week_" + strings.TrimPrefix(mr05Digest([]string{f.root, c.Name}), "sha256:")[:32]
	cc := taskrun.NewControlContract("week", id, mr05Digest(c))
	binary, e := os.ReadFile(f.binary)
	if e != nil {
		f.t.Fatal(e)
	}
	h := sha256.Sum256(binary)
	now := time.Now().UTC()
	b := weekrun.Binding{Contract: c, ID: id, SHA256: mr05Digest(c), ControlContract: cc, StartedAt: now.Format(time.RFC3339Nano), Deadline: now.Add(time.Hour).Format(time.RFC3339Nano), Control: taskrun.Executable{Path: f.binary, SHA256: "sha256:" + hex.EncodeToString(h[:])}, Environment: map[string]string{}}
	s := weekrun.State{Kind: "ply.workflow.week", SchemaVersion: 1, ID: id, ContractSHA256: b.SHA256, Phase: phase, StartedAt: b.StartedAt, Deadline: b.Deadline, UpdatedAt: b.StartedAt, LastProgressAt: b.StartedAt, DesiredAction: "continue", Entries: []weekrun.Entry{}, StartContract: &c, ControlContract: cc}
	dir := filepath.Join(f.root, ".ply", "week-runs", id)
	f.put(filepath.Join(dir, "binding.json"), b)
	f.put(filepath.Join(dir, "state.json"), s)
	f.put(filepath.Join(dir, "control-gate.json"), taskrun.ParentControlGate{Contract: cc, Deadline: b.Deadline, DesiredAction: s.DesiredAction})
	return id, dir
}

func mr05UnavailableWeekJSON(t *testing.T) {
	for _, scenario := range []string{"absent", "corrupt_state"} {
		t.Run(scenario, func(t *testing.T) {
			f := mr05Candidate(t)
			id := "week_" + strings.Repeat("a", 32)
			want := "control_not_found"
			if scenario == "corrupt_state" {
				var dir string
				id, dir = f.week("ready")
				if e := os.WriteFile(filepath.Join(dir, "state.json"), []byte(`{"secret":"must-not-leak-control-fixture"`), 0600); e != nil {
					t.Fatal(e)
				}
				want = "control_outcome_unknown"
			}
			before := mr05Snapshot(t, f.root)
			raw, stderr, e := f.run("workflow", "control", "show", "--week", id, "--workspace", f.root, "--format", "json")
			if e == nil {
				t.Error("unavailable week returned success")
			}
			var got map[string]any
			if e = json.Unmarshal(raw, &got); e != nil {
				t.Fatal("no JSON response", string(raw), stderr, e)
			}
			failure, ok := got["error"].(map[string]any)
			if !ok || failure["code"] != want || failure["detail"] == "" {
				t.Error("missing safe JSON error", got["error"])
			}
			if got["state_revision"] != nil || got["run_id"] != nil || got["owner_generation"] != nil || got["control_contract"] != nil || len(got["actions"].([]any)) != 0 || got["unavailable_reason"] == nil {
				t.Error("unknown authority was not explicitly closed", got)
			}
			if strings.Contains(string(raw)+stderr, "must-not-leak-control-fixture") {
				t.Error("private diagnostic leaked")
			}
			if !reflect.DeepEqual(before, mr05Snapshot(t, f.root)) {
				t.Error("failure read mutated files")
			}
			mr05Schema(t, "ply.workflow.control-capabilities", raw)
		})
	}
}

func mr05PublicObservationStates(t *testing.T) {
	t.Run("standalone_blocked", func(t *testing.T) {
		f := mr05Candidate(t)
		id := f.oldRun()
		path := filepath.Join(f.root, ".ply", "workflow-runs", "requests", id+".json")
		raw, e := os.ReadFile(path)
		if e != nil {
			t.Fatal(e)
		}
		var stored map[string]any
		if e = json.Unmarshal(raw, &stored); e != nil {
			t.Fatal(e)
		}
		w := taskrun.WorkerStatus{SessionID: "ply-worker:" + id, Phase: "blocked", Stage: "reviewer", Generation: 2}
		mr05EmptySlices(reflect.ValueOf(&w).Elem())
		var request taskrun.WorkflowRequest
		if e = json.Unmarshal(mustMR05JSON(stored["request"]), &request); e != nil {
			t.Fatal(e)
		}
		request.Envelope = taskrun.Envelope{Kind: "ply.workflow.headless-run-request", SchemaVersion: 2}
		request.Runtime.Mode = "headless"
		request.Worker = &taskrun.WorkerContract{SchemaVersion: 1, MaxAttempts: 1, AgentTimeoutSeconds: 300}
		request.Delivery = &taskrun.DeliveryContract{}
		stored["request"] = request
		stored["result"].(map[string]any)["request_sha256"] = mr05Digest(request)
		delivery := taskrun.DeliveryState{Phase: "working"}
		mr05EmptySlices(reflect.ValueOf(&delivery).Elem())
		stored["result"].(map[string]any)["delivery"] = delivery
		stored["result"].(map[string]any)["worker"] = w
		f.put(path, stored)
		before := mr05Snapshot(t, f.root)
		raw, stderr, e := f.run("workflow", "control", "show", id, "--workspace", f.root, "--format", "json")
		if e != nil {
			t.Fatal(e, stderr, string(raw))
		}
		var got map[string]any
		if e = json.Unmarshal(raw, &got); e != nil {
			t.Fatal(e)
		}
		if got["phase"] != "blocked" || got["runtime_observation"].(map[string]any)["state"] != "unknown" {
			t.Error("standalone internal blocked phase escaped", got)
		}
		mr05Schema(t, "ply.workflow.control-capabilities", raw)
		if !reflect.DeepEqual(before, mr05Snapshot(t, f.root)) {
			t.Error("standalone read wrote files")
		}
	})
	for _, phase := range []string{"ready", "pausing", "stopping", "blocked", "expired", "completed", "waiting", "paused", "stopped", "running"} {
		t.Run(phase, func(t *testing.T) {
			f := mr05Candidate(t)
			id, _ := f.week(phase)
			// A historical receipt must not override the separate present observation.
			basis := taskrun.NewControlBasis(f.root, taskrun.ControlTarget{Kind: "week", ID: id}, time.Now())
			basis.StateRevision = mr05Digest("historical-state")
			r := taskrun.ControlRequest{Kind: "ply.workflow.control-request", SchemaVersion: 1, RequestID: "state-observation-0001", Workspace: f.root}
			r.Input.Target = basis.Scope
			r.Input.Action = "pause"
			r.Validity.IssuedAt = time.Now().UTC().Format(time.RFC3339Nano)
			r.Validity.ExpiresAt = time.Now().Add(time.Minute).UTC().Format(time.RFC3339Nano)
			bound := basis
			bound.Observation.ObservedAt = ""
			bound.OwnerLock = ""
			bound.SupervisorLock = ""
			r.Binding = mr05Digest(struct {
				Request taskrun.ControlRequest
				Basis   taskrun.ControlBasis
			}{r, bound})
			record := taskrun.ControlRecord{Request: r, Digest: mr05Digest(r), Sequence: mr05Digest(r), Basis: basis, AcceptedAt: r.Validity.IssuedAt}
			dir := filepath.Join(f.root, ".ply", "runtime-controls", "v1", r.RequestID)
			f.put(filepath.Join(dir, "request.json"), record)
			f.put(filepath.Join(dir, "receipt.json"), taskrun.ControlReceipt{RequestDigest: record.Digest, Sequence: record.Sequence, Effect: "paused", ObservedAt: r.Validity.IssuedAt})
			before := mr05Snapshot(t, f.root)
			for _, verb := range []string{"show", "status"} {
				args := []string{"workflow", "control", verb}
				if verb == "show" {
					args = append(args, "--week", id)
				} else {
					args = append(args, r.RequestID)
				}
				args = append(args, "--workspace", f.root, "--format", "json")
				raw, stderr, e := f.run(args...)
				if e != nil {
					t.Fatal(verb, e, stderr, string(raw))
				}
				var got map[string]any
				if e = json.Unmarshal(raw, &got); e != nil {
					t.Fatal(e)
				}
				state := got["runtime_observation"].(map[string]any)["state"]
				allowed := map[any]bool{"running": true, "waiting": true, "paused": true, "stopped": true, "completed": true, "unknown": true}
				if !allowed[state] {
					t.Error("public observation escaped promised enum", verb, phase, state)
				}
				want := map[string]string{"ready": "unknown", "pausing": "running", "stopping": "running", "blocked": "unknown", "expired": "stopped"}[phase]
				if want == "" {
					want = phase
				}
				if state != want {
					t.Error("nonconservative observation", verb, phase, state, "want", want)
				}
				kind := "ply.workflow.control-capabilities"
				if verb == "status" {
					kind = "ply.workflow.control-status"
					if got["control_effect"] != "paused" {
						t.Error("historical receipt changed")
					}
				} else if got["phase"] != phase {
					t.Error("internal phase lost")
				}
				mr05Schema(t, kind, raw)
			}
			if !reflect.DeepEqual(before, mr05Snapshot(t, f.root)) {
				t.Error("show/status wrote files")
			}
		})
	}
}
