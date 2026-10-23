package taskrun

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// A disposition is a bounded decision to proceed under retained uncertainty.
// Neither declared resources nor source inspection establish runtime absence.
// Records live beside the Run, keeping strict historical state decoders valid.
type DeliveryDispositionInput struct {
	Kind          string                         `json:"kind"`
	SchemaVersion int                            `json:"schema_version"`
	EventID       string                         `json:"event_id"`
	RunID         string                         `json:"run_id"`
	RequestSHA256 string                         `json:"request_sha256"`
	MandateSHA256 string                         `json:"mandate_sha256"`
	NativeSession string                         `json:"native_session"`
	SourceOID     string                         `json:"source_oid"`
	SourceTree    string                         `json:"source_tree"`
	Retained      []DeliveryDispositionResidual  `json:"retained"`
	Operations    []DeliveryDispositionOperation `json:"operations"`
}

type DeliveryDispositionResidual struct {
	Stop     FileBinding `json:"stop"`
	Unknowns []string    `json:"unknowns"`
	// Concrete possible overlaps, not an exhaustive inventory of old effects.
	Conflicts  []string      `json:"conflicts"`
	Assessment string        `json:"assessment"`
	Provenance string        `json:"provenance"`
	Evidence   []FileBinding `json:"evidence"`
}

type DeliveryDispositionOperation struct {
	Operation  string   `json:"operation"`
	Action     string   `json:"action"`
	Reads      []string `json:"reads"`
	Writes     []string `json:"writes"`
	Children   []string `json:"children"`
	Assessment string   `json:"assessment"`
	Provenance string   `json:"provenance"`
	// Fixed native mechanisms, not an arbitrary command or permission allowlist.
	Boundary string        `json:"boundary"`
	Inputs   []FileBinding `json:"inputs"`
}

type deliveryDisposition struct {
	Kind          string                   `json:"kind"`
	SchemaVersion int                      `json:"schema_version"`
	Input         DeliveryDispositionInput `json:"input"`
	InputSource   FileBinding              `json:"request_input"`
	Context       FileBinding              `json:"context"`
	Source        ReportRetainedSource     `json:"source"`
	Effects       []DeliveryRetainedEffect `json:"retained_effects"`
	Previous      *FileBinding             `json:"previous"`
	RecordedAtUTC string                   `json:"recorded_at_utc"`
}

type DeliveryDispositionStatus struct {
	RequestInput    FileBinding                   `json:"request_input"`
	Record          FileBinding                   `json:"record"`
	Input           DeliveryDispositionInput      `json:"input"`
	Source          ReportRetainedSource          `json:"source"`
	RetainedEffects []DeliveryRetainedEffect      `json:"retained_effects"`
	Operations      []DeliveryDispositionDecision `json:"operations"`
	Cleanup         string                        `json:"cleanup"`
	RuntimeDecision string                        `json:"runtime_decision"`
}

type DeliveryDispositionDecision struct {
	Operation string `json:"operation"`
	State     string `json:"state"`
	Reason    string `json:"reason"`
}

func dispositionBoundary(op string) string {
	switch op {
	case "prepare", "report":
		return "fresh_source_snapshot"
	case "review":
		return "exact_clean_candidate"
	case "verify":
		return "candidate_and_isolated_verifier"
	case "integrate":
		return "qualified_candidate_and_exact_target"
	}
	return ""
}

// Unknown historical dimensions are preserved even when no concrete collision
// is identified. These labels never attest completion or absence.
var dispositionUnknowns = []string{"child_wait_and_descendants", "temp_socket_and_cleanup", "runtime_input_bytes", "filesystem_writes", "pending_writes"}

func dispositionResource(r string) bool {
	if strings.HasPrefix(r, "path:") {
		p := strings.TrimPrefix(r, "path:")
		return filepath.IsAbs(p) && filepath.Clean(p) == p && p != "/" && !strings.ContainsAny(p, "*?[]\x00")
	}
	for _, prefix := range []string{"process:", "network:", "native:"} {
		if strings.HasPrefix(r, prefix) {
			suffix := strings.TrimPrefix(r, prefix)
			return key(suffix) && suffix != "all" && suffix != "any"
		}
	}
	return false
}

