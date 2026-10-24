package taskrun

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

// The stand-in decodes each real agent.start -c as TOML and only exposes
// initialized prompt evidence when the selected footer includes thread identity.
func TestManagedThreadFooterOwner(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "initialized", true: "missing_marker"}[missing], func(t *testing.T) {
			f := workflowTestFixture(t)
			f.D.HerdrTimeout = 4 * time.Second
			workflowTestModel(t, f, map[string]any{"safe_require_thread_profile": true, "safe_omit_thread_marker": missing, "bootstrap_mode": "normal", "interactive_ready": true, "bootstrap_run_root": workflowPaths(f.R, 0).RunRoot, "bootstrap_expected_prompt": workflowExpectedBootstrapPrompt})
			o, err := workflowReadinessStart(t, f)
			if missing {
				if err == nil || len(safeStartupKinds(t, f)) != 0 {
					t.Fatalf("missing thread marker gained input: %v", err)
				}
			} else if err != nil || strings.Join(safeStartupKinds(t, f), ",") != "bootstrap,task" {
				t.Fatalf("initialized owner did not bootstrap and submit: %v", err)
			}
			if workflowUnguardedCalls(t, f) != 0 {
				t.Fatal("unguarded startup")
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}

func TestManagedThreadFooterReviewer(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	f.D.HerdrTimeout = 6 * time.Second
	reviewRoundCandidate(t, f, "fixed")
	reviewRoundIn(t, "w-fixture:root")
	workflowTestModel(t, f.workflowFixture, map[string]any{"safe_require_thread_profile": true, "review_sessionless": true})
	r, err := WorkflowDeliveryRequestReview(reviewRoundAsCoordinator(f), f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root"})
	if err != nil || reviewRoundLatest(t, r).ReviewerStart != "started" || strings.Join(safeStartupKinds(t, f.workflowFixture), ",") != "task,bootstrap,task" {
		t.Fatalf("initialized reviewer did not bootstrap and receive brief: %v", err)
	}
	if workflowUnguardedCalls(t, f.workflowFixture) != 0 {
		t.Fatal("unguarded review")
	}
}

func TestManagedThreadFooterResume(t *testing.T) {
	f := newOwnerResumeFixture(t, "codex", false)
	f.D.Now = func() time.Time { return time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC) }
	workflowTestModel(t, f.workflowFixture, map[string]any{"safe_require_thread_profile": true})
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil || p.State != "ready" {
		t.Fatalf("resume preview: %+v %v", p, err)
	}
	ownerResumePreserve(t, f, p)
	before := len(safeStartupKinds(t, f.workflowFixture))
	_, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
	if err != nil {
		t.Fatalf("initialized resume did not continue: %v", err)
	}
	if len(safeStartupKinds(t, f.workflowFixture)) != before+1 || workflowUnguardedCalls(t, f.workflowFixture) != 0 {
		t.Fatal("resume continuation not guarded")
	}
	calls := workflowProviderCalls(t, f.workflowFixture, "agent start")
	argv := calls[len(calls)-1]
	child := argv[0:0]
	for i, arg := range argv {
		if arg == "--" {
			child = argv[i+1:]
			break
		}
	}
	if len(child) < 2 || child[0] != "resume" || child[1] != f.Owner.Transport.AgentSessionID {
		t.Fatalf("resume ordering/session changed: %v", child)
	}
}

// Crosses the built native CLI boundary with a modeled producer; no provider runs.
func TestManagedThreadFooterCLI(t *testing.T) {
	binary := automaticCLIBinary(t)
	f := newDeliveryFixture(t, "codex")
	f.R.Runtime.PlyExecutable = Executable{Path: binary, SHA256: hashFileTest(t, binary)}
	f.R.Runtime.ConfigProfile = ptr("selected-profile")
	f.File = writeAny(t, f.R.WorkspaceRoot, "thread-footer-request.json", f.R)
	workflowTestModel(t, f.workflowFixture, map[string]any{"safe_require_thread_profile": true})
	_, out, err := automaticCLI(binary, f, "run", "start", "--file", f.File, "--check")
	var p WorkflowPreview
	if err != nil || json.Unmarshal(out, &p) != nil || p.Confirmation == nil {
		t.Fatalf("CLI preview: %v %s", err, out)
	}
	o := automaticCLIOK(t, binary, f, "run", "start", "--file", f.File, "--apply", "--confirm", *p.Confirmation)
	if o.Transport.State != "working" || strings.Join(safeStartupKinds(t, f.workflowFixture), ",") != "task" || workflowUnguardedCalls(t, f.workflowFixture) != 0 {
		t.Fatalf("CLI startup did not reach guarded task: %+v", o.Transport)
	}
	starts := workflowProviderCalls(t, f.workflowFixture, "agent start")
	if len(starts) != 1 {
		t.Fatalf("CLI start count: %d", len(starts))
	}
	argv := starts[0]
	for _, pair := range [][2]string{
		{"--model", f.R.Runtime.Model}, {"-C", f.Prepared.Preparation.Preparation.Plan.WorktreePath},
		{"--profile", "selected-profile"}, {"-a", "on-request"},
		{"-c", `approvals_reviewer="auto_review"`},
		{"-c", "default_permissions=" + workflowJSON(f.R.Runtime.PermissionBinding.ProfileID)},
		{"-c", `model_reasoning_effort="medium"`},
		{"-c", `tui.status_line=["thread-id","context-remaining"]`},
	} {
		found := false
		for i := 0; i+1 < len(argv); i++ {
			if argv[i] == pair[0] && argv[i+1] == pair[1] {
				found = true
			}
		}
		if !found {
			t.Fatalf("CLI lost launch option %q in %q", pair, argv)
		}
	}
	before, _ := os.ReadFile(f.Calls)
	automaticCLIOK(t, binary, f, "run", "show", o.RunID)
	after, _ := os.ReadFile(f.Calls)
	if string(before) != string(after) {
		t.Fatal("CLI historical readback replayed transport")
	}
}
