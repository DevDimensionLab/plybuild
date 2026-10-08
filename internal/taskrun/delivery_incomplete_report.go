package taskrun

import (
	"os"
	"path/filepath"
)

// This deliberately does not call workflowFresh or deliveryCallback. Reporting
// an unavailable runtime dependency must not require that dependency. It checks
// identity and immutable history, never effective permission for Task effects.
func deliveryIncompleteReportCallback(d Dependencies, s workflowState, contextPath string) error {
	if !deliveryRun(s.Request) || s.Request.Envelope != deliveryEnv("herdr-run-request") || s.Result.Delivery == nil {
		return deliveryContinuityError("report_unsupported", "incomplete reporting requires the supported delivery-owner run contract")
	}
	if s.Result.Delivery.OwnershipRelease != nil || s.Result.Delivery.Phase == "completed" {
		return deliveryContinuityError("report_owner", "the original delivery owner no longer owns this active mandate")
	}
	if contextPath == "" || contextPath != s.Result.Paths.Context || s.ContextSHA256 == "" {
		return deliveryContinuityError("report_context", "incomplete reporting requires the exact original private context")
	}
	cwd, err := d.CWD()
	if err != nil {
		return err
	}
	if cwd != s.Observed.Target.WorktreeLocator {
		return deliveryContinuityError("report_owner", "incomplete reporting requires the original owner's exact Task cwd")
	}
	if err = physical(cwd, false); err != nil {
		return err
	}
	target, err := d.Workspace.IntegrationGit.ObserveIntegrationWorktree(cwd, s.Observed.Target.Ref)
	if err != nil {
		return err
	}
	if target.Ref != s.Observed.Target.Ref || target.GitCommonDir != s.Observed.Target.GitCommonDir || target.Locator != cwd {
		return deliveryContinuityError("report_owner", "original Task identity changed; the report cannot substitute another owner")
	}
	runtime, err := workflowEffectiveRuntimeMode(s, false)
	if err != nil {
		return err
	}
	if err = verifyExecutable(runtime.PlyExecutable); err != nil {
		return deliveryContinuityError("report_integrity", "preserved callback control is missing or changed: "+err.Error())
	}
	actual, err := d.Executable()
	if err == nil {
		actual, err = filepath.EvalSymlinks(actual)
	}
	if err != nil {
		return err
	}
	if err = physical(actual, false); err != nil {
		return err
	}
	info, err := os.Stat(actual)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || info.Mode().Perm()&0111 == 0 {
		return deliveryContinuityError("report_integrity", "incomplete reporter must be the actual installed executable")
	}
	if err = workflowFreshArtifacts(d, s, false, false); err != nil {
		return deliveryContinuityError("report_integrity", "required incomplete-report history: "+err.Error())
	}
	return deliveryLiveOwner(d, s)
}

func deliveryIncompleteReadback(s workflowState) (WorkflowAction, bool) {
	d := s.Result.Delivery
	if d == nil || d.Phase != "stopped" && d.Phase != "needs_input" || len(d.Events) == 0 {
		return WorkflowAction{}, false
	}
	event := d.Events[len(d.Events)-1]
	if event.Kind != "report" || d.LastEventSHA256 == nil || *d.LastEventSHA256 != event.Binding.SHA256 {
		return WorkflowAction{}, false
	}
	raw, err := workflowBound(event.Binding, 1<<20)
	if err != nil {
		return WorkflowAction{}, false
	}
	var report DeliveryReport
	if decode(raw, 1<<20, &report) != nil || validateDeliveryReportContent(report) != nil || report.EventID != event.ID || report.Phase != d.Phase || workflowClaim(s, report.RunID, report.RequestSHA256, report.SessionID) != nil {
		return WorkflowAction{}, false
	}
	message := report.Meaning
	if report.Phase == "needs_input" {
		message = *report.Question
	}
	return WorkflowAction{"user", message}, true
}
