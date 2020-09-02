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
	"co-pilot/pkg/config"
	"co-pilot/pkg/logger"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var log = logger.Context()

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "co-pilot",
	Short: "Co-pilot er et støtteverktøy i utvikling der kjente oppgaver er automatisert",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			log.Fatalln(err)
		}
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
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

	if !config.LocalConfigExists() {
		err := config.TouchLocalConfigFile()
		if err != nil {
			log.Error(err)
		}
	} else {
		f, err := config.LocalConfigFilePath()
		if err != nil {
			log.Error(err)
		}
		log.Infof("using config file %s", f)
	}
}
