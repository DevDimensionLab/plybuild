package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// FactoryAuthorization is a local mandate, not cryptographic human attestation.
// Each slot binds its effective policy. All slots share the same immutable budget
// ledger under OutputRoot, so changing the authorization bytes cannot reset it.
type FactoryAuthorization struct {
	Kind                  string     `json:"kind"`
	ActorClaim            string     `json:"actor_claim"`
	Activity              string     `json:"activity"`
	OutputRoot            string     `json:"output_root"`
	WorkRoot              string     `json:"work_root"`
	TasksPerIteration     int        `json:"tasks_per_iteration"`
	Iterations            int        `json:"iterations"`
	MaxAgentStarts        int        `json:"max_agent_starts"`
	StartedUTC            string     `json:"started_at_utc"`
	DeadlineUTC           string     `json:"deadline_utc"`
	Model                 string     `json:"model"`
	ScenarioSHA256        string     `json:"scenario_sha256"`
	Ply                   Executable `json:"ply"`
	Codex                 Executable `json:"codex"`
	Git                   Executable `json:"git"`
	Python                Executable `json:"python"`
	EffectivePolicySHA256 string     `json:"effective_policy_sha256"`
}
type FactoryPolicy struct {
	Kind       string   `json:"kind"`
	ReadRoots  []string `json:"read_roots"`
	WriteRoots []string `json:"write_roots"`
	TempRoot   string   `json:"temp_root"`
	Network    bool     `json:"network"`
}
type factoryReservation struct {
	AuthorizationSHA256 string `json:"authorization_sha256"`
	RequestSHA256       string `json:"request_sha256"`
	RequestKey          string `json:"request_key"`
	Iteration           int    `json:"iteration"`
	TaskSlot            int    `json:"task_slot"`
	ReservedUTC         string `json:"reserved_at_utc"`
}

