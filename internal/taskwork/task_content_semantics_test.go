package taskwork

import (
	"errors"
	"fmt"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

// These tests pin validateTaskContentSemantics: which rule answers when
// several are broken at once, and the text of its answer. Every case is a
// complete draft or manifest built in memory. Content that the schema accepts
// is also decoded with decodeTaskContent, which classes every answer as
// task_content_invalid_input; content that the schema rejects is given to the
// validator alone.

// members are the members of a content object.
type members = map[string]canonicaljson.Value

func list(values ...canonicaljson.Value) []canonicaljson.Value {
	return append([]canonicaljson.Value{}, values...)
}

var semanticsDigest = "sha256:" + strings.Repeat("a", 64)

const semanticsOID = "0123456789abcdef0123456789abcdef01234567"

// semanticsShapes names, per operation, the kind of its draft and of its
// manifest, and the members that only one of them has.
var semanticsShapes = map[string]struct {
	draftKind, manifestKind string
	draft, manifest         members
}{
	"problem_record": {"WorkspaceTaskProblemDraft@1", "WorkspaceTaskProblemRevision@1", members{"expected_previous": nil}, members{"previous": nil, "revision": int64(1)}},
	"spec_record":    {"WorkspaceTaskSpecDraft@1", "WorkspaceTaskSpecRevision@1", members{"expected_previous": nil}, members{"previous": nil, "revision": int64(1)}},
	"spec_assess":    {"WorkspaceTaskSpecAssessmentDraft@1", "WorkspaceTaskSpecAssessment@1", members{"expected_previous_assessment": nil}, members{"id": "asm_" + strings.Repeat("a", 32), "ordinal": int64(1), "previous": nil}},
	"spec_select":    {"WorkspaceTaskSolutionSelectionDraft@1", "WorkspaceTaskSolutionSelection@1", members{"expected_previous_selection": nil}, members{"id": "sel_" + strings.Repeat("a", 32), "ordinal": int64(1), "previous": nil}},
	"spec_withdraw":  {"WorkspaceTaskSolutionSelectionDraft@1", "WorkspaceTaskSolutionSelection@1", members{"expected_previous_selection": nil}, members{"id": "sel_" + strings.Repeat("a", 32), "ordinal": int64(1), "previous": nil}},
}

// semanticsCase is one draft or manifest and the variant it is validated as.
// A test changes the fields and asks for the answer.
type semanticsCase struct {
	t         *testing.T
	operation string
	published bool
	internal  bool
	fields    members
}

// newSemanticsCase completes the fields of an operation into a draft or, when
// published, a manifest.
func newSemanticsCase(t *testing.T, operation string, published bool, fields members) *semanticsCase {
	t.Helper()
	shape := semanticsShapes[operation]
	kind, own := shape.draftKind, shape.draft
	if published {
		kind, own = shape.manifestKind, shape.manifest
		fields["source_draft_sha256"], fields["recorded_at_utc"] = semanticsDigest, "2026-11-20T12:00:00Z"
	} else {
		fields["registry_upgrade"] = nil
	}
	for name, value := range own {
		fields[name] = value
	}
	fields["publication_key"], fields["task_id"] = "semantics/case", "task"
	fields["recorder"] = contentObject(members{"actor_claim": "fixture agent", "control_surface": "isolated test fixture", "recorded_at_utc": "2026-11-20T12:00:00Z"})
	return &semanticsCase{t: t, operation: operation, published: published, fields: contentFields(contentEnvelope(kind, fields))}
}

// declare makes the given IDs the documents of the case: files in a draft,
// descriptors in a manifest.
func (c *semanticsCase) declare(ids ...string) {
	documents := list()
	for _, id := range ids {
		locator := "/fixture/" + id + ".md"
		if c.published {
			documents = append(documents, contentObject(members{"id": id, "locator": locator, "sha256": semanticsDigest, "size_bytes": int64(1), "media_type": "text/markdown", "provenance": contentObject(members{"origin_locator": nil, "git": nil})}))
		} else {
			documents = append(documents, contentObject(members{"id": id, "source": contentObject(members{"kind": "file", "locator": locator, "sha256": semanticsDigest, "size_bytes": int64(1), "media_type": "text/markdown", "git_provenance": nil})}))
		}
	}
	c.fields["documents"] = documents
}

// declareAlso adds one document after the declared ones.
func (c *semanticsCase) declareAlso(id string) {
	ids := []string{}
	for _, document := range contentArray(c.fields, "documents") {
		ids = append(ids, contentIDKey(document))
	}
	c.declare(append(ids, id)...)
}

const semanticsAccepted = "accepted"

// validator names the answer of the validator alone: "accepted" or the text
// of its error.
func (c *semanticsCase) validator() string {
	if err := validateTaskContentSemantics(c.fields, c.operation, c.published, c.internal); err != nil {
		return err.Error()
	}
	return semanticsAccepted
}

// decoded names the answer of decoding the case as its operation: "accepted"
// or the class and detail of the error.
func (c *semanticsCase) decoded() string {
	c.t.Helper()
	raw, err := canonicaljson.Marshal(contentObject(c.fields))
	if err != nil {
		c.t.Fatal(err)
	}
	if _, err = decodeTaskContent(raw, c.operation, c.published, c.internal); err != nil {
		var classified *TaskContentError
		if !errors.As(err, &classified) {
			c.t.Fatalf("unclassified answer: %v", err)
		}
		return err.Error()
	}
	return semanticsAccepted
}

// answer is the validator's answer to content that the schema accepts. The
// case fails unless decoding gives the same answer, classed as invalid input.
func (c *semanticsCase) answer() string {
	c.t.Helper()
	answer, want := c.validator(), semanticsAccepted
	if answer != semanticsAccepted {
		want = "task_content_invalid_input: " + answer
	}
	if got := c.decoded(); got != want {
		c.t.Fatalf("decoded as %q, but the validator answers %q", got, answer)
	}
	return answer
}

// changeEntry changes object i of the list member key.
func changeEntry(fields members, key string, i int, change func(entry members)) {
	entries := slices.Clone(contentArray(fields, key))
	entry := contentFields(entries[i])
	change(entry)
	entries[i] = contentObject(entry)
	fields[key] = entries
}

// set returns a change that sets one member.
func set(name string, value canonicaljson.Value) func(members) {
	return func(object members) { object[name] = value }
}

func partRef(document string, section canonicaljson.Value) canonicaljson.Value {
	return contentObject(members{"document_id": document, "section": section})
}

func specPart(state string, reason canonicaljson.Value, refs ...canonicaljson.Value) canonicaljson.Value {
	return contentObject(members{"state": state, "reason": reason, "documents": list(refs...)})
}

func problemSource(id string, document, withholding canonicaljson.Value) canonicaljson.Value {
	return contentObject(members{"id": id, "origin_kind": "email", "origin_claim": "Fixture origin.", "received_at_utc": nil, "document_id": document, "withholding_reason": withholding})
}

func problemClaim(id, source, verification string, reason canonicaljson.Value, evidence ...canonicaljson.Value) canonicaljson.Value {
	return contentObject(members{"id": id, "source_id": source, "text": "Fixture claim.", "verification": verification, "reason": reason, "evidence_document_ids": list(evidence...)})
}

func specRequirement(id string, functional, technical, verification []canonicaljson.Value) canonicaljson.Value {
	return contentObject(members{"id": id, "functional_refs": functional, "acceptance": "Fixture behaviour is preserved.", "verification_ids": verification, "technical_refs": technical})
}

func specPhase(id string, requirements, verification []canonicaljson.Value) canonicaljson.Value {
	return contentObject(members{"id": id, "purpose": "Fixture phase.", "requirement_ids": requirements, "entry_criteria": list(), "exit_criteria": list("Fixture verified."), "verification_ids": verification})
}

func assessmentCheck(id, outcome string, evidence ...canonicaljson.Value) canonicaljson.Value {
	return contentObject(members{"id": id, "outcome": outcome, "reason": "Fixture check.", "evidence_document_ids": list(evidence...)})
}

// problemCase is an authored problem with the documents "evidence", "letter"
// and "problem": the source "letter" holds a document and the source "rumour"
// is withheld, the claim "confirmed" has evidence and the claim "unverified"
// has none, and the deadline names the confirmed claim.
func problemCase(t *testing.T, published bool) *semanticsCase {
	t.Helper()
	c := newSemanticsCase(t, "problem_record", published, members{
		"origin": contentObject(members{"kind": "authored"}), "title": "Fixture problem", "summary": "Explain the fixture need.", "change_reason": "Fixture problem.",
		"problem_document_id": "problem",
		"sources":             list(problemSource("letter", "letter", nil), problemSource("rumour", nil, "Fixture: withheld.")),
		"claims":              list(problemClaim("confirmed", "letter", "confirmed", "Fixture check.", "evidence"), problemClaim("unverified", "rumour", "unverified", nil)),
		"deadline":            contentObject(members{"date": "2026-12-24", "claim_id": "confirmed"}),
	})
	c.declare("evidence", "letter", "problem")
	return c
}

// generatedProblemCase is the draft of a problem that Ply writes from a Task
// summary: it has no documents, sources, claims or deadline.
func generatedProblemCase(t *testing.T, origin string) *semanticsCase {
	t.Helper()
	c := problemCase(t, false)
	c.fields["origin"] = contentObject(members{"kind": "legacy_summary_import", "legacy_summary_sha256": semanticsDigest})
	if origin == "task_create_summary" {
		c.fields["origin"], c.internal = contentObject(members{"kind": origin}), true
	}
	c.fields["sources"], c.fields["claims"], c.fields["deadline"] = list(), list(), nil
	c.declare()
	return c
}

// specCase is a Spec with three present parts over the documents "abstract",
// "functional" and "design", and the supporting document "notes". Requirement
// r-01 is verified by "unit" and r-02 by "e2e"; the phase "build" names both.
func specCase(t *testing.T, published bool) *semanticsCase {
	t.Helper()
	c := bareSpecCase(t, published)
	c.fields["parts"] = contentObject(members{
		"abstract":   specPart("present", nil, partRef("abstract", nil)),
		"functional": specPart("present", nil, partRef("functional", nil), partRef("functional", "scope")),
		"technical":  specPart("present", nil, partRef("design", nil)),
	})
	c.fields["supporting"] = list(contentObject(members{"document_id": "notes", "usage": "informative"}))
	c.fields["requirements"] = list(
		specRequirement("r-01", list(partRef("functional", "scope")), list(partRef("design", nil)), list("unit")),
		specRequirement("r-02", list(partRef("functional", nil)), list(partRef("design", nil)), list("e2e")),
	)
	c.fields["phases"] = list(specPhase("build", list("r-01", "r-02"), list("e2e", "unit")))
	c.declare("abstract", "design", "functional", "notes")
	return c
}

// bareSpecCase is a Spec with nothing resolved: no documents, requirements or
// phases, and three unresolved parts.
func bareSpecCase(t *testing.T, published bool) *semanticsCase {
	t.Helper()
	unresolved := specPart("unresolved", "Fixture: not yet written.")
	c := newSemanticsCase(t, "spec_record", published, members{
		"spec_id": "solution", "problem": contentRefValue(&TaskRevisionRef{1, semanticsDigest}), "title": "Fixture solution", "change_reason": "Fixture solution.",
		"parts":      contentObject(members{"abstract": unresolved, "functional": unresolved, "technical": unresolved}),
		"supporting": list(), "requirements": list(), "removed_requirement_ids": list(), "phases": list(), "dependencies": list(),
		"implementation_basis": contentObject(members{
			"project_id": "project", "repo_id": "repo", "git_common_dir": "/fixture/repo/.git", "epic_id": "epic",
			"parent_worktree_id": "wt_" + strings.Repeat("a", 32), "parent_ref": "refs/heads/epic",
			"parent_oid": semanticsOID, "parent_tree": semanticsOID, "start_oid": semanticsOID, "start_tree": semanticsOID,
		}),
	})
	c.declare()
	return c
}

var assessmentCheckOrder = []string{"acceptance_coverage", "implementation_basis", "problem_coverage", "recovery", "scope_and_phases", "three_parts"}

// assessmentCase is a ready assessment with the six checks passing, the first
// on the evidence of the document "evidence".
func assessmentCase(t *testing.T, published bool) *semanticsCase {
	t.Helper()
	checks := list(assessmentCheck(assessmentCheckOrder[0], "pass", "evidence"))
	for _, id := range assessmentCheckOrder[1:] {
		checks = append(checks, assessmentCheck(id, "pass"))
	}
	c := newSemanticsCase(t, "spec_assess", published, members{
		"spec_id": "solution", "spec": contentRefValue(&TaskRevisionRef{1, semanticsDigest}), "outcome": "ready", "reason": "Fixture readiness claim.",
		"open_questions": list(), "checks": checks,
	})
	c.declare("evidence")
	return c
}

// selectionCase is the selection of a solution or, as "spec_withdraw", the
// withdrawal of the selected one.
func selectionCase(t *testing.T, operation string, published bool) *semanticsCase {
	t.Helper()
	action, solution := "withdraw", canonicaljson.Value(nil)
	if operation == "spec_select" {
		action = "select"
		solution = contentObject(members{"spec_id": "solution", "spec": contentRefValue(&TaskRevisionRef{1, semanticsDigest}), "problem": contentRefValue(&TaskRevisionRef{1, semanticsDigest}), "assessment": contentRefValue(&TaskDecisionRef{"asm_" + strings.Repeat("a", 32), semanticsDigest})})
	}
	return newSemanticsCase(t, operation, published, members{
		"action": action, "reason": "Fixture selection.", "solution": solution,
		"human_decision": contentObject(members{"actor_claim": "fixture human, not actual approval", "decided_at_utc": "2026-11-20T12:00:00Z", "source": "explicit_human_instruction", "statement": "Fixture only: decide for the isolated test."}),
	})
}

// semanticsRule breaks one rule of a case and gives the validator's answer.
type semanticsRule struct {
	name   string
	breach func(c *semanticsCase)
	answer string
}

func runSemanticsRules(t *testing.T, base func(t *testing.T) *semanticsCase, rules []semanticsRule) {
	t.Helper()
	for _, rule := range rules {
		t.Run(rule.name, func(t *testing.T) {
			c := base(t)
			rule.breach(c)
			if got := c.answer(); got != rule.answer {
				t.Fatalf("got %q, want %q", got, rule.answer)
			}
		})
	}
}

// semanticsOrder lists rules of one case that can be broken together, in the
// order the validator answers them.
type semanticsOrder struct {
	name  string
	base  func(t *testing.T, published bool) *semanticsCase
	rules []semanticsRule
}

func unusedDocumentRule() semanticsRule {
	return semanticsRule{name: "document that nothing refers to", breach: func(c *semanticsCase) { c.declareAlso("spare") }, answer: "unused document spare"}
}

const (
	withheldAnswer      = "withheld source cannot have a document"
	unknownSourceAnswer = "claim source does not exist"
	verifiedAnswer      = "verified claim requires reason and evidence"
	unverifiedAnswer    = "unverified claim cannot have verification evidence"
	generatedDraft      = "generated problem requires empty documents, sources and claims, null deadline and problem ID"
	outsidePartAnswer   = "requirement reference is outside its part"
	unknownRequirement  = "phase refers to an unknown requirement"
	sixChecksAnswer     = "assessment must contain exactly six checks"
	failingCheckAnswer  = "ready assessment requires all checks to pass"
)

func problemRuleOrder() semanticsOrder {
	source := func(i int, change func(members)) func(*semanticsCase) {
		return func(c *semanticsCase) { changeEntry(c.fields, "sources", i, change) }
	}
	claim := func(i int, change func(members)) func(*semanticsCase) {
		return func(c *semanticsCase) { changeEntry(c.fields, "claims", i, change) }
	}
	return semanticsOrder{name: "problem", base: problemCase, rules: []semanticsRule{
		{name: "problem document that is not declared", breach: func(c *semanticsCase) { c.fields["problem_document_id"] = "absent" }, answer: "document absent is missing"},
		{name: "first source withheld with a document", breach: source(0, set("withholding_reason", "Fixture: withheld.")), answer: withheldAnswer},
		{name: "first source's document not declared", breach: source(0, set("document_id", "gone")), answer: "document gone is missing"},
		{name: "second source withheld with a document", breach: source(1, set("document_id", "evidence")), answer: withheldAnswer},
		{name: "first claim of an unknown source", breach: claim(0, set("source_id", "nobody")), answer: unknownSourceAnswer},
		{name: "first claim confirmed without a reason", breach: claim(0, set("reason", nil)), answer: verifiedAnswer},
		{name: "first claim's evidence not declared", breach: claim(0, set("evidence_document_ids", list("lost"))), answer: "document lost is missing"},
		{name: "second claim of an unknown source", breach: claim(1, set("source_id", "nobody")), answer: unknownSourceAnswer},
		{name: "second claim unverified with a reason", breach: claim(1, set("reason", "Fixture reason.")), answer: unverifiedAnswer},
		{name: "deadline of an unknown claim", breach: func(c *semanticsCase) { changeObject(c.fields, "deadline", set("claim_id", "nothing")) }, answer: "deadline claim does not exist"},
		unusedDocumentRule(),
	}}
}

// specRuleOrders are the rules of a Spec. One requirement meets its rules in
// the order functional refs, technical refs and verification IDs; the first
// order breaks them as missing, the two after it as refs outside the part.
func specRuleOrders() []semanticsOrder {
	part := func(name string, change func(members)) func(*semanticsCase) {
		return func(c *semanticsCase) {
			changeObject(c.fields, "parts", func(parts members) { changeObject(parts, name, change) })
		}
	}
	requirement := func(i int, change func(members)) func(*semanticsCase) {
		return func(c *semanticsCase) { changeEntry(c.fields, "requirements", i, change) }
	}
	phase := func(change func(members)) func(*semanticsCase) {
		return func(c *semanticsCase) { changeEntry(c.fields, "phases", 0, change) }
	}
	partRules := func(name string) []semanticsRule {
		return []semanticsRule{
			{name: name + " part present with a reason", breach: part(name, set("reason", "Fixture reason.")), answer: fmt.Sprintf("present %s needs documents and null reason", name)},
			{name: name + " part's document not declared", breach: part(name, set("documents", list(partRef("gone-"+name, nil)))), answer: fmt.Sprintf("document gone-%s is missing", name)},
		}
	}
	withoutVerification := semanticsRule{name: "first requirement without verification IDs", breach: requirement(0, set("verification_ids", list())), answer: "requirement needs a verification ID"}
	rules := slices.Concat(partRules("abstract"), partRules("functional"), partRules("technical"), []semanticsRule{
		{name: "supporting document that a part holds", breach: func(c *semanticsCase) {
			c.fields["supporting"] = append(list(contentObject(members{"document_id": "design", "usage": "normative"})), contentArray(c.fields, "supporting")...)
		}, answer: "supporting document duplicates a part"},
		{name: "supporting document not declared", breach: func(c *semanticsCase) {
			changeEntry(c.fields, "supporting", len(contentArray(c.fields, "supporting"))-1, set("document_id", "zz-gone"))
		}, answer: "document zz-gone is missing"},
		{name: "first requirement without functional refs", breach: requirement(0, set("functional_refs", list())), answer: "requirement needs functional refs"},
		{name: "first requirement without technical refs", breach: requirement(0, set("technical_refs", list())), answer: "requirement needs technical refs"},
		withoutVerification,
		{name: "second requirement's functional ref outside the part", breach: requirement(1, set("functional_refs", list(partRef("design", nil)))), answer: outsidePartAnswer},
		{name: "phase without the verifier of its first requirement", breach: phase(set("verification_ids", list("e2e"))), answer: "phase does not cover requirement verifier unit"},
		{name: "phase's second requirement unknown", breach: phase(set("requirement_ids", list("r-01", "r-09"))), answer: unknownRequirement},
		{name: "later phase without the verifier of its requirement", breach: func(c *semanticsCase) {
			c.fields["phases"] = append(contentArray(c.fields, "phases"), specPhase("later", list("r-02"), list()))
		}, answer: "phase does not cover requirement verifier e2e"},
		unusedDocumentRule(),
	})
	return []semanticsOrder{
		{name: "spec", base: specCase, rules: rules},
		{name: "spec/functional ref outside the part", base: specCase, rules: []semanticsRule{
			{name: "first requirement's functional ref outside the part", breach: requirement(0, set("functional_refs", list(partRef("design", nil)))), answer: outsidePartAnswer},
			{name: "first requirement without technical refs", breach: requirement(0, set("technical_refs", list())), answer: "requirement needs technical refs"},
			withoutVerification,
		}},
		{name: "spec/technical ref outside the part", base: specCase, rules: []semanticsRule{
			{name: "first requirement without functional refs", breach: requirement(0, set("functional_refs", list())), answer: "requirement needs functional refs"},
			{name: "first requirement's technical ref outside the part", breach: requirement(0, set("technical_refs", list(partRef("functional", nil)))), answer: outsidePartAnswer},
			withoutVerification,
		}},
		{name: "spec/phase's first requirement unknown", base: specCase, rules: []semanticsRule{
			{name: "phase's first requirement unknown", breach: phase(set("requirement_ids", list("r-00", "r-02"))), answer: unknownRequirement},
			{name: "phase without the verifier of its second requirement", breach: phase(set("verification_ids", list("unit"))), answer: "phase does not cover requirement verifier e2e"},
		}},
	}
}

func assessmentRuleOrder() semanticsOrder {
	check := func(i int, change func(members)) func(*semanticsCase) {
		return func(c *semanticsCase) { changeEntry(c.fields, "checks", i, change) }
	}
	return semanticsOrder{name: "assessment", base: assessmentCase, rules: []semanticsRule{
		{name: "five checks", breach: func(c *semanticsCase) { c.fields["checks"] = contentArray(c.fields, "checks")[:5] }, answer: sixChecksAnswer},
		{name: "first check under another ID", breach: check(0, set("id", "aaa")), answer: "assessment check must be acceptance_coverage"},
		{name: "first check failing", breach: check(0, set("outcome", "fail")), answer: failingCheckAnswer},
		{name: "first check's evidence not declared", breach: check(0, set("evidence_document_ids", list("gone"))), answer: "document gone is missing"},
		{name: "second check under another ID", breach: check(1, set("id", "implementation")), answer: "assessment check must be implementation_basis"},
		{name: "second check unknown", breach: check(1, set("outcome", "unknown")), answer: failingCheckAnswer},
		{name: "open question", breach: func(c *semanticsCase) { c.fields["open_questions"] = list("Fixture question?") }, answer: "ready assessment cannot have open questions"},
		unusedDocumentRule(),
	}}
}

func TestTaskContentSemanticsAcceptsEveryOperation(t *testing.T) {
	t.Parallel()
	for _, published := range []bool{false, true} {
		for name, c := range map[string]*semanticsCase{
			"problem": problemCase(t, published), "spec": specCase(t, published), "bare spec": bareSpecCase(t, published), "assessment": assessmentCase(t, published),
			"selection": selectionCase(t, "spec_select", published), "withdrawal": selectionCase(t, "spec_withdraw", published),
		} {
			if got := c.answer(); got != semanticsAccepted {
				t.Errorf("%s, published %t: got %q", name, published, got)
			}
		}
	}
}

// Each rule answers when it alone is broken, and before every rule after it
// when both are broken. Drafts and manifests meet the same rules.
func TestTaskContentSemanticsEarlierRuleAnswersFirst(t *testing.T) {
	t.Parallel()
	for _, order := range append(specRuleOrders(), problemRuleOrder(), assessmentRuleOrder()) {
		for _, published := range []bool{false, true} {
			for i, earlier := range order.rules {
				t.Run(fmt.Sprintf("%s/published %t/%s", order.name, published, earlier.name), func(t *testing.T) {
					c := order.base(t, published)
					earlier.breach(c)
					if got := c.answer(); got != earlier.answer {
						t.Fatalf("alone: got %q, want %q", got, earlier.answer)
					}
					for _, later := range order.rules[i+1:] {
						c = order.base(t, published)
						earlier.breach(c)
						later.breach(c)
						if got := c.answer(); got != earlier.answer {
							t.Errorf("with %s: got %q, want %q", later.name, got, earlier.answer)
						}
					}
				})
			}
		}
	}
}

// The draft of a generated problem is empty: Ply adds its problem document
// when it publishes. Only that one rule answers for such a draft.
func TestGeneratedProblemDraftHasNoContent(t *testing.T) {
	t.Parallel()
	for _, origin := range []string{"task_create_summary", "legacy_summary_import"} {
		base := func(t *testing.T) *semanticsCase { return generatedProblemCase(t, origin) }
		t.Run(origin, func(t *testing.T) {
			if got := base(t).answer(); got != semanticsAccepted {
				t.Fatalf("got %q", got)
			}
			runSemanticsRules(t, base, []semanticsRule{
				{name: "declared problem document", breach: func(c *semanticsCase) { c.declare("problem") }, answer: generatedDraft},
				{name: "source", breach: func(c *semanticsCase) { c.fields["sources"] = list(problemSource("letter", nil, nil)) }, answer: generatedDraft},
				{name: "claim of an unknown source", breach: func(c *semanticsCase) { c.fields["claims"] = list(problemClaim("claim", "nobody", "unverified", nil)) }, answer: generatedDraft},
				{name: "deadline of an unknown claim", breach: func(c *semanticsCase) {
					c.fields["deadline"] = contentObject(members{"date": "2026-12-24", "claim_id": "nothing"})
				}, answer: generatedDraft},
				{name: "problem document under another ID", breach: func(c *semanticsCase) { c.fields["problem_document_id"] = "other" }, answer: generatedDraft},
			})
		})
	}
}

// A generated problem's manifest holds the problem document that Ply added
// and meets the rules of an authored problem.
func TestGeneratedProblemManifestDeclaresItsDocument(t *testing.T) {
	t.Parallel()
	for _, origin := range []canonicaljson.Value{
		contentObject(members{"kind": "task_create_summary"}),
		contentObject(members{"kind": "legacy_summary_import", "legacy_summary_sha256": semanticsDigest}),
	} {
		base := func(t *testing.T) *semanticsCase {
			c := problemCase(t, true)
			c.fields["origin"], c.internal = origin, true
			return c
		}
		t.Run(contentString(contentFields(origin), "kind"), func(t *testing.T) {
			if got := base(t).answer(); got != semanticsAccepted {
				t.Fatalf("with sources, claims and a deadline: got %q", got)
			}
			runSemanticsRules(t, base, []semanticsRule{
				{name: "problem document that is not declared", breach: func(c *semanticsCase) {
					c.fields["sources"], c.fields["claims"], c.fields["deadline"] = list(), list(), nil
					c.declare()
				}, answer: "document problem is missing"},
				{name: "claim of an unknown source", breach: func(c *semanticsCase) { changeEntry(c.fields, "claims", 0, set("source_id", "nobody")) }, answer: unknownSourceAnswer},
			})
		})
	}
}

// Only the schema knows that a Task's creation is internal: it rejects the
// summary origin elsewhere, and the validator treats it as generated anyway.
func TestTaskCreateSummaryOriginIsInternalToTheSchemaOnly(t *testing.T) {
	t.Parallel()
	c := generatedProblemCase(t, "task_create_summary")
	c.internal = false
	if got, want := c.decoded(), "task_content_invalid_input: origin: invalid problem origin"; got != want {
		t.Fatalf("decoded: got %q, want %q", got, want)
	}
	if got := c.validator(); got != semanticsAccepted {
		t.Fatalf("validator: got %q", got)
	}
	c.declare("problem")
	if got := c.validator(); got != generatedDraft {
		t.Fatalf("validator with a document: got %q, want %q", got, generatedDraft)
	}
}

func TestProblemSourceHoldsADocumentOrIsWithheld(t *testing.T) {
	t.Parallel()
	base := func(t *testing.T) *semanticsCase {
		c := problemCase(t, false)
		c.fields["claims"], c.fields["deadline"] = list(), nil
		c.declare("letter", "problem")
		return c
	}
	sources := func(sources ...canonicaljson.Value) func(*semanticsCase) {
		return func(c *semanticsCase) { c.fields["sources"] = list(sources...) }
	}
	runSemanticsRules(t, base, []semanticsRule{
		{name: "document", breach: sources(problemSource("letter", "letter", nil)), answer: semanticsAccepted},
		{name: "withheld", breach: sources(problemSource("letter", nil, "Fixture: withheld."), problemSource("other", "letter", nil)), answer: semanticsAccepted},
		{name: "neither document nor withholding reason", breach: sources(problemSource("bare", nil, nil), problemSource("letter", "letter", nil)), answer: semanticsAccepted},
		{name: "two sources of one document", breach: sources(problemSource("copy", "letter", nil), problemSource("letter", "letter", nil)), answer: semanticsAccepted},
		{name: "the problem document", breach: sources(problemSource("letter", "letter", nil), problemSource("self", "problem", nil)), answer: semanticsAccepted},
		{name: "withheld with a document", breach: sources(problemSource("letter", "letter", "Fixture: withheld.")), answer: withheldAnswer},
		{name: "no sources", breach: sources(), answer: "unused document letter"},
	})
}

// An unverified claim has neither reason nor evidence; a confirmed or refuted
// claim has both.
func TestProblemClaimEvidenceFollowsItsVerification(t *testing.T) {
	t.Parallel()
	for _, verification := range []string{"unverified", "confirmed", "refuted"} {
		for _, reasoned := range []bool{false, true} {
			for _, evidenced := range []bool{false, true} {
				want := semanticsAccepted
				switch {
				case verification == "unverified" && (reasoned || evidenced):
					want = unverifiedAnswer
				case verification != "unverified" && (!reasoned || !evidenced):
					want = verifiedAnswer
				}
				c := problemCase(t, false)
				c.fields["deadline"] = nil
				reason, evidence, documents := canonicaljson.Value(nil), list(), []string{"letter", "problem"}
				if reasoned {
					reason = "Fixture reason."
				}
				if evidenced {
					evidence, documents = list("evidence"), []string{"evidence", "letter", "problem"}
				}
				c.fields["claims"] = list(problemClaim("claim", "letter", verification, reason, evidence...))
				c.declare(documents...)
				if got := c.answer(); got != want {
					t.Errorf("%s, reason %t, evidence %t: got %q, want %q", verification, reasoned, evidenced, got, want)
				}
			}
		}
	}
}

func TestProblemClaimsAndDeadline(t *testing.T) {
	t.Parallel()
	base := func(t *testing.T) *semanticsCase { return problemCase(t, false) }
	runSemanticsRules(t, base, []semanticsRule{
		{name: "no deadline", breach: func(c *semanticsCase) { c.fields["deadline"] = nil }, answer: semanticsAccepted},
		{name: "deadline of the last claim", breach: func(c *semanticsCase) { changeObject(c.fields, "deadline", set("claim_id", "unverified")) }, answer: semanticsAccepted},
		{name: "deadline of a source", breach: func(c *semanticsCase) { changeObject(c.fields, "deadline", set("claim_id", "letter")) }, answer: "deadline claim does not exist"},
		{name: "deadline without claims", breach: func(c *semanticsCase) {
			c.fields["claims"] = list()
			c.declare("letter", "problem")
		}, answer: "deadline claim does not exist"},
		{name: "claim of a claim", breach: func(c *semanticsCase) { changeEntry(c.fields, "claims", 1, set("source_id", "confirmed")) }, answer: unknownSourceAnswer},
		{name: "claim without sources", breach: func(c *semanticsCase) {
			c.fields["sources"] = list()
			c.declare("evidence", "problem")
		}, answer: unknownSourceAnswer},
		{name: "evidence that a source and the problem also hold", breach: func(c *semanticsCase) {
			changeEntry(c.fields, "claims", 0, set("evidence_document_ids", list("evidence", "letter", "problem")))
		}, answer: semanticsAccepted},
		{name: "second evidence document not declared", breach: func(c *semanticsCase) {
			changeEntry(c.fields, "claims", 0, set("evidence_document_ids", list("evidence", "lost")))
		}, answer: "document lost is missing"},
	})
}

// A present part has documents and no reason. An unresolved part has a reason
// and no documents, and so has a part that is not applicable, which only the
// technical part may be.
func TestSpecPartState(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"abstract", "functional", "technical"} {
		invalid, incomplete := fmt.Sprintf("invalid %s part state", name), fmt.Sprintf("present %s needs documents and null reason", name)
		notApplicable := invalid
		if name == "technical" {
			notApplicable = semanticsAccepted
		}
		for _, tc := range []struct {
			state                string
			reasoned, documented bool
			want                 string
		}{
			{"present", false, true, semanticsAccepted},
			{"present", true, true, incomplete},
			{"present", false, false, incomplete},
			{"present", true, false, incomplete},
			{"unresolved", true, false, semanticsAccepted},
			{"unresolved", false, false, invalid},
			{"unresolved", true, true, invalid},
			{"unresolved", false, true, invalid},
			{"not_applicable", true, false, notApplicable},
			{"not_applicable", false, false, invalid},
			{"not_applicable", true, true, invalid},
			{"not_applicable", false, true, invalid},
		} {
			c := bareSpecCase(t, false)
			reason, refs := canonicaljson.Value(nil), list()
			if tc.reasoned {
				reason = "Fixture reason."
			}
			if tc.documented {
				refs = list(partRef("part", nil))
				c.declare("part")
			}
			changeObject(c.fields, "parts", set(name, specPart(tc.state, reason, refs...)))
			if got := c.answer(); got != tc.want {
				t.Errorf("%s %s, reason %t, documents %t: got %q, want %q", name, tc.state, tc.reasoned, tc.documented, got, tc.want)
			}
		}
	}
}

