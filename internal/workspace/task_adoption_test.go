package workspace

import (
	"bytes"
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

func TestAdoptionPredecessorFollowsCausalReleaseNotRandomID(t *testing.T) {
	// Registry order is random operation-ID order, opposite to causal order.
	a := WorktreeOperationRecord{ID: "wop_ffffffffffffffffffffffffffffffff", TaskID: "task", Adoption: &TaskAdoptionIntent{}, Release: &TaskAdoptionRelease{Reason: "first release"}}
	b := WorktreeOperationRecord{ID: "wop_00000000000000000000000000000000", TaskID: "task", Adoption: &TaskAdoptionIntent{Plan: TaskAdoptionPlan{Predecessor: queueDigest(a.Release)}}, Release: &TaskAdoptionRelease{Reason: "second release"}}
	r := WorkItemRegistry{WorktreeOperations: []WorktreeOperationRecord{b, a}}
	if got := adoptionPredecessor(r, "task"); got != queueDigest(b.Release) {
		t.Fatalf("fresh same-bytes adoption would revive an old generation: predecessor=%s want latest causal release=%s", got, queueDigest(b.Release))
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

func TestAdoptionAcceptsExplicitSourceDocumentWithoutExecutionProcess(t *testing.T) {
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
	if _, e = RecordTaskSpec(f.dependencies, draft); e != nil {
		t.Fatal(e)
	}
	if got := runLocalGit(t, in.Path, "rev-parse", "HEAD"); got != a.Plan.Source.OID {
		t.Fatal("document publication changed adopted source")
	}

}

func TestAdoptionRegistryDeliveredUpgradePreservesRecords(t *testing.T) {
	directory, pin := deliveredRegistryReader(t)
	f, _ := epicCreateFixture(t)
	root := filepath.Dir(f.wrapper)
	beforeRegistry, e := os.ReadFile(workItemsPath(root))
	if e != nil {
		t.Fatal(e)
	}
	beforeFixture := queueStateBytes(t, f)
	p, e := UpgradeAdoptionRegistry(f.dependencies, directory, pin, "", false)
	if e != nil {
		t.Fatal(e)
	}
	if queueStateBytes(t, f) != beforeFixture {
		t.Fatal("upgrade preview changed the fixture")
	}
	if _, e = UpgradeAdoptionRegistry(f.dependencies, directory, pin, p.Confirmation, true); e != nil {
		t.Fatal(e)
	}
	publishedFixture := queueStateBytes(t, f)
	again, e := UpgradeAdoptionRegistry(f.dependencies, directory, pin, p.Confirmation, true)
	if e != nil || again.State != "already_upgraded" {
		t.Fatalf("upgrade lost-reply recovery %+v %v", again, e)
	}
	if queueStateBytes(t, f) != publishedFixture {
		t.Fatal("upgrade retry changed the published fixture")
	}
	backup, e := os.ReadFile(p.Backup)
	if e != nil || !bytes.Equal(backup, beforeRegistry) {
		t.Fatal("upgrade lost byte-exact backup")
	}
	r := mustQueueRegistry(t, f)
	if r.FormatVersion != 7 || r.AdoptionBoundary == nil {
		t.Fatal("no native adoption format boundary")
	}
	r.FormatVersion = 6
	r.AdoptionBoundary = nil
	restored, e := encodeWorkItemRegistry(r)
	if e != nil || !bytes.Equal(restored, beforeRegistry) {
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

// r001-f01: every legacy create/recovery caller shares the same source-mode
// fence. A lost adoption publication reply cannot turn Git cleanliness into
// fresh ownership-transfer evidence.
func TestAdoptionLegacyCreateCannotFinalizeReservation(t *testing.T) {
	for _, entrypoint := range []string{"public-create", "shared-preparation-create"} {
		t.Run(entrypoint, func(t *testing.T) {
			f, in := adoptionFixture(t)
			preview, e := AdoptTaskWorktree(f.dependencies, in, false, "")
			if e != nil {
				t.Fatal(e)
			}
			publications := 0
			f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(stage string) error {
				if stage == "directory-sync" {
					publications++
					if publications == 1 {
						return errInjected
					}
				}
				return nil
			}}}
			if _, e = AdoptTaskWorktree(f.dependencies, in, true, preview.Confirmation); e == nil {
				t.Fatal("lost adoption reply was not injected")
			}
			f.dependencies.WorkItems = newSystemWorkItemStore()
			r := mustQueueRegistry(t, f)
			op, _ := findOperation(r, in.TaskID)
			if op == nil || op.State != "creating" || op.Adoption == nil {
				t.Fatal("fixture did not retain reserved adoption")
			}
			operationID := op.ID
			writerCalls, reserveCalls := 0, 0
			writerState := "present"
			f.dependencies.SourceWriters = func(string) (SourceWriterObservation, error) {
				writerCalls++
				return SourceWriterObservation{State: writerState, Method: "modeled review regression", Limitation: "Injected observations; no actual OS/provider acceptance."}, nil
			}
			before := queueStateBytes(t, f)
			input := TaskWorktreeCreateInput{TaskID: in.TaskID, Branch: "useful", Path: in.Path, ExpectedParentOID: f.oid}
			var result TaskWorktreeMutationResult
			if entrypoint == "public-create" {
				result, e = CreateTaskWorktree(f.dependencies, input)
			} else {
				root := filepath.Dir(f.wrapper)
				e = f.dependencies.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
					return f.dependencies.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
						registry, e := session.Snapshot()
						if e != nil {
							return e
						}
						reserve := func(*WorkItemRegistry, WorktreeOperationRecord) error { reserveCalls++; return nil }
						return createTaskWorktreeLocked(f.dependencies, root, projects, session, registry, input, in.Path, true, &result, reserve)
					})
				})
			}
			if e == nil {
				t.Fatalf("legacy %s finalized reserved adoption despite known writer: outcome=%s state=%s", entrypoint, result.Outcome, result.Operation.State)
			}
			if !strings.Contains(e.Error(), "task_adoption_recovery_required") {
				t.Fatalf("wrong source-mode fence: %v", e)
			}
			if writerCalls != 0 || reserveCalls != 0 || queueStateBytes(t, f) != before {
				t.Fatal("rejected legacy recovery observed transfer or changed native/Git state")
			}
			if _, e = AdoptTaskWorktree(f.dependencies, in, true, preview.Confirmation); e == nil || !strings.Contains(e.Error(), "adoption_writer_present") {
				t.Fatalf("native adoption recovery ignored fresh writer evidence: %v", e)
			}
			writerState = "absent"
			recovered, e := AdoptTaskWorktree(f.dependencies, in, true, preview.Confirmation)
			if e != nil {
				t.Fatal(e)
			}
			if recovered.Phase != "adopted" || recovered.Operation.ID != operationID || len(mustQueueRegistry(t, f).WorktreeOperations) != 1 {
				t.Fatal("native adoption did not recover the exact reserved generation")
			}
		})
	}
}
