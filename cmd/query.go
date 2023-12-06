package cmd

import (
	"fmt"
	"github.com/devdimensionlab/ply/pkg/config"
	"github.com/devdimensionlab/ply/pkg/maven"
	"github.com/spf13/cobra"
)

type QueryOpts struct {
	groupId    string
	artifactId string
}

func (opts QueryOpts) Any() bool {
	return opts.groupId != "" && opts.artifactId != ""
}

var queryOpts QueryOpts

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query dependencies in a project",
	Long:  `Query dependencies in a project`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return OkHelp(cmd, queryOpts.Any)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		ctx.Recursive = true
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		desc := fmt.Sprintf("Search for dependency %s:%s", queryOpts.groupId, queryOpts.artifactId)
		ctx.OnEachMavenProject(desc, func(repository maven.Repository, project config.Project) error {
			dep, err := project.Type.Model().FindDependency(queryOpts.groupId, queryOpts.artifactId)
			if err != nil {
				return nil
			}
			depVersion, err := project.Type.Model().GetDependencyVersion(dep)
			if err != nil {
				depVersion = dep.Version
			}
			log.Infof("Pom-file %s has dependency: %s:%s:%s",
				project.Type.FilePath(), dep.GroupId, dep.ArtifactId, depVersion)

			return nil
		})
	},
}

func init() {
	RootCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringVarP(&queryOpts.groupId, "groupId", "g", "", "groupId")
	queryCmd.Flags().StringVarP(&queryOpts.artifactId, "artifactId", "a", "", "artifactId")
	queryCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
}
