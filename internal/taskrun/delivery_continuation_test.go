package taskrun

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func continuationTestInput(t *testing.T, f deliveryFixture, o WorkflowRun) (Dependencies, DeliveryContinuationInput, DeliveryContinuationPreview) {
	t.Helper()
	source := filepath.Join(f.R.WorkspaceRoot, "compatible-installed-ply")
	if err := os.WriteFile(source, []byte("#!/bin/sh\n# compatible reader fixture\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	d := f.D
	d.Executable = func() (string, error) { return source, nil }
	in := DeliveryContinuationInput{RunID: o.RunID, ContextPath: o.Paths.Context, ControlExecutable: Executable{source, hashFileTest(t, source)}, RuntimeEvidencePath: deliveryTestRuntimeObservation(t, f, o)}
	p, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in)
	if err != nil {
		t.Fatal(err)
	}
	if err = privateDir(p.Directory); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(p.ControlExecutable.Path, raw, 0500); err != nil {
		t.Fatal(err)
	}
	return d, in, p
}

func TestDeliveryContinuationEquivalentObservationKeepsOneControl(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	d, in, p := continuationTestInput(t, f, o)
	if err := os.Remove(f.R.Runtime.Executable.Path); err != nil {
		t.Fatal(err)
	}
	first, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	var observation DeliveryRuntimeObservation
	if err = readValue(in.RuntimeEvidencePath, 1<<20, &observation); err != nil {
		t.Fatal(err)
	}
	observation.ObservedAtUTC = d.Now().Add(time.Second).UTC().Format(time.RFC3339Nano)
	in.RuntimeEvidencePath = writeAny(t, f.R.WorkspaceRoot, "fresh-equivalent-runtime.json", observation)
	again, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in)
	if err != nil || again.State != "existing" || again.Generation != 1 || again.ControlExecutable != first.ControlExecutable || !equal(again.Proof, first.Proof) || !equal(again.RuntimeObservation, first.RuntimeObservation) {
		t.Fatalf("equivalent current observation duplicated its control or changed readback: %+v %v", again, err)
	}
	observation.DeliveryPermission.PermissionConfirmed = false
	in.RuntimeEvidencePath = writeAny(t, f.R.WorkspaceRoot, "unconfirmed-runtime.json", observation)
	if _, err = WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in); err == nil {
		t.Fatal("existing proof bypassed an actual unconfirmed current authority observation")
	}
}

func TestDeliveryContinuationConcurrentAndInterruptedPublication(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	d, in, p := continuationTestInput(t, f, o)
	fault := d
	fault.Fault = func(point string) error {
		if point == "delivery_after_continuation_proof" {
			return errors.New("simulated lost publication reply")
		}
		return nil
	}
	if _, err := WorkflowContinueDelivery(fault, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
		t.Fatal("fault did not interrupt publication")
	}
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation)
			if err == nil && result.State != "continued" && result.State != "existing" {
				err = errors.New("unexpected continuation state")
			}
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
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || deliveryContinuationGeneration(s) != 1 {
		t.Fatalf("duplicate continuation generation: %d %v", deliveryContinuationGeneration(s), err)
	}
	entries, err := os.ReadDir(filepath.Dir(p.Directory))
	if err != nil || len(entries) != 1 {
		t.Fatalf("duplicate continuation slots: %v %v", entries, err)
	}
}

func TestDeliveryContinuationRecoversPublishedProofAfterObservationExpires(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	d, in, p := continuationTestInput(t, f, o)
	if err := os.Remove(f.R.Runtime.Executable.Path); err != nil {
		t.Fatal(err)
	}
	fault := d
	fault.Fault = func(point string) error {
		if point == "delivery_after_continuation_proof" {
			return errors.New("simulated interrupted proof publication")
		}
		return nil
	}
	if _, err := WorkflowContinueDelivery(fault, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
		t.Fatal("publication fault did not interrupt the transition")
	}
	proof := filepath.Join(p.Directory, "continuation.json")
	before, err := os.ReadFile(proof)
	if err != nil {
		t.Fatal(err)
	}
	now := d.Now()
	d.Now = func() time.Time { return now.Add(11 * time.Minute) }
	recovered, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation)
	if err != nil || recovered.State != "continued" || recovered.Generation != 1 {
		t.Fatalf("durably published transition was stranded by expired observation: %+v %v", recovered, err)
	}
	after, err := os.ReadFile(proof)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("recovery changed the original transition proof: %v", err)
	}
	if n := workflowTestCalls(t, f.workflowFixture, "agent start"); n != 1 {
		t.Fatalf("proof recovery restarted provider: %d", n)
	}
}

