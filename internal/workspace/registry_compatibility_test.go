package workspace

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestRegistryDiagnosticsUseDeclaredFormatAndValidation(t *testing.T) {
	valid := "format_version: 1\nepics: []\ntasks: []\nworktree_operations: []\n"
	for name, tc := range map[string]struct{ input, code string }{
		"future":             {strings.Replace(valid, "version: 1", "version: 99", 1) + "future_field: true\n", "registry_unsupported_format"},
		"authorization":      {valid + "acceptance_override: true\n", "registry_unknown_field"},
		"syntax":             {valid + "broken: [", "registry_malformed_yaml"},
		"missing version":    {strings.TrimPrefix(valid, "format_version: 1\n"), "registry_invalid_content"},
		"wrong version type": {strings.Replace(valid, "version: 1", "version: '99'", 1), "registry_invalid_content"},
		"duplicate":          {valid + "format_version: 99\n", "registry_invalid_content"},
		"multiple documents": {valid + "---\nformat_version: 99\n", "registry_malformed_yaml"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := decodeWorkItemRegistry([]byte(tc.input))
			d := RegistryError(err)
			if d == nil || d.Code != tc.code || d.NextAction == "" {
				t.Fatalf("diagnostic: %+v %v", d, err)
			}
		})
	}
	f := newWorkItemJourneyFixture(t)
	r, _ := f.dependencies.WorkItems.Snapshot(filepath.Dir(f.wrapper))
	b, _ := encodeWorkItemRegistry(r)
	b = bytes.Replace(b, []byte("parent_epic_id: epic"), []byte("parent_epic_id: missing"), 1)
	_, err := decodeWorkItemRegistry(b)
	if d := RegistryError(err); d == nil || d.Code != "registry_invalid_content" {
		t.Fatalf("invalid reference: %v", err)
	}
}

func TestRegistryRoundtripPreservesAutomaticAuthorizationAndHistory(t *testing.T) {
	f, in := newAutomaticIntegrationFixture(t, false)
	check, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	in.Apply, in.Confirmation = true, integrationDigest(t, check.Readback)
	if _, err = ApplyTaskIntegration(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(workItemsPath(f.root))
	if err != nil {
		t.Fatal(err)
	}
	registry, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	logicalBefore, err := encodeWorkItemRegistry(registry)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"acceptance:", "responsible_actor:", "integration_results:"} {
		if !bytes.Contains(logicalBefore, []byte(field)) {
			t.Fatalf("fixture missing %s", field)
		}
	}
	err = f.dependencies.WorkItems.WithLock(f.root, func(s WorkItemStoreSession) error {
		r, err := s.Snapshot()
		if err != nil {
			return err
		}
		return s.Publish(r)
	})
	if err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(workItemsPath(f.root))
	if !bytes.Equal(before, after) {
		t.Fatal("ordinary roundtrip changed authority or history")
	}
	bad := bytes.Replace(logicalBefore, []byte("responsible_actor:"), []byte("unknown_authorizer:"), 1)
	_, err = decodeWorkItemRegistry(bad)
	if d := RegistryError(err); d == nil || d.Code != "registry_unknown_field" {
		t.Fatalf("unknown authority not rejected: %v", err)
	}
}

func TestRegistryFormatChangeRequiresNativePermit(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	queueSetFixture(t, f, "format-fixture", []QueueEntry{})
	root := filepath.Dir(f.wrapper)
	before, _ := os.ReadFile(workItemsPath(root))
	err := f.dependencies.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
		r, err := s.Snapshot()
		if err != nil {
			return err
		}
		r.FormatVersion = 5
		return s.Publish(r)
	})
	if d := RegistryError(err); d == nil || d.Code != "registry_upgrade_required" {
		t.Fatalf("missing permit: %v", err)
	}
	after, _ := os.ReadFile(workItemsPath(root))
	if !bytes.Equal(before, after) {
		t.Fatal("rejection changed registry")
	}
}

// The maintainer acceptance command supplies the actual delivered package. An
// ordinary portable unit run does not build or invent a release.
func deliveredRegistryReader(t *testing.T) (string, string) {
	t.Helper()
	directory, pin := os.Getenv("PLY_REGISTRY_TEST_RELEASE"), os.Getenv("PLY_REGISTRY_TEST_MANIFEST")
	if directory == "" || pin == "" {
		t.Skip("run release acceptance with PLY_REGISTRY_TEST_RELEASE and PLY_REGISTRY_TEST_MANIFEST")
	}
	return directory, pin
}

