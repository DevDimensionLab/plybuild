//go:build darwin

package taskrun

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type darwinTerminal struct{}

func (darwinTerminal) Check() error {
	for _, f := range []*os.File{os.Stdin, os.Stdout, os.Stderr} {
		if !term.IsTerminal(int(f.Fd())) {
			return failure("task_run_tty_required", 3, "Interactive start requires stdin, stdout and stderr on the terminal.")
		}
	}
	p, e := unix.IoctlGetInt(int(os.Stdin.Fd()), unix.TIOCGPGRP)
	if e != nil || p != unix.Getpgrp() {
		return failure("task_run_tty_required", 3, "Ply must own the foreground terminal.")
	}
	return nil
}
func foreground(p int) error {
	ignored := signal.Ignored(syscall.SIGTTOU)
	signal.Ignore(syscall.SIGTTOU)
	if !ignored {
		defer signal.Reset(syscall.SIGTTOU)
	}
	return unix.IoctlSetPointerInt(int(os.Stdin.Fd()), unix.TIOCSPGRP, p)
}
func (darwinTerminal) Foreground(p int) error { return foreground(p) }
func (darwinTerminal) Capture() (func() error, error) {
	fd := int(os.Stdin.Fd())
	state, e := term.GetState(fd)
	if e != nil {
		return nil, e
	}
	group, e := unix.IoctlGetInt(fd, unix.TIOCGPGRP)
	if e != nil {
		return nil, e
	}
	return func() error {
		e := foreground(group)
		r := term.Restore(fd, state)
		if e != nil {
			return e
		}
		return r
	}, nil
}

type darwinChild struct {
	cmd    *exec.Cmd
	bound  Process
	mu     sync.Mutex
	waited bool
	birth  func(int) (string, error)
	group  func(int) (int, error)
	send   func(int, unix.Signal) error
}

func processBirth(pid int) (string, error) {
	boot, e := unix.Sysctl("kern.bootsessionuuid")
	if e != nil {
		return "", e
	}
	p, e := unix.SysctlKinfoProc("kern.proc.pid", pid)
	if e != nil {
		return "", e
	}
	if int(p.Proc.P_pid) != pid || p.Proc.P_starttime.Sec == 0 {
		return "", fmt.Errorf("unobservable process birth")
	}
	return fmt.Sprintf("darwin:%s:%d:%d:%d", boot, pid, p.Proc.P_starttime.Sec, p.Proc.P_starttime.Usec), nil
}
func spawnDarwin(s LaunchSpec) (Child, error) {
	c := exec.Command(s.Executable.Path, s.Argv[1:]...)
	c.Dir = s.CWD
	c.Env = append(os.Environ(), "PLY_TASK_RUN_CONTEXT="+s.ContextPath)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Foreground: true, Ctty: int(os.Stdin.Fd())}
	if e := c.Start(); e != nil {
		return nil, e
	}
	pid := c.Process.Pid
	p := Process{State: "unknown", PID: &pid, ProcessGroup: &pid}
	birth, e := processBirth(pid)
	if e == nil {
		group, gerr := unix.Getpgid(pid)
		if gerr == nil && group == pid {
			p.State = "running"
			p.StartIdentity = &birth
			p.Quiescence = ptr(false)
		}
	}
	return &darwinChild{cmd: c, bound: p, birth: processBirth, group: unix.Getpgid, send: unix.Kill}, nil
}
func (c *darwinChild) Process() Process { return c.bound }
func (c *darwinChild) Forward(s os.Signal) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.waited || c.bound.StartIdentity == nil || c.bound.PID == nil || c.bound.ProcessGroup == nil {
		return fmt.Errorf("unbound signal")
	}
	birth, e := c.birth(*c.bound.PID)
	if e != nil || birth != *c.bound.StartIdentity {
		return fmt.Errorf("process identity changed")
	}
	g, e := c.group(*c.bound.PID)
	if e != nil || g != *c.bound.ProcessGroup {
		return fmt.Errorf("process group changed")
	}
	sig, ok := s.(syscall.Signal)
	if !ok || (sig != syscall.SIGINT && sig != syscall.SIGTERM && sig != syscall.SIGHUP) {
		return fmt.Errorf("unsupported signal")
	}
	return c.send(-g, unix.Signal(sig))
}
func (c *darwinChild) Wait() (Process, error) {
	// Do not reap concurrently with group signaling. An unreaped owned child
	// pins its PID, so the birth check and group signal cannot target a reused PID.
	// On an unobservable state we disable forwarding before entering Wait.
	for c.bound.StartIdentity != nil {
		info, err := unix.SysctlKinfoProc("kern.proc.pid", *c.bound.PID)
		if err != nil || info.Proc.P_stat == 5 {
			break
		} // Darwin SZOMB, sys/proc.h.
		time.Sleep(20 * time.Millisecond)
	}
	c.mu.Lock()
	c.waited = true
	e := c.cmd.Wait()
	c.mu.Unlock()
	p := c.bound
	p.State = "exited"
	if c.cmd.ProcessState == nil {
		p.State = "unknown"
		p.Quiescence = nil
		return p, e
	}
	code := c.cmd.ProcessState.ExitCode()
	if code >= 0 {
		p.ExitCode = &code
	}
	if st, ok := c.cmd.ProcessState.Sys().(syscall.WaitStatus); ok && st.Signaled() {
		p.Signal = ptr(st.Signal().String())
	}
	if p.StartIdentity == nil {
		p.State = "unknown"
		p.Quiescence = nil
	} else {
		err := unix.Kill(-*p.ProcessGroup, 0)
		if errors.Is(err, unix.ESRCH) {
			p.Quiescence = ptr(true)
		} else if err == nil {
			p.Quiescence = ptr(false)
		} else {
			p.Quiescence = nil
		}
	}
	return p, e
}
func systemRunner() ProcessRunner {
	return &InteractiveRunner{Terminal: darwinTerminal{}, Spawn: spawnDarwin, Signals: []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP}}
}
