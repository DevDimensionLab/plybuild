package humaninput

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// HI01: an ordinary named activity, without a new Task, publishes, lists and
// shows one valid question; an existing Task scope is checked against the registry.
func TestHI01OrdinaryActivityPublishesListsAndShowsWithoutTask(t *testing.T) {
	f := newFixture(t)
	before, err := os.ReadFile(filepath.Join(f.root, ".ply", "work-items.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	q := f.ask(f.choice("conversation-1/format"))
	if !questionIDPattern.MatchString(q.QuestionID) || q.State != StateWaiting || q.NextAction.Actor != ActorHuman || q.Scope.TaskID != nil || q.Scope.RepoID != nil || q.Workspace != f.root {
		t.Fatalf("unexpected published question: %+v", q)
	}
	if q.Source.Provenance != ProvenanceAgent || q.Source.Provider != "claude" || q.Answer != nil || q.Receipt != nil || q.Continuation != nil || q.Withdrawal != nil {
		t.Fatalf("source or nullable observations are wrong: %+v", q)
	}
	after, err := os.ReadFile(filepath.Join(f.root, ".ply", "work-items.yaml"))
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("publishing a question must not create or rewrite a Task registration")
	}
	shown, err := Show(f.d, f.root, q.QuestionID)
	if err != nil || shown.QuestionSHA256 != q.QuestionSHA256 || shown.Text != q.Text {
		t.Fatalf("show: %+v %v", shown, err)
	}
	list, err := List(f.d, f.root, ListQuery{ProjectID: "ply"})
	if err != nil || len(list.Questions) != 1 || list.Questions[0].QuestionID != q.QuestionID || !list.Complete || list.HasMore || list.NextCursor != nil || list.Limit != 20 {
		t.Fatalf("list: %+v %v", list, err)
	}
	if other, err := List(f.d, f.root, ListQuery{ProjectID: "other"}); err != nil || len(other.Questions) != 0 || !other.Complete {
		t.Fatalf("project scope leaked: %+v %v", other, err)
	}
}

func TestHI01ExistingTaskScopeIsValidatedAgainstRegistry(t *testing.T) {
	f := newFixture(t)
	ok := f.choice("task/one")
	ok.RepoID, ok.TaskID = ptr("ply"), ptr("task")
	if q := f.ask(ok); q.Scope.TaskID == nil || *q.Scope.TaskID != "task" || *q.Scope.RepoID != "ply" {
		t.Fatalf("task scope not preserved: %+v", q.Scope)
	}
	cases := map[string]func(*QuestionDraft){
		"unknown task":      func(q *QuestionDraft) { q.RepoID, q.TaskID = ptr("ply"), ptr("missing") },
		"task without repo": func(q *QuestionDraft) { q.TaskID = ptr("task") },
		"task in other project": func(q *QuestionDraft) {
			q.ProjectID, q.RepoID, q.TaskID, q.OriginCWD = "other", ptr("other"), ptr("task"), f.otherOrigin
		},
		"repo not in project": func(q *QuestionDraft) { q.RepoID = ptr("other") },
		"unknown project":     func(q *QuestionDraft) { q.ProjectID = "missing" },
		"origin outside":      func(q *QuestionDraft) { q.OriginCWD = f.otherOrigin },
		"origin missing":      func(q *QuestionDraft) { q.OriginCWD = filepath.Join(f.origin, "absent") },
		"origin relative":     func(q *QuestionDraft) { q.OriginCWD = "ply/main" },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			q := f.choice("task/" + strings.ReplaceAll(name, " ", "-"))
			change(&q)
			before := snapshot(t, filepath.Join(f.root, ".ply"))
			_, err := AskDraft(f.d, f.root, q)
			wantCode(t, err, CodeScopeInvalid)
			sameSnapshot(t, before, snapshot(t, filepath.Join(f.root, ".ply")))
		})
	}
	planning := filepath.Join(f.root, "ply", "planning")
	if err := os.Mkdir(planning, 0o755); err != nil {
		t.Fatal(err)
	}
	companion := f.choice("task/wrapper-companion")
	companion.OriginCWD = planning
	f.ask(companion)
}

