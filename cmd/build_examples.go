package cmd

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var examplesCmd = &cobra.Command{
	Use:   "example",
	Short: "Builds example from cloud-config",
	Long:  `Builds example example from cloud-config`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		exampleName, _ := cmd.Flags().GetString("name")
		if exampleName == "" {
			log.Fatalln("please enter example --name")
		}

		// sync cloud config
		if err := ctx.CloudConfig.Refresh(ctx.LocalConfig); err != nil {
			log.Fatalln(err)
		}

		// force defaults
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			log.Fatalln(err)
		}

		examples, err := ctx.CloudConfig.Examples()
		if err != nil {
			log.Warnln(err)
			return
		}

		var foundExample = false
		for _, example := range examples {
			if example == exampleName {
				foundExample = true
			}
		}
		if !foundExample {
			log.Fatalf("could not find %s in examples", exampleName)
		}

		jsonConfigFile := file.Path("%s/examples/%s/ply.json", ctx.CloudConfig.Implementation().Dir(), exampleName)

		// legacy support for older co-pilot.json files
		if !file.Exists(jsonConfigFile) {
			jsonConfigFile = file.Path("%s/examples/%s/co-pilot.json", ctx.CloudConfig.Implementation().Dir(), exampleName)
		}

		orderConfig, err := config.InitProjectConfigurationFromFile(jsonConfigFile)

		cmd.Flags().String("config-file", jsonConfigFile, "Optional config file")
		projectConfig, err := config.InitProjectConfigurationFromFile(jsonConfigFile)
		if err != nil {
			log.Fatalln(err)
		}

		groupId, err := promptFor("groupId", projectConfig.GroupId, force)
		if err != nil {
			log.Fatalln(err)
		}
		orderConfig.GroupId = groupId

		artifactId, err := promptFor("artifactId", projectConfig.ArtifactId, force)
		if err != nil {
			log.Fatalln(err)
		}
		orderConfig.ArtifactId = artifactId

		packageName, err := promptFor("package", projectConfig.Package, force)
		if err != nil {
			log.Fatalln(err)
		}
		orderConfig.Package = packageName

		applicationName, err := promptFor("application-name", projectConfig.Name, force)
		if err != nil {
			log.Fatalln(err)
		}
		orderConfig.ApplicationName = applicationName

		build(orderConfig, "", false)
		return
	},
}

func promptFor(value, defaultValue string, force bool) (string, error) {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Enter %s: [%s]", value, defaultValue),
		Templates: templates,
	}

	if force {
		return defaultValue, nil
	}
	newValue, err := prompt.Run()
	if err != nil {
		return "", err
	}
	if newValue == "" {
		return defaultValue, err
	}
	return newValue, nil
}

func init() {
	buildCmd.AddCommand(examplesCmd)

	examplesCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
	examplesCmd.PersistentFlags().BoolVar(&ctx.DryRun, "dry-run", false, "dry run does not write to pom.xml")
	examplesCmd.Flags().String("name", "", "Example name to use")

	examplesCmd.Flags().String("boot-version", "", "Defines spring-boot version to use")
	//examplesCmd.Flags().String("group-id", "", "Overrides groupId from config file")
	//examplesCmd.Flags().String("artifact-id", "", "Overrides artifactId from config file")
	//examplesCmd.Flags().String("package", "", "Overrides package from config file")
	//examplesCmd.Flags().String("application-name", "", "Overrides application-name from config file")
}
