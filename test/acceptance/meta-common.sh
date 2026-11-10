#!/usr/bin/env bash

acceptance_meta_control=0

acceptance_meta_run() {
	local expected=$1
	local label=$2
	shift 2
	acceptance_meta_control=$((acceptance_meta_control + 1))
	local output="$work_root/control-${acceptance_meta_control}.log"
	local actual
	set +e
	"$@" >"$output" 2>&1
	actual=$?
	set -e
	if [[ "$expected" == pass && "$actual" -ne 0 ]] ||
		[[ "$expected" == fail && "$actual" -eq 0 ]]; then
		printf 'meta-control %s: FAIL (expected %s, exit %s)\n' \
			"$label" "$expected" "$actual" >&2
		sed -n '1,160p' "$output" >&2
		return 1
	fi
	printf 'meta-control %s: PASS\n' "$label"
}

acceptance_meta_expect_pass() {
	acceptance_meta_run pass "$@"
}

acceptance_meta_expect_fail() {
	acceptance_meta_run fail "$@"
}
