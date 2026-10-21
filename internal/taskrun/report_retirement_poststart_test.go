package taskrun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Raw CLI streams/exits are optional durable evidence, separate from modeled
// provider observations. Every invocation gets its own directory, including retry.
func al14Command(t *testing.T, c *exec.Cmd) ([]byte, []byte, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run()
	exit := 0
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			exit = e.ExitCode()
		} else {
			t.Fatal(err)
		}
	}
	if dir := os.Getenv("PLY_AL14_EVIDENCE_DIR"); dir != "" {
		runDir, e := os.MkdirTemp(dir, strings.ReplaceAll(t.Name(), "/", "_")+"-")
		if e != nil {
			t.Fatal(e)
		}
		argv, _ := json.Marshal(c.Args)
		for name, raw := range map[string][]byte{"stdout": stdout.Bytes(), "stderr": stderr.Bytes(), "exit": []byte(fmt.Sprintln(exit)), "argv.json": argv} {
			if e := os.WriteFile(filepath.Join(runDir, name), raw, 0600); e != nil {
				t.Fatal(e)
			}
		}
	}
	return stdout.Bytes(), stderr.Bytes(), err
}
func al14Combined(t *testing.T, c *exec.Cmd) ([]byte, error) {
	stdout, stderr, err := al14Command(t, c)
	return append(stdout, stderr...), err
}
func al14CLIOK(t *testing.T, binary string, f workflowFixture, args ...string) WorkflowRun {
	t.Helper()
	c := exec.Command(binary, append([]string{"workflow"}, append(args, "--format", "json")...)...)
	c.Dir, _ = f.D.CWD()
	stdout, stderr, err := al14Command(t, c)
	if err != nil {
		t.Fatalf("CLI %v: %v %s %s", args, err, stdout, stderr)
	}
	var out WorkflowRun
	if err := json.Unmarshal(stdout, &out); err != nil {
		t.Fatalf("JSON stdout: %v %s", err, stdout)
	}
	return out
}

// Herdr and provider are modeled; callbacks and retirement use the actual CLI
// against disposable native workspace/Git stores. No real historical Run is used.
func al14Started(t *testing.T, binary string, legacy bool) (workflowFixture, WorkflowRun) {
	t.Helper()
	f := workflowTestFixture(t)
	b := mustReadAL01(t, f.R.Herdr.Executable.Path)
	b = bytes.Replace(b, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+retireOpenPaneStandin+closedTransportInspectStandin+integrationOwnerRecoveryStandin), 1)
	if err := os.WriteFile(f.R.Herdr.Executable.Path, b, 0700); err != nil {
		t.Fatal(err)
	}
	f.R.Herdr.Executable.SHA256 = hash(b)
	if legacy {
		f.R.Runtime.PlyExecutable = Executable{binary, hashFileTest(t, binary)}
		f.File = writeAny(t, f.R.WorkspaceRoot, "al14-legacy.json", f.R)
	} else {
		f = al01ReceiptFixtureFrom(t, f, binary)
	}
	o := workflowTestStart(t, f)
	f.Acceptance.RunID, f.Acceptance.RequestSHA256, f.Acceptance.SessionID = o.RunID, o.RequestSHA256, o.SessionID
	f.Report.RunID, f.Report.RequestSHA256, f.Report.SessionID = o.RunID, o.RequestSHA256, o.SessionID
	if legacy {
		o = al14CLIOK(t, binary, f, "run", "accept", o.RunID, "--context", o.Paths.Context, "--file", writeAny(t, f.R.WorkspaceRoot, "legacy-accept.json", f.Acceptance))
	} else {
		o = al14CLIOK(t, binary, f, "run", "check-in", o.RunID, "--context", o.Paths.Context, "--provider", "codex", "--native-session", o.Transport.AgentSessionID)
	}
	return f, o
}

func al14Stop(t *testing.T, binary string, f workflowFixture, o WorkflowRun, effect string) WorkflowRun {
	t.Helper()
	proof := writeAny(t, f.R.WorkspaceRoot, "al14-denial.json", map[string]any{"modeled_provider": true, "action": "report transport", "effect_status": effect})
	f.D.CWD = func() (string, error) { return f.R.WorkspaceRoot, nil }
	return al14CLIOK(t, binary, f, "execute", "observe-stop", o.RunID, "--event-id", "al14-stop", "--provider", "codex", "--native-session", o.Transport.AgentSessionID, "--stage", "report", "--observer", "modeled coordinator", "--reason", "Provider stopped after actual local work", "--action", "report transport", "--effect-status", effect, "--evidence", proof)
}
func al14Retire(t *testing.T, binary string, f workflowFixture, o WorkflowRun) WorkflowRun {
	t.Helper()
	f.D.CWD = func() (string, error) { return f.R.WorkspaceRoot, nil }
	return al14CLIOK(t, binary, f, "run", "retire", o.RunID, "--observer", "modeled coordinator", "--reason", "Retain stopped report source")
}

