#!/usr/bin/env bash

set -euo pipefail
export LC_ALL=C LANG=C GOENV=off GOWORK=off GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
temp_root=$(mktemp -d "${TMPDIR:-/tmp}/ply-task-integration.XXXXXX")
temp_root=$(cd "$temp_root" && pwd -P)
cleanup() { find "$temp_root" -type d -exec chmod u+rwx {} + 2>/dev/null || true; rm -rf -- "$temp_root"; }
trap cleanup EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM
mkdir -p "$temp_root/home" "$temp_root/gocache"
binary="$temp_root/ply"
go_bin=${GO:-go}

fail() { printf 'workspace task integration roundtrip: %s\n' "$*" >&2; exit 1; }

(
	cd "$repo_root"
	GOCACHE="$temp_root/gocache" "$go_bin" build -o "$binary" ./cmd/ply
) || fail 'could not build fresh binary'

check_help() {
	local label=$1
	shift
	HOME="$temp_root/home" "$binary" "$@" >"$temp_root/$label.stdout" 2>"$temp_root/$label.stderr" ||
		fail "$label returned non-zero"
	[[ ! -s "$temp_root/$label.stderr" ]] || fail "$label wrote to stderr"
	grep -F 'Usage:' "$temp_root/$label.stdout" >/dev/null || fail "$label has no usage block"
}

check_help root --help
check_help workspace workspace --help
check_help task workspace task --help
check_help result workspace task result --help
check_help result-record workspace task result record --help
check_help qa workspace task qa --help
check_help qa-record workspace task qa record --help
check_help integrate workspace task integrate --help
check_help show workspace task show --help

grep -F 'Record immutable delivery evidence as typed results for workspace Tasks.' "$temp_root/result.stdout" >/dev/null || fail 'Task result help changed'
grep -F 'Record a human product QA outcome for an exact controlled Task result.' "$temp_root/qa.stdout" >/dev/null || fail 'Task QA help changed'
grep -F 'Check or apply one confirmed local fast-forward from an exact Task result to its registered Epic parent.' "$temp_root/integrate.stdout" >/dev/null || fail 'Task integration help changed'

file_sha256() {
	if command -v shasum >/dev/null 2>&1; then shasum -a 256 "$1" | awk '{print "sha256:"$1}'; else sha256sum "$1" | awk '{print "sha256:"$1}'; fi
}

file_mtime() {
	if [[ $(uname -s) == Darwin ]]; then stat -f '%m' "$1"; else stat -c '%Y' "$1"; fi
}

helper_source="$temp_root/lifecycle_helper.go"
cat >"$helper_source" <<'GO'
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

