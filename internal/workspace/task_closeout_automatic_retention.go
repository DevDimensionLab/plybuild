package workspace

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

type closeoutSnapshotBinding struct {
	Locator string `json:"locator"`
	SHA256  string `json:"sha256"`
}

// automaticCloseoutSnapshots resolves only the original test inputs of an
// already integrated automatic candidate. The native TaskResult binds the
// verification receipt bytes; that receipt binds both the original reference
// and an equal snapshot. Retention keeps the original reference and bytes.
// This resolver grants no execution authority and is never used by the live
// pre-integration acceptance gate.
func automaticCloseoutSnapshots(r WorkItemRegistry, tr TaskResultRecord) (map[string]closeoutEvidence, error) {
	var auth *DeliveryAuthorization
	for _, effect := range r.IntegrationResults {
		if effect.RecoveryStatus != "complete" || (effect.Outcome != "exact_effect" && effect.Outcome != "already_integrated") {
			continue
		}
		a := findAuthority(r, effect.AuthorityID)
		if a == nil || a.TaskID != tr.TaskID || a.TaskResultID != tr.ID || a.Plan.Task.ResultOID != tr.ResultOID || a.Plan.Task.ResultTree != tr.ResultTree || !automaticIntegrationAuthorization(a.Plan.DeliveryAuthorization) || a.Plan.DeliveryAuthorization.CandidateRunID != tr.RunID {
			continue
		}
		if auth != nil && !contentTypedEqual(auth, a.Plan.DeliveryAuthorization) {
			return nil, fmt.Errorf("automatic retention has conflicting integrated authority")
		}
		auth = a.Plan.DeliveryAuthorization
	}
	if auth == nil {
		return nil, nil
	}
	for _, artifact := range tr.Artifacts {
		if artifact.ArtifactID != "verification-receipt" {
			continue
		}
		raw, err := readCloseoutEvidence(artifact.Locator, artifact.SHA256)
		if err != nil {
			return nil, err
		}
		if int64(len(raw)) != artifact.SizeBytes {
			return nil, fmt.Errorf("automatic retention verification receipt size differs")
		}
		var receipt struct {
			Kind               string                  `json:"kind"`
			SchemaVersion      int                     `json:"schema_version"`
			RunID              string                  `json:"run_id"`
			RequestSHA256      string                  `json:"request_sha256"`
			CandidateOID       string                  `json:"candidate_oid"`
			CandidateTree      string                  `json:"candidate_tree"`
			CWD                string                  `json:"cwd"`
			Acceptance         closeoutSnapshotBinding `json:"acceptance"`
			AcceptanceSnapshot closeoutSnapshotBinding `json:"acceptance_snapshot"`
			Review             closeoutSnapshotBinding `json:"review"`
			Exit               *int                    `json:"exit"`
			Error              string                  `json:"error"`
			Outcome            string                  `json:"outcome"`
			Automatic          *struct {
				Policy                  DeliveryAcceptancePolicy `json:"policy"`
				CandidateBinary         closeoutSnapshotBinding  `json:"candidate_binary"`
				CandidateBinarySnapshot closeoutSnapshotBinding  `json:"candidate_binary_snapshot"`
				ExecutedArgv            []string                 `json:"executed_argv"`
			} `json:"automatic"`
		}
		if err := json.Unmarshal(raw, &receipt); err != nil || receipt.Kind != "PlyDeliveryVerification@2" || receipt.SchemaVersion != 2 || receipt.Automatic == nil || receipt.RunID != auth.RunID || receipt.RequestSHA256 != auth.RequestSHA256 || receipt.CandidateOID != tr.ResultOID || receipt.CandidateTree != tr.ResultTree || receipt.CWD != tr.SourceLocator || receipt.Outcome != "pass" || receipt.Exit == nil || *receipt.Exit != 0 || receipt.Error != "" || !contentTypedEqual(receipt.Automatic.Policy, *auth.Agreement.Acceptance) {
			return nil, fmt.Errorf("automatic retention receipt differs from the integrated candidate")
		}
		automatic := receipt.Automatic
		if len(automatic.ExecutedArgv) != 2 || automatic.ExecutedArgv[0] != "/bin/sh" || automatic.ExecutedArgv[1] != receipt.AcceptanceSnapshot.Locator || automatic.CandidateBinarySnapshot.Locator != filepath.Join(filepath.Dir(receipt.AcceptanceSnapshot.Locator), "candidate-executable") {
			return nil, fmt.Errorf("automatic retention snapshots differ from the executed invocation")
		}
		var review *TaskArtifactRecord
		for i := range tr.Artifacts {
			if tr.Artifacts[i].ArtifactID == "candidate-review" {
				review = &tr.Artifacts[i]
			}
		}
		if review == nil || review.SHA256 != receipt.Review.SHA256 {
			return nil, fmt.Errorf("automatic retention review differs from the native candidate artifact")
		}
		out := map[string]closeoutEvidence{}
		pairs := []struct {
			origin, snapshot closeoutSnapshotBinding
			size             int64
		}{
			{receipt.Acceptance, receipt.AcceptanceSnapshot, -1},
			{automatic.CandidateBinary, automatic.CandidateBinarySnapshot, -1},
			{receipt.Review, closeoutSnapshotBinding{Locator: review.Locator, SHA256: review.SHA256}, review.SizeBytes},
		}
		for _, pair := range pairs {
			origin, snapshot := pair.origin, pair.snapshot
			if !validAbsoluteCleanPath(origin.Locator) || !digestPattern.MatchString(origin.SHA256) || origin.SHA256 != snapshot.SHA256 {
				return nil, fmt.Errorf("automatic retention snapshot differs from its original binding")
			}
			raw, err := readCloseoutEvidence(snapshot.Locator, snapshot.SHA256)
			if err != nil {
				return nil, err
			}
			if pair.size >= 0 && int64(len(raw)) != pair.size {
				return nil, fmt.Errorf("automatic retention snapshot size differs")
			}
			if previous, ok := out[origin.Locator]; ok && previous.sha != origin.SHA256 {
				return nil, fmt.Errorf("automatic retention original has conflicting snapshots")
			}
			out[origin.Locator] = closeoutEvidence{origin: origin.Locator, sha: origin.SHA256, bytes: raw}
		}
		return out, nil
	}
	return nil, nil
}
