// Package knowledge reads a curated knowledge archive that another repository
// owns and finds posts in it. A post is a k-*.md file anywhere under the
// archive directory with frontmatter between two --- lines; its theme is its
// first directory under the archive. Ply tolerates the format and does not
// enforce it: only a k-*.md that is no regular file, unreadable frontmatter
// and a missing id make a file invalid, and other missing fields are empty.
// Nothing is written.
package knowledge

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// SituationPrefix starts the line of a post's body that is its situation.
const SituationPrefix = "**Situasjon:**"

// DefaultStatuses are the statuses shown when none are chosen; stale and
// superseded posts are left out.
var DefaultStatuses = []string{"candidate", "confirmed"}

// Query selects and ranks posts. Empty fields select everything.
type Query struct {
	// Words hit a post when one is a substring of a lowercase token of its
	// id, tags, title, situation or body. Without words every post is
	// listed in theme and id order, and Limit does not apply.
	Words []string
	// Project keeps posts whose scope.project is Project or "*".
	Project string
	// Repo ranks posts whose scope.repos is empty or holds Repo first among
	// posts with as many distinct words; it filters nothing.
	Repo string
	// Theme keeps posts of this theme, Tag posts with exactly this tag.
	Theme, Tag string
	// Statuses keeps posts with one of these statuses; nil keeps all.
	Statuses []string
	// Limit bounds the items of a word search; 0 is no bound.
	Limit int
}

// Scope is where a post applies.
type Scope struct {
	Project string   `json:"project"`
	Repos   []string `json:"repos"`
	Roles   []string `json:"roles"`
}

// Item is one post as found. Path is its file under the searched directory.
type Item struct {
	ID        string   `json:"id"`
	Theme     string   `json:"theme"`
	Type      string   `json:"type"`
	Status    string   `json:"status"`
	Scope     Scope    `json:"scope"`
	Tags      []string `json:"tags"`
	Title     string   `json:"title"`
	Situation string   `json:"situation"`
	Path      string   `json:"path"`
}

// InvalidFile is a k-*.md file that is no readable post; File is relative to
// the archive directory, with forward slashes.
type InvalidFile struct {
	File  string `json:"file"`
	Error string `json:"error"`
}

// Result is the ordered items of a search, the unreadable files, how many
// posts matched before Limit and how many posts the words hit but the
// Project, Theme, Tag or Statuses filters hid.
type Result struct {
	Items   []Item        `json:"items"`
	Invalid []InvalidFile `json:"invalid"`
	Matched int           `json:"matched"`
	Hidden  int           `json:"hidden"`
}

// Search reads every post under dir and returns those query selects, in
// ranked order. An error means dir itself cannot be read.
func Search(dir string, query Query) (Result, error) {
	posts, invalid, err := read(dir)
	if err != nil {
		return Result{}, err
	}
	words := tokens(strings.Join(query.Words, " "))
	result := Result{Items: []Item{}, Invalid: invalid}
	var hits []hit
	for _, p := range posts {
		h := p.score(words, query.Repo)
		if len(words) > 0 && h.words == 0 {
			continue
		}
		if !query.keeps(p.item) {
			result.Hidden++
			continue
		}
		hits = append(hits, h)
	}
	if len(words) == 0 {
		sort.Slice(hits, func(i, j int) bool {
			a, b := hits[i].item, hits[j].item
			if a.Theme != b.Theme {
				return a.Theme < b.Theme
			}
			return a.ID < b.ID
		})
	} else {
		sort.Slice(hits, func(i, j int) bool { return hits[i].before(hits[j]) })
	}
	result.Matched = len(hits)
	if len(words) > 0 && query.Limit > 0 && len(hits) > query.Limit {
		hits = hits[:query.Limit]
	}
	for _, h := range hits {
		result.Items = append(result.Items, h.item)
	}
	return result, nil
}

// keeps reports whether item passes the filters. A post without a project or
// status is unknown there, not excluded.
func (q Query) keeps(item Item) bool {
	if q.Project != "" && item.Scope.Project != "" && item.Scope.Project != "*" && item.Scope.Project != q.Project {
		return false
	}
	if q.Theme != "" && item.Theme != q.Theme {
		return false
	}
	if q.Tag != "" && !contains(item.Tags, q.Tag) {
		return false
	}
	return q.Statuses == nil || item.Status == "" || contains(q.Statuses, item.Status)
}

// Field weights: id, tags and title count most, then the situation, then the body.
const (
	weightBody = iota + 1
	weightSituation
	weightHead
)

type post struct {
	item                  Item
	head, situation, body []string
}

type hit struct {
	item          Item
	words, weight int
	repoFits      bool
}

