//go:build !windows

package agentwork

import (
	"os"
	"path/filepath"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestQualifiedCallbackRejectsNonregularBeforeClaimWithoutBlocking(t *testing.T) {
	for _, mode := range []string{"fifo", "directory", "swapped-fifo", "swapped-symlink"} {
		t.Run(mode, func(t *testing.T) {
			s, d, calls := fixture(t)
			d.NextCheckAtUTC = s.now()
			f, e := s.Register(d)
			if e != nil {
				t.Fatal(e)
			}
			path := filepath.Join(s.Root, "qualified-ply")
			if e = os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0700); e != nil {
				t.Fatal(e)
			}
			qualifyTestCallback(t, s.Root, path)
			replace := func() {
				t.Helper()
				if e := os.Remove(path); e != nil {
					t.Fatal(e)
				}
				switch mode {
				case "directory":
					e = os.Mkdir(path, 0700)
				case "swapped-symlink":
					target := filepath.Join(s.Root, "other-ply")
					e = os.WriteFile(target, []byte("#!/bin/sh\nexit 0\n"), 0700)
					if e == nil {
						e = os.Symlink(target, path)
					}
				default:
					e = unix.Mkfifo(path, 0700)
				}
				if e != nil {
					t.Fatal(e)
				}
			}
			if mode == "fifo" || mode == "directory" {
				replace()
			}
			var claims atomic.Int32
			var replacements atomic.Int32
			s.Fault = func(stage string) error {
				if stage == "before_callback_open" && (mode == "swapped-fifo" || mode == "swapped-symlink") {
					replace()
					replacements.Add(1)
				}
				if stage == "after_claim" {
					claims.Add(1)
				}
				return nil
			}
			before, e := s.read()
			if e != nil {
				t.Fatal(e)
			}
			done := make(chan error, 1)
			go func() { _, err := s.Dispatch(f.ID); done <- err }()
			blocked := false
			select {
			case e = <-done:
			case <-time.After(500 * time.Millisecond):
				blocked = true
				fd, err := unix.Open(path, unix.O_RDWR|unix.O_NONBLOCK, 0)
				if err != nil {
					t.Fatal(err)
				}
				unix.Close(fd)
				select {
				case e = <-done:
				case <-time.After(time.Second):
					t.Fatal("dispatch did not exit after FIFO rescue")
				}
			}
			if blocked {
				t.Error("invalid nonregular executor blocked dispatch and held the store lock")
			}
			if e == nil || e.Error() != "followup_qualified_callback_unavailable" {
				t.Errorf("invalid executor result: %v", e)
			}
			if (mode == "swapped-fifo" || mode == "swapped-symlink") && replacements.Load() != 1 {
				t.Error("open-boundary replacement was not exercised")
			}
			after, err := s.read()
			if err != nil || !reflect.DeepEqual(before, after) || calls.Load() != 0 || claims.Load() != 0 {
				t.Error("invalid executor changed store or claimed/sent", err, calls.Load(), claims.Load())
			}
		})
	}
}
