package config

import (
	"co-pilot/pkg/file"
	"co-pilot/pkg/http"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/shell"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
)

var globalConfigDir = "cloud-config"

func GlobalConfigDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", home, localConfigDir, globalConfigDir), nil
}

func Refresh() error {
	c, err := GetLocalConfig()
	if err != nil {
		log.Fatalln(err)
	}

	target, err := GlobalConfigDir()
	if err != nil {
		log.Fatalln(err)
	}

	if file.Exists(fmt.Sprintf("%s/.git", target)) {
		msg := logger.Info(fmt.Sprintf("pulling cloud config on %s", target))
		log.Info(msg)
		out, err := shell.GitPull(target)
		if err != nil {
			return errors.New(fmt.Sprintf("pulling cloud config failed:\n%s, %v", out, err))
		}
	} else {
		msg := logger.Info(fmt.Sprintf("cloning %s to %s", c.CloudConfig.Git.Url, target))
		log.Info(msg)
		out, err := shell.GitClone(c.CloudConfig.Git.Url, target)
		if err != nil {
			return errors.New(fmt.Sprintf("ploning cloud config failed:\n%s, %v", out, err))
		}
	}

	return nil
}

func getServices() (CloudServices, error) {
	var deprecated CloudServices

	path, err := GetCloudConfigFilePath("services.json")
	if err != nil {
		return deprecated, err
	}

	err = file.ReadJson(path, &deprecated)
	if err != nil {
		return deprecated, err
	}

	return deprecated, nil
}

func GetDataFromService(groupId string, artifactId string, linkKey string) (data map[string]interface{}, err error) {
	services, err := getServices()
	if err != nil {
		return
	}

	if services.Data != nil {
		for _, service := range services.Data {
			if service.GroupID == groupId && service.ArtifactID == artifactId {
				url, err := GetDefaultServiceEnvironmentUrl(service, linkKey)
				if err != nil {
					return data, err
				}
				log.Debugf("found service %s:%s with link-key %s, requesting %s", groupId, artifactId, linkKey, url)
				err = http.GetJson(url, &data)
				return data, err
			}
		}
	}

	return data, errors.New(fmt.Sprintf("could not get cloud config service information on %s:%s", groupId, artifactId))
}

func GetDefaultServiceEnvironmentUrl(service CloudService, key string) (url string, err error) {
	for _, environment := range service.Environments {
		if environment.Name == service.DefaultEnvironment {
			if link, ok := environment.Links[key]; ok {
				url = link.Href
				return
			}
		}
	}

	return url, errors.New(fmt.Sprintf("could not find environment with link key %s", key))
}

func GetDeprecated() (CloudDeprecated, error) {
	var deprecated CloudDeprecated

	path, err := GetCloudConfigFilePath("deprecated.json")
	if err != nil {
		return deprecated, err
	}

	err = file.ReadJson(path, &deprecated)
	if err != nil {
		return deprecated, err
	}

	return deprecated, nil
}

func ListDeprecated() error {
	deprecated, err := GetDeprecated()
	if err != nil {
		return err
	}

	for _, dep := range deprecated.Data.Dependencies {
		log.Infof("== deprecated dependency %s:%s ==", dep.GroupId, dep.ArtifactId)
		if dep.Associated.Dependencies != nil {
			for _, assoc := range dep.Associated.Dependencies {
				log.Infof("\t <= associated deprecated dependency %s:%s", assoc.GroupId, assoc.ArtifactId)
			}
		}
		if dep.ReplacementTemplates != nil {
			for _, repTemp := range dep.ReplacementTemplates {
				log.Infof("\t <= replacement template %s", repTemp)
			}
		}
	}

	return nil
}

func GetCloudConfigFilePath(fileName string) (string, error) {
	err := Refresh()
	if err != nil {
		return "", err
	}

	cloudConfigDir, err := GlobalConfigDir()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/%s", cloudConfigDir, fileName)
	if !file.Exists(path) {
		return "", errors.New(fmt.Sprintf("could not find %s in cloud config", fileName))
	}

	return path, nil
}
