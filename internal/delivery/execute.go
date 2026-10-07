package delivery

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func (s *Service) stop(root string, r *Receipt, state, reason, next string) error {
	r.State = state
	r.Reasons = []string{reason}
	r.NextAction = next
	s.event(r, "delivery."+state, r.AttemptID+"/"+state+"/"+shortHash(reason), reason, "")
	return s.save(root, r)
}

func (s *Service) Execute(cwd, id string) (Receipt, error) {
	d, root, e := s.context(cwd)
	if e != nil {
		return Receipt{}, e
	}
	var out Receipt
	e = withLock(root, func() error {
		r, e := readReceipt(root, id)
		out = r
		if e != nil {
			return e
		}
		defer func() { out = r }()
		if r.State == "delivered" {
			return nil
		}
		r.Attempts++
		r.AttemptID = fmt.Sprintf("%s/attempt/%06d", r.ID, r.Attempts)
		s.event(&r, "delivery.attempt_started", r.AttemptID, "Explicit delivery attempt started; no background worker.", "")
		if e = s.save(root, &r); e != nil {
			return e
		}
		if r.Manifest.Agreement.Mode == workspace.DeliveryPullRequest && r.Metadata == nil {
			choice, e := resolvePreferences(s.options.PreferencesPath, r.Manifest.Agreement.GitHubRepository, r.Registration.Metadata, s.now())
			if e != nil {
				return s.stop(root, &r, "blocked", e.Error(), "Correct the local PR preferences, then explicitly execute the same Delivery ID.")
			}
			r.Metadata = &choice
			s.event(&r, "delivery.metadata_frozen", "metadata/1", stringMustCanonical(choice), "")
			if e = s.save(root, &r); e != nil {
				return e
			}
		}
		auth, e := s.gate(d, root, &r)
		if e != nil {
			return s.stop(root, &r, "blocked", e.Error(), "Inspect the precise evidence or authority conflict; changed candidate, mode or target requires a new TaskResult and authorized Delivery.")
		}
		s.event(&r, "delivery.gate_observed", r.AttemptID+"/gate", r.State+": "+strings.Join(r.Reasons, "; "), "")
		if r.State != "ready" {
			return s.save(root, &r)
		}
		fx, e := readEffect(root, r.Manifest.EffectKey)
		if os.IsNotExist(e) {
			fx = effect{Kind: "ply.delivery.effect", SchemaVersion: 1, Key: r.Manifest.EffectKey, OwnerDeliveryID: r.ID}
			if e = saveEffect(root, fx); e != nil {
				return e
			}
		} else if e != nil {
			return e
		}
		if fx.OwnerDeliveryID != r.ID {
			owner, e := readReceipt(root, fx.OwnerDeliveryID)
			if e != nil {
				return e
			}
			if owner.Manifest.EffectKey != r.Manifest.EffectKey {
				return fmt.Errorf("candidate effect reservation conflicts")
			}
			r.ReusedDeliveryID = owner.ID
			if fx.Complete {
				return s.reuseCompleted(root, &r, &fx, owner)
			}
			r.PR = fx.PR
			r.Local = fx.Local
			return s.stop(root, &r, "blocked", "This candidate/target has an unfinished effect reserved by "+owner.ID, "Execute Delivery "+owner.ID+" to observe and resume that preserved attempt; no second effect is started.")
		}
		r.State = "delivering"
		r.Reasons = []string{}
		r.NextAction = "This invocation owns the reserved effect; if interrupted, explicitly execute the same Delivery ID to inspect and resume."
		if e = s.save(root, &r); e != nil {
			return e
		}
		if r.Manifest.Agreement.Mode == workspace.DeliveryPullRequest {
			return s.executePR(d, root, &r, &fx)
		}
		return s.executeLocal(d, root, &r, &fx, auth)
	})
	return out, e
}

