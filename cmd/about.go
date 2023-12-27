package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "About ply",
	Long:  description(),
}

func header() string {
	return fmt.Sprintf(`
           __      __          _ __    __
    ____  / /_  __/ /_  __  __(_) /___/ /
   / __ \/ / / / / __ \/ / / / / / __  / 
  / /_/ / / /_/ / /_/ / /_/ / / / /_/ /  
 / .___/_/\__, /_.___/\__,_/_/_/\__,_/   
/_/      /____/                          

== version: %s ==
`, version)
}

func description() string {
	about := `## About plybuild - https://devdimensionlab.github.io/
- Plybuild is a command line tool that supports developers using Spring Boot and Maven
- ply upgrade all | 2party | 3party  -> upgrades maven version of maven dependencies to latest 
- ply build ply.json -> generates a new application ready for business logic
- Authors:
  - Alexander Skjolden, Runwith AS
  - Per Otto Christensen, Codify Consulting AS
`
	return fmt.Sprintf("%s\n%s", header(), about)
}

func completionCommand() *cobra.Command {
	return &cobra.Command{
		Use: "completion",
	}
}

func init() {
	RootCmd.AddCommand(aboutCmd)

	completion := completionCommand()
	completion.Hidden = true
	RootCmd.AddCommand(completion)
}
