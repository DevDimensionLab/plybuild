#!/usr/bin/env bash

set -euo pipefail
export LC_ALL=C LANG=C GOENV=off GOWORK=off

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
temp_root=$(mktemp -d "${TMPDIR:-/tmp}/ply-cli-surface.XXXXXX")
cleanup_temp_root() {
	find "$temp_root" -type d -exec chmod u+rwx {} + 2>/dev/null || true
	rm -rf -- "$temp_root"
}
trap cleanup_temp_root EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM
go_bin=${GO:-go}
binary="$temp_root/ply"
mkdir -p "$temp_root/gocache" "$temp_root/home"

fail() {
	printf 'cli surface contract: %s\n' "$*" >&2
	exit 1
}

(
	cd "$repo_root"
	GOCACHE="$temp_root/gocache" "$go_bin" build -o "$binary" ./cmd/ply
) || fail 'could not build fresh CLI artifact'

for surface in root status upgrade build; do
	case $surface in
	root) arguments=(--help); expected='== version: v1.0.1 ==' ;;
	status) arguments=(status --help); expected='Status functionality for a project' ;;
	upgrade) arguments=(upgrade --help); expected='Perform upgrade on existing projects' ;;
	build) arguments=(build --help); expected='Builds a ply project' ;;
	esac
	HOME="$temp_root/home" "$binary" "${arguments[@]}" \
		>"$temp_root/$surface.stdout" 2>"$temp_root/$surface.stderr" ||
		fail "$surface help returned non-zero"
	grep -F "$expected" "$temp_root/$surface.stdout" >/dev/null ||
		fail "$surface help lost its expected text"
	[[ ! -s "$temp_root/$surface.stderr" ]] || fail "$surface help wrote to stderr"
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
	[[ "$exit_status" -eq 0 ]] || fail "$label help exit is $exit_status, expected 0"
	[[ ! -s "$temp_root/$label.stderr" ]] || fail "$label help wrote to stderr"
}

run_help_surface root-help --help
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
run_help_surface workspace-task-create-help workspace task create --help
run_help_surface workspace-task-show-help workspace task show --help
run_help_surface workspace-task-list-help workspace task list --help
run_help_surface workspace-task-worktree-help workspace task worktree --help
run_help_surface workspace-task-worktree-create-help workspace task worktree create --help
run_help_surface workflow-help workflow --help
run_help_surface workflow-handoff-help workflow handoff --help
run_help_surface workflow-handoff-create-help workflow handoff create --help
run_help_surface workflow-handoff-show-help workflow handoff show --help
run_help_surface workflow-handoff-inspect-help workflow handoff inspect --help
run_help_surface workflow-handoff-submit-start-help workflow handoff submit-start --help
run_help_surface workflow-handoff-submit-result-help workflow handoff submit-result --help
run_help_surface workflow-handoff-cancel-help workflow handoff cancel --help
run_help_surface workflow-handoff-supersede-help workflow handoff supersede --help
run_help_surface workflow-handoff-abandon-help workflow handoff abandon --help

grep -F '  workspace   Manage Ply workspaces' "$temp_root/root-help.stdout" >/dev/null ||
	fail 'root help does not expose the workspace parent'
grep -F '  workflow    Manage agent workflows' "$temp_root/root-help.stdout" >/dev/null ||
	fail 'root help does not expose the workflow parent'
grep -F 'Manage explicit local Ply workspaces.' "$temp_root/workspace-help.stdout" >/dev/null ||
	fail 'workspace help lost its long description'
grep -F 'ply workspace [command]' "$temp_root/workspace-help.stdout" >/dev/null ||
	fail 'workspace help usage changed'
grep -F '  init        Initialize a Ply workspace in the current directory' \
	"$temp_root/workspace-help.stdout" >/dev/null || fail 'workspace help does not expose init'
grep -F '  project     Manage projects in a Ply workspace' \
	"$temp_root/workspace-help.stdout" >/dev/null || fail 'workspace help does not expose project'
grep -F '  epic        Manage Epics in a Ply workspace' \
	"$temp_root/workspace-help.stdout" >/dev/null || fail 'workspace help does not expose Epic'
grep -F '  task        Manage Tasks in a Ply workspace' \
	"$temp_root/workspace-help.stdout" >/dev/null || fail 'workspace help does not expose Task'
grep -F '  ply workspace init' "$temp_root/workspace-help.stdout" >/dev/null ||
	fail 'workspace help lost its example'
