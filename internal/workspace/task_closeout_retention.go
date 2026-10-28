package workspace

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

const closeoutEvidenceLimit = 64 << 20

type TaskRetainedEvidence struct {
	Role            string `json:"role"`
	OriginLocator   string `json:"origin_locator"`
	SHA256          string `json:"sha256"`
	SizeBytes       int64  `json:"size_bytes"`
	RetainedLocator string `json:"retained_locator"`
}

type TaskRetentionManifest struct {
	Kind            string                 `json:"kind"`
	SchemaVersion   int                    `json:"schema_version"`
	TaskID          TaskID                 `json:"task_id"`
	TaskResultID    TaskResultID           `json:"task_result_id"`
	CandidateOID    string                 `json:"candidate_oid"`
	CandidateTree   string                 `json:"candidate_tree"`
	ArchiveRef      string                 `json:"archive_ref"`
	Entries         []TaskRetainedEvidence `json:"entries"`
	ManifestLocator string                 `json:"manifest_locator,omitempty"`
	ManifestSHA256  string                 `json:"manifest_sha256,omitempty"`
}

type closeoutEvidence struct {
	role, origin, sha string
	bytes             []byte
}

func readCloseoutEvidence(locator, sha string) ([]byte, error) {
	if !validAbsoluteCleanPath(locator) || !digestPattern.MatchString(sha) {
		return nil, fmt.Errorf("invalid evidence locator or hash")
	}
	b, err := contentReadBounded(locator, closeoutEvidenceLimit, false)
	if err != nil {
		return nil, err
	}
	if digestTaskBytes(b) != sha {
		return nil, fmt.Errorf("evidence hash differs: %s", locator)
	}
	return b, nil
}

// ResolveTaskRetainedEvidence resolves an exact historical absolute reference
// without rewriting it. A matching validated retention manifest, origin, digest
// and readable retained object are required; an ID or missing source is not proof.
func ResolveTaskRetainedEvidence(root string, id TaskID, originLocator, sha string) (string, error) {
	r, err := ReadTaskCloseoutAt(root, id)
	if err != nil {
		return "", err
	}
	if r == nil || r.Retention == nil {
		return "", fmt.Errorf("retention manifest is unavailable")
	}
	m := *r.Retention
	if m.TaskID != id || m.TaskResultID != r.Plan.TaskResultID || m.CandidateOID != r.Plan.ResultOID || m.CandidateTree != r.Plan.ResultTree || m.ManifestLocator != taskContentPath(root, "closeout", m.ManifestSHA256) {
		return "", fmt.Errorf("retention binding differs")
	}
	b, err := readCloseoutEvidence(m.ManifestLocator, m.ManifestSHA256)
	if err != nil {
		return "", err
	}
	copy := m
	copy.ManifestLocator, copy.ManifestSHA256 = "", ""
	want, err := json.Marshal(copy)
	if err != nil || string(want) != string(b) {
		return "", fmt.Errorf("retention manifest differs")
	}
	for _, e := range m.Entries {
		if e.OriginLocator != originLocator || e.SHA256 != sha {
			continue
		}
		if e.RetainedLocator != taskContentPath(root, "closeout", sha) || pathWithin(r.Plan.Source.Locator, e.RetainedLocator) {
			return "", fmt.Errorf("retained path is not bound outside source")
		}
		b, err := readCloseoutEvidence(e.RetainedLocator, e.SHA256)
		if err != nil {
			return "", err
		}
		if int64(len(b)) != e.SizeBytes {
			return "", fmt.Errorf("retained size differs")
		}
		return e.RetainedLocator, nil
	}
	return "", fmt.Errorf("no retained evidence matches the exact original path and hash")
}

func ReadTaskRetainedEvidence(root string, id TaskID, originLocator, sha string) ([]byte, error) {
	path, err := ResolveTaskRetainedEvidence(root, id, originLocator, sha)
	if err != nil {
		return nil, err
	}
	return readCloseoutEvidence(path, sha)
}

