//go:build !windows

package taskrun

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/sys/unix"
)

func TestWorkerGitKeepsOwnershipAfterControllerIsKilled(t *testing.T) {
	if dir := os.Getenv("PLY_WORKER_GIT_LOCK_CHILD"); dir != "" {
		s := workflowState{Request: WorkflowRequest{WorkspaceRoot: dir}, Result: WorkflowRun{RunID: "git-test", Paths: WorkflowPaths{RunRoot: dir}, Worker: &WorkerStatus{Generation: 1}}}
		s.Observed.Target.WorktreeLocator = dir
		lock, err := workerLock(filepath.Join(workerDirectory(s), "owner.lock"), true)
		if err != nil {
			t.Fatal(err)
		}
		defer lock.Close()
		d := Dependencies{WorkerOwner: &workerOwner{root: dir, runID: "git-test", generation: 1, lock: lock}}
		if _, err = workerGit(d, s, "-c", "alias.ply-worker-lock-test=!./git-effect.sh", "ply-worker-lock-test"); err != nil {
			t.Fatal(err)
		}
		return
	}
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err = privateDir(filepath.Join(dir, "worker")); err != nil {
		t.Fatal(err)
	}
	git := exec.Command("git", "init", "--quiet", dir)
	if b, e := git.CombinedOutput(); e != nil {
		t.Fatalf("git init: %s %v", b, e)
	}
	script := "#!/bin/sh\nset -eu\nprintf ready >git.ready\nwhile [ ! -f git.release ]; do sleep 0.02; done\ngit config --local worker.testEffect completed\nprintf done >git.done\n"
	if err = os.WriteFile(filepath.Join(dir, "git-effect.sh"), []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	c := exec.Command(os.Args[0], "-test.run=^TestWorkerGitKeepsOwnershipAfterControllerIsKilled$", "-test.count=1")
	c.Env = append(os.Environ(), "PLY_WORKER_GIT_LOCK_CHILD="+dir)
	if err = c.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.WriteFile(filepath.Join(dir, "git.release"), nil, 0600)
		c.Process.Kill()
		c.Wait()
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			if _, e := os.Stat(filepath.Join(dir, "git.done")); e == nil {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err = os.Stat(filepath.Join(dir, "git.ready")); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("trusted Git effect did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err = c.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	c.Wait()
	lock, err := workerLock(filepath.Join(dir, "worker", "owner.lock"), false)
	if lock != nil {
		lock.Close()
	}
	if !errors.Is(err, unix.EWOULDBLOCK) {
		t.Fatalf("controller death released ownership while trusted Git was still writing: %v", err)
	}
	if err = os.WriteFile(filepath.Join(dir, "git.release"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	deadline = time.Now().Add(5 * time.Second)
	for {
		lock, err = workerLock(filepath.Join(dir, "worker", "owner.lock"), false)
		if err == nil {
			lock.Close()
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("Git never released ownership")
		}
		time.Sleep(10 * time.Millisecond)
	}
	c = exec.Command("git", "-C", dir, "config", "--local", "worker.testEffect")
	if b, e := c.Output(); e != nil || strings.TrimSpace(string(b)) != "completed" {
		t.Fatalf("actual Git effect missing: %s %v", b, e)
	}
}

func TestWorkerAttemptWithUnknownProcessCannotBeConsumed(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	s := workflowState{Result: WorkflowRun{RunID: "test", Paths: WorkflowPaths{RunRoot: dir}, Worker: &WorkerStatus{}}}
	d := Dependencies{}
	jobDir := filepath.Join(workerDirectory(s), "attempts", "developer-0001")
	if err = privateDir(jobDir); err != nil {
		t.Fatal(err)
	}
	in, err := workflowKeep(d, filepath.Join(jobDir, "intent.json"), workerAttemptIntent{RunID: "test", JobID: "developer-0001", Kind: "developer", Argv: []string{}, Environment: []string{}})
	if err != nil {
		t.Fatal(err)
	}
	s.Result.Worker.Job = &WorkerJob{ID: "developer-0001", Kind: "developer", Directory: jobDir, IntentHash: in.SHA256}
	stream := filepath.Join(jobDir, "stdout.txt")
	if err = d.writeOnce(stream, nil); err != nil {
		t.Fatal(err)
	}
	r := workerAttemptResult{RunID: "test", JobID: "developer-0001", IntentSHA256: in.SHA256, FinishedAt: "2026-10-09T00:00:01Z", Exit: ptr(0), Stdout: FileBinding{stream, hash(nil)}, Stderr: FileBinding{stream, hash(nil)}}
	for _, process := range []*Process{nil, {State: "exited", ExitCode: ptr(0)}, {State: "running", Quiescence: ptr(true)}, {State: "exited", ExitCode: ptr(0), Quiescence: ptr(false)}} {
		r.Process = process
		if err = d.replaceValue(filepath.Join(jobDir, "result.json"), r); err != nil {
			t.Fatal(err)
		}
		if _, err = workerRunReserved(d, s); err == nil {
			t.Fatalf("accepted unknown/live process: %+v", process)
		}
	}
	r.Process = &Process{State: "exited", ExitCode: ptr(0), Quiescence: ptr(true)}
	if err = d.replaceValue(filepath.Join(jobDir, "result.json"), r); err != nil {
		t.Fatal(err)
	}
	if _, err = workerRunReserved(d, s); err != nil {
		t.Fatal(err)
	}
}

// Opt-in evidence against the installed CLI, without invoking a model or MCP server.
func TestWorkerActualProviderToolPolicy(t *testing.T) {
	provider := os.Getenv("PLY_WORKER_TEST_CODEX")
	if provider == "" {
		t.Skip("set PLY_WORKER_TEST_CODEX for the installed CLI configuration check")
	}
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	s := workflowState{}
	s.Observed.Target.WorktreeLocator = dir
	provider, err = filepath.EvalSymlinks(provider)
	if err != nil {
		t.Fatal(err)
	}
	rawProvider, err := os.ReadFile(provider)
	if err != nil {
		t.Fatal(err)
	}
	s.Request.Runtime.Executable = Executable{Path: provider, SHA256: hash(rawProvider)}
	prefix := []string{provider, "--no-daemon", "-a", "never"}
	prefix = append(prefix, workerCodexPolicy(dir, filepath.Join(dir, "parent.git"), filepath.Join(dir, "output"), filepath.Join(dir, "tmp"), false)...)
	extra, _, err := workerToolPolicy(s, prefix, true)
	if err != nil {
		t.Fatal(err)
	}
	prefix = append(prefix, extra...)
	_, proof, err := workerToolPolicy(s, prefix, false)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(proof)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(raw))
}

func TestWorkerDisabledMCPOverridesRemainValidWithoutUserConfig(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	provider := filepath.Join(dir, "configuration-only-provider")
	script := "#!/bin/sh\ncat <<'JSON'\n[{\"name\":\"fixture_stdio\",\"enabled\":true,\"transport\":{\"type\":\"stdio\"}},{\"name\":\"fixture_http\",\"enabled\":true,\"transport\":{\"type\":\"streamable_http\"}}]\nJSON\n"
	if err := os.WriteFile(provider, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	s := workflowState{}
	s.Observed.Target.WorktreeLocator = dir
	s.Request.Runtime.Executable = Executable{Path: provider, SHA256: hash([]byte(script))}
	args, _, err := workerToolPolicy(s, []string{provider}, true)
	if err != nil {
		t.Fatal(err)
	}
	var config struct {
		Servers map[string]struct {
			Enabled bool
			Command string
			URL     string
		} `toml:"mcp_servers"`
	}
	lines := []string{}
	for i := 1; i < len(args); i += 2 {
		lines = append(lines, args[i])
	}
	if err = toml.Unmarshal([]byte(strings.Join(lines, "\n")), &config); err != nil {
		t.Fatal(err)
	}
	if config.Servers["fixture_stdio"].Command != "/usr/bin/false" || config.Servers["fixture_http"].URL != "https://ply-worker-disabled.invalid/" {
		t.Fatalf("exec --ignore-user-config receives transportless MCP entries: %+v", config.Servers)
	}
	for name, server := range config.Servers {
		if server.Enabled {
			t.Fatalf("enabled MCP: %s", name)
		}
	}
}

func TestWorkerPolicySharedScratchRemainsValidTOML(t *testing.T) {
	for _, review := range []bool{false, true} {
		args := workerCodexPolicy("/fixture/task", "/fixture/main/.git", "/scratch", "/scratch", review)
		var policy struct {
			Permissions map[string]struct{ Filesystem map[string]string }
		}
		if err := toml.Unmarshal([]byte(args[3]), &policy); err != nil {
			t.Fatalf("actual Codex -c value is invalid: %v", err)
		}
		fs := policy.Permissions["ply-worker"].Filesystem
		if fs["/scratch"] != "write" || fs["/fixture/main/.git"] != "read" || fs["/fixture/task/.git"] != "read" {
			t.Fatalf("write bounds lost: %#v", fs)
		}
		if review && fs["/fixture/task"] == "write" {
			t.Fatal("reviewer acquired Task writes")
		}
	}
}

func TestWorkerCandidateWithOrdinaryExecutableModeBecomesPrivate(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "candidate")
	raw := []byte("#!/bin/sh\nexit 0\n")
	if err = os.WriteFile(path, raw, 0755); err != nil {
		t.Fatal(err)
	}
	if err = os.Chmod(path, 0755); err != nil {
		t.Fatal(err)
	}
	if err = workerSealCandidate(path); err != nil {
		t.Fatal(err)
	}
	binding, data, err := deliveryReadBinding(path, 64<<20)
	if err != nil || binding.SHA256 != hash(raw) || string(data) != string(raw) {
		t.Fatalf("candidate bytes changed or not private: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0700 {
		t.Fatalf("candidate is not executable and private: %v %v", info, err)
	}
	link := filepath.Join(dir, "link")
	if err = os.Symlink(path, link); err != nil {
		t.Fatal(err)
	}
	if err = workerSealCandidate(link); err == nil {
		t.Fatal("candidate symlink was accepted")
	}
}

func TestWorkerThreadIdentityMustBeActualAndUnique(t *testing.T) {
	for _, bad := range []string{"", `{"type":"item.completed","thread_id":"invented"}`, "not json", "{\"type\":\"thread.started\",\"thread_id\":\"one\"}\n", "{\"type\":\"thread.started\",\"thread_id\":\"one\"}\n{\"type\":\"turn.started\"}\n{\"type\":\"turn.failed\",\"error\":{\"message\":\"failed\"}}\n", "{\"type\":\"thread.started\",\"thread_id\":\"one\"}\n{\"type\":\"thread.started\",\"thread_id\":\"two\"}"} {
		if _, err := workerThreadID([]byte(bad)); err == nil {
			t.Fatalf("accepted %q", bad)
		}
	}
	id, err := workerThreadID([]byte("{\"type\":\"thread.started\",\"thread_id\":\"actual-thread\"}\n{\"type\":\"turn.started\"}\n{\"type\":\"turn.completed\",\"usage\":{\"input_tokens\":1,\"cached_input_tokens\":0,\"output_tokens\":1}}\n"))
	if err != nil || id != "actual-thread" {
		t.Fatalf("lost actual identity: %s %v", id, err)
	}
}

func TestWorkerInheritedLockSurvivesControllerClose(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "owner.lock")
	lock, err := workerLock(path, true)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("/bin/sh", "-c", "printf ready; read answer")
	cmd.ExtraFiles = []*os.File{lock}
	in, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { in.Close(); cmd.Wait() }()
	b := make([]byte, 5)
	if _, err = out.Read(b); err != nil || string(b) != "ready" {
		t.Fatal("child did not inherit")
	}
	lock.Close()
	other, err := workerLock(path, false)
	if other != nil {
		other.Close()
	}
	if !errors.Is(err, unix.EWOULDBLOCK) {
		t.Fatalf("controller exit released a live attempt: %v", err)
	}
	in.Write([]byte("finish\n"))
	if err = cmd.Wait(); err != nil {
		t.Fatal(err)
	}
	other, err = workerLock(path, false)
	if err != nil {
		t.Fatal(err)
	}
	other.Close()
}

func TestWorkerOwnerRejectsOldGenerationClosedOrReplacedLock(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	s := workflowState{Request: WorkflowRequest{WorkspaceRoot: dir}, Result: WorkflowRun{RunID: "test", Paths: WorkflowPaths{RunRoot: dir}, Worker: &WorkerStatus{Generation: 2}}}
	if err := os.Mkdir(workerDirectory(s), 0700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(workerDirectory(s), "owner.lock")
	lock, err := workerLock(path, true)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Close()
	d := Dependencies{WorkerOwner: &workerOwner{root: dir, runID: "test", generation: 1, lock: lock}}
	if err = workerLiveOwner(d, s); err == nil {
		t.Fatal("old callback generation accepted")
	}
	d.WorkerOwner.generation = 2
	if err = workerLiveOwner(d, s); err != nil {
		t.Fatal(err)
	}
	if err = os.Rename(path, path+".old"); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, nil, 0600); err != nil {
		t.Fatal(err)
	}
	if err = workerLiveOwner(d, s); err == nil {
		t.Fatal("replaced lock accepted")
	}
	lock.Close()
	if err = workerLiveOwner(d, s); err == nil || !strings.Contains(err.Error(), "closed") {
		t.Fatalf("closed owner accepted: %v", err)
	}
}

func TestWorkerInheritedCapabilityRejectsSeparatelyOpenedLock(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "owner.lock")
	owner, err := workerLock(path, true)
	if err != nil {
		t.Fatal(err)
	}
	defer owner.Close()
	separate, err := physicalOpenFile(path, os.O_RDWR|unix.O_NOFOLLOW, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer separate.Close()
	if err = workerConfirmInheritedLock(owner, path); err != nil {
		t.Fatal(err)
	}
	if err = workerConfirmInheritedLock(separate, path); err == nil {
		t.Fatal("a separately opened file descriptor impersonated the held ownership capability")
	}
}
