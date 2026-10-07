#!/bin/sh
set -eu

delivery_repo_root=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd -P)
delivery_test_root=$(mktemp -d "${TMPDIR:-/tmp}/ply-delivery-acceptance.XXXXXX")
delivery_test_root=$(CDPATH='' cd -- "$delivery_test_root" && pwd -P)
trap 'delivery_status=$?; if [ "$delivery_status" -eq 0 ]; then rm -rf "$delivery_test_root"; else printf "Acceptance did not pass; temporary evidence retained at %s\n" "$delivery_test_root"; fi' EXIT
trap 'exit 130' HUP INT TERM
cd "$delivery_repo_root"

go build -o "$delivery_test_root/ply" ./cmd/ply
go test -c -o "$delivery_test_root/taskrun.test" ./internal/taskrun
# Taskrun fixtures change process CWD/environment and exceed a single package's
# time budget. Separate processes preserve isolation and run every listed test.
export PLY_DELIVERY_TEST_BINARY="$delivery_test_root/ply"
export PLY_TASK_RUN_TEST_BINARY="$delivery_test_root/ply"
python3 ./test/delivery_acceptance.py "$delivery_repo_root" "$delivery_test_root"
go vet ./...
git diff --check
