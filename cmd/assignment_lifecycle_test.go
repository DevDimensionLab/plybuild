package cmd

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/assignment"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

var takenLine = regexp.MustCompile(`^ASSIGNMENT-2 taken claim (c_[0-9a-f]{32}) revision 1\n$`)

// lifecycleWorkspace has ASSIGNMENT-1..3 in Project a and returns the a wrapper.
func lifecycleWorkspace(t *testing.T) (string, string) {
	t.Helper()
	root := newAssignmentWorkspace(t)
	a := filepath.Join(root, "a")
	for i := 1; i <= 3; i++ {
		if _, err := runAssignment(t, a, "x", "add", "--title", fmt.Sprintf("Title %d", i), "--repo", "a-app"); err != nil {
			t.Fatal(err)
		}
	}
	return root, a
}

func mustRun(t *testing.T, cwd string, args ...string) string {
	t.Helper()
	out, err := runAssignment(t, cwd, "", args...)
	if err != nil {
		t.Fatalf("%v: %v", args, err)
	}
	return out
}

func listIDs(t *testing.T, cwd string, args ...string) []string {
	t.Helper()
	ids := []string{}
	for _, line := range strings.Split(strings.TrimSuffix(mustRun(t, cwd, append([]string{"list"}, args...)...), "\n"), "\n") {
		if line != "" {
			ids = append(ids, strings.Fields(line)[0])
		}
	}
	return ids
}

func treeState(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		out[rel] = fmt.Sprintf("%x", sha256.Sum256(contents))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func changedFiles(before, after map[string]string) []string {
	changed := []string{}
	for path, sum := range after {
		if before[path] != sum && filepath.Base(path) != ".lock" {
			changed = append(changed, path)
		}
	}
	for path := range before {
		if _, ok := after[path]; !ok {
			changed = append(changed, path)
		}
	}
	sort.Strings(changed)
	return changed
}

func TestAssignmentLifecycleJourneyInText(t *testing.T) {
	root, a := lifecycleWorkspace(t)
	out := mustRun(t, a, "list")
	if lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n"); len(lines) != 3 || !strings.HasSuffix(lines[1], " a-app Title 2") {
		t.Fatalf("v1 list %q", out)
	}
	if out := mustRun(t, a, "show", "ASSIGNMENT-2"); !regexp.MustCompile(`\ncreated: \S+\nstate: available \(revision 0\)\n\nx$`).MatchString(out) {
		t.Fatalf("show %q", out)
	}
	original := treeState(t, filepath.Join(a, "assignments"))["ASSIGNMENT-2.json"]

	out = mustRun(t, a, "take", "ASSIGNMENT-2", "--revision", "0", "--as", "agent a")
	match := takenLine.FindStringSubmatch(out)
	if match == nil {
		t.Fatalf("take %q", out)
	}
	claim := match[1]
	if ids := listIDs(t, a); fmt.Sprint(ids) != "[ASSIGNMENT-1 ASSIGNMENT-3]" {
		t.Fatalf("available %v", ids)
	}
	takenView := mustRun(t, a, "list", "--view", "taken")
	if !regexp.MustCompile(`^ASSIGNMENT-2 \S+ a-app taken \S+Z Title 2\n$`).MatchString(takenView) || strings.Contains(takenView, "agent a") {
		t.Fatalf("taken view %q", takenView)
	}
	if out := mustRun(t, a, "show", "ASSIGNMENT-2"); !regexp.MustCompile(`\nstate: taken by agent a since \S+Z \(revision 1, claim ` + claim + `\)\n\n`).MatchString(out) {
		t.Fatalf("show taken %q", out)
	}

	before := treeState(t, root)
	if out := mustRun(t, a, "done", "ASSIGNMENT-2", "--claim", claim); out != "ASSIGNMENT-2 done revision 2\n" {
		t.Fatalf("done %q", out)
	}
	if changed := changedFiles(before, treeState(t, root)); fmt.Sprint(changed) != "["+filepath.Join("a", "assignments", ".lifecycle", "ASSIGNMENT-2.json")+"]" {
		t.Fatalf("done changed %v", changed)
	}
	if ids := listIDs(t, a); fmt.Sprint(ids) != "[ASSIGNMENT-1 ASSIGNMENT-3]" {
		t.Fatalf("available after done %v", ids)
	}
	if ids := listIDs(t, a, "--view", "taken"); len(ids) != 0 {
		t.Fatalf("taken after done %v", ids)
	}
	archive := mustRun(t, a, "list", "--view", "archive")
	if !regexp.MustCompile(`^ASSIGNMENT-2 \S+ a-app done \S+Z Title 2\n$`).MatchString(archive) {
		t.Fatalf("archive %q", archive)
	}
	all := mustRun(t, a, "list", "--view", "all")
	if !regexp.MustCompile(`^ASSIGNMENT-1 \S+ a-app available - Title 1\nASSIGNMENT-2 \S+ a-app done \S+Z Title 2\nASSIGNMENT-3 \S+ a-app available - Title 3\n$`).MatchString(all) {
		t.Fatalf("all %q", all)
	}
	if out := mustRun(t, a, "show", "ASSIGNMENT-2"); !regexp.MustCompile(`\nstate: done by agent a at \S+Z \(revision 2\)\n\n`).MatchString(out) {
		t.Fatalf("show done %q", out)
	}
	if treeState(t, filepath.Join(a, "assignments"))["ASSIGNMENT-2.json"] != original {
		t.Fatal("original changed")
	}
	_, err := runAssignment(t, a, "", "take", "ASSIGNMENT-2", "--revision", "2", "--as", "b")
	requireAssignmentError(t, err, "assignment_conflict: ASSIGNMENT-2 is done")
}