grep -F '  ply workspace project list' "$temp_root/workspace-help.stdout" >/dev/null ||
	fail 'workspace help lost its project example'
grep -F 'Initialize a Ply workspace in the current directory by creating .ply/workspace.yaml.' \
	"$temp_root/workspace-init-help.stdout" >/dev/null || fail 'workspace init help lost its action'
grep -F 'The current directory does not need to be a Git repository.' \
	"$temp_root/workspace-init-help.stdout" >/dev/null || fail 'workspace init help lost its Git boundary'
grep -F 'ply workspace init [flags]' "$temp_root/workspace-init-help.stdout" >/dev/null ||
	fail 'workspace init help usage changed'
grep -F '  ply workspace init' "$temp_root/workspace-init-help.stdout" >/dev/null ||
	fail 'workspace init help lost its example'
grep -F -- '-h, --help' "$temp_root/workspace-init-help.stdout" >/dev/null ||
	fail 'workspace init help flag is unavailable'
if grep -F -- '--name' "$temp_root/workspace-init-help.stdout" >/dev/null; then
	fail 'workspace init help exposes --name'
fi
grep -F 'Register and inspect projects owned by the containing Ply workspace.' \
	"$temp_root/workspace-project-help.stdout" >/dev/null || fail 'project help lost its long description'
grep -F 'ply workspace project [command]' "$temp_root/workspace-project-help.stdout" >/dev/null ||
	fail 'project help usage changed'
grep -F '  add         Register a project and its explicit repository members' \
	"$temp_root/workspace-project-help.stdout" >/dev/null || fail 'project help does not expose add'
grep -F '  list        List registered projects' "$temp_root/workspace-project-help.stdout" >/dev/null ||
	fail 'project help does not expose list'
grep -F '  show        Show a registered project' "$temp_root/workspace-project-help.stdout" >/dev/null ||
	fail 'project help does not expose show'
grep -F 'ply workspace project add <project-id> [flags]' \
	"$temp_root/workspace-project-add-help.stdout" >/dev/null || fail 'project add usage changed'
grep -F -- '-n, --name string' "$temp_root/workspace-project-add-help.stdout" >/dev/null ||
	fail 'project add name flag is unavailable'
grep -F -- '-w, --wrapper string' "$temp_root/workspace-project-add-help.stdout" >/dev/null ||
	fail 'project add wrapper flag is unavailable'
grep -F -- '-r, --repo stringArray' "$temp_root/workspace-project-add-help.stdout" >/dev/null ||
	fail 'project add repo flag is not a repeatable string array'
grep -F 'ply workspace project show <project-id> [flags]' \
	"$temp_root/workspace-project-show-help.stdout" >/dev/null || fail 'project show usage changed'
grep -F 'ply workspace project list [flags]' \
	"$temp_root/workspace-project-list-help.stdout" >/dev/null || fail 'project list usage changed'
grep -F 'ply workspace epic [command]' "$temp_root/workspace-epic-help.stdout" >/dev/null || fail 'Epic usage changed'
grep -F 'ply workspace epic adopt <epic-id> [flags]' "$temp_root/workspace-epic-adopt-help.stdout" >/dev/null || fail 'Epic adopt usage changed'
for flag in '--title string' '--project string' '--repo string' '--worktree string' '--ref string' '--expected-oid string'; do
	grep -F -- "$flag" "$temp_root/workspace-epic-adopt-help.stdout" >/dev/null || fail "Epic adopt $flag is unavailable"
done
grep -F 'ply workspace epic show <epic-id> [flags]' "$temp_root/workspace-epic-show-help.stdout" >/dev/null || fail 'Epic show usage changed'
grep -F -- '--format string' "$temp_root/workspace-epic-show-help.stdout" >/dev/null || fail 'Epic show format is unavailable'
grep -F 'ply workspace epic list [flags]' "$temp_root/workspace-epic-list-help.stdout" >/dev/null || fail 'Epic list usage changed'
grep -F 'ply workspace task [command]' "$temp_root/workspace-task-help.stdout" >/dev/null || fail 'Task usage changed'
grep -F 'ply workspace task create <task-id> [flags]' "$temp_root/workspace-task-create-help.stdout" >/dev/null || fail 'Task create usage changed'
for flag in '--title string' '--description string' '--epic string' '--project string' '--repo string'; do
	grep -F -- "$flag" "$temp_root/workspace-task-create-help.stdout" >/dev/null || fail "Task create $flag is unavailable"
