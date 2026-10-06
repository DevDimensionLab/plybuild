package taskrun

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// This grant belongs to the Herdr request, never to native TaskRun Runtime.
type CodexProjectTrust struct {
	Mode           string `json:"mode"`
	RepositoryRoot string `json:"repository_root"`
}

type workflowTrustPath struct {
	Path   string `json:"path"`
	Kind   string `json:"kind"`
	Device uint64 `json:"device"`
	Inode  uint64 `json:"inode"`
	SHA256 string `json:"sha256,omitempty"`
}
type workflowTrustFacts struct {
	RepositoryRoot        string              `json:"repository_root"`
	WorktreeRoot          string              `json:"worktree_root"`
	CWD                   string              `json:"cwd"`
	GitDir                string              `json:"git_dir"`
	CommonDir             string              `json:"common_dir"`
	CodexHome             string              `json:"codex_home"`
	UserHome              string              `json:"user_home"`
	Paths                 []workflowTrustPath `json:"paths"`
	PermissionProfile     string              `json:"permission_profile"`
	EffectivePolicySHA256 string              `json:"effective_policy_sha256"`
}

func workflowValidateTrust(r WorkflowRequest) error {
	if e := workflowValidateClaudeTrust(r); e != nil {
		return e
	}
	g := r.CodexProjectTrust
	if g == nil {
		return nil
	}
	if r.Runtime.Provider != "codex" || g.Mode != "process-local" || !plain(g.RepositoryRoot, 1, 4096) || !filepath.IsAbs(g.RepositoryRoot) || filepath.Clean(g.RepositoryRoot) != g.RepositoryRoot {
		return workflowError(2, "codex_project_trust requires Codex, mode process-local and an absolute physical repository_root")
	}
	return nil
}

func workflowTrustStop(path, reason string) error {
	return workflowError(4, "Codex project trust cannot be confirmed at "+path+": "+reason+". Have the human inspect this path and policy before requesting a new preview; no trust was granted.")
}

func workflowTrustPathAt(path string) (workflowTrustPath, error) {
	f := workflowTrustPath{Path: path, Kind: "absent"}
	if e := physical(path, true); e != nil {
		return f, workflowTrustStop(path, "nonphysical or inaccessible search path")
	}
	st, e := os.Lstat(path)
	if os.IsNotExist(e) {
		return f, nil
	}
	if e != nil {
		return f, e
	}
	f.Device, f.Inode, e = workflowFileIdentity(st)
	if e != nil {
		return f, workflowTrustStop(path, "filesystem identity is unavailable")
	}
	switch {
	case st.IsDir():
		f.Kind = "directory"
	case st.Mode().IsRegular():
		f.Kind = "file"
		b, e := readFile(path, 4<<20, false)
		if e != nil {
			return f, e
		}
		f.SHA256 = hash(b)
	default:
		return f, workflowTrustStop(path, "unsupported file type")
	}
	return f, nil
}

