package inbox

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// C4: the store imports no notification, workflow, network or process package.
func TestInboxImportsNoNotificationWorkflowNetworkOrProcess(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	allowedModule := map[string]bool{
		"github.com/devdimensionlab/plybuild/internal/atomicfile":    true,
		"github.com/devdimensionlab/plybuild/internal/canonicaljson": true,
		"golang.org/x/sys/unix":    true,
		"golang.org/x/sys/windows": true,
	}
	checked := 0
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
			checked++
			switch {
			case path == "net" || strings.HasPrefix(path, "net/") || path == "os/exec" || path == "os/signal" || path == "plugin":
				t.Errorf("%s imports %s", file, path)
			case strings.Contains(path, "."):
				if !allowedModule[path] {
					t.Errorf("%s imports %s", file, path)
				}
			}
		}
	}
	if checked == 0 {
		t.Fatal("no imports checked")
	}
}
