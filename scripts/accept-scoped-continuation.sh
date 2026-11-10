#!/bin/sh
set -eu
: "${PLY_CANDIDATE_BINARY:?Native verification must supply its exact candidate snapshot.}"
: "${PLY_CANDIDATE_OID:?Native verification must supply its candidate OID.}"
: "${PLY_CANDIDATE_TREE:?Native verification must supply its candidate tree.}"
: "${PLY_DEP02_OLD_BINARY:?Select the exact preserved original PWA control.}"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain)"
"$PLY_CANDIDATE_BINARY" version --format json > "$TMPDIR/candidate-version.json"
python3 - "$TMPDIR/candidate-version.json" "$PLY_CANDIDATE_OID" "$PLY_DEP02_OLD_BINARY" <<'PY'
import hashlib,json,pathlib,sys
v=json.loads(pathlib.Path(sys.argv[1]).read_text())
assert v['source_commit']==sys.argv[2] and v['dirty'] is False, 'Expected exact clean make candidate'
assert hashlib.sha256(pathlib.Path(sys.argv[3]).read_bytes()).hexdigest()=='7044a9996acd208a2bacf9bc849e3f486b72a765aec5f4b1feb888458934cb32', 'Original d58 control changed'
PY
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
# Only DEP02 and the old stop/resolution/retirement invariants; no full suite,
# browser, live provider, host process probe, coverage or installation gate.
go test -mod=readonly ./internal/taskrun -run '^(TestDEP02.*|TestDeliveryResolvingLatestStopDoesNotForgetEarlierDenial|TestDeliveryObservedStopKeepsLegacyStateReadable|TestDeliveryObservedStopInterruptedPublicationRetainsBlock|TestDeliveryObservedStopResolutionReplayWithPendingDenials|TestDeliveryReceiptResolutionPreservesUncertainExternalEffects|TestDeliveryResolveStopContinuesNamedChangedActionWithoutPermissionClaim|TestAL40RetainsUncertainEffects|TestDeliveryRetireStopsOnActiveWriterOrUnknownEffect|TestDeliveryOwnerUnknownVerificationNeverReplays|TestDeliveryOwnerFailedVerifierReviewAndCandidateDrift)$' -count=1
