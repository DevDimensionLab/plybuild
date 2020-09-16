package service

import (
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/maven"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"strings"
)

type Context struct {
	Recursive       bool
	Overwrite       bool
	DryRun          bool
	TargetDirectory string
	PomModels       map[string]*pom.Model
	Err             error
}

var log = logger.Context()

func Write(overwrite bool, filename string, model *pom.Model) error {
	var writeToFile = filename
	if !overwrite {
		writeToFile = filename + ".new"
	}
	if err := maven.SortAndWritePom(model, writeToFile); err != nil {
		return err
	}

	return nil
}

func PomFileToTargetDirectory(pomFile string) string {
	pomFilePathParts := strings.Split(pomFile, "/")
	return strings.Join(pomFilePathParts[:len(pomFilePathParts)-1], "/")
}

func (ctx *Context) FindAndPopulatePomModels() {
	excludes := []string{
		"flattened-pom.xml",
		"/target/",
	}

	if ctx.PomModels == nil {
		ctx.PomModels = make(map[string]*pom.Model)
	}

	if ctx.Recursive {
		pomFiles, err := file.FindAll("pom.xml", excludes, ctx.TargetDirectory)
		if err != nil {
			log.Fatalln(err)
		}
		for _, pomFile := range pomFiles {
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Warnln(err)
				continue
			}
			ctx.PomModels[pomFile] = model
		}
	} else {
		pomFile := fmt.Sprintf("%s/pom.xml", ctx.TargetDirectory)
		model, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Warnln(err)
			return
		}
		ctx.PomModels[pomFile] = model
	}
}

func (ctx Context) OnEachPomProject(description string, do func(model *pom.Model, args ...interface{}) error) {
	if ctx.PomModels == nil {
		log.Errorln("could not find any pom models in the context")
		return
	}

	for pomFile, model := range ctx.PomModels {
		log.Info(logger.White(fmt.Sprintf("%s for pom file %s", description, pomFile)))

		if err := do(model); err != nil {
			log.Warnln(err)
			continue
		}

		if !ctx.DryRun {
			if err := Write(ctx.Overwrite, pomFile, model); err != nil {
				log.Warnln(err)
			}
		}
	}
}
