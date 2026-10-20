package taskrun

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Review rounds are ordinary transitions of the same delivery Run. Each round
// binds one exact clean candidate; its files are write-once beside the Run and
// never enter the frozen state.json that older controls read strictly. Ply
// derives identities, IDs, ordering and the qualification record itself, so a
// coordinator or reviewer supplies only actual judgment.

type DeliveryReviewRequestInput struct {
	Coordinator     string `json:"coordinator"`
	CoordinatorPane string `json:"coordinator_pane"`
	HerdrWorkspace  string `json:"herdr_workspace"`
	NewRound        bool   `json:"new_round"`
}

type deliveryReviewContext struct {
	Kind          string `json:"kind"`
	SchemaVersion int    `json:"schema_version"`
	RunID         string `json:"run_id"`
	RequestSHA256 string `json:"request_sha256"`
	Round         int    `json:"round"`
	CandidateOID  string `json:"candidate_oid"`
	CandidateTree string `json:"candidate_tree"`
}

type deliveryReviewRound struct {
	Kind            string       `json:"kind"`
	SchemaVersion   int          `json:"schema_version"`
	RunID           string       `json:"run_id"`
	RequestSHA256   string       `json:"request_sha256"`
	Round           int          `json:"round"`
	CandidateOID    string       `json:"candidate_oid"`
	CandidateTree   string       `json:"candidate_tree"`
	Coordinator     string       `json:"coordinator"`
	CoordinatorPane string       `json:"coordinator_pane"`
	OwnerSession    string       `json:"owner_session"`
	Previous        *FileBinding `json:"previous"`
	Supersedes      *FileBinding `json:"supersedes,omitempty"` // absent in rounds written by 8066ecb
	Brief           FileBinding  `json:"brief"`
	Context         FileBinding  `json:"context"`
	Template        FileBinding  `json:"template"`
	RequestedAtUTC  string       `json:"requested_at_utc"`
}

type DeliveryReviewFinding struct {
	ID       string `json:"id"`
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
	Evidence string `json:"evidence"`
}

type deliveryReviewResult struct {
	Kind            string                  `json:"kind"`
	SchemaVersion   int                     `json:"schema_version"`
	RunID           string                  `json:"run_id"`
	RequestSHA256   string                  `json:"request_sha256"`
	Round           int                     `json:"round"`
	CandidateOID    string                  `json:"candidate_oid"`
	CandidateTree   string                  `json:"candidate_tree"`
	RoundRequest    FileBinding             `json:"round_request"`
	ReviewState     string                  `json:"review_state"`
	Decision        string                  `json:"decision"`
	Summary         string                  `json:"summary"`
	ReviewerClaim   string                  `json:"reviewer_claim"`
	ReviewerSession string                  `json:"reviewer_session_id"`
	Findings        []DeliveryReviewFinding `json:"findings"`
	EvidenceGaps    []string                `json:"evidence_gaps"`
	Source          FileBinding             `json:"source"`
	RecordedAtUTC   string                  `json:"recorded_at_utc"`
}

type deliveryReviewRecord struct {
	Superseded bool
	Dir        string
	Request    deliveryReviewRound
	Binding    FileBinding
	Result     *deliveryReviewResult
	Return     *FileBinding
}

// DeliveryReviewStatus is a read-only projection of the preserved rounds.
type DeliveryReviewStatus struct {
	CurrentCandidate string                      `json:"current_candidate"`
	Rounds           []DeliveryReviewRoundStatus `json:"rounds"`
	Next             WorkflowAction              `json:"next"`
}

type DeliveryReviewRoundStatus struct {
	Round           int                     `json:"round"`
	CandidateOID    string                  `json:"candidate_oid"`
	CandidateTree   string                  `json:"candidate_tree"`
	Coordinator     string                  `json:"coordinator"`
	State           string                  `json:"state"`
	Decision        string                  `json:"decision,omitempty"`
	Summary         string                  `json:"summary,omitempty"`
	ReviewerClaim   string                  `json:"reviewer_claim,omitempty"`
	ReviewerSession string                  `json:"reviewer_session_id,omitempty"`
	Findings        []DeliveryReviewFinding `json:"findings"`
	EvidenceGaps    []string                `json:"evidence_gaps"`
	Current         bool                    `json:"current"`
	ReviewerStart   string                  `json:"reviewer_start"`
	ReviewerPane    string                  `json:"reviewer_pane,omitempty"`
	StartDetail     string                  `json:"start_detail,omitempty"`
	Previous        *FileBinding            `json:"previous"`
	Supersedes      *FileBinding            `json:"supersedes"`
	Request         FileBinding             `json:"request"`
	Brief           FileBinding             `json:"brief"`
	Context         FileBinding             `json:"context"`
	Template        FileBinding             `json:"template"`
	Result          *FileBinding            `json:"result"`
}