func TestRegistryDeliveredTransitionFaultsAndRetry(t *testing.T) {
	directory, pin := deliveredRegistryReader(t)
	for _, stage := range []string{"temp-open", "write", "file-sync", "replace", "directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			queueSetFixture(t, f, "format-fixture", []QueueEntry{})
			root := filepath.Dir(f.wrapper)
			before, _ := os.ReadFile(workItemsPath(root))
			p, err := UpgradeRegistry(f.dependencies, directory, pin, "", false)
			if err != nil {
				t.Fatal(err)
			}
			f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(s string) error {
				if s == stage {
					return errors.New("interrupted")
				}
				return nil
			}}}
			if _, err = UpgradeRegistry(f.dependencies, directory, pin, p.Confirmation, true); err == nil {
				t.Fatal("fault reported success")
			}
			after, _ := os.ReadFile(workItemsPath(root))
			if stage != "directory-sync" && !bytes.Equal(before, after) {
				t.Fatal("pre-commit failure changed registry")
			}
			f.dependencies.WorkItems = newSystemWorkItemStore()
			got, err := UpgradeRegistry(f.dependencies, directory, pin, p.Confirmation, true)
			if err != nil || got.State != "upgraded" && got.State != "already_upgraded" {
				t.Fatalf("recovery: %+v %v", got, err)
			}
			backup, err := os.ReadFile(got.Backup)
			if err != nil || !bytes.Equal(before, backup) {
				t.Fatal("missing exact backup")
			}
			r, err := f.dependencies.WorkItems.Snapshot(root)
			if err != nil {
				t.Fatal(err)
			}
			r.FormatVersion, r.FormatBoundary = 4, nil
			restored, err := encodeWorkItemRegistry(r)
			if err != nil || !bytes.Equal(restored, before) {
				t.Fatal("migration changed records")
			}
			queueSetFixture(t, f, "compatible-format-5-write", []QueueEntry{})
		})
	}
}

func TestRegistryDeliveredTransitionStaleAndConcurrent(t *testing.T) {
	directory, pin := deliveredRegistryReader(t)
	f := newWorkItemJourneyFixture(t)
	queueSetFixture(t, f, "format-fixture", []QueueEntry{})
	p, err := UpgradeRegistry(f.dependencies, directory, pin, "", false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = UpgradeRegistry(f.dependencies, directory, pin, "", true); err == nil {
		t.Fatal("missing confirmation accepted")
	}
	queueSetFixture(t, f, "competing-write", []QueueEntry{})
	if _, err = UpgradeRegistry(f.dependencies, directory, pin, p.Confirmation, true); err == nil {
		t.Fatal("stale check accepted")
	}
	p, err = UpgradeRegistry(f.dependencies, directory, pin, "", false)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	outcomes := make(chan string, 2)
	for n := 0; n < 2; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, e := UpgradeRegistry(f.dependencies, directory, pin, p.Confirmation, true)
			if e != nil {
				outcomes <- e.Error()
			} else {
				outcomes <- got.State
			}
		}()
	}
	wg.Wait()
	close(outcomes)
	counts := map[string]int{}
	for state := range outcomes {
		counts[state]++
	}
	if counts["upgraded"] != 1 || counts["already_upgraded"] != 1 {
		t.Fatalf("concurrent outcomes: %v", counts)
	}
}

func TestRegistryDeliveredReaderMustRemainAvailableAtPublication(t *testing.T) {
	original, pin := deliveredRegistryReader(t)
	for _, stage := range []string{"before-apply", "replace"} {
		t.Run(stage, func(t *testing.T) {
			directory := physicalPath(t, t.TempDir())
			for _, name := range []string{"ply", "ply-install", "manifest.json"} {
				b, err := os.ReadFile(filepath.Join(original, name))
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(directory, name), b, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			f := newWorkItemJourneyFixture(t)
			queueSetFixture(t, f, "format-fixture", []QueueEntry{})
			root := filepath.Dir(f.wrapper)
			before, _ := os.ReadFile(workItemsPath(root))
			p, err := UpgradeRegistry(f.dependencies, directory, pin, "", false)
			if err != nil {
				t.Fatal(err)
			}
			remove := func() error { return os.Remove(filepath.Join(directory, "ply")) }
			if stage == "before-apply" {
				if err := remove(); err != nil {
					t.Fatal(err)
				}
			} else {
				f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(s string) error {
					if s == stage {
						return remove()
					}
					return nil
				}}}
			}
			if _, err := UpgradeRegistry(f.dependencies, directory, pin, p.Confirmation, true); RegistryError(err) == nil || RegistryError(err).Code != "registry_upgrade_required" {
				t.Fatalf("missing delivered reader accepted: %v", err)
			}
			after, _ := os.ReadFile(workItemsPath(root))
			if !bytes.Equal(before, after) {
				t.Fatal("reader disappeared but new format was published")
			}
		})
	}
}

