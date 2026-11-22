package taskwork

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/filelock"
	"github.com/devdimensionlab/plybuild/internal/workitem"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"gopkg.in/yaml.v3"
)

// WorkItemStore reads and writes the work-item registry of one Project.
type WorkItemStore interface {
	Snapshot(workspace.ProjectWork) (WorkItemRegistry, error)
	SnapshotRegistrations(workspace.ProjectWork) (WorkItemRegistry, error)
	WithLock(workspace.ProjectWork, func(WorkItemStoreSession) error) error
}

type WorkItemStoreSession interface {
	Snapshot() (WorkItemRegistry, error)
	Publish(WorkItemRegistry) error
}

type workItemStoreFaults struct {
	fail       func(string) error
	shortWrite bool
}

type systemWorkItemStore struct {
	faults       *workItemStoreFaults
	scope        *taskContentScope
	contentScope *taskContentScope
}
type systemWorkItemStoreSession struct {
	store             *systemWorkItemStore
	work              workspace.ProjectWork
	confirmed         *WorkItemRegistry
	observed          os.FileInfo
	closed            bool
	definitionsDigest string
}

func newSystemWorkItemStore() *systemWorkItemStore { return &systemWorkItemStore{} }
func emptyWorkItemRegistry() WorkItemRegistry {
	return WorkItemRegistry{Epics: []EpicRecord{}, Tasks: []TaskRecord{}, WorktreeOperations: []WorktreeOperationRecord{}, TaskResults: []HistoricalRecord{}, HumanQARecords: []HistoricalRecord{}, IntegrationAuthorities: []IntegrationAuthority{}, IntegrationIntents: []IntegrationIntent{}, IntegrationAttempts: []IntegrationAttempt{}, IntegrationResults: []IntegrationResult{}, TaskProblemRevisions: []TaskProblemReference{}, TaskSpecRevisions: []TaskSpecReference{}, TaskSpecAssessments: []TaskAssessmentReference{}, TaskSolutionSelections: []TaskSelectionReference{}, TaskResultSpecBindings: []HistoricalRecord{}, TaskContentPublications: []TaskContentPublication{}, EpicBaseVersions: []EpicBaseVersion{}, EpicBaseUpdates: []EpicBaseUpdate{}, TaskQueueEvents: []TaskQueueEvent{}, TaskPreparations: []HistoricalRecord{}, EpicCreateOperations: []EpicCreateOperation{}}
}
func (store *systemWorkItemStore) Snapshot(work workspace.ProjectWork) (WorkItemRegistry, error) {
	return store.read(work)
}

func (store *systemWorkItemStore) WithLock(work workspace.ProjectWork, operation func(WorkItemStoreSession) error) error {
	lock, err := store.acquire(work)
	if err != nil {
		return err
	}
	session := &systemWorkItemStoreSession{store: store, work: work}
	_, session.definitionsDigest, err = workspace.ReadDefinitions(work.Workspace)
	var opErr error
	if err != nil {
		opErr = err
	} else {
		opErr = operation(session)
	}
	session.closed = true
	session.confirmed, session.observed = nil, nil
	releaseErr := store.release(lock)
	if opErr != nil {
		if releaseErr != nil {
			return workspace.WorkError(workspace.ErrorWorkIO, fmt.Sprintf("release after failure: %v; %v", opErr, releaseErr), opErr)
		}
		return opErr
	}
	return releaseErr
}

func (session *systemWorkItemStoreSession) Snapshot() (WorkItemRegistry, error) {
	if session.closed {
		return session.store.read(session.work)
	}
	info, err := os.Lstat(session.work.Registry())
	if session.confirmed != nil && err == nil && session.observed != nil && os.SameFile(info, session.observed) && sameRegistryFileState(info, session.observed) {
		if err := session.store.verifyIntegrationEvidence(session.work, *session.confirmed); err != nil {
			session.confirmed, session.observed = nil, nil
			return WorkItemRegistry{}, err
		}
		return cloneWorkItemRegistry(*session.confirmed), nil
	}
	session.confirmed, session.observed = nil, nil
	registry, err := session.store.read(session.work)
	if err != nil {
		return WorkItemRegistry{}, err
	}
	session.confirm(registry, info)
	return registry, nil
}
func (session *systemWorkItemStoreSession) Publish(registry WorkItemRegistry) error {
	if session.closed {
		return workspace.WorkError(workspace.ErrorWorkStoreConflict, "work-item session lock period has ended", nil)
	}
	_, digest, err := workspace.ReadDefinitions(session.work.Workspace)
	if err != nil {
		return err
	}
	if digest != session.definitionsDigest {
		return workspace.WorkError(workspace.ErrorWorkStoreConflict, "definitions changed while holding runtime lock; retry", nil)
	}
	err = session.publish(registry)
	if err != nil {
		session.confirmed, session.observed = nil, nil
	}
	return err
}

func publishWorkItemRegistryRecover(session WorkItemStoreSession, registry WorkItemRegistry) (WorkItemRegistry, error) {
	err := session.Publish(registry)
	if err == nil {
		// Publish accepts a value; return the observed generation, including its
		// physical digest and external references, rather than the caller's old one.
		return session.Snapshot()
	}
	observed, readErr := session.Snapshot()
	if readErr != nil {
		return WorkItemRegistry{}, err
	}
	want, wantErr := encodeWorkItemRegistry(registry)
	got, gotErr := encodeWorkItemRegistry(observed)
	if wantErr == nil && gotErr == nil && bytes.Equal(want, got) {
		return observed, nil
	}
	return WorkItemRegistry{}, err
}

func (store *systemWorkItemStore) acquire(work workspace.ProjectWork) (*filelock.Lock, error) {
	if err := work.Create(); err != nil {
		return nil, err
	}
	path := work.Lock()
	lock, err := filelock.Acquire(path, filelock.Options{Fail: store.fail})
	var lockErr *filelock.Error
	if errors.As(err, &lockErr) {
		return nil, workItemLockError(path, lockErr)
	}
	return lock, err
}

func workItemLockError(path string, err *filelock.Error) error {
	switch err.Op {
	case "inspect":
		return workspace.WorkError(workspace.ErrorWorkStoreConflict, fmt.Sprintf("cannot inspect work-item lock %s", path), err.Err)
	case "not-regular":
		return workspace.WorkError(workspace.ErrorWorkStoreConflict, fmt.Sprintf("work-item lock %s is not a regular file", path), nil)
	case "open":
		if err.Injected {
			return workspace.WorkError(workspace.ErrorWorkIO, "open work-item lock", err.Err)
		}
		return workspace.WorkError(workspace.ErrorWorkIO, fmt.Sprintf("open work-item lock %s", path), err.Err)
	case "stat-opened", "lstat-path", "changed":
		return workspace.WorkError(workspace.ErrorWorkStoreConflict, fmt.Sprintf("work-item lock %s changed or is not a regular file", path), err.Err)
	default:
		return workspace.WorkError(workspace.ErrorWorkIO, "lock work-item registry", err.Err)
	}
}

func (store *systemWorkItemStore) release(lock *filelock.Lock) error {
	if err := lock.Release(); err != nil {
		return workItemReleaseError(err)
	}
	return nil
}

func workItemReleaseError(err *filelock.ReleaseError) error {
	if err.Unlock != nil {
		return workspace.WorkError(workspace.ErrorWorkIO, fmt.Sprintf("unlock work-item lock %s", err.Path), err.Unlock)
	}
	return workspace.WorkError(workspace.ErrorWorkIO, fmt.Sprintf("close work-item lock %s", err.Path), err.Close)
}

func (store *systemWorkItemStore) read(work workspace.ProjectWork) (WorkItemRegistry, error) {
	return store.readRegistry(work)
}

// SnapshotRegistrations reads and validates the registry as Snapshot does,
// including the integration evidence objects it references; a registry
// Snapshot rejects is rejected here too.
func (store *systemWorkItemStore) SnapshotRegistrations(work workspace.ProjectWork) (WorkItemRegistry, error) {
	return store.readRegistry(work)
}

func (store *systemWorkItemStore) readRegistry(work workspace.ProjectWork) (WorkItemRegistry, error) {
	return store.readRegistryAs(work, func(contents []byte) (WorkItemRegistry, error) { return store.decodeStoredRegistry(work, contents) })
}

// RepoUses implements the register's RepoUses with the registrations of work.
func (store *systemWorkItemStore) RepoUses(work workspace.ProjectWork) ([]workspace.RepoUse, error) {
	registry, err := store.SnapshotRegistrations(work)
	if err != nil {
		return nil, err
	}
	lifecycle, err := readWorkItemLifecycle(work)
	if err != nil {
		return nil, err
	}
	return repoUses(registry, WorkItemLifecycles{lifecycle}), nil
}

