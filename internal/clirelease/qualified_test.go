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
