package agentwork

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

// HerdrRuntime uses supported agent list/get/prompt and public tab-list metadata.
// The adapter bounds time and bytes and never exports raw output or errors.
type HerdrRuntime struct {
	Executable string
	Timeout    time.Duration
}
type limitedOutput struct {
	data     []byte
	overflow bool
}

func (b *limitedOutput) Write(p []byte) (int, error) {
	n := len(p)
	left := (8 << 20) - len(b.data)
	if n > left {
		b.overflow = true
		p = p[:left]
	}
	b.data = append(b.data, p...)
	return n, nil
}
func (h HerdrRuntime) call(args ...string) (json.RawMessage, error) {
	limit := h.Timeout
	if limit <= 0 {
		limit = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()
	c := exec.CommandContext(ctx, h.Executable, args...)
	c.WaitDelay = time.Second
	var out, errout limitedOutput
	c.Stdout = &out
	c.Stderr = &errout
	if e := c.Run(); e != nil || out.overflow || errout.overflow {
		return nil, fmt.Errorf("herdr_source_unavailable")
	}
	var env struct {
		Result json.RawMessage `json:"result"`
		Error  json.RawMessage `json:"error"`
	}
	if !json.Valid(out.data) || json.Unmarshal(out.data, &env) != nil || len(env.Result) == 0 || string(env.Result) == "null" || len(env.Error) > 0 && string(env.Error) != "null" {
		return nil, fmt.Errorf("herdr_response_invalid")
	}
	return env.Result, nil
}

type wireAgent struct {
	Provider      string `json:"agent"`
	Status        string `json:"agent_status"`
	Workspace     string `json:"workspace_id"`
	Tab           string `json:"tab_id"`
	Pane          string `json:"pane_id"`
	Terminal      string `json:"terminal_id"`
	Ready         *bool  `json:"interactive_ready"`
	LaunchPending bool   `json:"launch_pending"`
	Session       *struct {
		Agent string `json:"agent"`
		Kind  string `json:"kind"`
		Value string `json:"value"`
	} `json:"agent_session"`
}

func (w wireAgent) observation() Observation {
	i := Identity{Provider: w.Provider, Workspace: w.Workspace, Tab: w.Tab, Pane: w.Pane, Terminal: w.Terminal}
	if w.Session != nil && w.Session.Agent == w.Provider && w.Session.Kind == "id" {
		i.Session = w.Session.Value
	}
	return Observation{Identity: i, Status: w.Status, Ready: w.Ready, LaunchPending: w.LaunchPending}
}
func (h HerdrRuntime) List() ([]Observation, error) {
	raw, e := h.call("agent", "list")
	if e != nil {
		return nil, e
	}
	var v struct {
		Agents *[]wireAgent `json:"agents"`
	}
	if json.Unmarshal(raw, &v) != nil || v.Agents == nil {
		return nil, fmt.Errorf("herdr_list_invalid")
	}
	o := []Observation{}
	for _, a := range *v.Agents {
		o = append(o, a.observation())
	}
	return o, nil
}
func (h HerdrRuntime) Get(i Identity) (Observation, error) {
	raw, e := h.call("agent", "get", i.Pane)
	if e != nil {
		return Observation{}, e
	}
	var v struct {
		Agent *wireAgent `json:"agent"`
	}
	if json.Unmarshal(raw, &v) != nil || v.Agent == nil {
		return Observation{}, fmt.Errorf("herdr_agent_invalid")
	}
	return v.Agent.observation(), nil
}
func (h HerdrRuntime) Prompt(i Identity, message string) error {
	raw, e := h.call("agent", "prompt", i.Pane, message)
	if e != nil {
		return e
	}
	var v struct {
		Type string `json:"type"`
	}
	if json.Unmarshal(raw, &v) != nil || v.Type != "agent_prompted" {
		return fmt.Errorf("herdr_ack_unknown")
	}
	return nil
}
