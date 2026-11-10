package server

import (
	"context"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

type recordedShutdownRequest struct {
	server  *http.Server
	context context.Context
}

type safeServerContextKey struct{}

type recordingServerOperations struct {
	listenServers    []*http.Server
	shutdownRequests []recordedShutdownRequest
	listenErr        error
	shutdownErr      error
}

func (recording *recordingServerOperations) ListenAndServe(server *http.Server) error {
	recording.listenServers = append(recording.listenServers, server)
	return recording.listenErr
}

func (recording *recordingServerOperations) Shutdown(server *http.Server, ctx context.Context) error {
	recording.shutdownRequests = append(recording.shutdownRequests, recordedShutdownRequest{
		server:  server,
		context: ctx,
	})
	return recording.shutdownErr
}

func (recording *recordingServerOperations) assertedListenServers() ([]*http.Server, error) {
	if len(recording.listenServers) == 0 {
		return nil, errors.New("recorded server listen population is empty")
	}
	return recording.listenServers, nil
}

func (recording *recordingServerOperations) assertedShutdownRequests() ([]recordedShutdownRequest, error) {
	if len(recording.shutdownRequests) == 0 {
		return nil, errors.New("recorded server shutdown population is empty")
	}
	return recording.shutdownRequests, nil
}

func TestDependenciesDefaultToSafeNoServerOperation(t *testing.T) {
	ctx := context.WithValue(context.Background(), safeServerContextKey{}, "must remain unread")
	selectServer := func() *http.Server {
		t.Fatal("zero-value server dependencies selected a server")
		return nil
	}

	listenErr := ListenAndServe(Dependencies{}, selectServer)
	shutdownErr := Shutdown(Dependencies{}, selectServer, ctx)
	nilSelectorListenErr := ListenAndServe(Dependencies{Operations: &recordingServerOperations{}}, nil)
	nilSelectorShutdownErr := Shutdown(Dependencies{Operations: &recordingServerOperations{}}, nil, ctx)

	if listenErr != nil || shutdownErr != nil || nilSelectorListenErr != nil || nilSelectorShutdownErr != nil {
		t.Fatalf("safe-zero server dependencies returned listen=%v shutdown=%v nil-selector-listen=%v nil-selector-shutdown=%v, want deterministic nil errors",
			listenErr, shutdownErr, nilSelectorListenErr, nilSelectorShutdownErr)
	}
}

func TestListenAndServePassesExactServerOnceAndReturnsExactError(t *testing.T) {
	supplied := &http.Server{Addr: "arbitrary supplied server identity"}
	sentinel := errors.New("complete server listen dependency error")
	recording := &recordingServerOperations{listenErr: sentinel}
	dependencies := Dependencies{Operations: recording}
	selections := 0
	selectServer := func() *http.Server {
		selections++
		return supplied
	}
	if dependencies.Operations != recording {
		t.Fatalf("server listen dependency lost its complete value: %#v", dependencies)
	}

	err := ListenAndServe(dependencies, selectServer)

	if err != sentinel {
		t.Fatalf("server listen returned %v, want exact dependency error %v", err, sentinel)
	}
	if selections != 1 {
		t.Fatalf("server listen selected its supplied server %d times, want exactly 1", selections)
	}
	servers, populationErr := recording.assertedListenServers()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if len(servers) != 1 || servers[0] != supplied {
		t.Fatalf("server listen received %#v, want one exact server identity %p", servers, supplied)
	}
	if len(recording.shutdownRequests) != 0 {
		t.Fatalf("server listen invoked unrelated shutdown operations: %#v", recording.shutdownRequests)
	}
}

func TestShutdownPassesExactServerAndContextOnceAndReturnsExactError(t *testing.T) {
	supplied := &http.Server{Addr: "another arbitrary supplied server identity"}
	ctx := context.WithValue(context.Background(), struct{ label string }{"context identity"}, "complete")
	sentinel := errors.New("complete server shutdown dependency error")
	recording := &recordingServerOperations{shutdownErr: sentinel}
	dependencies := Dependencies{Operations: recording}
	selections := 0
	selectServer := func() *http.Server {
		selections++
		return supplied
	}
	if dependencies.Operations != recording {
		t.Fatalf("server shutdown dependency lost its complete value: %#v", dependencies)
	}

	err := Shutdown(dependencies, selectServer, ctx)

	if err != sentinel {
		t.Fatalf("server shutdown returned %v, want exact dependency error %v", err, sentinel)
	}
	if selections != 1 {
		t.Fatalf("server shutdown selected its supplied server %d times, want exactly 1", selections)
	}
	requests, populationErr := recording.assertedShutdownRequests()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if len(requests) != 1 || requests[0].server != supplied || requests[0].context != ctx {
		t.Fatalf("server shutdown received %#v, want one exact server %p and context %#v",
			requests, supplied, ctx)
	}
	if len(recording.listenServers) != 0 {
		t.Fatalf("server shutdown invoked unrelated listen operations: %#v", recording.listenServers)
	}
}

func TestSystemSelectsOneCompleteDirectServerOperationDependency(t *testing.T) {
	dependencies := System()

	if dependencies.Operations == nil {
		t.Fatal("system server operations dependency is empty")
	}
	if reflect.TypeOf(dependencies.Operations) != reflect.TypeOf(systemServer{}) {
		t.Fatalf("system server operations dependency is %T, want %T",
			dependencies.Operations, systemServer{})
	}
	assertSystemServerDirectDelegation(t, "ListenAndServe", false)
	assertSystemServerDirectDelegation(t, "Shutdown", true)
}

func TestRecordedServerOperationsRejectEmptyPopulations(t *testing.T) {
	recording := &recordingServerOperations{}

	if _, err := recording.assertedListenServers(); err == nil {
		t.Fatal("empty recorded server listen population passed")
	}
	if _, err := recording.assertedShutdownRequests(); err == nil {
		t.Fatal("empty recorded server shutdown population passed")
	}
}

func assertSystemServerDirectDelegation(t *testing.T, method string, withContext bool) {
	t.Helper()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate server adapter contract")
	}
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filepath.Join(filepath.Dir(testFile), "server.go"), nil, 0)
	if err != nil {
		t.Fatalf("parse server adapter source: %v", err)
	}
	var function *ast.FuncDecl
	for _, declaration := range parsed.Decls {
		candidate, ok := declaration.(*ast.FuncDecl)
		if !ok || candidate.Recv == nil || candidate.Name.Name != method {
			continue
		}
		receiver, receiverOK := candidate.Recv.List[0].Type.(*ast.Ident)
		if receiverOK && receiver.Name == "systemServer" {
			function = candidate
			break
		}
	}
	wantParameters := 1
	if withContext {
		wantParameters = 2
	}
	if function == nil || function.Type.Params == nil || len(function.Type.Params.List) != wantParameters ||
		function.Type.Results == nil || len(function.Type.Results.List) != 1 ||
		function.Body == nil || len(function.Body.List) != 1 {
		t.Fatalf("system server %s signature or single-return body changed", method)
	}
	returned, ok := function.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(returned.Results) != 1 {
		t.Fatalf("system server %s no longer directly returns one operation", method)
	}
	call, ok := returned.Results[0].(*ast.CallExpr)
	if !ok {
		t.Fatalf("system server %s no longer directly returns one method call", method)
	}
	selector, selectorOK := call.Fun.(*ast.SelectorExpr)
	var supplied *ast.Ident
	suppliedOK := false
	if selectorOK {
		supplied, suppliedOK = selector.X.(*ast.Ident)
	}
	wantArguments := 0
	if withContext {
		wantArguments = 1
	}
	if !selectorOK || !suppliedOK || supplied.Name != "supplied" || selector.Sel.Name != method ||
		len(call.Args) != wantArguments {
		t.Fatalf("system server %s no longer delegates directly to the exact supplied server", method)
	}
	if withContext {
		ctx, ok := call.Args[0].(*ast.Ident)
		if !ok || ctx.Name != "ctx" {
			t.Fatal("system server Shutdown no longer passes the exact supplied context")
		}
	}
}
