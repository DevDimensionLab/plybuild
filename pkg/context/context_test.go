package context

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/devdimensionlab/plybuild/internal/testutil"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/logger"
	"github.com/devdimensionlab/plybuild/pkg/maven"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)

type contextLogRecord struct {
	level   logrus.Level
	message string
	data    logrus.Fields
}

type contextLogHook struct {
	records []contextLogRecord
}

func (*contextLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *contextLogHook) Fire(entry *logrus.Entry) error {
	data := make(logrus.Fields, len(entry.Data))
	for key, value := range entry.Data {
		data[key] = value
	}
	hook.records = append(hook.records, contextLogRecord{
		level: entry.Level, message: entry.Message, data: data,
	})
	return nil
}

func (hook *contextLogHook) assertedRecords() ([]contextLogRecord, error) {
	if len(hook.records) == 0 {
		return nil, errors.New("recorded Context log population is empty")
	}
	return append([]contextLogRecord(nil), hook.records...), nil
}

type contextMessageOnlyFormatter struct{}

func (contextMessageOnlyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

type contextCloudDefaults struct {
	config.CloudConfig
	defaults config.CloudProjectDefaults
	err      error
	calls    int
}

func (cloud *contextCloudDefaults) ProjectDefaults() (config.CloudProjectDefaults, error) {
	cloud.calls++
	return cloud.defaults, cloud.err
}

type contextCloudConfigOpener struct {
	cloud        config.CloudConfig
	profilePaths []string
}

func (opener *contextCloudConfigOpener) Open(profilePath string) config.CloudConfig {
	opener.profilePaths = append(opener.profilePaths, profilePath)
	return opener.cloud
}

func (opener *contextCloudConfigOpener) assertedProfilePaths() ([]string, error) {
	if len(opener.profilePaths) == 0 {
		return nil, errors.New("recorded cloud-config opener profile-path population is empty")
	}
	return append([]string(nil), opener.profilePaths...), nil
}

type contextEnvironmentValue struct {
	value string
	set   bool
}

func TestContextPackageLoggerInitializationAndSetLoggerPreserveExactSelection(t *testing.T) {
	initialized, ok := log.(*logrus.Entry)
	if !ok {
		t.Fatalf("Context package logger initialized as %T, want *logrus.Entry", log)
	}
	loggerContext, ok := logger.Context().(*logrus.Entry)
	if !ok {
		t.Fatalf("logger.Context returned %T, want *logrus.Entry", logger.Context())
	}
	if initialized.Logger != loggerContext.Logger {
		t.Fatalf("Context package logger uses %p, want logger package private identity %p",
			initialized.Logger, loggerContext.Logger)
	}
	if initialized.Data == nil || len(initialized.Data) != 0 {
		t.Fatalf("Context package logger fields were %#v, want allocated empty fields", initialized.Data)
	}

	original := log
	tests := []struct {
		name     string
		selected logrus.FieldLogger
	}{
		{name: "logger pointer", selected: logrus.New()},
		{name: "entry pointer", selected: logrus.New().WithField("selection", "exact entry")},
		{name: "nil interface", selected: nil},
	}
	if len(tests) == 0 {
		t.Fatal("Context SetLogger characterization population is empty")
	}
	for _, test := range tests {
		runContextStateSubtest(t, test.name, func(t *testing.T) {
			SetLogger(test.selected)
			if !sameContextFieldLogger(log, test.selected) {
				t.Fatalf("SetLogger stored %T %#x, want exact %T %#x",
					log, contextInterfacePointer(log), test.selected, contextInterfacePointer(test.selected))
			}
		})
	}
	if !sameContextFieldLogger(log, original) {
		t.Fatalf("Context logger restored as %T %#x, want %T %#x",
			log, contextInterfacePointer(log), original, contextInterfacePointer(original))
	}
}

func TestContextPublicStatePreservesZeroValuesAndDirectFieldOwnership(t *testing.T) {
	var zero Context
	if zero.Recursive || zero.DryRun || zero.DisableGit || zero.ForceCloudSync ||
		zero.OpenInBrowser || zero.StealthMode || zero.TargetDirectory != "" ||
		zero.ProfilesPath != "" || zero.Projects != nil || zero.Err != nil ||
		zero.LocalConfig.Implementation().Path != "" || zero.CloudConfig != nil {
		t.Fatalf("zero Context state changed: %#v", zero)
	}

	sentinelError := errors.New("complete Context state error")
	projects := []config.Project{{Path: "first direct project"}, {Path: "second direct project"}}
	if len(projects) == 0 {
		t.Fatal("Context public-state project population is empty")
	}
	local := config.OpenLocalConfig("/complete/direct/profile")
	cloud := &contextCloudDefaults{}
	ctx := Context{
		Recursive:       true,
		DryRun:          true,
		TargetDirectory: "/complete/direct/target",
		DisableGit:      true,
		ForceCloudSync:  true,
		OpenInBrowser:   true,
		StealthMode:     true,
		Projects:        projects,
		Err:             sentinelError,
		ProfilesPath:    "/complete/direct/profiles",
		LocalConfig:     local,
		CloudConfig:     cloud,
	}
	projects[0].Path = "mutated through shared slice"
	if ctx.Projects[0].Path != projects[0].Path || ctx.Err != sentinelError ||
		ctx.LocalConfig.Implementation().Path != local.Implementation().Path || ctx.CloudConfig != cloud ||
		!ctx.Recursive || !ctx.DryRun || !ctx.DisableGit || !ctx.ForceCloudSync ||
		!ctx.OpenInBrowser || !ctx.StealthMode {
		t.Fatalf("direct Context field state or slice ownership changed: %#v", ctx)
	}
}

func TestFindAndPopulateMavenProjectsPreservesNonRecursiveAppendSuccessAndDirectError(t *testing.T) {
	tests := []struct {
		name      string
		prepare   func(*testing.T, string) string
		wantError string
	}{
		{
			name: "success appends after existing project",
			prepare: func(t *testing.T, root string) string {
				return writeContextProject(t, root, "single", true)
			},
		},
		{
			name: "configuration population error returns before append",
			prepare: func(t *testing.T, root string) string {
				return writeBrokenContextProject(t, root, "broken")
			},
			wantError: "directory detected, but language was not set",
		},
	}
	if len(tests) == 0 {
		t.Fatal("non-recursive Context discovery characterization population is empty")
	}

	for _, test := range tests {
		runContextStateSubtest(t, test.name, func(t *testing.T) {
			profile := isolateContextHome(t)
			root := t.TempDir()
			target := test.prepare(t, root)
			existing := config.Project{Path: "existing project remains first"}
			ctx := Context{TargetDirectory: target, Projects: []config.Project{existing}}

			err := ctx.FindAndPopulateMavenProjects()

			if test.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantError) {
					t.Fatalf("non-recursive discovery error was %v, want text containing %q", err, test.wantError)
				}
				if !reflect.DeepEqual(ctx.Projects, []config.Project{existing}) {
					t.Fatalf("non-recursive error changed projects: %#v", ctx.Projects)
				}
				return
			}
			if err != nil {
				t.Fatalf("non-recursive discovery returned an error: %v", err)
			}
			if len(ctx.Projects) != 2 || ctx.Projects[0].Path != existing.Path {
				t.Fatalf("non-recursive append order was %#v", ctx.Projects)
			}
			project := ctx.Projects[1]
			if project.Path != target || project.Config.Name != "single" || !project.IsMavenProject() {
				t.Fatalf("non-recursive project was %#v, want exact loaded target", project)
			}
			cloud, ok := project.CloudConfig.(config.GitCloudConfig)
			if !ok || cloud.Implementation().Dir() != filepath.Join(profile, "cloud-config") {
				t.Fatalf("non-recursive cloud config was %T %#v, want profile %q", project.CloudConfig, project.CloudConfig, profile)
			}
		})
	}
}

