package git

import "os/exec"

func Clone(url string, target string) error {
	return exec.Command("git", "clone", url, target).Run()
}

func Pull(target string) error {
	return exec.Command("git", "-C", target, "pull", "origin").Run()
}
