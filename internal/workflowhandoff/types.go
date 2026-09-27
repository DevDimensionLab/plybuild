package workflowhandoff

import (
	"io"
	"io/fs"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type HandoffID string

type CreateInput struct{ DraftPath string }

type InspectInput struct {
	HandoffLocator    string
	Format            string
	Raw               string
	AcknowledgeSecret bool
}

type SubmitInput struct{ HandoffLocator, DraftPath string }

type ControlInput struct {
	HandoffID                 HandoffID
	Reason                    string
	AcknowledgeEffectsUnknown bool
}

type SupersedeInput struct {
	HandoffID HandoffID
	DraftPath string
	Reason    string
}

type CreateResult struct {
	HandoffID HandoffID
	Purpose   string
	Worktree  string
	Locator   string
	Created   bool
}

type ShowResult struct{ Status, Result, Meaning, NextAction string }

type InspectResult struct {
	Bytes    []byte
	AppendLF bool
}

type SubmitResult struct {
	Phase      string
	DocumentID string
	Locator    string
	SHA256     string
	Created    bool
}

type ControlResult struct {
	HandoffID            HandoffID
	ReplacementHandoffID HandoffID
	Reason               string
	Purpose              string
	Worktree             string
	Locator              string
}

type File interface {
	Write([]byte) (int, error)
	Chmod(fs.FileMode) error
	Sync() error
	Close() error
	Stat() (fs.FileInfo, error)
	Name() string
}

type FileSystem interface {
	Getwd() (string, error)
	EvalSymlinks(string) (string, error)
	Lstat(string) (fs.FileInfo, error)
	Stat(string) (fs.FileInfo, error)
	ReadFile(string) ([]byte, error)
	ReadDir(string) ([]fs.DirEntry, error)
	Mkdir(string, fs.FileMode) error
	MkdirAll(string, fs.FileMode) error
	MkdirTemp(string, string) (string, error)
	OpenFile(string, int, fs.FileMode) (File, error)
	Chmod(string, fs.FileMode) error
	Rename(string, string) error
	Remove(string) error
}

type WorkspaceSnapshot struct {
	Observation  workspace.WorkspaceObservation
	Projects     []workspace.ProjectRecord
	Repositories []workspace.RepoRecord
}

type WorkspaceObserver interface {
	ObserveContaining() (WorkspaceSnapshot, error)
	ObserveRoot(string) (WorkspaceSnapshot, error)
}

type StatusEntry struct {
	Path, IndexState, WorktreeState, Kind, SHA256 string
}

type StatusPolicy struct {
	Mode    string
	Entries []StatusEntry
	Value   canonicaljson.Object
}

type TargetRequest struct{ Worktree, Ref string }

type TargetObservation struct {
	Worktree, GitCommonDir, Ref, OID, Tree string
	Status                                 StatusPolicy
}

type GitBinding struct{ GitCommonDir, Ref, OID, Blob string }

type InputGitRequest struct {
	Locator string
	Binding GitBinding
}

type InputGitObservation struct {
	Binding         GitBinding
	MatchesExpected bool
}

type GitObserver interface {
	ObserveTarget(TargetRequest) (TargetObservation, error)
	ObserveInputGit(InputGitRequest) (InputGitObservation, error)
}

type Clock interface{ Now() time.Time }

type Dependencies struct {
	Files     FileSystem
	Workspace WorkspaceObserver
	Git       GitObserver
	Clock     Clock
	Random    io.Reader
	Store     Store
}

type identity struct {
	ActivityID, RunID, HandoffID, StartReceiptID, TerminalResultID string
	PublicationKey, ActivityKey, CreatedAtUTC                      string
}

type bindingRequest struct {
	ProjectID      workspace.ProjectID
	RepoID         workspace.RepoID
	TargetWorktree string
	TargetRef      string
	ExpectedOID    string
	Status         StatusPolicy
}

type inputSpec struct {
	ID, Role, Locator, SHA256, MediaType string
	SizeBytes                            int64
	GitBinding                           *GitBinding
}

type handoffDraft struct {
	Value, Goal, Recipient, Authority, Budget, Reporting canonicaljson.Object
	Binding                                              bindingRequest
	Inputs                                               []inputSpec
	Procedure, Verifiers, StopConditions                 []canonicaljson.Value
	PublicationKey, ActivityKey, GoalTitle               string
	Canonical                                            []byte
	Digest                                               string
	MaxRounds                                            int64
	PrincipalID                                          string
}

type startDraft struct {
	Value      canonicaljson.Object
	Canonical  []byte
	ReceiptID  string
	Acceptance string
	Principal  canonicaljson.Object
}

type terminalDraft struct {
	Value            canonicaljson.Object
	Canonical        []byte
	ResultID         string
	Outcome          string
	RoundsUsed       int64
	Summary, Meaning string
	Principal        canonicaljson.Object
	Artifacts        []canonicaljson.Value
}

type handoffDocument struct {
	Value                  canonicaljson.Object
	Bytes                  []byte
	SHA256                 string
	Locator                string
	Identity               identity
	SourceDraftSHA256      string
	GoalTitle              string
	Target                 TargetObservation
	Workspace              workspace.WorkspaceObservation
	ProjectID              workspace.ProjectID
	RepoID                 workspace.RepoID
	RegisteredLocator      string
	RegisteredGitCommonDir string
	ReplyCapabilityID      string
	ReplySecret            string
	ReplyRoot              string
	MaxRounds              int64
}

type acceptedDocument struct {
	Value      canonicaljson.Object
	Bytes      []byte
	SHA256     string
	Locator    string
	DocumentID string
	Outcome    string
	Summary    string
	Meaning    string
}

type Attempt struct {
	Phase, Classification, DocumentID, Digest, Locator, ErrorClass, Detail string
	SizeBytes                                                              int64
	Sequence                                                               int64
}

type Snapshot struct {
	WorkspaceRoot    string
	Handoff          handoffDocument
	Start            *acceptedDocument
	Terminal         *acceptedDocument
	HeadState        string
	ReplacementID    string
	Attempts         []Attempt
	Artifacts        []canonicaljson.Value
	IntegrityReasons []string
}

type CreateStoreInput struct {
	WorkspaceRoot, PublicationKey, SourceDraftSHA256 string
	Handoff                                          handoffDocument
	DeferActivityPublication                         bool
	StoreAlreadyLocked                               bool
	Revalidate                                       func() error
}

type CreateStoreResult struct {
	Snapshot Snapshot
	Created  bool
}

type SupersedeStoreInput struct {
	WorkspaceRoot string
	OldHandoffID  HandoffID
	Reason        string
	Create        CreateStoreInput
}

type SubmitStoreInput struct {
	Snapshot   Snapshot
	Phase      string
	Document   acceptedDocument
	RawInput   []byte
	Artifacts  []managedArtifact
	Revalidate func() error
}

type SubmitStoreResult struct {
	Document acceptedDocument
	Created  bool
}

type ControlStoreInput struct {
	WorkspaceRoot string
	HandoffID     HandoffID
	Reason        string
	State         string
}

type ControlStoreResult struct{ Snapshot Snapshot }

type Store interface {
	Create(CreateStoreInput) (CreateStoreResult, error)
	Supersede(SupersedeStoreInput) (CreateStoreResult, error)
	SubmitStart(SubmitStoreInput) (SubmitStoreResult, error)
	SubmitResult(SubmitStoreInput) (SubmitStoreResult, error)
	Cancel(ControlStoreInput) (ControlStoreResult, error)
	Abandon(ControlStoreInput) (ControlStoreResult, error)
	ReadByID(workspaceRoot string, id HandoffID) (Snapshot, error)
	ReadByLocator(locator string) (Snapshot, error)
}

type managedArtifact struct {
	ID, Locator, SHA256 string
	Bytes               []byte
}
