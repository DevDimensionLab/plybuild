// Package weekrun supervises an explicitly selected set of native Goal executions.
// Its files are operational control/observation records, never a second Task store.
package weekrun

import (
	"fmt"
	"regexp"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const MaxDurationSeconds = 7 * 24 * 60 * 60

// Contract is immutable after first start. Dependencies name selected Task IDs.
// The Goal revision/hash, native executor and delivery agreement are rechecked
// before every launch; later queue additions never enlarge this authorization.
type Contract struct {
	Kind                string              `json:"kind"`
	SchemaVersion       int                 `json:"schema_version"`
	Name                string              `json:"name"`
	Workspace           string              `json:"workspace"`
	ReturnWorktree      string              `json:"return_worktree"`
	ReturnRef           string              `json:"return_ref"`
	ProjectID           workspace.ProjectID `json:"project_id"`
	RepoID              workspace.RepoID    `json:"repo_id"`
	EpicID              workspace.EpicID    `json:"epic_id"`
	Selected            []Selection         `json:"selected"`
	Model               string              `json:"model"`
	Effort              string              `json:"effort"`
	DurationSeconds     int                 `json:"duration_seconds"`
	MaxAttempts         int                 `json:"max_attempts"`
	AgentTimeoutSeconds int                 `json:"agent_timeout_seconds"`
	MaxProviderRetries  int                 `json:"max_provider_retries"`
}

type Selection struct {
	TaskID    workspace.TaskID      `json:"task_id"`
	Goal      workspace.TaskGoalRef `json:"goal"`
	DependsOn []workspace.TaskID    `json:"depends_on"`
}

type Binding struct {
	Contract    Contract           `json:"contract"`
	SHA256      string             `json:"contract_sha256"`
	ID          string             `json:"id"`
	StartedAt   string             `json:"started_at"`
	Deadline    string             `json:"deadline"`
	Control     taskrun.Executable `json:"control"`
	Environment map[string]string  `json:"environment"`
}

type Entry struct {
	TaskID       workspace.TaskID `json:"task_id"`
	State        string           `json:"state"`
	RunID        string           `json:"run_id,omitempty"`
	Reason       string           `json:"reason,omitempty"`
	BaseOID      string           `json:"base_oid,omitempty"`
	BaseRevision int              `json:"base_revision,omitempty"`
}

type State struct {
	Kind            string           `json:"kind"`
	SchemaVersion   int              `json:"schema_version"`
	ID              string           `json:"id"`
	ContractSHA256  string           `json:"contract_sha256"`
	Phase           string           `json:"phase"`
	StartedAt       string           `json:"started_at"`
	Deadline        string           `json:"deadline"`
	LastProgressAt  string           `json:"last_progress_at"`
	UpdatedAt       string           `json:"updated_at"`
	NextCheckAt     string           `json:"next_check_at,omitempty"`
	NextAction      string           `json:"next_action"`
	Reason          string           `json:"reason,omitempty"`
	Entries         []Entry          `json:"entries"`
	ActiveTaskID    workspace.TaskID `json:"active_task_id,omitempty"`
	DesiredAction   string           `json:"desired_action"`
	SeenControl     string           `json:"seen_control,omitempty"`
	ResumeRequested bool             `json:"resume_requested,omitempty"`
	Generation      int              `json:"generation"`
	ControllerPID   int              `json:"controller_pid"`
	OwnerState      string           `json:"owner_state,omitempty"`
}

type Control struct {
	ID          string `json:"id"`
	Action      string `json:"action"`
	RequestedAt string `json:"requested_at"`
}

type Run struct {
	ID                  string
	Phase               string
	Stage               string
	OwnerState          string
	SupervisorState     string
	ControlAction       string
	Reason              string
	RetryNotBefore      string
	LastProviderFailure string
	BaseOID             string
	BaseRevision        int
	Completed           bool
	// Releasable means a concrete completed work outcome, never an unknown
	// process/provider/authentication failure. Native.Release proves it again.
	Releasable bool
}

// Native only reads or invokes the established native APIs. Implementations
// must preserve a Task's native Run identity after an interrupted Start.
type Native interface {
	Queue(Contract) (workspace.WorkspaceTaskQueueReadback, error)
	Preview(Contract, Selection, taskrun.Executable) error
	Start(Contract, Selection, taskrun.Executable, bool) (Run, error)
	Observe(Contract, string) (Run, error)
	Resume(Contract, string, bool) (Run, error)
	Control(Contract, string, string) (Run, error)
	Release(Contract, string, workspace.TaskID) error
}

type Manager struct {
	Native     Native
	Now        func() time.Time
	Executable func() (string, error)
	Launch     func(Binding) error
	Launchctl  Launchctl
}

var slug = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,95}$`)
var digest = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

func (c Contract) Target() workspace.QueueTargetInput {
	return workspace.QueueTargetInput{ProjectID: c.ProjectID, RepoID: c.RepoID, EpicID: c.EpicID}
}

func Validate(c Contract) error {
	if c.Kind != "ply.workflow.week-contract" || c.SchemaVersion != 1 || !slug.MatchString(c.Name) {
		return fmt.Errorf("week requires ply.workflow.week-contract schema 1 and a short lowercase name")
	}
	for _, p := range []string{c.Workspace, c.ReturnWorktree} {
		if err := physical(p); err != nil {
			return err
		}
	}
	if !slug.MatchString(string(c.ProjectID)) || !slug.MatchString(string(c.RepoID)) || !slug.MatchString(string(c.EpicID)) || c.ReturnRef == "" {
		return fmt.Errorf("week requires an explicit native Project/Repo/Epic and return ref")
	}
	if c.Model == "" || c.Effort == "" || c.DurationSeconds < 1 || c.DurationSeconds > MaxDurationSeconds || c.MaxAttempts < 1 || c.MaxAttempts > 20 || c.AgentTimeoutSeconds < 10 || c.AgentTimeoutSeconds > 86400 || c.MaxProviderRetries < 1 || c.MaxProviderRetries > 1000 {
		return fmt.Errorf("week requires an explicit model/effort, duration 1..604800, attempts 1..20, provider timeout 10..86400 and provider retries 1..1000")
	}
	if len(c.Selected) > 256 {
		return fmt.Errorf("week selects at most 256 native Goals")
	}
	ids := map[workspace.TaskID]Selection{}
	specs := map[string]bool{}
	for _, s := range c.Selected {
		if !slug.MatchString(string(s.TaskID)) || !slug.MatchString(s.Goal.SpecID) || s.Goal.Spec.Revision < 1 || !digest.MatchString(s.Goal.Spec.ManifestSHA256) {
			return fmt.Errorf("invalid exact Task/Goal binding")
		}
		if _, ok := ids[s.TaskID]; ok || specs[s.Goal.SpecID] {
			return fmt.Errorf("duplicate Task or ambiguous Spec in selected work")
		}
		ids[s.TaskID], specs[s.Goal.SpecID] = s, true
	}
	seen, active := map[workspace.TaskID]bool{}, map[workspace.TaskID]bool{}
	var visit func(workspace.TaskID) error
	visit = func(id workspace.TaskID) error {
		if active[id] {
			return fmt.Errorf("week dependency cycle at %s", id)
		}
		if seen[id] {
			return nil
		}
		s, ok := ids[id]
		if !ok {
			return fmt.Errorf("dependency %s is outside the selected work", id)
		}
		active[id] = true
		deps := map[workspace.TaskID]bool{}
		for _, dep := range s.DependsOn {
			if deps[dep] {
				return fmt.Errorf("duplicate dependency %s", dep)
			}
			deps[dep] = true
			if err := visit(dep); err != nil {
				return err
			}
		}
		delete(active, id)
		seen[id] = true
		return nil
	}
	for id := range ids {
		if err := visit(id); err != nil {
			return err
		}
	}
	return nil
}

func ReadContract(path string) (Contract, error) {
	var c Contract
	err := readJSON(path, &c)
	if err == nil {
		err = Validate(c)
	}
	return c, err
}
