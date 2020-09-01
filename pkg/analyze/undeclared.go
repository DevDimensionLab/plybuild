package analyze

import (
	"co-pilot/pkg/plugin"
	log "github.com/sirupsen/logrus"
)

func Undeclared(pomFile string) error {
	analyze, err := plugin.DependencyAnalyzeRaw(pomFile)
	if err != nil {
		return err
	}

	deps := plugin.DependencyAnalyze(analyze)

	for _, unused := range deps.UnusedDeclared {
		log.Infof("unused declared dependencies %s:%s", unused.GroupId, unused.ArtifactId)
	}

	for _, used := range deps.UsedUndeclared {
		log.Infof("used undeclared dependencies %s:%s", used.GroupId, used.ArtifactId)
	}

	return nil
}
