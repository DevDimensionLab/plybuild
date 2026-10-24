package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func capacityRows(session, cwd string) []map[string]any {
	return []map[string]any{
		{"timestamp": "2026-10-20T18:58:00Z", "type": "session_meta", "payload": map[string]any{"id": session, "cwd": cwd}},
		{"timestamp": "2026-10-20T18:58:01Z", "type": "event_msg", "payload": map[string]any{"type": "task_started", "turn_id": "turn-1"}},
		{"timestamp": "2026-10-20T18:58:02Z", "type": "response_item", "payload": map[string]any{"type": "custom_tool_call", "name": "exec", "call_id": "call-1", "input": `text(await tools.exec_command({cmd:"synthetic work",yield_time_ms:1}));`}},
		{"timestamp": "2026-10-20T18:58:03Z", "type": "response_item", "payload": map[string]any{"type": "custom_tool_call_output", "call_id": "call-1", "output": []any{map[string]any{"type": "input_text", "text": "Script completed\nWall time 1 seconds\nOutput:\n"}, map[string]any{"type": "input_text", "text": `{"chunk_id":"running","session_id":41951,"output":""}`}}}},
		{"timestamp": "2026-10-20T18:58:28.503Z", "type": "event_msg", "payload": map[string]any{"type": "task_complete", "turn_id": "turn-1", "root_turn_id": "turn-1", "error": map[string]any{"codex_error_info": "server_overloaded"}}},
		{"timestamp": "2026-10-20T18:58:38.152Z", "type": "event_msg", "payload": map[string]any{"type": "item_completed", "thread_id": session, "turn_id": "turn-1", "started_at_ms": 1792522682000, "completed_at_ms": 1792522718152, "item": map[string]any{"type": "CommandExecution", "id": "exec-1", "process_id": "41951", "cwd": "file://" + cwd, "command": []string{"/bin/sh", "-c", "synthetic work"}, "source": "unified_exec_startup", "status": "completed", "exit_code": 0}}},
	}
}

func capacityTestInput(t *testing.T) DeliveryQuotaReconcileInput {
	t.Helper()
	var in DeliveryQuotaReconcileInput
	if err := json.Unmarshal([]byte(`{"provider":"codex","native_session":"synthetic-session","quota_event":"2026-10-20T18:58:28.503Z","last_tool_result":"exec-1","process_id":"41951","command_sha256":"`+hash([]byte(`["/bin/sh","-c","synthetic work"]`))+`"}`), &in); err != nil {
		t.Fatal(err)
	}
	return in
}

func capacityRaw(t *testing.T, rows []map[string]any) []byte {
	t.Helper()
	var raw []byte
	for _, r := range rows {
		b, e := json.Marshal(r)
		if e != nil {
			t.Fatal(e)
		}
		raw = append(append(raw, b...), '\n')
	}
	return raw
}

func TestCapacityReconcileCompletionAfterOverload(t *testing.T) {
	var s workflowState
	s.Observed.Target.WorktreeLocator = "/synthetic/task"
	in := capacityTestInput(t)
	if err := deliveryValidateQuotaTranscript(s, in, capacityRaw(t, capacityRows(in.NativeSession, s.Observed.Target.WorktreeLocator))); err != nil {
		t.Fatalf("native eventual completion must support retirement-only audit: %v", err)
	}
}

