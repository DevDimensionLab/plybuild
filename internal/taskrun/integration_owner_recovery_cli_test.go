package taskrun

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const integrationOwnerRecoveryStandin = `
if model.get("closed_transport") and cmd == ["agent", "get"]:
    print(json.dumps({"id":"fixture", "error":{"code":"agent_not_found", "message":"closed fixture agent"}}))
    raise SystemExit(1)
if model.get("closed_transport") and cmd == ["agent", "list"]:
    print(json.dumps({"id":"fixture", "result":{"agents":model.get("recovery_agents", [])}}))
    raise SystemExit(0)
`

func integrationOwnerRecoveryFixture(t *testing.T, binary string) (deliveryFixture, WorkflowRun, string) {
	t.Helper()
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{SchemaVersion: 2, IntegrationOwner: workspace.IntegrationOwnerHuman, Mode: workspace.DeliveryLocalBranch}, ParentRef: "main"})
	herdr, err := os.ReadFile(f.R.Herdr.Executable.Path)
	if err != nil {
		t.Fatal(err)
	}
	herdr = bytes.Replace(herdr, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+closedTransportInspectStandin+integrationOwnerRecoveryStandin), 1)
	if err = os.WriteFile(f.R.Herdr.Executable.Path, herdr, 0700); err != nil {
		t.Fatal(err)
	}
	f.R.Herdr.Executable.SHA256 = hash(herdr)
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	o := deliveryCLIQualified(t, f)
	source := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	id := deliveryCLIString(t, deliveryCLIOK(t, binary, source, "register", "--file", deliveryCLIRegisterInput(t, f, o, "closed-owner")), "id")
	o = reportCandidateQuestion(t, f, o, "qa-before-closed-pane")
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true})
	return f, o, id
}

func ownerRecoverySnapshot(t *testing.T, f deliveryFixture) string {
	t.Helper()
	var kept []string
	for _, entry := range strings.Split(deliveryCLISnapshot(t, f.R.WorkspaceRoot), "\n") {
		// The synthetic Herdr records read calls; it is not native Ply state.
		if !strings.HasPrefix(entry, f.Calls+":") {
			kept = append(kept, entry)
		}
	}
	return strings.Join(kept, "\n")
}

func TestIntegrationCLIRecoversClosedOwnerWithExplicitHumanRelease(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f, o, id := integrationOwnerRecoveryFixture(t, binary)
	before := ownerRecoverySnapshot(t, f)
	args := []string{"release", "--delivery", id, "--recover", "--actor", "Synthetic fixture human", "--observation", "All fixture writers are stopped."}
	p := integrationCLIOK(t, binary, f.Parent, "", append(args, "--check")...)
	if deliveryCLIString(t, p, "state") != "ready" {
		t.Fatalf("closed qualified owner cannot be recovered: %s", p)
	}
	for _, answer := range []string{"", "pass\n", "released"} {
		if _, _, err := integrationCLIRun(binary, f.Parent, answer, args...); err == nil {
			t.Fatalf("release accepted missing or wrong explicit answer %q", answer)
		}
	}
	if before != ownerRecoverySnapshot(t, f) {
		t.Fatal("preview or rejected input mutated native workspace")
	}
	integrationCLIOK(t, binary, f.Parent, "released\n", args...)
	state, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || state.Result.Delivery.OwnershipRelease == nil {
		t.Fatalf("native release missing: %v", err)
	}
	if state.Result.Delivery.Candidates[0].HumanQA != nil || state.Result.Delivery.Candidates[0].Integration != nil {
		t.Fatal("ownership release invented QA or integrated")
	}
	if _, err = WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "revoked.json", DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "late-owner", PreviousEventSHA256: state.Result.Delivery.LastEventSHA256, Phase: "working", Summary: "Stale owner", Meaning: "Synthetic revoked callback", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}})); err == nil || !strings.Contains(err.Error(), "owner_released") {
		t.Fatalf("old callback was not revoked: %v", err)
	}
	beforeRetry := ownerRecoverySnapshot(t, f)
	integrationCLIOK(t, binary, f.Parent, "", args...)
	if beforeRetry != ownerRecoverySnapshot(t, f) {
		t.Fatal("release retry wrote another event")
	}
	preview := integrationCLIOK(t, binary, f.Parent, "", "--delivery", id, "--check")
	if deliveryCLIString(t, preview, "state") != "ready" {
		t.Fatalf("recovered release did not enable human integration: %s", preview["reasons"])
	}
	// Synthetic QA in a disposable fixture is not product acceptance.
	done := integrationCLIOK(t, binary, f.Parent, "pass\n", "--delivery", id)
	if deliveryCLIString(t, done, "state") != "completed" || gitOutput(t, f.Parent, "rev-parse", "HEAD") != o.Delivery.Candidates[0].OID {
		t.Fatal("native human integration did not complete exact candidate")
	}
	if _, err := os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); !os.IsNotExist(err) {
		t.Fatal("completed integration did not clean the Task worktree")
	}
}

