package taskrun

import (
	"os"
	"os/exec"
)

func automaticVerificationSignals() []os.Signal { return []os.Signal{os.Interrupt} }
func configureAutomaticProcess(c *exec.Cmd)     {}
