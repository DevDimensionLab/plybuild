package backlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
)

// Archive error classes shared with the CLI contract.
const (
	ClassArchiveConflict = "backlog_archive_conflict"
	ClassArchiveInvalid  = "backlog_archive_invalid"
	ClassGroupNotFound   = "backlog_group_not_found"
	ClassItemNotDone     = "backlog_item_not_done"
)

// ArchiveFormat is the only group file format v1 reads and writes.
const ArchiveFormat = "ply-backlog-archive-group/v1"

// ArchiveGroup is the exact content of <archive>/<group>.json.
type ArchiveGroup struct {
	Format string         `json:"format"`
	Group  string         `json:"group"`
	Items  []ArchivedItem `json:"items"`
}

// ArchivedItem is an archived done item; it is never an Item so the active format can change freely.
type ArchivedItem struct {
	ID         string         `json:"id"`
	Title      string         `json:"title"`
	Text       string         `json:"text"`
	CreatedAt  string         `json:"created_at"`
	Notes      []ArchivedNote `json:"notes"`
	Repos      []string       `json:"repos,omitempty"`
	Topics     []string       `json:"topics,omitempty"`
	File       *ArchivedFile  `json:"file,omitempty"`
	ArchivedAt string         `json:"archived_at"`
}

// ArchivedNote is one note of an archived item.
type ArchivedNote struct {
	At     string   `json:"at"`
	Text   string   `json:"text"`
	Author string   `json:"author,omitempty"`
	Repos  []string `json:"repos,omitempty"`
	Topics []string `json:"topics,omitempty"`
}

// ArchivedFile is the metadata of the archive's own copy at <archive>/files/<id>/<name>.
type ArchivedFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// GroupSummary is one group as shown by ListGroups.
type GroupSummary struct {
	Group          string `json:"group"`
	ItemCount      int    `json:"item_count"`
	LastArchivedAt string `json:"last_archived_at"`
}

// GroupListResult holds readable groups and unreadable *.json files; both are never nil.
type GroupListResult struct {
	Groups  []GroupSummary `json:"groups"`
	Invalid []InvalidFile  `json:"invalid"`
}

// Archive reads and writes the group files in one directory below root.
type Archive struct {
	dir, root string
	replace   func(string, []byte, fs.FileMode) error
	sync      func(string) error // every sync of Store.Archive; a private fault seam
}

// OpenArchive returns the archive in directory; both paths must be absolute and
// root (the project wrapper) is the last directory synced. Reads never create directory.
func OpenArchive(directory, root string) Archive {
	return Archive{dir: directory, root: root, replace: atomicfile.Replace, sync: syncDirectory}
}

