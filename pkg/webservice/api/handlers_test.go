package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/spring"
)

type staticGenerateCloudConfig struct {
	config.CloudConfig
	templates []config.CloudTemplate
	calls     int
}

func (cloud *staticGenerateCloudConfig) Templates() ([]config.CloudTemplate, error) {
	cloud.calls++
	return append([]config.CloudTemplate(nil), cloud.templates...), nil
}

type recordedCallbacks struct {
	values []bool
}

func (recording *recordedCallbacks) receive(callbacks <-chan bool) {
	recording.values = append(recording.values, <-callbacks)
}

func (recording *recordedCallbacks) assertedValues() ([]bool, error) {
	if len(recording.values) == 0 {
		return nil, errors.New("recorded webservice API callback population is empty")
	}
	return append([]bool(nil), recording.values...), nil
}

type packageGlobals struct {
	generateOptions GenerateOptions
	callbackChannel chan bool
	currentProject  config.Project
}

func TestGetGenerateRendersExactStatusAndBodyBytesAndRestoresGlobals(t *testing.T) {
	restored := preservePackageGlobalsAroundTest(t)

	t.Run("representative existing options", func(t *testing.T) {
		projectConfig := &config.ProjectConfiguration{
			MavenProjectConfiguration: config.MavenProjectConfiguration{
				Artifact: config.Artifact{
					GroupId:    "com.example<&",
					ArtifactId: `artifact "quoted"`,
				},
				Language: "java",
				Package:  "com.example.raw-package",
			},
			Name:        "Name & <Raw>",
			Description: `Description "quoted" & <raw>`,
		}
		cloud := &staticGenerateCloudConfig{templates: []config.CloudTemplate{
			{Name: "base<&>"},
			{Name: `second "template"`},
		}}
		if len(cloud.templates) == 0 {
			t.Fatal("Generate GET cloud-template population is empty")
		}
		root := representativeRootResponse(t)
		if len(root.Dependencies.Values) == 0 {
			t.Fatal("Generate GET dependency-group population is empty")
		}
		options := GenerateOptions{
			ProjectConfig: projectConfig,
			CloudConfig:   cloud,
			IoResponse:    root,
		}
		callbackChannel := make(chan bool)
		currentProject := config.Project{Path: "complete Generate GET current project"}
		installPackageGlobals(t, packageGlobals{
			generateOptions: options,
			callbackChannel: callbackChannel,
			currentProject:  currentProject,
		})

		request := httptest.NewRequest(http.MethodGet, "/ui/generate?ignored=complete", nil)
		response := httptest.NewRecorder()

		GetGenerate(response, request)

		assertExactHashedResponse(t, response, http.StatusOK, 2957,
			"9a9eb38ab026e751b06dee004e84d8db9a427454c5e33f4b864ebc83635399c1")
		if cloud.calls != 1 {
			t.Fatalf("Generate GET requested cloud templates %d times, want exactly 1", cloud.calls)
		}
		rawRenderedValues := []string{
			`value=com.example<&>`,
			`value=artifact "quoted">`,
			`value=com.example.raw-package>`,
			`value=Name & <Raw>>`,
			`value=Description "quoted" & <raw>>`,
			`<option>base<&></option>`,
			`<option>second "template"</option>`,
			`<option disabled>Messaging & <Queue></option>`,
			`<option value="queue<&quot;">Queue & <Primary></option>`,
		}
		if len(rawRenderedValues) == 0 {
			t.Fatal("Generate GET raw rendered-value population is empty")
		}
		body := response.Body.String()
		for _, value := range rawRenderedValues {
			if !strings.Contains(body, value) {
				t.Fatalf("Generate GET exact body omitted raw rendered bytes %q", value)
			}
		}
		assertInstalledPackageGlobals(t, packageGlobals{
			generateOptions: options,
			callbackChannel: callbackChannel,
			currentProject:  currentProject,
		})
	})

	assertInstalledPackageGlobals(t, restored)
}

