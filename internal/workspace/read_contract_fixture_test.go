package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

// These expectation types deliberately do not use the read-contract serializers.
// Their values come from the seed inputs, never from candidate CLI output.
type readContractExpectedProject struct {
	ProjectID string `json:"project_id"`
	Name      string `json:"name"`
	Wrapper   string `json:"wrapper"`
	RepoCount int    `json:"repo_count"`
}

type readContractExpectedRepo struct {
	RepoID       string `json:"repo_id"`
	Locator      string `json:"locator"`
	GitCommonDir string `json:"git_common_dir"`
}

type readContractExpectedTask struct {
	TaskID        string  `json:"task_id"`
	ProjectID     string  `json:"project_id"`
	RepoID        string  `json:"repo_id"`
	ParentEpicID  string  `json:"parent_epic_id"`
	WorktreeState string  `json:"worktree_state"`
	Title         *string `json:"title"`
	TitleSource   string  `json:"title_source"`
	TitleStatus   string  `json:"title_status"`
}

type readContractFixtureManifest struct {
	Kind                string                                `json:"kind"`
	Root                string                                `json:"root"`
	EmptyRoot           string                                `json:"empty_root"`
	OutsideRoot         string                                `json:"outside_root"`
	CorruptProjectsRoot string                                `json:"corrupt_projects_root"`
	CorruptTasksRoot    string                                `json:"corrupt_tasks_root"`
	Home                string                                `json:"home"`
	Projects            []readContractExpectedProject         `json:"projects"`
	Repositories        map[string][]readContractExpectedRepo `json:"repositories"`
	Tasks               []readContractExpectedTask            `json:"tasks"`
	ProtectedPaths      []string                              `json:"protected_paths"`
	Provenance          []string                              `json:"provenance"`
}

