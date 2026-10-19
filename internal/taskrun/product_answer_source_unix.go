//go:build !windows

package taskrun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const productAnswerQualification = "sha256:378849e8cf9825acce41a3170063b1095e267b613145dd1f6f6ae0d723c16337"

type productQuestionProposal struct {
	Purpose      string                  `json:"purpose"`
	Text         string                  `json:"text"`
	Reason       string                  `json:"reason"`
	Options      []ProductQuestionOption `json:"options"`
	Continuation struct {
		Type            string `json:"type"`
		SourceIteration int    `json:"source_iteration"`
		SourceJobID     string `json:"source_job_id"`
	} `json:"continuation"`
}
type productDeveloperAnswer struct {
	Status          string                   `json:"status"`
	Summary         string                   `json:"summary"`
	BlockReason     string                   `json:"block_reason"`
	ProductQuestion *productQuestionProposal `json:"product_question"`
}

type productSessionBinding struct {
	Path      string          `json:"path"`
	Prefix    FileBinding     `json:"prefix"`
	Length    int             `json:"length"`
	SessionID string          `json:"session_id"`
	CWD       string          `json:"cwd"`
	Model     string          `json:"model"`
	Effort    string          `json:"effort"`
	Policy    json.RawMessage `json:"policy"`
}

type productQuestionRecord struct {
	Question            ProductQuestion       `json:"question"`
	Scope               BoundAnswerScope      `json:"scope"`
	SourceJob           WorkerJob             `json:"source_job"`
	SourceResult        FileBinding           `json:"source_result"`
	Session             productSessionBinding `json:"session"`
	TaskFilesSHA256     string                `json:"task_files_sha256"`
	TargetSHA256        string                `json:"target_sha256"`
	QueueSHA256         string                `json:"queue_sha256"`
	OriginalDeadlineUTC string                `json:"original_deadline_utc"`
	QualificationSHA256 string                `json:"qualification_sha256"`
	PublishedAtUTC      string                `json:"published_at_utc"`
}

