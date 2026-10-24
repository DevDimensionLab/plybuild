//go:build !windows

package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type productAnswerBasis struct {
	Public               ProductQuestionReadback `json:"public"`
	QuestionRecordSHA256 string                  `json:"question_record_sha256"`
	Parent               ParentControlGate       `json:"parent"`
	QueueSHA256          string                  `json:"queue_sha256"`
	TaskFilesSHA256      string                  `json:"task_files_sha256"`
	WorkerRevision       uint64                  `json:"worker_revision"`
	WorkerPhase          string                  `json:"worker_phase"`
	RetryNotBefore       string                  `json:"retry_not_before"`
	ProviderRetries      int                     `json:"provider_retries"`
	MaxAttempts          int                     `json:"max_attempts"`
	PreparationSHA256    string                  `json:"preparation_sha256"`
}

func productAnswerBinding(b productAnswerBasis, request BoundAnswerRequest) string {
	b.Public.ObservedAtUTC = ""
	b.Public.Evidence = []BoundAnswerEvidence{}
	raw, _ := Canonical(request)
	var without map[string]json.RawMessage
	_ = json.Unmarshal(raw, &without)
	delete(without, "binding")
	return digest(struct {
		Domain  string                     `json:"domain"`
		Basis   productAnswerBasis         `json:"basis"`
		Request map[string]json.RawMessage `json:"request"`
	}{"ply.workflow.bound-answer/v2", b, without})
}

func productReadBase(d Dependencies, root, id string) ProductQuestionReadback {
	return ProductQuestionReadback{Kind: "ply.workflow.product-question", SchemaVersion: 2, Mode: BoundAnswerMode, Workspace: root, RunID: id, QuestionState: "unknown", Eligibility: "unknown", Reasons: []BoundAnswerReason{}, ObservedAtUTC: ControlTimestamp(controlNow(d)), Freshness: "unknown", Evidence: []BoundAnswerEvidence{}}
}
func productUnavailable(out *ProductQuestionReadback, eligibility string, e error) {
	out.Eligibility = eligibility
	out.Reasons = append(out.Reasons, BoundAnswerFailureOf(e))
}

