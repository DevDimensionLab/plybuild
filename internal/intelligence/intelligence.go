// Package intelligence stores what agents learn in a Project: immutable
// posts, one canonical JSON file each, whose schema field is the kind. Ply
// owns storage, identity, origin, validation and reads; the owners of a kind
// own its content. Errors are *Error with a fixed Class and a detail without
// paths.
package intelligence

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Kinds and the schema that identifies each in a post.
const (
	KindObservation      = "observation"
	KindCodeHealthReview = "code-health-review"
)

var schemas = map[string]string{
	KindObservation:      "ply.observation/v1",
	KindCodeHealthReview: "ply.code-health-review/v1",
}

// Kinds lists the kinds in a stable order.
var Kinds = []string{KindObservation, KindCodeHealthReview}

var limits = map[string]int{KindObservation: 32 << 10, KindCodeHealthReview: 4 << 20}

// Error classes shared with the CLI contract.
const (
	ClassInvalidArguments = "intelligence_invalid_arguments"
	ClassInvalidPost      = "intelligence_invalid_post"
	ClassConflict         = "intelligence_conflict"
	ClassNotFound         = "intelligence_not_found"
	ClassIOError          = "intelligence_io_error"
)

// Error is one `<class>: <detail>` failure.
type Error struct{ Class, Detail string }

func (e *Error) Error() string { return e.Class + ": " + e.Detail }

// Receipt states of Observe.
const (
	StateRecorded        = "recorded"
	StateAlreadyRecorded = "already_recorded"
)

// Origin is who recorded a post, in which repo and logical case.
type Origin struct {
	Activity string  `json:"activity"`
	Project  string  `json:"project"`
	Provider string  `json:"provider"`
	Repo     string  `json:"repo"`
	Role     string  `json:"role"`
	Session  *string `json:"session"`
}

// Receipt reports one Observe.
type Receipt struct {
	State  string `json:"state"`
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	SHA256 string `json:"sha256"`
}

// Item is one stored post with its kind and the hash of its bytes.
type Item struct {
	ID         string          `json:"id"`
	Kind       string          `json:"kind"`
	RecordedAt string          `json:"recorded_at"`
	SHA256     string          `json:"sha256"`
	Post       json.RawMessage `json:"post"`
}

// InvalidFile is a stored file that cannot be read; File is only its name.
type InvalidFile struct {
	File  string `json:"file"`
	Error string `json:"error"`
}

// ListResult holds matching posts in recorded order and unreadable files.
type ListResult struct {
	Items   []Item        `json:"items"`
	Invalid []InvalidFile `json:"invalid"`
}

// Filter selects posts; empty fields match everything. After is the ID of a
// stored post; only posts ordered after it match.
type Filter struct {
	Kind, Repo, Role, Activity, Since, After, LinksTo string
}

// Store holds the posts of one Project.
type Store struct {
	dir     string
	project string
	repos   map[string]bool
}

// Open returns the store in directory, the Project's intelligence directory,
// for project and its repository IDs. Nothing is created until Observe.
func Open(directory, project string, repos []string) Store {
	known := map[string]bool{}
	for _, repo := range repos {
		known[repo] = true
	}
	return Store{dir: directory, project: project, repos: known}
}

const timeLayout = "2006-01-02T15:04:05Z"

var (
	idPattern        = regexp.MustCompile(`^obs-[0-9a-f]{8}(-[0-9a-f]{4}){3}-[0-9a-f]{12}$`)
	knowledgePattern = regexp.MustCompile(`^k-[0-9]{8}-[a-z0-9]+(-[a-z0-9]+)*$`)
	repoPattern      = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	shaPattern       = regexp.MustCompile(`^[0-9a-f]{40}$`)
)

