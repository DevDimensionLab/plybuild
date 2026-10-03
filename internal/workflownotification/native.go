package workflownotification

import (
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

func checkNative(s Source) error {
	v, err := workflowhandoff.Inspect(workflowhandoff.SystemDependencies(), workflowhandoff.InspectInput{HandoffLocator: s.Handoff.Path, Format: "json"})
	if err != nil {
		return fail(2, "native_binding", "Native handoff integrity or binding is invalid.")
	}
	var r struct {
		Integrity struct {
			Valid bool `json:"valid"`
		} `json:"integrity"`
		Identities struct {
			Activity string `json:"activity_id"`
			Run      string `json:"run_id"`
			Handoff  string `json:"handoff_id"`
		} `json:"identities"`
		Documents map[string]*struct {
			Locator string `json:"locator"`
			SHA256  string `json:"sha256"`
			Content struct {
				SchemaVersion int `json:"schema_version"`
				Target        struct {
					Worktree string `json:"worktree"`
				} `json:"target_binding"`
			} `json:"content"`
		} `json:"documents"`
		Validations map[string]struct {
			Checked bool `json:"checked"`
			Valid   bool `json:"valid"`
		} `json:"validations"`
		Basis *struct {
			TaskID string `json:"task_id"`
		} `json:"task_spec_basis"`
	}
	if json.Unmarshal(v.Bytes, &r) != nil || !r.Integrity.Valid || r.Identities.Activity != s.Activity || r.Identities.Run != s.Run || s.HandoffID == nil || r.Identities.Handoff != *s.HandoffID {
		return fail(2, "native_binding", "Native handoff integrity or binding is invalid.")
	}
	for name, loc := range map[string]Locator{"handoff": s.Handoff, "start_receipt": s.Start, "terminal_result": s.Report} {
		d := r.Documents[name]
		if d == nil || d.Locator != loc.Path || d.SHA256 != loc.SHA256 {
			return fail(2, "native_binding", "Native document locators or digests do not match.")
		}
	}
	if r.Documents["handoff"].Content.Target.Worktree != s.Worktree {
		return invalid()
	}
	for _, name := range []string{"schema", "digest", "lifecycle", "binding", "capability", "principal_session"} {
		v, ok := r.Validations[name]
		if !ok || !v.Checked || !v.Valid {
			return invalid()
		}
	}
	var taskID *string
	if json.Unmarshal(s.TaskID, &taskID) != nil {
		return invalid()
	}
	if r.Documents["handoff"].Content.SchemaVersion == 2 {
		val := r.Validations["task_spec"]
		if !val.Checked || !val.Valid || r.Basis == nil && taskID != nil || r.Basis != nil && (taskID == nil || *taskID != r.Basis.TaskID) {
			return invalid()
		}
	} else if taskID != nil {
		return invalid()
	}
	// policy/evidence gaps and a reported failure are deliberately not product judgements.
	return nil
}