func TestAssignmentReleaseAndRevokeInText(t *testing.T) {
	_, a := lifecycleWorkspace(t)
	claim := takenLine.FindStringSubmatch(mustRun(t, a, "take", "ASSIGNMENT-2", "--revision", "0", "--as", "a"))[1]
	if out := mustRun(t, a, "release", "ASSIGNMENT-2", "--claim", claim); out != "ASSIGNMENT-2 released revision 2\n" {
		t.Fatalf("release %q", out)
	}
	if out := mustRun(t, a, "release", "ASSIGNMENT-2", "--claim", claim); out != "ASSIGNMENT-2 released revision 2\n" {
		t.Fatalf("repeated release %q", out)
	}
	_, err := runAssignment(t, a, "", "take", "ASSIGNMENT-2", "--revision", "1", "--as", "b")
	requireAssignmentError(t, err, "assignment_conflict: stale revision 1, current 2")
	mustRun(t, a, "take", "ASSIGNMENT-2", "--revision", "2", "--as", "b")
	_, err = runAssignment(t, a, "", "release", "ASSIGNMENT-2", "--revoke", "--revision", "2", "--as", "human", "--reason", "b is gone")
	requireAssignmentError(t, err, "assignment_conflict: stale revision 2, current 3")
	if out := mustRun(t, a, "release", "ASSIGNMENT-2", "--revoke", "--revision", "3", "--as", "human", "--reason", "b is gone"); out != "ASSIGNMENT-2 revoked revision 4\n" {
		t.Fatalf("revoke %q", out)
	}
	if ids := listIDs(t, a); fmt.Sprint(ids) != "[ASSIGNMENT-1 ASSIGNMENT-2 ASSIGNMENT-3]" {
		t.Fatalf("available after revoke %v", ids)
	}
	_, err = runAssignment(t, a, "", "done", "ASSIGNMENT-2", "--claim", claim)
	requireAssignmentError(t, err, "assignment_conflict: claim does not hold ASSIGNMENT-2")
}

