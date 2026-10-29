package taskexecute

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/mitchellh/go-homedir"
)

func projectDataProfile(t *testing.T) string {
	t.Helper()
	home := physicalTemp(t)
	t.Setenv("HOME", home)
	homedir.Reset()
	t.Cleanup(homedir.Reset)
	profiles := filepath.Join(home, ".ply", "profiles")
	profile := filepath.Join(profiles, "fixture")
	if err := os.MkdirAll(profile, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(profiles, ".active_profile"), []byte("fixture\n"), 0600); err != nil {
		t.Fatal(err)
	}
	return filepath.Join(profile, "local-config.yaml")
}

func setProjectDataRoot(t *testing.T, profile, root string) {
	t.Helper()
	if err := os.WriteFile(profile, []byte("projects:\n  fixture:\n    data_root: "+root+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestExecuteFreezesConfiguredDataRootBeforeIntent(t *testing.T) {
	profile := projectDataProfile(t)
	d, in, root, epic := launcherFixtureForProvider(t, "codex")
	if err := os.WriteFile(in.Runtime.HerdrExecutable, []byte(recoveryLauncherHerdr), 0700); err != nil {
		t.Fatal(err)
	}
	data := filepath.Join(root, "data-one")
	setProjectDataRoot(t, profile, data)
	in.Check = true
	preview, err := Execute(d, in)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(preview.Goal.AcceptancePath, data+"/task-executions/") {
		t.Fatalf("configured data root ignored: %s", preview.Goal.AcceptancePath)
	}
	if _, err = os.Stat(data); !os.IsNotExist(err) {
		t.Fatal("check created data root")
	}
	in.Check = false
	d.Fault = func(point string) error {
		if point == "execute_before_preparation" {
			return errors.New("fixture before preparation")
		}
		return nil
	}
	first, err := Execute(d, in)
	if err == nil || !strings.Contains(err.Error(), "fixture before preparation") {
		t.Fatal(err)
	}
	r, err := d.Workspace.WorkItems.Snapshot(root)
	if err != nil || len(r.TaskSpecRevisions) != 2 || len(r.TaskSpecAssessments) != 0 || len(r.TaskPreparations) != 0 {
		t.Fatalf("Spec alone must precede intent: specs=%d assessments=%d preparations=%d err=%v", len(r.TaskSpecRevisions), len(r.TaskSpecAssessments), len(r.TaskPreparations), err)
	}
	if _, err = os.Stat(first.Goal.WorktreePath); !os.IsNotExist(err) {
		t.Fatal("early Spec created Task worktree")
	}
	id := workspace.TaskID("goal")
	prepared, err := workspace.PrepareTask(d.Workspace, workspace.TaskPrepareInput{Selector: workspace.QueueSelector{Kind: "task", TaskID: &id}})
	if err != nil || prepared.State != "blocked" || prepared.Plan != nil || prepared.Preparation != nil {
		t.Fatalf("early Spec allowed direct preparation: %+v %v", prepared, err)
	}
	intentPath := filepath.Join(filepath.Dir(first.RequestPath), "intent.json")
	before, err := os.ReadFile(intentPath)
	if err != nil {
		t.Fatal(err)
	}
	setProjectDataRoot(t, profile, filepath.Join(root, "data-two"))
	if err := os.WriteFile(filepath.Join(root, "codex-home", "config.toml"), []byte("default_permissions = 'changed-after-intent'\n"), 0600); err != nil {
		t.Fatal(err)
	}
	subdir := filepath.Join(epic, "nested")
	if err = os.Mkdir(subdir, 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.Chdir(subdir); err != nil {
		t.Fatal(err)
	}
	second, err := Execute(d, in)
	if err == nil || !strings.Contains(err.Error(), "fixture before preparation") || first.RequestPath != second.RequestPath || first.Goal.Confirmation != second.Goal.Confirmation {
		t.Fatalf("root change lost frozen execution: %v", err)
	}
	after, err := os.ReadFile(intentPath)
	if err != nil || string(after) != string(before) {
		t.Fatal("retry changed preserved intent")
	}
	d.Fault = nil
	started, err := Execute(d, in)
	if err != nil || started.State != "started" || started.Run == nil || started.RequestPath != first.RequestPath {
		t.Fatalf("frozen attempt could not resume: %+v %v", started, err)
	}
	repeated, err := Execute(d, in)
	if err != nil || repeated.State != "existing" || repeated.Run == nil || repeated.Run.RunID != started.Run.RunID {
		t.Fatalf("recovery restarted Run: %+v %v", repeated, err)
	}
	var model map[string]any
	raw, err := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json"))
	if err != nil || json.Unmarshal(raw, &model) != nil || model["starts"] != float64(1) || model["prompts"] != float64(1) {
		t.Fatalf("recovery duplicated provider: %s %v", raw, err)
	}
	for _, absent := range []string{filepath.Join(root, "data-two"), filepath.Join(root, ".ply", "task-executions")} {
		if _, err = os.Stat(absent); !os.IsNotExist(err) {
			t.Fatalf("retry wrote outside frozen data root: %s", absent)
		}
	}
}

func TestExecuteRejectsInvalidDataRootBeforePublication(t *testing.T) {
	profile := projectDataProfile(t)
	d, in, root, _ := launcherFixtureForProvider(t, "codex")
	link := filepath.Join(root, "linked-data")
	if err := os.Symlink(root, link); err != nil {
		t.Fatal(err)
	}
	before, err := d.Workspace.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"relative/data", root + "/../unclean", link, filepath.Join(root, "main", "README.md")} {
		setProjectDataRoot(t, profile, bad)
		if _, err := Execute(d, in); err == nil {
			t.Fatalf("accepted invalid data_root %q", bad)
		}
		after, err := d.Workspace.WorkItems.Snapshot(root)
		if err != nil || before.RawSHA256 != after.RawSHA256 {
			t.Fatalf("invalid data_root published state: %v", err)
		}
	}
}

func TestConcurrentExecutePublishesOneConfiguredAttempt(t *testing.T) {
	d, in, root, _ := launcherFixtureForProvider(t, "codex")
	ready, release := make(chan struct{}, 2), make(chan struct{})
	results := make(chan error, 2)
	for _, name := range []string{"first", "second"} {
		copy := d
		data := filepath.Join(root, name)
		copy.Workspace.ProjectDataRoot = func(string) (string, error) { return data, nil }
		copy.Executable = func() (string, error) {
			ready <- struct{}{}
			<-release
			return d.Executable()
		}
		copy.Fault = func(point string) error {
			if point == "execute_before_preparation" {
				return errors.New("fixture before preparation")
			}
			return nil
		}
		go func() { _, err := Execute(copy, in); results <- err }()
	}
	<-ready
	<-ready
	close(release)
	frozen, conflicted := 0, 0
	for i := 0; i < 2; i++ {
		err := <-results
		if err != nil && strings.Contains(err.Error(), "fixture before preparation") {
			frozen++
		} else if err != nil && strings.Contains(err.Error(), "task_goal_execution_conflict") {
			conflicted++
		} else {
			t.Fatalf("unexpected concurrent result: %v", err)
		}
	}
	if frozen != 1 || conflicted != 1 {
		t.Fatalf("concurrent outcomes: frozen=%d conflicted=%d", frozen, conflicted)
	}
	r, err := d.Workspace.WorkItems.Snapshot(root)
	if err != nil || len(r.TaskSpecRevisions) != 2 || len(r.TaskPreparations) != 0 {
		t.Fatalf("duplicate publication: %v", err)
	}
	areas := 0
	for _, name := range []string{"first", "second"} {
		if _, err := os.Stat(filepath.Join(root, name)); err == nil {
			areas++
		}
	}
	if areas != 1 {
		t.Fatalf("concurrent calls wrote %d attempt roots", areas)
	}
}

func TestBoundStartUsesFrozenRootAfterIntentPublication(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	profile := projectDataProfile(t)
	firstRoot, secondRoot := filepath.Join(root, "first"), filepath.Join(root, "second")
	setProjectDataRoot(t, profile, firstRoot)
	p := boundPreview(t, d, in)
	d.Fault = func(point string) error {
		if point == "bound_start_after_intent" {
			setProjectDataRoot(t, profile, secondRoot)
		}
		return nil
	}
	started, err := ApplyBoundStart(d, p.Request)
	if err != nil || started.Run == nil || started.Run.RunID != boundRunID(p.Request) {
		t.Fatalf("bound start lost frozen basis: %+v %v", started, err)
	}
	if _, err := os.Stat(firstRoot); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(secondRoot); !os.IsNotExist(err) {
		t.Fatal("bound start wrote into new profile root")
	}
	d.Fault = nil
	repeated, err := ApplyBoundStart(d, p.Request)
	if err != nil || repeated.Run == nil || repeated.Run.RunID != started.Run.RunID {
		t.Fatalf("bound replay changed Run: %+v %v", repeated, err)
	}
}
