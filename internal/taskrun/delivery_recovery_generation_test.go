package taskrun

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// A failed durable prompt reservation leaves a real session-bound fixture with
// no Task prompt attempt. This models the boundary without rewriting run state.
type deliveryRecoveryRejectTaskPrompt struct{}

func (deliveryRecoveryRejectTaskPrompt) WriteOnce(path string, data []byte) error {
	if filepath.Base(path) == "prompt-attempt.json" {
		return errors.New("fixture stopped before durable Task prompt reservation")
	}
	return (systemFileSystem{}).WriteOnce(path, data)
}

func (deliveryRecoveryRejectTaskPrompt) Replace(path string, data []byte) error {
	return (systemFileSystem{}).Replace(path, data)
}

func TestDeliveryRecoveryManagedPendingAgentIsNotAnExitedProvider(t *testing.T) {
	f, original, in := deliveryRecoveryFixture(t)
	workflowTestModel(t, f.workflowFixture, map[string]any{
		"recovery_managed_present": true,
		"launch_pending":           true,
		"interactive_ready":        false,
	})
	before := deliveryRecoveryFiles(t, original.Paths.RunRoot)
	calls := map[string]int{}
	for _, command := range []string{"tab create", "agent start", "agent prompt"} {
		calls[command] = workflowTestCalls(t, f.workflowFixture, command)
	}
	// Herdr can report an empty pane and a sole foreground shell while a managed
	// agent start still awaits input processing. Its agent read remains present.
	if p, err := WorkflowPreviewDeliveryStartRecovery(f.D, f.R.WorkspaceRoot, in); err == nil || p.State == "ready" {
		t.Errorf("managed launch-pending agent was promoted to provider exit: %+v %v", p, err)
	}
	deliveryRecoveryUnchanged(t, before, "")
	if len(deliveryRecoveryFiles(t, original.Paths.RunRoot)) != len(before) {
		t.Error("blocked managed-agent preview created recovery artifacts")
	}
	for command, count := range calls {
		if actual := workflowTestCalls(t, f.workflowFixture, command); actual != count {
			t.Errorf("blocked managed-agent preview emitted %q: before=%d after=%d", command, count, actual)
		}
	}
}

func TestDeliveryRecoveryGenerationOldReadyOwnerCannotReserveTaskPrompt(t *testing.T) {
	f, original, in := deliveryRecoveryFixture(t)
	p := deliveryRecoveryTestPreview(t, f, in)
	deliveryRecoveryPreserveControl(t, in, p)
	reserve := f.D
	reserve.Files = deliveryRecoveryRejectTaskPrompt{}
	if _, err := WorkflowRecoverDeliveryStart(reserve, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
		t.Fatal("fixture did not stop before reserving its Task prompt")
	}
	state, err := workflowRead(f.R.WorkspaceRoot, original.RunID)
	if err != nil || state.Recovery == nil || state.Phase != "session_bound" || state.Result.Transport.AgentSessionID == "" {
		t.Fatalf("replacement did not reach pre-Task session readiness: %+v %v", state, err)
	}
	if _, err = os.Lstat(filepathForRound(state, "prompt-attempt.json")); !os.IsNotExist(err) {
		t.Fatalf("fixture already reserved a Task prompt: %v", err)
	}
	before := deliveryRecoveryFiles(t, original.Paths.RunRoot)
	calls := workflowTestCalls(t, f.workflowFixture, "agent prompt")
	// The original readiness owner already captured "original" before recovery.
	// Entering the prompt helper later must preserve that ownership boundary.
	if err = workflowPromptGenerationUntil(f.D, f.R.WorkspaceRoot, original.RunID, nil, time.Now().Add(5*time.Second), "original"); err == nil {
		t.Error("older ready owner adopted replacement generation for the Task prompt")
	}
	deliveryRecoveryUnchanged(t, before, "")
	if len(deliveryRecoveryFiles(t, original.Paths.RunRoot)) != len(before) {
		t.Error("older ready owner created a replacement Task prompt artifact")
	}
	if actual := workflowTestCalls(t, f.workflowFixture, "agent prompt"); actual != calls {
		t.Errorf("older ready owner sent a Task prompt: before=%d after=%d", calls, actual)
	}
}

