package cmd

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
)

const upgradeAcceptanceUser = "acceptance-fixture-user"
const upgradeAcceptancePassword = "acceptance-fixture-password"

// Each CLI invocation gets a new process: Cobra flags, profile discovery, log
// collection (also used by the stealth writer), and exit handling all run as
// they do at the executable boundary. The only substituted effect is an HTTP
// guard, which rejects and records traffic outside the loopback fixture.
func TestUpgradeMetadataAcceptanceCLIProcess(t *testing.T) {
	if os.Getenv("PLY_UPGRADE_ACCEPTANCE_CHILD") != "1" {
		return
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	http.DefaultTransport = upgradeAcceptanceTransport{
		host:      os.Getenv("PLY_UPGRADE_ACCEPTANCE_HOST"),
		requests:  os.Getenv("PLY_UPGRADE_ACCEPTANCE_REQUESTS"),
		transport: transport,
	}
	var args []string
	if err := json.Unmarshal([]byte(os.Getenv("PLY_UPGRADE_ACCEPTANCE_ARGS")), &args); err != nil {
		t.Fatal(err)
	}
	RootCmd.SetArgs(args)
	if err := ExecuteE(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

type upgradeAcceptanceTransport struct {
	host      string
	requests  string
	transport http.RoundTripper
}

func (transport upgradeAcceptanceTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	file, err := os.OpenFile(transport.requests, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	_, err = fmt.Fprintln(file, request.URL.String())
	closeErr := file.Close()
	if err != nil {
		return nil, err
	}
	if closeErr != nil {
		return nil, closeErr
	}
	if request.URL.Host != transport.host {
		return nil, fmt.Errorf("acceptance fixture rejected unexpected repository host %s", request.URL.Host)
	}
	return transport.transport.RoundTrip(request)
}

type upgradeAcceptanceResponse struct {
	status     int
	body       string
	disconnect bool
}

type upgradeAcceptanceRequest struct {
	path     string
	username string
	password string
}

type upgradeAcceptanceFixture struct {
	t         *testing.T
	home      string
	server    *httptest.Server
	mu        sync.Mutex
	requests  []upgradeAcceptanceRequest
	responses map[string][]upgradeAcceptanceResponse
}

func newUpgradeAcceptanceFixture(t *testing.T, responses map[string][]upgradeAcceptanceResponse) *upgradeAcceptanceFixture {
	t.Helper()
	fixture := &upgradeAcceptanceFixture{t: t, home: t.TempDir(), responses: responses}
	fixture.server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		username, password, _ := request.BasicAuth()
		fixture.mu.Lock()
		fixture.requests = append(fixture.requests, upgradeAcceptanceRequest{path: request.URL.Path, username: username, password: password})
		responses := fixture.responses[request.URL.Path]
		response := upgradeAcceptanceResponse{status: 599, body: "unexpected fixture request"}
		if len(responses) > 0 {
			response = responses[0]
			fixture.responses[request.URL.Path] = responses[1:]
		}
		fixture.mu.Unlock()
		if response.disconnect {
			connection, _, err := writer.(http.Hijacker).Hijack()
			if err == nil {
				_ = connection.Close()
			}
			return
		}
		status := response.status
		if status == 0 {
			status = http.StatusOK
		}
		writer.Header().Set("Content-Type", "application/xml")
		writer.WriteHeader(status)
		_, _ = writer.Write([]byte(response.body))
	}))
	t.Cleanup(fixture.server.Close)
	profile := filepath.Join(fixture.home, ".ply", "profiles", "acceptance")
	writeUpgradeAcceptanceFile(t, filepath.Join(fixture.home, ".ply", "profiles", ".active_profile"), "acceptance\n")
	// An explicit Nexus setting makes the production selector choose this
	// repository before touching os/user.Current().HomeDir Maven settings.
	writeUpgradeAcceptanceFile(t, filepath.Join(profile, "local-config.yaml"), fmt.Sprintf("nexus:\n  url: %q\n  username: %q\n  password: %q\n", fixture.server.URL+"/selected/releases", upgradeAcceptanceUser, upgradeAcceptancePassword))
	writeUpgradeAcceptanceFile(t, filepath.Join(profile, "cloud-config", "project-defaults.json"), `{ "settings": {} }`)
	return fixture
}

func writeUpgradeAcceptanceFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}
}

func (fixture *upgradeAcceptanceFixture) project(name, body, settings string) string {
	fixture.t.Helper()
	directory := filepath.Join(fixture.home, "projects", name)
	writeUpgradeAcceptanceFile(fixture.t, filepath.Join(directory, "pom.xml"), upgradeAcceptancePOM(body))
	if settings == "" {
		settings = "{}"
	}
	writeUpgradeAcceptanceFile(fixture.t, filepath.Join(directory, "ply.json"), `{"groupId":"com.example","artifactId":"sample","name":"sample","language":"java","package":"com.example","applicationName":"Application","settings":`+settings+`}`)
	return directory
}

