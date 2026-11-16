package intelligence

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// canonicalObservation is a synthetic post in the exact form capture.py
// wrote: sorted keys, two-space indent, raw UTF-8 and <, > and & unescaped.
const canonicalObservation = `{
  "effect": {
    "detail": "Bytte til <wrapper>/data & lese via ply.",
    "kind": "changed_action"
  },
  "evidence": [
    {
      "basis": "direct",
      "detail": "Kommandoen ga «ok» og exit 0.",
      "ref": "ply/main 8f7f601"
    }
  ],
  "id": "obs-00000000-0000-4000-8000-000000000001",
  "interpretation": null,
  "limitations": "Én kjøring.",
  "links": [],
  "observation": "plyapp leste plans/ direkte i stedet for via ply.",
  "origin": {
    "activity": "PLAN-75/enforce",
    "project": "ply",
    "provider": "claude",
    "repo": "ply",
    "role": "coordinator",
    "session": null
  },
  "prior_refs": [],
  "recorded_at": "2026-11-16T10:00:00Z",
  "schema": "ply.observation/v1",
  "situation": "Kartlegging av direkte filkoblinger."
}
`

func testStore(t *testing.T) Store {
	t.Helper()
	return Open(filepath.Join(t.TempDir(), "intelligence"), "ply", []string{"ply", "plyapp-server"})
}

func requireClass(t *testing.T, err error, class string) {
	t.Helper()
	var e *Error
	if !errors.As(err, &e) || e.Class != class {
		t.Fatalf("error = %v, want class %s", err, class)
	}
}

func edited(t *testing.T, raw string, edit func(map[string]any)) []byte {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		t.Fatal(err)
	}
	edit(value)
	b, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestObserveKeepsCanonicalBytesAndIsIdempotent(t *testing.T) {
	s := testStore(t)
	receipt, err := s.Observe([]byte(canonicalObservation))
	if err != nil || receipt.State != StateRecorded || receipt.Kind != KindObservation {
		t.Fatalf("observe %+v %v", receipt, err)
	}
	stored, err := os.ReadFile(filepath.Join(s.dir, "observations", receipt.ID+".json"))
	if err != nil || string(stored) != canonicalObservation {
		t.Fatalf("stored bytes differ:\n%s", stored)
	}
	// The same post in another layout is the same content.
	compact := edited(t, canonicalObservation, func(map[string]any) {})
	again, err := s.Observe(compact)
	if err != nil || again.State != StateAlreadyRecorded || again.SHA256 != receipt.SHA256 {
		t.Fatalf("retry %+v %v", again, err)
	}
	changed := edited(t, canonicalObservation, func(v map[string]any) { v["limitations"] = "To kjøringer." })
	_, err = s.Observe(changed)
	requireClass(t, err, ClassConflict)
	shown, err := s.Show(receipt.ID)
	if err != nil || string(shown.Post) != canonicalObservation || shown.SHA256 != receipt.SHA256 {
		t.Fatalf("show %+v %v", shown, err)
	}
}

func TestObserveValidatesOriginAndLinks(t *testing.T) {
	s := testStore(t)
	for name, edit := range map[string]func(map[string]any){
		"unknown repo":    func(v map[string]any) { v["origin"].(map[string]any)["repo"] = "trip" },
		"other project":   func(v map[string]any) { v["origin"].(map[string]any)["project"] = "trip" },
		"unknown field":   func(v map[string]any) { v["severity"] = "high" },
		"empty situation": func(v map[string]any) { v["situation"] = " " },
		"missing post": func(v map[string]any) {
			v["links"] = []any{map[string]any{"relation": "supports", "id": "obs-00000000-0000-4000-8000-000000000009"}}
		},
		"self link": func(v map[string]any) {
			v["links"] = []any{map[string]any{"relation": "corrects", "id": v["id"]}}
		},
		"uses observation": func(v map[string]any) {
			v["links"] = []any{map[string]any{"relation": "uses", "id": "obs-00000000-0000-4000-8000-000000000009"}}
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := s.Observe(edited(t, canonicalObservation, edit))
			requireClass(t, err, ClassInvalidPost)
		})
	}
	first, err := s.Observe([]byte(canonicalObservation))
	if err != nil {
		t.Fatal(err)
	}
	linked := edited(t, canonicalObservation, func(v map[string]any) {
		v["id"] = "obs-00000000-0000-4000-8000-000000000002"
		v["links"] = []any{
			map[string]any{"relation": "supports", "id": first.ID},
			map[string]any{"relation": "uses", "id": "k-20261115-arkivet-naar-ikke-agentene"},
		}
	})
	if _, err := s.Observe(linked); err != nil {
		t.Fatalf("valid links: %v", err)
	}
}

