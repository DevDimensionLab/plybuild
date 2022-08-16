package context

import (
	"fmt"
	"github.com/devdimensionlab/co-pilot/pkg/config"
	"github.com/devdimensionlab/co-pilot/pkg/file"
	"github.com/devdimensionlab/co-pilot/pkg/logger"
	"github.com/devdimensionlab/co-pilot/pkg/maven"
)

type Context struct {
	Recursive       bool
	DryRun          bool
	TargetDirectory string
	DisableGit      bool
	ForceCloudSync  bool
	OpenInBrowser   bool
	Projects        []config.Project
	Err             error
	ProfilesPath    string
	LocalConfig     config.LocalConfig
	CloudConfig     config.CloudConfig
}

func (ctx *Context) FindAndPopulateMavenProjects() error {
	excludes := []string{
		"flattened-pom.xml",
		"/target/",
	}

	if ctx.Recursive {
		pomFiles, err := file.FindAll("pom.xml", excludes, ctx.TargetDirectory)
		if err != nil {
			return err
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
			return err
		}
		ctx.Projects = append(ctx.Projects, project)
	}

	return nil
}

func (ctx Context) OnEachMavenProject(description string, do ...func(repository maven.Repository, project config.Project) error) {
	if ctx.Projects == nil || len(ctx.Projects) == 0 {
		log.Errorln("could not find any pom models in the context")
		return
	}

	repo, err := maven.DefaultRepository()
	if err != nil {
		log.Warnln(err)
	}

	for _, p := range ctx.Projects {
		if p.CloudConfig != nil {
			projectDefaults, err := p.CloudConfig.ProjectDefaults()
			if err != nil {
				log.Warnf("could not find a project-defaults.json file in cloud-config")
				log.Debugf("%v", err)
			}
			p.Config.Settings.MergeProjectDefaults(projectDefaults)
		}
		if p.Type == nil {
			log.Warnf("no project type defined for path: %s", p.Path)
			continue
		}

		log.Info(logger.White(fmt.Sprintf("%s in %s", description, p.Path)))

		if p.IsDirtyGitRepo() {
			log.Debugf("operating on a dirty git repo")
		}

		if do != nil {
			for _, job := range do {
				if job == nil {
					continue
				}
				err := job(repo, p)
				if err != nil {
					log.Warnln(err)
					continue
				}
			}
		}

		if !ctx.DryRun {
			if err := p.SortAndWritePom(); err != nil {
				log.Warnln(err)
			}
		}
	}
}

func (ctx Context) OnRootProject(description string, do ...func(project config.Project) error) {
	if ctx.Projects == nil || len(ctx.Projects) == 0 {
		log.Errorln("could not find any pom models in the context")
		return
	}

	rootProject := ctx.Projects[0]
	if rootProject.Type == nil {
		log.Fatalln(fmt.Sprintf("no project type defined for path: %s", rootProject.Path))
	}
	log.Info(logger.White(fmt.Sprintf("%s for file %s", description, rootProject.Type.FilePath())))

	if rootProject.IsDirtyGitRepo() {
		log.Warnf("operating on a dirty git repo")
	}

	if do != nil {
		for _, job := range do {
			if job == nil {
				continue
			}
			err := job(rootProject)
			if err != nil {
				log.Warnln(err)
				continue
			}
		}
	}

	if !ctx.DryRun {
		if err := rootProject.SortAndWritePom(); err != nil {
			log.Warnln(err)
		}
	}
}

func (ctx *Context) LoadProfile(profilePath string) {
	ctx.LocalConfig = config.NewLocalConfig(profilePath)
	ctx.CloudConfig = config.OpenGitCloudConfig(profilePath)
	if !ctx.LocalConfig.Exists() {
		err := ctx.LocalConfig.TouchFile()
		if err != nil {
			log.Error(err)
		}
	}
}