func TestFindAndPopulateMavenProjectsPreservesRecursiveWalkAppendPartialAndErrorBehavior(t *testing.T) {
	runContextStateSubtest(t, "ordered matches append partial project and exclude legacy paths", func(t *testing.T) {
		isolateContextHome(t)
		root := t.TempDir()
		firstDir := writeContextProject(t, root, "a-first", true)
		writeBrokenContextProject(t, root, "b-broken")
		writeContextProject(t, root, "target/excluded", true)
		writeContextFile(t, filepath.Join(root, "flattened-pom.xml"), []byte(contextPomXML("excluded-flattened")))
		writeContextProjectFiles(t, root, "root")
		existing := config.Project{Path: "existing recursive project"}
		ctx := Context{Recursive: true, TargetDirectory: root, Projects: []config.Project{existing}}
		_, hook := captureContextLog(t)

		err := ctx.FindAndPopulateMavenProjects()

		if err != nil {
			t.Fatalf("recursive discovery returned an error: %v", err)
		}
		if len(ctx.Projects) == 0 {
			t.Fatal("recursive discovery project population is empty")
		}
		if len(ctx.Projects) != 4 {
			t.Fatalf("recursive discovery returned %d projects, want existing, first, partial, root: %#v",
				len(ctx.Projects), ctx.Projects)
		}
		wantPaths := []string{existing.Path, firstDir + string(os.PathSeparator), "", root + string(os.PathSeparator)}
		if len(wantPaths) == 0 {
			t.Fatal("recursive project path expectation population is empty")
		}
		for index, want := range wantPaths {
			if ctx.Projects[index].Path != want {
				t.Fatalf("recursive project %d path was %q, want ordered %q", index, ctx.Projects[index].Path, want)
			}
		}
		if ctx.Projects[1].Config.Name != "a-first" || !ctx.Projects[1].IsMavenProject() {
			t.Fatalf("first recursive project was incomplete: %#v", ctx.Projects[1])
		}
		if !reflect.DeepEqual(ctx.Projects[2], config.Project{}) {
			t.Fatalf("recursive project-load error no longer appends the exact zero partial project: %#v", ctx.Projects[2])
		}
		if ctx.Projects[3].Config.Name != "root" || !ctx.Projects[3].IsMavenProject() {
			t.Fatalf("root recursive project was incomplete: %#v", ctx.Projects[3])
		}
		records, populationErr := hook.assertedRecords()
		if populationErr != nil {
			t.Fatal(populationErr)
		}
		if len(records) != 1 || records[0].level != logrus.WarnLevel ||
			!strings.Contains(records[0].message, "b-broken") ||
			!strings.Contains(records[0].message, "directory detected, but language was not set") {
			t.Fatalf("recursive partial-result logs were %#v", records)
		}
	})

	runContextStateSubtest(t, "empty directory returns nil without appending", func(t *testing.T) {
		isolateContextHome(t)
		existing := config.Project{Path: "existing empty-walk project"}
		ctx := Context{Recursive: true, TargetDirectory: t.TempDir(), Projects: []config.Project{existing}}

		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			t.Fatalf("empty recursive discovery returned an error: %v", err)
		}
		if !reflect.DeepEqual(ctx.Projects, []config.Project{existing}) {
			t.Fatalf("empty recursive discovery changed projects: %#v", ctx.Projects)
		}
	})

	runContextStateSubtest(t, "missing recursive root is suppressed by the existing walk callback", func(t *testing.T) {
		isolateContextHome(t)
		missing := filepath.Join(t.TempDir(), "missing recursive root")
		existing := config.Project{Path: "existing walk-error project"}
		ctx := Context{Recursive: true, TargetDirectory: missing, Projects: []config.Project{existing}}

		err := ctx.FindAndPopulateMavenProjects()

		if err != nil {
			t.Fatalf("missing recursive root returned %v, want suppressed nil", err)
		}
		if !reflect.DeepEqual(ctx.Projects, []config.Project{existing}) {
			t.Fatalf("missing recursive root changed projects: %#v", ctx.Projects)
		}
	})
}

