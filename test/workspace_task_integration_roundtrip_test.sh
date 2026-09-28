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
	case "handoff":
		out,target,ref,oid,input:=os.Args[2],os.Args[3],os.Args[4],os.Args[5],os.Args[6];inputDigest,inputSize:=digestFile(input);v:=envelope("ply.workflow.handoff-draft")
		v["publication_key"],v["activity_key"]="ws04/blackbox/happy","ws04/blackbox/happy"
		v["goal"]=map[string]any{"title":"Workspace Task integration blackbox","recipient_role":"Delivery agent","objective":"Preserve exact delivery evidence.","done_when":"The verifier passes and the exact result is reported."}
		v["recipient"]=map[string]any{"principal_id":"codex-delivery-agent","principal_kind":"human_started_agent","runtime_constraints":[]any{"local"}}
		v["binding_request"]=map[string]any{"project_id":"ply","repo_id":"ply","target_worktree":target,"target_ref":ref,"expected_oid":oid,"status_policy":map[string]any{"mode":"clean"}}
		v["inputs"]=[]any{map[string]any{"id":"fixture","role":"verification_fixture","locator":input,"sha256":inputDigest,"size_bytes":inputSize,"media_type":"text/plain","git_binding":nil}}
		v["authority"]=map[string]any{"allowed_effects":[]any{map[string]any{"id":"write-result","type":"filesystem_write","scope":map[string]any{"kind":"filesystem","paths":[]any{"delivery.txt"},"directory_prefixes":[]any{}},"max_occurrences":1,"sequence":1},map[string]any{"id":"run-tests","type":"command_execute","scope":map[string]any{"kind":"command","procedure_ids":[]any{"verify"},"verifier_ids":[]any{"tests"}},"max_occurrences":1,"sequence":2}},"forbidden_effects":[]any{map[string]any{"type":"network","reason":"Network is forbidden."},map[string]any{"type":"push","reason":"Push is forbidden."}},"human_gates":[]any{"human_task_qa","local_integration"}}
		v["budget"]=map[string]any{"max_rounds":2,"round_definition":map[string]any{"unit":"implementation_or_review_fix_iteration","command_retry_consumes_round":false,"retry_condition":"only_if_no_effect_started"}}
		v["procedure"]=[]any{map[string]any{"id":"verify","instruction":"Verify the exact local result.","required_before":[]any{}}}
		v["verifiers"]=[]any{map[string]any{"id":"tests","argv":[]any{"git","status","--short"},"cwd":target,"env":[]any{},"expected_exit":0,"stop_on_failure":true,"evidence":map[string]any{"capture_stdout":true,"capture_stderr":false,"classification":"workspace_internal","binding":"target_oid"}}}
		stops:=[]string{"product_decision_required","scope_or_authority_expansion","target_or_input_drift","unknown_or_partial_effect","unexpected_sensitive_data","round_budget_exhausted"};values:=[]any{};for _,s:=range stops{values=append(values,map[string]any{"type":s,"description":"Stop safely when this condition occurs."})};v["stop_conditions"]=values
		v["reporting"]=map[string]any{"summary_max_codepoints":240,"meaning_max_codepoints":600,"required_start_fields":stringsAny([]string{"acceptance","binding","contract_digests","issues","observed_inputs","observed_project","observed_target","observed_workspace","principal","receipt_id","sandbox"}),"required_terminal_fields":stringsAny([]string{"artifacts","binding","evidence_gaps","final_target","forbidden_effects_observed","meaning","observed_effects","principal","reported_outcome","result_id","review","rounds_used","start_binding","stop_reasons","summary","verifier_results"})};emit(out,v)
	case "start":
		out,handoffPath,tempRoot:=os.Args[2],os.Args[3],os.Args[4];h:=read(handoffPath);identity:=h["identity"].(map[string]any);workspace:=h["workspace_binding"].(map[string]any);project:=h["project_binding"].(map[string]any);target:=h["target_binding"].(map[string]any);reply:=h["reply_capability"].(map[string]any);principal:=map[string]any{"expected_principal_id":"codex-delivery-agent","human_start_principal":"blackbox-human","start_surface":"local shell","session_id":"blackbox-session","runtime_id":"codex","model_id":"test-model","started_at_utc":"2026-09-28T12:00:00Z"};v:=envelope("ply.workflow.start-receipt-draft")
		v["receipt_id"]=identity["start_receipt_id"];v["binding"]=map[string]any{"activity_id":identity["activity_id"],"run_id":identity["run_id"],"handoff_id":identity["handoff_id"],"handoff_sha256":digestJSON(handoffPath)};v["principal"]=principal;v["observed_workspace"]=map[string]any{"root":workspace["root"],"marker_format_version":workspace["marker_format_version"],"marker_sha256":workspace["marker_sha256"],"matches_expected":true};v["observed_project"]=map[string]any{"project_id":project["project_id"],"repo_id":project["repo_id"],"registered_locator":project["registered_locator"],"registered_git_common_dir":project["registered_git_common_dir"],"matches_expected":true};observedTarget:=map[string]any{};for k,x:=range target{observedTarget[k]=x};observedTarget["matches_expected"]=true;v["observed_target"]=observedTarget
		observedInputs:=[]any{};readRoots:=[]string{handoffPath,target["worktree"].(string)};for _,raw:=range h["inputs"].([]any){input:=raw.(map[string]any);observedInputs=append(observedInputs,map[string]any{"id":input["id"],"locator":input["locator"],"sha256":input["sha256"],"size_bytes":input["size_bytes"],"matches_expected":true});readRoots=append(readRoots,input["locator"].(string))};sort.Strings(readRoots);writeRoots:=[]string{reply["reply_root"].(string),target["worktree"].(string)};sort.Strings(writeRoots);v["observed_inputs"]=observedInputs;v["contract_digests"]=map[string]any{"authority_sha256":digestBytes(mustJSON(h["authority"])),"budget_sha256":digestBytes(mustJSON(h["budget"])),"verifiers_sha256":digestBytes(mustJSON(h["verifiers"])),"stop_conditions_sha256":digestBytes(mustJSON(h["stop_conditions"]))};v["sandbox"]=map[string]any{"read_roots":stringsAny(readRoots),"write_roots":stringsAny(writeRoots),"temp_root":tempRoot,"matches_contract":true};v["acceptance"]="started";v["issues"]=[]any{};emit(out,v)
	case "terminal":
		out,handoffPath,startPath,artifactPath:=os.Args[2],os.Args[3],os.Args[4],os.Args[5];h:=read(handoffPath);s:=read(startPath);identity:=h["identity"].(map[string]any);target:=h["target_binding"].(map[string]any);artifactDigest,artifactSize:=digestFile(artifactPath);v:=envelope("ply.workflow.terminal-result-draft");v["result_id"]=identity["terminal_result_id"];v["binding"]=map[string]any{"activity_id":identity["activity_id"],"run_id":identity["run_id"],"handoff_id":identity["handoff_id"],"handoff_sha256":digestJSON(handoffPath)};v["start_binding"]=map[string]any{"receipt_id":s["receipt_id"],"start_receipt_sha256":digestJSON(startPath)};v["principal"]=s["principal"];v["reported_outcome"]="complete";v["rounds_used"]=1;v["stop_reasons"]=[]any{};v["summary"]="The exact Task result is complete.";v["meaning"]="Technical delivery evidence is complete; human gates remain separate.";final:=map[string]any{};for k,x:=range target{final[k]=x};final["matches_expected"]=true;v["final_target"]=final
		allowed:=h["authority"].(map[string]any)["allowed_effects"].([]any);v["observed_effects"]=[]any{effect(allowed[0].(map[string]any),0),effect(allowed[1].(map[string]any),1)};verifier:=h["verifiers"].([]any)[0].(map[string]any);v["verifier_results"]=[]any{map[string]any{"verifier_id":"tests","argv":verifier["argv"],"cwd":verifier["cwd"],"exit":0,"bound_oid_or_sha256":target["oid"],"stdout_artifact_id":"verifier-stdout","stderr_artifact_id":nil}};v["review"]=map[string]any{"findings":[]any{},"fixes":[]any{},"open_actionable_findings":[]any{}};v["artifacts"]=[]any{map[string]any{"artifact_id":"verifier-stdout","kind":"managed","description":"Verifier standard output.","media_type":"text/plain","classification":"workspace_internal","size_bytes":artifactSize,"sha256":artifactDigest,"locator":artifactPath}};v["evidence_gaps"]=[]any{};v["forbidden_effects_observed"]=[]any{};emit(out,v)
	case "task-result":
		out,showPath,handoffPath,startPath,terminalPath,inspectionPath:=os.Args[2],os.Args[3],os.Args[4],os.Args[5],os.Args[6],os.Args[7];show,h,s,tr:=read(showPath),read(handoffPath),read(startPath),read(terminalPath);task:=show["task"].(map[string]any);repo:=show["repository"].(map[string]any);target:=h["target_binding"].(map[string]any);ids:=h["identity"].(map[string]any);artifact:=tr["artifacts"].([]any)[0].(map[string]any);v:=envelope("WorkspaceTaskResultRecordDraft@1");v["publication_key"]="task/result-blackbox";v["task_id"]="task";v["task_worktree_id"]=task["task_worktree_id"];v["handoff"]=map[string]any{"activity_id":ids["activity_id"],"run_id":ids["run_id"],"handoff_id":ids["handoff_id"],"handoff_locator":handoffPath,"handoff_sha256":digestJSON(handoffPath),"start_receipt_id":s["receipt_id"],"start_receipt_locator":startPath,"start_receipt_sha256":digestJSON(startPath),"terminal_result_id":tr["result_id"],"terminal_result_locator":terminalPath,"terminal_result_sha256":digestJSON(terminalPath),"inspection_sha256":digestJSON(inspectionPath)};v["source"]=map[string]any{"project_id":"ply","repo_id":"ply","git_common_dir":repo["git_common_dir"],"worktree_locator":target["worktree"],"source_ref":target["ref"],"result_oid":target["oid"],"result_tree":target["tree"]};v["technical_assessment"]=map[string]any{"gate":"passed","required_verifier_ids":[]any{"tests"},"accepted_debt":[]any{}};v["evidence_artifacts"]=[]any{map[string]any{"artifact_id":"verifier-stdout","role":"verifier_stdout","locator":artifact["locator"],"sha256":artifact["sha256"],"size_bytes":artifact["size_bytes"]}};v["recorder"]=map[string]any{"actor_claim":"blackbox recorder","control_surface":"fresh Ply CLI","recorded_at_utc":"2026-09-28T12:03:00Z"};emit(out,v)
	case "qa":
		out,resultPath,reportPath:=os.Args[2],os.Args[3],os.Args[4];result:=read(resultPath)["record"].(map[string]any);digest,size:=digestFile(reportPath);v:=envelope("WorkspaceTaskHumanQARecordDraft@1");v["publication_key"]="task/qa-blackbox";v["task_id"]="task";v["task_result_id"]=result["id"];v["result_oid"]=result["result_oid"];v["result_tree"]=result["result_tree"];v["outcome"]="pass";v["actor"]=map[string]any{"actor_claim":"blackbox human","start_surface":"local fixture","started_at_utc":"2026-09-28T12:04:00Z","completed_at_utc":"2026-09-28T12:05:00Z"};v["evidence"]=[]any{map[string]any{"id":"report","role":"report","locator":reportPath,"sha256":digest,"size_bytes":size}};v["observation"]="The isolated product journey is correct.";v["accepted_residual_risks"]=[]any{};emit(out,v)
	default: panic("unknown mode")
	}
}
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