func TestHI01DraftValidationRejectsBeforeAnyStoreIsCreated(t *testing.T) {
	f := newFixture(t)
	long := func(n int) string { return strings.Repeat("é", n) }
	cases := map[string]func(*QuestionDraft){
		"kind":              func(q *QuestionDraft) { q.Kind = "other" },
		"schema":            func(q *QuestionDraft) { q.SchemaVersion = 2 },
		"key":               func(q *QuestionDraft) { q.PublicationKey = "../escape" },
		"activity":          func(q *QuestionDraft) { q.ActivityID = "" },
		"title long":        func(q *QuestionDraft) { q.ActivityTitle = long(161) },
		"title control":     func(q *QuestionDraft) { q.ActivityTitle = "a\nb" },
		"provider":          func(q *QuestionDraft) { q.Provider = "Claude Code" },
		"text long":         func(q *QuestionDraft) { q.Text = long(2001) },
		"text blank":        func(q *QuestionDraft) { q.Text = " \n\t " },
		"text nul":          func(q *QuestionDraft) { q.Text = "a\x00b" },
		"reason long":       func(q *QuestionDraft) { q.Reason = long(1001) },
		"kind unknown":      func(q *QuestionDraft) { q.AnswerKind = "number" },
		"one option":        func(q *QuestionDraft) { q.Options = q.Options[:1] },
		"seven options":     func(q *QuestionDraft) { q.Options = append(q.Options, sevenMore()...) },
		"duplicate id":      func(q *QuestionDraft) { q.Options[1].ID = q.Options[0].ID },
		"duplicate label":   func(q *QuestionDraft) { q.Options[1].Label = q.Options[0].Label },
		"label long":        func(q *QuestionDraft) { q.Options[0].Label = long(121) },
		"consequence blank": func(q *QuestionDraft) { q.Options[0].Consequence = " " },
		"consequence long":  func(q *QuestionDraft) { q.Options[0].Consequence = long(1001) },
		"text with options": func(q *QuestionDraft) { q.AnswerKind = AnswerText },
		"option id":         func(q *QuestionDraft) { q.Options[0].ID = "Short" },
		"invalid utf8":      func(q *QuestionDraft) { q.Text = "a\xffb" },
		"session control":   func(q *QuestionDraft) { q.SessionID = ptr("a\tb") },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			q := f.choice("draft/x")
			q.Options = append([]Option{}, q.Options...)
			change(&q)
			_, err := AskDraft(f.d, f.root, q)
			wantCode(t, err, CodeInvalidDraft)
			if _, err := os.Lstat(storeRoot(f.root)); !os.IsNotExist(err) {
				t.Fatal("a rejected draft created the store")
			}
		})
	}
	max := f.choice("draft/max")
	max.Text, max.Reason, max.ActivityTitle = long(2000), long(1000), long(160)
	f.ask(max)
}

func sevenMore() []Option {
	var out []Option
	for i := 0; i < 5; i++ {
		out = append(out, Option{ID: fmt.Sprintf("o%d", i), Label: fmt.Sprintf("Option %d", i), Consequence: "x"})
	}
	return out
}

func TestHI01ReadDraftRejectsUnsafeOrMalformedFiles(t *testing.T) {
	f := newFixture(t)
	dir := filepath.Join(f.root, "drafts")
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	good := filepath.Join(dir, "good.json")
	if err := os.WriteFile(good, mustJSON(t, f.choice("file/good")), 0o600); err != nil {
		t.Fatal(err)
	}
	if q, err := ReadDraft(good); err != nil || q.PublicationKey != "file/good" {
		t.Fatalf("valid draft: %+v %v", q, err)
	}
	link := filepath.Join(dir, "link.json")
	if err := os.Symlink(good, link); err != nil {
		t.Fatal(err)
	}
	raw := string(mustJSON(t, f.choice("file/x")))
	write := func(name, body string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
		return p
	}
	for name, path := range map[string]string{
		"symlink":        link,
		"directory":      dir,
		"relative":       "drafts/good.json",
		"unknown field":  write("unknown.json", strings.Replace(raw, `{`, `{"extra":1,`, 1)),
		"missing null":   write("missing.json", strings.Replace(raw, `"task_id":null,`, ``, 1)),
		"duplicate":      write("dup.json", strings.Replace(raw, `{`, `{"text":"a",`, 1)),
		"lone surrogate": write("surrogate.json", strings.Replace(raw, `Which output`, `\ud800 output`, 1)),
		"trailing":       write("trailing.json", raw+" {}"),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := ReadDraft(path)
			wantCode(t, err, CodeInvalidDraft)
		})
	}
}

// HI02: idempotent publication, conflicts, concurrency and complete pagination.
func TestHI02IdempotentPublicationAndConflictingContent(t *testing.T) {
	f := newFixture(t)
	q := f.ask(f.choice("idem/one"))
	resumed := f.choice("idem/one")
	resumed.SessionID, resumed.Provider = ptr("resumed-session"), "codex"
	again := f.ask(resumed)
	if again.QuestionID != q.QuestionID || again.QuestionSHA256 != q.QuestionSHA256 || again.CreatedAtUTC != q.CreatedAtUTC || again.Source.Provider != "claude" {
		t.Fatalf("same key and content must find the original question: %+v", again)
	}
	before := snapshot(t, storeRoot(f.root))
	changed := f.choice("idem/one")
	changed.Text = "Which output should this task produce now?"
	_, err := AskDraft(f.d, f.root, changed)
	wantCode(t, err, CodePublicationConflict)
	sameSnapshot(t, before, snapshot(t, storeRoot(f.root)))
	if list, _ := List(f.d, f.root, ListQuery{ProjectID: "ply"}); len(list.Questions) != 1 {
		t.Fatalf("conflict created a duplicate: %d", len(list.Questions))
	}
}

func TestHI02ConcurrentPublicationKeepsEveryQuestionOnce(t *testing.T) {
	f := newFixture(t)
	var wg sync.WaitGroup
	ids := make([]string, 16)
	errs := make([]error, 16)
	for i := range ids {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("concurrent/%d", i%8) // two writers per key
			q, err := AskDraft(f.d, f.root, f.choice(key))
			ids[i], errs[i] = q.QuestionID, err
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Fatalf("writer %d: %v", i, err)
		}
		if ids[i] != ids[(i+8)%16] {
			t.Fatal("concurrent writers of one key got different questions")
		}
	}
	list, err := List(f.d, f.root, ListQuery{ProjectID: "ply"})
	if err != nil || len(list.Questions) != 8 {
		t.Fatalf("want 8 distinct questions, got %d %v", len(list.Questions), err)
	}
}

