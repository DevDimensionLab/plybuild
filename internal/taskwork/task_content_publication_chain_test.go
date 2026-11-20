package taskwork

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// These tests pin validateContentPublicationChain's own rules: which error
// wins when several rules are broken at once, and its text and class. The
// history of one Task is published once through the real commands and then
// held in memory; every test validates a changed copy of it.

// Publication keys of publicationChainHistory, in the order published.
const (
	createKey         = "task-create/task"
	problemKey        = "chain/problem"
	specKey           = "chain/spec"
	snapshotSpecKey   = "chain/snapshot-spec"
	assessmentKey     = "chain/assessment"
	nextAssessmentKey = "chain/next-assessment"
	selectionKey      = "chain/selection"
	reselectionKey    = "chain/reselection"
	withdrawalKey     = "chain/withdrawal"
)

// legacyImportKey names a case that is not in the history: the problem of the
// Task's creation, turned into the import of a summary that was registered
// before problems existed.
const legacyImportKey = "legacy import"

var publicationChainKeys = []string{createKey, problemKey, specKey, snapshotSpecKey, assessmentKey, nextAssessmentKey, selectionKey, reselectionKey, withdrawalKey}

// Positions in the history's lists of problem revisions, Spec revisions,
// assessments and selections, which are each in the order published.
const (
	chainFirst = iota
	chainSecond
)

// publicationChainHistory is the registry of Task "task" with every kind of
// publication, and the stored bytes of their manifests and requests: the
// problem from the Task's creation and an authored revision with a file, a
// Spec with a file and a revision that snapshots it, two ready assessments of
// that revision, and a selection, a reselection and a withdrawal.
type publicationChainHistory struct {
	work      workspace.ProjectWork
	registry  WorkItemRegistry
	manifests map[string][]byte
	requests  map[string][]byte
}

var (
	publicationChainOnce   sync.Once
	publicationChainShared *publicationChainHistory
)

// sharedPublicationChainHistory publishes the history once per test binary.
// Tests copy what they change; see newPublicationChainCase.
func sharedPublicationChainHistory(t *testing.T) publicationChainHistory {
	t.Helper()
	publicationChainOnce.Do(func() {
		history := publishPublicationChainHistory(t)
		publicationChainShared = &history
	})
	if publicationChainShared == nil {
		t.Fatal("the publication chain history was not published; see the first failing test")
	}
	return *publicationChainShared
}

func publishPublicationChainHistory(t *testing.T) publicationChainHistory {
	t.Helper()
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	work := testWork(t, root)
	record := func(publish func(Dependencies, TaskContentInput) (TaskContentMutationResult, error), in TaskContentInput) TaskContentOutcomeRef {
		t.Helper()
		out, err := publish(f.dependencies, in)
		if err != nil || out.OutcomeRef == nil {
			t.Fatalf("publish %s: %#v %v", in.File, out, err)
		}
		return *out.OutcomeRef
	}
	redraft := func(in TaskContentInput, change func(map[string]canonicaljson.Value)) TaskContentInput {
		t.Helper()
		raw, err := os.ReadFile(in.File)
		if err != nil {
			t.Fatal(err)
		}
		v, err := canonicaljson.DecodeStrict(raw)
		if err != nil {
			t.Fatal(err)
		}
		m := contentFields(v)
		change(m)
		if raw, err = canonicaljson.Marshal(contentObject(m)); err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(in.File, raw, 0o600); err != nil {
			t.Fatal(err)
		}
		return in
	}

	record(RecordTaskProblem, contentFixtureProblem(t, f, "chain/problem"))
	spec := record(RecordTaskSpec, contentFixtureSpec(t, f, "chain/spec", nil))
	specRef := &TaskRevisionRef{*spec.Revision, spec.ManifestSHA256}
	spec = record(RecordTaskSpec, redraft(contentFixtureSpec(t, f, "chain/snapshot-spec", specRef), func(m map[string]canonicaljson.Value) {
		source := contentObject(map[string]canonicaljson.Value{"kind": "snapshot", "manifest_sha256": specRef.ManifestSHA256, "document_id": "solution"})
		m["documents"] = []canonicaljson.Value{contentObject(map[string]canonicaljson.Value{"id": "solution", "source": source})}
	}))
	specRef = &TaskRevisionRef{*spec.Revision, spec.ManifestSHA256}

	checks := []canonicaljson.Value{}
	for _, id := range []string{"acceptance_coverage", "implementation_basis", "problem_coverage", "recovery", "scope_and_phases", "three_parts"} {
		checks = append(checks, contentObject(map[string]canonicaljson.Value{"id": id, "outcome": "pass", "reason": "Fixture assessment of the preserved solution.", "evidence_document_ids": []canonicaljson.Value{}}))
	}
	assess := func(key string, previous *TaskDecisionRef) *TaskDecisionRef {
		t.Helper()
		out := record(AssessTaskSpec, contentFixtureDraft(t, root, "WorkspaceTaskSpecAssessmentDraft@1", key, map[string]canonicaljson.Value{"spec_id": "solution", "spec": contentRefValue(specRef), "expected_previous_assessment": contentRefValue(previous), "outcome": "ready", "reason": "Fixture readiness claim.", "open_questions": []canonicaljson.Value{}, "checks": checks, "documents": []canonicaljson.Value{}}))
		return &TaskDecisionRef{*out.ID, out.ManifestSHA256}
	}
	assessment := assess("chain/assessment", nil)
	assessment = assess("chain/next-assessment", assessment)

	r, err := f.dependencies.WorkItems.Snapshot(work)
	if err != nil {
		t.Fatal(err)
	}
	solution := contentObject(map[string]canonicaljson.Value{"spec_id": "solution", "spec": contentRefValue(specRef), "problem": contentRefValue(taskContentState(r, "task").ProblemHead), "assessment": contentRefValue(assessment)})
	decide := func(publish func(Dependencies, TaskContentInput) (TaskContentMutationResult, error), key, action string, solution canonicaljson.Value, previous *TaskDecisionRef) *TaskDecisionRef {
		t.Helper()
		out := record(publish, contentFixtureDraft(t, root, "WorkspaceTaskSolutionSelectionDraft@1", key, map[string]canonicaljson.Value{"expected_previous_selection": contentRefValue(previous), "action": action, "solution": solution, "reason": "Fixture selection.", "human_decision": contentObject(map[string]canonicaljson.Value{"actor_claim": "fixture human, not actual approval", "decided_at_utc": "2026-09-29T12:00:00Z", "source": "explicit_human_instruction", "statement": "Fixture only: decide for the isolated test."})}))
		return &TaskDecisionRef{*out.ID, out.ManifestSHA256}
	}
	selection := decide(SelectTaskSolution, "chain/selection", "select", solution, nil)
	selection = decide(SelectTaskSolution, "chain/reselection", "select", solution, selection)
	decide(WithdrawTaskSolution, "chain/withdrawal", "withdraw", nil, selection)

	if r, err = f.dependencies.WorkItems.Snapshot(work); err != nil {
		t.Fatal(err)
	}
	if len(r.TaskContentPublications) != len(publicationChainKeys) {
		t.Fatalf("published %d publications, want %d", len(r.TaskContentPublications), len(publicationChainKeys))
	}
	stored := func(kind, digest string) []byte {
		t.Helper()
		raw, err := f.dependencies.TaskContent.Read(work, kind, digest)
		if err != nil {
			t.Fatal(err)
		}
		return raw
	}
	history := publicationChainHistory{work: work, registry: r, manifests: map[string][]byte{}, requests: map[string][]byte{}}
	for _, pub := range r.TaskContentPublications {
		history.manifests[pub.OutcomeRef.ManifestSHA256] = stored("manifests", pub.OutcomeRef.ManifestSHA256)
		history.requests[pub.RequestSHA256] = stored("requests", pub.RequestSHA256)
	}
	return history
}

