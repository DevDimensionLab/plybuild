package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
)

const commandPackagePath = "github.com/devdimensionlab/plybuild/cmd"

func TestExecutableEntryBoundary(t *testing.T) {
	repositoryRoot := testRepositoryRoot(t)
	entrypoints := []string{"main.go", "cmd/ply/main.go"}

	if len(entrypoints) == 0 {
		t.Fatal("executable entrypoint population is empty")
	}

	for _, relativePath := range entrypoints {
		relativePath := relativePath
		t.Run(relativePath, func(t *testing.T) {
			assertEntrypoint(t, repositoryRoot, relativePath)
		})
	}

	assertExecutionFunctions(t, repositoryRoot)

	violations, exitCounts := processTerminationViolations(t, repositoryRoot, entrypoints)
	if len(violations) != 0 {
		t.Fatalf("process termination escaped the executable boundaries (%d calls):\n%s",
			len(violations), strings.Join(violations, "\n"))
	}
	for _, relativePath := range entrypoints {
		if exitCounts[relativePath] != 1 {
			t.Errorf("%s main has %d os.Exit calls, want 1", relativePath, exitCounts[relativePath])
		}
	}
}

func testRepositoryRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate entry-boundary test")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), ".."))
}

func assertEntrypoint(t *testing.T, repositoryRoot, relativePath string) {
	t.Helper()
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filepath.Join(repositoryRoot, relativePath), nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", relativePath, err)
	}

	aliases := make(map[string]string)
	for _, imported := range parsed.Imports {
		path, err := strconv.Unquote(imported.Path.Value)
		if err != nil {
			t.Fatalf("unquote import in %s: %v", relativePath, err)
		}
		name := filepath.Base(path)
		if imported.Name != nil {
			name = imported.Name.Name
		}
		aliases[name] = path
	}

	executionPath := ""
	var mainFunction *ast.FuncDecl
	for _, declaration := range parsed.Decls {
		switch declaration := declaration.(type) {
		case *ast.GenDecl:
			for _, specification := range declaration.Specs {
				value, ok := specification.(*ast.ValueSpec)
				if !ok || len(value.Names) != 1 || value.Names[0].Name != "execute" || len(value.Values) != 1 {
					continue
				}
				selector, ok := value.Values[0].(*ast.SelectorExpr)
				if !ok {
					continue
				}
				packageName, ok := selector.X.(*ast.Ident)
				if ok {
					executionPath = aliases[packageName.Name] + "." + selector.Sel.Name
				}
			}
		case *ast.FuncDecl:
			if declaration.Recv == nil && declaration.Name.Name == "main" {
				mainFunction = declaration
			}
		}
	}

	wantPath := commandPackagePath + ".ExecuteE"
	if executionPath != wantPath {
		t.Errorf("%s delegates to %q, want %q", relativePath, executionPath, wantPath)
	}
	if mainFunction == nil || !handlesExecutionError(mainFunction) {
		t.Errorf("%s main does not convert the shared execution error to exit 1", relativePath)
	}
}

func handlesExecutionError(function *ast.FuncDecl) bool {
	handled := false
	ast.Inspect(function.Body, func(node ast.Node) bool {
		statement, ok := node.(*ast.IfStmt)
		if !ok || statement.Init == nil {
			return true
		}
		assignment, ok := statement.Init.(*ast.AssignStmt)
		if !ok || len(assignment.Lhs) != 1 || len(assignment.Rhs) != 1 {
			return true
		}
		errName, ok := assignment.Lhs[0].(*ast.Ident)
		call, callOK := assignment.Rhs[0].(*ast.CallExpr)
		callee, calleeOK := call.Fun.(*ast.Ident)
		if !ok || !callOK || !calleeOK || errName.Name != "err" || callee.Name != "execute" {
			return true
		}
		condition, ok := statement.Cond.(*ast.BinaryExpr)
		left, leftOK := condition.X.(*ast.Ident)
		right, rightOK := condition.Y.(*ast.Ident)
		if !ok || !leftOK || !rightOK || left.Name != errName.Name || right.Name != "nil" {
			return true
		}
		hasExitOne := false
		ast.Inspect(statement.Body, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok || !isSelectorCall(call, "os", "Exit") || len(call.Args) != 1 {
				return true
			}
			literal, ok := call.Args[0].(*ast.BasicLit)
			hasExitOne = ok && literal.Value == "1"
			return true
		})
		handled = hasExitOne
		return true
	})
	return handled
}

func assertExecutionFunctions(t *testing.T, repositoryRoot string) {
	t.Helper()
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filepath.Join(repositoryRoot, "cmd/root.go"), nil, 0)
	if err != nil {
		t.Fatalf("parse cmd/root.go: %v", err)
	}

	foundExecute := false
	foundExecuteE := false
	for _, declaration := range parsed.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Recv != nil {
			continue
		}
		switch function.Name.Name {
		case "Execute":
			foundExecute = function.Type.Results == nil
		case "ExecuteE":
			if function.Type.Results == nil || len(function.Type.Results.List) != 1 {
				continue
			}
			result, ok := function.Type.Results.List[0].Type.(*ast.Ident)
			foundExecuteE = ok && result.Name == "error"
		}
	}
	if !foundExecute {
		t.Error("cmd.Execute no longer has its exported func() compatibility signature")
	}
	if !foundExecuteE {
		t.Error("cmd.ExecuteE error-returning command path is missing")
	}
}

func processTerminationViolations(t *testing.T, repositoryRoot string, entrypoints []string) ([]string, map[string]int) {
	t.Helper()
	if len(entrypoints) == 0 {
		t.Fatal("process-termination entrypoint population is empty")
	}
	allowedEntrypoints := make(map[string]bool, len(entrypoints))
	for _, entrypoint := range entrypoints {
		allowedEntrypoints[entrypoint] = true
	}
	exitCounts := make(map[string]int, len(entrypoints))
	violations := make([]string, 0)

	err := filepath.Walk(repositoryRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "target" || info.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		relativePath, err := filepath.Rel(repositoryRoot, path)
		if err != nil {
			return err
		}
		fileSet := token.NewFileSet()
		parsed, err := parser.ParseFile(fileSet, path, nil, 0)
		if err != nil {
			return err
		}
		allowedExits := make(map[token.Pos]bool)
		for _, declaration := range parsed.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || function.Body == nil || !allowedEntrypoints[relativePath] ||
				function.Recv != nil || function.Name.Name != "main" {
				continue
			}
			ast.Inspect(function.Body, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if ok && isSelectorCall(call, "os", "Exit") {
					allowedExits[call.Pos()] = true
				}
				return true
			})
		}
		ast.Inspect(parsed, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			position := fileSet.Position(call.Pos())
			location := fmt.Sprintf("%s:%d", relativePath, position.Line)
			if isFatalCall(call) {
				violations = append(violations, location+" log.Fatal")
			}
			if isSelectorCall(call, "os", "Exit") {
				if allowedExits[call.Pos()] {
					exitCounts[relativePath]++
				} else {
					violations = append(violations, location+" os.Exit")
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("scan production Go files: %v", err)
	}
	sort.Strings(violations)
	return violations, exitCounts
}

func isFatalCall(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	receiver, ok := selector.X.(*ast.Ident)
	return ok && receiver.Name == "log" && strings.HasPrefix(selector.Sel.Name, "Fatal")
}

func isSelectorCall(call *ast.CallExpr, receiverName, methodName string) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != methodName {
		return false
	}
	receiver, ok := selector.X.(*ast.Ident)
	return ok && receiver.Name == receiverName
}
