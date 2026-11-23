package cmd

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/publication"
)

func moduleRuns(items []publication.View) map[string]string {
	out := map[string]string{}
	for _, v := range items {
		out[v.Repo+":"+v.Module+":"+v.Category] = v.Run
	}
	return out
}

// A module publication stands beside the repo's own in the same category; only
// --module and --modules read it, so the default list is unchanged.
func TestPublicationModuleStandsBesideTheRepo(t *testing.T) {
	root := newAssignmentWorkspace(t)
	dir := producerDir(t, root, "map", map[string]string{"index.html": "x"})
	args := []string{"--repo", "a-app", "--category", "code-health", "--dir", dir, "--entry", "index.html", "--title", "t"}
	rootRun := publishJSON(t, root, args...)
	module := publishJSON(t, root, append(args, "--module", "dom/app")...)
	other := publishJSON(t, root, "--repo", "a-lib", "--module", "core", "--category", "code-health", "--dir", dir, "--entry", "index.html", "--title", "t")
	publishJSON(t, root, "--repo", "a-lib", "--module", "core", "--category", "code-map", "--dir", dir, "--entry", "index.html", "--title", "t")

	publications := filepath.Join(root, "a", "data", "publications")
	if want := filepath.Join(publications, "a-app", "_modules", "dom%2Fapp", "code-health", module.Run); module.Path != want {
		t.Fatalf("module path %q, want %q", module.Path, want)
	}
	if rootRun.Module != "." || module.Module != "dom/app" {
		t.Fatalf("module fields %q %q", rootRun.Module, module.Module)
	}
	stored, err := os.ReadFile(filepath.Join(rootRun.Path, "publication.json"))
	if err != nil || strings.Contains(string(stored), `"module"`) {
		t.Fatalf("root publication.json carries module: %s %v", stored, err)
	}
	stored, err = os.ReadFile(filepath.Join(module.Path, "publication.json"))
	if err != nil || !strings.Contains(string(stored), `"module": "dom/app"`) {
		t.Fatalf("module publication.json: %s %v", stored, err)
	}

	if got := moduleRuns(listJSON(t, root, "--project", "a", "--repo", "a-app").Items); len(got) != 1 || got["a-app:.:code-health"] != rootRun.Run {
		t.Fatalf("default list %+v", got)
	}
	if got := moduleRuns(listJSON(t, root, "--project", "a").Items); len(got) != 1 {
		t.Fatalf("default list without --repo %+v", got)
	}
	if got := moduleRuns(listJSON(t, root, "--project", "a", "--repo", "a-app", "--module", "dom/app").Items); len(got) != 1 || got["a-app:dom/app:code-health"] != module.Run {
		t.Fatalf("--module %+v", got)
	}
	if got := moduleRuns(listJSON(t, root, "--project", "a", "--repo", "a-app", "--module", "dom").Items); len(got) != 0 {
		t.Fatalf("--module gives descendants %+v", got)
	}
	if got := moduleRuns(listJSON(t, root, "--project", "a", "--repo", "a-app", "--modules").Items); len(got) != 2 || got["a-app:.:code-health"] != rootRun.Run || got["a-app:dom/app:code-health"] != module.Run {
		t.Fatalf("--modules %+v", got)
	}
	if got := moduleRuns(listJSON(t, root, "--project", "a", "--modules", "--category", "code-health").Items); len(got) != 3 || got["a-lib:core:code-health"] != other.Run {
		t.Fatalf("--modules across repos %+v", got)
	}
	if got := moduleRuns(listJSON(t, root, "--project", "a", "--repo", "a-app", "--module", ".").Items); len(got) != 1 || got["a-app:.:code-health"] != rootRun.Run {
		t.Fatalf("--module . %+v", got)
	}

	out, err := runPublication(t, root, "show", "--project", "a", "--repo", "a-app", "--module", "dom/app", "--category", "code-health", "--format", "json")
	var shown publication.View
	if err != nil || json.Unmarshal([]byte(out), &shown) != nil || shown.Run != module.Run || shown.Module != "dom/app" || shown.Path != module.Path {
		t.Fatalf("show %q %v", out, err)
	}
	if text, err := runPublication(t, root, "show", "--project", "a", "--repo", "a-app", "--module", "dom/app", "--category", "code-health"); err != nil || !strings.Contains(text, "a-app/dom/app/code-health") {
		t.Fatalf("show text %q %v", text, err)
	}
}

