#!/bin/sh
set -eu

delivery_repo_root=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd -P)
# This exact Task alone has the preserved AW-01 focused acceptance agreement.
delivery_task_ref=$(GIT_OPTIONAL_LOCKS=0 git -c core.fsmonitor=false -c core.hooksPath=/dev/null -C "$delivery_repo_root" symbolic-ref --quiet HEAD 2>/dev/null || true)
if [ "$delivery_task_ref" = 'refs/heads/ply-task/aw-01-current-resume-guidance-038968865f8f' ]; then
    cd "$delivery_repo_root"
    exec /bin/sh ./test/resume-next-action-acceptance.sh
fi
# MR-05 replacement only; root must match this expected ref to fresh native preview.
if [ "$delivery_task_ref" = 'refs/heads/ply-task/mr-05-bound-runtime-control-replacement-ae7cec4dc362' ]; then
    cd "$delivery_repo_root"
    exec /bin/sh ./test/mr05-acceptance.sh
fi
# MR-05 project/repo projection correction only; fresh native preview must match.
if [ "$delivery_task_ref" = 'refs/heads/ply-task/mr-05-standalone-project-binding-correct-20951f628816' ]; then
    cd "$delivery_repo_root"
    exec /bin/sh ./test/mr05-projection-acceptance.sh
fi
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
