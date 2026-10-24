package agentwork

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

// Count real inventory invocations at their existing filesystem boundary. This
// avoids adding production hooks or replacing native/workflowstatus composition.
type inventoryCountingFiles struct {
	workspace.FileSystem
	root  string
	reads *int
}

func (f inventoryCountingFiles) Getwd() (string, error) { return f.root, nil }
func (f inventoryCountingFiles) ReadFile(path string) ([]byte, error) {
	if path == filepath.Join(f.root, workspace.MarkerDirectory, workspace.MarkerFile) {
		pcs := make([]uintptr, 40)
		frames := runtime.CallersFrames(pcs[:runtime.Callers(2, pcs)])
		for {
			frame, more := frames.Next()
			if strings.HasSuffix(frame.Function, "/taskrun.ReadInventory") {
				*f.reads++
				break
			}
			if !more {
				break
			}
		}
	}
	return f.FileSystem.ReadFile(path)
}
func TestNativeReusesWorkflowStatusInventory(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	d := workspace.SystemDependencies()
	count := 0
	d.Files = inventoryCountingFiles{d.Files, root, &count}
	if _, err := workspace.Init(d); err != nil {
		t.Fatal(err)
	}
	for request := 1; request <= 2; request++ {
		n := readNative(d, root)
		if !n.Available || n.Partial || !n.VetoCoverageComplete {
			t.Fatalf("empty valid native source lost: %+v", n)
		}
		if count != request {
			t.Errorf("request %d inventory reads=%d; want total %d", request, count, request)
		}
	}
	// New incomplete coverage on a later read remains visible to both native
	// binding and the source projection; the captured empty inventory is gone.
	requests := filepath.Join(root, ".ply", "workflow-runs", "requests")
	if err := os.MkdirAll(requests, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(requests, "unrecognized.json"), []byte("damaged"), 0600); err != nil {
		t.Fatal(err)
	}
	n := readNative(d, root)
	if count != 3 || !n.Available || !n.Partial || !n.CoverageUnknown || n.VetoCoverageComplete {
		t.Fatalf("changed coverage became fresh/clear: reads=%d native=%+v", count, n)
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
				basis.Vetoes = map[string]string{s.Runtime.(*fakeRuntime).agents[0].Identity.sessionKey(): "native_question"}
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
