package workflowhandoff

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func humanQAValidationAttestation(result workspace.TaskResultRecord) DeliveryHumanAttestation {
	return DeliveryHumanAttestation{Kind: "DeliveryHumanAttestation@1", SchemaVersion: 1, TaskID: string(result.TaskID), TaskResultID: string(result.ID), ResultOID: result.ResultOID, ResultTree: result.ResultTree, Outcome: "pass", ActorClaim: "Synthetic fixture human; not actual product approval", StartSurface: "Isolated test", StartedAtUTC: "2026-10-09T01:00:00Z", CompletedAtUTC: "2026-10-09T01:01:00Z", Answer: "pass", Observation: "Synthetic answer checks the native provenance contract."}
}

func humanQAValidationDraft(result workspace.TaskResultRecord, a DeliveryHumanAttestation, evidence string, raw []byte, key string) map[string]any {
	return map[string]any{"kind": "WorkspaceTaskHumanQARecordDraft@1", "schema_version": 1, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "publication_key": key, "task_id": result.TaskID, "task_result_id": result.ID, "result_oid": result.ResultOID, "result_tree": result.ResultTree, "outcome": a.Outcome, "actor": workspace.TaskHumanActorRecord{ActorClaim: a.ActorClaim, StartSurface: a.StartSurface, StartedAtUTC: a.StartedAtUTC, CompletedAtUTC: a.CompletedAtUTC}, "evidence": []workspace.TaskHumanQAEvidenceRecord{{ID: "human-attestation", Role: "report", Locator: evidence, SHA256: digestBytes(raw), SizeBytes: int64(len(raw))}}, "observation": a.Observation, "accepted_residual_risks": []workspace.TaskResidualRiskRecord{}}
}

func TestDeliveryHumanAttestationMatchesWorkspaceProvenance(t *testing.T) {
	f := prepareDeliveryFixture(t)
	qualified, err := QualifyDeliveryCandidate(f.d, f.candidate(t, "qa-provenance"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		change func(*DeliveryHumanAttestation)
		valid  bool
	}{
		{"utc-z", func(*DeliveryHumanAttestation) {}, true},
		{"start-zero-offset", func(a *DeliveryHumanAttestation) { a.StartedAtUTC = "2026-10-09T01:00:00+00:00" }, false},
		{"completed-zero-offset", func(a *DeliveryHumanAttestation) { a.CompletedAtUTC = "2026-10-09T01:01:00+00:00" }, false},
		{"both-zero-offset", func(a *DeliveryHumanAttestation) {
			a.StartedAtUTC, a.CompletedAtUTC = "2026-10-09T01:00:00+00:00", "2026-10-09T01:01:00+00:00"
		}, false},
		{"negative-zero-offset", func(a *DeliveryHumanAttestation) { a.StartedAtUTC = "2026-10-09T01:00:00-00:00" }, false},
		{"nonzero-offset", func(a *DeliveryHumanAttestation) { a.StartedAtUTC = "2026-10-09T01:00:00+02:00" }, false},
		{"fractional-z", func(a *DeliveryHumanAttestation) {
			a.StartedAtUTC, a.CompletedAtUTC = "2026-10-09T01:00:00.123456789Z", "2026-10-09T01:01:00.000000001Z"
		}, true},
		{"reversed-time", func(a *DeliveryHumanAttestation) { a.CompletedAtUTC = "2026-10-09T00:59:59Z" }, false},
		{"unicode-actor-boundary", func(a *DeliveryHumanAttestation) { a.ActorClaim = strings.Repeat("ø", 256) }, true},
		{"actor-whitespace", func(a *DeliveryHumanAttestation) { a.ActorClaim += " " }, false},
		{"surface-control", func(a *DeliveryHumanAttestation) { a.StartSurface += "\u0085" }, false},
		{"observation-too-long", func(a *DeliveryHumanAttestation) { a.Observation = strings.Repeat("x", 2001) }, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			a := humanQAValidationAttestation(qualified.TaskResult)
			tc.change(&a)
			raw, err := json.Marshal(a)
			if err != nil {
				t.Fatal(err)
			}
			_, upperErr := ValidateDeliveryHumanAttestation(raw, "task", qualified.TaskResult, "pass")
			path := filepath.Join(f.root, "attestation-"+tc.name+".json")
			if err := os.WriteFile(path, raw, 0600); err != nil {
				t.Fatal(err)
			}
			draft := filepath.Join(f.root, "qa-draft-"+tc.name+".json")
			deliveryFixtureWrite(t, draft, humanQAValidationDraft(qualified.TaskResult, a, path, raw, "fixture/provenance/"+tc.name))
			draftRaw, err := os.ReadFile(draft)
			if err != nil {
				t.Fatal(err)
			}
			if err := workspace.ValidateTaskHumanQADraftBytes(draftRaw); (err == nil) != tc.valid {
				t.Errorf("read-only native schema differs from provenance contract: want_valid=%t err=%v", tc.valid, err)
			}
			before, err := f.w.WorkItems.Snapshot(f.root)
			if err != nil {
				t.Fatal(err)
			}
			_, lowerErr := workspace.RecordTaskHumanQA(f.w, workspace.TaskHumanQARecordInput{TaskID: "task", File: draft})
			if (upperErr == nil) != tc.valid || (lowerErr == nil) != tc.valid {
				t.Errorf("pre-validation/native publication disagree with provenance contract: want_valid=%t upper=%v lower=%v", tc.valid, upperErr, lowerErr)
			}
			if !tc.valid {
				after, err := f.w.WorkItems.Snapshot(f.root)
				if err != nil || after.RawSHA256 != before.RawSHA256 {
					t.Fatalf("rejected provenance changed native records: %v", err)
				}
			}
			if kept, err := os.ReadFile(path); err != nil || !bytes.Equal(kept, raw) {
				t.Fatalf("provenance validation rewrote human evidence: %v", err)
			}
		})
	}
}

