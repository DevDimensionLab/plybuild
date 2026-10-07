package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode/utf16"

	"github.com/devdimensionlab/plybuild/internal/workflowtrace"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newWorkflowTraceCommand(d workspace.Dependencies) *cobra.Command {
	return newWorkflowTraceCommandWithReader(func(taskID string) (workflowtrace.Result, error) {
		return workflowtrace.Read(d, taskID)
	})
}

func newWorkflowTraceCommandWithReader(read func(string) (workflowtrace.Result, error)) *cobra.Command {
	var format string
	var details bool
	invalid := func(reason string) error { return fmt.Errorf("workflow_trace_invalid_arguments: %s", reason) }
	command := &cobra.Command{
		Use:               "trace <task-id>",
		Short:             "Trace one Task's recorded process across its runs",
		Long:              "Draw one registered Task's recorded workflow as ASCII step boxes with Human/Agent/Ply/Unknown actor labels. Arrows within a source chain mean recorded order, not causality; independent chains stay disconnected. Every event is shown. Long text excerpts are marked and counted. Current and frozen declarations remain separate from recorded steps.\n\n--details appends full evidence and a key from local diagram labels to exact IDs, hashes, declarations and payloads. --format json is always the unchanged full WorkflowTraceReadback@1 object and newline; --details has no effect on JSON. --json is a local result shortcut, including before workflow. --json=false selects no format; --json with explicit --format text is an error.\n\nThis read-only projection performs no Git, provider, network or process probes, initialization, locks or workflow transitions. It grants no start, retry, QA or integration authority. Fresh sources do not establish live activity or active effort. Fatal errors write stderr with empty stdout. See docs/workflow-trace.md for summary, ordering, evidence and coverage semantics.",
		Example:           "  ply workflow trace task-id\n  ply workflow trace task-id --details\n  ply workflow trace task-id --format json\n  ply --json workflow trace task-id",
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 || strings.TrimSpace(args[0]) == "" {
				return invalid("expected one registered Task ID")
			}
			if _, err := workspace.ParseTaskID(args[0]); err != nil {
				return invalid(err.Error())
			}
			if format != "text" && format != "json" {
				return invalid("--format must be text or json")
			}
			if workflowStatusJSON(cmd) && cmd.Flags().Changed("format") && format == "text" {
				return invalid("--json conflicts with explicit --format text")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			previousOutput := logrus.StandardLogger().Out
			logrus.SetOutput(cmd.ErrOrStderr())
			defer logrus.SetOutput(previousOutput)
			result, err := read(args[0])
			if err != nil {
				return err
			}
			if format == "json" || workflowStatusJSON(cmd) {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
			}
			if details {
				_, err = fmt.Fprint(cmd.OutOrStdout(), workflowTraceDetails(result))
				return err
			}
			_, err = fmt.Fprint(cmd.OutOrStdout(), workflowTraceText(result))
			return err
		},
	}
	command.Flags().StringVar(&format, "format", "text", "output format (text or json); --json is a local result shortcut")
	command.Flags().BoolVar(&details, "details", false, "append full evidence and diagram label key; JSON is always full")
	command.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return invalid(err.Error()) })
	return command
}

