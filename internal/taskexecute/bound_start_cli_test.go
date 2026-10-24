package taskexecute

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestBoundStartPublicCLIContract(t *testing.T) {
	source, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	compiled := os.Getenv("PLY_CANDIDATE_BINARY")
	if compiled == "" {
		compiled = os.Getenv("PLY_DELIVERY_TEST_BINARY")
	}
	if compiled == "" {
		compiled = filepath.Join(physicalTemp(t), "ply")
		build := exec.Command("go", "build", "-o", compiled, ".")
		build.Dir = source
		if output, err := build.CombinedOutput(); err != nil {
			t.Fatalf("build public CLI: %s %v", output, err)
		}
	} else if !filepath.IsAbs(compiled) {
		t.Fatal("candidate binary must be absolute")
	}
	_, _, root, epic := boundFixture(t)
	responses := filepath.Join(root, "responses")
	if err := os.Mkdir(responses, 0700); err != nil {
		t.Fatal(err)
	}
	call := func(name string, wantSuccess bool, args ...string) []byte {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		defer cancel()
		command := exec.CommandContext(ctx, compiled, args...)
		command.Dir = epic
		var stderr strings.Builder
		command.Stderr = &stderr
		output, err := command.Output()
		if (err == nil) != wantSuccess {
			t.Fatalf("CLI %v: %v\nstdout=%s\nstderr=%s", args, err, output, stderr.String())
		}
		if name != "" {
			if !json.Valid(output) {
				t.Fatalf("CLI did not return JSON: %s", output)
			}
			if err := os.WriteFile(filepath.Join(responses, name+".json"), output, 0600); err != nil {
				t.Fatal(err)
			}
		}
		return output
	}
	capabilities := call("capabilities", true, "capabilities", "--format", "json")
	var catalog struct {
		Operations []any `json:"operations"`
		Workflow   []struct {
			ID, Mode      string
			Command       []string
			ResultSchemas []struct {
				Kind    string
				Version int `json:"schema_version"`
			} `json:"result_schemas"`
		} `json:"workflow_extensions"`
	}
	if err := json.Unmarshal(capabilities, &catalog); err != nil {
		t.Fatal(err)
	}
	if len(catalog.Operations) != 5 {
		t.Fatal("changed the original core capability catalog")
	}
	found := map[string]bool{}
	for _, operation := range catalog.Workflow {
		if strings.HasPrefix(operation.ID, "workflow.execute.bound-start") {
			found[operation.ID] = true
		}
	}
	if len(found) != 3 {
		t.Fatalf("missing discoverable bound contract: %s", capabilities)
	}
	before := recoveryFiles(t, filepath.Join(root, ".ply"))
	raw := call("preview", true, "workflow", "execute", "preview", "--spec", "observable-goal", "--request-id", "public-cli-request-0001", "--herdr-workspace", "w-fixture", "--format", "json")
	var preview BoundStartPreview
	if err := json.Unmarshal(raw, &preview); err != nil {
		t.Fatal(err)
	}
	call("status-absent", true, "workflow", "execute", "status", preview.Request.RequestID, "--workspace", root, "--format", "json")
	call("ordinary-check", true, "workflow", "execute", "--spec", "observable-goal", "--check", "--format", "json")
	if !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatal("public preview/status/check wrote native state")
	}
	request := filepath.Join(responses, "request.json")
	preview.Request.Binding = "sha256:" + strings.Repeat("0", 64)
	fixtureJSON(t, request, preview.Request)
	changed := preview.Request
	changed.Input.SpecID = "absent-goal"
	bad := filepath.Join(root, "changed-request.json")
	fixtureJSON(t, bad, changed)
	rejected := call("rejected", false, "workflow", "execute", "apply", "--file", bad, "--format", "json")
	var rejection BoundStartStatus
	if err := json.Unmarshal(rejected, &rejection); err != nil || rejection.Error == nil || rejection.Error.Code != "bound_start_stale" || rejection.ExecutionEffect != "none" {
		t.Fatalf("unstable JSON rejection: %s %v", rejected, err)
	}
	if !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatal("rejected public apply wrote native state")
	}
	requestBytes, _ := os.ReadFile(request)
	for _, change := range []string{"duplicate", "unknown"} {
		prefix := `{"kind":"ply.workflow.bound-start-request",`
		if change == "unknown" {
			prefix = `{"unexpected":true,`
		}
		path := filepath.Join(root, "invalid-"+change+".json")
		if err := os.WriteFile(path, append([]byte(prefix), requestBytes[1:]...), 0600); err != nil {
			t.Fatal(err)
		}
		command := exec.Command(compiled, "workflow", "execute", "apply", "--file", path, "--format", "json")
		command.Dir = epic
		output, err := command.CombinedOutput()
		if err == nil || !strings.Contains(string(output), "bound_start_invalid_request") {
			t.Fatalf("public parser accepted %s: %s %v", change, output, err)
		}
	}
	if !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatal("invalid public request wrote native state")
	}
	var simultaneous sync.WaitGroup
	for i := 0; i < 8; i++ {
		simultaneous.Add(1)
		go func(n int) {
			defer simultaneous.Done()
			call("apply-"+strconv.Itoa(n), true, "workflow", "execute", "apply", "--file", request, "--format", "json")
		}(i)
	}
	simultaneous.Wait()
	call("status", true, "workflow", "execute", "status", preview.Request.RequestID, "--workspace", root, "--format", "json")
	call("replay", true, "workflow", "execute", "apply", "--file", request, "--format", "json")
	for _, verb := range []string{"preview", "apply", "status"} {
		call("", true, "workflow", "execute", verb, "--help")
	}
	call("", false, "workflow", "execute", "preview", "--spec", "observable-goal", "--request-id", "unsupported-mode", "--start-mode", "headless", "--herdr-workspace", "w-fixture", "--format", "json")
	call("", false, "workflow", "execute", "apply", "--file", request, "--unknown-bound-option")
	var model struct{ Tabs, Starts, Prompts int }
	transport, _ := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json"))
	if json.Unmarshal(transport, &model) != nil || model.Tabs != 1 || model.Starts != 1 || model.Prompts != 1 {
		t.Fatalf("public CLI duplicated transport: %s", transport)
	}
	if python := os.Getenv("PLY_BOUND_START_SCHEMA_PYTHON"); python != "" {
		command := exec.Command(python, filepath.Join(source, "test/bound_start_schema_test.py"), "--responses", responses)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("actual public JSON violates delivered schema: %s %v", output, err)
		} else {
			t.Log(strings.TrimSpace(string(output)))
		}
	} else {
		t.Log("JSON Schema validation is a separate check; set PLY_BOUND_START_SCHEMA_PYTHON to its test-only venv interpreter")
	}
}