func (s *Service) reuseCompleted(root string, r *Receipt, fx *effect, owner Receipt) error {
	if owner.Manifest.TaskID != r.Manifest.TaskID || owner.Manifest.WorkflowRunID != r.Manifest.WorkflowRunID || owner.Manifest.PreparationID != r.Manifest.PreparationID || owner.Manifest.TaskResult.ID != r.Manifest.TaskResult.ID || owner.Manifest.TaskResultSHA256 != r.Manifest.TaskResultSHA256 || !reflect.DeepEqual(owner.Manifest.Agreement, r.Manifest.Agreement) {
		return s.stop(root, r, "blocked", "The existing candidate/target receipt belongs to a different native Task, run, candidate evidence or source agreement.", "Inspect Delivery "+owner.ID+"; its observed effect cannot be claimed as this run's native closure.")
	}
	if owner.State != "delivered" || !owner.NativeClosed || owner.PR != nil && (owner.Metadata == nil || owner.Metadata.State != "applied" || !owner.PR.MetadataApplied) || !reflect.DeepEqual(fx.PR, owner.PR) || !reflect.DeepEqual(fx.Local, owner.Local) {
		return s.stop(root, r, "blocked", "The effect owner's completion or metadata adjustment is not fully observed.", "Execute Delivery "+owner.ID+" to finish its preserved outcome before reusing it.")
	}
	if metadataRevision(r) > 1 || r.Registration.Metadata != nil && !sameMetadataSelection(r.Registration.Metadata, owner.Metadata) {
		return s.stop(root, r, "blocked", "This candidate/target already has a completed receipt with different metadata.", "Use delivery metadata on "+owner.ID+" to record an explicit metadata adjustment; the completed effect will be reused.")
	}
	if fx.PR != nil {
		copy := *fx.PR
		copy.Disposition = "reused"
		r.PR = &copy
	}
	r.Local = fx.Local
	r.Metadata = owner.Metadata
	r.NativeClosed = true
	s.event(r, "delivery.effect_reused", "effect_reused", "Reused the observed completed candidate/target receipt from "+owner.ID, urlOf(r.PR))
	return s.complete(root, r, fx)
}

func sameMetadataSelection(selection *MetadataSelection, choice *MetadataChoice) bool {
	if selection == nil {
		return true
	}
	if choice == nil {
		return false
	}
	if selection.Assignees != nil && !reflect.DeepEqual(*selection.Assignees, choice.Assignees) {
		return false
	}
	if selection.Reviewers != nil && !reflect.DeepEqual(*selection.Reviewers, choice.Reviewers) {
		return false
	}
	return reflect.DeepEqual(normalizeEmpty(selection.RemoveAssignees), normalizeEmpty(choice.RemoveAssignees)) && reflect.DeepEqual(normalizeEmpty(selection.RemoveReviewers), normalizeEmpty(choice.RemoveReviewers))
}
func normalizeEmpty(v []string) []string {
	if len(v) == 0 {
		return []string{}
	}
	return v
}
func stringMustCanonical(v any) string { raw, _ := canonical(v); return string(raw) }
func urlOf(p *PRReceipt) string {
	if p == nil {
		return ""
	}
	return p.URL
}

func (s *Service) complete(root string, r *Receipt, fx *effect) error {
	r.State = "delivered"
	r.Reasons = []string{}
	r.NextAction = "Delivery is complete. Its receipt preserves the observed outcome; future publishing or merging is a separate action."
	s.event(r, "delivery.completed", "completed/"+fmt.Sprint(metadataRevision(r)), "Selected delivery effect and native Task closure observed.", urlOf(r.PR))
	if e := s.save(root, r); e != nil {
		return e
	}
	if fx.Complete && fx.OwnerDeliveryID != r.ID {
		return nil
	}
	fx.Complete = true
	fx.PR = r.PR
	fx.Local = r.Local
	return saveEffect(root, *fx)
}
func metadataRevision(r *Receipt) int {
	if r.Metadata == nil {
		return 0
	}
	return r.Metadata.Revision
}

