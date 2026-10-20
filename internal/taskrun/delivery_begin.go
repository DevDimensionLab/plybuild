package taskrun

import (
	"encoding/json"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// DeliveryBeginInput contains only recipient observations. Run identities,
// agreement bindings, evidence hashes and private receipt files belong to Ply.
// Policy is the recipient's observation, never an attestation by the host.
type DeliveryBeginInput struct {
	Provider              string   `json:"provider"`
	NativeSession         string   `json:"native_session"`
	Model                 string   `json:"model"`
	Profile               string   `json:"profile"`
	Policy                string   `json:"policy"`
	ReadRoots             []string `json:"read_roots"`
	WriteRoots            []string `json:"write_roots"`
	PermissionConfirmed   bool     `json:"permission_confirmed"`
	NotificationRoute     string   `json:"notification_route,omitempty"`
	NotificationAuthority string   `json:"notification_authority,omitempty"`
	NotificationTimezone  string   `json:"notification_timezone,omitempty"`
}

func WorkflowDeliveryBegin(d Dependencies, root, id, contextPath string, in DeliveryBeginInput) (WorkflowRun, error) {
	if err := containing(d, root); err != nil {
		return WorkflowRun{}, err
	}
	s, err := workflowRead(root, id)
	if err != nil {
		return WorkflowRun{}, err
	}
	if !deliveryRun(s.Request) || headlessRun(s.Request) {
		return WorkflowRun{}, workflowError(4, "begin requires an interactive delivery owner")
	}
	if deliveryObservedStop(s) {
		return deliveryReadback(d, root, id, workflowError(4, "resolve the preserved provider stop before begin; do not repeat denied input"))
	}
	if s.Request.Delivery.NotificationContext != nil && in.NotificationRoute != "" {
		return deliveryReadback(d, root, id, workflowError(4, "reuse the frozen notification context; begin cannot create a second notification identity"))
	}
	if err = workflowCallback(d, s, contextPath); err != nil {
		return deliveryReadback(d, root, id, err)
	}
	if in.Provider != s.Request.Runtime.Provider || in.NativeSession == "" || in.NativeSession != s.Result.Transport.AgentSessionID {
		return deliveryReadback(d, root, id, workflowError(4, "actual provider and native session must match the bound owner"))
	}
	receipt := mandateReceipt(s.Request)
	if receipt && in.PermissionConfirmed {
		return deliveryReadback(d, root, id, workflowError(2, "mandate receipt does not accept blanket permission confirmation"))
	}
	if !receipt && (!in.PermissionConfirmed || strings.TrimSpace(in.Policy) == "" || len(in.Policy) > 32768) {
		return deliveryReadback(d, root, id, workflowError(4, "actual permission is unconfirmed: preserve the stop; no target writes are authorized"))
	}
	if len(in.Policy) > 32768 {
		return WorkflowRun{}, workflowError(2, "observation exceeds limit")
	}
	if !receipt && (len(in.ReadRoots) == 0 || len(in.WriteRoots) == 0) {
		return deliveryReadback(d, root, id, workflowError(2, "actual read and write roots are required"))
	}
	if err = validateBeginNotification(in); err != nil {
		return deliveryReadback(d, root, id, err)
	}
	for _, roots := range [][]string{in.ReadRoots, in.WriteRoots} {
		for _, path := range roots {
			if !filepath.IsAbs(path) || filepath.Clean(path) != path {
				return WorkflowRun{}, workflowError(2, "actual roots must be absolute physical paths")
			}
		}
	}
	in.ReadRoots = slices.Clone(in.ReadRoots)
	in.WriteRoots = slices.Clone(in.WriteRoots)
	slices.Sort(in.ReadRoots)
	slices.Sort(in.WriteRoots)
	in.ReadRoots = slices.Compact(in.ReadRoots)
	in.WriteRoots = slices.Compact(in.WriteRoots)
	// Content addressing makes a lost response replay the same observations.
	// Different observations cannot replace an already accepted native receipt.
	dir := filepath.Join(s.Result.Paths.RunRoot, "delivery", "begin", strings.TrimPrefix(digest(in), "sha256:"))
	if err = privateDir(dir); err != nil {
		return WorkflowRun{}, err
	}
	evidenceName, evidenceRole := "actual-policy.json", "effective_policy"
	if receipt {
		evidenceName, evidenceRole = "runtime-observation.json", "runtime_observation"
	}
	policy, err := workflowKeep(d, filepath.Join(dir, evidenceName), in)
	if err != nil {
		return WorkflowRun{}, err
	}
	tmp := filepath.Join(dir, "tmp")
	if err = privateDir(tmp); err != nil {
		return WorkflowRun{}, err
	}
	var matches any = true
	if receipt {
		matches = nil
	}
	reads, writes := in.ReadRoots, in.WriteRoots
	if reads == nil {
		reads = []string{}
	}
	if writes == nil {
		writes = []string{}
	}
	sandbox, err := Canonical(map[string]any{"read_roots": reads, "write_roots": writes, "temp_root": tmp, "matches_contract": matches})
	if err != nil {
		return WorkflowRun{}, err
	}
	claim := RuntimeClaim{RuntimeID: ptr(in.Provider), NativeSessionID: ptr(in.NativeSession), EffectivePolicySHA256: ptr(policy.SHA256)}
	if in.Model != "" {
		claim.ModelID = ptr(in.Model)
	}
	if in.Profile != "" {
		claim.ProfileID = ptr(in.Profile)
	}
	permission := DeliveryPermissionAcceptance{LaunchContractSHA256: s.Request.Runtime.PermissionBinding.EffectivePolicySHA256, PermissionConfirmed: true, ActualPolicyEvidence: []Evidence{{Locator: policy.Locator, SHA256: policy.SHA256, Role: evidenceRole}}}
	if receipt {
		permission.PermissionConfirmed = false
		permission.StartupProtocol = MandateReceiptV1
		permission.MandateSHA256 = s.Result.Handoff.SHA256
	}
	if a := s.Request.Delivery.Agreement; a != nil {
		permission.DeliveryAgreementSHA256 = workspace.DeliveryAgreementDigest(*a)
		permission.AllowedEffects = workflowhandoff.DeliveryAllowedEffects(*a)
	}
	a := DeliveryAcceptance{Acceptance: Acceptance{Envelope: deliveryEnv("run-acceptance"), RunID: id, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID, RuntimeClaim: claim, Sandbox: sandbox, Acceptance: "started", Issues: json.RawMessage("[]")}, DeliveryPermission: permission}
	input, err := workflowKeep(d, filepath.Join(dir, "acceptance.json"), a)
	if err != nil {
		return WorkflowRun{}, err
	}
	out, err := WorkflowAccept(d, root, id, contextPath, input.Locator)
	if err == nil && in.NotificationRoute != "" {
		err = deliveryBindNotification(d, s, in)
		if err == nil {
			out.NextAction.Message += " Notification binding: " + filepath.Join(s.Result.Paths.RunRoot, "delivery", "notification-task.json")
		}
	}
	if err == nil && s.Request.Delivery.NotificationContext != nil {
		out.NextAction.Message += " Notification binding: " + s.Request.Delivery.NotificationContext.Locator
	}
	return out, err
}

func validateBeginNotification(in DeliveryBeginInput) error {
	if in.NotificationRoute == "" {
		if in.NotificationAuthority != "" || in.NotificationTimezone != "" {
			return workflowError(2, "notification authority requires an existing route")
		}
		return nil
	}
	if in.NotificationTimezone == "" || in.NotificationTimezone == "Local" {
		return workflowError(2, "notification route requires the recipient's explicit IANA timezone")
	}
	if _, err := time.LoadLocation(in.NotificationTimezone); err != nil {
		return workflowError(2, "invalid notification timezone")
	}
	if strings.TrimSpace(in.NotificationAuthority) == "" || len(in.NotificationAuthority) > 4096 {
		return workflowError(2, "the existing user authorization for this notification route is required")
	}
	b, err := readFile(in.NotificationRoute, 16<<10, true)
	if err != nil {
		return err
	}
	var route struct {
		Kind          string `json:"kind"`
		SchemaVersion int    `json:"schema_version"`
		Name          string `json:"name"`
		Transport     string `json:"transport"`
		WebhookEnv    string `json:"webhook_env"`
	}
	if err = json.Unmarshal(b, &route); err != nil {
		return err
	}
	if route.Kind != "ply.workflow.notification-route" || route.SchemaVersion != 1 || route.Name == "" || route.Transport != "slack_incoming_webhook" || route.WebhookEnv != "PLY_SLACK_WEBHOOK_URL" {
		return workflowError(2, "unsupported existing notification route; no setup or global change is performed")
	}
	return nil
}

func deliveryBindNotification(d Dependencies, s workflowState, in DeliveryBeginInput) error {
	view, err := workflowhandoff.TaskRunHandoffView(d.Workflow, s.Result.Handoff.Locator)
	if err != nil {
		return err
	}
	var handoff struct {
		Identity struct {
			ActivityID string `json:"activity_id"`
		} `json:"identity"`
	}
	if err = json.Unmarshal(view, &handoff); err != nil {
		return err
	}
	if handoff.Identity.ActivityID == "" {
		return workflowError(4, "native handoff has no notification activity identity")
	}
	dir := filepath.Join(s.Result.Paths.RunRoot, "delivery")
	authority, err := workflowKeep(d, filepath.Join(dir, "notification-authority.json"), map[string]any{"kind": "ply.workflow.notification-authority", "run_id": s.Result.RunID, "owner_claim": s.Request.Delivery.OwnerClaim, "provider": in.Provider, "source": in.NotificationAuthority, "route": in.NotificationRoute, "goal": s.Request.Delivery.Goal, "human_authority": s.Request.HumanAuthority})
	if err != nil {
		return err
	}
	_, err = workflowKeep(d, filepath.Join(dir, "notification-task.json"), map[string]any{"kind": "PlyAgentNotificationTask@2", "activity": handoff.Identity.ActivityID, "run": s.Result.RunID, "worktree": s.Observed.Target.WorktreeLocator, "handoff": authority.Locator, "route": in.NotificationRoute, "event_root": filepath.Join(dir, "notification-events"), "actor_claim": s.Request.Delivery.OwnerClaim, "notification_authorized": true, "provider": in.Provider, "origin_cwd": s.Observed.Target.WorktreeLocator, "context_root": s.Request.WorkspaceRoot, "timezone": in.NotificationTimezone})
	return err
}

// WorkflowDeliveryReceipt routes by the frozen contract, never by a new permission claim.
func WorkflowDeliveryReceipt(d Dependencies, root, id, contextPath string, in DeliveryBeginInput, receipt bool) (WorkflowRun, error) {
	s, err := workflowRead(root, id)
	if err != nil {
		return WorkflowRun{}, err
	}
	if reportReceipt(s.Request) && receipt {
		return WorkflowReportReceipt(d, root, id, contextPath, in)
	}
	if mandateReceipt(s.Request) != receipt {
		return deliveryReadback(d, root, id, workflowError(4, "startup entrypoint differs from frozen contract; use its callback guide"))
	}
	return WorkflowDeliveryBegin(d, root, id, contextPath, in)
}
