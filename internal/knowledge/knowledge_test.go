package knowledge

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const archive = "testdata/archive"

func ids(items []Item) []string {
	out := []string{}
	for _, item := range items {
		out = append(out, item.ID)
	}
	return out
}

func requireIDs(t *testing.T, items []Item, want ...string) {
	t.Helper()
	if got := ids(items); !reflect.DeepEqual(got, want) {
		t.Fatalf("items %v, want %v", got, want)
	}
}

func TestCatalogReadsEveryPostAndReportsUnreadableOnes(t *testing.T) {
	result, err := Search(archive, Query{})
	if err != nil {
		t.Fatal(err)
	}
	// Sorted on theme and id; a subtheme directory keeps its first directory as theme.
	requireIDs(t, result.Items, "k-20260106-future-schema", "k-20260107-no-situation", "k-20260108-bare", "k-20260101-headless-chrome", "k-20260105-map-order")
	want := Item{ID: "k-20260105-map-order", Theme: "verification", Type: "semantic", Status: "confirmed",
		Scope: Scope{Project: "*", Repos: []string{}, Roles: []string{}}, Tags: []string{"go"},
		Title: "Map-rekkefølge gir et sett av svar", Situation: "Du tester kode som går gjennom et map.",
		Path: filepath.Join(archive, "verification", "go", "k-20260105-map-order.md")}
	if !reflect.DeepEqual(result.Items[4], want) {
		t.Fatalf("item %+v\nwant %+v", result.Items[4], want)
	}
	// Unknown schema and missing fields are shown with empty values.
	if item := result.Items[0]; item.Title != "Fremtidig format" || item.Situation != "Arkivet har byttet format." {
		t.Fatalf("future schema %+v", item)
	}
	if item := result.Items[1]; item.Situation != "" || item.Status != "stale" {
		t.Fatalf("no situation line %+v", item)
	}
	bare := Item{ID: "k-20260108-bare", Theme: "domain", Scope: Scope{Repos: []string{}, Roles: []string{}}, Tags: []string{},
		Situation: "Bare id er oppgitt.", Path: filepath.Join(archive, "domain", "k-20260108-bare.md")}
	if !reflect.DeepEqual(result.Items[2], bare) {
		t.Fatalf("bare %+v", result.Items[2])
	}
	var files []string
	for _, invalid := range result.Invalid {
		if invalid.Error == "" {
			t.Fatalf("invalid without error: %+v", invalid)
		}
		files = append(files, invalid.File)
	}
	if want := []string{"broken/k-20260109-broken-yaml.md", "broken/k-20260110-no-id.md", "broken/k-20260111-no-frontmatter.md"}; !reflect.DeepEqual(files, want) {
		t.Fatalf("invalid %v (%+v)", files, result.Invalid)
	}
	if result.Matched != 5 || result.Hidden != 0 {
		t.Fatalf("matched %d hidden %d", result.Matched, result.Hidden)
	}
}

func TestSearchReportsAMissingDirectory(t *testing.T) {
	_, err := Search(filepath.Join(t.TempDir(), "absent"), Query{})
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("error %v", err)
	}
}

// fixture is a small archive entry for ranking and filter tests.
type fixture struct {
	path, id, title, status, project, situation, body string
	repos, tags                                       []string
}

func quoted(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	return `["` + strings.Join(values, `", "`) + `"]`
}

