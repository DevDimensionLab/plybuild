package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fixedTaskEvidenceReader struct {
	evidence TaskHandoffEvidence
	reads    int
}

func (reader *fixedTaskEvidenceReader) ReadTaskEvidence(request TaskHandoffEvidenceRequest) (TaskHandoffEvidence, error) {
	reader.reads++
	return reader.evidence, nil
}

type taskDeliveryFixture struct {
	workItemJourneyFixture
	root, taskPath, draftPath, resultOID, resultTree string
	reader                                           *fixedTaskEvidenceReader
}

func newTaskDeliveryFixture(t *testing.T, legacy bool) taskDeliveryFixture {
	t.Helper()
	base := newWorkItemJourneyFixtureVersion(t, legacy)
	root := filepath.Dir(base.wrapper)
	taskPath := filepath.Join(base.wrapper, "task")
	created, err := CreateTaskWorktree(base.dependencies, TaskWorktreeCreateInput{TaskID: "task", Branch: "task", Path: taskPath, ExpectedParentOID: base.oid})
	if err != nil {
		t.Fatal(err)
	}
	var basis *TaskSpecBasis
	if !legacy {
		selected := fixtureSelectTaskSpec(t, base)
		basis = &selected
	}
	if err := os.WriteFile(filepath.Join(taskPath, "delivery.txt"), []byte("delivered\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, taskPath, "add", "delivery.txt")
	runLocalGit(t, taskPath, "commit", "-m", "delivery")
	resultOID := runLocalGit(t, taskPath, "rev-parse", "HEAD")
	resultTree := runLocalGit(t, taskPath, "rev-parse", "HEAD^{tree}")
	digest := "sha256:" + strings.Repeat("a", 64)
	handoff := TaskHandoffEvidenceRequest{
		ActivityID: "act_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", RunID: "run_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", HandoffID: "hnd_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		HandoffLocator: filepath.Join(root, "handoff.json"), HandoffSHA256: digest,
		StartReceiptID: "rcp_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", StartReceiptLocator: filepath.Join(root, "start.json"), StartReceiptSHA256: digest,
		TerminalResultID: "res_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", TerminalResultLocator: filepath.Join(root, "terminal.json"), TerminalResultSHA256: digest,
		InspectionSHA256: digest,
	}
	evidence := TaskHandoffEvidence{
		TaskSpecBasis: basis, TaskSpecValid: basis != nil, TaskRequirementsValid: true, AcceptedStartOutcome: "started",
		ActivityID: handoff.ActivityID, RunID: handoff.RunID, HandoffID: handoff.HandoffID, HandoffLocator: handoff.HandoffLocator, HandoffSHA256: handoff.HandoffSHA256,
		StartReceiptID: handoff.StartReceiptID, StartReceiptLocator: handoff.StartReceiptLocator, StartReceiptSHA256: handoff.StartReceiptSHA256,
		TerminalResultID: handoff.TerminalResultID, TerminalResultLocator: handoff.TerminalResultLocator, TerminalResultSHA256: handoff.TerminalResultSHA256, InspectionSHA256: handoff.InspectionSHA256,
		ReportedOutcome: "complete", TargetWorktree: taskPath, TargetRef: "refs/heads/task", ResultOID: resultOID, ResultTree: resultTree, GitCommonDir: created.Task.GitCommonDir,
		SchemaValid: true, DigestValid: true, LifecycleValid: true, BindingValid: true, CapabilityValid: true, PrincipalSessionValid: true, PolicyValid: true, EvidenceCoverageValid: true,
		VerifierResults: []TaskVerifierResultRecord{}, ExpectedVerifierIDs: []string{}, Review: TaskReviewRecord{Findings: []TaskReviewEntry{}, Fixes: []TaskReviewEntry{}, OpenActionableFindings: []TaskReviewEntry{}}, Artifacts: []TaskArtifactRecord{}, EvidenceGaps: []string{},
	}
	draft := fmt.Sprintf(`{"kind":"WorkspaceTaskResultRecordDraft@1","schema_version":1,"format":"json","format_version":1,"canonicalization":"RFC8785","publication_key":"task/result-delivery","task_id":"task","task_worktree_id":%q,"handoff":{"activity_id":%q,"run_id":%q,"handoff_id":%q,"handoff_locator":%q,"handoff_sha256":%q,"start_receipt_id":%q,"start_receipt_locator":%q,"start_receipt_sha256":%q,"terminal_result_id":%q,"terminal_result_locator":%q,"terminal_result_sha256":%q,"inspection_sha256":%q},"source":{"project_id":"ply","repo_id":"ply","git_common_dir":%q,"worktree_locator":%q,"source_ref":"refs/heads/task","result_oid":%q,"result_tree":%q},"technical_assessment":{"gate":"passed","required_verifier_ids":[],"accepted_debt":[]},"evidence_artifacts":[],"recorder":{"actor_claim":"local operator","control_surface":"test CLI","recorded_at_utc":"2026-09-28T12:00:00Z"}}`,
		created.Task.Worktree.ID, handoff.ActivityID, handoff.RunID, handoff.HandoffID, handoff.HandoffLocator, handoff.HandoffSHA256, handoff.StartReceiptID, handoff.StartReceiptLocator, handoff.StartReceiptSHA256, handoff.TerminalResultID, handoff.TerminalResultLocator, handoff.TerminalResultSHA256, handoff.InspectionSHA256, created.Task.GitCommonDir, taskPath, resultOID, resultTree)
	draftPath := filepath.Join(root, "task-result.json")
	if err := os.WriteFile(draftPath, []byte(draft), 0o600); err != nil {
		t.Fatal(err)
	}
	reader := &fixedTaskEvidenceReader{evidence: evidence}
	base.dependencies.HandoffEvidence = reader
	base.dependencies.TaskLifecycleIDs = &sequenceTaskLifecycleIDs{}
	return taskDeliveryFixture{workItemJourneyFixture: base, root: root, taskPath: taskPath, draftPath: draftPath, resultOID: resultOID, resultTree: resultTree, reader: reader}
}

func writeHumanQADraft(t *testing.T, fixture taskDeliveryFixture, result TaskResultRecord, outcome, role string) (string, string) {
	t.Helper()
	evidencePath := filepath.Join(fixture.root, "human-qa-report.txt")
	evidenceBytes := []byte("fixture human QA observation\n")
	if err := os.WriteFile(evidencePath, evidenceBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	draft := fmt.Sprintf(`{"kind":"WorkspaceTaskHumanQARecordDraft@1","schema_version":1,"format":"json","format_version":1,"canonicalization":"RFC8785","publication_key":"task/qa-delivery","task_id":"task","task_result_id":%q,"result_oid":%q,"result_tree":%q,"outcome":%q,"actor":{"actor_claim":"fixture human tester, not actual approval","start_surface":"test CLI","started_at_utc":"2026-09-28T12:01:00Z","completed_at_utc":"2026-09-28T12:02:00Z"},"evidence":[{"id":"report","role":%q,"locator":%q,"sha256":%q,"size_bytes":%d}],"observation":"The product behavior is correct.","accepted_residual_risks":[]}`,
		result.ID, result.ResultOID, result.ResultTree, outcome, role, evidencePath, digestTaskBytes(evidenceBytes), len(evidenceBytes))
	draftPath := filepath.Join(fixture.root, "human-qa.json")
	if err := os.WriteFile(draftPath, []byte(draft), 0o600); err != nil {
		t.Fatal(err)
	}
	return draftPath, evidencePath
}

func TestTaskLifecycleIDsUseTypedRandomIdentitiesAndFailClosed(t *testing.T) {
	source := cryptoTaskLifecycleIDSource{Reader: strings.NewReader(strings.Repeat("a", 96))}
	got := []string{}
	result, err := source.NewTaskResultID()
	if err != nil {
		t.Fatal(err)
	}
	got = append(got, string(result))
	qa, _ := source.NewHumanQARecordID()
	authority, _ := source.NewIntegrationAuthorityID()
	intent, _ := source.NewIntegrationIntentID()
	attempt, _ := source.NewIntegrationAttemptID()
	integrationResult, _ := source.NewIntegrationResultID()
	got = append(got, string(qa), string(authority), string(intent), string(attempt), string(integrationResult))
	wantPrefixes := []string{"trs_", "hqa_", "iauth_", "iint_", "iat_", "ires_"}
	for i := range got {
		if !strings.HasPrefix(got[i], wantPrefixes[i]) || len(got[i]) != len(wantPrefixes[i])+32 {
			t.Fatalf("identity %d = %q", i, got[i])
		}
	}
	short := cryptoTaskLifecycleIDSource{Reader: strings.NewReader("short")}
	if _, err := short.NewTaskResultID(); err == nil || !hasWorkErrorClass(err, ErrorWorkIO) {
		t.Fatalf("short RNG error = %v", err)
	}
}

func TestTaskResultDraftStrictCanonicalSchema(t *testing.T) {
	directory := physicalPath(t, t.TempDir())
	path := filepath.Join(directory, "result.json")
	digest := "sha256:" + strings.Repeat("a", 64)
	oid := strings.Repeat("b", 40)
	tree := strings.Repeat("c", 40)
	common := filepath.Join(directory, "repo.git")
	worktree := filepath.Join(directory, "task")
	draft := `{"kind":"WorkspaceTaskResultRecordDraft@1","schema_version":1,"format":"json","format_version":1,"canonicalization":"RFC8785","publication_key":"task/result-1","task_id":"task","task_worktree_id":"wt_11111111111111111111111111111111","handoff":{"activity_id":"act_11111111111111111111111111111111","run_id":"run_11111111111111111111111111111111","handoff_id":"hnd_11111111111111111111111111111111","handoff_locator":"` + filepath.Join(directory, "handoff.json") + `","handoff_sha256":"` + digest + `","start_receipt_id":"rcp_11111111111111111111111111111111","start_receipt_locator":"` + filepath.Join(directory, "start.json") + `","start_receipt_sha256":"` + digest + `","terminal_result_id":"res_11111111111111111111111111111111","terminal_result_locator":"` + filepath.Join(directory, "terminal.json") + `","terminal_result_sha256":"` + digest + `","inspection_sha256":"` + digest + `"},"source":{"project_id":"ply","repo_id":"ply","git_common_dir":"` + common + `","worktree_locator":"` + worktree + `","source_ref":"refs/heads/task","result_oid":"` + oid + `","result_tree":"` + tree + `"},"technical_assessment":{"gate":"passed","required_verifier_ids":[],"accepted_debt":[]},"evidence_artifacts":[],"recorder":{"actor_claim":"local operator","control_surface":"CLI","recorded_at_utc":"2026-09-28T12:00:00Z"}}`
	if err := os.WriteFile(path, []byte(draft), 0o600); err != nil {
		t.Fatal(err)
	}
	files := projectCwdFileSystem{FileSystem: systemFileSystem{}, cwd: directory}
	parsed, err := readTaskResultDraft(files, path)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.ResultOID != oid || parsed.ResultTree != tree || parsed.TechnicalGate != "passed" || parsed.Digest != digestTaskBytes(parsed.Canonical) {
		t.Fatalf("parsed = %#v", parsed)
	}
	canonical := append([]byte(nil), parsed.Canonical...)
	if err := os.WriteFile(path, append([]byte(" \n"), []byte(draft)...), 0o600); err != nil {
		t.Fatal(err)
	}
	if spaced, err := readTaskResultDraft(files, path); err != nil || string(spaced.Canonical) != string(canonical) || spaced.Digest != parsed.Digest {
		t.Fatalf("canonical draft changed with whitespace: %#v, %v", spaced, err)
	}
	for name, invalid := range map[string]string{
		"unknown":            strings.Replace(draft, `"publication_key":`, `"unknown":true,"publication_key":`, 1),
		"duplicate":          strings.Replace(draft, `"publication_key":"task/result-1"`, `"publication_key":"task/result-1","publication_key":"task/result-1"`, 1),
		"missing":            strings.Replace(draft, `"technical_assessment":{"gate":"passed","required_verifier_ids":[],"accepted_debt":[]},`, "", 1),
		"wrong envelope":     strings.Replace(draft, "WorkspaceTaskResultRecordDraft@1", "Other@1", 1),
		"invalid gate":       strings.Replace(draft, `"gate":"passed"`, `"gate":"maybe"`, 1),
		"unsorted verifiers": strings.Replace(draft, `"required_verifier_ids":[]`, `"required_verifier_ids":["z","a"]`, 1),
		"non UTC recorder":   strings.Replace(draft, "2026-09-28T12:00:00Z", "2026-09-28T14:00:00+02:00", 1),
		"oversized text":     strings.Replace(draft, "local operator", strings.Repeat("x", 257), 1),
		"oversized artifact": strings.Replace(draft, `"evidence_artifacts":[]`, `"evidence_artifacts":[{"artifact_id":"large","role":"other","locator":"`+filepath.Join(directory, "large")+`","sha256":"`+digest+`","size_bytes":67108865}]`, 1),
	} {
		t.Run(name, func(t *testing.T) {
			if err := os.WriteFile(path, []byte(invalid), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := readTaskResultDraft(files, path); err == nil || !hasWorkErrorClass(err, ErrorTaskResultSchemaInvalid) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestTaskEvidenceGateRejectsMissingVerifierAndOpenFinding(t *testing.T) {
	draft := taskResultDraft{TechnicalGate: "passed", RequiredVerifierIDs: []string{"unit"}, EvidenceArtifacts: []TaskArtifactRecord{}}
	evidence := TaskHandoffEvidence{ReportedOutcome: "complete", SchemaValid: true, DigestValid: true, LifecycleValid: true, BindingValid: true, CapabilityValid: true, PrincipalSessionValid: true, PolicyValid: true, EvidenceCoverageValid: true, Artifacts: []TaskArtifactRecord{}}
	if err := validateTaskEvidence(draft, evidence); err == nil || !hasWorkErrorClass(err, ErrorTaskResultEvidenceConflict) {
		t.Fatalf("missing verifier error = %v", err)
	}
	evidence.VerifierResults = []TaskVerifierResultRecord{{VerifierID: "unit", Outcome: "passed"}}
	evidence.Review.OpenActionableFindings = []TaskReviewEntry{{ID: "finding", Severity: "high", Summary: "unsafe"}}
	if err := validateTaskEvidence(draft, evidence); err == nil || !hasWorkErrorClass(err, ErrorTaskResultEvidenceConflict) {
		t.Fatalf("open finding error = %v", err)
	}
}

func TestTaskEvidenceGateUnionAndDebtRiskBoundaries(t *testing.T) {
	base := TaskHandoffEvidence{
		ReportedOutcome: "complete", SchemaValid: true, DigestValid: true, LifecycleValid: true, BindingValid: true,
		CapabilityValid: true, PrincipalSessionValid: true, PolicyValid: true, EvidenceCoverageValid: true,
		ExpectedVerifierIDs: []string{}, VerifierResults: []TaskVerifierResultRecord{},
		Review:    TaskReviewRecord{Findings: []TaskReviewEntry{}, Fixes: []TaskReviewEntry{}, OpenActionableFindings: []TaskReviewEntry{}},
		Artifacts: []TaskArtifactRecord{}, EvidenceGaps: []string{},
	}
	for _, gate := range []string{"passed", "failed", "unknown"} {
		t.Run(gate, func(t *testing.T) {
			draft := taskResultDraft{TechnicalGate: gate, RequiredVerifierIDs: []string{}, EvidenceArtifacts: []TaskArtifactRecord{}, AcceptedDebt: []TaskAcceptedDebtRecord{}}
			if err := validateTaskEvidence(draft, base); err != nil {
				t.Fatalf("gate %s: %v", gate, err)
			}
		})
	}
	path := filepath.Join(physicalPath(t, t.TempDir()), "control")
	digest := "sha256:" + strings.Repeat("a", 64)
	goodEnough := taskResultDraft{TechnicalGate: "good_enough_with_known_debt", RequiredVerifierIDs: []string{}, EvidenceArtifacts: []TaskArtifactRecord{{ArtifactID: "control", Role: "debt_control", Locator: path, SHA256: digest, SizeBytes: 1}}, AcceptedDebt: []TaskAcceptedDebtRecord{{ID: "debt", RiskClass: "evidence", Severity: "low", Summary: "bounded gap", Control: "tracked", EvidenceArtifactIDs: []string{"control"}}}}
	debtEvidence := base
	debtEvidence.EvidenceCoverageValid = false
	debtEvidence.EvidenceGaps = []string{"missing detail"}
	debtEvidence.Artifacts = append([]TaskArtifactRecord(nil), goodEnough.EvidenceArtifacts...)
	if err := validateTaskEvidence(goodEnough, debtEvidence); err != nil {
		t.Fatalf("good-enough gate: %v", err)
	}
	for _, severity := range []string{"high", "critical"} {
		risky := debtEvidence
		risky.Review.Findings = []TaskReviewEntry{{ID: "risk", Severity: severity, Summary: "unsafe", EvidenceArtifactIDs: []string{}}}
		if err := validateTaskEvidence(goodEnough, risky); err == nil || !hasWorkErrorClass(err, ErrorTaskResultEvidenceConflict) {
			t.Fatalf("%s risk error = %v", severity, err)
		}
	}
}

func TestTaskTextLimitsCountRunesOnce(t *testing.T) {
	if !validTaskText(strings.Repeat("x", 256), 1, 256) {
		t.Fatal("256-rune bounded text was rejected")
	}
	if validTaskText(strings.Repeat("x", 257), 1, 256) {
		t.Fatal("257-rune bounded text was accepted")
	}
}

func TestTaskEvidenceGateRequiresExactArtifactRole(t *testing.T) {
	path := filepath.Join(physicalPath(t, t.TempDir()), "control.txt")
	digest := "sha256:" + strings.Repeat("a", 64)
	draft := taskResultDraft{
		TechnicalGate:       "good_enough_with_known_debt",
		RequiredVerifierIDs: []string{},
		AcceptedDebt: []TaskAcceptedDebtRecord{{
			ID: "debt", RiskClass: "evidence", Severity: "low", Summary: "known gap", Control: "tracked",
			EvidenceArtifactIDs: []string{"control"},
		}},
		EvidenceArtifacts: []TaskArtifactRecord{{ArtifactID: "control", Role: "debt_control", Locator: path, SHA256: digest, SizeBytes: 1}},
	}
	evidence := TaskHandoffEvidence{
		ReportedOutcome: "complete", SchemaValid: true, DigestValid: true, LifecycleValid: true,
		BindingValid: true, CapabilityValid: true, PrincipalSessionValid: true, PolicyValid: true,
		EvidenceCoverageValid: false, ExpectedVerifierIDs: []string{}, VerifierResults: []TaskVerifierResultRecord{},
		Review:       TaskReviewRecord{Findings: []TaskReviewEntry{}, Fixes: []TaskReviewEntry{}, OpenActionableFindings: []TaskReviewEntry{}},
		EvidenceGaps: []string{"missing journal"},
		Artifacts:    []TaskArtifactRecord{{ArtifactID: "control", Role: "review", Locator: path, SHA256: digest, SizeBytes: 1}},
	}
	if err := validateTaskEvidence(draft, evidence); err == nil || !hasWorkErrorClass(err, ErrorTaskResultEvidenceConflict) {
		t.Fatalf("artifact role mismatch error = %v", err)
	}
}

func TestHumanQADraftUsesHumanQASchemaErrorForFileShape(t *testing.T) {
	directory := physicalPath(t, t.TempDir())
	target := filepath.Join(directory, "target.json")
	if err := os.WriteFile(target, []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(directory, "qa.json")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	files := projectCwdFileSystem{FileSystem: systemFileSystem{}, cwd: directory}
	if _, err := readTaskHumanQADraft(files, link); err == nil || !hasWorkErrorClass(err, ErrorTaskQASchemaInvalid) {
		t.Fatalf("human QA symlink error = %v", err)
	}
}

func TestTaskResultRecordMigratesFormatOneAtomicallyAndRetriesWithoutRewrite(t *testing.T) {
	fixture := newTaskDeliveryFixture(t, true)
	result, err := RecordTaskResult(fixture.dependencies, TaskResultRecordInput{TaskID: "task", File: fixture.draftPath})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Created || result.Record.StoreTransition != "format_1_to_2" || result.Record.ResultOID != fixture.resultOID || fixture.reader.reads != 1 {
		t.Fatalf("result = %#v, evidence reads=%d", result, fixture.reader.reads)
	}
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil || registry.FormatVersion != 2 || len(registry.TaskResults) != 1 || registry.HumanQARecords == nil || registry.IntegrationAuthorities == nil || registry.IntegrationIntents == nil || registry.IntegrationAttempts == nil || registry.IntegrationResults == nil {
		t.Fatalf("migrated registry = %#v, %v", registry, err)
	}
	before, err := os.ReadFile(workItemsPath(fixture.root))
	if err != nil {
		t.Fatal(err)
	}
	beforeInfo, err := os.Stat(workItemsPath(fixture.root))
	if err != nil {
		t.Fatal(err)
	}
	retry, err := RecordTaskResult(fixture.dependencies, TaskResultRecordInput{TaskID: "task", File: fixture.draftPath})
	if err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(workItemsPath(fixture.root))
	afterInfo, _ := os.Stat(workItemsPath(fixture.root))
	if retry.Created || retry.Record.ID != result.Record.ID || string(after) != string(before) || !afterInfo.ModTime().Equal(beforeInfo.ModTime()) {
		t.Fatalf("idempotent retry changed state: retry=%#v bytes_equal=%t mtime_equal=%t", retry, string(after) == string(before), afterInfo.ModTime().Equal(beforeInfo.ModTime()))
	}
}

func TestHumanQADraftStrictSchemaAndOutcomeUnion(t *testing.T) {
	fixture := newTaskDeliveryFixture(t, false)
	result, err := RecordTaskResult(fixture.dependencies, TaskResultRecordInput{TaskID: "task", File: fixture.draftPath})
	if err != nil {
		t.Fatal(err)
	}
	path, _ := writeHumanQADraft(t, fixture, result.Record, "pass", "report")
	files := fixture.dependencies.Files
	validBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if parsed, err := readTaskHumanQADraft(files, path); err != nil || parsed.Outcome != "pass" || parsed.Digest != digestTaskBytes(parsed.Canonical) {
		t.Fatalf("parsed=%#v error=%v", parsed, err)
	}
	valid := string(validBytes)
	for name, invalid := range map[string]string{
		"unknown":        strings.Replace(valid, `"publication_key":`, `"unknown":true,"publication_key":`, 1),
		"duplicate":      strings.Replace(valid, `"publication_key":"task/qa-delivery"`, `"publication_key":"task/qa-delivery","publication_key":"task/qa-delivery"`, 1),
		"missing":        strings.Replace(valid, `,"observation":"The product behavior is correct."`, "", 1),
		"invalid union":  strings.Replace(valid, `"outcome":"pass"`, `"outcome":"maybe"`, 1),
		"pass no report": strings.Replace(valid, `"role":"report"`, `"role":"log"`, 1),
	} {
		t.Run(name, func(t *testing.T) {
			if err := os.WriteFile(path, []byte(invalid), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := readTaskHumanQADraft(files, path); err == nil || !hasWorkErrorClass(err, ErrorTaskQASchemaInvalid) {
				t.Fatalf("error = %v", err)
			}
		})
	}
	for _, outcome := range []string{"fail", "blocked"} {
		candidate := strings.Replace(valid, `"outcome":"pass"`, `"outcome":"`+outcome+`"`, 1)
		candidate = strings.Replace(candidate, `"role":"report"`, `"role":"log"`, 1)
		if err := os.WriteFile(path, []byte(candidate), 0o600); err != nil {
			t.Fatal(err)
		}
		if parsed, err := readTaskHumanQADraft(files, path); err != nil || parsed.Outcome != outcome {
			t.Fatalf("%s parsed=%#v error=%v", outcome, parsed, err)
		}
	}
}

func TestHumanQARecordRehashesAndBindsExactTaskResult(t *testing.T) {
	fixture := newTaskDeliveryFixture(t, false)
	result, err := RecordTaskResult(fixture.dependencies, TaskResultRecordInput{TaskID: "task", File: fixture.draftPath})
	if err != nil {
		t.Fatal(err)
	}
	path, evidencePath := writeHumanQADraft(t, fixture, result.Record, "pass", "report")
	recorded, err := RecordTaskHumanQA(fixture.dependencies, TaskHumanQARecordInput{TaskID: "task", File: path})
	if err != nil || !recorded.Created || recorded.Record.TaskResultID != result.Record.ID || recorded.Record.Outcome != "pass" {
		t.Fatalf("recorded=%#v error=%v", recorded, err)
	}
	before, _ := os.ReadFile(workItemsPath(fixture.root))
	retry, err := RecordTaskHumanQA(fixture.dependencies, TaskHumanQARecordInput{TaskID: "task", File: path})
	after, _ := os.ReadFile(workItemsPath(fixture.root))
	if err != nil || retry.Created || retry.Record.ID != recorded.Record.ID || string(before) != string(after) {
		t.Fatalf("retry=%#v error=%v bytes_equal=%t", retry, err, string(before) == string(after))
	}

	t.Run("artifact rehash", func(t *testing.T) {
		other := newTaskDeliveryFixture(t, false)
		otherResult, err := RecordTaskResult(other.dependencies, TaskResultRecordInput{TaskID: "task", File: other.draftPath})
		if err != nil {
			t.Fatal(err)
		}
		otherDraft, otherEvidence := writeHumanQADraft(t, other, otherResult.Record, "pass", "report")
		if err := os.WriteFile(otherEvidence, []byte("changed\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := RecordTaskHumanQA(other.dependencies, TaskHumanQARecordInput{TaskID: "task", File: otherDraft}); err == nil || !hasWorkErrorClass(err, ErrorTaskQAEvidenceConflict) {
			t.Fatalf("rehash error = %v", err)
		}
	})

	t.Run("wrong result tree", func(t *testing.T) {
		other := newTaskDeliveryFixture(t, false)
		otherResult, err := RecordTaskResult(other.dependencies, TaskResultRecordInput{TaskID: "task", File: other.draftPath})
		if err != nil {
			t.Fatal(err)
		}
		otherDraft, _ := writeHumanQADraft(t, other, otherResult.Record, "pass", "report")
		contents, _ := os.ReadFile(otherDraft)
		wrong := strings.Replace(string(contents), otherResult.Record.ResultTree, strings.Repeat("f", len(otherResult.Record.ResultTree)), 1)
		if err := os.WriteFile(otherDraft, []byte(wrong), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := RecordTaskHumanQA(other.dependencies, TaskHumanQARecordInput{TaskID: "task", File: otherDraft}); err == nil || !hasWorkErrorClass(err, ErrorTaskQAEvidenceConflict) {
			t.Fatalf("binding error = %v", err)
		}
	})

	if _, err := os.Stat(evidencePath); err != nil {
		t.Fatal(err)
	}
}
