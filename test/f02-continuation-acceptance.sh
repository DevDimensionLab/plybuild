#!/bin/sh
# TB58A978A: registered continuation reuse and retained delivery boundaries.
set -eu
: "${PLY_CANDIDATE_BINARY:?exact candidate is required}"
: "${PLY_CANDIDATE_OID:?candidate commit is required}"
: "${PLY_CANDIDATE_TREE:?candidate tree is required}"
"$PLY_CANDIDATE_BINARY" version --format json > "$TMPDIR/f02-version.json"
python3 - "$TMPDIR/f02-version.json" "$PLY_CANDIDATE_OID" <<'PY'
import json, sys
identity = json.load(open(sys.argv[1]))
assert identity['source_commit'] == sys.argv[2], identity
assert identity['dirty'] is False, identity
PY
printf 'TB58A978A: exact built candidate identity checked\n'
go test ./internal/taskrun -v -count=1 -timeout=10m -run '^(TestAL40RetiredContinuationRejectsBeforeValidation|TestDeliveryContinuationConcurrentAndInterruptedPublication|TestDeliveryContinuationEquivalentObservationKeepsOneControl|TestDeliveryContinuationRecoversBeforeStateAfterIncompleteReport|TestDeliveryContinuationRecoversPublishedProofAfterIncompleteReport|TestDeliveryContinuationRecoversPublishedProofAfterObservationExpires|TestDeliveryContinuationRejectsTamperedPriorControlAndAuthority|TestDeliveryContinuationRejectsUnknownVerificationAndStaleSession|TestDeliveryOwnerUnknownVerificationNeverReplays|TestDeliveryProviderUpgradeRetainsAuthorityAndLiveOwnerGuards|TestDeliveryReceiptLaterRuntimeRefusalPreservesStop|TestPlyUpgradeStoredContinuation|TestPublishedContinuationRecoveryChangedCaller|TestPublishedContinuationRecoveryGuardsUnderLock|TestPublishedContinuationRecoveryRejectsDifferentRuntimeBinding|TestPublishedContinuationRecoveryUsesOriginalExpiredRuntime|TestR5CallbackBindingsAndImmutableAcceptance|TestRegisteredContinuationAddsMissingRuntimeEvidence|TestRegisteredContinuationReusesAcrossCLIUpgrades|TestRegisteredContinuationValidatesExplicitRuntimeEvidence)$'
printf 'Retained concrete regression tests passed; removed process journeys are excluded.\n'
