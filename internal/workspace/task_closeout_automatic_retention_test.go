package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAutomaticCloseoutRetentionPreservesExactOriginalReferences(t *testing.T) {
	for _, defect := range []string{"", "before integration", "historical agreement", "different candidate", "tampered binary snapshot", "tampered review artifact"} {
		t.Run(defect, func(t *testing.T) {
			root, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			write := func(name string, raw []byte) closeoutSnapshotBinding {
				t.Helper()
				path := filepath.Join(root, name)
				if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, raw, 0600); err != nil {
					t.Fatal(err)
				}
				return closeoutSnapshotBinding{Locator: path, SHA256: digestTaskBytes(raw)}
			}
			scriptBytes, binaryBytes, reviewBytes := []byte("#!/bin/sh\nexit 0\n"), []byte("fixture executable bytes"), []byte("fixture review bytes")
			script := write("attempt/acceptance.sh", scriptBytes)
			binary := write("attempt/candidate-executable", binaryBytes)
			review := write("candidate/candidate-review", reviewBytes)
			original := func(name string, snapshot closeoutSnapshotBinding) closeoutSnapshotBinding {
				return closeoutSnapshotBinding{Locator: filepath.Join(root, name), SHA256: snapshot.SHA256}
			}
			originalScript, originalBinary, originalReview := original("original-acceptance.sh", script), original("original-binary", binary), original("original-review.json", review)
			tr := automaticResultFixture()
			tr.RunID, tr.SourceLocator = "candidate-run", filepath.Join(root, "task")
			a := automaticAgreementFixture()
			a.SourceRef = "refs/heads/task"
			auth := &DeliveryAuthorization{Agreement: a, AgreementSHA256: DeliveryAgreementDigest(a), RunID: "owner-run", CandidateRunID: tr.RunID, RequestSHA256: "sha256:" + strings.Repeat("d", 64), MandateSHA256: "sha256:" + strings.Repeat("e", 64), PermissionConfirmed: true, AllowedEffects: a.AllowedEffects()}
			if err := ValidateDeliveryAuthorization(*auth); err != nil {
				t.Fatal(err)
			}
			verification := map[string]any{
				"kind": "PlyDeliveryVerification@2", "schema_version": 2, "run_id": auth.RunID, "request_sha256": auth.RequestSHA256,
				"candidate_oid": tr.ResultOID, "candidate_tree": tr.ResultTree, "cwd": tr.SourceLocator,
				"acceptance": originalScript, "acceptance_snapshot": script, "review": originalReview, "outcome": "pass", "exit": 0,
				"automatic": map[string]any{"policy": a.Acceptance, "candidate_binary": originalBinary, "candidate_binary_snapshot": binary, "executed_argv": []string{"/bin/sh", script.Locator}},
			}
			if defect == "different candidate" {
				verification["candidate_oid"] = strings.Repeat("f", 40)
			}
			raw, err := json.Marshal(verification)
			if err != nil {
				t.Fatal(err)
			}
			receipt := write("candidate/verification-receipt", raw)
			tr.Artifacts = []TaskArtifactRecord{
				{ArtifactID: "candidate-review", Role: "review", Locator: review.Locator, SHA256: review.SHA256, SizeBytes: int64(len(reviewBytes))},
				{ArtifactID: "verification-receipt", Role: "other", Locator: receipt.Locator, SHA256: receipt.SHA256, SizeBytes: int64(len(raw))},
			}
			owner := write("owner", []byte("fixture native control evidence"))
			tr.HandoffLocator, tr.HandoffSHA256 = owner.Locator, owner.SHA256
			tr.StartReceiptLocator, tr.StartReceiptSHA256 = owner.Locator, owner.SHA256
			tr.TerminalResultLocator, tr.TerminalResultSHA256 = owner.Locator, owner.SHA256
			in := TaskCloseoutInput{Ownership: TaskCloseoutOwnership{EvidenceLocator: owner.Locator, EvidenceSHA256: owner.SHA256}}
			// The collector receives already validated native records. This small
			// fixture isolates its read-only graph traversal from effect execution.
			r := WorkItemRegistry{
				IntegrationAuthorities: []IntegrationAuthority{{ID: "authority", TaskID: tr.TaskID, TaskResultID: tr.ID, Plan: WorkspaceTaskIntegrationPlan{Task: IntegrationPlanTask{ResultOID: tr.ResultOID, ResultTree: tr.ResultTree}, DeliveryAuthorization: auth}}},
				IntegrationResults:     []IntegrationResult{{AuthorityID: "authority", Outcome: "exact_effect", RecoveryStatus: "complete"}},
			}
			write(filepath.Join(MarkerDirectory, WorkItemsFile), []byte("fixture native registry bytes"))
			switch defect {
			case "before integration":
				r.IntegrationResults = nil
			case "historical agreement":
				auth.Agreement.SchemaVersion, auth.Agreement.Acceptance = 1, nil
				auth.AgreementSHA256 = DeliveryAgreementDigest(auth.Agreement)
			case "tampered binary snapshot":
				write("attempt/candidate-executable", []byte("changed binary"))
			case "tampered review artifact":
				write("candidate/candidate-review", []byte("changed review"))
			}
			entries, err := collectCloseoutRetention(Dependencies{}, root, r, tr, in)
			if defect != "" {
				if err == nil {
					t.Fatal("unbound or unavailable snapshot satisfied retention")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			want := map[string][]byte{originalScript.Locator: scriptBytes, originalBinary.Locator: binaryBytes, originalReview.Locator: reviewBytes, receipt.Locator: raw}
			for _, entry := range entries {
				if expected, ok := want[entry.origin]; ok {
					if string(entry.bytes) != string(expected) || entry.sha != digestTaskBytes(expected) {
						t.Fatalf("retention rewrote original evidence: %s", entry.origin)
					}
					delete(want, entry.origin)
				}
			}
			if len(want) != 0 {
				t.Fatalf("retention lost original references: %v", want)
			}
		})
	}
}