// score counts the distinct words that hit p and sums, per word, the weight
// of the heaviest field it hits.
func (p post) score(words []string, repo string) hit {
	h := hit{item: p.item, repoFits: repo == "" || len(p.item.Scope.Repos) == 0 || contains(p.item.Scope.Repos, repo)}
	for _, word := range words {
		weight := 0
		switch {
		case hits(p.head, word):
			weight = weightHead
		case hits(p.situation, word):
			weight = weightSituation
		case hits(p.body, word):
			weight = weightBody
		}
		if weight > 0 {
			h.words++
			h.weight += weight
		}
	}
	return h
}

func (h hit) before(other hit) bool {
	switch {
	case h.words != other.words:
		return h.words > other.words
	case h.repoFits != other.repoFits:
		return h.repoFits
	case h.weight != other.weight:
		return h.weight > other.weight
	}
	return h.item.ID < other.item.ID
}

func hits(field []string, word string) bool {
	for _, token := range field {
		if strings.Contains(token, word) {
			return true
		}
	}
	return false
}

// tokens splits text on everything that is not a Unicode letter or digit,
// in lower case and without repeats.
func tokens(text string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, token := range strings.FieldsFunc(strings.ToLower(text), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }) {
		if !seen[token] {
			seen[token] = true
			out = append(out, token)
		}
	}
	return out
}

func contains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

// read walks dir for k-*.md files in lexical order. Directories that are
// symbolic links are not followed.
func read(dir string) ([]post, []InvalidFile, error) {
	posts := []post{}
	invalid := []InvalidFile{}
	err := filepath.WalkDir(dir, func(file string, entry fs.DirEntry, err error) error {
		if file == dir && err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir, file)
		rel = filepath.ToSlash(rel)
		if err != nil {
			invalid = append(invalid, InvalidFile{File: rel, Error: err.Error()})
			return nil
		}
		if file == dir || !strings.HasPrefix(entry.Name(), "k-") || path.Ext(entry.Name()) != ".md" {
			return nil
		}
		p, err := parse(file, rel)
		if err != nil {
			invalid = append(invalid, InvalidFile{File: rel, Error: err.Error()})
			return nil
		}
		posts = append(posts, p)
		return nil
	})
	return posts, invalid, err
}

type frontmatter struct {
	ID     string   `yaml:"id"`
	Title  string   `yaml:"title"`
	Type   string   `yaml:"type"`
	Status string   `yaml:"status"`
	Scope  Scope    `yaml:"scope"`
	Tags   []string `yaml:"tags"`
}

// parse reads file as a post. Only a regular file is opened: a symbolic link
// could lead out of the archive, and a FIFO or device could block the read.
func parse(file, rel string) (post, error) {
	info, err := os.Lstat(file)
	if err != nil {
		return post{}, err
	}
	if !info.Mode().IsRegular() {
		return post{}, errors.New("not a regular file")
	}
	raw, err := os.ReadFile(file)
	if err != nil {
		return post{}, err
	}
	lines := strings.Split(strings.ReplaceAll(string(raw), "\r\n", "\n"), "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return post{}, errors.New("no frontmatter")
	}
	end := -1
	for i := 1; i < len(lines) && end < 0; i++ {
		if lines[i] == "---" {
			end = i
		}
	}
	if end < 0 {
		return post{}, errors.New("frontmatter is not closed")
	}
	var meta frontmatter
	if err := yaml.Unmarshal([]byte(strings.Join(lines[1:end], "\n")), &meta); err != nil {
		return post{}, errors.New("unreadable frontmatter: " + err.Error())
	}
	if meta.ID == "" {
		return post{}, errors.New("missing id")
	}
	body := lines[end+1:]
	situation := ""
	for _, line := range body {
		if strings.HasPrefix(line, SituationPrefix) {
			situation = strings.TrimSpace(strings.TrimPrefix(line, SituationPrefix))
			break
		}
	}
	theme, _, _ := strings.Cut(rel, "/")
	if !strings.Contains(rel, "/") {
		theme = ""
	}
	item := Item{ID: meta.ID, Theme: theme, Type: meta.Type, Status: meta.Status,
		Scope: Scope{Project: meta.Scope.Project, Repos: nonNil(meta.Scope.Repos), Roles: nonNil(meta.Scope.Roles)},
		Tags:  nonNil(meta.Tags), Title: meta.Title, Situation: situation, Path: file}
	head := meta.ID + " " + strings.Join(meta.Tags, " ") + " " + meta.Title
	return post{item: item, head: tokens(head), situation: tokens(situation), body: tokens(strings.Join(body, "\n"))}, nil
}

func nonNil(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}