done
grep -F 'ply workspace task show <task-id> [flags]' "$temp_root/workspace-task-show-help.stdout" >/dev/null || fail 'Task show usage changed'
grep -F 'ply workspace task list [flags]' "$temp_root/workspace-task-list-help.stdout" >/dev/null || fail 'Task list usage changed'
grep -F 'ply workspace task worktree [command]' "$temp_root/workspace-task-worktree-help.stdout" >/dev/null || fail 'Task worktree usage changed'
grep -F 'ply workspace task worktree create <task-id> [flags]' "$temp_root/workspace-task-worktree-create-help.stdout" >/dev/null || fail 'Task worktree create usage changed'
for flag in '--branch string' '--path string' '--expected-parent-oid string'; do
	grep -F -- "$flag" "$temp_root/workspace-task-worktree-create-help.stdout" >/dev/null || fail "Task worktree create $flag is unavailable"
done

grep -F 'Manage explicit local agent workflow transitions.' "$temp_root/workflow-help.stdout" >/dev/null ||
	fail 'workflow help lost its long description'
grep -F 'ply workflow [command]' "$temp_root/workflow-help.stdout" >/dev/null ||
	fail 'workflow help usage changed'
grep -F '  handoff     Manage file-based agent handoffs' "$temp_root/workflow-help.stdout" >/dev/null ||
	fail 'workflow help does not expose handoff'
grep -F '  ply workflow handoff show hnd_0123456789abcdef0123456789abcdef' "$temp_root/workflow-help.stdout" >/dev/null ||
	fail 'workflow help lost its example'
grep -F 'Create, control, receive, and inspect immutable local agent handoffs.' "$temp_root/workflow-handoff-help.stdout" >/dev/null ||
	fail 'handoff help lost its long description'
for leaf in create show inspect submit-start submit-result cancel supersede abandon; do
	grep -F "  $leaf" "$temp_root/workflow-handoff-help.stdout" >/dev/null ||
		fail "handoff help does not expose $leaf"
done
grep -F 'Validate a handoff draft and publish one immutable handoff in the containing Ply workspace.' "$temp_root/workflow-handoff-create-help.stdout" >/dev/null ||
	fail 'create help lost its long description'
grep -F 'ply workflow handoff create [flags]' "$temp_root/workflow-handoff-create-help.stdout" >/dev/null ||
	fail 'create help usage changed'
grep -F -- '--file string' "$temp_root/workflow-handoff-create-help.stdout" >/dev/null || fail 'create file flag is unavailable'
grep -F 'ply workflow handoff show <handoff-id> [flags]' "$temp_root/workflow-handoff-show-help.stdout" >/dev/null || fail 'show usage changed'
grep -F 'ply workflow handoff inspect [flags]' "$temp_root/workflow-handoff-inspect-help.stdout" >/dev/null || fail 'inspect usage changed'
for flag in '--handoff string' '--format string' '--raw string' '--acknowledge-secret-exposure'; do
	grep -F -- "$flag" "$temp_root/workflow-handoff-inspect-help.stdout" >/dev/null || fail "inspect $flag is unavailable"
done
grep -F 'ply workflow handoff submit-start [flags]' "$temp_root/workflow-handoff-submit-start-help.stdout" >/dev/null || fail 'submit-start usage changed'
grep -F 'start receipt draft JSON file' "$temp_root/workflow-handoff-submit-start-help.stdout" >/dev/null || fail 'submit-start file help changed'
grep -F 'ply workflow handoff submit-result [flags]' "$temp_root/workflow-handoff-submit-result-help.stdout" >/dev/null || fail 'submit-result usage changed'
grep -F 'terminal result draft JSON file' "$temp_root/workflow-handoff-submit-result-help.stdout" >/dev/null || fail 'submit-result file help changed'
grep -F 'ply workflow handoff cancel <handoff-id> [flags]' "$temp_root/workflow-handoff-cancel-help.stdout" >/dev/null || fail 'cancel usage changed'
grep -F 'human-readable cancellation reason' "$temp_root/workflow-handoff-cancel-help.stdout" >/dev/null || fail 'cancel reason help changed'
grep -F 'ply workflow handoff supersede <handoff-id> [flags]' "$temp_root/workflow-handoff-supersede-help.stdout" >/dev/null || fail 'supersede usage changed'
grep -F 'replacement handoff draft JSON file' "$temp_root/workflow-handoff-supersede-help.stdout" >/dev/null || fail 'supersede file help changed'
grep -F 'human-readable superseding reason' "$temp_root/workflow-handoff-supersede-help.stdout" >/dev/null || fail 'supersede reason help changed'
grep -F 'ply workflow handoff abandon <handoff-id> [flags]' "$temp_root/workflow-handoff-abandon-help.stdout" >/dev/null || fail 'abandon usage changed'
grep -F -- '--acknowledge-effects-unknown' "$temp_root/workflow-handoff-abandon-help.stdout" >/dev/null || fail 'abandon acknowledgement is unavailable'
[[ ! -e "$temp_root/home/.ply" ]] || fail 'help created a global Ply profile'

