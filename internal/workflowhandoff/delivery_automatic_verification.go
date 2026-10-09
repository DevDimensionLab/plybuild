package workflowhandoff

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// AutomaticCandidateBinaryLimit fits the native closeout evidence object bound.
// Reject larger input before any test or integration can start.
const AutomaticCandidateBinaryLimit = 64 << 20

// AutomaticVerificationInstructions describes the actual invocation. It is
// preserved before starting the command and is never a human attestation.
type AutomaticVerificationInstructions struct {
	Policy                  workspace.DeliveryAcceptancePolicy `json:"policy"`
	CandidateBinary         AutomaticBinaryBinding             `json:"candidate_binary"`
	CandidateBinarySnapshot AutomaticBinaryBinding             `json:"candidate_binary_snapshot"`
	ExecutedArgv            []string                           `json:"executed_argv"`
	Environment             []string                           `json:"environment"`
	ExecutorClaim           string                             `json:"executor_claim"`
	ExecutorSessionID       string                             `json:"executor_session_id"`
	ExpectedExit            int                                `json:"expected_exit"`
	InstructionsSHA256      string                             `json:"instructions_sha256"`
}

type AutomaticBinaryBinding struct {
	Locator string `json:"locator"`
	SHA256  string `json:"sha256"`
}

func AutomaticInstructionsDigest(v AutomaticVerificationInstructions, argv []string, cwd, scriptSHA, oid, tree string) string {
	v.InstructionsSHA256 = ""
	raw, _ := canonicaljson.Marshal(bridgeValue(struct {
		Instructions AutomaticVerificationInstructions `json:"instructions"`
		Argv         []string                          `json:"argv"`
		CWD          string                            `json:"cwd"`
		ScriptSHA256 string                            `json:"script_sha256"`
		OID          string                            `json:"candidate_oid"`
		Tree         string                            `json:"candidate_tree"`
	}{v, argv, cwd, scriptSHA, oid, tree}))
	return digestBytes(raw)
}

func ValidateAutomaticInstructions(v AutomaticVerificationInstructions, agreement *workspace.DeliveryAgreement, argv []string, cwd, scriptSHA, oid, tree string) error {
	if agreement == nil || !agreement.AutomaticAcceptance() || !canonicalEqual(bridgeValue(v.Policy), bridgeValue(agreement.Acceptance)) || v.ExpectedExit != 0 {
		return fmt.Errorf("automatic invocation differs from its frozen acceptance policy")
	}
	if !filepath.IsAbs(v.CandidateBinary.Locator) || filepath.Clean(v.CandidateBinary.Locator) != v.CandidateBinary.Locator || !validateDigest(v.CandidateBinary.SHA256) || validatePlainText("executor_claim", v.ExecutorClaim, 1, 256) != nil || validatePlainText("executor_session_id", v.ExecutorSessionID, 1, 256) != nil {
		return fmt.Errorf("automatic invocation lacks an exact candidate binary or actual executor identity")
	}
	if !filepath.IsAbs(v.CandidateBinarySnapshot.Locator) || filepath.Clean(v.CandidateBinarySnapshot.Locator) != v.CandidateBinarySnapshot.Locator || v.CandidateBinarySnapshot.SHA256 != v.CandidateBinary.SHA256 || len(v.ExecutedArgv) != 2 || v.ExecutedArgv[0] != "/bin/sh" || !filepath.IsAbs(v.ExecutedArgv[1]) {
		return fmt.Errorf("automatic execution requires the preserved script and candidate binary snapshots")
	}
	if len(v.Environment) == 0 || !slices.IsSorted(v.Environment) {
		return fmt.Errorf("automatic test environment must be explicit and sorted")
	}
	seen := map[string]bool{}
	for _, entry := range v.Environment {
		key, _, ok := strings.Cut(entry, "=")
		if !ok || seen[key] || strings.ContainsRune(entry, 0) {
			return fmt.Errorf("invalid automatic test environment")
		}
		switch key {
		case "HOME", "PATH", "TMPDIR", "LC_ALL", "PLY_CANDIDATE_BINARY", "PLY_CANDIDATE_OID", "PLY_CANDIDATE_TREE":
		default:
			return fmt.Errorf("automatic test environment contains an undeclared key %s", key)
		}
		seen[key] = true
	}
	for _, entry := range []string{"PLY_CANDIDATE_BINARY=" + v.CandidateBinarySnapshot.Locator, "PLY_CANDIDATE_OID=" + oid, "PLY_CANDIDATE_TREE=" + tree} {
		if !slices.Contains(v.Environment, entry) {
			return fmt.Errorf("automatic environment differs from candidate binding")
		}
	}
	if v.InstructionsSHA256 != AutomaticInstructionsDigest(v, argv, cwd, scriptSHA, oid, tree) {
		return fmt.Errorf("automatic instruction identity differs")
	}
	return nil
}
