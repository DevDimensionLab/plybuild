package backlog

import (
	"errors"
	"io/fs"
	"os"
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	maxTags       = 32
	maxLabelRunes = 100
)

var repoID = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,99}$`)

// SetTags replaces the item's repo ids and topics; nil or empty clears a list.
func (s Store) SetTags(id string, repos, topics []string) (Item, error) {
	if err := s.checkDirectory(); err != nil {
		return Item{}, err
	}
	if !itemID.MatchString(id) {
		return Item{}, invalidArguments("invalid item id")
	}
	repos, topics, err := checkTags(repos, topics)
	if err != nil {
		return Item{}, err
	}
	if _, err := os.Stat(s.dir); errors.Is(err, fs.ErrNotExist) {
		return Item{}, notFound(id)
	}
	unlock, err := s.lock()
	if err != nil {
		return Item{}, err
	}
	defer unlock()
	item, err := s.read(id)
	if err != nil || (slices.Equal(item.Repos, repos) && slices.Equal(item.Topics, topics)) {
		return item, err
	}
	item.Repos, item.Topics = repos, topics
	return item, s.write(item)
}

// AddNoteTagged is AddNoteAs that also sets the new note's repo IDs and topics.
func (s Store) AddNoteTagged(id, text, author string, repos, topics []string) (Item, error) {
	if err := s.checkDirectory(); err != nil {
		return Item{}, err
	}
	if !itemID.MatchString(id) {
		return Item{}, invalidArguments("invalid item id")
	}
	if err := checkText(text); err != nil {
		return Item{}, err
	}
	author = strings.TrimSpace(author)
	if author != "" && !validLabel(author) {
		return Item{}, invalidArguments("author must be one line of at most 100 characters")
	}
	repos, topics, err := checkTags(repos, topics)
	if err != nil {
		return Item{}, err
	}
	if _, err := os.Stat(s.dir); errors.Is(err, fs.ErrNotExist) {
		return Item{}, notFound(id)
	}
	unlock, err := s.lock()
	if err != nil {
		return Item{}, err
	}
	defer unlock()
	item, err := s.read(id)
	if err != nil {
		return Item{}, err
	}
	item.Notes = append(item.Notes, Note{At: s.now().UTC().Format(timestampStyle), Text: text, Author: author, Repos: repos, Topics: topics})
	return item, s.write(item)
}

// SetNoteTags replaces the repo IDs and topics of the note at index (0-based, stored order);
// nil or empty clears a list. When nothing changes, no item file is written.
func (s Store) SetNoteTags(id string, index int, repos, topics []string) (Item, error) {
	if err := s.checkDirectory(); err != nil {
		return Item{}, err
	}
	if !itemID.MatchString(id) {
		return Item{}, invalidArguments("invalid item id")
	}
	if index < 0 {
		return Item{}, errNoSuchNote
	}
	repos, topics, err := checkTags(repos, topics)
	if err != nil {
		return Item{}, err
	}
	if _, err := os.Stat(s.dir); errors.Is(err, fs.ErrNotExist) {
		return Item{}, notFound(id)
	}
	unlock, err := s.lock()
	if err != nil {
		return Item{}, err
	}
	defer unlock()
	item, err := s.read(id)
	if err != nil {
		return Item{}, err
	}
	if index >= len(item.Notes) {
		return Item{}, errNoSuchNote
	}
	note := &item.Notes[index]
	if slices.Equal(note.Repos, repos) && slices.Equal(note.Topics, topics) {
		return item, nil
	}
	note.Repos, note.Topics = repos, topics
	return item, s.write(item)
}

var errNoSuchNote = &Error{Class: ClassInvalidArguments, Detail: "no such note"}

// checkTags normalizes repo ids and topics for every tag writer.
func checkTags(repos, topics []string) ([]string, []string, error) {
	repos, ok := normalizeTags(repos, validRepo)
	if !ok {
		return nil, nil, invalidArguments("repos must be at most 32 lowercase ASCII slugs of at most 100 characters")
	}
	topics, ok = normalizeTags(topics, validLabel)
	if !ok {
		return nil, nil, invalidArguments("topics must be at most 32 non-blank lines of at most 100 characters")
	}
	return repos, topics, nil
}

func normalizeTags(values []string, valid func(string) bool) ([]string, bool) {
	var result []string
	for _, value := range values {
		value = strings.TrimSpace(value)
		if !valid(value) {
			return nil, false
		}
		if !slices.ContainsFunc(result, func(kept string) bool { return strings.EqualFold(kept, value) }) {
			result = append(result, value)
		}
	}
	return result, len(result) <= maxTags
}

func validTags(values []string, valid func(string) bool) bool {
	return len(values) <= maxTags && !slices.ContainsFunc(values, func(value string) bool { return !valid(value) })
}

func validRepo(value string) bool { return repoID.MatchString(value) }

func validLabel(value string) bool {
	return value != "" && value == strings.TrimSpace(value) && utf8.ValidString(value) &&
		utf8.RuneCountInString(value) <= maxLabelRunes && strings.IndexFunc(value, unicode.IsControl) < 0
}