func (s *Service) executeLocal(d taskrun.Dependencies, root string, r *Receipt, fx *effect, auth taskrun.DeliveryAuthority) error {
	if !fx.LocalStarted {
		preview, e := workspace.CheckTaskIntegration(d.Workspace, integrationInput(*r, auth))
		if e != nil {
			return s.stop(root, r, "blocked", e.Error(), "Resolve the native integration precondition without changing the accepted candidate.")
		}
		if preview.Readback.Classification != "ready" && preview.Readback.RecoveryStatus != "complete" {
			return s.stop(root, r, "blocked", "native integration: "+preview.Readback.Classification, "Inspect the target and recorded integration plan; a moved base requires new control.")
		}
		fx.LocalStarted = true
		if e = saveEffect(root, *fx); e != nil {
			return e
		}
		s.event(r, "delivery.local_intent", "local_intent", "Reserved native integration of the exact candidate to "+r.Manifest.Agreement.TargetRef, "")
		if e = s.save(root, r); e != nil {
			return e
		}
	}
	// Native Check/Apply reconciles its own intent, attempt and observed Git effect.
	result, e := taskrun.CompleteLocalDelivery(d, root, r.Manifest.WorkflowRunID, r.Manifest.TaskResult.ID, r.HumanQA.ID)
	if e != nil {
		return s.stop(root, r, "unknown_effect", e.Error(), "Inspect native integration, queue and Epic base, then explicitly execute this Delivery ID to resume the preserved attempt.")
	}
	if e = s.fault("after_local_effect"); e != nil {
		return e
	}
	if !result.Completed || result.Integration.Readback.RecoveryStatus != "complete" || result.Integration.Readback.ParentOID != r.Manifest.TaskResult.ResultOID {
		return s.stop(root, r, "unknown_effect", "native local completion is not fully observed", "Inspect the native integration and base records, then resume this Delivery ID.")
	}
	r.Local = &LocalReceipt{TargetRef: r.Manifest.Agreement.TargetRef, TargetWorktree: r.Manifest.Agreement.TargetWorktree, BeforeOID: r.Manifest.ExpectedParentOID, AfterOID: result.Integration.Readback.ParentOID, Integration: result.Integration.Readback, Base: result.Base, Queue: result.Queue, ObservedAtUTC: s.now()}
	r.NativeClosed = true
	fx.Local = r.Local
	if e = saveEffect(root, *fx); e != nil {
		return e
	}
	s.event(r, "delivery.local_observed", "local_observed", "Exact local integration, native base update and this Task's queue closure observed.", "")
	return s.complete(root, r, fx)
}

