#!/bin/sh
set -eu
: "${PLY_CANDIDATE_BINARY:?exact packaged candidate required}"
: "${PLY_CANDIDATE_OID:?exact candidate commit required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
: "${PLY_REGRESSION_GO:?existing Go 1.26.7 executable required}"
: "${PLY_MR05_SCHEMA_PYTHON:?existing schema Python required}"
for executable in "$PLY_CANDIDATE_BINARY" "$PLY_REGRESSION_GO" "$PLY_MR05_SCHEMA_PYTHON"; do
    case "$executable" in /*) ;; *) echo "absolute executable required: $executable" >&2; exit 1 ;; esac
    test -x "$executable"
done
test "$("$PLY_REGRESSION_GO" version)" = 'go version go1.26.7 darwin/arm64'
"$PLY_MR05_SCHEMA_PYTHON" -I -B -c 'import jsonschema; print("schema Python:", jsonschema.__file__)'
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain)"
scratch=$(mktemp -d "${TMPDIR:-/tmp}/TB2546BF5.XXXXXX")
export TMPDIR="$scratch" GOTMPDIR="$scratch"
export GOFLAGS=-mod=readonly GOPROXY=off GOSUMDB=off GOENV=off GOWORK=off GOTOOLCHAIN=local
export PATH="$(dirname "$PLY_REGRESSION_GO"):$PATH"
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY" PLY_TASK_RUN_TEST_BINARY="$PLY_CANDIDATE_BINARY"
export PLY_SL01_SCHEMA_PYTHON="$PLY_MR05_SCHEMA_PYTHON"
binary_before=$(shasum -a 256 "$PLY_CANDIDATE_BINARY")
"$PLY_CANDIDATE_BINARY" version --format json > "$scratch/version.json"
"$PLY_MR05_SCHEMA_PYTHON" -I -B - "$scratch/version.json" "$PLY_CANDIDATE_OID" <<'PY'
import json, sys
version = json.load(open(sys.argv[1]))
assert version['source_commit'] == sys.argv[2], version
assert version['dirty'] is False, version
print('TB2546BF5 packaged candidate:', version)
PY
"$PLY_REGRESSION_GO" version -m "$PLY_CANDIDATE_BINARY"
"$PLY_REGRESSION_GO" test ./cmd -run '^(TestRegistryFailureProducesExactlyOneJSONDiagnostic|TestMR05.*|TestSL01AnswerCLI.*|TestHI07.*)$' -count=1 -timeout=180s
"$PLY_REGRESSION_GO" test ./internal/workspace -run '^TestRegistry(DiagnosticsUseDeclaredFormatAndValidation|RoundtripPreservesAutomaticAuthorizationAndHistory|FormatChangeRequiresNativePermit)$' -count=1 -timeout=180s
"$PLY_REGRESSION_GO" test ./internal/taskexecute -run '^(TestExecuteKeepsEarlierManualLaunchWhenAutoBecomesDefault|TestSharedEpicUniqueSpecPublicCLIPreparesAndPreservesRetry|TestBackupStartRejectsChangedRuntimeBeforeRunPublication|TestBackupStartPreRunReadbackShowsPreservedEffectiveProfile|TestExecuteRejectsChangedPreservedIntent|TestMR05ProjectionCandidateCLI)$' -count=1 -timeout=180s
"$PLY_REGRESSION_GO" test ./internal/taskrun -run '^(TestAL09.*|TestDeliveryCLIBinaryBuildIgnoresFixtureCWD)$' -count=1 -timeout=180s
"$PLY_REGRESSION_GO" vet ./cmd ./internal/taskexecute ./internal/taskrun
git diff --check
test -z "$(git status --porcelain)"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test "$(shasum -a 256 "$PLY_CANDIDATE_BINARY")" = "$binary_before"
printf 'TB2546BF5 focused acceptance passed; full package timeout remains separate. Evidence scratch: %s\n' "$scratch"
