// Package currentstate interprets current Run/Worker/Week observations. It has
// no stores or executors: reading an interpretation never consumes its facts.
package currentstate

type Action struct {
	Actor   string `json:"actor"`
	Message string `json:"message"`
}

const (
	ProviderWait          = "provider_wait"
	Blocked               = "blocked"
	ProductChoiceRequired = "One product choice is required; the same developer session and Task reservation are retained."
)

type Worker struct {
	Phase, Stage, Reason, ReasonCode, RetryNotBefore, ControlAction string
}

// Answer contains positively observed continuation facts, not instructions to
// dispatch. Detail is only a diagnostic; it never selects a transition.
type Answer struct {
	State, RequestID, Detail, ParentDetail string
}

type Run struct {
	Worker      *Worker
	Answer      *Answer
	Completed   bool
	HasReasons  bool
	SavedAction Action
}

// RunAction owns precedence and actor selection for all current worker reads.
// Fresh controls, concrete obstacles and parent parking win over answer advice.
func RunAction(r Run) Action {
	a, w := r.SavedAction, r.Worker
	if w == nil || r.HasReasons {
		return a
	}
	wait := w.ReasonCode == ProviderWait || w.ReasonCode == "" && LegacyProviderWait(*w)
	if w.Phase == "blocked" || w.Reason != "" && !wait {
		return a
	}
	// Other actors carry native delivery/startup/closeout decisions. Their
	// precedence is established by the read adapter before interpreting a worker.
	if a.Actor != "human" && a.Actor != "worker_wait" && a.Actor != "worker" && a.Actor != "product_answer" && a.Actor != "parent_control" {
		return a
	}
	if (w.ControlAction == "resume" && (w.Phase == "paused" || w.Phase == "stopped" || w.Phase == "running" || w.Phase == "waiting")) ||
		(w.ControlAction == "pause" && (w.Phase == "running" || w.Phase == "waiting")) ||
		(w.ControlAction == "stop" && (w.Phase == "running" || w.Phase == "waiting" || w.Phase == "paused")) {
		return Action{"worker", w.ControlAction + " requested for this same Run; consumption at a safe boundary is not yet confirmed. Existing attempts retain ownership until their completion is known."}
	}
	if w.Phase == "paused" || w.Phase == "stopped" {
		if r.Answer != nil && a.Actor == "product_answer" {
			return answerAction(r)
		}
		return a
	}
	// A later reviewer/provider wait is current; confirmed answer use is
	// history. An unresolved answer outcome still requires inspection.
	if w.Phase == "waiting" && w.RetryNotBefore != "" && (w.Reason == "" || wait) && (r.Answer == nil || r.Answer.State == "continuation_completed") {
		return Action{"worker_wait", "The same Run waits until " + w.RetryNotBefore + "; continuing still requires confirmed attempt completion."}
	}
	if r.Answer != nil {
		return answerAction(r)
	}
	switch w.Phase {
	case "awaiting_product_answer":
		return Action{"product_answer", "The same developer session and Task await one bound product answer; no model turn is dispatched while waiting."}
	case "running":
		if w.Reason == "" && w.RetryNotBefore == "" {
			return Action{"worker", "Observe the same Run; cached worker state does not attest provider liveness or attempt completion."}
		}
	}
	return a
}

func answerAction(r Run) Action {
	f, w := r.Answer, r.Worker
	inspect := func(detail string) Action {
		return Action{"human", detail + " Do not answer again or resend; inspect the preserved Run."}
	}
	if f.State == "unanswered" {
		if w.Phase != "awaiting_product_answer" {
			return inspect("The product question has no recorded answer while the worker is " + w.Phase + ".")
		}
		return Action{"product_answer", ProductChoiceRequired}
	}
	if f.RequestID == "" {
		return inspect(f.Detail)
	}
	answer := "Answer " + f.RequestID
	switch f.State {
	case "answer_recorded":
		return Action{"worker", answer + " is received; the worker owns its single same-session continuation, which has not started. Answer use is not confirmed; no new human answer is needed." + f.ParentDetail}
	case "continuation_deferred":
		return Action{"parent_control", answer + " is received; its continuation waits for explicit parent resume. Answer use is not confirmed; no new human answer is needed."}
	case "continuation_reserved":
		return Action{"worker", answer + " is received; the worker has reserved its single same-session continuation. Answer use is not confirmed; no new human answer is needed." + f.ParentDetail}
	case "continuation_observed":
		return Action{"worker", answer + " is received and the same developer session is continuing with it. Answer use is not confirmed until that turn completes; no new human answer is needed." + f.ParentDetail}
	case "continuation_completed":
		if r.Completed || w.Phase == "completed" || w.Stage == "completed" {
			return Action{"user", answer + " use is confirmed by the same developer session and the Task delivery is completed. No human answer or further action is needed for this answer."}
		}
		return Action{"worker", answer + " use is confirmed by the same developer session; the worker owns the " + w.Stage + " stage of this Run. No human answer is needed." + f.ParentDetail}
	}
	return inspect(answer + " is received, but its continuation is not confirmed: " + f.Detail)
}
