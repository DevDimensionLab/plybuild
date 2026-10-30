package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDeliveryBeginCannotCreateSecondFrozenNotificationBinding(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	path := writeAny(t, f.R.WorkspaceRoot, "existing-notification.json", map[string]string{"kind": "fixture binding"})
	f.R.Delivery.NotificationContext = &FileBinding{path, hashFileTest(t, path)}
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	o := workflowTestStart(t, f.workflowFixture)
	in := beginFixtureInput(f, o)
	in.NotificationRoute = writeAny(t, f.R.WorkspaceRoot, "second-route.json", map[string]any{"kind": "ply.workflow.notification-route", "schema_version": 1, "name": "fixture", "transport": "slack_incoming_webhook", "webhook_env": "PLY_SLACK_WEBHOOK_URL"})
	in.NotificationAuthority = "existing user authorization"
	in.NotificationTimezone = "Europe/Oslo"
	if _, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err == nil {
		t.Fatal("begin created a second notification identity despite the frozen existing binding")
	}
}

func TestDeliveryBeginBindsExistingAuthorizedNotificationNatively(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := workflowTestStart(t, f.workflowFixture)
	in := beginFixtureInput(f, o)
	in.NotificationRoute = writeAny(t, f.R.WorkspaceRoot, "route.json", map[string]any{"kind": "ply.workflow.notification-route", "schema_version": 1, "name": "fixture", "channel_label": "fixture", "state_root": filepath.Join(f.R.WorkspaceRoot, "notifications"), "transport": "slack_incoming_webhook", "webhook_env": "PLY_SLACK_WEBHOOK_URL"})
	in.NotificationAuthority = "Synthetic user authorized this existing route; owner handles completion and questions."
	in.NotificationTimezone = "Europe/Oslo"
	if _, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(o.Paths.RunRoot, "delivery", "notification-task.json")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	var task map[string]any
	if err = json.Unmarshal(b, &task); err != nil {
		t.Fatal(err)
	}
	if task["provider"] != "codex" || task["run"] != o.RunID || task["origin_cwd"] != f.Prepared.Preparation.Preparation.Plan.WorktreePath || task["route"] != in.NotificationRoute || task["notification_authorized"] != true {
		t.Fatalf("wrong native binding: %s", b)
	}
	if _, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(p)
	if err != nil || string(after) != string(b) {
		t.Fatal("replay changed notification identity")
	}
	if _, err = os.Stat(filepath.Join(f.R.WorkspaceRoot, "notifications")); !os.IsNotExist(err) {
		t.Fatal("binding invoked notification transport")
	}
}

func beginFixtureInput(f deliveryFixture, o WorkflowRun) DeliveryBeginInput {
	return DeliveryBeginInput{Provider: f.R.Runtime.Provider, NativeSession: o.Transport.AgentSessionID, Profile: "observed-fixture-policy", Policy: "Synthetic recipient observation: this isolated workspace is writable. No host permission attestation.", ReadRoots: []string{f.R.WorkspaceRoot}, WriteRoots: []string{f.R.WorkspaceRoot}, PermissionConfirmed: true}
}

func TestDeliveryStartupDoesNotConfuseLaunchAcceptanceAndWork(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := workflowTestStart(t, f.workflowFixture)
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	startup := shown.DeliveryStatus.Startup
	if startup == nil || startup.Terminal != "created" || startup.TaskReceived != "unknown" || startup.RuntimePermission != "pending_runtime_acceptance" || startup.TargetWork != "no_changes_observed" {
		t.Fatalf("launch fabricated acceptance or work: %+v", startup)
	}
	shown, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, beginFixtureInput(f, o))
	if err != nil {
		t.Fatal(err)
	}
	startup = shown.DeliveryStatus.Startup
	if startup.TaskReceived != "owner_callback_observed" || startup.RuntimePermission != "recipient_confirmed_contract" || startup.TargetWork != "no_changes_observed" {
		t.Fatalf("acceptance fabricated target work: %+v", startup)
	}
	if err = os.WriteFile(filepath.Join(f.Prepared.Preparation.Preparation.Plan.WorktreePath, "README.md"), []byte("actual target change\n"), 0600); err != nil {
		t.Fatal(err)
	}
	shown, err = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if shown.DeliveryStatus.Startup.TargetWork != "changes_observed" {
		t.Fatal("target work was not distinguished from startup")
	}
}

func TestDeliveryBeginNativeAcceptanceAndLostReply(t *testing.T) {
	for _, provider := range []string{"codex", "claude"} {
		t.Run(provider, func(t *testing.T) {
			f := newDeliveryFixture(t, provider)
			o := workflowTestStart(t, f.workflowFixture)
			in := beginFixtureInput(f, o)
			first, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in)
			if err != nil {
				t.Fatal(err)
			}
			if first.Delivery.PermissionState != "recipient_confirmed_contract" {
				t.Fatal("missing native acceptance")
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			binding := *s.Acceptance
			if _, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err != nil {
				t.Fatal(err)
			}
			s, _ = workflowRead(f.R.WorkspaceRoot, o.RunID)
			if *s.Acceptance != binding || len(workflowProviderCalls(t, f.workflowFixture, "agent start")) != 1 {
				t.Fatal("lost reply replaced acceptance or started another agent")
			}
			in.Policy = "changed observation"
			if _, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err == nil {
				t.Fatal("changed observation overwrote acceptance")
			}
		})
	}
}

func TestDeliveryBeginUnconfirmedAndWrongOwnerCannotAccept(t *testing.T) {
	for _, bad := range []string{"permission", "provider", "session", "policy", "scope", "context"} {
		t.Run(bad, func(t *testing.T) {
			f := newDeliveryFixture(t, "claude")
			o := workflowTestStart(t, f.workflowFixture)
			in, context := beginFixtureInput(f, o), o.Paths.Context
			switch bad {
			case "permission":
				in.PermissionConfirmed = false
			case "provider":
				in.Provider = "codex"
			case "session":
				in.NativeSession = "wrong"
			case "policy":
				in.Policy = ""
			case "scope":
				scratch := filepath.Join(f.R.WorkspaceRoot, "insufficient")
				if err := os.Mkdir(scratch, 0700); err != nil {
					t.Fatal(err)
				}
				in.WriteRoots = []string{scratch}
			case "context":
				context += ".wrong"
			}
			if _, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, context, in); err == nil {
				t.Fatal("invalid observation accepted")
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if s.Acceptance != nil || s.StartSHA256 != nil || s.Result.Delivery.PermissionState != "pending_runtime_acceptance" {
				t.Fatal("invalid observation granted authority")
			}
		})
	}
}
