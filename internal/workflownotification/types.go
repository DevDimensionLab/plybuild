// Package workflownotification delivers an explicitly claimed report-ready gate.
// Transport acknowledgement never attests product correctness or human approval.
package workflownotification

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Reason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type Error struct {
	Exit          int
	Code, Message string
}

func (e *Error) Error() string { return e.Message }
func fail(exit int, code, message string) error {
	return &Error{Exit: exit, Code: code, Message: message}
}
func Digest(b []byte) string   { h := sha256.Sum256(b); return "sha256:" + hex.EncodeToString(h[:]) }
func digestValue(v any) string { b, _ := json.Marshal(v); return Digest(b) }

type Locator struct {
	Path   string  `json:"path"`
	SHA256 string  `json:"sha256"`
	Kind   *string `json:"kind,omitempty"`
}
type Source struct {
	Kind      string          `json:"kind"`
	Activity  string          `json:"activity"`
	Run       string          `json:"run"`
	Worktree  string          `json:"worktree"`
	Handoff   Locator         `json:"handoff"`
	Start     Locator         `json:"start_receipt"`
	Report    Locator         `json:"report"`
	HandoffID *string         `json:"handoff_id,omitempty"`
	TaskID    json.RawMessage `json:"task_id,omitempty"`
}
type Gate struct {
	ID       string `json:"id"`
	Revision int    `json:"revision"`
	OpenedAt string `json:"opened_at"`
	Reason   string `json:"reason"`
	State    string `json:"state"`
}
type Public struct {
	TaskTitle  string `json:"task_title"`
	NextAction string `json:"next_action"`
}
type Request struct {
	Kind          string `json:"kind"`
	SchemaVersion int    `json:"schema_version"`
	Route         string `json:"route"`
	Source        Source `json:"source"`
	Gate          Gate   `json:"gate"`
	Sender        struct {
		ActorClaim string `json:"actor_claim"`
	} `json:"sender"`
	Public Public `json:"public"`
}
type Route struct {
	Kind          string `json:"kind"`
	SchemaVersion int    `json:"schema_version"`
	Name          string `json:"name"`
	Transport     string `json:"transport"`
	ChannelLabel  string `json:"channel_label"`
	WebhookEnv    string `json:"webhook_env"`
	StateRoot     string `json:"state_root"`
}
type RouteView struct {
	Name                string `json:"name"`
	ChannelLabel        string `json:"channel_label"`
	ChannelVerification string `json:"channel_verification"`
	SHA256              string `json:"route_sha256,omitempty"`
}
type Payload struct {
	Text        string `json:"text"`
	Mrkdwn      bool   `json:"mrkdwn"`
	Parse       string `json:"parse"`
	UnfurlLinks bool   `json:"unfurl_links"`
	UnfurlMedia bool   `json:"unfurl_media"`
}
type Preview struct {
	Kind          string    `json:"kind"`
	SchemaVersion int       `json:"schema_version"`
	ID            string    `json:"notification_id"`
	RequestSHA    string    `json:"request_sha256"`
	Confirmation  *string   `json:"confirmation"`
	Operation     string    `json:"operation"`
	Allowed       bool      `json:"allowed"`
	Reasons       []Reason  `json:"reasons"`
	Route         RouteView `json:"route"`
	StateRoot     string    `json:"state_root"`
	ExistingState *string   `json:"existing_state"`
	NextAttempt   int       `json:"next_attempt"`
	Payload       Payload   `json:"payload"`
	PayloadSHA    string    `json:"payload_sha256"`
	Effects       []string  `json:"effects"`
	SourceKind    string    `json:"source_kind"`
}
type Attempt struct {
	Number         int     `json:"number"`
	State          string  `json:"state"`
	ReservedAt     string  `json:"reserved_at"`
	CompletedAt    *string `json:"completed_at"`
	HTTPStatus     *int    `json:"http_status"`
	ResponseCode   *string `json:"response_code"`
	RetryNotBefore *string `json:"retry_not_before"`
	Confirmation   string  `json:"confirmation"`
	Dispatch       string  `json:"dispatch"`
}
type Result struct {
	Kind           string    `json:"kind"`
	SchemaVersion  int       `json:"schema_version"`
	ID             string    `json:"notification_id"`
	RequestSHA     string    `json:"request_sha256"`
	Source         Source    `json:"source_binding"`
	Route          RouteView `json:"route"`
	Gate           Gate      `json:"gate"`
	Knowledge      string    `json:"knowledge"`
	State          string    `json:"state"`
	Freshness      string    `json:"freshness"`
	PayloadSHA     string    `json:"payload_sha256"`
	Attempts       []Attempt `json:"attempts"`
	RetryNotBefore *string   `json:"retry_not_before"`
	Reasons        []Reason  `json:"reasons"`
	NextAction     string    `json:"next_action"`
	Persistence    string    `json:"persistence"`
}
type ErrorResult struct {
	Kind          string   `json:"kind"`
	SchemaVersion int      `json:"schema_version"`
	ID            *string  `json:"notification_id"`
	State         string   `json:"state"`
	Reasons       []Reason `json:"reasons"`
}

func ErrorOutput(err error) ErrorResult {
	e, ok := err.(*Error)
	if !ok {
		e = &Error{Exit: 1, Code: "local_io", Message: "Local notification storage could not be read or written."}
	}
	return ErrorResult{"ply.workflow.notification-error", 1, nil, "not_attempted", []Reason{{e.Code, e.Message}}}
}

type Input struct {
	Operation, File, Route, ID, Confirm string
	Apply                               bool
}

// Dependencies are an internal composition boundary. Production uses SystemDependencies;
// test executables can replace external transport, time and durable-write fault boundaries.
type Dependencies struct {
	Now       func() time.Time
	LookupEnv func(string) (string, bool)
	Transport func(context.Context, string, []byte) Observation
	Fault     func(string) error
}
type Observation struct {
	Status     int
	Body       []byte
	RetryAfter []string
	Dispatch   string
	Failed     bool
	Oversize   bool
}

func stamp(t time.Time) string { return t.UTC().Format(time.RFC3339) }
