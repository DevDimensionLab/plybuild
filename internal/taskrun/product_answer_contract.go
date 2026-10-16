package taskrun

import (
	"errors"
	"os"
	"regexp"
	"time"
)

// ProductQuestionContract opts a newly registered worker into one same-session
// answer turn. It never upgrades an existing schema-1 request.
type ProductQuestionContract struct {
	SchemaVersion        int    `json:"schema_version"`
	Mode                 string `json:"mode"`
	MaxContinuationTurns int    `json:"max_continuation_turns"`
}

func ValidateProductQuestionContract(c *ProductQuestionContract) error {
	if c == nil {
		return nil
	}
	if c.SchemaVersion != 1 || c.Mode != BoundAnswerMode || c.MaxContinuationTurns != 1 {
		return boundAnswerError("unsupported", "Product questions require schema 1, same-owner turn mode, and exactly one continuation.")
	}
	return nil
}

func (e *BoundAnswerReason) Error() string { return e.Code + ": " + e.Detail }
func boundAnswerError(code, detail string) error {
	return &BoundAnswerReason{"bound_answer_" + code, detail}
}
func BoundAnswerFailureOf(err error) BoundAnswerReason {
	var reason *BoundAnswerReason
	if errors.As(err, &reason) {
		return *reason
	}
	return BoundAnswerReason{"bound_answer_outcome_unknown", "Preserved answer evidence is unavailable or inconsistent; retain this request identity."}
}

var boundAnswerRequestID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{7,127}$`)
var boundAnswerRunID = regexp.MustCompile(`^wfr_[a-f0-9]{64}$`)
var boundAnswerQuestionID = regexp.MustCompile(`^pqw_[a-f0-9]{32}$`)
var boundAnswerOptionID = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)
var boundAnswerTimestamp = regexp.MustCompile(`^\d{4}-\d\d-\d\dT\d\d:\d\d:\d\d(?:\.\d{1,9})?Z$`)

func boundAnswerTime(s string) (time.Time, error) {
	if !boundAnswerTimestamp.MatchString(s) {
		return time.Time{}, boundAnswerError("invalid_arguments", "Expected an explicit UTC timestamp.")
	}
	t, e := time.Parse(time.RFC3339Nano, s)
	if e != nil {
		return t, boundAnswerError("invalid_arguments", "Invalid UTC timestamp.")
	}
	return t, nil
}

// ParseBoundAnswerRequest rejects duplicate/unknown/missing keys, noncanonical
// number tokens, malformed UTF-8 and extra JSON before any store is consulted.
func ParseBoundAnswerRequest(raw []byte) (BoundAnswerRequest, error) {
	var r BoundAnswerRequest
	if e := decode(raw, 64<<10, &r); e != nil {
		return r, boundAnswerError("invalid_request", "Expected one closed bound-answer request of at most 64 KiB.")
	}
	if r.Kind != "ply.workflow.bound-answer-request" || r.SchemaVersion != BoundAnswerSchemaVersion || r.Mode != BoundAnswerMode || r.Action != "answer_product_question" || !boundAnswerRequestID.MatchString(r.RequestID) || !boundAnswerRunID.MatchString(r.RunID) || !boundAnswerQuestionID.MatchString(r.QuestionID) || r.QuestionRevision != 1 || !boundAnswerOptionID.MatchString(r.OptionID) || !digestPattern.MatchString(r.Binding) || (r.DispatchPolicy != "when_parent_allows" && r.DispatchPolicy != "deferred_until_explicit_parent_resume") {
		return r, boundAnswerError("invalid_request", "Request identity, mode or dispatch policy is invalid.")
	}
	if len(r.Workspace) > 4096 || physical(r.Workspace, false) != nil {
		return r, boundAnswerError("invalid_request", "Workspace must name the absolute physical workspace directory.")
	}
	if info, e := os.Stat(r.Workspace); e != nil || !info.IsDir() {
		return r, boundAnswerError("invalid_request", "Workspace must be an existing physical directory.")
	}
	issued, e := boundAnswerTime(r.Validity.IssuedAtUTC)
	if e != nil {
		return r, boundAnswerError("invalid_request", "Invalid request timestamp.")
	}
	expires, e := boundAnswerTime(r.Validity.ExpiresAtUTC)
	if e != nil {
		return r, boundAnswerError("invalid_request", "Invalid request timestamp.")
	}
	deadline, e := boundAnswerTime(r.OriginalDeadlineUTC)
	if e != nil {
		return r, boundAnswerError("invalid_request", "Invalid request timestamp.")
	}
	if !issued.Before(expires) || expires.After(deadline) || expires.Sub(issued) > 300*time.Second {
		return r, boundAnswerError("invalid_request", "Request validity must be positive, at most 300 seconds, and within the original deadline.")
	}
	return r, nil
}

func ReadBoundAnswerRequest(path string) (BoundAnswerRequest, error) {
	raw, e := readFile(path, 64<<10, false)
	if e != nil {
		return BoundAnswerRequest{}, boundAnswerError("invalid_request", "Cannot read the bounded physical request file.")
	}
	return ParseBoundAnswerRequest(raw)
}
