package webservice

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"testing"
)

func TestInteractiveWebserviceCallersPreserveEndpointServerBrowserAndCallbackComposition(t *testing.T) {
	endpoints := []struct {
		endpoint Endpoint
		want     string
	}{
		{endpoint: Generate, want: "http://localhost:7999/ui/generate"},
		{endpoint: Upgrade, want: "http://localhost:7999/ui/upgrade"},
	}
	if len(endpoints) == 0 {
		t.Fatal("interactive endpoint contract population is empty")
	}
	for _, endpoint := range endpoints {
		if got := endpoint.endpoint.Uri(); got != endpoint.want {
			t.Fatalf("interactive endpoint %q URI is %q, want exact %q", endpoint.endpoint, got, endpoint.want)
		}
	}

	parsed := parseInteractiveWebserviceSource(t)
	standalone := findInteractiveWebserviceFunction(t, parsed, "InitAndBlockStandalone")
	if standalone.Type.Params == nil || len(standalone.Type.Params.List) != 2 ||
		standalone.Type.Results != nil || standalone.Body == nil || len(standalone.Body.List) != 3 {
		t.Fatal("InitAndBlockStandalone signature or three-step composition changed")
	}
	start, ok := standalone.Body.List[0].(*ast.GoStmt)
	if !ok || interactiveWebserviceCallName(start.Call.Fun) != "StartWebServer" || len(start.Call.Args) != 0 {
		t.Fatal("InitAndBlockStandalone no longer starts StartWebServer in its existing goroutine")
	}
	openAssignment, ok := standalone.Body.List[1].(*ast.AssignStmt)
	if !ok || openAssignment.Tok != token.ASSIGN || len(openAssignment.Lhs) != 1 ||
		len(openAssignment.Rhs) != 1 || !interactiveWebserviceIdentifier(openAssignment.Lhs[0], "_") {
		t.Fatal("InitAndBlockStandalone browser error-ignore assignment changed")
	}
	openCall, ok := openAssignment.Rhs[0].(*ast.CallExpr)
	if !ok || interactiveWebserviceCallName(openCall.Fun) != "OpenBrowser" || len(openCall.Args) != 1 {
		t.Fatal("InitAndBlockStandalone no longer makes one OpenBrowser request")
	}
	uriCall, ok := openCall.Args[0].(*ast.CallExpr)
	if !ok || interactiveWebserviceCallName(uriCall.Fun) != "endpoint.Uri" || len(uriCall.Args) != 0 {
		t.Fatal("InitAndBlockStandalone no longer opens the exact selected endpoint URI")
	}
	blocking, ok := standalone.Body.List[2].(*ast.ExprStmt)
	receive, receiveOK := blocking.X.(*ast.UnaryExpr)
	if !ok || !receiveOK || receive.Op != token.ARROW ||
		!interactiveWebserviceIdentifier(receive.X, "blockingChannel") {
		t.Fatal("InitAndBlockStandalone no longer blocks on the exact caller channel")
	}

	project := findInteractiveWebserviceFunction(t, parsed, "InitAndBlockProject")
	if project.Type.Params == nil || len(project.Type.Params.List) != 2 ||
		project.Type.Results == nil || len(project.Type.Results.List) != 1 ||
		project.Body == nil || len(project.Body.List) != 1 {
		t.Fatal("InitAndBlockProject signature or returned callback composition changed")
	}
	returned, ok := project.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(returned.Results) != 1 {
		t.Fatal("InitAndBlockProject no longer returns one project callback")
	}
	callback, ok := returned.Results[0].(*ast.FuncLit)
	if !ok || callback.Body == nil || len(callback.Body.List) != 3 {
		t.Fatal("InitAndBlockProject returned callback no longer has three exact steps")
	}
	assignment, ok := callback.Body.List[0].(*ast.AssignStmt)
	if !ok || assignment.Tok != token.ASSIGN || len(assignment.Lhs) != 1 || len(assignment.Rhs) != 1 ||
		interactiveWebserviceCallName(assignment.Lhs[0]) != "api.CurrentProject" ||
		!interactiveWebserviceIdentifier(assignment.Rhs[0], "project") {
		t.Fatal("InitAndBlockProject no longer installs the exact callback project first")
	}
	standaloneCall, ok := callback.Body.List[1].(*ast.ExprStmt)
	call, callOK := standaloneCall.X.(*ast.CallExpr)
	if !ok || !callOK || interactiveWebserviceCallName(call.Fun) != "InitAndBlockStandalone" ||
		len(call.Args) != 2 || !interactiveWebserviceIdentifier(call.Args[0], "endpoint") ||
		!interactiveWebserviceIdentifier(call.Args[1], "blockingChannel") {
		t.Fatal("InitAndBlockProject no longer delegates exact endpoint and callback channel")
	}
	callbackReturn, ok := callback.Body.List[2].(*ast.ReturnStmt)
	if !ok || len(callbackReturn.Results) != 1 || !interactiveWebserviceIdentifier(callbackReturn.Results[0], "nil") {
		t.Fatal("InitAndBlockProject callback no longer returns exact nil after blocking")
	}
}

func parseInteractiveWebserviceSource(t *testing.T) *ast.File {
	t.Helper()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate interactive webservice caller contract")
	}
	parsed, err := parser.ParseFile(token.NewFileSet(), filepath.Join(filepath.Dir(testFile), "init.go"), nil, 0)
	if err != nil {
		t.Fatalf("parse interactive webservice callers: %v", err)
	}
	return parsed
}

func findInteractiveWebserviceFunction(t *testing.T, parsed *ast.File, name string) *ast.FuncDecl {
	t.Helper()
	for _, declaration := range parsed.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && function.Recv == nil && function.Name.Name == name {
			return function
		}
	}
	t.Fatalf("interactive webservice function %s is missing", name)
	return nil
}

func interactiveWebserviceCallName(expression ast.Expr) string {
	switch value := expression.(type) {
	case *ast.Ident:
		return value.Name
	case *ast.SelectorExpr:
		return interactiveWebserviceCallName(value.X) + "." + value.Sel.Name
	default:
		return ""
	}
}

func interactiveWebserviceIdentifier(expression ast.Expr, want string) bool {
	identifier, ok := expression.(*ast.Ident)
	return ok && identifier.Name == want
}