func review(t *testing.T, edit func(map[string]any)) []byte {
	t.Helper()
	draft, err := Draft(KindCodeHealthReview, Origin{Project: "ply", Repo: "ply", Provider: "claude", Role: "author", Activity: "ply:1e2a460..5b89025"}, time.Date(2026, 11, 16, 10, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	return edited(t, string(draft), func(v map[string]any) {
		v["base"] = strings.Repeat("a", 40)
		v["head"] = strings.Repeat("b", 40)
		v["tool"] = "ply-code-quality version main-ce7694c"
		v["command"] = "plycq review --format=diff"
		v["exit"] = 0
		v["reviewer"] = map[string]any{"provider": "claude", "model": "unknown", "task": nil}
		v["first_review_question"] = "Is the new branch covered?"
		v["ratings"] = []any{
			map[string]any{"file": "cmd/x.go", "line": 12, "function": "run", "marks": []any{"HIGH 15->20"}, "rating": "used", "note": "Split the branch chain.", "effect": "changed_action"},
			map[string]any{"file": "cmd/x_test.go", "line": 3, "function": "TestRun", "marks": []any{"bumpy introduced"}, "rating": "irrelevant", "note": "Test setup loop.", "effect": "none"},
		}
		v["signals"], v["used"], v["irrelevant"], v["wrong"] = 2, 1, 1, 0
		v["counting"] = "code-health-review/v2"
		v["diff"] = "diff output\n"
		if edit != nil {
			edit(v)
		}
	})
}

func TestCodeHealthReviewCountsAndSize(t *testing.T) {
	s := testStore(t)
	receipt, err := s.Observe(review(t, nil))
	if err != nil || receipt.Kind != KindCodeHealthReview {
		t.Fatalf("observe review %+v %v", receipt, err)
	}
	for name, edit := range map[string]func(map[string]any){
		"count mismatch": func(v map[string]any) { v["used"], v["irrelevant"] = 2, 0 },
		"signals":        func(v map[string]any) { v["signals"] = 3 },
		"same commits":   func(v map[string]any) { v["head"] = v["base"] },
		"gate exit":      func(v map[string]any) { v["exit"] = 1 },
		"second role":    func(v map[string]any) { v["role"] = "author" },
		"review role":    func(v map[string]any) { v["origin"].(map[string]any)["role"] = "coordinator" },
		"too large":      func(v map[string]any) { v["diff"] = strings.Repeat("x", 4<<20) },
	} {
		t.Run(name, func(t *testing.T) {
			_, err := s.Observe(review(t, func(v map[string]any) {
				v["id"] = "obs-00000000-0000-4000-8000-0000000000aa"
				edit(v)
			}))
			requireClass(t, err, ClassInvalidPost)
		})
	}
}

func TestListFiltersAndAfter(t *testing.T) {
	s := testStore(t)
	first, _ := s.Observe([]byte(canonicalObservation))
	second, err := s.Observe(edited(t, canonicalObservation, func(v map[string]any) {
		v["id"] = "obs-00000000-0000-4000-8000-000000000002"
		v["recorded_at"] = "2026-11-16T11:00:00Z"
		v["origin"].(map[string]any)["repo"] = "plyapp-server"
		v["links"] = []any{map[string]any{"relation": "same_case", "id": first.ID}}
	}))
	if err != nil {
		t.Fatal(err)
	}
	reviewed, err := s.Observe(review(t, nil))
	if err != nil {
		t.Fatal(err)
	}
	ids := func(f Filter) []string {
		t.Helper()
		result, err := s.List(f)
		if err != nil {
			t.Fatal(err)
		}
		out := []string{}
		for _, item := range result.Items {
			if item.SHA256 == "" || len(item.Post) == 0 {
				t.Fatalf("item without post or hash: %+v", item)
			}
			out = append(out, item.ID)
		}
		return out
	}
	if got := ids(Filter{}); strings.Join(got, ",") != first.ID+","+reviewed.ID+","+second.ID {
		t.Fatalf("all %v", got)
	}
	if got := ids(Filter{After: first.ID}); strings.Join(got, ",") != reviewed.ID+","+second.ID {
		t.Fatalf("after %v", got)
	}
	if got := ids(Filter{Kind: KindCodeHealthReview}); len(got) != 1 || got[0] != reviewed.ID {
		t.Fatalf("kind %v", got)
	}
	if got := ids(Filter{Repo: "plyapp-server"}); len(got) != 1 || got[0] != second.ID {
		t.Fatalf("repo %v", got)
	}
	if got := ids(Filter{LinksTo: first.ID}); len(got) != 1 || got[0] != second.ID {
		t.Fatalf("links-to %v", got)
	}
	if got := ids(Filter{Since: "2026-11-16T10:30:00Z"}); strings.Join(got, ",") != second.ID {
		t.Fatalf("since %v", got)
	}
	if _, err := s.List(Filter{After: "obs-00000000-0000-4000-8000-0000000000ff"}); err == nil {
		t.Fatal("unknown --after accepted")
	}
}
