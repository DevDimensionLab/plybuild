package taskrun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// These scenarios exercise the installed command boundary against genuine
// native records. The provider and human attestations are explicitly synthetic
// test data, never acceptance of the product candidate running this test.
func deliveryCLIBinary(t *testing.T) string {
	t.Helper()
	if binary := os.Getenv("PLY_DELIVERY_TEST_BINARY"); binary != "" {
		if !filepath.IsAbs(binary) {
			t.Fatal("PLY_DELIVERY_TEST_BINARY must be absolute")
		}
		return binary
	}
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(t.TempDir(), "ply")
	c := exec.Command("go", "build", "-o", binary, "./cmd/ply")
	c.Dir = repo
	if output, err := c.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, output)
	}
	return binary
}

func deliveryCLIRun(binary, cwd string, args ...string) (map[string]json.RawMessage, []byte, error) {
	c := exec.Command(binary, append([]string{"workflow", "delivery"}, append(args, "--format", "json")...)...)
	c.Dir = cwd
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run()
	var result map[string]json.RawMessage
	if decodeErr := json.Unmarshal(stdout.Bytes(), &result); decodeErr != nil {
		return nil, append(stdout.Bytes(), stderr.Bytes()...), fmt.Errorf("CLI output: %w (exit: %v)", decodeErr, err)
	}
	return result, append(stdout.Bytes(), stderr.Bytes()...), err
}

func deliveryCLIString(t *testing.T, r map[string]json.RawMessage, key string) string {
	t.Helper()
	var value string
	if err := json.Unmarshal(r[key], &value); err != nil {
		t.Fatalf("field %s: %v", key, err)
	}
	return value
}

func deliveryCLIOK(t *testing.T, binary, cwd string, args ...string) map[string]json.RawMessage {
	t.Helper()
	r, output, err := deliveryCLIRun(binary, cwd, args...)
	if err != nil {
		t.Fatalf("delivery %v: %v\n%s", args, err, output)
	}
	return r
}

func deliveryCLIRegisterInput(t *testing.T, f deliveryFixture, o WorkflowRun, key string) string {
	t.Helper()
	candidate := o.Delivery.Candidates[len(o.Delivery.Candidates)-1]
	return writeAny(t, f.R.WorkspaceRoot, "registration-"+key+".json", map[string]any{
		"kind": "ply.delivery.registration", "schema_version": 1,
		"publication_key": "cli/" + key, "task_id": "task",
		"task_result_id": candidate.TaskResult.ID, "workflow_run_id": o.RunID,
	})
}

func deliveryCLIQualified(t *testing.T, f deliveryFixture) WorkflowRun {
	t.Helper()
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "observable delivery fixture")
	o, err := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "cli"))
	if err != nil {
		t.Fatal(err)
	}
	return o
}

func deliveryCLISnapshot(t *testing.T, root string) string {
	t.Helper()
	var entries []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		entries = append(entries, path+":"+hash(data))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(entries)
	return strings.Join(entries, "\n")
}

func deliveryCLILatestBase(base workspace.EpicBaseResult) workspace.EpicBaseVersion {
	var latest workspace.EpicBaseVersion
	for _, version := range base.Versions {
		if version.Revision > latest.Revision {
			latest = version
		}
	}
	return latest
}

