package taskexecute

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestExecutePreviewIgnoresUnrelatedRemovedTaskWorktrees(t *testing.T) {
	for _, otherRepo := range []bool{false, true} {
		name := "same_repo"
		if otherRepo {
			name = "other_repo"
		}
		t.Run(name, func(t *testing.T) {
			d, in, root, epic := launcherFixture(t)
			projectID, repoID, epicID := workspace.ProjectID("fixture"), workspace.RepoID("fixture"), workspace.EpicID("epic")
			historyRoot := root
			if otherRepo {
				historyRoot = filepath.Join(root, "other")
				if err := os.Mkdir(historyRoot, 0700); err != nil {
					t.Fatal(err)
				}
				repo := filepath.Join(historyRoot, "main")
				fixtureGit(t, root, "clone", "--no-local", filepath.Join(root, "main"), repo)
				epic = filepath.Join(historyRoot, "epic")
				fixtureGit(t, repo, "worktree", "add", "-b", "other-epic", epic, "HEAD")
				projectID, repoID, epicID = "other", "other", "other-epic"
				if _, err := workspace.AddProject(d.Workspace, workspace.ProjectAddInput{ProjectID: projectID, Name: "Other repository", Wrapper: historyRoot, Repos: []workspace.RepoInput{{RepoID: repoID, Path: repo}}}); err != nil {
					t.Fatal(err)
				}
				if _, err := workspace.AdoptEpic(d.Workspace, workspace.EpicAdoptInput{EpicID: epicID, Title: "Other Epic", ProjectID: projectID, RepoID: repoID, Worktree: epic, Ref: "refs/heads/other-epic", ExpectedOID: fixtureGit(t, epic, "rev-parse", "HEAD")}); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := workspace.CreateTask(d.Workspace, workspace.TaskCreateInput{TaskID: "historical", Title: "Historical Task", Description: "Unrelated source removed after use.", ParentEpicID: epicID, ProjectID: projectID, RepoID: repoID}); err != nil {
				t.Fatal(err)
			}
			removed := filepath.Join(historyRoot, "historical")
			created, err := workspace.CreateTaskWorktree(d.Workspace, workspace.TaskWorktreeCreateInput{TaskID: "historical", Branch: "historical", Path: removed, ExpectedParentOID: fixtureGit(t, epic, "rev-parse", "HEAD")})
			if err != nil {
				t.Fatal(err)
			}
			if created.Task.WorktreeState == workspace.WorkItemRetired {
				t.Fatal("fixture must retain a non-retired Task registration")
			}
			fixtureGit(t, epic, "worktree", "remove", removed)
			beforeReservation := errors.New("fixture stops before native reservation or provider startup")
			reservationBoundary := 0
			d.Fault = func(point string) error {
				if point == "workflow_before_reservation" {
					reservationBoundary++
					return beforeReservation
				}
				return nil
			}
			out, executeErr := Execute(d, in)
			// The request exists only after genuine Goal preparation, including
			// worktree creation, expanded Spec selection and queue reservation.
			var request taskrun.WorkflowRequest
			if _, err := readPreservedJSON(out.RequestPath, &request); err != nil || request.PreparationID == "" {
				t.Fatalf("did not reach post-preparation startup: %v (execute: %v)", err, executeErr)
			}
			registryPath := filepath.Join(root, ".ply", "work-items.yaml")
			before, err := os.ReadFile(registryPath)
			if err != nil {
				t.Fatal(err)
			}
			preview, previewErr := taskrun.WorkflowPreviewStart(d, out.RequestPath)
			after, err := os.ReadFile(registryPath)
			if err != nil || !bytes.Equal(before, after) {
				t.Fatal("startup preview rewrote historical registrations")
			}
			r, err := d.Workspace.WorkItems.Snapshot(root)
			if err != nil {
				t.Fatal(err)
			}
			foundHistory := false
			for _, task := range r.Tasks {
				if task.ID == "historical" {
					foundHistory = true
					if !reflect.DeepEqual(task, created.Task) {
						t.Fatal("removed historical source was rewritten or claimed retired")
					}
				}
			}
			if !foundHistory {
				t.Fatal("startup removed the historical Task registration")
			}
			if _, err := os.Lstat(removed); !os.IsNotExist(err) {
				t.Fatal("historical worktree was recreated")
			}
			if previewErr != nil {
				t.Fatalf("post-preparation WorkflowPreviewStart rejected an unrelated removed Task: %v", previewErr)
			}
			p, ok := preview.(taskrun.WorkflowPreview)
			if !ok || p.Confirmation == nil || reservationBoundary != 1 || !errors.Is(executeErr, beforeReservation) {
				t.Fatalf("valid target did not reach the controlled start boundary: preview=%T reservations=%d error=%v", preview, reservationBoundary, executeErr)
			}
			if _, err := os.Stat(filepath.Join(root, ".ply", "workflow-runs")); !os.IsNotExist(err) {
				t.Fatal("preview or guarded execution reserved a Run")
			}
		})
	}
}
