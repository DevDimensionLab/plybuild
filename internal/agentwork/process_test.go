package agentwork

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// A retained test child exercises real process locks and process death. It never
// invokes Herdr or a provider; its only transport appends a local fixture marker.
func TestFollowupDispatchProcessChild(t *testing.T) {
	root := os.Getenv("PLY_FOLLOWUP_TEST_ROOT")
	if root == "" {
		t.Skip("test-child entrypoint")
	}
	now := time.Date(2026, 10, 23, 13, 0, 0, 0, time.UTC)
	i := Identity{Provider: "codex", Session: "session-fixture", Workspace: "w-fixture", Tab: "tab-fixture", Pane: "pane-fixture", Terminal: "terminal-fixture"}
	s := &Service{Root: root, Now: func() time.Time { return now }, Runtime: &fakeRuntime{agents: []Observation{{Identity: i, Status: "idle", Ready: boolp(true)}}}, Native: func() NativeSnapshot {
		return NativeSnapshot{Available: true, VetoCoverageComplete: true, Vetoes: map[string]string{}}
	}, Send: func(Identity, string) error {
		f, e := os.OpenFile(filepath.Join(root, "input-count"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if e != nil {
			return e
		}
		defer f.Close()
		_, e = f.WriteString("input\n")
		if e != nil {
			return e
		}
		return f.Sync()
	}}
	if os.Getenv("PLY_FOLLOWUP_TEST_CRASH") == "after_claim" {
		s.Fault = func(stage string) error {
			if stage == "after_claim" {
				os.Exit(23)
			}
			return nil
		}
	}
	if _, e := s.Dispatch(os.Getenv("PLY_FOLLOWUP_TEST_ID")); e != nil {
		t.Fatal(e)
	}
}
func TestFollowupConcurrentProcessesAndProcessDeath(t *testing.T) {
	binary, e := os.Executable()
	if e != nil {
		t.Fatal(e)
	}
	child := func(root, id, crash string) *exec.Cmd {
		c := exec.Command(binary, "-test.run=^TestFollowupDispatchProcessChild$", "-test.timeout=15s")
		c.Env = append(os.Environ(), "PLY_FOLLOWUP_TEST_ROOT="+root, "PLY_FOLLOWUP_TEST_ID="+id, "PLY_FOLLOWUP_TEST_CRASH="+crash)
		return c
	}
	for _, crash := range []bool{false, true} {
		t.Run(fmt.Sprint(crash), func(t *testing.T) {
			s, d, _ := fixture(t)
			d.NextCheckAtUTC = s.now()
			f, e := s.Register(d)
			if e != nil {
				t.Fatal(e)
			}
			if crash {
				c := child(s.Root, f.ID, "after_claim")
				if e = c.Run(); e == nil || c.ProcessState.ExitCode() != 23 {
					t.Fatal("crash boundary not reached", e)
				}
			}
			a, b := child(s.Root, f.ID, ""), child(s.Root, f.ID, "")
			if e = a.Start(); e != nil {
				t.Fatal(e)
			}
			if e = b.Start(); e != nil {
				_ = a.Wait()
				t.Fatal(e)
			}
			if e = a.Wait(); e != nil {
				t.Fatal(e)
			}
			if e = b.Wait(); e != nil {
				t.Fatal(e)
			}
			bytes, e := os.ReadFile(filepath.Join(s.Root, "input-count"))
			got, _ := s.Show(f.ID)
			if crash {
				if !os.IsNotExist(e) || got.Transport.State != "unknown" {
					t.Fatal("process death replayed claim", string(bytes), e, got.Transport)
				}
			} else {
				if e != nil || strings.Count(string(bytes), "input\n") != 1 || got.Transport.State != "sent" {
					t.Fatal("multiple process inputs", string(bytes), e)
				}
			}
		})
	}
}
func TestFollowupMalformedStorePreservesOtherStatusSources(t *testing.T) {
	s, d, _ := fixture(t)
	f, e := s.Register(d)
	if e != nil {
		t.Fatal(e)
	}
	db, e := s.read()
	if e != nil {
		t.Fatal(e)
	}
	p := filepath.Join(s.dir(), "revisions", fmt.Sprintf("%012d.json", db.Sequence))
	if e = os.WriteFile(p, []byte("{truncated"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Show(f.ID); e == nil {
		t.Fatal("malformed source accepted")
	}
	v, e := s.Overview()
	if e != nil || len(v.Agents) != 1 || v.Sources[2].Status != "unavailable" {
		b, _ := json.Marshal(v)
		t.Fatal(string(b), e)
	}
}
