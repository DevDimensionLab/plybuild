package cmd

import (
	"fmt"
	"io"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func writeDeliveryStatus(out io.Writer, status *taskrun.DeliveryStatus) {
	if status == nil {
		return
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
	for _, reason := range status.Reasons {
		fmt.Fprintln(out, "Needs attention: "+reason.Detail)
	}
}
