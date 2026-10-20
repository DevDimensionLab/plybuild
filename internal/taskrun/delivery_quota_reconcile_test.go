package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Native source guards require physical ancestors. macOS TMPDIR can traverse
// /var -> /private/var; canonicalize only this synthetic provider HOME.
func quotaFixtureHome(t *testing.T) string {
	t.Helper()
	home, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return home
}

func quotaTranscript(t *testing.T, f deliveryFixture, o WorkflowRun) string {
	t.Helper()
	path := ownerResumeTranscript(t, f, o)
	cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	rows := []map[string]any{
		{"type": "assistant", "uuid": "tool-call", "message": map[string]any{"content": []any{map[string]any{"type": "tool_use", "id": "tool-1", "name": "Bash", "input": map[string]any{"command": "go test ./fixture"}}}}},
		{"type": "user", "uuid": "tool-completed", "sourceToolAssistantUUID": "tool-call", "toolUseResult": map[string]any{"interrupted": false}, "message": map[string]any{"content": []any{map[string]any{"type": "tool_result", "tool_use_id": "tool-1", "is_error": false}}}},
		{"type": "assistant", "uuid": "quota-rejected", "apiError": "usage_limit_reached", "isApiErrorMessage": true, "apiErrorStatus": 429, "message": map[string]any{"content": []any{map[string]any{"type": "text", "text": "quota reached"}}}},
		{"type": "system", "uuid": "cancelled", "subtype": "informational"},
	}
	var raw []byte
	for _, row := range rows {
		row["sessionId"], row["cwd"], row["timestamp"] = o.Transport.AgentSessionID, cwd, "2026-10-20T05:12:07.329Z"
		b, err := json.Marshal(row)
		if err != nil {
			t.Fatal(err)
		}
		raw = append(append(raw, b...), '\n')
	}
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func quotaInput(t *testing.T, o WorkflowRun, path string) DeliveryQuotaReconcileInput {
	t.Helper()
	return DeliveryQuotaReconcileInput{EventID: "quota-correction", StopSHA256: o.Delivery.ActiveStop.SHA256, Provider: "claude", NativeSession: o.Transport.AgentSessionID, Action: "synthetic refused action", QuotaEvent: "quota-rejected", LastToolResult: "tool-completed", SourceSHA256: hashFileTest(t, path), Observer: "fixture root", Reason: "Unfinished work was mistaken for a started continuation after quota"}
}

func quotaRows(t *testing.T, path string) []map[string]any {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var rows []map[string]any
	for _, line := range bytes.Split(bytes.TrimSpace(raw), []byte{'\n'}) {
		var row map[string]any
		if err := json.Unmarshal(line, &row); err != nil {
			t.Fatal(err)
		}
		rows = append(rows, row)
	}
	return rows
}

func quotaWriteRows(t *testing.T, path string, rows []map[string]any) []byte {
	t.Helper()
	var raw []byte
	for _, row := range rows {
		b, err := json.Marshal(row)
		if err != nil {
			t.Fatal(err)
		}
		raw = append(append(raw, b...), '\n')
	}
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	return raw
}

// Native async launch results do not complete the delegated effect. Only a
// matching native session-task completion before quota satisfies this boundary;
// assistant text, peer handbacks and tool output cannot substitute for it.
func TestQuotaReconcileTranscriptBoundaries(t *testing.T) {
	t.Setenv("HOME", quotaFixtureHome(t))
	t.Setenv("CLAUDE_CONFIG_DIR", "")
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectPartial)
	path := quotaTranscript(t, f, o)
	in := quotaInput(t, o, path)
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	base := quotaRows(t, path)
	clone := func() []map[string]any {
		b, _ := json.Marshal(base)
		var rows []map[string]any
		if err := json.Unmarshal(b, &rows); err != nil {
			t.Fatal(err)
		}
		return rows
	}
	for _, tc := range []struct {
		name   string
		mutate func([]map[string]any) []map[string]any
	}{
		{"text_only", func(r []map[string]any) []map[string]any { delete(r[2], "apiError"); return r }},
		{"wrong_session", func(r []map[string]any) []map[string]any { r[2]["sessionId"] = "another-session"; return r }},
		{"wrong_cwd", func(r []map[string]any) []map[string]any { r[1]["cwd"] = "/other-task"; return r }},
		{"wrong_selected_call_cwd", func(r []map[string]any) []map[string]any { r[0]["cwd"] = "/other-task"; return r }},
		{"rule_refusal", func(r []map[string]any) []map[string]any { r[2]["apiError"] = "permission_denied"; return r }},
		{"interrupted", func(r []map[string]any) []map[string]any {
			r[1]["toolUseResult"] = map[string]any{"interrupted": true}
			return r
		}},
		{"started_unknown_bash_inside_window", func(r []map[string]any) []map[string]any {
			b, _ := json.Marshal(r[:2])
			var earlier []map[string]any
			if err := json.Unmarshal(b, &earlier); err != nil {
				t.Fatal(err)
			}
			earlier[0]["uuid"] = "earlier-call"
			earlier[0]["message"].(map[string]any)["content"].([]any)[0].(map[string]any)["id"] = "earlier-use"
			earlier[1]["uuid"] = "earlier-result"
			earlier[1]["sourceToolAssistantUUID"] = "earlier-call"
			earlier[1]["message"].(map[string]any)["content"].([]any)[0].(map[string]any)["tool_use_id"] = "earlier-use"
			delete(earlier[1], "toolUseResult")
			return append(r[:2], append(earlier, r[2:]...)...)
		}},
		{"started_error_without_completion_inside_window", func(r []map[string]any) []map[string]any {
			b, _ := json.Marshal(r[:2])
			var earlier []map[string]any
			if err := json.Unmarshal(b, &earlier); err != nil {
				t.Fatal(err)
			}
			earlier[0]["uuid"] = "earlier-error-call"
			earlier[0]["message"].(map[string]any)["content"].([]any)[0].(map[string]any)["id"] = "earlier-error-use"
			earlier[1]["uuid"] = "earlier-error-result"
			earlier[1]["sourceToolAssistantUUID"] = "earlier-error-call"
			block := earlier[1]["message"].(map[string]any)["content"].([]any)[0].(map[string]any)
			block["tool_use_id"], block["is_error"], block["content"] = "earlier-error-use", true, "Exit code 1: failed"
			earlier[1]["toolUseResult"] = "Exit code 1: failed"
			return append(r[:2], append(earlier, r[2:]...)...)
		}},
		{"outstanding_tool", func(r []map[string]any) []map[string]any { return append(r[:1], r[2:]...) }},
		{"orphan_result", func(r []map[string]any) []map[string]any { return r[1:] }},
		{"wrong_tool_source", func(r []map[string]any) []map[string]any { r[1]["sourceToolAssistantUUID"] = "other-call"; return r }},
		{"later_tool", func(r []map[string]any) []map[string]any { r[0]["uuid"] = "later-tool"; return append(r[1:], r[0]) }},
		{"later_assistant", func(r []map[string]any) []map[string]any { r[3]["type"] = "assistant"; return r }},
		{"sidechain", func(r []map[string]any) []map[string]any { r[0]["isSidechain"] = true; return r }},
		{"duplicate_tool", func(r []map[string]any) []map[string]any { return append(r[:1], r...) }},
		{"background_bash", func(r []map[string]any) []map[string]any {
			r[0]["message"].(map[string]any)["content"].([]any)[0].(map[string]any)["input"] = map[string]any{"run_in_background": true}
			return r
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw := quotaWriteRows(t, path, tc.mutate(clone()))
			if err := deliveryValidateQuotaTranscript(s, in, raw); err == nil {
				t.Fatal("unsafe source established pre-execution proof")
			}
		})
	}
	for _, mode := range []string{"completed", "completed_native_string_background", "no_completion", "peer_text", "wrong_task", "failed", "after_quota"} {
		t.Run("async_"+mode, func(t *testing.T) {
			r := clone()
			call := map[string]any{"type": "assistant", "uuid": "async-call", "sessionId": in.NativeSession, "cwd": s.Observed.Target.WorktreeLocator,
				"message": map[string]any{"content": []any{map[string]any{"type": "tool_use", "id": "async-use", "name": "Agent", "input": map[string]any{"run_in_background": true}}}}}
			if mode == "completed_native_string_background" {
				call["message"].(map[string]any)["content"].([]any)[0].(map[string]any)["input"] = map[string]any{"run_in_background": "true"}
			}
			result := map[string]any{"type": "user", "uuid": "async-result", "sessionId": in.NativeSession, "cwd": s.Observed.Target.WorktreeLocator, "sourceToolAssistantUUID": "async-call",
				"toolUseResult": map[string]any{"isAsync": true, "status": "async_launched", "agentId": "agent-a"}, "message": map[string]any{"content": []any{map[string]any{"type": "tool_result", "tool_use_id": "async-use"}}}}
			status, task, origin := "completed", "agent-a", "task-notification"
			if mode == "failed" {
				status = "failed"
			}
			if mode == "wrong_task" {
				task = "agent-b"
			}
			if mode == "peer_text" {
				origin = "peer"
			}
			notice := map[string]any{"type": "attachment", "uuid": "async-notice", "sessionId": in.NativeSession, "cwd": s.Observed.Target.WorktreeLocator,
				"attachment": map[string]any{"type": "queued_command", "commandMode": "task-notification", "origin": map[string]any{"kind": origin, "producer": "session-task"},
					"prompt": "<task-notification><task-id>" + task + "</task-id><tool-use-id>async-use</tool-use-id><status>" + status + "</status></task-notification>"}}
			r = append([]map[string]any{call, result}, r...)
			if mode == "after_quota" {
				r = append(r, notice)
			} else if mode != "no_completion" {
				r = append(r[:2], append([]map[string]any{notice}, r[2:]...)...)
			}
			raw := quotaWriteRows(t, path, r)
			err := deliveryValidateQuotaTranscript(s, in, raw)
			wantPass := mode == "completed" || mode == "completed_native_string_background"
			if (err == nil) != wantPass {
				t.Fatalf("async boundary mode %s: %v", mode, err)
			}
			if !wantPass {
				attempt := in
				attempt.SourceSHA256 = hash(raw)
				if _, err := WorkflowReconcileQuotaStop(f.D, f.R.WorkspaceRoot, o.RunID, attempt); err == nil {
					t.Fatal("unresolved async effect recorded a correction")
				}
				if workflowTestCalls(t, f.workflowFixture, "pane close") != 0 || retireQueueCurrent(t, f) == nil {
					t.Fatal("unresolved async effect released writer")
				}
			}
		})
	}
}

