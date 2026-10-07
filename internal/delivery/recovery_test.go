package delivery

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type controlledPR struct {
	observation                 PRObservation
	pushes, creates, edits      int
	observationErr, errorOnEdit error
	ff                          bool
}

func (a *controlledPR) Observe(PRTarget) (PRObservation, error) {
	return a.observation, a.observationErr
}
func (a *controlledPR) CanFastForward(PRTarget, string) (bool, error) { return a.ff, nil }
func (a *controlledPR) Push(t PRTarget) error {
	a.pushes++
	a.observation.RemoteOID = t.OID
	return nil
}
func (a *controlledPR) Create(t PRTarget, title, body string) error {
	a.creates++
	p := testPR(t)
	p.Body = body + marker(t)
	a.observation.PullRequests = []PullRequest{p}
	return nil
}
func (a *controlledPR) ApplyMetadata(t PRTarget, p PullRequest, m MetadataChoice) error {
	a.edits++
	if a.errorOnEdit != nil {
		return a.errorOnEdit
	}
	for _, v := range m.Assignees {
		if !includes(p.Assignees, v) {
			p.Assignees = append(p.Assignees, v)
		}
	}
	for _, v := range m.Reviewers {
		if !includes(p.Reviewers, v) {
			p.Reviewers = append(p.Reviewers, v)
		}
	}
	a.observation.PullRequests = []PullRequest{p}
	return nil
}

// This fixture tests the effect journal below the native authority gate. The
// external-facing registration/execute tests in taskrun use actual native
// qualified candidates; this fixture deliberately cannot close a native run.
func journalReceipt(t *testing.T) (*Service, string, Receipt, effect, *controlledPR) {
	t.Helper()
	root, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	a := &controlledPR{observation: PRObservation{BaseExists: true}, ff: true}
	s := NewService(taskrun.Dependencies{}, Options{PR: a, Now: func() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }})
	target := testTarget()
	in := Registration{Kind: "ply.delivery.registration", SchemaVersion: 1, PublicationKey: "unit-effect-journal", TaskID: "task", TaskResultID: workspace.TaskResultID("trs_" + strings.Repeat("1", 32)), WorkflowRunID: "wfr_" + strings.Repeat("2", 64)}
	m := Manifest{Workspace: root, TaskID: "task", OwnerClaim: "test adapter", TaskResult: workspace.TaskResultRecord{ID: in.TaskResultID, SourceLocator: target.Worktree, SourceRef: target.SourceRef, ResultOID: target.OID, ResultTree: strings.Repeat("c", 40)}, Agreement: workspace.DeliveryAgreement{SchemaVersion: 1, Mode: workspace.DeliveryPullRequest, ProjectID: "project", RepoID: "repo", EpicID: "epic", SourceRef: target.SourceRef, TargetRef: target.BaseRef, Remote: target.Remote, GitHubRepository: target.Repository}, EffectKey: target.EffectKey}
	r := Receipt{Kind: "ply.delivery.receipt", SchemaVersion: 1, ID: "dlv_" + shortHash(in.PublicationKey), Registration: in, RegistrationSHA256: digest(in), Manifest: m, ManifestSHA256: digest(m), State: "delivering", Reasons: []string{}, Events: []Event{}, CreatedAtUTC: s.now(), AttemptID: "unit-attempt", Metadata: &MetadataChoice{Assignees: []string{"requested-owner"}, Reviewers: []string{"requested-reviewer"}, RemoveAssignees: []string{}, RemoveReviewers: []string{}, Revision: 1, State: "pending"}}
	s.event(&r, "delivery.registered", "registered", "unit journal fixture", "")
	if e = s.save(root, &r); e != nil {
		t.Fatal(e)
	}
	fx := effect{Kind: "ply.delivery.effect", SchemaVersion: 1, Key: m.EffectKey, OwnerDeliveryID: r.ID}
	if e = saveEffect(root, fx); e != nil {
		t.Fatal(e)
	}
	return s, root, r, fx, a
}

