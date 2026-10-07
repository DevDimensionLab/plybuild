package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/workflowtrace"
	"github.com/mattn/go-runewidth"
)

func traceDiagramDelivery(t *testing.T) workflowtrace.Result {
	t.Helper()
	r := workflowTraceFixture(t, "empty")
	r.HistoryState = "recorded"
	chain := "wfr_" + strings.Repeat("a", 64) + ":delivery"
	candidate, result := "candidate_"+strings.Repeat("b", 64), "trs_"+strings.Repeat("c", 32)
	r.Candidates = []workflowtrace.Candidate{{ID: candidate, OID: strings.Repeat("d", 40), Tree: strings.Repeat("e", 40), ResultIDs: []string{result}}}
	question := "Can this installed candidate pass the product check?"
	roles := []string{"agent", "ply", "agent", "agent", "human", "ply"}
	types := []string{"delivery_report", "verification", "delivery_report", "delivery_report", "human_qa", "integration"}
	outcomes := []string{"working", "passed", "needs_input", "needs_input", "pass", "exact_effect"}
	for i := range types {
		data := map[string]any{"summary": "Recorded work", "raw_only_marker": "evidence must remain available"}
		if i == 2 || i == 3 {
			data["question"] = question
		}
		payload, _ := json.Marshal(data)
		e := workflowtrace.Event{ID: fmt.Sprintf("%s/event-%d", chain, i+1), Type: types[i], Title: types[i], Role: roles[i], Outcome: &outcomes[i], EvidenceClass: "reported", Data: payload, Positions: []workflowtrace.Position{{ChainID: chain, Sequence: i + 1}}}
		if i == 1 || i >= 4 {
			e.CandidateID, e.ResultID = &candidate, &result
		}
		r.Events = append(r.Events, e)
	}
	ids := []string{}
	for _, e := range r.Events {
		ids = append(ids, e.ID)
	}
	r.Chains = []workflowtrace.Chain{{ID: chain, Ordering: "source sequence", EventIDs: ids}}
	return r
}

func TestWorkflowTraceDiagramKeepsPrintableUnicodeAndEscapesUnsafeText(t *testing.T) {
	r := traceDiagramDelivery(t)
	r.Events[2].Data = json.RawMessage(`{"question":"Når du er tilbake: kjør café 🧭 界 é.\n\u001b[31m\u202ereversed\u2066direction"}`)
	text := workflowTraceText(r)
	for _, printable := range []string{"Når du er tilbake: kjør", "café", "🧭", "界", "é."} {
		if !strings.Contains(text, printable) {
			t.Fatalf("printable user text %q must be legible in boxes:\n%s", printable, text)
		}
	}
	if !utf8.ValidString(text) || strings.ContainsAny(text, "\x1b\u202e\u2066") || !strings.Contains(text, `\n`) || !strings.Contains(text, `\u202e`) {
		t.Fatal("unsafe controls or directional characters were not escaped")
	}
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, " | ") && strings.HasSuffix(line, " |") || strings.HasPrefix(line, "         +") {
			if width := runewidth.StringWidth(line); width != 75 {
				t.Fatalf("Unicode box width=%d, want75: %q", width, line)
			}
		}
	}
}

func TestWorkflowTraceDiagramRoundsSubsecondDurationsWithoutNoise(t *testing.T) {
	for input, want := range map[float64]string{87.29266: "1m27.293s", 59.9999: "1m0s", 22367: "6h12m47s", 0.0004: "0s"} {
		if got := traceDuration(input); got != want {
			t.Errorf("%g: got %q want %q", input, got, want)
		}
	}
}

func TestWorkflowTraceDiagramShowsReadableReportedDuration(t *testing.T) {
	r := traceDiagramDelivery(t)
	seconds := float64(22367)
	r.Analysis = append(r.Analysis, workflowtrace.Analysis{ID: "reported_qa/test", Kind: "interval", Value: &seconds, Coverage: "exact", EvidenceIDs: []string{r.Events[4].ID}})
	text := workflowTraceText(r)
	if !strings.Contains(text, "6h12m47s") || !strings.Contains(text, "may include waiting") {
		t.Fatal("long reported interval needs a legible duration and waiting caveat")
	}
	if *r.Analysis[0].Value != 22367 {
		t.Fatal("presentation changed exact seconds")
	}
}

