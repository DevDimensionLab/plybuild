package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/intelligence"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func runIntelligence(t *testing.T, cwd, stdin string, args ...string) (string, error) {
	t.Helper()
	t.Chdir(cwd)
	command := newIntelligenceCommand(workspace.SystemDependencies())
	command.SilenceUsage, command.SilenceErrors = true, true // as RootCmd under ExecuteE
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&bytes.Buffer{})
	command.SetIn(strings.NewReader(stdin))
	command.SetArgs(args)
	err := command.Execute()
	return stdout.String(), err
}

// An agent drafts, fills and records an observation through ply, and reads it
// back with its hash; the post lives under <wrapper>/data/intelligence.
func TestIntelligenceDraftObserveListShow(t *testing.T) {
	root := newAssignmentWorkspace(t)
	draft, err := runIntelligence(t, root, "", "draft", "--project", "a", "--kind", "observation",
		"--repo", "a-app", "--activity", "PLAN-1/s1", "--role", "author", "--provider", "claude")
	if err != nil {
		t.Fatal(err)
	}
	var post map[string]any
	if err := json.Unmarshal([]byte(draft), &post); err != nil || post["schema"] != "ply.observation/v1" {
		t.Fatalf("draft %s %v", draft, err)
	}
	if origin := post["origin"].(map[string]any); origin["project"] != "a" || origin["repo"] != "a-app" {
		t.Fatalf("draft origin %v", origin)
	}
	post["situation"], post["observation"], post["limitations"] = "Testing.", "It works.", "One run."
	post["evidence"] = []any{map[string]any{"basis": "direct", "ref": "cmd test", "detail": "Ran it."}}
	filled, _ := json.Marshal(post)

	recorded, err := runIntelligence(t, root, string(filled), "observe", "--project", "a", "--file", "-", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	var receipt struct {
		Kind    string               `json:"kind"`
		Receipt intelligence.Receipt `json:"receipt"`
	}
	if err := json.Unmarshal([]byte(recorded), &receipt); err != nil || receipt.Kind != "PlyIntelligenceReceipt@1" || receipt.Receipt.State != intelligence.StateRecorded {
		t.Fatalf("observe %s %v", recorded, err)
	}
	id := receipt.Receipt.ID
	if _, err := os.Stat(filepath.Join(root, "a", "data", "intelligence", "observations", id+".json")); err != nil {
		t.Fatalf("post not under data/intelligence: %v", err)
	}
	draftFile := filepath.Join(t.TempDir(), "draft.json")
	if err := os.WriteFile(draftFile, filled, 0o600); err != nil {
		t.Fatal(err)
	}
	if again, err := runIntelligence(t, root, "", "observe", "--project", "a", "--file", draftFile); err != nil || again != "already_recorded "+id+" "+receipt.Receipt.SHA256+"\n" {
		t.Fatalf("retry %q %v", again, err)
	}

	listed, err := runIntelligence(t, root, "", "observations", "--project", "a", "--repo", "a-app", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	var list struct {
		Kind  string              `json:"kind"`
		Items []intelligence.Item `json:"items"`
	}
	if err := json.Unmarshal([]byte(listed), &list); err != nil || list.Kind != "PlyIntelligenceList@1" || len(list.Items) != 1 ||
		list.Items[0].ID != id || list.Items[0].SHA256 != receipt.Receipt.SHA256 {
		t.Fatalf("list %s %v", listed, err)
	}
	if text, err := runIntelligence(t, root, "", "observations", "--project", "a"); err != nil || !strings.HasPrefix(text, id+" ") {
		t.Fatalf("list text %q %v", text, err)
	}
	shown, err := runIntelligence(t, root, "", "show", "--project", "a", id)
	if err != nil || !strings.Contains(shown, `"observation": "It works."`) {
		t.Fatalf("show %q %v", shown, err)
	}

	_, err = runIntelligence(t, root, "", "draft", "--project", "a", "--kind", "observation",
		"--repo", "b-app", "--activity", "x", "--role", "author", "--provider", "claude")
	var failure *intelligence.Error
	if !errors.As(err, &failure) || failure.Class != intelligence.ClassInvalidArguments {
		t.Fatalf("repo of another project: %v", err)
	}
}
