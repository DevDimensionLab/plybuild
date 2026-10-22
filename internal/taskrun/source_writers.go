package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// ObserveSourceWriters uses the same managed inventories and bounded cwd table
// as delivery recovery, without inventing a session for an unregistered source.
// Tests inject the inventory transport and LocalProcesses; production calls
// only the native read-only pane/agent lists, never a provider.
func ObserveSourceWriters(d Dependencies, path string, inventory func(string) (json.RawMessage, error)) (workspace.SourceWriterObservation, error) {
	out := workspace.SourceWriterObservation{State: "unknown", Method: writerLivenessMethod, Limitation: "Bounded managed-pane/agent and local cwd observations. External writers that leave this cwd and have no known session may be unobservable."}
	if inventory == nil {
		executable, e := exec.LookPath("herdr")
		if e != nil {
			return out, e
		}
		executable, e = filepath.EvalSymlinks(executable)
		if e != nil {
			return out, e
		}
		b, e := os.ReadFile(executable)
		if e != nil {
			return out, e
		}
		r := WorkflowRequest{}
		r.Herdr.Executable = Executable{Path: executable, SHA256: hash(b)}
		inventory = func(kind string) (json.RawMessage, error) { return workflowCall(d, r, kind, "list") }
	}
	for _, kind := range []string{"pane", "agent"} {
		raw, e := inventory(kind)
		if e != nil {
			return out, e
		}
		var rows struct {
			Panes  *[]workflowAgent `json:"panes"`
			Agents *[]workflowAgent `json:"agents"`
		}
		if e = json.Unmarshal(raw, &rows); e != nil {
			return out, e
		}
		v := rows.Panes
		if kind == "agent" {
			v = rows.Agents
		}
		if v == nil {
			return out, fmt.Errorf("incomplete managed %s inventory", kind)
		}
		seen := map[string]bool{}
		for _, a := range *v {
			if a.PaneID == "" || a.TerminalID == "" || a.WorkspaceID == "" || a.TabID == "" || seen[a.PaneID] {
				return out, fmt.Errorf("ambiguous managed identity")
			}
			seen[a.PaneID] = true
			if sourceContains(path, a.CWD) || sourceContains(path, a.ForegroundCWD) {
				out.State = "present"
				return out, fmt.Errorf("managed source writer remains in pane %s", a.PaneID)
			}
		}
	}
	local, e := observePathWriterLiveness(d, path, "", "", 0)
	if e != nil {
		return out, e
	}
	if len(local.RemainingPIDs) > 0 {
		out.State = "present"
		return out, fmt.Errorf("source writer processes remain: %v", local.RemainingPIDs)
	}
	out.State = "absent"
	return out, nil
}