func TestSpecPartAndSupportingDocuments(t *testing.T) {
	t.Parallel()
	base := func(t *testing.T) *semanticsCase { return specCase(t, false) }
	supporting := func(ids ...string) func(*semanticsCase) {
		return func(c *semanticsCase) {
			entries := list()
			for _, id := range ids {
				entries = append(entries, contentObject(members{"document_id": id, "usage": "normative"}))
			}
			c.fields["supporting"] = entries
		}
	}
	runSemanticsRules(t, base, []semanticsRule{
		{name: "parts sharing a document", breach: func(c *semanticsCase) {
			changeObject(c.fields, "parts", set("abstract", specPart("present", nil, partRef("design", nil))))
			c.declare("design", "functional", "notes")
		}, answer: semanticsAccepted},
		{name: "part's second document not declared", breach: func(c *semanticsCase) {
			changeObject(c.fields, "parts", set("abstract", specPart("present", nil, partRef("abstract", nil), partRef("gone", nil))))
		}, answer: "document gone is missing"},
		{name: "no supporting documents", breach: supporting(), answer: "unused document notes"},
		{name: "supporting document of the abstract part", breach: supporting("abstract"), answer: "supporting document duplicates a part"},
		{name: "supporting document of a part after one of its own", breach: supporting("functional", "notes"), answer: "supporting document duplicates a part"},
	})
}

