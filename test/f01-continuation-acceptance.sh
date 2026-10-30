#!/bin/sh
# T93C57C19: protocol regressions in isolated native Git/Herdr fixtures.
set -eu
: "${PLY_CANDIDATE_BINARY:?exact candidate is required}"
: "${PLY_CANDIDATE_OID:?candidate commit is required}"
: "${PLY_CANDIDATE_TREE:?candidate tree is required}"
"$PLY_CANDIDATE_BINARY" version --format json > "$TMPDIR/f01-version.json"
python3 - "$TMPDIR/f01-version.json" "$PLY_CANDIDATE_OID" <<'PY'
import json, sys
identity=json.load(open(sys.argv[1]))
assert identity['source_commit']==sys.argv[2], identity
assert identity['dirty'] is False, identity
PY
printf 'T93C57C19: exact built candidate identity checked\n'
go test ./internal/taskrun -v -count=1 -timeout=10m -run '^(TestAL40RetiredContinuationRejectsBeforeValidation|TestDeliveryContinuationConcurrentAndInterruptedPublication|TestDeliveryContinuationEquivalentObservationKeepsOneControl|TestDeliveryContinuationRecoversBeforeStateAfterIncompleteReport|TestDeliveryContinuationRecoversPublishedProofAfterIncompleteReport|TestDeliveryContinuationRecoversPublishedProofAfterObservationExpires|TestDeliveryContinuationRejectsTamperedPriorControlAndAuthority|TestDeliveryContinuationRejectsUnknownVerificationAndStaleSession|TestDeliveryOwnerUnknownVerificationNeverReplays|TestDeliveryReceiptLaterRuntimeRefusalPreservesStop|TestPlyUpgradeStoredContinuation|TestPublishedContinuationRecoveryChangedCaller|TestPublishedContinuationRecoveryGuardsUnderLock|TestPublishedContinuationRecoveryRejectsDifferentRuntimeBinding|TestPublishedContinuationRecoveryUsesOriginalExpiredRuntime|TestR5CallbackBindingsAndImmutableAcceptance)$'
printf 'Retained concrete regression tests passed; removed process journeys are excluded.\n'
