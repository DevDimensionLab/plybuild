#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

makefile_under_test=${MAKEFILE_UNDER_TEST:-"$repo_root/Makefile"}
fake_go="$repo_root/test/fixtures/recording-go"
calls_file="$tmp_dir/go-calls"
expected_file="$tmp_dir/expected-calls"
make_workspace="$tmp_dir/make workspace"
gobin="$tmp_dir/missing bin"
gopath="$tmp_dir/go path"

mkdir -p "$make_workspace"
touch "$make_workspace/all" "$make_workspace/build" "$make_workspace/install"

run_make() {
	local target=${1-}
	local environment=${2:-explicit}
	local -a make_command=(make --no-print-directory -f "$makefile_under_test")

	if [[ -n "$target" ]]; then
		make_command+=("$target")
	fi

	: >"$calls_file"
	(
		cd "$make_workspace"
		case "$environment" in
			explicit)
				FAKE_GO_CALLS="$calls_file" GOBIN="$gobin" GOPATH="$gopath" GO="$fake_go" \
					"${make_command[@]}"
				;;
			unset)
				env -u GOBIN -u GOPATH FAKE_GO_CALLS="$calls_file" GO="$fake_go" \
					"${make_command[@]}"
				;;
			gopath)
				FAKE_GO_CALLS="$calls_file" GOBIN= GOPATH="$gopath" GO="$fake_go" \
					"${make_command[@]}"
				;;
			*)
				printf 'unknown environment fixture: %s\n' "$environment" >&2
				exit 1
				;;
		esac
	)
}

assert_single_call() {
	local expected_gobin=$1
	local expected_gopath=$2
	shift 2

	{
		printf 'CALL\n'
		printf 'GOBIN=%s\n' "$expected_gobin"
		printf 'GOPATH=%s\n' "$expected_gopath"
		printf 'argc=%s\n' "$#"
		printf 'arg=%s\n' "$@"
	} >"$expected_file"

	diff -u "$expected_file" "$calls_file"
}

run_make install
assert_single_call "$gobin" "$gopath" install ./cmd/ply

if [[ -e "$gobin" ]]; then
	printf 'fake go unexpectedly created GOBIN: %s\n' "$gobin" >&2
	exit 1
fi

run_make install unset
assert_single_call "" "" install ./cmd/ply

run_make install gopath
assert_single_call "" "$gopath" install ./cmd/ply

run_make
assert_single_call "$gobin" "$gopath" build -o ply ./cmd/ply

repo_recording="$repo_root/test/fixtures/.recording-go-leak-probe"
if FAKE_GO_CALLS="$repo_recording" "$fake_go" install ./cmd/ply 2>"$tmp_dir/leak-error"; then
	printf 'recording double accepted a repository-local output path\n' >&2
	exit 1
fi
if [[ -e "$repo_recording" ]]; then
	printf 'recording double leaked into repository: %s\n' "$repo_recording" >&2
	exit 1
fi
grep -F 'refusing to write recording inside repository' "$tmp_dir/leak-error" >/dev/null

printf 'make install contract: PASS (4 invocations and leak guard checked)\n'
