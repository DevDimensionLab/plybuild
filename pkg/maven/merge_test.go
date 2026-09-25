package maven

import (
	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"testing"
)

func TestMergeBuildPlugins(t *testing.T) {

	from, err := pom.GetModelFrom("test/merge/mergeBuildPluginFrom.xml")
	if err != nil {
		t.Fatalf("load source POM fixture: %v", err)
	}
	to, err := pom.GetModelFrom("test/merge/mergeBuildPluginTo.xml")
	if err != nil {
		t.Fatalf("load target POM fixture: %v", err)
	}

	err = mergeBuildPlugins(from, to)
	if err != nil {
		t.Fatalf("merge build plugins: %v", err)
	}

	var failed = true
	for _, mergedPlugin := range to.Build.Plugins.Plugin {
		if mergedPlugin.GroupId == "org.springframework.boot" && mergedPlugin.ArtifactId == "spring-boot-maven-plugin" {
			if mergedPlugin.Configuration.AnyElements[0].XMLName.Local == "layers" {
				failed = false
			}
		}
	}

	if failed {
		t.Error("Failed to merge build plugin")
	}
}
