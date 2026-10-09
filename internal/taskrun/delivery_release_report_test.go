package taskrun

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func releaseReportCandidate(t *testing.T) (deliveryFixture, WorkflowRun, string, string) {
	t.Helper()
	a := workspace.DeliveryAgreement{SchemaVersion: 2, Mode: workspace.DeliveryLocalBranch, IntegrationOwner: workspace.IntegrationOwnerHuman}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &a, ParentRef: "main"})
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "qualified source handover")
	counter := filepath.Join(f.R.WorkspaceRoot, "acceptance-count.txt")
	script, err := os.ReadFile(f.AcceptancePath)
	if err != nil {
		t.Fatal(err)
	}
	script = append(script, []byte("printf 'executed\\n' >> "+ShellQuote(counter)+"\n")...)
	if err = os.WriteFile(f.AcceptancePath, script, 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "source-handover")
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review)
	if err != nil {
		t.Fatal(err)
	}
	return f, o, review, counter
}

func TestDeliveryQualifiedReportReleaseStillRejectsChangedAuthority(t *testing.T) {
	for _, change := range []string{"dirty-source", "new-source", "acceptance", "review", "publication", "verification", "negative-verification", "unknown-verification", "negative-human", "wrong-context", "changed-context", "wrong-cwd", "wrong-session"} {
		t.Run(change, func(t *testing.T) {
			f, o, review, counter := releaseReportCandidate(t)
			c := o.Delivery.Candidates[0]
			o = releaseProgressReport(t, f, o, "working", "qualified-progress")
			if _, err := nativeDeliveryCandidate(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID); err != nil {
				t.Fatalf("negative fixture lacks the positive qualification it intends to invalidate: %v", err)
			}
			d, context := f.D, o.Paths.Context
			write := func(path string, data []byte) {
				t.Helper()
				if err := os.WriteFile(path, data, 0600); err != nil {
					t.Fatal(err)
				}
			}
			appendFile := func(path string) {
				t.Helper()
				raw, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				write(path, append(raw, '\n'))
			}
			switch change {
			case "dirty-source", "new-source":
				cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
				write(filepath.Join(cwd, "README.md"), []byte("source changed after qualification\n"))
				if change == "new-source" {
					runGit(t, cwd, "add", "README.md")
					runGit(t, cwd, "commit", "-m", "unverified source correction")
				}
			case "acceptance":
				appendFile(f.AcceptancePath)
			case "review":
				appendFile(review)
			case "verification":
				appendFile(c.Verification.Locator)
			case "publication":
				registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
				if err != nil {
					t.Fatal(err)
				}
				removed := false
				for _, p := range registry.TaskContentPublications {
					if p.TaskID == c.TaskResult.TaskID {
						path := filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "requests", "sha256", strings.TrimPrefix(p.RequestSHA256, "sha256:")+".json")
						if err := os.Remove(path); err != nil {
							t.Fatal(err)
						}
						removed = true
						break
					}
				}
				if !removed {
					t.Fatal("fixture has no required native publication")
				}
			case "negative-verification", "unknown-verification":
				verifier := f.D
				if change == "negative-verification" {
					write(f.AcceptancePath, []byte("#!/bin/sh\nprintf 'executed\\n' >> "+ShellQuote(counter)+"\nexit 7\n"))
				} else {
					verifier.Fault = func(point string) error {
						if point == "delivery_after_verification_reservation" {
							return errors.New("synthetic unknown verification after reservation")
						}
						return nil
					}
				}
				var err error
				o, err = WorkflowDeliveryVerify(verifier, f.R.WorkspaceRoot, o.RunID, context, review)
				if err == nil {
					t.Fatal("new negative or unknown verification unexpectedly qualified")
				}
				question := "Who will inspect this unresolved verification without retrying it?"
				r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "preserve-later-verification", PreviousEventSHA256: o.Delivery.LastEventSHA256, Phase: "needs_input", Summary: "A later verification is not qualified.", Meaning: "Preserve the exact outcome and stop source handover.", Question: &question, Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
				o, err = WorkflowDeliveryIncompleteReport(f.D, f.R.WorkspaceRoot, o.RunID, context, writeAny(t, f.R.WorkspaceRoot, "later-verification.json", r))
				if err != nil {
					t.Fatal(err)
				}
			case "negative-human":
				var err error
				o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, context, "fail", deliveryTestHuman(t, f, o, "fail"))
				if err != nil {
					t.Fatal(err)
				}
				o = releaseProgressReport(t, f, o, "needs_input", "after-negative-human")
			case "wrong-context":
				context = writeAny(t, f.R.WorkspaceRoot, "other-context.json", map[string]any{})
			case "changed-context":
				appendFile(context)
			case "wrong-cwd":
				d.CWD = func() (string, error) { return f.Parent, nil }
			case "wrong-session":
				workflowTestModel(t, f.workflowFixture, map[string]any{"agent_session_id": "another-native-session"})
			}
			statePath := filepath.Join(o.Paths.RunRoot, "state.json")
			before := hashFileTest(t, statePath)
			if _, err := ReleaseDeliveryOwnership(d, f.R.WorkspaceRoot, o.RunID, context, "Synthetic prohibited source handover."); err == nil {
				t.Fatalf("%s failed to prevent source ownership release", change)
			}
			if hashFileTest(t, statePath) != before {
				t.Fatal("rejected release modified native workflow state")
			}
		})
	}
}

