package analyze

import (
	"errors"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"strings"
)

func GetSecondPartyGroupId(model *pom.Model) (string, error) {
	if model.GroupId != "" {
		return GetFirstTwoPartsOfGroupId(model.GroupId)
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