// A requirement refers into the functional part and, unless that part is not
// applicable, into the technical part; a reference is the pair of a document
// and a section that the part itself holds.
func TestSpecRequirementRefsStayInsideTheirPart(t *testing.T) {
	t.Parallel()
	requirement := func(change func(members)) func(*semanticsCase) {
		return func(c *semanticsCase) { changeEntry(c.fields, "requirements", 0, change) }
	}
	technical := func(state string) func(*semanticsCase) {
		return func(c *semanticsCase) {
			changeObject(c.fields, "parts", set("technical", specPart(state, "Fixture reason.")))
			c.declare("abstract", "functional", "notes")
		}
	}
	withoutTechnical := func(state string, refs ...canonicaljson.Value) func(*semanticsCase) {
		return func(c *semanticsCase) {
			technical(state)(c)
			changeEntry(c.fields, "requirements", 0, set("technical_refs", list(refs...)))
			changeEntry(c.fields, "requirements", 1, set("technical_refs", list(refs...)))
		}
	}
	base := func(t *testing.T) *semanticsCase { return specCase(t, false) }
	runSemanticsRules(t, base, []semanticsRule{
		{name: "section that the part does not hold", breach: requirement(set("functional_refs", list(partRef("functional", "other")))), answer: outsidePartAnswer},
		{name: "whole document where the part holds a section", breach: func(c *semanticsCase) {
			changeObject(c.fields, "parts", set("technical", specPart("present", nil, partRef("design", "storage"))))
		}, answer: outsidePartAnswer},
		{name: "second reference outside the part", breach: requirement(set("functional_refs", list(partRef("functional", nil), partRef("notes", nil)))), answer: outsidePartAnswer},
		{name: "every reference of the part", breach: requirement(set("functional_refs", list(partRef("functional", nil), partRef("functional", "scope")))), answer: semanticsAccepted},
		{name: "no technical refs into a technical part that is not applicable", breach: withoutTechnical("not_applicable"), answer: semanticsAccepted},
		{name: "technical refs into a technical part that is not applicable", breach: technical("not_applicable"), answer: outsidePartAnswer},
		{name: "no technical refs into an unresolved technical part", breach: withoutTechnical("unresolved"), answer: "requirement needs technical refs"},
		{name: "technical refs into an unresolved technical part", breach: technical("unresolved"), answer: outsidePartAnswer},
		{name: "no functional refs into an unresolved functional part", breach: func(c *semanticsCase) {
			changeObject(c.fields, "parts", set("functional", specPart("unresolved", "Fixture reason.")))
			changeEntry(c.fields, "requirements", 0, set("functional_refs", list()))
		}, answer: "requirement needs functional refs"},
	})
}

