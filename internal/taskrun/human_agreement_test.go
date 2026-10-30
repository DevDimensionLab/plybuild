package taskrun

import "testing"

func TestTaskRunUsesTheHumanSuppliedAgreement(t *testing.T) {
	d, r, _ := fixture(t)
	r.Agreement = Agreement{"human-selected", 5, 7200, 4}
	file := writeAny(t, r.WorkspaceRoot, "human-agreement.json", r)
	before := treeState(t, r.WorkspaceRoot)
	if _, err := PreviewStart(d, file); err != nil {
		t.Fatalf("human-selected limits were rejected: %v", err)
	}
	if before != treeState(t, r.WorkspaceRoot) {
		t.Fatal("agreement preview wrote state")
	}
	for _, tc := range []struct {
		usage BudgetUsage
		want  bool
	}{
		{BudgetUsage{true, 5, 7200, 4}, true},
		{BudgetUsage{true, 6, 7200, 4}, false},
		{BudgetUsage{true, 5, 7201, 4}, false},
		{BudgetUsage{true, 5, 7200, 5}, false},
	} {
		got := budgetValid(&tc.usage, r.Agreement)
		if got == nil || *got != tc.want {
			t.Errorf("usage %+v: got %v, want %t", tc.usage, got, tc.want)
		}
	}
}