grep -F '  workspace   Manage Ply workspaces' "$repo_root/README.md" >/dev/null ||
	fail 'README root command overview does not expose workspace'
grep -F '  workflow    Manage agent workflows' "$repo_root/README.md" >/dev/null ||
	fail 'README root command overview does not expose workflow'
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
	'`ply workspace project list`.' \
	'' \
	'Retrying the same registration is idempotent. Changing membership, relocating repositories, and' \
	'`project init` are not part of this command.' \
	'' >"$temp_root/readme-workspace.expected"
sed -n '/^## Workspace$/,/^## Workflow handoffs$/p' "$repo_root/README.md" | sed '$d' \
	>"$temp_root/readme-workspace.actual"
expected_workspace_bytes=$(wc -c <"$temp_root/readme-workspace.expected" | tr -d ' ')
dd if="$temp_root/readme-workspace.actual" of="$temp_root/readme-workspace.prefix" bs=1 count="$expected_workspace_bytes" 2>/dev/null
cmp -s "$temp_root/readme-workspace.expected" "$temp_root/readme-workspace.prefix" ||
	fail 'README existing workspace documentation changed'
for text in \
	'ply workspace epic adopt ply-agentic-workflow-support \' \
	'It performs no Git change.' \
	'ply workspace task create workspace-work-item-bootstrap \' \
	'binds exactly one registered' \
	'ply workspace task worktree create workspace-work-item-bootstrap \' \
	'--expected-parent-oid 54f3631cbea789f25a4134945c7ca16d343139df' \
	'Ply persists a durable create intent before the additive branch/worktree operation.' \
	'are never reset, removed,' \
	'`ply workspace task show <id>`' \
	'Both show commands accept `--format json`' \
	'it does not start an agent, select a workflow, or grant execution authority.' \
	'Sub-tasks, additional repository anchors on an Epic, editing, rebinding, refreshing, cleanup, and'; do
	grep -F -- "$text" "$temp_root/readme-workspace.actual" >/dev/null || fail "README work-item text is missing: $text"
done
grep -F '## Workflow handoffs' "$repo_root/README.md" >/dev/null || fail 'README workflow section is missing'
grep -F 'The human—not Ply—opens the recipient in that exact working directory' "$repo_root/README.md" >/dev/null ||
	fail 'README does not state the human-start boundary'
grep -F '`submit-start` before any target effect' "$repo_root/README.md" >/dev/null ||
	fail 'README does not state the start-before-effect boundary'
grep -F '`--acknowledge-secret-exposure`' "$repo_root/README.md" >/dev/null ||
	fail 'README does not explain raw secret acknowledgement'
grep -F '`complete` result never authorizes QA' "$repo_root/README.md" >/dev/null ||
	fail 'README does not preserve downstream human gates'
grep -F '"kind":"ply.workflow.handoff-draft"' "$repo_root/README.md" >/dev/null ||
	fail 'README does not include a full draft example'

set +e
HOME="$temp_root/home" "$binary" definitely-not-a-command \
	>"$temp_root/invalid.stdout" 2>"$temp_root/invalid.stderr"
invalid_exit=$?
set -e
[[ "$invalid_exit" -eq 1 ]] || fail "unknown command exit is $invalid_exit, expected 1"
grep -F 'unknown command "definitely-not-a-command" for "ply"' \
	"$temp_root/invalid.stderr" >/dev/null || fail 'unknown-command stderr changed'
grep -F "Run 'ply --help' for usage." "$temp_root/invalid.stderr" >/dev/null ||
	fail 'unknown-command usage hint changed'
