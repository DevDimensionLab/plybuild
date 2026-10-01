package taskrun

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReturnRepairLegacyNegativeCacheKeepsClaimProvenance(t *testing.T) {
	d, r, f := fixture(t)
	runner := d.Runner.(*fakeRunner)
	runner.inside = func(LaunchSpec) error {
		a := acceptance(t, d, r)
		a.Acceptance = "unknown"
		a.Issues = json.RawMessage(`[{"type":"unknown","detail":"Synthetic old negative receipt."}]`)
		if _, e := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "claim.json", a)); e != nil {
			return e
		}
		j, e := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
		if e != nil {
			return e
		}
		// The previous producer collapsed negative event states to rejected.
		last := len(j.Events) - 1
		ev := j.Events[last]
		var payload map[string]any
		json.Unmarshal(ev.Payload, &payload)
		payload["state"] = "rejected"
		ev.Payload, _ = Canonical(payload)
		raw, _ := Canonical(ev)
		if e = os.WriteFile(filepath.Join(runPaths(r).RunRoot, "events", fmt.Sprintf("%06d.json", ev.Sequence)), raw, 0600); e != nil {
			return e
		}
		prefix := journal{Request: r, Binding: j.Binding, Result: initial(r)}
		for i, event := range j.Events {
			if e = fold(&prefix, event); e != nil {
				return e
			}
			prefix.Result.LastEventSHA256 = &j.Hashes[i]
		}
		j.Result = prefix.Result
		j.Result.Acceptance.State = "rejected"
		j.Result.LastEventSHA256 = ptr(hash(raw))
		if e = replaceValue(filepath.Join(runPaths(r).RunRoot, "result.json"), resultCacheValue(j.Result, 1)); e != nil {
			return e
		}
		before := treeState(t, r.WorkspaceRoot)
		out, e := Show(d, r.WorkspaceRoot, RunID(r.RequestKey))
		if e != nil {
			return e
		}
		if out.Acceptance.State != "unknown" || out.RuntimeFacts.Source == nil {
			return errors.New("legacy negative claim provenance lost")
		}
		if before != treeState(t, r.WorkspaceRoot) {
			return errors.New("legacy cache rewritten")
		}
		return nil
	}
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	Start(d, f, *p.(Preview).Confirmation)
	if runner.lastErr != nil {
		t.Fatal(runner.lastErr)
	}
}

func TestReturnRepairNegativeClaimPreservedAfterDraftBindingDrift(t *testing.T) {
	d, r, f := fixture(t)
	runner := d.Runner.(*fakeRunner)
	runner.inside = func(LaunchSpec) error {
		a := acceptance(t, d, r)
		a.SchemaVersion = 2
		a.Acceptance = "conflict"
		a.Issues = json.RawMessage(`[{"type":"target","detail":"Synthetic target changed during interrupted publication."}]`)
		j, e := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
		if e != nil {
			return e
		}
		draft, e := workflowhandoff.BuildTaskRunStart(d.Workflow, j.Binding.Handoff.Locator, r.HumanAuthority.ActorClaim, r.HumanAuthority.StartSurface, a.SessionID, *a.RuntimeClaim.RuntimeID, *a.RuntimeClaim.ModelID, a.Sandbox, a.Issues, a.Acceptance)
		if e != nil {
			return e
		}
		if e = writeOnce(acceptanceDraftPath(r, digest(a)), draft); e != nil {
			return e
		}
		if e = os.WriteFile(filepath.Join(j.Binding.Preparation.Plan.WorktreePath, "synthetic-drift.txt"), []byte("synthetic target drift"), 0644); e != nil {
			return e
		}
		out, e := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "claim.json", a))
		if e != nil {
			return e
		}
		if out.Acceptance.State != "conflict" || out.Acceptance.ReceiptSHA256 != nil {
			return errors.New("unrepresentable negative claim lost or fabricated receipt")
		}
		return nil
	}
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	Start(d, f, *p.(Preview).Confirmation)
	if runner.lastErr != nil {
		t.Fatal(runner.lastErr)
	}
}

