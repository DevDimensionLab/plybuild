package clirelease

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
)

func TestInstallerReplacePreservesAtomicfileError(t *testing.T) {
	root := t.TempDir()
	for _, tc := range []struct{ name, path, op string }{
		{"missing parent", filepath.Join(root, "missing", "file"), "open temp"},
		{"directory destination", filepath.Join(root, "directory"), "replace"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.op == "replace" {
				if err := os.Mkdir(tc.path, 0700); err != nil {
					t.Fatal(err)
				}
			}
			err := (Installer{}).replace(tc.path, []byte("complete"), 0600)
			var failure *atomicfile.Error
			if !errors.As(err, &failure) || failure.Op != tc.op || failure.Replaced {
				t.Fatalf("lost replacement context: %T %v", err, err)
			}
			if tc.op == "open temp" && !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("lost absence cause: %v", err)
			}
			if tc.op == "replace" {
				if st, e := os.Stat(tc.path); e != nil || !st.IsDir() {
					t.Fatalf("directory changed: %v", e)
				}
			}
			leftovers, _ := filepath.Glob(filepath.Join(root, ".atomicfile-*.tmp"))
			if len(leftovers) != 0 {
				t.Fatalf("temporary remains: %v", leftovers)
			}
		})
	}
}

func TestInstallerAtomicfileBackupReceiptAndRollback(t *testing.T) {
	dest, m, b := installFixture(t)
	old, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	receipt, err := (Installer{}).install(m, b, Digest([]byte("atomicfile manifest")), dest)
	if err != nil || receipt.State != "installed" {
		t.Fatalf("install: %+v %v", receipt, err)
	}
	check := func(path string, want []byte, mode os.FileMode) {
		t.Helper()
		got, e := os.ReadFile(path)
		if e != nil || !bytes.Equal(got, want) {
			t.Fatalf("bytes %s: %v", path, e)
		}
		st, e := os.Stat(path)
		if e != nil {
			t.Fatal(e)
		}
		if runtime.GOOS != "windows" && st.Mode().Perm() != mode {
			t.Fatalf("mode %s: %v", path, st.Mode())
		}
	}
	check(dest, b, 0755)
	check(receipt.Backup, old, 0500)
	receiptPath := dest + ".install.json"
	raw, err := os.ReadFile(receiptPath)
	if err != nil {
		t.Fatal(err)
	}
	check(receiptPath, raw, 0600)
	receipt, err = (Installer{}).Rollback(dest, m.SHA256)
	if err != nil || receipt.State != "rolled_back" {
		t.Fatalf("rollback: %+v %v", receipt, err)
	}
	check(dest, old, 0755)
	check(receipt.Backup, old, 0500)
	saved, err := readReceipt(receiptPath)
	if err != nil || saved.State != "rolled_back" {
		t.Fatalf("receipt: %+v %v", saved, err)
	}
}