func writeArchive(t *testing.T, posts ...fixture) string {
	t.Helper()
	dir := t.TempDir()
	for _, p := range posts {
		project := p.project
		if project == "" {
			project = "ply"
		}
		status := p.status
		if status == "" {
			status = "candidate"
		}
		text := "---\n" +
			`schema: "ply.knowledge/v1"` + "\n" +
			`id: "` + p.id + `"` + "\n" +
			`title: "` + p.title + `"` + "\n" +
			`status: "` + status + `"` + "\n" +
			`scope: {"project": "` + project + `", "repos": ` + quoted(p.repos) + `, "roles": []}` + "\n" +
			`tags: ` + quoted(p.tags) + "\n" +
			"---\n\n**Situasjon:** " + p.situation + "\n\n" + p.body + "\n"
		file := filepath.Join(dir, filepath.FromSlash(p.path), p.id+".md")
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestSearchRanksOnDistinctWordsRepoFieldAndID(t *testing.T) {
	dir := writeArchive(t,
		fixture{path: "a", id: "k-1-body", body: "skjermbilde chrome"},
		fixture{path: "a", id: "k-2-title", title: "Chrome"},
		fixture{path: "a", id: "k-3-situation", situation: "chrome"},
		fixture{path: "b", id: "k-4-tag", tags: []string{"chrome"}},
		fixture{path: "b", id: "k-0-id-chrome"},
		fixture{path: "b", id: "k-5-repo-body", body: "chrome", repos: []string{"ply"}},
		fixture{path: "b", id: "k-6-other-repo", title: "chrome", repos: []string{"trip"}},
		fixture{path: "b", id: "k-7-none", title: "uten treff"},
	)
	result, err := Search(dir, Query{Words: []string{"chrome", "skjermbilde", "chrome"}})
	if err != nil {
		t.Fatal(err)
	}
	// Two distinct words first; then id, tags and title before the situation line before the body, then id.
	requireIDs(t, result.Items, "k-1-body", "k-0-id-chrome", "k-2-title", "k-4-tag", "k-6-other-repo", "k-3-situation", "k-5-repo-body")
	if result.Matched != 7 || result.Hidden != 0 {
		t.Fatalf("matched %d hidden %d", result.Matched, result.Hidden)
	}
	// A post whose scope fits --repo comes first among equal word counts, before field weight;
	// an empty repos list fits every repository.
	result, err = Search(dir, Query{Words: []string{"chrome"}, Repo: "ply"})
	if err != nil {
		t.Fatal(err)
	}
	requireIDs(t, result.Items, "k-0-id-chrome", "k-2-title", "k-4-tag", "k-3-situation", "k-1-body", "k-5-repo-body", "k-6-other-repo")
}

func TestWordsMatchSubstringsOfLowercaseTokens(t *testing.T) {
	dir := writeArchive(t,
		fixture{path: "a", id: "k-1-letters", title: "Ærfugl på Østlandet"},
		fixture{path: "a", id: "k-2-split", body: "kjør go-test/TMPDIR"},
		fixture{path: "a", id: "k-3-across", body: "go test"},
	)
	for _, c := range []struct {
		words []string
		want  []string
	}{
		{[]string{"ærf", "ØSTL"}, []string{"k-1-letters"}},
		{[]string{"tmpdir"}, []string{"k-2-split"}},
		{[]string{"go-test"}, []string{"k-2-split", "k-3-across"}},
		// A word never spans two tokens.
		{[]string{"gotest"}, []string{}},
	} {
		result, err := Search(dir, Query{Words: c.words})
		if err != nil {
			t.Fatal(err)
		}
		requireIDs(t, result.Items, c.want...)
	}
}

func TestFiltersHideMatchesAndSayHowMany(t *testing.T) {
	dir := writeArchive(t,
		fixture{path: "tooling", id: "k-1-candidate", title: "chrome", tags: []string{"chrome"}},
		fixture{path: "tooling", id: "k-2-confirmed", title: "chrome", status: "confirmed"},
		fixture{path: "tooling", id: "k-3-stale", title: "chrome", status: "stale"},
		fixture{path: "domain/sub", id: "k-4-superseded", title: "chrome", status: "superseded"},
		fixture{path: "domain", id: "k-5-other-project", title: "chrome", project: "trip"},
		fixture{path: "domain", id: "k-6-everyone", title: "chrome", project: "*"},
		fixture{path: "domain", id: "k-7-no-hit", title: "annet"},
	)
	for _, c := range []struct {
		name   string
		query  Query
		want   []string
		hidden int
	}{
		{"default statuses", Query{Statuses: DefaultStatuses}, []string{"k-1-candidate", "k-2-confirmed", "k-5-other-project", "k-6-everyone"}, 2},
		{"all statuses", Query{}, []string{"k-1-candidate", "k-2-confirmed", "k-3-stale", "k-4-superseded", "k-5-other-project", "k-6-everyone"}, 0},
		{"one status", Query{Statuses: []string{"stale"}}, []string{"k-3-stale"}, 5},
		{"project or star", Query{Project: "ply"}, []string{"k-1-candidate", "k-2-confirmed", "k-3-stale", "k-4-superseded", "k-6-everyone"}, 1},
		{"theme is the first directory", Query{Theme: "domain"}, []string{"k-4-superseded", "k-5-other-project", "k-6-everyone"}, 3},
		{"exact tag", Query{Tag: "chrome"}, []string{"k-1-candidate"}, 5},
		{"no tag prefix", Query{Tag: "chr"}, []string{}, 6},
	} {
		t.Run(c.name, func(t *testing.T) {
			c.query.Words = []string{"chrome"}
			result, err := Search(dir, c.query)
			if err != nil {
				t.Fatal(err)
			}
			requireIDs(t, result.Items, c.want...)
			if result.Matched != len(c.want) || result.Hidden != c.hidden {
				t.Fatalf("matched %d hidden %d, want %d %d", result.Matched, result.Hidden, len(c.want), c.hidden)
			}
		})
	}
}

func TestMissingValuesDoNotFilterAPostOut(t *testing.T) {
	result, err := Search(archive, Query{Words: []string{"bare"}, Project: "ply", Statuses: DefaultStatuses})
	if err != nil {
		t.Fatal(err)
	}
	requireIDs(t, result.Items, "k-20260108-bare")
}

func TestLimitCutsMatchesButNotTheCatalog(t *testing.T) {
	dir := writeArchive(t,
		fixture{path: "b", id: "k-1", title: "chrome"},
		fixture{path: "a", id: "k-2", title: "chrome"},
		fixture{path: "a", id: "k-3", title: "chrome"},
	)
	result, err := Search(dir, Query{Words: []string{"chrome"}, Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	requireIDs(t, result.Items, "k-1", "k-2")
	if result.Matched != 3 {
		t.Fatalf("matched %d", result.Matched)
	}
	result, err = Search(dir, Query{Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	requireIDs(t, result.Items, "k-2", "k-3", "k-1")
}

func TestARepeatedWordCountsOnce(t *testing.T) {
	dir := writeArchive(t,
		fixture{path: "a", id: "k-1-body", body: "chrome"},
		fixture{path: "a", id: "k-2-title", title: "tmpdir"},
	)
	result, err := Search(dir, Query{Words: []string{"chrome", "Chrome", "tmpdir"}})
	if err != nil {
		t.Fatal(err)
	}
	requireIDs(t, result.Items, "k-2-title", "k-1-body")
}

func symlink(t *testing.T, target, link string) {
	t.Helper()
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
}

// A k-*.md that is a symbolic link or a directory is not opened: it is listed
// as invalid, and the other posts are still shown.
func TestOnlyRegularFilesAreReadAsPosts(t *testing.T) {
	outside := writeArchive(t, fixture{path: "x", id: "k-9-outside", title: "chrome"})
	dir := writeArchive(t, fixture{path: "a", id: "k-1-kept", title: "chrome"})
	symlink(t, filepath.Join(outside, "x", "k-9-outside.md"), filepath.Join(dir, "a", "k-2-link.md"))
	if err := os.Mkdir(filepath.Join(dir, "a", "k-3-dir.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	result, err := Search(dir, Query{Words: []string{"chrome"}})
	if err != nil {
		t.Fatal(err)
	}
	requireIDs(t, result.Items, "k-1-kept")
	want := []InvalidFile{{File: "a/k-2-link.md", Error: "not a regular file"}, {File: "a/k-3-dir.md", Error: "not a regular file"}}
	if !reflect.DeepEqual(result.Invalid, want) {
		t.Fatalf("invalid %+v", result.Invalid)
	}
}

// A directory under the archive that is a symbolic link is not followed.
func TestSymlinkedDirectoriesAreNotFollowed(t *testing.T) {
	outside := writeArchive(t, fixture{path: "x", id: "k-9-outside", title: "chrome"})
	dir := writeArchive(t, fixture{path: "a", id: "k-1-kept", title: "chrome"})
	symlink(t, filepath.Join(outside, "x"), filepath.Join(dir, "a", "linked"))
	symlink(t, outside, filepath.Join(dir, "linked-theme"))
	result, err := Search(dir, Query{Words: []string{"chrome"}})
	if err != nil {
		t.Fatal(err)
	}
	requireIDs(t, result.Items, "k-1-kept")
	if len(result.Invalid) != 0 {
		t.Fatalf("invalid %+v", result.Invalid)
	}
}
