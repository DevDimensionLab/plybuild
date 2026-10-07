package delivery

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type preferenceFields struct {
	Assignees *[]string `yaml:"assignees"`
	Reviewers *[]string `yaml:"reviewers"`
}
type preferences struct {
	SchemaVersion int                         `yaml:"schema_version"`
	Defaults      preferenceFields            `yaml:"defaults"`
	Repositories  map[string]preferenceFields `yaml:"repositories"`
}

var loginPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]{0,38}$`)

func people(list []string, reviewers bool) ([]string, error) {
	out := []string{}
	seen := map[string]bool{}
	for _, person := range list {
		parts := strings.Split(person, "/")
		if len(parts) > 2 || len(parts) == 2 && !reviewers {
			return nil, fmt.Errorf("invalid account %q for role", person)
		}
		for _, part := range parts {
			if !loginPattern.MatchString(part) {
				return nil, fmt.Errorf("invalid GitHub account %q", person)
			}
		}
		key := strings.ToLower(person)
		if seen[key] {
			return nil, fmt.Errorf("duplicate GitHub account %q", person)
		}
		seen[key] = true
		out = append(out, person)
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i]) < strings.ToLower(out[j]) })
	return out, nil
}

func validateMetadata(v *MetadataSelection) error {
	if v == nil {
		return nil
	}
	for _, p := range []struct {
		list      *[]string
		reviewers bool
	}{{v.Assignees, false}, {v.Reviewers, true}, {&v.RemoveAssignees, false}, {&v.RemoveReviewers, true}} {
		if p.list != nil {
			if _, e := people(*p.list, p.reviewers); e != nil {
				return e
			}
		}
	}
	for _, p := range []struct {
		add    *[]string
		remove []string
	}{{v.Assignees, v.RemoveAssignees}, {v.Reviewers, v.RemoveReviewers}} {
		if p.add != nil {
			for _, name := range *p.add {
				if includes(p.remove, name) {
					return fmt.Errorf("account %q both added and removed", name)
				}
			}
		}
	}
	return nil
}

func includes(list []string, name string) bool {
	for _, v := range list {
		if strings.EqualFold(v, name) {
			return true
		}
	}
	return false
}

func resolvePreferences(path, repository string, explicit *MetadataSelection, at string) (MetadataChoice, error) {
	out := MetadataChoice{Assignees: []string{}, Reviewers: []string{}, RemoveAssignees: []string{}, RemoveReviewers: []string{}, AssigneesSource: "none", ReviewersSource: "none", FrozenAtUTC: at, Revision: 1, State: "pending"}
	if e := validateMetadata(explicit); e != nil {
		return out, e
	}
	if path == "" {
		home, e := os.UserHomeDir()
		if e != nil {
			return out, e
		}
		path = filepath.Join(home, ".agents", "local", "pr-preferences.yaml")
	}
	raw, e := readBounded(path, 256<<10)
	if e != nil && !os.IsNotExist(e) {
		return out, fmt.Errorf("PR preferences: %w", e)
	}
	if e == nil {
		out.PreferencesSHA256 = hash(raw)
		var p preferences
		d := yaml.NewDecoder(bytes.NewReader(raw))
		d.KnownFields(true)
		if e = d.Decode(&p); e != nil {
			return out, fmt.Errorf("invalid PR preferences: %w", e)
		}
		if e = d.Decode(new(any)); e != io.EOF {
			return out, fmt.Errorf("PR preferences must contain exactly one document")
		}
		if p.SchemaVersion != 1 {
			return out, fmt.Errorf("PR preferences require schema_version: 1")
		}
		if e = applyPreferenceFields(&out, p.Defaults, "defaults"); e != nil {
			return out, e
		}
		seen := map[string]bool{}
		for repo, fields := range p.Repositories {
			key := strings.ToLower(repo)
			if !repositoryPattern.MatchString(repo) || seen[key] {
				return out, fmt.Errorf("invalid or duplicate case-insensitive PR preference repository %q", repo)
			}
			seen[key] = true
			// Validate all configured fields so a typo cannot be silently ignored.
			tmp := MetadataChoice{}
			if e = applyPreferenceFields(&tmp, fields, "repository:"+repo); e != nil {
				return out, e
			}
			if strings.EqualFold(repo, repository) {
				if e = applyPreferenceFields(&out, fields, "repository:"+repo); e != nil {
					return out, e
				}
			}
		}
	}
	if explicit != nil {
		if e = applyPreferenceFields(&out, preferenceFields{explicit.Assignees, explicit.Reviewers}, "explicit"); e != nil {
			return out, e
		}
		out.RemoveAssignees, _ = people(explicit.RemoveAssignees, false)
		out.RemoveReviewers, _ = people(explicit.RemoveReviewers, true)
	}
	if e = validateChoice(out); e != nil {
		return out, e
	}
	return out, nil
}

func applyPreferenceFields(out *MetadataChoice, fields preferenceFields, source string) error {
	if fields.Assignees != nil {
		v, e := people(*fields.Assignees, false)
		if e != nil {
			return fmt.Errorf("%s assignees: %w", source, e)
		}
		out.Assignees = v
		out.AssigneesSource = source
	}
	if fields.Reviewers != nil {
		v, e := people(*fields.Reviewers, true)
		if e != nil {
			return fmt.Errorf("%s reviewers: %w", source, e)
		}
		out.Reviewers = v
		out.ReviewersSource = source
	}
	return nil
}
func validateChoice(v MetadataChoice) error {
	return validateMetadata(&MetadataSelection{Assignees: &v.Assignees, Reviewers: &v.Reviewers, RemoveAssignees: v.RemoveAssignees, RemoveReviewers: v.RemoveReviewers})
}

func reviseChoice(old MetadataChoice, v MetadataSelection, at string) (MetadataChoice, error) {
	if e := validateMetadata(&v); e != nil {
		return old, e
	}
	out := old
	if e := applyPreferenceFields(&out, preferenceFields{v.Assignees, v.Reviewers}, "explicit_revision"); e != nil {
		return old, e
	}
	out.RemoveAssignees, _ = people(v.RemoveAssignees, false)
	out.RemoveReviewers, _ = people(v.RemoveReviewers, true)
	out.Revision++
	out.FrozenAtUTC = at
	out.State = "pending"
	out.Error = ""
	return out, validateChoice(out)
}
