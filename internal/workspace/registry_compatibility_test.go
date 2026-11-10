package workspace

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
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

func TestRegistryFormatChangeRequiresNativePermit(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	queueSetFixture(t, f, "format-fixture", []TaskGoalQueueEntry{})
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
			queueSetFixture(t, f, "format-fixture", []TaskGoalQueueEntry{})
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
			queueSetFixture(t, f, "compatible-format-5-write", []TaskGoalQueueEntry{})
		})
	}
}

func TestRegistryDeliveredTransitionStaleAndConcurrent(t *testing.T) {
	directory, pin := deliveredRegistryReader(t)
	f := newWorkItemJourneyFixture(t)
	queueSetFixture(t, f, "format-fixture", []TaskGoalQueueEntry{})
	p, err := UpgradeRegistry(f.dependencies, directory, pin, "", false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = UpgradeRegistry(f.dependencies, directory, pin, "", true); err == nil {
		t.Fatal("missing confirmation accepted")
	}
	queueSetFixture(t, f, "competing-write", []TaskGoalQueueEntry{})
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
			queueSetFixture(t, f, "format-fixture", []TaskGoalQueueEntry{})
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
			queueSetFixture(t, f, "format-fixture", []TaskGoalQueueEntry{})
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
