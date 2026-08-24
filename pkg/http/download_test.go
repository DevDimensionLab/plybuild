package http

import (
	"errors"
	"io"
	"io/fs"
	stdhttp "net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/internal/adapter/httpclient"
	"github.com/devdimensionlab/plybuild/internal/testutil"
)

type recordedDownloadFile struct {
	lifecycle *[]string
	closed    bool
}

func (file *recordedDownloadFile) Write([]byte) (int, error) {
	return 0, errors.New("recording download file must be copied through its filesystem dependency")
}

func (file *recordedDownloadFile) Close() error {
	*file.lifecycle = append(*file.lifecycle, "file-close")
	file.closed = true
	return nil
}

type recordingDownloadFilesystem struct {
	lifecycle *[]string
	path      string
	data      []byte
	file      *recordedDownloadFile
	createErr error
	copyErr   error
}

func (recording *recordingDownloadFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected read")
}

func (recording *recordingDownloadFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected stat")
}

func (recording *recordingDownloadFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected mkdir")
}

func (recording *recordingDownloadFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected write file")
}

func (recording *recordingDownloadFilesystem) Create(path string) (filesystem.File, error) {
	*recording.lifecycle = append(*recording.lifecycle, "create")
	recording.path = path
	if recording.createErr != nil {
		return nil, recording.createErr
	}
	recording.file = &recordedDownloadFile{lifecycle: recording.lifecycle}
	return recording.file, nil
}

func (recording *recordingDownloadFilesystem) Copy(destination filesystem.File, source io.Reader) (int64, error) {
	*recording.lifecycle = append(*recording.lifecycle, "copy")
	if destination != recording.file {
		return 0, errors.New("copy destination is not the created archive")
	}
	data, err := io.ReadAll(source)
	recording.data = append([]byte(nil), data...)
	if err != nil {
		return int64(len(data)), err
	}
	if recording.copyErr != nil {
		return int64(len(data)), recording.copyErr
	}
	return int64(len(data)), nil
}

type recordingDownloadClient struct {
	lifecycle *[]string
	requests  []httpclient.Request
	response  *stdhttp.Response
	err       error
}

func (client *recordingDownloadClient) Do(request httpclient.Request) (*stdhttp.Response, error) {
	*client.lifecycle = append(*client.lifecycle, "request")
	client.requests = append(client.requests, request)
	return client.response, client.err
}

type lifecycleReadCloser struct {
	reader    io.Reader
	lifecycle *[]string
	closed    bool
}

type downloadRoundTripFunc func(*stdhttp.Request) (*stdhttp.Response, error)

func (function downloadRoundTripFunc) RoundTrip(request *stdhttp.Request) (*stdhttp.Response, error) {
	return function(request)
}

func (body *lifecycleReadCloser) Read(buffer []byte) (int, error) {
	return body.reader.Read(buffer)
}

func (body *lifecycleReadCloser) Close() error {
	*body.lifecycle = append(*body.lifecycle, "response-close")
	body.closed = true
	return nil
}

func TestWgetPassesCompleteAnonymousRequestBeforeArchivePathAndCopiesBody(t *testing.T) {
	lifecycle := []string{}
	body := &lifecycleReadCloser{
		reader:    strings.NewReader("complete response bytes\nwith a second line\n"),
		lifecycle: &lifecycle,
	}
	client := &recordingDownloadClient{lifecycle: &lifecycle, response: &stdhttp.Response{
		StatusCode: stdhttp.StatusServiceUnavailable,
		Status:     "503 Must Still Be Copied",
		Body:       body,
	}}
	files := &recordingDownloadFilesystem{lifecycle: &lifecycle}
	dependencies := wgetDependencies{
		HTTP:  httpclient.Dependencies{Client: client},
		Files: filesystem.Dependencies{FileSystem: files},
	}
	downloadURL := "https://spring.example.invalid/starter.zip?a=first%2Fvalue&z=last+value"
	archivePath := "/complete working directory/spring-1700000000.zip"

	err := wget(dependencies, downloadURL, archivePath)

	if err != nil {
		t.Fatalf("Wget returned an error: %v", err)
	}
	if !reflect.DeepEqual(client.requests, []httpclient.Request{{URL: downloadURL}}) {
		t.Fatalf("anonymous download requests differ:\n got: %#v\nwant: %#v", client.requests, []httpclient.Request{{URL: downloadURL}})
	}
	if files.path != archivePath || string(files.data) != "complete response bytes\nwith a second line\n" {
		t.Fatalf("archive path/bytes were (%q, %q)", files.path, files.data)
	}
	wantLifecycle := []string{"request", "create", "copy", "file-close", "response-close"}
	if !reflect.DeepEqual(lifecycle, wantLifecycle) || files.file == nil || !files.file.closed || !body.closed {
		t.Fatalf("download lifecycle differs:\n got: %#v\nwant: %#v", lifecycle, wantLifecycle)
	}
}

