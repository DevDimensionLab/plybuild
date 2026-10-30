package humaninput

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// observe accepts only the physical Ply workspace root. It reads the marker
// and never discovers a workspace from cwd.
func observe(d Dependencies, root string) error {
	if absoluteInput(root) != nil {
		return fail(CodeInvalidArguments, "--workspace must name the physical absolute normalized workspace root")
	}
	if ws, e := workspace.ObserveRoot(d.Workspace, root); e != nil || ws.Root != root {
		return fail(CodeInvalidArguments, "--workspace is not a physical Ply workspace root")
	}
	return nil
}

type registrations struct {
	projects []workspace.ProjectRecord
	repos    []workspace.RepoRecord
	registry workspace.WorkItemRegistry
}

func readRegistrations(d Dependencies, root string) (registrations, error) {
	var r registrations
	var e error
	if d.Workspace.Projects == nil || d.Workspace.WorkItems == nil {
		return r, invalidState("workspace registration readers are unavailable")
	}
	if r.projects, r.repos, e = d.Workspace.Projects.Snapshot(root); e != nil {
		return r, invalidState("project registrations cannot be read")
	}
	if r.registry, e = d.Workspace.WorkItems.SnapshotRegistrations(root); e != nil {
		return r, invalidState("work-item registrations cannot be read")
	}
	return r, nil
}

// checkScope validates the project, optional repository and optional existing
// Task against the registers. It never creates or rewrites a registration.
func (r registrations) checkScope(s Scope) (*workspace.ProjectRecord, string) {
	var project *workspace.ProjectRecord
	for i := range r.projects {
		if string(r.projects[i].ID) == s.ProjectID {
			project = &r.projects[i]
		}
	}
	if project == nil {
		return nil, "project " + s.ProjectID + " is not registered in this workspace"
	}
	if s.RepoID != nil {
		found := false
		for _, id := range project.RepoIDs {
			found = found || string(id) == *s.RepoID
		}
		if !found {
			return nil, "repository " + *s.RepoID + " does not belong to project " + s.ProjectID
		}
	}
	if s.TaskID != nil {
		if s.RepoID == nil {
			return nil, "a Task question also names the Task's repository"
		}
		found := false
		for _, t := range r.registry.Tasks {
			if string(t.ID) == *s.TaskID {
				found = true
				if string(t.ProjectID) != s.ProjectID || string(t.RepoID) != *s.RepoID {
					return nil, "Task " + *s.TaskID + " belongs to another project or repository"
				}
			}
		}
		if !found {
			return nil, "Task " + *s.TaskID + " is not registered; a question does not create a Task"
		}
	}
	return project, ""
}

func within(path, base string) bool {
	return base != "" && (path == base || strings.HasPrefix(path, strings.TrimSuffix(base, "/")+"/"))
}

// checkOrigin requires the agent's physical cwd to belong to the project: its
// wrapper, repositories, companions or registered Epic/Task worktrees.
func (r registrations) checkOrigin(d Dependencies, root string, project *workspace.ProjectRecord, origin string) string {
	if absoluteInput(origin) != nil {
		return "origin_cwd must be a physical absolute normalized directory"
	}
	if info, e := os.Lstat(origin); e != nil || !info.IsDir() {
		return "origin_cwd must be an existing directory"
	}
	bases := []string{project.Wrapper}
	for _, repo := range r.repos {
		for _, id := range project.RepoIDs {
			if repo.ID == id {
				bases = append(bases, repo.Locator)
			}
		}
	}
	if meta, e := workspace.ReadProjectMetadata(d.Workspace, root); e == nil {
		for _, c := range meta.Companions(project.ID) {
			bases = append(bases, c.Locator)
		}
	}
	for _, epic := range r.registry.Epics {
		if epic.ProjectID == project.ID {
			for _, b := range epic.RepoBindings {
				bases = append(bases, b.Worktree.Locator)
			}
		}
	}
	for _, t := range r.registry.Tasks {
		if t.ProjectID == project.ID && t.Worktree != nil {
			bases = append(bases, t.Worktree.Locator)
		}
	}
	for _, base := range bases {
		if within(origin, base) {
			return ""
		}
	}
	return "origin_cwd does not belong to project " + string(project.ID) + "'s wrapper, repositories, companions or worktrees"
}

