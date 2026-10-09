package taskrun

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestDeliveryCLIReleaseReadsHistoricalQualifiedReports(t *testing.T) {
	// af5ad56 is the installed .4 source that overwrites conversational phase
	// on report. Its real callbacks write all candidate/receipt/report evidence;
	// only the release command runs the newer candidate binary.
	const historicalSource = "af5ad56cf1b799a241cdfde8b5a327ea30714165"
	old := continuitySourceBinary(t, historicalSource, "PLY_RELEASE_REPORT_OLD_BINARY")
	newer := automaticCLIBinary(t)
	t.Logf("historical source=%s binary=%s; candidate release binary=%s", historicalSource, hashFileTest(t, old), hashFileTest(t, newer))
	for _, reuse := range []bool{false, true} {
		t.Run(fmt.Sprintf("reuse=%t", reuse), func(t *testing.T) {
			a := workspace.DeliveryAgreement{SchemaVersion: 2, Mode: workspace.DeliveryLocalBranch, IntegrationOwner: workspace.IntegrationOwnerHuman}
			f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &a, ParentRef: "main"})
			f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
			f.R.Runtime.PlyExecutable = Executable{Path: old, SHA256: hashFileTest(t, old)}
			f.D.Executable = func() (string, error) { return old, nil }
			f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
			o := workflowTestStart(t, f.workflowFixture)
			acceptance := writeAny(t, f.R.WorkspaceRoot, "cli-runtime-acceptance.json", deliveryTestAcceptance(t, f, o))
			o = automaticCLIOK(t, old, f, "run", "accept", o.RunID, "--context", o.Paths.Context, "--file", acceptance)
			deliveryTestCandidate(t, f, "historical source release")
			counter := filepath.Join(f.R.WorkspaceRoot, "acceptance-count.txt")
			script, err := os.ReadFile(f.AcceptancePath)
			if err != nil {
				t.Fatal(err)
			}
			script = append(script, []byte("printf 'executed\\n' >> "+ShellQuote(counter)+"\n")...)
			if err = os.WriteFile(f.AcceptancePath, script, 0600); err != nil {
				t.Fatal(err)
			}
			review := deliveryTestReview(t, f, "historical-cli")
			o = automaticCLIOK(t, old, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review)
			c := o.Delivery.Candidates[0]
			requestSHA, contextSHA, controlSHA := hashFileTest(t, f.File), hashFileTest(t, o.Paths.Context), hashFileTest(t, old)
			parent := gitOutput(t, f.Parent, "rev-parse", "HEAD")
			for i, phase := range []string{"working", "needs_input"} {
				r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID,
					EventID: fmt.Sprintf("historical-progress-%d", i), PreviousEventSHA256: o.Delivery.LastEventSHA256, Phase: phase,
					Summary: "The same qualified candidate is ready for source handover.", Meaning: "Preserved progress only; no change to candidate, review, acceptance or human QA.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
				if phase == "needs_input" {
					r.Question = ptr("Which supported owner operation releases this exact qualified candidate?")
				}
				file := writeAny(t, f.R.WorkspaceRoot, r.EventID+".json", r)
				o = automaticCLIOK(t, old, f, "execute", "report", o.RunID, "--context", o.Paths.Context, "--file", file)
			}
			history := append([]DeliveryEventRecord{}, o.Delivery.Events...)
			if reuse {
				o = automaticCLIOK(t, old, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review, "--reuse", c.Key)
			}
			if o.Delivery.Phase != "needs_input" || !equal(o.Delivery.Events, history) {
				t.Fatal("historical CLI did not preserve the reported failure sequence")
			}
			state := filepath.Join(o.Paths.RunRoot, "state.json")
			before := hashFileTest(t, state)
			if _, output, err := automaticCLI(old, f, "execute", "release", o.RunID, "--context", o.Paths.Context, "--reason", "Synthetic owner source handover."); err == nil || !strings.Contains(string(output), "delivery candidate has not reached technical qualification") {
				t.Fatalf("historical release no longer exhibits the consumer failure: %v %s", err, output)
			}
			if hashFileTest(t, state) != before {
				t.Fatal("historical failed release modified native state")
			}
			released := automaticCLIOK(t, newer, f, "execute", "release", o.RunID, "--context", o.Paths.Context, "--reason", "Synthetic owner stopped writing; retain the same candidate and separate human QA.")
			if released.Delivery.OwnershipRelease == nil || released.RunID != o.RunID || len(released.Delivery.Events) != len(history)+1 || !equal(released.Delivery.Events[:len(history)], history) {
				t.Fatal("new ordinary release lost the historical Run, report chain or ownership event")
			}
			assertReleaseCandidatePreserved(t, f, released, c, counter, parent)
			if hashFileTest(t, f.File) != requestSHA || hashFileTest(t, o.Paths.Context) != contextSHA || hashFileTest(t, old) != controlSHA {
				t.Fatal("new ordinary release replaced historical request, context or frozen control")
			}
			if _, output, err := automaticCLI(old, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review, "--reuse", c.Key); err == nil || !strings.Contains(string(output), "owner_released") {
				t.Fatalf("old controller acquired authority after release: %v %s", err, output)
			}
		})
	}
}