func productQuestionsEnabled(s workflowState) bool {
	return s.Request.Worker != nil && s.Request.Worker.SchemaVersion == 2 && s.Request.Worker.ProductQuestions != nil && ValidateProductQuestionContract(s.Request.Worker.ProductQuestions) == nil && s.Request.ControlContract != nil && s.Request.ControlContract.ParentKind == "week" && ValidateControlContract(s.Request.ControlContract) == nil
}
func productRuntimeQualification(d Dependencies, s workflowState) (string, error) {
	if d.ProductAnswerQualification != nil {
		return d.ProductAnswerQualification(s.Request.Runtime, s.Request.Delivery.ReasoningEffort)
	}
	if s.Request.Runtime.Model != "gpt-6-astra" || s.Request.Delivery.ReasoningEffort != "max" {
		return "", boundAnswerError("unsupported", "This provider model/effort has no preserved same-session adapter qualification.")
	}
	return productAnswerQualification, nil
}
func productNativeQueue(d Dependencies, s workflowState) (string, error) {
	if s.Result.Worker.QueueReleased {
		return "", boundAnswerError("owner_unknown", "The native queue reservation was released.")
	}
	if d.ControlQueue != nil {
		return d.ControlQueue(s.Request.WorkspaceRoot, s.Result.RunID)
	}
	p, e := workerPreparation(d, s)
	if e != nil {
		return "", e
	}
	if e = workerPreparationCurrent(d, s); e != nil {
		return "", e
	}
	return digest(p), nil
}
func productScope(s workflowState) (BoundAnswerScope, error) {
	var draft struct {
		Binding struct {
			ProjectID string `json:"project_id"`
			RepoID    string `json:"repo_id"`
		} `json:"binding_request"`
		Task workspace.TaskSpecBasis `json:"task_spec_binding"`
	}
	if e := json.Unmarshal(s.Request.HandoffDraft, &draft); e != nil {
		return BoundAnswerScope{}, e
	}
	c := s.Request.Delivery
	if c == nil || c.Agreement == nil {
		return BoundAnswerScope{}, boundAnswerError("binding_unknown", "The delivery agreement is unavailable.")
	}
	// The handoff binds the materialized execution Spec. Its revision is not
	// the selected base-free Goal revision. Join the immutable goal_origin to
	// the exact Goal manifest preserved by the native delivery request.
	goalRaw, e := workflowBound(c.Goal, 4<<20)
	if e != nil {
		return BoundAnswerScope{}, e
	}
	var goal struct {
		Kind          string           `json:"kind"`
		SchemaVersion int              `json:"schema_version"`
		ContractKind  string           `json:"contract_kind"`
		TaskID        workspace.TaskID `json:"task_id"`
		SpecID        string           `json:"spec_id"`
		Revision      int              `json:"revision"`
	}
	if json.Unmarshal(goalRaw, &goal) != nil || goal.Kind != "WorkspaceTaskSpecRevision@2" || goal.SchemaVersion != 2 || goal.ContractKind != "goal" || goal.TaskID != draft.Task.TaskID || goal.SpecID != draft.Task.SpecID || goal.Revision < 1 || !digestPattern.MatchString(draft.Task.Spec.ManifestSHA256) {
		return BoundAnswerScope{}, boundAnswerError("binding_unknown", "The request does not preserve the exact native Goal manifest.")
	}
	executionPath := filepath.Join(s.Request.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(draft.Task.Spec.ManifestSHA256, "sha256:")+".json")
	executionRaw, e := workflowBound(FileBinding{executionPath, draft.Task.Spec.ManifestSHA256}, 4<<20)
	if e != nil {
		return BoundAnswerScope{}, e
	}
	var execution struct {
		ContractKind string                `json:"contract_kind"`
		TaskID       workspace.TaskID      `json:"task_id"`
		SpecID       string                `json:"spec_id"`
		Revision     int                   `json:"revision"`
		Origin       workspace.TaskGoalRef `json:"goal_origin"`
	}
	if json.Unmarshal(executionRaw, &execution) != nil || execution.ContractKind != "execution" || execution.TaskID != goal.TaskID || execution.SpecID != draft.Task.SpecID || execution.Revision != draft.Task.Spec.Revision || execution.Origin.SpecID != goal.SpecID || execution.Origin.Spec.Revision != goal.Revision || execution.Origin.Spec.ManifestSHA256 != c.Goal.SHA256 {
		return BoundAnswerScope{}, boundAnswerError("binding_unknown", "The execution Spec does not retain this exact selected Goal origin.")
	}
	out := BoundAnswerScope{ProjectID: draft.Binding.ProjectID, RepoID: draft.Binding.RepoID, EpicID: string(c.Agreement.EpicID), TaskID: string(draft.Task.TaskID), GoalID: goal.SpecID, GoalRevision: goal.Revision, GoalSHA256: c.Goal.SHA256, RunRequestSHA256: s.Result.RequestSHA256, DeliveryContractSHA256: digest(c), TaskWorktree: s.Observed.Target.WorktreeLocator}
	for _, v := range []string{out.ProjectID, out.RepoID, out.EpicID, out.TaskID, out.GoalID} {
		if !key(v) {
			return out, boundAnswerError("binding_unknown", "A complete native Project/Repo/Epic/Task/Goal binding is required.")
		}
	}
	if out.GoalRevision < 1 || !digestPattern.MatchString(out.GoalSHA256) || out.TaskWorktree == s.Request.WorkspaceRoot || physical(out.TaskWorktree, false) != nil {
		return out, boundAnswerError("binding_unknown", "The Goal revision or distinct physical Task worktree is unavailable.")
	}
	if _, e := workflowBound(FileBinding{Locator: s.Result.Handoff.Locator, SHA256: s.Result.Handoff.SHA256}, 4<<20); e != nil {
		return out, e
	}
	return out, nil
}

func productQuestionDigest(q ProductQuestion) string {
	raw, _ := Canonical(q)
	var m map[string]json.RawMessage
	_ = json.Unmarshal(raw, &m)
	delete(m, "question_sha256")
	return digest(m)
}
func validateProductProposal(q *productQuestionProposal, s workflowState) error {
	if q == nil || q.Purpose != "product_choice" || !plain(q.Text, 1, 2048) || !plain(q.Reason, 1, 1024) || len(q.Options) < 2 || len(q.Options) > 5 || q.Continuation.Type != "resume_same_developer_thread" || q.Continuation.SourceIteration != s.Result.Worker.Iteration || s.Result.Worker.Job == nil || q.Continuation.SourceJobID != s.Result.Worker.Job.ID {
		return boundAnswerError("invalid_question", "Expected a bounded product choice tied to this developer turn.")
	}
	seen := map[string]bool{}
	for _, o := range q.Options {
		if !boundAnswerOptionID.MatchString(o.OptionID) || !plain(o.Label, 1, 120) || !plain(o.Consequence, 1, 1024) || seen[o.OptionID] {
			return boundAnswerError("invalid_question", "Question options must have unique identifiers and bounded text.")
		}
		seen[o.OptionID] = true
	}
	return nil
}

