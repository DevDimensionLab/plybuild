package taskrun

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestContinuityJourneyRecoversLostReplyWithoutAnotherTransition(t *testing.T) {
	driver, err := filepath.Abs("../../test/continuity_journey.py")
	if err != nil {
		t.Fatal(err)
	}
	bytecodeBefore := continuityDriverBytecode(t, driver)
	t.Cleanup(func() {
		if !equal(bytecodeBefore, continuityDriverBytecode(t, driver)) {
			t.Error("importing the journey driver changed source-tree bytecode artifacts")
		}
	})
	current, old := deliveryCLIBinary(t), continuityOldBinary(t)
	before := continuitySourceBinary(t, continuityBeforeUpgradeSource, "PLY_CONTINUITY_BEFORE_UPGRADE_BINARY")
	f := newContinuityFixture(t, current, old, "")
	continuityUpgradeProvider(t, &f, before, "removed")
	manifest := writeAny(t, f.R.WorkspaceRoot, "human-journey.json", continuityJourneyManifest(t, f))
	// Lose only the wrapper's return after the real native command completed.
	// No native state or immutable proof is edited to manufacture recovery.
	program := `import importlib.util, pathlib, sys
sys.dont_write_bytecode = True
spec = importlib.util.spec_from_file_location("journey_driver", sys.argv[1])
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)
journey = module.Journey(pathlib.Path(sys.argv[2]))
native_call = journey.call
def lost_reply(binary, *args, **kwargs):
    result = native_call(binary, *args, **kwargs)
    if args[:3] == ("workflow", "execute", "continue") and "--check" not in args:
        raise RuntimeError("Synthetic lost native continuation response")
    return result
journey.call = lost_reply
try:
    journey.continue_delivery()
except RuntimeError as error:
    if str(error) == "Synthetic lost native continuation response":
        raise SystemExit(23)
    raise
raise SystemExit("fixture never lost the native response")
`
	command := exec.Command("python3", "-c", program, driver, manifest)
	command.Dir = f.CWD
	if output, err := command.CombinedOutput(); err == nil || command.ProcessState.ExitCode() != 23 {
		t.Fatalf("fixture did not interrupt the journey reply: %v\n%s", err, output)
	}
	if _, err := os.Stat(filepath.Join(f.R.WorkspaceRoot, "continued.json")); !os.IsNotExist(err) {
		t.Fatalf("lost reply fabricated a local continuation receipt: %v", err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, f.Run.RunID)
	if err != nil || deliveryContinuationGeneration(s) != 1 || s.Continuation == nil {
		t.Fatalf("native transition did not actually complete before the lost reply: %v", err)
	}
	proof, err := os.ReadFile(s.Continuation.Locator)
	if err != nil {
		t.Fatal(err)
	}
	command = exec.Command("python3", driver, "--fixture", manifest, "continue")
	command.Dir = f.CWD
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("same journey could not observe and recover its native transition: %v\n%s", err, output)
	}
	after, err := workflowRead(f.R.WorkspaceRoot, f.Run.RunID)
	if err != nil || deliveryContinuationGeneration(after) != 1 || !equal(s.Continuation, after.Continuation) {
		t.Fatalf("same journey retry published another continuation: before=%+v after=%+v err=%v", s.Continuation, after.Continuation, err)
	}
	retained, err := os.ReadFile(after.Continuation.Locator)
	if err != nil || !bytes.Equal(proof, retained) {
		t.Fatalf("journey retry changed the original immutable proof: %v", err)
	}
	if continuityCounter(t, f.Counter) != 2 || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 || workflowTestCalls(t, f.workflowFixture, "tab create") != 1 {
		t.Fatal("journey retry duplicated a Task, check, provider start or prompt")
	}
}

func continuityDriverBytecode(t *testing.T, driver string) map[string]string {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(filepath.Dir(driver), "__pycache__", "continuity_journey.*.pyc"))
	if err != nil {
		t.Fatal(err)
	}
	artifacts := map[string]string{}
	for _, path := range paths {
		artifacts[path] = hashFileTest(t, path)
	}
	return artifacts
}