func TestCapacityReconcileRejectsUnfinishedProcesses(t *testing.T) {
	for _, mode := range []string{"mixed_completed_and_running", "running_before_selected_turn"} {
		t.Run(mode, func(t *testing.T) {
			var s workflowState
			s.Observed.Target.WorktreeLocator = "/synthetic/task"
			in := capacityTestInput(t)
			rows := capacityRows(in.NativeSession, s.Observed.Target.WorktreeLocator)
			output := []any{map[string]any{"type": "input_text", "text": "Script completed\nWall time 1 seconds\nOutput:\n"}, map[string]any{"type": "input_text", "text": `{"chunk_id":"selected","session_id":41951,"output":""}`}, map[string]any{"type": "input_text", "text": `{"chunk_id":"finished","exit_code":0,"output":""}`}, map[string]any{"type": "input_text", "text": `{"chunk_id":"running","session_id":777,"output":""}`}}
			if mode == "mixed_completed_and_running" {
				rows[3]["payload"].(map[string]any)["output"] = output
				rows[2]["payload"].(map[string]any)["input"] = `text(await tools.exec_command({cmd:"synthetic work"}));text(await tools.exec_command({cmd:"completed sibling"}));text(await tools.exec_command({cmd:"running sibling"}));`
			} else {
				earlier := map[string]any{"timestamp": "2026-10-20T18:58:00.500Z", "type": "response_item", "payload": map[string]any{"type": "custom_tool_call_output", "call_id": "historical-call", "output": output}}
				priorCall := map[string]any{"timestamp": "2026-10-20T18:58:00.100Z", "type": "response_item", "payload": map[string]any{"type": "custom_tool_call", "name": "exec", "call_id": "historical-call", "input": `text(await tools.exec_command({cmd:"synthetic earlier work"}));`}}
				rows = append(rows[:1], append([]map[string]any{priorCall, earlier}, rows[1:]...)...)
			}
			if err := deliveryValidateQuotaTranscript(s, in, capacityRaw(t, rows)); err == nil {
				t.Fatal("one completed command hid an actually started process without native completion")
			}
		})
	}
}

func TestCapacityReconcileRejectsUnknownEvidence(t *testing.T) {
	cases := []string{"wrong_session", "wrong_cwd", "missing_capacity", "duplicate_capacity", "wrong_error", "missing_start", "missing_return", "duplicate_return", "wrong_process", "wrong_command", "failed_selected", "running_command", "wrong_command_session", "unsupported_tool", "unknown_tail", "permission_denied", "unmatched_tool_return", "duplicate_tool_return", "missing_tool_return", "unsupported_executor_output", "incomplete_jsonl"}
	for _, mode := range cases {
		t.Run(mode, func(t *testing.T) {
			var s workflowState
			s.Observed.Target.WorktreeLocator = "/synthetic/task"
			in := capacityTestInput(t)
			rows := capacityRows(in.NativeSession, s.Observed.Target.WorktreeLocator)
			cmd := rows[5]["payload"].(map[string]any)
			item := cmd["item"].(map[string]any)
			switch mode {
			case "wrong_session":
				rows[0]["payload"].(map[string]any)["id"] = "other"
			case "wrong_cwd":
				rows[0]["payload"].(map[string]any)["cwd"] = "/other"
			case "missing_capacity":
				rows[4]["timestamp"] = "2026-10-20T18:58:28.504Z"
			case "duplicate_capacity":
				rows = append(rows[:5], append([]map[string]any{rows[4]}, rows[5:]...)...)
			case "wrong_error":
				rows[4]["payload"].(map[string]any)["error"] = map[string]any{"codex_error_info": "permission_denied"}
			case "missing_start":
				rows = append(rows[:1], rows[2:]...)
			case "missing_return":
				rows = rows[:5]
			case "duplicate_return":
				rows = append(rows, rows[5])
			case "wrong_process":
				item["process_id"] = "other"
			case "wrong_command":
				item["command"] = []string{"different"}
			case "failed_selected":
				item["status"] = "failed"
				item["exit_code"] = 4
			case "running_command":
				item["status"] = "in_progress"
				delete(item, "exit_code")
			case "wrong_command_session":
				cmd["thread_id"] = "other"
			case "unsupported_tool":
				rows[2]["payload"].(map[string]any)["name"] = "unknown"
			case "unknown_tail":
				rows = append(rows, map[string]any{"timestamp": "2026-10-20T18:59:00Z", "type": "event_msg", "payload": map[string]any{"type": "unknown_activity"}})
			case "permission_denied":
				rows = append(rows, map[string]any{"timestamp": "2026-10-20T18:59:00Z", "type": "event_msg", "payload": map[string]any{"type": "task_complete", "error": map[string]any{"codex_error_info": "permission_denied"}}})
			case "unmatched_tool_return":
				rows[3]["payload"].(map[string]any)["call_id"] = "other"
			case "duplicate_tool_return":
				rows = append(rows[:4], append([]map[string]any{rows[3]}, rows[4:]...)...)
			case "missing_tool_return":
				rows = append(rows[:3], rows[4:]...)
			case "unsupported_executor_output":
				rows[3]["payload"].(map[string]any)["output"] = "completed according to assistant"
			}
			raw := capacityRaw(t, rows)
			if mode == "incomplete_jsonl" {
				raw = raw[:len(raw)-1]
			}
			if err := deliveryValidateQuotaTranscript(s, in, raw); err == nil {
				t.Fatal("invalid native evidence allowed retirement audit")
			}
		})
	}
}

