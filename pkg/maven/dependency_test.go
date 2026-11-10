package maven

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/devdimensionlab/plybuild/internal/testutil"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/file"
)

func TestUpgradeDependency(t *testing.T) {
	projectDir := copyMavenFixtureDir(t, "test/dependency")
	pomFile := filepath.Join(projectDir, "pom.xml")
	pomContents, err := os.ReadFile(pomFile)
	if err != nil {
		t.Fatalf("read dependency POM fixture: %v", err)
	}
	const currentVersion = "<flyway.version>8.0.3</flyway.version>"
	const oldVersion = "<flyway.version>7.15.0</flyway.version>"
	if matches := strings.Count(string(pomContents), currentVersion); matches != 1 {
		t.Fatalf("dependency fixture contained %d current versions, want exactly 1", matches)
	}
	pomContents = []byte(strings.Replace(string(pomContents), currentVersion, oldVersion, 1))
	if strings.Count(string(pomContents), oldVersion) != 1 || strings.Contains(string(pomContents), currentVersion) {
		t.Fatal("dependency fixture was not downgraded before exercising the upgrade")
	}
	writeMavenFixtureFile(t, pomFile, pomContents)

	projectConfig, err := config.InitProjectConfigurationFromDir(projectDir)
	if err != nil {
		t.Fatalf("load project configuration fixture: %v", err)
	}
	model, err := pom.GetModelFrom(pomFile)
	if err != nil {
		t.Fatalf("load dependency POM fixture: %v", err)
	}
	project := config.Project{
		Path:   projectDir,
		Config: projectConfig,
		Type: config.MavenProject{
			PomFile:  pomFile,
			PomModel: model,
		},
	}

	deps := model.Dependencies.Dependency
	if len(deps) != 1 || deps[0].GroupId != "org.flywaydb" || deps[0].ArtifactId != "flyway-core" {
		t.Fatalf("dependency population was %+v, want one org.flywaydb:flyway-core dependency", deps)
	}

	const repositoryURL = "https://packages.example.test/releases"
	wantMetadataURL := repositoryURL + "/org/flywaydb/flyway-core/maven-metadata.xml"
	var requestedURLs []string
	stubHTTPTransport(t, func(request *http.Request) (*http.Response, error) {
		requestedURLs = append(requestedURLs, request.URL.String())
		if request.Method != http.MethodGet {
			return nil, fmt.Errorf("unexpected method %s", request.Method)
		}
		metadata := `<metadata><groupId>org.flywaydb</groupId><artifactId>flyway-core</artifactId><versioning><release>9.2.1</release><versions><version>9.2.1</version></versions></versioning></metadata>`
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(metadata)),
		}, nil
	})

	repo := Repository{Url: repositoryURL}
	repo.upgradeDependencies(model, deps, project.Config.Settings, func(groupId string) bool {
		return true
	}, model.SetDependencyVersion)
	if len(requestedURLs) != 1 || requestedURLs[0] != wantMetadataURL {
		t.Fatalf("metadata requests were %v, want [%s]", requestedURLs, wantMetadataURL)
	}

	if err := project.SortAndWritePom(); err != nil {
		t.Fatalf("write upgraded dependency POM: %v", err)
	}

	originPomFile := "test/dependency/origin.pom.xml"
	equal, err := file.Equal(originPomFile, pomFile)
	if err != nil {
		t.Fatalf("compare upgraded dependency POM: %v", err)
	}
	if !equal {
		t.Errorf("%s is not equal to %s", originPomFile, pomFile)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func stubHTTPTransport(t *testing.T, transport roundTripFunc) {
	t.Helper()
	previous := http.DefaultTransport
	http.DefaultTransport = transport
	t.Cleanup(func() {
		http.DefaultTransport = previous
	})
}

func copyMavenFixtureDir(t *testing.T, source string) string {
	t.Helper()
	target := filepath.Join(t.TempDir(), filepath.Base(source))
	if err := testutil.CopyFSOutsideWorkingTree(target, os.DirFS(source)); err != nil {
		t.Fatalf("copy Maven fixture directory %s: %v", source, err)
	}
	return target
}

func writeMavenFixtureFile(t *testing.T, target string, contents []byte) {
	t.Helper()
	if err := testutil.WriteFileOutsideWorkingTree(target, contents, 0o644); err != nil {
		t.Fatalf("write fixture %s: %v", target, err)
	}
}
