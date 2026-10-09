package taskrun

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func existingAcceptanceFixture(t *testing.T) (deliveryFixture, WorkflowRun, string, string, string) {
	t.Helper()
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{SchemaVersion: 1, Mode: workspace.DeliveryLocalEpic}})
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.File = writeAny(t, f.R.WorkspaceRoot, "selection-native-request.json", f.R)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "historical choice candidate")
	review := deliveryTestReview(t, f, "before-selection")
	o, err := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review)
	if err != nil {
		t.Fatal(err)
	}
	choice := acceptanceChoiceFile(t, f, o, "choice.json")
	binary := filepath.Join(f.R.WorkspaceRoot, "private-candidate")
	if err = os.WriteFile(binary, []byte("#!/bin/sh\nprintf 'candidate\\n'\n"), 0700); err != nil {
		t.Fatal(err)
	}
	return f, o, choice, review, binary
}

func TestAcceptanceSelectionMissingPointerIsDiagnosed(t *testing.T) {
	f, o, choice, _, _ := existingAcceptanceFixture(t)
	var err error
	o, err = WorkflowDeliverySelectAcceptance(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, choice)
	if err != nil {
		t.Fatal(err)
	}
	if err = workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error { s.Result.Delivery.AcceptanceSelection = nil; return nil }); err != nil {
		t.Fatal(err)
	}
	if _, err = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID); err == nil {
		t.Fatal("orphan native choice silently reverted policy")
	}
	run := traceRun(t, traceRead(t, f.D, f.R.WorkspaceRoot), o.RunID)
	if run.Coverage == "complete" {
		t.Fatal("orphan selection was neither diagnosed nor preserved as damaged evidence")
	}
}

func TestAcceptanceSelectionInterruptedConcurrentAndTampered(t *testing.T) {
	f, o, choice, _, _ := existingAcceptanceFixture(t)
	d := f.D
	d.Fault = func(at string) error {
		if at == "delivery_after_acceptance_selection_record" {
			return errors.New("lost publication reply")
		}
		return nil
	}
	if _, err := WorkflowDeliverySelectAcceptance(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, choice); err == nil {
		t.Fatal("interruption not exercised")
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || s.Result.Delivery.AcceptanceSelection != nil {
		t.Fatalf("partial selection published authority: %v", err)
	}
	path := filepath.Join(deliverySelectionDirectory(s), "selection.json")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := WorkflowDeliverySelectAcceptance(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, choice)
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	after, _ := os.ReadFile(path)
	if string(before) != string(after) {
		t.Fatal("retry replaced the original timestamp or choice")
	}
	s, err = workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Result.Delivery.Events) != 2 || len(s.Result.Delivery.Candidates) != 1 || s.Result.Delivery.Candidates[0].HumanQA != nil {
		t.Fatal("choice repeated or created candidate/QA")
	}
	if err = os.Remove(choice); err != nil {
		t.Fatal(err)
	}
	if _, err = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID); err != nil {
		t.Fatalf("managed choice depends on original input: %v", err)
	}
	if err = os.WriteFile(filepath.Join(deliverySelectionDirectory(s), "choice.json"), []byte("changed"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID); err == nil {
		t.Fatal("tampered managed choice accepted")
	}
	if _, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err == nil {
		t.Fatal("tampered choice permitted integration")
	}
}