[[ ! -s "$temp_root/invalid.stdout" ]] || fail 'unknown command wrote to stdout'

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
[[ ! -e "$work/.git" ]] || fail 'fresh workspace unexpectedly contains .git'
work_root=$(cd "$work" && pwd -P)
(
	cd "$work"
	HOME="$temp_root/home" "$binary" workspace init
) >"$temp_root/work-first.stdout" 2>"$temp_root/work-first.stderr" ||
	fail 'first workspace init returned non-zero'
printf 'Initialized Ply workspace at %s.\n' "$work_root" >"$temp_root/work-first.expected"
cmp -s "$temp_root/work-first.expected" "$temp_root/work-first.stdout" ||
	fail 'first workspace init stdout changed'
[[ ! -s "$temp_root/work-first.stderr" ]] || fail 'first workspace init wrote to stderr'
printf 'format_version: 1\nroot: %s\n' "$work_root" >"$temp_root/marker.expected"
cmp -s "$temp_root/marker.expected" "$work/.ply/workspace.yaml" ||
	fail 'workspace marker bytes changed'
[[ $(file_mode "$work/.ply") == 755 ]] || fail 'workspace directory mode is not 0755'
[[ $(file_mode "$work/.ply/workspace.yaml") == 644 ]] || fail 'workspace marker mode is not 0644'
touch -t 200001010000 "$work/.ply/workspace.yaml"
first_mtime=$(file_mtime "$work/.ply/workspace.yaml")
cp "$work/.ply/workspace.yaml" "$temp_root/marker.before"
(
	cd "$work"
	HOME="$temp_root/home" "$binary" workspace init
) >"$temp_root/work-second.stdout" 2>"$temp_root/work-second.stderr" ||
	fail 'second workspace init returned non-zero'
printf 'Ply workspace already initialized at %s.\n' "$work_root" >"$temp_root/work-second.expected"
cmp -s "$temp_root/work-second.expected" "$temp_root/work-second.stdout" ||
	fail 'second workspace init stdout changed'
[[ ! -s "$temp_root/work-second.stderr" ]] || fail 'second workspace init wrote to stderr'
cmp -s "$temp_root/marker.before" "$work/.ply/workspace.yaml" ||
	fail 'idempotent init changed marker bytes'
[[ $(file_mtime "$work/.ply/workspace.yaml") == "$first_mtime" ]] ||
	fail 'idempotent init changed marker mtime'
[[ ! -e "$temp_root/home/.ply" ]] || fail 'workspace init created a global Ply profile'
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
	fail 'workspace init through symlink returned non-zero'
printf 'Initialized Ply workspace at %s.\n' "$physical_root" >"$temp_root/symlink.expected"
cmp -s "$temp_root/symlink.expected" "$temp_root/symlink.stdout" ||
	fail 'symlink init did not report the physical root'
printf 'format_version: 1\nroot: %s\n' "$physical_root" >"$temp_root/symlink-marker.expected"
cmp -s "$temp_root/symlink-marker.expected" "$physical/.ply/workspace.yaml" ||
	fail 'symlink init marker did not use the physical root'
[[ ! -s "$temp_root/symlink.stderr" ]] || fail 'symlink init wrote to stderr'

nested_parent="$temp_root/nested-parent"
nested_child="$nested_parent/child"
mkdir -p "$nested_child"
nested_parent_root=$(cd "$nested_parent" && pwd -P)
nested_child_root=$(cd "$nested_child" && pwd -P)
(
	cd "$nested_parent"
	HOME="$temp_root/home" "$binary" workspace init
) >"$temp_root/nested-parent.stdout" 2>"$temp_root/nested-parent.stderr" ||
	fail 'nested parent init returned non-zero'
set +e
(
	cd "$nested_child"
	HOME="$temp_root/home" "$binary" workspace init
) >"$temp_root/nested.stdout" 2>"$temp_root/nested.stderr"
nested_exit=$?
set -e
[[ "$nested_exit" -eq 1 ]] || fail "nested workspace exit is $nested_exit, expected 1"
[[ ! -s "$temp_root/nested.stdout" ]] || fail 'nested workspace wrote to stdout'
printf 'Error: workspace_init_nested: %s is inside Ply workspace %s\n' \
	"$nested_child_root" "$nested_parent_root" >"$temp_root/nested.expected"
cmp -s "$temp_root/nested.expected" "$temp_root/nested.stderr" ||
	fail 'nested workspace stderr changed'
