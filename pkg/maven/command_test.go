package maven

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/sirupsen/logrus"
)

const mavenProcessMutation = "2. Maven command keeps executable before arguments"

type recordingMavenProcess struct {
	commands []process.Command
	err      error
	sequence *[]string
}

func (recording *recordingMavenProcess) dependencies(dependencies process.Dependencies) process.Dependencies {
	dependencies.Runner = recording
	return dependencies
}

func (recording *recordingMavenProcess) Run(command process.Command) error {
	if recording.sequence != nil {
		*recording.sequence = append(*recording.sequence, "process")
	}
	command.Args = append([]string(nil), command.Args...)
	recording.commands = append(recording.commands, command)
	return recording.err
}

type recordingMavenLogHook struct {
	entries  []*logrus.Entry
	sequence *[]string
}

func (hook *recordingMavenLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *recordingMavenLogHook) Fire(entry *logrus.Entry) error {
	if hook.sequence != nil {
		*hook.sequence = append(*hook.sequence, "log")
	}
	copy := *entry
	hook.entries = append(hook.entries, &copy)
	return nil
}

func (recording *recordingMavenProcess) assertedCommands() ([]process.Command, error) {
	if len(recording.commands) == 0 {
		return nil, errors.New("recorded Maven process population is empty")
	}
	return recording.commands, nil
}

func TestMavenCommandKeepsExecutableBeforeArguments(t *testing.T) {
	if mavenProcessMutation == "" {
		t.Fatal("maven-process mutation label is empty")
	}
	tests := []struct {
		name       string
		level      logrus.Level
		wantStdout bool
	}{
		{name: "panic level", level: logrus.PanicLevel, wantStdout: false},
		{name: "fatal level", level: logrus.FatalLevel, wantStdout: false},
		{name: "error level", level: logrus.ErrorLevel, wantStdout: false},
		{name: "warning level", level: logrus.WarnLevel, wantStdout: false},
		{name: "info level", level: logrus.InfoLevel, wantStdout: false},
		{name: "debug level", level: logrus.DebugLevel, wantStdout: true},
		{name: "trace level", level: logrus.TraceLevel, wantStdout: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setGlobalLogrusLevel(t, test.level)
			injectedStdout := &bytes.Buffer{}
			sentinel := errors.New("complete Maven process dependency result")
			recording := &recordingMavenProcess{err: sentinel}
			dependencies := recording.dependencies(process.Dependencies{Stdout: injectedStdout})
			if dependencies.Runner != recording || dependencies.Stdout != injectedStdout {
				t.Fatalf("Maven process dependency lost its complete value: %#v", dependencies)
			}
			project := config.Project{Path: filepath.Join(t.TempDir(), "project with spaces")}
			arguments := []string{
				"org.codehaus.mojo:versions-maven-plugin:2.8.1:update-property",
				"-Dproperty=complete.dependency.value",
				"-DnewVersion=[1.2.3-complete value]",
				"-DallowDowngrade=true",
			}

			callback := runOn(dependencies, "mvn", arguments...)
			err := callback(Repository{Url: "https://unused.example.invalid/complete"}, project)

			if err != sentinel {
				t.Fatalf("Maven callback returned %v, want exact error %v", err, sentinel)
			}
			var wantStdout io.Writer
			if test.wantStdout {
				wantStdout = injectedStdout
			}
			assertRecordedMavenCommands(t, recording, []process.Command{{
				Name:   "mvn",
				Args:   arguments,
				Dir:    project.Path,
				Stdin:  nil,
				Stdout: wantStdout,
				Stderr: nil,
				Start:  false,
			}})
			if len(recording.commands) != 1 {
				t.Fatalf("Maven callback made %d process requests, want exactly 1", len(recording.commands))
			}
			if test.wantStdout && recording.commands[0].Stdout != injectedStdout {
				t.Fatal("enabled injected stdout writer did not reach the process dependency unchanged")
			}
		})
	}
}

