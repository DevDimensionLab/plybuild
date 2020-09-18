package config

import (
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/shell"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
)

type GitCloudConfig struct {
	configDirName string
}

type CloudConfig interface {
	Dir() string
	Refresh(localConfig LocalConfigFile) error
	Services() func() (CloudServices, error)
	LinkFromService(services func() (CloudServices, error), groupId string, artifactId string, linkKey string) (url string, err error)
	DefaultServiceEnvironmentUrl(service CloudService, key string) (url string, err error)
	Deprecated() (CloudDeprecated, error)
	ListDeprecated() error
	FilePath(fileName string) (string, error)
}

func InitGitCloudConfig(cloudConfigDirName string) (cfg GitCloudConfig, err error) {
	home, err := homedir.Dir()
	if err != nil {
		return cfg, err
	}

	cfg.configDirName = fmt.Sprintf("%s/%s/%s", home, localConfigDir, cloudConfigDirName)
	return
}

func (gitCfg GitCloudConfig) Dir() string {
	return gitCfg.configDirName
}

func (gitCfg GitCloudConfig) Refresh(localConfig LocalConfigFile) error {
	localCfg, err := localConfig.Config()
	if err != nil {
		log.Fatalln(err)
	}

	target := gitCfg.Dir()
	if file.Exists(fmt.Sprintf("%s/.git", target)) {
		msg := logger.Info(fmt.Sprintf("pulling cloud config on %s", target))
		log.Info(msg)
		out, err := shell.GitPull(target)
		if err != nil {
			return errors.New(fmt.Sprintf("pulling cloud config failed:\n%s, %v", out, err))
		}
	} else {
		msg := logger.Info(fmt.Sprintf("cloning %s to %s", localCfg.CloudConfig.Git.Url, target))
		log.Info(msg)
		out, err := shell.GitClone(localCfg.CloudConfig.Git.Url, target)
		if err != nil {
			return errors.New(fmt.Sprintf("ploning cloud config failed:\n%s, %v", out, err))
		}
	}

	return nil
}

func (gitCfg GitCloudConfig) Services() func() (CloudServices, error) {
	var services CloudServices

	path, err := gitCfg.FilePath("services.json")
	if err != nil {
		return func() (CloudServices, error) {
			return services, err
		}
	}

	err = file.ReadJson(path, &services)
	if err != nil {
		return func() (CloudServices, error) {
			return services, err
		}
	}

	return func() (CloudServices, error) {
		return services, nil
	}
}

func (gitCfg GitCloudConfig) LinkFromService(services func() (CloudServices, error), groupId string, artifactId string, linkKey string) (url string, err error) {
	s, err := services()
	if err != nil {
		return
	}

	if s.Data != nil {
		for _, service := range s.Data {
			if service.GroupID == groupId && service.ArtifactID == artifactId {
				url, err := gitCfg.DefaultServiceEnvironmentUrl(service, linkKey)
				if err != nil {
					return url, err
				}
				log.Debugf("found service %s:%s with link-key %s, requesting %s", groupId, artifactId, linkKey, url)
				return url, err
			}
		}
	}

	return url, errors.New(fmt.Sprintf("could not get cloud config service information on %s:%s", groupId, artifactId))
}

func (gitCfg GitCloudConfig) DefaultServiceEnvironmentUrl(service CloudService, key string) (url string, err error) {
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

func (gitCfg GitCloudConfig) Deprecated() (CloudDeprecated, error) {
	var deprecated CloudDeprecated

	path, err := gitCfg.FilePath("deprecated.json")
	if err != nil {
		return deprecated, err
	}

	err = file.ReadJson(path, &deprecated)
	if err != nil {
		return deprecated, err
	}

	return deprecated, nil
}

func (gitCfg GitCloudConfig) ListDeprecated() error {
	deprecated, err := gitCfg.Deprecated()
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

func (gitCfg GitCloudConfig) FilePath(fileName string) (string, error) {
	path := fmt.Sprintf("%s/%s", gitCfg.Dir(), fileName)
	if !file.Exists(path) {
		return "", errors.New(fmt.Sprintf("could not find %s in cloud config", fileName))
	}

	return path, nil
}
