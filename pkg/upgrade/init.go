package upgrade

import (
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
)

func Init(directory string) error {
	pomFile := directory + "/pom.xml"
	model, err := pom.GetModelFrom(pomFile)
	if err != nil {
		return err
	}

	log.Infof("Initializes project and writes to: %s", pomFile)
	return SortAndWrite(model, pomFile)
}
