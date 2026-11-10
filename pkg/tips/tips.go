package tips

import (
	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"os"
	"strings"
)

const TipsDir = "tips"

type listDependencies struct {
	Files filesystem.Dependencies
}

func systemListDependencies() listDependencies {
	return listDependencies{Files: filesystem.System()}
}

func LocalDir(gitCfg config.CloudConfig) string {
	return file.Path("%s/%s", gitCfg.Implementation().Dir(), TipsDir)
}

func List(gitCfg config.CloudConfig) (tips []os.DirEntry, err error) {
	return list(systemListDependencies(), gitCfg)
}

func list(dependencies listDependencies, gitCfg config.CloudConfig) (tips []os.DirEntry, err error) {
	items, err := filesystem.ReadDirEntries(dependencies.Files, LocalDir(gitCfg))
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
