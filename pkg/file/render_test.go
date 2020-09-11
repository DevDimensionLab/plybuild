package file

import (
	"testing"
)

func TestRender(t *testing.T) {
	err := Render("test_Dockerfile.render", "test_Dockerfile", struct {
		ArtifactId string
		GroupId    string
	}{
		ArtifactId: "testy",
		GroupId:    "testi",
	})
	if err != nil {
		t.Errorf("%v\n", err)
	}
}
