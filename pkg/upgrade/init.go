package upgrade

import (
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
)

func Init(model *pom.Model, pomFile string) error {
	log.Infof("Initializes project and writes to: %s", pomFile)
	return SortAndWrite(model, pomFile)
}
