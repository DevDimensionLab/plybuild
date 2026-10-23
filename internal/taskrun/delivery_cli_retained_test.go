package taskrun

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/clirelease"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// TD125516B: runtime/session observations are simulated; qualification, Git,
// registry and executable bytes belong only to this disposable workspace.
func retainedInstallFixture(t *testing.T, binary string) (deliveryFixture, WorkflowRun) {
	t.Helper()
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic coordinator", TimeoutSeconds: 30}}
	f := receiptFixture(t, "codex", deliveryFixtureOptions{Agreement: &a})
	reviewRoundPatchHerdr(t, &f)
	parallelRuntime(t, &f)
	o := workflowTestStart(t, f.workflowFixture)
	o, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}, true)
	if err != nil {
		t.Fatal(err)
	}
	reviewRoundCandidate(t, f, "qualified retained source")
	if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\n\"$PLY_CANDIDATE_BINARY\" version --format json\n"), 0600); err != nil {
		t.Fatal(err)
	}
	o = parallelQualify(t, f, o, binary)
	o, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if !parallelOwnerHandoff(err) || o.Delivery.Phase != "closing" {
		t.Fatalf("fixture: %v %+v", err, o.Delivery)
	}
	return f, o
}

func TestQualifiedCLIRetainedSourcePublicRoute(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := retainedInstallFixture(t, binary)
	dest := filepath.Join(f.R.WorkspaceRoot, "installed-ply")
	before := branchFixtureState(t, f.R.WorkspaceRoot)
	cmd := exec.Command(binary, "workflow", "execute", "install-cli", o.RunID, "--destination", dest, "--coordinator", "synthetic authorized coordinator", "--allow-retained-source", "--check", "--format", "json")
	cmd.Dir = f.Parent
	raw, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("TD125516B explicit cleanup-only preview: %v\n%s", err, raw)
	}
	var result struct {
		State                string `json:"state"`
		SourceCleanupPending bool   `json:"source_cleanup_pending"`
	}
	if err := json.Unmarshal(raw, &result); err != nil || result.State != "ready" || !result.SourceCleanupPending {
		t.Fatalf("cleanup readback: %s %v", raw, err)
	}
	if !equal(before, branchFixtureState(t, f.R.WorkspaceRoot)) {
		t.Fatal("preview wrote native state")
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Fatal("preview installed destination")
	}
}

// Snapshot mutations below are injected test observations, never private-store
// patches to the live workspace or claims of human/product approval.
type retainedRegistryObservation struct {
	workspace.WorkItemStore
	change func(*workspace.WorkItemRegistry)
}

func (s retainedRegistryObservation) Snapshot(root string) (workspace.WorkItemRegistry, error) {
	r, err := s.WorkItemStore.Snapshot(root)
	if err == nil {
		s.change(&r)
	}
	return r, err
}

