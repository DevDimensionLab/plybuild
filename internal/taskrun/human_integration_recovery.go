package taskrun

import (
	"path/filepath"
	"strings"
)

type HumanOwnerRecoveryObservation struct {
	Transport WorkflowTransport `json:"transport"`
	// InventoryObservedAt dates the completed Herdr inventory inspection only.
	// Transport retains its historical state/time; old receipts omit this field.
	InventoryObservedAt string `json:"inventory_observed_at,omitempty"`
	PaneInventorySHA    string `json:"pane_inventory_sha256"`
	AgentInventorySHA   string `json:"agent_inventory_sha256"`
	ProcessLiveness     string `json:"process_liveness"`
}

func sourceContains(source, path string) bool {
	if path == "" {
		return false
	}
	if physical, err := filepath.EvalSymlinks(path); err == nil {
		path = physical
	}
	path = filepath.Clean(path)
	return path == source || strings.HasPrefix(path, source+string(filepath.Separator))
}

// The confirmation binds only this candidate and its native history. Complete
// inventories are rechecked after input, without requiring unrelated panes to
// remain byte-identical while the human reads the prompt.
