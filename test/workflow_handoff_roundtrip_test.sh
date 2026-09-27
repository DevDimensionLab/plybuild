#!/usr/bin/env bash

set -euo pipefail
export LC_ALL=C LANG=C GOENV=off GOWORK=off GOPROXY=off
unset CDPATH

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
temp_root=$(mktemp -d "${TMPDIR:-/tmp}/ply-workflow-handoff.XXXXXX")
temp_root=$(cd "$temp_root" && pwd -P)
cleanup_temp_root() {
	find "$temp_root" -type d -exec chmod u+rwx {} + 2>/dev/null || true
	rm -rf -- "$temp_root"
}
trap cleanup_temp_root EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM

fail() {
	printf 'workflow handoff roundtrip: %s\n' "$*" >&2
	exit 1
}

file_mode() {
	if [[ $(uname -s) == Darwin ]]; then stat -f '%Lp' "$1"; else stat -c '%a' "$1"; fi
}

file_mtime() {
	if [[ $(uname -s) == Darwin ]]; then stat -f '%m' "$1"; else stat -c '%Y' "$1"; fi
}

file_sha256() {
	if command -v shasum >/dev/null 2>&1; then shasum -a 256 "$1" | awk '{print "sha256:"$1}'; else sha256sum "$1" | awk '{print "sha256:"$1}'; fi
}

mkdir -m 700 "$temp_root/home" "$temp_root/gocache" "$temp_root/recipient-temp"
binary="$temp_root/ply"
(
	cd "$repo_root"
	GOCACHE="$temp_root/gocache" "${GO:-go}" build -o "$binary" ./cmd/ply
) || fail 'could not build a fresh Ply binary'

