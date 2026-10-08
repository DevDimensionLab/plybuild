package taskrun

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// A provider update can remove the original launcher while its already
// accepted process and conversation remain alive in Herdr. The replacement on
// PATH is deliberately different and must never become the old run's identity.
func providerUpgradeFixture(t *testing.T, f deliveryFixture, change string) string {
	t.Helper()
	probe := filepath.Join(f.R.WorkspaceRoot, "unexpected-provider-starts")
	replacement := []byte("#!/bin/sh\nprintf 'unexpected start\\n' >> " + ShellQuote(probe) + "\nexit 0\n")
	switch change {
	case "removed":
		if err := os.Remove(f.R.Runtime.Executable.Path); err != nil {
			t.Fatal(err)
		}
	case "relocated":
		if err := os.Rename(f.R.Runtime.Executable.Path, f.R.Runtime.Executable.Path+".old"); err != nil {
			t.Fatal(err)
		}
	case "changed":
		if err := os.WriteFile(f.R.Runtime.Executable.Path, replacement, 0700); err != nil {
			t.Fatal(err)
		}
	default:
		t.Fatalf("unknown provider upgrade fixture %q", change)
	}
	installed := filepath.Join(f.R.WorkspaceRoot, "new-installed-provider")
	if err := os.WriteFile(installed, replacement, 0700); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(f.R.WorkspaceRoot, "bin", f.R.Runtime.Provider)
	if err := os.Remove(alias); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(installed, alias); err != nil {
		t.Fatal(err)
	}
	return probe
}

func providerUpgradeNoRestart(t *testing.T, f deliveryFixture, probe string) {
	t.Helper()
	if n := workflowTestCalls(t, f.workflowFixture, "agent start"); n != 1 {
		t.Fatalf("accepted delivery started %d provider sessions, want 1", n)
	}
	if n := workflowTestCalls(t, f.workflowFixture, "agent prompt"); n != 1 {
		t.Fatalf("accepted delivery sent %d Task prompts, want 1", n)
	}
	if _, err := os.Stat(probe); !os.IsNotExist(err) {
		t.Fatalf("provider update executed a historical or replacement launcher: %v", err)
	}
}

func TestDeliveryProviderUpgradePreservesAcceptedContinuation(t *testing.T) {
	for _, provider := range []string{"codex", "claude"} {
		for _, change := range []string{"removed", "relocated", "changed"} {
			t.Run(provider+"/"+change, func(t *testing.T) {
				f := newDeliveryFixture(t, provider)
				o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
				d, in, _ := continuationTestInput(t, f, o)
				frozen := map[string][]byte{}
				for _, path := range []string{workflowIndex(f.R.WorkspaceRoot, o.RunID), o.Paths.Context, f.R.Runtime.PlyExecutable.Path, filepath.Join(o.Paths.RunRoot, "acceptance.json"), filepath.Join(o.Paths.RunRoot, "mandate.json")} {
					raw, err := os.ReadFile(path)
					if err != nil {
						t.Fatal(err)
					}
					frozen[path] = raw
				}
				probe := providerUpgradeFixture(t, f, change)
				p, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in)
				if err != nil || p.State != "ready" {
					t.Fatalf("historical launcher blocked accepted continuation: %+v %v", p, err)
				}
				continued, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation)
				if err != nil || continued.State != "continued" {
					t.Fatalf("accepted continuation failed after provider update: %+v %v", continued, err)
				}
				d.Executable = func() (string, error) { return continued.ControlExecutable.Path, nil }
				r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "after-provider-upgrade", Phase: "working", Summary: "The accepted session continues after launcher cleanup.", Meaning: "Isolated fixture progress preserves the same native owner and authority.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
				o, err = WorkflowDeliveryReport(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "progress-after-upgrade.json", r))
				if err != nil || o.Delivery.Phase != "working" {
					t.Fatalf("same owner callback rejected after provider cleanup: %+v %v", o, err)
				}
				candidate := deliveryTestCandidate(t, f, "accepted provider continuity")
				o, err = WorkflowDeliveryVerify(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "after-provider-upgrade"))
				if err != nil || len(o.Delivery.Candidates) != 1 || o.Delivery.Candidates[0].OID != candidate || o.Delivery.Phase != "awaiting_human_qa" {
					t.Fatalf("accepted verification blocked after launcher cleanup: %+v %v", o.Delivery, err)
				}
				s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
				if err != nil {
					t.Fatal(err)
				}
				runtime, err := workflowEffectiveRuntime(s)
				if err != nil || runtime.Executable != f.R.Runtime.Executable || s.Result.SessionID != o.SessionID || s.Result.Transport.AgentSessionID != "fixture-session" {
					t.Fatalf("provider history or accepted owner was substituted: %+v %v", runtime, err)
				}
				for path, want := range frozen {
					got, err := os.ReadFile(path)
					if err != nil || !bytes.Equal(got, want) {
						t.Fatalf("provider update rewrote frozen artifact %s: %v", path, err)
					}
				}
				providerUpgradeNoRestart(t, f, probe)
			})
		}
	}
}

