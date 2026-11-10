package cmd

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

func TestInteractiveBuildPreservesRootOptionsEndpointCallbackAndValidationOrder(t *testing.T) {
	parsed := parseInteractiveCommandSource(t, "build.go")
	function := findInteractiveCommandFunction(t, parsed, "interactiveWebService")
	if function.Type.Params == nil || len(function.Type.Params.List) != 1 ||
		function.Type.Results == nil || len(function.Type.Results.List) != 1 ||
		function.Body == nil || len(function.Body.List) != 5 {
		t.Fatal("interactiveWebService signature or five-step composition changed")
	}

	rootAssignment, ok := function.Body.List[0].(*ast.AssignStmt)
	if !ok || rootAssignment.Tok != token.DEFINE || len(rootAssignment.Lhs) != 2 ||
		len(rootAssignment.Rhs) != 1 || !interactiveCommandIdentifier(rootAssignment.Lhs[0], "ioResp") ||
		!interactiveCommandIdentifier(rootAssignment.Lhs[1], "err") ||
		interactiveCommandCallName(rootAssignment.Rhs[0]) != "spring.GetRoot" {
		t.Fatal("interactiveWebService no longer selects the Spring root response first")
	}
	if !interactiveCommandErrorReturn(function.Body.List[1], "err") {
		t.Fatal("interactiveWebService no longer returns the exact Spring root error before mutation")
	}

	optionsAssignment, ok := function.Body.List[2].(*ast.AssignStmt)
	if !ok || optionsAssignment.Tok != token.ASSIGN || len(optionsAssignment.Lhs) != 1 ||
		len(optionsAssignment.Rhs) != 1 || interactiveCommandCallName(optionsAssignment.Lhs[0]) != "api.GOptions" {
		t.Fatal("interactiveWebService Generate options assignment changed")
	}
	options, ok := optionsAssignment.Rhs[0].(*ast.CompositeLit)
	if !ok || interactiveCommandCallName(options.Type) != "api.GenerateOptions" || len(options.Elts) != 3 {
		t.Fatal("interactiveWebService no longer installs one complete GenerateOptions value")
	}
	wantOptions := []struct {
		field string
		value string
	}{
		{field: "ProjectConfig", value: "orderConfig"},
		{field: "CloudConfig", value: "ctx.CloudConfig"},
		{field: "IoResponse", value: "ioResp"},
	}
	if len(wantOptions) == 0 {
		t.Fatal("interactive build option contract population is empty")
	}
	for index, want := range wantOptions {
		field, ok := options.Elts[index].(*ast.KeyValueExpr)
		if !ok || interactiveCommandCallName(field.Key) != want.field ||
			interactiveCommandCallName(field.Value) != want.value {
			t.Fatalf("interactiveWebService option %d is %#v, want %s: %s", index, options.Elts[index], want.field, want.value)
		}
	}

	standaloneStatement, ok := function.Body.List[3].(*ast.ExprStmt)
	standalone, callOK := standaloneStatement.X.(*ast.CallExpr)
	if !ok || !callOK || interactiveCommandCallName(standalone.Fun) != "webservice.InitAndBlockStandalone" ||
		len(standalone.Args) != 2 || interactiveCommandCallName(standalone.Args[0]) != "webservice.Generate" ||
		interactiveCommandCallName(standalone.Args[1]) != "api.CallbackChannel" {
		t.Fatal("interactiveWebService no longer blocks on the exact Generate endpoint and API callback channel")
	}
	if !interactiveCommandNilReturn(function.Body.List[4]) {
		t.Fatal("interactiveWebService no longer returns exact nil after the blocking callback")
	}

	runE := findInteractiveCommandCallback(t, parsed, "buildCmd", "RunE")
	interactiveIndex := -1
	for index, statement := range runE.Body.List {
		candidate, ok := statement.(*ast.IfStmt)
		if ok && interactiveCommandIdentifier(candidate.Cond, "interactive") {
			if interactiveIndex != -1 {
				t.Fatal("build command contains more than one interactive branch")
			}
			interactiveIndex = index
		}
	}
	if interactiveIndex < 0 || interactiveIndex+1 >= len(runE.Body.List) {
		t.Fatal("build command interactive branch or following validation is missing")
	}
	interactiveBranch := runE.Body.List[interactiveIndex].(*ast.IfStmt)
	if len(interactiveBranch.Body.List) != 1 {
		t.Fatal("build command interactive branch no longer has one error-returning call")
	}
	interactiveError, ok := interactiveBranch.Body.List[0].(*ast.IfStmt)
	if !ok || !interactiveCommandErrorReturn(interactiveError, "err") {
		t.Fatal("build command no longer returns the exact interactiveWebService error")
	}
	interactiveInit, ok := interactiveError.Init.(*ast.AssignStmt)
	if !ok || interactiveInit.Tok != token.DEFINE || len(interactiveInit.Lhs) != 1 ||
		len(interactiveInit.Rhs) != 1 || !interactiveCommandIdentifier(interactiveInit.Lhs[0], "err") {
		t.Fatal("build command interactive error initialization changed")
	}
	interactiveCall, ok := interactiveInit.Rhs[0].(*ast.CallExpr)
	if !ok || interactiveCommandCallName(interactiveCall.Fun) != "interactiveWebService" || len(interactiveCall.Args) != 1 {
		t.Fatal("build command no longer calls interactiveWebService exactly once")
	}
	orderConfig, ok := interactiveCall.Args[0].(*ast.UnaryExpr)
	if !ok || orderConfig.Op != token.AND || !interactiveCommandIdentifier(orderConfig.X, "orderConfig") {
		t.Fatal("build command no longer passes the exact orderConfig pointer to interactiveWebService")
	}
	validation, ok := runE.Body.List[interactiveIndex+1].(*ast.IfStmt)
	validationInit, initOK := validation.Init.(*ast.AssignStmt)
	if !ok || !initOK || validationInit.Tok != token.ASSIGN || len(validationInit.Rhs) != 1 ||
		interactiveCommandCallName(validationInit.Rhs[0]) != "orderConfig.Validate" ||
		!interactiveCommandErrorReturn(validation, "err") {
		t.Fatal("build command no longer validates immediately after the optional interactive caller")
	}
}

