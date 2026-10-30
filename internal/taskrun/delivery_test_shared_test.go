package taskrun

import (
	"bytes"
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

func deliveryCLIBinary(t *testing.T) string {
	t.Helper()
	if binary := os.Getenv("PLY_DELIVERY_TEST_BINARY"); binary != "" {
		if !filepath.IsAbs(binary) {
			t.Fatal("PLY_DELIVERY_TEST_BINARY must be absolute")
		}
		return binary
	}
	_, source, _, ok := runtime.Caller(0)
	if !ok || !filepath.IsAbs(source) {
		t.Fatal("CLI helper requires an absolute source filename or PLY_DELIVERY_TEST_BINARY")
	}
	repo := filepath.Clean(filepath.Join(filepath.Dir(source), "../.."))
	binary := filepath.Join(t.TempDir(), "ply")
	c := exec.Command("go", "build", "-o", binary, "./cmd/ply")
	c.Dir = repo
	if output, err := c.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, output)
	}
	return binary
}

func automaticCLIBinary(t *testing.T) string {
	t.Helper()
	path, err := filepath.EvalSymlinks(deliveryCLIBinary(t))
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func automaticCLI(binary string, f deliveryFixture, args ...string) (WorkflowRun, []byte, error) {
	c := exec.Command(binary, append([]string{"workflow"}, append(args, "--format", "json")...)...)
	c.Dir = f.Prepared.Preparation.Preparation.Plan.WorktreePath
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run()
	var o WorkflowRun
	if decodeErr := json.Unmarshal(stdout.Bytes(), &o); decodeErr != nil && err == nil {
		err = decodeErr
	}
	return o, append(stdout.Bytes(), stderr.Bytes()...), err
}

func automaticCLIOK(t *testing.T, binary string, f deliveryFixture, args ...string) WorkflowRun {
	t.Helper()
	o, output, err := automaticCLI(binary, f, args...)
	if err != nil {
		t.Fatalf("CLI %v: %v\n%s", args, err, output)
	}
	return o
}

func ownerPaneStandin(t *testing.T, f *deliveryFixture) {
	t.Helper()
	herdr, err := os.ReadFile(f.R.Herdr.Executable.Path)
	if err != nil {
		t.Fatal(err)
	}
	herdr = bytes.Replace(herdr, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+al21CloseoutRaceStandin+retireOpenPaneStandin+closedTransportInspectStandin+integrationOwnerRecoveryStandin), 1)
	if err = os.WriteFile(f.R.Herdr.Executable.Path, herdr, 0700); err != nil {
		t.Fatal(err)
	}
	f.R.Herdr.Executable.SHA256 = hash(herdr)
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
}

// cleanupOwnerLeaves models the released owner's actual exit: its bound pane,
// managed agent and Task processes are gone (requires ownerPaneStandin).

func cleanupOwnerLeaves(t *testing.T, f deliveryFixture) {
	t.Helper()
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true})
}

// cleanupOwnerHandoff asserts the owner's integrate ended in the CL02 handoff
// with the integration observed and the source still in place.

const integrationOwnerRecoveryStandin = `
if model.get("closed_transport") and cmd == ["agent", "get"]:
    print(json.dumps({"id":"fixture", "error":{"code":"agent_not_found", "message":"closed fixture agent"}}))
    raise SystemExit(1)
if model.get("closed_transport") and cmd == ["agent", "list"]:
    print(json.dumps({"id":"fixture", "result":{"agents":model.get("recovery_agents", [])}}))
    raise SystemExit(0)
`

func deliveryCLISnapshot(t *testing.T, root string) string {
	t.Helper()
	var entries []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		entries = append(entries, path+":"+hash(data))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(entries)
	return strings.Join(entries, "\n")
}

func autoScratchFixture(t *testing.T, script string) (deliveryFixture, WorkflowRun, string, string) {
	t.Helper()
	agreement := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "fixture orchestrator", TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &agreement})
	ownerPaneStandin(t, &f)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "scratch fixture candidate")
	binary := filepath.Join(f.R.WorkspaceRoot, "candidate-executable")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\nprintf 'fixture candidate\\n'\n"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	return f, o, binary, deliveryTestReview(t, f, "scratch")
}

func automaticTestReceipt(t *testing.T, o WorkflowRun) deliveryVerificationReceipt {
	t.Helper()
	var r deliveryVerificationReceipt
	raw, err := os.ReadFile(filepath.Join(o.Delivery.Attempt.Path, "verification.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &r); err != nil {
		t.Fatal(err)
	}
	return r
}

func deliveryCLILatestBase(base workspace.EpicBaseResult) workspace.EpicBaseVersion {
	var latest workspace.EpicBaseVersion
	for _, v := range base.Versions {
		if v.Revision > latest.Revision {
			latest = v
		}
	}
	return latest
}
