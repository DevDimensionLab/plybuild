// The acceptance driver is a separate test executable. It composes the actual
// production command graph with a transport that cannot perform network I/O.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/cmd"
	notification "github.com/devdimensionlab/plybuild/internal/workflownotification"
)

func die() {
	fmt.Fprintln(os.Stderr, "Invalid or unsafe notification fixture configuration.")
	os.Exit(2)
}
func driverMain() {
	if len(os.Args) == 2 && os.Args[1] == "--capabilities" {
		fmt.Println(`{"kind":"PlyNotificationAcceptanceDriver@1","production_command_graph":true,"live_network":false,"features":["fake_transport","fixed_clock","post_log","fault_after_reservation","fault_receipt_commit"]}`)
		return
	}
	if len(os.Args) < 6 || os.Args[1] != "--config" || os.Args[3] != "--" {
		die()
	}
	args := os.Args[4:]
	if len(args) < 3 || args[0] != "workflow" || args[1] != "notification" || args[2] != "send" && args[2] != "show" && args[2] != "retry" {
		die()
	}
	var cfg struct {
		Kind      string `json:"kind"`
		Root      string `json:"artifacts_root"`
		Now       string `json:"now"`
		Log       string `json:"post_log"`
		Transport struct {
			Kind    string            `json:"kind"`
			Status  int               `json:"status"`
			Body    string            `json:"body"`
			Headers map[string]string `json:"headers"`
			Code    string            `json:"code"`
		} `json:"transport"`
		Fault *string `json:"fault"`
	}
	b, e := os.ReadFile(os.Args[2])
	if e != nil || json.Unmarshal(b, &cfg) != nil || cfg.Kind != "PlyNotificationAcceptanceCase@1" {
		die()
	}
	root, e := filepath.EvalSymlinks(cfg.Root)
	if e != nil || root != cfg.Root || !filepath.IsAbs(root) {
		die()
	}
	for _, p := range []string{os.Args[2], cfg.Log} {
		parent, e := filepath.EvalSymlinks(filepath.Dir(p))
		if e != nil || parent != filepath.Dir(p) || !strings.HasPrefix(p, root+string(filepath.Separator)) {
			die()
		}
	}
	if !fixturePath(root, cfg.Log) {
		die()
	}
	for i, arg := range args {
		var path string
		if (arg == "--route" || arg == "--file") && i+1 < len(args) {
			path = args[i+1]
		}
		if strings.HasPrefix(arg, "--route=") {
			path = strings.TrimPrefix(arg, "--route=")
		}
		if strings.HasPrefix(arg, "--file=") {
			path = strings.TrimPrefix(arg, "--file=")
		}
		if path != "" {
			absolute, err := filepath.Abs(path)
			if err != nil || !fixturePath(root, absolute) {
				die()
			}
			data, err := os.ReadFile(absolute)
			if err == nil {
				var v any
				if json.Unmarshal(data, &v) == nil && !fixtureLocators(root, v) {
					die()
				}
			}
		}
	}
	now, e := time.Parse(time.RFC3339, cfg.Now)
	if e != nil {
		die()
	}
	switch cfg.Transport.Kind {
	case "http", "not_sent", "unknown":
	default:
		die()
	}
	if cfg.Fault != nil && *cfg.Fault != "after_reservation" && *cfg.Fault != "receipt_commit" {
		die()
	}
	d := notification.Dependencies{Now: func() time.Time { return now }, LookupEnv: os.LookupEnv}
	d.Transport = func(_ context.Context, _ string, payload []byte) notification.Observation {
		dispatch := "started"
		if cfg.Transport.Kind == "not_sent" {
			dispatch = "not_started"
		}
		entry := struct {
			Kind     string          `json:"kind"`
			Payload  json.RawMessage `json:"payload"`
			Dispatch string          `json:"dispatch"`
		}{"PlyNotificationFakePost@1", payload, dispatch}
		data, _ := json.Marshal(entry)
		data = append(data, '\n')
		f, e := os.OpenFile(cfg.Log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
		if e != nil {
			die()
		}
		n, e := f.Write(data)
		if e != nil || n != len(data) || f.Sync() != nil || f.Close() != nil {
			die()
		}
		o := notification.Observation{Status: cfg.Transport.Status, Body: []byte(cfg.Transport.Body), Dispatch: dispatch, Failed: cfg.Transport.Kind != "http"}
		if s, ok := cfg.Transport.Headers["Retry-After"]; ok {
			o.RetryAfter = []string{s}
		}
		return o
	}
	d.Fault = func(point string) error {
		if cfg.Fault != nil && point == *cfg.Fault {
			if point == "after_reservation" {
				os.Exit(86)
			}
			return errors.New("injected durable receipt failure")
		}
		return nil
	}
	workflow, _, e := cmd.RootCmd.Find([]string{"workflow"})
	if e != nil {
		die()
	}
	for _, child := range workflow.Commands() {
		if child.Name() == "notification" {
			workflow.RemoveCommand(child)
		}
	}
	workflow.AddCommand(cmd.NewWorkflowNotificationCommand(d))
	cmd.RootCmd.SetArgs(os.Args[4:])
	if e = cmd.ExecuteE(); e != nil {
		os.Exit(cmd.ExitCode(e))
	}
}

func TestMain(m *testing.M) {
	if len(os.Args) > 1 && (os.Args[1] == "--capabilities" || os.Args[1] == "--config") {
		driverMain()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestDriverRefusesCommandsOutsideNotification(t *testing.T) {
	if os.Getenv("PLY_DRIVER_SCOPE_CHILD") == "1" {
		driverMain()
		return
	}
	root, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	config := filepath.Join(root, "case.json")
	b, _ := json.Marshal(map[string]any{"kind": "PlyNotificationAcceptanceCase@1", "artifacts_root": root, "now": "2026-10-03T19:00:00Z", "post_log": filepath.Join(root, "posts"), "transport": map[string]any{"kind": "http", "status": 200, "body": "ok"}, "fault": nil})
	os.WriteFile(config, b, 0600)
	command := exec.Command(os.Args[0], "--config", config, "--", "workspace", "init")
	command.Dir = root
	command.Env = append(os.Environ(), "PLY_DRIVER_SCOPE_CHILD=1")
	output, e := command.CombinedOutput()
	if e == nil {
		t.Fatalf("non-notification command was executed: %s", output)
	}
	if _, e = os.Stat(filepath.Join(root, ".ply")); !os.IsNotExist(e) {
		t.Fatal("driver created unrelated workspace state")
	}
}

func fixturePath(root, path string) bool {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path || !strings.HasPrefix(path, root+string(filepath.Separator)) {
		return false
	}
	parent, e := filepath.EvalSymlinks(filepath.Dir(path))
	if e != nil || parent != filepath.Dir(path) {
		return false
	}
	info, e := os.Lstat(path)
	return os.IsNotExist(e) || e == nil && info.Mode()&os.ModeSymlink == 0
}
func fixtureLocators(root string, value any) bool {
	switch v := value.(type) {
	case map[string]any:
		for key, child := range v {
			if key == "path" || key == "worktree" || key == "state_root" {
				if path, ok := child.(string); ok && !fixturePath(root, path) {
					return false
				}
			}
			if !fixtureLocators(root, child) {
				return false
			}
		}
	case []any:
		for _, child := range v {
			if !fixtureLocators(root, child) {
				return false
			}
		}
	}
	return true
}
