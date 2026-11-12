package plan

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
	"github.com/devdimensionlab/plybuild/internal/filelock"
)

// specExample is the plan of the JSON readback example in the Planverk v1 spec.
const specExample = `{
  "title": "Planverk i Ply", "goal": "Agent og menneske leser samme plan",
  "scope": "Ply CLI og lagring; ikke server/PWA", "owner": "ply-epic",
  "selection": "selected", "status": "active",
  "summary": "Behov samlet; spec skrives", "next_action": "Ferdigstill spec",
  "refs": ["backlog:20261111-c97f7518"],
  "steps": [
    {"id":"spec", "title":"Avklar kontrakt", "outcome":"Implementerbar spec", "owner":"ply-epic", "status":"active", "next_action":"Review spec"},
    {"id":"implement", "title":"Implementer", "outcome":"Testet kandidat", "owner":"Claude", "status":"planned", "next_action":"Start valgt Assignment", "depends_on":["spec"]},
    {"id":"review", "title":"Kontroller", "outcome":"Review og CLI-pilot grønne", "owner":"ply-epic", "status":"planned", "next_action":"Review kandidat", "depends_on":["implement"]},
    {"id":"deliver", "title":"Lever lokalt", "outcome":"Avstemt masterleveranse", "owner":"ply-master", "status":"planned", "next_action":"Integrer kvalifisert kandidat", "depends_on":["review"]}
  ]
}`

func testStore(t *testing.T) (Store, string) {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "plans")
	store := Open(dir)
	clock := time.Date(2026, 11, 12, 7, 0, 0, 0, time.UTC)
	var mu sync.Mutex
	store.now = func() time.Time {
		mu.Lock()
		defer mu.Unlock()
		clock = clock.Add(time.Second)
		return clock
	}
	return store, dir
}

func requireClass(t *testing.T, err error, class string) *Error {
	t.Helper()
	var e *Error
	if !errors.As(err, &e) || e.Class != class {
		t.Fatalf("error %v, want class %s", err, class)
	}
	return e
}

// edit returns specExample with change applied to its decoded form.
func edit(t *testing.T, change func(map[string]any)) []byte {
	t.Helper()
	var document map[string]any
	if err := json.Unmarshal([]byte(specExample), &document); err != nil {
		t.Fatal(err)
	}
	change(document)
	encoded, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

func step(document map[string]any, i int) map[string]any {
	return document["steps"].([]any)[i].(map[string]any)
}

func encodeView(t *testing.T, v View) string {
	t.Helper()
	encoded, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}

func TestAddShowUpdateHistoryKeepEveryRevision(t *testing.T) {
	store, dir := testStore(t)
	added, err := store.Add([]byte(specExample), "ply-epic", "Brukeren valgte planverk i Ply")
	if err != nil {
		t.Fatal(err)
	}
	if added.ID != "PLAN-1" || added.Revision != 1 || added.CreatedAt != "2026-11-12T07:00:01Z" || added.UpdatedAt != added.CreatedAt ||
		added.Actor != "ply-epic" || added.Reason != "Brukeren valgte planverk i Ply" || added.NeedsHuman {
		t.Fatalf("added %+v", added)
	}
	if want := (Progress{Total: 4, Planned: 3, Active: 1}); added.Progress != want {
		t.Fatalf("progress %+v", added.Progress)
	}
	if got := stepIDs(added.Plan); !reflect.DeepEqual(got, []string{"spec", "implement", "review", "deliver"}) {
		t.Fatalf("steps %v", got)
	}
	shown, err := store.Show("PLAN-1", 0)
	if err != nil || encodeView(t, shown) != encodeView(t, added) {
		t.Fatalf("show %+v %v", shown, err)
	}
	firstBytes := encodeView(t, shown)

	next := edit(t, func(d map[string]any) {
		step(d, 0)["status"] = "done"
		step(d, 0)["result"] = "Spec peerreviewet"
		step(d, 1)["status"] = "active"
		d["summary"] = "Implementasjon pågår"
	})
	updated, err := store.Update("PLAN-1", 1, next, "ply-epic", "Spec ferdig; implementasjon starter")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Revision != 2 || updated.CreatedAt != added.CreatedAt || updated.UpdatedAt != "2026-11-12T07:00:02Z" ||
		updated.Plan.Summary != "Implementasjon pågår" || updated.Plan.Steps[0].Result != "Spec peerreviewet" {
		t.Fatalf("updated %+v", updated)
	}
	if got := stepIDs(updated.Plan); !reflect.DeepEqual(got, []string{"spec", "implement", "review", "deliver"}) {
		t.Fatalf("step order after status change %v", got)
	}
	if want := (Progress{Total: 4, Planned: 2, Active: 1, Done: 1}); updated.Progress != want {
		t.Fatalf("progress %+v", updated.Progress)
	}
	// An identical snapshot is still a new, deliberate revision.
	same, err := store.Update("PLAN-1", 2, next, "ply-epic", "Bekreftet uendret")
	if err != nil || same.Revision != 3 || !reflect.DeepEqual(same.Plan, updated.Plan) {
		t.Fatalf("identical update %+v %v", same, err)
	}

	old, err := store.Show("PLAN-1", 1)
	if err != nil || encodeView(t, old) != firstBytes {
		t.Fatalf("revision 1 changed after updates: %v\n%s\n%s", err, encodeView(t, old), firstBytes)
	}
	latest, err := store.Show("PLAN-1", 0)
	if err != nil || encodeView(t, latest) != encodeView(t, same) {
		t.Fatalf("latest %+v %v", latest, err)
	}
	if _, err := store.Show("PLAN-1", 4); err == nil {
		t.Fatal("future revision shown")
	} else {
		requireClass(t, err, ClassNotFound)
	}

	history, err := store.History("PLAN-1")
	want := []HistoryEntry{
		{Revision: 1, At: "2026-11-12T07:00:01Z", Actor: "ply-epic", Reason: "Brukeren valgte planverk i Ply"},
		{Revision: 2, At: "2026-11-12T07:00:02Z", Actor: "ply-epic", Reason: "Spec ferdig; implementasjon starter"},
		{Revision: 3, At: "2026-11-12T07:00:03Z", Actor: "ply-epic", Reason: "Bekreftet uendret"},
	}
	if err != nil || !reflect.DeepEqual(history, want) {
		t.Fatalf("history %+v %v", history, err)
	}

	// One file holds the plan and its whole history; nothing else is stored.
	names := dirNames(t, dir)
	if !reflect.DeepEqual(names, []string{".lock", "PLAN-1.json"}) {
		t.Fatalf("stored files %v", names)
	}
	var file struct {
		SchemaVersion int               `json:"schema_version"`
		ID            string            `json:"id"`
		Revisions     []json.RawMessage `json:"revisions"`
	}
	contents, err := os.ReadFile(filepath.Join(dir, "PLAN-1.json"))
	if err != nil || json.Unmarshal(contents, &file) != nil || file.SchemaVersion != 1 || file.ID != "PLAN-1" || len(file.Revisions) != 3 {
		t.Fatalf("file %s %v", contents, err)
	}
	if bytes.Contains(contents, []byte("needs_human")) || bytes.Contains(contents, []byte("progress")) {
		t.Fatalf("derived values stored: %s", contents)
	}
}

func TestAddNumbersAboveHighestFileNameIncludingCorruptFiles(t *testing.T) {
	store, dir := testStore(t)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "PLAN-7.json"), []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	added, err := store.Add([]byte(specExample), "a", "r")
	if err != nil || added.ID != "PLAN-8" {
		t.Fatalf("add %+v %v", added, err)
	}
	if got, _ := os.ReadFile(filepath.Join(dir, "PLAN-7.json")); string(got) != "not json" {
		t.Fatalf("corrupt file changed: %q", got)
	}
}

