package integration

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func adapterPhysicalTemp(t *testing.T) string {
	t.Helper()
	p, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func installFixture(t *testing.T) (InstallProfile, []byte) {
	t.Helper()
	root := adapterPhysicalTemp(t)
	artifact := []byte("exact candidate build\n")
	h := sha256.Sum256(artifact)
	cp, err := exec.LookPath("cp")
	if err != nil {
		t.Fatal(err)
	}
	cmp, err := exec.LookPath("cmp")
	if err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(root, "candidate")
	if err = os.WriteFile(source, artifact, 0600); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "installed")
	return InstallProfile{Kind: "ply.integration.install-profile", SchemaVersion: 1, Name: "selected-local", CandidateOID: strings.Repeat("a", 40), ArtifactPath: target, ArtifactSHA256: "sha256:" + hex.EncodeToString(h[:]), Install: CommandSpec{Executable: cp, Arguments: []string{source, target}, CWD: root}, Verify: CommandSpec{Executable: cmp, Arguments: []string{source, target}, CWD: root}}, artifact
}

func TestLocalInstallRunsOnlySelectedCommandAndVerifiesArtifact(t *testing.T) {
	p, _ := installFixture(t)
	a := LocalInstallAdapter{}
	o, err := a.Install(p)
	if err != nil || o.State != "verified" || o.ArtifactSHA256 != p.ArtifactSHA256 || o.CandidateOID != p.CandidateOID || o.ProfileSHA256 != InstallProfileDigest(p) {
		t.Fatalf("install was not artifact-verified: %+v %v", o, err)
	}
	o, err = a.Observe(p)
	if err != nil || o.State != "verified" {
		t.Fatalf("readback cannot verify installed version: %+v %v", o, err)
	}
}

func TestInstallationNeverEvaluatesArgumentTextAsShellCode(t *testing.T) {
	p, artifact := installFixture(t)
	p.Install.Arguments = []string{"$(touch should-not-exist)", "`literal`", "line one\nline two"}
	var calls []CommandSpec
	a := LocalInstallAdapter{Run: func(cwd, program string, args ...string) ([]byte, error) {
		calls = append(calls, CommandSpec{Executable: program, Arguments: append([]string(nil), args...), CWD: cwd})
		if len(calls) == 1 {
			return nil, os.WriteFile(p.ArtifactPath, artifact, 0600)
		}
		return nil, nil
	}}
	o, err := a.Install(p)
	if err != nil || o.State != "verified" || !reflect.DeepEqual(calls, []CommandSpec{p.Install, p.Verify}) {
		t.Fatalf("profile argv changed: %+v %v calls=%v", o, err, calls)
	}
	if _, err = os.Stat(filepath.Join(p.Install.CWD, "should-not-exist")); !os.IsNotExist(err) {
		t.Fatal("argument text executed")
	}
}

func TestInstallationExitSuccessWithoutArtifactProofDoesNotQualify(t *testing.T) {
	p, _ := installFixture(t)
	if err := os.WriteFile(p.ArtifactPath, []byte("wrong build"), 0600); err != nil {
		t.Fatal(err)
	}
	calls := 0
	a := LocalInstallAdapter{Run: func(string, string, ...string) ([]byte, error) { calls++; return nil, nil }}
	o, err := a.Install(p)
	if err != nil || o.State != "install_failed" || o.Reason != "installed_artifact_mismatch" || calls != 1 {
		t.Fatalf("exit-only success qualified: %+v %v calls=%d", o, err, calls)
	}
	// Recovery cannot distinguish an unstarted install from a partial effect.
	o, err = a.Observe(p)
	if err != nil || o.State != "effect_unknown" || calls != 1 {
		t.Fatalf("unknown installation retried: %+v %v calls=%d", o, err, calls)
	}
}

func TestInstallationUnknownRunUsesObservationAndNeverBlindRetry(t *testing.T) {
	for _, installed := range []bool{false, true} {
		t.Run(map[bool]string{false: "unproved", true: "artifact-observed"}[installed], func(t *testing.T) {
			p, artifact := installFixture(t)
			installs := 0
			a := LocalInstallAdapter{Run: func(_ string, program string, _ ...string) ([]byte, error) {
				if program == p.Install.Executable {
					installs++
					if installed {
						if err := os.WriteFile(p.ArtifactPath, artifact, 0600); err != nil {
							return nil, err
						}
					}
					return nil, errors.New("connection to process interrupted")
				}
				return nil, nil
			}}
			o, _ := a.Install(p)
			want := "effect_unknown"
			if installed {
				want = "verified"
			}
			if o.State != want {
				t.Fatalf("unknown outcome falsely classified: %+v", o)
			}
			for i := 0; i < 2; i++ {
				o, _ = a.Observe(p)
				if o.State != want {
					t.Fatalf("readback: %+v", o)
				}
			}
			if installs != 1 {
				t.Fatalf("recovery repeated install %d times", installs)
			}
		})
	}
}

