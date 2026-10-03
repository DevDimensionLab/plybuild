package taskrun

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type workspaceObservation = workspace.PlanWorktreeObservation

func storeRoot(root string) string { return filepath.Join(root, ".ply", "task-runs", "v1") }
func privateDir(path string) error {
	root, e := physicalRoot(path, true)
	if e != nil {
		return e
	}
	defer root.Close()
	info, e := root.Stat(".")
	if e != nil {
		return e
	}
	if !info.IsDir() || info.Mode().Perm()&0077 != 0 {
		return integrity("store directory must be private")
	}
	return nil
}
func syncDir(path string) error {
	f, e := os.Open(path)
	if e != nil {
		return e
	}
	defer f.Close()
	return f.Sync()
}

// writeOnce publishes only complete, synced bytes. Link is atomic and cannot
// clobber a concurrent writer; the temporary is owned exclusively by this call.
func writeOnce(path string, b []byte) error { return publishFile(path, b, false) }
func publishFile(path string, b []byte, replace bool) error {
	if e := physical(path, true); e != nil {
		return e
	}
	if !replace {
		if old, e := readFile(path, 8<<20, true); e == nil {
			if bytes.Equal(old, b) {
				return nil
			}
			return conflict("immutable bytes differ: " + filepath.Base(path))
		} else if !os.IsNotExist(e) {
			return e
		}
	}
	parent := filepath.Dir(path)
	if e := privateDir(parent); e != nil {
		return e
	}
	before, e := os.Lstat(parent)
	if e != nil {
		return e
	}
	root, e := physicalRoot(parent, false)
	if e != nil {
		return e
	}
	defer root.Close()
	actual, e := root.Stat(".")
	if e != nil || !os.SameFile(before, actual) {
		return conflict("publication directory changed")
	}
	nonce := make([]byte, 16)
	if _, e = rand.Read(nonce); e != nil {
		return e
	}
	name := ".publish-" + hex.EncodeToString(nonce)
	f, e := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		return e
	}
	defer root.Remove(name)
	_, e = f.Write(b)
	if e == nil {
		e = f.Sync()
	}
	c := f.Close()
	if e != nil {
		return e
	}
	if c != nil {
		return c
	}
	destination := filepath.Base(path)
	if replace {
		e = root.Rename(name, destination)
	} else {
		e = root.Link(name, destination)
	}
	if e != nil {
		if !replace && os.IsExist(e) {
			old, err := root.ReadFile(destination)
			if err == nil && bytes.Equal(old, b) {
				return nil
			}
		}
		return e
	}
	dir, e := root.Open(".")
	if e != nil {
		return e
	}
	defer dir.Close()
	return dir.Sync()
}

func writeValue(path string, v any) error {
	b, e := Canonical(v)
	if e != nil {
		return e
	}
	return writeOnce(path, b)
}
func readValue(path string, max int, v any) error {
	b, e := readFile(path, max, true)
	if e != nil {
		return e
	}
	if e = decode(b, max, v); e != nil {
		return integrity(e.Error())
	}
	return nil
}
func withStore(root string, fn func() error) error {
	p := storeRoot(root)
	if e := privateDir(p); e != nil {
		return e
	}
	f, e := lockFile(filepath.Join(p, "store.lock"))
	if e != nil {
		return e
	}
	defer unlockFile(f)
	return fn()
}

type journal struct {
	LaunchArgvSHA256     string
	Request              Request
	Binding              *Binding
	Events               []Event
	Hashes               []string
	Result               Result
	Claim                *Acceptance
	LatestActivityAt     time.Time
	LatestStatusAt       time.Time
	LatestReportSequence int
	LatestStatusSequence int
	LastTaskStatus       *TaskStatus
}