func TestAcceptanceSelectionRecoveryPreservesInterveningWork(t *testing.T) {
	f, o, choice, _, _ := existingAcceptanceFixture(t)
	d := f.D
	d.Fault = func(at string) error {
		if at == "delivery_after_acceptance_selection_record" {
			return errors.New("lost publication reply")
		}
		return nil
	}
	if _, err := WorkflowDeliverySelectAcceptance(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, choice); err == nil {
		t.Fatal("fault not exercised")
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(deliverySelectionDirectory(s), "selection.json")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	report := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "selection-interrupted", PreviousEventSHA256: s.Result.Delivery.LastEventSHA256, Phase: "working", Summary: "Selection publication was interrupted.", Meaning: "The same owner is correcting the candidate while preserving the already requested policy choice.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	if _, err = WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "selection-interrupted.json", report)); err != nil {
		t.Fatal(err)
	}
	deliveryTestCandidate(t, f, "corrected after interrupted selection")
	o, err = WorkflowDeliverySelectAcceptance(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, choice)
	if err != nil {
		t.Fatalf("same policy choice cannot resume after intervening report and code correction: %v", err)
	}
	after, _ := os.ReadFile(path)
	if string(before) != string(after) || len(o.Delivery.Events) != 3 || o.DeliveryStatus.Acceptance.Outcome == "pass" || o.DeliveryStatus.QualifiedCandidate.Current {
		t.Fatal("recovery rewrote history or approved changed source")
	}
}

func TestAcceptanceSelectionRejectsStaleOrUnresolvedAuthority(t *testing.T) {
	for _, kind := range []string{"source", "agreement", "runtime", "context", "control", "attempt", "released"} {
		t.Run(kind, func(t *testing.T) {
			f, o, choice, _, _ := existingAcceptanceFixture(t)
			d := f.D
			contextPath := o.Paths.Context
			switch kind {
			case "source":
				deliveryTestCandidate(t, f, "changed after choice")
			case "agreement":
				var c DeliveryAcceptanceChoice
				if err := readValue(choice, 1<<20, &c); err != nil {
					t.Fatal(err)
				}
				c.PreviousAgreementSHA256 = "sha256:" + strings.Repeat("0", 64)
				choice = writeAny(t, f.R.WorkspaceRoot, "wrong-agreement.json", c)
			case "runtime":
				if err := os.WriteFile(filepath.Join(f.R.WorkspaceRoot, "actual-policy.json"), []byte("changed runtime"), 0600); err != nil {
					t.Fatal(err)
				}
			case "context":
				contextPath = filepath.Join(f.R.WorkspaceRoot, "unbound-context")
			case "control":
				d.Executable = func() (string, error) { return filepath.Join(f.R.WorkspaceRoot, "other-control"), nil }
			case "attempt", "released":
				if err := workflowUpdate(d, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
					if kind == "attempt" {
						s.Result.Delivery.Attempt.State = "attempted"
					} else {
						s.Result.Delivery.OwnershipRelease = &FileBinding{Locator: "/unrelated", SHA256: "sha256:" + strings.Repeat("0", 64)}
					}
					return nil
				}); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := WorkflowDeliverySelectAcceptance(d, f.R.WorkspaceRoot, o.RunID, contextPath, choice); err == nil {
				t.Fatal("stale/unresolved authority selected automatic acceptance")
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if s.Result.Delivery.AcceptanceSelection != nil {
				t.Fatal("rejected selection changed policy")
			}
		})
	}
}

func TestAcceptanceSelectionPreservesActualHumanVeto(t *testing.T) {
	for _, outcome := range []string{"fail", "blocked"} {
		t.Run(outcome, func(t *testing.T) {
			f, o, choice, review, binary := existingAcceptanceFixture(t)
			var err error
			o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, outcome, deliveryTestHuman(t, f, o, outcome))
			if err != nil {
				t.Fatal(err)
			}
			o, err = WorkflowDeliverySelectAcceptance(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, choice)
			if err != nil {
				t.Fatal(err)
			}
			o, err = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
			if err != nil {
				t.Fatal(err)
			}
			if o.DeliveryStatus.Acceptance.Outcome != outcome || o.DeliveryStatus.HumanJudgment == nil || o.DeliveryStatus.HumanJudgment.Outcome != outcome || len(o.Delivery.Candidates) != 2 {
				t.Fatalf("choice erased original human veto: %+v", o.DeliveryStatus)
			}
			if _, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err == nil {
				t.Fatal("human veto bypassed")
			}
		})
	}
}

