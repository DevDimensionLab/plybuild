// Package publication stores runs that other applications publish through
// Ply: <root>/<repo-id|_project>/<category>/[<agent>/]<run>/ holding the
// producer's files unchanged and publication.json written by Ply. A run is
// copied to a dot-named temporary directory and renamed when complete, so
// readers never see half a run; a published run is never changed. Writes take
// an exclusive lock on <root>/.lock; reads take none and create nothing.
// Errors are *Error with a fixed Class and a detail without absolute paths.
package publication

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/filelock"
)

// DataPathCategory is the data directory category of a Project's publications.
const DataPathCategory = "publications"

// ProjectSegment is the directory of runs for the whole Project.
const ProjectSegment = "_project"

// Fixed JSON contract names.
const (
	Kind          = "PlyPublication@1"
	ListKind      = "PlyPublicationList@1"
	SchemaVersion = 1
)

// Error classes shared with the CLI contract.
const (
	ClassInvalidArguments = "publication_invalid_arguments"
	ClassNotFound         = "publication_not_found"
	ClassTooLarge         = "publication_too_large"
	ClassIOError          = "publication_io_error"
)

// Error is one `<class>: <detail>` failure.
type Error struct{ Class, Detail string }

func (e *Error) Error() string { return e.Class + ": " + e.Detail }

// File is one published file; Name is slash-separated and relative to the run.
type File struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

// Publication is the content of publication.json.
type Publication struct {
	Kind          string `json:"kind"`
	SchemaVersion int    `json:"schema_version"`
	Project       string `json:"project"`
	Repo          string `json:"repo,omitempty"`
	Category      string `json:"category"`
	Agent         string `json:"agent,omitempty"`
	Plan          string `json:"plan,omitempty"`
	Run           string `json:"run"`
	PublishedAt   string `json:"published_at"`
	Actor         string `json:"actor,omitempty"`
	Title         string `json:"title"`
	Entry         string `json:"entry"`
	Source        string `json:"source,omitempty"`
	Files         []File `json:"files"`
}

// View is one stored run with the absolute path of its directory.
type View struct {
	Publication
	Path string `json:"path"`
}

// Request is one publish. Repo and Plan are checked against the Project by
// the caller; the store checks their format only.
type Request struct {
	Project, Repo, Category, Agent, Plan, Dir, Entry, Title, Source, Actor string
}

// Filter selects runs for List; empty fields match everything.
type Filter struct {
	Repo, Category, Agent, Plan string
	All                         bool
}

// InvalidRun is a run directory whose publication.json cannot be read; Run
// is its path relative to the root.
type InvalidRun struct {
	Run   string `json:"run"`
	Error string `json:"error"`
}

// ListResult holds readable runs and unreadable ones; both are never nil.
type ListResult struct {
	Items   []View       `json:"items"`
	Invalid []InvalidRun `json:"invalid"`
}

// Store reads and writes the publications of one root directory.
type Store struct {
	root string
	now  func() time.Time
}

// Open returns the store for root, the absolute publications directory.
// Nothing is created until a publish.
func Open(root string) Store { return Store{root: root, now: time.Now} }

const (
	metadataName = "publication.json"
	lockName     = ".lock"
	tempSuffix   = ".tmp"
	runLayout    = "20060102T150405.000Z"
	maxBytes     = 20 << 20
	maxFiles     = 1000
	keepWithout  = 5
	maxTitle     = 200
)

var (
	namePattern   = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	repoPattern   = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)
	planPattern   = regexp.MustCompile(`^PLAN-[1-9][0-9]{0,8}$`)
	sourcePattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)
	runPattern    = regexp.MustCompile(`^[0-9]{8}T[0-9]{9}Z-[0-9a-f]{8}$`)
)

func invalid(detail string) error { return &Error{Class: ClassInvalidArguments, Detail: detail} }

func ioError(op string, err error) error {
	var pathError *fs.PathError
	if errors.As(err, &pathError) {
		err = pathError.Err // the path is absolute
	}
	return &Error{Class: ClassIOError, Detail: fmt.Sprintf("%s: %v", op, err)}
}

// runName is the run ID without the dot of the millisecond layout.
func runName(t time.Time, suffix string) string {
	return strings.Replace(t.UTC().Format(runLayout), ".", "", 1) + "-" + suffix
}

