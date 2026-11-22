package knowledge

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// paddedArchive writes posts plus fillers up to size posts. Fillers share
// the words "fyll", "og", "i" and "plan".
func paddedArchive(t *testing.T, size int, posts ...fixture) string {
	t.Helper()
	for i := len(posts); i < size; i++ {
		posts = append(posts, fixture{id: fmt.Sprintf("k-20260200-filler-%02d", i), title: "Fyll", body: "Fyll og i plan."})
	}
	return writeArchive(t, posts...)
}

var chrome = fixture{id: "k-20260101-headless-chrome", title: "Headless Chrome tar skjermbilde for bredt", tags: []string{"chrome"}, body: "Bruk en annen innstilling."}

func requireRelevant(t *testing.T, dir, text string, repos []string, want ...string) {
	t.Helper()
	items, err := Relevant(dir, text, repos, "ply")
	if err != nil {
		t.Fatal(err)
	}
	if len(want) == 0 {
		if len(items) != 0 {
			t.Fatalf("items %v, want none", ids(items))
		}
		return
	}
	requireIDs(t, items, want...)
}

func TestRelevantNeedsTwoUncommonWordsInIDTagsOrTitle(t *testing.T) {
	t.Parallel()
	dir := paddedArchive(t, 10, chrome)
	requireRelevant(t, dir, "Skjermbilde i 320 px med headless Chrome og plan", nil, "k-20260101-headless-chrome")
	// One word is not enough.
	requireRelevant(t, dir, "Chrome i plan", nil)
}

func TestRelevantGivesNothingWithoutUncommonWords(t *testing.T) {
	t.Parallel()
	dir := paddedArchive(t, 10, chrome)
	// No words must never become the catalog of the whole archive.
	requireRelevant(t, dir, "", nil)
	requireRelevant(t, dir, "og i plan", nil)
}

func TestRelevantDropsWordsInMoreThanAFifthOfThePosts(t *testing.T) {
	t.Parallel()
	// "bredde" is in three posts out of ten, so it is common and does not count.
	dir := paddedArchive(t, 10, chrome,
		fixture{id: "k-20260102-bredde-a", title: "Bredde vindu skala"},
		fixture{id: "k-20260103-bredde-b", title: "Bredde vindu skala"},
		fixture{id: "k-20260104-bredde-c", title: "Bredde"})
	requireRelevant(t, dir, "bredde vindu skala", nil, "k-20260102-bredde-a", "k-20260103-bredde-b")
	requireRelevant(t, dir, "bredde vindu", nil)
}

func TestRelevantIgnoresHitsInTheBody(t *testing.T) {
	t.Parallel()
	dir := paddedArchive(t, 10, fixture{id: "k-20260101-other", title: "Noe annet", body: "Headless Chrome tar skjermbilde."})
	requireRelevant(t, dir, "headless chrome skjermbilde", nil)
}

func TestRelevantKeepsDefaultStatusesAndTheProject(t *testing.T) {
	t.Parallel()
	stale := chrome
	stale.status = "stale"
	other := chrome
	other.id, other.project = "k-20260102-headless-chrome", "trip"
	dir := paddedArchive(t, 10, stale, other)
	requireRelevant(t, dir, "headless chrome", nil)
}

func TestRelevantRanksOnWordsThenRepoAndShowsThree(t *testing.T) {
	t.Parallel()
	// Twenty posts, so four posts with a word are not more than a fifth.
	dir := paddedArchive(t, 20,
		fixture{id: "k-20260101-a", title: "alfa beta"},
		fixture{id: "k-20260102-b", title: "alfa beta", repos: []string{"trip"}},
		fixture{id: "k-20260103-c", title: "alfa beta gamma", repos: []string{"trip"}},
		fixture{id: "k-20260104-d", title: "alfa beta", repos: []string{"ply"}})
	// Three words first; then posts whose scope fits ply (or has no repos), in id order.
	requireRelevant(t, dir, "alfa beta gamma", []string{"ply"}, "k-20260103-c", "k-20260101-a", "k-20260104-d")
}

func TestRelevantToleratesBrokenPosts(t *testing.T) {
	t.Parallel()
	dir := paddedArchive(t, 10, chrome)
	if err := os.WriteFile(filepath.Join(dir, "k-20260109-broken.md"), []byte("no frontmatter"), 0o600); err != nil {
		t.Fatal(err)
	}
	requireRelevant(t, dir, "headless chrome", nil, "k-20260101-headless-chrome")
}
