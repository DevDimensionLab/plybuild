package cmd

import (
	"co-pilot/pkg/file"
	"co-pilot/pkg/upgrade"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
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