func TestReturnRepairLateClaimNativeMismatchRejectedBeforePublication(t *testing.T) {
	d, r, f := fixture(t)
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	Start(d, f, *p.(Preview).Confirmation)
	s := syntheticStatus(t, d, r, "inactive")
	s.NativeSessionID = ptr("observed-session")
	if _, e = collectSyntheticStatus(t, d, r, s); e != nil {
		t.Fatal(e)
	}
	a := acceptance(t, d, r)
	a.RuntimeClaim.NativeSessionID = ptr("different-session")
	path := writeAny(t, r.WorkspaceRoot, "claim.json", a)
	before := treeState(t, storeRoot(r.WorkspaceRoot))
	if _, e = Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), path); e == nil {
		t.Fatal("mismatched native identity accepted")
	}
	if before != treeState(t, storeRoot(r.WorkspaceRoot)) {
		t.Fatal("mismatched identity published claim before rejection")
	}
}

func TestReturnRepairCallbackTransportRejectsDriftBeforeWrite(t *testing.T) {
	d, r, f := fixture(t)
	runner := d.Runner.(*fakeRunner)
	runner.inside = func(s LaunchSpec) error {
		path := writeAny(t, r.WorkspaceRoot, "claim.json", acceptance(t, d, r))
		for _, mode := range []string{"empty", "missing", "conflicting-env", "wrong-path", "symlink", "cwd", "executable", "context-session"} {
			bad := d
			explicit := s.ContextPath
			bad.CallbackContext = &explicit
			bad.ContextPath = func() string { return "" }
			contextBytes, _ := os.ReadFile(s.ContextPath)
			switch mode {
			case "empty":
				explicit = ""
			case "missing":
				bad.CallbackContext = nil
			case "conflicting-env":
				bad.ContextPath = func() string { return "/different" }
			case "wrong-path":
				explicit = path
			case "symlink":
				explicit = filepath.Join(r.WorkspaceRoot, "context-link")
				if e := os.Symlink(s.ContextPath, explicit); e != nil {
					return e
				}
			case "cwd":
				bad.CWD = func() (string, error) { return r.WorkspaceRoot, nil }
			case "executable":
				bad.Executable = func() (string, error) { return r.Runtime.Executable.Path, nil }
			case "context-session":
				var ctx Context
				json.Unmarshal(contextBytes, &ctx)
				ctx.SessionID = "different"
				writeAny(t, filepath.Dir(s.ContextPath), filepath.Base(s.ContextPath), ctx)
			}
			before := treeState(t, r.WorkspaceRoot)
			if _, e := Accept(bad, r.WorkspaceRoot, RunID(r.RequestKey), path); e == nil {
				return errors.New("accepted " + mode)
			}
			if before != treeState(t, r.WorkspaceRoot) {
				return errors.New("callback wrote before rejecting " + mode)
			}
			if mode == "context-session" {
				if e := os.WriteFile(s.ContextPath, contextBytes, 0600); e != nil {
					return e
				}
			}
		}
		// Old env-only transport remains supported, with the same binding checks.
		_, e := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), path)
		return e
	}
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	Start(d, f, *p.(Preview).Confirmation)
	if runner.lastErr != nil {
		t.Fatal(runner.lastErr)
	}
}

