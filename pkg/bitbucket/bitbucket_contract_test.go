package bitbucket

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/internal/adapter/httpclient"
	"github.com/devdimensionlab/plybuild/pkg/shell"
	"github.com/sirupsen/logrus"
)

type recordedBitbucketQuery struct {
	Request     httpclient.Request
	Destination interface{}
}

type recordingBitbucketQueryHTTP struct {
	queries   []recordedBitbucketQuery
	responses map[string]string
	errors    map[string]error
}

func (recording *recordingBitbucketQueryHTTP) dependencies() queryDependencies {
	return queryDependencies{Client: recording}
}

func (recording *recordingBitbucketQueryHTTP) GetBearerJSON(request httpclient.Request, parsed interface{}) error {
	if request.BasicAuth != nil {
		copied := *request.BasicAuth
		request.BasicAuth = &copied
	}
	if request.BearerJSON != nil {
		copied := *request.BearerJSON
		request.BearerJSON = &copied
	}
	recording.queries = append(recording.queries, recordedBitbucketQuery{
		Request:     request,
		Destination: parsed,
	})
	if err := recording.errors[request.URL]; err != nil {
		return err
	}
	return json.Unmarshal([]byte(recording.responses[request.URL]), parsed)
}

func (recording *recordingBitbucketQueryHTTP) assertedQueries() ([]recordedBitbucketQuery, error) {
	if len(recording.queries) == 0 {
		return nil, errors.New("recorded Bitbucket query population is empty")
	}
	return recording.queries, nil
}

type recordingBitbucketRepositoryFilesystem struct {
	statPaths []string
	statErr   error
}

func (*recordingBitbucketRepositoryFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected Bitbucket repository read")
}

func (*recordingBitbucketRepositoryFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, errors.New("unexpected Bitbucket repository read directory")
}

func (recording *recordingBitbucketRepositoryFilesystem) Stat(path string) (fs.FileInfo, error) {
	recording.statPaths = append(recording.statPaths, path)
	return nil, recording.statErr
}

func (*recordingBitbucketRepositoryFilesystem) Mkdir(string, fs.FileMode) error {
	return errors.New("unexpected Bitbucket repository mkdir")
}

func (*recordingBitbucketRepositoryFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected Bitbucket repository mkdir")
}

func (*recordingBitbucketRepositoryFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected Bitbucket repository write")
}

func (*recordingBitbucketRepositoryFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected Bitbucket repository open file")
}

func (*recordingBitbucketRepositoryFilesystem) Remove(string) error {
	return errors.New("unexpected Bitbucket repository remove")
}

func (*recordingBitbucketRepositoryFilesystem) RemoveAll(string) error {
	return errors.New("unexpected Bitbucket repository recursive remove")
}

func (*recordingBitbucketRepositoryFilesystem) Rename(string, string) error {
	return errors.New("unexpected Bitbucket repository rename")
}

func (*recordingBitbucketRepositoryFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected Bitbucket repository glob")
}

func (*recordingBitbucketRepositoryFilesystem) Walk(string, filepath.WalkFunc) error {
	return errors.New("unexpected Bitbucket repository walk")
}

func (*recordingBitbucketRepositoryFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected Bitbucket repository create")
}

func (*recordingBitbucketRepositoryFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected Bitbucket repository copy")
}

func (recording *recordingBitbucketRepositoryFilesystem) assertedStatPaths() ([]string, error) {
	if len(recording.statPaths) == 0 {
		return nil, errors.New("recorded Bitbucket repository stat population is empty")
	}
	return recording.statPaths, nil
}

type recordedBitbucketRepositoryGitCall struct {
	Operation string
	Values    []string
}

type recordingBitbucketRepositoryGit struct {
	calls       []recordedBitbucketRepositoryGitCall
	cloneOutput shell.Output
	pullOutput  shell.Output
}

func (recording *recordingBitbucketRepositoryGit) dependencies(files *recordingBitbucketRepositoryFilesystem) repositoryDependencies {
	return repositoryDependencies{
		Files: filesystem.Dependencies{FileSystem: files},
		Git:   recording,
	}
}

func (recording *recordingBitbucketRepositoryGit) Clone(url string, target string) shell.Output {
	recording.calls = append(recording.calls, recordedBitbucketRepositoryGitCall{
		Operation: "clone",
		Values:    []string{url, target},
	})
	return recording.cloneOutput
}

