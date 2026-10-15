//go:build !windows

package taskrun

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

type workerAttemptIntent struct {
	RunID           string       `json:"run_id"`
	JobID           string       `json:"job_id"`
	Kind            string       `json:"kind"`
	Generation      int          `json:"generation"`
	CandidateOID    string       `json:"candidate_oid"`
	CandidateTree   string       `json:"candidate_tree"`
	Argv            []string     `json:"argv"`
	Environment     []string     `json:"environment"`
	Prompt          *FileBinding `json:"prompt,omitempty"`
	Schema          *FileBinding `json:"schema,omitempty"`
	Review          *FileBinding `json:"review,omitempty"`
	CandidateBinary *FileBinding `json:"candidate_binary,omitempty"`
	TimeoutSeconds  int          `json:"timeout_seconds"`
}

type workerAttemptResult struct {
	RunID           string                 `json:"run_id"`
	JobID           string                 `json:"job_id"`
	IntentSHA256    string                 `json:"intent_sha256"`
	StartedAt       string                 `json:"started_at"`
	FinishedAt      string                 `json:"finished_at"`
	Exit            *int                   `json:"exit"`
	Error           string                 `json:"error,omitempty"`
	ThreadID        string                 `json:"thread_id,omitempty"`
	Stdout          FileBinding            `json:"stdout"`
	Stderr          FileBinding            `json:"stderr"`
	Answer          *FileBinding           `json:"answer,omitempty"`
	NativeCompleted bool                   `json:"native_completed"`
	Process         *Process               `json:"process,omitempty"`
	ProviderFailure *WorkerProviderFailure `json:"provider_failure,omitempty"`
}

func workerIntent(s workflowState, job WorkerJob) (workerAttemptIntent, error) {
	var in workerAttemptIntent
	if job.Directory != filepath.Join(workerDirectory(s), "attempts", job.ID) || !key(job.ID) {
		return in, fmt.Errorf("invalid worker attempt locator")
	}
	raw, err := workflowBound(FileBinding{Locator: filepath.Join(job.Directory, "intent.json"), SHA256: job.IntentHash}, 1<<20)
	if err != nil {
		return in, err
	}
	if err = decode(raw, 1<<20, &in); err != nil {
		return in, err
	}
	if in.RunID != s.Result.RunID || in.JobID != job.ID || in.Kind != job.Kind {
		return in, fmt.Errorf("worker attempt identity differs")
	}
	return in, nil
}

