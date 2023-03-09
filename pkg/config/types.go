package config

type GlobalCloudConfig struct {
	CloudConfigSource CloudConfigSource `yaml:"cloudConfigSource"`
}

type CloudConfigSource struct {
	RootUrl        string `yaml:"rootUrl"`
	RelativFileUrl string `yaml:"relativFileUrl"`
}

type LocalGitConfig struct {
	Git Git `yaml:"git"`
}

type Git struct {
	Url string `yaml:"url"`
}

type GitInfo struct {
	IsRepo       bool
	IsDirty      bool
	EnableCommit bool
}

type SourceProvider struct {
	Host            string   `yaml:"host"`
	AccessToken     string   `yaml:"access_token"`
	ExcludeProjects []string `yaml:"exclude_projects"`
}

type Nexus struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type TerminalConfig struct {
	Width  int    `yaml:"width"`
	Format string `yaml:"format"`
}

type CloudDeprecated struct {
	Type string `json:"type"`
	Data struct {
		Dependencies []CloudDeprecatedDependency `json:"dependencies"`
	} `json:"data"`
}

type Links struct {
	Href      string `json:"href"`
	Title     string `json:"title,omitempty"`
	Templated bool   `json:"templated"`
}

type CloudServices struct {
	Type string         `json:"type"`
	Data []CloudService `json:"data"`
}

type CloudService struct {
	GroupID            string `json:"groupId"`
	ArtifactID         string `json:"artifactId"`
	BuildInfo          string `json:"build-info"`
	DefaultEnvironment string `json:"defaultEnvironment"`
	Environments       []struct {
		Name  string           `json:"name"`
		Links map[string]Links `json:"_links"`
	} `json:"environments"`
}

type CloudDeprecatedDependency struct {
	GroupId    string   `json:"groupId"`
	ArtifactId string   `json:"artifactId"`
	Files      []string `json:"files"`
	Associated struct {
		Files        []string                    `json:"files"`
		Dependencies []CloudDeprecatedDependency `json:"dependencies"`
	} `json:"associated"`
	ReplacementTemplates []string `json:"replacement_templates"`
}

type Directory interface {
	Dir() string
	FilePath(fileName string) (string, error)
}

type DirConfig struct {
	Path string
}

type CloudTemplate struct {
	Name    string
	Project Project
}

type CloudScript struct {
	Name string
	Path string
}

type CloudProjectDefaults struct {
	Type     string          `json:"type"`
	Settings ProjectSettings `json:"settings"`
}
