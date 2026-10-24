package taskexecute

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/pelletier/go-toml/v2"
)

const BoundStartRequestKind = "ply.workflow.bound-start-request"
const BoundStartPreviewKind = "ply.workflow.bound-start-preview"
const BoundStartResultKind = "ply.workflow.bound-start-status"

type BoundStartInput struct {
	SpecID            string                     `json:"spec_id"`
	Target            workspace.QueueTargetInput `json:"target"`
	StartMode         string                     `json:"start_mode"`
	HerdrWorkspace    string                     `json:"herdr_workspace"`
	PermissionProfile string                     `json:"permission_profile"`
}

type BoundStartRequest struct {
	Kind          string          `json:"kind"`
	SchemaVersion int             `json:"schema_version"`
	RequestID     string          `json:"request_id"`
	Workspace     string          `json:"workspace"`
	Input         BoundStartInput `json:"input"`
	Binding       string          `json:"binding"`
}

type boundPolicySource struct {
	Path   string `json:"path"`
	State  string `json:"state"`
	SHA256 string `json:"sha256,omitempty"`
}

type BoundStartBasis struct {
	Goal                 workspace.TaskGoalExecutePreview `json:"goal"`
	Runtime              RuntimePreview                   `json:"runtime"`
	Environment          taskrun.GoalStartEnvironment     `json:"environment"`
	PolicySources        []boundPolicySource              `json:"policy_sources"`
	PermissionDefinition any                              `json:"permission_definition"`
	DependencySHA256     string                           `json:"dependency_sha256,omitempty"`
	RegistrySHA256       string                           `json:"registry_sha256"`
	ProjectSHA256        string                           `json:"project_sha256"`
	Trust                string                           `json:"trust"`
	RuntimeAcceptance    string                           `json:"runtime_acceptance"`
	Effects              []string                         `json:"effects"`
}

type BoundStartPreview struct {
	Kind          string            `json:"kind"`
	SchemaVersion int               `json:"schema_version"`
	State         string            `json:"state"`
	Request       BoundStartRequest `json:"request"`
	Basis         BoundStartBasis   `json:"basis"`
	Error         *BoundStartError  `json:"error,omitempty"`
}

type BoundStartError struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func (e *BoundStartError) Error() string   { return e.Code + ": " + e.Detail }
func boundError(code, detail string) error { return &BoundStartError{code, detail} }
func boundReason(err error, fallback string) *BoundStartError {
	var e *BoundStartError
	if errors.As(err, &e) {
		return e
	}
	return &BoundStartError{fallback, err.Error()}
}

type BoundStartStatus struct {
	Kind            string                `json:"kind"`
	SchemaVersion   int                   `json:"schema_version"`
	RequestID       string                `json:"request_id"`
	Workspace       string                `json:"workspace"`
	State           string                `json:"state"`
	ExecutionEffect string                `json:"execution_effect"`
	Binding         string                `json:"binding,omitempty"`
	ExpectedRunID   string                `json:"expected_run_id,omitempty"`
	Run             *taskrun.InventoryRun `json:"run,omitempty"`
	NextAction      string                `json:"next_action"`
	Error           *BoundStartError      `json:"error,omitempty"`
}

type boundStartIntent struct {
	Kind          string            `json:"kind"`
	SchemaVersion int               `json:"schema_version"`
	Request       BoundStartRequest `json:"request"`
	Basis         BoundStartBasis   `json:"basis"`
	RunID         string            `json:"run_id"`
	RecordedAtUTC string            `json:"recorded_at_utc"`
}

type boundStartReceipt struct {
	RequestSHA256 string           `json:"request_sha256"`
	Result        Result           `json:"result"`
	Error         *BoundStartError `json:"error,omitempty"`
}

