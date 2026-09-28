package workspace

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"gopkg.in/yaml.v3"
)

type WorkItemStore interface {
	Snapshot(string) (WorkItemRegistry, error)
	WithLock(string, func(WorkItemStoreSession) error) error
}

type WorkItemStoreSession interface {
	Snapshot() (WorkItemRegistry, error)
	Publish(WorkItemRegistry) error
}

type workItemStoreFaults struct {
	fail       func(string) error
	shortWrite bool
}

type systemWorkItemStore struct{ faults *workItemStoreFaults }
type systemWorkItemStoreSession struct {
	store *systemWorkItemStore
	root  string
}

func newSystemWorkItemStore() *systemWorkItemStore { return &systemWorkItemStore{} }
func emptyWorkItemRegistry() WorkItemRegistry {
	return WorkItemRegistry{FormatVersion: FormatVersion, Epics: []EpicRecord{}, Tasks: []TaskRecord{}, WorktreeOperations: []WorktreeOperationRecord{}}
}
func workItemsPath(root string) string { return filepath.Join(root, MarkerDirectory, WorkItemsFile) }

func (store *systemWorkItemStore) Snapshot(root string) (WorkItemRegistry, error) {
	return store.read(root)
}

func (store *systemWorkItemStore) WithLock(root string, operation func(WorkItemStoreSession) error) error {
	lock, err := store.acquire(root)
	if err != nil {
		return err
	}
	opErr := operation(&systemWorkItemStoreSession{store: store, root: root})
	releaseErr := store.release(lock)
	if opErr != nil {
		if releaseErr != nil {
			return workError(ErrorWorkIO, fmt.Sprintf("release after failure: %v; %v", opErr, releaseErr), opErr)
		}
		return opErr
	}
	return releaseErr
}

func (session *systemWorkItemStoreSession) Snapshot() (WorkItemRegistry, error) {
	return session.store.read(session.root)
}
func (session *systemWorkItemStoreSession) Publish(registry WorkItemRegistry) error {
	return session.store.publish(session.root, registry)
}

func (store *systemWorkItemStore) acquire(root string) (*os.File, error) {
	path := filepath.Join(root, MarkerDirectory, workItemsLock)
	for {
		info, err := os.Lstat(path)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, workError(ErrorWorkStoreConflict, fmt.Sprintf("cannot inspect work-item lock %s", path), err)
		}
		if err == nil && (!info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0) {
			return nil, workError(ErrorWorkStoreConflict, fmt.Sprintf("work-item lock %s is not a regular file", path), nil)
		}
		flags := os.O_RDWR
		if errors.Is(err, fs.ErrNotExist) {
			flags |= os.O_CREATE | os.O_EXCL
		}
		if fault := store.fail("lock-open"); fault != nil {
			return nil, workError(ErrorWorkIO, "open work-item lock", fault)
		}
		file, openErr := os.OpenFile(path, flags, 0o600)
		if openErr != nil && flags&os.O_EXCL != 0 && errors.Is(openErr, fs.ErrExist) {
			continue
		}
		if openErr != nil {
			return nil, workError(ErrorWorkIO, fmt.Sprintf("open work-item lock %s", path), openErr)
		}
		opened, statErr := file.Stat()
		current, pathErr := os.Lstat(path)
		if statErr != nil || pathErr != nil || !opened.Mode().IsRegular() || !current.Mode().IsRegular() || current.Mode()&fs.ModeSymlink != 0 || !os.SameFile(opened, current) {
			_ = file.Close()
			return nil, workError(ErrorWorkStoreConflict, fmt.Sprintf("work-item lock %s changed or is not a regular file", path), firstError(statErr, pathErr))
		}
		if fault := store.fail("lock"); fault != nil {
			_ = file.Close()
			return nil, workError(ErrorWorkIO, "lock work-item registry", fault)
		}
		if err := platformLock(file); err != nil {
			_ = file.Close()
			return nil, workError(ErrorWorkIO, "lock work-item registry", err)
		}
		return file, nil
	}
}

func (store *systemWorkItemStore) release(file *os.File) error {
	unlockErr := platformUnlock(file)
	if fault := store.fail("unlock"); fault != nil && unlockErr == nil {
		unlockErr = fault
	}
	closeErr := file.Close()
	if fault := store.fail("lock-close"); fault != nil && closeErr == nil {
		closeErr = fault
	}
	if unlockErr != nil {
		return workError(ErrorWorkIO, fmt.Sprintf("unlock work-item lock %s", file.Name()), unlockErr)
	}
	if closeErr != nil {
		return workError(ErrorWorkIO, fmt.Sprintf("close work-item lock %s", file.Name()), closeErr)
	}
	return nil
}