[[ ! -e "$nested_child/.ply" ]] || fail 'nested workspace created child .ply'

compatible="$temp_root/compatible"
mkdir -p "$compatible/.ply"
compatible_root=$(cd "$compatible" && pwd -P)
printf '# retained\nroot: "%s"\nformat_version: 1\n' "$compatible_root" \
	>"$compatible/.ply/workspace.yaml"
(
	cd "$compatible"
	HOME="$temp_root/home" "$binary" workspace init
) >"$temp_root/compatible.stdout" 2>"$temp_root/compatible.stderr" ||
	fail 'compatible workspace init returned non-zero'
printf 'Ply workspace already initialized at %s.\n' "$compatible_root" \
	>"$temp_root/compatible.expected"
cmp -s "$temp_root/compatible.expected" "$temp_root/compatible.stdout" ||
	fail 'compatible workspace stdout changed'
[[ ! -s "$temp_root/compatible.stderr" ]] || fail 'compatible workspace wrote to stderr'

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
[[ "$conflict_exit" -eq 1 ]] || fail "conflict exit is $conflict_exit, expected 1"
[[ ! -s "$temp_root/conflict.stdout" ]] || fail 'conflict wrote to stdout'
printf 'Error: workspace_init_conflict: %s/.ply/workspace.yaml: workspace marker is missing\n' \
	"$conflict_root" >"$temp_root/conflict.expected"
cmp -s "$temp_root/conflict.expected" "$temp_root/conflict.stderr" ||
	fail 'conflict stderr changed'

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
[[ "$name_exit" -eq 1 ]] || fail "--name exit is $name_exit, expected 1"
[[ "$argument_exit" -eq 1 ]] || fail "positional argument exit is $argument_exit, expected 1"
[[ ! -s "$temp_root/name.stdout" && ! -s "$temp_root/argument.stdout" ]] ||
	fail 'invalid workspace input wrote to stdout'
printf 'Error: workspace_init_invalid_arguments: unknown flag: --name\n' >"$temp_root/name.expected"
printf 'Error: workspace_init_invalid_arguments: expected no positional arguments, got 1\n' \
	>"$temp_root/argument.expected"
cmp -s "$temp_root/name.expected" "$temp_root/name.stderr" || fail '--name stderr changed'
cmp -s "$temp_root/argument.expected" "$temp_root/argument.stderr" ||
	fail 'positional argument stderr changed'

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
	[[ "$exit_status" -eq 1 ]] || fail "$label exit is $exit_status, expected 1"
	[[ ! -s "$temp_root/$label.stdout" ]] || fail "$label wrote to stdout"
	grep -F "Error: $expected_class:" "$temp_root/$label.stderr" >/dev/null ||
		fail "$label stderr lost class $expected_class"
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
	fail 'empty project list returned non-zero'
printf 'No projects are registered in Ply workspace %s.\n' "$work_root" \
	>"$temp_root/project-empty.expected"
cmp -s "$temp_root/project-empty.expected" "$temp_root/project-empty.stdout" ||
	fail 'empty project list stdout changed'
[[ ! -s "$temp_root/project-empty.stderr" ]] || fail 'empty project list wrote to stderr'

(
	cd "$work"
	HOME="$temp_root/home" "$binary" workspace project add single --name Single \
		--wrapper "$single_wrapper_root" --repo "single=$single_repo_root"
) >"$temp_root/project-single.stdout" 2>"$temp_root/project-single.stderr" ||
	fail 'single-repository project add returned non-zero'
printf 'Added project single (Single) to Ply workspace %s.\nWrapper: %s\nRepositories:\n  single:\n    locator: %s\n    git common directory: %s\n' \
	"$work_root" "$single_wrapper_root" "$single_repo_root" "$single_common" \
	>"$temp_root/project-single.expected"
cmp -s "$temp_root/project-single.expected" "$temp_root/project-single.stdout" ||
	fail 'single-repository project add stdout changed'
[[ ! -s "$temp_root/project-single.stderr" ]] || fail 'single-repository project add wrote to stderr'
printf 'format_version: 1\nprojects:\n  - id: single\n    name: Single\n    wrapper: %s\n    repo_ids:\n      - single\nrepos:\n  - id: single\n    locator: %s\n    git_common_dir: %s\n' \
	"$single_wrapper_root" "$single_repo_root" "$single_common" >"$temp_root/projects-single.expected"