// publicationChainCase is one publication with its manifest and request, as
// validateContentPublicationChain receives them, on its own copy of the
// history. A test changes the case and asks for the answer.
type publicationChainCase struct {
	t *testing.T
	publicationChainHistory
	pub      TaskContentPublication
	manifest map[string]canonicaljson.Value
	request  map[string]canonicaljson.Value
	// requestTail follows the canonical bytes of request.
	requestTail string
	// unreadable holds the read fault of a stored manifest, by digest.
	unreadable map[string]error
}

func newPublicationChainCase(t *testing.T, key string) *publicationChainCase {
	t.Helper()
	if key == legacyImportKey {
		return legacyImport(newPublicationChainCase(t, createKey))
	}
	c := &publicationChainCase{t: t, publicationChainHistory: sharedPublicationChainHistory(t), unreadable: map[string]error{}}
	for _, pub := range c.registry.TaskContentPublications {
		if pub.PublicationKey == key {
			c.pub = pub
			c.manifest = c.fields(c.manifests[pub.OutcomeRef.ManifestSHA256])
			c.request = c.fields(c.requests[pub.RequestSHA256])
		}
	}
	if c.manifest == nil {
		t.Fatalf("the history has no publication %s", key)
	}
	c.manifests = maps.Clone(c.manifests)
	c.registry.Tasks = slices.Clone(c.registry.Tasks)
	c.registry.TaskProblemRevisions = slices.Clone(c.registry.TaskProblemRevisions)
	c.registry.TaskSpecRevisions = slices.Clone(c.registry.TaskSpecRevisions)
	c.registry.TaskSpecAssessments = slices.Clone(c.registry.TaskSpecAssessments)
	c.registry.TaskSolutionSelections = slices.Clone(c.registry.TaskSolutionSelections)
	return c
}

func legacyImport(c *publicationChainCase) *publicationChainCase {
	c.t.Helper()
	summary, err := contentCanonical(map[string]any{"task_id": "task", "title": "Task", "description": "Description"})
	if err != nil {
		c.t.Fatal(err)
	}
	origin := contentObject(map[string]canonicaljson.Value{"kind": "legacy_summary_import", "legacy_summary_sha256": digestTaskBytes(summary)})
	c.pub.Operation, c.pub.IntentSHA256 = "problem_record", c.pub.RequestSHA256
	c.request["origin"], c.manifest["origin"] = origin, origin
	return c
}

func (c *publicationChainCase) fields(raw []byte) map[string]canonicaljson.Value {
	c.t.Helper()
	v, err := canonicaljson.DecodeStrict(raw)
	if err != nil {
		c.t.Fatal(err)
	}
	return contentFields(v)
}

func (c *publicationChainCase) canonical(fields map[string]canonicaljson.Value) []byte {
	c.t.Helper()
	raw, err := canonicaljson.Marshal(contentObject(fields))
	if err != nil {
		c.t.Fatal(err)
	}
	return raw
}

// read serves the stored manifests of the case; it does not hash them, so a
// test may store changed bytes under the digest of the original.
func (c *publicationChainCase) read(_ workspace.ProjectWork, kind, digest string) ([]byte, error) {
	if err := c.unreadable[digest]; err != nil {
		return nil, err
	}
	raw, stored := c.manifests[digest]
	if kind != "manifests" || !stored {
		return nil, fs.ErrNotExist
	}
	return raw, nil
}

// changeStored changes a stored manifest and keeps it under the same digest.
func (c *publicationChainCase) changeStored(digest string, change func(manifest map[string]canonicaljson.Value)) {
	c.t.Helper()
	manifest := c.fields(c.manifests[digest])
	change(manifest)
	c.manifests[digest] = c.canonical(manifest)
}

func (c *publicationChainCase) task() *TaskRecord {
	task, _ := findTask(c.registry, c.pub.TaskID)
	return task
}

// answerTo validates the case with the given request bytes and names the
// result as chainAnswer does.
func (c *publicationChainCase) answerTo(request []byte) string {
	return chainAnswer(validateContentPublicationChain(&TaskContentStorage{readCaptured: c.read}, c.work, c.registry, c.pub, c.manifest, request))
}

func (c *publicationChainCase) answer() string {
	c.t.Helper()
	return c.answerTo(append(c.canonical(c.request), c.requestTail...))
}

const chainAccepted = "accepted"

// chainAnswer names the class and text of the validation result: "accepted",
// a classified content error as "code: detail", and any other error as
// "unclassified: text". The callers add a class to the unclassified ones.
func chainAnswer(err error) string {
	if err == nil {
		return chainAccepted
	}
	if _, classified := err.(*TaskContentError); classified {
		return err.Error()
	}
	return "unclassified: " + err.Error()
}

var foreignDigest = "sha256:" + strings.Repeat("9", 64)

var errChainRead = errors.New("fixture read fault")

const chainNextAction = "Inspect the required artifact and its preserved publication before retrying."

// unreadableAnswer is the answer when a stored manifest fails to read with
// cause: the read and the manifest reader each wrap it under the given class.
func unreadableAnswer(class, digest, cause string) string {
	return fmt.Sprintf("%[1]s: read publication required manifests %[2]s: %[1]s: read content required manifests %[2]s: %[3]s; %[4]s; %[4]s", class, digest, cause, chainNextAction)
}

// rejectedAnswer is the answer when the manifest reader rejects the stored
// bytes itself, with the given class and detail.
func rejectedAnswer(class, digest, detail string) string {
	return fmt.Sprintf("%[1]s: read publication required manifests %[2]s: %[1]s: %[3]s; %[4]s", class, digest, detail, chainNextAction)
}

func faultAnswer(digest string) string {
	return unreadableAnswer("task_content_observation_unknown", digest, errChainRead.Error())
}

// changeObject changes the object member key of a manifest or a request.
func changeObject(fields map[string]canonicaljson.Value, key string, change func(object map[string]canonicaljson.Value)) {
	object := contentFields(fields[key])
	change(object)
	fields[key] = contentObject(object)
}

