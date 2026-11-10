package workspace

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestPriorReaderPreservesPendingRefreshAuthority(t *testing.T) {
	binary := os.Getenv("PLY_PRIOR_READER_BINARY")
	if binary == "" {
		t.Skip("opt-in compatibility fixture: set PLY_PRIOR_READER_BINARY to the actual prior release executable")
	}
	for _, format := range []string{"inline6", "inline", "9", "10"} {
		t.Run(format, func(t *testing.T) {
			f, pending := pendingRefreshFixture(t, format)
			root := f.root
			epicPath := f.epicPath
			epicBefore := runLocalGit(t, epicPath, "rev-parse", "HEAD")
			stored := storedIntegrationEvidence(t, root)
			if stored.format != strings.TrimSuffix(format, "6") {
				t.Fatalf("fixture format = %s", stored.format)
			}
			prior := priorReader{binary: binary, root: root, home: filepath.Join(root, "prior-home")}
			if err := os.Mkdir(prior.home, 0o700); err != nil {
				t.Fatal(err)
			}
			base := []string{"workspace", "task", "integrate", "second", "--expected-result-oid", pending.Plan.Task.ResultOID, "--expected-parent-oid", pending.Plan.Epic.ExpectedParentOID}

			out, stderr, exit := prior.run(t, append(append([]string{}, base...), "--check", "--format", "json")...)
			var check struct {
				Classification string                `json:"classification"`
				NextAction     IntegrationNextAction `json:"next_action"`
			}
			if exit != 0 || json.Unmarshal([]byte(out), &check) != nil || check.Classification != "blocked" || check.NextAction.Reason != historicalPendingReason {
				t.Fatalf("prior check did not stop on pending refresh: exit=%d out=%s err=%s", exit, out, stderr)
			}
			out, stderr, exit = prior.run(t, append(append([]string{}, base...), "--apply", "--confirm", pending.PlanSHA256, "--format", "json")...)
			if exit == 0 || !strings.Contains(out+stderr, "historical integration authority has no recorded outcome") {
				t.Fatalf("prior apply did not fail closed: exit=%d out=%s err=%s", exit, out, stderr)
			}
			if out, stderr, exit = prior.run(t, "workspace", "task", "show", "second", "--format", "json"); exit != 0 || strings.Contains(out, `"classification":"ready"`) {
				t.Fatalf("prior show: exit=%d out=%s err=%s", exit, out, stderr)
			}
			if out, stderr, exit = prior.run(t, "workspace", "task", "create", "later", "--title", "Later", "--description", "Prior writer", "--epic", "epic", "--project", "ply", "--repo", "ply"); exit != 0 {
				t.Fatalf("prior writer: exit=%d out=%s err=%s", exit, out, stderr)
			}
			if got := runLocalGit(t, epicPath, "rev-parse", "HEAD"); got != epicBefore {
				t.Fatalf("prior reader moved the Epic parent to %s", got)
			}
			after := storedIntegrationEvidence(t, root)
			if after.format != stored.format || !bytes.Equal(after.raw, stored.raw) {
				t.Fatalf("prior writer changed integration evidence:\n%s\n---\n%s", stored.raw, after.raw)
			}
			r := mustQueueRegistry(t, f.workItemJourneyFixture)
			if task, _ := findTask(r, "later"); task == nil {
				t.Fatal("prior writer did not publish")
			}
			resume := f.input(true)
			resume.Confirmation = pending.PlanSHA256
			done, err := ApplyTaskIntegration(f.dependencies, resume)
			if err != nil || done.Readback.Classification != "exact_effect" || runLocalGit(t, epicPath, "rev-parse", "HEAD") != f.c {
				t.Fatalf("current reader did not complete the preserved authority: %s %v", done.Readback.Classification, err)
			}
		})
	}
}

type priorReader struct{ binary, root, home string }

func (p priorReader) run(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	command := exec.Command(p.binary, args...)
	command.Dir = p.root
	env := []string{}
	for _, entry := range os.Environ() {
		if key, _, _ := strings.Cut(entry, "="); key != "HOME" && !strings.HasPrefix(key, "GIT_") {
			env = append(env, entry)
		}
	}
	command.Env = append(env, "HOME="+p.home, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+os.DevNull, "LC_ALL=C")
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	err := command.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		return stdout.String(), stderr.String(), exitErr.ExitCode()
	}
	if err != nil {
		t.Fatal(err)
	}
	return stdout.String(), stderr.String(), 0
}

