package workflowhandoff

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// DeliveryCandidateInput contains observations from the delivery runtime's actual
// verifier invocation. This bridge does not execute the command or infer review.
type DeliveryCandidateInput struct {
	ParentHandoffLocator, CandidateKey, Summary          string
	RunID, RequestSHA256                                 string
	CandidateOID, CandidateTree, VerifierID, CWD         string
	Argv                                                 []string
	Exit                                                 int
	StdoutPath, StderrPath, ReviewPath, VerificationPath string
	RequireIndependentReview                             bool
}

type DeliveryCandidateResult struct {
	Handoff         TaskRunLink
	Start, Terminal SubmitResult
	TaskResult      workspace.TaskResultRecord
}

type deliveryCandidateArtifact struct {
	id, media string
	raw       []byte
}

// DeliveryCandidateReviewTemplate is a schema-shaped starting point, not review
// evidence. The explicit unobserved decision and empty reviewer keep it invalid
// until the reviewer records an actual candidate-bound assessment.
func DeliveryCandidateReviewTemplate(oid, tree string) ([]byte, error) {
	return canonicaljson.Marshal(deliveryObject(map[string]any{
		"kind": "DeliveryCandidateReview@1", "schema_version": 1,
		"candidate_oid": oid, "candidate_tree": tree,
		"reviewer_claim": "", "reviewer_session_id": nil, "decision": "unobserved",
		"findings": []any{}, "fixes": []any{}, "open_actionable_findings": []any{},
	}))
}

