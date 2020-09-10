package upgrade

import (
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
)

func Init(model *pom.Model, pomFile string) error {
	log.Infof("initializes project and writes to: %s", pomFile)
	return SortAndWrite(model, pomFile)
}
