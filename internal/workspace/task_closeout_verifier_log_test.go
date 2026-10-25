package workspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Verifier stdout/stderr are text data. Their exact bytes are retained, but a
// historical JSON citation printed by a test is no required cleanup evidence.
func TestCloseoutRetentionVerifierLogsAreOpaque(t *testing.T) {
	type testCase struct{ name, stream, edge, role, defect, want string }
	cases := []testCase{}
	for _, stream := range []string{"stdout", "stderr"} {
		for _, version := range []string{"1", "2"} {
			for _, edge := range []string{"projected", "receipt", "both"} {
				for _, role := range []string{"verifier_" + stream, "review"} {
					cases = append(cases, testCase{name: fmt.Sprintf("%s v%s %s %s", stream, version, edge, role), stream: stream, edge: edge + version, role: role})
				}
			}
		}
	}
	cases = append(cases,
		testCase{name: "missing original log", stream: "stdout", edge: "both2", defect: "missing original", want: "verify/stdout: no such file"},
		testCase{name: "missing log copy", stream: "stderr", edge: "both1", defect: "missing copy", want: "candidate/verifier-stderr: no such file"},
		testCase{name: "changed log", stream: "stdout", edge: "both2", defect: "changed", want: "evidence hash differs"},
		testCase{name: "wrong log size", stream: "stderr", edge: "projected2", defect: "size", want: "evidence size differs"},
		testCase{name: "independent native edge to same bytes", stream: "stdout", edge: "both2", defect: "native edge", want: "removed-test-scratch"},
		testCase{name: "arbitrary stdout field", stream: "stdout", edge: "both1", defect: "arbitrary field", want: "removed-test-scratch"},
		testCase{name: "self-declared receipt kind", stream: "stderr", edge: "both2", defect: "declared kind", want: "removed-test-scratch"},
		testCase{name: "artifact name without verifier binding", stream: "stdout", edge: "projected1", defect: "unbound name", want: "removed-test-scratch"},
	)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			write := func(name string, raw []byte) (string, string) {
				t.Helper()
				p := filepath.Join(root, name)
				if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(p, raw, 0600); err != nil {
					t.Fatal(err)
				}
				return p, digestTaskBytes(raw)
			}
			marshal := func(v any) []byte {
				t.Helper()
				raw, err := json.Marshal(v)
				if err != nil {
					t.Fatal(err)
				}
				return raw
			}
			// The log looks like a native document and cites removed test scratch.
			citation := filepath.Join(root, "removed-test-scratch", "citation.json")
			cited := map[string]any{"locator": citation, "sha256": digestTaskBytes([]byte("historical scratch"))}
			log := marshal(map[string]any{"kind": "PlyDeliveryVerification@2", "schema_version": 2, "stdout": cited, "diagnostic": cited})
			original, logSHA := write("verify/"+tc.stream, log)
			copy, _ := write("candidate/verifier-"+tc.stream, log)
			id := "verifier-" + tc.stream
			version := tc.edge[len(tc.edge)-1:]
			edge := tc.edge[:len(tc.edge)-1]
			receipt := marshal(map[string]any{"kind": "PlyDeliveryVerification@" + version, "schema_version": int(version[0] - '0'), tc.stream: map[string]any{"locator": original, "sha256": logSHA}})
			receiptPath, receiptSHA := write("candidate/verification-receipt", receipt)
			descriptors := []any{map[string]any{"artifact_id": id, "kind": "managed", "media_type": "text/plain", "locator": copy, "sha256": logSHA, "size_bytes": len(log)}}
			tr := automaticResultFixture()
			tr.Artifacts = nil
			if edge != "receipt" {
				size := int64(len(log))
				if tc.defect == "size" {
					size++
				}
				tr.Artifacts = append(tr.Artifacts, TaskArtifactRecord{ArtifactID: id, Role: tc.role, Locator: copy, SHA256: logSHA, SizeBytes: size})
			}
			if edge != "projected" {
				descriptors = append(descriptors, map[string]any{"artifact_id": "verification-receipt", "kind": "managed", "media_type": "application/json", "locator": receiptPath, "sha256": receiptSHA, "size_bytes": len(receipt)})
				tr.Artifacts = append(tr.Artifacts, TaskArtifactRecord{ArtifactID: "verification-receipt", Role: "other", Locator: receiptPath, SHA256: receiptSHA, SizeBytes: int64(len(receipt))})
			}
			if tc.defect != "unbound name" {
				if tc.stream == "stdout" {
					tr.VerifierResults[0].StdoutArtifactID = &id
				} else {
					tr.VerifierResults[0].StderrArtifactID = &id
				}
			}
			terminal := map[string]any{"artifacts": descriptors}
			switch tc.defect {
			case "native edge":
				terminal["native"] = map[string]any{"locator": copy, "sha256": logSHA}
			case "arbitrary field":
				terminal["stdout"] = map[string]any{"locator": original, "sha256": logSHA}
			case "declared kind":
				fake := marshal(map[string]any{"kind": "PlyDeliveryVerification@2", "schema_version": 2, "fake": true, tc.stream: map[string]any{"locator": original, "sha256": logSHA}})
				path, sha := write("declared-receipt", fake)
				terminal["evidence"] = map[string]any{"locator": path, "sha256": sha}
			}
			terminalPath, terminalSHA := write("terminal", marshal(terminal))
			owner, ownerSHA := write("owner", marshal(map[string]any{"fixture": "native control"}))
			write(filepath.Join(MarkerDirectory, WorkItemsFile), []byte("native registry fixture"))
			switch tc.defect {
			case "missing original":
				if err := os.Remove(original); err != nil {
					t.Fatal(err)
				}
			case "missing copy":
				if err := os.Remove(copy); err != nil {
					t.Fatal(err)
				}
			case "changed":
				write("verify/"+tc.stream, []byte("changed"))
			}
			tr.HandoffLocator, tr.HandoffSHA256 = owner, ownerSHA
			tr.StartReceiptLocator, tr.StartReceiptSHA256 = owner, ownerSHA
			tr.TerminalResultLocator, tr.TerminalResultSHA256 = terminalPath, terminalSHA
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
			retained := map[string]bool{}
			for _, e := range entries {
				if strings.HasPrefix(e.origin, filepath.Dir(citation)) {
					t.Fatal("historical log citation became required evidence")
				}
				if e.origin == original || e.origin == copy {
					if !bytes.Equal(e.bytes, log) || e.sha != logSHA {
						t.Fatal("log bytes/hash changed")
					}
					retained[e.origin] = true
				}
			}
			if !retained[copy] || edge != "projected" && !retained[original] {
				t.Fatalf("exact log not retained: %v", retained)
			}
			for _, p := range []string{original, copy} {
				if actual, err := os.ReadFile(p); err != nil || !bytes.Equal(actual, log) {
					t.Fatal("log rewritten", err)
				}
			}
		})
	}
}
