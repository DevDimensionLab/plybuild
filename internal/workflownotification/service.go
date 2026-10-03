package workflownotification

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type prepared struct {
	uncertain        bool
	input            Input
	route            Route
	routeSHA, secret string
	db               *database
	prior            *record
	current          *record
	preview          Preview
}

func localError(e error) error {
	if _, ok := e.(*Error); ok {
		return e
	}
	return fail(1, "local_io", "Local notification storage could not be read or written.")
}
func conflict(code, message string) error { return fail(4, code, message) }
func sourceError(e error, existing, apply bool) error {
	if existing || apply {
		return conflict("source_changed", "The confirmed source is no longer current.")
	}
	return fail(2, "invalid_source", "Source files or their identities could not be validated.")
}
func noEffects() []string { return []string{} }
func prepare(d Dependencies, in Input) (*prepared, error) {
	p := &prepared{input: in}
	routePath, e := filepath.Abs(in.Route)
	if e != nil {
		return p, invalid()
	}
	routeBytes, e := readFile(routePath, 16<<10, false)
	if e != nil {
		return p, sourceError(e, false, false)
	}
	p.route, e = parseRoute(routeBytes, "")
	if e != nil {
		return p, e
	}
	p.routeSHA = Digest(routeBytes)
	p.db, e = loadDatabase(p.route)
	if e != nil {
		if typed, ok := localError(e).(*Error); ok && typed.Code != "invalid_state_parent" {
			p.uncertain = true
		}
		return p, localError(e)
	}
	if in.Operation == "show" || in.Operation == "retry" {
		if !idPattern.MatchString(in.ID) || p.db == nil || p.db.Records[in.ID] == nil {
			return p, conflict("unknown_id", "No notification with this ID exists in this route.")
		}
		p.prior = p.db.Records[in.ID]
		p.current = p.prior
		if in.Operation == "show" {
			return p, nil
		}
	}
	var raw []byte
	var req Request
	var requestPath string
	if in.Operation == "send" {
		requestPath, e = filepath.Abs(in.File)
		if e != nil {
			return p, invalid()
		}
		raw, e = readFile(requestPath, 64<<10, false)
		if e != nil {
			return p, sourceError(e, false, in.Apply)
		}
		req, e = parseRequest(raw, "")
		if e != nil {
			return p, e
		}
		if p.db != nil {
			p.prior = p.db.Records[identity(req)]
		}
	} else {
		raw = p.prior.RequestBytes
		req = p.prior.Request
		requestPath = p.prior.RequestPath
	}
	p.secret, _ = d.LookupEnv(p.route.WebhookEnv)
	if !credential(p.secret) {
		return p, fail(2, "invalid_credential", "The named webhook credential is missing or invalid.")
	}
	if _, e = parseRoute(routeBytes, p.secret); e != nil {
		return p, e
	}
	if _, e = parseRequest(raw, p.secret); e != nil {
		return p, e
	}
	if req.Route != p.route.Name {
		return p, invalid()
	}
	if req.Source.Worktree == p.route.StateRoot {
		return p, invalid()
	}
	for _, path := range []string{req.Source.Handoff.Path, req.Source.Start.Path, req.Source.Report.Path, requestPath, routePath, req.Source.Worktree} {
		if strings.HasPrefix(path, p.route.StateRoot+string(filepath.Separator)) {
			return p, invalid()
		}
	}
	if e = checkSources(req); e != nil {
		return p, sourceError(e, p.prior != nil, in.Apply)
	}
	// Retry is bound to the original file too, not just a deserialized copy.
	if in.Operation == "retry" {
		current, e := readFile(requestPath, 64<<10, false)
		if e != nil || Digest(current) != Digest(raw) {
			return p, conflict("request_changed", "The original request file has changed or is unavailable.")
		}
	}
	payload := message(req)
	basis := digestValue([]any{Digest(raw), requestPath, p.routeSHA, p.route.StateRoot, Digest([]byte(p.secret)), digestValue(payload), req.Source})
	p.current = &record{RequestBytes: raw, RequestPath: requestPath, Request: req, Route: p.route, RouteSHA: p.routeSHA, CredentialSHA: Digest([]byte(p.secret)), Basis: basis, Payload: payload, Attempts: []Attempt{}}
	if p.db != nil && p.db.RouteSHA != p.routeSHA {
		return p, conflict("route_changed", "Route bytes differ from the original state binding.")
	}
	if p.prior != nil && p.prior.Basis != basis {
		return p, conflict("binding_changed", "The immutable notification or credential binding has changed.")
	}
	id := identity(req)
	p.preview = Preview{Kind: "ply.workflow.notification-preview", SchemaVersion: 1, ID: id, RequestSHA: Digest(raw), Operation: in.Operation, Allowed: true, Reasons: []Reason{}, Route: RouteView{p.route.Name, p.route.ChannelLabel, "claimed", p.routeSHA}, StateRoot: p.route.StateRoot, NextAttempt: 1, Payload: payload, PayloadSHA: digestValue(payload), Effects: []string{"Preserve a private durable attempt reservation and receipt.", "Perform at most one Slack POST."}, SourceKind: req.Source.Kind}
	if p.prior != nil {
		s := observedState(p.prior)
		p.preview.ExistingState = &s
	}
	if in.Operation == "send" && p.prior != nil {
		p.preview.Confirmation = &p.prior.SendConfirmation
		p.preview.Effects = noEffects()
		p.current = p.prior
		return p, nil
	}
	if p.db != nil {
		for otherID, r := range p.db.Records {
			if otherID != id && gateFamily(r.Request) == gateFamily(req) && observedState(r) == "unknown" {
				return blocked(p, "gate_unknown", "An earlier revision of this gate has unknown delivery; do not resend.")
			}
		}
	}
	var attempts []Attempt
	var deadline *string
	if in.Operation == "retry" {
		p.current = p.prior
		attempts = p.prior.Attempts
		p.preview.NextAttempt = len(attempts) + 1
		last := attempts[len(attempts)-1]
		deadline = last.RetryNotBefore
		// Previously registered retry confirmations are readback authority, even when
		// the current terminal state cannot be retried again.
		if in.Apply {
			for _, a := range attempts[1:] {
				if in.Confirm == a.Confirmation {
					p.preview.Confirmation = &a.Confirmation
					p.preview.Effects = noEffects()
					return p, nil
				}
			}
		}
		state := observedState(p.prior)
		if state != "rejected" && state != "not_sent" && state != "rate_limited" {
			return blocked(p, "retry_forbidden", "Only a documented rejection or proven non-delivery permits retry. Unknown and acknowledged delivery cannot be retried.")
		}
		if state == "rate_limited" {
			if deadline == nil {
				return p, fail(1, "state_corrupt", "A rate-limited attempt has no valid deadline.")
			}
			t, e := time.Parse(time.RFC3339, *deadline)
			if e != nil {
				return p, fail(1, "state_corrupt", "A rate-limited attempt has no valid deadline.")
			}
			if d.Now().Before(t) {
				return blocked(p, "retry_too_early", "The retry deadline has not been reached.")
			}
		}
	}
	confirmation := digestValue([]any{basis, in.Operation, p.preview.NextAttempt, attempts, deadline})
	p.preview.Confirmation = &confirmation
	if in.Operation == "send" {
		p.current.SendConfirmation = confirmation
	}
	return p, nil
}
func blocked(p *prepared, code, msg string) (*prepared, error) {
	p.preview.Allowed = false
	p.preview.Confirmation = nil
	p.preview.Effects = noEffects()
	p.preview.Reasons = []Reason{{code, msg}}
	return p, conflict(code, msg)
}
func freshness(p *prepared) string {
	r := p.prior
	if r == nil {
		return "unavailable"
	}
	changed := p.routeSHA != r.RouteSHA
	for path, want := range map[string]string{r.RequestPath: Digest(r.RequestBytes), r.Request.Source.Handoff.Path: r.Request.Source.Handoff.SHA256, r.Request.Source.Start.Path: r.Request.Source.Start.SHA256, r.Request.Source.Report.Path: r.Request.Source.Report.SHA256} {
		b, e := readFile(path, 16<<20, false)
		if e != nil {
			return "unavailable"
		}
		if Digest(b) != want {
			changed = true
		}
	}
	if changed {
		return "changed"
	}
	if checkSources(r.Request) != nil {
		return "changed"
	}
	return "current"
}
func failureOutput(p *prepared, err error) any {
	if p != nil && p.uncertain {
		id := p.input.ID
		var source any
		var gate any
		var requestSHA any
		if p.input.Operation == "send" {
			path, e := filepath.Abs(p.input.File)
			if e == nil {
				if raw, e := readFile(path, 64<<10, false); e == nil {
					if req, e := parseRequest(raw, ""); e == nil {
						id = identity(req)
						source = req.Source
						gate = req.Gate
						requestSHA = Digest(raw)
					}
				}
			}
		}
		return map[string]any{"kind": "ply.workflow.notification", "schema_version": 1, "notification_id": id, "request_sha256": requestSHA, "source_binding": source, "route": RouteView{p.route.Name, p.route.ChannelLabel, "claimed", p.routeSHA}, "gate": gate, "knowledge": "reported", "state": "unknown", "freshness": "unavailable", "payload_sha256": nil, "attempts": []Attempt{}, "retry_not_before": nil, "reasons": []Reason{{"state_unavailable", "Existing attempt state cannot be established. No new receipt was created; do not resend."}}, "next_action": "Preserve the local return and inspect the existing state.", "persistence": "not_created"}
	}

	if p != nil && p.prior != nil {
		r := result(identity(p.prior.Request), p.prior, freshness(p))
		e := ErrorOutput(err)
		r.Reasons = append(r.Reasons, e.Reasons...)
		return r
	}
	return ErrorOutput(err)
}

