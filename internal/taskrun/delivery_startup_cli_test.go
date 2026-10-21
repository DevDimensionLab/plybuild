package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// Exercises the actual supplied candidate executable with modeled provider and
// Herdr fixtures. Notification binding does not send Slack or attest runtime authority.
func TestDeliveryStartupCLI(t *testing.T) {
	binary := automaticCLIBinary(t)
	for _, provider := range []string{"codex", "claude"} {
		t.Run(provider, func(t *testing.T) {
			f := newDeliveryFixture(t, provider)
			f.R.Runtime.PlyExecutable = Executable{Path: binary, SHA256: hashFileTest(t, binary)}
			f.D.Executable = func() (string, error) { return binary, nil }
			f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
			o := workflowTestStart(t, f.workflowFixture)
			args := []string{"execute", "begin", o.RunID, "--context", o.Paths.Context, "--provider", provider, "--native-session", o.Transport.AgentSessionID, "--policy", "Synthetic isolated workspace permission observation", "--read-root", f.R.WorkspaceRoot, "--write-root", f.R.WorkspaceRoot}
			if _, _, err := automaticCLI(binary, f, args...); err == nil {
				t.Fatal("unknown authority accepted through CLI")
			}
			args = append(args, "--permission-confirmed")
			route := writeAny(t, f.R.WorkspaceRoot, "route.json", map[string]any{"kind": "ply.workflow.notification-route", "schema_version": 1, "name": "fixture", "transport": "slack_incoming_webhook", "webhook_env": "PLY_SLACK_WEBHOOK_URL"})
			args = append(args, "--notification-route", route, "--notification-authority", "Synthetic standalone owner has the user's standing route authorization", "--notification-timezone", "Europe/Oslo")
			o = automaticCLIOK(t, binary, f, args...)
			if o.Delivery.PermissionState != "recipient_confirmed_contract" {
				t.Fatal("native begin missing acceptance")
			}
			notificationPath := filepath.Join(o.Paths.RunRoot, "delivery", "notification-task.json")
			binding, err := os.ReadFile(notificationPath)
			if err != nil {
				t.Fatal(err)
			}
			var notification map[string]any
			if err = json.Unmarshal(binding, &notification); err != nil {
				t.Fatal(err)
			}
			if len(notification) != 13 || notification["kind"] != "PlyAgentNotificationTask@2" || notification["provider"] != provider || notification["origin_cwd"] != f.Prepared.Preparation.Preparation.Plan.WorktreePath || notification["run"] != o.RunID || notification["activity"] == "" || notification["activity"] == nil {
				t.Fatalf("standalone owner has unusable notification binding: %s", binding)
			}
			proof := writeAny(t, f.R.WorkspaceRoot, "denial.json", map[string]string{"observed": "Synthetic provider rejected report transport"})
			observe := []string{"execute", "observe-stop", o.RunID, "--event-id", "report-denied", "--provider", provider, "--native-session", o.Transport.AgentSessionID, "--stage", "report", "--observer", "fixture coordinator", "--reason", "Provider rejected report transport", "--evidence", proof}
			o = automaticCLIOK(t, binary, f, observe...)
			o = automaticCLIOK(t, binary, f, observe...)
			if o.Delivery.Phase != "stopped" || len(o.Delivery.Events) != 1 {
				t.Fatal("denial lost or duplicated")
			}
			if _, _, err := automaticCLI(binary, f, args...); err == nil {
				t.Fatal("begin ignored active provider stop")
			}
			stop := *o.Delivery.ActiveStop
			denialEvent := o.Delivery.Events[0]
			cause := "Fixture coordinator clarified the synthetic report transport denial; no real provider action was refused"
			action := "Continue isolated fixture startup readback while preserving the rejected report"
			resolve := []string{"execute", "resolve-stop", o.RunID, "--context", o.Paths.Context, "--stop", stop.SHA256, "--resolution", cause, "--action", action}
			o = automaticCLIOK(t, binary, f, resolve...)
			o = automaticCLIOK(t, binary, f, resolve...)
			o = automaticCLIOK(t, binary, f, args...)
			if o.Delivery.Phase != "working" || o.Delivery.ActiveStop != nil || len(o.Delivery.Events) != 2 {
				t.Fatal("supported continuation failed or lost history")
			}
			var record DeliveryStopResolution
			raw, err := workflowBound(o.Delivery.Events[1].Binding, 1<<20)
			if err != nil {
				t.Fatal(err)
			}
			if err = json.Unmarshal(raw, &record); err != nil {
				t.Fatal(err)
			}
			if !equal(o.Delivery.Events[0], denialEvent) || record.Stop != stop || record.RunID != o.RunID || record.NativeSession != o.Transport.AgentSessionID || record.RequestSHA256 != o.RequestSHA256 || record.Resolution != cause || record.ContinuationAction != action || record.PermissionConfirmed {
				t.Fatal("continuation changed the refusal or lost its cause/action binding")
			}
			if len(workflowProviderCalls(t, f.workflowFixture, "agent start")) != 1 {
				t.Fatal("continuation launched another provider")
			}
			after, err := os.ReadFile(notificationPath)
			if err != nil || string(after) != string(binding) {
				t.Fatal("continuation changed standalone notification identity")
			}
		})
	}
}
