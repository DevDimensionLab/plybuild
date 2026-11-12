//go:build !windows

package plan

import (
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"testing"
	"time"
)

// TestLockFreeReadSeesWholeFileAcrossConcurrentReplacement replaces the plan
// file exactly between a reader's inspection and its open, as a concurrent
// writer may, and checks the reader still gets one complete revision while
// symlinks and FIFOs swapped in at that moment stay refused without hanging.
func TestLockFreeReadSeesWholeFileAcrossConcurrentReplacement(t *testing.T) {
	setup := func(t *testing.T) (Store, Store, string) {
		writer, dir := testStore(t)
		if _, err := writer.Add([]byte(specExample), "a", "r"); err != nil {
			t.Fatal(err)
		}
		reader := Open(dir)
		return reader, writer, filepath.Join(dir, "PLAN-1.json")
	}
	readWithin := func(t *testing.T, read func() error) error {
		t.Helper()
		result := make(chan error, 1)
		go func() { result <- read() }()
		select {
		case err := <-result:
			return err
		case <-time.After(5 * time.Second):
			t.Fatal("read did not return")
			return nil
		}
	}

	t.Run("atomic update", func(t *testing.T) {
		reader, writer, _ := setup(t)
		updated := edit(t, func(d map[string]any) { d["summary"] = "Ny revisjon" })
		once := sync.Once{}
		reader.beforeOpen = func(string) {
			once.Do(func() {
				if _, err := writer.Update("PLAN-1", 1, updated, "a", "under read"); err != nil {
					t.Error(err)
				}
			})
		}
		v, err := reader.Show("PLAN-1", 0)
		if err != nil || v.Revision != 2 || v.Plan.Summary != "Ny revisjon" {
			t.Fatalf("show during replacement %+v %v", v, err)
		}
		once = sync.Once{}
		reader.beforeOpen = func(string) {
			once.Do(func() {
				if _, err := writer.Update("PLAN-1", 2, []byte(specExample), "a", "under list"); err != nil {
					t.Error(err)
				}
			})
		}
		list, err := reader.List("")
		if err != nil || len(list.Items) != 1 || list.Items[0].Revision != 3 || len(list.Invalid) != 0 {
			t.Fatalf("list during replacement %+v %v", list, err)
		}
	})

	t.Run("symlink swapped in", func(t *testing.T) {
		reader, _, path := setup(t)
		target := filepath.Join(t.TempDir(), "elsewhere.json")
		if err := os.WriteFile(target, []byte(readFile(t, path)), 0o644); err != nil {
			t.Fatal(err)
		}
		reader.beforeOpen = func(path string) {
			if err := os.Remove(path); err != nil {
				t.Error(err)
			}
			if err := os.Symlink(target, path); err != nil {
				t.Error(err)
			}
			reader.beforeOpen = nil
		}
		_, err := reader.Show("PLAN-1", 0)
		requireClass(t, err, ClassInvalid)
	})

	t.Run("FIFO swapped in", func(t *testing.T) {
		reader, _, _ := setup(t)
		reader.beforeOpen = func(path string) {
			if err := os.Remove(path); err != nil {
				t.Error(err)
			}
			if err := syscall.Mkfifo(path, 0o644); err != nil {
				t.Error(err)
			}
		}
		err := readWithin(t, func() error { _, err := reader.Show("PLAN-1", 0); return err })
		requireClass(t, err, ClassInvalid)
		// A FIFO found already at inspection is refused the same way.
		reader.beforeOpen = nil
		err = readWithin(t, func() error { _, err := reader.Show("PLAN-1", 0); return err })
		requireClass(t, err, ClassInvalid)
		list, err := reader.List("")
		if err != nil || len(list.Invalid) != 1 || list.Invalid[0].File != "PLAN-1.json" {
			t.Fatalf("list %+v %v", list, err)
		}
	})
}
