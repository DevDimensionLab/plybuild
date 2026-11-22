#!/bin/sh
# Run through scripts/with-go so that go is the go.mod toolchain.
set -eu
: "${PLY_CANDIDATE_BINARY:?exact packaged candidate required}"
: "${PLY_CANDIDATE_OID:?exact candidate commit required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
: "${PLY_MR05_SCHEMA_PYTHON:?existing schema Python required}"
: "${TMPDIR:?private scratch parent required}"
for executable in "$PLY_CANDIDATE_BINARY" "$PLY_MR05_SCHEMA_PYTHON"; do
    case "$executable" in /*) ;; *) echo "absolute executable required: $executable" >&2; exit 1 ;; esac
    test -x "$executable"
done
"$PLY_MR05_SCHEMA_PYTHON" -I -B -c 'import jsonschema; print("schema Python:", jsonschema.__file__)'
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain)"
scratch=$(mktemp -d "$TMPDIR/TB2546BF5.XXXXXX")
export TMPDIR="$scratch" GOTMPDIR="$scratch"
export GOFLAGS=-mod=readonly GOPROXY=off GOSUMDB=off GOENV=off GOWORK=off GOTOOLCHAIN=local
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
go version -m "$PLY_CANDIDATE_BINARY"
"$PLY_MR05_SCHEMA_PYTHON" -I -B test/read_contract_schema_test.py
test/go_test_each.sh ./cmd 'TestHI07ArgumentErrorsPrintNoResult$' 'TestHI07ContractBundleIsLockedAndMatchesActualOutput$' 'TestHI07InputCapabilitiesMatchActualCommands$' 'TestHI07ValidRequestExamplesCarryNativeDigests$' 'TestHI07WriteContractExamples$' 'TestRegistryFailureProducesExactlyOneJSONDiagnostic$' -- -count=1 -timeout=180s
test/go_test_each.sh ./internal/taskwork 'TestRegistryDiagnosticsUseDeclaredFormatAndValidation$' -- -count=1 -timeout=180s
go vet ./cmd
git diff --check
test -z "$(git status --porcelain)"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test "$(shasum -a 256 "$PLY_CANDIDATE_BINARY")" = "$binary_before"
printf 'TB2546BF5 focused acceptance passed; full package timeout remains separate. Evidence scratch: %s\n' "$scratch"
printf 'Retained concrete regression tests passed; removed process journeys are excluded.\n'