// Only ordinary repositories and Git's standard linked-worktree layout are
// understood. In the latter Codex uses the main repository for shared trust.
// Verify both directions of the link, not just a guessed parent of the cwd.
func workflowTrustLayout(cwd string) (*workflowTrustFacts, error) {
	if e := physical(cwd, false); e != nil {
		return nil, workflowTrustStop(cwd, "cwd must be physical")
	}
	for _, key := range []string{"GIT_DIR", "GIT_WORK_TREE", "GIT_COMMON_DIR", "GIT_INDEX_FILE", "GIT_OBJECT_DIRECTORY", "GIT_ALTERNATE_OBJECT_DIRECTORIES"} {
		if os.Getenv(key) != "" {
			return nil, workflowTrustStop(cwd, "Git environment redirects the repository: "+key)
		}
	}
	git := func(arg string) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		c := exec.CommandContext(ctx, "git", "--no-optional-locks", "-c", "core.fsmonitor=false", "-c", "core.hooksPath=/dev/null", "-C", cwd, "rev-parse", "--path-format=absolute", arg)
		b, e := c.Output()
		if e != nil {
			return "", workflowTrustStop(cwd, "Git identity could not be observed")
		}
		s := strings.TrimSuffix(string(b), "\n")
		if !plain(s, 1, 4096) {
			return "", workflowTrustStop(cwd, "Git returned an unknown path")
		}
		return s, nil
	}
	wt, e := git("--show-toplevel")
	if e != nil {
		return nil, e
	}
	gd, e := git("--git-dir")
	if e != nil {
		return nil, e
	}
	common, e := git("--git-common-dir")
	if e != nil {
		return nil, e
	}
	root := filepath.Dir(common)
	if filepath.Base(common) != ".git" || common != filepath.Join(root, ".git") {
		return nil, workflowTrustStop(common, "nonstandard or bare Git layout")
	}
	f := &workflowTrustFacts{RepositoryRoot: root, WorktreeRoot: wt, CWD: cwd, GitDir: gd, CommonDir: common}
	paths := []string{cwd, wt, root, gd, common, filepath.Join(root, ".git/config")}
	if gd == common {
		if wt != root {
			return nil, workflowTrustStop(wt, "main worktree does not match the common directory")
		}
	} else {
		if filepath.Dir(gd) != filepath.Join(common, "worktrees") {
			return nil, workflowTrustStop(gd, "nonstandard linked-worktree metadata")
		}
		dotgit := filepath.Join(wt, ".git")
		forward, e := readFile(dotgit, 8192, false)
		if e != nil {
			return nil, e
		}
		if strings.TrimSpace(string(forward)) != "gitdir: "+gd {
			return nil, workflowTrustStop(dotgit, "linked-worktree pointer differs")
		}
		backPath := filepath.Join(gd, "gitdir")
		back, e := readFile(backPath, 8192, false)
		if e != nil {
			return nil, e
		}
		if strings.TrimSpace(string(back)) != dotgit {
			return nil, workflowTrustStop(backPath, "reverse worktree pointer differs")
		}
		commonPath := filepath.Join(gd, "commondir")
		b, e := readFile(commonPath, 8192, false)
		if e != nil {
			return nil, e
		}
		linked := strings.TrimSpace(string(b))
		if !filepath.IsAbs(linked) {
			linked = filepath.Join(gd, linked)
		}
		if filepath.Clean(linked) != common {
			return nil, workflowTrustStop(commonPath, "common directory pointer differs")
		}
		paths = append(paths, dotgit, backPath, commonPath, filepath.Dir(gd), filepath.Join(gd, "config.worktree"))
	}
	for _, p := range paths {
		v, e := workflowTrustPathAt(p)
		if e != nil {
			return nil, e
		}
		f.Paths = append(f.Paths, v)
	}
	return f, nil
}

func workflowObserveTrust(r WorkflowRequest, cwd string) (*workflowTrustFacts, error) {
	if e := workflowValidateTrust(r); e != nil {
		return nil, e
	}
	if r.CodexProjectTrust == nil {
		return nil, nil
	}
	f, e := workflowTrustLayout(cwd)
	if e != nil {
		return nil, e
	}
	if f.RepositoryRoot != r.CodexProjectTrust.RepositoryRoot {
		return nil, workflowTrustStop(r.CodexProjectTrust.RepositoryRoot, "grant differs from the observed Codex repository root "+f.RepositoryRoot)
	}
	home, e := os.UserHomeDir()
	if e != nil {
		return nil, e
	}
	f.UserHome = home
	f.CodexHome = filepath.Join(home, ".codex")
	if v := os.Getenv("CODEX_HOME"); v != "" {
		f.CodexHome = v
	}
	if e = physical(f.CodexHome, false); e != nil {
		return nil, workflowTrustStop(f.CodexHome, "Codex config home must be known and physical")
	}
	// Examine both physical ancestry chains. This also rejects a parent .codex
	// layer instead of silently assuming that a sibling linked worktree hides it.
	seen := map[string]bool{}
	for _, start := range []string{cwd, f.RepositoryRoot} {
		for p := start; ; p = filepath.Dir(p) {
			if !seen[p] {
				seen[p] = true
				for _, name := range []string{".codex", ".agents", "AGENTS.md", "AGENTS.override.md"} {
					path := filepath.Join(p, name)
					// The user's global layers are independent of project trust.
					if p == home {
						continue
					}
					v, e := workflowTrustPathAt(path)
					if e != nil {
						return nil, e
					}
					if v.Kind != "absent" {
						return nil, workflowTrustStop(path, "project config, hooks, rules or instructions could be activated")
					}
					f.Paths = append(f.Paths, v)
				}
			}
			if p == filepath.Dir(p) {
				break
			}
		}
	}
	if e = workflowTrustConfig(r, f); e != nil {
		return nil, e
	}
	f.PermissionProfile = r.Runtime.PermissionBinding.ProfileID
	f.EffectivePolicySHA256 = r.Runtime.PermissionBinding.EffectivePolicySHA256
	sort.Slice(f.Paths, func(i, j int) bool { return f.Paths[i].Path < f.Paths[j].Path })
	return f, nil
}

