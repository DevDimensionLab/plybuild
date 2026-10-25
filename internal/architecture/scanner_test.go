package architecture

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/scanner"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const modulePath = "github.com/devdimensionlab/plybuild"
const physicalReadFile = "internal/workspace/content_physical_read.go"

var ownedTypes = []string{"ProjectID", "RepoID", "EpicID", "TaskID", "WorktreeID", "WorktreeOperationID", "WorkItemState", "EpicWorktreeBinding", "EpicRepoBinding", "EpicRecord", "TaskWorktreeBinding", "TaskRecord"}
var affectedPackages = []string{"internal/workitem", "internal/workspace", "internal/workspaceview", "internal/taskrun", "internal/workflowhandoff"}

// This scanner is test-local: parse all production source files individually,
// including inactive build variants. Never use go/packages' host build filtering.
// target is the reserved build output documented in docs/workitem-core.md.
// Test helpers and deliberately invalid fixtures belong in _test.go/testdata.
type scanIO struct {
	walkDir  func(string, fs.WalkDirFunc) error
	readFile func(string) ([]byte, error)
}
type scanResult struct {
	problems []string
	edges    []string
}

func scanRepository(root string, access scanIO) scanResult {
	if access.walkDir == nil {
		access.walkDir = filepath.WalkDir
	}
	if access.readFile == nil {
		access.readFile = os.ReadFile
	}
	result := scanResult{}
	population := map[string]int{}
	canonical, aliases := map[string]bool{}, map[string]bool{}
	readFileFound := false
	readDefinitions := map[string]int{}
	readOwners := map[string]token.Token{
		"contentPhysical": token.FUNC, "contentReadBounded": token.FUNC,
		"contentReadBoundedWith": token.FUNC, "contentReadFile": token.TYPE, "contentReadOps": token.TYPE,
	}
	edgeSet := map[string]bool{}
	reserved := map[string]bool{}
	for _, name := range ownedTypes {
		reserved[name] = true
	}
	report := func(file string, line, column int, rule, detail string) {
		result.problems = append(result.problems, fmt.Sprintf("%s:%d:%d: %s %s", file, line, column, rule, detail))
	}
	err := access.walkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("relative path %s: %w", path, err)
		}
		file := filepath.ToSlash(relative)
		if walkErr != nil {
			report(file, 1, 1, "WI01-SCAN", "walk: "+walkErr.Error())
			return nil
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "vendor", "testdata", "target":
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Name() == "go.mod" && file != "go.mod" {
			report(file, 1, 1, "WI01-SCAN", "nested module in production scan")
			return nil
		}
		if !strings.HasSuffix(file, ".go") || strings.HasSuffix(file, "_test.go") {
			return nil
		}
		readFileFound = readFileFound || file == physicalReadFile
		directory := filepath.ToSlash(filepath.Dir(relative))
		population[directory]++
		pkg := modulePath + "/" + directory
		data, err := access.readFile(path)
		if err != nil {
			report(file, 1, 1, "WI01-SCAN", "read: "+err.Error())
			return nil
		}
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, file, data, parser.AllErrors)
		if err != nil {
			if parseErrors, ok := err.(scanner.ErrorList); ok {
				for _, problem := range parseErrors {
					report(file, problem.Pos.Line, problem.Pos.Column, "WI01-SCAN", "parse "+pkg+": "+problem.Msg)
				}
			} else {
				report(file, 1, 1, "WI01-SCAN", "parse "+pkg+": "+err.Error())
			}
			return nil
		}
		imports := map[string]string{}
		for _, spec := range node.Imports {
			imported, err := strconv.Unquote(spec.Path.Value)
			position := fset.Position(spec.Pos())
			if err != nil {
				report(file, position.Line, position.Column, "WI01-SCAN", "unquote import: "+err.Error())
				continue
			}
			name := filepath.Base(imported)
			if spec.Name != nil {
				name = spec.Name.Name
			}
			imports[name] = imported
			edgeSet[pkg+" -> "+imported] = true
			if file == physicalReadFile {
				switch imported {
				case "fmt", "io", "os", "path/filepath":
				default:
					report(file, position.Line, position.Column, "WR01-R1", "read file imports "+imported)
				}
			}
			if within(directory, "internal/workitem") && imported != "fmt" {
				report(file, position.Line, position.Column, "WI01-R1", pkg+" imports "+imported)
			}
			if within(directory, "internal/workspace") {
				for _, blocked := range []string{"internal/workspaceview", "internal/taskrun", "internal/workflowhandoff"} {
					if within(imported, modulePath+"/"+blocked) {
						report(file, position.Line, position.Column, "WI01-R2", pkg+" imports "+imported)
					}
				}
			}
		}
		topLevel := map[*ast.TypeSpec]bool{}
		readDefinition := func(name string, kind token.Token, pos token.Pos) {
			if readOwners[name] != kind {
				return
			}
			position := fset.Position(pos)
			if file != physicalReadFile {
				report(file, position.Line, position.Column, "WR01-R1", "definition "+name+" must be in "+physicalReadFile)
				return
			}
			readDefinitions[name]++
			if readDefinitions[name] > 1 {
				report(file, position.Line, position.Column, "WR01-R1", "duplicate definition "+name)
			}
		}
		for _, declaration := range node.Decls {
			if function, ok := declaration.(*ast.FuncDecl); ok && function.Recv == nil {
				readDefinition(function.Name.Name, token.FUNC, function.Pos())
			}
			if group, ok := declaration.(*ast.GenDecl); ok && group.Tok == token.TYPE {
				for _, declaration := range group.Specs {
					topLevel[declaration.(*ast.TypeSpec)] = true
				}
			}
		}
		ast.Inspect(node, func(n ast.Node) bool {
			spec, ok := n.(*ast.TypeSpec)
			if ok && topLevel[spec] {
				readDefinition(spec.Name.Name, token.TYPE, spec.Pos())
			}
			if !ok || !reserved[spec.Name.Name] {
				return true
			}
			name := spec.Name.Name
			position := fset.Position(spec.Pos())
			switch directory {
			case "internal/workitem":
				if spec.Assign.IsValid() {
					report(file, position.Line, position.Column, "WI01-R3", pkg+" canonical type "+name+" must be a definition")
				} else {
					canonical[name] = canonical[name] || topLevel[spec]
				}
			case "internal/workspace":
				selector, ok := spec.Type.(*ast.SelectorExpr)
				valid := false
				if ok {
					qualifier, isIdentifier := selector.X.(*ast.Ident)
					valid = isIdentifier && imports[qualifier.Name] == modulePath+"/internal/workitem" && qualifier.Name != "." && qualifier.Name != "_" && selector.Sel.Name == name
				}
				if !spec.Assign.IsValid() || !valid {
					var rhs bytes.Buffer
					if err := format.Node(&rhs, fset, spec.Type); err != nil {
						report(file, position.Line, position.Column, "WI01-SCAN", "format alias: "+err.Error())
					}
					report(file, position.Line, position.Column, "WI01-R3", pkg+" invalid alias "+name+" = "+strings.Join(strings.Fields(rhs.String()), " "))
				} else {
					aliases[name] = aliases[name] || topLevel[spec]
				}
			default:
				report(file, position.Line, position.Column, "WI01-R3", pkg+" declares reserved type "+name+" outside its owner")
			}
			return true
		})
		return nil
	})
	if err != nil {
		report(".", 1, 1, "WI01-SCAN", "walk: "+err.Error())
	}
	if !readFileFound {
		report(physicalReadFile, 1, 1, "WR01-R1", "missing production read file")
	} else {
		for name := range readOwners {
			if readDefinitions[name] == 0 {
				report(physicalReadFile, 1, 1, "WR01-R1", "missing definition "+name)
			}
		}
	}
	for _, pkg := range affectedPackages {
		if population[pkg] == 0 {
			report(pkg, 1, 1, "WI01-SCAN", modulePath+"/"+pkg+" has no production Go files")
		}
	}
	for _, name := range ownedTypes {
		if !canonical[name] {
			report("internal/workitem", 1, 1, "WI01-R3", modulePath+"/internal/workitem missing canonical type "+name)
		}
		if !aliases[name] {
			report("internal/workspace", 1, 1, "WI01-R3", modulePath+"/internal/workspace missing direct alias "+name)
		}
	}
	sort.Strings(result.problems)
	for edge := range edgeSet {
		result.edges = append(result.edges, edge)
	}
	sort.Strings(result.edges)
	return result
}
func within(path, prefix string) bool { return path == prefix || strings.HasPrefix(path, prefix+"/") }
