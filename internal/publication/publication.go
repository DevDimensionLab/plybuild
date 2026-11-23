// Package publication stores runs that other applications publish through
// Ply: <root>/<repo-id|_project>/<category>/[<agent>/]<run>/ holding the
// producer's files unchanged and publication.json written by Ply. A run of a
// module of a repo lives under <repo-id>/_modules/<escaped module path>/ with
// the same layout. Readers of the default paths skip directories that start
// with "_", so module runs are listed only on request. A run is
// copied to a dot-named temporary directory and renamed when complete, so
// readers never see half a run; a published run is never changed. A
// withdrawal writes withdrawn.json in the path directory and hides the runs up
// to and including the newest one at that time; nothing is removed. Writes take
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
	"net/url"
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

// ProjectSegment is the directory of runs for the whole Project.
const ProjectSegment = "_project"

// ModulesSegment is the directory of a repo that holds its module paths.
const ModulesSegment = "_modules"

// RootModule is the module of a repo's own runs.
const RootModule = "."

// Fixed JSON contract names.
const (
	Kind           = "PlyPublication@1"
	ListKind       = "PlyPublicationList@1"
	WithdrawalKind = "PlyPublicationWithdrawal@1"
	SchemaVersion  = 1
)

// Error classes shared with the CLI contract.
const (
	ClassInvalidArguments = "publication_invalid_arguments"
	ClassNotFound         = "publication_not_found"
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
	Module        string `json:"module,omitempty"`
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

// Withdrawal is the content of withdrawn.json: the runs of the path up to and
// including Through are hidden from List and from Show without a run.
type Withdrawal struct {
	Kind          string `json:"kind"`
	SchemaVersion int    `json:"schema_version"`
	Project       string `json:"project"`
	Repo          string `json:"repo,omitempty"`
	Module        string `json:"module,omitempty"`
	Category      string `json:"category"`
	Agent         string `json:"agent,omitempty"`
	Through       string `json:"through"`
	WithdrawnAt   string `json:"withdrawn_at"`
	Actor         string `json:"actor,omitempty"`
}

// View is one stored run with its module (RootModule for the repo's own) and
// the absolute path of its directory.
type View struct {
	Publication
	Module string `json:"module"`
	Path   string `json:"path"`
}

func newView(p Publication, path string) View {
	return View{Publication: p, Module: DisplayModule(p.Module), Path: path}
}

// DisplayModule is the module as readbacks show it: RootModule for none.
func DisplayModule(module string) string {
	if module == "" {
		return RootModule
	}
	return module
}

// Request is one publish. Repo and Plan are checked against the Project by
// the caller; the store checks their format only.
type Request struct {
	Project, Repo, Module, Category, Agent, Plan, Dir, Entry, Title, Source, Actor string
}

// Filter selects runs for List; empty fields match everything. Without
// Module or Modules only the runs of repos and the Project are listed; Module
// selects one module of Repo, and Modules adds every module.
type Filter struct {
	Repo, Module, Category, Agent, Plan string
	Modules, All                        bool
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
	withdrawName = "withdrawn.json"
	lockName     = ".lock"
	tempSuffix   = ".tmp"
	runLayout    = "20060102T150405.000Z"
	maxFiles     = 1000
	keepWithout  = 5
	maxTitle     = 200
	maxModule    = 255
)

// invalidWithdrawal is the detail of a withdrawn.json that cannot be parsed.
const invalidWithdrawal = withdrawName + " is not a valid " + WithdrawalKind

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

// CheckModule validates the form of a module path: relative, slash-separated
// and clean, without segments that start with a dot, and at most 255 bytes
// when escaped. RootModule is valid. Whether the module exists is not checked.
func CheckModule(module string) error {
	if module == RootModule {
		return nil
	}
	if module == "" || !utf8.ValidString(module) || strings.ContainsRune(module, '\\') || strings.IndexFunc(module, unicode.IsControl) >= 0 ||
		path.Clean(module) != module || strings.HasPrefix(module, "/") {
		return invalid("--module must be a clean relative path with /")
	}
	for _, segment := range strings.Split(module, "/") {
		if strings.HasPrefix(segment, ".") {
			return invalid("--module must not have a segment that starts with a dot")
		}
	}
	if len(url.PathEscape(module)) > maxModule {
		return invalid(fmt.Sprintf("--module must be at most %d bytes when escaped", maxModule))
	}
	return nil
}

// checkModuleOf checks a module selection; empty means none.
func checkModuleOf(repo, module string) error {
	if module == "" {
		return nil
	}
	if repo == "" {
		return invalid("--module requires --repo")
	}
	return CheckModule(module)
}

// storedModule is the module as stored: empty for the repo's own runs.
func storedModule(module string) string {
	if module == RootModule {
		return ""
	}
	return module
}

// location validates the identity parts of a path and returns it relative to the root.
func location(repo, module, category, agent string) (string, error) {
	segment := ProjectSegment
	if repo != "" {
		if !repoPattern.MatchString(repo) {
			return "", invalid("--repo is not a repo ID")
		}
		segment = repo
	}
	if err := checkModuleOf(repo, module); err != nil {
		return "", err
	}
	if module != "" && module != RootModule {
		segment = filepath.Join(segment, ModulesSegment, url.PathEscape(module))
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
	relative, err := location(r.Repo, r.Module, r.Category, r.Agent)
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
	p := Publication{Kind: Kind, SchemaVersion: SchemaVersion, Project: r.Project, Repo: r.Repo, Module: storedModule(r.Module), Category: r.Category, Agent: r.Agent,
		Plan: r.Plan, Run: run, PublishedAt: at.Format(time.RFC3339), Actor: r.Actor, Title: r.Title, Entry: r.Entry, Source: r.Source}
	if p.Files, err = copyRun(temp, files); err == nil {
		err = writeMetadata(filepath.Join(temp, metadataName), p)
	}
	if err == nil {
		err = os.Rename(temp, filepath.Join(directory, run))
	}
	if err != nil {
		_ = os.RemoveAll(temp)
		return View{}, err
	}
	v := newView(p, filepath.Join(directory, run))
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
	for _, f := range files {
		target := filepath.Join(temp, filepath.FromSlash(f.name))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return nil, ioError("create run", err)
		}
		size, sum, err := copyFile(f.path, target)
		if err != nil {
			return nil, err
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
	defer func() { _ = in.Close() }()
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return 0, "", ioError("write run", err)
	}
	hash := sha256.New()
	size, err := io.Copy(io.MultiWriter(out, hash), in)
	if closeErr := out.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return 0, "", ioError("write run", err)
	}
	return size, hex.EncodeToString(hash.Sum(nil)), nil
}

func writeMetadata(target string, p Publication) error {
	return writeJSON(target, metadataName, p)
}

func writeJSON(target, name string, value any) error {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		return ioError("encode "+name, err)
	}
	if err := os.WriteFile(target, buffer.Bytes(), 0o644); err != nil {
		return ioError("write "+name, err)
	}
	return nil
}

