#!/usr/bin/env bash

set -euo pipefail
export LC_ALL=C LANG=C GOENV=off GOWORK=off

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
. "$repo_root/test/roundtrip_root.sh"
roundtrip_root 'cli surface contract' ply-cli-surface
temp_root=$run_root
go_bin=${GO:-go}
binary="$temp_root/ply"
mkdir -p "$temp_root/gocache" "$temp_root/home"

fail() {
	local file
	printf 'cli surface contract: %s\n' "$1" >&2
	for file in "${@:2}"; do
		[[ -f $file ]] || continue
		printf '== %s ==\n' "$file" >&2 || true
		sed -n '1,80p' "$file" >&2 || true
	done
	exit 1
}

if (
	cd "$repo_root"
	GOCACHE="$temp_root/gocache" "$go_bin" build -o "$binary" ./cmd/ply
) >"$temp_root/build.stdout" 2>"$temp_root/build.stderr"; then
	:
else
	build_exit=$?
	cat "$temp_root/build.stderr" >&2 || true
	printf 'cli surface contract: could not build fresh CLI artifact\n' >&2 || true
	exit "$build_exit"
fi

for surface in root status upgrade build; do
	case $surface in
root) arguments=(--help); expected='== version: dev ==' ;;
	status) arguments=(status --help); expected='Status functionality for a project' ;;
	upgrade) arguments=(upgrade --help); expected='Perform upgrade on existing projects' ;;
	build) arguments=(build --help); expected='Builds a ply project' ;;
	esac
	HOME="$temp_root/home" "$binary" "${arguments[@]}" \
		>"$temp_root/$surface.stdout" 2>"$temp_root/$surface.stderr" ||
		fail "$surface help returned non-zero" "$temp_root/home" "$temp_root/$surface.stdout" "$temp_root/$surface.stderr"
	grep -F "$expected" "$temp_root/$surface.stdout" >/dev/null ||
		fail "$surface help lost its expected text" "$temp_root/$surface.stdout"
	[[ ! -s "$temp_root/$surface.stderr" ]] || fail "$surface help wrote to stderr" "$temp_root/$surface.stderr"
done

run_help_surface() {
	local label=$1
	shift
	set +e
	HOME="$temp_root/home" "$binary" "$@" \
		>"$temp_root/$label.stdout" 2>"$temp_root/$label.stderr"
	local exit_status=$?
	set -e
	printf '%s\n' "$exit_status" >"$temp_root/$label.exit"
	[[ "$exit_status" -eq 0 ]] || fail "$label help exit is $exit_status, expected 0" "$temp_root/$label.stdout" "$temp_root/$label.stderr"
	[[ ! -s "$temp_root/$label.stderr" ]] || fail "$label help wrote to stderr" "$temp_root/$label.stderr"
}

run_help_surface root-help --help
run_help_surface capabilities-help capabilities --help
run_help_surface workspace-help workspace --help
run_help_surface workspace-init-help workspace init --help
run_help_surface workspace-project-help workspace project --help
run_help_surface workspace-project-add-help workspace project add --help
run_help_surface workspace-project-show-help workspace project show --help
run_help_surface workspace-project-list-help workspace project list --help
run_help_surface workspace-epic-help workspace epic --help
run_help_surface workspace-epic-adopt-help workspace epic adopt --help
run_help_surface workspace-epic-show-help workspace epic show --help
run_help_surface workspace-epic-list-help workspace epic list --help
run_help_surface workspace-task-help workspace task --help
run_help_surface workspace-task-journal-help workspace task journal --help
for journal_leaf in show record observe; do
	run_help_surface "workspace-task-journal-$journal_leaf-help" workspace task journal "$journal_leaf" --help
done
grep -F 'journal' "$temp_root/workspace-task-help.stdout" >/dev/null || fail 'task help lost journal' "$temp_root/workspace-task-help.stdout"
grep -F 'immutable' "$temp_root/workspace-task-journal-record-help.stdout" >/dev/null || fail 'journal record help lost immutable contract' "$temp_root/workspace-task-journal-record-help.stdout"
grep -F -- '--snapshot' "$temp_root/workspace-task-journal-show-help.stdout" >/dev/null || fail 'journal show help lost offline snapshot' "$temp_root/workspace-task-journal-show-help.stdout"
grep -F -- '--file' "$temp_root/workspace-task-journal-observe-help.stdout" >/dev/null || fail 'journal observe help lost file input' "$temp_root/workspace-task-journal-observe-help.stdout"
run_help_surface workspace-task-create-help workspace task create --help
run_help_surface workspace-task-show-help workspace task show --help
run_help_surface workspace-task-list-help workspace task list --help
run_help_surface workspace-task-worktree-help workspace task worktree --help
run_help_surface workspace-task-worktree-create-help workspace task worktree create --help
run_help_surface workspace-task-integrate-help workspace task integrate --help
run_help_surface workflow-help workflow --help

grep -F '  workspace    Manage Ply workspaces' "$temp_root/root-help.stdout" >/dev/null ||
	fail 'root help does not expose the workspace parent' "$temp_root/root-help.stdout"
grep -F '  workflow     Manage agent workflows' "$temp_root/root-help.stdout" >/dev/null ||
	fail 'root help does not expose the workflow parent' "$temp_root/root-help.stdout"
grep -F '  capabilities Show supported workspace core read contracts' "$temp_root/root-help.stdout" >/dev/null ||
	fail 'root help does not expose capabilities' "$temp_root/root-help.stdout"
