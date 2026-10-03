package workflownotification

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

type record struct {
	RequestBytes     []byte    `json:"request_bytes"`
	RequestPath      string    `json:"request_path"`
	Request          Request   `json:"request"`
	Route            Route     `json:"route"`
	RouteSHA         string    `json:"route_sha256"`
	CredentialSHA    string    `json:"credential_sha256"`
	Basis            string    `json:"basis"`
	SendConfirmation string    `json:"send_confirmation"`
	Payload          Payload   `json:"payload"`
	Attempts         []Attempt `json:"attempts"`
}
type database struct {
	Integrity string             `json:"integrity_sha256"`
	Kind      string             `json:"kind"`
	Route     string             `json:"route"`
	RouteSHA  string             `json:"route_sha256"`
	Records   map[string]*record `json:"records"`
}

func timeNonce() int64 { return time.Now().UnixNano() }
func newDatabase(route Route, sha string) *database {
	return &database{Kind: "ply.workflow.notification-state", Route: route.Name, RouteSHA: sha, Records: map[string]*record{}}
}
func readDatabase(dir *os.File) (*database, error) {
	b, e := readAt(dir, "state.json", 64<<20, true)
	if errors.Is(e, os.ErrNotExist) {
		entries, err := dir.ReadDir(-1)
		if err != nil {
			return nil, err
		}
		// The committed file can appear between the failed open and this
		// directory snapshot. That is a concurrent writer, not corruption.
		for _, entry := range entries {
			if entry.Name() == "state.json" {
				return nil, fail(4, "state_changed", "Notification state appeared during readback; inspect it again.")
			}
		}
		for _, entry := range entries {
			if entry.Name() != ".lock" && !strings.HasPrefix(entry.Name(), ".pending-") {
				return nil, fail(1, "state_corrupt", "Notification state is incomplete or corrupt.")
			}
		}
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	var db database
	if _, e = strict(b, &db); e != nil || db.Kind != "ply.workflow.notification-state" || !routePattern.MatchString(db.Route) || !digestPattern.MatchString(db.RouteSHA) || db.Records == nil {
		return nil, fail(1, "state_corrupt", "Notification state is incomplete or corrupt.")
	}
	integrity := db.Integrity
	db.Integrity = ""
	if integrity != digestValue(db) {
		return nil, fail(1, "state_corrupt", "Notification state integrity does not match.")
	}
	db.Integrity = integrity
	binding, bindErr := readAt(dir, "binding", 1024, true)
	if bindErr != nil && !errors.Is(bindErr, os.ErrNotExist) {
		return nil, bindErr
	}
	if bindErr == nil && string(binding) != db.Route+"\n"+db.RouteSHA+"\n" {
		return nil, fail(4, "state_changed", "The state root binding is incomplete or changed.")
	}
	for id, r := range db.Records {
		if r == nil || id != identity(r.Request) || r.Route.Name != db.Route || r.RouteSHA != db.RouteSHA || r.Basis == "" || !digestPattern.MatchString(r.CredentialSHA) || !digestPattern.MatchString(r.SendConfirmation) || len(r.Attempts) == 0 || digestValue(r.Payload) != digestValue(message(r.Request)) {
			return nil, fail(1, "state_corrupt", "Notification state is incomplete or corrupt.")
		}
		req, e := parseRequest(r.RequestBytes, "")
		if e != nil || digestValue(req) != digestValue(r.Request) {
			return nil, fail(1, "state_corrupt", "Notification state is incomplete or corrupt.")
		}
		for i, a := range r.Attempts {
			if a.Number != i+1 || !digestPattern.MatchString(a.Confirmation) || a.ReservedAt == "" || !validState(a.State) || a.Dispatch != "not_started" && a.Dispatch != "started" && a.Dispatch != "unknown" || a.State == "reserved" && a.CompletedAt != nil || a.State != "reserved" && a.CompletedAt == nil {
				return nil, fail(1, "state_corrupt", "Notification state is incomplete or corrupt.")
			}
		}
	}
	return &db, nil
}
func validState(s string) bool {
	switch s {
	case "reserved", "transport_acknowledged", "rejected", "not_sent", "rate_limited", "unknown":
		return true
	}
	return false
}
func loadDatabase(route Route) (*database, error) {
	dir, e := stateDirectory(route.StateRoot, false)
	if errors.Is(e, os.ErrNotExist) {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	defer dir.Close()
	db, e := readDatabase(dir)
	if e != nil {
		return nil, e
	}
	if db != nil && db.Route != route.Name {
		return nil, fail(4, "route_conflict", "This state root belongs to a different route.")
	}
	return db, nil
}
func writeDatabase(dir *os.File, db *database) error {
	db.Integrity = ""
	db.Integrity = digestValue(db)
	b, e := json.Marshal(db)
	if e != nil {
		return e
	}
	return atomicState(dir, b)
}
func observedState(r *record) string {
	s := r.Attempts[len(r.Attempts)-1].State
	if s == "reserved" {
		return "unknown"
	}
	return s
}
func result(id string, r *record, freshness string) Result {
	attempts := append([]Attempt(nil), r.Attempts...)
	last := &attempts[len(attempts)-1]
	reasons := []Reason{}
	persistence := "stored"
	if last.State == "reserved" {
		last.State = "unknown"
		last.Dispatch = "unknown"
		reasons = append(reasons, Reason{"attempt_incomplete", "A durable reservation has no completed receipt; delivery is unknown. Do not resend."})
		persistence = "reservation_only"
	}
	if last.State == "unknown" && len(reasons) == 0 {
		reasons = append(reasons, Reason{"transport_unconfirmed", "Delivery is unknown. Preserve the local return and check Slack separately; do not resend."})
	}
	return Result{"ply.workflow.notification", 1, id, Digest(r.RequestBytes), r.Request.Source, RouteView{r.Route.Name, r.Route.ChannelLabel, "claimed", r.RouteSHA}, r.Request.Gate, "reported", last.State, freshness, digestValue(r.Payload), attempts, last.RetryNotBefore, reasons, r.Request.Public.NextAction, persistence}
}