func TestDeliveryProviderUpgradeDoesNotGrantInitialAcceptance(t *testing.T) {
	for _, change := range []string{"removed", "relocated", "changed"} {
		t.Run(change, func(t *testing.T) {
			f := newDeliveryFixture(t, "codex")
			o := workflowTestStart(t, f.workflowFixture)
			a := deliveryTestAcceptance(t, f, o)
			probe := providerUpgradeFixture(t, f, change)
			if _, err := WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "initial-acceptance-after-upgrade.json", a)); err == nil {
				t.Fatal("initial runtime acceptance ignored its launch executable binding")
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil || s.Acceptance != nil || s.StartSHA256 != nil || s.Result.Delivery.PermissionState != "pending_runtime_acceptance" {
				t.Fatalf("provider history substituted before positive acceptance: %+v %v", s.Result.Delivery, err)
			}
			providerUpgradeNoRestart(t, f, probe)
		})
	}
}

func TestDeliveryProviderUpgradeDoesNotGrantNewLaunch(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	preview, err := WorkflowPreviewStart(f.D, f.File)
	if err != nil {
		t.Fatal(err)
	}
	probe := providerUpgradeFixture(t, f, "removed")
	if _, err := WorkflowStart(f.D, f.File, *preview.(WorkflowPreview).Confirmation); err == nil {
		t.Fatal("new launch substituted the installed provider for its removed bound launcher")
	}
	if n := workflowTestCalls(t, f.workflowFixture, "agent start"); n != 0 {
		t.Fatalf("missing bound launcher started %d providers", n)
	}
	if _, err := os.Stat(probe); !os.IsNotExist(err) {
		t.Fatalf("new launch executed the unbound installed provider: %v", err)
	}
}

func TestDeliveryProviderUpgradeRequiresCurrentRuntimeObservation(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	d, in, _ := continuationTestInput(t, f, o)
	probe := providerUpgradeFixture(t, f, "removed")
	r := providerUpgradeIncomplete(t, f, o, "stopped")
	r.Phase, r.EventID = "working", "unobserved-working"
	if _, err := WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "unobserved-working.json", r)); err == nil {
		t.Fatal("ordinary callback granted continuation without current runtime evidence")
	}
	missing := in
	missing.RuntimeEvidencePath = ""
	if _, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, missing); err == nil || !strings.Contains(err.Error(), "runtime_observation_required") {
		t.Fatalf("historical launcher absence granted unobserved runtime authority: %v", err)
	}
	var original DeliveryRuntimeObservation
	if err := readValue(in.RuntimeEvidencePath, 1<<20, &original); err != nil {
		t.Fatal(err)
	}
	for _, mutation := range []string{"unconfirmed", "reduced_policy", "wrong_session", "wrong_provider", "wrong_context", "wrong_acceptance", "stale", "future"} {
		t.Run(mutation, func(t *testing.T) {
			observation := original
			switch mutation {
			case "unconfirmed":
				observation.DeliveryPermission.PermissionConfirmed = false
			case "reduced_policy":
				path := filepath.Join(f.R.WorkspaceRoot, "reduced-current-policy.txt")
				if err := os.WriteFile(path, []byte("The current fixture policy no longer covers Task effects.\n"), 0600); err != nil {
					t.Fatal(err)
				}
				observation.RuntimeClaim.EffectivePolicySHA256 = ptr(hashFileTest(t, path))
				observation.DeliveryPermission.ActualPolicyEvidence = []Evidence{{Locator: path, SHA256: hashFileTest(t, path), Role: "effective_policy"}}
			case "wrong_session":
				observation.RuntimeClaim.NativeSessionID = ptr("another-native-session")
			case "wrong_provider":
				observation.RuntimeClaim.RuntimeID = ptr("claude")
			case "wrong_context":
				observation.Context.SHA256 = hash([]byte("another context"))
			case "wrong_acceptance":
				observation.Acceptance.SHA256 = hash([]byte("another acceptance"))
			case "stale":
				observation.ObservedAtUTC = d.Now().Add(-11 * time.Minute).UTC().Format(time.RFC3339Nano)
			case "future":
				observation.ObservedAtUTC = d.Now().Add(time.Minute).UTC().Format(time.RFC3339Nano)
			}
			bad := in
			bad.RuntimeEvidencePath = writeAny(t, f.R.WorkspaceRoot, "runtime-observation-"+mutation+".json", observation)
			if _, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, bad); err == nil {
				t.Fatalf("invalid current runtime observation granted continuation: %s", mutation)
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil || s.Continuation != nil || len(s.Result.Delivery.Events) != 0 {
				t.Fatalf("invalid observation changed accepted delivery: %+v %v", s.Result.Delivery, err)
			}
		})
	}
	if got, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in); err != nil || got.State != "ready" {
		t.Fatalf("unchanged current runtime could not continue: %+v %v", got, err)
	}
	providerUpgradeNoRestart(t, f, probe)
}

