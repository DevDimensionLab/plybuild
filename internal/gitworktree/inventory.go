package gitworktree

import (
	"bytes"
	"fmt"
	"strings"
)

type Entry struct {
	Locator, Ref, OID string
	Locked, Prunable  bool
	Bare, Detached    bool
}

func ParseInventory(raw []byte) ([]Entry, error) {
	entries := []Entry{}
	var current *Entry
	seen := map[string]bool{}
	finish := func() error {
		if current == nil {
			return nil
		}
		if current.Locator == "" {
			return fmt.Errorf("incomplete worktree record: missing locator")
		}
		entries = append(entries, *current)
		current = nil
		return nil
	}
	for _, field := range bytes.Split(raw, []byte{0}) {
		text := string(field)
		if text == "" {
			if err := finish(); err != nil {
				return nil, err
			}
			continue
		}
		key, value, _ := strings.Cut(text, " ")
		if key == "worktree" {
			if err := finish(); err != nil {
				return nil, err
			}
			current = &Entry{Locator: value}
			seen = map[string]bool{"worktree": true}
			continue
		}
		if current == nil {
			return nil, fmt.Errorf("%s outside worktree record", key)
		}
		if seen[key] {
			return nil, fmt.Errorf("duplicate %s in worktree record", key)
		}
		seen[key] = true
		switch key {
		case "HEAD":
			current.OID = value
		case "branch":
			current.Ref = value
		case "locked":
			current.Locked = true
		case "prunable":
			current.Prunable = true
		case "bare":
			current.Bare = true
		case "detached":
			current.Detached = true
		default:
			return nil, fmt.Errorf("unknown worktree record field %s", key)
		}
	}
	if err := finish(); err != nil {
		return nil, err
	}
	return entries, nil
}