// collectCloseoutRetention performs bounded read-only traversal. Historical
// manifests keep their original bytes and absolute references; the manifest maps
// each origin to an immutable copy outside the source instead of rewriting it.
func collectCloseoutRetention(d Dependencies, root string, r WorkItemRegistry, tr TaskResultRecord, in TaskCloseoutInput) ([]closeoutEvidence, error) {
	snapshots, err := automaticCloseoutSnapshots(r, tr)
	if err != nil {
		return nil, err
	}
	entries := []closeoutEvidence{}
	seen := map[string]closeoutEvidence{}
	expanded := map[string]bool{}
	type inputReference struct{ id, path, sha string }
	handoffInputs := map[string]map[inputReference]string{}
	// Receipts echo input identity without its type. Resolve only those echoes
	// after native declarations have been read, regardless of traversal order.
	type inputObservation struct {
		handoff string
		inputs  any
	}
	observations := []inputObservation{}
	// Verifier stdout/stderr are text data. Only their exact native bindings
	// are retained without expansion: the TaskResult verifier artifact IDs and
	// the log edges of the bound verification receipt. Names, roles and kinds
	// grant nothing; a separate native edge to the same bytes still expands.
	logs := map[string]bool{}
	for _, v := range tr.VerifierResults {
		for _, id := range []*string{v.StdoutArtifactID, v.StderrArtifactID} {
			if id != nil {
				logs[*id] = true
			}
		}
	}
	receipts := map[string]bool{}
	for _, a := range tr.Artifacts {
		if a.ArtifactID == "verification-receipt" {
			receipts[a.SHA256] = true
		}
	}
	total := 0
	var add func(string, string, string, int64, bool) error
	var walk func(any, map[inputReference]string) error
	add = func(role, path, sha string, size int64, follow bool) error {
		e, ok := seen[path]
		if ok {
			if e.sha != sha {
				return fmt.Errorf("conflicting evidence references: %s", path)
			}
		} else {
			if len(seen) >= 4096 {
				return fmt.Errorf("retention graph exceeds 4096 objects")
			}
			b, err := readCloseoutEvidence(path, sha)
			if err != nil {
				if preserved, ok := snapshots[path]; ok && preserved.sha == sha {
					b, err = preserved.bytes, nil
				}
			}
			if err != nil {
				return err
			}
			total += len(b)
			if total > 256<<20 {
				return fmt.Errorf("retained evidence exceeds 256 MiB")
			}
			e = closeoutEvidence{role, path, sha, b}
			seen[path] = e
			entries = append(entries, e)
		}
		if size >= 0 && int64(len(e.bytes)) != size {
			return fmt.Errorf("evidence size differs: %s", path)
		}
		// A document edge preserves data without interpreting its contents as
		// native control. A separate required evidence edge must still expand
		// the same bytes, even if the document was encountered first.
		if !follow || expanded[path] {
			return nil
		}
		expanded[path] = true
		var value any
		if json.Unmarshal(e.bytes, &value) == nil {
			if document, ok := value.(map[string]any); ok {
				switch document["kind"] {
				case "ply.workflow.handoff":
					inputs := map[inputReference]string{}
					declared, _ := document["inputs"].([]any)
					for _, v := range declared {
						in, ok := v.(map[string]any)
						if !ok {
							continue
						}
						id, _ := in["id"].(string)
						path, _ := in["locator"].(string)
						sha, _ := in["sha256"].(string)
						media, _ := in["media_type"].(string)
						inputs[inputReference{id, path, sha}] = media
					}
					handoffInputs[e.sha] = inputs
				case "ply.workflow.start-receipt":
					// Native observations omit the type. Only the exact bound
					// handoff input can supply it; arbitrary native refs cannot.
					binding, _ := document["binding"].(map[string]any)
					sha, _ := binding["handoff_sha256"].(string)
					observations = append(observations, inputObservation{sha, document["observed_inputs"]})
					// Only the decoded copy changes. Signed receipt bytes remain exact.
					delete(document, "observed_inputs")
				}
				if receipts[e.sha] {
					for _, key := range []string{"stdout", "stderr"} {
						log, ok := document[key].(map[string]any)
						if !ok {
							continue
						}
						path, _ := log["locator"].(string)
						sha, _ := log["sha256"].(string)
						if err := add("verifier_log", path, sha, -1, false); err != nil {
							return err
						}
						delete(document, key)
					}
				}
			}
			return walk(value, nil)
		}
		return nil
	}
	walk = func(value any, observed map[inputReference]string) error {
		switch v := value.(type) {
		case []any:
			for _, member := range v {
				if err := walk(member, observed); err != nil {
					return err
				}
			}
		case map[string]any:
			// Only hashed file references are followed; commands, working
			// directories and unverified historical origin paths are not files.
			keys := make([]string, 0, len(v))
			for key := range v {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				path, ok := v[key].(string)
				if !ok || !filepath.IsAbs(path) {
					continue
				}
				shaKey := ""
				if key == "locator" {
					shaKey = "sha256"
				} else if strings.HasSuffix(key, "_locator") && key != "origin_locator" {
					shaKey = strings.TrimSuffix(key, "_locator") + "_sha256"
				}
				sha, _ := v[shaKey].(string)
				if shaKey != "" && digestPattern.MatchString(sha) {
					// Document media types are bound by the enclosing native
					// descriptor. JSON-looking report bytes are still document data.
					media, _ := v["media_type"].(string)
					follow := key != "locator" || (media != "text/plain" && media != "text/markdown")
					if key == "locator" {
						id, _ := v["id"].(string)
						if observed != nil {
							boundMedia, ok := observed[inputReference{id, path, sha}]
							if !ok || (media != "" && media != boundMedia) {
								return fmt.Errorf("observed input differs from its bound handoff: %s", path)
							}
							follow = boundMedia != "text/plain" && boundMedia != "text/markdown"
						}
					}
					size := int64(-1)
					if key == "locator" {
						if declared, ok := v["size_bytes"].(float64); ok {
							size = int64(declared)
						}
					}
					if err := add("referenced_evidence", path, sha, size, follow); err != nil {
						return err
					}
				}
			}
			// A bare manifest_sha256 field is informative: it names no store.
			// Native Spec manifests are required from the TaskResult basis below.
			for _, key := range keys {
				if err := walk(v[key], nil); err != nil {
					return err
				}
			}
		}
		return nil
	}
	for _, e := range []struct{ role, path, sha string }{
		{"handoff", tr.HandoffLocator, tr.HandoffSHA256}, {"acceptance", tr.StartReceiptLocator, tr.StartReceiptSHA256}, {"terminal", tr.TerminalResultLocator, tr.TerminalResultSHA256}, {"ownership", in.Ownership.EvidenceLocator, in.Ownership.EvidenceSHA256},
	} {
		if err := add(e.role, e.path, e.sha, -1, true); err != nil {
			return nil, err
		}
	}
	for _, a := range tr.Artifacts {
		if err := add(a.Role, a.Locator, a.SHA256, a.SizeBytes, !logs[a.ArtifactID]); err != nil {
			return nil, err
		}
	}
	for _, qa := range r.HumanQARecords {
		if qa.TaskResultID != tr.ID {
			continue
		}
		for _, a := range qa.Evidence {
			if err := add("human_qa", a.Locator, a.SHA256, a.SizeBytes, true); err != nil {
				return nil, err
			}
		}
	}
	for _, b := range r.TaskResultSpecBindings {
		if b.TaskResultID != tr.ID {
			continue
		}
		for _, sha := range []string{b.Basis.Problem.ManifestSHA256, b.Basis.Spec.ManifestSHA256, b.Basis.Assessment.ManifestSHA256, b.Basis.Selection.ManifestSHA256} {
			if err := add("spec_manifest", taskContentPath(root, "manifests", sha), sha, -1, true); err != nil {
				return nil, err
			}
		}
	}
	if in.ObservedPR != nil {
		if err := add("pr_integration", in.ObservedPR.EvidenceLocator, in.ObservedPR.EvidenceSHA256, -1, true); err != nil {
			return nil, err
		}
	}
	for i := 0; i < len(observations); i++ {
		observation := observations[i]
		inputs, ok := handoffInputs[observation.handoff]
		if !ok {
			return nil, fmt.Errorf("observed inputs have no retained hash-bound handoff")
		}
		if err := walk(observation.inputs, inputs); err != nil {
			return nil, err
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].origin != entries[j].origin {
			return entries[i].origin < entries[j].origin
		}
		return entries[i].role < entries[j].role
	})
	return entries, nil
}

