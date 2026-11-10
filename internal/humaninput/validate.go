package humaninput

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/privatestore"
)

const (
	CodeInvalidArguments    = "human_input_invalid_arguments"
	CodeInvalidDraft        = "human_input_invalid_draft"
	CodeInvalidRequest      = "human_input_invalid_request"
	CodeInvalidAnswer       = "human_input_invalid_answer"
	CodeScopeInvalid        = "human_input_scope_invalid"
	CodeScopeUnavailable    = "human_input_scope_unavailable"
	CodePublicationConflict = "human_input_publication_conflict"
	CodeNotFound            = "human_input_not_found"
	CodeRequestMismatch     = "human_input_request_mismatch"
	CodeRequestConflict     = "human_input_request_conflict"
	CodeAlreadyAnswered     = "human_input_already_answered"
	CodeWithdrawn           = "human_input_withdrawn"
	CodeNotAnswered         = "human_input_not_answered"
	CodeAnswerMismatch      = "human_input_answer_mismatch"
	CodeReportConflict      = "human_input_report_conflict"
	CodeStateInvalid        = "human_input_state_invalid"
	CodeOutcomeUnknown      = "human_input_outcome_unknown"
)

// Error is a stable, semantic refusal. Question is the current question when
// known; it is never synthesized from the refusal itself.
type Error struct {
	Code     string
	Detail   string
	Question *Question
}

func (e *Error) Error() string { return e.Code + ": " + e.Detail }
func fail(code, detail string) *Error {
	return &Error{Code: code, Detail: detail}
}
func failOn(q Question, code, detail string) *Error {
	return &Error{Code: code, Detail: detail, Question: &q}
}