"$helper" handoff "$workspace/handoff-draft.json" "$task" refs/heads/task "$result_oid" "$workspace/input.txt"
run_ok handoff-create "$workspace" workflow handoff create --file "$workspace/handoff-draft.json"
handoff_locator=$(sed -n 's/^Handoff: //p' "$temp_root/handoff-create.stdout")
"$helper" start "$workspace/start-draft.json" "$handoff_locator" "$temp_root"
run_ok handoff-start "$task" workflow handoff submit-start --handoff "$handoff_locator" --file "$workspace/start-draft.json"
start_locator=$(sed -n 's/^Start receipt: //p' "$temp_root/handoff-start.stdout")
reply_root=$("$helper" field "$handoff_locator" reply_capability.reply_root)
: >"$reply_root/staging/verifier.stdout"
"$helper" terminal "$workspace/terminal-draft.json" "$handoff_locator" "$start_locator" "$reply_root/staging/verifier.stdout"
run_ok handoff-result "$task" workflow handoff submit-result --handoff "$handoff_locator" --file "$workspace/terminal-draft.json"
terminal_locator=$(sed -n 's/^Terminal result: //p' "$temp_root/handoff-result.stdout")
run_ok handoff-inspect "$workspace" workflow handoff inspect --handoff "$handoff_locator" --format json
cp "$temp_root/handoff-inspect.stdout" "$workspace/inspection.json"
run_ok task-show-before "$workspace" workspace task show task --format json