func TestAL14CLIPostStartRetainsWorkAndHistory(t *testing.T) {
	binary := automaticCLIBinary(t)
	old := continuitySourceBinary(t, "2a2fa2c0070d4084ffc000caaaa318dfe6a9ab4d", "PLY_AL14_OLD_BINARY")
	for _, legacy := range []bool{false, true} {
		t.Run(fmt.Sprint("legacy=", legacy), func(t *testing.T) {
			f, o := al14Started(t, binary, legacy)
			o = al14CLIOK(t, binary, f, "run", "report", o.RunID, "--context", o.Paths.Context, "--file", writeAny(t, f.R.WorkspaceRoot, "prior-report.json", f.Report))
			source := mustReportWorktree(t, f, o)
			work := filepath.Join(source, "unique-work")
			if err := os.WriteFile(work, []byte("actual retained work\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if out, err := exec.Command("git", "-C", source, "add", "unique-work").CombinedOutput(); err != nil {
				t.Fatalf("%v %s", err, out)
			}
			if out, err := exec.Command("git", "-C", source, "commit", "-m", "Disposable unique report work").CombinedOutput(); err != nil {
				t.Fatalf("%v %s", err, out)
			}
			if err := os.WriteFile(work, []byte("actual retained work\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(source, "partial-draft"), []byte("retained dirty draft"), 0600); err != nil {
				t.Fatal(err)
			}
			o = al14Stop(t, binary, f, o, DeliveryEffectNotStarted)
			s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
			preserved := map[string]string{}
			for _, p := range []string{s.Acceptance.Locator, s.StartDraft.Locator, s.ReportProviderStop.Locator, o.Paths.Context, s.Result.Handoff.Locator, s.Records[0].Report.Locator, s.Records[0].Terminal.Locator} {
				preserved[p] = hashFileTest(t, p)
			}
			if !legacy {
				stateHash := hashFileTest(t, filepath.Join(o.Paths.RunRoot, "state.json"))
				c := exec.Command(old, "workflow", "run", "retire", o.RunID, "--observer", "modeled coordinator", "--reason", "Retain stopped report source")
				c.Dir = f.R.WorkspaceRoot
				out, err := al14Combined(t, c)
				if err == nil || !strings.Contains(string(out), "native start/terminal exists") {
					t.Fatalf("historical refusal missing: %v %s", err, out)
				}
				if hashFileTest(t, filepath.Join(o.Paths.RunRoot, "state.json")) != stateHash {
					t.Fatal("old refusal changed history")
				}
				t.Logf("historical post-start refusal (modeled provider, actual CLI): %s", out)
			}
			// Freeze launcher bytes as history: upgrading the consumed provider must not
			// force a new owner callback or launch for coordinator retirement.
			if err := os.WriteFile(f.R.Runtime.Executable.Path, []byte("changed historical provider"), 0700); err != nil {
				t.Fatal(err)
			}
			workflowTestModel(t, f, map[string]any{"closed_transport": true})
			retired := al14Retire(t, binary, f, o)
			if retired.FinalReturn.State != "retired" || retired.FinalReturn.TerminalSHA256 != nil || retired.TaskResultState != "not_published" || retired.Review.Decision != nil || retired.Delivery != nil {
				t.Fatalf("false delivery: %+v", retired)
			}
			r, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
			record, err := workflowReportRetirement(r)
			if err != nil {
				t.Fatal(err)
			}
			if record.ReportSource.OID != gitOutput(t, source, "rev-parse", "HEAD") || record.ReportSource.OID == r.Observed.Target.OID {
				t.Fatal("unique commit was lost or misreported")
			}
			var retained struct {
				ReportSource *struct {
					Clean        bool   `json:"clean"`
					SourceSHA256 string `json:"source_sha256"`
				} `json:"report_source"`
			}
			json.Unmarshal(mustReadAL01(t, r.ReportRetirement.Locator), &retained)
			if retained.ReportSource == nil || retained.ReportSource.Clean || retained.ReportSource.SourceSHA256 == "" {
				t.Fatal("retained dirty source was not recorded truthfully")
			}
			for p, h := range preserved {
				if hashFileTest(t, p) != h {
					t.Fatalf("history changed: %s", p)
				}
			}
			if string(mustReadAL01(t, work)) != "actual retained work\n" {
				t.Fatal("work lost")
			}
			digest := hashFileTest(t, r.ReportRetirement.Locator)
			al14Retire(t, binary, f, o) // identical lost reply reads the immutable retirement
			if hashFileTest(t, r.ReportRetirement.Locator) != digest {
				t.Fatal("retirement replayed")
			}
		})
	}
}

func TestAL14PostStartGuardsAndRecovery(t *testing.T) {
	binary := automaticCLIBinary(t)
	for _, which := range []string{"live", "inventory_unknown", "writer", "wrong_ref", "wrong_session", "unknown", "partial", "source_after_resolution", "lost_reply", "accepted"} {
		t.Run(which, func(t *testing.T) {
			f, o := al14Started(t, binary, false)
			source := mustReportWorktree(t, f, o)
			effect := DeliveryEffectNotStarted
			if which == "unknown" || which == "source_after_resolution" {
				effect = DeliveryEffectUnknown
			}
			if which == "partial" {
				effect = DeliveryEffectPartial
			}
			if which == "accepted" {
				o = workflowTestSettle(t, f, workflowTestReport(t, f, o))
				decision := workflowTestReview(t, f, o, "accepted")
				var err error
				o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, decision)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				o = al14Stop(t, binary, f, o, effect)
			}
			workflowTestModel(t, f, map[string]any{"closed_transport": true})
			f.D.CWD = func() (string, error) { return f.R.WorkspaceRoot, nil }
			f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
				return []LocalProcess{{PID: os.Getpid(), CWD: f.R.WorkspaceRoot}}, nil
			}
			switch which {
			case "wrong_session":
				proof := writeAny(t, f.R.WorkspaceRoot, "wrong-session-proof.json", map[string]any{"modeled": true})
				c := exec.Command(binary, "workflow", "execute", "observe-stop", o.RunID, "--event-id", "wrong-session", "--provider", "codex", "--native-session", "wrong-owner", "--stage", "report", "--observer", "modeled coordinator", "--reason", "Wrong identity guard", "--action", "report transport", "--effect-status", DeliveryEffectUnknown, "--evidence", proof)
				c.Dir = f.R.WorkspaceRoot
				out, err := al14Combined(t, c)
				if err == nil || !strings.Contains(string(out), "provider/session differs") {
					t.Fatalf("wrong identity accepted: %v %s", err, out)
				}
				return
			case "live":
				workflowTestModel(t, f, map[string]any{"closed_transport": false})
			case "inventory_unknown":
				workflowTestModel(t, f, map[string]any{"retire_inventory_denied": true})
			case "writer":
				f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
					return []LocalProcess{{PID: os.Getpid(), CWD: f.R.WorkspaceRoot}, {PID: 999999, CWD: source}}, nil
				}
			case "wrong_ref":
				if out, err := exec.Command("git", "-C", source, "checkout", "-b", "wrong-source").CombinedOutput(); err != nil {
					t.Fatalf("%v %s", err, out)
				}
			case "source_after_resolution":
				s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
				sha, _ := workflowTarget(f.D, s)
				proof := writeAny(t, f.R.WorkspaceRoot, "effect-result.json", map[string]any{"modeled_action_result": "not started"})
				_, err := WorkflowResolveReportStop(f.D, f.R.WorkspaceRoot, o.RunID, ReportStopResolutionInput{StopSHA256: s.ReportProviderStop.SHA256, SourceSHA256: sha, EffectStatus: DeliveryEffectNotStarted, Basis: "action_result", Observer: "modeled coordinator", Reason: "Observed callback rejection", EvidencePath: proof})
				if err != nil {
					t.Fatal(err)
				}
				os.WriteFile(filepath.Join(source, "later"), []byte("unreconciled"), 0600)
			case "lost_reply":
				f.D.Fault = func(at string) error {
					if at == "report_after_retirement" {
						return fmt.Errorf("modeled lost reply")
					}
					return nil
				}
			}
			in := DeliveryRetireInput{Observer: "modeled coordinator", Reason: "Retain stopped report source"}
			var err error
			var cliOutput []byte
			if which == "writer" || which == "lost_reply" {
				_, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, in)
			} else {
				c := exec.Command(binary, "workflow", "run", "retire", o.RunID, "--observer", in.Observer, "--reason", in.Reason)
				c.Dir = f.R.WorkspaceRoot
				out, e := al14Combined(t, c)
				cliOutput = out
				err = e
				if e != nil {
					t.Logf("actual CLI guard: %s", out)
				}
			}
			if which == "unknown" || which == "partial" || which == "source_after_resolution" {
				if err != nil {
					t.Fatal(err)
				}
				s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
				r, err := workflowReportRetirement(s)
				if err != nil {
					t.Fatal(err)
				}
				if bytes.Count(cliOutput, []byte("no replay or qualification")) != 1 || !bytes.Contains(cliOutput, []byte("original_effect="+effect)) || !bytes.Contains(cliOutput, []byte("retained_record=")) {
					t.Fatalf("missing sourced warning: %s", cliOutput)
				}
				c := exec.Command(binary, "workflow", "run", "retire", o.RunID, "--observer", in.Observer, "--reason", in.Reason, "--format", "json")
				c.Dir = f.R.WorkspaceRoot
				stdout, stderr, err := al14Command(t, c)
				if err != nil {
					t.Fatalf("retry: %v %s", err, stderr)
				}
				if bytes.Contains(stderr, []byte("no replay or qualification")) {
					t.Fatal("retry repeated warning")
				}
				var readback WorkflowRun
				if err := json.Unmarshal(stdout, &readback); err != nil {
					t.Fatalf("JSON stdout contaminated: %v", err)
				}
				if len(r.ReportEffects) != 1 || r.ReportEffects[0].OriginalEffect != effect {
					t.Fatal("uncertainty erased")
				}
				if s.Result.FinalReturn.TerminalSHA256 != nil || s.Result.TaskResultState != "not_published" {
					t.Fatal("uncertainty qualified")
				}
				return
			}
			if err == nil {
				t.Fatal("guard did not block retirement")
			}
			if which == "lost_reply" {
				f.D.Fault = nil
				if _, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil {
					t.Fatal(err)
				}
				return
			}
			s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if s.ReportRetirement != nil {
				t.Fatal("guard recorded retirement")
			}
			p, err := workspace.ShowTaskPreparation(f.D.Workspace, s.Request.PreparationID)
			if err != nil || p.Disposition != "current" {
				t.Fatalf("guard released queue: %s %v", p.Disposition, err)
			}
		})
	}
}