func TestOnEachMavenProjectPreservesEmptyErrorAndCompleteTraversalJobOrder(t *testing.T) {
	runContextStateSubtest(t, "empty projects log and return before repository selection", func(t *testing.T) {
		ctx := Context{LocalConfig: config.OpenLocalConfig(filepath.Join(t.TempDir(), "must not be read"))}
		_, hook := captureContextLog(t)

		ctx.OnEachMavenProject("must not run", func(maven.Repository, config.Project) error {
			t.Fatal("empty Context invoked a Maven project job")
			return nil
		})

		assertContextLogs(t, hook, []contextLogRecord{{
			level: logrus.ErrorLevel, message: "could not find any pom models in the context",
		}})
	})

	runContextStateSubtest(t, "projects defaults skips jobs and errors continue in delivered order", func(t *testing.T) {
		repositoryURL := "https://repository.example.invalid/context"
		local := writeContextLocalConfig(t, repositoryURL, "context-user", "context-password")
		firstOutput := filepath.Join(t.TempDir(), "first-pom.xml")
		thirdOutput := filepath.Join(t.TempDir(), "third-pom.xml")
		first := contextMavenProject("first project", firstOutput)
		first.GitInfo = config.GitInfo{IsRepo: true, IsDirty: true}
		first.Config.Settings.DisableUpgradesFor = []config.Artifact{{GroupId: "existing.group", ArtifactId: "existing"}}
		first.Config.Settings.MaxVersionForDependencies = []config.MaxArtifact{{
			Artifact: config.Artifact{GroupId: "existing.group", ArtifactId: "capped"}, MaxVersion: "1",
		}}
		firstCloud := &contextCloudDefaults{defaults: config.CloudProjectDefaults{Settings: config.ProjectSettings{
			DisableDependencySort:    true,
			DisableKotlinUpgrade:     true,
			DisableSpringBootUpgrade: true,
			DisableUpgradesFor: []config.Artifact{{
				GroupId: "default.group", ArtifactId: "default",
			}},
			MaxVersionForDependencies: []config.MaxArtifact{{
				Artifact: config.Artifact{GroupId: "default.group", ArtifactId: "capped"}, MaxVersion: "2",
			}},
		}}}
		first.CloudConfig = firstCloud
		nilTypeError := errors.New("complete project-defaults dependency error")
		nilTypeCloud := &contextCloudDefaults{err: nilTypeError}
		nilType := config.Project{Path: "second nil-type project", CloudConfig: nilTypeCloud}
		third := contextMavenProject("third project", thirdOutput)
		projects := []config.Project{first, nilType, third}
		if len(projects) == 0 {
			t.Fatal("OnEachMavenProject project population is empty")
		}
		ctx := Context{DryRun: true, StealthMode: true, LocalConfig: local, Projects: projects}
		_, hook := captureContextLog(t)

		var events []string
		jobError := errors.New("complete OnEachMavenProject job error")
		firstJob := func(repository maven.Repository, project config.Project) error {
			events = append(events, "first:"+project.Path)
			if repository.Url != repositoryURL || repository.Auth == nil ||
				repository.Auth.Username != "context-user" || repository.Auth.Password != "context-password" {
				t.Fatalf("OnEachMavenProject repository was %#v", repository)
			}
			if project.Path == first.Path {
				if !project.Config.Settings.DisableDependencySort || !project.Config.Settings.DisableKotlinUpgrade ||
					!project.Config.Settings.DisableSpringBootUpgrade || !project.Config.Settings.UseStealthMode ||
					len(project.Config.Settings.DisableUpgradesFor) != 2 ||
					len(project.Config.Settings.MaxVersionForDependencies) != 2 {
					t.Fatalf("OnEachMavenProject did not merge defaults then stealth into delivered value: %#v",
						project.Config.Settings)
				}
			}
			project.Config.Name = "job-local mutation"
			return nil
		}
		errorJob := func(_ maven.Repository, project config.Project) error {
			events = append(events, "error:"+project.Path)
			return jobError
		}
		lastJob := func(_ maven.Repository, project config.Project) error {
			events = append(events, "last:"+project.Path)
			return nil
		}
		jobs := []func(maven.Repository, config.Project) error{firstJob, nil, errorJob, lastJob}
		if len(jobs) == 0 {
			t.Fatal("OnEachMavenProject job population is empty")
		}

		ctx.OnEachMavenProject("characterizing", jobs...)

		wantEvents := []string{
			"first:first project", "error:first project", "last:first project",
			"first:third project", "error:third project", "last:third project",
		}
		if !reflect.DeepEqual(events, wantEvents) {
			t.Fatalf("OnEachMavenProject job traversal order was %#v, want %#v", events, wantEvents)
		}
		if firstCloud.calls != 1 || nilTypeCloud.calls != 1 {
			t.Fatalf("project-default calls were first=%d nil-type=%d, want one each", firstCloud.calls, nilTypeCloud.calls)
		}
		if ctx.Projects[0].Config.Name != first.Config.Name ||
			ctx.Projects[0].Config.Settings.DisableDependencySort ||
			ctx.Projects[0].Config.Settings.UseStealthMode ||
			len(ctx.Projects[0].Config.Settings.DisableUpgradesFor) != 1 ||
			len(ctx.Projects[0].Config.Settings.MaxVersionForDependencies) != 1 {
			t.Fatalf("OnEachMavenProject leaked value-copy settings or job mutations: %#v", ctx.Projects[0])
		}

		wantLogs := []contextLogRecord{
			{level: logrus.DebugLevel, message: "using maven repository from local config"},
			{level: logrus.InfoLevel, message: "characterizing in first project"},
			{level: logrus.DebugLevel, message: "operating on a dirty git repo"},
			{level: logrus.WarnLevel, message: jobError.Error()},
			{level: logrus.WarnLevel, message: "could not find a project-defaults.json file in cloud-config"},
			{level: logrus.DebugLevel, message: nilTypeError.Error()},
			{level: logrus.WarnLevel, message: "no project type defined for path: second nil-type project"},
			{level: logrus.InfoLevel, message: "characterizing in third project"},
			{level: logrus.WarnLevel, message: jobError.Error()},
		}
		assertContextLogs(t, hook, wantLogs)
	})
}

