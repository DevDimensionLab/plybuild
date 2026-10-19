package currentstate

type Origin string

const (
	OtherOrigin Origin = "other"
	NativeRun   Origin = "native_run"
	AnswerBlock Origin = "answer_continuation"

	ProductAnswerWait = "Wait for one bound product answer on the same native Task and developer session"
	NativeNextPrefix  = "Native Run next actor "
	WeekStartSelected = "Start only the selected native Goals"
	WeekQueueComplete = "Selected queue complete; no additional work is authorized"
	WeekRunCompleted  = "Native Run completed; the week controller selects the next eligible Goal"
)

type Week struct {
	Phase, DesiredAction, NextAction string
	GuidanceOrigin                   Origin
	OtherEntriesComplete             bool
	UnresolvedDispatch               bool
	ControllerBlocked                bool
	ReconciliationConfirmed          bool
}

type Entry struct {
	State       string
	BlockOrigin Origin
}

type Native struct {
	Phase       string
	Completed   bool
	Action      Action
	Unavailable bool
	Failure     string
}

// WeekEntry projects the same observed Run in Show and Tick. Only this
// interpretation creates and clears answer-derived blocks. Ordinary blocks
// (including failed queue release) remain owned by their execution path.
func WeekEntry(s Week, e Entry, r Native) (Week, Entry) {
	if s.UnresolvedDispatch || s.ControllerBlocked && !s.ReconciliationConfirmed {
		return s, e
	}
	if r.Unavailable {
		s.Phase, e.State = "unknown", "unknown"
		s.NextAction = "Native Run observation failed: " + r.Failure + "; inspect the same Run before answering"
		return s, e
	}
	if r.Completed {
		e.State, e.BlockOrigin = "completed", OtherOrigin
		if s.DesiredAction == "continue" {
			s.Phase, s.NextAction, s.GuidanceOrigin = "running", WeekRunCompleted, NativeRun
			if s.OtherEntriesComplete {
				s.Phase, s.NextAction = "completed", WeekQueueComplete
			}
		}
		return s, e
	}
	if e.State == "completed" || e.State == "blocked" && e.BlockOrigin != AnswerBlock {
		return s, e
	}
	if s.DesiredAction != "continue" {
		return s, e
	}
	if r.Phase == "awaiting_product_answer" && (r.Action.Actor == "product_answer" || r.Action.Actor == "") {
		s.Phase, e.State, e.BlockOrigin = "awaiting_product_answer", "awaiting_product_answer", OtherOrigin
		s.NextAction, s.GuidanceOrigin = ProductAnswerWait, NativeRun
		return s, e
	}
	if r.Phase == "awaiting_product_answer" && r.Action.Actor == "human" {
		s.Phase, e.State, e.BlockOrigin = "unknown", "blocked", AnswerBlock
		s.NextAction, s.GuidanceOrigin = nativeAction(r.Action), NativeRun
		return s, e
	}
	if e.BlockOrigin == AnswerBlock || e.State == "awaiting_product_answer" {
		e.State, e.BlockOrigin = "running", OtherOrigin
	}
	s.Phase = r.Phase
	if s.Phase == "awaiting_product_answer" {
		// A recorded answer with a worker/parent actor has left the unanswered
		// wait even while the provider's cached phase still names that wait.
		s.Phase = "running"
	}
	if s.GuidanceOrigin == NativeRun || r.Action.Actor != "" {
		s.NextAction, s.GuidanceOrigin = nativeAction(r.Action), NativeRun
	}
	return s, e
}

func nativeAction(a Action) string {
	if a.Actor == "" {
		return "Observe the same native Run"
	}
	return NativeNextPrefix + a.Actor + ": " + a.Message
}

func WeekQueue(s Week) Week {
	if s.ControllerBlocked {
		return s
	}
	if s.GuidanceOrigin == NativeRun {
		s.NextAction, s.GuidanceOrigin = WeekStartSelected, OtherOrigin
		if s.Phase == "completed" {
			s.NextAction = WeekQueueComplete
		}
	}
	return s
}
