#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

makefile_under_test=${MAKEFILE_UNDER_TEST:-"$repo_root/Makefile"}
calls_file="$tmp_dir/preflight-calls"
fake_bin="$tmp_dir/bin"
complete_scripts="$tmp_dir/complete/scripts"
empty_scripts="$tmp_dir/empty/scripts"
incomplete_scripts="$tmp_dir/incomplete/scripts"
orphan_scripts="$tmp_dir/orphan/scripts"
omitted_scripts="$tmp_dir/omitted/scripts"
make_workspace="$tmp_dir/make workspace"

mkdir -p "$fake_bin" "$complete_scripts" "$empty_scripts" \
	"$incomplete_scripts" "$orphan_scripts" "$make_workspace"
touch "$make_workspace/preflight"

fail() {
	printf 'make preflight contract: %s\n' "$*" >&2
	exit 1
}

cat >"$fake_bin/go" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
{
	printf 'go'
	printf ' <%s>' "$@"
	printf '\n'
} >>"$FAKE_PREFLIGHT_CALLS"
case ${1-} in
	build|test|vet) exit 0 ;;
	*) exit 91 ;;
esac
EOF

cat >"$fake_bin/golangci-lint" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
{
	printf 'lint'
	printf ' <%s>' "$@"
	printf '\n'
} >>"$FAKE_PREFLIGHT_CALLS"
case "$*" in
	'version --short') printf '2.12.2\n' ;;
	'config verify --config .golangci.yml') ;;
	'run --config .golangci.yml --modules-download-mode=readonly ./...') ;;
	*) exit 92 ;;
esac
EOF

cat >"$fake_bin/bash" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
[[ $# -eq 1 ]] || exit 93
printf 'bash <%s>\n' "$1" >>"$FAKE_PREFLIGHT_CALLS"
EOF

chmod +x "$fake_bin/go" "$fake_bin/golangci-lint" "$fake_bin/bash"

for script in alpha beta test-alpha test-beta; do
	touch "$complete_scripts/$script"
done
for script in alpha beta test-alpha; do
	touch "$incomplete_scripts/$script"
done
for script in alpha test-alpha test-orphan; do
	touch "$orphan_scripts/$script"
done

run_preflight() {
	local scripts_dir=$1
	local output=$2
	: >"$calls_file"
	(
		cd "$make_workspace"
		FAKE_PREFLIGHT_CALLS="$calls_file" make --no-print-directory \
			-f "$makefile_under_test" GO="$fake_bin/go" BASH="$fake_bin/bash" \
			GOLANGCI_LINT="$fake_bin/golangci-lint" SCRIPTS_DIR="$scripts_dir" preflight
	) >"$output" 2>&1
}

assert_once() {
	local expected=$1
	[[ $(grep -Fxc "$expected" "$calls_file") -eq 1 ]] ||
		fail "expected exactly one call: $expected"
}

if ! run_preflight "$complete_scripts" "$tmp_dir/complete-output"; then
	sed -n '1,200p' "$tmp_dir/complete-output" >&2
	fail 'complete population did not pass'
fi
if [[ $(wc -l <"$calls_file" | tr -d '[:space:]') -ne 17 ]]; then
	sed -n '1,200p' "$calls_file" >&2
	fail 'complete preflight did not execute exactly 17 required calls'
fi
assert_once 'go <build> <./...>'
assert_once 'go <test> <./...> <-count=1>'
assert_once 'go <vet> <./...>'
assert_once 'lint <version> <--short>'
assert_once 'lint <config> <verify> <--config> <.golangci.yml>'
assert_once 'lint <run> <--config> <.golangci.yml> <--modules-download-mode=readonly> <./...>'
assert_once "bash <$repo_root/test/codex_dev_start_test.sh>"
assert_once "bash <$repo_root/test/makefile_distribution_test.sh>"
assert_once "bash <$repo_root/test/makefile_lint_test.sh>"
assert_once "bash <$repo_root/test/makefile_install_test.sh>"
assert_once "bash <$repo_root/test/toolchain_declarations_test.sh>"
assert_once "bash <$repo_root/scripts/check-api-compat.sh>"
assert_once "bash <$repo_root/scripts/check-cli-compat.sh>"
assert_once "bash <$repo_root/test/cli_surface_contract_test.sh>"
assert_once "bash <$complete_scripts/test-alpha>"
assert_once "bash <$complete_scripts/test-beta>"
assert_once "bash <$repo_root/.quality/tools/test-quality-audit.sh>"
if grep -Ei 'release|publish|goreleaser|docker' "$calls_file" >/dev/null; then
	fail 'preflight invoked a publishing or packaging path'
fi

if run_preflight "$omitted_scripts" "$tmp_dir/omitted-output"; then
	fail 'preflight accepted an omitted scripts directory'
fi
grep -F 'preflight: script directory does not exist' "$tmp_dir/omitted-output" >/dev/null ||
	fail 'omitted population did not fail for the script directory'

if run_preflight "$empty_scripts" "$tmp_dir/empty-output"; then
	fail 'preflight accepted an empty script population'
fi
grep -F 'preflight: production script population is empty' "$tmp_dir/empty-output" >/dev/null ||
	fail 'empty population did not fail closed'

if run_preflight "$incomplete_scripts" "$tmp_dir/incomplete-output"; then
	fail 'preflight accepted an incomplete script population'
fi
grep -F 'preflight: missing meta-test for beta' "$tmp_dir/incomplete-output" >/dev/null ||
	fail 'incomplete population did not identify the missing meta-test'

if run_preflight "$orphan_scripts" "$tmp_dir/orphan-output"; then
	fail 'preflight accepted an orphan script meta-test'
fi
grep -F 'preflight: meta-test has no production script: test-orphan' "$tmp_dir/orphan-output" >/dev/null ||
	fail 'orphan population did not identify the extra meta-test'

printf 'make preflight contract: PASS (complete, empty, omitted, incomplete, orphan, and non-publishing populations checked)\n'