func TestContinuityJourneyRefreshesOnlyObservedUnpublishedTransitions(t *testing.T) {
	driver, err := filepath.Abs("../../test/continuity_journey.py")
	if err != nil {
		t.Fatal(err)
	}
	current, old := deliveryCLIBinary(t), continuityOldBinary(t)
	before := continuitySourceBinary(t, continuityBeforeUpgradeSource, "PLY_CONTINUITY_BEFORE_UPGRADE_BINARY")
	current, err = filepath.EvalSymlinks(current)
	if err != nil {
		t.Fatal(err)
	}
	for _, publication := range []string{"none", "proof", "before_state", "unknown_artifact"} {
		t.Run(publication, func(t *testing.T) {
			f := newContinuityFixture(t, current, old, "")
			continuityUpgradeProvider(t, &f, before, "removed")
			manifest := writeAny(t, f.R.WorkspaceRoot, "human-journey.json", continuityJourneyManifest(t, f))
			var observation DeliveryRuntimeObservation
			if err := readValue(f.Upgrade.RuntimeObservation.Locator, 1<<20, &observation); err != nil {
				t.Fatal(err)
			}
			observedAt := time.Now().UTC().Add(-11 * time.Minute)
			observation.ObservedAtUTC = observedAt.Format(time.RFC3339Nano)
			selectedPath := writeAny(t, f.R.WorkspaceRoot, "selected-stale-runtime.json", observation)
			selected := FileBinding{selectedPath, hashFileTest(t, selectedPath)}
			selectionPath := writeAny(t, f.R.WorkspaceRoot, "continuation-runtime.json", selected)
			var preservedPath, preservedDigest string
			if publication != "none" {
				// The supported native transition is interrupted while this
				// observation is fresh. Later the real CLI sees an expired file.
				d := f.D
				d.Now = func() time.Time { return observedAt.Add(time.Second) }
				d.Executable = func() (string, error) { return current, nil }
				in := DeliveryContinuationInput{RunID: f.Run.RunID, ContextPath: f.Run.Paths.Context, ControlExecutable: Executable{current, hashFileTest(t, current)}, RuntimeEvidencePath: selectedPath}
				p, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in)
				if err != nil {
					t.Fatal(err)
				}
				// Supply the same immutable candidate copy that the CLI's
				// PreserveControl step supplies before invoking this native API.
				if err = privateDir(p.Directory); err != nil {
					t.Fatal(err)
				}
				control, err := os.ReadFile(current)
				if err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(p.ControlExecutable.Path, control, 0500); err != nil {
					t.Fatal(err)
				}
				if publication == "unknown_artifact" {
					preservedPath = filepath.Join(p.Directory, "unknown-effect")
					if err = os.WriteFile(preservedPath, []byte("An unobserved fixture effect must not be replayed.\n"), 0600); err != nil {
						t.Fatal(err)
					}
				} else {
					preservedPath = filepath.Join(p.Directory, "continuation.json")
					if publication == "before_state" {
						d.Files = continuationRejectProof{}
						preservedPath = filepath.Join(p.Directory, "before-state.json")
					} else {
						d.Fault = func(point string) error {
							if point == "delivery_after_continuation_proof" {
								return errors.New("synthetic interrupted continuation publication")
							}
							return nil
						}
					}
					if _, err = WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
						t.Fatal("native publication did not interrupt at the chosen boundary")
					}
				}
				preservedDigest = hashFileTest(t, preservedPath)
			}
			command := exec.Command("python3", driver, "--fixture", manifest, "continue")
			command.Dir = f.CWD
			output, commandErr := command.CombinedOutput()
			var afterSelection FileBinding
			if err := readValue(selectionPath, 1<<20, &afterSelection); err != nil {
				t.Fatal(err)
			}
			state, err := workflowRead(f.R.WorkspaceRoot, f.Run.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if publication == "unknown_artifact" {
				if commandErr == nil || state.Continuation != nil || afterSelection != selected {
					t.Fatalf("reserved transition regenerated evidence or retried an unknown effect: %v\n%s", commandErr, output)
				}
			} else {
				if commandErr != nil || deliveryContinuationGeneration(state) != 1 {
					t.Fatalf("journey could not continue the observed publication state: %v\n%s", commandErr, output)
				}
				if publication == "proof" && afterSelection != selected {
					t.Fatal("published proof recovery replaced its original observation")
				}
				if publication != "proof" && afterSelection == selected {
					t.Fatal("known unpublished transition did not refresh its expired observation")
				}
			}
			if hashFileTest(t, selectedPath) != selected.SHA256 || preservedPath != "" && hashFileTest(t, preservedPath) != preservedDigest {
				t.Fatal("journey rewrote immutable observation or interrupted publication evidence")
			}
			if continuityCounter(t, f.Counter) != 2 || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 || workflowTestCalls(t, f.workflowFixture, "tab create") != 1 {
				t.Fatal("observation expiry duplicated a provider or acceptance effect")
			}
		})
	}
}
