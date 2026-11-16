package intelligence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// validated is a checked post: its decoded value, as stored, and its links.
type validated struct {
	value any
	id    string
	links []string
}

type link struct {
	Relation string `json:"relation"`
	ID       string `json:"id"`
}

type observationPost struct {
	Schema         string     `json:"schema"`
	ID             string     `json:"id"`
	RecordedAt     string     `json:"recorded_at"`
	Origin         Origin     `json:"origin"`
	Situation      string     `json:"situation"`
	Observation    string     `json:"observation"`
	Interpretation *string    `json:"interpretation"`
	Limitations    string     `json:"limitations"`
	Evidence       []evidence `json:"evidence"`
	PriorRefs      []string   `json:"prior_refs"`
	Effect         struct {
		Kind   string `json:"kind"`
		Detail string `json:"detail"`
	} `json:"effect"`
	Links []link `json:"links"`
}

type evidence struct {
	Basis  string `json:"basis"`
	Ref    string `json:"ref"`
	Detail string `json:"detail"`
}

type reviewPost struct {
	Schema     string `json:"schema"`
	ID         string `json:"id"`
	RecordedAt string `json:"recorded_at"`
	Origin     Origin `json:"origin"`
	Base       string `json:"base"`
	Head       string `json:"head"`
	Tool       string `json:"tool"`
	Command    string `json:"command"`
	Exit       int    `json:"exit"`
	Reviewer   struct {
		Provider string  `json:"provider"`
		Model    string  `json:"model"`
		Task     *string `json:"task"`
	} `json:"reviewer"`
	FirstReviewQuestion string   `json:"first_review_question"`
	Ratings             []rating `json:"ratings"`
	Signals             int      `json:"signals"`
	Used                int      `json:"used"`
	Irrelevant          int      `json:"irrelevant"`
	Wrong               int      `json:"wrong"`
	OtherObservations   *string  `json:"other_observations"`
	ProductSuggestion   *string  `json:"product_suggestion"`
	FollowUpCommit      *string  `json:"follow_up_commit"`
	Backlog             []string `json:"backlog"`
	Counting            string   `json:"counting"`
	Diff                string   `json:"diff"`
	DiffStderr          *string  `json:"diff_stderr"`
	Links               []link   `json:"links"`
}

type rating struct {
	File     string   `json:"file"`
	Line     int      `json:"line"`
	Function string   `json:"function"`
	Marks    []string `json:"marks"`
	Rating   string   `json:"rating"`
	Note     string   `json:"note"`
	Effect   string   `json:"effect"`
}

// validate checks raw against the schema it names and returns the post and kind.
func (s Store) validate(raw []byte) (validated, string, error) {
	var head struct {
		Schema string `json:"schema"`
	}
	if json.Unmarshal(raw, &head) != nil {
		return validated{}, "", invalidPost("the post is not a JSON object")
	}
	kind := ""
	for k, schema := range schemas {
		if schema == head.Schema {
			kind = k
		}
	}
	if kind == "" {
		return validated{}, "", invalidPost("schema must be one of ply.observation/v1, ply.code-health-review/v1")
	}
	var value any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if decoder.Decode(&value) != nil || decoder.More() {
		return validated{}, "", invalidPost("the post is not one JSON object")
	}
	var problem string
	var post validated
	switch kind {
	case KindObservation:
		var p observationPost
		if problem = strict(raw, &p); problem == "" {
			problem = s.checkObservation(p)
			post = validated{value: value, id: p.ID, links: linkIDs(p.Links)}
		}
	case KindCodeHealthReview:
		var p reviewPost
		if problem = strict(raw, &p); problem == "" {
			problem = s.checkReview(p)
			post = validated{value: value, id: p.ID, links: linkIDs(p.Links)}
		}
	}
	if problem != "" {
		return validated{}, "", invalidPost(problem)
	}
	return post, kind, nil
}

// strict decodes raw into p; every field must be present and known.
func strict(raw []byte, p any) string {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(p); err != nil {
		return "invalid post: " + strings.TrimPrefix(err.Error(), "json: ")
	}
	var fields map[string]json.RawMessage
	_ = json.Unmarshal(raw, &fields)
	expected, _ := json.Marshal(p)
	var want map[string]json.RawMessage
	_ = json.Unmarshal(expected, &want)
	for key := range want {
		if _, present := fields[key]; !present {
			return "missing field " + key
		}
	}
	return ""
}

