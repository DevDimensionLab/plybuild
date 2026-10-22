package workspace

import (
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func adoptionFixture(t *testing.T) (workItemJourneyFixture, TaskAdoptionInput) {
	t.Helper()
	f := newWorkItemJourneyFixture(t)
	f.dependencies.Files = projectCwdFileSystem{FileSystem: f.dependencies.Files, cwd: filepath.Join(f.wrapper, "epic")}
	a := automaticAgreementFixture()
	a.ProjectID, a.RepoID, a.EpicID = "ply", "ply", "epic"
	a.TargetRef, a.TargetWorktree = "refs/heads/epic", filepath.Join(f.wrapper, "epic")
	g := goalFixture(t, f, "task", "adopt-goal", a)
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &g}})
	seedAdoptionRegistry(t, f)
	path := filepath.Join(f.wrapper, "useful")
	runLocalGit(t, filepath.Join(f.wrapper, "main"), "worktree", "add", "-b", "useful", path, f.oid)
	if e := os.WriteFile(filepath.Join(path, "useful.txt"), []byte("useful inherited change\n"), 0600); e != nil {
		t.Fatal(e)
	}
	runLocalGit(t, path, "add", "useful.txt")
	runLocalGit(t, path, "commit", "-m", "useful inherited work")
	f.dependencies.SourceWriters = func(string) (SourceWriterObservation, error) {
		return SourceWriterObservation{State: "absent", Method: "modeled_idle", Limitation: "Injected fixture; no actual process or provider acceptance."}, nil
	}
	return f, TaskAdoptionInput{TaskID: "task", Target: queueTargetFixture(), Path: path, Ref: "refs/heads/useful", ExpectedOID: runLocalGit(t, path, "rev-parse", "HEAD")}
}

func TestAdoptionExistingSourceJourney(t *testing.T) {
	f, in := adoptionFixture(t)
	before := runLocalGit(t, filepath.Join(f.wrapper, "main"), "show-ref")
	queueBefore := queueDigest(mustQueueRegistry(t, f).TaskQueueEvents)
	p, e := AdoptTaskWorktree(f.dependencies, in, false, "")
	if e != nil {
		t.Fatalf("clean useful existing S != B must preview for the fresh Goal: %v", e)
	}
	a, e := AdoptTaskWorktree(f.dependencies, in, true, p.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	if a.Phase != "adopted" || a.Operation == nil {
		t.Fatalf("not adopted: %+v", a)
	}
	if got := runLocalGit(t, filepath.Join(f.wrapper, "main"), "show-ref"); got != before {
		t.Fatal("adoption mutated Git refs")
	}
	if queueDigest(mustQueueRegistry(t, f).TaskQueueEvents) != queueBefore {
		t.Fatal("adoption mutated queue")
	}
	p2, e := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: in.Target, Next: true})
	if e != nil {
		t.Fatal(e)
	}
	if p2.WorktreePath != in.Path || p2.ParentOID != f.oid {
		t.Fatal("execution lost S/B or existing source")
	}
	prepared, e := PrepareTaskGoalExecution(f.dependencies, p2.Input, p2.Confirmation, goalHumanFixture())
	if e != nil {
		t.Fatal(e)
	}
	if prepared.Preparation.Preparation.Outcome.ObservedOID != in.ExpectedOID {
		t.Fatal("preparation did not start at adopted S")
	}
	root := filepath.Dir(f.wrapper)
	e = WithTaskSpecSnapshot(f.dependencies, root, func(s *TaskSpecSession) error { _, e := s.CurrentDeliveryTaskSpec(in.TaskID); return e })
	if e != nil {
		t.Fatal(e)
	}
	// A clean sibling C of S (same tree, parent B) must not inherit S's mandate.
	sibling := runLocalGit(t, in.Path, "commit-tree", a.Plan.Source.Tree, "-p", f.oid, "-m", "sibling candidate outside adopted ancestry")
	runLocalGit(t, in.Path, "update-ref", in.Ref, sibling, in.ExpectedOID)
	e = WithTaskSpecSnapshot(f.dependencies, root, func(s *TaskSpecSession) error {
		_, e := s.ValidateDeliveryCandidateTarget(in.Path, in.Ref, a.Plan.Source.GitCommonDir, &prepared.Basis)
		return e
	})
	if e == nil || !strings.Contains(e.Error(), "not descended from the frozen execution start") {
		t.Fatalf("clean sibling candidate must fail adopted S ancestry: %v", e)
	}
}

