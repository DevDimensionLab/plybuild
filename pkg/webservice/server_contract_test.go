package webservice

import (
	"context"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	serveradapter "github.com/devdimensionlab/plybuild/internal/adapter/server"
	"github.com/devdimensionlab/plybuild/pkg/webservice/api"
)

type recordedWebServerShutdown struct {
	server  *http.Server
	context context.Context
}

type recordingWebServerOperations struct {
	listenServers    []*http.Server
	shutdownRequests []recordedWebServerShutdown
	listenErr        error
	shutdownErr      error
	sequence         *[]string
}

func (recording *recordingWebServerOperations) startDependencies(
	dependencies startWebServerDependencies,
) startWebServerDependencies {
	dependencies.ServerOperations.Operations = recording
	return dependencies
}

func (recording *recordingWebServerOperations) stopDependencies(
	dependencies stopWebServerDependencies,
) stopWebServerDependencies {
	dependencies.ServerOperations.Operations = recording
	return dependencies
}

func (recording *recordingWebServerOperations) ListenAndServe(server *http.Server) error {
	recording.listenServers = append(recording.listenServers, server)
	if recording.sequence != nil {
		*recording.sequence = append(*recording.sequence, "listen")
	}
	return recording.listenErr
}

func (recording *recordingWebServerOperations) Shutdown(server *http.Server, ctx context.Context) error {
	recording.shutdownRequests = append(recording.shutdownRequests, recordedWebServerShutdown{
		server: server, context: ctx,
	})
	if recording.sequence != nil {
		*recording.sequence = append(*recording.sequence, "shutdown")
	}
	return recording.shutdownErr
}

func (recording *recordingWebServerOperations) assertedListenServers() ([]*http.Server, error) {
	if len(recording.listenServers) == 0 {
		return nil, errors.New("recorded webservice listen population is empty")
	}
	return recording.listenServers, nil
}

func (recording *recordingWebServerOperations) assertedShutdownRequests() ([]recordedWebServerShutdown, error) {
	if len(recording.shutdownRequests) == 0 {
		return nil, errors.New("recorded webservice shutdown population is empty")
	}
	return recording.shutdownRequests, nil
}

type recordedHandlerRegistration struct {
	path    string
	handler func(http.ResponseWriter, *http.Request)
}

type recordingWebServerHandlers struct {
	registrations []recordedHandlerRegistration
	sequence      *[]string
}

func (recording *recordingWebServerHandlers) dependencies(
	dependencies startWebServerDependencies,
) startWebServerDependencies {
	dependencies.HandleFunc = recording.HandleFunc
	return dependencies
}

func (recording *recordingWebServerHandlers) HandleFunc(
	path string,
	handler func(http.ResponseWriter, *http.Request),
) {
	recording.registrations = append(recording.registrations, recordedHandlerRegistration{
		path: path, handler: handler,
	})
	if recording.sequence != nil {
		*recording.sequence = append(*recording.sequence, "handle "+path)
	}
}

func (recording *recordingWebServerHandlers) assertedRegistrations() ([]recordedHandlerRegistration, error) {
	if len(recording.registrations) == 0 {
		return nil, errors.New("recorded webservice handler population is empty")
	}
	return recording.registrations, nil
}

type recordingWebServerLogger struct {
	arguments [][]interface{}
	sequence  *[]string
}

func (recording *recordingWebServerLogger) dependencies(
	dependencies startWebServerDependencies,
) startWebServerDependencies {
	dependencies.Print = recording.Print
	return dependencies
}

func (recording *recordingWebServerLogger) Print(arguments ...interface{}) {
	recording.arguments = append(recording.arguments, append([]interface{}(nil), arguments...))
	if recording.sequence != nil {
		*recording.sequence = append(*recording.sequence, "log")
	}
}

func (recording *recordingWebServerLogger) assertedArguments() ([][]interface{}, error) {
	if len(recording.arguments) == 0 {
		return nil, errors.New("recorded webservice log population is empty")
	}
	return recording.arguments, nil
}

