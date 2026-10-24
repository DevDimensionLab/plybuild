//go:build !windows

package taskrun

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func workerProductSchema(continuation bool) any {
	text := map[string]any{"type": "string"}
	properties := map[string]any{"status": map[string]any{"type": "string", "enum": []string{"completed", "blocked", "needs_product_answer"}}, "summary": text, "block_reason": text}
	required := []string{"status", "summary", "block_reason", "product_question"}
	if continuation {
		properties["status"] = map[string]any{"type": "string", "enum": []string{"completed", "blocked"}}
		properties["consumed_answer_request_id"] = text
		properties["consumed_answer_request_sha256"] = text
		required = []string{"status", "summary", "block_reason", "consumed_answer_request_id", "consumed_answer_request_sha256"}
	} else {
		bounded := func(n int) any { return map[string]any{"type": "string", "minLength": 1, "maxLength": n} }
		option := map[string]any{"type": "object", "properties": map[string]any{"option_id": map[string]any{"type": "string", "pattern": "^[a-z0-9][a-z0-9_-]{0,63}$"}, "label": bounded(120), "consequence": bounded(1024)}, "required": []string{"option_id", "label", "consequence"}, "additionalProperties": false}
		next := map[string]any{"type": "object", "properties": map[string]any{"type": map[string]any{"const": "resume_same_developer_thread"}, "source_iteration": map[string]any{"type": "integer", "minimum": 1}, "source_job_id": bounded(128)}, "required": []string{"type", "source_iteration", "source_job_id"}, "additionalProperties": false}
		properties["product_question"] = map[string]any{"type": []string{"object", "null"}, "properties": map[string]any{"purpose": map[string]any{"const": "product_choice"}, "text": bounded(2048), "reason": bounded(1024), "options": map[string]any{"type": "array", "minItems": 2, "maxItems": 5, "items": option}, "continuation": next}, "required": []string{"purpose", "text", "reason", "options", "continuation"}, "additionalProperties": false}
	}
	return map[string]any{"type": "object", "properties": properties, "required": required, "additionalProperties": false}
}
func workerProductPrompt(s workflowState) string {
	id := fmt.Sprintf("developer-%04d", s.Result.Worker.Iteration)
	if s.Result.Worker.ProviderRetries > 0 {
		id += fmt.Sprintf("-retry-%04d", s.Result.Worker.ProviderRetries)
	}
	return fmt.Sprintf("This newly opted-in Goal permits one bounded product choice. If the Goal needs a human product decision, finish this model turn with status needs_product_answer and product_question purpose product_choice, a concrete question and reason, and 2–5 distinct options with option_id, label and consequence. Set continuation.type to resume_same_developer_thread, source_job_id to %s and source_iteration to %d. Return product_question null for completed or blocked. Authentication, quota, QA, release and changed authority are not product choices. A question preserves this same session and Task; do not start a background writer or claim global writer quiescence. The controller waits for a bound answer before it can resume your session. No candidate executable is required for needs_product_answer. Only one question and one answer turn are authorized for this Run.\n", id, s.Result.Worker.Iteration)
}

