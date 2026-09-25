package clock

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

type recordingClock struct {
	now      time.Time
	attempts int
}

func (recording *recordingClock) Now() time.Time {
	recording.attempts++
	return recording.now
}

func (recording *recordingClock) assertedAttempts() (int, bool) {
	return recording.attempts, recording.attempts > 0
}

type recordingSleeper struct {
	durations []time.Duration
}

func (recording *recordingSleeper) Sleep(duration time.Duration) {
	recording.durations = append(recording.durations, duration)
}

func (recording *recordingSleeper) assertedDurations() ([]time.Duration, bool) {
	return recording.durations, len(recording.durations) > 0
}

func TestDependenciesDefaultToDeterministicZeroTime(t *testing.T) {
	got := Now(Dependencies{})

	if got != (time.Time{}) {
		t.Fatalf("zero-value clock dependencies returned %#v, want exact deterministic zero time", got)
	}
}

func TestNowReturnsExactInjectedTimeAfterOneRead(t *testing.T) {
	want := time.Date(1968, time.February, 29, 23, 59, 58, 765432100,
		time.FixedZone("arbitrary injected zone", -7*60*60-31*60))
	recording := &recordingClock{now: want}
	dependencies := Dependencies{Clock: recording}
	if dependencies.Clock != recording {
		t.Fatalf("clock dependency lost its complete value: %#v", dependencies)
	}

	got := Now(dependencies)

	if got != want {
		t.Fatalf("clock read returned %#v, want exact injected value %#v", got, want)
	}
	attempts, populated := recording.assertedAttempts()
	if !populated {
		t.Fatal("recorded clock-read population is empty")
	}
	if attempts != 1 {
		t.Fatalf("clock dependency was read %d times, want exactly 1", attempts)
	}
}

func TestSleepDependenciesDefaultToSafeNoOp(t *testing.T) {
	Sleep(Dependencies{}, time.Duration(1<<63-1))
}

func TestSleepPassesExactArbitraryDurationAfterOneRequest(t *testing.T) {
	want := -73*time.Hour + 41*time.Nanosecond
	recording := &recordingSleeper{}
	clockRecording := &recordingClock{now: time.Unix(-1, 987654321)}
	dependencies := Dependencies{Clock: clockRecording, Sleeper: recording}
	if dependencies.Clock != clockRecording || dependencies.Sleeper != recording {
		t.Fatalf("sleep dependency lost its complete clock value: %#v", dependencies)
	}

	Sleep(dependencies, want)

	durations, populated := recording.assertedDurations()
	if !populated {
		t.Fatal("recorded clock-sleep population is empty")
	}
	if !reflect.DeepEqual(durations, []time.Duration{want}) {
		t.Fatalf("clock sleep durations were %#v, want one exact duration %#v", durations, []time.Duration{want})
	}
	if clockRecording.attempts != 0 {
		t.Fatalf("clock sleep read the current time %d times, want 0", clockRecording.attempts)
	}
}

func TestSystemReturnsDirectTimeNowValueWithinBoundedObservation(t *testing.T) {
	dependencies := System()
	if dependencies.Clock == nil {
		t.Fatal("system clock dependency is empty")
	}
	if dependencies.Sleeper == nil {
		t.Fatal("system clock sleep dependency is empty")
	}
	if reflect.TypeOf(dependencies.Sleeper) != reflect.TypeOf(dependencies.Clock) {
		t.Fatalf("system clock capabilities are %T and %T, want one exact system dependency type",
			dependencies.Clock, dependencies.Sleeper)
	}
	before := time.Now()

	got := Now(dependencies)

	after := time.Now()
	if got.Before(before) || got.After(after) {
		t.Fatalf("system clock returned %s outside direct time.Now observation [%s, %s]", got, before, after)
	}
	assertSystemClockDirectlyReturnsTimeNow(t)
	assertSystemClockDirectlySleepsExactDuration(t)
}

func TestRecordedClockRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingClock{}

	if _, populated := recording.assertedAttempts(); populated {
		t.Fatal("empty recorded clock-read population passed")
	}
}

func TestRecordedSleeperRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingSleeper{}

	if _, populated := recording.assertedDurations(); populated {
		t.Fatal("empty recorded clock-sleep population passed")
	}
}

