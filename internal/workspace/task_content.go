package workspace

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

// Revision and decision references always identify immutable canonical manifests.
type TaskRevisionRef struct {
	Revision       int    `yaml:"revision" json:"revision"`
	ManifestSHA256 string `yaml:"manifest_sha256" json:"manifest_sha256"`
}
type TaskDecisionRef struct {
	ID             string `yaml:"id" json:"id"`
	ManifestSHA256 string `yaml:"manifest_sha256" json:"manifest_sha256"`
}
type TaskSpecBasis struct {
	TaskID                TaskID          `yaml:"task_id" json:"task_id"`
	Problem               TaskRevisionRef `yaml:"problem" json:"problem"`
	SpecID                string          `yaml:"spec_id" json:"spec_id"`
	Spec                  TaskRevisionRef `yaml:"spec" json:"spec"`
	Assessment            TaskDecisionRef `yaml:"assessment" json:"assessment"`
	Selection             TaskDecisionRef `yaml:"selection" json:"selection"`
	TaskWorktreeID        WorktreeID      `yaml:"task_worktree_id" json:"task_worktree_id"`
	DeliveryBindingSHA256 string          `yaml:"delivery_binding_sha256" json:"delivery_binding_sha256"`
	Dependencies          []string        `yaml:"dependencies" json:"dependencies"`
}
type TaskSpecPolicy struct {
	TaskID                   TaskID         `yaml:"task_id" json:"task_id"`
	Mode                     string         `yaml:"mode" json:"mode"`
	ActivationPublicationKey *string        `yaml:"activation_publication_key" json:"activation_publication_key"`
	LegacyResultIDs          []TaskResultID `yaml:"legacy_result_ids" json:"legacy_result_ids"`
}
type TaskProblemReference struct {
	TaskID         TaskID `yaml:"task_id" json:"task_id"`
	Revision       int    `yaml:"revision" json:"revision"`
	ManifestSHA256 string `yaml:"manifest_sha256" json:"manifest_sha256"`
}
type TaskSpecReference struct {
	TaskID         TaskID `yaml:"task_id" json:"task_id"`
	SpecID         string `yaml:"spec_id" json:"spec_id"`
	Revision       int    `yaml:"revision" json:"revision"`
	ManifestSHA256 string `yaml:"manifest_sha256" json:"manifest_sha256"`
}
type TaskAssessmentReference struct {
	TaskID         TaskID `yaml:"task_id" json:"task_id"`
	SpecID         string `yaml:"spec_id" json:"spec_id"`
	SpecRevision   int    `yaml:"spec_revision" json:"spec_revision"`
	Ordinal        int    `yaml:"ordinal" json:"ordinal"`
	ID             string `yaml:"id" json:"id"`
	ManifestSHA256 string `yaml:"manifest_sha256" json:"manifest_sha256"`
}
type TaskSelectionReference struct {
	TaskID         TaskID `yaml:"task_id" json:"task_id"`
	Ordinal        int    `yaml:"ordinal" json:"ordinal"`
	ID             string `yaml:"id" json:"id"`
	ManifestSHA256 string `yaml:"manifest_sha256" json:"manifest_sha256"`
}
type TaskResultSpecBinding struct {
	TaskResultID       TaskResultID  `yaml:"task_result_id" json:"task_result_id"`
	TaskID             TaskID        `yaml:"task_id" json:"task_id"`
	Basis              TaskSpecBasis `yaml:"basis" json:"basis"`
	HandoffSHA256      string        `yaml:"handoff_sha256" json:"handoff_sha256"`
	StartReceiptSHA256 string        `yaml:"start_receipt_sha256" json:"start_receipt_sha256"`
}
type TaskContentOutcomeRef struct {
	Kind           string  `yaml:"kind" json:"kind"`
	TaskID         TaskID  `yaml:"task_id" json:"task_id"`
	SpecID         *string `yaml:"spec_id" json:"spec_id"`
	Revision       *int    `yaml:"revision" json:"revision"`
	ID             *string `yaml:"id" json:"id"`
	ManifestSHA256 string  `yaml:"manifest_sha256" json:"manifest_sha256"`
}
type TaskContentPublication struct {
	PublicationKey       string                `yaml:"publication_key" json:"publication_key"`
	TaskID               TaskID                `yaml:"task_id" json:"task_id"`
	Operation            string                `yaml:"operation" json:"operation"`
	IntentSHA256         string                `yaml:"intent_sha256" json:"intent_sha256"`
	RequestSHA256        string                `yaml:"request_sha256" json:"request_sha256"`
	OutcomeRef           TaskContentOutcomeRef `yaml:"outcome_ref" json:"outcome_ref"`
	StoreTransition      string                `yaml:"store_transition" json:"store_transition"`
	RegistryBeforeSHA256 *string               `yaml:"registry_before_sha256" json:"registry_before_sha256"`
	BackupSHA256         *string               `yaml:"backup_sha256" json:"backup_sha256"`
	RecordedAtUTC        string                `yaml:"recorded_at_utc" json:"recorded_at_utc"`
}
type TaskRegistryUpgrade struct {
	FromVersion    int    `yaml:"from_version" json:"from_version"`
	RegistrySHA256 string `yaml:"registry_sha256" json:"registry_sha256"`
}
type TaskContentInput struct {
	TaskID TaskID
	File   string
}
type TaskContentQuery struct {
	TaskID         TaskID
	SpecID         string
	Revision       int
	PublicationKey string
}
type TaskContentState struct {
	RegistrySHA256 *string          `json:"registry_sha256"`
	ProblemHead    *TaskRevisionRef `json:"problem_head"`
	SpecHead       *TaskRevisionRef `json:"spec_head"`
	SelectionEvent *TaskDecisionRef `json:"selection_event"`
	Freshness      string           `json:"freshness"`
}
type TaskContentAttempt struct {
	ObjectPublicationStarted bool   `json:"object_publication_started"`
	RegistryReplace          string `json:"registry_replace"`
}
type TaskNextAction struct {
	Kind   string   `json:"kind"`
	Reason string   `json:"reason"`
	Argv   []string `json:"argv"`
}
type TaskContentMutationResult struct {
	Workspace      string                 `json:"workspace"`
	TaskID         TaskID                 `json:"task_id"`
	PublicationKey string                 `json:"publication_key"`
	IntentSHA256   string                 `json:"intent_sha256"`
	Classification string                 `json:"classification"`
	PreState       TaskContentState       `json:"pre_state"`
	Attempted      TaskContentAttempt     `json:"attempted"`
	PostState      TaskContentState       `json:"post_state"`
	Durability     string                 `json:"durability"`
	OutcomeRef     *TaskContentOutcomeRef `json:"outcome_ref"`
	NextAction     TaskNextAction         `json:"next_action"`
}

