#!/bin/sh
# Run through scripts/with-go so that go is the go.mod toolchain.
set -eu
: "${PLY_CANDIDATE_BINARY:?exact candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate commit required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
export GOTOOLCHAIN=local
"$PLY_CANDIDATE_BINARY" version --format json | python3 -c 'import json,os,sys; x=json.load(sys.stdin); assert x["source_commit"]==os.environ["PLY_CANDIDATE_OID"] and not x["dirty"]'
test "$(git rev-parse HEAD^{tree})" = "$PLY_CANDIDATE_TREE"
test/go_test_each.sh ./internal/taskwork TestWorkItemSession TestWorkItemRegistry TestWorkItemStore TestGoal TestQueue TestTaskQueue TestRegistry TestReadContract TestTaskScoped -- -mod=readonly -timeout=15m -count=1
go vet -mod=readonly ./internal/workspace ./internal/taskwork ./cmd
git diff --check