helper_source="$temp_root/json_helper.go"
cat >"$helper_source" <<'GO'
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func must(err error) { if err != nil { panic(err) } }
func read(path string) map[string]any { b,err:=os.ReadFile(path);must(err);var v map[string]any;must(json.Unmarshal(b,&v));return v }
func emit(path string,v any){b,err:=json.Marshal(v);must(err);must(os.WriteFile(path,b,0600))}
func digestValue(v any)string{b,err:=json.Marshal(v);must(err);s:=sha256.Sum256(b);return "sha256:"+hex.EncodeToString(s[:])}
func digestFile(path string)(string,int64){b,err:=os.ReadFile(path);must(err);s:=sha256.Sum256(b);return "sha256:"+hex.EncodeToString(s[:]),int64(len(b))}
func stringsAny(values []string)[]any{r:=make([]any,len(values));for i,v:=range values{r[i]=v};return r}
func envelope(kind string)map[string]any{return map[string]any{"kind":kind,"schema_version":1,"format":"json","format_version":1,"canonicalization":"RFC8785"}}
func main(){
	if len(os.Args)<2{panic("mode required")}
	switch os.Args[1]{
	case "field":
		v:=any(read(os.Args[2]));for _,part:=range strings.Split(os.Args[3],"."){v=v.(map[string]any)[part]};fmt.Print(v)
	case "draft":
		if len(os.Args)!=11{panic("draft output target ref oid git-common input publication activity title")}
		out,target,ref,oid,common,input,pub,activity,title:=os.Args[2],os.Args[3],os.Args[4],os.Args[5],os.Args[6],os.Args[7],os.Args[8],os.Args[9],os.Args[10]
		inputDigest,inputSize:=digestFile(input);v:=envelope("ply.workflow.handoff-draft")
		v["publication_key"],v["activity_key"]=pub,activity
		v["goal"]=map[string]any{"title":title,"recipient_role":"Delivery agent","objective":"Perform the bounded local change.","done_when":"The declared verifier passes and the result is reported."}
		v["recipient"]=map[string]any{"principal_id":"codex-delivery-agent","principal_kind":"human_started_agent","runtime_constraints":[]any{"local"}}
		v["binding_request"]=map[string]any{"project_id":"demo","repo_id":"demo","target_worktree":target,"target_ref":ref,"expected_oid":oid,"status_policy":map[string]any{"mode":"clean"}}
		v["inputs"]=[]any{map[string]any{"id":"fixture","role":"verification_fixture","locator":input,"sha256":inputDigest,"size_bytes":inputSize,"media_type":"text/plain","git_binding":nil}}
		allowed:=[]any{
			map[string]any{"id":"write-readme","type":"filesystem_write","scope":map[string]any{"kind":"filesystem","paths":[]any{"README.md"},"directory_prefixes":[]any{}},"max_occurrences":1,"sequence":1},
			map[string]any{"id":"run-tests","type":"command_execute","scope":map[string]any{"kind":"command","procedure_ids":[]any{"implement"},"verifier_ids":[]any{"tests"}},"max_occurrences":1,"sequence":2},
		}
		v["authority"]=map[string]any{"allowed_effects":allowed,"forbidden_effects":[]any{map[string]any{"type":"network","reason":"Network is outside this local task."},map[string]any{"type":"push","reason":"Push requires a later human gate."}},"human_gates":[]any{"human_task_qa","local_integration"}}
		v["budget"]=map[string]any{"max_rounds":2,"round_definition":map[string]any{"unit":"implementation_or_review_fix_iteration","command_retry_consumes_round":false,"retry_condition":"only_if_no_effect_started"}}
		v["procedure"]=[]any{map[string]any{"id":"implement","instruction":"Perform and review the bounded local change.","required_before":[]any{}}}
		v["verifiers"]=[]any{map[string]any{"id":"tests","argv":[]any{"git","status","--short"},"cwd":target,"env":[]any{},"expected_exit":0,"stop_on_failure":true,"evidence":map[string]any{"capture_stdout":true,"capture_stderr":false,"classification":"workspace_internal","binding":"target_oid"}}}
		stops:=[]string{"product_decision_required","scope_or_authority_expansion","target_or_input_drift","unknown_or_partial_effect","unexpected_sensitive_data","round_budget_exhausted"};stopValues:=[]any{};for _,s:=range stops{stopValues=append(stopValues,map[string]any{"type":s,"description":"Stop safely when this condition occurs."})};v["stop_conditions"]=stopValues
		v["reporting"]=map[string]any{"summary_max_codepoints":240,"meaning_max_codepoints":600,"required_start_fields":stringsAny([]string{"acceptance","binding","contract_digests","issues","observed_inputs","observed_project","observed_target","observed_workspace","principal","receipt_id","sandbox"}),"required_terminal_fields":stringsAny([]string{"artifacts","binding","evidence_gaps","final_target","forbidden_effects_observed","meaning","observed_effects","principal","reported_outcome","result_id","review","rounds_used","start_binding","stop_reasons","summary","verifier_results"})}
		_ = common
		emit(out,v)
	case "start":
		if len(os.Args)!=5{panic("start output handoff temp-root")};out,handoffPath,tempRoot:=os.Args[2],os.Args[3],os.Args[4];h:=read(handoffPath);identity:=h["identity"].(map[string]any);workspace:=h["workspace_binding"].(map[string]any);project:=h["project_binding"].(map[string]any);target:=h["target_binding"].(map[string]any);reply:=h["reply_capability"].(map[string]any)
		principal:=map[string]any{"expected_principal_id":"codex-delivery-agent","human_start_principal":"roundtrip-human","start_surface":"local shell","session_id":"roundtrip-session","runtime_id":"codex","model_id":"test-model","started_at_utc":"2026-09-27T00:00:00Z"}
		v:=envelope("ply.workflow.start-receipt-draft");v["receipt_id"]=identity["start_receipt_id"];v["binding"]=map[string]any{"activity_id":identity["activity_id"],"run_id":identity["run_id"],"handoff_id":identity["handoff_id"],"handoff_sha256":digestValue(h)};v["principal"]=principal
		v["observed_workspace"]=map[string]any{"root":workspace["root"],"marker_format_version":workspace["marker_format_version"],"marker_sha256":workspace["marker_sha256"],"matches_expected":true}
		v["observed_project"]=map[string]any{"project_id":project["project_id"],"repo_id":project["repo_id"],"registered_locator":project["registered_locator"],"registered_git_common_dir":project["registered_git_common_dir"],"matches_expected":true}
		observedTarget:=map[string]any{};for k,x:=range target{observedTarget[k]=x};observedTarget["matches_expected"]=true;v["observed_target"]=observedTarget
		observedInputs:=[]any{};readRoots:=[]string{handoffPath,target["worktree"].(string)};for _,raw:=range h["inputs"].([]any){input:=raw.(map[string]any);observedInputs=append(observedInputs,map[string]any{"id":input["id"],"locator":input["locator"],"sha256":input["sha256"],"size_bytes":input["size_bytes"],"matches_expected":true});readRoots=append(readRoots,input["locator"].(string))};sort.Strings(readRoots);v["observed_inputs"]=observedInputs
		v["contract_digests"]=map[string]any{"authority_sha256":digestValue(h["authority"]),"budget_sha256":digestValue(h["budget"]),"verifiers_sha256":digestValue(h["verifiers"]),"stop_conditions_sha256":digestValue(h["stop_conditions"])}
		writeRoots:=[]string{reply["reply_root"].(string),target["worktree"].(string)};sort.Strings(writeRoots);v["sandbox"]=map[string]any{"read_roots":stringsAny(readRoots),"write_roots":stringsAny(writeRoots),"temp_root":tempRoot,"matches_contract":true};v["acceptance"]="started";v["issues"]=[]any{};emit(out,v)
	case "result":
		if len(os.Args)!=7{panic("result output handoff accepted-start artifact summary")};out,handoffPath,startPath,artifactPath,summary:=os.Args[2],os.Args[3],os.Args[4],os.Args[5],os.Args[6];h:=read(handoffPath);s:=read(startPath);identity:=h["identity"].(map[string]any);target:=h["target_binding"].(map[string]any);artifactDigest,artifactSize:=digestFile(artifactPath)
		v:=envelope("ply.workflow.terminal-result-draft");v["result_id"]=identity["terminal_result_id"];v["binding"]=map[string]any{"activity_id":identity["activity_id"],"run_id":identity["run_id"],"handoff_id":identity["handoff_id"],"handoff_sha256":digestValue(h)};v["start_binding"]=map[string]any{"receipt_id":s["receipt_id"],"start_receipt_sha256":digestValue(s)};v["principal"]=s["principal"];v["reported_outcome"]="complete";v["rounds_used"]=1;v["stop_reasons"]=[]any{};v["summary"]=summary;v["meaning"]="The bounded report is preserved; downstream human gates remain closed."
		finalTarget:=map[string]any{};for k,x:=range target{finalTarget[k]=x};finalTarget["matches_expected"]=true;v["final_target"]=finalTarget
		effects:=[]any{};for _,raw:=range h["authority"].(map[string]any)["allowed_effects"].([]any){effect:=raw.(map[string]any);effects=append(effects,map[string]any{"effect_id":effect["id"],"type":effect["type"],"scope":effect["scope"],"occurrences":1,"within_authority":true})};v["observed_effects"]=effects
		verifier:=h["verifiers"].([]any)[0].(map[string]any);v["verifier_results"]=[]any{map[string]any{"verifier_id":verifier["id"],"argv":verifier["argv"],"cwd":verifier["cwd"],"exit":0,"bound_oid_or_sha256":target["oid"],"stdout_artifact_id":"verifier-stdout","stderr_artifact_id":nil}}
		v["review"]=map[string]any{"findings":[]any{},"fixes":[]any{},"open_actionable_findings":[]any{}};v["artifacts"]=[]any{map[string]any{"artifact_id":"verifier-stdout","kind":"managed","description":"Verifier standard output.","media_type":"text/plain","classification":"workspace_internal","size_bytes":artifactSize,"sha256":artifactDigest,"locator":artifactPath}};v["evidence_gaps"]=[]any{};v["forbidden_effects_observed"]=[]any{};emit(out,v)
	case "size": info,err:=os.Stat(os.Args[2]);must(err);fmt.Print(strconv.FormatInt(info.Size(),10))
	default: panic("unknown mode")
	}
}
GO
GOCACHE="$temp_root/gocache" "${GO:-go}" build -o "$temp_root/json-helper" "$helper_source" || fail 'could not build JSON helper'
json_helper="$temp_root/json-helper"

