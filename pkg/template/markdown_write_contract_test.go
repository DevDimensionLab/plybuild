package template

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/pkg/config"
)

func (*recordingTemplateMarkdownFilesystem) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, errors.New("unexpected template-markdown archive open")
}

func (*recordingTemplateMarkdownFilesystem) Close(filesystem.File) error {
	return errors.New("unexpected template-markdown close")
}

type recordedTemplateMarkdownWrite struct {
	path string
	data []byte
	mode fs.FileMode
}

type recordingTemplateMarkdownFilesystem struct {
	writes   []recordedTemplateMarkdownWrite
	writeErr error
}

func (*recordingTemplateMarkdownFilesystem) WorkingDirectory() (string, error) {
	return "", errors.New("unexpected template-markdown working directory")
}

func (*recordingTemplateMarkdownFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected template-markdown read")
}

func (*recordingTemplateMarkdownFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, errors.New("unexpected template-markdown read directory")
}

func (*recordingTemplateMarkdownFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, errors.New("unexpected template-markdown read directory entries")
}

func (*recordingTemplateMarkdownFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected template-markdown stat")
}

func (*recordingTemplateMarkdownFilesystem) Mkdir(string, fs.FileMode) error {
	return errors.New("unexpected template-markdown mkdir")
}

func (*recordingTemplateMarkdownFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected template-markdown mkdir")
}

func (recording *recordingTemplateMarkdownFilesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	recording.writes = append(recording.writes, recordedTemplateMarkdownWrite{
		path: path,
		data: append([]byte{}, data...),
		mode: mode,
	})
	return recording.writeErr
}

func (*recordingTemplateMarkdownFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected template-markdown open file")
}

func (*recordingTemplateMarkdownFilesystem) Remove(string) error {
	return errors.New("unexpected template-markdown remove")
}

func (*recordingTemplateMarkdownFilesystem) RemoveAll(string) error {
	return errors.New("unexpected template-markdown recursive remove")
}

func (*recordingTemplateMarkdownFilesystem) Rename(string, string) error {
	return errors.New("unexpected template-markdown rename")
}

func (*recordingTemplateMarkdownFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected template-markdown glob")
}

func (*recordingTemplateMarkdownFilesystem) Walk(string, filepath.WalkFunc) error {
	return errors.New("unexpected template-markdown walk")
}

func (*recordingTemplateMarkdownFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected template-markdown create")
}

func (*recordingTemplateMarkdownFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected template-markdown copy")
}

func (recording *recordingTemplateMarkdownFilesystem) dependencies() templateMarkdownWriteDependencies {
	return templateMarkdownWriteDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingTemplateMarkdownFilesystem) assertedWrites() ([]recordedTemplateMarkdownWrite, error) {
	if len(recording.writes) == 0 {
		return nil, errors.New("recorded template-markdown write population is empty")
	}
	return recording.writes, nil
}

type recordingTemplateMarkdownDirectory struct {
	path          string
	dirCalls      int
	filePathCalls int
}

func (directory *recordingTemplateMarkdownDirectory) Dir() string {
	directory.dirCalls++
	return directory.path
}

func (directory *recordingTemplateMarkdownDirectory) FilePath(string) (string, error) {
	directory.filePathCalls++
	return "", errors.New("unexpected template-markdown FilePath call")
}

type recordingTemplateMarkdownCloudConfig struct {
	config.GitCloudConfig
	directory           config.Directory
	implementationCalls int
}

func (cloudConfig *recordingTemplateMarkdownCloudConfig) Implementation() config.Directory {
	cloudConfig.implementationCalls++
	return cloudConfig.directory
}

func TestSaveTemplateListMarkdownSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemTemplateMarkdownWriteDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("SaveTemplateListMarkdown selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("SaveTemplateListMarkdown filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestSaveTemplateListMarkdownPreservesExactReceiverPathDocumentBytesModeReturnedPathAndWriteError(t *testing.T) {
	writeError := errors.New("complete template-markdown write dependency error")
	tests := []struct {
		name             string
		markdownDocument string
		writeErr         error
	}{
		{name: "empty document succeeds", markdownDocument: ""},
		{
			name:             "non-ASCII and arbitrary bytes fail after delivery",
			markdownDocument: string([]byte{0x00, 'r', 0xc3, 0xb8, 0xff, '\n'}),
			writeErr:         writeError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := &recordingTemplateMarkdownDirectory{path: `/complete receiver root//with spaces/../backslash\segment`}
			cloudConfig := &recordingTemplateMarkdownCloudConfig{directory: directory}
			recording := &recordingTemplateMarkdownFilesystem{writeErr: test.writeErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("template-markdown dependency lost its complete filesystem value: %#v", dependencies)
			}
			wantPath := directory.path + "/" + TemplatesDir + "/README.md"

			path, err := saveTemplateListMarkdown(dependencies, cloudConfig, test.markdownDocument)

			if path != wantPath || err != test.writeErr {
				t.Fatalf("template-markdown result was (%q, %v), want exact result (%q, %v)", path, err, wantPath, test.writeErr)
			}
			writes, populationErr := recording.assertedWrites()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			wantWrites := []recordedTemplateMarkdownWrite{{
				path: wantPath,
				data: []byte(test.markdownDocument),
				mode: 0644,
			}}
			if !reflect.DeepEqual(writes, wantWrites) {
				t.Fatalf("recorded template-markdown writes differ:\n got: %#v\nwant: %#v", writes, wantWrites)
			}
			if cloudConfig.implementationCalls != 1 || directory.dirCalls != 1 || directory.filePathCalls != 0 {
				t.Fatalf("receiver path evaluation calls were Implementation=%d Dir=%d FilePath=%d, want 1, 1, 0",
					cloudConfig.implementationCalls, directory.dirCalls, directory.filePathCalls)
			}
		})
	}
}

func TestSaveTemplateListMarkdownDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	directory := &recordingTemplateMarkdownDirectory{path: "/developer/home/project/must-not-be-accessed"}
	cloudConfig := &recordingTemplateMarkdownCloudConfig{directory: directory}
	wantPath := directory.path + "/" + TemplatesDir + "/README.md"

	path, err := saveTemplateListMarkdown(templateMarkdownWriteDependencies{}, cloudConfig, "must not be written")

	if path != wantPath || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe template-markdown dependency default returned (%q, %v), want (%q, %v)", path, err, wantPath, filesystem.ErrNoFilesystem)
	}
	if cloudConfig.implementationCalls != 1 || directory.dirCalls != 1 || directory.filePathCalls != 0 {
		t.Fatalf("safe receiver path evaluation calls were Implementation=%d Dir=%d FilePath=%d, want 1, 1, 0",
			cloudConfig.implementationCalls, directory.dirCalls, directory.filePathCalls)
	}
}

func TestRecordedSaveTemplateListMarkdownRejectsEmptyWritePopulation(t *testing.T) {
	recording := &recordingTemplateMarkdownFilesystem{}

	if _, err := recording.assertedWrites(); err == nil {
		t.Fatal("empty recorded template-markdown write population passed")
	}
}