func assertSystemClockDirectlyReturnsTimeNow(t *testing.T) {
	t.Helper()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate clock adapter contract")
	}
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filepath.Join(filepath.Dir(testFile), "clock.go"), nil, 0)
	if err != nil {
		t.Fatalf("parse clock adapter source: %v", err)
	}
	var function *ast.FuncDecl
	for _, declaration := range parsed.Decls {
		candidate, ok := declaration.(*ast.FuncDecl)
		if ok && candidate.Recv != nil && candidate.Name.Name == "Now" {
			function = candidate
			break
		}
	}
	if function == nil || function.Type.Params == nil || len(function.Type.Params.List) != 0 ||
		function.Type.Results == nil || len(function.Type.Results.List) != 1 ||
		function.Body == nil || len(function.Body.List) != 1 {
		t.Fatal("system clock Now signature or single-return body changed")
	}
	result, resultOK := function.Type.Results.List[0].Type.(*ast.SelectorExpr)
	var resultPackage *ast.Ident
	resultPackageOK := false
	if resultOK {
		resultPackage, resultPackageOK = result.X.(*ast.Ident)
	}
	if !resultOK || !resultPackageOK || resultPackage.Name != "time" || result.Sel.Name != "Time" {
		t.Fatal("system clock Now no longer returns time.Time")
	}
	returned, ok := function.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(returned.Results) != 1 {
		t.Fatal("system clock Now no longer directly returns one value")
	}
	nowCall, ok := returned.Results[0].(*ast.CallExpr)
	if !ok || len(nowCall.Args) != 0 {
		t.Fatal("system clock Now no longer directly returns one zero-argument call")
	}
	selector, ok := nowCall.Fun.(*ast.SelectorExpr)
	var packageName *ast.Ident
	packageOK := false
	if ok {
		packageName, packageOK = selector.X.(*ast.Ident)
	}
	if !ok || !packageOK || packageName.Name != "time" || selector.Sel.Name != "Now" {
		t.Fatal("system clock Now no longer directly returns time.Now()")
	}
}

func assertSystemClockDirectlySleepsExactDuration(t *testing.T) {
	t.Helper()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate clock sleep adapter contract")
	}
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filepath.Join(filepath.Dir(testFile), "clock.go"), nil, 0)
	if err != nil {
		t.Fatalf("parse clock adapter source: %v", err)
	}
	var function *ast.FuncDecl
	for _, declaration := range parsed.Decls {
		candidate, ok := declaration.(*ast.FuncDecl)
		if ok && candidate.Recv != nil && candidate.Name.Name == "Sleep" {
			function = candidate
			break
		}
	}
	if function == nil || function.Type.Params == nil || len(function.Type.Params.List) != 1 ||
		function.Type.Results != nil || function.Body == nil || len(function.Body.List) != 1 {
		t.Fatal("system clock Sleep signature or single-call body changed")
	}
	receiver, receiverOK := function.Recv.List[0].Type.(*ast.Ident)
	parameter := function.Type.Params.List[0]
	durationType, durationTypeOK := parameter.Type.(*ast.SelectorExpr)
	var durationPackage *ast.Ident
	durationPackageOK := false
	if durationTypeOK {
		durationPackage, durationPackageOK = durationType.X.(*ast.Ident)
	}
	if !receiverOK || receiver.Name != "systemClock" || len(parameter.Names) != 1 ||
		parameter.Names[0].Name != "duration" || !durationTypeOK || !durationPackageOK ||
		durationPackage.Name != "time" || durationType.Sel.Name != "Duration" {
		t.Fatal("system clock Sleep no longer accepts one exact time.Duration on systemClock")
	}
	statement, ok := function.Body.List[0].(*ast.ExprStmt)
	if !ok {
		t.Fatal("system clock Sleep no longer contains one direct call")
	}
	call, ok := statement.X.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		t.Fatal("system clock Sleep no longer directly calls one operation with one duration")
	}
	selector, selectorOK := call.Fun.(*ast.SelectorExpr)
	var packageName *ast.Ident
	packageOK := false
	if selectorOK {
		packageName, packageOK = selector.X.(*ast.Ident)
	}
	duration, durationOK := call.Args[0].(*ast.Ident)
	if !selectorOK || !packageOK || packageName.Name != "time" || selector.Sel.Name != "Sleep" ||
		!durationOK || duration.Name != "duration" {
		t.Fatal("system clock Sleep no longer directly delegates the exact duration to time.Sleep")
	}
}
