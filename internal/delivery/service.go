package delivery

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type cwdFiles struct {
	workspace.FileSystem
	cwd string
}

func (f cwdFiles) Getwd() (string, error) { return f.cwd, nil }

type workspaceObserver struct{ dependencies workspace.Dependencies }

func (o workspaceObserver) ObserveContaining() (workflowhandoff.WorkspaceSnapshot, error) {
	v, e := workspace.ObserveContaining(o.dependencies)
	if e != nil {
		return workflowhandoff.WorkspaceSnapshot{}, e
	}
	return o.ObserveRoot(v.Root)
}
func (o workspaceObserver) ObserveRoot(root string) (workflowhandoff.WorkspaceSnapshot, error) {
	v, e := workspace.ObserveRoot(o.dependencies, root)
	if e != nil {
		return workflowhandoff.WorkspaceSnapshot{}, e
	}
	p, r, e := o.dependencies.Projects.Snapshot(root)
	return workflowhandoff.WorkspaceSnapshot{Observation: v, Projects: p, Repositories: r}, e
}

func (s *Service) context(cwd string) (taskrun.Dependencies, string, error) {
	d := s.dependencies
	abs, e := filepath.Abs(cwd)
	if e != nil {
		return d, "", e
	}
	physical, e := filepath.EvalSymlinks(abs)
	if e != nil {
		return d, "", e
	}
	d.Workspace.Files = cwdFiles{d.Workspace.Files, physical}
	d.CWD = func() (string, error) { return physical, nil }
	d.Workflow.TaskWorkspace = &d.Workspace
	d.Workflow.Workspace = workspaceObserver{d.Workspace}
	o, e := workspace.ObserveContaining(d.Workspace)
	return d, o.Root, e
}

func validateRegistration(r Registration) error {
	if r.Kind != "ply.delivery.registration" || r.SchemaVersion != 1 {
		return fmt.Errorf("registration requires kind ply.delivery.registration and schema_version 1")
	}
	if strings.TrimSpace(r.PublicationKey) == "" || len(r.PublicationKey) > 256 || strings.ContainsAny(r.PublicationKey, "\x00\r\n") {
		return fmt.Errorf("a nonempty publication_key of at most 256 bytes is required")
	}
	if _, e := workspace.ParseTaskID(string(r.TaskID)); e != nil {
		return e
	}
	if _, e := workspace.ParseTaskResultID(string(r.TaskResultID)); e != nil {
		return e
	}
	if r.WorkflowRunID == "" || len(r.WorkflowRunID) > 128 || strings.ContainsAny(r.WorkflowRunID, "/\\\x00\r\n") {
		return fmt.Errorf("workflow_run_id is required")
	}
	if r.PredecessorID != "" && !idPattern.MatchString(r.PredecessorID) {
		return fmt.Errorf("invalid predecessor_id")
	}
	if len(r.Title) > 256 || strings.ContainsAny(r.Title, "\x00\r\n") || len(r.Body) > 65536 || strings.ContainsRune(r.Body, 0) {
		return fmt.Errorf("invalid PR title or body")
	}
	return validateMetadata(r.Metadata)
}

func registryCandidate(d taskrun.Dependencies, root string, taskID workspace.TaskID, resultID workspace.TaskResultID) (workspace.WorkItemRegistry, workspace.TaskRecord, workspace.TaskResultRecord, workspace.TaskSpecBasis, error) {
	r, e := d.Workspace.WorkItems.Snapshot(root)
	if e != nil {
		return r, workspace.TaskRecord{}, workspace.TaskResultRecord{}, workspace.TaskSpecBasis{}, e
	}
	var task workspace.TaskRecord
	var result workspace.TaskResultRecord
	var basis workspace.TaskSpecBasis
	for _, v := range r.Tasks {
		if v.ID == taskID {
			task = v
		}
	}
	for _, v := range r.TaskResults {
		if v.ID == resultID {
			result = v
		}
	}
	for _, v := range r.TaskResultSpecBindings {
		if v.TaskResultID == resultID {
			basis = v.Basis
		}
	}
	if task.ID == "" || task.Worktree == nil || result.ID == "" || result.TaskID != task.ID || result.TaskWorktreeID != task.Worktree.ID || result.ProjectID != task.ProjectID || result.RepoID != task.RepoID || result.GitCommonDir != task.GitCommonDir || result.SourceLocator != task.Worktree.Locator || result.SourceRef != task.Worktree.Ref {
		return r, task, result, basis, fmt.Errorf("TaskResult does not bind the registered Task, source and repository")
	}
	if basis.TaskID != taskID || basis.Spec.Revision < 1 {
		return r, task, result, basis, fmt.Errorf("TaskResult has no immutable result-bound Spec")
	}
	return r, task, result, basis, nil
}