func (store *systemWorkItemStore) read(root string) (WorkItemRegistry, error) {
	path := workItemsPath(root)
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return emptyWorkItemRegistry(), nil
	}
	if err != nil {
		return WorkItemRegistry{}, workError(ErrorWorkStoreConflict, fmt.Sprintf("cannot inspect work-item registry %s", path), err)
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
		return WorkItemRegistry{}, workError(ErrorWorkStoreConflict, fmt.Sprintf("work-item registry %s is not a regular file", path), nil)
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return WorkItemRegistry{}, workError(ErrorWorkStoreConflict, fmt.Sprintf("cannot read work-item registry %s", path), err)
	}
	registry, err := decodeWorkItemRegistry(contents)
	if err != nil {
		return WorkItemRegistry{}, workError(ErrorWorkStoreConflict, fmt.Sprintf("invalid work-item registry %s: %v", path, err), err)
	}
	return registry, nil
}

func (store *systemWorkItemStore) publish(root string, registry WorkItemRegistry) error {
	contents, err := encodeWorkItemRegistry(registry)
	if err != nil {
		return workError(ErrorWorkStoreConflict, fmt.Sprintf("cannot encode work-item registry: %v", err), err)
	}
	directory := filepath.Join(root, MarkerDirectory)
	destination := workItemsPath(root)
	if fault := store.fail("temp-open"); fault != nil {
		return workError(ErrorWorkIO, "open work-item temporary file", fault)
	}
	temporary, err := os.CreateTemp(directory, ".work-items-*.tmp")
	if err != nil {
		return workError(ErrorWorkIO, "open work-item temporary file", err)
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
		return workError(ErrorWorkIO, operation+" work-item registry", primary)
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
	if err := temporary.Close(); err != nil {
		return abort("close", err)
	}
	if fault := store.fail("close"); fault != nil {
		return abort("close", fault)
	}
	if fault := store.fail("replace"); fault != nil {
		return abort("replace", fault)
	}
	if err := platformReplace(temporaryPath, destination); err != nil {
		return abort("replace", err)
	}
	replaced = true
	if fault := store.fail("directory-sync"); fault != nil {
		return workError(ErrorWorkIO, "sync work-item registry directory", fault)
	}
	if err := platformSyncDirectory(directory); err != nil {
		return workError(ErrorWorkIO, "sync work-item registry directory", err)
	}
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
	return output.Bytes(), nil
}

func decodeWorkItemRegistry(contents []byte) (WorkItemRegistry, error) {
	if len(contents) == 0 {
		return WorkItemRegistry{}, errors.New("registry is empty")
	}
	decoder := yaml.NewDecoder(bytes.NewReader(contents))
	decoder.KnownFields(true)
	var registry WorkItemRegistry
	if err := decoder.Decode(&registry); err != nil {
		return WorkItemRegistry{}, err
	}
	var extra yaml.Node
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return WorkItemRegistry{}, errors.New("registry contains multiple YAML documents")
		}
		return WorkItemRegistry{}, err
	}
	if registry.Epics == nil || registry.Tasks == nil || registry.WorktreeOperations == nil {
		return WorkItemRegistry{}, errors.New("epics, tasks, and worktree_operations must be explicit arrays")
	}
	if err := validateWorkItemRegistry(registry); err != nil {
		return WorkItemRegistry{}, err
	}
	return registry, nil
}

var worktreeIDPattern = regexp.MustCompile(`^wt_[0-9a-f]{32}$`)
var operationIDPattern = regexp.MustCompile(`^wop_[0-9a-f]{32}$`)
var digestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

