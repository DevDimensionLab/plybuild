package workflownotification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// IntegrationEvent is emitted by the native integration service only after an
// exact Git/PR effect has been observed. It is deliberately not an agent event.
type IntegrationEvent struct {
	Kind          string   `json:"kind"`
	SchemaVersion int      `json:"schema_version"`
	ID            string   `json:"id"`
	IntegrationID string   `json:"integration_id"`
	TaskID        string   `json:"task_id"`
	DeliveryID    string   `json:"delivery_id,omitempty"`
	CandidateOID  string   `json:"candidate_oid"`
	TargetRef     string   `json:"target_ref"`
	IntegratedOID string   `json:"integrated_oid"`
	ObservedAtUTC string   `json:"observed_at_utc"`
	PRURL         string   `json:"pr_url,omitempty"`
	Remaining     []string `json:"remaining"`
}

// Route and its exact source bytes are frozen in the displayed integration
// plan. This namespace leaves the existing agent/human-gate state untouched.
type IntegrationRouteBinding struct {
	Path      string `json:"path"`
	SHA256    string `json:"sha256"`
	Route     Route  `json:"route"`
	StateRoot string `json:"state_root"`
}

type IntegrationObservation struct {
	ID             string    `json:"id"`
	EventID        string    `json:"event_id"`
	State          string    `json:"state"`
	Reason         string    `json:"reason,omitempty"`
	NextAction     string    `json:"next_action"`
	Attempts       []Attempt `json:"attempts"`
	RetryNotBefore *string   `json:"retry_not_before,omitempty"`
}

type integrationTransportRecord struct {
	Event         IntegrationEvent `json:"event"`
	RouteSHA256   string           `json:"route_sha256"`
	CredentialSHA string           `json:"credential_sha256"`
	Payload       Payload          `json:"payload"`
	Attempts      []Attempt        `json:"attempts"`
}

type integrationTransportState struct {
	Kind          string                                 `json:"kind"`
	SchemaVersion int                                    `json:"schema_version"`
	Integrity     string                                 `json:"integrity_sha256"`
	Route         string                                 `json:"route"`
	RouteSHA256   string                                 `json:"route_sha256"`
	Records       map[string]*integrationTransportRecord `json:"records"`
}

var integrationOIDPattern = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)

func LoadIntegrationRoute(path string) (IntegrationRouteBinding, error) {
	var b IntegrationRouteBinding
	if !validPath(path) {
		return b, invalid()
	}
	raw, err := readFile(path, 16<<10, false)
	if err != nil {
		return b, err
	}
	route, err := parseRoute(raw, "")
	if err != nil {
		return b, err
	}
	b = IntegrationRouteBinding{Path: path, SHA256: Digest(raw), Route: route, StateRoot: route.StateRoot + "-integration-events"}
	// Read-only planning needs no credentials, locks, directory creation or POST.
	return b, nil
}

func validateIntegrationRoute(b IntegrationRouteBinding) error {
	current, err := LoadIntegrationRoute(b.Path)
	if err != nil {
		return err
	}
	if digestValue(current) != digestValue(b) {
		return conflict("route_changed", "The selected integration notification route changed; prepare a new plan.")
	}
	return nil
}

func validIntegrationEvent(e IntegrationEvent) bool {
	if e.Kind != "ply.integration.observed" || e.SchemaVersion != 1 || !cleanText(e.ID, 200) || !cleanText(e.IntegrationID, 200) || !cleanText(e.TaskID, 200) ||
		!integrationOIDPattern.MatchString(e.CandidateOID) || !integrationOIDPattern.MatchString(e.IntegratedOID) || !cleanText(e.TargetRef, 512) ||
		!strings.HasPrefix(e.TargetRef, "refs/heads/") || len(e.Remaining) > 16 {
		return false
	}
	if _, err := time.Parse(time.RFC3339Nano, e.ObservedAtUTC); err != nil {
		return false
	}
	for _, s := range append([]string{e.ID, e.IntegrationID, e.TaskID, e.DeliveryID, e.TargetRef, e.PRURL}, e.Remaining...) {
		if s != "" && (!cleanText(s, 1024) || !publicText(s) || webhookText.MatchString(s)) {
			return false
		}
	}
	if e.PRURL != "" && !strings.HasPrefix(e.PRURL, "https://github.com/") {
		return false
	}
	return true
}