// The adapter validates the observed permission profile, including every root.
// A closed model turn is not evidence that detached tool writers are absent.
func productPolicy(s workflowState, job WorkerJob, raw json.RawMessage) error {
	var p struct {
		Type       string `json:"type"`
		FileSystem struct {
			Type    string `json:"type"`
			Entries []struct {
				Path struct {
					Type  string `json:"type"`
					Path  string `json:"path"`
					Value struct {
						Kind string `json:"kind"`
					} `json:"value"`
				} `json:"path"`
				Access string `json:"access"`
			} `json:"entries"`
		} `json:"file_system"`
		Network string `json:"network"`
	}
	if json.Unmarshal(raw, &p) != nil || p.Type != "managed" || p.FileSystem.Type != "restricted" || p.Network != "restricted" {
		return boundAnswerError("policy_changed", "Observed provider permissions are not the qualified managed profile.")
	}
	want := map[string]string{":root": "read", ":minimal": "read", s.Observed.Target.WorktreeLocator: "write", filepath.Join(s.Observed.Target.WorktreeLocator, ".git"): "read", s.Observed.Target.GitCommonDir: "read", filepath.Join(job.Directory, "output"): "write", filepath.Join(job.Directory, "tmp"): "write"}
	got := map[string]string{}
	for _, v := range p.FileSystem.Entries {
		k := v.Path.Path
		if v.Path.Type == "special" {
			k = ":" + v.Path.Value.Kind
		} else if v.Path.Type != "path" {
			return boundAnswerError("policy_changed", "Unknown filesystem policy entry.")
		}
		if _, exists := got[k]; exists {
			return boundAnswerError("policy_changed", "Duplicate filesystem policy entry.")
		}
		got[k] = v.Access
	}
	if !equal(got, want) {
		return boundAnswerError("policy_changed", "Observed filesystem permissions differ from this source attempt.")
	}
	return nil
}

var productSessionID = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)

