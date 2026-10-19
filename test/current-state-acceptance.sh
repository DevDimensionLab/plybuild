#!/bin/sh
# AR01 deterministic acceptance; provider facts come from the owned stand-in.
set -eu
: "${PLY_CANDIDATE_BINARY:?exact private candidate executable is required}"
: "${PLY_MR05_SCHEMA_PYTHON:?existing schema interpreter is required}"
test -x "$PLY_CANDIDATE_BINARY"
"$PLY_MR05_SCHEMA_PYTHON" -I -B -c 'import jsonschema'
ar01_go=${AR01_GO:-go}
ar01_scratch=$(mktemp -d /private/tmp/ply-ar01-state.XXXXXX)
chmod 700 "$ar01_scratch"
export TMPDIR="$ar01_scratch"
export GOTMPDIR="$ar01_scratch"
export GOFLAGS=-mod=readonly
export GOPROXY=off
export GOSUMDB=off
export GOENV=off
export GOWORK=off
export GOTOOLCHAIN=local
export PLY_BIN="$PLY_CANDIDATE_BINARY"
"$ar01_go" test ./internal/currentstate ./internal/taskrun ./internal/weekrun ./internal/taskexecute ./cmd \
    -run 'Test(AR01|AW01Guidance|NA01|SL01|MR05)' -count=1 -timeout=600s
"$ar01_go" vet ./internal/currentstate ./internal/taskrun ./internal/weekrun ./internal/taskexecute ./cmd
git diff --check
printf 'AR01 state regressions and exact-candidate public CLI journey passed. Fixtures: %s\n' "$ar01_scratch"
