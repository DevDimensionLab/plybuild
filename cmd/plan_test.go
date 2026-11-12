package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/plan"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const planV1Snapshot = `{
  "title": "Planverk i Ply", "goal": "Agent og menneske leser samme plan",
  "scope": "Ply CLI og lagring; ikke server/PWA", "owner": "ply-epic",
  "selection": "selected", "status": "active",
  "summary": "Behov samlet; spec skrives", "next_action": "Ferdigstill spec",
  "details": "Detaljtekst", "decisions": ["Kun Ply først"], "learning": "Les spec helt",
  "refs": ["backlog:20261111-c97f7518"],
  "steps": [
    {"id":"spec", "title":"Avklar kontrakt", "outcome":"Implementerbar spec", "owner":"ply-epic", "status":"done", "next_action":"", "result":"Spec peerreviewet"},
    {"id":"implement", "title":"Implementer", "outcome":"Testet kandidat", "owner":"Claude", "status":"waiting", "next_action":"Svar på lagringsvalg",
     "attention": {"reason":"Velg lagringsmodell", "who":"bruker", "needs_human":true}, "depends_on":["spec"]},
    {"id":"review", "title":"Kontroller", "outcome":"Review og CLI-pilot grønne", "owner":"ply-epic", "status":"planned", "next_action":"Review kandidat", "depends_on":["implement"]}
  ]
}`

func runPlan(t *testing.T, cwd, stdin string, args ...string) (string, error) {
	t.Helper()
	t.Chdir(cwd)
	command := newPlanCommand(workspace.SystemDependencies())
	command.SilenceUsage, command.SilenceErrors = true, true // as RootCmd under ExecuteE
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&bytes.Buffer{})
	command.SetIn(strings.NewReader(stdin))
	command.SetArgs(args)
	err := command.Execute()
	return stdout.String(), err
}

func planKeys(t *testing.T, out string) []string {
	t.Helper()
	decoder := json.NewDecoder(strings.NewReader(out))
	if token, err := decoder.Token(); err != nil || token != json.Delim('{') {
		t.Fatalf("json %q: %v", out, err)
	}
	var keys []string
	for decoder.More() {
		key, _ := decoder.Token()
		keys = append(keys, key.(string))
		var skip json.RawMessage
		if err := decoder.Decode(&skip); err != nil {
			t.Fatal(err)
		}
	}
	return keys
}

func requirePlanClass(t *testing.T, err error, class string) {
	t.Helper()
	var e *plan.Error
	if !errors.As(err, &e) || e.Class != class {
		t.Fatalf("error %v, want %s", err, class)
	}
}

func TestPlanAddShowUpdateListHistoryInJSON(t *testing.T) {
	root := newAssignmentWorkspace(t)
	added, err := runPlan(t, root, planV1Snapshot, "add", "--project", "a", "--actor", "ply-epic", "--reason", "Brukeren valgte oppdraget", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	wantKeys := []string{"kind", "schema_version", "project", "id", "revision", "created_at", "updated_at", "actor", "reason", "needs_human", "plan", "progress"}
	if keys := planKeys(t, added); !reflect.DeepEqual(keys, wantKeys) {
		t.Fatalf("readback keys %v", keys)
	}
	var readback struct {
		Kind          string        `json:"kind"`
		SchemaVersion int           `json:"schema_version"`
		Project       string        `json:"project"`
		ID            string        `json:"id"`
		Revision      int           `json:"revision"`
		NeedsHuman    bool          `json:"needs_human"`
		Plan          plan.Plan     `json:"plan"`
		Progress      plan.Progress `json:"progress"`
	}
	if err := json.Unmarshal([]byte(added), &readback); err != nil {
		t.Fatal(err)
	}
	if readback.Kind != "PlyPlanReadback@1" || readback.SchemaVersion != 1 || readback.Project != "a" || readback.ID != "PLAN-1" ||
		readback.Revision != 1 || !readback.NeedsHuman || readback.Plan.Steps[1].Attention == nil ||
		readback.Progress != (plan.Progress{Total: 3, Planned: 1, Waiting: 1, Done: 1}) {
		t.Fatalf("readback %+v", readback)
	}
	if _, err := os.Stat(filepath.Join(root, "a", "plans", "PLAN-1.json")); err != nil {
		t.Fatalf("plan not stored in the Project's plans data path: %v", err)
	}

	// From inside the Project, --project may be left out.
	app := filepath.Join(root, "a", "app")
	shown, err := runPlan(t, app, "", "show", "PLAN-1", "--format", "json")
	if err != nil || shown != added {
		t.Fatalf("show %q %v", shown, err)
	}
	updated, err := runPlan(t, app, strings.Replace(planV1Snapshot, `"status":"waiting"`, `"status":"active"`, 1), "update", "PLAN-1", "--revision", "1", "--actor", "ply-epic", "--reason", "Svar mottatt", "--format", "json")
	if err != nil || !strings.Contains(updated, `"revision": 2`) {
		t.Fatalf("update %q %v", updated, err)
	}
	old, err := runPlan(t, app, "", "show", "PLAN-1", "--revision", "1", "--format", "json")
	if err != nil || old != added {
		t.Fatalf("revision 1 after update %q %v", old, err)
	}
	_, err = runPlan(t, app, planV1Snapshot, "update", "PLAN-1", "--revision", "1", "--actor", "a", "--reason", "stale")
	requirePlanClass(t, err, plan.ClassConflict)

	listed, err := runPlan(t, root, "", "list", "--project", "a", "--format", "json")
	if keys := planKeys(t, listed); err != nil || !reflect.DeepEqual(keys, []string{"kind", "schema_version", "project", "items", "invalid"}) {
		t.Fatalf("list %q %v", listed, err)
	}
	var list struct {
		Kind  string          `json:"kind"`
		Items []plan.ListItem `json:"items"`
	}
	if err := json.Unmarshal([]byte(listed), &list); err != nil || list.Kind != "PlyPlanList@1" || len(list.Items) != 1 || !list.Items[0].NeedsHuman || list.Items[0].Revision != 2 {
		t.Fatalf("list %+v %v", list, err)
	}
	filtered, err := runPlan(t, root, "", "list", "--project", "a", "--status", "done", "--format", "json")
	if err != nil || !strings.Contains(filtered, `"items": []`) {
		t.Fatalf("filtered %q %v", filtered, err)
	}

	history, err := runPlan(t, root, "", "history", "PLAN-1", "--project", "a", "--format", "json")
	if keys := planKeys(t, history); err != nil || !reflect.DeepEqual(keys, []string{"kind", "schema_version", "project", "id", "revisions"}) {
		t.Fatalf("history %q %v", history, err)
	}
	var h struct {
		Kind      string              `json:"kind"`
		Revisions []plan.HistoryEntry `json:"revisions"`
	}
	if err := json.Unmarshal([]byte(history), &h); err != nil || h.Kind != "PlyPlanHistory@1" || len(h.Revisions) != 2 || h.Revisions[1].Reason != "Svar mottatt" {
		t.Fatalf("history %+v %v", h, err)
	}
}

func TestPlanShortTextShowsEveryNeededHumanAction(t *testing.T) {
	root := newAssignmentWorkspace(t)
	// The plan is active and has no attention of its own; only a step asks.
	if _, err := runPlan(t, root, planV1Snapshot, "add", "--project", "a", "--actor", "ply-epic", "--reason", "r"); err != nil {
		t.Fatal(err)
	}
	short, err := runPlan(t, root, "", "show", "PLAN-1", "--project", "a")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"PLAN-1 Planverk i Ply", "status: active (selected)", "goal: Agent og menneske leser samme plan",
		"summary: Behov samlet; spec skrives", "owner: ply-epic", "revision 1", "next: Ferdigstill spec",
		"needs human: yes",
		"1. [done] spec: Avklar kontrakt (owner ply-epic)",
		"2. [waiting] implement: Implementer (owner Claude)", "next: Svar på lagringsvalg",
		"attention: Velg lagringsmodell (who: bruker, needs human: yes)",
		"3. [planned] review: Kontroller (owner ply-epic)", "next: Review kandidat",
	} {
		if !strings.Contains(short, want) {
			t.Fatalf("short text lacks %q:\n%s", want, short)
		}
	}
	for _, hidden := range []string{"Detaljtekst", "scope:", "Spec peerreviewet", "Kun Ply først", "backlog:"} {
		if strings.Contains(short, hidden) {
			t.Fatalf("short text shows detail %q:\n%s", hidden, short)
		}
	}
	detailed, err := runPlan(t, root, "", "show", "PLAN-1", "--project", "a", "--details")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"scope: Ply CLI og lagring", "details: Detaljtekst", "result: Spec peerreviewet", "decisions:", "Kun Ply først", "learning: Les spec helt", "refs: backlog:20261111-c97f7518", "depends on: spec"} {
		if !strings.Contains(detailed, want) {
			t.Fatalf("details lack %q:\n%s", want, detailed)
		}
	}
	jsonShort, _ := runPlan(t, root, "", "show", "PLAN-1", "--project", "a", "--format", "json")
	jsonDetails, _ := runPlan(t, root, "", "show", "PLAN-1", "--project", "a", "--format", "json", "--details")
	if jsonShort != jsonDetails || !strings.Contains(jsonShort, "Detaljtekst") {
		t.Fatal("JSON depends on --details")
	}

	listed, err := runPlan(t, root, "", "list", "--project", "a")
	if err != nil || !strings.Contains(listed, "PLAN-1") || !strings.Contains(listed, "needs human") {
		t.Fatalf("list text %q %v", listed, err)
	}
	if err := os.WriteFile(filepath.Join(root, "a", "plans", "PLAN-2.json"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	listed, err = runPlan(t, root, "", "list", "--project", "a")
	if err != nil || !strings.Contains(listed, "PLAN-1") || !strings.Contains(listed, "warning: invalid plan file PLAN-2.json: invalid JSON") {
		t.Fatalf("list with invalid %q %v", listed, err)
	}
	history, err := runPlan(t, root, "", "history", "PLAN-1", "--project", "a")
	if err != nil || !strings.Contains(history, "1 ") || !strings.Contains(history, "ply-epic: r") {
		t.Fatalf("history text %q %v", history, err)
	}
}

func TestPlanReadsCreateNothing(t *testing.T) {
	root := newAssignmentWorkspace(t)
	for _, args := range [][]string{
		{"list", "--project", "a"}, {"list", "--project", "a", "--format", "json"},
		{"show", "PLAN-1", "--project", "a"}, {"history", "PLAN-1", "--project", "a"},
	} {
		_, _ = runPlan(t, root, "", args...)
	}
	if _, err := os.Lstat(filepath.Join(root, "a", "plans")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatal("read created the plans directory")
	}
	out, err := runPlan(t, root, "", "list", "--project", "a", "--format", "json")
	if err != nil || !strings.Contains(out, `"items": []`) || !strings.Contains(out, `"invalid": []`) {
		t.Fatalf("empty list %q %v", out, err)
	}
}

func TestPlanRejectsInvalidArguments(t *testing.T) {
	root := newAssignmentWorkspace(t)
	for name, args := range map[string][]string{
		"add without actor":       {"add", "--project", "a", "--reason", "r"},
		"add without reason":      {"add", "--project", "a", "--actor", "a"},
		"update without revision": {"update", "PLAN-1", "--project", "a", "--actor", "a", "--reason", "r"},
		"update zero revision":    {"update", "PLAN-1", "--project", "a", "--revision", "0", "--actor", "a", "--reason", "r"},
		"show bad revision":       {"show", "PLAN-1", "--project", "a", "--revision", "0"},
		"bad format":              {"list", "--project", "a", "--format", "xml"},
		"bad status":              {"list", "--project", "a", "--status", "paused"},
		"unknown flag":            {"list", "--project", "a", "--secret-input"},
		"extra argument":          {"show", "PLAN-1", "PLAN-2", "--project", "a"},
		"bad id":                  {"show", "PLAN-01", "--project", "a"},
		"unknown subcommand":      {"remove"},
	} {
		_, err := runPlan(t, root, planV1Snapshot, args...)
		requirePlanClass(t, err, plan.ClassInvalidArguments)
		if strings.Contains(err.Error(), "secret-input") || strings.Contains(err.Error(), root) {
			t.Fatalf("%s: error echoes input: %v", name, err)
		}
	}
	_, err := runPlan(t, root, planV1Snapshot+strings.Repeat(" ", 128*1024), "add", "--project", "a", "--actor", "a", "--reason", "r")
	requirePlanClass(t, err, plan.ClassInvalidArguments)
	_, err = runPlan(t, root, planV1Snapshot, "add", "--project", "missing", "--actor", "a", "--reason", "r")
	var projectError *workspace.ProjectError
	if !errors.As(err, &projectError) || projectError.Class != workspace.ErrorProjectNotFound {
		t.Fatalf("missing project %v", err)
	}
}

func TestPlanIsRegisteredBesideExistingCommands(t *testing.T) {
	find := func(path ...string) bool {
		found, rest, err := RootCmd.Find(path)
		return err == nil && len(rest) == 0 && found.Name() == path[len(path)-1]
	}
	for _, path := range [][]string{
		{"plan", "add"}, {"plan", "update"}, {"plan", "show"}, {"plan", "list"}, {"plan", "history"},
		{"assignment", "show"}, {"workflow", "epic", "plan"},
	} {
		if !find(path...) {
			t.Fatalf("command %v not registered", path)
		}
	}
}