func upgradeAcceptancePOM(body string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
  <modelVersion>4.0.0</modelVersion>
  <groupId>com.example</groupId>
  <artifactId>sample</artifactId>
  <version>1.0.0</version>
  <name>sample</name>
` + body + `
  <repositories>
    <repository><id>unused-pom-repository</id><url>https://unused-pom.example.invalid/packages</url></repository>
  </repositories>
</project>
`
}

func upgradeAcceptanceDependency(group, artifact, version string) string {
	return fmt.Sprintf("    <dependency>\n      <groupId>%s</groupId>\n      <artifactId>%s</artifactId>\n      <version>%s</version>\n    </dependency>\n", group, artifact, version)
}

func upgradeAcceptanceMetadataPath(group, artifact string) string {
	return "/selected/releases/" + strings.ReplaceAll(group, ".", "/") + "/" + artifact + "/maven-metadata.xml"
}

const upgradeAcceptanceFallback = `<metadata><versioning><latest>5.13.0-SNAPSHOT</latest><versions><version>5.9.0</version><version>5.12.0</version><version>5.11.4</version><version>5.13.0-SNAPSHOT</version><version>5.14.0-RC1</version></versions></versioning></metadata>`
const upgradeAcceptanceNoRelease = `<metadata><versioning><release>5.13.0-SNAPSHOT</release><versions><version>5.13.0-SNAPSHOT</version><version>5.14.0-RC1</version></versions></versioning></metadata>`

type upgradeAcceptanceLog struct {
	Level   string `json:"level"`
	Message string `json:"msg"`
}

func (fixture *upgradeAcceptanceFixture) run(target string, args ...string) []upgradeAcceptanceLog {
	return fixture.runExecutable("", target, args...)
}

func (fixture *upgradeAcceptanceFixture) runExecutable(binary, target string, args ...string) []upgradeAcceptanceLog {
	fixture.t.Helper()
	executable, err := os.Executable()
	if err != nil {
		fixture.t.Fatal(err)
	}
	requestFile := filepath.Join(fixture.t.TempDir(), "requests.txt")
	parsed, err := url.Parse(fixture.server.URL)
	if err != nil {
		fixture.t.Fatal(err)
	}
	args = append(args, "--target", target, "--json")
	encoded, err := json.Marshal(args)
	if err != nil {
		fixture.t.Fatal(err)
	}
	processContext, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	command := exec.CommandContext(processContext, executable, "-test.run=^TestUpgradeMetadataAcceptanceCLIProcess$")
	if binary != "" {
		command = exec.CommandContext(processContext, binary, args...)
	}
	command.Dir = fixture.home
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if key == "HOME" || strings.HasPrefix(key, "PLY_") || strings.HasSuffix(strings.ToUpper(key), "_PROXY") {
			continue
		}
		command.Env = append(command.Env, entry)
	}
	command.Env = append(command.Env, "HOME="+fixture.home,
		"PLY_UPGRADE_ACCEPTANCE_CHILD=1",
		"PLY_UPGRADE_ACCEPTANCE_HOST="+parsed.Host,
		"PLY_UPGRADE_ACCEPTANCE_REQUESTS="+requestFile,
		"PLY_UPGRADE_ACCEPTANCE_ARGS="+string(encoded))
	if binary != "" {
		// The built executable cannot use the in-process guard. A local proxy
		// captures and rejects any POM/Central fallback without forwarding it.
		// Go sends selected loopback requests directly, so normal HTTP is real.
		command.Env = append(command.Env, "HTTP_PROXY="+fixture.server.URL, "HTTPS_PROXY="+fixture.server.URL, "NO_PROXY=", "http_proxy="+fixture.server.URL, "https_proxy="+fixture.server.URL, "no_proxy=")
	}
	output, err := command.CombinedOutput()
	if err != nil {
		fixture.t.Fatalf("CLI command %v exited unsuccessfully: %v\n%s", args, err, output)
	}
	for _, forbidden := range []string{upgradeAcceptanceUser, upgradeAcceptancePassword, "acceptance-url-user", "acceptance-url-password", "Authorization", base64.StdEncoding.EncodeToString([]byte(upgradeAcceptanceUser + ":" + upgradeAcceptancePassword))} {
		if bytes.Contains(output, []byte(forbidden)) {
			fixture.t.Errorf("CLI output disclosed synthetic credentials or header %q: %s", forbidden, output)
		}
	}
	var logs []upgradeAcceptanceLog
	for _, line := range bytes.Split(bytes.TrimSpace(output), []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		var entry upgradeAcceptanceLog
		if err := json.Unmarshal(line, &entry); err != nil {
			fixture.t.Fatalf("CLI output was not a JSON log record: %s (%v)", line, err)
		}
		logs = append(logs, entry)
	}
	if binary == "" {
		requests, err := os.ReadFile(requestFile)
		if err != nil {
			fixture.t.Fatal(err)
		}
		for _, requested := range strings.Split(strings.TrimSpace(string(requests)), "\n") {
			requestURL, err := url.Parse(requested)
			if err != nil || requestURL.Host != parsed.Host {
				fixture.t.Errorf("unexpected repository request (no POM/Central fallback allowed): %q", requested)
			}
		}
	}
	return logs
}

func (fixture *upgradeAcceptanceFixture) assertRequests(paths ...string) {
	fixture.t.Helper()
	fixture.mu.Lock()
	defer fixture.mu.Unlock()
	var got []string
	for _, request := range fixture.requests {
		got = append(got, request.path)
		if request.username != upgradeAcceptanceUser || request.password != upgradeAcceptancePassword {
			fixture.t.Errorf("request %s did not use configured synthetic Basic Auth", request.path)
		}
	}
	if !reflect.DeepEqual(got, paths) {
		fixture.t.Errorf("repository requests = %v, want %v", got, paths)
	}
}

func upgradeAcceptanceModel(t *testing.T, directory string) *pom.Model {
	t.Helper()
	model, err := pom.GetModelFrom(filepath.Join(directory, "pom.xml"))
	if err != nil {
		t.Fatal(err)
	}
	return model
}

func assertUpgradeAcceptanceVersion(t *testing.T, model *pom.Model, artifact, expected string, managed bool) {
	t.Helper()
	var dependencies []pom.Dependency
	if managed {
		if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
			dependencies = model.DependencyManagement.Dependencies.Dependency
		}
	} else if model.Dependencies != nil {
		dependencies = model.Dependencies.Dependency
	}
	for _, dependency := range dependencies {
		if dependency.ArtifactId != artifact {
			continue
		}
		version, err := model.GetDependencyVersion(dependency)
		if err != nil || version != expected {
			t.Errorf("%s version = %q (%v), want %q", artifact, version, err, expected)
		}
		return
	}
	t.Errorf("dependency %s missing (managed=%v)", artifact, managed)
}

func assertUpgradeAcceptanceFailure(t *testing.T, logs []upgradeAcceptanceLog, artifact, repository, cause string) {
	t.Helper()
	count := 0
	for _, entry := range logs {
		if strings.Contains(entry.Message, artifact) && entry.Level == "warning" {
			count++
			if !strings.Contains(entry.Message, repository) || !strings.Contains(strings.ToLower(entry.Message), strings.ToLower(cause)) {
				t.Errorf("warning lacks requested repository/cause (%s): %+v", cause, entry)
			}
		}
		if strings.Contains(entry.Message, artifact) && (strings.Contains(entry.Message, "outdated") || strings.Contains(entry.Message, "latest version") || strings.Contains(entry.Message, "up to date")) {
			t.Errorf("failed lookup was described as an upgrade/success: %+v", entry)
		}
	}
	if count != 1 {
		t.Errorf("got %d warnings for %s, want exactly one; logs: %+v", count, artifact, logs)
	}
}

func TestUpgradeMetadataAcceptanceReportsNoReleaseAndContinues(t *testing.T) {
	failed := upgradeAcceptanceMetadataPath("com.example", "no-release")
	later := upgradeAcceptanceMetadataPath("com.example", "later")
	fixture := newUpgradeAcceptanceFixture(t, map[string][]upgradeAcceptanceResponse{
		failed: {{body: upgradeAcceptanceNoRelease}},
		later:  {{body: `<metadata><versioning><release>5.12.0</release></versioning></metadata>`}},
	})
	directory := fixture.project("batch", "<dependencies>\n"+upgradeAcceptanceDependency("com.example", "no-release", "5.11.4")+upgradeAcceptanceDependency("com.example", "later", "5.11.4")+"</dependencies>", "")
	logs := fixture.run(directory, "upgrade", "2party", "--stealth")
	assertUpgradeAcceptanceFailure(t, logs, "com.example:no-release", fixture.server.URL+"/selected/releases", "suitable release")
	model := upgradeAcceptanceModel(t, directory)
	assertUpgradeAcceptanceVersion(t, model, "no-release", "5.11.4", false)
	assertUpgradeAcceptanceVersion(t, model, "later", "5.12.0", false)
	fixture.assertRequests(failed, later)
}

func TestUpgradeMetadataAcceptanceMainJourneyWritesAndDryRun(t *testing.T) {
	for _, mode := range []string{"normal", "stealth", "dry-run"} {
		for _, layout := range []string{"direct", "property", "dependencyManagement"} {
			t.Run(mode+"/"+layout, func(t *testing.T) {
				failed := upgradeAcceptanceMetadataPath("com.example", "no-release")
				client := upgradeAcceptanceMetadataPath("com.example", "client")
				fixture := newUpgradeAcceptanceFixture(t, map[string][]upgradeAcceptanceResponse{
					failed: {{body: upgradeAcceptanceNoRelease}}, client: {{body: upgradeAcceptanceFallback}},
				})
				version, properties := "5.11.4", ""
				if layout == "property" {
					version = "${client.version}"
					properties = "<properties>\n  <client.version>5.11.4</client.version>\n</properties>\n"
				}
				body := "<dependencies>\n" + upgradeAcceptanceDependency("com.example", "no-release", "5.11.4") + upgradeAcceptanceDependency("com.example", "client", version) + "</dependencies>\n"
				if layout == "dependencyManagement" {
					body = "<dependencyManagement>\n" + body + "</dependencyManagement>\n"
				}
				directory := fixture.project("main", properties+body, "")
				before := readUpgradeAcceptancePOM(t, directory)
				args := []string{"upgrade", "2party"}
				if mode != "normal" {
					args = append(args, "--"+mode)
				}
				logs := fixture.run(directory, args...)
				assertUpgradeAcceptanceFailure(t, logs, "com.example:no-release", fixture.server.URL+"/selected/releases", "suitable release")
				assertUpgradeAcceptanceWarningCount(t, logs, 1)
				foundUpgrade := false
				for _, entry := range logs {
					if entry.Level == "info" && strings.Contains(entry.Message, "com.example:client") && strings.Contains(entry.Message, "5.11.4") && strings.Contains(entry.Message, "5.12.0") {
						foundUpgrade = true
					}
				}
				if !foundUpgrade {
					t.Errorf("successful fallback was not reported: %+v", logs)
				}
				expected := "5.12.0"
				if mode == "dry-run" {
					expected = "5.11.4"
					if after := readUpgradeAcceptancePOM(t, directory); !bytes.Equal(before, after) {
						t.Error("dry-run changed POM bytes")
					}
				}
				model := upgradeAcceptanceModel(t, directory)
				assertUpgradeAcceptanceVersion(t, model, "client", expected, layout == "dependencyManagement")
				assertUpgradeAcceptanceVersion(t, model, "no-release", "5.11.4", layout == "dependencyManagement")
				if layout == "property" && !bytes.Contains(readUpgradeAcceptancePOM(t, directory), []byte("${client.version}")) {
					t.Error("property reference was replaced rather than its value")
				}
				fixture.assertRequests(failed, client)
			})
		}
	}
}

func TestUpgradeMetadataAcceptanceBatchReleaseFailuresAndAllFailed(t *testing.T) {
	for _, allFailed := range []bool{false, true} {
		t.Run(fmt.Sprintf("allFailed=%v", allFailed), func(t *testing.T) {
			failures := []struct{ artifact, body, cause string }{
				{"no-release", upgradeAcceptanceNoRelease, "suitable release"},
				{"invalid-xml", `<metadata><sensitive-response-marker></metadata>`, "XML"},
				{"invalid-release", `<metadata><versioning><release>broken</release><versions><version>5.12.0</version></versions></versioning></metadata>`, "version"},
				{"invalid-list", `<metadata><versioning><versions><version>5.12.0</version><version>broken</version></versions></versioning></metadata>`, "version"},
			}
			responses := map[string][]upgradeAcceptanceResponse{}
			var paths []string
			body := "<dependencies>\n"
			for _, failure := range failures {
				path := upgradeAcceptanceMetadataPath("com.example", failure.artifact)
				responses[path] = []upgradeAcceptanceResponse{{body: failure.body}}
				paths = append(paths, path)
				body += upgradeAcceptanceDependency("com.example", failure.artifact, "5.11.4")
			}
			if !allFailed {
				path := upgradeAcceptanceMetadataPath("com.example", "later")
				responses[path] = []upgradeAcceptanceResponse{{body: upgradeAcceptanceFallback}}
				paths = append(paths, path)
				body += upgradeAcceptanceDependency("com.example", "later", "5.11.4")
			}
			fixture := newUpgradeAcceptanceFixture(t, responses)
			directory := fixture.project("batch", body+"</dependencies>\n", "")
			logs := fixture.run(directory, "upgrade", "2party", "--stealth")
			model := upgradeAcceptanceModel(t, directory)
			for _, failure := range failures {
				assertUpgradeAcceptanceFailure(t, logs, "com.example:"+failure.artifact, fixture.server.URL+"/selected/releases", failure.cause)
				assertUpgradeAcceptanceVersion(t, model, failure.artifact, "5.11.4", false)
			}
			assertUpgradeAcceptanceWarningCount(t, logs, len(failures))
			if !allFailed {
				assertUpgradeAcceptanceVersion(t, model, "later", "5.12.0", false)
			}
			fixture.assertRequests(paths...)
		})
	}
}

func TestUpgradeMetadataAcceptanceHTTPFailuresContinue(t *testing.T) {
	cases := []struct {
		name, cause string
		response    upgradeAcceptanceResponse
	}{
		{"unauthorized", "HTTP 401", upgradeAcceptanceResponse{status: 401, body: "sensitive-response-marker"}},
		{"forbidden", "HTTP 403", upgradeAcceptanceResponse{status: 403, body: "sensitive-response-marker"}},
		{"not-found", "HTTP 404", upgradeAcceptanceResponse{status: 404, body: "sensitive-response-marker"}},
		{"unavailable", "HTTP 503", upgradeAcceptanceResponse{status: 503, body: "sensitive-response-marker"}},
		{"transport", "transport", upgradeAcceptanceResponse{disconnect: true}},
		{"xml", "XML", upgradeAcceptanceResponse{body: `<metadata><sensitive-response-marker></metadata>`}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			failed, later := upgradeAcceptanceMetadataPath("com.example", "failed"), upgradeAcceptanceMetadataPath("com.example", "later")
			fixture := newUpgradeAcceptanceFixture(t, map[string][]upgradeAcceptanceResponse{failed: {test.response}, later: {{body: upgradeAcceptanceFallback}}})
			directory := fixture.project("http", "<dependencies>\n"+upgradeAcceptanceDependency("com.example", "failed", "5.11.4")+upgradeAcceptanceDependency("com.example", "later", "5.11.4")+"</dependencies>\n", "")
			logs := fixture.run(directory, "upgrade", "2party")
			assertUpgradeAcceptanceFailure(t, logs, "com.example:failed", fixture.server.URL+"/selected/releases", test.cause)
			assertUpgradeAcceptanceWarningCount(t, logs, 1)
			for _, entry := range logs {
				if strings.Contains(entry.Message, "sensitive-response-marker") {
					t.Errorf("raw response content leaked in log: %+v", entry)
				}
			}
			model := upgradeAcceptanceModel(t, directory)
			assertUpgradeAcceptanceVersion(t, model, "failed", "5.11.4", false)
			assertUpgradeAcceptanceVersion(t, model, "later", "5.12.0", false)
			fixture.assertRequests(failed, later)
		})
	}
}

func TestUpgradeMetadataAcceptancePreservesDependencySelectionAndMaximum(t *testing.T) {
	for _, party := range []string{"2party", "3party"} {
		t.Run(party, func(t *testing.T) {
			selected, excluded := "com.example", "org.external"
			if party == "3party" {
				selected, excluded = excluded, selected
			}
			responses := map[string][]upgradeAcceptanceResponse{}
			var paths []string
			for _, artifact := range []string{"equal", "newer", "maximum", "later"} {
				path := upgradeAcceptanceMetadataPath(selected, artifact)
				responses[path] = []upgradeAcceptanceResponse{{body: upgradeAcceptanceFallback}}
				paths = append(paths, path)
			}
			fixture := newUpgradeAcceptanceFixture(t, responses)
			body := "<properties><revision>5.11.4</revision></properties>\n<dependencies>\n" +
				upgradeAcceptanceDependency(selected, "ignored", "5.11.4") +
				upgradeAcceptanceDependency(excluded, "wrong-party", "5.11.4") +
				upgradeAcceptanceDependency(selected, "project-version", "${project.version}") +
				upgradeAcceptanceDependency(selected, "revision-version", "${revision}") +
				upgradeAcceptanceDependency(selected, "equal", "5.12.0") +
				upgradeAcceptanceDependency(selected, "newer", "5.15.0") +
				upgradeAcceptanceDependency(selected, "maximum", "5.11.4") +
				upgradeAcceptanceDependency(selected, "later", "5.11.4") + "</dependencies>\n"
			settings := fmt.Sprintf(`{"disableUpgradesFor":[{"groupId":%q,"artifactId":"ignored"}],"maxVersionForDependencies":[{"groupId":%q,"artifactId":"maximum","maxVersion":"5.11.5"}]}`, selected, selected)
			directory := fixture.project("selection", body, settings)
			logs := fixture.run(directory, "upgrade", party)
			model := upgradeAcceptanceModel(t, directory)
			for artifact, version := range map[string]string{"ignored": "5.11.4", "wrong-party": "5.11.4", "equal": "5.12.0", "newer": "5.15.0", "maximum": "5.11.5", "later": "5.12.0"} {
				assertUpgradeAcceptanceVersion(t, model, artifact, version, false)
			}
			for _, dependency := range model.Dependencies.Dependency {
				if dependency.ArtifactId == "project-version" && dependency.Version != "${project.version}" {
					t.Error("project.version reference changed")
				}
				if dependency.ArtifactId == "revision-version" && dependency.Version != "${revision}" {
					t.Error("revision reference changed")
				}
			}
			// The existing configured maximum notice remains a warning; none of
			// the intentionally excluded or already current items gains one.
			assertUpgradeAcceptanceWarningCount(t, logs, 1)
			for _, entry := range logs {
				if entry.Level == "warning" && (!strings.Contains(entry.Message, selected+":maximum") || !strings.Contains(entry.Message, "held back")) {
					t.Errorf("unexpected selection warning: %+v", entry)
				}
			}
			fixture.assertRequests(paths...)
		})
	}
}

func TestUpgradeMetadataAcceptanceThirdPartyFailures(t *testing.T) {
	for _, managed := range []bool{false, true} {
		t.Run(fmt.Sprintf("managed=%v", managed), func(t *testing.T) {
			failed, later := upgradeAcceptanceMetadataPath("org.external", "failed"), upgradeAcceptanceMetadataPath("org.external", "later")
			fixture := newUpgradeAcceptanceFixture(t, map[string][]upgradeAcceptanceResponse{failed: {{body: upgradeAcceptanceNoRelease}}, later: {{body: upgradeAcceptanceFallback}}})
			body := "<dependencies>\n" + upgradeAcceptanceDependency("org.external", "failed", "5.11.4")
			if managed {
				body += "</dependencies>\n<dependencyManagement>\n<dependencies>\n"
			}
			body += upgradeAcceptanceDependency("org.external", "later", "5.11.4") + "</dependencies>\n"
			if managed {
				body += "</dependencyManagement>\n"
			}
			directory := fixture.project("third-party", body, "")
			logs := fixture.run(directory, "upgrade", "3party", "--stealth")
			assertUpgradeAcceptanceFailure(t, logs, "org.external:failed", fixture.server.URL+"/selected/releases", "suitable release")
			assertUpgradeAcceptanceWarningCount(t, logs, 1)
			model := upgradeAcceptanceModel(t, directory)
			assertUpgradeAcceptanceVersion(t, model, "failed", "5.11.4", false)
			assertUpgradeAcceptanceVersion(t, model, "later", "5.12.0", managed)
			fixture.assertRequests(failed, later)
		})
	}
}

func TestUpgradeMetadataAcceptanceTargetedDependencyKeepsEachSectionFailure(t *testing.T) {
	cases := []struct {
		name, directMetadata, managedMetadata string
		managedMatch                          bool
		directVersion, managedVersion         string
		warnings                              int
	}{
		{"fallback", upgradeAcceptanceFallback, "", false, "5.12.0", "5.11.4", 0},
		{"failure-before-unmatched-section", upgradeAcceptanceNoRelease, "", false, "5.11.3", "5.11.4", 1},
		{"failure-before-success", upgradeAcceptanceNoRelease, upgradeAcceptanceFallback, true, "5.11.3", "5.12.0", 1},
		{"success-before-failure", upgradeAcceptanceFallback, upgradeAcceptanceNoRelease, true, "5.12.0", "5.11.4", 1},
		{"both-fail", upgradeAcceptanceNoRelease, upgradeAcceptanceNoRelease, true, "5.11.3", "5.11.4", 2},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			path := upgradeAcceptanceMetadataPath("com.example", "client")
			responses := []upgradeAcceptanceResponse{{body: test.directMetadata}}
			managedArtifact := "unrelated"
			paths := []string{path}
			if test.managedMatch {
				responses = append(responses, upgradeAcceptanceResponse{body: test.managedMetadata})
				managedArtifact = "client"
				paths = append(paths, path)
			}
			fixture := newUpgradeAcceptanceFixture(t, map[string][]upgradeAcceptanceResponse{path: responses})
			body := "<dependencies>\n" + upgradeAcceptanceDependency("com.example", "client", "5.11.3") + "</dependencies>\n<dependencyManagement>\n<dependencies>\n" + upgradeAcceptanceDependency("com.example", managedArtifact, "5.11.4") + "</dependencies>\n</dependencyManagement>\n"
			directory := fixture.project("targeted", body, "")
			logs := fixture.run(directory, "upgrade", "dependency", "-g", "com.example", "-a", "client")
			assertUpgradeAcceptanceWarningCount(t, logs, test.warnings)
			for _, entry := range logs {
				if entry.Level == "warning" && (!strings.Contains(entry.Message, "com.example:client") || !strings.Contains(entry.Message, fixture.server.URL+"/selected/releases") || !strings.Contains(entry.Message, "suitable release")) {
					t.Errorf("targeted section lost lookup failure context: %+v", entry)
				}
			}
			model := upgradeAcceptanceModel(t, directory)
			assertUpgradeAcceptanceVersion(t, model, "client", test.directVersion, false)
			assertUpgradeAcceptanceVersion(t, model, managedArtifact, test.managedVersion, true)
			fixture.assertRequests(paths...)
		})
	}
}

func upgradeAcceptancePlugin(artifact string) string {
	return "<plugin>\n  <groupId>com.example</groupId>\n  <artifactId>" + artifact + "</artifactId>\n  <version>5.11.4</version>\n</plugin>\n"
}

const upgradeAcceptanceParent = "<parent>\n  <groupId>com.example</groupId>\n  <artifactId>parent</artifactId>\n  <version>5.11.4</version>\n</parent>\n"
const upgradeAcceptanceKotlin = "<properties>\n  <kotlin.version>5.11.4</kotlin.version>\n</properties>\n"

func TestUpgradeMetadataAcceptancePluginParentAndKotlin(t *testing.T) {
	cases := []struct{ command, group, artifact, body string }{
		{"plugins", "com.example", "build-plugin", "<build>\n<plugins>\n" + upgradeAcceptancePlugin("build-plugin") + "</plugins>\n</build>\n"},
		{"spring-boot", "com.example", "parent", upgradeAcceptanceParent},
		{"kotlin", "org.jetbrains.kotlin", "kotlin-maven-plugin", upgradeAcceptanceKotlin},
	}
	for _, test := range cases {
		for _, failed := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/failed=%v", test.command, failed), func(t *testing.T) {
				metadata := upgradeAcceptanceFallback
				if failed {
					metadata = upgradeAcceptanceNoRelease
				}
				path := upgradeAcceptanceMetadataPath(test.group, test.artifact)
				fixture := newUpgradeAcceptanceFixture(t, map[string][]upgradeAcceptanceResponse{path: {{body: metadata}}})
				directory := fixture.project("other-upgrades", test.body, "")
				logs := fixture.run(directory, "upgrade", test.command, "--stealth")
				expected := "5.12.0"
				if failed {
					expected = "5.11.4"
					assertUpgradeAcceptanceFailure(t, logs, test.group+":"+test.artifact, fixture.server.URL+"/selected/releases", "suitable release")
					assertUpgradeAcceptanceWarningCount(t, logs, 1)
				}
				model := upgradeAcceptanceModel(t, directory)
				var version string
				switch test.command {
				case "plugins":
					version = model.Build.Plugins.Plugin[0].Version
				case "spring-boot":
					version = model.Parent.Version
				case "kotlin":
					var err error
					version, err = model.Properties.FindKey("kotlin.version")
					if err != nil {
						t.Fatal(err)
					}
				}
				if version != expected {
					t.Errorf("%s version = %q, want %q", test.command, version, expected)
				}
				fixture.assertRequests(path)
			})
		}
	}
}

func TestUpgradeMetadataAcceptanceUpgradeAllContinuesJobsAndProjects(t *testing.T) {
	kotlin := upgradeAcceptanceMetadataPath("org.jetbrains.kotlin", "kotlin-maven-plugin")
	parent := upgradeAcceptanceMetadataPath("com.example", "parent")
	secondParty := upgradeAcceptanceMetadataPath("com.example", "second-party")
	thirdParty := upgradeAcceptanceMetadataPath("org.external", "third-party")
	failedPlugin := upgradeAcceptanceMetadataPath("com.example", "failed-plugin")
	laterPlugin := upgradeAcceptanceMetadataPath("com.example", "later-plugin")
	responses := map[string][]upgradeAcceptanceResponse{}
	for _, path := range []string{kotlin, parent, secondParty, failedPlugin} {
		responses[path] = []upgradeAcceptanceResponse{{body: upgradeAcceptanceNoRelease}, {body: upgradeAcceptanceFallback}}
	}
	for _, path := range []string{thirdParty, laterPlugin} {
		responses[path] = []upgradeAcceptanceResponse{{body: upgradeAcceptanceFallback}, {body: upgradeAcceptanceFallback}}
	}
	fixture := newUpgradeAcceptanceFixture(t, responses)
	body := upgradeAcceptanceParent + upgradeAcceptanceKotlin + "<dependencies>\n" +
		upgradeAcceptanceDependency("com.example", "second-party", "5.11.4") +
		upgradeAcceptanceDependency("org.external", "third-party", "5.11.4") + "</dependencies>\n" +
		"<build>\n<plugins>\n" + upgradeAcceptancePlugin("failed-plugin") + upgradeAcceptancePlugin("later-plugin") + "</plugins>\n</build>\n"
	first := fixture.project("a-first", body, "")
	second := fixture.project("b-second", body, "")
	logs := fixture.run(filepath.Dir(first), "upgrade", "all", "--recursive")
	assertUpgradeAcceptanceWarningCount(t, logs, 7) // Four errors plus three existing plugin upgrade notices.
	for _, artifact := range []string{"org.jetbrains.kotlin:kotlin-maven-plugin", "com.example:parent", "com.example:second-party", "com.example:failed-plugin"} {
		count := 0
		for _, entry := range logs {
			if entry.Level == "warning" && strings.Contains(entry.Message, "could not determine release for "+artifact) {
				count++
				if !strings.Contains(entry.Message, fixture.server.URL+"/selected/releases") || !strings.Contains(entry.Message, "suitable release") {
					t.Errorf("upgrade all warning lost context: %+v", entry)
				}
			}
		}
		if count != 1 {
			t.Errorf("upgrade all error count for %s = %d, want 1; logs: %+v", artifact, count, logs)
		}
	}
	for index, directory := range []string{first, second} {
		model := upgradeAcceptanceModel(t, directory)
		expected := "5.11.4"
		if index == 1 {
			expected = "5.12.0"
		}
		if model.Parent.Version != expected {
			t.Errorf("project %d parent = %s, want %s", index, model.Parent.Version, expected)
		}
		kotlinVersion, err := model.Properties.FindKey("kotlin.version")
		if err != nil || kotlinVersion != expected {
			t.Errorf("project %d kotlin = %s (%v), want %s", index, kotlinVersion, err, expected)
		}
		assertUpgradeAcceptanceVersion(t, model, "second-party", expected, false)
		assertUpgradeAcceptanceVersion(t, model, "third-party", "5.12.0", false)
		for _, plugin := range model.Build.Plugins.Plugin {
			want := "5.12.0"
			if plugin.ArtifactId == "failed-plugin" {
				want = expected
			}
			if plugin.Version != want {
				t.Errorf("project %d plugin %s = %s, want %s", index, plugin.ArtifactId, plugin.Version, want)
			}
		}
	}
	fixture.assertRequests(kotlin, parent, secondParty, thirdParty, failedPlugin, laterPlugin, kotlin, parent, secondParty, thirdParty, failedPlugin, laterPlugin)
}

func TestUpgradeMetadataAcceptanceBuiltExecutableJourney(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "ply")
	buildContext, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	build := exec.CommandContext(buildContext, "go", "build", "-o", binary, "./cmd/ply")
	build.Dir = testRepositoryRoot(t)
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build actual CLI: %v\n%s", err, output)
	}
	failed, later := upgradeAcceptanceMetadataPath("com.example", "failed"), upgradeAcceptanceMetadataPath("com.example", "later")
	fixture := newUpgradeAcceptanceFixture(t, map[string][]upgradeAcceptanceResponse{failed: {{body: upgradeAcceptanceNoRelease}}, later: {{body: upgradeAcceptanceFallback}}})
	directory := fixture.project("binary-journey", "<dependencies>\n"+upgradeAcceptanceDependency("com.example", "failed", "5.11.4")+upgradeAcceptanceDependency("com.example", "later", "5.11.4")+"</dependencies>\n", "")
	logs := fixture.runExecutable(binary, directory, "upgrade", "2party", "--stealth")
	assertUpgradeAcceptanceFailure(t, logs, "com.example:failed", fixture.server.URL+"/selected/releases", "suitable release")
	assertUpgradeAcceptanceWarningCount(t, logs, 1)
	model := upgradeAcceptanceModel(t, directory)
	assertUpgradeAcceptanceVersion(t, model, "failed", "5.11.4", false)
	assertUpgradeAcceptanceVersion(t, model, "later", "5.12.0", false)
	fixture.assertRequests(failed, later)
}

func TestUpgradeMetadataAcceptanceRedactsRepositoryURLCredentials(t *testing.T) {
	failed, later := upgradeAcceptanceMetadataPath("com.example", "failed"), upgradeAcceptanceMetadataPath("com.example", "later")
	fixture := newUpgradeAcceptanceFixture(t, map[string][]upgradeAcceptanceResponse{failed: {{status: http.StatusUnauthorized, body: "sensitive-response-marker"}}, later: {{body: upgradeAcceptanceFallback}}})
	repositoryURL, err := url.Parse(fixture.server.URL + "/selected/releases")
	if err != nil {
		t.Fatal(err)
	}
	repositoryURL.User = url.UserPassword("acceptance-url-user", "acceptance-url-password")
	profile := filepath.Join(fixture.home, ".ply", "profiles", "acceptance")
	writeUpgradeAcceptanceFile(t, filepath.Join(profile, "local-config.yaml"), fmt.Sprintf("nexus:\n  url: %q\n  username: %q\n  password: %q\n", repositoryURL.String(), upgradeAcceptanceUser, upgradeAcceptancePassword))
	directory := fixture.project("redaction", "<dependencies>\n"+upgradeAcceptanceDependency("com.example", "failed", "5.11.4")+upgradeAcceptanceDependency("com.example", "later", "5.11.4")+"</dependencies>\n", "")
	logs := fixture.run(directory, "upgrade", "2party", "--stealth", "--debug")
	assertUpgradeAcceptanceFailure(t, logs, "com.example:failed", fixture.server.URL+"/selected/releases", "HTTP 401")
	assertUpgradeAcceptanceWarningCount(t, logs, 1)
	model := upgradeAcceptanceModel(t, directory)
	assertUpgradeAcceptanceVersion(t, model, "failed", "5.11.4", false)
	assertUpgradeAcceptanceVersion(t, model, "later", "5.12.0", false)
	fixture.assertRequests(failed, later)
}

func readUpgradeAcceptancePOM(t *testing.T, directory string) []byte {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(directory, "pom.xml"))
	if err != nil {
		t.Fatal(err)
	}
	return contents
}

func assertUpgradeAcceptanceWarningCount(t *testing.T, logs []upgradeAcceptanceLog, expected int) {
	t.Helper()
	count := 0
	for _, entry := range logs {
		if entry.Level == "warning" {
			count++
		}
	}
	if count != expected {
		t.Errorf("warning count = %d, want %d; logs: %+v", count, expected, logs)
	}
}
