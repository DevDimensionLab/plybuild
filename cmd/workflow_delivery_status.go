package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func writeDeliveryStatus(out io.Writer, status *taskrun.DeliveryStatus) {
	if status == nil {
		return
	}
	if s := status.Startup; s != nil {
		fmt.Fprintf(out, "Startup: terminal %s; task prompt %s; task receipt %s; target work %s\n", s.Terminal, s.TaskPrompt, s.TaskReceived, s.TargetWork)
		if s.Stop != nil {
			fmt.Fprintf(out, "Observed stop: %s\n", s.Stop.Locator)
		}
	}
	if scope := status.ScopedContinuation; scope != nil {
		fmt.Fprintf(out, "Scoped continuation: %s (original stops and unknown effects retained)\n", scope.Record.Locator)
		for _, effect := range scope.RetainedEffects {
			fmt.Fprintf(out, "Retained effect: %s — %s\n", effect.Stop.SHA256, effect.OriginalEffect)
		}
		for _, op := range scope.Operations {
			fmt.Fprintf(out, "Scoped %s: %s — %s\n", op.Operation, op.State, op.Reason)
		}
		fmt.Fprintf(out, "Cleanup: %s\n", scope.Cleanup)
	}
	short := func(oid string) string {
		if len(oid) > 12 {
			return oid[:12]
		}
		return oid
	}
	if status.Source.State == "observed" {
		fmt.Fprintf(out, "Current source: %s (clean: %t)\n", short(status.Source.OID), status.Source.Clean)
	} else {
		fmt.Fprintln(out, "Current source: unknown")
	}
	if v := status.Verification; v != nil {
		fmt.Fprintf(out, "Verification: %s for %s (%s", v.Outcome, short(v.CandidateOID), v.AttemptID)
		if v.Exit != nil {
			fmt.Fprintf(out, ", exit %d", *v.Exit)
		}
		fmt.Fprintf(out, "); qualification: %s\n", v.Qualification)
	} else {
		fmt.Fprintln(out, "Verification: not recorded")
	}
	if c := status.QualifiedCandidate; c != nil {
		current := "historical; not ready for the current source"
		if c.Current {
			current = "current"
		}
		fmt.Fprintf(out, "Qualified candidate: %s (%s)\n", short(c.OID), current)
	} else {
		fmt.Fprintln(out, "Qualified candidate: none")
	}
	if a := status.Acceptance; a != nil {
		fmt.Fprintf(out, "Acceptance: %s %s (current: %t) — %s\n", a.Mode, a.Outcome, a.Current, a.Reason)
	}
	if qa := status.HumanJudgment; qa != nil {
		fmt.Fprintf(out, "Human judgment: %s for %s (current: %t)\n", qa.Outcome, short(qa.CandidateOID), qa.Current)
	} else {
		fmt.Fprintln(out, "Human judgment: not recorded for the latest qualified candidate")
	}
	if status.FinalDelivery.State == "delivered" {
		fmt.Fprintf(out, "Final delivery: observed for %s\n", short(status.FinalDelivery.CandidateOID))
	} else {
		fmt.Fprintf(out, "Final delivery: %s\n", status.FinalDelivery.State)
	}
	if r := status.Review; r != nil {
		for _, round := range r.Rounds {
			current := "historical"
			if round.Current {
				current = "current"
			}
			fmt.Fprintf(out, "Review round %d for %s (%s): %s", round.Round, short(round.CandidateOID), current, round.State)
			if round.Decision != "" {
				fmt.Fprintf(out, " — %s", round.Decision)
			}
			fmt.Fprintln(out)
			if diag := round.StartupDiagnostics; diag != nil {
				fmt.Fprintf(out, "  Reviewer startup: %s; start=%s session=%s task_prompt=%s provider_liveness=%s\n", round.ReviewerStart, diag.StartAttempt, diag.SessionBinding, diag.TaskPrompt, diag.ProviderLiveness)
				if id := diag.Identity; id != nil {
					fmt.Fprintf(out, "  Herdr: %s / %s / %s / %s\n", id.WorkspaceID, id.TabID, id.PaneID, id.TerminalID)
				}
				if f := diag.FirstFailure; f != nil {
					fmt.Fprintf(out, "  First failure: %s class=%s exit=%s provider_code=%s stdout_bytes=%d stderr_bytes=%d truncated=%t\n", f.Phase, f.Class, f.Exit, f.ProviderCode, f.StdoutBytes, f.StderrBytes, f.Truncated)
					if b := diag.FirstFailureEvidence; b != nil {
						fmt.Fprintf(out, "  Failure evidence: %s (%s)\n", b.Locator, b.SHA256)
					}
				}
				if o := diag.LastObservation; o != nil {
					fmt.Fprintf(out, "  Session observation: %s", o.Outcome)
					if o.Failure != nil {
						fmt.Fprintf(out, "; class=%s provider_code=%s", o.Failure.Class, o.Failure.ProviderCode)
					}
					if o.FirstFailure != nil {
						fmt.Fprintf(out, "; first_read_failure=%s provider_code=%s", o.FirstFailure.Class, o.FirstFailure.ProviderCode)
					}
					if o.SessionID != "" {
						fmt.Fprintf(out, "; last_observed_session=%s agent_status=%s", o.SessionID, o.AgentStatus)
					}
					fmt.Fprintln(out)
				}
				if round.ReviewerStart != "started" {
					for _, argv := range diag.InspectionArgv {
						quoted := make([]string, len(argv))
						for i, arg := range argv {
							quoted[i] = taskrun.ShellQuote(arg)
						}
						fmt.Fprintln(out, "  Inspect (read-only): "+strings.Join(quoted, " "))
					}
				}
			}
			for _, f := range round.Findings {
				fmt.Fprintf(out, "  Finding %s [%s]: %s\n", f.ID, f.Severity, f.Summary)
			}
			for _, g := range round.EvidenceGaps {
				fmt.Fprintf(out, "  Evidence gap: %s\n", g)
			}
		}
		fmt.Fprintf(out, "Review next: %s — %s\n", r.Next.Actor, r.Next.Message)
	}
	for _, reason := range status.Reasons {
		fmt.Fprintln(out, "Needs attention: "+reason.Detail)
	}
}
