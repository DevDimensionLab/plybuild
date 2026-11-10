//go:build !windows

package backlog

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestAddFileRejectsFIFOWithoutBlocking(t *testing.T) {
	store := testStore(t)
	fifo := filepath.Join(t.TempDir(), "pipe")
	if err := unix.Mkfifo(fifo, 0o644); err != nil {
		t.Fatal(err)
	}
	before := backlogTree(t, store)
	result := make(chan error, 1)
	go func() { _, err := store.AddFile("t", fifo); result <- err }()
	select {
	case err := <-result:
		if err == nil || err.Error() != ClassInvalidArguments+": file must be a regular file" {
			t.Fatalf("add file = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("AddFile blocked on FIFO")
	}
	if fmt.Sprint(backlogTree(t, store)) != fmt.Sprint(before) {
		t.Fatal("backlog changed")
	}
}

func TestAttachRejectsFIFOWithoutBlocking(t *testing.T) {
	store, item := textItem(t)
	fifo := filepath.Join(t.TempDir(), "pipe")
	if err := unix.Mkfifo(fifo, 0o644); err != nil {
		t.Fatal(err)
	}
	before := backlogTree(t, store)
	result := make(chan error, 1)
	go func() { _, err := store.Attach(item.ID, fifo); result <- err }()
	select {
	case err := <-result:
		if err == nil || err.Error() != ClassInvalidArguments+": file must be a regular file" {
			t.Fatalf("attach = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Attach blocked on FIFO")
	}
	requireTree(t, store, before)
}
