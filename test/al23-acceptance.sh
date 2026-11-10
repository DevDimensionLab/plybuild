#!/bin/sh
set -eu
export PATH=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin
export GOTOOLCHAIN=local GOFLAGS=-mod=readonly
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
: "${PLY_CANDIDATE_BINARY:?exact candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate OID required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain)"
git merge-base --is-ancestor 796383021488c7d4a6c153802146f7b7fa3642e5 HEAD
test "$(go env GOVERSION)" = go1.26.7

# Print raw streams and actual exit into the native verifier's retained output.
# The CLI test likewise logs each subprocess's exact argv, cwd and streams.
run() {
    name=$1
    shift
    printf 'argv:'
    printf ' <%s>' "$@"
    printf '\n'
    result=0
    "$@" >"$TMPDIR/al23-$name.stdout" 2>"$TMPDIR/al23-$name.stderr" || result=$?
    cat "$TMPDIR/al23-$name.stdout"
    cat "$TMPDIR/al23-$name.stderr" >&2
    printf 'exit=%s\n' "$result"
    test "$result" = 0
}
run version "$PLY_CANDIDATE_BINARY" version --format json
run build-metadata go version -m "$PLY_CANDIDATE_BINARY"
run hash shasum -a 256 "$PLY_CANDIDATE_BINARY"
run focused go test -count=1 -timeout 5m -v ./internal/taskrun -run '^(TestAL23SameInodeReadBindsConsumedBytes|TestAL23ReadInstabilityControls|TestAL23CandidateCLIRejectsFinalBytesForMixedDigest|TestAL05RetainedReportRequirementsNameInputAndCause|TestAL05OrdinaryAcceptanceAndReviewQualify|TestAL05InPlaceOverwriteAfterBindingDoesNotQualify|TestAL05NativePrivateRolesStayProtected)$'
for name in TestAL23SameInodeReadBindsConsumedBytes TestAL23ReadInstabilityControls TestAL23CandidateCLIRejectsFinalBytesForMixedDigest TestAL05RetainedReportRequirementsNameInputAndCause TestAL05OrdinaryAcceptanceAndReviewQualify TestAL05InPlaceOverwriteAfterBindingDoesNotQualify TestAL05NativePrivateRolesStayProtected; do
    grep -q "^--- PASS: $name " "$TMPDIR/al23-focused.stdout"
done
run vet go vet ./internal/taskrun ./cmd
run diffcheck git diff --check 796383021488c7d4a6c153802146f7b7fa3642e5 "$PLY_CANDIDATE_OID"
echo 'AL23 acceptance passed: truthful mixed-byte binding, rejected final bytes, stable continuation, retained protections and snapshots.'
