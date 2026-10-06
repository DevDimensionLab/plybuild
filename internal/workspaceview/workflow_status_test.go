package workspaceview

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func workflowFixture() *Snapshot {
	s := progressFixture()
	s.ObservedAtUTC = "2026-10-06T12:00:00Z"
	s.Projects.Projects[0].RepoIDs = []workspace.RepoID{"repo", "second-repo"}
	s.Projects.Projects = append(s.Projects.Projects, workspace.ProjectRecord{ID: "other", Name: "Other", RepoIDs: []workspace.RepoID{"other-repo"}}, workspace.ProjectRecord{ID: "empty", Name: "Empty"})
	s.Projects.Repos = append(s.Projects.Repos, workspace.RepoRecord{ID: "second-repo"}, workspace.RepoRecord{ID: "other-repo"})
	s.Registry.Epics[0].RepoBindings = []workspace.EpicRepoBinding{{RepoID: "repo"}, {RepoID: "second-repo"}}
	s.Registry.Epics = append(s.Registry.Epics, workspace.EpicRecord{ID: "empty-epic", ProjectID: "other", Title: "Empty Epic", RepoBindings: []workspace.EpicRepoBinding{{RepoID: "other-repo"}}})
	return s
}

func workflowBuild(t *testing.T, s *Snapshot, all bool, runs ...taskrun.InventoryRun) WorkflowStatus {
	t.Helper()
	out, err := BuildWorkflowStatus(s, WorkflowStatusOptions{IncludeAll: all}, nil, taskrun.Inventory{Runs: runs}, workspace.RegisteredTaskQueues{Freshness: "fresh"})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func workflowReportedRun(phase string) taskrun.InventoryRun {
	run := taskrun.InventoryRun{RunID: "run", TaskID: stringValue("task"), Provider: stringValue("codex"), State: "unknown", StateFreshness: "unknown", HistoricalState: stringValue("working"), Unresolved: true, Freshness: "fresh", Sources: []taskrun.FileBinding{}, Reasons: []taskrun.Reason{}}
	run.Delivery = &taskrun.InventoryDelivery{Phase: phase, Freshness: "fresh", Reports: []taskrun.InventoryDeliveryReport{}, Reasons: []taskrun.Reason{}}
	if phase == "needs_input" {
		report := taskrun.InventoryDeliveryReport{EventID: "question", Phase: phase, Question: stringValue("Which output should be the default?"), Current: true}
		run.Delivery.Report = &report
		run.Delivery.Reports = append(run.Delivery.Reports, report)
		run.Delivery.NextAction = taskrun.InventoryDeliveryAction{Actor: "human", RecordedActor: "user", Kind: "answer_question", Reason: *report.Question, EventID: &report.EventID, EvidenceIDs: []string{"run", "question"}}
	}
	if phase == "stopped" {
		run.Delivery.NextAction = taskrun.InventoryDeliveryAction{Actor: "agent", RecordedActor: "recipient", Kind: "follow_stopped_delivery", Reason: "Investigate the preserved stop", EvidenceIDs: []string{"run"}}
	}
	return run
}

func TestWorkflowStatusWorkspaceFiltersAndEmptyContainers(t *testing.T) {
	s := workflowFixture()
	out := workflowBuild(t, s, false)
	if len(out.Projects) != 3 || len(out.Epics) != 2 || out.Counts.Total != 1 || out.Counts.Hidden.Backlog != 1 || len(out.Items) != 0 {
		t.Fatalf("whole workspace or empty containers lost: %+v", out)
	}
	project, repo, epic := workspace.ProjectID("other"), workspace.RepoID("repo"), workspace.EpicID("empty-epic")
	out, err := BuildWorkflowStatus(s, WorkflowStatusOptions{Filters: workspace.TaskListFilters{ProjectID: &project, RepoID: &repo, EpicID: &epic}}, nil, taskrun.Inventory{}, workspace.RegisteredTaskQueues{})
	if err != nil || len(out.Items) != 0 || len(out.Projects) != 0 || len(out.Epics) != 0 {
		t.Fatalf("known disjoint AND filters: %+v %v", out, err)
	}
	epic = "missing"
	if _, err = BuildWorkflowStatus(s, WorkflowStatusOptions{Filters: workspace.TaskListFilters{ProjectID: &project, RepoID: &repo, EpicID: &epic}}, nil, taskrun.Inventory{}, workspace.RegisteredTaskQueues{}); err == nil {
		t.Fatal("unknown filter hidden by empty intersection")
	}
	project, repo, epic = "project", "second-repo", "epic"
	out, err = BuildWorkflowStatus(s, WorkflowStatusOptions{Filters: workspace.TaskListFilters{ProjectID: &project, RepoID: &repo, EpicID: &epic}}, nil, taskrun.Inventory{}, workspace.RegisteredTaskQueues{})
	if err != nil || len(out.Items) != 0 || len(out.Projects) != 1 || len(out.Epics) != 1 || len(out.Epics[0].Base) != 1 {
		t.Fatalf("empty registered target disappeared: %+v %v", out, err)
	}
}

func TestWorkflowStatusReportsKeepQuestionsActorsAndUnknownTimes(t *testing.T) {
	for _, phase := range []string{"working", "needs_input", "stopped"} {
		t.Run(phase, func(t *testing.T) {
			s := workflowFixture()
			run := workflowReportedRun(phase)
			out := workflowBuild(t, s, false, run)
			if len(out.Items) != 1 {
				t.Fatalf("reported work hidden: %+v", out)
			}
			row := out.Items[0]
			want := map[string]string{"working": "in_progress", "needs_input": "needs_you", "stopped": "follow_up"}[phase]
			if row.Category != want || row.Progress.TechnicalGate != nil || row.Progress.HumanQAOutcome != nil || row.Progress.IntegrationResultID != nil || row.LastActivityUTC != nil || row.Runs[0].StateFreshness != "unknown" {
				t.Fatalf("report became technical, QA, age or liveness fact: %+v", row)
			}
			if phase == "needs_input" && (len(row.NextActions) != 1 || row.NextActions[0].Actor != "human" || row.NextActions[0].Reason != "Which output should be the default?" || row.NextActions[0].SinceUTC != nil || row.NextActions[0].RunID == nil || !row.NextActions[0].Current) {
				t.Fatalf("question provenance lost: %+v", row)
			}
		})
	}
}

func TestWorkflowStatusNativeCandidateGatesAndCompetingResults(t *testing.T) {
	for _, tc := range []struct{ qa, integration, category, actor string }{
		{"", "", "needs_you", "human"}, {"fail", "", "follow_up", "agent"}, {"blocked", "", "follow_up", "agent"}, {"pass", "", "follow_up", "agent"}, {"pass", "exact_effect", "completed", ""},
	} {
		t.Run(tc.qa+tc.integration, func(t *testing.T) {
			s := workflowFixture()
			s.Registry.TaskResults = []workspace.TaskResultRecord{fixtureResult("passed")}
			if tc.qa != "" {
				s.Registry.HumanQARecords = []workspace.TaskHumanQARecord{fixtureQA(tc.qa)}
			}
			if tc.integration != "" {
				s.Registry.IntegrationAuthorities = []workspace.IntegrationAuthority{{ID: "authority", TaskID: "task", TaskResultID: "result", HumanQARecordID: "qa"}}
				s.Registry.IntegrationResults = []workspace.IntegrationResult{{ID: "integration", AuthorityID: "authority", Outcome: tc.integration}}
			}
			out := workflowBuild(t, s, true)
			if out.Items[0].Category != tc.category {
				t.Fatalf("gates collapsed: %+v", out.Items[0])
			}
			if tc.actor != "" && out.Items[0].NextActions[0].Actor != tc.actor {
				t.Fatalf("explicit next actor changed: %+v", out.Items[0])
			}
			if tc.integration != "" {
				hidden := workflowBuild(t, s, false)
				if hidden.Counts.Hidden.Completed != 1 || len(hidden.Items) != 0 {
					t.Fatalf("completed work not hidden: %+v", hidden)
				}
			}
		})
	}
	s := workflowFixture()
	first, second := fixtureResult("passed"), fixtureResult("failed")
	second.ID = "competing"
	s.Registry.TaskResults = []workspace.TaskResultRecord{first, second}
	out := workflowBuild(t, s, false)
	if out.Items[0].Category != "needs_you" || out.Items[0].Progress.TechnicalGate != nil || out.Items[0].NextActions[0].Kind != "resolve_integration_conflict" {
		t.Fatalf("competing results selected a winner: %+v", out)
	}
}

func TestWorkflowStatusSupersededAndStaleQuestionsRemainHistory(t *testing.T) {
	for _, stale := range []bool{false, true} {
		s := workflowFixture()
		run := workflowReportedRun("needs_input")
		if stale {
			s.Registry.TaskProblemRevisions = []workspace.TaskProblemReference{{TaskID: "task", Revision: 2, ManifestSHA256: "new"}}
			run.Delivery.Basis = &workspace.TaskSpecBasis{TaskID: "task", Problem: workspace.TaskRevisionRef{Revision: 1, ManifestSHA256: "old"}}
		} else {
			run.Delivery.Phase = "working"
			run.Delivery.Report.Current = false
			run.Delivery.Reports[0].Current = false
		}
		out := workflowBuild(t, s, false, run)
		if out.Items[0].Category == "needs_you" {
			t.Fatalf("historic question became a current human action: %+v", out.Items[0])
		}
		found := false
		for _, a := range out.Items[0].NextActions {
			if a.Kind == "answer_question" {
				found = true
				if a.Current {
					t.Fatal("historical question marked current")
				}
			}
		}
		if !found {
			t.Fatalf("historical question dropped: %+v", out.Items[0])
		}
	}
}

func TestWorkflowStatusQueueReadinessAndChangedSourceBasis(t *testing.T) {
	for _, state := range []string{"ready", "blocked", "unknown"} {
		for _, freshness := range []string{"fresh", "unknown"} {
			s := workflowFixture()
			s.Freshness = freshness
			queues := workspace.RegisteredTaskQueues{Freshness: "fresh", Queues: []workspace.RegisteredTaskQueue{{QueueID: "queue", Revision: 1, Target: workspace.QueueTarget{ProjectID: "project", RepoID: "repo", EpicID: "epic"}, Freshness: "fresh", Pending: []workspace.RegisteredQueuePending{{QueuePending: workspace.QueuePending{Rank: 1, TaskID: "task", State: state, SpecID: stringValue("goal")}, Spec: &workspace.TaskRevisionRef{Revision: 1, ManifestSHA256: "digest"}, Freshness: "fresh"}}}}}
			out, err := BuildWorkflowStatus(s, WorkflowStatusOptions{}, nil, taskrun.Inventory{}, queues)
			if err != nil {
				t.Fatal(err)
			}
			want := "follow_up"
			if state == "ready" && freshness == "fresh" {
				want = "ready_next"
			}
			if len(out.Items) != 1 || out.Items[0].Category != want || out.Items[0].Progress.HasProgress {
				t.Fatalf("state=%s source=%s false queue readiness/progress: %+v", state, freshness, out)
			}
		}
	}
}

func TestWorkflowStatusHiddenCountsDoNotOverlapOrHideDiagnostics(t *testing.T) {
	d, items, _ := viewDiskFixture(t)
	items.registry.TaskResults = []workspace.TaskResultRecord{fixtureResult("passed")}
	items.registry.HumanQARecords = []workspace.TaskHumanQARecord{fixtureQA("pass")}
	items.registry.IntegrationAuthorities = []workspace.IntegrationAuthority{{ID: "authority", TaskID: "task", TaskResultID: "result", HumanQARecordID: "qa"}}
	items.registry.IntegrationResults = []workspace.IntegrationResult{{ID: "integration", AuthorityID: "authority", Outcome: "exact_effect"}}
	if _, err := workspace.SetWorkItemLifecycle(d, workspace.WorkItemLifecycleInput{SubjectKind: "epic", SubjectID: "epic", State: workspace.LifecycleParked, ActorClaim: "fixture owner"}); err != nil {
		t.Fatal(err)
	}
	s, err := LoadSnapshot(d)
	if err != nil {
		t.Fatal(err)
	}
	run := workflowReportedRun("needs_input")
	run.Freshness = "unknown"
	run.Reasons = []taskrun.Reason{{Code: "broken_evidence", Detail: "damaged report"}}
	out := workflowBuild(t, s, false, run)
	if out.Counts.Hidden.Inactive != 1 || out.Counts.Hidden.Completed != 0 || out.Counts.Hidden.Backlog != 0 || len(out.Items) != 0 || len(out.Diagnostics) == 0 {
		t.Fatalf("hidden precedence or independent diagnostics lost: %+v", out)
	}
	all := workflowBuild(t, s, true, run)
	if all.Counts.Hidden != (WorkflowStatusHidden{}) || len(all.Items) != 1 || all.Items[0].Category != "inactive" || all.Items[0].TaskLifecycle != workspace.LifecycleActive || all.Items[0].Progress.State != "integrated" {
		t.Fatalf("--all changed lifecycle or native history: %+v", all)
	}
	for _, a := range all.Items[0].NextActions {
		if a.Current {
			t.Fatalf("historical action marked current: %+v", a)
		}
	}
}

func TestWorkflowStatusCorruptOrphanRunKeepsValidSibling(t *testing.T) {
	s := workflowFixture()
	valid := workflowReportedRun("needs_input")
	bad := taskrun.InventoryRun{RunID: "orphan", State: "unknown", StateFreshness: "unknown", Freshness: "unknown", Reasons: []taskrun.Reason{{Code: "run_evidence_unavailable", Detail: "invalid source"}}}
	out := workflowBuild(t, s, false, bad, valid)
	if len(out.Items) != 1 || out.Items[0].Category != "needs_you" || len(out.Items[0].Runs) != 1 || len(out.Diagnostics) < 2 || out.Freshness != "unknown" {
		t.Fatalf("bad source hid valid sibling or invented a Task: %+v", out)
	}
	if out.Diagnostics[0].TaskID != nil || out.Diagnostics[0].SinceUTC != nil {
		t.Fatalf("invented orphan identity/time: %+v", out.Diagnostics)
	}
}

func TestWorkflowStatusRepeatedReadsAreBatchedAndDoNotWrite(t *testing.T) {
	d, items, projects := viewDiskFixture(t)
	projects.projects.Projects[0].RepoIDs = []workspace.RepoID{"repo"}
	root, _ := d.Files.Getwd()
	subdir := filepath.Join(root, "arbitrary", "nested")
	if err := os.MkdirAll(subdir, 0700); err != nil {
		t.Fatal(err)
	}
	d.Files = viewTestFiles{d.Files, subdir}
	before := viewTree(t, root)
	for n := 0; n < 2; n++ {
		out, err := ReadWorkflowStatus(d, WorkflowStatusOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if out.Workspace.Root != root || out.Counts.Hidden.Backlog != 1 || items.reads != n+1 || projects.reads != n+1 || items.fullReads != 0 || items.locks != 0 {
			t.Fatalf("not one read per registry or accidental live/lock dependency: %+v items=%+v projects=%+v", out, items, projects)
		}
	}
	if !reflect.DeepEqual(before, viewTree(t, root)) {
		t.Fatal("status altered workspace sources")
	}
}

func TestWorkflowStatusJSONNullsCollectionsAndDeterminism(t *testing.T) {
	s := workflowFixture()
	run := workflowReportedRun("working")
	a, b := workflowBuild(t, s, true, run), workflowBuild(t, s, true, run)
	x, _ := json.Marshal(a)
	y, _ := json.Marshal(b)
	if string(x) != string(y) {
		t.Fatal("same captured facts yielded nondeterministic output")
	}
	var object map[string]any
	if err := json.Unmarshal(x, &object); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"projects", "epics", "items", "diagnostics"} {
		if _, ok := object[k].([]any); !ok {
			t.Fatalf("%s is not an array: %s", k, x)
		}
	}
	row := object["items"].([]any)[0].(map[string]any)
	if row["last_activity_utc"] != nil {
		t.Fatalf("read time invented activity: %s", x)
	}
	if !strings.Contains(string(x), `"observed_at_utc":"2026-10-06T12:00:00Z"`) {
		t.Fatal("missing separate source read time")
	}
}

// This measures the actual registered bulk pipeline, including journal/run
// inventories, at multiple-project scale. It is not a throughput improvement claim.
func BenchmarkWorkflowStatusMultiProject(b *testing.B) {
	root, err := filepath.EvalSymlinks(b.TempDir())
	if err != nil {
		b.Fatal(err)
	}
	d := workspace.SystemDependencies()
	d.Files = viewTestFiles{d.Files, root}
	if _, err = workspace.Init(d); err != nil {
		b.Fatal(err)
	}
	registry := workspace.WorkItemRegistry{}
	projects := workspace.ProjectSnapshot{}
	for p := 0; p < 4; p++ {
		project, repo := workspace.ProjectID(fmt.Sprintf("project-%d", p)), workspace.RepoID(fmt.Sprintf("repo-%d", p))
		projects.Projects = append(projects.Projects, workspace.ProjectRecord{ID: project, RepoIDs: []workspace.RepoID{repo}})
		projects.Repos = append(projects.Repos, workspace.RepoRecord{ID: repo})
		for e := 0; e < 3; e++ {
			epic := workspace.EpicID(fmt.Sprintf("epic-%d-%d", p, e))
			registry.Epics = append(registry.Epics, workspace.EpicRecord{ID: epic, ProjectID: project})
			for n := 0; n < 25; n++ {
				registry.Tasks = append(registry.Tasks, workspace.TaskRecord{ID: workspace.TaskID(fmt.Sprintf("task-%d-%d-%d", p, e, n)), ProjectID: project, RepoID: repo, ParentEpicID: epic, Title: "Planned task", WorktreeState: workspace.WorkItemUnbound})
			}
		}
	}
	d.WorkItems = &viewTestItems{registry: registry}
	d.Projects = &viewTestProjects{projects: projects}
	d.WorkGit, d.IntegrationGit, d.Repos, d.ProjectLocks = nil, nil, nil, nil
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := ReadWorkflowStatus(d, WorkflowStatusOptions{})
		if err != nil || out.Counts.Total != 300 {
			b.Fatalf("%+v %v", out.Counts, err)
		}
	}
}