func TestInterruptedPushAndCreateAreObservedWithoutReplay(t *testing.T) {
	for _, point := range []string{"after_push_effect", "after_pr_effect", "after_metadata_effect"} {
		t.Run(point, func(t *testing.T) {
			s, root, r, fx, a := journalReceipt(t)
			s.options.Fault = func(at string) error {
				if at == point {
					return errors.New("simulated process loss")
				}
				return nil
			}
			if e := s.executePR(taskrun.Dependencies{}, root, &r, &fx); e == nil {
				t.Fatal("interruption was not reached")
			}
			stored, e := readReceipt(root, r.ID)
			if e != nil {
				t.Fatal(e)
			}
			preserved, e := readEffect(root, fx.Key)
			if e != nil {
				t.Fatal(e)
			}
			resumed := NewService(taskrun.Dependencies{}, Options{PR: a})
			if e = resumed.executePR(taskrun.Dependencies{}, root, &stored, &preserved); e != nil {
				t.Fatal(e)
			}
			if a.pushes != 1 || a.creates != 1 || a.edits != 1 {
				t.Fatalf("effects replayed: push=%d create=%d metadata=%d", a.pushes, a.creates, a.edits)
			}
			if stored.State == "delivered" || stored.NativeClosed {
				t.Fatal("adapter alone claimed native completion")
			}
			created := []Event{}
			for _, event := range stored.Events {
				if event.Kind == "delivery.pr_created" {
					created = append(created, event)
				}
			}
			if len(created) != 1 || created[0].URL != "https://github.com/Example/Project/pull/7" || stored.PR.CreationEventID != created[0].ID {
				t.Fatalf("creation event not preserved: %+v", created)
			}
			if e = resumed.executePR(taskrun.Dependencies{}, root, &stored, &preserved); e != nil {
				t.Fatal(e)
			}
			count := 0
			for _, event := range stored.Events {
				if event.Kind == "delivery.pr_created" {
					count++
				}
			}
			if count != 1 || a.creates != 1 || stored.PR.CreationEventID != created[0].ID {
				t.Fatal("retry duplicated logical creation event")
			}
		})
	}
}

func TestUnknownCreationNeverBlindlyRepeatsAndRecoveryKeepsEvent(t *testing.T) {
	s, root, r, fx, a := journalReceipt(t)
	s.options.Fault = func(point string) error {
		if point == "after_pr_effect" {
			return errors.New("lost process")
		}
		return nil
	}
	if e := s.executePR(taskrun.Dependencies{}, root, &r, &fx); e == nil {
		t.Fatal("missing interruption")
	}
	s.options.Fault = nil
	a.observationErr = errors.New("GitHub unavailable")
	if e := s.executePR(taskrun.Dependencies{}, root, &r, &fx); e != nil || r.State != "unknown_effect" || a.creates != 1 {
		t.Fatalf("unknown effect replayed: %+v %v", r, e)
	}
	for _, event := range r.Events {
		if event.Kind == "delivery.pr_created" {
			t.Fatal("creation claimed before observing its effect")
		}
	}
	a.observationErr = nil
	a.observation.PullRequests = nil
	if e := s.executePR(taskrun.Dependencies{}, root, &r, &fx); e != nil || r.State != "unknown_effect" || a.creates != 1 {
		t.Fatalf("absent PR replayed reserved create: %s %v", r.State, e)
	}
	p := testPR(prTarget(r))
	p.Body = marker(prTarget(r))
	a.observation.PullRequests = []PullRequest{p}
	if e := s.executePR(taskrun.Dependencies{}, root, &r, &fx); e != nil {
		t.Fatal(e)
	}
	if r.PR == nil || r.PR.CreationEventID == "" || a.creates != 1 {
		t.Fatal("recovered PR did not preserve creation identity")
	}
}

