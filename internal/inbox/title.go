package inbox

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// DerivedTitleRule names the only title fallback rule; stored items keep the
// rule they were written with.
const DerivedTitleRule = "derived-title-v1"

// FallbackTitle is the derived title of a body without text.
const FallbackTitle = "Svar uten tekst"

// DeriveTitle is derived-title-v1: the first line with text after markdown
// markers, bold markers and backticks are removed, cut to 60 code points.
func DeriveTitle(body string) string {
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSuffix(line, "\r")
		for {
			line = strings.TrimLeftFunc(line, unicode.IsSpace)
			marker := markerLength(line)
			if marker == 0 {
				break
			}
			line = line[marker:]
		}
		line = strings.ReplaceAll(line, "**", "")
		line = strings.ReplaceAll(line, "`", "")
		if candidate := strings.Join(strings.Fields(line), " "); candidate != "" {
			return shorten(candidate)
		}
	}
	return FallbackTitle
}

// markerLength is the byte length of one leading #+, >, -, *, + or <digits>.
// marker followed by whitespace or the end, or 0.
func markerLength(line string) int {
	end := 0
	switch {
	case strings.HasPrefix(line, "#"):
		end = len(line) - len(strings.TrimLeft(line, "#"))
	case strings.HasPrefix(line, ">"), strings.HasPrefix(line, "-"), strings.HasPrefix(line, "*"), strings.HasPrefix(line, "+"):
		end = 1
	default:
		digits := len(line) - len(strings.TrimLeft(line, "0123456789"))
		if digits == 0 || !strings.HasPrefix(line[digits:], ".") {
			return 0
		}
		end = digits + 1
	}
	if end == len(line) {
		return end
	}
	if next, _ := utf8.DecodeRuneInString(line[end:]); unicode.IsSpace(next) {
		return end
	}
	return 0
}

func shorten(candidate string) string {
	runes := []rune(candidate)
	if len(runes) <= maxTitleRunes {
		return candidate
	}
	cut := string(runes[:maxTitleRunes-1])
	if space := strings.LastIndex(cut, " "); space >= 0 {
		cut = strings.TrimRightFunc(cut[:space], unicode.IsSpace)
	}
	return cut + "…"
}
