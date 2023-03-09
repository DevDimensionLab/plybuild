package config

import (
	"encoding/json"
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
	TerminalConfig TerminalConfig `yaml:"terminal"`
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

func (localCfgDir LocalConfigDir) Implementation() DirConfig {
	return localCfgDir.impl
}

func (localCfgDir LocalConfigDir) FilePath() string {
	return file.Path("%s/%s", localCfgDir.impl.Path, localConfigFileName)
}

func (localCfgDir LocalConfigDir) CheckOrCreateConfigDir() error {
	dir := localCfgDir.Implementation().Path

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func (localCfgDir LocalConfigDir) TouchFile() error {
	err := localCfgDir.CheckOrCreateConfigDir()
	if err != nil {
		return err
	}

	configFilePath := localCfgDir.FilePath()

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

func (localCfgDir LocalConfigDir) UpdateLocalConfig(config LocalConfiguration) error {
	err := localCfgDir.CheckOrCreateConfigDir()
	if err != nil {
		return err
	}

	configFilePath := localCfgDir.FilePath()

	d, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	log.Infof("update config file %s", configFilePath)

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

func (localCfgDir LocalConfigDir) Config() (LocalConfiguration, error) {
	config := LocalConfiguration{}
	localConfigFile := localCfgDir.FilePath()

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

func (localCfgDir LocalConfigDir) ConfigAsMap() (map[string]any, error) {
	config, err := localCfgDir.Config()
	if err != nil {
		return nil, err
	}

	serializedBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	var deserialized map[string]any
	err = json.Unmarshal(serializedBytes, &deserialized)
	return deserialized, err
}

func (localCfgDir LocalConfigDir) Print() error {
	c, err := localCfgDir.Config()
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	log.Infof("using: %s", localCfgDir.FilePath())
	log.Infof("\n%s\n", string(b))

	return nil
}

func (localCfgDir LocalConfigDir) Exists() bool {
	return file.Exists(localCfgDir.FilePath())
}

func (localCfgDir LocalConfigDir) GetTerminalConfig() (TerminalConfig, error) {
	cfg, err := localCfgDir.Config()
	if err != nil {
		return TerminalConfig{}, err
	}

	terminalConfig := cfg.TerminalConfig
	if 0 == terminalConfig.Width {
		terminalConfig.Width = 80
	}

	return terminalConfig, nil
}
