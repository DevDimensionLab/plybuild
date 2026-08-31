#!/usr/bin/env bash

set -euo pipefail
export LC_ALL=C LANG=C
unset CDPATH

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)
tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

makefile_under_test=${MAKEFILE_UNDER_TEST:-"$repo_root/Makefile"}
fake_bin="$tmp_dir/bin"
fake_scripts="$tmp_dir/scripts"
calls_file="$tmp_dir/quality-calls"
manual_evidence="$tmp_dir/manual-evidence.json"
gocache="$tmp_dir/gocache"
gomodcache="$tmp_dir/gomodcache"
python_bin=$(realpath "$(command -v python3)")
bash_bin=$(realpath /bin/bash)
mutation_names=(
	cli-context config-cloud maven-sorting template
	file-shell spring http interactive-build
)

mkdir -p "$fake_bin" "$fake_scripts" "$gocache" "$gomodcache"
printf '{}\n' >"$manual_evidence"

fail() {
	printf 'make quality contract: %s\n' "$*" >&2
	exit 1
}

for tool in go goreleaser apidiff golangci-lint docker; do
	cat >"$fake_bin/$tool" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
exit 0
EOF
	chmod +x "$fake_bin/$tool"
done

cat >"$fake_bin/bash" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
{
	printf 'bash'
	printf ' <%s>' "$@"
	printf '\n'
} >>"$QUALITY_TEST_CALLS"
exec /bin/bash "$@"
EOF
chmod +x "$fake_bin/bash"

cat >"$fake_bin/make" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
{
	printf 'make'
	printf ' <%s>' "$@"
	printf '\n'
} >>"$QUALITY_TEST_CALLS"

target=
for argument in "$@"; do
	case "$argument" in
	preflight|acceptance|acceptance-snapshot|acceptance-docker) target=$argument ;;
	esac
done

case "$target" in
preflight)
	[[ -d ${GOLANGCI_LINT_CACHE:-} ]] || exit 94
	;;
acceptance)
	printf '%s\n' \
		'verify-install: PASS' \
		'verify-status: PASS' \
		'verify-upgrade: PASS' \
		'verify-build: PASS'
	;;
acceptance-snapshot)
	if [[ ${FAKE_QUALITY_MODE:-good} != skip-snapshot ]]; then
		mkdir "$PLY_SNAPSHOT_EVIDENCE_ROOT"
		printf 'snapshot-acceptance: PASS\n' >"$PLY_SNAPSHOT_EVIDENCE_ROOT/report.txt"
	fi
	;;
acceptance-docker)
	if [[ ${FAKE_QUALITY_MODE:-good} != skip-docker ]]; then
		mkdir "$PLY_DOCKER_EVIDENCE_ROOT"
		printf 'docker-acceptance: PASS\n' >"$PLY_DOCKER_EVIDENCE_ROOT/report.txt"
	fi
	;;
*)
	printf 'unexpected recursive make argv: %s\n' "$*" >&2
	exit 91
	;;
esac
EOF
chmod +x "$fake_bin/make"

cat >"$fake_bin/quality-audit.sh" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

[[ $# -eq 9 ]] || exit 92
[[ $2 == --baseline && $4 == --manual-evidence && $6 == --out &&
	$8 == --only && $9 == 'Q0.*,Q1.*,Q2.*' ]] || exit 93

case ${FAKE_QUALITY_MODE:-good} in
audit-exit-1) exit 1 ;;
audit-exit-2) exit 2 ;;
esac

mkdir "$7"
"${QUALITY_TEST_PYTHON:?}" - "$7/scorecard.json" "${FAKE_QUALITY_MODE:-good}" <<'PY'
import json
import sys

path, mode = sys.argv[1:]
criteria = [
    {"id": "Q0." + str(value), "verdict": "PASS"}
    for value in range(1, 9)
] + [
    {"id": "Q1." + str(value), "verdict": "PASS"}
    for value in range(1, 10)
] + [
    {"id": "Q2." + str(value), "verdict": "PASS"}
    for value in range(1, 11)
]
if mode == "missing-criterion":
    criteria.pop()
document = {
    "attained_level": "L1" if mode == "not-l2" else "L2",
    "criteria": criteria,
    "denominators": {
        "declared_subjects": 8,
        "mutation_harnesses": 8,
        "declared_features": 4,
        "acceptance_scripts": 4,
        "test_functions": 1,
    },
    "manual_evidence": {"status": "valid", "criterion_receipts": 6},
    "ratchet": {
        "held": 1 if mode == "held-debt" else 0,
        "regressed": 0,
        "current_not_comparable": 0,
    },
    "repository": {"tree": {"measurement_clean": True, "dirty_paths": []}},
}
with open(path, "w", encoding="utf-8") as output:
    json.dump(document, output, sort_keys=True)
PY
EOF
chmod +x "$fake_bin/quality-audit.sh"

for name in "${mutation_names[@]}"; do
	cat >"$fake_scripts/test-mutate-$name" <<EOF
#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' 'test-mutate-$name: PASS (T1-T10, declared=10 killed=10 survived=0 unusable=0)'
EOF
	cat >"$fake_scripts/mutate-$name" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
printf 'declared=10 killed=10 survived=0 unusable=0\n'
EOF
	chmod +x "$fake_scripts/test-mutate-$name" "$fake_scripts/mutate-$name"
done

for name in install status upgrade build; do
	for prefix in verify test-verify; do
		cat >"$fake_scripts/$prefix-$name" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
