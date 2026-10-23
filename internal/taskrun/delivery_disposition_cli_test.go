package taskrun

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// DEP02: modeled provider/Herdr, real native CLI/store/Git. Historical command
// invocation and child/wait receipts deliberately do not exist in this fixture.
func dep02Request(t *testing.T, f deliveryFixture, o WorkflowRun, event string) string {
	t.Helper()
	target := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	proof := writeAny(t, f.R.WorkspaceRoot, "dep02-assessment.json", map[string]string{
		"provenance":   "modeled owner report and static operation inspection; not runtime absence",
		"residual":     "Unknown fixture child/wait/descendants, temp/socket/cleanup, exact opened runtime inputs, exhaustive filesystem writes and pending writes; browser launch refused",
		"independence": "New native source snapshot binds current dirty inputs; fresh clean candidate, independent review, isolated acceptance input snapshots and exact target guards bound later actions. This is no claim of historical write/child absence",
	})
	return writeAny(t, f.R.WorkspaceRoot, event+".json", map[string]any{
		"kind": "ply.workflow.scoped-disposition", "schema_version": 1, "event_id": event,
		"run_id": o.RunID, "request_sha256": o.RequestSHA256, "mandate_sha256": o.Handoff.SHA256,
		"native_session": o.Transport.AgentSessionID,
		"source_oid":     gitOutput(t, target, "rev-parse", "HEAD"), "source_tree": gitOutput(t, target, "rev-parse", "HEAD^{tree}"),
		"retained":   []any{map[string]any{"stop": *o.Delivery.ActiveStop, "unknowns": []string{"child_wait_and_descendants", "temp_socket_and_cleanup", "runtime_input_bytes", "filesystem_writes", "pending_writes"}, "conflicts": []string{}, "provenance": "owner_and_static_assessment", "assessment": "Historical fixture effects remain unknown; missing invocation/child receipts are not reconstructed", "evidence": []FileBinding{{proof, hashFileTest(t, proof)}}}},
		"operations": []any{map[string]any{"operation": "prepare", "boundary": "fresh_source_snapshot", "provenance": "owner_and_static_assessment", "action": "Prepare the authorized file-only candidate without using the historical fixture", "reads": []string{"path:" + target}, "writes": []string{"path:" + target, "path:" + f.AcceptancePath, "native:git-index", "native:ignored-build-output", "native:verification-scratch"}, "children": []string{"process:node-check", "process:node-build", "process:node-precache", "process:git-commit"}, "assessment": "This decision accepts retained unlocated historical writes and child uncertainty for a newly observed source snapshot. Known concrete input/write/child overlaps block the operation; subsequent source/candidate/input drift invalidates it", "inputs": []FileBinding{{proof, hashFileTest(t, proof)}}}},
	})
}

func dep02Stop(t *testing.T, binary string, f deliveryFixture, o WorkflowRun) WorkflowRun {
	t.Helper()
	proof := writeAny(t, f.R.WorkspaceRoot, "dep02-denial.json", map[string]string{"source": "modeled refusal after fixture start; invocation and child wait unknown"})
	return automaticCLIOK(t, binary, f, "execute", "observe-stop", o.RunID, "--event-id", "dep02-denied", "--provider", "codex", "--native-session", o.Transport.AgentSessionID, "--stage", "work", "--action", "Launch historical browser fixture", "--effect-status", DeliveryEffectPartial, "--observer", "modeled coordinator", "--reason", "Modeled browser refused after fixture launch; old child outcome unknown", "--evidence", proof)
}

