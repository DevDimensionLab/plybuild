#!/bin/sh
set -eu
export LC_ALL=C LANG=C
cb01a_root=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd -P)
test "$#" -eq 2 || { printf 'Usage: %s <release-package-directory> <manifest-sha256>\n' "$0" >&2; exit 2; }
export PLY_REGISTRY_TEST_RELEASE="$1"
export PLY_REGISTRY_TEST_MANIFEST="$2"
export PLY_CANDIDATE_BINARY="$1/ply"
cd "$cb01a_root"
python3 test/registry_cli_acceptance.py "$1" "$2"
go test ./internal/workspace ./internal/clirelease ./cmd \
  -run 'TestRegistry|TestInstall|TestExecutableEntryBoundary' -count=1
test/delivery_acceptance.sh
