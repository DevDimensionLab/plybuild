package workspace

// Measures published text payloads through isolated native fixtures.
import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

type writeAuditSection struct {
	Name      string `json:"name"`
	Before    int    `json:"before_bytes"`
	After     int    `json:"after_bytes"`
	Unchanged bool   `json:"unchanged"`
}
type writeAuditStep struct {
	Name           string              `json:"name"`
	Before         int                 `json:"registry_before_bytes"`
	After          int                 `json:"registry_after_bytes"`
	PublishedBytes int64               `json:"published_bytes"`
	PublishedFiles int                 `json:"published_files"`
	Sections       []writeAuditSection `json:"sections"`
}

// Slice the actual physical YAML at node line boundaries, preserving indentation
// and headers. The separate envelope prefix makes section totals exact.
func writeAuditSections(t *testing.T, b []byte) (map[string][]byte, []string) {
	t.Helper()
	var doc yaml.Node
	if err := yaml.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	root := doc.Content[0]
	version := registryNodeField(root, "format_version")
	if version == nil || (version.Value != "9" && version.Value != "10") {
		t.Fatal("measurement requires physical format 9 or 10")
	}
	reg := registryNodeField(root, "registry")
	lines := bytes.SplitAfter(b, []byte("\n"))
	offsets := []int{0}
	for _, line := range lines {
		offsets = append(offsets, offsets[len(offsets)-1]+len(line))
	}
	out := map[string][]byte{"envelope": b[:offsets[reg.Content[0].Line-1]]}
	order := []string{"envelope"}
	for i := 0; i < len(reg.Content); i += 2 {
		key := reg.Content[i]
		start, end := offsets[key.Line-1], len(b)
		if i+2 < len(reg.Content) {
			end = offsets[reg.Content[i+2].Line-1]
		}
		out[key.Value] = b[start:end]
		order = append(order, key.Value)
	}
	total := 0
	for _, v := range out {
		total += len(v)
	}
	if total != len(b) {
		t.Fatal("section accounting mismatch")
	}
	return out, order
}

