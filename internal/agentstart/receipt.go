package agentstart

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
)

const receiptMode fs.FileMode = 0600

type stored struct {
	Receipt
	Plan  json.RawMessage `json:"plan"`
	State string          `json:"state"`
	Error string          `json:"error"`
}

var legacyStates = map[string]struct{ reached, pending string }{
	"reserved":          {ReachedReserved, ""},
	"tab_created":       {ReachedTabCreated, ""},
	"agent_ready":       {ReachedAgentStarted, ""},
	"prompt_attempting": {ReachedAgentStarted, StagePrompt},
	"prompt_submitted":  {ReachedPromptSubmitted, ""},
}

func ReadReceipt(path string) (Receipt, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Receipt{}, &Error{Code: CodeReceiptUnavailable, Err: err}
	}
	var s stored
	if err = json.Unmarshal(raw, &s); err != nil {
		return Receipt{}, &Error{Code: CodeIncompleteReservation, Err: fmt.Errorf("receipt is empty or unreadable: %w", err)}
	}
	if s.InputSHA256 != "" && s.Reached != "" {
		s.Receipt.Legacy = false
		return s.Receipt, nil
	}
	if state, ok := legacyStates[s.State]; ok && len(s.Plan) > 0 && s.InputSHA256 == "" {
		legacy := Receipt{Reached: state.reached, Pending: state.pending, Pane: s.Pane, Legacy: true}
		if s.Error != "" {
			legacy.Failure = &StartFailure{Stage: "legacy", Class: "legacy_error", EffectUnknown: true}
		}
		return legacy, nil
	}
	return Receipt{}, &Error{Code: CodeIncompleteReservation, Err: errors.New("receipt has no recognized content")}
}

func writeReceipt(path string, r Receipt) error {
	raw, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return atomicfile.Replace(path, raw, receiptMode)
}

func (r Receipt) EffectUnknown() bool {
	if r.Pending != "" {
		return true
	}
	if r.Failure != nil {
		return r.Failure.EffectUnknown
	}
	return r.Reached != ReachedPromptSubmitted
}

func (r Receipt) complete() bool {
	return r.Reached == ReachedPromptSubmitted && r.Pending == "" && r.Failure == nil
}

// resumable is a finished start that stopped before any prompt marker, with a
// recorded tab whose agent can be observed again.
func (r Receipt) resumable() bool {
	switch r.Reached {
	case ReachedTabCreated, ReachedAgentStarted, ReachedSessionBound:
	default:
		return false
	}
	return r.Failure != nil && (r.Pending == "" || r.Pending == StageAgentStart) && r.Tab != "" && r.Pane != "" && r.Terminal != ""
}
