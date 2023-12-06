package file

import (
	"testing"
)

func TestPath(t *testing.T) {
	path1 := Path("%s/ply.json", "test")

	expected := "test/ply.json"

	if path1 != expected {
		t.Errorf("expected %s, but got: %s", expected, path1)
	}
}
