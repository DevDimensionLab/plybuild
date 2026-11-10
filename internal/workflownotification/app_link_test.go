package workflownotification

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

const fixtureAppLink = "https://plyapp-pwa.web.app"

// A request-v6 is v4, optionally with the v5 Herdr origin, plus one owner-selected app link.
func appLinkFixture(t *testing.T, herdr bool) (*fixture, map[string]any) {
	t.Helper()
	f, m := documentlessFixture(t)
	m["schema_version"] = 6
	presentation := m["presentation"].(map[string]any)
	presentation["app_link"] = fixtureAppLink
	if herdr {
		presentation["herdr"] = map[string]any{"workspace": "ply-master", "tab": "spec"}
	}
	writeJSON(t, f.input.File, m)
	return f, m
}

func TestOlderPayloadsKeepTheirFrozenBytes(t *testing.T) {
	for name, build := range map[string]func(*testing.T) *fixture{
		"v4": func(t *testing.T) *fixture { f, _ := documentlessFixture(t); return f },
		"v5": func(t *testing.T) *fixture { f, _ := herdrFixture(t, "ply-master", "spec"); return f },
	} {
		t.Run(name, func(t *testing.T) {
			p := build(t).preview(t)
			b, _ := json.Marshal(p.Payload)
			if strings.Contains(string(b), "blocks") || p.PayloadSHA != frozenPayloadSHA[name] {
				t.Fatal(p.PayloadSHA, string(b))
			}
		})
	}
}

// Recorded from the code before request-v6.
var frozenPayloadSHA = map[string]string{
	"v4": "sha256:c431d9aabdf0005f39271190770f81a7989ce553d56e67b6f18188e37866bc52",
	"v5": "sha256:8b367cb798c290eb7908c19ade45b4197ff50cfd183ec8ad744a0109adb1bb3f",
}

func TestAppLinkIsAnIsolatedPlyBuiltBlock(t *testing.T) {
	for herdr, first := range map[bool]string{false: "16:00:00 · claude", true: "16:00:00 · claude · ply-master › spec"} {
		f, _ := appLinkFixture(t, herdr)
		var body []byte
		f.deps.Transport = func(_ context.Context, _ string, b []byte) Observation {
			f.posts++
			body = b
			return Observation{Status: 200, Body: []byte("ok"), Dispatch: "started"}
		}
		p := f.preview(t)
		text := first + "\nDone: Varsler — Ferdig.\nNext: No action required."
		want := `{"text":` + quote(text) + `,"mrkdwn":false,"parse":"none","unfurl_links":false,"unfurl_media":false,"blocks":[` +
			`{"type":"section","text":{"type":"plain_text","text":` + quote(text) + `,"emoji":false}},` +
			`{"type":"context","elements":[{"type":"mrkdwn","text":` + quote("<"+fixtureAppLink+"|Åpne plyapp>") + `,"verbatim":true}]}]}`
		if out, err := f.apply(t, p); err != nil || out.(Result).State != "transport_acknowledged" || string(body) != want {
			t.Fatalf("%v %v\n%s\n%s", out, err, body, want)
		}
		f.input.Operation, f.input.ID = "show", p.ID
		if out, err := Execute(f.deps, f.input); err != nil || out.(Result).Freshness != "current" || f.posts != 1 {
			t.Fatal("stored v6 record does not read back", out, err)
		}
	}
}
func quote(s string) string { b, _ := json.Marshal(s); return string(b) }