type recordingWebServerBackground struct {
	context  context.Context
	attempts int
	sequence *[]string
}

func (recording *recordingWebServerBackground) dependencies(
	dependencies stopWebServerDependencies,
) stopWebServerDependencies {
	dependencies.Background = recording.Background
	return dependencies
}

func (recording *recordingWebServerBackground) Background() context.Context {
	recording.attempts++
	if recording.sequence != nil {
		*recording.sequence = append(*recording.sequence, "background")
	}
	return recording.context
}

func (recording *recordingWebServerBackground) assertedAttempts() (int, error) {
	if recording.attempts == 0 {
		return 0, errors.New("recorded webservice background population is empty")
	}
	return recording.attempts, nil
}

type recordingWebServerTimeout struct {
	parents     []context.Context
	durations   []time.Duration
	context     context.Context
	cancelCalls int
	sequence    *[]string
}

func (recording *recordingWebServerTimeout) dependencies(
	dependencies stopWebServerDependencies,
) stopWebServerDependencies {
	dependencies.WithTimeout = recording.WithTimeout
	return dependencies
}

func (recording *recordingWebServerTimeout) WithTimeout(
	parent context.Context,
	duration time.Duration,
) (context.Context, context.CancelFunc) {
	recording.parents = append(recording.parents, parent)
	recording.durations = append(recording.durations, duration)
	if recording.sequence != nil {
		*recording.sequence = append(*recording.sequence, "timeout")
	}
	return recording.context, func() {
		recording.cancelCalls++
		if recording.sequence != nil {
			*recording.sequence = append(*recording.sequence, "cancel")
		}
	}
}

func (recording *recordingWebServerTimeout) assertedRequests() ([]context.Context, []time.Duration, error) {
	if len(recording.parents) == 0 || len(recording.durations) == 0 {
		return nil, nil, errors.New("recorded webservice timeout population is empty")
	}
	return recording.parents, recording.durations, nil
}

func TestInteractiveWebServerBindsOnlyToExactIPv4Loopback(t *testing.T) {
	if port != 7999 {
		t.Fatalf("interactive web server port is %d, want exactly 7999", port)
	}
	if server == nil || server.Addr != "127.0.0.1:7999" {
		t.Fatalf("interactive package-level server is %#v, want exact IPv4 loopback address 127.0.0.1:7999", server)
	}
	if selectWebServer() != server {
		t.Fatalf("interactive server selector returned %p, want exact package-level server %p", selectWebServer(), server)
	}
}

func TestStartWebServerSelectsExactProductionDependenciesAndPublicComposition(t *testing.T) {
	dependencies := systemStartWebServerDependencies()
	systemOperations := serveradapter.System()

	if server == nil || server.Addr != "127.0.0.1:7999" {
		t.Fatalf("StartWebServer package-level server is %#v, want exact 127.0.0.1:7999 address", server)
	}
	assertSameFunction(t, dependencies.Server, serveradapter.Selector(selectWebServer), "StartWebServer server selector")
	if dependencies.Server() != server {
		t.Fatalf("StartWebServer selected server %p, want exact package-level server %p", dependencies.Server(), server)
	}
	if dependencies.ServerOperations.Operations == nil ||
		reflect.TypeOf(dependencies.ServerOperations.Operations) != reflect.TypeOf(systemOperations.Operations) {
		t.Fatalf("StartWebServer operations are %#v, want exact system dependency %#v",
			dependencies.ServerOperations, systemOperations)
	}
	assertSameFunction(t, dependencies.HandleFunc, http.HandleFunc, "StartWebServer handler registrar")
	assertSameFunction(t, dependencies.Print, log.Print, "StartWebServer logger")
	assertPublicWebServerComposition(t, "StartWebServer", "startWebServer", "systemStartWebServerDependencies")
}

