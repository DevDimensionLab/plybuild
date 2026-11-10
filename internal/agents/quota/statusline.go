package quota

import (
	"fmt"
	"path"
	"strings"
)

// RelayName is the executable basename that identifies a configured relay.
const RelayName = "claude-statusline-relay"

// isRelayCommand reports whether the first POSIX shell word of a status line
// command has the relay basename. The command is parsed, never evaluated.
func isRelayCommand(command string) (bool, error) {
	word, err := firstShellWord(command)
	if err != nil {
		return false, err
	}
	return word != "" && path.Base(word) == RelayName, nil
}

// firstShellWord handles single quotes, double quotes and backslash escapes.
func firstShellWord(s string) (string, error) {
	invalid := fmt.Errorf("%w: status line command quoting", ErrInvalidResponse)
	s = strings.TrimLeft(s, " \t\n")
	var word strings.Builder
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case ' ', '\t', '\n':
			return word.String(), nil
		case '\'':
			end := strings.IndexByte(s[i+1:], '\'')
			if end < 0 {
				return "", invalid
			}
			word.WriteString(s[i+1 : i+1+end])
			i += end + 1
		case '"':
			for i++; ; i++ {
				if i >= len(s) {
					return "", invalid
				}
				if s[i] == '"' {
					break
				}
				if s[i] == '\\' && i+1 < len(s) && strings.IndexByte("\"\\$`\n", s[i+1]) >= 0 {
					i++
				}
				word.WriteByte(s[i])
			}
		case '\\':
			if i+1 >= len(s) {
				return "", invalid
			}
			i++
			if s[i] != '\n' {
				word.WriteByte(s[i])
			}
		default:
			word.WriteByte(c)
		}
	}
	return word.String(), nil
}
