//go:build unix

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/backlog"
)

// The Downloads journey runs the shipped main package, or the exact binary in
// PLY_BACKLOG_TEST_BINARY, against an isolated registered workspace.
func TestBacklogFileImportAcceptance(t *testing.T) {
	binary := os.Getenv("PLY_BACKLOG_TEST_BINARY")
	if binary == "" {
		binary = filepath.Join(t.TempDir(), "ply")
		if out, err := exec.Command("go", "build", "-o", binary, "./ply").CombinedOutput(); err != nil {
			t.Fatalf("build actual CLI: %v\n%s", err, out)
		}
	} else if !filepath.IsAbs(binary) {
		t.Fatal("PLY_BACKLOG_TEST_BINARY must be absolute")
	}
	f := newPlanCLI(t, binary)
	ws, downloads := filepath.Join(f.root, "ws"), filepath.Join(f.root, "Downloads")
	repo := filepath.Join(ws, "repo")
	for _, dir := range []string{repo, downloads} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	f.git(t, repo, "init", "-b", "main")
	f.git(t, repo, "commit", "--allow-empty", "-m", "fixture")
	f.success(t, ws, "workspace", "init")
	definitions := `{"format_version":1,"projects":[{"id":"ply","name":"Ply","wrapper":".","data_paths":{"backlog":"workflow/backlog"},"repos":[{"id":"code","path":"repo"}],"epics":[]}]}`
	mustPlanWrite(t, filepath.Join(ws, "ply-workspace.json"), definitions, 0o644)
	backlogDir := filepath.Join(ws, "workflow", "backlog")
	original := []byte("# Min spec\n\xff not UTF-8\n")
	mustPlanWrite(t, filepath.Join(downloads, "min spec.md"), string(original), 0o600)
	backlogAt := func(root string, args ...string) []string {
		return append([]string{"backlog", "--workspace", root}, args...)
	}
	add := []string{"add", "--project", "ply", "--title", "Spec for ny app", "--file", "min spec.md", "--format", "json"}

	f.failure(t, downloads, nil, "workspace_not_found", append([]string{"backlog"}, add...)...)
	link := filepath.Join(f.root, "link")
	if err := os.Symlink(ws, link); err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"ws", link} {
		f.failure(t, downloads, nil, "backlog_invalid_arguments", backlogAt(bad, add...)...)
	}
	f.failure(t, downloads, nil, "backlog_invalid_arguments", backlogAt(ws, "add", "--project", "ply", "--title", "t", "--file", "")...)

	raw := f.success(t, downloads, backlogAt(ws, add...)...)
	item := decodeBacklogItem(t, raw)
	copied := filepath.Join(backlogDir, "files", item.ID, "min spec.md")
	if item.Text != "" || item.File == nil || *item.File != (backlog.File{Name: "min spec.md", Size: int64(len(original))}) || len(raw) > 1024 {
		t.Fatalf("item %s", raw)
	}
	stored, err := os.ReadFile(filepath.Join(backlogDir, item.ID+".json"))
	if err != nil || string(stored) != raw {
		t.Fatalf("item file %q differs from output: %v", stored, err)
	}
	if got, _ := os.ReadFile(filepath.Join(downloads, "min spec.md")); !bytes.Equal(got, original) {
		t.Fatal("original changed")
	}
	note := exec.Command(binary, backlogAt(ws, "note", "--project", "ply", item.ID)...)
	note.Dir, note.Env, note.Stdin = downloads, f.environment(nil), strings.NewReader("Notat\n")
	if out, err := note.CombinedOutput(); err != nil {
		t.Fatalf("note: %v %s", err, out)
	}
	f.success(t, downloads, backlogAt(ws, "status", "--project", "ply", "--value", "ready", item.ID)...)
	if err := os.Remove(filepath.Join(downloads, "min spec.md")); err != nil {
		t.Fatal(err)
	}
	text := f.success(t, downloads, backlogAt(ws, "show", "--project", "ply", item.ID)...)
	wantHead := fmt.Sprintf("%s  Spec for ny app\nCreated: %s\nStatus: ready\nFile: min spec.md (%d bytes, present)\nPath: %s\n\n\n\nNote 1 ", item.ID, item.CreatedAt, len(original), copied)
	if !strings.HasPrefix(text, wantHead) || !strings.Contains(text, "Notat\n") || strings.Contains(text, "Min spec") {
		t.Fatalf("show text %q, want prefix %q", text, wantHead)
	}
	if got, _ := os.ReadFile(copied); !bytes.Equal(got, original) {
		t.Fatal("copy does not hold the original bytes")
	}

	old := "{\n  \"id\": \"20261101-0000000a\",\n  \"title\": \"Old\",\n  \"text\": \"Old text\\n\",\n  \"created_at\": \"2026-11-01T10:00:00Z\",\n  \"status\": \"done\",\n  \"notes\": [\n    {\n      \"at\": \"2026-11-02T10:00:00Z\",\n      \"text\": \"n\"\n    }\n  ]\n}\n"
	writeBacklogItem(t, backlogDir, "20261101-0000000a", old)
	if got := f.success(t, downloads, backlogAt(ws, "show", "--project", "ply", "--format", "json", "20261101-0000000a")...); got != old {
		t.Fatalf("old item JSON %q, file %q", got, old)
	}

	pattern := bytes.Repeat([]byte("ply-backlog-128MiB-pattern-0123\n"), 1<<15) // 1 MiB
	large := filepath.Join(downloads, "large.bin")
	writer, err := os.Create(large)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 128; i++ {
		if _, err := writer.Write(pattern); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	mustPlanWrite(t, filepath.Join(downloads, "small.txt"), "small\n", 0o600)
	baseline, _ := maxRSSImport(t, f, downloads, ws, "small.txt")
	peak, raw := maxRSSImport(t, f, downloads, ws, "large.bin")
	t.Logf("maxrss small=%d large=%d bytes", baseline, peak)
	item = decodeBacklogItem(t, raw)
	if peak > baseline+32<<20 || len(raw) > 1024 || item.File.Size != 128<<20 {
		t.Fatalf("large import maxrss %d over baseline %d + 32 MiB, item %s", peak, baseline, raw)
	}
	for _, args := range [][]string{{"list"}, {"show", item.ID}, {"show", "--format", "json", item.ID}} {
		if out := f.success(t, downloads, backlogAt(ws, append(args, "--project", "ply")...)...); strings.Contains(out, "pattern") {
			t.Fatalf("%v shows file contents", args)
		}
	}
	if _, err := os.Lstat(filepath.Join(ws, "backlog")); err == nil {
		t.Fatal("default backlog directory used instead of data_paths.backlog")
	}
	a, _ := os.Open(large)
	b, _ := os.Open(filepath.Join(backlogDir, "files", item.ID, "large.bin"))
	defer a.Close()
	defer b.Close()
	left, right := make([]byte, 1<<20), make([]byte, 1<<20)
	for {
		n, errA := io.ReadFull(a, left)
		m, _ := io.ReadFull(b, right)
		if n != m || !bytes.Equal(left[:n], right[:m]) {
			t.Fatal("large copy differs from source")
		}
		if errA != nil {
			break
		}
	}
}

// maxRSSImport imports name and returns the child's peak resident set in bytes.
func maxRSSImport(t *testing.T, f planCLIFixture, cwd, ws, name string) (int64, string) {
	t.Helper()
	c := exec.Command(f.binary, "backlog", "add", "--workspace", ws, "--project", "ply", "--title", name, "--file", name, "--format", "json")
	c.Dir, c.Env = cwd, f.environment(nil)
	out, err := c.Output()
	if err != nil {
		t.Fatalf("import %s: %v", name, err)
	}
	rss := c.ProcessState.SysUsage().(*syscall.Rusage).Maxrss
	if runtime.GOOS != "darwin" {
		rss *= 1024 // Linux and the BSDs report KiB
	}
	return int64(rss), string(out)
}
