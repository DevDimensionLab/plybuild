package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/agentwork"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAgentWorkCLIReadAndErrors(t *testing.T) {
	root := &cobra.Command{Use: "ply"}
	root.AddCommand(newAgentWorkCommand(func(string) (*agentwork.Service, error) { return nil, errAgentFixture{} }))
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&bytes.Buffer{})
	root.SetArgs([]string{"agents", "--format", "json"})
	if e := root.Execute(); e == nil {
		t.Fatal("expected fixture failure")
	}
	if out.Len() != 0 {
		t.Fatal("error must not leak source bytes")
	}
}

type errAgentFixture struct{}

func (errAgentFixture) Error() string { return "fixture unavailable" }
func TestAgentWorkCapabilityDiscovery(t *testing.T) {
	c := coreCapabilities(capabilityBuild{})
	found := false
	for _, op := range c.ReadExtensions {
		if op.ID == "workflow.agents" {
			found = true
			b, _ := json.Marshal(op)
			if len(b) == 0 {
				t.Fatal("empty")
			}
		}
	}
	if !found {
		t.Fatal("agent work read missing")
	}
}
func TestAgentWorkStrictRequestDecoder(t *testing.T) {
	for _, data := range []string{`{"owner_ref":"a","owner_ref":"b"}`, `{"unknown":1}`, `{"owner_ref":"\ud800"}`, `{} {}`, strings.Repeat("x", 16385)} {
		p := filepath.Join(t.TempDir(), "input.json")
		if e := os.WriteFile(p, []byte(data), 0600); e != nil {
			t.Fatal(e)
		}
		var d agentwork.Draft
		if e := readAgentWorkInput(p, &d); e == nil {
			t.Fatalf("invalid input accepted: %.60s", data)
		}
	}
}
func TestAgentWorkMutationSchemaRejectsExplicitNull(t *testing.T) {
	p := filepath.Join(t.TempDir(), "mutation.json")
	if e := os.WriteFile(p, []byte(`{"request_id":"request","owner_ref":"owner","expected_revision":1,"generation":1,"summary":null}`), 0600); e != nil {
		t.Fatal(e)
	}
	var m agentwork.Mutation
	if e := readAgentWorkInput(p, &m); e == nil {
		t.Fatal("schema-disallowed null accepted as omitted report text")
	}
}
func TestAgentWorkMutationCapabilitiesAdvertiseRequiredInputs(t *testing.T) {
	c := coreCapabilities(capabilityBuild{})
	for _, op := range c.WorkflowExtensions {
		want := ""
		switch op.ID {
		case "workflow.followup.register", "workflow.followup.request", "workflow.followup.cancel", "workflow.followup.recheck", "workflow.followup.retarget":
			want = "--file"
		case "workflow.followup.list":
			want = "--due"
		case "workflow.followup.request-status":
			want = "--request-id"
		}
		if want != "" && !strings.Contains(strings.Join(op.Selectors, " "), want) {
			t.Errorf("%s omits %s", op.ID, want)
		}
	}
}
