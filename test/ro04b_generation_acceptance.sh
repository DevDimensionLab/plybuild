#!/bin/sh
set -eu
: "${PLY_CANDIDATE_BINARY:?exact candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate commit required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
RO04B_GO=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go
export GOTOOLCHAIN=local
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
"$PLY_CANDIDATE_BINARY" version --format json | python3 -c 'import json,os,sys; x=json.load(sys.stdin); assert x["source_commit"]==os.environ["PLY_CANDIDATE_OID"] and not x["dirty"]'
test "$(git rev-parse HEAD^{tree})" = "$PLY_CANDIDATE_TREE"
"$RO04B_GO" test -mod=readonly ./internal/workspace -run '^Test(WorkItemSession|WorkItemRegistry|WorkItemStore|Goal|TaskGoal|Queue|TaskQueue|Registry|ReadContract|TaskScoped|TaskSpecGuard|StartupDependency|TaskPreparationP5|EpicBaseB2)' -timeout=15m -count=1
"$RO04B_GO" test -mod=readonly ./internal/taskrun -run '^Test(WorkflowCodexTrustLegacyAndConflict|Goal|StartupDependency)' -timeout=10m -count=1
"$RO04B_GO" vet -mod=readonly ./internal/workspace ./internal/taskrun ./cmd
git diff --check