func checkName(kind, value string) error {
	if len(value) > 32 || !namePattern.MatchString(value) {
		return invalid(fmt.Sprintf("%s must match %s and have at most 32 characters", kind, namePattern))
	}
	return nil
}

// location validates the identity parts of a path and returns it relative to the root.
func location(repo, category, agent string) (string, error) {
	segment := ProjectSegment
	if repo != "" {
		if !repoPattern.MatchString(repo) {
			return "", invalid("--repo is not a repo ID")
		}
		segment = repo
	}
	if err := checkName("--category", category); err != nil {
		return "", err
	}
	if agent == "" {
		return filepath.Join(segment, category), nil
	}
	if err := checkName("--agent", agent); err != nil {
		return "", err
	}
	return filepath.Join(segment, category, agent), nil
}

func checkLine(kind, value string, max int, required bool) error {
	if value == "" {
		if required {
			return invalid(kind + " is required")
		}
		return nil
	}
	if !utf8.ValidString(value) || strings.TrimSpace(value) != value || utf8.RuneCountInString(value) > max ||
		strings.IndexFunc(value, unicode.IsControl) >= 0 {
		return invalid(fmt.Sprintf("%s must be one line of at most %d characters without surrounding space", kind, max))
	}
	return nil
}

// sourceFile is one regular file found under the producer directory.
type sourceFile struct {
	name, path string
}

// collect lists the regular files under dir, refusing anything else before any copy.
func collect(dir string) ([]sourceFile, error) {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return nil, invalid("--dir must be an existing directory")
	}
	var files []sourceFile
	var total int64
	err = filepath.WalkDir(dir, func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return ioError("read --dir", err)
		}
		if p == dir || entry.IsDir() {
			return nil
		}
		relative, _ := filepath.Rel(dir, p)
		name := filepath.ToSlash(relative)
		if !entry.Type().IsRegular() {
			return invalid(fmt.Sprintf("%s in --dir is not a regular file", name))
		}
		if name == metadataName {
			return invalid(metadataName + " in --dir is reserved for Ply")
		}
		info, err := entry.Info()
		if err != nil {
			return ioError("read --dir", err)
		}
		total += info.Size()
		if total > maxBytes {
			return &Error{Class: ClassTooLarge, Detail: fmt.Sprintf("--dir holds more than %d MiB", maxBytes>>20)}
		}
		if files = append(files, sourceFile{name: name, path: p}); len(files) > maxFiles {
			return invalid(fmt.Sprintf("--dir holds more than %d files", maxFiles))
		}
		return nil
	})
	sort.Slice(files, func(i, j int) bool { return files[i].name < files[j].name })
	return files, err
}

// Publish validates r and its files before any effect, then stores a new run
// and removes runs without plan beyond the newest five of the same path. A
// failure of that removal is returned with the stored view.
func (s Store) Publish(r Request) (View, error) {
	relative, err := location(r.Repo, r.Category, r.Agent)
	if err != nil {
		return View{}, err
	}
	if r.Plan != "" && !planPattern.MatchString(r.Plan) {
		return View{}, invalid("--plan must be PLAN-<n>")
	}
	if r.Source != "" && !sourcePattern.MatchString(r.Source) {
		return View{}, invalid("--source must match " + sourcePattern.String())
	}
	if err := checkLine("--title", r.Title, maxTitle, true); err != nil {
		return View{}, err
	}
	if err := checkLine("--actor", r.Actor, maxTitle, false); err != nil {
		return View{}, err
	}
	files, err := collect(r.Dir)
	if err != nil {
		return View{}, err
	}
	if !containsFile(files, r.Entry) {
		return View{}, invalid("--entry must name a file in --dir")
	}
	if err := os.MkdirAll(filepath.Join(s.root, relative), 0o755); err != nil {
		return View{}, ioError("create publication directory", err)
	}
	lock, err := filelock.Acquire(filepath.Join(s.root, lockName), filelock.Options{})
	if err != nil {
		return View{}, ioError("lock publications", err)
	}
	v, err := s.store(r, relative, files)
	if releaseErr := lock.Release(); releaseErr != nil && err == nil {
		err = &Error{Class: ClassIOError, Detail: fmt.Sprintf("run %s was stored; releasing the lock failed", v.Run)}
	}
	return v, err
}

