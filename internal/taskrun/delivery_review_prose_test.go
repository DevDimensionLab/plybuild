package taskrun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func al13Prose() string {
	return strings.Repeat("Observed candidate output preserves æøå, 日本語 and useful context. ", 40) + "\n\nThe provider journey is modeled; no live provider approval is claimed.\n\tThe native CLI and stored return are real."
}

// Consume the supplied candidate CLI. Herdr/provider are the existing modeled
// fixture, not a live reviewer. Optional raw process logs survive fixture cleanup.
func al13CLI(t *testing.T, binary string, f deliveryFixture, pane, label string, args ...string) (WorkflowRun, string, int) {
	t.Helper()
	c := exec.Command(binary, append([]string{"workflow", "execute"}, append(args, "--format", "json")...)...)
	c.Dir = f.Parent
	c.Env = append(os.Environ(), "HERDR_PANE_ID="+pane)
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run()
	exit := 0
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			exit = e.ExitCode()
		} else {
			t.Fatal(err)
		}
	}
	if dir := os.Getenv("PLY_AL13_EVIDENCE_DIR"); dir != "" {
		prefix := filepath.Join(dir, strings.ReplaceAll(t.Name(), "/", "_")+"-"+label)
		for suffix, raw := range map[string][]byte{".stdout.json": stdout.Bytes(), ".stderr": stderr.Bytes(), ".exit": []byte(fmt.Sprintln(exit))} {
			if e := os.WriteFile(prefix+suffix, raw, 0600); e != nil {
				t.Fatal(e)
			}
		}
	}
	t.Logf("%s: CLI exit %d", label, exit)
	var run WorkflowRun
	if exit == 0 {
		if e := json.Unmarshal(stdout.Bytes(), &run); e != nil {
			t.Fatalf("JSON stdout: %v\n%s", e, stdout.String())
		}
	}
	return run, stderr.String(), exit
}

func TestAL13ReviewReturnProseThroughCandidateCLI(t *testing.T) {
	binary := deliveryCLIBinary(t)
	for _, decision := range []string{"blocked", "changes_requested", "passed"} {
		t.Run(decision, func(t *testing.T) {
			f, o, _ := reviewRoundFixture(t)
			oid, tree := reviewRoundCandidate(t, f, "fixed prose candidate")
			run, detail, exit := al13CLI(t, binary, f, "w-fixture:root", "request", "review", o.RunID, "--coordinator", "fixture coordinator")
			if exit != 0 {
				t.Fatal(detail)
			}
			round := reviewRoundLatest(t, run)
			prose := al13Prose()
			findings := []map[string]string{}
			if decision == "changes_requested" {
				findings = append(findings, map[string]string{"severity": "high", "summary": "Concrete product defect", "evidence": prose})
			}
			raw, err := json.Marshal(map[string]any{"decision": decision, "summary": "Complete review", "findings": findings, "evidence_gaps": []string{prose}})
			if err != nil {
				t.Fatal(err)
			}
			file := reviewRoundReturnFile(t, f, "prose.json", string(raw))
			returned, detail, exit := al13CLI(t, binary, f, round.ReviewerPane, "return", "review-return", o.RunID, "--context", round.Context.Locator, "--file", file)
			if exit != 0 {
				t.Fatalf("ordinary prose refused: exit %d %s", exit, detail)
			}
			got := reviewRoundLatest(t, returned)
			if got.Decision != decision || got.CandidateOID != oid || got.CandidateTree != tree || got.ReviewerSession != "reviewer-session-1" || len(got.EvidenceGaps) != 1 || got.EvidenceGaps[0] != prose {
				t.Fatalf("return changed judgment, identity or prose: %+v", got)
			}
			if decision == "changes_requested" && (len(got.Findings) != 1 || got.Findings[0].Evidence != prose) {
				t.Fatalf("finding evidence changed: %+v", got.Findings)
			}
			if returned.Delivery.Attempt != nil || returned.DeliveryStatus.FinalDelivery.State != "not_delivered" {
				t.Fatal("review return performed acceptance/delivery")
			}
			// Model loss of the successful response by calling again with identical bytes.
			before := map[string][]byte{}
			dir := filepath.Dir(round.Context.Locator)
			for _, path := range []string{round.Request.Locator, round.Context.Locator, filepath.Join(dir, "return-input.json"), filepath.Join(dir, "return.json"), filepath.Join(o.Paths.RunRoot, "state.json")} {
				b, e := os.ReadFile(path)
				if e != nil {
					t.Fatal(e)
				}
				before[path] = b
			}
			replay, detail, exit := al13CLI(t, binary, f, round.ReviewerPane, "retry", "review-return", o.RunID, "--context", round.Context.Locator, "--file", file)
			if exit != 0 || len(replay.DeliveryStatus.Review.Rounds) != 1 {
				t.Fatalf("identical retry: %d %s", exit, detail)
			}
			changed := reviewRoundReturnFile(t, f, "changed.json", `{"decision":"passed","summary":"Changed judgment","findings":[],"evidence_gaps":[]}`)
			_, detail, exit = al13CLI(t, binary, f, round.ReviewerPane, "sealed", "review-return", o.RunID, "--context", round.Context.Locator, "--file", changed)
			if exit != 4 || !strings.Contains(detail, "sealed return") {
				t.Fatalf("changed sealed judgment: %d %s", exit, detail)
			}
			for path, b := range before {
				after, e := os.ReadFile(path)
				if e != nil || !bytes.Equal(b, after) {
					t.Fatalf("sealed bytes changed: %s %v", path, e)
				}
			}
		})
	}
}

