package assignment

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/filelock"
)

var fixedNow = time.Date(2026, 11, 12, 9, 14, 3, 0, time.FixedZone("CET", 3600))

func testStore(t *testing.T) (Store, string) {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "assignments")
	store := Open(dir)
	store.now = func() time.Time { return fixedNow }
	return store, dir
}

func mustAdd(t *testing.T, s Store, title string) Assignment {
	t.Helper()
	a, err := s.Add(title, "text\n", []string{"ply"}, nil)
	if err != nil {
		t.Fatalf("add %q: %v", title, err)
	}
	return a
}

func writeRaw(t *testing.T, dir, name, contents string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func validFile(id string) string {
	return fmt.Sprintf(`{"id":%q,"title":"T","text":"x","created_at":"2026-11-12T09:14:03Z","repos":["ply"]}`, id)
}

func requireClass(t *testing.T, err error, class string) *Error {
	t.Helper()
	var e *Error
	if !errors.As(err, &e) || e.Class != class {
		t.Fatalf("error %v, want class %s", err, class)
	}
	if e.Error() != e.Class+": "+e.Detail {
		t.Fatalf("Error() = %q", e.Error())
	}
	return e
}

func requireNoPath(t *testing.T, s Store, err error, more ...string) {
	t.Helper()
	for _, forbidden := range append([]string{s.dir, filepath.Dir(s.dir)}, more...) {
		if strings.Contains(err.Error(), forbidden) {
			t.Fatalf("error %q contains %q", err, forbidden)
		}
	}
}

func requireAbsent(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("%s exists: %v", path, err)
	}
}

func TestAddFirstInMissingDirectoryRoundTrips(t *testing.T) {
	s, dir := testStore(t)
	text := "Vis repo-tags.\n<b>&</b>\r\nBakgrunn: Task x.\n"
	added, err := s.Add("  Vis repo-tags  ", text, []string{" ply ", "app"}, []string{"backlog", "Web app"})
	if err != nil {
		t.Fatal(err)
	}
	want := Assignment{ID: "ASSIGNMENT-1", Title: "Vis repo-tags", Text: text, CreatedAt: "2026-11-12T08:14:03Z", Repos: []string{"ply", "app"}, Topics: []string{"backlog", "Web app"}}
	if !reflect.DeepEqual(added, want) {
		t.Fatalf("add %+v", added)
	}
	read, err := s.Read("ASSIGNMENT-1")
	if err != nil || !reflect.DeepEqual(read, want) {
		t.Fatalf("read %+v %v", read, err)
	}
	info, err := os.Stat(filepath.Join(dir, "ASSIGNMENT-1.json"))
	if err != nil || info.Mode().Perm() != 0o644 {
		t.Fatalf("file mode %v %v", info, err)
	}
	contents, _ := os.ReadFile(filepath.Join(dir, "ASSIGNMENT-1.json"))
	if !strings.Contains(string(contents), "<b>&</b>") || !strings.HasSuffix(string(contents), "}\n") {
		t.Fatalf("contents %s", contents)
	}
	noTopics, err := s.Add("T", "x", []string{"ply"}, nil)
	if err != nil || noTopics.Topics != nil {
		t.Fatalf("no topics %+v %v", noTopics, err)
	}
	if contents, _ := os.ReadFile(filepath.Join(dir, "ASSIGNMENT-2.json")); strings.Contains(string(contents), "topics") {
		t.Fatalf("empty topics written: %s", contents)
	}
}

func TestAddNumbersFromHighestFileName(t *testing.T) {
	for name, tc := range map[string]struct {
		setup func(t *testing.T, dir string)
		want  string
	}{
		"gap keeps max": {func(t *testing.T, dir string) {
			writeRaw(t, dir, "ASSIGNMENT-1.json", validFile("ASSIGNMENT-1"))
			writeRaw(t, dir, "ASSIGNMENT-3.json", validFile("ASSIGNMENT-3"))
		}, "ASSIGNMENT-4"},
		"corrupt file keeps its number": {func(t *testing.T, dir string) {
			writeRaw(t, dir, "ASSIGNMENT-5.json", "{broken")
		}, "ASSIGNMENT-6"},
		"directory keeps its number": {func(t *testing.T, dir string) {
			if err := os.MkdirAll(filepath.Join(dir, "ASSIGNMENT-9.json"), 0o755); err != nil {
				t.Fatal(err)
			}
		}, "ASSIGNMENT-10"},
		"non-canonical names do not count": {func(t *testing.T, dir string) {
			writeRaw(t, dir, "ASSIGNMENT-01.json", validFile("ASSIGNMENT-1"))
			writeRaw(t, dir, "x.json", "{}")
			writeRaw(t, dir, ".atomicfile-1.tmp", "")
			writeRaw(t, dir, "ASSIGNMENT-7.json.bak", "")
		}, "ASSIGNMENT-1"},
	} {
		t.Run(name, func(t *testing.T) {
			s, dir := testStore(t)
			tc.setup(t, dir)
			if got := mustAdd(t, s, "T"); got.ID != tc.want {
				t.Fatalf("id %s, want %s", got.ID, tc.want)
			}
		})
	}
}

