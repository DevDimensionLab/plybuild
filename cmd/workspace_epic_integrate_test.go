package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestEpicIntegratePassesExactInputAndRendersJSON(t *testing.T) {
	var got workspace.EpicIntegrateInput
	c := newWorkspaceEpicIntegrateCommandWith(func(input workspace.EpicIntegrateInput) (workspace.EpicIntegrateResult, error) {
		got = input
		return workspace.EpicIntegrateResult{EpicID: input.EpicID, Applied: true, Effect: "fast_forwarded"}, nil
	})
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&bytes.Buffer{})
	a, b := strings.Repeat("a", 40), strings.Repeat("b", 40)
	c.SetArgs([]string{"code-map-d1", "--repo", "ply-code-map", "--target-ref", "refs/heads/main", "--expected-target-oid", a, "--expected-epic-oid", b, "--apply", "--format", "json"})
	if err := c.Execute(); err != nil {
		t.Fatal(err)
	}
	want := workspace.EpicIntegrateInput{EpicID: "code-map-d1", RepoID: "ply-code-map", TargetRef: "refs/heads/main", ExpectedTargetOID: a, ExpectedEpicOID: b, Apply: true}
	if got != want || !strings.Contains(out.String(), `"kind": "WorkspaceEpicIntegrate@1"`) {
		t.Fatalf("input %+v output %s", got, out.String())
	}
	called := false
	c = newWorkspaceEpicIntegrateCommandWith(func(workspace.EpicIntegrateInput) (workspace.EpicIntegrateResult, error) {
		called = true
		return workspace.EpicIntegrateResult{}, nil
	})
	c.SetOut(&bytes.Buffer{})
	c.SetErr(&bytes.Buffer{})
	c.SetArgs([]string{"code-map-d1", "--repo", "ply-code-map", "--target-ref", "refs/heads/main", "--expected-target-oid", a})
	if err := c.Execute(); err == nil || called {
		t.Fatal("missing --expected-epic-oid accepted")
	}
}
