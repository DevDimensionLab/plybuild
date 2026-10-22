package taskrun

import (
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"strings"
	"testing"
)

func TestAL32GeneratedReviewReturnLimits(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	reviewRoundCandidate(t, f, "fixed guidance candidate")
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	// Generate the actual reviewer inputs without launching a review or inspecting processes.
	brief := deliveryReviewBrief(s, deliveryReviewRound{Round: 1}, nil, f.R.WorkspaceRoot, "/fixture/ply", "/fixture/template.json", "/fixture/context.json", workspace.WorkspaceTaskPreparePlan{})
	for _, want := range []string{"summary: 1-2000", "findings[].summary: 1-600", "Unicode codepoints", "one line", "256 KiB", "bytes"} {
		if !strings.Contains(brief, want) {
			t.Errorf("generated brief lacks %q", want)
		}
	}
	raw, err := json.Marshal(deliveryReviewTemplate())
	if err != nil {
		t.Fatal(err)
	}
	var input map[string]any
	if err = json.Unmarshal(raw, &input); err != nil {
		t.Fatal(err)
	}
	input["decision"] = "changes_requested"
	input["summary"] = strings.Repeat("界", 2000)
	finding := input["findings"].([]any)[0].(map[string]any)
	finding["severity"], finding["summary"], finding["evidence"] = "medium", strings.Repeat("ø", 600), al13Prose()
	raw, err = json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	got, err := deliveryParseReviewReturn(raw, 1)
	if err != nil {
		t.Fatalf("documented boundaries rejected: %v", err)
	}
	if got.Summary != input["summary"] || len(got.Findings) != 1 || got.Findings[0].Summary != finding["summary"] || got.Findings[0].Evidence != al13Prose() {
		t.Fatalf("return changed short fields or explanatory prose: %+v", got)
	}
	// Existing field-specific diagnostics still reject overflow and multiline summaries.
	for _, tc := range []struct{ field, value, diagnostic string }{
		{"summary", strings.Repeat("界", 2001), "summary must be 1-2000"},
		{"summary", "Two\nlines", "summary must be 1-2000"},
		{"finding", strings.Repeat("ø", 601), "findings[0]"},
		{"finding", "Two\nlines", "findings[0]"},
	} {
		input["summary"], finding["summary"] = "Complete review", "Concrete defect"
		if tc.field == "summary" {
			input["summary"] = tc.value
		} else {
			finding["summary"] = tc.value
		}
		raw, err = json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = deliveryParseReviewReturn(raw, 1); err == nil || !strings.Contains(err.Error(), tc.diagnostic) {
			t.Fatalf("%s diagnostic: %v", tc.field, err)
		}
	}
	// The native return reader applies a separate byte budget before parsing.
	file := reviewRoundReturnFile(t, f, "oversized.json", string(raw)+strings.Repeat(" ", 256<<10))
	if _, err = readFile(file, 256<<10, false); err == nil || !strings.Contains(err.Error(), "at most 262144 bytes") {
		t.Fatalf("total file limit diagnostic: %v", err)
	}
}
