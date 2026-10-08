package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// This source includes native continuation but still requires the original
// launcher on disk. The older bound run control predates continuation itself.
const continuityBeforeUpgradeSource = "6472bb87e11de6512adcde65c03b1be41ce46646"

type continuityProviderUpgrade struct {
	Scenario           string        `json:"scenario"`
	HistoricalProvider Executable    `json:"historical_provider"`
	InstalledProvider  Executable    `json:"installed_provider"`
	ProviderLookup     string        `json:"provider_lookup"`
	RelocatedProvider  string        `json:"relocated_provider,omitempty"`
	BeforeControl      Executable    `json:"before_control"`
	BeforeSource       string        `json:"before_source_revision"`
	FailureEvidence    []FileBinding `json:"failure_evidence"`
	IncompleteReport   string        `json:"incomplete_report"`
	RuntimeObservation FileBinding   `json:"runtime_observation"`
}

type continuityCommandObservation struct {
	Argv   []string `json:"argv"`
	CWD    string   `json:"cwd"`
	Exit   int      `json:"exit"`
	Stdout string   `json:"stdout"`
	Stderr string   `json:"stderr"`
}

func continuityObserveCommand(t *testing.T, f continuityFixture, name, binary string, args ...string) (continuityCommandObservation, FileBinding) {
	t.Helper()
	argv := append([]string{binary}, append(args, "--format", "json")...)
	c := exec.Command(argv[0], argv[1:]...)
	c.Dir = f.CWD
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			t.Fatalf("observe native command: %v", err)
		}
	}
	o := continuityCommandObservation{Argv: argv, CWD: f.CWD, Exit: c.ProcessState.ExitCode(), Stdout: stdout.String(), Stderr: stderr.String()}
	path := writeAny(t, f.R.WorkspaceRoot, name+".json", o)
	return o, FileBinding{path, hashFileTest(t, path)}
}