func TestWgetReturnsDependencyErrorsAndPreservesCloseOrder(t *testing.T) {
	requestError := errors.New("anonymous request failed")
	createError := errors.New("archive create failed")
	copyError := errors.New("archive body copy failed")
	tests := []struct {
		name          string
		requestErr    error
		createErr     error
		copyErr       error
		wantError     error
		wantLifecycle []string
	}{
		{name: "request before create", requestErr: requestError, wantError: requestError, wantLifecycle: []string{"request"}},
		{name: "create closes response", createErr: createError, wantError: createError, wantLifecycle: []string{"request", "create", "response-close"}},
		{name: "copy closes file before response", copyErr: copyError, wantError: copyError, wantLifecycle: []string{"request", "create", "copy", "file-close", "response-close"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lifecycle := []string{}
			body := &lifecycleReadCloser{reader: strings.NewReader("partial or complete bytes"), lifecycle: &lifecycle}
			client := &recordingDownloadClient{
				lifecycle: &lifecycle,
				response:  &stdhttp.Response{StatusCode: stdhttp.StatusOK, Status: "200 OK", Body: body},
				err:       test.requestErr,
			}
			files := &recordingDownloadFilesystem{
				lifecycle: &lifecycle, createErr: test.createErr, copyErr: test.copyErr,
			}

			err := wget(wgetDependencies{
				HTTP:  httpclient.Dependencies{Client: client},
				Files: filesystem.Dependencies{FileSystem: files},
			}, "https://error.example.invalid/starter.zip?complete=value", "/complete/error/archive.zip")

			if !errors.Is(err, test.wantError) {
				t.Fatalf("Wget dependency error was %v, want %v", err, test.wantError)
			}
			if !reflect.DeepEqual(lifecycle, test.wantLifecycle) {
				t.Fatalf("error lifecycle differs:\n got: %#v\nwant: %#v", lifecycle, test.wantLifecycle)
			}
		})
	}
}

func TestWgetDependenciesDefaultToNoRequestOrArchiveMutation(t *testing.T) {
	archivePath := filepath.Join(t.TempDir(), "safe-default", "archive.zip")

	err := wget(wgetDependencies{}, "https://must-not-request.example.invalid/starter.zip", archivePath)

	if !errors.Is(err, httpclient.ErrNoClient) {
		t.Fatalf("zero-value Wget dependencies returned %v, want %v", err, httpclient.ErrNoClient)
	}
	if _, statErr := os.Stat(archivePath); !os.IsNotExist(statErr) {
		t.Fatalf("zero-value Wget dependencies mutated the archive: %v", statErr)
	}
}

func TestWgetSystemPreservesAnonymousTransportAndCreateTruncate(t *testing.T) {
	previousTransport := stdhttp.DefaultTransport
	t.Cleanup(func() { stdhttp.DefaultTransport = previousTransport })
	body := &trackingReadCloser{reader: strings.NewReader("short system body")}
	var recorded *stdhttp.Request
	stdhttp.DefaultTransport = downloadRoundTripFunc(func(request *stdhttp.Request) (*stdhttp.Response, error) {
		recorded = request
		return &stdhttp.Response{
			StatusCode: stdhttp.StatusBadGateway,
			Status:     "502 Must Still Be Copied",
			Body:       body,
		}, nil
	})
	archivePath := filepath.Join(t.TempDir(), "existing archive.zip")
	if err := testutil.WriteFileOutsideWorkingTree(archivePath, []byte("long existing archive contents that must be truncated"), 0600); err != nil {
		t.Fatalf("create existing archive fixture: %v", err)
	}

	err := Wget("https://system.example.invalid/starter.zip?complete=value", archivePath)

	if err != nil {
		t.Fatalf("system Wget returned an error: %v", err)
	}
	data, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatalf("read system archive: %v", err)
	}
	if string(data) != "short system body" {
		t.Fatalf("system archive contained %q", data)
	}
	if recorded == nil || recorded.Method != stdhttp.MethodGet || recorded.URL.String() != "https://system.example.invalid/starter.zip?complete=value" {
		t.Fatalf("system anonymous request was %#v", recorded)
	}
	if !body.closed {
		t.Fatal("system response body was not closed")
	}
}
