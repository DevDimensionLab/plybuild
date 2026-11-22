package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// docs/ is not in Git, so help text must not send a clone's user there.
func TestHelpTextNamesNoLocalDocs(t *testing.T) {
	var visit func(*cobra.Command)
	visit = func(c *cobra.Command) {
		for _, text := range []string{c.Short, c.Long, c.Example} {
			if strings.Contains(text, "docs/") {
				t.Errorf("%q help names docs/", c.CommandPath())
			}
		}
		for _, child := range c.Commands() {
			visit(child)
		}
	}
	visit(RootCmd)
}
