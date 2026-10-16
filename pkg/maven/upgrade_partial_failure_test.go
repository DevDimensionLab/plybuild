package maven

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/sirupsen/logrus"
)

type mavenPartialFailureLog struct {
	level   logrus.Level
	message string
}

func TestUpgradeDependenciesReportsEveryPerItemFailureExactlyAndContinues(t *testing.T) {
	const repositoryURL = "https://dependencies.example.test/releases"
	dependencies := []pom.Dependency{
		{GroupId: "com.example", ArtifactId: "missing-version", Version: "${missing.version}"},
		{GroupId: "com.example", ArtifactId: "invalid-maximum", Version: "1.0.0"},
		{GroupId: "com.example", ArtifactId: "invalid-current", Version: "not-a-version"},
		{GroupId: "com.example", ArtifactId: "metadata-failure", Version: "1.0.0"},
		{GroupId: "com.example", ArtifactId: "later-success", Version: "1.0.0"},
	}
	if len(dependencies) == 0 {
		t.Fatal("dependency partial-failure population is empty")
	}
	model := &pom.Model{Dependencies: &pom.Dependencies{Dependency: dependencies}}
	settings := config.ProjectSettings{MaxVersionForDependencies: []config.MaxArtifact{{
		Artifact:   config.Artifact{GroupId: "com.example", ArtifactId: "invalid-maximum"},
		MaxVersion: "not-a-version",
	}}}
	hook := captureMavenPartialFailureLogs(t)
	var requestedURLs []string
	stubHTTPTransport(t, func(request *http.Request) (*http.Response, error) {
		requestedURLs = append(requestedURLs, request.URL.String())
		if strings.HasSuffix(request.URL.Path, "/metadata-failure/maven-metadata.xml") {
			return &http.Response{
				StatusCode: http.StatusServiceUnavailable,
				Status:     "503 Complete Dependency Status",
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("dependency status body")),
			}, nil
		}
		metadata := `<metadata><versioning><release>1.1.0</release><versions><version>1.1.0</version></versions></versioning></metadata>`
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(metadata)),
		}, nil
	})
	var actions []string
	action := func(dependency pom.Dependency, version string) error {
		actions = append(actions, dependency.GroupId+":"+dependency.ArtifactId+"="+version)
		return nil
	}

	Repository{Url: repositoryURL}.upgradeDependencies(model, dependencies, settings, func(string) bool { return true }, action)

	statusURL := repositoryURL + "/com/example/metadata-failure/maven-metadata.xml"
	assertMavenPartialFailureLogs(t, hook, []mavenPartialFailureLog{
		{level: logrus.InfoLevel, message: "failed to get version for com.example:missing-version"},
		{level: logrus.DebugLevel, message: "failed to get max version for com.example:invalid-maximum"},
		{level: logrus.WarnLevel, message: `unable to parse version:not-a-version due to strconv.Atoi: parsing "not": invalid syntax`},
		{level: logrus.WarnLevel, message: "could not determine release for com.example:metadata-failure from " + repositoryURL + ": HTTP 503 Service Unavailable"},
	})
	wantURLs := []string{
		repositoryURL + "/com/example/invalid-maximum/maven-metadata.xml",
		statusURL,
		repositoryURL + "/com/example/later-success/maven-metadata.xml",
	}
	if !reflect.DeepEqual(requestedURLs, wantURLs) {
		t.Fatalf("dependency metadata requests were %#v, want ordered continuation %#v", requestedURLs, wantURLs)
	}
	wantActions := []string{
		"com.example:invalid-maximum=1.1.0",
		"com.example:later-success=1.1.0",
	}
	if !reflect.DeepEqual(actions, wantActions) {
		t.Fatalf("dependency upgrade actions were %#v, want exact continued actions %#v", actions, wantActions)
	}
}

