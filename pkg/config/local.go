package config

import (
	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"gopkg.in/yaml.v2"
	"io/fs"
	"os"
)

var localConfigFileName = "local-config.yaml"
var defaultCloudConfigUrl = "https://github.com/devdimensionlab/plybuild-config.git"

type localConfigDirectoryStatDependencies struct {
	Files filesystem.Dependencies
}

func systemLocalConfigDirectoryStatDependencies() localConfigDirectoryStatDependencies {
	return localConfigDirectoryStatDependencies{Files: filesystem.System()}
}

func statLocalConfigDirectory(dependencies localConfigDirectoryStatDependencies, dir string) (fs.FileInfo, error) {
	return filesystem.Stat(dependencies.Files, dir)
}

type localConfigTouchWriteDependencies struct {
	Files filesystem.Dependencies
}

type localConfigTouchCreateDependencies struct {
	Files filesystem.Dependencies
}

func systemLocalConfigTouchCreateDependencies() localConfigTouchCreateDependencies {
	return localConfigTouchCreateDependencies{Files: filesystem.System()}
}

func createLocalConfigTouch(dependencies localConfigTouchCreateDependencies, configFilePath string) (filesystem.File, error) {
	return filesystem.Create(dependencies.Files, configFilePath)
}

func systemLocalConfigTouchWriteDependencies() localConfigTouchWriteDependencies {
	return localConfigTouchWriteDependencies{Files: filesystem.System()}
}

func writeLocalConfigTouch(dependencies localConfigTouchWriteDependencies, configFilePath string, data []byte) error {
	return filesystem.WriteFile(dependencies.Files, configFilePath, data, 0644)
}

type localConfigUpdateWriteDependencies struct {
	Files filesystem.Dependencies
}

type localConfigUpdateCreateDependencies struct {
	Files filesystem.Dependencies
}

func systemLocalConfigUpdateCreateDependencies() localConfigUpdateCreateDependencies {
	return localConfigUpdateCreateDependencies{Files: filesystem.System()}
}

func createLocalConfigUpdate(dependencies localConfigUpdateCreateDependencies, configFilePath string) (filesystem.File, error) {
	return filesystem.Create(dependencies.Files, configFilePath)
}

func systemLocalConfigUpdateWriteDependencies() localConfigUpdateWriteDependencies {
	return localConfigUpdateWriteDependencies{Files: filesystem.System()}
}

func writeLocalConfigUpdate(dependencies localConfigUpdateWriteDependencies, configFilePath string, data []byte) error {
	return filesystem.WriteFile(dependencies.Files, configFilePath, data, 0644)
}

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

	if _, err := statLocalConfigDirectory(systemLocalConfigDirectoryStatDependencies(), dir); os.IsNotExist(err) {
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

	f, err := createLocalConfigTouch(systemLocalConfigTouchCreateDependencies(), configFilePath)
	if err != nil {
		return err
	}

	err = writeLocalConfigTouch(systemLocalConfigTouchWriteDependencies(), configFilePath, d)
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

	f, err := createLocalConfigUpdate(systemLocalConfigUpdateCreateDependencies(), configFilePath)
	if err != nil {
		return err
	}

	err = writeLocalConfigUpdate(systemLocalConfigUpdateWriteDependencies(), configFilePath, d)
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
	if terminalConfig.Width == 0 {
		terminalConfig.Width = 80
	}

	return terminalConfig, nil
}