cmp -s "$temp_root/projects-single.expected" "$work/.ply/projects.yaml" ||
	fail 'single-repository registry bytes changed'
[[ $(file_mode "$work/.ply/projects.yaml") == 644 ]] || fail 'project registry mode is not 0644'
[[ $(file_mode "$work/.ply/projects.lock") == 600 ]] || fail 'project lock mode is not 0600'

(
	cd "$work"
	HOME="$temp_root/home" "$binary" workspace project show single
) >"$temp_root/project-show.stdout" 2>"$temp_root/project-show.stderr" ||
	fail 'project show returned non-zero'
printf 'Project single (Single) in Ply workspace %s.\nWrapper: %s\nRepositories:\n  single:\n    locator: %s\n    git common directory: %s\n' \
	"$work_root" "$single_wrapper_root" "$single_repo_root" "$single_common" \
	>"$temp_root/project-show.expected"
cmp -s "$temp_root/project-show.expected" "$temp_root/project-show.stdout" ||
	fail 'project show stdout changed'
[[ ! -s "$temp_root/project-show.stderr" ]] || fail 'project show wrote to stderr'

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
	fail 'idempotent project retry returned non-zero'
sed '1s/^Added project/Project/; 1s/ to Ply workspace/ is already registered in Ply workspace/' \
	"$temp_root/project-single.expected" >"$temp_root/project-retry.expected"
cmp -s "$temp_root/project-retry.expected" "$temp_root/project-retry.stdout" ||
	fail 'idempotent project retry stdout changed'
[[ ! -s "$temp_root/project-retry.stderr" ]] || fail 'idempotent project retry wrote to stderr'
cmp -s "$temp_root/projects.before-retry" "$work/.ply/projects.yaml" ||
	fail 'idempotent project retry changed registry bytes'
[[ $(file_mtime "$work/.ply/projects.yaml") == "$project_mtime" ]] ||
	fail 'idempotent project retry changed registry mtime'

(
	cd "$work"
	HOME="$temp_root/home" "$binary" workspace project add multi --name Multi \
		--wrapper "$multi_wrapper_root" \
		--repo "service=$multi_repo_root" --repo "api=$external_repo_root"
) >"$temp_root/project-multi.stdout" 2>"$temp_root/project-multi.stderr" ||
	fail 'multi-repository project add returned non-zero'
[[ ! -s "$temp_root/project-multi.stderr" ]] || fail 'multi-repository project add wrote to stderr'
grep -F 'Added project multi (Multi)' "$temp_root/project-multi.stdout" >/dev/null ||
	fail 'multi-repository add lost success line'
grep -F '  api:' "$temp_root/project-multi.stdout" >/dev/null || fail 'multi add lost api member'
grep -F '  service:' "$temp_root/project-multi.stdout" >/dev/null || fail 'multi add lost service member'

(
	cd "$work"
	HOME="$temp_root/home" "$binary" workspace project list
) >"$temp_root/project-list.stdout" 2>"$temp_root/project-list.stderr" ||
	fail 'project list returned non-zero'
printf 'Projects in Ply workspace %s:\n  multi: Multi (2 repositories) %s\n  single: Single (1 repository) %s\n' \
	"$work_root" "$multi_wrapper_root" "$single_wrapper_root" >"$temp_root/project-list.expected"
cmp -s "$temp_root/project-list.expected" "$temp_root/project-list.stdout" ||
	fail 'project list stdout changed'
[[ ! -s "$temp_root/project-list.stderr" ]] || fail 'project list wrote to stderr'

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
cmp -s "$temp_root/single.before" "$temp_root/single.after" || fail 'Ply changed the single wrapper or repository'
cmp -s "$temp_root/multi.before" "$temp_root/multi.after" || fail 'Ply changed the multi wrapper or repository'
cmp -s "$temp_root/external.before" "$temp_root/external.after" || fail 'Ply changed the external repository'
[[ $(cd "$work" && find . -mindepth 1 -print | LC_ALL=C sort) == $'./.ply\n./.ply/projects.lock\n./.ply/projects.yaml\n./.ply/workspace.yaml' ]] ||
	fail 'project journey created unexpected workspace files'
[[ ! -e "$temp_root/home/.ply" ]] || fail 'project journey created a global Ply profile'

printf 'cli surface contract: PASS (legacy behavior, workspace, project, and workflow handoff)\n'