// Withdraw hides the runs of the path up to its newest run. Withdrawing a path
// whose newest run is already withdrawn returns the stored withdrawal unchanged.
func (s Store) Withdraw(project, repo, module, category, agent, actor string) (Withdrawal, error) {
	relative, err := location(repo, module, category, agent)
	if err != nil {
		return Withdrawal{}, err
	}
	if err := checkLine("--actor", actor, maxTitle, false); err != nil {
		return Withdrawal{}, err
	}
	directory := filepath.Join(s.root, relative)
	if _, err := os.Stat(directory); errors.Is(err, fs.ErrNotExist) {
		return Withdrawal{}, &Error{Class: ClassNotFound, Detail: "no run in " + filepath.ToSlash(relative)}
	}
	lock, err := filelock.Acquire(filepath.Join(s.root, lockName), filelock.Options{})
	if err != nil {
		return Withdrawal{}, ioError("lock publications", err)
	}
	w, err := s.withdraw(project, repo, storedModule(module), category, agent, actor, relative)
	if releaseErr := lock.Release(); releaseErr != nil && err == nil {
		err = &Error{Class: ClassIOError, Detail: "the withdrawal was stored; releasing the lock failed"}
	}
	return w, err
}

func (s Store) withdraw(project, repo, module, category, agent, actor, relative string) (Withdrawal, error) {
	directory := filepath.Join(s.root, relative)
	runs, err := runNames(directory)
	if err != nil {
		return Withdrawal{}, err
	}
	if len(runs) == 0 {
		return Withdrawal{}, &Error{Class: ClassNotFound, Detail: "no run in " + filepath.ToSlash(relative)}
	}
	newest := runs[len(runs)-1]
	// An unparsable withdrawal is replaced; another read failure stops here.
	previous, found, err := readWithdrawal(directory)
	var readError *Error
	if err != nil && !(errors.As(err, &readError) && readError.Detail == invalidWithdrawal) {
		return Withdrawal{}, err
	}
	if found && previous.Through == newest {
		return previous, nil
	}
	w := Withdrawal{Kind: WithdrawalKind, SchemaVersion: SchemaVersion, Project: project, Repo: repo, Module: module, Category: category, Agent: agent,
		Through: newest, WithdrawnAt: s.now().UTC().Format(time.RFC3339), Actor: actor}
	temp := filepath.Join(directory, "."+withdrawName+tempSuffix)
	if err := writeJSON(temp, withdrawName, w); err != nil {
		_ = os.Remove(temp)
		return Withdrawal{}, err
	}
	if err := os.Rename(temp, filepath.Join(directory, withdrawName)); err != nil {
		_ = os.Remove(temp)
		return Withdrawal{}, ioError("write "+withdrawName, err)
	}
	return w, nil
}