func capacityFixture(t *testing.T, resolvedHistory ...bool) (deliveryFixture, WorkflowRun, string, DeliveryQuotaReconcileInput) {
	t.Helper()
	home, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	t.Setenv("CODEX_HOME", "")
	f, o := retireProviderFixture(t, "codex")
	if len(resolvedHistory) > 0 && resolvedHistory[0] {
		o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
		var err error
		o, err = WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, o.Delivery.ActiveStop.SHA256, "Earlier runtime denial inspected; no effects started", "Use a permitted scratch destination without replaying denied input")
		if err != nil {
			t.Fatal(err)
		}
	}
	o = retireStopFixture(t, f, o, DeliveryEffectPartial)
	path := ownerResumeTranscript(t, f, o)
	rows := capacityRows(o.Transport.AgentSessionID, f.Prepared.Preparation.Preparation.Plan.WorktreePath)
	if err := os.WriteFile(path, capacityRaw(t, rows), 0600); err != nil {
		t.Fatal(err)
	}
	in := capacityTestInput(t)
	in.EventID = "capacity-audit"
	in.StopSHA256 = o.Delivery.ActiveStop.SHA256
	in.NativeSession = o.Transport.AgentSessionID
	in.Action = "synthetic refused action"
	in.SourceSHA256 = hashFileTest(t, path)
	in.Observer = "fixture coordinator"
	in.Reason = "Native eventual completion after model capacity"
	return f, o, path, in
}

func TestCapacityReconcileNativeRetirement(t *testing.T) {
	f, o, path, in := capacityFixture(t)
	stop := *o.Delivery.ActiveStop
	original, _ := os.ReadFile(stop.Locator)
	retireAsCoordinator(&f)
	if _, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: in.Observer, Reason: in.Reason}); err == nil {
		t.Fatal("partial stop allowed retirement before audit")
	}
	out, err := WorkflowReconcileQuotaStop(f.D, f.R.WorkspaceRoot, o.RunID, in)
	if err != nil {
		t.Fatal(err)
	}
	retry, err := WorkflowReconcileQuotaStop(f.D, f.R.WorkspaceRoot, o.RunID, in)
	if err != nil || len(retry.Delivery.Events) != len(out.Delivery.Events) {
		t.Fatal("lost-response retry changed audit", err)
	}
	after, _ := os.ReadFile(stop.Locator)
	if !bytes.Equal(original, after) || out.Delivery.ActiveStop == nil || *out.Delivery.ActiveStop != stop || retireQueueCurrent(t, f) == nil || workflowTestCalls(t, f.workflowFixture, "pane close") != 0 {
		t.Fatal("audit changed stop or released writer")
	}
	if _, err := WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in.StopSHA256, "capacity clarified", "new action"); err == nil {
		t.Fatal("audit reopened owner")
	}
	conflicting := in
	conflicting.Reason = "different"
	if _, err := WorkflowReconcileQuotaStop(f.D, f.R.WorkspaceRoot, o.RunID, conflicting); err == nil {
		t.Fatal("conflicting retry accepted")
	}
	out, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: in.Observer, Reason: in.Reason})
	if err != nil {
		t.Fatal(err)
	}
	if retireQueueCurrent(t, f) != nil || workflowTestCalls(t, f.workflowFixture, "pane close") != 1 || out.Delivery.ActiveStop == nil {
		t.Fatal("ordinary retirement did not preserve stop and release queue")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("retirement removed native evidence", err)
	}
}

