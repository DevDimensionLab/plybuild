package file

import (
	"testing"
)

func TestPath(t *testing.T) {
	path1 := Path("%s/co-pilot.json", "test")

	expected := "test/co-pilot.json"

	if path1 != expected {
		t.Errorf("expected %s, but got: %s", expected, path1)
	}
}
