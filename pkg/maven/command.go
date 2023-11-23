package maven

import (
	"fmt"
	"github.com/devdimensionlab/ply/pkg/config"
	"github.com/devdimensionlab/ply/pkg/logger"
	"os/exec"
	"strings"
)

const versionsPlugin = "org.codehaus.mojo:versions-maven-plugin:2.8.1"

func RunOn(cmd string, args ...string) func(repository Repository, project config.Project) error {
	return func(repository Repository, project config.Project) error {
		log.Infof("running: [%s] => %s %s", project.Path, cmd, strings.Join(args, " "))
		cmd := exec.Command(cmd, args...)
		cmd.Dir = project.Path
		cmd.Stdout = logger.StdOut()
		return cmd.Run()
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
