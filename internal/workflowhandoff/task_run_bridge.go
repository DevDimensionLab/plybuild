package workflowhandoff

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"path/filepath"
	"strings"
)

// TaskRunDraft is a redacted structural projection, never a reply capability.
type TaskRunDraft struct {
	Basis                                     *workspace.TaskSpecBasis
	Worktree, Ref, OID, Title, PublicationKey string
}

func ValidateTaskRunDraft(raw []byte) (TaskRunDraft, error) {
	d, e := decodeHandoffDraft(raw)
	if e != nil {
		return TaskRunDraft{}, e
	}
	b, e := taskSpecBasis(d.Value)
	if e != nil {
		return TaskRunDraft{}, e
	}
	if taskSpecVersion(d.Value) != 2 || b == nil || d.MaxRounds != 4 || d.Binding.Status.Mode != "clean" {
		return TaskRunDraft{}, fmt.Errorf("Task run requires Handoff@2, an exact Task Spec, clean target and four total rounds")
	}
	allowed, _ := objectMember(d.Authority, "allowed_effects")
	for _, v := range allowed.([]canonicaljson.Value) {
		k := objectString(v.(canonicaljson.Object), "type")
		if !setOf("filesystem_write", "git_ref_write", "git_index_write", "git_commit", "command_execute")[k] {
			return TaskRunDraft{}, fmt.Errorf("Task run does not authorize agent effect %s", k)
		}
	}
	return TaskRunDraft{b, d.Binding.TargetWorktree, d.Binding.TargetRef, d.Binding.ExpectedOID, d.GoalTitle, d.PublicationKey}, nil
}

// WithTaskRunSnapshot preserves lock order: taskrun -> Project -> WorkItem -> WF.
func WithTaskRunSnapshot(d Dependencies, root string, fn func(Dependencies) error) error {
	return withTaskSpecWorkspace(d, root, fn)
}

type TaskRunObservation struct {
	Preparation                  workspace.TaskPreparation
	MarkerSHA256, RegistrySHA256 string
	Target, Epic                 workspace.PlanWorktreeObservation
}

