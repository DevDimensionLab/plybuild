package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/testutil"
	"github.com/devdimensionlab/plybuild/pkg/shell"
)

const cloudCloneMutation = "3. cloud clone keeps URL before target directory"

type recordedCloudGitCall struct {
	Operation string
	Values    []string
}

type recordingCloudGit struct {
	calls       []recordedCloudGitCall
	cloneOutput shell.Output
	pullOutput  shell.Output
}

func (recording *recordingCloudGit) dependencies() refreshGitDependencies {
	return refreshGitDependencies{Git: recording}
}

func (recording *recordingCloudGit) Clone(url string, target string) shell.Output {
	recording.calls = append(recording.calls, recordedCloudGitCall{
		Operation: "clone",
		Values:    []string{url, target},
	})
	return recording.cloneOutput
}

func (recording *recordingCloudGit) Pull(target string) shell.Output {
	recording.calls = append(recording.calls, recordedCloudGitCall{
		Operation: "pull",
		Values:    []string{target},
	})
	return recording.pullOutput
}

func (recording *recordingCloudGit) assertedCalls() ([]recordedCloudGitCall, error) {
	if len(recording.calls) == 0 {
		return nil, errors.New("recorded cloud Git call population is empty")
	}
	return recording.calls, nil
}

type staticLocalConfig struct {
	configuration LocalConfiguration
	err           error
}

func localConfigWithCloudURL(url string) staticLocalConfig {
	localConfig := staticLocalConfig{}
	localConfig.configuration.CloudConfig.Git.Url = url
	return localConfig
}

func createExistingGitPath(t *testing.T, target string) {
	t.Helper()
	if err := testutil.CopyFSOutsideWorkingTree(filepath.Join(target, ".git"), os.DirFS(t.TempDir())); err != nil {
		t.Fatalf("create existing Git path: %v", err)
	}
}

func (localConfig staticLocalConfig) Implementation() DirConfig {
	return DirConfig{}
}

func (localConfig staticLocalConfig) FilePath() string {
	return ""
}

func (localConfig staticLocalConfig) CheckOrCreateConfigDir() error {
	return nil
}

func (localConfig staticLocalConfig) TouchFile() error {
	return nil
}

func (localConfig staticLocalConfig) Config() (LocalConfiguration, error) {
	return localConfig.configuration, localConfig.err
}

func (localConfig staticLocalConfig) Print() error {
	return nil
}

func (localConfig staticLocalConfig) Exists() bool {
	return true
}

func TestCloudCloneKeepsCompleteURLBeforeCompleteTargetDirectory(t *testing.T) {
	if cloudCloneMutation == "" {
		t.Fatal("cloud-clone mutation label is empty")
	}
	recording := &recordingCloudGit{}
	url := "ssh://git@example.invalid/team/complete-repository.git?ref=complete-value"
	target := filepath.Join(t.TempDir(), "complete cloud target directory")
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: target}}

	err := gitCfg.refresh(recording.dependencies(), localConfigWithCloudURL(url))

	if err != nil {
		t.Fatalf("refresh returned an error: %v", err)
	}
	assertRecordedCloudGitCalls(t, recording, []recordedCloudGitCall{{
		Operation: "clone",
		Values:    []string{url, target},
	}})
}

func TestCloudRefreshExistingGitPathSelectsPullWithCompleteDependency(t *testing.T) {
	recording := &recordingCloudGit{}
	target := filepath.Join(t.TempDir(), "existing cloud target directory")
	createExistingGitPath(t, target)
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: target}}

	err := gitCfg.refresh(recording.dependencies(), localConfigWithCloudURL("unused://clone-url"))

	if err != nil {
		t.Fatalf("refresh returned an error: %v", err)
	}
	assertRecordedCloudGitCalls(t, recording, []recordedCloudGitCall{{
		Operation: "pull",
		Values:    []string{target},
	}})
}

func TestCloudRefreshFormatsGitDependencyErrors(t *testing.T) {
	tests := []struct {
		name      string
		existing  bool
		operation string
	}{
		{name: "clone", operation: "clone"},
		{name: "pull", existing: true, operation: "pull"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			target := filepath.Join(t.TempDir(), test.name+" cloud target")
			if test.existing {
				createExistingGitPath(t, target)
			}
			output := shell.Output{Err: errors.New(test.name + " dependency failed")}
			_, _ = output.StdOut.WriteString(test.name + " complete stdout\n")
			_, _ = output.StdErr.WriteString(test.name + " complete stderr\n")
			recording := &recordingCloudGit{cloneOutput: output, pullOutput: output}
			gitCfg := GitCloudConfig{Impl: DirConfig{Path: target}}

			err := gitCfg.refresh(recording.dependencies(), localConfigWithCloudURL("https://example.invalid/complete.git"))

			if err == nil {
				t.Fatal("Git dependency error was lost")
			}
			if err.Error() != output.FormatError().Error() {
				t.Fatalf("formatted Git error changed:\n got: %q\nwant: %q", err.Error(), output.FormatError().Error())
			}
			calls, callsErr := recording.assertedCalls()
			if callsErr != nil {
				t.Fatal(callsErr)
			}
			if len(calls) != 1 || calls[0].Operation != test.operation {
				t.Fatalf("Git dependency received unexpected calls: %#v", calls)
			}
		})
	}
}

func TestCloudRefreshGitDependenciesDefaultToNoMutation(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe default cloud target")
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: target}}

	err := gitCfg.refresh(refreshGitDependencies{}, localConfigWithCloudURL("https://example.invalid/must-not-clone.git"))

	if err != nil {
		t.Fatalf("safe Git dependency default returned an error: %v", err)
	}
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("safe Git dependency default mutated the target: %v", statErr)
	}
}

func TestRecordedCloudGitCallsRejectEmptyPopulation(t *testing.T) {
	recording := &recordingCloudGit{}

	if _, err := recording.assertedCalls(); err == nil {
		t.Fatal("empty recorded cloud Git population passed")
	}
}

func assertRecordedCloudGitCalls(t *testing.T, recording *recordingCloudGit, want []recordedCloudGitCall) {
	t.Helper()
	calls, err := recording.assertedCalls()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("recorded cloud Git calls differ:\n got: %#v\nwant: %#v", calls, want)
	}
}