func TestPublicationWithdrawModuleHidesOnlyTheModule(t *testing.T) {
	root := newAssignmentWorkspace(t)
	dir := producerDir(t, root, "map", map[string]string{"index.html": "x"})
	args := []string{"--repo", "a-app", "--category", "code-health", "--dir", dir, "--entry", "index.html", "--title", "t"}
	rootRun := publishJSON(t, root, args...)
	publishJSON(t, root, append(args, "--module", "dom/app")...)
	w := withdrawJSON(t, root, "--repo", "a-app", "--module", "dom/app", "--category", "code-health")
	if w.Module != "dom/app" {
		t.Fatalf("withdrawal %+v", w)
	}
	if got := moduleRuns(listJSON(t, root, "--project", "a", "--repo", "a-app", "--modules").Items); len(got) != 1 || got["a-app:.:code-health"] != rootRun.Run {
		t.Fatalf("after withdraw %+v", got)
	}
	if w := withdrawJSON(t, root, "--repo", "a-app", "--category", "code-health"); w.Module != "." {
		t.Fatalf("root withdrawal %+v", w)
	}
}

func TestPublicationRetentionCountsPerModule(t *testing.T) {
	root := newAssignmentWorkspace(t)
	dir := producerDir(t, root, "map", map[string]string{"index.html": "x"})
	args := []string{"--repo", "a-app", "--category", "code-health", "--dir", dir, "--entry", "index.html", "--title", "t"}
	for i := 0; i < 3; i++ {
		publishJSON(t, root, args...)
	}
	for i := 0; i < 6; i++ {
		publishJSON(t, root, append(args, "--module", "dom")...)
	}
	all := listJSON(t, root, "--project", "a", "--repo", "a-app", "--modules", "--all").Items
	count := map[string]int{}
	for _, v := range all {
		count[v.Module]++
	}
	if count["."] != 3 || count["dom"] != 5 {
		t.Fatalf("kept %v", count)
	}
}

// Ply checks the form of a module path, not that the module is in the tree.
func TestPublicationModuleForm(t *testing.T) {
	root := newAssignmentWorkspace(t)
	dir := producerDir(t, root, "map", map[string]string{"index.html": "x"})
	publish := func(module string, extra ...string) error {
		args := append([]string{"publish", "--project", "a", "--repo", "a-app", "--module", module, "--category", "c", "--dir", dir, "--entry", "index.html", "--title", "t"}, extra...)
		_, err := runPublication(t, root, args...)
		return err
	}
	exact := strings.Repeat("abcdefgh/", 20)
	for len(url.PathEscape(exact)) < 255 {
		exact += "x"
	}
	if n := len(url.PathEscape(exact)); n != 255 {
		t.Fatalf("fixture escapes to %d bytes", n)
	}
	if err := publish(exact); err != nil {
		t.Fatalf("255 bytes encoded: %v", err)
	}
	if err := publish("no/such/module"); err != nil {
		t.Fatalf("a module outside the tree: %v", err)
	}
	for _, bad := range []string{"../x", "/x", "a//b", "a/./b", ".x/y", "a/.x", "a/", "a\\b", "a\x01b", "..", "", exact + "y"} {
		requirePublicationClass(t, publish(bad), publication.ClassInvalidArguments)
	}
	_, err := runPublication(t, root, "publish", "--project", "a", "--module", "dom", "--category", "c", "--dir", dir, "--entry", "index.html", "--title", "t")
	requirePublicationClass(t, err, publication.ClassInvalidArguments)
	for _, args := range [][]string{
		{"list", "--project", "a", "--module", "dom"},
		{"show", "--project", "a", "--module", "dom", "--category", "c"},
		{"withdraw", "--project", "a", "--module", "dom", "--category", "c"},
		{"list", "--project", "a", "--repo", "a-app", "--module", "dom", "--modules"},
		{"list", "--project", "a", "--repo", "a-app", "--module", "../x"},
	} {
		_, err := runPublication(t, root, args...)
		requirePublicationClass(t, err, publication.ClassInvalidArguments)
	}
}

// A directory under _modules that is not canonically encoded, or whose
// publication.json names another module, is listed as invalid.
func TestPublicationListReportsMismatchedModuleDirectories(t *testing.T) {
	root := newAssignmentWorkspace(t)
	dir := producerDir(t, root, "map", map[string]string{"index.html": "x"})
	v := publishJSON(t, root, "--repo", "a-app", "--module", "dom/app", "--category", "c", "--dir", dir, "--entry", "index.html", "--title", "t")
	modules := filepath.Join(root, "a", "data", "publications", "a-app", "_modules")
	if err := os.MkdirAll(filepath.Join(modules, "other"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Dir(v.Path), filepath.Join(modules, "other", "c")); err != nil {
		t.Fatal(err)
	}
	// %70 is p: it decodes to dom/app but is not how Ply encodes it.
	if err := os.MkdirAll(filepath.Join(modules, "dom%2Fap%70", "c"), 0o755); err != nil {
		t.Fatal(err)
	}
	list := listJSON(t, root, "--project", "a", "--repo", "a-app", "--modules")
	if len(list.Items) != 0 || len(list.Invalid) != 2 {
		t.Fatalf("list %+v", list)
	}
}
