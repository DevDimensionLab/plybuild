package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"strings"
	"testing"
)

func TestTaskAdoptionPublicCommands(t *testing.T) {
	for _, name := range []string{"adopt", "release-adoption"} {
		t.Run(name, func(t *testing.T) {
			out, err := executeTaskCommand(t, workspaceTaskServices{}, "worktree", name, "--help")
			if err != nil || !strings.Contains(out, "--confirmation") || !strings.Contains(out, "--apply") || !strings.Contains(out, "--format") {
				t.Fatalf("missing usable native %s preview/apply contract: err=%v output=%s", name, err, out)
			}
		})
	}
}

func TestTaskAdoptionNormativeArgv(t *testing.T) {
	calls := 0
	services := workspaceTaskServices{adopt: func(in workspace.TaskAdoptionInput, apply bool, confirmation string) (workspace.TaskAdoptionReadback, error) {
		calls++
		if in.TaskID != "task" || in.Path != "/source" || in.Ref != "refs/heads/useful" || in.Target.ProjectID != "ply" || in.Target.RepoID != "ply" || in.Target.EpicID != "epic" || apply && confirmation != "sha256:exact" {
			t.Fatalf("wrong native input %+v %v %s", in, apply, confirmation)
		}
		return workspace.TaskAdoptionReadback{Kind: "ply.workspace.task-worktree-adoption", Phase: "preview"}, nil
	}}
	args := []string{"worktree", "adopt", "task", "--project", "ply", "--repo", "ply", "--epic", "epic", "--path", "/source", "--ref", "refs/heads/useful", "--expected-oid", strings.Repeat("a", 40), "--format", "json"}
	if out, e := executeTaskCommand(t, services, args...); e != nil || !strings.Contains(out, `"phase":"preview"`) {
		t.Fatalf("normative preview failed: %v %s", e, out)
	}
	if _, e := executeTaskCommand(t, services, append(args, "--apply", "--confirmation", "sha256:exact")...); e != nil {
		t.Fatal(e)
	}
	if calls != 2 {
		t.Fatal("service calls differ")
	}
}

func TestTaskReleaseAdoptionNormativeArgv(t *testing.T) {
	calls := 0
	id := workspace.WorktreeOperationID("wop_" + strings.Repeat("a", 32))
	s := workspaceTaskServices{releaseAdoption: func(got workspace.WorktreeOperationID, reason string, apply bool, confirmation string) (workspace.TaskAdoptionReadback, error) {
		calls++
		if got != id || reason != "source changed" || apply && confirmation != "sha256:exact" {
			t.Fatal("wrong release identity/input")
		}
		return workspace.TaskAdoptionReadback{Kind: "ply.workspace.task-worktree-adoption", Phase: "released", MetadataOnly: true, WriterAbsence: "not_required"}, nil
	}}
	args := []string{"worktree", "release-adoption", string(id), "--reason", "source changed", "--format", "json"}
	if out, e := executeTaskCommand(t, s, args...); e != nil || !strings.Contains(out, `"metadata_only":true`) {
		t.Fatalf("release preview %v %s", e, out)
	}
	if _, e := executeTaskCommand(t, s, append(args, "--apply", "--confirmation", "sha256:exact")...); e != nil {
		t.Fatal(e)
	}
	if calls != 2 {
		t.Fatal("release did not use exactly one service call per request")
	}
}