func TestHI02PaginationTraversesEveryStateWithoutLossOrDuplicate(t *testing.T) {
	f := newFixture(t)
	want := map[string]string{}
	for i := 0; i < 7; i++ {
		q := f.ask(f.choice(fmt.Sprintf("page/%d", i)))
		want[q.QuestionID] = StateWaiting
	}
	ids := sortedKeys(want)
	r := f.answer(ids[0], "request-answered-0", "short")
	want[ids[0]] = StateAnswered
	r = f.answer(ids[1], "request-received-1", "full")
	if _, err := Acknowledge(f.d, f.root, ids[1], r.RequestID, digestOf(r)); err != nil {
		t.Fatal(err)
	}
	want[ids[1]] = StateReceived
	r = f.answer(ids[2], "request-continued-2", "full")
	if _, err := Acknowledge(f.d, f.root, ids[2], r.RequestID, digestOf(r)); err != nil {
		t.Fatal(err)
	}
	if _, err := Continue(f.d, f.root, ids[2], r.RequestID, digestOf(r), "Wrote the full version."); err != nil {
		t.Fatal(err)
	}
	want[ids[2]] = StateContinued
	if _, err := Withdraw(f.d, f.root, ids[3], "No longer needed."); err != nil {
		t.Fatal(err)
	}
	want[ids[3]] = StateWithdrawn
	for _, limit := range []int{1, 2, 3, 20} {
		seen := map[string]string{}
		var cursor *string
		var last string
		for pages := 0; ; pages++ {
			page, err := List(f.d, f.root, ListQuery{ProjectID: "ply", Limit: limit, Cursor: cursor})
			if err != nil || pages > 10 || len(page.Questions) > limit {
				t.Fatalf("page: %+v %v", page, err)
			}
			for _, q := range page.Questions {
				if _, dup := seen[q.QuestionID]; dup || q.QuestionID <= last {
					t.Fatalf("duplicate or unordered row %s", q.QuestionID)
				}
				seen[q.QuestionID], last = q.State, q.QuestionID
			}
			if page.HasMore == page.Complete || (page.NextCursor == nil) != page.Complete {
				t.Fatalf("inconsistent paging flags: %+v", page)
			}
			if page.Complete {
				break
			}
			cursor = page.NextCursor
		}
		if fmt.Sprint(seen) != fmt.Sprint(want) {
			t.Fatalf("limit %d lost states: %v want %v", limit, seen, want)
		}
	}
	answered := StateAnswered
	if page, err := List(f.d, f.root, ListQuery{ProjectID: "ply", State: &answered}); err != nil || len(page.Questions) != 1 || page.Questions[0].QuestionID != ids[0] {
		t.Fatalf("state filter: %+v %v", page, err)
	}
	for _, bad := range []ListQuery{{ProjectID: "ply", Limit: 21}, {ProjectID: "ply", Limit: -1}, {ProjectID: "Bad"}, {ProjectID: "ply", Cursor: ptr("not-a-cursor")}, {ProjectID: "ply", State: ptr("pending")}} {
		_, err := List(f.d, f.root, bad)
		wantCode(t, err, CodeInvalidArguments)
	}
}

// The consumer output bound is checked with the agreed maximal text fields.
func TestHI02MaximalPageFitsConsumerOutputLimit(t *testing.T) {
	f := newFixture(t)
	worst := strings.Repeat("\U0001F600", 1)
	for i := 0; i < 20; i++ {
		q := f.choice(fmt.Sprintf("max/%02d", i))
		q.Text, q.Reason, q.ActivityTitle = strings.Repeat("\U0001F600", 2000), strings.Repeat("\"", 1000), strings.Repeat(worst, 160)
		q.Options = nil
		for o := 0; o < 6; o++ {
			q.Options = append(q.Options, Option{ID: fmt.Sprintf("o%d", o), Label: strings.Repeat(string(rune('a'+o)), 120), Consequence: strings.Repeat("\\", 1000)})
		}
		if i%2 == 1 {
			q.AnswerKind, q.Options = AnswerText, []Option{}
		}
		created := f.ask(q)
		if i%2 == 1 {
			text := strings.Repeat("\x01", 4000)
			p := f.preview(created.QuestionID, fmt.Sprintf("request-max-%02d", i), AnswerValue{AnswerKind: AnswerText, Text: &text})
			if _, err := Apply(f.d, f.root, p.Request); err != nil {
				t.Fatal(err)
			}
			if _, err := Acknowledge(f.d, f.root, created.QuestionID, p.Request.RequestID, p.RequestSHA256); err != nil {
				t.Fatal(err)
			}
			if _, err := Continue(f.d, f.root, created.QuestionID, p.Request.RequestID, p.RequestSHA256, strings.Repeat("\"", 1000)); err != nil {
				t.Fatal(err)
			}
		}
	}
	page, err := List(f.d, f.root, ListQuery{ProjectID: "ply"})
	if err != nil || len(page.Questions) != 20 {
		t.Fatalf("max page: %d %v", len(page.Questions), err)
	}
	body, err := Encode(page)
	if err != nil || len(body) > MaxOutputBytes {
		t.Fatalf("maximal page is %d bytes (limit %d): %v", len(body), MaxOutputBytes, err)
	}
}

