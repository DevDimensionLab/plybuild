//go:build !windows

package taskrun

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func workerCodexPolicy(task, gitdir, output, temporary string, readOnly bool) []string {
	entries := []string{`":root"="read"`, `":minimal"="read"`}
	quote := func(s string) string { b, _ := json.Marshal(s); return string(b) }
	if !readOnly {
		entries = append(entries, quote(task)+`="write"`)
		if output != temporary {
			entries = append(entries, quote(output)+`="write"`)
		}
	}
	entries = append(entries, quote(temporary)+`="write"`, quote(filepath.Join(task, ".git"))+`="read"`, quote(gitdir)+`="read"`)
	profile := "{filesystem={" + strings.Join(entries, ",") + "},network={enabled=false}}"
	args := []string{"-c", `default_permissions="ply-worker"`, "-c", "permissions.ply-worker=" + profile, "-c", `web_search="disabled"`}
	for _, feature := range workerDisabledFeatures {
		args = append(args, "-c", "features."+feature+"=false")
	}
	return args
}

var workerDisabledFeatures = []string{"hooks", "multi_agent", "apps", "plugins", "remote_plugin", "browser_use", "browser_use_external", "browser_use_full_cdp_access", "computer_use", "in_app_browser", "image_generation", "skill_mcp_dependency_install", "skill_search"}

type workerMCPSetting struct {
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	Transport struct {
		Type string `json:"type"`
	} `json:"transport"`
}

// Query only configuration, never connect to an MCP server. Empty table
// overrides are merged by Codex, so every discovered server needs its own
// explicit disabled entry. Record only names/booleans, never server secrets.
func workerToolPolicy(s workflowState, prefix []string, disable bool) ([]string, any, error) {
	if len(prefix) == 0 || prefix[0] != s.Request.Runtime.Executable.Path {
		return nil, nil, fmt.Errorf("provider configuration executable differs from the frozen runtime")
	}
	if err := verifyExecutable(s.Request.Runtime.Executable); err != nil {
		return nil, nil, err
	}
	run := func(tail ...string) ([]byte, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		c := exec.CommandContext(ctx, prefix[0], append(append([]string{}, prefix[1:]...), tail...)...)
		c.Dir = s.Observed.Target.WorktreeLocator
		c.Env = workerEnvironment()
		return c.Output()
	}
	raw, err := run("mcp", "list", "--json")
	if err != nil {
		return nil, nil, fmt.Errorf("cannot inspect configured MCP surface: %w", err)
	}
	var servers []workerMCPSetting
	if err = json.Unmarshal(raw, &servers); err != nil {
		return nil, nil, err
	}
	extra := []string{}
	for _, server := range servers {
		if !plain(server.Name, 1, 128) {
			return nil, nil, fmt.Errorf("invalid MCP server name")
		}
		for _, c := range server.Name {
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '-') {
				return nil, nil, fmt.Errorf("cannot safely disable MCP server name %q", server.Name)
			}
		}
		if disable {
			extra = append(extra, "-c", "mcp_servers."+server.Name+".enabled=false")
			// exec ignores user config. Keep an inert but complete transport in
			// the override, so user-only names remain parseable without copying
			// credentials, environment, headers, arguments or executable paths.
			switch server.Transport.Type {
			case "stdio":
				extra = append(extra, "-c", "mcp_servers."+server.Name+`.command="/usr/bin/false"`)
			case "streamable_http":
				extra = append(extra, "-c", "mcp_servers."+server.Name+`.url="https://ply-worker-disabled.invalid/"`)
			default:
				return nil, nil, fmt.Errorf("cannot safely disable unknown MCP transport for %s", server.Name)
			}
		} else if server.Enabled {
			return nil, nil, fmt.Errorf("MCP server remains enabled: %s", server.Name)
		}
	}
	if disable {
		return extra, nil, nil
	}
	raw, err = run("features", "list")
	if err != nil {
		return nil, nil, fmt.Errorf("cannot inspect effective provider features: %w", err)
	}
	features := map[string]bool{}
	for _, line := range strings.Split(string(raw), "\n") {
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			last := parts[len(parts)-1]
			if last == "true" || last == "false" {
				features[parts[0]] = last == "true"
			}
		}
	}
	for _, name := range workerDisabledFeatures {
		enabled, found := features[name]
		if !found || enabled {
			return nil, nil, fmt.Errorf("provider feature %s is not verifiably disabled", name)
		}
	}
	return nil, map[string]any{"kind": "PlyWorkerToolPolicy@1", "mcp_servers": servers, "disabled_features": workerDisabledFeatures, "user_config": "exec --ignore-user-config; discovered MCP names have transport-complete inert overrides, independently disabled", "exec_rules": "ignored by exec --ignore-rules", "scope": "Shell/filesystem enforcement is separately probed; app, browser, computer-use, plugin and MCP surfaces are disabled"}, nil
}

