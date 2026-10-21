package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type ReportStartupStatus struct {
	SourceSHA256      string        `json:"source_sha256"`
	EffectResolutions []FileBinding `json:"effect_resolutions"`
	Protocol          string        `json:"protocol"`
	Receipt           *FileBinding  `json:"receipt"`
	ProviderStop      *FileBinding  `json:"provider_stop"`
	Stop              *FileBinding  `json:"stop"`
	Retirement        *FileBinding  `json:"retirement"`
	SandboxCoverage   string        `json:"sandbox_coverage"`
	RuntimeDecision   string        `json:"runtime_decision"`
}

const ReportLegacyAcceptanceV1 = "legacy_acceptance_v1"

// Supported native Claude permission modes for interactive Herdr report Runs.
const (
	ClaudePermissionManual = "manual"
	ClaudePermissionAuto   = "auto"
)

func reportReceipt(r WorkflowRequest) bool {
	return r.SchemaVersion == 1 && r.ReturnMode == "reviewed_report_only" && r.Delivery == nil && r.StartupProtocol == MandateReceiptV1
}

// Report receipts reuse the native identity/input/start binding. No delivery
// contract or claim of sandbox coverage is introduced.
func WorkflowReportReceipt(d Dependencies, root, id, contextPath string, in DeliveryBeginInput) (WorkflowRun, error) {
	s, err := workflowRead(root, id)
	if err != nil {
		return WorkflowRun{}, err
	}
	if !reportReceipt(s.Request) {
		return WorkflowRun{}, workflowError(4, "report receipt differs from frozen contract")
	}
	if in.PermissionConfirmed || in.NotificationRoute != "" {
		return WorkflowRun{}, workflowError(2, "report receipt cannot confirm permissions or acquire delivery notification authority")
	}
	if err = workflowCallback(d, s, contextPath); err != nil {
		return WorkflowRun{}, err
	}
	if in.Provider != s.Request.Runtime.Provider || in.NativeSession == "" || in.NativeSession != s.Result.Transport.AgentSessionID {
		return WorkflowRun{}, workflowError(4, "provider or native session differs from bound report owner")
	}
	if s.ReportStop != nil || s.ReportProviderStop != nil {
		return WorkflowRun{}, workflowError(4, "preserved stop blocks work; coordinator must retire this owner")
	}
	dir := filepath.Join(s.Result.Paths.RunRoot, "receipt", digest(in)[7:])
	if err = privateDir(dir); err != nil {
		return WorkflowRun{}, err
	}
	tmp := filepath.Join(dir, "tmp")
	if err = privateDir(tmp); err != nil {
		return WorkflowRun{}, err
	}
	reads, writes := in.ReadRoots, in.WriteRoots
	if reads == nil {
		reads = []string{}
	}
	if writes == nil {
		writes = []string{}
	}
	sandbox, _ := Canonical(map[string]any{"read_roots": reads, "write_roots": writes, "temp_root": tmp, "matches_contract": nil})
	claim := RuntimeClaim{RuntimeID: ptr(in.Provider), NativeSessionID: ptr(in.NativeSession)}
	if in.Model != "" {
		claim.ModelID = ptr(in.Model)
	}
	if in.Profile != "" {
		claim.ProfileID = ptr(in.Profile)
	}
	if in.Policy != "" {
		evidence, e := workflowKeep(d, filepath.Join(dir, "runtime-observation.json"), in)
		if e != nil {
			return WorkflowRun{}, e
		}
		claim.EffectivePolicySHA256 = ptr(evidence.SHA256)
	}
	a := Acceptance{Envelope: workflowEnv("run-acceptance"), RunID: id, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID, RuntimeClaim: claim, Sandbox: sandbox, Acceptance: "started", Issues: json.RawMessage(`[]`)}
	binding, err := workflowKeep(d, filepath.Join(dir, "acceptance.json"), a)
	if err != nil {
		return WorkflowRun{}, err
	}
	return WorkflowAccept(d, root, id, contextPath, binding.Locator)
}

