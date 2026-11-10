package clirelease

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestQualifiedInstallRecovery(t *testing.T) {
	for _, stage := range []string{"before-backup", "after-backup", "after-intent", "before-replace", "after-replace"} {
		t.Run(stage, func(t *testing.T) {
			dest := filepath.Join(t.TempDir(), "ply")
			old, candidate := []byte("old"), []byte("qualified")
			if err := os.WriteFile(dest, old, 0755); err != nil {
				t.Fatal(err)
			}
			q := QualifiedCLI{RunID: "run", Project: "ply", Repo: "ply", Epic: "epic", OID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Tree: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", BinarySHA256: Digest(candidate), VerificationSHA256: Digest([]byte("verification"))}
			installer := Installer{fault: func(s string) error {
				if s == stage {
					return errors.New("interrupted")
				}
				return nil
			}}
			if _, err := installer.InstallQualified(q, candidate, dest); err == nil {
				t.Fatal("fault did not fail")
			}
			replaced, _ := os.Stat(dest)
			r, err := (Installer{}).InstallQualified(q, candidate, dest)
			if err != nil || r.State != "installed" || r.BeforeSHA256 != Digest(old) {
				t.Fatalf("recovery: %+v %v", r, err)
			}
			after, _ := os.Stat(dest)
			if stage == "after-replace" && !os.SameFile(replaced, after) {
				t.Fatal("recovery copied already-installed bytes again")
			}
			backup, err := os.ReadFile(r.Backup)
			if err != nil || string(backup) != string(old) {
				t.Fatal("lost backup", err)
			}
			again, err := (Installer{}).InstallQualified(q, candidate, dest)
			if err != nil || again.Backup != r.Backup {
				t.Fatal("repetition replaced backup", err)
			}
			if err := os.WriteFile(dest, []byte("later result"), 0755); err != nil {
				t.Fatal(err)
			}
			if _, err := (Installer{}).RollbackQualified(q, dest); err == nil {
				t.Fatal("rollback overwrote unrelated result")
			}
		})
	}
}

func TestQualifiedInstallSupersededOperationCannotReplay(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "ply")
	candidate := []byte("candidate")
	q := QualifiedCLI{RunID: "first", Project: "ply", Repo: "ply", Epic: "epic", OID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Tree: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", BinarySHA256: Digest(candidate), VerificationSHA256: Digest([]byte("verification"))}
	if err := os.WriteFile(dest, []byte("old"), 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := (Installer{}).InstallQualified(q, candidate, dest); err != nil {
		t.Fatal(err)
	}
	later := q
	later.RunID = "later"
	if _, err := (Installer{}).InstallQualified(later, candidate, dest); err != nil {
		t.Fatal(err)
	}
	if _, err := (Installer{}).InstallQualified(q, candidate, dest); err == nil {
		t.Fatal("old operation replayed over later separately chosen installation")
	}
	if _, err := (Installer{}).RollbackQualified(q, dest); err == nil {
		t.Fatal("old rollback crossed operation identity")
	}
}

func TestQualifiedInstallUnknownAndChangedInputs(t *testing.T) {
	for _, change := range []string{"input", "destination", "backup"} {
		t.Run(change, func(t *testing.T) {
			dest := filepath.Join(t.TempDir(), "ply")
			candidate := []byte("candidate")
			q := QualifiedCLI{RunID: "first", Project: "ply", Repo: "ply", Epic: "epic", OID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Tree: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", BinarySHA256: Digest(candidate), VerificationSHA256: Digest([]byte("verification"))}
			if err := os.WriteFile(dest, []byte("old"), 0755); err != nil {
				t.Fatal(err)
			}
			r, err := (Installer{fault: func(s string) error {
				if s == "after-intent" {
					return errors.New("interrupted")
				}
				return nil
			}}).InstallQualified(q, candidate, dest)
			if err == nil {
				t.Fatal("fault missing")
			}
			switch change {
			case "input":
				q.VerificationSHA256 = Digest([]byte("other verification"))
			case "destination":
				if err := os.WriteFile(dest, []byte("unknown"), 0755); err != nil {
					t.Fatal(err)
				}
			case "backup":
				if err := os.Chmod(r.Backup, 0600); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(r.Backup, []byte("corrupted backup"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := (Installer{}).InstallQualified(q, candidate, dest); err == nil {
				t.Fatal("changed input or unknown effect allowed blind retry")
			}
		})
	}
}

func TestQualifiedInstallCannotReplayAfterLegacyRelease(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "ply")
	candidate := []byte("qualified candidate")
	q := QualifiedCLI{RunID: "first", Project: "ply", Repo: "ply", Epic: "epic", OID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Tree: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", BinarySHA256: Digest(candidate), VerificationSHA256: Digest([]byte("verification"))}
	if err := os.WriteFile(dest, []byte("old"), 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := (Installer{}).InstallQualified(q, candidate, dest); err != nil {
		t.Fatal(err)
	}
	// This is the shared transaction called by the existing release-package entrypoint after Load validates the package.
	later := []byte("later deliberately selected release")
	m := Manifest{Kind: "PlyRelease@1", SchemaVersion: 1, SHA256: Digest(later), SizeBytes: int64(len(later))}
	if _, err := (Installer{}).install(m, later, Digest([]byte("later release manifest")), dest); err != nil {
		t.Fatal(err)
	}
	r, err := (Installer{}).InstallQualified(q, candidate, dest)
	actual, readErr := os.ReadFile(dest)
	t.Logf("old qualified retry err=%v state=%s installed=%q", err, r.State, actual)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if err == nil || string(actual) != string(later) {
		t.Fatalf("old qualified operation overwrote later separately selected release: err=%v installed=%q", err, actual)
	}

	// An explicit legacy rollback/reinstall remains legitimate. The qualified
	// tombstone must not turn a release package into a single-use operation.
	if _, err := (Installer{}).Rollback(dest, m.SHA256); err != nil {
		t.Fatal(err)
	}
	for attempt := 0; attempt < 2; attempt++ {
		if _, err := (Installer{}).install(m, later, Digest([]byte("later release manifest")), dest); err != nil {
			t.Fatalf("legitimate legacy reinstall %d: %v", attempt, err)
		}
	}
	actual, err = os.ReadFile(dest)
	if err != nil || string(actual) != string(later) {
		t.Fatalf("legacy reinstall lost selected bytes: %q %v", actual, err)
	}
	if _, err := (Installer{}).RollbackQualified(q, dest); err == nil {
		t.Fatal("old qualified rollback crossed later legacy selection")
	}
}