for label in capabilities-help workspace-project-list-help workspace-project-show-help workspace-task-list-help; do
	grep -F -- '--format string' "$temp_root/$label.stdout" >/dev/null || fail "$label lost its output format" "$temp_root/$label.stdout"
done
grep -F 'Manage explicit local Ply workspaces.' "$temp_root/workspace-help.stdout" >/dev/null ||
	fail 'workspace help lost its long description' "$temp_root/workspace-help.stdout"
grep -F 'ply workspace [command]' "$temp_root/workspace-help.stdout" >/dev/null ||
	fail 'workspace help usage changed' "$temp_root/workspace-help.stdout"
grep -F '  init        Initialize a Ply workspace in the current directory' \
	"$temp_root/workspace-help.stdout" >/dev/null || fail 'workspace help does not expose init' "$temp_root/workspace-help.stdout"
grep -F '  project     Manage projects in a Ply workspace' \
	"$temp_root/workspace-help.stdout" >/dev/null || fail 'workspace help does not expose project' "$temp_root/workspace-help.stdout"
grep -F '  epic        Manage Epics in a Ply workspace' \
	"$temp_root/workspace-help.stdout" >/dev/null || fail 'workspace help does not expose Epic' "$temp_root/workspace-help.stdout"
grep -F '  task        Manage Tasks in a Ply workspace' \
	"$temp_root/workspace-help.stdout" >/dev/null || fail 'workspace help does not expose Task' "$temp_root/workspace-help.stdout"
grep -F '  ply workspace init' "$temp_root/workspace-help.stdout" >/dev/null ||
	fail 'workspace help lost its example' "$temp_root/workspace-help.stdout"
grep -F '  ply workspace project list' "$temp_root/workspace-help.stdout" >/dev/null ||
	fail 'workspace help lost its project example' "$temp_root/workspace-help.stdout"
grep -F 'Initialize a Ply workspace in the current directory by creating .ply/workspace.yaml.' \
	"$temp_root/workspace-init-help.stdout" >/dev/null || fail 'workspace init help lost its action' "$temp_root/workspace-init-help.stdout"
grep -F 'The current directory does not need to be a Git repository.' \
	"$temp_root/workspace-init-help.stdout" >/dev/null || fail 'workspace init help lost its Git boundary' "$temp_root/workspace-init-help.stdout"
grep -F 'ply workspace init [flags]' "$temp_root/workspace-init-help.stdout" >/dev/null ||
	fail 'workspace init help usage changed' "$temp_root/workspace-init-help.stdout"
grep -F '  ply workspace init' "$temp_root/workspace-init-help.stdout" >/dev/null ||
	fail 'workspace init help lost its example' "$temp_root/workspace-init-help.stdout"
grep -F -- '-h, --help' "$temp_root/workspace-init-help.stdout" >/dev/null ||
	fail 'workspace init help flag is unavailable' "$temp_root/workspace-init-help.stdout"
if grep -F -- '--name' "$temp_root/workspace-init-help.stdout" >/dev/null; then
	fail 'workspace init help exposes --name' "$temp_root/workspace-init-help.stdout"
fi
grep -F 'Register and inspect projects owned by the containing Ply workspace.' \
	"$temp_root/workspace-project-help.stdout" >/dev/null || fail 'project help lost its long description' "$temp_root/workspace-project-help.stdout"
grep -F 'ply workspace project [command]' "$temp_root/workspace-project-help.stdout" >/dev/null ||
	fail 'project help usage changed' "$temp_root/workspace-project-help.stdout"
grep -E '^  add +Register a project and its explicit repository members$' \
	"$temp_root/workspace-project-help.stdout" >/dev/null || fail 'project help does not expose add' "$temp_root/workspace-project-help.stdout"
grep -E '^  list +List registered projects$' "$temp_root/workspace-project-help.stdout" >/dev/null ||
	fail 'project help does not expose list' "$temp_root/workspace-project-help.stdout"
grep -E '^  show +Show a registered project$' "$temp_root/workspace-project-help.stdout" >/dev/null ||
	fail 'project help does not expose show' "$temp_root/workspace-project-help.stdout"
grep -F 'ply workspace project add <project-id> [flags]' \
	"$temp_root/workspace-project-add-help.stdout" >/dev/null || fail 'project add usage changed' "$temp_root/workspace-project-add-help.stdout"
grep -F -- '-n, --name string' "$temp_root/workspace-project-add-help.stdout" >/dev/null ||
	fail 'project add name flag is unavailable' "$temp_root/workspace-project-add-help.stdout"
grep -F -- '-w, --wrapper string' "$temp_root/workspace-project-add-help.stdout" >/dev/null ||
	fail 'project add wrapper flag is unavailable' "$temp_root/workspace-project-add-help.stdout"
grep -F -- '-r, --repo stringArray' "$temp_root/workspace-project-add-help.stdout" >/dev/null ||
	fail 'project add repo flag is not a repeatable string array' "$temp_root/workspace-project-add-help.stdout"
grep -F 'ply workspace project show <project-id> [flags]' \
	"$temp_root/workspace-project-show-help.stdout" >/dev/null || fail 'project show usage changed' "$temp_root/workspace-project-show-help.stdout"
grep -F 'ply workspace project list [flags]' \
	"$temp_root/workspace-project-list-help.stdout" >/dev/null || fail 'project list usage changed' "$temp_root/workspace-project-list-help.stdout"
