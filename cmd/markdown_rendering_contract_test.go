package cmd

import (
	"strings"
	"testing"

	markdown "github.com/MichaelMure/go-term-markdown"
)

func TestTerminalMarkdownReferenceDefinitionsRemainRenderable(t *testing.T) {
	const source = "Read the [guide][guide].\n\n[guide]: https://example.com/guide \"Guide\"\n"
	rendered := string(markdown.Render(source, 120, 2))

	if expected := "Read the [guide](https://example.com/guide Guide)."; !strings.Contains(rendered, expected) {
		t.Fatalf("terminal Markdown output %q does not contain %q", rendered, expected)
	}
	if expected := "https://example.com/guide"; !strings.Contains(rendered, expected) {
		t.Fatalf("terminal Markdown output %q does not contain %q", rendered, expected)
	}
}
