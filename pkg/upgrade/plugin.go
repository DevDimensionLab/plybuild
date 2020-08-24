package upgrade

import (
	"co-pilot/pkg/maven"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
)

func Plugin(model *pom.Model) error {
	if model.Build == nil || model.Build.Plugins == nil {
		return nil
	}

	for _, plugin := range model.Build.Plugins.Plugin {
		if plugin.Version != "" {
			err := PluginUpgrade(model, plugin)
			if err != nil {
				log.Errorf("%v", err)
			}
		}
	}

	return nil
}

func PluginUpgrade(model *pom.Model, plugin pom.Plugin) error {
	currentVersionString, err := model.GetPluginVersion(plugin)
	if err != nil {
		return err
	}

	currentVersion, err := maven.ParseVersion(currentVersionString)
	if err != nil {
		return err
	}

	metaData, err := maven.GetMetaData(plugin.GroupId, plugin.ArtifactId)
	if err != nil {
		return err
	}

	latestRelease, err := metaData.LatestRelease()
	if err != nil {
		return err
	}

	if currentVersion != latestRelease {
		log.Warnf("outdated plugin %s:%s [%s => %s] \n", plugin.GroupId, plugin.ArtifactId, currentVersion, latestRelease)
		_ = model.SetPluginVersion(plugin, latestRelease.ToString())
	}
	return nil
}