func (recording *recordingBitbucketRepositoryGit) Pull(target string) shell.Output {
	recording.calls = append(recording.calls, recordedBitbucketRepositoryGitCall{
		Operation: "pull",
		Values:    []string{target},
	})
	return recording.pullOutput
}

func (recording *recordingBitbucketRepositoryGit) assertedCalls() ([]recordedBitbucketRepositoryGitCall, error) {
	if len(recording.calls) == 0 {
		return nil, errors.New("recorded Bitbucket repository Git population is empty")
	}
	return recording.calls, nil
}

func TestBitbucketWithSelectsCompleteSystemRepositoryDependencies(t *testing.T) {
	client := With(logrus.New(), "https://system.bitbucket.example.invalid", "system-token")
	systemFiles := filesystem.System()

	if client.repositories.Files.FileSystem == nil || client.repositories.Git == nil {
		t.Fatalf("With selected incomplete repository dependencies: %#v", client.repositories)
	}
	if reflect.TypeOf(client.repositories.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("With filesystem dependency is %T, want %T", client.repositories.Files.FileSystem, systemFiles.FileSystem)
	}
	if reflect.TypeOf(client.repositories.Git) != reflect.TypeOf(packageRepositoryGit{}) {
		t.Fatalf("With Git dependency is %T, want %T", client.repositories.Git, packageRepositoryGit{})
	}
}

func TestBitbucketMissingRepositorySelectsCloneWithCompletePathAndLogging(t *testing.T) {
	host := "ssh://git@bitbucket.example.invalid:7999/complete-base"
	workspace := "/complete workspace with spaces"
	repository := "/COMPLETE-PROJECT/complete repository"
	repositoryPath := workspace + repository
	files := &recordingBitbucketRepositoryFilesystem{statErr: &os.PathError{
		Op: "stat", Path: repositoryPath, Err: fs.ErrNotExist,
	}}
	git := &recordingBitbucketRepositoryGit{}
	logOutput := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(logOutput)
	logger.SetFormatter(&logrus.TextFormatter{DisableColors: true, DisableTimestamp: true})
	client := Bitbucket{
		host: host, log: logger, repositories: git.dependencies(files),
	}

	err := client.cloneOrPull(workspace, repository)

	if err != nil {
		t.Fatalf("missing Bitbucket repository selection returned an error: %v", err)
	}
	assertRecordedBitbucketRepositoryStatPaths(t, files, []string{repositoryPath})
	wantURL := host + "/scm" + repository + ".git"
	assertRecordedBitbucketRepositoryGitCalls(t, git, []recordedBitbucketRepositoryGitCall{{
		Operation: "clone",
		Values:    []string{wantURL, repositoryPath},
	}})
	wantLog := "clone [" + wantURL + "] -> [" + repositoryPath + "]"
	if !strings.Contains(logOutput.String(), wantLog) {
		t.Fatalf("Bitbucket clone log lost exact values:\n got: %s\nwant message: %s", logOutput.String(), wantLog)
	}
}

func TestBitbucketEveryNonMissingRepositoryResultSelectsPullWithExactPathAndLogging(t *testing.T) {
	statError := errors.New("permission denied while probing complete repository")
	tests := []struct {
		name    string
		statErr error
	}{
		{name: "existing repository"},
		{name: "other stat error", statErr: statError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			workspace := "/complete pull workspace"
			repository := "/complete-project/complete-repository"
			repositoryPath := workspace + repository
			pullPath := workspace + "/" + repository
			files := &recordingBitbucketRepositoryFilesystem{statErr: test.statErr}
			git := &recordingBitbucketRepositoryGit{}
			logOutput := &bytes.Buffer{}
			logger := logrus.New()
			logger.SetLevel(logrus.DebugLevel)
			logger.SetOutput(logOutput)
			logger.SetFormatter(&logrus.TextFormatter{DisableColors: true, DisableTimestamp: true})
			client := Bitbucket{log: logger, repositories: git.dependencies(files)}

			err := client.cloneOrPull(workspace, repository)

			if err != nil {
				t.Fatalf("non-missing Bitbucket repository selection returned an error: %v", err)
			}
			assertRecordedBitbucketRepositoryStatPaths(t, files, []string{repositoryPath})
			assertRecordedBitbucketRepositoryGitCalls(t, git, []recordedBitbucketRepositoryGitCall{{
				Operation: "pull",
				Values:    []string{pullPath},
			}})
			wantLog := " pull [" + pullPath + "]"
			if !strings.Contains(logOutput.String(), wantLog) {
				t.Fatalf("Bitbucket pull log lost exact values:\n got: %s\nwant message: %s", logOutput.String(), wantLog)
			}
		})
	}
}