func TestDeliveryContinuationRecoversPublishedProofAfterIncompleteReport(t *testing.T) {
	testDeliveryContinuationRecoveryAfterReport(t, false)
}

type continuationRejectProof struct{ systemFileSystem }

func (continuationRejectProof) WriteOnce(path string, raw []byte) error {
	if filepath.Base(path) == "continuation.json" {
		return errors.New("simulated publication failure before proof")
	}
	return (systemFileSystem{}).WriteOnce(path, raw)
}

func TestDeliveryContinuationRecoversBeforeStateAfterIncompleteReport(t *testing.T) {
	testDeliveryContinuationRecoveryAfterReport(t, true)
}

func testDeliveryContinuationRecoveryAfterReport(t *testing.T, beforeProof bool) {
	t.Helper()
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	d, in, p := continuationTestInput(t, f, o)
	if err := os.Remove(f.R.Runtime.Executable.Path); err != nil {
		t.Fatal(err)
	}
	fault := d
	if beforeProof {
		fault.Files = continuationRejectProof{}
	}
	fault.Fault = func(point string) error {
		if point == "delivery_after_continuation_proof" {
			return errors.New("simulated interrupted proof publication")
		}
		return nil
	}
	if _, err := WorkflowContinueDelivery(fault, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
		t.Fatal("publication fault did not interrupt the transition")
	}
	proofPath := filepath.Join(p.Directory, "continuation.json")
	if beforeProof {
		proofPath = filepath.Join(p.Directory, "before-state.json")
	}
	before, err := os.ReadFile(proofPath)
	if err != nil {
		t.Fatal(err)
	}
	report := providerUpgradeIncomplete(t, f, o, "stopped")
	if _, err = WorkflowDeliveryIncompleteReport(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "interrupted-continuation-report.json", report)); err != nil {
		t.Fatal(err)
	}
	p, err = WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation); err != nil {
		t.Fatalf("report stranded a durably published transition: %v", err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || len(s.Result.Delivery.Events) != 1 || s.Result.Delivery.Events[0].ID != report.EventID {
		t.Fatalf("recovery lost the intervening report: %+v %v", s.Result.Delivery, err)
	}
	after, err := os.ReadFile(proofPath)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("recovery rewrote the original proof: %v", err)
	}
}

func TestDeliveryContinuationRejectsUnknownVerificationAndStaleSession(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	d, in, _ := continuationTestInput(t, f, o)
	model := workflowProviderDocument(t, f.Model)
	model["agent_session_id"] = "different-session"
	writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.Model), model)
	if _, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in); err == nil || !strings.Contains(err.Error(), "delivery_continuation_session") {
		t.Fatalf("stale session granted continuation: %v", err)
	}
	model["agent_session_id"] = o.Transport.AgentSessionID
	writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.Model), model)
	deliveryTestCandidate(t, f, "unknown command")
	fault := f.D
	fault.Fault = func(point string) error {
		if point == "delivery_after_verification_reservation" {
			return errors.New("interrupted before known execution")
		}
		return nil
	}
	if _, err := WorkflowDeliveryVerify(fault, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err == nil {
		t.Fatal("fault did not reserve unknown command")
	}
	if _, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in); err == nil || !strings.Contains(err.Error(), "delivery_effect_unknown") {
		t.Fatalf("unknown effect was replayable: %v", err)
	}
}