run_ok() {
	local label=$1 directory=$2; shift 2
	set +e
	(cd "$directory" && HOME="$temp_root/home" TMPDIR="$temp_root/recipient-temp" "$binary" "$@") >"$temp_root/$label.stdout" 2>"$temp_root/$label.stderr"
	local status=$?
	set -e
	[[ $status -eq 0 ]] || fail "$label exited $status"
	[[ ! -s "$temp_root/$label.stderr" ]] || fail "$label wrote stderr"
}

run_fail() {
	local label=$1 directory=$2 class=$3; shift 3
	set +e
	(cd "$directory" && HOME="$temp_root/home" TMPDIR="$temp_root/recipient-temp" "$binary" "$@") >"$temp_root/$label.stdout" 2>"$temp_root/$label.stderr"
	local status=$?
	set -e
	[[ $status -eq 1 ]] || fail "$label exited $status instead of 1"
	[[ ! -s "$temp_root/$label.stdout" ]] || fail "$label wrote stdout"
	grep -F "$class:" "$temp_root/$label.stderr" >/dev/null || fail "$label did not report $class"
}

work="$temp_root/workspace"
target="$work/target"
mkdir -p "$target"
git -C "$target" init -b main >/dev/null
git -C "$target" config user.name 'WF01 Roundtrip'
git -C "$target" config user.email 'wf01@example.invalid'
printf 'base\n' >"$target/README.md"
git -C "$target" add README.md
git -C "$target" commit -m base >/dev/null
target=$(cd "$target" && pwd -P)
work=$(cd "$work" && pwd -P)
ref=$(git -C "$target" symbolic-ref HEAD)
oid=$(git -C "$target" rev-parse HEAD)
common=$(git -C "$target" rev-parse --path-format=absolute --git-common-dir)
printf 'immutable input\n' >"$work/input.txt"
input="$work/input.txt"

