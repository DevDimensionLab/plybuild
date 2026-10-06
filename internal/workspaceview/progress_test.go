package workspaceview

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func progressFixture() *Snapshot {
	return &Snapshot{
		Workspace: workspace.WorkspaceObservation{Root: "/work"},
		Projects:  workspace.ProjectSnapshot{Projects: []workspace.ProjectRecord{{ID: "project"}}, Repos: []workspace.RepoRecord{{ID: "repo"}}},
		Registry:  workspace.WorkItemRegistry{Tasks: []workspace.TaskRecord{{ID: "task", ParentEpicID: "epic", ProjectID: "project", RepoID: "repo", Title: "Feature", WorktreeState: workspace.WorkItemUnbound}}, Epics: []workspace.EpicRecord{{ID: "epic", ProjectID: "project", Title: "Epic"}}},
		Titles:    map[workspace.TaskID]workspace.TaskListTitle{"task": {Title: stringValue("Feature"), Source: "registration", Status: "available"}},
		Freshness: "fresh", Reasons: []string{},
	}
}

func fixtureResult(gate string) workspace.TaskResultRecord {
	return workspace.TaskResultRecord{ID: "result", TaskID: "task", ResultOID: "result-oid", ResultTree: "result-tree", TechnicalGate: gate, Recorder: workspace.TaskRecorderRecord{RecordedAtUTC: "2026-10-06T08:00:00Z"}}
}

func fixtureQA(outcome string) workspace.TaskHumanQARecord {
	return workspace.TaskHumanQARecord{ID: "qa", TaskID: "task", TaskResultID: "result", ResultOID: "result-oid", ResultTree: "result-tree", Outcome: outcome, Actor: workspace.TaskHumanActorRecord{CompletedAtUTC: "2026-10-06T09:00:00Z"}}
}

func TestProgressDoesNotTurnUntouchedOrPlannedTasksIntoNeedsYou(t *testing.T) {
	for _, planned := range []bool{false, true} {
		s := progressFixture()
		if planned {
			s.Registry.TaskProblemRevisions = []workspace.TaskProblemReference{{TaskID: "task"}}
			s.Registry.TaskSpecRevisions = []workspace.TaskSpecReference{{TaskID: "task"}}
		}
		rows, err := BuildTaskList(s, workspace.TaskListFilters{}, nil)
		if err != nil {
			t.Fatal(err)
		}
		progress := rows.Tasks[0].Progress
		if progress.State != "nothing" || progress.HasProgress || progress.TechnicalGate != nil || progress.HumanQAOutcome != nil || len(progress.NextActions) != 0 {
			t.Fatalf("planned=%v produced false execution facts: %+v", planned, progress)
		}
		attention, err := BuildAttention(s, workspace.TaskListFilters{}, nil)
		if err != nil || len(attention.Items) != 0 {
			t.Fatalf("attention=%+v err=%v", attention, err)
		}
	}
}