// A phase carries every verification ID of each requirement it names.
func TestSpecPhaseCoversTheVerifiersOfItsRequirements(t *testing.T) {
	t.Parallel()
	phases := func(phases ...canonicaljson.Value) func(*semanticsCase) {
		return func(c *semanticsCase) { c.fields["phases"] = list(phases...) }
	}
	base := func(t *testing.T) *semanticsCase { return specCase(t, false) }
	runSemanticsRules(t, base, []semanticsRule{
		{name: "no phases", breach: phases(), answer: semanticsAccepted},
		{name: "phase without requirements", breach: phases(specPhase("idle", list(), list())), answer: semanticsAccepted},
		{name: "phase with a verifier of no requirement", breach: phases(specPhase("build", list("r-01"), list("manual", "unit"))), answer: semanticsAccepted},
		{name: "one phase per requirement", breach: phases(specPhase("z-first", list("r-02"), list("e2e")), specPhase("a-second", list("r-01"), list("unit"))), answer: semanticsAccepted},
		{name: "requirement in two phases", breach: phases(specPhase("first", list("r-01"), list("unit")), specPhase("second", list("r-01", "r-02"), list("e2e", "unit"))), answer: semanticsAccepted},
		{name: "phase with the verifier of another requirement only", breach: phases(specPhase("build", list("r-01"), list("e2e"))), answer: "phase does not cover requirement verifier unit"},
		{name: "phase without verifiers", breach: phases(specPhase("build", list("r-02"), list())), answer: "phase does not cover requirement verifier e2e"},
		{name: "verifier covered by an earlier phase only", breach: phases(specPhase("first", list("r-01"), list("unit")), specPhase("second", list("r-01", "r-02"), list("e2e"))), answer: "phase does not cover requirement verifier unit"},
		{name: "phase of a verification ID as requirement", breach: phases(specPhase("build", list("unit"), list("unit"))), answer: unknownRequirement},
		{name: "phase without the second verifier of a requirement", breach: func(c *semanticsCase) {
			changeEntry(c.fields, "requirements", 0, set("verification_ids", list("manual", "unit")))
			c.fields["phases"] = list(specPhase("build", list("r-01"), list("unit")))
		}, answer: "phase does not cover requirement verifier manual"},
		{name: "phases without requirements to name", breach: func(c *semanticsCase) {
			c.fields["requirements"] = list()
		}, answer: unknownRequirement},
	})
}

