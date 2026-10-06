package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspaceview"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func workflowStatusTestRoot(t *testing.T, status *cobra.Command) (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	root := &cobra.Command{Use: "ply", PersistentPreRunE: func(*cobra.Command, []string) error {
		t.Error("workflow status invoked profile initialization")
		return errors.New("profile initialization must not run")
	}}
	root.PersistentFlags().Bool("json", false, "logging")
	root.PersistentFlags().Bool("debug", false, "debug logging")
	workflow := &cobra.Command{Use: "workflow"}
	workflow.AddCommand(status)
	root.AddCommand(workflow)
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SilenceErrors, root.SilenceUsage = true, true
	return root, &stdout, &stderr
}

func emptyWorkflowStatus() workspaceview.WorkflowStatus {
	return workspaceview.WorkflowStatus{
		Kind: workspaceview.WorkflowStatusKind, SchemaVersion: 1,
		Workspace:     workspaceview.WorkspaceRef{Root: "/work"},
		ObservedAtUTC: "2026-10-06T18:00:00Z", SourceBasis: "registered", Freshness: "fresh",
		Projects: []workspaceview.WorkflowStatusProject{}, Epics: []workspaceview.EpicRow{},
		Items: []workspaceview.WorkflowStatusItem{}, Diagnostics: []workspaceview.WorkflowDiagnostic{},
		Counts: workspaceview.WorkflowStatusCounts{Categories: map[string]int{}},
	}
}

func TestWorkflowStatusJSONFlagPlacementAndFormat(t *testing.T) {
	want := emptyWorkflowStatus()
	for _, args := range [][]string{
		{"workflow", "status", "--format", "json"},
		{"workflow", "status", "--json"},
		{"--json", "workflow", "status"},
		{"workflow", "--json", "status"},
		{"--json", "workflow", "status", "--format", "json"},
		{"workflow", "status", "--format", "json", "--json"},
		{"workflow", "status", "--json=false", "--format", "json"},
		{"--debug", "--json", "workflow", "status"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			calls := 0
			root, out, stderr := workflowStatusTestRoot(t, newWorkflowStatusCommandWithReader(func(workspaceview.WorkflowStatusOptions) (workspaceview.WorkflowStatus, error) {
				calls++
				return want, nil
			}))
			root.SetArgs(args)
			if err := root.Execute(); err != nil || calls != 1 || stderr.Len() != 0 {
				t.Fatalf("err=%v calls=%d stderr=%s", err, calls, stderr)
			}
			if !strings.HasSuffix(out.String(), "\n") || strings.Count(out.String(), "\n") != 1 {
				t.Fatalf("expected exactly one JSON line: %q", out.String())
			}
			var got workspaceview.WorkflowStatus
			decoder := json.NewDecoder(out)
			if err := decoder.Decode(&got); err != nil || !reflect.DeepEqual(got, want) {
				t.Fatalf("decode=%v got=%+v", err, got)
			}
			var extra any
			if err := decoder.Decode(&extra); err != io.EOF {
				t.Fatalf("extra stdout after result: %v", err)
			}
		})
	}
	for _, args := range [][]string{
		{"workflow", "status"}, {"--json=false", "workflow", "status"},
		{"workflow", "status", "--json=false", "--format", "text"},
	} {
		root, out, _ := workflowStatusTestRoot(t, newWorkflowStatusCommandWithReader(func(workspaceview.WorkflowStatusOptions) (workspaceview.WorkflowStatus, error) { return want, nil }))
		root.SetArgs(args)
		if err := root.Execute(); err != nil || !strings.HasPrefix(out.String(), "Workflow status — /work\n") {
			t.Fatalf("args=%v err=%v text=%s", args, err, out)
		}
	}
}