func (s *Service) executePR(d taskrun.Dependencies, root string, r *Receipt, fx *effect) error {
	t := prTarget(*r)
	o, e := s.options.PR.Observe(t)
	if e != nil {
		state := "blocked"
		if fx.PushStarted || fx.CreateStarted || fx.PR != nil {
			state = "unknown_effect"
		}
		return s.stop(root, r, state, e.Error(), "Restore remote/PR observation and explicitly execute this Delivery ID; no publication is blindly repeated.")
	}
	p, e := validateObservation(t, o)
	if e != nil {
		return s.stop(root, r, "blocked", e.Error(), "Resolve the conflicting PR/head/base selection; a new target requires a new authorized agreement.")
	}
	if fx.PR != nil {
		if p == nil || p.Number != fx.PR.Number || p.URL != fx.PR.URL {
			return s.stop(root, r, "unknown_effect", "Previously observed PR is no longer an exact open match.", "Inspect the preserved PR receipt and its current state. Delivery will not create another PR.")
		}
	}
	if fx.CreateStarted && p == nil {
		return s.stop(root, r, "unknown_effect", "PR creation was reserved but its outcome cannot be established.", "Inspect the GitHub PR creation outcome; do not repeat the create operation or use another Delivery ID.")
	}
	if o.RemoteOID != t.OID {
		if p != nil {
			return s.stop(root, r, "blocked", "PR and published branch head disagree.", "Inspect the remote candidate; delivery will not overwrite a conflicting branch.")
		}
		if fx.PushStarted {
			return s.stop(root, r, "unknown_effect", "A reserved push is not observed at the exact candidate.", "Inspect the published source branch; no repeated or force push is attempted.")
		}
		if o.RemoteOID != "" {
			ff, e := s.options.PR.CanFastForward(t, o.RemoteOID)
			if e != nil {
				return s.stop(root, r, "blocked", e.Error(), "Resolve source ancestry observation before publication.")
			}
			if !ff {
				return s.stop(root, r, "blocked", "Source publication would require force-push.", "Select a nonconflicting source or prepare a newly authorized candidate; delivery never force-pushes.")
			}
		}
		fx.PushStarted = true
		if e = saveEffect(root, *fx); e != nil {
			return e
		}
		s.event(r, "delivery.push_intent", "push_intent", "Reserved publication of only "+t.SourceRef, "")
		if e = s.save(root, r); e != nil {
			return e
		}
		pushErr := s.options.PR.Push(t)
		if e = s.fault("after_push_effect"); e != nil {
			return e
		}
		o, e = s.options.PR.Observe(t)
		if e != nil || o.RemoteOID != t.OID {
			return s.stop(root, r, "unknown_effect", joinError("Push outcome is not observed at the exact candidate", pushErr, e), "Restore source-ref observation and execute this Delivery ID to reconcile the reserved push.")
		}
		p, e = validateObservation(t, o)
		if e != nil {
			return s.stop(root, r, "blocked", e.Error(), "Resolve the observed PR conflict; the published source branch is preserved.")
		}
	}
	if !fx.PushObserved {
		fx.PushObserved = true
		if e = saveEffect(root, *fx); e != nil {
			return e
		}
		s.event(r, "delivery.push_observed", "push_observed", "Exact selected source branch commit observed; base branch unchanged.", "")
		if e = s.save(root, r); e != nil {
			return e
		}
	}
	if p == nil {
		fx.CreateStarted = true
		if e = saveEffect(root, *fx); e != nil {
			return e
		}
		s.event(r, "delivery.pr_create_intent", "pr_create_intent", "Reserved one PR creation for the exact repository, source and base.", "")
		if e = s.save(root, r); e != nil {
			return e
		}
		title := r.Registration.Title
		if title == "" {
			title = string(r.Manifest.TaskID)
		}
		createErr := s.options.PR.Create(t, title, r.Registration.Body)
		if e = s.fault("after_pr_effect"); e != nil {
			return e
		}
		o, e = s.options.PR.Observe(t)
		if e != nil {
			return s.stop(root, r, "unknown_effect", joinError("PR create response requires observation", createErr, e), "Restore PR observation and resume this Delivery ID; the create operation will not be repeated.")
		}
		p, e = validateObservation(t, o)
		if e != nil || p == nil {
			return s.stop(root, r, "unknown_effect", joinError("Reserved PR creation is not observed at the exact candidate", createErr, e), "Inspect the PR creation result; delivery will not create a second PR.")
		}
	}
	if o.RemoteOID != t.OID {
		return s.stop(root, r, "unknown_effect", "Published branch changed during PR creation observation.", "Inspect the candidate and remote before retry.")
	}
	if e = s.recordPR(root, r, fx, *p); e != nil {
		return e
	}
	if r.Metadata == nil {
		return fmt.Errorf("frozen PR metadata is missing")
	}
	if !metadataSatisfied(*p, *r.Metadata) {
		r.Metadata.State = "applying"
		r.Metadata.Error = ""
		s.event(r, "delivery.metadata_intent", r.AttemptID+"/metadata", stringMustCanonical(*r.Metadata), p.URL)
		if e = s.save(root, r); e != nil {
			return e
		}
		applyErr := s.options.PR.ApplyMetadata(t, *p, *r.Metadata)
		if e = s.fault("after_metadata_effect"); e != nil {
			return e
		}
		o, e = s.options.PR.Observe(t)
		if e != nil {
			r.Metadata.State = "unknown"
			r.Metadata.Error = joinError("Metadata result could not be observed", applyErr, e)
			return s.stop(root, r, "unknown_effect", r.Metadata.Error, "Restore GitHub observation, then resume; the existing PR is preserved.")
		}
		p, e = validateObservation(t, o)
		if e != nil || p == nil || p.Number != r.PR.Number || o.RemoteOID != t.OID {
			return s.stop(root, r, "unknown_effect", joinError("PR identity changed during metadata application", e), "Inspect the preserved PR; delivery will not create a replacement.")
		}
		if applyErr != nil || !metadataSatisfied(*p, *r.Metadata) {
			r.Metadata.State = "failed"
			r.Metadata.Error = joinError("PR exists; metadata is not fully applied", applyErr)
			return s.stop(root, r, "failed", r.Metadata.Error, "Correct GitHub permissions or accounts; retry reuses this PR and frozen choices. Use delivery metadata for an explicit recorded adjustment.")
		}
	}
	r.Metadata.State = "applied"
	r.Metadata.Error = ""
	r.PR.MetadataApplied = true
	s.event(r, "delivery.metadata_applied", fmt.Sprintf("metadata/%d/applied", r.Metadata.Revision), stringMustCanonical(*r.Metadata), r.PR.URL)
	if e = s.save(root, r); e != nil {
		return e
	}
	// Preserve a content-addressed observation before native workflow closure.
	observed := struct {
		Kind           string          `json:"kind"`
		SchemaVersion  int             `json:"schema_version"`
		DeliveryID     string          `json:"delivery_id"`
		ManifestSHA256 string          `json:"manifest_sha256"`
		PR             *PRReceipt      `json:"pull_request"`
		Metadata       *MetadataChoice `json:"metadata"`
	}{"ply.delivery.pr-observation", 1, r.ID, r.ManifestSHA256, r.PR, r.Metadata}
	raw, e := canonical(observed)
	if e != nil {
		return e
	}
	path := filepath.Join(storeRoot(root), "observations", strings.TrimPrefix(hash(raw), "sha256:")+".json")
	if e = atomicJSON(path, observed, true); e != nil {
		return e
	}
	_, e = taskrun.CompletePullRequestDelivery(d, root, r.Manifest.WorkflowRunID, r.Manifest.TaskResult.ID, taskrun.FileBinding{Locator: path, SHA256: hash(raw)}, taskrun.PullRequestDeliveryObservation{DeliveryID: r.ID, Repository: t.Repository, HeadRef: t.SourceRef, HeadOID: t.OID, BaseRef: t.BaseRef, URL: r.PR.URL, Number: r.PR.Number, MetadataApplied: true})
	if e != nil {
		return s.stop(root, r, "failed", "PR and metadata observed; native task closure remains incomplete: "+e.Error(), "Resume this Delivery ID to finish native closure; the PR will be reused.")
	}
	if e = s.fault("after_pr_native_close"); e != nil {
		return e
	}
	r.NativeClosed = true
	return s.complete(root, r, fx)
}

