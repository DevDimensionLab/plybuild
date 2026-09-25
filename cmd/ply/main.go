package main

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/cmd"
	"os"
)

var execute = cmd.ExecuteE

func main() {
	if err := execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