// Execute returns a public value independently of the executable exit status.
func Execute(d Dependencies, in Input) (any, error) {
	if d.Now == nil || d.LookupEnv == nil || d.Transport == nil {
		return ErrorOutput(invalid()), invalid()
	}
	p, e := prepare(d, in)
	if e != nil {
		if p != nil && p.preview.Kind != "" && !p.preview.Allowed && !in.Apply {
			return p.preview, e
		}
		return failureOutput(p, e), e
	}
	if in.Operation == "show" {
		return result(in.ID, p.prior, freshness(p)), nil
	}
	if !in.Apply {
		return p.preview, nil
	}
	if p.preview.Confirmation == nil || in.Confirm != *p.preview.Confirmation {
		e = conflict("confirmation_changed", "Confirmation does not match the current notification basis.")
		return failureOutput(p, e), e
	}
	if len(p.preview.Effects) == 0 {
		return result(p.preview.ID, p.prior, "current"), nil
	}
	dir, e := stateDirectory(p.route.StateRoot, true)
	if e != nil {
		return failureOutput(p, localError(e)), localError(e)
	}
	defer dir.Close()
	unlock, e := lockState(dir)
	if e != nil {
		return failureOutput(p, localError(e)), localError(e)
	}
	defer unlock()
	// Revalidate under the interprocess lock. Concurrent winners become readbacks;
	// source changes never inherit the already-checked confirmation.
	p, e = prepare(d, in)
	if e != nil {
		return failureOutput(p, e), e
	}
	if p.preview.Confirmation == nil || in.Confirm != *p.preview.Confirmation {
		e = conflict("confirmation_changed", "Confirmation changed before reservation.")
		return failureOutput(p, e), e
	}
	if len(p.preview.Effects) == 0 {
		return result(p.preview.ID, p.prior, "current"), nil
	}
	currentDir, e := openDirectory(p.route.StateRoot)
	if e != nil {
		return failureOutput(p, localError(e)), localError(e)
	}
	a, _ := dir.Stat()
	b, _ := currentDir.Stat()
	currentDir.Close()
	if a == nil || b == nil || !os.SameFile(a, b) {
		e = conflict("state_changed", "The physical state root changed.")
		return failureOutput(p, e), e
	}
	if d.Fault != nil {
		if e = d.Fault("before_reservation"); e != nil {
			return failureOutput(p, localError(e)), localError(e)
		}
	}
	if p.db == nil {
		p.db = newDatabase(p.route, p.routeSHA)
	}
	r := p.current
	r.Attempts = append(append([]Attempt(nil), r.Attempts...), Attempt{Number: p.preview.NextAttempt, State: "reserved", ReservedAt: stamp(d.Now()), Confirmation: in.Confirm, Dispatch: "not_started"})
	p.db.Records[p.preview.ID] = r
	if e = writeDatabase(dir, p.db); e != nil {
		// A failed directory sync may follow a successful rename. Inspect that fact
		// conservatively; never equate a write error with absence of a reservation.
		saved, readErr := readDatabase(dir)
		if readErr == nil && saved != nil && saved.Records[p.preview.ID] != nil {
			return result(p.preview.ID, saved.Records[p.preview.ID], "current"), localError(e)
		}
		out := result(p.preview.ID, r, "current")
		out.Persistence = "not_created"
		out.State = "unknown"
		return out, localError(e)
	}
	if e = bindRoot(dir, p.route.Name, p.routeSHA); e != nil {
		return result(p.preview.ID, r, "current"), localError(e)
	}
	if d.Fault != nil {
		if e = d.Fault("after_reservation"); e != nil {
			return result(p.preview.ID, r, "current"), localError(e)
		}
	}

	// Recheck the paths and immutable input after the durable reservation. A path
	// change at this boundary must preserve uncertainty without dispatching.
	guard, guardErr := prepare(d, in)
	latestDir, dirErr := openDirectory(p.route.StateRoot)
	sameRoot := false
	if dirErr == nil {
		oldInfo, _ := dir.Stat()
		newInfo, _ := latestDir.Stat()
		sameRoot = oldInfo != nil && newInfo != nil && os.SameFile(oldInfo, newInfo)
		latestDir.Close()
	}
	if guardErr != nil || !sameRoot || guard == nil || guard.prior == nil || guard.prior.Basis != r.Basis {
		out := result(p.preview.ID, r, "changed")
		out.Reasons = append(out.Reasons, Reason{"binding_changed", "The reserved source or state root changed before dispatch."})
		return out, conflict("binding_changed", "The reserved source or state root changed before dispatch.")
	}
	data, _ := json.Marshal(p.preview.Payload)
	observation := d.Transport(context.Background(), p.secret, data)
	if d.Fault != nil {
		if e = d.Fault("after_transport"); e != nil {
			return result(p.preview.ID, r, "current"), localError(e)
		}
	}
	finished := classify(observation, d.Now())
	last := &r.Attempts[len(r.Attempts)-1]
	finished.Number = last.Number
	finished.ReservedAt = last.ReservedAt
	finished.Confirmation = last.Confirmation
	*last = finished
	if d.Fault != nil {
		e = d.Fault("receipt_commit")
	}
	if e == nil {
		e = writeDatabase(dir, p.db)
	}
	if e != nil {
		last.State = "reserved"
		last.CompletedAt = nil
		last.HTTPStatus = nil
		last.ResponseCode = nil
		last.RetryNotBefore = nil
		out := result(p.preview.ID, r, "current")
		return out, fail(1, "receipt_not_stored", "Delivery is unknown because the completed receipt could not be durably preserved.")
	}
	out := result(p.preview.ID, r, "current")
	if out.State != "transport_acknowledged" {
		return out, fail(5, "transport_incomplete", "The attempt is preserved; transport was not acknowledged.")
	}
	return out, nil
}