func TestCapacityReconcileNativeSafetyGates(t *testing.T) {
	for _, mode := range []string{"wrong_digest", "wrong_session", "wrong_action", "wrong_stop", "native_attempt", "owner_resume_attempt", "rule_stop", "earlier_unknown_stop", "source_drift", "tail_drift_during_close", "remaining_writer", "tampered_copy"} {
		t.Run(mode, func(t *testing.T) {
			f, o, path, in := capacityFixture(t)
			blockedAudit := false
			switch mode {
			case "wrong_digest":
				in.SourceSHA256 = hash([]byte("wrong"))
				blockedAudit = true
			case "wrong_session":
				in.NativeSession = "other"
				blockedAudit = true
			case "wrong_action":
				in.Action = "other"
				blockedAudit = true
			case "wrong_stop":
				in.StopSHA256 = hash([]byte("wrong"))
				blockedAudit = true
			case "native_attempt":
				if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
					s.Result.Delivery.Attempt = &DeliveryAttempt{Kind: "verification", State: "attempted"}
					return nil
				}); err != nil {
					t.Fatal(err)
				}
				blockedAudit = true
			case "owner_resume_attempt":
				if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
					s.OwnerResume = &FileBinding{}
					s.OwnerResumeStep = ownerResumeTabAttempted
					return nil
				}); err != nil {
					t.Fatal(err)
				}
				blockedAudit = true
			case "rule_stop":
				proof := writeAny(t, f.R.WorkspaceRoot, "rule.json", map[string]string{"rule": "permission_denied"})
				var err error
				o, err = WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "rule-stop", Provider: "codex", NativeSession: in.NativeSession, Stage: "work", Observer: in.Observer, Reason: in.Reason, EvidencePath: proof, Action: in.Action, EffectStatus: DeliveryEffectPartial, Rule: "permission_denied"})
				if err != nil {
					t.Fatal(err)
				}
				in.StopSHA256 = o.Delivery.ActiveStop.SHA256
				blockedAudit = true
			case "earlier_unknown_stop":
				o = retireStopFixture(t, f, o, DeliveryEffectUnknown)
				in.StopSHA256 = o.Delivery.ActiveStop.SHA256
			}
			retireAsCoordinator(&f)
			out, err := WorkflowReconcileQuotaStop(f.D, f.R.WorkspaceRoot, o.RunID, in)
			if blockedAudit {
				if err == nil {
					t.Fatal("unsafe audit accepted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			switch mode {
			case "source_drift":
				b, _ := os.ReadFile(path)
				if err := os.WriteFile(path, append(b, '\n'), 0600); err != nil {
					t.Fatal(err)
				}
			case "tail_drift_during_close":
				workflowTestModel(t, f.workflowFixture, map[string]any{"quota_source_path": path, "quota_append_on_close": "\n"})
			case "remaining_writer":
				workflowTestModel(t, f.workflowFixture, map[string]any{"recovery_agents": []any{map[string]any{"workspace_id": "w-fixture", "tab_id": "other", "pane_id": "other", "terminal_id": "other", "agent": "codex", "agent_status": "working", "foreground_cwd": f.Prepared.Preparation.Preparation.Plan.WorktreePath}}})
			case "tampered_copy":
				s, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
				if e != nil {
					t.Fatal(e)
				}
				r, e := deliveryReadQuotaReconciliation(s, out.Delivery.Events[len(out.Delivery.Events)-1].Binding)
				if e != nil {
					t.Fatal(e)
				}
				if e = os.WriteFile(r.Evidence.Locator, []byte("tampered\n"), 0600); e != nil {
					t.Fatal(e)
				}
			}
			_, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: in.Observer, Reason: in.Reason})
			if err == nil {
				t.Fatal("unsafe retirement accepted")
			}
			if retireQueueCurrent(t, f) == nil {
				t.Fatal("blocked retirement released queue")
			}
			if mode == "tail_drift_during_close" && !strings.Contains(err.Error(), "after_close") {
				t.Fatal("source not rechecked after close", err)
			}
		})
	}
}

