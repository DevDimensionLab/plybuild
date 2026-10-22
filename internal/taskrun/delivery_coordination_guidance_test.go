package taskrun

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Exercise the actual generators with private files only: no workflow, provider,
// Herdr, policy probes or subprocesses are needed to inspect these instructions.
func TestAL41LaunchDelegationGuidance(t *testing.T) {
	for _, protocol := range []string{"", MandateReceiptV1} {
		for _, allowed := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/allowed=%t", protocol, allowed), func(t *testing.T) {
				root := t.TempDir()
				goal := []byte(`{"title":"Synthetic guidance task","objective":"Inspect generated launch text."}`)
				path := filepath.Join(root, "goal.json")
				if err := os.WriteFile(path, goal, 0600); err != nil {
					t.Fatal(err)
				}
				s := workflowState{Request: WorkflowRequest{
					Runtime: Runtime{Provider: "codex", PlyExecutable: Executable{Path: filepath.Join(root, "control")}},
					Delivery: &DeliveryContract{StartupProtocol: protocol, AllowSubagents: allowed,
						Goal: FileBinding{path, hash(goal)}},
				}, Result: WorkflowRun{RunID: "synthetic-run", Paths: WorkflowPaths{RunRoot: root}}}
				prompt := deliveryInstructions(s)
				for _, want := range []string{
					fmt.Sprintf("Delegation allowed: %t", allowed),
					"technical capability within the frozen Goal/Spec",
					"Follow the selected delegation route and owner/coordinator roles",
					"does not authorize replacing the selected Ply/Herdr route with internal agents",
					"Internal delegation remains available when explicitly authorized by the frozen Goal/Spec",
				} {
					if !strings.Contains(prompt, want) {
						t.Errorf("generated launch lacks %q: %s", want, prompt)
					}
				}
				if s.Request.Delivery.AllowSubagents != allowed {
					t.Fatal("generating guidance changed delegation capability")
				}
			})
		}
	}
}

func TestAL41ReviewCoordinatorReturnGuidance(t *testing.T) {
	s := workflowState{Request: WorkflowRequest{Delivery: &DeliveryContract{}}}
	r := deliveryReviewRound{Coordinator: "synthetic backlog coordinator", CoordinatorPane: "w-fixture:parent", RunID: "synthetic-run"}
	brief := deliveryReviewBrief(s, r, nil, "/fixture", "/fixture/control", "/fixture/template.json", "/fixture/context.json")
	for _, want := range []string{
		"This is an internal review round",
		"Return through native review-return to coordinator synthetic backlog coordinator (pane w-fixture:parent)",
		"When the parent owns the job and its notifications, send no separate Slack notification",
		"Honor an explicit separate notification agreement",
		"For the routine parent-owned return, do not create a notification route or add an approval gate",
		"workflow execute review-return synthetic-run",
		"read-only for you",
	} {
		if !strings.Contains(brief, want) {
			t.Errorf("generated brief lacks %q", want)
		}
	}
}