func TestDeliveryProviderUpgradeRetainsAuthorityAndLiveOwnerGuards(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	previewD, in, _ := continuationTestInput(t, f, o)
	probe := providerUpgradeFixture(t, f, "removed")
	p, err := WorkflowPreviewDeliveryContinuation(previewD, f.R.WorkspaceRoot, in)
	if err != nil {
		t.Fatal(err)
	}
	continued, err := WorkflowContinueDelivery(previewD, f.R.WorkspaceRoot, in, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	d := f.D
	d.Executable = func() (string, error) { return continued.ControlExecutable.Path, nil }
	r := providerUpgradeIncomplete(t, f, o, "stopped")
	r.Phase, r.EventID = "working", "guarded-working"
	file := writeAny(t, f.R.WorkspaceRoot, "guarded-working.json", r)
	for name, path := range map[string]string{
		"context":          o.Paths.Context,
		"request":          workflowIndex(f.R.WorkspaceRoot, o.RunID),
		"original_control": f.R.Runtime.PlyExecutable.Path,
		"current_control":  continued.ControlExecutable.Path,
		"acceptance":       filepath.Join(o.Paths.RunRoot, "acceptance.json"),
		"goal":             f.R.Delivery.Goal.Locator,
		"launch_policy":    f.R.Runtime.PermissionBinding.Evidence[0].Locator,
		"actual_policy":    filepath.Join(f.R.WorkspaceRoot, "actual-policy.json"),
		"current_runtime":  in.RuntimeEvidencePath,
		"herdr":            f.R.Herdr.Executable.Path,
	} {
		t.Run(name, func(t *testing.T) {
			providerUpgradeCorruptFile(t, path, false)
			if _, err := WorkflowPreviewDeliveryContinuation(previewD, f.R.WorkspaceRoot, in); err == nil {
				t.Fatalf("continuation ignored changed %s after provider cleanup", name)
			}
			if _, err := WorkflowDeliveryReport(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); err == nil {
				t.Fatalf("continued callback ignored changed %s", name)
			}
		})
	}
	for _, mutation := range []string{"missing_session", "wrong_session", "missing_foreground", "wrong_foreground", "missing_status", "unknown_status", "pending_launch"} {
		t.Run(mutation, func(t *testing.T) {
			changes := map[string]any{}
			switch mutation {
			case "missing_session":
				changes["observations"] = []any{map[string]any{"drop_fields": []string{"agent_session"}}}
			case "wrong_session":
				changes["agent_session_id"] = "another-native-session"
			case "missing_foreground":
				changes["observations"] = []any{map[string]any{"drop_fields": []string{"foreground_cwd"}}}
			case "wrong_foreground":
				changes["foreground_cwd"] = f.Parent
			case "missing_status":
				changes["observations"] = []any{map[string]any{"drop_fields": []string{"agent_status"}}}
			case "unknown_status":
				changes["agent_status"] = "unknown"
			case "pending_launch":
				changes["observations"] = []any{map[string]any{"launch_pending": true}}
			}
			providerUpgradeChangeModel(t, f, changes)
			if _, err := WorkflowPreviewDeliveryContinuation(previewD, f.R.WorkspaceRoot, in); err == nil {
				t.Fatalf("continuation granted without a known live owner: %s", mutation)
			}
		})
	}
	if got, err := WorkflowDeliveryReport(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); err != nil || got.Delivery.Phase != "working" {
		t.Fatalf("guard fixture has another hidden callback blocker: %+v %v", got, err)
	}
	providerUpgradeNoRestart(t, f, probe)
}

func providerUpgradeCorruptFile(t *testing.T, path string, remove bool) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chmod(path, info.Mode().Perm()|0200); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, raw, info.Mode().Perm()|0200); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(path, info.Mode().Perm()); err != nil {
			t.Fatal(err)
		}
	})
	if remove {
		err = os.Remove(path)
	} else {
		if err = os.Chmod(path, info.Mode().Perm()|0200); err == nil {
			err = os.WriteFile(path, append(append([]byte{}, raw...), []byte("tampered\n")...), info.Mode().Perm()|0200)
		}
	}
	if err != nil {
		t.Fatal(err)
	}
}

