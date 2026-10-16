//go:build !windows

package taskrun

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// A claim is the immutable winning answer, not a wake. The question-keyed file
// is published before its request index. Both carry the same canonical digest;
// a partial publication is discoverable and blocks every other answer.
type productAnswerClaim struct {
	Request        BoundAnswerRequest `json:"request"`
	RequestSHA256  string             `json:"request_sha256"`
	Basis          productAnswerBasis `json:"basis"`
	QuestionRecord FileBinding        `json:"question_record"`
	RecordedAtUTC  string             `json:"recorded_at_utc"`
}

func productAnswerRoot(root string) string {
	return filepath.Join(root, ".ply", "workflow-bound-answers", "v2")
}
func productClaimPath(root, run, question string) string {
	return filepath.Join(productAnswerRoot(root), "claims", strings.TrimPrefix(digest([]string{run, question}), "sha256:")+".json")
}
func productIndexPath(root, id string) string {
	return filepath.Join(productAnswerRoot(root), "requests", strings.TrimPrefix(hash([]byte(id)), "sha256:")+".json")
}
func productReadClaim(path string) (productAnswerClaim, error) {
	var c productAnswerClaim
	if e := readValue(path, 256<<10, &c); e != nil {
		return c, e
	}
	if c.RequestSHA256 != digest(c.Request) || c.Basis.Public.Question == nil || c.Basis.Public.Scope == nil || c.Basis.Public.Worker == nil || c.Basis.Public.ParentWeek == nil || c.Request.RunID != c.Basis.Public.RunID || c.Request.QuestionID != c.Basis.Public.Question.QuestionID || c.Request.Binding != productAnswerBinding(c.Basis, c.Request) || c.Basis.Public.Eligibility != "ready" || c.QuestionRecord.SHA256 != c.Basis.QuestionRecordSHA256 || path != productClaimPath(c.Request.Workspace, c.Request.RunID, c.Request.QuestionID) {
		return c, boundAnswerError("outcome_unknown", "The immutable answer claim is inconsistent.")
	}
	return c, nil
}
func productClaimForQuestion(root, run, question string) (productAnswerClaim, error) {
	return productReadClaim(productClaimPath(root, run, question))
}
func productAnswerIndexed(root string, c productAnswerClaim) (FileBinding, error) {
	var index FileBinding
	e := readValue(productIndexPath(root, c.Request.RequestID), 16<<10, &index)
	if e != nil || index.Locator != productClaimPath(root, c.Request.RunID, c.Request.QuestionID) || index.SHA256 != digest(c) {
		return index, boundAnswerError("outcome_unknown", "The immutable winner exists but its request index is incomplete or inconsistent.")
	}
	if _, e = workflowBound(index, 256<<10); e != nil {
		return index, e
	}
	return index, nil
}
func productClaimByID(root, id string) (*productAnswerClaim, error) {
	var index FileBinding
	e := readValue(productIndexPath(root, id), 16<<10, &index)
	if e == nil {
		raw, e := workflowBound(index, 256<<10)
		if e != nil {
			return nil, e
		}
		var c productAnswerClaim
		if e = decode(raw, 256<<10, &c); e != nil {
			return nil, e
		}
		if c.Request.RequestID != id || c.Request.Workspace != root {
			return nil, boundAnswerError("outcome_unknown", "Request index identity differs.")
		}
		checked, e := productReadClaim(index.Locator)
		if e != nil {
			return nil, e
		}
		return &checked, nil
	}
	if !os.IsNotExist(e) {
		return nil, e
	}
	dir := filepath.Join(productAnswerRoot(root), "claims")
	if e = physical(dir, false); os.IsNotExist(e) {
		return nil, nil
	} else if e != nil {
		return nil, e
	}
	entries, e := os.ReadDir(dir)
	if e != nil {
		return nil, e
	}
	if len(entries) > 16384 {
		return nil, boundAnswerError("outcome_unknown", "The bounded answer claim inventory cannot be inspected.")
	}
	var found *productAnswerClaim
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			return nil, boundAnswerError("outcome_unknown", "Unexpected answer claim inventory.")
		}
		c, e := productReadClaim(filepath.Join(dir, entry.Name()))
		if e != nil {
			return nil, e
		}
		if c.Request.RequestID == id {
			if found != nil {
				return nil, boundAnswerError("outcome_unknown", "Request identity is ambiguous.")
			}
			copy := c
			found = &copy
		}
	}
	return found, nil
}
func productStatusBase(d Dependencies, root, id string) BoundAnswerStatus {
	return BoundAnswerStatus{Kind: "ply.workflow.bound-answer-status", SchemaVersion: 2, Mode: BoundAnswerMode, Workspace: root, RequestID: id, Action: "answer_product_question", State: "outcome_unknown", ExecutionEffect: "unknown", Reasons: []BoundAnswerReason{}, ObservedAtUTC: ControlTimestamp(controlNow(d)), Freshness: "unknown", Evidence: []BoundAnswerEvidence{}}
}
func productUnknownStatus(out BoundAnswerStatus, e error) BoundAnswerStatus {
	out.State = "outcome_unknown"
	out.ExecutionEffect = "unknown"
	out.Freshness = "retained"
	out.Reasons = append(out.Reasons, BoundAnswerFailureOf(e))
	return out
}
func productClaimStatus(d Dependencies, root string, c productAnswerClaim) BoundAnswerStatus {
	r, b := c.Request, c.Basis.Public
	o := productStatusBase(d, root, r.RequestID)
	o.RequestSHA256 = &c.RequestSHA256
	o.Binding = &r.Binding
	o.RunID = &r.RunID
	o.Scope = b.Scope
	o.ParentWeek = b.ParentWeek
	o.Worker = b.Worker
	o.OriginalDeadlineUTC = &r.OriginalDeadlineUTC
	o.QuestionID = &r.QuestionID
	o.QuestionRevision = &r.QuestionRevision
	o.QuestionSHA256 = &b.Question.QuestionSHA256
	o.OptionID = &r.OptionID
	o.Evidence = append(o.Evidence, b.Evidence...)
	if _, e := productAnswerIndexed(root, c); e != nil {
		return productUnknownStatus(o, e)
	}
	o.State = "answer_recorded"
	o.AnswerRecorded = ptr(true)
	o.ContinuationRequested = ptr(false)
	o.ExecutionEffect = "not_confirmed"
	o.Freshness = "retained"
	o.Receipt = &BoundAnswerReceipt{RequestID: r.RequestID, RequestSHA256: c.RequestSHA256, QuestionID: r.QuestionID, QuestionRevision: r.QuestionRevision, QuestionSHA256: b.Question.QuestionSHA256, OptionID: r.OptionID, SourceJobID: b.Worker.SourceJobID, SourceIntentSHA256: b.Worker.SourceIntentSHA256, SourceSessionID: b.Worker.ProviderSessionID, RecordedAtUTC: c.RecordedAtUTC}
	o.Evidence = append(o.Evidence, BoundAnswerEvidence{r.RequestID, "answer", c.RequestSHA256, c.RecordedAtUTC})
	return productContinuationStatus(d, c, o)
}
func ShowBoundAnswer(d Dependencies, root, requestID string) (BoundAnswerStatus, error) {
	o := productStatusBase(d, root, requestID)
	if physical(root, false) != nil || !boundAnswerRequestID.MatchString(requestID) {
		return o, boundAnswerError("invalid_arguments", "Expected a physical workspace and exact request ID.")
	}
	c, e := productClaimByID(root, requestID)
	if e != nil {
		return productUnknownStatus(o, e), nil
	}
	if c == nil {
		o.State = "not_recorded"
		o.Reasons = []BoundAnswerReason{BoundAnswerFailureOf(boundAnswerError("not_recorded", "No authoritative request record was found; this is not proof of no earlier effects."))}
		return o, nil
	}
	return productClaimStatus(d, root, *c), nil
}
func productRejected(d Dependencies, r BoundAnswerRequest, b *productAnswerBasis, e error) BoundAnswerStatus {
	o := productStatusBase(d, r.Workspace, r.RequestID)
	o.State = "rejected"
	o.RequestSHA256 = ptr(digest(r))
	o.Binding = &r.Binding
	o.RunID = &r.RunID
	o.QuestionID = &r.QuestionID
	o.QuestionRevision = &r.QuestionRevision
	o.OptionID = &r.OptionID
	o.OriginalDeadlineUTC = &r.OriginalDeadlineUTC
	o.AnswerRecorded = ptr(false)
	o.ContinuationRequested = ptr(false)
	o.ExecutionEffect = "not_confirmed"
	o.Freshness = "current"
	o.Reasons = []BoundAnswerReason{BoundAnswerFailureOf(e)}
	if b != nil {
		o.Scope = b.Public.Scope
		o.ParentWeek = b.Public.ParentWeek
		o.Worker = b.Public.Worker
		if b.Public.Question != nil {
			o.QuestionSHA256 = &b.Public.Question.QuestionSHA256
		}
		o.Evidence = b.Public.Evidence
	}
	return o
}
func ApplyBoundAnswer(d Dependencies, request BoundAnswerRequest) (BoundAnswerStatus, error) {
	raw, e := Canonical(request)
	if e != nil {
		return BoundAnswerStatus{}, boundAnswerError("invalid_request", "Cannot canonicalize request.")
	}
	r, e := ParseBoundAnswerRequest(raw)
	if e != nil {
		return BoundAnswerStatus{}, e
	}
	root, id := r.Workspace, r.RequestID
	// Exact canonical replay precedes freshness and any lock/directory creation.
	replay := func() (BoundAnswerStatus, bool, error) {
		c, e := productClaimByID(root, id)
		if e != nil {
			return productUnknownStatus(productStatusBase(d, root, id), e), true, nil
		}
		if c == nil {
			return BoundAnswerStatus{}, false, nil
		}
		out := productClaimStatus(d, root, *c)
		if c.RequestSHA256 != digest(r) {
			out.Reasons = append(out.Reasons, BoundAnswerFailureOf(boundAnswerError("request_conflict", "This request ID already names different canonical bytes.")))
			return out, true, boundAnswerError("request_conflict", "This request ID already names different canonical bytes.")
		}
		return out, true, nil
	}
	if out, found, e := replay(); found {
		return out, e
	}
	s, e := workflowRead(root, r.RunID)
	if e != nil {
		e := boundAnswerError("not_found", "Native Run is unavailable.")
		return productRejected(d, r, nil, e), e
	}
	if !productQuestionsEnabled(s) {
		e := boundAnswerError("unsupported", "This Run did not opt into product answers.")
		return productRejected(d, r, nil, e), e
	}
	var out BoundAnswerStatus
	e = WithWorkerParentGuard(root, s.Request.ControlContract, func() error {
		return withStore(root, func() error {
			if existing, found, e := replay(); found {
				out = existing
				return e
			}
			b, e := productRead(d, root, r.RunID)
			if e != nil {
				out = productRejected(d, r, nil, e)
				return e
			}
			issued, _ := boundAnswerTime(r.Validity.IssuedAtUTC)
			expires, _ := boundAnswerTime(r.Validity.ExpiresAtUTC)
			now := controlNow(d)
			if now.Before(issued) || !now.Before(expires) {
				e := boundAnswerError("expired", "The first apply must be within the approved validity window.")
				out = productRejected(d, r, &b, e)
				return e
			}
			if e = productRequestMatches(r, b); e != nil {
				out = productRejected(d, r, &b, e)
				return e
			}
			current, e := workflowRead(root, r.RunID)
			if e != nil {
				out = productRejected(d, r, &b, e)
				return e
			}
			c := productAnswerClaim{Request: r, RequestSHA256: digest(r), Basis: b, QuestionRecord: *current.ProductQuestion, RecordedAtUTC: ControlTimestamp(now)}
			path := productClaimPath(root, r.RunID, r.QuestionID)
			if e = d.fault("bound_answer_before_claim"); e != nil {
				return e
			}
			now = controlNow(d)
			if now.Before(issued) || !now.Before(expires) {
				e := boundAnswerError("expired", "The first answer publication crossed its approved validity window.")
				out = productRejected(d, r, &b, e)
				return e
			}
			c.RecordedAtUTC = ControlTimestamp(now)
			binding, e := workflowKeep(d, path, c)
			if e != nil {
				return e
			}
			if e = d.fault("bound_answer_after_claim"); e != nil {
				out = productUnknownStatus(productStatusBase(d, root, id), e)
				return e
			}
			if e = d.writeValue(productIndexPath(root, id), binding); e != nil {
				return e
			}
			if e = d.fault("bound_answer_after_index"); e != nil {
				return e
			}
			out = productClaimStatus(d, root, c)
			return nil
		})
	})
	if e != nil {
		if c, lookup := productClaimByID(root, id); lookup == nil && c != nil {
			return productClaimStatus(d, root, *c), e
		}
		if out.Kind != "" {
			return out, e
		}
		return productUnknownStatus(productStatusBase(d, root, id), e), e
	}
	// The retained native owner polls this durable answer. Apply never starts a
	// replacement owner or provider, and replay never wakes or repairs anything.
	return out, nil
}

func productContinuationDirectory(c productAnswerClaim) string {
	return filepath.Join(filepath.Dir(filepath.Dir(c.QuestionRecord.Locator)), "..", "attempts", productContinuationJobID(c))
}
func productContinuationJobID(c productAnswerClaim) string {
	return fmt.Sprintf("developer-%04d-answer-%s", c.Basis.Public.Worker.Iteration, strings.TrimPrefix(c.RequestSHA256, "sha256:")[:16])
}