// A later independent review (and any ordinary correction) must assess the
// selected solution, even when the Goal permits multiple valid alternatives.
// The immutable prompt becomes part of the actual native attempt intent.
func workerProductChoiceContext(d Dependencies, s workflowState) (string, error) {
	if s.ProductQuestion == nil {
		return "", nil
	}
	q, e := productQuestionRecordRead(s)
	if e != nil {
		return "", e
	}
	c, e := productClaimForQuestion(s.Request.WorkspaceRoot, s.Result.RunID, q.Question.QuestionID)
	if e != nil {
		return "", e
	}
	status := productClaimStatus(d, s.Request.WorkspaceRoot, c)
	if status.State != "continuation_completed" || status.Receipt == nil || status.Receipt.ResultSHA256 == nil {
		return "", boundAnswerError("continuation_unknown", "Review or correction requires the validated consumed result for the selected answer.")
	}
	var option ProductQuestionOption
	for _, v := range q.Question.Options {
		if v.OptionID == c.Request.OptionID {
			option = v
		}
	}
	if option.OptionID == "" {
		return "", boundAnswerError("binding_unknown", "The winning option is absent from the preserved question.")
	}
	return fmt.Sprintf("The Goal has a completed bound product choice. Independently verify the selected product choice in the actual candidate, not merely either alternative allowed by the Goal. Question %s (digest %s): %s. Selected option ID %q, label %q, consequence %q. Winning request ID %s, canonical digest %s. Immutable claim %s (digest %s). Validated same-session consumed continuation result %s (digest %s), source session %s; consumed request ID %s, digest %s. Candidate commit %s, tree %s. Check the implementation, tests and review findings against this chosen solution. Preserve this choice during any correction. No additional product question or answer continuation is authorized; a correction returns product_question null and keeps the existing ordinary attempt budget.\n", q.Question.QuestionID, q.Question.QuestionSHA256, q.Question.Text, option.OptionID, option.Label, option.Consequence, c.Request.RequestID, c.RequestSHA256, productClaimPath(s.Request.WorkspaceRoot, s.Result.RunID, q.Question.QuestionID), digest(c), filepath.Join(s.ProductContinuation.Directory, "result.json"), *status.Receipt.ResultSHA256, q.Session.SessionID, *status.Receipt.ConsumedAnswerRequestID, *status.Receipt.ConsumedAnswerRequestSHA256, s.Result.Worker.CandidateOID, s.Result.Worker.CandidateTree), nil
}

func productCheckParent(d Dependencies, s workflowState, q productQuestionRecord) error {
	gate, e := ReadParentControl(s.Request.WorkspaceRoot, s.Request.ControlContract)
	if e != nil {
		return e
	}
	if gate.ActiveRun != s.Result.RunID || gate.Deadline != q.OriginalDeadlineUTC {
		return boundAnswerError("parent_changed", "The parent Run or original deadline changed.")
	}
	if e = CheckParentDispatch(s.Request.WorkspaceRoot, s.Request.ControlContract, controlNow(d)); e != nil {
		return e
	}
	if s.Result.Worker.RetryNotBefore != "" || s.Result.Worker.Iteration > s.Request.Worker.MaxAttempts || s.Result.Worker.QueueReleased {
		return boundAnswerError("budget_or_backoff", "The existing budget, backoff or native reservation prevents dispatch.")
	}
	return nil
}

func productBuildContinuation(d Dependencies, s workflowState, c productAnswerClaim, q productQuestionRecord) (workerAttemptIntent, error) {
	source, e := workerIntent(s, q.SourceJob)
	if e != nil {
		return workerAttemptIntent{}, e
	}
	split := -1
	for i, arg := range source.Argv {
		if arg == "exec" {
			split = i
			break
		}
	}
	if split < 1 || split+1 < len(source.Argv) && source.Argv[split+1] == "resume" {
		return workerAttemptIntent{}, boundAnswerError("intent_unknown", "Source intent is not a first developer exec turn.")
	}
	dir := productContinuationDirectory(c)
	option := ProductQuestionOption{}
	for _, v := range q.Question.Options {
		if v.OptionID == c.Request.OptionID {
			option = v
		}
	}
	prompt := fmt.Sprintf("Continue this exact original developer session for the same registered Goal. The native owner retains the same Task, Git, control, review and delivery authority. The product answer is option %q (%s): %s. Request ID: %s. Request SHA-256: %s. Complete the authorized implementation and acceptance work with this chosen option. Do not ask another question, delegate, run in the background, alter permissions or run Ply mutations. Build the candidate executable at the original path %s. Return completed only when it exists and your work is ready for the separate acceptance/reviewer. Return consumed_answer_request_id=%s and consumed_answer_request_sha256=%s exactly; these fields attest use of this input, not integration or global writer quiescence.\n", option.OptionID, option.Label, option.Consequence, c.Request.RequestID, c.RequestSHA256, filepath.Join(q.SourceJob.Directory, "output", "candidate"), c.Request.RequestID, c.RequestSHA256)
	pp := filepath.Join(dir, "prompt.txt")
	if e = d.writeOnce(pp, []byte(prompt)); e != nil {
		return workerAttemptIntent{}, e
	}
	schema, e := workflowKeep(d, filepath.Join(dir, "schema.json"), workerProductSchema(true))
	if e != nil {
		return workerAttemptIntent{}, e
	}
	// Reuse the source prefix and environment exactly. Output/tmp remain the
	// original permitted roots; the new attempt's evidence is controller-only.
	argv := append([]string{}, source.Argv[:split]...)
	argv = append(argv, "exec", "resume", "--ignore-user-config", "--ignore-rules", "--model", s.Request.Runtime.Model, "--json", "--output-schema", schema.Locator, "-o", filepath.Join(dir, "answer.json"), q.Session.SessionID, "-")
	seconds := source.TimeoutSeconds
	deadline, _ := boundAnswerTime(q.OriginalDeadlineUTC)
	remaining := int(deadline.Sub(controlNow(d)) / time.Second)
	if remaining < seconds {
		seconds = remaining
	}
	if seconds < 1 {
		return workerAttemptIntent{}, boundAnswerError("expired", "No provider turn fits before the original deadline.")
	}
	return workerAttemptIntent{BoundAnswer: &FileBinding{Locator: productClaimPath(c.Request.Workspace, c.Request.RunID, c.Request.QuestionID), SHA256: digest(c)}, RunID: s.Result.RunID, JobID: productContinuationJobID(c), Kind: "developer_answer", Generation: s.Result.Worker.Generation, CandidateOID: source.CandidateOID, CandidateTree: source.CandidateTree, Argv: argv, Environment: append([]string{}, source.Environment...), Prompt: &FileBinding{Locator: pp, SHA256: hash([]byte(prompt))}, Schema: &schema, TimeoutSeconds: seconds}, nil
}

