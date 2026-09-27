package cmd

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestWorkflowHandoffCommandShape(t *testing.T) {
	command := newWorkflowHandoffCommandWithServices(successfulHandoffServices())
	root := &cobra.Command{Use: "ply"}
	workflow := &cobra.Command{Use: "workflow"}
	root.AddCommand(workflow)
	workflow.AddCommand(command)
	if command.Use != "handoff" || command.Short != "Manage file-based agent handoffs" || command.Long != "Create, control, receive, and inspect immutable local agent handoffs." || command.Runnable() {
		t.Fatalf("handoff metadata = %#v", command)
	}
	wantUses := map[string]string{"create": "create", "show": "show <handoff-id>", "inspect": "inspect", "submit-start": "submit-start", "submit-result": "submit-result", "cancel": "cancel <handoff-id>", "supersede": "supersede <handoff-id>", "abandon": "abandon <handoff-id>"}
	for name, wantUse := range wantUses {
		child, _, err := command.Find([]string{name})
		if err != nil || child == command || child.Use != wantUse || !child.Runnable() || child.Args == nil {
			t.Fatalf("%s = %#v, %v", name, child, err)
		}
	}
	create, _, _ := command.Find([]string{"create"})
	create.InitDefaultHelpFlag()
	flags := []string{}
	create.LocalNonPersistentFlags().VisitAll(func(flag *pflag.Flag) { flags = append(flags, flag.Name) })
	if !reflect.DeepEqual(flags, []string{"file", "help"}) {
		t.Fatalf("create flags = %v", flags)
	}
}

func TestWorkflowHandoffRendersServices(t *testing.T) {
	services := successfulHandoffServices()
	tests := []struct {
		args []string
		want string
	}{
		{[]string{"create", "--file", "/draft.json"}, "Created agent handoff hnd_0123456789abcdef0123456789abcdef.\nPurpose: Purpose\nWorking directory: /target\nHandoff: /handoff.json\nNext action: Open a fresh recipient agent in the working directory and tell it: \"Read and execute the handoff at /handoff.json.\"\n"},
		{[]string{"show", "hnd_0123456789abcdef0123456789abcdef"}, "Status: status\nResult: result\nMeaning: meaning\nNext action: action\n"},
		{[]string{"inspect", "--handoff", "/handoff.json", "--format", "json"}, "{}\n"},
		{[]string{"submit-start", "--handoff", "/handoff.json", "--file", "/start.json"}, "Accepted start receipt rcp_test.\nStart receipt: /start.json\nSHA-256: sha256:abc\n"},
	}
	for _, test := range tests {
		stdout, stderr, err := executeHandoffCommand(test.args, services)
		if err != nil || stderr != "" || stdout != test.want {
			t.Fatalf("%v = stdout %q stderr %q err %v, want %q", test.args, stdout, stderr, err, test.want)
		}
	}
}

func TestWorkflowHandoffRejectsArgumentsBeforeService(t *testing.T) {
	calls := 0
	services := successfulHandoffServices()
	services.create = func(workflowhandoff.CreateInput) (workflowhandoff.CreateResult, error) {
		calls++
		return workflowhandoff.CreateResult{}, nil
	}
	for _, args := range [][]string{{"create"}, {"create", "--unknown"}, {"show"}, {"inspect", "--handoff", "relative", "--format", "json"}} {
		stdout, _, err := executeHandoffCommand(args, services)
		if err == nil || stdout != "" {
			t.Fatalf("%v = stdout %q err %v", args, stdout, err)
		}
	}
	if calls != 0 {
		t.Fatalf("service calls = %d", calls)
	}
}

func successfulHandoffServices() workflowHandoffServices {
	return workflowHandoffServices{
		create: func(workflowhandoff.CreateInput) (workflowhandoff.CreateResult, error) {
			return workflowhandoff.CreateResult{HandoffID: "hnd_0123456789abcdef0123456789abcdef", Purpose: "Purpose", Worktree: "/target", Locator: "/handoff.json", Created: true}, nil
		},
		show: func(workflowhandoff.HandoffID) (workflowhandoff.ShowResult, error) {
			return workflowhandoff.ShowResult{Status: "status", Result: "result", Meaning: "meaning", NextAction: "action"}, nil
		},
		inspect: func(workflowhandoff.InspectInput) (workflowhandoff.InspectResult, error) {
			return workflowhandoff.InspectResult{Bytes: []byte("{}"), AppendLF: true}, nil
		},
		submitStart: func(workflowhandoff.SubmitInput) (workflowhandoff.SubmitResult, error) {
			return workflowhandoff.SubmitResult{Phase: "start", DocumentID: "rcp_test", Locator: "/start.json", SHA256: "sha256:abc", Created: true}, nil
		},
		submitResult: func(workflowhandoff.SubmitInput) (workflowhandoff.SubmitResult, error) {
			return workflowhandoff.SubmitResult{Phase: "terminal", DocumentID: "res_test", Locator: "/result.json", SHA256: "sha256:def", Created: true}, nil
		},
		cancel: func(workflowhandoff.ControlInput) (workflowhandoff.ControlResult, error) {
			return workflowhandoff.ControlResult{}, nil
		},
		supersede: func(workflowhandoff.SupersedeInput) (workflowhandoff.CreateResult, error) {
			return workflowhandoff.CreateResult{}, nil
		},
		abandon: func(workflowhandoff.ControlInput) (workflowhandoff.ControlResult, error) {
			return workflowhandoff.ControlResult{}, nil
		},
	}
}

func executeHandoffCommand(args []string, services workflowHandoffServices) (string, string, error) {
	command := newWorkflowHandoffCommandWithServices(services)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	command.SetOut(stdout)
	command.SetErr(stderr)
	command.SetArgs(args)
	command.SilenceErrors = true
	command.SilenceUsage = true
	err := command.Execute()
	return stdout.String(), stderr.String(), err
}