func TestAdoptionPredecessorFollowsCausalReleaseNotRandomID(t *testing.T) {
	// Registry order is random operation-ID order, opposite to causal order.
	a := WorktreeOperationRecord{ID: "wop_ffffffffffffffffffffffffffffffff", TaskID: "task", Adoption: &TaskAdoptionIntent{}, Release: &TaskAdoptionRelease{Reason: "first release"}}
	b := WorktreeOperationRecord{ID: "wop_00000000000000000000000000000000", TaskID: "task", Adoption: &TaskAdoptionIntent{Plan: TaskAdoptionPlan{Predecessor: queueDigest(a.Release)}}, Release: &TaskAdoptionRelease{Reason: "second release"}}
	r := WorkItemRegistry{WorktreeOperations: []WorktreeOperationRecord{b, a}}
	if got := adoptionPredecessor(r, "task"); got != queueDigest(b.Release) {
		t.Fatalf("fresh same-bytes adoption would revive an old generation: predecessor=%s want latest causal release=%s", got, queueDigest(b.Release))
	}
}

func TestAdoptionReleaseFreshGenerationAndExecutionFence(t *testing.T) {
	f, in := adoptionFixture(t)
	// Each generation consumes 16 worktree-ID and 16 operation-ID bytes.
	// Deliberately reverse/permute chronological IDs with unchanged source bytes.
	f.dependencies.WorkIDs = cryptoWorkItemIDSource{Reader: strings.NewReader(strings.Repeat("\xff", 32) + strings.Repeat("\x11", 32) + strings.Repeat("\x88", 32) + strings.Repeat("\x44", 32))}
	var old []TaskAdoptionReadback
	for cycle := 0; cycle < 3; cycle++ {
		p, e := AdoptTaskWorktree(f.dependencies, in, false, "")
		if e != nil {
			t.Fatal(e)
		}
		for _, a := range old {
			if p.Confirmation == a.Confirmation {
				t.Fatal("fresh same-bytes preview revived older generation")
			}
		}
		a, e := AdoptTaskWorktree(f.dependencies, in, true, p.Confirmation)
		if e != nil {
			t.Fatal(e)
		}
		ex, e := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: in.Target, Next: true})
		if e != nil {
			t.Fatal(e)
		}
		before := queueDigest(mustQueueRegistry(t, f).TaskQueueEvents)
		// Release must not depend on source availability, dirty state, base, or writer absence.
		f.dependencies.SourceWriters = func(string) (SourceWriterObservation, error) {
			t.Fatal("release inspected writers")
			return SourceWriterObservation{}, nil
		}
		r, e := ReleaseTaskAdoption(f.dependencies, in.TaskID, a.Operation.ID, "fresh source", false, "")
		if e != nil {
			t.Fatal(e)
		}
		r, e = ReleaseTaskAdoption(f.dependencies, in.TaskID, a.Operation.ID, "fresh source", true, r.Confirmation)
		if e != nil {
			t.Fatal(e)
		}
		q, e := ListTaskQueue(f.dependencies, in.Target, false)
		if e != nil {
			t.Fatal(e)
		}
		if len(q.Pending) != 1 || q.Pending[0].State != "adoption_released" {
			t.Fatal("released Goal can fall through to creation")
		}
		if r.Phase != "released" || !r.MetadataOnly || r.WriterAbsence != "not_required" {
			t.Fatalf("untruthful release %+v", r)
		}
		if queueDigest(mustQueueRegistry(t, f).TaskQueueEvents) != before {
			t.Fatal("release changed queue")
		}
		if _, e = PrepareTaskGoalExecution(f.dependencies, ex.Input, ex.Confirmation, goalHumanFixture()); e == nil {
			t.Fatal("released execution preview materialized")
		}
		again, e := AdoptTaskWorktree(f.dependencies, in, true, a.Confirmation)
		if e != nil || again.Phase != "released" {
			t.Fatalf("old adopt revived: %+v %v", again, e)
		}
		f.dependencies.SourceWriters = func(string) (SourceWriterObservation, error) {
			return SourceWriterObservation{State: "absent", Method: "modeled", Limitation: "modeled observations only"}, nil
		}
		old = append(old, a)
	}
	p, e := AdoptTaskWorktree(f.dependencies, in, false, "")
	if e != nil {
		t.Fatal(e)
	}
	a, e := AdoptTaskWorktree(f.dependencies, in, true, p.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	ex, e := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: in.Target, Next: true})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = PrepareTaskGoalExecution(f.dependencies, ex.Input, ex.Confirmation, goalHumanFixture()); e != nil {
		t.Fatal(e)
	}
	if _, e = ReleaseTaskAdoption(f.dependencies, in.TaskID, a.Operation.ID, "too late", false, ""); e == nil {
		t.Fatal("release passed execution intent fence")
	}
	for _, o := range old {
		if _, e = ReleaseTaskAdoption(f.dependencies, in.TaskID, o.Operation.ID, "replacement", true, "sha256:bad"); e == nil {
			t.Fatal("stale release affected replacement")
		}
	}
}