func TestDeliveryCLIRegistrationRejectsConflictAndPreservesCandidateIsolation(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{Mode: workspace.DeliveryLocalBranch}, ParentRef: "main"})
	o := deliveryCLIQualified(t, f)
	cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	parentBefore := gitOutput(t, f.Parent, "rev-parse", "HEAD")
	registration := deliveryCLIRegisterInput(t, f, o, "original")
	registered := deliveryCLIOK(t, binary, cwd, "register", "--file", registration)
	id := deliveryCLIString(t, registered, "id")
	reject := func(name string, input map[string]any) {
		t.Helper()
		path := writeAny(t, f.R.WorkspaceRoot, name+".json", input)
		if _, output, err := deliveryCLIRun(binary, cwd, "register", "--file", path); err == nil {
			t.Fatalf("invalid registration %s succeeded: %s", name, output)
		}
	}
	for _, tc := range []struct{ name, field, value string }{
		{"same-key-conflict", "task_id", "other-task"},
		{"wrong-task", "task_id", "other-task"},
		{"wrong-result", "task_result_id", "trs_" + strings.Repeat("0", 32)},
	} {
		in := workflowProviderDocument(t, registration)
		if tc.name != "same-key-conflict" {
			in["publication_key"] = tc.name
		}
		in[tc.field] = tc.value
		reject(tc.name, in)
	}
	candidate := o.Delivery.Candidates[len(o.Delivery.Candidates)-1]
	artifact := candidate.TaskResult.Artifacts[0]
	original, err := os.ReadFile(artifact.Locator)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(artifact.Locator)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Chmod(artifact.Locator, 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(artifact.Locator, []byte("tampered fixture artifact"), 0600); err != nil {
		t.Fatal(err)
	}
	in := workflowProviderDocument(t, registration)
	in["publication_key"] = "tampered-evidence"
	reject("tampered-evidence", in)
	if err = os.Remove(artifact.Locator); err != nil {
		t.Fatal(err)
	}
	in["publication_key"] = "missing-evidence"
	reject("missing-evidence", in)
	if err = os.WriteFile(artifact.Locator, original, info.Mode().Perm()); err != nil {
		t.Fatal(err)
	}
	// Native QA may arrive outside the owner conversation. Its old candidate
	// identity must remain exact when a later candidate is qualified.
	if _, err = workflowhandoff.RecordDeliveryHumanQA(f.D.Workflow, "task", candidate.TaskResult, "pass", deliveryTestHuman(t, f, o, "pass")); err != nil {
		t.Fatal(err)
	}
	if ready := deliveryCLIOK(t, binary, cwd, "check", id); deliveryCLIString(t, ready, "state") != "ready" {
		t.Fatalf("restored exact evidence and native pass not ready: %+v", ready)
	}
	deliveryTestCandidate(t, f, "replacement candidate with distinct bytes")
	if prior, _, err := deliveryCLIRun(binary, cwd, "execute", id); err == nil || deliveryCLIString(t, prior, "state") != "blocked" {
		t.Fatal("old Delivery accepted a moved candidate")
	}
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "replacement"))
	if err != nil {
		t.Fatal(err)
	}
	replacementInput := workflowProviderDocument(t, deliveryCLIRegisterInput(t, f, o, "replacement"))
	replacementInput["predecessor_id"] = id
	replacementPath := writeAny(t, f.R.WorkspaceRoot, "replacement-with-predecessor.json", replacementInput)
	replacement := deliveryCLIOK(t, binary, cwd, "register", "--file", replacementPath)
	replacementID := deliveryCLIString(t, replacement, "id")
	if replacementID == id {
		t.Fatal("replacement overwrote the original Delivery")
	}
	if checked := deliveryCLIOK(t, binary, cwd, "check", replacementID); deliveryCLIString(t, checked, "state") != "awaiting_human" {
		t.Fatalf("new candidate inherited the prior candidate's pass: %+v", checked)
	}
	if waiting, _, err := deliveryCLIRun(binary, cwd, "execute", replacementID); err == nil || deliveryCLIString(t, waiting, "state") != "awaiting_human" {
		t.Fatal("new candidate delivered using the old pass")
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != parentBefore {
		t.Fatal("registration or rejected execution changed the local target")
	}
	listed := deliveryCLIOK(t, binary, cwd, "list", "--task", "task")
	var history []json.RawMessage
	if err = json.Unmarshal(listed["deliveries"], &history); err != nil || len(history) != 2 {
		t.Fatalf("registration conflicts or replacement corrupted Delivery history: %v count=%d", err, len(history))
	}
}

