package config

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/sirupsen/logrus"
)

var _ ProjectConfig = (*ProjectConfiguration)(nil)

type recordedProjectConfigWrite struct {
	path           string
	data           []byte
	mode           fs.FileMode
	logBeforeWrite string
}

type recordingProjectConfigFilesystem struct {
	writes    []recordedProjectConfigWrite
	writeErr  error
	logOutput *bytes.Buffer
}

func (*recordingProjectConfigFilesystem) WorkingDirectory() (string, error) {
	return "", errors.New("unexpected project-config working directory")
}

func (*recordingProjectConfigFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected project-config read")
}

func (*recordingProjectConfigFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, errors.New("unexpected project-config read directory")
}

func (*recordingProjectConfigFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, errors.New("unexpected project-config read directory entries")
}

func (*recordingProjectConfigFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected project-config stat")
}

func (*recordingProjectConfigFilesystem) Mkdir(string, fs.FileMode) error {
	return errors.New("unexpected project-config mkdir")
}

func (*recordingProjectConfigFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected project-config mkdir")
}

func (recording *recordingProjectConfigFilesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	logBeforeWrite := ""
	if recording.logOutput != nil {
		logBeforeWrite = recording.logOutput.String()
	}
	recording.writes = append(recording.writes, recordedProjectConfigWrite{
		path:           path,
		data:           append([]byte{}, data...),
		mode:           mode,
		logBeforeWrite: logBeforeWrite,
	})
	return recording.writeErr
}

func (*recordingProjectConfigFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected project-config open file")
}

func (*recordingProjectConfigFilesystem) Remove(string) error {
	return errors.New("unexpected project-config remove")
}

func (*recordingProjectConfigFilesystem) RemoveAll(string) error {
	return errors.New("unexpected project-config recursive remove")
}

func (*recordingProjectConfigFilesystem) Rename(string, string) error {
	return errors.New("unexpected project-config rename")
}

func (*recordingProjectConfigFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected project-config glob")
}

func (*recordingProjectConfigFilesystem) Walk(string, filepath.WalkFunc) error {
	return errors.New("unexpected project-config walk")
}

func (*recordingProjectConfigFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected project-config create")
}

func (*recordingProjectConfigFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected project-config copy")
}

func (recording *recordingProjectConfigFilesystem) dependencies() projectConfigWriteDependencies {
	return projectConfigWriteDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingProjectConfigFilesystem) assertedWrites() ([]recordedProjectConfigWrite, error) {
	if len(recording.writes) == 0 {
		return nil, errors.New("recorded project-config write population is empty")
	}
	return recording.writes, nil
}

type projectConfigMessageOnlyFormatter struct{}

