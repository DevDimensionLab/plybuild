package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func workflowProviderFixture(t *testing.T, provider string) workflowFixture {
	t.Helper()
	f := workflowTestFixture(t)
	if err := os.Symlink(f.R.Runtime.Executable.Path, filepath.Join(f.R.WorkspaceRoot, "bin", "claude")); err != nil {
		t.Fatal(err)
	}
	f.R.Runtime.Provider = provider
	if provider == "claude" {
		f.R.Runtime.PermissionBinding.ProfileID = "manual"
	}
	workflowProviderRebind(t, &f)
	workflowTestModel(t, f, map[string]any{"provider": provider})
	return f
}

func workflowProviderRebind(t *testing.T, f *workflowFixture) {
	t.Helper()
	writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.File), f.R)
	id, request := workflowID(f.R), digest(f.R)
	f.Acceptance.RunID, f.Acceptance.RequestSHA256, f.Acceptance.SessionID = id, request, "ply:"+id
	f.Acceptance.RuntimeClaim.RuntimeID = ptr(f.R.Runtime.Provider)
	f.Acceptance.RuntimeClaim.ProfileID = ptr(f.R.Runtime.PermissionBinding.ProfileID)
	f.Acceptance.RuntimeClaim.EffectivePolicySHA256 = ptr(f.R.Runtime.PermissionBinding.EffectivePolicySHA256)
	f.Report.RunID, f.Report.RequestSHA256, f.Report.SessionID = id, request, "ply:"+id
}

func workflowProviderDocument(t *testing.T, file string) map[string]any {
	t.Helper()
	b, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	return doc
}

func workflowProviderReadback(t *testing.T, value any, expected string) {
	t.Helper()
	b, err := Canonical(value)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	if doc["provider"] != expected {
		t.Fatalf("provider readback = %v, want %s", doc["provider"], expected)
	}
}

func workflowProviderCalls(t *testing.T, f workflowFixture, command string) [][]string {
	t.Helper()
	b, err := os.ReadFile(f.Calls)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var calls [][]string
	for _, line := range bytes.Split(b, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		var call struct{ Argv []string }
		if err := json.Unmarshal(line, &call); err != nil {
			t.Fatal(err)
		}
		if len(call.Argv) >= 2 && strings.Join(call.Argv[:2], " ") == command {
			calls = append(calls, call.Argv)
		}
	}
	return calls
}

func workflowProviderNoReservation(t *testing.T, f workflowFixture) {
	t.Helper()
	if _, err := os.Stat(workflowIndex(f.R.WorkspaceRoot, workflowID(f.R))); !os.IsNotExist(err) {
		t.Fatalf("rejected request reserved a run: %v", err)
	}
	for _, command := range []string{"tab create", "agent start", "agent prompt"} {
		if n := workflowTestCalls(t, f, command); n != 0 {
			t.Fatalf("rejected request sent %s %d times", command, n)
		}
	}
}

func TestWorkflowProviderRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		name, provider string
		omit           bool
	}{{"codex", "codex", false}, {"claude", "claude", false}, {"default_codex", "codex", true}} {
		t.Run(tc.name, func(t *testing.T) {
			f := workflowProviderFixture(t, tc.provider)
			explicit, err := os.ReadFile(f.File)
			if err != nil {
				t.Fatal(err)
			}
			if tc.omit {
				doc := workflowProviderDocument(t, f.File)
				delete(doc["runtime"].(map[string]any), "provider")
				writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.File), doc)
			}
			request, err := ReadWorkflowRequest(f.File)
			if err != nil {
				t.Fatal(err)
			}
			canonical, err := Canonical(request)
			if err != nil || !bytes.Equal(canonical, explicit) {
				t.Fatalf("provider selection changed the expected canonical request: %v", err)
			}
			v, err := WorkflowPreviewStart(f.D, f.File)
			if err != nil {
				t.Fatal(err)
			}
			p := v.(WorkflowPreview)
			workflowProviderReadback(t, p, tc.provider)
			if p.RequestSHA256 != digest(f.R) {
				t.Fatal("preview did not bind the selected provider and model")
			}
			if tc.provider == "codex" {
				oldConfirmation := digest(struct {
					Request  WorkflowRequest
					Observed Observed
					Paths    WorkflowPaths
					Effects  []string
				}{f.R, p.Observed, p.Paths, p.Effects})
				if *p.Confirmation != oldConfirmation || p.Effects[1] != "Create one background Herdr tab and one Codex session" {
					t.Fatal("legacy Codex confirmation or effects changed")
				}
			}
			o, err := WorkflowStart(f.D, f.File, *p.Confirmation)
			if err != nil {
				t.Fatal(err)
			}
			workflowProviderReadback(t, o, tc.provider)
			if o.Transport.AgentSessionID != "fixture-session" || o.RequestSHA256 != p.RequestSHA256 || o.NextAction.Actor != "recipient" {
				t.Fatalf("start did not bind the selected session: %+v", o)
			}
			cwd, _ := f.D.CWD()
			tabs := workflowProviderCalls(t, f, "tab create")
			if len(tabs) != 1 {
				t.Fatalf("tab count = %d", len(tabs))
			}
			wantTab := []string{"tab", "create", "--workspace", "w-fixture", "--cwd", cwd, "--label", "run " + f.R.Herdr.TabLabel, "--env", "PATH=" + os.Getenv("PATH"), "--no-focus"}
			if !reflect.DeepEqual(tabs[0], wantTab) {
				t.Fatalf("wrong tab cwd/argv: %q", tabs[0])
			}
			starts := workflowProviderCalls(t, f, "agent start")
			wantStart := []string{"agent", "start", "ply-" + o.RunID[4:32], "--kind", tc.provider, "--pane", "w-fixture:p1", "--timeout", "30000", "--"}
			if tc.provider == "claude" {
				wantStart = append(wantStart, "--model", f.R.Runtime.Model, "--permission-mode", "manual")
			} else {
				wantStart = append(wantStart, "-C", cwd, "--model", f.R.Runtime.Model, "-a", "on-request", "-c", "approvals_reviewer=\"auto_review\"", "-c", "default_permissions=\""+f.R.Runtime.PermissionBinding.ProfileID+"\"")
			}
			if len(starts) != 1 || !reflect.DeepEqual(starts[0], wantStart) {
				t.Fatalf("wrong selected provider argv: %q; want %q", starts, wantStart)
			}
			o = workflowTestAccept(t, f, o)
			o = workflowTestSettle(t, f, workflowTestReport(t, f, o))
			review := workflowTestReview(t, f, o, "changes_requested")
			o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, review)
			if err != nil || o.Round.Number != 1 || o.Budget.Floor.CorrectionRounds != 1 {
				t.Fatalf("correction did not retain the selected session: %+v %v", o, err)
			}
			if o.Transport.AgentSessionID != "fixture-session" || workflowTestCalls(t, f, "agent start") != 1 || workflowTestCalls(t, f, "agent prompt") != 2 {
				t.Fatal("correction started a new session or sent the wrong number of prompts")
			}
			o = workflowTestSettle(t, f, workflowTestReport(t, f, o))
			o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted"))
			if err != nil || o.FinalReturn.State != "accepted" || o.FinalReturn.TerminalSHA256 == nil || o.TaskResultState != "not_published" {
				t.Fatalf("selected provider round did not seal its native return: %+v %v", o, err)
			}
			workflowProviderReadback(t, o, tc.provider)
			workflowAssertNoReplay(t, f, o)
		})
	}
}

func TestWorkflowProviderRejectsInvalidSelectionBeforeEffects(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value any
	}{{"empty", ""}, {"null", nil}, {"unknown", "other"}, {"wrong_case", "Claude"}, {"wrong_type", true}} {
		t.Run(tc.name, func(t *testing.T) {
			f := workflowProviderFixture(t, "codex")
			doc := workflowProviderDocument(t, f.File)
			doc["runtime"].(map[string]any)["provider"] = tc.value
			writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.File), doc)
			if _, err := WorkflowPreviewStart(f.D, f.File); err == nil {
				t.Fatal("invalid provider selection accepted")
			}
			if _, err := WorkflowStart(f.D, f.File, "invalid"); err == nil {
				t.Fatal("invalid provider selection started")
			}
			workflowProviderNoReservation(t, f)
		})
	}
	for _, name := range []string{"codex_profile", "wrong_permission_mode", "codex_trust"} {
		t.Run(name, func(t *testing.T) {
			f := workflowProviderFixture(t, "claude")
			switch name {
			case "codex_profile":
				f.R.Runtime.ConfigProfile = ptr("codex-only")
			case "wrong_permission_mode":
				f.R.Runtime.PermissionBinding.ProfileID = "other"
			case "codex_trust":
				f.R.CodexProjectTrust = &CodexProjectTrust{"process-local", filepath.Join(f.R.WorkspaceRoot, "target")}
			}
			workflowProviderRebind(t, &f)
			if _, err := WorkflowPreviewStart(f.D, f.File); err == nil {
				t.Fatal("Claude accepted Codex-specific or unsupported runtime configuration")
			}
			if _, err := WorkflowStart(f.D, f.File, "invalid"); err == nil {
				t.Fatal("unsupported Claude configuration started")
			}
			workflowProviderNoReservation(t, f)
		})
	}
}