// Text is an unabridged presentation of the same projection. Quoted source
// strings preserve whitespace and escape controls/non-ASCII for a safe ASCII
// terminal view. JSON payloads are retained, not reduced to selected fields.
func workflowTraceEvidenceText(result workflowtrace.Result) string {
	var out strings.Builder
	line := func(label string, value any) { fmt.Fprintf(&out, "    %s: %s\n", label, workflowTraceJSON(value)) }
	fmt.Fprintf(&out, "Workflow trace - %s: %s\n", workflowTraceJSON(result.Task.ID), workflowTraceJSON(result.Task.Title))
	fmt.Fprintf(&out, "Contract: %s / schema %d\n", result.Kind, result.SchemaVersion)
	line("workspace", result.Workspace)
	line("Task", result.Task)
	fmt.Fprintf(&out, "Observed: %s | source freshness: %s | live state: %s | history: %s\n", workflowTraceJSON(result.ObservedAtUTC), workflowTraceJSON(result.Freshness), workflowTraceJSON(result.LiveState), workflowTraceJSON(result.HistoryState))
	line("ordering", result.Ordering)
	line("next_transition_authorized", result.NextTransitionAuthorized)
	fmt.Fprintln(&out, "Independent source chains are not a global chronology. Declarations are not executed steps.")
	if result.HistoryState == "empty" {
		fmt.Fprintln(&out, "No execution recorded. Registered Task declarations and state may still be present.")
	}

	fmt.Fprintf(&out, "\nAnalysis summary (%d) - definitions and evidence below\n", len(result.Analysis))
	fmt.Fprintln(&out, "Coverage: exact covers the defined selected records; lower_bound means at least these records; unknown means insufficient evidence. None proves complete real work history.")
	fmt.Fprintln(&out, "Elapsed time: reported QA/wait intervals may include waiting; verifier elapsed is not active agent or human effort. No time savings are measured.")
	for _, analysis := range result.Analysis {
		fmt.Fprintf(&out, "  %s: %s %s | coverage %s", workflowTraceJSON(analysis.ID), workflowTraceJSON(analysis.Value), workflowTraceJSON(analysis.Unit), workflowTraceJSON(analysis.Coverage))
		if analysis.Question != nil {
			fmt.Fprintf(&out, " | %s", workflowTraceJSON(analysis.Question))
		}
		fmt.Fprintln(&out)
	}

	fmt.Fprintf(&out, "\nSource chains (%d) - role map; positions are not elapsed time\n", len(result.Chains))
	eventIndex := map[string]workflowtrace.Event{}
	for _, event := range result.Events {
		eventIndex[event.ID] = event
	}
	for _, chain := range result.Chains {
		fmt.Fprintf(&out, "  Chain %s | ordering %s\n", workflowTraceJSON(chain.ID), workflowTraceJSON(chain.Ordering))
		for _, id := range chain.EventIDs {
			event, found := eventIndex[id]
			if !found {
				fmt.Fprintf(&out, "    [Unknown] %s | unresolved event reference\n", workflowTraceJSON(id))
				continue
			}
			positions := []int{}
			for _, position := range event.Positions {
				if position.ChainID == chain.ID {
					positions = append(positions, position.Sequence)
				}
			}
			fmt.Fprintf(&out, "    %s | [%s] %s | %s | %s\n", workflowTraceJSON(positions), workflowTraceRole(event.Role), workflowTraceJSON(event.Title), workflowTraceJSON(event.Type), workflowTraceJSON(event.ID))
		}
	}
	for _, event := range result.Events {
		if len(event.Positions) == 0 {
			fmt.Fprintf(&out, "  No source position | [%s] %s | %s | %s\n", workflowTraceRole(event.Role), workflowTraceJSON(event.Title), workflowTraceJSON(event.Type), workflowTraceJSON(event.ID))
		}
	}

	fmt.Fprintf(&out, "\nCurrent goal declarations (%d)\n", len(result.CurrentGoals))
	for _, goal := range result.CurrentGoals {
		fmt.Fprintf(&out, "  Goal %s revision %d | %s\n", workflowTraceJSON(goal.SpecID), goal.Revision, workflowTraceJSON(goal.Status))
		line("manifest_sha256", goal.ManifestSHA256)
		line("source_ids", goal.SourceIDs)
		line("declaration", goal.Declaration)
	}
	fmt.Fprintf(&out, "\nFrozen run contracts (%d)\n", len(result.Runs))
	if len(result.Runs) == 0 {
		fmt.Fprintln(&out, "  Declared process unknown: no frozen run contract is recorded.")
	}
	for _, run := range result.Runs {
		fmt.Fprintf(&out, "  Run %s | family %s | coverage %s\n", workflowTraceJSON(run.ID), workflowTraceJSON(run.Family), workflowTraceJSON(run.Coverage))
		line("provider", run.Provider)
		line("session_id", run.SessionID)
		line("request_sha256", run.RequestSHA256)
		line("frozen_basis", run.FrozenBasis)
		line("frozen_goal", run.FrozenGoal)
		line("declared_process", run.DeclaredProcess)
		if !workflowTraceHasDeclaration(run.DeclaredProcess) {
			fmt.Fprintln(&out, "    Declared process unknown: no preserved process declaration is available for this run.")
		}
		line("source_ids", run.SourceIDs)
		line("event_ids", run.EventIDs)
	}
	fmt.Fprintf(&out, "\nRecorded history (%d) - lanes: Human | Agent | Ply | Unknown\n", len(result.Events))
	if len(result.Events) == 0 {
		fmt.Fprintln(&out, "  No safely bound recorded events. This does not prove no unrecorded work occurred.")
	}
	for _, event := range result.Events {
		fmt.Fprintf(&out, "  [%s] %s | %s | %s\n", workflowTraceRole(event.Role), workflowTraceJSON(event.ID), workflowTraceJSON(event.Type), workflowTraceJSON(event.Title))
		line("native_id", event.NativeID)
		line("actor_claim", event.ActorClaim)
		line("recorder", event.Recorder)
		line("evidence_class", event.EvidenceClass)
		line("run_ids", event.RunIDs)
		line("candidate_id", event.CandidateID)
		line("result_id", event.ResultID)
		line("outcome", event.Outcome)
		line("occurred_at_utc", event.OccurredAtUTC)
		line("reported_at_utc", event.ReportedAtUTC)
		line("registered_at_utc", event.RegisteredAtUTC)
		line("time_basis", event.TimeBasis)
		line("source_ids", event.SourceIDs)
		line("positions", event.Positions)
		line("data", event.Data)
	}
	fmt.Fprintf(&out, "\nSource relations (%d)\n", len(result.Relations))
	for _, relation := range result.Relations {
		fmt.Fprintf(&out, "  %s --%s--> %s\n", workflowTraceJSON(relation.From), workflowTraceJSON(relation.Type), workflowTraceJSON(relation.To))
		line("basis", relation.Basis)
		line("source_ids", relation.SourceIDs)
	}
	fmt.Fprintf(&out, "\nCandidate facts (%d) - reported, technical, human QA and integration remain separate\n", len(result.Candidates))
	for _, candidate := range result.Candidates {
		fmt.Fprintf(&out, "  Candidate %s | oid %s | tree %s\n", workflowTraceJSON(candidate.ID), workflowTraceJSON(candidate.OID), workflowTraceJSON(candidate.Tree))
		line("run_ids", candidate.RunIDs)
		line("result_ids", candidate.ResultIDs)
		line("event_ids", candidate.EventIDs)
		keys := make([]string, 0, len(candidate.Axes))
		for key := range candidate.Axes {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			line("axis "+workflowTraceJSON(key), candidate.Axes[key])
		}
	}
	fmt.Fprintf(&out, "\nEvidence-linked analysis (%d)\n", len(result.Analysis))
	for _, analysis := range result.Analysis {
		fmt.Fprintf(&out, "  %s | %s | coverage %s\n", workflowTraceJSON(analysis.ID), workflowTraceJSON(analysis.Kind), workflowTraceJSON(analysis.Coverage))
		line("definition", analysis.Definition)
		line("value", analysis.Value)
		line("unit", analysis.Unit)
		line("evidence_ids", analysis.EvidenceIDs)
		line("source_ids", analysis.SourceIDs)
		line("question", analysis.Question)
		line("unknowns", analysis.Unknowns)
	}

	fmt.Fprintf(&out, "\nCoverage and unknowns - %s\n", workflowTraceJSON(result.Coverage.State))
	line("read", result.Coverage.Read)
	line("unknowns", result.Coverage.Unknowns)
	fmt.Fprintf(&out, "\nSources (%d)\n", len(result.Sources))
	for _, source := range result.Sources {
		fmt.Fprintf(&out, "  %s | %s | %s\n", workflowTraceJSON(source.ID), workflowTraceJSON(source.Kind), workflowTraceJSON(source.Status))
		line("locator", source.Locator)
		line("sha256", source.SHA256)
	}
	fmt.Fprintf(&out, "\nDiagnostics (%d)\n", len(result.Diagnostics))
	for _, diagnostic := range result.Diagnostics {
		fmt.Fprintf(&out, "  %s: %s\n", workflowTraceJSON(diagnostic.Code), workflowTraceJSON(diagnostic.Detail))
		line("source_ids", diagnostic.SourceIDs)
		line("event_ids", diagnostic.EventIDs)
	}
	return out.String()
}

