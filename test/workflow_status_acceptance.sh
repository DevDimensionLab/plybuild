#!/bin/sh
# Technical acceptance only. Preserve all generated evidence and fixture state.
set -eu
export PYTHONDONTWRITEBYTECODE=1
task_repo=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)
cd "$task_repo"
task_python=${PLY_TEST_PYTHON:-python3}
task_go=${GO:-go}
task_evidence=$(mktemp -d "${TMPDIR:-/tmp}/ply-workflow-status-acceptance.XXXXXX")
printf 'Workflow status technical acceptance evidence: %s\n' "$task_evidence"
"$task_python" -c 'from jsonschema import Draft202012Validator'
"$task_python" test/workflow_status_go_tests.py "$task_evidence" "$task_go"
"$task_go" build -o "$task_evidence/ply" ./cmd/ply
"$task_python" test/read_contract_schema_test.py
"$task_python" test/workspace_views_schema_test.py
"$task_python" test/workflow_status_schema_test.py
if [ -n "${PLY_STATUS_BASELINE:-}" ]; then
    "$task_python" test/workflow_status_acceptance.py --ply "$task_evidence/ply" --baseline "$PLY_STATUS_BASELINE" --root "$task_evidence/journey"
else
    "$task_python" test/workflow_status_acceptance.py --ply "$task_evidence/ply" --root "$task_evidence/journey"
fi
"$task_go" test ./internal/workspaceview -run '^$' -bench BenchmarkWorkflowStatusMultiProject -benchtime=3x -count=1
git diff --check
printf 'Technical acceptance passed; human QA is not performed by this script. Evidence: %s\n' "$task_evidence"
