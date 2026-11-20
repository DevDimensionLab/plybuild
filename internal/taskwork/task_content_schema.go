package taskwork

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

const taskManifestLimit = 256 << 10
const taskDocumentLimit = 4 << 20

type contentRule func(canonicaljson.Value) error

type contentUnknownFields struct {
	names  []string
	detail string
}

func (e *contentUnknownFields) Error() string { return e.detail }

func contentExact(fields map[string]contentRule) contentRule {
	return func(v canonicaljson.Value) error {
		o, ok := v.(canonicaljson.Object)
		if !ok {
			return fmt.Errorf("object has missing or unknown fields")
		}
		unknown := []string{}
		for _, m := range o {
			if _, known := fields[m.Name]; !known {
				unknown = append(unknown, m.Name)
			}
		}
		if len(o) != len(fields) {
			if len(unknown) > 0 {
				return &contentUnknownFields{unknown, "object has missing or unknown fields"}
			}
			return fmt.Errorf("object has missing or unknown fields")
		}
		seen := map[string]bool{}
		for _, m := range o {
			rule, ok := fields[m.Name]
			if !ok || seen[m.Name] {
				if !ok {
					return &contentUnknownFields{[]string{m.Name}, fmt.Sprintf("unknown or duplicate field %s", m.Name)}
				}
				return fmt.Errorf("unknown or duplicate field %s", m.Name)
			}
			seen[m.Name] = true
			if err := rule(m.Value); err != nil {
				return fmt.Errorf("%s: %w", m.Name, err)
			}
		}
		return nil
	}
}
func contentNullable(rule contentRule) contentRule {
	return func(v canonicaljson.Value) error {
		if v == nil {
			return nil
		}
		return rule(v)
	}
}

// Drafts keep their published registry_upgrade member, which must be null:
// the registry has a single format and is never upgraded by a command.
func contentNoRegistryUpgrade(v canonicaljson.Value) error {
	if v != nil {
		return fmt.Errorf("must be null; registry format %d has no upgrade", RegistryFormat)
	}
	return nil
}
func contentText(max int) contentRule {
	return func(v canonicaljson.Value) error {
		s, ok := v.(string)
		if !ok {
			return fmt.Errorf("must be text")
		}
		return validateWorkText("text", s, max)
	}
}
func contentEnum(values ...string) contentRule {
	return func(v canonicaljson.Value) error {
		s, ok := v.(string)
		if ok {
			for _, a := range values {
				if a == s {
					return nil
				}
			}
		}
		return fmt.Errorf("must be one of %s", strings.Join(values, ", "))
	}
}
func contentInteger(min, max int64) contentRule {
	return func(v canonicaljson.Value) error {
		n, ok := v.(int64)
		if !ok || n < min || n > max {
			return fmt.Errorf("must be an integer in %d..%d", min, max)
		}
		return nil
	}
}
func contentDigest(v canonicaljson.Value) error {
	s, ok := v.(string)
	if !ok || !digestPattern.MatchString(s) {
		return fmt.Errorf("must be sha256: and 64 lowercase hexadecimal digits")
	}
	return nil
}
func contentOID(v canonicaljson.Value) error {
	s, ok := v.(string)
	if !ok || !validOIDText(s) {
		return fmt.Errorf("must be a full lowercase Git OID")
	}
	return nil
}
func contentSlug(v canonicaljson.Value) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("must be a slug")
	}
	return validateWorkIdentifier("content", s)
}

var contentKeyPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._/-]*$`)

func contentKey(v canonicaljson.Value) error {
	s, ok := v.(string)
	if !ok || len(s) < 1 || len(s) > 128 || !contentKeyPattern.MatchString(s) || strings.Contains(s, "//") || strings.HasSuffix(s, "/") {
		return fmt.Errorf("must be a 1..128 character publication key")
	}
	for _, p := range strings.Split(s, "/") {
		if p == ".." {
			return fmt.Errorf("key contains a parent segment")
		}
	}
	return nil
}
func contentUTC(v canonicaljson.Value) error {
	s, ok := v.(string)
	if !ok || !strings.HasSuffix(s, "Z") {
		return fmt.Errorf("must be a UTC timestamp ending in Z")
	}
	_, e := time.Parse(time.RFC3339Nano, s)
	return e
}
func contentPath(v canonicaljson.Value) error {
	s, ok := v.(string)
	if !ok || !filepath.IsAbs(s) || filepath.Clean(s) != s || strings.ContainsRune(s, 0) {
		return fmt.Errorf("must be an absolute clean physical path")
	}
	return nil
}
func contentList(max int, item contentRule, key func(canonicaljson.Value) string, ordered bool) contentRule {
	return func(v canonicaljson.Value) error {
		a, ok := v.([]canonicaljson.Value)
		if !ok || len(a) > max {
			return fmt.Errorf("must be an explicit array with at most %d entries", max)
		}
		seen := map[string]bool{}
		last := ""
		for i, x := range a {
			if e := item(x); e != nil {
				return fmt.Errorf("entry %d: %w", i, e)
			}
			if key != nil {
				k := key(x)
				if seen[k] || (!ordered && i > 0 && k <= last) {
					return fmt.Errorf("entries must be unique%s", map[bool]string{true: "", false: " and sorted"}[ordered])
				}
				seen[k] = true
				last = k
			}
		}
		return nil
	}
}
func contentIDKey(v canonicaljson.Value) string     { return contentString(contentFields(v), "id") }
func contentStringKey(v canonicaljson.Value) string { s, _ := v.(string); return s }
func contentStringList(max int, rule contentRule) contentRule {
	return contentList(max, rule, contentStringKey, false)
}
func contentRevisionRule() contentRule {
	return contentExact(map[string]contentRule{"revision": contentInteger(1, 2147483647), "manifest_sha256": contentDigest})
}
func contentDecisionRule(prefix string) contentRule {
	return contentExact(map[string]contentRule{"id": func(v canonicaljson.Value) error {
		s, ok := v.(string)
		if !ok || !regexp.MustCompile("^"+prefix+"[0-9a-f]{32}$").MatchString(s) {
			return fmt.Errorf("invalid %s identity", prefix)
		}
		return nil
	}, "manifest_sha256": contentDigest})
}
func contentRecorderRule() contentRule {
	return contentExact(map[string]contentRule{"actor_claim": contentText(256), "control_surface": contentText(256), "recorded_at_utc": contentUTC})
}
func contentPartRefRule() contentRule {
	return contentExact(map[string]contentRule{"document_id": contentKey, "section": contentNullable(contentText(128))})
}
func contentPartKey(v canonicaljson.Value) string {
	m := contentFields(v)
	s, ok := m["section"].(string)
	if !ok {
		return contentString(m, "document_id") + "\x00"
	}
	return contentString(m, "document_id") + "\x01" + s
}
func contentPartRefs() contentRule {
	return contentList(32, contentPartRefRule(), contentPartKey, false)
}
func contentGitRule() contentRule {
	return contentExact(map[string]contentRule{"git_common_dir": contentPath, "ref": func(v canonicaljson.Value) error {
		s, ok := v.(string)
		if !ok || !validRefText(s) || !strings.HasPrefix(s, "refs/") {
			return fmt.Errorf("must be a full Git ref")
		}
		return nil
	}, "oid": contentOID, "tree": contentOID, "repo_path": func(v canonicaljson.Value) error {
		s, ok := v.(string)
		if !ok || s == "." || path.Clean(s) != s || strings.HasPrefix(s, "/") || strings.ContainsAny(s, "\\\x00") || strings.HasPrefix(s, "../") {
			return fmt.Errorf("invalid repository path")
		}
		return nil
	}, "blob": contentOID})
}
func contentDescriptorRule() contentRule {
	return contentExact(map[string]contentRule{"id": contentKey, "locator": contentPath, "sha256": contentDigest, "size_bytes": contentInteger(1, taskDocumentLimit), "media_type": contentEnum("text/markdown", "text/plain"), "provenance": contentExact(map[string]contentRule{"origin_locator": contentNullable(contentPath), "git": contentNullable(contentGitRule())})})
}
func contentDocumentInputRule() contentRule {
	return contentExact(map[string]contentRule{"id": contentKey, "source": func(v canonicaljson.Value) error {
		switch contentString(contentFields(v), "kind") {
		case "file":
			return contentExact(map[string]contentRule{"kind": contentEnum("file"), "locator": contentPath, "sha256": contentDigest, "size_bytes": contentInteger(1, taskDocumentLimit), "media_type": contentEnum("text/markdown", "text/plain"), "git_provenance": contentNullable(contentGitRule())})(v)
		case "snapshot":
			rules := map[string]contentRule{"kind": contentEnum("snapshot"), "manifest_sha256": contentDigest, "document_id": contentKey}
			if _, ok := contentFields(v)["task_id"]; ok {
				rules["task_id"] = contentSlug
			}
			return contentExact(rules)(v)
		default:
			return fmt.Errorf("source kind must be file or snapshot")
		}
	}})
}
func contentBasisRule() contentRule {
	return contentExact(map[string]contentRule{
		"project_id": contentSlug, "repo_id": contentSlug, "git_common_dir": contentPath, "epic_id": contentSlug,
		"parent_worktree_id": func(v canonicaljson.Value) error {
			s, ok := v.(string)
			if !ok || !worktreeIDPattern.MatchString(s) {
				return fmt.Errorf("invalid worktree ID")
			}
			return nil
		},
		"parent_ref": func(v canonicaljson.Value) error {
			s, ok := v.(string)
			if !ok || !validFullBranchRef(s) {
				return fmt.Errorf("invalid parent ref")
			}
			return nil
		},
		"parent_oid": contentOID, "parent_tree": contentOID, "start_oid": contentOID, "start_tree": contentOID,
	})
}
func taskContentSchema(operation string, published bool, internal bool) map[string]contentRule {
	kinds := map[string]string{"problem_record": "WorkspaceTaskProblemDraft@1", "spec_record": "WorkspaceTaskSpecDraft@1", "spec_assess": "WorkspaceTaskSpecAssessmentDraft@1", "spec_select": "WorkspaceTaskSolutionSelectionDraft@1", "spec_withdraw": "WorkspaceTaskSolutionSelectionDraft@1"}
	if published {
		kinds = map[string]string{"problem_record": "WorkspaceTaskProblemRevision@1", "spec_record": "WorkspaceTaskSpecRevision@1", "spec_assess": "WorkspaceTaskSpecAssessment@1", "spec_select": "WorkspaceTaskSolutionSelection@1", "spec_withdraw": "WorkspaceTaskSolutionSelection@1"}
	}
	s := map[string]contentRule{"kind": contentEnum(kinds[operation]), "schema_version": contentInteger(1, 1), "format": contentEnum("json"), "format_version": contentInteger(1, 1), "canonicalization": contentEnum("RFC8785"), "publication_key": contentKey, "task_id": contentSlug, "recorder": contentRecorderRule()}
	if !published {
		s["registry_upgrade"] = contentNoRegistryUpgrade
	} else {
		s["source_draft_sha256"] = contentDigest
		s["recorded_at_utc"] = contentUTC
	}
	doc := contentDocumentInputRule()
	if published {
		doc = contentDescriptorRule()
	}
	if operation == "problem_record" || operation == "spec_record" {
		s["title"] = contentText(128)
		s["change_reason"] = contentText(2000)
		if published {
			s["previous"] = contentNullable(contentRevisionRule())
			s["revision"] = contentInteger(1, 2147483647)
		} else {
			s["expected_previous"] = contentNullable(contentRevisionRule())
		}
	}
	if operation == "problem_record" || operation == "spec_record" || operation == "spec_assess" {
		s["documents"] = contentList(32, doc, contentIDKey, false)
	}
	switch operation {
	case "problem_record":
		s["summary"] = contentText(2048)
		s["problem_document_id"] = contentKey
		s["origin"] = func(v canonicaljson.Value) error {
			k := contentString(contentFields(v), "kind")
			if k == "legacy_summary_import" {
				return contentExact(map[string]contentRule{"kind": contentEnum(k), "legacy_summary_sha256": contentDigest})(v)
			}
			if k == "authored" || (internal && k == "task_create_summary") {
				return contentExact(map[string]contentRule{"kind": contentEnum(k)})(v)
			}
			return fmt.Errorf("invalid problem origin")
		}
		s["sources"] = contentList(32, contentExact(map[string]contentRule{"id": contentKey, "origin_kind": contentEnum("email", "issue", "document", "observation", "url"), "origin_claim": contentText(2000), "received_at_utc": contentNullable(contentUTC), "document_id": contentNullable(contentKey), "withholding_reason": contentNullable(contentText(2000))}), contentIDKey, false)
		s["claims"] = contentList(128, contentExact(map[string]contentRule{"id": contentKey, "source_id": contentKey, "text": contentText(2000), "verification": contentEnum("unverified", "confirmed", "refuted"), "reason": contentNullable(contentText(2000)), "evidence_document_ids": contentStringList(32, contentKey)}), contentIDKey, false)
		s["deadline"] = contentNullable(contentExact(map[string]contentRule{"date": func(v canonicaljson.Value) error {
			t, ok := v.(string)
			if !ok {
				return fmt.Errorf("invalid date")
			}
			_, e := time.Parse("2006-01-02", t)
			return e
		}, "claim_id": contentKey}))
	case "spec_record":
		s["spec_id"] = contentSlug
		s["problem"] = contentRevisionRule()
		part := contentExact(map[string]contentRule{"state": contentEnum("present", "unresolved", "not_applicable"), "reason": contentNullable(contentText(2000)), "documents": contentPartRefs()})
		s["parts"] = contentExact(map[string]contentRule{"abstract": part, "functional": part, "technical": part})
		s["supporting"] = contentList(32, contentExact(map[string]contentRule{"document_id": contentKey, "usage": contentEnum("normative", "informative")}), func(v canonicaljson.Value) string { return contentString(contentFields(v), "document_id") }, false)
		s["requirements"] = contentList(128, contentExact(map[string]contentRule{"id": contentKey, "functional_refs": contentPartRefs(), "acceptance": contentText(2000), "verification_ids": contentStringList(128, contentKey), "technical_refs": contentPartRefs()}), contentIDKey, false)
		s["removed_requirement_ids"] = contentStringList(128, contentKey)
		s["phases"] = contentList(32, contentExact(map[string]contentRule{"id": contentKey, "purpose": contentText(2000), "requirement_ids": contentStringList(128, contentKey), "entry_criteria": contentStringList(128, contentText(2000)), "exit_criteria": contentStringList(128, contentText(2000)), "verification_ids": contentStringList(128, contentKey)}), contentIDKey, true)
		s["implementation_basis"] = contentBasisRule()
		s["dependencies"] = contentList(0, contentKey, nil, false)
	case "spec_assess":
		s["spec_id"] = contentSlug
		s["spec"] = contentRevisionRule()
		s["outcome"] = contentEnum("ready", "needs_work", "rejected")
		s["reason"] = contentText(2000)
		s["open_questions"] = contentStringList(128, contentText(2000))
		s["checks"] = contentList(128, contentExact(map[string]contentRule{"id": contentKey, "outcome": contentEnum("pass", "fail", "unknown"), "reason": contentText(2000), "evidence_document_ids": contentStringList(32, contentKey)}), contentIDKey, false)
		if published {
			s["id"] = func(v canonicaljson.Value) error {
				return contentDecisionRule("asm_")(contentObject(map[string]canonicaljson.Value{"id": v, "manifest_sha256": "sha256:" + strings.Repeat("0", 64)}))
			}
			s["ordinal"] = contentInteger(1, 2147483647)
			s["previous"] = contentNullable(contentDecisionRule("asm_"))
		} else {
			s["expected_previous_assessment"] = contentNullable(contentDecisionRule("asm_"))
		}
	case "spec_select", "spec_withdraw":
		action := "select"
		if operation == "spec_withdraw" {
			action = "withdraw"
		}
		s["action"] = contentEnum(action)
		s["reason"] = contentText(2000)
		s["solution"] = contentNullable(contentExact(map[string]contentRule{"spec_id": contentSlug, "spec": contentRevisionRule(), "problem": contentRevisionRule(), "assessment": contentDecisionRule("asm_")}))
		s["human_decision"] = contentExact(map[string]contentRule{"actor_claim": contentText(256), "decided_at_utc": contentUTC, "source": contentEnum("human_cli", "explicit_human_instruction"), "statement": contentText(2000)})
		if published {
			s["id"] = func(v canonicaljson.Value) error {
				return contentDecisionRule("sel_")(contentObject(map[string]canonicaljson.Value{"id": v, "manifest_sha256": "sha256:" + strings.Repeat("0", 64)}))
			}
			s["ordinal"] = contentInteger(1, 2147483647)
			s["previous"] = contentNullable(contentDecisionRule("sel_"))
		} else {
			s["expected_previous_selection"] = contentNullable(contentDecisionRule("sel_"))
		}
	}
	return s
}

func decodeTaskContent(raw []byte, operation string, published, internal bool) (canonicaljson.Object, error) {
	if len(raw) > taskManifestLimit {
		return nil, contentError("task_content_invalid_input", "draft or manifest exceeds 256 KiB", nil)
	}
	v, err := canonicaljson.DecodeStrict(raw)
	if err != nil {
		return nil, contentError("task_content_invalid_input", err.Error(), err)
	}
	schema := taskContentSchema(operation, published, internal)
	if operation == "spec_record" && contentInt(contentFields(v), "schema_version") == 2 {
		if isTaskGoalSpec(v) {
			schema = taskGoalSpecSchema(published)
		}
	}
	if operation == "spec_record" {
		if _, present := contentFields(v)["delivery"]; present {
			schema["delivery"] = historicalGoalDelivery
		}
	}
	if err = contentExact(schema)(v); err == nil {
		if isTaskGoalSpec(v) {
			err = validateTaskGoalSemantics(contentFields(v))
		} else {
			err = validateTaskContentSemantics(contentFields(v), operation, published, internal)
		}
	}
	if err != nil {
		return nil, contentError("task_content_invalid_input", err.Error(), err)
	}
	b, err := canonicaljson.Marshal(v)
	if err != nil || len(b) > taskManifestLimit {
		return nil, contentError("task_content_invalid_input", "canonical manifest exceeds 256 KiB", err)
	}
	return v.(canonicaljson.Object), nil
}

// validateTaskContentSemantics checks the rules that relate the members of a
// draft or manifest whose shape the schema has accepted. The rules of the
// operation run in a fixed order and the first one that fails answers; that
// every declared document is referred to is checked last. Only the schema
// tells internal content from other content.
func validateTaskContentSemantics(m map[string]canonicaljson.Value, operation string, published, internal bool) error {
	docs := declaredContentDocuments(m)
	if err := validateOperationSemantics(m, operation, published, docs); err != nil {
		return err
	}
	return docs.unused()
}

// A selection has a solution and a withdrawal has none; the other operations
// have rules of their own.
func validateOperationSemantics(m map[string]canonicaljson.Value, operation string, published bool, docs contentDocuments) error {
	hasSolution := m["solution"] != nil
	switch {
	case operation == "problem_record":
		return validateProblemSemantics(m, published, docs)
	case operation == "spec_record":
		return validateSpecSemantics(m, docs)
	case operation == "spec_assess":
		return validateAssessmentSemantics(m, docs)
	case operation == "spec_select" && !hasSolution:
		return fmt.Errorf("select requires an exact solution")
	case operation == "spec_withdraw" && hasSolution:
		return fmt.Errorf("withdraw requires null solution")
	}
	return nil
}

// contentDocuments are the documents that a draft or manifest declares, and
// those among them that its content has referred to so far.
type contentDocuments struct {
	declared, used map[string]bool
}

func declaredContentDocuments(m map[string]canonicaljson.Value) contentDocuments {
	docs := contentDocuments{declared: map[string]bool{}, used: map[string]bool{}}
	for _, d := range contentArray(m, "documents") {
		docs.declared[contentIDKey(d)] = true
	}
	return docs
}

// use records a reference to the document id, which must be declared.
func (docs contentDocuments) use(id string) error {
	if !docs.declared[id] {
		return fmt.Errorf("document %s is missing", id)
	}
	docs.used[id] = true
	return nil
}

// useAll records a reference to each of the document IDs, in their order.
func (docs contentDocuments) useAll(ids []canonicaljson.Value) error {
	for _, id := range ids {
		if err := docs.use(id.(string)); err != nil {
			return err
		}
	}
	return nil
}

// unused rejects a declared document that nothing refers to. Documents are
// visited in map order, so which of several unused documents the error names
// varies between runs.
func (docs contentDocuments) unused() error {
	for id := range docs.declared {
		if !docs.used[id] {
			return fmt.Errorf("unused document %s", id)
		}
	}
	return nil
}

// contentStringSet is the set of the texts in a list.
func contentStringSet(values []canonicaljson.Value) map[string]bool {
	set := map[string]bool{}
	for _, v := range values {
		set[v.(string)] = true
	}
	return set
}

// A problem refers to its problem document. Its sources may each hold a
// document, its claims each name a source, and its deadline names a claim.
func validateProblemSemantics(m map[string]canonicaljson.Value, published bool, docs contentDocuments) error {
	if err := validateProblemDocument(m, published, docs); err != nil {
		return err
	}
	sources, err := problemSources(m, docs)
	if err != nil {
		return err
	}
	claims, err := problemClaims(m, sources, docs)
	if err != nil {
		return err
	}
	if m["deadline"] != nil && !claims[contentString(contentFields(m["deadline"]), "claim_id")] {
		return fmt.Errorf("deadline claim does not exist")
	}
	return nil
}

// The draft of a problem that Ply generates from a Task summary has no
// content of its own: its problem document "problem" is added on publication.
// Every other problem, and the manifest of a generated one, declares its
// problem document.
func validateProblemDocument(m map[string]canonicaljson.Value, published bool, docs contentDocuments) error {
	origin := contentString(contentFields(m["origin"]), "kind")
	generated := origin == "legacy_summary_import" || origin == "task_create_summary"
	if !generated || published {
		return docs.use(contentString(m, "problem_document_id"))
	}
	if len(docs.declared) != 0 || len(contentArray(m, "sources")) != 0 || len(contentArray(m, "claims")) != 0 || m["deadline"] != nil || contentString(m, "problem_document_id") != "problem" {
		return fmt.Errorf("generated problem requires empty documents, sources and claims, null deadline and problem ID")
	}
	return nil
}

// problemSources returns the IDs of the problem's sources, each validated
// before the next.
func problemSources(m map[string]canonicaljson.Value, docs contentDocuments) (map[string]bool, error) {
	ids := map[string]bool{}
	for _, v := range contentArray(m, "sources") {
		s := contentFields(v)
		ids[contentString(s, "id")] = true
		if err := useSourceDocument(s, docs); err != nil {
			return nil, err
		}
	}
	return ids, nil
}

// A source may hold a document, unless it is withheld.
func useSourceDocument(s map[string]canonicaljson.Value, docs contentDocuments) error {
	if s["document_id"] == nil {
		return nil
	}
	if s["withholding_reason"] != nil {
		return fmt.Errorf("withheld source cannot have a document")
	}
	return docs.use(contentString(s, "document_id"))
}

// problemClaims returns the IDs of the problem's claims, each validated
// against the given sources before the next.
func problemClaims(m map[string]canonicaljson.Value, sources map[string]bool, docs contentDocuments) (map[string]bool, error) {
	ids := map[string]bool{}
	for _, v := range contentArray(m, "claims") {
		c := contentFields(v)
		ids[contentString(c, "id")] = true
		if err := validateProblemClaim(c, sources, docs); err != nil {
			return nil, err
		}
	}
	return ids, nil
}

// A claim names a source of its problem. An unverified claim has neither
// reason nor evidence; any other claim has both.
func validateProblemClaim(c map[string]canonicaljson.Value, sources map[string]bool, docs contentDocuments) error {
	if !sources[contentString(c, "source_id")] {
		return fmt.Errorf("claim source does not exist")
	}
	evidence := contentArray(c, "evidence_document_ids")
	reasoned, evidenced := c["reason"] != nil, len(evidence) != 0
	unverified := contentString(c, "verification") == "unverified"
	if unverified && (reasoned || evidenced) {
		return fmt.Errorf("unverified claim cannot have verification evidence")
	}
	if !unverified && (!reasoned || !evidenced) {
		return fmt.Errorf("verified claim requires reason and evidence")
	}
	return docs.useAll(evidence)
}

// A Spec's parts hold document references, its supporting documents are not
// those of a part, its requirements refer into the parts, and its phases
// cover the verification of the requirements they name.
func validateSpecSemantics(m map[string]canonicaljson.Value, docs contentDocuments) error {
	parts := contentFields(m["parts"])
	refs, err := specPartRefs(parts, docs)
	if err != nil {
		return err
	}
	if err = validateSupportingDocuments(m, docs); err != nil {
		return err
	}
	technicalApplies := contentString(contentFields(parts["technical"]), "state") != "not_applicable"
	verifiers, err := specRequirementVerifiers(m, refs, technicalApplies)
	if err != nil {
		return err
	}
	return validatePhaseCoverage(m, verifiers)
}

// specPartRefs returns the keys of the document references that each part
// holds, by the part's name. The parts are validated in the order abstract,
// functional and technical.
func specPartRefs(parts map[string]canonicaljson.Value, docs contentDocuments) (map[string]map[string]bool, error) {
	refs := map[string]map[string]bool{}
	for _, name := range []string{"abstract", "functional", "technical"} {
		held, err := partRefKeys(name, contentFields(parts[name]), docs)
		if err != nil {
			return nil, err
		}
		refs[name] = held
	}
	return refs, nil
}

// partRefKeys returns the keys of the document references that the part p
// holds. The part's state is valid and its documents are declared.
func partRefKeys(name string, p map[string]canonicaljson.Value, docs contentDocuments) (map[string]bool, error) {
	if err := validatePartState(name, p); err != nil {
		return nil, err
	}
	held := map[string]bool{}
	for _, ref := range contentArray(p, "documents") {
		if err := docs.use(contentString(contentFields(ref), "document_id")); err != nil {
			return nil, err
		}
		held[contentPartKey(ref)] = true
	}
	return held, nil
}

// A present part has documents and no reason. Any other part has a reason and
// no documents, and only the technical part may be not applicable.
func validatePartState(name string, p map[string]canonicaljson.Value) error {
	state := contentString(p, "state")
	reasoned, documented := p["reason"] != nil, len(contentArray(p, "documents")) != 0
	if state == "present" {
		if reasoned || !documented {
			return fmt.Errorf("present %s needs documents and null reason", name)
		}
		return nil
	}
	if !reasoned || documented || state == "not_applicable" && name != "technical" {
		return fmt.Errorf("invalid %s part state", name)
	}
	return nil
}

func validateSupportingDocuments(m map[string]canonicaljson.Value, docs contentDocuments) error {
	for _, v := range contentArray(m, "supporting") {
		if err := useSupportingDocument(contentString(contentFields(v), "document_id"), docs); err != nil {
			return err
		}
	}
	return nil
}

// A supporting document is declared, and neither a part nor an earlier
// supporting entry refers to it.
func useSupportingDocument(id string, docs contentDocuments) error {
	if docs.used[id] {
		return fmt.Errorf("supporting document duplicates a part")
	}
	return docs.use(id)
}

// specRequirementVerifiers returns the verification IDs of each requirement,
// by the requirement's ID; a requirement that is listed twice keeps the last.
// The requirements are validated one at a time.
func specRequirementVerifiers(m map[string]canonicaljson.Value, partRefs map[string]map[string]bool, technicalApplies bool) (map[string]map[string]bool, error) {
	verifiers := map[string]map[string]bool{}
	for _, v := range contentArray(m, "requirements") {
		r := contentFields(v)
		ids, err := requirementVerifiers(r, partRefs, technicalApplies)
		if err != nil {
			return nil, err
		}
		verifiers[contentString(r, "id")] = ids
	}
	return verifiers, nil
}

// requirementVerifiers returns the verification IDs of the requirement r. A
// requirement refers into the functional part and, unless that part is not
// applicable, into the technical part, and has a verification ID.
func requirementVerifiers(r map[string]canonicaljson.Value, partRefs map[string]map[string]bool, technicalApplies bool) (map[string]bool, error) {
	if err := validateRequirementRefs(r, "functional", true, partRefs["functional"]); err != nil {
		return nil, err
	}
	if err := validateRequirementRefs(r, "technical", technicalApplies, partRefs["technical"]); err != nil {
		return nil, err
	}
	ids := contentArray(r, "verification_ids")
	if len(ids) == 0 {
		return nil, fmt.Errorf("requirement needs a verification ID")
	}
	return contentStringSet(ids), nil
}

// The references of a requirement into the part named part are among those
// that the part holds, and there is at least one if the part requires it.
func validateRequirementRefs(r map[string]canonicaljson.Value, part string, required bool, held map[string]bool) error {
	refs := contentArray(r, part+"_refs")
	if len(refs) == 0 && required {
		return fmt.Errorf("requirement needs %s refs", part)
	}
	for _, ref := range refs {
		if !held[contentPartKey(ref)] {
			return fmt.Errorf("requirement reference is outside its part")
		}
	}
	return nil
}

// A phase covers each requirement it names with its own verification IDs.
func validatePhaseCoverage(m map[string]canonicaljson.Value, verifiers map[string]map[string]bool) error {
	for _, v := range contentArray(m, "phases") {
		p := contentFields(v)
		covered := contentStringSet(contentArray(p, "verification_ids"))
		for _, id := range contentArray(p, "requirement_ids") {
			if err := validateCoveredRequirement(id.(string), verifiers, covered); err != nil {
				return err
			}
		}
	}
	return nil
}

// The requirement id is known and each of its verification IDs is covered.
// They are visited in map order, so which of several uncovered ones the error
// names varies between runs.
func validateCoveredRequirement(id string, verifiers map[string]map[string]bool, covered map[string]bool) error {
	required, known := verifiers[id]
	if !known {
		return fmt.Errorf("phase refers to an unknown requirement")
	}
	for verifier := range required {
		if !covered[verifier] {
			return fmt.Errorf("phase does not cover requirement verifier %s", verifier)
		}
	}
	return nil
}

// An assessment holds the six checks in their fixed order. A ready assessment
// has passed them all and has no open questions.
func validateAssessmentSemantics(m map[string]canonicaljson.Value, docs contentDocuments) error {
	checks := contentArray(m, "checks")
	expected := []string{"acceptance_coverage", "implementation_basis", "problem_coverage", "recovery", "scope_and_phases", "three_parts"}
	if len(checks) != len(expected) {
		return fmt.Errorf("assessment must contain exactly six checks")
	}
	ready := contentString(m, "outcome") == "ready"
	for i, v := range checks {
		if err := validateAssessmentCheck(contentFields(v), expected[i], ready, docs); err != nil {
			return err
		}
	}
	if ready && len(contentArray(m, "open_questions")) != 0 {
		return fmt.Errorf("ready assessment cannot have open questions")
	}
	return nil
}

// The check c is the one named id, has passed if the assessment is ready, and
// its evidence is declared documents.
func validateAssessmentCheck(c map[string]canonicaljson.Value, id string, ready bool, docs contentDocuments) error {
	if contentString(c, "id") != id {
		return fmt.Errorf("assessment check must be %s", id)
	}
	if ready && contentString(c, "outcome") != "pass" {
		return fmt.Errorf("ready assessment requires all checks to pass")
	}
	return docs.useAll(contentArray(c, "evidence_document_ids"))
}

func taskSpecStructurallyReady(spec canonicaljson.Object) error {
	if isTaskGoalSpec(spec) {
		return contentError("task_goal_is_document", "a Goal document does not declare a selected technical solution", nil)
	}
	m := contentFields(spec)
	p := contentFields(m["parts"])
	for _, name := range []string{"abstract", "functional", "technical"} {
		state := contentString(contentFields(p[name]), "state")
		if state != "present" && !(name == "technical" && state == "not_applicable") {
			return contentError("task_spec_not_ready", name+" is unresolved", nil)
		}
	}
	reqs := contentArray(m, "requirements")
	phases := contentArray(m, "phases")
	if len(reqs) == 0 || len(phases) == 0 {
		return contentError("task_spec_not_ready", "ready requires at least one requirement and phase", nil)
	}
	covered := map[string]bool{}
	for _, v := range phases {
		for _, id := range contentArray(contentFields(v), "requirement_ids") {
			covered[id.(string)] = true
		}
	}
	for _, r := range reqs {
		if !covered[contentIDKey(r)] {
			return contentError("task_spec_not_ready", "a requirement is not covered by any phase", nil)
		}
	}
	return nil
}
func taskDocumentBytesValid(b []byte) error {
	if len(b) < 1 || len(b) > taskDocumentLimit || !utf8.Valid(b) || strings.ContainsRune(string(b), 0) || strings.HasPrefix(string(b), "\ufeff") {
		return contentError("task_content_invalid_input", "document must be 1 byte..4 MiB of UTF-8 without NUL or BOM", nil)
	}
	return nil
}
func sortedContentStrings(values []string) []canonicaljson.Value {
	sort.Strings(values)
	out := []canonicaljson.Value{}
	for _, s := range values {
		out = append(out, s)
	}
	return out
}

func sortTaskContentRegistry(r *WorkItemRegistry) {
	sort.Slice(r.TaskProblemRevisions, func(i, j int) bool {
		a, b := r.TaskProblemRevisions[i], r.TaskProblemRevisions[j]
		if a.TaskID != b.TaskID {
			return a.TaskID < b.TaskID
		}
		return a.Revision < b.Revision
	})
	sort.Slice(r.TaskSpecRevisions, func(i, j int) bool {
		a, b := r.TaskSpecRevisions[i], r.TaskSpecRevisions[j]
		if a.TaskID != b.TaskID {
			return a.TaskID < b.TaskID
		}
		if a.SpecID != b.SpecID {
			return a.SpecID < b.SpecID
		}
		return a.Revision < b.Revision
	})
	sort.Slice(r.TaskSpecAssessments, func(i, j int) bool {
		a, b := r.TaskSpecAssessments[i], r.TaskSpecAssessments[j]
		if a.TaskID != b.TaskID {
			return a.TaskID < b.TaskID
		}
		if a.SpecID != b.SpecID {
			return a.SpecID < b.SpecID
		}
		if a.SpecRevision != b.SpecRevision {
			return a.SpecRevision < b.SpecRevision
		}
		return a.Ordinal < b.Ordinal
	})
	sort.Slice(r.TaskSolutionSelections, func(i, j int) bool {
		a, b := r.TaskSolutionSelections[i], r.TaskSolutionSelections[j]
		if a.TaskID != b.TaskID {
			return a.TaskID < b.TaskID
		}
		return a.Ordinal < b.Ordinal
	})
	sort.Slice(r.TaskContentPublications, func(i, j int) bool {
		return r.TaskContentPublications[i].PublicationKey < r.TaskContentPublications[j].PublicationKey
	})
}
func validateTaskContentRegistry(r WorkItemRegistry) error {
	if r.TaskProblemRevisions == nil || r.TaskSpecRevisions == nil || r.TaskSpecAssessments == nil || r.TaskSolutionSelections == nil || r.TaskResultSpecBindings == nil || r.TaskContentPublications == nil {
		return fmt.Errorf("Task content arrays must be explicit")
	}
	tasks := map[TaskID]TaskRecord{}
	for _, t := range r.Tasks {
		tasks[t.ID] = t
	}
	outcomes := map[string]TaskContentOutcomeRef{}
	problemHeads := map[TaskID]int{}
	specHeads := map[TaskID]int{}
	specIDs := map[TaskID]string{}
	assessmentHeads := map[string]int{}
	selectionHeads := map[TaskID]int{}
	decisionIDs := map[string]bool{}
	add := func(o TaskContentOutcomeRef) error {
		if tasks[o.TaskID].ID == "" || !digestPattern.MatchString(o.ManifestSHA256) {
			return fmt.Errorf("invalid content reference")
		}
		if _, ok := outcomes[o.ManifestSHA256]; ok {
			return fmt.Errorf("duplicate manifest reference")
		}
		outcomes[o.ManifestSHA256] = o
		return nil
	}
	last := ""
	for _, p := range r.TaskProblemRevisions {
		key := fmt.Sprintf("%s\x00%010d", p.TaskID, p.Revision)
		if key <= last || p.Revision != problemHeads[p.TaskID]+1 || p.Revision > 2147483647 {
			return fmt.Errorf("invalid problem revision chain")
		}
		last = key
		problemHeads[p.TaskID] = p.Revision
		n := p.Revision
		if e := add(TaskContentOutcomeRef{Kind: "problem", TaskID: p.TaskID, Revision: &n, ManifestSHA256: p.ManifestSHA256}); e != nil {
			return e
		}
	}
	last = ""
	for _, p := range r.TaskSpecRevisions {
		key := fmt.Sprintf("%s\x00%s\x00%010d", p.TaskID, p.SpecID, p.Revision)
		if key <= last || p.Revision != specHeads[p.TaskID]+1 || p.Revision > 2147483647 || contentSlug(p.SpecID) != nil || specIDs[p.TaskID] != "" && specIDs[p.TaskID] != p.SpecID || problemHeads[p.TaskID] == 0 {
			return fmt.Errorf("invalid Spec revision chain")
		}
		last = key
		specHeads[p.TaskID] = p.Revision
		specIDs[p.TaskID] = p.SpecID
		n, id := p.Revision, p.SpecID
		if e := add(TaskContentOutcomeRef{Kind: "spec", TaskID: p.TaskID, SpecID: &id, Revision: &n, ManifestSHA256: p.ManifestSHA256}); e != nil {
			return e
		}
	}
	last = ""
	for _, p := range r.TaskSpecAssessments {
		group := fmt.Sprintf("%s\x00%s\x00%010d", p.TaskID, p.SpecID, p.SpecRevision)
		key := fmt.Sprintf("%s\x00%010d", group, p.Ordinal)
		if key <= last || p.Ordinal != assessmentHeads[group]+1 || p.Ordinal > 2147483647 || p.SpecRevision < 1 || p.SpecRevision > specHeads[p.TaskID] || p.SpecID != specIDs[p.TaskID] || decisionIDs[p.ID] || !regexp.MustCompile(`^asm_[0-9a-f]{32}$`).MatchString(p.ID) {
			return fmt.Errorf("invalid assessment chain")
		}
		last = key
		assessmentHeads[group] = p.Ordinal
		decisionIDs[p.ID] = true
		n, s, id := p.SpecRevision, p.SpecID, p.ID
		if e := add(TaskContentOutcomeRef{Kind: "assessment", TaskID: p.TaskID, SpecID: &s, Revision: &n, ID: &id, ManifestSHA256: p.ManifestSHA256}); e != nil {
			return e
		}
	}
	last = ""
	for _, p := range r.TaskSolutionSelections {
		key := fmt.Sprintf("%s\x00%010d", p.TaskID, p.Ordinal)
		if key <= last || p.Ordinal != selectionHeads[p.TaskID]+1 || p.Ordinal > 2147483647 || decisionIDs[p.ID] || !regexp.MustCompile(`^sel_[0-9a-f]{32}$`).MatchString(p.ID) {
			return fmt.Errorf("invalid selection chain")
		}
		last = key
		selectionHeads[p.TaskID] = p.Ordinal
		decisionIDs[p.ID] = true
		id := p.ID
		if e := add(TaskContentOutcomeRef{Kind: "selection", TaskID: p.TaskID, ID: &id, ManifestSHA256: p.ManifestSHA256}); e != nil {
			return e
		}
	}
	pubs := map[string]TaskContentPublication{}
	seenOutcomes := map[string]bool{}
	last = ""
	for _, p := range r.TaskContentPublications {
		if p.PublicationKey <= last || contentKey(p.PublicationKey) != nil || !digestPattern.MatchString(p.IntentSHA256) || !digestPattern.MatchString(p.RequestSHA256) || contentUTC(p.RecordedAtUTC) != nil {
			return fmt.Errorf("invalid content publication")
		}
		last = p.PublicationKey
		o, ok := outcomes[p.OutcomeRef.ManifestSHA256]
		if !ok || seenOutcomes[o.ManifestSHA256] || !contentTypedEqual(o, p.OutcomeRef) || o.TaskID != p.TaskID {
			return fmt.Errorf("publication outcome does not resolve uniquely")
		}
		seenOutcomes[o.ManifestSHA256] = true
		validOp := map[string]string{"task_create": "problem", "problem_record": "problem", "spec_record": "spec", "spec_assess": "assessment", "spec_select": "selection", "spec_withdraw": "selection"}
		if validOp[p.Operation] != o.Kind {
			return fmt.Errorf("publication operation differs from outcome")
		}
		if p.RegistryBeforeSHA256 != nil && !digestPattern.MatchString(*p.RegistryBeforeSHA256) {
			return fmt.Errorf("invalid prior registry digest")
		}
		if p.StoreTransition != "none" || p.BackupSHA256 != nil {
			return fmt.Errorf("publication cannot carry a registry migration")
		}
		pubs[p.PublicationKey] = p
	}
	if len(seenOutcomes) != len(outcomes) {
		return fmt.Errorf("content record is missing its publication")
	}
	return nil
}
