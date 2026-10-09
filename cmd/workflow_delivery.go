package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/delivery"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

type workflowDeliveryService interface {
	Register(cwd, file string) (delivery.Receipt, error)
	Check(cwd, id string) (delivery.Receipt, error)
	Execute(cwd, id string) (delivery.Receipt, error)
	Read(cwd, id string) (delivery.Receipt, error)
	List(cwd, taskID, epicID string) (delivery.ListResult, error)
	Metadata(cwd, id, file string) (delivery.Receipt, error)
}

func newWorkflowDeliveryCommand(service workflowDeliveryService) *cobra.Command {
	return newWorkflowDeliveryCommandWithFactory(func(string) workflowDeliveryService { return service })
}

func newWorkflowDeliveryCommandWithFactory(factory func(string) workflowDeliveryService) *cobra.Command {
	parent := &cobra.Command{
		Use: "delivery", Short: "Register and deliver an exact verified Task candidate",
		Long:    "Preserve a durable Delivery for the Task's frozen agreement: a pull request, local Epic integration, or explicitly selected local branch integration. Registration preserves references without Git or PR effects. Execution rechecks native authority, evidence and the frozen automatic or human acceptance gate for the exact candidate. A pending Delivery requires another explicit invocation; no background worker is started.",
		Example: "  ply workflow delivery register --file /absolute/delivery.json --format json\n  ply workflow delivery check dlv_<digest>\n  ply workflow delivery execute dlv_<digest>\n  ply workflow delivery show dlv_<digest> --format json",
		Args:    cobra.NoArgs,
		RunE:    func(c *cobra.Command, _ []string) error { return c.Help() },
	}
	for _, verb := range []string{"register", "check", "execute", "show", "list", "metadata"} {
		verb := verb
		var file, format, taskID, epicID, preferencesFile string
		use := verb + " DELIVERY_ID"
		if verb == "register" || verb == "list" {
			use = verb
		}
		short := map[string]string{
			"register": "Register an exact technical candidate without delivery effects",
			"check":    "Check the current candidate, authority and acceptance gate without effects",
			"execute":  "Perform or resume the frozen delivery after its exact acceptance gate",
			"show":     "Read a Delivery receipt and its preserved history",
			"list":     "Find Deliveries by Task or Epic",
			"metadata": "Preserve an explicit revision of PR people metadata",
		}[verb]
		detail := map[string]string{
			"register": "Read a private ply.delivery.registration@1 JSON file. The native TaskResult and frozen workflow supply the candidate, Spec, agreement and authority. Reusing a publication key with identical content returns the same Delivery; changed content is a conflict. Registration is not completed delivery or human approval.",
			"check":    "Read-only. Recheck the frozen candidate and evidence, actual runtime authority, exact target, and candidate-bound automatic acceptance or required human QA. The receipt explains the next action. This command neither records human judgment nor publishes a branch or integrates code.",
			"execute":  "Recheck all gates and perform only the frozen mode and target. Pull requests publish only the source branch and stop before merge. Local modes update the exact registered local return and preserve queue/base evidence. After an interrupted effect, use this same ID to observe and resume it. Unknown outcomes stop blind repetition. PR merge, auto-merge, force-push, and target-branch push are excluded.",
			"show":     "Read-only. Show the durable manifest, events, actual effect receipts, blockers and next action. A completed receipt preserves delivery-time observations; it does not assert the current state of a remote PR.",
			"list":     "Read-only. Filter by the registered Task or Epic, or show all Deliveries in the containing workspace. Filters combine as exact matches. Reading never starts or resumes a Delivery.",
			"metadata": "Read an explicit metadata revision from a private JSON object with optional assignees, reviewers, remove_assignees and remove_reviewers string arrays. For example, {\"assignees\":[],\"reviewers\":[]} clears the desired additions and preserves existing PR people. Preserve the revised choices in history; execute the same Delivery to apply them. Only explicit remove arrays remove people. A metadata revision does not change the candidate, mode, target or human QA.",
		}[verb]
		example := "  ply workflow delivery " + verb
		if verb != "register" && verb != "list" {
			example += " dlv_<digest>"
		}
		if verb == "register" || verb == "metadata" {
			example += " --file /absolute/" + verb + ".json"
		}
		example += " --format json"
		child := &cobra.Command{Use: use, Short: short, Long: short + ". " + detail, Example: example,
			Args: func(c *cobra.Command, args []string) error {
				want := 1
				if verb == "register" || verb == "list" {
					want = 0
				}
				if len(args) != want {
					return workspace.WorkInvalidArguments(fmt.Sprintf("expected %d positional arguments", want))
				}
				if err := executeFormat(format); err != nil {
					return err
				}
				if verb == "register" || verb == "metadata" {
					if !filepath.IsAbs(file) {
						return workspace.WorkInvalidArguments("--file requires an absolute path to private JSON input")
					}
				}
				if c.Flags().Changed("preferences-file") && !filepath.IsAbs(preferencesFile) {
					return workspace.WorkInvalidArguments("--preferences-file requires an absolute path")
				}
				if c.Flags().Changed("task") {
					if _, err := workspace.ParseTaskID(taskID); err != nil {
						return err
					}
				}
				if c.Flags().Changed("epic") {
					if _, err := workspace.ParseEpicID(epicID); err != nil {
						return err
					}
				}
				return nil
			},
			RunE: func(c *cobra.Command, args []string) error {
				cwd, err := os.Getwd()
				if err != nil {
					return err
				}
				service := factory(preferencesFile)
				if verb == "list" {
					result, callErr := service.List(cwd, taskID, epicID)
					if err = writeDeliveryOutput(c, format, result, delivery.ListText(result)); err != nil {
						return err
					}
					return callErr
				}
				var result delivery.Receipt
				var callErr error
				switch verb {
				case "register":
					result, callErr = service.Register(cwd, file)
				case "check":
					result, callErr = service.Check(cwd, args[0])
				case "execute":
					result, callErr = service.Execute(cwd, args[0])
				case "show":
					result, callErr = service.Read(cwd, args[0])
				case "metadata":
					result, callErr = service.Metadata(cwd, args[0], file)
				}
				// Even when a later operation failed, an observed PR or partial
				// local return must remain visible to the caller.
				if err = writeDeliveryOutput(c, format, result, delivery.Text(result)); err != nil {
					return err
				}
				if callErr == nil && verb == "execute" && result.State != "delivered" {
					return &taskrun.Error{Code: "delivery_" + result.State, Detail: result.NextAction, Exit: 4}
				}
				return callErr
			},
		}
		child.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		if verb == "register" || verb == "metadata" {
			child.Flags().StringVar(&file, "file", "", "absolute path to a private JSON input file (required)")
		}
		if verb == "list" {
			child.Flags().StringVar(&taskID, "task", "", "filter by exact registered Task ID")
			child.Flags().StringVar(&epicID, "epic", "", "filter by exact registered Epic ID")
		}
		if verb == "execute" {
			child.Flags().StringVar(&preferencesFile, "preferences-file", "", "explicit local PR preference file for the first attempt (default: ~/.agents/local/pr-preferences.yaml); preserved choices win on retry")
		}
		setWorkFlagErrors(child)
		parent.AddCommand(child)
	}
	setWorkFlagErrors(parent)
	return parent
}

func writeDeliveryOutput(c *cobra.Command, format string, result any, rendered string) error {
	if format == "json" {
		b, err := taskrun.Canonical(result)
		if err != nil {
			return err
		}
		_, err = c.OutOrStdout().Write(append(b, '\n'))
		return err
	}
	_, err := fmt.Fprint(c.OutOrStdout(), rendered)
	return err
}
