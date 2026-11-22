package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/agentstart"
	"github.com/devdimensionlab/plybuild/internal/knowledge"
)

// sharedExperienceWorkspace is knowledgeWorkspace with ten filler posts, so a
// word in one or two posts is not common.
func sharedExperienceWorkspace(t *testing.T, pointer string) string {
	t.Helper()
	root := knowledgeWorkspace(t, pointer)
	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("k-20260200-filler-%02d", i)
		text := "---\nid: \"" + id + "\"\ntitle: \"Fyll\"\nstatus: \"confirmed\"\nscope: {\"project\": \"*\", \"repos\": [], \"roles\": []}\n---\n\nFyll og i plan.\n"
		path := filepath.Join(root, "a", "lib", "knowledge", "filler", id+".md")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func planSnapshot(title, goal string) string {
	return `{"title": "` + title + `", "goal": "` + goal + `", "scope": "Kun cmd", "owner": "ply-planner",
  "selection": "selected", "status": "active", "summary": "Start", "next_action": "Begynn",
  "steps": [{"id": "a", "title": "Prøv i plan", "outcome": "Ferdig", "owner": "ply-planner", "status": "planned", "next_action": "Begynn"}]}`
}

type sharedExperienceOutput struct {
	Knowledge []knowledge.Item `json:"knowledge"`
}

func TestPlanShowGivesSharedExperience(t *testing.T) {
	root := sharedExperienceWorkspace(t, knowledgePointer)
	for _, args := range [][]string{{"--format", "json"}, {}} {
		out, err := runPlan(t, root, planSnapshot("Skjermbilde i 320 px", "Headless Chrome tar skjermbildet"), append([]string{"add", "--project", "a", "--actor", "test", "--reason", "test"}, args...)...)
		if err != nil {
			t.Fatal(err)
		}
		// add and update give the stored revision only.
		if strings.Contains(out, "knowledge") || strings.Contains(out, "Shared experience") {
			t.Fatalf("add shows shared experience:\n%s", out)
		}
	}
	out, err := runPlan(t, root, "", "show", "PLAN-1", "--project", "a", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	var parsed sharedExperienceOutput
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatal(err)
	}
	if ids := knowledgeIDs(parsed.Knowledge); len(ids) != 1 || ids[0] != "k-20260101-headless" || parsed.Knowledge[0].Path == "" {
		t.Fatalf("knowledge %+v", parsed.Knowledge)
	}
	if keys := planKeys(t, out); keys[len(keys)-1] != "knowledge" {
		t.Fatalf("keys %v", keys)
	}
	text, err := runPlan(t, root, "", "show", "PLAN-1", "--project", "a")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "Shared experience (open the file at path):\n  k-20260101-headless  tooling/candidate  a-app\n    Headless Chrome — Du jobber med k-20260101-headless.\n    path: "+parsed.Knowledge[0].Path+"\n") {
		t.Fatalf("text:\n%s", text)
	}
}

func TestPlanShowWithoutSharedExperienceIsUnchanged(t *testing.T) {
	for name, pointer := range map[string]string{"no hits": knowledgePointer, "no source": "", "missing directory": `, "knowledge": {"repo": "a-lib", "path": "absent"}`} {
		t.Run(name, func(t *testing.T) {
			root := sharedExperienceWorkspace(t, pointer)
			if _, err := runPlan(t, root, planSnapshot("Rydd byggkataloger", "Mindre rot"), "add", "--project", "a", "--actor", "test", "--reason", "test"); err != nil {
				t.Fatal(err)
			}
			out, err := runPlan(t, root, "", "show", "PLAN-1", "--project", "a", "--format", "json")
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(out, `"knowledge"`) {
				t.Fatalf("json:\n%s", out)
			}
			text, err := runPlan(t, root, "", "show", "PLAN-1", "--project", "a")
			if err != nil || strings.Contains(text, "Shared experience") {
				t.Fatalf("text %v:\n%s", err, text)
			}
		})
	}
}

