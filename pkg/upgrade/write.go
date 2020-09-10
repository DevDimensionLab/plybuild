package upgrade

import (
	"co-pilot/pkg/analyze"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"sort"
)

func SortAndWrite(model *pom.Model, file string) error {
	log.Infof("sorting and writing to pom file: %s", file)
	secondPartyGroupId, err := analyze.GetSecondPartyGroupId(model)
	if err != nil {
		return err
	}

	if model.Dependencies != nil {
		sort.Sort(DependencySort{
			deps:               model.Dependencies.Dependency,
			secondPartyGroupId: secondPartyGroupId})
	}

	return model.WriteToFile(file)
}