// The runner inherits the controller's open flock description. Closing a dead
// controller's descriptor cannot release ownership while this runner is alive.
func RunWorkerAttempt(d Dependencies, root, id, jobID string) error {
	// This command is a dedicated native child process. Provider-created answer
	// files must satisfy the same private-file contract as other native evidence.
	previousMask := unix.Umask(0077)
	defer unix.Umask(previousMask)
	s, err := workflowRead(root, id)
	if err != nil {
		return err
	}
	if err = workerValidateIdentity(s); err != nil {
		return err
	}
	job := s.Result.Worker.Job
	if job == nil || job.ID != jobID {
		return fmt.Errorf("worker attempt is not reserved")
	}
	in, err := workerIntent(s, *job)
	if err != nil {
		return err
	}
	if in.Generation != s.Result.Worker.Generation {
		return fmt.Errorf("stale worker attempt generation")
	}
	if in.Kind == "verify" || in.Kind == "integrate" {
		if err = workerNativeIntentMatches(s, in); err != nil {
			return err
		}
	}
	lock := os.NewFile(3, "inherited-worker-owner")
	if lock == nil {
		return fmt.Errorf("missing inherited worker ownership")
	}
	defer lock.Close()
	d.WorkerOwner = &workerOwner{root: root, runID: id, generation: in.Generation, lock: lock}
	if err = workerLiveOwner(d, s); err != nil {
		return err
	}
	// A separately opened descriptor must see an already held lock. Possession
	// of an arbitrary descriptor for this path is not a worker capability.
	if err = workerConfirmInheritedLock(lock, filepath.Join(workerDirectory(s), "owner.lock")); err != nil {
		return err
	}
	var e error
	if err = os.Chdir(s.Observed.Target.WorktreeLocator); err != nil {
		return err
	}
	resultPath := filepath.Join(job.Directory, "result.json")
	if _, err = os.Lstat(resultPath); err == nil {
		return fmt.Errorf("worker attempt already has a result")
	} else if !os.IsNotExist(err) {
		return err
	}
	r := workerAttemptResult{RunID: id, JobID: job.ID, IntentSHA256: job.IntentHash, StartedAt: d.Now().UTC().Format(time.RFC3339Nano)}
	var stdout, stderr workflowLimitedOutput
	if in.Kind == "commit" {
		err = workerCommit(d, s)
		r.NativeCompleted = err == nil
		if err != nil {
			r.Error = err.Error()
		}
	} else if in.Kind == "verify" || in.Kind == "integrate" {
		var o WorkflowRun
		if in.Kind == "verify" {
			if s.Result.Worker.Review == nil {
				return fmt.Errorf("review missing")
			}
			o, err = WorkflowDeliveryVerifyWithBinary(d, root, id, s.Result.Paths.Context, s.Result.Worker.Review.Locator, s.Result.Worker.CandidateBinary)
			r.NativeCompleted = err == nil && o.Delivery != nil && o.Delivery.Phase == "automatic_acceptance_passed"
		} else {
			o, err = WorkflowDeliveryIntegrate(d, root, id, s.Result.Paths.Context)
			r.NativeCompleted = err == nil && o.Delivery != nil && o.Delivery.Phase == "completed"
		}
		data, _ := Canonical(o)
		stdout.Write(data)
		if err != nil {
			r.Error = err.Error()
		}
	} else {
		if len(in.Argv) == 0 || in.TimeoutSeconds < 1 {
			return fmt.Errorf("invalid worker command intent")
		}
		if in.Kind == "developer" || in.Kind == "reviewer" {
			end := -1
			for i, a := range in.Argv {
				if a == "exec" {
					end = i
					break
				}
			}
			if end < 1 {
				return fmt.Errorf("invalid provider exec intent")
			}
			_, proof, e := workerToolPolicy(s, in.Argv[:end], false)
			if e != nil {
				return e
			}
			if _, e = workflowKeep(d, filepath.Join(job.Directory, "tool-policy.json"), proof); e != nil {
				return e
			}
		}
		var stdin *os.File
		if in.Prompt != nil {
			_, e := workflowBound(*in.Prompt, 4<<20)
			if e != nil {
				return e
			}
			stdin, e = physicalOpenFile(in.Prompt.Locator, os.O_RDONLY|unix.O_NOFOLLOW, 0)
			if e != nil {
				return e
			}
			defer stdin.Close()
		}
		if in.Schema != nil {
			if _, e = workflowBound(*in.Schema, 1<<20); e != nil {
				return e
			}
		}
		program := s.Request.Runtime.Executable
		if in.Kind == "pretest" {
			b, e := readFile(in.Argv[0], 128<<20, false)
			if e != nil {
				return e
			}
			program = Executable{Path: in.Argv[0], SHA256: hash(b)}
		}
		completion := ProviderCompletion{}
		process, runErr := workerOwnedExec(LaunchSpec{Executable: program, Argv: in.Argv, CWD: s.Observed.Target.WorktreeLocator, Timeout: time.Duration(in.TimeoutSeconds) * time.Second, StreamsRoot: filepath.Join(job.Directory, "streams"), Completion: &completion, Environment: in.Environment, Stdin: stdin}, func(p Process) error { return d.writeValue(filepath.Join(job.Directory, "process.json"), p) })
		r.Process = &process
		r.Exit = process.ExitCode
		r.Error = workerAttemptError(in.Kind, runErr)
		if !quiescent(process) || process.State != "exited" || !completion.StreamComplete || completion.StdoutTruncated || completion.StderrTruncated || completion.TimedOut {
			r.Error = "attempt process or output is not positively complete; preserve this Task and inspect process.json and streams"
		}
		if completion.Stdout.Locator != "" {
			raw, e := workflowBound(completion.Stdout, providerStreamLimit)
			if e != nil {
				return e
			}
			stdout.data = raw
		}
		if completion.Stderr.Locator != "" {
			raw, e := workflowBound(completion.Stderr, providerStreamLimit)
			if e != nil {
				return e
			}
			stderr.data = raw
		}
		if in.Kind == "developer" || in.Kind == "reviewer" {
			if !completion.SequenceValid || completion.TurnsCompleted < 1 {
				r.Error = "provider event sequence did not record a completed turn"
			}
			r.ThreadID, e = workerThreadID(stdout.data)
			if e != nil {
				r.Error = e.Error()
			}
			binding, _, e := deliveryReadBinding(filepath.Join(job.Directory, "answer.json"), 1<<20)
			if e != nil {
				r.Error = "provider answer missing or invalid: " + e.Error()
			} else {
				r.Answer = &binding
			}
			// Only complete, quiescent attempts may enter automatic waiting.
			// Tool output and truncated streams never authorize a retry.
			if quiescent(process) && process.State == "exited" && process.ExitCode != nil && *process.ExitCode != 0 && completion.StreamComplete && !completion.StdoutTruncated && !completion.StderrTruncated && !completion.TimedOut {
				r.ProviderFailure = workerProviderFailure(stdout.data, d.Now())
			}
		}
	}
	r.FinishedAt = d.Now().UTC().Format(time.RFC3339Nano)
	r.Stdout = FileBinding{Locator: filepath.Join(job.Directory, "stdout.txt"), SHA256: hash(stdout.data)}
	r.Stderr = FileBinding{Locator: filepath.Join(job.Directory, "stderr.txt"), SHA256: hash(stderr.data)}
	if err = d.writeOnce(r.Stdout.Locator, stdout.data); err != nil {
		return err
	}
	if err = d.writeOnce(r.Stderr.Locator, stderr.data); err != nil {
		return err
	}
	return d.writeValue(resultPath, r)
}

