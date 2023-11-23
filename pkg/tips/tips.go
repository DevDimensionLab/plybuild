package tips

import (
	"github.com/devdimensionlab/ply/pkg/config"
	"github.com/devdimensionlab/ply/pkg/file"
	"os"
	"strings"
)

const TipsDir = "tips"

func LocalDir(gitCfg config.CloudConfig) string {
	return file.Path("%s/%s", gitCfg.Implementation().Dir(), TipsDir)
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