func TestPublicationEvidenceWriteMeasurement(t *testing.T) {
	f := newIntegrationJourneyFixture(t)
	f.dependencies.TaskLifecycleIDs = cryptoTaskLifecycleIDSource{}
	f.dependencies.IntegrationGit = &noEffectIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	initialInventory := writeAuditInventory(t, f.root)
	var steps []writeAuditStep
	// Capture every publication, not just the final size of each command.
	native := f.dependencies.WorkItems
	observer := &writeAuditStore{WorkItemStore: native, t: t, steps: &steps}
	f.dependencies.WorkItems = observer
	in := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
	for i := 0; i < 8; i++ {
		observer.name = fmt.Sprintf("retry-%d", i+1)
		in.Apply, in.Confirmation = false, ""
		check, err := CheckTaskIntegration(f.dependencies, in)
		if err != nil {
			t.Fatal(err)
		}
		in.Apply, in.Confirmation = true, integrationDigest(t, check.Readback)
		out, err := ApplyTaskIntegration(f.dependencies, in)
		if err != nil || out.Readback.Classification != "no_effect" {
			t.Fatalf("retry: %s %v", out.Readback.Classification, err)
		}
		r, err := native.Snapshot(f.root)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range r.IntegrationResults {
			if (in.RetryAfterResultID == nil && v.RetryAfterResultID == nil) || (in.RetryAfterResultID != nil && v.RetryAfterResultID != nil && *in.RetryAfterResultID == *v.RetryAfterResultID) {
				id := v.ID
				in.RetryAfterResultID = &id
				break
			}
		}
	}
	beforeFiles := evidenceFileSnapshot(t, f.root)
	observer.name = "create-task"
	create, err := ParseTaskCreateInput("another-task", "Another Task", "Synthetic measurement", "epic", "ply", "ply")
	if err != nil {
		t.Fatal(err)
	}
	result, err := CreateTask(f.dependencies, create)
	if err != nil || !result.Created {
		t.Fatalf("create: %v", err)
	}
	for path, info := range beforeFiles {
		if filepath.Base(path) == WorkItemsFile {
			continue
		}
		after, err := os.Lstat(path)
		if err != nil || !os.SameFile(info, after) || !sameRegistryFileState(info, after) {
			t.Fatalf("historical payload rewritten: %s", path)
		}
	}
	observer.name = "create-task-replay"
	oldCount := len(steps)
	result, err = CreateTask(f.dependencies, create)
	if err != nil || result.Created || len(steps) != oldCount {
		t.Fatalf("replay must publish nothing: %v", err)
	}
	var stock int64
	for _, info := range evidenceFileSnapshot(t, f.root) {
		stock += info.Size()
	}
	final, err := native.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	report := struct {
		Steps              []writeAuditStep                    `json:"steps"`
		InitialInventory   map[string]writeAuditInventoryEntry `json:"initial_inventory"`
		FinalInventory     map[string]writeAuditInventoryEntry `json:"final_inventory"`
		FinalStock         int64                               `json:"final_ply_payload_stock_bytes"`
		Counts             map[string]int                      `json:"final_counts"`
		ReplayPublications int                                 `json:"task_replay_publications"`
	}{steps, initialInventory, writeAuditInventory(t, f.root), stock, map[string]int{"tasks": len(final.Tasks), "task_results": len(final.TaskResults), "human_qa_records": len(final.HumanQARecords), "authorities": len(final.IntegrationAuthorities), "intents": len(final.IntegrationIntents), "attempts": len(final.IntegrationAttempts), "results": len(final.IntegrationResults)}, len(steps) - oldCount}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if dest := os.Getenv("PLY_PUBLICATION_MEASUREMENT_OUTPUT"); dest != "" {
		if err := os.WriteFile(dest, append(data, '\n'), 0600); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("observed %d native publications; final payload stock %d; task replay publications %d", len(steps), stock, len(steps)-oldCount)
}

type writeAuditStore struct {
	WorkItemStore
	t     *testing.T
	name  string
	steps *[]writeAuditStep
}

func (s *writeAuditStore) WithLock(root string, fn func(WorkItemStoreSession) error) error {
	return s.WorkItemStore.WithLock(root, func(session WorkItemStoreSession) error {
		return fn(&writeAuditSession{WorkItemStoreSession: session, store: s, root: root})
	})
}

type writeAuditSession struct {
	WorkItemStoreSession
	store *writeAuditStore
	root  string
}

func (s *writeAuditSession) Publish(r WorkItemRegistry) error {
	t := s.store.t
	before, err := os.ReadFile(workItemsPath(s.root))
	if err != nil {
		t.Fatal(err)
	}
	files := evidenceFileSnapshot(t, s.root)
	if err := s.WorkItemStoreSession.Publish(r); err != nil {
		return err
	}
	for path, info := range files {
		if filepath.Base(path) == WorkItemsFile {
			continue
		}
		after, err := os.Lstat(path)
		if err != nil || !os.SameFile(info, after) || !sameRegistryFileState(info, after) {
			t.Fatalf("publication rewrote historical payload: %s %v", path, err)
		}
	}
	after, err := os.ReadFile(workItemsPath(s.root))
	if err != nil {
		t.Fatal(err)
	}
	n, count := evidencePublishedDelta(t, files, evidenceFileSnapshot(t, s.root))
	b, order := writeAuditSections(t, before)
	a, afterOrder := writeAuditSections(t, after)
	if len(order) != len(afterOrder) {
		t.Fatal("section shape changed")
	}
	step := writeAuditStep{Name: s.store.name, Before: len(before), After: len(after), PublishedBytes: n, PublishedFiles: count}
	for _, key := range order {
		step.Sections = append(step.Sections, writeAuditSection{key, len(b[key]), len(a[key]), bytes.Equal(b[key], a[key])})
	}
	if s.store.name == "create-task" {
		if dest := os.Getenv("PLY_PUBLICATION_MEASUREMENT_OUTPUT"); dest != "" {
			for name, payload := range map[string][]byte{"before-create.yaml": before, "after-create.yaml": after} {
				if err := os.WriteFile(filepath.Join(filepath.Dir(dest), name), payload, 0600); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
	*s.store.steps = append(*s.store.steps, step)
	return nil
}

// Classify retained .ply payloads separately from Go build output and Git data.
type writeAuditInventoryEntry struct {
	Files     int   `json:"files"`
	Bytes     int64 `json:"bytes"`
	TextFiles int   `json:"utf8_without_nul_files"`
}

func writeAuditInventory(t *testing.T, root string) map[string]writeAuditInventoryEntry {
	t.Helper()
	out := map[string]writeAuditInventoryEntry{}
	for path, info := range evidenceFileSnapshot(t, root) {
		relative, err := filepath.Rel(filepath.Join(root, MarkerDirectory), path)
		if err != nil {
			t.Fatal(err)
		}
		parts := strings.Split(filepath.ToSlash(relative), "/")
		group := parts[0]
		if len(parts) > 2 && group == "task-content" {
			group = strings.Join(parts[:3], "/")
		}
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		v := out[group]
		v.Files++
		v.Bytes += info.Size()
		if utf8.Valid(b) && !bytes.ContainsRune(b, 0) {
			v.TextFiles++
		}
		out[group] = v
	}
	return out
}
