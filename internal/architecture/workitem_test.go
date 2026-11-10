package architecture

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestWI01Architecture(t *testing.T) {
	_, here, _, _ := runtime.Caller(0)
	result := scanRepository(filepath.Clean(filepath.Join(filepath.Dir(here), "../..")), scanIO{})
	for _, problem := range result.problems {
		t.Error(problem)
	}
	for _, edge := range result.edges {
		parts := strings.Split(edge, " -> ")
		source, target := false, false
		for _, pkg := range affectedPackages {
			source = source || parts[0] == modulePath+"/"+pkg
			target = target || parts[1] == modulePath+"/"+pkg
		}
		if source && (target || parts[0] == modulePath+"/internal/workitem") {
			t.Log(edge)
		}
	}
}
