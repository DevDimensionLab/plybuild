#!/bin/sh
# AR02: deterministic Worker/Week regressions with modeled provider facts.
set -eu
: "${PLY_CANDIDATE_BINARY:?exact private candidate executable is required}"
: "${PLY_MR05_SCHEMA_PYTHON:?existing schema interpreter is required}"
test -x "$PLY_CANDIDATE_BINARY"
"$PLY_MR05_SCHEMA_PYTHON" -I -B -c 'import jsonschema'
ar02_go=${AR02_GO:-go}
ar02_scratch=$(mktemp -d /private/tmp/ply-ar02-acceptance.XXXXXX)
chmod 700 "$ar02_scratch"
export TMPDIR="$ar02_scratch"
export GOTMPDIR="$ar02_scratch"
export GOFLAGS=-mod=readonly
export GOPROXY=off
export GOSUMDB=off
export GOENV=off
export GOWORK=off
export GOTOOLCHAIN=local
export PLY_BIN="$PLY_CANDIDATE_BINARY"
"$ar02_go" test ./internal/currentstate ./internal/taskrun ./internal/weekrun ./internal/taskexecute ./cmd \
    -run 'Test(AR02|AR01|Worker|Week|AW01|NA01|SL01|MR05)' -count=1 -timeout=600s
"$ar02_go" test ./internal/workspace -run 'Test.*(Queue|Preparation)' -count=1 -timeout=600s
"$ar02_go" vet ./internal/currentstate ./internal/taskrun ./internal/weekrun ./internal/taskexecute ./internal/workspace ./cmd
git diff --check
printf 'AR02 outcome/release, AR01, queue and exact-candidate CLI checks passed. Scratch: %s\n' "$ar02_scratch"