func must(err error) { if err != nil { panic(err) } }
func read(path string) map[string]any { b,err:=os.ReadFile(path);must(err);var v map[string]any;must(json.Unmarshal(b,&v));return v }
func emit(path string,v any){b,err:=json.Marshal(v);must(err);must(os.WriteFile(path,b,0600))}
func digestBytes(b []byte)string{s:=sha256.Sum256(b);return "sha256:"+hex.EncodeToString(s[:])}
func digestFile(path string)(string,int64){b,err:=os.ReadFile(path);must(err);return digestBytes(b),int64(len(b))}
func digestJSON(path string)string{v:=read(path);b,err:=json.Marshal(v);must(err);return digestBytes(b)}
func stringsAny(values []string)[]any{r:=make([]any,len(values));for i,v:=range values{r[i]=v};return r}
func envelope(kind string)map[string]any{return map[string]any{"kind":kind,"schema_version":1,"format":"json","format_version":1,"canonicalization":"RFC8785"}}
func field(v any,path string)any{for _,part:=range strings.Split(path,"."){switch x:=v.(type){case map[string]any:v=x[part];case []any:var i int;_,e:=fmt.Sscan(part,&i);must(e);v=x[i]}};return v}
func main(){
	if len(os.Args)<2{panic("mode required")}
	switch os.Args[1]{
	case "field": fmt.Print(field(read(os.Args[2]),os.Args[3]))

	case "spec":
		out,showPath,root,revision,oldPath:=os.Args[2],os.Args[3],os.Args[4],os.Args[5],os.Args[6];show:=read(showPath);resource:=show["resource"].(map[string]any);task:=resource["persisted"].(map[string]any);wt:=task["worktree"].(map[string]any);observed:=resource["observed"].(map[string]any)["target"].(map[string]any);v:=taskDraft("WorkspaceTaskSpecDraft@1","fixture/spec-"+revision)
		v["spec_id"],v["expected_previous"],v["problem"],v["title"]="solution",nil,show["problem"].(map[string]any)["head"],"Fixture solution"
		var previous map[string]any;if oldPath!="-"{previous=read(oldPath)["task_spec_binding"].(map[string]any)["spec"].(map[string]any);v["expected_previous"]=previous}
		docs:=[]any{};parts:=map[string]any{};refs:=map[string][]any{}
		for _,name:=range []string{"abstract","functional","technical"}{path:=root+"/"+name+".md";var source map[string]any
			if previous!=nil&&name!="technical"{source=map[string]any{"kind":"snapshot","manifest_sha256":previous["manifest_sha256"],"document_id":name}}else{must(os.WriteFile(path,[]byte("# Fixture "+name+" revision "+revision+"\n\nPreserve the fixture behavior. Scope: delivery.txt. Verify locally. Retain evidence on failure.\n"),0600));hash,size:=digestFile(path);source=map[string]any{"kind":"file","locator":path,"sha256":hash,"size_bytes":size,"media_type":"text/markdown","git_provenance":nil}}
			docs=append(docs,map[string]any{"id":name,"source":source});refs[name]=[]any{map[string]any{"document_id":name,"section":nil}};parts[name]=map[string]any{"state":"present","reason":nil,"documents":refs[name]}}
		v["parts"],v["documents"],v["supporting"]=parts,docs,[]any{};v["requirements"]=[]any{map[string]any{"id":"f-01","functional_refs":refs["functional"],"acceptance":"The isolated fixture preserves its exact delivery evidence.","verification_ids":[]any{"tests"},"technical_refs":refs["technical"]}};v["removed_requirement_ids"]=[]any{};v["phases"]=[]any{map[string]any{"id":"verify","purpose":"Verify the isolated fixture.","requirement_ids":[]any{"f-01"},"entry_criteria":[]any{},"exit_criteria":[]any{"Exact fixture evidence is preserved."},"verification_ids":[]any{"tests"}}};v["implementation_basis"]=map[string]any{"project_id":"ply","repo_id":"ply","git_common_dir":task["git_common_dir"],"epic_id":"epic","parent_worktree_id":wt["parent_worktree_id"],"parent_ref":wt["parent_ref"],"parent_oid":wt["parent_oid"],"parent_tree":wt["parent_tree"],"start_oid":observed["oid"],"start_tree":observed["tree"]};v["dependencies"]=[]any{};v["change_reason"]="Fixture revision "+revision;emit(out,v)
	case "assessment":
		v:=taskDraft("WorkspaceTaskSpecAssessmentDraft@1","fixture/assessment");v["spec_id"],v["spec"],v["expected_previous_assessment"]="solution",revisionRef(read(os.Args[3])["outcome_ref"].(map[string]any)),nil;v["outcome"],v["reason"],v["open_questions"],v["documents"]="ready","Fixture readiness claim.",[]any{},[]any{};checks:=[]any{};for _,id:=range []string{"acceptance_coverage","implementation_basis","problem_coverage","recovery","scope_and_phases","three_parts"}{checks=append(checks,map[string]any{"id":id,"outcome":"pass","reason":"Fixture assessment, not actual human approval.","evidence_document_ids":[]any{}})};v["checks"]=checks;emit(os.Args[2],v)
	case "selection":
		v:=taskDraft("WorkspaceTaskSolutionSelectionDraft@1","fixture/selection");v["expected_previous_selection"],v["action"],v["reason"]=nil,"select","Fixture only: choose revision 1.";v["solution"]=map[string]any{"spec_id":"solution","spec":revisionRef(read(os.Args[3])["outcome_ref"].(map[string]any)),"problem":read(os.Args[5])["problem"].(map[string]any)["head"],"assessment":decisionRef(read(os.Args[4])["outcome_ref"].(map[string]any))};v["human_decision"]=map[string]any{"actor_claim":"fixture human, not actual approval","decided_at_utc":"2026-09-29T12:00:00Z","source":"explicit_human_instruction","statement":"Fixture only: select solution revision 1."};emit(os.Args[2],v)
	case "problem2":
		v:=taskDraft("WorkspaceTaskProblemDraft@1","fixture/problem-2");v["expected_previous"]=read(os.Args[3])["problem"].(map[string]any)["head"];v["origin"]=map[string]any{"kind":"authored"};v["title"],v["summary"],v["problem_document_id"],v["change_reason"]="Changed fixture problem","Fixture changed after the historical start.","problem","Fixture change.";path:=os.Args[2]+".md";must(os.WriteFile(path,[]byte("# Changed fixture problem\n\nNeed a new selected basis.\n"),0600));hash,size:=digestFile(path);v["documents"]=[]any{map[string]any{"id":"problem","source":map[string]any{"kind":"file","locator":path,"sha256":hash,"size_bytes":size,"media_type":"text/markdown","git_provenance":nil}}};v["sources"],v["claims"],v["deadline"]=[]any{},[]any{},nil;emit(os.Args[2],v)
	case "other-handoff":v:=read(os.Args[3]);v["publication_key"],v["activity_key"]="fixture/negative-start","fixture/negative-start";emit(os.Args[2],v)
	case "invalid-handoff":
		v:=read(os.Args[3]);mode:=os.Args[4];v["publication_key"],v["activity_key"]="fixture/invalid-"+mode,"fixture/invalid-"+mode
		switch mode { case "v1":v["schema_version"]=1;delete(v,"task_spec_binding");case "null":v["task_spec_binding"]=nil;case "wrong-worktree":v["task_spec_binding"].(map[string]any)["task_worktree_id"]="wt_00000000000000000000000000000000";case "missing-input":inputs:=[]any{};for _,raw:=range v["inputs"].([]any){if raw.(map[string]any)["id"]!="task-basis/problem"{inputs=append(inputs,raw)}};v["inputs"]=inputs;case "forged-input":for _,raw:=range v["inputs"].([]any){m:=raw.(map[string]any);if m["id"]=="task-basis/problem"{m["role"]="design"}} }
		emit(os.Args[2],v)
	case "negative-start":v:=read(os.Args[3]);v["acceptance"]="conflict";v["issues"]=[]any{map[string]any{"type":"task_spec","detail":"Fixture problem changed; the preserved selection is stale."}};emit(os.Args[2],v)
	case "requirements":
		h:=read(os.Args[3]);b:=h["task_spec_binding"].(map[string]any);target:=h["target_binding"].(map[string]any);v:=envelope("WorkspaceTaskRequirementEvidence@1");v["task_id"],v["spec_id"],v["spec"],v["result_oid"],v["result_tree"]=b["task_id"],b["spec_id"],b["spec"],target["oid"],target["tree"];v["requirements"]=[]any{map[string]any{"id":"f-01","outcome":"passed","verifier_ids":[]any{"tests"},"artifact_ids":[]any{"verifier-stdout"},"reason":"Isolated fixture evidence only; not actual human QA."}};emit(os.Args[2],v)
	case "handoff":
		out,target,ref,oid,input:=os.Args[2],os.Args[3],os.Args[4],os.Args[5],os.Args[6];inputDigest,inputSize:=digestFile(input);v:=envelope("ply.workflow.handoff-draft")
		v["publication_key"],v["activity_key"]="ws05/blackbox/happy","ws05/blackbox/happy"
		v["goal"]=map[string]any{"title":"Workspace Task integration blackbox","recipient_role":"Delivery agent","objective":"Preserve exact delivery evidence.","done_when":"The verifier passes and the exact result is reported."}
		v["recipient"]=map[string]any{"principal_id":"codex-delivery-agent","principal_kind":"human_started_agent","runtime_constraints":[]any{"local"}}
		v["binding_request"]=map[string]any{"project_id":"ply","repo_id":"ply","target_worktree":target,"target_ref":ref,"expected_oid":oid,"status_policy":map[string]any{"mode":"clean"}}
		v["inputs"]=[]any{map[string]any{"id":"fixture","role":"verification_fixture","locator":input,"sha256":inputDigest,"size_bytes":inputSize,"media_type":"text/plain","git_binding":nil}}
		v["authority"]=map[string]any{"allowed_effects":[]any{map[string]any{"id":"write-result","type":"filesystem_write","scope":map[string]any{"kind":"filesystem","paths":[]any{"delivery.txt"},"directory_prefixes":[]any{}},"max_occurrences":1,"sequence":1},map[string]any{"id":"run-tests","type":"command_execute","scope":map[string]any{"kind":"command","procedure_ids":[]any{"verify"},"verifier_ids":[]any{"tests"}},"max_occurrences":1,"sequence":2}},"forbidden_effects":[]any{map[string]any{"type":"network","reason":"Network is forbidden."},map[string]any{"type":"push","reason":"Push is forbidden."}},"human_gates":[]any{"human_task_qa","local_integration"}}
		v["budget"]=map[string]any{"max_rounds":2,"round_definition":map[string]any{"unit":"implementation_or_review_fix_iteration","command_retry_consumes_round":false,"retry_condition":"only_if_no_effect_started"}}
		v["procedure"]=[]any{map[string]any{"id":"verify","instruction":"Verify the exact local result. Report the mandatory managed application/json task-requirements artifact covering every selected requirement.","required_before":[]any{}}}
		v["verifiers"]=[]any{map[string]any{"id":"tests","argv":[]any{"git","status","--short"},"cwd":target,"env":[]any{},"expected_exit":0,"stop_on_failure":true,"evidence":map[string]any{"capture_stdout":true,"capture_stderr":false,"classification":"workspace_internal","binding":"target_oid"}}}
		stops:=[]string{"product_decision_required","scope_or_authority_expansion","target_or_input_drift","unknown_or_partial_effect","unexpected_sensitive_data","round_budget_exhausted"};values:=[]any{};for _,s:=range stops{values=append(values,map[string]any{"type":s,"description":"Stop safely when this condition occurs."})};v["stop_conditions"]=values
		v["reporting"]=map[string]any{"summary_max_codepoints":240,"meaning_max_codepoints":600,"required_start_fields":stringsAny([]string{"acceptance","binding","contract_digests","issues","observed_inputs","observed_project","observed_target","observed_workspace","principal","receipt_id","sandbox"}),"required_terminal_fields":stringsAny([]string{"artifacts","binding","evidence_gaps","final_target","forbidden_effects_observed","meaning","observed_effects","principal","reported_outcome","result_id","review","rounds_used","start_binding","stop_reasons","summary","verifier_results"})};selected:=read(os.Args[7]);v["schema_version"]=2;v["task_spec_binding"]=selected["task_spec_binding"];inputs:=v["inputs"].([]any);inputs=append(inputs,selected["required_inputs"].([]any)...);sort.Slice(inputs,func(i,j int)bool{return inputs[i].(map[string]any)["id"].(string)<inputs[j].(map[string]any)["id"].(string)});v["inputs"]=inputs;emit(out,v)
	case "start":
		out,handoffPath,tempRoot:=os.Args[2],os.Args[3],os.Args[4];h:=read(handoffPath);identity:=h["identity"].(map[string]any);workspace:=h["workspace_binding"].(map[string]any);project:=h["project_binding"].(map[string]any);target:=h["target_binding"].(map[string]any);reply:=h["reply_capability"].(map[string]any);principal:=map[string]any{"expected_principal_id":"codex-delivery-agent","human_start_principal":"fixture human, not actual approval","start_surface":"local shell","session_id":"blackbox-session","runtime_id":"codex","model_id":"test-model","started_at_utc":"2026-09-28T12:00:00Z"};v:=envelope("ply.workflow.start-receipt-draft")
		v["receipt_id"]=identity["start_receipt_id"];v["binding"]=map[string]any{"activity_id":identity["activity_id"],"run_id":identity["run_id"],"handoff_id":identity["handoff_id"],"handoff_sha256":digestJSON(handoffPath)};v["principal"]=principal;v["observed_workspace"]=map[string]any{"root":workspace["root"],"marker_format_version":workspace["marker_format_version"],"marker_sha256":workspace["marker_sha256"],"matches_expected":true};v["observed_project"]=map[string]any{"project_id":project["project_id"],"repo_id":project["repo_id"],"registered_locator":project["registered_locator"],"registered_git_common_dir":project["registered_git_common_dir"],"matches_expected":true};observedTarget:=map[string]any{};for k,x:=range target{observedTarget[k]=x};observedTarget["matches_expected"]=true;v["observed_target"]=observedTarget
		observedInputs:=[]any{};readRoots:=[]string{handoffPath,target["worktree"].(string)};for _,raw:=range h["inputs"].([]any){input:=raw.(map[string]any);observedInputs=append(observedInputs,map[string]any{"id":input["id"],"locator":input["locator"],"sha256":input["sha256"],"size_bytes":input["size_bytes"],"matches_expected":true});readRoots=append(readRoots,input["locator"].(string))};sort.Strings(readRoots);writeRoots:=[]string{reply["reply_root"].(string),target["worktree"].(string)};sort.Strings(writeRoots);v["observed_inputs"]=observedInputs;v["contract_digests"]=map[string]any{"authority_sha256":digestBytes(mustJSON(h["authority"])),"budget_sha256":digestBytes(mustJSON(h["budget"])),"verifiers_sha256":digestBytes(mustJSON(h["verifiers"])),"stop_conditions_sha256":digestBytes(mustJSON(h["stop_conditions"]))};v["sandbox"]=map[string]any{"read_roots":stringsAny(readRoots),"write_roots":stringsAny(writeRoots),"temp_root":tempRoot,"matches_contract":true};v["acceptance"]="started";v["issues"]=[]any{};v["schema_version"]=2;v["task_spec_binding"]=h["task_spec_binding"];v["sandbox"].(map[string]any)["read_roots"]=[]any{workspace["root"]};emit(out,v)
	case "terminal":
		out,handoffPath,startPath,artifactPath:=os.Args[2],os.Args[3],os.Args[4],os.Args[5];h:=read(handoffPath);s:=read(startPath);identity:=h["identity"].(map[string]any);target:=h["target_binding"].(map[string]any);artifactDigest,artifactSize:=digestFile(artifactPath);v:=envelope("ply.workflow.terminal-result-draft");v["result_id"]=identity["terminal_result_id"];v["binding"]=map[string]any{"activity_id":identity["activity_id"],"run_id":identity["run_id"],"handoff_id":identity["handoff_id"],"handoff_sha256":digestJSON(handoffPath)};v["start_binding"]=map[string]any{"receipt_id":s["receipt_id"],"start_receipt_sha256":digestJSON(startPath)};v["principal"]=s["principal"];v["reported_outcome"]="complete";v["rounds_used"]=1;v["stop_reasons"]=[]any{};v["summary"]="The exact Task result is complete.";v["meaning"]="Technical delivery evidence is complete; human gates remain separate.";final:=map[string]any{};for k,x:=range target{final[k]=x};final["matches_expected"]=true;v["final_target"]=final
		allowed:=h["authority"].(map[string]any)["allowed_effects"].([]any);v["observed_effects"]=[]any{effect(allowed[0].(map[string]any),0),effect(allowed[1].(map[string]any),1)};verifier:=h["verifiers"].([]any)[0].(map[string]any);v["verifier_results"]=[]any{map[string]any{"verifier_id":"tests","argv":verifier["argv"],"cwd":verifier["cwd"],"exit":0,"bound_oid_or_sha256":target["oid"],"stdout_artifact_id":"verifier-stdout","stderr_artifact_id":nil}};v["review"]=map[string]any{"findings":[]any{},"fixes":[]any{},"open_actionable_findings":[]any{}};v["artifacts"]=[]any{map[string]any{"artifact_id":"verifier-stdout","kind":"managed","description":"Verifier standard output.","media_type":"text/plain","classification":"workspace_internal","size_bytes":artifactSize,"sha256":artifactDigest,"locator":artifactPath}};reqPath:=os.Args[6];reqHash,reqSize:=digestFile(reqPath);v["artifacts"]=append([]any{map[string]any{"artifact_id":"task-requirements","kind":"managed","description":"Fixture requirement coverage.","media_type":"application/json","classification":"workspace_internal","size_bytes":reqSize,"sha256":reqHash,"locator":reqPath}},v["artifacts"].([]any)...);v["evidence_gaps"]=[]any{};v["forbidden_effects_observed"]=[]any{};emit(out,v)
	case "task-result":
		out,showPath,handoffPath,startPath,terminalPath,inspectionPath:=os.Args[2],os.Args[3],os.Args[4],os.Args[5],os.Args[6],os.Args[7];show,h,s,tr:=read(showPath),read(handoffPath),read(startPath),read(terminalPath);persisted:=show["resource"].(map[string]any)["persisted"].(map[string]any);task:=map[string]any{"task_worktree_id":persisted["worktree"].(map[string]any)["id"]};repo:=map[string]any{"git_common_dir":persisted["git_common_dir"]};target:=h["target_binding"].(map[string]any);ids:=h["identity"].(map[string]any);artifact:=tr["artifacts"].([]any)[1].(map[string]any);v:=envelope("WorkspaceTaskResultRecordDraft@1");v["publication_key"]="task/result-blackbox";v["task_id"]="task";v["task_worktree_id"]=task["task_worktree_id"];v["handoff"]=map[string]any{"activity_id":ids["activity_id"],"run_id":ids["run_id"],"handoff_id":ids["handoff_id"],"handoff_locator":handoffPath,"handoff_sha256":digestJSON(handoffPath),"start_receipt_id":s["receipt_id"],"start_receipt_locator":startPath,"start_receipt_sha256":digestJSON(startPath),"terminal_result_id":tr["result_id"],"terminal_result_locator":terminalPath,"terminal_result_sha256":digestJSON(terminalPath),"inspection_sha256":digestJSON(inspectionPath)};v["source"]=map[string]any{"project_id":"ply","repo_id":"ply","git_common_dir":repo["git_common_dir"],"worktree_locator":target["worktree"],"source_ref":target["ref"],"result_oid":target["oid"],"result_tree":target["tree"]};v["technical_assessment"]=map[string]any{"gate":"passed","required_verifier_ids":[]any{"tests"},"accepted_debt":[]any{}};v["evidence_artifacts"]=[]any{map[string]any{"artifact_id":"verifier-stdout","role":"verifier_stdout","locator":artifact["locator"],"sha256":artifact["sha256"],"size_bytes":artifact["size_bytes"]}};req:=tr["artifacts"].([]any)[0].(map[string]any);v["evidence_artifacts"]=append([]any{map[string]any{"artifact_id":"task-requirements","role":"other","locator":req["locator"],"sha256":req["sha256"],"size_bytes":req["size_bytes"]}},v["evidence_artifacts"].([]any)...);v["recorder"]=map[string]any{"actor_claim":"blackbox fixture recorder","control_surface":"fresh Ply CLI","recorded_at_utc":"2026-09-28T12:03:00Z"};emit(out,v)
	case "qa":
		out,resultPath,reportPath:=os.Args[2],os.Args[3],os.Args[4];result:=read(resultPath)["record"].(map[string]any);digest,size:=digestFile(reportPath);v:=envelope("WorkspaceTaskHumanQARecordDraft@1");v["publication_key"]="task/qa-blackbox";v["task_id"]="task";v["task_result_id"]=result["id"];v["result_oid"]=result["result_oid"];v["result_tree"]=result["result_tree"];v["outcome"]="pass";v["actor"]=map[string]any{"actor_claim":"fixture human, not actual approval","start_surface":"local fixture","started_at_utc":"2026-09-28T12:04:00Z","completed_at_utc":"2026-09-28T12:05:00Z"};v["evidence"]=[]any{map[string]any{"id":"report","role":"report","locator":reportPath,"sha256":digest,"size_bytes":size}};v["observation"]="The isolated product journey is correct.";v["accepted_residual_risks"]=[]any{};emit(out,v)
	default: panic("unknown mode")
	}
}
func taskDraft(kind,key string)map[string]any{v:=envelope(kind);v["publication_key"],v["task_id"],v["registry_upgrade"]=key,"task",nil;v["recorder"]=map[string]any{"actor_claim":"fixture recorder","control_surface":"isolated test","recorded_at_utc":"2026-09-29T12:00:00Z"};return v}
func revisionRef(v map[string]any)map[string]any{return map[string]any{"revision":v["revision"],"manifest_sha256":v["manifest_sha256"]}}
func decisionRef(v map[string]any)map[string]any{return map[string]any{"id":v["id"],"manifest_sha256":v["manifest_sha256"]}}
func mustJSON(v any)[]byte{b,e:=json.Marshal(v);must(e);return b}
func effect(v map[string]any,n int)map[string]any{return map[string]any{"effect_id":v["id"],"type":v["type"],"scope":v["scope"],"occurrences":n,"within_authority":true}}
GO
GOCACHE="$temp_root/gocache" "$go_bin" build -o "$temp_root/lifecycle-helper" "$helper_source" || fail 'could not build lifecycle JSON helper'
helper="$temp_root/lifecycle-helper"