func TestOnEachMavenProjectPreservesDryRunWriteAndWriteErrorBehavior(t *testing.T) {
	tests := []struct {
		name       string
		dryRun     bool
		wantChange bool
	}{
		{name: "dry run skips write", dryRun: true},
		{name: "non dry run writes mutated shared model", wantChange: true},
	}
	if len(tests) == 0 {
		t.Fatal("OnEachMavenProject dry-run characterization population is empty")
	}
	for _, test := range tests {
		runContextStateSubtest(t, test.name, func(t *testing.T) {
			local := writeContextLocalConfig(t, "https://write.example.invalid/maven", "", "")
			pomPath := filepath.Join(t.TempDir(), "pom.xml")
			initial := []byte(contextPomXML("before-each-write"))
			writeContextFile(t, pomPath, initial)
			project := contextMavenProject("write project", pomPath)
			ctx := Context{DryRun: test.dryRun, LocalConfig: local, Projects: []config.Project{project}}
			_, hook := captureContextLog(t)
			jobs := 0

			ctx.OnEachMavenProject("writing", func(_ maven.Repository, delivered config.Project) error {
				jobs++
				delivered.Type.Model().Name = "after-each-write"
				return nil
			})

			if jobs == 0 {
				t.Fatal("OnEachMavenProject recorded job population is empty")
			}
			contents, err := os.ReadFile(pomPath)
			if err != nil {
				t.Fatalf("read OnEachMavenProject output: %v", err)
			}
			if test.wantChange {
				if !strings.Contains(string(contents), "<name>after-each-write</name>") {
					t.Fatalf("non-dry-run output did not contain mutated model:\n%s", contents)
				}
			} else if !bytes.Equal(contents, initial) {
				t.Fatalf("dry-run output changed:\n got: %q\nwant: %q", contents, initial)
			}
			if project.Type.Model().Name != "after-each-write" {
				t.Fatal("project value copy no longer shares the Maven model pointer")
			}
			if _, populationErr := hook.assertedRecords(); populationErr != nil {
				t.Fatal(populationErr)
			}
		})
	}

	runContextStateSubtest(t, "write error is logged and not propagated", func(t *testing.T) {
		local := writeContextLocalConfig(t, "https://write-error.example.invalid/maven", "", "")
		pomPath := filepath.Join(t.TempDir(), "missing-parent", "pom.xml")
		project := contextMavenProject("write-error project", pomPath)
		ctx := Context{LocalConfig: local, Projects: []config.Project{project}}
		_, hook := captureContextLog(t)
		calls := 0

		ctx.OnEachMavenProject("writing errors", func(maven.Repository, config.Project) error {
			calls++
			return nil
		})

		if calls == 0 {
			t.Fatal("OnEachMavenProject write-error job population is empty")
		}
		wantError := (&os.PathError{Op: "open", Path: pomPath, Err: syscall.ENOENT}).Error()
		records, populationErr := hook.assertedRecords()
		if populationErr != nil {
			t.Fatal(populationErr)
		}
		if records[len(records)-1].level != logrus.WarnLevel || records[len(records)-1].message != wantError {
			t.Fatalf("OnEachMavenProject write error log was %#v, want %q", records[len(records)-1], wantError)
		}
	})
}

