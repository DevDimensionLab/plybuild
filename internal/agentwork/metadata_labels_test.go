package agentwork

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMetadataLabelsRefreshAndExactJoin(t *testing.T) {
	h, root := tabProcessFixture(t)
	var rows []any
	for n, w := range []string{"w1", "w1", "w2"} {
		r := metadataRow(fmt.Sprintf("p%d", n))
		r["workspace_id"] = w
		rows = append(rows, r)
	}
	putFixture(t, root, "agents", metadataRaw(map[string]any{"result": map[string]any{"type": "agent_list", "agents": rows}}))
	for _, name := range []string{"cmd", "1", "é漢字 🌱", "<script>", ""} {
		putFixture(t, root, "w1", labelEnvelope(labelRow("w1", "t", name), labelRow("w2", "t", "Wrong")))
		putFixture(t, root, "w2", labelEnvelope(labelRow("w2", "t", "Same")))
		got := StatusMetadata(h)
		if got.Sources[0].Coverage != "complete" || len(got.Agents) != 3 {
			t.Fatal(got)
		}
		for i, a := range got.Agents {
			want := name
			if i == 2 {
				want = "Same"
			}
			if a.TabLabel != want || a.Pane != fmt.Sprintf("p%d", i) || a.Status != "unknown" {
				t.Fatal(a, want)
			}
		}
	}
	for _, call := range fixtureCalls(t, root) {
		if call != "agents" && call != "w1" && call != "w2" {
			t.Fatal(call)
		}
	}
}

func TestMetadataLabelsPublicPolicy(t *testing.T) {
	h, root := tabProcessFixture(t)
	putFixture(t, root, "agents", `{"result":{"type":"agent_list","agents":[{"agent":"codex","pane_id":"p","tab_id":"t","workspace_id":"w","agent_status":"idle"}]}}`)
	for _, label := range []string{"cmd", "1", "é漢字", "<script>", strings.Repeat("🌱", 128), strings.Repeat("a", 256), "", " ", " leading", "trailing ", "a/b", `a\b`, "a@b", "https:", "https://host", "SESSION_ID x", "Terminal_ID x", "CWD=x", "TOKEN=x", "SECRET=x", "PASSWORD=x", "x\n", "x\u007f", "x\u009b", "x\u202e", "x\u200d", "x\u2028", "x\u2029", "x\ufffd", strings.Repeat("🌱", 128) + "x", strings.Repeat("a", 257)} {
		t.Run(fmt.Sprintf("%q", label), func(t *testing.T) {
			putFixture(t, root, "w", labelEnvelope(labelRow("w", "t", label)))
			got := StatusMetadata(h)
			accepted := label == "cmd" || label == "1" || label == "é漢字" || label == "<script>" || label == "https:" || label == strings.Repeat("🌱", 128) || label == strings.Repeat("a", 256)
			want := ""
			if accepted {
				want = label
			}
			if len(got.Agents) != 1 || got.Sources[0].Coverage != "complete" || got.Agents[0].TabLabel != want {
				t.Fatal(got)
			}
			b, _ := json.Marshal(got.Agents[0])
			var fields map[string]any
			_ = json.Unmarshal(b, &fields)
			if _, ok := fields["tab_label"]; ok != accepted {
				t.Fatal("optional field", string(b))
			}
		})
	}
}