func workflowPreservePreStartStop(d Dependencies, s *workflowState, r WorkflowReport) error {
	if r.TechnicalAssessment.Gate != "unknown" || r.BudgetUsage != nil && (r.BudgetUsage.InitialExecutionStarted || r.BudgetUsage.CorrectionRounds != 0 || r.BudgetUsage.EnvironmentMeasures != 0) {
		return workflowError(4, "pre-start observation cannot claim work or technical qualification")
	}
	for _, raw := range []json.RawMessage{r.ObservedEffects, r.VerifierResults, r.Artifacts} {
		var entries []json.RawMessage
		if json.Unmarshal(raw, &entries) != nil || len(entries) != 0 {
			return workflowError(4, "pre-start stop cannot claim effects, tests or artifacts")
		}
	}
	if s.ReportRetirement != nil {
		return workflowError(4, "owner_retired")
	}
	if s.ReportStop != nil {
		if s.ReportStop.SHA256 == digest(r) {
			_, err := workflowBound(*s.ReportStop, 1<<20)
			return err
		}
		return workflowError(4, "pre-start stop observation is immutable")
	}
	binding, err := workflowKeep(d, filepathForRound(*s, "pre-start-stop.json"), r)
	if err != nil {
		return err
	}
	s.ReportStop = &binding
	s.Result.Round.State = "stopped"
	s.Result.NextAction = WorkflowAction{"coordinator", "Pre-start stop preserved without work acceptance. Inspect the stop; use workflow run retire for an absent clean owner and a fresh replacement mandate."}
	return nil
}

func workflowReportRetirement(s workflowState) (*DeliveryRetirement, error) {
	var raw []byte
	var err error
	if s.ReportRetirement == nil {
		raw, err = readFile(filepath.Join(s.Result.Paths.RunRoot, "report-retirement.json"), 1<<20, true)
		if os.IsNotExist(err) {
			return nil, nil
		}
	} else {
		raw, err = workflowBound(*s.ReportRetirement, 1<<20)
	}
	if err != nil {
		return nil, err
	}
	var r DeliveryRetirement
	if err = decode(raw, 1<<20, &r); err != nil {
		return nil, err
	}
	if r.RunID != s.Result.RunID || r.RequestSHA256 != s.Result.RequestSHA256 || r.NativeSession != s.Result.Transport.AgentSessionID || r.Kind != "ply.workflow.owner-retirement" || r.SchemaVersion != 1 {
		return nil, workflowError(4, "report retirement binding differs")
	}
	return &r, nil
}