func providerUpgradeChangeModel(t *testing.T, f deliveryFixture, changes map[string]any) {
	t.Helper()
	raw, err := os.ReadFile(f.Model)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.WriteFile(f.Model, raw, 0600); err != nil {
			t.Fatal(err)
		}
	})
	workflowTestModel(t, f.workflowFixture, changes)
}

func providerUpgradeIncomplete(t *testing.T, f deliveryFixture, o WorkflowRun, phase string) DeliveryReport {
	t.Helper()
	r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "runtime-dependency-" + phase, PreviousEventSHA256: o.Delivery.LastEventSHA256, Phase: phase, Summary: "The current runtime dependency cannot be confirmed.", Meaning: "The isolated owner preserves an incomplete outcome without claiming Task effects, verification, or human approval.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{{ID: "acceptance", Outcome: "not_run", Reason: "Runtime authority remains unresolved.", Argv: []string{}, Evidence: []FileBinding{}}}}
	if phase == "needs_input" {
		r.Question = ptr("Which actual policy currently applies to this same provider session?")
	}
	return r
}

func TestDeliveryIncompleteReportWithoutRuntimeAuthority(t *testing.T) {
	for _, accepted := range []bool{false, true} {
		for _, phase := range []string{"stopped", "needs_input"} {
			name := "before_acceptance/" + phase
			if accepted {
				name = "after_acceptance/" + phase
			}
			t.Run(name, func(t *testing.T) {
				f := newDeliveryFixture(t, "codex")
				o := workflowTestStart(t, f.workflowFixture)
				if accepted {
					o = deliveryTestAccept(t, f, o)
					if err := os.Remove(filepath.Join(f.R.WorkspaceRoot, "actual-policy.json")); err != nil {
						t.Fatal(err)
					}
				}
				before, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
				if err != nil {
					t.Fatal(err)
				}
				probe := providerUpgradeFixture(t, f, "removed")
				if err := os.Remove(f.R.Runtime.PermissionBinding.Evidence[0].Locator); err != nil {
					t.Fatal(err)
				}
				d := f.D
				installed := filepath.Join(f.R.WorkspaceRoot, "installed-report-reader")
				if err := os.WriteFile(installed, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
					t.Fatal(err)
				}
				d.Executable = func() (string, error) { return installed, nil }
				r := providerUpgradeIncomplete(t, f, o, phase)
				file := writeAny(t, f.R.WorkspaceRoot, "runtime-dependency-report.json", r)
				got, err := WorkflowDeliveryIncompleteReport(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file)
				if err != nil || got.Delivery.Phase != phase || got.NextAction.Actor != "user" || len(got.Delivery.Events) != 1 || len(got.Delivery.Candidates) != 0 {
					t.Fatalf("truthful incomplete report was blocked or qualified: delivery=%+v next=%+v err=%v", got.Delivery, got.NextAction, err)
				}
				wantAction := r.Meaning
				if r.Question != nil {
					wantAction = *r.Question
				}
				if got.NextAction.Message != wantAction {
					t.Fatalf("incomplete readback replaced the owner's exact question or stopped meaning: got %q want %q", got.NextAction.Message, wantAction)
				}
				if _, err = WorkflowDeliveryIncompleteReport(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); err != nil {
					t.Fatalf("same incomplete report retry failed: %v", err)
				}
				after, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
				if err != nil || !equal(before.Acceptance, after.Acceptance) || !equal(before.StartSHA256, after.StartSHA256) || after.Result.Delivery.PermissionState != before.Result.Delivery.PermissionState || after.Result.Delivery.Attempt != nil || after.Continuation != nil || len(after.Result.Delivery.Events) != 1 {
					t.Fatalf("incomplete report acquired authority or reserved an effect: %+v %v", after.Result.Delivery, err)
				}
				var stored DeliveryReport
				if err := readValue(after.Result.Delivery.Events[0].Binding.Locator, 1<<20, &stored); err != nil || !equal(stored, r) || stored.VerifierResults[0].Exit != nil {
					t.Fatalf("incomplete report lost exact observation or invented a verifier exit: %+v %v", stored, err)
				}
				registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
				if err != nil || len(registry.TaskResults) != 0 || len(registry.HumanQARecords) != 0 {
					t.Fatalf("incomplete report fabricated a candidate or human verdict: %+v %v", registry, err)
				}
				providerUpgradeNoRestart(t, f, probe)
			})
		}
	}
}

