package maven

import (
	"sort"
)

func SortAndWritePom(wrapper PomWrapper, overwrite bool) error {
	var disableDepSort = wrapper.ProjectConfig.Settings.DisableDependencySort

	secondPartyGroupId, err := GetSecondPartyGroupId(wrapper.Model)
	if err != nil {
		return err
	}

	if wrapper.Model.Dependencies != nil && !disableDepSort {
		log.Infof("sorting pom file with default dependencySort")
		sort.Sort(DependencySort{
			deps:               wrapper.Model.Dependencies.Dependency,
			secondPartyGroupId: secondPartyGroupId})
	}

	var writeToFile = wrapper.PomFile
	if !overwrite {
		writeToFile = writeToFile + ".new"
	}

	log.Infof("writing model to pom file: %s", writeToFile)
	return wrapper.Model.WriteToFile(writeToFile)
}