var deliveryReviewRoundName = regexp.MustCompile(`^[0-9]{3}$`)

func deliveryReviewDir(s workflowState) string {
	return filepath.Join(s.Result.Paths.RunRoot, "review", "rounds")
}

func deliveryReviewRounds(s workflowState) ([]deliveryReviewRecord, error) {
	root := deliveryReviewDir(s)
	if e := physical(root, true); e != nil {
		return nil, e
	}
	entries, e := os.ReadDir(root)
	if os.IsNotExist(e) {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	out := []deliveryReviewRecord{}
	for i, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".publish-") {
			continue
		}
		if !entry.IsDir() || !deliveryReviewRoundName.MatchString(entry.Name()) || entry.Name() != workflowRoundName(len(out)+1) {
			return nil, workflowError(4, "unrecognized review round entry: "+entry.Name())
		}
		dir := filepath.Join(root, entry.Name())
		raw, e := readFile(filepath.Join(dir, "request.json"), 1<<20, true)
		if os.IsNotExist(e) && i == len(entries)-1 {
			// An interrupted request published only its inputs; the same
			// request completes this number, a different one conflicts.
			break
		}
		if e != nil {
			return nil, e
		}
		r := deliveryReviewRecord{Dir: dir, Binding: FileBinding{filepath.Join(dir, "request.json"), hash(raw)}}
		if e = decode(raw, 1<<20, &r.Request); e != nil {
			return nil, e
		}
		if r.Request.Kind != "ply.workflow.review-round" || r.Request.SchemaVersion != 1 || r.Request.RunID != s.Result.RunID || r.Request.RequestSHA256 != s.Result.RequestSHA256 || r.Request.Round != len(out)+1 {
			return nil, workflowError(4, "review round binding changed: "+entry.Name())
		}
		raw, e = readFile(filepath.Join(dir, "return.json"), 1<<20, true)
		if e == nil {
			var result deliveryReviewResult
			if e = decode(raw, 1<<20, &result); e != nil {
				return nil, e
			}
			if result.Kind != "ply.workflow.review-result" || result.SchemaVersion != 1 || result.Round != r.Request.Round || result.CandidateOID != r.Request.CandidateOID || result.CandidateTree != r.Request.CandidateTree || result.RoundRequest != r.Binding {
				return nil, workflowError(4, "review return binding changed: "+entry.Name())
			}
			r.Result, r.Return = &result, &FileBinding{filepath.Join(dir, "return.json"), hash(raw)}
		} else if !os.IsNotExist(e) {
			return nil, e
		}
		out = append(out, r)
	}
	for _, r := range out {
		for i := range out {
			if r.Request.Supersedes != nil && *r.Request.Supersedes == out[i].Binding {
				out[i].Superseded = true
			}
		}
	}
	return out, nil
}