func TestExistingPRIsReusedAndMetadataFailureIsPartial(t *testing.T) {
	s, root, r, fx, a := journalReceipt(t)
	target := prTarget(r)
	p := testPR(target)
	p.Body = "an older PR"
	a.observation = PRObservation{RemoteOID: target.OID, BaseExists: true, PullRequests: []PullRequest{p}}
	a.errorOnEdit = errors.New("account cannot be assigned")
	if e := s.executePR(taskrun.Dependencies{}, root, &r, &fx); e != nil {
		t.Fatal(e)
	}
	if a.pushes != 0 || a.creates != 0 || r.State != "failed" || r.PR == nil || r.PR.Disposition != "reused" || r.PR.MetadataApplied || r.NativeClosed || r.Metadata.State != "failed" {
		t.Fatalf("partial metadata facts lost: %+v", r)
	}
	for _, event := range r.Events {
		if event.Kind == "delivery.pr_created" {
			t.Fatal("older PR falsely claimed as created")
		}
	}
	a.errorOnEdit = nil
	if e := s.executePR(taskrun.Dependencies{}, root, &r, &fx); e != nil {
		t.Fatal(e)
	}
	if a.creates != 0 || !r.PR.MetadataApplied || r.Metadata.State != "applied" || !includes(a.observation.PullRequests[0].Assignees, "existing-owner") || !includes(a.observation.PullRequests[0].Reviewers, "existing-reviewer") {
		t.Fatal("retry replaced PR, people or role choices")
	}
}

func TestConflictingPRAndNonFastForwardStopAllPublication(t *testing.T) {
	for _, scenario := range []string{"wrong-base", "wrong-head", "ambiguous", "force-push"} {
		t.Run(scenario, func(t *testing.T) {
			s, root, r, fx, a := journalReceipt(t)
			target := prTarget(r)
			p := testPR(target)
			a.observation.RemoteOID = target.OID
			a.observation.PullRequests = []PullRequest{p}
			switch scenario {
			case "wrong-base":
				a.observation.PullRequests[0].BaseRef = "refs/heads/other"
			case "wrong-head":
				a.observation.PullRequests[0].HeadOID = strings.Repeat("f", 40)
			case "ambiguous":
				a.observation.PullRequests = append(a.observation.PullRequests, p)
			case "force-push":
				a.observation.PullRequests = nil
				a.observation.RemoteOID = strings.Repeat("f", 40)
				a.ff = false
			}
			if e := s.executePR(taskrun.Dependencies{}, root, &r, &fx); e != nil {
				t.Fatal(e)
			}
			if r.State != "blocked" || a.pushes+a.creates+a.edits != 0 {
				t.Fatalf("conflicting target had effect: %s %+v", r.State, a)
			}
		})
	}
}

func TestCreationEventIsOwnedOnceAcrossDeliveryIDs(t *testing.T) {
	s, root, r, fx, _ := journalReceipt(t)
	other := r
	other.ID = "dlv_" + shortHash("second-key")
	other.Registration.PublicationKey = "second-key"
	other.RegistrationSHA256 = digest(other.Registration)
	other.Events = []Event{}
	if e := s.save(root, &other); e != nil {
		t.Fatal(e)
	}
	fx.CreateStarted = true
	p := testPR(prTarget(r))
	p.Body = marker(prTarget(r))
	if e := s.recordPR(root, &other, &fx, p); e != nil {
		t.Fatal(e)
	}
	owner, e := readReceipt(root, r.ID)
	if e != nil {
		t.Fatal(e)
	}
	if owner.PR == nil || owner.PR.CreationEventID == "" || other.PR.CreationEventID != owner.PR.CreationEventID {
		t.Fatal("logical PR identity changed across IDs")
	}
	for _, event := range other.Events {
		if event.Kind == "delivery.pr_created" {
			t.Fatal("alias manufactured a second creation event")
		}
	}
}

func TestStoreSerializesWritersAndRejectsChangedHistory(t *testing.T) {
	s, root, r, _, _ := journalReceipt(t)
	const writers = 12
	var wg sync.WaitGroup
	errs := make(chan error, writers)
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs <- withLock(root, func() error {
				v, e := readReceipt(root, r.ID)
				if e != nil {
					return e
				}
				s.event(&v, "test.observed", fmt.Sprint(i), "serialized observation", "")
				return s.save(root, &v)
			})
		}(i)
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	stored, e := readReceipt(root, r.ID)
	if e != nil {
		t.Fatal(e)
	}
	if len(stored.Events) != writers+1 {
		t.Fatalf("lost history: %d events", len(stored.Events))
	}
	stored.Events[0].Detail = "rewritten"
	if e = atomicJSON(recordPath(root, r.ID), stored, false); e != nil {
		t.Fatal(e)
	}
	if _, e = readReceipt(root, r.ID); e == nil {
		t.Fatal("rewritten event accepted")
	}
}

