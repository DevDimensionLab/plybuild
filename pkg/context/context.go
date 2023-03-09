package context

import (
	"fmt"
	"github.com/devdimensionlab/co-pilot/pkg/config"
	"github.com/devdimensionlab/co-pilot/pkg/file"
	"github.com/devdimensionlab/co-pilot/pkg/maven"
)

type Context struct {
	Recursive       bool
	DryRun          bool
	TargetDirectory string
	DisableGit      bool
	ForceCloudSync  bool
	OpenInBrowser   bool
	StealthMode     bool
	Projects        []config.Project
	Err             error
	ProfilesPath    string
	LocalConfig     config.LocalConfigDir
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

func (ctx *Context) OnEachMavenProject(description string, do ...func(repository maven.Repository, project config.Project) error) {
	if ctx.Projects == nil || len(ctx.Projects) == 0 {
		log.Errorln("could not find any pom models in the context")
		return
	}

	mavenRepository := ctx.GetMavenRepository()

	for _, p := range ctx.Projects {
		if ctx.CloudConfig != nil {
			projectDefaults, err := ctx.CloudConfig.ProjectDefaults()
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

		if ctx.StealthMode {
			p.Config.Settings.UseStealthMode = true
		}

		log.Info(fmt.Sprintf("%s in %s", description, p.Path))

		if p.IsDirtyGitRepo() {
			log.Debugf("operating on a dirty git repo")
		}

		if do != nil {
			for _, job := range do {
				if job == nil {
					continue
				}
				err := job(mavenRepository, p)
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

func (ctx *Context) OnRootProject(description string, do ...func(project config.Project) error) {
	if ctx.Projects == nil || len(ctx.Projects) == 0 {
		log.Errorln("could not find any pom models in the context")
		return
	}

	rootProject := ctx.Projects[0]
	if rootProject.Type == nil {
		log.Fatalln(fmt.Sprintf("no project type defined for path: %s", rootProject.Path))
	}
	log.Info(fmt.Sprintf("%s for file %s", description, rootProject.Type.FilePath()))

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
	ctx.LocalConfig = config.OpenLocalConfig(profilePath)
	ctx.CloudConfig = config.OpenGitCloudConfig(profilePath)
	if !ctx.LocalConfig.Exists() {
		log.Debugf("localConfig does not exists, touching a new file")
		err := ctx.LocalConfig.TouchFile()
		if err != nil {
			log.Error(err)
		}
	}
}

func (ctx *Context) GetMavenRepository() maven.Repository {
	cfg, err := ctx.LocalConfig.Config()
	if err != nil {
		log.Warnln(err)
	}

	var repository = maven.Repository{}

	if cfg.Nexus.Url != "" {
		log.Debugf("using maven repository from local config %s\n", cfg.Nexus.Url)
		repository = maven.RepositoryFrom(cfg.Nexus.Url, cfg.Nexus.Username, cfg.Nexus.Password)
	} else {
		log.Debugf("search for maven repository in .m2 folder \n")
		repository, err = maven.DefaultRepository()
		if err != nil {
			log.Warnln(err)
		}
	}

	return repository
}
