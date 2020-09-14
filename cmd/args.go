package cmd

import (
	"co-pilot/pkg/file"
	"co-pilot/pkg/upgrade"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

type CommonArgs struct {
	Recursive       bool
	Overwrite       bool
	DryRun          bool
	TargetDirectory string
	PomFiles        []string
	Err             error
}

var cArgs CommonArgs

func write(filename string, model *pom.Model) error {
	var writeToFile = filename
	if !cArgs.Overwrite {
		writeToFile = filename + ".new"
	}
	if err := upgrade.SortAndWrite(model, writeToFile); err != nil {
		return err
	}

	return nil
}

func pomFileToTargetDirectory(pomFile string) string {
	pomFilePathParts := strings.Split(pomFile, "/")
	return strings.Join(pomFilePathParts[:len(pomFilePathParts)-1], "/")
}

func populatePomFiles() {
	excludes := []string{
		"flattened-pom.xml",
		"/target/",
	}

	if cArgs.Recursive {
		if cArgs.PomFiles, cArgs.Err = file.FindAll("pom.xml", excludes, cArgs.TargetDirectory); cArgs.Err != nil {
			log.Fatalln(cArgs.Err)
		}
	} else {
		cArgs.PomFiles = append(cArgs.PomFiles, cArgs.TargetDirectory+"/pom.xml")
	}
}

func EnableDebug(cmd *cobra.Command, args []string) error {
	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return err
	}
	if debug {
		fmt.Println("== debug mode enabled ==")
		logrus.SetLevel(logrus.DebugLevel)
	}

	return nil
}
