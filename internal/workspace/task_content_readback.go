package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

type TaskContentReadbackResult struct{ Value canonicaljson.Object }

func MarshalTaskContentReadback(r TaskContentReadbackResult) ([]byte, error) {
	return canonicaljson.Marshal(r.Value)
}
func MarshalTaskContentMutation(r TaskContentMutationResult) ([]byte, error) {
	v, e := contentValue(r)
	if e != nil {
		return nil, e
	}
	m := contentFields(v)
	m["next_transition_authorized"] = false
	return canonicaljson.Marshal(contentEnvelope("WorkspaceTaskContentMutationReadback@1", m))
}
func contentIntegrity(err error) string {
	if err == nil {
		return "valid"
	}
	if errors.Is(err, fs.ErrNotExist) {
		return "missing"
	}
	if strings.Contains(err.Error(), "conflict") || strings.Contains(err.Error(), "invalid") {
		return "conflict"
	}
	return "unknown"
}

func contentIntegrityReason(integrity string) string {
	switch integrity {
	case "missing":
		return "task_content_missing"
	case "conflict":
		return "task_content_integrity_conflict"
	default:
		return "task_content_observation_unknown"
	}
}
func readContentStart(d Dependencies, q TaskContentQuery, missingTask bool) (string, WorkItemRegistry, *TaskRecord, error) {
	if _, e := ParseTaskID(string(q.TaskID)); e != nil {
		return "", WorkItemRegistry{}, nil, e
	}
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return "", WorkItemRegistry{}, nil, e
	}
	if d.WorkItems == nil {
		return root, WorkItemRegistry{}, nil, contentError("task_content_observation_unknown", "registry reader unavailable", nil)
	}
	r, e := d.WorkItems.Snapshot(root)
	if e != nil {
		return root, r, nil, e
	}
	t, _ := findTask(r, q.TaskID)
	if t == nil && !missingTask {
		return root, r, nil, workError(ErrorWorkNotFound, "Task is not registered", nil)
	}
	return root, r, t, nil
}
func contentReadFreshness(d Dependencies, root string, r WorkItemRegistry) string {
	after, e := d.WorkItems.Snapshot(root)
	if e != nil || after.RawSHA256 != r.RawSHA256 {
		return "unknown"
	}
	return "fresh"
}
func contentReadEnvelope(kind string, m map[string]canonicaljson.Value) TaskContentReadbackResult {
	m["next_transition_authorized"] = false
	return TaskContentReadbackResult{contentEnvelope(kind, m)}
}
func contentActionValue(action TaskNextAction) canonicaljson.Value {
	v, _ := contentValue(action)
	return v
}
func contentRefValue(value any) canonicaljson.Value { v, _ := contentValue(value); return v }
func contentTaskAction(id TaskID) canonicaljson.Value {
	return contentActionValue(TaskNextAction{"inspect_task", "Inspect the Task and its current basis.", []string{"ply", "workspace", "task", "show", string(id)}})
}
func ShowTaskProblem(d Dependencies, q TaskContentQuery) (TaskContentReadbackResult, error) {
	if q.Revision < 0 || q.Revision > 2147483647 {
		return TaskContentReadbackResult{}, contentError("task_content_invalid_input", "revision must be a positive integer", nil)
	}
	root, r, t, e := readContentStart(d, q, false)
	if e != nil {
		return TaskContentReadbackResult{}, e
	}
	var selected *TaskProblemReference
	for _, p := range r.TaskProblemRevisions {
		if p.TaskID == q.TaskID && (q.Revision == 0 || p.Revision == q.Revision) {
			v := p
			selected = &v
		}
	}
	var revision, legacy canonicaljson.Value
	integrity := "not_recorded"
	reasons := []string{}
	if selected == nil {
		if q.Revision != 0 {
			return TaskContentReadbackResult{}, workError(ErrorWorkNotFound, "problem revision is not registered", nil)
		}
		legacy = contentObject(map[string]canonicaljson.Value{"title": t.Title, "description": t.Description})
	} else {
		o, err := readRegisteredTaskManifest(d.TaskContent, root, r, selected.ManifestSHA256)
		integrity = contentIntegrity(err)
		if err == nil {
			revision = o
			err = verifyManifestDocuments(d.TaskContent, root, o)
			integrity = contentIntegrity(err)
		}
		if err != nil {
			reasons = append(reasons, contentIntegrityReason(integrity))
		}
	}
	return contentReadEnvelope("WorkspaceTaskProblemReadback@1", map[string]canonicaljson.Value{"workspace": root, "task_id": string(q.TaskID), "registry_sha256": nullableDigest(r.RawSHA256), "revision": revision, "legacy_problem_summary": legacy, "integrity": integrity, "freshness": contentReadFreshness(d, root, r), "reasons": sortedContentStrings(reasons), "next_action": contentTaskAction(q.TaskID)}), nil
}
func nullableDigest(s string) canonicaljson.Value {
	if s == "" {
		return nil
	}
	return s
}
func verifyManifestDocuments(s *TaskContentStorage, root string, o canonicaljson.Object) error {
	for _, v := range contentArray(contentFields(o), "documents") {
		m := contentFields(v)
		digest := contentString(m, "sha256")
		b, e := s.Read(root, "objects", digest)
		if e != nil {
			return e
		}
		if len(b) != contentInt(m, "size_bytes") || contentString(m, "locator") != taskContentPath(root, "objects", digest) {
			return contentError("task_content_integrity_conflict", "snapshot descriptor differs", nil)
		}
	}
	return nil
}
func selectionReadback(d Dependencies, root string, r WorkItemRegistry, id TaskID) (canonicaljson.Value, string) {
	state := taskContentState(r, id)
	m := map[string]canonicaljson.Value{"event": contentRefValue(state.SelectionEvent), "action": nil, "solution": nil, "freshness": "none"}
	if state.SelectionEvent == nil {
		return contentObject(m), "none"
	}
	o, e := readRegisteredTaskManifest(d.TaskContent, root, r, state.SelectionEvent.ManifestSHA256)
	if e != nil {
		m["freshness"] = "unknown"
		return contentObject(m), "unknown"
	}
	v := contentFields(o)
	m["action"] = v["action"]
	m["solution"] = v["solution"]
	fresh := "current"
	if contentString(v, "action") == "withdraw" {
		fresh = "withdrawn"
	} else {
		sol := contentFields(v["solution"])
		spec := valueRevision(sol["spec"])
		asm := valueDecision(sol["assessment"])
		last := lastTaskAssessment(r, id, contentString(sol, "spec_id"), spec.Revision)
		if !contentTypedEqual(state.ProblemHead, valueRevision(sol["problem"])) || last == nil || last.ID != asm.ID || last.ManifestSHA256 != asm.ManifestSHA256 {
			fresh = "stale"
		}
	}
	m["freshness"] = fresh
	return contentObject(m), fresh
}
func ListTaskSpecs(d Dependencies, q TaskContentQuery) (TaskContentReadbackResult, error) {
	root, r, _, e := readContentStart(d, q, false)
	if e != nil {
		return TaskContentReadbackResult{}, e
	}
	selected, fresh := selectionReadback(d, root, r, q.TaskID)
	sol := contentFields(contentFields(selected)["solution"])
	selectedRef := valueRevision(sol["spec"])
	revisions := []canonicaljson.Value{}
	reasons := []string{}
	var specID canonicaljson.Value
	for _, p := range r.TaskSpecRevisions {
		if p.TaskID != q.TaskID {
			continue
		}
		specID = p.SpecID
		o, e := readRegisteredTaskManifest(d.TaskContent, root, r, p.ManifestSHA256)
		m := contentFields(o)
		if e != nil {
			reasons = append(reasons, contentIntegrityReason(contentIntegrity(e)))
		}
		var assessment canonicaljson.Value
		if last := lastTaskAssessment(r, q.TaskID, p.SpecID, p.Revision); last != nil {
			a, e := readRegisteredTaskManifest(d.TaskContent, root, r, last.ManifestSHA256)
			if e == nil {
				assessment = contentObject(map[string]canonicaljson.Value{"id": last.ID, "manifest_sha256": last.ManifestSHA256, "outcome": contentFields(a)["outcome"]})
			} else {
				reasons = append(reasons, contentIntegrityReason(contentIntegrity(e)))
			}
		}
		revisions = append(revisions, contentObject(map[string]canonicaljson.Value{"revision": int64(p.Revision), "manifest_sha256": p.ManifestSHA256, "title": m["title"], "problem": m["problem"], "assessment": assessment, "selected": fresh != "withdrawn" && selectedRef != nil && selectedRef.Revision == p.Revision && selectedRef.ManifestSHA256 == p.ManifestSHA256}))
	}
	return contentReadEnvelope("WorkspaceTaskSpecListReadback@1", map[string]canonicaljson.Value{"workspace": root, "task_id": string(q.TaskID), "registry_sha256": nullableDigest(r.RawSHA256), "spec_id": specID, "revisions": revisions, "selection": selected, "freshness": contentReadFreshness(d, root, r), "reasons": sortedContentStrings(sortedReasons(reasons)), "next_action": contentTaskAction(q.TaskID)}), nil
}
func ShowTaskSpec(d Dependencies, q TaskContentQuery) (TaskContentReadbackResult, error) {
	if contentSlug(q.SpecID) != nil || q.Revision < 1 || q.Revision > 2147483647 {
		return TaskContentReadbackResult{}, contentError("task_content_invalid_input", "--spec and --revision (1..2147483647) are required", nil)
	}
	root, r, t, e := readContentStart(d, q, false)
	if e != nil {
		return TaskContentReadbackResult{}, e
	}
	var ref *TaskRevisionRef
	for _, p := range r.TaskSpecRevisions {
		if p.TaskID == q.TaskID && p.SpecID == q.SpecID && p.Revision == q.Revision {
			ref = &TaskRevisionRef{p.Revision, p.ManifestSHA256}
		}
	}
	if ref == nil {
		return TaskContentReadbackResult{}, workError(ErrorWorkNotFound, "Spec revision is not registered", nil)
	}
	o, err := readRegisteredTaskManifest(d.TaskContent, root, r, ref.ManifestSHA256)
	integrity := contentIntegrity(err)
	reasons := []string{}
	if err == nil {
		err = verifyManifestDocuments(d.TaskContent, root, o)
		integrity = contentIntegrity(err)
	}
	if err != nil {
		reasons = append(reasons, contentIntegrityReason(integrity))
	}
	assessments := []canonicaljson.Value{}
	selections := []canonicaljson.Value{}
	for _, p := range r.TaskSpecAssessments {
		if p.TaskID == q.TaskID && p.SpecID == q.SpecID && p.SpecRevision == q.Revision {
			v, e := readRegisteredTaskManifest(d.TaskContent, root, r, p.ManifestSHA256)
			if e == nil {
				assessments = append(assessments, v)
			} else {
				integrity = contentIntegrity(e)
				reasons = append(reasons, contentIntegrityReason(integrity))
			}
		}
	}
	for _, p := range r.TaskSolutionSelections {
		if p.TaskID == q.TaskID {
			v, e := readRegisteredTaskManifest(d.TaskContent, root, r, p.ManifestSHA256)
			if e == nil {
				selections = append(selections, v)
			} else {
				integrity = contentIntegrity(e)
				reasons = append(reasons, contentIntegrityReason(integrity))
			}
		}
	}
	_, selectionFresh := selectionReadback(d, root, r, q.TaskID)
	targetFresh := "unknown"
	inputs := []TaskSpecRequiredInput{}
	var basis *TaskSpecBasis
	projects, repos, e := d.Projects.Snapshot(root)
	if e == nil {
		eval, e := currentTaskSpec(d, root, r, ProjectSnapshot{Projects: projects, Repos: repos}, *t, true)
		reasons = append(reasons, eval.Reasons...)
		if e != nil {
			selectionFresh = "unknown"
			reasons = append(reasons, contentIntegrityReason(contentIntegrity(e)))
		} else {
			selectionFresh = eval.SelectionFreshness
		}
		if e == nil && contentString(contentFields(eval.Spec), "spec_id") == q.SpecID && contentInt(contentFields(eval.Spec), "revision") == q.Revision {
			inputs = eval.RequiredInputs
			targetFresh = eval.TargetFreshness
			if eval.SelectionFreshness == "current" && eval.ContentIntegrity == "valid" && targetFresh == "fresh" {
				basis = eval.Basis
			}
		}
	} else {
		selectionFresh = "unknown"
		reasons = append(reasons, "task_content_observation_unknown")
	}
	if contentReadFreshness(d, root, r) != "fresh" {
		selectionFresh = "unknown"
		targetFresh = "unknown"
		basis = nil
		reasons = append(reasons, "task_content_observation_unknown")
	}
	var revision canonicaljson.Value
	if o != nil {
		revision = o
	}
	return contentReadEnvelope("WorkspaceTaskSpecReadback@1", map[string]canonicaljson.Value{"workspace": root, "task_id": string(q.TaskID), "registry_sha256": nullableDigest(r.RawSHA256), "revision": revision, "assessments": assessments, "selection_history": selections, "integrity": integrity, "selection_freshness": selectionFresh, "target_freshness": targetFresh, "required_inputs": contentRefValue(inputs), "task_spec_binding": contentRefValue(basis), "reasons": sortedContentStrings(sortedReasons(reasons)), "next_action": contentTaskAction(q.TaskID)}), nil
}
func ShowTaskPublication(d Dependencies, q TaskContentQuery) (TaskContentReadbackResult, error) {
	if e := contentKey(q.PublicationKey); e != nil {
		return TaskContentReadbackResult{}, e
	}
	root, r, _, err := readContentStart(d, q, true)
	if root == "" {
		return TaskContentReadbackResult{}, err
	}
	classification := "not_published"
	integrity := "not_recorded"
	var publication, outcome canonicaljson.Value
	if err != nil {
		classification = "unknown"
		integrity = "unknown"
	} else if p := taskPublication(r, q.PublicationKey); p != nil {
		if p.TaskID != q.TaskID {
			return TaskContentReadbackResult{}, contentError("task_content_publication_conflict", "publication key belongs to a different Task", nil)
		}
		publication = contentRefValue(p)
		o, e := readRegisteredTaskManifest(d.TaskContent, root, r, p.OutcomeRef.ManifestSHA256)
		if e == nil {
			outcome = o
			e = verifyManifestDocuments(d.TaskContent, root, o)
		}
		integrity = contentIntegrity(e)
		classification = "published"
		if e != nil {
			classification = "unknown"
		}
	}
	if err == nil && contentReadFreshness(d, root, r) != "fresh" {
		classification = "unknown"
	}
	return contentReadEnvelope("WorkspaceTaskPublicationReadback@1", map[string]canonicaljson.Value{"workspace": root, "task_id": string(q.TaskID), "publication_key": q.PublicationKey, "registry_sha256": nullableDigest(r.RawSHA256), "publication": publication, "outcome": outcome, "integrity": integrity, "classification": classification, "durability": "not_observed", "next_action": contentActionValue(contentInspectAction(q.TaskID, q.PublicationKey))}), nil
}

