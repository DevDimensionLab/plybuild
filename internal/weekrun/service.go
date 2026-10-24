package weekrun

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ServiceStatus struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	PlistPath     string `json:"plist_path"`
	State         string `json:"state"`
	PID           int    `json:"pid,omitempty"`
	LastExitCode  string `json:"last_exit_code,omitempty"`
	ControlPath   string `json:"control_path"`
	ControlSHA256 string `json:"control_sha256"`
	Detail        string `json:"detail,omitempty"`
}

type Launchctl func(...string) (string, error)

func runLaunchctl(args ...string) (string, error) {
	timeout := 15 * time.Second
	if len(args) > 0 && args[0] == "kickstart" {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	b, err := exec.CommandContext(ctx, "/bin/launchctl", args...).CombinedOutput()
	if err != nil {
		return string(b), fmt.Errorf("launchctl %s: %w: %s", args[0], err, strings.TrimSpace(string(b)))
	}
	return string(b), nil
}
func serviceLabel(b Binding) string { return "com.ply.week." + b.ID }
func servicePath(b Binding) string {
	return filepath.Join(b.Environment["HOME"], "Library", "LaunchAgents", serviceLabel(b)+".plist")
}
func serviceDomain() string { return "gui/" + strconv.Itoa(os.Getuid()) }
func serviceInstalled(b Binding) bool {
	st, e := os.Lstat(servicePath(b))
	return e == nil && st.Mode().IsRegular()
}
func wakeService(b Binding) error { return wakeServiceUsing(b, runLaunchctl) }
func wakeServiceUsing(b Binding, ctl Launchctl) error {
	st, err := ServiceReadback(b, ctl)
	if err != nil {
		return err
	}
	if st.State == "absent" {
		return fmt.Errorf("managed service is not installed")
	}
	// launchd can still report a PID after the old controller's final control
	// read and owner-lock release. Only a held lock proves that it will consume
	// the already-published request; otherwise wait for acquisition or real exit.
	dir, err := runDir(b.Contract.Workspace, b.ID)
	if err != nil {
		return err
	}
	until := time.Now().Add(15 * time.Second)
	for st.State == "running" {
		switch ownerState(dir) {
		case "lock_held":
			return nil
		case "unowned":
		default:
			return fmt.Errorf("service owner is unknown; pending control is preserved")
		}
		if time.Now().After(until) {
			return fmt.Errorf("service process has not acquired ownership or exited; pending control is preserved")
		}
		time.Sleep(25 * time.Millisecond)
		st, err = ServiceReadback(b, ctl)
		if err != nil {
			return err
		}
		if st.State == "absent" {
			return fmt.Errorf("managed service is not installed")
		}
	}
	target := serviceDomain() + "/" + serviceLabel(b)
	if st.State == "installed" {
		_, err = ctl("bootstrap", serviceDomain(), servicePath(b))
		return err
	}
	// No -k: an active controller must never be killed to process a control.
	_, err = ctl("kickstart", target)
	return err
}
func escape(s string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

// Plist resolves ordinary installed Ply in the narrowly selected login environment.
// Successful terminal states remain exited; a crashed controller is restarted.
func Plist(b Binding) []byte { return weekPlist(b, []string{"/usr/bin/env", "ply"}) }

func weekPlist(b Binding, command []string) []byte {
	dir, _ := runDir(b.Contract.Workspace, b.ID)
	var s strings.Builder
	s.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\"><dict>\n")
	kv := func(k, v string) { fmt.Fprintf(&s, "<key>%s</key><string>%s</string>\n", escape(k), escape(v)) }
	kv("Label", serviceLabel(b))
	s.WriteString("<key>ProgramArguments</key><array>\n")
	for _, a := range append(command, "workflow", "week", "_supervise", b.ID, "--workspace", b.Contract.Workspace) {
		fmt.Fprintf(&s, "<string>%s</string>\n", escape(a))
	}
	s.WriteString("</array>\n")
	kv("WorkingDirectory", b.Contract.ReturnWorktree)
	kv("StandardOutPath", filepath.Join(dir, "service.log"))
	kv("StandardErrorPath", filepath.Join(dir, "service.log"))
	s.WriteString("<key>RunAtLoad</key><true/><key>KeepAlive</key><dict><key>SuccessfulExit</key><false/></dict>\n<key>ThrottleInterval</key><integer>15</integer>\n<key>EnvironmentVariables</key><dict>\n")
	keys := []string{}
	for k := range b.Environment {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if b.Environment[k] != "" {
			kv(k, b.Environment[k])
		}
	}
	s.WriteString("</dict>\n</dict></plist>\n")
	return []byte(s.String())
}

func ServiceReadback(b Binding, ctl Launchctl) (ServiceStatus, error) {
	out := ServiceStatus{ID: b.ID, Label: serviceLabel(b), PlistPath: servicePath(b), State: "absent", ControlPath: b.Control.Path, ControlSHA256: b.Control.SHA256}
	if err := physical(out.PlistPath); err != nil {
		return out, err
	}
	data, err := os.ReadFile(out.PlistPath)
	if os.IsNotExist(err) {
		return out, nil
	}
	if err != nil {
		return out, err
	}
	// Compare every non-runtime setting. Older services serialized the historical
	// Ply path; reading their unchanged configuration does not require that file.
	if !bytes.Equal(data, Plist(b)) && !bytes.Equal(data, weekPlist(b, []string{b.Control.Path})) {
		return out, fmt.Errorf("installed launchd service differs from the frozen week binding")
	}
	out.State = "installed"
	text, err := ctl("print", serviceDomain()+"/"+out.Label)
	if err != nil {
		out.Detail = "Installed but not loaded in the user login domain"
		return out, nil
	}
	out.State = "loaded"
	if p := regexp.MustCompile(`(?m)^\s*pid = ([0-9]+)`).FindStringSubmatch(text); len(p) == 2 {
		out.PID, _ = strconv.Atoi(p[1])
		out.State = "running"
	}
	if p := regexp.MustCompile(`(?m)^\s*last exit code = ([^\r\n]+)`).FindStringSubmatch(text); len(p) == 2 {
		out.LastExitCode = strings.TrimSpace(p[1])
	}
	return out, nil
}

func (m Manager) Service(root, id, action string) (ServiceStatus, error) {
	if runtime.GOOS != "darwin" {
		return ServiceStatus{}, fmt.Errorf("named week service requires macOS launchd")
	}
	ctl := m.Launchctl
	if ctl == nil {
		ctl = runLaunchctl
	}
	b, dir, err := readBinding(root, id)
	if err != nil {
		return ServiceStatus{}, err
	}
	path := servicePath(b)
	target := serviceDomain() + "/" + serviceLabel(b)
	if action == "status" {
		return ServiceReadback(b, ctl)
	}
	switch action {
	case "install":
		if err = directory(filepath.Dir(path)); err != nil {
			return ServiceStatus{}, err
		}
		if err = publishBytes(path, Plist(b), 0600, false); err != nil {
			return ServiceStatus{}, err
		}
		if _, err = ctl("print", target); err != nil {
			_, err = ctl("bootstrap", serviceDomain(), path)
		}
	case "start":
		var st ServiceStatus
		st, err = ServiceReadback(b, ctl)
		if err == nil && st.State == "absent" {
			err = fmt.Errorf("install this named service first")
		}
		if err == nil {
			err = wakeServiceUsing(b, ctl)
		}
	case "stop", "uninstall":
		st, e := ServiceReadback(b, ctl)
		if e != nil {
			return st, e
		}
		if st.State == "absent" {
			return st, nil
		}
		// Publication of the worker's stop request must precede unloading launchd.
		// Its native supervisor independently waits for the current attempt boundary.
		parked, parkErr := m.stopIdleTerminal(b, dir)
		if parkErr != nil {
			return st, parkErr
		}
		if !parked {
			_, err = m.Control(root, id, "stop")
		}
		if err == nil && !parked {
			until := time.Now().Add(25 * time.Second)
			for {
				var s State
				s, err = m.Show(root, id)
				if err != nil {
					break
				}
				if s.DesiredAction == "stop" && (s.Phase == "stopping" || s.Phase == "stopped" || s.Phase == "expired" || s.Phase == "completed" || s.Phase == "completed_with_blocks") {
					break
				}
				if time.Now().After(until) {
					err = fmt.Errorf("stop is preserved but the controller has not confirmed worker control; service remains loaded")
					break
				}
				time.Sleep(200 * time.Millisecond)
			}
		}
		if err == nil {
			if _, e := ctl("print", target); e == nil {
				_, err = ctl("bootout", target)
			}
		}
		if err == nil && action == "uninstall" {
			var st ServiceStatus
			st, err = ServiceReadback(b, ctl)
			if err == nil && st.State != "absent" {
				err = os.Remove(path)
			}
		}
		if err == nil {
			_ = os.Remove(filepath.Join(dir, "handover.json"))
		}
	default:
		return ServiceStatus{}, fmt.Errorf("service action must be install, start, status, stop or uninstall")
	}
	out, readErr := ServiceReadback(b, ctl)
	if err != nil {
		return out, err
	}
	return out, readErr
}

// A terminal scheduler without an active native Run needs no launchd restart
// merely to stop its service. Hold its actual owner lock while publishing stop;
// preserve completed entry/result truth and never guess from a cached PID.
func (m Manager) stopIdleTerminal(b Binding, dir string) (bool, error) {
	owner, err := lock(filepath.Join(dir, "owner.lock"))
	if errors.Is(err, errOwned) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer owner.Close()
	s, err := readState(b, dir)
	if err != nil {
		return false, err
	}
	if s.ActiveTaskID != "" || !terminal(s) {
		return false, nil
	}
	guard, err := lock(filepath.Join(dir, "control.lock"))
	if err != nil {
		return false, err
	}
	defer guard.Close()
	now := stamp(m.now())
	ctrl := Control{ID: hash(canonical([]string{"stop", now})), Action: "stop", RequestedAt: now}
	if err = publish(filepath.Join(dir, "control.json"), ctrl, true); err != nil {
		return false, err
	}
	s.SeenControl = ctrl.ID
	s.DesiredAction = "stop"
	s.ResumeRequested = false
	// Completed native work remains completed when only its service is removed.
	if s.Phase != "completed" && s.Phase != "completed_with_blocks" && s.Phase != "expired" {
		s.Phase = "stopped"
	}
	s.NextAction = "Service stop preserved; no active Run, completed work and original deadline retained"
	if err = m.save(b, dir, &s); err != nil {
		return false, err
	}
	return true, nil
}