func TestGetUpgradeRendersExactStatusAndBodyBytesAndRestoresGlobals(t *testing.T) {
	restored := preservePackageGlobalsAroundTest(t)

	t.Run("representative existing options", func(t *testing.T) {
		projectConfig := &config.ProjectConfiguration{
			MavenProjectConfiguration: config.MavenProjectConfiguration{
				Artifact: config.Artifact{GroupId: "unused<&", ArtifactId: `unused "quoted"`},
			},
			Name: "unused Upgrade GET option",
		}
		cloud := &staticGenerateCloudConfig{templates: []config.CloudTemplate{{Name: "unused template"}}}
		options := GenerateOptions{ProjectConfig: projectConfig, CloudConfig: cloud}
		callbackChannel := make(chan bool)
		currentProject := config.Project{Path: "complete Upgrade GET current project"}
		installPackageGlobals(t, packageGlobals{
			generateOptions: options,
			callbackChannel: callbackChannel,
			currentProject:  currentProject,
		})

		request := httptest.NewRequest(http.MethodGet, "/ui/upgrade?ignored=complete", nil)
		response := httptest.NewRecorder()

		GetUpgrade(response, request)

		assertExactHashedResponse(t, response, http.StatusOK, 788,
			"ae9194d9a6734d1fe568692e11327fc1b31b3f0d797fd558fd6ceba3ffa58ecc")
		if cloud.calls != 0 {
			t.Fatalf("Upgrade GET unexpectedly requested Generate cloud templates %d times", cloud.calls)
		}
		assertInstalledPackageGlobals(t, packageGlobals{
			generateOptions: options,
			callbackChannel: callbackChannel,
			currentProject:  currentProject,
		})
	})

	assertInstalledPackageGlobals(t, restored)
}