// readWithdrawal returns the withdrawal of directory and whether one exists.
func readWithdrawal(directory string) (Withdrawal, bool, error) {
	contents, err := os.ReadFile(filepath.Join(directory, withdrawName))
	if errors.Is(err, fs.ErrNotExist) {
		return Withdrawal{}, false, nil
	}
	if err != nil {
		return Withdrawal{}, false, ioError("read "+withdrawName, err)
	}
	var w Withdrawal
	if err := json.Unmarshal(contents, &w); err != nil || w.Kind != WithdrawalKind || !runPattern.MatchString(w.Through) {
		return Withdrawal{}, false, &Error{Class: ClassIOError, Detail: invalidWithdrawal}
	}
	return w, true, nil
}

// visibleRuns returns the runs of directory after its withdrawal, ascending.
func visibleRuns(directory string) ([]string, error) {
	runs, err := runNames(directory)
	if err != nil {
		return nil, err
	}
	w, found, err := readWithdrawal(directory)
	if err != nil || !found {
		return runs, err
	}
	return after(runs, w.Through), nil
}

// after returns the runs that sort after through; runs is ascending.
func after(runs []string, through string) []string {
	visible := runs[:0]
	for _, run := range runs {
		if run > through {
			visible = append(visible, run)
		}
	}
	return visible
}

