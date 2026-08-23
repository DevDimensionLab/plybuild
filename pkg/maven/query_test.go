package maven

import "testing"

func TestGetRepo(t *testing.T) {
	const fallbackURL = "https://repo.example.test/maven2"
	repos := Repositories{Fallback: Repository{Url: fallbackURL}}
	localRepo, err := repos.GetDefaultRepository()
	if err != nil {
		t.Fatalf("get fallback repository: %v", err)
	}

	if localRepo.Url != fallbackURL {
		t.Errorf("default repository URL was %q, want %q", localRepo.Url, fallbackURL)
	}
}
