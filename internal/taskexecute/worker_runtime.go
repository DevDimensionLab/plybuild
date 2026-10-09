package taskexecute

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func mustWorkerJSON(w *taskrun.WorkerContract) []byte { b, _ := json.Marshal(w); return b }
func workerTabLabel(w *taskrun.WorkerContract, title string) string {
	if w != nil {
		return ""
	}
	return tabLabel(title)
}

func previewWorkerRuntime(a workspace.TaskExecutorAssignment, control string, opts RuntimeOptions) (RuntimePreview, error) {
	var out RuntimePreview
	if runtime.GOOS != "darwin" {
		return out, fmt.Errorf("headless worker currently requires the verified macOS sandbox adapter")
	}
	if stringChoice(a.Provider, "codex") != "codex" || opts.PermissionProfile != "" || opts.HerdrWorkspace != "" {
		return out, fmt.Errorf("headless uses Codex and its explicit Task-file sandbox; no Herdr or interactive permission override")
	}
	provider, err := lookupExecutable(opts.ProviderExecutable, "codex")
	if err != nil {
		return out, err
	}
	bin, err := bindExecutable(control)
	if err != nil {
		return out, err
	}
	out.ReasoningEffort = stringChoice(a.Effort, "max")
	out.Runtime = taskrun.Runtime{Provider: "codex", Mode: "headless", Model: stringChoice(a.Model, "gpt-6-astra"), Executable: provider, PlyExecutable: bin,
		PermissionBinding: taskrun.Permission{AuthorityKind: "launch_contract_pending_runtime_acceptance", ProfileID: "ply-worker", Evidence: []taskrun.Evidence{}}}
	out.Policy = LaunchPolicy{Kind: "ply.workflow.worker-launch-policy", SchemaVersion: 1, Provider: "codex", PermissionProfile: "ply-worker", ApprovalPolicy: "never", Sources: []taskrun.FileBinding{}, Assurance: "Requested Task-file writes only; Git, native state and controller remain read-only to agents. A native sandbox canary must pass before development. The worker owns commits and exact native local delivery."}
	return out, nil
}
