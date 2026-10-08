package taskrun

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	notification "github.com/devdimensionlab/plybuild/internal/workflownotification"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestIntegrationCLIReconsiderationPreservesBothAnswersAndAllowsNewExactPlan(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f, o, id := integrationCLIFixture(t, binary, "main", true, false)
	trueBin, err := exec.LookPath("true")
	if err != nil {
		t.Fatal(err)
	}
	binaryBytes, err := os.ReadFile(binary)
	if err != nil {
		t.Fatal(err)
	}
	// Installing over the current executable pauses at controller_required,
	// before integration dispatch. The fixture never executes this installation.
	profile := writeAny(t, f.R.WorkspaceRoot, "reconsider-install.json", map[string]any{"kind": "ply.integration.install-profile", "schema_version": 1, "name": "fixture", "candidate_oid": o.Delivery.Candidates[0].OID, "artifact_path": binary, "artifact_sha256": hash(binaryBytes), "install": map[string]any{"executable": trueBin, "arguments": []string{}, "cwd": f.Parent}, "verify": map[string]any{"executable": trueBin, "arguments": []string{}, "cwd": f.Parent}})
	parentBefore := gitOutput(t, f.Parent, "rev-parse", "HEAD")
	r, output, err := integrationCLIRun(binary, f.Parent, "pass\n", "--delivery", id, "--keep", "--install-profile", profile)
	if err == nil || r == nil || deliveryCLIString(t, r, "state") != "controller_required" || string(r["integration_started"]) != "false" {
		t.Fatalf("fixture did not pause before an integration effect: %v %s", err, output)
	}
	operation := deliveryCLIString(t, r, "id")
	decision := append([]byte{}, r["decision"]...)
	var bound FileBinding
	if err = json.Unmarshal(r["decision_file"], &bound); err != nil {
		t.Fatal(err)
	}
	originalDecision, err := os.ReadFile(bound.Locator)
	if err != nil {
		t.Fatal(err)
	}
	before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	p := integrationCLIOK(t, binary, f.Parent, "", "reconsider", operation, "--check")
	if string(p["allowed"]) != "true" {
		t.Fatalf("never-dispatched plan cannot be reconsidered: %s", p)
	}
	for _, input := range []string{"", "pass\n", " blocked\n"} {
		if _, _, err = integrationCLIRun(binary, f.Parent, input, "reconsider", operation, "--observation", "The selected installation target needs a revised plan."); err == nil {
			t.Fatalf("invalid later negative answer accepted: %q", input)
		}
	}
	if after := deliveryCLISnapshot(t, f.R.WorkspaceRoot); before != after {
		t.Fatal("reconsideration preview or rejected input wrote evidence")
	}
	cancelled := integrationCLIOK(t, binary, f.Parent, "blocked\n", "reconsider", operation, "--observation", "The selected installation target needs a revised plan.")
	if deliveryCLIString(t, cancelled, "state") != "cancelled" || !bytes.Equal(decision, cancelled["decision"]) || !bytes.Contains(cancelled["reconsideration"], []byte(`"answer":"blocked"`)) {
		t.Fatalf("original pass or later actual negative decision lost: %s", cancelled)
	}
	if raw, err := os.ReadFile(bound.Locator); err != nil || !bytes.Equal(raw, originalDecision) {
		t.Fatalf("original decision bytes overwritten: %v", err)
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != parentBefore {
		t.Fatal("reconsideration integrated the candidate")
	}
	if _, output, err = integrationCLIRun(binary, f.Parent, "", "resume", operation); err == nil || !strings.Contains(string(output), "plan_revoked") {
		t.Fatalf("old pass was replayed after later negative: %v %s", err, output)
	}
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	qa, err := workspace.LatestTaskHumanQA(registry.HumanQARecords, o.Delivery.Candidates[0].TaskResult.ID, o.Delivery.Candidates[0].OID, o.Delivery.Candidates[0].Tree)
	if err != nil || qa == nil || qa.Outcome != "blocked" {
		t.Fatalf("native latest exact answer was not revoked: %+v %v", qa, err)
	}
	if _, err = ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "The original source owner explicitly releases the unchanged candidate after the plan correction."); err != nil {
		t.Fatal(err)
	}
	replacement := integrationCLIOK(t, binary, f.Parent, "pass\n", "--delivery", id, "--keep")
	if deliveryCLIString(t, replacement, "state") != "completed" || deliveryCLIString(t, replacement, "id") == operation {
		t.Fatalf("revised exact plan did not finish: %s", replacement)
	}
	old := integrationCLIOK(t, binary, f.Parent, "", "show", operation)
	if deliveryCLIString(t, old, "state") != "cancelled" || !bytes.Equal(old["decision"], decision) {
		t.Fatal("replacement altered the old cancelled plan")
	}
	if _, output, err = integrationCLIRun(binary, f.Parent, "blocked\n", "reconsider", deliveryCLIString(t, replacement, "id"), "--observation", "Cannot cancel completed integration."); err == nil || !strings.Contains(string(output), "reconsider_forbidden") {
		t.Fatalf("completed integration could be cancelled: %v %s", err, output)
	}
}

