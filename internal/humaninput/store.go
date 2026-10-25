package humaninput

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

// Layout below <workspace>/.ply/human-input/v1. Each fact is its own
// write-once file, so no crash can leave a half-applied multi-file change:
//
//	questions/<id>.json      immutable question
//	answers/<id>.json        the one winning answer request
//	receipts/<id>.json       agent-reported receipt of that answer
//	continuations/<id>.json  agent-reported use of that answer
//	withdrawals/<id>.json    agent withdrawal of an unanswered question
//	store.lock               held only while one write decides; never while waiting
func storeRoot(root string) string { return filepath.Join(root, ".ply", "human-input", "v1") }
func factPath(root, kind, id string) string {
	return filepath.Join(storeRoot(root), kind, id+".json")
}
func questionPath(root, id string) string     { return factPath(root, "questions", id) }
func answerPath(root, id string) string       { return factPath(root, "answers", id) }
func receiptPath(root, id string) string      { return factPath(root, "receipts", id) }
func continuationPath(root, id string) string { return factPath(root, "continuations", id) }
func withdrawalPath(root, id string) string   { return factPath(root, "withdrawals", id) }

const maxFactBytes = 512 << 10

type questionRecord struct {
	Workspace      string   `json:"workspace"`
	QuestionID     string   `json:"question_id"`
	QuestionSHA256 string   `json:"question_sha256"`
	PublicationKey string   `json:"publication_key"`
	Scope          Scope    `json:"scope"`
	Source         Source   `json:"source"`
	Text           string   `json:"text"`
	Reason         string   `json:"reason"`
	AnswerKind     string   `json:"answer_kind"`
	Options        []Option `json:"options"`
	CreatedAtUTC   string   `json:"created_at_utc"`
}

type answerRecord struct {
	Request       AnswerRequest `json:"request"`
	RequestSHA256 string        `json:"request_sha256"`
	RecordedAtUTC string        `json:"recorded_at_utc"`
}

func canonical(v any) []byte {
	b, e := taskrun.Canonical(v)
	if e != nil {
		panic(e)
	}
	return b
}
func digestOf(v any) string {
	x := sha256.Sum256(canonical(v))
	return "sha256:" + hex.EncodeToString(x[:])
}

// The question identity is stable for one workspace, project and key, so a
// resumed agent re-publishing the same key finds the same question.
func questionIDFor(root, project, key string) string {
	x := sha256.Sum256(canonical(map[string]string{"purpose": "human_input_question_identity_v1", "workspace": root, "project_id": project, "publication_key": key}))
	return "hiq_" + hex.EncodeToString(x[:16])
}

// question_sha256 binds scope, title, text, reason, kind and options; agent
// source facts and later status observations are deliberately outside it.
func questionDigest(r questionRecord) string {
	return digestOf(struct {
		Purpose        string   `json:"purpose"`
		Workspace      string   `json:"workspace"`
		QuestionID     string   `json:"question_id"`
		PublicationKey string   `json:"publication_key"`
		Scope          Scope    `json:"scope"`
		Text           string   `json:"text"`
		Reason         string   `json:"reason"`
		AnswerKind     string   `json:"answer_kind"`
		Options        []Option `json:"options"`
	}{"human_input_question_v1", r.Workspace, r.QuestionID, r.PublicationKey, r.Scope, r.Text, r.Reason, r.AnswerKind, r.Options})
}

func bindingOf(r AnswerRequest) string {
	return digestOf(struct {
		Kind           string      `json:"kind"`
		SchemaVersion  int         `json:"schema_version"`
		Mode           string      `json:"mode"`
		Purpose        string      `json:"purpose"`
		Workspace      string      `json:"workspace"`
		QuestionID     string      `json:"question_id"`
		QuestionSHA256 string      `json:"question_sha256"`
		Scope          Scope       `json:"scope"`
		RequestID      string      `json:"request_id"`
		Answer         AnswerValue `json:"answer"`
	}{r.Kind, r.SchemaVersion, r.Mode, r.Purpose, r.Workspace, r.QuestionID, r.QuestionSHA256, r.Scope, r.RequestID, r.Answer})
}

// RequestSHA256 is the digest of the complete canonical request.
func RequestSHA256(r AnswerRequest) string { return digestOf(r) }