func workflowRetireReport(d Dependencies, s *workflowState, in DeliveryRetireInput) error {
	if s.Request.SchemaVersion != 1 || s.Request.ReturnMode != "reviewed_report_only" || headlessRun(s.Request) {
		return workflowError(4, "retirement requires an interactive report or delivery owner")
	}
	if prior, err := workflowReportRetirement(*s); err != nil {
		return err
	} else if prior != nil {
		if prior.Input != in {
			return workflowError(4, "retirement input is immutable")
		}
		return nil
	}
	// Installed CLI/provider launch bytes are historical here. This coordinator
	// action never launches or impersonates that owner. Native artifacts, current
	// Herdr inspection and local processes still must be verifiable.
	if err := workflowFreshArtifacts(d, *s, false, false); err != nil {
		return err
	}
	facts, err := workflowhandoff.ReadTaskRunReturnFacts(d.Workflow, s.Result.Handoff.Locator)
	if err != nil {
		return err
	}
	if facts.Start != nil || facts.Terminal != nil {
		return workflowError(4, "native start/terminal exists; reconcile prior effects before pre-work retirement")
	}
	if s.StartSHA256 != nil || len(s.Records) != 0 || s.Result.Transport.AgentSessionID == "" || s.Phase != "following" && s.Phase != "prompt_attempted" && s.Phase != "prompt_unknown" {
		return workflowError(4, "report retirement requires a bound pre-work owner; reconcile existing work or unknown startup effects")
	}
	if err := workflowReportEffectsPermitRetirement(d, *s); err != nil {
		return err
	}

	clean := func() error {
		target, err := d.Workspace.WorkGit.ObserveWorktree(s.Observed.Target.WorktreeLocator)
		if err != nil {
			return err
		}
		if !target.Clean || target.OID != s.Observed.Target.OID || target.Ref != s.Observed.Target.Ref || target.GitCommonDir != s.Observed.Target.GitCommonDir {
			return workflowError(4, "report source changed; reconcile actual effects before retirement")
		}
		return nil
	}
	if err := clean(); err != nil {
		return err
	}
	// Report retirement observes absence; it does not close a potentially live
	// legacy writer or repeat a denied owner callback.
	writers := (taskWriters{d, *s}).cleanup()
	if writers.state != taskWritersAbsent {
		return writers.err
	}
	if err := clean(); err != nil {
		return err
	}
	r := DeliveryRetirement{Kind: "ply.workflow.owner-retirement", SchemaVersion: 1, RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, NativeSession: s.Result.Transport.AgentSessionID, Input: in, Absence: writers.absence, WriterLiveness: writers.liveness, RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
	binding, err := workflowKeep(d, filepath.Join(s.Result.Paths.RunRoot, "report-retirement.json"), r)
	if err != nil {
		return err
	}
	if err = d.fault("report_after_retirement"); err != nil {
		return err
	}
	s.ReportRetirement = &binding
	s.Result.Round.State, s.Result.FinalReturn.State = "stopped", "retired"
	s.Result.NextAction = WorkflowAction{"coordinator", "Report owner retired with history preserved. Release only its exact preparation; replacement requires a fresh mandate."}
	return nil
}

func workflowReleaseReportQueue(d Dependencies, s workflowState) error {
	p, err := workspace.ShowTaskPreparation(d.Workspace, s.Request.PreparationID)
	if err != nil {
		return err
	}
	if p.Plan == nil || p.Plan.WorktreePath != s.Observed.Target.WorktreeLocator || p.Plan.ParentOID != s.Observed.Target.OID {
		return workflowError(4, "retired report preparation binding differs")
	}
	if p.Disposition == "released" {
		return nil
	}
	t := p.Plan.Target
	target := workspace.QueueTargetInput{ProjectID: t.ProjectID, RepoID: t.RepoID, EpicID: t.EpicID}
	q, err := workspace.ListTaskQueue(d.Workspace, target, false)
	if err != nil {
		return err
	}
	if q.Current == nil || q.Current.PreparationID != s.Request.PreparationID {
		return workflowError(4, "retired report does not own current queue slot; no release performed")
	}
	_, err = workspace.CloseTaskQueue(d.Workspace, target, s.Request.PreparationID, q.Revision, "Retired pre-work report owner", "release")
	return err
}

// Coordinator observations preserve denied callback transport without speaking
// for the owner. Report owners with a stop can retire; receipt never clears it.
func workflowObserveReportStop(d Dependencies, s *workflowState, in DeliveryStopInput, source FileBinding, raw []byte) error {
	if s.Request.SchemaVersion != 1 || s.Request.ReturnMode != "reviewed_report_only" || s.ReportRetirement != nil {
		return workflowError(4, "report owner is closed or unsupported")
	}
	if in.Provider != s.Request.Runtime.Provider || in.NativeSession == "" || in.NativeSession != s.Result.Transport.AgentSessionID {
		return workflowError(4, "stop provider/session differs")
	}
	path := filepath.Join(s.Result.Paths.RunRoot, "report-provider-stop.json")
	if s.ReportProviderStop != nil {
		b, err := workflowBound(*s.ReportProviderStop, 2<<20)
		if err != nil {
			return err
		}
		var prior DeliveryStopRecord
		if err = decode(b, 2<<20, &prior); err != nil {
			return err
		}
		if !equal(prior.Input, in) || prior.Source != source {
			return workflowError(4, "report provider stop is immutable")
		}
		return nil
	}
	evidencePath := filepath.Join(s.Result.Paths.RunRoot, "report-provider-stop-evidence")
	if err := d.writeOnce(evidencePath, raw); err != nil {
		return err
	}
	record := DeliveryStopRecord{Kind: "ply.workflow.coordinator-stop", SchemaVersion: 1, RunID: s.Result.RunID, Input: in, Source: source, Evidence: FileBinding{evidencePath, source.SHA256}, ObservedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
	binding, err := workflowKeep(d, path, record)
	if err != nil {
		return err
	}
	s.ReportProviderStop = &binding
	s.Result.Round.State = "stopped"
	s.Result.NextAction = workflowReportStopAction(*s)
	return nil
}

func workflowLoadReportObservations(s *workflowState) error {
	if deliveryRun(s.Request) {
		return nil
	}
	for _, item := range []struct {
		path     string
		slot     **FileBinding
		provider bool
	}{
		{filepathForRound(*s, "pre-start-stop.json"), &s.ReportStop, false},
		{filepath.Join(s.Result.Paths.RunRoot, "report-provider-stop.json"), &s.ReportProviderStop, true},
	} {
		raw, err := readFile(item.path, 2<<20, true)
		if os.IsNotExist(err) && *item.slot == nil {
			continue
		}
		if err != nil {
			return err
		}
		binding := FileBinding{item.path, hash(raw)}
		if *item.slot != nil && **item.slot != binding {
			return workflowError(4, "report stop binding changed")
		}
		if item.provider {
			var r DeliveryStopRecord
			if err = decode(raw, 2<<20, &r); err != nil {
				return err
			}
			if r.Kind != "ply.workflow.coordinator-stop" || r.SchemaVersion != 1 || r.RunID != s.Result.RunID || r.Input.Provider != s.Request.Runtime.Provider || r.Input.NativeSession != s.Result.Transport.AgentSessionID {
				return workflowError(4, "report provider stop binding differs")
			}
			if _, err = workflowBound(r.Evidence, 1<<20); err != nil {
				return err
			}
		} else {
			r, err := workflowParseReport(raw)
			if err != nil {
				return err
			}
			if err = workflowClaim(*s, r.RunID, r.RequestSHA256, r.SessionID); err != nil {
				return err
			}
			if r.Round != s.Result.Round.Number || r.Outcome != "blocked" && r.Outcome != "stopped" {
				return workflowError(4, "pre-start stop binding differs")
			}
		}
		*item.slot = &binding
	}
	return nil
}

// Called under the existing workspace snapshot lock: observe the validated
// append-only registry without reacquiring the workspace lock.
func workflowReportQueueReleased(d Dependencies, root, id string) (bool, error) {
	r, err := d.Workspace.WorkItems.Snapshot(root)
	if err != nil {
		return false, err
	}
	queueID := ""
	for _, p := range r.TaskPreparations {
		if p.ID == id {
			queueID = p.Plan.QueueID
		}
	}
	if queueID == "" {
		return false, workflowError(4, "retired report preparation is missing")
	}
	for _, event := range r.TaskQueueEvents {
		if event.QueueID == queueID && event.Kind == "release" && event.Request["preparation_id"] == id {
			return true, nil
		}
	}
	return false, nil
}

func workflowReportStopped(s workflowState) bool {
	return !deliveryRun(s.Request) && (s.ReportStop != nil || s.ReportProviderStop != nil || s.ReportRetirement != nil)
}
func workflowReportStopAction(s workflowState) WorkflowAction {
	if s.ReportRetirement != nil {
		return WorkflowAction{"coordinator", "Report owner retired with history preserved; replacement needs a fresh prepared Task mandate."}
	}
	if s.ReportStop == nil && s.StartSHA256 != nil {
		return WorkflowAction{"coordinator", "Provider stop observed after mandate receipt is preserved; owner work may exist and pre-work retirement does not apply. Reconcile actual effects with workflow run resolve-stop and preserve the source; do not replay the refused callback."}
	}
	if s.ReportStop == nil {
		return WorkflowAction{"coordinator", "Provider stop preserved without acceptance. Preserve the refusal; retire only after known effects and positive writer absence, then use a fresh mandate."}
	}
	return WorkflowAction{"coordinator", "Preserved pre-start stop blocks owner work and successful round review. Reconcile any unknown effects with workflow run resolve-stop, then retire only after source and writer-absence checks; do not replay the refused callback."}
}
