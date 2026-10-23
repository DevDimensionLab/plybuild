package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParallelSourceLostMergeReplyWithoutHeadChangeCannotReplay(t *testing.T) {
	f, prepared, in := sourceReconcileFixture(t, false)
	preview, e := ReconcileTaskSource(f.dependencies, in)
	if e != nil {
		t.Fatal(e)
	}
	git := f.dependencies.IntegrationGit.(*systemTaskIntegrationGit)
	runner, merges := git.run, 0
	git.run = func(args, env []string) GitCommandOutcome {
		for _, arg := range args {
			if arg == "merge" {
				merges++
				if merges == 1 {
					return GitCommandOutcome{Exit: -1, Err: errInjected, Stderr: []byte("simulated lost Git effect reply")}
				}
			}
		}
		return runner(args, env)
	}
	f.dependencies.SourceReconciliationFault = func(stage, id string) error {
		if stage == "after_git" {
			return errInjected
		}
		return nil
	}
	in.Apply, in.Confirmation = true, preview.Confirmation
	if _, e = ReconcileTaskSource(f.dependencies, in); e == nil {
		t.Fatal("publication interruption not observed")
	}
	f.dependencies.SourceReconciliationFault = nil
	source := prepared.Preparation.Plan.WorktreePath
	before := runLocalGit(t, source, "rev-parse", "HEAD")
	result, e := ReconcileTaskSource(f.dependencies, in)
	t.Logf("simulated lost reply: retries requested merges=%d state=%s err=%v", merges, result.State, e)
	if e == nil || result.State != "effect_unknown" || merges != 1 || runLocalGit(t, source, "rev-parse", "HEAD") != before {
		t.Fatal("ambiguous attempted merge was replayed without positive unstarted evidence")
	}
}

func TestParallelCancelUnstartedRefusesUnprovenOrForeignEffect(t *testing.T) {
	for _, scenario := range []string{"wrong-run", "wrong-request", "wrong-mandate", "wrong-operation", "no-operation", "dirty", "source-moved", "wrong-ref", "marker", "mismatched-marker", "missing-positive-fact", "attempted-missing-marker", "effect-unknown", "conflict"} {
		t.Run(scenario, func(t *testing.T) {
			f, p, in := sourceReconcileFixture(t, scenario == "conflict")
			preview, e := ReconcileTaskSource(f.dependencies, in)
			if e != nil {
				t.Fatal(e)
			}
			in.Apply, in.Confirmation = true, preview.Confirmation
			stage := "after_intent"
			if scenario == "attempted-missing-marker" || scenario == "effect-unknown" || scenario == "conflict" {
				stage = "after_git"
			}
			f.dependencies.SourceReconciliationFault = func(at, id string) error {
				if at == stage {
					return errInjected
				}
				return nil
			}
			intent, e := ReconcileTaskSource(f.dependencies, in)
			if e == nil || intent.Operation == nil {
				t.Fatalf("fixture fault: %+v %v", intent, e)
			}
			f.dependencies.SourceReconciliationFault = nil
			source := p.Preparation.Plan.WorktreePath
			marker := runLocalGit(t, source, "rev-parse", "--path-format=absolute", "--git-path", "ply-source-"+intent.Operation.ID)
			if scenario == "attempted-missing-marker" || scenario == "effect-unknown" {
				if e = os.Remove(marker); e != nil {
					t.Fatal(e)
				}
			}
			if scenario == "effect-unknown" || scenario == "conflict" {
				_, _ = ReconcileTaskSource(f.dependencies, in)
			}
			in.OperationID, in.Action, in.Confirmation = intent.Operation.ID, "cancel-unstarted", ""
			switch scenario {
			case "wrong-run":
				in.RunID = "foreign-run"
			case "wrong-request":
				in.RequestSHA256 = "sha256:" + strings.Repeat("a", 64)
			case "wrong-mandate":
				in.MandateSHA256 = "sha256:" + strings.Repeat("a", 64)
			case "wrong-operation":
				in.OperationID = "src_" + strings.Repeat("a", 64)
			case "no-operation":
				in.OperationID = ""
			case "dirty":
				if e = os.WriteFile(filepath.Join(source, "README.md"), []byte("owner edits\n"), 0600); e != nil {
					t.Fatal(e)
				}
			case "source-moved":
				runLocalGit(t, source, "commit", "--allow-empty", "-m", "Owner source moved")
			case "wrong-ref":
				runLocalGit(t, source, "checkout", "-b", "foreign-source")
			case "marker", "mismatched-marker":
				b := []byte(intent.Operation.PlanSHA256 + "\n")
				if scenario == "mismatched-marker" {
					b = []byte("unknown marker\n")
				}
				if e = os.WriteFile(marker, b, 0600); e != nil {
					t.Fatal(e)
				}
			case "missing-positive-fact":
				// Model historical intent without the new positive native fact.
				e = f.dependencies.WorkItems.WithLock(filepath.Dir(f.wrapper), func(s WorkItemStoreSession) error {
					r, e := s.Snapshot()
					if e != nil {
						return e
					}
					raw, e := json.Marshal(r.SourceReconciliations[0])
					if e != nil {
						return e
					}
					var fields map[string]json.RawMessage
					if e = json.Unmarshal(raw, &fields); e != nil {
						return e
					}
					delete(fields, "unstarted")
					raw, e = json.Marshal(fields)
					if e != nil {
						return e
					}
					r.SourceReconciliations[0] = TaskSourceReconciliation{}
					if e = json.Unmarshal(raw, &r.SourceReconciliations[0]); e != nil {
						return e
					}
					return s.Publish(r)
				})
				if e != nil {
					t.Fatal(e)
				}
			}
			before := queueStateBytes(t, f)
			head := runLocalGit(t, source, "rev-parse", "HEAD")
			markerBefore, markerErr := os.ReadFile(marker)
			result, e := ReconcileTaskSource(f.dependencies, in)
			if e == nil || result.State == "aborted" || queueStateBytes(t, f) != before || runLocalGit(t, source, "rev-parse", "HEAD") != head {
				t.Fatalf("unproven cancellation changed owned source/history: %+v %v", result, e)
			}
			markerAfter, afterErr := os.ReadFile(marker)
			if (markerErr == nil) != (afterErr == nil) || string(markerBefore) != string(markerAfter) {
				t.Fatal("refused cancellation changed Git marker")
			}
			t.Logf("isolated %s cancellation refused without mutation: %v", scenario, e)
		})
	}
}