func TestDeliveryHumanAttestationRejectsOffsetsBeforeStaging(t *testing.T) {
	f := prepareDeliveryFixture(t)
	qualified, err := QualifyDeliveryCandidate(f.d, f.candidate(t, "qa-no-staging"))
	if err != nil {
		t.Fatal(err)
	}
	a := humanQAValidationAttestation(qualified.TaskResult)
	a.CompletedAtUTC = "2026-10-09T01:01:00+00:00"
	raw, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(f.root, "invalid-offset-attestation.json")
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = RecordDeliveryHumanQA(f.d, "task", qualified.TaskResult, "pass", path); err == nil {
		t.Fatal("offset timestamp was published as human QA")
	}
	control, err := f.d.Store.ReadByLocator(qualified.Handoff.Locator)
	if err != nil {
		t.Fatal(err)
	}
	area := filepath.Join(control.Handoff.ReplyRoot, "staging", "human-qa-"+digestBytes(raw)[7:])
	if _, err = os.Lstat(area); !os.IsNotExist(err) {
		t.Fatalf("pre-validation allowed rejected provenance to reserve immutable QA staging: %v", err)
	}
}

func stageLegacyHumanQA(t *testing.T, f deliveryFixture, result workspace.TaskResultRecord, a DeliveryHumanAttestation, name string) (string, DeliveryRejectedHumanQA) {
	t.Helper()
	raw, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(f.root, "legacy-attestation-"+name+".json")
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	control, err := f.d.Store.ReadByLocator(result.HandoffLocator)
	if err != nil {
		t.Fatal(err)
	}
	area := filepath.Join(control.Handoff.ReplyRoot, "staging", "human-qa-"+digestBytes(raw)[7:])
	if err = os.MkdirAll(area, 0700); err != nil {
		t.Fatal(err)
	}
	proof := DeliveryRejectedHumanQA{AttestationPath: filepath.Join(area, "attestation.json"), DraftPath: filepath.Join(area, "qa-draft.json"), PublicationKey: "delivery-qa/" + string(result.ID) + "/" + digestBytes(raw)[7:]}
	if err = os.WriteFile(proof.AttestationPath, raw, 0600); err != nil {
		t.Fatal(err)
	}
	draft, err := canonicaljson.Marshal(deliveryObject(humanQAValidationDraft(result, a, proof.AttestationPath, raw, proof.PublicationKey)))
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(proof.DraftPath, draft, 0600); err != nil {
		t.Fatal(err)
	}
	return path, proof
}

