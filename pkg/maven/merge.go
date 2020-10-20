package maven

import (
	"co-pilot/pkg/config"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"strings"
)

func MergePoms(from *pom.Model, to *pom.Model) error {
	if err := mergeDependencies(from, to); err != nil {
		log.Warnln(err)
	}

	if err := mergeManagementDependencies(from, to); err != nil {
		log.Warnln(err)
	}

	if err := mergeBuild(from, to); err != nil {
		log.Warnln(err)
	}

	return nil
}

func mergeDependencies(from *pom.Model, to *pom.Model) error {
	if from.Dependencies == nil {
		log.Debug("from dependencies is nil")
		return nil
	}

	if to.Dependencies == nil {
		log.Infof("inserting dependencies block into project")
		to.Dependencies = from.Dependencies
		return nil
	}

	for _, fromDep := range from.Dependencies.Dependency {
		if fromDep.GroupId == from.GetGroupId() {
			log.Infof("ignoring merge of dependency %s:%s due to dependency groupId equals project groupId", fromDep.GroupId, fromDep.ArtifactId)
			continue
		}
		var hasDependency = false
		for _, toDep := range to.Dependencies.Dependency {
			if fromDep.GroupId == toDep.GroupId && fromDep.ArtifactId == toDep.ArtifactId {
				hasDependency = true
			}
		}
		if !hasDependency {
			log.Infof("inserting dependency %s:%s into project", fromDep.GroupId, fromDep.ArtifactId)
			to.Dependencies.Dependency = append(to.Dependencies.Dependency, fromDep)
			mergePropertyKey(from, to, fromDep.Version)
		}
	}

	return nil
}

func mergeManagementDependencies(from *pom.Model, to *pom.Model) error {
	if from.DependencyManagement == nil {
		log.Debug("from dependencyManagement is nil")
		return nil
	}

	if to.DependencyManagement == nil {
		log.Infof("inserting dependency management block into project")
		from.DependencyManagement = to.DependencyManagement
		return nil
	}

	for _, fromDepMan := range from.DependencyManagement.Dependencies.Dependency {
		var hasManagementDependency = false
		for _, toDepMan := range to.DependencyManagement.Dependencies.Dependency {
			if fromDepMan.GroupId == toDepMan.GroupId && fromDepMan.ArtifactId == toDepMan.ArtifactId {
				hasManagementDependency = true
			}
		}
		if !hasManagementDependency {
			log.Infof("inserting management dependency %s:%s into project", fromDepMan.GroupId, fromDepMan.ArtifactId)
			to.DependencyManagement.Dependencies.Dependency = append(to.DependencyManagement.Dependencies.Dependency, fromDepMan)
			mergePropertyKey(from, to, fromDepMan.Version)
		}
	}

	return nil
}

func mergeBuild(from *pom.Model, to *pom.Model) error {
	if from.Build == nil {
		log.Debug("from build is nil")
		return nil
	}

	if to.Build == nil {
		to.Build = from.Build
		log.Infof("inserting build block into project")
		return nil
	}

	if to.Build.FinalName == "" && from.Build.FinalName != "" {
		to.Build.FinalName = from.Build.FinalName
		log.Infof("inserting <finalName>%s</finalName>", from.Build.FinalName)
	}

	if err := mergeBuildPlugins(from, to); err != nil {
		return err
	}

	if err := mergeBuildPluginManagement(from, to); err != nil {
		return err
	}

	return nil
}

