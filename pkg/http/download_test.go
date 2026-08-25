package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	stdhttp "net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/internal/adapter/httpclient"
	"github.com/devdimensionlab/plybuild/internal/testutil"
	"github.com/sirupsen/logrus"
)

type recordedDownloadFile struct {
	lifecycle *[]string
	closed    bool
	closeErr  error
}

func (file *recordedDownloadFile) Write([]byte) (int, error) {
	return 0, errors.New("recording download file must be copied through its filesystem dependency")
}

func (file *recordedDownloadFile) Close() error {
	*file.lifecycle = append(*file.lifecycle, "file-close")
	file.closed = true
	return file.closeErr
}

type recordingDownloadFilesystem struct {
	lifecycle  *[]string
	operations []string
	path       string
	data       []byte
	file       *recordedDownloadFile
	createErr  error
	copyErr    error
	closeErr   error
}

func (recording *recordingDownloadFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected read")
}

func (recording *recordingDownloadFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, errors.New("unexpected read directory")
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

func (recording *recordingDownloadFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected open file")
}

func (recording *recordingDownloadFilesystem) Remove(string) error {
	return errors.New("unexpected remove")
}

func (recording *recordingDownloadFilesystem) RemoveAll(string) error {
	return errors.New("unexpected recursive remove")
}

func (recording *recordingDownloadFilesystem) Rename(string, string) error {
	return errors.New("unexpected rename")
}

func (recording *recordingDownloadFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected glob")
}

func (recording *recordingDownloadFilesystem) Walk(string, filepath.WalkFunc) error {
	return errors.New("unexpected walk")
}

func (recording *recordingDownloadFilesystem) Create(path string) (filesystem.File, error) {
	*recording.lifecycle = append(*recording.lifecycle, "create")
	recording.operations = append(recording.operations, "create")
	recording.path = path
	if recording.createErr != nil {
		return nil, recording.createErr
	}
	recording.file = &recordedDownloadFile{lifecycle: recording.lifecycle, closeErr: recording.closeErr}
	return recording.file, nil
}

func (recording *recordingDownloadFilesystem) Copy(destination filesystem.File, source io.Reader) (int64, error) {
	*recording.lifecycle = append(*recording.lifecycle, "copy")
	recording.operations = append(recording.operations, "copy")
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

func (recording *recordingDownloadFilesystem) assertedOperations() ([]string, error) {
	if len(recording.operations) == 0 {
		return nil, errors.New("recorded download filesystem population is empty")
	}
	return recording.operations, nil
}

type recordingDownloadClient struct {
	lifecycle *[]string
	requests  []httpclient.Request
	response  *stdhttp.Response
	err       error
}

func (client *recordingDownloadClient) Do(request httpclient.Request) (*stdhttp.Response, error) {
	*client.lifecycle = append(*client.lifecycle, "request")
	if request.POST != nil {
		copied := *request.POST
		copied.Body = append([]byte(nil), request.POST.Body...)
		copied.Header = request.POST.Header.Clone()
		request.POST = &copied
	}
	client.requests = append(client.requests, request)
	return client.response, client.err
}

func (client *recordingDownloadClient) assertedRequests() ([]httpclient.Request, error) {
	if len(client.requests) == 0 {
		return nil, errors.New("recorded download HTTP population is empty")
	}
	return client.requests, nil
}

type lifecycleReadCloser struct {
	reader    io.Reader
	lifecycle *[]string
	closed    bool
	closeErr  error
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
	return body.closeErr
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

func TestWpostPassesCompleteFormRequestBeforeFileAndCopiesNonSuccessBody(t *testing.T) {
	lifecycle := []string{}
	body := &lifecycleReadCloser{
		reader:    strings.NewReader("complete form response bytes\nwith a second line\n"),
		lifecycle: &lifecycle,
	}
	client := &recordingDownloadClient{lifecycle: &lifecycle, response: &stdhttp.Response{
		StatusCode: stdhttp.StatusServiceUnavailable,
		Status:     "503 Must Still Be Copied",
		Body:       body,
	}}
	files := &recordingDownloadFilesystem{lifecycle: &lifecycle}
	dependencies := wpostDependencies{
		HTTP:  httpclient.Dependencies{Client: client},
		Files: filesystem.Dependencies{FileSystem: files},
	}
	downloadURL := "https://form.example.invalid:8443/download/archive?complete=query"
	filePath := "/complete target directory/downloaded archive.zip"
	formData := url.Values{}
	formData.Add("z-last", "last value")
	formData.Add("a-first", "first/value")
	formData.Add("a-first", "second value")
	logOutput := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(logOutput)
	logger.SetFormatter(&logrus.TextFormatter{DisableColors: true, DisableTimestamp: true})
	previousLog := log
	log = logger
	t.Cleanup(func() { log = previousLog })

	err := wpost(dependencies, downloadURL, filePath, formData)

	if err != nil {
		t.Fatalf("Wpost returned an error: %v", err)
	}
	wantRequest := httpclient.Request{
		URL: downloadURL,
		POST: &httpclient.POST{
			Body:             []byte(formData.Encode()),
			Header:           stdhttp.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}},
			UseDefaultClient: true,
		},
	}
	requests, assertErr := client.assertedRequests()
	if assertErr != nil {
		t.Fatal(assertErr)
	}
	if !reflect.DeepEqual(requests, []httpclient.Request{wantRequest}) {
		t.Fatalf("form download requests differ:\n got: %#v\nwant: %#v", requests, []httpclient.Request{wantRequest})
	}
	if files.path != filePath || string(files.data) != "complete form response bytes\nwith a second line\n" {
		t.Fatalf("form archive path/bytes were (%q, %q)", files.path, files.data)
	}
	wantLifecycle := []string{"request", "create", "copy", "file-close", "response-close"}
	if !reflect.DeepEqual(lifecycle, wantLifecycle) || files.file == nil || !files.file.closed || !body.closed {
		t.Fatalf("form download lifecycle differs:\n got: %#v\nwant: %#v", lifecycle, wantLifecycle)
	}
	wantDebug := fmt.Sprintf("downloading %s to %s with %s", downloadURL, filePath, formData)
	if !strings.Contains(logOutput.String(), wantDebug) {
		t.Fatalf("Wpost debug log lost exact message or values:\n got: %s\nwant message: %s", logOutput.String(), wantDebug)
	}
}

func TestWpostReturnsDependencyErrorsAfterRecordingAndPreservesCloseOrder(t *testing.T) {
	requestError := errors.New("form request failed")
	createError := errors.New("form archive create failed")
	copyError := errors.New("form archive body copy failed")
	ignoredFileCloseError := errors.New("ignored form file close error")
	ignoredResponseCloseError := errors.New("ignored form response close error")
	tests := []struct {
		name             string
		requestErr       error
		createErr        error
		copyErr          error
		fileCloseErr     error
		responseCloseErr error
		wantError        error
		wantLifecycle    []string
	}{
		{name: "request before create", requestErr: requestError, wantError: requestError, wantLifecycle: []string{"request"}},
		{name: "create closes response", createErr: createError, wantError: createError, wantLifecycle: []string{"request", "create", "response-close"}},
		{name: "copy closes file before response", copyErr: copyError, wantError: copyError, wantLifecycle: []string{"request", "create", "copy", "file-close", "response-close"}},
		{
			name: "close errors are ignored", fileCloseErr: ignoredFileCloseError,
			responseCloseErr: ignoredResponseCloseError,
			wantLifecycle:    []string{"request", "create", "copy", "file-close", "response-close"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lifecycle := []string{}
			body := &lifecycleReadCloser{
				reader: strings.NewReader("partial or complete form bytes"), lifecycle: &lifecycle,
				closeErr: test.responseCloseErr,
			}
			client := &recordingDownloadClient{
				lifecycle: &lifecycle,
				response:  &stdhttp.Response{StatusCode: stdhttp.StatusUnauthorized, Status: "401 Must Still Be Copied", Body: body},
				err:       test.requestErr,
			}
			files := &recordingDownloadFilesystem{
				lifecycle: &lifecycle, createErr: test.createErr, copyErr: test.copyErr, closeErr: test.fileCloseErr,
			}
			formData := url.Values{"complete": {"error value"}}
			downloadURL := "https://error.example.invalid/form-download?complete=query"

			err := wpost(wpostDependencies{
				HTTP:  httpclient.Dependencies{Client: client},
				Files: filesystem.Dependencies{FileSystem: files},
			}, downloadURL, "/complete/error/form-archive.zip", formData)

			if !errors.Is(err, test.wantError) {
				t.Fatalf("Wpost dependency error was %v, want %v", err, test.wantError)
			}
			if !reflect.DeepEqual(lifecycle, test.wantLifecycle) {
				t.Fatalf("form error lifecycle differs:\n got: %#v\nwant: %#v", lifecycle, test.wantLifecycle)
			}
			requests, assertErr := client.assertedRequests()
			if assertErr != nil {
				t.Fatal(assertErr)
			}
			wantRequest := httpclient.Request{
				URL: downloadURL,
				POST: &httpclient.POST{
					Body:             []byte(formData.Encode()),
					Header:           stdhttp.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}},
					UseDefaultClient: true,
				},
			}
			if !reflect.DeepEqual(requests, []httpclient.Request{wantRequest}) {
				t.Fatalf("form error request differs:\n got: %#v\nwant: %#v", requests, []httpclient.Request{wantRequest})
			}
		})
	}
}

