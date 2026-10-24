#!/bin/sh
# T965226AF: focused producer checks; no live agents, install or server activation.
set -eu
: "${PLY_CANDIDATE_BINARY:?exact make-built candidate is required}"
: "${PLY_CANDIDATE_OID:?exact source commit is required}"
: "${PLY_CANDIDATE_TREE:?exact source tree is required}"
"$PLY_CANDIDATE_BINARY" version --format json > "$TMPDIR/tab-label-version.json"
python3 - "$TMPDIR/tab-label-version.json" <<'PY'
import json, os, sys
v=json.load(open(sys.argv[1]))
assert v['source_commit']==os.environ['PLY_CANDIDATE_OID'] and v['dirty'] is False, v
print(json.dumps(v, sort_keys=True))
PY
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain)"
go version -m "$PLY_CANDIDATE_BINARY"
go test -race ./internal/agentwork -count=1 -timeout 120s
go test ./cmd -run '^TestAgentWork' -count=1 -timeout 120s
go test ./internal/taskrun -run '^(TestStatusFollowupCandidateJourney|TestInventoryStatusFollowupVeto|TestReadInventoryStopHistoryIdentityContracts|TestReadInventoryDeliveryReportsKeepExactQuestionAndSeparateNativeOutcomes|TestReadInventoryDeliveryCorruptionAndOrphansPreserveValidSiblings)$' -count=1 -timeout 180s
go vet ./internal/agentwork ./cmd
printf '%s\n' 'T965226AF producer acceptance passed; installed timings and consumer QA remain parent-owned.'
