package maven

import (
	"strings"
	"testing"
)

func TestGetRepo(t *testing.T) {
	settings, _ := NewSettings()
	repos, err := settings.GetRepositories()
	if err != nil {
		t.Errorf("%v", err)
	}
	localRepo, err := repos.GetDefaultRepository()

	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(localRepo.Url, "http") {
		t.Errorf("local repo does not contain http")
	}

}