run_ok() {
	local label=$1 directory=$2
	shift 2
	set +e
	(cd "$directory" && HOME="$temp_root/home" TMPDIR="$temp_root" "$binary" "$@") >"$temp_root/$label.stdout" 2>"$temp_root/$label.stderr"
	local status=$?
	set -e
	if [[ $status -ne 0 ]]; then sed -n '1,20p' "$temp_root/$label.stderr" >&2; fail "$label exited $status"; fi
	if [[ -s "$temp_root/$label.stderr" ]]; then sed -n '1,20p' "$temp_root/$label.stderr" >&2; fail "$label wrote stderr"; fi
}

workspace="$temp_root/workspace"
repository="$workspace/ply/main"
epic="$workspace/ply/epic"
task="$workspace/ply/task"
mkdir -p "$repository"
run_ok init "$workspace" workspace init
git -C "$repository" init -b main >/dev/null
git -C "$repository" -c user.name='Ply tests' -c user.email=tests@example.invalid commit --allow-empty -m base >/dev/null
base_oid=$(git -C "$repository" rev-parse HEAD)
main_before=$base_oid
git -C "$repository" worktree add -b epic "$epic" "$base_oid" >/dev/null
run_ok project "$workspace" workspace project add ply --name Ply --wrapper "$workspace/ply" --repo "ply=$repository"
run_ok epic "$workspace" workspace epic adopt epic --title Epic --project ply --repo ply --worktree "$epic" --ref refs/heads/epic --expected-oid "$base_oid"
run_ok task-create "$workspace" workspace task create task --title Task --description 'Integration Task' --epic epic --project ply --repo ply
run_ok task-worktree "$workspace" workspace task worktree create task --branch task --path "$task" --expected-parent-oid "$base_oid"
printf 'delivered\n' >"$task/delivery.txt"
git -C "$task" add delivery.txt
git -C "$task" -c user.name='Ply tests' -c user.email=tests@example.invalid commit -m delivery >/dev/null
result_oid=$(git -C "$task" rev-parse HEAD)
result_tree=$(git -C "$task" rev-parse 'HEAD^{tree}')
git -C "$epic" -c user.name='Ply tests' -c user.email=tests@example.invalid commit --allow-empty -m 'seed integration reflog' >/dev/null
seed_oid=$(git -C "$epic" rev-parse HEAD)
git -C "$epic" update-ref -m 'restore integration base' refs/heads/epic "$base_oid" "$seed_oid"
printf 'immutable input\n' >"$workspace/input.txt"