func within(root, path string) bool {
	rel, e := filepath.Rel(root, path)
	return e == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
func factorySlotRoot(r Request, a FactoryAuthorization) string {
	return filepath.Join(a.OutputRoot, "evidence", "provider", fmt.Sprintf("iteration-%d-slot-%d", r.FactoryTest.Iteration, r.FactoryTest.TaskSlot))
}
func factoryAuthorization(r Request) (a FactoryAuthorization, err error) {
	b, e := readFile(r.FactoryTest.AuthorizationPath, 64<<10, true)
	if e != nil {
		return a, e
	}
	if hash(b) != r.FactoryTest.AuthorizationSHA256 {
		return a, conflict("factory authorization bytes changed")
	}
	if e = decode(b, 64<<10, &a); e != nil {
		return a, e
	}
	if a.Kind != "PlyFactoryTestAuthorization@1" || !plain(a.Activity, 1, 256) || a.ActorClaim != r.HumanAuthority.ActorClaim || a.TasksPerIteration != 2 || a.Iterations < 1 || a.Iterations > 2 || a.MaxAgentStarts != a.Iterations*2 || a.Model != r.Runtime.Model || a.Ply != r.Runtime.PlyExecutable || a.Codex != r.Runtime.Executable || !digestPattern.MatchString(a.ScenarioSHA256) || a.EffectivePolicySHA256 != r.Runtime.PermissionBinding.EffectivePolicySHA256 {
		return a, invalid("invalid factory authorization bindings")
	}
	if !within("/private/tmp", a.OutputRoot) || a.WorkRoot != filepath.Join(a.OutputRoot, "work") || !within(filepath.Join(a.OutputRoot, "evidence"), r.FactoryTest.AuthorizationPath) {
		return a, conflict("factory roots or authorization escape")
	}
	for _, p := range []string{a.OutputRoot, a.WorkRoot, r.FactoryTest.AuthorizationPath} {
		if e = physical(p, false); e != nil {
			return a, e
		}
	}
	if info, e := os.Stat(a.OutputRoot); e != nil || !info.IsDir() || info.Mode().Perm()&0077 != 0 {
		return a, conflict("factory output root must be private")
	}
	start, se := time.Parse(time.RFC3339Nano, a.StartedUTC)
	end, ee := time.Parse(time.RFC3339Nano, a.DeadlineUTC)
	if se != nil || ee != nil || end.Sub(start) <= 0 || end.Sub(start) > 1200*time.Second {
		return a, invalid("factory deadline must be within 1200 seconds of first effect")
	}
	return a, nil
}
func validateFactory(r Request, p workspace.TaskPreparation, requestPath string, now time.Time) (a FactoryAuthorization, policy FactoryPolicy, err error) {
	a, err = factoryAuthorization(r)
	if err != nil {
		return
	}
	f := r.FactoryTest
	start, _ := time.Parse(time.RFC3339Nano, a.StartedUTC)
	deadline, _ := time.Parse(time.RFC3339Nano, a.DeadlineUTC)
	if f.Iteration > a.Iterations || now.Before(start) || !now.Before(deadline) {
		err = conflict("factory slot or total deadline exceeded")
		return
	}
	iteration := filepath.Join(a.WorkRoot, fmt.Sprintf("iteration-%d", f.Iteration))
	for _, path := range []string{requestPath, r.WorkspaceRoot, p.Plan.WorktreePath, p.Plan.Target.GitCommonDir, p.Plan.Target.ParentLocator} {
		if !within(iteration, path) {
			err = conflict("request/preparation/Git path escapes authorized iteration")
			return
		}
		if err = physical(path, false); err != nil {
			return
		}
	}
	common := p.Plan.Target.GitCommonDir
	if e := filepath.WalkDir(common, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return conflict("fixture Git metadata contains a symlink")
		}
		return nil
	}); e != nil {
		err = e
		return
	}
	if filepath.Base(common) != ".git" {
		err = conflict("fixture requires an independent Git common directory")
		return
	}
	for _, name := range []string{"objects/info/alternates", "objects/info/http-alternates"} {
		if _, e := os.Lstat(filepath.Join(common, name)); !os.IsNotExist(e) {
			err = conflict("Git alternates are forbidden")
			return
		}
	}
	for _, tool := range []Executable{a.Git, a.Python} {
		if err = verifyExecutable(tool); err != nil {
			return
		}
	}
	cmd := exec.Command(a.Git.Path, "--no-optional-locks", "-c", "core.hooksPath=/dev/null", "-c", "core.fsmonitor=false", "-C", p.Plan.WorktreePath, "remote")
	cmd.Env = factoryGitEnv()
	b, e := cmd.Output()
	if e != nil || len(b) != 0 {
		err = conflict("fixture remotes must be absent")
		return
	}
	found := false
	for _, ev := range r.Runtime.PermissionBinding.Evidence {
		if ev.Role != "effective_policy" {
			continue
		}
		if found || ev.SHA256 != a.EffectivePolicySHA256 {
			err = invalid("exactly one hash-bound effective factory policy is required")
			return
		}
		found = true
		b, e := readFile(ev.Locator, 64<<10, true)
		if e != nil {
			err = e
			return
		}
		if hash(b) != ev.SHA256 {
			err = conflict("effective policy changed")
			return
		}
		if err = decode(b, 64<<10, &policy); err != nil {
			return
		}
	}
	if !found || policy.Kind != "PlyFactoryPolicy@1" || policy.Network || policy.TempRoot != filepath.Join(iteration, "temp") {
		err = invalid("invalid narrow factory policy")
		return
	}
	// Callback stores are private to this disposable workspace. Provider evidence,
	// slot reservations, authorization, rig source and oracle are never writable.
	expected := []string{p.Plan.WorktreePath, common, filepath.Join(r.WorkspaceRoot, ".ply/task-runs/v1"), filepath.Join(r.WorkspaceRoot, ".ply/workflow/handoffs/v1"), filepath.Join(r.WorkspaceRoot, ".ply/projects.lock"), filepath.Join(r.WorkspaceRoot, ".ply/work-items.lock"), policy.TempRoot}
	sort.Strings(expected)
	if !equal(expected, policy.WriteRoots) {
		err = conflict("factory write roots differ from the exact Task/Git/return/temp scope")
		return
	}
	for _, path := range append(append([]string{}, policy.ReadRoots...), policy.WriteRoots...) {
		if err = physical(path, true); err != nil {
			return
		}
	}
	if !equal(policy.ReadRoots, factoryReadRoots(r.WorkspaceRoot, a)) {
		err = conflict("factory read roots differ from the named fixture, evidence and tool runtimes")
		return
	}
	return
}
func factoryGitEnv() []string {
	out := []string{}
	for _, s := range os.Environ() {
		if !strings.HasPrefix(s, "GIT_") {
			out = append(out, s)
		}
	}
	return append(out, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_OPTIONAL_LOCKS=0")
}
func reserveFactory(d Dependencies, r Request, p workspace.TaskPreparation, file string) error {
	a, _, e := validateFactory(r, p, file, d.Now())
	if e != nil {
		return e
	}
	root := filepath.Join(a.OutputRoot, "evidence", "factory-slots")
	if e = privateDir(root); e != nil {
		return e
	}
	lock, e := lockFile(filepath.Join(root, "ledger.lock"))
	if e != nil {
		return e
	}
	defer unlockFile(lock)
	// Immutable common authorization prevents a new file, key or deadline from
	// granting a second budget at the same output root.
	common := a
	common.EffectivePolicySHA256 = ""
	ledger := filepath.Join(root, "authorization.json")
	var previous FactoryAuthorization
	if e = readValue(ledger, 64<<10, &previous); os.IsNotExist(e) {
		e = d.writeValue(ledger, common)
	} else if e == nil && !equal(previous, common) {
		e = conflict("factory common authorization is immutable")
	}
	if e != nil {
		return e
	}
	slot := filepath.Join(root, fmt.Sprintf("%d-%d.json", r.FactoryTest.Iteration, r.FactoryTest.TaskSlot))
	if _, e = os.Lstat(slot); !os.IsNotExist(e) {
		return conflict("factory slot is already reserved; never relaunch even with a new request key")
	}
	entries, e := os.ReadDir(root)
	if e != nil {
		return e
	}
	count := 0
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".json") && entry.Name() != "authorization.json" {
			count++
		}
	}
	if count >= a.MaxAgentStarts {
		return conflict("factory agent-start budget exhausted")
	}
	return d.writeValue(slot, factoryReservation{r.FactoryTest.AuthorizationSHA256, digest(r), r.RequestKey, r.FactoryTest.Iteration, r.FactoryTest.TaskSlot, d.Now().UTC().Format(time.RFC3339Nano)})
}
func tomlQuote(s string) string { b, _ := json.Marshal(s); return string(b) }