func TestDeliveryIncompleteReportRetainsIdentityAndArtifactGuards(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	probe := providerUpgradeFixture(t, f, "removed")
	r := providerUpgradeIncomplete(t, f, o, "stopped")
	file := writeAny(t, f.R.WorkspaceRoot, "guarded-incomplete-report.json", r)
	stateBefore, err := os.ReadFile(filepath.Join(o.Paths.RunRoot, "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, mutation := range []string{"working", "fabricated_exit", "wrong_session", "wrong_request", "wrong_context_path", "wrong_cwd", "changed_context", "changed_request", "changed_control", "missing_herdr", "changed_herdr", "changed_acceptance", "changed_goal", "missing_session", "replaced_session", "replaced_provider", "replaced_pane", "replaced_owner", "missing_foreground", "pending_launch", "unknown_status"} {
		t.Run(mutation, func(t *testing.T) {
			d, context, input := f.D, o.Paths.Context, file
			bad := r
			switch mutation {
			case "working":
				bad.Phase = "working"
			case "fabricated_exit":
				bad.VerifierResults = append([]DeliveryVerifierResult{}, r.VerifierResults...)
				bad.VerifierResults[0].Exit = ptr(0)
			case "wrong_session":
				bad.SessionID = "another-owner"
			case "wrong_request":
				bad.RequestSHA256 = hash([]byte("another request"))
			case "wrong_context_path":
				context = filepath.Join(o.Paths.RunRoot, "different-context.json")
			case "wrong_cwd":
				d.CWD = func() (string, error) { return f.Parent, nil }
			case "changed_context":
				providerUpgradeCorruptFile(t, o.Paths.Context, false)
			case "changed_request":
				providerUpgradeCorruptFile(t, workflowIndex(f.R.WorkspaceRoot, o.RunID), false)
			case "changed_control":
				providerUpgradeCorruptFile(t, f.R.Runtime.PlyExecutable.Path, false)
			case "missing_herdr":
				providerUpgradeCorruptFile(t, f.R.Herdr.Executable.Path, true)
			case "changed_herdr":
				providerUpgradeCorruptFile(t, f.R.Herdr.Executable.Path, false)
			case "changed_acceptance":
				providerUpgradeCorruptFile(t, filepath.Join(o.Paths.RunRoot, "acceptance.json"), false)
			case "changed_goal":
				providerUpgradeCorruptFile(t, f.R.Delivery.Goal.Locator, false)
			case "missing_session":
				providerUpgradeChangeModel(t, f, map[string]any{"observations": []any{map[string]any{"drop_fields": []string{"agent_session"}}}})
			case "replaced_session":
				providerUpgradeChangeModel(t, f, map[string]any{"agent_session_id": "another-native-session"})
			case "replaced_provider":
				providerUpgradeChangeModel(t, f, map[string]any{"provider": "claude"})
			case "replaced_pane":
				providerUpgradeChangeModel(t, f, map[string]any{"observations": []any{map[string]any{"pane_id": "replacement-pane"}}})
			case "replaced_owner":
				providerUpgradeChangeModel(t, f, map[string]any{"name": "another-task-owner"})
			case "missing_foreground":
				providerUpgradeChangeModel(t, f, map[string]any{"observations": []any{map[string]any{"drop_fields": []string{"foreground_cwd"}}}})
			case "pending_launch":
				providerUpgradeChangeModel(t, f, map[string]any{"observations": []any{map[string]any{"launch_pending": true}}})
			case "unknown_status":
				providerUpgradeChangeModel(t, f, map[string]any{"agent_status": "unknown"})
			}
			if !equal(bad, r) {
				input = writeAny(t, f.R.WorkspaceRoot, "invalid-report-"+mutation+".json", bad)
			}
			if _, err := WorkflowDeliveryIncompleteReport(d, f.R.WorkspaceRoot, o.RunID, context, input); err == nil {
				t.Fatalf("incomplete reporting accepted %s", mutation)
			}
			state, err := os.ReadFile(filepath.Join(o.Paths.RunRoot, "state.json"))
			if err != nil || !bytes.Equal(stateBefore, state) {
				t.Fatalf("rejected %s changed delivery state: %v", mutation, err)
			}
		})
	}
	if got, err := WorkflowDeliveryIncompleteReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); err != nil || got.Delivery.Phase != "stopped" {
		t.Fatalf("guard fixture has another hidden blocker: %+v %v", got, err)
	}
	providerUpgradeNoRestart(t, f, probe)
}

