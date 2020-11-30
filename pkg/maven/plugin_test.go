package maven

import (
	"github.com/co-pilot-cli/co-pilot/pkg/file"
	"encoding/json"
	"testing"
)

func TestDependencyAnalyzeRawOutput(t *testing.T) {

	output := runAnalyze(file.Path("test/analyze/pom.xml"))

	if output.Err != nil {
		t.Errorf("%v-> %s", output.Err, output.String())
	}

	analyze := DependencyAnalyze(output.StdOut.String())

	_, err := json.MarshalIndent(analyze, "", "  ")
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}

}
