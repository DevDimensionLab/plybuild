package config

import (
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/file"
	"github.com/mitchellh/go-homedir"
)

const coPilotHomePath = ".co-pilot"

func GetCoPilotHomePath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", home, coPilotHomePath), nil
}

func GetProfilesPath() (string, error) {
	home, err := GetCoPilotHomePath()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/profiles", home), nil
}

func GetActiveProfilePath() (string, error) {
	profilesPath, err := GetProfilesPath()
	if err != nil {
		return "", err
	}

	activeProfile, err := file.OpenLinesStrict(fmt.Sprintf("%s/.active_profile", profilesPath))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", profilesPath, activeProfile[0]), nil
}

func SwitchProfile(newProfile string) error {
	profilesPath, err := GetProfilesPath()
	if err != nil {
		return err
	}

	return file.CreateFile(fmt.Sprintf("%s/.active_profile", profilesPath), newProfile)
}

func MigrateToProfiles() error {
	profilesPath, err := GetProfilesPath()
	if err != nil {
		return err
	}
	homePath, err := GetCoPilotHomePath()
	if err != nil {
		return err
	}

	if err := file.CreateDirectory(fmt.Sprintf("%s/default", profilesPath)); err != nil {
		return err
	}
	for _, f := range []string{"local-config.yaml", "cloud-config"} {
		if err := file.Move(fmt.Sprintf("%s/%s", homePath, f), fmt.Sprintf("%s/default/%s", profilesPath, f)); err != nil {
			log.Debugf(err.Error())
		}
	}

	return file.CreateFile(fmt.Sprintf("%s/.active_profile", profilesPath), "default")
}
