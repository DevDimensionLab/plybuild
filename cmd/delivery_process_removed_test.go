package cmd

import (
	"strings"
	"testing"
)

// The old delivery process is gone from the CLI: an unknown command or flag
// gets the ordinary Cobra error, without aliases or a special removal message.
func TestOldDeliveryProcessCommandsAndFlagsDoNotExist(t *testing.T) {
	for _, path := range [][]string{
		{"workspace", "task", "result"},
		{"workspace", "task", "qa"},
		{"workspace", "task", "prepare"},
		{"workspace", "task", "preparation"},
		{"workspace", "task", "queue", "advance"},
		{"workspace", "task", "queue", "release"},
		{"workspace", "task", "run"},
		{"workflow", "run"},
		{"workflow", "handoff"},
	} {
		if found, _, err := RootCmd.Find(path); err == nil && found.Name() == path[len(path)-1] {
			t.Errorf("%s still exists", strings.Join(path, " "))
		}
	}
	for command, flags := range map[string][]string{
		"workspace task goal register": {"delivery", "acceptance", "responsible-actor", "acceptance-timeout", "require-human-qa", "github-repository", "remote", "target-ref"},
		"workspace task integrate":     {"result", "qa"},
	} {
		found, _, err := RootCmd.Find(strings.Fields(command))
		if err != nil {
			t.Fatalf("%s: %v", command, err)
		}
		for _, flag := range flags {
			if found.Flags().Lookup(flag) != nil {
				t.Errorf("%s still takes --%s", command, flag)
			}
		}
	}
	for _, path := range [][]string{
		{"workspace", "task", "integrate"},
		{"workspace", "task", "queue", "list"},
		{"workspace", "task", "queue", "set"},
		{"workflow", "goal-start"},
		{"workflow", "status"},
		{"workflow", "agents"},
		{"workflow", "input"},
		{"workflow", "notification", "send"},
	} {
		if found, _, err := RootCmd.Find(path); err != nil || found.Name() != path[len(path)-1] {
			t.Errorf("%s is missing", strings.Join(path, " "))
		}
	}
}
