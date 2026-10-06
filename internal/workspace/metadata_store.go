package workspace

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

func readWorkspaceMetadata(root, name string) ([]byte, bool, error) {
	dir := filepath.Join(root, MarkerDirectory)
	if err := contentPhysical(dir); err != nil {
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
	return workError(ErrorWorkStoreConflict, fmt.Sprintf("cannot read metadata %s: %v", path, err), err)
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

func (p workspaceMetadataPublisher) publish(root, name string, b []byte) (err error) {
	dir, destination := filepath.Join(root, MarkerDirectory), filepath.Join(root, MarkerDirectory, name)
	if len(b) > workspaceMetadataLimit {
		return WorkInvalidArguments("workspace metadata exceeds 16 MiB")
	}
	// Check the destination even for a first write: a dangling symlink is not absence.
	if _, _, err = readWorkspaceMetadata(root, name); err != nil {
		return err
	}
	if err = p.check("temp-open"); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".metadata-*.tmp")
	if err != nil {
		return workError(ErrorWorkIO, "create metadata temporary file", err)
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
			err = workError(ErrorWorkIO, detail+": "+err.Error(), err)
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
	if _, _, err = readWorkspaceMetadata(root, name); err != nil {
		return err
	}
	if err = platformReplace(temp, destination); err != nil {
		return err
	}
	replaced = true
	if err = p.check("directory-sync"); err != nil {
		return err
	}
	return platformSyncDirectory(dir)
}