// HI03: preview and record choice and exact short text; invalid answers are
// rejected before any effect.
func TestHI03ChoiceAndExactTextPreviewAndRecord(t *testing.T) {
	f := newFixture(t)
	c := f.ask(f.choice("answer/choice"))
	before := snapshot(t, f.root)
	p := f.preview(c.QuestionID, "request-choice-1", AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("full")})
	sameSnapshot(t, before, snapshot(t, f.root))
	if p.Request.Purpose != PurposeAnswer || p.Request.QuestionSHA256 != c.QuestionSHA256 || p.RequestSHA256 != digestOf(p.Request) || p.Request.BindingSHA256 != bindingOf(p.Request) || p.Request.Answer.Text != nil {
		t.Fatalf("preview request is not bound: %+v", p)
	}
	s, err := Apply(f.d, f.root, p.Request)
	if err != nil || s.State != StatusRecorded || s.Question == nil || s.Question.State != StateAnswered || s.Question.NextAction.Actor != ActorAgent || *s.Question.Answer.OptionID != "full" {
		t.Fatalf("apply: %+v %v", s, err)
	}

	tq := f.ask(f.text("answer/text"))
	exact := "  Línea uno\r\nline twó  \n\n"
	tp := f.preview(tq.QuestionID, "request-text-1", AnswerValue{AnswerKind: AnswerText, Text: &exact})
	if _, err := Apply(f.d, f.root, tp.Request); err != nil {
		t.Fatal(err)
	}
	shown, err := Show(f.d, f.root, tq.QuestionID)
	if err != nil || shown.Answer == nil || shown.Answer.Text == nil || *shown.Answer.Text != exact || shown.Answer.OptionID != nil {
		t.Fatalf("text answer was not preserved exactly: %+v %v", shown.Answer, err)
	}
	maxText := strings.Repeat("é", 4000) // 8000 bytes
	tq2 := f.ask(f.text("answer/text-max"))
	f.preview(tq2.QuestionID, "request-text-max", AnswerValue{AnswerKind: AnswerText, Text: &maxText})
	bytesMax := strings.Repeat("\U0001F600", 4000) // 4000 codepoints, 16000 bytes
	f.preview(tq2.QuestionID, "request-text-bytes", AnswerValue{AnswerKind: AnswerText, Text: &bytesMax})
}

func TestHI03InvalidAnswersAreRejectedBeforeChange(t *testing.T) {
	f := newFixture(t)
	c := f.ask(f.choice("reject/choice"))
	tq := f.ask(f.text("reject/text"))
	cases := []struct {
		name, id string
		v        AnswerValue
	}{
		{"unknown option", c.QuestionID, AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("medium")}},
		{"text for choice", c.QuestionID, AnswerValue{AnswerKind: AnswerText, Text: ptr("x")}},
		{"option for text", tq.QuestionID, AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("short")}},
		{"blank", tq.QuestionID, AnswerValue{AnswerKind: AnswerText, Text: ptr("")}},
		{"whitespace", tq.QuestionID, AnswerValue{AnswerKind: AnswerText, Text: ptr(" \r\n\t　")}},
		{"too many codepoints", tq.QuestionID, AnswerValue{AnswerKind: AnswerText, Text: ptr(strings.Repeat("a", 4001))}},
		{"too many bytes", tq.QuestionID, AnswerValue{AnswerKind: AnswerText, Text: ptr(strings.Repeat("\U0001F600", 4000) + "a")}},
		{"nul", tq.QuestionID, AnswerValue{AnswerKind: AnswerText, Text: ptr("a\x00")}},
		{"invalid utf8", tq.QuestionID, AnswerValue{AnswerKind: AnswerText, Text: ptr("a\xc3")}},
		{"both values", tq.QuestionID, AnswerValue{AnswerKind: AnswerText, Text: ptr("a"), OptionID: ptr("short")}},
	}
	before := snapshot(t, f.root)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := PreviewAnswer(f.d, f.root, tc.id, "request-reject-1", tc.v)
			wantCode(t, err, CodeInvalidAnswer)
		})
	}
	sameSnapshot(t, before, snapshot(t, f.root))

	good := f.preview(c.QuestionID, "request-tamper-1", AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("short")})
	tamper := map[string]func(*AnswerRequest){
		"changed option keeps binding": func(r *AnswerRequest) { r.Answer.OptionID = ptr("full") },
		"changed question digest": func(r *AnswerRequest) {
			r.QuestionSHA256 = "sha256:" + strings.Repeat("0", 64)
			r.BindingSHA256 = bindingOf(*r)
		},
		"changed scope":          func(r *AnswerRequest) { r.Scope.ActivityID = "other/activity"; r.BindingSHA256 = bindingOf(*r) },
		"other question":         func(r *AnswerRequest) { r.QuestionID = tq.QuestionID; r.BindingSHA256 = bindingOf(*r) },
		"unknown option rebound": func(r *AnswerRequest) { r.Answer.OptionID = ptr("medium"); r.BindingSHA256 = bindingOf(*r) },
		"wrong purpose":          func(r *AnswerRequest) { r.Purpose = "product_answer_v1"; r.BindingSHA256 = bindingOf(*r) },
	}
	for name, change := range tamper {
		t.Run(name, func(t *testing.T) {
			r := good.Request
			r.Scope = good.Request.Scope
			change(&r)
			s, err := Apply(f.d, f.root, r)
			if err == nil || s.State == StatusRecorded {
				t.Fatalf("tampered request was accepted: %+v", s)
			}
		})
	}
	sameSnapshot(t, before, snapshot(t, f.root))
	if q, _ := Show(f.d, f.root, c.QuestionID); q.State != StateWaiting {
		t.Fatal("a rejected request changed the question")
	}
}

