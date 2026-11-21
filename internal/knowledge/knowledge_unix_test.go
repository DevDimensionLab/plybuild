//go:build unix

package knowledge

import (
	"path/filepath"
	"reflect"
	"syscall"
	"testing"
)

// A FIFO named k-*.md is not opened, so it cannot block the answer.
func TestAFIFOPostIsNotOpened(t *testing.T) {
	dir := writeArchive(t, fixture{path: "a", id: "k-1-kept", title: "chrome"})
	if err := syscall.Mkfifo(filepath.Join(dir, "a", "k-2-fifo.md"), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := Search(dir, Query{Words: []string{"chrome"}})
	if err != nil {
		t.Fatal(err)
	}
	requireIDs(t, result.Items, "k-1-kept")
	if want := []InvalidFile{{File: "a/k-2-fifo.md", Error: "not a regular file"}}; !reflect.DeepEqual(result.Invalid, want) {
		t.Fatalf("invalid %+v", result.Invalid)
	}
}
