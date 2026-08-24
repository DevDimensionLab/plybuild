package maven

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/logger"
	"github.com/sirupsen/logrus"
)

const mavenProcessMutation = "2. Maven command keeps executable before arguments"

type recordingMavenProcess struct {
	commands []process.Command
	err      error
}

func (recording *recordingMavenProcess) dependencies() process.Dependencies {
	return process.Dependencies{Runner: recording}
}

func (recording *recordingMavenProcess) Run(command process.Command) error {
	command.Args = append([]string(nil), command.Args...)
	recording.commands = append(recording.commands, command)
	return recording.err
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
	previousLevel := logrus.GetLevel()
	logrus.SetLevel(logrus.DebugLevel)
	t.Cleanup(func() {
		logrus.SetLevel(previousLevel)
	})
	wantStdout := logger.StdOut()
	if wantStdout == nil {
		t.Fatal("debug logger stdout writer is nil")
	}
	sentinel := errors.New("complete dependency result")
	recording := &recordingMavenProcess{err: sentinel}
	project := config.Project{Path: filepath.Join(t.TempDir(), "project with spaces")}
	arguments := []string{
		"org.codehaus.mojo:versions-maven-plugin:2.8.1:update-property",
		"-Dproperty=complete.dependency.value",
		"-DnewVersion=[1.2.3-complete value]",
		"-DallowDowngrade=true",
	}

	callback := runOn(recording.dependencies(), "mvn", arguments...)
	err := callback(Repository{}, project)

	if !errors.Is(err, sentinel) {
		t.Fatalf("Maven callback returned %v, want %v", err, sentinel)
	}
	assertRecordedMavenCommands(t, recording, []process.Command{{
		Name:   "mvn",
		Args:   arguments,
		Dir:    project.Path,
		Stdout: wantStdout,
		Stderr: nil,
	}})
	if recording.commands[0].Stdout != wantStdout {
		t.Fatal("logger stdout writer did not reach the process dependency unchanged")
	}
}

func TestMavenProcessDependenciesDefaultToNoProcess(t *testing.T) {
	callback := runOn(process.Dependencies{}, "this-process-must-not-exist", "complete argument")

	err := callback(Repository{}, config.Project{Path: t.TempDir()})

	if err != nil {
		t.Fatalf("safe Maven process default returned an error: %v", err)
	}
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
