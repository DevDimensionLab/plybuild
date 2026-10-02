//go:build darwin

package taskrun

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

type execRunner struct{}

func systemExecRunner() ProcessRunner { return execRunner{} }
func (execRunner) Check() error       { return nil }

type streamCopy struct {
	truncated bool
	complete  bool
	err       error
}

func copyProviderStream(dst *os.File, src *os.File, done chan<- streamCopy) {
	var size int64
	result := streamCopy{complete: true}
	buf := make([]byte, 32<<10)
	for {
		n, e := src.Read(buf)
		if n > 0 {
			keep := n
			if remain := int64(providerStreamLimit) - size; int64(keep) > remain {
				keep = int(remain)
				result.truncated = true
			}
			if keep > 0 {
				if _, we := dst.Write(buf[:keep]); we != nil {
					result.err = we
					result.complete = false
				}
				size += int64(keep)
			}
		}
		if e != nil {
			if e != io.EOF {
				result.err = e
				result.complete = false
			}
			break
		}
	}
	if e := dst.Sync(); e != nil {
		result.err = e
		result.complete = false
	}
	if e := dst.Close(); e != nil {
		result.err = e
		result.complete = false
	}
	_ = src.Close()
	done <- result
}
func (execRunner) Run(s LaunchSpec, started func(Process) error) (Process, error) {
	p := Process{State: "not_started", Quiescence: ptr(true)}
	if e := verifyExecutable(s.Executable); e != nil {
		return p, e
	}
	if s.Timeout <= 0 || s.Timeout > 240*time.Second || s.Completion == nil {
		return p, invalid("exec timeout or evidence destination missing")
	}
	if e := privateDir(s.StreamsRoot); e != nil {
		return p, e
	}
	out, e := os.OpenFile(filepath.Join(s.StreamsRoot, "stdout.jsonl"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		return p, e
	}
	errout, e := os.OpenFile(filepath.Join(s.StreamsRoot, "stderr.txt"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		out.Close()
		return p, e
	}
	stdout, wout, e := os.Pipe()
	if e != nil {
		out.Close()
		errout.Close()
		return p, e
	}
	stderr, werr, e := os.Pipe()
	if e != nil {
		stdout.Close()
		wout.Close()
		out.Close()
		errout.Close()
		return p, e
	}
	c := exec.Command(s.Executable.Path, s.Argv[1:]...)
	c.Dir = s.CWD
	c.Env = []string{}
	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, "PLY_TASK_RUN_CONTEXT=") {
			c.Env = append(c.Env, v)
		}
	}
	c.Env = append(c.Env, "PLY_TASK_RUN_CONTEXT="+s.ContextPath)
	c.Stdout = wout
	c.Stderr = werr
	c.Stdin = nil
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if e = c.Start(); e != nil {
		stdout.Close()
		stderr.Close()
		wout.Close()
		werr.Close()
		out.Close()
		errout.Close()
		return p, e
	}
	_ = wout.Close()
	_ = werr.Close()
	outDone := make(chan streamCopy, 1)
	errDone := make(chan streamCopy, 1)
	go copyProviderStream(out, stdout, outDone)
	go copyProviderStream(errout, stderr, errDone)
	pid := c.Process.Pid
	// Successful Start returns an owned, unreaped child. Its PID cannot be reused
	// until this sole waiter reaps it. The nonce identifies this concrete launch,
	// rather than claiming an OS-wide birth observation or looking up other tasks.
	nonce := make([]byte, 16)
	_, randomErr := rand.Read(nonce)
	launchIdentity := fmt.Sprintf("owned-child:%d:%d:%s", os.Getpid(), pid, hex.EncodeToString(nonce))
	group, groupErr := unix.Getpgid(pid)
	p = Process{State: "running", PID: &pid, ProcessGroup: &pid, StartIdentity: &launchIdentity, Quiescence: ptr(false)}
	var observationErr error
	if randomErr != nil || groupErr != nil || group != pid {
		observationErr = failure("task_run_identity_unknown", 5, "Spawned child group could not be bound.")
	}
	if observationErr == nil {
		observationErr = started(p)
	}
	deadline := time.Now().Add(s.Timeout)
	timedOut := false
	termAt := time.Time{}
	killed := false
	var status unix.WaitStatus
	reaped := false
	// This goroutine is the only waiter and signal owner. Wait4(WNOHANG) and every
	// possible signal are serialized. No signal is sent after reaping, even if a
	// descendant remains: that becomes unknown and leaves the environment intact.
	for {
		observed, err := unix.Wait4(pid, &status, unix.WNOHANG, nil)
		if err == unix.EINTR {
			continue
		}
		if err != nil {
			observationErr = failure("task_run_process_unknown", 5, "Owned child wait failed: "+err.Error())
			break
		}
		if observed == pid {
			reaped = true
			break
		}
		if time.Now().After(deadline) || observationErr != nil {
			timedOut = true
			pg, gerr := unix.Getpgid(pid)
			if gerr != nil || pg != pid {
				observationErr = failure("task_run_identity_unknown", 5, "Child group changed before timeout signal; preserve environment.")
				break
			}
			if termAt.IsZero() {
				if err = unix.Kill(-pid, unix.SIGTERM); err != nil {
					observationErr = err
					break
				}
				termAt = time.Now()
			} else if time.Since(termAt) >= 5*time.Second && !killed {
				if err = unix.Kill(-pid, unix.SIGKILL); err != nil {
					observationErr = err
					break
				}
				killed = true
			} else if killed && time.Since(termAt) > 6*time.Second {
				observationErr = failure("task_run_process_unknown", 5, "Child not reaped after timeout signals; preserve environment.")
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	var waitErr error
	if reaped {
		_ = c.Process.Release()
		p.State = "exited"
		if status.Exited() {
			p.ExitCode = ptr(status.ExitStatus())
			if status.ExitStatus() != 0 {
				waitErr = fmt.Errorf("provider exit %d", status.ExitStatus())
			}
		}
		if status.Signaled() {
			p.Signal = ptr(status.Signal().String())
			waitErr = fmt.Errorf("provider signal %s", status.Signal())
		}
		err := unix.Kill(-pid, 0)
		if errors.Is(err, unix.ESRCH) {
			p.Quiescence = ptr(true)
		} else if err == nil {
			p.Quiescence = ptr(false)
		} else {
			p.Quiescence = nil
		}
	} else {
		p.State = "unknown"
		p.Quiescence = nil
	}
	if observationErr != nil {
		p.State = "unknown"
		p.Quiescence = nil
	}
	getStream := func(ch <-chan streamCopy, rd *os.File) streamCopy {
		select {
		case x := <-ch:
			return x
		case <-time.After(250 * time.Millisecond):
			_ = rd.Close()
			x := <-ch
			x.complete = false
			return x
		}
	}
	outs := getStream(outDone, stdout)
	errs := getStream(errDone, stderr)
	raw, oe := readFile(filepath.Join(s.StreamsRoot, "stdout.jsonl"), providerStreamLimit, true)
	eraw, ee := readFile(filepath.Join(s.StreamsRoot, "stderr.txt"), providerStreamLimit, true)
	seq := parseProviderEvents(raw)
	reason := seq.Reason
	if timedOut {
		reason = "timeout or interrupted observation"
	}
	if observationErr != nil {
		reason = observationErr.Error()
	}
	*s.Completion = ProviderCompletion{Kind: "ProviderCompletion@1", Stdout: FileBinding{filepath.Join(s.StreamsRoot, "stdout.jsonl"), hash(raw)}, Stderr: FileBinding{filepath.Join(s.StreamsRoot, "stderr.txt"), hash(eraw)}, StdoutTruncated: outs.truncated, StderrTruncated: errs.truncated, ThreadID: seq.ThreadID, TurnsStarted: seq.Started, TurnsCompleted: seq.Completed, StreamComplete: outs.complete && errs.complete && oe == nil && ee == nil, SequenceValid: seq.Valid, TimedOut: timedOut, Process: p, TokenUsage: seq.Usage, Reason: reason}
	if observationErr != nil {
		return p, observationErr
	}
	return p, waitErr
}
