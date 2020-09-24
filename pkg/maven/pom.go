package maven

import (
	"co-pilot/pkg/config"
	"sort"
)

func SortAndWritePom(project config.Project, overwrite bool) error {
	var disableDepSort = project.Config.Settings.DisableDependencySort

	secondPartyGroupId, err := GetSecondPartyGroupId(project.PomModel)
	if err != nil {
		return err
	}

	if project.PomModel.Dependencies != nil && !disableDepSort {
		log.Infof("sorting pom file with default dependencySort")
		sort.Sort(DependencySort{
			deps:               project.PomModel.Dependencies.Dependency,
			secondPartyGroupId: secondPartyGroupId})
	}

	var writeToFile = project.PomFile
	if !overwrite {
		writeToFile = writeToFile + ".new"
	}

	log.Infof("writing model to pom file: %s", writeToFile)
	return project.PomModel.WriteToFile(writeToFile)
}