func TestDeliveryCLILocalJourneys(t *testing.T) {
	binary := deliveryCLIBinary(t)
	for _, tc := range []struct{ mode, parent string }{
		{workspace.DeliveryLocalEpic, "epic"},
		{workspace.DeliveryLocalBranch, "main"},
		{workspace.DeliveryLocalBranch, "master"},
	} {
		t.Run(tc.mode+"/"+tc.parent, func(t *testing.T) {
			f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{Mode: tc.mode}, ParentRef: tc.parent})
			o := deliveryCLIQualified(t, f)
			cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
			parentBefore := gitOutput(t, f.Parent, "rev-parse", "HEAD")
			registration := deliveryCLIRegisterInput(t, f, o, "first")
			r := deliveryCLIOK(t, binary, cwd, "register", "--file", registration)
			id := deliveryCLIString(t, r, "id")
			if id == "" || deliveryCLIString(t, r, "state") == "delivered" {
				t.Fatal("registration manufactured delivery")
			}
			repeated := deliveryCLIOK(t, binary, cwd, "register", "--file", registration)
			if deliveryCLIString(t, repeated, "id") != id {
				t.Fatal("identical registration changed identity")
			}
			second := deliveryCLIOK(t, binary, cwd, "register", "--file", deliveryCLIRegisterInput(t, f, o, "second-key"))
			secondID := deliveryCLIString(t, second, "id")
			if secondID == id {
				t.Fatal("different publication key did not retain a separate Delivery")
			}
			beforeRead := deliveryCLISnapshot(t, filepath.Join(f.R.WorkspaceRoot, ".ply"))
			for _, args := range [][]string{{"check", id}, {"show", id}, {"list", "--task", "task", "--epic", "epic"}} {
				deliveryCLIOK(t, binary, cwd, args...)
			}
			if afterRead := deliveryCLISnapshot(t, filepath.Join(f.R.WorkspaceRoot, ".ply")); afterRead != beforeRead {
				t.Fatal("read/check mutated native state")
			}
			waiting, _, _ := deliveryCLIRun(binary, cwd, "execute", id)
			if waiting != nil && deliveryCLIString(t, waiting, "state") == "delivered" {
				t.Fatal("missing human answer delivered")
			}
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") != parentBefore {
				t.Fatal("registration/check/missing QA changed parent")
			}
			var err error
			o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
			if err != nil {
				t.Fatal(err)
			}
			ready := deliveryCLIOK(t, binary, cwd, "check", id)
			if deliveryCLIString(t, ready, "state") != "ready" {
				t.Fatalf("matching pass not ready: %s", ready["state"])
			}
			// Separate processes must serialize this effect, even across Delivery IDs.
			var wg sync.WaitGroup
			outcomes := make([]map[string]json.RawMessage, 2)
			failures := make([]error, 2)
			outputs := make([][]byte, 2)
			for i, deliveryID := range []string{id, secondID} {
				wg.Add(1)
				go func(i int, deliveryID string) {
					defer wg.Done()
					outcomes[i], outputs[i], failures[i] = deliveryCLIRun(binary, cwd, "execute", deliveryID)
				}(i, deliveryID)
			}
			wg.Wait()
			for i := range outcomes {
				if failures[i] != nil || deliveryCLIString(t, outcomes[i], "state") != "delivered" {
					t.Fatalf("concurrent execute %d: %v\n%s", i, failures[i], outputs[i])
				}
			}
			candidate := o.Delivery.Candidates[len(o.Delivery.Candidates)-1]
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") != candidate.OID {
				t.Fatal("wrong local candidate integrated")
			}
			status, statusErr := exec.Command(testGit(t), "-C", f.Parent, "status", "--porcelain").Output()
			if statusErr != nil || len(status) != 0 {
				t.Fatalf("local return left target dirty: %v %s", statusErr, status)
			}
			refLog := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/"+tc.parent)
			for _, deliveryID := range []string{id, secondID} {
				done := deliveryCLIOK(t, binary, cwd, "execute", deliveryID)
				if deliveryCLIString(t, done, "state") != "delivered" || string(done["native_closed"]) != "true" {
					t.Fatalf("completion lost: %+v", done)
				}
			}
			if gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/"+tc.parent) != refLog {
				t.Fatal("retry repeated Git integration")
			}
			run, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
			if err != nil || run.Delivery.Phase != "completed" {
				t.Fatalf("native workflow not closed: %v phase=%s", err, run.Delivery.Phase)
			}
			base, err := workspace.ShowEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
			if err != nil {
				t.Fatal(err)
			}
			if latest := deliveryCLILatestBase(base); latest.OID != candidate.OID || latest.Tree != candidate.Tree {
				t.Fatalf("native base did not advance: %+v", latest)
			}
		})
	}
}

