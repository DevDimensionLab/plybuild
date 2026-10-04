package taskrun

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

var completedEvents = []byte("{\"type\":\"thread.started\",\"thread_id\":\"fixture-thread\"}\n{\"type\":\"turn.started\"}\n{\"type\":\"item.started\",\"item\":{\"id\":\"cmd\",\"type\":\"command_execution\",\"status\":\"in_progress\"}}\n{\"type\":\"item.completed\",\"item\":{\"id\":\"cmd\",\"type\":\"command_execution\",\"status\":\"completed\"}}\n{\"type\":\"turn.completed\",\"usage\":{\"input_tokens\":12,\"output_tokens\":8}}\n")

func TestFactoryExecIncludesPlatformRuntimeWithoutBroadWrites(t *testing.T) {
	var r Request
	r.Runtime.PermissionBinding.ProfileID = "factory-test"
	p := FactoryPolicy{ReadRoots: []string{"/runtime"}, WriteRoots: []string{"/fixture/task"}}
	args, err := execArgv(r, p, "/fixture/task", "/fixture/instructions.md")
	if err != nil {
		t.Fatal(err)
	}
	want := `permissions.factory-test={filesystem={":root"="none",":minimal"="read","/runtime"="read","/fixture/task"="write"},network={enabled=false}}`
	for i, arg := range args {
		if arg == want && i > 0 && args[i-1] == "-c" {
			return
		}
	}
	t.Fatalf("missing scoped profile with the platform runtime baseline: %q", args)
}