// QualifyDeliveryCandidate creates a separate immutable native control handoff
// for this candidate. Its one round controls preserved verifier evidence; it
// does not backdate a start or reclassify the owner's prior command execution as
// a new effect. It never limits the owner's implementation or correction work.
func QualifyDeliveryCandidate(d Dependencies, in DeliveryCandidateInput) (DeliveryCandidateResult, error) {
	var out DeliveryCandidateResult
	if in.Exit != 0 || !validateDigest(in.RequestSHA256) || validatePlainText("run_id", in.RunID, 1, 256) != nil || !validateOID(in.CandidateOID) || !validateOID(in.CandidateTree) || validateKey("candidate_key", in.CandidateKey) != nil || validatePlainText("summary", in.Summary, 1, 240) != nil {
		return out, fmt.Errorf("candidate qualification requires a successful observed verifier and exact candidate")
	}
	parent, e := d.Store.ReadByLocator(in.ParentHandoffLocator)
	if e != nil {
		return out, e
	}
	if taskSpecVersion(parent.Handoff.Value) != 3 || deliveryMode(parent.Handoff.Value) != "owner" || parent.Start == nil || parent.Start.Outcome != "started" || parent.Terminal != nil || parent.HeadState != "ready" {
		return out, fmt.Errorf("active accepted delivery owner is required")
	}
	principal, _ := objectMember(parent.Start.Value, "principal")
	p := principal.(canonicaljson.Object)
	verifiers, _ := objectMember(parent.Handoff.Value, "verifiers")
	vs := verifiers.([]canonicaljson.Value)
	if len(vs) != 1 {
		return out, fmt.Errorf("candidate control requires one acceptance verifier")
	}
	verifier := vs[0].(canonicaljson.Object)
	argv, _ := objectMember(verifier, "argv")
	if in.VerifierID != objectString(verifier, "id") || in.CWD != objectString(verifier, "cwd") || !canonicalEqual(argv, bridgeValue(in.Argv)) {
		return out, fmt.Errorf("observed verifier differs from the owner contract")
	}
	target, e := d.Git.ObserveTarget(TargetRequest{parent.Handoff.Target.Worktree, parent.Handoff.Target.Ref})
	if e != nil {
		return out, e
	}
	if target.OID != in.CandidateOID || target.Tree != in.CandidateTree || target.Status.Mode != "clean" {
		return out, fmt.Errorf("candidate changed after verifier execution")
	}
	reviewRaw, e := deliveryReadFile(d.Files, in.ReviewPath, 256<<10)
	if e != nil {
		return out, e
	}
	review, e := deliveryCandidateReview(reviewRaw, in, objectString(p, "session_id"))
	if e != nil {
		return out, e
	}
	stdout, e := deliveryReadFile(d.Files, in.StdoutPath, 64<<20)
	if e != nil {
		return out, e
	}
	stderr, e := deliveryReadFile(d.Files, in.StderrPath, 64<<20)
	if e != nil {
		return out, e
	}
	verification, e := deliveryReadFile(d.Files, in.VerificationPath, 256<<10)
	if e != nil {
		return out, e
	}
	acceptanceSnapshot, e := validateDeliveryVerification(d.Files, verification, in, stdout, stderr, reviewRaw)
	if e != nil {
		return out, e
	}
	basis, e := taskSpecBasis(parent.Handoff.Value)
	if e != nil || basis == nil {
		return out, fmt.Errorf("candidate execution Spec basis missing")
	}
	_, eval, e := historicalHandoffTaskSpec(d, parent)
	if e != nil {
		return out, e
	}
	// Assemble the actual managed evidence before any candidate handoff or
	// staging write. Review references must resolve against this same set.
	artifactInputs := []deliveryCandidateArtifact{{"acceptance-script", "text/plain", acceptanceSnapshot}, {"candidate-review", "application/json", reviewRaw}, {"verification-receipt", "application/json", verification}, {"verifier-stderr", "text/plain", stderr}, {"verifier-stdout", "text/plain", stdout}}
	reqs, _ := objectMember(eval.Spec, "requirements")
	coverage := []canonicaljson.Value{}
	for _, r := range reqs.([]canonicaljson.Value) {
		coverage = append(coverage, deliveryObject(map[string]any{"id": objectString(r.(canonicaljson.Object), "id"), "outcome": "passed", "verifier_ids": []string{in.VerifierID}, "artifact_ids": []string{"candidate-review", "verifier-stderr", "verifier-stdout"}, "reason": "The goal's executable acceptance passed for this exact candidate; explicit candidate review is preserved."}))
	}
	coverageRaw, e := canonicaljson.Marshal(deliveryObject(map[string]any{"kind": "WorkspaceTaskRequirementEvidence@1", "schema_version": 1, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "task_id": basis.TaskID, "spec_id": basis.SpecID, "spec": basis.Spec, "result_oid": in.CandidateOID, "result_tree": in.CandidateTree, "requirements": coverage}))
	if e != nil {
		return out, e
	}
	artifactInputs = append(artifactInputs, deliveryCandidateArtifact{"task-requirements", "application/json", coverageRaw})
	if e = deliveryCandidateReviewEvidence(review, artifactInputs); e != nil {
		return out, e
	}
	draft := bridgeEnvelope("ply.workflow.handoff-draft")
	draft = replaceObjectMember(draft, "schema_version", int64(3))
	for _, key := range []string{"recipient", "inputs", "verifiers", "stop_conditions", "reporting", "task_spec_binding"} {
		v, _ := objectMember(parent.Handoff.Value, key)
		draft = append(draft, canonicaljson.Member{Name: key, Value: v})
	}
	binding, _ := objectMember(parent.Handoff.Value, "delivery_binding")
	ob := binding.(canonicaljson.Object)
	publication := parent.Handoff.Identity.PublicationKey + "/candidate/" + in.CandidateKey
	fields := map[string]any{
		"publication_key": publication, "activity_key": parent.Handoff.Identity.ActivityKey,
		"goal":            map[string]any{"title": "Verify delivery candidate", "recipient_role": "Delivery owner performing candidate control", "objective": "Record the observed acceptance command and explicit review for this exact candidate.", "done_when": "Native technical evidence qualifies this exact candidate for actual human QA."},
		"binding_request": map[string]any{"project_id": parent.Handoff.ProjectID, "repo_id": parent.Handoff.RepoID, "target_worktree": target.Worktree, "target_ref": target.Ref, "expected_oid": target.OID, "status_policy": map[string]any{"mode": "clean"}},
		"authority":       map[string]any{"allowed_effects": []canonicaljson.Value{}, "forbidden_effects": []canonicaljson.Value{}, "human_gates": []string{"human_task_qa"}},
		"budget":          deliveryBudget(1), "procedure": []canonicaljson.Value{deliveryObject(map[string]any{"id": "verify-candidate", "instruction": "Control the preserved acceptance invocation executed under the already accepted delivery-owner handoff, explicit review and task-requirements artifact. This control starts now; prior command execution is not a new effect.", "required_before": []string{}})},
		"delivery_binding": map[string]any{"mode": "candidate", "human_actor": objectString(ob, "human_actor"), "owner_claim": objectString(ob, "owner_claim"), "parent_handoff_locator": parent.Handoff.Locator, "parent_handoff_sha256": parent.Handoff.SHA256, "candidate_key": in.CandidateKey},
	}
	for k, v := range fields {
		draft = append(draft, canonicaljson.Member{Name: k, Value: bridgeValue(v)})
	}
	raw, e := canonicaljson.Marshal(draft)
	if e != nil {
		return out, e
	}
	// The owner's reply area is already managed and writable under its accepted
	// runtime. Deterministic files make retries compare bytes, never overwrite.
	workarea := filepath.Join(parent.Handoff.ReplyRoot, "staging", "candidate-"+digestBytes([]byte(in.CandidateKey))[7:])
	if e = deliveryEnsureDir(d.Files, workarea); e != nil {
		return out, e
	}
	draftPath := filepath.Join(workarea, "handoff.json")
	if e = deliveryWriteOnce(d.Files, draftPath, raw); e != nil {
		return out, e
	}
	created, e := Create(d, CreateInput{DraftPath: draftPath})
	if e != nil {
		return out, e
	}
	out.Handoff, e = ReadTaskRunLink(d, created.Locator)
	if e != nil {
		return out, e
	}
	candidate, e := d.Store.ReadByLocator(created.Locator)
	if e != nil {
		return out, e
	}
	if candidate.Start == nil {
		path := filepath.Join(workarea, "start.json")
		if _, err := d.Files.Lstat(path); os.IsNotExist(err) {
			sandbox, _ := objectMember(parent.Start.Value, "sandbox")
			sb, _ := canonicaljson.Marshal(sandbox)
			start, err := BuildDeliveryTaskRunStart(d, created.Locator, objectString(p, "human_start_principal"), objectString(p, "start_surface"), objectString(p, "session_id"), objectString(p, "runtime_id"), objectString(p, "model_id"), sb, []byte("[]"), "started")
			if err != nil {
				return out, err
			}
			if err = deliveryWriteOnce(d.Files, path, start); err != nil {
				return out, err
			}
		} else if err != nil {
			return out, err
		}
		out.Start, e = SubmitStart(d, SubmitInput{HandoffLocator: created.Locator, DraftPath: path})
		if e != nil {
			return out, e
		}
	} else {
		out.Start = SubmitResult{Phase: "start", DocumentID: candidate.Start.DocumentID, Locator: candidate.Start.Locator, SHA256: candidate.Start.SHA256}
	}
	candidate, e = d.Store.ReadByLocator(created.Locator)
	if e != nil {
		return out, e
	}
	staging := filepath.Join(candidate.Handoff.ReplyRoot, "staging")
	if e = deliveryEnsureDir(d.Files, staging); e != nil {
		return out, e
	}
	artifacts := []canonicaljson.Value{}
	addArtifact := func(id, media string, b []byte) error {
		path := filepath.Join(staging, id)
		if e := deliveryWriteOnce(d.Files, path, b); e != nil {
			return e
		}
		artifacts = append(artifacts, deliveryObject(map[string]any{"artifact_id": id, "kind": "managed", "description": "Observed delivery candidate " + id, "media_type": media, "classification": "workspace_internal", "size_bytes": len(b), "sha256": digestBytes(b), "locator": path}))
		return nil
	}
	for _, a := range artifactInputs {
		if e = addArtifact(a.id, a.media, a.raw); e != nil {
			return out, e
		}
	}
	sort.Slice(artifacts, func(i, j int) bool {
		return objectString(artifacts[i].(canonicaljson.Object), "artifact_id") < objectString(artifacts[j].(canonicaljson.Object), "artifact_id")
	})
	report := deliveryObject(map[string]any{"outcome": "complete", "summary": in.Summary, "meaning": "This exact candidate passed the preserved acceptance invocation under its delivery owner and explicit technical review. Actual human QA is still required.", "stop_reasons": []canonicaljson.Value{}, "observed_effects": []canonicaljson.Value{}, "verifier_results": []canonicaljson.Value{deliveryObject(map[string]any{"verifier_id": in.VerifierID, "argv": in.Argv, "cwd": in.CWD, "exit": in.Exit, "bound_oid_or_sha256": in.CandidateOID, "stdout_artifact_id": "verifier-stdout", "stderr_artifact_id": "verifier-stderr"})}, "review": review, "artifacts": artifacts, "evidence_gaps": []canonicaljson.Value{}, "forbidden_effects_observed": []canonicaljson.Value{}})
	reportRaw, e := canonicaljson.Marshal(report)
	if e != nil {
		return out, e
	}
	terminal, e := BuildTaskRunTerminal(d, created.Locator, reportRaw, 1)
	if e != nil {
		return out, e
	}
	technical, _ := canonicaljson.Marshal(deliveryObject(map[string]any{"gate": "passed", "required_verifier_ids": []string{in.VerifierID}, "accepted_debt": []canonicaljson.Value{}}))
	if e = ValidateWorkflowRound(d, created.Locator, terminal); e != nil {
		return out, e
	}
	if e = ValidateWorkflowPassedAssessment(d, created.Locator, terminal, technical); e != nil {
		return out, e
	}
	terminalPath := filepath.Join(workarea, "terminal.json")
	if e = deliveryWriteOnce(d.Files, terminalPath, terminal); e != nil {
		return out, e
	}
	if candidate.Terminal != nil {
		recovered, err := RecoverTaskRunTerminal(d, created.Locator, terminal)
		if err != nil {
			return out, fmt.Errorf("recover candidate terminal: %w", err)
		}
		if recovered == nil {
			return out, fmt.Errorf("candidate terminal disappeared during recovery")
		}
		out.Terminal = *recovered
	} else {
		out.Terminal, e = SubmitResultDocument(d, SubmitInput{HandoffLocator: created.Locator, DraftPath: terminalPath})
		if e != nil {
			return out, e
		}
	}
	// A stable timestamp is needed when recovering the same immutable publication.
	cp, _ := objectMember(candidate.Start.Value, "principal")
	recorder := workspace.TaskRecorderRecord{ActorClaim: objectString(ob, "owner_claim"), ControlSurface: "ply workspace task execute verify", RecordedAtUTC: objectString(cp.(canonicaljson.Object), "started_at_utc")}
	resultRaw, _, e := BuildDeliveryTaskRunResult(d, created.Locator, publication, technical, recorder)
	if e != nil {
		return out, e
	}
	resultPath := filepath.Join(workarea, "task-result.json")
	if e = deliveryWriteOnce(d.Files, resultPath, resultRaw); e != nil {
		return out, e
	}
	if d.TaskWorkspace == nil {
		return out, fmt.Errorf("Task workspace dependency missing")
	}
	wd := *d.TaskWorkspace
	wd.HandoffEvidence = NewTaskHandoffEvidenceReader(d)
	recorded, e := workspace.RecordTaskResult(wd, workspace.TaskResultRecordInput{TaskID: basis.TaskID, File: resultPath})
	if e != nil {
		return out, e
	}
	out.TaskResult = recorded.Record
	return out, nil
}

