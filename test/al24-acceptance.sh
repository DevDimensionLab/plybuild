#!/bin/sh
# AL24: the CLI surface regression, and a real successful surface run that removes
# its own run root and leaves an unrelated sibling alone. Run from the repository
# root with the selected go first in PATH (scripts/with-go) and TMPDIR outside Git
# and Ply workspaces.
set -eu
root=$(mktemp -d "${TMPDIR:?set TMPDIR outside Git and Ply workspaces}/ply-al24-acceptance.XXXXXX")
trap 'rm -rf -- "$root"' EXIT
bash -n test/cli_surface_contract_test.sh
bash -n test/cli_surface_evidence_test.sh
sh -n test/al24-acceptance.sh
mkdir "$root/regression" "$root/regression-fixtures"
PLY_ROUNDTRIP_ROOT="$root/regression-fixtures" bash test/cli_surface_evidence_test.sh "$root/regression"
mkdir "$root/success"
printf 'unrelated\n' >"$root/success/sentinel"
PLY_ROUNDTRIP_ROOT="$root/success" bash test/cli_surface_contract_test.sh
test "$(cat "$root/success/sentinel")" = unrelated
test "$(find "$root/success" -mindepth 1 -print)" = "$root/success/sentinel"
printf 'real-success cleanup: own root removed; sibling retained\n'
printf 'AL24 acceptance passed: surface regression and real CLI success removed only its own root.\n'
