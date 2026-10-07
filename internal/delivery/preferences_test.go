package delivery

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func physicalTemp(t *testing.T) string {
	t.Helper()
	p, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	return p
}

func stringList(v ...string) *[]string { return &v }

func TestPreferenceResolutionPreservesRolesAndEmptyOverrides(t *testing.T) {
	path := filepath.Join(physicalTemp(t), "preferences.yaml")
	if e := os.WriteFile(path, []byte("schema_version: 1\ndefaults:\n  assignees: [global-owner]\n  reviewers: [global-reviewer]\nrepositories:\n  Example/Project:\n    assignees: [repo-owner]\n    reviewers: []\n"), 0600); e != nil {
		t.Fatal(e)
	}
	got, e := resolvePreferences(path, "example/PROJECT", nil, "2026-01-01T00:00:00Z")
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(got.Assignees, []string{"repo-owner"}) || len(got.Reviewers) != 0 || got.ReviewersSource != "repository:Example/Project" {
		t.Fatalf("inheritance: %+v", got)
	}
	got, e = resolvePreferences(path, "example/project", &MetadataSelection{Assignees: stringList(), Reviewers: stringList("explicit-reviewer")}, "now")
	if e != nil {
		t.Fatal(e)
	}
	if len(got.Assignees) != 0 || !reflect.DeepEqual(got.Reviewers, []string{"explicit-reviewer"}) || got.AssigneesSource != "explicit" {
		t.Fatalf("explicit: %+v", got)
	}
	missing, e := resolvePreferences(filepath.Join(physicalTemp(t), "absent.yaml"), "example/project", nil, "now")
	if e != nil || len(missing.Assignees)+len(missing.Reviewers) != 0 || missing.AssigneesSource != "none" {
		t.Fatalf("missing invented people: %+v, %v", missing, e)
	}
}

func TestInvalidPreferencesDoNotFallBack(t *testing.T) {
	for name, raw := range map[string]string{
		"unknown":        "schema_version: 1\ndefaults:\n  reviewers: [right]\n  reviewer: [wrong]\n",
		"version":        "schema_version: 2\n",
		"bad-account":    "schema_version: 1\ndefaults:\n  assignees: ['not a login']\n",
		"assignee-team":  "schema_version: 1\ndefaults:\n  assignees: [owner/team]\n",
		"case-duplicate": "schema_version: 1\nrepositories:\n  Example/Project: {}\n  example/project: {}\n",
		"multiple":       "schema_version: 1\n---\nschema_version: 1\n",
		"duplicate-role": "schema_version: 1\ndefaults:\n  reviewers: [same, SAME]\n",
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(physicalTemp(t), "prefs.yaml")
			if e := os.WriteFile(path, []byte(raw), 0600); e != nil {
				t.Fatal(e)
			}
			if _, e := resolvePreferences(path, "example/project", nil, "now"); e == nil {
				t.Fatalf("accepted invalid preferences: %s", raw)
			}
		})
	}
}

func TestFrozenChoicesSurviveFileChangesAndExplicitRevisionIsSeparate(t *testing.T) {
	path := filepath.Join(physicalTemp(t), "prefs.yaml")
	if e := os.WriteFile(path, []byte("schema_version: 1\ndefaults:\n  reviewers: [first-reviewer]\n"), 0600); e != nil {
		t.Fatal(e)
	}
	first, e := resolvePreferences(path, "example/project", nil, "first")
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(path, []byte("schema_version: 1\ndefaults:\n  reviewers: [second-reviewer]\n"), 0600); e != nil {
		t.Fatal(e)
	}
	change, e := reviseChoice(first, MetadataSelection{Assignees: stringList("explicit-owner")}, "later")
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(change.Reviewers, []string{"first-reviewer"}) || change.ReviewersSource != "defaults" || change.Revision != 2 || change.PreferencesSHA256 != first.PreferencesSHA256 {
		t.Fatalf("revision reread defaults: %+v", change)
	}
	if first.Revision != 1 || !reflect.DeepEqual(first.Reviewers, []string{"first-reviewer"}) {
		t.Fatal("original choice rewritten")
	}
	_, e = reviseChoice(first, MetadataSelection{Reviewers: stringList("same"), RemoveReviewers: []string{"SAME"}}, "later")
	if e == nil || !strings.Contains(e.Error(), "both added and removed") {
		t.Fatalf("conflicting role mutation: %v", e)
	}
}
