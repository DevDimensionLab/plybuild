package cmd

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var examplesCmd = &cobra.Command{
	Use:     "example",
	Short:   "Examples found in cloud-config",
	Long:    `Examples found in cloud-config`,
	Aliases: []string{"examples"},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}

var examplesInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install example from cloud-config",
	Long:  `install example from cloud-config`,
	Run: func(cmd *cobra.Command, args []string) {
		exampleName, _ := cmd.Flags().GetString("example-name")
		overrideGroupId, _ := cmd.Flags().GetString("group-id")
		overrideArtifactId, _ := cmd.Flags().GetString("artifact-id")
		overridePackage, _ := cmd.Flags().GetString("package")
		overrideName, _ := cmd.Flags().GetString("name")

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

		for _, example := range examples {
			if example == exampleName {

				path := file.Path("%s/examples/%s/ply.json", ctx.CloudConfig.Implementation().Dir(), example)

				// legacy support for older co-pilot.json files
				if !file.Exists(path) {
					path = file.Path("%s/examples/%s/co-pilot.json", ctx.CloudConfig.Implementation().Dir(), example)
				}

				cmd.Flags().String("config-file", path, "Optional config file")
				projectConfig, err := config.InitProjectConfigurationFromFile(path)
				if err != nil {
					log.Fatalln(err)
				}

				if overrideGroupId == "" {
					groupId, err := promptFor("groupId", projectConfig.GroupId, force)
					if err != nil {
						log.Fatalln(err)
					}
					if err := cmd.Flags().Set("group-id", groupId); err != nil {
						log.Fatalln(err)
					}
				}

				if overrideArtifactId == "" {
					artifactId, err := promptFor("artifactId", projectConfig.ArtifactId, force)
					if err != nil {
						log.Fatalln(err)
					}
					if err := cmd.Flags().Set("artifact-id", artifactId); err != nil {
						log.Fatalln(err)
					}
				}

				if overridePackage == "" {
					packageName, err := promptFor("package", projectConfig.Package, force)
					if err != nil {
						log.Fatalln(err)
					}
					if err := cmd.Flags().Set("package", packageName); err != nil {
						log.Fatalln(err)
					}
				}

				if overrideName == "" {
					name, err := promptFor("name", projectConfig.Name, force)
					if err != nil {
						log.Fatalln(err)
					}
					if err := cmd.Flags().Set("name", name); err != nil {
						log.Fatalln(err)
					}
					if err := cmd.Flags().Set("application-name", fmt.Sprintf("%sApplication", name)); err != nil {
						log.Fatalln(err)
					}
				}

				buildCmd.Run(cmd, args)
				return
			}
		}

		log.Fatalf("could not find %s in examples", exampleName)
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

	examplesCmd.AddCommand(examplesInstallCmd)
	examplesInstallCmd.Flags().StringP("example-name", "n", "", "Example name")

	examplesInstallCmd.Flags().String("boot-version", "", "Defines spring-boot version to use")
	examplesInstallCmd.Flags().String("group-id", "", "Overrides groupId from config file")
	examplesInstallCmd.Flags().String("artifact-id", "", "Overrides artifactId from config file")
	examplesInstallCmd.Flags().String("package", "", "Overrides package from config file")
	examplesInstallCmd.Flags().String("name", "", "Overrides name from config file")
	examplesInstallCmd.Flags().String("application-name", "", "Overrides application-name from config file")
}
