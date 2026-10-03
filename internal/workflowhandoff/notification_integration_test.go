package workflowhandoff_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	handoff "github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	notification "github.com/devdimensionlab/plybuild/internal/workflownotification"
)

func TestNotificationUsesNativeInspectionWithoutPromotingReportedOutcome(t *testing.T) {
	for _, version := range []int{1, 2, 0} {
		for _, blocked := range []bool{false, true} {
			t.Run(fmt.Sprintf("v%d-blocked-%t", version, blocked), func(t *testing.T) {
				inspection := handoff.NotificationFixture(t, version, blocked)
				var v struct {
					Identities map[string]string `json:"identities"`
					Documents  map[string]struct {
						Locator string         `json:"locator"`
						SHA     string         `json:"sha256"`
						Content map[string]any `json:"content"`
					} `json:"documents"`
					Basis map[string]any `json:"task_spec_basis"`
				}
				if e := json.Unmarshal(inspection, &v); e != nil {
					t.Fatal(e)
				}
				root, e := filepath.EvalSymlinks(t.TempDir())
				if e != nil {
					t.Fatal(e)
				}
				source := map[string]any{"kind": "native", "activity": v.Identities["activity_id"], "run": v.Identities["run_id"], "worktree": v.Documents["handoff"].Content["target_binding"].(map[string]any)["worktree"], "handoff_id": v.Identities["handoff_id"], "task_id": nil}
				if v.Basis != nil {
					source["task_id"] = v.Basis["task_id"]
				}
				for requestName, docName := range map[string]string{"handoff": "handoff", "start_receipt": "start_receipt", "report": "terminal_result"} {
					doc := v.Documents[docName]
					source[requestName] = map[string]string{"path": doc.Locator, "sha256": doc.SHA}
				}
				request := map[string]any{"kind": "ply.workflow.notification-request", "schema_version": 1, "route": "native-fixture", "source": source, "gate": map[string]any{"id": "result-control", "revision": 1, "opened_at": "2026-10-03T19:00:00Z", "reason": "result_control", "state": "waiting_for_human"}, "sender": map[string]string{"actor_claim": "fixture sender"}, "public": map[string]string{"task_title": "Native fixture", "next_action": "Start result control from the local return."}}
				route := notification.Route{Kind: "ply.workflow.notification-route", SchemaVersion: 1, Name: "native-fixture", Transport: "slack_incoming_webhook", ChannelLabel: "#fixture", WebhookEnv: "FIXTURE_WEBHOOK", StateRoot: filepath.Join(root, "state")}
				raw, _ := json.Marshal(request)
				reqPath := filepath.Join(root, "request.json")
				os.WriteFile(reqPath, raw, 0600)
				b, _ := json.Marshal(route)
				routePath := filepath.Join(root, "route.json")
				os.WriteFile(routePath, b, 0600)
				workspaceRoot := v.Documents["handoff"].Content["workspace_binding"].(map[string]any)["root"].(string)
				nativeBefore := notificationNativeSnapshot(t, workspaceRoot)
				reportPath := v.Documents["terminal_result"].Locator
				before, _ := os.ReadFile(reportPath)
				posts := 0
				d := notification.Dependencies{Now: time.Now, LookupEnv: func(string) (string, bool) { return "https://hooks.slack.com/services/TFIX/BFIX/test", true }, Transport: func(context.Context, string, []byte) notification.Observation {
					posts++
					return notification.Observation{Status: 200, Body: []byte("ok"), Dispatch: "started"}
				}}
				in := notification.Input{Operation: "send", File: reqPath, Route: routePath}
				out, e := notification.Execute(d, in)
				if e != nil {
					t.Fatalf("native preview: %v (%+v)", e, out)
				}
				p := out.(notification.Preview)
				in.Apply = true
				in.Confirm = *p.Confirmation
				out, e = notification.Execute(d, in)
				if e != nil {
					t.Fatal(e)
				}
				r := out.(notification.Result)
				if r.Knowledge != "reported" || r.State != "transport_acknowledged" || posts != 1 {
					t.Fatal("native result promoted or lost", r)
				}
				after, _ := os.ReadFile(reportPath)
				if !bytes.Equal(before, after) {
					t.Fatal("native terminal was mutated")
				}
				nativeAfter := notificationNativeSnapshot(t, workspaceRoot)
				if !bytes.Equal(nativeBefore, nativeAfter) {
					t.Fatal("native notification created or changed workspace files")
				}

				source["run"] = "unbound-run"
				bad, _ := json.Marshal(request)
				os.WriteFile(reqPath, bad, 0600)
				in.Apply = false
				in.Confirm = ""
				if _, e = notification.Execute(d, in); e == nil {
					t.Fatal("native identity mismatch accepted")
				}
				os.WriteFile(reqPath, raw, 0600)
				source["run"] = v.Identities["run_id"]
				source["task_id"] = "unbound-task"
				bad, _ = json.Marshal(request)
				os.WriteFile(reqPath, bad, 0600)
				if _, e = notification.Execute(d, in); e == nil {
					t.Fatal("native Task mismatch accepted")
				}
			})
		}
	}
}

func notificationNativeSnapshot(t *testing.T, root string) []byte {
	t.Helper()
	files := map[string]any{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		entry := []any{info.Mode().String()}
		if info.Mode().IsRegular() {
			b, e := os.ReadFile(path)
			if e != nil {
				return e
			}
			entry = append(entry, notification.Digest(b))
		}
		if info.Mode()&os.ModeSymlink != 0 {
			target, e := os.Readlink(path)
			if e != nil {
				return e
			}
			entry = append(entry, target)
		}
		files[path] = entry
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	b, e := json.Marshal(files)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