func validateDeliveryVerification(files FileSystem, raw []byte, in DeliveryCandidateInput, stdout, stderr, review []byte) ([]byte, error) {
	v, e := canonicaljson.DecodeStrict(raw)
	if e != nil {
		return nil, e
	}
	o, ok := v.(canonicaljson.Object)
	if !ok {
		return nil, fmt.Errorf("verification receipt must be an object")
	}
	keys := []string{"kind", "schema_version", "run_id", "request_sha256", "attempt_id", "candidate_oid", "candidate_tree", "argv", "cwd", "acceptance", "acceptance_snapshot", "review", "exit", "stdout", "stderr", "started_at", "finished_at"}
	if errorValue, found := objectMember(o, "error"); found {
		keys = append(keys, "error")
		if errorValue != "" {
			return nil, fmt.Errorf("verification receipt contains an execution error")
		}
	}
	m, e := exactObject(o, "verification receipt", keys...)
	if e != nil {
		return nil, e
	}
	version, e := intField(m, "schema_version", "verification receipt")
	exit, xe := intField(m, "exit", "verification receipt")
	if e != nil || xe != nil || version != 1 || exit != 0 || objectMapString(m, "kind") != "PlyDeliveryVerification@1" || objectMapString(m, "candidate_oid") != in.CandidateOID || objectMapString(m, "candidate_tree") != in.CandidateTree || objectMapString(m, "cwd") != in.CWD || !canonicalEqual(m["argv"], bridgeValue(in.Argv)) || objectMapString(m, "request_sha256") != in.RequestSHA256 || objectMapString(m, "run_id") != in.RunID || objectMapString(m, "attempt_id") != in.CandidateKey {
		return nil, fmt.Errorf("verification receipt differs from the observed candidate invocation")
	}
	start, se := time.Parse(time.RFC3339Nano, objectMapString(m, "started_at"))
	finish, fe := time.Parse(time.RFC3339Nano, objectMapString(m, "finished_at"))
	if se != nil || fe != nil || finish.Before(start) {
		return nil, fmt.Errorf("verification receipt time binding invalid")
	}
	fileBinding := func(key, path string, b []byte) (map[string]canonicaljson.Value, error) {
		f, e := exactObject(m[key], key, "locator", "sha256")
		if e != nil {
			return nil, e
		}
		if objectMapString(f, "locator") != path || objectMapString(f, "sha256") != digestBytes(b) {
			return nil, fmt.Errorf("verification %s binding differs", key)
		}
		return f, nil
	}
	if _, e = fileBinding("stdout", in.StdoutPath, stdout); e != nil {
		return nil, e
	}
	if _, e = fileBinding("stderr", in.StderrPath, stderr); e != nil {
		return nil, e
	}
	if _, e = fileBinding("review", in.ReviewPath, review); e != nil {
		return nil, e
	}
	a, e := exactObject(m["acceptance"], "acceptance", "locator", "sha256")
	if e != nil {
		return nil, e
	}
	s, e := exactObject(m["acceptance_snapshot"], "acceptance snapshot", "locator", "sha256")
	if e != nil {
		return nil, e
	}
	if len(in.Argv) != 2 || objectMapString(a, "locator") != in.Argv[1] || objectMapString(a, "sha256") != objectMapString(s, "sha256") {
		return nil, fmt.Errorf("acceptance snapshot differs from the executed script")
	}
	b, e := deliveryReadFile(files, objectMapString(s, "locator"), 4<<20)
	if e != nil {
		return nil, e
	}
	if digestBytes(b) != objectMapString(s, "sha256") {
		return nil, fmt.Errorf("acceptance snapshot hash differs")
	}
	return b, nil
}