func TestDeliveryContinuationRejectsTamperedPriorControlAndAuthority(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	d, in, p := continuationTestInput(t, f, o)
	if _, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f.R.Runtime.PlyExecutable.Path, []byte("#!/bin/sh\nexit 9\n"), 0700); err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = workflowEffectiveRuntime(s); err != nil {
		t.Fatalf("historical Ply upgrade blocked continuation: %v", err)
	}
	s.Request.SchemaVersion = 9
	if err = deliveryContinuationAuthority(s); err == nil || !strings.Contains(err.Error(), "delivery_continuation_unsupported") {
		t.Fatalf("unsupported contract granted continuation: %v", err)
	}
}

func TestPublishedContinuationRecoveryChangedCaller(t *testing.T) {
	for _, change := range []string{"path", "bytes"} {
		t.Run(change, func(t *testing.T) {
			f := newDeliveryFixture(t, "codex")
			o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
			d, in, p := continuationTestInput(t, f, o)
			beforeState, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			prompts := workflowTestCalls(t, f.workflowFixture, "agent prompt")
			tabs := workflowTestCalls(t, f.workflowFixture, "tab create")
			fault := d
			fault.Fault = func(point string) error {
				if point == "delivery_after_continuation_proof" {
					return errors.New("lost proof publication reply")
				}
				return nil
			}
			if _, err := WorkflowContinueDelivery(fault, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
				t.Fatal("fault did not interrupt")
			}
			proof := filepath.Join(p.Directory, "continuation.json")
			original, err := os.ReadFile(proof)
			if err != nil {
				t.Fatal(err)
			}
			frozen := map[string][]byte{}
			for _, file := range []string{p.ControlExecutable.Path, filepath.Join(p.Directory, "before-state.json"), o.Paths.Context, filepath.Join(o.Paths.RunRoot, "acceptance.json"), filepath.Join(o.Paths.RunRoot, "mandate.json")} {
				data, e := os.ReadFile(file)
				if e != nil {
					t.Fatal(e)
				}
				frozen[file] = data
			}
			path := in.ControlExecutable.Path
			raw, _ := os.ReadFile(path)
			if change == "path" {
				path += "-new"
			} else {
				raw = append(raw, []byte("# upgraded reader\n")...)
			}
			if err = os.WriteFile(path, raw, 0700); err != nil {
				t.Fatal(err)
			}
			d.Executable = func() (string, error) { return path, nil }
			in.ControlExecutable = Executable{path, hashFileTest(t, path)}
			preview, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in)
			if err != nil {
				t.Fatalf("installed Ply preview: %v", err)
			}
			if preview.ControlExecutable != p.ControlExecutable || preview.Proof == nil {
				t.Fatalf("lost original proof/control: %+v", preview)
			}
			// Another upgrade between preview and apply is independent of the
			// published proof's confirmation. Rebind the actual current caller.
			if change == "path" {
				path += "-next"
			} else {
				raw = append(raw, []byte("# next installed reader\n")...)
			}
			if err = os.WriteFile(path, raw, 0700); err != nil {
				t.Fatal(err)
			}
			in.ControlExecutable = Executable{path, hashFileTest(t, path)}
			var wg sync.WaitGroup
			errs := make(chan error, 2)
			for i := 0; i < 2; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, e := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, preview.Confirmation)
					errs <- e
				}()
			}
			wg.Wait()
			close(errs)
			successes := 0
			for e := range errs {
				if e == nil {
					successes++
				} else {
					t.Fatal(e)
				}
			}
			if successes != 2 {
				t.Fatalf("unexpected successful recoveries/reuses: %d", successes)
			}
			// A lost reply after pointer publication reuses the old confirmation.
			_, err = WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, preview.Confirmation)
			if err != nil {
				t.Fatalf("post-pointer repeat: %v", err)
			}
			again, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in)
			if err != nil {
				t.Fatal(err)
			}
			if again.State != "existing" || again.Generation != 1 || again.Proof == nil {
				t.Fatalf("registered proof was not reused: %+v", again)
			}
			state, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil || deliveryContinuationGeneration(state) != 1 || !equal(state.Result.Delivery, beforeState.Result.Delivery) {
				t.Fatalf("recovery changed generation/QA/authority: %v", err)
			}
			entries, err := os.ReadDir(filepath.Dir(p.Directory))
			if err != nil || len(entries) != 1 {
				t.Fatal("recovery published another generation")
			}
			if workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != prompts || workflowTestCalls(t, f.workflowFixture, "tab create") != tabs {
				t.Fatal("recovery caused new native effects")
			}
			for file, want := range frozen {
				got, e := os.ReadFile(file)
				if e != nil || !bytes.Equal(got, want) {
					t.Fatalf("original artifact changed: %s %v", file, e)
				}
			}
			after, err := os.ReadFile(proof)
			if err != nil || !bytes.Equal(original, after) {
				t.Fatal("proof bytes changed")
			}
		})
	}
}

