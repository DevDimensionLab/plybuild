package maven

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
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

func (recording *recordingMavenProcess) dependencies(stdout io.Writer) process.Dependencies {
	return process.Dependencies{Runner: recording, Stdout: stdout}
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
			dependencies := recording.dependencies(injectedStdout)
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
	dependencies := process.System()
	dependencies.Runner = recording

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

	err := runOn(recording.dependencies(nil), "complete-maven-tool", arguments...)(Repository{}, project)

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