// The target stop names only the next unstarted continuation after the exact
// selected successful result. Matched historical returns are not that action,
// and accepting this window must make no statement about their earlier effects.
func TestQuotaReconcileWindowPreservesHistoricalEffects(t *testing.T) {
	t.Setenv("HOME", quotaFixtureHome(t))
	t.Setenv("CLAUDE_CONFIG_DIR", "")
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectPartial)
	path := quotaTranscript(t, f, o)
	in := quotaInput(t, o, path)
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	base := quotaRows(t, path)
	for _, errorResult := range []bool{false, true} {
		b, _ := json.Marshal(base[:2])
		var historical []map[string]any
		if err := json.Unmarshal(b, &historical); err != nil {
			t.Fatal(err)
		}
		historical[0]["uuid"] = "historic-call"
		historical[0]["message"].(map[string]any)["content"].([]any)[0].(map[string]any)["id"] = "historic-use"
		historical[1]["uuid"], historical[1]["sourceToolAssistantUUID"] = "historic-result", "historic-call"
		historical[0]["cwd"], historical[1]["cwd"] = "/historical-cwd-outside-the-stop-window", "/historical-cwd-outside-the-stop-window"
		block := historical[1]["message"].(map[string]any)["content"].([]any)[0].(map[string]any)
		block["tool_use_id"], block["is_error"] = "historic-use", errorResult
		historical[1]["toolUseResult"] = "historical result text grants no effect attestation"
		raw := quotaWriteRows(t, path, append(historical, base...))
		if err := deliveryValidateQuotaTranscript(s, in, raw); err != nil {
			t.Fatalf("historical metadata gap outside exact continuation window blocked proof (is_error=%v): %v", errorResult, err)
		}
		// Selecting that historical result itself must not establish a boundary.
		bad := in
		bad.LastToolResult = "historic-result"
		if err := deliveryValidateQuotaTranscript(s, bad, raw); err == nil {
			t.Fatal("unknown historical result established the stop boundary")
		}
		// Even a known successful earlier result cannot be selected if more tools
		// ran afterwards: the named continuation would no longer be unstarted.
		historical[1]["toolUseResult"] = map[string]any{"interrupted": false}
		block["is_error"] = false
		raw = quotaWriteRows(t, path, append(historical, base...))
		if err := deliveryValidateQuotaTranscript(s, bad, raw); err == nil {
			t.Fatal("wrong earlier successful boundary hid subsequent tools")
		}
	}
	// Exercise ordinary retirement through the same audited annotation, while
	// preserving the historical return as source rather than effect evidence.
	historical := quotaRows(t, path)
	historical[1]["toolUseResult"] = "historical metadata gap"
	quotaWriteRows(t, path, historical)
	in.SourceSHA256 = hashFileTest(t, path)
	retireAsCoordinator(&f)
	out, err := WorkflowReconcileQuotaStop(f.D, f.R.WorkspaceRoot, o.RunID, in)
	if err != nil {
		t.Fatal(err)
	}
	if out.Delivery.ActiveStop == nil || retireQueueCurrent(t, f) == nil {
		t.Fatal("window annotation released owner")
	}
	if _, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "fixture root", Reason: "replace quota stopped owner"}); err != nil {
		t.Fatal(err)
	}
}