// repoUses lists the uses of repositories in registry in the order RepoUses
// promises. A use is live while its work item holds a worktree or an open
// operation: an Epic that is not archived, a Task whose worktree is being
// created, ready or awaiting reconciliation, a worktree operation of such a
// Task and an Epic create that is not completed.
func repoUses(registry WorkItemRegistry, lifecycles WorkItemLifecycles) []workspace.RepoUse {
	var uses []workspace.RepoUse
	for _, e := range registry.Epics {
		live := lifecycles.Epic(e.ID) != LifecycleArchived
		for _, b := range e.RepoBindings {
			uses = append(uses, workspace.RepoUse{Kind: workspace.RepoUseEpic, ID: string(e.ID), RepoID: b.RepoID, GitCommonDir: b.GitCommonDir, Worktree: b.Worktree.Locator, Live: live})
		}
	}
	liveTasks := map[TaskID]bool{}
	for _, task := range registry.Tasks {
		live := task.WorktreeState == WorkItemCreating || task.WorktreeState == WorkItemReady || task.WorktreeState == WorkItemReconciliationRequired
		liveTasks[task.ID] = live
		uses = append(uses, workspace.RepoUse{Kind: workspace.RepoUseTask, ID: string(task.ID), RepoID: task.RepoID, GitCommonDir: task.GitCommonDir, Live: live})
	}
	for _, operation := range registry.WorktreeOperations {
		uses = append(uses, workspace.RepoUse{Kind: workspace.RepoUseWorktreeOperation, ID: string(operation.ID), RepoID: operation.RepoID, Live: liveTasks[operation.TaskID]})
	}
	for _, operation := range registry.EpicCreateOperations {
		uses = append(uses, workspace.RepoUse{Kind: workspace.RepoUseEpicCreate, ID: string(operation.ID), RepoID: operation.Request.RepoID, Live: operation.State != "completed"})
	}
	return uses
}

func (store *systemWorkItemStore) readRegistryAs(work workspace.ProjectWork, decode func([]byte) (WorkItemRegistry, error)) (WorkItemRegistry, error) {
	path := work.Registry()
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return emptyWorkItemRegistry(), nil
	}
	if err != nil {
		return WorkItemRegistry{}, workspace.WorkError(workspace.ErrorWorkStoreConflict, fmt.Sprintf("cannot inspect work-item registry %s", path), err)
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
		return WorkItemRegistry{}, workspace.WorkError(workspace.ErrorWorkStoreConflict, fmt.Sprintf("work-item registry %s is not a regular file", path), nil)
	}
	if fault := store.fail("registry-read"); fault != nil {
		return WorkItemRegistry{}, fault
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return WorkItemRegistry{}, workspace.WorkError(workspace.ErrorWorkStoreConflict, fmt.Sprintf("cannot read work-item registry %s", path), err)
	}
	registry, err := decode(contents)
	if err != nil {
		if d := RegistryError(err); d != nil {
			d.Artifact = path
		}
		return WorkItemRegistry{}, workspace.WorkError(workspace.ErrorWorkStoreConflict, fmt.Sprintf("invalid work-item registry %s: %v", path, err), err)
	}
	registry.RawSHA256 = digestTaskBytes(contents)
	return registry, nil
}

func (session *systemWorkItemStoreSession) publish(registry WorkItemRegistry) error {
	store, work := session.store, session.work
	if store.contentScope != nil && store.contentScope.task == "" {
		return queueError("task_goal_execution_conflict", "operation content dependencies are not bound")
	}
	old, err := session.Snapshot()
	if err != nil {
		return err
	}
	if err := validateEpicCreateTransition(old, registry); err != nil {
		return err
	}
	if err := validateQueueTransition(old, registry); err != nil {
		return workspace.WorkError(workspace.ErrorWorkStoreConflict, err.Error(), err)
	}
	if err := validateTaskScopeTransition(old, registry, store.scope); err != nil {
		return err
	}
	contents, err := store.encodeStoredRegistry(work, &registry, old)
	if err != nil {
		return workspace.WorkError(workspace.ErrorWorkStoreConflict, fmt.Sprintf("cannot encode work-item registry: %v", err), err)
	}
	directory := work.Dir
	destination := work.Registry()
	if fault := store.fail("temp-open"); fault != nil {
		return workspace.WorkError(workspace.ErrorWorkIO, "open work-item temporary file", fault)
	}
	temporary, err := os.CreateTemp(directory, ".work-items-*.tmp")
	if err != nil {
		return workspace.WorkError(workspace.ErrorWorkIO, "open work-item temporary file", err)
	}
	temporaryPath := temporary.Name()
	replaced := false
	abort := func(operation string, primary error) error {
		if closeErr := temporary.Close(); closeErr != nil && !errors.Is(closeErr, os.ErrClosed) {
			primary = fmt.Errorf("%w; close: %v", primary, closeErr)
		}
		if !replaced {
			removeErr := store.fail("cleanup")
			if removeErr == nil {
				removeErr = os.Remove(temporaryPath)
			}
			if removeErr != nil && !errors.Is(removeErr, fs.ErrNotExist) {
				primary = fmt.Errorf("%w; cleanup: %v", primary, removeErr)
			}
		}
		return workspace.WorkError(workspace.ErrorWorkIO, operation+" work-item registry", primary)
	}
	if fault := store.fail("write"); fault != nil {
		return abort("write", fault)
	}
	written := 0
	if store.faults != nil && store.faults.shortWrite && len(contents) > 0 {
		written, err = temporary.Write(contents[:len(contents)-1])
	} else {
		written, err = temporary.Write(contents)
	}
	if err != nil {
		return abort("write", err)
	}
	if written != len(contents) {
		return abort("write", io.ErrShortWrite)
	}
	if fault := store.fail("chmod"); fault != nil {
		return abort("chmod", fault)
	}
	if err := temporary.Chmod(0o644); err != nil {
		return abort("chmod", err)
	}
	if fault := store.fail("file-sync"); fault != nil {
		return abort("sync", fault)
	}
	if err := temporary.Sync(); err != nil {
		return abort("sync", err)
	}
	writtenInfo, err := temporary.Stat()
	if err != nil {
		return abort("stat", err)
	}
	if err := temporary.Close(); err != nil {
		return abort("close", err)
	}
	if fault := store.fail("close"); fault != nil {
		return abort("close", fault)
	}
	if fault := store.fail("replace"); fault != nil {
		return abort("replace", fault)
	}
	if err := workspace.PlatformReplace(temporaryPath, destination); err != nil {
		return abort("replace", err)
	}
	replaced = true
	// Bind confirmation to the file we wrote, before any post-rename hook or
	// directory sync. Rename may update native change time, so compare identity,
	// size and modification time to the temporary file, then retain the new stat.
	publishedInfo, statErr := os.Lstat(destination)
	if statErr != nil || !os.SameFile(writtenInfo, publishedInfo) || writtenInfo.Size() != publishedInfo.Size() || !writtenInfo.ModTime().Equal(publishedInfo.ModTime()) || writtenInfo.Mode() != publishedInfo.Mode() {
		publishedInfo = nil
	}
	if fault := store.fail("directory-sync"); fault != nil {
		return workspace.WorkError(workspace.ErrorWorkIO, "sync work-item registry directory", fault)
	}
	if err := workspace.PlatformSyncDirectory(directory); err != nil {
		return workspace.WorkError(workspace.ErrorWorkIO, "sync work-item registry directory", err)
	}
	registry.RawSHA256 = digestTaskBytes(contents)
	session.confirm(registry, publishedInfo)
	return nil
}

func (store *systemWorkItemStore) fail(operation string) error {
	if store.faults == nil || store.faults.fail == nil {
		return nil
	}
	return store.faults.fail(operation)
}
func firstError(values ...error) error {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

func encodeWorkItemRegistry(registry WorkItemRegistry) ([]byte, error) {
	if err := validateWorkItemRegistry(registry); err != nil {
		return nil, err
	}
	var output bytes.Buffer
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)
	if err := encoder.Encode(registry); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	contents := output.Bytes()
	if err := validateWorkItemRegistryYAMLShape(contents); err != nil {
		return nil, err
	}
	return contents, nil
}

func decodeWorkItemRegistry(contents []byte) (registry WorkItemRegistry, err error) {
	root, err := parseRegistryYAML(contents)
	if err != nil {
		return registry, err
	}
	defer func() { err = classifyRegistryDecode(err, RegistryFormat) }()
	if err := rejectRemovedTaskAdoption(root); err != nil {
		return registry, err
	}
	// yaml can truncate floats when decoding into int. Check this contract
	// marker's raw type before typed decoding can erase that distinction.
	if tasks := registryNodeField(root, "tasks"); tasks != nil && tasks.Kind == yaml.SequenceNode {
		for _, task := range tasks.Content {
			if binding := registryNodeField(task, "worktree"); binding != nil {
				if marker := registryNodeField(binding, "native_receipt_version"); marker != nil && (marker.Kind != yaml.ScalarNode || marker.Tag != "!!int") {
					return WorkItemRegistry{}, errors.New("native_receipt_version must be a YAML integer")
				}
			}
		}
	}
	if err := decodeRegistryNode(root, &registry); err != nil {
		return WorkItemRegistry{}, err
	}
	if registry.Epics == nil || registry.Tasks == nil || registry.WorktreeOperations == nil {
		return WorkItemRegistry{}, errors.New("epics, tasks, and worktree_operations must be explicit arrays")
	}
	if registry.TaskResults == nil || registry.HumanQARecords == nil || registry.IntegrationAuthorities == nil || registry.IntegrationIntents == nil || registry.IntegrationAttempts == nil || registry.IntegrationResults == nil {
		return WorkItemRegistry{}, errors.New("Task lifecycle collections must be explicit arrays")
	}
	if err := validateWorkItemRegistryNodeShape(root); err != nil {
		return WorkItemRegistry{}, err
	}
	if err := validateWorkItemRegistry(registry); err != nil {
		return WorkItemRegistry{}, err
	}
	return registry, nil
}

