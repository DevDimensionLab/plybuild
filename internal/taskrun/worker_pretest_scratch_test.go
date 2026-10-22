//go:build darwin

package taskrun

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/sirupsen/logrus"
)

// Only the native reservation runs here: /usr/bin/true stands in for the
// dedicated runner. No provider, process observer or sandbox is invoked.
func reservePretestScratch(t *testing.T) (deliveryFixture, workflowState, workerAttemptIntent) {
	t.Helper()
	f, o := newWorkerDeliveryFixture(t)
	if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
		s.Result.Worker.Stage = "pretest"
		s.Result.Worker.Job = nil
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	s.Request.Runtime.PlyExecutable.Path = "/usr/bin/true"
	t.Setenv("TMPDIR", f.R.WorkspaceRoot)
	if err = workerPrepareAttempt(f.D, s); err != nil {
		t.Fatal(err)
	}
	s, err = workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	in, err := workerIntent(s, *s.Result.Worker.Job)
	if err != nil {
		t.Fatal(err)
	}
	return f, s, in
}

func TestPretestScratchReservedWorkspaceFixture(t *testing.T) {
	_, s, in := reservePretestScratch(t)
	var tmp string
	for _, e := range in.Environment {
		if v, ok := strings.CutPrefix(e, "TMPDIR="); ok {
			tmp = v
		}
	}
	t.Cleanup(func() {
		if !insidePlyWorkspace(tmp) {
			_ = os.RemoveAll(tmp)
		}
	})
	child := filepath.Join(tmp, "ordinary-workspace")
	if err := os.Mkdir(child, 0700); err != nil {
		t.Fatal(err)
	}
	deps := workspace.SystemDependencies()
	old, _ := os.Getwd()
	defer os.Chdir(old)
	if err := os.Chdir(child); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.Init(deps); err != nil {
		t.Fatalf("reserved pretest TMPDIR cannot host an ordinary workspace: %v", err)
	}
	if insidePlyWorkspace(tmp) || in.RunID != s.Result.RunID || in.JobID != s.Result.Worker.Job.ID || in.CandidateOID != s.Result.Worker.CandidateOID || in.CandidateTree != s.Result.Worker.CandidateTree {
		t.Fatal("scratch or candidate binding differs")
	}
	if tmp != in.PretestScratch {
		t.Fatal("intent does not bind TMPDIR")
	}
	physical, err := filepath.EvalSymlinks(tmp)
	if info, e := os.Stat(tmp); e != nil || info.Mode().Perm() != 0700 || err != nil || physical != tmp {
		t.Fatal("scratch is not private and physical")
	}
	policy, err := os.ReadFile(in.Argv[2])
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{tmp, filepath.Join(s.Result.Worker.Job.Directory, "home")} {
		if !strings.Contains(string(policy), "(allow file-write* (subpath "+strconv.Quote(path)+"))") {
			t.Fatalf("missing bounded policy: %s", policy)
		}
	}
	if strings.Count(string(policy), "(allow file-write*") != 3 || strings.Contains(string(policy), "(subpath "+strconv.Quote(filepath.Dir(tmp))+")") || in.Argv[0] != "/usr/bin/sandbox-exec" {
		t.Fatalf("policy changed scope: %s", policy)
	}
}

