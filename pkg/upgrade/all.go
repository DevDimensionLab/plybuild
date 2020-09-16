package upgrade

import "github.com/perottobc/mvn-pom-mutator/pkg/pom"

func All(model *pom.Model) error {
	if err := Kotlin(model); err != nil {
		log.Warn(err)
	}
	if err := SpringBoot(model); err != nil {
		log.Warn(err)
	}
	if err := Dependencies(model, true); err != nil {
		log.Warn(err)
	}
	if err := Dependencies(model, false); err != nil {
		log.Warn(err)
	}
	if err := Plugin(model); err != nil {
		log.Warn(err)
	}

	return nil
}