// AskDraft creates, or finds again, the one question for this publication key.
// The same key with the same content is idempotent; other content conflicts.
func AskDraft(d Dependencies, root string, draft QuestionDraft) (Question, error) {
	if e := observe(d, root); e != nil {
		return Question{}, e
	}
	if e := validateDraft(draft); e != nil {
		return Question{}, e
	}
	regs, e := readRegistrations(d, root)
	if e != nil {
		return Question{}, e
	}
	scope := Scope{ProjectID: draft.ProjectID, RepoID: draft.RepoID, TaskID: draft.TaskID, ActivityID: draft.ActivityID, ActivityTitle: draft.ActivityTitle}
	project, problem := regs.checkScope(scope)
	if problem == "" {
		problem = regs.checkOrigin(d, root, project, draft.OriginCWD)
	}
	if problem != "" {
		return Question{}, fail(CodeScopeInvalid, problem)
	}
	record := questionRecord{Workspace: root, QuestionID: questionIDFor(root, draft.ProjectID, draft.PublicationKey), PublicationKey: draft.PublicationKey, Scope: scope, Source: Source{draft.Provider, draft.SessionID, draft.OriginCWD, ProvenanceAgent}, Text: draft.Text, Reason: draft.Reason, AnswerKind: draft.AnswerKind, Options: draft.Options}
	record.QuestionSHA256 = questionDigest(record)
	existing := func() (Question, bool, error) {
		q, e := view(root, record.QuestionID)
		if errorCode(e) == CodeNotFound {
			return q, false, nil
		}
		if e != nil {
			return q, true, e
		}
		if q.QuestionSHA256 != record.QuestionSHA256 {
			return q, true, failOn(q, CodePublicationConflict, "this publication key already names a different question; publish a changed question under a new key")
		}
		return q, true, nil
	}
	if q, found, e := existing(); found {
		return q, e
	}
	var out Question
	e = withLock(root, func() error {
		q, found, e := existing()
		if found {
			out = q
			return e
		}
		record.CreatedAtUTC = d.now()
		if e = publish(questionPath(root, record.QuestionID), record); e != nil {
			return e
		}
		out, e = view(root, record.QuestionID)
		return e
	})
	if e != nil && errorCode(e) == "" {
		return out, fail(CodeOutcomeUnknown, "publication did not complete: "+e.Error()+"; ask again with the same publication key")
	}
	return out, e
}

func errorCode(e error) string {
	var x *Error
	if errors.As(e, &x) {
		return x.Code
	}
	return ""
}

var states = map[string]bool{StateWaiting: true, StateAnswered: true, StateReceived: true, StateContinued: true, StateWithdrawn: true}

// List returns one page of a project's questions in immutable ID order. A
// failed scan returns no page; it never implies a missing question is gone.
func List(d Dependencies, root string, in ListQuery) (ListResult, error) {
	out := ListResult{Kind: ListKind, SchemaVersion: SchemaVersion, Mode: Mode, Workspace: root, ProjectID: in.ProjectID, State: in.State, Limit: in.Limit, Cursor: in.Cursor, Questions: []Question{}}
	if out.Limit == 0 {
		out.Limit = DefaultListLimit
	}
	switch {
	case !projectPattern.MatchString(in.ProjectID):
		return out, fail(CodeInvalidArguments, "--project must be a project ID")
	case out.Limit < 1 || out.Limit > MaxListLimit:
		return out, fail(CodeInvalidArguments, "--limit must be between 1 and 20")
	case in.State != nil && !states[*in.State]:
		return out, fail(CodeInvalidArguments, "--state must be waiting_for_answer, answered, agent_received, agent_reported_continued or withdrawn")
	case in.Cursor != nil && !questionIDPattern.MatchString(*in.Cursor):
		return out, fail(CodeInvalidArguments, "--cursor must be a next_cursor value returned by list")
	}
	if e := observe(d, root); e != nil {
		return out, e
	}
	ids, e := questionIDs(root)
	if e != nil {
		return out, e
	}
	for _, id := range ids {
		if in.Cursor != nil && id <= *in.Cursor {
			continue
		}
		q, e := view(root, id)
		if e != nil {
			return out, e
		}
		if q.Scope.ProjectID != in.ProjectID || (in.State != nil && q.State != *in.State) {
			continue
		}
		if len(out.Questions) == out.Limit {
			out.HasMore = true
			out.NextCursor = ptr(out.Questions[len(out.Questions)-1].QuestionID)
			break
		}
		out.Questions = append(out.Questions, q)
	}
	out.Complete = !out.HasMore
	out.ObservedAtUTC = d.now()
	return out, nil
}

