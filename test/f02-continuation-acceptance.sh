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
go test ./internal/taskrun ./internal/taskexecute -v -count=1 -timeout=10m -run 'TestRegisteredContinuation|TestPublishedContinuationRecovery|TestDeliveryContinuation|TestDeliveryProviderUpgradeRetainsAuthorityAndLiveOwnerGuards|TestOwnerResumeBoundariesStopBeforeEffects|TestPlyUpgradeStoredContinuation|TestR5CallbackBindingsAndImmutableAcceptance|TestDeliveryOwnerUnknownVerificationNeverReplays|TestDeliveryReceiptLaterRuntimeRefusalPreservesStop|TestAL40RetiredContinuationRejectsBeforeValidation'
printf 'TB58A978A: continuation reuse, F01 recovery, owner resume and retained guards passed\n'
go test ./cmd -count=1 -timeout=2m -run 'TestExecute(CheckIn|Begin|ResolveStop)'
printf 'TB58A978A: command contracts passed\n'
