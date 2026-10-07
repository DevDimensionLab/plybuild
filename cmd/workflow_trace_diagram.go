package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/devdimensionlab/plybuild/internal/workflowtrace"
	"github.com/mattn/go-runewidth"
)

const traceExcerptLimit = 180

// Diagram labels are local lookup keys, never generation or time identities.
// The projection is not mutated: JSON and the evidence view retain every fact.
type traceDiagram struct {
	result                               workflowtrace.Result
	events                               map[string]workflowtrace.Event
	eventLabels, candidateLabels         map[string]string
	resultLabels, runLabels, chainLabels map[string]string
	chains                               []workflowtrace.Chain
	eventOrder                           []string
	excerpts                             int
}

func newTraceDiagram(r workflowtrace.Result) *traceDiagram {
	d := &traceDiagram{result: r, events: map[string]workflowtrace.Event{}, eventLabels: map[string]string{}, candidateLabels: map[string]string{}, resultLabels: map[string]string{}, runLabels: map[string]string{}, chainLabels: map[string]string{}}
	for _, e := range r.Events {
		d.events[e.ID] = e
	}
	d.chains = append([]workflowtrace.Chain(nil), r.Chains...)
	// Put recorded delivery steps first for readability, not global chronology.
	sort.SliceStable(d.chains, func(i, j int) bool {
		a, b := strings.HasSuffix(d.chains[i].ID, ":delivery"), strings.HasSuffix(d.chains[j].ID, ":delivery")
		if a != b {
			return a
		}
		return d.chains[i].ID < d.chains[j].ID
	})
	label := func(id string) {
		if _, ok := d.events[id]; ok && d.eventLabels[id] == "" {
			d.eventOrder = append(d.eventOrder, id)
			d.eventLabels[id] = fmt.Sprintf("E%d", len(d.eventOrder))
		}
	}
	for i, c := range d.chains {
		d.chainLabels[c.ID] = fmt.Sprintf("S%d", i+1)
		for _, id := range c.EventIDs {
			label(id)
		}
	}
	for _, e := range r.Events {
		label(e.ID)
	}
	results := map[string]bool{}
	for i, c := range r.Candidates {
		d.candidateLabels[c.ID] = fmt.Sprintf("C%d", i+1)
		for _, id := range c.ResultIDs {
			results[id] = true
		}
	}
	for _, e := range r.Events {
		if e.ResultID != nil {
			results[*e.ResultID] = true
		}
	}
	ids := []string{}
	for id := range results {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for i, id := range ids {
		d.resultLabels[id] = fmt.Sprintf("R%d", i+1)
	}
	for i, run := range r.Runs {
		d.runLabels[run.ID] = fmt.Sprintf("Run%d", i+1)
	}
	return d
}

func workflowTraceText(r workflowtrace.Result) string { return newTraceDiagram(r).render() }

func workflowTraceDetails(r workflowtrace.Result) string {
	d := newTraceDiagram(r)
	var out strings.Builder
	out.WriteString(d.render())
	fmt.Fprintln(&out, "\nDiagram label key (local labels; exact native identities follow)")
	for _, id := range d.eventOrder {
		fmt.Fprintf(&out, "  %s = %s\n", d.eventLabels[id], workflowTraceJSON(id))
	}
	for _, c := range d.chains {
		fmt.Fprintf(&out, "  %s = %s\n", d.chainLabels[c.ID], workflowTraceJSON(c.ID))
	}
	for _, c := range r.Candidates {
		fmt.Fprintf(&out, "  %s = %s\n", d.candidateLabels[c.ID], workflowTraceJSON(c.ID))
	}
	for _, pair := range []map[string]string{d.resultLabels, d.runLabels} {
		ids := []string{}
		for id := range pair {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		for _, id := range ids {
			fmt.Fprintf(&out, "  %s = %s\n", pair[id], workflowTraceJSON(id))
		}
	}
	fmt.Fprintln(&out, "\nFull evidence (unabridged)")
	out.WriteString(workflowTraceEvidenceText(r))
	return out.String()
}

func (d *traceDiagram) render() string {
	var out strings.Builder
	r := d.result
	fmt.Fprintf(&out, "Workflow trace: %s\n", d.excerpt(r.Task.Title))
	fmt.Fprintf(&out, "Task: %s | recorded-source coverage: %s\n", tracePlain(r.Task.ID), tracePlain(r.Coverage.State))
	fmt.Fprintf(&out, "Source freshness: %s | live activity: unknown\n", tracePlain(r.Freshness))
	fmt.Fprintln(&out, "Arrows: recorded order within one source, not causality or elapsed time.")
	fmt.Fprintln(&out, "Sections are independent; their placement is not a global timeline.")
	fmt.Fprintln(&out, "Time: unknown unless recorded. QA/wait intervals may include waiting;")
	fmt.Fprintln(&out, "verifier elapsed is not active agent or human effort. No savings measured.")
	fmt.Fprintf(&out, "Records: %s runs | %s verifier executions | %s candidates | %s QA records\n", d.count("recorded_runs"), d.count("verifier_attempts"), d.count("distinct_observed_candidates"), d.count("human_qa_records"))
	fmt.Fprintf(&out, "%d unique events, %d source chains; all events shown. Full evidence: --details\n", len(r.Events), len(r.Chains))
	if r.HistoryState == "empty" {
		fmt.Fprintln(&out, "No execution recorded. Registered facts and declarations may still exist.")
	}
	shown := map[string]bool{}
	questions := map[string]string{}
	for _, c := range d.chains {
		fmt.Fprintf(&out, "\n%s / %s\n", d.chainLabels[c.ID], d.chainTitle(c))
		previous, hadPrevious := 0, false
		occurrences := map[string]int{}
		for _, id := range c.EventIDs {
			e, found := d.events[id]
			if !found {
				fmt.Fprintln(&out, "  Unresolved event reference; inspect --details.")
				continue
			}
			sequence, positioned := tracePosition(e, c.ID, occurrences[id])
			occurrences[id]++
			if hadPrevious && positioned && sequence > previous {
				fmt.Fprintln(&out, "                                 |")
				fmt.Fprintln(&out, "                                 v  recorded order")
			} else if hadPrevious && positioned && sequence == previous {
				fmt.Fprintln(&out, "         Same source position: no order between these facts.")
			}
			if !positioned {
				fmt.Fprintln(&out, "         Source position unknown; no ordering arrow.")
			}
			lines := d.eventLines(e, questions)
			if shown[id] {
				lines = append(lines, "Same exact event at another source position; not another action.")
			}
			traceBox(&out, workflowTraceRole(e.Role), lines)
			shown[id] = true
			previous, hadPrevious = sequence, positioned
		}
	}
	unpositioned := []workflowtrace.Event{}
	for _, id := range d.eventOrder {
		if !shown[id] {
			unpositioned = append(unpositioned, d.events[id])
		}
	}
	if len(unpositioned) > 0 {
		fmt.Fprintln(&out, "\nOther registered records / no recorded cross-source order")
		for _, e := range unpositioned {
			traceBox(&out, workflowTraceRole(e.Role), d.eventLines(e, questions))
		}
	}
	links := 0
	for _, link := range r.Relations {
		if link.Type == "sequence" {
			continue
		}
		if links == 0 {
			fmt.Fprintln(&out, "\nExplicit recorded links (source claims; not inferred from proximity)")
		}
		fmt.Fprintf(&out, "  [%s] --%s--> [%s]\n", d.eventLabels[link.From], tracePlain(link.Type), d.eventLabels[link.To])
		links++
	}
	if len(r.Candidates) > 0 {
		fmt.Fprintln(&out, "\nCandidate and result key (local labels, not generation or time order)")
		for _, c := range r.Candidates {
			results := []string{}
			for _, id := range c.ResultIDs {
				results = append(results, d.resultLabels[id])
			}
			fmt.Fprintf(&out, "  %s: commit %.8s | %d native results: %s\n", d.candidateLabels[c.ID], c.OID, len(results), strings.Join(results, ", "))
		}
		fmt.Fprintln(&out, "C = exact commit/tree candidate; R = native result, not a run or QA verdict.")
	}
	fmt.Fprintln(&out, "\nDeclarations / not executed steps")
	for _, g := range r.CurrentGoals {
		fmt.Fprintf(&out, "  Current goal: %s revision %d (%s)\n", d.excerpt(g.SpecID), g.Revision, tracePlain(g.Status))
	}
	if len(r.CurrentGoals) == 0 {
		fmt.Fprintln(&out, "  Current goal: unknown")
	}
	for _, run := range r.Runs {
		fmt.Fprintf(&out, "  %s frozen declaration: %s\n", d.runLabels[run.ID], d.frozenDeclaration(run))
	}
	if len(r.Runs) == 0 {
		fmt.Fprintln(&out, "  Declared process unknown: no frozen run contract.")
	}
	patterns := []workflowtrace.Analysis{}
	for _, a := range r.Analysis {
		if a.Question != nil && !strings.HasPrefix(a.ID, "question/") {
			patterns = append(patterns, a)
		}
	}
	if len(patterns) > 0 {
		fmt.Fprintln(&out, "\nQuestions to investigate / not measured causes or improvements")
		for _, a := range patterns {
			fmt.Fprintf(&out, "  [%s] %s\n", d.labels(a.EvidenceIDs), d.excerpt(*a.Question))
		}
	}
	fmt.Fprintf(&out, "\nCoverage: %s selected records; not complete work or conversation history.\n", tracePlain(r.Coverage.State))
	fmt.Fprintln(&out, "Counts: exact = defined records; >= = lower bound; ? = unknown.")
	for _, unknown := range r.Coverage.Unknowns {
		fmt.Fprintf(&out, "  - %s\n", d.excerpt(unknown))
	}
	if len(r.Diagnostics) > 0 {
		fmt.Fprintf(&out, "Diagnostics: %d (all retained in --details / JSON)\n", len(r.Diagnostics))
		for _, diagnostic := range r.Diagnostics {
			fmt.Fprintf(&out, "  - %s: %s\n", d.excerpt(diagnostic.Code), d.excerpt(diagnostic.Detail))
		}
	}
	fmt.Fprintf(&out, "Summary policy: no events hidden; %d text field excerpts (180 chars each).\n", d.excerpts)
	fmt.Fprintln(&out, "IDs, hashes, declarations and raw payloads: --details; unchanged full JSON: --json.")
	return out.String()
}

func (d *traceDiagram) eventLines(e workflowtrace.Event, questions map[string]string) []string {
	var data map[string]json.RawMessage
	_ = json.Unmarshal(e.Data, &data)
	field := func(name string) string { var v string; _ = json.Unmarshal(data[name], &v); return v }
	q := field("question")
	if q == "" && (e.Type == "question" || e.Type == "decision_requested") {
		q = e.Title
	}
	title := traceEventTitle(e, q != "")
	lines := []string{d.eventLabels[e.ID] + " " + d.excerpt(title)}
	class := e.EvidenceClass
	if class == "" {
		class = "unknown"
	}
	binding := ""
	if e.CandidateID != nil {
		binding = d.candidateLabels[*e.CandidateID]
	}
	if e.ResultID != nil {
		if binding != "" {
			binding += " / "
		}
		binding += d.resultLabels[*e.ResultID]
	}
	if binding != "" {
		class += " | " + binding
	}
	lines = append(lines, tracePlain(strings.ReplaceAll(class, "_", " ")))
	if q != "" {
		if prior := questions[q]; prior != "" && prior != d.eventLabels[e.ID] {
			lines = append(lines, "Same question text as "+prior+"; separate source record.")
		} else {
			questions[q] = d.eventLabels[e.ID]
		}
		lines = append(lines, "Q: "+d.excerpt(q))
	} else if summary := field("summary"); summary != "" && e.Type == "delivery_report" {
		lines = append(lines, d.excerpt(summary))
	}
	var waiting map[string]string
	_ = json.Unmarshal(data["waiting"], &waiting)
	for _, key := range []string{"reason", "dependency", "next_actor"} {
		if v := waiting[key]; v != "" {
			lines = append(lines, strings.ReplaceAll(key, "_", " ")+": "+d.excerpt(v))
		}
	}
	clock := "Time: unknown"
	if e.OccurredAtUTC != nil {
		clock = "Occurred: " + tracePlain(*e.OccurredAtUTC)
	} else if e.ReportedAtUTC != nil {
		clock = "Reported: " + tracePlain(*e.ReportedAtUTC)
	} else if e.RegisteredAtUTC != nil {
		clock = "Registered: " + tracePlain(*e.RegisteredAtUTC)
	}
	lines = append(lines, clock)
	for _, a := range d.result.Analysis {
		if a.Kind != "interval" || !traceIncludes(a.EvidenceIDs, e.ID) {
			continue
		}
		label := "Reported interval"
		if strings.HasPrefix(a.ID, "verifier_elapsed/") {
			label = "Verifier elapsed"
		} else if strings.HasPrefix(a.ID, "reported_qa/") {
			label = "Reported QA interval (may include waiting)"
		}
		value := "unknown"
		if a.Value != nil && a.Coverage != "unknown" {
			value = traceDuration(*a.Value)
		}
		lines = append(lines, label+": "+value)
	}
	return lines
}

func traceEventTitle(e workflowtrace.Event, question bool) string {
	outcome := "unknown"
	if e.Outcome != nil {
		outcome = strings.ReplaceAll(*e.Outcome, "_", " ")
	}
	switch e.Type {
	case "delivery_report":
		if question {
			return "Question for human"
		}
		return "Work reported: " + outcome
	case "verification":
		return "Verification: " + outcome
	case "human_qa":
		return "Human QA: " + outcome
	case "integration":
		return "Integration: " + outcome
	case "task_result_technical":
		return "Result technical gate: " + outcome
	case "task_result_reported":
		return "Result reported: " + outcome
	case "run_request":
		return "Frozen request registered (declaration)"
	case "run_acceptance":
		return "Runtime acceptance: " + outcome
	}
	return strings.ReplaceAll(e.Title, "_", " ")
}

func (d *traceDiagram) chainTitle(c workflowtrace.Chain) string {
	name := "Recorded source chain"
	for suffix, title := range map[string]string{":delivery": "Delivery", ":acceptance": "Runtime acceptance", ":candidates": "Candidate facts", ":request": "Declaration records"} {
		if strings.HasSuffix(c.ID, suffix) {
			name = title
		}
	}
	for _, run := range d.result.Runs {
		if strings.HasPrefix(c.ID, run.ID+":") {
			return d.runLabels[run.ID] + " / " + name
		}
	}
	return name
}

func (d *traceDiagram) frozenDeclaration(run workflowtrace.Run) string {
	if !workflowTraceHasDeclaration(run.DeclaredProcess) {
		return "Declared process unknown"
	}
	var goal struct {
		SpecID string `json:"spec_id"`
		Spec   struct {
			Revision int `json:"revision"`
		} `json:"spec"`
	}
	_ = json.Unmarshal(run.FrozenGoal, &goal)
	if goal.SpecID != "" {
		return fmt.Sprintf("%s revision %d; process details: --details", d.excerpt(goal.SpecID), goal.Spec.Revision)
	}
	return "preserved; process details: --details"
}

func (d *traceDiagram) labels(ids []string) string {
	labels := []string{}
	for _, id := range ids {
		labels = append(labels, d.eventLabels[id])
	}
	return strings.Join(labels, ", ")
}

func (d *traceDiagram) count(id string) string {
	for _, a := range d.result.Analysis {
		if a.ID == id && a.Value != nil && a.Coverage != "unknown" {
			prefix := ""
			if a.Coverage == "lower_bound" {
				prefix = ">="
			}
			return fmt.Sprintf("%s%g", prefix, *a.Value)
		}
	}
	return "?"
}

func (d *traceDiagram) excerpt(text string) string {
	runes := []rune(text)
	if len(runes) <= traceExcerptLimit {
		return tracePlain(text)
	}
	d.excerpts++
	return tracePlain(string(runes[:traceExcerptLimit])) + fmt.Sprintf(" [excerpt: %d chars omitted; --details]", len(runes)-traceExcerptLimit)
}

func tracePlain(text string) string {
	var out strings.Builder
	for _, r := range text {
		if r == '\\' {
			out.WriteString(`\\`)
		} else if unicode.IsPrint(r) {
			out.WriteRune(r)
		} else {
			quoted := strconv.QuoteRuneToASCII(r)
			out.WriteString(quoted[1 : len(quoted)-1])
		}
	}
	return out.String()
}

func traceBox(out *strings.Builder, actor string, lines []string) {
	border := "         +" + strings.Repeat("-", 64) + "+\n"
	out.WriteString(border)
	for _, line := range lines {
		for _, wrapped := range traceWrap(line, 62) {
			fmt.Fprintf(out, "%-8s | %s%s |\n", actor, wrapped, strings.Repeat(" ", 62-runewidth.StringWidth(wrapped)))
			actor = ""
		}
	}
	out.WriteString(border)
}

func traceWrap(s string, width int) []string {
	var lines []string
	for runewidth.StringWidth(s) > width {
		end := len(runewidth.Truncate(s, width, ""))
		cut := strings.LastIndexByte(s[:end], ' ')
		if cut <= 0 {
			cut = end
		}
		lines = append(lines, s[:cut])
		s = strings.TrimPrefix(s[cut:], " ")
	}
	return append(lines, s)
}

func traceIncludes(ids []string, id string) bool {
	for _, candidate := range ids {
		if candidate == id {
			return true
		}
	}
	return false
}

func traceDuration(seconds float64) string {
	if seconds < 0 || math.IsNaN(seconds) || math.IsInf(seconds, 0) {
		return "unknown"
	}
	seconds = math.Round(seconds*1000) / 1000
	hours := math.Floor(seconds / 3600)
	minutes := math.Floor(math.Mod(seconds, 3600) / 60)
	rest := math.Mod(seconds, 60)
	var out strings.Builder
	if hours > 0 {
		fmt.Fprintf(&out, "%.0fh", hours)
	}
	if minutes > 0 || hours > 0 {
		fmt.Fprintf(&out, "%.0fm", minutes)
	}
	fmt.Fprintf(&out, "%ss", strings.TrimRight(strings.TrimRight(strconv.FormatFloat(rest, 'f', 3, 64), "0"), "."))
	return out.String()
}

func tracePosition(e workflowtrace.Event, chain string, occurrence int) (int, bool) {
	positions := []int{}
	for _, p := range e.Positions {
		if p.ChainID == chain {
			positions = append(positions, p.Sequence)
		}
	}
	sort.Ints(positions)
	if occurrence < len(positions) {
		return positions[occurrence], true
	}
	return 0, false
}
