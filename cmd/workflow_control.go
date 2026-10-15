package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/weekrun"
	"github.com/spf13/cobra"
)

func newWorkflowControlCommand(d taskrun.Dependencies) *cobra.Command {
	parent := &cobra.Command{Use: "control", Short: "Preview, approve and read request-bound native runtime controls", Long: "Read request-bound native runtime controls for the exact Run or week. Preview and status are read-only. Apply requires the complete unchanged approved control-request v1. A week-owned Run controls the entire preserved workset. Accepted means durably reserved; a separate receipt confirms the safe boundary. Old frozen and interactive runtimes remain unsupported. Keep the same request ID when the outcome is uncertain."}
	invalid := func(detail string) error { return taskrun.ControlError("control_invalid_arguments", detail) }
	service := func(root string, target taskrun.ControlTarget) (taskrun.ControlService, string, error) {
		if !taskrun.ValidateControlTarget(target) {
			return taskrun.ControlService{}, "", invalid("Use an exact Run ID or --week WEEK_ID.")
		}
		if target.Kind == "run" {
			resolved, e := taskrun.ControlTargetForRun(root, target.ID)
			if e != nil {
				return taskrun.ControlService{}, "", e
			}
			target = resolved
		}
		if target.Kind == "week" {
			return weekrun.New(d).Controls(), target.ID, nil
		}
		return taskrun.WorkerControls(d), target.ID, nil
	}
	write := func(c *cobra.Command, v any) error { return json.NewEncoder(c.OutOrStdout()).Encode(v) }
	for _, verb := range []string{"show", "preview", "apply", "status"} {
		v := verb
		var root, week, action, id, file, format string
		c := &cobra.Command{Use: v, Short: map[string]string{"show": "Show actions, scope and unavailable reasons without effects", "preview": "Bind an action to the current native revision for at most five minutes", "apply": "Reserve the exact approved request once; never treat receipt as liveness", "status": "Read the retained request and historical receipt without recovery or wake"}[v]}
		c.Args = func(c *cobra.Command, args []string) error {
			if format != "json" && format != "text" {
				return invalid("--format must be text or json")
			}
			if v == "apply" {
				if len(args) != 0 || file == "" {
					return invalid("apply requires only --file")
				}
				return nil
			}
			if root == "" || !filepath.IsAbs(root) || filepath.Clean(root) != root {
				return invalid("--workspace must name the physical absolute workspace")
			}
			if v == "status" {
				if len(args) != 1 || !taskrun.ValidateControlID(args[0]) {
					return invalid("status requires one valid request ID")
				}
				return nil
			}
			if (len(args) != 1 || week != "") && (len(args) != 0 || week == "") {
				return invalid("Choose exactly one Run ID or --week WEEK_ID")
			}
			return nil
		}
		c.RunE = func(c *cobra.Command, args []string) error {
			var value any
			var err error
			if v == "apply" {
				if !filepath.IsAbs(file) || filepath.Clean(file) != file {
					return taskrun.ControlError("control_invalid_request", "--file must be an absolute regular request file")
				}
				info, e := os.Lstat(file)
				if e != nil || !info.Mode().IsRegular() || info.Size() > 64<<10 {
					return taskrun.ControlError("control_invalid_request", "Cannot read a bounded regular request file")
				}
				f, e := os.Open(file)
				if e != nil {
					return taskrun.ControlError("control_invalid_request", "Cannot read the request file")
				}
				raw, e := io.ReadAll(io.LimitReader(f, 64<<10+1))
				f.Close()
				if e != nil {
					return taskrun.ControlError("control_invalid_request", "Cannot read the request file")
				}
				r, e := taskrun.ParseControlRequest(raw)
				if e != nil {
					return e
				}
				// Historical identity lookup must survive changed or missing owners.
				var svc taskrun.ControlService
				if r.Input.Target.Kind == "week" {
					svc = weekrun.New(d).Controls()
				} else {
					svc = taskrun.WorkerControls(d)
				}
				value, err = svc.Apply(r)
			} else if v == "status" {
				// The journal discovers its own original scope; no Run lookup aliases it.
				svc := taskrun.WorkerControls(d)
				s, e := svc.Status(root, args[0])
				if s.Scope != nil && s.Scope.Kind == "week" {
					s, e = weekrun.New(d).Controls().Status(root, args[0])
				}
				value, err = s, e
			} else {
				target := taskrun.ControlTarget{Kind: "week", ID: week}
				if week == "" {
					target = taskrun.ControlTarget{Kind: "run", ID: args[0]}
				}
				svc, targetID, e := service(root, target)
				if e != nil {
					return e
				}
				if v == "show" {
					value, err = svc.Show(root, targetID)
				} else {
					value, err = svc.Preview(root, targetID, action, id)
				}
			}
			if e := write(c, value); e != nil {
				return e
			}
			return err
		}
		c.Flags().StringVar(&format, "format", "json", "output format (json or text; both retain the complete structured contract)")
		if v == "apply" {
			c.Flags().StringVar(&file, "file", "", "absolute approved control-request v1 JSON file")
		} else {
			c.Flags().StringVar(&root, "workspace", "", "physical workspace root")
		}
		if v == "show" || v == "preview" {
			c.Flags().StringVar(&week, "week", "", "exact parent week ID, including between Goals; mutually exclusive with Run ID")
		}
		if v == "preview" {
			c.Flags().StringVar(&action, "action", "", "pause, stop or resume")
			c.Flags().StringVar(&id, "request-id", "", "stable request identity: 8-128 ASCII letters, digits, dots, underscores or hyphens")
		}
		c.SetFlagErrorFunc(func(c *cobra.Command, e error) error { return invalid(e.Error()) })
		parent.AddCommand(c)
	}
	parent.SetFlagErrorFunc(func(c *cobra.Command, e error) error { return invalid(fmt.Sprint(e)) })
	return parent
}
