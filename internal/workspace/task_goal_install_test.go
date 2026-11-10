package workspace

import (
	"encoding/json"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func TestGoalInstallSelectionStrictAndHistoricalBytes(t *testing.T) {
	rule := taskGoalSpecSchema(false)["executor"]
	for _, tc := range []struct {
		name, raw string
		valid     bool
	}{
		{"historical", `{"provider":"codex","model":null,"effort":null}`, true},
		{"reserved", `{"provider":"codex","model":null,"effort":null,"allow_local_install":false}`, true},
		{"allowed", `{"provider":"codex","model":null,"effort":null,"allow_local_install":true}`, true},
		{"null", `{"provider":"codex","model":null,"effort":null,"allow_local_install":null}`, false},
		{"text", `{"provider":"codex","model":null,"effort":null,"allow_local_install":"true"}`, false},
		{"unknown", `{"provider":"codex","model":null,"effort":null,"install":true}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			value, err := canonicaljson.DecodeStrict([]byte(tc.raw))
			if err != nil {
				t.Fatal(err)
			}
			if err := rule(value); (err == nil) != tc.valid {
				t.Fatalf("validation = %v, want valid %t", err, tc.valid)
			}
			if !tc.valid {
				return
			}
			var assignment TaskExecutorAssignment
			if err := json.Unmarshal([]byte(tc.raw), &assignment); err != nil {
				t.Fatal(err)
			}
			actual, err := contentCanonical(assignment)
			if err != nil {
				t.Fatal(err)
			}
			expected, err := canonicaljson.Marshal(value)
			if err != nil {
				t.Fatal(err)
			}
			if string(actual) != string(expected) {
				t.Fatalf("canonical assignment bytes changed: %s != %s", actual, expected)
			}
		})
	}
}
