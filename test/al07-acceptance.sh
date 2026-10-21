#!/bin/sh
set -eu
: "${PLY_CANDIDATE_BINARY:?exact private candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate OID required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse HEAD^{tree})" = "$PLY_CANDIDATE_TREE"
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
export GOTOOLCHAIN=local GOFLAGS=-mod=readonly
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
cached_go=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin
export PATH="$cached_go:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin"
# Named consumed readbacks and real candidate CLI processes; no provider/network.
go test -v -count=1 -timeout=240s ./internal/taskrun -run '^(TestAL07StartupActionMatchesFrozenProtocol|TestAL07ConsumedInventoryRecovery|TestAL07CandidateCLIReadActionsDoNotWrite|TestAL07RetiredReviewHistoryHasNoAvailableReviewer|TestHumanCloseoutRetainsOnlyItsOwnSourceReservation)$'
go test -v -count=1 -timeout=60s ./internal/workspaceview -run '^(TestAL07RequiredInstallationFailureRemainsVisibleAcrossWorkspaceViews|TestAL07CloseoutPreservesExplicitHumanWait|TestAL07MissingOrChangedNativeActorStaysUnknown|TestCloseoutProgressDistinguishesRetiredKeptAndPending|TestRegisteredProgressSeparatesTechnicalHumanAndIntegrationFacts|TestProcessProgressAndHumanWaitNeverManufactureQAPass)$'
# Preserve existing runtime, human-gate, refusal and unknown-effect boundaries.
go test -v -count=1 -timeout=360s ./internal/taskrun -run '^(TestDeliveryReceiptDoesNotRequireBlanketPermission|TestAutomaticDeliveryPreservesRequiredHumanAndIntegrationOwnerGates|TestAutomaticDeliveryLatestNativeHumanRejectionBlocks|TestOwnerResumeBoundariesStopBeforeEffects|TestReadInventoryStopHistoryIdentityContracts|TestReadInventoryNativeRetirementNeedsNoAgentReport)$'
go vet ./internal/taskrun ./internal/workspaceview
git diff --check
"$PLY_CANDIDATE_BINARY" workflow execute check-in --help >/dev/null
"$PLY_CANDIDATE_BINARY" workflow execute review-return --help >/dev/null
