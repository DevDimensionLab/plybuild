package springio

type Dependency struct {
	GroupId    string `json:"groupId"`
	ArtifactId string `json:"artifactId"`
	Scope      string `json:"scope"`
	Version    string `json:"version,omitempty"`
	Bom        string `json:"bom,omitempty"`
}

type Bom struct {
	GroupID      string   `json:"groupId"`
	ArtifactID   string   `json:"artifactId"`
	Version      string   `json:"version"`
	Repositories []string `json:"repositories"`
}

type Repository struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	SnapshotEnabled bool   `json:"snapshotEnabled"`
}

type IoDependenciesResponse struct {
	BootVersion  string                `json:"bootVersion"`
	Dependencies map[string]Dependency `json:"dependencies"`
	Repositories map[string]Repository `json:"repositories"`
}

type LinksResponse struct {
	Href      string `json:"href"`
	Title     string `json:"title,omitempty"`
	Templated bool   `json:"templated"`
}

type ValueResponse struct {
	Id           string `json:"id,omitempty"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	VersionRange string `json:"versionRange,omitempty"`
	Action       string `json:"action,omitempty"`
	//Links        map[string]LinksResponse `json:"_links,omitempty"`
	Tags map[string]string `json:"tags,omitempty"`
}

type IoRootResponse struct {
	Links      map[string]LinksResponse `json:"_links"`
	ArtifactId struct {
		Default string `json:"default"`
		Type    string `json:"text"`
	} `json:"artifactId"`
	BootVersion struct {
		Default string          `json:"default"`
		Type    string          `json:"type"`
		Values  []ValueResponse `json:"values"`
	} `json:"bootVersion"`
	Dependencies struct {
		Type   string `json:"type"`
		Values []struct {
			Name   string          `json:"name"`
			Values []ValueResponse `json:"values"`
		}
	} `json:"dependencies"`
	Description struct {
		Default string `json:"default"`
		Type    string `json:"type"`
	} `json:"description"`
	JavaVersion struct {
		Default string          `json:"default"`
		Type    string          `json:"type"`
		Values  []ValueResponse `json:"values"`
	} `json:"javaVersion"`
	Language struct {
		Default string          `json:"default"`
		Type    string          `json:"type"`
		Values  []ValueResponse `json:"values"`
	} `json:"language"`
	Name struct {
		Default string `json:"default"`
		Type    string `json:"type"`
	} `json:"name"`
	GroupId struct {
		Default string `json:"default"`
		Type    string `json:"type"`
	} `json:"groupId"`
	PackageName struct {
		Default string `json:"default"`
		Type    string `json:"type"`
	} `json:"packageName"`
	Packaging struct {
		Default string          `json:"default"`
		Type    string          `json:"type"`
		Values  []ValueResponse `json:"values"`
	} `json:"packaging"`
	Type struct {
		Default string          `json:"default"`
		Type    string          `json:"type"`
		Values  []ValueResponse `json:"values"`
	} `json:"type"`
	Version struct {
		Default string `json:"default"`
		Type    string `json:"type"`
	} `json:"version"`
}

type IoInfoResponse struct {
	BomRanges map[string]map[string]string `json:"bom-ranges"`
	Build     struct {
		Artifact string            `json:"artifact"`
		Group    string            `json:"group"`
		Name     string            `json:"name"`
		Time     string            `json:"time"`
		Version  string            `json:"version"`
		Versions map[string]string `json:"versions"`
	} `json:"build"`
	DependencyRanges map[string]map[string]string `json:"dependency-ranges"`
	Git              struct {
		Branch string            `json:"branch"`
		Commit map[string]string `json:"commit"`
	} `json:"git"`
}