// Draft returns a post of kind with a new ID, recordedAt and origin, and
// empty content fields to fill in. The ID and time stay on every retry.
func Draft(kind string, origin Origin, recordedAt time.Time) ([]byte, error) {
	schema, known := schemas[kind]
	if !known {
		return nil, invalidArguments("kind must be one of " + strings.Join(Kinds, ", "))
	}
	post := map[string]any{"schema": schema, "id": newID(),
		"recorded_at": recordedAt.UTC().Format(timeLayout), "origin": origin, "links": []any{}}
	switch kind {
	case KindObservation:
		for key, value := range map[string]any{"situation": "", "observation": "", "interpretation": nil, "limitations": "",
			"evidence": []any{map[string]any{"basis": "direct", "ref": "", "detail": ""}}, "prior_refs": []any{},
			"effect": map[string]any{"kind": "unknown", "detail": ""}} {
			post[key] = value
		}
	case KindCodeHealthReview:
		for key, value := range map[string]any{"base": "", "head": "", "tool": "", "command": "", "exit": 0,
			"reviewer":              map[string]any{"provider": origin.Provider, "model": "unknown", "task": nil},
			"first_review_question": "", "ratings": []any{}, "signals": 0, "used": 0, "irrelevant": 0, "wrong": 0,
			"other_observations": nil, "product_suggestion": nil, "follow_up_commit": nil, "backlog": []any{},
			"counting": "code-health-review/v2", "diff": "", "diff_stderr": nil} {
			post[key] = value
		}
	}
	return canonical(post)
}

// Observe validates raw and stores its canonical form as a new post. The same
// ID with the same content is already_recorded; with other content a conflict.
func (s Store) Observe(raw []byte) (Receipt, error) {
	if err := s.checkDirectory(); err != nil {
		return Receipt{}, err
	}
	post, kind, err := s.validate(raw)
	if err != nil {
		return Receipt{}, err
	}
	data, err := canonical(post.value)
	if err != nil {
		return Receipt{}, invalidPost("cannot encode the post")
	}
	if len(data) > limits[kind] {
		return Receipt{}, invalidPost(fmt.Sprintf("a %s post is at most %d bytes; link larger reports", kind, limits[kind]))
	}
	receipt := Receipt{ID: post.id, Kind: kind, SHA256: digest(data)}
	if existing, err := os.ReadFile(s.path(post.id)); err == nil {
		return s.replay(receipt, existing, data)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return Receipt{}, ioError("read post", err)
	}
	for _, link := range post.links {
		if idPattern.MatchString(link) {
			if _, err := os.Lstat(s.path(link)); err != nil {
				return Receipt{}, invalidPost("link target " + link + " is not stored")
			}
		}
	}
	if err := s.publish(post.id, data); err != nil {
		if errors.Is(err, fs.ErrExist) {
			existing, readErr := os.ReadFile(s.path(post.id))
			if readErr != nil {
				return Receipt{}, ioError("read post", readErr)
			}
			return s.replay(receipt, existing, data)
		}
		return Receipt{}, ioError("store post", err)
	}
	receipt.State = StateRecorded
	return receipt, nil
}

func (s Store) replay(receipt Receipt, existing, data []byte) (Receipt, error) {
	if !bytes.Equal(existing, data) {
		return Receipt{}, &Error{Class: ClassConflict, Detail: "post " + receipt.ID + " is stored with other content; read it before correcting"}
	}
	receipt.State = StateAlreadyRecorded
	return receipt, nil
}

// publish creates the post file without ever replacing one.
func (s Store) publish(id string, data []byte) error {
	folder := filepath.Join(s.dir, "observations")
	if err := os.MkdirAll(folder, 0o755); err != nil {
		return err
	}
	staged, err := os.CreateTemp(folder, ".staged-*")
	if err != nil {
		return err
	}
	defer os.Remove(staged.Name())
	if _, err = staged.Write(data); err == nil {
		err = staged.Sync()
	}
	if closeErr := staged.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Link(staged.Name(), s.path(id))
}

// Show returns the stored post id.
func (s Store) Show(id string) (Item, error) {
	if !idPattern.MatchString(id) {
		return Item{}, invalidArguments("id must be obs-<uuid>")
	}
	item, detail, err := s.read(id + ".json")
	if errors.Is(err, fs.ErrNotExist) {
		return Item{}, &Error{Class: ClassNotFound, Detail: "post " + id + " is not stored"}
	}
	if err != nil {
		return Item{}, ioError("read post", err)
	}
	if detail != "" {
		return Item{}, invalidPost("stored post " + id + ": " + detail)
	}
	return item, nil
}

