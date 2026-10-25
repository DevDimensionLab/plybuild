package taskexecute

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

// Exercise the CLI composition with a real isolated registry and modeled Herdr.
func TestPublishedContinuationRecoveryInstalledPly(t *testing.T) {
	for _, change := range []string{"path", "bytes"} {
		t.Run(change, func(t *testing.T) {
			d, start, root, _ := boundFixture(t)
			p := boundPreview(t, d, start)
			out, err := ApplyBoundStart(d, p.Request)
			if err != nil {
				t.Fatal(err)
			}
			run, err := taskrun.WorkflowShow(d, root, out.ExpectedRunID)
			if err != nil {
				t.Fatal(err)
			}
			var model struct {
				CWD string `json:"cwd"`
			}
			modelRaw, _ := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json"))
			if err = json.Unmarshal(modelRaw, &model); err != nil {
				t.Fatal(err)
			}
			if err = os.Chdir(model.CWD); err != nil {
				t.Fatal(err)
			}
			_, err = taskrun.WorkflowDeliveryReceipt(d, root, run.RunID, run.Paths.Context, taskrun.DeliveryBeginInput{Provider: "codex", NativeSession: run.Transport.AgentSessionID}, true)
			if err != nil {
				t.Fatal(err)
			}
			input := ContinuationInput{RunID: run.RunID, ContextPath: run.Paths.Context}
			fault := d
			fault.Fault = func(point string) error {
				if point == "delivery_after_continuation_proof" {
					return errors.New("lost proof publication reply")
				}
				return nil
			}
			if _, err = ContinueDelivery(fault, input); err == nil {
				t.Fatal("fault did not interrupt")
			}
			dir := filepath.Join(run.Paths.RunRoot, "delivery", "continuations", "001")
			original := recoveryFiles(t, dir)
			path, _ := d.Executable()
			raw, _ := os.ReadFile(path)
			if change == "path" {
				path = filepath.Join(root, "bin", "ply-new-path")
			} else {
				raw = append(raw, []byte("# upgraded reader\n")...)
			}
			if err = os.WriteFile(path, raw, 0700); err != nil {
				t.Fatal(err)
			}
			d.Executable = func() (string, error) { return path, nil }
			checked := input
			checked.Check = true
			preview, err := ContinueDelivery(d, checked)
			if err != nil || preview.Preview == nil || preview.Preview.Proof == nil {
				t.Fatalf("installed Ply check: %+v %v", preview, err)
			}
			recovered, err := ContinueDelivery(d, input)
			if err != nil || recovered.State != "continued" {
				t.Fatalf("installed Ply recovery: %+v %v", recovered, err)
			}
			if after := recoveryFiles(t, dir); !reflect.DeepEqual(original, after) {
				t.Fatal("recovery rewrote original proof/control")
			}
		})
	}
}

// Recovery through the ordinary executable and Cobra command, with no saved
// controller invocation. Only Herdr/provider observation is synthetic.
func TestPublishedContinuationRecoveryPublicCLI(t *testing.T) {
	source, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	binary := os.Getenv("PLY_CANDIDATE_BINARY")
	if binary == "" {
		binary = filepath.Join(physicalTemp(t), "ply")
		build := exec.Command("go", "build", "-o", binary, ".")
		build.Dir = source
		if raw, err := build.CombinedOutput(); err != nil {
			t.Fatalf("build CLI: %s %v", raw, err)
		}
	}
	d, start, root, _ := boundFixture(t)
	p := boundPreview(t, d, start)
	out, err := ApplyBoundStart(d, p.Request)
	if err != nil {
		t.Fatal(err)
	}
	run, err := taskrun.WorkflowShow(d, root, out.ExpectedRunID)
	if err != nil {
		t.Fatal(err)
	}
	var model struct {
		CWD                   string `json:"cwd"`
		Tabs, Starts, Prompts int
	}
	modelPath := filepath.Join(root, "bin", "recovery-model.json")
	modelRaw, _ := os.ReadFile(modelPath)
	if err = json.Unmarshal(modelRaw, &model); err != nil {
		t.Fatal(err)
	}
	if err = os.Chdir(model.CWD); err != nil {
		t.Fatal(err)
	}
	if _, err = taskrun.WorkflowDeliveryReceipt(d, root, run.RunID, run.Paths.Context, taskrun.DeliveryBeginInput{Provider: "codex", NativeSession: run.Transport.AgentSessionID}, true); err != nil {
		t.Fatal(err)
	}
	fault := d
	fault.Fault = func(point string) error {
		if point == "delivery_after_continuation_proof" {
			return errors.New("lost publication reply")
		}
		return nil
	}
	if _, err = ContinueDelivery(fault, ContinuationInput{RunID: run.RunID, ContextPath: run.Paths.Context}); err == nil {
		t.Fatal("fault did not interrupt")
	}
	dir := filepath.Join(run.Paths.RunRoot, "delivery", "continuations", "001")
	original := recoveryFiles(t, dir)
	call := func(check bool) ContinuationResult {
		t.Helper()
		args := []string{"workflow", "execute", "continue", run.RunID, "--context", run.Paths.Context, "--format", "json"}
		if check {
			args = append(args, "--check")
		}
		command := exec.Command(binary, args...)
		command.Dir = model.CWD
		raw, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("installed CLI: %s %v", raw, err)
		}
		var result ContinuationResult
		if err = json.Unmarshal(raw, &result); err != nil {
			t.Fatalf("CLI JSON: %s %v", raw, err)
		}
		return result
	}
	if got := call(true); got.State != "ready" || got.Preview == nil || got.Preview.Proof == nil {
		t.Fatalf("check: %+v", got)
	}
	if got := call(false); got.State != "continued" || got.Preview.Generation != 1 {
		t.Fatalf("recover: %+v", got)
	}
	// The installed binary has a different digest: after pointer publication,
	// --check describes a new transition under the existing rules. Do not apply it.
	if got := call(true); got.State != "ready" || got.Preview.Generation != 2 || got.Preview.Proof != nil {
		t.Fatalf("post-pointer check: %+v", got)
	}
	if !reflect.DeepEqual(original, recoveryFiles(t, dir)) {
		t.Fatal("CLI rewrote original proof/control")
	}
	var after struct{ Tabs, Starts, Prompts int }
	modelRaw, _ = os.ReadFile(modelPath)
	if err = json.Unmarshal(modelRaw, &after); err != nil {
		t.Fatal(err)
	}
	if model.Tabs != after.Tabs || model.Starts != after.Starts || model.Prompts != after.Prompts {
		t.Fatal("CLI caused new native effects")
	}
}
