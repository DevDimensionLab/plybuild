package taskrun

// Bound-answer v2 is the frozen same-owner turn wire contract. Nullable fields
// remain explicit JSON null; public objects deliberately have no omitempty.
const BoundAnswerMode = "week_codex_same_owner_turn_v1"
const BoundAnswerSchemaVersion = 2

type BoundAnswerReason struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

type BoundAnswerEvidence struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	SHA256        string `json:"sha256"`
	ObservedAtUTC string `json:"observed_at_utc"`
}

type BoundAnswerValidity struct {
	IssuedAtUTC  string `json:"issued_at_utc"`
	ExpiresAtUTC string `json:"expires_at_utc"`
}

type BoundAnswerScope struct {
	ProjectID              string `json:"project_id"`
	RepoID                 string `json:"repo_id"`
	EpicID                 string `json:"epic_id"`
	TaskID                 string `json:"task_id"`
	GoalID                 string `json:"goal_id"`
	GoalRevision           int    `json:"goal_revision"`
	GoalSHA256             string `json:"goal_sha256"`
	RunRequestSHA256       string `json:"run_request_sha256"`
	DeliveryContractSHA256 string `json:"delivery_contract_sha256"`
	TaskWorktree           string `json:"task_worktree"`
}

type BoundAnswerParentWeek struct {
	WeekID              string `json:"week_id"`
	ContractSHA256      string `json:"contract_sha256"`
	WorksetSHA256       string `json:"workset_sha256"`
	ControlRevision     int    `json:"control_revision"`
	DesiredAction       string `json:"desired_action"`
	OriginalDeadlineUTC string `json:"original_deadline_utc"`
}

type BoundAnswerWorker struct {
	OwnerID                 string `json:"owner_id"`
	OwnerGeneration         int    `json:"owner_generation"`
	OwnerState              string `json:"owner_state"`
	Iteration               int    `json:"iteration"`
	SourceJobID             string `json:"source_job_id"`
	SourceIntentSHA256      string `json:"source_intent_sha256"`
	Provider                string `json:"provider"`
	ProviderSessionID       string `json:"provider_session_id"`
	ProviderBinarySHA256    string `json:"provider_binary_sha256"`
	Model                   string `json:"model"`
	Effort                  string `json:"effort"`
	PolicySHA256            string `json:"policy_sha256"`
	SessionBindingSHA256    string `json:"session_binding_sha256"`
	QualificationSHA256     string `json:"qualification_sha256"`
	ContinuationLimit       int    `json:"continuation_limit"`
	ContinuationsDispatched int    `json:"continuations_dispatched"`
}

type ProductQuestionOption struct {
	OptionID    string `json:"option_id"`
	Label       string `json:"label"`
	Consequence string `json:"consequence"`
}

type ProductQuestionContinuation struct {
	Type                   string `json:"type"`
	SourceJobID            string `json:"source_job_id"`
	SourceIteration        int    `json:"source_iteration"`
	SourceSessionID        string `json:"source_session_id"`
	DispatchPolicy         string `json:"dispatch_policy"`
	MaxDispatches          int    `json:"max_dispatches"`
	Boundary               string `json:"boundary"`
	GlobalWriterQuiescence string `json:"global_writer_quiescence"`
}

type ProductQuestion struct {
	QuestionID       string                      `json:"question_id"`
	QuestionRevision int                         `json:"question_revision"`
	QuestionSHA256   string                      `json:"question_sha256"`
	WaitGeneration   int                         `json:"wait_generation"`
	Purpose          string                      `json:"purpose"`
	Text             string                      `json:"text"`
	Reason           string                      `json:"reason"`
	Options          []ProductQuestionOption     `json:"options"`
	Continuation     ProductQuestionContinuation `json:"continuation"`
}

type ProductQuestionReadback struct {
	Kind                string                 `json:"kind"`
	SchemaVersion       int                    `json:"schema_version"`
	Mode                string                 `json:"mode"`
	Workspace           string                 `json:"workspace"`
	RunID               string                 `json:"run_id"`
	Scope               *BoundAnswerScope      `json:"scope"`
	ParentWeek          *BoundAnswerParentWeek `json:"parent_week"`
	Worker              *BoundAnswerWorker     `json:"worker"`
	OriginalDeadlineUTC *string                `json:"original_deadline_utc"`
	QuestionState       string                 `json:"question_state"`
	Eligibility         string                 `json:"eligibility"`
	Question            *ProductQuestion       `json:"question"`
	Reasons             []BoundAnswerReason    `json:"reasons"`
	ObservedAtUTC       string                 `json:"observed_at_utc"`
	Freshness           string                 `json:"freshness"`
	Evidence            []BoundAnswerEvidence  `json:"evidence"`
}

