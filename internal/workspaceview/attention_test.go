package workspaceview

import (
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
)

func TestUnknownProcessHistoryDoesNotInventProgressOrHideRecordedFacts(t *testing.T) {
	s := progressFixture()
	process := map[taskwork.TaskID]ProcessFacts{"task": {Freshness: "unknown", Reasons: []string{"journal_unavailable"}}}
	p, err := BuildTaskProgress(s, "task", process)
	if err != nil || p.HasProgress || p.State != "nothing" || p.Freshness != "unknown" || len(p.Reasons) != 1 {
		t.Fatalf("unknown became invented work: %+v %v", p, err)
	}
	s.Registry.IntegrationAuthorities = []taskwork.IntegrationAuthority{fixtureAuthority()}
	p, err = BuildTaskProgress(s, "task", process)
	if err != nil || !p.HasProgress || p.Freshness != "unknown" {
		t.Fatalf("optional source hid native facts: %+v %v", p, err)
	}
}