func workerReserveProductContinuation(d Dependencies, s workflowState, c productAnswerClaim) error {
	q, e := productClaimQuestion(c)
	if e != nil {
		return e
	}
	return WithWorkerParentGuard(s.Request.WorkspaceRoot, s.Request.ControlContract, func() error {
		return workflowUpdate(d, s.Request.WorkspaceRoot, s.Result.RunID, func(current *workflowState) error {
			if e := workerLiveOwner(d, *current); e != nil {
				return e
			}
			if current.ProductContinuation != nil {
				return nil
			}
			if _, e := os.Lstat(filepath.Join(productContinuationDirectory(c), "intent.json")); !os.IsNotExist(e) {
				return boundAnswerError("outcome_unknown", "A continuation intent exists without its native reservation; preserve it for inspection.")
			}
			if _, e := productAnswerIndexed(s.Request.WorkspaceRoot, c); e != nil {
				return e
			}
			if current.ProductQuestion == nil || *current.ProductQuestion != c.QuestionRecord || current.Result.Worker.Job == nil || *current.Result.Worker.Job != q.SourceJob {
				return boundAnswerError("stale", "Source wait reservation changed.")
			}
			if e := productCheckParent(d, *current, q); e != nil {
				return e
			}
			if e := productSourceCurrent(d, *current, q); e != nil {
				return e
			}
			in, e := productBuildContinuation(d, *current, c, q)
			if e != nil {
				return e
			}
			b, e := workflowKeep(d, filepath.Join(productContinuationDirectory(c), "intent.json"), in)
			if e != nil {
				return e
			}
			job := &WorkerJob{ID: in.JobID, Kind: in.Kind, Directory: productContinuationDirectory(c), IntentHash: b.SHA256}
			current.ProductContinuation = job
			current.ProductContinuationGeneration = current.Result.Worker.Generation
			current.Result.Worker.Job = job
			return d.fault("bound_answer_after_continuation_intent")
		})
	})
}

