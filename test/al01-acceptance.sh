#!/bin/sh
set -eu
: "${PLY_CANDIDATE_BINARY:?exact private candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate OID required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse HEAD^{tree})" = "$PLY_CANDIDATE_TREE"
"$PLY_CANDIDATE_BINARY" workflow run check-in --help >/dev/null
"$PLY_CANDIDATE_BINARY" workflow run retire --help >/dev/null
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
export GOTOOLCHAIN=local
export GOFLAGS=-mod=readonly
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
cached_go=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin
export PATH="$cached_go:$PATH"
go test -count=1 -v ./internal/taskrun -run '^TestAL01'
go test -count=1 ./internal/workflowhandoff
go test -count=1 ./internal/taskrun -run '^(TestWorkflowRuntimeDriftAndNegativeAuthority|TestWorkflowProviderRoundTrip|TestWorkflowProviderAcceptanceRequiresFreshIdentity|TestDeliveryReceipt|TestDeliveryRetire|TestOwnerResumeBoundariesStopBeforeEffects)'
go test -count=1 ./cmd -run '^Test(ExecuteBegin|ExecuteCheckIn|ExecuteRetire|DeliveryReceipt|WorkflowRun)'
go vet ./internal/taskrun ./internal/workflowhandoff ./cmd
git diff --check