run_ok content-initial "$workspace" workspace task show task --format json
"$helper" spec "$workspace/spec-1.json" "$temp_root/content-initial.stdout" "$workspace" 1 -
run_ok spec-record "$workspace" workspace task spec record task --file "$workspace/spec-1.json" --format json
"$helper" assessment "$workspace/assessment.json" "$temp_root/spec-record.stdout"
run_ok spec-assess "$workspace" workspace task spec assess task --file "$workspace/assessment.json" --format json
"$helper" selection "$workspace/selection.json" "$temp_root/spec-record.stdout" "$temp_root/spec-assess.stdout" "$temp_root/content-initial.stdout"
run_ok spec-select "$workspace" workspace task spec select task --file "$workspace/selection.json" --format json
run_ok selected-spec "$workspace" workspace task spec show task --spec solution --revision 1 --format json
"$helper" spec "$workspace/spec-2.json" "$temp_root/content-initial.stdout" "$workspace" 2 "$temp_root/selected-spec.stdout"
run_ok spec-record-2 "$workspace" workspace task spec record task --file "$workspace/spec-2.json" --format json
rm "$workspace/abstract.md" "$workspace/functional.md" "$workspace/technical.md"
run_ok spec-retry-without-source "$workspace" workspace task spec record task --file "$workspace/spec-1.json" --format json
run_ok selected-after-draft "$workspace" workspace task show task --format json
[[ $("$helper" field "$temp_root/selected-after-draft.stdout" solution.selected_revision.revision) == 1 && $("$helper" field "$temp_root/selected-after-draft.stdout" solution.newer_draft_available) == true ]] || fail 'r2 silently changed the chosen revision'
"$helper" handoff "$workspace/handoff-draft.json" "$task" refs/heads/task "$result_oid" "$workspace/input.txt" "$temp_root/selected-spec.stdout"
for invalid_mode in v1 null wrong-worktree missing-input forged-input; do
	"$helper" invalid-handoff "$workspace/invalid-handoff.json" "$workspace/handoff-draft.json" "$invalid_mode"
	if (cd "$workspace" && HOME="$temp_root/home" "$binary" workflow handoff create --file "$workspace/invalid-handoff.json") >"$temp_root/invalid-$invalid_mode.stdout" 2>"$temp_root/invalid-$invalid_mode.stderr"; then fail "accepted invalid Task handoff: $invalid_mode"; fi
	grep -F 'task_spec' "$temp_root/invalid-$invalid_mode.stderr" >/dev/null || fail "invalid $invalid_mode failed outside the Task guard"