type storedEvidence struct {
	format string
	raw    []byte
}

func storedIntegrationEvidence(t *testing.T, root string) storedEvidence {
	t.Helper()
	raw, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	node, err := parseRegistryYAML(raw)
	if err != nil {
		t.Fatal(err)
	}
	out := storedEvidence{format: "inline"}
	registry := node
	if inner := registryNodeField(node, "registry"); inner != nil {
		out.format = registryNodeField(node, "format_version").Value
		registry = inner
	}
	for _, name := range evidenceCollections(false) {
		collection := registryNodeField(registry, name)
		encoded, err := marshalRegistryNode(collection)
		if err != nil {
			t.Fatal(err)
		}
		out.raw = append(out.raw, encoded...)
		for _, item := range collection.Content {
			if item.Kind != yaml.ScalarNode {
				continue
			}
			object, err := evidenceStorage(newSystemWorkItemStore()).Read(root, integrationEvidenceKind, item.Value)
			if err != nil {
				t.Fatal(err)
			}
			out.raw = append(out.raw, object...)
		}
	}
	return out
}

func setIntegrationStoreFormat(t *testing.T, root, format string) {
	t.Helper()
	s := newSystemWorkItemStore()
	r, err := s.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	var raw []byte
	switch format {
	case "10":
		return
	case "inline6":
		upgradeEpicCreateFixture(t, SystemDependencies(), root)
		if r, err := s.Snapshot(root); err != nil || r.FormatVersion != 6 {
			t.Fatalf("inline fixture is not logical format 6: %v", err)
		}
		return
	case "9":
		old := r
		old.publicationEvidence = false
		raw, err = s.encodeStoredRegistry(root, &r, old)
	default:
		raw, err = encodeWorkItemRegistry(r)
	}
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(workItemsPath(root), raw, 0o600); err != nil {
		t.Fatal(err)
	}
}

func pendingRefreshFixture(t *testing.T, format string) (refreshFixture, IntegrationAuthority) {
	t.Helper()
	f := newRefreshFixture(t)
	setIntegrationStoreFormat(t, f.root, format)
	var pending IntegrationAuthority
	err := f.dependencies.ProjectLocks.WithSnapshotLock(f.root, func(projects ProjectSnapshot) error {
		return f.dependencies.WorkItems.WithLock(f.root, func(session WorkItemStoreSession) error {
			r, err := session.Snapshot()
			if err != nil {
				return err
			}
			plan, digest, _, err := buildIntegrationPlan(f.dependencies, f.root, projects, r, f.input(true))
			if err != nil {
				return err
			}
			if plan.Readiness != "ready" || plan.Kind != refreshIntegrationPlanKind {
				t.Fatalf("refresh plan = %s %s %v", plan.Kind, plan.Readiness, plan.Reasons)
			}
			authorityID, _ := f.dependencies.TaskLifecycleIDs.NewIntegrationAuthorityID()
			intentID, _ := f.dependencies.TaskLifecycleIDs.NewIntegrationIntentID()
			pending = IntegrationAuthority{ID: authorityID, CreatedAtUTC: f.dependencies.WorkClock.Now().UTC().Format(time.RFC3339Nano), PlanSHA256: digest, Plan: plan, TaskID: "second", AllowedEffect: IntegrationAllowedEffect{Kind: "local_ff_only", ParentRef: plan.Epic.ParentRef, ExpectedParentOID: plan.Epic.ExpectedParentOID, ResultOID: plan.Task.ResultOID, MaxOccurrences: 1}}
			intent := intentFromPlan(intentID, authorityID, plan, digest)
			intent.IntentSHA256 = integrationIntentDigest(intent)
			r.IntegrationAuthorities = append(r.IntegrationAuthorities, pending)
			r.IntegrationIntents = append(r.IntegrationIntents, intent)
			sortWorkRegistry(&r)
			return session.Publish(r)
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	if a := findAuthority(r, pending.ID); a == nil || a.historical != nil {
		t.Fatal("pending refresh authority is not an active record")
	}
	return f, pending
}
