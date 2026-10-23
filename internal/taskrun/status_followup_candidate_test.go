package taskrun

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// This runs only the exact caller-supplied candidate. The Herdr subprocess is a
// retained disposable test child; no real agent or production state is used.
func TestStatusFollowupCandidateJourney(t *testing.T) {
	binary := os.Getenv("PLY_CANDIDATE_BINARY")
	if binary == "" {
		t.Skip("exact candidate is supplied by acceptance")
	}
	f := newDeliveryFixture(t, "codex")
	s := deliveryInventoryState(t, f, "status-followup-fixture")
	s.Result.Transport.TerminalID = "term-private-status-fixture"
	deliveryInventoryReport(t, f, &s, "working", "working")
	// R1: known target B has damaged historical event bytes in the same Task.
	// It must remain diagnostic without vetoing the independent current target A.
	foreign := deliveryInventoryState(t, f, "foreign-target-B")
	deliveryInventoryReport(t, f, &foreign, "foreign-report", "working")
	if e := os.WriteFile(foreign.Result.Delivery.Events[0].Binding.Locator, []byte("corrupt historical fixture event"), 0600); e != nil {
		t.Fatal(e)
	}
	binDir := filepath.Join(f.R.WorkspaceRoot, "fake-bin")
	if e := os.MkdirAll(binDir, 0700); e != nil {
		t.Fatal(e)
	}
	countPath := filepath.Join(binDir, "prompts")
	agent := map[string]any{"workspace_id": s.Result.Transport.WorkspaceID, "tab_id": s.Result.Transport.TabID, "pane_id": s.Result.Transport.PaneID, "terminal_id": s.Result.Transport.TerminalID, "agent": "codex", "agent_status": "done", "interactive_ready": true, "agent_session": map[string]string{"agent": "codex", "kind": "id", "value": s.Result.Transport.AgentSessionID}, "title": "synthetic-private-title", "cwd": "/synthetic/private/path", "tokens": map[string]string{"token": "synthetic-private-token"}}
	a, _ := json.Marshal(agent)
	script := "#!/usr/bin/env python3\nimport json,sys\na=json.loads(" + string(mustJSON(string(a))) + ")\nu=dict(a);u['pane_id']='unbound-pane';u['terminal_id']='unbound-private-terminal';u['agent_session']={'agent':'codex','kind':'id','value':'unbound-private-session'}\nif sys.argv[1:3]==['agent','list']: r={'type':'agent_list','agents':[a,u]}\nelif sys.argv[1:3]==['agent','get']: r={'type':'agent_info','agent':u if sys.argv[3]=='unbound-pane' else a}\nelif sys.argv[1:3]==['agent','prompt']:\n with open(" + string(mustJSON(countPath)) + ",'a') as f: f.write('one\\n')\n r={'type':'agent_prompted','agent':a}\nelse: sys.exit(2)\nprint(json.dumps({'result':r}))\n"
	if e := os.WriteFile(filepath.Join(binDir, "herdr"), []byte(script), 0700); e != nil {
		t.Fatal(e)
	}
	run := func(args ...string) map[string]any {
		t.Helper()
		c := exec.Command(binary, append(args, "--workspace", f.R.WorkspaceRoot, "--format", "json")...)
		c.Dir = f.R.WorkspaceRoot
		c.Env = append(os.Environ(), "PATH="+binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
		b, e := c.CombinedOutput()
		if e != nil {
			t.Fatalf("candidate %v: %v\n%s", args, e, b)
		}
		var o map[string]any
		if e = json.Unmarshal(b, &o); e != nil {
			t.Fatalf("candidate JSON: %s", b)
		}
		for _, secret := range []string{"synthetic-private", s.Result.Transport.AgentSessionID, s.Result.Transport.TerminalID} {
			if strings.Contains(string(b), secret) {
				t.Fatal("candidate leaked private metadata")
			}
		}
		return o
	}
	overview := run("workflow", "agents")
	agents := overview["agents"].([]any)
	if len(agents) != 2 {
		t.Fatal(overview)
	}
	row := agents[0].(map[string]any)
	var unbound map[string]any
	for _, entry := range agents {
		r := entry.(map[string]any)
		if r["work"] != nil {
			row = r
		} else {
			unbound = r
		}
	}
	if row["work"] == nil {
		t.Fatal("verified native binding missing", overview)
	}
	work := row["work"].(map[string]any)
	if work["phase"] != "working" {
		t.Fatal("completed turn changed native phase", work)
	}
	target := row["agent_ref"].(string)
	write := func(name string, v any) string {
		t.Helper()
		p := filepath.Join(binDir, name)
		if e := os.WriteFile(p, mustJSON(v), 0600); e != nil {
			t.Fatal(e)
		}
		return p
	}
	draft := map[string]any{"publication_key": "fixture-status", "owner_ref": "fixture-owner", "mandate_ref": "fixture-mandate", "registrar": map[string]string{"kind": "human_requested", "actor_ref": "fixture-human"}, "target_ref": target, "wait_kind": "status", "wait_ref": target, "next_check_at_utc": "2000-01-01T00:00:00Z"}
	p := write("register.json", draft)
	followup := run("workflow", "followup", "register", "--file", p)
	id := followup["id"].(string)
	// Recover a dropped register response by list + the exact canonical tuple.
	listed := run("workflow", "followup", "list")
	if len(listed["items"].([]any)) != 1 {
		t.Fatal(listed)
	}
	again := run("workflow", "followup", "register", "--file", p)
	if again["id"] != id {
		t.Fatal("register replay changed identity")
	}
	sent := run("workflow", "followup", "dispatch", id)
	if sent["transport"].(map[string]any)["state"] != "sent" || sent["recheck"] != nil {
		t.Fatal("ACK not separate from receipt", sent)
	}
	run("workflow", "followup", "dispatch", id)
	b, e := os.ReadFile(countPath)
	if e != nil || string(b) != "one\n" {
		t.Fatal("cycle replay", string(b), e)
	}
	mutation := map[string]any{"request_id": "fixture-fulfilled", "owner_ref": "fixture-owner", "expected_revision": sent["revision"], "generation": sent["generation"], "target_ref": target, "outcome": "fulfilled", "done": "Checked actual fixture sources", "next": "No further status check", "blocker": ""}
	p = write("recheck.json", mutation)
	done := run("workflow", "followup", "recheck", id, "--file", p)
	if done["state"] != "fulfilled" {
		t.Fatal(done)
	}
	readback := run("workflow", "followup", "request-status", id, "--request-id", "fixture-fulfilled")
	if readback["revision"] != done["revision"] {
		t.Fatal("lost reply readback changed result")
	}
	run("workflow", "followup", "dispatch", id)
	b, _ = os.ReadFile(countPath)
	if string(b) != "one\n" {
		t.Fatal("closed cycle replay")
	}
	ut := unbound["agent_ref"].(string)
	draft["publication_key"] = "unbound-self-status"
	draft["target_ref"] = ut
	draft["wait_ref"] = ut
	draft["registrar"] = map[string]string{"kind": "agent_self_registered", "actor_ref": ut}
	up := write("unbound-register.json", draft)
	uf := run("workflow", "followup", "register", "--file", up)
	if uf["work"] != nil {
		t.Fatal("invented unbound work")
	}
	uf = run("workflow", "followup", "dispatch", uf["id"].(string))
	if uf["transport"].(map[string]any)["state"] != "sent" {
		t.Fatal("complete unbound veto coverage cannot send", uf)
	}
	overview = run("workflow", "agents")
	for _, entry := range overview["agents"].([]any) {
		r := entry.(map[string]any)
		if r["work"] != nil {
			work = r["work"].(map[string]any)
		}
	}
	if work["phase"] != "working" {
		t.Fatal("self report changed native work")
	}
}
func mustJSON(v any) []byte {
	b, e := json.Marshal(v)
	if e != nil {
		panic(e)
	}
	return b
}