func TestAdoptionExecutionMaterializationFaultFencesRelease(t *testing.T) {
	f, in := adoptionFixture(t)
	p, e := AdoptTaskWorktree(f.dependencies, in, false, "")
	if e != nil {
		t.Fatal(e)
	}
	a, e := AdoptTaskWorktree(f.dependencies, in, true, p.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	ex, e := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: in.Target, Next: true})
	if e != nil {
		t.Fatal(e)
	}
	f.dependencies.TaskContent.fault = func(stage string) error {
		if stage == "before-registry-replace" {
			return errInjected
		}
		return nil
	}
	if _, e = PrepareTaskGoalExecution(f.dependencies, ex.Input, ex.Confirmation, goalHumanFixture()); e == nil {
		t.Fatal("fault did not interrupt materialization")
	}
	f.dependencies.TaskContent.fault = nil
	if _, e = ReleaseTaskAdoption(f.dependencies, in.TaskID, a.Operation.ID, "possible execution intent", false, ""); e == nil {
		t.Fatal("release lost a possible execution materialization after publication began")
	}
	retry, e := PrepareTaskGoalExecution(f.dependencies, ex.Input, ex.Confirmation, goalHumanFixture())
	if e != nil {
		t.Fatalf("same-request retry stranded durable execution intent: %v", e)
	}
	if retry.Preparation.Preparation.Outcome.ObservedOID != in.ExpectedOID || len(mustQueueRegistry(t, f).TaskPreparations) != 1 {
		t.Fatal("retry duplicated preparation or changed source")
	}
	again, e := PrepareTaskGoalExecution(f.dependencies, ex.Input, ex.Confirmation, goalHumanFixture())
	if e != nil || again.Preparation.Preparation.ID != retry.Preparation.Preparation.ID {
		t.Fatal("lost reply did not recover same preparation")
	}

}

