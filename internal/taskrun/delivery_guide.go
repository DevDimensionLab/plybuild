package taskrun

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// These files describe input, not evidence. Runtime facts and human judgments
// stay unobserved until the recipient fills a copy from actual observations.
func workflowDeliveryGuide(d Dependencies, s workflowState) error {
	if !deliveryRun(s.Request) {
		return nil
	}
	dir := filepath.Join(workflowDeliveryGuideDirectory(s), "templates")
	a := DeliveryAcceptance{
		Acceptance:         Acceptance{Envelope: deliveryEnv("run-acceptance"), RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID, RuntimeClaim: RuntimeClaim{NativeSessionID: ptr(s.Result.Transport.AgentSessionID)}, Sandbox: json.RawMessage("null"), Issues: json.RawMessage("[]")},
		DeliveryPermission: DeliveryPermissionAcceptance{LaunchContractSHA256: s.Request.Runtime.PermissionBinding.EffectivePolicySHA256, ActualPolicyEvidence: []Evidence{}},
	}
	if s.Request.Delivery.Agreement != nil {
		a.DeliveryPermission.DeliveryAgreementSHA256 = workspace.DeliveryAgreementDigest(*s.Request.Delivery.Agreement)
	}
	r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID, Phase: "working", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	for _, f := range []struct {
		name  string
		value any
	}{{"acceptance.json", a}, {"report.json", r}} {
		if mandateReceipt(s.Request) && f.name == "acceptance.json" {
			continue
		}
		if _, err := workflowKeep(d, filepath.Join(dir, f.name), f.value); err != nil {
			return err
		}
	}
	guide := fmt.Sprintf("# Delivery observations\n\nRun: %s\nSource: %s\n\nPreserve actual command results and resource identities. Choose reviews, testing, integration, installation and cleanup explicitly. Observing a Run never closes a pane or removes a worktree.\n", s.Result.RunID, s.Observed.Target.WorktreeLocator)
	return d.writeOnce(filepath.Join(workflowDeliveryGuideDirectory(s), "callback-guide.md"), []byte(guide))
}
