package upgrade

import (
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"strings"
)

type DependencySort struct {
	deps               []pom.Dependency
	secondPartyGroupId string
}

func (a DependencySort) Len() int      { return len(a.deps) }
func (a DependencySort) Swap(i, j int) { a.deps[i], a.deps[j] = a.deps[j], a.deps[i] }
func (a DependencySort) Less(i, j int) bool {
	return dependencySort(a.deps[i], a.deps[j], a.secondPartyGroupId)
}

func dependencySort(a pom.Dependency, b pom.Dependency, secondPartyGroupId string) bool {
	return concat(a, secondPartyGroupId) < concat(b, secondPartyGroupId)
}

func concat(dep pom.Dependency, secondPartyGroupId string) string {
	return fmt.Sprintf("%d:%s:%s", scopeWeight(dep.Scope), groupIdWeight(dep.GroupId, secondPartyGroupId), dep.ArtifactId)
}

func groupIdWeight(groupId string, secondPartyGroupId string) string {
	// implement custom groupId prefixing. example:
	if strings.Contains(groupId, secondPartyGroupId) {
		return fmt.Sprintf("%d-%s", 1, groupId)
	}

	return fmt.Sprintf("%d-%s", 100, groupId)
}

func scopeWeight(scope string) int {
	switch scope {
	case "":
		return 10
	case "compile":
		return 20
	case "provided":
		return 30
	case "runtime":
		return 40
	case "system":
		return 50
	case "import":
		return 60
	case "test":
		return 70

	default:
		return 100
	}
}
