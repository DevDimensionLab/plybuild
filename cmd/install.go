package cmd

import (
	"github.com/devdimensionlab/ply/pkg/config"
	"github.com/devdimensionlab/ply/pkg/file"
	"github.com/spf13/cobra"
)

type InstallOpts struct {
	AutoComplete bool
}

var installOpts InstallOpts

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Various install options for generating autocompletion etc",
	Long:  `Various install options for generating autocompletion etc`,
	Run: func(cmd *cobra.Command, args []string) {
		if installOpts.AutoComplete {
			generateAutoCompleteFiles(cmd)
		}
	},
}

func generateAutoCompleteFiles(cmd *cobra.Command) {
	completeDir, err := configPathFor("autocomplete")
	if err != nil {
		log.Warnln(err)
	}
	if err := file.CreateDirectory(completeDir); err != nil {
		log.Fatalln(err)
	}

	if err := generateCompleteFor("autocomplete/autocomplete.bash", cmd.Root().GenBashCompletionFile); err != nil {
		log.Warnln(err.Error())
	}
	if err := generateCompleteFor("autocomplete/autocomplete.zsh", cmd.Root().GenZshCompletionFile); err != nil {
		log.Warnln(err.Error())
	}
	if err := generateCompleteFor("autocomplete/autocomplete.pshell", cmd.Root().GenPowerShellCompletionFileWithDesc); err != nil {
		log.Warnln(err.Error())
	}

	fishFile, err := configPathFor("autocomplete/autocomplete.fish")
	if err != nil {
		log.Warnln(err.Error())
	}
	if err := cmd.Root().GenFishCompletionFile(fishFile, true); err != nil {
		log.Warnln(err.Error())
	}
}

func generateCompleteFor(relPath string, generator func(string) error) error {
	f, err := configPathFor(relPath)
	if err != nil {
		return err
	}
	return generator(f)
}

func configPathFor(fileName string) (string, error) {
	configPath, err := config.GetPlyHomePath()
	if err != nil {
		return "", err
	}
	return file.Path("%s/%s", configPath, fileName), nil
}

func init() {
	RootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolVar(&installOpts.AutoComplete, "autocomplete", false, "Generate autocomplete files")
}