func workerVerifierPolicy(d Dependencies, s workflowState, directory string) (string, error) {
	quote := func(s string) string { b, _ := json.Marshal(s); return string(b) }
	// Tests can read their candidate and dependencies but write only the native
	// attempt's scratch space. They cannot alter Task, Git, worker or evidence.
	policy := "(version 1)\n(deny default)\n(allow process*)\n(allow signal)\n(allow sysctl-read)\n(allow mach-lookup)\n(allow file-read*)\n(allow file-write* (literal \"/dev/null\"))\n"
	for _, name := range []string{"home", "tmp"} {
		path := filepath.Join(directory, name)
		if err := privateDir(path); err != nil {
			return "", err
		}
		policy += "(allow file-write* (subpath " + quote(path) + "))\n"
	}
	path := filepath.Join(directory, "verifier.sb")
	if err := d.writeOnce(path, []byte(policy)); err != nil {
		return "", err
	}
	return path, nil
}

// The probe runs the provider's actual sandbox adapter, without a model call.
// Requested flags alone never become evidence of enforced Task write bounds.
func workerProbe(d Dependencies, s workflowState) (FileBinding, error) {
	if err := verifyExecutable(s.Request.Runtime.Executable); err != nil {
		return FileBinding{}, err
	}
	dir := filepath.Join(workerDirectory(s), "policy")
	if err := privateDir(dir); err != nil {
		return FileBinding{}, err
	}
	probeDir := filepath.Join(dir, "scratch")
	if err := privateDir(probeDir); err != nil {
		return FileBinding{}, err
	}
	guard := filepath.Join(dir, "protected")
	if err := d.writeOnce(guard, []byte("protected\n")); err != nil {
		return FileBinding{}, err
	}
	gitGuard := filepath.Join(s.Observed.Target.GitCommonDir, "ply-worker-probe-"+s.Result.RunID[4:20])
	if _, err := os.Lstat(gitGuard); !os.IsNotExist(err) {
		return FileBinding{}, fmt.Errorf("worker probe Git sentinel already exists")
	}
	python, err := exec.LookPath("python3")
	if err != nil {
		return FileBinding{}, err
	}
	code := `from pathlib import Path
import socket, sys
scratch=Path(sys.argv[1]); (scratch/'legal').write_text('ok')
for name in sys.argv[2:]:
 try: Path(name).write_text('forbidden')
 except PermissionError: pass
 else: raise SystemExit('protected write was permitted: '+name)
try:
 s=socket.socket(); s.bind(('127.0.0.1',0))
except PermissionError: pass
else: raise SystemExit('network operation was permitted')
print('worker sandbox probe passed')
`
	argv := []string{s.Request.Runtime.Executable.Path, "--no-daemon", "-a", "never"}
	argv = append(argv, workerCodexPolicy(s.Observed.Target.WorktreeLocator, s.Observed.Target.GitCommonDir, probeDir, probeDir, false)...)
	argv = append(argv, "sandbox", "--", python, "-I", "-B", "-c", code, probeDir, guard, gitGuard)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir = s.Observed.Target.WorktreeLocator
	cmd.Env = workerEnvironment()
	output, runErr := cmd.CombinedOutput()
	before, readErr := os.ReadFile(guard)
	_, gitErr := os.Lstat(gitGuard)
	passed := runErr == nil && readErr == nil && string(before) == "protected\n" && os.IsNotExist(gitErr)
	reviewGuard := filepath.Join(s.Observed.Target.WorktreeLocator, ".ply-review-probe-"+s.Result.RunID[4:20])
	if _, e := os.Lstat(reviewGuard); !os.IsNotExist(e) {
		return FileBinding{}, fmt.Errorf("reviewer probe sentinel already exists")
	}
	reviewArgv := []string{s.Request.Runtime.Executable.Path, "--no-daemon", "-a", "never"}
	reviewArgv = append(reviewArgv, workerCodexPolicy(s.Observed.Target.WorktreeLocator, s.Observed.Target.GitCommonDir, probeDir, probeDir, true)...)
	reviewArgv = append(reviewArgv, "sandbox", "--", python, "-I", "-B", "-c", code, probeDir, guard, gitGuard, reviewGuard)
	reviewCmd := exec.CommandContext(ctx, reviewArgv[0], reviewArgv[1:]...)
	reviewCmd.Dir = s.Observed.Target.WorktreeLocator
	reviewCmd.Env = workerEnvironment()
	reviewOutput, reviewErr := reviewCmd.CombinedOutput()
	_, reviewGuardErr := os.Lstat(reviewGuard)
	passed = passed && reviewErr == nil && os.IsNotExist(reviewGuardErr)
	value := map[string]any{"kind": "PlyWorkerPolicyProbe@1", "argv": argv, "output": string(output), "passed": passed, "controller": "Native Ply process with the invoking user's host permissions; only delegates and test subprocesses are sandboxed", "developer": "Task and private output files writable; Git, controller and evidence read-only; tool network denied", "reviewer": "read-only candidate and evidence; private scratch only"}
	value["reviewer_argv"], value["reviewer_output"] = reviewArgv, string(reviewOutput)
	if reviewErr != nil {
		value["reviewer_error"] = reviewErr.Error()
	}
	if cmd.ProcessState != nil {
		value["exit"] = cmd.ProcessState.ExitCode()
	}
	if runErr != nil {
		value["error"] = runErr.Error()
	}
	// Each generation preserves its own actual probe, including failed probes.
	path := filepath.Join(dir, fmt.Sprintf("probe-%06d.json", s.Result.Worker.Generation))
	binding, err := workflowKeep(d, path, value)
	if err != nil {
		return binding, err
	}
	if !passed {
		return binding, fmt.Errorf("worker sandbox probe failed; inspect %s", path)
	}
	return binding, nil
}

