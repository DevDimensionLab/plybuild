package taskexecute

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func installEffect(t *testing.T, raw []byte, section string) bool {
	t.Helper()
	var draft struct {
		Authority map[string]json.RawMessage `json:"authority"`
	}
	if err := json.Unmarshal(raw, &draft); err != nil {
		t.Fatal(err)
	}
	var effects []struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(draft.Authority[section], &effects); err != nil {
		t.Fatal(err)
	}
	for _, effect := range effects {
		if effect.Type == "install" {
			return true
		}
	}
	return false
}

func TestInstallAuthorityRequiresExplicitFrozenGoalSelection(t *testing.T) {
	allowed, reserved := true, false
	for _, tc := range []struct {
		name   string
		choice *bool
		want   bool
	}{
		{"omitted", nil, false}, {"reserved", &reserved, false}, {"allowed", &allowed, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d, in, _, _ := launcherFixtureWithInstall(t, "codex", "", tc.choice)
			d.Fault = func(point string) error {
				if point == "workflow_after_reservation" {
					return errors.New("synthetic stop before provider effects")
				}
				return nil
			}
			out, err := Execute(d, in)
			if err == nil || out.Run == nil {
				t.Fatalf("missing reserved Run: %+v %v", out, err)
			}
			var req taskrun.WorkflowRequest
			if _, err := readPreservedJSON(out.RequestPath, &req); err != nil {
				t.Fatal(err)
			}
			if req.Delivery.AllowLocalInstall != tc.want || installEffect(t, req.HandoffDraft, "allowed_effects") != tc.want || installEffect(t, req.HandoffDraft, "forbidden_effects") == tc.want {
				t.Fatalf("owner installation authority not selected by frozen Goal: want %t, contract %+v, draft %s", tc.want, req.Delivery, req.HandoffDraft)
			}
			// Neither local integration nor interactive provider policy grants install.
			if req.Runtime.PermissionBinding.ProfileID != "fixture" {
				t.Fatal("provider policy changed")
			}
			// A direct draft caller must not widen a Goal-derived execution.
			if !tc.want {
				contract := *req.Delivery
				contract.AllowLocalInstall = true
				_, err := taskrun.BuildDeliveryWorkflowRequest(taskrun.DeliveryRequestInput{RequestKey: req.RequestKey, WorkspaceRoot: req.WorkspaceRoot, PreparationID: req.PreparationID, PreparationSHA256: req.PreparationSHA256, HandoffDraft: req.HandoffDraft, Runtime: req.Runtime, HumanAuthority: req.HumanAuthority, HerdrExecutable: req.Herdr.Executable, HerdrWorkspaceID: req.Herdr.WorkspaceID, TabLabel: req.Herdr.TabLabel, Delivery: contract})
				if err == nil {
					t.Fatal("new request option widened reserved Goal installation")
				}
				_, err = workflowhandoff.BuildDeliveryHandoffDraft(d.Workflow, req.PreparationID, "fixture human", "fixture owner", "fixture/overgrant", req.Delivery.AcceptancePath, workflowhandoff.DeliveryDraftAuthority{AllowLocalInstall: true, Agreement: req.Delivery.Agreement})
				if err == nil {
					t.Fatal("direct draft option widened reserved Goal installation")
				}
			}
		})
	}
}

