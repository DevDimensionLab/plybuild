package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Every GitHub call is handled by this local program. The remote merge itself
// updates only the fixture's bare repository, preserving the local return ref.
const integrationCLIAPIStandin = `import json, os, signal, subprocess, sys
from pathlib import Path
p = Path(__file__).resolve().parent.parent / "pr-fixture.json"
m = json.loads(p.read_text())
a = sys.argv[1:]
if a[:1] == ["api"]:
    def flag(name): return a[a.index(name)+1]
    if flag("--hostname") != "github.com": raise SystemExit("wrong host")
    method = flag("--method")
    endpoint = next(v for v in a if v.startswith("repos/"))
    prpath = "repos/fixture/product/pulls/17"
    base = subprocess.check_output([m["git"], "-C", m["remote"], "rev-parse", "refs/heads/main"], text=True).strip()
    def reply(status, value):
        body=json.dumps(value)
        sys.stdout.write("HTTP/1.1 "+str(status)+" Fixture\r\nContent-Type: application/json\r\nContent-Length: "+str(len(body))+"\r\n\r\n"+body)
        raise SystemExit(0)
    if method == "GET" and endpoint == prpath:
        merged = bool(m.get("merge_oid"))
        reply(200, {"number":17,"html_url":m["pr"]["url"],"state":"closed" if merged else "open","draft":False,
            "head":{"ref":m["source_ref"][len("refs/heads/"):],"sha":m["oid"],"repo":{"full_name":"fixture/product"}},
            "base":{"ref":"main","sha":base,"repo":{"full_name":"fixture/product"}},
            "merged":merged,"merged_at":"2026-10-08T09:00:00Z" if merged else None,
            "merge_commit_sha":m.get("merge_oid"),"mergeable":True,"mergeable_state":"clean","stack":None})
    if method == "GET" and endpoint == "repos/fixture/product":
        reply(200,{"full_name":"fixture/product","allow_merge_commit":True,"allow_squash_merge":True,"permissions":{"push":True}})
    details={"uuid":"00000000-0000-4000-8000-000000000017","merge_method":m["method"],"merge_action":"direct_merge","expected_head_sha":m["oid"],"bypass_rules":False}
    if method == "GET" and endpoint == prpath+"/merge-async/"+details["uuid"]:
        reply(200,{"status":"pending","details":details})
    if method == "PUT" and endpoint == prpath+"/merge-async":
        if any(v not in a for v in ["sha="+m["oid"],"merge_method="+m["method"],"merge_action=direct_merge","bypass_rules=false"]):
            raise SystemExit("unconfirmed merge parameters")
        if m.get("merge_count",0): raise SystemExit("duplicate merge request")
        m["merge_count"]=1
        if m.get("pending"):
            p.write_text(json.dumps(m))
            reply(202,{"status":"pending","details":details})
        tree=subprocess.check_output([m["git"],"-C",m["remote"],"rev-parse",m["oid"]+"^{tree}"],text=True).strip()
        parents=["-p",base]
        if m["method"] == "merge": parents += ["-p",m["oid"]]
        oid=subprocess.check_output([m["git"],"-C",m["remote"],"-c","user.name=Fixture","-c","user.email=fixture@example.invalid","commit-tree",tree,*parents,"-m","Isolated "+m["method"]],text=True).strip()
        subprocess.run([m["git"],"-C",m["remote"],"update-ref","refs/heads/main",oid,base],check=True)
        m["merge_oid"]=oid
        crash=m.pop("kill_after_merge",False)
        p.write_text(json.dumps(m))
        if crash:
            os.kill(os.getppid(),signal.SIGKILL)
            raise SystemExit(9)
        details["sha"]=oid
        reply(200,{"status":"merged","details":details})
    raise SystemExit("unexpected API effect: "+repr(a))
`

