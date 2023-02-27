// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"github.com/devdimensionlab/co-pilot/pkg/config"
	"github.com/devdimensionlab/co-pilot/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const version = "v0.6.24"

var log = logger.Context()

var RootCmd = &cobra.Command{
	Use:   "co-pilot",
	Short: "Co-pilot is a developer tool for automating common tasks on a spring boot project",
	Long:  header(),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		initConfig()
	})

	logrus.SetOutput(os.Stdout)
	RootCmd.PersistentFlags().Bool("debug", false, "turn on debug output")
	RootCmd.PersistentFlags().Bool("json", false, "turn on json output logging")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// migration step from old version without profiles
	_, err := config.GetActiveProfilePath()
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		if err := config.InstallOrMigrateToProfiles(); err != nil {
			log.Fatalln(err)
		}
	}

	ctx.ProfilesPath, err = config.GetActiveProfilePath()
	if err != nil {
		log.Fatalln(err)
	}
	ctx.LoadProfile(ctx.ProfilesPath)

	viper.SetEnvPrefix("COPILOT")
	viper.SetConfigFile(ctx.LocalConfig.FilePath())
	_ = viper.ReadInConfig()
	viper.AutomaticEnv() // read in environment variables that match
}
