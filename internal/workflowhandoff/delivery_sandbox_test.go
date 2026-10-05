package workflowhandoff

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeliveryRuntimeFilesystemRootCoverage(t *testing.T) {
	f := prepareDeliveryFixture(t)
	for _, tc := range []struct {
		name, readRoot, writeRoot, wantError string
	}{
		{"filesystem root covers the bound fixture", "/", "/", ""},
		{"task-only reads do not cover handoff", f.target, "/", "sandbox read roots do not cover"},
		{"task-only writes do not cover reply", "/", f.target, "sandbox write roots do not cover"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sandbox, err := json.Marshal(map[string]any{
				"read_roots": []string{tc.readRoot}, "write_roots": []string{tc.writeRoot},
				"temp_root": filepath.Join(f.root, "recipient-temp"), "matches_contract": true,
			})
			if err != nil {
				t.Fatal(err)
			}
			start, err := BuildDeliveryTaskRunStart(f.d, f.locator, "fixture human, not actual approval", "synthetic test", "fixture-owner-session", "codex", "fixture-model", sandbox, []byte("[]"), "started")
			if err == nil {
				err = ValidateTaskRunStart(f.d, f.locator, start)
			}
			if tc.wantError == "" {
				if err != nil {
					t.Fatalf("reported filesystem root should cover its descendants: %v", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tc.wantError) {
				t.Fatalf("want %q, got %v", tc.wantError, err)
			}
		})
	}
}