// WorkflowDeliveryRequestReview is the coordinator's single operation to
// review the owner's current clean candidate: it binds the round and starts
// the separate reviewer session through Ply's Herdr launcher. It needs no
// review Task, queue move, preparation or owner callback, and the owner keeps
// working. A reserved start is never replayed.
func WorkflowDeliveryRequestReview(d Dependencies, root, id string, in DeliveryReviewRequestInput) (WorkflowRun, error) {
	if e := containing(d, root); e != nil {
		return WorkflowRun{}, e
	}
	if !plain(in.Coordinator, 1, 256) {
		return WorkflowRun{}, workflowError(2, "--coordinator must name the actual review coordinator (1-256 characters)")
	}
	executable, e := d.Executable()
	if e != nil {
		return WorkflowRun{}, e
	}
	if resolved, err := filepath.EvalSymlinks(executable); err == nil {
		executable = resolved
	}
	var s workflowState
	var round deliveryReviewRecord
	var launch *deliveryReviewLaunch
	e = withStore(root, func() error {
		var e error
		s, e = workflowRead(root, id)
		if e != nil {
			return e
		}
		if !deliveryRun(s.Request) || s.Result.Delivery == nil {
			return workflowError(4, "review rounds require a delivery-owner Run")
		}
		if s.Result.Delivery.Phase == "completed" {
			return workflowError(4, "this delivery is complete; its review history is preserved")
		}
		if r, e := deliveryRetirement(s); e != nil {
			return e
		} else if r != nil {
			return workflowError(4, "owner_retired: review the fresh Task copy's Run instead")
		}
		coordinatorPane, e := deliveryCallerPane(d, s, in.CoordinatorPane)
		if e != nil {
			return e
		}
		if coordinatorPane == s.Result.Transport.PaneID {
			return workflowError(4, "the owner cannot coordinate its own review; the coordinator runs this from its own session")
		}
		x, e := deliveryTarget(d, s)
		if e != nil {
			return workflowError(4, "review needs a clean committed candidate; the owner commits first: "+e.Error())
		}
		rounds, e := deliveryReviewRounds(s)
		if e != nil {
			return e
		}
		n := len(rounds)
		same := n > 0 && rounds[n-1].Request.CandidateOID == x.OID && rounds[n-1].Request.CandidateTree == x.Tree
		if in.NewRound && (!same || rounds[n-1].Result != nil || deliveryReviewStartState(rounds[n-1].Dir).State == "started") {
			return workflowError(4, "--new-round only replaces an unknown reviewer start for the same candidate before any return")
		}
		if same && !in.NewRound {
			round = rounds[n-1]
			return deliveryReviewReserve(d, s, round, in.HerdrWorkspace, &launch)
		}
		number := len(rounds) + 1
		dir := filepath.Join(deliveryReviewDir(s), workflowRoundName(number))
		r := deliveryReviewRound{Kind: "ply.workflow.review-round", SchemaVersion: 1, RunID: id, RequestSHA256: s.Result.RequestSHA256, Round: number, CandidateOID: x.OID, CandidateTree: x.Tree, Coordinator: in.Coordinator, CoordinatorPane: coordinatorPane, OwnerSession: deliveryOwnerSession(s), RequestedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
		if number > 1 {
			r.Previous = &rounds[number-2].Binding
		}
		if in.NewRound {
			// The coordinator replaces an unknown reviewer start. The old
			// round keeps its history but can no longer return or govern.
			r.Supersedes = &rounds[n-1].Binding
		}
		if r.Context, e = workflowKeep(d, filepath.Join(dir, "context.json"), deliveryReviewContext{"ply.workflow.review-context", 1, id, s.Result.RequestSHA256, number, x.OID, x.Tree}); e != nil {
			return e
		}
		if r.Template, e = workflowKeep(d, filepath.Join(dir, "return-template.json"), deliveryReviewTemplate()); e != nil {
			return e
		}
		brief := deliveryReviewBrief(s, r, rounds, root, executable, filepath.Join(dir, "return-template.json"), filepath.Join(dir, "context.json"))
		r.Brief = FileBinding{filepath.Join(dir, "brief.md"), hash([]byte(brief))}
		if e = d.writeOnce(r.Brief.Locator, []byte(brief)); e != nil {
			return e
		}
		binding, e := workflowKeep(d, filepath.Join(dir, "request.json"), r)
		if e != nil {
			return e
		}
		round = deliveryReviewRecord{Dir: dir, Request: r, Binding: binding}
		return deliveryReviewReserve(d, s, round, in.HerdrWorkspace, &launch)
	})
	if e == nil && launch != nil {
		e = deliveryLaunchReviewer(d, s, round, *launch)
	}
	o, err := WorkflowShow(d, root, id)
	if e != nil {
		return o, e
	}
	return o, err
}

// deliveryCallerPane binds the coordinator's actual Herdr pane. It is the
// pane whose shell holds this process; an inherited HERDR_PANE_ID is only the
// first pane checked, since it can name a pane the caller left. A detached
// caller (for example a provider daemon's tool process) names its verified
// pane explicitly; Ply checks it against fresh Herdr state and the observed
// ancestry. Focus is never used, and an unknown or ambiguous caller stops
// before any review effect.
func deliveryCallerPane(d Dependencies, s workflowState, explicit string) (string, error) {
	next := "; run `ply workflow execute review` from the coordinator's own Herdr pane shell, or pass --coordinator-pane with its verified live pane ID"
	inferred, ancestry := deliveryCallerPaneFromAncestry(d, s)
	if explicit != "" {
		if len(inferred) > 0 && !slices.Contains(inferred, explicit) {
			return "", workflowError(4, fmt.Sprintf("--coordinator-pane %s contradicts pane %s, whose shell holds this process", explicit, strings.Join(inferred, ", ")))
		}
		if invokerShellPID(d, s, explicit) == 0 {
			return "", workflowError(4, fmt.Sprintf("--coordinator-pane %s is not a live Herdr pane; pass the coordinator's verified live pane ID", explicit))
		}
		return explicit, nil
	}
	switch {
	case len(inferred) == 1:
		return inferred[0], nil
	case len(inferred) > 1:
		return "", workflowError(4, fmt.Sprintf("the coordinator's Herdr pane is ambiguous: %s all hold this process%s", strings.Join(inferred, ", "), next))
	}
	return "", workflowError(4, fmt.Sprintf("the coordinator's Herdr pane is unknown: %s (HERDR_PANE_ID %q)%s", ancestry, os.Getenv("HERDR_PANE_ID"), next))
}

// deliveryCallerPaneFromAncestry returns the panes whose shell is an ancestor
// of this process, or why none could be observed.
func deliveryCallerPaneFromAncestry(d Dependencies, s workflowState) ([]string, string) {
	observe := d.LocalProcesses
	if observe == nil {
		observe = systemLocalProcesses
	}
	table, e := observe("", "")
	if e != nil {
		return nil, "local process inspection failed: " + e.Error()
	}
	parent := map[int]int{}
	for _, p := range table {
		parent[p.PID] = p.ParentPID
	}
	ancestors := map[int]bool{}
	for pid, ok := parent[os.Getpid()]; ok && pid > 1 && !ancestors[pid]; pid, ok = parent[pid] {
		ancestors[pid] = true
	}
	if hint := os.Getenv("HERDR_PANE_ID"); hint != "" && ancestors[invokerShellPID(d, s, hint)] {
		return []string{hint}, ""
	}
	raw, e := workflowCall(d, s.Request, "pane", "list")
	var list struct {
		Panes []struct {
			PaneID string `json:"pane_id"`
		} `json:"panes"`
	}
	if e == nil {
		e = json.Unmarshal(raw, &list)
	}
	if e != nil {
		return nil, "Herdr pane list failed: " + e.Error()
	}
	matches := []string{}
	for _, p := range list.Panes {
		if ancestors[invokerShellPID(d, s, p.PaneID)] {
			matches = append(matches, p.PaneID)
		}
	}
	return matches, "no pane shell holds this process"
}

// deliveryReviewReserve records the one reviewer start of a round before any
// Herdr effect. An existing reservation is only observed, never replayed.
func deliveryReviewReserve(d Dependencies, s workflowState, r deliveryReviewRecord, workspaceID string, out **deliveryReviewLaunch) error {
	if deliveryReviewStartState(r.Dir).State != "not_started" || r.Result != nil {
		return nil
	}
	l, e := deliveryReviewPlan(d, s, r.Request.Round, workspaceID)
	if e != nil {
		return e
	}
	if _, e = workflowKeep(d, filepath.Join(r.Dir, "launch.json"), l); e != nil {
		return e
	}
	*out = &l
	return nil
}

func deliveryReviewTemplate() map[string]any {
	return map[string]any{"kind": "ply.workflow.review-return", "schema_version": 1, "reviewer_claim": "", "reviewer_session_id": "", "decision": "passed | changes_requested | blocked", "summary": "", "findings": []any{map[string]any{"severity": "low | medium | high | critical", "summary": "", "evidence": ""}}, "evidence_gaps": []any{}}
}

func deliveryReviewBrief(s workflowState, r deliveryReviewRound, earlier []deliveryReviewRecord, root, executable, template, context string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Independent review round %d\n\n", r.Round)
	fmt.Fprintf(&b, "Coordinator: %s. You are a separate reviewer session, not the developer (owner session %s).\n\n", r.Coordinator, r.OwnerSession)
	fmt.Fprintf(&b, "- Run: %s\n- Goal and Spec: %s\n- Task worktree (read-only for you): %s\n- Task base: %s\n- Candidate: %s (tree %s)\n- Owner acceptance entrypoint: %s\n", r.RunID, s.Request.Delivery.Goal.Locator, s.Observed.Target.WorktreeLocator, s.Observed.Target.OID, r.CandidateOID, r.CandidateTree, s.Request.Delivery.AcceptancePath)
	for _, e := range earlier {
		if e.Return != nil {
			fmt.Fprintf(&b, "- Earlier round %d (%s, candidate %s): %s\n", e.Request.Round, e.Result.Decision, e.Request.CandidateOID, e.Return.Locator)
		}
	}
	b.WriteString("\nReview exactly this candidate against the Goal and Spec: the diff from the base, tests and actual evidence. Do not edit the Task worktree, commit, verify or integrate.\n")
	fmt.Fprintf(&b, "Keep the candidate exact while testing. This pane runs with GOFLAGS=%s: Go modules stay readonly, so never pass -mod=mod, run go mod tidy, go get or change global Go or provider configuration. After the tests, `git -C %s status --porcelain` must print nothing and `git -C %s rev-parse HEAD` must print %s. If not, Ply cannot confirm an exact-candidate pass from the current source: Never reset, clean or stash it; return blocked or changes_requested with the observed status and diff. Ply refuses passed while the Task worktree is dirty or at another HEAD.\n", deliveryReviewGoFlags(), ShellQuote(s.Observed.Target.WorktreeLocator), ShellQuote(s.Observed.Target.WorktreeLocator), r.CandidateOID)
	b.WriteString("A complete review with product findings is a completed review: use changes_requested with concrete findings. Use blocked when evidence is insufficient to judge. passed requires no findings. Keep truthful evidence gaps; they never block the return.\n\n")
	fmt.Fprintf(&b, "Copy %s to your own file and fill in the decision, summary, findings and evidence gaps (reviewer_claim and reviewer_session_id are optional: Ply binds the session it started for you). Return from this reviewer pane, inside the workspace %s:\n\n", template, root)
	fmt.Fprintf(&b, "    %s workflow execute review-return %s --context %s --file <your-review.json>\n\n", ShellQuote(executable), r.RunID, ShellQuote(context))
	b.WriteString("Ply assigns finding IDs and ordering. Repeating the same file is safe; a changed judgment for this candidate needs a new round from the coordinator.\n")
	return b.String()
}

type deliveryReviewReturnInput struct {
	Kind              *string `json:"kind"`
	SchemaVersion     *int    `json:"schema_version"`
	ReviewerClaim     string  `json:"reviewer_claim"`
	ReviewerSessionID string  `json:"reviewer_session_id"`
	Decision          string  `json:"decision"`
	Summary           string  `json:"summary"`
	Findings          []struct {
		Severity string `json:"severity"`
		Summary  string `json:"summary"`
		Evidence string `json:"evidence"`
	} `json:"findings"`
	EvidenceGaps []string `json:"evidence_gaps"`
}

func deliveryParseReviewReturn(raw []byte, round int) (deliveryReviewResult, error) {
	var in deliveryReviewReturnInput
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if e := dec.Decode(&in); e != nil {
		return deliveryReviewResult{}, workflowError(2, "review return must be one JSON object with reviewer_claim, reviewer_session_id, decision, summary, findings and evidence_gaps: "+e.Error())
	}
	if in.Kind != nil && *in.Kind != "ply.workflow.review-return" || in.SchemaVersion != nil && *in.SchemaVersion != 1 {
		return deliveryReviewResult{}, workflowError(2, "review return kind must be ply.workflow.review-return schema_version 1 when given")
	}
	norm := func(v string) string {
		return strings.NewReplacer("-", "_", " ", "_").Replace(strings.ToLower(strings.TrimSpace(v)))
	}
	out := deliveryReviewResult{ReviewState: "completed", Decision: norm(in.Decision), Summary: strings.TrimSpace(in.Summary), ReviewerClaim: strings.TrimSpace(in.ReviewerClaim), ReviewerSession: strings.TrimSpace(in.ReviewerSessionID), Findings: []DeliveryReviewFinding{}, EvidenceGaps: []string{}}
	if out.ReviewerClaim != "" && !plain(out.ReviewerClaim, 1, 256) || out.ReviewerSession != "" && !plain(out.ReviewerSession, 1, 256) {
		return out, workflowError(2, "reviewer_claim and reviewer_session_id are optional, at most 256 characters on one line")
	}
	if out.Decision != "passed" && out.Decision != "changes_requested" && out.Decision != "blocked" {
		return out, workflowError(2, "decision must be passed, changes_requested or blocked")
	}
	if !plain(out.Summary, 1, 2000) {
		return out, workflowError(2, "summary must be 1-2000 characters on one line")
	}
	if len(in.Findings) > 100 || len(in.EvidenceGaps) > 100 {
		return out, workflowError(2, "at most 100 findings and 100 evidence gaps")
	}
	for i, f := range in.Findings {
		x := DeliveryReviewFinding{ID: fmt.Sprintf("r%03d-f%02d", round, i+1), Severity: norm(f.Severity), Summary: strings.TrimSpace(f.Summary), Evidence: strings.TrimSpace(f.Evidence)}
		if x.Severity != "low" && x.Severity != "medium" && x.Severity != "high" && x.Severity != "critical" {
			return out, workflowError(2, fmt.Sprintf("findings[%d].severity must be low, medium, high or critical", i))
		}
		if !plain(x.Summary, 1, 600) || x.Evidence != "" && !plain(x.Evidence, 1, 2000) {
			return out, workflowError(2, fmt.Sprintf("findings[%d] needs a one-line summary (1-600 characters) and optional one-line evidence (at most 2000)", i))
		}
		out.Findings = append(out.Findings, x)
	}
	for i, g := range in.EvidenceGaps {
		g = strings.TrimSpace(g)
		if !plain(g, 1, 600) {
			return out, workflowError(2, fmt.Sprintf("evidence_gaps[%d] must be 1-600 characters on one line", i))
		}
		out.EvidenceGaps = append(out.EvidenceGaps, g)
	}
	switch {
	case out.Decision == "passed" && len(out.Findings) != 0:
		return out, workflowError(2, "passed requires no open findings; use changes_requested")
	case out.Decision == "changes_requested" && len(out.Findings) == 0:
		return out, workflowError(2, "changes_requested needs at least one concrete finding")
	case out.Decision == "blocked" && len(out.Findings) == 0 && len(out.EvidenceGaps) == 0:
		return out, workflowError(2, "blocked needs a finding or an evidence gap that explains it")
	}
	return out, nil
}

// WorkflowDeliveryReviewReturn seals one reviewer's completed judgment for the
// round's own candidate. A later owner commit does not invalidate it; status
// shows the drift and the new candidate needs its own round.
func WorkflowDeliveryReviewReturn(d Dependencies, root, id, contextPath, file string) (WorkflowRun, error) {
	if e := containing(d, root); e != nil {
		return WorkflowRun{}, e
	}
	raw, e := readFile(file, 256<<10, false)
	if e != nil {
		return WorkflowRun{}, e
	}
	e = withStore(root, func() error {
		s, e := workflowRead(root, id)
		if e != nil {
			return e
		}
		if !deliveryRun(s.Request) || s.Result.Delivery == nil {
			return workflowError(4, "review rounds require a delivery-owner Run")
		}
		rounds, e := deliveryReviewRounds(s)
		if e != nil {
			return e
		}
		var r *deliveryReviewRecord
		for i := range rounds {
			if rounds[i].Request.Context.Locator == contextPath {
				r = &rounds[i]
			}
		}
		if r == nil {
			return workflowError(4, "--context does not name a review round of this Run; use the context from the round's brief")
		}
		if _, e = workflowBound(r.Request.Context, 1<<20); e != nil {
			return e
		}
		if r.Superseded && r.Result == nil {
			return workflowError(4, fmt.Sprintf("review round %d was superseded by the coordinator's fresh round for the same candidate; its reviewer cannot return", r.Request.Round))
		}
		start := deliveryReviewStartState(r.Dir)
		if start.Agent == nil {
			return workflowError(4, fmt.Sprintf("review round %d has no reviewer session started and bound by Ply (%s); the coordinator starts it with execute review", r.Request.Round, start.State))
		}
		if pane := os.Getenv("HERDR_PANE_ID"); pane != start.Agent.Pane.PaneID {
			return workflowError(4, fmt.Sprintf("review-return must come from the bound reviewer pane %s (caller pane %q)", start.Agent.Pane.PaneID, pane))
		}
		limit := 15 * time.Second
		if d.HerdrTimeout > 0 {
			limit = d.HerdrTimeout
		}
		a, e := deliveryReviewAgentGet(d, s, start.Agent.Pane, start.Agent.AgentName, start.Agent.Provider, time.Now().Add(limit))
		if e != nil {
			return e
		}
		if a.Session == nil || a.Session.Value != start.Agent.SessionID {
			return workflowError(4, "the bound reviewer pane no longer shows the reviewer session started by Ply")
		}
		source := FileBinding{filepath.Join(r.Dir, "return-input.json"), hash(raw)}
		if r.Result != nil {
			if r.Result.Source.SHA256 == source.SHA256 {
				return nil
			}
			return workflowError(4, fmt.Sprintf("review round %d already has a different sealed return; a changed judgment needs a new round from the coordinator", r.Request.Round))
		}
		result, e := deliveryParseReviewReturn(raw, r.Request.Round)
		if e != nil {
			return e
		}
		if result.Decision == "passed" {
			// Ply confirms a pass only while the observed current source is
			// this round's exact candidate; a dirty worktree or another HEAD
			// cannot confirm it. Sealed returns stay as they are and a
			// negative return is still sealed.
			x, e := deliveryTarget(d, s)
			if e != nil {
				paths := []string{}
				for _, entry := range x.StatusEntries {
					path, _ := base64.StdEncoding.DecodeString(entry.PathBase64)
					paths = append(paths, strconv.Quote(string(path)))
				}
				return workflowError(4, fmt.Sprintf("passed is refused: the Task worktree is not the clean exact candidate %s (changed: %s); Ply does not reset it. Return blocked or changes_requested with the observed status and diff", r.Request.CandidateOID, strings.Join(paths, ", ")))
			}
			if x.OID != r.Request.CandidateOID || x.Tree != r.Request.CandidateTree {
				return workflowError(4, fmt.Sprintf("passed is refused: the Task worktree is now at %s, not round %d's candidate %s, so Ply cannot confirm an exact-candidate pass from the current source; Ply does not reset it. Return blocked or changes_requested for this round; the coordinator requests a new round for %s", x.OID, r.Request.Round, r.Request.CandidateOID, x.OID))
			}
		}
		if result.ReviewerSession != "" && result.ReviewerSession != start.Agent.SessionID {
			return workflowError(4, "reviewer_session_id differs from the reviewer session Ply started and observed; omit it or give that session")
		}
		result.ReviewerSession = start.Agent.SessionID
		if result.ReviewerClaim == "" {
			result.ReviewerClaim = "Independent " + start.Agent.Provider + " reviewer started by Ply"
		}
		if result.ReviewerSession == r.Request.OwnerSession || result.ReviewerSession == deliveryOwnerSession(s) {
			return workflowError(4, "independent review requires a reviewer session distinct from the owner session")
		}
		if e = d.writeOnce(source.Locator, raw); e != nil {
			return e
		}
		result.Kind, result.SchemaVersion, result.RunID, result.RequestSHA256 = "ply.workflow.review-result", 1, id, s.Result.RequestSHA256
		result.Round, result.CandidateOID, result.CandidateTree, result.RoundRequest = r.Request.Round, r.Request.CandidateOID, r.Request.CandidateTree, r.Binding
		result.Source, result.RecordedAtUTC = source, d.Now().UTC().Format(time.RFC3339Nano)
		_, e = workflowKeep(d, filepath.Join(r.Dir, "return.json"), result)
		return e
	})
	o, err := WorkflowShow(d, root, id)
	if e != nil {
		return o, e
	}
	return o, err
}

// deliveryReviewJudge decides one exact candidate across all its rounds. Any
// sealed negative return for it governs. Otherwise its latest round that the
// coordinator has not superseded decides: open rounds are awaited.
func deliveryReviewJudge(rounds []deliveryReviewRecord, oid, tree string) (string, *deliveryReviewRecord) {
	var latest, negative *deliveryReviewRecord
	for i := range rounds {
		r := &rounds[i]
		if r.Request.CandidateOID != oid || r.Request.CandidateTree != tree {
			continue
		}
		if r.Result != nil && r.Result.Decision != "passed" {
			negative = r
		}
		if r.Superseded {
			continue
		}
		latest = r
	}
	switch {
	case latest == nil:
		return "none", nil
	case negative != nil:
		return "negative", negative
	case latest.Result == nil:
		return "awaiting", latest
	}
	return "passed", latest
}

func deliveryReviewFindingIDs(r *deliveryReviewRecord) string {
	ids := []string{}
	for _, f := range r.Result.Findings {
		ids = append(ids, f.ID)
	}
	if len(ids) == 0 {
		return "no findings, see its evidence gaps"
	}
	return strings.Join(ids, ", ")
}

// deliveryNativeCandidateReview turns the passed judgment of exactly this
// candidate into the existing qualification record. Findings of earlier
// candidates stay as fixes. Another candidate never inherits a pass.
func deliveryNativeCandidateReview(d Dependencies, s workflowState, oid, tree string) (string, error) {
	rounds, e := deliveryReviewRounds(s)
	if e != nil {
		return "", e
	}
	verdict, match := deliveryReviewJudge(rounds, oid, tree)
	switch verdict {
	case "none":
		return "", workflowError(4, fmt.Sprintf("candidate %s has no independent review; the coordinator requests one with `ply workflow execute review %s`", oid, s.Result.RunID))
	case "awaiting":
		return "", workflowError(4, fmt.Sprintf("review round %d for candidate %s has not returned yet; its judgment governs this candidate", match.Request.Round, oid))
	case "negative":
		return "", workflowError(4, fmt.Sprintf("review round %d for candidate %s is %s (%s); correct them and commit a new candidate for a new review round", match.Request.Round, oid, match.Result.Decision, deliveryReviewFindingIDs(match)))
	}
	if _, e = workflowBound(*match.Return, 1<<20); e != nil {
		return "", e
	}
	if match.Result.ReviewerSession == deliveryOwnerSession(s) {
		return "", workflowError(4, "independent review requires a reviewer session distinct from the owner session")
	}
	fixes := []any{}
	for _, r := range rounds {
		if r.Request.Round < match.Request.Round && r.Result != nil && (r.Request.CandidateOID != oid || r.Request.CandidateTree != tree) {
			for _, f := range r.Result.Findings {
				fixes = append(fixes, map[string]any{"id": f.ID, "severity": f.Severity, "summary": f.Summary, "evidence_ids": []string{"candidate-review"}})
			}
		}
	}
	review := map[string]any{"kind": "DeliveryCandidateReview@1", "schema_version": 1, "candidate_oid": oid, "candidate_tree": tree, "reviewer_claim": match.Result.ReviewerClaim, "reviewer_session_id": match.Result.ReviewerSession, "decision": "passed", "findings": []any{}, "fixes": fixes, "open_actionable_findings": []any{}}
	b, e := workflowKeep(d, filepath.Join(match.Dir, "candidate-review.json"), review)
	return b.Locator, e
}

func deliveryReviewStatus(s workflowState, source DeliverySourceStatus) (*DeliveryReviewStatus, error) {
	rounds, e := deliveryReviewRounds(s)
	if e != nil || len(rounds) == 0 {
		return nil, e
	}
	out := &DeliveryReviewStatus{CurrentCandidate: "not_reviewed", Rounds: []DeliveryReviewRoundStatus{}}
	for _, r := range rounds {
		x := DeliveryReviewRoundStatus{Round: r.Request.Round, CandidateOID: r.Request.CandidateOID, CandidateTree: r.Request.CandidateTree, Coordinator: r.Request.Coordinator, State: "requested", Findings: []DeliveryReviewFinding{}, EvidenceGaps: []string{}, Previous: r.Request.Previous, Supersedes: r.Request.Supersedes, Request: r.Binding, Brief: r.Request.Brief, Context: r.Request.Context, Template: r.Request.Template, Result: r.Return}
		x.Current = source.State == "observed" && source.Clean && source.OID == x.CandidateOID && source.Tree == x.CandidateTree
		start := deliveryReviewStartState(r.Dir)
		x.ReviewerStart, x.StartDetail = start.State, start.Detail
		if start.Pane != nil {
			x.ReviewerPane = start.Pane.PaneID
		}
		if r.Superseded {
			x.State = "superseded"
		}
		if res := r.Result; res != nil {
			x.State, x.Decision, x.Summary, x.ReviewerClaim, x.ReviewerSession, x.Findings, x.EvidenceGaps = res.ReviewState, res.Decision, res.Summary, res.ReviewerClaim, res.ReviewerSession, res.Findings, res.EvidenceGaps
		}
		out.Rounds = append(out.Rounds, x)
	}
	var current *DeliveryReviewRoundStatus
	if source.State == "observed" && source.Clean {
		if _, r := deliveryReviewJudge(rounds, source.OID, source.Tree); r != nil {
			current = &out.Rounds[r.Request.Round-1]
		}
	}
	id := s.Result.RunID
	switch {
	case s.Result.Delivery.Phase == "completed":
		out.Next = WorkflowAction{"none", "Delivery is complete; review rounds remain history."}
		if current != nil && current.Decision == "passed" {
			out.CurrentCandidate = "review_passed"
		}
	case current == nil && !source.Clean:
		out.Next = WorkflowAction{"owner", "Commit a clean candidate; the coordinator then requests its review round."}
	case current == nil:
		out.Next = WorkflowAction{"coordinator", fmt.Sprintf("Candidate %s has no review; request one with `ply workflow execute review %s --coordinator <claim>`.", source.OID, id)}
	case current.State == "requested" && current.ReviewerStart == "not_started":
		out.CurrentCandidate = "review_requested"
		out.Next = WorkflowAction{"coordinator", fmt.Sprintf("Review round %d has no reviewer yet; start it with `ply workflow execute review %s --coordinator <claim>`.", current.Round, id)}
	case current.State == "requested" && current.ReviewerStart != "started":
		out.CurrentCandidate = "review_requested"
		out.Next = WorkflowAction{"coordinator", fmt.Sprintf("Review round %d stopped: %s. After inspection, `ply workflow execute review %s --coordinator <claim> --new-round` starts a fresh reviewer for this candidate.", current.Round, current.StartDetail, id)}
	case current.State == "requested":
		out.CurrentCandidate = "review_requested"
		out.Next = WorkflowAction{"reviewer", fmt.Sprintf("Review round %d is open in pane %s; the reviewer follows %s and returns with review-return.", current.Round, current.ReviewerPane, current.Brief.Locator)}
	case current.Decision == "passed":
		out.CurrentCandidate = "review_passed"
		out.Next = WorkflowAction{"owner", fmt.Sprintf("Round %d passed this exact candidate; verify without --review to use it.", current.Round)}
	case current.Decision == "changes_requested":
		out.CurrentCandidate = "findings_open"
		out.Next = WorkflowAction{"owner", fmt.Sprintf("Correct the %d finding(s) from round %d and commit; the coordinator then requests a new round.", len(current.Findings), current.Round)}
	default:
		out.CurrentCandidate = "review_blocked"
		out.Next = WorkflowAction{"coordinator", fmt.Sprintf("Round %d is blocked; resolve its stated findings or evidence gaps before a new round.", current.Round)}
	}
	return out, nil
}

func deliveryAttemptReview(s workflowState) (string, error) {
	raw, e := readFile(filepath.Join(s.Result.Delivery.Attempt.Path, "attempt.json"), 256<<10, true)
	if e != nil {
		return "", workflowError(4, "verification attempt evidence is missing; no replay is safe")
	}
	var initial deliveryVerificationReceipt
	if e = decode(raw, 256<<10, &initial); e != nil {
		return "", e
	}
	return initial.Review.Locator, nil
}
