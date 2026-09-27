package main

import (
	"github.com/devdimensionlab/plybuild/cmd"
	"os"
)

var execute = cmd.ExecuteE

func main() {
	if err := execute(); err != nil {
		os.Exit(1)
	}
}