const deliveryCLILocalCrashGit = `import json, os, signal, subprocess, sys
from pathlib import Path
p = Path(__file__).resolve().parent.parent / "local-crash-fixture.json"
m = json.loads(p.read_text())
a = sys.argv[1:]
if "merge" in a:
    if "-C" not in a or a[a.index("-C")+1] != m["parent"] or a[-1] != m["oid"]:
        raise SystemExit("unexpected integration target")
    r = subprocess.run([m["git"], *a])
    if r.returncode == 0:
        m["merge_count"] = m.get("merge_count", 0) + 1
        crash = m.pop("crash", False)
        p.write_text(json.dumps(m))
        if crash:
            os.kill(os.getppid(), signal.SIGKILL)
            raise SystemExit(9)
    raise SystemExit(r.returncode)
if any(v in a for v in ["push", "fetch", "pull", "clone", "ls-remote"]):
    raise SystemExit("local delivery attempted a network-capable command")
os.execv(m["git"], [m["git"], *a])
`

func TestDeliveryCLILocalCrashRecovery(t *testing.T) {
	binary := deliveryCLIBinary(t)
	for _, tc := range []struct{ mode, parent string }{
		{workspace.DeliveryLocalEpic, "epic"},
		{workspace.DeliveryLocalBranch, "main"},
	} {
		t.Run(tc.mode, func(t *testing.T) {
			f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{Mode: tc.mode}, ParentRef: tc.parent})
			o := deliveryCLIQualified(t, f)
			var err error
			o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
			if err != nil {
				t.Fatal(err)
			}
			candidate := o.Delivery.Candidates[len(o.Delivery.Candidates)-1]
			cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
			git, err := exec.LookPath("git")
			if err != nil {
				t.Fatal(err)
			}
			python, err := exec.LookPath("python3")
			if err != nil {
				t.Fatal(err)
			}
			model := writeAny(t, f.R.WorkspaceRoot, "local-crash-fixture.json", map[string]any{"git": git, "parent": f.Parent, "oid": candidate.OID, "crash": true})
			bin := filepath.Join(f.R.WorkspaceRoot, "local-adapters")
			if err = os.Mkdir(bin, 0700); err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(filepath.Join(bin, "git"), []byte("#!"+python+"\n"+deliveryCLILocalCrashGit), 0700); err != nil {
				t.Fatal(err)
			}
			t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
			registered := deliveryCLIOK(t, binary, cwd, "register", "--file", deliveryCLIRegisterInput(t, f, o, "local-crash"))
			id := deliveryCLIString(t, registered, "id")
			if _, output, err := deliveryCLIRun(binary, cwd, "execute", id); err == nil {
				t.Fatalf("interrupted invocation reported success: %s", output)
			}
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") != candidate.OID {
				t.Fatal("fixture did not interrupt after the actual Git effect")
			}
			interrupted, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
			if err != nil || interrupted.Delivery.Phase == "completed" {
				t.Fatalf("fixture did not interrupt before native closure: %v", err)
			}
			base, err := workspace.ShowEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
			if err != nil {
				t.Fatal(err)
			}
			if deliveryCLILatestBase(base).OID == candidate.OID {
				t.Fatal("fixture did not interrupt before the registered base update")
			}
			refLog := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/"+tc.parent)
			done := deliveryCLIOK(t, binary, cwd, "execute", id)
			if deliveryCLIString(t, done, "state") != "delivered" || string(done["native_closed"]) != "true" {
				t.Fatalf("interrupted local delivery did not recover: %+v", done)
			}
			deliveryCLIOK(t, binary, cwd, "execute", id)
			if state := workflowProviderDocument(t, model); state["merge_count"] != float64(1) {
				t.Fatalf("local integration was repeated: %+v", state)
			}
			if gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/"+tc.parent) != refLog {
				t.Fatal("recovery or retry changed the target again")
			}
			run, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
			if err != nil || run.Delivery.Phase != "completed" {
				t.Fatalf("recovered workflow not completed: %v", err)
			}
			base, err = workspace.ShowEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
			if err != nil {
				t.Fatal(err)
			}
			if latest := deliveryCLILatestBase(base); latest.OID != candidate.OID || latest.Tree != candidate.Tree {
				t.Fatalf("recovery left the registered base stale: %+v", latest)
			}
		})
	}
}