type TaskContentIDSource interface {
	NewAssessmentID() (string, error)
	NewSelectionID() (string, error)
}
type cryptoTaskContentIDs struct{ Reader io.Reader }

func (s cryptoTaskContentIDs) next(prefix string) (string, error) {
	r := s.Reader
	if r == nil {
		r = rand.Reader
	}
	return randomWorkItemID(r, prefix)
}
func (s cryptoTaskContentIDs) NewAssessmentID() (string, error) { return s.next("asm_") }
func (s cryptoTaskContentIDs) NewSelectionID() (string, error)  { return s.next("sel_") }

type TaskContentError struct {
	Code, Detail string
	Err          error
}

func (e *TaskContentError) Error() string               { return e.Code + ": " + e.Detail }
func (e *TaskContentError) Unwrap() error               { return e.Err }
func contentError(code, detail string, err error) error { return &TaskContentError{code, detail, err} }

func contentValue(value any) (canonicaljson.Value, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return canonicaljson.DecodeStrict(b)
}
func contentCanonical(value any) ([]byte, error) {
	v, err := contentValue(value)
	if err != nil {
		return nil, err
	}
	return canonicaljson.Marshal(v)
}
func contentDecode(value canonicaljson.Value, out any) error {
	b, err := canonicaljson.Marshal(value)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}
