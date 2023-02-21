package mustache_render

import (
	"github.com/devdimensionlab/co-pilot/pkg/config"
	"github.com/devdimensionlab/co-pilot/pkg/file"
)

func LocalDir(gitCfg config.CloudConfig) string {
	return file.Path("%s/%s", gitCfg.Implementation().Dir(), "mustache")
}

func MarkdownMustacheTemplateString(gitCfg config.CloudConfig, templateName string) (string, error) {
	mustacheDir := LocalDir(gitCfg)
	content, err := file.Open(mustacheDir + "/" + templateName)
	if err != nil {
		return "", err
	}
	return string(content[:]), nil
}