func continuityUpgradeProvider(t *testing.T, f *continuityFixture, before, scenario string) {
	t.Helper()
	root := f.R.WorkspaceRoot
	u := &continuityProviderUpgrade{Scenario: scenario, HistoricalProvider: f.R.Runtime.Executable, ProviderLookup: filepath.Join(root, "bin", "codex"), BeforeSource: continuityBeforeUpgradeSource}
	control := filepath.Join(root, "pre-upgrade-ply-control")
	raw, err := os.ReadFile(before)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(control, raw, 0500); err != nil {
		t.Fatal(err)
	}
	u.BeforeControl = Executable{control, hashFileTest(t, control)}
	current := filepath.Join(root, "provider-versions", "0.2", "codex")
	if err = os.MkdirAll(filepath.Dir(current), 0700); err != nil {
		t.Fatal(err)
	}
	// This is a new synthetic installation, never a claim about the provider
	// version or authority of the already accepted native session.
	if err = os.WriteFile(current, []byte("#!/bin/sh\n# Synthetic current provider 0.2; never invoked by this journey.\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	u.InstalledProvider = Executable{current, hashFileTest(t, current)}
	if u.InstalledProvider.SHA256 == u.HistoricalProvider.SHA256 {
		t.Fatal("upgrade fixture did not install a distinct provider")
	}
	switch scenario {
	case "removed":
		err = os.Remove(u.HistoricalProvider.Path)
	case "relocated":
		u.RelocatedProvider = filepath.Join(root, "provider-versions", "historical-relocated")
		err = os.Rename(u.HistoricalProvider.Path, u.RelocatedProvider)
	case "replaced":
		err = os.WriteFile(u.HistoricalProvider.Path, []byte("#!/bin/sh\n# Different bytes at the historical path are not the accepted launcher.\nexit 0\n"), 0700)
	default:
		t.Fatalf("unknown isolated upgrade scenario: %s", scenario)
	}
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(u.ProviderLookup); err != nil {
		t.Fatal(err)
	}
	if err = os.Symlink(current, u.ProviderLookup); err != nil {
		t.Fatal(err)
	}
	resolved, err := exec.LookPath("codex")
	if err == nil {
		resolved, err = filepath.EvalSymlinks(resolved)
	}
	if err != nil || resolved != current {
		t.Fatalf("ordinary provider lookup does not reach the distinct installation: %q %v", resolved, err)
	}
	report := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: f.Run.RunID, RequestSHA256: f.Run.RequestSHA256, SessionID: f.Run.SessionID, EventID: "provider-upgrade-incomplete", PreviousEventSHA256: f.Run.Delivery.LastEventSHA256, Phase: "stopped", Summary: "Historical provider cleanup blocked the original callback.", Meaning: "Synthetic fixture: the preserved native callback could not validate its historical launcher after provider upgrade; this incomplete report grants no verification, human QA or integration authority.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	u.IncompleteReport = writeAny(t, root, "provider-upgrade-incomplete-report.json", report)
	statePath := filepath.Join(f.Run.Paths.RunRoot, "state.json")
	stateDigest := hashFileTest(t, statePath)
	for _, step := range []struct {
		name, binary string
		args         []string
	}{
		{"before-upgrade-continuation", control, []string{"workflow", "execute", "continue", f.Run.RunID, "--context", f.Run.Paths.Context, "--check"}},
		{"before-upgrade-report", f.Old, []string{"workflow", "execute", "report", f.Run.RunID, "--context", f.Run.Paths.Context, "--file", u.IncompleteReport}},
	} {
		observation, evidence := continuityObserveCommand(t, *f, step.name, step.binary, step.args...)
		combined := observation.Stdout + observation.Stderr
		if observation.Exit == 0 || !strings.Contains(combined, "executable bytes changed") && !strings.Contains(combined, "no such file or directory") {
			t.Fatalf("%s did not reproduce the real launcher dependency: %+v", step.name, observation)
		}
		if hashFileTest(t, statePath) != stateDigest {
			t.Fatalf("rejected %s changed native delivery state", step.name)
		}
		u.FailureEvidence = append(u.FailureEvidence, evidence)
	}
	observation := deliveryTestRuntimeObservation(t, f.deliveryFixture, f.Run)
	u.RuntimeObservation = FileBinding{observation, hashFileTest(t, observation)}
	f.Upgrade = u
	t.Logf("Synthetic provider %s: historical %s remains bound; current lookup resolves distinct %s; pre-fix continuation and original callback both failed without recording an effect", scenario, u.HistoricalProvider.Path, current)
}

func TestContinuityNativeProviderUpgradePreservesAcceptedOwner(t *testing.T) {
	current, old := deliveryCLIBinary(t), continuityOldBinary(t)
	before := continuitySourceBinary(t, continuityBeforeUpgradeSource, "PLY_CONTINUITY_BEFORE_UPGRADE_BINARY")
	for _, scenario := range []string{"removed", "replaced", "relocated"} {
		t.Run(scenario, func(t *testing.T) {
			f := newContinuityFixture(t, current, old, "")
			protected := map[string]string{}
			for _, path := range []string{f.Old, f.Run.Paths.Context, f.Run.Handoff.Locator, workflowIndex(f.R.WorkspaceRoot, f.Run.RunID), filepath.Join(f.Run.Paths.RunRoot, "acceptance.json"), filepath.Join(f.Run.Delivery.Attempt.Path, "verification.json")} {
				protected[path] = hashFileTest(t, path)
			}
			continuityUpgradeProvider(t, &f, before, scenario)
			missing, err := continuityCLI(current, f.CWD, "workflow", "execute", "continue", f.Run.RunID, "--context", f.Run.Paths.Context, "--check")
			if err == nil || !strings.Contains(string(missing), "delivery_runtime_observation_required") || !strings.Contains(string(missing), f.Upgrade.HistoricalProvider.Path) {
				t.Fatalf("missing current authority was not distinguished from the historical launcher: %v\n%s", err, missing)
			}
			originalEvents := len(f.Run.Delivery.Events)
			o := continuityRun(t, continuityCLIOK(t, current, f.CWD, "workflow", "execute", "report", f.Run.RunID, "--context", f.Run.Paths.Context, "--file", f.Upgrade.IncompleteReport, "--incomplete"))
			o = continuityRun(t, continuityCLIOK(t, current, f.CWD, "workflow", "execute", "report", f.Run.RunID, "--context", f.Run.Paths.Context, "--file", f.Upgrade.IncompleteReport, "--incomplete"))
			if o.Delivery.Phase != "stopped" || len(o.Delivery.Events) != originalEvents+1 || len(o.Delivery.Candidates) != 1 || o.Delivery.Candidates[0].OID != f.FirstOID || o.Delivery.Attempt.ID != f.Run.Delivery.Attempt.ID {
				t.Fatalf("bounded incomplete report changed qualification or lost evidence: %+v", o.Delivery)
			}
			state, err := workflowRead(f.R.WorkspaceRoot, f.Run.RunID)
			if err != nil || state.Continuation != nil {
				t.Fatalf("bounded incomplete report changed the active control: %v", err)
			}
			preview := continuityCLIOK(t, current, f.CWD, "workflow", "execute", "continue", f.Run.RunID, "--context", f.Run.Paths.Context, "--runtime-evidence", f.Upgrade.RuntimeObservation.Locator, "--check")
			var check struct {
				State string `json:"state"`
			}
			if err := json.Unmarshal(preview, &check); err != nil || check.State != "ready" {
				t.Fatalf("upgraded provider prevented continuation preview: %v\n%s", err, preview)
			}
			control := continuityContinue(t, f)
			if again := continuityContinue(t, f); again != control {
				t.Fatal("repeated continuation preserved a different control")
			}
			o = continuityRun(t, continuityCLIOK(t, control, f.CWD, "workflow", "execute", "verify", f.Run.RunID, "--context", f.Run.Paths.Context, "--review", f.Review, "--reuse", f.Run.Delivery.Attempt.ID))
			if len(o.Delivery.Candidates) != 2 || o.Delivery.Candidates[0].HumanQA == nil || o.Delivery.Candidates[0].HumanQA.Outcome != "fail" || o.Delivery.Candidates[1].OID != f.CorrectionOID || o.Delivery.Candidates[1].HumanQA != nil || o.Delivery.Phase != "awaiting_human_qa" {
				t.Fatalf("provider cleanup erased results or implied actual human QA: %+v", o.Delivery)
			}
			continuityCLIOK(t, control, f.CWD, "workflow", "execute", "verify", f.Run.RunID, "--context", f.Run.Paths.Context, "--review", f.Review, "--reuse", f.Run.Delivery.Attempt.ID)
			if _, err := continuityCLI(control, f.CWD, "workflow", "execute", "integrate", f.Run.RunID, "--context", f.Run.Paths.Context); err == nil {
				t.Fatal("provider continuation supplied a human pass")
			}
			if continuityCounter(t, f.Counter) != 2 || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "tab create") != 1 || o.Transport.AgentSessionID != f.Run.Transport.AgentSessionID {
				t.Fatal("provider cleanup reran checks, restarted the owner, or changed its native session")
			}
			if hashFileTest(t, f.Upgrade.InstalledProvider.Path) != f.Upgrade.InstalledProvider.SHA256 {
				t.Fatal("continuation changed the current provider installation")
			}
			if scenario == "replaced" {
				if hashFileTest(t, f.Upgrade.HistoricalProvider.Path) == f.Upgrade.HistoricalProvider.SHA256 {
					t.Fatal("continuation restored the historical launcher")
				}
			} else if _, err := os.Stat(f.Upgrade.HistoricalProvider.Path); !os.IsNotExist(err) {
				t.Fatalf("continuation restored the missing historical path: %v", err)
			}
			if scenario == "relocated" && hashFileTest(t, f.Upgrade.RelocatedProvider) != f.Upgrade.HistoricalProvider.SHA256 {
				t.Fatal("continuation modified relocated historical evidence")
			}
			for path, digest := range f.Unrelated {
				protected[path] = digest
			}
			for path, digest := range protected {
				if hashFileTest(t, path) != digest {
					t.Fatalf("provider cleanup rewrote immutable history: %s", path)
				}
			}
		})
	}
}
