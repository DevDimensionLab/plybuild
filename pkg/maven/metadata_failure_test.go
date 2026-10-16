package maven

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/sirupsen/logrus"
)

func TestTargetedDependencySkipsPomWithoutDependencySections(t *testing.T) {
	for _, model := range []*pom.Model{{}, {DependencyManagement: &pom.DependencyManagement{}}} {
		hook := captureMavenPartialFailureLogs(t)
		stubHTTPTransport(t, func(*http.Request) (*http.Response, error) {
			t.Fatal("dependency-free POM must not perform metadata requests")
			return nil, nil
		})
		project := config.Project{Type: config.MavenProject{PomModel: model}}
		if err := UpgradeDependency("com.example", "client")(Repository{Url: "https://packages.example.test/releases"}, project); err != nil {
			t.Fatalf("dependency-free POM returned a new job error: %v", err)
		}
		for _, entry := range hook.entries {
			if entry.Level <= logrus.WarnLevel {
				t.Errorf("dependency-free POM produced a warning: %s", entry.Message)
			}
		}
	}
}

func TestMetadataFailureDiagnosticsPreserveContextAndRedactUntrustedDetails(t *testing.T) {
	const repositoryURL = "https://url-user:url-password@packages.example.test/releases?token=query-secret#fragment-secret"
	cases := []struct {
		name   string
		status int
		body   string
		err    error
		reason string
	}{
		{name: "unauthorized", status: 401, reason: "HTTP 401 Unauthorized"},
		{name: "forbidden", status: 403, reason: "HTTP 403 Forbidden"},
		{name: "not found", status: 404, reason: "HTTP 404 Not Found"},
		{name: "unavailable", status: 503, reason: "HTTP 503 Service Unavailable"},
		{name: "transport", err: errors.New("untrusted transport query-secret url-password"), reason: "transport error"},
		{name: "xml", status: 200, body: `<metadata><response-secret></metadata>`, reason: "invalid XML"},
		{name: "invalid release", status: 200, body: `<metadata><versioning><release>response-secret</release><versions><version>5.12.0</version></versions></versioning></metadata>`, reason: "invalid version"},
		{name: "invalid list", status: 200, body: `<metadata><versioning><versions><version>5.12.0</version><version>response-secret</version></versions></versioning></metadata>`, reason: "invalid version"},
		{name: "no release", status: 200, body: `<metadata><versioning><latest>5.12.0</latest><versions><version>5.13.0-RC1</version></versions></versioning></metadata>`, reason: "suitable release"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hook := captureMavenPartialFailureLogs(t)
			stubHTTPTransport(t, func(request *http.Request) (*http.Response, error) {
				user, password, ok := request.BasicAuth()
				if !ok || user != "auth-user" || password != "auth-password" {
					t.Fatal("metadata request did not preserve configured Basic Auth")
				}
				if tc.err != nil {
					return nil, tc.err
				}
				return &http.Response{StatusCode: tc.status, Status: "untrusted response-secret", Body: io.NopCloser(strings.NewReader(tc.body))}, nil
			})
			repository := Repository{Url: repositoryURL, Auth: &RepositoryAuth{Username: "auth-user", Password: "auth-password"}}
			_, err := repository.latestRelease("com.requested", "requested-client")
			if err == nil {
				t.Fatal("expected metadata failure")
			}
			for _, want := range []string{"com.requested:requested-client", "https://packages.example.test/releases", tc.reason} {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("diagnostic %q missing %q", err.Error(), want)
				}
			}
			if tc.err != nil && !errors.Is(err, tc.err) {
				t.Error("transport cause identity was lost")
			}
			allMessages := err.Error()
			for _, entry := range hook.entries {
				allMessages += "\n" + entry.Message
				if entry.Level == logrus.WarnLevel {
					t.Error("lookup logged its own warning; caller owns reporting")
				}
			}
			for _, secret := range []string{"url-user", "url-password", "query-secret", "fragment-secret", "auth-user", "auth-password", "response-secret", "Authorization"} {
				if strings.Contains(allMessages, secret) {
					t.Errorf("diagnostics exposed %q: %s", secret, allMessages)
				}
			}
		})
	}
}

