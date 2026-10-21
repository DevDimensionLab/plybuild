package clirelease

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
)

type Receipt struct {
	Kind           string   `json:"kind"`
	SchemaVersion  int      `json:"schema_version"`
	Destination    string   `json:"destination"`
	Manifest       Manifest `json:"release"`
	ManifestSHA256 string   `json:"manifest_sha256"`
	BeforeSHA256   string   `json:"before_sha256"`
	Backup         string   `json:"backup"`
	State          string   `json:"state"`
}

// Installer owns only the destination binary and sibling receipt/backup files.
// Faults are per-instance test dependencies, never production environment flags.
type Installer struct{ fault func(string) error }

func (i Installer) check(stage string) error {
	if i.fault != nil {
		return i.fault(stage)
	}
	return nil
}

func syncDirectory(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return f.Sync()
}

func (i Installer) replace(path string, b []byte, mode fs.FileMode) (err error) {
	f, err := os.CreateTemp(filepath.Dir(path), ".ply-install-*.tmp")
	if err != nil {
		return err
	}
	defer func() { _ = f.Close(); _ = os.Remove(f.Name()) }()
	if err = f.Chmod(mode); err != nil {
		return err
	}
	if _, err = f.Write(b); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	if err = os.Rename(f.Name(), path); err != nil {
		return err
	}
	return syncDirectory(filepath.Dir(path))
}

func destinationPaths(destination string) (string, string, error) {
	if !filepath.IsAbs(destination) || filepath.Clean(destination) != destination {
		return "", "", fmt.Errorf("install_invalid_destination: use an absolute clean path")
	}
	parent := filepath.Dir(destination)
	physical, err := filepath.EvalSymlinks(parent)
	if err != nil || physical != parent {
		return "", "", fmt.Errorf("install_invalid_destination: parent must be an existing physical directory")
	}
	if st, err := os.Lstat(destination); err == nil && !st.Mode().IsRegular() {
		return "", "", fmt.Errorf("install_invalid_destination: existing destination must be a regular file")
	} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return "", "", err
	}
	return destination + ".install.json", destination + ".backups", nil
}

func installedDigest(path string) (string, []byte, error) {
	b, err := ReadRegular(path, 128<<20)
	if errors.Is(err, fs.ErrNotExist) {
		return "absent", nil, nil
	}
	if err != nil {
		return "", nil, err
	}
	return Digest(b), b, nil
}

func saveReceipt(i Installer, path string, r Receipt) error {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return i.replace(path, append(b, '\n'), 0o600)
}

func readReceipt(path string) (Receipt, error) {
	var r Receipt
	b, err := ReadRegular(path, 64<<10)
	if err != nil {
		return r, err
	}
	err = DecodeStrict(b, &r)
	if err == nil && (r.Kind != "PlyInstallReceipt@1" || r.SchemaVersion != 1 || (r.State != "pending" && r.State != "installed" && r.State != "rolled_back")) {
		err = fmt.Errorf("install_receipt_invalid: unrecognized receipt")
	}
	return r, err
}

func (i Installer) Install(directory, expected, destination string) (Receipt, error) {
	m, b, err := Load(directory, expected)
	if err != nil {
		return Receipt{}, err
	}
	return i.install(m, b, expected, destination)
}

