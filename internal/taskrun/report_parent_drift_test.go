package taskrun

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReportStartAllowsEpicAdvanceButKeepsExactTaskSource(t *testing.T) {
	for _, changeTask := range []bool{false, true} {
		name := "epic-advance"
		if changeTask {
			name = "task-advance"
		}
		t.Run(name, func(t *testing.T) {
			f := workflowTestFixture(t)
			before, err := WorkflowPreviewStart(f.D, f.File)
			if err != nil {
				t.Fatal(err)
			}
			obs := before.(WorkflowPreview).Observed
			path := obs.Epic.WorktreeLocator
			if changeTask {
				path = obs.Target.WorktreeLocator
			}
			if err := os.WriteFile(filepath.Join(path, "independent.txt"), []byte("independent content\n"), 0600); err != nil {
				t.Fatal(err)
			}
			runGit(t, path, "add", "independent.txt")
			runGit(t, path, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "--allow-empty", "-m", "independent advance")
			after, err := WorkflowPreviewStart(f.D, f.File)
			if changeTask {
				if err == nil {
					t.Fatal("changed task source was accepted")
				}
			} else {
				if err != nil {
					t.Fatalf("report-only start blocked by an independent Epic commit: %v", err)
				}
				if after.(WorkflowPreview).Confirmation == nil {
					t.Fatal("no start confirmation")
				}
			}
		})
	}
}

func TestDeliveryStartStillRejectsEpicAdvance(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	before, err := WorkflowPreviewStart(f.D, f.File)
	if err != nil {
		t.Fatal(err)
	}
	path := before.(WorkflowPreview).Observed.Epic.WorktreeLocator
	if err := os.WriteFile(filepath.Join(path, "independent.txt"), []byte("independent content\n"), 0600); err != nil {
		t.Fatal(err)
	}
	runGit(t, path, "add", "independent.txt")
	runGit(t, path, "commit", "-m", "independent advance")
	if _, err := WorkflowPreviewStart(f.D, f.File); err == nil || !strings.Contains(err.Error(), "target or Epic changed since preparation") {
		t.Fatalf("expected exact parent guard, got %v", err)
	}
}