func TestQuotaReconcileRetainsSafetyGates(t *testing.T) {
	for _, mode := range []string{"earlier_unknown", "changed_source", "missing_source", "tampered_copy", "remaining_writer", "ignored_close", "native_attempt", "wrong_action", "wrong_session", "wrong_digest", "rule_stop"} {
		t.Run(mode, func(t *testing.T) {
			t.Setenv("HOME", quotaFixtureHome(t))
			t.Setenv("CLAUDE_CONFIG_DIR", "")
			f, o := retireFixture(t)
			if mode == "earlier_unknown" {
				retireStopFixture(t, f, o, DeliveryEffectUnknown)
			}
			o = retireStopFixture(t, f, o, DeliveryEffectPartial)
			path := quotaTranscript(t, f, o)
			in := quotaInput(t, o, path)
			reconcileBlocked := false
			switch mode {
			case "native_attempt":
				if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
					s.Result.Delivery.Attempt = &DeliveryAttempt{Kind: "verification", State: "attempted"}
					return nil
				}); err != nil {
					t.Fatal(err)
				}
				reconcileBlocked = true
			case "wrong_action":
				in.Action = "a different action"
				reconcileBlocked = true
			case "wrong_session":
				in.NativeSession = "different-session"
				reconcileBlocked = true
			case "wrong_digest":
				in.SourceSHA256 = hash([]byte("wrong"))
				reconcileBlocked = true
			case "rule_stop":
				proof := writeAny(t, f.R.WorkspaceRoot, "rule.json", map[string]string{"refusal": "runtime rule"})
				var err error
				o, err = WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "rule-stop", Provider: "claude", NativeSession: in.NativeSession, Action: in.Action, Stage: "work", EffectStatus: DeliveryEffectPartial, Rule: "runtime-denied", Observer: in.Observer, Reason: "Rule refused action", EvidencePath: proof})
				if err != nil {
					t.Fatal(err)
				}
				in.StopSHA256 = o.Delivery.ActiveStop.SHA256
				reconcileBlocked = true
			}
			retireAsCoordinator(&f)
			out, err := WorkflowReconcileQuotaStop(f.D, f.R.WorkspaceRoot, o.RunID, in)
			if reconcileBlocked {
				if err == nil {
					t.Fatal("invalid correction succeeded")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				switch mode {
				case "changed_source":
					r := quotaRows(t, path)
					r[3]["type"] = "assistant"
					quotaWriteRows(t, path, r)
				case "missing_source":
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				case "tampered_copy":
					s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
					if err != nil {
						t.Fatal(err)
					}
					r, err := deliveryReadQuotaReconciliation(s, out.Delivery.Events[1].Binding)
					if err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(r.Evidence.Locator, []byte("tampered\n"), 0600); err != nil {
						t.Fatal(err)
					}
				case "remaining_writer":
					workflowTestModel(t, f.workflowFixture, map[string]any{"recovery_agents": []any{map[string]any{"workspace_id": "w-fixture", "tab_id": "other", "pane_id": "other", "terminal_id": "other", "agent": "codex", "agent_status": "working", "foreground_cwd": f.Prepared.Preparation.Preparation.Plan.WorktreePath}}})
				case "ignored_close":
					workflowTestModel(t, f.workflowFixture, map[string]any{"close_ignored": true})
				}
			}
			if _, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "fixture root", Reason: "selected replacement"}); err == nil {
				t.Fatal("unsafe retirement succeeded")
			}
			wantCloses := 0
			if mode == "remaining_writer" || mode == "ignored_close" {
				wantCloses = 1
			}
			if workflowTestCalls(t, f.workflowFixture, "pane close") != wantCloses || retireQueueCurrent(t, f) == nil {
				t.Fatal("unsafe pane effect or queue release")
			}
		})
	}
}