func TestPublishedContinuationRecoveryRejectsDifferentRuntimeBinding(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	d, in, p := continuationTestInput(t, f, o)
	fault := d
	fault.Fault = func(point string) error {
		if point == "delivery_after_continuation_proof" {
			return errors.New("interrupted proof")
		}
		return nil
	}
	if _, err := WorkflowContinueDelivery(fault, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
		t.Fatal("fault did not interrupt")
	}
	var observation DeliveryRuntimeObservation
	if err := readValue(in.RuntimeEvidencePath, 1<<20, &observation); err != nil {
		t.Fatal(err)
	}
	observation.ObservedAtUTC = d.Now().Add(time.Second).UTC().Format(time.RFC3339Nano)
	in.RuntimeEvidencePath = writeAny(t, f.R.WorkspaceRoot, "different-valid-runtime.json", observation)
	if _, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in); err == nil {
		t.Fatal("pending recovery accepted a different valid runtime binding")
	}
}

func TestPublishedContinuationRecoveryGuardsUnderLock(t *testing.T) {
	for _, registered := range []bool{false, true} {
		for _, guard := range []string{"proof", "acceptance", "control", "before-state", "context", "task", "owner", "session", "policy", "runtime", "released", "unknown", "caller-path", "caller-bytes"} {
			name := "pending/" + guard
			if registered {
				name = "registered/" + guard
			}
			// Registered controls need not remain on disk; pending publication
			// still checks the control it will adopt.
			if registered && guard == "control" {
				continue
			}
			t.Run(name, func(t *testing.T) {
				f := newDeliveryFixture(t, "codex")
				o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
				d, in, p := continuationTestInput(t, f, o)
				fault := d
				fault.Fault = func(point string) error {
					if point == "delivery_after_continuation_proof" {
						return errors.New("interrupted proof")
					}
					return nil
				}
				if registered {
					if _, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(in.ControlExecutable.Path, []byte("#!/bin/sh\n# upgraded reader\nexit 0\n"), 0700); err != nil {
						t.Fatal(err)
					}
					in.ControlExecutable.SHA256 = hashFileTest(t, in.ControlExecutable.Path)
				} else if _, err := WorkflowContinueDelivery(fault, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
					t.Fatal("fault did not interrupt")
				}
				// Preview is valid, then the dependency changes before apply's locked read.
				preview, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in)
				if err != nil {
					t.Fatal(err)
				}
				state, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
				if err != nil {
					t.Fatal(err)
				}
				switch guard {
				case "proof", "before-state", "acceptance":
					name := "continuation.json"
					if guard == "before-state" {
						name = "before-state.json"
					}
					path := filepath.Join(p.Directory, name)
					if guard == "acceptance" {
						path = filepath.Join(o.Paths.RunRoot, "acceptance.json")
					}
					raw, err := os.ReadFile(path)
					if err != nil {
						t.Fatal(err)
					}
					// Valid JSON still has to match the original binding/confirmation.
					if err = os.WriteFile(path, append(raw, ' '), 0600); err != nil {
						t.Fatal(err)
					}
					if guard == "proof" { // A semantic proof edit, not just a serialization change.
						var record deliveryContinuationRecord
						if err = readValue(path, 1<<20, &record); err != nil {
							t.Fatal(err)
						}
						record.SourceControl.SHA256 = "sha256:" + strings.Repeat("0", 64)
						writeAny(t, p.Directory, name, record)
					}
				case "control":
					if err = os.Chmod(p.ControlExecutable.Path, 0700); err != nil {
						t.Fatal(err)
					}
					if err = os.WriteFile(p.ControlExecutable.Path, []byte("#!/bin/sh\n# tampered control\n"), 0500); err != nil {
						t.Fatal(err)
					}
				case "context":
					in.ContextPath += "-wrong"
				case "task":
					d.CWD = func() (string, error) { return f.R.WorkspaceRoot, nil }
				case "owner", "session":
					model := workflowProviderDocument(t, f.Model)
					if guard == "session" {
						model["agent_session_id"] = "other-session"
					} else {
						model["foreground_cwd"] = f.R.WorkspaceRoot
						model["cwd"] = f.R.WorkspaceRoot
					}
					writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.Model), model)
				case "policy":
					v := "sha256:" + strings.Repeat("0", 64)
					state.Result.Delivery.ActualPolicySHA256 = &v
				case "runtime":
					var observation DeliveryRuntimeObservation
					if err = readValue(in.RuntimeEvidencePath, 1<<20, &observation); err != nil {
						t.Fatal(err)
					}
					v := "other-session"
					observation.RuntimeClaim.NativeSessionID = &v
					writeAny(t, filepath.Dir(in.RuntimeEvidencePath), filepath.Base(in.RuntimeEvidencePath), observation)
				case "released":
					state.Result.Delivery.OwnershipRelease = &FileBinding{}
				case "unknown":
					state.Result.Delivery.Attempt = &DeliveryAttempt{Kind: "integration", State: "attempted", ID: "unknown-effect"}
				case "caller-path":
					in.ControlExecutable.Path += "-not-caller"
				case "caller-bytes":
					if err = os.WriteFile(in.ControlExecutable.Path, []byte("#!/bin/sh\n# changed since caller binding\n"), 0700); err != nil {
						t.Fatal(err)
					}
				}
				if guard == "policy" || guard == "released" || guard == "unknown" {
					if err = workflowSave(d, state); err != nil {
						t.Fatal(err)
					}
				}
				if _, err = WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, preview.Confirmation); err == nil {
					t.Fatalf("locked recovery bypassed %s", guard)
				}
				state, err = workflowRead(f.R.WorkspaceRoot, o.RunID)
				if err == nil && !registered && state.Continuation != nil {
					t.Fatal("rejected recovery published pointer")
				}
			})
		}
	}
}

