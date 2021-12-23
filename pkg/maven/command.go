package maven

import (
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"io"
	"os/exec"
	"strings"
)

func RunOn(stdout io.Writer, cmd string, args ...string) func(project config.Project) error {
	return func(project config.Project) error {
		log.Infof("running: [%s] => %s %s", project.Path, cmd, strings.Join(args, " "))
		cmd := exec.Command(cmd, args...)
		cmd.Dir = project.Path
		cmd.Stdout = stdout
		return cmd.Run()
	}
}
