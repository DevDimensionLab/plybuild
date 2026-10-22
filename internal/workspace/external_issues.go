package workspace

type ExternalProvider string

const (
	ExternalJira         ExternalProvider = "jira"
	ExternalGitHub       ExternalProvider = "github"
	ExternalIssueAdd                      = "add"
	ExternalIssueUnlink                   = "unlink"
	ExternalIssuePrimary                  = "primary"
)

type ExternalIssue struct {
	JiraContextPath string           `json:"jira_context_path"`
	ID              string           `json:"id"`
	Provider        ExternalProvider `json:"provider"`
	Host            string           `json:"host"`
	Key             string           `json:"key"`
	Owner           string           `json:"owner"`
	Repository      string           `json:"repository"`
	Number          int64            `json:"number"`
	URL             string           `json:"url"`
}
type ExternalCommitReference struct {
	RefID       string `json:"ref_id"`
	CommitKey   string `json:"commit_key"`
	SubjectKind string `json:"subject_kind"`
	SubjectID   string `json:"subject_id"`
}
type ExternalIssuesReadback struct {
	Refs           []ExternalIssue          `json:"refs"`
	Primary        *ExternalCommitReference `json:"primary"`
	Effective      *ExternalCommitReference `json:"effective"`
	Revision       int                      `json:"revision"`
	Kind           string                   `json:"kind"`
	SchemaVersion  int                      `json:"schema_version"`
	Workspace      MetadataWorkspace        `json:"workspace"`
	SubjectKind    string                   `json:"subject_kind"`
	SubjectID      string                   `json:"subject_id"`
	ProviderStatus string                   `json:"provider_status"`
	Events         []ExternalIssueEvent     `json:"events"`
	Source         MetadataSource           `json:"source"`
}
type ExternalIssuesMutation struct {
	ExternalIssuesReadback
	Changed bool `json:"changed"`
}
type ExternalIssueInput struct {
	SubjectKind, SubjectID, Operation, RefID, ActorClaim string
	Issue                                                *ExternalIssue
	Primary                                              bool
	ExpectedRevision                                     *int
}