func TestCapacityReconcileCLIJourney(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o, path, in := capacityFixture(t)
	args := []string{"execute", "reconcile-capacity-stop", o.RunID, "--event-id", in.EventID, "--stop", in.StopSHA256, "--provider", "codex", "--native-session", in.NativeSession, "--action", in.Action, "--capacity-event", in.QuotaEvent, "--command-id", in.LastToolResult, "--process-id", in.ProcessID, "--command-sha256", in.CommandSHA256, "--source-sha256", hashFileTest(t, path), "--observer", in.Observer, "--reason", in.Reason}
	out := automaticCLIOK(t, binary, f, args...)
	out = automaticCLIOK(t, binary, f, args...)
	if len(out.Delivery.Events) != 2 || out.Delivery.ActiveStop == nil || !strings.Contains(out.NextAction.Message, "retire") || retireQueueCurrent(t, f) == nil {
		t.Fatal("CLI audit released ownership or changed stop")
	}
	automaticCLIOK(t, binary, f, "execute", "retire", o.RunID, "--observer", in.Observer, "--reason", in.Reason)
	if retireQueueCurrent(t, f) != nil {
		t.Fatal("CLI ordinary retirement failed")
	}
}

func TestCapacityReconcileRejectedResolutionOnly(t *testing.T) {
	var s workflowState
	s.Result.RunID = "synthetic-run"
	s.Result.Paths.Context = "/synthetic/context.json"
	s.Request.Runtime.PlyExecutable.Path = "/synthetic/ply-control"
	s.Result.Delivery = &DeliveryState{ActiveStop: &FileBinding{Locator: "/synthetic/stop.json"}}
	in := capacityTestInput(t)
	in.StopSHA256 = hash([]byte("synthetic stop"))
	argv := []string{"/bin/zsh", "-lc", s.Request.Runtime.PlyExecutable.Path + " workflow execute resolve-stop " + s.Result.RunID + " --context " + s.Result.Paths.Context + " --stop " + in.StopSHA256 + " --resolution 'capacity observed' --action 'new action'"}
	output := "Run: " + s.Result.RunID + "\nAgent: codex\nDelivery: stopped\nObserved stop: /synthetic/stop.json\nError: workflow_run: unknown or partial external effects remain unresolved; text and permission confirmation cannot clear this stop\n"
	if !capacityRejectedStopResolution(s, in, argv, output, 4) {
		t.Fatal("exact native stop rejection not recognized")
	}
	for _, mode := range []string{"permission_denial", "extra_command", "shell_substitution", "wrong_stop", "wrong_run", "generic_failure", "executable_substitution", "executable_assignment", "wrong_context"} {
		t.Run(mode, func(t *testing.T) {
			a := append([]string(nil), argv...)
			o := output
			exit := 4
			switch mode {
			case "permission_denial":
				o = "Error: permission denied\n"
			case "extra_command":
				a[2] += "; unexpected-effect"
			case "shell_substitution":
				a[2] = strings.Replace(a[2], " --context ", " --context $(effect) ", 1)
			case "wrong_stop":
				a[2] = strings.Replace(a[2], in.StopSHA256, hash([]byte("other")), 1)
			case "wrong_run":
				o = strings.Replace(o, s.Result.RunID, "other", 1)
			case "generic_failure":
				exit = 1
			case "executable_substitution":
				a[2] = strings.Replace(a[2], s.Request.Runtime.PlyExecutable.Path, "$(unexpected-effect)", 1)
			case "executable_assignment":
				a[2] = strings.Replace(a[2], s.Request.Runtime.PlyExecutable.Path, "ENV=value", 1)
			case "wrong_context":
				a[2] = strings.Replace(a[2], s.Result.Paths.Context, "/another/context.json", 1)
			}
			if capacityRejectedStopResolution(s, in, a, o, exit) {
				t.Fatal("unsupported failed command allowed")
			}
		})
	}
}

