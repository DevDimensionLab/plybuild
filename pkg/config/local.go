package config

import (
	"github.com/devdimensionlab/co-pilot/pkg/file"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var localConfigFileName = "local-config.yaml"
var defaultCloudConfigUrl = "https://github.com/devdimensionlab/co-pilot-config.git"

type LocalConfig struct {
	impl DirConfig
}

type LocalConfiguration struct {
	CloudConfig    LocalGitConfig `yaml:"cloudConfig"`
	SourceProvider SourceProvider `yaml:"sourceProvider"`
}

type LocalConfigFile interface {
	Implementation() DirConfig
	FilePath() string
	CheckOrCreateConfigDir() error
	TouchFile() error
	Config() (LocalConfiguration, error)
	Print() error
	Exists() bool
}

func NewLocalConfig(absConfigDir string) (cfg LocalConfig) {
	cfg.impl.Path = file.Path(absConfigDir)
	return
}

func (localCfg LocalConfig) Implementation() DirConfig {
	return localCfg.impl
}

func (localCfg LocalConfig) FilePath() string {
	return file.Path("%s/%s", localCfg.impl.Path, localConfigFileName)
}

func (localCfg LocalConfig) CheckOrCreateConfigDir() error {
	dir := localCfg.Implementation().Path

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func (localCfg LocalConfig) TouchFile() error {
	err := localCfg.CheckOrCreateConfigDir()
	if err != nil {
		return err
	}

	configFile := localCfg.FilePath()

	config := LocalConfiguration{}
	config.CloudConfig.Git.Url = defaultCloudConfigUrl
	d, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	log.Infof("creating new config file %s", configFile)

	f, err := os.Create(configFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFile, d, 0644)
	if err != nil {
		return err
	}

	return f.Close()
}

func (localCfg LocalConfig) Config() (LocalConfiguration, error) {
	config := LocalConfiguration{}
	localConfigFile := localCfg.FilePath()

	b, err := file.Open(localConfigFile)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (localCfg LocalConfig) Print() error {
	c, err := localCfg.Config()
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	log.Infof("using: %s", localCfg.FilePath())
	log.Infof("\n%s\n", string(b))

	return nil
}

func (localCfg LocalConfig) Exists() bool {
	return file.Exists(localCfg.FilePath())
}
