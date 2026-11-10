package cmd

import (
	"errors"
	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/devdimensionlab/plybuild/pkg/maven"
	"github.com/devdimensionlab/plybuild/pkg/template"
	"github.com/spf13/cobra"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		fromPomFile, err := cmd.Flags().GetString("from")
		if err != nil {
			return err
		}
		if fromPomFile == "" {
			return errors.New("missing valid --from flag for pom.xml to add from")
		}

		importModel, err := pom.GetModelFrom(fromPomFile)
		if err != nil {
			return err
		}

		targetProject, err := config.InitProjectFromDirectory(ctx.TargetDirectory)
		if err != nil {
			return err
		}

		if err = maven.MergePoms(importModel, targetProject.Type.Model()); err != nil {
			return err
		}

		if err = targetProject.SortAndWritePom(); err != nil {
			return err
		}
		return nil
	},
}

var addTextCmd = &cobra.Command{
	Use:   "text",
	Short: "Merges two text files",
	Long:  `Merges two text files`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fromFile, err := cmd.Flags().GetString("from")
		if err != nil {
			return err
		}
		if fromFile == "" {
			return errors.New("missing valid --from file flag")
		}

		toFile, err := cmd.Flags().GetString("to")
		if err != nil {
			return err
		}
		if toFile == "" {
			return errors.New("missing valid --to file flag")
		}

		if err := file.MergeTextFiles(fromFile, toFile); err != nil {
			return err
		}
		return nil
	},
}

var addTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Adds a template from ply-config",
	Long:  `Adds a template from ply-config`,
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		if templateName == "" {
			return errors.New("missing template --name")
		}

		project, err := config.InitProjectFromDirectory(ctx.TargetDirectory)
		if err != nil {
			return err
		}

		cloudTemplate, err := project.CloudConfig.Template(templateName)
		if err != nil {
			return err
		}

		if err := template.MergeTemplate(cloudTemplate, project, false); err != nil {
			return err
		}
		return nil
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