func TestTargetedDependencyRetainsFailuresAcrossPomSections(t *testing.T) {
	for _, failureCount := range []int{1, 2} {
		t.Run(string(rune('0'+failureCount))+" failures", func(t *testing.T) {
			hook := captureMavenPartialFailureLogs(t)
			log.(*logrus.Logger).SetLevel(logrus.InfoLevel)
			dep := pom.Dependency{GroupId: "com.example", ArtifactId: "client", Version: "5.11.4"}
			model := &pom.Model{Dependencies: &pom.Dependencies{Dependency: []pom.Dependency{dep}}, DependencyManagement: &pom.DependencyManagement{Dependencies: &pom.Dependencies{Dependency: []pom.Dependency{dep}}}}
			requests := 0
			stubHTTPTransport(t, func(*http.Request) (*http.Response, error) {
				requests++
				body := `<metadata><versioning><release>5.11.4</release></versioning></metadata>`
				if requests <= failureCount {
					body = `<metadata><versioning><release>broken</release></versioning></metadata>`
				}
				return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
			})
			project := config.Project{Type: config.MavenProject{PomModel: model}}
			if err := UpgradeDependency(dep.GroupId, dep.ArtifactId)(Repository{Url: "https://packages.example.test/releases"}, project); err != nil {
				t.Fatalf("per-lookup errors should be reported individually: %v", err)
			}
			warnings := 0
			for _, entry := range hook.entries {
				if entry.Level == logrus.WarnLevel {
					warnings++
					if !strings.Contains(entry.Message, "com.example:client") || !strings.Contains(entry.Message, "invalid version") {
						t.Errorf("unexpected warning %q", entry.Message)
					}
				}
			}
			if requests != 2 || warnings != failureCount {
				t.Errorf("requests=%d warnings=%d; want 2 requests and %d warnings", requests, warnings, failureCount)
			}
			if model.Dependencies.Dependency[0].Version != dep.Version || model.DependencyManagement.Dependencies.Dependency[0].Version != dep.Version {
				t.Fatal("failed or equal-version lookup changed a dependency")
			}
		})
	}
}

func TestUpgradeDependenciesWarnsOnceForFailedReleaseAndContinues(t *testing.T) {
	const repositoryURL = "https://packages.example.test/releases"
	model := &pom.Model{Dependencies: &pom.Dependencies{Dependency: []pom.Dependency{
		{GroupId: "com.example", ArtifactId: "no-release", Version: "5.11.4"},
		{GroupId: "com.example", ArtifactId: "later-success", Version: "5.11.4"},
	}}}
	hook := captureMavenPartialFailureLogs(t)
	log.(*logrus.Logger).SetLevel(logrus.InfoLevel)
	var requested []string
	stubHTTPTransport(t, func(request *http.Request) (*http.Response, error) {
		requested = append(requested, request.URL.Path)
		metadata := `<metadata><versioning><release>5.12.0</release></versioning></metadata>`
		if strings.Contains(request.URL.Path, "/no-release/") {
			metadata = `<metadata><versioning><release>5.13.0-SNAPSHOT</release><versions><version>5.13.0-SNAPSHOT</version></versions></versioning></metadata>`
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(metadata))}, nil
	})
	var actions []string
	(Repository{Url: repositoryURL}).upgradeDependencies(model, model.Dependencies.Dependency, config.ProjectSettings{}, func(string) bool { return true }, func(dep pom.Dependency, version string) error {
		actions = append(actions, dep.ArtifactId+"="+version)
		return nil
	})
	if len(requested) != 2 || len(actions) != 1 || actions[0] != "later-success=5.12.0" {
		t.Fatalf("requests=%v actions=%v; want failed item unchanged and later upgrade", requested, actions)
	}
	var warnings []string
	for _, entry := range hook.entries {
		if entry.Level == logrus.WarnLevel {
			warnings = append(warnings, entry.Message)
		}
	}
	if len(warnings) != 1 {
		t.Fatalf("got warnings %v, want exactly one visible release failure", warnings)
	}
	for _, text := range []string{"com.example:no-release", repositoryURL, "suitable release"} {
		if !strings.Contains(warnings[0], text) {
			t.Errorf("warning %q missing %q", warnings[0], text)
		}
	}
}