// TestReadContractAcceptanceFixture is opt-in because the read-phase harness needs
// a preserved fixture after go test exits. The destination must not already exist:
//
//	PLY_READ_CONTRACT_FIXTURE_DIR=/absolute/new/directory go test ./internal/workspace -run '^TestReadContractAcceptanceFixture$' -count=1 -v
//
// Run test/read_contract_acceptance.py with the resulting manifest.json and two
// independently built binaries. queue-healthy/manifest.json and
// unavailable-v4/manifest.json additionally cover actual queue registrations.
// The old base rejects all format-4 reads when any problem manifest is unavailable,
// so its success-only text comparison requires the primary format-3 fixture.
// The fixtures contain no user data and are not removed.
func TestReadContractAcceptanceFixture(t *testing.T) {
	base := os.Getenv("PLY_READ_CONTRACT_FIXTURE_DIR")
	if base == "" {
		t.Skip("set PLY_READ_CONTRACT_FIXTURE_DIR to preserve an acceptance fixture")
	}
	if !filepath.IsAbs(base) || filepath.Clean(base) != base {
		t.Fatal("PLY_READ_CONTRACT_FIXTURE_DIR must be an absolute clean path")
	}
	if err := os.MkdirAll(filepath.Dir(base), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(base, 0o755); err != nil {
		t.Fatalf("fixture destination must be new (existing fixtures are preserved): %v", err)
	}
	base = physicalPath(t, base)
	readContractBuildFixture(t, base, false, true)
	for _, scenario := range []struct {
		name        string
		unavailable bool
	}{{"queue-healthy", false}, {"unavailable-v4", true}} {
		dir := filepath.Join(base, scenario.name)
		if err := os.Mkdir(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		readContractBuildFixture(t, dir, true, scenario.unavailable)
	}
}

func readContractBuildFixture(t *testing.T, base string, withQueue, unavailableContent bool) {
	t.Helper()
	f := readContractFixtureManifest{
		Kind: "PlyReadContractFixture@1", Root: filepath.Join(base, "registered"),
		EmptyRoot: filepath.Join(base, "empty"), OutsideRoot: filepath.Join(base, "outside"),
		CorruptProjectsRoot: filepath.Join(base, "corrupt-projects"),
		CorruptTasksRoot:    filepath.Join(base, "corrupt-tasks"), Home: filepath.Join(base, "home"),
		Repositories: map[string][]readContractExpectedRepo{}, ProtectedPaths: []string{},
		Provenance: []string{
			"Expected Project, Repository and Task values are seeded in read_contract_fixture_test.go independently of CLI output.",
			"Workspace markers use Init; project registrations use the validated ProjectStore.Add path.",
			"Three Epics use AdoptEpic against local committed Git repositories; current Tasks use CreateTask.",
			"The legacy Task and reconciliation state use WorkItemStoreSession.Publish, including registry and content-closure validation.",
			"The revised problem uses RecordTaskProblem; creating uses the existing noEffectGit fixture; ready uses CreateTaskWorktree.",
			"All registry/content validators pass before any deliberate manifest faults; original faulted manifests are preserved beside the fixture.",
			"All existing checkouts and Git directories are inside the watched registered root; the missing checkout is never created.",
		},
	}
	if withQueue {
		f.Provenance = append(f.Provenance, "SetTaskQueue publishes two pending entries in format 4; the other six Tasks remain outside the queue.")
	} else {
		f.Provenance = append(f.Provenance, "Format 3 preserves old successful text behavior with unavailable titles; the separate queue fixtures cover in/out queue membership.")
	}
	if unavailableContent {
		f.Provenance = append(f.Provenance, "The corrupt and missing problem manifests intentionally make their two Task titles unavailable.")
	}
	for _, root := range []string{f.Root, f.EmptyRoot, f.OutsideRoot, f.CorruptProjectsRoot, f.CorruptTasksRoot, f.Home} {
		if err := os.Mkdir(root, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, root := range []string{f.Root, f.EmptyRoot, f.CorruptProjectsRoot, f.CorruptTasksRoot} {
		if _, err := Init(readContractDependencies(root)); err != nil {
			t.Fatal(err)
		}
	}
	d := readContractDependencies(f.Root)
	alphaWrapper := filepath.Join(f.Root, "Ålpha, project")
	zetaWrapper := filepath.Join(f.Root, "Zeta & project")
	alphaRepo, alphaOID := readContractRepository(t, alphaWrapper, "alpha-main")
	zetaRepo, zetaOID := readContractRepository(t, zetaWrapper, "zeta-main")
	missingRepo := RepoRecord{
		ID: "alpha-missing", Locator: filepath.Join(alphaWrapper, "missing checkout, ø"),
		GitCommonDir: filepath.Join(alphaWrapper, "missing checkout, ø", ".git"),
	}
	projectSeeds := []struct {
		project ProjectRecord
		repos   []RepoRecord
	}{
		{ProjectRecord{ID: "alpha", Name: "Ålpha: \"research\", one", Wrapper: alphaWrapper, RepoIDs: []RepoID{"alpha-main", "alpha-missing"}}, []RepoRecord{alphaRepo, missingRepo}},
		{ProjectRecord{ID: "zeta", Name: "Zeta & 日本語", Wrapper: zetaWrapper, RepoIDs: []RepoID{"zeta-main"}}, []RepoRecord{zetaRepo}},
	}
	// Reverse insertion order ensures expectations do not inherit write order.
	for i := len(projectSeeds) - 1; i >= 0; i-- {
		seed := projectSeeds[i]
		if _, _, _, err := d.Projects.Add(f.Root, seed.project, seed.repos); err != nil {
			t.Fatal(err)
		}
	}
	for _, seed := range projectSeeds {
		p := seed.project
		f.Projects = append(f.Projects, readContractExpectedProject{string(p.ID), p.Name, p.Wrapper, len(p.RepoIDs)})
		for _, r := range seed.repos {
			f.Repositories[string(p.ID)] = append(f.Repositories[string(p.ID)], readContractExpectedRepo{string(r.ID), r.Locator, r.GitCommonDir})
		}
	}
	for _, seed := range []struct {
		id      string
		project string
		repo    RepoRecord
		oid     string
	}{
		{"epic-zeta", "zeta", zetaRepo, zetaOID},
		{"epic-alpha-later", "alpha", alphaRepo, alphaOID},
		{"epic-alpha", "alpha", alphaRepo, alphaOID},
	} {
		path := filepath.Join(filepath.Dir(seed.repo.Locator), seed.id)
		runLocalGit(t, seed.repo.Locator, "worktree", "add", "-b", seed.id, path, seed.oid)
		in, err := ParseEpicAdoptInput(seed.id, "Fixture "+seed.id, seed.project, string(seed.repo.ID), path, "refs/heads/"+seed.id, seed.oid)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = AdoptEpic(d, in); err != nil {
			t.Fatal(err)
		}
	}
	f.Tasks = []readContractExpectedTask{
		{"task-a-current", "alpha", "alpha-main", "epic-alpha", "unbound", stringPointer("Revised: løse \"read\", safely"), "problem_revision", "available"},
		{"task-b-legacy", "alpha", "alpha-main", "epic-alpha", "unbound", stringPointer("Legacy: registrert tittel, æøå"), "registration", "available"},
		{"task-c-creating", "alpha", "alpha-main", "epic-alpha", "creating", stringPointer("Creating a worktree"), "problem_revision", "available"},
		{"task-d-reconcile", "alpha", "alpha-main", "epic-alpha", "reconciliation_required", stringPointer("Reconcile preserved intent"), "problem_revision", "available"},
		{"task-e-ready", "zeta", "zeta-main", "epic-zeta", "worktree_ready", stringPointer("Ready worktree: 日本語"), "problem_revision", "available"},
		{"task-f-corrupt", "zeta", "zeta-main", "epic-zeta", "unbound", nil, "problem_revision", "unavailable"},
		{"task-g-sentinel", "alpha", "alpha-main", "epic-alpha-later", "unbound", stringPointer("Problem content unavailable"), "problem_revision", "available"},
		{"task-h-missing", "zeta", "zeta-main", "epic-zeta", "unbound", nil, "problem_revision", "unavailable"},
	}
	if !unavailableContent {
		for _, i := range []int{5, 7} {
			f.Tasks[i].Title = stringPointer("Registered title before unavailable content")
			f.Tasks[i].TitleStatus = "available"
		}
	}
	for i := len(f.Tasks) - 1; i >= 0; i-- {
		seed := f.Tasks[i]
		title := "Registered title before unavailable content"
		if seed.Title != nil {
			title = *seed.Title
		}
		if seed.TaskID == "task-a-current" {
			title = "Original registration title"
		}
		if seed.TitleSource == "registration" {
			readContractPublish(t, d, f.Root, func(r *WorkItemRegistry) {
				r.Tasks = append(r.Tasks, TaskRecord{ID: TaskID(seed.TaskID), Title: title, Description: "Legacy fixture registration", ParentEpicID: EpicID(seed.ParentEpicID), ProjectID: ProjectID(seed.ProjectID), RepoID: RepoID(seed.RepoID), GitCommonDir: alphaRepo.GitCommonDir, WorktreeState: WorkItemUnbound})
				r.TaskSpecPolicies = append(r.TaskSpecPolicies, TaskSpecPolicy{TaskID: TaskID(seed.TaskID), Mode: "legacy", LegacyResultIDs: []TaskResultID{}})
			})
			continue
		}
		in, err := ParseTaskCreateInput(seed.TaskID, title, "Isolated RL-01 read fixture", seed.ParentEpicID, seed.ProjectID, seed.RepoID)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = CreateTask(d, in); err != nil {
			t.Fatal(err)
		}
	}
	readContractReviseProblem(t, d, f.Root, "task-a-current", *f.Tasks[0].Title)
	for _, seed := range []struct {
		id      string
		wrapper string
		oid     string
	}{
		{"task-c-creating", alphaWrapper, alphaOID},
		{"task-d-reconcile", alphaWrapper, alphaOID},
		{"task-e-ready", zetaWrapper, zetaOID},
	} {
		in, err := ParseTaskWorktreeCreateInput(seed.id, seed.id, filepath.Join(seed.wrapper, seed.id), seed.oid)
		if err != nil {
			t.Fatal(err)
		}
		deps := d
		if seed.id != "task-e-ready" {
			deps.WorkGit = noEffectGit{WorkItemGit: d.WorkGit}
		}
		_, err = CreateTaskWorktree(deps, in)
		if seed.id == "task-e-ready" && err != nil || seed.id != "task-e-ready" && !hasWorkErrorClass(err, ErrorWorkGitEffect) {
			t.Fatalf("seed %s worktree state: %v", seed.id, err)
		}
	}
	readContractPublish(t, d, f.Root, func(r *WorkItemRegistry) {
		task, _ := findTask(*r, "task-d-reconcile")
		op, _ := findOperation(*r, task.ID)
		task.WorktreeState = WorkItemReconciliationRequired
		op.State = "reconciliation_required"
		op.LastObservation = emptyCreateObservation("partial_or_unknown")
		op.LastObservation.Reasons = []string{"task_worktree_reconciliation_required"}
	})
	if withQueue {
		r := readContractRegistry(t, d, f.Root)
		queue := WorkspaceTaskQueueDraft{
			Kind: "WorkspaceTaskQueueDraft@1", SchemaVersion: 1, PublicationKey: "rl01/queue",
			ProjectID: "alpha", RepoID: "alpha-main", EpicID: "epic-alpha", ExpectedRevision: 0,
			Entries:         []QueueEntry{{TaskID: "task-b-legacy"}, {TaskID: "task-a-current"}},
			HumanDecision:   QueueHumanDecision{ActorClaim: "isolated fixture", DecidedAtUTC: "2026-10-05T00:00:00Z", Source: "explicit_human_instruction", Statement: "Synthetic test priority; not a human product verdict."},
			RegistryUpgrade: queueUpgrade(r),
		}
		queueBytes, err := contentCanonical(queue)
		if err != nil {
			t.Fatal(err)
		}
		queuePath := filepath.Join(f.Root, "queue-draft.json")
		readContractWrite(t, queuePath, queueBytes, 0o600)
		if _, err = SetTaskQueue(d, queuePath); err != nil {
			t.Fatal(err)
		}
	}
	r := readContractRegistry(t, d, f.Root)
	if err := validateTaskContentClosure(d.TaskContent, f.Root, r, false); err != nil {
		t.Fatal(err)
	}
	if err := validateQueueClosure(d, f.Root, r); err != nil {
		t.Fatal(err)
	}
	if withQueue && len(r.TaskQueueEvents) != 1 || !withQueue && len(r.TaskQueueEvents) != 0 || len(r.Tasks) != 8 || len(r.Epics) != 3 {
		t.Fatal("fixture lost queue, Task or Epic seed coverage")
	}
	if head := taskContentState(r, "task-a-current").ProblemHead; head == nil || head.Revision != 2 {
		t.Fatal("current title must come from the second problem revision")
	}
	if _, err := os.Lstat(missingRepo.Locator); !os.IsNotExist(err) {
		t.Fatalf("missing checkout unexpectedly exists: %v", err)
	}
	// Preserve original immutable bytes, then make one head corrupt and one absent.
	if unavailableContent {
		preserved := filepath.Join(base, "preserved-manifests")
		if err := os.Mkdir(preserved, 0o700); err != nil {
			t.Fatal(err)
		}
		for _, id := range []TaskID{"task-f-corrupt", "task-h-missing"} {
			head := taskContentState(r, id).ProblemHead
			path := taskContentPath(f.Root, "manifests", head.ManifestSHA256)
			if err := os.Rename(path, filepath.Join(preserved, string(id)+".json")); err != nil {
				t.Fatal(err)
			}
			if id == "task-f-corrupt" {
				readContractWrite(t, path, []byte("{\"deliberately_corrupt\":true}\n"), 0o400)
			}
		}
	}
	readContractValidateSeeds(t, d, f)
	readContractWrite(t, filepath.Join(f.CorruptProjectsRoot, MarkerDirectory, ProjectsFile), []byte("format_version: 999\nprojects: []\nrepos: []\n"), 0o644)
	readContractWrite(t, workItemsPath(f.CorruptTasksRoot), []byte("format_version: 999\nepics: []\ntasks: []\nworktree_operations: []\n"), 0o644)
	if _, _, err := d.Projects.Snapshot(f.CorruptProjectsRoot); err == nil {
		t.Fatal("corrupt project registry must be rejected")
	}
	if _, err := d.WorkItems.Snapshot(f.CorruptTasksRoot); err == nil {
		t.Fatal("corrupt Task registry must be rejected")
	}
	manifest, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(base, "manifest.json")
	readContractWrite(t, path, append(manifest, '\n'), 0o644)
	t.Logf("preserved RL-01 acceptance fixture: %s", path)
}

func readContractDependencies(root string) Dependencies {
	d := SystemDependencies()
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: root}
	return d
}

func readContractRepository(t *testing.T, wrapper, id string) (RepoRecord, string) {
	t.Helper()
	path := filepath.Join(wrapper, "main")
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, path, "init", "-b", "main")
	runLocalGit(t, path, "config", "user.email", "read-contract@example.invalid")
	runLocalGit(t, path, "config", "user.name", "Read contract fixture")
	readContractWrite(t, filepath.Join(path, "README.md"), []byte("# Isolated RL-01 fixture\n"), 0o644)
	runLocalGit(t, path, "add", "README.md")
	runLocalGit(t, path, "commit", "-m", "Seed isolated read-contract fixture")
	return RepoRecord{ID: RepoID(id), Locator: path, GitCommonDir: filepath.Join(path, ".git")}, runLocalGit(t, path, "rev-parse", "HEAD")
}

func readContractWrite(t *testing.T, path string, contents []byte, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, contents, mode); err != nil {
		t.Fatal(err)
	}
}

func readContractRegistry(t *testing.T, d Dependencies, root string) WorkItemRegistry {
	t.Helper()
	r, err := d.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if err = validateWorkItemRegistry(r); err != nil {
		t.Fatal(err)
	}
	return r
}

func readContractPublish(t *testing.T, d Dependencies, root string, mutate func(*WorkItemRegistry)) {
	t.Helper()
	err := d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
		r, err := session.Snapshot()
		if err != nil {
			return err
		}
		mutate(&r)
		sortWorkRegistry(&r)
		return session.Publish(r)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func readContractReviseProblem(t *testing.T, d Dependencies, root string, id TaskID, title string) {
	t.Helper()
	head := taskContentState(readContractRegistry(t, d, root), id).ProblemHead
	doc := filepath.Join(root, "revised-problem.md")
	contents := []byte("# " + title + "\n\nIndependently seeded revised fixture problem.\n")
	readContractWrite(t, doc, contents, 0o600)
	fields := map[string]canonicaljson.Value{
		"publication_key": "rl01/revised-problem", "task_id": string(id), "registry_upgrade": nil,
		"recorder":          contentObject(map[string]canonicaljson.Value{"actor_claim": "isolated fixture", "control_surface": "canonical Go fixture", "recorded_at_utc": "2026-10-05T00:00:00Z"}),
		"expected_previous": contentRefValue(head), "origin": contentObject(map[string]canonicaljson.Value{"kind": "authored"}),
		"title": title, "summary": "Independently seeded revised fixture problem.", "problem_document_id": "problem",
		"documents": []canonicaljson.Value{contentObject(map[string]canonicaljson.Value{"id": "problem", "source": contentObject(map[string]canonicaljson.Value{"kind": "file", "locator": doc, "sha256": digestTaskBytes(contents), "size_bytes": int64(len(contents)), "media_type": "text/markdown", "git_provenance": nil})})},
		"sources":   []canonicaljson.Value{}, "claims": []canonicaljson.Value{}, "deadline": nil, "change_reason": "Exercise the current title after revision.",
	}
	raw, err := canonicaljson.Marshal(contentEnvelope("WorkspaceTaskProblemDraft@1", fields))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "revised-problem-draft.json")
	readContractWrite(t, path, raw, 0o600)
	if _, err = RecordTaskProblem(d, TaskContentInput{TaskID: id, File: path}); err != nil {
		t.Fatal(err)
	}
}

func readContractValidateSeeds(t *testing.T, d Dependencies, f readContractFixtureManifest) {
	t.Helper()
	projects, repos, err := d.Projects.Snapshot(f.Root)
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != len(f.Projects) || len(repos) != 3 {
		t.Fatal("registered project/repository count differs from seed")
	}
	for i, p := range projects {
		if got := (readContractExpectedProject{string(p.ID), p.Name, p.Wrapper, len(p.RepoIDs)}); got != f.Projects[i] {
			t.Fatalf("project seed differs: %#v", got)
		}
	}
	// Validate registry bytes independently of format-4's deliberately strict
	// Snapshot content-closure check. The read contract owns per-Task title faults.
	raw, err := os.ReadFile(workItemsPath(f.Root))
	if err != nil {
		t.Fatal(err)
	}
	r, err := decodeWorkItemRegistry(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Tasks) != len(f.Tasks) {
		t.Fatal("registered Task count differs from seed")
	}
	for i, task := range r.Tasks {
		want := f.Tasks[i]
		got := readContractExpectedTask{string(task.ID), string(task.ProjectID), string(task.RepoID), string(task.ParentEpicID), string(task.WorktreeState), &task.Title, "registration", "available"}
		if head := taskContentState(r, task.ID).ProblemHead; head != nil {
			got.TitleSource = "problem_revision"
			manifest, err := readRegisteredTaskManifest(d.TaskContent, f.Root, r, head.ManifestSHA256)
			if err != nil {
				got.Title = nil
				got.TitleStatus = "unavailable"
			} else {
				title := contentString(contentFields(manifest), "title")
				got.Title = &title
			}
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Task seed differs for %s: got %s; want %s", task.ID, readContractJSON(got), readContractJSON(want))
		}
	}
}

func readContractJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%#v", v)
	}
	return string(b)
}