func dispositionOverlap(a, b string) bool {
	if strings.HasPrefix(a, "path:") && strings.HasPrefix(b, "path:") {
		a, b = strings.TrimPrefix(a, "path:"), strings.TrimPrefix(b, "path:")
		return sourceContains(a, b) || sourceContains(b, a)
	}
	return a == b
}

func dispositionBindings(inputs []FileBinding) error {
	if len(inputs) == 0 || len(inputs) > 64 {
		return workflowError(4, "scoped disposition requires bounded actual input/evidence bindings")
	}
	for _, b := range inputs {
		if _, e := inputBound(b, inputReportEvidence, 128<<20); e != nil {
			return e
		}
	}
	return nil
}

func dispositionValidate(s workflowState, in DeliveryDispositionInput, checkInputs bool) error {
	if in.Kind != "ply.workflow.scoped-disposition" || in.SchemaVersion != 1 || !key(in.EventID) || len(in.Operations) == 0 || len(in.Operations) > 5 {
		return workflowError(2, "invalid bounded scoped disposition")
	}
	if in.RunID != s.Result.RunID || in.RequestSHA256 != s.Result.RequestSHA256 || in.MandateSHA256 != s.Result.Handoff.SHA256 || in.NativeSession != s.Result.Transport.AgentSessionID {
		return workflowError(4, "disposition Run/request/mandate/owner session differs")
	}
	if len(in.Retained) == 0 || len(in.Retained) > 128 {
		return workflowError(4, "disposition must preserve every unresolved stop")
	}
	for _, r := range in.Retained {
		if !plain(r.Assessment, 1, 2000) || r.Provenance != "owner_and_static_assessment" || !equal(r.Unknowns, dispositionUnknowns) || r.Conflicts == nil {
			return workflowError(4, "retain all unknown dimensions and label owner/static assessment; no runtime absence claim")
		}
		if e := dispositionBindings(r.Evidence); e != nil {
			return e
		}
		for _, c := range r.Conflicts {
			if !dispositionResource(c) {
				return workflowError(4, "invalid concrete retained conflict")
			}
		}
	}
	seen := map[string]bool{}
	for _, op := range in.Operations {
		if dispositionBoundary(op.Operation) == "" || op.Boundary != dispositionBoundary(op.Operation) || seen[op.Operation] || !plain(op.Action, 1, 2000) || !plain(op.Assessment, 1, 2000) || op.Provenance != "owner_and_static_assessment" || len(op.Reads) == 0 || op.Writes == nil || op.Children == nil {
			return workflowError(4, "unsupported, broad or incomplete operation; only fixed native boundaries are supported")
		}
		seen[op.Operation] = true
		if len(op.Inputs) == 0 {
			return workflowError(4, "operation requires bound inputs")
		}
		if checkInputs {
			if e := dispositionBindings(op.Inputs); e != nil {
				return e
			}
		}
		resources := append(append(append([]string{}, op.Reads...), op.Writes...), op.Children...)
		if len(resources) > 128 {
			return workflowError(4, "operation resources exceed bound")
		}
		for _, resource := range resources {
			if !dispositionResource(resource) || !dispositionResourceAllowed(s, op, resource) {
				return workflowError(4, "operation resource is invalid or unbounded")
			}
			if strings.HasPrefix(resource, "path:") {
				p := strings.TrimPrefix(resource, "path:")
				// Refuse lexical/symlink scope tricks. Existing parent chains only; paths
				// for future outputs are checked against their nearest existing ancestor.
				parent := p
				for {
					if _, e := os.Lstat(parent); !os.IsNotExist(e) {
						break
					}
					next := filepath.Dir(parent)
					if next == parent {
						break
					}
					parent = next
				}
				if e := physical(parent, false); e != nil {
					return e
				}
			}
		}
		if op.Operation == "review" {
			for _, write := range op.Writes {
				if write != "native:review-scratch" && write != "native:review-evidence" && write != "native:ignored-build-output" {
					return workflowError(4, "review disposition cannot grant a new source writer")
				}
			}
		}
		for _, effect := range in.Retained {
			raw, e := workflowBound(effect.Stop, 2<<20)
			if e != nil {
				return e
			}
			var stop DeliveryStopRecord
			if e = decode(raw, 2<<20, &stop); e != nil {
				return e
			}
			if deliveryActionText(op.Action) == deliveryActionText(stop.Input.Action) {
				return workflowError(4, "the refused action remains refused")
			}
		}
	}
	return nil
}