done
run_ok handoff-create "$workspace" workflow handoff create --file "$workspace/handoff-draft.json"
handoff_locator=$(sed -n 's/^Handoff: //p' "$temp_root/handoff-create.stdout")
"$helper" other-handoff "$workspace/other-handoff.json" "$workspace/handoff-draft.json"
run_ok other-handoff-create "$workspace" workflow handoff create --file "$workspace/other-handoff.json"
other_handoff=$(sed -n 's/^Handoff: //p' "$temp_root/other-handoff-create.stdout")
"$helper" start "$workspace/start-draft.json" "$handoff_locator" "$temp_root"
run_ok handoff-start "$task" workflow handoff submit-start --handoff "$handoff_locator" --file "$workspace/start-draft.json"
start_locator=$(sed -n 's/^Start receipt: //p' "$temp_root/handoff-start.stdout")
reply_root=$("$helper" field "$handoff_locator" reply_capability.reply_root)
: >"$reply_root/staging/verifier.stdout"
"$helper" requirements "$reply_root/staging/task-requirements.json" "$handoff_locator"
"$helper" terminal "$workspace/terminal-draft.json" "$handoff_locator" "$start_locator" "$reply_root/staging/verifier.stdout" "$reply_root/staging/task-requirements.json"
run_ok handoff-result "$task" workflow handoff submit-result --handoff "$handoff_locator" --file "$workspace/terminal-draft.json"
terminal_locator=$(sed -n 's/^Terminal result: //p' "$temp_root/handoff-result.stdout")
run_ok handoff-inspect "$workspace" workflow handoff inspect --handoff "$handoff_locator" --format json
cp "$temp_root/handoff-inspect.stdout" "$workspace/inspection.json"
run_ok task-show-before "$workspace" workspace task show task --format json

