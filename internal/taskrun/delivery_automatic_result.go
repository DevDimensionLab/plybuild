package taskrun

import (
	"encoding/base64"
	"os"
	"path/filepath"
)

// The completion spool is written atomically before output artifacts and native
// candidate registration. Recovery materializes this same result, never a test.
type automaticExecutionResult struct {
	Kind         string                      `json:"kind"`
	Receipt      deliveryVerificationReceipt `json:"receipt"`
	StdoutBase64 string                      `json:"stdout_base64"`
	StderrBase64 string                      `json:"stderr_base64"`
}

func preserveAutomaticResult(d Dependencies, path string, receipt deliveryVerificationReceipt, stdout, stderr []byte) error {
	if receipt.Automatic == nil {
		return nil
	}
	return d.writeValue(filepath.Join(path, "execution-result.json"), automaticExecutionResult{Kind: "PlyAutomaticExecution@1", Receipt: receipt, StdoutBase64: base64.StdEncoding.EncodeToString(stdout), StderrBase64: base64.StdEncoding.EncodeToString(stderr)})
}

func recoverAutomaticResult(d Dependencies, s workflowState) error {
	a := s.Result.Delivery.Attempt
	if a == nil || a.Kind != "verification" || s.Request.Delivery.Agreement == nil || !s.Request.Delivery.Agreement.AutomaticAcceptance() {
		return nil
	}
	if _, err := readFile(filepath.Join(a.Path, "verification.json"), 256<<10, true); err == nil || !os.IsNotExist(err) {
		return err
	}
	raw, err := readFile(filepath.Join(a.Path, "execution-result.json"), 16<<20, true)
	if os.IsNotExist(err) {
		return nil
	} // No durable completion: ordinary unknown guard remains authoritative.
	if err != nil {
		return err
	}
	var result automaticExecutionResult
	if err = decode(raw, 16<<20, &result); err != nil {
		return err
	}
	stdout, err := base64.StdEncoding.Strict().DecodeString(result.StdoutBase64)
	if err != nil {
		return err
	}
	stderr, err := base64.StdEncoding.Strict().DecodeString(result.StderrBase64)
	if err != nil {
		return err
	}
	initialRaw, err := readFile(filepath.Join(a.Path, "attempt.json"), 256<<10, true)
	if err != nil || hash(initialRaw) != a.InputSHA256 {
		return workflowError(4, "automatic completion differs from the reserved attempt")
	}
	var initial deliveryVerificationReceipt
	if err = decode(initialRaw, 256<<10, &initial); err != nil {
		return err
	}
	r := result.Receipt
	copy := r
	copy.Exit, copy.Stdout, copy.Stderr, copy.FinishedAt, copy.Error, copy.Outcome = nil, FileBinding{}, FileBinding{}, "", "", ""
	if result.Kind != "PlyAutomaticExecution@1" || !equal(initial, copy) || !deliveryVerificationVersion(r) || r.Automatic == nil || r.RunID != s.Result.RunID || r.RequestSHA256 != s.Result.RequestSHA256 || r.AttemptID != a.ID || a.Path != filepath.Join(s.Result.Paths.RunRoot, "delivery", "attempts", a.ID) || r.Stdout != (FileBinding{Locator: filepath.Join(a.Path, "stdout.txt"), SHA256: hash(stdout)}) || r.Stderr != (FileBinding{Locator: filepath.Join(a.Path, "stderr.txt"), SHA256: hash(stderr)}) {
		return workflowError(4, "preserved automatic completion is inconsistent; no re-execution is allowed")
	}
	if err = d.writeOnce(r.Stdout.Locator, stdout); err == nil {
		err = d.writeOnce(r.Stderr.Locator, stderr)
	}
	if err == nil {
		err = d.writeValue(filepath.Join(a.Path, "verification.json"), r)
	}
	return err
}