func retainCloseoutEvidence(d Dependencies, root string, p TaskCloseoutPlan, entries []closeoutEvidence) (TaskRetentionManifest, error) {
	m := TaskRetentionManifest{Kind: "WorkspaceTaskRetentionManifest@1", SchemaVersion: 1, TaskID: p.TaskID, TaskResultID: p.TaskResultID, CandidateOID: p.ResultOID, CandidateTree: p.ResultTree, ArchiveRef: p.ArchiveRef, Entries: []TaskRetainedEvidence{}}
	if d.TaskContent == nil {
		return m, fmt.Errorf("content store unavailable")
	}
	if pathWithin(p.Source.Locator, taskContentRoot(root)) {
		return m, fmt.Errorf("retention store is inside source")
	}
	for _, e := range entries {
		sha, err := d.TaskContent.retainCloseout(root, e.origin, e.bytes)
		if err != nil {
			return m, err
		}
		if sha != e.sha {
			return m, fmt.Errorf("retention digest differs")
		}
		m.Entries = append(m.Entries, TaskRetainedEvidence{e.role, e.origin, sha, int64(len(e.bytes)), taskContentPath(root, "closeout", sha)})
	}
	if err := closeoutArchiveCandidate(p); err != nil {
		return m, err
	}
	b, err := json.Marshal(m)
	if err != nil {
		return m, err
	}
	sha, err := d.TaskContent.Publish(root, "closeout", b)
	if err != nil {
		return m, err
	}
	m.ManifestSHA256, m.ManifestLocator = sha, taskContentPath(root, "closeout", sha)
	return m, verifyCloseoutRetention(d, root, p, m)
}