func TestCapacityReconcilePreCapacityTerminalFailureIsHistory(t *testing.T) {
	for _, mode := range []string{"terminal_before_capacity", "missing_completion", "completion_after_capacity", "return_timestamp_drift", "process_identity_drift"} {
		t.Run(mode, func(t *testing.T) {
			var s workflowState
			s.Observed.Target.WorktreeLocator = "/synthetic/task"
			in := capacityTestInput(t)
			rows := capacityRows(in.NativeSession, s.Observed.Target.WorktreeLocator)
			// The same executor launches the selected command and an unrelated command.
			// The unrelated command's actual exit1 is retained; no shell semantics are inferred.
			out := rows[3]["payload"].(map[string]any)["output"].([]any)
			out = append(out, map[string]any{"type": "input_text", "text": `{"chunk_id":"historical-start","session_id":777,"output":""}`})
			rows[3]["payload"].(map[string]any)["output"] = out
			rows[2]["payload"].(map[string]any)["input"] = `text(await tools.exec_command({cmd:"synthetic work"}));text(await tools.exec_command({cmd:"synthetic-command"}));`
			end := int64(1792522684000)
			timestamp := "2026-10-20T18:58:04Z"
			if mode == "completion_after_capacity" {
				end = 1792522710000
				timestamp = "2026-10-20T18:58:30Z"
			}
			item := map[string]any{"type": "CommandExecution", "id": "historical-failed", "process_id": "777", "cwd": "file:///synthetic/task", "command": []string{"synthetic-command"}, "source": "unified_exec_startup", "status": "failed", "exit_code": 1}
			payload := map[string]any{"type": "item_completed", "thread_id": in.NativeSession, "turn_id": "turn-1", "started_at_ms": 1792522682050, "completed_at_ms": end, "item": item}
			row := map[string]any{"timestamp": timestamp, "type": "event_msg", "payload": payload}
			switch mode {
			case "return_timestamp_drift":
				payload["completed_at_ms"] = end + 1
			case "process_identity_drift":
				item["process_id"] = "other"
			}
			if mode != "missing_completion" {
				at := 4
				if mode == "completion_after_capacity" {
					at = 5
				}
				rows = append(rows[:at], append([]map[string]any{row}, rows[at:]...)...)
			}
			err := deliveryValidateQuotaTranscript(s, in, capacityRaw(t, rows))
			if mode == "terminal_before_capacity" {
				if err != nil {
					t.Fatal("historical native exit1 must remain completed history", err)
				}
			} else if err == nil {
				t.Fatal("unresolved or altered completion allowed")
			}
		})
	}
}

func TestCapacityReconcileBindsResolvedHistory(t *testing.T) {
	f, o, _, in := capacityFixture(t, true)
	retireAsCoordinator(&f)
	before, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	out, err := WorkflowReconcileCapacityStop(f.D, f.R.WorkspaceRoot, o.RunID, in)
	if err != nil {
		t.Fatal(err)
	}
	after, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	r, err := deliveryReadQuotaReconciliation(after, out.Delivery.Events[len(out.Delivery.Events)-1].Binding)
	if err != nil {
		t.Fatal(err)
	}
	if !equal(r.History, before.Result.Delivery.Events) || len(r.History) != 3 || r.History[0].Kind != "coordinator_stop" || r.History[1].Kind != "stop_resolution" || after.Result.Delivery.ActiveStop == nil {
		t.Fatal("audit lost native denial/resolution/capacity history")
	}
	altered := after
	altered.Result.Delivery = &DeliveryState{Events: append([]DeliveryEventRecord(nil), after.Result.Delivery.Events...)}
	altered.Result.Delivery.Events[0].ID = "altered-history"
	if _, err := deliveryReadQuotaReconciliation(altered, out.Delivery.Events[len(out.Delivery.Events)-1].Binding); err == nil {
		t.Fatal("audit accepted changed historical event identity")
	}
}

