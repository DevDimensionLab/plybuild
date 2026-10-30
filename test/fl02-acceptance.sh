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
: "${PLY_MR05_SCHEMA_PYTHON:?recorded JSON Schema interpreter required}"
test -x "$PLY_MR05_SCHEMA_PYTHON"
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
"$PLY_CANDIDATE_BINARY" workflow run check-in --help
"$PLY_FL02_GO" test ./internal/workspace -run '^(TestAutomaticCloseoutRetentionPreservesExactOriginalReferences|TestAutomaticNativeIntegrationRecoversAfterGitWithoutHumanQA|TestDeliveryAgreementHumanOwnershipIsExplicitAndVersioned|TestDeliveryAgreementStrictShapeAndHistoricalAbsence|TestDeliveryAuthorizationRegistryShapePreservesLegacyAndReceipt|TestRegistryRoundtripPreservesAutomaticAuthorizationAndHistory)$' -count=1 -timeout=180s
"$PLY_FL02_GO" test ./internal/taskrun -run '^(TestDeliveryBeginBindsExistingAuthorizedNotificationNatively|TestDeliveryBeginCannotCreateSecondFrozenNotificationBinding|TestDeliveryBeginNativeAcceptanceAndLostReply|TestDeliveryBeginUnconfirmedAndWrongOwnerCannotAccept|TestDeliveryContinuationRejectsTamperedPriorControlAndAuthority|TestDeliveryContinuationRejectsUnknownVerificationAndStaleSession|TestDeliveryObservedStopInterruptedPublicationRetainsBlock|TestDeliveryObservedStopKeepsLegacyStateReadable|TestDeliveryObservedStopRequiresOwnerResolutionBeforeProgress|TestDeliveryObservedStopResolutionReplayWithPendingDenials|TestDeliveryReceiptDoesNotDowngradeLegacyContract|TestDeliveryReceiptDoesNotRequireBlanketPermission|TestDeliveryReceiptLaterRuntimeRefusalPreservesStop|TestDeliveryReceiptReportedRootsNeverAttestSandbox|TestDeliveryReceiptResolutionPreservesUncertainExternalEffects|TestDeliveryReceiptResolutionSelectsFrozenStartupEntrypoint|TestDeliveryReceiptReuseRejectsChangedInputsWithoutRunningAcceptance|TestDeliveryResolvingLatestStopDoesNotForgetEarlierDenial|TestDeliveryStartupDoesNotConfuseLaunchAcceptanceAndWork)$' -count=1 -timeout=180s
git diff --check
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain=v1 --untracked-files=all)"
printf 'Retained startup, refusal and Git-evidence tests passed\n'
printf 'Retained concrete regression tests passed; removed process journeys are excluded.\n'