func deliveryCandidateReview(raw []byte, in DeliveryCandidateInput, ownerSession string) (canonicaljson.Object, error) {
	v, e := canonicaljson.DecodeStrict(raw)
	if e != nil {
		return nil, e
	}
	m, e := exactObject(v, "candidate review", "kind", "schema_version", "candidate_oid", "candidate_tree", "reviewer_claim", "reviewer_session_id", "decision", "findings", "fixes", "open_actionable_findings")
	if e != nil {
		return nil, e
	}
	version, e := intField(m, "schema_version", "candidate review")
	if e != nil || version != 1 || objectMapString(m, "kind") != "DeliveryCandidateReview@1" || objectMapString(m, "candidate_oid") != in.CandidateOID || objectMapString(m, "candidate_tree") != in.CandidateTree || objectMapString(m, "decision") != "passed" {
		return nil, fmt.Errorf("candidate review must explicitly pass this exact candidate")
	}
	claim, e := stringField(m, "reviewer_claim", "candidate review")
	if e != nil {
		return nil, e
	}
	if e = validatePlainText("candidate review reviewer_claim", claim, 1, 256); e != nil {
		return nil, e
	}
	session, ok := m["reviewer_session_id"].(string)
	if m["reviewer_session_id"] != nil && (!ok || validatePlainText("reviewer_session_id", session, 1, 256) != nil) {
		return nil, fmt.Errorf("invalid reviewer session")
	}
	if in.RequireIndependentReview && (!ok || session == ownerSession) {
		return nil, fmt.Errorf("independent review requires a distinct actual reviewer session")
	}
	review := canonicaljson.Object{}
	for _, k := range []string{"findings", "fixes", "open_actionable_findings"} {
		review = append(review, canonicaljson.Member{Name: k, Value: m[k]})
	}
	if e = validateReview(review); e != nil {
		return nil, e
	}
	if len(m["open_actionable_findings"].([]canonicaljson.Value)) != 0 {
		return nil, fmt.Errorf("candidate review has open actionable findings")
	}
	return review, nil
}

