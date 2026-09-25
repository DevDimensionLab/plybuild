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

set +e
HOME="$temp_root/home" "$binary" definitely-not-a-command \
	>"$temp_root/invalid.stdout" 2>"$temp_root/invalid.stderr"
invalid_exit=$?
set -e
[[ "$invalid_exit" -eq 1 ]] || fail "unknown command exit is $invalid_exit, expected 1"
grep -F 'unknown command "definitely-not-a-command" for "ply"' \
	"$temp_root/invalid.stderr" >/dev/null || fail 'unknown-command stderr changed'
grep -F 'unknown command "definitely-not-a-command" for "ply"' \
	"$temp_root/invalid.stdout" >/dev/null || fail 'Execute stdout error changed'

printf 'cli surface contract: PASS (root, status, upgrade, build, and exit behavior)\n'