func TestStartWebServerRecordingDoublesPreserveCompleteCallerDependencies(t *testing.T) {
	callerOperations := &recordingWebServerOperations{}
	callerServer := &http.Server{Addr: "complete caller-owned start server"}
	callerServerSelector := func() *http.Server { return callerServer }
	callerHandle := func(string, func(http.ResponseWriter, *http.Request)) {}
	callerPrint := func(...interface{}) {}
	initial := startWebServerDependencies{
		ServerOperations: serveradapter.Dependencies{Operations: callerOperations},
		Server:           callerServerSelector,
		HandleFunc:       callerHandle,
		Print:            callerPrint,
	}
	serverRecording := &recordingWebServerOperations{}
	handlerRecording := &recordingWebServerHandlers{}
	loggerRecording := &recordingWebServerLogger{}

	withServer := serverRecording.startDependencies(initial)
	withHandler := handlerRecording.dependencies(initial)
	withLogger := loggerRecording.dependencies(initial)

	if withServer.ServerOperations.Operations != serverRecording {
		t.Fatalf("start server recorder lost the caller-owned operations dependency: %#v", withServer)
	}
	assertSameFunction(t, withServer.Server, serveradapter.Selector(callerServerSelector), "preserved start server selector")
	assertSameFunction(t, withServer.HandleFunc, callerHandle, "preserved start handler registrar")
	assertSameFunction(t, withServer.Print, callerPrint, "preserved start logger")
	if withHandler.ServerOperations.Operations != callerOperations {
		t.Fatalf("start handler recorder lost the caller-owned operations dependency: %#v", withHandler)
	}
	assertSameFunction(t, withHandler.Server, serveradapter.Selector(callerServerSelector), "handler-preserved start server selector")
	assertSameFunction(t, withHandler.HandleFunc, handlerRecording.HandleFunc, "recording start handler registrar")
	assertSameFunction(t, withHandler.Print, callerPrint, "handler-preserved start logger")
	if withLogger.ServerOperations.Operations != callerOperations {
		t.Fatalf("start logger recorder lost the caller-owned operations dependency: %#v", withLogger)
	}
	assertSameFunction(t, withLogger.Server, serveradapter.Selector(callerServerSelector), "logger-preserved start server selector")
	assertSameFunction(t, withLogger.HandleFunc, callerHandle, "logger-preserved start handler registrar")
	assertSameFunction(t, withLogger.Print, loggerRecording.Print, "recording start logger")
}

func TestStartWebServerRegistersExactHandlersBeforeOneListenAndLogsOnlyNonNilError(t *testing.T) {
	listenError := errors.New("complete webservice listen dependency error")
	tests := []struct {
		name string
		err  error
	}{
		{name: "nil error prints nothing"},
		{name: "non-nil error is printed exactly", err: listenError},
	}

	if len(tests) == 0 {
		t.Fatal("TestStartWebServerRegistersExactHandlersBeforeOneListenAndLogsOnlyNonNilError test-case population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sequence := []string{}
			supplied := &http.Server{Addr: "arbitrary recorded start server identity"}
			serverSelections := 0
			selectServer := func() *http.Server {
				serverSelections++
				return supplied
			}
			operations := &recordingWebServerOperations{listenErr: test.err, sequence: &sequence}
			handlers := &recordingWebServerHandlers{sequence: &sequence}
			logger := &recordingWebServerLogger{sequence: &sequence}
			dependencies := logger.dependencies(handlers.dependencies(operations.startDependencies(
				startWebServerDependencies{Server: selectServer},
			)))

			startWebServer(dependencies)

			registrations, populationErr := handlers.assertedRegistrations()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			assertExactWebServerRegistrations(t, registrations)
			servers, populationErr := operations.assertedListenServers()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if serverSelections != 1 || len(servers) != 1 || servers[0] != supplied {
				t.Fatalf("StartWebServer selected %d times and listened on %#v, want one exact server identity %p",
					serverSelections, servers, supplied)
			}
			if len(operations.shutdownRequests) != 0 {
				t.Fatalf("StartWebServer invoked unrelated shutdown operations: %#v", operations.shutdownRequests)
			}
			wantSequence := []string{
				"handle /ui/generate",
				"handle /api/generate",
				"handle /ui/upgrade",
				"handle /api/upgrade",
				"listen",
			}
			if test.err != nil {
				wantSequence = append(wantSequence, "log")
			}
			if !reflect.DeepEqual(sequence, wantSequence) {
				t.Fatalf("StartWebServer effect order was %#v, want %#v", sequence, wantSequence)
			}
			if test.err == nil {
				if len(logger.arguments) != 0 {
					t.Fatalf("StartWebServer logged nil-listen success arguments: %#v", logger.arguments)
				}
				return
			}
			arguments, populationErr := logger.assertedArguments()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if len(arguments) != 1 || len(arguments[0]) != 1 || arguments[0][0] != listenError {
				t.Fatalf("StartWebServer log arguments were %#v, want one exact error %v", arguments, listenError)
			}
		})
	}
}