store="$workspace/.ply/work-items.yaml"

"$helper" task-result "$workspace/task-result.json" "$temp_root/task-show-before.stdout" "$handoff_locator" "$start_locator" "$terminal_locator" "$workspace/inspection.json"
run_ok result-record "$workspace" workspace task result record task --file "$workspace/task-result.json" --format json
[[ $("$helper" field "$temp_root/result-record.stdout" spec_binding.basis.spec.revision) == 1 ]] || fail 'result lost exact selected r1 binding'
grep -F '"kind":"WorkspaceTaskResultRecordReadback@2"' "$temp_root/result-record.stdout" >/dev/null || fail 'v3 result readback is missing'
touch -t 200001010000 "$store"
result_mtime=$(file_mtime "$store")
result_store_sha=$(file_sha256 "$store")
run_ok result-retry "$workspace" workspace task result record task --file "$workspace/task-result.json" --format json
[[ $(file_mtime "$store") == "$result_mtime" && $(file_sha256 "$store") == "$result_store_sha" ]] || fail 'Task result retry rewrote the store'

printf 'Fixture human QA passed; not actual human approval.\n' >"$workspace/qa-report.txt"
"$helper" qa "$workspace/qa.json" "$temp_root/result-record.stdout" "$workspace/qa-report.txt"
run_ok qa-record "$workspace" workspace task qa record task --file "$workspace/qa.json" --format json
task_result_id=$("$helper" field "$temp_root/result-record.stdout" record.id)
qa_id=$("$helper" field "$temp_root/qa-record.stdout" record.id)
run_ok integrate-check "$workspace" workspace task integrate task --result "$task_result_id" --qa "$qa_id" --expected-result-oid "$result_oid" --expected-parent-oid "$base_oid" --check --format json
confirm=$("$helper" field "$temp_root/integrate-check.stdout" plan.sha256)
[[ $(git -C "$epic" rev-parse HEAD) == "$base_oid" ]] || fail 'read-only check changed the Epic parent'
run_ok integrate-apply "$workspace" workspace task integrate task --result "$task_result_id" --qa "$qa_id" --expected-result-oid "$result_oid" --expected-parent-oid "$base_oid" --apply --confirm "$confirm" --format json
grep -F '"classification":"exact_effect"' "$temp_root/integrate-apply.stdout" >/dev/null || fail 'confirmed apply was not exact_effect'
[[ $(git -C "$epic" rev-parse HEAD) == "$result_oid" && $(git -C "$epic" rev-parse 'HEAD^{tree}') == "$result_tree" ]] || fail 'Epic parent did not fast-forward to the exact result'
[[ $(git -C "$repository" rev-parse refs/heads/main) == "$main_before" ]] || fail 'ordinary main changed'
touch -t 200001010000 "$store"
apply_store_sha=$(file_sha256 "$store")
apply_store_mtime=$(file_mtime "$store")
run_ok integrate-retry "$workspace" workspace task integrate task --result "$task_result_id" --qa "$qa_id" --expected-result-oid "$result_oid" --expected-parent-oid "$base_oid" --apply --confirm "$confirm" --format json
[[ $(file_sha256 "$store") == "$apply_store_sha" && $(file_mtime "$store") == "$apply_store_mtime" ]] || fail 'confirmed retry rewrote the store'
run_ok task-show-after "$workspace" workspace task show task --format json
grep -F '"kind":"WorkspaceTaskIntegrationReadback@2"' "$temp_root/task-show-after.stdout" >/dev/null || fail 'Spec-aware Task integration readback is missing'
grep -F '"classification":"exact_effect"' "$temp_root/task-show-after.stdout" >/dev/null || fail 'Task show lost the exact integration result'
"$helper" problem2 "$workspace/problem-2.json" "$temp_root/task-show-after.stdout"
run_ok problem-2 "$workspace" workspace task problem record task --file "$workspace/problem-2.json" --format json
run_ok historical-inspection "$workspace" workflow handoff inspect --handoff "$handoff_locator" --format json
[[ $(file_sha256 "$temp_root/historical-inspection.stdout") == $(file_sha256 "$workspace/inspection.json") ]] || fail 'P2 changed historical inspection bytes'
run_ok historical-result-retry "$workspace" workspace task result record task --file "$workspace/task-result.json" --format json
"$helper" start "$workspace/stale-start.json" "$other_handoff" "$temp_root"
set +e
(cd "$task" && HOME="$temp_root/home" "$binary" workflow handoff submit-start --handoff "$other_handoff" --file "$workspace/stale-start.json") >"$temp_root/stale-start.stdout" 2>"$temp_root/stale-start.stderr"
stale_status=$?
set -e
[[ $stale_status -ne 0 ]] || fail 'old unstarted handoff was accepted after P2'
"$helper" negative-start "$workspace/negative-start.json" "$workspace/stale-start.json"
run_ok accepted-negative "$task" workflow handoff submit-start --handoff "$other_handoff" --file "$workspace/negative-start.json"
run_ok stale-integration "$workspace" workspace task integrate task --result "$task_result_id" --qa "$qa_id" --expected-result-oid "$result_oid" --expected-parent-oid "$base_oid" --check --format json
[[ $("$helper" field "$temp_root/stale-integration.stdout" task_spec_relevance.relevance) == stale ]] || fail 'old result silently gained current integration authority'
[[ ! -e "$temp_root/home/.ply" ]] || fail 'blackbox wrote global HOME state'

