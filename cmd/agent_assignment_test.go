package cmd

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
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
		if in.InitialPrompt != "/ply-assignment ASSIGNMENT-7 --project ply" || in.Agent != "claude" || in.Executable != exe || in.Workspace != "w7" || in.Label != "ASSIGNMENT-7" {
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

func TestAgentAssignmentStartsCodexWithSkillMention(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	codex := filepath.Join(t.TempDir(), "codex")
	if err = os.WriteFile(codex, nil, 0700); err != nil {
		t.Fatal(err)
	}
	called := false
	command := newAgentAssignmentCommand(func(_ context.Context, in agentlaunch.Request) (agentlaunch.Result, error) {
		called = true
		if in.InitialPrompt != "$ply-assignment ASSIGNMENT-7 --project ply" || in.Agent != "codex" || in.Executable != codex {
			t.Fatalf("%+v", in)
		}
		return agentlaunch.Result{State: "agent_observed", Agent: "codex"}, nil
	})
	var out bytes.Buffer
	command.SetOut(&out)
	command.SetArgs([]string{"ASSIGNMENT-7", "--project", "ply", "--herdr-workspace", "w7", "--provider", "codex", "--claude-bin", exe, "--codex-bin", codex, "--herdr-bin", exe})
	if err = command.Execute(); err != nil {
		t.Fatal(err)
	}
	if !called || !bytes.Contains(out.Bytes(), []byte(`"agent":"codex"`)) {
		t.Fatal(out.String())
	}
}

func TestAgentAssignmentRejectsUnknownProviderBeforeLookup(t *testing.T) {
	for _, provider := range []string{"", "gemini", "Codex", "claude "} {
		command := newAgentAssignmentCommand(func(context.Context, agentlaunch.Request) (agentlaunch.Result, error) {
			t.Fatal("start")
			return agentlaunch.Result{}, nil
		})
		var out bytes.Buffer
		command.SetOut(&out)
		command.SetErr(&out)
		command.SetArgs([]string{"ASSIGNMENT-7", "--project", "ply", "--herdr-workspace", "w7", "--provider", provider, "--claude-bin", "/absent/claude", "--codex-bin", "/absent/codex", "--herdr-bin", "/absent/herdr"})
		err := command.Execute()
		if err == nil || !strings.Contains(err.Error(), "--provider") {
			t.Fatalf("provider %q: %v", provider, err)
		}
	}
}

func TestAssignmentSkillLabelFollowsActualProvider(t *testing.T) {
	skill, err := os.ReadFile(filepath.Join("..", "skills", "ply-assignment", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(skill)
	for _, rule := range []string{
		"`claude:${CLAUDE_SESSION_ID}`",
		"`codex:` plus random hex",
		"never another provider's prefix",
		"Never send done in the same batch",
	} {
		if !strings.Contains(text, rule) {
			t.Fatalf("skill lacks %q", rule)
		}
	}
}
