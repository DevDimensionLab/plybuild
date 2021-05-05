package maven

import (
	"errors"
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/shell"
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

func isSecondPartyGroupId(groupId string, secondPartyGroupId string) (bool, error) {
	groupIdParts := strings.Split(groupId, ".")
	secondPartyGroupIdParts := strings.Split(secondPartyGroupId, ".")

	if len(groupIdParts) <= 1 {
		return false, errors.New(fmt.Sprintf(
			"secondParty groupId (%s) must at least contain two punctuations for comparison",
			groupId,
		))
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
