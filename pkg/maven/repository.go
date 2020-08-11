package maven

import (
	"co-pilot/pkg/file"
	"errors"
	"fmt"
	"os/user"
)

func GetRepositories() ([]string, error) {
	var repos []string
	defaultMavenRepo := "https://repo1.maven.org/maven2"
	settingsFile, err := GetSettingsFile()

	if err == nil {
		var settings Settings
		err = file.ReadXml(settingsFile, &settings)
		if err != nil {
			return repos, err
		}

		for _, profile := range settings.Profiles.Profile {
			for _, repo := range profile.Repositories.Repository {
				if repo.Releases.Enabled && repo.URL != "" {
					repos = append(repos, repo.URL)
				}
			}
		}

		for _, mirror := range settings.Mirrors {
			if mirror.Mirror.URL != "" {
				repos = append(repos, mirror.Mirror.URL)
			}
		}

	} else {
		// could not find settings.xml, adding defaultMavenRepo
		repos = append(repos, defaultMavenRepo)
	}

	return repos, nil
}

func GetSettingsFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	home := usr.HomeDir
	m2Settings := fmt.Sprintf("%s/.m2/settings.xml", home)
	confSettings := fmt.Sprintf("%s/conf/settings.xml", home)

	if file.Exists(m2Settings) {
		return m2Settings, nil
	} else if file.Exists(confSettings) {
		return confSettings, nil
	}

	return "", errors.New("could not find settings.xml")
}
