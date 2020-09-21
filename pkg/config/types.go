package config

type LocalConfiguration struct {
	CloudConfig    LocalGitConfig `yaml:"cloudConfig"`
	SourceProvider SourceProvider `yaml:"sourceProvider"`
}

type LocalGitConfig struct {
	Git Git `yaml:"git"`
}

type Git struct {
	Url string `yaml:"url"`
}

type SourceProvider struct {
	Host        string `yaml:"host"`
	AccessToken string `yaml:"access_token"`
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
	Name string
	Impl DirConfig
}
