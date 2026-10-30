package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Exercise the actual writer without starting a workflow, provider or process.
// This fixture remains runnable after delivery as a bounded guide smoke check.
func TestDeliveryGuideManagedRequirements(t *testing.T) {
	for _, mode := range []string{"automatic", "explicit-review", "legacy-explicit-review", "legacy-report-only"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			s := workflowState{Request: WorkflowRequest{
				Envelope:     deliveryEnv("workflow-request"),
				HandoffDraft: json.RawMessage(`{"task_spec_binding":{"task_id":"synthetic-guide-task"}}`),
				Runtime:      Runtime{Provider: "codex", PlyExecutable: Executable{Path: filepath.Join(root, "unused-control")}},
				Delivery:     &DeliveryContract{StartupProtocol: MandateReceiptV1, AcceptancePath: filepath.Join(root, "acceptance.sh")},
			}, Result: WorkflowRun{RunID: "synthetic-guide-run", Paths: WorkflowPaths{RunRoot: root, Context: filepath.Join(root, "context.json")}}}
			if mode == "automatic" {
				s.Request.Delivery.Agreement = &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic,
					Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic guide fixture", TimeoutSeconds: 30}}
			}
			if mode == "legacy-explicit-review" {
				s.Request.Delivery.StartupProtocol = ""
			}
			if mode == "legacy-report-only" {
				s.Request.SchemaVersion = 1
				s.Request.Delivery = nil
			}
			if err := workflowDeliveryGuide(SystemDependencies(workspace.SystemDependencies()), s); err != nil {
				t.Fatal(err)
			}
			b, err := os.ReadFile(filepath.Join(root, "delivery", "callback-guide.md"))
			if mode == "legacy-report-only" {
				if !os.IsNotExist(err) {
					t.Fatalf("report-only run acquired delivery guidance: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			guide := string(b)
			for _, want := range []string{"Preserve actual command results", "Choose reviews, testing, integration, installation and cleanup explicitly", "Observing a Run never closes a pane or removes a worktree"} {
				if !strings.Contains(guide, want) {
					t.Errorf("guide lacks %q", want)
				}
			}
			if strings.Contains(guide, "workflow execute") {
				t.Fatal("guide recommends a removed process command")
			}

			_, acceptanceErr := os.Stat(filepath.Join(root, "delivery", "templates", "acceptance.json"))
			if mode == "legacy-explicit-review" {
				if acceptanceErr != nil {
					t.Fatal("legacy acceptance contract lost")
				}
			} else if !os.IsNotExist(acceptanceErr) || strings.Contains(guide, "--permission-confirmed") {
				t.Fatal("receipt path acquired legacy acceptance requirements")
			}
		})
	}
}