func TestHI03RemovedScopeBlocksNewAnswersAndKeepsHistory(t *testing.T) {
	f := newFixture(t)
	other := func(key string) QuestionDraft {
		q := f.choice(key)
		q.ProjectID, q.OriginCWD = "other", f.otherOrigin
		return q
	}
	answered := f.ask(other("scope/answered"))
	r := f.answer(answered.QuestionID, "request-scope-1", "short")
	waiting := f.ask(other("scope/waiting"))
	p := f.preview(waiting.QuestionID, "request-scope-2", AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("short")})
	path := filepath.Join(f.root, ".ply", "projects.yaml")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	// Remove project other (it has no Tasks) as an operator would.
	removed := strings.Replace(string(raw), "  - id: other\n    name: Other\n    wrapper: "+filepath.Join(f.root, "other")+"\n    repo_ids:\n      - other\n", "", 1)
	removed = strings.Replace(removed, "  - id: other\n    locator: "+f.otherOrigin+"\n    git_common_dir: "+filepath.Join(f.otherOrigin, ".git")+"\n", "", 1)
	if removed == string(raw) {
		t.Fatalf("fixture registry layout changed:\n%s", raw)
	}
	if err := os.WriteFile(path, []byte(removed), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err = PreviewAnswer(f.d, f.root, waiting.QuestionID, "request-scope-3", AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("short")})
	wantCode(t, err, CodeScopeUnavailable)
	s, err := Apply(f.d, f.root, p.Request)
	if err == nil || s.State != StatusNotRecorded || codeOfReasons(s) != CodeScopeUnavailable {
		t.Fatalf("removed scope accepted an answer: %+v %v", s, err)
	}
	if s, err := Status(f.d, f.root, answered.QuestionID, r.RequestID); err != nil || s.State != StatusRecorded {
		t.Fatalf("history lost after scope removal: %+v %v", s, err)
	}
	if list, err := List(f.d, f.root, ListQuery{ProjectID: "other"}); err != nil || len(list.Questions) != 2 {
		t.Fatalf("history not listed after scope removal: %+v %v", list, err)
	}
}

func codeOfReasons(s AnswerStatus) string {
	if len(s.Reasons) == 0 {
		return ""
	}
	return s.Reasons[0].Code
}

// HI04: one immutable winner under concurrency, idempotent replay and recovery.
func TestHI04ConcurrentApplyHasOneImmutableWinner(t *testing.T) {
	f := newFixture(t)
	q := f.ask(f.choice("winner/one"))
	var requests []AnswerRequest
	for i := 0; i < 12; i++ {
		option := []string{"short", "full"}[i%2]
		requests = append(requests, f.preview(q.QuestionID, fmt.Sprintf("request-race-%02d", i), AnswerValue{AnswerKind: AnswerChoice, OptionID: &option}).Request)
	}
	results := make([]AnswerStatus, len(requests))
	var wg sync.WaitGroup
	for i := range requests {
		wg.Add(1)
		go func(i int) { defer wg.Done(); results[i], _ = Apply(f.d, f.root, requests[i]) }(i)
	}
	wg.Wait()
	winners := 0
	var winner AnswerRequest
	for i, s := range results {
		switch s.State {
		case StatusRecorded:
			winners++
			winner = requests[i]
		case StatusOtherWinner:
			if s.Winner == nil {
				t.Fatalf("loser does not name the winner: %+v", s)
			}
		default:
			t.Fatalf("unexpected concurrent outcome: %+v", s)
		}
	}
	if winners != 1 {
		t.Fatalf("want one winner, got %d", winners)
	}
	for i, s := range results {
		if s.State == StatusOtherWinner && (s.Winner.RequestID != winner.RequestID || s.Winner.RequestSHA256 != digestOf(winner) || s.RequestID != requests[i].RequestID) {
			t.Fatalf("loser names a wrong winner: %+v", s)
		}
	}
	before := snapshot(t, storeRoot(f.root))
	if s, err := Apply(f.d, f.root, winner); err != nil || s.State != StatusRecorded {
		t.Fatalf("replay of the winner is not idempotent: %+v %v", s, err)
	}
	for _, r := range requests {
		if r.RequestID == winner.RequestID {
			continue
		}
		if s, err := Apply(f.d, f.root, r); err == nil || s.State != StatusOtherWinner {
			t.Fatalf("another request overwrote the winner: %+v", s)
		}
		if s, _ := Status(f.d, f.root, q.QuestionID, r.RequestID); s.State != StatusOtherWinner || s.Winner.RequestID != winner.RequestID {
			t.Fatalf("status of a loser: %+v", s)
		}
	}
	sameRequestOtherBytes := winner
	sameRequestOtherBytes.Answer.OptionID = ptr(map[string]string{"short": "full", "full": "short"}[*winner.Answer.OptionID])
	sameRequestOtherBytes.BindingSHA256 = bindingOf(sameRequestOtherBytes)
	if s, err := Apply(f.d, f.root, sameRequestOtherBytes); err == nil || s.State != StatusOtherWinner || codeOfReasons(s) != CodeRequestConflict {
		t.Fatalf("same request ID with other bytes: %+v %v", s, err)
	}
	sameSnapshot(t, before, snapshot(t, storeRoot(f.root)))
}

