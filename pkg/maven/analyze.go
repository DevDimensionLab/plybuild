package maven

import (
	"co-pilot/pkg/shell"
	"errors"
	"strings"
)

func ListUnusedAndUndeclared(pomFile string) error {
	analyze := runAnalyze(pomFile)
	if analyze.Err != nil {
		return analyze.FormatError()
	}

	deps := DependencyAnalyze(analyze.StdOut.String())

	for _, unused := range deps.UnusedDeclared {
		log.Infof("unused declared dependencies %s:%s", unused.GroupId, unused.ArtifactId)
	}

	for _, used := range deps.UsedUndeclared {
		log.Infof("used undeclared dependencies %s:%s", used.GroupId, used.ArtifactId)
	}

	return nil
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

func runAnalyze(pomFile string) shell.Output {
	return shell.Run("mvn", "-f", pomFile, "dependency:analyze")
}