func integrationCLIPRFixture(t *testing.T, binary, method string, crash, pending bool) (deliveryFixture, WorkflowRun, string, string) {
	t.Helper()
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{SchemaVersion: 2, IntegrationOwner: workspace.IntegrationOwnerHuman, Mode: workspace.DeliveryPullRequest, GitHubRepository: "fixture/product", Remote: "origin"}, ParentRef: "main"})
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	o := deliveryCLIQualified(t, f)
	var e error
	o, e = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
	if e != nil {
		t.Fatal(e)
	}
	candidate := o.Delivery.Candidates[0]
	model := deliveryCLIPRStandins(t, f, candidate.OID, false)
	state := workflowProviderDocument(t, model)
	state["method"], state["kill_after_merge"], state["pending"] = method, crash, pending
	writeAny(t, f.R.WorkspaceRoot, "pr-fixture.json", state)
	python, e := exec.LookPath("python3")
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(f.R.WorkspaceRoot, "pr-adapters", "gh"), []byte("#!"+python+"\n"+integrationCLIAPIStandin+deliveryCLIGHStandin), 0700); e != nil {
		t.Fatal(e)
	}
	source := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	id := deliveryCLIString(t, deliveryCLIOK(t, binary, source, "register", "--file", deliveryCLIRegisterInput(t, f, o, "pr-integration")), "id")
	preferences := filepath.Join(f.R.WorkspaceRoot, "empty-preferences.yaml")
	if e = os.WriteFile(preferences, []byte("schema_version: 1\ndefaults:\n  assignees: []\n  reviewers: []\n"), 0600); e != nil {
		t.Fatal(e)
	}
	done := deliveryCLIOK(t, binary, source, "execute", id, "--preferences-file", preferences)
	if deliveryCLIString(t, done, "state") != "delivered" {
		t.Fatal("PR publication did not complete")
	}
	o, e = ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic fixture releases the completed PR publication owner.")
	if e != nil {
		t.Fatal(e)
	}
	return f, o, id, model
}

func integrationCLIPrepareSuccessor(t *testing.T, f deliveryFixture) {
	t.Helper()
	var entries []workspace.TaskGoalQueueEntry
	for _, id := range []workspace.TaskID{"next-task", "later-task"} {
		if _, e := workspace.CreateTask(f.D.Workspace, workspace.TaskCreateInput{TaskID: id, Title: string(id), Description: "Independent synthetic queued Task", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
			t.Fatal(e)
		}
		r, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
		if e != nil {
			t.Fatal(e)
		}
		var problem workspace.TaskRevisionRef
		for _, p := range r.TaskProblemRevisions {
			if p.TaskID == id {
				problem = workspace.TaskRevisionRef{Revision: p.Revision, ManifestSHA256: p.ManifestSHA256}
			}
		}
		draft := workflowProviderDocument(t, filepath.Join(f.R.WorkspaceRoot, "goal-draft.json"))
		draft["task_id"], draft["publication_key"], draft["problem"] = id, "fixture/"+string(id), problem
		publication, e := workspace.RecordTaskSpec(f.D.Workspace, workspace.TaskContentInput{TaskID: id, File: writeAny(t, f.R.WorkspaceRoot, string(id)+"-goal.json", draft)})
		if e != nil {
			t.Fatal(e)
		}
		goal := workspace.TaskGoalRef{SpecID: "solution", Spec: workspace.TaskRevisionRef{Revision: *publication.OutcomeRef.Revision, ManifestSHA256: publication.OutcomeRef.ManifestSHA256}}
		entries = append(entries, workspace.TaskGoalQueueEntry{TaskID: id, Goal: &goal})
	}
	current, e := workspace.ListTaskQueue(f.D.Workspace, workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, false)
	if e != nil {
		t.Fatal(e)
	}
	draft := workspace.WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: "fixture/next-queue", ProjectID: "ply", RepoID: "ply", EpicID: "epic", ExpectedRevision: current.Revision, Entries: entries, Recorder: workspace.TaskGoalQueueRecorder{ActorClaim: "Fixture planner", ControlSurface: "Synthetic queue isolation test", RecordedAtUTC: "2026-10-08T05:00:00Z"}}
	if _, e = workspace.SetTaskQueue(f.D.Workspace, writeAny(t, f.R.WorkspaceRoot, "next-queue.json", draft)); e != nil {
		t.Fatal(e)
	}
	preview, e := workspace.PreviewTaskGoalExecution(f.D.Workspace, workspace.TaskGoalExecuteInput{Target: workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, Next: true, AcceptancePath: filepath.Join(f.R.WorkspaceRoot, "next-acceptance.sh")})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = workspace.PrepareTaskGoalExecution(f.D.Workspace, preview.Input, preview.Confirmation, workspace.QueueHumanDecision{ActorClaim: "Synthetic fixture operator", DecidedAtUTC: "2026-10-08T05:00:00Z", Source: "explicit_human_instruction", Statement: "Prepare only the unrelated next fixture Task."}); e != nil {
		t.Fatal(e)
	}
}