func productSessionPath(in workerAttemptIntent, thread string) (string, error) {
	if !productSessionID.MatchString(thread) {
		return "", boundAnswerError("session_unknown", "Source provider session identifier is invalid.")
	}
	env := map[string]string{}
	for _, v := range in.Environment {
		k, v, ok := strings.Cut(v, "=")
		if ok {
			env[k] = v
		}
	}
	home := env["CODEX_HOME"]
	if home == "" && env["HOME"] != "" {
		home = filepath.Join(env["HOME"], ".codex")
	}
	if home == "" || physical(home, false) != nil {
		return "", boundAnswerError("session_unknown", "The source intent does not retain a physical Codex home.")
	}
	paths, e := filepath.Glob(filepath.Join(home, "sessions", "*", "*", "*", "*"+thread+".jsonl"))
	if e != nil || len(paths) != 1 {
		return "", boundAnswerError("session_unknown", "Exactly one physical source session artifact is required.")
	}
	return paths[0], nil
}
func productReadSession(s workflowState, job WorkerJob, thread, path string) (productSessionBinding, []byte, error) {
	return productReadSessionTurns(s, job, thread, path, 1)
}
func productReadSessionTurns(s workflowState, job WorkerJob, thread, path string, turns int) (productSessionBinding, []byte, error) {
	out := productSessionBinding{Path: path, SessionID: thread}
	raw, e := readFile(path, 16<<20, true)
	if e != nil {
		return out, nil, e
	}
	if len(raw) == 0 || raw[len(raw)-1] != '\n' {
		return out, nil, boundAnswerError("session_unknown", "Provider session artifact is incomplete.")
	}
	meta, contexts := 0, 0
	for _, line := range bytes.Split(bytes.TrimSuffix(raw, []byte{'\n'}), []byte{'\n'}) {
		var item struct {
			Type    string          `json:"type"`
			Payload json.RawMessage `json:"payload"`
		}
		if e = json.Unmarshal(line, &item); e != nil {
			return out, nil, e
		}
		if item.Type == "session_meta" {
			var p struct {
				ID  string `json:"id"`
				CWD string `json:"cwd"`
			}
			if json.Unmarshal(item.Payload, &p) != nil || p.ID != thread || p.CWD != s.Observed.Target.WorktreeLocator {
				return out, nil, boundAnswerError("session_changed", "Session metadata does not match the source turn.")
			}
			meta++
		}
		if item.Type == "turn_context" {
			var p struct {
				CWD      string          `json:"cwd"`
				Approval string          `json:"approval_policy"`
				Model    string          `json:"model"`
				Effort   string          `json:"effort"`
				Policy   json.RawMessage `json:"permission_profile"`
			}
			if json.Unmarshal(item.Payload, &p) != nil || p.CWD != s.Observed.Target.WorktreeLocator || p.Approval != "never" || p.Model != s.Request.Runtime.Model || p.Effort != s.Request.Delivery.ReasoningEffort {
				return out, nil, boundAnswerError("session_changed", "Observed session runtime differs from the bound model, effort, cwd or approval policy.")
			}
			if e = productPolicy(s, job, p.Policy); e != nil {
				return out, nil, e
			}
			out.CWD, out.Model, out.Effort, out.Policy = p.CWD, p.Model, p.Effort, p.Policy
			contexts++
		}
	}
	if meta != 1 || contexts != turns {
		return out, nil, boundAnswerError("session_unknown", "A first source turn requires one session metadata record and one observed turn context.")
	}
	out.Length = len(raw)
	return out, raw, nil
}
func productVerifySession(q productQuestionRecord, allowAppend bool) error {
	prefix, e := workflowBound(q.Session.Prefix, 16<<20)
	if e != nil {
		return e
	}
	raw, e := readFile(q.Session.Path, 32<<20, true)
	if e != nil {
		return e
	}
	if len(prefix) != q.Session.Length || len(raw) < len(prefix) || !bytes.Equal(raw[:len(prefix)], prefix) || !allowAppend && len(raw) != len(prefix) {
		return boundAnswerError("session_changed", "The preserved provider session prefix changed or another turn was appended.")
	}
	return nil
}
func productTaskFiles(root string) (string, error) {
	type entry struct {
		Path   string `json:"path"`
		Mode   uint32 `json:"mode"`
		SHA256 string `json:"sha256"`
	}
	entries := []entry{}
	var total int64
	e := filepath.WalkDir(root, func(path string, de fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		rel, e := filepath.Rel(root, path)
		if e != nil {
			return e
		}
		if rel == "." {
			return nil
		}
		if rel == ".git" {
			if de.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if len(entries) > 200000 {
			return fmt.Errorf("Task file inventory is too large")
		}
		info, e := de.Info()
		if e != nil {
			return e
		}
		if info.IsDir() {
			return nil
		}
		var raw []byte
		if info.Mode()&os.ModeSymlink != 0 {
			v, err := os.Readlink(path)
			if err != nil {
				return err
			}
			raw = []byte(v)
		} else if info.Mode().IsRegular() {
			raw, e = readFile(path, 64<<20, false)
			if e != nil {
				return e
			}
		} else {
			return fmt.Errorf("unsupported Task file type")
		}
		total += int64(len(raw))
		if total > 1<<30 {
			return fmt.Errorf("Task file inventory exceeds 1 GiB")
		}
		entries = append(entries, entry{rel, uint32(info.Mode()), hash(raw)})
		return nil
	})
	if e != nil {
		return "", e
	}
	return digest(entries), nil
}
func productObserveTarget(d Dependencies, s workflowState) (workspace.IntegrationWorktreeObservation, error) {
	if d.Workspace.IntegrationGit == nil {
		return workspace.IntegrationWorktreeObservation{}, boundAnswerError("binding_unknown", "Native Task Git observation is unavailable.")
	}
	x, e := d.Workspace.IntegrationGit.ObserveIntegrationWorktree(s.Observed.Target.WorktreeLocator, s.Observed.Target.Ref)
	if e != nil {
		return x, e
	}
	if x.Locator != s.Observed.Target.WorktreeLocator || x.Ref != s.Observed.Target.Ref || x.GitCommonDir != s.Observed.Target.GitCommonDir || len(x.InProgress) != 0 {
		return x, boundAnswerError("file_drift", "Task Git identity changed or an operation is in progress.")
	}
	return x, nil
}

func workerPublishProductQuestion(d Dependencies, s workflowState, r workerAttemptResult, proposal *productQuestionProposal) error {
	if !productQuestionsEnabled(s) || s.ProductQuestion != nil || s.Result.Worker.Stage != "developer" || s.Result.Worker.Job == nil || s.Result.Worker.Job.Kind != "developer" {
		return boundAnswerError("unsupported", "This worker cannot publish a product question.")
	}
	if e := validateProductProposal(proposal, s); e != nil {
		return e
	}
	qualification, e := productRuntimeQualification(d, s)
	if e != nil {
		return e
	}
	// Reparse the hash-bound stream at the publication boundary; a copied result
	// field or PID is never a substitute for the completed source turn.
	stream, e := workflowBound(r.Stdout, 16<<20)
	if e != nil {
		return e
	}
	thread, e := workerThreadID(stream)
	if e != nil || thread != r.ThreadID || r.Error != "" || r.Exit == nil || *r.Exit != 0 || r.ProviderFailure != nil {
		return boundAnswerError("source_unknown", "The source has no complete successful model-turn result.")
	}
	scope, e := productScope(s)
	if e != nil {
		return e
	}
	in, e := workerIntent(s, *s.Result.Worker.Job)
	if e != nil {
		return e
	}
	sessionPath, e := productSessionPath(in, thread)
	if e != nil {
		return e
	}
	session, prefix, e := productReadSession(s, *s.Result.Worker.Job, thread, sessionPath)
	if e != nil {
		return e
	}
	queue, e := productNativeQueue(d, s)
	if e != nil {
		return e
	}
	target, e := productObserveTarget(d, s)
	if e != nil {
		return e
	}
	files, e := productTaskFiles(target.Locator)
	if e != nil {
		return e
	}
	source, e := func() (FileBinding, error) {
		b, _, err := deliveryReadBinding(filepath.Join(s.Result.Worker.Job.Directory, "result.json"), 1<<20)
		return b, err
	}()
	if e != nil {
		return e
	}
	q := ProductQuestion{QuestionID: "pqw_" + strings.TrimPrefix(digest([]string{s.Result.RunID, s.Result.Worker.Job.IntentHash, r.Answer.SHA256}), "sha256:")[:32], QuestionRevision: 1, WaitGeneration: s.Result.Worker.Generation, Purpose: proposal.Purpose, Text: proposal.Text, Reason: proposal.Reason, Options: proposal.Options, Continuation: ProductQuestionContinuation{Type: "resume_same_developer_thread", SourceJobID: s.Result.Worker.Job.ID, SourceIteration: s.Result.Worker.Iteration, SourceSessionID: thread, DispatchPolicy: "subject_to_parent_control", MaxDispatches: 1, Boundary: "completed_model_turn_same_owner", GlobalWriterQuiescence: "not_attested"}}
	q.QuestionSHA256 = productQuestionDigest(q)
	return WithWorkerParentGuard(s.Request.WorkspaceRoot, s.Request.ControlContract, func() error {
		return workflowUpdate(d, s.Request.WorkspaceRoot, s.Result.RunID, func(current *workflowState) error {
			if e := workerLiveOwner(d, *current); e != nil {
				return e
			}
			if current.ProductQuestion != nil || !equal(current.Result.Worker.Job, s.Result.Worker.Job) {
				return boundAnswerError("stale", "Source reservation changed before question publication.")
			}
			gate, e := ReadParentControl(current.Request.WorkspaceRoot, current.Request.ControlContract)
			if e != nil {
				return e
			}
			if gate.ActiveRun != current.Result.RunID {
				return boundAnswerError("owner_unknown", "The parent does not retain this exact Run and dispatch reservation.")
			}
			dir := filepath.Join(workerDirectory(*current), "product-questions", q.QuestionID)
			session.Prefix = FileBinding{Locator: filepath.Join(dir, "source-session-prefix.jsonl"), SHA256: hash(prefix)}
			if e = d.writeOnce(session.Prefix.Locator, prefix); e != nil {
				return e
			}
			rec := productQuestionRecord{q, scope, *current.Result.Worker.Job, source, session, files, digest(target), queue, gate.Deadline, qualification, ControlTimestamp(controlNow(d))}
			b, e := workflowKeep(d, filepath.Join(dir, "question.json"), rec)
			if e != nil {
				return e
			}
			current.ProductQuestion = &b
			current.Result.Worker.Phase = "awaiting_product_answer"
			current.Result.Worker.Reason = ""
			current.Result.NextAction = productAnswerRequired
			return nil
		})
	})
}

func productQuestionRecordRead(s workflowState) (productQuestionRecord, error) {
	var q productQuestionRecord
	if s.ProductQuestion == nil {
		return q, os.ErrNotExist
	}
	raw, e := workflowBound(*s.ProductQuestion, 2<<20)
	if e != nil {
		return q, e
	}
	if e = decode(raw, 2<<20, &q); e != nil {
		return q, e
	}
	if q.Question.QuestionSHA256 != productQuestionDigest(q.Question) || q.Scope.RunRequestSHA256 != s.Result.RequestSHA256 {
		return q, boundAnswerError("binding_unknown", "Question identity changed.")
	}
	return q, nil
}
