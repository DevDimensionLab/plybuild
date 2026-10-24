#!/bin/sh
# TC0E6E1FE: synthetic consumer fixtures; never a claim of live provider QA.
set -eu
: "${PLY_CANDIDATE_BINARY:?exact candidate is required}"
: "${PLY_CANDIDATE_OID:?candidate commit is required}"
: "${PLY_CANDIDATE_TREE:?candidate tree is required}"
# Keep Unix sockets short and isolate all writable compiler/test state.
: "${TMPDIR:?private physical scratch is required}"
export GOCACHE="$TMPDIR/gocache"
export GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local
"$PLY_CANDIDATE_BINARY" version --format json > "$TMPDIR/pv01-version.json"
python3 - "$TMPDIR/pv01-version.json" "$PLY_CANDIDATE_OID" <<'PY'
import json, sys
identity=json.load(open(sys.argv[1]))
assert identity['source_commit']==sys.argv[2], identity
assert identity['dirty'] is False, identity
PY
printf 'TC0E6E1FE: exact built candidate identity checked\n'
go test ./internal/taskexecute ./internal/taskrun ./internal/weekrun -count=1 -timeout=5m -run 'TestPlyUpgrade|TestBoundStartRejectsInvalidCurrentInputsBeforeAnyIntent|TestBoundStartInterruptedAttemptNeverResumesOnReplay|TestBoundStartNativeLockAndReservationRacesRejectChangedBasis|TestBoundStartRejectsNativeWriteBetweenContentAndWorktree|TestR5CallbackBindingsAndImmutableAcceptance|TestAL04IdentityHistoryAndArtifactGuardsAfterDrift|TestDeliveryOwnerUnknownVerificationNeverReplays|TestDeliveryReceiptLaterRuntimeRefusalPreservesStop|TestAutomaticDeliveryRevalidatesLiveInputsAtEveryGate|TestAutomaticDeliveryRuntimeDriftInvalidatesCurrentAcceptance|TestDeliveryContinuationRejectsUnknownVerificationAndStaleSession|TestWeekPlyUpgrade|TestWeekDetachedRecoveryUsesFrozenLoginEnvironment|TestMR05WeekControlChildAndLegacyAuthority|TestDeliveryContinuationRejectsTamperedPriorControlAndAuthority|TestParallelActivationHumanReleaseRequiresEveryRemainingControl|TestMR05ControlStaleBasisAndABA|TestMR05ControlFaultsPreserveUnknown|TestParallelNativeTwoOwners'
printf 'TC0E6E1FE: upgrade consumers and retained authority/effect/candidate guards passed\n'
