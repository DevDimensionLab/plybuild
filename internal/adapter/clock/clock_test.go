package clock

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
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

func TestSystemReturnsDirectTimeNowValueWithinBoundedObservation(t *testing.T) {
	dependencies := System()
	if dependencies.Clock == nil {
		t.Fatal("system clock dependency is empty")
	}
	before := time.Now()

	got := Now(dependencies)

	after := time.Now()
	if got.Before(before) || got.After(after) {
		t.Fatalf("system clock returned %s outside direct time.Now observation [%s, %s]", got, before, after)
	}
	assertSystemClockDirectlyReturnsTimeNow(t)
}

func TestRecordedClockRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingClock{}

	if _, populated := recording.assertedAttempts(); populated {
		t.Fatal("empty recorded clock-read population passed")
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