func TestAssignmentHasNoArchiveOperation(t *testing.T) {
	root, a := lifecycleWorkspace(t)
	before := treeState(t, root)
	for _, args := range [][]string{{"archive", "ASSIGNMENT-1"}, {"archive"}} {
		_, err := runAssignment(t, a, "", args...)
		requireAssignmentError(t, err, "assignment_invalid_arguments: unknown assignment command")
	}
	for _, view := range []string{"available", "taken", "archive", "all"} {
		mustRun(t, a, "list", "--view", view)
	}
	if changed := changedFiles(before, treeState(t, root)); len(changed) != 0 {
		t.Fatalf("changed %v", changed)
	}
}

func TestAssignmentLifecycleJSONShape(t *testing.T) {
	root := newAssignmentWorkspace(t)
	a := filepath.Join(root, "a")
	dir := filepath.Join(a, "assignments")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	original := `{"id":"ASSIGNMENT-7","title":"T <&>","text":"x\n","created_at":"2026-11-10T18:00:00Z","repos":["a-app"],"topics":["web"]}`
	if err := os.WriteFile(filepath.Join(dir, "ASSIGNMENT-7.json"), []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	v1Assignment := `  "assignment": {
    "id": "ASSIGNMENT-7",
    "title": "T <&>",
    "text": "x\n",
    "created_at": "2026-11-10T18:00:00Z",
    "repos": [
      "a-app"
    ],
    "topics": [
      "web"
    ]
  },
`
	out := mustRun(t, a, "show", "ASSIGNMENT-7", "--format", "json")
	want := "{\n  \"project\": \"a\",\n" + v1Assignment + `  "lifecycle": {
    "state": "available",
    "revision": 0,
    "holder": null,
    "done": null,
    "events": []
  }
}
`
	if out != want {
		t.Fatalf("show json\n%s\nwant\n%s", out, want)
	}

	out = mustRun(t, a, "take", "ASSIGNMENT-7", "--revision", "0", "--as", "claude:9f2c", "--format", "json")
	var taken struct {
		Project   string               `json:"project"`
		ID        string               `json:"id"`
		Lifecycle assignment.Lifecycle `json:"lifecycle"`
	}
	keys := decodeAssignmentJSON(t, out, &taken)
	if len(keys) != 3 || taken.Project != "a" || taken.ID != "ASSIGNMENT-7" || taken.Lifecycle.State != assignment.StateTaken || taken.Lifecycle.Holder == nil {
		t.Fatalf("take json %s", out)
	}
	claim := taken.Lifecycle.Holder.Claim
	var lifecycleKeys map[string]json.RawMessage
	if err := json.Unmarshal(keys["lifecycle"], &lifecycleKeys); err != nil || len(lifecycleKeys) != 5 {
		t.Fatalf("lifecycle keys %v %v", lifecycleKeys, err)
	}
	if !strings.Contains(string(lifecycleKeys["holder"]), `"taken_at"`) || string(lifecycleKeys["done"]) != "null" {
		t.Fatalf("holder %s done %s", lifecycleKeys["holder"], lifecycleKeys["done"])
	}
	out = mustRun(t, a, "show", "ASSIGNMENT-7", "--format", "json")
	if !strings.Contains(out, v1Assignment) || !strings.Contains(out, `"claim": "`+claim+`"`) {
		t.Fatalf("show taken json %s", out)
	}
	for _, view := range []string{"available", "taken", "archive", "all"} {
		out := mustRun(t, a, "list", "--view", view, "--format", "json")
		var listed map[string]json.RawMessage
		if err := json.Unmarshal([]byte(out), &listed); err != nil || len(listed) != 4 || string(listed["view"]) != `"`+view+`"` || string(listed["project"]) != `"a"` {
			t.Fatalf("list %s %s", view, out)
		}
		if strings.Contains(out, claim) || strings.Contains(out, `"claim"`) {
			t.Fatalf("list %s carries a claim: %s", view, out)
		}
	}
	out = mustRun(t, a, "list", "--view", "taken", "--format", "json")
	if !strings.Contains(out, `"state": "taken"`) || !strings.Contains(out, `"actor": "claude:9f2c"`) || !strings.Contains(out, `"taken_at"`) || !strings.Contains(out, `"revision": 1`) {
		t.Fatalf("taken list %s", out)
	}
	out = mustRun(t, a, "done", "ASSIGNMENT-7", "--claim", claim, "--format", "json")
	if !strings.Contains(out, `"state": "done"`) || !strings.Contains(out, `"holder": null`) {
		t.Fatalf("done json %s", out)
	}
	out = mustRun(t, a, "list", "--view", "archive", "--format", "json")
	if !strings.Contains(out, `"done_at"`) || !strings.Contains(out, `"id": "ASSIGNMENT-7"`) {
		t.Fatalf("archive json %s", out)
	}
}

func TestAssignmentLifecycleArgumentErrors(t *testing.T) {
	_, a := lifecycleWorkspace(t)
	const claim = "c_0123456789abcdef0123456789abcdef"
	for name, tc := range map[string]struct {
		args []string
		want string
	}{
		"take without revision":   {[]string{"take", "ASSIGNMENT-1", "--as", "a"}, "assignment_invalid_arguments: --revision is required"},
		"take without as":         {[]string{"take", "ASSIGNMENT-1", "--revision", "0"}, "assignment_invalid_arguments: --as is required"},
		"take bad revision":       {[]string{"take", "ASSIGNMENT-1", "--revision", "x", "--as", "a"}, "assignment_invalid_arguments: invalid flag"},
		"take negative revision":  {[]string{"take", "ASSIGNMENT-1", "--revision", "-1", "--as", "a"}, "assignment_invalid_arguments: revision must be a non-negative integer"},
		"take two-line actor":     {[]string{"take", "ASSIGNMENT-1", "--revision", "0", "--as", "a\nhemmelig"}, "assignment_invalid_arguments: actor must be one line of 1 to 100 characters without control characters"},
		"take without id":         {[]string{"take", "--revision", "0", "--as", "a"}, "assignment_invalid_arguments: expected 1 positional arguments, got 0"},
		"take bad id":             {[]string{"take", "ASSIGNMENT-01", "--revision", "0", "--as", "a"}, "assignment_invalid_arguments: invalid assignment id"},
		"done without claim":      {[]string{"done", "ASSIGNMENT-1"}, "assignment_invalid_arguments: --claim is required"},
		"done malformed claim":    {[]string{"done", "ASSIGNMENT-1", "--claim", "c_hemmelig"}, "assignment_invalid_arguments: claim must be c_ followed by 32 lowercase hex digits"},
		"release neither":         {[]string{"release", "ASSIGNMENT-1"}, "assignment_invalid_arguments: exactly one of --claim and --revoke is required"},
		"release both":            {[]string{"release", "ASSIGNMENT-1", "--claim", claim, "--revoke"}, "assignment_invalid_arguments: exactly one of --claim and --revoke is required"},
		"release claim with as":   {[]string{"release", "ASSIGNMENT-1", "--claim", claim, "--as", "a"}, "assignment_invalid_arguments: --revision, --as and --reason require --revoke"},
		"revoke without revision": {[]string{"release", "ASSIGNMENT-1", "--revoke", "--as", "h", "--reason", "r"}, "assignment_invalid_arguments: --revision is required"},
		"revoke without as":       {[]string{"release", "ASSIGNMENT-1", "--revoke", "--revision", "1", "--reason", "r"}, "assignment_invalid_arguments: --as is required"},
		"revoke without reason":   {[]string{"release", "ASSIGNMENT-1", "--revoke", "--revision", "1", "--as", "h"}, "assignment_invalid_arguments: --reason is required"},
		"revoke long reason":      {[]string{"release", "ASSIGNMENT-1", "--revoke", "--revision", "1", "--as", "h", "--reason", strings.Repeat("r", 201)}, "assignment_invalid_arguments: reason must be one line of 1 to 200 characters without control characters"},
		"bad view":                {[]string{"list", "--view", "done"}, "assignment_invalid_arguments: --view must be available, taken, archive or all"},
		"show list flag":          {[]string{"show", "ASSIGNMENT-1", "--view", "all"}, "assignment_invalid_arguments: invalid flag"},
		"done missing":            {[]string{"done", "ASSIGNMENT-9", "--claim", claim}, "assignment_not_found: ASSIGNMENT-9"},
		"release available":       {[]string{"release", "ASSIGNMENT-1", "--claim", claim}, "assignment_conflict: claim does not hold ASSIGNMENT-1"},
		"revoke available":        {[]string{"release", "ASSIGNMENT-1", "--revoke", "--revision", "0", "--as", "hemmelig", "--reason", "hemmelig"}, "assignment_conflict: ASSIGNMENT-1 is not taken"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := runAssignment(t, a, "", tc.args...)
			requireAssignmentError(t, err, tc.want)
			if strings.Contains(err.Error(), "hemmelig") || strings.Contains(err.Error(), a) {
				t.Fatalf("error leaks input: %v", err)
			}
		})
	}
	if _, err := os.Lstat(filepath.Join(a, "assignments", ".lifecycle")); err == nil {
		t.Fatal("argument errors created .lifecycle")
	}
}

