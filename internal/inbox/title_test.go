package inbox

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// S8: derived-title-v1.
func TestDerivedTitleV1(t *testing.T) {
	sixty := strings.Repeat("a", 60)
	cases := []struct{ name, body, want string }{
		{"empty body", "", FallbackTitle},
		{"only blank and marker lines", "\n   \n#\n- \n> \n1.\n", FallbackTitle},
		{"heading", "## Spec ferdig", "Spec ferdig"},
		{"nested markers", "> - **Spec** `ferdig`", "Spec ferdig"},
		{"numbered list", "12. Første punkt", "Første punkt"},
		{"plus and star", "+ * Punkt", "Punkt"},
		{"marker without space is text", "#tag og >pil", "#tag og >pil"},
		{"bold start is not a list marker", "**Viktig:** les", "Viktig: les"},
		{"leading blank lines", "\n\n   \n  Tekst  her\nNeste", "Tekst her"},
		{"crlf", "\r\n# Tittel\r\nmer", "Tittel"},
		{"collapse whitespace", "a \t  b", "a b"},
		{"exactly 60", sixty, sixty},
		{"61 without space", strings.Repeat("b", 61), strings.Repeat("b", 59) + "…"},
		{"61 with space", "ord " + strings.Repeat("c", 57), "ord…"},
		{"cut at last space and trim", strings.Repeat("d", 50) + "  " + strings.Repeat("e", 20), strings.Repeat("d", 50) + "…"},
		{"emoji counted as code points", strings.Repeat("🙂", 61), strings.Repeat("🙂", 59) + "…"},
		{"backticks only", "``` \n`code`", "code"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := DeriveTitle(c.body)
			if got != c.want {
				t.Fatalf("DeriveTitle(%q) = %q, want %q", c.body, got, c.want)
			}
			if n := utf8.RuneCountInString(got); n < 1 || n > maxTitleRunes {
				t.Fatalf("title length %d", n)
			}
		})
	}
}