func TestQuotaReconcileNeverReopensOwnerAndRecoversLostState(t *testing.T) {
	t.Setenv("HOME", quotaFixtureHome(t))
	t.Setenv("CLAUDE_CONFIG_DIR", "")
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectPartial)
	path := quotaTranscript(t, f, o)
	in := quotaInput(t, o, path)
	fault := f.D
	fault.Files = deliveryRejectStopState{statePath: filepath.Join(o.Paths.RunRoot, "state.json")}
	if _, err := WorkflowReconcileQuotaStop(fault, f.R.WorkspaceRoot, o.RunID, in); err == nil {
		t.Fatal("fault not reached")
	}
	out, err := WorkflowReconcileQuotaStop(f.D, f.R.WorkspaceRoot, o.RunID, in)
	if err != nil || len(out.Delivery.Events) != 2 {
		t.Fatalf("lost-state recovery failed: %v", err)
	}
	if _, err := WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in.StopSHA256, "quota clarified", "different action"); err == nil {
		t.Fatal("correction enabled owner resolution")
	}
	if _, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "claude", NativeSession: in.NativeSession}, true); err == nil {
		t.Fatal("correction enabled receipt")
	}
	report := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "late-owner", PreviousEventSHA256: out.Delivery.LastEventSHA256, Phase: "working", Summary: "Attempted owner continuation", Meaning: "Must stay blocked after retirement-only correction", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	if _, err := WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "late-owner.json", report)); err == nil {
		t.Fatal("correction enabled owner report")
	}
	if _, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, DeliveryOwnerResumeInput{RunID: o.RunID, ProviderExecutable: f.R.Runtime.Executable, ControlExecutable: f.R.Runtime.PlyExecutable}); err == nil {
		t.Fatal("correction enabled owner resume")
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if s.Result.Delivery.ActiveStop == nil || s.Result.Delivery.OwnershipRelease != nil || s.Result.Delivery.Attempt != nil {
		t.Fatal("correction changed authority")
	}
}