func contentFields(value canonicaljson.Value) map[string]canonicaljson.Value {
	m := map[string]canonicaljson.Value{}
	if o, ok := value.(canonicaljson.Object); ok {
		for _, p := range o {
			m[p.Name] = p.Value
		}
	}
	return m
}
func contentObject(m map[string]canonicaljson.Value) canonicaljson.Object {
	o := canonicaljson.Object{}
	for k, v := range m {
		o = append(o, canonicaljson.Member{Name: k, Value: v})
	}
	sort.Slice(o, func(i, j int) bool { return o[i].Name < o[j].Name })
	return o
}
func contentString(m map[string]canonicaljson.Value, key string) string {
	s, _ := m[key].(string)
	return s
}
func contentInt(m map[string]canonicaljson.Value, key string) int {
	n, _ := m[key].(int64)
	return int(n)
}
func contentArray(m map[string]canonicaljson.Value, key string) []canonicaljson.Value {
	a, _ := m[key].([]canonicaljson.Value)
	return a
}
func contentEqual(a, b canonicaljson.Value) bool {
	aa, e := canonicaljson.Marshal(a)
	bb, f := canonicaljson.Marshal(b)
	return e == nil && f == nil && string(aa) == string(bb)
}
func contentTypedEqual(a, b any) bool {
	aa, e := contentCanonical(a)
	bb, f := contentCanonical(b)
	return e == nil && f == nil && string(aa) == string(bb)
}
func contentEnvelope(kind string, m map[string]canonicaljson.Value) canonicaljson.Object {
	m["kind"] = kind
	m["schema_version"] = int64(1)
	m["format"] = "json"
	m["format_version"] = int64(1)
	m["canonicalization"] = "RFC8785"
	return contentObject(m)
}
func taskContentState(r WorkItemRegistry, id TaskID) TaskContentState {
	s := TaskContentState{Freshness: "fresh"}
	if r.RawSHA256 != "" {
		s.RegistrySHA256 = &r.RawSHA256
	}
	for _, p := range r.TaskProblemRevisions {
		if p.TaskID == id {
			s.ProblemHead = &TaskRevisionRef{p.Revision, p.ManifestSHA256}
		}
	}
	for _, p := range r.TaskSpecRevisions {
		if p.TaskID == id {
			s.SpecHead = &TaskRevisionRef{p.Revision, p.ManifestSHA256}
		}
	}
	for _, p := range r.TaskSolutionSelections {
		if p.TaskID == id {
			s.SelectionEvent = &TaskDecisionRef{p.ID, p.ManifestSHA256}
		}
	}
	return s
}
func taskRequiresSpec(r WorkItemRegistry, id TaskID) bool {
	for _, p := range r.TaskSpecPolicies {
		if p.TaskID == id {
			return p.Mode == "spec_required"
		}
	}
	return false
}
func taskPublication(r WorkItemRegistry, key string) *TaskContentPublication {
	for _, p := range r.TaskContentPublications {
		if p.PublicationKey == key {
			return &p
		}
	}
	return nil
}
func contentInspectAction(id TaskID, key string) TaskNextAction {
	return TaskNextAction{"inspect_publication", "Inspect the publication before retrying.", []string{"ply", "workspace", "task", "publication", "show", string(id), "--key", key}}
}

func RecordTaskProblem(d Dependencies, in TaskContentInput) (TaskContentMutationResult, error) {
	return mutateTaskContent(d, in, "problem_record")
}
func RecordTaskSpec(d Dependencies, in TaskContentInput) (TaskContentMutationResult, error) {
	return mutateTaskContent(d, in, "spec_record")
}
func AssessTaskSpec(d Dependencies, in TaskContentInput) (TaskContentMutationResult, error) {
	return mutateTaskContent(d, in, "spec_assess")
}
func SelectTaskSolution(d Dependencies, in TaskContentInput) (TaskContentMutationResult, error) {
	return mutateTaskContent(d, in, "spec_select")
}
func WithdrawTaskSolution(d Dependencies, in TaskContentInput) (TaskContentMutationResult, error) {
	return mutateTaskContent(d, in, "spec_withdraw")
}

func contentConflict(code string, expected, actual any) error {
	return contentError(code, fmt.Sprintf("expected predecessor %v; actual head %v; inspect the current Task before proposing a new revision", expected, actual), nil)
}

func mutateTaskContent(d Dependencies, in TaskContentInput, operation string) (TaskContentMutationResult, error) {
	if _, e := ParseTaskID(string(in.TaskID)); e != nil {
		return TaskContentMutationResult{}, e
	}
	if d.TaskContent == nil || d.TaskContentIDs == nil || d.WorkClock == nil {
		return TaskContentMutationResult{}, contentError("task_content_observation_unknown", "content storage, IDs and clock dependencies are required", nil)
	}
	raw, e := d.TaskContent.ReadSource(in.File, taskManifestLimit)
	if e != nil {
		return TaskContentMutationResult{}, e
	}
	draft, e := decodeTaskContent(raw, operation, false, false)
	if e != nil {
		return TaskContentMutationResult{}, e
	}
	m := contentFields(draft)
	if contentString(m, "task_id") != string(in.TaskID) {
		return TaskContentMutationResult{}, contentError("task_content_invalid_input", "draft task_id differs from command Task", nil)
	}
	if e = requireWorkMutationDependencies(d); e != nil {
		return TaskContentMutationResult{}, e
	}
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return TaskContentMutationResult{}, e
	}
	b, e := canonicaljson.Marshal(draft)
	if e != nil {
		return TaskContentMutationResult{}, e
	}
	result := newTaskContentMutation(root, in.TaskID, contentString(m, "publication_key"), digestTaskBytes(b))
	e = d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			r, err := session.Snapshot()
			if err != nil {
				return err
			}
			result.PreState = taskContentState(r, in.TaskID)
			result.PostState = result.PreState
			task, _ := findTask(r, in.TaskID)
			if task == nil {
				return workError(ErrorWorkNotFound, "Task is not registered", nil)
			}
			if _, _, err = projectAndRepo(projects, task.ProjectID, task.RepoID); err != nil {
				return err
			}
			return publishTaskContent(d, root, session, r, *task, projects, draft, operation, b, &result, false)
		})
	})
	return finishTaskContentMutation(d, result, e)
}
func newTaskContentMutation(root string, id TaskID, key, digest string) TaskContentMutationResult {
	return TaskContentMutationResult{Workspace: root, TaskID: id, PublicationKey: key, IntentSHA256: digest, Classification: "not_published", Durability: "not_observed", PreState: TaskContentState{Freshness: "unknown"}, PostState: TaskContentState{Freshness: "unknown"}, Attempted: TaskContentAttempt{RegistryReplace: "not_attempted"}, NextAction: contentInspectAction(id, key)}
}
func finishTaskContentMutation(d Dependencies, result TaskContentMutationResult, err error) (TaskContentMutationResult, error) {
	r, e := d.WorkItems.Snapshot(result.Workspace)
	if e == nil {
		result.PostState = taskContentState(r, result.TaskID)
	} else {
		result.PostState.Freshness = "unknown"
		if err == nil {
			err = e
		}
	}
	if err != nil {
		if result.Attempted.RegistryReplace != "not_attempted" || result.Classification == "existing" {
			result.Classification = "unknown"
			result.Durability = "unknown"
			return result, contentError("task_content_publication_unknown", fmt.Sprintf("Task %s publication %s may have changed; inspect with ply workspace task publication show %s --key %s: %v", result.TaskID, result.PublicationKey, result.TaskID, result.PublicationKey, err), err)
		}
		if strings.Contains(err.Error(), "conflict") {
			result.Classification = "conflict"
		}
		return result, err
	}
	result.Durability = "confirmed"
	result.NextAction = TaskNextAction{"inspect_task", "Inspect the preserved Task basis.", []string{"ply", "workspace", "task", "show", string(result.TaskID)}}
	return result, nil
}
func syncTaskRegistry(root string) error {
	f, e := os.Open(workItemsPath(root))
	if e != nil {
		return e
	}
	e = f.Sync()
	c := f.Close()
	if e = firstError(e, c); e != nil {
		return e
	}
	return platformSyncDirectory(filepathJoinMarker(root))
}
func filepathJoinMarker(root string) string { return root + string(os.PathSeparator) + MarkerDirectory }

