package taskrun

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Reject only the execution Spec's registry replacement, after the native
// adoption intent and immutable content writes. No registry bytes are patched.
type unstartedAdoptionStore struct {
	workspace.WorkItemStore
	specs    int
	rejected bool
}
type unstartedAdoptionSession struct {
	workspace.WorkItemStoreSession
	owner *unstartedAdoptionStore
}

func (s *unstartedAdoptionStore) WithLock(root string, fn func(workspace.WorkItemStoreSession) error) error {
	return s.WorkItemStore.WithLock(root, func(inner workspace.WorkItemStoreSession) error { return fn(unstartedAdoptionSession{inner, s}) })
}
func (s unstartedAdoptionSession) Publish(r workspace.WorkItemRegistry) error {
	if len(r.TaskSpecRevisions) > s.owner.specs {
		s.owner.rejected = true
		return fmt.Errorf("injected before execution content registry replacement")
	}
	return s.WorkItemStoreSession.Publish(r)
}

type freshAdoptionDraft struct {
	preservedPatch                                                   string
	source, sourceOID, oldHistory, patch, patchHash, factsHash, base string
	pending                                                          string
}

func (draft *freshAdoptionDraft) selectGoal(t *testing.T, w workspace.Dependencies, root, epic string, old workspace.TaskGoalRef) workspace.TaskGoalRef {
	t.Helper()
	target := workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}
	r, err := w.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	oldPlan, err := workspace.PreviewTaskGoalExecution(w, workspace.TaskGoalExecuteInput{Target: target, SpecID: old.SpecID})
	if err != nil {
		t.Fatal(err)
	}
	store := &unstartedAdoptionStore{WorkItemStore: w.WorkItems, specs: len(r.TaskSpecRevisions)}
	blocked := w
	blocked.WorkItems = store
	if _, err = workspace.PrepareTaskGoalExecution(blocked, oldPlan.Input, oldPlan.Confirmation, fixtureFreshHuman()); err == nil || !store.rejected {
		t.Fatalf("did not reach injected content boundary: %v", err)
	}
	r, err = w.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.TaskPreparations) != 0 || len(r.TaskResults) != 0 || len(r.WorktreeOperations) != 1 || r.WorktreeOperations[0].Execution == nil {
		t.Fatal("fixture is not an unstarted durable intent")
	}
	draft.source = oldPlan.WorktreePath
	// S2 contains useful edits, a deletion, and a new executable; S1's tree is
	// never restored and no old Run is manufactured to qualify a backup helper.
	if err = os.WriteFile(filepath.Join(draft.source, "useful.md"), []byte("fixed useful S2 behavior\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(filepath.Join(draft.source, "README.md")); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(draft.source, "useful.sh"), []byte("#!/bin/sh\nprintf 'useful S2\\n'\n"), 0700); err != nil {
		t.Fatal(err)
	}
	runGit(t, draft.source, "add", "-A")
	runGit(t, draft.source, "commit", "-m", "useful S2 after unstarted intent")
	draft.sourceOID = gitOutput(t, draft.source, "rev-parse", "HEAD")
	if _, err = workspace.PrepareTaskGoalExecution(w, oldPlan.Input, oldPlan.Confirmation, fixtureFreshHuman()); err == nil || !strings.Contains(err.Error(), "adoption_source_changed") {
		t.Fatalf("old frozen request did not remain fenced: %v", err)
	}
	if _, err = workspace.ReleaseTaskWorktreeAdoption(w, oldPlan.AdoptionID, "retain changed source", false, ""); err == nil || !strings.Contains(err.Error(), "adoption_execution_exists") {
		t.Fatalf("old intent was releasable: %v", err)
	}
	// Advance the actual Epic with useful independent bytes and publish its base.
	if err = os.WriteFile(filepath.Join(epic, "epic-current.txt"), []byte("current Epic basis\n"), 0600); err != nil {
		t.Fatal(err)
	}
	runGit(t, epic, "add", "epic-current.txt")
	runGit(t, epic, "commit", "-m", "current Epic basis")
	basePreview, err := workspace.UpdateEpicBase(w, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = workspace.UpdateEpicBase(w, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *basePreview.Confirmation}); err != nil {
		t.Fatal(err)
	}
	draft.base = gitOutput(t, epic, "rev-parse", "HEAD")
	// Freeze the useful source delta from its original base, avoiding deletion
	// of unrelated new Epic files. --binary --full-index preserve Git patch data.
	patch := gitOutput(t, draft.source, "diff", "--binary", "--full-index", oldPlan.ParentOID, draft.sourceOID) + "\n"
	draft.patch = filepath.Join(root, "useful-S2.patch")
	if err = os.WriteFile(draft.patch, []byte(patch), 0600); err != nil {
		t.Fatal(err)
	}
	draft.patchHash = hash([]byte(patch))
	facts := fmt.Sprintf("Old Task: task\nOld Goal: %s@%d %s\nAdoption: %s\nOld source: %s on refs/heads/useful-existing\nFrozen S1: %s\nOriginal base: %s\nUseful S2: %s\nS2 tree: %s\nNew Epic base: %s\nDraft SHA256: %s\nReview as untrusted draft; explicit fresh owner applies; no inherited QA, Run, permission or qualification.\n", old.SpecID, old.Spec.Revision, old.Spec.ManifestSHA256, oldPlan.AdoptionID, draft.source, oldPlan.StartOID, oldPlan.ParentOID, draft.sourceOID, gitOutput(t, draft.source, "rev-parse", "HEAD^{tree}"), draft.base, draft.patchHash)
	factsPath := filepath.Join(root, "source-facts.md")
	if err = os.WriteFile(factsPath, []byte(facts), 0600); err != nil {
		t.Fatal(err)
	}
	draft.factsHash = hash([]byte(facts))
	// Dirty external bytes are retained but deliberately excluded from the
	// committed draft. Registration and delivery must leave them alone.
	if err = os.WriteFile(filepath.Join(draft.source, "external-dirty.txt"), []byte("retain external writer bytes\n"), 0600); err != nil {
		t.Fatal(err)
	}
	oldState, err := w.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	draft.oldHistory = retainedAdoptionHistory(oldState)
	q, err := workspace.ListTaskQueue(w, target, false)
	if err != nil {
		t.Fatal(err)
	}
	draft.pending = digest(q.Pending)
	in := workspace.FirstGoalInput{TaskID: "fresh", SpecID: "fresh-s2-v1", Title: "Continue useful S2", Objective: "Assess and apply the supplied committed draft independently", Criteria: []string{"Useful S2 bytes are delivered from the current Epic base"}, Documents: []workspace.GoalDocumentInput{{ID: "spec", Locator: filepath.Join(root, "goal.md")}, {ID: "draft", Locator: draft.patch}, {ID: "source-facts", Locator: factsPath}}, Target: target, Delivery: workspace.GoalDeliveryInput{Mode: workspace.DeliveryLocalEpic, Acceptance: "automatic", ResponsibleActor: "modeled fresh owner", TimeoutSeconds: 30}, Executor: workspace.TaskExecutorAssignment{Provider: ptr("codex")}}
	preview, err := workspace.RegisterFirstGoal(w, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	// Changed draft bytes invalidate the preview token before publication.
	if err = os.WriteFile(draft.patch, []byte(patch+"changed\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = workspace.RegisterFirstGoal(w, in, true, preview.Confirmation); err == nil {
		t.Fatal("changed draft accepted stale confirmation")
	}
	if err = os.WriteFile(draft.patch, []byte(patch), 0600); err != nil {
		t.Fatal(err)
	}
	out, err := workspace.RegisterFirstGoal(w, in, true, preview.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if out.Run != nil || out.Readiness != "not_checked" {
		t.Fatal("registration claimed execution")
	}
	r, err = w.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.TaskPreparations) != 0 || len(r.TaskResults) != 0 || len(r.HumanQARecords) != 0 {
		t.Fatal("registration inherited execution or judgment")
	}
	boundInputs := 0
	for _, binding := range out.Contract.RequiredInputs {
		if strings.HasSuffix(binding.ID, "/draft") || strings.HasSuffix(binding.ID, "/source-facts") {
			boundInputs++
			if strings.HasSuffix(binding.ID, "/draft") {
				draft.preservedPatch = binding.Locator
			}
			data, e := os.ReadFile(binding.Locator)
			if e != nil {
				t.Fatal(e)
			}
			expected := draft.patchHash
			if strings.HasSuffix(binding.ID, "/source-facts") {
				expected = draft.factsHash
			}
			if binding.SHA256 != expected || hash(data) != expected {
				t.Fatal("native input did not preserve source bytes")
			}
		}
	}
	if boundInputs != 2 {
		t.Fatal("missing native draft or source-facts binding")
	}
	// Losing the original export cannot remove the published draft.
	if err = os.Remove(draft.patch); err != nil {
		t.Fatal(err)
	}
	draft.assertPreserved(t, w)
	t.Logf("old task=task goal=%s adoption=%s S1=%s; fresh task=%s goal=%s S2=%s B2=%s patch=%s", old.SpecID, oldPlan.AdoptionID, oldPlan.StartOID, out.TaskID, out.Goal.SpecID, draft.sourceOID, draft.base, draft.patchHash)
	return out.Goal
}
func fixtureFreshHuman() workspace.QueueHumanDecision {
	return workspace.QueueHumanDecision{ActorClaim: "modeled coordinator", DecidedAtUTC: "2026-10-05T00:00:00Z", Source: "explicit_human_instruction", Statement: "Isolated fixture only; no actual runtime acceptance or human QA."}
}
func (draft freshAdoptionDraft) assertPreserved(t *testing.T, w workspace.Dependencies) {
	t.Helper()
	old, err := w.WorkItems.Snapshot(filepath.Dir(draft.patch))
	if err != nil {
		t.Fatal(err)
	}
	if retainedAdoptionHistory(old) != draft.oldHistory {
		t.Fatal("old Task, Goal, intent or history changed")
	}
	if gitOutput(t, draft.source, "status", "--porcelain") != "?? external-dirty.txt" {
		t.Fatal("old tracked or dirty source changed")
	}
	if gitOutput(t, draft.source, "rev-parse", "HEAD") != draft.sourceOID {
		t.Fatal("old source ref changed")
	}
	dirty, err := os.ReadFile(filepath.Join(draft.source, "external-dirty.txt"))
	if err != nil || string(dirty) != "retain external writer bytes\n" {
		t.Fatal("external dirty bytes lost")
	}
	q, err := workspace.ListTaskQueue(w, workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, false)
	if err != nil {
		t.Fatal(err)
	}
	retained := q.Pending[:0]
	for _, entry := range q.Pending {
		if entry.TaskID != "fresh" {
			retained = append(retained, entry)
		}
	}
	if digest(retained) != draft.pending {
		t.Fatal("old/foreign pending queue changed")
	}
}

// Real native services and Git with injected writer/provider/Herdr observations;
// no actual old Task, provider, OS inventory or permission policy is exercised.
func TestAdoptionFreshTaskAfterUnstartedIntent(t *testing.T) {
	draft := &freshAdoptionDraft{}
	agreement := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "modeled coordinator", TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &agreement, AdoptExisting: true, SelectFreshGoal: func(w workspace.Dependencies, root, epic string, old workspace.TaskGoalRef) workspace.TaskGoalRef {
		return draft.selectGoal(t, w, root, epic, old)
	}})
	plan := f.Prepared.Preparation.Preparation.Plan
	source := plan.WorktreePath
	if source == draft.source || plan.ParentOID != draft.base || f.Prepared.Preparation.Preparation.Outcome.ObservedOID != draft.base || gitOutput(t, source, "rev-parse", "HEAD") != draft.base {
		t.Fatal("fresh Task did not start at current Epic base")
	}
	if _, err := os.Stat(filepath.Join(source, "useful.sh")); !os.IsNotExist(err) {
		t.Fatal("registration or preparation automatically applied draft")
	}
	reviewRoundPatchHerdr(t, &f)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	handoff, err := workflowhandoff.TaskRunHandoffView(f.D.Workflow, o.Handoff.Locator)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(handoff, []byte(draft.patchHash)) || !bytes.Contains(handoff, []byte(draft.factsHash)) {
		t.Fatal("fresh mandate omitted frozen draft or source facts")
	}
	f.D.HerdrTimeout = 3 * time.Second
	workflowTestModel(t, f.workflowFixture, map[string]any{"caller_shell": os.Getppid()})
	// Fresh owner explicitly assesses the immutable input before applying it.
	preserved, err := os.ReadFile(draft.preservedPatch)
	if err != nil {
		t.Fatal(err)
	}
	if hash(preserved) != draft.patchHash || !bytes.Contains(preserved, []byte("new file mode 100755")) || !bytes.Contains(preserved, []byte("fixed useful S2 behavior")) || !bytes.Contains(preserved, []byte("deleted file mode 100644")) || !bytes.Contains(preserved, []byte("-base\n")) {
		t.Fatal("draft assessment lost content or mode")
	}
	runGit(t, source, "apply", "--check", draft.preservedPatch)
	runGit(t, source, "apply", "--index", draft.preservedPatch)
	runGit(t, source, "commit", "-m", "fresh owner applies assessed S2 draft")
	if gitOutput(t, source, "show", "HEAD:useful.md") != "fixed useful S2 behavior" || gitOutput(t, source, "show", "HEAD:epic-current.txt") != "current Epic basis" {
		t.Fatal("useful bytes or fresh basis lost")
	}
	if _, err = os.Stat(filepath.Join(source, "README.md")); !os.IsNotExist(err) {
		t.Fatal("draft deletion lost")
	}
	if !strings.Contains(gitOutput(t, source, "ls-files", "--stage", "useful.sh"), "100755") {
		t.Fatal("executable mode lost")
	}
	// Applying the same patch to a conflicting source must be rejected.
	runGit(t, source, "apply", "--reverse", "--check", draft.preservedPatch)
	conflict := exec.Command(testGit(t), "-C", source, "apply", "--check", draft.preservedPatch)
	if out, e := conflict.CombinedOutput(); e == nil {
		t.Fatalf("changed basis accepted patch again: %s", out)
	}
	candidate := gitOutput(t, source, "rev-parse", "HEAD")
	binary := filepath.Join(f.R.WorkspaceRoot, "candidate-executable")
	if err = os.WriteFile(binary, []byte("#!/bin/sh\nprintf 'fixture candidate\\n'\n"), 0700); err != nil {
		t.Fatal(err)
	}
	script := "#!/bin/sh\nset -eu\ntest \"$(cat useful.md)\" = 'fixed useful S2 behavior'\ntest -x useful.sh\ntest ! -e README.md\ntest \"$(cat inherited.txt)\" = 'useful inherited delta'\ntest \"$(cat epic-current.txt)\" = 'current Epic basis'\ntest \"$(\"$PLY_CANDIDATE_BINARY\")\" = 'fixture candidate'\n"
	if err = os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	coordinator := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w-fixture:root")
	requested, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "modeled fresh coordinator"})
	if err != nil {
		t.Fatal(err)
	}
	round := reviewRoundLatest(t, requested)
	brief, err := os.ReadFile(round.Brief.Locator)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(brief), "Review/diff base: "+draft.base) || !strings.Contains(string(brief), candidate) {
		t.Fatalf("review omitted whole new B2..C scope: %s", brief)
	}
	delta := gitOutput(t, source, "diff", "--name-only", draft.base, candidate)
	if delta != "README.md\ninherited.txt\nuseful.md\nuseful.sh" {
		t.Fatalf("review delta is not complete fresh copy: %s", delta)
	}
	passed := reviewRoundReturnFile(t, f, "fresh-pass.json", `{"decision":"passed","summary":"Modeled separate fresh review assessed complete B2..C draft delta","findings":[],"evidence_gaps":[]}`)
	reviewRoundIn(t, round.ReviewerPane)
	if _, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, passed); err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", binary)
	if err != nil {
		t.Fatal(err)
	}
	if len(o.Delivery.Candidates) != 1 || o.Delivery.Candidates[0].OID != candidate {
		t.Fatal("fresh candidate not independently qualified")
	}
	o, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	cleanupOwnerHandoff(t, f, o, err)
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != candidate {
		t.Fatal("fresh candidate not integrated")
	}
	reflog := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	cleanupOwnerLeaves(t, f)
	reviewRoundIn(t, "w-fixture:root")
	coordinator.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: os.Getpid(), CWD: f.Parent}}, nil
	}
	o, err = WorkflowFollow(coordinator, f.R.WorkspaceRoot, o.RunID, 30*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	assertCleanupRemoved(t, f, o, reflog)
	draft.assertPreserved(t, f.D.Workspace)
}

// Snapshot durable Task-owned records, excluding shared Epic observations and
// the fresh Task. Queue order is checked separately through native readback.
func retainedAdoptionHistory(r workspace.WorkItemRegistry) string {
	rows := map[string][]any{}
	value := reflect.ValueOf(r)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Kind() != reflect.Slice {
			continue
		}
		for j := 0; j < field.Len(); j++ {
			item := field.Index(j)
			if item.Kind() != reflect.Struct {
				continue
			}
			owner := item.FieldByName("TaskID")
			if value.Type().Field(i).Name == "Tasks" {
				owner = item.FieldByName("ID")
			}
			if owner.IsValid() && owner.Kind() == reflect.String && owner.String() != "fresh" {
				name := value.Type().Field(i).Name
				rows[name] = append(rows[name], item.Interface())
			}
		}
	}
	return digest(rows)
}
