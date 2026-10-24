//go:build !windows

package weekrun

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"syscall"
	"time"
)

func launchDetached(b Binding, current ...string) error {
	control, err := os.Executable()
	if len(current) > 0 {
		control, err = current[0], nil
	}
	if err != nil {
		return err
	}
	dir, _ := runDir(b.Contract.Workspace, b.ID)
	// Loaded launchd owns restart supervision. Kickstart does not kill an active
	// service, unlike -k; an already-running service will observe its control file.
	if serviceInstalled(b) {
		data, readErr := os.ReadFile(servicePath(b))
		if readErr != nil {
			return readErr
		}
		if bytes.Equal(data, Plist(b)) {
			return wakeService(b)
		}
		if !bytes.Equal(data, weekPlist(b, []string{b.Control.Path})) {
			return fmt.Errorf("installed launchd service differs from the frozen week binding")
		}
		// An old launchd entry may name removed Ply. The ordinary CLI can
		// run the same controller under its existing locks without rewriting it.
	}
	guard, err := lock(filepath.Join(dir, "launch.lock"))
	if errors.Is(err, errOwned) {
		return nil
	}
	if err != nil {
		return err
	}
	defer guard.Close()
	if ownerState(dir) == "lock_held" {
		return nil
	}
	log, err := os.OpenFile(filepath.Join(dir, "controller.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer log.Close()
	cmd := exec.Command(control, "workflow", "week", "_serve", b.ID, "--workspace", b.Contract.Workspace)
	cmd.Dir = b.Contract.ReturnWorktree
	keys := []string{}
	for k := range b.Environment {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	cmd.Env = []string{}
	for _, k := range keys {
		if b.Environment[k] != "" {
			cmd.Env = append(cmd.Env, k+"="+b.Environment[k])
		}
	}
	cmd.Stdout = log
	cmd.Stderr = log
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err = cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}

// Supervise is the launchd entry. If a detached controller already exists,
// wait for its boundary handover; no second scheduler or provider can acquire
// the same owner lock. Once this service owns the controller, a crash is a
// nonzero launchd exit and the same frozen control will inspect native effects.
func (m Manager) Supervise(ctx context.Context, root, id string) error {
	b, dir, err := readBinding(root, id)
	if err != nil {
		return err
	}
	for {
		if err = ctx.Err(); err != nil {
			return err
		}
		err = m.serve(ctx, root, id, true)
		if !errors.Is(err, errOwned) {
			return err
		}
		// A detached controller observes this request between native API effects.
		if err = publish(filepath.Join(dir, "handover.json"), map[string]string{"service": serviceLabel(b)}, true); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
}