func sameString(a, b *string) bool { return (a == nil) == (b == nil) && (a == nil || *a == *b) }
func sameScope(a, b Scope) bool {
	return a.ProjectID == b.ProjectID && sameString(a.RepoID, b.RepoID) && sameString(a.TaskID, b.TaskID) && a.ActivityID == b.ActivityID && a.ActivityTitle == b.ActivityTitle
}

func invalidState(detail string) error {
	return fail(CodeStateInvalid, detail+"; nothing was repaired or changed")
}

// readFact decodes one optional private fact. A missing file is (false, nil);
// a symlink, wrong type, public mode or malformed content is an error.
func readFact(path string, out any) (bool, error) {
	raw, e := taskrun.PrivateRead(path, maxFactBytes)
	if os.IsNotExist(e) {
		return false, nil
	}
	if e != nil {
		return false, invalidState("unsafe or unreadable human-input file " + filepath.Base(filepath.Dir(path)) + "/" + filepath.Base(path) + ": " + e.Error())
	}
	if e = taskrun.StrictDecode(raw, maxFactBytes, out); e != nil {
		return false, invalidState("malformed human-input file " + filepath.Base(filepath.Dir(path)) + "/" + filepath.Base(path))
	}
	return true, nil
}

func readQuestionRecord(root, id string) (questionRecord, error) {
	var r questionRecord
	found, e := readFact(questionPath(root, id), &r)
	if e != nil {
		return r, e
	}
	if !found {
		return r, fail(CodeNotFound, "no question "+id+" exists in this workspace")
	}
	draft := QuestionDraft{Kind: DraftKind, SchemaVersion: SchemaVersion, PublicationKey: r.PublicationKey, ProjectID: r.Scope.ProjectID, RepoID: r.Scope.RepoID, TaskID: r.Scope.TaskID, ActivityID: r.Scope.ActivityID, ActivityTitle: r.Scope.ActivityTitle, Provider: r.Source.Provider, SessionID: r.Source.SessionID, OriginCWD: r.Source.OriginCWD, Text: r.Text, Reason: r.Reason, AnswerKind: r.AnswerKind, Options: r.Options}
	if r.Workspace != root || r.QuestionID != id || questionIDFor(root, r.Scope.ProjectID, r.PublicationKey) != id || r.QuestionSHA256 != questionDigest(r) || r.Source.Provenance != ProvenanceAgent || validateDraft(draft) != nil {
		return r, invalidState("question " + id + " does not match its immutable identity or digest")
	}
	return r, nil
}

// view derives the current public question from its immutable facts. Reads
// take no lock and create nothing.
func view(root, id string) (Question, error) {
	r, e := readQuestionRecord(root, id)
	if e != nil {
		return Question{}, e
	}
	q := Question{Kind: QuestionKind, SchemaVersion: SchemaVersion, Mode: Mode, Workspace: root, QuestionID: id, QuestionSHA256: r.QuestionSHA256, PublicationKey: r.PublicationKey, Scope: r.Scope, Source: r.Source, Text: r.Text, Reason: r.Reason, AnswerKind: r.AnswerKind, Options: r.Options, CreatedAtUTC: r.CreatedAtUTC, State: StateWaiting}
	var a answerRecord
	if found, e := readFact(answerPath(root, id), &a); e != nil {
		return q, e
	} else if found {
		req := a.Request
		if validateRequestShape(req) != nil || req.Workspace != root || req.QuestionID != id || req.QuestionSHA256 != r.QuestionSHA256 || !sameScope(req.Scope, r.Scope) || a.RequestSHA256 != digestOf(req) || validateAnswer(q, req.Answer) != nil {
			return q, invalidState("the recorded answer of " + id + " does not bind this question")
		}
		q.Answer = &Answer{RequestID: req.RequestID, RequestSHA256: a.RequestSHA256, AnswerKind: req.Answer.AnswerKind, OptionID: req.Answer.OptionID, Text: req.Answer.Text, RecordedAtUTC: a.RecordedAtUTC}
		q.State = StateAnswered
	}
	var receipt Receipt
	if found, e := readFact(receiptPath(root, id), &receipt); e != nil {
		return q, e
	} else if found {
		if q.Answer == nil || receipt.RequestID != q.Answer.RequestID || receipt.RequestSHA256 != q.Answer.RequestSHA256 || receipt.Provenance != ProvenanceAgent {
			return q, invalidState("the receipt of " + id + " does not match its recorded answer")
		}
		q.Receipt, q.State = &receipt, StateReceived
	}
	var cont Continuation
	if found, e := readFact(continuationPath(root, id), &cont); e != nil {
		return q, e
	} else if found {
		if q.Receipt == nil || cont.RequestID != q.Answer.RequestID || cont.RequestSHA256 != q.Answer.RequestSHA256 || cont.Provenance != ProvenanceAgent || !ValidReport(cont.Summary) {
			return q, invalidState("the continuation report of " + id + " does not match its receipt")
		}
		q.Continuation, q.State = &cont, StateContinued
	}
	var w Withdrawal
	if found, e := readFact(withdrawalPath(root, id), &w); e != nil {
		return q, e
	} else if found {
		if q.Answer != nil || w.Provenance != ProvenanceAgent || !ValidReport(w.Reason) {
			return q, invalidState("question " + id + " is both answered and withdrawn or has an invalid withdrawal")
		}
		q.Withdrawal, q.State = &w, StateWithdrawn
	}
	q.NextAction = nextAction(q)
	return q, nil
}

