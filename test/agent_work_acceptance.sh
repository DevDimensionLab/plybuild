#!/bin/sh
# T7CF442AE: bounded producer conformance; no live Herdr/provider/server.
set -eu
: "${PLY_CANDIDATE_BINARY:?an exact make-built candidate is required}"
export GOTOOLCHAIN=local
"$PLY_CANDIDATE_BINARY" version --format json > "$TMPDIR/agent-work-candidate-version.json"
python3 - "$TMPDIR/agent-work-candidate-version.json" <<'PY'
import json,os,sys
v=json.load(open(sys.argv[1]))
assert v['source_commit'] and v['dirty'] is False, v
if os.environ.get('PLY_CANDIDATE_OID'):
    assert v['source_commit']==os.environ['PLY_CANDIDATE_OID'],v
print(json.dumps(v,sort_keys=True))
PY
go version
go version -m "$PLY_CANDIDATE_BINARY"
go test ./internal/agentwork -count=1 -timeout 120s
go test ./cmd -run '^(TestAgentWork|TestWorkflowStatus|TestCapabilities)' -count=1 -timeout 120s
go test ./internal/taskrun -run '^(TestStatusFollowupCandidateJourney|TestInventoryStatusFollowupVeto|TestReadInventoryStopHistoryIdentityContracts|TestReadInventoryDeliveryReportsKeepExactQuestionAndSeparateNativeOutcomes|TestReadInventoryDeliveryCorruptionAndOrphansPreserveValidSiblings)$' -count=1 -timeout 180s