func TestMavenDebugSystemDependencyForwardsExactOSStdoutWithoutLaunchingProcess(t *testing.T) {
	setGlobalLogrusLevel(t, logrus.DebugLevel)
	recording := &recordingMavenProcess{}
	dependencies := systemRunOnDependencies()
	runnerOnly := process.SystemRunner()
	standardOutput := process.SystemStdout()
	if dependencies.Runner == nil || dependencies.Runner != runnerOnly.Runner {
		t.Fatalf("Maven production process runner is %T, want exact runner %T",
			dependencies.Runner, runnerOnly.Runner)
	}
	if runnerOnly.Stdout != nil {
		t.Fatalf("runner-only system dependency has standard output %T, want nil", runnerOnly.Stdout)
	}
	if dependencies.Stdout == nil || dependencies.Stdout != standardOutput || dependencies.Stdout != os.Stdout {
		t.Fatalf("Maven production standard output is %T, want exact system os.Stdout", dependencies.Stdout)
	}
	assertPublicRunOnComposition(t)
	dependencies = recording.dependencies(dependencies)
	if dependencies.Runner != recording || dependencies.Stdout != standardOutput {
		t.Fatalf("Maven recorder lost the complete production dependency: %#v", dependencies)
	}

	err := runOn(dependencies, "mvn", "complete argument")(Repository{}, config.Project{Path: t.TempDir()})

	if err != nil {
		t.Fatalf("Maven system-composition callback returned an error: %v", err)
	}
	commands, populationErr := recording.assertedCommands()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if len(commands) != 1 {
		t.Fatalf("Maven system-composition callback made %d process requests, want exactly 1", len(commands))
	}
	if commands[0].Stdout != os.Stdout {
		t.Fatalf("Maven system-composition stdout is %T, want exact os.Stdout", commands[0].Stdout)
	}
}

func assertPublicRunOnComposition(t *testing.T) {
	t.Helper()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate Maven RunOn contract")
	}
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filepath.Join(filepath.Dir(testFile), "command.go"), nil, 0)
	if err != nil {
		t.Fatalf("parse Maven command source: %v", err)
	}
	var function *ast.FuncDecl
	for _, declaration := range parsed.Decls {
		candidate, ok := declaration.(*ast.FuncDecl)
		if ok && candidate.Recv == nil && candidate.Name.Name == "RunOn" {
			function = candidate
			break
		}
	}
	if function == nil || function.Type.Params == nil || len(function.Type.Params.List) != 2 ||
		function.Type.Results == nil || len(function.Type.Results.List) != 1 ||
		function.Body == nil || len(function.Body.List) != 1 {
		t.Fatal("public RunOn signature or single-return body changed")
	}
	commandParameter := function.Type.Params.List[0]
	argumentsParameter := function.Type.Params.List[1]
	commandType, commandTypeOK := commandParameter.Type.(*ast.Ident)
	variadicType, variadicTypeOK := argumentsParameter.Type.(*ast.Ellipsis)
	var variadicElement *ast.Ident
	variadicElementOK := false
	if variadicTypeOK {
		variadicElement, variadicElementOK = variadicType.Elt.(*ast.Ident)
	}
	callbackType, callbackTypeOK := function.Type.Results.List[0].Type.(*ast.FuncType)
	if len(commandParameter.Names) != 1 || commandParameter.Names[0].Name != "cmd" ||
		!commandTypeOK || commandType.Name != "string" || len(argumentsParameter.Names) != 1 ||
		argumentsParameter.Names[0].Name != "args" || !variadicTypeOK || !variadicElementOK ||
		variadicElement.Name != "string" || !callbackTypeOK || callbackType.Params == nil ||
		len(callbackType.Params.List) != 2 || callbackType.Results == nil || len(callbackType.Results.List) != 1 {
		t.Fatal("public RunOn is no longer func RunOn(cmd string, args ...string) func(Repository, config.Project) error")
	}
	repositoryParameter := callbackType.Params.List[0]
	projectParameter := callbackType.Params.List[1]
	repositoryType, repositoryTypeOK := repositoryParameter.Type.(*ast.Ident)
	projectType, projectTypeOK := projectParameter.Type.(*ast.SelectorExpr)
	var projectPackage *ast.Ident
	projectPackageOK := false
	if projectTypeOK {
		projectPackage, projectPackageOK = projectType.X.(*ast.Ident)
	}
	callbackResult, callbackResultOK := callbackType.Results.List[0].Type.(*ast.Ident)
	if len(repositoryParameter.Names) != 1 || repositoryParameter.Names[0].Name != "repository" ||
		!repositoryTypeOK || repositoryType.Name != "Repository" || len(projectParameter.Names) != 1 ||
		projectParameter.Names[0].Name != "project" || !projectTypeOK || !projectPackageOK ||
		projectPackage.Name != "config" || projectType.Sel.Name != "Project" ||
		!callbackResultOK || callbackResult.Name != "error" {
		t.Fatal("public RunOn callback no longer preserves exact repository, project, and error shape")
	}
	returned, ok := function.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(returned.Results) != 1 {
		t.Fatal("public RunOn no longer directly returns one private command callback")
	}
	call, ok := returned.Results[0].(*ast.CallExpr)
	if !ok || len(call.Args) != 3 || call.Ellipsis == token.NoPos {
		t.Fatal("public RunOn no longer directly returns runOn with dependencies, command, and variadic arguments")
	}
	callee, calleeOK := call.Fun.(*ast.Ident)
	selector, selectorOK := call.Args[0].(*ast.CallExpr)
	if !selectorOK {
		t.Fatal("public RunOn no longer selects its private production dependencies at the call site")
	}
	selectorName, selectorNameOK := selector.Fun.(*ast.Ident)
	command, commandOK := call.Args[1].(*ast.Ident)
	arguments, argumentsOK := call.Args[2].(*ast.Ident)
	if !calleeOK || callee.Name != "runOn" || !selectorNameOK ||
		selectorName.Name != "systemRunOnDependencies" || len(selector.Args) != 0 ||
		!commandOK || command.Name != "cmd" || !argumentsOK || arguments.Name != "args" {
		t.Fatal("public RunOn no longer selects exact private production dependencies beside caller command values")
	}
}

