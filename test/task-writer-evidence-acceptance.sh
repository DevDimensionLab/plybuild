#!/bin/sh
set -eu
: "${PLY_CANDIDATE_BINARY:?exact private candidate executable required}"
: "${PLY_CANDIDATE_OID:?exact candidate commit required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
ar03_go=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go
export PLY_MR05_SCHEMA_PYTHON=/Users/perottochristensen/.local/share/ply/gs01-20261009/schema-venv/bin/python
ar03_scratch=$(mktemp -d /private/tmp/ply-ar03-acceptance.XXXXXX)
chmod 700 "$ar03_scratch"
export TMPDIR="$ar03_scratch" GOTMPDIR="$ar03_scratch"
export GOFLAGS=-mod=readonly GOPROXY=off GOSUMDB=off GOENV=off GOWORK=off GOTOOLCHAIN=local
export GOPATH=/Users/perottochristensen/go GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
export PATH="$(dirname "$ar03_go"):$PATH"
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
export PLY_TASK_RUN_TEST_BINARY="$PLY_CANDIDATE_BINARY"
export PLY_BIN="$PLY_CANDIDATE_BINARY"
test -x "$PLY_CANDIDATE_BINARY"
"$PLY_MR05_SCHEMA_PYTHON" -I -B -c 'import jsonschema'
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
# Writer composition, real local observers, distinct caller policies, historical
"$ar03_go" test ./internal/taskrun -run '^(TestAR03RetirementDoesNotReplayDeniedHerdrInspection|TestCodexResumeWriterHelper|TestDeliveryRetireEndsSelectedOwnerAndFreesQueueWithoutOwnerCallback|TestDeliveryRetireStopsOnActiveWriterOrUnknownEffect|TestFS02ReviewRetirePreservesEarlierUnknownEffect|TestFS02ReviewRetireRecoversDurableEventAfterStateFailure|TestOwnerResumeAfterStartupRecoveryKeepsRecoveryHistory|TestSystemLocalProcessesObservesRealWriters)$' -count=1 -timeout=600s
"$ar03_go" vet ./internal/taskrun
git diff --check
test -z "$(git status --porcelain)"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
printf 'Retained concrete regression tests passed; removed process journeys are excluded.\n'
