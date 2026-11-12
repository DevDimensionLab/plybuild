package plan

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// TestImportsOnlyStandardLibraryAndFileHelpers documents the seam: the package
// takes a plans directory and knows no Project, workspace, Task, Assignment,
// backlog, Herdr, CLI or pkg/* code.
func TestImportsOnlyStandardLibraryAndFileHelpers(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	internal := map[string]bool{}
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatal(err)
		}
		for _, spec := range parsed.Imports {
			path, _ := strconv.Unquote(spec.Path.Value)
			first := strings.SplitN(path, "/", 2)[0]
			if strings.Contains(first, ".") {
				internal[path] = true
			}
		}
	}
	var got []string
	for path := range internal {
		got = append(got, path)
	}
	sort.Strings(got)
	want := []string{"github.com/devdimensionlab/plybuild/internal/atomicfile", "github.com/devdimensionlab/plybuild/internal/filelock"}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Fatalf("non-standard imports %v, want %v", got, want)
	}
}
