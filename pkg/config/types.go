package config

type InitConfiguration struct {
	GroupId      string   `json:"groupId"`
	ArtifactId   string   `json:"artifactId"`
	Package      string   `json:"package"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Dependencies []string `json:"dependencies"`
}

type GlobalConfiguration struct {
	BannedPomUrl         string `yaml:"banned_pom_url"`
	BitBucketHost        string `yaml:"bitbucket_host"`
	BitBucketAccessToken string `yaml:"bitbucket_personal_access_token"`
}
