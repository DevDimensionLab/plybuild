#!/usr/bin/env bash

set -euo pipefail
export LC_ALL=C LANG=C

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
checker="$repo_root/scripts/check-cli-compat.sh"
temp_root=$(mktemp -d "${TMPDIR:-/tmp}/ply-cli-compat-test.XXXXXX")
trap 'rm -rf "$temp_root"' EXIT
fake_repo="$temp_root/repo"
mkdir -p "$fake_repo/test/compat"
printf 'module example.invalid/cli\n\ngo 1.18\n' >"$fake_repo/go.mod"
cp "$repo_root/test/compat/export_cli.go" \
	"$repo_root/test/compat/check_cli_report.py" "$fake_repo/test/compat/"

fail() {
	printf 'cli compatibility meta-test: %s\n' "$*" >&2
	exit 1
}

write_tree() {
	local path=$1
	local commands=$2
	printf '{"schema_version":1,"root":"ply","commands":%s}\n' "$commands" >"$path"
}

write_allowlist() {
	local deltas=$1
	printf '{"schema_version":1,"base_ref":"v1.0.1","allowed_deltas":%s}\n' \
		"$deltas" >"$fake_repo/test/compat/cli-allowlist.json"
}

run_checker() {
	CLI_COMPAT_REPO="$fake_repo" \
	CLI_COMPAT_BASELINE="$temp_root/stored-base.json" \
	CLI_COMPAT_BASE_GENERATED="$temp_root/generated-base.json" \
	CLI_COMPAT_CURRENT_GENERATED="$temp_root/current.json" \
	CLI_COMPAT_REPORT_OUT="$temp_root/report.json" \
		"$checker"
}

expect_failure() {
	local expected=$1
	shift
	if "$@" >"$temp_root/failure.stdout" 2>"$temp_root/failure.stderr"; then
		fail "expected failure containing: $expected"
	fi
	grep -F "$expected" "$temp_root/failure.stderr" >/dev/null || {
		sed -n '1,80p' "$temp_root/failure.stderr" >&2
		fail "missing failure text: $expected"
	}
}

command_a='{"path":"ply","use":"ply","flags":[{"scope":"persistent","name":"debug","default":"false","usage":"debug"}]}'
command_b='{"path":"ply status","use":"status","flags":[{"scope":"local","name":"show","default":"false","usage":"show"}]}'
write_tree "$temp_root/stored-base.json" "[$command_a,$command_b]"
cp "$temp_root/stored-base.json" "$temp_root/generated-base.json"
write_tree "$temp_root/current.json" "[$command_b,$command_a]"
write_allowlist '[]'
run_checker >"$temp_root/pass.stdout"
grep -F 'cli compatibility: PASS' "$temp_root/pass.stdout" >/dev/null ||
	fail 'reordered equivalent tree did not pass'

write_tree "$temp_root/current.json" "[$command_a]"
expect_failure 'CLI deltas do not match allowlist' run_checker

changed_flag='{"path":"ply","use":"ply","flags":[{"scope":"persistent","name":"debug","default":"true","usage":"changed"}]}'
write_tree "$temp_root/current.json" "[$changed_flag,$command_b]"
expect_failure 'CLI deltas do not match allowlist' run_checker

cp "$temp_root/stored-base.json" "$temp_root/current.json"
printf '\n' >>"$temp_root/generated-base.json"
expect_failure 'stored v1.0.1 CLI baseline is stale' run_checker
cp "$temp_root/stored-base.json" "$temp_root/generated-base.json"

printf 'cli compatibility meta-test: PASS (order, removal, flag, and stale-baseline controls)\n'
