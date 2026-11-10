#!/usr/bin/env bash

set -euo pipefail
export LC_ALL=C LANG=C

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
. "$repo_root/scripts/compat_root.sh"
. "$repo_root/scripts/compat_root_test.sh"
compat_run_root 'api compatibility meta-test' ply-api-compat-test
temp_root=$run_root
fake_repo="$temp_root/repo"
fake_base="$temp_root/base"
fake_bin="$temp_root/bin"
checker="$fake_repo/scripts/check-api-compat.sh"
prefix=ply-api-compat
mkdir -p "$fake_repo/test/compat" "$fake_repo/scripts" "$fake_base" "$fake_bin"
cp "$repo_root/scripts/check-api-compat.sh" "$repo_root/scripts/compat_root.sh" "$fake_repo/scripts/"

fail() {
	printf 'api compatibility meta-test: %s\n' "$*" >&2
	exit 1
}

printf 'module example.invalid/compat\n\ngo 1.18\n' >"$fake_repo/go.mod"
cp "$repo_root/test/compat/check_api_report.py" "$fake_repo/test/compat/"
printf 'module example.invalid/compat\n\ngo 1.18\n' >"$fake_base/go.mod"

cat >"$fake_bin/go" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
[[ "${1-}" == version && "${2-}" == -m ]] || exit 90
printf '%s\n' \
	'fake-apidiff: go1.test' \
	$'\tpath\tgolang.org/x/exp/cmd/apidiff' \
	$'\tmod\tgolang.org/x/exp\t'"${FAKE_TOOL_VERSION:-v0.0.0-20260709172345-9ea1abe57597}"$'\th1:test'
EOF

cat >"$fake_bin/apidiff" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
if [[ "${1-}" == -m && "${2-}" == -w ]]; then
	printf '%s %s\n' "$TMPDIR" "$GOTMPDIR" >>"${FAKE_LOG:-/dev/null}"
	if [[ "${FAKE_MODE-}" == hang ]]; then printf 'hang %s\n' "$$" >>"$FAKE_LOG"; exec sleep 60; fi
	printf 'fake export\n' >"$3"
	exit 0
fi
[[ "${1-}" == -m && $# -eq 3 ]] || exit 91
cat "${FAKE_API_REPORT:?}"
EOF
chmod +x "$fake_bin/go" "$fake_bin/apidiff"

write_allowlist() {
	local compatible=$1
	local incompatible=$2
	cat >"$fake_repo/test/compat/api-allowlist.json" <<EOF
{
  "base_ref": "v1.0.1",
  "compatible_changes": $compatible,
  "incompatible_changes": $incompatible,
  "module": "example.invalid/compat",
  "schema_version": 1,
  "tool_module": "golang.org/x/exp",
  "tool_version": "v0.0.0-20260709172345-9ea1abe57597"
}
EOF
}

run_checker() {
	FAKE_API_REPORT="$temp_root/report" \
	API_COMPAT_REPO="$fake_repo" \
	API_COMPAT_BASE_DIR="$fake_base" \
	API_COMPAT_REPORT_OUT="${REPORT_OUT-$temp_root/output.json}" \
	APIDIFF="$fake_bin/apidiff" GO="$fake_bin/go" \
		${checker_exec:-} "$checker"
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

printf 'Compatible changes:\n- package example.invalid/compat/new: added\n' >"$temp_root/report"
write_allowlist '["package example.invalid/compat/new: added"]' '[]'
run_checker >"$temp_root/pass.stdout"
grep -F 'api compatibility: PASS' "$temp_root/pass.stdout" >/dev/null ||
	fail 'valid allowlist did not pass'
grep -F 'example.invalid/compat/new' "$temp_root/output.json" >/dev/null ||
	fail 'machine-readable report omitted the allowed change'

FAKE_TOOL_VERSION=v0.0.0-wrong expect_failure 'apidiff version mismatch' run_checker

printf 'Incompatible changes:\n- Exported: removed\n' >"$temp_root/report"
expect_failure 'incompatible_changes expected [], got' run_checker

printf 'Compatible changes:\n- package example.invalid/compat/other: added\n' >"$temp_root/report"
expect_failure 'compatible_changes expected' run_checker

run_pass() {
	printf 'Compatible changes:\n- package example.invalid/compat/new: added\n' >"$temp_root/report"
	run_checker
}
run_fail() {
	printf 'Incompatible changes:\n- Exported: removed\n' >"$temp_root/report"
	run_checker
}
check_compat_root

printf 'Compatible changes:\n- package example.invalid/compat/new: added\n' >"$temp_root/report"
write_allowlist \
	'["package example.invalid/compat/new: added", "package example.invalid/compat/new: added"]' '[]'
expect_failure 'must be sorted and unique' run_checker

printf 'api compatibility meta-test: PASS (valid, version, incompatible, unexpected, and duplicate controls; owned root, temp, cleanup)\n'