func TestMetadataLabelsAmbiguousAndMalformed(t *testing.T) {
	h, root := tabProcessFixture(t)
	putFixture(t, root, "agents", `{"result":{"type":"agent_list","agents":[{"agent":"codex","pane_id":"p","tab_id":"t","workspace_id":"w","agent_status":"idle"}]}}`)
	good := `{"workspace_id":"w","tab_id":"t","label":"Good"}`
	for _, bad := range []string{
		`{"workspace_id":"w","tab_id":"t","label":42}`,
		`{"workspace_id":"w","tab_id":"t","label":"\\ud800"}`,
		`{"workspace_id":"w","tab_id":"t","label":"\ud800"}`,
		`{"workspace_id":"w","tab_id":"t","label":42,"label":"Hidden"}`,
		`{"workspace_id":"w","tab_id":"t","tab_id":"other","label":"Hidden"}`,
		`{"workspace_id":"w","workspace_id":"other","tab_id":"t","label":"Hidden"}`,
	} {
		for _, reverse := range []bool{false, true} {
			rows := []string{good, bad}
			if reverse {
				rows[0], rows[1] = rows[1], rows[0]
			}
			putFixture(t, root, "w", `{"result":{"tabs":[`+strings.Join(rows, ",")+`]}}`)
			got := StatusMetadata(h)
			if len(got.Agents) != 1 || got.Agents[0].TabLabel != "" || got.Sources[0].Coverage != "complete" {
				t.Fatal(got)
			}
		}
	}
}

func TestMetadataSharedBudgetCancellationAndRecovery(t *testing.T) {
	agents := `{"result":{"type":"agent_list","agents":[{"agent":"codex","pane_id":"p","tab_id":"t","workspace_id":"w","agent_status":"idle"}]}}`
	parent, exhaust := context.WithCancel(context.Background())
	var mu sync.Mutex
	var deadlines []time.Time
	stalled := true
	h := fakeHerdr(func(ctx context.Context, args []string) (string, error) {
		deadline, ok := ctx.Deadline()
		mu.Lock()
		deadlines = append(deadlines, deadline)
		mu.Unlock()
		if !ok {
			return "", fmt.Errorf("call without a budget: %v", args)
		}
		if args[0] == "agent" {
			return agents, nil
		}
		if stalled {
			// The label read uses up what is left of the shared budget.
			exhaust()
			<-ctx.Done()
			return "", ctx.Err()
		}
		return labelEnvelope(labelRow("w", "t", "Known")), nil
	})
	got := StatusMetadataContext(parent, h)
	if len(got.Agents) != 1 || got.Agents[0].TabLabel != "" || got.Sources[0].Status != "available" {
		t.Fatal("exhausted label read", got)
	}
	if len(deadlines) != 2 || !deadlines[0].Equal(deadlines[1]) {
		t.Fatal("label read did not share the request budget", deadlines)
	}
	got = StatusMetadataContext(parent, h)
	if got.Sources[0].Status != "unavailable" || len(got.Agents) != 0 {
		t.Fatal("cancelled request", got)
	}
	stalled = false
	got = StatusMetadata(h)
	if len(got.Agents) != 1 || got.Agents[0].TabLabel != "Known" {
		t.Fatal("recovery", got)
	}
}

func TestMetadataWorkspaceLimitPreservesMembership(t *testing.T) {
	var rows []any
	labels := map[string]string{}
	for w := 33; w >= 0; w-- {
		id := fmt.Sprintf("w%02d", w)
		r := metadataRow(fmt.Sprintf("p%02d", w))
		r["workspace_id"] = id
		rows = append(rows, r)
		labels[id] = labelEnvelope(labelRow(id, "t", "Known"))
	}
	agents := metadataRaw(map[string]any{"result": map[string]any{"type": "agent_list", "agents": rows}})
	var mu sync.Mutex
	var calls []string
	h := fakeHerdr(func(_ context.Context, args []string) (string, error) {
		mu.Lock()
		defer mu.Unlock()
		if args[0] == "agent" {
			calls = append(calls, "agents")
			return agents, nil
		}
		calls = append(calls, args[3])
		return labels[args[3]], nil
	})
	got := StatusMetadata(h)
	if len(got.Agents) != 34 || got.Sources[0].Coverage != "complete" {
		t.Fatal(got)
	}
	for i, a := range got.Agents {
		want := "Known"
		if i >= 32 {
			want = ""
		}
		if a.TabLabel != want {
			t.Fatal(i, a)
		}
	}
	if len(calls) != 33 || calls[0] != "agents" {
		t.Fatal(calls)
	}
}
