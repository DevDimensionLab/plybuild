package architecture

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

const fixtureModule = "github.com/devdimensionlab/plybuild"

var fixtureNames = []string{"ProjectID", "RepoID", "EpicID", "TaskID", "WorktreeID", "WorktreeOperationID", "WorkItemState", "EpicWorktreeBinding", "EpicRepoBinding", "EpicRecord", "TaskWorktreeBinding", "TaskRecord"}

func legalFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	put(t, root, "go.mod", "module "+fixtureModule+"\n")
	core := "package workitem\nimport \"fmt\"\nvar _ = fmt.Errorf\n"
	aliases := "package workspace\nimport wi \"" + fixtureModule + "/internal/workitem\"\n"
	for _, name := range fixtureNames {
		core += "type " + name + " string\n"
		aliases += "type " + name + " = wi." + name + "\n"
	}
	put(t, root, "internal/workitem/model.go", core)
	put(t, root, "internal/workspace/aliases.go", aliases)
	put(t, root, "internal/workspace/content_physical_read.go", wr01FixtureRead)
	put(t, root, "internal/workspaceview/view.go", "package workspaceview\nimport _ \""+fixtureModule+"/internal/workitem\"\n")
	return root
}
func put(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
func expectProblems(t *testing.T, root string, want ...string) {
	t.Helper()
	sort.Strings(want)
	got := scanRepository(root, scanIO{}).problems
	if len(got) == 0 && len(want) == 0 {
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("problems:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}
func TestScannerLegalAndExcludedFiles(t *testing.T) {
	root := legalFixture(t)
	for _, path := range []string{"internal/workitem/illegal_test.go", "testdata/illegal.go", "vendor/illegal.go", ".git/illegal.go", "target/illegal.go", "internal/workspace/testdata/illegal.go"} {
		put(t, root, path, "invalid Go and forbidden types/imports")
	}
	put(t, root, "vendor/nested/go.mod", "module ignored")
	expectProblems(t, root)
}
func TestScannerR1AllImportForms(t *testing.T) {
	for _, path := range []string{fixtureModule + "/internal/workspace", "os", "os/exec", "net/http", "example.org/thirdparty", "unsafe", fixtureModule + "/internal/workitem/subpackage"} {
		for _, alias := range []string{"", "renamed ", ". ", "_ "} {
			t.Run(path+"/"+alias, func(t *testing.T) {
				root := legalFixture(t)
				put(t, root, "internal/workitem/bad.go", "package workitem\nimport "+alias+"\""+path+"\"\n")
				expectProblems(t, root, "internal/workitem/bad.go:2:8: WI01-R1 "+fixtureModule+"/internal/workitem imports "+path)
			})
		}
	}
}
func TestScannerR2PathsAndBoundaries(t *testing.T) {
	for _, p := range []string{"workspaceview"} {
		for _, suffix := range []string{"", "/sub"} {
			for _, alias := range []string{"", "alias ", ". ", "_ "} {
				t.Run(p+suffix+alias, func(t *testing.T) {
					root := legalFixture(t)
					path := fixtureModule + "/internal/" + p + suffix
					put(t, root, "internal/workspace/bad.go", "package workspace\nimport "+alias+"\""+path+"\"\n")
					expectProblems(t, root, "internal/workspace/bad.go:2:8: WI01-R2 "+fixtureModule+"/internal/workspace imports "+path)
				})
			}
		}
	}
	root := legalFixture(t)
	put(t, root, "internal/workspace/near.go", "package workspace\nimport _ \""+fixtureModule+"/internal/workspaceviews\"\n")
	expectProblems(t, root)
}
func TestScannerR3Ownership(t *testing.T) {
	for _, definition := range []string{"type TaskID string", "type TaskID struct{}", "type TaskID wi.TaskID", "type TaskID = wi.RepoID", "type TaskID = wi.TaskID"} {
		t.Run(definition, func(t *testing.T) {
			root := legalFixture(t)
			put(t, root, "other/bad.go", "package other\nimport wi \""+fixtureModule+"/internal/workitem\"\n"+definition+"\n")
			expectProblems(t, root, "other/bad.go:3:6: WI01-R3 "+fixtureModule+"/other declares reserved type TaskID outside its owner")
		})
	}
	for _, rhs := range []string{"wi.RepoID", "TaskID", "TaskIDCopy"} {
		t.Run(rhs, func(t *testing.T) {
			root := legalFixture(t)
			p := filepath.Join(root, "internal/workspace/aliases.go")
			b, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			put(t, root, "internal/workspace/aliases.go", strings.Replace(string(b), "type TaskID = wi.TaskID", "type TaskID = "+rhs, 1))
			expectProblems(t, root,
				"internal/workspace:1:1: WI01-R3 "+fixtureModule+"/internal/workspace missing direct alias TaskID",
				"internal/workspace/aliases.go:6:6: WI01-R3 "+fixtureModule+"/internal/workspace invalid alias TaskID = "+rhs)
		})
	}
	root := legalFixture(t)
	p := filepath.Join(root, "internal/workspace/aliases.go")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	put(t, root, "internal/workspace/aliases.go", strings.Replace(string(b), "type TaskID = wi.TaskID", "type TaskID wi.TaskID", 1))
	expectProblems(t, root, "internal/workspace:1:1: WI01-R3 "+fixtureModule+"/internal/workspace missing direct alias TaskID", "internal/workspace/aliases.go:6:6: WI01-R3 "+fixtureModule+"/internal/workspace invalid alias TaskID = wi.TaskID")
}
func TestScannerMissingDefinitionsAndDotAlias(t *testing.T) {
	for _, tc := range []struct{ path, from, to, want string }{
		{"internal/workitem/model.go", "type TaskID string\n", "", "internal/workitem:1:1: WI01-R3 " + fixtureModule + "/internal/workitem missing canonical type TaskID"},
		{"internal/workspace/aliases.go", "type TaskID = wi.TaskID\n", "", "internal/workspace:1:1: WI01-R3 " + fixtureModule + "/internal/workspace missing direct alias TaskID"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			root := legalFixture(t)
			b, err := os.ReadFile(filepath.Join(root, tc.path))
			if err != nil {
				t.Fatal(err)
			}
			put(t, root, tc.path, strings.Replace(string(b), tc.from, tc.to, 1))
			expectProblems(t, root, tc.want)
		})
	}
	root := legalFixture(t)
	put(t, root, "internal/workspace/dot.go", "package workspace\nimport . \""+fixtureModule+"/internal/workitem\"\ntype TaskRecord = TaskRecord\n")
	expectProblems(t, root, "internal/workspace/dot.go:3:6: WI01-R3 "+fixtureModule+"/internal/workspace invalid alias TaskRecord = TaskRecord")
}
func TestScannerIncludesEveryBuildVariantAndGeneratedFile(t *testing.T) {
	root := legalFixture(t)
	for _, tc := range []struct {
		file, header string
		line         int
	}{
		{"bad_windows.go", "", 2}, {"bad_linux.go", "", 2}, {"bad_custom.go", "//go:build wi01_inactive\n\n", 4}, {"bad_generated.go", "// Code generated by WI01 fixture. DO NOT EDIT.\n", 3},
	} {
		put(t, root, "internal/workitem/sub/"+tc.file, tc.header+"package workitem\nimport \"os\"\ntype TaskID string\n")
	}
	got := scanRepository(root, scanIO{}).problems
	// Both import and type scanning must cover the same four files.
	if len(got) != 8 {
		t.Fatalf("variant findings = %v", got)
	}
	for _, want := range []string{
		"internal/workitem/sub/bad_custom.go:4:8: WI01-R1", "internal/workitem/sub/bad_custom.go:5:6: WI01-R3",
		"internal/workitem/sub/bad_generated.go:3:8: WI01-R1", "internal/workitem/sub/bad_generated.go:4:6: WI01-R3",
		"internal/workitem/sub/bad_linux.go:2:8: WI01-R1", "internal/workitem/sub/bad_linux.go:3:6: WI01-R3",
		"internal/workitem/sub/bad_windows.go:2:8: WI01-R1", "internal/workitem/sub/bad_windows.go:3:6: WI01-R3",
	} {
		found := false
		for _, p := range got {
			found = found || strings.HasPrefix(p, want)
		}
		if !found {
			t.Errorf("missing %s in %v", want, got)
		}
	}
}
func TestScannerRejectsIncompleteAndUnreadableRepositories(t *testing.T) {
	t.Run("missing root", func(t *testing.T) {
		r := scanRepository(filepath.Join(t.TempDir(), "missing"), scanIO{})
		if !containsProblem(r.problems, "WI01-SCAN walk", "no such file") {
			t.Fatal(r.problems)
		}
	})
	for _, pkg := range []string{"workitem", "workspace", "workspaceview"} {
		for _, removeDir := range []bool{false, true} {
			t.Run(pkg+map[bool]string{true: "missing", false: "empty"}[removeDir], func(t *testing.T) {
				root := legalFixture(t)
				dir := filepath.Join(root, "internal", pkg)
				if err := os.RemoveAll(dir); err != nil {
					t.Fatal(err)
				}
				if !removeDir {
					if err := os.MkdirAll(dir, 0755); err != nil {
						t.Fatal(err)
					}
				}
				r := scanRepository(root, scanIO{})
				if !containsProblem(r.problems, "WI01-SCAN "+fixtureModule+"/internal/"+pkg+" has no production Go files") {
					t.Fatal(r.problems)
				}
			})
		}
	}
	t.Run("invalid Go", func(t *testing.T) {
		root := legalFixture(t)
		put(t, root, "other/bad.go", "package other\ntype\n")
		r := scanRepository(root, scanIO{})
		if !containsProblem(r.problems, "other/bad.go:2:", "WI01-SCAN parse") {
			t.Fatal(r.problems)
		}
	})
	t.Run("nested module", func(t *testing.T) {
		root := legalFixture(t)
		put(t, root, "nested/go.mod", "module example.org/nested\n")
		expectProblems(t, root, "nested/go.mod:1:1: WI01-SCAN nested module in production scan")
	})
	t.Run("read failure", func(t *testing.T) {
		root := legalFixture(t)
		r := scanRepository(root, scanIO{readFile: func(p string) ([]byte, error) {
			if strings.HasSuffix(p, "model.go") {
				return nil, errors.New("injected read failure")
			}
			return os.ReadFile(p)
		}})
		if !containsProblem(r.problems, "internal/workitem/model.go:1:1:", "WI01-SCAN read: injected read failure") {
			t.Fatal(r.problems)
		}
	})
	t.Run("walk failure", func(t *testing.T) {
		root := legalFixture(t)
		r := scanRepository(root, scanIO{walkDir: func(_ string, fn fs.WalkDirFunc) error {
			return fn(filepath.Join(root, "broken"), nil, errors.New("injected walk failure"))
		}})
		if !containsProblem(r.problems, "broken:1:1: WI01-SCAN walk: injected walk failure") {
			t.Fatal(r.problems)
		}
	})
	t.Run("walk return failure", func(t *testing.T) {
		root := legalFixture(t)
		r := scanRepository(root, scanIO{walkDir: func(string, fs.WalkDirFunc) error { return errors.New("injected walker failure") }})
		if !containsProblem(r.problems, "WI01-SCAN walk: injected walker failure") {
			t.Fatal(r.problems)
		}
	})
}
func containsProblem(problems []string, parts ...string) bool {
	for _, p := range problems {
		found := true
		for _, part := range parts {
			found = found && strings.Contains(p, part)
		}
		if found {
			return true
		}
	}
	return false
}

func TestScannerVariantDefinitionsRemainInTheirOwner(t *testing.T) {
	root := legalFixture(t)
	// Presence is checked across variants; duplicate declarations in mutually
	// exclusive source files belong to Go's compiler, not this ownership rule.
	put(t, root, "internal/workitem/id_windows.go", "package workitem\ntype TaskID string\n")
	put(t, root, "internal/workspace/id_linux.go", "package workspace\nimport core \""+fixtureModule+"/internal/workitem\"\ntype TaskID = core.TaskID\n")
	expectProblems(t, root)
}

func TestScannerCanonicalAliasAndLocalTypeCannotSupplyDefinition(t *testing.T) {
	for _, replacement := range []string{"type TaskID = string", "func local() { type TaskID string }"} {
		t.Run(replacement, func(t *testing.T) {
			root := legalFixture(t)
			b, err := os.ReadFile(filepath.Join(root, "internal/workitem/model.go"))
			if err != nil {
				t.Fatal(err)
			}
			put(t, root, "internal/workitem/model.go", strings.Replace(string(b), "type TaskID string", replacement, 1))
			want := []string{"internal/workitem:1:1: WI01-R3 " + fixtureModule + "/internal/workitem missing canonical type TaskID"}
			if strings.Contains(replacement, "=") {
				want = append(want, "internal/workitem/model.go:7:6: WI01-R3 "+fixtureModule+"/internal/workitem canonical type TaskID must be a definition")
			}
			expectProblems(t, root, want...)
		})
	}
}

func TestScannerR2IncludesWorkspaceSubpackages(t *testing.T) {
	root := legalFixture(t)
	put(t, root, "internal/workspace/sub/bad.go", "package sub\nimport _ \""+fixtureModule+"/internal/workspaceview\"\n")
	expectProblems(t, root, "internal/workspace/sub/bad.go:2:8: WI01-R2 "+fixtureModule+"/internal/workspace/sub imports "+fixtureModule+"/internal/workspaceview")
}

func TestScannerUnquotesActualImportPaths(t *testing.T) {
	root := legalFixture(t)
	put(t, root, "internal/workitem/escaped.go", "package workitem\nimport \"\\x6fs\"\n")
	put(t, root, "internal/workspace/escaped.go", "package workspace\nimport _ \""+fixtureModule+"/internal/\\x77orkspaceview\"\n")
	expectProblems(t, root,
		"internal/workitem/escaped.go:2:8: WI01-R1 "+fixtureModule+"/internal/workitem imports os",
		"internal/workspace/escaped.go:2:8: WI01-R2 "+fixtureModule+"/internal/workspace imports "+fixtureModule+"/internal/workspaceview")
}