func TestAssignmentInvalidLifecycleIsReportedAndBlocks(t *testing.T) {
	_, a := lifecycleWorkspace(t)
	dir := filepath.Join(a, "assignments", ".lifecycle")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ASSIGNMENT-2.json"), []byte("{broken"), 0o644); err != nil {
		t.Fatal(err)
	}
	out := mustRun(t, a, "list", "--view", "all")
	if out != strings.Join([]string{
		strings.Split(out, "\n")[0],
		strings.Split(out, "\n")[1],
		"invalid: .lifecycle/ASSIGNMENT-2.json: assignment_lifecycle_invalid: invalid JSON", "",
	}, "\n") || strings.Contains(out, "ASSIGNMENT-2 ") {
		t.Fatalf("list %q", out)
	}
	for _, args := range [][]string{
		{"show", "ASSIGNMENT-2"},
		{"take", "ASSIGNMENT-2", "--revision", "0", "--as", "a"},
		{"done", "ASSIGNMENT-2", "--claim", "c_0123456789abcdef0123456789abcdef"},
	} {
		_, err := runAssignment(t, a, "", args...)
		requireAssignmentError(t, err, "assignment_lifecycle_invalid: invalid JSON")
	}
}

// TestHelperAssignmentCLI is not a test: it runs one assignment command as a
// separate process for TestAssignmentConcurrentTakeProcesses.
func TestHelperAssignmentCLI(t *testing.T) {
	args := os.Getenv("PLY_ASSIGNMENT_CLI_ARGS")
	if args == "" {
		t.Skip("helper process only")
	}
	t.Chdir(os.Getenv("PLY_ASSIGNMENT_CLI_DIR"))
	command := newAssignmentCommand(workspace.SystemDependencies())
	command.SilenceUsage, command.SilenceErrors = true, true
	var argv []string
	if err := json.Unmarshal([]byte(args), &argv); err != nil {
		t.Fatal(err)
	}
	command.SetArgs(argv)
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func TestAssignmentConcurrentTakeProcesses(t *testing.T) {
	_, a := lifecycleWorkspace(t)
	const processes = 8
	commands := make([]*exec.Cmd, processes)
	stderr := make([]*strings.Builder, processes)
	for i := range commands {
		args, _ := json.Marshal([]string{"take", "ASSIGNMENT-2", "--revision", "0", "--as", fmt.Sprintf("agent-%d", i)})
		commands[i] = exec.Command(os.Args[0], "-test.run=^TestHelperAssignmentCLI$")
		commands[i].Env = append(os.Environ(), "PLY_ASSIGNMENT_CLI_ARGS="+string(args), "PLY_ASSIGNMENT_CLI_DIR="+a)
		stderr[i] = &strings.Builder{}
		commands[i].Stderr = stderr[i]
	}
	for _, command := range commands {
		if err := command.Start(); err != nil {
			t.Fatal(err)
		}
	}
	won := 0
	for i, command := range commands {
		if err := command.Wait(); err == nil {
			won++
			continue
		}
		if !strings.HasPrefix(stderr[i].String(), "assignment_conflict: ") {
			t.Fatalf("process %d: %q", i, stderr[i].String())
		}
	}
	if won != 1 {
		t.Fatalf("%d processes took ASSIGNMENT-2", won)
	}
	var shown struct {
		Lifecycle assignment.Lifecycle `json:"lifecycle"`
	}
	decodeAssignmentJSON(t, mustRun(t, a, "show", "ASSIGNMENT-2", "--format", "json"), &shown)
	if shown.Lifecycle.Revision != 1 || len(shown.Lifecycle.Events) != 1 {
		t.Fatalf("lifecycle %+v", shown.Lifecycle)
	}
}

func TestAssignmentSkillAndReadmeDescribeLifecycle(t *testing.T) {
	skill, err := os.ReadFile(filepath.Join("..", "skills", "ply-assignment", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(skill)
	order := []string{
		"ply assignment show ASSIGNMENT-<n>",
		"ply assignment take ASSIGNMENT-<n> --revision <r> --as <label>",
		"ply assignment done ASSIGNMENT-<n>",
		"ply assignment release ASSIGNMENT-<n>",
	}
	last := -1
	for _, step := range order {
		at := strings.Index(text, step)
		if at <= last {
			t.Fatalf("skill step %q missing or out of order", step)
		}
		last = at
	}
	for _, rule := range []string{"Never revoke", "never start work without a claim"} {
		if !strings.Contains(text, rule) {
			t.Fatalf("skill lacks %q", rule)
		}
	}
	readme, err := os.ReadFile(filepath.Join("..", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	example := regexp.MustCompile(`--allowedTools [^\n]*`).FindAllString(string(readme), -1)
	found := false
	for _, line := range example {
		if strings.Contains(line, `"Bash(ply assignment take:*)"`) && strings.Contains(line, `"Bash(ply assignment done:*)"`) && strings.Contains(line, `"Bash(ply assignment release:*)"`) {
			found = true
		}
	}
	if !found {
		t.Fatalf("README --allowedTools lacks take, done and release: %v", example)
	}
}

// TestAssignmentSkillStepJumpsFollowTheFlow checks that every "go to step N"
// in the skill names an existing step, and that an adopted claim still reads
// the Task references before working.
func TestAssignmentSkillStepJumpsFollowTheFlow(t *testing.T) {
	skill, err := os.ReadFile(filepath.Join("..", "skills", "ply-assignment", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	steps := map[string]string{}
	for _, match := range regexp.MustCompile(`(?m)^## (\d+)\. (.+)$`).FindAllStringSubmatch(string(skill), -1) {
		steps[match[1]] = match[2]
	}
	// A bullet and its indented continuation lines are one line here.
	var lines []string
	for _, line := range strings.Split(string(skill), "\n") {
		if strings.HasPrefix(line, "  ") && len(lines) > 0 && strings.HasPrefix(lines[len(lines)-1], "- ") {
			lines[len(lines)-1] += " " + strings.TrimSpace(line)
			continue
		}
		lines = append(lines, line)
	}
	jumps := map[string]string{}
	for _, line := range lines {
		if match := regexp.MustCompile(`go to step (\d+)`).FindStringSubmatch(line); match != nil {
			if _, ok := steps[match[1]]; !ok {
				t.Fatalf("jump to missing step %s: %q", match[1], line)
			}
			jumps[line] = steps[match[1]]
		}
	}
	targets := map[string]string{"own label": "", "`available`": ""}
	for line, step := range jumps {
		for key := range targets {
			if strings.Contains(line, key) {
				targets[key] = step
			}
		}
	}
	if !strings.HasPrefix(targets["own label"], "Follow references") {
		t.Fatalf("adopted claim goes to %q, want the references step", targets["own label"])
	}
	if !strings.HasPrefix(targets["`available`"], "Take it") {
		t.Fatalf("available goes to %q, want the take step", targets["`available`"])
	}
}
