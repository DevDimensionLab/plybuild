package agentstart

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

// startupDialog is a provider screen that blocks startup and that Ply answers
// itself. It matches the last non-empty trimmed lines of a pane read: the
// header lines, the options (one marked as selected) and the footer.
type startupDialog struct {
	name     string
	provider string
	// source is the pane read source that shows the dialog.
	source   string
	header   func(cwd string) []func(string) bool
	options  []func(string) bool
	selected string
	footer   string
	// target is the option Ply selects; it may be reached from the option
	// above it with one down.
	target int
	record func(*Receipt)
}

func equals(want string) func(string) bool {
	return func(line string) bool { return line == want }
}

func prefixed(want string) func(string) bool {
	return func(line string) bool { return strings.HasPrefix(line, want) }
}

var codexUpdateHeading = regexp.MustCompile(`^(✨ )?Update available · \d+\.\d+\.\d+ → \d+\.\d+\.\d+$`)

var startupDialogs = []startupDialog{
	// Claude's folder trust dialog, observed with Claude Code 2.1.309. The
	// question line is matched by its first sentence; the rest is advice.
	{
		name: "Claude folder trust", provider: ProviderClaude, source: "recent-unwrapped",
		header: func(cwd string) []func(string) bool {
			return []func(string) bool{equals("Accessing workspace:"), equals(cwd),
				prefixed("Quick safety check: Is this a project you created or one you trust? "),
				equals("Claude Code'll be able to read, edit, and execute files here."), equals("Security guide")}
		},
		options:  []func(string) bool{equals("No, exit"), equals("Yes, I trust this folder")},
		selected: "❯ ", footer: "Enter to confirm · Esc to cancel", target: 1,
		record: func(r *Receipt) { r.FolderTrustAccepted = true },
	},
	// Codex's update prompt, observed with codex-cli 0.164.0. Skip leaves the
	// installed version and the saved settings unchanged; it is drawn on the
	// alternate screen.
	{
		name: "Codex update", provider: ProviderCodex, source: "visible",
		header: func(string) []func(string) bool {
			return []func(string) bool{codexUpdateHeading.MatchString, equals("Release notes: https://github.com/openai/codex/releases/latest")}
		},
		options:  []func(string) bool{prefixed("1. Update now"), equals("2. Skip"), equals("3. Skip until next version")},
		selected: "› ", footer: "enter continue · esc skip", target: 1,
		record: func(r *Receipt) { r.CodexUpdateSkipped = true },
	},
}

// choice reads the pane and returns the selected option when its end is
// exactly this dialog for the requested cwd.
func (s *session) choice(d startupDialog) (int, bool) {
	screen, err := s.o.text(s.ctx, s.in.Herdr, []string{"pane", "read", s.r.Pane, "--source", d.source, "--lines", "60"})
	if err != nil {
		return 0, false
	}
	var lines []string
	for _, line := range strings.Split(screen, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}
	header := d.header(s.in.Cwd)
	n := len(header) + len(d.options) + 1
	if len(lines) < n {
		return 0, false
	}
	lines = lines[len(lines)-n:]
	for i, match := range header {
		if !match(lines[i]) {
			return 0, false
		}
	}
	selected := -1
	for i, match := range d.options {
		line := lines[len(header)+i]
		if marked, ok := strings.CutPrefix(line, d.selected); ok {
			if selected >= 0 {
				return 0, false
			}
			selected, line = i, marked
		}
		if !match(line) {
			return 0, false
		}
	}
	return selected, selected >= 0 && lines[n-1] == d.footer
}

// dialog returns the startup dialog the pane shows for this provider.
func (s *session) dialog() (startupDialog, int, bool) {
	for _, d := range startupDialogs {
		if d.provider != s.in.Provider {
			continue
		}
		if selected, ok := s.choice(d); ok && (selected == d.target || selected == d.target-1) {
			return d, selected, true
		}
	}
	return startupDialog{}, 0, false
}

// blockedOnDialog reports whether a failed agent start left the reserved agent
// blocked on a startup dialog that Ply answers.
func (s *session) blockedOnDialog() bool {
	raw, err := s.o.call(s.ctx, s.in.Herdr, []string{"agent", "get", s.r.Pane})
	if err != nil {
		return false
	}
	var got struct {
		Agent agentView `json:"agent"`
	}
	if json.Unmarshal(raw, &got) != nil || s.identity(got.Agent) != nil || got.Agent.Status != "blocked" {
		return false
	}
	_, _, ok := s.dialog()
	return ok
}

// answerDialog selects the target option once in each startup dialog Ply
// answers. Any other blocked screen is left to the human.
func (s *session) answerDialog() error {
	d, selected, ok := s.dialog()
	if !ok || s.answered[d.name] {
		return nil
	}
	if s.answered == nil {
		s.answered = map[string]bool{}
	}
	s.answered[d.name] = true
	if selected != d.target {
		if _, err := s.o.text(s.ctx, s.in.Herdr, []string{"pane", "send-keys", s.r.Pane, "down"}); err != nil {
			return s.fail(StageAwaitReady, ClassAgentResponse, false, err)
		}
		if selected, ok = s.choice(d); !ok || selected != d.target {
			return s.fail(StageAwaitReady, ClassNotReady, false, errors.New(d.name+" dialog did not select the intended option"))
		}
	}
	if _, err := s.o.text(s.ctx, s.in.Herdr, []string{"pane", "send-keys", s.r.Pane, "enter"}); err != nil {
		return s.fail(StageAwaitReady, ClassAgentResponse, false, err)
	}
	d.record(&s.r)
	return nil
}