func TestBitbucketCloneOrPullPropagatesRecordedOperationErrors(t *testing.T) {
	tests := []struct {
		name      string
		statErr   error
		operation string
	}{
		{name: "clone", statErr: fs.ErrNotExist, operation: "clone"},
		{name: "pull", operation: "pull"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sentinel := errors.New("complete " + test.name + " dependency failure")
			output := shell.Output{Err: sentinel}
			_, _ = output.StdOut.WriteString(test.name + " complete stdout\n")
			_, _ = output.StdErr.WriteString(test.name + " complete stderr\n")
			files := &recordingBitbucketRepositoryFilesystem{statErr: test.statErr}
			git := &recordingBitbucketRepositoryGit{cloneOutput: output, pullOutput: output}
			client := Bitbucket{
				host: "ssh://git@errors.bitbucket.example.invalid:7999",
				log:  logrus.New(), repositories: git.dependencies(files),
			}

			err := client.cloneOrPull("/complete error workspace", "/complete-project/complete-repository")

			if err == nil || err.Error() != output.FormatError().Error() {
				t.Fatalf("Bitbucket %s error was %v, want %v", test.name, err, output.FormatError())
			}
			calls, callsErr := git.assertedCalls()
			if callsErr != nil {
				t.Fatal(callsErr)
			}
			if len(calls) != 1 || calls[0].Operation != test.operation {
				t.Fatalf("Bitbucket operation calls were %#v, want one %s", calls, test.operation)
			}
		})
	}
}

func TestBitbucketRepositoryDependenciesDefaultToNoGitProcessOrMutation(t *testing.T) {
	workspace := t.TempDir()
	repository := "/safe-default/must-not-clone"
	repositoryPath := filepath.Join(workspace, repository)
	client := Bitbucket{
		host: "ssh://must-not-run.example.invalid", log: logrus.New(),
	}
	wantOutput := shell.Output{Err: filesystem.ErrNoFilesystem}

	err := client.cloneOrPull(workspace, repository)

	if err == nil || err.Error() != wantOutput.FormatError().Error() {
		t.Fatalf("safe Bitbucket repository default returned %v, want %v", err, wantOutput.FormatError())
	}
	if _, statErr := os.Stat(repositoryPath); !os.IsNotExist(statErr) {
		t.Fatalf("safe Bitbucket repository default mutated the target: %v", statErr)
	}
}

func TestRecordedBitbucketRepositoryDependenciesRejectEmptyPopulations(t *testing.T) {
	files := &recordingBitbucketRepositoryFilesystem{}
	git := &recordingBitbucketRepositoryGit{}

	if _, err := files.assertedStatPaths(); err == nil {
		t.Fatal("empty recorded Bitbucket repository stat population passed")
	}
	if _, err := git.assertedCalls(); err == nil {
		t.Fatal("empty recorded Bitbucket repository Git population passed")
	}
}