func TestDEP02PublicDispositionRetainsPartialStop(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, true, 30)
	o = dep02Stop(t, binary, f, o)
	stop := *o.Delivery.ActiveStop
	original, err := os.ReadFile(stop.Locator)
	if err != nil {
		t.Fatal(err)
	}
	request := dep02Request(t, f, o, "dep02-preparation")
	args := []string{"execute", "disposition", o.RunID, "--context", o.Paths.Context, "--file", request}
	o = automaticCLIOK(t, binary, f, args...)
	o = automaticCLIOK(t, binary, f, args...)
	if o.Delivery.ActiveStop == nil || *o.Delivery.ActiveStop != stop {
		t.Fatal("disposition cleared original stop")
	}
	after, _ := os.ReadFile(stop.Locator)
	if string(original) != string(after) {
		t.Fatal("stop bytes changed")
	}
	raw, _ := json.Marshal(o.DeliveryStatus)
	var status map[string]json.RawMessage
	if err := json.Unmarshal(raw, &status); err != nil {
		t.Fatal(err)
	}
	if string(status["scoped_continuation"]) == "null" || len(status["scoped_continuation"]) == 0 {
		t.Fatalf("missing disposition/readback: %s", raw)
	}
	automaticCLIOK(t, binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--check-operation", "prepare")
	// Files unrelated to the preserved stop can change only under an exact source
	// assessment. A later operation must register its own newly observed source.
	target := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	if err := os.WriteFile(filepath.Join(target, "README.md"), []byte("changed source\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := automaticCLI(binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--check-operation", "prepare"); err == nil {
		t.Fatal("source drift ignored")
	}
	if _, _, err := automaticCLI(binary, f, "execute", "resolve-stop", o.RunID, "--context", o.Paths.Context, "--stop", stop.SHA256, "--resolution", "Scoped work is independent", "--action", "Prepare candidate"); err == nil {
		t.Fatal("ordinary resolution cleared partial effect")
	}
}

func dep02CandidateRequest(t *testing.T, f deliveryFixture, o WorkflowRun, event string, binary string) string {
	t.Helper()
	path := dep02Request(t, f, o, event)
	var in map[string]any
	raw, e := os.ReadFile(path)
	if e != nil {
		t.Fatal(e)
	}
	if e = json.Unmarshal(raw, &in); e != nil {
		t.Fatal(e)
	}
	first := in["operations"].([]any)[0].(map[string]any)
	ops := []any{}
	for _, operation := range []string{"report", "review", "verify", "integrate"} {
		op := map[string]any{}
		for k, v := range first {
			op[k] = v
		}
		op["operation"] = operation
		op["boundary"] = map[string]string{"report": "fresh_source_snapshot", "review": "exact_clean_candidate", "verify": "candidate_and_isolated_verifier", "integrate": "qualified_candidate_and_exact_target"}[operation]
		op["action"] = "Perform independent native " + operation + " under the original local agreement"
		op["writes"] = []string{}
		op["children"] = []string{}
		if operation == "verify" {
			op["writes"] = []string{"native:verification-scratch", "native:run-evidence"}
			op["children"] = []string{"process:node-check", "process:node-build", "process:node-precache"}
			op["inputs"] = append([]any{map[string]string{"locator": f.AcceptancePath, "sha256": hashFileTest(t, f.AcceptancePath)}, map[string]string{"locator": binary, "sha256": hashFileTest(t, binary)}}, first["inputs"].([]any)...)
		}
		if operation == "integrate" {
			op["writes"] = []string{"native:local-target", "native:run-evidence"}
		}
		ops = append(ops, op)
	}
	if log := filepath.Join(f.R.WorkspaceRoot, "node-preparation.log"); dep02FileExists(log) {
		for _, value := range ops {
			op := value.(map[string]any)
			op["inputs"] = append(op["inputs"].([]any), map[string]string{"locator": log, "sha256": hashFileTest(t, log)})
		}
	}
	in["operations"] = ops
	return writeAny(t, f.R.WorkspaceRoot, event+"-candidate.json", in)
}

func TestDEP02PublicCandidateReviewVerifyLocalBranch(t *testing.T) {
	binary := automaticCLIBinary(t)
	original := continuitySourceBinary(t, "d58dbd4def24bee385307db7dc18fe1583958533", "PLY_DEP02_OLD_BINARY")
	original, e := filepath.EvalSymlinks(original)
	if e != nil {
		t.Fatal(e)
	}
	for name, control := range map[string]string{"current": binary, "frozen_d58dbd4": original} {
		t.Run(name, func(t *testing.T) { dep02CLIJourney(t, binary, control) })
	}
}
func dep02CLIJourney(t *testing.T, binary, control string) {
	agreement := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalBranch, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "modeled coordinator", TimeoutSeconds: 30}}
	f := receiptFixture(t, "codex", deliveryFixtureOptions{Agreement: &agreement, ParentRef: "main"})
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.R.Runtime.PlyExecutable = Executable{control, hashFileTest(t, control)}
	f.D.Executable = func() (string, error) { return control, nil }
	reviewRoundPatchHerdr(t, &f)
	dep02ProcessModel(t, f)
	workflowTestModel(t, f.workflowFixture, map[string]any{"caller_shell": os.Getppid()})
	o := workflowTestStart(t, f.workflowFixture)
	o = automaticCLIOK(t, control, f, "execute", "check-in", o.RunID, "--context", o.Paths.Context, "--provider", "codex", "--native-session", o.Transport.AgentSessionID)
	if o.Delivery.PermissionState != RuntimePerAction || o.DeliveryStatus.Startup.StartupProtocol != MandateReceiptV1 {
		t.Fatal("PWA fixture changed the actual receipt protocol")
	}
	stateAtReceipt, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	var receipt DeliveryAcceptance
	if e = json.Unmarshal(dep02Read(t, stateAtReceipt.Acceptance.Locator), &receipt); e != nil {
		t.Fatal(e)
	}
	if receipt.DeliveryPermission.PermissionConfirmed || receipt.DeliveryPermission.StartupProtocol != MandateReceiptV1 || receipt.DeliveryPermission.MandateSHA256 != o.Handoff.SHA256 {
		t.Fatal("receipt asserted permission or lost frozen mandate")
	}
	dep02Evidence(t, "mandate-receipt", receipt)
	target := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	if reg, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot); e != nil || reg.FormatVersion != 4 {
		t.Fatalf("expected unchanged format4 fixture: %v", e)
	}
	if e := os.WriteFile(filepath.Join(target, "README.md"), []byte("dirty original work\n"), 0600); e != nil {
		t.Fatal(e)
	}
	dep02ModelFiles(t, target, false)
	o = dep02Stop(t, control, f, o)
	stop := *o.Delivery.ActiveStop
	requestSHA, contextSHA, controlSHA := hashFileTest(t, f.File), hashFileTest(t, o.Paths.Context), hashFileTest(t, control)
	if control != binary {
		if _, output, e := automaticCLI(control, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--candidate-binary", binary); e == nil || !strings.Contains(string(output), "preserved provider stop") {
			t.Fatalf("historical guard unexpectedly changed: %v %s", e, output)
		}
		if e := os.Remove(f.R.Runtime.Executable.Path); e != nil {
			t.Fatal(e)
		}
	}
	request := dep02Request(t, f, o, "dirty-preparation")
	automaticCLIOK(t, binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--file", request)
	automaticCLIOK(t, binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--check-operation", "prepare")
	// Intentional acceptance/validator/source writes invalidate preparation-1.
	// Refresh before executing the newly prepared Node/build/test operation.
	dep02ModelFiles(t, target, true)
	if e := os.WriteFile(filepath.Join(target, "README.md"), []byte("independent expected\n"), 0600); e != nil {
		t.Fatal(e)
	}
	script := "#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = 'independent expected'\nnpm --offline --ignore-scripts run check\nnpm --offline --ignore-scripts run build\nnpm --offline --ignore-scripts run test:build\nprintf 'actual Node build/test children waited; historical effects remain unknown\\n'\n"
	if e := os.WriteFile(f.AcceptancePath, []byte(script), 0600); e != nil {
		t.Fatal(e)
	}
	refreshed := dep02Request(t, f, o, "prepared-node-inputs")
	automaticCLIOK(t, binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--file", refreshed)
	automaticCLIOK(t, binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--check-operation", "prepare")
	scratch := t.TempDir()
	cmd := exec.Command("/bin/sh", f.AcceptancePath)
	cmd.Dir = target
	cmd.Env = append(os.Environ(), "HOME="+scratch, "TMPDIR="+scratch)
	output, e := cmd.CombinedOutput()
	if e != nil {
		t.Fatalf("actual modeled Node preparation: %v %s", e, output)
	}
	if e = os.WriteFile(filepath.Join(f.R.WorkspaceRoot, "node-preparation.log"), output, 0600); e != nil {
		t.Fatal(e)
	}
	runGit(t, target, "add", "README.md", ".gitignore", "package.json", "scripts/model-artifact.mjs")
	runGit(t, target, "commit", "-m", "prepare modeled Node build and artifact checks")
	request = dep02CandidateRequest(t, f, o, "exact-candidate", binary)
	o = automaticCLIOK(t, binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--file", request)
	callback := control
	if control != binary {
		observation := deliveryTestRuntimeObservation(t, f, o)
		if raw, e := continuityCLI(binary, target, "workflow", "execute", "continue", o.RunID, "--context", o.Paths.Context); e == nil || !strings.Contains(string(raw), "runtime_observation_required") {
			t.Fatalf("missing launcher gained authority without runtime evidence: %v %s", e, raw)
		}
		continuityCLIOK(t, binary, target, "workflow", "execute", "continue", o.RunID, "--context", o.Paths.Context, "--runtime-evidence", observation, "--check")
		raw := continuityCLIOK(t, binary, target, "workflow", "execute", "continue", o.RunID, "--context", o.Paths.Context, "--runtime-evidence", observation)
		dep02Evidence(t, "runtime-observation", json.RawMessage(dep02Read(t, observation)))
		dep02Evidence(t, "continuation", json.RawMessage(raw))
		var result struct {
			Preview struct {
				Control Executable `json:"control_executable"`
			} `json:"preview"`
		}
		if e := json.Unmarshal(raw, &result); e != nil || result.Preview.Control.Path == "" {
			t.Fatalf("compatible continuation: %v %s", e, raw)
		}
		callback = result.Preview.Control.Path
		// The original control still reads the same strict frozen state;
		// only the supported continuation becomes the current callback.
		automaticCLIOK(t, control, f, "execute", "show", o.RunID)
	}
	if _, output, e := automaticCLI(callback, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--candidate-binary", binary); e == nil || !strings.Contains(string(output), "review") {
		t.Fatal("missing separate review qualified")
	}
	review, detail, exit := al13CLI(t, binary, f, "w-fixture:root", "dep02-review", "review", o.RunID, "--coordinator", "modeled coordinator")
	if exit != 0 {
		t.Fatalf("review: %s", detail)
	}
	round := reviewRoundLatest(t, review)
	returned := reviewRoundReturnFile(t, f, "dep02-review.json", `{"decision":"passed","summary":"Modeled independent review inspects the exact file-only candidate; retained effects remain unknown","findings":[],"evidence_gaps":["Historical runtime inputs/filesystem/pending writes and fixture child remain unknown; no live provider approval."]}`)
	review, detail, exit = al13CLI(t, binary, f, round.ReviewerPane, "dep02-return", "review-return", o.RunID, "--context", round.Context.Locator, "--file", returned)
	if exit != 0 {
		t.Fatalf("review return: %s", detail)
	}
	unbound := filepath.Join(f.R.WorkspaceRoot, "unbound-artifact")
	if e := os.WriteFile(unbound, []byte("#!/bin/sh\nexit 0\n"), 0700); e != nil {
		t.Fatal(e)
	}
	if _, _, e := automaticCLI(callback, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--candidate-binary", unbound); e == nil {
		t.Fatal("scoped verifier accepted an undeclared candidate executable")
	}
	originalAcceptance := dep02Read(t, f.AcceptancePath)
	if e := os.WriteFile(f.AcceptancePath, append(append([]byte{}, originalAcceptance...), []byte("# changed input\n")...), 0600); e != nil {
		t.Fatal(e)
	}
	if _, _, e := automaticCLI(callback, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--candidate-binary", binary); e == nil {
		t.Fatal("acceptance input drift ignored")
	}
	if e := os.WriteFile(f.AcceptancePath, originalAcceptance, 0600); e != nil {
		t.Fatal(e)
	}
	o = automaticCLIOK(t, callback, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--candidate-binary", binary)
	v := o.DeliveryStatus.Verification
	if v == nil || !v.Current || !v.InputsMatch || v.Qualification != "qualified" || v.Exit == nil || *v.Exit != 0 {
		t.Fatalf("not exact qualified: %+v", v)
	}
	c := o.Delivery.Candidates[len(o.Delivery.Candidates)-1]
	if _, _, e := automaticCLI(binary, f, "execute", "release", o.RunID, "--context", o.Paths.Context, "--reason", "Modeled unauthorized writer transfer under retained effects"); e == nil {
		t.Fatal("scoped integration authority leaked into writer release")
	}
	state, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if _, e := deliveryAuthorityState(f.D, state); e == nil {
		t.Fatal("generic authority bypassed original stop")
	}
	// Native reservation with a deliberately lost pre-effect reply, not a
	// handwritten pending-attempt or store edit. A different action is blocked;
	// the same native integration ledger alone may resume its reserved effect.
	parentBefore := gitOutput(t, f.Parent, "rev-parse", "HEAD")
	reserved := f.D
	reserved.Fault = func(stage string) error {
		if stage == "delivery_after_integration_reservation" {
			return errors.New("modeled loss after native reservation")
		}
		return nil
	}
	if _, e := CompleteLocalDelivery(reserved, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, ""); e == nil || !strings.Contains(e.Error(), "modeled loss") {
		t.Fatalf("did not preserve native reservation: %v", e)
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != parentBefore {
		t.Fatal("reservation unexpectedly integrated")
	}
	if _, _, e := automaticCLI(binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--check-operation", "report"); e == nil {
		t.Fatal("pending attempt granted another operation")
	}
	if _, _, e := automaticCLI(callback, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--candidate-binary", binary); e == nil {
		t.Fatal("pending integration allowed verifier execution")
	}
	pendingRequest := dep02CandidateRequest(t, f, o, "pending-new-registration", binary)
	if _, _, e := automaticCLI(binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--file", pendingRequest); e == nil {
		t.Fatal("pending attempt allowed another registration")
	}
	integrated, output, e := automaticCLI(callback, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context)
	if e == nil || !strings.Contains(string(output), "retained_effects_cleanup_pending") {
		t.Fatalf("expected honest retained cleanup: %v %s", e, output)
	}
	if integrated.DeliveryStatus.FinalDelivery.State != "delivered" || gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
		t.Fatalf("native local branch delivery not observed: %s", output)
	}
	reflog := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/main")
	if _, out, e := automaticCLI(callback, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context); e == nil || !strings.Contains(string(out), "retained_effects_cleanup_pending") {
		t.Fatalf("same delivery readback lost residual blocker: %v %s", e, out)
	}
	if gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/main") != reflog {
		t.Fatal("lost reply replayed integration")
	}
	if integrated.Delivery.ActiveStop == nil || *integrated.Delivery.ActiveStop != stop {
		t.Fatal("delivery hid original stop")
	}
	if _, e := os.Stat(target); e != nil {
		t.Fatal("retained source removed", e)
	}
	if hashFileTest(t, f.File) != requestSHA || hashFileTest(t, o.Paths.Context) != contextSHA || hashFileTest(t, control) != controlSHA {
		t.Fatal("rewrote historical request/context/control")
	}
	if _, _, e := automaticCLI(binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--check-operation", "report"); e == nil {
		t.Fatal("delivered source reopened for work")
	}
	queueRaw := continuityCLIOK(t, binary, target, "workspace", "task", "queue", "list", "--project", "ply", "--repo", "ply", "--epic", "epic", "--schema-version", "2")
	var queue workspace.WorkspaceTaskQueueReadback
	if e := json.Unmarshal(queueRaw, &queue); e != nil || queue.Current != nil || len(queue.Active) > 0 {
		t.Fatalf("queue closure absent: %v %s", e, queueRaw)
	}
	baseRaw := continuityCLIOK(t, binary, target, "workspace", "epic", "base", "show", "epic", "--repo", "ply")
	var base workspace.EpicBaseResult
	if e := json.Unmarshal(baseRaw, &base); e != nil || deliveryCLILatestBase(base).OID != c.OID {
		t.Fatalf("base update absent: %v %s", e, baseRaw)
	}
	if reg, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot); e != nil || reg.FormatVersion != 4 {
		t.Fatalf("format4 unexpectedly changed: %v", e)
	}
	if integrated.Delivery.OwnershipRelease != nil {
		t.Fatal("automatic scoped delivery transferred original writer")
	}
	dep02Evidence(t, "queue", json.RawMessage(queueRaw))
	dep02Evidence(t, "base", json.RawMessage(baseRaw))
	dep02Evidence(t, "request", json.RawMessage(dep02Read(t, request)))
	dep02Evidence(t, "acceptance", string(dep02Read(t, f.AcceptancePath)))
	dep02Evidence(t, "node-preparation", string(dep02Read(t, filepath.Join(f.R.WorkspaceRoot, "node-preparation.log"))))
	dep02Evidence(t, "review-return", review)
	dep02Evidence(t, "verification-receipt", json.RawMessage(dep02Read(t, c.Verification.Locator)))
	dep02Evidence(t, "qualified", o)
	dep02Evidence(t, "delivered", integrated)
	dep02Evidence(t, "identity", map[string]any{"original_control": control, "original_control_sha256": controlSHA, "original_context_sha256": contextSHA, "original_request_sha256": requestSHA, "candidate_binary_sha256": hashFileTest(t, binary), "model": "Herdr/provider/judgment are modeled; actual native Git, review-return, child/wait, verification and integration executed in disposable workspace"})
	if workflowTestCalls(t, f.workflowFixture, "pane close") != 0 {
		t.Fatal("delivery closed original owner")
	}
}

func dep02ProcessModel(t *testing.T, f deliveryFixture) {
	t.Helper()
	script := `#!/usr/bin/python3
import os
# The fake lsof process is a child of this CLI; only owned PIDs are read.
print("p%d\nR%d\nn/modeled-cwd" % (os.getppid(), ` + fmt.Sprint(os.Getppid()) + `))
`
	if e := os.WriteFile(filepath.Join(f.R.WorkspaceRoot, "bin", "lsof"), []byte(script), 0700); e != nil {
		t.Fatal(e)
	}
}

func TestDEP02DispositionUnsafeScopesAndIdentity(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, true, 30)
	o = dep02Stop(t, binary, f, o)
	original := dep02Request(t, f, o, "safety")
	raw, e := os.ReadFile(original)
	if e != nil {
		t.Fatal(e)
	}
	cases := map[string]func(map[string]any){
		"wrong-run":            func(v map[string]any) { v["run_id"] = "foreign-run" },
		"wrong-request":        func(v map[string]any) { v["request_sha256"] = "sha256:" + strings.Repeat("0", 64) },
		"wrong-mandate":        func(v map[string]any) { v["mandate_sha256"] = "sha256:" + strings.Repeat("0", 64) },
		"wrong-session":        func(v map[string]any) { v["native_session"] = "another-owner" },
		"source-drift":         func(v map[string]any) { v["source_oid"] = strings.Repeat("0", 40) },
		"permission-expansion": func(v map[string]any) { v["operations"].([]any)[0].(map[string]any)["operation"] = "permissions" },
		"cleanup-expansion":    func(v map[string]any) { v["operations"].([]any)[0].(map[string]any)["operation"] = "cleanup" },
		"writer-transfer":      func(v map[string]any) { v["operations"].([]any)[0].(map[string]any)["operation"] = "release" },
		"retire-expansion":     func(v map[string]any) { v["operations"].([]any)[0].(map[string]any)["operation"] = "retire" },
		"unknowns-omitted": func(v map[string]any) {
			v["retained"].([]any)[0].(map[string]any)["unknowns"] = []string{"child_wait_and_descendants"}
		},
		"false-runtime-provenance": func(v map[string]any) { v["retained"].([]any)[0].(map[string]any)["provenance"] = "runtime_absent" },
		"refused-replay": func(v map[string]any) {
			v["operations"].([]any)[0].(map[string]any)["action"] = "Launch historical browser fixture"
		},
		"broad-write": func(v map[string]any) {
			v["operations"].([]any)[0].(map[string]any)["writes"] = []string{"path:/Users"}
		},
		"arbitrary-native-effect": func(v map[string]any) {
			v["operations"].([]any)[0].(map[string]any)["writes"] = []string{"native:permission-policy"}
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			var value map[string]any
			if e := json.Unmarshal(raw, &value); e != nil {
				t.Fatal(e)
			}
			mutate(value)
			value["event_id"] = "safety-" + name
			path := writeAny(t, f.R.WorkspaceRoot, "negative-"+name+".json", value)
			if _, _, e := automaticCLI(binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--file", path); e == nil {
				t.Fatal("unsafe disposition accepted")
			}
		})
	}
	if _, _, e := automaticCLI(binary, f, "execute", "disposition", o.RunID, "--context", f.AcceptancePath, "--file", original); e == nil {
		t.Fatal("wrong context accepted")
	}
	// The live provider identity must match, independently of input claims.
	workflowTestModel(t, f.workflowFixture, map[string]any{"agent_session_id": "changed-session"})
	if _, _, e := automaticCLI(binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--file", original); e == nil {
		t.Fatal("fresh session drift accepted")
	}
}

func TestDEP02DispositionConcurrentLostRepliesAndConflicts(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, true, 30)
	o = dep02Stop(t, binary, f, o)
	path := dep02Request(t, f, o, "concurrent")
	args := []string{"execute", "disposition", o.RunID, "--context", o.Paths.Context, "--file", path}
	type result struct {
		run WorkflowRun
		raw []byte
		err error
	}
	replies := make(chan result, 2)
	for i := 0; i < 2; i++ {
		go func() { run, raw, err := automaticCLI(binary, f, args...); replies <- result{run, raw, err} }()
	}
	a, b := <-replies, <-replies
	if a.err != nil || b.err != nil {
		t.Fatalf("identical concurrent registration: %v %v %s %s", a.err, b.err, a.raw, b.raw)
	}
	if !equal(a.run.DeliveryStatus.ScopedContinuation.Record, b.run.DeliveryStatus.ScopedContinuation.Record) {
		t.Fatal("concurrent calls created different observations")
	}
	again := automaticCLIOK(t, binary, f, args...)
	if !equal(again.DeliveryStatus.ScopedContinuation.Record, a.run.DeliveryStatus.ScopedContinuation.Record) {
		t.Fatal("lost response retry changed record")
	}
	// Even formatting-only changed bytes under the same ID are a conflict.
	whitespace := writeAny(t, f.R.WorkspaceRoot, "whitespace-placeholder.json", map[string]string{"unused": "placeholder"})
	requestBytes := dep02Read(t, path)
	if e := os.WriteFile(whitespace, append(requestBytes, '\n'), 0600); e != nil {
		t.Fatal(e)
	}
	if _, _, e := automaticCLI(binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--file", whitespace); e == nil {
		t.Fatal("changed request bytes reused event ID")
	}
	var changed map[string]any
	raw, _ := os.ReadFile(path)
	if e := json.Unmarshal(raw, &changed); e != nil {
		t.Fatal(e)
	}
	changed["operations"].([]any)[0].(map[string]any)["assessment"] = "Changed assessment under conflicting existing ID"
	different := writeAny(t, f.R.WorkspaceRoot, "conflicting-bytes.json", changed)
	if _, _, e := automaticCLI(binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--file", different); e == nil {
		t.Fatal("conflicting request replaced original")
	}
}