func ObserveTaskRun(d Dependencies, root, preparationID string, raw []byte) (TaskRunObservation, error) {
	var out TaskRunObservation
	if d.taskSpecSession == nil {
		e := WithTaskRunSnapshot(d, root, func(locked Dependencies) error {
			var e error
			out, e = ObserveTaskRun(locked, root, preparationID, raw)
			return e
		})
		return out, e
	}
	draft, e := decodeHandoffDraft(raw)
	if e != nil {
		return out, e
	}
	projection, e := ValidateTaskRunDraft(raw)
	if e != nil {
		return out, e
	}
	s := d.taskSpecSession
	p, e := s.TaskRunPreparation(preparationID)
	if e != nil {
		return out, e
	}
	// Preserve typed historical preparation and fresh observations for blocked
	// previews as well. These observations do not authorize a launch.
	if p.Preparation != nil {
		out.Preparation = *p.Preparation
		out.RegistrySHA256 = s.Registry.RawSHA256
		ws, err := d.Workspace.ObserveRoot(root)
		if err == nil {
			out.MarkerSHA256 = ws.Observation.MarkerSHA256
		}
		observe := func(path, ref string) workspace.PlanWorktreeObservation {
			x, _ := s.Dependencies.IntegrationGit.ObserveIntegrationWorktree(path, ref)
			return workspace.PlanWorktreeObservation{WorktreeLocator: x.Locator, Ref: x.Ref, OID: x.OID, Tree: x.Tree, GitCommonDir: x.GitCommonDir, ObjectFormat: x.ObjectFormat, RefFormat: x.RefFormat, Symbolic: x.Symbolic, Clean: x.Clean, StatusEntries: x.StatusEntries, InProgress: x.InProgress}
		}
		plan := p.Preparation.Plan
		out.Target = observe(plan.WorktreePath, "refs/heads/"+plan.Branch)
		out.Epic = observe(plan.Target.ParentLocator, plan.Target.ParentRef)
	}
	if p.Preparation == nil || p.State != "prepared" || p.Disposition != "current" || p.Freshness != "fresh" || len(p.Reasons) != 0 {
		return out, fmt.Errorf("preparation is not current, prepared and fresh")
	}
	plan := p.Preparation.Plan
	b := projection.Basis
	if p.Preparation.Outcome == nil || b.TaskID != plan.TaskID || b.TaskWorktreeID != p.Preparation.Outcome.WorktreeID || b.SpecID != plan.SpecID || b.Spec != plan.Spec || b.Problem != plan.Problem || b.Assessment != plan.Assessment || b.Selection != plan.Selection || projection.Worktree != plan.WorktreePath || projection.Ref != "refs/heads/"+plan.Branch || projection.OID != plan.ParentOID || draft.Binding.ProjectID != plan.Target.ProjectID || draft.Binding.RepoID != plan.Target.RepoID {
		return out, fmt.Errorf("handoff differs from the exact preparation")
	}
	obs, e := observeCreateBindings(d, draft)
	if e != nil {
		return out, e
	}
	if obs.workspace.Observation.Root != root {
		return out, fmt.Errorf("containing workspace differs from request")
	}
	if len(draft.Inputs) != len(plan.RequiredInputs) {
		return out, fmt.Errorf("handoff inputs differ from preparation")
	}
	for i, in := range plan.RequiredInputs {
		v, e := taskSpecJSONValue(in)
		if e != nil {
			return out, e
		}
		inputs, _ := objectMember(draft.Value, "inputs")
		if !canonicalEqual(v, inputs.([]canonicaljson.Value)[i]) {
			return out, fmt.Errorf("handoff input differs from preparation")
		}
	}
	observe := func(path, ref string) (workspace.PlanWorktreeObservation, error) {
		x, e := s.Dependencies.IntegrationGit.ObserveIntegrationWorktree(path, ref)
		return workspace.PlanWorktreeObservation{WorktreeLocator: x.Locator, Ref: x.Ref, OID: x.OID, Tree: x.Tree, GitCommonDir: x.GitCommonDir, ObjectFormat: x.ObjectFormat, RefFormat: x.RefFormat, Symbolic: x.Symbolic, Clean: x.Clean, StatusEntries: x.StatusEntries, InProgress: x.InProgress}, e
	}
	target, e := observe(plan.WorktreePath, "refs/heads/"+plan.Branch)
	if e != nil {
		return out, e
	}
	epic, e := observe(plan.Target.ParentLocator, plan.Target.ParentRef)
	if e != nil {
		return out, e
	}
	for _, x := range []workspace.PlanWorktreeObservation{target, epic} {
		if !x.Clean || !x.Symbolic || len(x.InProgress) != 0 || x.OID != plan.ParentOID || x.Tree != plan.ParentTree || x.GitCommonDir != plan.Target.GitCommonDir {
			return out, fmt.Errorf("target or Epic changed since preparation")
		}
	}
	return TaskRunObservation{*p.Preparation, obs.workspace.Observation.MarkerSHA256, s.Registry.RawSHA256, target, epic}, nil
}

// TaskRunHandoffView omits the secret and is safe in the recipient instructions.
func TaskRunHandoffView(d Dependencies, locator string) ([]byte, error) {
	s, e := d.Store.ReadByLocator(locator)
	if e != nil {
		return nil, e
	}
	return canonicaljson.Marshal(redactHandoff(s.Handoff.Value))
}

type TaskRunLink struct{ ActivityID, RunID, HandoffID, Locator, SHA256 string }

func ReadTaskRunLink(d Dependencies, locator string) (TaskRunLink, error) {
	s, e := d.Store.ReadByLocator(locator)
	h := s.Handoff
	return TaskRunLink{h.Identity.ActivityID, h.Identity.RunID, h.Identity.HandoffID, h.Locator, h.SHA256}, e
}
func bridgeObject(v any) canonicaljson.Object {
	b, _ := json.Marshal(v)
	x, _ := canonicaljson.DecodeStrict(b)
	o, _ := x.(canonicaljson.Object)
	return o
}
func bridgeEnvelope(kind string) canonicaljson.Object {
	return canonicaljson.Object{{Name: "kind", Value: kind}, {Name: "schema_version", Value: int64(1)}, {Name: "format", Value: "json"}, {Name: "format_version", Value: int64(1)}, {Name: "canonicalization", Value: "RFC8785"}}
}
func bridgeBinding(h handoffDocument) canonicaljson.Object {
	return bridgeObject(map[string]any{"activity_id": h.Identity.ActivityID, "run_id": h.Identity.RunID, "handoff_id": h.Identity.HandoffID, "handoff_sha256": h.SHA256})
}