func TestCapacityReconcileEarlierRunningProcessMustComplete(t *testing.T) {
	for _, completed := range []bool{false, true} {
		t.Run(map[bool]string{false: "missing", true: "completed"}[completed], func(t *testing.T) {
			var s workflowState
			s.Observed.Target.WorktreeLocator = "/synthetic/task"
			in := capacityTestInput(t)
			rows := capacityRows(in.NativeSession, s.Observed.Target.WorktreeLocator)
			earlier := map[string]any{"timestamp": "2026-10-20T18:58:00.500Z", "type": "response_item", "payload": map[string]any{"type": "custom_tool_call_output", "call_id": "historical-call", "output": []any{map[string]any{"type": "input_text", "text": "Script completed\nWall time 1 seconds\nOutput:\n"}, map[string]any{"type": "input_text", "text": `{"chunk_id":"earlier","session_id":777,"output":""}`}}}}
			priorCall := map[string]any{"timestamp": "2026-10-20T18:58:00.100Z", "type": "response_item", "payload": map[string]any{"type": "custom_tool_call", "name": "exec", "call_id": "historical-call", "input": `text(await tools.exec_command({cmd:"synthetic earlier work"}));`}}
			rows = append(rows[:1], append([]map[string]any{priorCall, earlier}, rows[1:]...)...)
			if completed {
				rows = append(rows, map[string]any{"timestamp": "2026-10-20T18:58:40Z", "type": "event_msg", "payload": map[string]any{"type": "item_completed", "thread_id": in.NativeSession, "turn_id": "earlier-turn", "started_at_ms": 1792522680100, "completed_at_ms": 1792522720000, "item": map[string]any{"type": "CommandExecution", "id": "earlier-completion", "process_id": "777", "cwd": "file:///synthetic/task", "command": []string{"synthetic earlier work"}, "source": "unified_exec_startup", "status": "completed", "exit_code": 0}}})
			}
			err := deliveryValidateQuotaTranscript(s, in, capacityRaw(t, rows))
			if completed && err != nil {
				t.Fatal(err)
			}
			if !completed && err == nil {
				t.Fatal("earlier active process was lost")
			}
		})
	}
}

func TestCapacityReconcileRejectsOmittedRunningResult(t *testing.T) {
	var s workflowState
	s.Observed.Target.WorktreeLocator = "/synthetic/task"
	in := capacityTestInput(t)
	rows := capacityRows(in.NativeSession, s.Observed.Target.WorktreeLocator)
	rows[2]["payload"].(map[string]any)["input"] = `const hidden = await tools.exec_command({cmd:"unreported running command"}); text(await tools.exec_command({cmd:"synthetic work"}));`
	if err := deliveryValidateQuotaTranscript(s, in, capacityRaw(t, rows)); err == nil {
		t.Fatal("script-selected finished output hid an unreported running execution")
	}
}

func TestCapacityReconcileLiteralExecutorInventory(t *testing.T) {
	for _, source := range []string{
		`text(await tools.exec_command({cmd:"opaque shell bytes", "max_output_tokens":1000}));`,
		`text(await tools.write_stdin({session_id:41951,chars:"",yield_time_ms:1000}));`,
		`const results=await Promise.allSettled([tools.exec_command({cmd:"first"}),tools.exec_command({cmd:"second"})]);results.forEach(text);`,
	} {
		if _, err := capacityExecutorInput(source); err != nil {
			t.Fatal(err)
		}
	}
	for _, source := range []string{
		`const hidden=await tools.exec_command({cmd:"hidden"});text(await tools.exec_command({cmd:"visible"}));`,
		`text(await tools.exec_command({cmd:computeCommand()}));`,
		`text(await tools.exec_command({...args}));`,
		`text(await tools.exec_command({cmd:"first",cmd:"second"}));`,
		`text(await tools.write_stdin({session_id:41951,chars:"interrupt"}));`,
		`const results=await Promise.allSettled([tools.exec_command({cmd:"hidden"})]);text(results[0]);`,
	} {
		if _, err := capacityExecutorInput(source); err == nil {
			t.Fatal("unsupported or hidden effects accepted", source)
		}
	}
	var s workflowState
	s.Observed.Target.WorktreeLocator = "/synthetic/task"
	in := capacityTestInput(t)
	rows := capacityRows(in.NativeSession, s.Observed.Target.WorktreeLocator)
	rows[2]["payload"].(map[string]any)["input"] = `text(await tools.exec_command({cmd:"synthetic work"}));text(await tools.exec_command({cmd:"second unreturned command"}));`
	if err := deliveryValidateQuotaTranscript(s, in, capacityRaw(t, rows)); err == nil {
		t.Fatal("supported source without every native result was accepted")
	}
}
