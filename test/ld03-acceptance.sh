#!/bin/sh
# Candidate-bound retention and same-Run closeout consumer checks.
set -eu
: "${PLY_CANDIDATE_BINARY:?native candidate binary required}"
: "${PLY_CANDIDATE_OID:?native candidate commit required}"
: "${PLY_CANDIDATE_TREE:?native candidate tree required}"
: "${PLY_LD03_GO:?cached Go executable required}"
: "${PLY_INFORMATIVE_OLD_READER:?preserved base reader required}"
: "${GOMODCACHE:?existing module cache required}"
: "${GOCACHE:?existing build cache required}"
test -x "$PLY_CANDIDATE_BINARY"
test -x "$PLY_LD03_GO"
test -x "$PLY_INFORMATIVE_OLD_READER"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain=v1 --untracked-files=all)"
ld03_tmp=$(mktemp -d /private/tmp/ply-ld03-acceptance.XXXXXX)
export TMPDIR="$ld03_tmp"
export GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOENV=off GOWORK=off
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
# Preserve the actual disposable workspaces for readback; no real Task is used.
export PLY_INFORMATIVE_JOURNEY_ROOT="$ld03_tmp"
"$PLY_LD03_GO" test ./internal/workspace -run '^(TestCloseoutRetention|TestAutomaticCloseoutRetention|TestTaskCloseout)' -count=1 -v -timeout=180s
"$PLY_LD03_GO" test ./internal/taskrun -run '^TestInformativeHistoricalReportCLI' -count=1 -v -timeout=240s
for ld03_selector in '^TestAutomaticDeliveryCloseoutInterruptionResumesWithoutIntegration$' '^TestAutomaticDeliveryCLICloseoutRecoveryFromReleasedOwner$' '^TestAutomaticAcceptanceCLIAfterMergeRecovery$'; do
    "$PLY_LD03_GO" test ./internal/taskrun -run "$ld03_selector" -count=1 -v -timeout=180s
done
git diff --check
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain=v1 --untracked-files=all)"
printf 'LD03 candidate retention/closeout acceptance passed; synthetic consumers: %s\n' "$ld03_tmp"