func checkEvidence(d taskrun.Dependencies, result workspace.TaskResultRecord, basis workspace.TaskSpecBasis) error {
	if result.TechnicalGate != "passed" && result.TechnicalGate != "good_enough_with_known_debt" {
		return fmt.Errorf("candidate is not technically qualified")
	}
	if len(result.VerifierResults) == 0 || len(result.Review.OpenActionableFindings) != 0 {
		return fmt.Errorf("candidate needs verifier results and review without open actionable findings")
	}
	reader := d.Workspace.HandoffEvidence
	if reader == nil {
		reader = workflowhandoff.NewTaskHandoffEvidenceReader(d.Workflow)
	}
	ev, e := reader.ReadTaskEvidence(workspace.TaskHandoffEvidenceRequest{ActivityID: result.ActivityID, RunID: result.RunID, HandoffID: result.HandoffID, HandoffLocator: result.HandoffLocator, HandoffSHA256: result.HandoffSHA256, StartReceiptID: result.StartReceiptID, StartReceiptLocator: result.StartReceiptLocator, StartReceiptSHA256: result.StartReceiptSHA256, TerminalResultID: result.TerminalResultID, TerminalResultLocator: result.TerminalResultLocator, TerminalResultSHA256: result.TerminalResultSHA256, InspectionSHA256: result.InspectionSHA256})
	if e != nil {
		return fmt.Errorf("immutable candidate evidence: %w", e)
	}
	if !ev.SchemaValid || !ev.DigestValid || !ev.LifecycleValid || !ev.BindingValid || !ev.CapabilityValid || !ev.PrincipalSessionValid || !ev.PolicyValid || !ev.TaskSpecValid || !ev.TaskRequirementsValid || ev.AcceptedStartOutcome != "started" || ev.ReportedOutcome != "complete" || result.ReportedOutcome != ev.ReportedOutcome || !reflect.DeepEqual(ev.TaskSpecBasis, &basis) || !reflect.DeepEqual(ev.VerifierResults, result.VerifierResults) || !reflect.DeepEqual(ev.Review, result.Review) || ev.ResultOID != result.ResultOID || ev.ResultTree != result.ResultTree || ev.TargetRef != result.SourceRef || ev.TargetWorktree != result.SourceLocator {
		return fmt.Errorf("TaskResult differs from revalidated native evidence or Spec")
	}
	if len(ev.ExpectedVerifierIDs) != len(result.VerifierResults) {
		return fmt.Errorf("native verifier set differs")
	}
	for _, id := range ev.ExpectedVerifierIDs {
		n := 0
		for _, v := range result.VerifierResults {
			if v.VerifierID == id && v.Outcome == "passed" {
				n++
			}
		}
		if n != 1 {
			return fmt.Errorf("required verifier %s has no exact passing candidate result", id)
		}
	}
	for _, a := range result.Artifacts {
		if e = rehash(a.Locator, a.SHA256, a.SizeBytes); e != nil {
			return fmt.Errorf("artifact %s: %w", a.ArtifactID, e)
		}
	}
	return nil
}

func rehash(path, expected string, size int64) error {
	raw, e := readBounded(path, 64<<20)
	if e != nil {
		return e
	}
	if int64(len(raw)) != size || hash(raw) != expected {
		return fmt.Errorf("preserved evidence bytes changed: %s", path)
	}
	return nil
}

func sourceUnchanged(d taskrun.Dependencies, result workspace.TaskResultRecord) error {
	obs, e := d.Workspace.WorkGit.ObserveWorktree(result.SourceLocator)
	if e != nil {
		return e
	}
	if !obs.Clean || !obs.InventoryMatch || obs.GitCommonDir != result.GitCommonDir || obs.Ref != result.SourceRef || obs.OID != result.ResultOID || obs.Tree != result.ResultTree {
		return fmt.Errorf("candidate source is dirty, moved or differs from its recorded commit/tree")
	}
	return nil
}

