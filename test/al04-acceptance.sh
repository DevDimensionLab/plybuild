#!/bin/sh
set -eu
: "${PLY_CANDIDATE_BINARY:?exact private candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate OID required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse HEAD^{tree})" = "$PLY_CANDIDATE_TREE"
"$PLY_CANDIDATE_BINARY" workflow run report --help >/dev/null
"$PLY_CANDIDATE_BINARY" workflow run show --help >/dev/null
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
export GOTOOLCHAIN=local
export GOFLAGS=-mod=readonly
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
cached_go=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin
export PATH="$cached_go:$PATH"
test "$(go env GOVERSION)" = go1.26.7
go test -count=1 -timeout 6m -v ./internal/taskrun -run '^(TestAL04IdentityHistoryAndArtifactGuardsAfterDrift|TestAL04NegativeAndStoppedHistoryIsNotRevived|TestAL04PartialPromptEffectAfterDrift|TestAL04PreLaunchFactsStayEffective|TestAL04ProviderDriftAfterLaunch|TestAL04ReportOwnerCompletesRoundAfterPlyUpdate)$'
go test -count=1 -timeout 12m -v ./internal/taskrun -run '^(TestAL01CoordinatorDenialBlocksReceiptWithoutImpersonation|TestAL01PreAcceptanceStop|TestAL01PreStartStopIsBoundImmutableAndNonQualifying|TestAL01ReportCLIJourney|TestAL01ReportReceiptCannotChangeNativeTerminalContract|TestAL01ReportReceiptUnknownFactsAndCompletedReview|TestDeliveryProviderUpgradeDoesNotGrantNewLaunch|TestDeliveryProviderUpgradePreservesAcceptedContinuation|TestDeliveryProviderUpgradeRequiresCurrentRuntimeObservation|TestWorkflowProviderAcceptanceRequiresFreshIdentity|TestWorkflowProviderRoundTrip|TestWorkflowRuntimeDriftAndNegativeAuthority)$'
go vet ./internal/taskrun ./cmd
# The whole Task diff from its AL03 base.
git diff --check f5535ef88b8bf2e2b75720285234c2fa71070a34 "$PLY_CANDIDATE_OID"
printf 'Retained concrete regression tests passed; removed process journeys are excluded.\n'