func TestRegistryDeliveredProcessKillRecovery(t *testing.T) {
	directory, pin := deliveredRegistryReader(t)
	if os.Getenv("PLY_REGISTRY_TEST_CHILD") == "1" {
		root, stage := os.Getenv("PLY_REGISTRY_TEST_ROOT"), os.Getenv("PLY_REGISTRY_TEST_STAGE")
		d := SystemDependencies()
		d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: root}
		d.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(s string) error {
			if s == stage {
				_ = os.WriteFile(filepath.Join(root, "checkpoint"), []byte(s), 0o600)
				select {}
			}
			return nil
		}}}
		_, _ = UpgradeRegistry(d, directory, pin, os.Getenv("PLY_REGISTRY_TEST_CONFIRM"), true)
		os.Exit(3)
	}
	for _, stage := range []string{"replace", "directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			queueSetFixture(t, f, "format-fixture", []QueueEntry{})
			root := filepath.Dir(f.wrapper)
			before, _ := os.ReadFile(workItemsPath(root))
			p, err := UpgradeRegistry(f.dependencies, directory, pin, "", false)
			if err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command(os.Args[0], "-test.run=^TestRegistryDeliveredProcessKillRecovery$")
			cmd.Env = append(os.Environ(), "PLY_REGISTRY_TEST_CHILD=1", "PLY_REGISTRY_TEST_ROOT="+root, "PLY_REGISTRY_TEST_STAGE="+stage, "PLY_REGISTRY_TEST_CONFIRM="+p.Confirmation)
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			defer func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }()
			deadline := time.Now().Add(15 * time.Second)
			for {
				if _, err := os.Stat(filepath.Join(root, "checkpoint")); err == nil {
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("child missed registry publication checkpoint")
				}
				time.Sleep(10 * time.Millisecond)
			}
			if err := cmd.Process.Kill(); err != nil {
				t.Fatal(err)
			}
			_ = cmd.Wait()
			after, err := os.ReadFile(workItemsPath(root))
			if err != nil {
				t.Fatal(err)
			}
			if stage == "replace" && !bytes.Equal(before, after) {
				t.Fatal("registry changed before replacement")
			}
			if _, err := decodeWorkItemRegistry(after); err != nil {
				t.Fatalf("partial registry after SIGKILL: %v", err)
			}
			got, err := UpgradeRegistry(f.dependencies, directory, pin, p.Confirmation, true)
			expected := "upgraded"
			if stage == "directory-sync" {
				expected = "already_upgraded"
			}
			if err != nil || got.State != expected {
				t.Fatalf("SIGKILL recovery: %+v %v", got, err)
			}
			backup, err := os.ReadFile(got.Backup)
			if err != nil || !bytes.Equal(before, backup) {
				t.Fatal("lost registry preimage after SIGKILL")
			}
		})
	}
}

