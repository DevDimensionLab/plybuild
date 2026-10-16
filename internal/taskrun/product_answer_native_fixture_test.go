//go:build darwin

package taskrun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"golang.org/x/sys/unix"
)

// These exports exist only in the test build, permitting the external-package
// test to compose the real week NativeAdapter without an import cycle. Native
// Task/Goal/preparation/worker stores are real. The owned provider is explicitly
// synthetic and cannot attest Codex, login or OS sandbox qualification.
type SL01NativeFixture struct {
	D                        Dependencies
	Root, RunID, Task, Turns string
	f                        deliveryFixture
}

func sl01ShellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }

func NewSL01NativeFixture(t *testing.T, parent func(Dependencies, string, string, workspace.TaskGoalRef) *ControlContract, worker ...*WorkerContract) *SL01NativeFixture {
	t.Helper()
	var control *ControlContract
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic SL01 composition", TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &a, BeforePrepare: func(w workspace.Dependencies, root, epic string, goal workspace.TaskGoalRef) {
		control = parent(SystemDependencies(w), root, epic, goal)
	}})
	root := f.R.WorkspaceRoot
	if e := os.Chmod(root, 0700); e != nil {
		t.Fatal(e)
	}
	home := filepath.Join(root, "synthetic-codex-home")
	if e := privateDir(home); e != nil {
		t.Fatal(e)
	}
	t.Setenv("CODEX_HOME", home)
	python, e := exec.LookPath("python3")
	if e != nil {
		t.Fatal(e)
	}
	python, e = filepath.EvalSymlinks(python)
	if e != nil {
		t.Fatal(e)
	}
	provider := "#!" + python + "\n" + sl01OwnedProvider
	if e = os.WriteFile(f.R.Runtime.Executable.Path, []byte(provider), 0700); e != nil {
		t.Fatal(e)
	}
	testBinary, e := os.Executable()
	if e != nil {
		t.Fatal(e)
	}
	launcher := "#!/bin/sh\nset -eu\ncase \"$3\" in\n_supervise) printf '%s' \"$$\" > " + sl01ShellQuote(filepath.Join(root, "synthetic-supervisor.pid")) + "; exit 0;;\n_attempt) exec " + sl01ShellQuote(testBinary) + " -test.run='^TestSL01NativeAttemptProcess$' -test.timeout=60s -- \"$@\";;\n*) exit 91;;\nesac\n"
	if e = os.WriteFile(f.R.Runtime.PlyExecutable.Path, []byte(launcher), 0700); e != nil {
		t.Fatal(e)
	}
	f.R.Envelope = deliveryEnv("headless-run-request")
	f.R.Worker = &WorkerContract{SchemaVersion: 2, MaxAttempts: 3, AgentTimeoutSeconds: 30, MaxProviderRetries: 2, ProductQuestions: &ProductQuestionContract{1, BoundAnswerMode, 1}}
	if len(worker) == 1 {
		copy := *worker[0]
		f.R.Worker = &copy
	}
	f.R.ControlContract = control
	f.R.Herdr = (WorkflowRequest{}).Herdr
	f.R.Runtime.Mode = "headless"
	f.R.Runtime.Executable.SHA256 = hashFileTest(t, f.R.Runtime.Executable.Path)
	f.R.Runtime.PlyExecutable.SHA256 = hashFileTest(t, f.R.Runtime.PlyExecutable.Path)
	f.R.Runtime.PermissionBinding.ProfileID = "ply-worker"
	f.R.HumanAuthority.StartSurface = "human_authorized_worker"
	f.R.Delivery.AllowLocalInstall = false
	f.D.ProductAnswerQualification = sl01SyntheticQualification
	f.R.HandoffDraft, e = workflowhandoff.BuildDeliveryHandoffDraft(f.D.Workflow, f.R.PreparationID, "synthetic human", f.R.Delivery.OwnerClaim, "fixture/delivery", f.AcceptancePath, workflowhandoff.DeliveryDraftAuthority{AllowLocalInstall: false, Agreement: f.Prepared.Delivery})
	if e != nil {
		t.Fatal(e)
	}
	f.File = writeAny(t, root, "headless-request.json", f.R)
	o := workflowTestStart(t, f.workflowFixture)
	// StartWorker releases its short-lived synthetic supervisor. Reap exactly
	// that child; never signal a process found through a name or stale PID.
	var pid int
	for until := time.Now().Add(3 * time.Second); time.Now().Before(until); {
		raw, er := os.ReadFile(filepath.Join(root, "synthetic-supervisor.pid"))
		if er == nil {
			pid, _ = strconv.Atoi(string(raw))
			if pid > 0 {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	if pid == 0 {
		t.Fatal("synthetic supervisor did not publish its owned PID")
	}
	var status unix.WaitStatus
	if _, e = unix.Wait4(pid, &status, 0, nil); e != nil {
		t.Fatal(e)
	}
	lock, e := workerLock(filepath.Join(o.Paths.RunRoot, "worker", "owner.lock"), true)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { lock.Close() })
	f.D.WorkerOwner = &workerOwner{root: root, runID: o.RunID, generation: 1, lock: lock}
	if e = workflowUpdate(f.D, root, o.RunID, func(s *workflowState) error {
		s.Result.Worker.Generation = 1
		s.Result.Worker.Phase = "running"
		return nil
	}); e != nil {
		t.Fatal(e)
	}
	acceptance := deliveryTestAcceptance(t, f, o)
	acceptance.RuntimeClaim.RuntimeID = ptr("ply-worker")
	acceptance.RuntimeClaim.NativeSessionID = ptr(o.Worker.SessionID)
	if _, e = WorkflowAccept(f.D, root, o.RunID, o.Paths.Context, writeAny(t, root, "worker-test-acceptance.json", acceptance)); e != nil {
		t.Fatal(e)
	}
	return &SL01NativeFixture{D: f.D, Root: root, RunID: o.RunID, Task: f.Prepared.Preparation.Preparation.Plan.WorktreePath, Turns: filepath.Join(home, "model-turns.jsonl"), f: f}
}

func sl01SyntheticQualification(Runtime, string) (string, error) {
	return "sha256:" + strings.Repeat("f", 64), nil
}

func TestSL01NativeAttemptProcess(t *testing.T) {
	args := os.Args
	cut := -1
	for i, a := range args {
		if a == "--" {
			cut = i
			break
		}
	}
	if cut < 0 {
		return
	}
	args = args[cut+1:]
	if len(args) != 8 || strings.Join(args[:3], " ") != "workflow worker _attempt" || args[4] != "--workspace" || args[6] != "--job" {
		t.Fatalf("unexpected isolated attempt args %q", args)
	}
	d := SystemDependencies(workspace.SystemDependencies())
	d.ProductAnswerQualification = sl01SyntheticQualification
	if e := RunWorkerAttempt(d, args[5], args[3], args[7]); e != nil {
		t.Fatal(e)
	}
}

func (f *SL01NativeFixture) Step(t *testing.T) WorkflowRun {
	t.Helper()
	s, e := workflowRead(f.Root, f.RunID)
	if e != nil {
		t.Fatal(e)
	}
	handled, e := workerProductWaiting(f.D, s)
	if e == nil && !handled {
		e = workerStep(f.D, s)
	}
	if e != nil {
		t.Fatal(e)
	}
	o, e := ShowWorker(f.D, f.Root, f.RunID)
	if e != nil {
		t.Fatal(e)
	}
	return o
}
func (f *SL01NativeFixture) AssertNativeSource(t *testing.T) {
	t.Helper()
	s, e := workflowRead(f.Root, f.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = workerPreparation(f.D, s); e != nil {
		t.Fatal(e)
	}
	if e = workerPreparationCurrent(f.D, s); e != nil {
		t.Fatal(e)
	}
	q, e := productQuestionRecordRead(s)
	if e != nil {
		t.Fatal(e)
	}
	if q.Scope.ProjectID != "ply" || q.Scope.RepoID != "ply" || q.Scope.TaskID != "task" || q.Scope.GoalID != "solution" || q.Scope.GoalRevision < 1 || q.Scope.GoalSHA256 == "" || q.Scope.TaskWorktree != f.Task {
		t.Fatal(q.Scope)
	}
	r, e := workerRunReserved(f.D, s)
	if e != nil {
		t.Fatal(e)
	}
	if r.Process == nil || r.Process.PID == nil || r.Process.StartIdentity == nil || r.ThreadID != sl01Thread || !quiescent(*r.Process) {
		t.Fatalf("source not an actually reaped owned stand-in: %+v", r)
	}
}
func (f *SL01NativeFixture) TurnCount(t *testing.T) int {
	t.Helper()
	raw, e := os.ReadFile(f.Turns)
	if e != nil {
		t.Fatal(e)
	}
	return len(strings.Fields(string(raw)))
}

// Mutate only this owned stand-in's actual session artifact, restoring the
// exact bytes before the positive flow continues. No provider is launched by
// either boundary under test.
func (f *SL01NativeFixture) AssertSourceSessionDriftBlocks(t *testing.T, request BoundAnswerRequest, boundary string) {
	t.Helper()
	s, e := workflowRead(f.Root, f.RunID)
	if e != nil {
		t.Fatal(e)
	}
	q, e := productQuestionRecordRead(s)
	if e != nil {
		t.Fatal(e)
	}
	original, e := os.ReadFile(q.Session.Path)
	if e != nil {
		t.Fatal(e)
	}
	for _, change := range []string{"prefix", "appended_turn"} {
		t.Run(change, func(t *testing.T) {
			drift := append([]byte{}, original...)
			if change == "prefix" {
				drift = bytes.Replace(drift, []byte(sl01Thread), []byte("01a145a7-9221-7ac0-a56d-d64f4c7bab12"), 1)
			} else {
				drift = append(drift, []byte("{\"type\":\"turn_context\",\"payload\":{}}\n")...)
			}
			if e := os.WriteFile(q.Session.Path, drift, 0600); e != nil {
				t.Fatal(e)
			}
			defer func() {
				if e := os.WriteFile(q.Session.Path, original, 0600); e != nil {
					t.Error(e)
				}
			}()
			starts := 0
			var blocked error
			switch boundary {
			case "apply":
				_, blocked = ApplyBoundAnswer(f.D, request)
				if _, er := os.Lstat(productClaimPath(f.Root, f.RunID, request.QuestionID)); !os.IsNotExist(er) {
					t.Errorf("drift published a winning claim: %v", er)
				}
			case "start":
				in, er := workerIntent(s, *s.ProductContinuation)
				if er != nil {
					t.Fatal(er)
				}
				blocked = workerProductFinalStart(f.D, s, in, func() error { starts++; return nil })
				current, er := workflowRead(f.Root, f.RunID)
				if er != nil || current.ProductLaunch != nil || current.ProductContinuation == nil || *current.ProductContinuation != *s.ProductContinuation {
					t.Errorf("rejected source drift changed launch/reservation: %+v %v", current.ProductLaunch, er)
				}
			default:
				t.Fatal("unsupported test boundary")
			}
			if blocked == nil || starts != 0 || f.TurnCount(t) != 1 {
				t.Fatalf("source %s drift passed %s: starts=%d err=%v", change, boundary, starts, blocked)
			}
		})
	}
}

func (f *SL01NativeFixture) AssertResumedContextDriftBlocksCompletion(t *testing.T, requestID string) {
	t.Helper()
	s, e := workflowRead(f.Root, f.RunID)
	if e != nil {
		t.Fatal(e)
	}
	q, e := productQuestionRecordRead(s)
	if e != nil {
		t.Fatal(e)
	}
	original, e := os.ReadFile(q.Session.Path)
	if e != nil {
		t.Fatal(e)
	}
	for _, field := range []string{"model", "effort", "approval_policy", "permission_profile"} {
		t.Run(field, func(t *testing.T) {
			lines := bytes.Split(bytes.TrimSuffix(original, []byte{'\n'}), []byte{'\n'})
			contexts := 0
			for i, line := range lines {
				var item map[string]any
				if e := json.Unmarshal(line, &item); e != nil {
					t.Fatal(e)
				}
				if item["type"] != "turn_context" {
					continue
				}
				contexts++
				if contexts != 2 {
					continue
				}
				payload := item["payload"].(map[string]any)
				if field == "permission_profile" {
					payload[field].(map[string]any)["network"] = "enabled"
				} else {
					payload[field] = "different-from-bound-source"
				}
				lines[i], e = Canonical(item)
				if e != nil {
					t.Fatal(e)
				}
			}
			if contexts != 2 {
				t.Fatalf("expected actual source and resumed contexts, got %d", contexts)
			}
			if e = os.WriteFile(q.Session.Path, append(bytes.Join(lines, []byte{'\n'}), '\n'), 0600); e != nil {
				t.Fatal(e)
			}
			defer func() {
				if e := os.WriteFile(q.Session.Path, original, 0600); e != nil {
					t.Error(e)
				}
			}()
			status, er := ShowBoundAnswer(f.D, f.Root, requestID)
			if er != nil || status.State != "outcome_unknown" || status.Receipt == nil || status.Receipt.ResultSHA256 != nil || status.Receipt.ConsumedAnswerRequestID != nil || f.TurnCount(t) != 2 {
				t.Fatalf("changed resumed %s retained a completed claim: %+v %v", field, status, er)
			}
		})
	}
}

func (f *SL01NativeFixture) StartHeldContinuation(t *testing.T) func() {
	t.Helper()
	home := filepath.Dir(f.Turns)
	if e := os.WriteFile(filepath.Join(home, "hold-resume"), nil, 0600); e != nil {
		t.Fatal(e)
	}
	done := make(chan struct{})
	go func() { defer close(done); f.Step(t) }()
	released := false
	release := func() {
		if released {
			return
		}
		released = true
		if e := os.WriteFile(filepath.Join(home, "release-resume"), nil, 0600); e != nil {
			t.Error(e)
		}
		select {
		case <-done:
		case <-time.After(8 * time.Second):
			t.Error("owned continuation did not return after release")
		}
	}
	t.Cleanup(release)
	return release
}

func (f *SL01NativeFixture) AssertMissingAnswerBlocksReviewer(t *testing.T) {
	t.Helper()
	s, e := workflowRead(f.Root, f.RunID)
	if e != nil {
		t.Fatal(e)
	}
	q, e := productQuestionRecordRead(s)
	if e != nil {
		t.Fatal(e)
	}
	path := productClaimPath(f.Root, f.RunID, q.Question.QuestionID)
	if e = os.Rename(path, path+".held"); e != nil {
		t.Fatal(e)
	}
	defer func() {
		if e := os.Rename(path+".held", path); e != nil {
			t.Fatal(e)
		}
	}()
	d := f.D
	reached := false
	d.Fault = func(point string) error {
		if point == "worker_before_attempt_reservation" {
			reached = true
			return fmt.Errorf("test stops before launch")
		}
		return nil
	}
	if e = workerPrepareAttempt(d, s); e == nil || reached {
		t.Errorf("missing winner reached reviewer reservation: reached=%v err=%v", reached, e)
	}
}

func (f *SL01NativeFixture) AssertReviewerChoice(t *testing.T) {
	t.Helper()
	s, e := workflowRead(f.Root, f.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if s.Result.Worker.Job == nil || s.Result.Worker.Job.Kind != "reviewer" {
		t.Fatal(s.Result.Worker)
	}
	in, e := workerIntent(s, *s.Result.Worker.Job)
	if e != nil || in.Prompt == nil {
		t.Fatal(in, e)
	}
	raw, e := workflowBound(*in.Prompt, 4<<20)
	if e != nil {
		t.Fatal(e)
	}
	q, e := productQuestionRecordRead(s)
	if e != nil {
		t.Fatal(e)
	}
	c, e := productClaimForQuestion(f.Root, f.RunID, q.Question.QuestionID)
	if e != nil {
		t.Fatal(e)
	}
	result, _, e := deliveryReadBinding(filepath.Join(s.ProductContinuation.Directory, "result.json"), 1<<20)
	if e != nil {
		t.Fatal(e)
	}
	for _, required := range []string{q.Question.QuestionID, q.Question.QuestionSHA256, c.Request.RequestID, c.RequestSHA256, digest(c), "task_first", "Task first", "Tasks precede commands", result.SHA256, s.Result.Worker.CandidateOID, s.Result.Worker.CandidateTree, "verify the selected product choice"} {
		if !strings.Contains(string(raw), required) {
			t.Errorf("actual reviewer intent omits required choice context %q", required)
		}
	}
}

func (f *SL01NativeFixture) AssertCorrectionRetainsChoiceAcrossOwnerRestart(t *testing.T) {
	t.Helper()
	s, e := workflowRead(f.Root, f.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if e = workerCorrection(f.D, s, "Synthetic reviewer requires an ordinary bounded correction"); e != nil {
		t.Fatal(e)
	}
	s, e = workflowRead(f.Root, f.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if e = f.D.WorkerOwner.lock.Close(); e != nil {
		t.Fatal(e)
	}
	lock, e := workerLock(filepath.Join(workerDirectory(s), "owner.lock"), false)
	if e != nil {
		t.Fatal(e)
	}
	defer lock.Close()
	f.D.WorkerOwner = nil
	if e = workerClaimOwner(&f.D, s, lock); e != nil {
		t.Fatalf("ordinary correction could not recover its native owner: %v", e)
	}
	s, e = workflowRead(f.Root, f.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if handled, e := workerProductWaiting(f.D, s); handled || e != nil {
		t.Fatalf("ordinary correction was absorbed into the prior answer turn: %v %v", handled, e)
	}
	d := f.D
	reached := false
	d.Fault = func(point string) error {
		if point == "worker_before_attempt_reservation" {
			reached = true
			return fmt.Errorf("bounded prompt inspection before new developer launch")
		}
		return nil
	}
	if e = workerPrepareAttempt(d, s); e == nil || !reached {
		t.Fatalf("correction failed before its normal dispatch boundary: %v", e)
	}
	path := filepath.Join(workerDirectory(s), "attempts", "developer-0002", "prompt.txt")
	raw, e := readFile(path, 4<<20, true)
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(string(raw), "task_first") || !strings.Contains(string(raw), "verify the selected product choice") || strings.Contains(string(raw), "permits one bounded product choice") {
		t.Fatal("correction forgot the answer or reopened product authority")
	}
}

func (f *SL01NativeFixture) AssertSingleAttemptForbidsCorrection(t *testing.T) {
	t.Helper()
	s, e := workflowRead(f.Root, f.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if s.Request.Worker.MaxAttempts != 1 || s.Request.Worker.MaxProviderRetries != 0 || s.Result.Worker.ProviderRetries != 0 || f.TurnCount(t) != 3 {
		t.Fatalf("one-attempt contract or three actual launches changed: %+v %+v", s.Request.Worker, s.Result.Worker)
	}
	if e = workerCorrection(f.D, s, "No additional correction is authorized"); e == nil {
		t.Fatal("single-attempt contract allowed another developer iteration")
	}
	t.Log("actual owned launches=3: initial + approved same-session answer + separate reviewer; extra retries=0; correction forbidden")
}
func (f *SL01NativeFixture) Export(t *testing.T, name string, value any) {
	t.Helper()
	if dir := os.Getenv("PLY_SL01_CORE_EXPORT"); dir != "" {
		if !filepath.IsAbs(dir) {
			t.Fatal("export must be an explicit absolute test evidence directory")
		}
		raw, e := Canonical(value)
		if e != nil {
			t.Fatal(e)
		}
		p := filepath.Join(dir, name+".json")
		fd, e := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e != nil {
			t.Fatal(e)
		}
		_, e = fd.Write(append(raw, '\n'))
		ce := fd.Close()
		if e != nil {
			t.Fatal(e)
		}
		if ce != nil {
			t.Fatal(ce)
		}
	}
}

const sl01OwnedProvider = `import json, os, pathlib, sys, time, tomllib
args=sys.argv[1:]
if 'mcp' in args:
 print('[]'); raise SystemExit(0)
if 'features' in args:
 for name in ['hooks','multi_agent','apps','plugins','remote_plugin','browser_use','browser_use_external','browser_use_full_cdp_access','computer_use','in_app_browser','image_generation','skill_mcp_dependency_install','skill_search']: print(name+' stable false')
 raise SystemExit(0)
assert 'exec' in args and '-o' in args
prompt=sys.stdin.read()
answer_path=pathlib.Path(args[args.index('-o')+1]); job=answer_path.parent
intent=json.loads((job/'intent.json').read_text())
home=pathlib.Path(os.environ['CODEX_HOME'])
failure=home/'usage-limit-role'
if failure.exists() and failure.read_text()==intent['kind']:
 with (home/'model-turns.jsonl').open('a') as f: f.write('usage-limit:'+intent['kind']+'\n')
 for event in [{'type':'thread.started','thread_id':'synthetic-usage-limit'},{'type':'turn.started'},{'type':'error','code':'usage_limit_reached'},{'type':'turn.failed','error':{'code':'usage_limit_reached'}}]: print(json.dumps(event),flush=True)
 raise SystemExit(1)
if intent['kind']=='reviewer':
 answer_path.write_text(json.dumps({'decision':'passed','summary':'Synthetic independent reviewer','findings':[]}))
 with (home/'model-turns.jsonl').open('a') as f: f.write('review\n')
 for event in [{'type':'thread.started','thread_id':'synthetic-independent-reviewer'},{'type':'turn.started'},{'type':'turn.completed'}]: print(json.dumps(event),flush=True)
 raise SystemExit(0)
settings={}
for i,a in enumerate(args):
 if a=='-c': settings.update(tomllib.loads(args[i+1]))
policy=settings['permissions']['ply-worker']
entries=[]
for name,access in policy['filesystem'].items():
 path={'type':'special','value':{'kind':name[1:]}} if name.startswith(':') else {'type':'path','path':name}
 entries.append({'path':path,'access':access})
thread='01a145a7-9221-7ac0-a56d-d64f4c7bab11'
cwd=os.getcwd(); model=args[args.index('--model')+1]
context={'type':'turn_context','payload':{'cwd':cwd,'approval_policy':'never','model':model,'effort':settings['model_reasoning_effort'],'permission_profile':{'type':'managed','network':'restricted','file_system':{'type':'restricted','entries':entries}}}}
home=pathlib.Path(os.environ['CODEX_HOME']); session=home/'sessions'/'2026'/'10'/'16'/('rollout-'+thread+'.jsonl')
session.parent.mkdir(parents=True,exist_ok=True)
resume=args[args.index('exec')+1]=='resume'
for event in [{'type':'thread.started','thread_id':thread},{'type':'turn.started'}]: print(json.dumps(event),flush=True)
if resume:
 assert args[-2:]==[thread,'-'] and '--last' not in args and '--fork' not in args
 if (home/'hold-resume').exists():
  until=time.monotonic()+8
  while not (home/'release-resume').exists():
   assert time.monotonic()<until, 'owned fixture continuation barrier timed out'
   time.sleep(0.01)
 claim=json.loads(pathlib.Path(intent['bound_answer']['locator']).read_text())
 q=json.loads(pathlib.Path(claim['question_record']['locator']).read_text())
 source=json.loads((pathlib.Path(q['source_job']['directory'])/'intent.json').read_text())
 assert intent['argv'][:intent['argv'].index('exec')]==source['argv'][:source['argv'].index('exec')]
 assert intent['environment']==source['environment']
 with session.open('a') as f: f.write(json.dumps(context)+'\n')
 pathlib.Path('reading-order.txt').write_text(claim['request']['option_id']+'\n')
 pathlib.Path('test').mkdir(exist_ok=True)
 pathlib.Path('test/delivery_acceptance.sh').write_text('#!/bin/sh\nset -eu\ntest -x "$PLY_CANDIDATE_BINARY"\ntest -s reading-order.txt\n')
 candidate=pathlib.Path(q['source_job']['directory'])/'output'/'candidate'; candidate.write_text('#!/bin/sh\nexit 0\n'); candidate.chmod(0o700)
 answer={'status':'completed','summary':'Applied the exact selected product choice','block_reason':'','consumed_answer_request_id':claim['request']['request_id'],'consumed_answer_request_sha256':claim['request_sha256']}
else:
 with session.open('x') as f:
  f.write(json.dumps({'type':'session_meta','payload':{'id':thread,'cwd':cwd}})+'\n'); f.write(json.dumps(context)+'\n')
 answer={'status':'needs_product_answer','summary':'Both reading orders satisfy the Goal','block_reason':'','product_question':{'purpose':'product_choice','text':'Which reading order?','reason':'Choose the intended product order','options':[{'option_id':'task_first','label':'Task first','consequence':'Tasks precede commands'},{'option_id':'command_first','label':'Commands first','consequence':'Commands precede tasks'}],'continuation':{'type':'resume_same_developer_thread','source_iteration':1,'source_job_id':job.name}}}
with (home/'model-turns.jsonl').open('a') as f: f.write(('resume' if resume else 'initial')+'\n')
answer_path.write_text(json.dumps(answer))
print(json.dumps({'type':'turn.completed','usage':{'input_tokens':1,'cached_input_tokens':0,'output_tokens':1}}),flush=True)
`