run_ok workspace-init "$work" workspace init
run_ok project-add "$work" workspace project add demo --name Demo --wrapper "$work" --repo "demo=$target"
[[ ! -e "$temp_root/home/.ply" ]] || fail 'workspace prerequisite created a global profile'

make_draft() { "$json_helper" draft "$1" "$target" "$ref" "$oid" "$common" "$input" "$2" "$3" "$4"; }
make_start() { "$json_helper" start "$1" "$2" "$temp_root/recipient-temp"; }
handoff_field() { "$json_helper" field "$1" "$2"; }

draft="$work/happy-draft.json"
make_draft "$draft" wf01/roundtrip/happy wf01/roundtrip/happy 'Roundtrip happy path'
run_ok create "$work" workflow handoff create --file "$draft"
locator=$(sed -n 's/^Handoff: //p' "$temp_root/create.stdout")
[[ -n "$locator" && -f "$locator" ]] || fail 'create did not return an immutable locator'
handoff_id=$(handoff_field "$locator" identity.handoff_id)
printf 'Created agent handoff %s.\nPurpose: Roundtrip happy path\nWorking directory: %s\nHandoff: %s\nNext action: Open a fresh recipient agent in the working directory and tell it: "Read and execute the handoff at %s."\n' "$handoff_id" "$target" "$locator" "$locator" >"$temp_root/create.expected"
cmp -s "$temp_root/create.expected" "$temp_root/create.stdout" || fail 'create output changed'
[[ $(file_mode "$locator") == 400 ]] || fail 'handoff mode is not 0400'
[[ $(file_mode "$(dirname "$locator")") == 700 ]] || fail 'run directory mode is not 0700'
locator_sha=$(file_sha256 "$locator"); locator_mtime=$(file_mtime "$locator")
run_ok create-retry "$work" workflow handoff create --file "$draft"
grep -F "Agent handoff $handoff_id already exists with identical content." "$temp_root/create-retry.stdout" >/dev/null || fail 'create retry was not idempotent'
[[ $(file_sha256 "$locator") == "$locator_sha" && $(file_mtime "$locator") == "$locator_mtime" && $(file_mode "$locator") == 400 ]] || fail 'create retry rewrote the handoff'

run_ok inspect-ready "$work" workflow handoff inspect --handoff "$locator" --format json
grep -F '"derived_state":"ready"' "$temp_root/inspect-ready.stdout" >/dev/null || fail 'ready inspection state is wrong'
grep -F '"secret":"[REDACTED]"' "$temp_root/inspect-ready.stdout" >/dev/null || fail 'inspection did not redact the secret'
secret=$(handoff_field "$locator" reply_capability.secret)
! grep -F "$secret" "$temp_root/inspect-ready.stdout" >/dev/null || fail 'inspection exposed the secret'
run_fail raw-without-ack "$work" workflow_handoff_invalid_arguments workflow handoff inspect --handoff "$locator" --raw handoff
run_ok raw-handoff "$work" workflow handoff inspect --handoff "$locator" --raw handoff --acknowledge-secret-exposure
cmp -s "$locator" "$temp_root/raw-handoff.stdout" || fail 'raw handoff bytes differ from stored bytes'

