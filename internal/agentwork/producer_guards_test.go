package agentwork

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func qualifyTestCallback(t *testing.T, root, path string) {
	t.Helper()
	path, e := filepath.EvalSymlinks(path)
	if e != nil {
		t.Fatal(e)
	}
	b, e := os.ReadFile(path)
	if e != nil {
		t.Fatal(e)
	}
	raw, _ := json.Marshal(map[string]any{"path": path, "sha256": fmt.Sprintf("sha256:%x", sha256.Sum256(b)), "qualification_ref": "TA2E037F9-tests"})
	dir := root
	if e = os.MkdirAll(dir, 0700); e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(dir, "ply-callback.local.json"), raw, 0600); e != nil {
		t.Fatal(e)
	}
}
func privateAnswerValues() []string {
	return []string{"/Users/me", "/uSeRs/me", "/private/file", "/home/me", "/tmp/build", "~/repo", `~\repo`, "xoxb-value", "xoxa-value", "xoxp-value", "xoxr-value", "xoxs-value", "xapp-value", "TOKEN=value", "password=value", "secret=value", "run 12345678-abcd-abcd-abcd-123456789abc"}
}
func TestProducerAnswerPrivacyBeforeStorage(t *testing.T) {
	for _, field := range []string{"summary", "done", "next", "blocker"} {
		for _, value := range privateAnswerValues() {
			t.Run(field+value, func(t *testing.T) {
				s, d, _ := fixture(t)
				f, e := s.Register(d)
				if e != nil {
					t.Fatal(e)
				}
				before, e := s.read()
				if e != nil {
					t.Fatal(e)
				}
				m := mutation(f, "answer")
				m.Outcome = "fulfilled"
				m.Summary = "TA2E037F9 implementing guards"
				switch field {
				case "summary":
					m.Summary = value
				case "done":
					m.Done = value
				case "next":
					m.Next = value
				case "blocker":
					m.Blocker = value
				}
				if _, e = s.Mutate(f.ID, "recheck", m); e == nil || e.Error() != "followup_invalid_recheck" {
					t.Fatalf("private answer accepted: %v", e)
				}
				after, e := s.read()
				if e != nil || !reflect.DeepEqual(before, after) {
					t.Fatal("rejected answer changed store", e)
				}
				if _, e = s.RequestStatus(f.ID, m.RequestID); e == nil {
					t.Fatal("rejected request cached")
				}
				m.Summary = "TA2E037F9 still implementing"
				m.Done = "Updated internal/agentwork"
				m.Next = "Review"
				m.Blocker = "none"
				got, e := s.Mutate(f.ID, "recheck", m)
				if e != nil || got.State != "fulfilled" || got.Cycle != f.Cycle || len(got.History) != 0 || got.Transport.State != f.Transport.State {
					t.Fatal(got, e)
				}
				again, e := s.Mutate(f.ID, "recheck", m)
				if e != nil || !reflect.DeepEqual(got, again) {
					t.Fatal("idempotence", e)
				}
				m.Summary = "changed"
				if _, e = s.Mutate(f.ID, "recheck", m); e == nil {
					t.Fatal("request conflict lost")
				}
			})
		}
	}
}
func TestProducerBothLabelProjections(t *testing.T) {
	labels := append(privateAnswerValues(), "", " ", "bad\nlabel", strings.Repeat("😀", 129), "Vanlig agentnavn", strings.Repeat("😀", 128))
	for _, label := range labels {
		t.Run(label, func(t *testing.T) {
			s, _, _ := fixture(t)
			p := filepath.Join(t.TempDir(), "label")
			writeLabelResponse(t, p, label)
			h := runtimeScript(t, `case "$1 $2" in
 "agent list") printf '%s' '{"result":{"type":"agent_list","agents":[{"agent":"codex","agent_status":"idle","workspace_id":"w-fixture","tab_id":"tab-fixture","pane_id":"pane-fixture","terminal_id":"terminal-fixture","agent_session":{"agent":"codex","kind":"id","value":"session-fixture"}}]}}';;
 "tab list") cat "`+p+`";; *) exit 9;; esac`)
			s.Runtime = h
			rich, e := s.Overview()
			if e != nil || len(rich.Agents) != 1 {
				t.Fatal(rich, e)
			}
			meta := StatusMetadata(h)
			if len(meta.Agents) != 1 {
				t.Fatal(meta)
			}
			want := ""
			if label == "Vanlig agentnavn" || label == strings.Repeat("😀", 128) {
				want = label
			}
			for _, got := range []string{rich.Agents[0].TabLabel, meta.Agents[0].TabLabel} {
				if got != want {
					t.Errorf("label %q want %q", got, want)
				}
			}
			raw := metadataRaw(map[string]any{"result": map[string]any{"type": "agent_list", "agents": []any{metadataRow("p")}}})
			if got := metadataFixture(t, raw); len(got.Agents) != 1 {
				t.Fatal(got)
			}
		})
	}
}
func TestProducerQualifiedCallbackBeforeClaim(t *testing.T) {
	for _, mode := range []string{"missing-config", "missing", "changed", "relative", "no-evidence", "not-executable", "valid"} {
		t.Run(mode, func(t *testing.T) {
			s, d, calls := fixture(t)
			d.NextCheckAtUTC = s.now()
			f, e := s.Register(d)
			if e != nil {
				t.Fatal(e)
			}
			path := filepath.Join(t.TempDir(), "qualified ply")
			if e = os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0700); e != nil {
				t.Fatal(e)
			}
			path, e = filepath.EvalSymlinks(path)
			if e != nil {
				t.Fatal(e)
			}
			qualifyTestCallback(t, s.Root, path)
			switch mode {
			case "missing-config":
				os.Remove(filepath.Join(s.Root, "ply-callback.local.json"))
			case "relative", "no-evidence":
				configPath := filepath.Join(s.Root, "ply-callback.local.json")
				raw, _ := os.ReadFile(configPath)
				var pin map[string]any
				json.Unmarshal(raw, &pin)
				if mode == "relative" {
					pin["path"] = "ply"
				} else {
					delete(pin, "qualification_ref")
				}
				raw, _ = json.Marshal(pin)
				os.WriteFile(configPath, raw, 0600)
			case "not-executable":
				os.Chmod(path, 0600)
			case "missing":
				os.Remove(path)
			case "changed":
				os.WriteFile(path, []byte("#!/bin/sh\nexit 1\n"), 0700)
			}
			old := t.TempDir()
			os.WriteFile(filepath.Join(old, "ply"), []byte("#!/bin/sh\nexit 99\n"), 0700)
			t.Setenv("PATH", old)
			var message string
			s.Send = func(_ Identity, m string) error { calls.Add(1); message = m; return nil }
			before, _ := s.read()
			got, e := s.Dispatch(f.ID)
			if mode != "valid" {
				after, _ := s.read()
				if e == nil || calls.Load() != 0 || !reflect.DeepEqual(before, after) {
					t.Fatal("unqualified executor consumed cycle", got, e)
				}
				return
			}
			if e != nil || calls.Load() != 1 {
				t.Fatal(got, e)
			}
			if !strings.HasPrefix(message, "Status question from the user:") {
				t.Error("useful question must come first", message)
			}
			for _, text := range []string{path, "outcome=fulfilled", "UTF-16", "summary", "done", "next", "blocker", "generation", "cycle", "waiting", "no authority"} {
				if !strings.Contains(message, text) {
					t.Errorf("missing %q", text)
				}
			}
			if bytes.Contains([]byte(message), []byte("`ply workflow")) {
				t.Fatal("PATH callback")
			}
		})
	}
}