// Closing a writer proves liveness, not what it did after the initial source
// check. The second check must stop release without replaying that close.
func TestQuotaReconcileSourceChangeDuringPaneClose(t *testing.T) {
	for _, mode := range []string{"tool_started", "benign_final_metadata"} {
		t.Run(mode, func(t *testing.T) {
			t.Setenv("HOME", quotaFixtureHome(t))
			t.Setenv("CLAUDE_CONFIG_DIR", "")
			f, o := retireFixture(t)
			o = retireStopFixture(t, f, o, DeliveryEffectPartial)
			path := quotaTranscript(t, f, o)
			in := quotaInput(t, o, path)
			retireAsCoordinator(&f)
			out, err := WorkflowReconcileQuotaStop(f.D, f.R.WorkspaceRoot, o.RunID, in)
			if err != nil {
				t.Fatal(err)
			}
			before := out.Delivery.Events[1].Binding
			row := map[string]any{"type": "system", "uuid": "after-close", "sessionId": in.NativeSession, "cwd": f.Prepared.Preparation.Preparation.Plan.WorktreePath, "message": "Automatic continue cancelled"}
			if mode == "tool_started" {
				row["type"] = "assistant"
				row["message"] = map[string]any{"content": []any{map[string]any{"type": "tool_use", "id": "late-tool", "name": "Bash", "input": map[string]any{"command": "fixture side effect"}}}}
			}
			b, err := json.Marshal(row)
			if err != nil {
				t.Fatal(err)
			}
			workflowTestModel(t, f.workflowFixture, map[string]any{"quota_source_path": path, "quota_append_on_close": string(b) + "\n"})
			retire := DeliveryRetireInput{Observer: "fixture root", Reason: "selected replacement"}
			if _, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, retire); err == nil {
				t.Fatal("source changed during pane close but retirement released queue")
			}
			if workflowTestCalls(t, f.workflowFixture, "pane close") != 1 || retireQueueCurrent(t, f) == nil {
				t.Fatal("source race close/release facts wrong")
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if len(s.Result.Delivery.Events) != 2 || s.Result.Delivery.Events[1].Binding != before {
				t.Fatal("source race rewrote correction or recorded retirement")
			}
			// Repeating retire sees the stale source before any second close.
			if _, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, retire); err == nil {
				t.Fatal("stale annotation bypassed source race")
			}
			if workflowTestCalls(t, f.workflowFixture, "pane close") != 1 {
				t.Fatal("source race replayed pane close")
			}
			in.EventID = "fresh-source-audit"
			in.SourceSHA256 = hashFileTest(t, path)
			in.Reason = "Fresh native source after observed close; assess only the same unstarted continuation window"
			fresh, err := WorkflowReconcileQuotaStop(f.D, f.R.WorkspaceRoot, o.RunID, in)
			if mode == "tool_started" {
				if err == nil {
					t.Fatal("new audit cleared genuine post-quota tool effect")
				}
				if retireQueueCurrent(t, f) == nil {
					t.Fatal("post-quota tool effect released queue")
				}
				return
			}
			if err != nil {
				t.Fatalf("fresh audit could not bind benign final metadata: %v", err)
			}
			if len(fresh.Delivery.Events) != 3 || fresh.Delivery.Events[1].Binding != before || fresh.Delivery.ActiveStop == nil {
				t.Fatal("fresh audit replaced history or cleared stop")
			}
			if _, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, retire); err != nil {
				t.Fatalf("retirement did not use newest exact audit: %v", err)
			}
			if retireQueueCurrent(t, f) != nil || workflowTestCalls(t, f.workflowFixture, "pane close") != 1 {
				t.Fatal("benign audit did not release safely or replayed close")
			}
		})
	}
}