// The immutable Spec supplies the original agreement. An explicit native
// acceptance amendment may affect its policy, but cannot rewrite this basis.
func originalDeliveryAgreement(d taskrun.Dependencies, root string, r workspace.WorkItemRegistry, task workspace.TaskRecord, result workspace.TaskResultRecord, basis workspace.TaskSpecBasis) (*workspace.DeliveryAgreement, error) {
	eval, e := workspace.ReadHistoricalTaskSpec(d.Workspace, root, basis)
	if e != nil {
		return nil, e
	}
	agreement, e := workspace.DeliveryAgreementFromSpec(eval.Spec)
	if e != nil {
		return nil, e
	}
	if agreement == nil {
		// Only the historical local Epic contract is retained when no new field
		// exists. The immutable native owner request must independently agree.
		for _, epic := range r.Epics {
			if epic.ID == task.ParentEpicID {
				for _, binding := range epic.RepoBindings {
					if binding.RepoID == task.RepoID {
						agreement = &workspace.DeliveryAgreement{SchemaVersion: 1, Mode: workspace.DeliveryLocalEpic, ProjectID: task.ProjectID, RepoID: task.RepoID, EpicID: task.ParentEpicID, SourceRef: result.SourceRef, TargetRef: binding.Worktree.Ref, TargetWorktree: binding.Worktree.Locator}
					}
				}
			}
		}
	}
	if agreement == nil {
		return nil, fmt.Errorf("no immutable delivery agreement or historical local Epic return context")
	}
	if e = workspace.ValidateDeliveryAgreement(*agreement); e != nil {
		return nil, e
	}
	if agreement.SourceRef != result.SourceRef || agreement.ProjectID != task.ProjectID || agreement.RepoID != task.RepoID || agreement.EpicID != task.ParentEpicID {
		return nil, fmt.Errorf("result-bound agreement differs from the Task source and scope")
	}
	return agreement, nil
}

func manifestAuthorityMatches(m Manifest, auth taskrun.DeliveryAuthority) error {
	if auth.Authorization == nil || auth.RequestSHA256 != m.RequestSHA256 || auth.MandateSHA256 != m.MandateSHA256 || auth.PreparationID != m.PreparationID || auth.ExpectedParentOID != m.ExpectedParentOID || !reflect.DeepEqual(auth.Agreement, m.Agreement) || !reflect.DeepEqual(auth.Authorization.AcceptanceAmendment, m.AcceptanceAmendment) {
		return fmt.Errorf("preserved native delivery authority or selected acceptance policy changed")
	}
	return nil
}

func (s *Service) buildManifest(d taskrun.Dependencies, root string, in Registration) (Manifest, error) {
	var out Manifest
	r, task, result, basis, e := registryCandidate(d, root, in.TaskID, in.TaskResultID)
	if e != nil {
		return out, e
	}
	agreement, e := originalDeliveryAgreement(d, root, r, task, result, basis)
	if e != nil {
		return out, e
	}
	if e = checkEvidence(d, result, basis); e != nil {
		return out, e
	}
	if e = sourceUnchanged(d, result); e != nil {
		return out, e
	}
	auth, e := taskrun.ValidateDeliveryAuthority(d, root, in.WorkflowRunID, result, *agreement)
	if e != nil {
		return out, e
	}
	// Only the native accepted run can supply the effective policy. The original
	// Spec remains in the manifest and is re-read before each authority check.
	agreement = &auth.Agreement
	if agreement.Mode != workspace.DeliveryPullRequest && (in.Metadata != nil || in.Title != "" || in.Body != "") {
		return out, fmt.Errorf("PR title, body and people metadata do not belong to local delivery")
	}
	out = Manifest{Workspace: root, TaskID: task.ID, EpicID: task.ParentEpicID, ProjectID: task.ProjectID, RepoID: task.RepoID, Spec: basis, TaskResult: result, TaskResultSHA256: digest(result), Agreement: *agreement, AcceptanceAmendment: auth.Authorization.AcceptanceAmendment, WorkflowRunID: in.WorkflowRunID, RequestSHA256: auth.RequestSHA256, MandateSHA256: auth.MandateSHA256, PreparationID: auth.PreparationID, OwnerClaim: auth.OwnerClaim, ExpectedParentOID: auth.ExpectedParentOID, RequiredGates: []string{"native_technical_qualification", "preserved_review", "exact_candidate_human_pass", "frozen_delivery_authority", "exact_target"}}
	if agreement.AutomaticAcceptance() {
		out.RequiredGates[2] = "exact_candidate_automatic_pass"
		if agreement.Acceptance.RequireHumanQA || agreement.HumanOwnedIntegration() {
			out.RequiredGates = append(out.RequiredGates, "exact_candidate_human_pass")
		}
	}
	// Independent publication keys still serialize on the actual candidate and
	// effect target. Neither workflow ID nor TaskResult ID creates a new effect.
	out.EffectKey = digest(struct{ Common, OID, Tree, Mode, TargetRef, TargetWorktree, Repository string }{result.GitCommonDir, result.ResultOID, result.ResultTree, agreement.Mode, agreement.TargetRef, agreement.TargetWorktree, strings.ToLower(agreement.GitHubRepository)})
	return out, nil
}

