package weekrun

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func mr05ParentReceiptRecovery(t *testing.T) {
	f := newMR05Week(t)
	r := weekPreview(t, f, "pause", "parent-receipt-recovery-0001")
	weekApply(t, f, r)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e := f.m.Serve(ctx, f.c.Workspace, f.b.ID); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	f.reload(t)
	f.m.Fault = func(point string) error {
		if point == "control_before_receipt" {
			return errors.New("receipt publication interrupted")
		}
		return nil
	}
	if e := f.m.Tick(f.b, f.dir, &f.s); e == nil {
		t.Fatal("publication fault missed")
	}
	f.reload(t)
	if f.s.ControlReceipt == nil {
		t.Fatal("no durable owner receipt to recover")
	}
	retained := *f.s.ControlReceipt
	status, _ := f.m.Controls().Status(f.c.Workspace, r.RequestID)
	if status.ControlEffect != "unknown" {
		t.Fatal("missing publication concealed", status)
	}
	f.m.Fault = nil
	if e := f.m.Serve(ctx, f.c.Workspace, f.b.ID); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	f.reload(t)
	if e := f.m.Tick(f.b, f.dir, &f.s); e != nil {
		t.Error("positive saved receipt cannot recover through a later owner", e)
	}
	status, e := f.m.Controls().Status(f.c.Workspace, r.RequestID)
	if e != nil || status.Receipt == nil || taskrun.ControlRevision(*status.Receipt) != taskrun.ControlRevision(retained) || len(f.n.starts) != 0 {
		t.Error("receipt recovery lost original evidence", status, e)
	}
}

func mr05ChildConsumptionDrift(t *testing.T) {
	for _, action := range []string{"resume", "pause"} {
		for _, mutation := range []string{"provider", "result.json", "stdout.jsonl", "stderr.log", "session", "job", "unchanged"} {
			t.Run(action+"/"+mutation, func(t *testing.T) {
				f, id, provider, job, _ := mr05NativeChild(t)
				if action == "pause" {
					g, e := taskrun.ReadParentControl(f.c.Workspace, f.b.ControlContract)
					if e != nil {
						t.Fatal(e)
					}
					f.s.Phase, f.s.DesiredAction = "running", "continue"
					g.DesiredAction = "continue"
					if e = f.m.ControlGuard(f.c.Workspace, f.b.ID, func() error { return f.m.saveBound(f.b, f.dir, &f.s, &g) }); e != nil {
						t.Fatal(e)
					}
				}
				r := weekPreview(t, f, action, "after-reservation-0001")
				weekApply(t, f, r)
				switch mutation {
				case "unchanged":
				case "provider":
					if e := os.WriteFile(provider, []byte("changed"), 0700); e != nil {
						t.Fatal(e)
					}
				case "session", "job":
					path := filepath.Join(f.c.Workspace, ".ply", "workflow-runs", "runs", id, "state.json")
					var state map[string]any
					if e := readJSON(path, &state); e != nil {
						t.Fatal(e)
					}
					w := state["result"].(map[string]any)["worker"].(map[string]any)
					if mutation == "session" {
						w["session_id"] = "different-worker-session"
					} else {
						w["job"] = nil
					}
					if e := publish(path, state, true); e != nil {
						t.Fatal(e)
					}
				default:
					if e := os.Remove(filepath.Join(job, mutation)); e != nil {
						t.Fatal(e)
					}
				}
				err := f.m.Tick(f.b, f.dir, &f.s)
				status, se := f.m.Controls().Status(f.c.Workspace, r.RequestID)
				g, ge := taskrun.ReadParentControl(f.c.Workspace, f.b.ControlContract)
				if ge != nil {
					t.Fatal(ge)
				}
				if mutation == "unchanged" {
					if err != nil || se != nil || status.Receipt == nil {
						t.Fatal("valid boundary refused", status, err, se)
					}
					return
				}
				if err == nil || status.Receipt != nil || status.ControlEffect != "unknown" || status.RequestState != "outcome_unknown" || g.ResumeAllowed {
					t.Error("drift produced a positive boundary or remained pending", mutation, status, err)
				}
				if *f.m.Native.(mr05NativeQueue).effects != 0 {
					t.Error("drift launched native effect")
				}
			})
		}
	}
}