func TestStopWebServerSelectsExactProductionDependenciesAndPublicComposition(t *testing.T) {
	dependencies := systemStopWebServerDependencies()
	systemOperations := serveradapter.System()

	if server == nil || server.Addr != "127.0.0.1:7999" {
		t.Fatalf("StopWebServer package-level server is %#v, want exact 127.0.0.1:7999 address", server)
	}
	assertSameFunction(t, dependencies.Server, serveradapter.Selector(selectWebServer), "StopWebServer server selector")
	if dependencies.Server() != server {
		t.Fatalf("StopWebServer selected server %p, want exact package-level server %p", dependencies.Server(), server)
	}
	if dependencies.ServerOperations.Operations == nil ||
		reflect.TypeOf(dependencies.ServerOperations.Operations) != reflect.TypeOf(systemOperations.Operations) {
		t.Fatalf("StopWebServer operations are %#v, want exact system dependency %#v",
			dependencies.ServerOperations, systemOperations)
	}
	assertSameFunction(t, dependencies.Background, context.Background, "StopWebServer background selector")
	assertSameFunction(t, dependencies.WithTimeout, context.WithTimeout, "StopWebServer timeout selector")
	assertPublicWebServerComposition(t, "StopWebServer", "stopWebServer", "systemStopWebServerDependencies")
}

func TestStopWebServerRecordingDoublesPreserveCompleteCallerDependencies(t *testing.T) {
	callerOperations := &recordingWebServerOperations{}
	callerServer := &http.Server{Addr: "complete caller-owned stop server"}
	callerServerSelector := func() *http.Server { return callerServer }
	callerBackground := func() context.Context { return context.TODO() }
	callerTimeout := func(parent context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
		return parent, func() {}
	}
	initial := stopWebServerDependencies{
		ServerOperations: serveradapter.Dependencies{Operations: callerOperations},
		Server:           callerServerSelector,
		Background:       callerBackground,
		WithTimeout:      callerTimeout,
	}
	serverRecording := &recordingWebServerOperations{}
	backgroundRecording := &recordingWebServerBackground{}
	timeoutRecording := &recordingWebServerTimeout{}

	withServer := serverRecording.stopDependencies(initial)
	withBackground := backgroundRecording.dependencies(initial)
	withTimeout := timeoutRecording.dependencies(initial)

	if withServer.ServerOperations.Operations != serverRecording {
		t.Fatalf("stop server recorder lost the caller-owned operations dependency: %#v", withServer)
	}
	assertSameFunction(t, withServer.Server, serveradapter.Selector(callerServerSelector), "preserved stop server selector")
	assertSameFunction(t, withServer.Background, callerBackground, "server-preserved stop background")
	assertSameFunction(t, withServer.WithTimeout, callerTimeout, "server-preserved stop timeout")
	if withBackground.ServerOperations.Operations != callerOperations {
		t.Fatalf("stop background recorder lost the caller-owned operations dependency: %#v", withBackground)
	}
	assertSameFunction(t, withBackground.Server, serveradapter.Selector(callerServerSelector), "background-preserved stop server selector")
	assertSameFunction(t, withBackground.Background, backgroundRecording.Background, "recording stop background")
	assertSameFunction(t, withBackground.WithTimeout, callerTimeout, "background-preserved stop timeout")
	if withTimeout.ServerOperations.Operations != callerOperations {
		t.Fatalf("stop timeout recorder lost the caller-owned operations dependency: %#v", withTimeout)
	}
	assertSameFunction(t, withTimeout.Server, serveradapter.Selector(callerServerSelector), "timeout-preserved stop server selector")
	assertSameFunction(t, withTimeout.Background, callerBackground, "timeout-preserved stop background")
	assertSameFunction(t, withTimeout.WithTimeout, timeoutRecording.WithTimeout, "recording stop timeout")
}

