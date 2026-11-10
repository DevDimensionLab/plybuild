#!/bin/sh
set -eu
cached_go=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin
export PATH="$cached_go:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin"
export GOTOOLCHAIN=local GOFLAGS=-mod=readonly
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
: "${PLY_CANDIDATE_BINARY:?exact private candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate OID required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
test -x "$PLY_CANDIDATE_BINARY"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain)"
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
# Persist raw CLI stdout/stderr/exits with the native acceptance evidence.
export PLY_AL13_EVIDENCE_DIR="$TMPDIR/al13-cli"
mkdir -p "$PLY_AL13_EVIDENCE_DIR"
# Exact candidate processes; provider and Herdr observations are modeled fixtures.
# Named compatibility cases preserve prior sealed history and short returns.
log="$TMPDIR/al13-tests.log"
if go test -v -count=1 -timeout=240s ./internal/taskrun -run '^(TestAL13ReviewReturnProseParser|TestAL13ReviewReturnProseThroughCandidateCLI|TestAL13ReviewReturnRejectsInvalidCLIInput|TestFS01ReviewRoundsThroughCandidateCLI|TestRV02EvidenceBlockedReassessmentThroughCandidateCLI)$' >"$log" 2>&1; then
 cat "$log"
else
 result=$?
 cat "$log"
 exit "$result"
fi
# A matching command with no tests is not acceptance.
for name in TestAL13ReviewReturnProseParser TestAL13ReviewReturnProseThroughCandidateCLI TestAL13ReviewReturnRejectsInvalidCLIInput TestFS01ReviewRoundsThroughCandidateCLI TestRV02EvidenceBlockedReassessmentThroughCandidateCLI; do
 grep -q "^--- PASS: $name " "$log"
done
go vet ./internal/taskrun ./cmd
git diff --check
test -z "$(git status --porcelain)"