(
	cd "$repo_root"
	GOCACHE="$temp_root/gocache" "$go_bin" test ./internal/workspace ./internal/workflowhandoff \
		-run 'Test(TaskIntegration(CheckAndSingleConfirmedApply|ConfirmationMismatchStartsNoAuthorityOrGit|RetriesAfterNamedNoEffectWithOneNewAttempt|RecoversLostResponseFromReflogWithoutSecondMerge|RecoversDurableAuthorityAndAttemptPublicationBoundaries|ConcurrentIdenticalApplyProducesOneAttempt|ClassifiesAlreadyIntegratedBlockedConflictAndPartial)|WorkspaceTaskEvidenceReaderProjectsExactImmutableEvidence)$' \
		-count=1
) >"$temp_root/lifecycle.stdout" 2>"$temp_root/lifecycle.stderr" || fail 'isolated lifecycle roundtrip failed'
[[ ! -s "$temp_root/lifecycle.stderr" ]] || fail 'isolated lifecycle roundtrip wrote stderr'
grep -F 'ok  ' "$temp_root/lifecycle.stdout" >/dev/null || fail 'isolated lifecycle roundtrip did not report success'
[[ ! -e "$temp_root/home/.ply" ]] || fail 'roundtrip wrote global HOME state'

printf 'workspace task integration roundtrip: PASS\n'