func (s Store) checkCommon(id, recordedAt string, origin Origin, links []link) string {
	if !idPattern.MatchString(id) {
		return "id must be obs-<uuid>"
	}
	if _, err := time.Parse(timeLayout, recordedAt); err != nil {
		return "recorded_at must be a UTC time such as 2026-11-16T10:00:00Z"
	}
	switch {
	case origin.Project != s.project:
		return "origin.project must be " + s.project
	case !s.repos[origin.Repo]:
		return "origin.repo " + origin.Repo + " is not a repository of Project " + s.project
	case origin.Provider != "codex" && origin.Provider != "claude" && origin.Provider != "other":
		return "origin.provider must be codex, claude or other"
	case blank(origin.Role) || blank(origin.Activity):
		return "origin.role and origin.activity are required"
	case origin.Session != nil && blank(*origin.Session):
		return "origin.session is a known ID or null"
	}
	seen := map[link]bool{}
	for _, l := range links {
		switch l.Relation {
		case "corrects", "contradicts", "supports", "same_case":
			if !idPattern.MatchString(l.ID) {
				return "a " + l.Relation + " link points to obs-<uuid>"
			}
		case "uses":
			if !knowledgePattern.MatchString(l.ID) {
				return "a uses link points to a knowledge post k-<date>-<slug>"
			}
		default:
			return "link relation must be corrects, contradicts, supports, same_case or uses"
		}
		if l.ID == id || seen[l] {
			return "self and duplicate links are not allowed"
		}
		seen[l] = true
	}
	return ""
}

func (s Store) checkObservation(p observationPost) string {
	if problem := s.checkCommon(p.ID, p.RecordedAt, p.Origin, p.Links); problem != "" {
		return problem
	}
	if blank(p.Situation) || blank(p.Observation) || blank(p.Limitations) {
		return "situation, observation and limitations are required"
	}
	if p.Interpretation != nil && blank(*p.Interpretation) {
		return "interpretation is text or null"
	}
	if len(p.Evidence) == 0 {
		return "at least one evidence reference is required"
	}
	for _, e := range p.Evidence {
		if (e.Basis != "direct" && e.Basis != "reported" && e.Basis != "replay") || blank(e.Ref) || blank(e.Detail) {
			return "evidence needs basis direct, reported or replay, a ref and a detail"
		}
	}
	for _, ref := range p.PriorRefs {
		if blank(ref) {
			return "prior_refs entries are required text"
		}
	}
	switch p.Effect.Kind {
	case "new_question", "confirmation", "changed_action", "none":
		if blank(p.Effect.Detail) {
			return "effect.detail is required unless effect.kind is unknown"
		}
	case "unknown":
	default:
		return "effect.kind must be new_question, confirmation, changed_action, none or unknown"
	}
	return ""
}

func (s Store) checkReview(p reviewPost) string {
	if problem := s.checkCommon(p.ID, p.RecordedAt, p.Origin, p.Links); problem != "" {
		return problem
	}
	switch {
	case !shaPattern.MatchString(p.Base) || !shaPattern.MatchString(p.Head) || p.Base == p.Head:
		return "base and head are two different full commit SHAs"
	case blank(p.Tool) || blank(p.Command) || blank(p.Counting):
		return "tool, command and counting are required"
	case p.Exit != 0 && p.Exit != 2:
		return "exit is 0 or 2"
	case blank(p.Reviewer.Provider) || blank(p.Reviewer.Model):
		return "reviewer.provider and reviewer.model are required; model may be unknown"
	case p.Origin.Role != "author" && p.Origin.Role != "independent_reviewer":
		return "origin.role is author or independent_reviewer"
	case blank(p.FirstReviewQuestion):
		return "first_review_question is required"
	case p.FollowUpCommit != nil && !shaPattern.MatchString(*p.FollowUpCommit):
		return "follow_up_commit is a full commit SHA or null"
	}
	counts := map[string]int{}
	for _, r := range p.Ratings {
		if blank(r.File) || r.Line < 1 || blank(r.Function) || len(r.Marks) == 0 || blank(r.Note) {
			return "every rating needs file, line, function, marks and note"
		}
		switch r.Effect {
		case "new_question", "confirmation", "changed_action", "none":
		default:
			return "rating effect must be new_question, confirmation, changed_action or none"
		}
		counts[r.Rating]++
	}
	if counts["used"]+counts["irrelevant"]+counts["wrong"] != len(p.Ratings) {
		return "rating must be used, irrelevant or wrong"
	}
	if p.Signals != len(p.Ratings) || p.Used != counts["used"] || p.Irrelevant != counts["irrelevant"] || p.Wrong != counts["wrong"] {
		return fmt.Sprintf("counts must match ratings: %d signals, %d used, %d irrelevant, %d wrong",
			len(p.Ratings), counts["used"], counts["irrelevant"], counts["wrong"])
	}
	return ""
}

func linkIDs(links []link) []string {
	ids := make([]string, 0, len(links))
	for _, l := range links {
		ids = append(ids, l.ID)
	}
	return ids
}

func blank(value string) bool { return strings.TrimSpace(value) == "" }
