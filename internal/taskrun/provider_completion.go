package taskrun

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

const providerStreamLimit = 10 << 20
const providerLineLimit = 1 << 20

type ExecLaunch struct {
	Executable          Executable `json:"executable"`
	Argv                []string   `json:"argv"`
	CWD                 string     `json:"cwd"`
	RequestSHA256       string     `json:"request_sha256"`
	TimeoutMilliseconds int64      `json:"timeout_milliseconds"`
}
type ProviderCompletion struct {
	Kind            string          `json:"kind"`
	RequestSHA256   string          `json:"request_sha256"`
	LaunchSHA256    string          `json:"launch_sha256"`
	Stdout          FileBinding     `json:"stdout"`
	Stderr          FileBinding     `json:"stderr"`
	StdoutTruncated bool            `json:"stdout_truncated"`
	StderrTruncated bool            `json:"stderr_truncated"`
	ThreadID        *string         `json:"thread_id"`
	TurnsStarted    int             `json:"turns_started"`
	TurnsCompleted  int             `json:"turns_completed"`
	StreamComplete  bool            `json:"stream_complete"`
	SequenceValid   bool            `json:"sequence_valid"`
	TimedOut        bool            `json:"timed_out"`
	Process         Process         `json:"process"`
	TokenUsage      json.RawMessage `json:"token_usage"`
	Reason          string          `json:"reason"`
}
type providerSequence struct {
	ThreadID           *string
	Started, Completed int
	Valid              bool
	Reason             string
	Usage              json.RawMessage
}

func parseProviderEvents(raw []byte) providerSequence {
	return parseProviderSequence(raw, true)
}

