#!/usr/bin/env bash
# go_test_each.sh checks each pattern argument against the listed test names only,
# so a pattern that matches the package summary line or one alternation branch that
# selects nothing is judged by the names it actually selects. Run with the selected
# go first in PATH (scripts/with-go) and TMPDIR outside Git and Ply workspaces.
set -euo pipefail
script=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)/go_test_each.sh
root=$(mktemp -d "${TMPDIR:?set TMPDIR outside Git and Ply workspaces}/go-test-each.XXXXXX")
trap 'rm -rf -- "$root"' EXIT
mkdir "$root/module"
printf 'module example.test/each\n\ngo 1.26\n' >"$root/module/go.mod"
printf 'package each\n\nimport "testing"\n\nfunc TestAlpha(t *testing.T) {}\n' >"$root/module/each_test.go"
failed=0
# check <name> <expected exit> <text expected in output> <pattern>...
check() {
	local name=$1 want=$2 text=$3 got=0
	shift 3
	(cd "$root/module" && GOWORK=off sh "$script" . "$@" -- -count=1 -v) >"$root/$name.out" 2>&1 || got=$?
	if [[ $got != "$want" ]] || ! grep -Fq -- "$text" "$root/$name.out"; then
		printf 'FAIL %s: exit %s (want %s), output lacks %s or differs:\n' "$name" "$got" "$want" "$text" >&2
		sed 's/^/  /' "$root/$name.out" >&2
		failed=1
	else
		printf 'ok   %s\n' "$name"
	fi
}
check existing-name 0 '--- PASS: TestAlpha' 'TestAlpha$'
check missing-name 1 'go_test_each: TestMissing selects no test in .' 'TestMissing'
check summary-line-is-no-test 1 'selects no test' '^ok[[:space:]]'
check alternation-is-one-pattern 0 '--- PASS: TestAlpha' 'TestAlpha$|TestMissing'
check every-argument-must-select 1 'go_test_each: TestMissing selects no test in .' 'TestAlpha' 'TestMissing'
exit "$failed"
