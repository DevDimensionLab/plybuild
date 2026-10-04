package taskjournal

import (
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRegistrationClockRegressionIsVisible(t *testing.T) {
	f := fixture(t)
	now := time.Date(2026, 10, 3, 10, 0, 0, 0, time.UTC)
	f.s.Now = func() time.Time { return now }
	f.append(t, f.input)
	now = now.Add(-time.Minute)
	in := f.input
	in.PublicationKey = "regressed-clock"
	r := f.append(t, in)
	s := f.show(t)
	found := false
	for _, id := range s.Coverage.TimeConflictEventIDs {
		found = found || id == r.EventID
	}
	if !found || s.Coverage.State != "partial" {
		t.Fatal("registration clock contradicted immutable sequence without warning")
	}
}

func TestOfflineRejectsInventedStepDuration(t *testing.T) {
	f := fixture(t)
	in := f.input
	in.Type = "step_started"
	in.Step = &StepRef{"step", "implementation", nil}
	f.append(t, in)
	s := f.show(t)
	s.Steps[0].State = "closed"
	s.Steps[0].DurationSeconds = ptr(600.0)
	s.SnapshotID = snapshotDigest(s)
	if _, e := ReadSnapshot(f.file(t, s), "task", Options{}); e == nil {
		t.Fatal("offline snapshot invented a closed interval without a finish")
	}
}

func TestNativeSnapshotBindsContentDocuments(t *testing.T) {
	f := fixture(t)
	s := f.show(t)
	found := false
	for _, src := range s.Sources {
		found = found || src.Kind == "task_content_document" && src.SHA256 != nil && src.Status == "valid"
	}
	if !found {
		t.Fatal("snapshot omitted the content bytes read by native validation")
	}
}

func TestMissingNativeRelationRemainsOfflineReadable(t *testing.T) {
	f := fixture(t)
	s := f.show(t)
	var target Event
	for _, e := range s.Events {
		if e.Type != "task_state" {
			target = e
			break
		}
	}
	if target.EventID == "" {
		t.Fatal("fixture content event missing")
	}
	in := f.input
	in.Relations = []Relation{{"responds_to", target.EventID}}
	r := f.append(t, in)
	var path string
	for _, src := range s.Sources {
		if src.SourceID == target.SourceIDs[0] {
			path = src.Locator
		}
	}
	if e := os.Rename(path, path+".diagnostic-missing"); e != nil {
		t.Fatal(e)
	}
	s = f.show(t)
	found := false
	for _, e := range s.Events {
		found = found || e.EventID == r.EventID
	}
	if !found || s.Coverage.State != "partial" {
		t.Fatal("independent contribution disappeared")
	}
	if _, e := ReadSnapshot(f.file(t, s), "task", Options{}); e != nil {
		t.Fatalf("partial snapshot lost reference closure: %v", e)
	}
}

func TestStorageCreationDoesNotRecreateMissingWorkspace(t *testing.T) {
	f := fixture(t)
	f.s.fault = func(stage string) error {
		if stage == "before_directories" {
			return os.Rename(filepath.Join(f.root, ".ply"), filepath.Join(f.root, "saved-ply"))
		}
		return nil
	}
	if _, e := f.s.Append("task", f.file(t, f.input), "event"); e == nil {
		t.Fatal("append proceeded after workspace marker disappeared")
	}
	if _, e := os.Lstat(filepath.Join(f.root, ".ply")); !os.IsNotExist(e) {
		t.Fatal("journal recreated a workspace ancestor")
	}
}

func TestOfflineRejectsCyclicRelationsAndFalseAxis(t *testing.T) {
	f := fixture(t)
	first := f.append(t, f.input)
	second := f.input
	second.PublicationKey = "second"
	last := f.append(t, second)
	for _, mode := range []string{"cycle", "false-native-axis", "mismatched-actor", "invalid-data", "invalid-observation-target"} {
		t.Run(mode, func(t *testing.T) {
			s := f.show(t)
			for i := range s.Events {
				e := &s.Events[i]
				if e.EventID == first.EventID {
					switch mode {
					case "cycle":
						e.Relations = []Relation{{"responds_to", last.EventID}}
					case "false-native-axis":
						s.Current.Axes["human"] = Axis{"pass", []string{e.EventID}, "current"}
					case "mismatched-actor":
						e.Actor = &Actor{"different", "observer", nil}
					case "invalid-data":
						e.Data = raw(map[string]any{"arbitrary_payload": "not a journal input"})
					case "invalid-observation-target":
						s.Observations = append(s.Observations, Observation{ObservationInput: ObservationInput{Common: f.input.Common, Type: "proposal", Targets: Targets{EventIDs: []string{}, StepIDs: []string{}, Interval: &Interval{"bad", "bad", []string{"missing"}, []string{}}}}, EventID: e.EventID})
					}
				}
				if mode == "cycle" && e.EventID == last.EventID {
					e.Relations = []Relation{{"responds_to", first.EventID}}
				}
			}
			s.SnapshotID = snapshotDigest(s)
			if _, e := ReadSnapshot(f.file(t, s), "task", Options{}); e == nil {
				t.Fatal("invalid digest-bound snapshot accepted")
			}
		})
	}
}

func TestCorruptSemanticRecordBlocksNewAppend(t *testing.T) {
	f := fixture(t)
	r := f.append(t, f.input)
	b, e := os.ReadFile(r.RecordLocator)
	if e != nil {
		t.Fatal(e)
	}
	var record Record
	json.Unmarshal(b, &record)
	in := f.input
	in.Type = "step_finished"
	in.Step = &StepRef{"no-start", "implementation", nil}
	record.Input = raw(in)
	record.InputSHA256 = hash(record.Input)
	write(t, r.RecordLocator, raw(record))
	next := f.input
	next.PublicationKey = "next"
	if _, e = f.s.Append("task", f.file(t, next), "event"); e == nil {
		t.Fatal("syntactically valid but semantically damaged log appended")
	}
}

func TestCandidateObjectAndTaskBinding(t *testing.T) {
	f := fixture(t)
	r, e := f.s.Workspace.WorkItems.Snapshot(f.root)
	if e != nil {
		t.Fatal(e)
	}
	base := r.Epics[0].RepoBindings[0].Worktree
	created, e := workspace.CreateTaskWorktree(f.s.Workspace, workspace.TaskWorktreeCreateInput{TaskID: "task", Branch: "task", Path: filepath.Join(f.root, "task-work"), ExpectedParentOID: base.OID})
	if e != nil {
		t.Fatal(e)
	}
	in := f.input
	in.Candidate = &Candidate{"repo", string(created.Task.Worktree.ID), strings.Repeat("f", 40), base.Tree}
	if _, e = f.s.Append("task", f.file(t, in), "event"); e == nil {
		t.Fatal("nonexistent Git candidate accepted as source-bound")
	}
	in.Candidate.OID = base.OID
	in.Candidate.Tree = strings.Repeat("e", 40)
	if _, e = f.s.Append("task", f.file(t, in), "event"); e == nil {
		t.Fatal("wrong candidate tree accepted")
	}
}

func TestSourceDisappearanceAtFreshnessInvalidatesOutcome(t *testing.T) {
	f := fixture(t)
	r, e := f.s.Workspace.WorkItems.Snapshot(f.root)
	if e != nil {
		t.Fatal(e)
	}
	base := r.Epics[0].RepoBindings[0].Worktree
	created, e := workspace.CreateTaskWorktree(f.s.Workspace, workspace.TaskWorktreeCreateInput{TaskID: "task", Branch: "task", Path: filepath.Join(f.root, "task-work"), ExpectedParentOID: base.OID})
	if e != nil {
		t.Fatal(e)
	}
	in := f.input
	in.Candidate = &Candidate{"repo", string(created.Task.Worktree.ID), base.OID, base.Tree}
	in.Outcome = ptr("pass")
	f.append(t, in)
	called := false
	f.s.fault = func(stage string) error {
		if stage == "before_freshness" {
			called = true
			return os.Remove(in.Sources[0].Locator)
		}
		return nil
	}
	s := f.show(t)
	if !called {
		t.Fatal("freshness fault boundary was not exercised")
	}
	if s.Coverage.State != "partial" || s.Current.Axes["reported"].Value != "unknown" {
		t.Fatal("source disappeared during read but current outcome stayed green")
	}
}
