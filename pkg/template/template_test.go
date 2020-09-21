package template

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"strings"
	"testing"
)

func newMockCloudConfig() (cfg config.GitCloudConfig) {
	cfg.Impl.Path = "test/cloud-config"
	return
}

func TestMergeTemplate(t *testing.T) {
	cfg := newMockCloudConfig()
	err := MergeTemplate(cfg, "test-template", "test/target-app")

	if err != nil {
		t.Errorf("%v\n", err)
	}
}

func TestReplacePathForSource(t *testing.T) {
	sourceDir := "test/cloud-config/templates/test-template"
	targetDir := "test/target-app"

	sourceConfig, err := config.InitProjectConfigurationFromDir(sourceDir)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	targetConfig, err := config.InitProjectConfigurationFromDir(targetDir)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	files, _ := FilesToCopy(sourceDir)
	for _, f := range files {
		if strings.Contains(f, ".kt") {
			sourceRelPath, err := file.RelPath(sourceDir, f)
			if err != nil {
				t.Errorf("%v\n", err)
			}
			sourceRelPath = ReplacePathForSource(sourceRelPath, sourceConfig, targetConfig)

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