var boundRequestID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{7,127}$`)
var boundDigest = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

func validateBoundRequest(r BoundStartRequest, requireBinding bool) error {
	if r.Kind != BoundStartRequestKind || r.SchemaVersion != 1 {
		return boundError("bound_start_unsupported", "only ply.workflow.bound-start-request schema_version 1 is supported")
	}
	if !boundRequestID.MatchString(r.RequestID) {
		return boundError("bound_start_invalid_request", "request_id must contain 8-128 letters, digits, dots, underscores or hyphens and begin with a letter or digit")
	}
	if r.Input.StartMode != "interactive" {
		return boundError("bound_start_unsupported", "bound start v1 supports only an explicit Goal with interactive Codex; no next, headless, restart or fallback")
	}
	if r.Input.SpecID == "" || strings.TrimSpace(r.Input.SpecID) != r.Input.SpecID {
		return boundError("bound_start_invalid_request", "an explicit Goal spec_id is required")
	}
	if err := physicalPath(r.Workspace); err != nil {
		return boundError("bound_start_invalid_request", err.Error())
	}
	if r.Input.HerdrWorkspace == "" || len(r.Input.HerdrWorkspace) > 128 || strings.ContainsAny(r.Input.HerdrWorkspace, "\x00\r\n") {
		return boundError("bound_start_invalid_request", "an explicit Herdr workspace is required")
	}
	if requireBinding && !boundDigest.MatchString(r.Binding) {
		return boundError("bound_start_invalid_binding", "the unchanged preview binding is required")
	}
	return nil
}

func boundHash(v any) string {
	b, err := taskrun.Canonical(v)
	if err != nil {
		return ""
	}
	return digestBytes(b)
}
func boundBinding(r BoundStartRequest, b BoundStartBasis) string {
	// Ply is the command handling this request, not part of its authority.
	b.Runtime.Runtime.PlyExecutable = taskrun.Executable{}
	return legacyBoundBinding(r, b)
}

// Existing stored intents retain their original digest and history. A legacy
// request with only an opaque digest cannot reconstruct a lost preview basis.
func legacyBoundBinding(r BoundStartRequest, b BoundStartBasis) string {
	r.Binding = ""
	if b.DependencySHA256 != "" {
		b.RegistrySHA256 = ""
	}
	return boundHash(struct {
		Request BoundStartRequest `json:"request"`
		Basis   BoundStartBasis   `json:"basis"`
	}{r, b})
}

func boundRuntimeOptions(in BoundStartInput) RuntimeOptions {
	return RuntimeOptions{HerdrWorkspace: in.HerdrWorkspace, PermissionProfile: in.PermissionProfile}
}

func boundEnvironment() (taskrun.GoalStartEnvironment, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return taskrun.GoalStartEnvironment{}, err
	}
	home, err = filepath.EvalSymlinks(home)
	if err != nil {
		return taskrun.GoalStartEnvironment{}, err
	}
	codexHome := os.Getenv("CODEX_HOME")
	if codexHome == "" {
		codexHome = filepath.Join(home, ".codex")
	}
	codexHome, err = filepath.EvalSymlinks(codexHome)
	if err != nil {
		return taskrun.GoalStartEnvironment{}, err
	}
	return taskrun.GoalStartEnvironment{Home: home, CodexHome: codexHome, Path: os.Getenv("PATH")}, nil
}

func observeBoundPolicy(p workspace.TaskGoalExecutePreview, env taskrun.GoalStartEnvironment, profile string) ([]boundPolicySource, any, error) {
	etc, err := filepath.EvalSymlinks("/etc")
	if err != nil {
		return nil, nil, err
	}
	paths := map[string]bool{}
	for _, dir := range []string{filepath.Join(etc, "codex"), env.CodexHome} {
		for _, name := range []string{"config.toml", "managed_config.toml", "requirements.toml", "AGENTS.md", "AGENTS.override.md", "rules"} {
			paths[filepath.Join(dir, name)] = true
		}
	}
	for _, path := range []string{filepath.Join("/Library/Managed Preferences", "com.openai.codex.plist"), filepath.Join("/Library/Managed Preferences", filepath.Base(env.Home), "com.openai.codex.plist"), filepath.Join(env.Home, "Library/Preferences/com.openai.codex.plist")} {
		paths[path] = true
	}
	// Parent Git tree binds the future Task's repository files. Bind additional
	// configuration/instruction layers outside that tree, including absent paths.
	for _, start := range []string{p.Target.ParentLocator, filepath.Dir(p.WorktreePath)} {
		for dir := start; ; dir = filepath.Dir(dir) {
			paths[filepath.Join(dir, ".codex/config.toml")] = true
			paths[filepath.Join(dir, ".codex/rules")] = true
			paths[filepath.Join(dir, "AGENTS.md")] = true
			paths[filepath.Join(dir, "AGENTS.override.md")] = true
			if filepath.Dir(dir) == dir {
				break
			}
		}
	}
	keys := make([]string, 0, len(paths))
	for path := range paths {
		keys = append(keys, path)
	}
	sort.Strings(keys)
	out := []boundPolicySource{}
	for _, path := range keys {
		source, err := observeBoundPolicySource(path)
		if err != nil {
			return nil, nil, err
		}
		out = append(out, source)
	}
	body, err := os.ReadFile(filepath.Join(env.CodexHome, "config.toml"))
	if err != nil {
		return nil, nil, err
	}
	var cfg map[string]any
	if err = toml.Unmarshal(body, &cfg); err != nil {
		return nil, nil, err
	}
	permissions, _ := cfg["permissions"].(map[string]any)
	definition := permissions[profile]
	if strings.HasPrefix(profile, ":") {
		definition = map[string]any{"builtin": profile}
	}
	return out, definition, nil
}

func observeBoundPolicySource(path string) (boundPolicySource, error) {
	out := boundPolicySource{Path: path, State: "absent"}
	if err := physicalPath(path); err != nil {
		return out, err
	}
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return out, nil
	}
	if err != nil {
		return out, err
	}
	if info.IsDir() {
		entries := []boundPolicySource{}
		total := int64(0)
		err = filepath.WalkDir(path, func(file string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				return nil
			}
			info, err := entry.Info()
			if err != nil {
				return err
			}
			total += info.Size()
			if len(entries) >= 256 || total > 16<<20 {
				return fmt.Errorf("bound policy directory exceeds 256 files or 16 MiB: %s", path)
			}
			item, err := observeBoundPolicySource(file)
			if err != nil {
				return err
			}
			item.Path, _ = filepath.Rel(path, file)
			entries = append(entries, item)
			return nil
		})
		if err != nil {
			return out, err
		}
		out.State, out.SHA256 = "directory", boundHash(entries)
		return out, nil
	}
	if !info.Mode().IsRegular() || info.Size() > 4<<20 {
		return out, fmt.Errorf("unsupported bound policy source: %s", path)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return out, err
	}
	out.State, out.SHA256 = "file", digestBytes(body)
	return out, nil
}

func buildBoundBasis(d taskrun.Dependencies, s *workspace.TaskSpecSession, r BoundStartRequest) (BoundStartBasis, error) {
	p, err := s.PreviewInitialTaskGoalExecution(workspace.TaskGoalExecuteInput{Target: r.Input.Target, SpecID: r.Input.SpecID})
	if err != nil {
		return BoundStartBasis{}, err
	}
	if err = checkOrigin(d, p.Target); err != nil {
		return BoundStartBasis{}, err
	}
	if err = rejectEarlierLaunch(p); err != nil {
		return BoundStartBasis{}, err
	}
	control, err := d.Executable()
	if err != nil {
		return BoundStartBasis{}, err
	}
	runtime, err := PreviewRuntime(p.Goal.Executor, control, boundRuntimeOptions(r.Input))
	if err != nil {
		return BoundStartBasis{}, err
	}
	if runtime.Runtime.Provider != "codex" {
		return BoundStartBasis{}, boundError("bound_start_unsupported", "bound start v1 supports interactive Codex only")
	}
	env, err := boundEnvironment()
	if err != nil {
		return BoundStartBasis{}, err
	}
	sources, definition, err := observeBoundPolicy(p, env, runtime.Policy.PermissionProfile)
	if err != nil {
		return BoundStartBasis{}, err
	}
	out := BoundStartBasis{Goal: p, Runtime: runtime, Environment: env, PolicySources: sources, PermissionDefinition: definition,
		DependencySHA256: workspace.TaskStartupDependencySHA256(s.Registry, p.Goal.TaskID), RegistrySHA256: s.Registry.RawSHA256, ProjectSHA256: boundHash(s.Projects),
		Trust: "Native Codex onboarding only; no persistent or process-local trust grant by Ply.", RuntimeAcceptance: "requested_only; the recipient must confirm actual runtime and permissions before Task writes",
		Effects: []string{"Preserve this immutable request identity before any execution effect", "Publish the exact Goal execution Spec, assessment and selection, and bind the queue", "Reserve and create the shown Task branch and worktree from the bound base", "Preserve launch policy, historical Ply observations, native Run and delivery handoff", "Create at most one Herdr tab in the shown workspace and start its assigned Codex agent", "Require actual runtime acceptance and the unchanged delivery/acceptance agreement; no new Git or delivery authority"}}
	if _, err := taskrun.Canonical(out); err != nil {
		return BoundStartBasis{}, boundError("bound_start_unsupported", "selected runtime or permission policy cannot be canonically bound: "+err.Error())
	}
	return out, nil
}

func rejectEarlierLaunch(p workspace.TaskGoalExecutePreview) error {
	path := filepath.Join(filepath.Dir(p.AcceptancePath), "attempts")
	if err := physicalPath(path); err != nil {
		return err
	}
	entries, err := os.ReadDir(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			return boundError("bound_start_existing_attempt", "an earlier launcher artifact must be inspected before a new initial approval")
		}
		for _, name := range []string{"intent.json", "workflow-request.json"} {
			if _, err := os.Lstat(filepath.Join(path, entry.Name(), name)); err == nil || !os.IsNotExist(err) {
				return boundError("bound_start_existing_attempt", "this Goal already has a preserved launcher attempt; inspect its original identity instead of approving a new start")
			}
		}
	}
	return nil
}

func PreviewBoundStart(d taskrun.Dependencies, requestID string, input BoundStartInput) (out BoundStartPreview, err error) {
	out = BoundStartPreview{Kind: BoundStartPreviewKind, SchemaVersion: 1, State: "rejected"}
	defer func() {
		if err != nil {
			out.Error = boundReason(err, "bound_start_unavailable")
		}
	}()
	if input.StartMode == "" {
		input.StartMode = "interactive"
	}
	if input.HerdrWorkspace == "" {
		input.HerdrWorkspace = os.Getenv("HERDR_WORKSPACE_ID")
	}
	r := BoundStartRequest{Kind: BoundStartRequestKind, SchemaVersion: 1, RequestID: requestID, Workspace: "/", Input: input}
	if err = validateBoundRequest(r, false); err != nil {
		return out, err
	}
	observed, err := d.Workflow.Workspace.ObserveContaining()
	if err != nil {
		return out, err
	}
	r.Workspace = observed.Observation.Root
	err = workspace.WithTaskSpecRegistrationSnapshot(d.Workspace, r.Workspace, func(s *workspace.TaskSpecSession) error {
		var err error
		out.Basis, err = buildBoundBasis(d, s, r)
		return err
	})
	if err != nil {
		return out, err
	}
	r.Binding = boundBinding(r, out.Basis)
	out.Request, out.State = r, "ready"
	return out, nil
}

func boundIntentPath(root, id string) string {
	return filepath.Join(root, ".ply", "bound-goal-starts", strings.TrimPrefix(digestBytes([]byte(id)), "sha256:"), "intent.json")
}
func boundStatus(root, id string) BoundStartStatus {
	return BoundStartStatus{Kind: BoundStartResultKind, SchemaVersion: 1, Workspace: root, RequestID: id, State: "not_recorded", ExecutionEffect: "not_confirmed", NextAction: "No intent was found. Absence alone is not proof of no execution effect; do not invent another identity to recover a lost response."}
}

func readBoundIntent(root, id string) (boundStartIntent, error) {
	var intent boundStartIntent
	_, err := readPreservedJSON(boundIntentPath(root, id), &intent)
	if err != nil {
		return intent, err
	}
	if intent.Kind != "ply.workflow.bound-start-intent" || intent.SchemaVersion != 1 || intent.Request.RequestID != id || intent.Request.Workspace != root || validateBoundRequest(intent.Request, true) != nil || (boundBinding(intent.Request, intent.Basis) != intent.Request.Binding && legacyBoundBinding(intent.Request, intent.Basis) != intent.Request.Binding) || intent.RunID != boundRunID(intent.Request) {
		return intent, boundError("bound_start_outcome_unknown", "preserved request identity or basis is damaged; never restart")
	}
	return intent, nil
}
func boundRequestKey(r BoundStartRequest) string {
	return "bound-start/" + strings.TrimPrefix(boundHash([]string{r.RequestID, r.Binding}), "sha256:")
}
func boundRunID(r BoundStartRequest) string {
	return "wfr_" + strings.TrimPrefix(boundHash([]string{r.Workspace, boundRequestKey(r)}), "sha256:")
}

func StatusBoundStart(d taskrun.Dependencies, root, id string) (out BoundStartStatus, err error) {
	out = boundStatus(root, id)
	if !boundRequestID.MatchString(id) {
		err = boundError("bound_start_invalid_request", "invalid request_id")
		out.Error = boundReason(err, "bound_start_invalid_request")
		return
	}
	observed, err := workspace.ObserveRoot(d.Workspace, root)
	if err != nil || observed.Root != root {
		if err == nil {
			err = fmt.Errorf("workspace must be its physical root")
		}
		out.Error = boundReason(err, "bound_start_invalid_request")
		return out, err
	}
	intent, err := readBoundIntent(root, id)
	if os.IsNotExist(err) {
		return out, nil
	}
	if err != nil {
		out.State, out.ExecutionEffect, out.Error = "outcome_unknown", "unknown", boundReason(err, "bound_start_outcome_unknown")
		return out, err
	}
	out.State, out.ExecutionEffect, out.Binding, out.ExpectedRunID = "retained_attempt", "unknown", intent.Request.Binding, intent.RunID
	out.NextAction = "Inspect this same identity and its preserved Run. Apply replay never resumes, retries, reopens, or creates another writer."
	out.Run, err = taskrun.ReadBoundStartRun(d, root, intent.RunID)
	if os.IsNotExist(err) {
		return out, nil
	}
	if err != nil {
		out.State, out.Error = "outcome_unknown", boundReason(err, "bound_start_outcome_unknown")
		return out, err
	}
	out.State, out.ExecutionEffect = "existing_run", "reserved"
	if out.Run.Unresolved || out.Run.Freshness == "unknown" {
		out.ExecutionEffect = "unknown"
	}
	return out, nil
}

type approvedStartGuard struct {
	d        taskrun.Dependencies
	request  BoundStartRequest
	basis    BoundStartBasis
	registry string
}

func (g *approvedStartGuard) Check(root string, r workspace.WorkItemRegistry) error {
	observed := r.RawSHA256
	if g.basis.DependencySHA256 != "" {
		observed = workspace.TaskStartupDependencySHA256(r, g.basis.Goal.Goal.TaskID)
	}
	if root != g.request.Workspace || observed != g.registry {
		return boundError("bound_start_stale", "native registry changed before the approved start effect")
	}
	projects, repos, err := g.d.Workspace.Projects.Snapshot(root)
	if err != nil {
		return err
	}
	if boundHash(workspace.ProjectSnapshot{Projects: projects, Repos: repos}) != g.basis.ProjectSHA256 {
		return boundError("bound_start_stale", "project or repository binding changed")
	}
	requireTask := false
	for _, task := range r.Tasks {
		if task.ID == g.basis.Goal.Goal.TaskID && task.Worktree != nil {
			requireTask = true
		}
	}
	if err := g.git(requireTask); err != nil {
		return err
	}
	if requireTask {
		for _, source := range g.basis.PolicySources {
			rel, err := filepath.Rel(g.basis.Goal.Target.ParentLocator, source.Path)
			if err != nil || (rel != ".codex/config.toml" && rel != ".codex/rules" && rel != "AGENTS.md" && rel != "AGENTS.override.md") {
				continue
			}
			actual, err := observeBoundPolicySource(filepath.Join(g.basis.Goal.WorktreePath, rel))
			if err != nil || actual.State != source.State || actual.SHA256 != source.SHA256 {
				return boundError("bound_start_stale", "Task policy or instructions differ from the approved source, including ignored files")
			}
		}
	}
	return g.runtime()
}

func (g *approvedStartGuard) git(requireTask bool) error {
	p := g.basis.Goal
	for _, target := range []struct {
		path, ref string
		optional  bool
	}{{p.Target.ParentLocator, p.Target.ParentRef, false}, {p.WorktreePath, "refs/heads/" + p.Branch, !requireTask}} {
		if _, err := os.Lstat(target.path); os.IsNotExist(err) && target.optional {
			continue
		} else if err != nil {
			return boundError("bound_start_stale", "approved Git worktree is unavailable")
		}
		observed, err := g.d.Workspace.WorkGit.ObserveWorktree(target.path)
		if err != nil || !observed.InventoryMatch || !observed.Clean || observed.Locator != target.path || observed.Ref != target.ref || observed.OID != p.ParentOID || observed.Tree != p.ParentTree || observed.GitCommonDir != p.Target.GitCommonDir {
			return boundError("bound_start_stale", "approved physical Git worktree, clean ref, OID or tree changed")
		}
		full, err := g.d.Workspace.IntegrationGit.ObserveIntegrationWorktree(target.path, target.ref)
		if err != nil || !full.Clean || !full.Symbolic || len(full.InProgress) != 0 || full.OID != p.ParentOID || full.Tree != p.ParentTree {
			return boundError("bound_start_stale", "approved Git start changed or has an operation in progress")
		}
	}
	return nil
}
func (g *approvedStartGuard) Committed(_ string, r workspace.WorkItemRegistry) {
	g.registry = r.RawSHA256
	if g.basis.DependencySHA256 != "" {
		g.registry = workspace.TaskStartupDependencySHA256(r, g.basis.Goal.Goal.TaskID)
	}
}
func (g *approvedStartGuard) runtime() error {
	control, err := g.d.Executable()
	if err != nil {
		return err
	}
	runtime, err := PreviewRuntime(g.basis.Goal.Goal.Executor, control, boundRuntimeOptions(g.request.Input))
	if err != nil {
		return boundError("bound_start_stale", err.Error())
	}
	env, err := boundEnvironment()
	if err != nil {
		return boundError("bound_start_stale", err.Error())
	}
	sources, definition, err := observeBoundPolicy(g.basis.Goal, env, runtime.Policy.PermissionProfile)
	if err != nil {
		return boundError("bound_start_stale", err.Error())
	}
	if !sameRuntimeAuthority(runtime, g.basis.Runtime) || env != g.basis.Environment || boundHash(sources) != boundHash(g.basis.PolicySources) || boundHash(definition) != boundHash(g.basis.PermissionDefinition) {
		return boundError("bound_start_stale", "runtime, policy, trust, executable or launch environment changed")
	}
	return nil
}

func ApplyBoundStart(d taskrun.Dependencies, r BoundStartRequest) (out BoundStartStatus, err error) {
	out = boundStatus(r.Workspace, r.RequestID)
	out.State, out.ExecutionEffect = "rejected", "none"
	defer func() {
		if err != nil && out.State == "rejected" {
			if retained, readErr := readBoundIntent(r.Workspace, r.RequestID); readErr == nil && boundHash(retained.Request) == boundHash(r) {
				out.State, out.ExecutionEffect, out.Binding, out.ExpectedRunID = "retained_attempt", "unknown", r.Binding, retained.RunID
				out.NextAction = "The immutable attempt is retained. Read this same identity; do not retry or create another writer."
			}
		}
		if err != nil && out.Error == nil {
			fallback := "bound_start_unavailable"
			if out.ExecutionEffect == "unknown" {
				fallback = "bound_start_outcome_unknown"
			}
			out.Error = boundReason(err, fallback)
		}
	}()
	if err = validateBoundRequest(r, true); err != nil {
		return out, err
	}
	// Identity lookup precedes all fresh-basis checks. The first start may have
	// changed its own queue, runtime trust and Run state legitimately.
	if old, readErr := readBoundIntent(r.Workspace, r.RequestID); readErr == nil {
		if boundHash(old.Request) != boundHash(r) {
			return out, boundError("bound_start_request_conflict", "request_id already has different immutable content")
		}
		return StatusBoundStart(d, r.Workspace, r.RequestID)
	} else if !os.IsNotExist(readErr) {
		out.State, out.ExecutionEffect = "outcome_unknown", "unknown"
		return out, readErr
	}
	if d.Fault != nil {
		if err = d.Fault("bound_start_before_lock"); err != nil {
			return out, err
		}
	}
	var intent boundStartIntent
	owned := false
	err = workspace.WithTaskSpecRegistrationSnapshot(d.Workspace, r.Workspace, func(s *workspace.TaskSpecSession) error {
		if old, e := readBoundIntent(r.Workspace, r.RequestID); e == nil {
			if boundHash(old.Request) != boundHash(r) {
				return boundError("bound_start_request_conflict", "request_id already has different immutable content")
			}
			intent = old
			return nil
		} else if !os.IsNotExist(e) {
			return e
		}
		basis, e := buildBoundBasis(d, s, r)
		if e != nil {
			return boundError("bound_start_stale", e.Error())
		}
		if boundBinding(r, basis) != r.Binding {
			return boundError("bound_start_stale", "approved preview no longer matches Goal, queue, base, runtime, policy, executable, target, trust or agreement")
		}
		intent = boundStartIntent{Kind: "ply.workflow.bound-start-intent", SchemaVersion: 1, Request: r, Basis: basis, RunID: boundRunID(r), RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
		body, e := taskrun.Canonical(intent)
		if e != nil {
			return e
		}
		if e = writeOnce(boundIntentPath(r.Workspace, r.RequestID), body, 0600); e != nil {
			return e
		}
		owned = true
		if d.Fault != nil {
			if e := d.Fault("bound_start_intent_published"); e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		return out, err
	}
	if !owned {
		return StatusBoundStart(d, r.Workspace, r.RequestID)
	}
	out.State, out.ExecutionEffect, out.Binding, out.ExpectedRunID = "retained_attempt", "unknown", r.Binding, intent.RunID
	if d.Fault != nil {
		if err = d.Fault("bound_start_after_intent"); err != nil {
			return out, err
		}
	}
	guard := &approvedStartGuard{d: d, request: r, basis: intent.Basis, registry: intent.Basis.RegistrySHA256}
	locked := d
	if intent.Basis.DependencySHA256 != "" {
		guard.registry = intent.Basis.DependencySHA256
	}
	locked.Workspace = workspace.WithTaskStartGuard(d.Workspace, guard)
	locked.Workflow.TaskWorkspace = &locked.Workspace
	locked.ApprovedGoalStart = func() error {
		return workspace.WithTaskSpecRegistrationSnapshot(locked.Workspace, r.Workspace, func(*workspace.TaskSpecSession) error { return nil })
	}
	in := Input{SpecID: r.Input.SpecID, Target: r.Input.Target, Runtime: boundRuntimeOptions(r.Input), ExpectedTaskID: intent.Basis.Goal.Goal.TaskID, ExpectedGoal: &intent.Basis.Goal.Goal.Goal, boundEnvironment: &intent.Basis.Environment, expectedRuntime: &intent.Basis.Runtime, boundRequestKey: boundRequestKey(r)}
	result, runErr := Execute(locked, in)
	receipt := boundStartReceipt{RequestSHA256: boundHash(r), Result: result}
	if runErr != nil {
		receipt.Error = boundReason(runErr, "bound_start_outcome_unknown")
	}
	body, e := taskrun.Canonical(receipt)
	if e == nil {
		e = writeOnce(filepath.Join(filepath.Dir(boundIntentPath(r.Workspace, r.RequestID)), "receipt.json"), body, 0600)
	}
	if e != nil {
		return out, boundError("bound_start_outcome_unknown", "attempt is retained but its reply receipt could not be preserved: "+e.Error())
	}
	out, err = StatusBoundStart(d, r.Workspace, r.RequestID)
	if err == nil && runErr != nil {
		out.Error = receipt.Error
		err = runErr
	}
	return out, err
}

func ReadBoundStartRequest(path string) (BoundStartRequest, error) {
	var r BoundStartRequest
	body, err := readPreservedJSON(path, &r)
	if err != nil {
		return r, boundError("bound_start_invalid_request", err.Error())
	}
	value, err := canonicaljson.DecodeStrict(body)
	if err != nil {
		return r, boundError("bound_start_invalid_request", err.Error())
	}
	canonical, err := canonicaljson.Marshal(value)
	expected, encodeErr := taskrun.Canonical(r)
	if err != nil || encodeErr != nil || !bytes.Equal(canonical, expected) {
		return r, boundError("bound_start_invalid_request", "all request fields must be present with their declared JSON types; null and omitted fields are not supported")
	}
	return r, validateBoundRequest(r, true)
}
