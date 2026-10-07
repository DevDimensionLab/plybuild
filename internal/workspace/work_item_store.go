package workspace

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type WorkItemStore interface {
	Snapshot(string) (WorkItemRegistry, error)
	SnapshotRegistrations(string) (WorkItemRegistry, error)
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

type workItemRegistryV1Wire struct {
	FormatVersion      int                       `yaml:"format_version"`
	Epics              []EpicRecord              `yaml:"epics"`
	Tasks              []TaskRecord              `yaml:"tasks"`
	WorktreeOperations []WorktreeOperationRecord `yaml:"worktree_operations"`
}
type workItemRegistryV2Wire struct {
	FormatVersion          int                       `yaml:"format_version"`
	Epics                  []EpicRecord              `yaml:"epics"`
	Tasks                  []TaskRecord              `yaml:"tasks"`
	WorktreeOperations     []WorktreeOperationRecord `yaml:"worktree_operations"`
	TaskResults            []TaskResultRecord        `yaml:"task_results"`
	HumanQARecords         []TaskHumanQARecord       `yaml:"human_qa_records"`
	IntegrationAuthorities []IntegrationAuthority    `yaml:"integration_authorities"`
	IntegrationIntents     []IntegrationIntent       `yaml:"integration_intents"`
	IntegrationAttempts    []IntegrationAttempt      `yaml:"integration_attempts"`
	IntegrationResults     []IntegrationResult       `yaml:"integration_results"`
}
type workItemRegistryV3Wire struct {
	FormatVersion           int                       `yaml:"format_version"`
	Epics                   []EpicRecord              `yaml:"epics"`
	Tasks                   []TaskRecord              `yaml:"tasks"`
	WorktreeOperations      []WorktreeOperationRecord `yaml:"worktree_operations"`
	TaskResults             []TaskResultRecord        `yaml:"task_results"`
	HumanQARecords          []TaskHumanQARecord       `yaml:"human_qa_records"`
	IntegrationAuthorities  []IntegrationAuthority    `yaml:"integration_authorities"`
	IntegrationIntents      []IntegrationIntent       `yaml:"integration_intents"`
	IntegrationAttempts     []IntegrationAttempt      `yaml:"integration_attempts"`
	IntegrationResults      []IntegrationResult       `yaml:"integration_results"`
	TaskSpecPolicies        []TaskSpecPolicy          `yaml:"task_spec_policies"`
	TaskProblemRevisions    []TaskProblemReference    `yaml:"task_problem_revisions"`
	TaskSpecRevisions       []TaskSpecReference       `yaml:"task_spec_revisions"`
	TaskSpecAssessments     []TaskAssessmentReference `yaml:"task_spec_assessments"`
	TaskSolutionSelections  []TaskSelectionReference  `yaml:"task_solution_selections"`
	TaskResultSpecBindings  []TaskResultSpecBinding   `yaml:"task_result_spec_bindings"`
	TaskContentPublications []TaskContentPublication  `yaml:"task_content_publications"`
}

func (registry WorkItemRegistry) MarshalYAML() (any, error) {
	if registry.FormatVersion == 4 {
		return workItemRegistryV4(registry), nil
	}
	if registry.FormatVersion == 1 {
		return workItemRegistryV1Wire{FormatVersion: 1, Epics: registry.Epics, Tasks: registry.Tasks, WorktreeOperations: registry.WorktreeOperations}, nil
	}
	if registry.FormatVersion == 3 {
		return workItemRegistryV3Wire{FormatVersion: registry.FormatVersion, Epics: registry.Epics, Tasks: registry.Tasks, WorktreeOperations: registry.WorktreeOperations, TaskResults: registry.TaskResults, HumanQARecords: registry.HumanQARecords, IntegrationAuthorities: registry.IntegrationAuthorities, IntegrationIntents: registry.IntegrationIntents, IntegrationAttempts: registry.IntegrationAttempts, IntegrationResults: registry.IntegrationResults, TaskSpecPolicies: registry.TaskSpecPolicies, TaskProblemRevisions: registry.TaskProblemRevisions, TaskSpecRevisions: registry.TaskSpecRevisions, TaskSpecAssessments: registry.TaskSpecAssessments, TaskSolutionSelections: registry.TaskSolutionSelections, TaskResultSpecBindings: registry.TaskResultSpecBindings, TaskContentPublications: registry.TaskContentPublications}, nil
	}
	return workItemRegistryV2Wire{FormatVersion: registry.FormatVersion, Epics: registry.Epics, Tasks: registry.Tasks, WorktreeOperations: registry.WorktreeOperations, TaskResults: registry.TaskResults, HumanQARecords: registry.HumanQARecords, IntegrationAuthorities: registry.IntegrationAuthorities, IntegrationIntents: registry.IntegrationIntents, IntegrationAttempts: registry.IntegrationAttempts, IntegrationResults: registry.IntegrationResults}, nil
}

func newSystemWorkItemStore() *systemWorkItemStore { return &systemWorkItemStore{} }
func emptyWorkItemRegistry() WorkItemRegistry {
	return WorkItemRegistry{FormatVersion: 3, TaskSpecPolicies: []TaskSpecPolicy{}, TaskProblemRevisions: []TaskProblemReference{}, TaskSpecRevisions: []TaskSpecReference{}, TaskSpecAssessments: []TaskAssessmentReference{}, TaskSolutionSelections: []TaskSelectionReference{}, TaskResultSpecBindings: []TaskResultSpecBinding{}, TaskContentPublications: []TaskContentPublication{}, Epics: []EpicRecord{}, Tasks: []TaskRecord{}, WorktreeOperations: []WorktreeOperationRecord{}, TaskResults: []TaskResultRecord{}, HumanQARecords: []TaskHumanQARecord{}, IntegrationAuthorities: []IntegrationAuthority{}, IntegrationIntents: []IntegrationIntent{}, IntegrationAttempts: []IntegrationAttempt{}, IntegrationResults: []IntegrationResult{}}
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

func publishWorkItemRegistryRecover(session WorkItemStoreSession, registry WorkItemRegistry) (WorkItemRegistry, error) {
	err := session.Publish(registry)
	if err == nil {
		return registry, nil
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
	return store.readRegistry(root, true)
}

// SnapshotRegistrations validates the registry itself without requiring the
// content of unrelated workflow/queue objects. Registered lists verify current
// problem heads individually so an unavailable title does not hide its Task.
// Snapshot and every mutation retain their existing full closure checks.
func (store *systemWorkItemStore) SnapshotRegistrations(root string) (WorkItemRegistry, error) {
	return store.readRegistry(root, false)
}

func (store *systemWorkItemStore) readRegistry(root string, requireClosure bool) (WorkItemRegistry, error) {
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
	if requireClosure && registry.FormatVersion == 4 {
		storage := &TaskContentStorage{}
		if err := validateTaskContentClosure(storage, root, registry, false); err != nil {
			return WorkItemRegistry{}, err
		}
		if err := validateQueueClosure(Dependencies{TaskContent: storage}, root, registry); err != nil {
			return WorkItemRegistry{}, err
		}
	}
	registry.RawSHA256 = digestTaskBytes(contents)
	return registry, nil
}

func (store *systemWorkItemStore) publish(root string, registry WorkItemRegistry) error {
	old, err := store.read(root)
	if err != nil {
		return err
	}
	if err := validateQueueTransition(old, registry); err != nil {
		return workError(ErrorWorkStoreConflict, err.Error(), err)
	}
	if err := validateQueueClosure(Dependencies{TaskContent: &TaskContentStorage{}}, root, registry); err != nil {
		return err
	}

	if err := validateTaskContentClosure(&TaskContentStorage{}, root, registry, false); err != nil {
		return err
	}
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
	contents := output.Bytes()
	if err := validateWorkItemRegistryYAMLShape(contents, registry.FormatVersion); err != nil {
		return nil, err
	}
	return contents, nil
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
	if registry.FormatVersion == 1 {
		if registry.TaskResults != nil || registry.HumanQARecords != nil || registry.IntegrationAuthorities != nil || registry.IntegrationIntents != nil || registry.IntegrationAttempts != nil || registry.IntegrationResults != nil {
			return WorkItemRegistry{}, errors.New("format 1 must not contain Task lifecycle collections")
		}
	} else if registry.FormatVersion >= 2 {
		if registry.TaskResults == nil || registry.HumanQARecords == nil || registry.IntegrationAuthorities == nil || registry.IntegrationIntents == nil || registry.IntegrationAttempts == nil || registry.IntegrationResults == nil {
			return WorkItemRegistry{}, errors.New("format 2 Task lifecycle collections must be explicit arrays")
		}
	}
	if err := validateWorkItemRegistryYAMLShape(contents, registry.FormatVersion); err != nil {
		return WorkItemRegistry{}, err
	}
	if err := validateWorkItemRegistry(registry); err != nil {
		return WorkItemRegistry{}, err
	}
	return registry, nil
}

type yamlShapeField struct {
	name   string
	typeOf reflect.Type
}

func validateWorkItemRegistryYAMLShape(contents []byte, version int) error {
	var document yaml.Node
	if err := yaml.Unmarshal(contents, &document); err != nil || len(document.Content) != 1 {
		if err == nil {
			err = errors.New("registry must contain one YAML document")
		}
		return err
	}
	fields := []yamlShapeField{
		{name: "format_version", typeOf: reflect.TypeOf(int(0))},
		{name: "epics", typeOf: reflect.TypeOf([]EpicRecord{})},
		{name: "tasks", typeOf: reflect.TypeOf([]TaskRecord{})},
		{name: "worktree_operations", typeOf: reflect.TypeOf([]WorktreeOperationRecord{})},
	}
	if version >= 2 {
		fields = append(fields,
			yamlShapeField{name: "task_results", typeOf: reflect.TypeOf([]TaskResultRecord{})},
			yamlShapeField{name: "human_qa_records", typeOf: reflect.TypeOf([]TaskHumanQARecord{})},
			yamlShapeField{name: "integration_authorities", typeOf: reflect.TypeOf([]IntegrationAuthority{})},
			yamlShapeField{name: "integration_intents", typeOf: reflect.TypeOf([]IntegrationIntent{})},
			yamlShapeField{name: "integration_attempts", typeOf: reflect.TypeOf([]IntegrationAttempt{})},
			yamlShapeField{name: "integration_results", typeOf: reflect.TypeOf([]IntegrationResult{})},
		)
	}
	if version >= 3 {
		fields = append(fields, yamlShapeField{name: "task_spec_policies", typeOf: reflect.TypeOf([]TaskSpecPolicy{})}, yamlShapeField{name: "task_problem_revisions", typeOf: reflect.TypeOf([]TaskProblemReference{})}, yamlShapeField{name: "task_spec_revisions", typeOf: reflect.TypeOf([]TaskSpecReference{})}, yamlShapeField{name: "task_spec_assessments", typeOf: reflect.TypeOf([]TaskAssessmentReference{})}, yamlShapeField{name: "task_solution_selections", typeOf: reflect.TypeOf([]TaskSelectionReference{})}, yamlShapeField{name: "task_result_spec_bindings", typeOf: reflect.TypeOf([]TaskResultSpecBinding{})}, yamlShapeField{name: "task_content_publications", typeOf: reflect.TypeOf([]TaskContentPublication{})})
	}
	if version == 4 {
		for _, f := range yamlStructShapeFields(reflect.TypeOf(WorkItemRegistry{})) {
			if f.name == "epic_base_versions" || f.name == "epic_base_updates" || f.name == "task_queue_events" || f.name == "task_preparations" {
				fields = append(fields, f)
			}
		}
	}
	return validateYAMLMappingShape(document.Content[0], fields, "work-item registry")
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
	switch expected.Kind() {
	case reflect.Map:
		return validateQueueYAMLMap(node)
	case reflect.Struct:
		fields := yamlStructShapeFields(expected)
		if expected == reflect.TypeOf(DeliveryAgreement{}) {
			present := map[string]bool{}
			for i := 0; i+1 < len(node.Content); i += 2 {
				present[node.Content[i].Value] = true
			}
			filtered := []yamlShapeField{}
			for _, field := range fields {
				optional := field.name == "source_ref" || field.name == "target_worktree" || field.name == "github_repository" || field.name == "remote"
				if !optional || present[field.name] {
					filtered = append(filtered, field)
				}
			}
			fields = filtered
		}
		if expected == reflect.TypeOf(IntegrationAuthority{}) {
			mode := ""
			for i := 0; i+1 < len(node.Content); i += 2 {
				if node.Content[i].Value == "mode" {
					mode = node.Content[i+1].Value
				}
			}
			if mode == "human_cli_start" {
				filtered := []yamlShapeField{}
				for _, f := range fields {
					if f.name != "delivery_owner" {
						filtered = append(filtered, f)
					}
				}
				fields = filtered
			}
		}
		if expected == reflect.TypeOf(WorkspaceTaskIntegrationPlan{}) {
			version := ""
			for i := 0; i+1 < len(node.Content); i += 2 {
				if node.Content[i].Value == "schema_version" {
					version = node.Content[i+1].Value
				}
			}
			if version != "3" {
				filtered := []yamlShapeField{}
				for _, field := range fields {
					if field.name != "delivery_authorization" {
						filtered = append(filtered, field)
					}
				}
				fields = filtered
			}
			if version == "1" {
				filtered := []yamlShapeField{}
				for _, f := range fields {
					if f.name != "task_spec_guard" {
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

func validateWorkItemRegistry(registry WorkItemRegistry) error {
	if err := validateQueueRegistry(registry); err != nil {
		return err
	}
	if err := validateTaskContentRegistry(registry); err != nil {
		return err
	}
	if registry.FormatVersion != 1 && registry.FormatVersion != 2 && registry.FormatVersion != 3 && registry.FormatVersion != 4 {
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
		if hasOperation && (operation.ProjectID != task.ProjectID || operation.RepoID != task.RepoID || operation.GitCommonDir != task.GitCommonDir || operation.ParentEpicID != task.ParentEpicID || operation.ParentWorktreeID != binding.Worktree.ID || operation.ParentRef != binding.Worktree.Ref || !historicalEpicBaseMatches(registry, task.ParentEpicID, task.RepoID, operation.ParentWorktreeID, operation.ParentRef, operation.ParentOID, operation.ParentTree)) {
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
	if registry.FormatVersion >= 2 {
		if err := validateTaskLifecycleRegistry(registry); err != nil {
			return err
		}
	}
	return nil
}

func validateTaskLifecycleRegistry(registry WorkItemRegistry) error {
	if !sortedTaskResults(registry.TaskResults) || !sortedHumanQA(registry.HumanQARecords) || !sortedAuthorities(registry.IntegrationAuthorities) || !sortedIntents(registry.IntegrationIntents) || !sortedAttempts(registry.IntegrationAttempts) || !sortedIntegrationResults(registry.IntegrationResults) {
		return errors.New("Task lifecycle arrays must be sorted by unique IDs")
	}
	tasks := map[TaskID]TaskRecord{}
	for _, task := range registry.Tasks {
		tasks[task.ID] = task
	}
	results := map[TaskResultID]TaskResultRecord{}
	publication := map[string]bool{}
	triples := map[string]bool{}
	for _, r := range registry.TaskResults {
		task, ok := tasks[r.TaskID]
		if !ok || task.Worktree == nil || task.Worktree.ID != r.TaskWorktreeID || task.ProjectID != r.ProjectID || task.RepoID != r.RepoID || task.GitCommonDir != r.GitCommonDir || task.Worktree.Locator != r.SourceLocator || task.Worktree.Ref != r.SourceRef {
			return fmt.Errorf("Task result %s has invalid Task binding", r.ID)
		}
		if !lifecycleIDPatterns["task result"].MatchString(string(r.ID)) || publication[r.PublicationKey] || validateTaskResultRecordShape(r) != nil {
			return fmt.Errorf("Task result %s has invalid identity or digest", r.ID)
		}
		if r.StoreTransition != "none" && r.StoreTransition != "format_1_to_2" {
			return fmt.Errorf("Task result %s has invalid store transition", r.ID)
		}
		if !setString("passed", "good_enough_with_known_debt", "failed", "unknown")[r.TechnicalGate] {
			return fmt.Errorf("Task result %s has invalid technical gate", r.ID)
		}
		key := string(r.TaskID) + "\x00" + r.TerminalResultID + "\x00" + r.ResultOID
		if triples[key] {
			return fmt.Errorf("duplicate Task result evidence")
		}
		publication[r.PublicationKey] = true
		triples[key] = true
		results[r.ID] = r
	}
	qaIDs := map[HumanQARecordID]TaskHumanQARecord{}
	qaPub := map[string]bool{}
	for _, q := range registry.HumanQARecords {
		r, ok := results[q.TaskResultID]
		if !ok || r.TaskID != q.TaskID || r.ResultOID != q.ResultOID || r.ResultTree != q.ResultTree {
			return fmt.Errorf("human QA %s has invalid Task result binding", q.ID)
		}
		if !lifecycleIDPatterns["human QA"].MatchString(string(q.ID)) || qaPub[q.PublicationKey] || validateTaskHumanQARecordShape(q) != nil {
			return fmt.Errorf("human QA %s has invalid identity", q.ID)
		}
		qaPub[q.PublicationKey] = true
		qaIDs[q.ID] = q
	}
	authorities := map[IntegrationAuthorityID]IntegrationAuthority{}
	planDigests := map[string]bool{}
	for _, a := range registry.IntegrationAuthorities {
		validMode := a.Mode == "human_cli_start" && a.DeliveryOwner == nil || a.Mode == "delivery_owner_after_human_pass" && validDeliveryIntegrationOwner(a.DeliveryOwner)
		if a.Plan.DeliveryAuthorization != nil && a.DeliveryOwner != nil && (a.DeliveryOwner.RunID != a.Plan.DeliveryAuthorization.RunID || a.DeliveryOwner.RequestSHA256 != a.Plan.DeliveryAuthorization.RequestSHA256) {
			validMode = false
		}
		if !lifecycleIDPatterns["integration authority"].MatchString(string(a.ID)) || !validMode || !validTaskUTC(a.CreatedAtUTC) || !digestPattern.MatchString(a.PlanSHA256) || integrationPlanDigest(a.Plan) != a.PlanSHA256 || planDigests[a.PlanSHA256] {
			return fmt.Errorf("integration authority %s is invalid", a.ID)
		}
		r, rok := results[a.TaskResultID]
		q, qok := qaIDs[a.HumanQARecordID]
		if !rok || !qok || r.TaskID != a.TaskID || q.TaskResultID != r.ID || q.Outcome != "pass" || !(r.TechnicalGate == "passed" || r.TechnicalGate == "good_enough_with_known_debt") {
			return fmt.Errorf("integration authority %s has invalid gate binding", a.ID)
		}
		if err := validateStoredIntegrationPlan(a.Plan, registry, r, q); err != nil || a.Plan.Task.TaskID != a.TaskID || a.Plan.TaskResult.ID != a.TaskResultID || a.Plan.HumanQA.ID != a.HumanQARecordID || !sameOptionalIntegrationResultID(a.RetryAfterResultID, a.Plan.RetryAfterResultID) || !sameAllowedEffect(a.AllowedEffect, a.Plan.Effect) {
			return fmt.Errorf("integration authority %s has invalid plan binding", a.ID)
		}
		planDigests[a.PlanSHA256] = true
		authorities[a.ID] = a
	}
	intents := map[IntegrationIntentID]IntegrationIntent{}
	intentByAuthority := map[IntegrationAuthorityID]bool{}
	for _, i := range registry.IntegrationIntents {
		a, ok := authorities[i.AuthorityID]
		if !ok || intentByAuthority[i.AuthorityID] || !lifecycleIDPatterns["integration intent"].MatchString(string(i.ID)) || i.State != "prepared" || i.MaxAttempts != 1 || integrationIntentDigest(i) != i.IntentSHA256 || i.PlanSHA256 != a.PlanSHA256 || !intentMatchesAuthority(i, a) {
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
			if !ok || a.ResultID == nil || *a.ResultID != r.ID || a.State != "result_recorded" || !reflect.DeepEqual(r.BeforeObservation, a.PreObservation) {
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
		if !ok || prior.Outcome != "no_effect" || retryUse[*authority.RetryAfterResultID] {
			return fmt.Errorf("integration authority %s has invalid retry predecessor", authority.ID)
		}
		priorAuthority := authorities[prior.AuthorityID]
		if priorAuthority.TaskID != authority.TaskID || priorAuthority.TaskResultID != authority.TaskResultID || priorAuthority.HumanQARecordID != authority.HumanQARecordID || priorAuthority.Plan.Epic.ParentRef != authority.Plan.Epic.ParentRef {
			return fmt.Errorf("integration authority %s retry binding differs", authority.ID)
		}
		retryUse[*authority.RetryAfterResultID] = true
	}
	return nil
}

func validateTaskResultRecordShape(r TaskResultRecord) error {
	if validatePublicationKey(r.PublicationKey) != nil || !digestPattern.MatchString(r.DraftSHA256) || !handoffEvidenceIDPatterns["activity"].MatchString(r.ActivityID) || !handoffEvidenceIDPatterns["run"].MatchString(r.RunID) || !handoffEvidenceIDPatterns["handoff"].MatchString(r.HandoffID) || !handoffEvidenceIDPatterns["start"].MatchString(r.StartReceiptID) || !handoffEvidenceIDPatterns["terminal"].MatchString(r.TerminalResultID) {
		return errors.New("invalid result identity")
	}
	if !validAbsoluteCleanPath(r.HandoffLocator) || !validAbsoluteCleanPath(r.StartReceiptLocator) || !validAbsoluteCleanPath(r.TerminalResultLocator) || !digestPattern.MatchString(r.HandoffSHA256) || !digestPattern.MatchString(r.StartReceiptSHA256) || !digestPattern.MatchString(r.TerminalResultSHA256) || !digestPattern.MatchString(r.InspectionSHA256) || !validOIDText(r.ResultOID) || !validOIDText(r.ResultTree) || len(r.ResultOID) != len(r.ResultTree) {
		return errors.New("invalid result evidence")
	}
	if !setString("complete", "blocked", "budget_exhausted", "unknown", "conflict")[r.ReportedOutcome] || !setString("passed", "good_enough_with_known_debt", "failed", "unknown")[r.TechnicalGate] || !validTaskText(r.Recorder.ActorClaim, 1, 256) || !validTaskText(r.Recorder.ControlSurface, 1, 256) || !validTaskUTC(r.Recorder.RecordedAtUTC) {
		return errors.New("invalid result assessment")
	}
	last := ""
	artifacts := map[string]bool{}
	var total int64
	for _, a := range r.Artifacts {
		total += a.SizeBytes
		if a.ArtifactID <= last || validatePublicationKey(a.ArtifactID) != nil || !setString("verifier_stdout", "verifier_stderr", "review", "debt_control", "other")[a.Role] || !validAbsoluteCleanPath(a.Locator) || !digestPattern.MatchString(a.SHA256) || a.SizeBytes < 0 || a.SizeBytes > 64<<20 || total > 256<<20 {
			return errors.New("invalid result artifact")
		}
		last, artifacts[a.ArtifactID] = a.ArtifactID, true
	}
	last = ""
	for _, v := range r.VerifierResults {
		if v.VerifierID <= last || validatePublicationKey(v.VerifierID) != nil || !setString("passed", "failed")[v.Outcome] || len(v.Argv) == 0 || !validAbsoluteCleanPath(v.CWD) || (!validOIDText(v.BoundOIDOrSHA256) && !digestPattern.MatchString(v.BoundOIDOrSHA256)) || (v.StdoutArtifactID != nil && !artifacts[*v.StdoutArtifactID]) || (v.StderrArtifactID != nil && !artifacts[*v.StderrArtifactID]) {
			return errors.New("invalid verifier result")
		}
		for _, arg := range v.Argv {
			if !validTaskText(arg, 1, 4096) {
				return errors.New("invalid verifier argv")
			}
		}
		last = v.VerifierID
	}
	if err := validateTaskReview(r.Review, artifacts); err != nil {
		return err
	}
	last = ""
	for _, d := range r.AcceptedDebt {
		if d.ID <= last || validatePublicationKey(d.ID) != nil || !setString("coverage", "evidence", "journal")[d.RiskClass] || !setString("low", "medium")[d.Severity] || !validTaskText(d.Summary, 1, 600) || !validTaskText(d.Control, 1, 2000) || !sortedNonemptyArtifactIDs(d.EvidenceArtifactIDs, artifacts) {
			return errors.New("invalid accepted debt")
		}
		last = d.ID
	}
	return nil
}

func validateTaskReview(review TaskReviewRecord, artifacts map[string]bool) error {
	for _, entries := range [][]TaskReviewEntry{review.Findings, review.Fixes, review.OpenActionableFindings} {
		last := ""
		for _, entry := range entries {
			if entry.ID <= last || validatePublicationKey(entry.ID) != nil || !setString("low", "medium", "high", "critical")[entry.Severity] || !validTaskText(entry.Summary, 1, 600) || !sortedArtifactIDs(entry.EvidenceArtifactIDs, artifacts) {
				return errors.New("invalid review entry")
			}
			last = entry.ID
		}
	}
	return nil
}

func sortedArtifactIDs(ids []string, known map[string]bool) bool {
	for i, id := range ids {
		if validatePublicationKey(id) != nil || !known[id] || (i > 0 && ids[i-1] >= id) {
			return false
		}
	}
	return ids != nil
}
func sortedNonemptyArtifactIDs(ids []string, known map[string]bool) bool {
	return len(ids) > 0 && sortedArtifactIDs(ids, known)
}

func validateTaskHumanQARecordShape(q TaskHumanQARecord) error {
	if validatePublicationKey(q.PublicationKey) != nil || !digestPattern.MatchString(q.DraftSHA256) || !validOIDText(q.ResultOID) || !validOIDText(q.ResultTree) || len(q.ResultOID) != len(q.ResultTree) || !setString("pass", "fail", "blocked")[q.Outcome] || !validTaskText(q.Observation, 1, 2000) {
		return errors.New("invalid QA record")
	}
	started, se := time.Parse(time.RFC3339Nano, q.Actor.StartedAtUTC)
	completed, ce := time.Parse(time.RFC3339Nano, q.Actor.CompletedAtUTC)
	if !validTaskText(q.Actor.ActorClaim, 1, 256) || !validTaskText(q.Actor.StartSurface, 1, 256) || se != nil || ce != nil || !strings.HasSuffix(q.Actor.StartedAtUTC, "Z") || !strings.HasSuffix(q.Actor.CompletedAtUTC, "Z") || completed.Before(started) {
		return errors.New("invalid QA actor")
	}
	if len(q.Evidence) == 0 {
		return errors.New("QA evidence is empty")
	}
	last := ""
	var total int64
	report := false
	for _, e := range q.Evidence {
		total += e.SizeBytes
		if e.ID <= last || validatePublicationKey(e.ID) != nil || !setString("script", "report", "screenshot", "log", "other")[e.Role] || !validAbsoluteCleanPath(e.Locator) || !digestPattern.MatchString(e.SHA256) || e.SizeBytes < 0 || e.SizeBytes > 64<<20 || total > 256<<20 {
			return errors.New("invalid QA evidence")
		}
		last = e.ID
		report = report || e.Role == "report"
	}
	if q.Outcome == "pass" && !report {
		return errors.New("QA pass has no report")
	}
	last = ""
	for _, risk := range q.AcceptedResidualRisks {
		if risk.ID <= last || validatePublicationKey(risk.ID) != nil || !setString("low", "medium")[risk.Severity] || !validTaskText(risk.Summary, 1, 600) {
			return errors.New("invalid QA residual risk")
		}
		last = risk.ID
	}
	return nil
}

func validateStoredIntegrationPlan(p WorkspaceTaskIntegrationPlan, registry WorkItemRegistry, result TaskResultRecord, qa TaskHumanQARecord) error {
	if validateStoredTaskSpecGuard(registry, p) != nil || p.Format != "json" || p.FormatVersion != 1 || p.Canonicalization != "RFC8785" || !validAbsoluteCleanPath(p.Workspace.Root) || !digestPattern.MatchString(p.Workspace.MarkerSHA256) || p.Project.ProjectID != result.ProjectID || p.Repository.RepoID != result.RepoID || !validAbsoluteCleanPath(p.Repository.RegisteredLocator) || p.Repository.GitCommonDir != result.GitCommonDir || !validAbsoluteCleanPath(p.Repository.GitCommonDir) || !setString("sha1", "sha256")[p.Repository.ObjectFormat] || p.Repository.RefFormat != "files" {
		return errors.New("invalid integration plan envelope")
	}
	task, _ := findTask(registry, result.TaskID)
	if task == nil || task.Worktree == nil {
		return errors.New("integration plan Task binding is missing")
	}
	epic, _ := findEpic(registry, task.ParentEpicID)
	if epic == nil {
		return errors.New("integration plan Epic binding is missing")
	}
	parent, _ := findEpicRepo(*epic, result.RepoID)
	if parent == nil || p.Task.TaskID != task.ID || p.Task.TaskWorktreeID != task.Worktree.ID || p.Task.SourceLocator != result.SourceLocator || p.Task.SourceRef != result.SourceRef || p.Task.ResultOID != result.ResultOID || p.Task.ResultTree != result.ResultTree || p.Epic.EpicID != epic.ID || p.Epic.ParentWorktreeID != parent.Worktree.ID || p.Epic.ParentLocator != parent.Worktree.Locator || p.Epic.ParentRef != parent.Worktree.Ref || p.Epic.ExpectedParentOID != task.Worktree.ParentOID || p.Epic.ExpectedParentTree != task.Worktree.ParentTree || p.TaskResult.ID != result.ID || p.TaskResult.DraftSHA256 != result.DraftSHA256 || p.TaskResult.TechnicalGate != result.TechnicalGate || p.HumanQA.ID != qa.ID || p.HumanQA.DraftSHA256 != qa.DraftSHA256 || p.HumanQA.Outcome != qa.Outcome || p.StoreTransition != result.StoreTransition || !digestPattern.MatchString(p.WorkItemsSHA256) {
		return errors.New("integration plan binding differs")
	}
	if p.DeliveryAuthorization != nil {
		auth := p.DeliveryAuthorization
		target := QueueTarget{ProjectID: p.Project.ProjectID, RepoID: p.Repository.RepoID, EpicID: p.Epic.EpicID, GitCommonDir: p.Repository.GitCommonDir, ParentWorktreeID: p.Epic.ParentWorktreeID, ParentLocator: p.Epic.ParentLocator, ParentRef: p.Epic.ParentRef}
		bound, err := BindDeliveryAgreement(auth.Agreement, target, p.Task.SourceRef)
		if err != nil || ValidateDeliveryAuthorization(*auth) != nil || auth.CandidateRunID != result.RunID || auth.Agreement.Mode == DeliveryPullRequest || !contentTypedEqual(bound, auth.Agreement) {
			return errors.New("integration plan delivery authority differs from exact result and target")
		}
	}
	if validatePlanObservation(p.ObservedSource) != nil || validatePlanObservation(p.ObservedParent) != nil || validateInventory(p.ObservedInventory) != nil || validateReflog(p.ObservedReflog) != nil || !setString("ready", "already_integrated", "blocked", "conflict", "unknown")[p.Readiness] || !sortedReasonTokens(p.Reasons) {
		return errors.New("invalid integration plan observation")
	}
	source := worktreeFromPlan(p.ObservedSource)
	parentObservation := worktreeFromPlan(p.ObservedParent)
	mainForbidden := ordinaryProductRef(p.Epic.ParentRef) && (p.DeliveryAuthorization == nil || p.DeliveryAuthorization.Agreement.Mode != DeliveryLocalBranch)
	if !setString("ready", "already_integrated")[p.Readiness] || len(p.Reasons) != 0 || !p.TechnicalGateReady || !p.HumanQAReady || !p.AncestryReady || p.Repository.Shallow || p.Repository.PartialClone || p.Repository.SparseCheckout || mainForbidden || p.ObservedSource.OID != p.Task.ResultOID || p.ObservedSource.Tree != p.Task.ResultTree || !p.ObservedSource.Clean || len(p.ObservedSource.StatusEntries) != 0 || len(p.ObservedSource.InProgress) != 0 || !p.ObservedParent.Clean || len(p.ObservedParent.StatusEntries) != 0 || len(p.ObservedParent.InProgress) != 0 || (p.Readiness == "ready" && (p.ObservedParent.OID != p.Epic.ExpectedParentOID || p.ObservedParent.Tree != p.Epic.ExpectedParentTree)) || (p.Readiness == "already_integrated" && (p.ObservedParent.OID != p.Task.ResultOID || p.ObservedParent.Tree != p.Task.ResultTree)) || p.ObservedSource.GitCommonDir != p.Repository.GitCommonDir || p.ObservedParent.GitCommonDir != p.Repository.GitCommonDir || p.ObservedSource.ObjectFormat != p.Repository.ObjectFormat || p.ObservedParent.ObjectFormat != p.Repository.ObjectFormat || p.ObservedSource.RefFormat != p.Repository.RefFormat || p.ObservedParent.RefFormat != p.Repository.RefFormat || !integrationInventoryContains(p.ObservedInventory, source) || !integrationInventoryContains(p.ObservedInventory, parentObservation) || !objectIDsMatchFormat(p.Repository.ObjectFormat, p.Epic.ExpectedParentOID, p.Epic.ExpectedParentTree, p.Task.ResultOID, p.Task.ResultTree, p.ObservedSource.OID, p.ObservedSource.Tree, p.ObservedParent.OID, p.ObservedParent.Tree) {
		return errors.New("integration authority plan is not a ready exact plan")
	}
	wantArgv := []string{"git", "-c", "core.hooksPath=" + os.DevNull, "-c", "merge.autoStash=false", "-c", "gc.auto=0", "-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false", "-c", "submodule.recurse=false", "-C", p.Epic.ParentLocator, "merge", "--ff-only", "--no-stat", "--no-autostash", p.Task.ResultOID}
	if p.TechnicalGateReady != (result.TechnicalGate == "passed" || result.TechnicalGate == "good_enough_with_known_debt") || p.HumanQAReady != (qa.Outcome == "pass") || p.Effect.Kind != "local_ff_only" || p.Effect.ParentRef != p.Epic.ParentRef || p.Effect.ExpectedParentOID != p.Epic.ExpectedParentOID || p.Effect.ResultOID != p.Task.ResultOID || p.Effect.MaxOccurrences != 1 || !slices.Equal(p.Effect.Argv, wantArgv) {
		return errors.New("invalid integration plan effect")
	}
	return nil
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
	return i.WorkspaceRoot == p.Workspace.Root && i.WorkspaceMarkerSHA256 == p.Workspace.MarkerSHA256 && i.ProjectID == p.Project.ProjectID && i.RepoID == p.Repository.RepoID && i.RegisteredRepoLocator == p.Repository.RegisteredLocator && i.GitCommonDir == p.Repository.GitCommonDir && i.ObjectFormat == p.Repository.ObjectFormat && i.RefFormat == p.Repository.RefFormat && i.EpicID == p.Epic.EpicID && i.ParentWorktreeID == p.Epic.ParentWorktreeID && i.ParentLocator == p.Epic.ParentLocator && i.ParentRef == p.Epic.ParentRef && i.ExpectedParentOID == p.Epic.ExpectedParentOID && i.ExpectedParentTree == p.Epic.ExpectedParentTree && i.TaskID == p.Task.TaskID && i.TaskWorktreeID == p.Task.TaskWorktreeID && i.SourceLocator == p.Task.SourceLocator && i.SourceRef == p.Task.SourceRef && i.ResultOID == p.Task.ResultOID && i.ResultTree == p.Task.ResultTree && i.TaskResultID == p.TaskResult.ID && i.HumanQARecordID == p.HumanQA.ID && sameAllowedEffect(i.Effect, p.Effect)
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
	emptyArgv := len(result.NextAction.Argv) == 0
	switch result.Outcome {
	case "exact_effect", "already_integrated":
		if result.RecoveryStatus != "complete" || result.NextAction.Kind != "none" || !emptyArgv {
			return errors.New("successful integration result has contradictory recovery guidance")
		}
	case "no_effect":
		want := []string{"ply", "workspace", "task", "integrate", string(authority.TaskID), "--result", string(authority.TaskResultID), "--qa", string(authority.HumanQARecordID), "--expected-result-oid", authority.Plan.Task.ResultOID, "--expected-parent-oid", authority.Plan.Epic.ExpectedParentOID, "--retry-after", string(result.ID), "--check"}
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

func sortedTaskResults(v []TaskResultRecord) bool {
	for i := 1; i < len(v); i++ {
		if v[i-1].ID >= v[i].ID {
			return false
		}
	}
	return true
}
func sortedHumanQA(v []TaskHumanQARecord) bool {
	for i := 1; i < len(v); i++ {
		if v[i-1].ID >= v[i].ID {
			return false
		}
	}
	return true
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

type workItemRegistryV4Wire struct {
	FormatVersion           int                       `yaml:"format_version"`
	Epics                   []EpicRecord              `yaml:"epics"`
	Tasks                   []TaskRecord              `yaml:"tasks"`
	WorktreeOperations      []WorktreeOperationRecord `yaml:"worktree_operations"`
	TaskResults             []TaskResultRecord        `yaml:"task_results"`
	HumanQARecords          []TaskHumanQARecord       `yaml:"human_qa_records"`
	IntegrationAuthorities  []IntegrationAuthority    `yaml:"integration_authorities"`
	IntegrationIntents      []IntegrationIntent       `yaml:"integration_intents"`
	IntegrationAttempts     []IntegrationAttempt      `yaml:"integration_attempts"`
	IntegrationResults      []IntegrationResult       `yaml:"integration_results"`
	TaskSpecPolicies        []TaskSpecPolicy          `yaml:"task_spec_policies"`
	TaskProblemRevisions    []TaskProblemReference    `yaml:"task_problem_revisions"`
	TaskSpecRevisions       []TaskSpecReference       `yaml:"task_spec_revisions"`
	TaskSpecAssessments     []TaskAssessmentReference `yaml:"task_spec_assessments"`
	TaskSolutionSelections  []TaskSelectionReference  `yaml:"task_solution_selections"`
	TaskResultSpecBindings  []TaskResultSpecBinding   `yaml:"task_result_spec_bindings"`
	TaskContentPublications []TaskContentPublication  `yaml:"task_content_publications"`
	EpicBaseVersions        []EpicBaseVersion         `yaml:"epic_base_versions"`
	EpicBaseUpdates         []EpicBaseUpdate          `yaml:"epic_base_updates"`
	TaskQueueEvents         []TaskQueueEvent          `yaml:"task_queue_events"`
	TaskPreparations        []TaskPreparation         `yaml:"task_preparations"`
}

func workItemRegistryV4(r WorkItemRegistry) workItemRegistryV4Wire {
	return workItemRegistryV4Wire{FormatVersion: r.FormatVersion, Epics: r.Epics, Tasks: r.Tasks, WorktreeOperations: r.WorktreeOperations, TaskResults: r.TaskResults, HumanQARecords: r.HumanQARecords, IntegrationAuthorities: r.IntegrationAuthorities, IntegrationIntents: r.IntegrationIntents, IntegrationAttempts: r.IntegrationAttempts, IntegrationResults: r.IntegrationResults, TaskSpecPolicies: r.TaskSpecPolicies, TaskProblemRevisions: r.TaskProblemRevisions, TaskSpecRevisions: r.TaskSpecRevisions, TaskSpecAssessments: r.TaskSpecAssessments, TaskSolutionSelections: r.TaskSolutionSelections, TaskResultSpecBindings: r.TaskResultSpecBindings, TaskContentPublications: r.TaskContentPublications, EpicBaseVersions: r.EpicBaseVersions, EpicBaseUpdates: r.EpicBaseUpdates, TaskQueueEvents: r.TaskQueueEvents, TaskPreparations: r.TaskPreparations}
}