func dispositionDir(s workflowState) string {
	return filepath.Join(s.Result.Paths.RunRoot, "delivery", "dispositions")
}

func dispositionHistory(s workflowState) ([]deliveryDisposition, []FileBinding, error) {
	records := []deliveryDisposition{}
	bindings := []FileBinding{}
	for n := 1; n <= 129; n++ {
		path := filepath.Join(dispositionDir(s), fmt.Sprintf("%03d.json", n))
		raw, e := readFile(path, 2<<20, true)
		if os.IsNotExist(e) {
			return records, bindings, nil
		}
		if e != nil {
			return nil, nil, e
		}
		if n > 128 {
			return nil, nil, workflowError(4, "disposition history limit exceeded")
		}
		var r deliveryDisposition
		if e = decode(raw, 2<<20, &r); e != nil {
			return nil, nil, e
		}
		var previous *FileBinding
		if len(bindings) > 0 {
			previous = &bindings[len(bindings)-1]
		}
		if r.Kind != "ply.workflow.scoped-disposition-record" || r.SchemaVersion != 1 || !equal(r.Previous, previous) || r.Context != (FileBinding{s.Result.Paths.Context, s.ContextSHA256}) || r.Input.RunID != s.Result.RunID || r.Input.RequestSHA256 != s.Result.RequestSHA256 || r.Input.MandateSHA256 != s.Result.Handoff.SHA256 || r.Input.NativeSession != s.Result.Transport.AgentSessionID {
			return nil, nil, workflowError(4, "scoped disposition history identity changed")
		}
		requestRaw, e := workflowBound(r.InputSource, 2<<20)
		if e != nil {
			return nil, nil, e
		}
		var request DeliveryDispositionInput
		if e = decode(requestRaw, 2<<20, &request); e != nil {
			return nil, nil, e
		}
		if !equal(request, r.Input) {
			return nil, nil, workflowError(4, "disposition request bytes no longer bind recorded input")
		}
		records = append(records, r)
		bindings = append(bindings, FileBinding{path, hash(raw)})
	}
	return records, bindings, nil
}

