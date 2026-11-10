package spring

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/httpclient"
	"github.com/sirupsen/logrus"
)

const springDownloadMutation = "5. Spring download keeps URL before archive path"

type recordedInitializerOperation struct {
	Name   string
	First  string
	Second string
}

type recordingInitializerIO struct {
	operations  []recordedInitializerOperation
	downloadErr error
	unzipErr    error
	deleteErr   error
}

func (recording *recordingInitializerIO) dependencies() initializerDependencies {
	return initializerDependencies{IO: recording}
}

func (recording *recordingInitializerIO) Download(downloadURL, archivePath string) error {
	recording.operations = append(recording.operations, recordedInitializerOperation{
		Name: "download", First: downloadURL, Second: archivePath,
	})
	return recording.downloadErr
}

func (recording *recordingInitializerIO) Unzip(archivePath, targetDirectory string) ([]string, error) {
	recording.operations = append(recording.operations, recordedInitializerOperation{
		Name: "unzip", First: archivePath, Second: targetDirectory,
	})
	return []string{"complete/unzipped/file"}, recording.unzipErr
}

func (recording *recordingInitializerIO) Delete(archivePath string) error {
	recording.operations = append(recording.operations, recordedInitializerOperation{
		Name: "delete", First: archivePath,
	})
	return recording.deleteErr
}

func (recording *recordingInitializerIO) assertedOperations() ([]recordedInitializerOperation, error) {
	if len(recording.operations) == 0 {
		return nil, errors.New("recorded Spring initializer population is empty")
	}
	return recording.operations, nil
}

func TestSpringDownloadKeepsCompleteEncodedURLBeforeCompleteResolvedArchivePath(t *testing.T) {
	if springDownloadMutation == "" {
		t.Fatal("spring-download mutation label is empty")
	}
	previousBaseURL := baseUrl
	baseUrl = "https://spring.example.invalid/complete-base"
	t.Cleanup(func() { baseUrl = previousBaseURL })
	formData := url.Values{}
	formData.Add("z-last", "last value")
	formData.Add("a-first", "first/value")
	archivePath := "/complete working directory/spring-1700000000.zip"
	targetDirectory := "/complete resolved initializer target"
	recording := &recordingInitializerIO{}
	logOutput := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(logOutput)
	previousLogger := log
	log = logger
	t.Cleanup(func() { log = previousLogger })

	err := downloadInitializer(recording.dependencies(), targetDirectory, formData, archivePath)

	if err != nil {
		t.Fatalf("Spring initializer flow returned an error: %v", err)
	}
	wantURL := baseUrl + "/starter.zip?" + formData.Encode()
	want := []recordedInitializerOperation{
		{Name: "download", First: wantURL, Second: archivePath},
		{Name: "unzip", First: archivePath, Second: targetDirectory},
		{Name: "delete", First: archivePath},
	}
	assertRecordedInitializerOperations(t, recording, want)
	assertOrderedLogValues(t, logOutput.String(), []string{
		"Downloading from " + baseUrl + " to " + archivePath,
		"Unzipping " + archivePath + " to " + targetDirectory,
		"Deleting archive file: " + archivePath,
	})
}

