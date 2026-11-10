package privatestore

import "path/filepath"

// PrivateWriteOnce publishes complete synced bytes once; identical bytes are a
// no-op and different bytes conflict. Missing parents are created private.
func PrivateWriteOnce(path string, b []byte) error { return writeOnce(path, b) }

// PrivateRead reads a bounded private regular file without following symlinks.
func PrivateRead(path string, max int) ([]byte, error) { return readFile(path, max, true) }

// PhysicalPath rejects relative, unclean or symlinked paths. A missing path is
// an os.IsNotExist error; nothing is created.
func PhysicalPath(path string) error { return physical(path, false) }

// StrictDecode rejects duplicate, unknown or missing members (including
// explicit nulls), noncanonical numbers, lone surrogates and trailing values.
func StrictDecode(b []byte, max int, out any) error { return decode(b, max, out) }

// WithPrivateLock creates the private directory and holds dir/store.lock
// exclusively while fn runs.
func WithPrivateLock(dir string, fn func() error) error {
	if e := privateDir(dir); e != nil {
		return e
	}
	f, e := lockFile(filepath.Join(dir, "store.lock"))
	if e != nil {
		return e
	}
	defer unlockFile(f)
	return fn()
}
