#!/bin/sh
set -eu
cached_go=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin
export PATH="$cached_go:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin"
export GOTOOLCHAIN=local GOFLAGS=-mod=readonly
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
: "${PLY_CANDIDATE_BINARY:?exact private candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate OID required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
test -x "$PLY_CANDIDATE_BINARY"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain)"
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
# Durable raw CLI streams/exits, including rejected historical/guard commands.
# Provider/Herdr are modeled; native callbacks, Git and candidate CLI are real.
PLY_AL14_EVIDENCE_DIR=$(mktemp -d "$(dirname "$PLY_CANDIDATE_BINARY")/al14-cli.XXXXXX")
export PLY_AL14_EVIDENCE_DIR
printf 'Raw CLI evidence: %s\n' "$PLY_AL14_EVIDENCE_DIR"
log="$PLY_AL14_EVIDENCE_DIR/tests.stdout"
if go test -v -count=1 -timeout=240s ./internal/taskrun -run '^(TestAL14CLIPostStartRetainsWorkAndHistory|TestAL14PostStartGuardsAndRecovery|TestAL14CLIReconcilesClosedPreparation|TestAL14BlockedNativeTerminalRetirementDiagnosisPreservesHistory|TestAL01RetirementGuardsAndDurableRetry|TestAL01ReportEffectReconciliationPreservesUncertaintyAndSource|TestDeliveryRetireStopsOnActiveWriterOrUnknownEffect)$' >"$log" 2>"$PLY_AL14_EVIDENCE_DIR/tests.stderr"; then
 printf '0\n' >"$PLY_AL14_EVIDENCE_DIR/tests.exit"
else
 result=$?
 printf '%s\n' "$result" >"$PLY_AL14_EVIDENCE_DIR/tests.exit"
 cat "$log" "$PLY_AL14_EVIDENCE_DIR/tests.stderr"
 exit "$result"
fi
cat "$log" "$PLY_AL14_EVIDENCE_DIR/tests.stderr"
for name in TestAL14CLIPostStartRetainsWorkAndHistory TestAL14PostStartGuardsAndRecovery TestAL14CLIReconcilesClosedPreparation TestAL14BlockedNativeTerminalRetirementDiagnosisPreservesHistory TestAL01RetirementGuardsAndDurableRetry TestAL01ReportEffectReconciliationPreservesUncertaintyAndSource TestDeliveryRetireStopsOnActiveWriterOrUnknownEffect; do
 grep -q "^--- PASS: $name " "$log"
done
go test -v -count=1 ./cmd -run '^TestAL14ReportRetirementWarningsStayOnStderr$'
go vet ./internal/taskrun ./cmd
git diff --check
test -z "$(git status --porcelain)"
