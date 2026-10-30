package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDeliveryResumedReadbackSeparatesCurrentStateFromStartupHistory(t *testing.T) {
	f, o := deliveryPausedStartup(t)
	historical := []Reason{{"herdr_start_response", "Synthetic start reply timed out."}, {"herdr_start_stopped", "Synthetic startup stopped; no Task prompt sent."}, {"delivery_start_resume_stopped", "Synthetic previous resume stopped."}}
	if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
		s.Result.Reasons = append(historical, Reason{"unrelated_observation", "Keep this current observation."})
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	before, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil || len(before.Reasons) != 4 || before.Round.State != "unknown" {
		t.Fatalf("unresolved startup lost: %+v %v", before, err)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"observations": nil, "agent_status": "idle"})
	o, err = WorkflowResumeDeliveryStart(f.D, f.R.WorkspaceRoot, o.RunID, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(o.Paths.RunRoot, "state.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, code := range []string{"herdr_start_response", "herdr_start_stopped", "delivery_start_resume_stopped"} {
		if !bytes.Contains(raw, []byte(code)) {
			t.Fatalf("historical reason %s was removed from state", code)
		}
	}
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(raw, after) {
		t.Fatal("readback rewrote historical state")
	}
	for _, got := range []WorkflowRun{o, shown} {
		if len(got.Reasons) != 1 || got.Reasons[0].Code != "unrelated_observation" || got.Round.State != "awaiting_acceptance" || got.Transport.State != "working" {
			t.Fatalf("resolved startup is still presented as current attention: %+v", got)
		}
	}
}

func TestDeliveryActualPolicyDoesNotInventRequestedProfileName(t *testing.T) {
	for _, tc := range []struct {
		provider string
		profile  *string
	}{{"codex", ptr("managed")}, {"claude", ptr("auto")}, {"codex", nil}} {
		t.Run(tc.provider+"/"+workflowJSON(tc.profile), func(t *testing.T) {
			f := newDeliveryFixture(t, tc.provider)
			o := workflowTestStart(t, f.workflowFixture)
			a := deliveryTestAcceptance(t, f, o)
			a.RuntimeClaim.ProfileID = tc.profile
			for _, mutation := range []string{"unconfirmed", "policy_digest", "provider", "native_session", "scope"} {
				bad := a
				switch mutation {
				case "unconfirmed":
					bad.DeliveryPermission.PermissionConfirmed = false
				case "policy_digest":
					bad.RuntimeClaim.EffectivePolicySHA256 = ptr(hash([]byte("unobserved policy")))
				case "provider":
					bad.RuntimeClaim.RuntimeID = ptr("another-provider")
				case "native_session":
					bad.RuntimeClaim.NativeSessionID = ptr("another-session")
				case "scope":
					var sandbox map[string]any
					if err := json.Unmarshal(a.Sandbox, &sandbox); err != nil {
						t.Fatal(err)
					}
					sandbox["write_roots"] = []string{}
					bad.Sandbox, _ = Canonical(sandbox)
				}
				if _, err := WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "bad-"+mutation+".json", bad)); err == nil {
					t.Fatalf("missing actual authority or identity accepted: %s", mutation)
				}
			}
			got, err := WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "actual-runtime.json", a))
			if err != nil {
				t.Fatalf("truthful actual policy label rejected against requested %q: %v", f.R.Runtime.PermissionBinding.ProfileID, err)
			}
			if got.Delivery.PermissionState != "recipient_confirmed_contract" || got.Delivery.ActualPolicySHA256 == nil || *got.Delivery.ActualPolicySHA256 != *a.RuntimeClaim.EffectivePolicySHA256 {
				t.Fatalf("actual evidence was not preserved: %+v", got.Delivery)
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil || s.Request.Runtime.PermissionBinding.ProfileID != f.R.Runtime.PermissionBinding.ProfileID {
				t.Fatalf("launch selector was rewritten: %+v %v", s.Request.Runtime.PermissionBinding, err)
			}
			var stored DeliveryAcceptance
			if err = readValue(s.Acceptance.Locator, 1<<20, &stored); err != nil || !equal(stored.RuntimeClaim.ProfileID, tc.profile) {
				t.Fatalf("actual profile label was replaced: %+v %v", stored.RuntimeClaim, err)
			}
		})
	}
}

func TestDeliveryClaudePermissionModeUsesBoundRequest(t *testing.T) {
	f := newDeliveryFixture(t, "claude")
	for _, mode := range []string{"manual", "auto"} {
		r := f.R
		r.Runtime.PermissionBinding.ProfileID = mode
		if err := validateDeliveryWorkflowRequest(r); err != nil {
			t.Errorf("documented Claude mode %s rejected: %v", mode, err)
		}
		argv := strings.Join(workflowStartArgv(r, f.Prepared.Preparation.Preparation.Plan.WorktreePath, "fixture-pane"), " ")
		if !strings.Contains(argv, "--permission-mode "+mode) || strings.Contains(argv, "default_permissions") || strings.Contains(argv, "auto_review") {
			t.Errorf("bound Claude mode %s not used: %s", mode, argv)
		}
	}
	r := f.R
	r.Runtime.PermissionBinding.ProfileID = "bypassPermissions"
	if err := validateDeliveryWorkflowRequest(r); err == nil {
		t.Fatal("unsupported mode granted")
	}
}

func TestDeliveryFollowExpiredObservationKeepsKnownRun(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := workflowTestStart(t, f.workflowFixture)
	gets := workflowTestCalls(t, f.workflowFixture, "agent get")
	got, err := WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, time.Nanosecond)
	if err == nil || !strings.Contains(err.Error(), "observation timed out without stopping or restarting") {
		t.Fatalf("observer deadline presented as an unknown startup: %+v %v", got, err)
	}
	if got.Transport.State != o.Transport.State || got.Transport.Observation != "cached" || got.Delivery.Phase != o.Delivery.Phase || len(got.Reasons) != len(o.Reasons) {
		t.Fatalf("observer deadline changed the last known run: before=%+v after=%+v", o, got)
	}
	if workflowTestCalls(t, f.workflowFixture, "agent get") != gets || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 {
		t.Fatal("expired observation performed an extra transport call")
	}
}