func TestAddExhaustedIDSpace(t *testing.T) {
	s, dir := testStore(t)
	writeRaw(t, dir, "ASSIGNMENT-999999999.json", validFile("ASSIGNMENT-999999999"))
	_, err := s.Add("T", "x", []string{"ply"}, nil)
	if e := requireClass(t, err, ClassIOError); e.Detail != "id space exhausted" {
		t.Fatalf("detail %q", e.Detail)
	}
}

func TestConcurrentAddsGetDistinctIDs(t *testing.T) {
	_, dir := testStore(t)
	const writers = 16
	ids := make(chan string, writers)
	var wg sync.WaitGroup
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			a, err := Open(dir).Add(fmt.Sprintf("T%d", i), "x", []string{"ply"}, nil)
			if err != nil {
				t.Error(err)
				return
			}
			ids <- a.ID
		}(i)
	}
	wg.Wait()
	close(ids)
	seen := map[string]bool{}
	for id := range ids {
		seen[id] = true
	}
	for n := 1; n <= writers; n++ {
		if !seen[fmt.Sprintf("ASSIGNMENT-%d", n)] {
			t.Fatalf("missing ASSIGNMENT-%d in %v", n, seen)
		}
	}
	list, err := Open(dir).List()
	if err != nil || len(list.Items) != writers || len(list.Invalid) != 0 {
		t.Fatalf("list %+v %v", list, err)
	}
}

func TestFailedWriteLeavesNoFileAndReusesNumber(t *testing.T) {
	s, dir := testStore(t)
	mustAdd(t, s, "first")
	failing := s
	failing.replace = func(string, []byte, fs.FileMode) error {
		return &fs.PathError{Op: "rename", Path: filepath.Join(dir, "x"), Err: fs.ErrPermission}
	}
	_, err := failing.Add("second", "secret text", []string{"ply"}, nil)
	requireClass(t, err, ClassIOError)
	requireNoPath(t, s, err, "secret text")
	requireAbsent(t, filepath.Join(dir, "ASSIGNMENT-2.json"))
	if got := mustAdd(t, s, "third"); got.ID != "ASSIGNMENT-2" {
		t.Fatalf("id %s", got.ID)
	}
}

func TestLockFailureIsIOErrorWithoutFile(t *testing.T) {
	s, dir := testStore(t)
	s.lock = filelock.Options{Fail: func(op string) error {
		if op == "lock" {
			return errors.New("injected")
		}
		return nil
	}}
	_, err := s.Add("T", "x", []string{"ply"}, nil)
	if e := requireClass(t, err, ClassIOError); e.Detail != "lock assignments" {
		t.Fatalf("detail %q", e.Detail)
	}
	requireAbsent(t, filepath.Join(dir, "ASSIGNMENT-1.json"))
}