func workerProductFinalStart(d Dependencies, s workflowState, in workerAttemptIntent, start func() error) error {
	if in.BoundAnswer == nil {
		return boundAnswerError("intent_unknown", "Continuation answer binding is absent.")
	}
	raw, e := workflowBound(*in.BoundAnswer, 256<<10)
	if e != nil {
		return e
	}
	var c productAnswerClaim
	if e = decode(raw, 256<<10, &c); e != nil {
		return e
	}
	q, e := productClaimQuestion(c)
	if e != nil {
		return e
	}
	if e = d.fault("bound_answer_before_final_start"); e != nil {
		return e
	}
	return WithWorkerParentGuard(s.Request.WorkspaceRoot, s.Request.ControlContract, func() error {
		return withStore(s.Request.WorkspaceRoot, func() error {
			current, e := workflowRead(s.Request.WorkspaceRoot, s.Result.RunID)
			if e != nil {
				return e
			}
			if e = workerLiveOwner(d, current); e != nil {
				return e
			}
			job := current.ProductContinuation
			if job == nil || current.Result.Worker.Job == nil || *job != *current.Result.Worker.Job || job.ID != in.JobID || job.IntentHash != digest(in) || current.ProductContinuationGeneration != current.Result.Worker.Generation {
				return boundAnswerError("owner_changed", "The continuation reservation or current native generation changed.")
			}
			if _, e = productAnswerIndexed(s.Request.WorkspaceRoot, c); e != nil {
				return e
			}
			if e = productCheckParent(d, current, q); e != nil {
				return e
			}
			if e = productSourceCurrent(d, current, q); e != nil {
				return e
			}
			launchPath := filepath.Join(job.Directory, "launch.json")
			if current.ProductLaunch != nil {
				return boundAnswerError("outcome_unknown", "The single provider launch was already authorized; no replay is permitted.")
			}
			if _, e = os.Lstat(launchPath); !os.IsNotExist(e) {
				return boundAnswerError("outcome_unknown", "A launch record already exists or is unreadable; no second dispatch is permitted.")
			}
			launch := productLaunchRecord{RunID: current.Result.RunID, JobID: job.ID, IntentSHA256: job.IntentHash, RequestSHA256: c.RequestSHA256, Generation: current.Result.Worker.Generation, AuthorizedAtUTC: ControlTimestamp(controlNow(d))}
			binding, e := workflowKeep(d, launchPath, launch)
			if e != nil {
				return e
			}
			current.ProductLaunch = &binding
			current.Result.Worker.Phase = "running"
			current.ControlRevision++
			if e = workflowSave(d, current); e != nil {
				return e
			}
			if e = d.fault("bound_answer_after_launch_authorization"); e != nil {
				return e
			}
			// Locks do not stop time. Inventory and durable writes above may have
			// crossed the original deadline; authorize the syscall only while it
			// remains live. The launch record stays immutable on late refusal.
			if e = productCheckParent(d, current, q); e != nil {
				return e
			}
			// The parent guard remains held through this immediate syscall only. No
			// provider wait, stream read, owner acquisition or completion happens here.
			return start()
		})
	})
}

func workerLaunchProductContinuation(d Dependencies, s workflowState) error {
	control, controlErr := currentPlyPath(d)
	if controlErr != nil {
		return controlErr
	}
	job := s.ProductContinuation
	if job == nil {
		return boundAnswerError("intent_unknown", "Continuation is not reserved.")
	}
	if s.ProductLaunch != nil {
		return boundAnswerError("outcome_unknown", "The continuation launch cannot be replayed.")
	}
	if e := d.fault("bound_answer_after_continuation_reservation"); e != nil {
		return e
	}
	cmd := exec.Command(control, "workflow", "worker", "_attempt", s.Result.RunID, "--workspace", s.Request.WorkspaceRoot, "--job", job.ID)
	cmd.Dir, cmd.Stdout, cmd.Stderr = s.Observed.Target.WorktreeLocator, os.Stdout, os.Stderr
	cmd.ExtraFiles = []*os.File{d.WorkerOwner.lock}
	e := cmd.Run()
	current, readErr := workflowRead(s.Request.WorkspaceRoot, s.Result.RunID)
	if readErr != nil {
		return readErr
	}
	if _, resultErr := os.Lstat(filepath.Join(job.Directory, "result.json")); resultErr == nil {
		return nil
	}
	if current.ProductLaunch != nil {
		return boundAnswerError("outcome_unknown", "The single continuation launch has no confirmed result; do not resend.")
	}
	// A rejected final gate is a preserved reservation. The next allowed turn
	// still uses this intent; it is not a new answer or a new provider retry.
	q, readErr := productQuestionRecordRead(current)
	if readErr != nil {
		return readErr
	}
	if check := productCheckParent(d, current, q); check != nil {
		return nil
	}
	if e != nil {
		return fmt.Errorf("continuation runner failed before launch: %w", e)
	}
	return boundAnswerError("outcome_unknown", "Continuation runner returned without a launch/result record.")
}

