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

grep -F '  workspace   Manage Ply workspaces' "$temp_root/root-help.stdout" >/dev/null ||
	fail 'root help does not expose the workspace parent'
grep -F 'Manage explicit local Ply workspaces.' "$temp_root/workspace-help.stdout" >/dev/null ||
	fail 'workspace help lost its long description'
grep -F 'ply workspace [command]' "$temp_root/workspace-help.stdout" >/dev/null ||
	fail 'workspace help usage changed'
grep -F '  init        Initialize a Ply workspace in the current directory' \
	"$temp_root/workspace-help.stdout" >/dev/null || fail 'workspace help does not expose init'
grep -F '  ply workspace init' "$temp_root/workspace-help.stdout" >/dev/null ||
	fail 'workspace help lost its example'
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
[[ ! -e "$temp_root/home/.ply" ]] || fail 'help created a global Ply profile'

grep -F '  workspace   Manage Ply workspaces' "$repo_root/README.md" >/dev/null ||
	fail 'README root command overview does not expose workspace'
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
	'safe and leaves an existing compatible marker unchanged. It does not create a Git repository,' \
	'register repositories, or create workflows.' \
	'' >"$temp_root/readme-workspace.expected"
sed -n '/^## Workspace$/,/^## Install$/p' "$repo_root/README.md" | sed '$d' \
	>"$temp_root/readme-workspace.actual"
cmp -s "$temp_root/readme-workspace.expected" "$temp_root/readme-workspace.actual" ||
	fail 'README workspace documentation changed'

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

printf 'cli surface contract: PASS (legacy help/error and workspace init behavior)\n'