func verifyCloseoutRetention(d Dependencies, root string, p TaskCloseoutPlan, m TaskRetentionManifest) error {
	if m.Kind != "WorkspaceTaskRetentionManifest@1" || m.SchemaVersion != 1 || m.TaskID != p.TaskID || m.TaskResultID != p.TaskResultID || m.CandidateOID != p.ResultOID || m.CandidateTree != p.ResultTree || m.ArchiveRef != p.ArchiveRef || len(m.Entries) == 0 || m.ManifestLocator != taskContentPath(root, "closeout", m.ManifestSHA256) {
		return fmt.Errorf("retention binding differs")
	}
	manifest, err := readCloseoutEvidence(m.ManifestLocator, m.ManifestSHA256)
	if err != nil {
		return err
	}
	copy := m
	copy.ManifestLocator, copy.ManifestSHA256 = "", ""
	b, err := json.Marshal(copy)
	if err != nil || string(b) != string(manifest) {
		return fmt.Errorf("retention manifest differs")
	}
	for _, e := range m.Entries {
		if e.RetainedLocator != taskContentPath(root, "closeout", e.SHA256) || pathWithin(p.Source.Locator, e.RetainedLocator) {
			return fmt.Errorf("retention path differs or is inside source")
		}
		b, err := readCloseoutEvidence(e.RetainedLocator, e.SHA256)
		if err != nil {
			return err
		}
		if int64(len(b)) != e.SizeBytes {
			return fmt.Errorf("retained evidence size differs")
		}
	}
	ref := closeoutGit(p, p.RepositoryLocator, false, "rev-parse", "--verify", "--quiet", p.ArchiveRef)
	if ref.Err != nil || strings.TrimSpace(string(ref.Stdout)) != p.ResultOID {
		return fmt.Errorf("candidate archive ref is missing or changed")
	}
	if symbolic := closeoutGit(p, p.RepositoryLocator, false, "symbolic-ref", "--quiet", p.ArchiveRef); symbolic.Exit != 1 {
		return fmt.Errorf("candidate archive must be a direct ref")
	}
	return closeoutCheckCandidate(p)
}