func releaseProgressReport(t *testing.T, f deliveryFixture, o WorkflowRun, phase, id string) WorkflowRun {
	t.Helper()
	r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID,
		EventID: id, PreviousEventSHA256: o.Delivery.LastEventSHA256, Phase: phase, Summary: "The exact qualified candidate is ready for source handover.",
		Meaning: "Progress only; the candidate, acceptance and review are unchanged. Human QA remains separate.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	if phase == "needs_input" {
		r.Question = ptr("Which supported source handover should this unchanged candidate use?")
	}
	out, err := WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, id+".json", r))
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func assertReleaseCandidatePreserved(t *testing.T, f deliveryFixture, o WorkflowRun, candidate DeliveryCandidate, counter, parent string) {
	t.Helper()
	if len(o.Delivery.Candidates) != 1 || !equal(o.Delivery.Candidates[0], candidate) || o.Delivery.Attempt.ID != candidate.Key || o.Delivery.Attempt.State != "recorded" {
		t.Fatalf("handover replaced candidate, TaskResult, receipt or attempt: %+v", o.Delivery)
	}
	raw, err := os.ReadFile(counter)
	if err != nil || string(raw) != "executed\n" {
		t.Fatalf("acceptance must execute exactly once: %q %v", raw, err)
	}
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(registry.TaskResults) != 1 || len(registry.HumanQARecords) != 0 {
		t.Fatalf("handover created a result or human judgment: results=%d QA=%d error=%v", len(registry.TaskResults), len(registry.HumanQARecords), err)
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != parent {
		t.Fatal("source release integrated the candidate")
	}
}

func TestDeliveryQualifiedReportReleasePreservesExactCandidate(t *testing.T) {
	for _, phases := range [][]string{{"working"}, {"needs_input"}, {"working", "needs_input"}} {
		for _, reuse := range []bool{false, true} {
			t.Run(strings.Join(phases, "-")+fmt.Sprintf("/reuse=%t", reuse), func(t *testing.T) {
				f, o, review, counter := releaseReportCandidate(t)
				candidate := o.Delivery.Candidates[0]
				parent := gitOutput(t, f.Parent, "rev-parse", "HEAD")
				requestSHA, contextSHA := hashFileTest(t, f.File), hashFileTest(t, o.Paths.Context)
				control := f.R.Runtime.PlyExecutable
				for i, phase := range phases {
					o = releaseProgressReport(t, f, o, phase, fmt.Sprintf("progress-%d", i))
				}
				reports := append([]DeliveryEventRecord{}, o.Delivery.Events...)
				if reuse {
					var err error
					o, err = WorkflowDeliveryRequalify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, candidate.Key)
					if err != nil {
						t.Fatal(err)
					}
					if o.Delivery.Phase != phases[len(phases)-1] || !equal(o.Delivery.Events, reports) {
						t.Fatal("receipt reuse rewrote progress phase or historical reports")
					}
				}
				// Release intentionally supports a newer ordinary CLI while old
				// controls and the reports they wrote remain immutable.
				newCLI := filepath.Join(f.R.WorkspaceRoot, "new-ordinary-ply")
				if err := os.WriteFile(newCLI, []byte("#!/bin/sh\n# synthetic newer release-only CLI\nexit 0\n"), 0700); err != nil {
					t.Fatal(err)
				}
				newer := f.D
				newer.Executable = func() (string, error) { return newCLI, nil }
				released, err := ReleaseDeliveryOwnership(newer, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic owner has stopped writing; preserve this exact candidate for separate human QA.")
				if err != nil {
					t.Fatalf("progress invalidated unchanged technical qualification (reuse=%t): %v", reuse, err)
				}
				if released.RunID != o.RunID || released.Delivery.OwnershipRelease == nil || len(released.Delivery.Events) != len(reports)+1 || !equal(released.Delivery.Events[:len(reports)], reports) {
					t.Fatal("release lost the original Run, report history or source revocation")
				}
				assertReleaseCandidatePreserved(t, f, released, candidate, counter, parent)
				if hashFileTest(t, f.File) != requestSHA || hashFileTest(t, o.Paths.Context) != contextSHA || hashFileTest(t, control.Path) != control.SHA256 {
					t.Fatal("newer release rewrote frozen request, context or control")
				}
				again, err := ReleaseDeliveryOwnership(newer, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Repeat the same already preserved source release.")
				if err != nil || !equal(again.Delivery.OwnershipRelease, released.Delivery.OwnershipRelease) || len(again.Delivery.Events) != len(released.Delivery.Events) {
					t.Fatalf("release retry duplicated the revocation: %v", err)
				}
				if _, err = WorkflowDeliveryRequalify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, candidate.Key); err == nil || !strings.Contains(err.Error(), "owner_released") {
					t.Fatalf("receipt reuse revived the released owner: %v", err)
				}
				if err = CheckAgentDeliveryExecution(f.D, f.R.WorkspaceRoot, o.RunID); err == nil {
					t.Fatal("source release granted agent integration authority")
				}
			})
		}
	}
}