func deliveryCandidateReviewEvidence(review canonicaljson.Object, artifacts []deliveryCandidateArtifact) error {
	known := make(map[string]bool, len(artifacts))
	ids := make([]string, 0, len(artifacts))
	for _, artifact := range artifacts {
		known[artifact.id] = true
		ids = append(ids, artifact.id)
	}
	sort.Strings(ids)
	// deliveryCandidateReview has already validated entry and array shapes.
	for _, list := range review {
		for i, value := range list.Value.([]canonicaljson.Value) {
			entry := value.(canonicaljson.Object)
			refs, _ := objectMember(entry, "evidence_ids")
			for _, ref := range refs.([]canonicaljson.Value) {
				if !known[ref.(string)] {
					return fmt.Errorf("candidate review %s[%d].evidence_ids references unavailable artifact %q; available artifacts: %s", list.Name, i, ref, strings.Join(ids, ", "))
				}
			}
		}
	}
	return nil
}

func deliveryReadFile(files FileSystem, path string, limit int64) ([]byte, error) {
	physical, e := physicalRegularFile(files, path)
	if e != nil {
		return nil, e
	}
	if physical != path {
		return nil, fmt.Errorf("evidence path must be physical")
	}
	st, e := files.Lstat(path)
	if e != nil {
		return nil, e
	}
	if st.Size() > limit {
		return nil, fmt.Errorf("evidence exceeds size limit")
	}
	b, e := files.ReadFile(path)
	if int64(len(b)) > limit {
		return nil, fmt.Errorf("evidence exceeds size limit")
	}
	return b, e
}