func TestOnRootProjectPreservesRootSelectionJobsDirtyLogDryRunAndWrites(t *testing.T) {
	runContextStateSubtest(t, "empty projects log and return", func(t *testing.T) {
		ctx := Context{}
		_, hook := captureContextLog(t)

		ctx.OnRootProject("must not run", func(config.Project) error {
			t.Fatal("empty Context invoked a root project job")
			return nil
		})

		assertContextLogs(t, hook, []contextLogRecord{{
			level: logrus.ErrorLevel, message: "could not find any pom models in the context",
		}})
	})

	runContextStateSubtest(t, "nil first project stops before later typed project", func(t *testing.T) {
		later := contextMavenProject("later project", filepath.Join(t.TempDir(), "later.xml"))
		ctx := Context{Projects: []config.Project{{Path: "nil root"}, later}}
		_, hook := captureContextLog(t)
		calls := 0

		ctx.OnRootProject("must not run", func(config.Project) error {
			calls++
			return nil
		})

		if calls != 0 {
			t.Fatalf("nil root invoked %d jobs", calls)
		}
		assertContextLogs(t, hook, []contextLogRecord{{
			level: logrus.ErrorLevel, message: "no project type defined for path: nil root",
		}})
	})

	tests := []struct {
		name       string
		dryRun     bool
		wantChange bool
	}{
		{name: "dry run skips root write", dryRun: true},
		{name: "non dry run writes root after all jobs", wantChange: true},
	}
	if len(tests) == 0 {
		t.Fatal("OnRootProject write characterization population is empty")
	}
	for _, test := range tests {
		runContextStateSubtest(t, test.name, func(t *testing.T) {
			pomPath := filepath.Join(t.TempDir(), "root-pom.xml")
			initial := []byte(contextPomXML("before-root-write"))
			writeContextFile(t, pomPath, initial)
			root := contextMavenProject("selected root", pomPath)
			root.Config.Name = "original root config"
			root.GitInfo = config.GitInfo{IsRepo: true, IsDirty: true}
			later := contextMavenProject("ignored later project", filepath.Join(t.TempDir(), "ignored.xml"))
			ctx := Context{DryRun: test.dryRun, Projects: []config.Project{root, later}}
			_, hook := captureContextLog(t)
			jobError := errors.New("complete OnRootProject job error")
			var events []string
			firstJob := func(project config.Project) error {
				events = append(events, "first:"+project.Path)
				project.Config.Name = "job-local root mutation"
				project.Type.Model().Name = "after-root-write"
				return nil
			}
			errorJob := func(project config.Project) error {
				events = append(events, "error:"+project.Path)
				return jobError
			}
			lastJob := func(project config.Project) error {
				events = append(events, "last:"+project.Path)
				return nil
			}
			jobs := []func(config.Project) error{firstJob, nil, errorJob, lastJob}
			if len(jobs) == 0 {
				t.Fatal("OnRootProject job population is empty")
			}

			ctx.OnRootProject("root action", jobs...)

			wantEvents := []string{"first:selected root", "error:selected root", "last:selected root"}
			if !reflect.DeepEqual(events, wantEvents) {
				t.Fatalf("OnRootProject job order was %#v, want %#v", events, wantEvents)
			}
			if ctx.Projects[0].Config.Name != "original root config" ||
				ctx.Projects[1].Type.Model().Name != "ignored later project" {
				t.Fatalf("OnRootProject changed a project value or selected a later project: %#v", ctx.Projects)
			}
			contents, err := os.ReadFile(pomPath)
			if err != nil {
				t.Fatalf("read OnRootProject output: %v", err)
			}
			if test.wantChange {
				if !strings.Contains(string(contents), "<name>after-root-write</name>") {
					t.Fatalf("non-dry-run root output did not contain mutated model:\n%s", contents)
				}
			} else if !bytes.Equal(contents, initial) {
				t.Fatalf("dry-run root output changed:\n got: %q\nwant: %q", contents, initial)
			}
			wantLogs := []contextLogRecord{
				{level: logrus.InfoLevel, message: "root action for file " + pomPath},
				{level: logrus.WarnLevel, message: "operating on a dirty git repo"},
				{level: logrus.WarnLevel, message: jobError.Error()},
			}
			assertContextLogs(t, hook, wantLogs)
		})
	}

	runContextStateSubtest(t, "root write error is logged", func(t *testing.T) {
		pomPath := filepath.Join(t.TempDir(), "missing-root-parent", "pom.xml")
		ctx := Context{Projects: []config.Project{contextMavenProject("root write error", pomPath)}}
		_, hook := captureContextLog(t)

		ctx.OnRootProject("root write error")

		wantError := (&os.PathError{Op: "open", Path: pomPath, Err: syscall.ENOENT}).Error()
		records, populationErr := hook.assertedRecords()
		if populationErr != nil {
			t.Fatal(populationErr)
		}
		if records[len(records)-1].level != logrus.WarnLevel || records[len(records)-1].message != wantError {
			t.Fatalf("OnRootProject write error log was %#v, want %q", records[len(records)-1], wantError)
		}
	})
}

