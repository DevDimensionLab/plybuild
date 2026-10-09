package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func TestExecuteRegistersQAReservationRecovery(t *testing.T) {
	command, _, err := newWorkflowExecuteCommand(taskrun.Dependencies{}).Find([]string{"recover-qa"})
	if err != nil || command.Name() != "recover-qa" {
		t.Fatalf("known rejected QA reservation has no supported recovery command: %v", err)
	}
}

func qaRecoveryCLIArgs() []string {
	return []string{"wfr_fixture", "--attempt", "qa-00000007", "--context", "/fixture/context.json", "--runtime-evidence", "/fixture/current-runtime.json"}
}

func runQARecoveryCLI(actions workflowQARecoveryActions, args ...string) ([]byte, error) {
	c := workflowExecuteQARecoveryCommand(actions)
	var output bytes.Buffer
	c.SetOut(&output)
	c.SetErr(&bytes.Buffer{})
	c.SilenceUsage, c.SilenceErrors = true, true
	c.SetArgs(args)
	err := c.Execute()
	return output.Bytes(), err
}

func TestExecuteQARecoveryRejectsAmbiguousInputsBeforeWorkspaceAccess(t *testing.T) {
	base := qaRecoveryCLIArgs()
	for _, args := range [][]string{
		{}, append(append([]string{}, base...), "extra", "--check"),
		{"wfr_fixture", "--check"},
		{"wfr_fixture", "--attempt", "qa-00000007", "--context", "/fixture/context.json", "--check"},
		{"wfr_fixture", "--runtime-evidence", "/fixture/runtime.json", "--context", "/fixture/context.json", "--check"},
		{"wfr_fixture", "--runtime-evidence", "/fixture/runtime.json", "--attempt", "qa-00000007", "--check"},
		base,
		append(append([]string{}, base...), "--confirm", "not-a-preview-digest"),
		append(append([]string{}, base...), "--confirm", "sha256:"+strings.Repeat("z", 64)),
		append(append([]string{}, base...), "--check", "--confirm", "sha256:"+strings.Repeat("a", 64)),
		append(append([]string{}, base...), "--check", "--format", "yaml"),
		append(append([]string{}, base...), "--check", "--outcome", "pass"),
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			actions := workflowQARecoveryActions{root: func() (string, error) {
				t.Fatal("invalid recovery input reached workspace observation")
				return "", nil
			}}
			if _, err := runQARecoveryCLI(actions, args...); err == nil {
				t.Fatal("invalid recovery input was accepted")
			}
		})
	}
}

func TestExecuteQARecoveryCheckDoesNotApplyAndExposesExactConfirmation(t *testing.T) {
	confirmation := "sha256:" + strings.Repeat("a", 64)
	previewCalls, applyCalls := 0, 0
	actions := workflowQARecoveryActions{
		root: func() (string, error) { return "/fixture/workspace", nil },
		preview: func(root string, input taskrun.DeliveryQARecoveryInput) (taskrun.DeliveryQARecoveryPreview, error) {
			previewCalls++
			if root != "/fixture/workspace" || input != (taskrun.DeliveryQARecoveryInput{RunID: "wfr_fixture", AttemptID: "qa-00000007", ContextPath: "/fixture/context.json", RuntimeEvidencePath: "/fixture/current-runtime.json"}) {
				t.Fatalf("preview changed the original run, attempt or owner evidence: %q %+v", root, input)
			}
			return taskrun.DeliveryQARecoveryPreview{RunID: input.RunID, AttemptID: input.AttemptID, State: "ready", Confirmation: confirmation, NextAction: "Inspect this exact rejected reservation before confirming recovery."}, nil
		},
		recover: func(string, taskrun.DeliveryQARecoveryInput, string) (taskrun.DeliveryQARecoveryPreview, error) {
			applyCalls++
			return taskrun.DeliveryQARecoveryPreview{}, errors.New("check invoked recovery")
		},
	}
	output, err := runQARecoveryCLI(actions, append(qaRecoveryCLIArgs(), "--check", "--format", "json")...)
	var result taskrun.DeliveryQARecoveryPreview
	if err != nil || json.Unmarshal(output, &result) != nil || result.State != "ready" || result.Confirmation != confirmation || result.Receipt != nil || previewCalls != 1 || applyCalls != 0 {
		t.Fatalf("read-only check applied recovery or lost confirmation: %v\n%s", err, output)
	}
	output, err = runQARecoveryCLI(actions, append(qaRecoveryCLIArgs(), "--check")...)
	if err != nil || !strings.Contains(string(output), "ready") || !strings.Contains(string(output), confirmation) || !strings.Contains(string(output), result.NextAction) || applyCalls != 0 {
		t.Fatalf("text preview hid its exact state, confirmation or next action: %v\n%s", err, output)
	}
}