func checkQuestionArgs(d Dependencies, root, id string) error {
	if !questionIDPattern.MatchString(id) {
		return fail(CodeInvalidArguments, "expected one exact hiq_ question ID")
	}
	return observe(d, root)
}

func Show(d Dependencies, root, id string) (Question, error) {
	if e := checkQuestionArgs(d, root, id); e != nil {
		return Question{}, e
	}
	return view(root, id)
}

func validateAnswer(q Question, v AnswerValue) error {
	bad := func(detail string) error { return failOn(q, CodeInvalidAnswer, detail) }
	if v.AnswerKind != q.AnswerKind {
		return bad("this question takes a " + q.AnswerKind + " answer")
	}
	if q.AnswerKind == AnswerChoice {
		if v.OptionID == nil || v.Text != nil {
			return bad("a choice answer has option_id and null text")
		}
		for _, o := range q.Options {
			if o.ID == *v.OptionID {
				return nil
			}
		}
		return bad("option " + *v.OptionID + " is not one of this question's options")
	}
	if v.Text == nil || v.OptionID != nil {
		return bad("a text answer has text and null option_id")
	}
	if !ValidAnswerText(*v.Text) {
		return bad("text must be non-blank UTF-8 without NUL, at most 4000 codepoints and 16000 bytes")
	}
	return nil
}

func openForAnswer(q Question) error {
	switch q.State {
	case StateWaiting:
		return nil
	case StateWithdrawn:
		return failOn(q, CodeWithdrawn, "the agent withdrew this question; it takes no answer")
	default:
		return failOn(q, CodeAlreadyAnswered, "this question already has its one recorded answer")
	}
}

func scopeAvailable(d Dependencies, root string, q Question) error {
	regs, e := readRegistrations(d, root)
	if e != nil {
		return e
	}
	if _, problem := regs.checkScope(q.Scope); problem != "" {
		return failOn(q, CodeScopeUnavailable, "the question's scope is no longer registered ("+problem+"); history is kept but it takes no new answer")
	}
	return nil
}

// PreviewAnswer returns the exact immutable request for one answer. It writes
// nothing; the caller keeps the request bytes and applies them unchanged.
func PreviewAnswer(d Dependencies, root, id, requestID string, v AnswerValue) (Preview, error) {
	if e := checkQuestionArgs(d, root, id); e != nil {
		return Preview{}, e
	}
	if !requestIDPattern.MatchString(requestID) {
		return Preview{}, fail(CodeInvalidArguments, "--request-id must be 8-128 ASCII letters, digits, dots, underscores or hyphens")
	}
	q, e := view(root, id)
	if e != nil {
		return Preview{}, e
	}
	if e = openForAnswer(q); e != nil {
		return Preview{}, e
	}
	if e = validateAnswer(q, v); e != nil {
		return Preview{}, e
	}
	if e = scopeAvailable(d, root, q); e != nil {
		return Preview{}, e
	}
	r := AnswerRequest{Kind: RequestKind, SchemaVersion: SchemaVersion, Mode: Mode, Purpose: PurposeAnswer, Workspace: root, QuestionID: id, QuestionSHA256: q.QuestionSHA256, Scope: q.Scope, RequestID: requestID, Answer: v}
	r.BindingSHA256 = bindingOf(r)
	return Preview{Kind: PreviewKind, SchemaVersion: SchemaVersion, Mode: Mode, Workspace: root, Question: q, Request: r, RequestSHA256: digestOf(r), ObservedAtUTC: d.now()}, nil
}

