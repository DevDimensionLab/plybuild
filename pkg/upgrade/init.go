package upgrade

import (
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"sort"
)

func Init(directory string) error {
	pomFile := directory + "pom.xml"
	model, err := pom.GetModelFrom(pomFile)
	if err != nil {
		return err
	}

	fmt.Printf("[INFO] Initializes project and writes to: %s.new\n", pomFile)
	sort.Sort(DependencySort(model.Dependencies.Dependency))
	return model.WriteToFile(pomFile + ".new")
}
