package workspaceview

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func workflowFixture() *Snapshot {
	s := progressFixture()
	s.ObservedAtUTC = "2026-10-06T12:00:00Z"
	s.Projects.Projects[0].RepoIDs = []workspace.RepoID{"repo", "second-repo"}
	s.Projects.Projects = append(s.Projects.Projects, workspace.ProjectRecord{ID: "other", Name: "Other", RepoIDs: []workspace.RepoID{"other-repo"}}, workspace.ProjectRecord{ID: "empty", Name: "Empty"})
	s.Projects.Repos = append(s.Projects.Repos, workspace.RepoRecord{ID: "second-repo"}, workspace.RepoRecord{ID: "other-repo"})
	s.Registry.Epics[0].RepoBindings = []taskwork.EpicRepoBinding{{RepoID: "repo"}, {RepoID: "second-repo"}}
	s.Registry.Epics = append(s.Registry.Epics, taskwork.EpicRecord{ID: "empty-epic", ProjectID: "other", Title: "Empty Epic", RepoBindings: []taskwork.EpicRepoBinding{{RepoID: "other-repo"}}})
	s.Registry.EpicBaseVersions = []taskwork.EpicBaseVersion{{ProjectID: "project", RepoID: "repo", EpicID: "epic", Revision: 1}, {ProjectID: "project", RepoID: "second-repo", EpicID: "epic", Revision: 1}, {ProjectID: "other", RepoID: "other-repo", EpicID: "empty-epic", Revision: 1}}
	return s
}