var (
	questionIDPattern = regexp.MustCompile(`^hiq_[a-f0-9]{32}$`)
	requestIDPattern  = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{7,127}$`)
	optionIDPattern   = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)
	keyPattern        = regexp.MustCompile(`^[a-z0-9][a-z0-9._/-]{0,127}$`)
	providerPattern   = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,31}$`)
	projectPattern    = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)
	repoTaskPattern   = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,127}$`)
	digestPattern     = regexp.MustCompile(`^sha256:[a-f0-9]{64}$`)
)

func ValidQuestionID(s string) bool { return questionIDPattern.MatchString(s) }
func ValidRequestID(s string) bool  { return requestIDPattern.MatchString(s) }
func ValidOptionID(s string) bool   { return optionIDPattern.MatchString(s) }
func ValidDigest(s string) bool     { return digestPattern.MatchString(s) }
func ValidProjectID(s string) bool  { return projectPattern.MatchString(s) }

func stableKey(s string) bool {
	if !keyPattern.MatchString(s) || strings.Contains(s, "//") || strings.HasSuffix(s, "/") {
		return false
	}
	for _, part := range strings.Split(s, "/") {
		if part == "." || part == ".." {
			return false
		}
	}
	return true
}

func blank(s string) bool { return strings.TrimFunc(s, unicode.IsSpace) == "" }

// line is a one-line display value: no control characters and no surrounding
// whitespace, so it cannot be confused with neighbouring UI text.
func line(s string, max int) bool {
	if !utf8.ValidString(s) || blank(s) || strings.TrimSpace(s) != s || utf8.RuneCountInString(s) > max {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

// prose is agent-written multi-line text: tab and newline are its only
// control characters. It is preserved exactly, never trimmed or normalized.
func prose(s string, max int) bool {
	if !utf8.ValidString(s) || blank(s) || utf8.RuneCountInString(s) > max {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return false
		}
	}
	return true
}

// ValidAnswerText checks a human text answer: valid UTF-8, not blank, no NUL,
// at most 4000 codepoints and 16000 bytes. Other content is kept byte-exact.
func ValidAnswerText(s string) bool {
	return utf8.ValidString(s) && !blank(s) && !strings.ContainsRune(s, 0) && utf8.RuneCountInString(s) <= 4000 && len(s) <= 16000
}

// ValidReport checks an agent-reported continue summary or withdraw reason.
func ValidReport(s string) bool { return prose(s, 1000) }

func validateDraft(q QuestionDraft) error {
	bad := func(detail string) error { return fail(CodeInvalidDraft, detail) }
	switch {
	case q.Kind != DraftKind || q.SchemaVersion != SchemaVersion:
		return bad("expected kind " + DraftKind + " with schema_version 1")
	case !stableKey(q.PublicationKey):
		return bad("publication_key must be a stable lowercase key of at most 128 characters")
	case !projectPattern.MatchString(q.ProjectID):
		return bad("project_id is not a valid project ID")
	case q.RepoID != nil && !repoTaskPattern.MatchString(*q.RepoID):
		return bad("repo_id is not a valid repository ID")
	case q.TaskID != nil && !repoTaskPattern.MatchString(*q.TaskID):
		return bad("task_id is not a valid Task ID")
	case !stableKey(q.ActivityID):
		return bad("activity_id must be a stable lowercase key of at most 128 characters")
	case !line(q.ActivityTitle, 160):
		return bad("activity_title must be one non-blank line of at most 160 codepoints")
	case !providerPattern.MatchString(q.Provider):
		return bad("provider must be a short lowercase provider name such as claude or codex")
	case q.SessionID != nil && !line(*q.SessionID, 256):
		return bad("session_id must be null or one line of at most 256 codepoints")
	case !prose(q.Text, 2000):
		return bad("text must be non-blank UTF-8 of at most 2000 codepoints without control characters other than tab and newline")
	case !prose(q.Reason, 1000):
		return bad("reason must be non-blank UTF-8 of at most 1000 codepoints without control characters other than tab and newline")
	case q.Options == nil:
		return bad("options must be an array")
	}
	switch q.AnswerKind {
	case AnswerText:
		if len(q.Options) != 0 {
			return bad("a text question has an empty options array")
		}
	case AnswerChoice:
		if len(q.Options) < 2 || len(q.Options) > 6 {
			return bad("a choice question has 2 to 6 options")
		}
		ids, labels := map[string]bool{}, map[string]bool{}
		for _, o := range q.Options {
			if !optionIDPattern.MatchString(o.ID) || ids[o.ID] {
				return bad("option IDs must be unique lowercase IDs of at most 64 characters")
			}
			if !line(o.Label, 120) || labels[o.Label] {
				return bad("option labels must be unique non-blank lines of at most 120 codepoints")
			}
			if !prose(o.Consequence, 1000) {
				return bad("option consequences must be non-blank text of at most 1000 codepoints")
			}
			ids[o.ID], labels[o.Label] = true, true
		}
	default:
		return bad("answer_kind must be choice or text")
	}
	return nil
}

// absoluteInput checks an explicit CLI input path before it is read.
func absoluteInput(path string) error {
	if len(path) == 0 || len(path) > 4096 || !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return fmt.Errorf("must be an absolute normalized path")
	}
	return privatestore.PhysicalPath(path)
}

// readInput reads one explicit regular file without following symlinks. Input
// files are caller-owned and need not be private.
func readInput(path string, max int) ([]byte, error) {
	if e := absoluteInput(path); e != nil {
		return nil, e
	}
	info, e := os.Lstat(path)
	if e != nil {
		return nil, e
	}
	if !info.Mode().IsRegular() || info.Size() > int64(max) {
		return nil, fmt.Errorf("must be a regular file of at most %d bytes", max)
	}
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	// Lstat before and Stat after open detect a swapped path or symlink.
	opened, e := f.Stat()
	if e != nil || !os.SameFile(info, opened) {
		return nil, fmt.Errorf("file changed while opening")
	}
	b, e := io.ReadAll(io.LimitReader(f, int64(max)+1))
	if e != nil {
		return nil, e
	}
	if len(b) > max {
		return nil, fmt.Errorf("file exceeds %d bytes", max)
	}
	return b, nil
}

const maxDocumentBytes = 256 << 10

// ReadDraft reads one closed question draft from an explicit physical file.
func ReadDraft(path string) (QuestionDraft, error) {
	var q QuestionDraft
	raw, e := readInput(path, maxDocumentBytes)
	if e != nil {
		return q, fail(CodeInvalidDraft, "--file "+e.Error())
	}
	if e = privatestore.StrictDecode(raw, maxDocumentBytes, &q); e != nil {
		return q, fail(CodeInvalidDraft, "expected one closed question draft: "+e.Error())
	}
	return q, nil
}

// ReadTextFile reads an exact text answer; bytes are never trimmed or normalized.
func ReadTextFile(path string) (string, error) {
	raw, e := readInput(path, 16000)
	if e != nil {
		return "", fail(CodeInvalidAnswer, "--text-file "+e.Error())
	}
	if !ValidAnswerText(string(raw)) {
		return "", fail(CodeInvalidAnswer, "text must be non-blank UTF-8 without NUL, at most 4000 codepoints and 16000 bytes")
	}
	return string(raw), nil
}

// ReadAnswerRequest reads one unchanged preview request from an explicit file.
func ReadAnswerRequest(path string) (AnswerRequest, error) {
	var r AnswerRequest
	raw, e := readInput(path, maxDocumentBytes)
	if e != nil {
		return r, fail(CodeInvalidRequest, "--file "+e.Error())
	}
	if e = privatestore.StrictDecode(raw, maxDocumentBytes, &r); e != nil {
		return r, fail(CodeInvalidRequest, "expected one closed answer request: "+e.Error())
	}
	return r, validateRequestShape(r)
}

func validateRequestShape(r AnswerRequest) error {
	bad := func(detail string) error { return fail(CodeInvalidRequest, detail) }
	switch {
	case r.Kind != RequestKind || r.SchemaVersion != SchemaVersion || r.Mode != Mode:
		return bad("expected kind " + RequestKind + ", schema_version 1 and mode " + Mode)
	case r.Purpose != PurposeAnswer:
		return bad("purpose must be " + PurposeAnswer)
	case absoluteInput(r.Workspace) != nil:
		return bad("workspace must be the physical absolute workspace")
	case !questionIDPattern.MatchString(r.QuestionID) || !digestPattern.MatchString(r.QuestionSHA256) || !requestIDPattern.MatchString(r.RequestID) || !digestPattern.MatchString(r.BindingSHA256):
		return bad("question_id, question_sha256, request_id or binding_sha256 is malformed")
	case r.BindingSHA256 != bindingOf(r):
		return bad("binding_sha256 does not bind this exact request")
	}
	return nil
}