func TestStopWebServerUsesExactBackgroundFiveSecondContextOnceDefersCancelAndIgnoresError(t *testing.T) {
	sequence := []string{}
	supplied := &http.Server{Addr: "arbitrary recorded stop server identity"}
	serverSelections := 0
	selectServer := func() *http.Server {
		serverSelections++
		return supplied
	}
	parent := context.WithValue(context.Background(), struct{ label string }{"parent"}, "complete")
	shutdownContext := context.WithValue(parent, struct{ label string }{"child"}, "exact identity")
	background := &recordingWebServerBackground{context: parent, sequence: &sequence}
	timeout := &recordingWebServerTimeout{context: shutdownContext, sequence: &sequence}
	shutdownError := errors.New("complete ignored webservice shutdown dependency error")
	operations := &recordingWebServerOperations{shutdownErr: shutdownError, sequence: &sequence}
	dependencies := timeout.dependencies(background.dependencies(operations.stopDependencies(
		stopWebServerDependencies{Server: selectServer},
	)))

	stopWebServer(dependencies)

	attempts, populationErr := background.assertedAttempts()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if attempts != 1 {
		t.Fatalf("StopWebServer requested a background context %d times, want exactly 1", attempts)
	}
	parents, durations, populationErr := timeout.assertedRequests()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if len(parents) != 1 || parents[0] != parent ||
		!reflect.DeepEqual(durations, []time.Duration{5 * time.Second}) {
		t.Fatalf("StopWebServer timeout requests were parents=%#v durations=%#v, want exact parent and five seconds",
			parents, durations)
	}
	requests, populationErr := operations.assertedShutdownRequests()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if serverSelections != 1 || len(requests) != 1 || requests[0].server != supplied || requests[0].context != shutdownContext {
		t.Fatalf("StopWebServer selected %d times and recorded %#v, want one exact server %p and context %#v",
			serverSelections, requests, supplied, shutdownContext)
	}
	if len(operations.listenServers) != 0 {
		t.Fatalf("StopWebServer invoked unrelated listen operations: %#v", operations.listenServers)
	}
	if timeout.cancelCalls != 1 {
		t.Fatalf("StopWebServer cancellation ran %d times, want exactly 1", timeout.cancelCalls)
	}
	wantSequence := []string{"background", "timeout", "shutdown", "cancel"}
	if !reflect.DeepEqual(sequence, wantSequence) {
		t.Fatalf("StopWebServer effect order was %#v, want %#v", sequence, wantSequence)
	}
}

func TestWebServerDependenciesDefaultToSafeNoOperation(t *testing.T) {
	startWebServer(startWebServerDependencies{})
	stopWebServer(stopWebServerDependencies{})
}