func TestHI04LostReturnRecoversTheSameOutcome(t *testing.T) {
	f := newFixture(t)
	q := f.ask(f.choice("recovery/one"))
	p := f.preview(q.QuestionID, "request-lost-1", AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("short")})
	if s, err := Status(f.d, f.root, q.QuestionID, p.Request.RequestID); err != nil || s.State != StatusNotRecorded || s.Winner != nil {
		t.Fatalf("before apply: %+v %v", s, err)
	}
	lost := f.d
	lost.Fault = func(point string) error {
		if point == "after_answer_publish" {
			return errors.New("simulated lost CLI return")
		}
		return nil
	}
	if _, err := Apply(lost, f.root, p.Request); err == nil {
		t.Fatal("fault did not surface")
	}
	s, err := Status(f.d, f.root, q.QuestionID, p.Request.RequestID)
	if err != nil || s.State != StatusRecorded || s.RequestSHA256 == nil || *s.RequestSHA256 != p.RequestSHA256 {
		t.Fatalf("status after lost return: %+v %v", s, err)
	}
	if s, err := Apply(f.d, f.root, p.Request); err != nil || s.State != StatusRecorded {
		t.Fatalf("same-ID recovery: %+v %v", s, err)
	}

	q2 := f.ask(f.choice("recovery/two"))
	p2 := f.preview(q2.QuestionID, "request-lost-2", AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("full")})
	early := f.d
	early.Fault = func(point string) error {
		if point == "before_answer_publish" {
			return errors.New("simulated crash before publication")
		}
		return nil
	}
	if _, err := Apply(early, f.root, p2.Request); err == nil {
		t.Fatal("fault did not surface")
	}
	if s, _ := Status(f.d, f.root, q2.QuestionID, p2.Request.RequestID); s.State != StatusNotRecorded {
		t.Fatalf("authoritative not_recorded after an early crash: %+v", s)
	}
	// Fresh confirmation of the SAME immutable request may now be applied.
	if s, err := Apply(f.d, f.root, p2.Request); err != nil || s.State != StatusRecorded {
		t.Fatalf("same request after not_recorded: %+v %v", s, err)
	}
}

// HI05: wait is a read-only wait; acknowledge/continue need the right winner.
func TestHI05WaitTimeoutIsReadOnlyAndReturnsTheAnswer(t *testing.T) {
	f := newFixture(t)
	q := f.ask(f.choice("wait/one"))
	before := snapshot(t, f.root)
	w, err := Wait(f.d, f.root, q.QuestionID, 30*time.Millisecond)
	if err != nil || w.Outcome != WaitStillWaiting || w.Question.State != StateWaiting || w.Question.NextAction.Actor != ActorHuman {
		t.Fatalf("timeout: %+v %v", w, err)
	}
	sameSnapshot(t, before, snapshot(t, f.root))
	go func() {
		time.Sleep(40 * time.Millisecond)
		p, _ := PreviewAnswer(f.d, f.root, q.QuestionID, "request-wait-1", AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("short")})
		_, _ = Apply(f.d, f.root, p.Request)
	}()
	w, err = Wait(f.d, f.root, q.QuestionID, 5*time.Second)
	if err != nil || w.Outcome != WaitAnswered || w.Question.Answer == nil || w.Question.Answer.RequestID != "request-wait-1" {
		t.Fatalf("wait did not return the answer: %+v %v", w, err)
	}
	gone := f.ask(f.choice("wait/withdrawn"))
	if _, err := Withdraw(f.d, f.root, gone.QuestionID, "Resolved from the goal."); err != nil {
		t.Fatal(err)
	}
	if w, err := Wait(f.d, f.root, gone.QuestionID, time.Second); err != nil || w.Outcome != WaitWithdrawn {
		t.Fatalf("withdrawn wait: %+v %v", w, err)
	}
}

