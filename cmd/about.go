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
  _____  _       
 |  __ \| |      
 | |__) | |_   _ 
 |  ___/| | | | |
 | |    | | |_| |
 |_|    |_|\__, |
            __/ |
           |___/  

== version: %s ==
`, version)
}

func description() string {
	about := `## About Ply - https://devdimensionlab.github.io/
- Ply is a command line tool that supports developers using Spring Boot and Maven
- ply upgrade all | 2party | 3party  -> upgrades maven version of maven dependencies to latest 
- ply generate ply.json -> generates a new application ready for business logic
- Authors:
  - Alexander Skjolden, Skjolden Frilans AS
  - Per Otto Christensen, Codify Consulting AS
`
	return fmt.Sprintf("%s\n%s", header(), about)
}

func init() {
	RootCmd.AddCommand(aboutCmd)
}