func TestRegisteredProgressSeparatesTechnicalHumanAndIntegrationFacts(t *testing.T) {
	for _, tc := range []struct{ gate, qa, integration, state, humanAction string }{
		{"passed", "", "", "in_progress", "record_human_qa"},
		{"good_enough_with_known_debt", "", "", "in_progress", "record_human_qa"},
		{"failed", "", "", "attention", ""},
		{"unknown", "", "", "in_progress", ""},
		{"passed", "fail", "", "attention", ""},
		{"passed", "blocked", "", "attention", ""},
		{"passed", "pass", "", "in_progress", ""},
		{"passed", "pass", "exact_effect", "integrated", ""},
		{"passed", "pass", "already_integrated", "integrated", ""},
		{"passed", "pass", "conflict", "attention", "resolve_integration_conflict"},
		{"passed", "pass", "partial", "attention", "resolve_integration_conflict"},
	} {
		t.Run(tc.gate+"/"+tc.qa+"/"+tc.integration, func(t *testing.T) {
			s := progressFixture()
			s.Registry.TaskResults = []workspace.TaskResultRecord{fixtureResult(tc.gate)}
			if tc.qa != "" {
				s.Registry.HumanQARecords = []workspace.TaskHumanQARecord{fixtureQA(tc.qa)}
			}
			if tc.integration != "" {
				s.Registry.IntegrationAuthorities = []workspace.IntegrationAuthority{{ID: "authority", TaskID: "task", TaskResultID: "result", HumanQARecordID: "qa"}}
				s.Registry.IntegrationResults = []workspace.IntegrationResult{{ID: "integration", AuthorityID: "authority", Outcome: tc.integration, RecordedAtUTC: "2026-10-06T10:00:00Z"}}
			}
			rows, err := BuildTaskList(s, workspace.TaskListFilters{}, nil)
			if err != nil {
				t.Fatal(err)
			}
			p := rows.Tasks[0].Progress
			if p.State != tc.state || !p.HasProgress || p.TechnicalGate == nil || *p.TechnicalGate != tc.gate || p.SourceBasis != "registered" || p.ResultCount != 1 {
				t.Fatalf("progress=%+v", p)
			}
			if (tc.qa == "") != (p.HumanQAOutcome == nil) || p.HumanQAOutcome != nil && *p.HumanQAOutcome != tc.qa {
				t.Fatalf("QA=%+v", p.HumanQAOutcome)
			}
			wantClass := tc.integration
			if wantClass == "" {
				wantClass = "unknown"
			}
			if p.IntegrationClassification != wantClass {
				t.Fatalf("unobserved Git readiness claimed: %+v", p)
			}
			a, err := BuildAttention(s, workspace.TaskListFilters{}, nil)
			if err != nil {
				t.Fatal(err)
			}
			if tc.humanAction == "" && len(a.Items) != 0 {
				t.Fatalf("agent work became human work: %+v", a)
			}
			if tc.humanAction != "" && (len(a.Items) != 1 || a.Items[0].Kind != tc.humanAction) {
				t.Fatalf("missing human action: %+v", a)
			}
		})
	}
}

func TestProcessProgressAndHumanWaitNeverManufactureQAPass(t *testing.T) {
	s := progressFixture()
	process := map[workspace.TaskID]ProcessFacts{"task": {HasProgress: true, Freshness: "fresh", Waiting: []ProcessWaiting{{Kind: "waiting", Actor: "human", Reason: "Choose the integration target", EventID: "event", SinceUTC: stringValue("2026-10-06T08:00:00Z")}}}}
	rows, err := BuildTaskList(s, workspace.TaskListFilters{}, process)
	if err != nil {
		t.Fatal(err)
	}
	p := rows.Tasks[0].Progress
	if p.State != "in_progress" || !p.HasProgress || p.TechnicalGate != nil || p.HumanQAOutcome != nil || p.ResultCount != 0 {
		t.Fatalf("process report became QA: %+v", p)
	}
	a, err := BuildAttention(s, workspace.TaskListFilters{}, process)
	if err != nil || len(a.Items) != 1 || a.Items[0].Source != "journal" || a.Items[0].Actor != "human" {
		t.Fatalf("attention=%+v err=%v", a, err)
	}
	show, err := BuildTaskProgress(s, "task", process)
	if err != nil || !reflect.DeepEqual(p, show) {
		t.Fatalf("list/show progress diverged: %+v %+v %v", p, show, err)
	}
}

func TestProgressJSONUsesNullsAndEmptyArrays(t *testing.T) {
	s := progressFixture()
	p, err := BuildTaskProgress(s, "task", nil)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err = json.Unmarshal(b, &fields); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"technical_gate", "human_qa_outcome", "last_activity_utc", "task_result_id", "integration_result_id"} {
		if v, present := fields[name]; !present || v != nil {
			t.Fatalf("%s not explicit null: %s", name, b)
		}
	}
	if next, ok := fields["next_actions"].([]any); !ok || len(next) != 0 {
		t.Fatalf("next actions not []: %s", b)
	}
}

