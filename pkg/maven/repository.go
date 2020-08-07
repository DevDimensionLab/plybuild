package maven

import (
	"errors"
	"fmt"
	"os/user"
	"spring-boot-co-pilot/pkg/file"
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
		for _, profile := range settings.Profiles {
			for _, repo := range profile.Profile.Repositories {
				repos = append(repos, repo.Repository.URL)
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
