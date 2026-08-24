package template

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/internal/testutil"
	"github.com/devdimensionlab/plybuild/pkg/config"
)

const templateCopyMutation = "6. template copy keeps source before destination"

type recordedTemplateCopyCall struct {
	Source      string
	Destination string
}

type recordingTemplateCopy struct {
	calls []recordedTemplateCopyCall
	err   error
}

func (recording *recordingTemplateCopy) dependencies() templateCopyDependencies {
	return templateCopyDependencies{Files: recording}
}

func (recording *recordingTemplateCopy) CopyOrMerge(source string, destination string) error {
	recording.calls = append(recording.calls, recordedTemplateCopyCall{Source: source, Destination: destination})
	return recording.err
}

func (recording *recordingTemplateCopy) assertedCalls() ([]recordedTemplateCopyCall, error) {
	if len(recording.calls) == 0 {
		return nil, errors.New("recorded template copy population is empty")
	}
	return recording.calls, nil
}

func TestTemplateCopyKeepsCompleteSourceBeforeCompleteResolvedDestination(t *testing.T) {
	if templateCopyMutation == "" {
		t.Fatal("template-copy mutation label is empty")
	}
	sourceProject, targetProject, source, destination := templateCopyProjects(t)
	recording := &recordingTemplateCopy{}

	err := merge(recording.dependencies(), sourceProject, targetProject, false)

	if err != nil {
		t.Fatalf("template merge returned an error: %v", err)
	}
	assertRecordedTemplateCopyCalls(t, recording, []recordedTemplateCopyCall{{
		Source: source, Destination: destination,
	}})
}

func TestTemplateCopyReturnsDependencyErrorAfterCompleteDelivery(t *testing.T) {
	sourceProject, targetProject, source, destination := templateCopyProjects(t)
	sentinel := errors.New("complete template copy dependency error")
	recording := &recordingTemplateCopy{err: sentinel}

	err := merge(recording.dependencies(), sourceProject, targetProject, false)

	if !errors.Is(err, sentinel) {
		t.Fatalf("template copy dependency error was %v, want %v", err, sentinel)
	}
	assertRecordedTemplateCopyCalls(t, recording, []recordedTemplateCopyCall{{
		Source: source, Destination: destination,
	}})
}

func TestTemplateCopyDependenciesDefaultToNoMutation(t *testing.T) {
	sourceProject, targetProject, _, destination := templateCopyProjects(t)

	err := merge(templateCopyDependencies{}, sourceProject, targetProject, false)

	if !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe template copy dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
	if _, statErr := os.Stat(destination); !os.IsNotExist(statErr) {
		t.Fatalf("safe template copy dependency default mutated the destination: %v", statErr)
	}
}

func TestRecordedTemplateCopyCallsRejectEmptyPopulation(t *testing.T) {
	recording := &recordingTemplateCopy{}

	if _, err := recording.assertedCalls(); err == nil {
		t.Fatal("empty recorded template copy population passed")
	}
}

func templateCopyProjects(t *testing.T) (config.Project, config.Project, string, string) {
	t.Helper()
	sourceDirectory := t.TempDir()
	relative := "complete source file.txt"
	source := filepath.Join(sourceDirectory, relative)
	if err := testutil.WriteFileOutsideWorkingTree(source, []byte("complete template bytes"), 0644); err != nil {
		t.Fatalf("create template copy fixture: %v", err)
	}
	targetDirectory := filepath.Join(t.TempDir(), "complete resolved target project")
	return config.Project{Path: sourceDirectory}, config.Project{Path: targetDirectory},
		source, filepath.Join(targetDirectory, relative)
}

func assertRecordedTemplateCopyCalls(t *testing.T, recording *recordingTemplateCopy, want []recordedTemplateCopyCall) {
	t.Helper()
	calls, err := recording.assertedCalls()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("recorded template copy calls differ:\n got: %#v\nwant: %#v", calls, want)
	}
}