func TestLoadProfilePreservesAssignmentsCreationExistingContentAndTouchError(t *testing.T) {
	runContextStateSubtest(t, "missing local config creates default profile file", func(t *testing.T) {
		profiles := filepath.Join(t.TempDir(), "profiles")
		createContextDirectory(t, profiles)
		profile := filepath.Join(profiles, "new-profile")
		ctx := Context{
			LocalConfig: config.OpenLocalConfig("/old/local"),
			CloudConfig: &contextCloudDefaults{},
		}
		_, hook := captureContextLog(t)

		ctx.LoadProfile(profile)

		assertContextProfileAssignments(t, ctx, profile)
		contents, err := os.ReadFile(filepath.Join(profile, "local-config.yaml"))
		if err != nil {
			t.Fatalf("read created local config: %v", err)
		}
		if len(contents) == 0 {
			t.Fatal("created local-config content population is empty")
		}
		cfg, err := ctx.LocalConfig.Config()
		if err != nil {
			t.Fatalf("parse created local config: %v", err)
		}
		if cfg.CloudConfig.Git.Url != "https://github.com/devdimensionlab/plybuild-config.git" {
			t.Fatalf("created local config cloud URL was %q", cfg.CloudConfig.Git.Url)
		}
		info, err := os.Stat(filepath.Join(profile, "local-config.yaml"))
		if err != nil || info.Mode().Perm() != 0644 {
			t.Fatalf("created local config metadata was (%v, %v), want mode 0644", info, err)
		}
		assertContextLogs(t, hook, []contextLogRecord{{
			level: logrus.DebugLevel, message: "localConfig does not exists, touching a new file",
		}})
	})

	runContextStateSubtest(t, "existing local config is retained without touch log", func(t *testing.T) {
		profile := filepath.Join(t.TempDir(), "existing-profile")
		marker := []byte("existing local config bytes remain exact\n")
		writeContextFile(t, filepath.Join(profile, "local-config.yaml"), marker)
		ctx := Context{}
		_, hook := captureContextLog(t)

		ctx.LoadProfile(profile)

		assertContextProfileAssignments(t, ctx, profile)
		contents, err := os.ReadFile(filepath.Join(profile, "local-config.yaml"))
		if err != nil || !bytes.Equal(contents, marker) {
			t.Fatalf("existing local config was (%q, %v), want exact %q", contents, err, marker)
		}
		if len(hook.records) != 0 {
			t.Fatalf("existing local config emitted Context logs: %#v", hook.records)
		}
	})

	runContextStateSubtest(t, "missing profile parent logs exact directory-create failure and retains assignments", func(t *testing.T) {
		profile := filepath.Join(t.TempDir(), "missing-parent", "new-profile")
		ctx := Context{}
		_, hook := captureContextLog(t)

		ctx.LoadProfile(profile)

		assertContextProfileAssignments(t, ctx, profile)
		wantError := (&os.PathError{Op: "mkdir", Path: profile, Err: syscall.ENOENT}).Error()
		assertContextLogs(t, hook, []contextLogRecord{
			{level: logrus.DebugLevel, message: "localConfig does not exists, touching a new file"},
			{level: logrus.ErrorLevel, message: wantError},
		})
	})
}

func TestLoadProfileCloudConfigOpenerPreservesProductionSelectionAndPublicCacheMapping(t *testing.T) {
	dependencies := systemLoadProfileDependencies()
	if reflect.TypeOf(dependencies.Opener) != reflect.TypeOf(gitCloudConfigOpener{}) {
		t.Fatalf("production cloud-config opener was %T, want %T", dependencies.Opener, gitCloudConfigOpener{})
	}

	profile := filepath.Join(t.TempDir(), `complete profile path with spaces and backslash\segment-ø`)
	cloud := dependencies.Open(profile)
	gitCloud, ok := cloud.(config.GitCloudConfig)
	if !ok {
		t.Fatalf("production cloud-config opener returned %T, want config.GitCloudConfig", cloud)
	}
	wantCache := filepath.Join(profile, "cloud-config")
	if gitCloud.Implementation().Dir() != wantCache {
		t.Fatalf("production cloud-config cache path was %q, want %q", gitCloud.Implementation().Dir(), wantCache)
	}
}

func TestLoadProfileHelperDeliversCompletePathOnceAndAssignsExactCloudConfigIdentity(t *testing.T) {
	runContextStateSubtest(t, "injected opener preserves state", func(t *testing.T) {
		profile := filepath.Join(t.TempDir(), `complete profile path with spaces and backslash\segment-β`)
		marker := []byte("existing local config remains exact through injected opener\n")
		writeContextFile(t, filepath.Join(profile, "local-config.yaml"), marker)
		cloud := &contextCloudDefaults{}
		opener := &contextCloudConfigOpener{cloud: cloud}
		ctx := Context{
			LocalConfig: config.OpenLocalConfig("/old/local"),
			CloudConfig: &contextCloudDefaults{},
		}
		_, hook := captureContextLog(t)

		ctx.loadProfile(loadProfileDependencies{Opener: opener}, profile)

		paths, err := opener.assertedProfilePaths()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(paths, []string{profile}) {
			t.Fatalf("cloud-config opener profile paths were %#v, want one complete path %q", paths, profile)
		}
		if ctx.CloudConfig != cloud {
			t.Fatalf("LoadProfile helper stored %T %#x, want exact returned %T %#x",
				ctx.CloudConfig, contextInterfacePointer(ctx.CloudConfig), cloud, contextInterfacePointer(cloud))
		}
		if ctx.LocalConfig.Implementation().Path != profile {
			t.Fatalf("LoadProfile helper local config path was %q, want %q", ctx.LocalConfig.Implementation().Path, profile)
		}
		contents, err := os.ReadFile(filepath.Join(profile, "local-config.yaml"))
		if err != nil || !bytes.Equal(contents, marker) {
			t.Fatalf("LoadProfile helper existing local config was (%q, %v), want exact %q", contents, err, marker)
		}
		if len(hook.records) != 0 {
			t.Fatalf("LoadProfile helper existing local config emitted logs: %#v", hook.records)
		}
	})
}

func TestLoadProfileDependenciesZeroValueReturnsNilWithoutTouchingPath(t *testing.T) {
	profile := filepath.Join(t.TempDir(), "zero opener must remain absent")
	if cloud := (loadProfileDependencies{}).Open(profile); cloud != nil {
		t.Fatalf("zero cloud-config opener returned %T %#v, want nil", cloud, cloud)
	}
	if _, err := os.Lstat(profile); !os.IsNotExist(err) {
		t.Fatalf("zero cloud-config opener touched profile path %q: %v", profile, err)
	}
}

