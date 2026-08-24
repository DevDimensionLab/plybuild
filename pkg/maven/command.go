package maven

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/logger"
	"strings"
)

const versionsPlugin = "org.codehaus.mojo:versions-maven-plugin:2.8.1"

func RunOn(cmd string, args ...string) func(repository Repository, project config.Project) error {
	return runOn(process.System(), cmd, args...)
}

func runOn(dependencies process.Dependencies, cmd string, args ...string) func(repository Repository, project config.Project) error {
	return func(repository Repository, project config.Project) error {
		log.Infof("running: [%s] => %s %s", project.Path, cmd, strings.Join(args, " "))
		return process.Execute(dependencies, process.Command{
			Name:   cmd,
			Args:   args,
			Dir:    project.Path,
			Stdout: logger.StdOut(),
		})
	}
}

func UpdateProperty(property, version string) []string {
	return []string{
		fmt.Sprintf("%s:update-property", versionsPlugin),
		fmt.Sprintf("-Dproperty=%s", property),
		fmt.Sprintf("-DnewVersion=[%s]", version),
		"-DallowDowngrade=true",
	}
}

func UseLatestVersion(groupId, artifactId string) []string {
	return []string{
		fmt.Sprintf("%s:use-latest-versions", versionsPlugin),
		fmt.Sprintf("-Dincludes=%s:%s", groupId, artifactId),
	}
}