func TestDeliveryCLIExportHumanJourney(t *testing.T) {
	root := os.Getenv("PLY_DELIVERY_JOURNEY_ROOT")
	if root == "" {
		t.Skip("export only when explicitly requested for an installed disposable journey")
	}
	if !filepath.IsAbs(root) {
		t.Fatal("journey root must be absolute")
	}
	binary := deliveryCLIBinary(t)
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{Mode: workspace.DeliveryLocalBranch}, ParentRef: "main", Root: root})
	o := deliveryCLIQualified(t, f)
	var err error
	o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
	if err != nil {
		t.Fatal(err)
	}
	cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	registration := deliveryCLIRegisterInput(t, f, o, "human-journey")
	r := deliveryCLIOK(t, binary, cwd, "register", "--file", registration)
	manifest := map[string]any{"binary": binary, "workspace": f.R.WorkspaceRoot, "task_worktree": cwd, "target_worktree": f.Parent, "target_ref": "refs/heads/main", "delivery_id": deliveryCLIString(t, r, "id"), "registration_file": registration, "candidate_oid": o.Delivery.Candidates[len(o.Delivery.Candidates)-1].OID, "fixture_notice": "Disposable product exercise; provider and QA records are synthetic test data, not human approval of the implementation."}
	path := writeAny(t, f.R.WorkspaceRoot, "human-journey.json", manifest)
	t.Log("Prepared installed product journey: " + path)
}

// The Git wrapper preserves real local Git observations/effects while routing
// the agreed GitHub remote to an isolated bare repository. Every gh response is
// controlled fixture data; no credential or network operation can escape here.
const deliveryCLIGitStandin = `import json, os, subprocess, sys
from pathlib import Path
state_path = Path(__file__).resolve().parent.parent / "pr-fixture.json"
m = json.loads(state_path.read_text())
args = sys.argv[1:]
config_args = []
while args[:1] == ["-c"]:
    config_args.extend(args[:2])
    args = args[2:]
if args[:2] == ["remote", "get-url"]:
    if args[-1] != "origin": raise SystemExit("unexpected remote")
    print("https://github.com/fixture/product.git")
elif args[:2] == ["ls-remote", "--heads"]:
    if args[2:] != ["origin", m["source_ref"], "refs/heads/main"]: raise SystemExit("unexpected ref query")
    raise SystemExit(subprocess.run([m["git"], "ls-remote", "--heads", m["remote"], *args[3:]]).returncode)
elif args[:1] == ["push"]:
    if args != ["push", "--no-follow-tags", "--recurse-submodules=no", "--no-verify", "--", "origin", m["oid"]+":"+m["source_ref"]]: raise SystemExit("unauthorized fixture push")
    if config_args != ["-c", "remote.origin.mirror=false", "-c", "core.hooksPath=/dev/null"]: raise SystemExit("push configuration not isolated")
    result = subprocess.run([m["git"], *config_args, *args[:-2], m["remote"], args[-1]])
    if result.returncode == 0:
        m["push_count"] = m.get("push_count", 0) + 1
        state_path.write_text(json.dumps(m))
    raise SystemExit(result.returncode)
else:
    if any(a in args for a in ["fetch", "pull", "clone"]): raise SystemExit("unexpected network-capable command")
    os.execv(m["git"], [m["git"], *args])
`