func TestPublishedContinuationRecoveryUsesOriginalExpiredRuntime(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	d, in, p := continuationTestInput(t, f, o)
	fault := d
	fault.Fault = func(point string) error {
		if point == "delivery_after_continuation_proof" {
			return errors.New("interrupted proof")
		}
		return nil
	}
	if _, err := WorkflowContinueDelivery(fault, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
		t.Fatal("fault did not interrupt")
	}
	if err := os.Remove(f.R.Runtime.Executable.Path); err != nil {
		t.Fatal(err)
	}
	originalRuntime := in.RuntimeEvidencePath
	in.RuntimeEvidencePath = ""
	path := in.ControlExecutable.Path
	if err := os.WriteFile(path, []byte("#!/bin/sh\n# new installed reader\n"), 0700); err != nil {
		t.Fatal(err)
	}
	in.ControlExecutable.SHA256 = hashFileTest(t, path)
	now := d.Now()
	d.Now = func() time.Time { return now.Add(11 * time.Minute) }
	preview, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in)
	if err != nil || preview.RuntimeObservation == nil || preview.RuntimeObservation.Locator != originalRuntime || preview.ControlExecutable != p.ControlExecutable {
		t.Fatalf("original runtime/control lost: %+v %v", preview, err)
	}
	result, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, preview.Confirmation)
	if err != nil || result.State != "continued" || result.Generation != 1 {
		t.Fatalf("recovery: %+v %v", result, err)
	}
}
