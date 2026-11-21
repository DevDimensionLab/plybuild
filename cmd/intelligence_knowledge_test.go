package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/knowledge"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

// knowledgeWorkspace is the assignment workspace whose Project a points to
// pointer and whose repository a-lib holds a small archive under knowledge/.
func knowledgeWorkspace(t *testing.T, pointer string) string {
	t.Helper()
	root := newAssignmentWorkspace(t)
	definitions := strings.Replace(assignmentDefinitions, `{"id": "a-lib", "path": "a/lib"}]}`, `{"id": "a-lib", "path": "a/lib"}]`+pointer+`}`, 1)
	if err := os.WriteFile(filepath.Join(root, workspacedef.File), []byte(definitions), 0o644); err != nil {
		t.Fatal(err)
	}
	posts := map[string]string{
		"tooling/k-20260101-headless.md":   `title: "Headless Chrome"` + "\n" + `status: "candidate"` + "\n" + `scope: {"project": "a", "repos": ["a-app"], "roles": []}` + "\n" + `tags: ["chrome"]`,
		"tooling/sub/k-20260102-tmpdir.md": `title: "TMPDIR og go test"` + "\n" + `status: "confirmed"` + "\n" + `scope: {"project": "*", "repos": [], "roles": []}` + "\n" + `tags: []`,
		"domain/k-20260103-stale.md":       `title: "Gammel chrome-regel"` + "\n" + `status: "stale"` + "\n" + `scope: {"project": "a", "repos": [], "roles": []}` + "\n" + `tags: []`,
		"domain/k-20260104-broken.md":      `tags: ["unclosed"`,
	}
	for file, meta := range posts {
		id := strings.TrimSuffix(filepath.Base(file), ".md")
		text := "---\nid: \"" + id + "\"\n" + meta + "\n---\n\n**Situasjon:** Du jobber med " + id + ".\n"
		path := filepath.Join(root, "a", "lib", "knowledge", filepath.FromSlash(file))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

const knowledgePointer = `, "knowledge": {"repo": "a-lib", "path": "knowledge"}`

type knowledgeOutput struct {
	Kind          string `json:"kind"`
	SchemaVersion int    `json:"schema_version"`
	Project       string `json:"project"`
	Source        struct {
		Repo, Path, Status, Diagnostic string
	} `json:"source"`
	Items   []knowledge.Item        `json:"items"`
	Invalid []knowledge.InvalidFile `json:"invalid"`
	Matched int                     `json:"matched"`
	Hidden  int                     `json:"hidden"`
}

func runKnowledgeJSON(t *testing.T, root string, args ...string) knowledgeOutput {
	t.Helper()
	out, err := runIntelligence(t, root, "", append([]string{"knowledge", "--project", "a", "--format", "json"}, args...)...)
	if err != nil {
		t.Fatalf("knowledge %v: %v", args, err)
	}
	var parsed knowledgeOutput
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("json %s: %v", out, err)
	}
	if strings.Contains(out, `"score"`) {
		t.Fatalf("score in %s", out)
	}
	return parsed
}

func knowledgeIDs(items []knowledge.Item) []string {
	ids := []string{}
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

func TestKnowledgeFindsPostsInTheProjectsArchive(t *testing.T) {
	root := knowledgeWorkspace(t, knowledgePointer)
	out := runKnowledgeJSON(t, root, "chrome", "tmpdir")
	if out.Kind != "PlyIntelligenceKnowledge@1" || out.SchemaVersion != 1 || out.Project != "a" {
		t.Fatalf("frame %+v", out)
	}
	if out.Source.Repo != "a-lib" || out.Source.Path != "knowledge" || out.Source.Status != "ok" || out.Source.Diagnostic != "" {
		t.Fatalf("source %+v", out.Source)
	}
	// The stale post is hidden by the default statuses.
	if got := knowledgeIDs(out.Items); !reflect.DeepEqual(got, []string{"k-20260101-headless", "k-20260102-tmpdir"}) || out.Matched != 2 || out.Hidden != 1 {
		t.Fatalf("items %v matched %d hidden %d", got, out.Matched, out.Hidden)
	}
	physical, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	if item := out.Items[1]; item.Theme != "tooling" || item.Path != filepath.Join(physical, "a", "lib", "knowledge", "tooling", "sub", "k-20260102-tmpdir.md") {
		t.Fatalf("item %+v", item)
	}
	if want := []knowledge.InvalidFile{{File: "domain/k-20260104-broken.md"}}; len(out.Invalid) != 1 || out.Invalid[0].File != want[0].File {
		t.Fatalf("invalid %+v", out.Invalid)
	}

	// --repo ranks a post whose scope fits first; it does not filter.
	if got := knowledgeIDs(runKnowledgeJSON(t, root, "chrome", "tmpdir", "--repo", "b-app").Items); !reflect.DeepEqual(got, []string{"k-20260102-tmpdir", "k-20260101-headless"}) {
		t.Fatalf("--repo b-app %v", got)
	}
	if got := knowledgeIDs(runKnowledgeJSON(t, root, "chrome", "--status", "all").Items); !reflect.DeepEqual(got, []string{"k-20260101-headless", "k-20260103-stale"}) {
		t.Fatalf("--status all %v", got)
	}
	if got := knowledgeIDs(runKnowledgeJSON(t, root, "chrome", "--status", "stale,confirmed").Items); !reflect.DeepEqual(got, []string{"k-20260103-stale"}) {
		t.Fatalf("--status stale,confirmed %v", got)
	}
	if got := knowledgeIDs(runKnowledgeJSON(t, root, "--theme", "tooling", "--tag", "chrome").Items); !reflect.DeepEqual(got, []string{"k-20260101-headless"}) {
		t.Fatalf("--theme --tag %v", got)
	}
	if got := runKnowledgeJSON(t, root, "jobber", "--limit", "1"); !reflect.DeepEqual(knowledgeIDs(got.Items), []string{"k-20260101-headless"}) || got.Matched != 2 {
		t.Fatalf("--limit %+v", got)
	}
	// Without words: the catalog in theme order, without --limit.
	if got := knowledgeIDs(runKnowledgeJSON(t, root, "--limit", "1").Items); !reflect.DeepEqual(got, []string{"k-20260101-headless", "k-20260102-tmpdir"}) {
		t.Fatalf("catalog %v", got)
	}
}

func TestKnowledgeText(t *testing.T) {
	root := knowledgeWorkspace(t, knowledgePointer)
	out, err := runIntelligence(t, root, "", "knowledge", "--project", "a", "chrome")
	if err != nil {
		t.Fatal(err)
	}
	want := "k-20260101-headless  tooling/candidate  a-app\n" +
		"  Headless Chrome — Du jobber med k-20260101-headless.\n" +
		"invalid domain/k-20260104-broken.md: "
	if !strings.HasPrefix(out, want) || !strings.HasSuffix(out, "\n1 more post hidden by --status, --theme, --tag or the Project scope\n") {
		t.Fatalf("text:\n%s", out)
	}
	out, err = runIntelligence(t, root, "", "knowledge", "--project", "a", "--status", "all", "regel")
	if err != nil || !strings.HasPrefix(out, "k-20260103-stale  domain/stale  -\n  Gammel chrome-regel — Du jobber med k-20260103-stale.\n") {
		t.Fatalf("text without repos %q %v", out, err)
	}
	out, err = runIntelligence(t, root, "", "knowledge", "--project", "a", "--status", "all", "nothing")
	if err != nil || !strings.HasPrefix(out, "no matches\n") {
		t.Fatalf("no matches %q %v", out, err)
	}
}

// A missing or bad pointer gives an empty answer with a diagnostic and exit 0,
// and never stops other reads of the workspace.
func TestKnowledgeWithoutAUsableSourceAnswersEmpty(t *testing.T) {
	for _, c := range []struct {
		name, pointer, status, diagnostic string
	}{
		{"not configured", "", "not_configured", "Project a has no knowledge source in ply-workspace.json"},
		{"repository outside the Project", `, "knowledge": {"repo": "b-app", "path": "knowledge"}`, "invalid", "knowledge repository b-app is not a repository of Project a"},
		{"path leaves the repository", `, "knowledge": {"repo": "a-lib", "path": "../app"}`, "invalid", "knowledge path ../app must be a relative path inside repository a-lib"},
		{"absolute path", `, "knowledge": {"repo": "a-lib", "path": "/knowledge"}`, "invalid", "knowledge path /knowledge must be a relative path inside repository a-lib"},
		{"missing directory", `, "knowledge": {"repo": "a-lib", "path": "absent"}`, "missing", "knowledge directory absent does not exist in repository a-lib"},
		{"not a directory", `, "knowledge": {"repo": "a-lib", "path": "knowledge/tooling/k-20260101-headless.md"}`, "missing", "knowledge directory knowledge/tooling/k-20260101-headless.md is not a directory in repository a-lib"},
	} {
		t.Run(c.name, func(t *testing.T) {
			root := knowledgeWorkspace(t, c.pointer)
			out := runKnowledgeJSON(t, root, "chrome")
			if out.Source.Status != c.status || out.Source.Diagnostic != c.diagnostic {
				t.Fatalf("source %+v", out.Source)
			}
			if out.Items == nil || len(out.Items) != 0 || out.Invalid == nil || len(out.Invalid) != 0 || out.Matched != 0 || out.Hidden != 0 {
				t.Fatalf("answer %+v", out)
			}
			text, err := runIntelligence(t, root, "", "knowledge", "--project", "a")
			if err != nil || text != "knowledge source "+c.status+": "+c.diagnostic+"\n" {
				t.Fatalf("text %q %v", text, err)
			}
			if _, err := runIntelligence(t, root, "", "observations", "--project", "a"); err != nil {
				t.Fatalf("observations with this pointer: %v", err)
			}
		})
	}
}

func TestKnowledgeReportsAnUnreadableDirectory(t *testing.T) {
	root := knowledgeWorkspace(t, knowledgePointer)
	dir := filepath.Join(root, "a", "lib", "knowledge")
	if err := os.Chmod(dir, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
	out := runKnowledgeJSON(t, root, "chrome")
	if out.Source.Status != "unreadable" || !strings.HasPrefix(out.Source.Diagnostic, "cannot read knowledge directory knowledge in repository a-lib: ") || len(out.Items) != 0 {
		t.Fatalf("answer %+v", out)
	}
}

func TestKnowledgeRejectsInvalidArguments(t *testing.T) {
	root := knowledgeWorkspace(t, knowledgePointer)
	for _, args := range [][]string{
		{"--limit", "0"},
		{"--format", "xml"},
		{"--status", ""},
		{"--status", "candidate,,confirmed"},
		{"--role", "author"},
	} {
		_, err := runIntelligence(t, root, "", append([]string{"knowledge", "--project", "a"}, args...)...)
		if err == nil || !strings.HasPrefix(err.Error(), "intelligence_invalid_arguments: ") {
			t.Fatalf("%v: %v", args, err)
		}
	}
}

// The pointer's directory is resolved physically: a directory that is, or is
// reached through, a symbolic link out of the repository is invalid, and one
// that is, or is reached through, a link that stays inside is read.
func TestKnowledgePointerIsCheckedThroughSymbolicLinks(t *testing.T) {
	for _, c := range []struct {
		name, path, status, diagnostic string
		items                          []string
	}{
		{"directory is a link out", "out", "invalid", "knowledge path out leaves repository a-lib through a symbolic link", []string{}},
		{"path goes through a link out", "via/archive", "invalid", "knowledge path via/archive leaves repository a-lib through a symbolic link", []string{}},
		{"path goes through a link inside", "shelf/tooling", "ok", "", []string{"k-20260101-headless"}},
		{"directory is a link inside", "alias", "ok", "", []string{"k-20260101-headless"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			root := knowledgeWorkspace(t, `, "knowledge": {"repo": "a-lib", "path": "`+c.path+`"}`)
			outside := t.TempDir()
			post := filepath.Join(outside, "archive", "tooling", "k-20260109-outside.md")
			if err := os.MkdirAll(filepath.Dir(post), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(post, []byte("---\nid: \"k-20260109-outside\"\ntitle: \"chrome\"\n---\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			lib := filepath.Join(root, "a", "lib")
			for link, target := range map[string]string{"out": filepath.Join(outside, "archive"), "via": outside, "shelf": "knowledge", "alias": filepath.Join("knowledge", "tooling")} {
				if err := os.Symlink(target, filepath.Join(lib, link)); err != nil {
					t.Fatal(err)
				}
			}
			out := runKnowledgeJSON(t, root, "chrome")
			if out.Source.Status != c.status || out.Source.Diagnostic != c.diagnostic {
				t.Fatalf("source %+v", out.Source)
			}
			if got := knowledgeIDs(out.Items); !reflect.DeepEqual(got, c.items) || len(out.Invalid) != 0 {
				t.Fatalf("items %v invalid %+v", got, out.Invalid)
			}
		})
	}
}
