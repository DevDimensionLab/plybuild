package upgrade

import (
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
)

type DependencySort []pom.Dependency

func (a DependencySort) Len() int      { return len(a) }
func (a DependencySort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a DependencySort) Less(i, j int) bool {
	return dependencySort(a[i], a[j])
}

func dependencySort(a pom.Dependency, b pom.Dependency) bool {
	return concat(a) < concat(b)
}

func concat(dep pom.Dependency) string {
	return fmt.Sprintf("%d:%s:%s", scopeWeight(dep.Scope), groupIdWeight(dep.GroupId), dep.ArtifactId)
}

func groupIdWeight(groupId string) string {
	// implement custom groupId prefixing. example:
	//if strings.Contains(groupId, "com.example") {
	//	return fmt.Sprintf("%d-%s", 1, groupId)
	//}

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
