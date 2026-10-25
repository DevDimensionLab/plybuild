package taskrun

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// F06: request, binding and journal own a legacy trn_ result; result.json is
// derived output only. Damaged or forged cache bytes neither block readback
// nor grant qualification, liveness or owner authority.
func damageResultCache(t *testing.T, r Request, mode string, journalResult Result) {
	t.Helper()
	cache := filepath.Join(runPaths(r).RunRoot, "result.json")
	var e error
	switch mode {
	case "missing":
		e = os.Remove(cache)
	case "malformed":
		e = os.WriteFile(cache, []byte("{broken"), 0600)
	case "forged_authority":
		forged := journalResult
		forged.Collection.State = "qualified"
		forged.TaskExecution = TaskExecution{State: "inactive", Source: "human_observation"}
		forged.Process.State = "exited"
		forged.Process.Quiescence = ptr(true)
		forged.Delivery.State = "received"
		e = replaceValue(cache, forged)
	default:
		t.Fatalf("unknown cache damage %q", mode)
	}
	if e != nil {
		t.Fatal(e)
	}
}

func TestF06JournalOwnsResultDespiteDamagedCache(t *testing.T) {
	for _, mode := range []string{"missing", "malformed", "forged_authority"} {
		t.Run(mode, func(t *testing.T) {
			d, r, f := fixture(t)
			completeInside(t, &d, r, nil)
			p, e := PreviewStart(d, f)
			if e != nil {
				t.Fatal(e)
			}
			_, _ = Start(d, f, *p.(Preview).Confirmation)
			id := RunID(r.RequestKey)
			baseline, e := Show(d, r.WorkspaceRoot, id)
			if e != nil {
				t.Fatal(e)
			}
			baselineCollect, e := PreviewCollect(d, r.WorkspaceRoot, id)
			if e != nil {
				t.Fatal(e)
			}
			if baseline.Collection.State == "qualified" {
				t.Fatalf("fixture must not be qualified before damage: %+v", baseline.Collection)
			}
			runner := d.Runner.(*fakeRunner)
			calls := runner.calls
			damageResultCache(t, r, mode, baseline)
			state := treeState(t, r.WorkspaceRoot)
			got, e := Show(d, r.WorkspaceRoot, id)
			if e != nil {
				t.Fatal(e)
			}
			if !equal(got, baseline) {
				t.Fatalf("derived cache changed the journal result:\n%+v\n%+v", got, baseline)
			}
			collect, e := PreviewCollect(d, r.WorkspaceRoot, id)
			if e != nil {
				t.Fatal(e)
			}
			if !equal(collect, baselineCollect) {
				t.Fatal("derived cache granted collection authority")
			}
			replay, e := Start(d, f, "replay")
			if e != nil && !isCode(e, "task_run_task_status_unknown") {
				t.Fatal(e)
			}
			if replay.RunID != id || replay.Collection.State == "qualified" || replay.Process.State != baseline.Process.State {
				t.Fatalf("replay did not return the same journal Run: %+v", replay)
			}
			if state != treeState(t, r.WorkspaceRoot) {
				t.Fatal("readback wrote or repaired files")
			}
			if runner.calls != calls {
				t.Fatal("replay launched another provider")
			}
		})
	}
}

func TestF06DamagedJournalStillRejectsWithDamagedCache(t *testing.T) {
	for _, mode := range []string{"broken_event", "chain_mismatch"} {
		t.Run(mode, func(t *testing.T) {
			d, r, f := fixture(t)
			p, e := PreviewStart(d, f)
			if e != nil {
				t.Fatal(e)
			}
			_, _ = Start(d, f, *p.(Preview).Confirmation)
			id := RunID(r.RequestKey)
			baseline, e := Show(d, r.WorkspaceRoot, id)
			if e != nil {
				t.Fatal(e)
			}
			calls := d.Runner.(*fakeRunner).calls
			j := journalEvents(t, r)
			damageResultCache(t, r, "malformed", baseline)
			event := filepath.Join(runPaths(r).RunRoot, "events", "000002.json")
			if mode == "broken_event" {
				e = os.WriteFile(event, []byte("{broken"), 0600)
			} else {
				ev := j.Events[1]
				ev.PreviousSHA256 = ptr(hash([]byte("other chain")))
				raw, _ := Canonical(ev)
				e = os.WriteFile(event, raw, 0600)
			}
			if e != nil {
				t.Fatal(e)
			}
			if _, e = Show(d, r.WorkspaceRoot, id); !isCode(e, "task_run_integrity") {
				t.Fatalf("damaged journal accepted: %v", e)
			}
			if _, e = Start(d, f, "replay"); e == nil {
				t.Fatal("damaged journal permitted replay")
			}
			if d.Runner.(*fakeRunner).calls != calls {
				t.Fatal("damaged journal launched another provider")
			}
		})
	}
}