func productRead(d Dependencies, root, id string) (productAnswerBasis, error) {
	b := productAnswerBasis{Public: productReadBase(d, root, id)}
	out := &b.Public
	if !boundAnswerRunID.MatchString(id) || len(root) > 4096 || physical(root, false) != nil {
		return b, boundAnswerError("invalid_arguments", "Expected a physical workspace and exact workflow Run ID.")
	}
	s, e := workflowRead(root, id)
	if e != nil {
		productUnavailable(out, "unknown", boundAnswerError("not_found", "The native Run could not be read."))
		return b, nil
	}
	if !productQuestionsEnabled(s) {
		out.QuestionState = "none"
		out.Freshness = "current"
		productUnavailable(out, "unsupported", boundAnswerError("unsupported", "The Run did not explicitly opt into same-owner product answers."))
		return b, nil
	}
	if s.ProductQuestion == nil {
		out.QuestionState = "none"
		out.Freshness = "current"
		productUnavailable(out, "blocked", boundAnswerError("no_question", "The Run has no published product question."))
		return b, nil
	}
	q, e := productQuestionRecordRead(s)
	if e != nil {
		productUnavailable(out, "unknown", e)
		return b, nil
	}
	out.Question = &q.Question
	out.Scope = &q.Scope
	out.OriginalDeadlineUTC = &q.OriginalDeadlineUTC
	out.QuestionState = "active"
	out.Freshness = "retained"
	out.Evidence = []BoundAnswerEvidence{{q.Question.QuestionID, "question", q.Question.QuestionSHA256, q.PublishedAtUTC}, {q.SourceJob.ID, "source_result", q.SourceResult.SHA256, q.PublishedAtUTC}, {"same-session-adapter", "qualification", q.QualificationSHA256, q.PublishedAtUTC}}
	w := s.Result.Worker
	out.Worker = &BoundAnswerWorker{OwnerID: "worker-" + id, OwnerGeneration: w.Generation, OwnerState: "unknown", Iteration: w.Iteration, SourceJobID: q.SourceJob.ID, SourceIntentSHA256: q.SourceJob.IntentHash, Provider: "codex", ProviderSessionID: q.Session.SessionID, ProviderBinarySHA256: s.Request.Runtime.Executable.SHA256, Model: q.Session.Model, Effort: q.Session.Effort, PolicySHA256: digest(q.Session.Policy), SessionBindingSHA256: digest(q.Session), QualificationSHA256: q.QualificationSHA256, ContinuationLimit: 1, ContinuationsDispatched: 0}
	owner := controlLockObservation(filepath.Join(workerDirectory(s), "owner.lock"))
	if owner == "lock_held" {
		out.Worker.OwnerState = "confirmed"
	} else if owner == "unowned" {
		out.Worker.OwnerState = "unowned"
	}
	gate, e := ReadParentControl(root, s.Request.ControlContract)
	if e != nil {
		productUnavailable(out, "unknown", e)
		return b, nil
	}
	out.ParentWeek = &BoundAnswerParentWeek{WeekID: *s.Request.ControlContract.WeekID, ContractSHA256: s.Request.ControlContract.Digest, WorksetSHA256: *s.Request.ControlContract.WeekBinding, ControlRevision: int(gate.Revision), DesiredAction: gate.DesiredAction, OriginalDeadlineUTC: gate.Deadline}
	b.Parent = gate
	b.QuestionRecordSHA256 = s.ProductQuestion.SHA256
	b.WorkerRevision = s.ControlRevision
	b.WorkerPhase = w.Phase
	b.RetryNotBefore = w.RetryNotBefore
	b.ProviderRetries = w.ProviderRetries
	b.MaxAttempts = s.Request.Worker.MaxAttempts
	b.PreparationSHA256 = s.Request.PreparationSHA256
	claim, claimErr := productClaimForQuestion(root, id, q.Question.QuestionID)
	if claimErr != nil && !os.IsNotExist(claimErr) {
		productUnavailable(out, "unknown", claimErr)
		return b, nil
	}
	if claimErr == nil {
		out.QuestionState = "answered"
		if _, e := productAnswerIndexed(root, claim); e != nil {
			out.QuestionState = "unknown"
			productUnavailable(out, "unknown", e)
			return b, nil
		}
		status := productClaimStatus(d, root, claim)
		if status.State == "continuation_observed" || status.State == "continuation_completed" {
			out.Worker.ContinuationsDispatched = 1
		}
		if status.State == "outcome_unknown" {
			productUnavailable(out, "unknown", boundAnswerError("outcome_unknown", "The answer is retained but its continuation outcome is not positively observed."))
			return b, nil
		}
		productUnavailable(out, "blocked", boundAnswerError("already_answered", "This question already has an immutable winning request."))
		return b, nil
	}
	if gate.Deadline != q.OriginalDeadlineUTC || gate.ActiveRun != id || gate.Revision > 9007199254740991 {
		productUnavailable(out, "unknown", boundAnswerError("parent_changed", "Parent deadline, active Run or dispatch ownership differs."))
		return b, nil
	}
	deadline, e := boundAnswerTime(gate.Deadline)
	if e != nil {
		productUnavailable(out, "unknown", e)
		return b, nil
	}
	if gate.DesiredAction == "stop" || !controlNow(d).Before(deadline) {
		productUnavailable(out, "blocked", boundAnswerError("deadline_or_stop", "The parent is stopped or its original deadline expired."))
		return b, nil
	}
	if gate.DesiredAction != "continue" && gate.DesiredAction != "pause" {
		productUnavailable(out, "unknown", boundAnswerError("parent_unknown", "Parent control action is unknown."))
		return b, nil
	}
	_, pending, e := ControlHistory(root, ControlTarget{"week", *s.Request.ControlContract.WeekID})
	if e != nil || pending != nil {
		productUnavailable(out, "blocked", boundAnswerError("control_pending", "A parent control operation is unresolved."))
		return b, nil
	}
	if out.Worker.OwnerState != "confirmed" || w.QueueReleased || w.Iteration != q.Question.Continuation.SourceIteration || w.Job == nil || *w.Job != q.SourceJob || (w.Phase != "awaiting_product_answer" && w.Phase != "paused") {
		productUnavailable(out, "blocked", boundAnswerError("owner_unknown", "The same live native owner, source attempt and wait reservation are not confirmed."))
		return b, nil
	}
	if w.RetryNotBefore != "" || w.Iteration > s.Request.Worker.MaxAttempts {
		productUnavailable(out, "blocked", boundAnswerError("budget_or_backoff", "An existing worker budget or provider backoff prevents continuation."))
		return b, nil
	}
	if e = productSourceCurrent(d, s, q); e != nil {
		productUnavailable(out, "blocked", e)
		return b, nil
	}
	b.QueueSHA256 = q.QueueSHA256
	b.TaskFilesSHA256 = q.TaskFilesSHA256
	out.Eligibility = "ready"
	out.Freshness = "current"
	return b, nil
}

