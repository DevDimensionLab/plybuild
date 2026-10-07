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

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workflowtrace"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func workflowTraceTestRoot(t *testing.T, trace *cobra.Command) (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	root := &cobra.Command{Use: "ply", PersistentPreRunE: func(*cobra.Command, []string) error {
		t.Error("workflow trace invoked profile initialization")
		return errors.New("profile initialization must not run")
	}}
	root.PersistentFlags().Bool("json", false, "logging")
	root.PersistentFlags().Bool("debug", false, "debug logging")
	workflow := &cobra.Command{Use: "workflow"}
	workflow.AddCommand(trace)
	root.AddCommand(workflow)
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SilenceErrors, root.SilenceUsage = true, true
	return root, &stdout, &stderr
}

func workflowTraceFixture(t *testing.T, name string) workflowtrace.Result {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "schemas", "read-contract", "workflow-trace-examples", name+".valid.json"))
	if err != nil {
		t.Fatal(err)
	}
	var value workflowtrace.Result
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatal(err)
	}
	return value
}

func TestWorkflowTraceJSONFlagsOneObjectAndSameFacts(t *testing.T) {
	want := workflowTraceFixture(t, "multi-family")
	for _, args := range [][]string{
		{"workflow", "trace", "trace-example", "--format", "json"},
		{"workflow", "trace", "trace-example", "--json"},
		{"--json", "workflow", "trace", "trace-example"},
		{"workflow", "--json", "trace", "trace-example"},
		{"--json", "workflow", "trace", "trace-example", "--format", "json"},
		{"workflow", "trace", "trace-example", "--format", "json", "--json"},
		{"workflow", "trace", "trace-example", "--json=false", "--format", "json"},
		{"--debug", "--json", "workflow", "trace", "trace-example"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var calls int
			root, out, stderr := workflowTraceTestRoot(t, newWorkflowTraceCommandWithReader(func(taskID string) (workflowtrace.Result, error) {
				if taskID != "trace-example" {
					t.Fatalf("wrong Task selected: %q", taskID)
				}
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
			var got workflowtrace.Result
			decoder := json.NewDecoder(out)
			if err := decoder.Decode(&got); err != nil {
				t.Fatal(err)
			}
			// Compare semantic JSON: RawMessage preserves formatting from its
			// input, while the encoder compacts payload whitespace.
			gotJSON, _ := json.Marshal(got)
			wantJSON, _ := json.Marshal(want)
			if !bytes.Equal(gotJSON, wantJSON) {
				t.Fatalf("JSON changed recorded facts:\n%s\n%s", gotJSON, wantJSON)
			}
			var extra any
			if err := decoder.Decode(&extra); err != io.EOF {
				t.Fatalf("extra stdout after result: %v", err)
			}
		})
	}
	for _, flags := range [][]string{{}, {"--json=false"}, {"--json=false", "--format", "text"}} {
		root, out, _ := workflowTraceTestRoot(t, newWorkflowTraceCommandWithReader(func(string) (workflowtrace.Result, error) { return want, nil }))
		root.SetArgs(append([]string{"workflow", "trace", "trace-example"}, flags...))
		if err := root.Execute(); err != nil || out.String() != workflowTraceText(want) {
			t.Fatalf("flags=%v err=%v stdout=%s", flags, err, out)
		}
	}
}

func TestWorkflowTraceRejectsInvalidArgumentsBeforeReading(t *testing.T) {
	for _, args := range [][]string{
		{"workflow", "trace"}, {"workflow", "trace", "one", "two"},
		{"workflow", "trace", ""}, {"workflow", "trace", "Bad"},
		{"workflow", "trace", "trace-example", "--format", "yaml"},
		{"workflow", "trace", "trace-example", "--format"},
		{"workflow", "trace", "trace-example", "--unknown"},
		{"--json", "workflow", "trace", "trace-example", "--format", "text"},
		{"workflow", "trace", "trace-example", "--format", "text", "--json"},
		{"workflow", "trace", "trace-example", "--json", "--format", "text"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			calls := 0
			root, out, _ := workflowTraceTestRoot(t, newWorkflowTraceCommandWithReader(func(string) (workflowtrace.Result, error) {
				calls++
				return workflowtrace.Result{}, nil
			}))
			root.SetArgs(args)
			if err := root.Execute(); err == nil || !strings.HasPrefix(err.Error(), "workflow_trace_invalid_arguments:") || calls != 0 || out.Len() != 0 {
				t.Fatalf("err=%v calls=%d stdout=%s", err, calls, out)
			}
		})
	}
}

func TestWorkflowTraceErrorsAndLogsStayOnStderr(t *testing.T) {
	want := workflowTraceFixture(t, "partial")
	for _, readError := range []error{nil, errors.New("trace source registry is unreadable")} {
		var legacyLog bytes.Buffer
		previous := logrus.StandardLogger().Out
		logrus.SetOutput(&legacyLog)
		root, out, stderr := workflowTraceTestRoot(t, newWorkflowTraceCommandWithReader(func(string) (workflowtrace.Result, error) {
			logrus.Warn("trace source diagnostic")
			return want, readError
		}))
		root.SetArgs([]string{"--debug", "--json", "workflow", "trace", "trace-example"})
		originalRoot := RootCmd
		RootCmd = root
		err := ExecuteE()
		RootCmd = originalRoot
		restored := logrus.StandardLogger().Out == &legacyLog
		logrus.SetOutput(previous)
		if !errors.Is(err, readError) || !restored || legacyLog.Len() != 0 || !strings.Contains(stderr.String(), "trace source diagnostic") {
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

func TestWorkflowTraceActualUnknownTaskAndMissingWorkspaceDoNotInitialize(t *testing.T) {
	d, path := viewFixture(t)
	root, out, _ := workflowTraceTestRoot(t, newWorkflowTraceCommand(d))
	root.SetArgs([]string{"workflow", "trace", "not-registered", "--json"})
	if err := root.Execute(); err == nil || out.Len() != 0 {
		t.Fatalf("unknown Task err=%v stdout=%s", err, out)
	}
	entries, err := os.ReadDir(filepath.Join(path, ".ply"))
	if err != nil || len(entries) != 1 || entries[0].Name() != "workspace.yaml" {
		t.Fatalf("read created state: %v, %v", entries, err)
	}
	outside := t.TempDir()
	d.Files = viewCWD{d.Files, outside}
	root, out, _ = workflowTraceTestRoot(t, newWorkflowTraceCommand(d))
	root.SetArgs([]string{"workflow", "trace", "not-registered", "--json"})
	if err := root.Execute(); err == nil || out.Len() != 0 {
		t.Fatalf("missing workspace: err=%v stdout=%s", err, out)
	}
	if entries, err := os.ReadDir(outside); err != nil || len(entries) != 0 {
		t.Fatalf("missing workspace initialized: %v, %v", entries, err)
	}
}

func TestWorkflowTraceTextPreservesQuestionsEvidenceAndUncertainty(t *testing.T) {
	value := workflowTraceFixture(t, "multi-family")
	text := workflowTraceText(value)
	for _, expected := range []string{
		"Current goal declarations", "Frozen run contracts", "Source chains",
		"Human | Agent | Ply | Unknown", "[Human]", "[Agent]", "[Ply]", "[Unknown]",
		"revision 2", "\"revision\":1", "which output", // checked case-insensitively below
		`Which output should be the default?\nKeep caf\u00e9 and \ud83e\udded unchanged.`,
		`"reason":"Waiting for fixture answer"`, `"dependency":"design-question"`, `"next_actor":"human"`,
		`"correction" --"corrects"--> "verify-c1-fail"`,
		`"verify-c3" --"handoff"--> "qa-c3-pass"`,
		`"qa-c2-fail"`, `"qa-c3-pass"`, `"result-c2"`, `"result-c3"`,
		`occurred_at_utc: null`, `reported_at_utc: null`, `registered_at_utc: null`,
		`"clock":null`, `value: null`, `"inspect-loop"`,
		"Could an earlier product check expose the missing question text?",
		"No effort or time savings was measured.", "Independent source chains are not a global chronology.",
	} {
		if !strings.Contains(strings.ToLower(text), strings.ToLower(expected)) {
			t.Errorf("text omitted %q", expected)
		}
	}
	for _, r := range text {
		if r > 127 || (r < 32 && r != '\n') {
			t.Fatalf("text is not safe ASCII: %U", r)
		}
	}
	// Every evidence identity and safe payload remains inspectable in text.
	for _, e := range value.Events {
		if !strings.Contains(text, workflowTraceJSON(e.ID)) || !strings.Contains(text, workflowTraceJSON(e.Data)) {
			t.Fatalf("event or full payload omitted: %s", e.ID)
		}
	}
	for _, s := range value.Sources {
		if !strings.Contains(text, workflowTraceJSON(s.ID)) || !strings.Contains(text, workflowTraceJSON(s.Locator)) {
			t.Fatalf("source omitted: %+v", s)
		}
	}
	partial := workflowTraceText(workflowTraceFixture(t, "partial"))
	for _, expected := range []string{"lower_bound", "trace_source_missing", "missing-history", "independent valid history is retained"} {
		if !strings.Contains(partial, expected) {
			t.Fatalf("partial text omitted %q", expected)
		}
	}
}

func TestWorkflowTraceTextDoesNotTruncateHistory(t *testing.T) {
	value := workflowTraceFixture(t, "multi-family")
	original := value.Events[0]
	value.Events = nil
	for i := 0; i < 250; i++ {
		event := original
		event.ID = strings.Repeat("x", i+1) + "-last-marker"
		value.Events = append(value.Events, event)
	}
	text := workflowTraceText(value)
	if strings.Count(text, "last-marker") != 250 || !strings.Contains(text, "Recorded history (250)") {
		t.Fatal("history was silently truncated")
	}
	if workflowTraceText(value) != text {
		t.Fatal("text presentation is unstable")
	}
}

func TestWorkflowTraceEmptyExecutionAndAbsentProcessStayExplicit(t *testing.T) {
	empty := workflowTraceFixture(t, "empty")
	// A native Task-state event does not mean a run was executed.
	empty.Events = workflowTraceFixture(t, "multi-family").Events[:1]
	text := workflowTraceText(empty)
	for _, expected := range []string{"No execution recorded", "Declared process unknown", "legacy-report"} {
		if !strings.Contains(text, expected) {
			t.Fatalf("empty execution omitted %q", expected)
		}
	}
	for _, raw := range []string{"null", "{}", `{"contract":null,"obligations":[]}`, `{"contract":{},"obligations":null}`} {
		value := workflowTraceFixture(t, "multi-family")
		value.Runs = value.Runs[2:]
		value.Runs[0].DeclaredProcess = json.RawMessage(raw)
		if !strings.Contains(workflowTraceText(value), "Declared process unknown") {
			t.Fatalf("absent process became known: %s", raw)
		}
	}
}

func TestWorkflowTraceSummaryExplainsCoverageAndElapsedTimeBeforeChains(t *testing.T) {
	value := workflowTraceFixture(t, "multi-family")
	elapsed := float64(22367)
	value.Analysis = append(value.Analysis, workflowtrace.Analysis{
		ID: "reported_qa/example", Kind: "interval", Definition: "Reported human QA attestation interval.",
		Value: &elapsed, Unit: "seconds", Coverage: "exact",
	})
	text := workflowTraceText(value)
	index := strings.Index(text, "Source chains (")
	if index < 0 {
		t.Fatal("source chains heading missing")
	}
	summary := text[:index]
	for _, expected := range []string{
		`22367 "seconds" | coverage "exact"`,
		"exact covers the defined selected records",
		"lower_bound means at least these records",
		"unknown means insufficient evidence",
		"None proves complete real work history",
		"reported QA/wait intervals may include waiting",
		"verifier elapsed is not active agent or human effort",
		"No time savings are measured",
	} {
		if !strings.Contains(summary, expected) {
			t.Errorf("leading summary omitted %q", expected)
		}
	}
}

func TestWorkflowTraceHelpAndCapabilitiesRemainAdditive(t *testing.T) {
	tree := &cobra.Command{Use: "ply"}
	tree.AddCommand(newWorkflowCommand(workflowhandoff.Dependencies{}))
	trace, rest, err := tree.Find([]string{"workflow", "trace"})
	if err != nil || len(rest) != 0 || trace.Use != "trace <task-id>" || trace.Flags().Lookup("format") == nil {
		t.Fatalf("trace not discoverable: %v %v", rest, err)
	}
	workflow, _, _ := tree.Find([]string{"workflow"})
	if !strings.Contains(workflow.Long, "trace one Task") || !strings.Contains(workflow.Example, "workflow trace") || !strings.Contains(RootCmd.Example, "workflow trace") {
		t.Fatal("workflow/root help omitted trace")
	}
	catalog := coreCapabilities(capabilityBuild{})
	if len(catalog.Operations) != 5 || catalog.Coverage != "workspace-core-read" {
		t.Fatal("changed legacy capability core")
	}
	for _, op := range catalog.ReadExtensions {
		if op.ID != "workflow.trace" {
			continue
		}
		if !reflect.DeepEqual(op.Command, []string{"workflow", "trace"}) || op.Effect != "read" || !reflect.DeepEqual(op.Formats, []string{"text", "json"}) || len(op.ResultSchemas) != 1 || op.ResultSchemas[0].Kind != workflowtrace.Kind || op.ResultSchemas[0].SchemaVersion == nil || *op.ResultSchemas[0].SchemaVersion != 1 {
			t.Fatalf("wrong trace capability: %+v", op)
		}
		return
	}
	t.Fatal("trace capability missing")
}