grep -F 'ply workspace epic [command]' "$temp_root/workspace-epic-help.stdout" >/dev/null || fail 'Epic usage changed' "$temp_root/workspace-epic-help.stdout"
grep -F 'ply workspace epic adopt <epic-id> [flags]' "$temp_root/workspace-epic-adopt-help.stdout" >/dev/null || fail 'Epic adopt usage changed' "$temp_root/workspace-epic-adopt-help.stdout"
for flag in '--title string' '--project string' '--repo string' '--worktree string' '--ref string' '--expected-oid string'; do
	grep -F -- "$flag" "$temp_root/workspace-epic-adopt-help.stdout" >/dev/null || fail "Epic adopt $flag is unavailable" "$temp_root/workspace-epic-adopt-help.stdout"
done
grep -F 'ply workspace epic show <epic-id> [flags]' "$temp_root/workspace-epic-show-help.stdout" >/dev/null || fail 'Epic show usage changed' "$temp_root/workspace-epic-show-help.stdout"
grep -F -- '--format string' "$temp_root/workspace-epic-show-help.stdout" >/dev/null || fail 'Epic show format is unavailable' "$temp_root/workspace-epic-show-help.stdout"
grep -F 'ply workspace epic list [flags]' "$temp_root/workspace-epic-list-help.stdout" >/dev/null || fail 'Epic list usage changed' "$temp_root/workspace-epic-list-help.stdout"
grep -F 'ply workspace task [command]' "$temp_root/workspace-task-help.stdout" >/dev/null || fail 'Task usage changed' "$temp_root/workspace-task-help.stdout"
for leaf in create integrate list problem publication show spec worktree; do
	grep -F "  $leaf" "$temp_root/workspace-task-help.stdout" >/dev/null || fail "Task help does not expose $leaf" "$temp_root/workspace-task-help.stdout"
done
grep -F 'ply workspace task create <task-id> [flags]' "$temp_root/workspace-task-create-help.stdout" >/dev/null || fail 'Task create usage changed' "$temp_root/workspace-task-create-help.stdout"
for flag in '--title string' '--description string' '--epic string' '--project string' '--repo string' '--format string' '--upgrade-store string'; do
	grep -F -- "$flag" "$temp_root/workspace-task-create-help.stdout" >/dev/null || fail "Task create $flag is unavailable" "$temp_root/workspace-task-create-help.stdout"
done
grep -F 'ply workspace task show <task-id> [flags]' "$temp_root/workspace-task-show-help.stdout" >/dev/null || fail 'Task show usage changed' "$temp_root/workspace-task-show-help.stdout"
grep -F 'ply workspace task list [flags]' "$temp_root/workspace-task-list-help.stdout" >/dev/null || fail 'Task list usage changed' "$temp_root/workspace-task-list-help.stdout"
grep -F 'ply workspace task worktree [command]' "$temp_root/workspace-task-worktree-help.stdout" >/dev/null || fail 'Task worktree usage changed' "$temp_root/workspace-task-worktree-help.stdout"
grep -F 'ply workspace task worktree create <task-id> [flags]' "$temp_root/workspace-task-worktree-create-help.stdout" >/dev/null || fail 'Task worktree create usage changed' "$temp_root/workspace-task-worktree-create-help.stdout"
for flag in '--branch string' '--path string' '--expected-parent-oid string'; do
	grep -F -- "$flag" "$temp_root/workspace-task-worktree-create-help.stdout" >/dev/null || fail "Task worktree create $flag is unavailable" "$temp_root/workspace-task-worktree-create-help.stdout"
done
grep -F 'ply workspace task integrate <task-id> [flags]' "$temp_root/workspace-task-integrate-help.stdout" >/dev/null || fail 'Task integrate usage changed' "$temp_root/workspace-task-integrate-help.stdout"
for flag in '--expected-result-oid string' '--expected-parent-oid string' '--check' '--apply' '--confirm string' '--retry-after string' '--format string'; do
	grep -F -- "$flag" "$temp_root/workspace-task-integrate-help.stdout" >/dev/null || fail "Task integrate $flag is unavailable" "$temp_root/workspace-task-integrate-help.stdout"
done

grep -F "Read registered workspace workflow status, trace one Task's recorded history, create local planning repositories, start an explicitly assigned agent, and manage notifications." "$temp_root/workflow-help.stdout" >/dev/null ||
	fail 'workflow help lost its long description' "$temp_root/workflow-help.stdout"
grep -F 'ply workflow [command]' "$temp_root/workflow-help.stdout" >/dev/null ||
	fail 'workflow help usage changed' "$temp_root/workflow-help.stdout"
[[ ! -e "$temp_root/home/.ply" ]] || fail 'help created a global Ply profile' "$temp_root/home/.ply"

grep -F '  workspace    Manage Ply workspaces' "$repo_root/README.md" >/dev/null ||
	fail 'README root command overview does not expose workspace' "$repo_root/README.md"
grep -F '  workflow     Manage agent workflows' "$repo_root/README.md" >/dev/null ||
	fail 'README root command overview does not expose workflow' "$repo_root/README.md"