func TestExecuteQARecoveryApplyForwardsExactAttemptAndConfirmation(t *testing.T) {
	confirmation := "sha256:" + strings.Repeat("b", 64)
	applied := 0
	actions := workflowQARecoveryActions{
		root: func() (string, error) { return "/fixture/workspace", nil },
		preview: func(string, taskrun.DeliveryQARecoveryInput) (taskrun.DeliveryQARecoveryPreview, error) {
			t.Fatal("apply silently substituted a new preview confirmation")
			return taskrun.DeliveryQARecoveryPreview{}, nil
		},
		recover: func(root string, input taskrun.DeliveryQARecoveryInput, gotConfirmation string) (taskrun.DeliveryQARecoveryPreview, error) {
			applied++
			if root != "/fixture/workspace" || input.RunID != "wfr_fixture" || input.AttemptID != "qa-00000007" || input.ContextPath != "/fixture/context.json" || input.RuntimeEvidencePath != "/fixture/current-runtime.json" || gotConfirmation != confirmation {
				t.Fatalf("apply changed the frozen repair selection: %q %+v %q", root, input, gotConfirmation)
			}
			return taskrun.DeliveryQARecoveryPreview{State: "recovered", Receipt: &taskrun.FileBinding{Locator: "/fixture/recovery.json", SHA256: "sha256:" + strings.Repeat("c", 64)}, NextAction: "Use the original control to submit the actual human answer."}, nil
		},
	}
	output, err := runQARecoveryCLI(actions, append(qaRecoveryCLIArgs(), "--confirm", confirmation, "--format", "json")...)
	var result taskrun.DeliveryQARecoveryPreview
	if err != nil || json.Unmarshal(output, &result) != nil || applied != 1 || result.State != "recovered" || result.Receipt == nil || result.Receipt.Locator != "/fixture/recovery.json" {
		t.Fatalf("explicit repair lost its actual native receipt: %v\n%s", err, output)
	}
}

func TestExecuteQARecoveryPreservesNativeRejectionAndObservationFailure(t *testing.T) {
	confirmation := "sha256:" + strings.Repeat("d", 64)
	native := &taskrun.Error{Code: "delivery_qa_recovery_stale", Detail: "The exact QA reservation changed after preview.", Exit: 4}
	actions := workflowQARecoveryActions{
		root: func() (string, error) { return "/fixture/workspace", nil },
		recover: func(string, taskrun.DeliveryQARecoveryInput, string) (taskrun.DeliveryQARecoveryPreview, error) {
			return taskrun.DeliveryQARecoveryPreview{State: "blocked", NextAction: native.Detail}, native
		},
	}
	output, err := runQARecoveryCLI(actions, append(qaRecoveryCLIArgs(), "--confirm", confirmation, "--format", "json")...)
	var result taskrun.DeliveryQARecoveryPreview
	if !errors.Is(err, native) || json.Unmarshal(output, &result) != nil || result.State != "blocked" || result.NextAction != native.Detail || result.Receipt != nil {
		t.Fatalf("native rejection was hidden or presented as successful repair: %v\n%s", err, output)
	}
	observationErr := errors.New("the original Task workspace is unavailable")
	actions.root = func() (string, error) { return "", observationErr }
	actions.recover = func(string, taskrun.DeliveryQARecoveryInput, string) (taskrun.DeliveryQARecoveryPreview, error) {
		t.Fatal("workspace observation failure reached recovery")
		return taskrun.DeliveryQARecoveryPreview{}, nil
	}
	output, err = runQARecoveryCLI(actions, append(qaRecoveryCLIArgs(), "--confirm", confirmation, "--format", "json")...)
	if !errors.Is(err, observationErr) || len(output) != 0 {
		t.Fatalf("workspace observation failure produced a native recovery claim: %v\n%s", err, output)
	}
}
