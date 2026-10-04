package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	notification "github.com/devdimensionlab/plybuild/internal/workflownotification"
	"github.com/spf13/cobra"
)

func TestNotificationArgumentErrorsAreTypedAndPrivate(t *testing.T) {
	for _, args := range [][]string{{"send", "--format", "json"}, {"send", "--unexpected", "private-argument", "--format", "json"}, {"show", "id", "--route", "route", "--format", "bad"}, {"send", "--route", "route", "--file", "file", "--check", "--apply", "--confirm", "x", "--format", "json"}, {"send", "--route", "route", "--file", "file", "--confirm", "x", "--format", "json"}, {"private-argument"}} {
		t.Run(args[0]+args[len(args)-1], func(t *testing.T) {
			root := &cobra.Command{Use: "ply", SilenceErrors: true, SilenceUsage: true}
			c := NewWorkflowNotificationCommand(notification.Dependencies{})
			root.AddCommand(c)
			var out, diagnostic bytes.Buffer
			root.SetOut(&out)
			root.SetErr(&diagnostic)
			full := append([]string{"notification"}, args...)
			root.SetArgs(full)
			old := os.Args
			os.Args = append([]string{"ply"}, full...)
			defer func() { os.Args = old }()
			e := root.Execute()
			if ExitCode(e) != 2 {
				t.Fatalf("wrong exit: %v", e)
			}
			if bytes.Contains(append(out.Bytes(), diagnostic.Bytes()...), []byte("private-argument")) || bytes.Contains([]byte(e.Error()), []byte("private-argument")) {
				t.Fatal("raw argument leaked")
			}
			if args[len(args)-1] == "json" {
				var v map[string]any
				if json.Unmarshal(out.Bytes(), &v) != nil || v["kind"] != "ply.workflow.notification-error" {
					t.Fatalf("missing JSON error: %s", out.String())
				}
			}
		})
	}
}
