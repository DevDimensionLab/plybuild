package cmd

import (
	"co-pilot/pkg/upgrade"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"strings"
)

type CommonArgs struct {
	Recursive       bool
	Overwrite       bool
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
	log.Infof("sorting and writing to pom file %s", filename)
	if err := upgrade.SortAndWrite(model, writeToFile); err != nil {
		return err
	}

	return nil
}

func pomFileToTargetDirectory(pomFile string) string {
	pomFilePathParts := strings.Split(pomFile, "/")
	return strings.Join(pomFilePathParts[:len(pomFilePathParts)-1], "/")
}