func TestReturnRepairNegativeClaimDriftAndRecovery(t *testing.T) {
	for _, point := range []string{"after_acceptance_claim", "after_submit_start", "after_event_acceptance_received"} {
		t.Run(point, func(t *testing.T) {
			d, r, f := fixture(t)
			runner := d.Runner.(*fakeRunner)
			var path string
			runner.inside = func(LaunchSpec) error {
				a := acceptance(t, d, r)
				a.SchemaVersion = 2
				a.RuntimeClaim.ModelID = ptr("different-model")
				path = writeAny(t, r.WorkspaceRoot, "claim.json", a)
				if _, e := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), path); e == nil {
					return errors.New("positive mismatch accepted")
				}
				a.Acceptance = "conflict"
				a.Issues = json.RawMessage(`[{"type":"principal","detail":"Synthetic runtime and policy drift."}]`)
				if point != "after_submit_start" {
					a.RuntimeClaim.RuntimeID = nil
					a.RuntimeClaim.ProfileID = nil
					a.RuntimeClaim.EffectivePolicySHA256 = nil
					a.Sandbox = json.RawMessage(`null`)
				}
				// Runtime evidence drift must not prevent preserving a bound rejection.
				if e := os.WriteFile(r.Runtime.PermissionBinding.Evidence[0].Locator, []byte("synthetic changed policy"), 0600); e != nil {
					return e
				}
				path = writeAny(t, r.WorkspaceRoot, "claim.json", a)
				d.Fault = func(p string) error {
					if p == point {
						return errors.New("synthetic acceptance crash")
					}
					return nil
				}
				if _, e := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), path); e == nil {
					return errors.New("fault not fired")
				}
				d.Fault = nil
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
			if out.Acceptance.State != "conflict" || out.RuntimeFacts.ModelState != "mismatch" {
				t.Fatalf("negative facts lost: %+v", out)
			}
			if point != "after_submit_start" && out.Acceptance.ReceiptSHA256 != nil {
				t.Fatal("fabricated receipt")
			}
			before := treeState(t, r.WorkspaceRoot)
			if _, e = Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), path); e != nil {
				t.Fatal(e)
			}
			if before != treeState(t, r.WorkspaceRoot) {
				t.Fatal("claim retry wrote")
			}
			var a Acceptance
			raw, _ := os.ReadFile(path)
			json.Unmarshal(raw, &a)
			a.Acceptance = "unknown"
			changed := writeAny(t, r.WorkspaceRoot, "different-claim.json", a)
			if _, e = Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), changed); e == nil {
				t.Fatal("different claim replaced first")
			}
			out, e = collectSyntheticStatus(t, d, r, syntheticStatus(t, d, r, "inactive"))
			if e != nil || out.Collection.State == "qualified" {
				t.Fatal("negative claim qualified")
			}
		})
	}
}

// Build a synthetic old-contract publication through the unchanged WF/registry
// validators, without invoking the new collect gate. This models existing disk
// history; production has no bypass for publishing it today.
func legacyResultFixture(t *testing.T, d Dependencies, r Request, registered bool) []byte {
	t.Helper()
	j, e := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
	if e != nil {
		t.Fatal(e)
	}
	raw, e := readFile(filepath.Join(runPaths(r).RunRoot, "reports", "accepted.json"), 4<<20, true)
	if e != nil {
		t.Fatal(e)
	}
	report, e := validateReport(raw)
	if e != nil {
		t.Fatal(e)
	}
	technical, _ := Canonical(report.TechnicalAssessment)
	draft, _, e := workflowhandoff.BuildTaskRunResult(d.Workflow, j.Binding.Handoff.Locator, RunID(r.RequestKey), technical)
	if e != nil {
		t.Fatal(e)
	}
	path := filepath.Join(runPaths(r).TempRoot, "task-result-draft.json")
	if e = writeOnce(path, draft); e != nil {
		t.Fatal(e)
	}
	if registered {
		w := d.Workspace
		w.HandoffEvidence = workflowhandoff.NewTaskHandoffEvidenceReader(d.Workflow)
		record, e := workspace.RecordTaskResult(w, workspace.TaskResultRecordInput{TaskID: j.Binding.Preparation.Plan.TaskID, File: path})
		if e != nil {
			t.Fatal(e)
		}
		if e = appendEvent(d, r, "collection", map[string]any{"collection": historicalCollection(&record.Record), "observed_target": historicalTarget(&record.Record), "reasons": []Reason{}}); e != nil {
			t.Fatal(e)
		}
	}
	j, e = readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
	if e != nil {
		t.Fatal(e)
	}
	if e = replaceValue(filepath.Join(runPaths(r).RunRoot, "result.json"), resultCacheValue(j.Result, 1)); e != nil {
		t.Fatal(e)
	}
	return draft
}

