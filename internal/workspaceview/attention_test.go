package workspaceview

import (
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestUnknownProcessHistoryDoesNotInventProgressOrHideRecordedFacts(t *testing.T) {
	s := progressFixture()
	process := map[workspace.TaskID]ProcessFacts{"task": {Freshness: "unknown", Reasons: []string{"journal_unavailable"}}}
	p, err := BuildTaskProgress(s, "task", process)
	if err != nil || p.HasProgress || p.State != "nothing" || p.Freshness != "unknown" || len(p.Reasons) != 1 {
		t.Fatalf("unknown became invented work: %+v %v", p, err)
	}
	s.Registry.IntegrationAuthorities = []workspace.IntegrationAuthority{fixtureAuthority()}
	p, err = BuildTaskProgress(s, "task", process)
	if err != nil || !p.HasProgress || p.Freshness != "unknown" {
		t.Fatalf("optional source hid native facts: %+v %v", p, err)
	}
}