func TestStoreRefusesSymlinkRecords(t *testing.T) {
	s, root, r, _, _ := journalReceipt(t)
	path := recordPath(root, r.ID)
	copy := filepath.Join(root, "elsewhere.json")
	raw, e := os.ReadFile(path)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(copy, raw, 0600); e != nil {
		t.Fatal(e)
	}
	if e = os.Remove(path); e != nil {
		t.Fatal(e)
	}
	if e = os.Symlink(copy, path); e != nil {
		t.Fatal(e)
	}
	if _, e = readReceipt(root, r.ID); e == nil {
		t.Fatal("symlink record accepted")
	}
	if e = s.save(root, &r); e == nil {
		t.Fatal("symlink record overwritten")
	}
}

func TestStoreRejectsSymlinkDirectoryBeforeCreatingAnything(t *testing.T) {
	root := physicalTemp(t)
	elsewhere := physicalTemp(t)
	if e := os.Mkdir(filepath.Join(root, ".ply"), 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.Symlink(elsewhere, filepath.Join(root, ".ply", "deliveries")); e != nil {
		t.Fatal(e)
	}
	if e := withLock(root, func() error { return nil }); e == nil {
		t.Fatal("symlink delivery root accepted")
	}
	files, e := os.ReadDir(elsewhere)
	if e != nil {
		t.Fatal(e)
	}
	if len(files) != 0 {
		t.Fatalf("storage rejected symlink after creating outside directories: %v", files)
	}
}

func TestMetadataOnReusedDeliveryReportsOwnerBeforeAcceptingChange(t *testing.T) {
	s, root, owner, fx, _ := journalReceipt(t)
	if e := os.WriteFile(filepath.Join(root, ".ply", "workspace.yaml"), []byte("format_version: 1\nroot: "+root+"\n"), 0600); e != nil {
		t.Fatal(e)
	}
	s.dependencies = taskrun.SystemDependencies(workspace.SystemDependencies())
	owner.State = "delivered"
	owner.NativeClosed = true
	owner.PR = &PRReceipt{URL: "https://github.com/Example/Project/pull/7", MetadataApplied: true}
	if e := s.save(root, &owner); e != nil {
		t.Fatal(e)
	}
	fx.Complete = true
	fx.PR = owner.PR
	if e := saveEffect(root, fx); e != nil {
		t.Fatal(e)
	}
	alias := owner
	alias.ID = "dlv_" + shortHash("alias")
	alias.Registration.PublicationKey = "alias"
	alias.RegistrationSHA256 = digest(alias.Registration)
	alias.Events = []Event{}
	alias.ReusedDeliveryID = owner.ID
	if e := s.save(root, &alias); e != nil {
		t.Fatal(e)
	}
	file := filepath.Join(root, "adjustment.json")
	if e := os.WriteFile(file, []byte(`{"reviewers":["new-explicit-reviewer"]}`), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Metadata(root, alias.ID, file); e == nil || !strings.Contains(e.Error(), owner.ID) {
		t.Fatalf("alias adjustment accepted then discarded: %v", e)
	}
	stored, e := readReceipt(root, alias.ID)
	if e != nil {
		t.Fatal(e)
	}
	if stored.State != "delivered" || stored.Metadata.Revision != 1 {
		t.Fatal("rejected alias metadata changed the receipt")
	}
}

func TestReadinessPreservesUnobservedCreateIntentWithoutWriting(t *testing.T) {
	s, root, r, fx, a := journalReceipt(t)
	fx.CreateStarted = true
	fx.PushObserved = true
	a.observation.RemoteOID = r.Manifest.TaskResult.ResultOID
	if e := saveEffect(root, fx); e != nil {
		t.Fatal(e)
	}
	before, e := os.ReadFile(recordPath(root, r.ID))
	if e != nil {
		t.Fatal(e)
	}
	r.State = "ready"
	r.Reasons = []string{}
	if e = s.checkPRReadiness(root, &r); e != nil {
		t.Fatal(e)
	}
	if r.State != "unknown_effect" || len(r.Reasons) == 0 || !strings.Contains(r.NextAction, "observe") {
		t.Fatalf("lost unresolved intent: %s, %v, %s", r.State, r.Reasons, r.NextAction)
	}
	after, e := os.ReadFile(recordPath(root, r.ID))
	if e != nil {
		t.Fatal(e)
	}
	if string(before) != string(after) {
		t.Fatal("readiness wrote a delivery record")
	}
	if a.pushes+a.creates+a.edits != 0 {
		t.Fatal("readiness performed publication")
	}
}

func TestCompletedAliasPreservesOwnAcceptanceAndOriginalPRObservation(t *testing.T) {
	s, root, owner, fx, _ := journalReceipt(t)
	owner.State = "delivered"
	owner.NativeClosed = true
	owner.Metadata.State = "applied"
	owner.HumanQA = &workspace.TaskHumanQARecord{ID: "owner-qa"}
	fx.Complete = true
	fx.PR = &PRReceipt{URL: "https://github.com/Example/Project/pull/7", Disposition: "created", CreationEventID: "original-created-event", MetadataApplied: true}
	owner.PR = fx.PR
	alias := owner
	alias.ID = "dlv_" + shortHash("another-key")
	alias.Registration.PublicationKey = "another-key"
	alias.RegistrationSHA256 = digest(alias.Registration)
	alias.Events = []Event{}
	alias.HumanQA = &workspace.TaskHumanQARecord{ID: "own-qa"}
	if e := s.reuseCompleted(root, &alias, &fx, owner); e != nil {
		t.Fatal(e)
	}
	if alias.HumanQA.ID != "own-qa" || alias.PR.Disposition != "reused" || alias.PR.CreationEventID != "original-created-event" || fx.PR.Disposition != "created" {
		t.Fatalf("alias rewrote acceptance or original creation observation: alias=%+v effect=%+v", alias.PR, fx.PR)
	}
	third := alias
	third.ID = "dlv_" + shortHash("third-key")
	third.Registration.PublicationKey = "third-key"
	third.RegistrationSHA256 = digest(third.Registration)
	third.Events = []Event{}
	if e := s.reuseCompleted(root, &third, &fx, owner); e != nil {
		t.Fatal(e)
	}
	if third.State != "delivered" || third.PR.Disposition != "reused" || fx.PR.Disposition != "created" {
		t.Fatalf("third alias lost original effect: %s %+v", third.State, fx.PR)
	}
}

func TestCompletedAliasCannotClaimPendingOwnerMetadataOrAnotherNativeClosure(t *testing.T) {
	for _, changed := range []string{"pending-owner-metadata", "different-native-run", "different-source-ref"} {
		t.Run(changed, func(t *testing.T) {
			s, root, owner, fx, _ := journalReceipt(t)
			owner.State = "delivered"
			owner.NativeClosed = true
			owner.Metadata.State = "applied"
			fx.Complete = true
			fx.PR = &PRReceipt{MetadataApplied: true}
			owner.PR = fx.PR
			alias := owner
			alias.ID = "dlv_" + shortHash("another-key")
			alias.Registration.PublicationKey = "another-key"
			alias.RegistrationSHA256 = digest(alias.Registration)
			alias.Events = []Event{}
			alias.NativeClosed = false
			if changed == "pending-owner-metadata" {
				owner.State = "registered"
				owner.Metadata.State = "pending"
			} else if changed == "different-native-run" {
				alias.Manifest.WorkflowRunID = "another-native-run"
				alias.ManifestSHA256 = digest(alias.Manifest)
			} else {
				alias.Manifest.Agreement.SourceRef = "refs/heads/another-task-source"
				alias.ManifestSHA256 = digest(alias.Manifest)
			}
			if e := s.reuseCompleted(root, &alias, &fx, owner); e != nil {
				t.Fatal(e)
			}
			if alias.State != "blocked" || alias.NativeClosed {
				t.Fatalf("alias claimed unobserved completion: %s %t", alias.State, alias.NativeClosed)
			}
		})
	}
}