exit 0
EOF
		chmod +x "$fake_scripts/$prefix-$name"
	done
done
for name in snapshot docker; do
	for prefix in accept test-accept; do
		cat >"$fake_scripts/$prefix-$name" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
exit 0
EOF
		chmod +x "$fake_scripts/$prefix-$name"
	done
done

run_quality() {
	local label=$1
	local mode=${2:-good}
	local makefile=${3:-$makefile_under_test}
	local output="$tmp_dir/output-$label"
	shift 3 || true
	: >"$calls_file"
	FAKE_QUALITY_MODE="$mode" QUALITY_TEST_CALLS="$calls_file" \
	QUALITY_TEST_PYTHON="$python_bin" \
	make --no-print-directory -f "$makefile" REPO_ROOT="$repo_root" \
		MAKE="$fake_bin/make" GO="$fake_bin/go" BASH="$fake_bin/bash" \
		GORELEASER="$fake_bin/goreleaser" APIDIFF="$fake_bin/apidiff" \
		GOLANGCI_LINT="$fake_bin/golangci-lint" DOCKER="$fake_bin/docker" \
		PYTHON="$python_bin" SCRIPTS_DIR="$fake_scripts" \
		QUALITY_AUDIT="$fake_bin/quality-audit.sh" \
		QUALITY_BASELINE="$repo_root/.quality/baseline/scorecard.json" \
		QUALITY_MANUAL_EVIDENCE="$manual_evidence" \
		QUALITY_OUTPUT_ROOT="$output" QUALITY_GOCACHE="$gocache" \
		QUALITY_GOMODCACHE="$gomodcache" "$@" quality
}

expect_failure() {
	local label=$1
	local mode=${2:-good}
	local makefile=${3:-$makefile_under_test}
	shift 3 || true
	if run_quality "$label" "$mode" "$makefile" "$@" \
		>"$tmp_dir/$label.stdout" 2>"$tmp_dir/$label.stderr"; then
		fail "accepted $label"
	fi
}

if ! run_quality good good "$makefile_under_test" \
	>"$tmp_dir/good.stdout" 2>"$tmp_dir/good.stderr"; then
	sed -n '1,240p' "$tmp_dir/good.stderr" >&2
	fail 'complete quality population did not pass'
fi
grep -Fqx 'quality: PASS (Q0-Q2 attained L2)' "$tmp_dir/good.stdout" ||
	fail 'complete population produced no terminal PASS'

expected_calls="$tmp_dir/expected-calls"
{
	printf 'make <--no-print-directory> <-C> <%s> <MAKEOVERRIDES=> <preflight>\n' "$repo_root"
	for name in "${mutation_names[@]}"; do
		printf 'bash <%s/test-mutate-%s>\n' "$fake_scripts" "$name"
	done
	for name in "${mutation_names[@]}"; do
		printf 'bash <%s/mutate-%s>\n' "$fake_scripts" "$name"
	done
	printf 'make <--no-print-directory> <-C> <%s> <acceptance> <GO=%s> <BASH=%s> <SCRIPTS_DIR=%s>\n' \
		"$repo_root" "$fake_bin/go" "$fake_bin/bash" "$fake_scripts"
	printf 'make <--no-print-directory> <-C> <%s> <acceptance-snapshot> <GO=%s> <BASH=%s> <GORELEASER=%s> <SCRIPTS_DIR=%s>\n' \
		"$repo_root" "$fake_bin/go" "$fake_bin/bash" "$fake_bin/goreleaser" "$fake_scripts"
	printf 'make <--no-print-directory> <-C> <%s> <acceptance-docker> <GO=%s> <BASH=%s> <SCRIPTS_DIR=%s>\n' \
		"$repo_root" "$fake_bin/go" "$fake_bin/bash" "$fake_scripts"
	printf 'bash <%s> <%s> <--baseline> <%s> <--manual-evidence> <%s> <--out> <%s> <--only> <Q0.*,Q1.*,Q2.*>\n' \
		"$fake_bin/quality-audit.sh" "$repo_root" "$repo_root/.quality/baseline/scorecard.json" \
		"$manual_evidence" "$tmp_dir/output-good/audit-q0-q2"
} >"$expected_calls"
diff -u "$expected_calls" "$calls_file" || fail 'stage argv or ordering changed'

expect_failure missing-evidence good "$makefile_under_test" QUALITY_MANUAL_EVIDENCE=
expect_failure local-output good "$makefile_under_test" QUALITY_OUTPUT_ROOT="$repo_root/quality-output"
mkdir "$tmp_dir/output-stale-output"
expect_failure stale-output

mv "$fake_scripts/mutate-http" "$tmp_dir/mutate-http"
expect_failure missing-mutation
mv "$tmp_dir/mutate-http" "$fake_scripts/mutate-http"

expect_failure skipped-snapshot skip-snapshot
expect_failure skipped-docker skip-docker
expect_failure audit-exit-1 audit-exit-1
expect_failure audit-exit-2 audit-exit-2
expect_failure not-l2 not-l2
expect_failure missing-criterion missing-criterion
expect_failure held-debt held-debt

wrong_scope="$tmp_dir/Makefile-wrong-scope"
sed "s/--only 'Q0\.\*,Q1\.\*,Q2\.\*'/--only 'Q0.*,Q1.*'/" \
	"$makefile_under_test" >"$wrong_scope"
expect_failure wrong-scope good "$wrong_scope"

printf 'make quality contract: PASS (complete ordered population and 12 fail-closed controls checked)\n'
