package shell

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

const (
	gitProcessMutation = "1. git clone keeps URL before target directory"
	gitCommitMutation  = "7. git commit keeps target directory before message"
)

type recordedProcessCall struct {
	Name string
	Args []string
	Dir  string
}

type recordingProcess struct {
	calls  []recordedProcessCall
	stdout string
	stderr string
}

func (recording *recordingProcess) dependencies() process.Dependencies {
	dependencies := process.System()
	dependencies.Runner = recording
	return dependencies
}

func (recording *recordingProcess) Run(command process.Command) error {
	recording.calls = append(recording.calls, recordedProcessCall{
		Name: command.Name,
		Args: append([]string(nil), command.Args...),
		Dir:  command.Dir,
	})
	stdout, ok := command.Stdout.(*bytes.Buffer)
	if !ok {
		return fmt.Errorf("stdout dependency is %T, want *bytes.Buffer", command.Stdout)
	}
	stderr, ok := command.Stderr.(*bytes.Buffer)
	if !ok {
		return fmt.Errorf("stderr dependency is %T, want *bytes.Buffer", command.Stderr)
	}
	_, _ = stdout.WriteString(recording.stdout)
	_, _ = stderr.WriteString(recording.stderr)
	return nil
}

func (recording *recordingProcess) assertedCalls() ([]recordedProcessCall, error) {
	if len(recording.calls) == 0 {
		return nil, errors.New("recorded process call population is empty")
	}
	return recording.calls, nil
}

func TestGitCloneKeepsURLBeforeTargetDirectory(t *testing.T) {
	if gitProcessMutation == "" {
		t.Fatal("git-process mutation label is empty")
	}
	recording := &recordingProcess{stdout: "cloned\n", stderr: "git diagnostic\n"}
	url := "ssh://git@example.invalid/team/repository.git"
	target := filepath.Join(t.TempDir(), "complete target directory")

	output := gitClone(recording.dependencies(), url, target)

	assertRecordedProcessCalls(t, recording, []recordedProcessCall{{
		Name: "git",
		Args: []string{"clone", url, target},
		Dir:  "",
	}})
	if output.StdOut.String() != recording.stdout || output.StdErr.String() != recording.stderr {
		t.Fatalf("shell output changed: stdout %q, stderr %q", output.StdOut.String(), output.StdErr.String())
	}
}

func TestGitPullAndInitKeepExactArgumentsAndWorkingDirectory(t *testing.T) {
	target := filepath.Join(t.TempDir(), "repository")
	tests := []struct {
		name string
		run  func(process.Dependencies) Output
		want recordedProcessCall
	}{
		{
			name: "pull",
			run: func(dependencies process.Dependencies) Output {
				return gitPull(dependencies, target)
			},
			want: recordedProcessCall{Name: "git", Args: []string{"-C", target, "pull", "origin"}},
		},
		{
			name: "init",
			run: func(dependencies process.Dependencies) Output {
				return gitInit(dependencies, target)
			},
			want: recordedProcessCall{Name: "git", Args: []string{"-C", target, "init"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingProcess{}
			test.run(recording.dependencies())
			assertRecordedProcessCalls(t, recording, []recordedProcessCall{test.want})
		})
	}
}

func TestGitAddAndCommitKeepsTargetBeforeCompleteMessage(t *testing.T) {
	if gitCommitMutation == "" {
		t.Fatal("git-commit mutation label is empty")
	}
	recording := &recordingProcess{}
	target := filepath.Join(t.TempDir(), "repository with spaces")
	message := "complete dependency value: com.example:library:1.2.3"

	gitAddAndCommit(recording.dependencies(), target, message)

	assertRecordedProcessCalls(t, recording, []recordedProcessCall{
		{Name: "git", Args: []string{"-C", target, "add", "."}},
		{Name: "git", Args: []string{"-C", target, "commit", "-m", fmt.Sprintf("\"%s\"", message)}},
	})
}

func TestGitProcessDependenciesDefaultToNoProcess(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe-default-repository")

	output := gitInit(process.Dependencies{}, target)

	if output.Err != nil {
		t.Fatalf("safe process default returned an error: %v", output.Err)
	}
	if _, err := os.Stat(filepath.Join(target, ".git")); !os.IsNotExist(err) {
		t.Fatalf("safe process default created git state: %v", err)
	}
}

func TestRecordedGitCallsRejectEmptyPopulation(t *testing.T) {
	recording := &recordingProcess{}

	if _, err := recording.assertedCalls(); err == nil {
		t.Fatal("empty recorded process population passed")
	}
}

func assertRecordedProcessCalls(t *testing.T, recording *recordingProcess, want []recordedProcessCall) {
	t.Helper()
	calls, err := recording.assertedCalls()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("recorded process calls differ:\n got: %#v\nwant: %#v", calls, want)
	}
}
