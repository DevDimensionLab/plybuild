package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// Opt-in, read-only validation of real frozen requests. This calls the normal
// request validator, never startup, receipt, verification or delivery callbacks.
func TestInstallAuthorityHistoricalFrozenRequests(t *testing.T) {
	paths := filepath.SplitList(os.Getenv("PLY_INSTALL_HISTORICAL_REQUESTS"))
	if len(paths) == 0 {
		t.Skip("supply preserved historical request files for read-only validation")
	}
	for _, path := range paths {
		t.Run(filepath.Base(filepath.Dir(path)), func(t *testing.T) {
			before, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			request, err := ReadWorkflowRequest(path)
			if err != nil {
				t.Fatal(err)
			}
			canonical, err := Canonical(request)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(before, canonical) {
				t.Fatal("historical request canonical bytes changed")
			}
			if request.Delivery == nil || !request.Delivery.AllowLocalInstall {
				t.Fatal("historical installation authority was not preserved")
			}
			goal, err := workflowBound(request.Delivery.Goal, 4<<20)
			if err != nil {
				t.Fatal(err)
			}
			var fields struct {
				Executor map[string]json.RawMessage `json:"executor"`
			}
			if err := json.Unmarshal(goal, &fields); err != nil {
				t.Fatal(err)
			}
			if _, present := fields.Executor["allow_local_install"]; present {
				t.Fatal("expected a historical Goal with omitted installation selection")
			}
			after, err := os.ReadFile(path)
			if err != nil || !bytes.Equal(before, after) {
				t.Fatalf("historical request changed: %v", err)
			}
			t.Logf("Frozen omitted-field request retained install=true and passed normal validation unchanged: %s digest %s", path, hash(before))
		})
	}
}