func TestWorkflowStatusValidatesBeforeReading(t *testing.T) {
	for _, args := range [][]string{
		{"workflow", "status", "extra"}, {"workflow", "status", "--format", "yaml"},
		{"workflow", "status", "--format"}, {"workflow", "status", "--unknown"},
		{"workflow", "status", "--project", ""}, {"workflow", "status", "--repo", "Bad"},
		{"workflow", "status", "--epic", ""},
		{"--json", "workflow", "status", "--format", "text"},
		{"workflow", "status", "--format", "text", "--json"},
		{"workflow", "status", "--json", "--format", "text"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			calls := 0
			root, out, _ := workflowStatusTestRoot(t, newWorkflowStatusCommandWithReader(func(workspaceview.WorkflowStatusOptions) (workspaceview.WorkflowStatus, error) {
				calls++
				return emptyWorkflowStatus(), nil
			}))
			root.SetArgs(args)
			if err := root.Execute(); err == nil || calls != 0 || out.Len() != 0 {
				t.Fatalf("err=%v calls=%d stdout=%s", err, calls, out)
			}
		})
	}
}

func TestWorkflowStatusFiltersAndAllReachReader(t *testing.T) {
	for _, flags := range [][]string{{}, {"--project", "p"}, {"--repo", "r"}, {"--epic", "e"}, {"--project", "p", "--repo", "r", "--epic", "e"}} {
		var got workspaceview.WorkflowStatusOptions
		root, _, _ := workflowStatusTestRoot(t, newWorkflowStatusCommandWithReader(func(o workspaceview.WorkflowStatusOptions) (workspaceview.WorkflowStatus, error) {
			got = o
			return emptyWorkflowStatus(), nil
		}))
		root.SetArgs(append([]string{"workflow", "status", "--all", "--json"}, flags...))
		if err := root.Execute(); err != nil || !got.IncludeAll {
			t.Fatalf("err=%v options=%+v", err, got)
		}
		for i := 0; i < len(flags); i += 2 {
			var value string
			switch flags[i] {
			case "--project":
				if got.Filters.ProjectID != nil {
					value = string(*got.Filters.ProjectID)
				}
			case "--repo":
				if got.Filters.RepoID != nil {
					value = string(*got.Filters.RepoID)
				}
			case "--epic":
				if got.Filters.EpicID != nil {
					value = string(*got.Filters.EpicID)
				}
			}
			if value != flags[i+1] {
				t.Fatalf("lost filter %s: %+v", flags[i], got)
			}
		}
	}
}

func TestWorkflowStatusErrorsAndLogsStayOnStderr(t *testing.T) {
	for _, readError := range []error{nil, errors.New("source registry is unreadable")} {
		var legacyLog bytes.Buffer
		previous := logrus.StandardLogger().Out
		logrus.SetOutput(&legacyLog)
		root, out, stderr := workflowStatusTestRoot(t, newWorkflowStatusCommandWithReader(func(workspaceview.WorkflowStatusOptions) (workspaceview.WorkflowStatus, error) {
			logrus.Warn("status source diagnostic")
			return emptyWorkflowStatus(), readError
		}))
		root.SetArgs([]string{"--debug", "--json", "workflow", "status"})
		originalRoot := RootCmd
		RootCmd = root
		err := ExecuteE()
		RootCmd = originalRoot
		restored := logrus.StandardLogger().Out == &legacyLog
		logrus.SetOutput(previous)
		if !errors.Is(err, readError) || !restored || legacyLog.Len() != 0 || !strings.Contains(stderr.String(), "status source diagnostic") {
			t.Fatalf("err=%v restored=%v log=%s stderr=%s", err, restored, &legacyLog, stderr)
		}
		if readError != nil {
			if out.Len() != 0 || !strings.Contains(stderr.String(), readError.Error()) {
				t.Fatalf("stdout=%s stderr=%s", out, stderr)
			}
		} else if !json.Valid(out.Bytes()) {
			t.Fatalf("contaminated JSON stdout: %s", out)
		}
	}
}