func TestIntegrationCLIClosedOwnerRecoveryRejectsKnownWritersAndUnknownInventory(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f, o, id := integrationOwnerRecoveryFixture(t, binary)
	source := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	for _, tc := range []struct {
		name          string
		panes, agents any
	}{
		{"unknown_panes", nil, []any{}},
		{"unknown_agents", []any{}, nil},
		{"moved_terminal", []any{map[string]any{"workspace_id": "w-fixture", "tab_id": "new-tab", "pane_id": "new-pane", "terminal_id": o.Transport.TerminalID}}, []any{}},
		{"resumed_session", []any{}, []any{map[string]any{"workspace_id": "w-fixture", "tab_id": "new-tab", "pane_id": "new-pane", "terminal_id": "new-terminal", "agent": "codex", "foreground_cwd": f.Parent, "agent_session": map[string]any{"agent": "codex", "kind": "id", "value": o.Transport.AgentSessionID}}}},
		{"other_source_writer", []any{}, []any{map[string]any{"workspace_id": "w-fixture", "tab_id": "new-tab", "pane_id": "new-pane", "terminal_id": "new-terminal", "agent": "codex", "agent_status": "idle", "foreground_cwd": source}}},
		{"source_shell", []any{map[string]any{"workspace_id": "w-fixture", "tab_id": "new-tab", "pane_id": "new-pane", "terminal_id": "new-terminal", "foreground_cwd": source}}, []any{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			workflowTestModel(t, f.workflowFixture, map[string]any{"inventory": tc.panes, "recovery_agents": tc.agents})
			before := ownerRecoverySnapshot(t, f)
			if _, _, err := integrationCLIRun(binary, f.Parent, "released\n", "release", "--delivery", id, "--recover"); err == nil {
				t.Fatal("unsafe recovery was accepted")
			}
			if before != ownerRecoverySnapshot(t, f) {
				t.Fatal("blocked recovery mutated workspace")
			}
		})
	}
}

func TestHumanOwnerRecoveryRechecksWritersAndSerializesRelease(t *testing.T) {
	f, o, _ := integrationOwnerRecoveryFixture(t, deliveryCLIBinary(t))
	c := o.Delivery.Candidates[0]
	p, err := PreviewHumanOwnerRecovery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID)
	if err != nil {
		t.Fatal(err)
	}
	in := HumanOwnerRecoveryAnswer{Confirmation: p.Confirmation, Actor: "Synthetic recovery human", Answer: "released", StartedAtUTC: "2026-10-08T00:00:00Z", CompletedAtUTC: "2026-10-08T00:00:01Z"}
	workflowTestModel(t, f.workflowFixture, map[string]any{"recovery_agents": []any{map[string]any{"workspace_id": "w-fixture", "tab_id": "new-tab", "pane_id": "new-pane", "terminal_id": "new-terminal", "foreground_cwd": c.TaskResult.SourceLocator}}})
	if _, err = RecoverHumanDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, in); err == nil {
		t.Fatal("a writer appearing after preview did not block recovery")
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"recovery_agents": []any{}})
	bad := in
	bad.Confirmation = hash([]byte("another candidate or history"))
	if _, err = RecoverHumanDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, bad); err == nil {
		t.Fatal("another preview authorized release")
	}
	var wg sync.WaitGroup
	errors := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, e := RecoverHumanDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, in)
			errors <- e
		}()
	}
	wg.Wait()
	close(errors)
	for err := range errors {
		if err != nil {
			t.Fatal(err)
		}
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || s.Result.Delivery.OwnershipRelease == nil {
		t.Fatalf("release missing: %v", err)
	}
	count := 0
	for _, e := range s.Result.Delivery.Events {
		if e.Kind == "ownership_release" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("concurrent recovery created %d release events", count)
	}
	trace, err := ReadTraceHistory(f.D, f.R.WorkspaceRoot, string(c.TaskResult.TaskID))
	if err != nil || len(trace.Runs) != 1 {
		t.Fatalf("recovery history unavailable: %v", err)
	}
	if trace.Runs[0].Coverage != "complete" {
		t.Fatalf("recovery history partial: %+v", trace.Runs[0].Reasons)
	}
	found := false
	for _, entry := range trace.Runs[0].Entries {
		if entry.Kind == "source_ownership_release" {
			found = entry.Role == "human" && entry.ActorClaim == in.Actor
		}
	}
	if !found {
		t.Fatal("trace did not attribute recovery to the actual human caller")
	}
}
