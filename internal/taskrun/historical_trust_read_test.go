package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestHistoricalTrustFieldsPreserveReservationOwnership(t *testing.T) {
	for _, field := range []string{"claude_project_trust", "codex_project_trust"} {
		t.Run(field, func(t *testing.T) {
			f := workflowTestFixture(t)
			p, err := WorkflowPreviewStart(f.D, f.File)
			if err != nil {
				t.Fatal(err)
			}
			o, err := WorkflowStart(f.D, f.File, *p.(WorkflowPreview).Confirmation)
			if err != nil {
				t.Fatal(err)
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			paths := []string{workflowIndex(f.R.WorkspaceRoot, o.RunID), filepath.Join(s.Result.Paths.RunRoot, "state.json")}
			snapshots := map[string][]byte{}
			for _, path := range paths {
				b, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				var v map[string]any
				if err := json.Unmarshal(b, &v); err != nil {
					t.Fatal(err)
				}
				request := v["request"].(map[string]any)
				v["result"].(map[string]any)[field] = map[string]any{"historical_effect": "retained"}
				v[field] = map[string]any{"historical_observation": "retained"}
				request[field] = map[string]any{"mode": "configuration", "historical_evidence": "retained"}
				v["result"].(map[string]any)["request_sha256"] = digest(request)
				b, err = Canonical(v)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, b, 0600); err != nil {
					t.Fatal(err)
				}
				snapshots[path] = b
			}
			other := workspace.PlanWorktreeObservation{WorktreeLocator: filepath.Join(f.R.WorkspaceRoot, "other-task"), GitCommonDir: s.Observed.Target.GitCommonDir, Ref: "refs/heads/other-task"}
			if err := workflowReservations(f.D, f.R.WorkspaceRoot, other); err != nil {
				t.Fatalf("historical trust must not block unrelated target: %v", err)
			}
			if err := workflowReservations(f.D, f.R.WorkspaceRoot, s.Observed.Target); err == nil {
				t.Fatal("historical owner was released")
			}
			for path, before := range snapshots {
				after, err := os.ReadFile(path)
				if err != nil || !bytes.Equal(before, after) {
					t.Fatalf("history changed: %s %v", path, err)
				}
			}
		})
	}
}