func workflowTrustConfig(r WorkflowRequest, f *workflowTrustFacts) error {
	etc, e := filepath.EvalSymlinks("/etc")
	if e != nil {
		return workflowTrustStop("/etc", "system config location is unknown")
	}
	paths := []string{filepath.Join(etc, "codex/config.toml"), filepath.Join(f.CodexHome, "config.toml")}
	if r.Runtime.ConfigProfile != nil {
		if !regexp.MustCompile(`^[A-Za-z0-9_-]+$`).MatchString(*r.Runtime.ConfigProfile) {
			return workflowTrustStop(f.CodexHome, "unsupported config profile name")
		}
		paths = append(paths, filepath.Join(f.CodexHome, *r.Runtime.ConfigProfile+".config.toml"))
	}
	paths = append(paths, filepath.Join(etc, "codex/managed_config.toml"), filepath.Join(etc, "codex/requirements.toml"), filepath.Join(f.CodexHome, "managed_config.toml"))
	relevant := func(p string) bool {
		return p == f.CWD || p == f.WorktreeRoot || p == f.RepositoryRoot || strings.HasPrefix(f.CWD, p+string(os.PathSeparator)) || strings.HasPrefix(f.RepositoryRoot, p+string(os.PathSeparator))
	}
	profileFound := strings.HasPrefix(r.Runtime.PermissionBinding.ProfileID, ":")
	for _, path := range paths {
		v, e := workflowTrustPathAt(path)
		if e != nil {
			return e
		}
		f.Paths = append(f.Paths, v)
		if v.Kind == "absent" {
			if r.Runtime.ConfigProfile != nil && path == filepath.Join(f.CodexHome, *r.Runtime.ConfigProfile+".config.toml") {
				return workflowTrustStop(path, "selected profile is missing")
			}
			continue
		}
		b, e := readFile(path, 4<<20, false)
		if e != nil {
			return e
		}
		var cfg map[string]any
		if e = toml.Unmarshal(b, &cfg); e != nil {
			return workflowTrustStop(path, "configuration is not valid TOML")
		}
		if markers, ok := cfg["project_root_markers"]; ok && !equal(markers, []string{".git"}) {
			return workflowTrustStop(path, "custom project root search is not supported for a new grant")
		}
		for p, entry := range workflowTOMLTable(cfg["projects"]) {
			if relevant(p) {
				level := workflowTOMLTable(entry)["trust_level"]
				if level == "untrusted" {
					return workflowTrustStop(path, "explicit untrusted project "+p)
				}
				if level != nil && level != "trusted" {
					return workflowTrustStop(path, "unknown project trust policy")
				}
			}
		}
		if workflowTOMLTable(cfg["permissions"])[r.Runtime.PermissionBinding.ProfileID] != nil {
			profileFound = true
		}
		if filepath.Base(path) == "requirements.toml" {
			if e := workflowTrustRequirements(r, path, cfg); e != nil {
				return e
			}
		}
		if filepath.Base(path) == "managed_config.toml" {
			for key, want := range map[string]string{"model": r.Runtime.Model, "default_permissions": r.Runtime.PermissionBinding.ProfileID, "approval_policy": "on-request", "approvals_reviewer": "auto_review"} {
				if val, ok := cfg[key]; ok && val != want {
					return workflowTrustStop(path, "managed default conflicts with bound "+key)
				}
			}
		}
	}
	if !profileFound {
		return workflowTrustStop(f.CodexHome, "bound permission profile is not defined by the observed configuration")
	}
	// MDM is a distinct high-priority source. A normal application preferences
	// file is not a managed policy. Inspect the two documented managed keys;
	// opaque managed payloads require a supported policy reader before a grant.
	for _, path := range []string{filepath.Join("/Library/Managed Preferences", "com.openai.codex.plist"), filepath.Join("/Library/Managed Preferences", filepath.Base(f.UserHome), "com.openai.codex.plist"), filepath.Join(f.UserHome, "Library/Preferences/com.openai.codex.plist")} {
		v, e := workflowTrustPathAt(path)
		if e != nil {
			return e
		}
		f.Paths = append(f.Paths, v)
		if v.Kind != "absent" {
			if e := workflowTrustPreferences(path); e != nil {
				return e
			}
		}
	}
	for _, path := range []string{filepath.Join(f.CodexHome, "hooks.json"), filepath.Join(f.CodexHome, "rules")} {
		v, e := workflowTrustPathAt(path)
		if e != nil {
			return e
		}
		f.Paths = append(f.Paths, v)
		if v.Kind == "directory" {
			e = filepath.WalkDir(path, func(p string, entry os.DirEntry, e error) error {
				if e != nil {
					return e
				}
				if p == path {
					return nil
				}
				v, e := workflowTrustPathAt(p)
				if e != nil {
					return e
				}
				f.Paths = append(f.Paths, v)
				return nil
			})
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func workflowTrustPreferences(path string) error {
	b, e := readFile(path, 4<<20, false)
	if e != nil {
		return e
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, "/usr/bin/plutil", "-convert", "xml1", "-o", "-", "--", "-")
	c.Stdin = bytes.NewReader(b)
	var out workflowLimitedOutput
	c.Stdout = &out
	if e = c.Run(); e != nil || out.overflow {
		return workflowTrustStop(path, "managed preferences cannot be decoded")
	}
	d := xml.NewDecoder(bytes.NewReader(out.data))
	for {
		token, e := d.Token()
		if e == io.EOF {
			break
		}
		if e != nil {
			return workflowTrustStop(path, "invalid preferences XML")
		}
		if start, ok := token.(xml.StartElement); ok && start.Name.Local == "key" {
			var key string
			if e = d.DecodeElement(&key, &start); e != nil {
				return e
			}
			if key == "config_toml_base64" || key == "requirements_toml_base64" {
				return workflowTrustStop(path, "managed preferences require policy inspection; effective constraints are unknown")
			}
		}
	}
	return nil
}
func workflowTOMLTable(v any) map[string]any { m, _ := v.(map[string]any); return m }

func workflowTrustFresh(s workflowState) error {
	if s.Request.CodexProjectTrust == nil {
		if s.CodexTrust != nil {
			return workflowError(4, "unexpected saved Codex trust binding")
		}
		return nil
	}
	fresh, e := workflowObserveTrust(s.Request, s.Observed.Target.WorktreeLocator)
	if e != nil {
		return e
	}
	if s.CodexTrust == nil || !equal(fresh, s.CodexTrust) {
		return workflowError(4, "Codex trust, repository identity or configuration changed since confirmation; no new effect is authorized")
	}
	return nil
}

func workflowTrustRequirements(r WorkflowRequest, path string, cfg map[string]any) error {
	if values, ok := cfg["allowed_permission_profiles"]; ok {
		if workflowTOMLTable(values)[r.Runtime.PermissionBinding.ProfileID] != true {
			return workflowTrustStop(path, "managed policy denies allowed_permission_profiles="+r.Runtime.PermissionBinding.ProfileID)
		}
	}
	for key, want := range map[string]string{"allowed_approval_policies": "on-request", "allowed_approvals_reviewers": "auto_review"} {
		if values, ok := cfg[key]; ok {
			allowed := false
			if a, ok := values.([]any); ok {
				for _, v := range a {
					if v == want {
						allowed = true
					}
				}
			}
			if !allowed {
				return workflowTrustStop(path, "managed policy denies "+key+"="+want)
			}
		}
	}

	return nil
}
