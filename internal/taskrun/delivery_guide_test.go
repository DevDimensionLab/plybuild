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
			for _, want := range []string{
				"Native verification assembles and preserves task-requirements",
				"from the frozen Spec requirements and this exact candidate's acceptance and review",
				"Do not hand-author a second requirements/status artifact",
				"Managed artifacts become available during successful candidate verification",
				"Keep actual evidence and truthful requirement outcomes",
				"a progress report or this guide does not qualify a candidate",
				"For an observation preserved in this review, cite candidate-review",
				"--review <copied-review.json>",
				"an empty template is not a review",
				"The callback itself executes /bin/sh",
				"For a clean review-ready candidate, return to the existing coordinator",
				"phase=working and question=null",
				"candidate OID/tree, test evidence and private binary path",
				"Keep the candidate unchanged until the separate review returns",
				"Do not use an interactive human question for this routine coordinator transition",
				"Real human QA, access dialogs and new product or scope decisions",
			} {
				if !strings.Contains(guide, want) {
					t.Errorf("generated guide lacks %q", want)
				}
			}
			if mode == "automatic" {
				if !strings.Contains(guide, "--candidate-binary <absolute-private-executable>") || !strings.Contains(guide, "## Automatic local Epic delivery") {
					t.Fatal("automatic candidate instructions lost")
				}
			} else if strings.Contains(guide, "--candidate-binary") || !strings.Contains(guide, "## Actual human QA and local integration") {
				t.Fatal("explicit-review path acquired automatic delivery")
			}
			_, acceptanceErr := os.Stat(filepath.Join(root, "delivery", "templates", "acceptance.json"))
			if mode == "legacy-explicit-review" {
				if acceptanceErr != nil || !strings.Contains(guide, "--permission-confirmed") {
					t.Fatal("legacy acceptance contract lost")
				}
			} else if !os.IsNotExist(acceptanceErr) || strings.Contains(guide, "--permission-confirmed") {
				t.Fatal("receipt path acquired legacy acceptance requirements")
			}
		})
	}
}
