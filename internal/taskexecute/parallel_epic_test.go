package taskexecute

import (
	"encoding/json"
	"errors"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Simulated visible panes/provider sessions. Every unexpected transport command
// fails; no real provider, Herdr, OS process inspection or browser is invoked.
const parallelHerdr = `import json,sys
from pathlib import Path
p=Path(__file__).parent/'parallel-model.json'
m=json.loads(p.read_text()) if p.exists() else {'panes':{},'agents':{},'starts':0,'prompts':0}
a=sys.argv[1:]; cmd=a[:2]
if cmd==['workspace','get']: r={'workspace':{'workspace_id':a[2]}}
elif cmd==['tab','create']:
 n=str(len(m['panes'])+1); pane={'workspace_id':'w-fixture','tab_id':'tab-'+n,'pane_id':'pane-'+n,'terminal_id':'terminal-'+n,'cwd':a[a.index('--cwd')+1],'foreground_cwd':a[a.index('--cwd')+1]}; m['panes'][pane['pane_id']]=pane; r={'root_pane':pane}
elif cmd==['tab','rename']: r={'ok':True}
elif cmd==['pane','list']: r={'panes':list(m['panes'].values())}
elif cmd==['pane','get']: r={'pane':m['panes'][a[2]]}
elif cmd==['pane','process-info']: r={'process_info':{'pane_id':a[2],'shell_pid':531,'foreground_process_group_id':531,'foreground_processes':[{'pid':531,'name':'zsh','argv0':'zsh','argv':['-zsh'],'cwd':m['panes'][a[2]]['cwd']}]}}
elif cmd==['agent','start']:
 pane=m['panes'][a[a.index('--pane')+1]]; agent={**pane,'name':a[2],'agent':'codex','agent_status':'idle','agent_session':{'agent':'codex','kind':'id','value':'session-'+pane['pane_id']}}; m['agents'][pane['pane_id']]=agent; m['starts']+=1; r={'agent':agent}
elif cmd==['agent','get']:
 if a[2] not in m['agents']:
  print(json.dumps({'error':{'code':'agent_not_found','message':'simulated unstarted agent'}}),file=sys.stderr);raise SystemExit(1)
 r={'agent':m['agents'][a[2]]}
elif cmd==['agent','prompt']:
 m['prompts']+=1; m['agents'][a[2]]['agent_status']='working'; r={'acknowledged':True}
else: raise SystemExit('UNEXPECTED SIMULATED COMMAND '+repr(a))
p.write_text(json.dumps(m)); print(json.dumps({'result':r}))
`

func parallelLauncherFixture(t *testing.T, candidate ...string) (taskrun.Dependencies, Input, string) {
	t.Helper()
	a := workspace.DeliveryAgreement{SchemaVersion: 1, Mode: workspace.DeliveryLocalEpic}
	var d taskrun.Dependencies
	var in Input
	var root, epic string
	if len(candidate) > 0 {
		d, in, root, epic = launcherFixtureWithCLI(t, "codex", candidate[0], a)
		d.Executable = func() (string, error) { return candidate[0], nil }
	} else {
		d, in, root, epic = launcherFixtureForProvider(t, "codex", a)
	}
	d.LocalProcesses = func(string, string) ([]taskrun.LocalProcess, error) {
		t.Fatal("unexpected real process inspection fallback")
		return nil, nil
	}
	d.StartupProcessTree = func(pid int) ([]taskrun.StartupProcess, error) {
		return []taskrun.StartupProcess{{PID: pid, ParentPID: 1, ProcessGroupID: pid}}, nil
	}
	python, err := exec.LookPath("python3")
	if err != nil {
		t.Fatal(err)
	}
	python, err = filepath.EvalSymlinks(python)
	if err != nil || !filepath.IsAbs(python) {
		t.Fatalf("absolute fixture Python: %v", err)
	}
	if err = os.WriteFile(in.Runtime.HerdrExecutable, []byte("#!"+python+"\n"+parallelHerdr), 0700); err != nil {
		t.Fatal(err)
	}
	w := d.Workspace
	r, err := w.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	r.ParallelBoundary = &workspace.RegistryFormatBoundary{FromVersion: 4, FromSHA256: r.RawSHA256, ManifestSHA256: "sha256:" + strings.Repeat("a", 64), ReaderVersion: "v0.0.0-fixture", ReaderCommit: fixtureGit(t, epic, "rev-parse", "HEAD"), ReaderSHA256: "sha256:" + strings.Repeat("b", 64)}
	r.FormatVersion = workspace.ParallelRegistryFormat
	r.SourceReconciliations = []workspace.TaskSourceReconciliation{}
	raw, err := yaml.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(root, ".ply", "work-items.yaml"), raw, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = workspace.CreateTask(w, workspace.TaskCreateInput{TaskID: "goal-b", Title: "B", Description: "B", ParentEpicID: "epic", ProjectID: "fixture", RepoID: "fixture"}); err != nil {
		t.Fatal(err)
	}
	r, err = w.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	var b map[string]any
	raw, err = os.ReadFile(filepath.Join(root, "goal.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &b); err != nil {
		t.Fatal(err)
	}
	b["task_id"] = "goal-b"
	b["publication_key"] = "fixture/goal-b"
	b["spec_id"] = "observable-b"
	for _, v := range r.TaskProblemRevisions {
		if v.TaskID == "goal-b" {
			b["problem"] = workspace.TaskRevisionRef{Revision: v.Revision, ManifestSHA256: v.ManifestSHA256}
		}
	}
	file := filepath.Join(root, "goal-b.json")
	fixtureJSON(t, file, b)
	bs, err := workspace.RecordTaskSpec(w, workspace.TaskContentInput{TaskID: "goal-b", File: file})
	if err != nil {
		t.Fatal(err)
	}
	q, err := workspace.ListTaskQueue(w, workspace.QueueTargetInput{ProjectID: "fixture", RepoID: "fixture", EpicID: "epic"}, false)
	if err != nil {
		t.Fatal(err)
	}
	draft := workspace.WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: "fixture/two", ProjectID: "fixture", RepoID: "fixture", EpicID: "epic", ExpectedRevision: q.Revision, Entries: []workspace.TaskGoalQueueEntry{{TaskID: "goal", Goal: q.Pending[0].Goal}, {TaskID: "goal-b", Goal: &workspace.TaskGoalRef{SpecID: *bs.OutcomeRef.SpecID, Spec: workspace.TaskRevisionRef{Revision: *bs.OutcomeRef.Revision, ManifestSHA256: bs.OutcomeRef.ManifestSHA256}}}}, Recorder: workspace.TaskGoalQueueRecorder{ActorClaim: "simulated planner", ControlSurface: "isolated synthetic test", RecordedAtUTC: "2026-10-05T00:00:00Z"}}
	file = filepath.Join(root, "two-queue.json")
	fixtureJSON(t, file, draft)
	if _, err = workspace.SetTaskQueue(w, file); err != nil {
		t.Fatal(err)
	}
	in.Next = false
	in.SpecID = "observable-goal"
	in.Parallel = true
	return d, in, root
}

func TestParallelLaunchReservationFaultsDoNotReplay(t *testing.T) {
	for _, point := range []string{"workflow_after_reservation", "workflow_before_tab_send", "workflow_before_agent_send", "workflow_before_prompt_send"} {
		t.Run(point, func(t *testing.T) {
			d, in, root := parallelLauncherFixture(t)
			in.SpecID = "observable-b"
			b, err := Execute(d, in)
			if err != nil {
				t.Fatal(err)
			}
			in.SpecID = "observable-goal"
			d.Fault = func(at string) error {
				if at == point {
					return errors.New("simulated launch interruption")
				}
				return nil
			}
			first, err := Execute(d, in)
			if err == nil {
				t.Fatal("interruption not observed")
			}
			d.Fault = nil
			before, _ := os.ReadFile(filepath.Join(root, "bin", "parallel-model.json"))
			again, retryErr := Execute(d, in)
			after, _ := os.ReadFile(filepath.Join(root, "bin", "parallel-model.json"))
			var prior, now struct{ Starts, Prompts int }
			if json.Unmarshal(before, &prior) != nil || json.Unmarshal(after, &now) != nil {
				t.Fatal("lost simulated transport evidence")
			}
			if point != "workflow_after_reservation" && (prior.Starts != now.Starts || prior.Prompts != now.Prompts) {
				t.Fatalf("unknown launch replayed: before=%s after=%s", before, after)
			}
			if first.Run != nil && again.Run != nil && first.Run.RunID != again.Run.RunID {
				t.Fatal("retry replaced Run")
			}
			r, err := d.Workspace.WorkItems.Snapshot(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(r.TaskPreparations) != 2 {
				t.Fatal("launch retry duplicated preparation")
			}
			if b.Run == nil {
				t.Fatal("lost independent B")
			}
			t.Logf("simulated stage=%s prior starts/prompts=%d/%d retry=%d/%d err=%v nativeRun=%v", point, prior.Starts, prior.Prompts, now.Starts, now.Prompts, retryErr, again.Run != nil)
		})
	}
}

func TestParallelRejectsSharedProviderSessionBeforePrompt(t *testing.T) {
	d, in, root := parallelLauncherFixture(t)
	path := in.Runtime.HerdrExecutable
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	raw = []byte(strings.ReplaceAll(string(raw), "'session-'+pane['pane_id']", "'shared-simulated-session'"))
	if err = os.WriteFile(path, raw, 0700); err != nil {
		t.Fatal(err)
	}
	if _, err = Execute(d, in); err != nil {
		t.Fatal(err)
	}
	in.SpecID = "observable-b"
	_, err = Execute(d, in)
	if err == nil {
		t.Fatal("second parallel owner accepted another active owner's native session")
	}
	raw, err = os.ReadFile(filepath.Join(root, "bin", "parallel-model.json"))
	if err != nil {
		t.Fatal(err)
	}
	var model struct{ Starts, Prompts int }
	if err = json.Unmarshal(raw, &model); err != nil {
		t.Fatal(err)
	}
	if model.Prompts != 1 {
		t.Fatalf("Task prompt sent into shared session: %s", raw)
	}
}

func TestParallelNativeTwoOwners(t *testing.T) {
	d, in, root := parallelLauncherFixture(t)
	var raw []byte
	first, err := Execute(d, in)
	if err != nil {
		t.Fatalf("first native launch: %v", err)
	}
	in.SpecID = "observable-b"
	second, err := Execute(d, in)
	if err != nil {
		t.Fatalf("second native launch: %v", err)
	}
	if first.Run == nil || second.Run == nil || first.Run.RunID == second.Run.RunID || first.Run.Transport.AgentSessionID == second.Run.Transport.AgentSessionID || first.Run.Transport.PaneID == second.Run.Transport.PaneID || first.RequestPath == second.RequestPath {
		t.Fatal("parallel launch shared identity")
	}
	controls := map[string]bool{}
	for _, result := range []Result{first, second} {
		var request taskrun.WorkflowRequest
		raw, err = os.ReadFile(result.RequestPath)
		if err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal(raw, &request); err != nil {
			t.Fatal(err)
		}
		if request.SchemaVersion != 3 || request.Delivery.AdmissionMode != workspace.AdmissionParallel {
			t.Fatal("launch lost frozen parallel authority")
		}
		control := request.Runtime.PlyExecutable.Path
		if controls[control] || !strings.HasPrefix(control, filepath.Dir(result.RequestPath)+string(filepath.Separator)) {
			t.Fatal("owners share immutable control path")
		}
		controls[control] = true
		t.Logf("simulated native owner Run=%s pane=%s session=%s control=%s", result.Run.RunID, result.Run.Transport.PaneID, result.Run.Transport.AgentSessionID, control)
	}
	in.SpecID = "observable-goal"
	again, err := Execute(d, in)
	if err != nil || again.Run == nil || again.Run.RunID != first.Run.RunID {
		t.Fatalf("same A retry: %v", err)
	}
	var model struct{ Starts, Prompts int }
	raw, err = os.ReadFile(filepath.Join(root, "bin", "parallel-model.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &model); err != nil {
		t.Fatal(err)
	}
	if model.Starts != 2 || model.Prompts != 2 {
		t.Fatalf("duplicate native launch effects: %+v", model)
	}
}

func TestParallelExclusiveRoutesRefuseBeforeEffects(t *testing.T) {
	d, in, root := parallelLauncherFixture(t)
	if _, e := Execute(d, in); e != nil {
		t.Fatal(e)
	}
	in.SpecID = "observable-b"
	if _, e := Execute(d, in); e != nil {
		t.Fatal(e)
	}
	before, e := os.ReadFile(filepath.Join(root, ".ply", "work-items.yaml"))
	if e != nil {
		t.Fatal(e)
	}
	transport, e := os.ReadFile(filepath.Join(root, "bin", "parallel-model.json"))
	if e != nil {
		t.Fatal(e)
	}
	for _, mode := range []string{"next", "serial", "restart", "worker", "control"} {
		bad := in
		switch mode {
		case "next":
			bad.Next = true
			bad.SpecID = ""
			bad.Parallel = false
		case "serial":
			bad.Parallel = false
		case "restart":
			bad.Restart = true
		case "worker":
			bad.Worker = &taskrun.WorkerContract{SchemaVersion: 1, MaxAttempts: 1, AgentTimeoutSeconds: 30}
		case "control":
			bad.ControlContract = taskrun.NewControlContract("week", "synthetic", "sha256:"+strings.Repeat("c", 64))
		}
		if _, e = Execute(d, bad); e == nil {
			t.Fatalf("exclusive %s route adopted parallel work", mode)
		}
	}
	after, e := os.ReadFile(filepath.Join(root, ".ply", "work-items.yaml"))
	if e != nil || string(before) != string(after) {
		t.Fatal("exclusive route wrote registry")
	}
	after, e = os.ReadFile(filepath.Join(root, "bin", "parallel-model.json"))
	if e != nil || string(transport) != string(after) {
		t.Fatal("exclusive route launched/inspected provider")
	}
}
