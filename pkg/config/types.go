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
		Dependencies []struct {
			GroupId    string `json:"groupId"`
			ArtifactId string `json:"artifactId"`
		} `json:"dependencies"`
	} `json:"data"`
}
