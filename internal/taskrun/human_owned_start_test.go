package taskrun

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestInteractiveStartKeepsProviderApprovalPolicy(t *testing.T) {
	r := WorkflowRequest{Runtime: Runtime{Provider: "codex", Model: "chosen-model", PermissionBinding: Permission{ProfileID: "chosen-profile"}}}
	argv := workflowStartArgv(r, "/chosen/worktree", "chosen-pane")
	for _, arg := range argv {
		if arg == "-a" || strings.HasPrefix(arg, "approvals_reviewer=") {
			t.Fatalf("Ply overrides the human's native approval configuration: %q", argv)
		}
	}
	joined := strings.Join(argv, " ")
	for _, selected := range []string{"--pane chosen-pane", "-C /chosen/worktree", "--model chosen-model", `default_permissions="chosen-profile"`} {
		if !strings.Contains(joined, selected) {
			t.Fatalf("explicit launch input was lost: %s in %q", selected, argv)
		}
	}
}

func TestWorkflowRejectsRemovedTrustOverrides(t *testing.T) {
	for _, field := range []string{"codex_project_trust", "claude_project_trust"} {
		t.Run(field, func(t *testing.T) {
			f := workflowTestFixture(t)
			encoded, err := json.Marshal(f.R)
			if err != nil {
				t.Fatal(err)
			}
			var request map[string]any
			if err := json.Unmarshal(encoded, &request); err != nil {
				t.Fatal(err)
			}
			request[field] = map[string]any{"mode": "configuration"}
			file := writeAny(t, f.R.WorkspaceRoot, "removed-trust.json", request)
			if _, err := ReadWorkflowRequest(file); err == nil || !strings.Contains(err.Error(), field) {
				t.Fatalf("removed trust override accepted: %v", err)
			}
			workflowProviderNoReservation(t, f)
		})
	}
}
