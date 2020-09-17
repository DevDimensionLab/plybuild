package maven

import (
	"co-pilot/pkg/logger"
	"co-pilot/pkg/shell"
	"errors"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"strings"
)

func ListUnusedAndUndeclared(pomFile string) error {
	analyze, err := runAnalyze(pomFile)
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

func GetSecondPartyGroupId(model *pom.Model) (string, error) {
	if model.GetGroupId() != "" {
		return GetFirstTwoPartsOfGroupId(model.GetGroupId())
	}

	return "", errors.New("could not extract 2party groupId")
}

func GetFirstTwoPartsOfGroupId(groupId string) (string, error) {
	parts := strings.Split(groupId, ".")

	if len(parts) <= 1 {
		return "", errors.New("groupId must at least contain two punctuations")
	} else {
		return strings.Join(parts[:2], "."), nil
	}
}

func IsSecondPartyGroupId(groupId string, secondPartyGroupId string) (bool, error) {
	groupIdParts := strings.Split(groupId, ".")
	secondPartyGroupIdParts := strings.Split(secondPartyGroupId, ".")

	if len(groupIdParts) <= 1 {
		return false, errors.New("groupId must at least contain two punctuations")
	} else {
		for i := range secondPartyGroupIdParts[:2] {
			if groupIdParts[i] != secondPartyGroupIdParts[i] {
				return false, nil
			}
		}
	}

	return true, nil
}

func runAnalyze(pomFile string) (string, error) {
	return shell.Run("mvn", "-f", pomFile, "dependency:analyze")
}
