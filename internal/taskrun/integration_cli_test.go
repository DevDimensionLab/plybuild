package taskrun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func integrationCLIRun(binary, cwd, input string, args ...string) (map[string]json.RawMessage, []byte, error) {
	c := exec.Command(binary, append([]string{"integration"}, append(args, "--format", "json")...)...)
	c.Dir = cwd
	c.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	var result map[string]json.RawMessage
	if stdout.Len() > 0 {
		if e := json.Unmarshal(stdout.Bytes(), &result); e != nil {
			return nil, append(stdout.Bytes(), stderr.Bytes()...), fmt.Errorf("invalid integration JSON: %w; exit %v", e, err)
		}
	}
	return result, append(stdout.Bytes(), stderr.Bytes()...), err
}
func integrationCLIOK(t *testing.T, binary, cwd, input string, args ...string) map[string]json.RawMessage {
	t.Helper()
	r, output, e := integrationCLIRun(binary, cwd, input, args...)
	if e != nil {
		if r != nil {
			detail := ""
			if n := strings.LastIndex(string(output), "\nError:"); n >= 0 {
				detail = string(output[n:])
			}
			t.Fatalf("integration %v: %v state=%s reasons=%s next=%s%s", args, e, r["state"], r["reasons"], r["next_action"], detail)
		}
		t.Fatalf("integration %v: %v\n%s", args, e, output)
	}
	return r
}
func integrationCLIFixture(t *testing.T, binary, parent string, human, remote bool, fixtureRoot ...string) (deliveryFixture, WorkflowRun, string) {
	t.Helper()
	mode := workspace.DeliveryLocalBranch
	if parent == "epic" {
		mode = workspace.DeliveryLocalEpic
	}
	a := workspace.DeliveryAgreement{Mode: mode}
	if human {
		a.SchemaVersion = 2
		a.IntegrationOwner = workspace.IntegrationOwnerHuman
	}
	root := ""
	if len(fixtureRoot) > 0 {
		root = fixtureRoot[0]
	}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &a, ParentRef: parent, Root: root})
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	if remote {
		runGit(t, f.Parent, "remote", "add", "origin", filepath.Join(f.R.WorkspaceRoot, "nonexistent-remote-no-network"))
	}
	o := deliveryCLIQualified(t, f)
	source := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	id := deliveryCLIString(t, deliveryCLIOK(t, binary, source, "register", "--file", deliveryCLIRegisterInput(t, f, o, "human-integration")), "id")
	var e error
	o, e = ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic fixture owner releases writing; actual human product acceptance is not claimed.")
	if e != nil {
		t.Fatal(e)
	}
	return f, o, id
}

func TestIntegrationCLIHumanLocalReturnAndReadOnlyChecks(t *testing.T) {
	binary := deliveryCLIBinary(t)
	for _, tc := range []struct {
		parent              string
		human, remote, keep bool
	}{{"main", true, false, true}, {"master", false, true, false}, {"epic", true, false, false}} {
		t.Run(tc.parent, func(t *testing.T) {
			f, o, id := integrationCLIFixture(t, binary, tc.parent, tc.human, tc.remote)
			source := f.Prepared.Preparation.Preparation.Plan.WorktreePath
			candidate := o.Delivery.Candidates[0].OID
			parentBefore := gitOutput(t, f.Parent, "rev-parse", "HEAD")
			args := []string{"--delivery", id}
			if tc.keep {
				args = append(args, "--keep")
			}
			before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
			p := integrationCLIOK(t, binary, f.Parent, "", append(args, "--check")...)
			if state := deliveryCLIString(t, p, "state"); state != "ready" {
				t.Fatalf("preview not ready: %s\n%s", state, p["reasons"])
			}
			integrationCLIOK(t, binary, f.Parent, "", "list", "--project", "ply")
			integrationCLIOK(t, binary, f.Parent, "", "scan", "--project", "ply")
			if _, output, e := integrationCLIRun(binary, f.Parent, "", args...); e == nil || !strings.Contains(string(output), "input ended") {
				t.Fatalf("EOF was accepted: %v %s", e, output)
			}
			if after := deliveryCLISnapshot(t, f.R.WorkspaceRoot); after != before {
				t.Fatal("read-only plan/list/scan or EOF mutated workspace")
			}
			r := integrationCLIOK(t, binary, f.Parent, "pass\n", args...)
			if deliveryCLIString(t, r, "state") != "completed" {
				t.Fatalf("incomplete receipt: %s", r)
			}
			if got := gitOutput(t, f.Parent, "rev-parse", "HEAD"); got != candidate || got == parentBefore {
				t.Fatalf("wrong local target: %s", got)
			}
			integrationID := deliveryCLIString(t, r, "id")
			integrationCLIOK(t, binary, f.Parent, "", "resume", integrationID)
			integrationCLIOK(t, binary, f.Parent, "", "show", integrationID)
			if tc.keep {
				if _, e := os.Stat(source); e != nil {
					t.Fatal("keep removed source")
				}
			} else {
				if _, e := os.Stat(source); !os.IsNotExist(e) {
					t.Fatalf("source was not retired: %v", e)
				}
			}
			closed, e := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, "task")
			if e != nil || closed == nil || closed.State != "complete" {
				t.Fatalf("closeout missing: %+v %v", closed, e)
			}
			if tc.remote {
				if got := gitOutput(t, f.Parent, "remote", "get-url", "origin"); !strings.Contains(got, "nonexistent-remote-no-network") {
					t.Fatal("local integration altered remote")
				}
			}
		})
	}
}

