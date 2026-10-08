package integration

import "fmt"

func exactRetryPR(t PRMergeTarget, o PRMergeObservation) bool {
	return o.Repository == t.Repository && o.Number == t.Number && o.HeadRef == t.HeadRef && o.HeadOID == t.HeadOID && o.BaseRef == t.BaseRef
}

// prepareMergeRetry is called only by an explicit resume --retry-merge after
// current native QA validation. It preserves the old request before clearing
// its operation ID, and never derives non-delivery from an open PR alone.
func (s *Service) prepareMergeRetry(r *Receipt) error {
	if r.Reconsideration != nil || r.Integrated || r.InstallationStarted || r.Closeout != nil || r.Plan.Flow != "github_pr_merge" || r.Plan.PR == nil || !r.IntegrationStarted || r.PR == nil {
		return fmt.Errorf("retry_forbidden: there is no failed PR attempt eligible for retry")
	}
	target := *r.Plan.PR
	target.OperationID = r.PR.OperationID
	observed, err := s.options.PR.Observe(target)
	if err != nil {
		return err
	}
	if !exactRetryPR(target, observed) || (observed.State != "ready" && observed.State != "blocked") {
		return fmt.Errorf("retry_forbidden: the exact PR has no proven current unmerged outcome; resume observation without retry")
	}
	prior := *r.PR
	if !prior.NoEffect {
		if !observed.NoEffect {
			return fmt.Errorf("retry_forbidden: an open PR does not prove that the previous request had no effect")
		}
		prior = observed
	}
	if !exactRetryPR(target, prior) || (prior.State != "ready" && prior.State != "blocked") {
		return fmt.Errorf("retry_forbidden: the preserved rejection does not prove this exact request had no effect")
	}
	target.OperationID = ""
	ready, err := s.options.PR.Observe(target)
	if err != nil {
		return err
	}
	if ready.State != "ready" || !exactRetryPR(target, ready) {
		return fmt.Errorf("merge_blocked: the exact frozen PR/method must pass a fresh preflight before retry")
	}
	r.MergeAttempts = append(r.MergeAttempts, prior)
	r.MergeRetryAtUTC = s.now()
	r.PR = &ready
	r.IntegrationStarted = false
	r.State, r.Reasons, r.NextAction = "retry_authorized", []string{}, "Perform one explicitly requested retry of the unchanged PR integration plan."
	return s.save(r)
}
