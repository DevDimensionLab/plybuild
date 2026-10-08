package integration

import (
	"errors"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

type rejectedMergeFixture struct {
	observes     int
	merges       int
	observation  PRMergeObservation
	failure      error
	operationIDs []string
}

func (a *rejectedMergeFixture) Observe(t PRMergeTarget) (PRMergeObservation, error) {
	a.observes++
	a.operationIDs = append(a.operationIDs, t.OperationID)
	o := a.observation
	o.Repository, o.Number, o.HeadRef, o.HeadOID, o.BaseRef = t.Repository, t.Number, t.HeadRef, t.HeadOID, t.BaseRef
	return o, a.failure
}
func (a *rejectedMergeFixture) Merge(t PRMergeTarget) (PRMergeObservation, error) {
	a.merges++
	o, _ := a.Observe(t)
	o.State, o.Reason, o.NoEffect = "blocked", "merge_blocked", true
	return o, nil
}

func rejectedMergeReceipt(t *testing.T) Receipt {
	r := serviceEffectReceipt(t, adapterPhysicalTemp(t))
	r.Integrated = false
	r.IntegratedOID = ""
	r.Plan.Flow = "github_pr_merge"
	target := mergeTargetFixture()
	r.Plan.PR = &target
	r.IntegrationStarted = true
	r.State = "blocked"
	r.PR = &PRMergeObservation{State: "blocked", Reason: "merge_blocked", NoEffect: true, Repository: target.Repository, Number: target.Number, HeadRef: target.HeadRef, HeadOID: target.HeadOID, BaseRef: target.BaseRef}
	return r
}

func TestServiceKeepsProvenRejectedMergeRetryableWithoutAutomaticallyRetrying(t *testing.T) {
	a := &rejectedMergeFixture{observation: PRMergeObservation{State: "ready"}}
	s := NewService(taskrun.Dependencies{}, Options{PR: a})
	r := rejectedMergeReceipt(t)
	if err := s.merge(&r); err != nil {
		t.Fatal(err)
	}
	if r.State != "merge_blocked" || r.PR == nil || !r.PR.NoEffect || a.merges != 0 {
		t.Fatalf("known rejection became unknown or retried implicitly: state=%s observation=%+v merges=%d", r.State, r.PR, a.merges)
	}
}

func TestServicePreservesKnownRejectionAcrossObservationFailure(t *testing.T) {
	a := &rejectedMergeFixture{observation: PRMergeObservation{State: "effect_unknown"}, failure: errors.New("temporary readback failure")}
	s := NewService(taskrun.Dependencies{}, Options{PR: a})
	r := rejectedMergeReceipt(t)
	r.PR.OperationID = "known-rejected-operation"
	if err := s.merge(&r); err != nil {
		t.Fatal(err)
	}
	a.failure, a.observation.State = nil, "ready"
	if err := s.merge(&r); err != nil {
		t.Fatal(err)
	}
	if r.State != "merge_blocked" || r.PR == nil || !r.PR.NoEffect || r.PR.OperationID != "known-rejected-operation" || a.merges != 0 {
		t.Fatalf("temporary observation failure erased the positively rejected request: state=%s request=%+v merges=%d", r.State, r.PR, a.merges)
	}
	if err := s.prepareMergeRetry(&r); err != nil {
		t.Fatal(err)
	}
	if len(r.MergeAttempts) != 1 || r.MergeAttempts[0].OperationID != "known-rejected-operation" || !r.MergeAttempts[0].NoEffect {
		t.Fatal("explicit retry lost the original proven rejected request")
	}
	if err := s.merge(&r); err != nil || a.merges != 1 {
		t.Fatalf("controlled retry did not dispatch exactly once: %v merges=%d", err, a.merges)
	}
}

func TestServicePreservesPendingOperationIDAcrossIncompleteObservation(t *testing.T) {
	a := &rejectedMergeFixture{failure: errors.New("lost read response")}
	s := NewService(taskrun.Dependencies{}, Options{PR: a})
	r := rejectedMergeReceipt(t)
	r.PR.NoEffect, r.PR.State, r.PR.OperationID = false, "remote_pending", "pending-operation"
	if err := s.merge(&r); err != nil {
		t.Fatal(err)
	}
	a.failure, a.observation.State = nil, "remote_pending"
	if err := s.merge(&r); err != nil {
		t.Fatal(err)
	}
	if r.PR.OperationID != "pending-operation" || a.operationIDs[1] != "pending-operation" || r.PR.NoEffect || a.merges != 0 {
		t.Fatalf("lost pending request identity across read failure: observation=%+v calls=%v", r.PR, a.operationIDs)
	}
}

func TestExplicitMergeRetryArchivesRejectedRequestAndDispatchesOnce(t *testing.T) {
	a := &rejectedMergeFixture{observation: PRMergeObservation{State: "ready"}}
	s := NewService(taskrun.Dependencies{}, Options{PR: a})
	r := rejectedMergeReceipt(t)
	r.PR.OperationID = "previous-rejected-operation"
	if err := s.prepareMergeRetry(&r); err != nil {
		t.Fatal(err)
	}
	if r.IntegrationStarted || len(r.MergeAttempts) != 1 || r.MergeAttempts[0].OperationID != "previous-rejected-operation" || !r.MergeAttempts[0].NoEffect || r.MergeRetryAtUTC == "" || r.PR.OperationID != "" || a.merges != 0 {
		t.Fatalf("retry lost original request or dispatched before its authority was durable: %+v", r)
	}
	if len(a.operationIDs) != 2 || a.operationIDs[0] != "previous-rejected-operation" || a.operationIDs[1] != "" {
		t.Fatalf("retry did not separately observe old request and preflight exact PR: %v", a.operationIDs)
	}
	if err := s.merge(&r); err != nil {
		t.Fatal(err)
	}
	if err := s.merge(&r); err != nil {
		t.Fatal(err)
	}
	if a.merges != 1 || r.State != "merge_blocked" || len(r.MergeAttempts) != 1 {
		t.Fatalf("one explicit retry repeated or lost history: merges=%d state=%s prior=%d", a.merges, r.State, len(r.MergeAttempts))
	}
}

func TestExplicitMergeRetryRequiresExactPositiveNoEffectProof(t *testing.T) {
	for _, name := range []string{"unknown-request", "pending", "unknown", "merged", "wrong-rejection-identity", "malformed-rejection-state", "observation-error", "cancelled", "installed", "integrated"} {
		t.Run(name, func(t *testing.T) {
			r := rejectedMergeReceipt(t)
			a := &rejectedMergeFixture{observation: PRMergeObservation{State: "ready"}}
			switch name {
			case "unknown-request":
				r.PR.NoEffect = false
			case "pending":
				a.observation.State = "remote_pending"
			case "unknown":
				a.observation.State = "effect_unknown"
			case "merged":
				a.observation.State = "merged"
			case "wrong-rejection-identity":
				r.PR.HeadOID = "different-candidate"
			case "malformed-rejection-state":
				r.PR.State = "remote_pending"
			case "observation-error":
				a.failure = errors.New("readback unavailable")
			case "cancelled":
				r.Reconsideration = &ReconsiderationReceipt{}
			case "installed":
				r.InstallationStarted = true
			case "integrated":
				r.Integrated = true
			}
			s := NewService(taskrun.Dependencies{}, Options{PR: a})
			if err := s.prepareMergeRetry(&r); err == nil {
				t.Fatal("retry accepted without exact positive no-effect evidence")
			}
			if !r.IntegrationStarted || a.merges != 0 || len(r.MergeAttempts) != 0 {
				t.Fatal("rejected retry changed effect authority")
			}
		})
	}
}
