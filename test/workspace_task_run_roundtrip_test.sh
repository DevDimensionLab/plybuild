#!/usr/bin/env bash
set -euo pipefail
export LC_ALL=C LANG=C GOENV=off GOWORK=off GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local
repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
temp_root=$(mktemp -d "${TMPDIR:-/tmp}/ply-task-run.XXXXXX")
temp_root=$(cd "$temp_root" && pwd -P)
go_bin=${GO:-go}
# Preserve the newly built CLI and its synthetic journey under the caller's temp
# root. There is no real Codex/Claude launch, TTY device access or product bypass.
(cd "$repo_root" && "$go_bin" build -o "$temp_root/ply" ./cmd/ply)
(cd "$repo_root" && PLY_TASK_RUN_TEST_BINARY="$temp_root/ply" "$go_bin" test ./internal/taskrun -run '^TestR13FreshBinaryCallbacks$' -count=1 -v)
printf 'workspace task run roundtrip: PASS (synthetic CLI; binary: %s)\n' "$temp_root/ply"