// SP01: the correction must cross the public CLI and leave the old refusal
// active until ordinary retirement establishes actual writer absence.
func TestQuotaReconcileCLIJourney(t *testing.T) {
	binary := automaticCLIBinary(t)
	t.Setenv("HOME", quotaFixtureHome(t))
	t.Setenv("CLAUDE_CONFIG_DIR", "")
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectPartial)
	stop := *o.Delivery.ActiveStop
	original, err := os.ReadFile(stop.Locator)
	if err != nil {
		t.Fatal(err)
	}
	path := quotaTranscript(t, f, o)
	if _, _, err = automaticCLI(binary, f, "execute", "retire", o.RunID, "--observer", "fixture root", "--reason", "replace quota stopped owner"); err == nil {
		t.Fatal("false partial did not block retirement before correction")
	}
	args := []string{"execute", "reconcile-quota-stop", o.RunID, "--event-id", "quota-correction", "--stop", stop.SHA256,
		"--provider", "claude", "--native-session", o.Transport.AgentSessionID, "--action", "synthetic refused action",
		"--quota-event", "quota-rejected", "--last-tool-result", "tool-completed", "--source-sha256", hashFileTest(t, path),
		"--observer", "fixture root", "--reason", "Partial described unfinished work; this continuation never started after quota"}
	out := automaticCLIOK(t, binary, f, args...)
	out = automaticCLIOK(t, binary, f, args...)
	if len(out.Delivery.Events) != 2 || out.Delivery.ActiveStop == nil || *out.Delivery.ActiveStop != stop || out.Delivery.Phase != "stopped" || !strings.Contains(out.NextAction.Message, "retire") {
		t.Fatalf("correction cleared stop, duplicated history or hid retirement: %+v", out.Delivery)
	}
	after, err := os.ReadFile(stop.Locator)
	if err != nil || !bytes.Equal(original, after) {
		t.Fatal("correction rewrote original stop")
	}
	if workflowTestCalls(t, f.workflowFixture, "pane close") != 0 || retireQueueCurrent(t, f) == nil {
		t.Fatal("correction released owner")
	}
	args[len(args)-1] = "conflicting correction"
	if _, _, err = automaticCLI(binary, f, args...); err == nil {
		t.Fatal("conflicting retry accepted")
	}
	retireAsCoordinator(&f)
	out = automaticCLIOK(t, binary, f, "execute", "retire", o.RunID, "--observer", "fixture root", "--reason", "replace quota stopped owner")
	if len(out.Delivery.Events) != 3 || retireQueueCurrent(t, f) != nil || workflowTestCalls(t, f.workflowFixture, "pane close") != 1 {
		t.Fatal("ordinary retirement failed after audited correction")
	}
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Fatal(err)
	}
}
