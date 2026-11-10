#!/usr/bin/env bash
set -euo pipefail
export LC_ALL=C LANG=C GOENV=off GOWORK=off GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local
repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
. "$repo_root/test/roundtrip_root.sh"
roundtrip_root 'workspace task queue removed source' ply-queue-removed-source
temp_root=$run_root
go_bin=${GO:-go}
cp "$repo_root/test/fixtures/workspace_task_queue/lifecycle_helper.go.txt" "$temp_root/lifecycle_helper.go"
(cd "$repo_root" && "$go_bin" build -o "$temp_root/ply" ./cmd/ply)
"$go_bin" build -o "$temp_root/lifecycle-helper" "$temp_root/lifecycle_helper.go"
python3 "$repo_root/test/fixtures/workspace_task_queue/journey.py" "$temp_root" --removed-source
