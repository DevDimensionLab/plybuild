// Copyright © 2020 NAME HERE <EMAIL ADDRESS>
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
	"co-pilot/pkg/bitbucket"
	"co-pilot/pkg/git"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var bitbucketCmd = &cobra.Command{
	Use:   "bitbucket",
	Short: "Clone projects from bitbucket",
	Long: `Clone projects from bitbucket, requires a $HOME/.co-pilot.yaml with:

bitbucket_host: <bitbucket_host>
bitbucket_personal_access_token: <bitbucket_personal_access_token> 
`,
	Run: func(cmd *cobra.Command, args []string) {
		bitbucketHost := viper.GetString("bitbucket_host")
		personalAccessToken := viper.GetString("bitbucket_personal_access_token")

		if("" ==  bitbucketHost) || ("" == personalAccessToken) {
			log.Fatalln( "Command requires $HOME/.co-pilot.yaml with bitbucketHost: <bitbucketHost> and personalAccessToken: <personalAccessToken>" )
			os.Exit(1)
		}

		projects, err := bitbucket.GetProjects(bitbucketHost, personalAccessToken )
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		for _, bitBucketProject := range projects.Values {
			projectKey := strings.ToLower(bitBucketProject.Key)
			log.Infoln( "project: " + projectKey )

			bitBucketProjectReposResponse , err := bitbucket.GetProjectRepos(bitbucketHost, personalAccessToken, projectKey)
			if err != nil {
				log.Warnln(err)
				continue
			}

			for _, bitBucketRepo := range bitBucketProjectReposResponse.BitBucketRepo {
				log.Infoln(  "  " + bitBucketRepo.Name )

				err := git.PullRepo(bitbucketHost, ".", "/" + projectKey + "/" + bitBucketRepo.Name)
				if err != nil {
					log.Warnln(err)
					continue
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(bitbucketCmd)
}