func TestDeliveryHumanQATimestampRejectionProofPreservesExactEvidence(t *testing.T) {
	f := prepareDeliveryFixture(t)
	qualified, err := QualifyDeliveryCandidate(f.d, f.candidate(t, "qa-rejection-proof"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, start, completed string }{
		{"start-offset", "2026-10-09T01:00:00+00:00", "2026-10-09T01:01:00Z"},
		{"completed-offset", "2026-10-09T01:00:00Z", "2026-10-09T01:01:00+00:00"},
		{"both-offset", "2026-10-09T01:00:00+00:00", "2026-10-09T01:01:00+00:00"},
		{"fractional-offset", "2026-10-09T01:00:00.123+00:00", "2026-10-09T01:01:00.456+00:00"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			a := humanQAValidationAttestation(qualified.TaskResult)
			a.StartedAtUTC, a.CompletedAtUTC = tc.start, tc.completed
			path, expected := stageLegacyHumanQA(t, f, qualified.TaskResult, a, tc.name)
			before := map[string][]byte{}
			for _, artifact := range []string{path, expected.AttestationPath, expected.DraftPath} {
				before[artifact], err = os.ReadFile(artifact)
				if err != nil {
					t.Fatal(err)
				}
			}
			registry, err := f.w.WorkItems.Snapshot(f.root)
			if err != nil {
				t.Fatal(err)
			}
			if _, err = workspace.RecordTaskHumanQA(f.w, workspace.TaskHumanQARecordInput{TaskID: "task", File: expected.DraftPath}); err == nil {
				t.Fatal("legacy fixture was not rejected by actual native publication")
			}
			proof, err := ValidateRejectedDeliveryHumanQA(f.d, "task", qualified.TaskResult, "pass", path)
			if err != nil || proof != expected {
				t.Fatalf("known staged UTC rejection was not proven: %+v %v", proof, err)
			}
			for artifact, original := range before {
				if kept, err := os.ReadFile(artifact); err != nil || !bytes.Equal(kept, original) {
					t.Fatalf("rejection proof rewrote %s: %v", artifact, err)
				}
			}
			entries, err := os.ReadDir(filepath.Dir(expected.DraftPath))
			if err != nil || len(entries) != 2 {
				t.Fatalf("rejection proof created another evidence artifact: %v %v", entries, err)
			}
			after, err := f.w.WorkItems.Snapshot(f.root)
			if err != nil || after.RawSHA256 != registry.RawSHA256 || len(after.HumanQARecords) != 0 {
				t.Fatalf("rejection proof published a human answer: %+v %v", after.HumanQARecords, err)
			}
		})
	}
}

func TestDeliveryHumanQATimestampRejectionProofRejectsOtherCauses(t *testing.T) {
	f := prepareDeliveryFixture(t)
	qualified, err := QualifyDeliveryCandidate(f.d, f.candidate(t, "qa-unknown-rejection"))
	if err != nil {
		t.Fatal(err)
	}
	for _, mutation := range []string{"already-z", "nonzero-offset", "negative-zero-offset", "invalid-time", "reversed-time", "actor-whitespace", "observation-control", "different-answer", "changed-attestation", "missing-attestation", "symlink-attestation", "changed-draft", "missing-draft", "noncanonical-draft", "invalid-native-result-id"} {
		t.Run(mutation, func(t *testing.T) {
			result := qualified.TaskResult
			if mutation == "invalid-native-result-id" {
				result.ID = "not-a-native-result-id"
			}
			a := humanQAValidationAttestation(result)
			a.CompletedAtUTC = "2026-10-09T01:01:00+00:00"
			a.Observation += " " + mutation
			switch mutation {
			case "already-z":
				a.CompletedAtUTC = "2026-10-09T01:01:00Z"
			case "nonzero-offset":
				a.CompletedAtUTC = "2026-10-09T02:01:00+01:00"
			case "negative-zero-offset":
				a.CompletedAtUTC = "2026-10-09T01:01:00-00:00"
			case "invalid-time":
				a.CompletedAtUTC = "not-a-timestamp+00:00"
			case "reversed-time":
				a.CompletedAtUTC = "2026-10-09T00:59:59+00:00"
			case "actor-whitespace":
				a.ActorClaim += " "
			case "observation-control":
				a.Observation += "\n"
			case "different-answer":
				a.Answer = "fail"
			}
			path, staged := stageLegacyHumanQA(t, f, result, a, mutation)
			switch mutation {
			case "changed-attestation":
				if err := os.WriteFile(staged.AttestationPath, []byte("{}"), 0600); err != nil {
					t.Fatal(err)
				}
			case "missing-attestation", "symlink-attestation":
				if err := os.Remove(staged.AttestationPath); err != nil {
					t.Fatal(err)
				}
				if mutation == "symlink-attestation" {
					if err := os.Symlink(path, staged.AttestationPath); err != nil {
						t.Fatal(err)
					}
				}
			case "changed-draft":
				if err := os.WriteFile(staged.DraftPath, []byte("{}"), 0600); err != nil {
					t.Fatal(err)
				}
			case "missing-draft":
				if err := os.Remove(staged.DraftPath); err != nil {
					t.Fatal(err)
				}
			case "noncanonical-draft":
				original, err := os.ReadFile(staged.DraftPath)
				if err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(staged.DraftPath, append(original, '\n'), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if proof, err := ValidateRejectedDeliveryHumanQA(f.d, "task", result, "pass", path); err == nil || proof != (DeliveryRejectedHumanQA{}) {
				t.Fatalf("unknown rejection was classified as the known UTC mismatch: %+v %v", proof, err)
			}
		})
	}
}