func TestAssessmentChecks(t *testing.T) {
	t.Parallel()
	check := func(i int, change func(members)) func(*semanticsCase) {
		return func(c *semanticsCase) { changeEntry(c.fields, "checks", i, change) }
	}
	notReady := func(outcome string, then func(*semanticsCase)) func(*semanticsCase) {
		return func(c *semanticsCase) {
			c.fields["outcome"] = outcome
			changeEntry(c.fields, "checks", 2, set("outcome", "fail"))
			changeEntry(c.fields, "checks", 5, set("outcome", "unknown"))
			c.fields["open_questions"] = list("Fixture question?")
			then(c)
		}
	}
	nothing := func(*semanticsCase) {}
	base := func(t *testing.T) *semanticsCase { return assessmentCase(t, false) }
	runSemanticsRules(t, base, []semanticsRule{
		{name: "no checks", breach: func(c *semanticsCase) { c.fields["checks"] = list() }, answer: sixChecksAnswer},
		{name: "seven checks", breach: func(c *semanticsCase) {
			c.fields["checks"] = append(contentArray(c.fields, "checks"), assessmentCheck("zzz", "pass"))
		}, answer: sixChecksAnswer},
		{name: "third check under another ID", breach: check(2, set("id", "problem")), answer: "assessment check must be problem_coverage"},
		{name: "fifth check under another ID", breach: check(4, set("id", "scope")), answer: "assessment check must be scope_and_phases"},
		{name: "last check under another ID", breach: check(5, set("id", "zzz")), answer: "assessment check must be three_parts"},
		{name: "last check failing", breach: check(5, set("outcome", "fail")), answer: failingCheckAnswer},
		{name: "last check's evidence not declared", breach: check(5, set("evidence_document_ids", list("gone"))), answer: "document gone is missing"},
		{name: "evidence shared by two checks", breach: check(3, set("evidence_document_ids", list("evidence"))), answer: semanticsAccepted},
		{name: "needs work with failing checks and open questions", breach: notReady("needs_work", nothing), answer: semanticsAccepted},
		{name: "rejected with failing checks and open questions", breach: notReady("rejected", nothing), answer: semanticsAccepted},
		{name: "rejected with five checks", breach: notReady("rejected", func(c *semanticsCase) { c.fields["checks"] = contentArray(c.fields, "checks")[:5] }), answer: sixChecksAnswer},
		{name: "rejected with a check under another ID", breach: notReady("rejected", check(3, set("id", "reco"))), answer: "assessment check must be recovery"},
		{name: "rejected on evidence that is not declared", breach: notReady("rejected", check(2, set("evidence_document_ids", list("gone")))), answer: "document gone is missing"},
		{name: "rejected with a document that nothing refers to", breach: notReady("rejected", func(c *semanticsCase) { c.declareAlso("spare") }), answer: "unused document spare"},
	})
}

