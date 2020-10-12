package template

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"strings"
	"testing"
)

func newMockCloudConfig() (cfg config.GitCloudConfig) {
	cfg.Impl.Path = file.Path("test/cloud-config")
	return
}

func TestMergeTemplate_test_template(t *testing.T) {
	cfg := newMockCloudConfig()
	project, _ := config.InitProjectFromDirectory(file.Path("test/target-test-template"))
	template, _ := cfg.Template("test-template")
	err := MergeTemplate(template, project)

	if err != nil {
		t.Errorf("%v\n", err)
	}
}

func TestMergeTemplate_simple_template(t *testing.T) {
	cfg := newMockCloudConfig()
	project, _ := config.InitProjectFromDirectory(file.Path("test/target-simple-template"))
	template, _ := cfg.Template("simple-template")
	err := MergeTemplate(template, project)

	if err != nil {
		t.Errorf("%v\n", err)
	}
}

func TestReplacePathForSource(t *testing.T) {
	sourceDir := file.Path("test/cloud-config/templates/test-template")
	targetDir := file.Path("test/target-test-template")

	sourceConfig, err := config.InitProjectConfigurationFromDir(sourceDir)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	targetConfig, err := config.InitProjectConfigurationFromDir(targetDir)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	files, _ := filesToCopy(sourceDir)
	for _, f := range files {
		if strings.Contains(f, ".kt") {
			sourceRelPath, err := file.RelPath(sourceDir, f)
			if err != nil {
				t.Errorf("%v\n", err)
			}
			sourceRelPath = replacePathForSource(sourceRelPath, sourceConfig, targetConfig)

			expectedContains1 := "java"
			if !strings.Contains(sourceRelPath, expectedContains1) {
				t.Errorf("expectedContains1 %s not found in %s", expectedContains1, sourceRelPath)
			}

			expectedContains2 := "no/copilot/template/target"
			if !strings.Contains(sourceRelPath, expectedContains2) {
				t.Errorf("expectedContains1 %s not found in %s", expectedContains2, sourceRelPath)
			}
		}
	}
}