func TestUpgradePluginsOnModelReportsEveryPerPluginFailureExactlyAndContinues(t *testing.T) {
	const repositoryURL = "https://plugins.example.test/releases"
	plugins := []pom.Plugin{
		{GroupId: "com.example", ArtifactId: "missing-version", Version: "${missing.plugin.version}"},
		{GroupId: "com.example", ArtifactId: "invalid-current", Version: "not-a-version"},
		{GroupId: "com.example", ArtifactId: "metadata-failure", Version: "1.0.0"},
		{GroupId: "com.example", ArtifactId: "release-failure", Version: "1.0.0"},
		{GroupId: "com.example", ArtifactId: "later-success", Version: "1.0.0"},
	}
	if len(plugins) == 0 {
		t.Fatal("plugin partial-failure population is empty")
	}
	model := &pom.Model{Build: &pom.Build{Plugins: &pom.BuildPlugins{Plugin: plugins}}}
	project := config.Project{Type: config.MavenProject{PomModel: model}}
	hook := captureMavenPartialFailureLogs(t)
	var requestedURLs []string
	stubHTTPTransport(t, func(request *http.Request) (*http.Response, error) {
		requestedURLs = append(requestedURLs, request.URL.String())
		switch {
		case strings.HasSuffix(request.URL.Path, "/metadata-failure/maven-metadata.xml"):
			return &http.Response{
				StatusCode: http.StatusBadGateway,
				Status:     "502 Complete Plugin Status",
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("plugin status body")),
			}, nil
		case strings.HasSuffix(request.URL.Path, "/release-failure/maven-metadata.xml"):
			metadata := `<metadata><versioning><release>1.1.0-SNAPSHOT</release><versions><version>1.1.0-SNAPSHOT</version></versions></versioning></metadata>`
			return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(metadata))}, nil
		default:
			metadata := `<metadata><versioning><release>1.1.0</release><versions><version>1.1.0</version></versions></versioning></metadata>`
			return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(metadata))}, nil
		}
	})

	err := (Repository{Url: repositoryURL}).upgradePluginsOnModel(&project)

	if err != nil {
		t.Fatalf("plugin partial-failure operation returned an error: %v", err)
	}
	statusURL := repositoryURL + "/com/example/metadata-failure/maven-metadata.xml"
	assertMavenPartialFailureLogs(t, hook, []mavenPartialFailureLog{
		{level: logrus.WarnLevel, message: "version points at ${missing.plugin.version}, but no properties are defined"},
		{level: logrus.WarnLevel, message: `unable to parse version:not-a-version due to strconv.Atoi: parsing "not": invalid syntax`},
		{level: logrus.WarnLevel, message: "could not determine release for com.example:metadata-failure from " + repositoryURL + ": HTTP 502 Bad Gateway"},
		{level: logrus.WarnLevel, message: "could not determine release for com.example:release-failure from " + repositoryURL + ": could not find a suitable release version"},
	})
	wantURLs := []string{
		statusURL,
		repositoryURL + "/com/example/release-failure/maven-metadata.xml",
		repositoryURL + "/com/example/later-success/maven-metadata.xml",
	}
	if !reflect.DeepEqual(requestedURLs, wantURLs) {
		t.Fatalf("plugin metadata requests were %#v, want ordered continuation %#v", requestedURLs, wantURLs)
	}
	if len(model.Build.Plugins.Plugin) != len(plugins) {
		t.Fatalf("plugin population changed from %d to %d entries", len(plugins), len(model.Build.Plugins.Plugin))
	}
	if actual := model.Build.Plugins.Plugin[4].Version; actual != "1.1.0" {
		t.Fatalf("later plugin version was %q, want 1.1.0 after every earlier failure", actual)
	}
}

func captureMavenPartialFailureLogs(t *testing.T) *recordingMavenLogHook {
	t.Helper()
	previousLogger := log
	hook := &recordingMavenLogHook{}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(io.Discard)
	logger.AddHook(hook)
	SetLogger(logger)
	t.Cleanup(func() { SetLogger(previousLogger) })
	return hook
}

func assertMavenPartialFailureLogs(t *testing.T, hook *recordingMavenLogHook, want []mavenPartialFailureLog) {
	t.Helper()
	if len(want) == 0 {
		t.Fatal("Maven partial-failure log expectation population is empty")
	}
	if len(hook.entries) == 0 {
		t.Fatal("Maven partial-failure log population is empty")
	}
	searchFrom := 0
	for _, expected := range want {
		found := -1
		matches := 0
		for index, entry := range hook.entries {
			if entry.Level == expected.level && entry.Message == expected.message {
				matches++
				if index >= searchFrom && found < 0 {
					found = index
				}
			}
		}
		if matches != 1 || found < 0 {
			actual := make([]mavenPartialFailureLog, 0, len(hook.entries))
			for _, entry := range hook.entries {
				actual = append(actual, mavenPartialFailureLog{level: entry.Level, message: entry.Message})
			}
			t.Fatalf("Maven exact log level=%s message=%q occurred %d times after index %d in %#v",
				expected.level, expected.message, matches, searchFrom, actual)
		}
		searchFrom = found + 1
	}
	if searchFrom == 0 {
		t.Fatal("Maven partial-failure assertions did not inspect any log entries")
	}
}