func TestWpostDependenciesDefaultToNoRequestOrFileMutation(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "safe-default", "form-archive.zip")

	err := wpost(wpostDependencies{}, "https://must-not-request.example.invalid/form-download", filePath, url.Values{"must": {"not request"}})

	if !errors.Is(err, httpclient.ErrNoClient) {
		t.Fatalf("zero-value Wpost dependencies returned %v, want %v", err, httpclient.ErrNoClient)
	}
	if _, statErr := os.Stat(filePath); !os.IsNotExist(statErr) {
		t.Fatalf("zero-value Wpost dependencies mutated the archive: %v", statErr)
	}
}

func TestWpostSystemPreservesDefaultClientRedirectsAndCreateTruncate(t *testing.T) {
	previousClient := stdhttp.DefaultClient
	previousTransport := stdhttp.DefaultTransport
	t.Cleanup(func() {
		stdhttp.DefaultClient = previousClient
		stdhttp.DefaultTransport = previousTransport
	})
	redirects := 0
	stdhttp.DefaultClient = &stdhttp.Client{CheckRedirect: func(request *stdhttp.Request, via []*stdhttp.Request) error {
		redirects++
		return nil
	}}
	lifecycle := []string{}
	redirectBody := &lifecycleReadCloser{reader: strings.NewReader("redirect"), lifecycle: &lifecycle}
	finalBody := &lifecycleReadCloser{reader: strings.NewReader("short redirected form body"), lifecycle: &lifecycle}
	requests := []*stdhttp.Request{}
	bodies := []string{}
	stdhttp.DefaultTransport = downloadRoundTripFunc(func(request *stdhttp.Request) (*stdhttp.Response, error) {
		requests = append(requests, request)
		body, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		bodies = append(bodies, string(body))
		if len(requests) == 1 {
			return &stdhttp.Response{
				StatusCode: stdhttp.StatusTemporaryRedirect,
				Status:     "307 Temporary Redirect",
				Header:     stdhttp.Header{"Location": []string{"https://system.example.invalid/final-form"}},
				Body:       redirectBody,
			}, nil
		}
		return &stdhttp.Response{
			StatusCode: stdhttp.StatusBadGateway,
			Status:     "502 Must Still Be Copied",
			Body:       finalBody,
		}, nil
	})
	filePath := filepath.Join(t.TempDir(), "existing form archive.zip")
	if err := testutil.WriteFileOutsideWorkingTree(filePath, []byte("long existing form archive contents that must be truncated"), 0600); err != nil {
		t.Fatalf("create existing form archive fixture: %v", err)
	}
	formData := url.Values{"z-last": {"last value"}, "a-first": {"first/value"}}

	err := Wpost("https://system.example.invalid/start-form", filePath, formData)

	if err != nil {
		t.Fatalf("system Wpost returned an error: %v", err)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read system form archive: %v", err)
	}
	if string(data) != "short redirected form body" {
		t.Fatalf("system form archive contained %q", data)
	}
	if redirects != 1 || len(requests) != 2 {
		t.Fatalf("default client redirects/requests were (%d, %d), want (1, 2)", redirects, len(requests))
	}
	wantURLs := []string{"https://system.example.invalid/start-form", "https://system.example.invalid/final-form"}
	for index, request := range requests {
		if request.Method != stdhttp.MethodPost || request.URL.String() != wantURLs[index] {
			t.Fatalf("system form request %d was %s %s", index, request.Method, request.URL)
		}
		if bodies[index] != formData.Encode() {
			t.Fatalf("system form request %d body was %q, want %q", index, bodies[index], formData.Encode())
		}
		if got := request.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Fatalf("system form request %d Content-Type was %q", index, got)
		}
	}
	if !redirectBody.closed || !finalBody.closed {
		t.Fatalf("system form redirect bodies were closed=(%t, %t)", redirectBody.closed, finalBody.closed)
	}
}

func TestRecordedWpostPopulationsRejectEmptyRequestsAndOperations(t *testing.T) {
	lifecycle := []string{}
	client := &recordingDownloadClient{lifecycle: &lifecycle}
	files := &recordingDownloadFilesystem{lifecycle: &lifecycle}

	if _, err := client.assertedRequests(); err == nil {
		t.Fatal("empty recorded Wpost HTTP population passed")
	}
	if _, err := files.assertedOperations(); err == nil {
		t.Fatal("empty recorded Wpost filesystem population passed")
	}
}
