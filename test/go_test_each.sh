#!/bin/sh
# Run the Go tests in one package whose names match each pattern, and fail
# before running when a pattern selects no test: a stale selector is no test.
# Each pattern is a regular expression anchored at the start of a top-level test
# name; end it with $ to select exactly one test. Run from the repository root
# with the selected go first in PATH (scripts/with-go).
set -eu
usage() {
	printf 'usage: %s <package> <pattern>... [-- <go test flags>...]\n' "$0" >&2
	exit 2
}
[ "$#" -ge 2 ] || usage
package=$1
shift
# Each argument is one pattern, also when it contains |; keep them one per line.
patterns=
each=
while [ "$#" -gt 0 ] && [ "$1" != -- ]; do
	patterns="${patterns:+$patterns|}$1"
	each="${each:+$each
}$1"
	shift
done
[ -n "$patterns" ] || usage
[ "$#" -eq 0 ] || shift
# -list prints the selected names and then a summary line such as "ok <pkg> <time>";
# only the names count.
listed=$(go test "$@" -list "^($patterns)" "$package")
listed=$(printf '%s\n' "$listed" | grep -E '^[[:alpha:]_][[:alnum:]_]*$' || true)
missing=0
set -f
old_ifs=$IFS
IFS='
'
for pattern in $each; do
	if ! printf '%s\n' "$listed" | grep -Eq "^($pattern)"; then
		printf 'go_test_each: %s selects no test in %s\n' "$pattern" "$package" >&2
		missing=1
	fi
done
IFS=$old_ifs
set +f
[ "$missing" -eq 0 ] || exit 1
exec go test "$@" -run "^($patterns)" "$package"
