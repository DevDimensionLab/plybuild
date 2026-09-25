package maven

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/testutil"
)

func TestDependencyAnalyzeRawOutput(t *testing.T) {
	fakeMaven := filepath.Join(t.TempDir(), "mvn")
	script := `#!/bin/sh
printf 'argv'
for argument in "$@"; do
  printf '\t%s' "$argument"
done
printf '\n'
printf '%s\n' '[WARNING] Used undeclared dependencies found:'
printf '%s\n' '[WARNING]    com.example:used-library:jar:1.2.3:compile'
printf '%s\n' '[WARNING] Unused declared dependencies found:'
printf '%s\n' '[WARNING]    org.example:first-unused:jar:4.5.6:compile'
printf '%s\n' '[WARNING]    net.example:second-unused:jar:7.8.9:test'
`
	if err := testutil.WriteFileOutsideWorkingTree(fakeMaven, []byte(script), 0o700); err != nil {
		t.Fatalf("write Maven double: %v", err)
	}
	t.Setenv("PATH", filepath.Dir(fakeMaven)+string(os.PathListSeparator)+os.Getenv("PATH"))

	const pomFile = "test/analyze/pom.xml"
	output := runAnalyze(pomFile)
	if output.Err != nil {
		t.Fatalf("run dependency analysis: %v\n%s", output.Err, output.String())
	}
	parts := strings.SplitN(output.StdOut.String(), "\n", 2)
	if len(parts) != 2 {
		t.Fatalf("Maven double output had no analysis payload: %q", output.StdOut.String())
	}
	wantArgv := "argv\t-f\t" + pomFile + "\tdependency:analyze"
	if parts[0] != wantArgv {
		t.Fatalf("Maven argv was %q, want %q", parts[0], wantArgv)
	}

	rawOutput := parts[1]
	if !strings.Contains(rawOutput, `[WARNING] Used undeclared dependencies found:`) {
		t.Fatalf("Maven double produced no dependency population: %q", rawOutput)
	}

	analyze := DependencyAnalyze(rawOutput)
	if len(analyze.UsedUndeclared) != 1 {
		t.Fatalf("parsed %d used undeclared dependencies, want 1", len(analyze.UsedUndeclared))
	}
	if len(analyze.UnusedDeclared) != 2 {
		t.Fatalf("parsed %d unused declared dependencies, want 2", len(analyze.UnusedDeclared))
	}

	used := analyze.UsedUndeclared[0]
	if used.GroupId != "com.example" || used.ArtifactId != "used-library" || used.Version != "1.2.3" || used.Scope != "compile" {
		t.Errorf("parsed used dependency as %+v", used)
	}
	firstUnused := analyze.UnusedDeclared[0]
	secondUnused := analyze.UnusedDeclared[1]
	if firstUnused.ArtifactId != "first-unused" || secondUnused.ArtifactId != "second-unused" {
		t.Errorf("parsed unused dependencies in the wrong order: %+v", analyze.UnusedDeclared)
	}
}
