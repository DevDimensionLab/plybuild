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
			log.Fatalln(err)
		}
		return OkHelp(cmd, infoOpts.Any)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if infoOpts.SpringInfo {
			springInfo()
		}
		if infoOpts.SpringManaged {
			showSpringManaged()
		}
		if infoOpts.MavenRepositories {
			showMavenRepositories()
		}
		if infoOpts.Templates {
			showTemplates()
		}
		if infoOpts.Examples {
			showExamples()
		}
	},
}

func springInfo() {
	repo, err := maven.DefaultRepository()
	if err != nil {
		log.Fatalln(err)
	}

	latestVersionMeta, err := repo.GetMetaData("org.springframework.boot", "spring-boot")
	if err != nil {
		log.Fatalln(err)
	}

	latestVersion, err := latestVersionMeta.LatestRelease()
	if err != nil {
		log.Fatalln(err)
	}

	root, err := spring.GetRoot()
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("Latest version of spring boot are: %s\n", latestVersion)

	log.Infof(fmt.Sprintf("Valid dependencies: "))
	for _, category := range root.Dependencies.Values {
		fmt.Println(fmt.Sprintf("%s", category.Name))
		fmt.Printf("================================\n")
		for _, dep := range category.Values {
			fmt.Printf("[%s]\n    %s, (%s)\n", dep.Id, dep.Name, dep.Description)
		}
		fmt.Printf("\n")
	}
}

func showSpringManaged() {
	deps, err := spring.GetDependencies()
	if err != nil {
		log.Fatalln(err)
	}

	log.Infof(fmt.Sprintf("Spring Boot managed dependencies:"))
	var organized = make(map[string][]pom.Dependency)
	for _, dep := range deps.Dependencies {
		mvnDep := pom.Dependency{
			GroupId:    dep.GroupId,
			ArtifactId: dep.ArtifactId,
		}
		organized[dep.GroupId] = append(organized[dep.GroupId], mvnDep)
	}

	for k, v := range organized {
		fmt.Println(fmt.Sprintf("GroupId: %s", k))
		fmt.Printf("================================\n")
		for _, mvnDep := range v {
			fmt.Printf("  ArtifactId: %s\n", mvnDep.ArtifactId)
		}
	}
}

func showMavenRepositories() {
	settings, _ := maven.NewSettings()
	if err := settings.ListRepositories(); err != nil {
		log.Fatalln(err)
	}
}

func showTemplates() {
	markdownFormat := false
	templates, err := ctx.CloudConfig.Templates()
	if err != nil {
		log.Fatalln(err)
	}
	terminalConfig, err := ctx.LocalConfig.GetTerminalConfig()
	if err != nil {
		log.Fatalln(err)
	}

	if markdownFormat || terminalConfig.Format == "markdown" {
		markdownDocument, err := template.ListAsMarkdown(ctx.CloudConfig, templates)
		if err != nil {
			log.Fatalln(err)
		}

		markdownForTerminal := markdown.Render(markdownDocument, terminalConfig.Width, 2)
		fmt.Println("\n" + string(markdownForTerminal))

		gCloudCfg, err := ctx.CloudConfig.GlobalCloudConfig()
		if err != nil {
			log.Fatalln(err)
		}
		cloudSource := gCloudCfg.SourceFor(template.TemplatesDir, "README.md")
		if err != nil {
			log.Fatalln(err)
		}
		log.Infoln("Cloud source: " + cloudSource)
	} else {
		for _, folder := range templates {
			log.Infof("%s - %s (%s)", folder.Name, folder.Project.Config.Description, folder.Project.Config.Language)
		}
	}
}

func showExamples() {
	// sync cloud config
	if err := ctx.CloudConfig.Refresh(ctx.LocalConfig); err != nil {
		log.Fatalln(err)
	}

	examples, err := ctx.CloudConfig.Examples()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Available examples are:")
	for _, example := range examples {
		fmt.Printf("\t* %s\n", example)
	}
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