func TestDeliveryIncompleteReportInterruptedAndConcurrentPublication(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	probe := providerUpgradeFixture(t, f, "removed")
	r := providerUpgradeIncomplete(t, f, o, "stopped")
	file := writeAny(t, f.R.WorkspaceRoot, "idempotent-incomplete-report.json", r)
	fault := f.D
	fault.Fault = func(point string) error {
		if point == "delivery_after_incomplete_report_event" {
			return errors.New("lost reply after immutable incomplete event publication")
		}
		return nil
	}
	if _, err := WorkflowDeliveryIncompleteReport(fault, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); err == nil {
		t.Fatal("fault did not interrupt incomplete report publication")
	}
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := WorkflowDeliveryIncompleteReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file)
			if err == nil && (got.Delivery.Phase != "stopped" || len(got.Delivery.Events) != 1) {
				err = errors.New("incomplete publication duplicated its event or lost its stopped phase")
			}
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	entries, err := os.ReadDir(filepath.Join(o.Paths.RunRoot, "delivery", "events"))
	if err != nil || len(entries) != 1 {
		t.Fatalf("interrupted retry created extra event artifacts: %v %v", entries, err)
	}
	r.Summary = "A conflicting report tries to reuse the same event identity."
	if _, err = WorkflowDeliveryIncompleteReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "conflicting-incomplete-report.json", r)); err == nil {
		t.Fatal("a changed incomplete report reused the published event ID")
	}
	providerUpgradeNoRestart(t, f, probe)
}

func TestDeliveryIncompleteReportPreservesUnresolvedVerification(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "unresolved verification")
	counter := filepath.Join(f.R.WorkspaceRoot, "acceptance-invocations.txt")
	if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nprintf 'run\\n' >> "+ShellQuote(counter)+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	fault := f.D
	fault.Fault = func(point string) error {
		if point == "delivery_after_verification_reservation" {
			return errors.New("verification interrupted after reservation")
		}
		return nil
	}
	interrupted, err := WorkflowDeliveryVerify(fault, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "unresolved"))
	if err == nil || interrupted.Delivery.Attempt == nil || interrupted.Delivery.Attempt.State != "attempted" {
		t.Fatalf("fixture did not preserve an unresolved verification: %+v %v", interrupted.Delivery, err)
	}
	probe := providerUpgradeFixture(t, f, "removed")
	r := providerUpgradeIncomplete(t, f, interrupted, "stopped")
	r.VerifierResults[0].Outcome = "unknown"
	r.VerifierResults[0].Reason = "The verifier reservation exists; no actual execution result is known."
	got, err := WorkflowDeliveryIncompleteReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "unresolved-verification-report.json", r))
	if err != nil || got.Delivery.Phase != "stopped" || !equal(got.Delivery.Attempt, interrupted.Delivery.Attempt) || len(got.Delivery.Candidates) != 0 {
		t.Fatalf("incomplete report replayed, cleared, or qualified unresolved verification: %+v %v", got.Delivery, err)
	}
	if _, err := os.Stat(counter); !os.IsNotExist(err) {
		t.Fatalf("reporting executed unresolved acceptance: %v", err)
	}
	providerUpgradeNoRestart(t, f, probe)
}
