package resources

import (
	"github.com/devdimensionlab/ply/pkg/config"
	"github.com/devdimensionlab/ply/pkg/file"
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