// prune removes runs without plan beyond the newest keepWithout in directory.
func (s Store) prune(directory string) error {
	runs, err := runNames(directory)
	if err != nil {
		return err
	}
	kept := 0
	for i := len(runs) - 1; i >= 0; i-- {
		p, err := decodeMetadata(filepath.Join(directory, runs[i]))
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

// readMetadata reads the publication.json of run, whose path belongs to module
// (empty for the repo's own runs).
func readMetadata(run, module string) (Publication, error) {
	p, err := decodeMetadata(run)
	if err == nil && p.Module != module {
		return Publication{}, &Error{Class: ClassIOError, Detail: metadataName + " names another module than its directory"}
	}
	return p, err
}

func decodeMetadata(run string) (Publication, error) {
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
func (s Store) Show(repo, module, category, agent, run string) (View, error) {
	relative, err := location(repo, module, category, agent)
	if err != nil {
		return View{}, err
	}
	directory := filepath.Join(s.root, relative)
	if run == "" {
		runs, err := visibleRuns(directory)
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
	p, err := readMetadata(filepath.Join(directory, run), storedModule(module))
	if err != nil {
		return View{}, err
	}
	return newView(p, filepath.Join(directory, run)), nil
}

// List returns the runs matching f, newest first; without f.All only the
// newest matching run of each path.
func (s Store) List(f Filter) (ListResult, error) {
	result := ListResult{Items: []View{}, Invalid: []InvalidRun{}}
	if f.Module != "" && f.Modules {
		return result, invalid("--module and --modules exclude each other")
	}
	if err := checkModuleOf(f.Repo, f.Module); err != nil {
		return result, err
	}
	directories, err := s.paths(f, &result)
	if err != nil {
		return result, err
	}
	for _, d := range directories {
		directory := d.directory
		runs, err := runNames(directory)
		if err != nil {
			return result, err
		}
		if w, found, err := readWithdrawal(directory); err != nil {
			// An unreadable withdrawal hides nothing and is reported like an unreadable run.
			result.Invalid = append(result.Invalid, s.invalidRun(filepath.Join(directory, withdrawName), err))
		} else if found {
			runs = after(runs, w.Through)
		}
		for i := len(runs) - 1; i >= 0; i-- {
			run := filepath.Join(directory, runs[i])
			p, err := readMetadata(run, d.module)
			if err != nil {
				result.Invalid = append(result.Invalid, s.invalidRun(run, err))
				continue
			}
			if f.Plan != "" && p.Plan != f.Plan {
				continue
			}
			result.Items = append(result.Items, newView(p, run))
			if !f.All {
				break
			}
		}
	}
	sort.SliceStable(result.Items, func(i, j int) bool { return result.Items[i].Run > result.Items[j].Run })
	return result, nil
}

func (s Store) invalidRun(target string, err error) InvalidRun {
	relative, _ := filepath.Rel(s.root, target)
	return InvalidRun{Run: filepath.ToSlash(relative), Error: err.Error()}
}

// runDirectory is a directory that can hold runs, with the stored module its
// runs must name.
type runDirectory struct {
	directory, module string
}

// paths returns the directories that can hold runs matching f. A directory
// under _modules that is not a canonically escaped module path is reported in
// result.Invalid.
func (s Store) paths(f Filter, result *ListResult) ([]runDirectory, error) {
	var found []runDirectory
	segments, err := subdirectories(s.root)
	if err != nil {
		return nil, err
	}
	for _, segment := range segments {
		if f.Repo != "" && segment != f.Repo {
			continue
		}
		base := filepath.Join(s.root, segment)
		module := storedModule(f.Module)
		if module == "" {
			if found, err = s.categoryPaths(found, base, "", f); err != nil {
				return nil, err
			}
		}
		if !f.Modules && module == "" {
			continue
		}
		names := []string{url.PathEscape(module)}
		if f.Modules {
			if names, err = subdirectories(filepath.Join(base, ModulesSegment)); err != nil {
				return nil, err
			}
		}
		for _, name := range names {
			directory := filepath.Join(base, ModulesSegment, name)
			decoded, err := url.PathUnescape(name)
			if err != nil || decoded == RootModule || CheckModule(decoded) != nil || url.PathEscape(decoded) != name {
				result.Invalid = append(result.Invalid, s.invalidRun(directory, &Error{Class: ClassIOError, Detail: "not an escaped module path"}))
				continue
			}
			if found, err = s.categoryPaths(found, directory, decoded, f); err != nil {
				return nil, err
			}
		}
	}
	return found, nil
}

// categoryPaths adds the category and agent directories of base whose runs
// belong to module. Directories that start with "_" are not categories.
func (s Store) categoryPaths(found []runDirectory, base, module string, f Filter) ([]runDirectory, error) {
	categories, err := subdirectories(base)
	if err != nil {
		return nil, err
	}
	for _, category := range categories {
		if strings.HasPrefix(category, "_") || f.Category != "" && category != f.Category {
			continue
		}
		directory := filepath.Join(base, category)
		if f.Agent == "" {
			found = append(found, runDirectory{directory, module})
		}
		agents, err := subdirectories(directory)
		if err != nil {
			return nil, err
		}
		for _, agent := range agents {
			if namePattern.MatchString(agent) && (f.Agent == "" || agent == f.Agent) {
				found = append(found, runDirectory{filepath.Join(directory, agent), module})
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
