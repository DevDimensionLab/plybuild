package agentwork

import (
	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspaceview"
	"time"
)

// System constructs the public producer adapter. Explicit roots are observed,
// never initialized; no runtime/process/cache or permission probe is performed.
func System(root string) (*Service, error) {
	d := taskwork.SystemDependencies()
	if root == "" {
		w, e := workspace.ObserveContaining(d.Dependencies)
		if e != nil {
			return nil, e
		}
		root = w.Root
	}
	rooted, e := workspace.AtRoot(d.Dependencies, root)
	if e != nil {
		return nil, e
	}
	d.Dependencies = rooted
	h := HerdrRuntime{Executable: "herdr"}
	return &Service{Root: root, Now: time.Now, Runtime: h, Native: func() NativeSnapshot { return readNative(d, root) }, Send: h.Prompt}, nil
}
func readNative(d taskwork.Dependencies, root string) NativeSnapshot {
	out := NativeSnapshot{Vetoes: map[string]string{}}
	v, e := workspaceview.ReadWorkflowStatus(d, workspaceview.WorkflowStatusOptions{SchemaVersion: 2, IncludeAll: true})
	if e != nil {
		return out
	}
	out.Available = true
	out.Partial = v.Freshness != "fresh" || len(v.Diagnostics) > 0
	out.VetoCoverageComplete = true
	return out
}
