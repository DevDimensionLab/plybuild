package upgrade

import (
	"co-pilot/pkg/analyze"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"sort"
)

func SortAndWrite(model *pom.Model, file string) error {
	secondPartyGroupId, err := analyze.GetSecondPartyGroupId(model)
	if err != nil {
		return err
	}
	sort.Sort(DependencySort{
		deps:               model.Dependencies.Dependency,
		secondPartyGroupId: secondPartyGroupId})

	return model.WriteToFile(file)
}
