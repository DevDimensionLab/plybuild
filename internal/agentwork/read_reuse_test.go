package agentwork

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestOverviewReusesOneNativeRead(t *testing.T) {
	for _, records := range []int{0, 1, 12} {
		t.Run(fmt.Sprint(records), func(t *testing.T) {
			s, draft, _ := fixture(t)
			for i := 0; i < records; i++ {
				d := draft
				d.PublicationKey = fmt.Sprintf("read/%d", i)
				d.WaitRef = fmt.Sprintf("wait-%d", i)
				if _, err := s.Register(d); err != nil {
					t.Fatal(err)
				}
			}
			basis := s.native()
			calls := 0
			s.Native = func() NativeSnapshot {
				calls++
				n := basis
				n.Bindings = append([]Binding(nil), basis.Bindings...)
				n.Bindings[0].Work.Phase = fmt.Sprintf("read-%d", calls)
				return n
			}
			for request := 1; request <= 2; request++ {
				out, err := s.Overview()
				if err != nil {
					t.Fatal(err)
				}
				if calls != request {
					t.Errorf("request %d read native %d total times; want %d", request, calls, request)
				}
				if len(out.Agents) != 1 || len(out.Followups) != records {
					t.Fatalf("lost rows: %+v", out)
				}
				phase := fmt.Sprintf("read-%d", request)
				if out.Agents[0].Work == nil || out.Agents[0].Work.Phase != phase {
					t.Fatal("agent used another request basis", out.Agents)
				}
				for _, f := range out.Followups {
					if f.Work == nil || f.Work.Phase != phase {
						t.Error("followup and agent used different native facts", f.Work)
					}
				}
			}
		})
	}
}

func ptr(s string) *string { return &s }

type rootedFiles struct {
	workspace.FileSystem
	root string
}

func (f rootedFiles) Getwd() (string, error) { return f.root, nil }

// Without the old Run engine the native snapshot is the workflow status read:
// files left under workflow-runs carry no meaning and veto nothing.
func TestNativeSnapshotIgnoresOldRunFiles(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	d := workspace.SystemDependencies()
	d.Files = rootedFiles{d.Files, root}
	if _, err := workspace.Init(d); err != nil {
		t.Fatal(err)
	}
	requests := filepath.Join(root, ".ply", "workflow-runs", "requests")
	if err := os.MkdirAll(requests, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(requests, "unrecognized.json"), []byte("damaged"), 0600); err != nil {
		t.Fatal(err)
	}
	n := readNative(d, root)
	if !n.Available || n.Partial || n.CoverageUnknown || !n.VetoCoverageComplete || len(n.Bindings) != 0 || len(n.Vetoes) != 0 {
		t.Fatalf("old Run files changed the native snapshot: %+v", n)
	}
}

func TestOverviewBasisDoesNotCrossOperationBoundary(t *testing.T) {
	for _, state := range []string{"partial", "unavailable", "same_session_veto"} {
		t.Run(state, func(t *testing.T) {
			s, d, sends := fixture(t)
			d.NextCheckAtUTC = s.now()
			f, err := s.Register(d)
			if err != nil {
				t.Fatal(err)
			}
			basis := s.native()
			reads := 0
			s.Native = func() NativeSnapshot { reads++; return basis }
			if _, err := s.Overview(); err != nil {
				t.Fatal(err)
			}
			basis.Bindings = append([]Binding(nil), basis.Bindings...)
			switch state {
			case "partial":
				basis.Partial, basis.CoverageUnknown = true, true
			case "unavailable":
				basis.Available, basis.Bindings = false, nil
			case "same_session_veto":
				basis.Vetoes = map[string]string{s.Runtime.(*fakeRuntime).agents[0].Identity.sessionKey(): "native_refusal"}
			}
			out, err := s.Overview()
			if err != nil {
				t.Fatal(err)
			}
			want := "available"
			if state == "partial" {
				want = "partial"
			}
			if state == "unavailable" {
				want = "unavailable"
			}
			if reads != 2 || out.Sources[1].Status != want {
				t.Fatalf("stale request basis reused: reads=%d source=%+v", reads, out.Sources[1])
			}
			if state != "same_session_veto" && out.Sources[1].Freshness != "unknown" {
				t.Fatal("unknown source became fresh", out.Sources[1])
			}
			before := reads
			result, err := s.Dispatch(f.ID)
			if err != nil || sends.Load() != 0 || result.Transport.State != "not_sent" || reads <= before {
				t.Fatalf("dispatch reused read authority: reads=%d sends=%d result=%+v err=%v", reads, sends.Load(), result, err)
			}
		})
	}
}
