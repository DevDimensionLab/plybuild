package weekrun

import "testing"

func TestAR02BoundWeekRetainsConcreteFailureForRecovery(t *testing.T) {
	f := newMR05Week(t)
	f.tick(t)
	f.tick(t)
	r := f.n.runs["first"]
	r.Phase, r.OwnerState, r.Releasable = "blocked", "unowned", true
	r.Reason = "Concrete failed attempt; diagnostic wording is arbitrary."
	f.n.runs["first"] = r
	f.tick(t)
	if len(f.n.released) != 0 || f.s.ActiveTaskID != "first" || len(f.n.starts) != 1 {
		t.Fatalf("bound Week adopted schema-1 automatic release: %+v", f.s)
	}
}
