package config

import (
	"github.com/devdimensionlab/co-pilot/pkg/file"
	"github.com/mitchellh/go-homedir"
)

const coPilotHomePath = ".co-pilot"

func GetCoPilotHomePath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return file.Path("%s/%s", home, coPilotHomePath), nil
}

func GetProfilesPath() (string, error) {
	home, err := GetCoPilotHomePath()
	if err != nil {
		return "", err
	}

	return file.Path("%s/profiles", home), nil
}

func GetProfilesPathFor(profile string) (string, error) {
	profilesPath, err := GetProfilesPath()
	if err != nil {
		return "", err
	}
	return file.Path("%s/%s", profilesPath, profile), nil
}

func GetActiveProfilePath() (string, error) {
	profilesPath, err := GetProfilesPath()
	if err != nil {
		return "", err
	}

	activeProfile, err := file.OpenLinesStrict(file.Path("%s/.active_profile", profilesPath))
	if err != nil {
		return "", err
	}

	return file.Path("%s/%s", profilesPath, activeProfile[0]), nil
}

func SwitchProfile(newProfile string) error {
	profilesPath, err := GetProfilesPath()
	if err != nil {
		return err
	}

	return file.CreateFile(file.Path("%s/.active_profile", profilesPath), newProfile)
}

func InstallOrMigrateToProfiles() error {
	profilesPath, err := GetProfilesPath()
	if err != nil {
		return err
	}
	homePath, err := GetCoPilotHomePath()
	if err != nil {
		return err
	}

	if err := file.CreateDirectory(file.Path("%s/default", profilesPath)); err != nil {
		return err
	}
	for _, f := range []string{"local-config.yaml", "cloud-config"} {
		if err := file.Move(file.Path("%s/%s", homePath, f), file.Path("%s/default/%s", profilesPath, f)); err != nil {
			log.Debugf(err.Error())
		}
	}

	return file.CreateFile(file.Path("%s/.active_profile", profilesPath), "default")
}
