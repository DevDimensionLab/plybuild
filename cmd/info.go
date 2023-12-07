package cmd

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/pkg/maven"
	"github.com/devdimensionlab/plybuild/pkg/spring"
	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
)

type InfoOpts struct {
	SpringManaged     bool
	SpringInfo        bool
	MavenRepositories bool
}

var infoOpts InfoOpts

func (infoOpts InfoOpts) Any() bool {
	return infoOpts.SpringManaged ||
		infoOpts.SpringInfo ||
		infoOpts.MavenRepositories
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Prints info on spring version, dependencies etc",
	Long:  `Prints info on spring version, dependencies etc`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return OkHelp(cmd, infoOpts.Any)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if infoOpts.SpringInfo {
			springInfo()
		}
		if infoOpts.SpringManaged {
			showSpringManaged()
		}
		if infoOpts.MavenRepositories {
			showMavenRepositories()
		}
	},
}

func springInfo() {
	repo, err := maven.DefaultRepository()
	if err != nil {
		log.Fatalln(err)
	}

	latestVersionMeta, err := repo.GetMetaData("org.springframework.boot", "spring-boot")
	if err != nil {
		log.Fatalln(err)
	}

	latestVersion, err := latestVersionMeta.LatestRelease()
	if err != nil {
		log.Fatalln(err)
	}

	root, err := spring.GetRoot()
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("Latest version of spring boot are: %s\n", latestVersion)

	log.Infof(fmt.Sprintf("Valid dependencies: "))
	for _, category := range root.Dependencies.Values {
		fmt.Println(fmt.Sprintf("%s", category.Name))
		fmt.Printf("================================\n")
		for _, dep := range category.Values {
			fmt.Printf("[%s]\n    %s, (%s)\n", dep.Id, dep.Name, dep.Description)
		}
		fmt.Printf("\n")
	}
}

func showSpringManaged() {
	deps, err := spring.GetDependencies()
	if err != nil {
		log.Fatalln(err)
	}

	log.Infof(fmt.Sprintf("Spring Boot managed dependencies:"))
	var organized = make(map[string][]pom.Dependency)
	for _, dep := range deps.Dependencies {
		mvnDep := pom.Dependency{
			GroupId:    dep.GroupId,
			ArtifactId: dep.ArtifactId,
		}
		organized[dep.GroupId] = append(organized[dep.GroupId], mvnDep)
	}

	for k, v := range organized {
		fmt.Println(fmt.Sprintf("GroupId: %s", k))
		fmt.Printf("================================\n")
		for _, mvnDep := range v {
			fmt.Printf("  ArtifactId: %s\n", mvnDep.ArtifactId)
		}
	}
}

func showMavenRepositories() {
	settings, _ := maven.NewSettings()
	if err := settings.ListRepositories(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	RootCmd.AddCommand(infoCmd)
	infoCmd.PersistentFlags().BoolVar(&infoOpts.SpringInfo, "spring-info", false, "show spring boot status")
	infoCmd.PersistentFlags().BoolVar(&infoOpts.SpringManaged, "spring-managed", false, "show spring boot managed dependencies info")
	infoCmd.PersistentFlags().BoolVar(&infoOpts.MavenRepositories, "maven-repositories", false, "show current maven repositories")

}
