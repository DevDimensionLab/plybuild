package taskrun

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWorkflowStandinInterruptedModelPublicationPreservesOriginal(t *testing.T) {
	python, err := exec.LookPath("python3")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	model := filepath.Join(root, "herdr-model.json")
	original := []byte(`{"agent_status":"idle","agent_session_id":"fixture-session","label":"original"}`)
	if err = os.WriteFile(model, original, 0600); err != nil {
		t.Fatal(err)
	}
	// Pause this temporary stand-in after its model write opens, when an
	// in-place writer has already truncated the published file. Waiting for a
	// marker makes cancellation deterministic; no production hook or timing
	// race is needed. A sibling publication file is paused at the same point.
	const barrier = `import io
model_open = io.open
def pause_model_write(file, mode="r", *args, **kwargs):
    stream = model_open(file, mode, *args, **kwargs)
    if "w" in mode and Path(file).name.startswith("herdr-model.json"):
        print("MODEL_WRITE_OPENED", flush=True)
        sys.stdin.buffer.read(1)
    return stream
io.open = pause_model_write
`
	const importLine = "from pathlib import Path\n"
	if strings.Count(workflowStandin, importLine) != 1 {
		t.Fatal("stand-in import boundary is missing or ambiguous")
	}
	interruptedScript := filepath.Join(root, "interrupted-standin.py")
	if err = os.WriteFile(interruptedScript, []byte(strings.Replace(workflowStandin, importLine, importLine+barrier, 1)), 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, python, interruptedScript, "tab", "rename", "tab-fixture", "changed")
	stdout, err := command.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdin, err := command.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	defer stdin.Close()
	var stderr bytes.Buffer
	command.Stderr = &stderr
	if err = command.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if command.ProcessState == nil {
			_ = command.Process.Kill()
			_ = command.Wait()
		}
	}()
	marker, err := bufio.NewReader(stdout).ReadString('\n')
	if err != nil || marker != "MODEL_WRITE_OPENED\n" {
		_ = command.Process.Kill()
		_ = command.Wait()
		t.Fatalf("stand-in did not reach its publication boundary: marker=%q err=%v stderr=%s", marker, err, stderr.String())
	}
	if err = command.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	if err = command.Wait(); err == nil {
		t.Fatal("interrupted stand-in unexpectedly completed")
	}
	after, err := os.ReadFile(model)
	if err != nil || !json.Valid(after) || !bytes.Equal(after, original) {
		t.Fatalf("interrupted stand-in corrupted its published model: bytes=%d valid_json=%t err=%v", len(after), json.Valid(after), err)
	}

	// A completed command must publish its new model, so retaining the old
	// file by dropping every update would not satisfy the regression.
	completeScript := filepath.Join(root, "complete-standin.py")
	if err = os.WriteFile(completeScript, []byte(workflowStandin), 0600); err != nil {
		t.Fatal(err)
	}
	completeCtx, completeCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer completeCancel()
	completed := exec.CommandContext(completeCtx, python, completeScript, "tab", "rename", "tab-fixture", "changed")
	if output, err := completed.CombinedOutput(); err != nil {
		t.Fatalf("complete stand-in failed: %v\n%s", err, output)
	}
	after, err = os.ReadFile(model)
	var saved struct {
		Label string `json:"label"`
	}
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(after, &saved); err != nil || saved.Label != "changed" {
		t.Fatalf("completed stand-in did not publish its new model: %+v %v", saved, err)
	}
}
