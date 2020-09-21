package config

import (
	"co-pilot/pkg/file"
	"errors"
	"fmt"
)

func (dirCfg DirConfig) Dir() string {
	return dirCfg.Path
}

func (dirCfg DirConfig) FilePath(fileName string) (string, error) {
	path := fmt.Sprintf("%s/%s", dirCfg.Dir(), fileName)
	if !file.Exists(path) {
		return "", errors.New(fmt.Sprintf("could not find %s in cloud config", fileName))
	}

	return path, nil
}
