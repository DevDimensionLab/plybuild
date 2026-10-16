#!/usr/bin/env bash
set -euo pipefail
export LC_ALL=C LANG=C GOENV=off GOWORK=off GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local
repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
temp_root=$(mktemp -d "${TMPDIR:-/tmp}/ply-queue-removed-source.XXXXXX")
temp_root=$(cd "$temp_root" && pwd -P)
go_bin=${GO:-go}
printf 'QR-01 fixture: %s\n' "$temp_root"
cp "$repo_root/test/fixtures/workspace_task_queue/lifecycle_helper.go.txt" "$temp_root/lifecycle_helper.go"
(cd "$repo_root" && "$go_bin" build -o "$temp_root/ply" ./cmd/ply)
"$go_bin" build -o "$temp_root/lifecycle-helper" "$temp_root/lifecycle_helper.go"
python3 "$repo_root/test/fixtures/workspace_task_queue/journey.py" "$temp_root" --removed-source