func containsFile(files []sourceFile, name string) bool {
	for _, f := range files {
		if f.name == name {
			return true
		}
	}
	return false
}

func (s Store) store(r Request, relative string, files []sourceFile) (View, error) {
	directory := filepath.Join(s.root, relative)
	runs, err := runNames(directory)
	if err != nil {
		return View{}, err
	}
	at := s.now().UTC().Truncate(time.Millisecond)
	if n := len(runs); n > 0 {
		// A run must sort after the newest one, also with a clock behind it.
		if latest, err := time.Parse(runLayout, runs[n-1][:15]+"."+runs[n-1][15:19]); err == nil && !at.After(latest) {
			at = latest.Add(time.Millisecond)
		}
	}
	random := make([]byte, 4)
	if _, err := rand.Read(random); err != nil {
		return View{}, ioError("create run ID", err)
	}
	run := runName(at, hex.EncodeToString(random))
	temp := filepath.Join(directory, "."+run+tempSuffix)
	p := Publication{Kind: Kind, SchemaVersion: SchemaVersion, Project: r.Project, Repo: r.Repo, Category: r.Category, Agent: r.Agent,
		Plan: r.Plan, Run: run, PublishedAt: at.Format(time.RFC3339), Actor: r.Actor, Title: r.Title, Entry: r.Entry, Source: r.Source}
	if p.Files, err = copyRun(temp, files); err == nil {
		err = writeMetadata(filepath.Join(temp, metadataName), p)
	}
	if err == nil {
		err = os.Rename(temp, filepath.Join(directory, run))
	}
	if err != nil {
		os.RemoveAll(temp)
		return View{}, err
	}
	v := View{Publication: p, Path: filepath.Join(directory, run)}
	if err := s.prune(directory); err != nil {
		return v, &Error{Class: ClassIOError, Detail: fmt.Sprintf("run %s was stored; removing older runs failed: %s", run, err)}
	}
	return v, nil
}

func copyRun(temp string, files []sourceFile) ([]File, error) {
	if err := os.Mkdir(temp, 0o755); err != nil {
		return nil, ioError("create run", err)
	}
	result := make([]File, 0, len(files))
	var total int64
	for _, f := range files {
		target := filepath.Join(temp, filepath.FromSlash(f.name))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return nil, ioError("create run", err)
		}
		size, sum, err := copyFile(f.path, target)
		if err != nil {
			return nil, err
		}
		if total += size; total > maxBytes {
			return nil, &Error{Class: ClassTooLarge, Detail: fmt.Sprintf("--dir grew past %d MiB during publish", maxBytes>>20)}
		}
		result = append(result, File{Name: f.name, Size: size, SHA256: sum})
	}
	return result, nil
}

func copyFile(source, target string) (int64, string, error) {
	in, err := os.Open(source)
	if err != nil {
		return 0, "", ioError("read --dir", err)
	}
	defer in.Close()
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return 0, "", ioError("write run", err)
	}
	hash := sha256.New()
	size, err := io.Copy(io.MultiWriter(out, hash), io.LimitReader(in, maxBytes+1))
	if closeErr := out.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return 0, "", ioError("write run", err)
	}
	return size, hex.EncodeToString(hash.Sum(nil)), nil
}

func writeMetadata(target string, p Publication) error {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(p); err != nil {
		return ioError("encode "+metadataName, err)
	}
	if err := os.WriteFile(target, buffer.Bytes(), 0o644); err != nil {
		return ioError("write "+metadataName, err)
	}
	return nil
}

// prune removes runs without plan beyond the newest keepWithout in directory.
func (s Store) prune(directory string) error {
	runs, err := runNames(directory)
	if err != nil {
		return err
	}
	kept := 0
	for i := len(runs) - 1; i >= 0; i-- {
		p, err := readMetadata(filepath.Join(directory, runs[i]))
		if err != nil || p.Plan != "" {
			continue // an unreadable run is reported by list, never removed
		}
		if kept++; kept > keepWithout {
			if err := os.RemoveAll(filepath.Join(directory, runs[i])); err != nil {
				return ioError("remove run "+runs[i], err)
			}
		}
	}
	return nil
}