func TestAL13ReviewReturnRejectsInvalidCLIInput(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f, o, _ := reviewRoundFixture(t)
	reviewRoundCandidate(t, f, "fixed input candidate")
	run, detail, exit := al13CLI(t, binary, f, "w-fixture:root", "request", "review", o.RunID, "--coordinator", "fixture coordinator")
	if exit != 0 {
		t.Fatal(detail)
	}
	round := reviewRoundLatest(t, run)
	valid := `{"decision":"passed","summary":"Short existing return","findings":[],"evidence_gaps":[]}`
	for _, tc := range []struct {
		name, body, want string
		exit             int
	}{
		{"trailing-object", valid + ` {"decision":"blocked"}`, "one JSON object", 2},
		{"trailing-garbage", valid + ` trailing`, "one JSON object", 2},
		{"malformed", `{"decision":`, "one JSON object", 2},
		{"schema", `{"schema_version":2,"decision":"passed","summary":"Review"}`, "schema_version", 2},
		{"unknown-field", `{"decision":"passed","summary":"Review","extra":true}`, "unknown field", 2},
		{"wrong-structure", `{"decision":"passed","summary":"Review","evidence_gaps":"text"}`, "one JSON object", 2},
		{"passed-findings", `{"decision":"passed","summary":"Review","findings":[{"severity":"high","summary":"Defect"}]}`, "passed requires no open findings", 2},
		{"changes-empty", `{"decision":"changes_requested","summary":"Review"}`, "at least one concrete finding", 2},
		{"blocked-empty", `{"decision":"blocked","summary":"Review"}`, "blocked needs", 2},
		{"blocked-blank", `{"decision":"blocked","summary":"Review","evidence_gaps":[" \n\t "]}`, "evidence_gaps[0]", 2},
		{"blocked-control", `{"decision":"blocked","summary":"Review","evidence_gaps":["\u0000"]}`, "evidence_gaps[0]", 2},
		{"wrong-session", `{"reviewer_session_id":"owner-or-stranger","decision":"passed","summary":"Review"}`, "reviewer_session_id differs", 4},
		{"oversized", valid + strings.Repeat(" ", 256<<10), "at most 262144 bytes", 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file := reviewRoundReturnFile(t, f, "invalid.json", tc.body)
			_, detail, exit := al13CLI(t, binary, f, round.ReviewerPane, "refused", "review-return", o.RunID, "--context", round.Context.Locator, "--file", file)
			if exit != tc.exit || !strings.Contains(detail, tc.want) {
				t.Fatalf("want exit %d containing %q; got %d %s", tc.exit, tc.want, exit, detail)
			}
			for _, name := range []string{"return.json", "return-input.json"} {
				if _, err := os.Stat(filepath.Join(filepath.Dir(round.Context.Locator), name)); !os.IsNotExist(err) {
					t.Fatalf("invalid input published %s: %v", name, err)
				}
			}
		})
		// A faulty parser can seal a trailing object; stop rather than masking that red.
		if t.Failed() {
			return
		}
	}
	file := reviewRoundReturnFile(t, f, "short.json", valid+"\n \t")
	for _, tc := range []struct{ name, pane, context, want string }{
		{"wrong-pane", "w-fixture:root", round.Context.Locator, "bound reviewer pane"},
		{"wrong-context", round.ReviewerPane, o.Paths.Context, "does not name a review round"},
	} {
		_, detail, exit = al13CLI(t, binary, f, tc.pane, tc.name, "review-return", o.RunID, "--context", tc.context, "--file", file)
		if exit != 4 || !strings.Contains(detail, tc.want) {
			t.Fatalf("binding accepted: %d %s", exit, detail)
		}
	}
	short, detail, exit := al13CLI(t, binary, f, round.ReviewerPane, "short-return", "review-return", o.RunID, "--context", round.Context.Locator, "--file", file)
	if exit != 0 || reviewRoundLatest(t, short).Decision != "passed" {
		t.Fatalf("short return: %d %s", exit, detail)
	}
}

func TestAL13ReviewReturnProseParser(t *testing.T) {
	for _, prose := range []string{strings.Repeat("Ordinary explanation. ", 100), al13Prose(), strings.Repeat("a", 250000)} {
		raw, err := json.Marshal(map[string]any{"decision": "blocked", "summary": "Evidence unavailable", "evidence_gaps": []string{prose}})
		if err != nil {
			t.Fatal(err)
		}
		got, err := deliveryParseReviewReturn(raw, 1)
		if err != nil || len(got.EvidenceGaps) != 1 || got.EvidenceGaps[0] != strings.TrimSpace(prose) {
			t.Fatalf("bounded prose rejected or changed: %v", err)
		}
	}
	for _, control := range []string{"\x00", "\x1b[31m", "\x7f", "\u0085"} {
		raw, err := json.Marshal(map[string]any{"decision": "blocked", "summary": "Evidence unavailable", "evidence_gaps": []string{"Explanation " + control + " text"}})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := deliveryParseReviewReturn(raw, 1); err == nil {
			t.Fatalf("control text accepted: %q", control)
		}
	}
}
