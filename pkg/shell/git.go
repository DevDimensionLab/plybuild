package shell

import (
	"os/exec"
)

func GitClone(url string, target string) (string, error) {
	return run(exec.Command("git", "clone", url, target))
}

func GitPull(target string) (string, error) {
	return run(exec.Command("git", "-C", target, "pull", "origin"))
}