func TestF06CallbackBindingAndInterruptionWithDamagedCache(t *testing.T) {
	d, r, f := fixture(t)
	runner := d.Runner.(*fakeRunner)
	id := RunID(r.RequestKey)
	runner.inside = func(LaunchSpec) error {
		j, e := readJournal(r.WorkspaceRoot, id)
		if e != nil {
			return e
		}
		damageResultCache(t, r, "malformed", j.Result)
		bad := acceptance(t, d, r)
		bad.SessionID = "other"
		if _, e = Accept(d, r.WorkspaceRoot, id, writeAny(t, r.WorkspaceRoot, "bad-claim.json", bad)); e == nil {
			return errors.New("accepted wrong callback session binding")
		}
		d.Fault = func(p string) error {
			if p == "after_event_acceptance_received" {
				return errors.New("synthetic interruption before cache publication")
			}
			return nil
		}
		_, e = Accept(d, r.WorkspaceRoot, id, writeAny(t, r.WorkspaceRoot, "claim.json", acceptance(t, d, r)))
		d.Fault = nil
		if e == nil {
			return errors.New("interruption fault did not fire")
		}
		return nil
	}
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	_, _ = Start(d, f, *p.(Preview).Confirmation)
	if runner.lastErr != nil {
		t.Fatal(runner.lastErr)
	}
	// The interrupted callback left the acceptance event with a stale and then
	// malformed cache. The journal still owns the accepted claim and the Run.
	damageResultCache(t, r, "malformed", Result{})
	state := treeState(t, r.WorkspaceRoot)
	out, e := Show(d, r.WorkspaceRoot, id)
	if e != nil {
		t.Fatal(e)
	}
	if out.Acceptance.State != "started" || out.RuntimeFacts.Source == nil || out.Collection.State == "qualified" {
		t.Fatalf("journal claim provenance lost: %+v", out)
	}
	replay, e := Start(d, f, "replay")
	if e != nil && !isCode(e, "task_run_task_status_unknown") {
		t.Fatal(e)
	}
	if replay.RunID != id || runner.calls != 1 {
		t.Fatalf("interrupted Run was not retained without relaunch: calls=%d %+v", runner.calls, replay)
	}
	if state != treeState(t, r.WorkspaceRoot) {
		t.Fatal("readback wrote or repaired files")
	}
}

func TestF06CachePublicationKeepsPhysicalWriteGuard(t *testing.T) {
	d, r, f := fixture(t)
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	_, _ = Start(d, f, *p.(Preview).Confirmation)
	id := RunID(r.RequestKey)
	baseline, e := Show(d, r.WorkspaceRoot, id)
	if e != nil {
		t.Fatal(e)
	}
	cache := filepath.Join(runPaths(r).RunRoot, "result.json")
	decoy := filepath.Join(r.WorkspaceRoot, "decoy.json")
	if e = os.WriteFile(decoy, []byte("decoy"), 0600); e != nil {
		t.Fatal(e)
	}
	if e = os.Remove(cache); e != nil {
		t.Fatal(e)
	}
	if e = os.Symlink(decoy, cache); e != nil {
		t.Fatal(e)
	}
	got, e := Show(d, r.WorkspaceRoot, id)
	if e != nil {
		t.Fatal(e)
	}
	if !equal(got, baseline) {
		t.Fatal("ignored cache changed the journal result")
	}
	// An ignored read cache is no exemption: the cache publication used by
	// appendEvent still refuses an unsafe physical path before writing.
	if e = d.replaceValue(cache, baseline); e == nil || !strings.Contains(e.Error(), "symlink path rejected") {
		t.Fatalf("published through a symlinked cache path: %v", e)
	}
	if b, _ := os.ReadFile(decoy); string(b) != "decoy" {
		t.Fatal("cache publication wrote through the symlink")
	}
	if info, e := os.Lstat(cache); e != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatal("unsafe cache path was replaced")
	}
}

func journalEvents(t *testing.T, r Request) journal {
	t.Helper()
	j, e := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
	if e != nil {
		t.Fatal(e)
	}
	return j
}

func isCode(e error, code string) bool {
	var typed *Error
	return errors.As(e, &typed) && typed.Code == code
}
