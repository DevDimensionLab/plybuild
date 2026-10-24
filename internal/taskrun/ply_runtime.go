package taskrun

import "path/filepath"

// A command must exist to launch a child. Its historical path, bytes and
// build identity are not authority for an existing Run.
func currentPlyPath(d Dependencies) (string, error) {
	path, err := d.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(path)
}
