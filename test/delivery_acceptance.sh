#!/bin/sh
set -eu

delivery_repo_root=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd -P)
delivery_test_root=$(mktemp -d "${TMPDIR:-/tmp}/ply-delivery-acceptance.XXXXXX")
trap 'rm -rf "$delivery_test_root"' EXIT HUP INT TERM
cd "$delivery_repo_root"

go build -o "$delivery_test_root/ply" ./cmd/ply
# The workflow suite includes real Git repositories and bounded provider
# protocol scenarios. Its full run exceeds Go's default ten-minute timeout.
PLY_DELIVERY_TEST_BINARY="$delivery_test_root/ply" go test ./... -count=1 -timeout=30m
go vet ./...
git diff --check
