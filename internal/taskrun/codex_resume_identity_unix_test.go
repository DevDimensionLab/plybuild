//go:build !windows

package taskrun

import (
	"bufio"

	"os"

	"syscall"
	"testing"
)

// This uses an actual OS-held Codex-format thread lock. Only Herdr is modeled;
// it deliberately reports no hook metadata until after the first turn.

func TestCodexResumeWriterHelper(t *testing.T) {
	path := os.Getenv("PLY_TEST_CODEX_LOCK")
	if path == "" {
		return
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if os.Getenv("PLY_TEST_CODEX_LOCK_MODE") != "open" {
		if err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
			t.Fatal(err)
		}
	}
	os.Stdout.WriteString("locked\n")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		mode := scanner.Text()
		if err = syscall.Flock(int(file.Fd()), syscall.LOCK_UN); err != nil {
			t.Fatal(err)
		}
		flags := 0
		switch mode {
		case "shared":
			flags = syscall.LOCK_SH
		case "exclusive":
			flags = syscall.LOCK_EX
		case "unlock":
		default:
			t.Fatal("unknown helper mode")
		}
		if flags != 0 {
			if err = syscall.Flock(int(file.Fd()), flags|syscall.LOCK_NB); err != nil {
				t.Fatal(err)
			}
		}
		os.Stdout.WriteString("locked\n")
	}
}

// A known no-input session_pending attempt has a coordinator-owned exit;
// positive Herdr absence must not hide a surviving unmanaged local writer.

// An unrelated fixture's real process must not look like this fixture's
// competing writer merely because a test adapter reused a native session ID.

// r001-f01: failed local inspection is an unavailable observation, not a
// positive observation of an exiting writer. Never poll a rejected read.
