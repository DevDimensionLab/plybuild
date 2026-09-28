#!/usr/bin/env bash

set -euo pipefail
export LC_ALL=C LANG=C GOENV=off GOWORK=off GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
temp_root=$(mktemp -d "${TMPDIR:-/tmp}/ply-work-items.XXXXXX")
cleanup() { find "$temp_root" -type d -exec chmod u+rwx {} + 2>/dev/null || true; rm -rf -- "$temp_root"; }
trap cleanup EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM
mkdir -p "$temp_root/home" "$temp_root/gocache" "$temp_root/workspace/ply/main"
binary="$temp_root/ply"
go_bin=${GO:-go}

fail() { printf 'workspace work-items roundtrip: %s\n' "$*" >&2; exit 1; }
run() { HOME="$temp_root/home" "$binary" "$@"; }
git_local() { git -c user.name='Ply tests' -c user.email=tests@example.invalid "$@"; }

(
	cd "$repo_root"
	GOCACHE="$temp_root/gocache" "$go_bin" build -o "$binary" ./cmd/ply
) || fail 'could not build fresh binary'

workspace="$temp_root/workspace"
repository="$workspace/ply/main"
epic="$workspace/ply/epic"
task_worktree="$workspace/ply/task"
(
	cd "$workspace"
	run workspace init >/dev/null
)
git_local -C "$repository" init -b main >/dev/null
printf 'fixture\n' >"$repository/README.md"
git_local -C "$repository" add README.md
git_local -C "$repository" commit -m fixture >/dev/null
oid=$(git -C "$repository" rev-parse HEAD)
tree=$(git -C "$repository" rev-parse 'HEAD^{tree}')
git -C "$repository" worktree add -b epic "$epic" "$oid" >/dev/null
(
	cd "$workspace"
	run workspace project add ply --name Ply --wrapper "$workspace/ply" --repo "ply=$repository" >/dev/null
)
marker_before=$(shasum -a 256 "$workspace/.ply/workspace.yaml")
projects_before=$(shasum -a 256 "$workspace/.ply/projects.yaml")

(
	cd "$workspace"
	run workspace epic adopt epic --title Epic --project ply --repo ply --worktree "$epic" --ref refs/heads/epic --expected-oid "$oid"
	run workspace task create task --title Task --description 'Task description' --epic epic --project ply --repo ply
	run workspace task worktree create task --branch task --path "$task_worktree" --expected-parent-oid "$oid"
	run workspace epic show epic --format json
	run workspace task show task --format json
	run workspace epic list
	run workspace task list --epic epic
) >"$temp_root/journey.stdout" 2>"$temp_root/journey.stderr" || fail 'happy path failed'
[[ ! -s "$temp_root/journey.stderr" ]] || fail 'happy path wrote stderr'
grep -F 'Adopted Epic epic: Epic' "$temp_root/journey.stdout" >/dev/null || fail 'adopt output missing'
grep -F 'Created Task task: Task' "$temp_root/journey.stdout" >/dev/null || fail 'Task output missing'
grep -F 'Created Task worktree for task.' "$temp_root/journey.stdout" >/dev/null || fail 'worktree output missing'
grep -F '"kind":"WorkspaceEpicReadback@1"' "$temp_root/journey.stdout" >/dev/null || fail 'Epic JSON missing'
grep -F '"kind":"WorkspaceTaskReadback@1"' "$temp_root/journey.stdout" >/dev/null || fail 'Task JSON missing'
grep -F '"ready_for_handoff":true' "$temp_root/journey.stdout" >/dev/null || fail 'Task is not handoff-ready'

[[ $(git -C "$epic" rev-parse HEAD) == "$oid" ]] || fail 'Epic HEAD changed'
[[ $(git -C "$epic" rev-parse 'HEAD^{tree}') == "$tree" ]] || fail 'Epic tree changed'
[[ -z $(git -C "$epic" status --porcelain=v1 --untracked-files=all) ]] || fail 'Epic became dirty'
[[ $(git -C "$task_worktree" rev-parse HEAD) == "$oid" ]] || fail 'Task worktree does not use pinned base'
[[ -z $(git -C "$task_worktree" status --porcelain=v1 --untracked-files=all) ]] || fail 'Task worktree is dirty'
[[ $(shasum -a 256 "$workspace/.ply/workspace.yaml") == "$marker_before" ]] || fail 'workspace marker changed'
[[ $(shasum -a 256 "$workspace/.ply/projects.yaml") == "$projects_before" ]] || fail 'Project registry changed'

store="$workspace/.ply/work-items.yaml"
grep -F 'format_version: 1' "$store" >/dev/null || fail 'work-item format version missing'
grep -F 'worktree_state: worktree_ready' "$store" >/dev/null || fail 'ready state missing'
grep -F 'classification: exact_effect' "$store" >/dev/null || fail 'exact-effect observation missing'
touch -t 200001010000 "$store"
store_before=$(shasum -a 256 "$store")
if [[ $(uname -s) == Darwin ]]; then
	store_mtime_before=$(stat -f '%m' "$store")
else
	store_mtime_before=$(stat -c '%Y' "$store")
fi
worktree_count_before=$(git -C "$repository" worktree list --porcelain | grep -c '^worktree ')
(
	cd "$workspace"
	run workspace epic adopt epic --title Epic --project ply --repo ply --worktree "$epic" --ref refs/heads/epic --expected-oid "$oid" >/dev/null
	run workspace task create task --title Task --description 'Task description' --epic epic --project ply --repo ply >/dev/null
	run workspace task worktree create task --branch task --path "$task_worktree" --expected-parent-oid "$oid" >"$temp_root/retry.stdout"
)
grep -F 'Task worktree for task is already ready.' "$temp_root/retry.stdout" >/dev/null || fail 'retry was not idempotent'
[[ $(shasum -a 256 "$store") == "$store_before" ]] || fail 'retry changed store bytes'
[[ $(git -C "$repository" worktree list --porcelain | grep -c '^worktree ') == "$worktree_count_before" ]] || fail 'retry added another worktree'
if [[ $(uname -s) == Darwin ]]; then
	[[ $(stat -f '%m' "$store") == "$store_mtime_before" ]] || fail 'retry changed store mtime'
else
	[[ $(stat -c '%Y' "$store") == "$store_mtime_before" ]] || fail 'retry changed store mtime'
fi

set +e
(
	cd "$workspace"
	run workspace task show missing
) >"$temp_root/error.stdout" 2>"$temp_root/error.stderr"
status=$?
set -e
[[ $status -eq 1 && ! -s "$temp_root/error.stdout" ]] || fail 'missing Task stream/exit contract failed'
grep -F 'Error: workspace_work_not_found:' "$temp_root/error.stderr" >/dev/null || fail 'missing Task class changed'

[[ ! -e "$temp_root/home/.ply" ]] || fail 'command wrote global HOME state'
printf 'workspace work-items roundtrip: PASS\n'
