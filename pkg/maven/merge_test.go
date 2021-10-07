package maven

import (
	"github.com/co-pilot-cli/mvn-pom-mutator/pkg/pom"
	"testing"
)

func TestMergeBuildPlugins(t *testing.T) {

	from, err := pom.GetModelFrom("test/merge/mergeBuildPluginFrom.xml")
	to, err := pom.GetModelFrom("test/merge/mergeBuildPluginTo.xml")

	if err != nil {
		t.Error(err)
	}

	err = mergeBuildPlugins(from, to)
	if err != nil {
		t.Error(err)
	}

	err = to.WriteToFile("test/merge/mergeBuildPluginMerged.xml")
	if err != nil {
		t.Error(err)
	}

	merged, err := pom.GetModelFrom("test/merge/mergeBuildPluginMerged.xml")
	if err != nil {
		t.Error(err)
	}

	var failed = true
	for _, mergedPlugin := range merged.Build.Plugins.Plugin {
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