// Reject the removed Task model before typed decoding can discard its fields.
// Epic origin=adopted and the historical format boundary remain supported.
func rejectRemovedTaskAdoption(root *yaml.Node) error {
	unsupported := func() error {
		return errors.New("Task adoption is no longer supported; preserve the registry and use a compatible historical reader to inspect it; do not strip adoption history")
	}
	if ops := registryNodeField(root, "worktree_operations"); ops != nil {
		for _, op := range ops.Content {
			for _, key := range []string{"adoption", "release", "adoption_observations"} {
				if registryNodeField(op, key) != nil {
					return unsupported()
				}
			}
		}
	}
	if tasks := registryNodeField(root, "tasks"); tasks != nil {
		for _, task := range tasks.Content {
			if b := registryNodeField(task, "worktree"); b != nil {
				if origin := registryNodeField(b, "origin"); origin != nil && origin.Value == "adopted" {
					return unsupported()
				}
			}
		}
	}
	if preparations := registryNodeField(root, "task_preparations"); preparations != nil {
		for _, preparation := range preparations.Content {
			if plan := registryNodeField(preparation, "plan"); plan != nil {
				if registryNodeField(plan, "adoption_id") != nil {
					return unsupported()
				}
				if mode := registryNodeField(plan, "source_mode"); mode != nil && mode.Value == "adopted" {
					return unsupported()
				}
			}
		}
	}
	return nil
}

type yamlShapeField struct {
	name   string
	typeOf reflect.Type
}

func validateWorkItemRegistryYAMLShape(contents []byte) error {
	var document yaml.Node
	if err := yaml.Unmarshal(contents, &document); err != nil || len(document.Content) != 1 {
		if err == nil {
			err = errors.New("registry must contain one YAML document")
		}
		return err
	}
	return validateWorkItemRegistryNodeShape(document.Content[0])
}

// The logical registry has one shape: every collection is an explicit array in
// declaration order; only the Task Spec policy list may be absent.
func validateWorkItemRegistryNodeShape(root *yaml.Node) error {
	present := false
	for i := 0; i < len(root.Content); i += 2 {
		if root.Content[i].Value == "task_spec_policies" {
			present = true
		}
	}
	fields := []yamlShapeField{}
	for _, f := range yamlStructShapeFields(reflect.TypeOf(WorkItemRegistry{})) {
		if f.name == "task_spec_policies" {
			if !present {
				continue
			}
			f.typeOf = reflect.TypeOf([]TaskSpecPolicy{})
		}
		fields = append(fields, f)
	}
	return validateYAMLMappingShape(root, fields, "work-item registry")
}

func validateYAMLNodeShape(node *yaml.Node, expected reflect.Type, context string) error {
	if expected.Kind() == reflect.Pointer {
		if node.Tag == "!!null" {
			if node.Value != "null" {
				return fmt.Errorf("%s must use explicit null", context)
			}
			return nil
		}
		return validateYAMLNodeShape(node, expected.Elem(), context)
	}
	// Records of the removed delivery process are kept, never interpreted.
	if expected == reflect.TypeOf(HistoricalRecord{}) {
		return nil
	}
	if expected == reflect.TypeOf(IntegrationAuthority{}) && historicalAuthorityNode(node) || expected == reflect.TypeOf(IntegrationIntent{}) && historicalIntegrationNode(node) {
		return nil
	}
	switch expected.Kind() {
	case reflect.Map:
		return validateQueueYAMLMap(node)
	case reflect.Struct:
		fields := yamlStructShapeFields(expected)
		if expected == reflect.TypeOf(WorkspaceEpicBaseUpdatePlan{}) {
			present := map[string]bool{}
			for i := 0; i+1 < len(node.Content); i += 2 {
				present[node.Content[i].Value] = true
			}
			filtered := []yamlShapeField{}
			for _, f := range fields {
				if f.name != "current_preparation_id" || present[f.name] {
					filtered = append(filtered, f)
				}
			}
			fields = filtered
		}
		if expected == reflect.TypeOf(TaskWorktreeBinding{}) || expected == reflect.TypeOf(WorktreeOperationRecord{}) {
			present := map[string]bool{}
			for i := 0; i+1 < len(node.Content); i += 2 {
				present[node.Content[i].Value] = true
			}
			filtered := []yamlShapeField{}
			for _, f := range fields {
				optional := f.name == "native_receipt_version" || f.name == "verification_receipt" || f.name == "admission_mode" || f.name == "source_reconciliation" || f.name == "execution" || f.name == "source_mode" || f.name == "start_oid" || f.name == "start_tree"
				if !optional || present[f.name] {
					filtered = append(filtered, f)
				}
			}
			fields = filtered
		}
		if expected == reflect.TypeOf(TaskExecutorAssignment{}) {
			present := false
			for i := 0; i+1 < len(node.Content); i += 2 {
				if node.Content[i].Value == "allow_local_install" {
					present = true
				}
			}
			if !present {
				filtered := []yamlShapeField{}
				for _, f := range fields {
					if f.name != "allow_local_install" {
						filtered = append(filtered, f)
					}
				}
				fields = filtered
			}
		}
		return validateYAMLMappingShape(node, fields, context)
	case reflect.Slice, reflect.Array:
		if node.Kind != yaml.SequenceNode {
			return fmt.Errorf("%s must be an explicit array", context)
		}
		for i, item := range node.Content {
			if err := validateYAMLNodeShape(item, expected.Elem(), fmt.Sprintf("%s[%d]", context, i)); err != nil {
				return err
			}
		}
		return nil
	default:
		if node.Kind != yaml.ScalarNode || node.Tag == "!!null" {
			return fmt.Errorf("%s must be a non-null scalar", context)
		}
		return nil
	}
}

func validateYAMLMappingShape(node *yaml.Node, fields []yamlShapeField, context string) error {
	if node.Kind != yaml.MappingNode || len(node.Content) != len(fields)*2 {
		return fmt.Errorf("%s has missing or extra fields", context)
	}
	for i, field := range fields {
		key := node.Content[i*2]
		if key.Kind != yaml.ScalarNode || key.Value != field.name {
			return fmt.Errorf("%s field %d must be %s", context, i, field.name)
		}
		if err := validateYAMLNodeShape(node.Content[i*2+1], field.typeOf, context+"."+field.name); err != nil {
			return err
		}
	}
	return nil
}

func yamlStructShapeFields(value reflect.Type) []yamlShapeField {
	fields := make([]yamlShapeField, 0, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.PkgPath != "" {
			continue
		}
		name := strings.Split(field.Tag.Get("yaml"), ",")[0]
		if name == "-" {
			continue
		}
		if name == "" {
			name = strings.ToLower(field.Name)
		}
		fields = append(fields, yamlShapeField{name: name, typeOf: field.Type})
	}
	return fields
}

var worktreeIDPattern = regexp.MustCompile(`^wt_[0-9a-f]{32}$`)
var operationIDPattern = regexp.MustCompile(`^wop_[0-9a-f]{32}$`)
var digestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

// validateWorkItemRegistry runs the registry's phases in a fixed order; the
// first failing rule answers. Epics, operations and Tasks are each validated
// in ID order, one record at a time.
func validateWorkItemRegistry(registry WorkItemRegistry) error {
	if err := validateEpicCreates(registry); err != nil {
		return err
	}
	if err := validateQueueRegistry(registry); err != nil {
		return err
	}
	if err := validateTaskContentRegistry(registry); err != nil {
		return err
	}
	if !sortedUniqueEpics(registry.Epics) || !sortedUniqueTasks(registry.Tasks) || !sortedUniqueOperations(registry.WorktreeOperations) {
		return errors.New("work-item arrays must be sorted by unique IDs")
	}
	validation := &workItemRegistryValidation{registry: registry, epics: map[EpicID]*EpicRecord{}, operations: map[TaskID]WorktreeOperationRecord{}, operationIDs: map[WorktreeOperationID]bool{}, worktreeOwners: map[WorktreeID]string{}, resourceOwners: map[string]string{}}
	if err := validation.validateEpics(); err != nil {
		return err
	}
	if err := validation.validateOperations(); err != nil {
		return err
	}
	if err := validation.validateTasks(); err != nil {
		return err
	}
	if err := validation.validateOperationsHaveTasks(); err != nil {
		return err
	}
	return validateTaskLifecycleRegistry(registry)
}

// workItemRegistryValidation carries what the records validated so far
// established for the ones after them: the Epics by ID, the operation of each
// Task, the operation IDs in use, and which record first claimed each
// worktree ID and Git resource.
type workItemRegistryValidation struct {
	registry       WorkItemRegistry
	epics          map[EpicID]*EpicRecord
	operations     map[TaskID]WorktreeOperationRecord
	operationIDs   map[WorktreeOperationID]bool
	worktreeOwners map[WorktreeID]string
	resourceOwners map[string]string
}

