package taskrun

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func backupRetiredFixture(t *testing.T) (deliveryFixture, WorkflowRun) {
	t.Helper()
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
	retireAsCoordinator(&f)
	var err error
	o, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "modeled coordinator", Reason: "Fresh backup fixture"})
	if err != nil {
		t.Fatal(err)
	}
	return f, o
}
func TestBackupGoalRegistration(t *testing.T) {
	f, o := backupRetiredFixture(t)
	config := t.TempDir()
	t.Setenv("CODEX_HOME", config)
	if err := os.WriteFile(filepath.Join(config, "config.toml"), []byte("model = \"gpt-5.4\"\nmodel_reasoning_effort = \"high\"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	in := BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "codex"}
	before, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	preview, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	after, _ := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if !equal(before, after) {
		t.Fatal("preview mutated registry")
	}
	if *preview.Executor.Provider != "codex" || *preview.Executor.Model != "gpt-5.4" || *preview.Executor.Effort != "high" || !strings.Contains(preview.ExecutorSource, "Codex settings") {
		t.Fatalf("source executor leaked or fresh defaults missing: %+v", preview)
	}
	applied, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, preview.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	retry, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, preview.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if applied.Goal != retry.Goal || len(retry.Queue.Pending) != 1 || retry.Queue.Pending[0].TaskID != "backup" {
		t.Fatalf("bad retry: %+v", retry)
	}
	show, err := workspace.ShowTaskSpec(f.D.Workspace, workspace.TaskContentQuery{TaskID: "backup", SpecID: "backup-goal", Revision: 1})
	if err != nil {
		t.Fatal(err)
	}
	if show.Value == nil {
		t.Fatal("missing readback")
	}
	source := workflowProviderDocument(t, f.R.Delivery.Goal.Locator)
	freshPath := filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(applied.Goal.Spec.ManifestSHA256, "sha256:")+".json")
	fresh := workflowProviderDocument(t, freshPath)
	for _, key := range []string{"requirements", "documents", "design", "objective", "delivery"} {
		if !equal(source[key], fresh[key]) {
			t.Fatalf("product binding %s changed", key)
		}
	}
	manifestPath := func(ref any) string {
		m := ref.(map[string]any)
		return filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(m["manifest_sha256"].(string), "sha256:")+".json")
	}
	oldProblem := workflowProviderDocument(t, manifestPath(source["problem"]))
	newProblem := workflowProviderDocument(t, manifestPath(fresh["problem"]))
	for _, key := range []string{"documents", "claims", "sources", "deadline", "summary"} {
		if !equal(oldProblem[key], newProblem[key]) {
			t.Fatalf("problem binding %s changed", key)
		}
	}
	if !strings.Contains(fmt.Sprint(fresh["constraints"]), "normative for runtime identity") {
		t.Fatal("runtime delta is not normative")
	}
	registry, _ := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	for _, a := range registry.TaskSpecAssessments {
		if a.TaskID == "backup" {
			t.Fatal("inherited assessment")
		}
	}
	for _, r := range registry.TaskResults {
		if r.TaskID == "backup" {
			t.Fatal("inherited result")
		}
	}
	in.Provider = "claude"
	if _, err = WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, preview.Confirmation); err == nil {
		t.Fatal("changed provider accepted")
	}
}

