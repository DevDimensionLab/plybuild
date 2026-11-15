package agentstart

import (
	"encoding/json"
	"errors"
	"strings"
)

const (
	trustSelected = "❯ "
	trustNo       = "No, exit"
	trustYes      = "Yes, I trust this folder"
)

// trustQuestion starts the question line of Claude's folder trust dialog
// (observed with Claude Code 2.1.309); the rest of that line is advice.
const trustQuestion = "Quick safety check: Is this a project you created or one you trust? "

// trustChoice reads the pane and returns the selected option when its end is
// exactly Claude's folder trust dialog for the requested cwd.
func (s *session) trustChoice() (string, bool) {
	screen, err := s.o.text(s.ctx, s.in.Herdr, []string{"pane", "read", s.r.Pane, "--source", "recent-unwrapped", "--lines", "60"})
	if err != nil {
		return "", false
	}
	var lines []string
	for _, line := range strings.Split(screen, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}
	want := []string{"Accessing workspace:", s.in.Cwd, trustQuestion, "Claude Code'll be able to read, edit, and execute files here.",
		"Security guide", "", "", "Enter to confirm · Esc to cancel"}
	if len(lines) < len(want) {
		return "", false
	}
	lines = lines[len(lines)-len(want):]
	question, no, yes := 2, len(want)-3, len(want)-2
	for i, line := range want {
		switch i {
		case question:
			if !strings.HasPrefix(lines[i], line) {
				return "", false
			}
		case no, yes:
		default:
			if lines[i] != line {
				return "", false
			}
		}
	}
	switch {
	case lines[no] == trustSelected+trustNo && lines[yes] == trustYes:
		return trustNo, true
	case lines[no] == trustNo && lines[yes] == trustSelected+trustYes:
		return trustYes, true
	}
	return "", false
}

// blockedOnTrust reports whether a failed Claude agent start left the reserved
// agent blocked on the folder trust dialog for the requested cwd.
func (s *session) blockedOnTrust() bool {
	if s.in.Provider != ProviderClaude {
		return false
	}
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
	_, ok := s.trustChoice()
	return ok
}

// answerTrust selects Yes once in Claude's folder trust dialog for the
// requested cwd. Any other blocked screen is left to the human.
func (s *session) answerTrust() error {
	if s.in.Provider != ProviderClaude || s.trustAnswered {
		return nil
	}
	choice, ok := s.trustChoice()
	if !ok {
		return nil
	}
	s.trustAnswered = true
	if choice == trustNo {
		if _, err := s.o.text(s.ctx, s.in.Herdr, []string{"pane", "send-keys", s.r.Pane, "down"}); err != nil {
			return s.fail(StageAwaitReady, ClassAgentResponse, false, err)
		}
		if choice, ok = s.trustChoice(); !ok || choice != trustYes {
			return s.fail(StageAwaitReady, ClassNotReady, false, errors.New("folder trust dialog did not select Yes"))
		}
	}
	if _, err := s.o.text(s.ctx, s.in.Herdr, []string{"pane", "send-keys", s.r.Pane, "enter"}); err != nil {
		return s.fail(StageAwaitReady, ClassAgentResponse, false, err)
	}
	s.r.FolderTrustAccepted = true
	return nil
}