func TestHI05AcknowledgeAndContinueRequireTheRecordedWinner(t *testing.T) {
	f := newFixture(t)
	q := f.ask(f.choice("report/one"))
	_, err := Acknowledge(f.d, f.root, q.QuestionID, "request-report-1", "sha256:"+strings.Repeat("a", 64))
	wantCode(t, err, CodeNotAnswered)
	r := f.answer(q.QuestionID, "request-report-1", "full")
	sha := digestOf(r)
	_, err = Continue(f.d, f.root, q.QuestionID, r.RequestID, sha, "Used it.")
	wantCode(t, err, CodeNotAcknowledged)
	_, err = Acknowledge(f.d, f.root, q.QuestionID, r.RequestID, "sha256:"+strings.Repeat("b", 64))
	wantCode(t, err, CodeAnswerMismatch)
	_, err = Acknowledge(f.d, f.root, q.QuestionID, "request-report-2", sha)
	wantCode(t, err, CodeAnswerMismatch)
	got, err := Acknowledge(f.d, f.root, q.QuestionID, r.RequestID, sha)
	if err != nil || got.State != StateReceived || got.Receipt == nil || got.Receipt.Provenance != ProvenanceAgent || got.NextAction.Actor != ActorAgent {
		t.Fatalf("acknowledge: %+v %v", got, err)
	}
	if again, err := Acknowledge(f.d, f.root, q.QuestionID, r.RequestID, sha); err != nil || again.Receipt.ReportedAtUTC != got.Receipt.ReportedAtUTC {
		t.Fatalf("acknowledge replay is not idempotent: %+v %v", again, err)
	}
	_, err = Continue(f.d, f.root, q.QuestionID, r.RequestID, sha, strings.Repeat("x", 1001))
	wantCode(t, err, CodeInvalidArguments)
	done, err := Continue(f.d, f.root, q.QuestionID, r.RequestID, sha, "Wrote the full version into README.md.")
	if err != nil || done.State != StateContinued || done.NextAction.Actor != ActorNone || done.Continuation.Provenance != ProvenanceAgent || !strings.Contains(done.NextAction.Message, "not tracked") {
		t.Fatalf("continue: %+v %v", done, err)
	}
	if _, err := Continue(f.d, f.root, q.QuestionID, r.RequestID, sha, "Wrote the full version into README.md."); err != nil {
		t.Fatalf("continue replay: %v", err)
	}
	_, err = Continue(f.d, f.root, q.QuestionID, r.RequestID, sha, "Something else.")
	wantCode(t, err, CodeReportConflict)
	_, err = Withdraw(f.d, f.root, q.QuestionID, "Too late.")
	wantCode(t, err, CodeAlreadyAnswered)
}

func TestHI05WithdrawIsTerminalAndKeepsTheQuestion(t *testing.T) {
	f := newFixture(t)
	q := f.ask(f.choice("withdraw/one"))
	p := f.preview(q.QuestionID, "request-withdraw-1", AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("short")})
	w, err := Withdraw(f.d, f.root, q.QuestionID, "The goal answered it.")
	if err != nil || w.State != StateWithdrawn || w.Withdrawal == nil || w.Withdrawal.Reason != "The goal answered it." || w.Reason != q.Reason || w.NextAction.Actor != ActorNone {
		t.Fatalf("withdraw: %+v %v", w, err)
	}
	if again, err := Withdraw(f.d, f.root, q.QuestionID, "The goal answered it."); err != nil || again.Withdrawal.WithdrawnAtUTC != w.Withdrawal.WithdrawnAtUTC {
		t.Fatalf("withdraw replay: %+v %v", again, err)
	}
	_, err = Withdraw(f.d, f.root, q.QuestionID, "Other reason.")
	wantCode(t, err, CodeReportConflict)
	s, err := Apply(f.d, f.root, p.Request)
	if err == nil || s.State != StatusNotRecorded || codeOfReasons(s) != CodeWithdrawn || s.Question == nil || s.Question.State != StateWithdrawn {
		t.Fatalf("answer after withdrawal: %+v %v", s, err)
	}
	_, err = PreviewAnswer(f.d, f.root, q.QuestionID, "request-withdraw-2", AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("short")})
	wantCode(t, err, CodeWithdrawn)
}

func TestHI05InterruptedAgentKeepsAnswerBinding(t *testing.T) {
	f := newFixture(t)
	q := f.ask(f.choice("interrupt/one"))
	r := f.answer(q.QuestionID, "request-interrupt-1", "short")
	// A new process (fresh dependencies, no in-memory state) resumes the work.
	fresh := Dependencies{Workspace: f.d.Workspace, Now: f.d.Now, Poll: time.Millisecond}
	again, err := AskDraft(fresh, f.root, f.choice("interrupt/one"))
	if err != nil || again.QuestionID != q.QuestionID || again.Answer == nil || again.Answer.RequestID != r.RequestID || again.State != StateAnswered {
		t.Fatalf("resumed agent lost the answer: %+v %v", again, err)
	}
	if _, err := Acknowledge(fresh, f.root, q.QuestionID, again.Answer.RequestID, again.Answer.RequestSHA256); err != nil {
		t.Fatal(err)
	}
}

// HI06: read surfaces create nothing; unsafe or invalid state is an
// understandable error that neither repairs nor damages preserved data.
func TestHI06ReadSurfacesCreateNoStateOrLocks(t *testing.T) {
	f := newFixture(t)
	before := snapshot(t, f.root)
	if l, err := List(f.d, f.root, ListQuery{ProjectID: "ply"}); err != nil || len(l.Questions) != 0 || !l.Complete {
		t.Fatalf("empty list: %+v %v", l, err)
	}
	id := "hiq_" + strings.Repeat("0", 32)
	_, err := Show(f.d, f.root, id)
	wantCode(t, err, CodeNotFound)
	if s, _ := Status(f.d, f.root, id, "request-none-1"); s.State != StatusUnavailable {
		t.Fatalf("status of an absent question: %+v", s)
	}
	_, err = Wait(f.d, f.root, id, time.Millisecond)
	wantCode(t, err, CodeNotFound)
	sameSnapshot(t, before, snapshot(t, f.root))

	q := f.ask(f.choice("reads/one"))
	f.answer(q.QuestionID, "request-reads-1", "short")
	before = snapshot(t, f.root)
	_, _ = List(f.d, f.root, ListQuery{ProjectID: "ply"})
	_, _ = Show(f.d, f.root, q.QuestionID)
	_, _ = Status(f.d, f.root, q.QuestionID, "request-reads-1")
	_, _ = Wait(f.d, f.root, q.QuestionID, time.Millisecond)
	_, _ = PreviewAnswer(f.d, f.root, q.QuestionID, "request-reads-2", AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("short")})
	sameSnapshot(t, before, snapshot(t, f.root))
}