func TestAdoptionSourceAdvanceReleaseAndCompetingExecution(t *testing.T) {
	f, in := adoptionFixture(t)
	p, e := AdoptTaskWorktree(f.dependencies, in, false, "")
	if e != nil {
		t.Fatal(e)
	}
	a, e := AdoptTaskWorktree(f.dependencies, in, true, p.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	old := a.Operation.ID
	if e = os.WriteFile(filepath.Join(in.Path, "third.txt"), []byte("third useful change\n"), 0600); e != nil {
		t.Fatal(e)
	}
	runLocalGit(t, in.Path, "add", "third.txt")
	runLocalGit(t, in.Path, "commit", "-m", "source advanced to S2")
	in.ExpectedOID = runLocalGit(t, in.Path, "rev-parse", "HEAD")
	if _, e = PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: in.Target, Next: true}); e == nil {
		t.Fatal("changed source accepted as old S")
	}
	before := runLocalGit(t, in.Path, "show-ref")
	r, e := ReleaseTaskWorktreeAdoption(f.dependencies, old, "source advanced", false, "")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = ReleaseTaskWorktreeAdoption(f.dependencies, old, "source advanced", true, r.Confirmation); e != nil {
		t.Fatal(e)
	}
	p, e = AdoptTaskWorktree(f.dependencies, in, false, "")
	if e != nil {
		t.Fatal(e)
	}
	a, e = AdoptTaskWorktree(f.dependencies, in, true, p.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	if a.Operation.ID == old || a.Plan.Source.OID != in.ExpectedOID || before != runLocalGit(t, in.Path, "show-ref") {
		t.Fatal("S2 adoption restored bytes or revived old claim")
	}
	ex, e := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: in.Target, Next: true})
	if e != nil {
		t.Fatal(e)
	}
	r, e = ReleaseTaskWorktreeAdoption(f.dependencies, a.Operation.ID, "race", false, "")
	if e != nil {
		t.Fatal(e)
	}
	start := make(chan struct{})
	results := make(chan error, 2)
	go func() {
		<-start
		_, e := PrepareTaskGoalExecution(f.dependencies, ex.Input, ex.Confirmation, goalHumanFixture())
		results <- e
	}()
	go func() {
		<-start
		_, e := ReleaseTaskWorktreeAdoption(f.dependencies, a.Operation.ID, "race", true, r.Confirmation)
		results <- e
	}()
	close(start)
	wins := 0
	for i := 0; i < 2; i++ {
		if <-results == nil {
			wins++
		}
	}
	if wins != 1 {
		t.Fatalf("execution/release CAS must have one winner; got %d", wins)
	}
	current := mustQueueRegistry(t, f)
	op, _ := adoptionOperation(current, a.Operation.ID)
	if op.Release != nil && (op.Execution != nil || len(current.TaskPreparations) > 0) {
		t.Fatal("released generation has execution effects")
	}
}

func TestAdoptionPublicationFaultRecovery(t *testing.T) {
	for _, stage := range []string{"intent-lost-reply", "binding-lost-reply", "release-before-write", "release-lost-reply"} {
		t.Run(stage, func(t *testing.T) {
			f, in := adoptionFixture(t)
			p, e := AdoptTaskWorktree(f.dependencies, in, false, "")
			if e != nil {
				t.Fatal(e)
			}
			count := 0
			fault := func(point string) error {
				if point == "directory-sync" {
					count++
					if stage == "intent-lost-reply" && count == 1 || stage == "binding-lost-reply" && count == 2 || stage == "release-lost-reply" && count == 1 {
						return errInjected
					}
				}
				if stage == "release-before-write" && point == "temp-open" {
					return errInjected
				}
				return nil
			}
			if stage == "intent-lost-reply" || stage == "binding-lost-reply" {
				f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: fault}}
				if _, e = AdoptTaskWorktree(f.dependencies, in, true, p.Confirmation); e == nil {
					t.Fatal("publication fault not reported")
				}
				f.dependencies.WorkItems = newSystemWorkItemStore()
			}
			a, e := AdoptTaskWorktree(f.dependencies, in, true, p.Confirmation)
			if e != nil {
				t.Fatal(e)
			}
			if a.Phase != "adopted" || len(mustQueueRegistry(t, f).WorktreeOperations) != 1 {
				t.Fatal("retry duplicated or lost adoption")
			}
			r, e := ReleaseTaskWorktreeAdoption(f.dependencies, a.Operation.ID, "same release", false, "")
			if e != nil {
				t.Fatal(e)
			}
			if stage == "release-before-write" || stage == "release-lost-reply" {
				f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: fault}}
				if _, e = ReleaseTaskWorktreeAdoption(f.dependencies, a.Operation.ID, "same release", true, r.Confirmation); e == nil {
					t.Fatal("release fault not reported")
				}
				f.dependencies.WorkItems = newSystemWorkItemStore()
			}
			released, e := ReleaseTaskWorktreeAdoption(f.dependencies, a.Operation.ID, "same release", true, r.Confirmation)
			if e != nil {
				t.Fatal(e)
			}
			if released.Phase != "released" || released.Operation.ID != a.Operation.ID {
				t.Fatal("release recovery changed identity")
			}
			if _, e = AdoptTaskWorktree(f.dependencies, in, true, p.Confirmation); e != nil {
				t.Fatal(e)
			}
			if op, _ := findOperation(mustQueueRegistry(t, f), in.TaskID); op != nil {
				t.Fatal("released request revived after retry")
			}
		})
	}
}