// List returns the posts matching f, ordered by recorded_at and then ID.
func (s Store) List(f Filter) (ListResult, error) {
	result := ListResult{Items: []Item{}, Invalid: []InvalidFile{}}
	if f.Kind != "" && schemas[f.Kind] == "" {
		return result, invalidArguments("kind must be one of " + strings.Join(Kinds, ", "))
	}
	if f.Since != "" {
		if _, err := time.Parse(timeLayout, f.Since); err != nil {
			return result, invalidArguments("since must be a UTC time such as 2026-11-16T10:00:00Z")
		}
	}
	entries, err := os.ReadDir(filepath.Join(s.dir, "observations"))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return result, ioError("read intelligence directory", err)
	}
	var all []Item
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		item, detail, err := s.read(name)
		if err != nil {
			detail = "cannot read file"
		}
		if detail != "" {
			result.Invalid = append(result.Invalid, InvalidFile{File: name, Error: detail})
			continue
		}
		all = append(all, item)
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].RecordedAt != all[j].RecordedAt {
			return all[i].RecordedAt < all[j].RecordedAt
		}
		return all[i].ID < all[j].ID
	})
	start := 0
	if f.After != "" {
		start = -1
		for i, item := range all {
			if item.ID == f.After {
				start = i + 1
			}
		}
		if start < 0 {
			return result, &Error{Class: ClassNotFound, Detail: "post " + f.After + " is not stored"}
		}
	}
	for _, item := range all[start:] {
		if f.matches(item) {
			result.Items = append(result.Items, item)
		}
	}
	return result, nil
}

func (f Filter) matches(item Item) bool {
	var head struct {
		Origin Origin `json:"origin"`
		Links  []struct {
			ID string `json:"id"`
		} `json:"links"`
	}
	_ = json.Unmarshal(item.Post, &head)
	linked := f.LinksTo == ""
	for _, link := range head.Links {
		linked = linked || link.ID == f.LinksTo
	}
	return (f.Kind == "" || item.Kind == f.Kind) && (f.Repo == "" || head.Origin.Repo == f.Repo) &&
		(f.Role == "" || head.Origin.Role == f.Role) && (f.Activity == "" || head.Origin.Activity == f.Activity) &&
		(f.Since == "" || item.RecordedAt >= f.Since) && linked
}

// read returns a stored file as an item, or a detail when it is not a valid post.
func (s Store) read(name string) (Item, string, error) {
	id := strings.TrimSuffix(name, ".json")
	if !strings.HasSuffix(name, ".json") || !idPattern.MatchString(id) {
		return Item{}, "invalid file name", nil
	}
	path := filepath.Join(s.dir, "observations", name)
	info, err := os.Lstat(path)
	if err != nil {
		return Item{}, "", err
	}
	if !info.Mode().IsRegular() {
		return Item{}, "not a regular file", nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Item{}, "", err
	}
	var head struct {
		Schema     string `json:"schema"`
		ID         string `json:"id"`
		RecordedAt string `json:"recorded_at"`
	}
	if json.Unmarshal(data, &head) != nil || head.ID != id {
		return Item{}, "invalid post", nil
	}
	kind := ""
	for k, schema := range schemas {
		if schema == head.Schema {
			kind = k
		}
	}
	if kind == "" {
		return Item{}, "unknown schema", nil
	}
	return Item{ID: id, Kind: kind, RecordedAt: head.RecordedAt, SHA256: digest(data), Post: data}, "", nil
}

func (s Store) path(id string) string {
	return filepath.Join(s.dir, "observations", id+".json")
}

func (s Store) checkDirectory() error {
	if !filepath.IsAbs(s.dir) {
		return invalidArguments("intelligence directory must be absolute")
	}
	return nil
}

// canonical is the stored form: sorted keys, two-space indent, raw UTF-8
// with <, > and & unescaped, and a final newline.
func canonical(value any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// newID is obs- and a random version 4 UUID.
func newID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	h := hex.EncodeToString(b[:])
	return "obs-" + h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:32]
}

func digest(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func invalidArguments(detail string) error {
	return &Error{Class: ClassInvalidArguments, Detail: detail}
}

func invalidPost(detail string) error {
	return &Error{Class: ClassInvalidPost, Detail: detail}
}

// ioError keeps only the error kind, never the path inside err.
func ioError(action string, err error) error {
	kind := "failed"
	switch {
	case errors.Is(err, fs.ErrPermission):
		kind = "permission denied"
	case errors.Is(err, fs.ErrNotExist):
		kind = "not found"
	}
	return &Error{Class: ClassIOError, Detail: action + ": " + kind}
}