func TestSelectionHasASolutionAndWithdrawalHasNone(t *testing.T) {
	t.Parallel()
	for _, published := range []bool{false, true} {
		selection, withdrawal := selectionCase(t, "spec_select", published), selectionCase(t, "spec_withdraw", published)
		selection.fields["solution"], withdrawal.fields["solution"] = withdrawal.fields["solution"], selection.fields["solution"]
		if got, want := selection.answer(), "select requires an exact solution"; got != want {
			t.Errorf("selection, published %t: got %q, want %q", published, got, want)
		}
		if got, want := withdrawal.answer(), "withdraw requires null solution"; got != want {
			t.Errorf("withdrawal, published %t: got %q, want %q", published, got, want)
		}
	}
}

// Documents and the verification IDs of a requirement are visited in map
// order, so which of several the answer names changes from run to run.
func TestTaskContentSemanticsNamesOneOfSeveral(t *testing.T) {
	t.Parallel()
	oneOf := func(c *semanticsCase, answers ...string) {
		t.Helper()
		if got := c.validator(); !slices.Contains(answers, got) {
			t.Errorf("the validator answers %q, want one of %q", got, answers)
		}
		if got := strings.TrimPrefix(c.decoded(), "task_content_invalid_input: "); !slices.Contains(answers, got) {
			t.Errorf("decoded as %q, want one of %q", got, answers)
		}
	}
	c := problemCase(t, false)
	c.declareAlso("spare")
	c.declareAlso("surplus")
	oneOf(c, "unused document spare", "unused document surplus")
	c = specCase(t, false)
	changeEntry(c.fields, "requirements", 0, set("verification_ids", list("manual", "unit")))
	changeEntry(c.fields, "phases", 0, set("verification_ids", list("e2e")))
	oneOf(c, "phase does not cover requirement verifier manual", "phase does not cover requirement verifier unit")
}