// changeDocument changes document i of a manifest or a request.
func changeDocument(fields map[string]canonicaljson.Value, i int, change func(document map[string]canonicaljson.Value)) {
	documents := slices.Clone(contentArray(fields, "documents"))
	document := contentFields(documents[i])
	change(document)
	documents[i] = contentObject(document)
	fields["documents"] = documents
}

// changeSource changes the source of document i of the request.
func (c *publicationChainCase) changeSource(i int, change func(source map[string]canonicaljson.Value)) {
	changeDocument(c.request, i, func(document map[string]canonicaljson.Value) { changeObject(document, "source", change) })
}

// snapshotEvidence gives an assessment one document per ID, in the given
// order, as evidence of its first check. Each snapshots the first Spec's
// document "solution" under the new ID.
func snapshotEvidence(c *publicationChainCase, ids ...string) {
	firstSpec := c.registry.TaskSpecRevisions[chainFirst].ManifestSHA256
	source := contentObject(map[string]canonicaljson.Value{"kind": "snapshot", "manifest_sha256": firstSpec, "document_id": "solution"})
	solution := contentFields(contentArray(c.fields(c.manifests[firstSpec]), "documents")[0])
	requested, published, used := []canonicaljson.Value{}, []canonicaljson.Value{}, []canonicaljson.Value{}
	for _, id := range ids {
		solution["id"] = id
		requested = append(requested, contentObject(map[string]canonicaljson.Value{"id": id, "source": source}))
		published = append(published, contentObject(solution))
		used = append(used, id)
	}
	checks := slices.Clone(contentArray(c.request, "checks"))
	check := contentFields(checks[0])
	check["evidence_document_ids"] = used
	checks[0] = contentObject(check)
	c.request["documents"], c.manifest["documents"] = requested, published
	c.request["checks"], c.manifest["checks"] = checks, checks
}

// addManifestDocument gives the manifest one document more than its request.
func addManifestDocument(c *publicationChainCase) {
	extra := contentObject(map[string]canonicaljson.Value{"id": "unrequested"})
	c.manifest["documents"] = append(slices.Clone(contentArray(c.manifest, "documents")), extra)
}

// chainRule breaks one rule of the publication named by key and gives the
// answer of that rule.
type chainRule struct {
	name   string
	key    string
	breach func(c *publicationChainCase)
	answer string
}

func runChainRules(t *testing.T, rules []chainRule) {
	t.Helper()
	for _, rule := range rules {
		t.Run(rule.name, func(t *testing.T) {
			c := newPublicationChainCase(t, rule.key)
			rule.breach(c)
			if got := c.answer(); got != rule.answer {
				t.Fatalf("got %q, want %q", got, rule.answer)
			}
		})
	}
}

const (
	predecessorAnswer  = "unclassified: predecessor differs from registered chain"
	identityAnswerSpec = "unclassified: Spec identity differs"
	identityAnswerAsm  = "unclassified: assessment identity differs"
	generatedAnswer    = "unclassified: generated problem differs from the original Task summary"
	generatedDocAnswer = "unclassified: generated problem document differs from the original Task summary"
	unownedSnapshot    = "unclassified: snapshot source is not owned by Task"
	unregisteredAsm    = "unclassified: selection assessment is not registered"
)

// requestRules are the rules that every publication meets first, in the order
// they answer. A Task's creation has its own intent rule.
func requestRules(draftKind string, intent chainRule) []chainRule {
	incompatible := fmt.Sprintf("task_content_reader_incompatible: required %s schema 1 exceeds reader capability %s; unsupported fields: future", draftKind, taskContentReaderCapability)
	return []chainRule{
		{name: "request with an unknown field", breach: func(c *publicationChainCase) { c.request["future"] = "unknown" }, answer: incompatible},
		{name: "request that is not canonical", breach: func(c *publicationChainCase) { c.requestTail = "\n" }, answer: "unclassified: request is not canonical"},
		intent,
		{name: "request field that the manifest changed", breach: func(c *publicationChainCase) { c.manifest["publication_key"] = "chain/other" }, answer: "unclassified: request field publication_key differs"},
		{name: "publication at another time than the manifest's", breach: func(c *publicationChainCase) { c.pub.RecordedAtUTC = "2026-11-20T00:00:00Z" }, answer: "unclassified: publication time differs"},
	}
}

var (
	requestIntentRule = chainRule{name: "intent other than the request", breach: func(c *publicationChainCase) { c.pub.IntentSHA256 = foreignDigest }, answer: "unclassified: request and intent digests differ"}
	createIntentRule  = chainRule{name: "Task changed after its creation", breach: func(c *publicationChainCase) { c.task().RepoID = "other" }, answer: "unclassified: create intent differs"}
)

// generatedProblemRules follow the request rules of a problem that Ply wrote
// from the Task's summary.
func generatedProblemRules() []chainRule {
	return []chainRule{
		{name: "outcome without a revision", breach: func(c *publicationChainCase) { c.pub.OutcomeRef.Revision = nil }, answer: "unclassified: problem revision differs"},
		{name: "problem registered before the first", breach: func(c *publicationChainCase) {
			c.registry.TaskProblemRevisions = append(c.registry.TaskProblemRevisions, TaskProblemReference{TaskID: "task", Revision: 0, ManifestSHA256: foreignDigest})
		}, answer: predecessorAnswer},
		{name: "title other than the Task's", breach: func(c *publicationChainCase) { c.request["title"], c.manifest["title"] = "Other", "Other" }, answer: generatedAnswer},
		{name: "document of another size than the summary", breach: func(c *publicationChainCase) {
			changeDocument(c.manifest, 0, func(document map[string]canonicaljson.Value) { document["size_bytes"] = int64(1) })
		}, answer: generatedDocAnswer},
	}
}