func TestIntegrationCLINegativeAnswerAndLatestQABlockOldPass(t *testing.T) {
	binary := deliveryCLIBinary(t)
	for _, answer := range []string{"fail", "blocked"} {
		t.Run(answer, func(t *testing.T) {
			f, _, id := integrationCLIFixture(t, binary, "main", true, false)
			before := gitOutput(t, f.Parent, "rev-parse", "HEAD")
			r, output, e := integrationCLIRun(binary, f.Parent, answer+"\n", "--delivery", id, "--keep")
			if e == nil || r == nil || deliveryCLIString(t, r, "state") != "human_"+answer {
				t.Fatalf("negative answer not preserved: %v %s", e, output)
			}
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") != before {
				t.Fatal("negative answer integrated")
			}
			integrationID := deliveryCLIString(t, r, "id")
			if _, _, e = integrationCLIRun(binary, f.Parent, "", "resume", integrationID); e == nil {
				t.Fatal("resume treated negative answer as pass")
			}
		})
	}
}

func TestIntegrationCLIConcurrentSameCandidateClosesOnce(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f, _, id := integrationCLIFixture(t, binary, "main", true, false)
	var wg sync.WaitGroup
	type result struct {
		row map[string]json.RawMessage
		out []byte
		err error
	}
	results := make([]result, 2)
	for i := range results {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i].row, results[i].out, results[i].err = integrationCLIRun(binary, f.Parent, "pass\n", "--delivery", id, "--keep")
		}(i)
	}
	wg.Wait()
	completed := 0
	for _, r := range results {
		if r.err == nil && deliveryCLIString(t, r.row, "state") == "completed" {
			completed++
		}
	}
	if completed != 1 {
		t.Fatalf("expected one integration owner, got %d: %s %s", completed, results[0].out, results[1].out)
	}
	registry, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil {
		t.Fatal(e)
	}
	if len(registry.IntegrationAttempts) != 1 {
		t.Fatalf("integration repeated: %d attempts", len(registry.IntegrationAttempts))
	}
}

func TestIntegrationCLIRejectsNotificationControlFileInsideCleanupSource(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{Mode: workspace.DeliveryLocalBranch}, ParentRef: "main"})
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	source := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	route := writeAny(t, source, "route.json", map[string]any{"kind": "ply.workflow.notification-route", "schema_version": 1, "name": "fixture", "transport": "slack_incoming_webhook", "channel_label": "Fixture only", "webhook_env": "UNSET_FIXTURE_WEBHOOK", "state_root": filepath.Join(f.R.WorkspaceRoot, "notifications")})
	runGit(t, source, "add", "route.json")
	deliveryTestCandidate(t, f, "notification control retention")
	var e error
	o, e = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "route"))
	if e != nil {
		t.Fatal(e)
	}
	id := deliveryCLIString(t, deliveryCLIOK(t, binary, source, "register", "--file", deliveryCLIRegisterInput(t, f, o, "route")), "id")
	if _, e = ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic source owner release."); e != nil {
		t.Fatal(e)
	}
	p := integrationCLIOK(t, binary, f.Parent, "", "--delivery", id, "--notification-route", route, "--check")
	if deliveryCLIString(t, p, "state") != "blocked" || !bytes.Contains(p["reasons"], []byte("control_input_inside_source")) {
		t.Fatalf("plan accepts deleting its notification control: state=%s reasons=%s", p["state"], p["reasons"])
	}
}

func TestIntegrationCLILegacyReconcileCannotBypassPendingInstallation(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f, o, id := integrationCLIFixture(t, binary, "main", true, false)
	falseBin, e := exec.LookPath("false")
	if e != nil {
		t.Fatal(e)
	}
	trueBin, e := exec.LookPath("true")
	if e != nil {
		t.Fatal(e)
	}
	profile := writeAny(t, f.R.WorkspaceRoot, "install.json", map[string]any{"kind": "ply.integration.install-profile", "schema_version": 1, "name": "fixture", "candidate_oid": o.Delivery.Candidates[0].OID, "artifact_path": filepath.Join(f.R.WorkspaceRoot, "installed"), "artifact_sha256": hash([]byte("expected artifact")), "install": map[string]any{"executable": falseBin, "arguments": []string{}, "cwd": f.Parent}, "verify": map[string]any{"executable": trueBin, "arguments": []string{}, "cwd": f.Parent}})
	r, output, e := integrationCLIRun(binary, f.Parent, "pass\n", "--delivery", id, "--install-profile", profile)
	if e == nil || r == nil || deliveryCLIString(t, r, "state") != "install_failed" {
		t.Fatalf("did not reach known failed install: %v %s", e, output)
	}
	operation := deliveryCLIString(t, r, "id")
	if _, output, e = integrationCLIRun(binary, f.Parent, "released\n", "release", "--task", "task"); e == nil || !strings.Contains(string(output), operation) {
		t.Fatalf("legacy release bypasses active integration %s: %v %s", operation, e, output)
	}
	scan := integrationCLIOK(t, binary, f.Parent, "", "scan", "--project", "ply")
	if !bytes.Contains(scan["entries"], []byte("install_failed")) || bytes.Contains(scan["entries"], []byte("--reconcile")) {
		t.Fatalf("scan recommends bypassing pending installation: %s", scan)
	}
}
