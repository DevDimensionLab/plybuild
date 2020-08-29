package plugin

import (
	"bytes"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"os/exec"
	"strings"
)

type DependencyAnalyzeResult struct {
	UsedUndeclared []pom.Dependency
	UnusedDeclared []pom.Dependency
}
func DependencyAnalyze(rawOutput string) DependencyAnalyzeResult {
	var usedUndeclared []pom.Dependency
	var unusedDeclared []pom.Dependency

	used := false
	unused := false
	for _, line := range strings.Split(rawOutput, "\n") {

		if strings.Contains( line, "Used undeclared dependencies found:") {
			used  = true
			unused = false
		}
		if strings.Contains( line, "Unused declared dependencies found:") {
			used = false
			unused  = true
		}

		trimedLine := strings.TrimSpace(strings.Replace(line, "[WARNING]", "", 1))
		artifactParts := strings.Split(trimedLine, ":")

		if strings.HasPrefix(line, "[WARNING]    ") && len(artifactParts) == 5 {
			dependency := pom.Dependency{
				GroupId:    artifactParts[0],
				ArtifactId: artifactParts[1],
				Type_:      artifactParts[2],
				Version:    artifactParts[3],
				Scope:      artifactParts[4],
			}

			if used {
				usedUndeclared = append(usedUndeclared, dependency)
			}
			if unused {
				unusedDeclared = append(unusedDeclared, dependency)
			}
		}
	}

	return DependencyAnalyzeResult{
		UsedUndeclared: usedUndeclared,
		UnusedDeclared: unusedDeclared,
	}
}

func DependencyAnalyzeRaw() (string, error) {
	cmd := exec.Command("mvn", "-f", "test/pom.xml", "dependency:analyze")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		println( "error: " + errOut.String())
		return out.String(), err
	}

	return out.String(), err
}
