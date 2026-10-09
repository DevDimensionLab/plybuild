package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/weekrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkflowWeekCommand(d taskrun.Dependencies) *cobra.Command {
	m := weekrun.New(d)
	parent := &cobra.Command{Use: "week", Short: "Run an explicitly selected set of native Goals through automatic local delivery", Long: "Preserve selected Task/Goal revisions, dependencies, model, return target and a deadline up to seven days. Native workers own development, tests, review and local return. The week controller chooses only approved work, continues independent work after proven blockers, and preserves active Run identity after interruption. No human test gate, global provider settings, push or notification is added."}
	var file, format string
	var check, service bool
	start := &cobra.Command{Use: "start", Args: cobra.NoArgs, Short: "Bind selected Goals and start the durable local controller", RunE: func(c *cobra.Command, _ []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		contract, err := weekrun.ReadContract(file)
		if err != nil {
			return err
		}
		s, err := m.Start(contract, check)
		if err == nil && !check {
			if service {
				_, err = m.Service(contract.Workspace, s.ID, "install")
			} else {
				err = m.LaunchRun(contract.Workspace, s.ID)
			}
		}
		if e := writeWeek(c, format, s); e != nil {
			return e
		}
		return err
	}}
	start.Flags().StringVar(&file, "file", "", "absolute JSON week contract path")
	start.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	start.Flags().BoolVar(&check, "check", false, "validate selected native bindings without writes or startup")
	start.Flags().BoolVar(&service, "service", false, "install and bootstrap a named per-user launchd service using the frozen control")
	_ = start.MarkFlagRequired("file")
	setWorkFlagErrors(start)
	parent.AddCommand(start)
	for _, verb := range []string{"show", "pause", "stop", "resume", "_serve", "_supervise"} {
		v := verb
		var root, f string
		c := &cobra.Command{Use: v + " WEEK_ID", Args: cobra.ExactArgs(1), Hidden: v[0] == '_', RunE: func(c *cobra.Command, args []string) error {
			if err := executeFormat(f); err != nil {
				return err
			}
			if root == "" {
				o, e := workspace.ObserveContaining(d.Workspace)
				if e != nil {
					return e
				}
				root = o.Root
			}
			if v == "_serve" || v == "_supervise" {
				ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
				defer stop()
				if v == "_supervise" {
					return m.Supervise(ctx, root, args[0])
				}
				return m.Serve(ctx, root, args[0])
			}
			var s weekrun.State
			var err error
			if v == "show" {
				s, err = m.Show(root, args[0])
			} else {
				s, err = m.Control(root, args[0], v)
			}
			if e := writeWeek(c, f, s); e != nil {
				return e
			}
			return err
		}}
		c.Flags().StringVar(&root, "workspace", "", "physical workspace root (default: containing workspace)")
		c.Flags().StringVar(&f, "format", "text", "output format (text or json)")
		setWorkFlagErrors(c)
		parent.AddCommand(c)
	}
	svc := &cobra.Command{Use: "service", Short: "Install, observe, start, stop or remove this week's named launchd service"}
	for _, verb := range []string{"install", "start", "status", "stop", "uninstall"} {
		v := verb
		var root, f string
		c := &cobra.Command{Use: v + " WEEK_ID", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
			if err := executeFormat(f); err != nil {
				return err
			}
			if root == "" {
				o, e := workspace.ObserveContaining(d.Workspace)
				if e != nil {
					return e
				}
				root = o.Root
			}
			out, err := m.Service(root, args[0], v)
			if e := writeWeek(c, f, out); e != nil {
				return e
			}
			return err
		}}
		c.Flags().StringVar(&root, "workspace", "", "physical workspace root (default: containing workspace)")
		c.Flags().StringVar(&f, "format", "text", "output format (text or json)")
		setWorkFlagErrors(c)
		svc.AddCommand(c)
	}
	parent.AddCommand(svc)
	return parent
}
func writeWeek(c *cobra.Command, format string, v any) error {
	if format == "json" {
		b, err := taskrun.Canonical(v)
		if err != nil {
			return err
		}
		_, err = c.OutOrStdout().Write(append(b, '\n'))
		return err
	}
	switch s := v.(type) {
	case weekrun.State:
		fmt.Fprintf(c.OutOrStdout(), "Week: %s\nState: %s\nStarted: %s\nDeadline: %s\nLast progress: %s\nNext: %s\n", s.ID, s.Phase, s.StartedAt, s.Deadline, s.LastProgressAt, s.NextAction)
		for _, e := range s.Entries {
			fmt.Fprintf(c.OutOrStdout(), "%s: %s", e.TaskID, e.State)
			if e.RunID != "" {
				fmt.Fprintf(c.OutOrStdout(), " (%s)", e.RunID)
			}
			if e.Reason != "" {
				fmt.Fprintf(c.OutOrStdout(), " — %s", e.Reason)
			}
			fmt.Fprintln(c.OutOrStdout())
		}
		if s.Reason != "" {
			fmt.Fprintln(c.OutOrStdout(), s.Reason)
		}
	case weekrun.ServiceStatus:
		fmt.Fprintf(c.OutOrStdout(), "Service: %s\nState: %s\nPlist: %s\nControl: %s\n", s.Label, s.State, s.PlistPath, s.ControlPath)
	}
	return nil
}
