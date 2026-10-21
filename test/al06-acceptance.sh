#!/bin/sh
# AL06 acceptance: new interactive Claude report Runs use the selected manual
# or auto mode through preview, launch, receipt, correction and report return;
# unsupported modes fail before effects and actual refusals stay preserved.
# Runs against the exact candidate commit and binary.
set -eu
# Same tools as the development regressions: cached Go first, then Homebrew Git
# before the system tools.
cached_go=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin
export PATH="$cached_go:/opt/homebrew/bin:$PATH"
: "${PLY_CANDIDATE_BINARY:?exact private candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate OID required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
BASE=c02e3ee474fb94da61a4d375f6b0d437c1c5c4df
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain)"
git merge-base --is-ancestor "$BASE" HEAD
"$PLY_CANDIDATE_BINARY" workflow run start --help >/dev/null
"$PLY_CANDIDATE_BINARY" workflow run check-in --help >/dev/null
# The candidate CLI journey executes this exact binary; it is only read and run.
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
export GOTOOLCHAIN=local
export GOFLAGS=-mod=readonly
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
test "$(go env GOVERSION)" = go1.26.7
test "$(command -v git)" = /opt/homebrew/bin/git

# Runs one go test group, keeps go's actual exit and prints the raw log on
# failure, or a short PASS summary on success.
group() {
	name=$1
	shift
	log="$TMPDIR/$name.log"
	if go test -count=1 "$@" >"$log" 2>&1; then
		grep -E '^(--- PASS|ok)' "$log" || true
	else
		status=$?
		echo "== $name failed with exit $status; raw log:" >&2
		cat "$log" >&2
		exit "$status"
	fi
}

group al06 -timeout 6m -v ./internal/taskrun -run '^TestAL06'
for name in \
	TestAL06ClaudeReportModeReceiptCorrectionAndReturn \
	TestAL06ClaudeAutoLostStartAndPreservedRefusal \
	TestAL06ClaudeReportUnsupportedModeFailsBeforeEffects \
	TestAL06ClaudeAutoReportCLIJourney \
	TestAL06ClaudeReceiptPreservesUnknownOrDifferentObservation; do
	grep -q "^--- PASS: $name " "$TMPDIR/al06.log"
done
# Named compatibility boundaries of the changed mode routing: existing Claude
# manual and Codex report starts, AL01 receipt with unknown facts and preserved
# coordinator denial, and Claude delivery mode routing.
group contrasts -timeout 8m -v ./internal/taskrun -run '^(TestWorkflowProviderRoundTrip|TestWorkflowProviderRejectsInvalidSelectionBeforeEffects|TestWorkflowProviderDoesNotOpenNativeTransports|TestAL01ReportReceiptUnknownFactsAndCompletedReview|TestAL01CoordinatorDenialBlocksReceiptWithoutImpersonation|TestAL01ReportReceiptCannotChangeNativeTerminalContract|TestDeliveryClaudePermissionModeUsesBoundRequest)$'
test "$(grep -c '^--- PASS: Test' "$TMPDIR/contrasts.log")" = 7
go vet ./internal/taskrun ./cmd
git diff --check "$BASE" "$PLY_CANDIDATE_OID"
echo "AL06 acceptance passed"