func TestReturnRepairLegacyCachePublicationAndReservation(t *testing.T) {
	for _, mode := range []string{"qualified", "registry-crash", "unregistered-draft"} {
		t.Run(mode, func(t *testing.T) {
			d, r, f := fixture(t)
			completeInside(t, &d, r, nil)
			p, e := PreviewStart(d, f)
			if e != nil {
				t.Fatal(e)
			}
			Start(d, f, *p.(Preview).Confirmation)
			draft := legacyResultFixture(t, d, r, mode != "unregistered-draft")
			if mode == "registry-crash" {
				j, e := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
				if e != nil {
					t.Fatal(e)
				}
				last := len(j.Events)
				if e = os.Remove(filepath.Join(runPaths(r).RunRoot, "events", fmt.Sprintf("%06d.json", last))); e != nil {
					t.Fatal(e)
				}
				if e = os.Remove(filepath.Join(runPaths(r).RunRoot, "result.json")); e != nil {
					t.Fatal(e)
				}
			}
			before := treeState(t, r.WorkspaceRoot)
			out, e := Show(d, r.WorkspaceRoot, RunID(r.RequestKey))
			if e != nil {
				t.Fatal(e)
			}
			preview, e := PreviewCollect(d, r.WorkspaceRoot, RunID(r.RequestKey))
			if e != nil {
				t.Fatal(e)
			}
			if before != treeState(t, r.WorkspaceRoot) {
				t.Fatal("historical readback wrote")
			}
			if out.TaskExecution.State != "unknown" || out.TaskExecution.HistoricalQualification != (mode != "unregistered-draft") {
				t.Fatal("historical task fact invented")
			}
			out, e = Collect(d, r.WorkspaceRoot, RunID(r.RequestKey), *preview.Confirmation)
			if e != nil {
				t.Fatal(e)
			}
			if (out.Collection.State == "qualified") != (mode != "unregistered-draft") {
				t.Fatal("legacy draft bypassed today's gate")
			}
			next := r
			next.RequestKey = "fixture/new"
			nextFile := writeAny(t, r.WorkspaceRoot, "next.json", next)
			if _, e = PreviewStart(d, nextFile); e == nil {
				t.Fatal("historical qualification released reservation")
			}
			j, e := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
			if e != nil {
				t.Fatal(e)
			}
			// Both durable request inventory and target cache independently retain it.
			var reserved struct {
				Target     workspace.PlanWorktreeObservation `json:"target"`
				RequestKey string                            `json:"request_key"`
			}
			if e = decode(j.Events[0].Payload, 64<<10, &reserved); e != nil {
				t.Fatal(e)
			}
			index := targetPath(r, reserved.Target)
			bytes, e := os.ReadFile(index)
			if e != nil {
				t.Fatal(e)
			}
			if e = os.Remove(index); e != nil {
				t.Fatal(e)
			}
			if _, e = PreviewStart(d, nextFile); e == nil {
				t.Fatal("request inventory failed to reserve")
			}
			if e = os.WriteFile(index, bytes, 0600); e != nil {
				t.Fatal(e)
			}
			reqIndex := filepath.Join(storeRoot(r.WorkspaceRoot), "requests", strings.TrimPrefix(hash([]byte(r.RequestKey)), "sha256:")+".json")
			reqBytes, e := os.ReadFile(reqIndex)
			if e != nil {
				t.Fatal(e)
			}
			if e = os.Remove(reqIndex); e != nil {
				t.Fatal(e)
			}
			if _, e = PreviewStart(d, nextFile); e == nil {
				t.Fatal("target inventory failed to reserve")
			}
			if e = os.WriteFile(reqIndex, reqBytes, 0600); e != nil {
				t.Fatal(e)
			}
			out, e = collectSyntheticStatus(t, d, r, syntheticStatus(t, d, r, "inactive"))
			if e != nil {
				t.Fatal(e)
			}
			if out.Collection.State != "qualified" {
				t.Fatalf("inactive recovery: %+v", out)
			}
			current, e := os.ReadFile(filepath.Join(runPaths(r).TempRoot, "task-result-draft.json"))
			if e != nil || string(current) != string(draft) {
				t.Fatal("historical draft changed")
			}
			reg, e := d.Workspace.WorkItems.Snapshot(r.WorkspaceRoot)
			if e != nil || len(reg.TaskResults) != 1 {
				t.Fatal("duplicate historical result")
			}
			if _, e = PreviewStart(d, nextFile); e != nil {
				t.Fatalf("inactive reservation not released: %v", e)
			}
		})
	}
}