func factoryProfileWriteRoots(p FactoryPolicy) ([]string, error) {
	// A linked worktree's .git pointer makes its resolved metadata directory
	// read-only in Codex unless it is named explicitly, even below a write root.
	// Only add descendants already authorized by the hash-bound factory policy.
	seen := map[string]bool{}
	for _, root := range p.WriteRoots {
		seen[root] = true
		info, err := os.Lstat(root)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			continue
		}
		marker := filepath.Join(root, ".git")
		info, err = os.Lstat(marker)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			continue
		}
		data, err := readFile(marker, 4096, false)
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(string(data), "gitdir: ") {
			return nil, invalid("invalid linked Git metadata pointer")
		}
		directory := strings.TrimSpace(string(data[len("gitdir: "):]))
		if directory == "" {
			return nil, invalid("empty linked Git metadata pointer")
		}
		if !filepath.IsAbs(directory) {
			directory = filepath.Join(root, directory)
		}
		directory = filepath.Clean(directory)
		allowed := false
		for _, authorized := range p.WriteRoots {
			allowed = allowed || directory == authorized || within(authorized, directory)
		}
		if !allowed {
			return nil, conflict("linked Git metadata is outside the named write roots")
		}
		if err := physical(directory, false); err != nil {
			return nil, err
		}
		info, err = os.Stat(directory)
		if err != nil || !info.IsDir() {
			return nil, conflict("linked Git metadata must be a directory")
		}
		seen[directory] = true
	}
	roots := make([]string, 0, len(seen))
	for root := range seen {
		roots = append(roots, root)
	}
	sort.Strings(roots)
	return roots, nil
}