// Handles only the retained question/developer boundary. Completed answers
// leave this branch through workerStep and the existing commit/review pipeline.
func workerProductWaiting(d Dependencies, s workflowState) (bool, error) {
	if s.ProductQuestion == nil || s.Result.Worker.Stage != "developer" {
		return false, nil
	}
	q, e := productQuestionRecordRead(s)
	if e != nil {
		return true, e
	}
	if s.Result.Worker.Iteration != q.Question.Continuation.SourceIteration {
		// A later ordinary correction has its own attempt and native controls.
		// The old answer is immutable review context, not another wait boundary.
		return false, nil
	}
	gate, e := ReadParentControl(s.Request.WorkspaceRoot, s.Request.ControlContract)
	if e != nil {
		return true, e
	}
	if gate.ActiveRun != s.Result.RunID || gate.Deadline != q.OriginalDeadlineUTC {
		return true, boundAnswerError("parent_changed", "The parent no longer retains this exact Run and deadline.")
	}
	deadline, e := boundAnswerTime(gate.Deadline)
	if e != nil {
		return true, e
	}
	action := gate.DesiredAction
	if !controlNow(d).Before(deadline) {
		action = "stop"
	}
	if action == "pause" || action == "stop" {
		// The live attempt runner holds the owner lock until it has actually ended.
		// A reserved but unstarted continuation needs no provider completion.
		if s.ProductLaunch != nil {
			if _, e = workerRunReserved(d, s); e != nil {
				return true, e
			}
		}
		phase := map[string]string{"pause": "paused", "stop": "stopped"}[action]
		if s.Result.Worker.Phase != phase {
			e = workerStateUpdate(d, s.Request.WorkspaceRoot, s.Result.RunID, func(current *workflowState) error {
				current.Result.Worker.Phase = phase
				current.Result.NextAction = WorkflowAction{"parent_control", "The product answer and same Task are retained under parent " + action + "."}
				return nil
			})
		}
		return true, e
	}
	if action != "continue" {
		return true, boundAnswerError("parent_unknown", "Unknown parent action.")
	}
	if s.Result.Worker.Phase == "paused" || s.Result.Worker.Phase == "stopped" {
		if !gate.ResumeAllowed {
			return true, nil
		}
		e = workerStateUpdate(d, s.Request.WorkspaceRoot, s.Result.RunID, func(current *workflowState) error {
			current.Result.Worker.Phase = "awaiting_product_answer"
			if current.ProductLaunch != nil {
				current.Result.Worker.Phase = "running"
			}
			return d.replaceValue(filepath.Join(workerDirectory(*current), "control.json"), workerControl{current.Result.RunID, "continue", ControlTimestamp(controlNow(d))})
		})
		return true, e
	}
	if s.Result.Worker.Phase == "blocked" {
		return true, nil
	}
	c, e := productClaimForQuestion(s.Request.WorkspaceRoot, s.Result.RunID, q.Question.QuestionID)
	if os.IsNotExist(e) {
		return true, d.fault("worker_product_waiting")
	}
	if e != nil {
		return true, e
	}
	if _, e = productAnswerIndexed(s.Request.WorkspaceRoot, c); e != nil {
		return true, e
	}
	if s.ProductContinuation == nil {
		return true, workerReserveProductContinuation(d, s, c)
	}
	if s.ProductLaunch == nil {
		return true, workerLaunchProductContinuation(d, s)
	}
	return false, nil
}