// BuildTaskRunStart fills only observations obtained from the existing observers.
func BuildTaskRunStart(d Dependencies, locator, actor, surface, session, runtime, model string, sandbox, issues []byte, acceptance string) ([]byte, error) {
	s, e := d.Store.ReadByLocator(locator)
	if e != nil {
		return nil, e
	}
	h := s.Handoff
	ws, e := d.Workspace.ObserveRoot(h.Workspace.Root)
	if e != nil {
		return nil, e
	}
	repo, member := findBoundRepository(ws, h.ProjectID, h.RepoID)
	target, e := d.Git.ObserveTarget(TargetRequest{h.Target.Worktree, h.Target.Ref})
	if e != nil {
		return nil, e
	}
	sb, e := canonicaljson.DecodeStrict(sandbox)
	if e != nil {
		return nil, e
	}
	is, e := canonicaljson.DecodeStrict(issues)
	if e != nil {
		return nil, e
	}
	inputs, _ := objectMember(h.Value, "inputs")
	observed := []canonicaljson.Value{}
	for _, v := range inputs.([]canonicaljson.Value) {
		in := v.(canonicaljson.Object)
		path := objectString(in, "locator")
		physical, e := physicalRegularFile(d.Files, path)
		if e != nil {
			return nil, e
		}
		digest, size, e := hashFile(d.Files, physical)
		if e != nil {
			return nil, e
		}
		observed = append(observed, bridgeObject(map[string]any{"id": objectString(in, "id"), "locator": physical, "sha256": digest, "size_bytes": size, "matches_expected": digest == objectString(in, "sha256") && size == objectInt(in, "size_bytes")}))
	}
	contract := canonicaljson.Object{}
	for _, name := range []string{"authority", "budget", "verifiers", "stop_conditions"} {
		v, _ := objectMember(h.Value, name)
		b, _ := canonicaljson.Marshal(v)
		contract = append(contract, canonicaljson.Member{Name: name + "_sha256", Value: digestBytes(b)})
	}
	p := bridgeObject(map[string]any{"expected_principal_id": handoffNestedString(h.Value, "recipient", "principal_id"), "human_start_principal": actor, "start_surface": surface, "session_id": session, "runtime_id": runtime, "model_id": model, "started_at_utc": d.Clock.Now().UTC().Format("2006-01-02T15:04:05.999999999Z07:00")})
	o := bridgeEnvelope("ply.workflow.start-receipt-draft")
	o = replaceObjectMember(o, "schema_version", int64(2))
	basis, _ := objectMember(h.Value, "task_spec_binding")
	for k, v := range map[string]canonicaljson.Value{"receipt_id": h.Identity.StartReceiptID, "binding": bridgeBinding(h), "principal": p, "observed_workspace": bridgeObject(map[string]any{"root": ws.Observation.Root, "marker_format_version": ws.Observation.MarkerFormatVersion, "marker_sha256": ws.Observation.MarkerSHA256, "matches_expected": ws.Observation == h.Workspace}), "observed_project": bridgeObject(map[string]any{"project_id": string(h.ProjectID), "repo_id": string(h.RepoID), "registered_locator": repo.Locator, "registered_git_common_dir": repo.GitCommonDir, "matches_expected": member && repo.Locator == h.RegisteredLocator && repo.GitCommonDir == h.RegisteredGitCommonDir}), "observed_target": append(targetValue(target), canonicaljson.Member{Name: "matches_expected", Value: sameTarget(target, h.Target)}), "observed_inputs": observed, "contract_digests": contract, "sandbox": sb, "acceptance": acceptance, "issues": is, "task_spec_binding": basis} {
		o = append(o, canonicaljson.Member{Name: k, Value: v})
	}
	b, e := canonicaljson.Marshal(o)
	if e != nil {
		return nil, e
	}
	_, e = decodeStartDraft(b)
	return b, e
}

// ValidateTaskRunStart checks the unchanged WF binding and sandbox authority
// before the transport freezes a positive claim. It publishes nothing.
func ValidateTaskRunStart(d Dependencies, locator string, raw []byte) error {
	s, e := d.Store.ReadByLocator(locator)
	if e != nil {
		return e
	}
	draft, e := decodeStartDraft(raw)
	if e != nil {
		return e
	}
	return WithTaskRunSnapshot(d, s.Handoff.Workspace.Root, func(locked Dependencies) error {
		return validateStartBinding(locked, s, draft)
	})
}