func TestDEP02ConcreteConflictsAndAdditionalStops(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, true, 30)
	deliveryTestCandidate(t, f, "independent expected")
	o = dep02Stop(t, binary, f, o)
	path := dep02CandidateRequest(t, f, o, "concrete-conflicts", binary)
	var in map[string]any
	raw, _ := os.ReadFile(path)
	if e := json.Unmarshal(raw, &in); e != nil {
		t.Fatal(e)
	}
	in["retained"].([]any)[0].(map[string]any)["conflicts"] = []string{"native:verification-scratch"}
	// Omit the actual scratch from the declaration: native dependencies must
	// still block verify, independently of chosen resource names.
	for _, op := range in["operations"].([]any) {
		m := op.(map[string]any)
		if m["operation"] == "verify" {
			m["writes"] = []string{}
			m["children"] = []string{}
		}
	}
	path = writeAny(t, f.R.WorkspaceRoot, "conflicts.json", in)
	o = automaticCLIOK(t, binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--file", path)
	if _, out, e := automaticCLI(binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--check-operation", "verify"); e == nil || !strings.Contains(string(out), "retained possible effect native:verification-scratch") {
		t.Fatalf("implicit verifier dependency ignored: %v %s", e, out)
	}
	automaticCLIOK(t, binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--check-operation", "report")
	// A new native stop invalidates the former operation decision. Registering
	// only its latest stop cannot hide the older unresolved partial effect.
	proof := writeAny(t, f.R.WorkspaceRoot, "extra-denial.json", map[string]string{"source": "second modeled refusal; source write completion unknown"})
	o = automaticCLIOK(t, binary, f, "execute", "observe-stop", o.RunID, "--event-id", "additional-stop", "--provider", "codex", "--native-session", o.Transport.AgentSessionID, "--stage", "work", "--action", "Different historical operation", "--effect-status", DeliveryEffectUnknown, "--observer", "modeled coordinator", "--reason", "Additional modeled effect remains unknown", "--evidence", proof)
	if _, _, e := automaticCLI(binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--check-operation", "report"); e == nil {
		t.Fatal("additional stop ignored")
	}
	latestOnly := dep02Request(t, f, o, "latest-only")
	if _, _, e := automaticCLI(binary, f, "execute", "disposition", o.RunID, "--context", o.Paths.Context, "--file", latestOnly); e == nil {
		t.Fatal("newest stop hid older partial stop")
	}
	shown := automaticCLIOK(t, binary, f, "execute", "show", o.RunID)
	if shown.Delivery.ActiveStop == nil || len(shown.Delivery.Events) != 2 {
		t.Fatal("independent read lost stops")
	}
}

func dep02Evidence(t *testing.T, name string, value any) {
	t.Helper()
	if dir := os.Getenv("PLY_DEP02_EVIDENCE_DIR"); dir != "" {
		raw, e := json.MarshalIndent(value, "", "  ")
		if e != nil {
			t.Fatal(e)
		}
		if e = os.MkdirAll(dir, 0700); e != nil {
			t.Fatal(e)
		}
		file := filepath.Join(dir, strings.ReplaceAll(t.Name(), "/", "_")+"-"+name+".json")
		if e = os.WriteFile(file, raw, 0600); e != nil {
			t.Fatal(e)
		}
	}
}

func dep02Read(t *testing.T, path string) []byte {
	t.Helper()
	raw, e := os.ReadFile(path)
	if e != nil {
		t.Fatal(e)
	}
	return raw
}
func dep02FileExists(path string) bool { _, e := os.Stat(path); return e == nil }

func dep02ModelFiles(t *testing.T, target string, complete bool) {
	t.Helper()
	files := map[string]string{
		".gitignore":                 "dist/\ndist-server/\n.cache/\n",
		"package.json":               `{"private":true,"scripts":{"check":"node scripts/model-artifact.mjs check","build":"node scripts/model-artifact.mjs build","test:build":"node scripts/model-artifact.mjs precache"}}`,
		"scripts/model-artifact.mjs": "throw new Error('incomplete inherited validator');\n",
	}
	if complete {
		files["scripts/model-artifact.mjs"] = `import assert from 'node:assert/strict';
import {mkdirSync,readFileSync,writeFileSync} from 'node:fs';
import {join} from 'node:path';
import {spawnSync} from 'node:child_process';
const expected='independent expected';
assert.equal(readFileSync('README.md','utf8').trim(),expected);
const task=process.argv[2];
mkdirSync('.cache',{recursive:true});
writeFileSync(join(process.env.TMPDIR,'node-'+task),task);
if(task==='check') {
 const child=spawnSync(process.execPath,['-e',"require('node:fs').writeFileSync('.cache/child-wait.json',JSON.stringify({exit:0,waited:true}))"],{stdio:'pipe'});
 assert.equal(child.status,0);assert.equal(child.signal,null);
 assert.equal(JSON.parse(readFileSync('.cache/child-wait.json','utf8')).waited,true);
} else if(task==='build') {
 for(const dir of ['dist','dist-server']) {mkdirSync(dir,{recursive:true});writeFileSync(join(dir,'artifact.json'),JSON.stringify({candidate:expected,modeled:true}));}
} else if(task==='precache') {
 for(const dir of ['dist','dist-server']) assert.equal(JSON.parse(readFileSync(join(dir,'artifact.json'),'utf8')).candidate,expected);
 writeFileSync('.cache/precache.json',JSON.stringify({validated:['dist','dist-server']}));
} else throw new Error('unexpected modeled operation');
console.log('modeled '+task+' actual child/write checks passed; historical writes and descendants remain unknown');
`
	}
	for rel, contents := range files {
		path := filepath.Join(target, rel)
		if e := os.MkdirAll(filepath.Dir(path), 0700); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(path, []byte(contents), 0600); e != nil {
			t.Fatal(e)
		}
	}
}

// The existing observation schema has no sandbox field. A mandate receipt
// must validate against its preserved matches_contract=null native acceptance,
// not fabricate a new sandbox attestation during compatible continuation.
func TestDEP02ReceiptRuntimeObservationPreservesUnattestedSandbox(t *testing.T) {
	f := receiptFixture(t, "codex")
	o := workflowTestStart(t, f.workflowFixture)
	o, e := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}, true)
	if e != nil {
		t.Fatal(e)
	}
	path := deliveryTestRuntimeObservation(t, f, o)
	state, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if e = deliveryRuntimeObservation(state, FileBinding{path, hashFileTest(t, path)}); e != nil {
		t.Fatalf("same receipt observation incorrectly requires sandbox attestation: %v", e)
	}
}
