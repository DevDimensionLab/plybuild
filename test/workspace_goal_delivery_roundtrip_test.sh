#!/usr/bin/env bash

set -euo pipefail
export LC_ALL=C LANG=C GOENV=off GOWORK=off GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
. "$repo_root/test/roundtrip_root.sh"
roundtrip_root 'workspace goal delivery roundtrip' ply-goal-delivery
temp_root=$run_root
mkdir -p "$temp_root/home" "$temp_root/gocache"
binary="$temp_root/ply"
go_bin=${GO:-go}

fail() { printf 'workspace goal delivery roundtrip: %s\n' "$*" >&2; exit 1; }

(
	cd "$repo_root"
	GOCACHE="$temp_root/gocache" "$go_bin" build -o "$binary" ./cmd/ply
) || fail 'could not build fresh binary'

run_ok() {
	local label=$1 cwd=$2
	shift 2
	(cd "$cwd" && HOME="$temp_root/home" "$binary" "$@") >"$temp_root/$label.stdout" 2>"$temp_root/$label.stderr" ||
		fail "$label failed: $(cat "$temp_root/$label.stderr" "$temp_root/$label.stdout")"
}
field() { python3 -c 'import json,sys; v=json.load(open(sys.argv[1]))
for k in sys.argv[2].split("."): v=v[k]
print(v)' "$1" "$2"; }

workspace="$temp_root/workspace"
repository="$workspace/repo"
epic="$workspace/epic"
task="$workspace/task"
store="$workspace/.ply/work-items.yaml"
mkdir -p "$repository"
git -C "$repository" init -q -b main
git -C "$repository" config user.email fixture@example.invalid
git -C "$repository" config user.name Fixture
printf 'fixture\n' >"$repository/README.md"
git -C "$repository" add README.md
git -C "$repository" commit -qm fixture
base_oid=$(git -C "$repository" rev-parse HEAD)
git -C "$repository" worktree add -q -b epic "$epic"
run_ok init "$workspace" workspace init
run_ok project "$workspace" workspace project add ply --name Ply --wrapper "$workspace" --repo "ply=$repository"
run_ok epic "$workspace" workspace epic adopt epic --title Epic --project ply --repo ply --worktree "$epic" --ref refs/heads/epic --expected-oid "$base_oid"

goal=(workspace task goal register task --spec goal-v1 --title 'Goal' --objective 'Deliver the Goal' --criterion 'The Epic has the result' --project ply --repo ply --epic epic --format json)
run_ok goal-preview "$workspace" "${goal[@]}"
run_ok goal-apply "$workspace" "${goal[@]}" --apply --confirmation "$(field "$temp_root/goal-preview.stdout" confirmation)"
run_ok worktree "$workspace" workspace task worktree create task --branch task --path "$task" --expected-parent-oid "$base_oid"

printf 'result\n' >"$task/result.txt"
git -C "$task" add result.txt
git -C "$task" commit -qm result
result_oid=$(git -C "$task" rev-parse HEAD)

integrate=(workspace task integrate task --expected-result-oid "$result_oid" --expected-parent-oid "$base_oid" --format json)
store_sha=$(git hash-object "$store")
run_ok check "$workspace" "${integrate[@]}" --check
[[ $(field "$temp_root/check.stdout" classification) == ready ]] || fail 'Goal delivery check was not ready'
[[ $(git hash-object "$store") == "$store_sha" && $(git -C "$epic" rev-parse HEAD) == "$base_oid" ]] || fail 'check had an effect'
confirm=$(field "$temp_root/check.stdout" plan.sha256)
run_ok apply "$workspace" "${integrate[@]}" --apply --confirm "$confirm"
[[ $(field "$temp_root/apply.stdout" classification) == exact_effect ]] || fail 'apply was not exact_effect'
[[ $(git -C "$epic" rev-parse HEAD) == "$result_oid" ]] || fail 'Epic did not reach the exact result'
[[ $(git -C "$repository" rev-parse refs/heads/main) == "$base_oid" ]] || fail 'ordinary main changed'
! grep -q 'task_results: \[[^]]' "$store" || fail 'Goal delivery wrote a TaskResult'
store_sha=$(git hash-object "$store")
run_ok replay "$workspace" "${integrate[@]}" --apply --confirm "$confirm"
[[ $(field "$temp_root/replay.stdout" classification) == exact_effect && $(git hash-object "$store") == "$store_sha" ]] || fail 'replay was not idempotent'
[[ $(git -C "$epic" reflog show --format=%H refs/heads/epic) == "$result_oid"$'\n'"$base_oid" ]] || fail 'reflog lacks one exact move'
run_ok show "$workspace" workspace task show task --format json
grep -F '"classification":"exact_effect"' "$temp_root/show.stdout" >/dev/null || fail 'Task show lost the delivery'
run_ok trace "$workspace" workflow trace task --format json
! grep -F '"result_id":""' "$temp_root/trace.stdout" >/dev/null || fail 'trace invented an empty result identity'
grep -F "$result_oid" "$temp_root/trace.stdout" >/dev/null || fail 'trace lost the delivered candidate'
[[ ! -e "$temp_root/home/.ply" ]] || fail 'blackbox wrote global HOME state'

printf 'workspace goal delivery roundtrip: PASS\n'
