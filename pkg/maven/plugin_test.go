package maven

import (
	"co-pilot/pkg/file"
	"encoding/json"
	"testing"
)

func TestDependencyAnalyzeRawOutput(t *testing.T) {

	output, err := runAnalyze(file.Path("test/analyze/pom.xml"))

	if err != nil {
		t.Errorf("%v-> %s", err, output)
	}

	analyze := DependencyAnalyze(output)

	_, err = json.MarshalIndent(analyze, "", "  ")
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}

}