type BoundAnswerRequest struct {
	Kind                string              `json:"kind"`
	SchemaVersion       int                 `json:"schema_version"`
	Mode                string              `json:"mode"`
	Action              string              `json:"action"`
	RequestID           string              `json:"request_id"`
	Workspace           string              `json:"workspace"`
	RunID               string              `json:"run_id"`
	QuestionID          string              `json:"question_id"`
	QuestionRevision    int                 `json:"question_revision"`
	OptionID            string              `json:"option_id"`
	OriginalDeadlineUTC string              `json:"original_deadline_utc"`
	Validity            BoundAnswerValidity `json:"validity"`
	Binding             string              `json:"binding"`
	DispatchPolicy      string              `json:"dispatch_policy"`
}

type BoundAnswerPreview struct {
	Kind          string                   `json:"kind"`
	SchemaVersion int                      `json:"schema_version"`
	Mode          string                   `json:"mode"`
	Workspace     string                   `json:"workspace"`
	RunID         string                   `json:"run_id"`
	Eligibility   string                   `json:"eligibility"`
	Basis         *ProductQuestionReadback `json:"basis"`
	Request       *BoundAnswerRequest      `json:"request"`
	RequestSHA256 *string                  `json:"request_sha256"`
	Reasons       []BoundAnswerReason      `json:"reasons"`
	ObservedAtUTC string                   `json:"observed_at_utc"`
	Freshness     string                   `json:"freshness"`
	Evidence      []BoundAnswerEvidence    `json:"evidence"`
}

type BoundAnswerReceipt struct {
	RequestID                   string  `json:"request_id"`
	RequestSHA256               string  `json:"request_sha256"`
	QuestionID                  string  `json:"question_id"`
	QuestionRevision            int     `json:"question_revision"`
	QuestionSHA256              string  `json:"question_sha256"`
	OptionID                    string  `json:"option_id"`
	SourceJobID                 string  `json:"source_job_id"`
	SourceIntentSHA256          string  `json:"source_intent_sha256"`
	SourceSessionID             string  `json:"source_session_id"`
	RecordedAtUTC               string  `json:"recorded_at_utc"`
	ContinuationJobID           *string `json:"continuation_job_id"`
	ContinuationIntentSHA256    *string `json:"continuation_intent_sha256"`
	ContinuationSessionID       *string `json:"continuation_session_id"`
	ContinuationStartedAtUTC    *string `json:"continuation_started_at_utc"`
	ContinuationCompletedAtUTC  *string `json:"continuation_completed_at_utc"`
	ResultSHA256                *string `json:"result_sha256"`
	ConsumedAnswerRequestID     *string `json:"consumed_answer_request_id"`
	ConsumedAnswerRequestSHA256 *string `json:"consumed_answer_request_sha256"`
}

type BoundAnswerStatus struct {
	Kind                  string                 `json:"kind"`
	SchemaVersion         int                    `json:"schema_version"`
	Mode                  string                 `json:"mode"`
	Workspace             string                 `json:"workspace"`
	RequestID             string                 `json:"request_id"`
	Action                string                 `json:"action"`
	RequestSHA256         *string                `json:"request_sha256"`
	Binding               *string                `json:"binding"`
	RunID                 *string                `json:"run_id"`
	Scope                 *BoundAnswerScope      `json:"scope"`
	ParentWeek            *BoundAnswerParentWeek `json:"parent_week"`
	Worker                *BoundAnswerWorker     `json:"worker"`
	OriginalDeadlineUTC   *string                `json:"original_deadline_utc"`
	QuestionID            *string                `json:"question_id"`
	QuestionRevision      *int                   `json:"question_revision"`
	QuestionSHA256        *string                `json:"question_sha256"`
	OptionID              *string                `json:"option_id"`
	State                 string                 `json:"state"`
	AnswerRecorded        *bool                  `json:"answer_recorded"`
	ContinuationRequested *bool                  `json:"continuation_requested"`
	ExecutionEffect       string                 `json:"execution_effect"`
	Receipt               *BoundAnswerReceipt    `json:"receipt"`
	Reasons               []BoundAnswerReason    `json:"reasons"`
	ObservedAtUTC         string                 `json:"observed_at_utc"`
	Freshness             string                 `json:"freshness"`
	Evidence              []BoundAnswerEvidence  `json:"evidence"`
}