func TestWorkflowProviderDoesNotOpenNativeTransports(t *testing.T) {
	f := workflowProviderFixture(t, "codex")
	interactive := workflowNativeRequest(f.R)
	_, factory, _, _ := factoryFixture(t)
	for _, tc := range []struct {
		name    string
		request Request
	}{{"interactive", interactive}, {"exec", factory}} {
		t.Run(tc.name, func(t *testing.T) {
			good, _ := Canonical(tc.request)
			if _, err := parseRequest(good); err != nil {
				t.Fatalf("control Codex native request is invalid: %v", err)
			}
			tc.request.Runtime.Provider = "claude"
			invalid, _ := Canonical(tc.request)
			if _, err := parseRequest(invalid); err == nil {
				t.Fatal("Herdr extension opened the Codex-only native transport")
			}
		})
	}
}

func TestWorkflowProviderMissingOrWrongBinaryDoesNotFallback(t *testing.T) {
	for _, provider := range []string{"codex", "claude"} {
		for _, scenario := range []string{"missing_path", "wrong_path", "missing_bound_binary"} {
			t.Run(provider+"/"+scenario, func(t *testing.T) {
				f := workflowProviderFixture(t, provider)
				v, err := WorkflowPreviewStart(f.D, f.File)
				if err != nil {
					t.Fatal(err)
				}
				bin := filepath.Join(f.R.WorkspaceRoot, "bin")
				t.Setenv("PATH", bin+":/usr/bin:/bin")
				selected := filepath.Join(bin, provider)
				if scenario == "missing_bound_binary" {
					if err := os.Remove(f.R.Runtime.Executable.Path); err != nil {
						t.Fatal(err)
					}
				} else {
					if err := os.Remove(selected); err != nil {
						t.Fatal(err)
					}
					if scenario == "wrong_path" {
						if err := os.WriteFile(selected, []byte("#!/bin/sh\nexit 99\n"), 0700); err != nil {
							t.Fatal(err)
						}
					}
				}
				if _, err := WorkflowStart(f.D, f.File, *v.(WorkflowPreview).Confirmation); err == nil {
					t.Fatal("unavailable or wrong selected provider started")
				}
				workflowProviderNoReservation(t, f)
			})
		}
	}
}

