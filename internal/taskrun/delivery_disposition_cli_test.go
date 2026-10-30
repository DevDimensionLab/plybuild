package taskrun

import (
	"encoding/json"

	"fmt"

	"os"

	"path/filepath"
	"strings"
	"testing"
)

// DEP02: modeled provider/Herdr, real native CLI/store/Git. Historical command
// invocation and child/wait receipts deliberately do not exist in this fixture.
func dep02Request(t *testing.T, f deliveryFixture, o WorkflowRun, event string) string {
	t.Helper()
	target := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	proof := writeAny(t, f.R.WorkspaceRoot, "dep02-assessment.json", map[string]string{
		"provenance":   "modeled owner report and static operation inspection; not runtime absence",
		"residual":     "Unknown fixture child/wait/descendants, temp/socket/cleanup, exact opened runtime inputs, exhaustive filesystem writes and pending writes; browser launch refused",
		"independence": "New native source snapshot binds current dirty inputs; fresh clean candidate, independent review, isolated acceptance input snapshots and exact target guards bound later actions. This is no claim of historical write/child absence",
	})
	return writeAny(t, f.R.WorkspaceRoot, event+".json", map[string]any{
		"kind": "ply.workflow.scoped-disposition", "schema_version": 1, "event_id": event,
		"run_id": o.RunID, "request_sha256": o.RequestSHA256, "mandate_sha256": o.Handoff.SHA256,
		"native_session": o.Transport.AgentSessionID,
		"source_oid":     gitOutput(t, target, "rev-parse", "HEAD"), "source_tree": gitOutput(t, target, "rev-parse", "HEAD^{tree}"),
		"retained":   []any{map[string]any{"stop": *o.Delivery.ActiveStop, "unknowns": []string{"child_wait_and_descendants", "temp_socket_and_cleanup", "runtime_input_bytes", "filesystem_writes", "pending_writes"}, "conflicts": []string{}, "provenance": "owner_and_static_assessment", "assessment": "Historical fixture effects remain unknown; missing invocation/child receipts are not reconstructed", "evidence": []FileBinding{{proof, hashFileTest(t, proof)}}}},
		"operations": []any{map[string]any{"operation": "prepare", "boundary": "fresh_source_snapshot", "provenance": "owner_and_static_assessment", "action": "Prepare the authorized file-only candidate without using the historical fixture", "reads": []string{"path:" + target}, "writes": []string{"path:" + target, "path:" + f.AcceptancePath, "native:git-index", "native:ignored-build-output", "native:verification-scratch"}, "children": []string{"process:node-check", "process:node-build", "process:node-precache", "process:git-commit"}, "assessment": "This decision accepts retained unlocated historical writes and child uncertainty for a newly observed source snapshot. Known concrete input/write/child overlaps block the operation; subsequent source/candidate/input drift invalidates it", "inputs": []FileBinding{{proof, hashFileTest(t, proof)}}}},
	})
}

func dep02Stop(t *testing.T, binary string, f deliveryFixture, o WorkflowRun) WorkflowRun {
	t.Helper()
	proof := writeAny(t, f.R.WorkspaceRoot, "dep02-denial.json", map[string]string{"source": "modeled refusal after fixture start; invocation and child wait unknown"})
	return automaticCLIOK(t, binary, f, "execute", "observe-stop", o.RunID, "--event-id", "dep02-denied", "--provider", "codex", "--native-session", o.Transport.AgentSessionID, "--stage", "work", "--action", "Launch historical browser fixture", "--effect-status", DeliveryEffectPartial, "--observer", "modeled coordinator", "--reason", "Modeled browser refused after fixture launch; old child outcome unknown", "--evidence", proof)
}

func dep02CandidateRequest(t *testing.T, f deliveryFixture, o WorkflowRun, event string, binary string) string {
	t.Helper()
	path := dep02Request(t, f, o, event)
	var in map[string]any
	raw, e := os.ReadFile(path)
	if e != nil {
		t.Fatal(e)
	}
	if e = json.Unmarshal(raw, &in); e != nil {
		t.Fatal(e)
	}
	first := in["operations"].([]any)[0].(map[string]any)
	ops := []any{}
	for _, operation := range []string{"report", "review", "verify", "integrate"} {
		op := map[string]any{}
		for k, v := range first {
			op[k] = v
		}
		op["operation"] = operation
		op["boundary"] = map[string]string{"report": "fresh_source_snapshot", "review": "exact_clean_candidate", "verify": "candidate_and_isolated_verifier", "integrate": "qualified_candidate_and_exact_target"}[operation]
		op["action"] = "Perform independent native " + operation + " under the original local agreement"
		op["writes"] = []string{}
		op["children"] = []string{}
		if operation == "verify" {
			op["writes"] = []string{"native:verification-scratch", "native:run-evidence"}
			op["children"] = []string{"process:node-check", "process:node-build", "process:node-precache"}
			op["inputs"] = append([]any{map[string]string{"locator": f.AcceptancePath, "sha256": hashFileTest(t, f.AcceptancePath)}, map[string]string{"locator": binary, "sha256": hashFileTest(t, binary)}}, first["inputs"].([]any)...)
		}
		if operation == "integrate" {
			op["writes"] = []string{"native:local-target", "native:run-evidence"}
		}
		ops = append(ops, op)
	}
	if log := filepath.Join(f.R.WorkspaceRoot, "node-preparation.log"); dep02FileExists(log) {
		for _, value := range ops {
			op := value.(map[string]any)
			op["inputs"] = append(op["inputs"].([]any), map[string]string{"locator": log, "sha256": hashFileTest(t, log)})
		}
	}
	in["operations"] = ops
	return writeAny(t, f.R.WorkspaceRoot, event+"-candidate.json", in)
}