func seedAdoptionRegistry(t *testing.T, f workItemJourneyFixture) {
	t.Helper()
	root := filepath.Dir(f.wrapper)
	r := mustQueueRegistry(t, f)
	boundary := func(from int) *RegistryFormatBoundary {
		return &RegistryFormatBoundary{FromVersion: from, FromSHA256: "sha256:" + strings.Repeat("a", 64), ManifestSHA256: "sha256:" + strings.Repeat("b", 64), ReaderVersion: "v1.0.0", ReaderCommit: strings.Repeat("c", 40), ReaderSHA256: "sha256:" + strings.Repeat("d", 64)}
	}
	r.FormatVersion = 7
	r.FormatBoundary = boundary(4)
	r.EpicCreateBoundary = boundary(5)
	r.AdoptionBoundary = boundary(6)
	r.EpicCreateOperations = []EpicCreateOperation{}
	b, e := encodeWorkItemRegistry(r)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(workItemsPath(root), b, 0600); e != nil {
		t.Fatal(e)
	}
}

func TestAdoptionMetadataReleaseDoesNotReadGoalBytes(t *testing.T) {
	f, in := adoptionFixture(t)
	p, e := AdoptTaskWorktree(f.dependencies, in, false, "")
	if e != nil {
		t.Fatal(e)
	}
	a, e := AdoptTaskWorktree(f.dependencies, in, true, p.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	reads := 0
	fault := func(stage string) error {
		if stage == "managed-read" {
			reads++
			return errInjected
		}
		return nil
	}
	f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: fault}}
	f.dependencies.TaskContent.fault = fault
	r, e := ReleaseTaskWorktreeAdoption(f.dependencies, a.Operation.ID, "Goal bytes unavailable before start", false, "")
	if e != nil {
		t.Fatalf("metadata-only unstarted release must not depend on Goal freshness: %v", e)
	}
	if _, e = ReleaseTaskWorktreeAdoption(f.dependencies, a.Operation.ID, "Goal bytes unavailable before start", true, r.Confirmation); e != nil {
		t.Fatal(e)
	}
	if reads != 0 {
		t.Fatal("metadata release read Goal bytes")
	}
}

type adoptionInProgressGit struct {
	TaskIntegrationGit
	source string
}

func (g adoptionInProgressGit) ObserveIntegrationWorktree(path, ref string) (IntegrationWorktreeObservation, error) {
	o, e := g.TaskIntegrationGit.ObserveIntegrationWorktree(path, ref)
	if path == g.source {
		o.InProgress = []string{"MERGE_HEAD"}
	}
	return o, e
}
func TestAdoptionCleanSourceWithGitOperationIsRejected(t *testing.T) {
	f, in := adoptionFixture(t)
	f.dependencies.IntegrationGit = adoptionInProgressGit{f.dependencies.IntegrationGit, in.Path}
	if _, e := AdoptTaskWorktree(f.dependencies, in, false, ""); e == nil {
		t.Fatal("clean status hid an in-progress source Git operation")
	}
}