func TestIntegrationCLIPRMergeSquashAndCrashRecovery(t *testing.T) {
	binary := deliveryCLIBinary(t)
	for _, tc := range []struct {
		method string
		crash  bool
	}{{"merge", false}, {"squash", true}} {
		t.Run(fmt.Sprintf("%s-crash-%t", tc.method, tc.crash), func(t *testing.T) {
			f, o, id, model := integrationCLIPRFixture(t, binary, tc.method, tc.crash, false)
			integrationCLIPrepareSuccessor(t, f)
			parentBefore := gitOutput(t, f.Parent, "rev-parse", "HEAD")
			queueBefore, e := workspace.ListTaskQueue(f.D.Workspace, workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, false)
			if e != nil {
				t.Fatal(e)
			}
			if queueBefore.Current == nil || queueBefore.Current.TaskID != "next-task" || len(queueBefore.Pending) != 1 || queueBefore.Pending[0].TaskID != "later-task" {
				t.Fatal("fixture did not preserve a newer current and pending Task")
			}
			args := []string{"--delivery", id, "--method", tc.method}
			before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
			preview := integrationCLIOK(t, binary, f.Parent, "", append(args, "--check")...)
			if deliveryCLIString(t, preview, "state") != "ready" {
				t.Fatalf("PR blocked: %s", preview["reasons"])
			}
			if deliveryCLISnapshot(t, f.R.WorkspaceRoot) != before {
				t.Fatal("PR preview mutated fixture")
			}
			operation := deliveryCLIString(t, preview, "id")
			var result map[string]json.RawMessage
			if tc.crash {
				if _, _, e = integrationCLIRun(binary, f.Parent, "pass\n", args...); e == nil {
					t.Fatal("crashed process returned success")
				}
				result = integrationCLIOK(t, binary, f.Parent, "", "resume", operation)
			} else {
				result = integrationCLIOK(t, binary, f.Parent, "pass\n", args...)
			}
			if deliveryCLIString(t, result, "state") != "completed" {
				t.Fatalf("PR not closed: %s", result["reasons"])
			}
			integrationCLIOK(t, binary, f.Parent, "", "resume", operation)
			state := workflowProviderDocument(t, model)
			if state["merge_count"] != float64(1) {
				t.Fatalf("merge repeated: %v", state["merge_count"])
			}
			if deliveryCLIString(t, result, "integrated_oid") != state["merge_oid"] {
				t.Fatal("receipt lost actual merge oid")
			}
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") != parentBefore {
				t.Fatal("PR merge changed local base checkout")
			}
			queueAfter, e := workspace.ListTaskQueue(f.D.Workspace, workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, false)
			if e != nil || queueBefore.QueueID != queueAfter.QueueID || queueBefore.Revision != queueAfter.Revision || !equal(queueBefore.Current, queueAfter.Current) || !equal(queueBefore.Pending, queueAfter.Pending) {
				t.Fatalf("PR closeout changed already-closed queue: %v", e)
			}
			if _, e = os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); !os.IsNotExist(e) {
				t.Fatal("PR Task source was not retired")
			}
			runGit(t, f.Parent, "gc", "--prune=now")
			if gitOutput(t, f.Parent, "cat-file", "-t", o.Delivery.Candidates[0].OID) != "commit" {
				t.Fatal("cleanup lost candidate after GC")
			}
		})
	}
}

func TestIntegrationCLIPRRemotePendingIsNotClosure(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f, _, id, model := integrationCLIPRFixture(t, binary, "squash", false, true)
	r, output, e := integrationCLIRun(binary, f.Parent, "pass\n", "--delivery", id, "--method", "squash")
	if e == nil || r == nil || deliveryCLIString(t, r, "state") != "remote_pending" {
		t.Fatalf("pending treated as merged: %v %s", e, output)
	}
	operation := deliveryCLIString(t, r, "id")
	for i := 0; i < 2; i++ {
		r, output, e = integrationCLIRun(binary, f.Parent, "", "resume", operation)
		if e == nil || r == nil || deliveryCLIString(t, r, "state") != "remote_pending" || string(r["integrated"]) != "false" {
			t.Fatalf("pending changed truth: %v %s", e, output)
		}
	}
	if state := workflowProviderDocument(t, model); state["merge_count"] != float64(1) {
		t.Fatal("pending merge was resubmitted")
	}
	if _, e = os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); e != nil {
		t.Fatal("pending merge removed source")
	}
	if receipt, e := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, "task"); e != nil || receipt != nil {
		t.Fatal("pending merge started closeout")
	}
}
