#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

makefile_under_test=${MAKEFILE_UNDER_TEST:-"$repo_root/Makefile"}
config_under_test=${GOLANGCI_CONFIG_UNDER_TEST:-"$repo_root/.golangci.yml"}
calls_file="$tmp_dir/lint-calls"
format_calls_file="$tmp_dir/format-calls"
explicit_bin_dir="$tmp_dir/explicit"
path_bin_dir="$tmp_dir/path"
make_workspace="$tmp_dir/make workspace"

mkdir -p "$explicit_bin_dir" "$path_bin_dir" "$make_workspace"
touch "$make_workspace/lint" "$make_workspace/format"

write_fake_linter() {
	local path=$1
	cat >"$path" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

{
	printf 'CALL\n'
	printf 'bin=%s\n' "$0"
	printf 'cwd=%s\n' "$PWD"
	printf 'GOPROXY=%s\n' "${GOPROXY-}"
	printf 'GOSUMDB=%s\n' "${GOSUMDB-}"
	printf 'GOTOOLCHAIN=%s\n' "${GOTOOLCHAIN-}"
	printf 'argc=%s\n' "$#"
	printf 'arg=%s\n' "$@"
} >>"$FAKE_LINT_CALLS"

case "$*" in
	'version --short')
		printf '%s\n' "${FAKE_LINT_VERSION:-2.12.2}"
		exit "${FAKE_LINT_VERSION_RC:-0}"
		;;
	'config verify --config .golangci.yml')
		exit 0
		;;
	'run --config .golangci.yml --modules-download-mode=readonly ./...')
		exit 0
		;;
	*)
		printf 'unexpected fake golangci-lint invocation: %s\n' "$*" >&2
		exit 97
		;;
esac
EOF
	chmod +x "$path"
}

write_fake_linter "$explicit_bin_dir/golangci-lint"
write_fake_linter "$path_bin_dir/golangci-lint"

cat >"$tmp_dir/recording-gofmt" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
{
	printf 'CALL\n'
	printf 'arg=%s\n' "$@"
} >>"$FAKE_FORMAT_CALLS"
EOF
chmod +x "$tmp_dir/recording-gofmt"

assert_config() {
	cat >"$tmp_dir/expected-golangci.yml" <<'EOF'
version: "2"

linters:
  default: none
  enable:
    - errcheck
    - govet
    - ineffassign
    - staticcheck
    - unused
  exclusions:
    rules:
      # Preserve the legacy exported Kibana result/error ordering.
      - path: pkg/kibana/kibana.go
        linters:
          - staticcheck
        text: "ST1008:"
      # These legacy messages are user-visible; preserve their exact text.
      - path: pkg/file/file.go
        linters:
          - staticcheck
        text: "ST1005:"
      # Keep the public Upgrade constant's declaration and inferred type stable.
      - path: pkg/webservice/init.go
        linters:
          - staticcheck
        text: "SA9004:"

formatters:
  enable:
    - gofmt
  settings:
    gofmt:
      simplify: false

issues:
  max-issues-per-linter: 0
  max-same-issues: 0
EOF
	diff -u "$tmp_dir/expected-golangci.yml" "$config_under_test"
}

run_lint() {
	local configured_bin=${1-}
	: >"$calls_file"
	(
		cd "$make_workspace"
		if [[ -n "$configured_bin" ]]; then
			PATH="$path_bin_dir:$PATH" FAKE_LINT_CALLS="$calls_file" \
				GOLANGCI_LINT="$configured_bin" \
				make --no-print-directory -f "$makefile_under_test" lint
		else
			PATH="$path_bin_dir:$PATH" FAKE_LINT_CALLS="$calls_file" \
				env -u GOLANGCI_LINT make --no-print-directory -f "$makefile_under_test" lint
		fi
	)
}

assert_three_read_only_calls() {
	local expected_bin=$1
	[[ $(grep -c '^CALL$' "$calls_file") -eq 3 ]]
	[[ $(grep -c -F "bin=$expected_bin" "$calls_file") -eq 3 ]]
	[[ $(grep -c -F "cwd=$repo_root" "$calls_file") -eq 3 ]]
	[[ $(grep -c '^GOPROXY=off$' "$calls_file") -eq 3 ]]
	[[ $(grep -c '^GOSUMDB=off$' "$calls_file") -eq 3 ]]
	[[ $(grep -c '^GOTOOLCHAIN=local$' "$calls_file") -eq 3 ]]
	[[ $(grep -c '^arg=version$' "$calls_file") -eq 1 ]]
	[[ $(grep -c '^arg=verify$' "$calls_file") -eq 1 ]]
	[[ $(grep -c '^arg=run$' "$calls_file") -eq 1 ]]
	if grep -E '^arg=(--fix|fmt)$|^arg=-w$' "$calls_file" >/dev/null; then
		printf 'lint target requested a source rewrite\n' >&2
		exit 1
	fi
}

assert_config

run_lint "$explicit_bin_dir/golangci-lint"
assert_three_read_only_calls "$explicit_bin_dir/golangci-lint"

run_lint
assert_three_read_only_calls "$path_bin_dir/golangci-lint"

: >"$calls_file"
if PATH="$path_bin_dir:$PATH" FAKE_LINT_CALLS="$calls_file" \
	GOLANGCI_LINT="$tmp_dir/missing-golangci-lint" \
	make --no-print-directory -f "$makefile_under_test" lint \
	>"$tmp_dir/missing-out" 2>"$tmp_dir/missing-error"; then
	printf 'lint accepted a missing explicitly configured binary\n' >&2
	exit 1
fi
[[ ! -s "$calls_file" ]]
grep -F 'golangci-lint 2.12.2 is required' "$tmp_dir/missing-error" >/dev/null
grep -F 'curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v2.12.2' \
	"$tmp_dir/missing-error" >/dev/null

: >"$calls_file"
if PATH="$path_bin_dir:$PATH" FAKE_LINT_CALLS="$calls_file" \
	FAKE_LINT_VERSION=2.12.1 GOLANGCI_LINT="$explicit_bin_dir/golangci-lint" \
	make --no-print-directory -f "$makefile_under_test" lint \
	>"$tmp_dir/version-out" 2>"$tmp_dir/version-error"; then
	printf 'lint accepted the wrong golangci-lint version\n' >&2
	exit 1
fi
[[ $(grep -c '^CALL$' "$calls_file") -eq 1 ]]
grep -F 'expected 2.12.2, got 2.12.1' "$tmp_dir/version-error" >/dev/null

: >"$format_calls_file"
(
	cd "$make_workspace"
	FAKE_FORMAT_CALLS="$format_calls_file" GOFMT="$tmp_dir/recording-gofmt" \
		make --no-print-directory -f "$makefile_under_test" format
)
[[ $(grep -c '^CALL$' "$format_calls_file") -gt 0 ]]
awk '
	$0 == "CALL" { expect_write = 1; next }
	expect_write { if ($0 != "arg=-w") exit 1; expect_write = 0 }
	END { if (expect_write) exit 1 }
' "$format_calls_file"
[[ $(grep -c '^arg=.*\.go$' "$format_calls_file") -gt 0 ]]
if grep '^arg=' "$format_calls_file" | grep -v '^arg=-w$' | grep -v '\.go$' >/dev/null; then
	printf 'format target passed a non-Go source to gofmt\n' >&2
	exit 1
fi

printf 'make lint contract: PASS (pin, precedence, offline read-only run, config, and format split checked)\n'
