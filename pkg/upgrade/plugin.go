package upgrade

import (
	"co-pilot/pkg/maven"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
	"sort"
)

func Plugin(targetDirectory string, dryRun bool) error {
	pomFile := targetDirectory + "/pom.xml"
	model, err := pom.GetModelFrom(pomFile)
	if err != nil {
		return err
	}

	for _, plugin := range model.Build.Plugins.Plugin {
		if plugin.Version != "" {
			err = PluginUpgrade(model, plugin)
			if err != nil {
				log.Errorf("%v", err)
			}
		}
	}

	if !dryRun {
		sort.Sort(DependencySort(model.Dependencies.Dependency))
		return model.WriteToFile(pomFile)
	} else {
		return nil
	}
}

func PluginUpgrade(model *pom.Model, plugin pom.Plugin) error {
	currentVersion, err := model.GetPluginVersion(plugin)
	metaData, err := maven.GetMetaData(plugin.GroupId, plugin.ArtifactId)
	if err == nil {
		latestRelease := metaData.Versioning.Release
		if currentVersion != latestRelease {

			log.Warnf("outdated plugin %s:%s [%s => %s] \n", plugin.GroupId, plugin.ArtifactId, currentVersion, latestRelease)
			_ = model.SetPluginVersion(plugin, latestRelease)
		}
		return nil
	} else {
		return err
	}
}