func workerConfirmInheritedLock(lock *os.File, path string) error {
	probe, err := workerLock(path, false)
	if err == nil {
		probe.Close()
		return fmt.Errorf("inherited ownership is not held")
	}
	if !errors.Is(err, unix.EWOULDBLOCK) {
		return err
	}
	if err = unix.Flock(int(lock.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		return fmt.Errorf("descriptor does not own the inherited lock: %w", err)
	}
	return nil
}

func workerAttemptError(kind string, err error) string {
	if err == nil {
		return ""
	}
	var exit processExitError
	if kind == "pretest" && errors.As(err, &exit) {
		return ""
	}
	return err.Error()
}

func workerThreadID(raw []byte) (string, error) {
	seq := parseProviderEvents(raw)
	if !seq.Valid || seq.ThreadID == nil || !plain(*seq.ThreadID, 1, 256) || seq.Started != 1 || seq.Completed != 1 {
		return "", fmt.Errorf("provider JSONL does not preserve one actual thread and completed turn: %s", seq.Reason)
	}
	return *seq.ThreadID, nil
}

func workerSealCandidate(path string) error {
	// Called only after a positively completed developer process. Normalize the
	// private output's mode without changing its bytes or following symlinks.
	f, err := physicalOpenFile(path, os.O_RDONLY|unix.O_NOFOLLOW, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || info.Mode().Perm()&0111 == 0 || info.Size() > 64<<20 {
		return fmt.Errorf("candidate must be a bounded regular executable")
	}
	if err = f.Chmod(0700); err != nil {
		return err
	}
	_, _, err = deliveryReadBinding(path, 64<<20)
	return err
}

func workerRunReserved(d Dependencies, s workflowState) (workerAttemptResult, error) {
	var r workerAttemptResult
	job := s.Result.Worker.Job
	if job == nil {
		return r, fmt.Errorf("worker attempt missing")
	}
	if _, err := workerIntent(s, *job); err != nil {
		return r, err
	}
	path := filepath.Join(job.Directory, "result.json")
	err := readValue(path, 1<<20, &r)
	if os.IsNotExist(err) {
		return r, fmt.Errorf("attempt %s has no durable result; it will not be replayed or replaced by a second writer", job.ID)
	}
	if err != nil {
		return r, err
	}
	if r.RunID != s.Result.RunID || r.JobID != job.ID || r.IntentSHA256 != job.IntentHash || r.FinishedAt == "" {
		return r, fmt.Errorf("worker attempt result binding differs")
	}
	if !workerNativeKind(job.Kind) {
		if r.Process == nil || r.Process.State != "exited" || !quiescent(*r.Process) || r.Process.ExitCode == nil || r.Exit == nil || *r.Process.ExitCode != *r.Exit {
			return r, fmt.Errorf("attempt %s has no positive process completion; pause, replacement and further effects remain blocked", job.ID)
		}
	}
	for _, b := range []FileBinding{r.Stdout, r.Stderr} {
		if _, err = workflowBound(b, 16<<20); err != nil {
			return r, err
		}
	}
	if r.Answer != nil {
		if _, err = workflowBound(*r.Answer, 1<<20); err != nil {
			return r, err
		}
	}
	return r, nil
}

func workerLaunchAttempt(d Dependencies, s workflowState, in workerAttemptIntent, dir string) error {
	in.RunID, in.Generation = s.Result.RunID, s.Result.Worker.Generation
	in.JobID = filepath.Base(dir)
	b, err := workflowKeep(d, filepath.Join(dir, "intent.json"), in)
	if err != nil {
		return err
	}
	job := &WorkerJob{ID: in.JobID, Kind: in.Kind, Directory: dir, IntentHash: b.SHA256}
	if err = workerControlDispatch(d, s, func(current *workflowState) error {
		if current.Result.Worker.Job != nil {
			return fmt.Errorf("another worker attempt is reserved")
		}
		current.Result.Worker.Job = job
		current.Result.Worker.Phase = "running"
		return nil
	}); err != nil {
		return err
	}
	// Reservation precedes launch. A crash in this gap is deliberately unknown,
	// not permission to launch the same command twice.
	cmd := exec.Command(s.Request.Runtime.PlyExecutable.Path, "workflow", "worker", "_attempt", s.Result.RunID, "--workspace", s.Request.WorkspaceRoot, "--job", job.ID)
	cmd.Dir, cmd.Stdout, cmd.Stderr = s.Observed.Target.WorktreeLocator, os.Stdout, os.Stderr
	cmd.ExtraFiles = []*os.File{d.WorkerOwner.lock}
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("attempt runner ended without confirmed completion: %w", err)
	}
	return nil
}

func workerGit(d Dependencies, s workflowState, args ...string) (string, error) {
	if err := workerLiveOwner(d, s); err != nil {
		return "", err
	}
	argv := append([]string{"-c", "core.hooksPath=/dev/null", "-c", "core.fsmonitor=false", "-c", "commit.gpgSign=false", "-c", "user.name=Ply worker", "-c", "user.email=ply-worker@localhost"}, args...)
	c := exec.Command("git", argv...)
	c.Dir = s.Observed.Target.WorktreeLocator
	c.Env = workerEnvironment()
	// Git is a trusted native effect. It retains ownership if the controller
	// dies between starting Git and recording the resulting commit.
	c.ExtraFiles = []*os.File{d.WorkerOwner.lock}
	b, e := c.CombinedOutput()
	if e != nil {
		return "", fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), b, e)
	}
	return strings.TrimSpace(string(b)), nil
}