func statusBase(d Dependencies, root, id, requestID string) AnswerStatus {
	return AnswerStatus{Kind: StatusKind, SchemaVersion: SchemaVersion, Mode: Mode, Workspace: root, QuestionID: id, RequestID: requestID, State: StatusUnknown, Reasons: []Reason{}, ObservedAtUTC: d.now()}
}
func reasonOf(e error) Reason {
	var x *Error
	if errors.As(e, &x) {
		return Reason{x.Code, x.Detail}
	}
	return Reason{CodeOutcomeUnknown, "The outcome could not be read reliably; keep the same request and read its status."}
}

// classify reads the authoritative outcome of (question, request) from a view.
func classify(s AnswerStatus, q Question, e error) AnswerStatus {
	if e != nil {
		s.State = StatusUnknown
		if errorCode(e) == CodeNotFound {
			s.State = StatusUnavailable
		}
		s.Reasons = append(s.Reasons, reasonOf(e))
		return s
	}
	s.Question = &q
	if q.Answer == nil {
		s.State = StatusNotRecorded
		return s
	}
	s.Winner = &WinnerRef{q.Answer.RequestID, q.Answer.RequestSHA256}
	if q.Answer.RequestID == s.RequestID {
		s.State, s.RequestSHA256 = StatusRecorded, &q.Answer.RequestSHA256
	} else {
		s.State = StatusOtherWinner
	}
	return s
}

// Status reads the outcome of one request without effects. not_recorded means
// the question is readable and this request is not its answer; unavailable and
// unknown mean that nothing reliable is known and open no new write.
func Status(d Dependencies, root, id, requestID string) (AnswerStatus, error) {
	s := statusBase(d, root, id, requestID)
	if !requestIDPattern.MatchString(requestID) {
		return s, fail(CodeInvalidArguments, "expected one stable request ID")
	}
	if e := checkQuestionArgs(d, root, id); e != nil {
		if errorCode(e) == CodeInvalidArguments && questionIDPattern.MatchString(id) {
			s.State = StatusUnavailable
			s.Reasons = append(s.Reasons, reasonOf(e))
			return s, nil
		}
		return s, e
	}
	q, e := view(root, id)
	return classify(s, q, e), nil
}

