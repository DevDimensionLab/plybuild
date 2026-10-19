package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestInstallAuthorityLegacyManualDesignDocumentStillBuilds(t *testing.T) {
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{AllowLocalInstall: true})
	raw, err := workflowBound(f.R.Delivery.Goal, 4<<20)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(raw), "# ") || !f.R.Delivery.AllowLocalInstall {
		t.Fatal("fixture did not exercise a legacy plaintext locator with selected installation")
	}
}

func TestInstallAuthorityStructuredGoalCannotFallBackToLegacyDocument(t *testing.T) {
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic structured Goal validation", TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &a})
	r := f.R
	r.Delivery.AllowLocalInstall = false
	var err error
	r.HandoffDraft, err = workflowhandoff.BuildDeliveryHandoffDraft(f.D.Workflow, r.PreparationID, "fixture human", r.Delivery.OwnerClaim, "fixture/structured", r.Delivery.AcceptancePath, workflowhandoff.DeliveryDraftAuthority{Agreement: r.Delivery.Agreement})
	if err != nil {
		t.Fatal(err)
	}
	for _, raw := range []string{"# legacy markdown cannot replace a structured Goal", `null`, `{"contract_kind":"other","executor":{"allow_local_install":true}}`} {
		path := filepath.Join(f.R.WorkspaceRoot, "invalid-structured-goal.json")
		if err := os.WriteFile(path, []byte(raw), 0600); err != nil {
			t.Fatal(err)
		}
		contract := *r.Delivery
		contract.Goal = FileBinding{path, hash([]byte(raw))}
		_, err := BuildDeliveryWorkflowRequest(DeliveryRequestInput{RequestKey: r.RequestKey, WorkspaceRoot: r.WorkspaceRoot, PreparationID: r.PreparationID, PreparationSHA256: r.PreparationSHA256, HandoffDraft: r.HandoffDraft, Runtime: r.Runtime, HumanAuthority: r.HumanAuthority, HerdrExecutable: r.Herdr.Executable, HerdrWorkspaceID: r.Herdr.WorkspaceID, TabLabel: r.Herdr.TabLabel, Delivery: contract})
		if err == nil {
			t.Fatalf("structured Goal silently fell back to a legacy document: %s", raw)
		}
	}
}

// Opt-in, read-only validation of real frozen requests. This calls the normal
// request validator, never startup, receipt, verification or delivery callbacks.
func TestInstallAuthorityHistoricalFrozenRequests(t *testing.T) {
	paths := filepath.SplitList(os.Getenv("PLY_INSTALL_HISTORICAL_REQUESTS"))
	if len(paths) == 0 {
		t.Skip("supply preserved historical request files for read-only validation")
	}
	for _, path := range paths {
		t.Run(filepath.Base(filepath.Dir(path)), func(t *testing.T) {
			before, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			request, err := ReadWorkflowRequest(path)
			if err != nil {
				t.Fatal(err)
			}
			canonical, err := Canonical(request)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(before, canonical) {
				t.Fatal("historical request canonical bytes changed")
			}
			if request.Delivery == nil || !request.Delivery.AllowLocalInstall {
				t.Fatal("historical installation authority was not preserved")
			}
			goal, err := workflowBound(request.Delivery.Goal, 4<<20)
			if err != nil {
				t.Fatal(err)
			}
			var fields struct {
				Executor map[string]json.RawMessage `json:"executor"`
			}
			if err := json.Unmarshal(goal, &fields); err != nil {
				t.Fatal(err)
			}
			if _, present := fields.Executor["allow_local_install"]; present {
				t.Fatal("expected a historical Goal with omitted installation selection")
			}
			after, err := os.ReadFile(path)
			if err != nil || !bytes.Equal(before, after) {
				t.Fatalf("historical request changed: %v", err)
			}
			t.Logf("Frozen omitted-field request retained install=true and passed normal validation unchanged: %s digest %s", path, hash(before))
		})
	}
}
