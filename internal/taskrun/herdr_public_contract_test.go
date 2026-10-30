package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// This synthetic public Herdr fixture has no private startup capability or
// operations. It exercises Ply consumers, not a live provider or atomic writer.
func publicHerdrFixture(t *testing.T, f *workflowFixture) {
	t.Helper()
	workflowTestModel(t, *f, map[string]any{"public_contract_only": true})
}

func assertPublicHerdrCalls(t *testing.T, f workflowFixture, prompts int) {
	t.Helper()
	b, err := os.ReadFile(f.Calls)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, line := range bytes.Split(bytes.TrimSpace(b), []byte("\n")) {
		var call struct{ Argv []string }
		if err = json.Unmarshal(line, &call); err != nil {
			t.Fatal(err)
		}
		argv := strings.Join(call.Argv, " ")
		if strings.Contains(argv, "agent startup") || strings.Contains(argv, "tui.status_line") {
			t.Fatalf("private startup requirement: %s", argv)
		}
		if len(call.Argv) >= 2 && call.Argv[0] == "agent" && call.Argv[1] == "prompt" {
			count++
		}
	}
	if count != prompts {
		t.Fatalf("public prompt count = %d, want %d", count, prompts)
	}
}

func TestPublicHerdrOwnerStartup(t *testing.T) {
	f := workflowTestFixture(t)
	publicHerdrFixture(t, &f)
	o, err := workflowReadinessStart(t, f)
	if err != nil {
		t.Fatal(err)
	}
	if o.Transport.AgentSessionID != "fixture-session" || o.Transport.State != "working" {
		t.Fatalf("owner not bound and prompted: %+v", o.Transport)
	}
	assertPublicHerdrCalls(t, f, 1)
	workflowAssertNoReplay(t, f, o)
}
