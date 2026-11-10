package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func mr05Snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	e := filepath.WalkDir(root, func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		i, e := d.Info()
		if e != nil {
			return e
		}
		v := i.Mode().String()
		if !d.IsDir() {
			b, e := os.ReadFile(p)
			if e != nil {
				return e
			}
			h := sha256.Sum256(b)
			v += hex.EncodeToString(h[:])
		}
		out[p] = v
		return nil
	})
	if e != nil {
		t.Fatal(e)
	}
	return out
}
func sl01Python(t *testing.T) string {
	t.Helper()
	py := os.Getenv("PLY_SL01_SCHEMA_PYTHON")
	if py == "" {
		py = os.Getenv("PLY_MR05_SCHEMA_PYTHON")
	}
	if py == "" || !filepath.IsAbs(py) {
		t.Fatal("PLY_SL01_SCHEMA_PYTHON (or existing PLY_MR05_SCHEMA_PYTHON) must name an absolute Python with jsonschema; no install or skip")
	}
	return py
}