func TestFactoryProfileLinkedGitMetadataStaysWithinWriteRoots(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	task, common := filepath.Join(root, "task"), filepath.Join(root, "main", ".git")
	metadata := filepath.Join(common, "worktrees", "task")
	for _, path := range []string{task, metadata} {
		if err := os.MkdirAll(path, 0700); err != nil {
			t.Fatal(err)
		}
	}
	marker := filepath.Join(task, ".git")
	p := FactoryPolicy{WriteRoots: []string{task, common}}
	for _, pointer := range []string{metadata, "../main/.git/worktrees/task"} {
		if err := os.WriteFile(marker, []byte("gitdir: "+pointer+"\n"), 0600); err != nil {
			t.Fatal(err)
		}
		roots, err := factoryProfileWriteRoots(p)
		want := []string{task, common, metadata}
		sort.Strings(want)
		if err != nil || !equal(roots, want) {
			t.Fatalf("roots=%q, err=%v; want %q", roots, err, want)
		}
	}
	if !equal(p.WriteRoots, []string{task, common}) {
		t.Fatal("profile generation mutated the bound policy")
	}
	escaped := filepath.Join(root, "main", ".git-outside", "worktrees", "task")
	if err := os.MkdirAll(escaped, 0700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(common, "worktrees", "redirect")
	if err := os.Symlink(metadata, link); err != nil {
		t.Fatal(err)
	}
	for _, content := range []string{"gitdir: " + escaped + "\n", "invalid pointer\n", "gitdir: \n", "gitdir: " + link + "\n"} {
		if err := os.WriteFile(marker, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := factoryProfileWriteRoots(p); err == nil {
			t.Fatalf("accepted Git metadata pointer %q", content)
		}
		if _, err := execArgv(Request{}, p, task, "/fixture/instructions.md"); err == nil {
			t.Fatal("exec adapter did not propagate the profile error")
		}
	}
}

func TestFactoryNativeSequence(t *testing.T) {
	for _, tc := range []struct {
		name  string
		raw   []byte
		valid bool
	}{
		{"complete", completedEvents, true},
		{"truncated", completedEvents[:len(completedEvents)-1], false},
		{"second turn", append(append([]byte{}, completedEvents...), []byte("{\"type\":\"turn.started\"}\n")...), false},
		{"late error", append(append([]byte{}, completedEvents...), []byte("{\"type\":\"error\"}\n")...), false},
		{"unfinished tool", bytes.Replace(completedEvents, []byte("{\"type\":\"item.completed\",\"item\":{\"id\":\"cmd\",\"type\":\"command_execution\",\"status\":\"completed\"}}\n"), nil, 1), false},
		{"thread mismatch", bytes.Replace(completedEvents, []byte(`"type":"turn.started"`), []byte(`"type":"turn.started","thread_id":"wrong"`), 1), false},
		{"invalid json", []byte("not JSON\n"), false},
		{"oversize", append(bytes.Repeat([]byte(" "), providerLineLimit+1), '\n'), false},
		{"failed turn", bytes.Replace(completedEvents, []byte("turn.completed"), []byte("turn.failed"), 1), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := parseProviderEvents(tc.raw)
			if got.Valid != tc.valid {
				t.Fatalf("%+v", got)
			}
		})
	}
}
func factoryFixture(t *testing.T) (Dependencies, Request, string, FactoryAuthorization) {
	t.Helper()
	// Factory authorization intentionally confines effects to /private/tmp.
	// macOS's ambient TMPDIR normally resolves under /private/var/folders, so
	// resolving t.TempDir() alone does not put this fixture in the allowed root.
	output, e := os.MkdirTemp("/private/tmp", "ply-factory-test-")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(output); err != nil {
			t.Errorf("remove private factory fixture: %v", err)
		}
	})
	iteration := filepath.Join(output, "work", "iteration-1")
	if e = os.MkdirAll(iteration, 0700); e != nil {
		t.Fatal(e)
	}
	t.Setenv("TMPDIR", iteration)
	d, r, _ := fixture(t, filepath.Join(iteration, "workspace"))
	r.SchemaVersion = 2
	r.Runtime.Mode = "exec"
	r.Runtime.PermissionBinding.ProfileID = "factory-test"
	r.HumanAuthority.StartSurface = "human_started_factory_test"
	p := dMustPreparation(t, d, r)
	temp := filepath.Join(iteration, "temp")
	if e = os.Mkdir(temp, 0700); e != nil {
		t.Fatal(e)
	}
	writes := []string{p.Plan.WorktreePath, p.Plan.Target.GitCommonDir, filepath.Join(r.WorkspaceRoot, ".ply/task-runs/v1"), filepath.Join(r.WorkspaceRoot, ".ply/workflow/handoffs/v1"), filepath.Join(r.WorkspaceRoot, ".ply/projects.lock"), filepath.Join(r.WorkspaceRoot, ".ply/work-items.lock"), temp}
	sort.Strings(writes)
	policy := FactoryPolicy{"PlyFactoryPolicy@1", []string{r.WorkspaceRoot}, writes, temp, false}
	ev := filepath.Join(output, "evidence")
	if e = os.Mkdir(ev, 0700); e != nil {
		t.Fatal(e)
	}
	policyPath := writeAny(t, ev, "policy.json", policy)
	policyHash := hashFileTest(t, policyPath)
	r.Runtime.PermissionBinding.EffectivePolicySHA256 = policyHash
	r.Runtime.PermissionBinding.Evidence = []Evidence{{policyPath, policyHash, "effective_policy"}}
	git := testGit(t)
	shell, _ := filepath.EvalSymlinks("/bin/sh")
	now := d.Now()
	a := FactoryAuthorization{Kind: "PlyFactoryTestAuthorization@1", ActorClaim: r.HumanAuthority.ActorClaim, Activity: "synthetic-factory", OutputRoot: output, WorkRoot: filepath.Join(output, "work"), TasksPerIteration: 2, Iterations: 1, MaxAgentStarts: 2, StartedUTC: now.Add(-time.Second).UTC().Format(time.RFC3339Nano), DeadlineUTC: now.Add(1199 * time.Second).UTC().Format(time.RFC3339Nano), Model: r.Runtime.Model, ScenarioSHA256: hash([]byte("synthetic scenario")), Ply: r.Runtime.PlyExecutable, Codex: r.Runtime.Executable, Git: Executable{git, hashFileTest(t, git)}, Python: Executable{shell, hashFileTest(t, shell)}, EffectivePolicySHA256: policyHash}
	policy.ReadRoots = factoryReadRoots(r.WorkspaceRoot, a)
	policyPath = writeAny(t, ev, "policy.json", policy)
	policyHash = hashFileTest(t, policyPath)
	a.EffectivePolicySHA256 = policyHash
	r.Runtime.PermissionBinding.EffectivePolicySHA256 = policyHash
	r.Runtime.PermissionBinding.Evidence = []Evidence{{policyPath, policyHash, "effective_policy"}}
	auth := writeAny(t, ev, "authorization.json", a)
	r.FactoryTest = &FactoryTest{auth, hashFileTest(t, auth), 1, 1, "low", 240}
	d.ExecRunner = &fakeRunner{}
	return d, r, writeAny(t, r.WorkspaceRoot, "request.json", r), a
}

