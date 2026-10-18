package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Attached history is retained through its hash-bound locators. A bare
// manifest_sha256 field in such JSON is informative, wherever the document
// lives; only the native Spec basis binding requires the local manifest store.
func TestCloseoutRetentionRequiresOnlyNativeSpecManifests(t *testing.T) {
	for _, tc := range []struct {
		name, historyIn, defect, wantErr string
		nativeBasis                      bool
	}{
		{name: "external informative history", historyIn: "external"},
		{name: "local informative history", historyIn: "local"},
		{name: "native basis manifest retained", historyIn: "local", nativeBasis: true},
		{name: "native basis manifest must exist locally", historyIn: "external", nativeBasis: true, defect: "missing local manifest", wantErr: "no such file"},
		{name: "native basis manifest hash checked", historyIn: "local", nativeBasis: true, defect: "tampered local manifest", wantErr: "evidence hash differs"},
		{name: "external history hash checked", historyIn: "external", defect: "tampered history", wantErr: "evidence hash differs"},
		{name: "external history must exist", historyIn: "external", defect: "missing history", wantErr: "no such file"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			physical := func() string {
				dir, err := filepath.EvalSymlinks(t.TempDir())
				if err != nil {
					t.Fatal(err)
				}
				return dir
			}
			root, external := physical(), physical()
			write := func(path string, raw []byte) string {
				t.Helper()
				if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, raw, 0600); err != nil {
					t.Fatal(err)
				}
				return digestTaskBytes(raw)
			}
			// The informative manifest exists, unchanged and hash-valid, only in
			// the external workspace that wrote the history.
			informative := []byte(`{"kind":"external spec manifest"}`)
			informativeSHA := digestTaskBytes(informative)
			write(taskContentPath(external, "manifests", informativeSHA), informative)
			historyRoot := external
			if tc.historyIn == "local" {
				historyRoot = root
			}
			history, err := json.Marshal(map[string]any{"kind": "stop history", "task_spec_binding": map[string]any{"spec": map[string]any{"manifest_sha256": informativeSHA, "revision": 1}}})
			if err != nil {
				t.Fatal(err)
			}
			historyPath := filepath.Join(historyRoot, "history", "stop.json")
			historySHA := write(historyPath, history)
			switch tc.defect {
			case "tampered history":
				write(historyPath, []byte(`{"changed":true}`))
			case "missing history":
				if err := os.Remove(historyPath); err != nil {
					t.Fatal(err)
				}
			}
			stop, err := json.Marshal(map[string]any{"history": map[string]any{"locator": historyPath, "sha256": historySHA}})
			if err != nil {
				t.Fatal(err)
			}
			stopPath := filepath.Join(root, "candidate", "stop-evidence")
			stopSHA := write(stopPath, stop)
			owner := filepath.Join(root, "owner")
			ownerSHA := write(owner, []byte("fixture native control evidence"))
			write(workItemsPath(root), []byte("fixture native registry bytes"))
			tr := automaticResultFixture()
			tr.HandoffLocator, tr.HandoffSHA256 = owner, ownerSHA
			tr.StartReceiptLocator, tr.StartReceiptSHA256 = owner, ownerSHA
			tr.TerminalResultLocator, tr.TerminalResultSHA256 = owner, ownerSHA
			tr.Artifacts = []TaskArtifactRecord{{ArtifactID: "stop-evidence", Role: "other", Locator: stopPath, SHA256: stopSHA, SizeBytes: int64(len(stop))}}
			in := TaskCloseoutInput{Ownership: TaskCloseoutOwnership{EvidenceLocator: owner, EvidenceSHA256: ownerSHA}}
			r := WorkItemRegistry{}
			native := []byte(`{"kind":"native spec manifest"}`)
			nativeSHA := digestTaskBytes(native)
			if tc.nativeBasis {
				write(taskContentPath(root, "manifests", nativeSHA), native)
				switch tc.defect {
				case "missing local manifest":
					if err := os.Remove(taskContentPath(root, "manifests", nativeSHA)); err != nil {
						t.Fatal(err)
					}
				case "tampered local manifest":
					write(taskContentPath(root, "manifests", nativeSHA), []byte(`{"changed":true}`))
				}
				ref := TaskRevisionRef{ManifestSHA256: nativeSHA, Revision: 1}
				decision := TaskDecisionRef{ManifestSHA256: nativeSHA}
				r.TaskResultSpecBindings = []TaskResultSpecBinding{{TaskResultID: tr.ID, TaskID: tr.TaskID, Basis: TaskSpecBasis{Problem: ref, Spec: ref, Assessment: decision, Selection: decision}}}
			}

			entries, err := collectCloseoutRetention(Dependencies{}, root, r, tr, in)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("want diagnostic %q, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			retained := map[string]closeoutEvidence{}
			for _, e := range entries {
				retained[e.origin] = e
			}
			if e, ok := retained[historyPath]; !ok || e.sha != historySHA || string(e.bytes) != string(history) {
				t.Fatalf("history was not retained with its original reference and bytes: %+v", e)
			}
			if _, ok := retained[taskContentPath(root, "manifests", informativeSHA)]; ok {
				t.Fatal("informative manifest_sha256 became a local manifest dependency")
			}
			if e, ok := retained[taskContentPath(root, "manifests", nativeSHA)]; tc.nativeBasis && (!ok || e.role != "spec_manifest" || string(e.bytes) != string(native)) {
				t.Fatalf("native basis manifest was not retained: %+v", e)
			}
			if b, err := os.ReadFile(historyPath); err != nil || string(b) != string(history) {
				t.Fatalf("history changed: %v", err)
			}
		})
	}
}