// Registration and startup use the real CLI; only provider/Herdr transport is
// synthetic. Both choices retain the same native automatic integration gates.
func TestInstallAuthorityPublicCLIConsumer(t *testing.T) {
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		t.Skip("requires private candidate PLY_DELIVERY_TEST_BINARY")
	}
	for _, allowed := range []bool{false, true} {
		name := "reserved"
		if allowed {
			name = "allowed"
		}
		t.Run(name, func(t *testing.T) {
			a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic installation consumer", TimeoutSeconds: 30}}
			_, in, root, epic := launcherFixtureWithInstall(t, "codex", binary, &allowed, a)
			if err := os.WriteFile(in.Runtime.HerdrExecutable, []byte(recoveryLauncherHerdr), 0700); err != nil {
				t.Fatal(err)
			}
			var started Result
			if err := json.Unmarshal(runInstallCLI(t, binary, epic, "workflow", "execute", "--spec", "observable-goal", "--format", "json"), &started); err != nil || started.Run == nil {
				t.Fatalf("startup: %+v %v", started, err)
			}
			var req taskrun.WorkflowRequest
			if _, err := readPreservedJSON(started.RequestPath, &req); err != nil {
				t.Fatal(err)
			}
			mandate, err := os.ReadFile(filepath.Join(started.Run.Paths.RunRoot, "mandate.json"))
			if err != nil {
				t.Fatal(err)
			}
			if req.Delivery.AllowLocalInstall != allowed || installEffect(t, mandate, "allowed_effects") != allowed || installEffect(t, mandate, "forbidden_effects") == allowed {
				t.Fatal("consumer authority differs from frozen choice")
			}
			// The exact text sent to Herdr is durable transport evidence.
			rawPrompt, err := os.ReadFile(filepath.Join(started.Run.Paths.RunRoot, "rounds", "000", "prompt-attempt.json"))
			if err != nil {
				t.Fatal(err)
			}
			var prompt struct {
				Argv []string `json:"argv"`
			}
			if err := json.Unmarshal(rawPrompt, &prompt); err != nil {
				t.Fatal(err)
			}
			want := "Local installation allowed: false"
			if allowed {
				want = "Local installation allowed: true"
			}
			if len(prompt.Argv) < 4 || !strings.Contains(prompt.Argv[3], want) {
				t.Fatalf("actual owner prompt missing installation choice: %s", rawPrompt)
			}

			taskPath := started.Goal.WorktreePath
			owner := func(args ...string) []byte {
				return runInstallCLI(t, req.Runtime.PlyExecutable.Path, taskPath, args...)
			}
			o := started.Run
			owner("workflow", "execute", "check-in", o.RunID, "--context", o.Paths.Context, "--provider", "codex", "--native-session", o.Transport.AgentSessionID, "--format", "json")
			if err := os.WriteFile(filepath.Join(taskPath, "README.md"), []byte("explicit install choice "+name+"\n"), 0644); err != nil {
				t.Fatal(err)
			}
			fixtureGit(t, taskPath, "add", "README.md")
			fixtureGit(t, taskPath, "commit", "-m", "synthetic consumer candidate")
			candidate := fixtureGit(t, taskPath, "rev-parse", "HEAD")
			review := filepath.Join(root, "synthetic-review.json")
			fixtureJSON(t, review, map[string]any{"kind": "DeliveryCandidateReview@1", "schema_version": 1, "candidate_oid": candidate, "candidate_tree": fixtureGit(t, taskPath, "rev-parse", "HEAD^{tree}"), "reviewer_claim": "synthetic fixture review; no actual provider review", "reviewer_session_id": nil, "decision": "passed", "findings": []any{}, "fixes": []any{}, "open_actionable_findings": []any{}})
			if err := os.WriteFile(started.Goal.AcceptancePath, []byte("#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = 'explicit install choice "+name+"'\n"), 0600); err != nil {
				t.Fatal(err)
			}
			base := fixtureGit(t, epic, "rev-parse", "HEAD")
			owner("workflow", "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review, "--candidate-binary", binary, "--format", "json")
			if fixtureGit(t, epic, "rev-parse", "HEAD") != base {
				t.Fatal("verification integrated the target")
			}
			var delivered taskrun.WorkflowRun
			if err := json.Unmarshal(owner("workflow", "execute", "integrate", o.RunID, "--context", o.Paths.Context, "--format", "json"), &delivered); err != nil || delivered.DeliveryStatus == nil || delivered.DeliveryStatus.FinalDelivery.State != "delivered" {
				t.Fatalf("native return did not close: %+v, %v", delivered, err)
			}
			if fixtureGit(t, epic, "rev-parse", "HEAD") != candidate {
				t.Fatal("wrong integrated candidate")
			}

		})
	}
}

func runInstallCLI(t *testing.T, binary, cwd string, args ...string) []byte {
	t.Helper()
	command := exec.Command(binary, args...)
	command.Dir = cwd
	var stderr strings.Builder
	command.Stderr = &stderr
	raw, err := command.Output()
	if err != nil {
		t.Fatalf("CLI %v: %v\n%s\n%s", args, err, raw, stderr.String())
	}
	return raw
}