func validateWorkItemRegistry(registry WorkItemRegistry) error {
	if registry.FormatVersion != FormatVersion {
		return fmt.Errorf("unsupported format_version %d", registry.FormatVersion)
	}
	if !sortedUniqueEpics(registry.Epics) || !sortedUniqueTasks(registry.Tasks) || !sortedUniqueOperations(registry.WorktreeOperations) {
		return errors.New("work-item arrays must be sorted by unique IDs")
	}
	epics := map[EpicID]EpicRecord{}
	worktreeIDs := map[WorktreeID]string{}
	resources := map[string]string{}
	for _, epic := range registry.Epics {
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
		seenRepos := map[RepoID]bool{}
		for _, binding := range epic.RepoBindings {
			if seenRepos[binding.RepoID] {
				return fmt.Errorf("Epic %s has duplicate repository %s", epic.ID, binding.RepoID)
			}
			seenRepos[binding.RepoID] = true
			if err := validateWorkIdentifier("repository", string(binding.RepoID)); err != nil {
				return err
			}
			if !validStoredPath(binding.GitCommonDir) {
				return fmt.Errorf("Epic %s has invalid common directory", epic.ID)
			}
			wt := binding.Worktree
			if !worktreeIDPattern.MatchString(string(wt.ID)) || wt.OwnerKind != "epic" || wt.OwnerID != epic.ID || wt.Origin != "adopted" || wt.GitCommonDir != binding.GitCommonDir || !validStoredPath(wt.Locator) || !validFullBranchRef(wt.Ref) || !validOIDText(wt.OID) || !validOIDText(wt.Tree) {
				return fmt.Errorf("Epic %s has invalid worktree binding", epic.ID)
			}
			if owner, ok := worktreeIDs[wt.ID]; ok {
				return fmt.Errorf("worktree ID %s is shared by %s and Epic %s", wt.ID, owner, epic.ID)
			}
			worktreeIDs[wt.ID] = "Epic " + string(epic.ID)
			for _, key := range []string{binding.GitCommonDir + "\x00path\x00" + wt.Locator, binding.GitCommonDir + "\x00ref\x00" + wt.Ref} {
				if owner, ok := resources[key]; ok {
					return fmt.Errorf("Git resource shared by %s and Epic %s", owner, epic.ID)
				}
				resources[key] = "Epic " + string(epic.ID)
			}
		}
		epics[epic.ID] = epic
	}
	operations := map[TaskID]WorktreeOperationRecord{}
	operationIDs := map[WorktreeOperationID]bool{}
	for _, operation := range registry.WorktreeOperations {
		if !operationIDPattern.MatchString(string(operation.ID)) || operationIDs[operation.ID] {
			return fmt.Errorf("invalid or duplicate operation ID %s", operation.ID)
		}
		operationIDs[operation.ID] = true
		if _, ok := operations[operation.TaskID]; ok {
			return fmt.Errorf("Task %s has multiple operations", operation.TaskID)
		}
		if !worktreeIDPattern.MatchString(string(operation.WorktreeID)) || !digestPattern.MatchString(operation.IntentDigest) || !validStoredPath(operation.GitCommonDir) || !validStoredPath(operation.TargetLocator) || !validFullBranchRef(operation.ParentRef) || !validFullBranchRef(operation.SourceRef) || !validOIDText(operation.ParentOID) || !validOIDText(operation.ParentTree) {
			return fmt.Errorf("operation %s has invalid intent", operation.ID)
		}
		if owner, duplicate := worktreeIDs[operation.WorktreeID]; duplicate {
			return fmt.Errorf("worktree ID %s is shared by %s and operation %s", operation.WorktreeID, owner, operation.ID)
		}
		worktreeIDs[operation.WorktreeID] = "operation " + string(operation.ID)
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
		for _, key := range []string{operation.GitCommonDir + "\x00path\x00" + operation.TargetLocator, operation.GitCommonDir + "\x00ref\x00" + operation.SourceRef} {
			if owner, ok := resources[key]; ok {
				return fmt.Errorf("Git resource shared by %s and operation %s", owner, operation.ID)
			}
			resources[key] = "operation " + string(operation.ID)
		}
		operations[operation.TaskID] = operation
	}
	for _, task := range registry.Tasks {
		if err := validateWorkIdentifier("Task", string(task.ID)); err != nil {
			return err
		}
		if err := validateWorkText("Task title", task.Title, 128); err != nil {
			return err
		}
		if err := validateWorkText("Task description", task.Description, 2048); err != nil {
			return err
		}
		epic, ok := epics[task.ParentEpicID]
		if !ok || epic.ProjectID != task.ProjectID {
			return fmt.Errorf("Task %s has invalid Epic/Project reference", task.ID)
		}
		binding, _ := findEpicRepo(epic, task.RepoID)
		if binding == nil || binding.GitCommonDir != task.GitCommonDir {
			return fmt.Errorf("Task %s has invalid repository binding", task.ID)
		}
		operation, hasOperation := operations[task.ID]
		if hasOperation && (operation.ProjectID != task.ProjectID || operation.RepoID != task.RepoID || operation.GitCommonDir != task.GitCommonDir || operation.ParentEpicID != task.ParentEpicID || operation.ParentWorktreeID != binding.Worktree.ID || operation.ParentRef != binding.Worktree.Ref || operation.ParentOID != binding.Worktree.OID || operation.ParentTree != binding.Worktree.Tree) {
			return fmt.Errorf("Task %s operation differs from its parent repository binding", task.ID)
		}
		switch task.WorktreeState {
		case WorkItemUnbound:
			if task.Worktree != nil || hasOperation {
				return fmt.Errorf("unbound Task %s has worktree state", task.ID)
			}
		case WorkItemCreating:
			if task.Worktree != nil || !hasOperation || operation.State != "creating" {
				return fmt.Errorf("Task %s has incomplete operation state", task.ID)
			}
		case WorkItemReconciliationRequired:
			if task.Worktree != nil || !hasOperation || operation.State != "reconciliation_required" {
				return fmt.Errorf("Task %s has incomplete operation state", task.ID)
			}
		case WorkItemReady:
			if task.Worktree == nil || !hasOperation || operation.State != "ready" {
				return fmt.Errorf("ready Task %s has invalid binding", task.ID)
			}
			wt := task.Worktree
			if wt.ID != operation.WorktreeID || wt.OperationID != operation.ID || wt.IntentDigest != operation.IntentDigest || wt.OwnerKind != "task" || wt.OwnerID != task.ID || wt.Origin != "ply_created" || wt.Locator != operation.TargetLocator || wt.Ref != operation.SourceRef || wt.OID != operation.ParentOID || wt.Tree != operation.ParentTree || wt.GitCommonDir != operation.GitCommonDir || wt.ParentEpicID != task.ParentEpicID || wt.ParentWorktreeID != operation.ParentWorktreeID || wt.ParentRef != operation.ParentRef || wt.ParentOID != operation.ParentOID || wt.ParentTree != operation.ParentTree {
				return fmt.Errorf("ready Task %s binding differs from operation", task.ID)
			}
			if owner, ok := worktreeIDs[wt.ID]; ok && owner != "operation "+string(operation.ID) {
				return fmt.Errorf("worktree ID %s reused", wt.ID)
			}
			observed := operation.LastObservation
			if *observed.ParentRefOID != operation.ParentOID || *observed.ParentRefTree != operation.ParentTree || *observed.ParentWorktreeLocator != binding.Worktree.Locator || *observed.ParentWorktreeRef != operation.ParentRef || *observed.ParentWorktreeOID != operation.ParentOID || *observed.ParentWorktreeTree != operation.ParentTree || *observed.ParentWorktreeGitCommonDir != operation.GitCommonDir || !*observed.ParentWorktreeClean || *observed.SourceRefOID != operation.ParentOID || *observed.SourceRefTree != operation.ParentTree || *observed.SourceRefCheckedOutAt != operation.TargetLocator || *observed.TargetLocator != operation.TargetLocator || *observed.TargetRef != operation.SourceRef || *observed.TargetOID != operation.ParentOID || *observed.TargetTree != operation.ParentTree || *observed.TargetGitCommonDir != operation.GitCommonDir || !*observed.TargetClean || !*observed.InventoryMatch {
				return fmt.Errorf("ready Task %s exact observation differs from operation", task.ID)
			}
		default:
			return fmt.Errorf("Task %s has invalid worktree_state", task.ID)
		}
	}
	for taskID := range operations {
		found := false
		for _, task := range registry.Tasks {
			if task.ID == taskID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("operation refers to missing Task %s", taskID)
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
	for name, value := range map[string]*string{"parent_ref_oid": observation.ParentRefOID, "parent_ref_tree": observation.ParentRefTree, "parent_worktree_oid": observation.ParentWorktreeOID, "parent_worktree_tree": observation.ParentWorktreeTree, "source_ref_oid": observation.SourceRefOID, "source_ref_tree": observation.SourceRefTree, "target_oid": observation.TargetOID, "target_tree": observation.TargetTree} {
		if value != nil && !validOIDText(*value) {
			return fmt.Errorf("invalid %s", name)
		}
	}
	for name, value := range map[string]*string{"parent_worktree_locator": observation.ParentWorktreeLocator, "parent_worktree_git_common_dir": observation.ParentWorktreeGitCommonDir, "source_ref_checked_out_at": observation.SourceRefCheckedOutAt, "target_locator": observation.TargetLocator, "target_git_common_dir": observation.TargetGitCommonDir} {
		if value != nil && !validStoredPath(*value) {
			return fmt.Errorf("invalid %s", name)
		}
	}
	for name, value := range map[string]*string{"parent_worktree_ref": observation.ParentWorktreeRef, "target_ref": observation.TargetRef} {
		if value != nil && !validFullBranchRef(*value) {
			return fmt.Errorf("invalid %s", name)
		}
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