func TestInteractiveUpgradePreservesProjectGuardsEndpointAndCallbackCaller(t *testing.T) {
	parsed := parseInteractiveCommandSource(t, "upgrade.go")
	runE := findInteractiveCommandCallback(t, parsed, "upgradeInteractiveCmd", "RunE")
	if runE.Body == nil || len(runE.Body.List) != 4 {
		t.Fatal("interactive upgrade guard, caller, or return flow changed")
	}
	if _, ok := runE.Body.List[0].(*ast.IfStmt); !ok {
		t.Fatal("interactive upgrade no longer checks for an empty project population first")
	}
	if _, ok := runE.Body.List[1].(*ast.IfStmt); !ok {
		t.Fatal("interactive upgrade no longer checks the root project type second")
	}
	rootStatement, ok := runE.Body.List[2].(*ast.ExprStmt)
	rootCall, callOK := rootStatement.X.(*ast.CallExpr)
	if !ok || !callOK || interactiveCommandCallName(rootCall.Fun) != "ctx.OnRootProject" || len(rootCall.Args) != 2 {
		t.Fatal("interactive upgrade no longer delegates one callback through OnRootProject after its guards")
	}
	description, ok := rootCall.Args[0].(*ast.BasicLit)
	if !ok || description.Kind != token.STRING || description.Value != `"starting interactive upgrade"` {
		t.Fatal("interactive upgrade root-project description changed")
	}
	callback, ok := rootCall.Args[1].(*ast.CallExpr)
	if !ok || interactiveCommandCallName(callback.Fun) != "webservice.InitAndBlockProject" ||
		len(callback.Args) != 2 || interactiveCommandCallName(callback.Args[0]) != "webservice.Upgrade" ||
		interactiveCommandCallName(callback.Args[1]) != "api.CallbackChannel" {
		t.Fatal("interactive upgrade no longer selects the exact Upgrade endpoint and API callback channel")
	}
	if !interactiveCommandNilReturn(runE.Body.List[3]) {
		t.Fatal("interactive upgrade no longer returns exact nil after scheduling its blocking callback")
	}
}

func parseInteractiveCommandSource(t *testing.T, relative string) *ast.File {
	t.Helper()
	parsed, err := parser.ParseFile(token.NewFileSet(), filepath.Join(testRepositoryRoot(t), "cmd", relative), nil, 0)
	if err != nil {
		t.Fatalf("parse interactive command source %s: %v", relative, err)
	}
	return parsed
}

func findInteractiveCommandFunction(t *testing.T, parsed *ast.File, name string) *ast.FuncDecl {
	t.Helper()
	for _, declaration := range parsed.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && function.Recv == nil && function.Name.Name == name {
			return function
		}
	}
	t.Fatalf("interactive command function %s is missing", name)
	return nil
}

func findInteractiveCommandCallback(t *testing.T, parsed *ast.File, variable, field string) *ast.FuncLit {
	t.Helper()
	for _, declaration := range parsed.Decls {
		general, ok := declaration.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, specification := range general.Specs {
			value, ok := specification.(*ast.ValueSpec)
			if !ok || len(value.Names) != 1 || value.Names[0].Name != variable || len(value.Values) != 1 {
				continue
			}
			address, ok := value.Values[0].(*ast.UnaryExpr)
			if !ok || address.Op != token.AND {
				continue
			}
			literal, ok := address.X.(*ast.CompositeLit)
			if !ok {
				continue
			}
			for _, element := range literal.Elts {
				pair, ok := element.(*ast.KeyValueExpr)
				if !ok || interactiveCommandCallName(pair.Key) != field {
					continue
				}
				callback, ok := pair.Value.(*ast.FuncLit)
				if ok {
					return callback
				}
			}
		}
	}
	t.Fatalf("interactive command callback %s.%s is missing", variable, field)
	return nil
}

func interactiveCommandErrorReturn(statement ast.Stmt, name string) bool {
	conditional, ok := statement.(*ast.IfStmt)
	if !ok {
		return false
	}
	condition, ok := conditional.Cond.(*ast.BinaryExpr)
	if !ok || condition.Op != token.NEQ || !interactiveCommandIdentifier(condition.X, name) ||
		!interactiveCommandIdentifier(condition.Y, "nil") || len(conditional.Body.List) != 1 {
		return false
	}
	returned, ok := conditional.Body.List[0].(*ast.ReturnStmt)
	return ok && len(returned.Results) == 1 && interactiveCommandIdentifier(returned.Results[0], name)
}

func interactiveCommandNilReturn(statement ast.Stmt) bool {
	returned, ok := statement.(*ast.ReturnStmt)
	return ok && len(returned.Results) == 1 && interactiveCommandIdentifier(returned.Results[0], "nil")
}

func interactiveCommandCallName(expression ast.Expr) string {
	switch value := expression.(type) {
	case *ast.Ident:
		return value.Name
	case *ast.SelectorExpr:
		return interactiveCommandCallName(value.X) + "." + value.Sel.Name
	case *ast.CallExpr:
		return interactiveCommandCallName(value.Fun)
	default:
		return ""
	}
}

func interactiveCommandIdentifier(expression ast.Expr, want string) bool {
	identifier, ok := expression.(*ast.Ident)
	return ok && identifier.Name == want
}