func TestQualifiedCLIRetainedSourceIsolationAndGuards(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := retainedInstallFixture(t, binary)
	root, c := f.R.WorkspaceRoot, o.Delivery.Candidates[0]
	d := reviewRoundAsCoordinator(f)
	forbidden := 0
	d.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		forbidden++
		return nil, fmt.Errorf("forbidden process observation")
	}
	d.StartupProcessTree = func(int) ([]StartupProcess, error) { forbidden++; return nil, fmt.Errorf("forbidden process tree") }
	d.Workspace.SourceWriters = func(string) (workspace.SourceWriterObservation, error) {
		forbidden++
		return workspace.SourceWriterObservation{}, fmt.Errorf("forbidden writer observation")
	}
	d.Fault = func(stage string) error {
		if !strings.HasPrefix(stage, "cli_install_") {
			forbidden++
			return fmt.Errorf("forbidden lifecycle stage %s", stage)
		}
		return nil
	}
	options := CLIInstallOptions{AllowRetainedSource: true}
	dest := filepath.Join(root, "installed-ply")
	call := func(dep Dependencies, destination, caller, action string) (CLIInstallation, error) {
		return InstallQualifiedCLIWithOptions(dep, root, o.RunID, destination, caller, action, options)
	}
	previewBefore := branchFixtureState(t, root)
	ready, err := call(d, dest, "synthetic coordinator", "check")
	if err != nil || !ready.SourceCleanupPending || !ready.AllowRetainedSource {
		t.Fatalf("preview: %+v %v", ready, err)
	}
	if !equal(previewBefore, branchFixtureState(t, root)) {
		t.Fatal("preview wrote native records")
	}
	d.Executable = func() (string, error) { return ready.Controller.Locator, nil }
	if _, err := InstallQualifiedCLI(d, root, o.RunID, dest, "synthetic coordinator", "install"); err == nil {
		t.Fatal("default route accepted closing")
	}
	if _, _, err := qualifiedCLIInstallation(d, root, o.RunID, dest, "synthetic coordinator"); err == nil {
		t.Fatal("activation qualifier accepted closing")
	}
	if _, err := ActivateParallelWorkspace(d, root, o.RunID, "synthetic coordinator", "", false); err == nil {
		t.Fatal("installation option leaked into activation")
	}

	assertBlocked := func(t *testing.T, dep Dependencies, destination, caller string) {
		t.Helper()
		before := branchFixtureState(t, root)
		destinationBefore, beforeErr := os.ReadFile(dest)
		if _, err := call(dep, destination, caller, "install"); err == nil {
			t.Fatal("invalid candidate installed")
		}
		if !equal(before, branchFixtureState(t, root)) {
			t.Fatal("rejected installation wrote native state")
		}
		destinationAfter, afterErr := os.ReadFile(dest)
		if string(destinationBefore) != string(destinationAfter) || os.IsNotExist(beforeErr) != os.IsNotExist(afterErr) {
			t.Fatal("rejected installation changed destination")
		}
	}
	t.Run("owner-cwd", func(t *testing.T) { dep := d; dep.CWD = f.D.CWD; assertBlocked(t, dep, dest, "synthetic coordinator") })
	t.Run("empty-attribution", func(t *testing.T) { assertBlocked(t, d, dest, "") })
	t.Run("relative-destination", func(t *testing.T) { assertBlocked(t, d, "relative-ply", "synthetic coordinator") })
	for name, change := range map[string]func(*workspace.WorkItemRegistry){
		"missing-result":      func(r *workspace.WorkItemRegistry) { r.TaskResults = nil },
		"missing-queue-close": func(r *workspace.WorkItemRegistry) { r.TaskQueueEvents = nil },
		"missing-base":        func(r *workspace.WorkItemRegistry) { r.EpicBaseVersions = nil },
		"stale-base": func(r *workspace.WorkItemRegistry) {
			r.EpicBaseVersions[len(r.EpicBaseVersions)-1].OID = strings.Repeat("1", 40)
		},
		"missing-ledger-result":     func(r *workspace.WorkItemRegistry) { r.IntegrationResults = nil },
		"forged-ledger-result":      func(r *workspace.WorkItemRegistry) { r.IntegrationResults[0].RecoveryStatus = "unknown" },
		"missing-ledger-authority":  func(r *workspace.WorkItemRegistry) { r.IntegrationAuthorities = nil },
		"missing-ledger-intent":     func(r *workspace.WorkItemRegistry) { r.IntegrationIntents = nil },
		"unresolved-target-attempt": func(r *workspace.WorkItemRegistry) { r.IntegrationAttempts[0].ResultID = nil },

		"qa-veto": func(r *workspace.WorkItemRegistry) {
			r.HumanQARecords = append(r.HumanQARecords, workspace.TaskHumanQARecord{TaskID: c.TaskResult.TaskID, TaskResultID: c.TaskResult.ID, ResultOID: c.OID, ResultTree: c.Tree, Outcome: "fail"})
		},
		"qa-blocked": func(r *workspace.WorkItemRegistry) {
			r.HumanQARecords = append(r.HumanQARecords, workspace.TaskHumanQARecord{TaskID: c.TaskResult.TaskID, TaskResultID: c.TaskResult.ID, ResultOID: c.OID, ResultTree: c.Tree, Outcome: "blocked"})
		},
	} {
		t.Run(name, func(t *testing.T) {
			dep := d
			dep.Workspace.WorkItems = retainedRegistryObservation{d.Workspace.WorkItems, change}
			assertBlocked(t, dep, dest, "synthetic coordinator")
		})
	}
	state, err := workflowRead(root, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	statePath := filepath.Join(o.Paths.RunRoot, "state.json")
	original, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	for name, change := range map[string]func(*workflowState){
		"unintegrated":      func(s *workflowState) { s.Result.Delivery.Candidates[0].Integration = nil },
		"unreleased":        func(s *workflowState) { s.Result.Delivery.OwnershipRelease = nil },
		"unresolved-effect": func(s *workflowState) { s.Result.Delivery.Attempt.State = "attempted" },
		"wrong-phase":       func(s *workflowState) { s.Result.Delivery.Phase = "working" },
		"wrong-candidate":   func(s *workflowState) { s.Result.Delivery.Candidates[0].OID = strings.Repeat("1", 40) },
	} {
		t.Run(name, func(t *testing.T) {
			var changed workflowState
			if err := json.Unmarshal(original, &changed); err != nil {
				t.Fatal(err)
			}
			change(&changed)
			raw, err := Canonical(changed)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(statePath, raw, 0600); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := os.WriteFile(statePath, original, 0600); err != nil {
					t.Fatal(err)
				}
			}()
			assertBlocked(t, d, dest, "synthetic coordinator")
		})
	}
	verificationRaw, err := workflowBound(c.Verification, 4<<20)
	if err != nil {
		t.Fatal(err)
	}
	var v deliveryVerificationReceipt
	if err := decode(verificationRaw, 4<<20, &v); err != nil {
		t.Fatal(err)
	}
	for name, binding := range map[string]FileBinding{"verification": c.Verification, "review": v.Review, "release": *state.Result.Delivery.OwnershipRelease, "integration": *c.Integration, "binary": ready.Controller} {
		for _, missing := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-missing-%t", name, missing), func(t *testing.T) {
				saved, err := os.ReadFile(binding.Locator)
				if err != nil {
					t.Fatal(err)
				}
				info, err := os.Stat(binding.Locator)
				if err != nil {
					t.Fatal(err)
				}
				if missing {
					err = os.Remove(binding.Locator)
				} else {
					err = os.WriteFile(binding.Locator, append(saved, '\n'), info.Mode().Perm())
				}
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					if err := os.WriteFile(binding.Locator, saved, info.Mode().Perm()); err != nil {
						t.Fatal(err)
					}
				}()
				assertBlocked(t, d, dest, "synthetic coordinator")
			})
		}
	}
	t.Run("manual-release", func(t *testing.T) {
		var release deliveryOwnershipRelease
		raw, err := workflowBound(*state.Result.Delivery.OwnershipRelease, 1<<20)
		if err != nil {
			t.Fatal(err)
		}
		if err := decode(raw, 1<<20, &release); err != nil {
			t.Fatal(err)
		}
		release.Origin = "bound_owner_release"
		path := writeAny(t, root, "synthetic-manual-release.json", release)
		changed := state
		var copyState workflowState
		raw, _ = Canonical(changed)
		if err := json.Unmarshal(raw, &copyState); err != nil {
			t.Fatal(err)
		}
		binding := FileBinding{path, hashFileTest(t, path)}
		old := *copyState.Result.Delivery.OwnershipRelease
		copyState.Result.Delivery.OwnershipRelease = &binding
		for i := range copyState.Result.Delivery.Events {
			if equal(copyState.Result.Delivery.Events[i].Binding, old) {
				copyState.Result.Delivery.Events[i].Binding = binding
			}
		}
		raw, _ = Canonical(copyState)
		if err := os.WriteFile(statePath, raw, 0600); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := os.WriteFile(statePath, original, 0600); err != nil {
				t.Fatal(err)
			}
		}()
		assertBlocked(t, d, dest, "synthetic coordinator")
	})
	t.Run("source-destination", func(t *testing.T) {
		assertBlocked(t, d, filepath.Join(c.TaskResult.SourceLocator, "ply"), "synthetic coordinator")
	})
	t.Run("native-control-destination", func(t *testing.T) { assertBlocked(t, d, filepath.Join(root, ".ply", "ply"), "synthetic coordinator") })

	for _, name := range []string{"failed-acceptance", "wrong-verification-run", "failed-review", "owner-review", "missing-review-session"} {
		t.Run(name, func(t *testing.T) {
			var forged deliveryVerificationReceipt
			if err := decode(verificationRaw, 4<<20, &forged); err != nil {
				t.Fatal(err)
			}
			switch name {
			case "failed-acceptance":
				forged.Outcome = "fail"
			case "wrong-verification-run":
				forged.RunID = "another-run"
			default:
				reviewRaw, err := workflowBound(forged.Review, 1<<20)
				if err != nil {
					t.Fatal(err)
				}
				var review map[string]any
				if err := json.Unmarshal(reviewRaw, &review); err != nil {
					t.Fatal(err)
				}
				switch name {
				case "failed-review":
					review["decision"] = "failed"
				case "owner-review":
					review["reviewer_session_id"] = deliveryOwnerSession(state)
				case "missing-review-session":
					review["reviewer_session_id"] = nil
				}
				path := writeAny(t, root, "synthetic-"+name+".json", review)
				forged.Review = FileBinding{path, hashFileTest(t, path)}
			}
			path := writeAny(t, root, "synthetic-verify-"+name+".json", forged)
			var changed workflowState
			if err := json.Unmarshal(original, &changed); err != nil {
				t.Fatal(err)
			}
			changed.Result.Delivery.Candidates[0].Verification = FileBinding{path, hashFileTest(t, path)}
			raw, err := Canonical(changed)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(statePath, raw, 0600); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := os.WriteFile(statePath, original, 0600); err != nil {
					t.Fatal(err)
				}
			}()
			assertBlocked(t, d, dest, "synthetic coordinator")
		})
	}
	t.Run("stale-target", func(t *testing.T) {
		runGit(t, f.Parent, "commit", "--allow-empty", "-m", "simulated later delivery")
		defer runGit(t, f.Parent, "reset", "--hard", c.OID)
		assertBlocked(t, d, dest, "synthetic coordinator")
	})
	t.Run("dirty-target", func(t *testing.T) {
		path := filepath.Join(f.Parent, "untracked-target")
		if err := os.WriteFile(path, []byte("drift"), 0600); err != nil {
			t.Fatal(err)
		}
		defer os.Remove(path)
		assertBlocked(t, d, dest, "synthetic coordinator")
	})

	// Seed two prior qualified installations so every history receipt and guarded
	// backup must survive the new native installation, its retry and rollback.
	old := []byte("old executable")
	if err := os.WriteFile(dest, old, 0755); err != nil {
		t.Fatal(err)
	}
	prior := ready.Source
	prior.RunID = "synthetic-earlier-1"
	prior.BinarySHA256 = clirelease.Digest([]byte("earlier one"))
	if _, err := (clirelease.Installer{}).InstallQualified(prior, []byte("earlier one"), dest); err != nil {
		t.Fatal(err)
	}
	prior.RunID = "synthetic-earlier-2"
	prior.BinarySHA256 = clirelease.Digest([]byte("earlier two"))
	priorReceipt, err := (clirelease.Installer{}).InstallQualified(prior, []byte("earlier two"), dest)
	if err != nil {
		t.Fatal(err)
	}
	priorManifest, err := os.ReadFile(dest + ".install.json")
	if err != nil {
		t.Fatal(err)
	}
	history, err := os.ReadDir(dest + ".install-history")
	if err != nil {
		t.Fatal(err)
	}
	preserved := map[string]string{}
	for _, entry := range history {
		path := filepath.Join(dest+".install-history", entry.Name())
		preserved[path] = hashFileTest(t, path)
	}
	before := branchFixtureState(t, root)
	sourceBefore, err := os.Stat(c.TaskResult.SourceLocator)
	if err != nil {
		t.Fatal(err)
	}
	installed, err := call(d, dest, "synthetic coordinator", "install")
	if err != nil || installed.State != "verified" || !installed.SourceCleanupPending || installed.ObservedSHA256 != ready.Source.BinarySHA256 {
		t.Fatalf("install: %+v %v", installed, err)
	}
	repeated, err := call(d, dest, "synthetic coordinator", "install")
	if err != nil || !equal(installed.Receipt, repeated.Receipt) {
		t.Fatalf("retry: %+v %v", repeated, err)
	}
	if _, err := call(d, dest+"-other", "synthetic coordinator", "install"); err == nil {
		t.Fatal("changed destination reused installation")
	}
	if _, err := call(d, dest, "another coordinator", "install"); err == nil {
		t.Fatal("changed attribution reused installation")
	}
	if _, err := call(d, dest, "synthetic coordinator", "observe"); err != nil {
		t.Fatal(err)
	}
	restored, err := call(d, dest, "synthetic coordinator", "rollback")
	if err != nil || restored.State != "rolled_back" || !restored.SourceCleanupPending {
		t.Fatalf("rollback: %+v %v", restored, err)
	}
	bytes, err := os.ReadFile(dest)
	if err != nil || string(bytes) != "earlier two" {
		t.Fatal("rollback lost preimage", err)
	}
	for path, digest := range preserved {
		if hashFileTest(t, path) != digest {
			t.Fatal("old history changed")
		}
	}
	archived, err := os.ReadFile(filepath.Join(dest+".install-history", strings.TrimPrefix(priorReceipt.ManifestSHA256, "sha256:")+".json"))
	if err != nil || string(archived) != string(priorManifest) {
		t.Fatal("superseded full manifest not preserved", err)
	}
	after := branchFixtureState(t, root)
	for path := range after {
		if strings.Contains(path, string(filepath.Separator)+"cli-installations"+string(filepath.Separator)) {
			delete(after, path)
		}
	}
	if !equal(before, after) {
		t.Fatal("installation changed Run/registry/cleanup ownership")
	}
	sourceAfter, err := os.Stat(c.TaskResult.SourceLocator)
	if err != nil || !os.SameFile(sourceBefore, sourceAfter) || gitOutput(t, f.Parent, "rev-parse", c.TaskResult.SourceRef) != c.OID {
		t.Fatal("installation removed source or branch", err)
	}
	if forbidden != 0 {
		t.Fatalf("forbidden process/cleanup calls: %d", forbidden)
	}
	// Independently finish only the disposable fixture's existing closeout
	// using simulated absence, then observe the SAME installation operation.
	coordinator := reviewRoundAsCoordinator(f)
	coordinator.LocalProcesses = f.D.LocalProcesses
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true, "inventory": []any{}})
	if _, handled, err := ResumeAutomaticLocalDelivery(coordinator, root, o.RunID, c.TaskResult.ID); err != nil || !handled {
		t.Fatalf("modeled completion: %v", err)
	}
	closedObservation, err := call(d, dest, "synthetic coordinator", "observe")
	if err != nil || closedObservation.State != "rolled_back" || closedObservation.SourceCleanupPending {
		t.Fatalf("closeout changed installation identity: %+v %v", closedObservation, err)
	}
	if forbidden != 0 {
		t.Fatalf("installation inspected processes after modeled closeout: %d", forbidden)
	}
	t.Run("active-stop", func(t *testing.T) {
		evidencePath := writeAny(t, root, "synthetic-stop-evidence.json", map[string]string{"source": "simulated refusal; no real action attempted"})
		evidence := FileBinding{evidencePath, hashFileTest(t, evidencePath)}
		if err := workflowUpdate(d, root, o.RunID, func(s *workflowState) error {
			_, err := deliveryAppendEvent(d, s, "synthetic-install-stop", "coordinator_stop", DeliveryStopRecord{Kind: "ply.workflow.coordinator-stop", SchemaVersion: 1, RunID: o.RunID, Input: DeliveryStopInput{EventID: "synthetic-install-stop", Stage: "delivery", Action: "synthetic forbidden effect", EffectStatus: DeliveryEffectNotStarted}, Source: evidence, Evidence: evidence})
			return err
		}); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := os.WriteFile(statePath, original, 0600); err != nil {
				t.Fatal(err)
			}
		}()
		assertBlocked(t, d, dest, "synthetic coordinator")
	})
	t.Log("TD125516B: exact retained bytes installed, smoke verified, full history/backup and closing source preserved; zero writer/process/lifecycle calls")
}