func dep02ProcessModel(t *testing.T, f deliveryFixture) {
	t.Helper()
	script := `#!/usr/bin/python3
import os
# The fake lsof process is a child of this CLI; only owned PIDs are read.
print("p%d\nR%d\nn/modeled-cwd" % (os.getppid(), ` + fmt.Sprint(os.Getppid()) + `))
`
	if e := os.WriteFile(filepath.Join(f.R.WorkspaceRoot, "bin", "lsof"), []byte(script), 0700); e != nil {
		t.Fatal(e)
	}
}

func dep02Evidence(t *testing.T, name string, value any) {
	t.Helper()
	if dir := os.Getenv("PLY_DEP02_EVIDENCE_DIR"); dir != "" {
		raw, e := json.MarshalIndent(value, "", "  ")
		if e != nil {
			t.Fatal(e)
		}
		if e = os.MkdirAll(dir, 0700); e != nil {
			t.Fatal(e)
		}
		file := filepath.Join(dir, strings.ReplaceAll(t.Name(), "/", "_")+"-"+name+".json")
		if e = os.WriteFile(file, raw, 0600); e != nil {
			t.Fatal(e)
		}
	}
}

func dep02Read(t *testing.T, path string) []byte {
	t.Helper()
	raw, e := os.ReadFile(path)
	if e != nil {
		t.Fatal(e)
	}
	return raw
}
func dep02FileExists(path string) bool { _, e := os.Stat(path); return e == nil }

func dep02ModelFiles(t *testing.T, target string, complete bool) {
	t.Helper()
	files := map[string]string{
		".gitignore":                 "dist/\ndist-server/\n.cache/\n",
		"package.json":               `{"private":true,"scripts":{"check":"node scripts/model-artifact.mjs check","build":"node scripts/model-artifact.mjs build","test:build":"node scripts/model-artifact.mjs precache"}}`,
		"scripts/model-artifact.mjs": "throw new Error('incomplete inherited validator');\n",
	}
	if complete {
		files["scripts/model-artifact.mjs"] = `import assert from 'node:assert/strict';
import {mkdirSync,readFileSync,writeFileSync} from 'node:fs';
import {join} from 'node:path';
import {spawnSync} from 'node:child_process';
const expected='independent expected';
assert.equal(readFileSync('README.md','utf8').trim(),expected);
const task=process.argv[2];
mkdirSync('.cache',{recursive:true});
writeFileSync(join(process.env.TMPDIR,'node-'+task),task);
if(task==='check') {
 const child=spawnSync(process.execPath,['-e',"require('node:fs').writeFileSync('.cache/child-wait.json',JSON.stringify({exit:0,waited:true}))"],{stdio:'pipe'});
 assert.equal(child.status,0);assert.equal(child.signal,null);
 assert.equal(JSON.parse(readFileSync('.cache/child-wait.json','utf8')).waited,true);
} else if(task==='build') {
 for(const dir of ['dist','dist-server']) {mkdirSync(dir,{recursive:true});writeFileSync(join(dir,'artifact.json'),JSON.stringify({candidate:expected,modeled:true}));}
} else if(task==='precache') {
 for(const dir of ['dist','dist-server']) assert.equal(JSON.parse(readFileSync(join(dir,'artifact.json'),'utf8')).candidate,expected);
 writeFileSync('.cache/precache.json',JSON.stringify({validated:['dist','dist-server']}));
} else throw new Error('unexpected modeled operation');
console.log('modeled '+task+' actual child/write checks passed; historical writes and descendants remain unknown');
`
	}
	for rel, contents := range files {
		path := filepath.Join(target, rel)
		if e := os.MkdirAll(filepath.Dir(path), 0700); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(path, []byte(contents), 0600); e != nil {
			t.Fatal(e)
		}
	}
}

// The existing observation schema has no sandbox field. A mandate receipt
// must validate against its preserved matches_contract=null native acceptance,
// not fabricate a new sandbox attestation during compatible continuation.
func TestDEP02ReceiptRuntimeObservationPreservesUnattestedSandbox(t *testing.T) {
	f := receiptFixture(t, "codex")
	o := workflowTestStart(t, f.workflowFixture)
	o, e := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}, true)
	if e != nil {
		t.Fatal(e)
	}
	path := deliveryTestRuntimeObservation(t, f, o)
	state, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if e = deliveryRuntimeObservation(state, FileBinding{path, hashFileTest(t, path)}); e != nil {
		t.Fatalf("same receipt observation incorrectly requires sandbox attestation: %v", e)
	}
}