const deliveryCLIGHStandin = `import json, os, signal, sys
from pathlib import Path
p = Path(__file__).resolve().parent.parent / "pr-fixture.json"
m = json.loads(p.read_text())
a = sys.argv[1:]
def flag(name): return a[a.index(name)+1]
if flag("--repo") != "github.com/fixture/product": raise SystemExit("wrong repository or host")
if a[:2] == ["pr", "list"]:
    print(json.dumps([m["pr"]] if m.get("pr") else []))
elif a[:2] == ["pr", "create"]:
    if m.get("pr"): raise SystemExit("duplicate PR creation")
    if flag("--head") != m["source_ref"][len("refs/heads/"):] or flag("--base") != "main": raise SystemExit("wrong PR target")
    m["pr"] = {"number":17,"url":"https://github.com/fixture/product/pull/17",
        "headRefName":flag("--head"),"headRefOid":m["oid"],"baseRefName":"main",
        "headRepository":{"name":"product"},"headRepositoryOwner":{"login":"fixture"},
        "isCrossRepository":False,"body":Path(flag("--body-file")).read_text(),
        "assignees":[],"reviewRequests":[]}
    m["create_count"] = m.get("create_count", 0) + 1
    kill = m.pop("kill_after_create", False)
    p.write_text(json.dumps(m))
    if kill:
        os.kill(os.getppid(), signal.SIGKILL)
        raise SystemExit(9)
    print(m["pr"]["url"])
elif a[:2] == ["pr", "edit"]:
    if a[2] != "17": raise SystemExit("wrong PR")
    for option, field in [("--add-assignee","assignees"),("--add-reviewer","reviewRequests")]:
        if option in a:
            for login in flag(option).split(","):
                if not any(v["login"].lower()==login.lower() for v in m["pr"][field]): m["pr"][field].append({"login":login})
    for option, field in [("--remove-assignee","assignees"),("--remove-reviewer","reviewRequests")]:
        if option in a:
            removing = [v.lower() for v in flag(option).split(",")]
            m["pr"][field] = [v for v in m["pr"][field] if v["login"].lower() not in removing]
    m["edit_count"] = m.get("edit_count", 0) + 1
    p.write_text(json.dumps(m))
else:
    raise SystemExit("unexpected gh effect: "+repr(a))
`

