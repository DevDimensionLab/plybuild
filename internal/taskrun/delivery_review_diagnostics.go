package taskrun

import (
	"errors"
	"os"
	"path/filepath"
)

// Only facts already classified by HerdrCallError cross this boundary. Child
// messages, stderr, terminal contents and attempted argv never enter readback.
type deliveryReviewFailure struct {
	Phase        string `json:"phase"`
	Class        string `json:"class"`
	Exit         string `json:"exit"`
	ProviderCode string `json:"provider_code"`
	StdoutBytes  int64  `json:"stdout_bytes"`
	StderrBytes  int64  `json:"stderr_bytes"`
	Truncated    bool   `json:"truncated"`
}

type deliveryReviewObservation struct {
	Outcome              string                 `json:"outcome"`
	TaskPromptStage      string                 `json:"task_prompt_stage,omitempty"`
	BootstrapPromptStage string                 `json:"bootstrap_prompt_stage,omitempty"`
	FirstFailure         *deliveryReviewFailure `json:"first_failure,omitempty"`
	SessionID            string                 `json:"last_observed_session_id,omitempty"`
	AgentStatus          string                 `json:"last_observed_agent_status,omitempty"`
	Failure              *deliveryReviewFailure `json:"failure,omitempty"`
}

type DeliveryReviewStartupDiagnostics struct {
	Identity             *deliveryReviewPane        `json:"identity,omitempty"`
	StartAttempt         string                     `json:"start_attempt"`
	SessionBinding       string                     `json:"session_binding"`
	BootstrapPrompt      string                     `json:"bootstrap_prompt"`
	TaskPrompt           string                     `json:"task_prompt"`
	ProviderLiveness     string                     `json:"provider_liveness"`
	FirstFailure         *deliveryReviewFailure     `json:"first_failure,omitempty"`
	FirstFailureEvidence *FileBinding               `json:"first_failure_evidence,omitempty"`
	LastObservation      *deliveryReviewObservation `json:"last_observation,omitempty"`
	ObservationEvidence  *FileBinding               `json:"observation_evidence,omitempty"`
	InspectionArgv       [][]string                 `json:"inspection_argv,omitempty"`
}

func deliveryReviewFailureFacts(err error) *deliveryReviewFailure {
	var call *HerdrCallError
	if !errors.As(err, &call) {
		return nil
	}
	return &deliveryReviewFailure{call.Phase, call.Class, call.Exit, call.ProviderCode, call.StdoutBytes, call.StderrBytes, call.Truncated}
}

func deliveryReviewDiagnostics(s workflowState, dir string, start deliveryReviewStart) *DeliveryReviewStartupDiagnostics {
	out := &DeliveryReviewStartupDiagnostics{Identity: start.Pane, StartAttempt: "not_recorded", SessionBinding: "not_recorded", BootstrapPrompt: "not_recorded", TaskPrompt: "not_recorded", ProviderLiveness: "unknown"}
	// Missing old evidence is unknown, never an invented outcome or cause.
	for name, dest := range map[string]*string{"start-attempt.json": &out.StartAttempt, "bootstrap-attempt.json": &out.BootstrapPrompt, "prompt-attempt.json": &out.TaskPrompt} {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			*dest = "attempted"
		}
	}
	if start.Agent != nil {
		out.SessionBinding = "bound"
	}
	// New rounds record an explicit pre-prompt observation; only those rounds
	// can distinguish no prompt attempt from incomplete historical evidence.
	read := func(name string, v any) *FileBinding {
		path := filepath.Join(dir, name)
		raw, err := readFile(path, 1<<20, true)
		if err != nil || decode(raw, 1<<20, v) != nil {
			return nil
		}
		return &FileBinding{path, hash(raw)}
	}
	var failure deliveryReviewFailure
	if binding := read("start-failure.json", &failure); binding != nil {
		out.FirstFailure = &failure
		out.FirstFailureEvidence = binding
	}
	var obs deliveryReviewObservation
	if binding := read("session-observation.json", &obs); binding != nil {
		if start.Agent == nil && obs.TaskPromptStage == "not_attempted" {
			out.SessionBinding = "not_bound"
		}
		out.LastObservation = &obs
		out.ObservationEvidence = binding
		if out.TaskPrompt == "not_recorded" && obs.TaskPromptStage == "not_attempted" {
			out.TaskPrompt = "not_attempted"
		}
		if out.BootstrapPrompt == "not_recorded" && obs.BootstrapPromptStage == "not_attempted" {
			out.BootstrapPrompt = "not_attempted"
		}
	}
	if start.State == "started" {
		out.TaskPrompt = "sent"
	}
	if start.Pane != nil {
		exe := s.Request.Herdr.Executable.Path
		pane := start.Pane.PaneID
		out.InspectionArgv = [][]string{{exe, "agent", "get", pane}, {exe, "pane", "get", pane}, {exe, "pane", "process-info", "--pane", pane}}
	}
	return out
}