// runNames returns the run directories of directory in ascending order.
func runNames(directory string) ([]string, error) {
	entries, err := os.ReadDir(directory)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, ioError("read publications", err)
	}
	var runs []string
	for _, entry := range entries {
		if entry.IsDir() && runPattern.MatchString(entry.Name()) {
			runs = append(runs, entry.Name())
		}
	}
	return runs, nil // ReadDir sorts by name
}

func readMetadata(run string) (Publication, error) {
	contents, err := os.ReadFile(filepath.Join(run, metadataName))
	if err != nil {
		return Publication{}, ioError("read "+metadataName, err)
	}
	var p Publication
	if err := json.Unmarshal(contents, &p); err != nil || p.Kind != Kind || p.Run != filepath.Base(run) {
		return Publication{}, &Error{Class: ClassIOError, Detail: metadataName + " is not a valid " + Kind}
	}
	return p, nil
}

// Show returns the run of the path; the newest one when run is empty.
func (s Store) Show(repo, category, agent, run string) (View, error) {
	relative, err := location(repo, category, agent)
	if err != nil {
		return View{}, err
	}
	directory := filepath.Join(s.root, relative)
	if run == "" {
		runs, err := runNames(directory)
		if err != nil {
			return View{}, err
		}
		if len(runs) == 0 {
			return View{}, &Error{Class: ClassNotFound, Detail: "no run in " + filepath.ToSlash(relative)}
		}
		run = runs[len(runs)-1]
	} else if !runPattern.MatchString(run) {
		return View{}, invalid("--run is not a run ID")
	}
	if _, err := os.Stat(filepath.Join(directory, run)); errors.Is(err, fs.ErrNotExist) {
		return View{}, &Error{Class: ClassNotFound, Detail: path.Join(filepath.ToSlash(relative), run)}
	}
	p, err := readMetadata(filepath.Join(directory, run))
	if err != nil {
		return View{}, err
	}
	return View{Publication: p, Path: filepath.Join(directory, run)}, nil
}

// List returns the runs matching f, newest first; without f.All only the
// newest matching run of each path.
func (s Store) List(f Filter) (ListResult, error) {
	result := ListResult{Items: []View{}, Invalid: []InvalidRun{}}
	directories, err := s.paths(f)
	if err != nil {
		return result, err
	}
	for _, directory := range directories {
		runs, err := runNames(directory)
		if err != nil {
			return result, err
		}
		for i := len(runs) - 1; i >= 0; i-- {
			run := filepath.Join(directory, runs[i])
			p, err := readMetadata(run)
			if err != nil {
				relative, _ := filepath.Rel(s.root, run)
				result.Invalid = append(result.Invalid, InvalidRun{Run: filepath.ToSlash(relative), Error: err.Error()})
				continue
			}
			if f.Plan != "" && p.Plan != f.Plan {
				continue
			}
			result.Items = append(result.Items, View{Publication: p, Path: run})
			if !f.All {
				break
			}
		}
	}
	sort.SliceStable(result.Items, func(i, j int) bool { return result.Items[i].Run > result.Items[j].Run })
	return result, nil
}

// paths returns the directories that can hold runs matching f.
func (s Store) paths(f Filter) ([]string, error) {
	var found []string
	segments, err := subdirectories(s.root)
	if err != nil {
		return nil, err
	}
	for _, segment := range segments {
		if f.Repo != "" && segment != f.Repo {
			continue
		}
		categories, err := subdirectories(filepath.Join(s.root, segment))
		if err != nil {
			return nil, err
		}
		for _, category := range categories {
			if f.Category != "" && category != f.Category {
				continue
			}
			directory := filepath.Join(s.root, segment, category)
			if f.Agent == "" {
				found = append(found, directory)
			}
			agents, err := subdirectories(directory)
			if err != nil {
				return nil, err
			}
			for _, agent := range agents {
				if namePattern.MatchString(agent) && (f.Agent == "" || agent == f.Agent) {
					found = append(found, filepath.Join(directory, agent))
				}
			}
		}
	}
	return found, nil
}

// subdirectories lists the directories of directory that are not dot-named.
func subdirectories(directory string) ([]string, error) {
	entries, err := os.ReadDir(directory)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, ioError("read publications", err)
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}