var groupName = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,62}$`)

// ListGroups returns every readable group sorted by name and every unreadable *.json file.
func (a Archive) ListGroups() (GroupListResult, error) {
	result := GroupListResult{Groups: []GroupSummary{}, Invalid: []InvalidFile{}}
	groups, invalid, err := a.readAll()
	if err != nil {
		return result, err
	}
	for _, group := range groups {
		summary := GroupSummary{Group: group.Group, ItemCount: len(group.Items)}
		if len(group.Items) > 0 {
			summary.LastArchivedAt = group.Items[len(group.Items)-1].ArchivedAt
		}
		result.Groups = append(result.Groups, summary)
	}
	result.Invalid = invalid
	return result, nil
}

// ReadGroup returns the validated group file of group.
func (a Archive) ReadGroup(group string) (ArchiveGroup, error) {
	read, _, err := a.ReadGroupFile(group)
	return read, err
}

// ReadGroupFile is ReadGroup that also returns the group file's bytes as stored.
func (a Archive) ReadGroupFile(group string) (ArchiveGroup, []byte, error) {
	if err := a.check(); err != nil {
		return ArchiveGroup{}, nil, err
	}
	if !groupName.MatchString(group) {
		return ArchiveGroup{}, nil, errInvalidGroup
	}
	return a.read(group)
}

// FileState is Store.FileState for the archive's own copy of item.
func (a Archive) FileState(item ArchivedItem) (string, FileState, error) {
	if item.File == nil {
		return "", "", nil
	}
	if err := a.check(); err != nil {
		return "", "", err
	}
	return Store{dir: a.dir}.FileState(Item{ID: item.ID, File: &File{Name: item.File.Name, Size: item.File.Size}})
}

var errInvalidGroup = invalidArguments("group must match ^[a-z0-9][a-z0-9-]{0,62}$")

func (a Archive) check() error {
	if !filepath.IsAbs(a.dir) || !filepath.IsAbs(a.root) {
		return invalidArguments("archive directory and root must be absolute paths")
	}
	return nil
}

func (a Archive) groupPath(group string) string { return filepath.Join(a.dir, group+fileSuffix) }

// readAll reads every *.json file in name order; a missing directory has none.
func (a Archive) readAll() ([]ArchiveGroup, []InvalidFile, error) {
	groups, invalid := []ArchiveGroup{}, []InvalidFile{}
	if err := a.check(); err != nil {
		return groups, invalid, err
	}
	entries, err := os.ReadDir(a.dir)
	if errors.Is(err, fs.ErrNotExist) {
		return groups, invalid, nil
	}
	if err != nil {
		return groups, invalid, ioError("read archive directory", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, fileSuffix) {
			continue
		}
		group, _, err := a.read(strings.TrimSuffix(name, fileSuffix))
		if err != nil {
			invalid = append(invalid, InvalidFile{File: name, Error: detail(err)})
			continue
		}
		groups = append(groups, group)
	}
	// An id in two groups makes every group holding it invalid; no group is picked as the first match.
	holders := map[string]int{}
	for _, group := range groups {
		for _, item := range group.Items {
			holders[item.ID]++
		}
	}
	valid := []ArchiveGroup{}
	for _, group := range groups {
		if slices.ContainsFunc(group.Items, func(item ArchivedItem) bool { return holders[item.ID] > 1 }) {
			invalid = append(invalid, InvalidFile{File: group.Group + fileSuffix, Error: "item id archived in another group"})
			continue
		}
		valid = append(valid, group)
	}
	sort.Slice(invalid, func(i, j int) bool { return invalid[i].File < invalid[j].File })
	return valid, invalid, nil
}

func (a Archive) read(group string) (ArchiveGroup, []byte, error) {
	path := a.groupPath(group)
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return ArchiveGroup{}, nil, &Error{Class: ClassGroupNotFound, Detail: "group " + group + " does not exist; see ply backlog groups"}
	}
	if err != nil {
		return ArchiveGroup{}, nil, ioError("read group", err)
	}
	if !info.Mode().IsRegular() {
		return ArchiveGroup{}, nil, archiveInvalid("group file is not a regular file")
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return ArchiveGroup{}, nil, ioError("read group", err)
	}
	// The decoder would replace invalid bytes, so they are rejected before decoding.
	if !utf8.Valid(contents) {
		return ArchiveGroup{}, nil, archiveInvalid("group file is not valid UTF-8")
	}
	var stored ArchiveGroup
	if err := json.Unmarshal(contents, &stored); err != nil {
		var syntax *json.SyntaxError
		switch {
		case errors.As(err, &syntax):
			return ArchiveGroup{}, nil, archiveInvalid(fmt.Sprintf("invalid JSON at offset %d", syntax.Offset))
		case errors.Is(err, io.ErrUnexpectedEOF):
			return ArchiveGroup{}, nil, archiveInvalid(fmt.Sprintf("invalid JSON at offset %d", len(contents)))
		}
		return ArchiveGroup{}, nil, archiveInvalid("invalid group JSON")
	}
	switch {
	case stored.Format != ArchiveFormat:
		return ArchiveGroup{}, nil, archiveInvalid("unknown format")
	case !groupName.MatchString(group):
		return ArchiveGroup{}, nil, archiveInvalid("invalid group file name")
	case stored.Group != group:
		return ArchiveGroup{}, nil, archiveInvalid("group does not match file name")
	}
	// Author is omitempty, so a present empty author is only visible before it decodes into a string.
	var authors struct {
		Items []struct {
			Notes []struct {
				Author *string `json:"author"`
			} `json:"notes"`
		} `json:"items"`
	}
	_ = json.Unmarshal(contents, &authors)
	for _, item := range authors.Items {
		for _, note := range item.Notes {
			if note.Author != nil && *note.Author == "" {
				return ArchiveGroup{}, nil, archiveInvalid("invalid note author")
			}
		}
	}
	seen := map[string]bool{}
	for i, item := range stored.Items {
		if problem := archivedProblem(item); problem != "" {
			return ArchiveGroup{}, nil, archiveInvalid(problem)
		}
		if seen[item.ID] {
			return ArchiveGroup{}, nil, archiveInvalid("duplicate item id")
		}
		seen[item.ID] = true
		if item.Notes == nil {
			stored.Items[i].Notes = []ArchivedNote{}
		}
	}
	if stored.Items == nil {
		stored.Items = []ArchivedItem{}
	}
	return stored, contents, nil
}

// archivedProblem applies the v1 limits frozen from 92a14b8; the shared helpers
// may serve only while they give exactly these limits.
func archivedProblem(item ArchivedItem) string {
	switch {
	case !itemID.MatchString(item.ID):
		return "invalid item id"
	case item.Title == "" || item.CreatedAt == "" || item.ArchivedAt == "":
		return "missing title, created_at or archived_at"
	case item.Text == "" && item.File == nil:
		return "missing text"
	case item.File != nil && (!validFileName(item.File.Name) || item.File.Size < 0):
		return "invalid file"
	case !validTags(item.Repos, validRepo):
		return "invalid repos"
	case !validTags(item.Topics, validLabel):
		return "invalid topics"
	}
	for _, note := range item.Notes {
		switch {
		case note.Author != "" && !validLabel(note.Author):
			return "invalid note author"
		case !validTags(note.Repos, validRepo):
			return "invalid note repos"
		case !validTags(note.Topics, validLabel):
			return "invalid note topics"
		}
	}
	return ""
}

// Archive moves the done item id into group, creating the group only when
// newGroup is set. A repeat after an interrupted or completed archiving
// finishes it without writing a duplicate; see docs for the order of steps.
func (s Store) Archive(id string, archive Archive, group string, newGroup bool) (ArchivedItem, error) {
	if err := s.checkDirectory(); err != nil {
		return ArchivedItem{}, err
	}
	if err := archive.check(); err != nil {
		return ArchivedItem{}, err
	}
	if !itemID.MatchString(id) {
		return ArchivedItem{}, invalidArguments("invalid item id")
	}
	if !groupName.MatchString(group) {
		return ArchivedItem{}, errInvalidGroup
	}
	if sameDirectory(s.dir, archive.dir) {
		return ArchivedItem{}, invalidArguments("backlog and archive directories must differ")
	}
	// Without a backlog directory only a finished archiving can be confirmed; nothing is locked or removed.
	locked := true
	if _, err := os.Stat(s.dir); errors.Is(err, fs.ErrNotExist) {
		locked = false
	} else {
		unlock, err := s.lock()
		if err != nil {
			return ArchivedItem{}, err
		}
		defer unlock()
	}
	if _, err := os.Lstat(s.attachmentsPath(id)); err == nil {
		return ArchivedItem{}, itemInvalid("item has attachments; archive does not support them yet")
	} else if !errors.Is(err, fs.ErrNotExist) {
		return ArchivedItem{}, ioError("check attachments", err)
	}
	groups, invalid, err := archive.readAll()
	if err != nil {
		return ArchivedItem{}, err
	}
	if len(invalid) > 0 {
		return ArchivedItem{}, archiveInvalid("group file " + invalid[0].File + " is invalid")
	}
	var target, foundIn *ArchiveGroup
	var found *ArchivedItem
	for i := range groups {
		if groups[i].Group == group {
			target = &groups[i]
		}
		for j, entry := range groups[i].Items {
			if entry.ID == id {
				found, foundIn = &groups[i].Items[j], &groups[i]
			}
		}
	}
	item, err := s.read(id)
	var backlogError *Error
	if errors.As(err, &backlogError) && backlogError.Class == ClassItemNotFound && found != nil {
		if foundIn != target {
			return ArchivedItem{}, archiveConflict("item " + id + " is archived in group " + foundIn.Group)
		}
		if err := archive.requireCopy(*found); err != nil {
			return ArchivedItem{}, err
		}
		return *found, s.finish(archive, group, *found, locked)
	}
	if err != nil {
		return ArchivedItem{}, err
	}
	if item.Status != StatusDone {
		return ArchivedItem{}, &Error{Class: ClassItemNotDone, Detail: "item " + id + " is not done"}
	}
	if item.File != nil {
		path, state, err := s.FileState(item)
		if err != nil {
			return ArchivedItem{}, err
		}
		if info, statErr := os.Lstat(path); state != FilePresent || statErr != nil || info.Size() != item.File.Size {
			return ArchivedItem{}, itemInvalid("file copy is missing or changed")
		}
	}
	if found != nil && foundIn != target {
		return ArchivedItem{}, archiveConflict("item " + id + " is archived in group " + foundIn.Group)
	}
	archived := toArchived(item, s.now().UTC().Format(timestampStyle))
	if found != nil {
		if !sameArchived(*found, archived) {
			return ArchivedItem{}, archiveConflict("archived item " + id + " differs from the active item")
		}
		if err := archive.requireCopy(*found); err != nil {
			return ArchivedItem{}, err
		}
		if item.File != nil {
			active, _, _ := s.FileState(item)
			copied, _, _ := archive.FileState(*found)
			if same, err := sameContent(active, copied); err != nil || !same {
				return ArchivedItem{}, errors.Join(err, archiveConflict("archived file copy of "+id+" differs from the active copy"))
			}
		}
		return *found, s.finish(archive, group, *found, locked)
	}
	switch {
	case target == nil && !newGroup:
		_, _, err := archive.read(group)
		return ArchivedItem{}, err
	case target != nil && newGroup:
		return ArchivedItem{}, archiveConflict("group " + group + " already exists")
	case target == nil:
		target = &ArchiveGroup{Format: ArchiveFormat, Group: group, Items: []ArchivedItem{}}
	}
	if item.File != nil {
		if err := s.copyToArchive(item, archive); err != nil {
			return ArchivedItem{}, err
		}
	}
	target.Items = append(target.Items, archived)
	contents, err := encodeIndented(target)
	if err != nil {
		return ArchivedItem{}, ioError("encode group", err)
	}
	if err := os.MkdirAll(archive.dir, 0o755); err != nil {
		return ArchivedItem{}, ioError("create archive directory", err)
	}
	if err := archive.replace(archive.groupPath(group), contents, 0o644); err != nil {
		return ArchivedItem{}, ioError("write group "+group, err)
	}
	return archived, s.finish(archive, group, archived, locked)
}

// sameDirectory reports whether two paths name one directory, also through a symlink or case alias.
func sameDirectory(a, b string) bool {
	if filepath.Clean(a) == filepath.Clean(b) {
		return true
	}
	first, errA := os.Stat(a)
	second, errB := os.Stat(b)
	return errA == nil && errB == nil && os.SameFile(first, second)
}

func toArchived(item Item, at string) ArchivedItem {
	notes := []ArchivedNote{}
	for _, note := range item.Notes {
		notes = append(notes, ArchivedNote{At: note.At, Text: note.Text, Author: note.Author, Repos: note.Repos, Topics: note.Topics})
	}
	archived := ArchivedItem{ID: item.ID, Title: item.Title, Text: item.Text, CreatedAt: item.CreatedAt, Notes: notes, Repos: item.Repos, Topics: item.Topics, ArchivedAt: at}
	if item.File != nil {
		archived.File = &ArchivedFile{Name: item.File.Name, Size: item.File.Size}
	}
	return archived
}

// sameArchived compares everything except archived_at.
func sameArchived(a, b ArchivedItem) bool {
	return a.ID == b.ID && a.Title == b.Title && a.Text == b.Text && a.CreatedAt == b.CreatedAt &&
		slices.EqualFunc(a.Notes, b.Notes, sameArchivedNote) && slices.Equal(a.Repos, b.Repos) && slices.Equal(a.Topics, b.Topics) &&
		(a.File == nil) == (b.File == nil) && (a.File == nil || *a.File == *b.File)
}

func sameArchivedNote(a, b ArchivedNote) bool {
	return a.At == b.At && a.Text == b.Text && a.Author == b.Author && slices.Equal(a.Repos, b.Repos) && slices.Equal(a.Topics, b.Topics)
}

// requireCopy proves the archive's copy of a file item is a regular file of the recorded size.
func (a Archive) requireCopy(item ArchivedItem) error {
	if item.File == nil {
		return nil
	}
	path, state, err := a.FileState(item)
	if err != nil {
		return err
	}
	if info, statErr := os.Lstat(path); state != FilePresent || statErr != nil || info.Size() != item.File.Size {
		return archiveInvalid("archived file copy of " + item.ID + " is missing or changed")
	}
	return nil
}

// copyToArchive places the archive's copy before the group file is published,
// reusing a byte-equal leftover and never overwriting anything else.
func (s Store) copyToArchive(item Item, archive Archive) error {
	files := filepath.Join(archive.dir, filesDirectoryName)
	directory := filepath.Join(files, item.ID)
	if err := os.MkdirAll(archive.dir, 0o755); err != nil {
		return ioError("create archive directory", err)
	}
	for _, path := range []string{files, directory} {
		if err := os.Mkdir(path, 0o755); err != nil && !errors.Is(err, fs.ErrExist) {
			return ioError("create archive file directory", err)
		}
		if info, err := os.Lstat(path); err != nil || !info.IsDir() {
			return &Error{Class: ClassIOError, Detail: "archive file directory is not a directory"}
		}
	}
	active, _, _ := s.FileState(item)
	target := filepath.Join(directory, item.File.Name)
	if info, err := os.Lstat(target); err == nil {
		if !info.Mode().IsRegular() {
			return archiveConflict("archived file copy of " + item.ID + " is not a regular file")
		}
		same, err := sameContent(active, target)
		if err == nil && !same {
			err = archiveConflict("archived file copy of " + item.ID + " differs from the active copy")
		}
		return err
	} else if !errors.Is(err, fs.ErrNotExist) {
		return ioError("check archived file copy", err)
	}
	source, err := openSource(active)
	if err != nil {
		return &Error{Class: ClassIOError, Detail: "open file copy"}
	}
	defer source.Close()
	temp, size, err := copyToTemp(directory, source)
	if err != nil {
		return err
	}
	if size != item.File.Size {
		_ = os.Remove(temp)
		return &Error{Class: ClassIOError, Detail: "file copy changed while copying"}
	}
	if err := os.Rename(temp, target); err != nil {
		_ = os.Remove(temp)
		return ioError("move archived file copy", err)
	}
	for _, path := range []string{directory, files} {
		if err := archive.sync(path); err != nil {
			return ioError("sync archive file directory", err)
		}
	}
	return nil
}

// finish makes the archived entry durable (H), then, when the backlog is
// locked, removes the active item file, syncs the backlog and cleans up (R).
func (s Store) finish(archive Archive, group string, item ArchivedItem, locked bool) error {
	paths := []string{archive.groupPath(group)}
	if item.File != nil {
		files := filepath.Join(archive.dir, filesDirectoryName)
		directory := filepath.Join(files, item.ID)
		paths = append(paths, filepath.Join(directory, item.File.Name), directory, files)
	}
	for directory := archive.dir; ; directory = filepath.Dir(directory) {
		paths = append(paths, directory)
		if directory == archive.root || directory == filepath.Dir(directory) {
			break
		}
	}
	for _, path := range paths {
		if err := archive.sync(path); err != nil {
			return ioError("sync archive", err)
		}
	}
	if !locked {
		return nil
	}
	if err := os.Remove(s.path(item.ID)); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return ioError("remove item "+item.ID, err)
	}
	if err := archive.sync(s.dir); err != nil {
		return ioError("sync backlog directory", err)
	}
	return s.cleanup(archive, item)
}

// cleanup removes the active copy of an archived file item when it is byte-equal to the archive's copy.
func (s Store) cleanup(archive Archive, archived ArchivedItem) error {
	if archived.File == nil {
		return nil
	}
	path, state, err := s.FileState(Item{ID: archived.ID, File: &File{Name: archived.File.Name, Size: archived.File.Size}})
	switch {
	case err != nil:
		return err
	case state == FileMissing:
		return nil
	case state == FileUnexpected:
		return archiveConflict("active file copy of " + archived.ID + " is not a regular file")
	}
	copied, _, _ := archive.FileState(archived)
	same, err := sameContent(path, copied)
	if err != nil {
		return err
	}
	directory := filepath.Dir(path)
	entries, err := os.ReadDir(directory)
	if err != nil {
		return ioError("read file directory", err)
	}
	if !same || len(entries) != 1 {
		return archiveConflict("active file copy of " + archived.ID + " differs from the archived copy")
	}
	if err := os.Remove(path); err != nil {
		return ioError("remove file copy", err)
	}
	if err := os.Remove(directory); err != nil {
		return ioError("remove file directory", err)
	}
	if err := archive.sync(s.filesDir()); err != nil {
		return ioError("sync files directory", err)
	}
	return nil
}

// sameContent streams two regular files with the copyToTemp buffer size.
func sameContent(first, second string) (bool, error) {
	a, err := openSource(first)
	if err != nil {
		return false, &Error{Class: ClassIOError, Detail: "open file copy"}
	}
	defer a.Close()
	b, err := openSource(second)
	if err != nil {
		return false, &Error{Class: ClassIOError, Detail: "open file copy"}
	}
	defer b.Close()
	bufferA, bufferB := make([]byte, copyBufferBytes), make([]byte, copyBufferBytes)
	for {
		n, errA := io.ReadFull(a, bufferA)
		m, errB := io.ReadFull(b, bufferB)
		if !bytes.Equal(bufferA[:n], bufferB[:m]) {
			return false, nil
		}
		endA := errors.Is(errA, io.EOF) || errors.Is(errA, io.ErrUnexpectedEOF)
		endB := errors.Is(errB, io.EOF) || errors.Is(errB, io.ErrUnexpectedEOF)
		if err := errors.Join(errA, errB); err != nil && !(endA || endB) {
			return false, ioError("read file copy", err)
		}
		if endA || endB {
			return endA && endB, nil
		}
	}
}

func archiveInvalid(detail string) error  { return &Error{Class: ClassArchiveInvalid, Detail: detail} }
func archiveConflict(detail string) error { return &Error{Class: ClassArchiveConflict, Detail: detail} }
