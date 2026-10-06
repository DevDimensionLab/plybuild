package workspaceview

import (
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestEpicsIncludeEmptyAndUseCurrentRecordedBaseAndSharedCounts(t *testing.T) {
	s := progressFixture()
	s.Registry.FormatVersion = 4
	s.Registry.Epics[0].RepoBindings = []workspace.EpicRepoBinding{{RepoID: "repo", Worktree: workspace.EpicWorktreeBinding{OID: "adopted"}}}
	s.Registry.Epics = append(s.Registry.Epics, workspace.EpicRecord{ID: "empty", ProjectID: "project", Title: "Empty Epic"})
	s.Registry.EpicBaseVersions = []workspace.EpicBaseVersion{{EpicID: "epic", RepoID: "repo", Revision: 1, OID: "adopted"}, {EpicID: "epic", RepoID: "repo", Revision: 2, OID: "current", Tree: "current-tree", Locator: "/work/project/epic", Ref: "refs/heads/epic", WorktreeID: "worktree"}}
	s.Registry.TaskResults = []workspace.TaskResultRecord{fixtureResult("passed")}
	list, err := BuildEpics(s, EpicFilters{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Epics) != 2 || list.Epics[0].EpicID != "empty" || list.Epics[0].TaskCount != 0 || list.Epics[0].LastActivityUTC != nil || list.Epics[0].Base == nil || list.Epics[0].Worktrees == nil {
		t.Fatalf("empty Epic missing: %+v", list)
	}
	row := list.Epics[1]
	if len(row.Base) != 1 || row.Base[0].Revision != 2 || row.Base[0].OID != "current" || row.Base[0].SourceBasis != "registered" || row.Worktrees[0].Locator != "/work/project/epic" {
		t.Fatalf("adoption confused with current base: %+v", row)
	}
	if row.TaskCount != 1 || row.ProgressCounts.InProgress != 1 || row.LastActivityUTC == nil || *row.LastActivityUTC != "2026-10-06T08:00:00Z" {
		t.Fatalf("counts or activity disagree with Task view: %+v", row)
	}
}

func TestPreparedTaskHasProgressButNoClaimOfRunningOrQAPass(t *testing.T) {
	s := progressFixture()
	s.Registry.TaskPreparations = []workspace.TaskPreparation{{Plan: workspace.WorkspaceTaskPreparePlan{TaskID: "task"}, CreatedAtUTC: "2026-10-06T08:00:00Z"}}
	p, err := BuildTaskProgress(s, "task", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !p.HasProgress || p.State != "in_progress" || p.IntegrationClassification != "unknown" || p.TechnicalGate != nil || p.HumanQAOutcome != nil {
		t.Fatalf("prepared progress=%+v", p)
	}
	a, err := BuildAttention(s, workspace.TaskListFilters{}, nil)
	if err != nil || len(a.Items) != 0 {
		t.Fatalf("prepared Task needs no human: %+v %v", a, err)
	}
}