// claimFirst records owner for key unless an earlier record holds it, and then
// reports that earlier owner instead.
func claimFirst[K comparable](owners map[K]string, key K, owner string) (earlier string, taken bool) {
	if earlier, taken = owners[key]; taken {
		return earlier, true
	}
	owners[key] = owner
	return "", false
}

// gitResourceKeys names the path and the ref that a worktree occupies in one
// repository, in the order they are claimed.
func gitResourceKeys(gitCommonDir, locator, ref string) []string {
	return []string{gitCommonDir + "\x00path\x00" + locator, gitCommonDir + "\x00ref\x00" + ref}
}

func (validation *workItemRegistryValidation) validateEpics() error {
	for i := range validation.registry.Epics {
		epic := &validation.registry.Epics[i]
		if err := validation.validateEpic(*epic); err != nil {
			return err
		}
		validation.epics[epic.ID] = epic
	}
	return nil
}

func (validation *workItemRegistryValidation) validateEpic(epic EpicRecord) error {
	if err := validateWorkIdentifier("Epic", string(epic.ID)); err != nil {
		return err
	}
	if err := validateWorkText("Epic title", epic.Title, 128); err != nil {
		return err
	}
	if err := validateWorkIdentifier("Project", string(epic.ProjectID)); err != nil {
		return err
	}
	if len(epic.RepoBindings) == 0 {
		return fmt.Errorf("Epic %s has no repository binding", epic.ID)
	}
	if !sort.SliceIsSorted(epic.RepoBindings, func(i, j int) bool { return epic.RepoBindings[i].RepoID < epic.RepoBindings[j].RepoID }) {
		return fmt.Errorf("Epic %s repository bindings are not sorted", epic.ID)
	}
	seenRepos := map[workspace.RepoID]bool{}
	for _, binding := range epic.RepoBindings {
		if err := validation.validateEpicRepoBinding(epic, binding, seenRepos); err != nil {
			return err
		}
	}
	return nil
}

func (validation *workItemRegistryValidation) validateEpicRepoBinding(epic EpicRecord, binding EpicRepoBinding, seenRepos map[workspace.RepoID]bool) error {
	if seenRepos[binding.RepoID] {
		return fmt.Errorf("Epic %s has duplicate repository %s", epic.ID, binding.RepoID)
	}
	seenRepos[binding.RepoID] = true
	if err := validateWorkIdentifier("repository", string(binding.RepoID)); err != nil {
		return err
	}
	if !workspace.ValidStoredPath(binding.GitCommonDir) {
		return fmt.Errorf("Epic %s has invalid common directory", epic.ID)
	}
	wt := binding.Worktree
	if !validEpicWorktree(epic, binding) {
		return fmt.Errorf("Epic %s has invalid worktree binding", epic.ID)
	}
	owner := "Epic " + string(epic.ID)
	if earlier, taken := claimFirst(validation.worktreeOwners, wt.ID, owner); taken {
		return fmt.Errorf("worktree ID %s is shared by %s and Epic %s", wt.ID, earlier, epic.ID)
	}
	// Epics may share an integration base, so a path or ref stays with the
	// first Epic that claimed it; operations and native receipts below may
	// not reuse any Epic path or ref.
	for _, key := range gitResourceKeys(binding.GitCommonDir, wt.Locator, wt.Ref) {
		claimFirst(validation.resourceOwners, key, owner)
	}
	return nil
}

func validEpicWorktree(epic EpicRecord, binding EpicRepoBinding) bool {
	wt := binding.Worktree
	return worktreeIDPattern.MatchString(string(wt.ID)) &&
		wt.OwnerKind == "epic" &&
		wt.OwnerID == epic.ID &&
		(wt.Origin == "adopted" || wt.Origin == "ply_created") &&
		wt.GitCommonDir == binding.GitCommonDir &&
		workspace.ValidStoredPath(wt.Locator) &&
		validFullBranchRef(wt.Ref) &&
		validOIDText(wt.OID) &&
		validOIDText(wt.Tree)
}

func (validation *workItemRegistryValidation) validateOperations() error {
	for _, operation := range validation.registry.WorktreeOperations {
		if err := validation.validateOperation(operation); err != nil {
			return err
		}
		validation.operations[operation.TaskID] = operation
	}
	return nil
}

func (validation *workItemRegistryValidation) validateOperation(operation WorktreeOperationRecord) error {
	if !operationIDPattern.MatchString(string(operation.ID)) || validation.operationIDs[operation.ID] {
		return fmt.Errorf("invalid or duplicate operation ID %s", operation.ID)
	}
	validation.operationIDs[operation.ID] = true
	if _, ok := validation.operations[operation.TaskID]; ok {
		return fmt.Errorf("Task %s has multiple operations", operation.TaskID)
	}
	if !wellFormedOperationIntent(operation) {
		return fmt.Errorf("operation %s has invalid intent", operation.ID)
	}
	owner := "operation " + string(operation.ID)
	if earlier, taken := claimFirst(validation.worktreeOwners, operation.WorktreeID, owner); taken {
		return fmt.Errorf("worktree ID %s is shared by %s and operation %s", operation.WorktreeID, earlier, operation.ID)
	}
	computed, err := intentDigest(operation)
	if err != nil || computed != operation.IntentDigest {
		return fmt.Errorf("operation %s intent digest does not match", operation.ID)
	}
	if operation.State != "creating" && operation.State != "ready" && operation.State != "reconciliation_required" {
		return fmt.Errorf("operation %s has invalid state", operation.ID)
	}
	if err := validateObservation(operation.LastObservation, operation.State); err != nil {
		return fmt.Errorf("operation %s: %w", operation.ID, err)
	}
	for _, key := range gitResourceKeys(operation.GitCommonDir, operation.TargetLocator, operation.SourceRef) {
		if earlier, taken := claimFirst(validation.resourceOwners, key, owner); taken {
			return fmt.Errorf("Git resource shared by %s and operation %s", earlier, operation.ID)
		}
	}
	return nil
}

func wellFormedOperationIntent(operation WorktreeOperationRecord) bool {
	return worktreeIDPattern.MatchString(string(operation.WorktreeID)) &&
		digestPattern.MatchString(operation.IntentDigest) &&
		workspace.ValidStoredPath(operation.GitCommonDir) &&
		workspace.ValidStoredPath(operation.TargetLocator) &&
		validFullBranchRef(operation.ParentRef) &&
		validFullBranchRef(operation.SourceRef) &&
		validOIDText(operation.ParentOID) &&
		validOIDText(operation.ParentTree)
}

func (validation *workItemRegistryValidation) validateTasks() error {
	for _, task := range validation.registry.Tasks {
		if err := validation.validateTask(task); err != nil {
			return err
		}
	}
	return nil
}

func (validation *workItemRegistryValidation) validateTask(task TaskRecord) error {
	if err := validateWorkIdentifier("Task", string(task.ID)); err != nil {
		return err
	}
	if err := validateWorkText("Task title", task.Title, 128); err != nil {
		return err
	}
	if err := validateWorkText("Task description", task.Description, 2048); err != nil {
		return err
	}
	epic := validation.epics[task.ParentEpicID]
	if err := workitem.ValidateTaskParent(task, epic); err != nil {
		return err
	}
	// ValidateTaskParent found the Epic and its binding for the Task's repository.
	binding, _ := findEpicRepo(*epic, task.RepoID)
	parent := binding.Worktree
	operation, hasOperation := validation.operations[task.ID]
	if hasOperation && !validation.operationMatchesParent(task, operation, parent) {
		return fmt.Errorf("Task %s operation differs from its parent repository binding", task.ID)
	}
	// A native receipt is the whole record of a completed worktree; without
	// one, the worktree state says which operation and binding must exist.
	if task.Worktree != nil && task.Worktree.NativeReceiptVersion != 0 {
		return validation.validateNativeReceiptClaims(task, hasOperation)
	}
	return validation.validateTaskWorktreeState(task, operation, hasOperation, parent)
}

func (validation *workItemRegistryValidation) validateTaskWorktreeState(task TaskRecord, operation WorktreeOperationRecord, hasOperation bool, parent EpicWorktreeBinding) error {
	switch task.WorktreeState {
	case WorkItemUnbound:
		if task.Worktree != nil || hasOperation {
			return fmt.Errorf("unbound Task %s has worktree state", task.ID)
		}
		return nil
	case WorkItemCreating:
		return validatePendingTask(task, operation, hasOperation, "creating")
	case WorkItemReconciliationRequired:
		return validatePendingTask(task, operation, hasOperation, "reconciliation_required")
	case WorkItemReady, WorkItemRetired:
		return validation.validateReadyTask(task, operation, hasOperation, parent)
	default:
		return fmt.Errorf("Task %s has invalid worktree_state", task.ID)
	}
}

func (validation *workItemRegistryValidation) operationMatchesParent(task TaskRecord, operation WorktreeOperationRecord, parent EpicWorktreeBinding) bool {
	return operation.ProjectID == task.ProjectID &&
		operation.RepoID == task.RepoID &&
		operation.GitCommonDir == task.GitCommonDir &&
		operation.ParentEpicID == task.ParentEpicID &&
		operation.ParentWorktreeID == parent.ID &&
		operation.ParentRef == parent.Ref &&
		historicalEpicBaseMatches(validation.registry, task.ParentEpicID, task.RepoID, operation.ParentWorktreeID, operation.ParentRef, operation.ParentOID, operation.ParentTree)
}

