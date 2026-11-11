package cmd

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/agentlaunch"
)

func TestAgentAssignmentBuildsPullPrompt(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	called := false
	command := newAgentAssignmentCommand(func(_ context.Context, in agentlaunch.Request) (agentlaunch.Result, error) {
		called = true
		if in.InitialPrompt != "/ply-assignment ASSIGNMENT-7 --project ply" || in.Workspace != "w7" || in.Label != "ASSIGNMENT-7" {
			t.Fatalf("%+v", in)
		}
		return agentlaunch.Result{State: "agent_observed", AgentStatus: "working"}, nil
	})
	var out bytes.Buffer
	command.SetOut(&out)
	command.SetArgs([]string{"ASSIGNMENT-7", "--project", "ply", "--herdr-workspace", "w7", "--claude-bin", exe, "--herdr-bin", exe})
	if err = command.Execute(); err != nil {
		t.Fatal(err)
	}
	if !called || !bytes.Contains(out.Bytes(), []byte("agent_observed")) {
		t.Fatal(out.String())
	}
}
func TestAgentAssignmentRejectsExtraInstruction(t *testing.T) {
	command := newAgentAssignmentCommand(func(context.Context, agentlaunch.Request) (agentlaunch.Result, error) {
		t.Fatal("start")
		return agentlaunch.Result{}, nil
	})
	command.SetArgs([]string{"ASSIGNMENT-7", "extra", "--project", "ply", "--herdr-workspace", "w7"})
	if command.Execute() == nil {
		t.Fatal("extra argument accepted")
	}
}