func deliveryEnsureDir(files FileSystem, path string) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return fmt.Errorf("delivery directory must be absolute and clean")
	}
	// Check the existing ancestor before MkdirAll so a pre-existing staging
	// symlink cannot create a directory outside the managed reply area.
	for ancestor := path; ; ancestor = filepath.Dir(ancestor) {
		if _, err := files.Lstat(ancestor); err == nil {
			physical, err := files.EvalSymlinks(ancestor)
			if err != nil {
				return err
			}
			if physical != ancestor {
				return fmt.Errorf("delivery directory ancestor must be physical")
			}
			break
		} else if !os.IsNotExist(err) {
			return err
		}
		if filepath.Dir(ancestor) == ancestor {
			return fmt.Errorf("delivery directory ancestor is missing")
		}
	}
	if e := files.MkdirAll(path, 0700); e != nil {
		return e
	}
	p, e := files.EvalSymlinks(path)
	if e != nil {
		return e
	}
	if p != path {
		return fmt.Errorf("delivery directory must be physical")
	}
	return nil
}

func deliveryWriteOnce(files FileSystem, path string, b []byte) error {
	if _, e := files.Lstat(path); e == nil {
		old, e := deliveryReadFile(files, path, int64(len(b))+1)
		if e != nil {
			return e
		}
		if !bytes.Equal(old, b) {
			return fmt.Errorf("immutable delivery artifact differs: %s", path)
		}
		return nil
	} else if !os.IsNotExist(e) {
		return e
	}
	p, e := files.EvalSymlinks(filepath.Dir(path))
	if e != nil {
		return e
	}
	if p != filepath.Dir(path) {
		return fmt.Errorf("delivery artifact parent must be physical")
	}
	f, e := files.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, fs.FileMode(0600))
	if e != nil {
		return e
	}
	n, e := f.Write(b)
	if e == nil && n != len(b) {
		e = fmt.Errorf("short delivery artifact write")
	}
	if e == nil {
		e = f.Sync()
	}
	ce := f.Close()
	if e != nil {
		return e
	}
	return ce
}