func TestRecordedWebServerDependenciesRejectEmptyPopulations(t *testing.T) {
	operations := &recordingWebServerOperations{}
	handlers := &recordingWebServerHandlers{}
	logger := &recordingWebServerLogger{}
	background := &recordingWebServerBackground{}
	timeout := &recordingWebServerTimeout{}

	if _, err := operations.assertedListenServers(); err == nil {
		t.Fatal("empty recorded webservice listen population passed")
	}
	if _, err := operations.assertedShutdownRequests(); err == nil {
		t.Fatal("empty recorded webservice shutdown population passed")
	}
	if _, err := handlers.assertedRegistrations(); err == nil {
		t.Fatal("empty recorded webservice handler population passed")
	}
	if _, err := logger.assertedArguments(); err == nil {
		t.Fatal("empty recorded webservice log population passed")
	}
	if _, err := background.assertedAttempts(); err == nil {
		t.Fatal("empty recorded webservice background population passed")
	}
	if _, _, err := timeout.assertedRequests(); err == nil {
		t.Fatal("empty recorded webservice timeout population passed")
	}
}

func assertExactWebServerRegistrations(t *testing.T, registrations []recordedHandlerRegistration) {
	t.Helper()
	wantPaths := []string{"/ui/generate", "/api/generate", "/ui/upgrade", "/api/upgrade"}
	wantHandlers := []func(http.ResponseWriter, *http.Request){
		api.GetGenerate,
		api.PostGenerate,
		api.GetUpgrade,
		api.PostUpgrade,
	}
	if len(registrations) != len(wantPaths) {
		t.Fatalf("StartWebServer registered %d handlers, want exactly %d", len(registrations), len(wantPaths))
	}
	for index, registration := range registrations {
		if registration.path != wantPaths[index] ||
			reflect.ValueOf(registration.handler).Pointer() != reflect.ValueOf(wantHandlers[index]).Pointer() {
			t.Fatalf("StartWebServer registration %d was path=%q handler=%p, want path=%q handler=%p",
				index, registration.path, registration.handler, wantPaths[index], wantHandlers[index])
		}
	}
}

func assertSameFunction(t *testing.T, actual, expected interface{}, label string) {
	t.Helper()
	if actual == nil || expected == nil || reflect.ValueOf(actual).Pointer() != reflect.ValueOf(expected).Pointer() {
		t.Fatalf("%s is %p, want exact function %p", label, actual, expected)
	}
}

func assertPublicWebServerComposition(t *testing.T, public, private, selector string) {
	t.Helper()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate webservice server contract")
	}
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filepath.Join(filepath.Dir(testFile), "api.go"), nil, 0)
	if err != nil {
		t.Fatalf("parse webservice source: %v", err)
	}
	var function *ast.FuncDecl
	for _, declaration := range parsed.Decls {
		candidate, ok := declaration.(*ast.FuncDecl)
		if ok && candidate.Recv == nil && candidate.Name.Name == public {
			function = candidate
			break
		}
	}
	if function == nil || function.Type.Params == nil || len(function.Type.Params.List) != 0 ||
		function.Type.Results != nil || function.Body == nil || len(function.Body.List) != 1 {
		t.Fatalf("public %s signature or single-call composition changed", public)
	}
	statement, ok := function.Body.List[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("public %s no longer directly invokes one private composition", public)
	}
	call, ok := statement.X.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		t.Fatalf("public %s no longer calls %s with one dependency value", public, private)
	}
	callee, calleeOK := call.Fun.(*ast.Ident)
	selectorCall, selectorOK := call.Args[0].(*ast.CallExpr)
	if !calleeOK || callee.Name != private || !selectorOK || len(selectorCall.Args) != 0 {
		t.Fatalf("public %s no longer composes exact private function %s", public, private)
	}
	selectorName, selectorNameOK := selectorCall.Fun.(*ast.Ident)
	if !selectorNameOK || selectorName.Name != selector {
		t.Fatalf("public %s no longer selects exact production dependencies through %s", public, selector)
	}
}
