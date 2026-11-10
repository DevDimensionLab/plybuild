package cmd

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/humaninput"
)

// HI07: every advertised cooperative_agent_v1 operation is an actual command
// whose flags are exactly the advertised selectors (plus --format).
func TestHI07InputCapabilitiesMatchActualCommands(t *testing.T) {
	want := map[string]struct {
		effect  string
		schemas []string
	}{
		"workflow.input.ask":            {"human_input_question_publication", []string{humaninput.QuestionKind, humaninput.ErrorKind}},
		"workflow.input.list":           {"read", []string{humaninput.ListKind, humaninput.ErrorKind}},
		"workflow.input.show":           {"read", []string{humaninput.QuestionKind, humaninput.ErrorKind}},
		"workflow.input.answer.preview": {"read", []string{humaninput.PreviewKind, humaninput.RequestKind, humaninput.ErrorKind}},
		"workflow.input.answer.apply":   {"human_input_answer_record", []string{humaninput.StatusKind, humaninput.ErrorKind}},
		"workflow.input.answer.status":  {"read", []string{humaninput.StatusKind}},
		"workflow.input.wait":           {"read", []string{humaninput.WaitKind, humaninput.ErrorKind}},
		"workflow.input.acknowledge":    {"agent_reported_receipt", []string{humaninput.QuestionKind, humaninput.ErrorKind}},
		"workflow.input.continue":       {"agent_reported_continuation", []string{humaninput.QuestionKind, humaninput.ErrorKind}},
		"workflow.input.withdraw":       {"human_input_question_withdrawal", []string{humaninput.QuestionKind, humaninput.ErrorKind}},
	}
	seen := map[string]int{}
	for _, op := range coreCapabilities(capabilityBuild{}).WorkflowExtensions {
		w, ok := want[op.ID]
		if !ok {
			continue
		}
		seen[op.ID]++
		if op.Mode != humaninput.Mode || op.Effect != w.effect || len(op.ResultSchemas) != len(w.schemas) {
			t.Fatalf("capability %s: %#v", op.ID, op)
		}
		for i, s := range op.ResultSchemas {
			if s.Kind != w.schemas[i] || s.SchemaVersion == nil || *s.SchemaVersion != humaninput.SchemaVersion {
				t.Fatalf("capability %s schema %d: %#v", op.ID, i, s)
			}
		}
		if op.Command[0] != "workflow" || op.Command[1] != "input" || strings.Join(op.Command[2:], ".") != strings.TrimPrefix(op.ID, "workflow.input.") {
			t.Fatalf("capability %s command %v", op.ID, op.Command)
		}
		found, rest, err := newWorkflowInputCommand(humaninput.Dependencies{}).Find(op.Command[2:])
		if err != nil || len(rest) != 0 || found.Name() != op.Command[len(op.Command)-1] {
			t.Fatalf("advertised command %v is absent", op.Command)
		}
		flags := map[string]bool{}
		for _, s := range op.Selectors {
			flags[strings.TrimPrefix(s, "--")] = true
			if found.Flags().Lookup(strings.TrimPrefix(s, "--")) == nil {
				t.Fatalf("%s advertises %s which is not a flag", op.ID, s)
			}
		}
		actual := regexp.MustCompile(`(?m)^\s+(?:-\w, )?--([a-z0-9-]+)`).FindAllStringSubmatch(found.LocalFlags().FlagUsages(), -1)
		if len(actual) != len(op.Selectors)+1 {
			t.Fatalf("%s has flags %v; advertised %v plus --format", op.ID, actual, op.Selectors)
		}
		for _, m := range actual {
			if m[1] != "format" && !flags[m[1]] {
				t.Fatalf("%s has unadvertised flag --%s", op.ID, m[1])
			}
		}
		var help bytes.Buffer
		found.SetOut(&help)
		if err := found.Help(); err != nil || !strings.Contains(help.String(), "--format") {
			t.Fatalf("help for %s: %v", op.ID, err)
		}
	}
	for id := range want {
		if seen[id] != 1 {
			t.Fatalf("capability %s appears %d times", id, seen[id])
		}
	}

}
