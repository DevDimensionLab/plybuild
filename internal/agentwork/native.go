package agentwork

import (
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspaceview"
	"time"
)

type rootedFiles struct {
	workspace.FileSystem
	root string
}

func (f rootedFiles) Getwd() (string, error) { return f.root, nil }

// System constructs the public producer adapter. Explicit roots are observed,
// never initialized; no runtime/process/cache or permission probe is performed.
func System(root string) (*Service, error) {
	d := workspace.SystemDependencies()
	if root == "" {
		w, e := workspace.ObserveContaining(d)
		if e != nil {
			return nil, e
		}
		root = w.Root
	} else {
		w, e := workspace.ObserveRoot(d, root)
		if e != nil || w.Root != root {
			return nil, workspace.InvalidArguments("workspace must be its physical root")
		}
	}
	d.Files = rootedFiles{d.Files, root}
	h := HerdrRuntime{Executable: "herdr"}
	return &Service{Root: root, Now: time.Now, Runtime: h, Native: func() NativeSnapshot { return readNative(d, root) }, Send: h.Prompt}, nil
}
func deref(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
func readNative(d workspace.Dependencies, root string) NativeSnapshot {
	out := NativeSnapshot{Vetoes: map[string]string{}}
	v, inventory, e := workspaceview.ReadWorkflowStatusWithInventory(d, workspaceview.WorkflowStatusOptions{SchemaVersion: 2, IncludeAll: true})
	if e != nil {
		return out
	}
	out.Available = true
	out.Partial = v.Freshness != "fresh" || len(v.Diagnostics) > 0
	scopeInventory(&out, inventory, nil)

	byID := map[string]taskrun.InventoryRun{}
	for _, run := range inventory.Runs {
		byID[run.RunID] = run
	}
	for _, row := range v.Items {
		for _, prior := range row.Runs {
			run, ok := byID[prior.RunID]
			if !ok {
				out.CoverageUnknown = true
				out.VetoCoverageComplete = false
				continue
			}
			if run.Herdr == nil || run.Provider == nil {
				continue
			}
			h := run.Herdr
			i := Identity{*run.Provider, deref(h.AgentSessionID), deref(h.WorkspaceID), deref(h.TabID), deref(h.PaneID), deref(h.TerminalID)}
			if !i.valid() {
				continue
			}
			w := Work{ProjectID: string(row.ProjectID), RepoID: string(row.RepoID), EpicID: string(row.EpicID), TaskID: string(row.TaskID), RunID: run.RunID, AttemptRef: run.RunID, Phase: row.Progress.State, NextActor: "unknown", Source: "registered", Freshness: run.Freshness, SourceObservedAtUTC: row.LastActivityUTC}
			veto := run.StatusFollowupVeto
			if run.Freshness != "fresh" {
				veto = "native_sources_unknown"
			}
			if run.Delivery != nil {
				w.Phase = run.Delivery.Phase
				w.NextActor = run.Delivery.NextAction.Actor
				w.NextAction = run.Delivery.NextAction.Kind
				w.SourceObservedAtUTC = run.Delivery.LastActivityUTC
			}
			w.Reason = "Read the registered native action before continuing."
			w.DetailRef = "workflow.trace:" + w.TaskID
			out.Bindings = append(out.Bindings, Binding{Identity: i, Work: w, Veto: veto})
		}
	}
	return out
}

// Enumeration failures or a genuinely unidentifiable/changing attempt can hide
// this target. A damaged event on a separately verified session is scoped to
// that session and remains visible in source diagnostics without vetoing others.
func scopeInventory(out *NativeSnapshot, inventory taskrun.Inventory, err error) {
	if out.Vetoes == nil {
		out.Vetoes = map[string]string{}
	}
	if err != nil || len(inventory.Reasons) > 0 {
		out.CoverageUnknown = true
		out.Partial = true
	}
	for _, run := range inventory.Runs {
		if !run.StatusIdentityKnown || run.Freshness == "stale" || run.Provider == nil || run.Herdr == nil || run.Herdr.AgentSessionID == nil {
			out.CoverageUnknown = true
			out.Partial = true
			continue
		}
		key := *run.Provider + ":" + *run.Herdr.AgentSessionID
		if out.Vetoes[key] != "" && out.Vetoes[key] != "registered_session_requires_exact_binding" {
			continue
		}
		if run.Freshness != "fresh" {
			out.Partial = true
			out.Vetoes[key] = "native_sources_unknown"
		} else if run.StatusFollowupVeto != "clear" {
			out.Vetoes[key] = "native_veto_unknown_or_active"
		} else {
			out.Vetoes[key] = "registered_session_requires_exact_binding"
		}
	}
	out.VetoCoverageComplete = !out.CoverageUnknown
}