func deliveryCLIPRStandins(t *testing.T, f deliveryFixture, oid string, kill bool) string {
	t.Helper()
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	git, err = filepath.EvalSymlinks(git)
	if err != nil {
		t.Fatal(err)
	}
	python, err := exec.LookPath("python3")
	if err != nil {
		t.Fatal(err)
	}
	remote := filepath.Join(f.R.WorkspaceRoot, "isolated-remote.git")
	runGit(t, f.Parent, "init", "--bare", remote)
	runGit(t, f.Parent, "push", remote, "HEAD:refs/heads/main")
	source := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	runGit(t, source, "remote", "add", "origin", "https://github.com/fixture/product.git")
	model := writeAny(t, f.R.WorkspaceRoot, "pr-fixture.json", map[string]any{"git": git, "remote": remote, "oid": oid, "source_ref": gitOutput(t, source, "symbolic-ref", "HEAD"), "kill_after_create": kill})
	bin := filepath.Join(f.R.WorkspaceRoot, "pr-adapters")
	if err = os.Mkdir(bin, 0700); err != nil {
		t.Fatal(err)
	}
	for name, script := range map[string]string{"git": deliveryCLIGitStandin, "gh": deliveryCLIGHStandin} {
		if err = os.WriteFile(filepath.Join(bin, name), []byte("#!"+python+"\n"+script), 0700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	return model
}

func TestDeliveryCLIPRJourneyAndCrashRecovery(t *testing.T) {
	binary := deliveryCLIBinary(t)
	for _, crash := range []bool{false, true} {
		t.Run(fmt.Sprintf("crash-after-create-%t", crash), func(t *testing.T) {
			f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{Mode: workspace.DeliveryPullRequest, GitHubRepository: "fixture/product", Remote: "origin"}, ParentRef: "main"})
			o := deliveryCLIQualified(t, f)
			var err error
			o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
			if err != nil {
				t.Fatal(err)
			}
			candidate := o.Delivery.Candidates[len(o.Delivery.Candidates)-1]
			cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
			parentBefore := gitOutput(t, f.Parent, "rev-parse", "HEAD")
			model := deliveryCLIPRStandins(t, f, candidate.OID, crash)
			preferences := filepath.Join(f.R.WorkspaceRoot, "private-pr-preferences.yaml")
			if err = os.WriteFile(preferences, []byte("schema_version: 1\ndefaults:\n  assignees: [fixture-assignee]\n  reviewers: [fixture-reviewer]\n"), 0600); err != nil {
				t.Fatal(err)
			}
			r := deliveryCLIOK(t, binary, cwd, "register", "--file", deliveryCLIRegisterInput(t, f, o, "pr"))
			id := deliveryCLIString(t, r, "id")
			ready := deliveryCLIOK(t, binary, cwd, "check", id)
			if deliveryCLIString(t, ready, "state") != "ready" {
				t.Fatalf("PR check not ready: %+v", ready)
			}
			if crash {
				_, _, err = deliveryCLIRun(binary, cwd, "execute", id, "--preferences-file", preferences)
				if err == nil {
					t.Fatal("crashed process reported success")
				}
			}
			done := deliveryCLIOK(t, binary, cwd, "execute", id, "--preferences-file", preferences)
			if deliveryCLIString(t, done, "state") != "delivered" || string(done["native_closed"]) != "true" {
				t.Fatalf("PR not delivered: %+v", done)
			}
			if string(done["local_integration"]) != "null" {
				t.Fatal("PR receipt fabricated local integration")
			}
			var pr struct {
				URL, HeadOID    string
				MetadataApplied bool
			}
			// JSON uses snake case; inspect the exact wire fields.
			var wire map[string]json.RawMessage
			if err = json.Unmarshal(done["pull_request"], &wire); err != nil {
				t.Fatal(err)
			}
			pr.URL = deliveryCLIString(t, wire, "url")
			pr.HeadOID = deliveryCLIString(t, wire, "head_oid")
			if pr.URL != "https://github.com/fixture/product/pull/17" || pr.HeadOID != candidate.OID || string(wire["metadata_applied"]) != "true" {
				t.Fatalf("wrong PR receipt: %s", done["pull_request"])
			}
			deliveryCLIOK(t, binary, cwd, "execute", id, "--preferences-file", preferences)
			state := workflowProviderDocument(t, model)
			if state["create_count"] != float64(1) || state["push_count"] != float64(1) {
				t.Fatalf("duplicate effects: create=%v push=%v", state["create_count"], state["push_count"])
			}
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") != parentBefore || gitOutput(t, state["remote"].(string), "rev-parse", "refs/heads/main") != parentBefore {
				t.Fatal("PR changed local or remote base")
			}
			var events []struct{ Kind, ID, URL string }
			if err = json.Unmarshal(done["events"], &events); err != nil {
				t.Fatal(err)
			}
			created := 0
			for _, event := range events {
				if event.Kind == "delivery.pr_created" {
					created++
					if event.ID == "" || event.URL != pr.URL {
						t.Fatal("creation event missing stable identity/URL")
					}
				}
			}
			if created != 1 {
				t.Fatalf("PR creation events=%d", created)
			}
			run, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
			if err != nil || run.Delivery.Phase != "completed" {
				t.Fatalf("PR workflow closure: %v phase=%s", err, run.Delivery.Phase)
			}
			registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
			if err != nil {
				t.Fatal(err)
			}
			if len(registry.IntegrationResults) != 0 {
				t.Fatal("PR workflow fabricated native local integration")
			}
		})
	}
}