func TestAdoptionRejectsLegacyExecutionPublicationBeforeFence(t *testing.T) {
	f, in := adoptionFixture(t)
	p, e := AdoptTaskWorktree(f.dependencies, in, false, "")
	if e != nil {
		t.Fatal(e)
	}
	a, e := AdoptTaskWorktree(f.dependencies, in, true, p.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	r := mustQueueRegistry(t, f)
	draft := contentFixtureSpec(t, f, "legacy-on-adopted", taskContentState(r, in.TaskID).SpecHead)
	raw, e := os.ReadFile(draft.File)
	if e != nil {
		t.Fatal(e)
	}
	v, e := canonicaljson.DecodeStrict(raw)
	if e != nil {
		t.Fatal(e)
	}
	m := contentFields(v)
	basis := contentFields(m["implementation_basis"])
	basis["start_oid"] = a.Plan.Source.OID
	basis["start_tree"] = a.Plan.Source.Tree
	m["implementation_basis"] = contentObject(basis)
	m["removed_requirement_ids"] = []canonicaljson.Value{"behavior"}
	raw, e = canonicaljson.Marshal(contentObject(m))
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(draft.File, raw, 0600); e != nil {
		t.Fatal(e)
	}
	before := queueStateBytes(t, f)
	if _, e = RecordTaskSpec(f.dependencies, draft); e == nil {
		t.Fatal("native RecordTaskSpec materialized an unfenced legacy execution on adopted S")
	}
	if !strings.Contains(e.Error(), "adoption_requires_goal_execute") {
		t.Fatalf("legacy publication failed for an unrelated guard: %v", e)
	}
	if queueStateBytes(t, f) != before {
		t.Fatal("legacy rejection wrote native metadata")
	}
}

func TestAdoptionChangedSourceAfterExecutionIntentRemainsBlocked(t *testing.T) {
	f, in := adoptionFixture(t)
	p, e := AdoptTaskWorktree(f.dependencies, in, false, "")
	if e != nil {
		t.Fatal(e)
	}
	a, e := AdoptTaskWorktree(f.dependencies, in, true, p.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	ex, e := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: in.Target, Next: true})
	if e != nil {
		t.Fatal(e)
	}
	f.dependencies.TaskContent.fault = func(stage string) error {
		if stage == "before-registry-replace" {
			return errInjected
		}
		return nil
	}
	if _, e = PrepareTaskGoalExecution(f.dependencies, ex.Input, ex.Confirmation, goalHumanFixture()); e == nil {
		t.Fatal("fault did not stop materialization")
	}
	f.dependencies.TaskContent.fault = nil
	runLocalGit(t, in.Path, "commit", "--allow-empty", "-m", "external source change after intent")
	if _, e = PrepareTaskGoalExecution(f.dependencies, ex.Input, ex.Confirmation, goalHumanFixture()); e == nil || !strings.Contains(e.Error(), "adoption_source_changed") {
		t.Fatalf("changed frozen source must block same request: %v", e)
	}
	if _, e = ReleaseTaskWorktreeAdoption(f.dependencies, a.Operation.ID, "changed after execution intent", false, ""); e == nil || !strings.Contains(e.Error(), "adoption_execution_exists") {
		t.Fatalf("release must not cancel a possible execution: %v", e)
	}
	r := mustQueueRegistry(t, f)
	if len(r.TaskPreparations) != 0 {
		t.Fatal("blocked intent created preparation")
	}
	o, _ := adoptionOperation(r, a.Operation.ID)
	if o.Execution == nil || o.Release != nil {
		t.Fatal("intent/refusal history lost")
	}
	t.Log("Supported boundary: same-request no-effect publication retries recover on exact S/B. After intent, changed S or B blocks preparation and metadata release; this feature adds no generic cancellation or restoration.")
}

func TestAdoptionRegistryDeliveredUpgradePreservesRecords(t *testing.T) {
	directory, pin := deliveredRegistryReader(t)
	f, _ := epicCreateFixture(t)
	root := filepath.Dir(f.wrapper)
	p, e := UpgradeAdoptionRegistry(f.dependencies, directory, pin, "", false)
	if e != nil {
		t.Fatal(e)
	}
	before := queueStateBytes(t, f)
	if _, e = UpgradeAdoptionRegistry(f.dependencies, directory, pin, p.Confirmation, true); e != nil {
		t.Fatal(e)
	}
	again, e := UpgradeAdoptionRegistry(f.dependencies, directory, pin, p.Confirmation, true)
	if e != nil || again.State != "already_upgraded" {
		t.Fatalf("upgrade lost-reply recovery %+v %v", again, e)
	}
	backup, e := os.ReadFile(p.Backup)
	if e != nil || string(backup) != before {
		t.Fatal("upgrade lost byte-exact backup")
	}
	r := mustQueueRegistry(t, f)
	if r.FormatVersion != 7 || r.AdoptionBoundary == nil {
		t.Fatal("no native adoption format boundary")
	}
	r.FormatVersion = 6
	r.AdoptionBoundary = nil
	restored, e := encodeWorkItemRegistry(r)
	if e != nil || string(restored) != before {
		t.Fatal("upgrade rewrote existing semantic records")
	}
	if e = validateRegistryFormatTransition(mustQueueRegistry(t, f), r); e == nil {
		t.Fatal("adoption registry downgrade accepted")
	}
	if _, e = os.Stat(filepath.Join(root, ".ply", "task-content")); e != nil {
		t.Fatal(e)
	}
}