func TestGetMavenRepositoryPreservesConfiguredAuthAndLegacyDefaultSelection(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantAuth bool
	}{
		{name: "both credentials create auth", username: "user-before-password", password: "password-after-user", wantAuth: true},
		{name: "username alone remains anonymous", username: "username-only"},
		{name: "password alone remains anonymous", password: "password-only"},
	}
	if len(tests) == 0 {
		t.Fatal("configured Maven repository characterization population is empty")
	}
	for _, test := range tests {
		runContextStateSubtest(t, test.name, func(t *testing.T) {
			url := "https://configured.example.invalid/" + strings.ReplaceAll(test.name, " ", "-")
			ctx := Context{LocalConfig: writeContextLocalConfig(t, url, test.username, test.password)}
			_, hook := captureContextLog(t)

			repository := ctx.GetMavenRepository()

			if repository.Url != url || repository.Id != "" {
				t.Fatalf("configured Maven repository was %#v, want URL %q and empty ID", repository, url)
			}
			if test.wantAuth {
				if repository.Auth == nil || repository.Auth.Username != test.username ||
					repository.Auth.Password != test.password || repository.Auth.Encrypted {
					t.Fatalf("configured Maven auth was %#v", repository.Auth)
				}
			} else if repository.Auth != nil {
				t.Fatalf("partial configured credentials created auth: %#v", repository.Auth)
			}
			assertContextLogs(t, hook, []contextLogRecord{{
				level: logrus.DebugLevel, message: "using maven repository from local config",
			}})
		})
	}

	runContextStateSubtest(t, "missing local config warns then uses existing default Maven lookup", func(t *testing.T) {
		isolateContextHome(t)
		assertNoCurrentUserMavenSettings(t)
		localDir := filepath.Join(t.TempDir(), "missing-local-config")
		ctx := Context{LocalConfig: config.OpenLocalConfig(localDir)}
		_, hook := captureContextLog(t)

		repository := ctx.GetMavenRepository()

		if repository != (maven.Repository{Url: "https://repo1.maven.org/maven2"}) {
			t.Fatalf("legacy default Maven repository was %#v, want exact fallback", repository)
		}
		configPath := filepath.Join(localDir, "local-config.yaml")
		wantConfigError := (&os.PathError{Op: "open", Path: configPath, Err: syscall.ENOENT}).Error()
		assertContextLogs(t, hook, []contextLogRecord{
			{level: logrus.WarnLevel, message: wantConfigError},
			{level: logrus.DebugLevel, message: "search for maven repository in .m2 folder \n"},
		})
	})
}

func runContextStateSubtest(t *testing.T, name string, action func(*testing.T)) {
	t.Helper()
	wantLogger := log
	wantEnvironment := snapshotContextEnvironment()
	wantDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("snapshot Context test working directory: %v", err)
	}
	t.Run(name, func(t *testing.T) {
		t.Cleanup(func() { SetLogger(wantLogger) })
		action(t)
	})
	if !sameContextFieldLogger(log, wantLogger) {
		t.Fatalf("Context logger leaked from %q: got %T %#x, want %T %#x", name,
			log, contextInterfacePointer(log), wantLogger, contextInterfacePointer(wantLogger))
	}
	if actual := snapshotContextEnvironment(); !reflect.DeepEqual(actual, wantEnvironment) {
		t.Fatalf("Context environment leaked from %q:\n got: %#v\nwant: %#v", name, actual, wantEnvironment)
	}
	actualDirectory, err := os.Getwd()
	if err != nil || actualDirectory != wantDirectory {
		t.Fatalf("Context working directory after %q was (%q, %v), want %q", name, actualDirectory, err, wantDirectory)
	}
}

func captureContextLog(t *testing.T) (*bytes.Buffer, *contextLogHook) {
	t.Helper()
	output := &bytes.Buffer{}
	hook := &contextLogHook{}
	testLogger := logrus.New()
	testLogger.SetLevel(logrus.DebugLevel)
	testLogger.SetOutput(output)
	testLogger.SetFormatter(contextMessageOnlyFormatter{})
	testLogger.AddHook(hook)
	SetLogger(testLogger.WithField("contract", "context"))
	return output, hook
}

func assertContextLogs(t *testing.T, hook *contextLogHook, want []contextLogRecord) {
	t.Helper()
	if len(want) == 0 {
		t.Fatal("Context expected-log population is empty")
	}
	actual, err := hook.assertedRecords()
	if err != nil {
		t.Fatal(err)
	}
	if len(actual) != len(want) {
		t.Fatalf("Context recorded %d logs, want %d:\n got: %#v\nwant: %#v", len(actual), len(want), actual, want)
	}
	for index := range want {
		if actual[index].level != want[index].level || actual[index].message != want[index].message {
			t.Fatalf("Context log %d was level=%s message=%q, want level=%s message=%q",
				index, actual[index].level, actual[index].message, want[index].level, want[index].message)
		}
		if !reflect.DeepEqual(actual[index].data, logrus.Fields{"contract": "context"}) {
			t.Fatalf("Context log %d fields were %#v, want exact selected logger field", index, actual[index].data)
		}
	}
}

func snapshotContextEnvironment() map[string]contextEnvironmentValue {
	values := make(map[string]contextEnvironmentValue)
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		value := ""
		if len(parts) == 2 {
			value = parts[1]
		}
		values[parts[0]] = contextEnvironmentValue{value: value, set: true}
	}
	return values
}

func isolateContextHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	homedir.Reset()
	t.Cleanup(homedir.Reset)
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, "xdg-config"))
	profile := filepath.Join(home, ".ply", "profiles", "default")
	createContextDirectory(t, profile)
	writeContextFile(t, filepath.Join(home, ".ply", "profiles", ".active_profile"), []byte("default"))
	return profile
}

func writeContextProject(t *testing.T, root, relative string, valid bool) string {
	t.Helper()
	directory := filepath.Join(root, filepath.FromSlash(relative))
	if valid {
		writeContextProjectFiles(t, directory, filepath.Base(directory))
	} else {
		writeBrokenContextProject(t, root, relative)
	}
	return directory
}

func writeContextProjectFiles(t *testing.T, directory, name string) {
	t.Helper()
	projectConfig := config.ProjectConfiguration{
		MavenProjectConfiguration: config.MavenProjectConfiguration{
			Artifact: config.Artifact{GroupId: "com.example", ArtifactId: name},
			Language: "java", Package: "com.example." + strings.ReplaceAll(name, "-", "_"),
			ApplicationName: "CompleteApplication",
		},
		Name: name,
	}
	contents, err := json.Marshal(projectConfig)
	if err != nil {
		t.Fatalf("marshal Context project config: %v", err)
	}
	writeContextFile(t, filepath.Join(directory, "ply.json"), contents)
	writeContextFile(t, filepath.Join(directory, "pom.xml"), []byte(contextPomXML(name)))
}

func writeBrokenContextProject(t *testing.T, root, relative string) string {
	t.Helper()
	directory := filepath.Join(root, filepath.FromSlash(relative))
	createContextDirectory(t, filepath.Join(directory, "src"))
	writeContextFile(t, filepath.Join(directory, "pom.xml"), []byte(contextPomXML("broken")))
	return directory
}

func contextPomXML(name string) string {
	return "<project>\n" +
		"  <modelVersion>4.0.0</modelVersion>\n" +
		"  <groupId>com.example</groupId>\n" +
		"  <artifactId>" + name + "</artifactId>\n" +
		"  <version>1.0.0</version>\n" +
		"  <name>" + name + "</name>\n" +
		"</project>\n"
}

func contextMavenProject(path, pomFile string) config.Project {
	return config.Project{
		Path: path,
		Type: config.MavenProject{
			PomFile: pomFile,
			PomModel: &pom.Model{
				ModelVersion: "4.0.0", GroupId: "com.example", ArtifactId: "context-project",
				Version: "1.0.0", Name: path,
			},
		},
		Config: config.ProjectConfiguration{Name: path},
	}
}

func writeContextLocalConfig(t *testing.T, url, username, password string) config.LocalConfigDir {
	t.Helper()
	directory := t.TempDir()
	contents := fmt.Sprintf("nexus:\n  url: %q\n  username: %q\n  password: %q\n", url, username, password)
	writeContextFile(t, filepath.Join(directory, "local-config.yaml"), []byte(contents))
	return config.OpenLocalConfig(directory)
}

func writeContextFile(t *testing.T, path string, contents []byte) {
	t.Helper()
	createContextDirectory(t, filepath.Dir(path))
	if err := testutil.WriteFileOutsideWorkingTree(path, contents, 0644); err != nil {
		t.Fatalf("write Context fixture %s: %v", path, err)
	}
}

func createContextDirectory(t *testing.T, path string) {
	t.Helper()
	if err := testutil.CopyFSOutsideWorkingTree(path, os.DirFS(t.TempDir())); err != nil {
		t.Fatalf("create Context fixture directory %s: %v", path, err)
	}
}

func assertContextProfileAssignments(t *testing.T, ctx Context, profile string) {
	t.Helper()
	if ctx.LocalConfig.Implementation().Path != profile ||
		ctx.LocalConfig.FilePath() != filepath.Join(profile, "local-config.yaml") {
		t.Fatalf("LoadProfile local config was %#v, want profile %q", ctx.LocalConfig, profile)
	}
	cloud, ok := ctx.CloudConfig.(config.GitCloudConfig)
	if !ok || cloud.Implementation().Dir() != filepath.Join(profile, "cloud-config") {
		t.Fatalf("LoadProfile cloud config was %T %#v, want profile %q", ctx.CloudConfig, ctx.CloudConfig, profile)
	}
}

func assertNoCurrentUserMavenSettings(t *testing.T) {
	t.Helper()
	current, err := user.Current()
	if err != nil {
		t.Fatalf("resolve current user before legacy Maven lookup: %v", err)
	}
	paths := []string{
		filepath.Join(current.HomeDir, ".m2", "settings.xml"),
		filepath.Join(current.HomeDir, "conf", "settings.xml"),
		filepath.Join(current.HomeDir, ".m2", "settings-security.xml"),
	}
	if len(paths) == 0 {
		t.Fatal("legacy Maven settings precondition population is empty")
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil || !os.IsNotExist(err) {
			t.Fatalf("legacy Maven lookup is not isolated: expected %q to be absent, stat returned %v", path, err)
		}
	}
}

func sameContextFieldLogger(actual, want logrus.FieldLogger) bool {
	if actual == nil || want == nil {
		return actual == nil && want == nil
	}
	actualValue := reflect.ValueOf(actual)
	wantValue := reflect.ValueOf(want)
	return actualValue.Type() == wantValue.Type() && actualValue.Pointer() == wantValue.Pointer()
}

func contextInterfacePointer(value interface{}) uintptr {
	if value == nil {
		return 0
	}
	return reflect.ValueOf(value).Pointer()
}