// chainRuleOrders lists, for one publication of each kind, rules that can be
// broken together, in the order validateContentPublicationChain answers them.
func chainRuleOrders(t *testing.T) map[string][]chainRule {
	t.Helper()
	h := sharedPublicationChainHistory(t)
	firstSpec, secondSpec := h.registry.TaskSpecRevisions[chainFirst].ManifestSHA256, h.registry.TaskSpecRevisions[chainSecond].ManifestSHA256
	secondAssessment := h.registry.TaskSpecAssessments[chainSecond].ManifestSHA256
	unrequestedDocument := chainRule{name: "manifest with an unrequested document", breach: addManifestDocument, answer: "unclassified: document set differs from request"}
	renamedDocument := chainRule{name: "document under another ID", breach: func(c *publicationChainCase) {
		changeDocument(c.manifest, 0, func(document map[string]canonicaljson.Value) { document["id"] = "renamed" })
	}, answer: "unclassified: document ID differs"}
	resized := func(name, answer string) chainRule {
		return chainRule{name: name, breach: func(c *publicationChainCase) {
			changeDocument(c.manifest, 0, func(document map[string]canonicaljson.Value) { document["size_bytes"] = int64(1) })
		}, answer: answer}
	}
	return map[string][]chainRule{
		createKey: append(requestRules("WorkspaceTaskProblemDraft@1", createIntentRule), generatedProblemRules()...),
		legacyImportKey: append(append(requestRules("WorkspaceTaskProblemDraft@1", requestIntentRule), generatedProblemRules()...),
			chainRule{name: "digest of another summary", breach: func(c *publicationChainCase) {
				changeObject(c.request, "origin", func(origin map[string]canonicaljson.Value) { origin["legacy_summary_sha256"] = foreignDigest })
				c.manifest["origin"] = c.request["origin"]
			}, answer: "unclassified: legacy summary digest differs"}),
		problemKey: append(requestRules("WorkspaceTaskProblemDraft@1", requestIntentRule),
			chainRule{name: "outcome without a revision", breach: func(c *publicationChainCase) { c.pub.OutcomeRef.Revision = nil }, answer: "unclassified: problem revision differs"},
			chainRule{name: "predecessor registered with another manifest", breach: func(c *publicationChainCase) {
				c.registry.TaskProblemRevisions[chainFirst].ManifestSHA256 = foreignDigest
			}, answer: predecessorAnswer},
			unrequestedDocument,
			renamedDocument,
			resized("document of another size than its file", "unclassified: document metadata differs"),
			chainRule{name: "document from another file", breach: func(c *publicationChainCase) {
				changeDocument(c.manifest, 0, func(document map[string]canonicaljson.Value) {
					changeObject(document, "provenance", func(provenance map[string]canonicaljson.Value) {
						provenance["origin_locator"] = "/elsewhere/problem.md"
					})
				})
			}, answer: "unclassified: document provenance differs"},
		),
		snapshotSpecKey: append(requestRules("WorkspaceTaskSpecDraft@1", requestIntentRule),
			chainRule{name: "outcome without a revision", breach: func(c *publicationChainCase) { c.pub.OutcomeRef.Revision = nil }, answer: identityAnswerSpec},
			chainRule{name: "predecessor registered with another manifest", breach: func(c *publicationChainCase) { c.registry.TaskSpecRevisions[chainFirst].ManifestSHA256 = foreignDigest }, answer: predecessorAnswer},
			chainRule{name: "problem registered with another manifest", breach: func(c *publicationChainCase) {
				c.registry.TaskProblemRevisions[chainSecond].ManifestSHA256 = foreignDigest
			}, answer: "unclassified: problem reference is not owned by Task"},
			unrequestedDocument,
			renamedDocument,
			chainRule{name: "snapshot of a manifest that the named Task does not own", breach: func(c *publicationChainCase) {
				c.changeSource(0, func(source map[string]canonicaljson.Value) { source["task_id"] = "other" })
			}, answer: unownedSnapshot},
			chainRule{name: "unreadable snapshot source", breach: func(c *publicationChainCase) { c.unreadable[firstSpec] = errChainRead }, answer: faultAnswer(firstSpec)},
			resized("snapshot of another size than its source", "unclassified: snapshot descriptor differs"),
		),
		nextAssessmentKey: append(requestRules("WorkspaceTaskSpecAssessmentDraft@1", requestIntentRule),
			chainRule{name: "outcome without an ID", breach: func(c *publicationChainCase) { c.pub.OutcomeRef.ID = nil }, answer: identityAnswerAsm},
			chainRule{name: "predecessor registered with another manifest", breach: func(c *publicationChainCase) {
				c.registry.TaskSpecAssessments[chainFirst].ManifestSHA256 = foreignDigest
			}, answer: predecessorAnswer},
			chainRule{name: "Spec registered with another manifest", breach: func(c *publicationChainCase) {
				c.registry.TaskSpecRevisions[chainSecond].ManifestSHA256 = foreignDigest
			}, answer: "unclassified: Spec reference is not owned by Task"},
			unrequestedDocument,
		),
		reselectionKey: append(requestRules("WorkspaceTaskSolutionSelectionDraft@1", requestIntentRule),
			chainRule{name: "outcome without an ID", breach: func(c *publicationChainCase) { c.pub.OutcomeRef.ID = nil }, answer: "unclassified: selection identity differs"},
			chainRule{name: "predecessor registered with another manifest", breach: func(c *publicationChainCase) {
				c.registry.TaskSolutionSelections[chainFirst].ManifestSHA256 = foreignDigest
			}, answer: predecessorAnswer},
			chainRule{name: "assessment registered with another manifest", breach: func(c *publicationChainCase) {
				c.registry.TaskSpecAssessments[chainSecond].ManifestSHA256 = foreignDigest
			}, answer: unregisteredAsm},
			chainRule{name: "unreadable assessment", breach: func(c *publicationChainCase) { c.unreadable[secondAssessment] = errChainRead }, answer: faultAnswer(secondAssessment)},
			chainRule{name: "assessment that is not ready", breach: func(c *publicationChainCase) {
				c.changeStored(secondAssessment, func(manifest map[string]canonicaljson.Value) { manifest["outcome"] = "rejected" })
			}, answer: "unclassified: selection does not bind a ready assessment"},
			chainRule{name: "unreadable Spec", breach: func(c *publicationChainCase) { c.unreadable[secondSpec] = errChainRead }, answer: faultAnswer(secondSpec)},
			chainRule{name: "Spec of another problem than the selected", breach: func(c *publicationChainCase) {
				c.changeStored(secondSpec, func(manifest map[string]canonicaljson.Value) {
					manifest["problem"] = contentRefValue(TaskRevisionRef{1, c.registry.TaskProblemRevisions[chainFirst].ManifestSHA256})
				})
			}, answer: "unclassified: selection problem and Spec differ"},
			unrequestedDocument,
		),
	}
}

func TestPublicationChainAcceptsEveryPublishedKind(t *testing.T) {
	t.Parallel()
	for _, key := range append(slices.Clone(publicationChainKeys), legacyImportKey) {
		if got := newPublicationChainCase(t, key).answer(); got != chainAccepted {
			t.Errorf("%s: got %q", key, got)
		}
	}
}

// Each rule answers when it alone is broken, and before every rule after it
// when both are broken.
func TestPublicationChainEarlierRuleAnswersFirst(t *testing.T) {
	t.Parallel()
	for key, rules := range chainRuleOrders(t) {
		for i, earlier := range rules {
			t.Run(key+"/"+earlier.name, func(t *testing.T) {
				c := newPublicationChainCase(t, key)
				earlier.breach(c)
				if got := c.answer(); got != earlier.answer {
					t.Fatalf("alone: got %q, want %q", got, earlier.answer)
				}
				for _, later := range rules[i+1:] {
					c = newPublicationChainCase(t, key)
					earlier.breach(c)
					later.breach(c)
					if got := c.answer(); got != earlier.answer {
						t.Errorf("with %s: got %q, want %q", later.name, got, earlier.answer)
					}
				}
			})
		}
	}
}

