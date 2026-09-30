package workspace

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

func contentExact(fields map[string]contentRule) contentRule {
	return func(v canonicaljson.Value) error {
		o, ok := v.(canonicaljson.Object)
		if !ok || len(o) != len(fields) {
			return fmt.Errorf("object has missing or unknown fields")
		}
		seen := map[string]bool{}
		for _, m := range o {
			rule, ok := fields[m.Name]
			if !ok || seen[m.Name] {
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
			return contentExact(map[string]contentRule{"kind": contentEnum("snapshot"), "manifest_sha256": contentDigest, "document_id": contentKey})(v)
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
		s["registry_upgrade"] = contentNullable(contentExact(map[string]contentRule{"from_version": contentInteger(1, 2), "registry_sha256": contentDigest}))
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
	if err = contentExact(taskContentSchema(operation, published, internal))(v); err == nil {
		err = validateTaskContentSemantics(contentFields(v), operation, published, internal)
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

func validateTaskContentSemantics(m map[string]canonicaljson.Value, operation string, published, internal bool) error {
	docs := map[string]bool{}
	used := map[string]bool{}
	for _, d := range contentArray(m, "documents") {
		docs[contentIDKey(d)] = true
	}
	use := func(id string) error {
		if !docs[id] {
			return fmt.Errorf("document %s is missing", id)
		}
		used[id] = true
		return nil
	}
	switch operation {
	case "problem_record":
		origin := contentString(contentFields(m["origin"]), "kind")
		generated := origin == "legacy_summary_import" || origin == "task_create_summary"
		if generated && !published {
			if len(docs) != 0 || len(contentArray(m, "sources")) != 0 || len(contentArray(m, "claims")) != 0 || m["deadline"] != nil || contentString(m, "problem_document_id") != "problem" {
				return fmt.Errorf("generated problem requires empty documents, sources and claims, null deadline and problem ID")
			}
		} else {
			if e := use(contentString(m, "problem_document_id")); e != nil {
				return e
			}
		}
		sources := map[string]bool{}
		for _, v := range contentArray(m, "sources") {
			s := contentFields(v)
			sources[contentString(s, "id")] = true
			if s["document_id"] != nil {
				if s["withholding_reason"] != nil {
					return fmt.Errorf("withheld source cannot have a document")
				}
				if e := use(contentString(s, "document_id")); e != nil {
					return e
				}
			}
		}
		claims := map[string]bool{}
		for _, v := range contentArray(m, "claims") {
			c := contentFields(v)
			claims[contentString(c, "id")] = true
			if !sources[contentString(c, "source_id")] {
				return fmt.Errorf("claim source does not exist")
			}
			evidence := contentArray(c, "evidence_document_ids")
			if contentString(c, "verification") == "unverified" {
				if c["reason"] != nil || len(evidence) != 0 {
					return fmt.Errorf("unverified claim cannot have verification evidence")
				}
			} else if c["reason"] == nil || len(evidence) == 0 {
				return fmt.Errorf("verified claim requires reason and evidence")
			}
			for _, e := range evidence {
				if err := use(e.(string)); err != nil {
					return err
				}
			}
		}
		if m["deadline"] != nil && !claims[contentString(contentFields(m["deadline"]), "claim_id")] {
			return fmt.Errorf("deadline claim does not exist")
		}
	case "spec_record":
		parts := contentFields(m["parts"])
		sets := map[string]map[string]bool{}
		for _, name := range []string{"abstract", "functional", "technical"} {
			p := contentFields(parts[name])
			state := contentString(p, "state")
			refs := contentArray(p, "documents")
			if state == "present" {
				if p["reason"] != nil || len(refs) == 0 {
					return fmt.Errorf("present %s needs documents and null reason", name)
				}
			} else if p["reason"] == nil || len(refs) != 0 || state == "not_applicable" && name != "technical" {
				return fmt.Errorf("invalid %s part state", name)
			}
			sets[name] = map[string]bool{}
			for _, v := range refs {
				if e := use(contentString(contentFields(v), "document_id")); e != nil {
					return e
				}
				sets[name][contentPartKey(v)] = true
			}
		}
		for _, v := range contentArray(m, "supporting") {
			id := contentString(contentFields(v), "document_id")
			if used[id] {
				return fmt.Errorf("supporting document duplicates a part")
			}
			if e := use(id); e != nil {
				return e
			}
		}
		reqs := map[string]map[string]bool{}
		for _, v := range contentArray(m, "requirements") {
			r := contentFields(v)
			for _, name := range []string{"functional", "technical"} {
				refs := contentArray(r, name+"_refs")
				if len(refs) == 0 && (name == "functional" || contentString(contentFields(parts["technical"]), "state") != "not_applicable") {
					return fmt.Errorf("requirement needs %s refs", name)
				}
				for _, ref := range refs {
					if !sets[name][contentPartKey(ref)] {
						return fmt.Errorf("requirement reference is outside its part")
					}
				}
			}
			vs := contentArray(r, "verification_ids")
			if len(vs) == 0 {
				return fmt.Errorf("requirement needs a verification ID")
			}
			reqs[contentString(r, "id")] = map[string]bool{}
			for _, id := range vs {
				reqs[contentString(r, "id")][id.(string)] = true
			}
		}
		for _, v := range contentArray(m, "phases") {
			p := contentFields(v)
			vs := map[string]bool{}
			for _, id := range contentArray(p, "verification_ids") {
				vs[id.(string)] = true
			}
			for _, id := range contentArray(p, "requirement_ids") {
				required, ok := reqs[id.(string)]
				if !ok {
					return fmt.Errorf("phase refers to an unknown requirement")
				}
				for vid := range required {
					if !vs[vid] {
						return fmt.Errorf("phase does not cover requirement verifier %s", vid)
					}
				}
			}
		}
	case "spec_assess":
		checks := contentArray(m, "checks")
		expected := []string{"acceptance_coverage", "implementation_basis", "problem_coverage", "recovery", "scope_and_phases", "three_parts"}
		if len(checks) != len(expected) {
			return fmt.Errorf("assessment must contain exactly six checks")
		}
		for i, v := range checks {
			c := contentFields(v)
			if contentString(c, "id") != expected[i] {
				return fmt.Errorf("assessment check must be %s", expected[i])
			}
			if contentString(m, "outcome") == "ready" && contentString(c, "outcome") != "pass" {
				return fmt.Errorf("ready assessment requires all checks to pass")
			}
			for _, id := range contentArray(c, "evidence_document_ids") {
				if e := use(id.(string)); e != nil {
					return e
				}
			}
		}
		if contentString(m, "outcome") == "ready" && len(contentArray(m, "open_questions")) != 0 {
			return fmt.Errorf("ready assessment cannot have open questions")
		}
	case "spec_select":
		if m["solution"] == nil {
			return fmt.Errorf("select requires an exact solution")
		}
	case "spec_withdraw":
		if m["solution"] != nil {
			return fmt.Errorf("withdraw requires null solution")
		}
	}
	for id := range docs {
		if !used[id] {
			return fmt.Errorf("unused document %s", id)
		}
	}
	return nil
}

func taskSpecStructurallyReady(spec canonicaljson.Object) error {
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
	sort.Slice(r.TaskSpecPolicies, func(i, j int) bool { return r.TaskSpecPolicies[i].TaskID < r.TaskSpecPolicies[j].TaskID })
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
	sort.Slice(r.TaskResultSpecBindings, func(i, j int) bool {
		return r.TaskResultSpecBindings[i].TaskResultID < r.TaskResultSpecBindings[j].TaskResultID
	})
	sort.Slice(r.TaskContentPublications, func(i, j int) bool {
		return r.TaskContentPublications[i].PublicationKey < r.TaskContentPublications[j].PublicationKey
	})
}
func validateTaskContentRegistry(r WorkItemRegistry) error {
	if r.FormatVersion < 3 {
		if r.TaskSpecPolicies != nil || r.TaskProblemRevisions != nil || r.TaskSpecRevisions != nil || r.TaskSpecAssessments != nil || r.TaskSolutionSelections != nil || r.TaskResultSpecBindings != nil || r.TaskContentPublications != nil {
			return fmt.Errorf("legacy registry cannot contain Task content collections")
		}
		return nil
	}
	if r.TaskSpecPolicies == nil || r.TaskProblemRevisions == nil || r.TaskSpecRevisions == nil || r.TaskSpecAssessments == nil || r.TaskSolutionSelections == nil || r.TaskResultSpecBindings == nil || r.TaskContentPublications == nil {
		return fmt.Errorf("format 3 requires explicit Task content arrays")
	}
	tasks := map[TaskID]TaskRecord{}
	for _, t := range r.Tasks {
		tasks[t.ID] = t
	}
	policies := map[TaskID]TaskSpecPolicy{}
	last := ""
	for _, p := range r.TaskSpecPolicies {
		if string(p.TaskID) <= last || tasks[p.TaskID].ID == "" || (p.Mode != "legacy" && p.Mode != "spec_required") || p.LegacyResultIDs == nil {
			return fmt.Errorf("invalid Task policy")
		}
		last = string(p.TaskID)
		if p.Mode == "legacy" && (p.ActivationPublicationKey != nil || len(p.LegacyResultIDs) != 0) || p.Mode == "spec_required" && p.ActivationPublicationKey == nil {
			return fmt.Errorf("invalid Task policy activation")
		}
		previous := ""
		for _, id := range p.LegacyResultIDs {
			if string(id) <= previous {
				return fmt.Errorf("legacy result IDs must be sorted unique")
			}
			found := false
			for _, v := range r.TaskResults {
				if v.ID == id && v.TaskID == p.TaskID {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("dangling legacy result ID")
			}
			previous = string(id)
		}
		policies[p.TaskID] = p
	}
	if len(policies) != len(tasks) {
		return fmt.Errorf("format 3 requires one policy per Task")
	}
	outcomes := map[string]TaskContentOutcomeRef{}
	problemHeads := map[TaskID]int{}
	specHeads := map[TaskID]int{}
	specIDs := map[TaskID]string{}
	assessmentHeads := map[string]int{}
	selectionHeads := map[TaskID]int{}
	decisionIDs := map[string]bool{}
	add := func(o TaskContentOutcomeRef) error {
		if tasks[o.TaskID].ID == "" || !digestPattern.MatchString(o.ManifestSHA256) || policies[o.TaskID].Mode != "spec_required" {
			return fmt.Errorf("invalid content reference")
		}
		if _, ok := outcomes[o.ManifestSHA256]; ok {
			return fmt.Errorf("duplicate manifest reference")
		}
		outcomes[o.ManifestSHA256] = o
		return nil
	}
	last = ""
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
		if p.StoreTransition == "none" {
			if p.BackupSHA256 != nil {
				return fmt.Errorf("non-migration publication cannot have backup")
			}
		} else if (p.StoreTransition != "format_1_to_3" && p.StoreTransition != "format_2_to_3") || p.BackupSHA256 == nil || p.RegistryBeforeSHA256 == nil || *p.BackupSHA256 != *p.RegistryBeforeSHA256 || !digestPattern.MatchString(*p.BackupSHA256) {
			return fmt.Errorf("invalid migration backup binding")
		}
		pubs[p.PublicationKey] = p
	}
	if len(seenOutcomes) != len(outcomes) {
		return fmt.Errorf("content record is missing its publication")
	}
	for _, p := range policies {
		if p.Mode == "spec_required" {
			pub, ok := pubs[*p.ActivationPublicationKey]
			if !ok || pub.TaskID != p.TaskID || pub.OutcomeRef.Kind != "problem" || pub.OutcomeRef.Revision == nil || *pub.OutcomeRef.Revision != 1 {
				return fmt.Errorf("policy does not bind initial problem publication")
			}
		}
	}
	links := map[TaskResultID]bool{}
	last = ""
	for _, link := range r.TaskResultSpecBindings {
		if string(link.TaskResultID) <= last || links[link.TaskResultID] || !digestPattern.MatchString(link.HandoffSHA256) || !digestPattern.MatchString(link.StartReceiptSHA256) {
			return fmt.Errorf("invalid result Spec link")
		}
		last = string(link.TaskResultID)
		var result *TaskResultRecord
		for _, v := range r.TaskResults {
			if v.ID == link.TaskResultID {
				c := v
				result = &c
			}
		}
		if result == nil || result.TaskID != link.TaskID || result.HandoffSHA256 != link.HandoffSHA256 || result.StartReceiptSHA256 != link.StartReceiptSHA256 || validateTaskSpecBasisRegistry(r, link.Basis) != nil || link.Basis.TaskID != link.TaskID {
			return fmt.Errorf("result Spec link does not match historical basis")
		}
		links[link.TaskResultID] = true
	}
	for _, v := range r.TaskResults {
		p := policies[v.TaskID]
		if p.Mode != "spec_required" || links[v.ID] {
			continue
		}
		legacy := false
		for _, id := range p.LegacyResultIDs {
			if id == v.ID {
				legacy = true
			}
		}
		if !legacy {
			return fmt.Errorf("required Task result is missing its atomic Spec link")
		}
	}
	return nil
}