func productSourceCurrent(d Dependencies, s workflowState, q productQuestionRecord) error {
	if verifyExecutable(s.Request.Runtime.Executable) != nil {
		return boundAnswerError("runtime_changed", "A frozen provider or controller executable changed.")
	}
	qual, e := productRuntimeQualification(d, s)
	if e != nil {
		return e
	}
	if qual != q.QualificationSHA256 {
		return boundAnswerError("runtime_changed", "Adapter qualification changed.")
	}
	scope, e := productScope(s)
	if e != nil {
		return e
	}
	if !equal(scope, q.Scope) {
		return boundAnswerError("binding_changed", "Native scope or delivery contract changed.")
	}
	queue, e := productNativeQueue(d, s)
	if e != nil {
		return e
	}
	if queue != q.QueueSHA256 {
		return boundAnswerError("owner_changed", "The native Task reservation changed.")
	}
	if _, e = workflowBound(q.SourceResult, 1<<20); e != nil {
		return e
	}
	sourceState := s
	copy := *s.Result.Worker
	copy.Job = &q.SourceJob
	sourceState.Result.Worker = &copy
	r, e := workerRunReserved(d, sourceState)
	if e != nil {
		return e
	}
	raw, e := workflowBound(r.Stdout, 16<<20)
	if e != nil {
		return e
	}
	thread, e := workerThreadID(raw)
	if e != nil || thread != q.Session.SessionID || r.ThreadID != thread || r.Exit == nil || *r.Exit != 0 || r.Error != "" || r.ProviderFailure != nil {
		return boundAnswerError("source_unknown", "Source model-turn completion is unavailable.")
	}
	if e = productVerifySession(q, false); e != nil {
		return e
	}
	x, e := productObserveTarget(d, s)
	if e != nil {
		return e
	}
	if digest(x) != q.TargetSHA256 {
		return boundAnswerError("file_drift", "The Task Git observation changed while waiting.")
	}
	files, e := productTaskFiles(q.Scope.TaskWorktree)
	if e != nil {
		return e
	}
	if files != q.TaskFilesSHA256 {
		return boundAnswerError("file_drift", "Task files changed after the question was published.")
	}
	return nil
}