// ValidateTaskRunClaimFields uses the WF sandbox and issue vocabulary even when
// missing runtime facts prevent constructing a complete StartReceipt.
func ValidateTaskRunClaimFields(sandbox, issues []byte, acceptance string, nullableSandbox bool) error {
	if !setOf("started", "rejected", "conflict", "unknown")[acceptance] {
		return fmt.Errorf("invalid acceptance")
	}
	sb, e := canonicaljson.DecodeStrict(sandbox)
	if e != nil {
		return e
	}
	if sb != nil || !nullableSandbox || acceptance == "started" {
		if e = validateSandbox(sb); e != nil {
			return e
		}
	}
	value, e := canonicaljson.DecodeStrict(issues)
	if e != nil {
		return e
	}
	entries, ok := value.([]canonicaljson.Value)
	if !ok || (acceptance == "started") != (len(entries) == 0) {
		return fmt.Errorf("started alone has no issues")
	}
	last := ""
	for _, item := range entries {
		fields, err := exactObject(item, "issue", "type", "detail")
		if err != nil {
			return err
		}
		kind, _ := stringField(fields, "type", "issue")
		detail, _ := stringField(fields, "detail", "issue")
		if !setOf("binding", "workspace", "project", "target", "input", "contract", "sandbox", "principal", "unknown")[kind] || validatePlainText("issue detail", detail, 1, 2000) != nil {
			return fmt.Errorf("invalid issue")
		}
		key := kind + "\x00" + detail
		if last != "" && key <= last {
			return fmt.Errorf("issues must be sorted")
		}
		last = key
	}
	return nil
}

// ValidateTaskRunSemantics reuses the terminal field validators without inventing
// missing verifier executions or requiring a representable WF terminal envelope.
func ValidateTaskRunSemantics(raw []byte) error {
	v, e := canonicaljson.DecodeStrict(raw)
	if e != nil {
		return e
	}
	o, ok := v.(canonicaljson.Object)
	if !ok {
		return fmt.Errorf("report must be an object")
	}
	get := func(k string) canonicaljson.Value { v, _ := objectMember(o, k); return v }
	outcome := objectString(o, "outcome")
	if !setOf("complete", "blocked", "conflict", "budget_exhausted", "unknown")[outcome] {
		return fmt.Errorf("invalid outcome")
	}
	if validatePlainText("summary", objectString(o, "summary"), 1, 240) != nil || validatePlainText("meaning", objectString(o, "meaning"), 1, 600) != nil {
		return fmt.Errorf("invalid summary or meaning")
	}
	stops, ok := get("stop_reasons").([]canonicaljson.Value)
	if !ok || (outcome == "complete") != (len(stops) == 0) {
		return fmt.Errorf("invalid stop reasons")
	}
	if e = validateStopReasons(stops, outcome); e != nil {
		return e
	}
	for k, f := range map[string]func(canonicaljson.Value) error{"observed_effects": validateObservedEffects, "verifier_results": validateVerifierResults, "review": validateReview} {
		if e = f(get(k)); e != nil {
			return e
		}
	}
	a, ok := get("artifacts").([]canonicaljson.Value)
	if !ok {
		return fmt.Errorf("artifacts must be an array")
	}
	if e = validateArtifacts(a); e != nil {
		return e
	}
	for _, k := range []string{"evidence_gaps", "forbidden_effects_observed"} {
		if e = validateGapEntries(get(k), k); e != nil {
			return e
		}
	}
	return nil
}
func BuildTaskRunTerminal(d Dependencies, locator string, raw []byte, rounds int) ([]byte, error) {
	if e := ValidateTaskRunSemantics(raw); e != nil {
		return nil, e
	}
	s, e := d.Store.ReadByLocator(locator)
	if e != nil {
		return nil, e
	}
	if s.Start == nil {
		return nil, fmt.Errorf("start receipt missing")
	}
	v, _ := canonicaljson.DecodeStrict(raw)
	report := v.(canonicaljson.Object)
	target, e := d.Git.ObserveTarget(TargetRequest{s.Handoff.Target.Worktree, s.Handoff.Target.Ref})
	if e != nil {
		return nil, e
	}
	p, _ := objectMember(s.Start.Value, "principal")
	o := bridgeEnvelope("ply.workflow.terminal-result-draft")
	for k, v := range map[string]canonicaljson.Value{"result_id": s.Handoff.Identity.TerminalResultID, "binding": bridgeBinding(s.Handoff), "start_binding": bridgeObject(map[string]any{"receipt_id": s.Start.DocumentID, "start_receipt_sha256": s.Start.SHA256}), "principal": p, "reported_outcome": objectString(report, "outcome"), "rounds_used": int64(rounds), "final_target": append(targetValue(target), canonicaljson.Member{Name: "matches_expected", Value: true})} {
		o = append(o, canonicaljson.Member{Name: k, Value: v})
	}
	for _, k := range []string{"summary", "meaning", "stop_reasons", "observed_effects", "verifier_results", "review", "artifacts", "evidence_gaps", "forbidden_effects_observed"} {
		v, _ := objectMember(report, k)
		o = append(o, canonicaljson.Member{Name: k, Value: v})
	}
	b, e := canonicaljson.Marshal(o)
	if e != nil {
		return nil, e
	}
	_, e = decodeTerminalDraft(b)
	return b, e
}