func TestBitbucketQueriesKeepCompleteURLsBearerValuesAndParsedResponses(t *testing.T) {
	host := "https://bitbucket.example.invalid:8443/complete-base"
	token := "complete-bitbucket-token"
	projectPath := "/rest/api/1.0/projects?limit=500"
	repositoryPath := "/rest/api/1.0/projects/complete-project/repos?limit=1000"
	recording := &recordingBitbucketQueryHTTP{responses: map[string]string{
		host + projectPath: `{
			"size":1,"limit":500,"isLastPage":true,"start":0,"nextPageStart":1,
			"values":[{"key":"COMPLETE-PROJECT","id":41,"name":"Complete Project","public":true,"type":"NORMAL"}]
		}`,
		host + repositoryPath: `{
			"size":1,"limit":1000,"isLastPage":true,"start":0,
			"values":[{"slug":"complete-repo","id":73,"name":"Complete Repo","scmId":"git","state":"AVAILABLE","project":{"key":"COMPLETE-PROJECT"}}]
		}`,
	}}
	client := With(logrus.New(), host, token)
	client.queries = recording.dependencies()

	projects, err := client.queryProjects()
	if err != nil {
		t.Fatalf("Bitbucket project query returned an error: %v", err)
	}
	repositories, err := client.queryRepos("complete-project")
	if err != nil {
		t.Fatalf("Bitbucket repository query returned an error: %v", err)
	}

	if projects.Size != 1 || projects.Limit != 500 || !projects.IsLastPage ||
		len(projects.Values) != 1 || projects.Values[0].Key != "COMPLETE-PROJECT" ||
		projects.Values[0].ID != 41 || projects.Values[0].Name != "Complete Project" ||
		!projects.Values[0].Public || projects.Values[0].Type != "NORMAL" {
		t.Fatalf("parsed Bitbucket projects response was incomplete: %#v", projects)
	}
	if repositories.Size != 1 || repositories.Limit != 1000 || !repositories.IsLastPage ||
		len(repositories.BitBucketRepo) != 1 || repositories.BitBucketRepo[0].Slug != "complete-repo" ||
		repositories.BitBucketRepo[0].ID != 73 || repositories.BitBucketRepo[0].Name != "Complete Repo" ||
		repositories.BitBucketRepo[0].ScmID != "git" || repositories.BitBucketRepo[0].State != "AVAILABLE" ||
		repositories.BitBucketRepo[0].Project.Key != "COMPLETE-PROJECT" {
		t.Fatalf("parsed Bitbucket repositories response was incomplete: %#v", repositories)
	}
	assertRecordedBitbucketQueries(t, recording, []recordedBitbucketQuery{
		{
			Request: httpclient.Request{
				URL: host + projectPath,
				BearerJSON: &httpclient.BearerJSON{
					AccessToken: token,
				},
			},
			Destination: (*ProjectList)(nil),
		},
		{
			Request: httpclient.Request{
				URL: host + repositoryPath,
				BearerJSON: &httpclient.BearerJSON{
					AccessToken: token,
				},
			},
			Destination: (*ProjectRepos)(nil),
		},
	})
}

func TestBitbucketQueriesReturnDependencyErrorsAfterCompleteDelivery(t *testing.T) {
	host := "https://errors.bitbucket.example.invalid/complete-base"
	token := "complete-error-token"
	projectURL := host + "/rest/api/1.0/projects?limit=500"
	repositoryURL := host + "/rest/api/1.0/projects/error-project/repos?limit=1000"
	sentinel := errors.New("Bitbucket query dependency failed")
	tests := []struct {
		name        string
		url         string
		destination interface{}
		invoke      func(Bitbucket) error
	}{
		{
			name:        "projects",
			url:         projectURL,
			destination: (*ProjectList)(nil),
			invoke: func(client Bitbucket) error {
				_, err := client.queryProjects()
				return err
			},
		},
		{
			name:        "repositories",
			url:         repositoryURL,
			destination: (*ProjectRepos)(nil),
			invoke: func(client Bitbucket) error {
				_, err := client.queryRepos("error-project")
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingBitbucketQueryHTTP{errors: map[string]error{test.url: sentinel}}
			client := With(logrus.New(), host, token)
			client.queries = recording.dependencies()

			err := test.invoke(client)

			if !errors.Is(err, sentinel) {
				t.Fatalf("Bitbucket query dependency error was %v, want %v", err, sentinel)
			}
			queries, assertErr := recording.assertedQueries()
			if assertErr != nil {
				t.Fatal(assertErr)
			}
			wantRequest := httpclient.Request{
				URL: test.url,
				BearerJSON: &httpclient.BearerJSON{
					AccessToken: token,
				},
			}
			if len(queries) != 1 || !reflect.DeepEqual(queries[0].Request, wantRequest) {
				t.Fatalf("Bitbucket error query was %#v, want %#v", queries, wantRequest)
			}
			if reflect.TypeOf(queries[0].Destination) != reflect.TypeOf(test.destination) {
				t.Fatalf("Bitbucket error destination was %T, want %T", queries[0].Destination, test.destination)
			}
		})
	}
}