func TestAssignmentShowGivesSharedExperience(t *testing.T) {
	root := sharedExperienceWorkspace(t, knowledgePointer)
	for stdin, title := range map[string]string{"Headless Chrome tar skjermbildet.": "Skjermbilde i 320 px", "Mindre rot.": "Rydd byggkataloger"} {
		if _, err := runAssignment(t, root, stdin, "add", "--project", "a", "--title", title, "--repo", "a-app"); err != nil {
			t.Fatal(err)
		}
	}
	ids := map[string]string{}
	for _, id := range []string{"ASSIGNMENT-1", "ASSIGNMENT-2"} {
		out, err := runAssignment(t, root, "", "show", id, "--project", "a")
		if err != nil {
			t.Fatal(err)
		}
		ids[id] = out
	}
	var chromeID, otherID string
	for id, out := range ids {
		if strings.Contains(out, "Skjermbilde") {
			chromeID = id
		} else {
			otherID = id
		}
	}
	out, err := runAssignment(t, root, "", "show", chromeID, "--project", "a", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	var parsed sharedExperienceOutput
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatal(err)
	}
	if got := knowledgeIDs(parsed.Knowledge); len(got) != 1 || got[0] != "k-20260101-headless" {
		t.Fatalf("knowledge %+v", parsed.Knowledge)
	}
	if !strings.Contains(ids[chromeID], "\nShared experience (open the file at path):\n  k-20260101-headless  tooling/candidate  a-app\n") {
		t.Fatalf("text:\n%s", ids[chromeID])
	}
	if strings.Contains(ids[otherID], "Shared experience") {
		t.Fatalf("text:\n%s", ids[otherID])
	}
	out, err = runAssignment(t, root, "", "show", otherID, "--project", "a", "--format", "json")
	if err != nil || strings.Contains(out, `"knowledge"`) {
		t.Fatalf("json %v:\n%s", err, out)
	}
}

// goalStartArchive points Project example of goalStartFixture to an archive
// in its repository. The fixture's Task and Goal say only "T" and "O", so the
// posts avoid both letters except in the one post they point at.
func goalStartArchive(t *testing.T, worktree string) {
	t.Helper()
	root := filepath.Dir(worktree)
	manifest := filepath.Join(root, "ply-workspace.json")
	raw, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	var definitions map[string]any
	if err := json.Unmarshal(raw, &definitions); err != nil {
		t.Fatal(err)
	}
	for _, p := range definitions["projects"].([]any) {
		if project := p.(map[string]any); project["id"] == "example" {
			project["knowledge"] = map[string]string{"repo": "code", "path": "knowledge"}
		}
	}
	if raw, err = json.MarshalIndent(definitions, "", "  "); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifest, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	posts := map[string]string{"k-20260101-x": "T O"}
	for i := 0; i < 9; i++ {
		posts[fmt.Sprintf("k-2026020%d-fyll", i)] = "Fyll"
	}
	for id, title := range posts {
		path := filepath.Join(root, "repo", "knowledge", id+".md")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		text := "---\nid: \"" + id + "\"\ntitle: \"" + title + "\"\nstatus: \"confirmed\"\nscope: {\"project\": \"*\", \"repos\": [], \"roles\": []}\n---\n\nFyll.\n"
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestGoalStartPreviewGivesSharedExperienceOutsideThePrompt(t *testing.T) {
	d, _, worktree, _, _ := goalStartFixture(t)
	before, err := previewGoalStart(t, d)
	if err != nil {
		t.Fatal(err)
	}
	goalStartArchive(t, worktree)
	c := newWorkflowGoalStartCommand(d)
	var out strings.Builder
	c.SetOut(&out)
	c.SetErr(&strings.Builder{})
	c.SilenceErrors, c.SilenceUsage = true, true
	c.SetArgs([]string{"moved", "--spec", "moved-v1", "--revision", "1", "--herdr-workspace", "w1"})
	if err := c.Execute(); err != nil {
		t.Fatal(err)
	}
	var parsed sharedExperienceOutput
	if err := json.Unmarshal([]byte(out.String()), &parsed); err != nil {
		t.Fatal(err)
	}
	if got := knowledgeIDs(parsed.Knowledge); len(got) != 1 || got[0] != "k-20260101-x" {
		t.Fatalf("knowledge %+v in %s", parsed.Knowledge, out.String())
	}
	after, err := previewGoalStart(t, d)
	if err != nil {
		t.Fatal(err)
	}
	// The plan, and so the prompt and the receipt key, are the same with and without the archive.
	if goalStartInput(before, "/r", agentstart.HerdrBinary{}, nil).Prompt != goalStartInput(after, "/r", agentstart.HerdrBinary{}, nil).Prompt || before.OID != after.OID {
		t.Fatalf("plan changed:\n%+v\n%+v", before, after)
	}
}
