package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/agentwork"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/spf13/cobra"
	"io"
	"os"
)

type agentWorkFactory func(string) (*agentwork.Service, error)

func agentWorkOutput(cmd *cobra.Command, format string, v any) error {
	if format == "json" || workflowStatusJSON(cmd) {
		return json.NewEncoder(cmd.OutOrStdout()).Encode(v)
	}
	switch x := v.(type) {
	case agentwork.MetadataReadback:
		fmt.Fprintf(cmd.OutOrStdout(), "Agent metadata · read %s · source %s · coverage %s\n", x.ObservedAtUTC, x.Sources[0].Status, x.Sources[0].Coverage)
		for _, a := range x.Agents {
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s · %s · tab %s · workspace %s\n", a.Pane, a.Provider, a.Status, a.Tab, a.Workspace)
		}
	case agentwork.Overview:
		fmt.Fprintf(cmd.OutOrStdout(), "Agent work · read %s\n", x.ObservedAtUTC)
		for _, s := range x.Sources {
			fmt.Fprintf(cmd.OutOrStdout(), "Source %s: %s · freshness %s\n", s.Source, s.Status, s.Freshness)
		}
		for _, a := range x.Agents {
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s · %s · pane %s · binding %s\n", a.AgentRef, a.Provider, a.Status, a.Pane, a.BindingStatus)
			if a.Work != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "  Task %s · Run %s · registered %s · next %s: %s\n", a.Work.TaskID, a.Work.RunID, a.Work.Phase, a.Work.NextActor, a.Work.NextAction)
			}
		}
		for _, f := range x.Followups {
			printFollowup(cmd, f)
		}
	case agentwork.List:
		fmt.Fprintf(cmd.OutOrStdout(), "Follow-ups · read %s\n", x.ObservedAtUTC)
		for _, f := range x.Items {
			printFollowup(cmd, f)
		}
	case agentwork.Followup:
		printFollowup(cmd, x)
	}
	return nil
}
func printFollowup(cmd *cobra.Command, f agentwork.Followup) {
	fmt.Fprintf(cmd.OutOrStdout(), "%s · revision %d · generation %d · cycle %d · %s\n  %s: %s · check %s · transport %s (%s)\n", f.ID, f.Revision, f.Generation, f.Cycle, f.State, f.WaitKind, f.WaitRef, f.NextCheckAtUTC, f.Transport.State, f.Transport.Reason)
	if f.Recheck != nil {
		r := f.Recheck
		fmt.Fprintf(cmd.OutOrStdout(), "  Self-report cycle %d %s at %s: summary=%s; done=%s; next=%s; blocker=%s\n", r.Cycle, r.Outcome, r.RecordedAtUTC, r.Summary, r.Done, r.Next, r.Blocker)
	}
}
func agentWorkFormat(cmd *cobra.Command, format string) error {
	if format != "text" && format != "json" {
		return fmt.Errorf("agent_work_invalid_format")
	}
	if workflowStatusJSON(cmd) && cmd.Flags().Changed("format") && format == "text" {
		return fmt.Errorf("agent_work_json_format_conflict")
	}
	return nil
}
func newAgentWorkCommand(factory agentWorkFactory) *cobra.Command {
	var root, format string
	var metadataOnly bool
	c := &cobra.Command{Use: "agents", SilenceUsage: true, Short: "Read public fresh agent observations and separate native work/follow-ups", Args: cobra.NoArgs, PersistentPreRunE: func(*cobra.Command, []string) error { return nil }, RunE: func(cmd *cobra.Command, _ []string) error {
		if e := agentWorkFormat(cmd, format); e != nil {
			return e
		}
		if metadataOnly {
			return agentWorkOutput(cmd, format, agentwork.StatusMetadata(agentwork.HerdrRuntime{Executable: "herdr"}))
		}
		s, e := factory(root)
		if e != nil {
			return e
		}
		o, e := s.Overview()
		if e != nil {
			return e
		}
		return agentWorkOutput(cmd, format, o)
	}}
	c.Flags().BoolVar(&metadataOnly, "metadata-only", false, "read only versioned public status metadata; no native work or history reads")
	c.Flags().StringVar(&root, "workspace", "", "physical Ply workspace root (default: discover)")
	c.Flags().StringVar(&format, "format", "text", "text or json")
	return c
}
func readAgentWorkInput(path string, out any) error {
	f, e := os.Open(path)
	if e != nil {
		return fmt.Errorf("agent_work_input_unavailable")
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, (16<<10)+1))
	if e != nil || len(b) > 16<<10 || !json.Valid(b) {
		return fmt.Errorf("agent_work_input_invalid")
	}
	if e := taskrun.StrictDecode(b, 16<<10, out); e != nil {
		return fmt.Errorf("agent_work_input_invalid")
	}

	return nil
}
func newFollowupCommand(factory agentWorkFactory) *cobra.Command {
	var root, format string
	c := &cobra.Command{Use: "followup", SilenceUsage: true, Short: "Register and dispatch durable fixed status-only follow-ups", PersistentPreRunE: func(*cobra.Command, []string) error { return nil }}
	c.PersistentFlags().StringVar(&root, "workspace", "", "physical Ply workspace root (default: discover)")
	c.PersistentFlags().StringVar(&format, "format", "text", "text or json")
	for _, op := range []string{"register", "show", "list", "request", "dispatch", "recheck", "cancel", "retarget", "request-status"} {
		op := op
		var path, requestID string
		var due bool
		verb := &cobra.Command{Use: op, Short: "Status follow-up " + op, Args: func(_ *cobra.Command, args []string) error {
			want := 1
			if op == "register" || op == "list" {
				want = 0
			}
			if len(args) != want {
				return fmt.Errorf("followup_expected_%d_arguments", want)
			}
			return nil
		}, RunE: func(cmd *cobra.Command, args []string) error {
			if e := agentWorkFormat(cmd, format); e != nil {
				return e
			}
			s, e := factory(root)
			if e != nil {
				return e
			}
			var result any
			switch op {
			case "register":
				var d agentwork.Draft
				if e = readAgentWorkInput(path, &d); e == nil {
					result, e = s.Register(d)
				}
			case "list":
				result, e = s.List(due)
			case "show":
				result, e = s.Show(args[0])
			case "request-status":
				result, e = s.RequestStatus(args[0], requestID)
			case "dispatch":
				result, e = s.Dispatch(args[0])
			default:
				var m agentwork.Mutation
				if e = readAgentWorkInput(path, &m); e == nil {
					result, e = s.Mutate(args[0], op, m)
				}
			}
			if e != nil {
				return e
			}
			return agentWorkOutput(cmd, format, result)
		}}
		if op == "list" {
			verb.Flags().BoolVar(&due, "due", false, "only waiting, due, unconsumed cycles")
		}
		if op == "register" || op == "request" || op == "recheck" || op == "cancel" || op == "retarget" {
			verb.Flags().StringVar(&path, "file", "", "bounded JSON request file")
			_ = verb.MarkFlagRequired("file")
		}
		if op == "request-status" {
			verb.Flags().StringVar(&requestID, "request-id", "", "exact stable mutation request ID")
			_ = verb.MarkFlagRequired("request-id")
		}
		c.AddCommand(verb)
	}
	return c
}