func publishTaskContent(d Dependencies, root string, session WorkItemStoreSession, r WorkItemRegistry, task TaskRecord, projects ProjectSnapshot, draft canonicaljson.Object, operation string, request []byte, result *TaskContentMutationResult, creating bool) error {
	m := contentFields(draft)
	key := result.PublicationKey
	if existing := taskPublication(r, key); existing != nil {
		if existing.TaskID != task.ID || existing.IntentSHA256 != result.IntentSHA256 {
			return contentError("task_content_publication_conflict", "publication key already identifies different input", nil)
		}
		result.Classification = "existing"
		result.OutcomeRef = &existing.OutcomeRef
		if e := validateTaskContentClosure(d.TaskContent, root, r, true); e != nil {
			return e
		}
		return syncTaskRegistry(root)
	}
	if e := validateTaskContentClosure(d.TaskContent, root, r, false); e != nil {
		return e
	}
	var upgrade *TaskRegistryUpgrade
	if m["registry_upgrade"] != nil {
		if e := contentDecode(m["registry_upgrade"], &upgrade); e != nil {
			return e
		}
	}
	transition := "none"
	var backup []byte
	var backupHash *string
	if r.FormatVersion < 3 {
		if !creating && operation != "problem_record" {
			return contentError("task_spec_required", "record a problem before recording a solution", nil)
		}
		if upgrade == nil {
			return contentError("task_content_upgrade_required", "acknowledge registry "+r.RawSHA256+" before migrating format 1 or 2", nil)
		}
		if upgrade.FromVersion != r.FormatVersion || upgrade.RegistrySHA256 != r.RawSHA256 {
			return contentError("task_content_upgrade_conflict", "upgrade acknowledgement does not match current registry", nil)
		}
		var e error
		backup, e = os.ReadFile(workItemsPath(root))
		if e != nil {
			return e
		}
		if digestTaskBytes(backup) != r.RawSHA256 {
			return contentError("task_content_upgrade_conflict", "registry changed before migration", nil)
		}
		transition = fmt.Sprintf("format_%d_to_3", r.FormatVersion)
		value := r.RawSHA256
		backupHash = &value
		upgradeTaskContentRegistry(&r)
	} else if upgrade != nil {
		return contentError("task_content_upgrade_conflict", "registry is already format 3", nil)
	}
	if creating {
		r.Tasks = append(r.Tasks, task)
	}
	state := taskContentState(r, task.ID)
	var previous canonicaljson.Value
	revision, ordinal := 0, 0
	id := ""
	specID := contentString(m, "spec_id")
	switch operation {
	case "problem_record":
		if !contentTypedEqual(state.ProblemHead, valueRevision(m["expected_previous"])) {
			return contentConflict("task_problem_revision_conflict", m["expected_previous"], state.ProblemHead)
		}
		if state.ProblemHead != nil {
			revision = state.ProblemHead.Revision
		}
		if revision == 2147483647 {
			return contentError("task_content_invalid_input", "problem revision exceeds 2147483647", nil)
		}
		revision++
		previous = m["expected_previous"]
		origin := contentFields(m["origin"])
		if contentString(origin, "kind") == "legacy_summary_import" {
			if state.ProblemHead != nil || taskRequiresSpec(r, task.ID) || contentString(m, "title") != task.Title || contentString(m, "summary") != task.Description {
				return contentError("task_content_invalid_input", "legacy import must match the original Task summary", nil)
			}
			b, _ := contentCanonical(map[string]any{"task_id": task.ID, "title": task.Title, "description": task.Description})
			if contentString(origin, "legacy_summary_sha256") != digestTaskBytes(b) {
				return contentError("task_content_integrity_conflict", "legacy summary digest differs", nil)
			}
		}
		if state.ProblemHead == nil && !creating {
			for _, a := range r.IntegrationAuthorities {
				if a.TaskID == task.ID && !taskIntegrationAuthorityComplete(r, a.ID) {
					return contentError("task_spec_binding_conflict", "cannot activate a Task with an unresolved integration authority", nil)
				}
			}
		}
	case "spec_record":
		if !taskRequiresSpec(r, task.ID) || state.ProblemHead == nil {
			return contentError("task_spec_required", "record a Task problem first", nil)
		}
		if !contentTypedEqual(state.ProblemHead, valueRevision(m["problem"])) {
			return contentError("task_spec_basis_stale", "solution must bind the current problem", nil)
		}
		for _, s := range r.TaskSpecRevisions {
			if s.TaskID == task.ID && s.SpecID != specID {
				return contentError("task_spec_alternative_not_supported", "Task already owns a different Spec ID", nil)
			}
		}
		if !contentTypedEqual(state.SpecHead, valueRevision(m["expected_previous"])) {
			return contentConflict("task_spec_revision_conflict", m["expected_previous"], state.SpecHead)
		}
		if state.SpecHead != nil {
			revision = state.SpecHead.Revision
		}
		if revision == 2147483647 {
			return contentError("task_content_invalid_input", "Spec revision exceeds 2147483647", nil)
		}
		revision++
		previous = m["expected_previous"]
		removed := []string{}
		if state.SpecHead != nil {
			old, e := readContentManifest(d.TaskContent, root, state.SpecHead.ManifestSHA256)
			if e != nil {
				return e
			}
			now := map[string]bool{}
			for _, v := range contentArray(m, "requirements") {
				now[contentIDKey(v)] = true
			}
			for _, v := range contentArray(contentFields(old), "requirements") {
				if !now[contentIDKey(v)] {
					removed = append(removed, contentIDKey(v))
				}
			}
		}
		if !contentEqual(sortedContentStrings(removed), m["removed_requirement_ids"]) {
			return contentError("task_content_invalid_input", "removed_requirement_ids must exactly describe removed predecessor requirements", nil)
		}
		if e := validateSpecImplementationBasis(d, projects, r, task, m["implementation_basis"], r.FormatVersion >= 4); e != nil {
			return e
		}
	case "spec_assess":
		spec, e := registeredTaskSpec(d, root, r, task.ID, specID, valueRevision(m["spec"]))
		if e != nil {
			return e
		}
		last := lastTaskAssessment(r, task.ID, specID, contentInt(contentFields(m["spec"]), "revision"))
		var expected *TaskDecisionRef
		if last != nil {
			expected = &TaskDecisionRef{last.ID, last.ManifestSHA256}
			ordinal = last.Ordinal
		}
		if !contentTypedEqual(expected, valueDecision(m["expected_previous_assessment"])) {
			return contentConflict("task_spec_assessment_conflict", m["expected_previous_assessment"], expected)
		}
		previous = m["expected_previous_assessment"]
		if contentString(m, "outcome") == "ready" {
			if !contentTypedEqual(state.ProblemHead, valueRevision(contentFields(spec)["problem"])) {
				return contentError("task_spec_basis_stale", "ready assessment must bind the current problem", nil)
			}
			if e = taskSpecStructurallyReady(spec); e != nil {
				return e
			}
			if e = validateSpecImplementationBasis(d, projects, r, task, contentFields(spec)["implementation_basis"], true); e != nil {
				return e
			}
		}
		id, e = d.TaskContentIDs.NewAssessmentID()
		if e != nil {
			return e
		}
	case "spec_select", "spec_withdraw":
		if !contentTypedEqual(state.SelectionEvent, valueDecision(m["expected_previous_selection"])) {
			return contentConflict("task_spec_selection_conflict", m["expected_previous_selection"], state.SelectionEvent)
		}
		previous = m["expected_previous_selection"]
		for _, s := range r.TaskSolutionSelections {
			if s.TaskID == task.ID {
				ordinal = s.Ordinal
			}
		}
		if operation == "spec_withdraw" {
			if state.SelectionEvent == nil {
				return contentError("task_spec_selection_missing", "there is no selection to withdraw", nil)
			}
			last, e := readContentManifest(d.TaskContent, root, state.SelectionEvent.ManifestSHA256)
			if e != nil {
				return e
			}
			if contentString(contentFields(last), "action") != "select" {
				return contentError("task_spec_selection_conflict", "the latest event already withdrew the selection", nil)
			}
		} else {
			sol := contentFields(m["solution"])
			specID = contentString(sol, "spec_id")
			spec, e := registeredTaskSpec(d, root, r, task.ID, specID, valueRevision(sol["spec"]))
			if e != nil {
				return e
			}
			if !contentTypedEqual(state.ProblemHead, valueRevision(sol["problem"])) || !contentEqual(sol["problem"], contentFields(spec)["problem"]) {
				return contentError("task_spec_selection_stale", "selection must bind the current problem and matching Spec", nil)
			}
			last := lastTaskAssessment(r, task.ID, specID, valueRevision(sol["spec"]).Revision)
			if last == nil || !contentTypedEqual(&TaskDecisionRef{last.ID, last.ManifestSHA256}, valueDecision(sol["assessment"])) {
				return contentError("task_spec_assessment_conflict", "selection must bind the latest assessment of this revision", nil)
			}
			assessment, e := readContentManifest(d.TaskContent, root, last.ManifestSHA256)
			if e != nil {
				return e
			}
			if contentString(contentFields(assessment), "outcome") != "ready" {
				return contentError("task_spec_not_ready", "the latest assessment is not ready", nil)
			}
			if e = taskSpecStructurallyReady(spec); e != nil {
				return e
			}
			if e = validateSpecImplementationBasis(d, projects, r, task, contentFields(spec)["implementation_basis"], true); e != nil {
				return e
			}
		}
		var e error
		id, e = d.TaskContentIDs.NewSelectionID()
		if e != nil {
			return e
		}
	}
	if ordinal == 2147483647 {
		return contentError("task_content_invalid_input", "decision ordinal exceeds 2147483647", nil)
	}
	ordinal++
	docs := []canonicaljson.Value{}
	objects := map[string][]byte{}
	if operation == "problem_record" && contentString(contentFields(m["origin"]), "kind") != "authored" {
		b := []byte("# " + task.Title + "\n\n" + task.Description + "\n")
		digest := digestTaskBytes(b)
		objects[digest] = b
		docs = append(docs, contentObject(map[string]canonicaljson.Value{"id": "problem", "locator": taskContentPath(root, "objects", digest), "sha256": digest, "size_bytes": int64(len(b)), "media_type": "text/markdown", "provenance": contentObject(map[string]canonicaljson.Value{"origin_locator": nil, "git": nil})}))
	} else {
		var e error
		docs, objects, e = contentResolveDocuments(d, root, r, task.ID, contentArray(m, "documents"))
		if e != nil {
			return e
		}
	}
	manifest := contentFields(draft)
	delete(manifest, "registry_upgrade")
	delete(manifest, "expected_previous")
	delete(manifest, "expected_previous_assessment")
	delete(manifest, "expected_previous_selection")
	manifest["previous"] = previous
	manifest["source_draft_sha256"] = result.IntentSHA256
	manifest["recorded_at_utc"] = d.WorkClock.Now().UTC().Format(time.RFC3339Nano)
	outcome := TaskContentOutcomeRef{TaskID: task.ID}
	if operation == "problem_record" || operation == "spec_record" {
		manifest["revision"] = int64(revision)
		outcome.Revision = &revision
	} else {
		manifest["ordinal"] = int64(ordinal)
		manifest["id"] = id
		outcome.ID = &id
	}
	switch operation {
	case "problem_record":
		manifest["kind"] = "WorkspaceTaskProblemRevision@1"
		outcome.Kind = "problem"
		manifest["documents"] = docs
	case "spec_record":
		manifest["kind"] = "WorkspaceTaskSpecRevision@1"
		outcome.Kind = "spec"
		outcome.SpecID = &specID
		manifest["documents"] = docs
	case "spec_assess":
		manifest["kind"] = "WorkspaceTaskSpecAssessment@1"
		outcome.Kind = "assessment"
		outcome.SpecID = &specID
		v := valueRevision(m["spec"]).Revision
		outcome.Revision = &v
		manifest["documents"] = docs
	default:
		manifest["kind"] = "WorkspaceTaskSolutionSelection@1"
		outcome.Kind = "selection"
	}
	manifestBytes, e := canonicaljson.Marshal(contentObject(manifest))
	if e != nil {
		return e
	}
	if _, e = decodeTaskContent(manifestBytes, operation, true, true); e != nil {
		return e
	}
	outcome.ManifestSHA256 = digestTaskBytes(manifestBytes)
	pub := TaskContentPublication{PublicationKey: key, TaskID: task.ID, Operation: operation, IntentSHA256: result.IntentSHA256, RequestSHA256: digestTaskBytes(request), OutcomeRef: outcome, StoreTransition: transition, RegistryBeforeSHA256: result.PreState.RegistrySHA256, BackupSHA256: backupHash, RecordedAtUTC: contentString(manifest, "recorded_at_utc")}
	if creating {
		pub.Operation = "task_create"
	}
	if outcome.Kind == "problem" {
		r.TaskProblemRevisions = append(r.TaskProblemRevisions, TaskProblemReference{task.ID, revision, outcome.ManifestSHA256})
		if revision == 1 {
			activateTaskSpecPolicy(&r, task.ID, key)
		}
	}
	if outcome.Kind == "spec" {
		r.TaskSpecRevisions = append(r.TaskSpecRevisions, TaskSpecReference{task.ID, specID, revision, outcome.ManifestSHA256})
	}
	if outcome.Kind == "assessment" {
		r.TaskSpecAssessments = append(r.TaskSpecAssessments, TaskAssessmentReference{task.ID, specID, *outcome.Revision, ordinal, id, outcome.ManifestSHA256})
	}
	if outcome.Kind == "selection" {
		r.TaskSolutionSelections = append(r.TaskSolutionSelections, TaskSelectionReference{task.ID, ordinal, id, outcome.ManifestSHA256})
	}
	r.TaskContentPublications = append(r.TaskContentPublications, pub)
	sortWorkRegistry(&r)
	if e = validateWorkItemRegistry(r); e != nil {
		return e
	}
	result.Attempted.ObjectPublicationStarted = true
	if backup != nil {
		if _, e = d.TaskContent.Publish(root, "backups", backup); e != nil {
			return e
		}
	}
	if _, e = d.TaskContent.Publish(root, "requests", request); e != nil {
		return e
	}
	keys := make([]string, 0, len(objects))
	for k := range objects {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if _, e = d.TaskContent.Publish(root, "objects", objects[k]); e != nil {
			return e
		}
	}
	if _, e = d.TaskContent.Publish(root, "manifests", manifestBytes); e != nil {
		return e
	}
	if e = d.TaskContent.check("before-registry-replace"); e != nil {
		return e
	}
	result.Attempted.RegistryReplace = "attempted"
	result.OutcomeRef = &outcome
	if e = session.Publish(r); e != nil {
		result.Attempted.RegistryReplace = "may_have_completed"
		return e
	}
	result.Attempted.RegistryReplace = "completed"
	result.Classification = "published"
	if e = d.TaskContent.check("postread"); e != nil {
		return e
	}
	observed, e := session.Snapshot()
	if e != nil {
		return e
	}
	result.PostState = taskContentState(observed, task.ID)
	actual := taskPublication(observed, key)
	if actual == nil || !contentTypedEqual(*actual, pub) {
		return contentError("task_content_publication_unknown", "published registry was not observed", nil)
	}
	if e = validateTaskContentClosure(d.TaskContent, root, observed, false); e != nil {
		return e
	}
	return platformSyncDirectory(filepathJoinMarker(root))
}
func valueRevision(v canonicaljson.Value) *TaskRevisionRef {
	if v == nil {
		return nil
	}
	var out TaskRevisionRef
	if contentDecode(v, &out) != nil {
		return nil
	}
	return &out
}
func valueDecision(v canonicaljson.Value) *TaskDecisionRef {
	if v == nil {
		return nil
	}
	var out TaskDecisionRef
	if contentDecode(v, &out) != nil {
		return nil
	}
	return &out
}
func lastTaskAssessment(r WorkItemRegistry, id TaskID, spec string, revision int) *TaskAssessmentReference {
	var last *TaskAssessmentReference
	for _, a := range r.TaskSpecAssessments {
		if a.TaskID == id && a.SpecID == spec && a.SpecRevision == revision {
			v := a
			last = &v
		}
	}
	return last
}
func registeredTaskSpec(d Dependencies, root string, r WorkItemRegistry, id TaskID, spec string, ref *TaskRevisionRef) (canonicaljson.Object, error) {
	if ref != nil {
		for _, s := range r.TaskSpecRevisions {
			if s.TaskID == id && s.SpecID == spec && s.Revision == ref.Revision && s.ManifestSHA256 == ref.ManifestSHA256 {
				return readContentManifest(d.TaskContent, root, s.ManifestSHA256)
			}
		}
	}
	return nil, contentError("task_spec_binding_conflict", "Spec reference is not registered to this Task", nil)
}
func taskIntegrationAuthorityComplete(r WorkItemRegistry, id IntegrationAuthorityID) bool {
	for _, v := range r.IntegrationResults {
		if v.AuthorityID == id {
			return v.Outcome == "exact_effect" || v.Outcome == "already_integrated" || v.Outcome == "no_effect" || v.Outcome == "conflict"
		}
	}
	return false
}
func upgradeTaskContentRegistry(r *WorkItemRegistry) {
	r.FormatVersion = 3
	if r.TaskResults == nil {
		r.TaskResults = []TaskResultRecord{}
		r.HumanQARecords = []TaskHumanQARecord{}
		r.IntegrationAuthorities = []IntegrationAuthority{}
		r.IntegrationIntents = []IntegrationIntent{}
		r.IntegrationAttempts = []IntegrationAttempt{}
		r.IntegrationResults = []IntegrationResult{}
	}
	r.TaskSpecPolicies = []TaskSpecPolicy{}
	r.TaskProblemRevisions = []TaskProblemReference{}
	r.TaskSpecRevisions = []TaskSpecReference{}
	r.TaskSpecAssessments = []TaskAssessmentReference{}
	r.TaskSolutionSelections = []TaskSelectionReference{}
	r.TaskResultSpecBindings = []TaskResultSpecBinding{}
	r.TaskContentPublications = []TaskContentPublication{}
	for _, t := range r.Tasks {
		r.TaskSpecPolicies = append(r.TaskSpecPolicies, TaskSpecPolicy{TaskID: t.ID, Mode: "legacy", LegacyResultIDs: []TaskResultID{}})
	}
}
func activateTaskSpecPolicy(r *WorkItemRegistry, id TaskID, key string) {
	p := TaskSpecPolicy{TaskID: id, Mode: "spec_required", ActivationPublicationKey: &key, LegacyResultIDs: []TaskResultID{}}
	for _, v := range r.TaskResults {
		if v.TaskID == id {
			p.LegacyResultIDs = append(p.LegacyResultIDs, v.ID)
		}
	}
	for i, v := range r.TaskSpecPolicies {
		if v.TaskID == id {
			r.TaskSpecPolicies[i] = p
			return
		}
	}
	r.TaskSpecPolicies = append(r.TaskSpecPolicies, p)
}