start_draft="$work/start.json"
make_start "$start_draft" "$locator"
run_ok submit-start "$target" workflow handoff submit-start --handoff "$locator" --file "$start_draft"
start_locator=$(sed -n 's/^Start receipt: //p' "$temp_root/submit-start.stdout")
start_sha=$(sed -n 's/^SHA-256: //p' "$temp_root/submit-start.stdout")
[[ -f "$start_locator" && $(file_sha256 "$start_locator") == "$start_sha" ]] || fail 'accepted start locator/digest is invalid'
start_mtime=$(file_mtime "$start_locator")
run_ok submit-start-retry "$target" workflow handoff submit-start --handoff "$locator" --file "$start_draft"
[[ $(file_mtime "$start_locator") == "$start_mtime" ]] || fail 'start retry rewrote accepted bytes'
run_ok show-started "$work" workflow handoff show "$handoff_id"
grep -Fx 'Status: The expected recipient reported a successful start.' "$temp_root/show-started.stdout" >/dev/null || fail 'started show status changed'

reply_root=$(handoff_field "$locator" reply_capability.reply_root)
artifact_staging="$reply_root/staging/verifier.stdout"
printf 'working tree clean\n' >"$artifact_staging"
result_draft="$work/result.json"
"$json_helper" result "$result_draft" "$locator" "$start_locator" "$artifact_staging" 'The bounded roundtrip is complete.'
run_ok submit-result "$target" workflow handoff submit-result --handoff "$locator" --file "$result_draft"
result_locator=$(sed -n 's/^Terminal result: //p' "$temp_root/submit-result.stdout")
result_sha=$(sed -n 's/^SHA-256: //p' "$temp_root/submit-result.stdout")
[[ -f "$result_locator" && $(file_sha256 "$result_locator") == "$result_sha" ]] || fail 'accepted result locator/digest is invalid'
artifact_sha=$(file_sha256 "$artifact_staging")
artifact_locator="$reply_root/artifacts/sha256/${artifact_sha#sha256:}"
cmp -s "$artifact_staging" "$artifact_locator" || fail 'managed artifact bytes were not preserved'
[[ $(file_mode "$artifact_locator") == 400 ]] || fail 'managed artifact mode is not 0400'
result_mtime=$(file_mtime "$result_locator")
run_ok submit-result-retry "$target" workflow handoff submit-result --handoff "$locator" --file "$result_draft"
[[ $(file_mtime "$result_locator") == "$result_mtime" ]] || fail 'result retry rewrote accepted bytes'
run_ok show-complete "$work" workflow handoff show "$handoff_id"
grep -Fx 'Status: The recipient reported a complete result for the expected run.' "$temp_root/show-complete.stdout" >/dev/null || fail 'complete show status changed'
grep -Fx 'Result: The bounded roundtrip is complete.' "$temp_root/show-complete.stdout" >/dev/null || fail 'complete show result changed'
run_ok inspect-complete "$work" workflow handoff inspect --handoff "$locator" --format json
grep -F '"derived_state":"complete"' "$temp_root/inspect-complete.stdout" >/dev/null || fail 'complete inspection state is wrong'
grep -F '"next_transition_authorized":false' "$temp_root/inspect-complete.stdout" >/dev/null || fail 'inspection opened a downstream transition'
run_ok raw-start "$work" workflow handoff inspect --handoff "$locator" --raw start
run_ok raw-result "$work" workflow handoff inspect --handoff "$locator" --raw result
[[ $(file_sha256 "$temp_root/raw-start.stdout") == "$start_sha" ]] || fail 'raw start digest differs from locator digest'
[[ $(file_sha256 "$temp_root/raw-result.stdout") == "$result_sha" ]] || fail 'raw result digest differs from locator digest'
"$json_helper" result "$work/competing-result.json" "$locator" "$start_locator" "$artifact_staging" 'A competing terminal report.'
run_fail competing-result "$target" workflow_handoff_conflict workflow handoff submit-result --handoff "$locator" --file "$work/competing-result.json"

