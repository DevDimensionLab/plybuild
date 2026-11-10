package workspace

import (
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func persistedPhysicalDigest(t *testing.T, root string) string {
	t.Helper()
	raw, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	return digestTaskBytes(raw)
}
func assertPersistedDigest(t *testing.T, v canonicaljson.Value, expected string) {
	t.Helper()
	obj := v.(canonicaljson.Object)
	if got := integrationReadbackMember(t, obj, "sha256"); got != expected {
		t.Fatalf("persisted hash=%v, physical=%s", got, expected)
	}
	if got := integrationReadbackMember(t, obj, "freshness"); got != "fresh" {
		t.Fatalf("observed snapshot freshness=%v", got)
	}
}
func TestIntegrationEvidencePersistedReadbackUsesPhysicalDigest(t *testing.T) {
	for _, legacy := range []bool{false} {
		name := "current-envelope"
		if legacy {
			name = "legacy"
		}
		t.Run(name, func(t *testing.T) {
			f := newIntegrationJourneyFixture(t)
			r, err := f.dependencies.WorkItems.Snapshot(f.root)
			if err != nil {
				t.Fatal(err)
			}
			physical := persistedPhysicalDigest(t, f.root)
			assertPersistedDigest(t, integrationPersistedView(r, nil, nil, nil, nil), physical)
			nullable := newNullableIntegrationReadback(f.dependencies, f.root, ProjectSnapshot{}, r, r.Tasks[0], r.Epics[0], r.Epics[0].RepoBindings[0], "blocked", IntegrationNextAction{})
			assertPersistedDigest(t, integrationReadbackMember(t, nullable.Value, "persisted_after"), physical)
			git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit, lost: true}
			f.dependencies.IntegrationGit = git
			in := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
			check, err := CheckTaskIntegration(f.dependencies, in)
			if err != nil {
				t.Fatal(err)
			}
			assertPersistedDigest(t, integrationReadbackMember(t, check.Readback.Value, "persisted_before"), physical)
			// A durable registry publication with a lost directory-sync response must
			// return its observed physical generation, alongside the lost Git reply.
			fired := false
			f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(stage string) error {
				if stage == "directory-sync" && !fired {
					fired = true
					return errInjected
				}
				return nil
			}}}
			in.Apply, in.Confirmation = true, integrationDigest(t, check.Readback)
			out, err := ApplyTaskIntegration(f.dependencies, in)
			if err != nil {
				t.Fatal(err)
			}
			assertPersistedDigest(t, integrationReadbackMember(t, out.Readback.Value, "persisted_before"), physical)
			after := persistedPhysicalDigest(t, f.root)
			assertPersistedDigest(t, integrationReadbackMember(t, out.Readback.Value, "persisted_after"), after)
			replay, err := ApplyTaskIntegration(f.dependencies, in)
			if err != nil {
				t.Fatal(err)
			}
			assertPersistedDigest(t, integrationReadbackMember(t, replay.Readback.Value, "persisted_after"), after)
			shown, err := ShowTask(f.dependencies, "task")
			if err != nil || shown.Integration == nil {
				t.Fatal("show", err)
			}
			assertPersistedDigest(t, integrationReadbackMember(t, shown.Integration.Value, "persisted_after"), after)
			if git.count() != 1 {
				t.Fatal("readback recovery repeated Git effect")
			}
		})
	}
}
func TestIntegrationEvidenceUnobservedReadbackIsUnknown(t *testing.T) {
	v := integrationPersistedView(emptyWorkItemRegistry(), nil, nil, nil, nil).(canonicaljson.Object)
	if integrationReadbackMember(t, v, "sha256") != nil || integrationReadbackMember(t, v, "freshness") != "unknown" {
		t.Fatal("unpublished model claimed physical bytes")
	}
}

// Cross-package journal regression fixture. The child exports only synthetic
// private workspace bytes; the parent restores their original private paths
// after testing's cleanup so Git/workspace identities remain exact.
func TestIntegrationEvidenceExportJournalFixture(t *testing.T) {
	dest := os.Getenv("PLY_INTEGRATION_JOURNAL_FIXTURE")
	if dest == "" {
		t.Skip("private cross-package fixture not requested")
	}
	f := newIntegrationJourneyFixture(t)
	in := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
	check, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	in.Apply, in.Confirmation = true, integrationDigest(t, check.Readback)
	if _, err := ApplyTaskIntegration(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	err = filepath.WalkDir(f.root, func(path string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(f.root, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, "archive", rel)
		info, err := e.Info()
		if err != nil {
			return err
		}
		if e.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		if e.Type()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(link, target)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, b, info.Mode().Perm())
	})
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(map[string]string{"root": f.root})
	if err := os.WriteFile(filepath.Join(dest, "fixture.json"), b, 0600); err != nil {
		t.Fatal(err)
	}
}
