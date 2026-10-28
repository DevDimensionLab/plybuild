package agentwork

import (
	"sort"
)

func (s *Service) Overview() (Overview, error) {
	out := Overview{Kind: OverviewKind, SchemaVersion: 1, ObservedAtUTC: s.now(), Sources: []Source{}, Agents: []Agent{}, Followups: []Followup{}}
	source := Source{Source: "herdr.agent.list", Status: "unavailable", Freshness: "unknown", Coverage: "Recognized codex/claude agents in Herdr's current user scope; maximum 4096. No transcript, cache or process probes."}
	n := s.native()
	ns := Source{Source: "ply.registered", Status: "unavailable", Freshness: "unknown", Coverage: "Exact provider/session/workspace/tab/pane/terminal binding to registered native attempts."}
	if n.Available {
		ns.Status = "available"
		ns.Freshness = "fresh"
		if len(n.Bindings) == 0 {
			ns.Status = "empty"
		}
		if n.Partial {
			ns.Status = "partial"
			ns.Freshness = "unknown"
		}
	}
	if s.Runtime != nil {
		agents, e := s.Runtime.List()
		if e == nil {
			source.Status = "available"
			source.Freshness = "fresh"
			if len(agents) == 0 {
				source.Status = "empty"
			}
			seen := map[string]bool{}
			for index, a := range agents {
				if index >= 4096 {
					source.Status = "partial"
					source.Freshness = "unknown"
					break
				}
				if a.Provider != "codex" && a.Provider != "claude" {
					continue
				}
				if !a.valid() {
					source.Status = "partial"
					source.Freshness = "unknown"
				}
				if !reference(a.Pane) || !reference(a.Tab) || !reference(a.Workspace) {
					source.Status = "partial"
					source.Freshness = "unknown"
					continue
				}
				ref := a.Ref()
				if seen[ref] {
					source.Status = "partial"
					source.Freshness = "unknown"
					continue
				}
				seen[ref] = true
				state := a.Status
				switch state {
				case "working", "idle", "done", "blocked", "unknown":
				default:
					state = "unknown"
				}
				w, veto := binding(n, a.Identity)
				bs := "unbound"
				if w != nil {
					bs = "verified"
				}
				if !a.valid() {
					w = nil
					bs = "unknown"
				}
				caps := map[string]Capability{"register": {State: "unknown", Reason: "Caller must supply authorized owner_ref, mandate_ref and explicit registrar provenance; target_ref is agent_ref."}, "first_request": {State: "unknown", Reason: "Explicitly register canonical wait_kind=status with wait_ref=agent_ref, then request its returned ID; refresh never registers."}, "dispatch": {State: "unknown", Reason: "No follow-up cycle selected."}}
				if !a.valid() {
					caps["register"] = Capability{State: "disabled", Reason: "Complete stable execution identity is unavailable."}
					caps["first_request"] = caps["register"]
				}
				if veto != "clear" {
					caps["dispatch"] = Capability{State: "unknown", Reason: "Native veto coverage unavailable or vetoed; follow-up dispatch must revalidate."}
				}
				out.Agents = append(out.Agents, Agent{AgentRef: ref, Provider: a.Provider, Status: state, Workspace: a.Workspace, Tab: a.Tab, Pane: a.Pane, BindingStatus: bs, Work: w, Capabilities: caps})
			}
		}
	}
	if len(out.Agents) == 0 && source.Status == "available" {
		source.Status = "empty"
	}
	out.Sources = append(out.Sources, source, ns)
	list, e := s.list(false, func() NativeSnapshot { return n })
	fs := Source{Source: "ply.followups", Status: "available", Freshness: "fresh", Coverage: "Local status-only follow-up store, maximum 512; no Task status copies."}
	if e != nil {
		fs.Status = "unavailable"
		fs.Freshness = "unknown"
	} else {
		out.Followups = list.Items
		if len(list.Items) == 0 {
			fs.Status = "empty"
		}
	}
	out.Sources = append(out.Sources, fs)
	for i := range out.Agents {
		a := &out.Agents[i]
		a.FollowupRefs = []Subject{}
		for _, f := range out.Followups {
			if f.TargetRef == a.AgentRef {
				a.FollowupRefs = append(a.FollowupRefs, subject(f))
			}
		}
		for op, c := range a.Capabilities {
			c.Subject = Subject{AgentRef: a.AgentRef}
			a.Capabilities[op] = c
		}
	}
	sort.Slice(out.Agents, func(i, j int) bool { return out.Agents[i].AgentRef < out.Agents[j].AgentRef })
	if labels, ok := s.Runtime.(interface {
		tabLabels([]Agent) map[tabKey]string
	}); ok {
		metadata := labels.tabLabels(out.Agents)
		for i := range out.Agents {
			a := &out.Agents[i]
			if label := metadata[tabKey{a.Workspace, a.Tab}]; validPublicTabLabel(label) {
				a.TabLabel = label
			}
		}
	}
	return out, nil
}
