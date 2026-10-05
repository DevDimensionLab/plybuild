package taskexecute

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/pelletier/go-toml/v2"
)

// RuntimeOptions names existing executables and an optional existing Codex
// profile. These are launch choices, never changes to permission configuration.
type RuntimeOptions struct {
	ProviderExecutable string
	HerdrExecutable    string
	PermissionProfile  string
	CodexHome          string
	HerdrWorkspace     string
}

type LaunchPolicy struct {
	Kind              string                `json:"kind"`
	SchemaVersion     int                   `json:"schema_version"`
	Provider          string                `json:"provider"`
	PermissionProfile string                `json:"permission_profile"`
	ApprovalPolicy    string                `json:"approval_policy"`
	ApprovalReviewer  string                `json:"approval_reviewer"`
	Sources           []taskrun.FileBinding `json:"sources"`
	Assurance         string                `json:"assurance"`
}

type RuntimePreview struct {
	Runtime         taskrun.Runtime    `json:"runtime"`
	ReasoningEffort string             `json:"reasoning_effort"`
	Herdr           taskrun.Executable `json:"herdr"`
	HerdrWorkspace  string             `json:"herdr_workspace"`
	Policy          LaunchPolicy       `json:"launch_policy"`
}

func stringChoice(v *string, fallback string) string {
	if v == nil {
		return fallback
	}
	return *v
}

func lookupExecutable(choice, fallback string) (taskrun.Executable, error) {
	if choice == "" {
		var err error
		choice, err = exec.LookPath(fallback)
		if err != nil {
			return taskrun.Executable{}, fmt.Errorf("required %s binary is missing: %w", fallback, err)
		}
	}
	return bindExecutable(choice)
}

// PreviewRuntime reads only local settings and executable bytes. Its policy is
// explicitly a requested launch contract. Recipient acceptance must supply the
// actual runtime policy; a config digest is not evidence of an OS write grant.
func PreviewRuntime(assignment workspace.TaskExecutorAssignment, control string, opts RuntimeOptions) (RuntimePreview, error) {
	var out RuntimePreview
	provider := stringChoice(assignment.Provider, "codex")
	if provider != "codex" && provider != "claude" {
		return out, fmt.Errorf("unsupported implementor %q; choose codex or claude", provider)
	}
	out.Policy = LaunchPolicy{
		Kind: "ply.workflow.launch-policy", SchemaVersion: 1, Provider: provider,
		Sources:   []taskrun.FileBinding{},
		Assurance: "Requested launch contract only; recipient must confirm actual permissions before Task writes. Runtime enforcement remains authoritative.",
	}
	model, effort, profile := "opus", "medium", "auto"
	if provider == "codex" {
		codexHome := opts.CodexHome
		if codexHome == "" {
			codexHome = os.Getenv("CODEX_HOME")
		}
		if codexHome == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return out, err
			}
			codexHome = filepath.Join(home, ".codex")
		}
		path := filepath.Join(codexHome, "config.toml")
		b, err := os.ReadFile(path)
		if err != nil {
			return out, fmt.Errorf("read existing Codex settings: %w", err)
		}
		if len(b) > 4<<20 {
			return out, fmt.Errorf("Codex settings exceed 4 MiB")
		}
		var cfg map[string]any
		if err = toml.Unmarshal(b, &cfg); err != nil {
			return out, fmt.Errorf("read Codex settings: %w", err)
		}
		model, _ = cfg["model"].(string)
		effort, _ = cfg["model_reasoning_effort"].(string)
		profile, _ = cfg["default_permissions"].(string)
		if opts.PermissionProfile != "" {
			profile = opts.PermissionProfile
		}
		if profile == "" {
			return out, fmt.Errorf("Codex has no selected permission profile; configure it in Codex or name an existing profile with --permission-profile")
		}
		permissions, _ := cfg["permissions"].(map[string]any)
		if permissions[profile] == nil && !strings.HasPrefix(profile, ":") {
			return out, fmt.Errorf("Codex permission profile %q is not defined in %s", profile, path)
		}
		physical, err := filepath.EvalSymlinks(path)
		if err != nil {
			return out, err
		}
		out.Policy.Sources = append(out.Policy.Sources, taskrun.FileBinding{Locator: physical, SHA256: digestBytes(b)})
		out.Policy.ApprovalPolicy = "on-request"
		out.Policy.ApprovalReviewer = "auto_review"
	} else if opts.PermissionProfile != "" && opts.PermissionProfile != "auto" {
		return out, fmt.Errorf("new Claude executions use native auto permission mode; omit --permission-profile or use auto")
	}
	model = stringChoice(assignment.Model, model)
	effort = stringChoice(assignment.Effort, effort)
	if model == "" || strings.TrimSpace(model) != model {
		return out, fmt.Errorf("the implementor model is unknown; record a model in the goal or existing provider settings")
	}
	if effort != "" && effort != "low" && effort != "medium" && effort != "high" && effort != "xhigh" && effort != "max" && effort != "minimal" && effort != "none" {
		return out, fmt.Errorf("unsupported reasoning effort %q", effort)
	}
	if provider == "claude" && (effort == "minimal" || effort == "none") {
		return out, fmt.Errorf("Claude does not support reasoning effort %q", effort)
	}
	var err error
	providerBin, err := lookupExecutable(opts.ProviderExecutable, provider)
	if err != nil {
		return out, err
	}
	controlBin, err := bindExecutable(control)
	if err != nil {
		return out, err
	}
	out.Herdr, err = lookupExecutable(opts.HerdrExecutable, "herdr")
	if err != nil {
		return out, err
	}
	out.HerdrWorkspace = opts.HerdrWorkspace
	if out.HerdrWorkspace == "" {
		out.HerdrWorkspace = os.Getenv("HERDR_WORKSPACE_ID")
	}
	if out.HerdrWorkspace == "" {
		return out, fmt.Errorf("Herdr workspace is unknown; run from a Herdr terminal or provide --herdr-workspace")
	}
	out.Policy.PermissionProfile = profile
	out.ReasoningEffort = effort
	out.Runtime = taskrun.Runtime{Provider: provider, Mode: "interactive", Model: model, Executable: providerBin, PlyExecutable: controlBin,
		PermissionBinding: taskrun.Permission{AuthorityKind: "launch_contract_pending_runtime_acceptance", ProfileID: profile, Evidence: []taskrun.Evidence{}},
	}
	return out, nil
}

func preserveRuntime(preview RuntimePreview, directory string) (RuntimePreview, error) {
	for _, source := range preview.Policy.Sources {
		b, err := os.ReadFile(source.Locator)
		if err != nil || digestBytes(b) != source.SHA256 {
			return preview, fmt.Errorf("provider configuration changed before launch: %s", source.Locator)
		}
	}
	policy, err := taskrun.Canonical(preview.Policy)
	if err != nil {
		return preview, err
	}
	policyPath := filepath.Join(directory, "launch-policy.json")
	if err = writeOnce(policyPath, policy, 0600); err != nil {
		return preview, err
	}
	preview.Runtime.PermissionBinding.EffectivePolicySHA256 = digestBytes(policy)
	preview.Runtime.PermissionBinding.Evidence = []taskrun.Evidence{{Locator: policyPath, SHA256: digestBytes(policy), Role: "runtime_contract"}}
	preview.Runtime.PlyExecutable, err = PreserveControl(preview.Runtime.PlyExecutable, directory)
	return preview, err
}
