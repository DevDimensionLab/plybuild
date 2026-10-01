//go:build darwin

package taskrun

import (
	"errors"
	"golang.org/x/sys/unix"
	"syscall"
	"testing"
)

func TestR11PIDReuseAndGroupScope(t *testing.T) {
	calls := 0
	c := &darwinChild{bound: Process{PID: ptr(41), ProcessGroup: ptr(41), StartIdentity: ptr("bound-birth")}, birth: func(int) (string, error) { return "reused-birth", nil }, group: func(int) (int, error) { return 41, nil }, send: func(group int, s unix.Signal) error {
		if group != -41 {
			t.Fatal("unbound group")
		}
		calls++
		return nil
	}}
	if c.Forward(syscall.SIGTERM) == nil || calls != 0 {
		t.Fatal("signaled reused PID")
	}
	c.birth = func(int) (string, error) { return "bound-birth", nil }
	c.group = func(int) (int, error) { return 42, nil }
	if c.Forward(syscall.SIGHUP) == nil || calls != 0 {
		t.Fatal("signaled unrelated group")
	}
	c.group = func(int) (int, error) { return 41, nil }
	if e := c.Forward(syscall.SIGTERM); e != nil || calls != 1 {
		t.Fatal(e)
	}
	c.waited = true
	if c.Forward(syscall.SIGTERM) == nil || calls != 1 {
		t.Fatal("signaled after Wait")
	}
	c.waited = false
	c.birth = func(int) (string, error) { return "", errors.New("unobservable") }
	if c.Forward(syscall.SIGINT) == nil || calls != 1 {
		t.Fatal("signaled unknown identity")
	}
}
