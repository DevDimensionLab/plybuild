package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/workflowtrace"
	"strings"
	"testing"
)

func TestReviewDiagramPreservesRepeatedSourcePositions(t *testing.T) {
	r := workflowTraceFixture(t, "empty")
	r.HistoryState = "recorded"
	r.Events = []workflowtrace.Event{
		{ID: "first", Type: "delivery_report", Title: "First", Role: "agent", EvidenceClass: "reported", Positions: []workflowtrace.Position{{ChainID: "one:delivery", Sequence: 1}, {ChainID: "one:delivery", Sequence: 3}}},
		{ID: "middle", Type: "delivery_report", Title: "Middle", Role: "agent", EvidenceClass: "reported", Positions: []workflowtrace.Position{{ChainID: "one:delivery", Sequence: 2}}},
	}
	r.Chains = []workflowtrace.Chain{{ID: "one:delivery", Ordering: "source sequence", EventIDs: []string{"first", "middle", "first"}}}
	got := workflowTraceText(r)
	if arrows := strings.Count(got, "v  recorded order"); arrows != 2 {
		t.Fatalf("source positions 1,2,3 require two recorded-order arrows; got %d:\n%s", arrows, got)
	}
	if strings.Contains(got, "Same exact event shown in another chain") {
		t.Fatal("same-chain repeated position falsely described as another chain")
	}
}
