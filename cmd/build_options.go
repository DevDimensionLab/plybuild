package cmd

import (
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/devdimensionlab/plybuild/pkg/maven"
	"github.com/devdimensionlab/plybuild/pkg/spring"
	"github.com/devdimensionlab/plybuild/pkg/template"
	"github.com/spf13/cobra"
)

type InfoOpts struct {
	SpringManaged     bool
	SpringInfo        bool
	MavenRepositories bool
	Templates         bool
	Examples          bool
}

var infoOpts InfoOpts

func (infoOpts InfoOpts) Any() bool {
	return infoOpts.SpringManaged ||
		infoOpts.SpringInfo ||
		infoOpts.MavenRepositories ||
		infoOpts.Templates ||
		infoOpts.Examples
}

var optionsCmd = &cobra.Command{
	Use:   "options",
	Short: "Prints options on spring version, dependencies etc",
	Long:  `Prints options on spring version, dependencies etc`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := InitGlobals(cmd); err != nil {
			return err
		}
		return OkHelp(cmd, infoOpts.Any)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if infoOpts.SpringInfo {
			if err := springInfo(); err != nil {
				return err
			}
		}
		if infoOpts.SpringManaged {
			if err := showSpringManaged(); err != nil {
				return err
			}
		}
		if infoOpts.MavenRepositories {
			if err := showMavenRepositories(); err != nil {
				return err
			}
		}
		if infoOpts.Templates {
			if err := showTemplates(); err != nil {
				return err
			}
		}
		if infoOpts.Examples {
			if err := showExamples(); err != nil {
				return err
			}
		}
		return nil
	},
}

func springInfo() error {
	repo, err := maven.DefaultRepository()
	if err != nil {
		return err
	}

	latestVersionMeta, err := repo.GetMetaData("org.springframework.boot", "spring-boot")
	if err != nil {
		return err
	}

	latestVersion, err := latestVersionMeta.LatestRelease()
	if err != nil {
		return err
	}

	root, err := spring.GetRoot()
	if err != nil {
		return err
	}
	log.Infof("Latest version of spring boot are: %s\n", latestVersion)

	log.Info("Valid dependencies: ")
	for _, category := range root.Dependencies.Values {
		fmt.Println(category.Name)
		fmt.Printf("================================\n")
		for _, dep := range category.Values {
			fmt.Printf("[%s]\n    %s, (%s)\n", dep.Id, dep.Name, dep.Description)
		}
		fmt.Printf("\n")
	}
	return nil
}

func showSpringManaged() error {
	deps, err := spring.GetDependencies()
	if err != nil {
		return err
	}

	log.Info("Spring Boot managed dependencies:")
	var organized = make(map[string][]pom.Dependency)
	for _, dep := range deps.Dependencies {
		mvnDep := pom.Dependency{
			GroupId:    dep.GroupId,
			ArtifactId: dep.ArtifactId,
		}
		organized[dep.GroupId] = append(organized[dep.GroupId], mvnDep)
	}

	for k, v := range organized {
		fmt.Printf("GroupId: %s\n", k)
		fmt.Printf("================================\n")
		for _, mvnDep := range v {
			fmt.Printf("  ArtifactId: %s\n", mvnDep.ArtifactId)
		}
	}
	return nil
}

func showMavenRepositories() error {
	settings, _ := maven.NewSettings()
	return settings.ListRepositories()
}

func showTemplates() error {
	markdownFormat := false
	templates, err := ctx.CloudConfig.Templates()
	if err != nil {
		return err
	}
	terminalConfig, err := ctx.LocalConfig.GetTerminalConfig()
	if err != nil {
		return err
	}

	if markdownFormat || terminalConfig.Format == "markdown" {
		markdownDocument, err := template.ListAsMarkdown(ctx.CloudConfig, templates)
		if err != nil {
			return err
		}

		markdownForTerminal := markdown.Render(markdownDocument, terminalConfig.Width, 2)
		fmt.Println("\n" + string(markdownForTerminal))

		gCloudCfg, err := ctx.CloudConfig.GlobalCloudConfig()
		if err != nil {
			return err
		}
		cloudSource := gCloudCfg.SourceFor(template.TemplatesDir, "README.md")
		log.Infoln("Cloud source: " + cloudSource)
	} else {
		for _, folder := range templates {
			log.Infof("%s - %s (%s)", folder.Name, folder.Project.Config.Description, folder.Project.Config.Language)
		}
	}
	return nil
}

func showExamples() error {
	// sync cloud config
	if err := ctx.CloudConfig.Refresh(ctx.LocalConfig); err != nil {
		return err
	}

	examples, err := ctx.CloudConfig.Examples()
	if err != nil {
		return err
	}

	fmt.Println("Available examples are:")
	for _, example := range examples {
		fmt.Printf("\t* %s\n", example)
	}
	return nil
}

func init() {
	buildCmd.AddCommand(optionsCmd)
	optionsCmd.PersistentFlags().BoolVar(&infoOpts.SpringInfo, "spring-dependencies", false, "show spring boot status")
	optionsCmd.PersistentFlags().BoolVar(&infoOpts.SpringManaged, "spring-managed", false, "show spring boot managed dependencies info")
	optionsCmd.PersistentFlags().BoolVar(&infoOpts.MavenRepositories, "maven-repositories", false, "show current maven repositories")
	optionsCmd.PersistentFlags().BoolVar(&infoOpts.Templates, "templates", false, "show plybuild templates")
	optionsCmd.PersistentFlags().BoolVar(&infoOpts.Examples, "examples", false, "show plybuild examples")

	//optionsCmd.Flags().Bool("markdown", false, "Outputs templates as markdown in the terminal")
	//optionsCmd.Flags().Bool("save", false, "Saves the template markdown doc to cloud-config template-folder")
}
