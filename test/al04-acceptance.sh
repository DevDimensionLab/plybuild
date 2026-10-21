#!/bin/sh
set -eu
: "${PLY_CANDIDATE_BINARY:?exact private candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate OID required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse HEAD^{tree})" = "$PLY_CANDIDATE_TREE"
"$PLY_CANDIDATE_BINARY" workflow run report --help >/dev/null
"$PLY_CANDIDATE_BINARY" workflow execute show --help >/dev/null
# The candidate CLI journey runs this exact binary after a real AL03-base install.
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
export GOTOOLCHAIN=local
export GOFLAGS=-mod=readonly
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
cached_go=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin
export PATH="$cached_go:$PATH"
test "$(go env GOVERSION)" = go1.26.7
# Explicit timeouts stay inside the 1800s acceptance budget: about 4x the
# measured 72-80s, 182s and <1s on 3744fc4 (review r001-f01).
go test -count=1 -timeout 6m -v ./internal/taskrun -run '^TestAL04'
# Consumers of the changed report callback and freshness boundary, the pre-launch
# trust contrast, and Delivery contrasts. Unchanged retirement and effect
# reconciliation branches stay with their own Task acceptance.
go test -count=1 -timeout 12m -v ./internal/taskrun -run '^(TestAL01(ReportReceiptUnknownFactsAndCompletedReview|ReportCLIJourney|PreAcceptanceStop|PreStartStopIsBoundImmutableAndNonQualifying|CoordinatorDenialBlocksReceiptWithoutImpersonation|ReportReceiptCannotChangeNativeTerminalContract)|TestWorkflow(ProviderRoundTrip|ProviderAcceptanceRequiresFreshIdentity|RuntimeDriftAndNegativeAuthority|CodexTrustDriftAfterReservationWithholdsGrant|CodexTrustCandidateCLIJourney)|TestDelivery(ProviderUpgradeDoesNotGrantNewLaunch|ProviderUpgradeRequiresCurrentRuntimeObservation|ProviderUpgradePreservesAcceptedContinuation|ReceiptCLIJourney))$'
go test -count=1 -timeout 2m ./cmd -run '^Test(ExecuteBegin|ExecuteCheckIn|ExecuteRetire|DeliveryReceipt|WorkflowRun|WorkflowExecute)'
go vet ./internal/taskrun ./cmd
# The whole Task diff from its AL03 base.
git diff --check f5535ef88b8bf2e2b75720285234c2fa71070a34 "$PLY_CANDIDATE_OID"
