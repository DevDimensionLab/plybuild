#!/usr/bin/env bash
set -euo pipefail
export LC_ALL=C LANG=C GOENV=off GOWORK=off GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local
repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
# One owned run dir holds the journey and all Go temp files; it is removed on every
# exit. The journey reports its result and any failing command streams on stdout/stderr.
. "$repo_root/test/roundtrip_root.sh"
roundtrip_root 'workspace task queue roundtrip' ply-task-queue
go_bin=${GO:-go}
cp "$repo_root/test/fixtures/workspace_task_queue/lifecycle_helper.go.txt" "$run_root/lifecycle_helper.go"
(cd "$repo_root" && "$go_bin" build -o "$run_root/ply" ./cmd/ply)
"$go_bin" build -o "$run_root/lifecycle-helper" "$run_root/lifecycle_helper.go"
python3 "$repo_root/test/fixtures/workspace_task_queue/journey.py" "$run_root"
printf 'workspace task queue roundtrip: PASS\n'