// A restart recovers ownership only after the held native lock and preserved
// parent/queue/source evidence agree. The caller then advances the generation;
// pre-restart previews and callbacks cannot inherit this capability.
func productReadoption(d Dependencies, s workflowState) error {
	if s.ProductQuestion == nil || s.Result.Worker.Stage != "developer" {
		return nil
	}
	q, e := productQuestionRecordRead(s)
	if e != nil {
		return e
	}
	gate, e := ReadParentControl(s.Request.WorkspaceRoot, s.Request.ControlContract)
	if e != nil {
		return e
	}
	if gate.ActiveRun != s.Result.RunID || gate.Deadline != q.OriginalDeadlineUTC {
		return boundAnswerError("owner_unknown", "Restart cannot adopt another parent reservation.")
	}
	if s.Result.Worker.Iteration != q.Question.Continuation.SourceIteration {
		if s.Result.Worker.Iteration <= q.Question.Continuation.SourceIteration {
			return boundAnswerError("owner_unknown", "The original developer iteration moved backwards.")
		}
		if _, e = workerProductChoiceContext(d, s); e != nil {
			return e
		}
		_, e = productNativeQueue(d, s)
		return e
	}
	if s.ProductLaunch == nil {
		if s.ProductContinuation != nil {
			return productUnstartedReservation(d, s)
		}
		c, e := productClaimForQuestion(s.Request.WorkspaceRoot, s.Result.RunID, q.Question.QuestionID)
		if e != nil && !os.IsNotExist(e) {
			return e
		}
		if e == nil {
			if _, e = productAnswerIndexed(s.Request.WorkspaceRoot, c); e != nil {
				return e
			}
			if _, e = os.Lstat(filepath.Join(productContinuationDirectory(c), "intent.json")); !os.IsNotExist(e) {
				return boundAnswerError("outcome_unknown", "A partial continuation reservation cannot authorize a new owner.")
			}
		}
		return productSourceCurrent(d, s, q)
	}
	c, e := productClaimForQuestion(s.Request.WorkspaceRoot, s.Result.RunID, q.Question.QuestionID)
	if e != nil {
		return e
	}
	if _, e = productAnswerIndexed(s.Request.WorkspaceRoot, c); e != nil {
		return e
	}
	if s.ProductContinuation == nil || s.Result.Worker.Job == nil || *s.ProductContinuation != *s.Result.Worker.Job {
		return boundAnswerError("outcome_unknown", "Started continuation reservation changed.")
	}
	r, e := workerRunReserved(d, s)
	if e != nil {
		return e
	}
	if _, e = productValidateCompleted(d, s, c, q, r); e != nil {
		return e
	}
	queue, e := productNativeQueue(d, s)
	if e != nil {
		return e
	}
	if queue != q.QueueSHA256 {
		return boundAnswerError("owner_changed", "Native reservation changed after continuation.")
	}
	return nil
}

// A positively reserved continuation with no launch authorization is a safe
// parent-control boundary. This is a native protocol fact, not writer quiescence.
func productUnstartedReservation(d Dependencies, s workflowState) error {
	if s.ProductContinuation == nil || s.ProductLaunch != nil || s.Result.Worker.Job == nil || *s.ProductContinuation != *s.Result.Worker.Job {
		return boundAnswerError("outcome_unknown", "Continuation is not at its reserved, unstarted boundary.")
	}
	in, e := workerIntent(s, *s.ProductContinuation)
	if e != nil {
		return e
	}
	if in.BoundAnswer == nil {
		return boundAnswerError("intent_unknown", "Bound answer is missing.")
	}
	if _, e = os.Lstat(filepath.Join(s.ProductContinuation.Directory, "launch.json")); !os.IsNotExist(e) {
		return boundAnswerError("outcome_unknown", "Launch authorization exists or cannot be observed.")
	}
	c, e := productReadClaim(in.BoundAnswer.Locator)
	if e != nil {
		return e
	}
	if digest(c) != in.BoundAnswer.SHA256 {
		return boundAnswerError("binding_changed", "Answer claim changed.")
	}
	if _, e = productAnswerIndexed(s.Request.WorkspaceRoot, c); e != nil {
		return e
	}
	q, e := productClaimQuestion(c)
	if e != nil {
		return e
	}
	return productSourceCurrent(d, s, q)
}