func (s *Service) Register(cwd, file string) (Receipt, error) {
	var out Receipt
	var in Registration
	if e := readJSON(file, &in); e != nil {
		return out, e
	}
	if e := validateRegistration(in); e != nil {
		return out, e
	}
	d, root, e := s.context(cwd)
	if e != nil {
		return out, e
	}
	id := "dlv_" + shortHash(in.PublicationKey)
	e = withLock(root, func() error {
		old, e := readReceipt(root, id)
		if e == nil {
			if old.RegistrationSHA256 != digest(in) {
				return fmt.Errorf("delivery publication_key already exists with different content; use a new key and predecessor_id")
			}
			out = old
			return nil
		}
		if !os.IsNotExist(e) {
			return e
		}
		manifest, e := s.buildManifest(d, root, in)
		if e != nil {
			return e
		}
		if in.PredecessorID != "" {
			previous, e := readReceipt(root, in.PredecessorID)
			if e != nil {
				return e
			}
			if previous.Manifest.TaskID != manifest.TaskID {
				return fmt.Errorf("predecessor belongs to another Task")
			}
		}
		at := s.now()
		out = Receipt{Kind: "ply.delivery.receipt", SchemaVersion: 1, ID: id, Registration: in, RegistrationSHA256: digest(in), Manifest: manifest, ManifestSHA256: digest(manifest), State: "registered", Reasons: []string{}, NextAction: "Run delivery check, then execute this Delivery ID after exact human pass; no background execution is scheduled.", CreatedAtUTC: at, UpdatedAtUTC: at, Events: []Event{}}
		if manifest.Agreement.AutomaticAcceptance() {
			out.NextAction = "Run delivery check, then execute this Delivery ID after its candidate-bound acceptance gates pass."
		}
		s.event(&out, "delivery.registered", "registered", "Preserved native candidate, evidence, frozen agreement and authority; no delivery effect.", "")
		return s.save(root, &out)
	})
	return out, e
}

func (s *Service) Read(cwd, id string) (Receipt, error) {
	_, root, e := s.context(cwd)
	if e != nil {
		return Receipt{}, e
	}
	return readReceipt(root, id)
}
func (s *Service) List(cwd, taskID, epicID string) (ListResult, error) {
	out := ListResult{Kind: "ply.delivery.list", SchemaVersion: 1, Deliveries: []Receipt{}}
	_, root, e := s.context(cwd)
	if e != nil {
		return out, e
	}
	rows, e := listReceipts(root)
	if e != nil {
		return out, e
	}
	for _, r := range rows {
		if (taskID == "" || string(r.Manifest.TaskID) == taskID) && (epicID == "" || string(r.Manifest.EpicID) == epicID) {
			out.Deliveries = append(out.Deliveries, r)
		}
	}
	return out, nil
}