func TestFactoryFixtureOwnsPrivateTempRoot(t *testing.T) {
	// A fresh subtest makes the fixture independent of both the parent's cached
	// testing.TempDir root and an unusable ambient temporary directory.
	unavailable := filepath.Join(t.TempDir(), "not-created")
	t.Run("ambient_temp_is_not_factory_authority", func(t *testing.T) {
		t.Setenv("TMPDIR", unavailable)
		d, _, file, a := factoryFixture(t)
		if !within("/private/tmp", a.OutputRoot) {
			t.Fatalf("factory output escaped its authorized root: %s", a.OutputRoot)
		}
		if _, err := PreviewStart(d, file); err != nil {
			t.Fatalf("valid private factory fixture rejected: %v", err)
		}
		if _, err := os.Stat(unavailable); !os.IsNotExist(err) {
			t.Fatal("factory fixture used the ambient temporary directory")
		}
	})
}

func fakeCompletion(t *testing.T, s LaunchSpec) {
	t.Helper()
	writeAny(t, s.StreamsRoot, "unused.json", map[string]string{"evidence": "synthetic provider only"})
	if e := os.WriteFile(filepath.Join(s.StreamsRoot, "stdout.jsonl"), completedEvents, 0600); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(filepath.Join(s.StreamsRoot, "stderr.txt"), nil, 0600); e != nil {
		t.Fatal(e)
	}
	seq := parseProviderEvents(completedEvents)
	*s.Completion = ProviderCompletion{Kind: "ProviderCompletion@1", Stdout: FileBinding{filepath.Join(s.StreamsRoot, "stdout.jsonl"), hash(completedEvents)}, Stderr: FileBinding{filepath.Join(s.StreamsRoot, "stderr.txt"), hash(nil)}, ThreadID: seq.ThreadID, TurnsStarted: 1, TurnsCompleted: 1, StreamComplete: true, SequenceValid: true, TokenUsage: seq.Usage}
}
func TestFactorySchemaSlotAndCompletion(t *testing.T) {
	d, r, file, a := factoryFixture(t)
	raw, _ := Canonical(r)
	if _, e := parseRequest(raw); e != nil {
		t.Fatal(e)
	}
	// Unknown fields are rejected, including on the new optional boundary.
	var m map[string]any
	_ = json.Unmarshal(raw, &m)
	m["extra"] = true
	bad, _ := json.Marshal(m)
	if _, e := parseRequest(bad); e == nil {
		t.Fatal("unknown field accepted")
	}
	r1 := r
	r1.SchemaVersion = 1
	r1.Runtime.Mode = "interactive"
	r1.FactoryTest = nil
	r1.HumanAuthority.StartSurface = "human_ordinary_terminal"
	b1, _ := Canonical(r1)
	if bytes.Contains(b1, []byte("factory_test")) {
		t.Fatal("schema 1 bytes extended")
	}
	if _, e := parseRequest(b1); e != nil {
		t.Fatal(e)
	}
	before := treeState(t, a.OutputRoot)
	preview, e := PreviewStart(d, file)
	if e != nil {
		t.Fatal(e)
	}
	if before != treeState(t, a.OutputRoot) {
		t.Fatal("preview wrote")
	}
	runner := d.ExecRunner.(*fakeRunner)
	runner.inside = func(s LaunchSpec) error {
		if !strings.Contains(strings.Join(s.Argv, " "), "--no-daemon") || !strings.Contains(strings.Join(s.Argv, " "), "--ignore-user-config") {
			t.Fatal("exec isolation flags missing")
		}
		claim := acceptance(t, d, r)
		claim.SchemaVersion = 2
		claim.RuntimeClaim.ModelID = nil
		if _, e := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "accept.json", claim)); e != nil {
			return e
		}
		if _, e := SubmitReport(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "report.json", factorySemanticReport(t, d, r))); e != nil {
			return e
		}
		fakeCompletion(t, s)
		return nil
	}
	result, e := Start(d, file, *preview.(Preview).Confirmation)
	if e != nil {
		t.Fatalf("%v; child: %v; %+v", e, runner.lastErr, result)
	}
	if result.SchemaVersion != 3 || result.Collection.State != "qualified" || result.ProviderCompletion == nil || result.TaskExecution.State != "inactive" || result.RuntimeFacts.ModelState != "unknown" {
		t.Fatalf("%+v", result)
	}
	if _, e = Start(d, file, *preview.(Preview).Confirmation); e != nil || runner.calls != 1 {
		t.Fatal("retry relaunched", e)
	}
	p, e := PreviewCollect(d, r.WorkspaceRoot, RunID(r.RequestKey))
	if e != nil || p.SchemaVersion != 3 || p.ProviderCompletion == nil {
		t.Fatal(p, e)
	}
	manual := d
	manual.TaskStatusPath = "/never-read"
	if _, e = PreviewCollect(manual, r.WorkspaceRoot, RunID(r.RequestKey)); e == nil {
		t.Fatal("manual exec status accepted")
	}
	// A new key cannot mint another start in the same durable slot.
	other := r
	other.RequestKey = "fixture/different-key"
	otherFile := writeAny(t, r.WorkspaceRoot, "other.json", other)
	if e = reserveFactory(d, other, dMustPreparation(t, d, r), otherFile); e == nil {
		t.Fatal("slot reuse accepted")
	}
}
func TestFactoryAuthEscapeAndDeadline(t *testing.T) {
	d, r, file, a := factoryFixture(t)
	p := dMustPreparation(t, d, r)
	if _, _, e := validateFactory(r, p, file, d.Now()); e != nil {
		t.Fatal(e)
	}
	if _, _, e := validateFactory(r, p, "/private/tmp/outside.json", d.Now()); e == nil {
		t.Fatal("request escape accepted")
	}
	if _, _, e := validateFactory(r, p, file, d.Now().Add(1300*time.Second)); e == nil {
		t.Fatal("deadline accepted")
	}
	a.MaxAgentStarts = 4
	writeAny(t, filepath.Dir(r.FactoryTest.AuthorizationPath), "authorization.json", a)
	r.FactoryTest.AuthorizationSHA256 = hashFileTest(t, r.FactoryTest.AuthorizationPath)
	if _, _, e := validateFactory(r, p, file, d.Now()); e == nil {
		t.Fatal("budget expansion accepted")
	}
}
func TestFactoryCompletionRequiredForInactive(t *testing.T) {
	p := Process{State: "exited", ExitCode: ptr(0), Quiescence: ptr(true)}
	c := ProviderCompletion{Kind: "ProviderCompletion@1", ThreadID: ptr("thread"), TurnsStarted: 1, TurnsCompleted: 1, StreamComplete: true, SequenceValid: true, Process: p}
	if !providerInactive(c) {
		t.Fatal("complete evidence rejected")
	}
	for _, mutate := range []func(*ProviderCompletion){func(c *ProviderCompletion) { c.TimedOut = true }, func(c *ProviderCompletion) { c.StdoutTruncated = true }, func(c *ProviderCompletion) { c.StderrTruncated = true }, func(c *ProviderCompletion) { c.Process.Quiescence = nil }, func(c *ProviderCompletion) { c.StreamComplete = false }, func(c *ProviderCompletion) { c.SequenceValid = false }, func(c *ProviderCompletion) { c.Process.Signal = ptr("term") }} {
		bad := c
		mutate(&bad)
		if providerInactive(bad) {
			t.Fatal("uncertain evidence qualified")
		}
	}
	for _, process := range []Process{
		{State: "unknown", ExitCode: ptr(0), Quiescence: ptr(true)},
		{State: "running", ExitCode: ptr(0), Quiescence: ptr(true)},
		{State: "exited", ExitCode: ptr(1), Quiescence: ptr(true)},
		{State: "exited", Quiescence: ptr(true)},
		{State: "exited", ExitCode: ptr(0), Quiescence: ptr(false)},
	} {
		bad := c
		bad.Process = process
		if providerInactive(bad) {
			t.Fatal("process uncertainty qualified", process)
		}
	}
}
func TestFactoryFreshBinaryCallbacks(t *testing.T) {
	bin := os.Getenv("PLY_TASK_RUN_TEST_BINARY")
	if bin == "" {
		t.Skip("make compat-cli supplies freshly built binary")
	}
	d, r, file, a := factoryFixture(t)
	r.Runtime.PlyExecutable = Executable{bin, hashFileTest(t, bin)}
	a.Ply = r.Runtime.PlyExecutable
	var policy FactoryPolicy
	if e := readValue(r.Runtime.PermissionBinding.Evidence[0].Locator, 64<<10, &policy); e != nil {
		t.Fatal(e)
	}
	policy.ReadRoots = factoryReadRoots(r.WorkspaceRoot, a)
	writeAny(t, filepath.Dir(r.Runtime.PermissionBinding.Evidence[0].Locator), "policy.json", policy)
	policyHash := hashFileTest(t, r.Runtime.PermissionBinding.Evidence[0].Locator)
	a.EffectivePolicySHA256 = policyHash
	r.Runtime.PermissionBinding.EffectivePolicySHA256 = policyHash
	r.Runtime.PermissionBinding.Evidence[0].SHA256 = policyHash
	writeAny(t, filepath.Dir(r.FactoryTest.AuthorizationPath), "authorization.json", a)
	r.FactoryTest.AuthorizationSHA256 = hashFileTest(t, r.FactoryTest.AuthorizationPath)
	file = writeAny(t, r.WorkspaceRoot, "request.json", r)
	runner := d.ExecRunner.(*fakeRunner)
	call := func(s LaunchSpec, operation, path string) error {
		cmd := exec.Command(bin, "workspace", "task", "run", operation, RunID(r.RequestKey), "--context", s.ContextPath, "--file", path, "--format", "json")
		cmd.Dir = s.CWD
		cmd.Env = []string{}
		for _, v := range os.Environ() {
			if !strings.HasPrefix(v, "PLY_TASK_RUN_CONTEXT=") {
				cmd.Env = append(cmd.Env, v)
			}
		}
		output, e := cmd.CombinedOutput()
		if e != nil {
			return errors.New(string(output))
		}
		var parsed map[string]any
		if e = json.Unmarshal(output, &parsed); e != nil {
			return e
		}
		return nil
	}
	runner.inside = func(s LaunchSpec) error {
		claim := acceptance(t, d, r)
		claim.SchemaVersion = 2
		claim.RuntimeClaim.ModelID = nil
		if e := call(s, "accept", writeAny(t, r.WorkspaceRoot, "claim.json", claim)); e != nil {
			return e
		}
		if e := call(s, "report", writeAny(t, r.WorkspaceRoot, "report.json", factorySemanticReport(t, d, r))); e != nil {
			return e
		}
		fakeCompletion(t, s)
		return nil
	}
	preview, e := PreviewStart(d, file)
	if e != nil {
		t.Fatal(e)
	}
	result, e := Start(d, file, *preview.(Preview).Confirmation)
	if e != nil {
		t.Fatalf("%v; callback: %v; %+v", e, runner.lastErr, result)
	}
	if result.Collection.State != "qualified" {
		t.Fatal(result)
	}
}