func TestAcceptanceSelectionFailureAndUnknownCannotDeliver(t *testing.T) {
	for _, kind := range []string{"fail", "unknown"} {
		t.Run(kind, func(t *testing.T) {
			f, o, choice, review, binary := existingAcceptanceFixture(t)
			var err error
			o, err = WorkflowDeliverySelectAcceptance(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, choice)
			if err != nil {
				t.Fatal(err)
			}
			marker := filepath.Join(f.R.WorkspaceRoot, "actual-test-count")
			script := "#!/bin/sh\nprintf 'once\\n' >> " + ShellQuote(marker) + "\nexit 7\n"
			if err = os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
				t.Fatal(err)
			}
			d := f.D
			if kind == "unknown" {
				d.Fault = func(at string) error {
					if at == "delivery_after_verification_reservation" {
						return errors.New("process start unobserved")
					}
					return nil
				}
			}
			o, err = WorkflowDeliveryVerifyWithBinary(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
			if err == nil {
				t.Fatal("nonpassing acceptance qualified")
			}
			if _, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err == nil {
				t.Fatal("nonpassing acceptance integrated")
			}
			if len(o.Delivery.Candidates) != 1 || o.DeliveryStatus.Acceptance.Outcome == "pass" {
				t.Fatalf("old result became current pass: %+v", o.DeliveryStatus)
			}
			if kind == "unknown" {
				if _, err = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary); err == nil {
					t.Fatal("unknown test replayed")
				}
				if _, err = os.Stat(marker); !os.IsNotExist(err) {
					t.Fatal("unknown reservation executed on retry")
				}
			} else if raw, err := os.ReadFile(marker); err != nil || string(raw) != "once\n" {
				t.Fatalf("unexpected test executions %q %v", raw, err)
			}
		})
	}
}

func TestAcceptanceSelectionNativeEffectRevalidatesDecision(t *testing.T) {
	for _, kind := range []string{"changed", "missing", "unregistered"} {
		t.Run(kind, func(t *testing.T) {
			f, o, choice, review, binary := existingAcceptanceFixture(t)
			var err error
			o, err = WorkflowDeliverySelectAcceptance(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, choice)
			if err != nil {
				t.Fatal(err)
			}
			o, err = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
			if err != nil {
				t.Fatal(err)
			}
			c := o.Delivery.Candidates[len(o.Delivery.Candidates)-1]
			a, err := ValidateDeliveryAuthority(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult, *f.R.Delivery.Agreement)
			if err != nil {
				t.Fatal(err)
			}
			in := autoGateInput(a)
			preview, err := workspace.CheckTaskIntegration(f.D.Workspace, in)
			if err != nil || preview.Readback.Classification != "ready" {
				t.Fatalf("fixture not ready: %v %+v", err, preview.Readback)
			}
			var plan struct {
				Plan struct {
					SHA256 string `json:"sha256"`
				} `json:"plan"`
			}
			raw, err := workspace.MarshalTaskIntegrationReadback(preview.Readback)
			if err != nil {
				t.Fatal(err)
			}
			if err = json.Unmarshal(raw, &plan); err != nil || plan.Plan.SHA256 == "" {
				t.Fatalf("native confirmation missing: %v", err)
			}
			in.Apply = true
			in.Confirmation = plan.Plan.SHA256
			binding := *o.Delivery.AcceptanceSelection
			switch kind {
			case "changed":
				err = os.WriteFile(binding.Locator, []byte("changed selection"), 0600)
			case "missing":
				err = os.Remove(binding.Locator)
			case "unregistered":
				// Preserve a hash-valid decision and derived authorization but
				// remove its current native selection. Caller input grants nothing.
				err = workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error { s.Result.Delivery.AcceptanceSelection = nil; return nil })
			}
			if err != nil {
				t.Fatal(err)
			}
			if _, err = workspace.ApplyTaskIntegration(f.D.Workspace, in); err == nil {
				t.Fatal("native effect trusted stale/forged caller selection after ready check")
			}
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") != a.ExpectedParentOID {
				t.Fatal("invalid selection changed native Git target")
			}
		})
	}
}
