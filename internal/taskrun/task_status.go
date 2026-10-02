package taskrun

import (
	"bytes"
	"path/filepath"
	"strings"
	"time"
)

func parseTaskStatus(raw []byte) (TaskStatus, error) {
	var s TaskStatus
	if e := decode(raw, 64<<10, &s); e != nil {
		return s, e
	}
	if e := checkEnvelope(s.Envelope, "task-status"); e != nil {
		return s, e
	}
	if s.Source != "human_observation" || !plain(s.ActorClaim, 1, 256) || !digestPattern.MatchString(s.BasisEventSHA256) {
		return s, invalid("invalid task observation source or basis")
	}
	if s.State != "inactive" && s.State != "active" && s.State != "unknown" {
		return s, invalid("invalid task execution state")
	}
	if s.TaskLabel != nil && !plain(*s.TaskLabel, 1, 256) || s.NativeSessionID != nil && !plain(*s.NativeSessionID, 1, 256) || s.MatchingTasks != nil && *s.MatchingTasks < 0 {
		return s, invalid("invalid task observation facts")
	}
	tm, e := time.Parse(time.RFC3339Nano, s.ObservedAtUTC)
	if e != nil || tm.Location() != time.UTC {
		return s, invalid("task observation must have a UTC timestamp")
	}
	canonical, _ := Canonical(s)
	if !bytes.Equal(raw, canonical) {
		return s, invalid("task status must use exact canonical JSON bytes")
	}
	return s, nil
}

func validateTaskStatus(s TaskStatus, j journal) error {
	if e := claimBinding(s.RunID, s.RequestSHA256, s.SessionID, j); e != nil {
		return e
	}
	if j.Binding == nil || j.Result.Launch.Attempts != 1 {
		return conflict("task observation requires a bound launch")
	}
	if j.Result.LastEventSHA256 == nil || s.BasisEventSHA256 != *j.Result.LastEventSHA256 {
		return conflict("task observation basis differs from the current journal tip; obtain a new preview")
	}
	plan := j.Binding.Preparation.Plan
	if s.WorktreePath != plan.WorktreePath {
		return conflict("task observation worktree differs")
	}
	if s.VisibleGroupPath != nil && *s.VisibleGroupPath != plan.WorktreePath && *s.VisibleGroupPath != filepath.Dir(plan.Target.GitCommonDir) {
		return conflict("task observation group differs")
	}
	if s.State == "inactive" && (s.VisibleGroupPath == nil || s.TaskLabel == nil || s.MatchingTasks == nil || *s.MatchingTasks != 1) {
		return invalid("inactive requires one unambiguous observed task and its matching group")
	}
	if j.Claim != nil && j.Claim.RuntimeClaim.NativeSessionID != nil && s.NativeSessionID != nil && *j.Claim.RuntimeClaim.NativeSessionID != *s.NativeSessionID {
		return conflict("task observation native session differs from recipient claim")
	}
	observed, _ := time.Parse(time.RFC3339Nano, s.ObservedAtUTC)
	if observed.Before(j.LatestActivityAt) || observed.Before(j.LatestStatusAt) {
		return conflict("task observation predates process, report or previous status evidence")
	}
	return nil
}

func executionFromStatus(s TaskStatus, h string) TaskExecution {
	return TaskExecution{State: s.State, Source: "human_observation", StatusSHA256: &h, ObservedAtUTC: &s.ObservedAtUTC}
}
func taskStatusPath(r Request, h string) string {
	return filepath.Join(runPaths(r).RunRoot, "task-status", strings.TrimPrefix(h, "sha256:")+".json")
}

// proposedStatus never writes. A previously journaled digest is an exact retry;
// all other observations must bind the current tip, even after a file-only crash.
func proposedStatus(d Dependencies, j journal) (*TaskStatus, error) {
	if j.Request.SchemaVersion == 2 && d.TaskStatusPath != "" {
		return nil, invalid("exec rejects --task-status; provider evidence is runner-owned")
	}
	if d.TaskStatusPath == "" {
		return nil, nil
	}
	raw, e := readFile(d.TaskStatusPath, 64<<10, true)
	if e != nil {
		return nil, e
	}
	s, e := parseTaskStatus(raw)
	if e != nil {
		return nil, e
	}
	h := hash(raw)
	for _, ev := range j.Events {
		if ev.Type == "task_status_observed" {
			var p struct {
				StatusSHA256 string `json:"status_sha256"`
			}
			if decode(ev.Payload, 64<<10, &p) == nil && p.StatusSHA256 == h {
				return &s, nil
			}
		}
	}
	if e = validateTaskStatus(s, j); e != nil {
		return nil, e
	}
	observed, _ := time.Parse(time.RFC3339Nano, s.ObservedAtUTC)
	if observed.After(d.Now()) {
		return nil, conflict("task observation is in the future")
	}
	return &s, nil
}
func statusAlreadyRecorded(j journal, h string) bool {
	for _, ev := range j.Events {
		if ev.Type == "task_status_observed" {
			var p struct {
				StatusSHA256 string `json:"status_sha256"`
			}
			if decode(ev.Payload, 64<<10, &p) == nil && p.StatusSHA256 == h {
				return true
			}
		}
	}
	return false
}
func publishTaskStatus(d Dependencies, j journal, confirmedDigest *string) error {
	s, e := proposedStatus(d, j)
	if e != nil || s == nil {
		return e
	}
	h := digest(s)
	if confirmedDigest == nil || h != *confirmedDigest {
		return conflict("task status bytes changed after the confirmed preview")
	}
	if statusAlreadyRecorded(j, h) {
		return nil
	}
	if e = d.writeValue(taskStatusPath(j.Request, h), s); e != nil {
		return e
	}
	if e = d.fault("after_task_status_file"); e != nil {
		return e
	}
	return appendEvent(d, j.Request, "task_status_observed", map[string]any{"status_sha256": h})
}
func taskInactive(j journal) bool {
	return j.Result.TaskExecution.State == "inactive" && j.LatestStatusSequence > j.LatestReportSequence
}
func reservationReleased(j journal) bool {
	return quiescent(j.Result.Process) && (j.Result.TaskExecution.State == "not_started" || taskInactive(j))
}