func TestPostGeneratePreservesFormMutationCallbackAndExactResponseAndRestoresGlobals(t *testing.T) {
	restored := preservePackageGlobalsAroundTest(t)

	t.Run("complete form", func(t *testing.T) {
		projectConfig := &config.ProjectConfiguration{
			MavenProjectConfiguration: config.MavenProjectConfiguration{
				Artifact:        config.Artifact{GroupId: "old group", ArtifactId: "old artifact"},
				Language:        "old language",
				Package:         "old package",
				ApplicationName: "preserved application",
			},
			Profile:      "preserved profile",
			Name:         "old name",
			Description:  "old description",
			Dependencies: []string{"existing dependency"},
			Templates:    []string{"existing template"},
			Settings: config.ProjectSettings{
				DisableDependencySort: true,
				MaxSpringBootVersion:  "preserved version",
			},
			Render: map[string]string{"preserved": "render value"},
		}
		projectConfig.Team.Name = "preserved team"
		initial := *projectConfig
		cloud := &staticGenerateCloudConfig{templates: []config.CloudTemplate{{Name: "unused POST template"}}}
		options := GenerateOptions{
			ProjectConfig: projectConfig,
			CloudConfig:   cloud,
			IoResponse:    representativeRootResponse(t),
		}
		callbackChannel := make(chan bool)
		if cap(callbackChannel) != 0 {
			t.Fatal("Generate POST test callback channel is not unbuffered")
		}
		currentProject := config.Project{Path: "complete Generate POST current project"}
		installPackageGlobals(t, packageGlobals{
			generateOptions: options,
			callbackChannel: callbackChannel,
			currentProject:  currentProject,
		})

		form := url.Values{
			"groupId":      {"first decoded group", "ignored second group"},
			"artifactId":   {"first/artifact", "ignored second artifact"},
			"package":      {"first package + decoded", "ignored second package"},
			"name":         {"first name & decoded", "ignored second name"},
			"description":  {"first description = decoded", "ignored second description"},
			"language":     {"kotlin", "ignored second language"},
			"templates":    {"first template", "second template", "first template"},
			"dependencies": {"first:dependency", "second&dependency", "first:dependency"},
			"unrelated":    {"must not mutate project"},
		}
		if len(form) == 0 {
			t.Fatal("Generate POST form population is empty")
		}
		request := httptest.NewRequest(http.MethodPost,
			"/api/generate?groupId=query-group&templates=query-template",
			strings.NewReader(form.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()

		PostGenerate(response, request)

		callbacks := &recordedCallbacks{}
		callbacks.receive(callbackChannel)
		assertOneExactCallback(t, callbacks, callbackChannel)
		assertExactResponse(t, response, http.StatusOK, "OK")
		if !reflect.DeepEqual(request.PostForm, form) {
			t.Fatalf("Generate POST parsed body form differs:\n got: %#v\nwant: %#v", request.PostForm, form)
		}

		want := initial
		want.GroupId = form["groupId"][0]
		want.ArtifactId = form["artifactId"][0]
		want.Package = form["package"][0]
		want.Name = form["name"][0]
		want.Description = form["description"][0]
		want.Language = form["language"][0]
		want.Templates = append(append([]string(nil), initial.Templates...), form["templates"]...)
		want.Dependencies = append(append([]string(nil), initial.Dependencies...), form["dependencies"]...)
		if !reflect.DeepEqual(*projectConfig, want) {
			t.Fatalf("Generate POST project mutation differs:\n got: %#v\nwant: %#v", *projectConfig, want)
		}
		scalarFields := []struct {
			name string
			got  string
			want string
		}{
			{name: "groupId", got: projectConfig.GroupId, want: form["groupId"][0]},
			{name: "artifactId", got: projectConfig.ArtifactId, want: form["artifactId"][0]},
			{name: "package", got: projectConfig.Package, want: form["package"][0]},
			{name: "name", got: projectConfig.Name, want: form["name"][0]},
			{name: "description", got: projectConfig.Description, want: form["description"][0]},
			{name: "language", got: projectConfig.Language, want: form["language"][0]},
		}
		if len(scalarFields) == 0 {
			t.Fatal("Generate POST scalar-field characterization population is empty")
		}
		for _, field := range scalarFields {
			if field.got != field.want {
				t.Fatalf("Generate POST %s was %q, want first body value %q", field.name, field.got, field.want)
			}
		}
		assertInstalledPackageGlobals(t, packageGlobals{
			generateOptions: options,
			callbackChannel: callbackChannel,
			currentProject:  currentProject,
		})
	})

	assertInstalledPackageGlobals(t, restored)
}

func TestPostUpgradeDeliversOneCallbackAndExactResponseAndRestoresGlobals(t *testing.T) {
	restored := preservePackageGlobalsAroundTest(t)

	t.Run("arbitrary request body", func(t *testing.T) {
		projectConfig := &config.ProjectConfiguration{Name: "unchanged Upgrade POST project config"}
		options := GenerateOptions{
			ProjectConfig: projectConfig,
			CloudConfig:   &staticGenerateCloudConfig{templates: []config.CloudTemplate{{Name: "unused"}}},
		}
		callbackChannel := make(chan bool)
		if cap(callbackChannel) != 0 {
			t.Fatal("Upgrade POST test callback channel is not unbuffered")
		}
		currentProject := config.Project{Path: "complete Upgrade POST current project"}
		installPackageGlobals(t, packageGlobals{
			generateOptions: options,
			callbackChannel: callbackChannel,
			currentProject:  currentProject,
		})
		request := httptest.NewRequest(http.MethodPost, "/api/upgrade?ignored=complete",
			strings.NewReader("malformed=%zz"))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()

		PostUpgrade(response, request)

		callbacks := &recordedCallbacks{}
		callbacks.receive(callbackChannel)
		assertOneExactCallback(t, callbacks, callbackChannel)
		assertExactResponse(t, response, http.StatusOK, "OK")
		if request.Form != nil || request.PostForm != nil {
			t.Fatalf("Upgrade POST unexpectedly parsed request forms: Form=%#v PostForm=%#v",
				request.Form, request.PostForm)
		}
		assertInstalledPackageGlobals(t, packageGlobals{
			generateOptions: options,
			callbackChannel: callbackChannel,
			currentProject:  currentProject,
		})
	})

	assertInstalledPackageGlobals(t, restored)
}

func TestRecordedWebserviceAPICallbacksRejectEmptyPopulation(t *testing.T) {
	recording := &recordedCallbacks{}

	if _, err := recording.assertedValues(); err == nil {
		t.Fatal("empty recorded webservice API callback population passed")
	}
}

func representativeRootResponse(t *testing.T) spring.IoRootResponse {
	t.Helper()
	contents := `{
		"dependencies": {
			"type": "hierarchical-multi-select",
			"values": [
				{
					"name": "Messaging & <Queue>",
					"values": [
						{"id": "queue<&quot;", "name": "Queue & <Primary>"},
						{"id": "second/id", "name": "Second > Queue"}
					]
				},
				{
					"name": "Storage",
					"values": [{"id": "store+id", "name": "Store"}]
				}
			]
		}
	}`
	var response spring.IoRootResponse
	if err := json.Unmarshal([]byte(contents), &response); err != nil {
		t.Fatalf("prepare representative Spring root response: %v", err)
	}
	return response
}

func preservePackageGlobalsAroundTest(t *testing.T) packageGlobals {
	t.Helper()
	original := snapshotPackageGlobals()
	t.Cleanup(func() { restorePackageGlobals(original) })

	sentinel := packageGlobals{
		generateOptions: GenerateOptions{
			ProjectConfig: &config.ProjectConfiguration{Name: "restored sentinel options"},
			CloudConfig: &staticGenerateCloudConfig{
				templates: []config.CloudTemplate{{Name: "restored sentinel template"}},
			},
		},
		callbackChannel: make(chan bool),
		currentProject:  config.Project{Path: "restored sentinel current project"},
	}
	restorePackageGlobals(sentinel)
	return sentinel
}

func installPackageGlobals(t *testing.T, installed packageGlobals) {
	t.Helper()
	previous := snapshotPackageGlobals()
	restorePackageGlobals(installed)
	t.Cleanup(func() { restorePackageGlobals(previous) })
}

func snapshotPackageGlobals() packageGlobals {
	return packageGlobals{
		generateOptions: GOptions,
		callbackChannel: CallbackChannel,
		currentProject:  CurrentProject,
	}
}

func restorePackageGlobals(globals packageGlobals) {
	GOptions = globals.generateOptions
	CallbackChannel = globals.callbackChannel
	CurrentProject = globals.currentProject
}

func assertInstalledPackageGlobals(t *testing.T, want packageGlobals) {
	t.Helper()
	if GOptions.ProjectConfig != want.generateOptions.ProjectConfig ||
		GOptions.CloudConfig != want.generateOptions.CloudConfig ||
		!reflect.DeepEqual(GOptions.IoResponse, want.generateOptions.IoResponse) {
		t.Fatalf("Generate options global differs:\n got: %#v\nwant: %#v", GOptions, want.generateOptions)
	}
	if CallbackChannel != want.callbackChannel {
		t.Fatalf("callback channel global identity is %p, want %p", CallbackChannel, want.callbackChannel)
	}
	if !reflect.DeepEqual(CurrentProject, want.currentProject) {
		t.Fatalf("current project global differs:\n got: %#v\nwant: %#v", CurrentProject, want.currentProject)
	}
}

func assertOneExactCallback(t *testing.T, recording *recordedCallbacks, channel <-chan bool) {
	t.Helper()
	callbacks, err := recording.assertedValues()
	if err != nil {
		t.Fatal(err)
	}
	if len(callbacks) != 1 {
		t.Fatalf("webservice API delivered %d callbacks, want exactly 1", len(callbacks))
	}
	for index, callback := range callbacks {
		if !callback {
			t.Fatalf("webservice API callback %d was false, want true", index)
		}
	}
	select {
	case extra := <-channel:
		t.Fatalf("webservice API delivered an additional callback %t", extra)
	default:
	}
}

func assertExactResponse(t *testing.T, response *httptest.ResponseRecorder, wantStatus int, wantBody string) {
	t.Helper()
	if response.Code != wantStatus {
		t.Fatalf("HTTP status was %d, want exactly %d", response.Code, wantStatus)
	}
	if body := response.Body.String(); body != wantBody {
		t.Fatalf("HTTP body was %q, want exact bytes %q", body, wantBody)
	}
}

func assertExactHashedResponse(
	t *testing.T,
	response *httptest.ResponseRecorder,
	wantStatus int,
	wantLength int,
	wantSHA256 string,
) {
	t.Helper()
	if response.Code != wantStatus {
		t.Fatalf("HTTP status was %d, want exactly %d", response.Code, wantStatus)
	}
	body := response.Body.Bytes()
	digest := sha256.Sum256(body)
	gotSHA256 := hex.EncodeToString(digest[:])
	if len(body) != wantLength || gotSHA256 != wantSHA256 {
		t.Fatalf("HTTP body has length %d and SHA-256 %s, want exact length %d and SHA-256 %s",
			len(body), gotSHA256, wantLength, wantSHA256)
	}
}