func TestAdoptionMetadataReleasePreservesDirtyUnavailableSourceAndScope(t *testing.T) {
	f, in := adoptionFixture(t)
	p, e := AdoptTaskWorktree(f.dependencies, in, false, "")
	if e != nil {
		t.Fatal(e)
	}
	a, e := AdoptTaskWorktree(f.dependencies, in, true, p.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(in.Path, "external-untracked"), []byte("external writer bytes\n"), 0600); e != nil {
		t.Fatal(e)
	}
	before := runLocalGit(t, in.Path, "status", "--porcelain")
	refs := runLocalGit(t, in.Path, "show-ref")
	f.dependencies.WorkGit = nil
	f.dependencies.IntegrationGit = nil
	f.dependencies.SourceWriters = nil
	r, e := ReleaseTaskWorktreeAdoption(f.dependencies, a.Operation.ID, "metadata only, source observation unavailable", false, "")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = ReleaseTaskWorktreeAdoption(f.dependencies, a.Operation.ID, "metadata only, source observation unavailable", true, r.Confirmation); e != nil {
		t.Fatal(e)
	}
	if before != runLocalGit(t, in.Path, "status", "--porcelain") || refs != runLocalGit(t, in.Path, "show-ref") {
		t.Fatal("release changed dirty source")
	}
	beforeRegistry := queueStateBytes(t, f)
	metadata := adoptionMetadataDependencies(f.dependencies, in.TaskID)
	e = metadata.WorkItems.WithLock(filepath.Dir(f.wrapper), func(s WorkItemStoreSession) error {
		r, e := s.Snapshot()
		if e != nil {
			return e
		}
		r.Tasks[0].Title = "unrelated mutation"
		return s.Publish(r)
	})
	if e == nil || queueStateBytes(t, f) != beforeRegistry {
		t.Fatal("metadata-only view allowed an unrelated write")
	}
}

type adoptionUnknownGit struct {
	WorkItemGit
	TaskContentGit
	source  string
	blocked *bool
}

func (g adoptionUnknownGit) ObserveWorktree(path string) (GitWorktreeObservation, error) {
	if path == g.source && *g.blocked {
		return GitWorktreeObservation{}, errInjected
	}
	return g.WorkItemGit.ObserveWorktree(path)
}
func TestAdoptionUnknownSourceObservationIsPreservedAndCanRelease(t *testing.T) {
	f, in := adoptionFixture(t)
	p, e := AdoptTaskWorktree(f.dependencies, in, false, "")
	if e != nil {
		t.Fatal(e)
	}
	blocked := false
	original := f.dependencies.WorkGit
	f.dependencies.WorkGit = adoptionUnknownGit{original, original.(TaskContentGit), in.Path, &blocked}
	f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(stage string) error {
		if stage == "directory-sync" {
			blocked = true
		}
		return nil
	}}}
	out, e := AdoptTaskWorktree(f.dependencies, in, true, p.Confirmation)
	if e == nil || out.Phase != "unknown" {
		t.Fatalf("unknown source became adoption success: %+v %v", out, e)
	}
	f.dependencies.WorkItems = newSystemWorkItemStore()
	r := mustQueueRegistry(t, f)
	if len(r.WorktreeOperations) != 1 || len(r.WorktreeOperations[0].AdoptionObservations) != 1 {
		t.Fatal("unknown source lost native observation history")
	}
	id := r.WorktreeOperations[0].ID
	release, e := ReleaseTaskWorktreeAdoption(f.dependencies, id, "unknown source but authoritative unstarted metadata", false, "")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = ReleaseTaskWorktreeAdoption(f.dependencies, id, "unknown source but authoritative unstarted metadata", true, release.Confirmation); e != nil {
		t.Fatal(e)
	}
	current := mustQueueRegistry(t, f).WorktreeOperations[0]
	if current.Release == nil || current.AdoptionObservations[0].Phase != "unknown" {
		t.Fatal("release erased unknown source history")
	}
}