func TestPublicationChainRequestDecodesAsItsOperation(t *testing.T) {
	t.Parallel()
	t.Run("empty request", func(t *testing.T) {
		if got, want := newPublicationChainCase(t, problemKey).answerTo(nil), "task_content_invalid_input: decode JSON: EOF"; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
	runChainRules(t, []chainRule{
		{name: "summary origin outside a Task's creation", key: createKey, breach: func(c *publicationChainCase) {
			c.pub.Operation, c.pub.IntentSHA256 = "problem_record", c.pub.RequestSHA256
		}, answer: "task_content_invalid_input: origin: invalid problem origin"},
		{name: "withdrawal published as a selection", key: withdrawalKey, breach: func(c *publicationChainCase) { c.pub.Operation = "spec_select" }, answer: "task_content_invalid_input: action: must be one of select"},
	})
}

func TestPublicationChainCreateIntentIsTheRegisteredTask(t *testing.T) {
	t.Parallel()
	if c := newPublicationChainCase(t, createKey); c.pub.RequestSHA256 == c.pub.IntentSHA256 {
		t.Fatal("the creation's request is its intent; the history no longer shows that they may differ")
	}
	runChainRules(t, []chainRule{
		{name: "Task is not registered", key: createKey, breach: func(c *publicationChainCase) { c.registry.Tasks = nil }, answer: "unclassified: missing Task"},
		{name: "intent of another creation", key: createKey, breach: func(c *publicationChainCase) { c.pub.IntentSHA256 = foreignDigest }, answer: "unclassified: create intent differs"},
	})
}

func TestPublicationChainRequestFieldsEqualTheManifest(t *testing.T) {
	t.Parallel()
	withoutPrevious := func(c *publicationChainCase) { c.manifest["previous"] = nil }
	const previousAnswer = "unclassified: request field previous differs"
	runChainRules(t, []chainRule{
		{name: "problem's expected previous is the manifest's previous", key: problemKey, breach: withoutPrevious, answer: previousAnswer},
		{name: "Spec's expected previous is the manifest's previous", key: snapshotSpecKey, breach: withoutPrevious, answer: previousAnswer},
		{name: "expected previous assessment is the manifest's previous", key: nextAssessmentKey, breach: withoutPrevious, answer: previousAnswer},
		{name: "expected previous selection is the manifest's previous", key: reselectionKey, breach: withoutPrevious, answer: previousAnswer},
		{name: "null request field equals a missing manifest field", key: problemKey, breach: func(c *publicationChainCase) { delete(c.manifest, "deadline") }, answer: chainAccepted},
		{name: "registry upgrade is not compared", key: problemKey, breach: func(c *publicationChainCase) { c.manifest["registry_upgrade"] = "ignored" }, answer: chainAccepted},
	})
}

// Request fields are compared in map order, so which of several differing
// fields the answer names changes from run to run.
func TestPublicationChainNamesOneOfSeveralDifferingFields(t *testing.T) {
	t.Parallel()
	c := newPublicationChainCase(t, problemKey)
	c.manifest["title"], c.manifest["summary"] = "Other", "Other"
	if got := c.answer(); got != "unclassified: request field title differs" && got != "unclassified: request field summary differs" {
		t.Fatalf("got %q", got)
	}
	c = newPublicationChainCase(t, problemKey)
	c.manifest = nil
	if got := c.answer(); !strings.HasPrefix(got, "unclassified: request field ") {
		t.Fatalf("nil manifest: got %q", got)
	}
}

func TestPublicationChainOutcomeIdentity(t *testing.T) {
	t.Parallel()
	outcome := func(change func(ref *TaskContentOutcomeRef)) func(*publicationChainCase) {
		return func(c *publicationChainCase) { change(&c.pub.OutcomeRef) }
	}
	runChainRules(t, []chainRule{
		{name: "problem of another revision", key: problemKey, breach: outcome(func(ref *TaskContentOutcomeRef) { ref.Revision = intPointer(3) }), answer: "unclassified: problem revision differs"},
		{name: "Spec without an ID", key: specKey, breach: outcome(func(ref *TaskContentOutcomeRef) { ref.SpecID = nil }), answer: identityAnswerSpec},
		{name: "Spec under another ID", key: specKey, breach: outcome(func(ref *TaskContentOutcomeRef) { ref.SpecID = stringPointer("elsewhere") }), answer: identityAnswerSpec},
		{name: "Spec of another revision", key: specKey, breach: outcome(func(ref *TaskContentOutcomeRef) { ref.Revision = intPointer(2) }), answer: identityAnswerSpec},
		{name: "assessment under another ID", key: assessmentKey, breach: outcome(func(ref *TaskContentOutcomeRef) { ref.ID = stringPointer("asm_other") }), answer: identityAnswerAsm},
		{name: "assessment without a Spec ID", key: assessmentKey, breach: outcome(func(ref *TaskContentOutcomeRef) { ref.SpecID = nil }), answer: identityAnswerAsm},
		{name: "assessment of another Spec", key: assessmentKey, breach: outcome(func(ref *TaskContentOutcomeRef) { ref.SpecID = stringPointer("elsewhere") }), answer: identityAnswerAsm},
		{name: "assessment without a Spec revision", key: assessmentKey, breach: outcome(func(ref *TaskContentOutcomeRef) { ref.Revision = nil }), answer: identityAnswerAsm},
		{name: "assessment with its ordinal as Spec revision", key: assessmentKey, breach: outcome(func(ref *TaskContentOutcomeRef) { ref.Revision = intPointer(1) }), answer: identityAnswerAsm},
		{name: "selection under another ID", key: selectionKey, breach: outcome(func(ref *TaskContentOutcomeRef) { ref.ID = stringPointer("sel_other") }), answer: "unclassified: selection identity differs"},
		{name: "outcome of an unknown kind", key: problemKey, breach: outcome(func(ref *TaskContentOutcomeRef) { ref.Kind = "delivery" }), answer: "unclassified: unknown publication outcome"},
		{name: "outcome without a kind", key: problemKey, breach: outcome(func(ref *TaskContentOutcomeRef) { ref.Kind = "" }), answer: "unclassified: unknown publication outcome"},
	})
}

// The predecessor is the last registered record just before the publication
// in its own chain: the Task's problems, one Spec's revisions, the assessments
// of one Spec revision, or the Task's selections.
func TestPublicationChainPredecessorIsTheLastRegisteredMatch(t *testing.T) {
	t.Parallel()
	registry := func(change func(r *WorkItemRegistry)) func(*publicationChainCase) {
		return func(c *publicationChainCase) { change(&c.registry) }
	}
	problem := func(task TaskID, revision int) TaskProblemReference {
		return TaskProblemReference{TaskID: task, Revision: revision, ManifestSHA256: foreignDigest}
	}
	spec := func(task TaskID, id string, revision int) TaskSpecReference {
		return TaskSpecReference{TaskID: task, SpecID: id, Revision: revision, ManifestSHA256: foreignDigest}
	}
	assessment := func(task TaskID, spec string, revision, ordinal int) TaskAssessmentReference {
		return TaskAssessmentReference{TaskID: task, SpecID: spec, SpecRevision: revision, Ordinal: ordinal, ID: "asm_other", ManifestSHA256: foreignDigest}
	}
	selection := func(task TaskID, ordinal int) TaskSelectionReference {
		return TaskSelectionReference{TaskID: task, Ordinal: ordinal, ID: "sel_other", ManifestSHA256: foreignDigest}
	}
	runChainRules(t, []chainRule{
		{name: "problem without its predecessor", key: problemKey, breach: registry(func(r *WorkItemRegistry) { r.TaskProblemRevisions = r.TaskProblemRevisions[1:] }), answer: predecessorAnswer},
		{name: "problem after a later duplicate of its predecessor", key: problemKey, breach: registry(func(r *WorkItemRegistry) { r.TaskProblemRevisions = append(r.TaskProblemRevisions, problem("task", 1)) }), answer: predecessorAnswer},
		{name: "problem after an earlier duplicate of its predecessor", key: problemKey, breach: registry(func(r *WorkItemRegistry) {
			r.TaskProblemRevisions = slices.Insert(r.TaskProblemRevisions, 0, problem("task", 1))
		}), answer: chainAccepted},
		{name: "problem beside other chains", key: problemKey, breach: registry(func(r *WorkItemRegistry) {
			r.TaskProblemRevisions = append(r.TaskProblemRevisions, problem("other", 1), problem("task", 0), problem("task", 3))
		}), answer: chainAccepted},

		{name: "Spec without its predecessor", key: snapshotSpecKey, breach: registry(func(r *WorkItemRegistry) { r.TaskSpecRevisions = r.TaskSpecRevisions[1:] }), answer: predecessorAnswer},
		{name: "Spec after a later duplicate of its predecessor", key: snapshotSpecKey, breach: registry(func(r *WorkItemRegistry) {
			r.TaskSpecRevisions = append(r.TaskSpecRevisions, spec("task", "solution", 1))
		}), answer: predecessorAnswer},
		{name: "Spec after an earlier duplicate of its predecessor", key: snapshotSpecKey, breach: registry(func(r *WorkItemRegistry) {
			r.TaskSpecRevisions = slices.Insert(r.TaskSpecRevisions, 0, spec("task", "solution", 1))
		}), answer: chainAccepted},
		{name: "Spec beside other chains", key: snapshotSpecKey, breach: registry(func(r *WorkItemRegistry) {
			r.TaskSpecRevisions = append(r.TaskSpecRevisions, spec("other", "solution", 1), spec("task", "elsewhere", 1), spec("task", "solution", 0), spec("task", "solution", 3))
		}), answer: chainAccepted},

		{name: "assessment without its predecessor", key: nextAssessmentKey, breach: registry(func(r *WorkItemRegistry) { r.TaskSpecAssessments = r.TaskSpecAssessments[1:] }), answer: predecessorAnswer},
		{name: "assessment after a later duplicate of its predecessor", key: nextAssessmentKey, breach: registry(func(r *WorkItemRegistry) {
			r.TaskSpecAssessments = append(r.TaskSpecAssessments, assessment("task", "solution", 2, 1))
		}), answer: predecessorAnswer},
		{name: "assessment after an earlier duplicate of its predecessor", key: nextAssessmentKey, breach: registry(func(r *WorkItemRegistry) {
			r.TaskSpecAssessments = slices.Insert(r.TaskSpecAssessments, 0, assessment("task", "solution", 2, 1))
		}), answer: chainAccepted},
		{name: "assessment beside other chains", key: nextAssessmentKey, breach: registry(func(r *WorkItemRegistry) {
			r.TaskSpecAssessments = append(r.TaskSpecAssessments, assessment("other", "solution", 2, 1), assessment("task", "elsewhere", 2, 1), assessment("task", "solution", 1, 1), assessment("task", "solution", 2, 0), assessment("task", "solution", 2, 3))
		}), answer: chainAccepted},

		{name: "selection without its predecessor", key: reselectionKey, breach: registry(func(r *WorkItemRegistry) { r.TaskSolutionSelections = r.TaskSolutionSelections[1:] }), answer: predecessorAnswer},
		{name: "selection after a later duplicate of its predecessor", key: reselectionKey, breach: registry(func(r *WorkItemRegistry) {
			r.TaskSolutionSelections = append(r.TaskSolutionSelections, selection("task", 1))
		}), answer: predecessorAnswer},
		{name: "selection after an earlier duplicate of its predecessor", key: reselectionKey, breach: registry(func(r *WorkItemRegistry) {
			r.TaskSolutionSelections = slices.Insert(r.TaskSolutionSelections, 0, selection("task", 1))
		}), answer: chainAccepted},
		{name: "selection beside other chains", key: reselectionKey, breach: registry(func(r *WorkItemRegistry) {
			r.TaskSolutionSelections = append(r.TaskSolutionSelections, selection("other", 1), selection("task", 0), selection("task", 4))
		}), answer: chainAccepted},
	})
}

func TestPublicationChainReferencesAreOwnedByTheTask(t *testing.T) {
	t.Parallel()
	const unownedProblem, unownedSpec = "unclassified: problem reference is not owned by Task", "unclassified: Spec reference is not owned by Task"
	problem := func(change func(ref *TaskProblemReference)) func(*publicationChainCase) {
		return func(c *publicationChainCase) { change(&c.registry.TaskProblemRevisions[chainSecond]) }
	}
	spec := func(change func(ref *TaskSpecReference)) func(*publicationChainCase) {
		return func(c *publicationChainCase) { change(&c.registry.TaskSpecRevisions[chainSecond]) }
	}
	runChainRules(t, []chainRule{
		{name: "Spec's problem registered to another Task", key: specKey, breach: problem(func(ref *TaskProblemReference) { ref.TaskID = "other" }), answer: unownedProblem},
		{name: "Spec's problem registered as another revision", key: specKey, breach: problem(func(ref *TaskProblemReference) { ref.Revision = 3 }), answer: unownedProblem},
		{name: "assessed Spec registered to another Task", key: assessmentKey, breach: spec(func(ref *TaskSpecReference) { ref.TaskID = "other" }), answer: unownedSpec},
		{name: "assessed Spec registered under another ID", key: assessmentKey, breach: spec(func(ref *TaskSpecReference) { ref.SpecID = "elsewhere" }), answer: unownedSpec},
		{name: "assessed Spec registered as another revision", key: assessmentKey, breach: spec(func(ref *TaskSpecReference) { ref.Revision = 3 }), answer: unownedSpec},
	})
}

func TestPublicationChainSelectionBindsAReadyAssessmentOfItsSpec(t *testing.T) {
	t.Parallel()
	h := sharedPublicationChainHistory(t)
	firstSpec, secondSpec := h.registry.TaskSpecRevisions[chainFirst].ManifestSHA256, h.registry.TaskSpecRevisions[chainSecond].ManifestSHA256
	secondAssessment := h.registry.TaskSpecAssessments[chainSecond].ManifestSHA256
	registered := func(change func(ref *TaskAssessmentReference)) func(*publicationChainCase) {
		return func(c *publicationChainCase) { change(&c.registry.TaskSpecAssessments[chainSecond]) }
	}
	stored := func(digest string, change func(manifest map[string]canonicaljson.Value)) func(*publicationChainCase) {
		return func(c *publicationChainCase) { c.changeStored(digest, change) }
	}
	const differs = "unclassified: selection problem and Spec differ"
	runChainRules(t, []chainRule{
		{name: "assessment registered to another Task", key: selectionKey, breach: registered(func(ref *TaskAssessmentReference) { ref.TaskID = "other" }), answer: unregisteredAsm},
		{name: "assessment registered to another Spec", key: selectionKey, breach: registered(func(ref *TaskAssessmentReference) { ref.SpecID = "elsewhere" }), answer: unregisteredAsm},
		{name: "assessment registered to another Spec revision", key: selectionKey, breach: registered(func(ref *TaskAssessmentReference) { ref.SpecRevision = 1 }), answer: unregisteredAsm},
		{name: "assessment registered under another ID", key: selectionKey, breach: registered(func(ref *TaskAssessmentReference) { ref.ID = "asm_other" }), answer: unregisteredAsm},
		{name: "assessment of another Spec revision than the selected", key: selectionKey, breach: stored(secondAssessment, func(manifest map[string]canonicaljson.Value) {
			manifest["spec"] = contentRefValue(TaskRevisionRef{1, firstSpec})
		}), answer: "unclassified: selection does not bind a ready assessment"},
		{name: "assessment missing from the store", key: selectionKey, breach: func(c *publicationChainCase) { delete(c.manifests, secondAssessment) }, answer: unreadableAnswer("task_content_missing", secondAssessment, "file does not exist")},
		{name: "selected Spec missing from the store", key: selectionKey, breach: func(c *publicationChainCase) { delete(c.manifests, secondSpec) }, answer: unreadableAnswer("task_content_missing", secondSpec, "file does not exist")},
		{name: "Spec of another Task", key: selectionKey, breach: stored(secondSpec, func(manifest map[string]canonicaljson.Value) { manifest["task_id"] = "other" }), answer: differs},
		{name: "Spec under another ID", key: selectionKey, breach: stored(secondSpec, func(manifest map[string]canonicaljson.Value) { manifest["spec_id"] = "elsewhere" }), answer: differs},
		{name: "withdrawal binds nothing", key: withdrawalKey, breach: func(c *publicationChainCase) {
			c.registry.TaskSpecAssessments, c.registry.TaskSpecRevisions, c.manifests = nil, nil, nil
		}, answer: chainAccepted},
	})
	t.Run("without storage", func(t *testing.T) {
		c := newPublicationChainCase(t, selectionKey)
		got := chainAnswer(validateContentPublicationChain(nil, c.work, c.registry, c.pub, c.manifest, c.canonical(c.request)))
		if want := unreadableAnswer("task_content_observation_unknown", secondAssessment, "task_content_observation_unknown: content storage dependency is missing"); got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}

// A problem that Ply wrote from the Task's summary is checked against the
// registered Task; its request has no documents to compare.
func TestPublicationChainGeneratedProblemIsTheTaskSummary(t *testing.T) {
	t.Parallel()
	runChainRules(t, []chainRule{
		{name: "Task is not registered", key: legacyImportKey, breach: func(c *publicationChainCase) { c.registry.Tasks = nil }, answer: generatedAnswer},
		{name: "Task renamed since", key: legacyImportKey, breach: func(c *publicationChainCase) { c.task().Title = "Renamed" }, answer: generatedAnswer},
		{name: "summary other than the Task's description", key: legacyImportKey, breach: func(c *publicationChainCase) { c.request["summary"], c.manifest["summary"] = "Other", "Other" }, answer: generatedAnswer},
		{name: "manifest without documents", key: legacyImportKey, breach: func(c *publicationChainCase) { c.manifest["documents"] = []canonicaljson.Value{} }, answer: generatedDocAnswer},
		{name: "manifest with a second document", key: legacyImportKey, breach: addManifestDocument, answer: generatedDocAnswer},
		{name: "document outside the managed store", key: createKey, breach: func(c *publicationChainCase) {
			changeDocument(c.manifest, 0, func(document map[string]canonicaljson.Value) { document["locator"] = "/elsewhere/problem.md" })
		}, answer: generatedDocAnswer},
	})
}

func TestPublicationChainFileDocumentsFollowTheRequest(t *testing.T) {
	t.Parallel()
	document := func(change func(document map[string]canonicaljson.Value)) func(*publicationChainCase) {
		return func(c *publicationChainCase) { changeDocument(c.manifest, 0, change) }
	}
	const metadata = "unclassified: document metadata differs"
	runChainRules(t, []chainRule{
		{name: "manifest without documents", key: problemKey, breach: func(c *publicationChainCase) { c.manifest["documents"] = []canonicaljson.Value{} }, answer: "unclassified: document set differs from request"},
		{name: "document with other content", key: problemKey, breach: document(func(document map[string]canonicaljson.Value) { document["sha256"] = foreignDigest }), answer: metadata},
		{name: "document of another media type", key: problemKey, breach: document(func(document map[string]canonicaljson.Value) { document["media_type"] = "text/plain" }), answer: metadata},
		{name: "document with Git provenance that the request lacks", key: problemKey, breach: document(func(document map[string]canonicaljson.Value) {
			changeObject(document, "provenance", func(provenance map[string]canonicaljson.Value) { provenance["git"] = "tracked" })
		}), answer: "unclassified: document provenance differs"},
	})
}

// A snapshot copies a document of a manifest that its Task, or the Task the
// source names, has registered as a problem, a Spec or an assessment.
func TestPublicationChainSnapshotsFollowTheirSource(t *testing.T) {
	t.Parallel()
	h := sharedPublicationChainHistory(t)
	firstSpec := h.registry.TaskSpecRevisions[chainFirst].ManifestSHA256
	ownedByOther := func(register func(r *WorkItemRegistry)) func(*publicationChainCase) {
		return func(c *publicationChainCase) {
			c.changeSource(0, func(source map[string]canonicaljson.Value) { source["task_id"] = "other" })
			register(&c.registry)
		}
	}
	evidence := func(c *publicationChainCase) { snapshotEvidence(c, "evidence") }
	runChainRules(t, []chainRule{
		{name: "source names its own Task", key: snapshotSpecKey, breach: func(c *publicationChainCase) {
			c.changeSource(0, func(source map[string]canonicaljson.Value) { source["task_id"] = "task" })
		}, answer: chainAccepted},
		{name: "source is a problem of the named Task", key: snapshotSpecKey, breach: ownedByOther(func(r *WorkItemRegistry) {
			r.TaskProblemRevisions = append(r.TaskProblemRevisions, TaskProblemReference{TaskID: "other", Revision: 1, ManifestSHA256: firstSpec})
		}), answer: chainAccepted},
		{name: "source is a Spec of the named Task", key: snapshotSpecKey, breach: ownedByOther(func(r *WorkItemRegistry) {
			r.TaskSpecRevisions = append(r.TaskSpecRevisions, TaskSpecReference{TaskID: "other", SpecID: "solution", Revision: 1, ManifestSHA256: firstSpec})
		}), answer: chainAccepted},
		{name: "source is an assessment of the named Task", key: snapshotSpecKey, breach: ownedByOther(func(r *WorkItemRegistry) {
			r.TaskSpecAssessments = append(r.TaskSpecAssessments, TaskAssessmentReference{TaskID: "other", SpecID: "solution", SpecRevision: 1, Ordinal: 1, ID: "asm_other", ManifestSHA256: firstSpec})
		}), answer: chainAccepted},
		{name: "source is a selection of the named Task", key: snapshotSpecKey, breach: ownedByOther(func(r *WorkItemRegistry) {
			r.TaskSolutionSelections = append(r.TaskSolutionSelections, TaskSelectionReference{TaskID: "other", Ordinal: 1, ID: "sel_other", ManifestSHA256: firstSpec})
		}), answer: unownedSnapshot},
		{name: "source missing from the store", key: snapshotSpecKey, breach: func(c *publicationChainCase) { delete(c.manifests, firstSpec) }, answer: unreadableAnswer("task_content_missing", firstSpec, "file does not exist")},
		{name: "source without the document", key: snapshotSpecKey, breach: func(c *publicationChainCase) {
			c.changeSource(0, func(source map[string]canonicaljson.Value) { source["document_id"] = "absent" })
		}, answer: "unclassified: snapshot source document is missing"},
		{name: "source without the document, and another descriptor", key: snapshotSpecKey, breach: func(c *publicationChainCase) {
			c.changeSource(0, func(source map[string]canonicaljson.Value) { source["document_id"] = "absent" })
			changeDocument(c.manifest, 0, func(document map[string]canonicaljson.Value) { document["size_bytes"] = int64(1) })
		}, answer: "unclassified: snapshot source document is missing"},
		{name: "snapshot under a new document ID", key: nextAssessmentKey, breach: evidence, answer: chainAccepted},
		{name: "snapshot under a new document ID with another descriptor", key: nextAssessmentKey, breach: func(c *publicationChainCase) {
			evidence(c)
			changeDocument(c.manifest, 0, func(document map[string]canonicaljson.Value) { document["media_type"] = "text/plain" })
		}, answer: "unclassified: snapshot descriptor differs"},
	})
}

// A stored manifest that the manifest reader rejects answers with the reader's
// own classified error, wherever the chain reads one.
func TestPublicationChainForwardsRejectedStoredManifests(t *testing.T) {
	t.Parallel()
	h := sharedPublicationChainHistory(t)
	firstSpec, secondSpec := h.registry.TaskSpecRevisions[chainFirst].ManifestSHA256, h.registry.TaskSpecRevisions[chainSecond].ManifestSHA256
	secondAssessment := h.registry.TaskSpecAssessments[chainSecond].ManifestSHA256
	malformed := func(digest string) func(*publicationChainCase) {
		return func(c *publicationChainCase) { c.manifests[digest] = []byte("{") }
	}
	notCanonical := func(digest string) func(*publicationChainCase) {
		return func(c *publicationChainCase) { c.manifests[digest] = append(slices.Clone(c.manifests[digest]), '\n') }
	}
	malformedAnswer := func(digest string) string {
		return rejectedAnswer("task_content_invalid_input", digest, "malformed publication: unterminated JSON object")
	}
	notCanonicalAnswer := func(digest string) string {
		return rejectedAnswer("task_content_integrity_conflict", digest, "manifest is not canonical")
	}
	runChainRules(t, []chainRule{
		{name: "malformed assessment", key: selectionKey, breach: malformed(secondAssessment), answer: malformedAnswer(secondAssessment)},
		{name: "assessment that is not canonical", key: selectionKey, breach: notCanonical(secondAssessment), answer: notCanonicalAnswer(secondAssessment)},
		{name: "malformed selected Spec", key: selectionKey, breach: malformed(secondSpec), answer: malformedAnswer(secondSpec)},
		{name: "selected Spec that is not canonical", key: selectionKey, breach: notCanonical(secondSpec), answer: notCanonicalAnswer(secondSpec)},
		{name: "malformed snapshot source", key: snapshotSpecKey, breach: malformed(firstSpec), answer: malformedAnswer(firstSpec)},
		{name: "snapshot source that is not canonical", key: snapshotSpecKey, breach: notCanonical(firstSpec), answer: notCanonicalAnswer(firstSpec)},
	})
}

// Documents are compared in request order; the first one that differs answers.
func TestPublicationChainFirstDifferingDocumentAnswers(t *testing.T) {
	t.Parallel()
	twoDocuments := func(c *publicationChainCase) { snapshotEvidence(c, "appendix", "evidence") }
	resize := func(c *publicationChainCase, i int) {
		changeDocument(c.manifest, i, func(document map[string]canonicaljson.Value) { document["size_bytes"] = int64(1) })
	}
	rename := func(c *publicationChainCase, i int) {
		changeDocument(c.manifest, i, func(document map[string]canonicaljson.Value) { document["id"] = "renamed" })
	}
	const descriptor, id = "unclassified: snapshot descriptor differs", "unclassified: document ID differs"
	runChainRules(t, []chainRule{
		{name: "two documents", key: nextAssessmentKey, breach: twoDocuments, answer: chainAccepted},
		{name: "second document differs", key: nextAssessmentKey, breach: func(c *publicationChainCase) { twoDocuments(c); resize(c, 1) }, answer: descriptor},
		{name: "first document's descriptor before the second's ID", key: nextAssessmentKey, breach: func(c *publicationChainCase) { twoDocuments(c); resize(c, 0); rename(c, 1) }, answer: descriptor},
		{name: "first document's ID before the second's descriptor", key: nextAssessmentKey, breach: func(c *publicationChainCase) { twoDocuments(c); rename(c, 0); resize(c, 1) }, answer: id},
	})
}

func TestPublicationChainEmptyRegistry(t *testing.T) {
	t.Parallel()
	answers := map[string]string{
		createKey:         "unclassified: missing Task",
		problemKey:        predecessorAnswer,
		specKey:           "unclassified: problem reference is not owned by Task",
		snapshotSpecKey:   predecessorAnswer,
		assessmentKey:     "unclassified: Spec reference is not owned by Task",
		nextAssessmentKey: predecessorAnswer,
		selectionKey:      unregisteredAsm,
		reselectionKey:    predecessorAnswer,
		withdrawalKey:     predecessorAnswer,
		legacyImportKey:   generatedAnswer,
	}
	for key, want := range answers {
		c := newPublicationChainCase(t, key)
		c.registry = WorkItemRegistry{}
		if got := c.answer(); got != want {
			t.Errorf("%s: got %q, want %q", key, got, want)
		}
	}
}