scenario_create() {
	local name=$1
	local scenario_draft="$work/$name-draft.json"
	make_draft "$scenario_draft" "wf01/roundtrip/$name" "wf01/roundtrip/$name" "Scenario $name"
	run_ok "$name-create" "$work" workflow handoff create --file "$scenario_draft"
	sed -n 's/^Handoff: //p' "$temp_root/$name-create.stdout"
}

drift_locator=$(scenario_create drift)
make_start "$work/drift-start.json" "$drift_locator"
printf 'drift\n' >"$target/README.md"
run_fail drift-start "$target" workflow_handoff_conflict workflow handoff submit-start --handoff "$drift_locator" --file "$work/drift-start.json"
git -C "$target" restore README.md

early_locator=$(scenario_create early)
printf '{}\n' >"$work/early-result.json"
run_fail terminal-before-start "$target" workflow_handoff_start_required workflow handoff submit-result --handoff "$early_locator" --file "$work/early-result.json"

cancel_locator=$(scenario_create cancel)
cancel_id=$(handoff_field "$cancel_locator" identity.handoff_id)
run_ok cancel "$work" workflow handoff cancel "$cancel_id" --reason 'No longer needed'
run_ok cancel-show "$work" workflow handoff show "$cancel_id"
grep -Fx 'Status: The unstarted agent handoff was cancelled.' "$temp_root/cancel-show.stdout" >/dev/null || fail 'cancel did not close the run'

supersede_locator=$(scenario_create supersede)
supersede_id=$(handoff_field "$supersede_locator" identity.handoff_id)
make_draft "$work/supersede-replacement.json" wf01/roundtrip/supersede-replacement wf01/roundtrip/supersede 'Scenario supersede replacement'
run_ok supersede "$work" workflow handoff supersede "$supersede_id" --file "$work/supersede-replacement.json" --reason 'Corrected input binding'
replacement_locator=$(sed -n 's/^Handoff: //p' "$temp_root/supersede.stdout")
[[ -f "$replacement_locator" ]] || fail 'supersede did not publish a replacement'
run_ok supersede-show "$work" workflow handoff show "$supersede_id"
grep -Fx 'Status: The agent handoff was superseded before start.' "$temp_root/supersede-show.stdout" >/dev/null || fail 'supersede did not close the old run'

abandon_locator=$(scenario_create abandon)
abandon_id=$(handoff_field "$abandon_locator" identity.handoff_id)
make_start "$work/abandon-start.json" "$abandon_locator"
run_ok abandon-start "$target" workflow handoff submit-start --handoff "$abandon_locator" --file "$work/abandon-start.json"
run_ok abandon "$work" workflow handoff abandon "$abandon_id" --reason 'Recipient disappeared after start' --acknowledge-effects-unknown
run_ok abandon-show "$work" workflow handoff show "$abandon_id"
grep -Fx 'Status: The started agent handoff was abandoned with target effects unknown.' "$temp_root/abandon-show.stdout" >/dev/null || fail 'abandon did not close the started run'

invalid_locator=$(scenario_create invalid)
printf '{"kind":"invalid"}\n' >"$work/invalid-start.json"
run_fail invalid-start "$target" workflow_handoff_schema_invalid workflow handoff submit-start --handoff "$invalid_locator" --file "$work/invalid-start.json"
invalid_reply=$(handoff_field "$invalid_locator" reply_capability.reply_root)
[[ $(find "$invalid_reply/start/rejected" -type f | wc -l | tr -d ' ') -eq 1 ]] || fail 'schema-invalid input was not bounded and preserved'

conflict_draft="$work/publication-conflict.json"
make_draft "$conflict_draft" wf01/roundtrip/happy wf01/roundtrip/happy 'Different bytes for the same publication'
run_fail publication-conflict "$work" workflow_handoff_conflict workflow handoff create --file "$conflict_draft"

[[ ! -e "$temp_root/home/.ply" ]] || fail 'workflow commands created a global profile'
case "$locator $start_locator $result_locator $artifact_locator $replacement_locator" in "$temp_root"/*) ;; *) fail 'a workflow file escaped the private test root' ;; esac
printf '%s\n' 'workflow handoff roundtrip: PASS'