func (validation *workItemRegistryValidation) validateNativeReceiptClaims(task TaskRecord, hasOperation bool) error {
	if err := validateNativeReceipt(validation.registry, task, hasOperation); err != nil {
		return err
	}
	wt := task.Worktree
	owner := "Task " + string(task.ID)
	if _, taken := claimFirst(validation.worktreeOwners, wt.ID, owner); taken || validation.operationIDs[wt.OperationID] {
		return fmt.Errorf("native receipt reuses worktree or operation ID")
	}
	validation.operationIDs[wt.OperationID] = true
	for _, key := range gitResourceKeys(wt.GitCommonDir, wt.Locator, wt.Ref) {
		if earlier, taken := claimFirst(validation.resourceOwners, key, owner); taken {
			return fmt.Errorf("native receipt shares Git resource with %s", earlier)
		}
	}
	return nil
}

// A Task whose worktree is being created or reconciled has no binding yet,
// only an operation in that same state.
func validatePendingTask(task TaskRecord, operation WorktreeOperationRecord, hasOperation bool, operationState string) error {
	if task.Worktree != nil || !hasOperation || operation.State != operationState {
		return fmt.Errorf("Task %s has incomplete operation state", task.ID)
	}
	return nil
}

func (validation *workItemRegistryValidation) validateReadyTask(task TaskRecord, operation WorktreeOperationRecord, hasOperation bool, parent EpicWorktreeBinding) error {
	if task.Worktree == nil || !hasOperation || operation.State != "ready" {
		return fmt.Errorf("ready Task %s has invalid binding", task.ID)
	}
	wt := task.Worktree
	if !readyBindingMatchesOperation(task, operation) {
		return fmt.Errorf("ready Task %s binding differs from operation", task.ID)
	}
	if owner, ok := validation.worktreeOwners[wt.ID]; ok && owner != "operation "+string(operation.ID) {
		return fmt.Errorf("worktree ID %s reused", wt.ID)
	}
	if !exactObservationMatchesOperation(operation, parent.Locator) {
		return fmt.Errorf("ready Task %s exact observation differs from operation", task.ID)
	}
	return nil
}

func readyBindingMatchesOperation(task TaskRecord, operation WorktreeOperationRecord) bool {
	wt := task.Worktree
	return wt.ID == operation.WorktreeID &&
		wt.OperationID == operation.ID &&
		wt.IntentDigest == operation.IntentDigest &&
		wt.OwnerKind == "task" &&
		wt.OwnerID == task.ID &&
		wt.Origin == operationOrigin(operation) &&
		wt.Locator == operation.TargetLocator &&
		wt.Ref == operation.SourceRef &&
		wt.OID == operationStartOID(operation) &&
		wt.Tree == operationStartTree(operation) &&
		wt.GitCommonDir == operation.GitCommonDir &&
		wt.ParentEpicID == task.ParentEpicID &&
		wt.ParentWorktreeID == operation.ParentWorktreeID &&
		wt.ParentRef == operation.ParentRef &&
		wt.ParentOID == operation.ParentOID &&
		wt.ParentTree == operation.ParentTree
}

// exactObservationMatchesOperation compares the last observation of a ready
// operation with its intent: parent and new worktree are clean and exactly
// where the operation put them. validateObservation has already required
// every observed field of a ready operation.
func exactObservationMatchesOperation(operation WorktreeOperationRecord, parentLocator string) bool {
	observed := operation.LastObservation
	return *observed.ParentRefOID == operation.ParentOID &&
		*observed.ParentRefTree == operation.ParentTree &&
		*observed.ParentWorktreeLocator == parentLocator &&
		*observed.ParentWorktreeRef == operation.ParentRef &&
		*observed.ParentWorktreeOID == operation.ParentOID &&
		*observed.ParentWorktreeTree == operation.ParentTree &&
		*observed.ParentWorktreeGitCommonDir == operation.GitCommonDir &&
		*observed.ParentWorktreeClean &&
		*observed.SourceRefOID == operationStartOID(operation) &&
		*observed.SourceRefTree == operationStartTree(operation) &&
		*observed.SourceRefCheckedOutAt == operation.TargetLocator &&
		*observed.TargetLocator == operation.TargetLocator &&
		*observed.TargetRef == operation.SourceRef &&
		*observed.TargetOID == operationStartOID(operation) &&
		*observed.TargetTree == operationStartTree(operation) &&
		*observed.TargetGitCommonDir == operation.GitCommonDir &&
		*observed.TargetClean &&
		*observed.InventoryMatch
}

func (validation *workItemRegistryValidation) validateOperationsHaveTasks() error {
	for _, taskID := range slices.Sorted(maps.Keys(validation.operations)) {
		if task, _ := findTask(validation.registry, taskID); task == nil {
			return fmt.Errorf("operation refers to missing Task %s", taskID)
		}
	}
	return nil
}

func validateTaskLifecycleRegistry(registry WorkItemRegistry) error {
	if !sortedAuthorities(registry.IntegrationAuthorities) || !sortedIntents(registry.IntegrationIntents) || !sortedAttempts(registry.IntegrationAttempts) || !sortedIntegrationResults(registry.IntegrationResults) {
		return errors.New("Task lifecycle arrays must be sorted by unique IDs")
	}
	authorities := map[IntegrationAuthorityID]IntegrationAuthority{}
	planDigests := map[string]bool{}
	for _, a := range registry.IntegrationAuthorities {
		if !lifecycleIDPatterns["integration authority"].MatchString(string(a.ID)) || planDigests[a.PlanSHA256] {
			return fmt.Errorf("integration authority %s is invalid", a.ID)
		}
		if a.historical == nil {
			if !validTaskUTC(a.CreatedAtUTC) || !digestPattern.MatchString(a.PlanSHA256) || integrationPlanDigest(a.Plan) != a.PlanSHA256 {
				return fmt.Errorf("integration authority %s is invalid", a.ID)
			}
			if err := validateStoredIntegrationPlan(a.Plan, registry); err != nil || a.Plan.Task.TaskID != a.TaskID || !sameOptionalIntegrationResultID(a.RetryAfterResultID, a.Plan.RetryAfterResultID) || !sameAllowedEffect(a.AllowedEffect, a.Plan.Effect) {
				return fmt.Errorf("integration authority %s has invalid plan binding", a.ID)
			}
		}
		planDigests[a.PlanSHA256] = true
		authorities[a.ID] = a
	}
	intents := map[IntegrationIntentID]IntegrationIntent{}
	intentByAuthority := map[IntegrationAuthorityID]bool{}
	for _, i := range registry.IntegrationIntents {
		a, ok := authorities[i.AuthorityID]
		if !ok || intentByAuthority[i.AuthorityID] || !lifecycleIDPatterns["integration intent"].MatchString(string(i.ID)) || i.State != "prepared" || i.MaxAttempts != 1 || i.PlanSHA256 != a.PlanSHA256 || a.historical == nil && (i.historical != nil || integrationIntentDigest(i) != i.IntentSHA256 || !intentMatchesAuthority(i, a)) {
			return fmt.Errorf("integration intent %s is invalid", i.ID)
		}
		intentByAuthority[i.AuthorityID] = true
		intents[i.ID] = i
	}
	attemptByAuthority := map[IntegrationAuthorityID]bool{}
	attempts := map[IntegrationAttemptID]IntegrationAttempt{}
	for _, a := range registry.IntegrationAttempts {
		authority, ok := authorities[a.AuthorityID]
		if !ok || attemptByAuthority[a.AuthorityID] || !lifecycleIDPatterns["integration attempt"].MatchString(string(a.ID)) || a.Ordinal != 1 || (a.State != "prepared" && a.State != "result_recorded") {
			return fmt.Errorf("integration attempt %s is invalid", a.ID)
		}
		intent, ok := intents[a.IntentID]
		expectedPre := integrationObservation(a.PreparedAtUTC, worktreeFromPlan(authority.Plan.ObservedParent), authority.Plan.ObservedInventory)
		if !ok || intent.AuthorityID != a.AuthorityID || a.CommandIdentity != "git-merge-ff-only" || a.ReflogAction != "ply-workspace-task-integrate:"+string(a.ID) || !validTaskUTC(a.PreparedAtUTC) || validateIntegrationObservation(a.PreObservation) != nil || !reflect.DeepEqual(a.PreObservation, expectedPre) || (a.State == "prepared") != (a.ResultID == nil) {
			return fmt.Errorf("integration attempt %s has missing intent", a.ID)
		}
		attemptByAuthority[a.AuthorityID] = true
		attempts[a.ID] = a
	}
	resultByAuthority := map[IntegrationAuthorityID]bool{}
	for _, r := range registry.IntegrationResults {
		if _, ok := authorities[r.AuthorityID]; !ok || resultByAuthority[r.AuthorityID] || !lifecycleIDPatterns["integration result"].MatchString(string(r.ID)) {
			return fmt.Errorf("integration result %s is invalid", r.ID)
		}
		intent, ok := intents[r.IntentID]
		authority := authorities[r.AuthorityID]
		if !ok || intent.AuthorityID != r.AuthorityID || r.RetryAfterResultID == nil != (authority.RetryAfterResultID == nil) || (r.RetryAfterResultID != nil && *r.RetryAfterResultID != *authority.RetryAfterResultID) || validateIntegrationResultShape(r) != nil || validateIntegrationResultOutcomeContract(r, authority) != nil {
			return fmt.Errorf("integration result %s has missing intent", r.ID)
		}
		if r.AttemptID != nil {
			a, ok := attempts[*r.AttemptID]
			if !ok || a.AuthorityID != r.AuthorityID || a.IntentID != r.IntentID || a.ResultID == nil || *a.ResultID != r.ID || a.State != "result_recorded" || !reflect.DeepEqual(r.BeforeObservation, a.PreObservation) {
				return fmt.Errorf("integration result %s has invalid attempt", r.ID)
			}
		}
		if r.Outcome == "exact_effect" {
			authority := authorities[r.AuthorityID]
			expectedParent := worktreeFromPlan(authority.Plan.ObservedParent)
			expectedParent.OID = authority.Plan.Task.ResultOID
			expectedParent.Tree = authority.Plan.Task.ResultTree
			expectedInventory := inventoryAfterFastForward(authority.Plan.ObservedInventory, authority.Plan.Epic.ParentLocator, authority.Plan.Epic.ParentRef, authority.Plan.Task.ResultOID)
			expectedAfter := integrationObservation(r.RecordedAtUTC, expectedParent, expectedInventory)
			if r.AttemptID == nil || !reflogProves(r.Reflog, *r.AttemptID, authority.Plan.Epic.ExpectedParentOID, authority.Plan.Task.ResultOID) || !reflect.DeepEqual(r.AfterObservation, expectedAfter) {
				return fmt.Errorf("integration result %s does not prove its exact effect", r.ID)
			}
		}
		if r.Outcome == "no_effect" {
			attempt := attempts[*r.AttemptID]
			if !sameIntegrationObservation(r.BeforeObservation, r.AfterObservation) || reflogMentionsAttempt(r.Reflog, attempt.ID) {
				return errors.New("no-effect integration result contains effect evidence")
			}
		}
		resultByAuthority[r.AuthorityID] = true
	}
	retryUse := map[IntegrationResultID]bool{}
	for _, authority := range registry.IntegrationAuthorities {
		if authority.RetryAfterResultID == nil {
			continue
		}
		prior, ok := integrationResultByID(registry.IntegrationResults, *authority.RetryAfterResultID)
		retired := authorityRetired(registry, &authority)
		if !ok || prior.Outcome != "no_effect" || !retired && retryUse[*authority.RetryAfterResultID] {
			return fmt.Errorf("integration authority %s has invalid retry predecessor", authority.ID)
		}
		priorAuthority := authorities[prior.AuthorityID]
		if priorAuthority.TaskID != authority.TaskID || priorAuthority.Plan.Epic.ParentRef != authority.Plan.Epic.ParentRef {
			return fmt.Errorf("integration authority %s retry binding differs", authority.ID)
		}
		retryUse[*authority.RetryAfterResultID] = retryUse[*authority.RetryAfterResultID] || !retired
	}
	return nil
}

