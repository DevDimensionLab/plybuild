#!/usr/bin/env bash

set -euo pipefail
export LC_ALL=C LANG=C GOENV=off GOWORK=off
unset CDPATH

repo_root=${CLI_COMPAT_REPO:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}
base_ref=${CLI_COMPAT_BASE_REF:-v1.0.1}
baseline=${CLI_COMPAT_BASELINE:-"$repo_root/test/compat/cli-v1.0.1.json"}
allowlist=${CLI_COMPAT_ALLOWLIST:-"$repo_root/test/compat/cli-allowlist.json"}
report_out=${CLI_COMPAT_REPORT_OUT:-"$repo_root/target/compatibility/cli-report.json"}
exporter=${CLI_COMPAT_EXPORTER:-"$repo_root/test/compat/export_cli.go"}
parser="$repo_root/test/compat/check_cli_report.py"
go_bin=${GO:-go}

fail() {
	printf 'cli compatibility: %s\n' "$*" >&2
	exit 1
}

[[ -f "$repo_root/go.mod" && -f "$baseline" && -f "$allowlist" &&
	-f "$exporter" && -f "$parser" ]] || fail 'CLI compatibility inputs are missing'

temp_root=$(mktemp -d "${TMPDIR:-/tmp}/ply-cli-compat.XXXXXX") ||
	fail 'could not create temporary CLI compatibility workspace'
cleanup_temp_root() {
	find "$temp_root" -type d -exec chmod u+rwx {} + 2>/dev/null || true
	rm -rf -- "$temp_root"
}
trap cleanup_temp_root EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM
mkdir -p "$temp_root/base" "$temp_root/gocache"

if [[ -n "${CLI_COMPAT_BASE_GENERATED:-}" || -n "${CLI_COMPAT_CURRENT_GENERATED:-}" ]]; then
	[[ -f "${CLI_COMPAT_BASE_GENERATED:-}" && -f "${CLI_COMPAT_CURRENT_GENERATED:-}" ]] ||
		fail 'both generated CLI test overrides are required'
	cp "$CLI_COMPAT_BASE_GENERATED" "$temp_root/base.json"
	cp "$CLI_COMPAT_CURRENT_GENERATED" "$temp_root/current.json"
else
	if [[ -n "${CLI_COMPAT_BASE_DIR:-}" ]]; then
		cp -R "$CLI_COMPAT_BASE_DIR/." "$temp_root/base/"
	else
		git -C "$repo_root" archive "$base_ref" | tar -xf - -C "$temp_root/base" ||
			fail "could not materialize CLI base $base_ref"
	fi
	mkdir -p "$temp_root/base/.compat-export"
	cp "$exporter" "$temp_root/base/.compat-export/export_cli.go"
	(
		cd "$temp_root/base"
		GOCACHE="$temp_root/gocache" GOPROXY=off GOSUMDB=off \
			"$go_bin" run .compat-export/export_cli.go
	) >"$temp_root/base.json" || fail 'could not export base CLI tree'
	(
		cd "$repo_root"
		GOCACHE="$temp_root/gocache" GOPROXY=off GOSUMDB=off \
			"$go_bin" run "$exporter"
	) >"$temp_root/current.json" || fail 'could not export current CLI tree'
fi

cmp -s "$temp_root/base.json" "$baseline" || fail 'stored v1.0.1 CLI baseline is stale'
python3 "$parser" \
	--base "$baseline" \
	--current "$temp_root/current.json" \
	--allowlist "$allowlist" \
	--output "$report_out" \
	--base-ref "$base_ref"

printf 'cli compatibility: PASS (Cobra tree against %s)\n' "$base_ref"
