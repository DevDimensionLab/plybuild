package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

func TestReturnClosureLegacyRejectionKeepsHistoryAndConcreteDiagnostic(t *testing.T) {
	d, r, f := fixture(t)
	runner := d.Runner.(*fakeRunner)
	runner.inside = func(s LaunchSpec) error {
		if _, err := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "accept.json", acceptance(t, d, r))); err != nil {
			return err
		}
		report := semanticReport(t, d, r)
		report.ObservedEffects = bytes.ReplaceAll(report.ObservedEffects, []byte(`"directory_prefixes":[],`), nil)
		raw, _ := Canonical(report)
		j, _ := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
		draft, err := workflowhandoff.BuildTaskRunTerminal(d.Workflow, j.Binding.Handoff.Locator, raw, 4)
		if err != nil {
			return err
		}
		// Fixture of an already reserved legacy report: production must never
		// rewrite it, even though new reports reject this before reservation.
		accepted := filepath.Join(runPaths(r).RunRoot, "reports", "accepted.json")
		if err = d.writeOnce(accepted, raw); err != nil {
			return err
		}
		if err = d.writeOnce(terminalDraftPath(r, hash(raw)), draft); err != nil {
			return err
		}
		path := writeAny(t, r.WorkspaceRoot, "legacy-report.json", report)
		out, err := SubmitReport(d, r.WorkspaceRoot, RunID(r.RequestKey), path)
		if err == nil || !strings.Contains(err.Error(), "observed effect") || out.Delivery.State != "not_representable" {
			t.Fatalf("missing concrete rejection: %v %+v", err, out.Delivery)
		}
		before := treeState(t, runPaths(r).RunRoot)
		_, _ = SubmitReport(d, r.WorkspaceRoot, RunID(r.RequestKey), path)
		if before != treeState(t, runPaths(r).RunRoot) {
			t.Fatal("retry changed legacy history")
		}
		preserved, _ := os.ReadFile(accepted)
		if !bytes.Equal(preserved, raw) {
			t.Fatal("legacy accepted report rewritten")
		}
		return nil
	}
	p, err := PreviewStart(d, f)
	if err != nil {
		t.Fatal(err)
	}
	out, _ := Start(d, f, *p.(Preview).Confirmation)
	if runner.lastErr != nil {
		t.Fatal(runner.lastErr)
	}
	if out.Collection.State == "qualified" {
		t.Fatal("legacy rejected terminal qualified")
	}
}

func TestReturnClosureInvalidReportCanBeCorrected(t *testing.T) {
	for _, mode := range []string{"missing-empty-prefixes", "scope-drift", "verifier-target", "unknown-budget"} {
		t.Run(mode, func(t *testing.T) {
			d, r, f := fixture(t)
			runner := d.Runner.(*fakeRunner)
			runner.inside = func(s LaunchSpec) error {
				if _, err := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "accept.json", acceptance(t, d, r))); err != nil {
					return err
				}
				good := semanticReport(t, d, r)
				bad := good
				want := "observed effect"
				switch mode {
				case "missing-empty-prefixes", "scope-drift":
					var effects []map[string]any
					json.Unmarshal(good.ObservedEffects, &effects)
					scope := effects[0]["scope"].(map[string]any)
					if mode == "missing-empty-prefixes" {
						delete(scope, "directory_prefixes")
					} else {
						scope["paths"] = []string{"outside.py"}
					}
					bad.ObservedEffects, _ = json.Marshal(effects)
				case "verifier-target":
					var results []map[string]any
					json.Unmarshal(good.VerifierResults, &results)
					results[0]["bound_oid_or_sha256"] = strings.Repeat("a", 40)
					bad.VerifierResults, _ = json.Marshal(results)
					want = "wrong target binding"
				case "unknown-budget":
					bad.BudgetUsage = nil
					want = "budget"
				}
				_, err := SubmitReport(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "bad.json", bad))
				if err == nil || !strings.Contains(err.Error(), want) {
					t.Fatalf("expected concrete %s rejection, got %v", want, err)
				}
				if _, err = os.Stat(filepath.Join(runPaths(r).RunRoot, "reports", "accepted.json")); !os.IsNotExist(err) {
					t.Fatal("invalid input occupied immutable report slot")
				}
				path := writeAny(t, r.WorkspaceRoot, "good.json", good)
				out, err := SubmitReport(d, r.WorkspaceRoot, RunID(r.RequestKey), path)
				if err != nil || out.Delivery.State != "received" || out.Delivery.TerminalSHA256 == nil {
					t.Fatalf("correction failed: %v %+v", err, out.Delivery)
				}
				before := treeState(t, runPaths(r).RunRoot)
				if _, err = SubmitReport(d, r.WorkspaceRoot, RunID(r.RequestKey), path); err != nil {
					return err
				}
				good.Summary = "Contradictory later report."
				if _, err = SubmitReport(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "different.json", good)); err == nil {
					t.Fatal("accepted report replaced")
				}
				if before != treeState(t, runPaths(r).RunRoot) {
					t.Fatal("accepted history changed")
				}
				return nil
			}
			p, err := PreviewStart(d, f)
			if err != nil {
				t.Fatal(err)
			}
			_, _ = Start(d, f, *p.(Preview).Confirmation)
			if runner.lastErr != nil {
				t.Fatal(runner.lastErr)
			}
		})
	}
}