func TestProducerAnswerTextBoundaries(t *testing.T) {
	for _, field := range []string{"summary", "done", "next", "blocker"} {
		for _, v := range []string{strings.Repeat("😀", 128), strings.Repeat("😀", 128) + "x", "bad\ntext"} {
			t.Run(field+fmt.Sprint(len(v)), func(t *testing.T) {
				s, d, _ := fixture(t)
				f, e := s.Register(d)
				if e != nil {
					t.Fatal(e)
				}
				m := mutation(f, "boundary")
				m.Outcome = "fulfilled"
				switch field {
				case "summary":
					m.Summary = v
				case "done":
					m.Done = v
				case "next":
					m.Next = v
				case "blocker":
					m.Blocker = v
				}
				_, e = s.Mutate(f.ID, "recheck", m)
				if (e == nil) != (v == strings.Repeat("😀", 128)) {
					t.Fatal("UTF-16/control boundary", e)
				}
			})
		}
	}
}

func TestProducerHistoricalIdempotentAnswerUnchanged(t *testing.T) {
	s, d, _ := fixture(t)
	f, e := s.Register(d)
	if e != nil {
		t.Fatal(e)
	}
	m := mutation(f, "historical")
	m.Outcome = "fulfilled"
	m.Summary = "old /Users/private"
	old := f
	old.State = "fulfilled"
	old.Revision++
	old.Recheck = &Recheck{Summary: m.Summary, Outcome: m.Outcome, ReporterRef: f.TargetRef, Cycle: f.Cycle}
	e = s.locked(func(db *database) error {
		r := db.Records[f.ID]
		r.Value = old
		r.Requests[m.RequestID] = operation{Digest: hash(struct {
			Op       string
			Mutation Mutation
		}{"recheck", m}), Result: old}
		db.Records[f.ID] = r
		return s.save(db)
	})
	if e != nil {
		t.Fatal(e)
	}
	before, _ := s.read()
	got, e := s.Mutate(f.ID, "recheck", m)
	after, _ := s.read()
	if e != nil || got.Recheck.Summary != m.Summary || !reflect.DeepEqual(before, after) {
		t.Fatal("historical idempotent result was rewritten", got, e)
	}
}