// The schema answers before the validator sees the content.
func TestTaskContentSchemaAnswersBeforeSemantics(t *testing.T) {
	t.Parallel()
	const invalid = "task_content_invalid_input: "
	for _, tc := range []struct {
		name            string
		c               *semanticsCase
		change          func(c *semanticsCase)
		schema, without string
	}{
		{"unknown member", problemCase(t, false), func(c *semanticsCase) { c.fields["future"] = true }, "object has missing or unknown fields", "document absent is missing"},
		{"unknown member of a claim", problemCase(t, false), func(c *semanticsCase) { changeEntry(c.fields, "claims", 0, set("future", true)) }, "claims: entry 0: object has missing or unknown fields", "document absent is missing"},
		{"missing member", problemCase(t, false), func(c *semanticsCase) { delete(c.fields, "deadline") }, "object has missing or unknown fields", "document absent is missing"},
		{"sources out of order", problemCase(t, false), func(c *semanticsCase) { slices.Reverse(contentArray(c.fields, "sources")) }, "sources: entries must be unique and sorted", "document absent is missing"},
		{"manifest decoded as a draft", problemCase(t, true), func(c *semanticsCase) { c.published = false }, "object has missing or unknown fields", "document absent is missing"},
	} {
		tc.c.fields["problem_document_id"] = "absent"
		tc.change(tc.c)
		if got := tc.c.decoded(); got != invalid+tc.schema {
			t.Errorf("%s: decoded as %q, want %q", tc.name, got, invalid+tc.schema)
		}
		if got := tc.c.validator(); got != tc.without {
			t.Errorf("%s: the validator answers %q, want %q", tc.name, got, tc.without)
		}
	}
}