// BuildTaskRunResult freezes the exact existing TaskResult draft, including the
// inspection digest and recorder time. RecordTaskResult remains the final gate.
func BuildTaskRunResult(d Dependencies, locator, runID string, technical []byte) ([]byte, string, error) {
	s, e := d.Store.ReadByLocator(locator)
	if e != nil {
		return nil, "", e
	}
	if s.Start == nil || s.Terminal == nil {
		return nil, "", fmt.Errorf("start or terminal receipt missing")
	}
	basis, e := taskSpecBasis(s.Handoff.Value)
	if e != nil || basis == nil {
		return nil, "", fmt.Errorf("Task Spec basis missing")
	}
	inspection, e := Inspect(d, InspectInput{HandoffLocator: locator, Format: "json"})
	if e != nil {
		return nil, "", e
	}
	digest := digestBytes(inspection.Bytes)
	h := s.Handoff
	start := s.Start
	end := s.Terminal
	evidence, e := NewTaskHandoffEvidenceReader(d).ReadTaskEvidence(workspace.TaskHandoffEvidenceRequest{ActivityID: h.Identity.ActivityID, RunID: h.Identity.RunID, HandoffID: h.Identity.HandoffID, HandoffLocator: h.Locator, HandoffSHA256: h.SHA256, StartReceiptID: start.DocumentID, StartReceiptLocator: start.Locator, StartReceiptSHA256: start.SHA256, TerminalResultID: end.DocumentID, TerminalResultLocator: end.Locator, TerminalResultSHA256: end.SHA256, InspectionSHA256: digest})
	if e != nil {
		return nil, "", e
	}
	var ta struct {
		AcceptedDebt []workspace.TaskAcceptedDebtRecord `json:"accepted_debt"`
	}
	if e = json.Unmarshal(technical, &ta); e != nil {
		return nil, "", e
	}
	for _, debt := range ta.AcceptedDebt {
		for _, id := range debt.EvidenceArtifactIDs {
			for i, a := range evidence.Artifacts {
				if a.ArtifactID == id && a.Role == "other" {
					evidence.Artifacts[i].Role = "debt_control"
				}
			}
		}
	}
	draft := map[string]any{"kind": "WorkspaceTaskResultRecordDraft@1", "schema_version": 1, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "publication_key": "task-run/" + runID, "task_id": basis.TaskID, "task_worktree_id": basis.TaskWorktreeID,
		"handoff": map[string]any{"activity_id": evidence.ActivityID, "run_id": evidence.RunID, "handoff_id": evidence.HandoffID, "handoff_locator": evidence.HandoffLocator, "handoff_sha256": evidence.HandoffSHA256, "start_receipt_id": evidence.StartReceiptID, "start_receipt_locator": evidence.StartReceiptLocator, "start_receipt_sha256": evidence.StartReceiptSHA256, "terminal_result_id": evidence.TerminalResultID, "terminal_result_locator": evidence.TerminalResultLocator, "terminal_result_sha256": evidence.TerminalResultSHA256, "inspection_sha256": digest},
		"source":  map[string]any{"project_id": h.ProjectID, "repo_id": h.RepoID, "git_common_dir": evidence.GitCommonDir, "worktree_locator": evidence.TargetWorktree, "source_ref": evidence.TargetRef, "result_oid": evidence.ResultOID, "result_tree": evidence.ResultTree}, "technical_assessment": json.RawMessage(technical), "evidence_artifacts": evidence.Artifacts, "recorder": map[string]any{"actor_claim": "ply-task-run", "control_surface": "ply workspace task run collect", "recorded_at_utc": d.Clock.Now().UTC().Format("2006-01-02T15:04:05.999999999Z07:00")}}
	b, e := json.Marshal(draft)
	if e != nil {
		return nil, "", e
	}
	v, e := canonicaljson.DecodeStrict(b)
	if e != nil {
		return nil, "", e
	}
	b, e = canonicaljson.Marshal(v)
	return b, digest, e
}

