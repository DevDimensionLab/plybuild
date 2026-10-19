package workspace

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCloseoutRetentionInformativeDocumentsAreOpaque(t *testing.T) {
	for _, tc := range []struct{ name, media, defect, want string }{
		{name: "plain historical JSON", media: "text/plain"},
		{name: "markdown historical JSON", media: "text/markdown"},
		{name: "missing document", media: "text/plain", defect: "missing", want: "no such file"},
		{name: "changed document", media: "text/plain", defect: "changed", want: "evidence hash differs"},
		{name: "wrong size", media: "text/plain", defect: "size", want: "evidence size differs"},
		{name: "native JSON remains required", media: "application/json", want: "conflicting evidence references"},
		{name: "untyped evidence remains required", want: "conflicting evidence references"},
		{name: "independent native reference still required", media: "text/plain", defect: "native", want: "conflicting evidence references"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			write := func(name string, raw []byte) string {
				t.Helper()
				p := filepath.Join(root, name)
				if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(p, raw, 0600); err != nil {
					t.Fatal(err)
				}
				return p
			}
			historical := []byte("historical acceptance version one")
			origin := write("old/acceptance.sh", historical)
			// The first citation exists and matches. The second is a legitimate older
			// version at the same origin, reproducing the observed path/hash conflict.
			report, err := json.Marshal(map[string]any{"kind": "PlyDeliveryVerification@2", "citations": []any{
				map[string]any{"locator": origin, "sha256": digestTaskBytes(historical)},
				map[string]any{"locator": origin, "sha256": digestTaskBytes([]byte("historical version two"))},
			}})
			if err != nil {
				t.Fatal(err)
			}
			document := write("frozen-report", report)
			descriptor := map[string]any{"role": "other", "locator": document, "sha256": digestTaskBytes(report), "size_bytes": len(report)}
			if tc.media != "" {
				descriptor["media_type"] = tc.media
			}
			if tc.defect == "size" {
				descriptor["size_bytes"] = len(report) + 1
			}
			control := map[string]any{"inputs": []any{descriptor}}
			if tc.defect == "native" {
				control["required_native"] = map[string]any{"locator": document, "sha256": digestTaskBytes(report)}
			}
			raw, err := json.Marshal(control)
			if err != nil {
				t.Fatal(err)
			}
			owner := write("owner", raw)
			write(filepath.Join(MarkerDirectory, WorkItemsFile), []byte("native registry fixture"))
			switch tc.defect {
			case "missing":
				if err := os.Remove(document); err != nil {
					t.Fatal(err)
				}
			case "changed":
				write("frozen-report", []byte("changed"))
			}
			tr := automaticResultFixture()
			tr.HandoffLocator, tr.HandoffSHA256 = owner, digestTaskBytes(raw)
			tr.StartReceiptLocator, tr.StartReceiptSHA256 = owner, digestTaskBytes(raw)
			tr.TerminalResultLocator, tr.TerminalResultSHA256 = owner, digestTaskBytes(raw)
			tr.Artifacts = nil
			in := TaskCloseoutInput{Ownership: TaskCloseoutOwnership{EvidenceLocator: owner, EvidenceSHA256: digestTaskBytes(raw)}}
			entries, err := collectCloseoutRetention(Dependencies{}, root, WorkItemRegistry{}, tr, in)
			if tc.want != "" {
				if err == nil || !strings.Contains(err.Error(), tc.want) {
					t.Fatalf("want %q, got %v", tc.want, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			found := false
			for _, e := range entries {
				if e.origin == origin {
					t.Fatal("historical citation became native evidence")
				}
				if e.origin == document {
					found = true
					if !bytes.Equal(e.bytes, report) || e.sha != digestTaskBytes(report) {
						t.Fatal("report bytes/hash changed")
					}
				}
			}
			if !found {
				t.Fatal("informative report not retained")
			}
			actual, err := os.ReadFile(document)
			if err != nil || !bytes.Equal(actual, report) {
				t.Fatal("original report changed", err)
			}
		})
	}
}

func TestCloseoutRetentionObservedInputUsesExactHandoff(t *testing.T) {
	for _, tc := range []struct{ name, defect, want string }{
		{name: "receipt before declaration"},
		{name: "declaration before receipt", defect: "handoff first"},
		{name: "missing historical citation is data", defect: "missing citation"},
		{name: "wrong input ID", defect: "id", want: "observed input differs"},
		{name: "wrong input path", defect: "path", want: "observed input differs"},
		{name: "wrong input hash", defect: "hash", want: "observed input differs"},
		{name: "wrong input type", defect: "type", want: "observed input differs"},
		{name: "native declared type cannot inherit data", defect: "native type", want: "conflicting evidence references"},
		{name: "wrong observed size", defect: "size", want: "evidence size differs"},
		{name: "wrong handoff hash", defect: "handoff hash", want: "no retained hash-bound handoff"},
		{name: "changed handoff", defect: "handoff bytes", want: "evidence hash differs"},
		{name: "separate native edge in same receipt", defect: "native edge", want: "conflicting evidence references"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			write := func(name string, v any) (string, string) {
				t.Helper()
				raw, err := json.Marshal(v)
				if err != nil {
					t.Fatal(err)
				}
				p := filepath.Join(root, name)
				if err := os.WriteFile(p, raw, 0600); err != nil {
					t.Fatal(err)
				}
				return p, digestTaskBytes(raw)
			}
			origin, first := write("historical", map[string]any{"version": 1})
			data, dataSHA := write("report", map[string]any{"kind": "ply.workflow.handoff", "citations": []any{map[string]any{"locator": origin, "sha256": first}, map[string]any{"locator": origin, "sha256": digestTaskBytes([]byte("second historical version"))}}})
			reportBytes, err := os.ReadFile(data)
			if err != nil {
				t.Fatal(err)
			}
			media := "text/plain"
			if tc.defect == "native type" {
				media = "application/json"
			}
			declaration := map[string]any{"id": "report", "role": "other", "locator": data, "sha256": dataSHA, "size_bytes": len(reportBytes), "media_type": media}
			handoff, handoffSHA := write("bound-handoff", map[string]any{"kind": "ply.workflow.handoff", "inputs": []any{declaration}})
			observed := map[string]any{"id": "report", "locator": data, "sha256": dataSHA, "size_bytes": len(reportBytes), "matches_expected": true}
			receipt := map[string]any{"kind": "ply.workflow.start-receipt", "binding": map[string]any{"handoff_sha256": handoffSHA}, "observed_inputs": []any{observed}}
			switch tc.defect {
			case "id":
				observed["id"] = "another"
			case "path":
				observed["locator"] = filepath.Join(root, "another")
			case "hash":
				observed["sha256"] = first
			case "type":
				observed["media_type"] = "application/json"
			case "size":
				observed["size_bytes"] = len(reportBytes) + 1
			case "handoff hash":
				receipt["binding"] = map[string]any{"handoff_sha256": first}
			case "handoff bytes":
				if err := os.WriteFile(handoff, []byte("changed handoff"), 0600); err != nil {
					t.Fatal(err)
				}
			case "native edge":
				receipt["other_native_evidence"] = map[string]any{"id": "report", "locator": data, "sha256": dataSHA}
			case "missing citation":
				if err := os.Remove(origin); err != nil {
					t.Fatal(err)
				}
			}
			receiptPath, receiptSHA := write("receipt", receipt)
			receiptRef := map[string]any{"locator": receiptPath, "sha256": receiptSHA}
			handoffRef := map[string]any{"locator": handoff, "sha256": handoffSHA}
			control := map[string]any{"a": receiptRef, "z": handoffRef}
			if tc.defect == "handoff first" {
				control["a"], control["z"] = handoffRef, receiptRef
			}
			owner, ownerSHA := write("owner", control)
			if err := os.MkdirAll(filepath.Dir(workItemsPath(root)), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(workItemsPath(root), []byte("native registry fixture"), 0600); err != nil {
				t.Fatal(err)
			}
			tr := automaticResultFixture()
			tr.HandoffLocator, tr.HandoffSHA256 = owner, ownerSHA
			tr.StartReceiptLocator, tr.StartReceiptSHA256 = owner, ownerSHA
			tr.TerminalResultLocator, tr.TerminalResultSHA256 = owner, ownerSHA
			tr.Artifacts = nil
			in := TaskCloseoutInput{Ownership: TaskCloseoutOwnership{EvidenceLocator: owner, EvidenceSHA256: ownerSHA}}
			entries, err := collectCloseoutRetention(Dependencies{}, root, WorkItemRegistry{}, tr, in)
			if tc.want != "" {
				if err == nil || !strings.Contains(err.Error(), tc.want) {
					t.Fatalf("want %q, got %v", tc.want, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			for _, e := range entries {
				if e.origin == origin {
					t.Fatal("observation expanded historical citation")
				}
			}
			actual, err := os.ReadFile(data)
			if err != nil || !bytes.Equal(actual, reportBytes) {
				t.Fatal("report rewritten", err)
			}
		})
	}
}