// Regression for the human's request for boxes and arrows: the original default
// emitted hundreds of lines of source IDs and payloads without a diagram.
func TestWorkflowTraceDefaultDrawsSixDeliveryBoxesInsteadOfRawEvidence(t *testing.T) {
	r := traceDiagramDelivery(t)
	text := workflowTraceText(r)
	if !strings.Contains(text, "+---") || !strings.Contains(text, "|\n") || strings.Count(text, "v  recorded order") != 5 {
		t.Fatalf("default needs six connected step boxes, not a raw evidence dump:\n%s", text)
	}
	previous := -1
	for _, label := range []string{"E1 Work reported", "E2 Verification: passed", "E3 Question for human", "E4 Question for human", "E5 Human QA: pass", "E6 Integration: exact effect"} {
		at := strings.Index(text, label)
		if at <= previous {
			t.Fatalf("missing or out-of-order step %q:\n%s", label, text)
		}
		previous = at
	}
	for _, required := range []string{"Agent", "Ply", "Human", "Same question text as E3", "C1 / R1", "Time: unknown", "--details", "6 unique events", "not causality", "not active agent or human effort"} {
		if !strings.Contains(text, required) {
			t.Errorf("diagram omitted %q", required)
		}
	}
	if strings.Count(text, "Can this installed candidate pass the product check?") != 2 {
		t.Fatal("both original question records need separate visible boxes")
	}
	for _, forbidden := range []string{r.Events[0].ID, r.Candidates[0].ID, r.Candidates[0].Tree, "raw_only_marker", "request_sha256", "projection_records"} {
		if strings.Contains(text, forbidden) {
			t.Errorf("default is drowned in machine evidence: %s", forbidden)
		}
	}
	if len(text) > 9000 {
		t.Fatalf("six simple delivery steps need a compact diagram, got %d bytes", len(text))
	}
}

func TestWorkflowTraceDiagramKeepsTiesChainsAndUnpositionedEventsUnlinked(t *testing.T) {
	r := traceDiagramDelivery(t)
	r.Events = r.Events[:3]
	r.Events[0].Positions = []workflowtrace.Position{{ChainID: "one:delivery", Sequence: 1}}
	r.Events[1].Positions = []workflowtrace.Position{{ChainID: "one:delivery", Sequence: 1}, {ChainID: "two:delivery", Sequence: 1}}
	r.Events[2].Positions = nil
	r.Chains = []workflowtrace.Chain{
		{ID: "one:delivery", EventIDs: []string{r.Events[0].ID, r.Events[1].ID}},
		{ID: "two:delivery", EventIDs: []string{r.Events[1].ID}},
	}
	text := workflowTraceText(r)
	for _, required := range []string{"Same source position: no order", "Same exact event at another source position", "3 unique events", "Other registered records", "E3 Question for human"} {
		if !strings.Contains(text, required) {
			t.Errorf("missing %q", required)
		}
	}
	if strings.Contains(text, "v  recorded order") || strings.Count(text, "E2 Verification") != 2 {
		t.Fatal("invented order or duplicated exact-event identity")
	}
}

func TestWorkflowTraceDiagramExcerptsTextNotEventsOrIdentity(t *testing.T) {
	r := traceDiagramDelivery(t)
	for i, suffix := range []string{"alpha", "beta"} {
		r.Events[i+2].Data, _ = json.Marshal(map[string]string{"question": strings.Repeat("界", 195) + suffix})
	}
	text := workflowTraceText(r)
	if strings.Count(text, "[excerpt:") != 2 || !strings.Contains(text, "2 text field excerpts") || strings.Contains(text, "Same question text") {
		t.Fatal("equal excerpts must not merge distinct full questions")
	}
	details := workflowTraceDetails(r)
	for _, e := range r.Events {
		if !strings.Contains(details, workflowTraceJSON(e.ID)) || !strings.Contains(details, workflowTraceJSON(e.Data)) {
			t.Fatal("details must preserve full evidence for every event")
		}
	}
	if !strings.Contains(details, "E1 = "+workflowTraceJSON(r.Events[0].ID)) {
		t.Fatal("diagram labels must resolve to native IDs")
	}
}

func TestWorkflowTraceDiagramPreservesCandidateJudgmentsAndExplicitCorrectionLinks(t *testing.T) {
	r := workflowTraceFixture(t, "multi-family")
	text := workflowTraceText(r)
	for _, required := range []string{"Human QA: fail", "Human QA: pass", "--corrects-->", "--handoff-->", "reason: Waiting for fixture answer", "dependency: design-question", "next actor: human", "not generation or time order", "Declarations / not executed steps"} {
		if !strings.Contains(text, required) {
			t.Errorf("lost judgment or recorded relationship %q", required)
		}
	}
}

func TestWorkflowTraceDetailsFlagProvidesFullEvidenceAndLabelKey(t *testing.T) {
	r := traceDiagramDelivery(t)
	root, out, stderr := workflowTraceTestRoot(t, newWorkflowTraceCommandWithReader(func(string) (workflowtrace.Result, error) { return r, nil }))
	root.SetArgs([]string{"workflow", "trace", "trace-example", "--details"})
	if err := root.Execute(); err != nil || stderr.Len() != 0 {
		t.Fatalf("details err=%v stderr=%s", err, stderr)
	}
	if out.String() != workflowTraceDetails(r) || !strings.Contains(out.String(), "raw_only_marker") || !strings.Contains(out.String(), "Full evidence (unabridged)") {
		t.Fatal("details flag lost the full projection")
	}
}
