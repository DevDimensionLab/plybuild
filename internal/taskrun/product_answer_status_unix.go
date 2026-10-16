//go:build !windows

package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
)

type productCompletedAnswer struct {
	Status                      string `json:"status"`
	Summary                     string `json:"summary"`
	BlockReason                 string `json:"block_reason"`
	ConsumedAnswerRequestID     string `json:"consumed_answer_request_id"`
	ConsumedAnswerRequestSHA256 string `json:"consumed_answer_request_sha256"`
}
type productLaunchRecord struct {
	RunID           string `json:"run_id"`
	JobID           string `json:"job_id"`
	IntentSHA256    string `json:"intent_sha256"`
	RequestSHA256   string `json:"request_sha256"`
	Generation      int    `json:"generation"`
	AuthorizedAtUTC string `json:"authorized_at_utc"`
}

func productClaimQuestion(c productAnswerClaim) (productQuestionRecord, error) {
	var q productQuestionRecord
	raw, e := workflowBound(c.QuestionRecord, 2<<20)
	if e != nil {
		return q, e
	}
	e = decode(raw, 2<<20, &q)
	if e == nil && (q.Question.QuestionSHA256 != productQuestionDigest(q.Question) || q.Question.QuestionID != c.Request.QuestionID || q.Scope.RunRequestSHA256 != c.Basis.Public.Scope.RunRequestSHA256) {
		e = boundAnswerError("binding_unknown", "The claimed source question changed.")
	}
	return q, e
}
func productContinuationStatus(d Dependencies, c productAnswerClaim, o BoundAnswerStatus) BoundAnswerStatus {
	s, e := workflowRead(c.Request.Workspace, c.Request.RunID)
	if e != nil {
		return productUnknownStatus(o, e)
	}
	q, e := productClaimQuestion(c)
	if e != nil {
		return productUnknownStatus(o, e)
	}
	gate, e := ReadParentControl(c.Request.Workspace, s.Request.ControlContract)
	if e != nil {
		o.Reasons = append(o.Reasons, BoundAnswerFailureOf(e))
	} else {
		deadline, er := boundAnswerTime(gate.Deadline)
		if gate.ActiveRun != c.Request.RunID || gate.Deadline != q.OriginalDeadlineUTC || er != nil || !controlNow(d).Before(deadline) || gate.DesiredAction == "stop" {
			o.Reasons = append(o.Reasons, BoundAnswerFailureOf(boundAnswerError("parent_blocked", "Current parent ownership, stop or original deadline prevents another dispatch.")))
		} else if gate.DesiredAction == "pause" {
			o.State = "continuation_deferred"
			o.Reasons = append(o.Reasons, BoundAnswerFailureOf(boundAnswerError("parent_paused", "The recorded answer waits for explicit parent resume.")))
		}
	}
	dir := productContinuationDirectory(c)
	path := filepath.Join(dir, "intent.json")
	raw, e := readFile(path, 1<<20, true)
	if os.IsNotExist(e) {
		if s.ProductContinuation != nil {
			return productUnknownStatus(o, boundAnswerError("outcome_unknown", "The retained continuation intent is missing."))
		}
		return o
	}
	if e != nil {
		return productUnknownStatus(o, e)
	}
	var in workerAttemptIntent
	if e = decode(raw, 1<<20, &in); e != nil {
		return productUnknownStatus(o, e)
	}
	job := WorkerJob{ID: productContinuationJobID(c), Kind: "developer_answer", Directory: dir, IntentHash: hash(raw)}
	if in.JobID != job.ID || in.RunID != c.Request.RunID || in.Kind != job.Kind || in.BoundAnswer == nil || in.BoundAnswer.SHA256 != digest(c) || in.BoundAnswer.Locator != productClaimPath(c.Request.Workspace, c.Request.RunID, c.Request.QuestionID) {
		return productUnknownStatus(o, boundAnswerError("outcome_unknown", "Continuation intent does not bind the winning request."))
	}
	o.ContinuationRequested = ptr(true)
	o.Receipt.ContinuationJobID = &job.ID
	o.Receipt.ContinuationIntentSHA256 = &job.IntentHash
	o.Evidence = append(o.Evidence, BoundAnswerEvidence{job.ID, "continuation_intent", job.IntentHash, c.RecordedAtUTC})
	if s.ProductContinuation == nil || *s.ProductContinuation != job {
		return productUnknownStatus(o, boundAnswerError("outcome_unknown", "Continuation intent exists without its complete native reservation."))
	}
	o.State = "continuation_reserved" // A later pause never erases a known intent.
	var launch productLaunchRecord
	launchBinding, launchRaw, e := deliveryReadBinding(filepath.Join(dir, "launch.json"), 64<<10)
	if os.IsNotExist(e) {
		if s.ProductLaunch != nil {
			return productUnknownStatus(o, boundAnswerError("outcome_unknown", "The retained launch authorization is missing."))
		}
		return o
	}
	if e != nil {
		return productUnknownStatus(o, e)
	}
	if s.ProductLaunch == nil || *s.ProductLaunch != launchBinding {
		return productUnknownStatus(o, boundAnswerError("outcome_unknown", "Launch authorization has no matching native state binding."))
	}
	if e = decode(launchRaw, 64<<10, &launch); e != nil {
		return productUnknownStatus(o, e)
	}
	if launch.RunID != in.RunID || launch.JobID != job.ID || launch.IntentSHA256 != job.IntentHash || launch.RequestSHA256 != c.RequestSHA256 {
		return productUnknownStatus(o, boundAnswerError("outcome_unknown", "Launch identity is inconsistent."))
	}
	started, startErr := boundAnswerTime(launch.AuthorizedAtUTC)
	recorded, recordErr := boundAnswerTime(c.RecordedAtUTC)
	if startErr != nil || recordErr != nil || started.Before(recorded) {
		return productUnknownStatus(o, boundAnswerError("continuation_unknown", "The continuation launch boundary predates its recorded answer or has an invalid time."))
	}
	parentReasons := append([]BoundAnswerReason{}, o.Reasons...)
	o = productUnknownStatus(o, boundAnswerError("start_outcome_unknown", "The single launch authorization exists; do not resend."))
	// A provider process alone does not establish which session it resumed.
	// Positive thread.started + turn.started in the actual stream does.
	stream, e := readFile(filepath.Join(dir, "streams", "stdout.jsonl"), 16<<20, true)
	if e == nil && productObservedThread(stream, q.Session.SessionID) {
		var process Process
		if readValue(filepath.Join(dir, "process.json"), 64<<10, &process) == nil && process.PID != nil && process.StartIdentity != nil {
			if controlNow(d).Before(started) {
				return productUnknownStatus(o, boundAnswerError("continuation_unknown", "The continuation launch boundary is later than its current observation."))
			}
			o.State = "continuation_observed"
			o.ExecutionEffect = "continuation_observed"
			o.ObservedAtUTC = ControlTimestamp(controlNow(d))
			o.Receipt.ContinuationSessionID = &q.Session.SessionID
			o.Receipt.ContinuationStartedAtUTC = &launch.AuthorizedAtUTC
			o.Reasons = parentReasons
			o.Evidence = append(o.Evidence, BoundAnswerEvidence{job.ID, "provider_start", digest(process), launch.AuthorizedAtUTC})
		}
	}
	if _, e = os.Lstat(filepath.Join(dir, "result.json")); os.IsNotExist(e) {
		return o
	} else if e != nil {
		return productUnknownStatus(o, e)
	}
	state := s
	w := *s.Result.Worker
	w.Job = &job
	state.Result.Worker = &w
	r, e := workerRunReserved(d, state)
	if e != nil {
		return productUnknownStatus(o, e)
	}
	completed, e := productValidateCompleted(d, state, c, q, r)
	if e != nil {
		return productUnknownStatus(o, e)
	}
	finished, finishErr := boundAnswerTime(r.FinishedAt)
	if finishErr != nil || finished.Before(started) {
		return productUnknownStatus(o, boundAnswerError("continuation_unknown", "The continuation result predates its single launch boundary."))
	}
	result, _, e := deliveryReadBinding(filepath.Join(dir, "result.json"), 1<<20)
	if e != nil {
		return productUnknownStatus(o, e)
	}
	o.State = "continuation_completed"
	o.ExecutionEffect = "continuation_completed"
	o.Reasons = parentReasons
	o.ObservedAtUTC = ControlTimestamp(controlNow(d))
	o.Receipt.ContinuationSessionID = &r.ThreadID
	// Keep the same immutable launch boundary used by observed. r.StartedAt
	// timestamps attempt preparation, before policy checks and physical Start.
	// The process/stream/result checks above establish execution separately.
	o.Receipt.ContinuationStartedAtUTC = &launch.AuthorizedAtUTC
	o.Receipt.ContinuationCompletedAtUTC = &r.FinishedAt
	o.Receipt.ResultSHA256 = &result.SHA256
	o.Receipt.ConsumedAnswerRequestID = &completed.ConsumedAnswerRequestID
	o.Receipt.ConsumedAnswerRequestSHA256 = &completed.ConsumedAnswerRequestSHA256
	o.Evidence = append(o.Evidence, BoundAnswerEvidence{job.ID, "continuation_result", result.SHA256, r.FinishedAt})
	return o
}
func productObservedThread(raw []byte, expected string) bool {
	thread, started := "", false
	for _, line := range bytes.Split(raw, []byte{'\n'}) {
		if len(line) == 0 {
			continue
		}
		var v struct {
			Type     string `json:"type"`
			ThreadID string `json:"thread_id"`
		}
		if json.Unmarshal(line, &v) != nil {
			break
		}
		if v.Type == "thread.started" {
			if thread != "" {
				return false
			}
			thread = v.ThreadID
		}
		if v.Type == "turn.started" {
			started = true
		}
	}
	return thread == expected && started
}
func productValidateCompleted(d Dependencies, s workflowState, c productAnswerClaim, q productQuestionRecord, r workerAttemptResult) (productCompletedAnswer, error) {
	var answer productCompletedAnswer
	started, startErr := boundAnswerTime(r.StartedAt)
	finished, finishErr := boundAnswerTime(r.FinishedAt)
	recorded, recordErr := boundAnswerTime(c.RecordedAtUTC)
	if startErr != nil || finishErr != nil || recordErr != nil || started.Before(recorded) || finished.Before(started) || controlNow(d).Before(finished) {
		return answer, boundAnswerError("continuation_unknown", "Continuation timestamps do not form an observed ordered result after this answer.")
	}
	if r.Error != "" || r.Exit == nil || *r.Exit != 0 || r.Answer == nil || r.ThreadID != q.Session.SessionID || r.ProviderFailure != nil {
		return answer, boundAnswerError("continuation_unknown", "Continuation has no successful result for the same provider session.")
	}
	raw, e := workflowBound(r.Stdout, 16<<20)
	if e != nil {
		return answer, e
	}
	thread, e := workerThreadID(raw)
	if e != nil || thread != q.Session.SessionID {
		return answer, boundAnswerError("continuation_unknown", "Continuation stream does not prove the exact source session.")
	}
	raw, e = workflowBound(*r.Answer, 1<<20)
	if e != nil {
		return answer, e
	}
	if e = decode(raw, 1<<20, &answer); e != nil {
		return answer, e
	}
	if answer.Status != "completed" || answer.ConsumedAnswerRequestID != c.Request.RequestID || answer.ConsumedAnswerRequestSHA256 != c.RequestSHA256 {
		return answer, boundAnswerError("consumed_mismatch", "The completed answer does not confirm consumption of this exact winning request.")
	}
	if e = productVerifySession(q, true); e != nil {
		return answer, e
	}
	if _, _, e = productReadSessionTurns(s, q.SourceJob, q.Session.SessionID, q.Session.Path, 2); e != nil {
		return answer, e
	}
	return answer, nil
}