func (i Installer) install(m Manifest, b []byte, expected, destination string) (Receipt, error) {
	var r Receipt
	receiptPath, backups, err := destinationPaths(destination)
	if err != nil {
		return r, err
	}
	unlock, err := installationLock(destination + ".install.lock")
	if err != nil {
		return r, err
	}
	defer unlock()
	if m.Qualification != nil {
		archived := filepath.Join(destination+".install-history", expected[len("sha256:"):]+".json")
		if _, err := readReceipt(archived); err == nil {
			return r, fmt.Errorf("install_superseded: this operation is historical; observe the later installation")
		} else if !errors.Is(err, fs.ErrNotExist) {
			return r, err
		}
	}
	current, old, err := installedDigest(destination)
	if err != nil {
		return r, err
	}
	previous, err := readReceipt(receiptPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return r, err
	}
	if err == nil && previous.State == "rolled_back" && previous.ManifestSHA256 == expected && m.Qualification != nil {
		return previous, fmt.Errorf("install_rolled_back: this operation was rolled back; select a new qualified operation")
	}
	if err == nil && previous.State == "pending" && (previous.Destination != destination || previous.ManifestSHA256 != expected) {
		return previous, fmt.Errorf("install_pending: recover the recorded package before another installation")
	}
	if err == nil && previous.ManifestSHA256 == expected && previous.Destination == destination && previous.State != "rolled_back" {
		r = previous
		if m.Qualification != nil && r.State == "installed" && current != m.SHA256 {
			return r, fmt.Errorf("install_destination_changed: completed installation no longer present; observe without retry")
		}
		if !reflect.DeepEqual(r.Manifest, m) || r.BeforeSHA256 != "absent" && !digestPattern.MatchString(r.BeforeSHA256) {
			return r, fmt.Errorf("install_receipt_invalid: release or preimage differs")
		}
		if r.BeforeSHA256 != "absent" {
			if r.Backup != filepath.Join(backups, r.BeforeSHA256[len("sha256:"):]) {
				return r, fmt.Errorf("install_backup_invalid: unexpected backup path")
			}
			backup, err := ReadRegular(r.Backup, 128<<20)
			if err != nil || Digest(backup) != r.BeforeSHA256 {
				return r, fmt.Errorf("install_backup_invalid: preserved preimage differs")
			}
		}
		if current != r.BeforeSHA256 && current != m.SHA256 {
			return r, fmt.Errorf("install_destination_changed: neither recorded preimage nor candidate is installed")
		}
	} else {
		if err == nil && previous.ManifestSHA256 != expected && m.Qualification != nil {
			if !digestPattern.MatchString(previous.ManifestSHA256) {
				return r, fmt.Errorf("install_receipt_invalid: invalid operation digest")
			}
			history := destination + ".install-history"
			if err := os.Mkdir(history, 0700); err != nil && !errors.Is(err, fs.ErrExist) {
				return r, err
			}
			st, err := os.Lstat(history)
			if err != nil || !st.IsDir() || st.Mode()&fs.ModeSymlink != 0 {
				return r, fmt.Errorf("install_history_invalid")
			}
			archive := filepath.Join(history, previous.ManifestSHA256[len("sha256:"):]+".json")
			if prior, err := readReceipt(archive); errors.Is(err, fs.ErrNotExist) {
				if err := saveReceipt(i, archive, previous); err != nil {
					return r, err
				}
			} else if err != nil || !reflect.DeepEqual(prior, previous) {
				return r, fmt.Errorf("install_history_invalid: previous receipt differs")
			}
		}
		r = Receipt{Kind: "PlyInstallReceipt@1", SchemaVersion: 1, Destination: destination, Manifest: m, ManifestSHA256: expected, BeforeSHA256: current, State: "pending"}
		if err := i.check("before-backup"); err != nil {
			return r, err
		}
		if current != "absent" {
			if err := os.Mkdir(backups, 0o700); err != nil && !errors.Is(err, fs.ErrExist) {
				return r, err
			}
			st, err := os.Lstat(backups)
			if err != nil || !st.IsDir() || st.Mode()&fs.ModeSymlink != 0 {
				return r, fmt.Errorf("install_backup_invalid: backup directory must be physical")
			}
			r.Backup = filepath.Join(backups, current[len("sha256:"):])
			if saved, err := ReadRegular(r.Backup, 128<<20); errors.Is(err, fs.ErrNotExist) {
				if err = i.replace(r.Backup, old, 0o500); err != nil {
					return r, err
				}
			} else if err != nil || Digest(saved) != current {
				return r, fmt.Errorf("install_backup_invalid: preserved preimage differs")
			}
		}
		if err := i.check("after-backup"); err != nil {
			return r, err
		}
		if err := saveReceipt(i, receiptPath, r); err != nil {
			return r, err
		}
	}
	if err := i.check("after-intent"); err != nil {
		return r, err
	}
	modeOK := false
	if st, err := os.Stat(destination); err == nil {
		modeOK = st.Mode().Perm() == 0o755
	}
	if current != m.SHA256 || !modeOK {
		// Stage in the destination directory; the old file stays intact until rename.
		f, err := os.CreateTemp(filepath.Dir(destination), ".ply-candidate-*.tmp")
		if err != nil {
			return r, err
		}
		defer func() { _ = f.Close(); _ = os.Remove(f.Name()) }()
		if err = f.Chmod(0o755); err != nil {
			return r, err
		}
		if _, err = f.Write(b); err != nil {
			return r, err
		}
		if err = f.Sync(); err != nil {
			return r, err
		}
		if err = f.Close(); err != nil {
			return r, err
		}
		if err = i.check("staged"); err != nil {
			return r, err
		}
		staged, err := ReadRegular(f.Name(), 128<<20)
		if err != nil || Digest(staged) != m.SHA256 {
			return r, fmt.Errorf("install_staging_failed: staged digest differs")
		}
		fresh, _, err := installedDigest(destination)
		if err != nil || fresh != current {
			return r, fmt.Errorf("install_destination_changed: recheck the installed preimage")
		}
		if err = i.check("before-replace"); err != nil {
			return r, err
		}
		if err = os.Rename(f.Name(), destination); err != nil {
			return r, err
		}
	}
	if err = i.check("after-replace"); err != nil {
		return r, fmt.Errorf("install_outcome_unknown: inspect the receipt and retry the same package: %w", err)
	}
	if err = syncDirectory(filepath.Dir(destination)); err != nil {
		return r, err
	}
	actual, _, err := installedDigest(destination)
	if err != nil || actual != m.SHA256 {
		return r, fmt.Errorf("install_verification_failed: installed digest differs")
	}
	st, err := os.Stat(destination)
	if err != nil || st.Mode().Perm() != 0o755 {
		return r, fmt.Errorf("install_verification_failed: installed executable permissions differ")
	}
	r.State = "installed"
	if err = saveReceipt(i, receiptPath, r); err != nil {
		r.State = "pending"
		return r, err
	}
	return r, nil
}

