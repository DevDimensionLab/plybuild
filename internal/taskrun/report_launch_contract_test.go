package taskrun

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
	"testing"
)

// The provider/Herdr are deterministic stand-ins; all commands in this journey
// execute the consumed candidate CLI against real disposable native stores.
func TestAL11ReportLaunchContractCLIJourney(t *testing.T) {
	binary := automaticCLIBinary(t)
	for _, mode := range []string{"codex", ClaudePermissionManual, ClaudePermissionAuto} {
		t.Run(mode, func(t *testing.T) {
			provider := "claude"
			if mode == "codex" {
				provider = "codex"
			}
			f := workflowProviderFixture(t, provider)
			if provider == "claude" {
				f.R.Runtime.PermissionBinding.ProfileID = mode
			}
			f = al01ReceiptFixtureFrom(t, f, binary)
			f.R.Runtime.PermissionBinding.AuthorityKind = "launch_contract_pending_runtime_acceptance"
			f.R.Runtime.PermissionBinding.EffectivePolicySHA256 = ""
			f.R.Runtime.PermissionBinding.Evidence = []Evidence{}
			f.File = writeAny(t, f.R.WorkspaceRoot, "al11-request.json", f.R)
			f.Report.RunID, f.Report.RequestSHA256, f.Report.SessionID = workflowID(f.R), digest(f.R), "ply:"+workflowID(f.R)
			if mode == ClaudePermissionManual {
				bad := f.R
				bad.Runtime.PermissionBinding.ProfileID = "invalid"
				file := writeAny(t, f.R.WorkspaceRoot, "al11-invalid-profile.json", bad)
				c := exec.Command(binary, "workflow", "run", "start", "--file", file, "--check", "--format", "json")
				c.Dir = f.R.WorkspaceRoot
				out, err := c.CombinedOutput()
				t.Logf("invalid profile CLI exit=%d output=%s", c.ProcessState.ExitCode(), out)
				if err == nil || !strings.Contains(string(out), "unsupported Claude permission mode") {
					t.Fatalf("profile rejection: %v %s", err, out)
				}
			}
			preview := al01PreviewCLI(t, binary, f.R.WorkspaceRoot, f.File)
			o := al11CLIOK(t, binary, f, "run", "start", "--file", f.File, "--apply", "--confirm", *preview.Confirmation)
			// Lost responses are observed by retrying exactly the same frozen input.
			retry := al11CLIOK(t, binary, f, "run", "start", "--file", f.File, "--apply", "--confirm", *preview.Confirmation)
			if retry.RunID != o.RunID || workflowTestCalls(t, f, "agent start") != 1 {
				t.Fatal("retry launched another owner")
			}
			for _, bad := range [][]string{{"--provider", "wrong", "--native-session", o.Transport.AgentSessionID}, {"--provider", provider, "--native-session", "wrong"}} {
				c := exec.Command(binary, append([]string{"workflow", "run", "check-in", o.RunID, "--context", o.Paths.Context}, bad...)...)
				c.Dir, _ = f.D.CWD()
				out, err := c.CombinedOutput()
				if err == nil || !strings.Contains(string(out), "provider or native session differs") {
					t.Fatalf("identity rejection: %v %s", err, out)
				}
				t.Logf("identity CLI raw exit=%d output=%s", c.ProcessState.ExitCode(), out)
			}
			args := []string{"run", "check-in", o.RunID, "--context", o.Paths.Context, "--provider", provider, "--native-session", o.Transport.AgentSessionID}
			o = al11CLIOK(t, binary, f, args...)
			first := o.ReportStartup.Receipt
			o = al11CLIOK(t, binary, f, args...)
			if first == nil || !equal(first, o.ReportStartup.Receipt) || o.ReportStartup.RuntimeDecision != "per_action_not_observed" || o.Delivery != nil {
				t.Fatalf("receipt manufactured authority: %+v", o)
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			raw, err := workflowBound(*s.Acceptance, 1<<20)
			if err != nil {
				t.Fatal(err)
			}
			var a Acceptance
			if err = json.Unmarshal(raw, &a); err != nil {
				t.Fatal(err)
			}
			if a.RuntimeClaim.ModelID != nil || a.RuntimeClaim.ProfileID != nil || a.RuntimeClaim.EffectivePolicySHA256 != nil || !strings.Contains(string(a.Sandbox), `"matches_contract":null`) {
				t.Fatalf("invented runtime facts: %s", raw)
			}
			report := al10Report(t, f, o, "al11-blocked", "blocked", "failed", nil, nil)
			o = al11CLIOK(t, binary, f, "run", "report", o.RunID, "--context", o.Paths.Context, "--file", report)
			if o.TaskResultState != "not_published" || o.FinalReturn.TerminalSHA256 != nil || o.Delivery != nil {
				t.Fatal("report acquired qualification or delivery authority")
			}
			t.Logf("candidate CLI exit=0: %s preview/start/retry/check-in/retry/report; policy document absent; runtime facts unknown", mode)
		})
	}
}

func TestAL11LaunchContractBoundaries(t *testing.T) {
	f := workflowProviderFixture(t, "claude")
	f.R.StartupProtocol = MandateReceiptV1
	f.R.Runtime.PermissionBinding.AuthorityKind = "launch_contract_pending_runtime_acceptance"
	f.R.Runtime.PermissionBinding.EffectivePolicySHA256 = ""
	f.R.Runtime.PermissionBinding.Evidence = []Evidence{}
	read := func() error {
		_, err := ReadWorkflowRequest(writeAny(t, f.R.WorkspaceRoot, "al11-boundary.json", f.R))
		return err
	}
	if err := read(); err != nil {
		t.Fatal(err)
	}
	f.R.Runtime.PermissionBinding.ProfileID = "invalid"
	if err := read(); err == nil || !strings.Contains(err.Error(), "unsupported Claude permission mode") {
		t.Fatalf("profile: %v", err)
	}
	f.R.Runtime.PermissionBinding.ProfileID = ClaudePermissionManual
	f.R.StartupProtocol = ReportLegacyAcceptanceV1
	// A handoff marker cannot override the top-level frozen legacy protocol.
	var draft map[string]any
	json.Unmarshal(f.R.HandoffDraft, &draft)
	draft["report_startup_protocol"] = MandateReceiptV1
	f.R.HandoffDraft, _ = Canonical(draft)
	if err := read(); err == nil || !strings.Contains(err.Error(), "effective policy evidence is required") {
		t.Fatalf("legacy weakened: %v", err)
	}
	f.R.StartupProtocol = MandateReceiptV1
	evidence := Evidence{Locator: f.File, SHA256: hashFileTest(t, f.File), Role: "runtime_contract"}
	f.R.Runtime.PermissionBinding.Evidence = []Evidence{evidence}
	f.R.Runtime.PermissionBinding.EffectivePolicySHA256 = evidence.SHA256
	if err := read(); err != nil {
		t.Fatalf("chosen evidence refused: %v", err)
	}
	if _, err := WorkflowPreviewStart(f.D, writeAny(t, f.R.WorkspaceRoot, "al11-optional.json", f.R)); err != nil {
		t.Fatal(err)
	}
	f.R.Runtime.PermissionBinding.EffectivePolicySHA256 = hash([]byte("unbound"))
	if err := read(); err == nil || !strings.Contains(err.Error(), "optional policy digest requires matching evidence") {
		t.Fatalf("unbound optional digest: %v", err)
	}
	f.R.Runtime.PermissionBinding.EffectivePolicySHA256 = ""
	f.R.Runtime.PermissionBinding.Evidence = []Evidence{{Locator: f.File, SHA256: hash([]byte("wrong")), Role: "runtime_contract"}}
	if err := read(); err != nil {
		t.Fatal(err)
	}
	if _, err := WorkflowPreviewStart(f.D, writeAny(t, f.R.WorkspaceRoot, "al11-evidence.json", f.R)); err == nil || !strings.Contains(err.Error(), "policy evidence changed") {
		t.Fatalf("optional evidence not checked: %v", err)
	}
	// Ordinary terminal parsing cannot use the pending report launch contract.
	native := workflowNativeRequest(f.R)
	native.Runtime.Provider = "codex"
	native.Runtime.PermissionBinding.Evidence = []Evidence{}
	raw, _ := Canonical(native)
	if _, err := parseRequest(raw); err == nil {
		t.Fatal("terminal gained pending report authority")
	}
}

func al11CLIOK(t *testing.T, binary string, f workflowFixture, args ...string) WorkflowRun {
	t.Helper()
	c := exec.Command(binary, append([]string{"workflow"}, append(args, "--format", "json")...)...)
	c.Dir, _ = f.D.CWD()
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run()
	t.Logf("CLI %v exit=%d stdout=%s stderr=%s", args, c.ProcessState.ExitCode(), stdout.Bytes(), stderr.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	var o WorkflowRun
	if err = json.Unmarshal(stdout.Bytes(), &o); err != nil {
		t.Fatal(err)
	}
	return o
}
