#!/bin/sh
set -eu
: "${PLY_REGRESSION_GO:?existing absolute Go required}"
: "${PLY_MR05_SCHEMA_PYTHON:?existing isolated schema interpreter required}"
: "${GOOS:?explicit target required}"
: "${GOARCH:?explicit architecture required}"
: "${CGO_ENABLED:?explicit cgo required}"
"$PLY_MR05_SCHEMA_PYTHON" -I -B test/coverage_preflight_test.py
scratch=$(mktemp -d "${TMPDIR:-/tmp}/T778B627C.XXXXXX")
export PATH="$(dirname "$PLY_REGRESSION_GO"):$PATH"
export GOENV=off GOWORK=off GOFLAGS=-mod=readonly GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY" PLY_TASK_RUN_TEST_BINARY="$PLY_CANDIDATE_BINARY"
export PLY_SL01_SCHEMA_PYTHON="$PLY_MR05_SCHEMA_PYTHON"
"$PLY_MR05_SCHEMA_PYTHON" -I -B test/coverage-preflight.py --output "$scratch/before.json" -- \
    "$PLY_REGRESSION_GO" test -json -count=1 -timeout=300s "-coverprofile=$scratch/coverage.out" \
    ./internal/taskrun ./internal/workspace -run '^(TestBackupGoal.*|TestRecordedProgress.*|TestAL05CandidateCLIJourneyWithOrdinaryModes|TestAL16UnknownIdentityQualifiedInstall|TestQualifiedCLIInstallationPublicConsumer|TestAdoptionSourceWriterObservations|TestAR03ComposedWriterEvidence|TestAR03HeadlessCleanupUsesActualSession|TestWorkflowDirtyContentDriftAndSymlinkAreNotTrusted)$'
set +e
"$PLY_REGRESSION_GO" test -json -count=1 -timeout=300s "-coverprofile=$scratch/coverage.out" \
    ./internal/taskrun ./internal/workspace -run '^(TestBackupGoal.*|TestRecordedProgress.*|TestAL05CandidateCLIJourneyWithOrdinaryModes|TestAL16UnknownIdentityQualifiedInstall|TestQualifiedCLIInstallationPublicConsumer|TestAdoptionSourceWriterObservations|TestAR03ComposedWriterEvidence|TestAR03HeadlessCleanupUsesActualSession|TestWorkflowDirtyContentDriftAndSymlinkAreNotTrusted)$' \
    > "$scratch/test.jsonl" 2> "$scratch/test.stderr"
result=$?
set -e
cat "$scratch/test.jsonl"
cat "$scratch/test.stderr" >&2
"$PLY_REGRESSION_GO" vet ./internal/taskrun
git diff --check
"$PLY_MR05_SCHEMA_PYTHON" -I -B test/coverage-preflight.py --after "$scratch/before.json" \
    --test-exit "$result" --output "$scratch/after.json"
printf 'T778B627C focused runtime acceptance passed; evidence: %s\n' "$scratch"
