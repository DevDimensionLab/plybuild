#!/bin/sh
set -eu
export PATH=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin
export GOTOOLCHAIN=local GOFLAGS=-mod=readonly
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
: "${PLY_CANDIDATE_BINARY:?exact candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate OID required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain)"
test "$(go env GOVERSION)" = go1.26.7
# Native verification removes its TMPDIR after exit. Keep raw child evidence
# in a separate private physical root so failed runs survive that cleanup.
artifact_root=$(mktemp -d /private/tmp/ply-al24-acceptance.XXXXXX)
printf 'AL24 raw evidence: %s\n' "$artifact_root"
run() {
    name=$1
    shift
    printf 'argv:' >"$artifact_root/$name.argv"
    printf ' <%s>' "$@" >>"$artifact_root/$name.argv"
    printf '\n' >>"$artifact_root/$name.argv"
    cat "$artifact_root/$name.argv"
    result=0
    "$@" >"$artifact_root/$name.stdout" 2>"$artifact_root/$name.stderr" || result=$?
    printf '%s\n' "$result" >"$artifact_root/$name.exit"
    cat "$artifact_root/$name.stdout"
    cat "$artifact_root/$name.stderr" >&2
    printf 'exit=%s\n' "$result"
    test "$result" = 0
}
run version "$PLY_CANDIDATE_BINARY" version --format json
run build-metadata go version -m "$PLY_CANDIDATE_BINARY"
run hash shasum -a 256 "$PLY_CANDIDATE_BINARY"
run contract-syntax bash -n test/cli_surface_contract_test.sh
run regression-syntax bash -n test/cli_surface_evidence_test.sh
run acceptance-syntax sh -n test/al24-acceptance.sh
run regression bash test/cli_surface_evidence_test.sh "$artifact_root/regression"
# This remains the existing real script, with its own fresh build and isolation.
# Keep an unrelated sibling to prove that successful cleanup owns only its root.
mkdir "$artifact_root/success-tmp"
printf 'unrelated\n' >"$artifact_root/success-tmp/sentinel"
run real-success env TMPDIR="$artifact_root/success-tmp" bash test/cli_surface_contract_test.sh
test "$(cat "$artifact_root/success-tmp/sentinel")" = unrelated
test "$(find "$artifact_root/success-tmp" -mindepth 1 -print)" = "$artifact_root/success-tmp/sentinel"
printf 'real-success cleanup: own root removed; sibling retained\n'
run diffcheck git diff --check 51d4bcf3c427c666029041d7b16903c1d5de5856 "$PLY_CANDIDATE_OID"
printf 'AL24 acceptance passed: retained failure streams and original exits; real CLI success cleaned its own root.\n'