// TaskRunReturnFacts exposes hashes and target facts, never secret material.
type TaskRunReturnFacts struct {
	Start, Terminal       *SubmitResult
	Acceptance, SessionID string
	FinalTarget           *workspace.PlanWorktreeObservation
}

func ReadTaskRunReturnFacts(d Dependencies, locator string) (TaskRunReturnFacts, error) {
	var out TaskRunReturnFacts
	s, e := d.Store.ReadByLocator(locator)
	if e != nil {
		return out, e
	}
	if s.Start != nil {
		out.Start = &SubmitResult{Phase: "start", DocumentID: s.Start.DocumentID, Locator: s.Start.Locator, SHA256: s.Start.SHA256}
		out.Acceptance = s.Start.Outcome
		p, _ := objectMember(s.Start.Value, "principal")
		out.SessionID = objectString(p.(canonicaljson.Object), "session_id")
	}
	if s.Terminal != nil {
		out.Terminal = &SubmitResult{Phase: "terminal", DocumentID: s.Terminal.DocumentID, Locator: s.Terminal.Locator, SHA256: s.Terminal.SHA256}
		v, _ := objectMember(s.Terminal.Value, "final_target")
		x := v.(canonicaljson.Object)
		status, _ := objectMember(x, "status_policy")
		sp, _ := validateStatusPolicy(status)
		out.FinalTarget = &workspace.PlanWorktreeObservation{WorktreeLocator: objectString(x, "worktree"), Ref: objectString(x, "ref"), OID: objectString(x, "oid"), Tree: objectString(x, "tree"), GitCommonDir: objectString(x, "git_common_dir"), Clean: sp.Mode == "clean"}
	}
	return out, nil
}

// RecoverTaskRunTerminal recognizes only the already accepted immutable draft.
// Managed locators are the deterministic WF projection, not a new observation.
func RecoverTaskRunTerminal(d Dependencies, locator string, raw []byte) (*SubmitResult, error) {
	s, e := d.Store.ReadByLocator(locator)
	if e != nil {
		return nil, e
	}
	if s.Terminal == nil {
		return nil, nil
	}
	draft, e := decodeTerminalDraft(raw)
	if e != nil {
		return nil, e
	}
	o := draft.Value
	artifacts := []canonicaljson.Value{}
	for _, v := range draft.Artifacts {
		a := v.(canonicaljson.Object)
		if objectString(a, "kind") == "managed" {
			a = replaceObjectMember(a, "locator", filepath.Join(s.Handoff.ReplyRoot, "artifacts", "sha256", strings.TrimPrefix(objectString(a, "sha256"), "sha256:")))
		}
		artifacts = append(artifacts, a)
	}
	o = replaceObjectMember(o, "artifacts", artifacts)
	expected := replaceObjectMember(removeObjectMember(s.Terminal.Value, "capability_proof"), "kind", "ply.workflow.terminal-result-draft")
	if !canonicalEqual(o, expected) {
		return nil, fmt.Errorf("preserved terminal draft differs from accepted terminal")
	}
	return &SubmitResult{Phase: "terminal", DocumentID: s.Terminal.DocumentID, Locator: s.Terminal.Locator, SHA256: s.Terminal.SHA256}, nil
}

func RecoverTaskRunStart(d Dependencies, locator string, raw []byte) (*SubmitResult, error) {
	s, e := d.Store.ReadByLocator(locator)
	if e != nil {
		return nil, e
	}
	if s.Start == nil {
		return nil, nil
	}
	draft, e := decodeStartDraft(raw)
	if e != nil {
		return nil, e
	}
	original := replaceObjectMember(removeObjectMember(s.Start.Value, "capability_proof"), "kind", "ply.workflow.start-receipt-draft")
	if taskSpecVersion(original) == 2 {
		original = removeObjectMember(original, "task_spec_observation")
	}
	if !canonicalEqual(original, draft.Value) {
		return nil, fmt.Errorf("preserved start draft differs from accepted receipt")
	}
	return &SubmitResult{Phase: "start", DocumentID: s.Start.DocumentID, Locator: s.Start.Locator, SHA256: s.Start.SHA256}, nil
}
