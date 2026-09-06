#!/usr/bin/env bash

set -euo pipefail
export LC_ALL=C LANG=C GOENV=off GOWORK=off
unset CDPATH

repo_root=${API_COMPAT_REPO:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}
base_ref=${API_COMPAT_BASE_REF:-v1.0.1}
allowlist=${API_COMPAT_ALLOWLIST:-"$repo_root/test/compat/api-allowlist.json"}
report_out=${API_COMPAT_REPORT_OUT:-"$repo_root/target/compatibility/api-report.json"}
apidiff_name=${APIDIFF:-apidiff}
expected_version=${APIDIFF_VERSION:-v0.0.0-20260709172345-9ea1abe57597}
go_bin=${GO:-go}
parser="$repo_root/test/compat/check_api_report.py"

fail() {
	printf 'api compatibility: %s\n' "$*" >&2
	exit 1
}

if [[ "$apidiff_name" == */* ]]; then
	apidiff_bin=$apidiff_name
else
	apidiff_bin=$(command -v -- "$apidiff_name" 2>/dev/null) ||
		fail "apidiff $expected_version is required"
fi
[[ -f "$apidiff_bin" && -x "$apidiff_bin" ]] ||
	fail 'configured apidiff is not an executable file'
[[ -f "$repo_root/go.mod" && -f "$allowlist" && -f "$parser" ]] ||
	fail 'repository, allowlist, or report parser is missing'

tool_metadata=$($go_bin version -m "$apidiff_bin" 2>/dev/null) ||
	fail 'could not inspect apidiff build metadata'
tool_path=$(printf '%s\n' "$tool_metadata" | awk '$1 == "path" { print $2 }')
tool_version=$(printf '%s\n' "$tool_metadata" |
	awk '$1 == "mod" && $2 == "golang.org/x/exp" { print $3 }')
[[ "$tool_path" == 'golang.org/x/exp/cmd/apidiff' && "$tool_version" == "$expected_version" ]] ||
	fail "apidiff version mismatch: expected $expected_version"

module=$(awk '$1 == "module" { print $2; exit }' "$repo_root/go.mod")
[[ -n "$module" ]] || fail 'module path is missing from go.mod'

temp_root=$(mktemp -d "${TMPDIR:-/tmp}/ply-api-compat.XXXXXX") ||
	fail 'could not create temporary compatibility workspace'
cleanup_temp_root() {
	find "$temp_root" -type d -exec chmod u+rwx {} + 2>/dev/null || true
	rm -rf -- "$temp_root"
}
trap cleanup_temp_root EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM
mkdir -p "$temp_root/base" "$temp_root/gocache"

if [[ -n "${API_COMPAT_BASE_DIR:-}" ]]; then
	[[ -d "$API_COMPAT_BASE_DIR" && ! -L "$API_COMPAT_BASE_DIR" ]] ||
		fail 'API_COMPAT_BASE_DIR must be a regular directory'
	cp -R "$API_COMPAT_BASE_DIR/." "$temp_root/base/"
else
	git -C "$repo_root" archive "$base_ref" | tar -xf - -C "$temp_root/base" ||
		fail "could not materialize API base $base_ref"
fi

base_module=$(awk '$1 == "module" { print $2; exit }' "$temp_root/base/go.mod")
[[ "$base_module" == "$module" ]] || fail 'base and current module paths differ'

if ! (
	cd "$temp_root/base"
	GOCACHE="$temp_root/gocache" GOPROXY=off GOSUMDB=off \
		"$apidiff_bin" -m -w "$temp_root/base.api" "$module"
) 2>"$temp_root/base.stderr"; then
	sed -n '1,80p' "$temp_root/base.stderr" >&2
	fail 'could not export base API'
fi
if ! (
	cd "$repo_root"
	GOCACHE="$temp_root/gocache" GOPROXY=off GOSUMDB=off \
		"$apidiff_bin" -m -w "$temp_root/current.api" "$module"
) 2>"$temp_root/current.stderr"; then
	sed -n '1,80p' "$temp_root/current.stderr" >&2
	fail 'could not export current API'
fi
if ! "$apidiff_bin" -m "$temp_root/base.api" "$temp_root/current.api" \
	>"$temp_root/apidiff.txt" 2>"$temp_root/report.stderr"; then
	sed -n '1,80p' "$temp_root/report.stderr" >&2
	fail 'apidiff comparison failed'
fi

python3 "$parser" \
	--allowlist "$allowlist" \
	--report "$temp_root/apidiff.txt" \
	--output "$report_out" \
	--base-ref "$base_ref" \
	--module "$module" \
	--tool-version "$expected_version"

printf 'api compatibility: PASS (%s against %s)\n' "$module" "$base_ref"
