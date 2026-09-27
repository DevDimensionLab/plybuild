// Copyright © 2019 plybuild.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/devdimensionlab/plybuild/pkg/bitbucket"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/context"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/devdimensionlab/plybuild/pkg/http"
	"github.com/devdimensionlab/plybuild/pkg/logger"
	"github.com/devdimensionlab/plybuild/pkg/maven"
	"github.com/devdimensionlab/plybuild/pkg/shell"
	"github.com/devdimensionlab/plybuild/pkg/spring"
	"github.com/devdimensionlab/plybuild/pkg/template"
	"github.com/devdimensionlab/plybuild/pkg/webservice"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const version = "v1.0.1"

var log = logger.Context()

var initializationErr error

var errCommandComplete = errors.New("command complete")

var templates = &promptui.PromptTemplates{
	Prompt:  "{{ . }} ",
	Valid:   "{{ . | green }} ",
	Invalid: "{{ . | red }} ",
	Success: "{{ . | bold }} ",
}

var RootCmd = &cobra.Command{
	Use:   "ply",
	Short: "Ply is a developer tool for automating common tasks on a spring boot project",
	Long:  header(),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return InitGlobals(cmd)
	},
}

// ExecuteE runs the command tree and returns any command error to the executable
// boundary. Cobra's legacy diagnostic text is preserved without letting Cobra
// select the process exit behavior.
func ExecuteE() error {
	silenceErrors := RootCmd.SilenceErrors
	silenceUsage := RootCmd.SilenceUsage
	RootCmd.SilenceErrors = true
	RootCmd.SilenceUsage = true
	executed, err := RootCmd.ExecuteC()
	RootCmd.SilenceErrors = silenceErrors
	RootCmd.SilenceUsage = silenceUsage

	if errors.Is(err, errCommandComplete) {
		return nil
	}
	if err != nil {
		RootCmd.PrintErrln("Error:", err.Error())
		if strings.HasPrefix(err.Error(), "unknown command") {
			RootCmd.PrintErrf("Run '%v --help' for usage.\n", executed.CommandPath())
		}
	}
	return err
}

// Execute preserves the original public API while leaving process termination
// to executable main packages.
func Execute() {
	if err := ExecuteE(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	logrus.SetOutput(os.Stdout)
	RootCmd.PersistentFlags().Bool("debug", false, "turn on debug output")
	RootCmd.PersistentFlags().Bool("json", false, "turn on json output logging")
	RootCmd.PersistentFlags().Bool("force", false, "uses default for prompts")
	RootCmd.PersistentFlags().Bool("doc", false, "open documentation website")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	initializationErr = initializeConfig()
}

func initializeConfig() error {
	// migration step from old version without profiles
	_, err := config.GetActiveProfilePath()
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		if err := config.InstallOrMigrateToProfiles(); err != nil {
			return err
		}
	}

	ctx.ProfilesPath, err = config.GetActiveProfilePath()
	if err != nil {
		return err
	}
	ctx.LoadProfile(ctx.ProfilesPath)

	viper.SetEnvPrefix("PLY")
	viper.SetConfigFile(ctx.LocalConfig.FilePath())
	_ = viper.ReadInConfig()
	viper.AutomaticEnv() // read in environment variables that match
	return nil
}

var ctx context.Context

func InitGlobals(cmd *cobra.Command) error {
	initConfig()
	if initializationErr != nil {
		return initializationErr
	}
	json, err := cmd.Flags().GetBool("json")
	if err != nil {
		return err
	}
	if json {
		logger.SetFieldLogger()
		logger.SetJsonLogging()
	}

	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return err
	}
	if debug {
		debugLogger := logger.DebugLogger()
		log = debugLogger
		bitbucket.SetLogger(debugLogger)
		config.SetLogger(debugLogger)
		context.SetLogger(debugLogger)
		file.SetLogger(debugLogger)
		http.SetLogger(debugLogger)
		maven.SetLogger(debugLogger)
		shell.SetLogger(debugLogger)
		spring.SetLogger(debugLogger)
		template.SetLogger(debugLogger)
	}

	stealth, _ := cmd.Flags().GetBool("stealth")
	if stealth {
		logger.SetFieldLogger()
	}

	return nil
}

func OkHelp(cmd *cobra.Command, depend func() bool) error {
	if !cmd.Flags().HasFlags() || !depend() {
		return flag.ErrHelp
	}
	return nil
}

func SyncActiveProfileCloudConfig() error {
	if ctx.ForceCloudSync {
		if err := ctx.CloudConfig.Refresh(ctx.LocalConfig); err != nil {
			return errors.New("failed to sync cloud config: " + err.Error())
		}
	}
	return nil
}

func OpenDocumentationWebsite(cmd *cobra.Command, path string) error {
	doc, err := cmd.Flags().GetBool("doc")
	if err != nil {
		return err
	}
	if doc {
		err = webservice.OpenBrowser("https://plybuild.io/" + path)
		if err != nil {
			return err
		}
		return errCommandComplete
	}
	return nil
}

func promptForValue(value, defaultValue string, force bool) (string, error) {
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

func promptForContinue(question string, force bool) (bool, error) {
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("%s, continue? [Y/n]", question),
		//Templates: templates,
		Validate: func(input string) error {
			if len(input) > 0 && (input != "y" && input != "Y" && input != "n") {
				return errors.New("please enter 'Y/y' or 'n'")
			}
			return nil
		},
	}

	if force {
		return true, nil
	}

	answer, err := prompt.Run()
	if err != nil {
		return false, err
	}

	if answer == "y" || answer == "Y" || answer == "" {
		return true, nil
	}

	return false, nil
}
