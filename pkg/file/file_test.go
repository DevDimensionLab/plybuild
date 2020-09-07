package file

import "testing"

func TestRelPath(t *testing.T) {
	relPath, err := RelPath(
		"/home/user/.co-pilot/cloud-config/templates/flyway-demo",
		"/home/user/.co-pilot/cloud-config/templates/flyway-demo/src/main/kotlin/no/copilot/template/demo/flyway/Queue.kt")

	expected := "src/main/kotlin/no/copilot/template/demo/flyway/Queue.kt"

	if err != nil {
		t.Errorf("%v\n", err)
	}

	if relPath != expected {
		t.Errorf("%s is not %s", relPath, expected)
	}
}
