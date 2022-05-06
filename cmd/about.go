package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "About co-pilot",
	Long:  description(),
}

func header() string {
	return fmt.Sprintf(`  _____                  _ _       _
 / ____|                (_) |     | |
| |     ___ ______ _ __  _| | ___ | |_
| |    / _ \______| '_ \| | |/ _ \| __|
| |___| (_) |     | |_) | | | (_) | |_
 \_____\___/      | .__/|_|_|\___/ \__|
                  | |
                  |_|
== version: %s ==
`, version)
}

func description() string {
	about := `## About Co-pilot - https://devdimensionlab.github.io/
- Co-pilot is a command line tool that supports developers using Spring Boot and Maven
- co-pilot upgrade all | 2party | 3party  -> upgrades maven version of maven dependencies to latest 
- co-pilot generate co-pilot.json -> generates a new application ready for business logic
- Authors:
  - Alexander Skjolden, Skjolden Frilans AS
  - Per Otto Christensen, Codify Consulting AS
`
	return fmt.Sprintf("%s\n%s", header(), about)
}

func init() {
	RootCmd.AddCommand(aboutCmd)
}
