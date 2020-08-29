package plugin

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDependencyAnalyzeRawOutput(t *testing.T) {

	output, err := DependencyAnalyzeRaw()

	if err != nil  {
		t.Errorf("%v-> %s", err, output )
	}

	analyze := DependencyAnalyze(output)

	e, err := json.MarshalIndent(analyze,"", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(e))
}