func TestWorkflowProviderRejectsWrongStartupIdentity(t *testing.T) {
	for _, provider := range []string{"codex", "claude"} {
		other := map[string]string{"codex": "claude", "claude": "codex"}[provider]
		cases := map[string]map[string]any{
			"outer_provider":   {"agent": other},
			"session_provider": {"agent_session": map[string]any{"agent": other, "kind": "id", "value": "fixture-session"}},
			"session_kind":     {"agent_session": map[string]any{"agent": provider, "kind": "other", "value": "fixture-session"}},
			"session_empty":    {"agent_session": map[string]any{"agent": provider, "kind": "id", "value": ""}},
		}
		if provider == "claude" {
			cases["cwd_changed"] = map[string]any{"foreground_cwd": "/wrong-worktree"}
			cases["cwd_missing"] = map[string]any{"drop_fields": []string{"foreground_cwd"}}
			for _, field := range []string{"workspace_id", "tab_id", "pane_id", "terminal_id", "name"} {
				cases[field] = map[string]any{field: "other"}
			}
		}
		for name, observation := range cases {
			t.Run(provider+"/"+name, func(t *testing.T) {
				f := workflowProviderFixture(t, provider)
				workflowTestModel(t, f, map[string]any{"observations": []any{observation}})
				o, err := workflowReadinessStart(t, f)
				if err == nil || o.NextAction.Actor != "coordinator" || workflowTestCalls(t, f, "agent prompt") != 0 {
					t.Fatalf("wrong provider/location/session permitted Task input: %+v %v", o, err)
				}
				workflowAssertNoReplay(t, f, o)
			})
		}
		t.Run(provider+"/session_changed_before_prompt", func(t *testing.T) {
			f := workflowProviderFixture(t, provider)
			workflowTestModel(t, f, map[string]any{"observations": []any{map[string]any{}, map[string]any{"agent_session": map[string]any{"agent": provider, "kind": "id", "value": "replacement-session"}}}})
			o, err := workflowReadinessStart(t, f)
			if err == nil || workflowTestCalls(t, f, "agent prompt") != 0 {
				t.Fatal("changed native session permitted the Task prompt")
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}

func TestWorkflowProviderAcceptanceAndFollowKeepIdentity(t *testing.T) {
	for _, provider := range []string{"codex", "claude"} {
		other := map[string]string{"codex": "claude", "claude": "codex"}[provider]
		for _, claim := range []string{"runtime", "session"} {
			t.Run(provider+"/accept_"+claim, func(t *testing.T) {
				f := workflowProviderFixture(t, provider)
				o := workflowTestStart(t, f)
				bad := f.Acceptance
				if claim == "runtime" {
					bad.RuntimeClaim.RuntimeID = ptr(other)
				} else {
					bad.RuntimeClaim.NativeSessionID = ptr("replacement-session")
				}
				file := writeAny(t, f.R.WorkspaceRoot, "wrong-claim.json", bad)
				if _, err := WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); err == nil {
					t.Fatal("acceptance accepted another provider or session")
				}
				s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
				if err != nil || s.Acceptance != nil || s.StartSHA256 != nil {
					t.Fatalf("wrong runtime claim published native acceptance: %v", err)
				}
				workflowTestAccept(t, f, o)
			})
		}
		cases := map[string]map[string]any{
			"outer_provider":   {"agent": other},
			"session_provider": {"agent_session": map[string]any{"agent": other, "kind": "id", "value": "fixture-session"}},
			"session_changed":  {"agent_session": map[string]any{"agent": provider, "kind": "id", "value": "replacement-session"}},
		}
		if provider == "claude" {
			cases["cwd_changed"] = map[string]any{"foreground_cwd": "/wrong-worktree"}
		}
		for name, observation := range cases {
			t.Run(provider+"/follow_"+name, func(t *testing.T) {
				f := workflowProviderFixture(t, provider)
				o := workflowTestReport(t, f, workflowTestAccept(t, f, workflowTestStart(t, f)))
				workflowTestModel(t, f, map[string]any{"agent_status": "idle", "observations": []any{observation}})
				o, err := WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, 2*time.Second)
				if err == nil || o.Round.State != "unknown" || o.NextAction.Actor != "coordinator" {
					t.Fatalf("changed provider identity became reviewable: %+v %v", o, err)
				}
				if workflowTestCalls(t, f, "agent start") != 1 || workflowTestCalls(t, f, "agent prompt") != 1 {
					t.Fatal("failed follow started or prompted an agent")
				}
				if _, err := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted")); err == nil {
					t.Fatal("unknown provider identity accepted review")
				}
			})
		}
	}
}

func TestWorkflowProviderAcceptanceRequiresFreshIdentity(t *testing.T) {
	cases := map[string]map[string]any{
		"outer_provider_changed":   {"agent": "codex"},
		"session_provider_changed": {"agent_session": map[string]any{"agent": "codex", "kind": "id", "value": "fixture-session"}},
		"session_replaced":         {"agent_session": map[string]any{"agent": "claude", "kind": "id", "value": "replacement-session"}},
		"cwd_changed":              {"foreground_cwd": "/wrong-worktree"},
	}
	for name, observation := range cases {
		t.Run(name, func(t *testing.T) {
			f := workflowProviderFixture(t, "claude")
			o := workflowTestStart(t, f)
			workflowTestModel(t, f, map[string]any{"observations": []any{observation}})
			file := writeAny(t, f.R.WorkspaceRoot, "valid-claim-after-drift.json", f.Acceptance)
			if _, err := WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); err == nil {
				t.Error("acceptance trusted the cached session after the live Claude identity changed")
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if s.Acceptance != nil || s.StartDraft != nil || s.StartSHA256 != nil {
				t.Error("changed live identity published acceptance or a native start")
			}
			for _, name := range []string{"acceptance.json", "start-draft.json"} {
				if _, err := os.Stat(filepath.Join(o.Paths.RunRoot, name)); !os.IsNotExist(err) {
					t.Errorf("changed live identity created %s: %v", name, err)
				}
			}
			if workflowTestCalls(t, f, "agent start") != 1 || workflowTestCalls(t, f, "agent prompt") != 1 {
				t.Fatal("acceptance drift caused another start or input")
			}
		})
	}
}

func TestWorkflowProviderLostStartIsObservedWithoutReplay(t *testing.T) {
	for _, provider := range []string{"codex", "claude"} {
		for _, ready := range []bool{true, false} {
			name := "unknown"
			if ready {
				name = "recovered_same_session"
			}
			t.Run(provider+"/"+name, func(t *testing.T) {
				f := workflowProviderFixture(t, provider)
				f.D.HerdrTimeout = 2 * time.Second
				model := map[string]any{"lost_command": "agent start"}
				if !ready {
					model["observations"] = []any{workflowPendingObservation()}
				}
				workflowTestModel(t, f, model)
				o, err := workflowReadinessStart(t, f)
				if ready {
					if err != nil || o.Transport.AgentSessionID != "fixture-session" || workflowTestCalls(t, f, "agent prompt") != 1 {
						t.Fatalf("lost reply did not follow the same valid session: %+v %v", o, err)
					}
				} else if err == nil || o.Transport.AgentSessionID != "" || o.NextAction.Actor != "coordinator" || workflowTestCalls(t, f, "agent prompt") != 0 {
					t.Fatalf("unknown start acquired readiness or Task authority: %+v %v", o, err)
				}
				if len(o.Reasons) == 0 || workflowTestCalls(t, f, "tab create") != 1 || workflowTestCalls(t, f, "agent start") != 1 {
					t.Fatal("lost start diagnosis or single attempt binding lost")
				}
				workflowAssertNoReplay(t, f, o)
			})
		}
	}
}

func TestWorkflowProviderLegacyCodexStateRemainsReadableAndFollowable(t *testing.T) {
	f := workflowProviderFixture(t, "codex")
	o := workflowTestStart(t, f)
	files := []string{workflowIndex(f.R.WorkspaceRoot, o.RunID), filepath.Join(o.Paths.RunRoot, "state.json")}
	before := map[string][]byte{}
	for _, file := range files {
		doc := workflowProviderDocument(t, file)
		delete(doc["result"].(map[string]any), "provider")
		writeAny(t, filepath.Dir(file), filepath.Base(file), doc)
		b, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		before[file] = b
	}
	calls := workflowTestCalls(t, f, "agent get")
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil || shown.RequestSHA256 != digest(f.R) || shown.Transport.AgentSessionID != "fixture-session" {
		t.Fatalf("legacy state no longer readable: %+v %v", shown, err)
	}
	workflowProviderReadback(t, shown, "codex")
	for _, file := range files {
		after, err := os.ReadFile(file)
		if err != nil || !bytes.Equal(before[file], after) {
			t.Fatal("readback rewrote legacy state")
		}
	}
	if calls != workflowTestCalls(t, f, "agent get") {
		t.Fatal("legacy show contacted Herdr")
	}
	o = workflowTestSettle(t, f, workflowTestReport(t, f, workflowTestAccept(t, f, shown)))
	o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted"))
	if err != nil || o.FinalReturn.State != "accepted" || o.FinalReturn.TerminalSHA256 == nil {
		t.Fatalf("legacy Codex return no longer follows or seals: %+v %v", o, err)
	}
}

func TestWorkflowProviderCannotReplaceAnExistingRequest(t *testing.T) {
	f := workflowProviderFixture(t, "codex")
	o := workflowTestStart(t, f)
	f.R.Runtime.Provider = "claude"
	f.R.Runtime.PermissionBinding.ProfileID = "manual"
	workflowProviderRebind(t, &f)
	if _, err := WorkflowPreviewStart(f.D, f.File); err == nil {
		t.Fatal("provider change reused an existing run key")
	}
	if _, err := WorkflowStart(f.D, f.File, "already started"); err == nil {
		t.Fatal("provider change replaced the bound native session")
	}
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil || shown.RequestSHA256 != o.RequestSHA256 || shown.Transport.AgentSessionID != o.Transport.AgentSessionID {
		t.Fatalf("rejected provider change damaged existing run: %v", err)
	}
	workflowProviderReadback(t, shown, "codex")
	for _, command := range []string{"tab create", "agent start", "agent prompt"} {
		if workflowTestCalls(t, f, command) != 1 {
			t.Fatal("provider change replayed " + command)
		}
	}
}
