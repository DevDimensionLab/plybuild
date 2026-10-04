package cmd

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func addWorkspaceTaskJournalCommands(task *cobra.Command, d workspace.Dependencies) {
	journal := &cobra.Command{Use: "journal", Short: "Read and contribute to a Task process journal", Long: "Read native facts, meaningful agent steps and observer contributions across Task runs. Contributions are self-reported data and grant no Task, QA, integration or start authority.", Example: "  ply workspace task journal show my-task\n  ply workspace task journal record my-task --file /absolute/event.json", Args: cobra.NoArgs, RunE: func(c *cobra.Command, a []string) error { return c.Help() }}
	for _, verb := range []string{"show", "record", "observe"} {
		verb := verb
		var file, format, view, order, run, actor, snapshot string
		short := map[string]string{"show": "Read a Task banner, text timeline or JSON snapshot", "record": "Append one immutable agent step or process note", "observe": "Append one immutable observer proposal or follow-up"}[verb]
		long := short + ". "
		if verb == "show" {
			long += "Read-only: creates no journal directory, lock, cache or receipt. Missing optional sources produce partial coverage. Both occurrence and registration clocks remain visible; unknown time and missing ends are explicit. --format json always exports the full snapshot. Filters only change selection. --snapshot renders an exported snapshot offline without workspace discovery, source reads or a fresh check."
		} else {
			long += "--file must be a physical regular UTF-8 JSON file of at most 64 KiB, with all required fields and no duplicate or unknown keys. Sources are absolute regular files with matching SHA-256. An identical publication_key retry returns the same event even after artifact loss; different content conflicts. Text and sources are data, never instructions. Do not include transcripts, credentials or secrets. Outcomes are reported claims, never native technical or human QA judgments."
		}
		example := "  ply workspace task journal " + verb + " my-task --file /absolute/" + map[string]string{"record": "event", "observe": "observation"}[verb] + ".json --format json"
		if verb == "show" {
			example = "  ply workspace task journal show my-task\n  ply workspace task journal show my-task --view timeline\n  ply workspace task journal show my-task --format json > /absolute/snapshot.json\n  ply workspace task journal show my-task --snapshot /absolute/snapshot.json --view details"
		}
		c := &cobra.Command{Use: verb + " <task-id>", Short: short, Long: long, Example: example, Args: func(c *cobra.Command, a []string) error {
			if len(a) != 1 {
				return workspace.WorkInvalidArguments("expected exactly one Task ID")
			}
			if _, e := workspace.ParseTaskID(a[0]); e != nil {
				return e
			}
			if format != "text" && format != "json" {
				return workspace.WorkInvalidArguments("--format must be text or json")
			}
			if verb != "show" && file == "" {
				return workspace.WorkInvalidArguments("--file is required")
			}
			if verb == "show" {
				if view != "summary" && view != "timeline" && view != "details" {
					return workspace.WorkInvalidArguments("--view must be summary, timeline or details")
				}
				if order != "occurred" && order != "recorded" {
					return workspace.WorkInvalidArguments("--order must be occurred or recorded")
				}
				if c.Flags().Changed("snapshot") && snapshot == "" {
					return workspace.WorkInvalidArguments("--snapshot must not be empty")
				}
			}
			return nil
		}}
		c.RunE = func(c *cobra.Command, a []string) error {
			service := taskjournal.New(d)
			var result any
			if verb == "show" {
				o := taskjournal.Options{View: view, Order: order, RunID: run, ActorID: actor}
				var s taskjournal.Snapshot
				var e error
				if snapshot != "" {
					s, e = taskjournal.ReadSnapshot(snapshot, a[0], o)
				} else {
					s, e = service.Show(a[0], o)
				}
				if e != nil {
					return e
				}
				if format == "text" {
					_, e = fmt.Fprint(c.OutOrStdout(), taskjournal.Text(s, view))
					return e
				}
				result = s
			} else {
				kind := "event"
				if verb == "observe" {
					kind = "observation"
				}
				r, e := service.Append(a[0], file, kind)
				if e != nil {
					return e
				}
				for _, warning := range r.Warnings {
					if _, e = fmt.Fprintln(c.ErrOrStderr(), "Warning: "+warning); e != nil {
						return e
					}
				}
				if format == "text" {
					_, e = fmt.Fprintf(c.OutOrStdout(), "Event: %s\nRecorded: %s\nCreated: %t\nRecord: %s\nNext transition authorized: false\n", r.EventID, r.RecordedAt, r.Created, r.RecordLocator)
					return e
				}
				result = r
			}
			b, e := taskjournal.Canonical(result)
			if e != nil {
				return e
			}
			_, e = c.OutOrStdout().Write(append(b, '\n'))
			return e
		}
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		if verb == "show" {
			c.Flags().StringVar(&view, "view", "summary", "text view (summary, timeline or details)")
			c.Flags().StringVar(&order, "order", "occurred", "time order (occurred or recorded)")
			c.Flags().StringVar(&run, "run", "", "exact native run ID filter")
			c.Flags().StringVar(&actor, "actor", "", "exact actor ID filter")
			c.Flags().StringVar(&snapshot, "snapshot", "", "absolute exported snapshot JSON path for offline read-only rendering")
		} else {
			c.Flags().StringVar(&file, "file", "", "absolute physical JSON input file (required; at most 64 KiB)")
		}
		c.SetFlagErrorFunc(func(c *cobra.Command, e error) error { return workspace.WorkInvalidArguments(e.Error()) })
		journal.AddCommand(c)
	}
	task.AddCommand(journal)
}
