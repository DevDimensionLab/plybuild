package maven

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"strings"
)

type Context struct {
	Recursive       bool
	Overwrite       bool
	DryRun          bool
	TargetDirectory string
	Poms            []PomWrapper
	Err             error
}

func Write(pair PomWrapper, overwrite bool) error {
	if err := SortAndWritePom(pair, overwrite); err != nil {
		return err
	}

	return nil
}

func PomFileToTargetDirectory(pomFile string) string {
	pomFilePathParts := strings.Split(pomFile, "/")
	return file.Path(strings.Join(pomFilePathParts[:len(pomFilePathParts)-1], "/"))
}

func (ctx *Context) FindAndPopulatePomModels() *Context {
	excludes := []string{
		"flattened-pom.xml",
		"/target/",
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
			projectConfigFile := file.Path(strings.Replace(pomFile, "pom.xml", "co-pilot.json", 1))
			projectConfig, err := config.InitProjectConfigurationFromFile(projectConfigFile)
			if err != nil {
				log.Warnln(err)
			}
			ctx.Poms = append(ctx.Poms, PomWrapper{
				Model:         model,
				PomFile:       pomFile,
				ProjectConfig: projectConfig,
			})
		}
	} else {
		pomFile := file.Path("%s/pom.xml", ctx.TargetDirectory)
		model, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Warnln(err)
			return ctx
		}
		projectConfigFile := file.Path(strings.Replace(pomFile, "pom.xml", "co-pilot.json", 1))
		projectConfig, err := config.InitProjectConfigurationFromFile(projectConfigFile)
		if err != nil {
			log.Warnln(err)
		}
		ctx.Poms = append(ctx.Poms, PomWrapper{
			Model:         model,
			PomFile:       pomFile,
			ProjectConfig: projectConfig,
		})
	}

	return ctx
}

func (ctx Context) OnEachPomProject(description string, do func(pomWrapper PomWrapper, args ...interface{}) error) {
	if ctx.Poms == nil || len(ctx.Poms) == 0 {
		log.Errorln("could not find any pom models in the context")
		return
	}

	for _, p := range ctx.Poms {
		log.Info(logger.White(fmt.Sprintf("%s for pom file %s", description, p.PomFile)))

		if do != nil {
			err := do(p)
			if err != nil {
				log.Warnln(err)
				continue
			}
		}

		if !ctx.DryRun {
			if err := Write(p, ctx.Overwrite); err != nil {
				log.Warnln(err)
			}
		}
	}
}