func validateStoredIntegrationPlan(p WorkspaceTaskIntegrationPlan, registry WorkItemRegistry) error {
	task, _ := findTask(registry, p.Task.TaskID)
	if task == nil || task.Worktree == nil {
		return errors.New("integration plan Task binding is missing")
	}
	if !(p.Kind == currentIntegrationPlanKind && p.SchemaVersion == 6 || refreshIntegrationPlan(p) && p.SchemaVersion == 7) || p.Format != "json" || p.FormatVersion != 1 || p.Canonicalization != "RFC8785" || !validAbsoluteCleanPath(p.Workspace.Root) || !digestPattern.MatchString(p.Workspace.MarkerSHA256) || p.Project.ProjectID != task.ProjectID || p.Repository.RepoID != task.RepoID || !validAbsoluteCleanPath(p.Repository.RegisteredLocator) || p.Repository.GitCommonDir != task.GitCommonDir || !validAbsoluteCleanPath(p.Repository.GitCommonDir) || !setString("sha1", "sha256")[p.Repository.ObjectFormat] || p.Repository.RefFormat != "files" {
		return errors.New("invalid integration plan envelope")
	}
	epic, _ := findEpic(registry, task.ParentEpicID)
	if epic == nil {
		return errors.New("integration plan Epic binding is missing")
	}
	parent, _ := findEpicRepo(*epic, task.RepoID)
	if parent == nil || p.Task.TaskWorktreeID != task.Worktree.ID || p.Task.SourceLocator != task.Worktree.Locator || p.Task.SourceRef != task.Worktree.Ref || p.Epic.EpicID != epic.ID || p.Epic.ParentWorktreeID != parent.Worktree.ID || p.Epic.ParentLocator != parent.Worktree.Locator || p.Epic.ParentRef != parent.Worktree.Ref || !storedIntegrationPlanBasis(registry, p, *task, *parent) || p.StoreTransition != "none" || !digestPattern.MatchString(p.WorkItemsSHA256) {
		return errors.New("integration plan binding differs")
	}
	if validatePlanObservation(p.ObservedSource) != nil || validatePlanObservation(p.ObservedParent) != nil || validateInventory(p.ObservedInventory) != nil || validateReflog(p.ObservedReflog) != nil || !setString("ready", "already_integrated", "blocked", "conflict", "unknown")[p.Readiness] || !sortedReasonTokens(p.Reasons) {
		return errors.New("invalid integration plan observation")
	}
	source := worktreeFromPlan(p.ObservedSource)
	parentObservation := worktreeFromPlan(p.ObservedParent)
	if !setString("ready", "already_integrated")[p.Readiness] || len(p.Reasons) != 0 || !p.AncestryReady || p.Repository.Shallow || p.Repository.PartialClone || p.Repository.SparseCheckout || p.ObservedSource.OID != p.Task.ResultOID || p.ObservedSource.Tree != p.Task.ResultTree || !p.ObservedSource.Clean || len(p.ObservedSource.StatusEntries) != 0 || len(p.ObservedSource.InProgress) != 0 || !p.ObservedParent.Clean || len(p.ObservedParent.StatusEntries) != 0 || len(p.ObservedParent.InProgress) != 0 || (p.Readiness == "ready" && (p.ObservedParent.OID != p.Epic.ExpectedParentOID || p.ObservedParent.Tree != p.Epic.ExpectedParentTree)) || (p.Readiness == "already_integrated" && (p.ObservedParent.OID != p.Task.ResultOID || p.ObservedParent.Tree != p.Task.ResultTree)) || p.ObservedSource.GitCommonDir != p.Repository.GitCommonDir || p.ObservedParent.GitCommonDir != p.Repository.GitCommonDir || p.ObservedSource.ObjectFormat != p.Repository.ObjectFormat || p.ObservedParent.ObjectFormat != p.Repository.ObjectFormat || p.ObservedSource.RefFormat != p.Repository.RefFormat || p.ObservedParent.RefFormat != p.Repository.RefFormat || !integrationInventoryContains(p.ObservedInventory, source) || !integrationInventoryContains(p.ObservedInventory, parentObservation) || !objectIDsMatchFormat(p.Repository.ObjectFormat, p.Epic.ExpectedParentOID, p.Epic.ExpectedParentTree, p.Task.ResultOID, p.Task.ResultTree, p.ObservedSource.OID, p.ObservedSource.Tree, p.ObservedParent.OID, p.ObservedParent.Tree) {
		return errors.New("integration authority plan is not a ready exact plan")
	}
	wantArgv := []string{"git", "-c", "core.hooksPath=" + os.DevNull, "-c", "merge.autoStash=false", "-c", "gc.auto=0", "-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false", "-c", "submodule.recurse=false", "-C", p.Epic.ParentLocator, "merge", "--ff-only", "--no-stat", "--no-autostash", p.Task.ResultOID}
	if p.Effect.Kind != "local_ff_only" || p.Effect.ParentRef != p.Epic.ParentRef || p.Effect.ExpectedParentOID != p.Epic.ExpectedParentOID || p.Effect.ResultOID != p.Task.ResultOID || p.Effect.MaxOccurrences != 1 || !slices.Equal(p.Effect.Argv, wantArgv) {
		return errors.New("invalid integration plan effect")
	}
	return nil
}

func storedIntegrationPlanBasis(registry WorkItemRegistry, p WorkspaceTaskIntegrationPlan, task TaskRecord, parent EpicRepoBinding) bool {
	if !refreshIntegrationPlan(p) {
		return p.Epic.ExpectedParentOID == task.Worktree.ParentOID && p.Epic.ExpectedParentTree == task.Worktree.ParentTree
	}
	return p.Epic.ExpectedParentOID != task.Worktree.ParentOID && historicalEpicBaseMatches(registry, task.ParentEpicID, task.RepoID, parent.Worktree.ID, parent.Worktree.Ref, p.Epic.ExpectedParentOID, p.Epic.ExpectedParentTree)
}

