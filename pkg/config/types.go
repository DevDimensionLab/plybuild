package config

type InitConfiguration struct {
	GroupId      string   `json:"groupId"`
	ArtifactId   string   `json:"artifactId"`
	Package      string   `json:"package"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Dependencies []string `json:"dependencies"`
}