store="$workspace/.ply/work-items.yaml"
sed -e '/^task_results: \[\]$/d' -e '/^human_qa_records: \[\]$/d' -e '/^integration_authorities: \[\]$/d' -e '/^integration_intents: \[\]$/d' -e '/^integration_attempts: \[\]$/d' -e '/^integration_results: \[\]$/d' "$store" >"$workspace/work-items-v1.yaml"
mv "$workspace/work-items-v1.yaml" "$store"
chmod 0644 "$store"
sed -i.bak 's/^format_version: 2$/format_version: 1/' "$store"
rm -f "$store.bak"
grep -F 'format_version: 1' "$store" >/dev/null || fail 'legacy fixture does not report format 1'

"$helper" task-result "$workspace/task-result.json" "$temp_root/task-show-before.stdout" "$handoff_locator" "$start_locator" "$terminal_locator" "$workspace/inspection.json"
run_ok result-record "$workspace" workspace task result record task --file "$workspace/task-result.json" --format json
grep -F '"store_transition":"format_1_to_2"' "$temp_root/result-record.stdout" >/dev/null || fail 'format-1 migration was not reported'
touch -t 200001010000 "$store"
result_mtime=$(file_mtime "$store")
result_store_sha=$(file_sha256 "$store")
run_ok result-retry "$workspace" workspace task result record task --file "$workspace/task-result.json" --format json
[[ $(file_mtime "$store") == "$result_mtime" && $(file_sha256 "$store") == "$result_store_sha" ]] || fail 'Task result retry rewrote the store'

printf 'Human QA passed.\n' >"$workspace/qa-report.txt"
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
grep -F '"kind":"WorkspaceTaskIntegrationReadback@1"' "$temp_root/task-show-after.stdout" >/dev/null || fail 'format-2 Task show schema is missing'
grep -F '"classification":"exact_effect"' "$temp_root/task-show-after.stdout" >/dev/null || fail 'Task show lost the exact integration result'
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