func TestMavenProcessDependenciesDefaultToNoProcess(t *testing.T) {
	setGlobalLogrusLevel(t, logrus.DebugLevel)
	callback := runOn(process.Dependencies{}, "this-process-must-not-exist", "complete argument")

	err := callback(Repository{}, config.Project{Path: t.TempDir()})

	if err != nil {
		t.Fatalf("safe Maven process default returned an error: %v", err)
	}
}

func TestMavenLogsExactCommandBeforeProcessExecution(t *testing.T) {
	sequence := []string{}
	hook := &recordingMavenLogHook{sequence: &sequence}
	testLogger := logrus.New()
	testLogger.SetOutput(io.Discard)
	testLogger.AddHook(hook)
	previousLog := log
	log = testLogger
	t.Cleanup(func() {
		log = previousLog
	})
	recording := &recordingMavenProcess{sequence: &sequence}
	project := config.Project{Path: `/arbitrary project path//with spaces/../backslash\segment-ø`}
	arguments := []string{"first complete argument", `second\argument`}
	repository := Repository{
		Id:  "complete repository id",
		Url: `https://example.invalid/arbitrary repository/\path-ø`,
		Auth: &RepositoryAuth{
			Username:  "complete user",
			Password:  "complete password",
			Encrypted: true,
		},
	}
	wantRepository := repository
	wantRepositoryAuth := *repository.Auth
	wantRepository.Auth = &wantRepositoryAuth
	callback := runOn(recording.dependencies(process.Dependencies{}), "complete-maven-tool", arguments...)
	if len(sequence) != 0 || len(hook.entries) != 0 || len(recording.commands) != 0 {
		t.Fatalf("Maven callback construction performed effects: sequence=%#v logs=%d commands=%d",
			sequence, len(hook.entries), len(recording.commands))
	}

	err := callback(repository, project)

	if err != nil {
		t.Fatalf("Maven callback returned an error: %v", err)
	}
	if !reflect.DeepEqual(sequence, []string{"log", "process"}) {
		t.Fatalf("Maven log/process order was %#v, want log before process", sequence)
	}
	if len(hook.entries) != 1 {
		t.Fatalf("Maven emitted %d log entries, want exactly 1", len(hook.entries))
	}
	wantMessage := "running: [" + project.Path + "] => complete-maven-tool first complete argument second\\argument"
	if hook.entries[0].Level != logrus.InfoLevel || hook.entries[0].Message != wantMessage {
		t.Fatalf("Maven log entry was level=%s message=%q, want info %q",
			hook.entries[0].Level, hook.entries[0].Message, wantMessage)
	}
	commands, populationErr := recording.assertedCommands()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if len(commands) != 1 {
		t.Fatalf("Maven logging contract recorded %d process requests, want exactly 1", len(commands))
	}
	if !reflect.DeepEqual(repository, wantRepository) {
		t.Fatalf("Maven callback changed the complete caller repository value:\n got: %#v\nwant: %#v",
			repository, wantRepository)
	}
}

func setGlobalLogrusLevel(t *testing.T, level logrus.Level) {
	t.Helper()
	previousLevel := logrus.GetLevel()
	logrus.SetLevel(level)
	t.Cleanup(func() {
		logrus.SetLevel(previousLevel)
	})
}

func TestRecordedMavenCommandsRejectEmptyPopulation(t *testing.T) {
	recording := &recordingMavenProcess{}

	if _, err := recording.assertedCommands(); err == nil {
		t.Fatal("empty recorded Maven process population passed")
	}
}

func assertRecordedMavenCommands(t *testing.T, recording *recordingMavenProcess, want []process.Command) {
	t.Helper()
	commands, err := recording.assertedCommands()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(commands, want) {
		t.Fatalf("recorded Maven process commands differ:\n got: %#v\nwant: %#v", commands, want)
	}
}