func integrationNotificationID(e IntegrationEvent) string {
	return "nti_" + strings.TrimPrefix(digestValue([]string{e.Kind, e.ID}), "sha256:")
}

func integrationPayload(e IntegrationEvent) Payload {
	status := "Integrated; closeout complete."
	if len(e.Remaining) != 0 {
		status = "Integrated; remaining: " + strings.Join(e.Remaining, ", ") + "."
	}
	link := ""
	if e.PRURL != "" {
		link = "\n" + e.PRURL
	}
	return Payload{Text: fmt.Sprintf("Task %s integrated into %s\nCandidate %s; observed %s\n%s%s", e.TaskID, e.TargetRef, e.CandidateOID, e.IntegratedOID, status, link), Mrkdwn: false, Parse: "none", UnfurlLinks: false, UnfurlMedia: false}
}

func readIntegrationState(dir *os.File, binding IntegrationRouteBinding) (*integrationTransportState, error) {
	raw, err := readAt(dir, "state.json", 64<<20, true)
	if errors.Is(err, os.ErrNotExist) {
		entries, err := dir.ReadDir(-1)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.Name() != ".lock" && !strings.HasPrefix(entry.Name(), ".pending-") {
				return nil, fail(1, "state_corrupt", "Integration notification state is missing; do not resend.")
			}
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var db integrationTransportState
	if _, err = strict(raw, &db); err != nil || db.Kind != "ply.integration.notification-state" || db.SchemaVersion != 1 || db.Route != binding.Route.Name || db.RouteSHA256 != binding.SHA256 || db.Records == nil {
		return nil, fail(1, "state_corrupt", "Integration notification state or route binding is invalid.")
	}
	integrity := db.Integrity
	db.Integrity = ""
	if integrity != digestValue(db) {
		return nil, fail(1, "state_corrupt", "Integration notification state integrity is invalid.")
	}
	db.Integrity = integrity
	tombstone, err := readAt(dir, "binding", 1024, true)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if err == nil && string(tombstone) != binding.Route.Name+"\n"+binding.SHA256+"\n" {
		return nil, fail(1, "state_corrupt", "Integration notification root binding is invalid.")
	}
	for id, r := range db.Records {
		if r == nil || !validIntegrationEvent(r.Event) || id != integrationNotificationID(r.Event) || r.RouteSHA256 != binding.SHA256 || !digestPattern.MatchString(r.CredentialSHA) || len(r.Attempts) < 1 || digestValue(r.Payload) != digestValue(integrationPayload(r.Event)) {
			return nil, fail(1, "state_corrupt", "Integration notification record is invalid.")
		}
		for i, a := range r.Attempts {
			if a.Number != i+1 || !validState(a.State) || !digestPattern.MatchString(a.Confirmation) || a.ReservedAt == "" ||
				(a.State == "reserved") != (a.CompletedAt == nil) || (a.Dispatch != "started" && a.Dispatch != "not_started" && a.Dispatch != "unknown") {
				return nil, fail(1, "state_corrupt", "Integration notification attempt is invalid.")
			}
		}
	}
	return &db, nil
}

func writeIntegrationState(dir *os.File, db *integrationTransportState) error {
	db.Integrity = ""
	db.Integrity = digestValue(db)
	raw, err := json.Marshal(db)
	if err != nil {
		return err
	}
	return atomicState(dir, raw)
}

func integrationTransportObservation(e IntegrationEvent, r *integrationTransportRecord) IntegrationObservation {
	o := IntegrationObservation{ID: integrationNotificationID(e), EventID: e.ID, State: "not_sent", Reason: "not_attempted", NextAction: "Send the single notification selected in the confirmed integration plan.", Attempts: []Attempt{}}
	if r == nil {
		return o
	}
	o.Attempts = append([]Attempt(nil), r.Attempts...)
	last := o.Attempts[len(o.Attempts)-1]
	o.State, o.RetryNotBefore = last.State, last.RetryNotBefore
	switch last.State {
	case "reserved", "unknown":
		o.State, o.Reason, o.NextAction = "unknown", "notification_unknown", "Read the preserved notification state and inspect Slack separately; do not resend an unknown delivery."
	case "transport_acknowledged":
		o.Reason, o.NextAction = "", "Slack acknowledged this integration notification; no further send is needed."
	default:
		o.Reason, o.NextAction = "notification_"+last.State, "The observed integration remains valid; explicitly retry notification only after proven non-delivery."
	}
	return o
}

func unknownIntegrationTransport(e IntegrationEvent, err error) (IntegrationObservation, error) {
	o := integrationTransportObservation(e, nil)
	o.State, o.Reason, o.NextAction = "unknown", "notification_unknown", "Inspect the preserved integration notification state; do not resend until its delivery state is known."
	return o, err
}

// ObserveIntegration performs no writes, credential reads, locks or network
// access. A missing/corrupt receipt after reservation remains unknown.
func ObserveIntegration(binding IntegrationRouteBinding, event IntegrationEvent) (IntegrationObservation, error) {
	if !validIntegrationEvent(event) {
		return unknownIntegrationTransport(event, invalid())
	}
	if err := validateIntegrationRoute(binding); err != nil {
		return unknownIntegrationTransport(event, err)
	}
	dir, err := stateDirectory(binding.StateRoot, false)
	if errors.Is(err, os.ErrNotExist) {
		return integrationTransportObservation(event, nil), nil
	}
	if err != nil {
		return unknownIntegrationTransport(event, err)
	}
	defer dir.Close()
	db, err := readIntegrationState(dir, binding)
	if err != nil {
		return unknownIntegrationTransport(event, err)
	}
	if db == nil {
		return integrationTransportObservation(event, nil), nil
	}
	r := db.Records[integrationNotificationID(event)]
	if r != nil && digestValue(r.Event) != digestValue(event) {
		return unknownIntegrationTransport(event, conflict("event_changed", "The integration event changed; its original notification must not be replaced."))
	}
	return integrationTransportObservation(event, r), nil
}

// SendIntegration executes at most one transport attempt for an event. Exact
// repeated calls return the durable result. RetryIntegration is deliberately a
// separate operation and requires a fresh explicit confirmation of non-delivery.
func SendIntegration(d Dependencies, binding IntegrationRouteBinding, event IntegrationEvent) (IntegrationObservation, error) {
	return sendIntegration(d, binding, event, "")
}

// IntegrationRetryConfirmation is read-only. Its digest binds the exact event,
// route and preceding attempts, so an old confirmation cannot trigger a later
// attempt. An unknown or acknowledged send never produces a retry confirmation.
func IntegrationRetryConfirmation(d Dependencies, binding IntegrationRouteBinding, event IntegrationEvent) (string, error) {
	o, err := ObserveIntegration(binding, event)
	if err != nil {
		return "", err
	}
	if o.State != "not_sent" && o.State != "rejected" && o.State != "rate_limited" || len(o.Attempts) == 0 {
		return "", conflict("retry_forbidden", "Only a preserved, proven non-delivery can be retried.")
	}
	if o.State == "rate_limited" {
		if d.Now == nil || o.RetryNotBefore == nil {
			return "", invalid()
		}
		deadline, err := time.Parse(time.RFC3339, *o.RetryNotBefore)
		if err != nil || d.Now().Before(deadline) {
			return "", conflict("retry_too_early", "The notification retry deadline has not been reached.")
		}
	}
	return digestValue([]any{"integration-notification-retry", event, binding, o.Attempts}), nil
}

func RetryIntegration(d Dependencies, binding IntegrationRouteBinding, event IntegrationEvent, confirmation string) (IntegrationObservation, error) {
	if !digestPattern.MatchString(confirmation) {
		return unknownIntegrationTransport(event, invalid())
	}
	return sendIntegration(d, binding, event, confirmation)
}

func sendIntegration(d Dependencies, binding IntegrationRouteBinding, event IntegrationEvent, retry string) (IntegrationObservation, error) {
	if d.Now == nil || d.LookupEnv == nil || d.Transport == nil || !validIntegrationEvent(event) {
		return unknownIntegrationTransport(event, invalid())
	}
	prior, err := ObserveIntegration(binding, event)
	if err != nil {
		return prior, err
	}
	if len(prior.Attempts) > 0 && retry == "" {
		return prior, nil
	}
	secret, _ := d.LookupEnv(binding.Route.WebhookEnv)
	if !credential(secret) {
		return prior, fail(2, "invalid_credential", "The locally selected Slack credential is unavailable or invalid.")
	}
	if _, err = parseRoute(mustIntegrationJSON(binding.Route), secret); err != nil {
		return prior, err
	}
	dir, err := stateDirectory(binding.StateRoot, true)
	if err != nil {
		return unknownIntegrationTransport(event, err)
	}
	defer dir.Close()
	unlock, err := lockState(dir)
	if err != nil {
		return unknownIntegrationTransport(event, err)
	}
	defer unlock()
	if err = validateIntegrationRoute(binding); err != nil {
		return unknownIntegrationTransport(event, err)
	}
	db, err := readIntegrationState(dir, binding)
	if err != nil {
		return unknownIntegrationTransport(event, err)
	}
	if db == nil {
		db = &integrationTransportState{Kind: "ply.integration.notification-state", SchemaVersion: 1, Route: binding.Route.Name, RouteSHA256: binding.SHA256, Records: map[string]*integrationTransportRecord{}}
	}
	id := integrationNotificationID(event)
	r := db.Records[id]
	if r != nil {
		if digestValue(r.Event) != digestValue(event) || r.CredentialSHA != Digest([]byte(secret)) {
			return unknownIntegrationTransport(event, conflict("binding_changed", "The immutable integration notification binding changed."))
		}
		if retry == "" {
			return integrationTransportObservation(event, r), nil
		}
		for _, a := range r.Attempts {
			if a.Confirmation == retry {
				return integrationTransportObservation(event, r), nil
			}
		}
		confirmation, err := IntegrationRetryConfirmation(d, binding, event)
		if err != nil || confirmation != retry {
			if err == nil {
				err = conflict("confirmation_changed", "The retry confirmation changed.")
			}
			return integrationTransportObservation(event, r), err
		}
	} else {
		if retry != "" {
			return unknownIntegrationTransport(event, conflict("retry_forbidden", "There is no preserved attempt to retry."))
		}
		r = &integrationTransportRecord{Event: event, RouteSHA256: binding.SHA256, CredentialSHA: Digest([]byte(secret)), Payload: integrationPayload(event), Attempts: []Attempt{}}
	}
	confirmation := retry
	if confirmation == "" {
		confirmation = digestValue([]any{event, binding, r.Payload, r.CredentialSHA})
	}
	r.Attempts = append(r.Attempts, Attempt{Number: len(r.Attempts) + 1, State: "reserved", ReservedAt: stamp(d.Now()), Confirmation: confirmation, Dispatch: "not_started"})
	db.Records[id] = r
	if err = writeIntegrationState(dir, db); err != nil {
		return unknownIntegrationTransport(event, err)
	}
	if err = bindRoot(dir, binding.Route.Name, binding.SHA256); err != nil {
		return integrationTransportObservation(event, r), err
	}
	if d.Fault != nil {
		if err = d.Fault("integration_after_reservation"); err != nil {
			return integrationTransportObservation(event, r), err
		}
	}
	currentDir, dirErr := openDirectory(binding.StateRoot)
	if dirErr != nil {
		return integrationTransportObservation(event, r), dirErr
	}
	before, _ := dir.Stat()
	after, _ := currentDir.Stat()
	currentDir.Close()
	if before == nil || after == nil || !os.SameFile(before, after) {
		return integrationTransportObservation(event, r), conflict("state_changed", "The physical notification state root changed.")
	}
	if err = validateIntegrationRoute(binding); err != nil {
		return integrationTransportObservation(event, r), err
	}
	observed := d.Transport(context.Background(), secret, mustIntegrationJSON(r.Payload))
	if d.Fault != nil {
		if err = d.Fault("integration_after_send"); err != nil {
			return integrationTransportObservation(event, r), err
		}
	}
	last := len(r.Attempts) - 1
	completed := classify(observed, d.Now())
	completed.Number, completed.ReservedAt, completed.Confirmation = r.Attempts[last].Number, r.Attempts[last].ReservedAt, r.Attempts[last].Confirmation
	r.Attempts[last] = completed
	if err = writeIntegrationState(dir, db); err != nil {
		return unknownIntegrationTransport(event, err)
	}
	return integrationTransportObservation(event, r), nil
}

func mustIntegrationJSON(v any) []byte { raw, _ := json.Marshal(v); return raw }

// IntegrationStatePath is the observable durable namespace, useful when a
// controller preserves route evidence before retiring the source worktree.
func IntegrationStatePath(binding IntegrationRouteBinding) string {
	return filepath.Join(binding.StateRoot, "state.json")
}
