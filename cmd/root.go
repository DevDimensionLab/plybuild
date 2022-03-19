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
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"github.com/co-pilot-cli/co-pilot/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var version = "v0.5.3"
var log = logger.Context()
var activeLocalConfig config.LocalConfig
var activeCloudConfig config.GitCloudConfig

var RootCmd = &cobra.Command{
	Use:   "co-pilot",
	Short: "Co-pilot is a developer tool for automating common tasks on a spring boot project",
	Long:  header(),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
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
	cobra.OnInitialize(initConfig)
	logrus.SetOutput(os.Stdout)
	RootCmd.PersistentFlags().Bool("debug", false, "turn on debug output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	_, err := config.GetActiveProfilePath()
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		if err := config.MigrateToProfiles(); err != nil {
			log.Fatalln(err)
		}
	}

	activeProfilePath, err := config.GetActiveProfilePath()
	if err != nil {
		log.Fatalln(err)
	}

	loadProfile(activeProfilePath)
}

func loadProfile(profilePath string) {

	activeLocalConfig = config.NewLocalConfig(profilePath)
	activeCloudConfig = config.OpenGitCloudConfig(profilePath)
	if !activeLocalConfig.Exists() {
		err := activeLocalConfig.TouchFile()
		if err != nil {
			log.Error(err)
		}
	}
}
