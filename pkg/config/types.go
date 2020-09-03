package config

type GlobalConfiguration struct {
	CloudConfig    CloudConfig    `yaml:"cloudConfig"`
	SourceProvider SourceProvider `yaml:"sourceProvider"`
}

type CloudConfig struct {
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

type CloudDeprecatedDependency struct {
	GroupId    string   `json:"groupId"`
	ArtifactId string   `json:"artifactId"`
	Files      []string `json:"files"`
	Associated struct {
		Files        []string                    `json:"files"`
		Dependencies []CloudDeprecatedDependency `json:"dependencies"`
	} `json:"associated"`
	Replacements struct {
		Files        []string                    `json:"files"`
		Dependencies []CloudDeprecatedDependency `json:"dependencies"`
	} `json:"replacements"`
}
