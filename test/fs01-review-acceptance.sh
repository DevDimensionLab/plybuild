#!/bin/sh
# FS01 candidate-bound acceptance: native review rounds on one existing Task
# through negative review, correction, rereview and agreed local Epic delivery.
# Fixtures are synthetic; they do not attest a real Claude or Codex journey.
set -eu
: "${PLY_CANDIDATE_BINARY:?native candidate binary required}"
: "${PLY_CANDIDATE_OID:?native candidate commit required}"
: "${PLY_CANDIDATE_TREE:?native candidate tree required}"
: "${HOME:?isolated native HOME required}"
: "${TMPDIR:?isolated native TMPDIR required}"
: "${PLY_FS01_GO:?recorded local Go executable required}"
: "${GOMODCACHE:?recorded module cache required}"
: "${PLY_MR05_SCHEMA_PYTHON:?recorded JSON Schema interpreter required}"
test -x "$PLY_CANDIDATE_BINARY"
test -x "$PLY_FS01_GO"
test -x "$PLY_MR05_SCHEMA_PYTHON"
# The worktree may be shared with a reviewer: any change fails, and is shown.
fs01_clean() {
    fs01_status=$(git status --porcelain=v1 --untracked-files=all)
    test -z "$fs01_status" && return 0
    printf 'Task worktree is not the clean candidate:\n%s\n' "$fs01_status"
    git diff --stat
    return 1
}
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
fs01_clean
fs01_test_root=$(mktemp -d /private/tmp/ply-fs01-tests.XXXXXX)
# Fixtures create native workspaces; they must be outside the real workspace.
trap 'fs01_exit=$?; if [ "$fs01_exit" -eq 0 ]; then rm -rf "$fs01_test_root"; else printf "Failed acceptance evidence: %s\n" "$fs01_test_root"; fi' EXIT
export TMPDIR="$fs01_test_root"
export GOCACHE="$fs01_test_root/go-cache" GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOENV=off GOWORK=off
# Go reads a private copy of the committed module files, read-only, so no
# acceptance step can rewrite the worktree's go.mod or go.sum.
mkdir "$fs01_test_root/mod"
git show "$PLY_CANDIDATE_OID:go.mod" > "$fs01_test_root/mod/go.mod"
git show "$PLY_CANDIDATE_OID:go.sum" > "$fs01_test_root/mod/go.sum"
export GOFLAGS="-mod=readonly -modfile=$fs01_test_root/mod/go.mod"
export GIT_AUTHOR_NAME="FS01 fixture" GIT_AUTHOR_EMAIL="fixture@example.invalid" GIT_COMMITTER_NAME="FS01 fixture" GIT_COMMITTER_EMAIL="fixture@example.invalid"
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
export PLY_TASK_RUN_TEST_BINARY="$PLY_CANDIDATE_BINARY"
# The exact candidate exposes the ordinary review surface and its contract.
"$PLY_CANDIDATE_BINARY" workflow execute review --help | grep -q -- '--coordinator'
"$PLY_CANDIDATE_BINARY" workflow execute review-return --help | grep -q -- '--context'
"$PLY_CANDIDATE_BINARY" capabilities --format json | grep -q '"workflow.execute.review-return"'
# Full cmd package: CLI surface, status text and the locked capabilities schema.
"$PLY_FS01_GO" test ./cmd -count=1 -timeout=300s
# Each selector runs in its own process: existing fixtures mutate CWD and env.
for fs01_selector in '^TestFS01' '^TestDeliveryRetire' '^TestAutomaticDeliveryIntegratesWithoutHumanRecord$' '^TestAutomaticDeliveryRevalidatesLiveInputsAtEveryGate$' '^TestDeliveryOwnerActualAcceptanceAndTruthfulProgress$' '^TestDeliveryReceipt'; do
    "$PLY_FS01_GO" test ./internal/taskrun -run "$fs01_selector" -count=1 -timeout=600s
done
git show "$PLY_CANDIDATE_OID:go.mod" | cmp -s - "$fs01_test_root/mod/go.mod"
git show "$PLY_CANDIDATE_OID:go.sum" | cmp -s - "$fs01_test_root/mod/go.sum"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
fs01_clean
printf 'FS01 candidate-bound acceptance passed\n'
