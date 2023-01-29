package tips

import (
	"github.com/devdimensionlab/co-pilot/pkg/config"
	"github.com/devdimensionlab/co-pilot/pkg/file"
	"os"
	"strings"
)

const tipsDir = "tips"

func LocalDir(gitCfg config.CloudConfig) string {
	return file.Path("%s/%s", gitCfg.Implementation().Dir(), tipsDir)
}

func CloudSource(name string, cloudConfig config.CloudConfig) (string, error) {
	glbConf, err := cloudConfig.GlobalCloudConfig()
	if err != nil {
		return "", err
	}
	link := glbConf.CloudConfigSource.RootUrl + glbConf.CloudConfigSource.RelativFileUrl + "/" + tipsDir + "/" + name + ".md"
	return link, nil
}

func List(gitCfg config.CloudConfig) (tips []os.DirEntry, err error) {
	items, err := os.ReadDir(LocalDir(gitCfg))
	if err != nil {
		return
	}

	for _, item := range items {
		if !item.IsDir() && strings.HasSuffix(item.Name(), ".md") {
			tips = append(tips, item)
		}
	}

	return
}