func initial(r Request) Result {
	envelope := env("result")
	if r.SchemaVersion == 2 {
		envelope.SchemaVersion = 3
	}
	return Result{Envelope: envelope, RunID: RunID(r.RequestKey), RequestSHA256: digest(r), RuntimeFacts: RuntimeFacts{RequestedModel: r.Runtime.Model, ModelState: "missing"}, TaskExecution: TaskExecution{State: "unknown", Source: "none"}, Launch: Launch{State: "reserved"}, Acceptance: AcceptanceState{State: "missing"}, Process: Process{State: "not_started"}, Delivery: Delivery{State: "missing"}, Collection: Collection{State: "pending"}, Reasons: []Reason{}, NextAction: "Start plan result control before separate human QA or integration."}
}
func readJournal(root, id string) (journal, error) {
	var j journal
	if !runPattern.MatchString(id) {
		return j, invalid("invalid run ID")
	}
	run := filepath.Join(storeRoot(root), "runs", id)
	b, e := readFile(filepath.Join(run, "request.json"), 1<<20, true)
	if e != nil {
		return j, e
	}
	r, e := parseRequest(b)
	if e != nil {
		return j, integrity(e.Error())
	}
	if RunID(r.RequestKey) != id || r.WorkspaceRoot != root {
		return j, integrity("request identity differs")
	}
	j.Request = r
	j.Result = initial(r)
	var binding Binding
	e = readValue(filepath.Join(run, "binding.json"), 2<<20, &binding)
	if e == nil {
		if binding.Envelope != env("binding") || binding.RunID != id || binding.RequestSHA256 != digest(r) || binding.SessionID != "ply:"+id || binding.RunRoot != run || binding.TempRoot != runPaths(r).TempRoot || binding.ReportPath != runPaths(r).ReportPath {
			return j, integrity("binding differs")
		}
		j.Binding = &binding
	} else if !os.IsNotExist(e) {
		return j, e
	}
	entries, e := os.ReadDir(filepath.Join(run, "events"))
	if e != nil && !os.IsNotExist(e) {
		return j, e
	}
	var previous *string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".publish-") {
			continue
		}
		want := fmt.Sprintf("%06d.json", len(j.Events)+1)
		if entry.Name() != want {
			return j, integrity("event sequence gap or unexpected file")
		}
		var event Event
		b, e := readFile(filepath.Join(run, "events", entry.Name()), 64<<10, true)
		if e != nil {
			return j, e
		}
		if e = decode(b, 64<<10, &event); e != nil {
			return j, integrity(e.Error())
		}
		if event.Envelope != env("event") || event.RunID != id || event.RequestSHA256 != digest(r) || event.Sequence != len(j.Events)+1 || !equal(event.PreviousSHA256, previous) {
			return j, integrity("event chain differs")
		}
		tm, e := time.Parse(time.RFC3339Nano, event.RecordedAtUTC)
		if e != nil || tm.Location() != time.UTC {
			return j, integrity("invalid event UTC")
		}
		if e = fold(&j, event); e != nil {
			return j, e
		}
		h := hash(b)
		previous = &h
		j.Events = append(j.Events, event)
		j.Hashes = append(j.Hashes, h)
		j.Result.LastEventSHA256 = &h
	}
	if len(j.Events) == 0 || j.Result.Launch.State == "reserved" || j.Result.Launch.State == "intent" {
		j.Result.Launch.State = "unknown"
		j.Result.Process.State = "unknown"
		j.Result.Collection.State = "unknown"
		j.Result.Reasons = append(j.Result.Reasons, Reason{"task_run_setup_unknown", "Preserved reservation has no proven completed launch sequence; it must never be relaunched."})
	}
	if j.Result.Process.State == "running" {
		// A durable start is not a liveness observation. In particular, parent
		// loss cannot be inferred away from a lock file or event age.
		j.Result.Process.State = "unknown"
		j.Result.Process.Quiescence = nil
		j.Result.Reasons = append(j.Result.Reasons, Reason{"task_run_process_unobserved", "A process start was preserved, but no bound Wait/group-end observation is available yet."})
	}
	// A cache must name a known chain prefix and exactly match its derived state.
	var cache Result
	rawCache, cacheErr := readFile(filepath.Join(run, "result.json"), 4<<20, true)
	e = cacheErr
	if e == nil {
		cache, e = decodeResultCache(rawCache)
	}
	if e == nil {
		if cache.LastEventSHA256 == nil {
			return j, integrity("cache has no event binding")
		}
		idx := -1
		for i, h := range j.Hashes {
			if h == *cache.LastEventSHA256 {
				idx = i
			}
		}
		if idx < 0 {
			return j, integrity("cache refers to an unknown chain")
		}
		prefix := journal{Request: r, Binding: j.Binding, Result: initial(r)}
		for i := 0; i <= idx; i++ {
			if e = fold(&prefix, j.Events[i]); e != nil {
				return j, e
			}
			prefix.Result.LastEventSHA256 = &j.Hashes[i]
		}
		expected := prefix.Result
		if cache.SchemaVersion == 1 {
			// The old producer collapsed negative states in the event/cache. Validate
			// those exact historical bytes, while current readback follows the claim.
			for _, ev := range j.Events[:idx+1] {
				if ev.Type == "acceptance_received" {
					var p struct {
						State string `json:"state"`
					}
					if err := json.Unmarshal(ev.Payload, &p); err != nil {
						return j, integrity(err.Error())
					}
					expected.Acceptance.State = p.State
				}
			}
		}
		if !equal(resultCacheValue(cache, cache.SchemaVersion), resultCacheValue(expected, cache.SchemaVersion)) {
			return j, integrity("cache contradicts its event prefix")
		}
	} else if !os.IsNotExist(e) {
		return j, e
	}
	return j, nil
}
func equal(a, b any) bool {
	x, e := Canonical(a)
	if e != nil {
		return false
	}
	y, e := Canonical(b)
	return e == nil && bytes.Equal(x, y)
}
func fold(j *journal, e Event) error {
	r := &j.Result
	if e.Sequence < 1 {
		return integrity("invalid event sequence")
	}
	switch e.Type {
	case "reserved":
		var p struct {
			Target     workspace.PlanWorktreeObservation `json:"target"`
			RequestKey string                            `json:"request_key"`
		}
		if decode(e.Payload, 64<<10, &p) != nil || e.Sequence != 1 || p.RequestKey != j.Request.RequestKey {
			return integrity("invalid reservation")
		}
	case "handoff_bound":
		var p struct {
			BindingSHA256 string `json:"binding_sha256"`
		}
		if decode(e.Payload, 64<<10, &p) != nil || j.Binding == nil || p.BindingSHA256 != digest(j.Binding) || r.Handoff != nil {
			return integrity("invalid handoff binding")
		}
		r.Handoff = &j.Binding.Handoff
	case "launch_intent":
		var p struct {
			ExecutableSHA256      string `json:"executable_sha256"`
			ArgvSHA256            string `json:"argv_sha256"`
			CWD                   string `json:"cwd"`
			SessionID             string `json:"session_id"`
			EffectivePolicySHA256 string `json:"effective_policy_sha256"`
		}
		if decode(e.Payload, 64<<10, &p) != nil || r.Handoff == nil || r.Launch.Attempts != 0 || p.ExecutableSHA256 != j.Request.Runtime.Executable.SHA256 || p.SessionID != "ply:"+r.RunID || p.EffectivePolicySHA256 != j.Request.Runtime.PermissionBinding.EffectivePolicySHA256 || !digestPattern.MatchString(p.ArgvSHA256) {
			return integrity("invalid launch intent")
		}
		j.LaunchArgvSHA256 = p.ArgvSHA256
		r.Launch = Launch{"intent", 1}
		r.Process.State = "unknown"
	case "launch_failed":
		var p struct {
			Code   string `json:"code"`
			Detail string `json:"detail"`
		}
		if decode(e.Payload, 64<<10, &p) != nil || r.Launch.State != "intent" {
			return integrity("invalid launch failure")
		}
		r.Launch.State = "failed"
		r.Process = Process{State: "not_started", Quiescence: ptr(true)}
		r.TaskExecution = TaskExecution{State: "not_started", Source: "spawn_failure"}
		r.Reasons = append(r.Reasons, Reason{p.Code, p.Detail})
	case "process_started", "process_exited":
		var p struct {
			Process Process `json:"process"`
		}
		if decode(e.Payload, 64<<10, &p) != nil || r.Launch.Attempts != 1 {
			return integrity("invalid process event")
		}
		x := p.Process
		if err := validateProcess(x); err != nil {
			return err
		}
		if e.Type == "process_started" {
			if r.Launch.State != "intent" || x.State != "running" || x.PID == nil || *x.PID <= 0 || x.ProcessGroup == nil || *x.ProcessGroup != *x.PID || x.StartIdentity == nil || !plain(*x.StartIdentity, 1, 1000) {
				return integrity("invalid process start identity")
			}
			r.Launch.State = "started"
		} else {
			if r.Process.State == "exited" || r.Launch.State == "failed" {
				return integrity("duplicate or contradictory process exit")
			}
			if x.State != "exited" && x.State != "unknown" {
				return integrity("invalid process exit")
			}
			if r.Process.PID != nil && (!equal(r.Process.PID, x.PID) || !equal(r.Process.StartIdentity, x.StartIdentity) || !equal(r.Process.ProcessGroup, x.ProcessGroup)) {
				return integrity("process identity changed")
			}
		}
		j.LatestActivityAt, _ = time.Parse(time.RFC3339Nano, e.RecordedAtUTC)
		r.Process = x
	case "acceptance_received":
		var p struct {
			ClaimSHA256   string  `json:"claim_sha256"`
			ReceiptSHA256 *string `json:"receipt_sha256"`
			State         string  `json:"state"`
		}
		if decode(e.Payload, 64<<10, &p) != nil || r.Acceptance.State != "missing" || !digestPattern.MatchString(p.ClaimSHA256) {
			return integrity("invalid acceptance event")
		}
		if p.State != "started" && p.State != "rejected" && p.State != "conflict" && p.State != "unknown" {
			return integrity("unknown acceptance state")
		}
		if p.ReceiptSHA256 != nil && !digestPattern.MatchString(*p.ReceiptSHA256) || p.State == "started" && p.ReceiptSHA256 == nil || r.Launch.Attempts != 1 {
			return integrity("invalid acceptance facts")
		}
		raw, err := readFile(filepath.Join(runPaths(j.Request).RunRoot, "claims", strings.TrimPrefix(p.ClaimSHA256, "sha256:")+".json"), 256<<10, true)
		if err != nil || hash(raw) != p.ClaimSHA256 {
			return integrity("accepted claim bytes differ")
		}
		claim, err := parseAcceptance(raw)
		if err != nil || claimBinding(claim.RunID, claim.RequestSHA256, claim.SessionID, *j) != nil {
			return integrity("invalid accepted claim")
		}
		// Older transports collapsed all negative event states to rejected. Preserve
		// those bytes, but derive the actual claim state for current readback.
		if claim.Acceptance != p.State && !(claim.SchemaVersion == 1 && p.State == "rejected" && claim.Acceptance != "started") {
			return integrity("acceptance state differs from claim")
		}
		if claim.Acceptance == "started" {
			if err = positiveClaim(claim, j.Request); err != nil {
				return integrity(err.Error())
			}
		}
		if j.LastTaskStatus != nil && j.LastTaskStatus.NativeSessionID != nil && claim.RuntimeClaim.NativeSessionID != nil && *j.LastTaskStatus.NativeSessionID != *claim.RuntimeClaim.NativeSessionID {
			return integrity("claim and task status native sessions differ")
		}
		j.Claim = &claim
		model := reportedModel(claim)
		state := "unknown"
		if model != nil {
			state = "reported"
			if *model != j.Request.Runtime.Model {
				state = "mismatch"
			}
		}
		r.RuntimeFacts = RuntimeFacts{RequestedModel: j.Request.Runtime.Model, ReportedModel: model, ModelState: state, Source: ptr("recipient_claim"), ClaimSHA256: &p.ClaimSHA256}
		r.Acceptance = AcceptanceState{claim.Acceptance, p.ReceiptSHA256}
	case "report_received":
		var p struct {
			ReportSHA256   string  `json:"report_sha256"`
			TerminalSHA256 *string `json:"terminal_sha256"`
			State          string  `json:"state"`
		}
		if decode(e.Payload, 64<<10, &p) != nil || r.Delivery.State != "missing" {
			return integrity("invalid report event")
		}
		if !digestPattern.MatchString(p.ReportSHA256) || p.State != "received" && p.State != "not_representable" && p.State != "conflict" || p.TerminalSHA256 != nil && !digestPattern.MatchString(*p.TerminalSHA256) || p.State == "received" && p.TerminalSHA256 == nil {
			return integrity("invalid report facts")
		}
		path := filepath.Join(runPaths(j.Request).RunRoot, "reports", strings.TrimPrefix(p.ReportSHA256, "sha256:")+".json")
		b, err := readFile(path, 4<<20, true)
		if err != nil || hash(b) != p.ReportSHA256 {
			return integrity("report bytes differ")
		}
		report, err := validateReport(b)
		if err != nil {
			return integrity(err.Error())
		}
		j.LatestActivityAt, _ = time.Parse(time.RFC3339Nano, e.RecordedAtUTC)
		j.LatestReportSequence = e.Sequence
		r.Delivery = Delivery{p.State, &p.ReportSHA256, p.TerminalSHA256, &report.Outcome}
		r.Budget.Reported = report.BudgetUsage
		r.Budget.WithinAgreement = budgetValid(report.BudgetUsage)
	case "provider_completion":
		return foldProviderCompletion(j, e)
	case "task_status_observed":
		if j.Request.SchemaVersion == 2 {
			return integrity("exec does not accept manual task status")
		}
		var p struct {
			StatusSHA256 string `json:"status_sha256"`
		}
		if decode(e.Payload, 64<<10, &p) != nil || !digestPattern.MatchString(p.StatusSHA256) {
			return integrity("invalid task status event")
		}
		raw, err := readFile(taskStatusPath(j.Request, p.StatusSHA256), 64<<10, true)
		if err != nil || hash(raw) != p.StatusSHA256 {
			return integrity("task status bytes differ")
		}
		status, err := parseTaskStatus(raw)
		if err != nil {
			return integrity(err.Error())
		}
		if err = validateTaskStatus(status, *j); err != nil {
			return integrity(err.Error())
		}
		observed, _ := time.Parse(time.RFC3339Nano, status.ObservedAtUTC)
		recorded, _ := time.Parse(time.RFC3339Nano, e.RecordedAtUTC)
		if recorded.Before(observed) {
			return integrity("task status event predates observation")
		}
		j.LatestStatusAt = observed
		j.LatestStatusSequence = e.Sequence
		j.LastTaskStatus = &status
		r.TaskExecution = executionFromStatus(status, p.StatusSHA256)
	case "collection":
		var p struct {
			Collection     Collection            `json:"collection"`
			ObservedTarget *workspaceObservation `json:"observed_target"`
			Reasons        []Reason              `json:"reasons"`
		}
		if decode(e.Payload, 64<<10, &p) != nil {
			return integrity("invalid collection event")
		}
		switch p.Collection.State {
		case "pending", "qualified", "blocked", "unknown", "conflict":
		default:
			return integrity("invalid collection state")
		}
		if p.Collection.State == "qualified" && (p.Collection.TaskResultID == nil || p.Collection.TaskResultDraftSHA256 == nil || !digestPattern.MatchString(*p.Collection.TaskResultDraftSHA256) || !quiescent(r.Process) || r.Process.State != "exited" || r.Process.ExitCode == nil || *r.Process.ExitCode != 0 || r.Process.Signal != nil || r.Acceptance.State != "started" || r.Delivery.State != "received" || r.Delivery.ReportedOutcome == nil || *r.Delivery.ReportedOutcome != "complete" || r.Budget.WithinAgreement == nil || !*r.Budget.WithinAgreement) {
			return integrity("collection contradicts qualifying facts")
		}
		if p.Collection.State == "qualified" && j.Request.SchemaVersion == 2 && !taskInactive(*j) {
			return integrity("exec collection lacks provider completion")
		}
		if p.Collection.State == "qualified" && r.TaskExecution.StatusSHA256 == nil && j.Request.SchemaVersion == 1 {
			r.TaskExecution.HistoricalQualification = true
		}
		r.Collection = p.Collection
		r.ObservedTarget = p.ObservedTarget
		r.Reasons = p.Reasons
	default:
		return integrity("unknown event type")
	}
	return nil
}
func replaceValue(path string, v any) error {
	b, e := Canonical(v)
	if e != nil {
		return e
	}
	return publishFile(path, b, true)
}

