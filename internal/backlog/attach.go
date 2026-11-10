package backlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ClassAttachmentExists: the name is already used on the item or present in its file directory.
const ClassAttachmentExists = "backlog_attachment_exists"

// AttachmentsFormat is the only sidecar format v1 reads and writes.
const AttachmentsFormat = "ply-backlog-attachments/v1"

// Attachment is the metadata of one owned copy at <backlog>/files/<id>/<name>.
type Attachment struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	AddedAt string `json:"added_at"`
}

// sidecar is the exact content of <backlog>/attachments/<id>.json.
type sidecar struct {
	Format      string       `json:"format"`
	ID          string       `json:"id"`
	Attachments []Attachment `json:"attachments"`
}

const (
	attachmentsDirectoryName = "attachments"
	maxAttachments           = 32
)

var errAttachmentExists = &Error{Class: ClassAttachmentExists, Detail: "attachment already exists"}

// Attach copies the regular file at the absolute path source to the existing item id
// and returns all its attachments. The item file is never written.
func (s Store) Attach(id, source string) ([]Attachment, error) {
	if err := s.checkDirectory(); err != nil {
		return nil, err
	}
	if !itemID.MatchString(id) {
		return nil, invalidArguments("invalid item id")
	}
	if !filepath.IsAbs(source) {
		return nil, invalidArguments("source must be an absolute path")
	}
	name := filepath.Base(source)
	if !validFileName(name) {
		return nil, invalidArguments("invalid file name")
	}
	if _, err := s.read(id); err != nil {
		return nil, err
	}
	file, err := openSource(source)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	files := s.filesDir()
	if err := ensureDirectory(files, "files directory"); err != nil {
		return nil, err
	}
	temp, size, err := copyToTemp(files, file)
	if err != nil {
		return nil, err
	}
	defer os.Remove(temp)
	unlock, err := s.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()
	item, err := s.read(id)
	if err != nil {
		return nil, err
	}
	list, err := s.readAttachments(id)
	if err != nil {
		return nil, err
	}
	if len(list) >= maxAttachments {
		return nil, invalidArguments("an item has at most 32 attachments")
	}
	if item.File != nil && strings.EqualFold(item.File.Name, name) || attachmentIndex(list, name, strings.EqualFold) >= 0 {
		return nil, errAttachmentExists
	}
	directory := filepath.Join(files, id)
	if err := ensureDirectory(directory, "file directory"); err != nil {
		return nil, err
	}
	path := filepath.Join(directory, name)
	// Link never replaces, so an existing or unlisted copy is never overwritten.
	if err := os.Link(temp, path); errors.Is(err, fs.ErrExist) {
		return nil, errAttachmentExists
	} else if err != nil {
		return nil, ioError("link file", err)
	}
	list = append(list, Attachment{Name: name, Size: size, AddedAt: s.now().UTC().Format(timestampStyle)})
	if err := s.publishAttachments(id, directory, list); err != nil {
		s.removeUnpublished(id, path, name)
		return nil, err
	}
	return list, nil
}

// Attachments returns the attachments of the existing item id, oldest first; never nil.
func (s Store) Attachments(id string) ([]Attachment, error) {
	if _, err := s.Read(id); err != nil {
		return nil, err
	}
	return s.readAttachments(id)
}

// AttachmentState is FileState for one attachment of item id.
func (s Store) AttachmentState(id string, attachment Attachment) (string, FileState, error) {
	if err := s.checkDirectory(); err != nil {
		return "", "", err
	}
	if !itemID.MatchString(id) {
		return "", "", invalidArguments("invalid item id")
	}
	if !validFileName(attachment.Name) {
		return "", "", invalidArguments("invalid file name")
	}
	return s.copyState(id, attachment.Name)
}

func (s Store) attachmentsPath(id string) string {
	return filepath.Join(s.dir, attachmentsDirectoryName, id+fileSuffix)
}

// readAttachments reads the sidecar of id; a missing sidecar is an empty list.
func (s Store) readAttachments(id string) ([]Attachment, error) {
	path := s.attachmentsPath(id)
	for _, component := range []string{filepath.Dir(path), path} {
		info, err := os.Lstat(component)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			return []Attachment{}, nil
		case err != nil:
			return nil, ioError("check attachments", err)
		case component != path && !info.IsDir(), component == path && !info.Mode().IsRegular():
			return nil, itemInvalid("attachments are not a regular file")
		}
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, ioError("read attachments", err)
	}
	var stored sidecar
	switch err := json.Unmarshal(contents, &stored); {
	case err != nil:
		return nil, itemInvalid("invalid attachments JSON")
	case stored.Format != AttachmentsFormat:
		return nil, itemInvalid("unknown attachments format")
	case stored.ID != id:
		return nil, itemInvalid("attachments id does not match file name")
	case len(stored.Attachments) > maxAttachments:
		return nil, itemInvalid("too many attachments")
	}
	for i, attachment := range stored.Attachments {
		switch {
		case !validFileName(attachment.Name):
			return nil, itemInvalid("invalid attachment name")
		case attachment.Size < 0:
			return nil, itemInvalid("invalid attachment size")
		case attachment.AddedAt == "":
			return nil, itemInvalid("missing attachment added_at")
		case attachmentIndex(stored.Attachments[:i], attachment.Name, strings.EqualFold) >= 0:
			return nil, itemInvalid("duplicate attachment name")
		}
	}
	if stored.Attachments == nil {
		return []Attachment{}, nil
	}
	return stored.Attachments, nil
}

// publishAttachments syncs the new copy and replaces the sidecar, the single publication point.
func (s Store) publishAttachments(id, directory string, list []Attachment) error {
	if err := syncDirectory(directory); err != nil {
		return ioError("sync file directory", err)
	}
	if err := syncDirectory(s.filesDir()); err != nil {
		return ioError("sync files directory", err)
	}
	path := s.attachmentsPath(id)
	if err := ensureDirectory(filepath.Dir(path), "attachments directory"); err != nil {
		return err
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(sidecar{Format: AttachmentsFormat, ID: id, Attachments: list}); err != nil {
		return ioError("encode attachments "+id, err)
	}
	if err := s.replace(path, buffer.Bytes(), 0o644); err != nil {
		return ioError("write attachments "+id, err)
	}
	return nil
}

// removeUnpublished removes the copy linked by a failed Attach only when the sidecar
// provably does not list it; an unreadable sidecar leaves everything as it is.
func (s Store) removeUnpublished(id, path, name string) {
	list, err := s.readAttachments(id)
	if err == nil && attachmentIndex(list, name, func(a, b string) bool { return a == b }) < 0 {
		_ = os.Remove(path)
	}
}

func attachmentIndex(list []Attachment, name string, equal func(string, string) bool) int {
	for i, attachment := range list {
		if equal(attachment.Name, name) {
			return i
		}
	}
	return -1
}

// ensureDirectory creates path if needed and requires a real directory, not a symlink.
func ensureDirectory(path, what string) error {
	if err := os.Mkdir(path, 0o755); err != nil && !errors.Is(err, fs.ErrExist) {
		return ioError("create "+what, err)
	}
	if info, err := os.Lstat(path); err != nil || !info.IsDir() {
		return &Error{Class: ClassIOError, Detail: what + " is not a directory"}
	}
	return nil
}
