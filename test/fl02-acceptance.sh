#!/bin/sh
# Candidate-bound startup, refusal, compatibility and native local delivery checks.
set -eu
: "${PLY_CANDIDATE_BINARY:?native candidate binary required}"
: "${PLY_CANDIDATE_OID:?native candidate commit required}"
: "${PLY_CANDIDATE_TREE:?native candidate tree required}"
: "${HOME:?isolated native HOME required}"
: "${TMPDIR:?isolated native TMPDIR required}"
: "${PLY_FL02_GO:?recorded local Go executable required}"
: "${GOMODCACHE:?recorded module cache required}"
test -x "$PLY_CANDIDATE_BINARY"
test -x "$PLY_FL02_GO"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain=v1 --untracked-files=all)"
fl02_test_root=$(mktemp -d /private/tmp/ply-fl02-tests.XXXXXX)
# Fixtures create native workspaces; they must be outside the real workspace.
trap 'fl02_exit=$?; if [ "$fl02_exit" -eq 0 ]; then rm -rf "$fl02_test_root"; else printf "Failed acceptance evidence: %s\n" "$fl02_test_root"; fi' EXIT
export TMPDIR="$fl02_test_root"
export GOCACHE="$fl02_test_root/go-cache" GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOENV=off GOWORK=off
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
export PLY_TASK_RUN_TEST_BINARY="$PLY_CANDIDATE_BINARY"
"$PLY_CANDIDATE_BINARY" workflow execute check-in --help
"$PLY_FL02_GO" test ./cmd -run '^(TestExecute(CheckIn|Begin)|TestStartupCapabilities)' -count=1 -timeout=120s
"$PLY_FL02_GO" test ./internal/workflowhandoff -run '^(TestDelivery|Test.*Sandbox)' -count=1 -timeout=180s
"$PLY_FL02_GO" test ./internal/workspace -run '^(TestDeliveryAgreement|TestDeliveryAuthorization|Test.*Automatic.*Retention)' -count=1 -timeout=180s
"$PLY_FL02_GO" test ./internal/taskexecute -count=1 -timeout=180s
# Each selector runs in its own process: existing fixtures mutate CWD and env.
for fl02_selector in '^TestDeliveryReceipt' '^TestDeliveryBegin' '^TestDeliveryStartup' '^TestDelivery(ResolvingLatestStop|ObservedStop)' '^TestDeliveryContinuationPreservesAcceptedSessionAndFrozenInputs$' '^TestDeliveryContinuationRejects' '^TestAutomaticDeliveryRevalidatesLiveInputsAtEveryGate$' '^TestAutomaticDeliveryRuntimeDriftInvalidatesCurrentAcceptance$' '^TestWorkerNativeReviewRejectsDeveloperIdentityAndForgedReview$' '^TestWorkerRecoveryCannotConsumeAnotherCandidate$'; do
    "$PLY_FL02_GO" test ./internal/taskrun -run "$fl02_selector" -count=1 -timeout=180s
done
git diff --check
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain=v1 --untracked-files=all)"
printf 'FL02 candidate-bound acceptance passed\n'
