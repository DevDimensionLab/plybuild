package maven

import (
	"co-pilot/pkg/file"
	"errors"
	"fmt"
	"os/user"
)

type Repositories struct {
	Fallback string
	Profile  []string
	Mirror   []string
}

func GetRepositories() (Repositories, error) {
	repos := Repositories{
		Fallback: "https://repo1.maven.org/maven2",
	}
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
					repos.Profile = append(repos.Profile, repo.URL)
				}
			}
		}

		for _, mirror := range settings.Mirrors {
			if mirror.Mirror.URL != "" {
				repos.Mirror = append(repos.Mirror, mirror.Mirror.URL)
			}
		}
	}

	return repos, nil
}

func ListRepositories() error {
	repos, err := GetRepositories()
	if err != nil {
		return err
	}

	for _, profileRepo := range repos.Profile {
		log.Infof("found profile repo: %s", profileRepo)
	}

	for _, mirrorRepo := range repos.Mirror {
		log.Infof("found mirror repo: %s", mirrorRepo)
	}

	log.Infof("fallback repo: %s", repos.Fallback)

	return nil
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
