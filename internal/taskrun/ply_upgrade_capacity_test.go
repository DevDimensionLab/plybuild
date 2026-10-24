package taskrun

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestPlyUpgradeCapacityRejectionUsesInstalledCLI(t *testing.T) {
	for _, control := range []string{"/synthetic/ply-control", "/ordinary/installed/ply", "ply"} {
		t.Run(control, func(t *testing.T) {
			var s workflowState
			s.Result.RunID = "synthetic-run"
			s.Result.Paths.Context = "/synthetic/context.json"
			s.Observed.Target.WorktreeLocator = "/synthetic/task"
			s.Request.Runtime.PlyExecutable.Path = "/synthetic/ply-control"
			s.Result.Delivery = &DeliveryState{ActiveStop: &FileBinding{Locator: "/synthetic/stop.json"}}
			in := capacityTestInput(t)
			in.StopSHA256 = hash([]byte("synthetic stop"))
			command := control + " workflow execute resolve-stop " + s.Result.RunID + " --context " + s.Result.Paths.Context + " --stop " + in.StopSHA256 + " --resolution 'capacity observed' --action 'new action'"
			output := "Run: " + s.Result.RunID + "\nAgent: codex\nDelivery: stopped\nObserved stop: /synthetic/stop.json\nError: workflow_run: unknown or partial external effects remain unresolved; text and permission confirmation cannot clear this stop\n"
			rows := capacityRows(in.NativeSession, s.Observed.Target.WorktreeLocator)
			returned, _ := json.Marshal(map[string]any{"chunk_id": "rejected", "exit_code": 4, "output": output})
			rows = append(rows,
				map[string]any{"timestamp": "2026-10-20T18:58:40Z", "type": "response_item", "payload": map[string]any{"type": "custom_tool_call", "name": "exec", "call_id": "call-rejection", "input": fmt.Sprintf("text(await tools.exec_command({cmd:%q}));", command)}},
				map[string]any{"timestamp": "2026-10-20T18:58:41Z", "type": "event_msg", "payload": map[string]any{"type": "item_completed", "thread_id": in.NativeSession, "turn_id": "turn-1", "started_at_ms": int64(1792522720000), "completed_at_ms": int64(1792522721000), "item": map[string]any{"type": "CommandExecution", "id": "exec-rejection", "process_id": "41952", "cwd": "file://" + s.Observed.Target.WorktreeLocator, "command": []string{"/bin/zsh", "-lc", command}, "source": "unified_exec_startup", "status": "failed", "exit_code": 4, "aggregated_output": output}}},
				map[string]any{"timestamp": "2026-10-20T18:58:42Z", "type": "response_item", "payload": map[string]any{"type": "custom_tool_call_output", "call_id": "call-rejection", "output": []any{map[string]any{"type": "input_text", "text": "Script completed\nWall time 1 seconds\nOutput:\n"}, map[string]any{"type": "input_text", "text": string(returned)}}}},
			)
			if err := deliveryValidateQuotaTranscript(s, in, capacityRaw(t, rows)); err != nil {
				t.Fatalf("unchanged native rejection through %s blocked reconciliation: %v", control, err)
			}
		})
	}
}