func TestRegistryDeliveredHistoryAndOrdinaryRegistrations(t *testing.T) {
	directory, pin := deliveredRegistryReader(t)
	f, in := newAutomaticIntegrationFixture(t, false)
	checked, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	in.Apply, in.Confirmation = true, integrationDigest(t, checked.Readback)
	if _, err = ApplyTaskIntegration(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	before, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	p, err := UpgradeRegistry(f.dependencies, directory, pin, "", false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = UpgradeRegistry(f.dependencies, directory, pin, p.Confirmation, true); err != nil {
		t.Fatal(err)
	}
	newPath := filepath.Join(f.wrapper, "second")
	runLocalGit(t, filepath.Join(f.wrapper, "main"), "worktree", "add", "-b", "second", newPath, f.resultOID)
	adopt, err := ParseEpicAdoptInput("second", "Unrelated Epic", "ply", "ply", newPath, "refs/heads/second", f.resultOID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = AdoptEpic(f.dependencies, adopt); err != nil {
		t.Fatal(err)
	}
	create, _ := ParseTaskCreateInput("fresh", "Fresh task", "Ordinary registration after migration", "epic", "ply", "ply")
	if _, err = CreateTask(f.dependencies, create); err != nil {
		t.Fatal(err)
	}
	goal := goalFixture(t, f.workItemJourneyFixture, "fresh", "fresh-goal")
	queue := goalQueueFixture(t, f.workItemJourneyFixture, []TaskGoalQueueEntry{{TaskID: "fresh", Goal: &goal}})
	if len(queue.Pending) != 1 || queue.Pending[0].State != "ready" {
		t.Fatalf("new goal is not ready: %+v", queue)
	}
	after, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	if after.FormatVersion != 5 || !reflect.DeepEqual(before.IntegrationAuthorities, after.IntegrationAuthorities) || !reflect.DeepEqual(before.TaskResults, after.TaskResults) || !reflect.DeepEqual(before.IntegrationResults, after.IntegrationResults) || !reflect.DeepEqual(before.TaskPreparations, after.TaskPreparations) {
		t.Fatal("ordinary writes changed preserved authority/results/preparations")
	}
	for _, old := range before.EpicBaseVersions {
		found := false
		for _, now := range after.EpicBaseVersions {
			if reflect.DeepEqual(old, now) {
				found = true
			}
		}
		if !found {
			t.Fatal("lost Epic base history")
		}
	}
	// Both attempts to reinterpret the migration or downgrade fail without bytes changing.
	saved, _ := os.ReadFile(workItemsPath(f.root))
	for _, mutate := range []func(*WorkItemRegistry){
		func(r *WorkItemRegistry) { r.FormatVersion = 4; r.FormatBoundary = nil },
		func(r *WorkItemRegistry) {
			copy := *r.FormatBoundary
			copy.ReaderVersion = "v99.0.0"
			r.FormatBoundary = &copy
		},
	} {
		err := f.dependencies.WorkItems.WithLock(f.root, func(s WorkItemStoreSession) error {
			r, e := s.Snapshot()
			if e != nil {
				return e
			}
			mutate(&r)
			return s.Publish(r)
		})
		if err == nil {
			t.Fatal("history mutation accepted")
		}
		got, _ := os.ReadFile(workItemsPath(f.root))
		if !bytes.Equal(saved, got) {
			t.Fatal("rejection changed bytes")
		}
	}
}

// Authorization source is additive: absent legacy fields must stay absent,
// while an explicit receipt source must survive the same strict YAML boundary.
func TestDeliveryAuthorizationRegistryShapePreservesLegacyAndReceipt(t *testing.T) {
	for _, source := range []string{"", DeliveryMandateReceiptV1} {
		t.Run(source, func(t *testing.T) {
			auth := DeliveryAuthorization{AuthorityKind: source, PermissionConfirmed: source == "", AllowedEffects: []string{DeliveryLocalEpic}}
			raw, err := yaml.Marshal(auth)
			if err != nil {
				t.Fatal(err)
			}
			var document yaml.Node
			if err = yaml.Unmarshal(raw, &document); err != nil {
				t.Fatal(err)
			}
			if err = validateYAMLNodeShape(document.Content[0], reflect.TypeOf(auth), "delivery_authorization"); err != nil {
				t.Fatalf("strict registry rejected supported authorization source %q: %v", source, err)
			}
			if bytes.Contains(raw, []byte("authority_kind:")) != (source != "") {
				t.Fatal("legacy authorization acquired a new source claim")
			}
			var read DeliveryAuthorization
			if err = yaml.Unmarshal(raw, &read); err != nil || !reflect.DeepEqual(auth, read) {
				t.Fatalf("authorization source changed on roundtrip: %v", err)
			}
			bad := append(append([]byte{}, raw...), []byte("unknown_authorizer: true\n")...)
			if err = yaml.Unmarshal(bad, &document); err != nil {
				t.Fatal(err)
			}
			if err = validateYAMLNodeShape(document.Content[0], reflect.TypeOf(auth), "delivery_authorization"); err == nil {
				t.Fatal("unknown authority member bypassed strict registry schema")
			}
		})
	}
}