func TestInstallationVerificationFailureAndArtifactSwapRemainIncomplete(t *testing.T) {
	for _, mode := range []string{"verification-exit", "verification-swaps-artifact", "symlink-artifact"} {
		t.Run(mode, func(t *testing.T) {
			p, artifact := installFixture(t)
			if err := os.WriteFile(p.ArtifactPath, artifact, 0600); err != nil {
				t.Fatal(err)
			}
			if mode == "symlink-artifact" {
				if err := os.Remove(p.ArtifactPath); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(p.Install.Arguments[0], p.ArtifactPath); err != nil {
					t.Fatal(err)
				}
			}
			a := LocalInstallAdapter{Run: func(string, string, ...string) ([]byte, error) {
				if mode == "verification-exit" {
					return nil, exec.Command("false").Run()
				}
				if mode == "verification-swaps-artifact" {
					return nil, os.WriteFile(p.ArtifactPath, []byte("different version"), 0600)
				}
				t.Fatal("verification executed with unsafe artifact")
				return nil, nil
			}}
			o, err := a.Observe(p)
			if err == nil || o.State == "verified" {
				t.Fatalf("unverified artifact accepted: %+v %v", o, err)
			}
		})
	}
}

func TestInstallProfileStrictBindingChangesItsDigest(t *testing.T) {
	p, _ := installFixture(t)
	path := filepath.Join(p.Install.CWD, "profile.json")
	raw, _ := json.Marshal(p)
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadInstallProfile(path)
	if err != nil || InstallProfileDigest(loaded) != InstallProfileDigest(p) {
		t.Fatalf("profile read: %+v %v", loaded, err)
	}
	if err = ValidateInstallProfile(p, strings.Repeat("f", 40)); err == nil {
		t.Fatal("wrong candidate accepted")
	}
	for name, change := range map[string]func(*InstallProfile){
		"cwd":        func(p *InstallProfile) { p.Install.CWD += "-different" },
		"executable": func(p *InstallProfile) { p.Install.Executable += "-different" },
		"argument":   func(p *InstallProfile) { p.Install.Arguments = append([]string{"different"}, p.Install.Arguments...) },
		"verify":     func(p *InstallProfile) { p.Verify.Arguments = []string{"another"} },
		"artifact":   func(p *InstallProfile) { p.ArtifactSHA256 = "sha256:" + strings.Repeat("f", 64) },
	} {
		t.Run(name, func(t *testing.T) {
			changed := p
			change(&changed)
			if InstallProfileDigest(changed) == InstallProfileDigest(p) {
				t.Fatal("changed effect profile retained old plan digest")
			}
		})
	}
	for _, bad := range []string{strings.TrimSuffix(string(raw), "}") + `,"unexpected":true}`, strings.TrimSuffix(string(raw), "}") + `,"name":"duplicated"}`} {
		if err = os.WriteFile(path, []byte(bad), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err = LoadInstallProfile(path); err == nil {
			t.Fatal("ambiguous install profile accepted")
		}
	}
}

func TestInstallControllerMustSurviveReplacementAndCleanup(t *testing.T) {
	p, _ := installFixture(t)
	source := filepath.Join(p.Install.CWD, "task")
	if err := os.Mkdir(source, 0700); err != nil {
		t.Fatal(err)
	}
	controller := filepath.Join(p.Install.CWD, "frozen-control")
	input := filepath.Join(p.Install.CWD, "frozen-input")
	for _, path := range []string{controller, input, filepath.Join(source, "controller"), filepath.Join(source, "input"), p.ArtifactPath} {
		if err := os.WriteFile(path, []byte("preserved"), 0700); err != nil {
			t.Fatal(err)
		}
	}
	if err := ValidateInstallController(p, controller, source, input); err != nil {
		t.Fatal(err)
	}
	for _, bad := range []struct{ controller, input string }{{filepath.Join(source, "controller"), input}, {p.ArtifactPath, input}, {controller, filepath.Join(source, "input")}} {
		if err := ValidateInstallController(p, bad.controller, source, bad.input); err == nil {
			t.Fatalf("unsafe controller/input accepted: %+v", bad)
		}
	}
	if err := os.Remove(p.ArtifactPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(controller, p.ArtifactPath); err != nil {
		t.Fatal(err)
	}
	if err := ValidateInstallController(p, controller, source, input); err == nil {
		t.Fatal("controller hardlink can be overwritten by installation")
	}
}