func TestAddRejectsInvalidArgumentsWithoutEffect(t *testing.T) {
	long := strings.Repeat("ø", 201)
	repos := func(n int) []string {
		out := []string{}
		for i := 0; i < n; i++ {
			out = append(out, fmt.Sprintf("r%d", i))
		}
		return out
	}
	topics := func(n int) []string {
		out := []string{}
		for i := 0; i < n; i++ {
			out = append(out, fmt.Sprintf("t%d", i))
		}
		return out
	}
	for name, tc := range map[string]struct {
		title, text   string
		repos, topics []string
	}{
		"blank title":          {"  ", "x", []string{"ply"}, nil},
		"two-line title":       {"a\nb", "x", []string{"ply"}, nil},
		"title 201 runes":      {long, "x", []string{"ply"}, nil},
		"blank text":           {"T", " \n\t", []string{"ply"}, nil},
		"text over 64 KiB":     {"T", strings.Repeat("a", 64*1024+1), []string{"ply"}, nil},
		"invalid UTF-8 text":   {"T", "\xff", []string{"ply"}, nil},
		"no repos":             {"T", "x", nil, nil},
		"33 repos":             {"T", "x", repos(33), nil},
		"invalid repo":         {"T", "x", []string{"Ply"}, nil},
		"duplicate repo":       {"T", "x", []string{"ply", " ply"}, nil},
		"duplicate topic":      {"T", "x", []string{"ply"}, []string{"Web", "web"}},
		"33 topics":            {"T", "x", []string{"ply"}, topics(33)},
		"control in topic":     {"T", "x", []string{"ply"}, []string{"a\tb"}},
		"blank topic":          {"T", "x", []string{"ply"}, []string{" "}},
		"topic over 100 runes": {"T", "x", []string{"ply"}, []string{strings.Repeat("t", 101)}},
	} {
		t.Run(name, func(t *testing.T) {
			s, dir := testStore(t)
			_, err := s.Add(tc.title, tc.text, tc.repos, tc.topics)
			requireClass(t, err, ClassInvalidArguments)
			requireAbsent(t, dir)
		})
	}
	s, _ := testStore(t)
	if _, err := s.Add(strings.Repeat("ø", 200), strings.Repeat("a", 64*1024), repos(32), topics(32)); err != nil {
		t.Fatalf("limits: %v", err)
	}
	relative := Open("assignments")
	_, err := relative.Add("T", "x", []string{"ply"}, nil)
	requireClass(t, err, ClassInvalidArguments)
	requireAbsent(t, "assignments")
}

func TestReadNotFoundAndInvalidID(t *testing.T) {
	s, dir := testStore(t)
	_, err := s.Read("ASSIGNMENT-1")
	if e := requireClass(t, err, ClassNotFound); e.Detail != "ASSIGNMENT-1" {
		t.Fatalf("detail %q", e.Detail)
	}
	requireAbsent(t, dir)
	mustAdd(t, s, "T")
	_, err = s.Read("ASSIGNMENT-2")
	requireClass(t, err, ClassNotFound)
	for _, id := range []string{"ASSIGNMENT-01", "ASSIGNMENT-0", "assignment-1", "ASSIGNMENT-1.json", "../ASSIGNMENT-1", "ASSIGNMENT-1234567890"} {
		_, err := s.Read(id)
		if e := requireClass(t, err, ClassInvalidArguments); e.Detail != "invalid assignment id" {
			t.Fatalf("%s: detail %q", id, e.Detail)
		}
	}
	if _, err := os.Lstat(filepath.Join(dir, ".lock")); err != nil {
		t.Fatalf("lock from add: %v", err)
	}
}