func TestAL14CLIReconcilesClosedPreparation(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := al14Started(t, binary, false)
	o = al14Stop(t, binary, f, o, DeliveryEffectUnknown)
	s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
	sha, _ := workflowTarget(f.D, s)
	proof := writeAny(t, f.R.WorkspaceRoot, "resolved-callback.json", map[string]any{"modeled_result": "rejected before publication"})
	f.D.CWD = func() (string, error) { return f.R.WorkspaceRoot, nil }
	al14CLIOK(t, binary, f, "run", "resolve-stop", o.RunID, "--stop", s.ReportProviderStop.SHA256, "--source", sha, "--effect-status", DeliveryEffectNotStarted, "--basis", "action_result", "--observer", "modeled coordinator", "--reason", "Actual modeled rejection result", "--evidence", proof)
	target := workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}
	q, err := workspace.ListTaskQueue(f.D.Workspace, target, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = workspace.CloseTaskQueue(f.D.Workspace, target, s.Request.PreparationID, q.Revision, "Historical supported closure", "advance"); err != nil {
		t.Fatal(err)
	}
	al01QueueReplacementTask(t, f)
	preview, err := workspace.PrepareTask(f.D.Workspace, workspace.TaskPrepareInput{Target: target, Selector: workspace.QueueSelector{Kind: "next"}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = workspace.PrepareTask(f.D.Workspace, workspace.TaskPrepareInput{Target: target, Selector: workspace.QueueSelector{Kind: "next"}, Apply: true, Confirmation: *preview.Confirmation})
	if err != nil {
		t.Fatal(err)
	}
	before, err := workspace.ListTaskQueue(f.D.Workspace, target, false)
	if err != nil {
		t.Fatal(err)
	}
	workflowTestModel(t, f, map[string]any{"closed_transport": true})
	al14Retire(t, binary, f, o)
	al14Retire(t, binary, f, o)
	after, err := workspace.ListTaskQueue(f.D.Workspace, target, false)
	if err != nil {
		t.Fatal(err)
	}
	if before.Revision != after.Revision || !equal(before.Current, after.Current) {
		t.Fatal("historical retirement touched new queue owner")
	}
	// Actual CLI rejects a conflicting retry without provider or Git replay.
	c := exec.Command(binary, "workflow", "run", "retire", o.RunID, "--observer", "other", "--reason", "conflict")
	c.Dir = f.R.WorkspaceRoot
	out, err := al14Combined(t, c)
	if err == nil || !strings.Contains(string(out), "immutable") {
		t.Fatalf("conflicting retirement: %v %s", err, out)
	}
}