func appendEvent(d Dependencies, r Request, typ string, payload any) error {
	j, e := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
	if e != nil {
		return e
	}
	events, hashes := j.Events, j.Hashes
	j = journal{Request: r, Binding: j.Binding, Result: initial(r), Events: events, Hashes: hashes}
	for i, x := range j.Events {
		if e = fold(&j, x); e != nil {
			return e
		}
		j.Result.LastEventSHA256 = &j.Hashes[i]
	}
	raw, e := Canonical(payload)
	if e != nil {
		return e
	}
	event := Event{env("event"), j.Result.RunID, digest(r), len(j.Events) + 1, j.Result.LastEventSHA256, typ, d.Now().UTC().Format(time.RFC3339Nano), raw}
	if e = fold(&j, event); e != nil {
		return e
	}
	b, e := Canonical(event)
	if e != nil {
		return e
	}
	if len(b) > 64<<10 {
		return invalid("event exceeds 64 KiB")
	}
	if e = d.writeOnce(filepath.Join(runPaths(r).RunRoot, "events", fmt.Sprintf("%06d.json", event.Sequence)), b); e != nil {
		return e
	}
	j.Result.LastEventSHA256 = ptr(hash(b))
	if e = d.fault("after_event_" + typ); e != nil {
		return e
	}
	return d.replaceValue(filepath.Join(runPaths(r).RunRoot, "result.json"), j.Result)
}

func validateProcess(p Process) error {
	switch p.State {
	case "not_started", "running", "exited", "unknown":
	default:
		return integrity("invalid process state")
	}
	if p.PID != nil && *p.PID <= 0 || p.ProcessGroup != nil && *p.ProcessGroup <= 0 || p.StartIdentity != nil && !plain(*p.StartIdentity, 1, 1000) || p.ExitCode != nil && (*p.ExitCode < 0 || *p.ExitCode > 255) || p.Signal != nil && !plain(*p.Signal, 1, 256) {
		return integrity("invalid process facts")
	}
	if p.State == "running" && (p.PID == nil || p.ProcessGroup == nil || p.StartIdentity == nil || p.ExitCode != nil || p.Signal != nil) {
		return integrity("running process lacks identity")
	}
	if p.State == "exited" && (p.PID == nil || p.ProcessGroup == nil || p.StartIdentity == nil || p.ExitCode == nil && p.Signal == nil) {
		return integrity("exit lacks Wait identity or status")
	}
	return nil
}
