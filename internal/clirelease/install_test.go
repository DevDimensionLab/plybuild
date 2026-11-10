package clirelease

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func installFixture(t *testing.T) (string, Manifest, []byte) {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(dir, "ply")
	if err := os.WriteFile(dest, []byte("known old executable"), 0o755); err != nil {
		t.Fatal(err)
	}
	b := bytes.Repeat([]byte("qualified new bytes\n"), 100)
	m := Manifest{SHA256: Digest(b), SizeBytes: int64(len(b))}
	return dest, m, b
}

func TestInstallAtomicFailureAndRecovery(t *testing.T) {
	for _, stage := range []string{"after-intent", "staged", "before-replace", "after-replace"} {
		t.Run(stage, func(t *testing.T) {
			dest, m, b := installFixture(t)
			old, _ := os.ReadFile(dest)
			i := Installer{fault: func(point string) error {
				if point == stage {
					return errors.New("injected interruption")
				}
				return nil
			}}
			r, err := i.install(m, b, Digest([]byte("manifest")), dest)
			if err == nil || r.State == "installed" {
				t.Fatalf("false success: %+v %v", r, err)
			}
			got, _ := os.ReadFile(dest)
			if !bytes.Equal(got, old) && !bytes.Equal(got, b) {
				t.Fatal("partial executable")
			}
			if stage != "after-replace" && !bytes.Equal(got, old) {
				t.Fatal("changed before commit point")
			}
			backup, err := os.ReadFile(r.Backup)
			if err != nil || !bytes.Equal(backup, old) {
				t.Fatal("lost identifiable preimage")
			}
			r, err = (Installer{}).install(m, b, Digest([]byte("manifest")), dest)
			if err != nil || r.State != "installed" {
				t.Fatalf("recovery: %+v %v", r, err)
			}
			r, err = (Installer{}).Rollback(dest, m.SHA256)
			got, _ = os.ReadFile(dest)
			if err != nil || r.State != "rolled_back" || !bytes.Equal(got, old) {
				t.Fatalf("rollback: %+v %v", r, err)
			}
		})
	}
}

func TestInstallRepairsMatchingBytesWithoutExecutablePermission(t *testing.T) {
	dest, m, b := installFixture(t)
	if err := os.WriteFile(dest, b, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dest, 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := (Installer{}).install(m, b, Digest([]byte("manifest")), dest)
	if err != nil {
		t.Fatal(err)
	}
	st, err := os.Stat(dest)
	if err != nil || st.Mode().Perm()&0o111 != 0o111 || r.State != "installed" {
		t.Fatalf("installer reported success without executable permission: %+v %v %v", r, st.Mode(), err)
	}
}

func TestInstallRefusesChangedDestinationAndBackup(t *testing.T) {
	for _, change := range []string{"destination", "backup"} {
		t.Run(change, func(t *testing.T) {
			dest, m, b := installFixture(t)
			i := Installer{fault: func(s string) error {
				if s == "after-intent" {
					return errors.New("stop")
				}
				return nil
			}}
			r, _ := i.install(m, b, Digest([]byte("manifest")), dest)
			target := dest
			if change == "backup" {
				target = r.Backup
				_ = os.Chmod(target, 0o600)
			}
			if err := os.WriteFile(target, []byte("unrelated bytes"), 0o600); err != nil {
				t.Fatal(err)
			}
			before, _ := os.ReadFile(dest)
			if _, err := (Installer{}).install(m, b, r.ManifestSHA256, dest); err == nil {
				t.Fatal("accepted drift")
			}
			after, _ := os.ReadFile(dest)
			if !bytes.Equal(before, after) {
				t.Fatal("rejection changed destination")
			}
		})
	}
}

func TestInstallRejectsInvalidDigestPlatformAndManifest(t *testing.T) {
	for _, raw := range []string{`{"kind":"a","kind":"b"}`, `{"unknown":true}`, `{} {}`} {
		var m Manifest
		if DecodeStrict([]byte(raw), &m) == nil {
			t.Fatalf("accepted %s", raw)
		}
	}
	if err := ValidateExecutable("unused", []byte("x"), Manifest{OS: "windows", Arch: "arm64"}); err == nil || !strings.Contains(err.Error(), "platform_mismatch") {
		t.Fatal(err)
	}
	if err := ValidateExecutable("unused", []byte("x"), Manifest{OS: "darwin", Arch: "arm64", SHA256: Digest([]byte("other"))}); err == nil {
		t.Fatal("wrong digest accepted")
	}
	if err := ValidateExecutable("unused", []byte("x"), Manifest{OS: "darwin", Arch: "arm64", SHA256: Digest([]byte("x")), SizeBytes: 1}); err == nil {
		t.Fatal("non-Mach-O accepted")
	}
}

// A real process dies at each side of rename. No defer/recovery handler runs.
func TestInstallProcessKillRecovery(t *testing.T) {
	if os.Getenv("PLY_INSTALL_TEST_CHILD") == "1" {
		dest := os.Getenv("PLY_INSTALL_TEST_DEST")
		stage := os.Getenv("PLY_INSTALL_TEST_STAGE")
		b := []byte("new executable")
		i := Installer{fault: func(s string) error {
			if s == stage {
				_ = os.WriteFile(dest+".checkpoint", []byte(s), 0o600)
				select {}
			}
			return nil
		}}
		_, _ = i.install(Manifest{SHA256: Digest(b), SizeBytes: int64(len(b))}, b, Digest([]byte("manifest")), dest)
		os.Exit(3)
	}
	for _, stage := range []string{"staged", "after-replace"} {
		t.Run(stage, func(t *testing.T) {
			dest, _, _ := installFixture(t)
			old, _ := os.ReadFile(dest)
			cmd := exec.Command(os.Args[0], "-test.run=^TestInstallProcessKillRecovery$")
			cmd.Env = append(os.Environ(), "PLY_INSTALL_TEST_CHILD=1", "PLY_INSTALL_TEST_DEST="+dest, "PLY_INSTALL_TEST_STAGE="+stage)
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			defer func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }()
			deadline := time.Now().Add(10 * time.Second)
			for {
				if _, err := os.Stat(dest + ".checkpoint"); err == nil {
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("child missed checkpoint")
				}
				time.Sleep(10 * time.Millisecond)
			}
			if err := cmd.Process.Kill(); err != nil {
				t.Fatal(err)
			}
			_ = cmd.Wait()
			b := []byte("new executable")
			got, _ := os.ReadFile(dest)
			if !bytes.Equal(got, old) && !bytes.Equal(got, b) {
				t.Fatal("partial bytes after SIGKILL")
			}
			r, err := (Installer{}).install(Manifest{SHA256: Digest(b), SizeBytes: int64(len(b))}, b, Digest([]byte("manifest")), dest)
			if err != nil || r.State != "installed" {
				t.Fatalf("recovery after SIGKILL: %+v %v", r, err)
			}
		})
	}
}