func TestSpringDownloadReturnsDependencyErrorsInFollowUpOrder(t *testing.T) {
	archivePath := "/complete working directory/spring-1700000001.zip"
	targetDirectory := "/complete resolved error target"
	formData := url.Values{"name": {"complete error flow"}}
	wantURL := baseUrl + "/starter.zip?" + formData.Encode()
	downloadError := errors.New("download dependency failed")
	unzipError := errors.New("unzip dependency failed")
	deleteError := errors.New("delete dependency failed")
	tests := []struct {
		name      string
		recording *recordingInitializerIO
		wantError error
		want      []recordedInitializerOperation
	}{
		{
			name:      "download prevents unzip and delete",
			recording: &recordingInitializerIO{downloadErr: downloadError},
			wantError: downloadError,
			want: []recordedInitializerOperation{
				{Name: "download", First: wantURL, Second: archivePath},
			},
		},
		{
			name:      "unzip prevents delete",
			recording: &recordingInitializerIO{unzipErr: unzipError},
			wantError: unzipError,
			want: []recordedInitializerOperation{
				{Name: "download", First: wantURL, Second: archivePath},
				{Name: "unzip", First: archivePath, Second: targetDirectory},
			},
		},
		{
			name:      "delete remains final error",
			recording: &recordingInitializerIO{deleteErr: deleteError},
			wantError: deleteError,
			want: []recordedInitializerOperation{
				{Name: "download", First: wantURL, Second: archivePath},
				{Name: "unzip", First: archivePath, Second: targetDirectory},
				{Name: "delete", First: archivePath},
			},
		},
	}

	if len(tests) == 0 {
		t.Fatal("TestSpringDownloadReturnsDependencyErrorsInFollowUpOrder test-case population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := downloadInitializer(test.recording.dependencies(), targetDirectory, formData, archivePath)

			if !errors.Is(err, test.wantError) {
				t.Fatalf("Spring initializer dependency error was %v, want %v", err, test.wantError)
			}
			assertRecordedInitializerOperations(t, test.recording, test.want)
		})
	}
}

func TestSpringDownloadDependenciesDefaultToNoRequestOrMutation(t *testing.T) {
	archivePath := filepath.Join(t.TempDir(), "safe default", "spring-1700000002.zip")

	err := downloadInitializer(initializerDependencies{}, t.TempDir(), url.Values{"name": {"must not request"}}, archivePath)

	if !errors.Is(err, httpclient.ErrNoClient) {
		t.Fatalf("safe Spring download dependency default returned %v, want %v", err, httpclient.ErrNoClient)
	}
	if _, statErr := os.Stat(archivePath); !os.IsNotExist(statErr) {
		t.Fatalf("safe Spring download dependency default mutated the archive: %v", statErr)
	}
}

func TestRecordedSpringDownloadOperationsRejectEmptyPopulation(t *testing.T) {
	recording := &recordingInitializerIO{}

	if _, err := recording.assertedOperations(); err == nil {
		t.Fatal("empty recorded Spring initializer population passed")
	}
}

func TestArchivePathRetainsWorkingDirectoryAndSpringUnixZipName(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("resolve test working directory: %v", err)
	}
	before := time.Now().Unix()

	archive, err := archivePath()

	after := time.Now().Unix()
	if err != nil {
		t.Fatalf("resolve Spring archive path: %v", err)
	}
	if filepath.Dir(archive) != workingDirectory {
		t.Fatalf("archive directory was %q, want %q", filepath.Dir(archive), workingDirectory)
	}
	matchedTimestamp := false
	for timestamp := before; timestamp <= after; timestamp++ {
		if archive == filepath.Join(workingDirectory, "spring-"+fmt.Sprint(timestamp)+".zip") {
			matchedTimestamp = true
			break
		}
	}
	if !matchedTimestamp {
		t.Fatalf("archive path %q did not preserve the public system composition in Unix-second range [%d, %d]",
			archive, before, after)
	}
}

func assertRecordedInitializerOperations(t *testing.T, recording *recordingInitializerIO, want []recordedInitializerOperation) {
	t.Helper()
	operations, err := recording.assertedOperations()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(operations, want) {
		t.Fatalf("recorded Spring initializer operations differ:\n got: %#v\nwant: %#v", operations, want)
	}
}

func assertOrderedLogValues(t *testing.T, output string, values []string) {
	t.Helper()
	if len(values) == 0 {
		t.Fatal("Spring initializer ordered-log expectation population is empty")
	}
	previous := -1
	for _, value := range values {
		index := strings.Index(output, value)
		if index < 0 || index <= previous {
			t.Fatalf("Spring initializer logging lost order for %q:\n%s", value, output)
		}
		previous = index
	}
}