// Registration is available from the same owner's actual newer reader before
// a compatible-control transition; the historical request/control stays intact.
func WorkflowDeliveryDisposition(d Dependencies, root, id, contextPath, file, operation string) (WorkflowRun, error) {
	if e := containing(d, root); e != nil {
		return WorkflowRun{}, e
	}
	if (file == "") == (operation == "") {
		return WorkflowRun{}, workflowError(2, "supply either --file or --check-operation")
	}
	var in DeliveryDispositionInput
	var requestRaw []byte
	if file != "" {
		raw, e := readFile(file, 2<<20, true)
		requestRaw = raw
		if e != nil {
			return WorkflowRun{}, e
		}
		if e = decode(raw, 2<<20, &in); e != nil {
			return WorkflowRun{}, e
		}
	}
	e := withStore(root, func() error {
		s, e := workflowRead(root, id)
		if e != nil {
			return e
		}
		if e = deliveryContinuationAuthority(s); e != nil {
			return e
		}
		if e = deliveryIncompleteReportCallback(d, s, contextPath); e != nil {
			return e
		}
		if operation != "" {
			if dispositionBoundary(operation) == "" {
				return workflowError(4, "operation is outside scoped disposition")
			}
			return deliveryDispositionAllows(d, s, operation)
		}
		records, bindings, e := dispositionHistory(s)
		if e != nil {
			return e
		}
		for _, r := range records {
			if r.Input.EventID == in.EventID {
				if !equal(r.Input, in) || r.InputSource.SHA256 != hash(requestRaw) {
					return workflowError(4, "disposition ID already has different immutable bytes")
				}
				return nil
			}
		}
		if s.Result.Delivery.Phase == "closing" {
			return workflowError(4, "delivery already observed; retained effects allow no further source work")
		}
		if !dispositionLocalAgreement(s) || !s.Request.Delivery.Agreement.AutomaticAcceptance() {
			return workflowError(4, "scoped disposition supports only the existing automatic local agreement")
		}
		if e = dispositionValidate(s, in, true); e != nil {
			return e
		}
		if a := s.Result.Delivery.Attempt; a != nil && a.State != "recorded" {
			return workflowError(4, "unresolved native attempt blocks scoped disposition")
		}
		effects, e := deliveryRetainedEffects(s)
		if e != nil {
			return e
		}
		if len(effects) != len(in.Retained) {
			return workflowError(4, "disposition must bind every older and current unresolved stop")
		}
		for i, effect := range effects {
			if effect.Stop != in.Retained[i].Stop {
				return workflowError(4, "disposition unresolved stop set differs")
			}
		}
		source, e := dispositionSource(d, s)
		if e != nil {
			return e
		}
		if source.OID != in.SourceOID || source.Tree != in.SourceTree {
			return workflowError(4, "disposition source candidate drift")
		}
		for _, op := range in.Operations {
			if op.Operation != "prepare" && op.Operation != "report" && !source.Clean {
				return workflowError(4, "review/verify/integrate require an exact clean committed candidate")
			}
			if op.Operation == "verify" {
				matched := false
				for _, input := range op.Inputs {
					if input.Locator == s.Request.Delivery.AcceptancePath {
						matched = true
					}
				}
				if !matched {
					return workflowError(4, "verify disposition must bind the actual frozen acceptance entrypoint bytes")
				}
			}
		}
		inputSource := FileBinding{filepath.Join(dispositionDir(s), "inputs", strings.TrimPrefix(hash(requestRaw), "sha256:")+".json"), hash(requestRaw)}
		if e = d.writeOnce(inputSource.Locator, requestRaw); e != nil {
			return e
		}
		r := deliveryDisposition{Kind: "ply.workflow.scoped-disposition-record", SchemaVersion: 1, Input: in, InputSource: inputSource, Context: FileBinding{contextPath, s.ContextSHA256}, Source: source, Effects: effects, RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
		if len(bindings) > 0 {
			r.Previous = &bindings[len(bindings)-1]
		}
		if len(records) >= 128 {
			return workflowError(4, "disposition history limit reached")
		}
		_, e = workflowKeep(d, filepath.Join(dispositionDir(s), fmt.Sprintf("%03d.json", len(records)+1)), r)
		return e
	})
	return deliveryReadback(d, root, id, e)
}

func deliveryDispositionAllows(d Dependencies, s workflowState, operation string) error {
	if !deliveryObservedStop(s) {
		return nil
	}
	if s.Result.Delivery.Phase == "closing" || s.Result.Delivery.Phase == "completed" {
		return workflowError(4, "delivery already observed; only its receipt is readable, retained effects block further work and cleanup")
	}
	records, _, e := dispositionHistory(s)
	if e != nil {
		return e
	}
	if len(records) == 0 {
		return workflowError(4, "active provider stop blocks "+operation+"; no scoped disposition")
	}
	r := records[len(records)-1]
	if retired, err := deliveryRetirement(s); err != nil {
		return err
	} else if retired != nil {
		return workflowError(4, "owner_retired: scoped continuation cannot revive ownership")
	}
	if e = deliveryAcceptedAuthority(s); e != nil {
		return e
	}
	if e = dispositionValidate(s, r.Input, false); e != nil {
		return e
	}
	effects, e := deliveryRetainedEffects(s)
	if e != nil {
		return e
	}
	if !equal(effects, r.Effects) {
		return workflowError(4, "additional or changed unresolved stop blocks "+operation)
	}
	var selected *DeliveryDispositionOperation
	for i := range r.Input.Operations {
		if r.Input.Operations[i].Operation == operation {
			selected = &r.Input.Operations[i]
		}
	}
	if selected == nil {
		return workflowError(4, "retained effects block "+operation+"; operation outside scoped disposition")
	}
	if e := dispositionBindings(selected.Inputs); e != nil {
		return e
	}
	if e := dispositionConflicts(s, r.Input, *selected); e != nil {
		return e
	}
	// A reserved verification/integration may finish only its exact candidate.
	// Existing native recovery checks determine whether it can resume; no new
	// operation may use a pending attempt as permission to start another effect.
	if a := s.Result.Delivery.Attempt; a != nil && a.State != "recorded" {
		if !((operation == "verify" && a.Kind == "verification" || operation == "integrate" && a.Kind == "integration") && a.CandidateOID == r.Source.OID && a.CandidateTree == r.Source.Tree) {
			return workflowError(4, "unresolved native attempt blocks "+operation)
		}
	}
	source, e := dispositionSource(d, s)
	if e != nil {
		return e
	}
	if !equal(source, r.Source) {
		return workflowError(4, "source snapshot drift blocks "+operation+"; retain old disposition and assess new source")
	}
	if operation != "prepare" && operation != "report" && !source.Clean {
		return workflowError(4, "scoped operation requires exact clean candidate")
	}
	if s.Result.Delivery.OwnershipRelease != nil {
		return workflowError(4, "scoped continuation cannot transfer writer ownership")
	}
	return nil
}

func deliveryDispositionStatus(d Dependencies, s workflowState) (*DeliveryDispositionStatus, error) {
	records, bindings, e := dispositionHistory(s)
	if e != nil || len(records) == 0 {
		return nil, e
	}
	r := records[len(records)-1]
	out := &DeliveryDispositionStatus{RequestInput: r.InputSource, Record: bindings[len(bindings)-1], Input: r.Input, Source: r.Source, RetainedEffects: r.Effects, Cleanup: "blocked: retained unknown child/writer effects; source retained", RuntimeDecision: "per_action; disposition is not runtime approval or absence evidence", Operations: []DeliveryDispositionDecision{}}
	for _, op := range []string{"prepare", "report", "review", "verify", "integrate"} {
		decision := DeliveryDispositionDecision{Operation: op, State: "eligible_with_retained_effects", Reason: "Exact native operation boundary only; candidate/review/acceptance/target gates still apply. Unknown effects remain unknown."}
		if e := deliveryDispositionAllows(d, s, op); e != nil {
			decision.State, decision.Reason = "blocked", e.Error()
		}
		out.Operations = append(out.Operations, decision)
	}
	return out, nil
}

// Avoid blanket callback relaxation. Only these call sites select a supported
// operation; all other callbacks keep the original stop guard.
func deliveryScopedCallback(d Dependencies, s workflowState, contextPath, operation string) error {
	if !deliveryObservedStop(s) {
		return deliveryCallback(d, s, contextPath, true)
	}
	if e := workflowCallback(d, s, contextPath); e != nil {
		return e
	}
	if !deliveryStarted(s) || s.Result.Delivery.Phase == "completed" || s.Result.Delivery.OwnershipRelease != nil {
		return workflowError(4, "scoped callback requires the same active accepted owner")
	}
	if e := deliveryDispositionAllows(d, s, operation); e != nil {
		return e
	}
	return deliveryLiveOwner(d, s)
}

func dispositionLocalAgreement(s workflowState) bool {
	a := s.Request.Delivery.Agreement
	return a != nil && (a.Mode == workspace.DeliveryLocalBranch || a.Mode == workspace.DeliveryLocalEpic)
}

func dispositionVerifierInput(s workflowState, binary string) error {
	records, _, e := dispositionHistory(s)
	if e != nil {
		return e
	}
	if len(records) == 0 {
		return workflowError(4, "missing scoped verifier inputs")
	}
	for _, op := range records[len(records)-1].Input.Operations {
		if op.Operation == "verify" {
			for _, input := range op.Inputs {
				if input.Locator == binary {
					return nil
				}
			}
		}
	}
	return workflowError(4, "scoped verify must bind the actual candidate binary bytes before execution")
}

func dispositionSource(d Dependencies, s workflowState) (ReportRetainedSource, error) {
	x, e := d.Workspace.IntegrationGit.ObserveIntegrationWorktree(s.Observed.Target.WorktreeLocator, s.Observed.Target.Ref)
	if e != nil {
		return ReportRetainedSource{}, e
	}
	if x.Locator != s.Observed.Target.WorktreeLocator || x.Ref != s.Observed.Target.Ref || x.GitCommonDir != s.Observed.Target.GitCommonDir || len(x.InProgress) > 0 {
		return ReportRetainedSource{}, workflowError(4, "scoped source identity or Git operation changed")
	}
	sha, e := workflowTargetFiles(d, s, true)
	if e != nil {
		return ReportRetainedSource{}, e
	}
	return ReportRetainedSource{Worktree: x.Locator, Ref: x.Ref, OID: x.OID, Tree: x.Tree, Clean: x.Clean, SourceSHA256: sha, GitCommonDir: x.GitCommonDir}, nil
}

func dispositionResourceAllowed(s workflowState, op DeliveryDispositionOperation, resource string) bool {
	if strings.HasPrefix(resource, "path:") {
		p := strings.TrimPrefix(resource, "path:")
		if sourceContains(s.Observed.Target.WorktreeLocator, p) || sourceContains(s.Result.Paths.RunRoot, p) || p == s.Request.Delivery.AcceptancePath {
			return true
		}
		if op.Operation == "integrate" && p == s.Request.Delivery.Agreement.TargetWorktree {
			return true
		}
		// External immutable operation inputs may be read, never written by this
		// disposition. The native binary/scratch snapshots handle verifier copies.
		for _, input := range op.Inputs {
			if input.Locator == p {
				for _, write := range op.Writes {
					if dispositionOverlap(write, resource) {
						return false
					}
				}
				return true
			}
		}
		return false
	}
	if strings.HasPrefix(resource, "native:") {
		allowed := map[string][]string{
			"prepare":   {"native:git-index", "native:run-evidence", "native:ignored-build-output", "native:verification-scratch"},
			"report":    {"native:run-evidence"},
			"review":    {"native:review-scratch", "native:review-evidence", "native:ignored-build-output"},
			"verify":    {"native:verification-scratch", "native:run-evidence", "native:ignored-build-output"},
			"integrate": {"native:local-target", "native:run-evidence"},
		}
		for _, r := range allowed[op.Operation] {
			if r == resource {
				return true
			}
		}
		return false
	}
	// Named actual children/network dependencies are assessments, never OS grants.
	return true
}

func dispositionConflicts(s workflowState, in DeliveryDispositionInput, op DeliveryDispositionOperation) error {
	resources := append(append(append([]string{}, op.Reads...), op.Writes...), op.Children...)
	// Native facts provide implicit dependencies even if an assessment omits them.
	resources = append(resources, "path:"+s.Observed.Target.WorktreeLocator, "path:"+s.Observed.Target.GitCommonDir, "path:"+s.Result.Paths.RunRoot, "native:run-evidence")
	for _, input := range op.Inputs {
		resources = append(resources, "path:"+input.Locator)
	}
	switch op.Operation {
	case "prepare":
		resources = append(resources, "path:"+s.Request.Delivery.AcceptancePath, "native:git-index")
	case "review":
		resources = append(resources, "native:review-scratch", "native:review-evidence")
	case "verify":
		resources = append(resources, "path:"+s.Request.Delivery.AcceptancePath, "native:verification-scratch")
	case "integrate":
		resources = append(resources, "path:"+s.Request.Delivery.Agreement.TargetWorktree, "native:local-target")
	}
	for _, retained := range in.Retained {
		for _, conflict := range retained.Conflicts {
			for _, resource := range resources {
				if dispositionOverlap(resource, conflict) {
					return workflowError(4, fmt.Sprintf("%s depends on retained possible effect %s from %s", op.Operation, conflict, retained.Stop.SHA256))
				}
			}
		}
	}
	return nil
}