func validatePlanObservation(o PlanWorktreeObservation) error {
	if !validAbsoluteCleanPath(o.WorktreeLocator) || !validFullBranchRef(o.Ref) || !validOIDText(o.OID) || !validOIDText(o.Tree) || len(o.OID) != len(o.Tree) || !validAbsoluteCleanPath(o.GitCommonDir) || !setString("sha1", "sha256")[o.ObjectFormat] || o.RefFormat != "files" || !o.Symbolic || validateStatusEntries(o.StatusEntries) != nil || o.InProgress == nil || !sort.StringsAreSorted(o.InProgress) {
		return errors.New("invalid worktree observation")
	}
	return nil
}
func validateIntegrationObservation(o IntegrationObservation) error {
	if !validTaskUTC(o.ObservedAtUTC) || validateInventory(o.InventoryEntries) != nil {
		return errors.New("invalid integration observation")
	}
	if o.WorktreeLocator == "" {
		if o.Ref != "" || o.OID != "" || o.Tree != "" || o.GitCommonDir != "" || o.ObjectFormat != "" || o.RefFormat != "" || o.Symbolic || o.Clean || o.StatusEntries == nil || o.InProgress == nil {
			return errors.New("invalid unknown integration observation")
		}
		return nil
	}
	if validatePlanObservation(PlanWorktreeObservation{WorktreeLocator: o.WorktreeLocator, Ref: o.Ref, OID: o.OID, Tree: o.Tree, GitCommonDir: o.GitCommonDir, ObjectFormat: o.ObjectFormat, RefFormat: o.RefFormat, Symbolic: o.Symbolic, Clean: o.Clean, StatusEntries: o.StatusEntries, InProgress: o.InProgress}) != nil {
		return errors.New("invalid integration observation")
	}
	return nil
}
func sameIntegrationObservation(a, b IntegrationObservation) bool {
	a.ObservedAtUTC, b.ObservedAtUTC = "", ""
	return reflect.DeepEqual(a, b)
}
func validateStatusEntries(entries []IntegrationStatusEntry) error {
	if entries == nil {
		return errors.New("status entries must be explicit")
	}
	var previous []byte
	for i, e := range entries {
		path, err := base64.StdEncoding.DecodeString(e.PathBase64)
		if err != nil || len(path) == 0 || (i > 0 && bytes.Compare(previous, path) >= 0) || !setString("ordinary", "renamed", "copied", "unmerged", "untracked")[e.RecordKind] {
			return errors.New("invalid status entry")
		}
		rename := e.RecordKind == "renamed" || e.RecordKind == "copied"
		if rename != (e.OriginalPathBase64 != nil) || (rename && func() bool { _, err := base64.StdEncoding.DecodeString(*e.OriginalPathBase64); return err != nil }()) {
			return errors.New("invalid status original path")
		}
		untracked := e.RecordKind == "untracked"
		if untracked != (e.IndexState == nil && e.WorktreeState == nil && e.SubmoduleState == nil) {
			return errors.New("invalid status state")
		}
		if !untracked && (e.IndexState == nil || e.WorktreeState == nil || e.SubmoduleState == nil || len(*e.IndexState) != 1 || len(*e.WorktreeState) != 1 || len(*e.SubmoduleState) != 4) {
			return errors.New("invalid status state")
		}
		previous = path
	}
	return nil
}
func validateInventory(entries []IntegrationInventoryEntry) error {
	if entries == nil {
		return errors.New("inventory must be explicit")
	}
	for i, e := range entries {
		if !validAbsoluteCleanPath(e.Locator) || !validOIDText(e.OID) || (!e.Bare && !e.Detached && !validFullBranchRef(e.Ref)) || (i > 0 && !inventoryLess(entries[i-1], e)) {
			return errors.New("invalid inventory")
		}
	}
	return nil
}
func inventoryLess(a, b IntegrationInventoryEntry) bool {
	if a.Locator != b.Locator {
		return a.Locator < b.Locator
	}
	if a.Ref != b.Ref {
		return a.Ref < b.Ref
	}
	return a.OID < b.OID
}
func validateReflog(entries []IntegrationReflogEntry) error {
	// A new parent has only its creation entry; exact effects use reflogProves.
	if len(entries) == 0 {
		return errors.New("reflog must contain at least one entry")
	}
	for i, e := range entries {
		if e.Ordinal != i || !validOIDText(e.OID) || !validTaskText(e.Selector, 1, 512) || !validTaskText(e.Action, 1, 2000) {
			return errors.New("invalid reflog")
		}
	}
	return nil
}
func sortedReasonTokens(reasons []string) bool {
	if reasons == nil {
		return false
	}
	for i, reason := range reasons {
		if !regexp.MustCompile(`^[a-z0-9_]+$`).MatchString(reason) || (i > 0 && reasons[i-1] >= reason) {
			return false
		}
	}
	return true
}
func sameAllowedEffect(a IntegrationAllowedEffect, e IntegrationPlanEffect) bool {
	return a.Kind == e.Kind && a.ParentRef == e.ParentRef && a.ExpectedParentOID == e.ExpectedParentOID && a.ResultOID == e.ResultOID && a.MaxOccurrences == e.MaxOccurrences
}
func sameOptionalIntegrationResultID(a, b *IntegrationResultID) bool {
	return (a == nil && b == nil) || (a != nil && b != nil && *a == *b)
}
func intentMatchesAuthority(i IntegrationIntent, a IntegrationAuthority) bool {
	p := a.Plan
	return i.WorkspaceRoot == p.Workspace.Root && i.WorkspaceMarkerSHA256 == p.Workspace.MarkerSHA256 && i.ProjectID == p.Project.ProjectID && i.RepoID == p.Repository.RepoID && i.RegisteredRepoLocator == p.Repository.RegisteredLocator && i.GitCommonDir == p.Repository.GitCommonDir && i.ObjectFormat == p.Repository.ObjectFormat && i.RefFormat == p.Repository.RefFormat && i.EpicID == p.Epic.EpicID && i.ParentWorktreeID == p.Epic.ParentWorktreeID && i.ParentLocator == p.Epic.ParentLocator && i.ParentRef == p.Epic.ParentRef && i.ExpectedParentOID == p.Epic.ExpectedParentOID && i.ExpectedParentTree == p.Epic.ExpectedParentTree && i.TaskID == p.Task.TaskID && i.TaskWorktreeID == p.Task.TaskWorktreeID && i.SourceLocator == p.Task.SourceLocator && i.SourceRef == p.Task.SourceRef && i.ResultOID == p.Task.ResultOID && i.ResultTree == p.Task.ResultTree && sameAllowedEffect(i.Effect, p.Effect)
}
func validateIntegrationResultShape(r IntegrationResult) error {
	if !setString("no_effect", "exact_effect", "already_integrated", "blocked", "conflict", "partial", "unknown")[r.Outcome] || !validTaskUTC(r.RecordedAtUTC) || validateIntegrationObservation(r.BeforeObservation) != nil || validateIntegrationObservation(r.AfterObservation) != nil || validateResultReflog(r.Reflog) != nil || !setString("complete", "safe-no-effect", "blocked", "conflict", "reconciliation-required", "unknown")[r.RecoveryStatus] || !validNextAction(r.NextAction) {
		return errors.New("invalid integration result")
	}
	requiresAttempt := setString("no_effect", "exact_effect", "partial", "unknown")[r.Outcome]
	if requiresAttempt != (r.AttemptID != nil) {
		return errors.New("invalid result attempt nullform")
	}
	if (r.Outcome == "exact_effect" && (r.GitChanged == nil || !*r.GitChanged)) || (setString("no_effect", "already_integrated", "blocked", "conflict")[r.Outcome] && (r.GitChanged == nil || *r.GitChanged)) || (setString("partial", "unknown")[r.Outcome] && r.GitChanged != nil) {
		return errors.New("invalid git_changed")
	}
	if (r.Command.Launched && (r.Command.ExitCode == nil || r.Command.LaunchError != nil || *r.Command.ExitCode < 0)) || (!r.Command.Launched && r.Command.ExitCode != nil) || (r.Command.LaunchError != nil && *r.Command.LaunchError != "launch_failed") || validateCommandStream(r.Command.Stdout) != nil || validateCommandStream(r.Command.Stderr) != nil {
		return errors.New("invalid command evidence")
	}
	return nil
}