func TestHI06UnsafeOrInvalidStateFailsWithoutRepair(t *testing.T) {
	setup := func(t *testing.T) (*hiFixture, Question) {
		f := newFixture(t)
		q := f.ask(f.choice("unsafe/one"))
		f.answer(q.QuestionID, "request-unsafe-1", "short")
		return f, q
	}
	reads := func(f *hiFixture, q Question) []error {
		_, e1 := List(f.d, f.root, ListQuery{ProjectID: "ply"})
		_, e2 := Show(f.d, f.root, q.QuestionID)
		s, _ := Status(f.d, f.root, q.QuestionID, "request-unsafe-1")
		var e3 error
		if s.State == StatusUnknown {
			e3 = errors.New("status is unknown")
		}
		return []error{e1, e2, e3}
	}
	cases := map[string]func(*testing.T, *hiFixture, Question){
		"symlinked store": func(t *testing.T, f *hiFixture, q Question) {
			moved := storeRoot(f.root) + ".real"
			mustDo(t, os.Rename(storeRoot(f.root), moved))
			mustDo(t, os.Symlink(moved, storeRoot(f.root)))
		},
		"question file is a directory": func(t *testing.T, f *hiFixture, q Question) {
			p := questionPath(f.root, q.QuestionID)
			mustDo(t, os.Rename(p, p+".moved"))
			mustDo(t, os.Mkdir(p, 0o700))
		},
		"corrupt answer": func(t *testing.T, f *hiFixture, q Question) {
			mustDo(t, os.WriteFile(answerPath(f.root, q.QuestionID), []byte(`{"truncated":`), 0o600))
		},
		"question digest drift": func(t *testing.T, f *hiFixture, q Question) {
			p := questionPath(f.root, q.QuestionID)
			raw, err := os.ReadFile(p)
			mustDo(t, err)
			mustDo(t, os.WriteFile(p, bytes.Replace(raw, []byte("Which output"), []byte("Which outpuT"), 1), 0o600))
		},
		"public permissions": func(t *testing.T, f *hiFixture, q Question) {
			mustDo(t, os.Chmod(answerPath(f.root, q.QuestionID), 0o644))
		},
		"answer for another question": func(t *testing.T, f *hiFixture, q Question) {
			other := f.ask(f.choice("unsafe/other"))
			f.answer(other.QuestionID, "request-unsafe-2", "full")
			raw, err := os.ReadFile(answerPath(f.root, other.QuestionID))
			mustDo(t, err)
			mustDo(t, os.WriteFile(answerPath(f.root, q.QuestionID), raw, 0o600))
		},
	}
	for name, damage := range cases {
		t.Run(name, func(t *testing.T) {
			f, q := setup(t)
			damage(t, f, q)
			before := snapshot(t, f.root)
			for i, err := range reads(f, q) {
				if err == nil {
					t.Fatalf("read %d accepted unsafe state", i)
				}
				if i < 2 && codeOf(err) != CodeStateInvalid {
					t.Fatalf("read %d: want %s, got %v", i, CodeStateInvalid, err)
				}
			}
			s, err := Apply(f.d, f.root, f.preview2(q.QuestionID))
			if err == nil || s.State == StatusRecorded {
				t.Fatalf("write accepted unsafe state: %+v", s)
			}
			sameSnapshot(t, before, snapshot(t, f.root))
		})
	}
}

func (f *hiFixture) preview2(id string) AnswerRequest {
	q, _ := readQuestionRecord(f.root, id)
	r := AnswerRequest{Kind: RequestKind, SchemaVersion: SchemaVersion, Mode: Mode, Purpose: PurposeAnswer, Workspace: f.root, QuestionID: id, QuestionSHA256: q.QuestionSHA256, Scope: q.Scope, RequestID: "request-unsafe-9", Answer: AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("full")}}
	r.BindingSHA256 = bindingOf(r)
	return r
}

func mustDo(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// HI03/HI07: apply takes the explicit workspace and refuses a request bound to
// another workspace before any state change.
func TestHI03ApplyRequiresTheExplicitMatchingWorkspace(t *testing.T) {
	f := newFixture(t)
	other := newFixture(t)
	q := f.ask(f.choice("workspace/one"))
	p := f.preview(q.QuestionID, "request-workspace-1", AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("short")})
	before, otherBefore := snapshot(t, f.root), snapshot(t, other.root)
	for _, root := range []string{other.root, "", "relative/root"} {
		s, err := Apply(f.d, root, p.Request)
		if err == nil || s.State == StatusRecorded {
			t.Fatalf("apply in %q accepted a request bound to %s: %+v", root, f.root, s)
		}
	}
	s, err := Apply(f.d, other.root, p.Request)
	if codeOf(err) != CodeRequestMismatch || s.State != StatusNotRecorded || s.Workspace != other.root {
		t.Fatalf("mismatched workspace: %+v %v", s, err)
	}
	sameSnapshot(t, before, snapshot(t, f.root))
	sameSnapshot(t, otherBefore, snapshot(t, other.root))
	if s, err := Apply(f.d, f.root, p.Request); err != nil || s.State != StatusRecorded {
		t.Fatalf("matching workspace: %+v %v", s, err)
	}
}
