package taskrun

import (
	"os"
	"os/signal"
	"sync"
)

// Terminal and Child are the narrow injected boundary used by synthetic tests.
// No test switch or environment bypass is part of the product CLI.
type Terminal interface {
	Check() error
	Capture() (func() error, error)
	Foreground(int) error
}
type Child interface {
	Process() Process
	Wait() (Process, error)
	Forward(os.Signal) error
}
type InteractiveRunner struct {
	Terminal Terminal
	Spawn    func(LaunchSpec) (Child, error)
	Signals  []os.Signal
}

func (r *InteractiveRunner) Check() error { return r.Terminal.Check() }
func (r *InteractiveRunner) Run(spec LaunchSpec, started func(Process) error) (Process, error) {
	unknown := Process{State: "unknown"}
	if e := r.Check(); e != nil {
		return Process{State: "not_started", Quiescence: ptr(true)}, e
	}
	if e := verifyExecutable(spec.Executable); e != nil {
		return Process{State: "not_started", Quiescence: ptr(true)}, e
	}
	restore, e := r.Terminal.Capture()
	if e != nil {
		return Process{State: "not_started", Quiescence: ptr(true)}, e
	}
	child, e := r.Spawn(spec)
	if e != nil {
		_ = restore()
		return Process{State: "not_started", Quiescence: ptr(true)}, e
	}
	p := child.Process()
	var observationErr error
	if p.State == "running" && p.StartIdentity != nil && p.ProcessGroup != nil {
		if e = r.Terminal.Foreground(*p.ProcessGroup); e != nil {
			observationErr = e
		}
		if e = started(p); e != nil {
			observationErr = e
		}
	} else {
		observationErr = failure("task_run_identity_unknown", 5, "Spawn succeeded but process birth identity could not be observed.")
	}
	signals := make(chan os.Signal, 8)
	done := make(chan struct{})
	var wg sync.WaitGroup
	if len(r.Signals) > 0 {
		signal.Notify(signals, r.Signals...)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case s := <-signals:
				_ = child.Forward(s)
			case <-done:
				return
			}
		}
	}()
	exited, waitErr := child.Wait()
	signal.Stop(signals)
	close(done)
	wg.Wait()
	restoreErr := restore()
	if restoreErr != nil {
		exited.State = "unknown"
		exited.Quiescence = nil
		return exited, failure("task_run_terminal_restore_failed", 5, "Process outcome preserved; terminal restoration was unsuccessful.")
	}
	if observationErr != nil {
		unknown = exited
		unknown.State = "unknown"
		unknown.Quiescence = nil
		return unknown, observationErr
	}
	return exited, waitErr
}
