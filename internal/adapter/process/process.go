package process

import (
	"io"
	"os/exec"
)

// Command is the complete process request passed to a dependency.
type Command struct {
	Name   string
	Args   []string
	Dir    string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Start  bool
}

// Runner performs one complete process request.
type Runner interface {
	Run(Command) error
}

// Dependencies contains the process effect used by a caller. Its zero value is
// a safe no-op; production callers must select System explicitly.
type Dependencies struct {
	Runner Runner
}

// Execute passes the complete command to the configured dependency.
func Execute(dependencies Dependencies, command Command) error {
	if dependencies.Runner == nil {
		return nil
	}
	return dependencies.Runner.Run(command)
}

// System returns the production dependency that executes an operating-system
// process.
func System() Dependencies {
	return Dependencies{Runner: systemRunner{}}
}

type systemRunner struct{}

func (systemRunner) Run(command Command) error {
	cmd := exec.Command(command.Name, command.Args...)
	cmd.Dir = command.Dir
	cmd.Stdin = command.Stdin
	cmd.Stdout = command.Stdout
	cmd.Stderr = command.Stderr
	if command.Start {
		return cmd.Start()
	}
	return cmd.Run()
}
