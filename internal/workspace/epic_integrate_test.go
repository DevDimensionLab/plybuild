package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

type epicIntegrateFixture struct {
	f              workItemJourneyFixture
	main, epic     string
	mainOID, epicO string
}

func newEpicIntegrateFixture(t *testing.T) epicIntegrateFixture {
	t.Helper()
	f := newWorkItemJourneyFixture(t)
	x := epicIntegrateFixture{f: f, main: physicalPath(t, filepath.Join(f.wrapper, "main")), epic: physicalPath(t, filepath.Join(f.wrapper, "epic"))}
	if err := os.WriteFile(filepath.Join(x.epic, "epic.txt"), []byte("epic work\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, x.epic, "add", "epic.txt")
	runLocalGit(t, x.epic, "commit", "-m", "epic work")
	x.mainOID, x.epicO = runLocalGit(t, x.main, "rev-parse", "HEAD"), runLocalGit(t, x.epic, "rev-parse", "HEAD")
	return x
}

func (x epicIntegrateFixture) input(apply bool) EpicIntegrateInput {
	return EpicIntegrateInput{EpicID: "epic", RepoID: "ply", TargetRef: "refs/heads/main", ExpectedTargetOID: x.mainOID, ExpectedEpicOID: x.epicO, Apply: apply}
}

func TestEpicIntegrateFastForwardsCheckedOutTarget(t *testing.T) {
	x := newEpicIntegrateFixture(t)
	check, err := IntegrateEpic(x.f.dependencies, x.input(false))
	if err != nil || check.Applied || check.Effect != "would_fast_forward" || check.TargetWorktree != x.main {
		t.Fatalf("check %+v %v", check, err)
	}
	if runLocalGit(t, x.main, "rev-parse", "HEAD") != x.mainOID {
		t.Fatal("check moved main")
	}
	applied, err := IntegrateEpic(x.f.dependencies, x.input(true))
	if err != nil || !applied.Applied || applied.Effect != "fast_forwarded" || applied.TargetOID != x.epicO {
		t.Fatalf("apply %+v %v", applied, err)
	}
	if runLocalGit(t, x.main, "rev-parse", "HEAD") != x.epicO || runLocalGit(t, x.main, "status", "--porcelain") != "" {
		t.Fatal("main worktree not fast-forwarded cleanly")
	}
	if reflog := runLocalGit(t, x.main, "reflog", "-1", "--format=%gs", "refs/heads/main"); reflog != "ply-workspace-epic-integrate:epic: Fast-forward" {
		t.Fatalf("reflog %q", reflog)
	}
	in := x.input(true)
	in.ExpectedTargetOID = x.epicO
	again, err := IntegrateEpic(x.f.dependencies, in)
	if err != nil || again.Effect != "already_integrated" || again.Applied {
		t.Fatalf("repeat %+v %v", again, err)
	}
}

func TestEpicIntegrateUpdatesTargetNotCheckedOut(t *testing.T) {
	x := newEpicIntegrateFixture(t)
	runLocalGit(t, x.main, "branch", "release", x.mainOID)
	in := x.input(true)
	in.TargetRef = "refs/heads/release"
	applied, err := IntegrateEpic(x.f.dependencies, in)
	if err != nil || applied.Effect != "fast_forwarded" || applied.TargetWorktree != "" || runLocalGit(t, x.main, "rev-parse", "release") != x.epicO {
		t.Fatalf("apply %+v %v", applied, err)
	}
}

func TestEpicIntegrateRefusesStaleDivergedOrDirtyTarget(t *testing.T) {
	cases := map[string]func(t *testing.T, x *epicIntegrateFixture, in *EpicIntegrateInput) WorkItemErrorClass{
		"stale target": func(t *testing.T, x *epicIntegrateFixture, in *EpicIntegrateInput) WorkItemErrorClass {
			in.ExpectedTargetOID = x.epicO
			return ErrorWorkParentStale
		},
		"stale epic": func(t *testing.T, x *epicIntegrateFixture, in *EpicIntegrateInput) WorkItemErrorClass {
			in.ExpectedEpicOID = x.mainOID
			return ErrorWorkParentStale
		},
		"diverged": func(t *testing.T, x *epicIntegrateFixture, in *EpicIntegrateInput) WorkItemErrorClass {
			_ = os.WriteFile(filepath.Join(x.main, "main.txt"), []byte("main\n"), 0o644)
			runLocalGit(t, x.main, "add", "main.txt")
			runLocalGit(t, x.main, "commit", "-m", "main work")
			x.mainOID = runLocalGit(t, x.main, "rev-parse", "HEAD")
			in.ExpectedTargetOID = x.mainOID
			return ErrorWorkRefConflict
		},
		"dirty": func(t *testing.T, x *epicIntegrateFixture, in *EpicIntegrateInput) WorkItemErrorClass {
			_ = os.WriteFile(filepath.Join(x.main, "untracked.txt"), []byte("x\n"), 0o644)
			return ErrorWorkParentDirty
		},
	}
	for name, setup := range cases {
		t.Run(name, func(t *testing.T) {
			x := newEpicIntegrateFixture(t)
			in := x.input(true)
			want := setup(t, &x, &in)
			_, err := IntegrateEpic(x.f.dependencies, in)
			requireWorkClass(t, err, want)
			if runLocalGit(t, x.main, "rev-parse", "HEAD") != x.mainOID {
				t.Fatal("refused integration moved main")
			}
		})
	}
}