func mergeBuildPlugins(from *pom.Model, to *pom.Model) error {
	if from.Build.Plugins == nil {
		log.Debug("from build.plugin is nil")
		return nil
	}

	if to.Build.Plugins == nil {
		to.Build.Plugins = from.Build.Plugins
		log.Infof("inserting build.plugins block into project")
		return nil
	}

	for i, fromPlugin := range from.Build.Plugins.Plugin {
		var hasPlugin = false
		for j, toPlugin := range to.Build.Plugins.Plugin {
			if fromPlugin.GroupId == toPlugin.GroupId && fromPlugin.ArtifactId == toPlugin.ArtifactId {
				hasPlugin = true
				mergeBuildPluginExecutions(&from.Build.Plugins.Plugin[i], &to.Build.Plugins.Plugin[j])
			}
		}
		if !hasPlugin {
			log.Infof("inserting plugin %s:%s into project", fromPlugin.GroupId, fromPlugin.ArtifactId)
			to.Build.Plugins.Plugin = append(to.Build.Plugins.Plugin, fromPlugin)
			mergePropertyKey(from, to, fromPlugin.Version)
			if fromPlugin.Dependencies != nil {
				for _, fromPlugDep := range fromPlugin.Dependencies.Dependency {
					mergePropertyKey(from, to, fromPlugDep.Version)
				}
			}
		}
	}

	return nil
}

func mergeBuildPluginManagement(from *pom.Model, to *pom.Model) error {
	if from.Build.PluginManagement == nil {
		log.Debug("from build.pluginManagement is nil")
		return nil
	}

	if to.Build.PluginManagement == nil {
		to.Build.PluginManagement = from.Build.PluginManagement
		log.Infof("inserting build.pluginManagement block into project")
		for _, fromPluginMan := range from.Build.PluginManagement.Plugins.Plugin {
			mergePropertyKey(from, to, fromPluginMan.Version)
		}
		return nil
	}

	for i, fromPluginMan := range from.Build.PluginManagement.Plugins.Plugin {
		var hasPlugin = false
		for j, toPluginMan := range to.Build.PluginManagement.Plugins.Plugin {
			if fromPluginMan.GroupId == toPluginMan.GroupId && fromPluginMan.ArtifactId == toPluginMan.ArtifactId {
				hasPlugin = true
				mergeBuildPluginExecutions(&from.Build.PluginManagement.Plugins.Plugin[i], &to.Build.PluginManagement.Plugins.Plugin[j])
			}
		}
		if !hasPlugin {
			log.Infof("inserting plugin management %s:%s into project", fromPluginMan.GroupId, fromPluginMan.ArtifactId)
			to.Build.PluginManagement.Plugins.Plugin = append(to.Build.PluginManagement.Plugins.Plugin, fromPluginMan)
			mergePropertyKey(from, to, fromPluginMan.Version)
		}
	}

	return nil
}

func mergeBuildPluginExecutions(from *pom.Plugin, to *pom.Plugin) {
	if from.Executions == nil {
		return
	}

	if to.Executions == nil {
		log.Infof("merging all plugin executions into plugin %s:%s", to.GroupId, to.ArtifactId)
		to.Executions = from.Executions
		return
	}

	for _, fromExecution := range from.Executions.Execution {
		var hasExecution = false
		for _, toExecution := range to.Executions.Execution {
			if fromExecution.Id == toExecution.Id {
				hasExecution = true
			}
		}

		if !hasExecution {
			log.Infof("merging execution %s into plugin %s:%s", fromExecution.Id, to.GroupId, to.ArtifactId)
		}
	}

	return
}

func mergePropertyKey(from *pom.Model, to *pom.Model, version string) {
	if version != "" && strings.HasPrefix(version, "${") {
		versionKey := strings.Trim(version, "${}")
		if from.Properties != nil {
			for _, version := range from.Properties.AnyElements {
				if version.XMLName.Local == versionKey {
					to.Properties.AnyElements = append(to.Properties.AnyElements, version)
				}
			}
		}
	}
}

func MergeAndWritePomFiles(source config.Project, target config.Project) error {
	log.Infof(fmt.Sprintf("merging %s into %s", source.Type.FilePath(), target.Type.FilePath()))
	if err := MergePoms(source.Type.Model(), target.Type.Model()); err != nil {
		return err
	}

	return target.SortAndWritePom()
}
