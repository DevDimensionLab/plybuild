package sorting

import (
	"fmt"
	"github.com/co-pilot-cli/mvn-pom-mutator/pkg/pom"
	"strings"
)

type DependencySort struct {
	Deps    []pom.Dependency
	SortKey string
}

func (a DependencySort) Len() int      { return len(a.Deps) }
func (a DependencySort) Swap(i, j int) { a.Deps[i], a.Deps[j] = a.Deps[j], a.Deps[i] }
func (a DependencySort) Less(i, j int) bool {
	return dependencySort(a.Deps[i], a.Deps[j], a.SortKey)
}

func dependencySort(a pom.Dependency, b pom.Dependency, sortKey string) bool {
	return concat(a, sortKey) < concat(b, sortKey)
}

func concat(dep pom.Dependency, sortKey string) string {
	return fmt.Sprintf("%d:%s:%s", scopeWeight(dep.Scope), groupIdWeight(dep.GroupId, sortKey), dep.ArtifactId)
}

func groupIdWeight(groupId string, sortKey string) string {
	if strings.Contains(groupId, sortKey) {
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