func TestReadRejectsInvalidFiles(t *testing.T) {
	const secret = "hemmelig-innhold"
	valid := func(edit func(string) string) string { return edit(validFile("ASSIGNMENT-1")) }
	replace := func(old, new string) func(string) string {
		return func(s string) string { return strings.Replace(s, old, new, 1) }
	}
	for name, tc := range map[string]struct {
		contents string
		detail   string
	}{
		"unknown field":    {valid(replace(`"repos"`, `"`+secret+`":1,"repos"`)), "unknown field"},
		"id differs":       {validFile("ASSIGNMENT-2"), "id does not match file name"},
		"trailing data":    {validFile("ASSIGNMENT-1") + "{}", "data after JSON value"},
		"syntax":           {`{"id":` + secret, "invalid JSON"},
		"wrong type":       {valid(replace(`"title":"T"`, `"title":7`)), "invalid JSON"},
		"missing title":    {valid(replace(`"title":"T",`, ``)), "missing title"},
		"missing text":     {valid(replace(`"text":"x",`, ``)), "missing text"},
		"missing created":  {valid(replace(`"created_at":"2026-11-12T09:14:03Z",`, ``)), "missing created_at"},
		"missing repos":    {valid(replace(`,"repos":["ply"]`, ``)), "missing repos"},
		"case-folded key":  {valid(replace(`"text":"x"`, `"TEXT":"x"`)), "unknown field"},
		"duplicate key":    {valid(replace(`"text":"x"`, `"text":"x","text":"y"`)), "duplicate field"},
		"local time":       {valid(replace(`09:14:03Z`, `10:14:03+01:00`)), "invalid created_at"},
		"untrimmed title":  {valid(replace(`"title":"T"`, `"title":" T"`)), "invalid title"},
		"blank text":       {valid(replace(`"text":"x"`, `"text":" "`)), "invalid text"},
		"duplicate repos":  {valid(replace(`["ply"]`, `["ply","ply"]`)), "invalid repos"},
		"empty repos":      {valid(replace(`["ply"]`, `[]`)), "invalid repos"},
		"invalid topic":    {valid(replace(`["ply"]`, `["ply"],"topics":[" x"]`)), "invalid topics"},
		"duplicate topics": {valid(replace(`["ply"]`, `["ply"],"topics":["A","a"]`)), "invalid topics"},
	} {
		t.Run(name, func(t *testing.T) {
			s, dir := testStore(t)
			writeRaw(t, dir, "ASSIGNMENT-1.json", tc.contents)
			_, err := s.Read("ASSIGNMENT-1")
			if e := requireClass(t, err, ClassInvalid); e.Detail != tc.detail {
				t.Fatalf("detail %q, want %q", e.Detail, tc.detail)
			}
			requireNoPath(t, s, err, secret)
		})
	}
	t.Run("symlink", func(t *testing.T) {
		s, dir := testStore(t)
		writeRaw(t, dir, "target", validFile("ASSIGNMENT-1"))
		if err := os.Symlink(filepath.Join(dir, "target"), filepath.Join(dir, "ASSIGNMENT-1.json")); err != nil {
			t.Fatal(err)
		}
		_, err := s.Read("ASSIGNMENT-1")
		if e := requireClass(t, err, ClassInvalid); e.Detail != "not a regular file" {
			t.Fatalf("detail %q", e.Detail)
		}
		requireNoPath(t, s, err)
	})
	t.Run("too large", func(t *testing.T) {
		s, dir := testStore(t)
		big := strings.Replace(validFile("ASSIGNMENT-1"), `"x"`, `"`+strings.Repeat("a", 1<<20)+`"`, 1)
		writeRaw(t, dir, "ASSIGNMENT-1.json", big)
		_, err := s.Read("ASSIGNMENT-1")
		if e := requireClass(t, err, ClassInvalid); e.Detail != "file too large" {
			t.Fatalf("detail %q", e.Detail)
		}
	})
}

func TestListOrdersNumericallyAndReportsInvalidFiles(t *testing.T) {
	s, dir := testStore(t)
	for _, n := range []int{11, 2, 10} {
		id := fmt.Sprintf("ASSIGNMENT-%d", n)
		writeRaw(t, dir, id+".json", validFile(id))
	}
	writeRaw(t, dir, "ASSIGNMENT-3.json", "{broken")
	writeRaw(t, dir, "notes.json", validFile("ASSIGNMENT-4"))
	writeRaw(t, dir, ".atomicfile-1.tmp", "")
	writeRaw(t, dir, ".hidden.json", "")
	writeRaw(t, dir, "README.md", "")
	writeRaw(t, dir, ".lock", "")
	result, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, item := range result.Items {
		ids = append(ids, item.ID)
	}
	if !reflect.DeepEqual(ids, []string{"ASSIGNMENT-2", "ASSIGNMENT-10", "ASSIGNMENT-11"}) {
		t.Fatalf("ids %v", ids)
	}
	want := []InvalidFile{{File: "ASSIGNMENT-3.json", Error: "invalid JSON"}, {File: "notes.json", Error: "invalid file name"}}
	if !reflect.DeepEqual(result.Invalid, want) {
		t.Fatalf("invalid %+v", result.Invalid)
	}
	if got := result.Items[0]; got.Title != "T" || got.CreatedAt != "2026-11-12T09:14:03Z" || !reflect.DeepEqual(got.Repos, []string{"ply"}) {
		t.Fatalf("summary %+v", got)
	}
}

func TestListMissingDirectoryIsEmpty(t *testing.T) {
	s, dir := testStore(t)
	result, err := s.List()
	if err != nil || result.Items == nil || result.Invalid == nil || len(result.Items)+len(result.Invalid) != 0 {
		t.Fatalf("list %+v %v", result, err)
	}
	requireAbsent(t, dir)
}