func workerProbeCurrentGeneration(d Dependencies, s workflowState) error {
	path := filepath.Join(workerDirectory(s), "policy", fmt.Sprintf("probe-%06d.json", s.Result.Worker.Generation))
	raw, err := readFile(path, 1<<20, false)
	if os.IsNotExist(err) {
		_, err = workerProbe(d, s)
		return err
	}
	if err != nil {
		return err
	}
	var probe struct {
		Passed bool `json:"passed"`
	}
	if err = json.Unmarshal(raw, &probe); err != nil {
		return err
	}
	if !probe.Passed {
		return fmt.Errorf("preserved runtime probe failed: %s", path)
	}
	return nil
}

func workerAccept(d Dependencies, s workflowState) error {
	if s.StartSHA256 != nil {
		return nil
	}
	acceptancePath := filepath.Join(workerDirectory(s), "acceptance.json")
	if _, err := os.Lstat(acceptancePath); err == nil {
		_, err = WorkflowAccept(d, s.Request.WorkspaceRoot, s.Result.RunID, s.Result.Paths.Context, acceptancePath)
		return err
	} else if !os.IsNotExist(err) {
		return err
	}
	proof, _, err := deliveryReadBinding(filepath.Join(workerDirectory(s), "policy", fmt.Sprintf("probe-%06d.json", s.Result.Worker.Generation)), 1<<20)
	if err != nil {
		return err
	}
	tmp := filepath.Join(workerDirectory(s), "tmp")
	if err = privateDir(tmp); err != nil {
		return err
	}
	sandbox, _ := Canonical(map[string]any{"read_roots": []string{s.Request.WorkspaceRoot}, "write_roots": []string{s.Request.WorkspaceRoot}, "temp_root": tmp, "matches_contract": true})
	a := DeliveryAcceptance{Acceptance: Acceptance{Envelope: deliveryEnv("run-acceptance"), RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID,
		RuntimeClaim: RuntimeClaim{RuntimeID: ptr("ply-worker"), ProfileID: ptr("ply-worker"), EffectivePolicySHA256: &proof.SHA256, NativeSessionID: ptr(deliveryOwnerSession(s))}, Sandbox: sandbox, Acceptance: "started", Issues: json.RawMessage("[]")},
		DeliveryPermission: DeliveryPermissionAcceptance{LaunchContractSHA256: s.Request.Runtime.PermissionBinding.EffectivePolicySHA256, PermissionConfirmed: true, ActualPolicyEvidence: []Evidence{{Locator: proof.Locator, SHA256: proof.SHA256, Role: "effective_policy"}}, DeliveryAgreementSHA256: workspace.DeliveryAgreementDigest(*s.Request.Delivery.Agreement), AllowedEffects: workflowhandoff.DeliveryAllowedEffects(*s.Request.Delivery.Agreement)}}
	binding, err := workflowKeep(d, filepath.Join(workerDirectory(s), "acceptance.json"), a)
	if err != nil {
		return err
	}
	_, err = WorkflowAccept(d, s.Request.WorkspaceRoot, s.Result.RunID, s.Result.Paths.Context, binding.Locator)
	return err
}

func workerEnvironment() []string {
	env := []string{}
	for _, item := range os.Environ() {
		name, _, _ := strings.Cut(item, "=")
		switch name {
		case "HOME", "USER", "LOGNAME", "PATH", "SHELL", "TMPDIR", "LANG", "LC_ALL", "LC_CTYPE", "TERM", "CODEX_HOME", "SSL_CERT_FILE", "SSL_CERT_DIR":
		default:
			continue
		}
		env = append(env, item)
	}
	env = append(env, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null")
	sort.Strings(env)
	return env
}