func (s *Service) gate(d taskrun.Dependencies, root string, r *Receipt) (_ taskrun.DeliveryAuthority, gateErr error) {
	r.Acceptance = nil
	defer func() {
		if gateErr != nil && r.Manifest.Agreement.AutomaticAcceptance() {
			r.Acceptance = &workspace.DeliveryAcceptanceDecision{Mode: "automatic", Outcome: "blocked", Reason: gateErr.Error(), TaskResultID: r.Manifest.TaskResult.ID, ResultOID: r.Manifest.TaskResult.ResultOID, ResultTree: r.Manifest.TaskResult.ResultTree, PolicySHA256: digest(*r.Manifest.Agreement.Acceptance)}
		}
	}()
	var empty taskrun.DeliveryAuthority
	registry, task, result, basis, e := registryCandidate(d, root, r.Manifest.TaskID, r.Manifest.TaskResult.ID)
	if e != nil {
		return empty, e
	}
	if digest(result) != r.Manifest.TaskResultSHA256 || !reflect.DeepEqual(basis, r.Manifest.Spec) {
		return empty, fmt.Errorf("registered native TaskResult or Spec changed")
	}
	original, e := originalDeliveryAgreement(d, root, registry, task, result, basis)
	if e != nil {
		return empty, e
	}
	if e = checkEvidence(d, result, basis); e != nil {
		return empty, e
	}
	if e = sourceUnchanged(d, result); e != nil {
		return empty, e
	}
	auth, e := taskrun.ValidateDeliveryAuthority(d, root, r.Manifest.WorkflowRunID, result, *original)
	if e != nil {
		return empty, e
	}
	if e = manifestAuthorityMatches(r.Manifest, auth); e != nil {
		return empty, e
	}
	// A newer fail/blocked supersedes a prior answer for this exact result. QA
	// for another result, commit or tree can never unlock the candidate.
	r.HumanQA, e = taskrun.LatestDeliveryAcceptanceHumanQA(registry.HumanQARecords, result, &r.Manifest.Agreement)
	if e != nil {
		return empty, e
	}
	if r.HumanQA != nil {
		for _, a := range r.HumanQA.Evidence {
			if e = rehash(a.Locator, a.SHA256, a.SizeBytes); e != nil {
				return empty, fmt.Errorf("human QA evidence: %w", e)
			}
		}
	}
	decision, e := workspace.EvaluateDeliveryAcceptance(&r.Manifest.Agreement, result, r.HumanQA)
	if e != nil {
		return empty, e
	}
	r.Acceptance = &decision
	if decision.Outcome != "pass" {
		r.State = "blocked"
		r.Reasons = []string{decision.Reason}
		r.NextAction = "Resolve the recorded acceptance findings, then verify and review any changed candidate before delivery."
		if r.HumanQA == nil && (!r.Manifest.Agreement.AutomaticAcceptance() || r.Manifest.Agreement.Acceptance.RequireHumanQA) {
			r.State = "awaiting_human"
			r.NextAction = "Obtain and record the actual candidate-bound human answer, then explicitly execute this Delivery ID."
		}
		return auth, nil
	}
	if !reflect.DeepEqual(auth.Candidate.HumanQA, r.HumanQA) {
		return empty, fmt.Errorf("human pass is not the exact answer preserved by the native workflow candidate")
	}
	if auth.Authorization.HumanIntegrationRequired || auth.Authorization.HumanIntegration != nil {
		r.State = "blocked"
		r.Reasons = []string{"human_integration_required: local delivery belongs to the explicitly selected human integration flow."}
		r.NextAction = "Continue the explicit human integration plan for this candidate."
		return auth, nil
	}
	r.State = "ready"
	r.Reasons = []string{}
	r.NextAction = "Explicitly execute this Delivery ID to perform only its frozen delivery mode."
	return auth, nil
}

func (s *Service) Check(cwd, id string) (Receipt, error) {
	d, root, e := s.context(cwd)
	if e != nil {
		return Receipt{}, e
	}
	r, e := readReceipt(root, id)
	if e != nil {
		return r, e
	}
	if r.State == "delivered" {
		return r, nil
	}
	if r.Manifest.Agreement.AutomaticAcceptance() {
		observed, found, err := taskrun.ObserveAutomaticLocalDelivery(d, root, r.Manifest.WorkflowRunID, r.Manifest.TaskResult.ID)
		if found {
			if err != nil && !observed.Completed {
				r.State, r.Reasons, r.NextAction = "blocked", []string{err.Error()}, "Restore the exact native integration evidence before resuming closeout."
				return r, nil
			}
			r.State, r.NativeClosed = "pending_closeout", false
			r.Reasons = []string{"Local integration, queue closure and Epic base update are observed; native Task closeout remains pending."}
			r.NextAction = "Execute this Delivery ID to resume native closeout while retaining the source worktree and branch. The Git integration will not be repeated."
			if observed.Completed {
				r.Local = &LocalReceipt{TargetRef: r.Manifest.Agreement.TargetRef, TargetWorktree: r.Manifest.Agreement.TargetWorktree, BeforeOID: r.Manifest.ExpectedParentOID, AfterOID: observed.Integration.Readback.ParentOID, Integration: observed.Integration.Readback, Base: observed.Base, Queue: observed.Queue, ObservedAtUTC: s.now(), Closeout: observed.Closeout}
			}
			if err != nil {
				r.Reasons = append(r.Reasons, err.Error())
			} else if observed.Closeout != nil && observed.Closeout.State == "complete" && observed.Closeout.LifecycleCompleted {
				r.State, r.NativeClosed, r.Reasons = "ready", true, []string{}
				r.NextAction = "Execute this Delivery ID to record the already observed native completion. The worktree and branch are retained; no integration will be repeated."
			}
			return r, nil
		}
		if err != nil {
			r.State, r.Reasons, r.NextAction = "blocked", []string{err.Error()}, "Restore the preserved automatic delivery evidence before continuing."
			return r, nil
		}
	}
	auth, e := s.gate(d, root, &r)
	if e != nil {
		r.State = "blocked"
		r.Reasons = []string{e.Error()}
		r.NextAction = "Inspect the precise conflict; mode or target changes require a new authorized run and Delivery."
		return r, nil
	}
	if r.State != "ready" {
		return r, nil
	}
	if r.Manifest.Agreement.Mode == workspace.DeliveryPullRequest {
		if e = s.checkPRReadiness(root, &r); e != nil {
			r.State = "blocked"
			r.Reasons = []string{e.Error()}
		}
	} else {
		in := integrationInput(r, auth)
		preview, e := workspace.CheckTaskIntegration(d.Workspace, in)
		if e != nil {
			r.State = "blocked"
			r.Reasons = []string{e.Error()}
		} else if preview.Readback.Classification != "ready" && preview.Readback.RecoveryStatus != "complete" {
			r.State = "blocked"
			r.Reasons = []string{"local integration: " + preview.Readback.Classification}
		}
	}
	return r, nil
}