// Completion and publication inputs are modeled; no process execution is claimed.
func TestPretestScratchCompletionAndLegacy(t *testing.T) {
	for _, mode := range []string{"pass", "fail", "incomplete", "unknown", "publication-failure", "legacy"} {
		t.Run(mode, func(t *testing.T) {
			base := t.TempDir()
			own, other := filepath.Join(base, "own"), filepath.Join(base, "other")
			for _, p := range []string{own, other} {
				if err := os.Mkdir(p, 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(p, "child.log"), []byte("evidence"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			in := workerAttemptIntent{Kind: "pretest", RunID: "run", JobID: "job", PretestScratch: own, Environment: []string{"TMPDIR=" + own}}
			r := workerAttemptResult{PretestComplete: true, Exit: ptr(0), Process: &Process{State: "exited", ExitCode: ptr(0), Quiescence: ptr(true)}}
			durable := true
			switch mode {
			case "fail":
				r.Exit = ptr(3)
				r.Process.ExitCode = ptr(3)
			case "incomplete":
				r.PretestComplete = false
			case "unknown":
				r.Process = nil
			case "publication-failure":
				durable = false
			case "legacy":
				in.PretestScratch = ""
			}
			var warnings bytes.Buffer
			logger := logrus.StandardLogger()
			previous := logger.Out
			logger.SetOutput(&warnings)
			defer logger.SetOutput(previous)
			workerFinishPretestScratch(in, r, durable)
			if _, err := os.Stat(own); mode == "pass" {
				if !os.IsNotExist(err) {
					t.Fatal("successful scratch remains")
				}
			} else {
				if raw, err := os.ReadFile(filepath.Join(own, "child.log")); err != nil || string(raw) != "evidence" {
					t.Fatal("child evidence lost")
				}
				if mode != "legacy" && (!strings.Contains(warnings.String(), own) || !strings.Contains(warnings.String(), "warning")) {
					t.Fatal("retained locator warning missing")
				}
			}
			if _, err := os.Stat(filepath.Join(other, "child.log")); err != nil {
				t.Fatal("foreign evidence changed")
			}
		})
	}
}

func TestPretestScratchReservedRecoveryAndCorrection(t *testing.T) {
	f, s, in := reservePretestScratch(t)
	defer os.RemoveAll(in.PretestScratch)
	evidence := filepath.Join(in.PretestScratch, "child.log")
	if err := os.WriteFile(evidence, []byte("retained"), 0600); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if _, err := workerRunReserved(f.D, s); err == nil || !strings.Contains(err.Error(), "will not be replayed") {
			t.Fatalf("unknown reservation replayable: %v", err)
		}
	}
	fresh, err := workflowRead(f.R.WorkspaceRoot, s.Result.RunID)
	if err != nil || *fresh.Result.Worker.Job != *s.Result.Worker.Job {
		t.Fatal("recovery replaced reservation")
	}
	stdout, err := workflowKeep(f.D, filepath.Join(s.Result.Worker.Job.Directory, "stdout.txt"), "modeled output")
	if err != nil {
		t.Fatal(err)
	}
	stderr, err := workflowKeep(f.D, filepath.Join(s.Result.Worker.Job.Directory, "stderr.txt"), "modeled error")
	if err != nil {
		t.Fatal(err)
	}
	r := workerAttemptResult{RunID: in.RunID, JobID: in.JobID, IntentSHA256: s.Result.Worker.Job.IntentHash, FinishedAt: "2026-10-22T00:00:00Z", PretestComplete: true, Exit: ptr(3), Process: &Process{State: "exited", ExitCode: ptr(3), Quiescence: ptr(true)}, Stdout: stdout, Stderr: stderr}
	if err = f.D.writeValue(filepath.Join(s.Result.Worker.Job.Directory, "result.json"), r); err != nil {
		t.Fatal(err)
	}
	if err = workerStep(f.D, s); err != nil {
		t.Fatal(err)
	}
	fresh, err = workflowRead(f.R.WorkspaceRoot, s.Result.RunID)
	if err != nil || fresh.Result.Worker.Stage != "developer" || fresh.Result.Worker.Iteration != 2 {
		t.Fatal("failed pretest did not correct")
	}
	if raw, err := os.ReadFile(evidence); err != nil || string(raw) != "retained" {
		t.Fatal("correction deleted child evidence")
	}
}

// Faults exercise actual preparation and native reservation publication. The
// fixture's runner is never reached and no sandbox/provider is started.
func TestPretestScratchPreparationFailures(t *testing.T) {
	for _, mode := range []string{"allocation", "script", "reservation", "intent-publication"} {
		t.Run(mode, func(t *testing.T) {
			f, o := newWorkerDeliveryFixture(t)
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			s.Result.Worker.Stage = "pretest"
			s.Result.Worker.Job = nil
			base, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			autoScratchBases(t, base)
			if mode == "allocation" {
				autoScratchBases(t, f.R.WorkspaceRoot)
			}
			dir := filepath.Join(workerDirectory(s), "attempts", "pretest-0001")
			if mode == "script" {
				if err = os.MkdirAll(dir, 0700); err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(filepath.Join(dir, "acceptance.sh"), []byte("conflicting"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "reservation" {
				f.D.Fault = func(point string) error {
					if point == "worker_before_attempt_reservation" {
						return errors.New("modeled reservation failure")
					}
					return nil
				}
			}
			if mode == "intent-publication" {
				f.D.Files = pretestIntentPublicationFailure{}
			}
			var warnings bytes.Buffer
			logger := logrus.StandardLogger()
			previous := logger.Out
			logger.SetOutput(&warnings)
			defer logger.SetOutput(previous)
			if err = workerPrepareAttempt(f.D, s); err == nil {
				t.Fatal("preparation unexpectedly succeeded")
			}
			entries, err := os.ReadDir(base)
			if err != nil {
				t.Fatal(err)
			}
			if mode == "intent-publication" {
				if len(entries) != 1 {
					t.Fatalf("published intent scratch lost: %v", entries)
				}
				retained := filepath.Join(base, entries[0].Name())
				defer os.RemoveAll(retained)
				if !strings.Contains(warnings.String(), retained) || !strings.Contains(warnings.String(), o.RunID) || !strings.Contains(warnings.String(), "pretest-0001") {
					t.Fatal("retained publication warning missing identity")
				}
			} else if len(entries) != 0 {
				t.Fatal("never-started preparation leaked scratch")
			}
		})
	}
}

type pretestIntentPublicationFailure struct{ systemFileSystem }

func (fs pretestIntentPublicationFailure) WriteOnce(path string, raw []byte) error {
	if err := fs.systemFileSystem.WriteOnce(path, raw); err != nil {
		return err
	}
	if filepath.Base(path) == "intent.json" {
		return errors.New("modeled publication response failure after durable intent")
	}
	return nil
}

func TestPretestScratchDurableResultPublication(t *testing.T) {
	for _, mode := range []string{"pass", "failed-result", "stdout", "stderr", "result"} {
		t.Run(mode, func(t *testing.T) {
			scratch := t.TempDir()
			if err := os.WriteFile(filepath.Join(scratch, "child.log"), []byte("evidence"), 0600); err != nil {
				t.Fatal(err)
			}
			dir := t.TempDir()
			if err := os.Chmod(dir, 0700); err != nil {
				t.Fatal(err)
			}
			in := workerAttemptIntent{Kind: "pretest", RunID: "run", JobID: "job", PretestScratch: scratch}
			r := workerAttemptResult{RunID: "run", JobID: "job", PretestComplete: true, Exit: ptr(0), Process: &Process{State: "exited", ExitCode: ptr(0), Quiescence: ptr(true)}}
			if mode == "failed-result" {
				r.Exit = ptr(7)
				r.Process.ExitCode = ptr(7)
			}
			d := Dependencies{Files: pretestResultPublicationFailure{fail: mode}}
			var warnings bytes.Buffer
			logger := logrus.StandardLogger()
			previous := logger.Out
			logger.SetOutput(&warnings)
			defer logger.SetOutput(previous)
			err := workerPublishAttemptResult(d, in, r, dir, []byte("actual stdout"), []byte("actual stderr"))
			if (mode == "pass" || mode == "failed-result") != (err == nil) {
				t.Fatalf("publication result: %v", err)
			}
			if mode == "pass" {
				if _, err = os.Stat(scratch); !os.IsNotExist(err) {
					t.Fatal("durable pass retained scratch")
				}
			} else {
				if _, err = os.Stat(filepath.Join(scratch, "child.log")); err != nil {
					t.Fatal("evidence deleted before durable success")
				}
				if !strings.Contains(warnings.String(), scratch) {
					t.Fatal("retained result locator not reported")
				}
			}
			if mode == "pass" || mode == "failed-result" {
				var kept workerAttemptResult
				if err = readValue(filepath.Join(dir, "result.json"), 1<<20, &kept); err != nil {
					t.Fatal(err)
				}
				for _, b := range []FileBinding{kept.Stdout, kept.Stderr} {
					if _, err = workflowBound(b, 1<<20); err != nil {
						t.Fatal("native streams not durable", err)
					}
				}
			}
		})
	}
}

type pretestResultPublicationFailure struct {
	systemFileSystem
	fail string
}

func (fs pretestResultPublicationFailure) WriteOnce(path string, raw []byte) error {
	if strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) == fs.fail {
		return errors.New("modeled native publication failure")
	}
	return fs.systemFileSystem.WriteOnce(path, raw)
}

// A separate, ordinary unsandboxed shell invokes the exact make-built CLI.
// This demonstrates workspace creation and evidence after shell exit, and
// establishes no claim about live OS sandbox enforcement or provider behavior.
func TestPretestScratchOrdinaryCLIWorkspace(t *testing.T) {
	binary := os.Getenv("PLY_CANDIDATE_BINARY")
	if binary == "" {
		t.Skip("requires exact make-built PLY_CANDIDATE_BINARY")
	}
	_, s, in := reservePretestScratch(t)
	defer os.RemoveAll(in.PretestScratch)
	cmd := exec.Command("/bin/sh", "-c", `set -eu
fixture=$(mktemp -d "$TMPDIR/workspace.XXXXXX")
cd "$fixture"
"$PLY_CANDIDATE_BINARY" workspace init
printf '%s\n' "$fixture" > "$TMPDIR/child-locator"
printf 'ordinary child evidence' > "$fixture/child.log"
`)
	cmd.Dir = s.Observed.Target.WorktreeLocator
	cmd.Env = append([]string{}, in.Environment...)
	for i, e := range cmd.Env {
		if strings.HasPrefix(e, "PLY_CANDIDATE_BINARY=") {
			cmd.Env[i] = "PLY_CANDIDATE_BINARY=" + binary
		}
	}
	if raw, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("ordinary CLI fixture failed: %v\n%s", err, raw)
	}
	workerFinishPretestScratch(in, workerAttemptResult{}, false)
	raw, err := os.ReadFile(filepath.Join(in.PretestScratch, "child-locator"))
	if err != nil {
		t.Fatal(err)
	}
	if evidence, err := os.ReadFile(filepath.Join(strings.TrimSpace(string(raw)), "child.log")); err != nil || string(evidence) != "ordinary child evidence" {
		t.Fatal("child evidence did not survive ordinary shell exit")
	}
}