func TestWorkflowStatusReassessmentRetiresOldQuestion(t *testing.T) {
	s := workflowFixture()
	s.Registry.TaskProblemRevisions = []workspace.TaskProblemReference{{TaskID: "task", Revision: 1, ManifestSHA256: "problem"}}
	s.Registry.TaskSolutionSelections = []workspace.TaskSelectionReference{{TaskID: "task", Ordinal: 1, ID: "selected", ManifestSHA256: "selection"}}
	s.Registry.TaskSpecRevisions = []workspace.TaskSpecReference{{TaskID: "task", SpecID: "execution", Revision: 1, ManifestSHA256: "spec"}}
	s.Registry.TaskSpecAssessments = []workspace.TaskAssessmentReference{
		{TaskID: "task", SpecID: "execution", SpecRevision: 1, Ordinal: 1, ID: "assessment-old", ManifestSHA256: "assessment-old-sha"},
		{TaskID: "task", SpecID: "execution", SpecRevision: 1, Ordinal: 2, ID: "assessment-new", ManifestSHA256: "assessment-new-sha"},
	}
	run := workflowReportedRun("needs_input")
	run.Delivery.Basis = &workspace.TaskSpecBasis{TaskID: "task", Problem: workspace.TaskRevisionRef{Revision: 1, ManifestSHA256: "problem"}, SpecID: "execution", Spec: workspace.TaskRevisionRef{Revision: 1, ManifestSHA256: "spec"}, Selection: workspace.TaskDecisionRef{ID: "selected", ManifestSHA256: "selection"}, Assessment: workspace.TaskDecisionRef{ID: "assessment-old", ManifestSHA256: "assessment-old-sha"}}
	out := workflowBuild(t, s, false, run)
	if len(out.Items) != 1 || out.Items[0].Category != "follow_up" {
		t.Fatalf("new assessment leaves stale run question current: category=%s actions=%+v", out.Items[0].Category, out.Items[0].NextActions)
	}
}

