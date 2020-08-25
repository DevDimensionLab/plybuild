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

		if ("" == bitbucketHost) || ("" == personalAccessToken) {
			log.Fatalln("Command requires $HOME/.co-pilot.yaml with bitbucket_host: <bitbucketHost> and bitbucket_personal_access_token: <personalAccessToken>")
			os.Exit(1)
		}

		projects, err := bitbucket.QueryProjects(bitbucketHost, personalAccessToken)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		for _, bitBucketProject := range projects.Values {
			projectKey := strings.ToLower(bitBucketProject.Key)
			log.Infoln("project: " + projectKey)

			bitBucketProjectReposResponse, err := bitbucket.QueryRepos(bitbucketHost, projectKey, personalAccessToken)
			if err != nil {
				log.Warnln(err)
				continue
			}

			for _, bitBucketRepo := range bitBucketProjectReposResponse.BitBucketRepo {
				log.Infoln("  " + bitBucketRepo.Name)

				err := git.PullRepo(bitbucketHost, ".", "/"+projectKey+"/"+bitBucketRepo.Name)
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
