#!/bin/sh
# AL05 acceptance: ordinary report, evidence, acceptance and review inputs
# continue at 0644 while private roles stay owner-only and every retained
# refusal names its input. Runs against the exact candidate commit and binary.
set -eu
: "${PLY_CANDIDATE_BINARY:?exact private candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate OID required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
BASE=7e049edfc69057b6f15e1187b2740cb47861f4bc
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain)"
git merge-base --is-ancestor "$BASE" HEAD
"$PLY_CANDIDATE_BINARY" workflow execute report --help >/dev/null
"$PLY_CANDIDATE_BINARY" workflow execute verify --help >/dev/null
# The candidate CLI journey runs this exact binary through report, verify and
# local delivery.
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
export GOTOOLCHAIN=local
export GOFLAGS=-mod=readonly
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
cached_go=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin
export PATH="$cached_go:$PATH"
test "$(go env GOVERSION)" = go1.26.7

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

# Explicit timeouts stay inside the 1800s budget: about 4x the measured ~80s
# for the AL05 tests, ~105s for the contrasts and ~1s for cmd.
group al05 -timeout 6m -v ./internal/taskrun -run '^TestAL05'
for name in \
	TestAL05OrdinaryReportAndEvidenceAreConsumed \
	TestAL05RetainedReportRequirementsNameInputAndCause \
	TestAL05BaselineOrdinaryCandidateBinary \
	TestAL05OrdinaryAcceptanceAndReviewQualify \
	TestAL05ChangedOrdinaryInputsDoNotQualify \
	TestAL05InPlaceOverwriteAfterBindingDoesNotQualify \
	TestAL05NativePrivateRolesStayProtected \
	TestAL05CandidateCLIJourneyWithOrdinaryModes; do
	grep -q "^--- PASS: $name " "$TMPDIR/al05.log"
done
# Existing consumers of the changed report, verify, receipt-reuse, status and
# begin (notification route) readers.
group contrasts -timeout 8m -v ./internal/taskrun -run '^(TestDeliveryOwnerActualAcceptanceAndTruthfulProgress|TestAutomaticVerificationMissingInputsCannotReusePriorPass|TestDeliveryReceiptReuseRejectsChangedInputsWithoutRunningAcceptance|TestDeliveryRecordedVerificationRequalificationDoesNotRepeatAcceptance|TestDeliveryStatusPreservesCompletedReceiptBeforeQualification|TestDeliveryBeginBindsExistingAuthorizedNotificationNatively|TestDeliveryBeginNativeAcceptanceAndLostReply|TestAutomaticAcceptanceCLIGreenAndFailureCorrection)$'
test "$(grep -c '^--- PASS' "$TMPDIR/contrasts.log")" = 8
# Unchanged CLI input and private-input boundaries.
group cmd -timeout 2m -v ./cmd -run '^(TestDeliveryCLIRejectsAmbiguousOrMissingInputBeforeEffects|TestExecuteAcceptanceRequiresExactPrivateInputsBeforeWorkspaceAccess|TestExecuteCheckInHasNoBlanketPermissionFlag)$'
test "$(grep -c '^--- PASS' "$TMPDIR/cmd.log")" = 3
go vet ./internal/taskrun ./cmd
git diff --check "$BASE" "$PLY_CANDIDATE_OID"
echo "AL05 acceptance passed"
