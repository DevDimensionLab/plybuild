package workspace

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParallelOptionalUpgradeTargetsStayFiveSixSeven(t *testing.T) {
	directory, pin := deliveredRegistryReader(t)
	f := newWorkItemJourneyFixture(t)
	queueSetFixture(t, f, "ordinary-four", []TaskGoalQueueEntry{})
	for i, upgrade := range []func(Dependencies, string, string, string, bool) (RegistryUpgradePlan, error){UpgradeRegistry, UpgradeEpicCreateRegistry, UpgradeAdoptionRegistry} {
		version := i + 5
		before := mustQueueRegistry(t, f)
		preview, e := upgrade(f.dependencies, directory, pin, "", false)
		if e != nil || preview.ToVersion != version || preview.FromVersion != version-1 {
			t.Fatalf("explicit upgrade target %d: %+v %v", version, preview, e)
		}
		if _, e = upgrade(f.dependencies, directory, pin, preview.Confirmation, true); e != nil {
			t.Fatal(e)
		}
		after := mustQueueRegistry(t, f)
		backup, e := os.ReadFile(preview.Backup)
		if e != nil || digestTaskBytes(backup) != before.RawSHA256 {
			t.Fatal("legacy optional upgrade lost exact preimage")
		}
		if after.FormatVersion != version || after.ParallelBoundary != nil || after.ParallelActivation != nil {
			t.Fatal("reader max8 selected an ordinary format8 upgrade")
		}
		registration := firstGoalFixture(t, f)
		registration.TaskID = TaskID(fmt.Sprintf("ordinary-%d", version))
		registration.SpecID = fmt.Sprintf("ordinary-goal-%d", version)
		p, e := RegisterFirstGoal(f.dependencies, registration, false, "")
		if e != nil {
			t.Fatal(e)
		}
		if _, e = RegisterFirstGoal(f.dependencies, registration, true, p.Confirmation); e != nil {
			t.Fatal(e)
		}
		if mustQueueRegistry(t, f).FormatVersion != version {
			t.Fatal("ordinary registration changed actual optional format")
		}
	}
}

func TestParallelReaderPreservesEveryLegacyPredecessor(t *testing.T) {
	for _, version := range []int{4, 5, 6, 7} {
		t.Run(string(rune('0'+version)), func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			queueSetFixture(t, f, "legacy-seed", []TaskGoalQueueEntry{})
			r := mustQueueRegistry(t, f)
			boundary := func(from int) *RegistryFormatBoundary {
				return &RegistryFormatBoundary{from, r.RawSHA256, "sha256:" + strings.Repeat("a", 64), "v0.0.0-simulated", f.oid, "sha256:" + strings.Repeat("b", 64)}
			}
			r.FormatVersion = version
			if version >= 5 {
				r.FormatBoundary = boundary(4)
			}
			if version >= 6 {
				r.EpicCreateBoundary = boundary(5)
				r.EpicCreateOperations = []EpicCreateOperation{}
			}
			if version >= 7 {
				r.AdoptionBoundary = boundary(6)
			}
			before, e := encodeWorkItemRegistry(r)
			if e != nil {
				t.Fatal(e)
			}
			legacy, e := decodeWorkItemRegistry(before)
			if e != nil {
				t.Fatal(e)
			}
			legacy.ParallelBoundary = boundary(version)
			legacy.FormatVersion = 8
			legacy.SourceReconciliations = []HistoricalRecord{}
			raw, e := encodeWorkItemRegistry(legacy)
			if e != nil {
				t.Fatal(e)
			}
			after, e := decodeWorkItemRegistry(raw)
			if e != nil {
				t.Fatal(e)
			}
			after.FormatVersion = version
			after.ParallelBoundary = nil
			after.SourceReconciliations = nil
			restored, e := encodeWorkItemRegistry(after)
			if e != nil || !bytes.Equal(before, restored) {
				t.Fatal("parallel reader rewrote/fabricated predecessor history")
			}
		})
	}
}

func TestParallelActivationReadOnlyCASAndPreservedFormatFour(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	queueSetFixture(t, f, "format-four", []TaskGoalQueueEntry{})
	observe := func(r WorkItemRegistry) (ParallelActivationEvidence, error) {
		ep, b := queueTargetBinding(r, QueueTarget{EpicID: "epic", RepoID: "ply"})
		return ParallelActivationEvidence{Target: QueueTarget{ep.ProjectID, b.RepoID, ep.ID, b.GitCommonDir, b.Worktree.ID, b.Worktree.Locator, b.Worktree.Ref}, ReaderVersion: "v0.0.0-synthetic", ReaderCommit: f.oid, ReaderSHA256: "sha256:" + strings.Repeat("a", 64), QualifiedManifestSHA256: "sha256:" + strings.Repeat("b", 64), InventorySHA256: "sha256:" + strings.Repeat("c", 64), Consumers: []ParallelConsumerCompatibility{}}, nil
	}
	before := queueStateBytes(t, f)
	preview, err := ActivateParallelRegistry(f.dependencies, "", false, observe)
	if err != nil {
		t.Fatalf("native qualified activation preview: %v", err)
	}
	if before != queueStateBytes(t, f) {
		t.Fatal("activation preview wrote state")
	}
	queueSetFixture(t, f, "planner-change", []TaskGoalQueueEntry{})
	if _, err = ActivateParallelRegistry(f.dependencies, preview.Confirmation, true, observe); err == nil {
		t.Fatal("stale activation digest wrote format8")
	}
	if mustQueueRegistry(t, f).FormatVersion != 4 {
		t.Fatal("failed activation upgraded format")
	}
	preview, err = ActivateParallelRegistry(f.dependencies, "", false, observe)
	if err != nil {
		t.Fatal(err)
	}
	current := mustQueueRegistry(t, f)
	result, err := ActivateParallelRegistry(f.dependencies, preview.Confirmation, true, observe)
	if err != nil || result.State != "activated" {
		t.Fatalf("exact qualified activation: %+v %v", result, err)
	}
	r := mustQueueRegistry(t, f)
	if r.FormatVersion != 8 || r.ParallelBoundary.FromVersion != 4 || r.FormatBoundary != nil || r.AdoptionBoundary != nil || !contentTypedEqual(current.TaskQueueEvents, r.TaskQueueEvents) {
		t.Fatal("activation rewrote old history or fabricated intermediate upgrades")
	}
	backup, err := os.ReadFile(result.Backup)
	if err != nil || digestTaskBytes(backup) != current.RawSHA256 {
		t.Fatalf("exact old-format backup: %v", err)
	}
	before = queueStateBytes(t, f)
	again, err := ActivateParallelRegistry(f.dependencies, preview.Confirmation, true, observe)
	if err != nil || again.State != "activated" || before != queueStateBytes(t, f) {
		t.Fatalf("activation retry repeated effects: %v", err)
	}
	if _, err = os.Stat(filepath.Join(f.wrapper, "main", ".git", "config")); err != nil {
		t.Fatal(err)
	}
}
