package taskrun

import (
	"reflect"
	"testing"
)

func TestStartupProcessTreeIncludesBackgroundDescendantsOnly(t *testing.T) {
	tree, err := startupProcessTreeFromTable(42, "1 0 1\n42 1 42\n77 42 77\n88 77 77\n91 1 91\n")
	want := []StartupProcess{{42, 1, 42}, {77, 42, 77}, {88, 77, 77}}
	if err != nil || !reflect.DeepEqual(tree, want) {
		t.Fatalf("shell and background descendants = %v, %v; want %v", tree, err, want)
	}
}

func TestStartupProcessTreeRejectsUnknownOrAmbiguousInspection(t *testing.T) {
	for name, table := range map[string]string{
		"empty":         "",
		"missing_shell": "1 0 1\n43 1 43\n",
		"malformed":     "42 1 42\n77 ? 77\n",
		"extra_fields":  "42 1 42 unexpected\n",
		"duplicate":     "42 1 42\n42 1 77\n",
		"cycle":         "42 77 42\n77 42 77\n",
	} {
		t.Run(name, func(t *testing.T) {
			if tree, err := startupProcessTreeFromTable(42, table); err == nil {
				t.Fatalf("unknown inspection accepted: %v", tree)
			}
		})
	}
}

func TestStartupProcessTreeAcceptsObservedIdleShell(t *testing.T) {
	tree, err := startupProcessTreeFromTable(42, "1 0 1\n42 1 42\n91 1 91\n")
	if err != nil || !reflect.DeepEqual(tree, []StartupProcess{{42, 1, 42}}) {
		t.Fatalf("idle shell inspection = %v, %v", tree, err)
	}
}