func TestBackupGoalBoundaries(t *testing.T) {
	f, o := backupRetiredFixture(t)
	in := BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "claude"}
	// The supported bootstrap adapter, never a real provider, created/retired
	// the fixture. Every negative below must leave registered facts unchanged.
	original, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
	for _, tc := range []struct {
		name   string
		change func(*workflowState)
		input  BackupGoalInput
	}{
		{"delivered", func(s *workflowState) { s.Result.Delivery.Phase = "completed" }, in},
		{"unsupported-contract", func(s *workflowState) { s.Request.Delivery.Agreement.Mode = workspace.DeliveryPullRequest }, in},
		{"wrong-target", func(s *workflowState) { s.Request.Delivery.Agreement.TargetRef = "refs/heads/other" }, in},
		{"foreign-task", nil, BackupGoalInput{TaskID: "task", SpecID: "fresh-goal", Provider: "claude"}},
		{"bad-provider", nil, BackupGoalInput{TaskID: "backup", SpecID: "fresh-goal", Provider: "other"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Isolated modeled state uses the ordinary native store; no real Run touched.
			if tc.change != nil {
				if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error { tc.change(s); return nil }); err != nil {
					t.Fatal(err)
				}
			}
			before, _ := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
			if _, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, tc.input, false, ""); err == nil {
				t.Fatal("unsupported source accepted")
			}
			after, _ := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
			if !equal(before, after) {
				t.Fatal("negative mutated registration")
			}
			if err := workflowSave(f.D, original); err != nil {
				t.Fatal(err)
			}
		})
	}
}

type backupFailStore struct {
	workspace.WorkItemStore
	remaining *int
}
type backupFailSession struct {
	workspace.WorkItemStoreSession
	remaining *int
}

func (s backupFailStore) WithLock(root string, action func(workspace.WorkItemStoreSession) error) error {
	return s.WorkItemStore.WithLock(root, func(session workspace.WorkItemStoreSession) error {
		return action(backupFailSession{session, s.remaining})
	})
}
func (s backupFailSession) Publish(r workspace.WorkItemRegistry) error {
	if err := s.WorkItemStoreSession.Publish(r); err != nil {
		return err
	}
	*s.remaining--
	if *s.remaining == 0 {
		return fmt.Errorf("modeled lost publication acknowledgement")
	}
	return nil
}
func TestBackupGoalPartialPublicationAndQueuePreservation(t *testing.T) {
	f, o := backupRetiredFixture(t)
	// Existing fixture helper only registers/prepares unrelated synthetic Goals;
	// it starts no provider and makes no PR/integration/closeout call.
	integrationCLIPrepareSuccessor(t, f)
	target := workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}
	before, err := workspace.ListTaskQueue(f.D.Workspace, target, false)
	if err != nil {
		t.Fatal(err)
	}
	in := BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "claude"}
	p, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	n := 2
	store := f.D.Workspace.WorkItems
	f.D.Workspace.WorkItems = backupFailStore{store, &n}
	if _, err = WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation); err == nil {
		t.Fatal("lost acknowledgement not reported")
	}
	f.D.Workspace.WorkItems = store
	result, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if !equal(before.Current, result.Queue.Current) || !equal(before.Pending, result.Queue.Pending[:len(before.Pending)]) || len(result.Queue.Pending) != len(before.Pending)+1 {
		t.Fatal("current or pending order/refs changed")
	}
	after, _ := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	retry, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	final, _ := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if !equal(after, final) || result.Goal != retry.Goal {
		t.Fatal("retry added history")
	}
}

func TestBackupGoalRequiresRetiredSource(t *testing.T) {
	f, o := retireFixture(t)
	before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	if _, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "claude"}, false, ""); err == nil || !strings.Contains(err.Error(), "source_not_retired") {
		t.Fatalf("live source: %v", err)
	}
	if before != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatal("live-source preview wrote state")
	}
}