func TestWorkflowStatusActualReadEmptyAndMissingWorkspace(t *testing.T) {
	d, path := viewFixture(t)
	for _, args := range [][]string{{"workflow", "status", "--json"}, {"--json", "workflow", "status", "--all"}} {
		root, out, stderr := workflowStatusTestRoot(t, newWorkflowStatusCommand(d))
		root.SetArgs(args)
		if err := root.Execute(); err != nil || stderr.Len() != 0 {
			t.Fatalf("err=%v stderr=%s", err, stderr)
		}
		var value map[string]any
		if err := json.Unmarshal(out.Bytes(), &value); err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"projects", "epics", "items", "diagnostics"} {
			if rows, ok := value[key].([]any); !ok || len(rows) != 0 {
				t.Fatalf("%s should be []: %s", key, out)
			}
		}
	}
	entries, err := os.ReadDir(filepath.Join(path, ".ply"))
	if err != nil || len(entries) != 1 || entries[0].Name() != "workspace.yaml" {
		t.Fatalf("read created state: %v, %v", entries, err)
	}
	outside := t.TempDir()
	d.Files = viewCWD{d.Files, outside}
	root, out, _ := workflowStatusTestRoot(t, newWorkflowStatusCommand(d))
	root.SetArgs([]string{"workflow", "status", "--json"})
	if err := root.Execute(); err == nil || out.Len() != 0 {
		t.Fatalf("missing workspace: err=%v stdout=%s", err, out)
	}
	if entries, err := os.ReadDir(outside); err != nil || len(entries) != 0 {
		t.Fatalf("missing workspace initialized: %v, %v", entries, err)
	}
}

func TestWorkflowStatusTextRetainsActorsHistoryAndUnknownTimes(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "schemas", "read-contract", "workflow-status-examples", "populated.valid.json"))
	if err != nil {
		t.Fatal(err)
	}
	var value workspaceview.WorkflowStatus
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatal(err)
	}
	text := workflowStatusText(value)
	for _, expected := range []string{"Needs you (1)", "Follow-up (1)", "In progress (recorded) (1)", "Ready next (registered) (1)", "Human: Which output should be the default?", "Agent:", "time unknown", "live state unknown", "start preflight required", "Hidden: 8 backlog · 3 inactive · 2 completed. Use --all to inspect."} {
		if !strings.Contains(text, expected) {
			t.Fatalf("missing %q:\n%s", expected, text)
		}
	}
	value.Items[0].Category = "inactive"
	value.Items[0].TaskLifecycle = "parked"
	value.Items[0].NextActions[0].Current = false
	value.Counts.Categories["needs_you"] = 0
	value.Counts.Categories["inactive"] = 1
	text = workflowStatusText(value)
	if !strings.Contains(text, "historical") || !strings.Contains(text, "Task lifecycle: parked") || strings.Contains(text, "Needs you (") {
		t.Fatalf("historical need became current:\n%s", text)
	}
}

func TestWorkflowStatusStdoutErrorsPropagate(t *testing.T) {
	for _, format := range []string{"text", "json"} {
		root, _, _ := workflowStatusTestRoot(t, newWorkflowStatusCommandWithReader(func(workspaceview.WorkflowStatusOptions) (workspaceview.WorkflowStatus, error) {
			return emptyWorkflowStatus(), nil
		}))
		root.SetOut(coreFailWriter{})
		root.SetArgs([]string{"workflow", "status", "--format", format})
		if err := root.Execute(); err == nil {
			t.Fatal("ignored output failure")
		}
	}
}

func TestWorkflowStatusTextIdentifiesOrphanDiagnostics(t *testing.T) {
	value := emptyWorkflowStatus()
	queue, run, task := "orphan-queue-a", "orphan-run-b", "damaged-task"
	value.Diagnostics = []workspaceview.WorkflowDiagnostic{
		{Code: "task_queue_orphaned", Source: "queue", Actor: "unknown", Severity: "attention", Reason: "Queue identity does not match a registered target.", QueueID: &queue, EvidenceIDs: []string{"queue-source-digest"}},
		{Code: "run_evidence_unavailable", Source: "run", Actor: "unknown", Severity: "attention", Reason: "Run source is unreadable.", RunID: &run, TaskID: &task, EvidenceIDs: []string{"run-source-digest"}},
	}
	text := workflowStatusText(value)
	for _, identity := range []string{queue, run, task, "queue-source-digest", "run-source-digest"} {
		if !strings.Contains(text, identity) {
			t.Fatalf("diagnostic cannot be identified: missing %q in\n%s", identity, text)
		}
	}
}