// Apply records the request as the question's one winner in the explicit
// workspace root, which must equal the request's immutable workspace. Replay of
// the same request is idempotent; any other request is refused without
// overwriting.
func Apply(d Dependencies, root string, r AnswerRequest) (AnswerStatus, error) {
	s := statusBase(d, root, r.QuestionID, r.RequestID)
	if e := validateRequestShape(r); e != nil {
		s.State = StatusNotRecorded
		s.Reasons = append(s.Reasons, reasonOf(e))
		return s, e
	}
	if root != r.Workspace {
		e := fail(CodeRequestMismatch, "--workspace differs from the request's immutable workspace; nothing was read or changed")
		s.State = StatusNotRecorded
		s.Reasons = append(s.Reasons, reasonOf(e))
		return s, e
	}
	sha := digestOf(r)
	if e := observe(d, r.Workspace); e != nil {
		s.State = StatusUnavailable
		s.Reasons = append(s.Reasons, reasonOf(e))
		return s, e
	}
	id := r.QuestionID
	// decide returns a final status, or (zero, nil) when this request may win.
	decide := func() (AnswerStatus, error) {
		q, e := view(root, id)
		out := classify(s, q, e)
		switch out.State {
		case StatusRecorded:
			if *out.RequestSHA256 != sha {
				e := fail(CodeRequestConflict, "this request ID already names different request bytes")
				out.State, out.RequestSHA256 = StatusOtherWinner, nil
				out.Reasons = append(out.Reasons, reasonOf(e))
				return out, e
			}
			return out, nil
		case StatusOtherWinner:
			e := fail(CodeAlreadyAnswered, "another request is this question's recorded answer; it is not overwritten")
			out.Reasons = append(out.Reasons, reasonOf(e))
			return out, e
		case StatusNotRecorded:
		default:
			return out, errors.New(out.Reasons[0].Detail)
		}
		refuse := func(e error) (AnswerStatus, error) {
			out.Reasons = append(out.Reasons, reasonOf(e))
			return out, e
		}
		if e := openForAnswer(q); e != nil {
			return refuse(e)
		}
		if r.QuestionSHA256 != q.QuestionSHA256 || !sameScope(r.Scope, q.Scope) {
			return refuse(failOn(q, CodeRequestMismatch, "the request was prepared for different question content or scope"))
		}
		if e := validateAnswer(q, r.Answer); e != nil {
			return refuse(e)
		}
		if e := scopeAvailable(d, root, q); e != nil {
			return refuse(e)
		}
		return AnswerStatus{}, nil
	}
	if out, e := decide(); out.Kind != "" {
		return out, e
	}
	var out AnswerStatus
	e := withLock(root, func() error {
		var e error
		if out, e = decide(); out.Kind != "" {
			return e
		}
		if e = d.fault("before_answer_publish"); e != nil {
			return e
		}
		if e = publish(answerPath(root, id), answerRecord{Request: r, RequestSHA256: sha, RecordedAtUTC: d.now()}); e != nil {
			return e
		}
		return d.fault("after_answer_publish")
	})
	if e == nil {
		q, e := view(root, id)
		return classify(s, q, e), nil
	}
	if out.Kind != "" {
		return out, e
	}
	// An interrupted write reports only what a fresh read proves.
	q, ve := view(root, id)
	if observed := classify(s, q, ve); observed.State == StatusRecorded || observed.State == StatusOtherWinner {
		return observed, e
	}
	s.Reasons = append(s.Reasons, Reason{CodeOutcomeUnknown, "The apply attempt was interrupted; read status for the same request before any new attempt."})
	return s, e
}

// Wait polls without a lock until the question leaves waiting_for_answer or
// the technical timeout ends. Timeout means still waiting, never a deadline.
func Wait(d Dependencies, root, id string, timeout time.Duration) (WaitResult, error) {
	if e := checkQuestionArgs(d, root, id); e != nil {
		return WaitResult{}, e
	}
	poll := d.Poll
	if poll <= 0 {
		poll = 500 * time.Millisecond
	}
	deadline := time.Now().Add(timeout)
	for {
		q, e := view(root, id)
		if e != nil {
			return WaitResult{}, e
		}
		out := WaitResult{Kind: WaitKind, SchemaVersion: SchemaVersion, Mode: Mode, Workspace: root, QuestionID: id, Outcome: WaitStillWaiting, TimeoutSeconds: int(timeout / time.Second), Question: q}
		switch q.State {
		case StateWaiting:
		case StateWithdrawn:
			out.Outcome = WaitWithdrawn
		default:
			out.Outcome = WaitAnswered
		}
		if out.Outcome != WaitStillWaiting || !time.Now().Before(deadline) {
			out.ObservedAtUTC = d.now()
			return out, nil
		}
		sleep := poll
		if left := time.Until(deadline); left < sleep {
			sleep = left
		}
		time.Sleep(sleep)
	}
}

// transition decides on a fresh view, first without and then under the store
// lock. decide returns the fact to publish, or a nil fact when the current
// question already holds it (idempotent replay).
func transition(d Dependencies, root, id string, decide func(Question) (any, string, error)) (Question, error) {
	if e := checkQuestionArgs(d, root, id); e != nil {
		return Question{}, e
	}
	q, e := view(root, id)
	if e != nil {
		return q, e
	}
	fact, _, e := decide(q)
	if e != nil || fact == nil {
		return q, e
	}
	e = withLock(root, func() error {
		if q, e = view(root, id); e != nil {
			return e
		}
		fact, path, e := decide(q)
		if e != nil || fact == nil {
			return e
		}
		if e = publish(path, fact); e != nil {
			return e
		}
		q, e = view(root, id)
		return e
	})
	if e != nil && errorCode(e) == "" {
		return q, fail(CodeOutcomeUnknown, "the report did not complete: "+e.Error()+"; show the question before repeating it")
	}
	return q, e
}