func validateIntegrationResultOutcomeContract(result IntegrationResult, authority IntegrationAuthority) error {
	if authority.historical != nil {
		return nil
	}
	emptyArgv := len(result.NextAction.Argv) == 0
	switch result.Outcome {
	case "exact_effect", "already_integrated":
		if result.RecoveryStatus != "complete" || result.NextAction.Kind != "none" || !emptyArgv {
			return errors.New("successful integration result has contradictory recovery guidance")
		}
	case "no_effect":
		want := integrationNoEffectNextAction(authority, result.ID).Argv
		if result.RecoveryStatus != "safe-no-effect" || result.NextAction.Kind != "retry_after_no_effect" || !slices.Equal(result.NextAction.Argv, want) {
			return errors.New("no-effect integration result has contradictory recovery guidance")
		}
	case "blocked":
		if result.RecoveryStatus != "blocked" || result.NextAction.Kind != "resolve_blocker" || !emptyArgv {
			return errors.New("blocked integration result has contradictory recovery guidance")
		}
	case "conflict":
		if result.RecoveryStatus != "conflict" || result.NextAction.Kind != "refresh_or_reverify_parent" || !emptyArgv {
			return errors.New("conflicting integration result has contradictory recovery guidance")
		}
	case "partial", "unknown":
		wantRecovery := "reconciliation-required"
		if result.Outcome == "unknown" {
			wantRecovery = "unknown"
		}
		want := []string{"ply", "workspace", "task", "show", string(authority.TaskID), "--format", "json"}
		if result.RecoveryStatus != wantRecovery || result.NextAction.Kind != "read_only_recovery_control" || !slices.Equal(result.NextAction.Argv, want) {
			return errors.New("unresolved integration result has contradictory recovery guidance")
		}
	}
	if result.AttemptID == nil && !reflect.DeepEqual(result.Command, emptyCommandEvidence()) {
		return errors.New("integration result without an Attempt has command evidence")
	}
	return nil
}

func validateCommandStream(stream IntegrationCommandStream) error {
	if stream.SizeBytes < 0 || stream.CapturedBytes < 0 || stream.CapturedBytes > 64<<10 || stream.CapturedBytes > stream.SizeBytes || !digestPattern.MatchString(stream.SHA256) {
		return errors.New("invalid command stream sizes or digest")
	}
	decoded, err := base64.StdEncoding.DecodeString(stream.Base64)
	if err != nil || int64(len(decoded)) != stream.CapturedBytes {
		return errors.New("invalid command stream base64")
	}
	wantCaptured := stream.SizeBytes
	if wantCaptured > 64<<10 {
		wantCaptured = 64 << 10
	}
	if stream.CapturedBytes != wantCaptured || stream.Truncated != (stream.SizeBytes > stream.CapturedBytes) {
		return errors.New("invalid command stream truncation")
	}
	if !stream.Truncated && digestTaskBytes(decoded) != stream.SHA256 {
		return errors.New("invalid command stream digest")
	}
	return nil
}
func validateResultReflog(entries []IntegrationReflogEntry) error {
	if entries == nil {
		return errors.New("result reflog must be explicit")
	}
	for i, e := range entries {
		if e.Ordinal != i || !validOIDText(e.OID) || !validTaskText(e.Selector, 1, 512) || !validTaskText(e.Action, 1, 2000) {
			return errors.New("invalid result reflog")
		}
	}
	return nil
}
func validNextAction(a IntegrationNextAction) bool {
	return setString("record_task_result", "record_human_qa", "apply_confirmed_plan", "new_check", "retry_after_no_effect", "refresh_or_reverify_parent", "resolve_blocker", "read_only_recovery_control", "none")[a.Kind] && validTaskText(a.Reason, 1, 2000) && a.Argv != nil
}
func integrationResultByID(results []IntegrationResult, id IntegrationResultID) (IntegrationResult, bool) {
	for _, result := range results {
		if result.ID == id {
			return result, true
		}
	}
	return IntegrationResult{}, false
}

func sortedAuthorities(v []IntegrationAuthority) bool {
	for i := 1; i < len(v); i++ {
		if v[i-1].ID >= v[i].ID {
			return false
		}
	}
	return true
}
func sortedIntents(v []IntegrationIntent) bool {
	for i := 1; i < len(v); i++ {
		if v[i-1].ID >= v[i].ID {
			return false
		}
	}
	return true
}
func sortedAttempts(v []IntegrationAttempt) bool {
	for i := 1; i < len(v); i++ {
		if v[i-1].ID >= v[i].ID {
			return false
		}
	}
	return true
}
func sortedIntegrationResults(v []IntegrationResult) bool {
	for i := 1; i < len(v); i++ {
		if v[i-1].ID >= v[i].ID {
			return false
		}
	}
	return true
}

// firstInvalidField rejects the first present field, by name, whose value the
// rule does not accept, so that several invalid fields always give one answer.
func firstInvalidField(fields map[string]*string, valid func(string) bool) error {
	for _, name := range slices.Sorted(maps.Keys(fields)) {
		if value := fields[name]; value != nil && !valid(*value) {
			return fmt.Errorf("invalid %s", name)
		}
	}
	return nil
}

func validateObservation(observation WorktreeObservationRecord, state string) error {
	if observation.Reasons == nil || !sort.StringsAreSorted(observation.Reasons) {
		return errors.New("observation reasons must be an explicit sorted array")
	}
	for i := 1; i < len(observation.Reasons); i++ {
		if observation.Reasons[i] == observation.Reasons[i-1] {
			return errors.New("observation reasons must be unique")
		}
	}
	knownReasons := map[string]bool{"task_worktree_reconciliation_required": true, "parent_ref_stale": true, "parent_worktree_stale": true, "parent_worktree_dirty": true, "parent_observation_unknown": true, "source_ref_stale": true, "source_ref_checked_out_elsewhere": true, "target_worktree_stale": true, "target_worktree_dirty": true, "target_observation_unknown": true, "worktree_inventory_stale": true, "worktree_inventory_unknown": true}
	for _, reason := range observation.Reasons {
		if !knownReasons[reason] {
			return fmt.Errorf("unknown observation reason %s", reason)
		}
	}
	if observation.TargetKind != "absent" && observation.TargetKind != "worktree" && observation.TargetKind != "other" && observation.TargetKind != "unknown" {
		return fmt.Errorf("invalid target kind %s", observation.TargetKind)
	}
	if err := firstInvalidField(map[string]*string{"parent_ref_oid": observation.ParentRefOID, "parent_ref_tree": observation.ParentRefTree, "parent_worktree_oid": observation.ParentWorktreeOID, "parent_worktree_tree": observation.ParentWorktreeTree, "source_ref_oid": observation.SourceRefOID, "source_ref_tree": observation.SourceRefTree, "target_oid": observation.TargetOID, "target_tree": observation.TargetTree}, validOIDText); err != nil {
		return err
	}
	if err := firstInvalidField(map[string]*string{"parent_worktree_locator": observation.ParentWorktreeLocator, "parent_worktree_git_common_dir": observation.ParentWorktreeGitCommonDir, "source_ref_checked_out_at": observation.SourceRefCheckedOutAt, "target_locator": observation.TargetLocator, "target_git_common_dir": observation.TargetGitCommonDir}, workspace.ValidStoredPath); err != nil {
		return err
	}
	if err := firstInvalidField(map[string]*string{"parent_worktree_ref": observation.ParentWorktreeRef, "target_ref": observation.TargetRef}, validFullBranchRef); err != nil {
		return err
	}
	switch state {
	case "creating":
		if observation.Classification != "no_effect" && observation.Classification != "branch_only" {
			return errors.New("creating operation has invalid classification")
		}
		if len(observation.Reasons) != 0 || observation.TargetKind != "absent" || observation.TargetLocator != nil || observation.TargetRef != nil || observation.TargetOID != nil || observation.TargetTree != nil || observation.TargetGitCommonDir != nil || observation.TargetClean != nil || observation.InventoryMatch != nil {
			return errors.New("creating operation has inconsistent safe observation")
		}
		if observation.Classification == "no_effect" && (observation.SourceRefOID != nil || observation.SourceRefTree != nil || observation.SourceRefCheckedOutAt != nil) {
			return errors.New("no_effect observation contains a source ref")
		}
		if observation.Classification == "branch_only" && (observation.SourceRefOID == nil || observation.SourceRefTree == nil || observation.SourceRefCheckedOutAt != nil) {
			return errors.New("branch_only observation does not contain an unbound source ref")
		}
	case "ready":
		if observation.Classification != "exact_effect" || len(observation.Reasons) != 0 || observation.TargetKind != "worktree" || observation.ParentRefOID == nil || observation.ParentRefTree == nil || observation.ParentWorktreeLocator == nil || observation.ParentWorktreeRef == nil || observation.ParentWorktreeOID == nil || observation.ParentWorktreeTree == nil || observation.ParentWorktreeGitCommonDir == nil || observation.ParentWorktreeClean == nil || observation.SourceRefOID == nil || observation.SourceRefTree == nil || observation.SourceRefCheckedOutAt == nil || observation.TargetLocator == nil || observation.TargetRef == nil || observation.TargetOID == nil || observation.TargetTree == nil || observation.TargetGitCommonDir == nil || observation.TargetClean == nil || observation.InventoryMatch == nil {
			return errors.New("ready operation must have exact_effect observation")
		}
	case "reconciliation_required":
		if observation.Classification != "partial_or_unknown" {
			return errors.New("reconciliation operation must have partial_or_unknown observation")
		}
	}
	return nil
}
func sortedUniqueEpics(values []EpicRecord) bool {
	for i := 1; i < len(values); i++ {
		if values[i-1].ID >= values[i].ID {
			return false
		}
	}
	return true
}
func sortedUniqueTasks(values []TaskRecord) bool {
	for i := 1; i < len(values); i++ {
		if values[i-1].ID >= values[i].ID {
			return false
		}
	}
	return true
}
func sortedUniqueOperations(values []WorktreeOperationRecord) bool {
	for i := 1; i < len(values); i++ {
		if values[i-1].ID >= values[i].ID {
			return false
		}
	}
	return true
}