// A recovery sender owns the generation it reserved. A later, explicitly
// selected recovery must not let that delayed sender adopt its new transport.
func TestDeliveryRecoveryGenerationLateOwnerCannotSendOrReplaceAcceptedState(t *testing.T) {
	f, original, in := deliveryRecoveryFixture(t)
	in.Timeout = 30 * time.Second
	first := deliveryRecoveryTestPreview(t, f, in)
	deliveryRecoveryPreserveControl(t, in, first)
	paused, release, finished := make(chan struct{}), make(chan struct{}), make(chan struct{})
	var once sync.Once
	resume := func() { once.Do(func() { close(release) }) }
	delayed := f.D
	delayed.Fault = func(point string) error {
		if point == "delivery_recovery_before_agent_send" {
			close(paused)
			<-release
		}
		return nil
	}
	var firstErr error
	go func() {
		_, firstErr = WorkflowRecoverDeliveryStart(delayed, f.R.WorkspaceRoot, in, first.Confirmation)
		close(finished)
	}()
	t.Cleanup(func() {
		resume()
		select {
		case <-finished:
		case <-time.After(10 * time.Second):
			t.Error("older recovery sender did not finish")
		}
	})
	select {
	case <-paused:
	case <-finished:
		t.Fatalf("first recovery failed before its send reservation: %v", firstErr)
	case <-time.After(10 * time.Second):
		t.Fatal("first recovery did not reach the deterministic send boundary")
	}
	firstArtifacts := deliveryRecoveryFiles(t, first.RecoveryDirectory)

	// This isolated fixture positively reports an idle shell with no descendants.
	// A fresh user confirmation selects the next recovery; no actual provider runs.
	second := deliveryRecoveryTestPreview(t, f, in)
	if first.Generation != 1 || second.Generation != 2 || second.Confirmation == first.Confirmation || second.RecoveryDirectory == first.RecoveryDirectory {
		t.Fatalf("fresh recovery did not bind a distinct generation: first=%+v second=%+v", first, second)
	}
	deliveryRecoveryPreserveControl(t, in, second)
	current, err := WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, second.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if current.RunID != original.RunID || current.StartupRecovery == nil || current.StartupRecovery.Generation != 2 {
		t.Fatalf("recovery did not preserve run identity and select generation two: %+v", current)
	}
	f.D.Executable = func() (string, error) { return second.ControlExecutable.Path, nil }
	current = deliveryTestAccept(t, f, current)
	if current.Delivery.PermissionState != "recipient_confirmed_contract" {
		t.Fatal("second recovery did not accept its own runtime")
	}
	deliveryRecoveryUnchanged(t, firstArtifacts, "")
	before := deliveryRecoveryFiles(t, current.Paths.RunRoot)
	calls := map[string]int{}
	for _, command := range []string{"tab create", "agent start", "agent prompt"} {
		calls[command] = workflowTestCalls(t, f.workflowFixture, command)
	}
	resume()
	select {
	case <-finished:
	case <-time.After(10 * time.Second):
		t.Fatal("older recovery sender did not stop after replacement")
	}
	if firstErr == nil {
		t.Error("older sender did not reject its superseded recovery generation")
	}
	deliveryRecoveryUnchanged(t, before, "")
	if len(deliveryRecoveryFiles(t, current.Paths.RunRoot)) != len(before) {
		t.Error("older recovery sender added artifacts after replacement acceptance")
	}
	for command, count := range calls {
		if actual := workflowTestCalls(t, f.workflowFixture, command); actual != count {
			t.Errorf("older recovery sender replayed %q: before=%d after=%d", command, count, actual)
		}
	}

	// Reusing the exact human confirmation observes the same completed attempt.
	// It cannot create a third attempt or emit another provider effect.
	again, err := WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, second.Confirmation)
	if err != nil || again.StartupRecovery == nil || again.StartupRecovery.Generation != 2 {
		t.Fatalf("same confirmation did not return its existing recovery: %+v %v", again, err)
	}
	deliveryRecoveryUnchanged(t, before, "")
	for command, count := range calls {
		if actual := workflowTestCalls(t, f.workflowFixture, command); actual != count {
			t.Errorf("same confirmation replayed %q: before=%d after=%d", command, count, actual)
		}
	}
	statePath := filepath.Join(current.Paths.RunRoot, "state.json")
	if state, err := os.ReadFile(statePath); err != nil || !bytes.Equal(state, before[statePath]) {
		t.Fatalf("same confirmation changed accepted state: %v", err)
	}
}