func matchesWinner(q Question, requestID, sha string) error {
	switch {
	case q.State == StateWithdrawn:
		return failOn(q, CodeWithdrawn, "the agent withdrew this question")
	case q.Answer == nil:
		return failOn(q, CodeNotAnswered, "this question has no recorded answer yet; wait for it")
	case q.Answer.RequestID != requestID || q.Answer.RequestSHA256 != sha:
		return failOn(q, CodeAnswerMismatch, "the request ID and digest must be those of the recorded answer")
	}
	return nil
}

func checkReportArgs(requestID, sha string) error {
	if !requestIDPattern.MatchString(requestID) || !digestPattern.MatchString(sha) {
		return fail(CodeInvalidArguments, "expected the recorded answer's --request-id and --request-sha256")
	}
	return nil
}

// Acknowledge records the agent's own report that it received exactly the
// recorded answer. It attests no new model turn and no delivery.
func Acknowledge(d Dependencies, root, id, requestID, sha string) (Question, error) {
	if e := checkReportArgs(requestID, sha); e != nil {
		return Question{}, e
	}
	return transition(d, root, id, func(q Question) (any, string, error) {
		if e := matchesWinner(q, requestID, sha); e != nil || q.Receipt != nil {
			return nil, "", e
		}
		return Receipt{requestID, sha, d.now(), ProvenanceAgent}, receiptPath(root, id), nil
	})
}

// Continue records the agent's report of how it used the answer after doing
// the work. Provenance stays reported_by_agent.
func Continue(d Dependencies, root, id, requestID, sha, summary string) (Question, error) {
	if e := checkReportArgs(requestID, sha); e != nil {
		return Question{}, e
	}
	if !ValidReport(summary) {
		return Question{}, fail(CodeInvalidArguments, "--summary must be non-blank text of at most 1000 codepoints without control characters other than tab and newline")
	}
	return transition(d, root, id, func(q Question) (any, string, error) {
		if e := matchesWinner(q, requestID, sha); e != nil {
			return nil, "", e
		}
		if q.Continuation != nil {
			if q.Continuation.Summary != summary {
				return nil, "", failOn(q, CodeReportConflict, "a different continuation summary is already recorded")
			}
			return nil, "", nil
		}
		return Continuation{requestID, sha, summary, d.now(), ProvenanceAgent}, continuationPath(root, id), nil
	})
}

// Withdraw ends an unanswered question without deleting it.
func Withdraw(d Dependencies, root, id, reason string) (Question, error) {
	if !ValidReport(reason) {
		return Question{}, fail(CodeInvalidArguments, "--reason must be non-blank text of at most 1000 codepoints without control characters other than tab and newline")
	}
	return transition(d, root, id, func(q Question) (any, string, error) {
		if q.Answer != nil {
			return nil, "", failOn(q, CodeAlreadyAnswered, "an answered question cannot be withdrawn")
		}
		if q.Withdrawal != nil {
			if q.Withdrawal.Reason != reason {
				return nil, "", failOn(q, CodeReportConflict, "the question was withdrawn with a different reason")
			}
			return nil, "", nil
		}
		return Withdrawal{reason, d.now(), ProvenanceAgent}, withdrawalPath(root, id), nil
	})
}

// Encode returns the canonical public JSON within the consumer output bound.
func Encode(v any) ([]byte, error) {
	b, e := taskrun.Canonical(v)
	if e != nil {
		return nil, e
	}
	if len(b)+1 > MaxOutputBytes {
		return nil, fail(CodeOutcomeUnknown, "the public result exceeds 4 MiB")
	}
	return b, nil
}
