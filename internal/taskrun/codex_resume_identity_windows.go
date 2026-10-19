package taskrun

import "time"

// Windows keeps the hook identity requirement; the local Unix lock observer
// cannot establish ownership there. No input is sent without a verified ID.
func observeCodexResumeIdentity(d Dependencies, s workflowState, a workflowAgent, deadline time.Time) (*CodexResumeIdentity, error) {
	return nil, workflowError(5, "local Codex thread-lock identity inspection is unavailable on Windows; await Herdr session metadata or retire the same attempt")
}
