package clirelease

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
)

// QualifiedCLI is supplied by the native delivery reader, never a profile or
// owner attestation. It keeps provenance inside the install/backup receipt.
type QualifiedCLI struct {
	RunID              string `json:"run_id"`
	Project            string `json:"project"`
	Repo               string `json:"repo"`
	Epic               string `json:"epic"`
	OID                string `json:"oid"`
	Tree               string `json:"tree"`
	BinarySHA256       string `json:"binary_sha256"`
	VerificationSHA256 string `json:"verification_sha256"`
}

func qualifiedManifest(q QualifiedCLI, destination string) (Manifest, string, error) {
	if q.Project != "ply" || q.Repo != "ply" || q.Epic == "" || q.RunID == "" || !commitPattern.MatchString(q.OID) || !commitPattern.MatchString(q.Tree) || !digestPattern.MatchString(q.BinarySHA256) || !digestPattern.MatchString(q.VerificationSHA256) {
		return Manifest{}, "", fmt.Errorf("qualified_cli_invalid: native Ply qualification is required")
	}
	raw, _ := json.Marshal(struct {
		Source      QualifiedCLI
		Destination string
	}{q, destination})
	return Manifest{Kind: "PlyQualifiedCLI@1", SchemaVersion: 1, SourceCommit: q.OID, OS: runtime.GOOS, Arch: runtime.GOARCH, SHA256: q.BinarySHA256, Qualification: &q}, Digest(raw), nil
}

// InstallQualified reuses the release installer's durable preimage, atomic
// replacement and same-operation reconciliation. It does not rebuild bytes.
func (i Installer) InstallQualified(q QualifiedCLI, binary []byte, destination string) (Receipt, error) {
	m, expected, err := qualifiedManifest(q, destination)
	if err != nil {
		return Receipt{}, err
	}
	if Digest(binary) != q.BinarySHA256 {
		return Receipt{}, fmt.Errorf("qualified_cli_changed: retained bytes differ")
	}
	m.SizeBytes = int64(len(binary))
	return i.install(m, binary, expected, destination)
}

func (i Installer) RollbackQualified(q QualifiedCLI, destination string) (Receipt, error) {
	_, expected, err := qualifiedManifest(q, destination)
	if err != nil {
		return Receipt{}, err
	}
	return i.rollback(destination, q.BinarySHA256, expected)
}

// ObserveQualified never copies bytes. A pending receipt and installed hash
// identify an interrupted replace; calling InstallQualified with identical
// inputs reconciles that operation. Unrelated hashes remain unknown.
func (i Installer) ObserveQualified(q QualifiedCLI, destination string) (Receipt, string, error) {
	_, expected, err := qualifiedManifest(q, destination)
	if err != nil {
		return Receipt{}, "", err
	}
	path, backups, err := destinationPaths(destination)
	if err != nil {
		return Receipt{}, "", err
	}
	unlock, err := installationLock(destination + ".install.lock")
	if err != nil {
		return Receipt{}, "", err
	}
	defer unlock()
	r, err := readReceipt(path)
	if err != nil {
		return r, "", err
	}
	if r.ManifestSHA256 != expected || r.Destination != destination {
		return r, "", fmt.Errorf("install_binding_changed: receipt belongs to another operation")
	}
	if r.BeforeSHA256 != "absent" {
		if !digestPattern.MatchString(r.BeforeSHA256) || r.Backup != filepath.Join(backups, r.BeforeSHA256[len("sha256:"):]) {
			return r, "", fmt.Errorf("install_backup_invalid")
		}
		backup, err := ReadRegular(r.Backup, 128<<20)
		if err != nil || Digest(backup) != r.BeforeSHA256 {
			return r, "", fmt.Errorf("install_backup_invalid: preserved preimage differs")
		}
	}
	actual, _, err := installedDigest(destination)
	return r, actual, err
}