func mr05ParentRuntimeDrift(t *testing.T) {
	for _, action := range []string{"pause", "resume"} {
		t.Run(action, func(t *testing.T) {
			f := newMR05Week(t)
			if action == "resume" {
				r := weekPreview(t, f, "pause", "before-runtime-drift-0001")
				weekApply(t, f, r)
				f.tick(t)
			}
			r := weekPreview(t, f, action, "parent-runtime-drift-0001")
			weekApply(t, f, r)
			if e := os.Chmod(f.b.Control.Path, 0700); e != nil {
				t.Fatal(e)
			}
			if e := os.WriteFile(f.b.Control.Path, []byte("changed runtime"), 0700); e != nil {
				t.Fatal(e)
			}
			err := f.m.Tick(f.b, f.dir, &f.s)
			status, e := f.m.Controls().Status(f.c.Workspace, r.RequestID)
			g, ge := taskrun.ReadParentControl(f.c.Workspace, f.b.ControlContract)
			if err == nil || e == nil || status.ControlEffect != "unknown" || status.Receipt != nil || ge != nil || g.ResumeAllowed || len(f.n.starts) != 0 {
				t.Error("changed parent runtime confirmed a boundary", status, err, e, ge)
			}
		})
	}
}

func mr05CompletedIntegrationProof(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "positive", true: "missing_result"}[missing], func(t *testing.T) {
			f, id, _, job, _ := mr05NativeChild(t)
			statePath := filepath.Join(f.c.Workspace, ".ply", "workflow-runs", "runs", id, "state.json")
			var state, intent, result map[string]any
			if e := readJSON(statePath, &state); e != nil {
				t.Fatal(e)
			}
			if e := readJSON(filepath.Join(job, "intent.json"), &intent); e != nil {
				t.Fatal(e)
			}
			if e := readJSON(filepath.Join(job, "result.json"), &result); e != nil {
				t.Fatal(e)
			}
			intent["kind"] = "integrate"
			intentHash := hash(canonical(intent))
			result["intent_sha256"] = intentHash
			result["native_completed"] = true
			w := state["result"].(map[string]any)["worker"].(map[string]any)
			w["stage"] = "integrate"
			w["job"].(map[string]any)["kind"] = "integrate"
			w["job"].(map[string]any)["intent_sha256"] = intentHash
			// The fixture is seeded before approval; post-approval mutation below is
			// the actual completion read by NativeAdapter, not a live integration.
			for path, value := range map[string]any{statePath: state, filepath.Join(job, "intent.json"): intent, filepath.Join(job, "result.json"): result} {
				if e := publish(path, value, true); e != nil {
					t.Fatal(e)
				}
			}
			g, e := taskrun.ReadParentControl(f.c.Workspace, f.b.ControlContract)
			if e != nil {
				t.Fatal(e)
			}
			f.s.Phase, f.s.DesiredAction = "running", "continue"
			g.DesiredAction = "continue"
			if e = f.m.ControlGuard(f.c.Workspace, f.b.ID, func() error { return f.m.saveBound(f.b, f.dir, &f.s, &g) }); e != nil {
				t.Fatal(e)
			}
			r := weekPreview(t, f, "pause", "completed-proof-0001")
			weekApply(t, f, r)
			w["job"] = nil
			w["stage"] = "completed"
			w["phase"] = "completed"
			state["result"].(map[string]any)["delivery"].(map[string]any)["phase"] = "completed"
			if e = publish(statePath, state, true); e != nil {
				t.Fatal(e)
			}
			if missing {
				if e = os.Remove(filepath.Join(job, "result.json")); e != nil {
					t.Fatal(e)
				}
			}
			err := f.m.Tick(f.b, f.dir, &f.s)
			status, se := f.m.Controls().Status(f.c.Workspace, r.RequestID)
			if missing {
				if err == nil || se == nil || status.Receipt != nil || status.ControlEffect != "unknown" {
					t.Error("completed child concealed lost approved integration proof", status, err, se)
				}
			} else if err != nil || se != nil || status.ControlEffect != "paused" || f.s.ActiveTaskID != "" {
				t.Error("completed native integration did not park between Goals", status, err, se)
			}
			if *f.m.Native.(mr05NativeQueue).effects != 0 {
				t.Error("completion launched another Goal")
			}
		})
	}
}

