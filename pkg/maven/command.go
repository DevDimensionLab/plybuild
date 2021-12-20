package maven

import (
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"os/exec"
	"strings"
)

func RunOn(cmd string, args ...string) func(project config.Project) error {
	return func(project config.Project) error {
		log.Debugf("running: %s => %s %s", project.Path, cmd, strings.Join(args, " "))
		cmd := exec.Command(cmd, args...)
		cmd.Dir = project.Path
		return cmd.Run()
	}
}
