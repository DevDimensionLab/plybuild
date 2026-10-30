package taskrun

import (
	"bytes"

	"os"
	"path/filepath"
	"testing"
)

// TC0E6E1FE: fixtures exercise real consumers with historical requests intact.
func TestPlyUpgradePreparedStart(t *testing.T) {
	for _, change := range []string{"changed", "removed"} {
		t.Run(change, func(t *testing.T) {
			f := newDeliveryFixture(t, "codex")
			p, err := WorkflowPreviewStart(f.D, f.File)
			if err != nil {
				t.Fatal(err)
			}
			if change == "removed" {
				err = os.Remove(f.R.Runtime.PlyExecutable.Path)
			} else {
				err = os.WriteFile(f.R.Runtime.PlyExecutable.Path, []byte("#!/bin/sh\n# upgraded Ply\nexit 0\n"), 0700)
			}
			if err != nil {
				t.Fatal(err)
			}
			o, err := WorkflowStart(f.D, f.File, *p.(WorkflowPreview).Confirmation)
			if err != nil {
				t.Fatalf("Ply upgrade invalidated prepared start: %v", err)
			}
			again, err := WorkflowStart(f.D, f.File, *p.(WorkflowPreview).Confirmation)
			if err != nil || again.RunID != o.RunID {
				t.Fatalf("replay: %+v %v", again, err)
			}
			if n := workflowTestCalls(t, f.workflowFixture, "agent start"); n != 1 {
				t.Fatalf("provider starts = %d", n)
			}
		})
	}
}

func TestPlyUpgradeExistingOwnerCallbacks(t *testing.T) {
	for _, change := range []string{"changed", "removed"} {
		t.Run(change, func(t *testing.T) {
			f := newDeliveryFixture(t, "codex")
			o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
			frozen := map[string][]byte{}
			for _, p := range []string{workflowIndex(f.R.WorkspaceRoot, o.RunID), o.Paths.Context, filepath.Join(o.Paths.RunRoot, "acceptance.json"), filepath.Join(o.Paths.RunRoot, "mandate.json")} {
				raw, err := os.ReadFile(p)
				if err != nil {
					t.Fatal(err)
				}
				frozen[p] = raw
			}
			if change == "removed" {
				if err := os.Remove(f.R.Runtime.PlyExecutable.Path); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := os.WriteFile(f.R.Runtime.PlyExecutable.Path, []byte("#!/bin/sh\n# replacement\nexit 0\n"), 0700); err != nil {
					t.Fatal(err)
				}
			}
			installed := filepath.Join(f.R.WorkspaceRoot, "upgraded-ply")
			if err := os.WriteFile(installed, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
				t.Fatal(err)
			}
			f.D.Executable = func() (string, error) { return installed, nil }
			r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "ply-upgrade", Phase: "working", Summary: "TC0E6E1FE upgraded Ply callback.", Meaning: "Same owner and mandate without rebinding.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
			var err error
			o, err = WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "ply-upgrade-report.json", r))
			if err != nil {
				t.Fatalf("upgraded owner callback: %v", err)
			}
			candidate := deliveryTestCandidate(t, f, "Ply upgrade candidate")
			o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
			if err != nil || len(o.Delivery.Candidates) != 1 || o.Delivery.Candidates[0].OID != candidate {
				t.Fatalf("upgraded verification: %v", err)
			}
			for p, want := range frozen {
				got, err := os.ReadFile(p)
				if err != nil || !bytes.Equal(got, want) {
					t.Fatalf("history rewritten: %s %v", p, err)
				}
			}
			if n := workflowTestCalls(t, f.workflowFixture, "agent start"); n != 1 {
				t.Fatalf("provider starts = %d", n)
			}
		})
	}
}

func TestPlyUpgradeStoredContinuation(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	d, in, _ := continuationTestInput(t, f, o)
	p, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in)
	if err != nil {
		t.Fatal(err)
	}
	continued, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{f.R.Runtime.PlyExecutable.Path, continued.ControlExecutable.Path} {
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}
	d.Executable = func() (string, error) { return in.ControlExecutable.Path, nil }
	r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "stored-ply-upgrade", Phase: "working", Summary: "TC0E6E1FE old continuation remains readable.", Meaning: "Removed historical controllers are evidence only.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	if _, err = WorkflowDeliveryReport(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "stored-ply-upgrade-report.json", r)); err != nil {
		t.Fatal(err)
	}
}

func TestPlyUpgradeTerminalCallback(t *testing.T) {
	d, r, file := fixture(t)
	d.Runner.(*fakeRunner).inside = func(LaunchSpec) error {
		if err := os.Remove(r.Runtime.PlyExecutable.Path); err != nil {
			return err
		}
		d.Executable = func() (string, error) { return "/ordinary/upgraded/ply", nil }
		a := acceptance(t, d, r)
		_, err := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "upgraded-terminal-claim.json", a))
		return err
	}
	p, err := PreviewStart(d, file)
	if err != nil {
		t.Fatal(err)
	}
	out, err := Start(d, file, *p.(Preview).Confirmation)
	if d.Runner.(*fakeRunner).lastErr != nil || out.Acceptance.State != "started" {
		t.Fatalf("terminal callback required old Ply: %+v %v %v", out.Acceptance, err, d.Runner.(*fakeRunner).lastErr)
	}
}