func workflowTraceRole(role string) string {
	switch role {
	case "human":
		return "Human"
	case "agent":
		return "Agent"
	case "ply":
		return "Ply"
	default:
		return "Unknown"
	}
}

func workflowTraceHasDeclaration(raw json.RawMessage) bool {
	var value any
	if json.Unmarshal(raw, &value) != nil {
		return false
	}
	var present func(any) bool
	present = func(v any) bool {
		switch v := v.(type) {
		case nil:
			return false
		case map[string]any:
			for _, item := range v {
				if present(item) {
					return true
				}
			}
			return false
		case []any:
			for _, item := range v {
				if present(item) {
					return true
				}
			}
			return false
		case string:
			return strings.TrimSpace(v) != ""
		default:
			return true
		}
	}
	return present(value)
}

func workflowTraceJSON(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		// Production values are validated JSON. Keep any broken payload explicit
		// if this renderer is reused, rather than silently dropping the field.
		return fmt.Sprintf("<invalid JSON: %s>", err)
	}
	var out strings.Builder
	for _, r := range string(data) {
		if r < 128 {
			out.WriteRune(r)
		} else if r <= 0xffff {
			fmt.Fprintf(&out, "\\u%04x", r)
		} else {
			hi, lo := utf16.EncodeRune(r)
			fmt.Fprintf(&out, "\\u%04x\\u%04x", hi, lo)
		}
	}
	return out.String()
}
