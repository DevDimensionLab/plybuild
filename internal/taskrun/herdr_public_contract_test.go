package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

// This synthetic public Herdr fixture has no private startup capability or
// operations. It exercises Ply consumers, not a live provider or atomic writer.
func publicHerdrFixture(t *testing.T, f *workflowFixture) {
	t.Helper()
	workflowTestModel(t, *f, map[string]any{"public_contract_only": true})
}

func assertPublicHerdrCalls(t *testing.T, f workflowFixture, prompts int) {
	t.Helper()
	b, err := os.ReadFile(f.Calls)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, line := range bytes.Split(bytes.TrimSpace(b), []byte("\n")) {
		var call struct{ Argv []string }
		if err = json.Unmarshal(line, &call); err != nil {
			t.Fatal(err)
		}
		argv := strings.Join(call.Argv, " ")
		if strings.Contains(argv, "agent startup") || strings.Contains(argv, "tui.status_line") {
			t.Fatalf("private startup requirement: %s", argv)
		}
		if len(call.Argv) >= 2 && call.Argv[0] == "agent" && call.Argv[1] == "prompt" {
			count++
		}
	}
	if count != prompts {
		t.Fatalf("public prompt count = %d, want %d", count, prompts)
	}
}

func TestPublicHerdrOwnerStartup(t *testing.T) {
	f := workflowTestFixture(t)
	publicHerdrFixture(t, &f)
	o, err := workflowReadinessStart(t, f)
	if err != nil {
		t.Fatal(err)
	}
	if o.Transport.AgentSessionID != "fixture-session" || o.Transport.State != "working" {
		t.Fatalf("owner not bound and prompted: %+v", o.Transport)
	}
	assertPublicHerdrCalls(t, f, 1)
	workflowAssertNoReplay(t, f, o)
}

func TestPublicHerdrIndependentReviewerStartup(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	publicHerdrFixture(t, &f.workflowFixture)
	reviewRoundCandidate(t, f, "fixed")
	reviewRoundIn(t, "w-fixture:root")
	r, err := WorkflowDeliveryRequestReview(reviewRoundAsCoordinator(f), f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root"})
	if err != nil {
		t.Fatal(err)
	}
	if reviewRoundLatest(t, r).ReviewerStart != "started" {
		t.Fatal("independent reviewer did not start")
	}
	if _, err = WorkflowDeliveryRequestReview(reviewRoundAsCoordinator(f), f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root"}); err != nil {
		t.Fatal(err)
	}
	assertPublicHerdrCalls(t, f.workflowFixture, 2)
}

func publicHerdrResumeFixture(t *testing.T) ownerResumeFixture {
	t.Helper()
	f := newOwnerResumeFixture(t, "codex", false)
	// Keep the unrelated preview timestamp/process inventory stable so this
	// regression reaches the real public-Herdr resume consumer.
	f.D.Now = func() time.Time { return time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC) }
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: os.Getpid(), ParentPID: os.Getppid()}}, nil
	}
	publicHerdrFixture(t, &f.workflowFixture)
	return f
}

func TestPublicHerdrOwnerNewTerminalResume(t *testing.T) {
	f := publicHerdrResumeFixture(t)
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil || p.State != "ready" {
		t.Fatalf("resume not ready: %+v %v", p, err)
	}
	ownerResumePreserve(t, f, p)
	o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if o.Transport.AgentSessionID != f.Owner.Transport.AgentSessionID || o.Transport.PaneID == f.Owner.Transport.PaneID || o.OwnerResume.Step != "prompt_sent" {
		t.Fatalf("resume lost owner binding: %+v", o)
	}
	if _, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation); err != nil {
		t.Fatal(err)
	}
	assertPublicHerdrCalls(t, f.workflowFixture, 2)
}

func TestPublicHerdrOwnerResumeRejectsDifferentSession(t *testing.T) {
	f := publicHerdrResumeFixture(t)
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil || p.State != "ready" {
		t.Fatalf("resume not ready: %+v %v", p, err)
	}
	ownerResumePreserve(t, f, p)
	workflowTestModel(t, f.workflowFixture, map[string]any{"agent_session_id": "foreign-session"})
	if _, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation); err == nil || !strings.Contains(err.Error(), "same bound session") {
		t.Fatalf("resume did not reject the changed session: %v", err)
	}
	assertPublicHerdrCalls(t, f.workflowFixture, 1)
}
