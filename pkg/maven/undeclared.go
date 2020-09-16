package maven

import (
	"co-pilot/pkg/logger"
)

func Undeclared(pomFile string) error {
	analyze, err := DependencyAnalyzeRaw(pomFile)
	if err != nil {
		return logger.ExternalError(err, analyze)
	}

	deps := DependencyAnalyze(analyze)

	for _, unused := range deps.UnusedDeclared {
		log.Infof("unused declared dependencies %s:%s", unused.GroupId, unused.ArtifactId)
	}

	for _, used := range deps.UsedUndeclared {
		log.Infof("used undeclared dependencies %s:%s", used.GroupId, used.ArtifactId)
	}

	return nil
}
