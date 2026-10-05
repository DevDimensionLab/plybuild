package workflowhandoff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func deliveryCandidateReviewFixture(t *testing.T, in DeliveryCandidateInput, change func(map[string]any)) []byte {
	t.Helper()
	raw, err := os.ReadFile(in.ReviewPath)
	if err != nil {
		t.Fatal(err)
	}
	var review map[string]any
	if err = json.Unmarshal(raw, &review); err != nil {
		t.Fatal(err)
	}
	change(review)
	deliveryFixtureWrite(t, in.ReviewPath, review)
	raw, err = os.ReadFile(in.ReviewPath)
	if err != nil {
		t.Fatal(err)
	}
	// This remains a faithful synthetic verifier receipt: bind the exact review
	// presented by this fixture before qualification, not an unrelated digest.
	receiptRaw, err := os.ReadFile(in.VerificationPath)
	if err != nil {
		t.Fatal(err)
	}
	var receipt map[string]any
	if err = json.Unmarshal(receiptRaw, &receipt); err != nil {
		t.Fatal(err)
	}
	receipt["review"] = map[string]any{"locator": in.ReviewPath, "sha256": digestBytes(raw)}
	deliveryFixtureWrite(t, in.VerificationPath, receipt)
	return raw
}

func deliveryHandoffDirectories(t *testing.T, locator string) []string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Dir(filepath.Dir(locator)))
	if err != nil {
		t.Fatal(err)
	}
	names := []string{}
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

func TestDeliveryReviewUnknownEvidenceRejectedBeforeCandidateEffects(t *testing.T) {
	for _, list := range []string{"findings", "fixes"} {
		t.Run(list, func(t *testing.T) {
			f := prepareDeliveryFixture(t)
			in := f.candidate(t, "unknown-evidence")
			deliveryCandidateReviewFixture(t, in, func(review map[string]any) {
				review[list] = []any{map[string]any{"id": "observed-item", "severity": "low", "summary": "Synthetic real review item; preserve it when its evidence is corrected.", "evidence_ids": []string{"review-notes-r1"}}}
			})
			before := deliveryHandoffDirectories(t, f.locator)
			registry, err := f.w.WorkItems.Snapshot(f.root)
			if err != nil {
				t.Fatal(err)
			}
			owner, err := f.d.Store.ReadByLocator(f.locator)
			if err != nil {
				t.Fatal(err)
			}
			workarea := filepath.Join(owner.Handoff.ReplyRoot, "staging", "candidate-"+digestBytes([]byte(in.CandidateKey))[7:])
			got, err := QualifyDeliveryCandidate(f.d, in)
			if err == nil || !strings.Contains(err.Error(), list) || !strings.Contains(err.Error(), "evidence_ids") || !strings.Contains(err.Error(), "review-notes-r1") {
				t.Errorf("missing specific unsupported evidence reference error: %v", err)
			}
			if got.Handoff.Locator != "" || got.Terminal.Locator != "" || !reflect.DeepEqual(before, deliveryHandoffDirectories(t, f.locator)) {
				t.Errorf("invalid review left a candidate handoff/terminal: %+v", got)
			}
			if _, err = os.Lstat(workarea); !os.IsNotExist(err) {
				t.Errorf("invalid review left candidate staging: %v", err)
			}
			after, err := f.w.WorkItems.Snapshot(f.root)
			if err != nil || after.RawSHA256 != registry.RawSHA256 || len(after.TaskResults) != 0 {
				t.Errorf("invalid review changed TaskResult state: %+v %v", after.TaskResults, err)
			}
		})
	}
}

func TestDeliveryReviewKnownEvidencePreservesNonemptyAssessment(t *testing.T) {
	f := prepareDeliveryFixture(t)
	in := f.candidate(t, "nonempty-review")
	raw := deliveryCandidateReviewFixture(t, in, func(review map[string]any) {
		review["findings"] = []any{map[string]any{"id": "finding-a", "severity": "low", "summary": "The synthetic reviewer found a wording issue before candidate completion.", "evidence_ids": []string{"candidate-review"}}}
		review["fixes"] = []any{map[string]any{"id": "fix-a", "severity": "low", "summary": "The wording was corrected and reviewed in this exact candidate.", "evidence_ids": []string{"candidate-review", "verifier-stdout"}}}
	})
	got, err := QualifyDeliveryCandidate(f.d, in)
	if err != nil {
		t.Fatal(err)
	}
	want := workspace.TaskReviewRecord{
		Findings:               []workspace.TaskReviewEntry{{ID: "finding-a", Severity: "low", Summary: "The synthetic reviewer found a wording issue before candidate completion.", EvidenceArtifactIDs: []string{"candidate-review"}}},
		Fixes:                  []workspace.TaskReviewEntry{{ID: "fix-a", Severity: "low", Summary: "The wording was corrected and reviewed in this exact candidate.", EvidenceArtifactIDs: []string{"candidate-review", "verifier-stdout"}}},
		OpenActionableFindings: []workspace.TaskReviewEntry{},
	}
	if !reflect.DeepEqual(got.TaskResult.Review, want) {
		t.Fatalf("nonempty findings/fixes were changed or hidden: %+v", got.TaskResult.Review)
	}
	preserved := false
	for _, artifact := range got.TaskResult.Artifacts {
		if artifact.ArtifactID == "candidate-review" {
			preserved = artifact.SHA256 == digestBytes(raw)
		}
	}
	if !preserved {
		t.Fatal("exact original review was not preserved as managed evidence")
	}
	r, err := f.w.WorkItems.Snapshot(f.root)
	if err != nil || len(r.TaskResults) != 1 || !reflect.DeepEqual(r.TaskResults[0].Review, want) {
		t.Fatalf("nonempty review was not published: %+v %v", r.TaskResults, err)
	}
}

func TestDeliveryReviewClaimLengthNamesTheField(t *testing.T) {
	in := DeliveryCandidateInput{CandidateOID: strings.Repeat("a", 40), CandidateTree: strings.Repeat("b", 40)}
	template, err := DeliveryCandidateReviewTemplate(in.CandidateOID, in.CandidateTree)
	if err != nil {
		t.Fatal(err)
	}
	var review map[string]any
	if err = json.Unmarshal(template, &review); err != nil {
		t.Fatal(err)
	}
	review["decision"], review["reviewer_claim"] = "passed", strings.Repeat("r", 305)
	raw, _ := json.Marshal(review)
	if _, err = deliveryCandidateReview(raw, in, "owner-session"); err == nil || !strings.Contains(err.Error(), "reviewer_claim") || !strings.Contains(err.Error(), "256") {
		t.Fatalf("reviewer field length rejected as an unrelated candidate mismatch: %v", err)
	}
}
