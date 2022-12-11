package config

import (
	"github.com/devdimensionlab/co-pilot/pkg/file"
	"gopkg.in/yaml.v2"
	"os"
)

var localConfigFileName = "local-config.yaml"
var defaultCloudConfigUrl = "https://github.com/devdimensionlab/co-pilot-config.git"

type LocalConfigDir struct {
	impl DirConfig
}

type LocalConfiguration struct {
	CloudConfig    LocalGitConfig `yaml:"cloudConfig"`
	SourceProvider SourceProvider `yaml:"sourceProvider"`
	Nexus          Nexus          `yaml:"nexus"`
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

func OpenLocalConfig(absConfigDir string) (cfg LocalConfigDir) {
	cfg.impl.Path = file.Path(absConfigDir)
	return
}

func (localCfg LocalConfigDir) Implementation() DirConfig {
	return localCfg.impl
}

func (localCfg LocalConfigDir) FilePath() string {
	return file.Path("%s/%s", localCfg.impl.Path, localConfigFileName)
}

func (localCfg LocalConfigDir) CheckOrCreateConfigDir() error {
	dir := localCfg.Implementation().Path

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func (localCfg LocalConfigDir) TouchFile() error {
	err := localCfg.CheckOrCreateConfigDir()
	if err != nil {
		return err
	}

	configFilePath := localCfg.FilePath()

	config := LocalConfiguration{}
	config.CloudConfig.Git.Url = defaultCloudConfigUrl
	d, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	log.Infof("creating new config file %s", configFilePath)

	f, err := os.Create(configFilePath)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, d, 0644)
	if err != nil {
		return err
	}

	return f.Close()
}

func (localCfg LocalConfigDir) Config() (LocalConfiguration, error) {
	config := LocalConfiguration{}
	localConfigFile := localCfg.FilePath()

	b, err := file.Open(localConfigFile)
	if err != nil {
		return config, err
	}

	b = []byte(os.ExpandEnv(string(b)))

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (localCfg LocalConfigDir) Print() error {
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

func (localCfg LocalConfigDir) Exists() bool {
	return file.Exists(localCfg.FilePath())
}