// The validator reads content that the schema would reject as far as it can:
// a member of the wrong type reads as absent, and members that it does not
// know are not read.
func TestTaskContentSemanticsOfContentTheSchemaRejects(t *testing.T) {
	t.Parallel()
	document := func(id canonicaljson.Value) canonicaljson.Value { return contentObject(members{"id": id}) }
	requirement := func(id string, verification ...canonicaljson.Value) canonicaljson.Value {
		return specRequirement(id, list(partRef("functional", nil)), list(), verification)
	}
	spec := func(change members) members {
		fields := members{
			"documents": list(document("functional")),
			"parts": contentObject(members{
				"abstract":   specPart("unresolved", "Fixture reason."),
				"functional": specPart("present", nil, partRef("functional", nil)),
				"technical":  specPart("not_applicable", "Fixture reason."),
			}),
		}
		for name, value := range change {
			fields[name] = value
		}
		return fields
	}
	for _, tc := range []struct {
		name      string
		operation string
		fields    members
		want      string
		// asManifest is the answer to the content as a manifest, where it
		// differs: a manifest is never a generated draft.
		asManifest string
	}{
		{name: "problem without members", operation: "problem_record", fields: nil, want: "document  is missing"},
		{name: "spec without members", operation: "spec_record", fields: nil, want: "invalid abstract part state"},
		{name: "assessment without members", operation: "spec_assess", fields: nil, want: sixChecksAnswer},
		{name: "selection without members", operation: "spec_select", fields: nil, want: "select requires an exact solution"},
		{name: "withdrawal without members", operation: "spec_withdraw", fields: nil, want: semanticsAccepted},
		{name: "empty problem", operation: "problem_record", fields: members{}, want: "document  is missing"},
		{name: "unknown operation", operation: "task_create", fields: members{"solution": nil, "checks": list(), "parts": nil}, want: semanticsAccepted},
		{name: "unknown operation with a document", operation: "task_create", fields: members{"documents": list(document("left"))}, want: "unused document left"},
		{name: "no operation", operation: "", fields: nil, want: semanticsAccepted},
		{name: "selection with a document", operation: "spec_select", fields: members{"solution": contentObject(members{}), "documents": list(document("left"))}, want: "unused document left"},
		{name: "selection without a solution and with a document", operation: "spec_select", fields: members{"documents": list(document("left"))}, want: "select requires an exact solution"},
		{name: "withdrawal of a solution that is not an object", operation: "spec_withdraw", fields: members{"solution": false}, want: "withdraw requires null solution"},
		{name: "withdrawal with a document", operation: "spec_withdraw", fields: members{"documents": list(document("left"))}, want: "unused document left"},
		{name: "documents that are not a list", operation: "problem_record", fields: members{"documents": document("problem"), "problem_document_id": "problem"}, want: "document problem is missing"},
		{name: "document without an ID is declared under the empty ID", operation: "problem_record", fields: members{"documents": list("problem", document(int64(7)))}, want: semanticsAccepted},
		{name: "problem document ID that is not text", operation: "problem_record", fields: members{"documents": list(document("problem")), "problem_document_id": int64(7)}, want: "document  is missing"},
		{name: "document declared twice", operation: "problem_record", fields: members{"documents": list(document("problem"), document("problem")), "problem_document_id": "problem"}, want: semanticsAccepted},
		{name: "origin that is not an object is authored", operation: "problem_record", fields: members{"origin": "task_create_summary", "documents": list(document("problem")), "problem_document_id": "problem"}, want: semanticsAccepted},
		{name: "generated draft with a document declared twice", operation: "problem_record", fields: members{"origin": contentObject(members{"kind": "task_create_summary"}), "documents": list(document("problem"), document("problem")), "problem_document_id": "problem"}, want: generatedDraft, asManifest: semanticsAccepted},
		{name: "generated draft whose lists are not lists", operation: "problem_record", fields: members{"origin": contentObject(members{"kind": "task_create_summary"}), "documents": "none", "sources": "none", "claims": "none", "problem_document_id": "problem"}, want: semanticsAccepted, asManifest: "document problem is missing"},
		{name: "generated draft with a false deadline", operation: "problem_record", fields: members{"origin": contentObject(members{"kind": "task_create_summary"}), "deadline": false, "problem_document_id": "problem"}, want: generatedDraft, asManifest: "document problem is missing"},
		{name: "source document ID that is not text", operation: "problem_record", fields: members{"documents": list(document("problem")), "problem_document_id": "problem", "sources": list(contentObject(members{"id": "letter", "document_id": int64(7)}))}, want: "document  is missing"},
		{name: "sources that are not objects", operation: "problem_record", fields: members{"documents": list(document("problem")), "problem_document_id": "problem", "sources": list("letter", nil), "claims": list(contentObject(members{"verification": "unverified"}))}, want: semanticsAccepted},
		{name: "claim whose evidence is not a list", operation: "problem_record", fields: members{"documents": list(document("problem")), "problem_document_id": "problem", "sources": list(contentObject(members{"id": "letter"})), "claims": list(contentObject(members{"source_id": "letter", "verification": "confirmed", "reason": "Fixture reason.", "evidence_document_ids": "problem"}))}, want: verifiedAnswer},
		{name: "claim without a verification is verified", operation: "problem_record", fields: members{"documents": list(document("problem")), "problem_document_id": "problem", "sources": list(contentObject(members{"id": "letter"})), "claims": list(contentObject(members{"source_id": "letter"}))}, want: verifiedAnswer},
		{name: "deadline that is not an object names the claim without an ID", operation: "problem_record", fields: members{"documents": list(document("problem")), "problem_document_id": "problem", "sources": list(contentObject(members{})), "claims": list(contentObject(members{"verification": "unverified"})), "deadline": "soon"}, want: semanticsAccepted},
		{name: "unknown members", operation: "problem_record", fields: members{"documents": list(contentObject(members{"id": "problem", "future": true})), "problem_document_id": "problem", "future": list(document("left")), "sources": list(contentObject(members{"id": "letter", "future": "left"}))}, want: semanticsAccepted},
		{name: "parts that are not an object", operation: "spec_record", fields: members{"parts": list()}, want: "invalid abstract part state"},
		{name: "part without a state", operation: "spec_record", fields: spec(members{"parts": contentObject(members{"abstract": contentObject(members{"reason": "Fixture reason."})})}), want: "invalid functional part state"},
		{name: "part of an unknown state", operation: "spec_record", fields: spec(members{"documents": list(), "parts": contentObject(members{"abstract": specPart("later", "Fixture reason."), "functional": specPart("later", "Fixture reason."), "technical": specPart("later", "Fixture reason.")})}), want: semanticsAccepted},
		{name: "supporting document named twice", operation: "spec_record", fields: spec(members{"documents": list(document("functional"), document("notes")), "supporting": list(contentObject(members{"document_id": "notes"}), contentObject(members{"document_id": "notes"}))}), want: "supporting document duplicates a part"},
		{name: "supporting entry that is not an object", operation: "spec_record", fields: spec(members{"supporting": list("notes")}), want: "document  is missing"},
		{name: "requirement whose refs are not a list", operation: "spec_record", fields: spec(members{"requirements": list(contentObject(members{"id": "r-01", "functional_refs": partRef("functional", nil)}))}), want: "requirement needs functional refs"},
		{name: "requirement ref that is not an object", operation: "spec_record", fields: spec(members{"requirements": list(contentObject(members{"id": "r-01", "functional_refs": list("functional")}))}), want: outsidePartAnswer},
		{name: "requirement ref whose section is not text", operation: "spec_record", fields: spec(members{"requirements": list(specRequirement("r-01", list(partRef("functional", int64(7))), list(), list("unit")))}), want: semanticsAccepted},
		{name: "requirement registered twice keeps its last verifiers", operation: "spec_record", fields: spec(members{"requirements": list(requirement("r-01", "unit"), requirement("r-01", "e2e")), "phases": list(specPhase("build", list("r-01"), list("e2e")))}), want: semanticsAccepted},
		{name: "requirement registered twice loses its first verifiers", operation: "spec_record", fields: spec(members{"requirements": list(requirement("r-01", "unit"), requirement("r-01", "e2e")), "phases": list(specPhase("build", list("r-01"), list("unit")))}), want: "phase does not cover requirement verifier e2e"},
		{name: "requirement without an ID", operation: "spec_record", fields: spec(members{"requirements": list(contentObject(members{"functional_refs": list(partRef("functional", nil)), "verification_ids": list("unit")})), "phases": list(specPhase("build", list(""), list("unit")))}), want: semanticsAccepted},
		{name: "phase that is not an object", operation: "spec_record", fields: spec(members{"phases": list("build", nil)}), want: semanticsAccepted},
		{name: "checks that are not a list", operation: "spec_assess", fields: members{"checks": assessmentCheck("acceptance_coverage", "pass")}, want: sixChecksAnswer},
		{name: "checks that are not objects", operation: "spec_assess", fields: members{"checks": list("acceptance_coverage", nil, nil, nil, nil, nil)}, want: "assessment check must be acceptance_coverage"},
		{name: "outcome that is not text is not ready", operation: "spec_assess", fields: members{"outcome": true, "open_questions": list("Fixture question?"), "checks": []canonicaljson.Value{
			assessmentCheck("acceptance_coverage", "fail"), assessmentCheck("implementation_basis", "fail"), assessmentCheck("problem_coverage", "fail"), assessmentCheck("recovery", "fail"), assessmentCheck("scope_and_phases", "fail"), assessmentCheck("three_parts", "fail"),
		}}, want: semanticsAccepted},
	} {
		for _, published := range []bool{false, true} {
			for _, internal := range []bool{false, true} {
				c := &semanticsCase{t: t, operation: tc.operation, published: published, internal: internal, fields: tc.fields}
				want := tc.want
				if published && tc.asManifest != "" {
					want = tc.asManifest
				}
				if got := c.validator(); got != want {
					t.Errorf("%s, published %t, internal %t: got %q, want %q", tc.name, published, internal, got, want)
				}
			}
		}
	}
}

// The validator takes the IDs in evidence, verification and requirement lists
// as text without looking: the schema has rejected anything else.
func TestTaskContentSemanticsTrustsTheSchemaForListsOfIDs(t *testing.T) {
	t.Parallel()
	const number = "claims: entry 0: evidence_document_ids: entry 0: must be a 1..128 character publication key"
	c := problemCase(t, false)
	changeEntry(c.fields, "claims", 0, set("evidence_document_ids", list(int64(7))))
	if got, want := c.decoded(), "task_content_invalid_input: "+number; got != want {
		t.Fatalf("decoded as %q, want %q", got, want)
	}
	for name, c := range map[string]*semanticsCase{
		"claim's evidence": func() *semanticsCase {
			c := problemCase(t, false)
			changeEntry(c.fields, "claims", 0, set("evidence_document_ids", list(int64(7))))
			return c
		}(),
		"check's evidence": func() *semanticsCase {
			c := assessmentCase(t, false)
			changeEntry(c.fields, "checks", 0, set("evidence_document_ids", list(int64(7))))
			return c
		}(),
		"requirement's verification IDs": func() *semanticsCase {
			c := specCase(t, false)
			changeEntry(c.fields, "requirements", 0, set("verification_ids", list(int64(7))))
			return c
		}(),
		"phase's verification IDs": func() *semanticsCase {
			c := specCase(t, false)
			changeEntry(c.fields, "phases", 0, set("verification_ids", list(int64(7))))
			return c
		}(),
		"phase's requirement IDs": func() *semanticsCase {
			c := specCase(t, false)
			changeEntry(c.fields, "phases", 0, set("requirement_ids", list(int64(7))))
			return c
		}(),
	} {
		func() {
			defer func() {
				if _, asserted := recover().(*runtime.TypeAssertionError); !asserted {
					t.Errorf("%s: the validator did not take the ID as text", name)
				}
			}()
			t.Errorf("%s: the validator answers %q", name, c.validator())
		}()
	}
}
