package knowledge

import "sort"

// Selection of the posts a plan, assignment or Goal points at. The values
// are starting points; measured use of the posts decides whether they hold.
const (
	relevantLimit = 3
	// A word in more than commonPercent of all posts says nothing about a
	// post, so it is dropped; this removes common words in any language
	// without a stop list.
	commonPercent = 20
	// A post is relevant when this many uncommon words hit its id, tags or title.
	relevantWords = 2
)

// Relevant returns at most three posts under dir that text points at. Words
// of text that hit more than a fifth of all posts in any field are dropped;
// a post is relevant when at least two remaining words hit its id, tags or
// title. Posts rank on distinct words, then on a scope that fits one of
// repos (or has no repos), then on id. Only posts of project and with one of
// DefaultStatuses are kept. Text without such words gives no posts, never
// the whole archive. An error means dir itself cannot be read; unreadable
// posts are left out.
func Relevant(dir, text string, repos []string, project string) ([]Item, error) {
	posts, _, err := read(dir)
	if err != nil {
		return nil, err
	}
	words := uncommon(tokens(text), posts)
	filter := Query{Project: project, Statuses: DefaultStatuses}
	var found []hit
	for _, p := range posts {
		if !filter.keeps(p.item) {
			continue
		}
		h := hit{item: p.item, words: p.headHits(words), repoFits: fitsAny(p.item.Scope.Repos, repos)}
		if h.words >= relevantWords {
			found = append(found, h)
		}
	}
	sort.Slice(found, func(i, j int) bool { return found[i].before(found[j]) })
	if len(found) > relevantLimit {
		found = found[:relevantLimit]
	}
	var items []Item
	for _, h := range found {
		items = append(items, h.item)
	}
	return items, nil
}

// uncommon keeps the words that hit at most commonPercent of posts.
func uncommon(words []string, posts []post) []string {
	var out []string
	for _, word := range words {
		if mentions(posts, word)*100 <= commonPercent*len(posts) {
			out = append(out, word)
		}
	}
	return out
}

// fitsAny reports whether a post scoped to scope fits one of repos; a post
// without repos fits every repository, and no repos rank nothing.
func fitsAny(scope, repos []string) bool {
	if len(repos) == 0 || len(scope) == 0 {
		return true
	}
	for _, repo := range repos {
		if contains(scope, repo) {
			return true
		}
	}
	return false
}

// headHits counts the words that hit p's id, tags or title.
func (p post) headHits(words []string) int {
	count := 0
	for _, word := range words {
		if hits(p.head, word) {
			count++
		}
	}
	return count
}

// mentions counts the posts that word hits in any field.
func mentions(posts []post, word string) int {
	count := 0
	for _, p := range posts {
		if hits(p.head, word) || hits(p.situation, word) || hits(p.body, word) {
			count++
		}
	}
	return count
}