// Agent text stays plain text; only the validated link reaches the mrkdwn element.
func TestAgentTextNeverReachesTheLinkElement(t *testing.T) {
	f, m := appLinkFixture(t, false)
	m["public"].(map[string]any)["task_title"] = "*fet* _x_ ~s~ `c` &lt;!here&gt; :+1:"
	writeJSON(t, f.input.File, m)
	p := f.preview(t)
	if len(p.Payload.Blocks) != 2 || !strings.Contains(p.Payload.Blocks[0].Text.Text, "*fet* _x_ ~s~ `c` &lt;!here&gt; :+1:") ||
		p.Payload.Blocks[0].Text.Type != "plain_text" || p.Payload.Blocks[1].Elements[0].Text != "<"+fixtureAppLink+"|Åpne plyapp>" {
		t.Fatalf("%+v", p.Payload.Blocks)
	}
	for _, unsafe := range []string{"<https://evil.example|Åpne plyapp>", "Hei @here"} {
		m["public"].(map[string]any)["task_title"] = unsafe
		writeJSON(t, f.input.File, m)
		if _, err := Execute(f.deps, f.input); exitOf(err) != 2 {
			t.Fatal(unsafe, err)
		}
	}
}
func TestAppLinkIsValidatedStrictly(t *testing.T) {
	for _, link := range []string{"", "http://plyapp-pwa.web.app", "https://plyapp-pwa.web.app:443", "https://u@plyapp-pwa.web.app",
		"https://plyapp-pwa.web.app?x=1", "https://plyapp-pwa.web.app#x", "https://plyapp-pwa.web.app/a|b", "https://plyapp-pwa.web.app/>",
		"https://Plyapp-pwa.web.app", "https://localhost", "https://-x.web.app", " https://plyapp-pwa.web.app", "https://plyapp-pwa.web.app//x",
		"https://hooks.slack.com/services/T0/B0/x", "https://" + strings.Repeat("a", 190) + ".web.app"} {
		f, m := appLinkFixture(t, false)
		m["presentation"].(map[string]any)["app_link"] = link
		writeJSON(t, f.input.File, m)
		if _, err := Execute(f.deps, f.input); exitOf(err) != 2 {
			t.Fatal(link, err)
		}
	}
	for _, link := range []string{"https://plyapp-pwa.web.app/", "https://ply.example/app/v1.2_x~y-z"} {
		f, m := appLinkFixture(t, false)
		m["presentation"].(map[string]any)["app_link"] = link
		writeJSON(t, f.input.File, m)
		f.preview(t)
	}
	for name, change := range map[string]func(m map[string]any){
		"v6 without link": func(m map[string]any) { delete(m["presentation"].(map[string]any), "app_link") },
		"v4 with link":    func(m map[string]any) { m["schema_version"] = 4 },
		"v5 with link": func(m map[string]any) {
			m["schema_version"] = 5
			m["presentation"].(map[string]any)["herdr"] = map[string]any{"workspace": "w", "tab": "t"}
		},
		"incomplete herdr": func(m map[string]any) { m["presentation"].(map[string]any)["herdr"] = map[string]any{"tab": "t"} },
		"link not string":  func(m map[string]any) { m["presentation"].(map[string]any)["app_link"] = []string{fixtureAppLink} },
	} {
		f, m := appLinkFixture(t, false)
		change(m)
		writeJSON(t, f.input.File, m)
		if _, err := Execute(f.deps, f.input); exitOf(err) != 2 {
			t.Fatal(name, err)
		}
	}
}

// v4, v5 and v6 for one event share the identity, so adding the link never sends twice.
func TestAppLinkSharesTheEventIdentity(t *testing.T) {
	f, v6 := appLinkFixture(t, false)
	raw, _ := os.ReadFile(f.input.File)
	var v4 map[string]any
	if err := json.Unmarshal(raw, &v4); err != nil {
		t.Fatal(err)
	}
	v4["schema_version"] = 4
	delete(v4["presentation"].(map[string]any), "app_link")
	writeJSON(t, f.input.File, v4)
	p := f.preview(t)
	if out, _ := f.apply(t, p); out.(Result).State != "transport_acknowledged" {
		t.Fatal(out)
	}
	writeJSON(t, f.input.File, v6)
	out, err := Execute(f.deps, f.input)
	if r, ok := out.(Result); exitOf(err) != 4 || !ok || r.ID != p.ID || f.posts != 1 {
		t.Fatal(out, err, f.posts)
	}
	in := f.input
	in.Apply, in.Confirm = true, *p.Confirmation
	if _, err = Execute(f.deps, in); exitOf(err) != 4 || f.posts != 1 {
		t.Fatal("second apply sent", err, f.posts)
	}
}
