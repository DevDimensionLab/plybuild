package taskwork

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Metadata is separate from the identity registries so presentation changes do
// not invalidate historical Task, preparation, or integration bindings.
const workspaceMetadataLimit = 16 << 20

type MetadataWorkspace struct {
	Root string `json:"root"`
}
type MetadataSource struct {
	Locator string  `json:"locator"`
	SHA256  *string `json:"sha256"`
	Exists  bool    `json:"exists"`
}

func metadataSource(locator, hash string, exists bool) MetadataSource {
	var digest *string
	if exists {
		digest = &hash
	}
	return MetadataSource{Locator: locator, SHA256: digest, Exists: exists}
}

// readWorkspaceMetadata reads file name in dir. A missing dir or file is
// absence.
func readWorkspaceMetadata(dir, name string) ([]byte, bool, error) {
	err := contentPhysical(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, metadataStoreError(dir, err)
	}
	path := filepath.Join(dir, name)
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, metadataStoreError(path, err)
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
		return nil, false, metadataStoreError(path, errors.New("metadata must be a regular file, not a symlink"))
	}
	b, err := contentReadBounded(path, workspaceMetadataLimit, false)
	if err != nil {
		return nil, false, metadataStoreError(path, err)
	}
	return b, true, nil
}

func decodeWorkspaceMetadata(b []byte, out any) error {
	if _, err := canonicaljson.DecodeStrict(b); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.DisallowUnknownFields()
	return decoder.Decode(out)
}

func metadataStoreError(path string, err error) error {
	return workspace.WorkError(workspace.ErrorWorkStoreConflict, fmt.Sprintf("cannot read metadata %s: %v", path, err), err)
}

// The domain owner holds its existing register lock while publishing. A single
// replacement publishes both current metadata and its complete change history.
type workspaceMetadataPublisher struct{ fault func(string) error }

func (p workspaceMetadataPublisher) check(stage string) error {
	if p.fault != nil {
		return p.fault(stage)
	}
	return nil
}

// publish replaces file name in dir, which must exist.
func (p workspaceMetadataPublisher) publish(dir, name string, b []byte) (err error) {
	destination := filepath.Join(dir, name)
	if len(b) > workspaceMetadataLimit {
		return workspace.WorkInvalidArguments("workspace metadata exceeds 16 MiB")
	}
	// Check the destination even for a first write: a dangling symlink is not absence.
	if _, _, err = readWorkspaceMetadata(dir, name); err != nil {
		return err
	}
	if err = p.check("temp-open"); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".metadata-*.tmp")
	if err != nil {
		return workspace.WorkError(workspace.ErrorWorkIO, "create metadata temporary file", err)
	}
	temp := f.Name()
	replaced := false
	defer func() {
		if e := f.Close(); err == nil && e != nil && !errors.Is(e, os.ErrClosed) {
			err = e
		}
		if !replaced {
			if e := os.Remove(temp); err == nil && e != nil && !errors.Is(e, fs.ErrNotExist) {
				err = e
			}
		}
		if err != nil {
			detail := "metadata publication failed; inspect " + destination + " before retrying"
			if replaced {
				detail = "metadata was replaced but durability could not be confirmed; inspect " + destination + " before retrying"
			}
			err = workspace.WorkError(workspace.ErrorWorkIO, detail+": "+err.Error(), err)
		}
	}()
	if err = p.check("write"); err != nil {
		return err
	}
	n, err := f.Write(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return io.ErrShortWrite
	}
	if err = p.check("file-sync"); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	if err = p.check("replace"); err != nil {
		return err
	}
	// Recheck the physical parent and destination before the effect.
	if _, _, err = readWorkspaceMetadata(dir, name); err != nil {
		return err
	}
	if err = workspace.PlatformReplace(temp, destination); err != nil {
		return err
	}
	replaced = true
	if err = p.check("directory-sync"); err != nil {
		return err
	}
	return workspace.PlatformSyncDirectory(dir)
}