func buildTaskContentReadback(d Dependencies, r WorkItemRegistry, projects ProjectSnapshot, result TaskReadbackResult) canonicaljson.Object {
	root, t := result.Workspace, result.Task
	state := taskContentState(r, t.ID)
	problem := map[string]canonicaljson.Value{"head": contentRefValue(state.ProblemHead), "title": nil, "summary": nil, "legacy_problem_summary": nil, "manifest_locator": nil, "integrity": "not_recorded"}
	reasons := []string{}
	if state.ProblemHead == nil {
		problem["title"] = t.Title
		problem["summary"] = t.Description
		problem["legacy_problem_summary"] = contentObject(map[string]canonicaljson.Value{"title": t.Title, "description": t.Description})
	} else {
		problem["manifest_locator"] = taskContentPath(root, "manifests", state.ProblemHead.ManifestSHA256)
		o, e := readRegisteredTaskManifest(d.TaskContent, root, r, state.ProblemHead.ManifestSHA256)
		if e == nil {
			e = verifyManifestDocuments(d.TaskContent, root, o)
		}
		if e == nil {
			problem["title"] = contentFields(o)["title"]
			problem["summary"] = contentFields(o)["summary"]
		}
		problem["integrity"] = contentIntegrity(e)
		if e != nil {
			reasons = append(reasons, contentIntegrityReason(contentIntegrity(e)))
		}
	}
	sel, sfresh := selectionReadback(d, root, r, t.ID)
	sol := contentFields(contentFields(sel)["solution"])
	var specID canonicaljson.Value
	for _, p := range r.TaskSpecRevisions {
		if p.TaskID == t.ID {
			specID = p.SpecID
		}
	}
	selected := valueRevision(sol["spec"])
	newer := selected != nil && state.SpecHead != nil && state.SpecHead.Revision > selected.Revision
	eval, e := currentTaskSpec(d, root, r, projects, t, true)
	if e != nil {
		eval.ContentIntegrity = contentIntegrity(e)
		eval.TargetFreshness = "unknown"
		sfresh = "unknown"
		reasons = append(reasons, contentIntegrityReason(contentIntegrity(e)))
	} else {
		sfresh = eval.SelectionFreshness
	}
	if eval.ContentIntegrity == "" {
		eval.ContentIntegrity = "not_recorded"
	}
	reasons = append(reasons, eval.Reasons...)
	results := []canonicaljson.Value{}
	for _, v := range r.TaskResults {
		if v.TaskID != t.ID {
			continue
		}
		var link *TaskResultSpecBinding
		for _, b := range r.TaskResultSpecBindings {
			if b.TaskResultID == v.ID {
				copy := b
				link = &copy
			}
		}
		status, relevance := "not_recorded", "not_recorded"
		if link != nil {
			status = "bound"
			relevance = "stale"
			if _, e := historicalTaskSpec(d, root, r, link.Basis); e != nil {
				status = "invalid"
				relevance = "unknown"
			} else {
				cur, e := currentTaskSpec(d, root, r, projects, t, false)
				if e != nil {
					relevance = "unknown"
				} else if cur.SelectionFreshness == "current" && contentTypedEqual(cur.Basis, &link.Basis) {
					relevance = "current"
				}
			}
		}
		var locator canonicaljson.Value
		for _, a := range v.Artifacts {
			if a.ArtifactID == "task-requirements" {
				locator = a.Locator
			}
		}
		results = append(results, contentObject(map[string]canonicaljson.Value{"task_result_id": string(v.ID), "result_oid": v.ResultOID, "result_tree": v.ResultTree, "technical_gate": v.TechnicalGate, "spec_binding": contentRefValue(link), "basis_status": status, "relevance": relevance, "requirement_artifact_locator": locator}))
	}
	fresh := contentReadFreshness(d, root, r)
	ready := fresh == "fresh" && sfresh == "current" && eval.ContentIntegrity == "valid" && eval.TargetFreshness == "fresh" && result.WorktreeReady
	if fresh != "fresh" {
		sfresh = "unknown"
		reasons = append(reasons, "task_content_observation_unknown")
	}
	selection := contentFields(sel)
	selection["freshness"] = sfresh
	sel = contentObject(selection)
	action := TaskNextAction{"record_solution", "Record a solution revision with three parts.", []string{}}
	if state.ProblemHead == nil {
		action = TaskNextAction{"record_problem", "Record the Task problem explicitly.", []string{}}
	} else if state.SpecHead != nil {
		action = TaskNextAction{"record_human_choice", "Record an explicit human choice after a ready assessment.", []string{}}
	}
	if sfresh == "stale" || sfresh == "withdrawn" {
		action = TaskNextAction{"resolve_blocker", "Clarify the current solution basis.", []string{}}
	}
	if sfresh == "current" && !result.WorktreeReady {
		action = TaskNextAction{"authorize_worktree", "Worktree creation requires a separate authorized action.", []string{}}
	}
	if ready {
		action = TaskNextAction{"prepare_handoff", "Prepare a handoff for the exact selected revision.", []string{}}
	}
	if len(results) > 0 {
		action = TaskNextAction{"separate_result_control", "Return the historical result for separate result control.", []string{}}
	}
	if fresh == "unknown" || sfresh == "unknown" || eval.ContentIntegrity == "missing" || eval.ContentIntegrity == "conflict" || eval.ContentIntegrity == "unknown" || problem["integrity"] != "valid" && problem["integrity"] != "not_recorded" {
		action = TaskNextAction{"inspect_task", "Resolve the reported content or basis conflict.", []string{"ply", "workspace", "task", "show", string(t.ID)}}
	}
	resource := result
	resource.Content = nil
	resource.Integration = nil
	b, _ := MarshalTaskReadback(resource)
	rv, _ := canonicaljson.DecodeStrict(b)
	var integration canonicaljson.Value
	if result.Integration != nil {
		bytes, e := MarshalTaskIntegrationReadback(*result.Integration)
		if e == nil {
			integration, _ = canonicaljson.DecodeStrict(bytes)
		}
	}
	mode := "legacy"
	if taskRequiresSpec(r, t.ID) {
		mode = "spec_required"
	}
	m := map[string]canonicaljson.Value{"workspace": root, "registry_sha256": nullableDigest(r.RawSHA256), "observed_at_utc": d.WorkClock.Now().UTC().Format(time.RFC3339Nano), "mode": mode, "resource": rv, "integration": integration, "problem": contentObject(problem), "solution": contentObject(map[string]canonicaljson.Value{"spec_id": specID, "head": contentRefValue(state.SpecHead), "selected_revision": contentRefValue(selected), "newer_draft_available": newer}), "selection": sel, "results": results, "ready_for_spec_handoff": ready, "freshness": contentObject(map[string]canonicaljson.Value{"registry": fresh, "content": eval.ContentIntegrity, "selection": sfresh, "target": eval.TargetFreshness}), "reasons": sortedContentStrings(sortedReasons(reasons)), "next_action": contentActionValue(action), "next_transition_authorized": false}
	o := contentFields(contentEnvelope("WorkspaceTaskReadback@3", m))
	o["schema_version"] = int64(3)
	return contentObject(o)
}
func TaskContentText(o canonicaljson.Object) string {
	m := contentFields(o)
	var out strings.Builder
	write := func(format string, args ...any) { fmt.Fprintf(&out, format, args...) }
	revision := func(v canonicaljson.Value) string {
		r := valueRevision(v)
		if r == nil {
			return "Not recorded"
		}
		return fmt.Sprintf("revision %d", r.Revision)
	}
	kind := contentString(m, "kind")
	if kind == "WorkspaceTaskReadback@3" {
		p, s, sel := contentFields(m["problem"]), contentFields(m["solution"]), contentFields(m["selection"])
		title := contentString(p, "title")
		if title == "" {
			title = "Problem content unavailable"
		}
		write("Task: %s\nProblem: %s\n", title, revision(p["head"]))
		if summary := contentString(p, "summary"); summary != "" {
			write("%s\n", summary)
		}
		if s["selected_revision"] != nil {
			write("Selected solution: %s, %s\n", contentString(s, "spec_id"), revision(s["selected_revision"]))
		} else {
			write("Selected solution: Not recorded\n")
		}
		if s["newer_draft_available"] == true {
			write("Newer draft: %s; the current choice has not changed\n", revision(s["head"]))
		}
		write("Selection: %s\n", contentString(sel, "freshness"))
		if m["ready_for_spec_handoff"] == true {
			write("Readiness: Ready for handoff preparation; no agent start authorized\n")
		} else {
			write("Readiness: Not ready for handoff preparation; no agent start authorized\n")
		}
		fresh := contentFields(m["freshness"])
		write("Documents: %s\n", contentString(fresh, "content"))
		results := contentArray(m, "results")
		if len(results) == 0 {
			write("Result: No controlled delivery result recorded\n")
		}
		for _, raw := range results {
			r := contentFields(raw)
			write("Result: %s; technical gate %s\n", contentString(r, "task_result_id"), contentString(r, "technical_gate"))
			link := contentFields(r["spec_binding"])
			if len(link) == 0 {
				write("Solution basis: Not recorded for this historical result.\n")
			} else {
				basis := contentFields(link["basis"])
				write("Implemented solution: %s, %s; problem %s; relevance %s\n", contentString(basis, "spec_id"), revision(basis["spec"]), revision(basis["problem"]), contentString(r, "relevance"))
			}
		}
	} else {
		write("Task: %s\n", contentString(m, "task_id"))
		switch kind {
		case "WorkspaceTaskProblemReadback@1", "WorkspaceTaskSpecReadback@1":
			r := contentFields(m["revision"])
			if len(r) == 0 {
				legacy := contentFields(m["legacy_problem_summary"])
				write("Problem: %s\n%s\n", contentString(legacy, "title"), contentString(legacy, "description"))
			} else {
				write("%s: revision %d\n", contentString(r, "title"), contentInt(r, "revision"))
				if summary := contentString(r, "summary"); summary != "" {
					write("%s\n", summary)
				}
				parts := contentFields(r["parts"])
				for _, name := range []string{"abstract", "functional", "technical"} {
					if part := contentFields(parts[name]); len(part) > 0 {
						write("%s: %s", strings.ToUpper(name[:1])+name[1:], contentString(part, "state"))
						if reason := contentString(part, "reason"); reason != "" {
							write(" (%s)", reason)
						}
						write("\n")
						for _, v := range contentArray(part, "documents") {
							ref := contentFields(v)
							write("  Document: %s\n", contentString(ref, "document_id"))
						}
					}
				}
				for _, raw := range contentArray(r, "documents") {
					d := contentFields(raw)
					write("Snapshot %s: %s\n", contentString(d, "id"), contentString(d, "locator"))
				}
				for _, raw := range contentArray(r, "requirements") {
					req := contentFields(raw)
					write("Requirement %s: %s\n", contentString(req, "id"), contentString(req, "acceptance"))
				}
			}
			write("Documents: %s\n", contentString(m, "integrity"))
			if kind == "WorkspaceTaskSpecReadback@1" {
				write("Selection: %s; target: %s\n", contentString(m, "selection_freshness"), contentString(m, "target_freshness"))
				write("Required handoff inputs: %d preserved files\n", len(contentArray(m, "required_inputs")))
			}
		case "WorkspaceTaskSpecListReadback@1":
			if len(contentArray(m, "revisions")) == 0 {
				write("Solution: Not recorded\n")
			}
			for _, raw := range contentArray(m, "revisions") {
				r := contentFields(raw)
				a := contentFields(r["assessment"])
				write("%s, revision %d: %s; assessment %s; selected %t\n", contentString(m, "spec_id"), contentInt(r, "revision"), contentString(r, "title"), contentString(a, "outcome"), r["selected"] == true)
			}
		case "WorkspaceTaskPublicationReadback@1":
			write("Publication %s: %s; integrity %s; durability %s\n", contentString(m, "publication_key"), contentString(m, "classification"), contentString(m, "integrity"), contentString(m, "durability"))
		}
	}
	write("Next action: %s\n", contentString(contentFields(m["next_action"]), "reason"))
	return out.String()
}
