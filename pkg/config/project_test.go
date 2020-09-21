package config

import "testing"

func TestProjectConfiguration_SourceMainPath(t *testing.T) {
	sourceConfig, err := InitProjectConfigurationFromDir("test/cloud-config/templates/test-template")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	targetConfig, err := InitProjectConfigurationFromDir("test/target-app")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	expectedSourceRoot := "src/main/kotlin/no/copilot/template/test"
	if sourceConfig.SourceMainPath() != expectedSourceRoot {
		t.Errorf("expected %s, but got instead %s", expectedSourceRoot, sourceConfig.SourceMainPath())
	}

	expectedTargetRoot := "src/main/java/no/copilot/template/target"
	if targetConfig.SourceMainPath() != expectedTargetRoot {
		t.Errorf("expected %s, but got instead %s", expectedTargetRoot, targetConfig.SourceMainPath())
	}
}
