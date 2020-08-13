package upgrade

import (
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"strings"
)

type DependencySort struct {
	deps         []pom.Dependency
	localGroupId string
}

func (a DependencySort) Len() int      { return len(a.deps) }
func (a DependencySort) Swap(i, j int) { a.deps[i], a.deps[j] = a.deps[j], a.deps[i] }
func (a DependencySort) Less(i, j int) bool {
	return dependencySort(a.deps[i], a.deps[j], a.localGroupId)
}

func dependencySort(a pom.Dependency, b pom.Dependency, localGroupId string) bool {
	return concat(a, localGroupId) < concat(b, localGroupId)
}

func concat(dep pom.Dependency, localGroupId string) string {
	return fmt.Sprintf("%d:%s:%s", scopeWeight(dep.Scope), groupIdWeight(dep.GroupId, localGroupId), dep.ArtifactId)
}

func groupIdWeight(groupId string, localGroupId string) string {
	// implement custom groupId prefixing. example:
	if strings.Contains(groupId, localGroupId) {
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
