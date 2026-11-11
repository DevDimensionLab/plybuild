package agentlaunch

import (
	"os/exec"
	"strings"
	"testing"
)

func TestNoDomainDependencies(t *testing.T) {
	output, err := exec.Command("go", "list", "-deps", ".").Output()
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Fields(string(output)) {
		if strings.Contains(line, "github.com/devdimensionlab/plybuild/") && line != "github.com/devdimensionlab/plybuild/internal/agentlaunch" {
			t.Fatalf("domain dependency: %s", line)
		}
	}
}