func TestWorkflowStatusChangedCurrentPreparationRequiresFollowUp(t *testing.T) {
	s := workflowFixture()
	run := workflowReportedRun("working")
	run.Delivery = nil // A preserved legacy workflow run has no delivery contract.
	queues := workspace.RegisteredTaskQueues{Freshness: "fresh", Queues: []workspace.RegisteredTaskQueue{{QueueID: "queue", Revision: 1, Target: workspace.QueueTarget{ProjectID: "project", RepoID: "repo", EpicID: "epic"}, Freshness: "fresh", Current: &workspace.RegisteredQueueCurrent{QueueCurrent: workspace.QueueCurrent{TaskID: "task", State: "prepared", Reasons: []workspace.QueueReason{{Code: "task_goal_problem_changed", Message: "Current Problem differs from this preserved preparation."}}}, Freshness: "fresh"}}}}
	out, err := BuildWorkflowStatus(s, WorkflowStatusOptions{}, nil, taskrun.Inventory{Runs: []taskrun.InventoryRun{run}}, queues)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Items) != 1 || out.Items[0].Category != "follow_up" {
		t.Fatalf("known current preparation drift has no follow-up: category=%s actions=%+v queue=%+v", out.Items[0].Category, out.Items[0].NextActions, out.Items[0].Queues)
	}
}