func (s *Service) checkPRReadiness(root string, r *Receipt) error {
	fx, e := readEffect(root, r.Manifest.EffectKey)
	if e != nil && !os.IsNotExist(e) {
		return e
	}
	o, e := s.options.PR.Observe(prTarget(*r))
	if e != nil {
		if fx.PushStarted || fx.CreateStarted || fx.PR != nil {
			readinessUnknown(r, e.Error())
			return nil
		}
		return e
	}
	p, e := validateObservation(prTarget(*r), o)
	if e != nil {
		return e
	}
	if fx.CreateStarted && p == nil {
		readinessUnknown(r, "PR creation is reserved but its result is not observed.")
		return nil
	}
	if fx.PushStarted && o.RemoteOID != r.Manifest.TaskResult.ResultOID {
		readinessUnknown(r, "Reserved source push is not observed at the exact candidate.")
		return nil
	}
	if fx.PR != nil && (p == nil || p.Number != fx.PR.Number || p.URL != fx.PR.URL) {
		readinessUnknown(r, "Previously observed PR no longer matches the preserved receipt.")
		return nil
	}
	if e == nil && o.RemoteOID != "" && o.RemoteOID != r.Manifest.TaskResult.ResultOID {
		var ff bool
		ff, e = s.options.PR.CanFastForward(prTarget(*r), o.RemoteOID)
		if e == nil && !ff {
			e = fmt.Errorf("source publication would require force-push")
		}
	}
	return e
}

func readinessUnknown(r *Receipt, reason string) {
	r.State = "unknown_effect"
	r.Reasons = []string{reason}
	r.NextAction = "Restore access and observe the reserved effect, then explicitly execute this Delivery ID to resume; no blind publication is allowed."
}

func integrationInput(r Receipt, a taskrun.DeliveryAuthority) workspace.TaskIntegrationInput {
	return workspace.TaskIntegrationInput{TaskID: r.Manifest.TaskID, TaskResultID: r.Manifest.TaskResult.ID, HumanQARecordID: deliveryQAID(r), ExpectedResultOID: r.Manifest.TaskResult.ResultOID, ExpectedParentOID: r.Manifest.ExpectedParentOID, DeliveryOwner: &workspace.DeliveryIntegrationOwner{RunID: a.RunID, RequestSHA256: a.RequestSHA256, ActorClaim: a.OwnerClaim, PreparationID: a.PreparationID}, DeliveryAuthorization: a.Authorization}
}

func deliveryQAID(r Receipt) workspace.HumanQARecordID {
	if r.HumanQA == nil || r.HumanQA.TaskResultID != r.Manifest.TaskResult.ID {
		return ""
	}
	return r.HumanQA.ID
}
func prTarget(r Receipt) PRTarget {
	a := r.Manifest.Agreement
	return PRTarget{Worktree: r.Manifest.TaskResult.SourceLocator, Remote: a.Remote, Repository: a.GitHubRepository, SourceRef: a.SourceRef, BaseRef: a.TargetRef, OID: r.Manifest.TaskResult.ResultOID, EffectKey: r.Manifest.EffectKey}
}
