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
go test ./internal/taskrun ./internal/taskexecute -v -count=1 -timeout=10m -run 'TestPublishedContinuationRecovery|TestDeliveryContinuation|TestPlyUpgradeStoredContinuation|TestBoundStartPublicCLIContract|TestR5CallbackBindingsAndImmutableAcceptance|TestDeliveryOwnerUnknownVerificationNeverReplays|TestAutomaticDeliveryRevalidatesLiveInputsAtEveryGate|TestAutomaticDeliveryRuntimeDriftInvalidatesCurrentAcceptance|TestDeliveryReceiptLaterRuntimeRefusalPreservesStop|TestAL40RetiredContinuationRejectsBeforeValidation'
printf 'T93C57C19: published recovery, CLI orchestration and retained guards passed\n'
go test ./cmd -count=1 -timeout=2m -run 'TestExecute(CheckIn|Begin|ResolveStop|Review|Acceptance)|TestWorkflowExecuteVerifyHelpManagedRequirements'
printf 'T93C57C19: command contracts passed\n'
