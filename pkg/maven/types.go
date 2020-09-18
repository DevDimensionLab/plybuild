package maven

import (
	"encoding/xml"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
)

type PomPair struct {
	PomFile string
	Model   *pom.Model
}

type Metadata struct {
	XMLName      xml.Name `xml:"metadata"`
	Text         string   `xml:",chardata"`
	ModelVersion string   `xml:"modelVersion,attr"`
	GroupId      string   `xml:"groupId"`
	ArtifactId   string   `xml:"artifactId"`
	Versioning   struct {
		Text     string `xml:",chardata"`
		Latest   string `xml:"latest"`
		Release  string `xml:"release"`
		Versions struct {
			Text    string   `xml:",chardata"`
			Version []string `xml:"version"`
		} `xml:"versions"`
		LastUpdated string `xml:"lastUpdated"`
	} `xml:"versioning"`
}

type Settings struct {
	XMLName           xml.Name `xml:"settings"`
	Text              string   `xml:",chardata"`
	LocalRepository   string   `xml:"localRepository"`
	InteractiveMode   string   `xml:"interactiveMode"`
	UsePluginRegistry string   `xml:"usePluginRegistry"`
	Offline           string   `xml:"offline"`
	Proxies           []struct {
		Text  string `xml:",chardata"`
		Proxy struct {
			Text          string `xml:",chardata"`
			Active        string `xml:"active"`
			Protocol      string `xml:"protocol"`
			Username      string `xml:"username"`
			Password      string `xml:"password"`
			Port          string `xml:"port"`
			Host          string `xml:"host"`
			NonProxyHosts string `xml:"nonProxyHosts"`
			ID            string `xml:"id"`
		} `xml:"proxy"`
	} `xml:"proxies"`
	Servers []struct {
		Text   string `xml:",chardata"`
		Server struct {
			Text                 string `xml:",chardata"`
			Username             string `xml:"username"`
			Password             string `xml:"password"`
			PrivateKey           string `xml:"privateKey"`
			Passphrase           string `xml:"passphrase"`
			FilePermissions      string `xml:"filePermissions"`
			DirectoryPermissions string `xml:"directoryPermissions"`
			Configuration        string `xml:"configuration"`
			ID                   string `xml:"id"`
		} `xml:"server"`
	} `xml:"servers"`
	Mirrors []struct {
		Text   string `xml:",chardata"`
		Mirror struct {
			Text     string `xml:",chardata"`
			MirrorOf string `xml:"mirrorOf"`
			Name     string `xml:"name"`
			URL      string `xml:"url"`
			ID       string `xml:"id"`
		} `xml:"mirror"`
	} `xml:"mirrors"`
	Profiles struct {
		Text    string `xml:",chardata"`
		Profile []struct {
			Text       string `xml:",chardata"`
			Activation struct {
				Text            string `xml:",chardata"`
				ActiveByDefault string `xml:"activeByDefault"`
				Jdk             string `xml:"jdk"`
				Os              struct {
					Text    string `xml:",chardata"`
					Name    string `xml:"name"`
					Family  string `xml:"family"`
					Arch    string `xml:"arch"`
					Version string `xml:"version"`
				} `xml:"os"`
				Property struct {
					Text  string `xml:",chardata"`
					Name  string `xml:"name"`
					Value string `xml:"value"`
				} `xml:"property"`
				File struct {
					Text    string `xml:",chardata"`
					Missing string `xml:"missing"`
					Exists  string `xml:"exists"`
				} `xml:"file"`
			} `xml:"activation"`
			Properties   string `xml:"properties"`
			Repositories struct {
				Text       string `xml:",chardata"`
				Repository []struct {
					Text     string `xml:",chardata"`
					Releases struct {
						Text           string `xml:",chardata"`
						Enabled        bool   `xml:"enabled"`
						UpdatePolicy   string `xml:"updatePolicy"`
						ChecksumPolicy string `xml:"checksumPolicy"`
					} `xml:"releases"`
					Snapshots struct {
						Text           string `xml:",chardata"`
						Enabled        string `xml:"enabled"`
						UpdatePolicy   string `xml:"updatePolicy"`
						ChecksumPolicy string `xml:"checksumPolicy"`
					} `xml:"snapshots"`
					ID     string `xml:"id"`
					Name   string `xml:"name"`
					URL    string `xml:"url"`
					Layout string `xml:"layout"`
				} `xml:"repository"`
			} `xml:"repositories"`
			PluginRepositories []struct {
				Text             string `xml:",chardata"`
				PluginRepository struct {
					Text     string `xml:",chardata"`
					Releases struct {
						Text           string `xml:",chardata"`
						Enabled        string `xml:"enabled"`
						UpdatePolicy   string `xml:"updatePolicy"`
						ChecksumPolicy string `xml:"checksumPolicy"`
					} `xml:"releases"`
					Snapshots struct {
						Text           string `xml:",chardata"`
						Enabled        string `xml:"enabled"`
						UpdatePolicy   string `xml:"updatePolicy"`
						ChecksumPolicy string `xml:"checksumPolicy"`
					} `xml:"snapshots"`
					ID     string `xml:"id"`
					Name   string `xml:"name"`
					URL    string `xml:"url"`
					Layout string `xml:"layout"`
				} `xml:"pluginRepository"`
			} `xml:"pluginRepositories"`
			ID string `xml:"id"`
		} `xml:"profile"`
	} `xml:"profiles"`
	ActiveProfiles []struct {
		Text          string `xml:",chardata"`
		ActiveProfile string `xml:"activeProfile"`
	} `xml:"activeProfiles"`
	PluginGroups []struct {
		Text        string `xml:",chardata"`
		PluginGroup string `xml:"pluginGroup"`
	} `xml:"pluginGroups"`
}

type JavaVersion struct {
	Major           int
	Minor           int
	Patch           int
	Suffix          string
	SuffixSeparator string
}