printf '%s\n' \
	'## Workspace' \
	'Initialize the current directory as an explicit Ply workspace:' \
	'' \
	'```shell script' \
	'ply workspace init' \
	'```' \
	'' \
	'The command creates `.ply/workspace.yaml` with format version 1 and the canonical physical' \
	'directory path. The directory does not need to be a Git repository. Re-running the command is' \
	'safe and leaves an existing compatible marker unchanged. It does not create a Git repository or' \
	'create workflows.' \
	'' \
	'From an initialized workspace, or any directory below it, register a project with an explicit' \
	'wrapper and one or more Git worktree roots:' \
	'' \
	'```shell script' \
	'ply workspace project add ply \' \
	'  --name Ply \' \
	'  --wrapper ../ply \' \
	'  --repo ply=../ply/main' \
	'```' \
	'' \
	'Repeat `--repo` to register a multi-repository project:' \
	'' \
	'```shell script' \
	'ply workspace project add trip \' \
	'  --name Trip \' \
	'  --wrapper ../trip \' \
	'  --repo trip-frontend=../trip/trip-frontend/main \' \
	'  --repo trip-openapi=../trip/trip-openapi/main \' \
	'  --repo trip-service=../trip/trip-service/main' \
	'```' \
	'' \
	'The wrapper and repository members are explicit and may be outside the workspace. Ply validates' \
	'only the nominated worktree roots. Dirty repositories are accepted, and registration performs no' \
	'discovery or Git mutation. Read registrations with `ply workspace project show <id>` and' \
	'`ply workspace project list`. These commands also accept `--format json`.' \
	'See the [workspace read contract](docs/read-contract.md) for independent Task filters,' \
	'versioned JSON schemas, and `ply capabilities`.' \
	'' \
	'Retrying the same registration is idempotent. Changing membership, relocating repositories, and' \
	'`project init` are not part of this command.' \
	'' >"$temp_root/readme-workspace.expected"
sed -n '/^## Workspace$/,/^## Agent quota$/p' "$repo_root/README.md" | sed '$d' \
	>"$temp_root/readme-workspace.actual"
expected_workspace_bytes=$(wc -c <"$temp_root/readme-workspace.expected" | tr -d ' ')
dd if="$temp_root/readme-workspace.actual" of="$temp_root/readme-workspace.prefix" bs=1 count="$expected_workspace_bytes" 2>/dev/null
cmp -s "$temp_root/readme-workspace.expected" "$temp_root/readme-workspace.prefix" ||
	fail 'README existing workspace documentation changed' "$temp_root/readme-workspace.expected" "$temp_root/readme-workspace.prefix"
for text in \
	'ply workspace epic adopt ply-agentic-workflow-support \' \
	'It performs no Git change' \
	'ply workspace task create workspace-work-item-bootstrap \' \
	'binds exactly one registered' \
	'ply workspace task worktree create workspace-work-item-bootstrap \' \
	'--expected-parent-oid 54f3631cbea789f25a4134945c7ca16d343139df' \
	'Ply persists a durable create intent before the additive branch/worktree operation.' \
	'are never reset, removed,' \
	'`ply workspace task show <id>`' \
	'confirmation digest' \
	'exactly one local `ff-only`' \
	'Git 2.45.0 or newer' \
	'Both show commands accept `--format json`' \
	'it does not start an agent, select a workflow, or grant execution authority.' \
	'Sub-tasks, additional repository anchors on an Epic, editing, rebinding, cleanup, and'; do
	grep -F -- "$text" "$temp_root/readme-workspace.actual" >/dev/null || fail "README work-item text is missing: $text" "$temp_root/readme-workspace.actual"
done
set +e
HOME="$temp_root/home" "$binary" definitely-not-a-command \
	>"$temp_root/invalid.stdout" 2>"$temp_root/invalid.stderr"
invalid_exit=$?
set -e
[[ "$invalid_exit" -eq 1 ]] || fail "unknown command exit is $invalid_exit, expected 1" "$temp_root/invalid.stdout" "$temp_root/invalid.stderr"
grep -F 'unknown command "definitely-not-a-command" for "ply"' \
	"$temp_root/invalid.stderr" >/dev/null || fail 'unknown-command stderr changed' "$temp_root/invalid.stderr"
grep -F "Run 'ply --help' for usage." "$temp_root/invalid.stderr" >/dev/null ||
	fail 'unknown-command usage hint changed' "$temp_root/invalid.stderr"
[[ ! -s "$temp_root/invalid.stdout" ]] || fail 'unknown command wrote to stdout' "$temp_root/invalid.stdout"

file_mode() {
	if [[ $(uname -s) == Darwin ]]; then
		stat -f '%Lp' "$1"
	else
		stat -c '%a' "$1"
	fi
}

file_mtime() {
	if [[ $(uname -s) == Darwin ]]; then
		stat -f '%m' "$1"
	else
		stat -c '%Y' "$1"
	fi
}

work="$temp_root/work"
mkdir "$work"
[[ ! -e "$work/.git" ]] || fail 'fresh workspace unexpectedly contains .git' "$work/.git"
work_root=$(cd "$work" && pwd -P)
(
	cd "$work"
	HOME="$temp_root/home" "$binary" workspace init
) >"$temp_root/work-first.stdout" 2>"$temp_root/work-first.stderr" ||
	fail 'first workspace init returned non-zero' "$temp_root/work-first.stdout" "$temp_root/work-first.stderr"
printf 'Initialized Ply workspace at %s.\n' "$work_root" >"$temp_root/work-first.expected"
cmp -s "$temp_root/work-first.expected" "$temp_root/work-first.stdout" ||
	fail 'first workspace init stdout changed' "$temp_root/work-first.expected" "$temp_root/work-first.stdout"
[[ ! -s "$temp_root/work-first.stderr" ]] || fail 'first workspace init wrote to stderr' "$temp_root/work-first.stderr"
printf 'format_version: 1\nroot: %s\n' "$work_root" >"$temp_root/marker.expected"
cmp -s "$temp_root/marker.expected" "$work/.ply/workspace.yaml" ||
	fail 'workspace marker bytes changed' "$temp_root/marker.expected" "$work/.ply/workspace.yaml"
[[ $(file_mode "$work/.ply") == 755 ]] || fail 'workspace directory mode is not 0755' "$work/.ply"
[[ $(file_mode "$work/.ply/workspace.yaml") == 644 ]] || fail 'workspace marker mode is not 0644' "$work/.ply/workspace.yaml"
touch -t 200001010000 "$work/.ply/workspace.yaml"
first_mtime=$(file_mtime "$work/.ply/workspace.yaml")
cp "$work/.ply/workspace.yaml" "$temp_root/marker.before"
(
	cd "$work"
	HOME="$temp_root/home" "$binary" workspace init
) >"$temp_root/work-second.stdout" 2>"$temp_root/work-second.stderr" ||
	fail 'second workspace init returned non-zero' "$temp_root/work-second.stdout" "$temp_root/work-second.stderr"
printf 'Ply workspace already initialized at %s.\n' "$work_root" >"$temp_root/work-second.expected"
cmp -s "$temp_root/work-second.expected" "$temp_root/work-second.stdout" ||
	fail 'second workspace init stdout changed' "$temp_root/work-second.expected" "$temp_root/work-second.stdout"
[[ ! -s "$temp_root/work-second.stderr" ]] || fail 'second workspace init wrote to stderr' "$temp_root/work-second.stderr"
cmp -s "$temp_root/marker.before" "$work/.ply/workspace.yaml" ||
	fail 'idempotent init changed marker bytes' "$temp_root/marker.before" "$work/.ply/workspace.yaml"
[[ $(file_mtime "$work/.ply/workspace.yaml") == "$first_mtime" ]] ||
	fail 'idempotent init changed marker mtime' "$work/.ply/workspace.yaml"
[[ ! -e "$temp_root/home/.ply" ]] || fail 'workspace init created a global Ply profile' "$temp_root/home/.ply"
[[ $(cd "$work" && find . -mindepth 1 -print | LC_ALL=C sort) == $'./.ply\n./.ply/workspace.yaml' ]] ||
	fail 'workspace init created unexpected persistent files'

physical="$temp_root/physical"
alias="$temp_root/physical-link"
mkdir "$physical"
ln -s "$physical" "$alias"
physical_root=$(cd "$physical" && pwd -P)
(
	cd "$alias"
	HOME="$temp_root/home" "$binary" workspace init
) >"$temp_root/symlink.stdout" 2>"$temp_root/symlink.stderr" ||
	fail 'workspace init through symlink returned non-zero' "$temp_root/symlink.stdout" "$temp_root/symlink.stderr"
printf 'Initialized Ply workspace at %s.\n' "$physical_root" >"$temp_root/symlink.expected"
cmp -s "$temp_root/symlink.expected" "$temp_root/symlink.stdout" ||
	fail 'symlink init did not report the physical root' "$temp_root/symlink.expected" "$temp_root/symlink.stdout"
printf 'format_version: 1\nroot: %s\n' "$physical_root" >"$temp_root/symlink-marker.expected"
cmp -s "$temp_root/symlink-marker.expected" "$physical/.ply/workspace.yaml" ||
	fail 'symlink init marker did not use the physical root' "$temp_root/symlink-marker.expected" "$physical/.ply/workspace.yaml"
[[ ! -s "$temp_root/symlink.stderr" ]] || fail 'symlink init wrote to stderr' "$temp_root/symlink.stderr"

nested_parent="$temp_root/nested-parent"
nested_child="$nested_parent/child"
mkdir -p "$nested_child"
nested_parent_root=$(cd "$nested_parent" && pwd -P)
nested_child_root=$(cd "$nested_child" && pwd -P)
(
	cd "$nested_parent"
	HOME="$temp_root/home" "$binary" workspace init
) >"$temp_root/nested-parent.stdout" 2>"$temp_root/nested-parent.stderr" ||
	fail 'nested parent init returned non-zero' "$temp_root/nested-parent.stdout" "$temp_root/nested-parent.stderr"
set +e
(
	cd "$nested_child"
	HOME="$temp_root/home" "$binary" workspace init
) >"$temp_root/nested.stdout" 2>"$temp_root/nested.stderr"
nested_exit=$?
set -e
[[ "$nested_exit" -eq 1 ]] || fail "nested workspace exit is $nested_exit, expected 1" "$temp_root/nested.stdout" "$temp_root/nested.stderr"
[[ ! -s "$temp_root/nested.stdout" ]] || fail 'nested workspace wrote to stdout' "$temp_root/nested.stdout"
printf 'Error: workspace_init_nested: %s is inside Ply workspace %s\n' \
	"$nested_child_root" "$nested_parent_root" >"$temp_root/nested.expected"
cmp -s "$temp_root/nested.expected" "$temp_root/nested.stderr" ||
	fail 'nested workspace stderr changed' "$temp_root/nested.expected" "$temp_root/nested.stderr"
[[ ! -e "$nested_child/.ply" ]] || fail 'nested workspace created child .ply' "$nested_child/.ply"

compatible="$temp_root/compatible"
mkdir -p "$compatible/.ply"
compatible_root=$(cd "$compatible" && pwd -P)
printf '# retained\nroot: "%s"\nformat_version: 1\n' "$compatible_root" \
	>"$compatible/.ply/workspace.yaml"
(
	cd "$compatible"
	HOME="$temp_root/home" "$binary" workspace init
) >"$temp_root/compatible.stdout" 2>"$temp_root/compatible.stderr" ||
	fail 'compatible workspace init returned non-zero' "$temp_root/compatible.stdout" "$temp_root/compatible.stderr"
printf 'Ply workspace already initialized at %s.\n' "$compatible_root" \
	>"$temp_root/compatible.expected"
cmp -s "$temp_root/compatible.expected" "$temp_root/compatible.stdout" ||
	fail 'compatible workspace stdout changed' "$temp_root/compatible.expected" "$temp_root/compatible.stdout"
[[ ! -s "$temp_root/compatible.stderr" ]] || fail 'compatible workspace wrote to stderr' "$temp_root/compatible.stderr"

conflict="$temp_root/conflict"
mkdir -p "$conflict/.ply"
conflict_root=$(cd "$conflict" && pwd -P)
set +e
(
	cd "$conflict"
	HOME="$temp_root/home" "$binary" workspace init
) >"$temp_root/conflict.stdout" 2>"$temp_root/conflict.stderr"
conflict_exit=$?
set -e
[[ "$conflict_exit" -eq 1 ]] || fail "conflict exit is $conflict_exit, expected 1" "$temp_root/conflict.stdout" "$temp_root/conflict.stderr"
[[ ! -s "$temp_root/conflict.stdout" ]] || fail 'conflict wrote to stdout' "$temp_root/conflict.stdout"
printf 'Error: workspace_init_conflict: %s/.ply/workspace.yaml: workspace marker is missing\n' \
	"$conflict_root" >"$temp_root/conflict.expected"
cmp -s "$temp_root/conflict.expected" "$temp_root/conflict.stderr" ||
	fail 'conflict stderr changed' "$temp_root/conflict.expected" "$temp_root/conflict.stderr"

set +e
(
	cd "$temp_root"
	HOME="$temp_root/home" "$binary" workspace init --name example
) >"$temp_root/name.stdout" 2>"$temp_root/name.stderr"
name_exit=$?
(
	cd "$temp_root"
	HOME="$temp_root/home" "$binary" workspace init unexpected
) >"$temp_root/argument.stdout" 2>"$temp_root/argument.stderr"
argument_exit=$?
set -e
[[ "$name_exit" -eq 1 ]] || fail "--name exit is $name_exit, expected 1" "$temp_root/name.stdout" "$temp_root/name.stderr"
[[ "$argument_exit" -eq 1 ]] || fail "positional argument exit is $argument_exit, expected 1" "$temp_root/argument.stdout" "$temp_root/argument.stderr"
[[ ! -s "$temp_root/name.stdout" && ! -s "$temp_root/argument.stdout" ]] ||
	fail 'invalid workspace input wrote to stdout' "$temp_root/name.stdout" "$temp_root/argument.stdout"
printf 'Error: workspace_init_invalid_arguments: unknown flag: --name\n' >"$temp_root/name.expected"
printf 'Error: workspace_init_invalid_arguments: expected no positional arguments, got 1\n' \
	>"$temp_root/argument.expected"
cmp -s "$temp_root/name.expected" "$temp_root/name.stderr" || fail '--name stderr changed' "$temp_root/name.expected" "$temp_root/name.stderr"
cmp -s "$temp_root/argument.expected" "$temp_root/argument.stderr" ||
	fail 'positional argument stderr changed' "$temp_root/argument.expected" "$temp_root/argument.stderr"

run_project_error() {
	local label=$1
	local expected_class=$2
	local directory=$3
	shift 3
	set +e
	(
		cd "$directory"
		HOME="$temp_root/home" "$binary" "$@"
	) >"$temp_root/$label.stdout" 2>"$temp_root/$label.stderr"
	local exit_status=$?
	set -e
	[[ "$exit_status" -eq 1 ]] || fail "$label exit is $exit_status, expected 1" "$temp_root/$label.stdout" "$temp_root/$label.stderr"
	[[ ! -s "$temp_root/$label.stdout" ]] || fail "$label wrote to stdout" "$temp_root/$label.stdout"
	grep -F "Error: $expected_class:" "$temp_root/$label.stderr" >/dev/null ||
		fail "$label stderr lost class $expected_class" "$temp_root/$label.stderr"
}

single_wrapper="$temp_root/single"
single_repo="$single_wrapper/main"
multi_wrapper="$temp_root/multi"
multi_repo="$multi_wrapper/service"
external_repo="$temp_root/external-api"
non_git="$temp_root/not-git"
mkdir -p "$single_repo" "$multi_repo" "$external_repo" "$non_git"
git -C "$single_repo" init -q
git -C "$multi_repo" init -q
git -C "$external_repo" init -q
printf 'dirty\n' >"$multi_repo/dirty.txt"
single_wrapper_root=$(cd "$single_wrapper" && pwd -P)
single_repo_root=$(cd "$single_repo" && pwd -P)
multi_wrapper_root=$(cd "$multi_wrapper" && pwd -P)
multi_repo_root=$(cd "$multi_repo" && pwd -P)
external_repo_root=$(cd "$external_repo" && pwd -P)
single_common=$(git -C "$single_repo_root" rev-parse --path-format=absolute --git-common-dir)
multi_common=$(git -C "$multi_repo_root" rev-parse --path-format=absolute --git-common-dir)
external_common=$(git -C "$external_repo_root" rev-parse --path-format=absolute --git-common-dir)
find "$single_wrapper" -mindepth 1 -print | LC_ALL=C sort >"$temp_root/single.before"
find "$multi_wrapper" -mindepth 1 -print | LC_ALL=C sort >"$temp_root/multi.before"
find "$external_repo" -mindepth 1 -print | LC_ALL=C sort >"$temp_root/external.before"

(
	cd "$work"
	HOME="$temp_root/home" "$binary" workspace project list
) >"$temp_root/project-empty.stdout" 2>"$temp_root/project-empty.stderr" ||
	fail 'empty project list returned non-zero' "$temp_root/project-empty.stdout" "$temp_root/project-empty.stderr"
printf 'No projects are registered in Ply workspace %s.\n' "$work_root" \
	>"$temp_root/project-empty.expected"
cmp -s "$temp_root/project-empty.expected" "$temp_root/project-empty.stdout" ||
	fail 'empty project list stdout changed' "$temp_root/project-empty.expected" "$temp_root/project-empty.stdout"
[[ ! -s "$temp_root/project-empty.stderr" ]] || fail 'empty project list wrote to stderr' "$temp_root/project-empty.stderr"

(
	cd "$work"
	HOME="$temp_root/home" "$binary" workspace project add single --name Single \
		--wrapper "$single_wrapper_root" --repo "single=$single_repo_root"
) >"$temp_root/project-single.stdout" 2>"$temp_root/project-single.stderr" ||
	fail 'single-repository project add returned non-zero' "$temp_root/project-single.stdout" "$temp_root/project-single.stderr"
printf 'Added project single (Single) to Ply workspace %s.\nWrapper: %s\nRepositories:\n  single:\n    locator: %s\n    git common directory: %s\n' \
	"$work_root" "$single_wrapper_root" "$single_repo_root" "$single_common" \
	>"$temp_root/project-single.expected"
cmp -s "$temp_root/project-single.expected" "$temp_root/project-single.stdout" ||
	fail 'single-repository project add stdout changed' "$temp_root/project-single.expected" "$temp_root/project-single.stdout"
[[ ! -s "$temp_root/project-single.stderr" ]] || fail 'single-repository project add wrote to stderr' "$temp_root/project-single.stderr"
printf 'format_version: 1\nprojects:\n  - id: single\n    name: Single\n    wrapper: %s\n    repo_ids:\n      - single\nrepos:\n  - id: single\n    locator: %s\n    git_common_dir: %s\n' \
	"$single_wrapper_root" "$single_repo_root" "$single_common" >"$temp_root/projects-single.expected"
cmp -s "$temp_root/projects-single.expected" "$work/.ply/projects.yaml" ||
	fail 'single-repository registry bytes changed' "$temp_root/projects-single.expected" "$work/.ply/projects.yaml"
[[ $(file_mode "$work/.ply/projects.yaml") == 644 ]] || fail 'project registry mode is not 0644' "$work/.ply/projects.yaml"
[[ $(file_mode "$work/.ply/projects.lock") == 600 ]] || fail 'project lock mode is not 0600' "$work/.ply/projects.lock"

(
	cd "$work"
	HOME="$temp_root/home" "$binary" workspace project show single
) >"$temp_root/project-show.stdout" 2>"$temp_root/project-show.stderr" ||
	fail 'project show returned non-zero' "$temp_root/project-show.stdout" "$temp_root/project-show.stderr"
printf 'Project single (Single) in Ply workspace %s.\nWrapper: %s\nRepositories:\n  single:\n    locator: %s\n    git common directory: %s\n' \
	"$work_root" "$single_wrapper_root" "$single_repo_root" "$single_common" \
	>"$temp_root/project-show.expected"
cmp -s "$temp_root/project-show.expected" "$temp_root/project-show.stdout" ||
	fail 'project show stdout changed' "$temp_root/project-show.expected" "$temp_root/project-show.stdout"
[[ ! -s "$temp_root/project-show.stderr" ]] || fail 'project show wrote to stderr' "$temp_root/project-show.stderr"

touch -t 200001010000 "$work/.ply/projects.yaml"
project_mtime=$(file_mtime "$work/.ply/projects.yaml")
cp "$work/.ply/projects.yaml" "$temp_root/projects.before-retry"
single_alias="$temp_root/single-alias"
ln -s "$single_repo_root" "$single_alias"
(
	cd "$work"
	HOME="$temp_root/home" "$binary" workspace project add single --name Single \
		--wrapper "$single_wrapper_root" --repo "single=$single_alias"
) >"$temp_root/project-retry.stdout" 2>"$temp_root/project-retry.stderr" ||
	fail 'idempotent project retry returned non-zero' "$temp_root/project-retry.stdout" "$temp_root/project-retry.stderr"
sed '1s/^Added project/Project/; 1s/ to Ply workspace/ is already registered in Ply workspace/' \
	"$temp_root/project-single.expected" >"$temp_root/project-retry.expected"
cmp -s "$temp_root/project-retry.expected" "$temp_root/project-retry.stdout" ||
	fail 'idempotent project retry stdout changed' "$temp_root/project-retry.expected" "$temp_root/project-retry.stdout"
[[ ! -s "$temp_root/project-retry.stderr" ]] || fail 'idempotent project retry wrote to stderr' "$temp_root/project-retry.stderr"
cmp -s "$temp_root/projects.before-retry" "$work/.ply/projects.yaml" ||
	fail 'idempotent project retry changed registry bytes' "$temp_root/projects.before-retry" "$work/.ply/projects.yaml"
[[ $(file_mtime "$work/.ply/projects.yaml") == "$project_mtime" ]] ||
	fail 'idempotent project retry changed registry mtime' "$work/.ply/projects.yaml"

(
	cd "$work"
	HOME="$temp_root/home" "$binary" workspace project add multi --name Multi \
		--wrapper "$multi_wrapper_root" \
		--repo "service=$multi_repo_root" --repo "api=$external_repo_root"
) >"$temp_root/project-multi.stdout" 2>"$temp_root/project-multi.stderr" ||
	fail 'multi-repository project add returned non-zero' "$temp_root/project-multi.stdout" "$temp_root/project-multi.stderr"
[[ ! -s "$temp_root/project-multi.stderr" ]] || fail 'multi-repository project add wrote to stderr' "$temp_root/project-multi.stderr"
grep -F 'Added project multi (Multi)' "$temp_root/project-multi.stdout" >/dev/null ||
	fail 'multi-repository add lost success line' "$temp_root/project-multi.stdout"
grep -F '  api:' "$temp_root/project-multi.stdout" >/dev/null || fail 'multi add lost api member' "$temp_root/project-multi.stdout"
grep -F '  service:' "$temp_root/project-multi.stdout" >/dev/null || fail 'multi add lost service member' "$temp_root/project-multi.stdout"

(
	cd "$work"
	HOME="$temp_root/home" "$binary" workspace project list
) >"$temp_root/project-list.stdout" 2>"$temp_root/project-list.stderr" ||
	fail 'project list returned non-zero' "$temp_root/project-list.stdout" "$temp_root/project-list.stderr"
printf 'Projects in Ply workspace %s:\n  multi: Multi (2 repositories) %s\n  single: Single (1 repository) %s\n' \
	"$work_root" "$multi_wrapper_root" "$single_wrapper_root" >"$temp_root/project-list.expected"
cmp -s "$temp_root/project-list.expected" "$temp_root/project-list.stdout" ||
	fail 'project list stdout changed' "$temp_root/project-list.expected" "$temp_root/project-list.stdout"
[[ ! -s "$temp_root/project-list.stderr" ]] || fail 'project list wrote to stderr' "$temp_root/project-list.stderr"

run_project_error missing-workspace workspace_not_found "$non_git" workspace project list
run_project_error invalid-project workspace_project_invalid_arguments "$work" \
	workspace project add Trip --name Trip --wrapper "$single_wrapper_root" --repo "single=$single_repo_root"
run_project_error empty-repositories workspace_project_invalid_arguments "$work" \
	workspace project add empty --name Empty --wrapper "$single_wrapper_root"
run_project_error non-git workspace_project_repo_invalid "$work" \
	workspace project add nongit --name NonGit --wrapper "$non_git" --repo "nongit=$non_git"
run_project_error duplicate-common workspace_project_duplicate_repo "$work" \
	workspace project add duplicate --name Duplicate --wrapper "$non_git" \
	--repo "one=$single_repo_root" --repo "two=$single_alias"
run_project_error unknown-show workspace_project_not_found "$work" workspace project show unknown
run_project_error changed-members workspace_project_conflict "$work" \
	workspace project add single --name Single --wrapper "$single_wrapper_root" --repo "single=$external_repo_root"
run_project_error duplicate-wrapper workspace_project_conflict "$work" \
	workspace project add other --name Other --wrapper "$single_wrapper_root" --repo "api=$external_repo_root"

find "$single_wrapper" -mindepth 1 -print | LC_ALL=C sort >"$temp_root/single.after"
find "$multi_wrapper" -mindepth 1 -print | LC_ALL=C sort >"$temp_root/multi.after"
find "$external_repo" -mindepth 1 -print | LC_ALL=C sort >"$temp_root/external.after"
cmp -s "$temp_root/single.before" "$temp_root/single.after" || fail 'Ply changed the single wrapper or repository' "$temp_root/single.before" "$temp_root/single.after"
cmp -s "$temp_root/multi.before" "$temp_root/multi.after" || fail 'Ply changed the multi wrapper or repository' "$temp_root/multi.before" "$temp_root/multi.after"
cmp -s "$temp_root/external.before" "$temp_root/external.after" || fail 'Ply changed the external repository' "$temp_root/external.before" "$temp_root/external.after"
[[ $(cd "$work" && find . -mindepth 1 -print | LC_ALL=C sort) == $'./.ply\n./.ply/projects.lock\n./.ply/projects.yaml\n./.ply/workspace.yaml' ]] ||
	fail 'project journey created unexpected workspace files'
[[ ! -e "$temp_root/home/.ply" ]] || fail 'project journey created a global Ply profile' "$temp_root/home/.ply"

printf 'cli surface contract: PASS (legacy behavior, workspace, project, and workflow)\n'


for path in 'task queue' 'task queue list' 'task queue set' 'epic base' 'epic base show' 'epic base update' 'epic base operation' 'epic base operation show' 'task spec' 'task spec record' 'task spec assess' 'task spec select'; do
 read -r -a words <<<"$path"
 label="ws06-${path// /-}"
 run_help_surface "$label" workspace "${words[@]}" --help
 grep -F "ply workspace $path" "$temp_root/$label.stdout" >/dev/null || fail "$path help lost its workspace prefix" "$temp_root/$label.stdout"
 done
for path in 'task queue list' 'task queue set' 'epic base show' 'epic base update' 'epic base operation show'; do
 label="ws06-${path// /-}"
 grep -F -- '--format' "$temp_root/$label.stdout" >/dev/null || fail "$path has no format flag" "$temp_root/$label.stdout"
done
for text in 'workspace task queue set --file /absolute/queue.json' 'workspace epic base update'; do
 grep -F "$text" "$repo_root/README.md" >/dev/null || fail "README queue contract is missing: $text" "$repo_root/README.md"
done