func joinError(prefix string, errs ...error) string {
	for _, e := range errs {
		if e != nil {
			prefix += "; " + e.Error()
		}
	}
	return prefix
}

func (s *Service) recordPR(root string, r *Receipt, fx *effect, p PullRequest) error {
	if fx.PR == nil {
		disposition := "reused"
		if fx.CreateStarted && strings.Contains(p.Body, marker(prTarget(*r))) {
			disposition = "created"
		}
		observed := &PRReceipt{Repository: p.Repository, Number: p.Number, URL: p.URL, HeadRef: p.HeadRef, HeadOID: p.HeadOID, BaseRef: p.BaseRef, Disposition: disposition, ObservedAtUTC: s.now()}
		owner := r
		if fx.OwnerDeliveryID != r.ID {
			v, e := readReceipt(root, fx.OwnerDeliveryID)
			if e != nil {
				return e
			}
			owner = &v
		}
		if disposition == "created" {
			s.event(owner, "delivery.pr_created", "pr_created", "PR creation observed for the reserved candidate/target.", p.URL)
			for _, event := range owner.Events {
				if event.Kind == "delivery.pr_created" {
					observed.CreationEventID = event.ID
				}
			}
		} else {
			s.event(owner, "delivery.pr_reused", "pr_reused", "Existing exact PR observed and reused; no new creation claimed.", p.URL)
		}
		owner.PR = observed
		if e := s.save(root, owner); e != nil {
			return e
		}
		fx.PR = observed
		if e := saveEffect(root, *fx); e != nil {
			return e
		}
	}
	copy := *fx.PR
	r.PR = &copy
	if fx.OwnerDeliveryID != r.ID {
		s.event(r, "delivery.pr_reused", "pr_reused", "Reused PR from candidate/target effect owned by "+fx.OwnerDeliveryID, p.URL)
	}
	return s.save(root, r)
}

func (s *Service) Metadata(cwd, id, file string) (Receipt, error) {
	var selection MetadataSelection
	if e := readJSON(file, &selection); e != nil {
		return Receipt{}, e
	}
	if e := validateMetadata(&selection); e != nil {
		return Receipt{}, e
	}
	_, root, e := s.context(cwd)
	if e != nil {
		return Receipt{}, e
	}
	var out Receipt
	e = withLock(root, func() error {
		r, e := readReceipt(root, id)
		if e != nil {
			return e
		}
		if r.Manifest.Agreement.Mode != workspace.DeliveryPullRequest {
			return fmt.Errorf("metadata adjustment is only available for pull_request delivery")
		}
		fx, effectErr := readEffect(root, r.Manifest.EffectKey)
		if effectErr == nil && fx.OwnerDeliveryID != r.ID {
			return fmt.Errorf("candidate/target effect belongs to Delivery %s; record the explicit metadata adjustment on that Delivery ID", fx.OwnerDeliveryID)
		}
		if effectErr != nil && !os.IsNotExist(effectErr) {
			return effectErr
		}
		if r.Metadata == nil {
			choice, e := resolvePreferences(s.options.PreferencesPath, r.Manifest.Agreement.GitHubRepository, r.Registration.Metadata, s.now())
			if e != nil {
				return e
			}
			r.Metadata = &choice
		}
		choice, e := reviseChoice(*r.Metadata, selection, s.now())
		if e != nil {
			return e
		}
		r.Metadata = &choice
		if r.PR != nil {
			copy := *r.PR
			copy.MetadataApplied = false
			r.PR = &copy
		}
		r.State = "registered"
		r.Reasons = []string{}
		r.NextAction = "Explicit metadata adjustment is preserved. Execute this Delivery ID to apply it to the existing PR."
		s.event(&r, "delivery.metadata_adjusted", fmt.Sprintf("metadata/%d", choice.Revision), stringMustCanonical(choice), urlOf(r.PR))
		e = s.save(root, &r)
		out = r
		return e
	})
	return out, e
}