func TestConflictingQAIsNotSilentlyPresentedAsPass(t *testing.T) {
	s := progressFixture()
	s.Registry.TaskResults = []workspace.TaskResultRecord{fixtureResult("passed")}
	pass, fail := fixtureQA("pass"), fixtureQA("fail")
	pass.ID, fail.ID = "pass", "fail"
	s.Registry.HumanQARecords = []workspace.TaskHumanQARecord{pass, fail}
	p, err := BuildTaskProgress(s, "task", nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.State != "attention" || p.IntegrationClassification != "conflict" || p.HumanQAOutcome != nil {
		t.Fatalf("contradictory human judgments became pass: %+v", p)
	}
	// A separate integration authority can select an exact QA/result binding.
	// The overview follows that recorded selection, without inventing one.
	s.Registry.IntegrationAuthorities = []workspace.IntegrationAuthority{{ID: "authority", TaskID: "task", TaskResultID: "result", HumanQARecordID: "pass"}}
	p, err = BuildTaskProgress(s, "task", nil)
	if err != nil || p.HumanQAOutcome == nil || *p.HumanQAOutcome != "pass" || p.IntegrationClassification != "unknown" {
		t.Fatalf("exact authority selection lost: %+v %v", p, err)
	}
}

func TestCompetingRecordedResultsNeverInviteQAOrClaimIntegration(t *testing.T) {
	for _, integrated := range []bool{false, true} {
		s := progressFixture()
		old, current := fixtureResult("passed"), fixtureResult("failed")
		current.ID, current.Recorder.RecordedAtUTC = "new-result", "2026-10-06T10:00:00Z"
		s.Registry.TaskResults = []workspace.TaskResultRecord{old, current}
		if integrated {
			s.Registry.HumanQARecords = []workspace.TaskHumanQARecord{fixtureQA("pass")}
			s.Registry.IntegrationAuthorities = []workspace.IntegrationAuthority{{ID: "authority", TaskID: "task", TaskResultID: "result", HumanQARecordID: "qa"}}
			s.Registry.IntegrationResults = []workspace.IntegrationResult{{ID: "integration", AuthorityID: "authority", Outcome: "exact_effect"}}
		}
		p, err := BuildTaskProgress(s, "task", nil)
		if err != nil {
			t.Fatal(err)
		}
		if p.State != "attention" || p.IntegrationClassification != "conflict" || p.TechnicalGate != nil || p.HumanQAOutcome != nil || p.TaskResultID != nil || p.IntegrationResultID != nil || p.ResultCount != 2 {
			t.Fatalf("integrated=%v: old success hid competing result: %+v", integrated, p)
		}
		for _, action := range p.NextActions {
			if action.Kind == "record_human_qa" || action.Kind == "check_integration" {
				t.Fatalf("invited next effect without unambiguous basis: %+v", action)
			}
		}
	}
}

func TestStaleRegisteredResultKeepsHistoryWithoutInvitingHumanQA(t *testing.T) {
	p := buildProgress(workspace.TaskRecord{ID: "task"}, workspace.TaskRecordedProgressFacts{HasProgress: true, ResultCount: 2, ResultBasisStale: true}, ProcessFacts{}, "fresh", nil)
	if p.State != "in_progress" || p.IntegrationClassification != "unknown" || p.TechnicalGate != nil || p.HumanQAOutcome != nil || p.ResultCount != 2 || p.TaskResultID != nil || p.IntegrationResultID != nil || !reflect.DeepEqual(p.Reasons, []string{"recorded_result_basis_stale"}) {
		t.Fatalf("stale history became current verification: %+v", p)
	}
	if len(p.NextActions) != 1 || p.NextActions[0].Actor != "agent" || p.NextActions[0].Kind != "inspect_task_basis" {
		t.Fatalf("stale basis did not produce a concrete agent follow-up: %+v", p.NextActions)
	}
}