// Real candidate CLI covers only the new registration/readback boundary.
// Bootstrap uses the existing modeled Herdr/process adapters; no OS absence
// claim or actual provider retirement/start is made by this test.
func TestBackupGoalCandidateCLI(t *testing.T) {
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		t.Skip("requires explicit make-built candidate CLI")
	}
	f, o := backupRetiredFixture(t)
	config := t.TempDir()
	t.Setenv("CODEX_HOME", config)
	if err := os.WriteFile(filepath.Join(config, "config.toml"), []byte("model='gpt-5.4'\nmodel_reasoning_effort='high'\n"), 0600); err != nil {
		t.Fatal(err)
	}
	evidence := os.Getenv("PLY_AL27_EVIDENCE_DIR")
	if evidence == "" {
		evidence = t.TempDir()
	}
	if err := os.MkdirAll(evidence, 0700); err != nil {
		t.Fatal(err)
	}
	sequence := 0
	lastStderr := ""
	cli := func(args ...string) ([]byte, error) {
		sequence++
		cmd := exec.Command(binary, args...)
		cmd.Dir = f.Parent
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		lastStderr = stderr.String()
		exit := 0
		if err != nil {
			exit = -1
			var ee *exec.ExitError
			if errors.As(err, &ee) {
				exit = ee.ExitCode()
			}
		}
		prefix := filepath.Join(evidence, fmt.Sprintf("cli-%02d", sequence))
		for suffix, data := range map[string][]byte{".stdout": stdout.Bytes(), ".stderr": stderr.Bytes(), ".command.json": mustBackupJSON(t, map[string]any{"argv": cmd.Args, "cwd": cmd.Dir, "exit": exit})} {
			if e := os.WriteFile(prefix+suffix, data, 0600); e != nil {
				t.Fatal(e)
			}
		}
		return stdout.Bytes(), err
	}
	oldRun := deliveryCLISnapshot(t, o.Paths.RunRoot)
	args := []string{"workflow", "execute", "backup-goal", o.RunID, "--task", "backup", "--spec", "backup-goal", "--provider", "codex", "--format", "json"}
	before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	raw, err := cli(args...)
	if err != nil {
		t.Fatalf("preview: %v %s", err, raw)
	}
	if before != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatal("CLI preview wrote workspace")
	}
	var preview workspace.BackupGoalRegistrationResult
	if err = json.Unmarshal(raw, &preview); err != nil {
		t.Fatal(err)
	}
	if *preview.Executor.Model != "gpt-5.4" {
		t.Fatal("provider-only CLI failed to resolve native defaults")
	}
	apply := append(append([]string{}, args...), "--apply", "--confirmation", preview.Confirmation)
	raw, err = cli(apply...)
	if err != nil {
		t.Fatalf("apply: %v %s", err, raw)
	}
	var result workspace.BackupGoalRegistrationResult
	if err = json.Unmarshal(raw, &result); err != nil {
		t.Fatal(err)
	}
	after := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	raw, err = cli(apply...)
	if err != nil {
		t.Fatalf("retry: %v %s", err, raw)
	}
	var retry workspace.BackupGoalRegistrationResult
	if err = json.Unmarshal(raw, &retry); err != nil {
		t.Fatal(err)
	}
	if retry.Goal != result.Goal || after != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatal("CLI lost-reply retry changed registration")
	}
	for _, read := range [][]string{
		{"workspace", "task", "spec", "show", "backup", "--spec", "backup-goal", "--revision", "1", "--format", "json"},
		{"workspace", "task", "show", "backup", "--format", "json"},
		{"workspace", "task", "queue", "list", "--format", "json"},
	} {
		raw, err = cli(read...)
		if err != nil {
			t.Fatalf("readback %v: %v %s", read, err, raw)
		}
	}
	// Ordinary execute check can truthfully report missing local launch settings.
	// Registration itself requires neither a Herdr workspace nor permissions.
	before = deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	raw, err = cli("workflow", "execute", "--spec", "backup-goal", "--check", "--format", "json")
	if err != nil && !strings.Contains(lastStderr, "Codex has no selected permission profile") {
		t.Fatalf("start-check had no truthful blocked readback: %v %s", err, raw)
	}
	if before != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatal("start-check wrote workspace")
	}
	if oldRun != deliveryCLISnapshot(t, o.Paths.RunRoot) {
		t.Fatal("registration changed retired source Run history")
	}
	t.Logf("New CLI preview/apply/retry/spec/task/queue/start-check only; modeled retired source; evidence %s; start-check exit error=%v", evidence, err)
}
func mustBackupJSON(t *testing.T, value any) []byte {
	t.Helper()
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return raw
}