func TestWorkflowStatusDeliveryRequiresRegisteredExactBasis(t *testing.T) {
	for _, change := range []string{"none", "problem", "selection", "assessment", "spec", "spec_digest", "task_owner"} {
		t.Run(change, func(t *testing.T) {
			s := workflowFixture()
			s.Registry.TaskProblemRevisions = []workspace.TaskProblemReference{{TaskID: "task", Revision: 1, ManifestSHA256: "problem"}}
			s.Registry.TaskSolutionSelections = []workspace.TaskSelectionReference{{TaskID: "task", Ordinal: 1, ID: "selected", ManifestSHA256: "selection"}}
			s.Registry.TaskSpecRevisions = []workspace.TaskSpecReference{{TaskID: "task", SpecID: "execution", Revision: 1, ManifestSHA256: "spec"}}
			s.Registry.TaskSpecAssessments = []workspace.TaskAssessmentReference{{TaskID: "task", SpecID: "execution", SpecRevision: 1, Ordinal: 1, ID: "assessed", ManifestSHA256: "assessment"}}
			run := workflowReportedRun("needs_input")
			run.Delivery.Basis = &workspace.TaskSpecBasis{TaskID: "task", Problem: workspace.TaskRevisionRef{Revision: 1, ManifestSHA256: "problem"}, SpecID: "execution", Spec: workspace.TaskRevisionRef{Revision: 1, ManifestSHA256: "spec"}, Selection: workspace.TaskDecisionRef{ID: "selected", ManifestSHA256: "selection"}, Assessment: workspace.TaskDecisionRef{ID: "assessed", ManifestSHA256: "assessment"}}
			switch change {
			case "problem":
				s.Registry.TaskProblemRevisions = nil
			case "selection":
				s.Registry.TaskSolutionSelections = nil
			case "assessment":
				s.Registry.TaskSpecAssessments = nil
			case "spec":
				s.Registry.TaskSpecRevisions = nil
			case "spec_digest":
				s.Registry.TaskSpecRevisions[0].ManifestSHA256 = "other"
			case "task_owner":
				run.Delivery.Basis.TaskID = "other"
			}
			out := workflowBuild(t, s, false, run)
			want := "follow_up"
			if change == "none" {
				want = "needs_you"
			}
			if len(out.Items) != 1 || out.Items[0].Category != want {
				t.Fatalf("category=%s actions=%+v", out.Items[0].Category, out.Items[0].NextActions)
			}
		})
	}
}