func TestAddRefusesExhaustedIDSpaceWithoutOverwriting(t *testing.T) {
	store, dir := testStore(t)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	last := filepath.Join(dir, "PLAN-999999999.json")
	if err := os.WriteFile(last, []byte("kept"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := store.Add([]byte(specExample), "a", "r")
	requireClass(t, err, ClassFull)
	if got, _ := os.ReadFile(last); string(got) != "kept" || !reflect.DeepEqual(dirNames(t, dir), []string{".lock", "PLAN-999999999.json"}) {
		t.Fatalf("files changed: %v", dirNames(t, dir))
	}
}

func TestUpdateWithStaleRevisionConflictsWithoutWriting(t *testing.T) {
	store, dir := testStore(t)
	_, err := store.Update("PLAN-1", 1, []byte(specExample), "a", "r")
	requireClass(t, err, ClassNotFound)
	if _, err := os.Lstat(dir); !errors.Is(err, fs.ErrNotExist) {
		t.Fatal("update of a missing plan created the plans directory")
	}
	if _, err := store.Add([]byte(specExample), "a", "r"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Update("PLAN-1", 1, []byte(specExample), "a", "r2"); err != nil {
		t.Fatal(err)
	}
	before := readFile(t, filepath.Join(dir, "PLAN-1.json"))
	for _, revision := range []int{1, 3} {
		_, err := store.Update("PLAN-1", revision, []byte(specExample), "a", "stale")
		requireClass(t, err, ClassConflict)
	}
	if after := readFile(t, filepath.Join(dir, "PLAN-1.json")); after != before {
		t.Fatal("conflict changed the file")
	}
	_, err = store.Update("PLAN-1", 0, []byte(specExample), "a", "r")
	requireClass(t, err, ClassInvalidArguments)
	_, err = store.Update("PLAN-2", 1, []byte(specExample), "a", "r")
	requireClass(t, err, ClassNotFound)
}

func TestConcurrentUpdatesWithSameRevisionWriteExactlyOnce(t *testing.T) {
	store, dir := testStore(t)
	if _, err := store.Add([]byte(specExample), "a", "r"); err != nil {
		t.Fatal(err)
	}
	const writers = 16
	errs := make([]error, writers)
	snapshots := make([][]byte, writers)
	for i := range snapshots {
		snapshots[i] = edit(t, func(d map[string]any) { d["summary"] = fmt.Sprintf("writer %d", i) })
	}
	var start, done sync.WaitGroup
	start.Add(1)
	for i := 0; i < writers; i++ {
		done.Add(1)
		go func(i int) {
			defer done.Done()
			start.Wait()
			_, errs[i] = store.Update("PLAN-1", 1, snapshots[i], "a", "race")
		}(i)
	}
	start.Done()
	done.Wait()
	won := -1
	for i, err := range errs {
		if err == nil {
			if won >= 0 {
				t.Fatalf("writers %d and %d both won", won, i)
			}
			won = i
			continue
		}
		requireClass(t, err, ClassConflict)
	}
	latest, err := store.Show("PLAN-1", 0)
	if won < 0 || err != nil || latest.Revision != 2 || latest.Plan.Summary != fmt.Sprintf("writer %d", won) {
		t.Fatalf("winner %d latest %+v %v", won, latest, err)
	}
	if names := dirNames(t, dir); !reflect.DeepEqual(names, []string{".lock", "PLAN-1.json"}) {
		t.Fatalf("files %v", names)
	}
}

func TestConcurrentAddsGetUniqueIDs(t *testing.T) {
	store, _ := testStore(t)
	const writers = 16
	ids := make([]string, writers)
	errs := make([]error, writers)
	var done sync.WaitGroup
	for i := 0; i < writers; i++ {
		done.Add(1)
		go func(i int) {
			defer done.Done()
			var v View
			v, errs[i] = store.Add([]byte(specExample), "a", "r")
			ids[i] = v.ID
		}(i)
	}
	done.Wait()
	seen := map[string]bool{}
	for i, id := range ids {
		if errs[i] != nil || seen[id] {
			t.Fatalf("add %d: %q %v", i, id, errs[i])
		}
		seen[id] = true
	}
	list, err := store.List("")
	if err != nil || len(list.Items) != writers || len(list.Invalid) != 0 {
		t.Fatalf("list %+v %v", list, err)
	}
}

func TestRejectsInvalidSnapshotsBeforeAnyEffect(t *testing.T) {
	long := strings.Repeat("x", 201)
	cases := map[string][]byte{
		"empty":              []byte(""),
		"not an object":      []byte(`[]`),
		"trailing JSON":      []byte(specExample + `{}`),
		"trailing text":      []byte(specExample + ` x`),
		"duplicate key":      []byte(strings.Replace(specExample, `"title": "Planverk i Ply",`, `"title": "Planverk i Ply", "title": "Annen",`, 1)),
		"duplicate step key": []byte(strings.Replace(specExample, `{"id":"spec",`, `{"id":"spec", "id":"spec2",`, 1)),
		"case-folded key":    []byte(strings.Replace(specExample, `"title": "Planverk`, `"Title": "Planverk`, 1)),
		"invalid UTF-8":      []byte(strings.Replace(specExample, "Planverk i Ply", "Plan\xffverk", 1)),
		"too large":          append([]byte(specExample), bytes.Repeat([]byte(" "), 128*1024)...),
		"unknown field":      edit(t, func(d map[string]any) { d["percent"] = 50 }),
		"metadata field":     edit(t, func(d map[string]any) { d["id"] = "PLAN-1" }),
		"unknown step field": edit(t, func(d map[string]any) { step(d, 0)["estimate"] = "2d" }),
		"unknown attention field": edit(t, func(d map[string]any) {
			d["attention"] = map[string]any{"reason": "r", "who": "w", "needs_human": true, "since": "x"}
		}),
		"null value":                    edit(t, func(d map[string]any) { d["details"] = nil }),
		"missing title":                 edit(t, func(d map[string]any) { delete(d, "title") }),
		"blank goal":                    edit(t, func(d map[string]any) { d["goal"] = "  " }),
		"blank summary":                 edit(t, func(d map[string]any) { d["summary"] = "" }),
		"multi-line title":              edit(t, func(d map[string]any) { d["title"] = "a\nb" }),
		"long title":                    edit(t, func(d map[string]any) { d["title"] = long }),
		"long owner":                    edit(t, func(d map[string]any) { d["owner"] = long }),
		"long summary":                  edit(t, func(d map[string]any) { d["summary"] = strings.Repeat("y", 501) }),
		"NUL":                           edit(t, func(d map[string]any) { d["details"] = "a\x00b" }),
		"terminal escape":               edit(t, func(d map[string]any) { d["details"] = "a\x1b[2Jb" }),
		"carriage return":               edit(t, func(d map[string]any) { d["details"] = "a\rb" }),
		"bad selection":                 edit(t, func(d map[string]any) { d["selection"] = "maybe" }),
		"bad status":                    edit(t, func(d map[string]any) { d["status"] = "paused" }),
		"bad step status":               edit(t, func(d map[string]any) { step(d, 1)["status"] = "started" }),
		"waiting plan, no attention":    edit(t, func(d map[string]any) { d["status"] = "waiting" }),
		"blocked step, no attention":    edit(t, func(d map[string]any) { step(d, 1)["status"] = "blocked" }),
		"unknown step, no attention":    edit(t, func(d map[string]any) { step(d, 1)["status"] = "unknown" }),
		"attention without who":         edit(t, func(d map[string]any) { d["attention"] = map[string]any{"reason": "r", "who": "", "needs_human": true} }),
		"attention without needs_human": edit(t, func(d map[string]any) { d["attention"] = map[string]any{"reason": "r", "who": "w"} }),
		"done plan, open step":          edit(t, func(d map[string]any) { d["status"] = "done" }),
		"no steps":                      edit(t, func(d map[string]any) { d["steps"] = []any{} }),
		"too many steps":                edit(t, func(d map[string]any) { d["steps"] = manySteps(129) }),
		"bad step id":                   edit(t, func(d map[string]any) { step(d, 0)["id"] = "Spec" }),
		"duplicate step id":             edit(t, func(d map[string]any) { step(d, 1)["id"] = "spec"; delete(step(d, 1), "depends_on") }),
		"blank step outcome":            edit(t, func(d map[string]any) { step(d, 0)["outcome"] = "" }),
		"multi-line step owner":         edit(t, func(d map[string]any) { step(d, 0)["owner"] = "a\nb" }),
		"missing dependency":            edit(t, func(d map[string]any) { step(d, 1)["depends_on"] = []any{"nowhere"} }),
		"self dependency":               edit(t, func(d map[string]any) { step(d, 1)["depends_on"] = []any{"implement"} }),
		"repeated dependency":           edit(t, func(d map[string]any) { step(d, 1)["depends_on"] = []any{"spec", "spec"} }),
		"dependency cycle":              edit(t, func(d map[string]any) { step(d, 0)["depends_on"] = []any{"deliver"} }),
		"blank ref":                     edit(t, func(d map[string]any) { d["refs"] = []any{""} }),
		"too many refs":                 edit(t, func(d map[string]any) { d["refs"] = manyStrings(129) }),
		"too many decisions":            edit(t, func(d map[string]any) { d["decisions"] = manyStrings(129) }),
		"wrong type":                    edit(t, func(d map[string]any) { d["steps"] = "spec" }),
	}
	for name, snapshot := range cases {
		snapshot := snapshot
		t.Run(name, func(t *testing.T) {
			store, dir := testStore(t)
			_, err := store.Add(snapshot, "ply-epic", "r")
			requireClass(t, err, ClassInvalidArguments)
			if _, statErr := os.Lstat(dir); !errors.Is(statErr, fs.ErrNotExist) {
				t.Fatalf("add with invalid snapshot created %s", dir)
			}
		})
	}
	t.Run("update keeps file", func(t *testing.T) {
		store, dir := testStore(t)
		if _, err := store.Add([]byte(specExample), "a", "r"); err != nil {
			t.Fatal(err)
		}
		before := readFile(t, filepath.Join(dir, "PLAN-1.json"))
		for name, snapshot := range cases {
			_, err := store.Update("PLAN-1", 1, snapshot, "a", "r")
			if e := requireClass(t, err, ClassInvalidArguments); strings.Contains(e.Detail, "Planverk") {
				t.Fatalf("%s: detail echoes input: %s", name, e.Detail)
			}
		}
		if readFile(t, filepath.Join(dir, "PLAN-1.json")) != before {
			t.Fatal("invalid update changed the file")
		}
	})
}

func TestAcceptsOptionalFieldsAndAttentionOnAnyStatus(t *testing.T) {
	store, _ := testStore(t)
	snapshot := edit(t, func(d map[string]any) {
		d["details"] = "Linje 1\n\tLinje 2"
		d["decisions"] = []any{"Kun Ply først"}
		d["learning"] = "Les spec helt"
		d["next_action"] = ""
		d["attention"] = map[string]any{"reason": "Åpent produktvalg", "who": "bruker", "needs_human": false}
		step(d, 0)["refs"] = []any{"commit:abc123"}
		step(d, 0)["attention"] = map[string]any{"reason": "Review venter", "who": "ply-epic", "needs_human": false}
		step(d, 0)["next_action"] = ""
	})
	v, err := store.Add(snapshot, "a", "r")
	if err != nil {
		t.Fatal(err)
	}
	if v.Plan.Attention == nil || v.Plan.Attention.Who != "bruker" || v.Plan.Details != "Linje 1\n\tLinje 2" || v.Plan.Steps[0].Attention == nil || v.NeedsHuman {
		t.Fatalf("view %+v", v)
	}
	all := edit(t, func(d map[string]any) {
		d["status"] = "done"
		for i := 0; i < 4; i++ {
			step(d, i)["status"] = "done"
		}
	})
	if v, err := store.Update(v.ID, 1, all, "a", "ferdig"); err != nil || v.Progress != (Progress{Total: 4, Done: 4}) {
		t.Fatalf("done plan %+v %v", v, err)
	}
	if _, err := store.Add(edit(t, func(d map[string]any) { d["steps"] = manySteps(128) }), "a", "r"); err != nil {
		t.Fatalf("128 steps: %v", err)
	}
}

func TestNeedsHumanIsDerivedFromPlanOrAnyStep(t *testing.T) {
	store, _ := testStore(t)
	stepAsks := edit(t, func(d map[string]any) {
		step(d, 1)["status"] = "waiting"
		step(d, 1)["attention"] = map[string]any{"reason": "Velg lagringsmodell", "who": "bruker", "needs_human": true}
	})
	v, err := store.Add(stepAsks, "a", "r")
	if err != nil || !v.NeedsHuman || v.Plan.Status != "active" || v.Plan.Attention != nil || v.Progress != (Progress{Total: 4, Planned: 2, Active: 1, Waiting: 1}) {
		t.Fatalf("step attention %+v %v", v, err)
	}
	planAsks := edit(t, func(d map[string]any) {
		d["status"] = "blocked"
		d["attention"] = map[string]any{"reason": "Mangler tilgang", "who": "ply-master", "needs_human": true}
	})
	if v, err := store.Add(planAsks, "a", "r"); err != nil || !v.NeedsHuman {
		t.Fatalf("plan attention %+v %v", v, err)
	}
	if v, err := store.Add([]byte(specExample), "a", "r"); err != nil || v.NeedsHuman {
		t.Fatalf("no attention %+v %v", v, err)
	}
	list, err := store.List("")
	if err != nil || len(list.Items) != 3 {
		t.Fatalf("list %+v %v", list, err)
	}
	needs := map[string]bool{}
	for _, item := range list.Items {
		needs[item.ID] = item.NeedsHuman
	}
	if !reflect.DeepEqual(needs, map[string]bool{"PLAN-1": true, "PLAN-2": true, "PLAN-3": false}) {
		t.Fatalf("list needs_human %v", needs)
	}
}

func TestListSortsByUpdateAndKeepsInvalidUnderFilter(t *testing.T) {
	store, dir := testStore(t)
	for i := 0; i < 3; i++ {
		if _, err := store.Add([]byte(specExample), "a", "r"); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := store.Update("PLAN-1", 1, edit(t, func(d map[string]any) { d["status"] = "planned" }), "a", "r"); err != nil {
		t.Fatal(err)
	}
	// PLAN-2 and PLAN-3 share an update time; the ID breaks the tie.
	frozen := store.now()
	store.now = func() time.Time { return frozen }
	for _, id := range []string{"PLAN-3", "PLAN-2"} {
		if _, err := store.Update(id, 1, []byte(specExample), "a", "r"); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "PLAN-9.json"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "notes.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	list, err := store.List("")
	if err != nil {
		t.Fatal(err)
	}
	if got := itemIDs(list); !reflect.DeepEqual(got, []string{"PLAN-2", "PLAN-3", "PLAN-1"}) {
		t.Fatalf("order %v", got)
	}
	wantInvalid := []InvalidFile{{File: "PLAN-9.json", Error: "invalid JSON"}, {File: "notes.json", Error: "invalid file name"}}
	if !reflect.DeepEqual(list.Invalid, wantInvalid) {
		t.Fatalf("invalid %+v", list.Invalid)
	}
	item := list.Items[2]
	if item.Status != "planned" || item.Revision != 2 || item.Title != "Planverk i Ply" || item.Selection != "selected" || item.Owner != "ply-epic" ||
		item.Summary == "" || item.NextAction != "Ferdigstill spec" || item.Progress.Total != 4 || item.UpdatedAt == "" {
		t.Fatalf("item %+v", item)
	}
	filtered, err := store.List("planned")
	if err != nil || !reflect.DeepEqual(itemIDs(filtered), []string{"PLAN-1"}) || !reflect.DeepEqual(filtered.Invalid, wantInvalid) {
		t.Fatalf("filtered %+v %v", filtered, err)
	}
	_, err = store.List("paused")
	requireClass(t, err, ClassInvalidArguments)
}

func TestReadsCreateAndChangeNothing(t *testing.T) {
	store, dir := testStore(t)
	list, err := store.List("")
	if err != nil || len(list.Items) != 0 || list.Items == nil || list.Invalid == nil {
		t.Fatalf("empty list %+v %v", list, err)
	}
	_, err = store.Show("PLAN-1", 0)
	requireClass(t, err, ClassNotFound)
	_, err = store.History("PLAN-1")
	requireClass(t, err, ClassNotFound)
	if _, err := os.Lstat(dir); !errors.Is(err, fs.ErrNotExist) {
		t.Fatal("reads created the plans directory")
	}

	if _, err := store.Add([]byte(specExample), "a", "r"); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(dir, ".lock")); err != nil {
		t.Fatal(err)
	}
	before := treeState(t, dir)
	_, _ = store.List("")
	_, _ = store.List("active")
	_, _ = store.Show("PLAN-1", 0)
	_, _ = store.Show("PLAN-1", 1)
	_, _ = store.Show("PLAN-2", 0)
	_, _ = store.History("PLAN-1")
	if after := treeState(t, dir); !reflect.DeepEqual(after, before) {
		t.Fatalf("reads changed files:\n%v\n%v", before, after)
	}
}

func TestRejectsInvalidIDsAndArguments(t *testing.T) {
	store, _ := testStore(t)
	for _, id := range []string{"", "PLAN-0", "PLAN-01", "plan-1", "PLAN-1000000000", "PLAN-1.json", "../PLAN-1", "PLAN-1/x"} {
		_, err := store.Show(id, 0)
		requireClass(t, err, ClassInvalidArguments)
		_, err = store.History(id)
		requireClass(t, err, ClassInvalidArguments)
		_, err = store.Update(id, 1, []byte(specExample), "a", "r")
		requireClass(t, err, ClassInvalidArguments)
	}
	_, err := store.Show("PLAN-1", -1)
	requireClass(t, err, ClassInvalidArguments)
	for name, args := range map[string][2]string{
		"blank actor":      {" ", "r"},
		"multi-line actor": {"a\nb", "r"},
		"long actor":       {strings.Repeat("a", 201), "r"},
		"blank reason":     {"a", ""},
		"long reason":      {"a", strings.Repeat("r", 1001)},
		"control reason":   {"a", "r\x07"},
	} {
		_, err := store.Add([]byte(specExample), args[0], args[1])
		if e := requireClass(t, err, ClassInvalidArguments); e == nil {
			t.Fatal(name)
		}
	}
	if _, err := Open("relative/plans").Add([]byte(specExample), "a", "r"); err == nil {
		t.Fatal("relative directory accepted")
	}
}

func TestStoredFileProblemsAreInvalidAndNeverOverwritten(t *testing.T) {
	valid := func(t *testing.T) (Store, string, string) {
		store, dir := testStore(t)
		if _, err := store.Add([]byte(specExample), "a", "r"); err != nil {
			t.Fatal(err)
		}
		if _, err := store.Add([]byte(specExample), "a", "r"); err != nil {
			t.Fatal(err)
		}
		return store, dir, readFile(t, filepath.Join(dir, "PLAN-1.json"))
	}
	cases := map[string]func(string) string{
		"unknown schema": func(s string) string { return strings.Replace(s, `"schema_version": 1`, `"schema_version": 2`, 1) },
		"newer schema field": func(s string) string {
			return strings.Replace(s, `"schema_version": 1,`, `"schema_version": 2, "extra": true,`, 1)
		},
		"other id": func(s string) string { return strings.Replace(s, `"id": "PLAN-1"`, `"id": "PLAN-2"`, 1) },
		"duplicate key": func(s string) string {
			return strings.Replace(s, `"id": "PLAN-1",`, `"id": "PLAN-1", "id": "PLAN-1",`, 1)
		},
		"unknown field": func(s string) string { return strings.Replace(s, `"actor": "a",`, `"actor": "a", "host": "x",`, 1) },
		"trailing data": func(s string) string { return s + "{}" },
		"revision gap":  func(s string) string { return strings.Replace(s, `"revision": 1`, `"revision": 2`, 1) },
		"bad time": func(s string) string {
			return strings.Replace(s, `"at": "2026-11-12T07:00:01Z"`, `"at": "yesterday"`, 1)
		},
		"invalid plan":  func(s string) string { return strings.Replace(s, `"status": "active"`, `"status": "paused"`, 1) },
		"no revisions":  func(s string) string { return `{"schema_version": 1, "id": "PLAN-1", "revisions": []}` },
		"invalid UTF-8": func(s string) string { return strings.Replace(s, "Planverk", "Plan\xffverk", 1) },
	}
	for name, corrupt := range cases {
		corrupt := corrupt
		t.Run(name, func(t *testing.T) {
			store, dir, contents := valid(t)
			changed := corrupt(contents)
			if changed == contents {
				t.Fatal("corruption did not apply")
			}
			path := filepath.Join(dir, "PLAN-1.json")
			if err := os.WriteFile(path, []byte(changed), 0o644); err != nil {
				t.Fatal(err)
			}
			requireInvalidAndKept(t, store, path, changed)
		})
	}
	t.Run("symlink", func(t *testing.T) {
		store, dir, contents := valid(t)
		target := filepath.Join(t.TempDir(), "elsewhere.json")
		if err := os.WriteFile(target, []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, "PLAN-1.json")
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, path); err != nil {
			t.Fatal(err)
		}
		requireInvalidAndKept(t, store, target, contents)
		if info, err := os.Lstat(path); err != nil || info.Mode()&fs.ModeSymlink == 0 {
			t.Fatal("symlink replaced")
		}
	})
	t.Run("directory", func(t *testing.T) {
		store, dir, _ := valid(t)
		path := filepath.Join(dir, "PLAN-1.json")
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(path, 0o755); err != nil {
			t.Fatal(err)
		}
		_, err := store.Show("PLAN-1", 0)
		requireClass(t, err, ClassInvalid)
		_, err = store.Update("PLAN-1", 1, []byte(specExample), "a", "r")
		requireClass(t, err, ClassInvalid)
	})
}

func requireInvalidAndKept(t *testing.T, store Store, path, contents string) {
	t.Helper()
	_, err := store.Show("PLAN-1", 0)
	requireClass(t, err, ClassInvalid)
	_, err = store.History("PLAN-1")
	requireClass(t, err, ClassInvalid)
	_, err = store.Update("PLAN-1", 1, []byte(specExample), "a", "r")
	requireClass(t, err, ClassInvalid)
	if readFile(t, path) != contents {
		t.Fatal("invalid file was overwritten")
	}
	list, err := store.List("")
	if err != nil || !reflect.DeepEqual(itemIDs(list), []string{"PLAN-2"}) || len(list.Invalid) != 1 || list.Invalid[0].File != "PLAN-1.json" {
		t.Fatalf("list %+v %v", list, err)
	}
	if strings.Contains(list.Invalid[0].Error, string(filepath.Separator)) {
		t.Fatalf("invalid detail has a path: %s", list.Invalid[0].Error)
	}
}

func TestFileLimitRefusesUpdateAndKeepsFile(t *testing.T) {
	store, dir := testStore(t)
	big := edit(t, func(d map[string]any) { d["details"] = strings.Repeat("d", 120*1024) })
	v, err := store.Add(big, "a", "r")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "PLAN-1.json")
	for {
		before := readFile(t, path)
		next, err := store.Update(v.ID, v.Revision, big, "a", "r")
		if err != nil {
			requireClass(t, err, ClassFull)
			if readFile(t, path) != before {
				t.Fatal("full update changed the file")
			}
			if len(before) > maxFileBytes || v.Revision < 60 {
				t.Fatalf("refused at %d bytes, revision %d", len(before), v.Revision)
			}
			return
		}
		v = next
	}
}

func TestWriteFailuresDescribeTheirEffect(t *testing.T) {
	failure := func(op string, replaced bool) func(string, []byte, fs.FileMode) error {
		return func(string, []byte, fs.FileMode) error {
			return &atomicfile.Error{Op: op, Path: "/secret/plans/PLAN-1.json", Replaced: replaced, Err: errors.New("disk")}
		}
	}
	for name, test := range map[string]struct {
		replace func(string, []byte, fs.FileMode) error
		unknown bool
	}{
		"before replace":         {failure("write", false), false},
		"replace error":          {failure("replace", false), true},
		"directory sync failure": {failure("sync directory", true), true},
	} {
		test := test
		t.Run(name, func(t *testing.T) {
			store, _ := testStore(t)
			store.replace = test.replace
			_, err := store.Add([]byte(specExample), "a", "r")
			addErr := requireClass(t, err, ClassIOError)
			if strings.Contains(addErr.Detail, "secret") || !strings.Contains(addErr.Detail, "PLAN-1") {
				t.Fatalf("add detail %q", addErr.Detail)
			}
			if got := strings.Contains(addErr.Detail, "unknown"); got != test.unknown {
				t.Fatalf("add detail %q, unknown %v", addErr.Detail, test.unknown)
			}

			store.replace = atomicfile.Replace
			if _, err := store.Add([]byte(specExample), "a", "r"); err != nil {
				t.Fatal(err)
			}
			store.replace = test.replace
			_, err = store.Update("PLAN-1", 1, []byte(specExample), "a", "r")
			updateErr := requireClass(t, err, ClassIOError)
			if got := strings.Contains(updateErr.Detail, "unknown"); got != test.unknown || !strings.Contains(updateErr.Detail, "PLAN-1") {
				t.Fatalf("update detail %q", updateErr.Detail)
			}
		})
	}
	t.Run("lock failure", func(t *testing.T) {
		store, dir := testStore(t)
		store.lock = filelock.Options{Fail: func(op string) error {
			if op == "lock" {
				return errors.New("busy")
			}
			return nil
		}}
		_, err := store.Add([]byte(specExample), "a", "r")
		e := requireClass(t, err, ClassIOError)
		if !strings.Contains(e.Detail, "no plan file changed") {
			t.Fatalf("lock detail %q", e.Detail)
		}
		if names := dirNames(t, dir); !reflect.DeepEqual(names, []string{".lock"}) {
			t.Fatalf("files %v", names)
		}
	})
	t.Run("lock release failure after write", func(t *testing.T) {
		store, _ := testStore(t)
		store.lock = filelock.Options{Fail: func(op string) error {
			if op == "unlock" {
				return errors.New("unlock")
			}
			return nil
		}}
		_, err := store.Add([]byte(specExample), "a", "r")
		e := requireClass(t, err, ClassIOError)
		if !strings.Contains(e.Detail, "PLAN-1 revision 1 was written") {
			t.Fatalf("release detail %q", e.Detail)
		}
		store.lock = filelock.Options{}
		if v, err := store.Show("PLAN-1", 0); err != nil || v.Revision != 1 {
			t.Fatalf("readback %+v %v", v, err)
		}
	})
}

func stepIDs(p Plan) []string {
	var ids []string
	for _, s := range p.Steps {
		ids = append(ids, s.ID)
	}
	return ids
}

func itemIDs(list ListResult) []string {
	ids := []string{}
	for _, item := range list.Items {
		ids = append(ids, item.ID)
	}
	return ids
}

func manySteps(n int) []any {
	var steps []any
	for i := 0; i < n; i++ {
		steps = append(steps, map[string]any{"id": fmt.Sprintf("s%d", i), "title": "t", "outcome": "o", "owner": "w", "status": "planned", "next_action": ""})
	}
	return steps
}

func manyStrings(n int) []any {
	var values []any
	for i := 0; i < n; i++ {
		values = append(values, fmt.Sprintf("v%d", i))
	}
	return values
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(contents)
}

func dirNames(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	return names
}

func treeState(t *testing.T, dir string) map[string]string {
	t.Helper()
	state := map[string]string{}
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		contents := ""
		if info.Mode().IsRegular() {
			contents = readFile(t, path)
		}
		state[path] = fmt.Sprintf("%v %v %s", info.Mode(), info.ModTime().UnixNano(), contents)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return state
}

// planWriterEnv makes the test binary act as one writer process with its own
// Store, as separate agents do; see TestProcessesWriteThroughSeparateStores.
const planWriterEnv = "PLY_PLAN_TEST_WRITER"

func TestPlanWriterProcess(t *testing.T) {
	op := os.Getenv(planWriterEnv)
	if op == "" {
		t.Skip("helper process only")
	}
	store := Open(os.Getenv("PLY_PLAN_TEST_DIR"))
	var v View
	var err error
	if op == "add" {
		v, err = store.Add([]byte(specExample), "process", "parallel add")
	} else {
		v, err = store.Update("PLAN-1", 1, []byte(strings.Replace(specExample, "Behov samlet; spec skrives", op, 1)), "process", "race")
	}
	var e *Error
	switch {
	case err == nil:
		fmt.Printf("ok %s %d\n", v.ID, v.Revision)
	case errors.As(err, &e):
		fmt.Printf("error %s\n", e.Class)
	default:
		fmt.Printf("error %v\n", err)
	}
}

func TestProcessesWriteThroughSeparateStores(t *testing.T) {
	const processes = 8
	run := func(t *testing.T, dir string, ops []string) []string {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		outputs := make([]string, len(ops))
		errs := make([]error, len(ops))
		var done sync.WaitGroup
		for i, op := range ops {
			command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestPlanWriterProcess$", "-test.count=1")
			command.Env = append(os.Environ(), planWriterEnv+"="+op, "PLY_PLAN_TEST_DIR="+dir)
			done.Add(1)
			go func(i int, command *exec.Cmd) {
				defer done.Done()
				out, err := command.Output()
				outputs[i], errs[i] = string(out), err
			}(i, command)
		}
		done.Wait()
		results := make([]string, len(ops))
		for i := range ops {
			if errs[i] != nil {
				t.Fatalf("process %d: %v\n%s", i, errs[i], outputs[i])
			}
			results[i] = strings.SplitN(outputs[i], "\n", 2)[0]
		}
		return results
	}

	t.Run("same revision", func(t *testing.T) {
		store, dir := testStore(t)
		if _, err := store.Add([]byte(specExample), "a", "r"); err != nil {
			t.Fatal(err)
		}
		var ops []string
		for i := 0; i < processes; i++ {
			ops = append(ops, fmt.Sprintf("writer %d", i))
		}
		results := run(t, dir, ops)
		won := -1
		for i, result := range results {
			switch result {
			case "ok PLAN-1 2":
				if won >= 0 {
					t.Fatalf("processes %d and %d both won: %v", won, i, results)
				}
				won = i
			case "error " + ClassConflict:
			default:
				t.Fatalf("process %d: %q", i, result)
			}
		}
		latest, err := store.Show("PLAN-1", 0)
		if won < 0 || err != nil || latest.Revision != 2 || latest.Plan.Summary != ops[won] {
			t.Fatalf("winner %d latest %+v %v", won, latest, err)
		}
	})

	t.Run("parallel add", func(t *testing.T) {
		store, dir := testStore(t)
		ops := make([]string, processes)
		for i := range ops {
			ops[i] = "add"
		}
		results := run(t, dir, ops)
		seen := map[string]bool{}
		for _, result := range results {
			if !strings.HasPrefix(result, "ok PLAN-") || !strings.HasSuffix(result, " 1") || seen[result] {
				t.Fatalf("results %v", results)
			}
			seen[result] = true
		}
		if list, err := store.List(""); err != nil || len(list.Items) != processes || len(list.Invalid) != 0 {
			t.Fatalf("list %+v %v", list, err)
		}
	})
}

func TestFailureAfterRealReplaceLeavesExactlyOneNewRevision(t *testing.T) {
	store, dir := testStore(t)
	if _, err := store.Add([]byte(specExample), "a", "r"); err != nil {
		t.Fatal(err)
	}
	calls := 0
	store.replace = func(path string, contents []byte, mode fs.FileMode) error {
		calls++
		if err := atomicfile.Replace(path, contents, mode); err != nil {
			t.Fatal(err)
		}
		return &atomicfile.Error{Op: "sync directory", Path: filepath.Dir(path), Replaced: true, Err: errors.New("sync")}
	}
	next := edit(t, func(d map[string]any) { d["summary"] = "Etter usikker skriving" })
	_, err := store.Update("PLAN-1", 1, next, "a", "uncertain")
	e := requireClass(t, err, ClassIOError)
	if !strings.Contains(e.Detail, "PLAN-1 revision 2") || !strings.Contains(e.Detail, "result unknown") || strings.Contains(e.Detail, dir) {
		t.Fatalf("detail %q", e.Detail)
	}
	if calls != 1 {
		t.Fatalf("replace called %d times", calls)
	}
	store.replace = atomicfile.Replace
	history, err := store.History("PLAN-1")
	if err != nil || len(history) != 2 || history[1].Reason != "uncertain" {
		t.Fatalf("history %+v %v", history, err)
	}
	latest, err := store.Show("PLAN-1", 0)
	if err != nil || latest.Revision != 2 || latest.Plan.Summary != "Etter usikker skriving" {
		t.Fatalf("latest %+v %v", latest, err)
	}
	// The writer must read before writing again: the old revision now conflicts.
	_, err = store.Update("PLAN-1", 1, next, "a", "retry")
	requireClass(t, err, ClassConflict)
}
