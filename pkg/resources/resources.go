package resources

import (
	"github.com/devdimensionlab/co-pilot/pkg/config"
	"github.com/devdimensionlab/co-pilot/pkg/file"
)

func LocalDir(gitCfg config.CloudConfig) string {
	return file.Path("%s/%s", gitCfg.Implementation().Dir(), "resources")
}

func ResourceAsString(gitCfg config.CloudConfig, filename string) (string, error) {
	resourcesDir := LocalDir(gitCfg)
	content, err := file.Open(resourcesDir + "/" + filename)
	if err != nil {
		return "", err
	}
	return string(content[:]), nil
}
