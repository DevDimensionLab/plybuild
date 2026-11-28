package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/agentlaunch"
	"github.com/devdimensionlab/plybuild/internal/publication"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type fakeLaunch struct {
	calls  []agentlaunch.Request
	result agentlaunch.Result
	err    error
}

func (f *fakeLaunch) start(_ context.Context, in agentlaunch.Request) (agentlaunch.Result, error) {
	f.calls = append(f.calls, in)
	return f.result, f.err
}

func runShare(t *testing.T, cwd string, launch *fakeLaunch, args ...string) (publication.Share, error) {
	t.Helper()
	t.Chdir(cwd)
	t.Setenv("HERDR_ENV", "1")
	command := newPublicationCommandWith(workspace.SystemDependencies(), launch.start)
	command.SilenceUsage, command.SilenceErrors = true, true
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&bytes.Buffer{})
	command.SetArgs(append(args, "--format", "json"))
	err := command.Execute()
	var share publication.Share
	if stdout.Len() > 0 {
		if decodeErr := json.Unmarshal(stdout.Bytes(), &share); decodeErr != nil {
			t.Fatalf("share json %q: %v", stdout.String(), decodeErr)
		}
	}
	return share, err
}

func TestPublicationShareStartsOneClaudeInTheRunDirectory(t *testing.T) {
	root := newAssignmentWorkspace(t)
	dir := producerDir(t, root, "graf", map[string]string{"index.html": "graf"})
	v := publishJSON(t, root, "--repo", "a-app", "--category", "code-health-history", "--dir", dir, "--entry", "index.html", "--title", "graf")
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	launch := &fakeLaunch{result: agentlaunch.Result{State: "agent_observed", Pane: "p1", Tab: "t1", Agent: "claude", Session: "s1"}}
	args := []string{"share", "--project", "a", "--repo", "a-app", "--category", "code-health-history", "--herdr-workspace", "w7", "--claude-bin", exe, "--herdr-bin", exe, "--actor", "plyapp-server"}
	share, err := runShare(t, root, launch, args...)
	if err != nil || share.Kind != publication.ShareKind || share.State != publication.ShareStatePending || share.Run != v.Run ||
		share.Session == nil || share.Session.Pane != "p1" || share.Session.ID != "s1" || share.Actor != "plyapp-server" {
		t.Fatalf("share %+v %v", share, err)
	}
	if len(launch.calls) != 1 {
		t.Fatalf("launches %d", len(launch.calls))
	}
	in := launch.calls[0]
	wantPrompt := "/ply-publication-share --project a --workspace " + root + " --repo a-app --category code-health-history --run " + v.Run
	if in.Agent != "claude" || in.Executable != exe || in.Cwd != v.Path || in.Workspace != "w7" || in.InitialPrompt != wantPrompt ||
		in.Label != "share a-app/code-health-history" || !slices.Contains(in.Env, "PLY_ACTOR=claude:publication-share") {
		t.Fatalf("launch %+v", in)
	}
	if _, err := runShare(t, root, launch, args...); err == nil {
		t.Fatal("second share while pending")
	} else {
		requirePublicationClass(t, err, publication.ClassSharePending)
	}
	result := []string{"share-result", "--project", "a", "--repo", "a-app", "--category", "code-health-history", "--run", v.Run, "--url", "https://claude.ai/artifact/Abc_1-2"}
	shared, err := runShare(t, root, launch, result...)
	if err != nil || shared.State != publication.ShareStateShared || *shared.URL != "https://claude.ai/artifact/Abc_1-2" {
		t.Fatalf("result %+v %v", shared, err)
	}
	if again, err := runShare(t, root, launch, args...); err != nil || again.State != publication.ShareStateShared || len(launch.calls) != 1 {
		t.Fatalf("same run again %+v %v %d", again, err, len(launch.calls))
	}
	list := listJSON(t, root, "--project", "a")
	if len(list.Items) != 1 || list.Items[0].Share == nil || *list.Items[0].Share.SharedRun != v.Run {
		t.Fatalf("list %+v", list.Items)
	}
}

func TestPublicationShareRecordsTheLaunchOutcome(t *testing.T) {
	root := newAssignmentWorkspace(t)
	dir := producerDir(t, root, "graf", map[string]string{"index.html": "graf"})
	publishJSON(t, root, "--category", "code-health-history", "--dir", dir, "--entry", "index.html", "--title", "graf")
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"share", "--project", "a", "--category", "code-health-history", "--herdr-workspace", "w7", "--claude-bin", exe, "--herdr-bin", exe}
	failed, err := runShare(t, root, &fakeLaunch{result: agentlaunch.Result{State: "not_started"}, err: errors.New("launch requires a Herdr-managed caller")}, args...)
	requirePublicationClass(t, err, publication.ClassShareLaunchFailed)
	if failed.State != publication.ShareStateFailed || failed.Error == nil || *failed.Error != "launch requires a Herdr-managed caller" || failed.Session != nil {
		t.Fatalf("failed %+v", failed)
	}
	unknown, err := runShare(t, root, &fakeLaunch{result: agentlaunch.Result{State: "submission_unknown", Pane: "p2", Tab: "t2", EffectUnknown: true}, err: errors.New("timeout")}, args...)
	requirePublicationClass(t, err, publication.ClassShareLaunchFailed)
	if unknown.State != publication.ShareStateUnknown || unknown.Session == nil || unknown.Session.Pane != "p2" {
		t.Fatalf("unknown %+v", unknown)
	}
	if _, err := runShare(t, root, &fakeLaunch{}, args...); err == nil {
		t.Fatal("share while unknown")
	} else {
		requirePublicationClass(t, err, publication.ClassSharePending)
	}
}

func TestPublicationShareChecksBeforeAnyEffect(t *testing.T) {
	root := newAssignmentWorkspace(t)
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	launch := &fakeLaunch{}
	_, err = runShare(t, root, launch, "share", "--project", "a", "--category", "code-health-history", "--herdr-workspace", "w7", "--claude-bin", exe, "--herdr-bin", exe)
	requirePublicationClass(t, err, publication.ClassNotFound)
	_, err = runShare(t, root, launch, "share", "--project", "a", "--repo", "b-app", "--category", "x", "--herdr-workspace", "w7", "--claude-bin", exe, "--herdr-bin", exe)
	requirePublicationClass(t, err, publication.ClassInvalidArguments)
	_, err = runShare(t, root, launch, "share-result", "--project", "a", "--category", "x", "--run", "20261128T000000000Z-00000000")
	requirePublicationClass(t, err, publication.ClassInvalidArguments) // neither --url nor --error
	dir := producerDir(t, root, "graf", map[string]string{"index.html": "graf"})
	publishJSON(t, root, "--category", "x", "--dir", dir, "--entry", "index.html", "--title", "graf")
	t.Chdir(root)
	command := newPublicationCommandWith(workspace.SystemDependencies(), launch.start)
	command.SilenceUsage, command.SilenceErrors = true, true
	command.SetOut(&bytes.Buffer{})
	command.SetArgs([]string{"share", "--project", "a", "--category", "x", "--herdr-workspace", "w7", "--claude-bin", exe, "--herdr-bin", exe})
	t.Setenv("HERDR_ENV", "")
	requirePublicationClass(t, command.Execute(), publication.ClassInvalidArguments)
	if v, err := runPublication(t, root, "show", "--project", "a", "--category", "x", "--format", "json"); err != nil || !strings.Contains(v, `"share": null`) {
		t.Fatalf("a refused caller reserved the share: %s %v", v, err)
	}
	if len(launch.calls) != 0 {
		t.Fatal("launched")
	}
}
