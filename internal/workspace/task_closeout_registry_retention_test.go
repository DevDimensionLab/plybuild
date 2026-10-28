package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func closeoutContentStats(t *testing.T, root string) (int, int64) {
	t.Helper()
	files, err := os.ReadDir(filepath.Join(taskContentRoot(root), "closeout", "sha256"))
	if err != nil {
		t.Fatal(err)
	}
	var size int64
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			t.Fatal(err)
		}
		size += info.Size()
	}
	return len(files), size
}

func TestTaskCloseoutRetentionPreservesExplicitRegistryEvidence(t *testing.T) {
	f, in := closeoutFixture(t, true, false)
	r, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(workItemsPath(f.root))
	if err != nil {
		t.Fatal(err)
	}
	tr := *findTaskResult(r, f.resultID)
	tr.Artifacts = append(tr.Artifacts, TaskArtifactRecord{ArtifactID: "registry-input", Role: "other", Locator: workItemsPath(f.root), SHA256: digestTaskBytes(raw), SizeBytes: int64(len(raw))})
	entries, err := collectCloseoutRetention(f.dependencies, f.root, r, tr, in)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.origin == workItemsPath(f.root) && entry.sha == digestTaskBytes(raw) && string(entry.bytes) == string(raw) {
			return
		}
	}
	t.Fatal("an explicit Task evidence reference to registry bytes was discarded")
}

func TestTaskCloseoutReadsLegacyRegistrySnapshot(t *testing.T) {
	f, in := closeoutFixture(t, true, false)
	legacyRegistry, err := os.ReadFile(workItemsPath(f.root))
	if err != nil {
		t.Fatal(err)
	}
	p, err := PreviewTaskCloseout(f.dependencies, in)
	if err != nil || !p.Ready {
		t.Fatalf("preview: %+v, %v", p, err)
	}
	receipt, err := ApplyTaskCloseout(f.dependencies, in, p.PlanSHA256)
	if err != nil {
		t.Fatal(err)
	}
	sha, err := f.dependencies.TaskContent.Publish(f.root, "closeout", legacyRegistry)
	if err != nil {
		t.Fatal(err)
	}
	manifest := *receipt.Retention
	manifest.ManifestLocator, manifest.ManifestSHA256 = "", ""
	manifest.Entries = append(manifest.Entries, TaskRetainedEvidence{Role: "native_registry", OriginLocator: workItemsPath(f.root), SHA256: sha, SizeBytes: int64(len(legacyRegistry)), RetainedLocator: taskContentPath(f.root, "closeout", sha)})
	encoded, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	manifest.ManifestSHA256, err = f.dependencies.TaskContent.Publish(f.root, "closeout", encoded)
	if err != nil {
		t.Fatal(err)
	}
	manifest.ManifestLocator = taskContentPath(f.root, "closeout", manifest.ManifestSHA256)
	store, err := readTaskCloseouts(f.root)
	if err != nil {
		t.Fatal(err)
	}
	store.receipt(in.TaskID).Retention = &manifest
	if err := store.publish(f.root); err != nil {
		t.Fatal(err)
	}
	retained, err := ReadTaskRetainedEvidence(f.root, in.TaskID, workItemsPath(f.root), sha)
	if err != nil || string(retained) != string(legacyRegistry) {
		t.Fatalf("legacy registry bytes were not preserved: %v", err)
	}
	if err := verifyCloseoutRetention(f.dependencies, f.root, p.Plan, manifest); err != nil {
		t.Fatal(err)
	}
	if _, err := ApplyTaskCloseout(f.dependencies, in, p.PlanSHA256); err != nil {
		t.Fatalf("legacy complete closeout could not be read again: %v", err)
	}
}

func TestTaskCloseoutRetentionDoesNotCopyUnrelatedRegistryChanges(t *testing.T) {
	f, in := closeoutFixture(t, true, false)
	p, err := PreviewTaskCloseout(f.dependencies, in)
	if err != nil || !p.Ready {
		t.Fatalf("preview: %+v, %v", p, err)
	}
	retain := func() TaskRetentionManifest {
		t.Helper()
		r, err := f.dependencies.WorkItems.Snapshot(f.root)
		if err != nil {
			t.Fatal(err)
		}
		entries, err := collectCloseoutRetention(f.dependencies, f.root, r, *findTaskResult(r, f.resultID), in)
		if err != nil {
			t.Fatal(err)
		}
		manifest, err := retainCloseoutEvidence(f.dependencies, f.root, p.Plan, entries)
		if err != nil {
			t.Fatal(err)
		}
		return manifest
	}
	first := retain()
	beforeCount, beforeBytes := closeoutContentStats(t, f.root)
	if _, err := CreateTask(f.dependencies, TaskCreateInput{TaskID: "unrelated", Title: "Unrelated task", Description: "Must not change another Task's retained evidence", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
		t.Fatal(err)
	}
	second := retain()
	afterCount, afterBytes := closeoutContentStats(t, f.root)
	t.Logf("closeout objects/bytes before unrelated Task: %d/%d; after: %d/%d", beforeCount, beforeBytes, afterCount, afterBytes)
	if first.ManifestSHA256 != second.ManifestSHA256 || beforeCount != afterCount || beforeBytes != afterBytes {
		t.Fatal("unrelated registry change created a new closeout evidence snapshot")
	}
	for _, entry := range second.Entries {
		if entry.Role == "native_registry" || entry.OriginLocator == workItemsPath(f.root) {
			t.Fatal("closeout retained the workspace registry without a Task evidence reference")
		}
	}
}
