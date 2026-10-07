package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/delivery"
	"github.com/devdimensionlab/plybuild/internal/taskexecute"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type recordingDeliveryService struct {
	calls   []string
	args    []string
	receipt delivery.Receipt
	err     error
}

func (s *recordingDeliveryService) record(verb string, args ...string) (delivery.Receipt, error) {
	s.calls = append(s.calls, verb)
	s.args = args
	return s.receipt, s.err
}
func (s *recordingDeliveryService) Register(cwd, file string) (delivery.Receipt, error) {
	return s.record("register", cwd, file)
}
func (s *recordingDeliveryService) Check(cwd, id string) (delivery.Receipt, error) {
	return s.record("check", cwd, id)
}
func (s *recordingDeliveryService) Execute(cwd, id string) (delivery.Receipt, error) {
	return s.record("execute", cwd, id)
}
func (s *recordingDeliveryService) Read(cwd, id string) (delivery.Receipt, error) {
	return s.record("show", cwd, id)
}
func (s *recordingDeliveryService) Metadata(cwd, id, file string) (delivery.Receipt, error) {
	return s.record("metadata", cwd, id, file)
}
func (s *recordingDeliveryService) List(cwd, task, epic string) (delivery.ListResult, error) {
	s.record("list", cwd, task, epic)
	return delivery.ListResult{Kind: "ply.delivery.list", SchemaVersion: 1, Deliveries: []delivery.Receipt{s.receipt}}, s.err
}

func TestDeliveryCLIRejectsAmbiguousOrMissingInputBeforeEffects(t *testing.T) {
	for _, args := range [][]string{
		{"register"}, {"register", "unexpected", "--file", "/private/draft.json"},
		{"register", "--file", "relative.json"}, {"check"}, {"show", "one", "two"},
		{"execute", "one", "--format", "yaml"}, {"execute", "one", "--force"},
		{"metadata", "one"}, {"metadata", "one", "--file", "relative.json"},
		{"list", "unexpected"}, {"list", "--task", ""}, {"list", "--epic", "../epic"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			s := &recordingDeliveryService{}
			c := newWorkflowDeliveryCommand(s)
			c.SetOut(&bytes.Buffer{})
			c.SetErr(&bytes.Buffer{})
			c.SetArgs(args)
			if err := c.Execute(); err == nil || len(s.calls) != 0 {
				t.Fatalf("err=%v calls=%v", err, s.calls)
			}
		})
	}
}

func TestDeliveryCLIReadCommandsNeverExecute(t *testing.T) {
	for _, args := range [][]string{{"show", "dlv_fixture"}, {"check", "dlv_fixture"}, {"list", "--task", "task", "--epic", "epic"}} {
		t.Run(args[0], func(t *testing.T) {
			s := &recordingDeliveryService{receipt: delivery.Receipt{Kind: "ply.delivery.receipt", SchemaVersion: 1, ID: "dlv_fixture", State: "waiting", NextAction: "Record actual candidate-bound human QA."}}
			c := newWorkflowDeliveryCommand(s)
			var out bytes.Buffer
			c.SetOut(&out)
			c.SetErr(&bytes.Buffer{})
			c.SetArgs(append(args, "--format", "json"))
			if err := c.Execute(); err != nil {
				t.Fatal(err)
			}
			if len(s.calls) != 1 || s.calls[0] != args[0] {
				t.Fatalf("calls=%v", s.calls)
			}
			if !json.Valid(out.Bytes()) {
				t.Fatalf("invalid JSON: %s", out.String())
			}
			if args[0] == "list" && (s.args[1] != "task" || s.args[2] != "epic") {
				t.Fatalf("lost exact filters: %v", s.args)
			}
		})
	}
}

func TestDeliveryCLIEmitsObservedPRBeforeReturningMetadataFailure(t *testing.T) {
	for _, format := range []string{"text", "json"} {
		t.Run(format, func(t *testing.T) {
			failure := errors.New("account cannot be assigned")
			s := &recordingDeliveryService{err: failure, receipt: delivery.Receipt{Kind: "ply.delivery.receipt", SchemaVersion: 1, ID: "dlv_fixture", State: "failed", NextAction: "Correct metadata and execute this same Delivery.", PR: &delivery.PRReceipt{URL: "https://github.com/example/product/pull/17", HeadOID: "exact-candidate", MetadataApplied: false}, NativeClosed: false}}
			c := newWorkflowDeliveryCommand(s)
			var out bytes.Buffer
			c.SetOut(&out)
			c.SetErr(&bytes.Buffer{})
			c.SetArgs([]string{"execute", "dlv_fixture", "--format", format})
			if err := c.Execute(); !errors.Is(err, failure) {
				t.Fatalf("err=%v", err)
			}
			if !strings.Contains(out.String(), s.receipt.PR.URL) || !strings.Contains(out.String(), "false") {
				t.Fatalf("partial receipt lost: %s", out.String())
			}
			if len(s.calls) != 1 || s.calls[0] != "execute" {
				t.Fatalf("calls=%v", s.calls)
			}
		})
	}
}

func TestExecutePreviewShowsExplicitModeAndDistinctEffects(t *testing.T) {
	for _, mode := range []string{workspace.DeliveryPullRequest, workspace.DeliveryLocalEpic, workspace.DeliveryLocalBranch} {
		t.Run(mode, func(t *testing.T) {
			a := &workspace.DeliveryAgreement{SchemaVersion: 1, Mode: mode, ProjectID: "project", RepoID: "repo", EpicID: "epic", SourceRef: "refs/heads/task", TargetRef: "refs/heads/main"}
			if mode == workspace.DeliveryPullRequest {
				a.GitHubRepository = "example/product"
				a.Remote = "origin"
			} else {
				a.TargetWorktree = "/workspace/product"
			}
			r := taskexecute.Result{Goal: &workspace.TaskGoalExecutePreview{Delivery: a}}
			c := newWorkflowDeliveryCommand(nil)
			var out bytes.Buffer
			c.SetOut(&out)
			if err := writeExecuteResult(c, "text", r); err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{mode, a.SourceRef, a.TargetRef, "exact candidate human pass", "actual runtime authority"} {
				if !strings.Contains(out.String(), want) {
					t.Fatalf("preview omits %s: %s", want, out.String())
				}
			}
			if mode == workspace.DeliveryPullRequest && !strings.Contains(out.String(), "stop before merge") {
				t.Fatal(out.String())
			}
		})
	}
}