func execArgv(r Request, p FactoryPolicy, cwd, instructions string) ([]string, error) {
	writes, err := factoryProfileWriteRoots(p)
	if err != nil {
		return nil, err
	}
	// The managed profile is passed explicitly; no inherited developer profile,
	// resume, search, MCP configuration, hook or daemon is part of this adapter.
	// The named roots do not replace platform startup requirements such as
	// macOS dyld's read of the root directory itself. Keep Codex's runtime
	// baseline without granting recursive root access or additional work roots.
	entries := []string{`":root"="none"`, `":minimal"="read"`}
	for _, path := range p.ReadRoots {
		entries = append(entries, tomlQuote(path)+`="read"`)
	}
	for _, path := range writes {
		entries = append(entries, tomlQuote(path)+`="write"`)
	}
	profile := r.Runtime.PermissionBinding.ProfileID
	args := []string{r.Runtime.Executable.Path, "--no-daemon", "-a", "never", "exec", "--json", "--ephemeral", "--ignore-user-config", "--cd", cwd, "--model", r.Runtime.Model}
	config := []string{"model_reasoning_effort=\"low\"", "default_permissions=" + tomlQuote(profile), "permissions." + profile + "={filesystem={" + strings.Join(entries, ",") + "},network={enabled=false}}", `web_search="disabled"`, `mcp_servers={}`, `apps._default.enabled=false`, `features.hooks=false`, `features.multi_agent=false`}
	for _, c := range config {
		args = append(args, "-c", c)
	}
	return append(args, "Read and execute the private Task run instructions at "+instructions+". Your first action is the typed accept callback before any target write. Complete one turn and stop."), nil
}

// A slot reserved before request publication is still spent. Same-request retry
// reads unknown instead of trying again; another request key is a conflict.
func existingFactoryReservation(r Request) (*Result, error) {
	a, e := factoryAuthorization(r)
	if e != nil {
		return nil, e
	}
	path := filepath.Join(a.OutputRoot, "evidence", "factory-slots", fmt.Sprintf("%d-%d.json", r.FactoryTest.Iteration, r.FactoryTest.TaskSlot))
	var slot factoryReservation
	e = readValue(path, 64<<10, &slot)
	if os.IsNotExist(e) {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	if slot.RequestKey != r.RequestKey || slot.RequestSHA256 != digest(r) || slot.AuthorizationSHA256 != r.FactoryTest.AuthorizationSHA256 {
		return nil, conflict("factory slot is already spent by another request")
	}
	out := initial(r)
	out.Launch.State = "unknown"
	out.Process.State = "unknown"
	out.Collection.State = "unknown"
	out.Reasons = []Reason{{"task_run_factory_slot_unknown", "Slot reserved before request publication; no second start is permitted."}}
	return &out, nil
}

func factoryReadRoots(workspaceRoot string, a FactoryAuthorization) []string {
	seen := map[string]bool{}
	for _, path := range []string{workspaceRoot, filepath.Join(a.OutputRoot, "evidence"), "/System", "/usr", "/bin", "/dev", "/Library", "/opt/homebrew", filepath.Dir(a.Ply.Path), filepath.Dir(a.Codex.Path), filepath.Dir(a.Git.Path), filepath.Dir(a.Python.Path)} {
		seen[path] = true
	}
	roots := make([]string, 0, len(seen))
	for path := range seen {
		roots = append(roots, path)
	}
	sort.Strings(roots)
	return roots
}