func workflowBuild(t *testing.T, s *Snapshot, all bool) WorkflowStatus {
	t.Helper()
	out, err := BuildWorkflowStatus(s, WorkflowStatusOptions{IncludeAll: all}, nil, taskwork.RegisteredTaskQueues{Freshness: "fresh"})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestWorkflowStatusWorkspaceFiltersAndEmptyContainers(t *testing.T) {
	s := workflowFixture()
	out := workflowBuild(t, s, false)
	if len(out.Projects) != 3 || len(out.Epics) != 2 || out.Counts.Total != 1 || out.Counts.Hidden.Backlog != 1 || len(out.Items) != 0 {
		t.Fatalf("whole workspace or empty containers lost: %+v", out)
	}
	project, repo, epic := workspace.ProjectID("other"), workspace.RepoID("repo"), taskwork.EpicID("empty-epic")
	out, err := BuildWorkflowStatus(s, WorkflowStatusOptions{Filters: taskwork.TaskListFilters{ProjectID: &project, RepoID: &repo, EpicID: &epic}}, nil, taskwork.RegisteredTaskQueues{})
	if err != nil || len(out.Items) != 0 || len(out.Projects) != 0 || len(out.Epics) != 0 {
		t.Fatalf("known disjoint AND filters: %+v %v", out, err)
	}
	epic = "missing"
	if _, err = BuildWorkflowStatus(s, WorkflowStatusOptions{Filters: taskwork.TaskListFilters{ProjectID: &project, RepoID: &repo, EpicID: &epic}}, nil, taskwork.RegisteredTaskQueues{}); err == nil {
		t.Fatal("unknown filter hidden by empty intersection")
	}
	project, repo, epic = "project", "second-repo", "epic"
	out, err = BuildWorkflowStatus(s, WorkflowStatusOptions{Filters: taskwork.TaskListFilters{ProjectID: &project, RepoID: &repo, EpicID: &epic}}, nil, taskwork.RegisteredTaskQueues{})
	if err != nil || len(out.Items) != 0 || len(out.Projects) != 1 || len(out.Epics) != 1 || len(out.Epics[0].Base) != 1 {
		t.Fatalf("empty registered target disappeared: %+v %v", out, err)
	}
}

func TestWorkflowStatusQueueReadinessAndChangedSourceBasis(t *testing.T) {
	for _, state := range []string{"ready", "blocked", "unknown"} {
		for _, freshness := range []string{"fresh", "unknown"} {
			s := workflowFixture()
			s.Freshness = freshness
			queues := taskwork.RegisteredTaskQueues{Freshness: "fresh", Queues: []taskwork.RegisteredTaskQueue{{QueueID: "queue", Revision: 1, Target: taskwork.QueueTarget{ProjectID: "project", RepoID: "repo", EpicID: "epic"}, Freshness: "fresh", Pending: []taskwork.RegisteredQueuePending{{QueuePending: taskwork.QueuePending{Rank: 1, TaskID: "task", State: state, SpecID: stringValue("goal")}, Spec: &taskwork.TaskRevisionRef{Revision: 1, ManifestSHA256: "digest"}, Freshness: "fresh"}}}}}
			out, err := BuildWorkflowStatus(s, WorkflowStatusOptions{}, nil, queues)
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
	items.registry.IntegrationAuthorities = []taskwork.IntegrationAuthority{fixtureAuthority()}
	items.registry.IntegrationResults = []taskwork.IntegrationResult{{ID: "integration", AuthorityID: "authority", Outcome: "exact_effect"}}
	if _, err := taskwork.SetWorkItemLifecycle(withViewLocks(d), taskwork.WorkItemLifecycleInput{SubjectKind: "epic", SubjectID: "epic", State: taskwork.LifecycleParked, ActorClaim: "fixture owner"}); err != nil {
		t.Fatal(err)
	}
	s, err := LoadSnapshot(d)
	if err != nil {
		t.Fatal(err)
	}
	out := workflowBuild(t, s, false)
	if out.Counts.Hidden.Inactive != 1 || out.Counts.Hidden.Completed != 0 || out.Counts.Hidden.Backlog != 0 || len(out.Items) != 0 {
		t.Fatalf("hidden precedence or independent diagnostics lost: %+v", out)
	}
	all := workflowBuild(t, s, true)
	if all.Counts.Hidden != (WorkflowStatusHidden{}) || len(all.Items) != 1 || all.Items[0].Category != "inactive" || all.Items[0].TaskLifecycle != taskwork.LifecycleActive || all.Items[0].Progress.State != "integrated" {
		t.Fatalf("--all changed lifecycle or native history: %+v", all)
	}
	for _, a := range all.Items[0].NextActions {
		if a.Current {
			t.Fatalf("historical action marked current: %+v", a)
		}
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
		if out.Workspace.Root != root || out.Counts.Hidden.Backlog != 1 || items.reads != n+1 || projects.reads != 0 || items.fullReads != 0 || items.locks != 0 {
			t.Fatalf("not one read per registry or accidental live/lock dependency: %+v items=%+v projects=%+v", out, items, projects)
		}
	}
	if !reflect.DeepEqual(before, viewTree(t, root)) {
		t.Fatal("status altered workspace sources")
	}
}

func TestWorkflowStatusJSONNullsCollectionsAndDeterminism(t *testing.T) {
	s := workflowFixture()
	a, b := workflowBuild(t, s, true), workflowBuild(t, s, true)
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
	d := taskwork.SystemDependencies()
	d.Files = viewTestFiles{d.Files, root}
	if _, err = workspace.Init(d.Dependencies); err != nil {
		b.Fatal(err)
	}
	registry := taskwork.WorkItemRegistry{}
	projects := workspace.ProjectSnapshot{}
	for p := 0; p < 4; p++ {
		project, repo := workspace.ProjectID(fmt.Sprintf("project-%d", p)), workspace.RepoID(fmt.Sprintf("repo-%d", p))
		projects.Projects = append(projects.Projects, workspace.ProjectRecord{ID: project, RepoIDs: []workspace.RepoID{repo}})
		projects.Repos = append(projects.Repos, workspace.RepoRecord{ID: repo})
		for e := 0; e < 3; e++ {
			epic := taskwork.EpicID(fmt.Sprintf("epic-%d-%d", p, e))
			registry.Epics = append(registry.Epics, taskwork.EpicRecord{ID: epic, ProjectID: project})
			for n := 0; n < 25; n++ {
				registry.Tasks = append(registry.Tasks, taskwork.TaskRecord{ID: taskwork.TaskID(fmt.Sprintf("task-%d-%d-%d", p, e, n)), ProjectID: project, RepoID: repo, ParentEpicID: epic, Title: "Planned task", WorktreeState: taskwork.WorkItemUnbound})
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
