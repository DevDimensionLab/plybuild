package taskrun

import (
	"strings"
	"testing"
)

func al26Report() DeliveryReport {
	return DeliveryReport{Envelope: deliveryEnv("delivery-report"), EventID: "progress", Phase: "working", Summary: "Focused checks observed.", Meaning: "Progress only; no candidate qualification.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{{ID: "focused-test", Outcome: "not_run", Reason: "Not executed yet."}}}
}

func TestDeliveryReportVerifierFieldDiagnostics(t *testing.T) {
	for _, tc := range []struct {
		name, id, reason, want string
		duplicate              bool
	}{
		{"id", "TestActualName", "Not executed yet.", "verifier_results[0].id", false},
		{"duplicate", "focused-test", "Not executed yet.", "verifier_results[1].id", true},
		{"missing_reason", "focused-test", "", "verifier_results[0].reason", false},
		{"padded_reason", "focused-test", " padded ", "verifier_results[0].reason", false},
		{"control_reason", "focused-test", "line\nbreak", "verifier_results[0].reason", false},
		{"long_reason", "focused-test", strings.Repeat("x", 2001), "verifier_results[0].reason", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := al26Report()
			r.VerifierResults[0].ID = tc.id
			r.VerifierResults[0].Reason = tc.reason
			if tc.duplicate {
				r.VerifierResults = append(r.VerifierResults, r.VerifierResults[0])
			}
			err := validateDeliveryReportContent(r)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected actionable %s diagnostic, got %v", tc.want, err)
			}
			if !strings.Contains(err.Error(), tc.id) {
				t.Fatalf("diagnostic lacks declaration ID: %v", err)
			}
			correction := "1..2000 Unicode codepoints"
			if tc.name == "id" {
				correction = "Put actual test names in argv/evidence"
			} else if tc.duplicate {
				correction = "duplicates verifier_results[0].id; use a unique declaration key"
			}
			if !strings.Contains(err.Error(), correction) || strings.Contains(err.Error(), r.Meaning) {
				t.Fatalf("diagnostic lacks correction or dumps report text: %v", err)
			}
		})
	}
}

func TestDeliveryReportVerifierOutcomeBoundaries(t *testing.T) {
	r := al26Report()
	r.VerifierResults = nil
	if err := validateDeliveryReportContent(r); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		outcome string
		exit    *int
		argv    []string
		cwd     string
		valid   bool
	}{
		{"not_run", nil, nil, "", true}, {"unknown", nil, nil, "", true},
		{"failed", ptr(1), []string{"go", "test", "-run", "^TestActualName$"}, "/tmp", true},
		{"passed", ptr(0), []string{"go", "test"}, "/tmp", true},
		{"not_run", ptr(0), nil, "", false}, {"unknown", ptr(1), nil, "", false},
		{"passed", ptr(1), []string{"go"}, "/tmp", false}, {"failed", ptr(0), []string{"go"}, "/tmp", false},
		{"passed", nil, []string{"go"}, "/tmp", false}, {"passed", ptr(0), nil, "/tmp", false},
		{"passed", ptr(0), []string{"go"}, "relative", false},
	} {
		r = al26Report()
		v := &r.VerifierResults[0]
		v.Outcome, v.Exit, v.Argv, v.CWD = tc.outcome, tc.exit, tc.argv, tc.cwd
		if err := validateDeliveryReportContent(r); (err == nil) != tc.valid {
			t.Fatalf("%+v: %v", tc, err)
		}
	}
}