// Rollback restores only the exact preserved binary. It never downgrades data.
func (i Installer) Rollback(destination, expected string) (Receipt, error) {
	return i.rollback(destination, expected, "")
}

func (i Installer) rollback(destination, expected, operation string) (Receipt, error) {
	var r Receipt
	receiptPath, backups, err := destinationPaths(destination)
	if err != nil {
		return r, err
	}
	unlock, err := installationLock(destination + ".install.lock")
	if err != nil {
		return r, err
	}
	defer unlock()
	r, err = readReceipt(receiptPath)
	if err != nil {
		return r, err
	}
	if operation != "" && r.ManifestSHA256 != operation {
		return r, fmt.Errorf("rollback_binding_changed: later installation has a different identity")
	}
	if !digestPattern.MatchString(expected) || expected != r.Manifest.SHA256 || r.Destination != destination || !digestPattern.MatchString(r.BeforeSHA256) || r.Backup != filepath.Join(backups, r.BeforeSHA256[len("sha256:"):]) {
		return r, fmt.Errorf("rollback_binding_invalid: exact installed digest and preserved preimage are required")
	}
	current, _, err := installedDigest(destination)
	if err != nil {
		return r, err
	}
	if current != expected && current != r.BeforeSHA256 {
		return r, fmt.Errorf("install_destination_changed: rollback refuses unrelated bytes")
	}
	b, err := ReadRegular(r.Backup, 128<<20)
	if err != nil || Digest(b) != r.BeforeSHA256 {
		return r, fmt.Errorf("install_backup_invalid: preserved preimage differs")
	}
	if current != r.BeforeSHA256 {
		if err = i.replace(destination, b, 0o755); err != nil {
			return r, err
		}
	}
	r.State = "rolled_back"
	return r, saveReceipt(i, receiptPath, r)
}