func (projectConfigMessageOnlyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

func TestProjectConfigurationWriteToSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemProjectConfigWriteDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("ProjectConfiguration.WriteTo selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("ProjectConfiguration.WriteTo filesystem dependency is %T, want %T",
			dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestProjectConfigurationWriteToPreservesLogMarshalTargetBytesModeErrorAndReceiver(t *testing.T) {
	logOutput := captureProjectConfigLog(t)
	writeError := errors.New("complete project-config write dependency error")
	tests := []struct {
		name       string
		config     *ProjectConfiguration
		wantConfig *ProjectConfiguration
		target     string
		wantJSON   string
		writeErr   error
	}{
		{
			name:     "nil receiver",
			config:   nil,
			target:   `/arbitrary project-config target//../backslash\nil-ø.json`,
			wantJSON: "null",
		},
		{
			name:       "zero-value receiver",
			config:     &ProjectConfiguration{},
			wantConfig: &ProjectConfiguration{},
			target:     `/arbitrary project-config target//../backslash\empty.json`,
			wantJSON: `{
    "groupId": "",
    "artifactId": "",
    "language": "",
    "package": "",
    "applicationName": "",
    "name": "",
    "description": "",
    "team": {
        "name": "",
        "email": ""
    },
    "dependencies": null,
    "templates": null,
    "settings": {
        "disableDependencySort": false,
        "disableSpringBootUpgrade": false,
        "disableKotlinUpgrade": false,
        "pomFileIndentation": "",
        "disableUpgradesFor": null,
        "maxSpringBootVersion": "",
        "maxVersionForDependencies": null
    },
    "render": null
}`,
		},
		{
			name:       "empty non-ASCII fields collections and map keys",
			config:     representativeProjectConfiguration(),
			wantConfig: representativeProjectConfiguration(),
			target:     `/arbitrary project-config target//../backslash\representative-ø.json`,
			wantJSON: `{
    "groupId": "",
    "artifactId": "artifakt-ø",
    "language": "",
    "package": "no.økonomi",
    "applicationName": "",
    "profile": "profil-æ",
    "name": "",
    "description": "linje α\nandre",
    "team": {
        "name": "Team Ø",
        "email": "",
        "slackChannel": "kanal-β"
    },
    "dependencies": [],
    "templates": null,
    "settings": {
        "disableDependencySort": false,
        "disableSpringBootUpgrade": true,
        "disableKotlinUpgrade": false,
        "pomFileIndentation": "",
        "disableUpgradesFor": [],
        "maxSpringBootVersion": "",
        "maxVersionForDependencies": null
    },
    "render": {
        "a-key": "",
        "ø-key": "verdi"
    }
}`,
			writeErr: writeError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logOutput.Reset()
			recording := &recordingProjectConfigFilesystem{writeErr: test.writeErr, logOutput: logOutput}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("project-config dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := test.config.writeTo(dependencies, test.target)

			if err != test.writeErr {
				t.Fatalf("project-config write error was %v, want exact dependency error %v", err, test.writeErr)
			}
			writes, populationErr := recording.assertedWrites()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			wantLog := "writes project config file to " + test.target + "\n"
			wantWrites := []recordedProjectConfigWrite{{
				path:           test.target,
				data:           []byte(test.wantJSON),
				mode:           0644,
				logBeforeWrite: wantLog,
			}}
			if !reflect.DeepEqual(writes, wantWrites) {
				t.Fatalf("recorded project-config writes differ:\n got: %#v\nwant: %#v", writes, wantWrites)
			}
			if logOutput.String() != wantLog {
				t.Fatalf("project-config log was %q, want exact first log %q", logOutput.String(), wantLog)
			}
			if !reflect.DeepEqual(test.config, test.wantConfig) {
				t.Fatalf("project-config receiver changed:\n got: %#v\nwant: %#v", test.config, test.wantConfig)
			}
		})
	}
}

func TestProjectConfigurationWriteToDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	logOutput := captureProjectConfigLog(t)
	config := representativeProjectConfiguration()
	wantConfig := representativeProjectConfiguration()
	target := "/developer/home/project/must-not-be-accessed/ply.json"

	err := config.writeTo(projectConfigWriteDependencies{}, target)

	if err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe project-config dependency default returned %v, want exact %v", err, filesystem.ErrNoFilesystem)
	}
	wantLog := "writes project config file to " + target + "\n"
	if logOutput.String() != wantLog {
		t.Fatalf("safe project-config log was %q, want exact first log %q", logOutput.String(), wantLog)
	}
	if !reflect.DeepEqual(config, wantConfig) {
		t.Fatalf("safe project-config write changed receiver:\n got: %#v\nwant: %#v", config, wantConfig)
	}
}

func TestRecordedProjectConfigurationWriteToRejectsEmptyWritePopulation(t *testing.T) {
	recording := &recordingProjectConfigFilesystem{}

	if _, err := recording.assertedWrites(); err == nil {
		t.Fatal("empty recorded project-config write population passed")
	}
}

func representativeProjectConfiguration() *ProjectConfiguration {
	config := &ProjectConfiguration{
		MavenProjectConfiguration: MavenProjectConfiguration{
			Artifact: Artifact{
				ArtifactId: "artifakt-ø",
			},
			Package: "no.økonomi",
		},
		Profile:      "profil-æ",
		Description:  "linje α\nandre",
		Dependencies: []string{},
		Templates:    nil,
		Settings: ProjectSettings{
			DisableSpringBootUpgrade:  true,
			DisableUpgradesFor:        []Artifact{},
			MaxVersionForDependencies: nil,
		},
		Render: map[string]string{
			"ø-key": "verdi",
			"a-key": "",
		},
	}
	config.Team.Name = "Team Ø"
	config.Team.SlackChannel = "kanal-β"
	return config
}

func captureProjectConfigLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	output := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(output)
	logger.SetFormatter(projectConfigMessageOnlyFormatter{})
	previousLogger := log
	log = logger
	t.Cleanup(func() { log = previousLogger })
	return output
}
