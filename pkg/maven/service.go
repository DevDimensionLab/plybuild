package maven

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"fmt"
)

type Context struct {
	Recursive       bool
	Overwrite       bool
	DryRun          bool
	TargetDirectory string
	Projects        []config.Project
	Err             error
}

func Write(project config.Project, overwrite bool) error {
	if err := SortAndWritePom(project, overwrite); err != nil {
		return err
	}

	return nil
}

func (ctx *Context) FindAndPopulatePomProjects() *Context {
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
			project, err := config.InitProjectFromPomFile(pomFile)
			if err != nil {
				log.Warnln(err)
			}
			ctx.Projects = append(ctx.Projects, project)
		}
	} else {
		project, err := config.InitProjectFromDirectory(ctx.TargetDirectory)
		if err != nil {
			log.Fatalln(err)
		}
		ctx.Projects = append(ctx.Projects, project)
	}

	return ctx
}

func (ctx Context) OnEachProject(description string, do func(project config.Project, args ...interface{}) error) {
	if ctx.Projects == nil || len(ctx.Projects) == 0 {
		log.Errorln("could not find any pom models in the context")
		return
	}

	for _, p := range ctx.Projects {
		log.Info(logger.White(fmt.Sprintf("%s for pom file %s", description, p.PomFile)))

		if p.IsDirtyGitRepo() {
			log.Warnf("operating on a dirty git repo")
		}

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
