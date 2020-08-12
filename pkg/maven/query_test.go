package maven

import (
	"strings"
	"testing"
)

func TestGetRepo(t *testing.T) {
	repos, err := GetRepositories()
	if err != nil {
		t.Errorf("%v", err)
	}
	localRepo, err := GetRepo(repos, true)
	if err != nil {
		t.Errorf("%v", err)
	}

	if !strings.Contains(localRepo, "http") {
		t.Errorf("local repo does not contain http")
	}

}