func TestBitbucketQueryDependenciesDefaultToSafeNoRequest(t *testing.T) {
	client := Bitbucket{
		host:        "https://must-not-request.example.invalid",
		accessToken: "must-not-request-token",
		log:         logrus.New(),
	}

	if _, err := client.queryProjects(); !errors.Is(err, httpclient.ErrNoClient) {
		t.Fatalf("safe Bitbucket project-query default returned %v, want %v", err, httpclient.ErrNoClient)
	}
	if _, err := client.queryRepos("must-not-request"); !errors.Is(err, httpclient.ErrNoClient) {
		t.Fatalf("safe Bitbucket repository-query default returned %v, want %v", err, httpclient.ErrNoClient)
	}
}

func TestBitbucketSynchronizationSelectsProjectThenLowercaseRepositoryQueryAndWarns(t *testing.T) {
	host := "https://sync.bitbucket.example.invalid"
	token := "synchronization-token"
	projectURL := host + "/rest/api/1.0/projects?limit=500"
	repositoryURL := host + "/rest/api/1.0/projects/mixed-project/repos?limit=1000"
	sentinel := errors.New("repository query warning sentinel")
	recording := &recordingBitbucketQueryHTTP{
		responses: map[string]string{
			projectURL: `{"values":[{"key":"MiXeD-PrOjEcT"}]}`,
		},
		errors: map[string]error{repositoryURL: sentinel},
	}
	logOutput := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetOutput(logOutput)
	client := With(logger, host, token)
	client.queries = recording.dependencies()

	err := client.SynchronizeAllRepos(nil)

	if err != nil {
		t.Fatalf("Bitbucket synchronization returned an error: %v", err)
	}
	queries, assertErr := recording.assertedQueries()
	if assertErr != nil {
		t.Fatal(assertErr)
	}
	if len(queries) != 2 || queries[0].Request.URL != projectURL || queries[1].Request.URL != repositoryURL {
		t.Fatalf("Bitbucket synchronization query order was %#v", queries)
	}
	if !strings.Contains(logOutput.String(), sentinel.Error()) {
		t.Fatalf("Bitbucket synchronization did not warn on repository query error:\n%s", logOutput.String())
	}
}

func TestRecordedBitbucketQueriesRejectEmptyPopulation(t *testing.T) {
	recording := &recordingBitbucketQueryHTTP{}

	if _, err := recording.assertedQueries(); err == nil {
		t.Fatal("empty recorded Bitbucket query population passed")
	}
}

func assertRecordedBitbucketQueries(t *testing.T, recording *recordingBitbucketQueryHTTP, want []recordedBitbucketQuery) {
	t.Helper()
	queries, err := recording.assertedQueries()
	if err != nil {
		t.Fatal(err)
	}
	gotRequests := make([]httpclient.Request, len(queries))
	wantRequests := make([]httpclient.Request, len(want))
	for index, query := range queries {
		gotRequests[index] = query.Request
		if index >= len(want) {
			continue
		}
		wantRequests[index] = want[index].Request
		if reflect.TypeOf(query.Destination) != reflect.TypeOf(want[index].Destination) {
			t.Fatalf("Bitbucket query destination %d was %T, want %T", index, query.Destination, want[index].Destination)
		}
	}
	if !reflect.DeepEqual(gotRequests, wantRequests) {
		t.Fatalf("recorded Bitbucket queries differ:\n got: %#v\nwant: %#v", gotRequests, wantRequests)
	}
}

func assertRecordedBitbucketRepositoryStatPaths(
	t *testing.T,
	recording *recordingBitbucketRepositoryFilesystem,
	want []string,
) {
	t.Helper()
	paths, err := recording.assertedStatPaths()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(paths, want) {
		t.Fatalf("recorded Bitbucket repository stat paths differ:\n got: %#v\nwant: %#v", paths, want)
	}
}

func assertRecordedBitbucketRepositoryGitCalls(
	t *testing.T,
	recording *recordingBitbucketRepositoryGit,
	want []recordedBitbucketRepositoryGitCall,
) {
	t.Helper()
	calls, err := recording.assertedCalls()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("recorded Bitbucket repository Git calls differ:\n got: %#v\nwant: %#v", calls, want)
	}
}
