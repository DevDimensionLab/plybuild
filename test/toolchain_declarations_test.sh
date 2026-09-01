#!/usr/bin/env bash

set -euo pipefail
export LC_ALL=C LANG=C

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)
go_mod=${GO_MOD_UNDER_TEST:-"$repo_root/go.mod"}
dockerfile=${DOCKERFILE_UNDER_TEST:-"$repo_root/Dockerfile"}
release_workflow=${RELEASE_WORKFLOW_UNDER_TEST:-"$repo_root/.github/workflows/release.yaml.deactivated"}
lint_workflow=${LINT_WORKFLOW_UNDER_TEST:-"$repo_root/.github/workflows/lint.yaml.deactivated"}
readme=${README_UNDER_TEST:-"$repo_root/README.md"}
baseline=${QUALITY_BASELINE_UNDER_TEST:-"$repo_root/.quality/baseline/scorecard.json"}
migration=${TOOLCHAIN_MIGRATION_UNDER_TEST:-"$repo_root/.quality/baseline/toolchain-migration.json"}

fail() {
	printf 'toolchain declaration contract: %s\n' "$*" >&2
	exit 1
}

for declaration in "$go_mod" "$dockerfile" "$release_workflow" "$lint_workflow" \
	"$readme" "$baseline" "$migration"; do
	[[ -f "$declaration" && ! -L "$declaration" ]] ||
		fail "declaration is missing, non-regular, or symlinked: $declaration"
done

[[ $(awk '$1 == "go" { print $2 }' "$go_mod") == 1.18 ]] ||
	fail 'go.mod does not preserve the Go 1.18 language-compatibility floor'
[[ $(awk '$1 == "go" { count++ } END { print count + 0 }' "$go_mod") -eq 1 ]] ||
	fail 'go.mod must contain exactly one go directive'
[[ $(awk '$1 == "toolchain" { print $2 }' "$go_mod") == go1.26.7 ]] ||
	fail 'go.mod does not declare exact preferred toolchain go1.26.7'
[[ $(awk '$1 == "toolchain" { count++ } END { print count + 0 }' "$go_mod") -eq 1 ]] ||
	fail 'go.mod must contain exactly one toolchain directive'

expected_builder='FROM golang:1.26.7-alpine3.24@sha256:28d89ee9cc0ff9fec75c82ca201e6bf7fdf9a679d4b7b24dfa04f2bb766bb468 AS build'
[[ $(awk '/^FROM / { print; exit }' "$dockerfile") == "$expected_builder" ]] ||
	fail 'Docker builder is not the exact accepted Go image manifest'
[[ $(grep -Ec '^FROM golang:' "$dockerfile") -eq 1 ]] ||
	fail 'Dockerfile must contain exactly one Go builder identity'

[[ $(grep -Fxc '        uses: actions/setup-go@v2' "$release_workflow") -eq 1 ]] ||
	fail 'deactivated release workflow changed its setup-go action contract'
[[ $(grep -Fxc "          go-version: '1.26.7'" "$release_workflow") -eq 1 ]] ||
	fail 'deactivated release workflow does not declare exact Go 1.26.7'
[[ $(grep -Ec '^[[:space:]]+go-version:' "$release_workflow") -eq 1 ]] ||
	fail 'deactivated release workflow has duplicate Go identities'
if grep -Eq 'setup-go|go-version:' "$lint_workflow"; then
	fail 'deactivated lint workflow unexpectedly acquired a second Go identity'
fi

grep -Fq 'Requirement: [Go 1.26.7](https://go.dev/doc/install)' "$readme" ||
	fail 'README does not document the exact supported Go toolchain'

python3 - "$baseline" "$migration" <<'PY' || exit 1
import hashlib
import json
import sys

baseline_path, migration_path = sys.argv[1:]
with open(baseline_path, encoding="utf-8") as stream:
    baseline = json.load(stream)
with open(migration_path, encoding="utf-8") as stream:
    migration = json.load(stream)

def fail(message):
    print("toolchain declaration contract: " + message, file=sys.stderr)
    raise SystemExit(1)

if baseline["tool"]["go_build"]["version"] != "go version go1.26.7 darwin/arm64":
    fail("baseline Go identity does not match the declared toolchain")
baseline_sha = hashlib.sha256(open(baseline_path, "rb").read()).hexdigest()
if baseline_sha != migration["new"]["scorecard_sha256"]:
    fail("toolchain migration does not bind the current baseline scorecard")
if migration.get("schema_version") != 1:
    fail("toolchain migration schema is not pinned")
if migration["old"]["version"] != "go version go1.26.2 darwin/arm64":
    fail("toolchain migration lost the old instrument identity")
if migration["old"]["scorecard_source_commit"] != "097a9f15782c773e47a61d903f5d09f5c96408f0":
    fail("toolchain migration lost the pre-P7 baseline source commit")
if migration["new"]["version"] != "go version go1.26.7 darwin/arm64":
    fail("toolchain migration lost the new instrument identity")
if migration["selection"]["selected"]["version"] != "go1.26.7":
    fail("toolchain decision does not select Go 1.26.7")
if migration["selection"]["rejected"]["version"] != "go1.27.0":
    fail("toolchain decision does not record the rejected Go 1.27 line")
comparison = migration["comparison"]
if comparison != {
    "criteria_equal": True,
    "denominators_equal": True,
    "numeric_debt_equal": True,
    "numeric_debt_leaves": 228,
    "raw_report_body_sha256": "cf23c9dca987acd3a966693f933c4c7eca9f7d11c51f0fbabfe1d697f3d7497f",
    "unexpected_differences": [],
}:
    fail("toolchain comparison does not preserve the complete numeric baseline")
if migration["declarations"] != {
    "docker_builder": "golang:1.26.7-alpine3.24@sha256:28d89ee9cc0ff9fec75c82ca201e6bf7fdf9a679d4b7b24dfa04f2bb766bb468",
    "go_language": "1.18",
    "go_toolchain": "go1.26.7",
    "release_workflow": "1.26.7",
}:
    fail("toolchain migration metadata disagrees with contractual declarations")
if migration["module_selection"] != {
    "non_standard_packages": 143,
    "packages_equal": True,
    "selected_modules": 233,
    "selected_modules_equal": True,
    "tidy_projection": {
        "new_go_mod_checksums": [],
        "post_declaration_lines": 207,
        "pre_declaration_lines": 207,
        "version_changes": [],
    },
}:
    fail("toolchain migration does not explain module or tidy projection drift")
PY

printf 'toolchain declaration contract: PASS (module, Docker, deactivated workflows, docs, and exact baseline identity checked)\n'
