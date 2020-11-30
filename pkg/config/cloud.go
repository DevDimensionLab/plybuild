package config

import (
	"github.com/co-pilot-cli/co-pilot/pkg/file"
	"github.com/co-pilot-cli/co-pilot/pkg/logger"
	"github.com/co-pilot-cli/co-pilot/pkg/shell"
	"errors"
	"fmt"
	"io/ioutil"
)

type GitCloudConfig struct {
	Impl DirConfig
}

type CloudConfig interface {
	Implementation() Directory
	Refresh(localConfig LocalConfigFile) error
	Services() func() (CloudServices, error)
	LinkFromService(services func() (CloudServices, error), groupId string, artifactId string, linkKey string) (url string, err error)
	DefaultServiceEnvironmentUrl(service CloudService, key string) (url string, err error)
	Deprecated() (CloudDeprecated, error)
	ListDeprecated() error

	HasTemplate(name string) bool
	ValidTemplatesFrom(list []string) (templates []CloudTemplate, err error)
	Templates() (templates []CloudTemplate, err error)
	Template(name string) (CloudTemplate, error)
}

func (gitCfg GitCloudConfig) Implementation() Directory {
	return gitCfg.Impl
}

func (gitCfg GitCloudConfig) Refresh(localConfig LocalConfigFile) error {
	localCfg, err := localConfig.Config()
	if err != nil {
		return err
	}

	target := gitCfg.Implementation().Dir()
	if file.Exists(file.Path("%s/.git", target)) {
		msg := logger.Info(fmt.Sprintf("pulling cloud config on %s", target))
		log.Info(msg)
		pull := shell.GitPull(target)
		if pull.Err != nil {
			return pull.FormatError()
		}
	} else {
		msg := logger.Info(fmt.Sprintf("cloning %s to %s", localCfg.CloudConfig.Git.Url, target))
		log.Info(msg)
		clone := shell.GitClone(localCfg.CloudConfig.Git.Url, target)
		if clone.Err != nil {
			return clone.FormatError()
		}
	}

	return nil
}

func (gitCfg GitCloudConfig) Services() func() (CloudServices, error) {
	var services CloudServices

	path, err := gitCfg.Implementation().FilePath("services.json")
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

	path, err := gitCfg.Implementation().FilePath("deprecated.json")
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

func (gitCfg GitCloudConfig) HasTemplate(name string) bool {
	templates, err := gitCfg.Templates()
	if err != nil {
		return false
	}
	for _, template := range templates {
		if template.Name == name {
			return true
		}
	}
	return false
}

func (gitCfg GitCloudConfig) Template(name string) (CloudTemplate, error) {
	templates, err := gitCfg.Templates()
	if err != nil {
		return CloudTemplate{}, err
	}

	for _, template := range templates {
		if template.Name == name {
			return template, nil
		}
	}

	return CloudTemplate{}, errors.New(fmt.Sprintf("could not find any valid templates with name: %s", name))
}

func (gitCfg GitCloudConfig) Templates() (templates []CloudTemplate, err error) {
	root := file.Path("%s/templates", gitCfg.Implementation().Dir())
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return
	}

	for _, f := range files {
		if f.IsDir() {
			project, err := InitProjectFromDirectory(file.Path("%s/%s", root, f.Name()))
			if err != nil {
				log.Error(err)
				continue
			}
			templates = append(templates, CloudTemplate{
				Name:    f.Name(),
				Project: project,
			})
		}
	}
	return
}

func (gitCfg GitCloudConfig) ValidTemplatesFrom(list []string) (templates []CloudTemplate, err error) {
	for _, t := range unique(list) {
		template, err := gitCfg.Template(t)
		if err != nil {
			log.Warnln(err)
		} else {
			templates = append(templates, template)
		}
	}

	return
}

func unique(input []string) (list []string) {
	keys := make(map[string]bool)
	for _, entry := range input {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return
}
