package taskrun

import (
	"os"
	"testing"
)

func TestAutomaticDeliveryCLIRetriesOnlyProvenNoEffectOnSameDelivery(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, true, 30)
	deliveryTestCandidate(t, f, "automatic no-effect retry")
	if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = 'automatic no-effect retry'\n"), 0600); err != nil {
		t.Fatal(err)
	}
	o = automaticCLIOK(t, binary, f, automaticCLIVerifyArgs(o, deliveryTestReview(t, f, "automatic-no-effect"), binary)...)
	cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	registered := deliveryCLIOK(t, binary, cwd, "register", "--file", deliveryCLIRegisterInput(t, f, o, "automatic-no-effect"))
	id := deliveryCLIString(t, registered, "id")
	lock := gitOutput(t, f.Parent, "rev-parse", "--path-format=absolute", "--git-path", "index.lock")
	if err := os.WriteFile(lock, []byte("fixture transient merge lock\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := deliveryCLIRun(binary, cwd, "execute", id); err == nil {
		t.Fatal("transient merge failure claimed completion")
	}
	r, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(r.IntegrationResults) != 1 || r.IntegrationResults[0].Outcome != "no_effect" || r.IntegrationResults[0].GitChanged == nil || *r.IntegrationResults[0].GitChanged {
		t.Fatalf("first attempt was not proven no-effect: results=%v err=%v", r.IntegrationResults, err)
	}
	prior := r.IntegrationResults[0].ID
	if err = os.Remove(lock); err != nil {
		t.Fatal(err)
	}
	done := deliveryCLIOK(t, binary, cwd, "execute", id)
	if deliveryCLIString(t, done, "state") != "delivered" {
		t.Fatalf("same Delivery did not resume a proven no-effect attempt: %+v", done)
	}
	r, err = f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(r.IntegrationAuthorities) != 2 || len(r.IntegrationAttempts) != 2 || len(r.IntegrationResults) != 2 || len(r.HumanQARecords) != 0 {
		t.Fatalf("retry history was not preserved: authorities=%d attempts=%d results=%d qa=%d err=%v", len(r.IntegrationAuthorities), len(r.IntegrationAttempts), len(r.IntegrationResults), len(r.HumanQARecords), err)
	}
	bound := false
	for _, a := range r.IntegrationAuthorities {
		if a.RetryAfterResultID != nil && *a.RetryAfterResultID == prior {
			bound = true
		}
	}
	if !bound || gitOutput(t, f.Parent, "rev-parse", "HEAD") != o.Delivery.Candidates[0].OID {
		t.Fatal("retry lost its exact predecessor or target")
	}
	before := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	deliveryCLIOK(t, binary, cwd, "execute", id)
	if after := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic"); after != before {
		t.Fatal("completed retry repeated integration")
	}
}