func createTaskProblem(d Dependencies, root string, session WorkItemStoreSession, r WorkItemRegistry, projects ProjectSnapshot, task TaskRecord, upgrade *TaskRegistryUpgrade) (TaskContentMutationResult, error) {
	if d.TaskContent == nil || d.WorkClock == nil {
		return TaskContentMutationResult{}, contentError("task_content_observation_unknown", "content storage and clock dependencies are required", nil)
	}
	if upgrade != nil {
		copy := *upgrade
		upgrade = &copy
		if upgrade.FromVersion == 0 {
			upgrade.FromVersion = r.FormatVersion
		}
		if upgrade.FromVersion < 1 || upgrade.FromVersion > 2 || !digestPattern.MatchString(upgrade.RegistrySHA256) {
			return TaskContentMutationResult{}, contentError("task_content_invalid_input", "invalid store upgrade acknowledgement", nil)
		}
	}
	now := d.WorkClock.Now().UTC().Format(time.RFC3339Nano)
	u, e := contentValue(upgrade)
	if e != nil {
		return TaskContentMutationResult{}, e
	}
	intent, e := contentCanonical(map[string]any{"task_id": task.ID, "title": task.Title, "description": task.Description, "parent_epic_id": task.ParentEpicID, "project_id": task.ProjectID, "repo_id": task.RepoID, "registry_upgrade": upgrade})
	if e != nil {
		return TaskContentMutationResult{}, e
	}
	draft := contentEnvelope("WorkspaceTaskProblemDraft@1", map[string]canonicaljson.Value{"publication_key": "task-create/" + string(task.ID), "task_id": string(task.ID), "expected_previous": nil, "registry_upgrade": u, "origin": contentObject(map[string]canonicaljson.Value{"kind": "task_create_summary"}), "title": task.Title, "summary": task.Description, "problem_document_id": "problem", "documents": []canonicaljson.Value{}, "sources": []canonicaljson.Value{}, "claims": []canonicaljson.Value{}, "deadline": nil, "change_reason": "Record the initial Task problem.", "recorder": contentObject(map[string]canonicaljson.Value{"actor_claim": "local_cli_user", "control_surface": "ply workspace task create", "recorded_at_utc": now})})
	request, e := canonicaljson.Marshal(draft)
	if e != nil {
		return TaskContentMutationResult{}, e
	}
	result := newTaskContentMutation(root, task.ID, "task-create/"+string(task.ID), digestTaskBytes(intent))
	result.PreState = taskContentState(r, task.ID)
	result.PostState = result.PreState
	e = publishTaskContent(d, root, session, r, task, projects, draft, "problem_record", request, &result, true)
	return result, e
}

// MarshalTaskSpecValue provides the canonical bridge without reversing package dependencies.
func MarshalTaskSpecValue(value any) ([]byte, error) { return contentCanonical(value) }