func mr05ParentTerminalConsumption(t *testing.T) {
	f := newMR05Week(t)
	r := weekPreview(t, f, "pause", "terminal-parent-0001")
	weekApply(t, f, r)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for i := 0; i < 2; i++ {
		_ = f.m.Serve(ctx, f.c.Workspace, f.b.ID)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := f.m.Serve(ctx, f.c.Workspace, f.b.ID)
	if err == nil || ctx.Err() != nil {
		t.Error("terminal parent retried instead of exiting promptly", time.Since(start), err)
	}
	f.reload(t)
	generation := f.s.Generation
	revision := f.s.ControlRevision
	for i := 0; i < 3; i++ {
		status, e := f.m.Controls().Status(f.c.Workspace, r.RequestID)
		if e == nil || status.RequestState != "outcome_unknown" || status.ControlEffect != "unknown" || status.Receipt != nil {
			t.Error("terminal failure still pending", status, e)
		}
		_, _ = f.m.Controls().Apply(r)
	}
	f.reload(t)
	if f.s.Generation != generation || f.s.ControlRevision != revision || len(f.n.starts) != 0 {
		t.Error("replay mutated terminal owner")
	}
}

func mr05StateBeforeGateGap(t *testing.T) {
	t.Run("check_is_read_only", func(t *testing.T) {
		f := newMR05Week(t)
		before := mr05ChildInventory(t, f.c.Workspace)
		if _, e := f.m.Start(f.c, true); e != nil {
			t.Fatal(e)
		}
		if taskrun.ControlRevision(before) != taskrun.ControlRevision(mr05ChildInventory(t, f.c.Workspace)) {
			t.Error("existing registration check wrote files")
		}
	})
	for _, scenario := range []string{"initial", "activity", "retained_control", "changed_gate", "changed_gate_after_activity", "changed_resume_gate"} {
		t.Run(scenario, func(t *testing.T) {
			f := newMR05Week(t)
			gatePath := taskrun.ParentControlPath(f.c.Workspace, f.b.ID)
			switch scenario {
			case "activity", "changed_gate_after_activity", "changed_resume_gate":
				f.tick(t)
			case "retained_control":
				weekApply(t, f, weekPreview(t, f, "pause", "gate-gap-pause-0001"))
			}
			changed := strings.HasPrefix(scenario, "changed_")
			if changed {
				g, e := taskrun.ReadParentControl(f.c.Workspace, f.b.ControlContract)
				if e != nil {
					t.Fatal(e)
				}
				if scenario == "changed_resume_gate" {
					g.ResumeAllowed = true
				} else {
					g.DesiredAction = "stop"
				}
				if e := publish(gatePath, g, true); e != nil {
					t.Fatal(e)
				}
			} else if e := os.Remove(gatePath); e != nil {
				t.Fatal(e)
			}
			stateBefore, e := os.ReadFile(filepath.Join(f.dir, "state.json"))
			if e != nil {
				t.Fatal(e)
			}
			*f.now = f.now.Add(time.Minute)
			s, err := f.m.Start(f.c, false)
			if scenario == "initial" {
				if err != nil || s.StartedAt != f.b.StartedAt || s.Deadline != f.b.Deadline {
					t.Fatal("original registration not recovered", s, err)
				}
				if _, e = taskrun.ReadParentControl(f.c.Workspace, f.b.ControlContract); e != nil {
					t.Error("state-before-gate gap was left unusable", e)
				}
			} else {
				if err == nil {
					t.Error("inconsistent retained authority accepted", scenario)
				}
				if !changed {
					if _, e = os.Stat(gatePath); !os.IsNotExist(e) {
						t.Error("reconstructed authority after activity", e)
					}
				}
			}
			stateAfter, e := os.ReadFile(filepath.Join(f.dir, "state.json"))
			if e != nil || string(stateBefore) != string(stateAfter) {
				t.Error("registration readback rewrote existing state", e)
			}
		})
	}
}
