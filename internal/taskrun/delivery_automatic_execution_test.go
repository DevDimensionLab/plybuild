package taskrun

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestAutomaticVerificationStartFailureHasNoInventedExit(t *testing.T) {
	r := deliveryVerificationReceipt{CWD: filepath.Join(t.TempDir(), "missing-directory"), Automatic: &workflowhandoff.AutomaticVerificationInstructions{Policy: workspace.DeliveryAcceptancePolicy{TimeoutSeconds: 2}, ExecutedArgv: []string{"/bin/sh", "missing-script"}}}
	_, _, stop := runDeliveryVerification(Dependencies{Now: time.Now}, &r)
	defer stop()
	if r.Outcome != "blocked" || r.Exit != nil || r.Error == "" || r.FinishedAt == "" {
		t.Fatalf("process start failure invented a result: %+v", r)
	}
}

func TestAutomaticVerificationBlocksUnretainableBinaryBeforeExecution(t *testing.T) {
	f, o, binary, review := autoGateFixture(t, false, "")
	if err := os.Truncate(binary, (64<<20)+1); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(f.R.WorkspaceRoot, "oversized-test-started")
	if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nprintf started > "+ShellQuote(marker)+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	blocked, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
	if err == nil || !strings.Contains(err.Error(), "64 MiB") || len(blocked.Delivery.Candidates) != 1 || blocked.DeliveryStatus.Acceptance.Current {
		t.Fatalf("unretainable binary did not block before execution: %v", err)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("oversized input started a test: %v", err)
	}
}

func TestAutomaticVerificationMissingInputsCannotReusePriorPass(t *testing.T) {
	for _, missing := range []string{"script", "binary"} {
		t.Run(missing, func(t *testing.T) {
			f, o, binary, review := autoGateFixture(t, false, "")
			path := f.AcceptancePath
			if missing == "binary" {
				path = binary
			}
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			blocked, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
			if err == nil || len(blocked.Delivery.Candidates) != 1 || blocked.DeliveryStatus.Acceptance.Outcome != "blocked" || blocked.DeliveryStatus.Acceptance.Current {
				t.Fatalf("missing %s reused prior evidence: %v %+v", missing, err, blocked.DeliveryStatus.Acceptance)
			}
			if _, err := WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err == nil {
				t.Fatalf("missing %s did not block integration", missing)
			}
		})
	}
}

type automaticPublicationFailure struct {
	systemFileSystem
	basename string
	failed   bool
}

func (f *automaticPublicationFailure) WriteOnce(path string, raw []byte) error {
	if !f.failed && filepath.Base(path) == f.basename {
		f.failed = true
		return errors.New("fixture artifact publication failure")
	}
	return f.systemFileSystem.WriteOnce(path, raw)
}

func TestAutomaticVerificationRecoversCompletedOutputPublication(t *testing.T) {
	f, o, binary, review := autoGateFixture(t, false, "")
	counter := filepath.Join(f.R.WorkspaceRoot, "execution-count")
	script := "#!/bin/sh\nprintf 'run\\n' >> " + ShellQuote(counter) + "\nprintf 'preserved result\\n'\n"
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	f.D.Files = &automaticPublicationFailure{basename: "stdout.txt"}
	interrupted, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
	if err == nil || interrupted.Delivery.Attempt == nil || len(interrupted.Delivery.Candidates) != 1 {
		t.Fatalf("fixture did not interrupt completed output publication: %v %+v", err, interrupted.Delivery)
	}
	f.D.Files = nil
	recovered, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
	if err != nil || recovered.DeliveryStatus.Verification.Outcome != "passed" || recovered.Delivery.Attempt.ID != interrupted.Delivery.Attempt.ID || len(recovered.Delivery.Candidates) != 2 {
		t.Fatalf("finished result could not be registered without repeating execution: %v %+v", err, recovered.Delivery)
	}
	raw, err := os.ReadFile(counter)
	if err != nil || string(raw) != "run\n" {
		t.Fatalf("completed test was repeated: %q %v", raw, err)
	}
	raw, err = os.ReadFile(filepath.Join(recovered.Delivery.Attempt.Path, "stdout.txt"))
	if err != nil || string(raw) != "preserved result\n" {
		t.Fatalf("actual output was lost: %q %v", raw, err)
	}
}

func TestAutomaticVerificationRetriesUnstartedLargeBinarySnapshot(t *testing.T) {
	f, o, binary, review := autoGateFixture(t, false, "")
	large := "#!/bin/sh\nprintf 'fixture candidate\\n'\n#" + strings.Repeat("x", 9<<20) + "\n"
	if err := os.WriteFile(binary, []byte(large), 0700); err != nil {
		t.Fatal(err)
	}
	f.D.Files = &automaticPublicationFailure{basename: "attempt.json"}
	if _, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary); err == nil {
		t.Fatal("fixture did not interrupt reservation publication before execution")
	}
	f.D.Files = nil
	recovered, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
	if err != nil || recovered.DeliveryStatus.Verification.Outcome != "passed" || len(recovered.Delivery.Candidates) != 2 {
		t.Fatalf("unstarted attempt could not reuse its existing binary snapshot: %v %+v", err, recovered.Delivery)
	}
}

func TestAutomaticVerificationExecutesReservedScript(t *testing.T) {
	f, o, binary, review := autoGateFixture(t, false, "")
	failedScript := []byte("#!/bin/sh\nexit 7\n")
	if err := os.WriteFile(f.AcceptancePath, failedScript, 0600); err != nil {
		t.Fatal(err)
	}
	backup := filepath.Join(f.R.WorkspaceRoot, "reserved-script")
	if err := os.WriteFile(backup, failedScript, 0600); err != nil {
		t.Fatal(err)
	}
	f.D.Fault = func(stage string) error {
		if stage == "delivery_after_verification_reservation" {
			// Replacement restores the original bytes before completion. A
			// post-execution hash alone cannot detect this different execution.
			return os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\ncp "+ShellQuote(backup)+" "+ShellQuote(f.AcceptancePath)+"\nexit 0\n"), 0600)
		}
		return nil
	}
	result, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
	if err == nil || result.DeliveryStatus.Verification.Exit == nil || *result.DeliveryStatus.Verification.Exit != 7 || len(result.Delivery.Candidates) != 1 {
		t.Fatalf("verification executed mutable replacement instead of reserved script: err=%v status=%+v candidates=%d", err, result.DeliveryStatus.Verification, len(result.Delivery.Candidates))
	}
}