func TestIntegrationCLINotificationRetryUsesNativePreservedAttemptOnly(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f, _, id := integrationCLIFixture(t, binary, "main", true, false)
	routePath := writeAny(t, f.R.WorkspaceRoot, "retry-route.json", notification.Route{Kind: "ply.workflow.notification-route", SchemaVersion: 1, Name: "local-fake-transport", Transport: "slack_incoming_webhook", ChannelLabel: "Fixture only", WebhookEnv: "PLY_TEST_PUBLIC_RETRY_NO_REAL_WEBHOOK", StateRoot: filepath.Join(f.R.WorkspaceRoot, "notification-fixture")})
	t.Setenv("PLY_TEST_PUBLIC_RETRY_NO_REAL_WEBHOOK", "")
	r := integrationCLIOK(t, binary, f.Parent, "pass\n", "--delivery", id, "--keep", "--notification-route", routePath)
	operation := deliveryCLIString(t, r, "id")
	var event notification.IntegrationEvent
	if err := json.Unmarshal(r["notification_event"], &event); err != nil {
		t.Fatal(err)
	}
	route, err := notification.LoadIntegrationRoute(routePath)
	if err != nil {
		t.Fatal(err)
	}
	posts := 0
	d := notification.Dependencies{Now: func() time.Time { return time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC) }, LookupEnv: func(string) (string, bool) { return "https://hooks.slack.com/services/T_TEST/B_TEST/fake-only", true }, Transport: func(context.Context, string, []byte) notification.Observation {
		posts++
		return notification.Observation{Status: 400, Body: []byte("invalid_payload"), Dispatch: "started"}
	}}
	if observed, err := notification.SendIntegration(d, route, event); err != nil || observed.State != "rejected" {
		t.Fatalf("fake native rejected attempt: %+v %v", observed, err)
	}
	before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	p := integrationCLIOK(t, binary, f.Parent, "", "retry-notification", operation, "--check")
	if string(p["allowed"]) != "true" {
		t.Fatalf("public native retry preview rejected known non-delivery: %s", p)
	}
	confirmation := deliveryCLIString(t, p, "confirmation")
	if after := deliveryCLISnapshot(t, f.R.WorkspaceRoot); before != after {
		t.Fatal("notification retry preview wrote state")
	}
	if _, output, err := integrationCLIRun(binary, f.Parent, "", "retry-notification", operation, "--apply", "--confirm", "sha256:"+strings.Repeat("0", 64)); err == nil || !strings.Contains(string(output), "confirmation_changed") {
		t.Fatalf("unconfirmed notification retry accepted: %v %s", err, output)
	}
	d.Transport = func(context.Context, string, []byte) notification.Observation {
		posts++
		return notification.Observation{Status: 200, Body: []byte("ok"), Dispatch: "started"}
	}
	if observed, err := notification.RetryIntegration(d, route, event, confirmation); err != nil || observed.State != "transport_acknowledged" {
		t.Fatalf("fake controlled retry: %+v %v", observed, err)
	}
	// Exact repeated public application must only read the native retry receipt;
	// the actual CLI has no credential and cannot perform a remote operation.
	ack := integrationCLIOK(t, binary, f.Parent, "", "retry-notification", operation, "--apply", "--confirm", confirmation)
	if deliveryCLIString(t, ack, "notification_state") != "transport_acknowledged" || deliveryCLIString(t, ack, "state") != "completed" || posts != 2 {
		t.Fatalf("retry lost native outcome or changed integration: %s posts=%d", ack, posts)
	}
	p = integrationCLIOK(t, binary, f.Parent, "", "retry-notification", operation, "--check")
	if string(p["allowed"]) != "false" {
		t.Fatal("acknowledged notification offered another retry")
	}
}

func TestIntegrationCLIExplicitMergeRetryKeepsRejectedAttemptAndMergesOnce(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f, _, id, model := integrationCLIPRFixture(t, binary, "merge", false, false)
	state := workflowProviderDocument(t, model)
	state["reject_once"] = true
	writeAny(t, f.R.WorkspaceRoot, "pr-fixture.json", state)
	adapter := filepath.Join(f.R.WorkspaceRoot, "pr-adapters", "gh")
	raw, err := os.ReadFile(adapter)
	if err != nil {
		t.Fatal(err)
	}
	marker := `        if m.get("merge_count",0): raise SystemExit("duplicate merge request")`
	reject := `        if m.pop("reject_once",False):
            m["rejection_count"]=m.get("rejection_count",0)+1
            p.write_text(json.dumps(m))
            reply(405,{"message":"Fixture rejects this request before merge"})
`
	if !strings.Contains(string(raw), marker) {
		t.Fatal("fixture merge insertion point changed")
	}
	if err = os.WriteFile(adapter, []byte(strings.Replace(string(raw), marker, reject+marker, 1)), 0700); err != nil {
		t.Fatal(err)
	}
	r, output, err := integrationCLIRun(binary, f.Parent, "pass\n", "--delivery", id, "--method", "merge", "--keep")
	if err == nil || r == nil || deliveryCLIString(t, r, "state") != "merge_blocked" || !bytes.Contains(r["pull_request"], []byte(`"no_effect":true`)) {
		t.Fatalf("known rejected request did not remain retryable: %v %s", err, output)
	}
	operation := deliveryCLIString(t, r, "id")
	r, output, err = integrationCLIRun(binary, f.Parent, "", "resume", operation)
	if err == nil || r == nil || deliveryCLIString(t, r, "state") != "merge_blocked" {
		t.Fatalf("ordinary resume erased known rejection or retried: %v %s", err, output)
	}
	if state := workflowProviderDocument(t, model); state["merge_count"] != nil || state["rejection_count"] != float64(1) {
		t.Fatalf("ordinary resume dispatched again: %v", state)
	}
	r = integrationCLIOK(t, binary, f.Parent, "", "resume", operation, "--retry-merge")
	if deliveryCLIString(t, r, "state") != "completed" || !bytes.Contains(r["merge_attempts"], []byte(`"no_effect":true`)) {
		t.Fatalf("explicit exact retry failed or lost rejection: %s", r)
	}
	integrationCLIOK(t, binary, f.Parent, "", "resume", operation)
	if state := workflowProviderDocument(t, model); state["merge_count"] != float64(1) || state["rejection_count"] != float64(1) {
		t.Fatalf("explicit merge retry was repeated: %v", state)
	}
}