func factorySemanticReport(t *testing.T, d Dependencies, r Request) Report {
	v := semanticReport(t, d, r)
	v.BudgetUsage = &BudgetUsage{true, 0, 1, 0}
	return v
}
func TestFactorySlotCrashCannotRelaunch(t *testing.T) {
	d, r, f, _ := factoryFixture(t)
	p := dMustPreparation(t, d, r)
	if e := reserveFactory(d, r, p, f); e != nil {
		t.Fatal(e)
	}
	out, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	result, ok := out.(*Result)
	if !ok || result.Launch.State != "unknown" {
		t.Fatalf("%T %+v", out, out)
	}
	if _, e = Start(d, f, "unused"); e != nil {
		t.Fatal(e)
	}
	if d.ExecRunner.(*fakeRunner).calls != 0 {
		t.Fatal("crashed slot restarted")
	}
}

func TestFactoryBroadReadPolicyRejected(t *testing.T) {
	d, r, f, a := factoryFixture(t)
	var policy FactoryPolicy
	path := r.Runtime.PermissionBinding.Evidence[0].Locator
	if e := readValue(path, 64<<10, &policy); e != nil {
		t.Fatal(e)
	}
	policy.ReadRoots = []string{"/"}
	writeAny(t, filepath.Dir(path), filepath.Base(path), policy)
	h := hashFileTest(t, path)
	r.Runtime.PermissionBinding.Evidence[0].SHA256 = h
	r.Runtime.PermissionBinding.EffectivePolicySHA256 = h
	a.EffectivePolicySHA256 = h
	writeAny(t, filepath.Dir(r.FactoryTest.AuthorizationPath), "authorization.json", a)
	r.FactoryTest.AuthorizationSHA256 = hashFileTest(t, r.FactoryTest.AuthorizationPath)
	if _, _, e := validateFactory(r, dMustPreparation(t, d, r), f, d.Now()); e == nil {
		t.Fatal("broad filesystem read root accepted")
	}
}
