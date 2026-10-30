package taskrun

import (
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
)

func TestRegisteredContinuationReusesAcrossCLIUpgrades(t *testing.T) {
	for _, change := range []string{"path", "bytes"} {
		for _, launcher := range []string{"present", "missing"} {
			t.Run(change+"/"+launcher, func(t *testing.T) {
				f := newDeliveryFixture(t, "codex")
				o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
				d, in, p := continuationTestInput(t, f, o)
				first, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation)
				if err != nil {
					t.Fatal(err)
				}
				if launcher == "missing" {
					if err = os.Remove(f.R.Runtime.Executable.Path); err != nil {
						t.Fatal(err)
					}
				}
				if err = os.Remove(first.ControlExecutable.Path); err != nil {
					t.Fatal(err)
				}
				original := deliveryRecoveryFiles(t, o.Paths.RunRoot)
				raw, err := os.ReadFile(in.ControlExecutable.Path)
				if err != nil {
					t.Fatal(err)
				}
				path := in.ControlExecutable.Path
				if change == "path" {
					path += "-new"
				} else {
					raw = append(raw, []byte("# upgrade\n")...)
				}
				if err = os.WriteFile(path, raw, 0700); err != nil {
					t.Fatal(err)
				}
				in.ControlExecutable = Executable{path, hashFileTest(t, path)}
				in.RuntimeEvidencePath = ""
				d.Executable = func() (string, error) { return path, nil }
				preview, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in)
				if err != nil || preview.State != "existing" {
					t.Fatalf("CLI-only upgrade: %+v %v", preview, err)
				}
				// Both repetition and locked concurrent apply reuse the original proof.
				var wg sync.WaitGroup
				for i := 0; i < 2; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						got, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, preview.Confirmation)
						if err != nil || got.State != "existing" || got.Generation != first.Generation || got.Directory != first.Directory || got.ControlExecutable != first.ControlExecutable || !reflect.DeepEqual(got.Proof, first.Proof) || got.Confirmation != first.Confirmation || !reflect.DeepEqual(got.RuntimeObservation, first.RuntimeObservation) {
							t.Errorf("registered reuse changed identity: %+v %v", got, err)
						}
					}()
				}
				wg.Wait()
				if !reflect.DeepEqual(original, deliveryRecoveryFiles(t, o.Paths.RunRoot)) {
					t.Fatal("reuse changed original artifacts or created a generation")
				}
				entries, err := os.ReadDir(filepath.Dir(first.Directory))
				if err != nil || len(entries) != 1 {
					t.Fatalf("generation directories: %v %v", entries, err)
				}
				for _, call := range []string{"tab create", "agent start", "agent prompt"} {
					if n := workflowTestCalls(t, f.workflowFixture, call); n != 1 {
						t.Fatalf("%s count: %d", call, n)
					}
				}
			})
		}
	}
}

func TestRegisteredContinuationValidatesExplicitRuntimeEvidence(t *testing.T) {
	for _, field := range []string{"policy", "session", "acceptance", "context", "run"} {
		t.Run(field, func(t *testing.T) {
			f := newDeliveryFixture(t, "codex")
			o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
			d, in, p := continuationTestInput(t, f, o)
			if _, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation); err != nil {
				t.Fatal(err)
			}
			var observation DeliveryRuntimeObservation
			if err := readValue(in.RuntimeEvidencePath, 1<<20, &observation); err != nil {
				t.Fatal(err)
			}
			wrong := "wrong"
			switch field {
			case "policy":
				observation.RuntimeClaim.EffectivePolicySHA256 = &wrong
			case "session":
				observation.RuntimeClaim.NativeSessionID = &wrong
			case "acceptance":
				observation.Acceptance.SHA256 = wrong
			case "context":
				observation.Context.SHA256 = wrong
			case "run":
				observation.RunID = wrong
			}
			in.RuntimeEvidencePath = writeAny(t, f.R.WorkspaceRoot, "invalid-runtime.json", observation)
			if _, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in); err == nil {
				t.Fatal("existing bypassed explicit invalid runtime evidence")
			}
		})
	}
}

func TestRegisteredContinuationAddsMissingRuntimeEvidence(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	d, in, _ := continuationTestInput(t, f, o)
	evidence := in.RuntimeEvidencePath
	in.RuntimeEvidencePath = ""
	p, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in)
	if err != nil {
		t.Fatal(err)
	}
	first, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation)
	if err != nil || first.RuntimeObservation != nil {
		t.Fatalf("first: %+v %v", first, err)
	}
	original := deliveryRecoveryFiles(t, first.Directory)
	if err = os.Remove(f.R.Runtime.Executable.Path); err != nil {
		t.Fatal(err)
	}
	in.RuntimeEvidencePath = evidence
	p, err = WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in)
	if err != nil || p.State != "ready" || p.Generation != 2 {
		t.Fatalf("needed runtime registration: %+v %v", p, err)
	}
	raw, _ := os.ReadFile(in.ControlExecutable.Path)
	if err = privateDir(p.Directory); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(p.ControlExecutable.Path, raw, 0500); err != nil {
		t.Fatal(err)
	}
	second, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation)
	if err != nil || second.RuntimeObservation == nil {
		t.Fatalf("registration: %+v %v", second, err)
	}
	if !reflect.DeepEqual(original, deliveryRecoveryFiles(t, first.Directory)) {
		t.Fatal("registration overwrote prior history")
	}
}