func ReadProductQuestion(d Dependencies, root, runID string) (ProductQuestionReadback, error) {
	b, e := productRead(d, root, runID)
	return b.Public, e
}
func PreviewBoundAnswer(d Dependencies, root, runID, questionID string, revision int, optionID, requestID string) (out BoundAnswerPreview, err error) {
	out = BoundAnswerPreview{Kind: "ply.workflow.bound-answer-preview", SchemaVersion: 2, Mode: BoundAnswerMode, Workspace: root, RunID: runID, Eligibility: "unknown", Reasons: []BoundAnswerReason{}, Freshness: "unknown", Evidence: []BoundAnswerEvidence{}}
	// The enclosing observation includes the basis read and any issued request.
	defer func() { out.ObservedAtUTC = ControlTimestamp(controlNow(d)) }()
	if !boundAnswerQuestionID.MatchString(questionID) || revision != 1 || !boundAnswerOptionID.MatchString(optionID) || !boundAnswerRequestID.MatchString(requestID) {
		return out, boundAnswerError("invalid_arguments", "Expected exact question revision, option and request identities.")
	}
	b, e := productRead(d, root, runID)
	if e != nil {
		return out, e
	}
	out.Basis = &b.Public
	out.Eligibility = b.Public.Eligibility
	out.Freshness = b.Public.Freshness
	out.Reasons = b.Public.Reasons
	out.Evidence = b.Public.Evidence
	if b.Public.Eligibility != "ready" {
		return out, nil
	}
	if b.Public.Question.QuestionID != questionID || b.Public.Question.QuestionRevision != revision {
		out.Eligibility = "blocked"
		out.Reasons = []BoundAnswerReason{BoundAnswerFailureOf(boundAnswerError("stale_question", "The requested question revision is no longer active."))}
		return out, nil
	}
	selected := false
	for _, o := range b.Public.Question.Options {
		if o.OptionID == optionID {
			selected = true
		}
	}
	if !selected {
		out.Eligibility = "blocked"
		out.Reasons = []BoundAnswerReason{BoundAnswerFailureOf(boundAnswerError("invalid_option", "The option does not belong to the active question."))}
		return out, nil
	}
	now := controlNow(d)
	deadline, _ := boundAnswerTime(*b.Public.OriginalDeadlineUTC)
	expires := now.Add(300 * time.Second)
	if deadline.Before(expires) {
		expires = deadline
	}
	if !now.Before(expires) {
		return out, boundAnswerError("expired", "No positive preview validity remains.")
	}
	policy := "when_parent_allows"
	if b.Parent.DesiredAction == "pause" {
		policy = "deferred_until_explicit_parent_resume"
	}
	r := BoundAnswerRequest{Kind: "ply.workflow.bound-answer-request", SchemaVersion: 2, Mode: BoundAnswerMode, Action: "answer_product_question", RequestID: requestID, Workspace: root, RunID: runID, QuestionID: questionID, QuestionRevision: revision, OptionID: optionID, OriginalDeadlineUTC: *b.Public.OriginalDeadlineUTC, Validity: BoundAnswerValidity{ControlTimestamp(now), ControlTimestamp(expires)}, Binding: "", DispatchPolicy: policy}
	r.Binding = productAnswerBinding(b, r)
	out.Request = &r
	out.RequestSHA256 = ptr(digest(r))
	return out, nil
}

func productRequestMatches(r BoundAnswerRequest, b productAnswerBasis) error {
	q := b.Public.Question
	if b.Public.Eligibility != "ready" {
		if len(b.Public.Reasons) > 0 {
			return &b.Public.Reasons[0]
		}
		return boundAnswerError("unavailable", "Question is not ready.")
	}
	if q == nil || r.QuestionID != q.QuestionID || r.QuestionRevision != q.QuestionRevision || r.OriginalDeadlineUTC != b.Parent.Deadline || r.Binding != productAnswerBinding(b, r) {
		return boundAnswerError("stale", "The approved question, native authority or worker generation changed.")
	}
	policy := "when_parent_allows"
	if b.Parent.DesiredAction == "pause" {
		policy = "deferred_until_explicit_parent_resume"
	}
	if r.DispatchPolicy != policy {
		return boundAnswerError("stale", "The approved parent dispatch policy changed.")
	}
	for _, o := range q.Options {
		if o.OptionID == r.OptionID {
			return nil
		}
	}
	return boundAnswerError("invalid_option", fmt.Sprintf("Option %q is not part of this question.", r.OptionID))
}