func parseProviderSequence(raw []byte, completedFailures bool) providerSequence {
	s := providerSequence{Valid: true, Usage: json.RawMessage("null")}
	fail := func(reason string) providerSequence { s.Valid = false; s.Reason = reason; return s }
	if len(raw) == 0 || raw[len(raw)-1] != '\n' {
		return fail("incomplete JSONL stream")
	}
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	scanner.Buffer(make([]byte, 4096), providerLineLimit+1)
	pending := map[string]string{}
	finished := map[string]bool{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) > providerLineLimit {
			return fail("JSONL line exceeds 1 MiB")
		}
		if _, e := canonicaljson.DecodeStrict(line); e != nil {
			return fail("unparsable JSONL event")
		}
		var ev struct {
			Type     string `json:"type"`
			ThreadID string `json:"thread_id"`
			Item     struct {
				ID     string `json:"id"`
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"item"`
			Usage json.RawMessage `json:"usage"`
		}
		if e := json.Unmarshal(line, &ev); e != nil {
			return fail("invalid event object")
		}
		if s.Completed > 0 {
			return fail("event after terminal turn")
		}
		switch ev.Type {
		case "thread.started":
			if s.ThreadID != nil || s.Started != 0 || !plain(ev.ThreadID, 1, 256) {
				return fail("inconsistent native thread")
			}
			s.ThreadID = ptr(ev.ThreadID)
		case "turn.started":
			if s.ThreadID == nil || s.Started != 0 {
				return fail("expected exactly one turn")
			}
			s.Started++
		case "item.started":
			if s.Started != 1 || ev.Item.ID == "" || ev.Item.Type == "" || pending[ev.Item.ID] != "" || finished[ev.Item.ID] {
				return fail("invalid item start")
			}
			pending[ev.Item.ID] = ev.Item.Type
		case "item.updated":
			if s.Started != 1 || pending[ev.Item.ID] != ev.Item.Type || ev.Item.Type == "" {
				return fail("update without matching start")
			}
		case "item.completed":
			if s.Started != 1 || ev.Item.ID == "" || finished[ev.Item.ID] {
				return fail("invalid item completion")
			}
			typ := ev.Item.Type
			if prior := pending[ev.Item.ID]; prior != "" {
				if prior != typ {
					return fail("item changed type")
				}
				delete(pending, ev.Item.ID)
			} else if typ != "agent_message" && typ != "reasoning" && typ != "todo_list" {
				return fail("tool completion without start")
			}
			// A completed failed command is still closed. Its failure stays in
			// the bound stream; verifier/report checks decide delivery quality.
			if !completedFailures && (ev.Item.Status == "in_progress" || ev.Item.Status == "failed") {
				return fail("unfinished or failed tool item")
			}
			if ev.Item.Status == "in_progress" {
				return fail("unfinished tool item")
			}
			finished[ev.Item.ID] = true
		case "turn.completed":
			if s.Started != 1 || len(pending) != 0 {
				return fail("turn completed with unfinished items")
			}
			s.Completed++
			if len(ev.Usage) > 0 {
				s.Usage = ev.Usage
			}
		default:
			return fail("unknown, failed or detached provider event: " + ev.Type)
		}
		if ev.Type != "thread.started" && ev.ThreadID != "" && (s.ThreadID == nil || ev.ThreadID != *s.ThreadID) {
			return fail("native thread changed")
		}
	}
	if scanner.Err() != nil {
		return fail("JSONL line exceeds limit or stream read failed")
	}
	if s.ThreadID == nil || s.Started != 1 || s.Completed != 1 || len(pending) != 0 {
		return fail("missing terminal turn")
	}
	return s
}
func providerInactive(c ProviderCompletion) bool {
	return c.Kind == "ProviderCompletion@1" && c.StreamComplete && c.SequenceValid && !c.StdoutTruncated && !c.StderrTruncated && !c.TimedOut && c.ThreadID != nil && c.TurnsStarted == 1 && c.TurnsCompleted == 1 && c.Process.State == "exited" && c.Process.ExitCode != nil && *c.Process.ExitCode == 0 && c.Process.Signal == nil && quiescent(c.Process)
}
func publishProviderCompletion(d Dependencies, r Request, s LaunchSpec, p Process) error {
	c := *s.Completion
	c.RequestSHA256 = digest(r)
	c.Process = p
	launch, e := readFile(filepath.Join(s.StreamsRoot, "launch.json"), 64<<10, true)
	if e != nil {
		return e
	}
	c.LaunchSHA256 = hash(launch)
	path := filepath.Join(s.StreamsRoot, "completion.json")
	if e = d.writeValue(path, c); e != nil {
		return e
	}
	return appendEvent(d, r, "provider_completion", FileBinding{path, digest(c)})
}
func foldProviderCompletion(j *journal, e Event) error {
	if j.Request.SchemaVersion != 2 || j.Result.ProviderCompletion != nil || j.Result.Process.State == "running" {
		return integrity("unexpected provider completion")
	}
	var ref FileBinding
	if err := decode(e.Payload, 64<<10, &ref); err != nil {
		return err
	}
	a, err := factoryAuthorization(j.Request)
	if err != nil {
		return err
	}
	root := factorySlotRoot(j.Request, a)
	if ref.Locator != filepath.Join(root, "completion.json") {
		return integrity("completion outside runner-owned slot")
	}
	raw, err := readFile(ref.Locator, 64<<10, true)
	if err != nil || hash(raw) != ref.SHA256 {
		return integrity("provider completion bytes differ")
	}
	var c ProviderCompletion
	if err = decode(raw, 64<<10, &c); err != nil {
		return err
	}
	if c.Kind != "ProviderCompletion@1" || c.RequestSHA256 != digest(j.Request) || !equal(c.Process, j.Result.Process) || c.Stdout.Locator != filepath.Join(root, "stdout.jsonl") || c.Stderr.Locator != filepath.Join(root, "stderr.txt") {
		return integrity("provider completion binding differs")
	}
	launch, err := readFile(filepath.Join(root, "launch.json"), 64<<10, true)
	if err != nil || hash(launch) != c.LaunchSHA256 {
		return integrity("launch bytes differ")
	}
	var l ExecLaunch
	if err = decode(launch, 64<<10, &l); err != nil {
		return err
	}
	if l.RequestSHA256 != c.RequestSHA256 || l.Executable != j.Request.Runtime.Executable || j.Binding == nil || l.CWD != j.Binding.Preparation.Plan.WorktreePath {
		return integrity("launch binding differs")
	}
	if j.LaunchArgvSHA256 != digest(l.Argv) {
		return integrity("completion lacks its launch intent")
	}
	stdout, err := readFile(c.Stdout.Locator, providerStreamLimit, true)
	if err != nil || hash(stdout) != c.Stdout.SHA256 {
		return integrity("provider stdout changed")
	}
	stderr, err := readFile(c.Stderr.Locator, providerStreamLimit, true)
	if err != nil || hash(stderr) != c.Stderr.SHA256 {
		return integrity("provider stderr changed")
	}
	seq := parseProviderEvents(stdout)
	// Old negative projections stopped at the first failed item. Validate those
	// original facts without rewriting their bytes or upgrading them to inactive.
	if !c.SequenceValid && c.Reason == "unfinished or failed tool item" {
		seq = parseProviderSequence(stdout, false)
	}
	if c.SequenceValid != seq.Valid || !equal(c.ThreadID, seq.ThreadID) || c.TurnsStarted != seq.Started || c.TurnsCompleted != seq.Completed || !equal(c.TokenUsage, seq.Usage) {
		return integrity("provider sequence projection differs")
	}
	if j.Claim != nil && j.Claim.RuntimeClaim.NativeSessionID != nil && c.ThreadID != nil && *j.Claim.RuntimeClaim.NativeSessionID != *c.ThreadID {
		return integrity("recipient native session differs from provider")
	}
	j.Result.ProviderCompletion = &ref
	state := "unknown"
	if providerInactive(c) {
		state = "inactive"
	}
	j.Result.TaskExecution = TaskExecution{State: state, Source: "provider_terminal_event", ObservedAtUTC: ptr(e.RecordedAtUTC)}
	j.LatestStatusSequence = e.Sequence
	j.LatestStatusAt, _ = time.Parse(time.RFC3339Nano, e.RecordedAtUTC)
	if state == "unknown" {
		j.Result.Reasons = append(j.Result.Reasons, Reason{"task_run_provider_completion_unknown", fmt.Sprintf("Provider completion is not qualifying: %s", c.Reason)})
	}
	return nil
}
