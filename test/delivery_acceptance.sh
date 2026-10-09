#!/bin/sh
set -eu

delivery_repo_root=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd -P)
# Native automatic verification supplies TMPDIR inside the Run. Disposable Ply
# workspace fixtures must live outside the real workspace to avoid nesting it.
delivery_test_root=$(mktemp -d "/tmp/ply-delivery-acceptance.XXXXXX")
delivery_test_root=$(CDPATH='' cd -- "$delivery_test_root" && pwd -P)
trap 'delivery_status=$?; if [ "$delivery_status" -eq 0 ]; then rm -rf "$delivery_test_root"; else printf "Acceptance did not pass; temporary evidence retained at %s\n" "$delivery_test_root"; fi' EXIT
trap 'exit 130' HUP INT TERM
cd "$delivery_repo_root"
export TMPDIR="$delivery_test_root"

delivery_candidate_binary=${PLY_CANDIDATE_BINARY:-}
if [ -z "$delivery_candidate_binary" ]; then
    delivery_candidate_binary="$delivery_test_root/ply"
    go build -o "$delivery_candidate_binary" ./cmd/ply
fi
test -x "$delivery_candidate_binary"
go test -c -o "$delivery_test_root/taskrun.test" ./internal/taskrun
# Taskrun fixtures change process CWD/environment and exceed a single package's
# time budget. Separate processes preserve isolation and run every listed test.
export PLY_DELIVERY_TEST_BINARY="$delivery_candidate_binary"
export PLY_TASK_RUN_TEST_BINARY="$delivery_candidate_binary"
python3 ./test/delivery_acceptance.py "$delivery_repo_root" "$delivery_test_root"
go vet ./...
git diff --check
