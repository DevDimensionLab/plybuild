package backlog

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// File is the metadata of the owned copy at <backlog>/files/<id>/<name>.
type File struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// FileState is what FileState finds at an item's copy path.
type FileState string

const (
	FilePresent    FileState = "present"    // regular file (Lstat)
	FileMissing    FileState = "missing"    // some path component is absent
	FileUnexpected FileState = "unexpected" // symlink, directory or other type
)

const (
	filesDirectoryName = "files"
	copyBufferBytes    = 32 * 1024
)

// AddFile creates an item with empty text and a streamed copy of the regular
// file at the absolute path source. The source is opened read-only and never changed.
// The item file is written last; a failed write removes the copy only when the item is proven absent.
func (s Store) AddFile(title, source string) (Item, error) {
	if err := s.checkDirectory(); err != nil {
		return Item{}, err
	}
	title, err := checkTitle(title)
	if err != nil {
		return Item{}, err
	}
	if !filepath.IsAbs(source) {
		return Item{}, invalidArguments("source must be an absolute path")
	}
	name := filepath.Base(source)
	if !validFileName(name) {
		return Item{}, invalidArguments("invalid file name")
	}
	file, err := openSource(source)
	if err != nil {
		return Item{}, err
	}
	defer file.Close()
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return Item{}, ioError("create backlog directory", err)
	}
	files := s.filesDir()
	if err := ensureDirectory(files, "files directory"); err != nil {
		return Item{}, err
	}
	temp, size, err := copyToTemp(files, file)
	if err != nil {
		return Item{}, err
	}
	moved := false
	defer func() {
		if !moved {
			_ = os.Remove(temp)
		}
	}()
	unlock, err := s.lock()
	if err != nil {
		return Item{}, err
	}
	defer unlock()
	now := s.now().UTC()
	for attempt := 0; attempt < maxIDAttempts; attempt++ {
		id, taken, err := s.candidateID(now)
		if err != nil {
			return Item{}, err
		}
		if taken {
			continue
		}
		// A fresh directory created under the lock is what keeps copies from being overwritten or reused.
		directory := filepath.Join(files, id)
		if err := os.Mkdir(directory, 0o755); errors.Is(err, fs.ErrExist) {
			continue
		} else if err != nil {
			return Item{}, ioError("create file directory", err)
		}
		path := filepath.Join(directory, name)
		if err := os.Rename(temp, path); err != nil {
			_ = os.Remove(directory)
			return Item{}, ioError("move file", err)
		}
		moved = true
		item := Item{ID: id, Title: title, Text: "", CreatedAt: now.Format(timestampStyle), Status: StatusOpen, Notes: []Note{}, File: &File{Name: name, Size: size}}
		if err := s.publish(item, directory); err != nil {
			// Replaced=false does not prove absence; only a missing item file does.
			if _, statErr := os.Lstat(s.path(id)); errors.Is(statErr, fs.ErrNotExist) {
				_ = os.Remove(path)
				_ = os.Remove(directory)
			}
			return Item{}, err
		}
		return item, nil
	}
	return Item{}, errNoUniqueID
}

// FileState returns the copy path of item and what is there, without reading
// content or following symlinks. An item without file gives ("", "", nil).
func (s Store) FileState(item Item) (string, FileState, error) {
	if item.File == nil {
		return "", "", nil
	}
	if err := s.checkDirectory(); err != nil {
		return "", "", err
	}
	if !itemID.MatchString(item.ID) {
		return "", "", invalidArguments("invalid item id")
	}
	if !validFileName(item.File.Name) {
		return "", "", invalidArguments("invalid file name")
	}
	return s.copyState(item.ID, item.File.Name)
}

// copyState Lstats files/, files/<id>/ and the copy; it neither reads content nor follows symlinks.
func (s Store) copyState(id, name string) (string, FileState, error) {
	files := s.filesDir()
	directory := filepath.Join(files, id)
	path := filepath.Join(directory, name)
	for _, component := range []string{files, directory, path} {
		info, err := os.Lstat(component)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			return path, FileMissing, nil
		case err != nil:
			return "", "", ioError("check file", err)
		case component != path && !info.IsDir(), component == path && !info.Mode().IsRegular():
			return path, FileUnexpected, nil
		}
	}
	return path, FilePresent, nil
}

func (s Store) filesDir() string { return filepath.Join(s.dir, filesDirectoryName) }

// copyToTemp streams source into a new temp file in directory and returns its path and byte count.
func copyToTemp(directory string, source *os.File) (string, int64, error) {
	temp, err := os.CreateTemp(directory, ".import-*")
	if err != nil {
		return "", 0, ioError("create temp file", err)
	}
	size, err := io.CopyBuffer(temp, source, make([]byte, copyBufferBytes))
	if err == nil {
		err = temp.Chmod(0o644)
	}
	if err == nil {
		err = temp.Sync()
	}
	if closeErr := temp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(temp.Name())
		return "", 0, ioError("copy file", err)
	}
	return temp.Name(), size, nil
}

func (s Store) publish(item Item, directory string) error {
	if err := syncDirectory(directory); err != nil {
		return ioError("sync file directory", err)
	}
	if err := syncDirectory(s.filesDir()); err != nil {
		return ioError("sync files directory", err)
	}
	contents, err := encode(item)
	if err != nil {
		return err
	}
	if err := s.replace(s.path(item.ID), contents, 0o644); err != nil {
		return ioError("write item "+item.ID, err)
	}
	return nil
}

// validFileName accepts one UTF-8 path element; the name is never a path.
func validFileName(name string) bool {
	return name != "" && name != "." && name != ".." && utf8.ValidString(name) && !strings.ContainsAny(name, "/\\\x00\r\n")
}

var errNotRegular = invalidArguments("file must be a regular file")

// openError keeps only the operating-system cause so details never carry paths.
func openError(err error) error {
	var pathError *fs.PathError
	if errors.As(err, &pathError) {
		err = pathError.Err
	}
	return invalidArguments("open file: " + err.Error())
}
