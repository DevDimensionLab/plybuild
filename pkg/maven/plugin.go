package maven

import (
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"strings"
)

type DependencyAnalyzeResult struct {
	UsedUndeclared []pom.Dependency
	UnusedDeclared []pom.Dependency
}

func DependencyAnalyze(rawOutput string) DependencyAnalyzeResult {
	var usedUndeclared []pom.Dependency
	var unusedDeclared []pom.Dependency

	used := false
	unused := false
	for _, line := range strings.Split(rawOutput, "\n") {

		if strings.Contains(line, "Used undeclared dependencies found:") {
			used = true
			unused = false
		}
		if strings.Contains(line, "Unused declared dependencies found:") {
			used = false
			unused = true
		}

		messageParts := strings.Split(line, "]")
		if len(messageParts) == 2 {
			artifactParts := strings.Split(strings.TrimSpace(messageParts[1]), ":")

			if len(artifactParts) == 5 {
				dependency := pom.Dependency{
					GroupId:    artifactParts[0],
					ArtifactId: artifactParts[1],
					Type_:      artifactParts[2],
					Version:    artifactParts[3],
					Scope:      artifactParts[4],
				}

				if used {
					usedUndeclared = append(usedUndeclared, dependency)
				}
				if unused {
					unusedDeclared = append(unusedDeclared, dependency)
				}
			}
		}
	}

	return DependencyAnalyzeResult{
		UsedUndeclared: usedUndeclared,
		UnusedDeclared: unusedDeclared,
	}
}

func UpgradePlugins() func(pair PomPair, args ...interface{}) error {
	return func(pair PomPair, args ...interface{}) error {
		return UpgradeKotlinOnModel(pair.Model)
	}
}

func UpgradePluginsOnModel(model *pom.Model) error {
	if model.Build == nil || model.Build.Plugins == nil {
		return nil
	}

	for _, plugin := range model.Build.Plugins.Plugin {
		if plugin.Version != "" {
			if err := upgradePlugin(model, plugin); err != nil {
				log.Warnf("%v", err)
			}
		}
	}
	return nil
}

func upgradePlugin(model *pom.Model, plugin pom.Plugin) error {
	currentVersionString, err := model.GetPluginVersion(plugin)
	if err != nil {
		return err
	}

	currentVersion, err := ParseVersion(currentVersionString)
	if err != nil {
		return err
	}

	metaData, err := GetMetaData(plugin.GroupId, plugin.ArtifactId)
	if err != nil {
		return err
	}

	latestRelease, err := metaData.LatestRelease()
	if err != nil {
		return err
	}

	if currentVersion != latestRelease {
		log.Warnf("outdated plugin %s:%s [%s => %s] \n", plugin.GroupId, plugin.ArtifactId, currentVersion.ToString(), latestRelease.ToString())
		if err := model.SetPluginVersion(plugin, latestRelease.ToString()); err != nil {
			return err
		}
	}
	return nil
}
