package cmd

import (
	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/devdimensionlab/plybuild/pkg/maven"
	"github.com/devdimensionlab/plybuild/pkg/template"
	"github.com/spf13/cobra"
	"os"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add templates and functionalities for files to a build",
	Long:  `Add templates and functionalities for files to a build`,
}

var addPomCmd = &cobra.Command{
	Use:   "pom",
	Short: "Adds a pom-file into a project",
	Long:  `Adds a pom-file into a project`,
	Run: func(cmd *cobra.Command, args []string) {
		fromPomFile, err := cmd.Flags().GetString("from")
		if err != nil {
			log.Fatalln(err)
		}
		if fromPomFile == "" {
			log.Errorln("missing valid --from flag for pom.xml to add from")
			os.Exit(-1)
		}

		importModel, err := pom.GetModelFrom(fromPomFile)
		if err != nil {
			log.Fatalln(err)
		}

		targetProject, err := config.InitProjectFromDirectory(ctx.TargetDirectory)
		if err != nil {
			log.Fatalln(err)
		}

		if err = maven.MergePoms(importModel, targetProject.Type.Model()); err != nil {
			log.Fatalln(err)
		}

		if err = targetProject.SortAndWritePom(); err != nil {
			log.Fatalln(err)
		}
	},
}

var addTextCmd = &cobra.Command{
	Use:   "text",
	Short: "Merges two text files",
	Long:  `Merges two text files`,
	Run: func(cmd *cobra.Command, args []string) {
		fromFile, err := cmd.Flags().GetString("from")
		if err != nil {
			log.Fatalln(err)
		}
		if fromFile == "" {
			log.Errorln("missing valid --from file flag")
			os.Exit(-1)
		}

		toFile, err := cmd.Flags().GetString("to")
		if err != nil {
			log.Fatalln(err)
		}
		if toFile == "" {
			log.Errorln("missing valid --to file flag")
			os.Exit(-1)
		}

		if err := file.MergeTextFiles(fromFile, toFile); err != nil {
			log.Fatalln(err)
		}
	},
}

var addTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Adds a template from ply-config",
	Long:  `Adds a template from ply-config`,
	Run: func(cmd *cobra.Command, args []string) {
		templateName, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatalln(err)
		}
		if templateName == "" {
			log.Fatalln("Missing template --name")
		}

		project, err := config.InitProjectFromDirectory(ctx.TargetDirectory)
		if err != nil {
			log.Fatalln(err)
		}

		cloudTemplate, err := project.CloudConfig.Template(templateName)
		if err != nil {
			log.Fatalln(err)
		}

		if err := template.MergeTemplate(cloudTemplate, project, false); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	buildCmd.AddCommand(addCmd)
	addCmd.AddCommand(addPomCmd)
	addCmd.AddCommand(addTextCmd)
	addCmd.AddCommand(addTemplateCmd)
	addCmd.PersistentFlags().String("from", "", "file to add")
	addPomCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
	addTextCmd.PersistentFlags().String("to", "", "target file to add to")
	addTemplateCmd.Flags().String("name", "", "template to add")
	addTemplateCmd.Flags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
}
