package resources

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/testutil"
	"github.com/devdimensionlab/plybuild/pkg/config"
)

var (
	_ config.Directory                                 = (*recordingResourceDirectory)(nil)
	_ config.CloudConfig                               = (*recordingResourceCloudConfig)(nil)
	_ func(config.CloudConfig) string                  = LocalDir
	_ func(config.CloudConfig, string) (string, error) = ResourceAsString
)

type recordingResourceDirectory struct {
	path          string
	dirCalls      int
	filePathCalls int
}

func (directory *recordingResourceDirectory) Dir() string {
	directory.dirCalls++
	return directory.path
}

func (directory *recordingResourceDirectory) FilePath(string) (string, error) {
	directory.filePathCalls++
	return "", errors.New("unexpected resource FilePath call")
}

type recordingResourceCloudConfig struct {
	config.GitCloudConfig
	directory           config.Directory
	implementationCalls int
}

func (cloudConfig *recordingResourceCloudConfig) Implementation() config.Directory {
	cloudConfig.implementationCalls++
	return cloudConfig.directory
}

func TestLocalDirPreservesExactImplementationDirectorySelectionAndPathBytes(t *testing.T) {
	separator := string(os.PathSeparator)
	tests := []struct {
		name string
		root string
		want string
	}{
		{
			name: "representative arbitrary supported bytes",
			root: "relative root #[]\x7f/inner",
			want: "relative root #[]\x7f" + separator + "inner" + separator + "resources",
		},
		{name: "empty root", root: "", want: separator + "resources"},
		{
			name: "trailing separator is retained",
			root: "trailing root" + separator,
			want: "trailing root" + separator + separator + "resources",
		},
	}
	if len(tests) == 0 {
		t.Fatal("LocalDir path characterization population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := &recordingResourceDirectory{path: test.root}
			cloudConfig := &recordingResourceCloudConfig{directory: directory}

			got := LocalDir(cloudConfig)

			if got != test.want {
				t.Fatalf("LocalDir returned %q, want exact path bytes %q", got, test.want)
			}
			assertExactResourceDirectorySelection(t, cloudConfig, directory)
		})
	}
}

func TestResourceAsStringPreservesExactFilenameSelectionAndContentBytes(t *testing.T) {
	root := filepath.Join(t.TempDir(), "resource root #[] ø")
	nestedParent := filepath.Join(root, "resources", "nested", "deeper")
	if err := os.MkdirAll(nestedParent, 0700); err != nil {
		t.Fatalf("prepare resource fixture parent: %v", err)
	}

	tests := []struct {
		name     string
		filename string
		content  []byte
	}{
		{name: "empty content", filename: "empty.resource", content: []byte{}},
		{
			name:     "arbitrary filename and content bytes",
			filename: "arbitrary name #[] ø.resource",
			content:  []byte{0x00, 'a', '\r', '\n', 0x7f, 0xc3, 0xb8},
		},
		{
			name:     "non UTF-8 content",
			filename: "non-utf8.resource",
			content:  []byte{0x00, 0xff, 0xfe, 'x', 0x80},
		},
		{
			name:     "nested resource name",
			filename: "nested/deeper/content.resource",
			content:  []byte("nested\ncontent"),
		},
	}
	if len(tests) == 0 {
		t.Fatal("ResourceAsString success characterization population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := &recordingResourceDirectory{path: root}
			cloudConfig := &recordingResourceCloudConfig{directory: directory}
			resourcePath := filepath.Join(root, "resources") + "/" + test.filename
			if err := testutil.WriteFileOutsideWorkingTree(resourcePath, test.content, 0600); err != nil {
				t.Fatalf("write temporary resource fixture: %v", err)
			}

			got, err := ResourceAsString(cloudConfig, test.filename)

			if err != nil {
				t.Fatalf("ResourceAsString returned error: %v", err)
			}
			if !bytes.Equal([]byte(got), test.content) {
				t.Fatalf("ResourceAsString returned bytes %v, want exact bytes %v", []byte(got), test.content)
			}
			assertExactResourceDirectorySelection(t, cloudConfig, directory)
		})
	}
}

func TestResourceAsStringReturnsDirectMissingReadErrorWithLiteralSlashPath(t *testing.T) {
	root := filepath.Join(t.TempDir(), "missing resource root #[]") + string(os.PathSeparator)
	filename := "nested/../missing name #[]\x7f.resource"
	directory := &recordingResourceDirectory{path: root}
	cloudConfig := &recordingResourceCloudConfig{directory: directory}
	wantResourcesDir := root + string(os.PathSeparator) + "resources"
	wantPath := wantResourcesDir + "/" + filename

	got, err := ResourceAsString(cloudConfig, filename)

	if got != "" {
		t.Fatalf("missing resource returned %q, want empty string", got)
	}
	pathError, ok := err.(*os.PathError)
	if !ok {
		t.Fatalf("missing resource error is %T (%v), want direct *os.PathError", err, err)
	}
	if pathError.Op != "open" || pathError.Path != wantPath || !errors.Is(pathError.Err, fs.ErrNotExist) {
		t.Fatalf("missing resource error was %#v, want open of exact literal path %q with missing-file cause", pathError, wantPath)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("missing resource error %v does not preserve fs.ErrNotExist", err)
	}
	assertExactResourceDirectorySelection(t, cloudConfig, directory)
}

func assertExactResourceDirectorySelection(t *testing.T, cloudConfig *recordingResourceCloudConfig, directory *recordingResourceDirectory) {
	t.Helper()
	if cloudConfig.implementationCalls != 1 || directory.dirCalls != 1 || directory.filePathCalls != 0 {
		t.Fatalf("resource directory selection calls were Implementation=%d Dir=%d FilePath=%d, want 1, 1, 0",
			cloudConfig.implementationCalls, directory.dirCalls, directory.filePathCalls)
	}
}