const notTracked = " Progress of the whole assignment and any later human steps are not tracked here."

func nextAction(q Question) NextAction {
	switch q.State {
	case StateWaiting:
		if q.AnswerKind == AnswerText {
			return NextAction{ActorHuman, "Write one short text answer to this question."}
		}
		return NextAction{ActorHuman, "Choose one of the listed options."}
	case StateAnswered:
		return NextAction{ActorAgent, "The answer is recorded. The agent has not reported receiving it; an interrupted agent must be resumed before it can use the answer."}
	case StateReceived:
		return NextAction{ActorAgent, "The agent reported receiving the answer and has not yet reported using it."}
	case StateContinued:
		return NextAction{ActorNone, "The agent reported that it continued its work with this answer." + notTracked}
	default:
		return NextAction{ActorNone, "The agent withdrew this unanswered question." + notTracked}
	}
}

var inventoryName = regexp.MustCompile(`^hiq_[a-f0-9]{32}\.json$`)

// questionIDs lists question IDs in their immutable sort order. Interrupted
// publication temporaries are not questions; any other entry is unsafe state.
func questionIDs(root string) ([]string, error) {
	dir := filepath.Join(storeRoot(root), "questions")
	e := taskrun.PhysicalPath(dir)
	if os.IsNotExist(e) {
		return []string{}, nil
	}
	if e != nil {
		return nil, invalidState("unsafe question inventory: " + e.Error())
	}
	info, e := os.Lstat(dir)
	if e != nil || !info.IsDir() || info.Mode().Perm()&0077 != 0 {
		return nil, invalidState("the question inventory is not a private directory")
	}
	entries, e := os.ReadDir(dir)
	if e != nil {
		return nil, invalidState("the question inventory cannot be read")
	}
	ids := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".publish-") {
			continue
		}
		if !inventoryName.MatchString(name) || !entry.Type().IsRegular() {
			return nil, invalidState("unexpected entry in the question inventory: " + name)
		}
		ids = append(ids, strings.TrimSuffix(name, ".json"))
	}
	sort.Strings(ids)
	return ids, nil
}

// ReadQuestionBlockers lists every question once for the agent follow-up veto.
// blocked holds the sources of open questions and of validated questions whose
// later facts cannot be read, even when that partial view is already terminal.
// An unidentifiable question or inventory, or a question to block without a
// session, could be any agent's, so coverage is unknown. partial reports that
// some history was invalid. Reads take no lock and create or repair nothing.
func ReadQuestionBlockers(root string) (blocked []Source, coverageUnknown, partial bool) {
	ids, e := questionIDs(root)
	if e != nil {
		return nil, true, true
	}
	for _, id := range ids {
		q, e := view(root, id)
		partial = partial || e != nil
		switch {
		case q.QuestionID == "":
			coverageUnknown = true
		case e == nil && (q.State == StateContinued || q.State == StateWithdrawn):
		case q.Source.SessionID == nil:
			coverageUnknown = true
		default:
			blocked = append(blocked, q.Source)
		}
	}
	return blocked, coverageUnknown, partial
}

// publish writes one immutable fact under the held store lock.
func publish(path string, v any) error {
	return taskrun.PrivateWriteOnce(path, canonical(v))
}
func withLock(root string, fn func() error) error {
	return taskrun.WithPrivateLock(storeRoot(root), fn)
}
