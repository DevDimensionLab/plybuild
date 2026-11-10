#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
. "$repo_root/scripts/compat_root.sh"
compat_run_root 'make install contract' ply-makefile-install
tmp_dir=$run_root

makefile_under_test=${MAKEFILE_UNDER_TEST:-"$repo_root/Makefile"}
fixture_go="$repo_root/test/fixtures/recording-go"
# The run root is inside this repository, which the recording double refuses; record through a copy.
fake_go="$tmp_dir/recording/test/fixtures/recording-go"
mkdir -p "$(dirname "$fake_go")"
cp "$fixture_go" "$fake_go"
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
	local output=${3-}
	local -a make_command=(make --no-print-directory -f "$makefile_under_test")

	if [[ -n "$target" ]]; then
		make_command+=("$target")
	fi

	if [[ -n "$output" ]]; then
		make_command+=("PLY_BUILD_OUTPUT=$output")
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

# Identity must come from the source checkout even when make is invoked from
# elsewhere. Explicit relative output paths still belong to the caller's working directory.
source_commit=
source_dirty=
if [[ "$(git -C "$repo_root" rev-parse --show-toplevel 2>/dev/null || true)" == "$repo_root" ]]; then
	source_commit=$(git -C "$repo_root" rev-parse --verify HEAD)
	source_dirty=false
	[[ -z "$(git -C "$repo_root" status --porcelain)" ]] || source_dirty=true
fi
identity_flags="-X github.com/devdimensionlab/plybuild/internal/buildinfo.SourceCommit=$source_commit -X github.com/devdimensionlab/plybuild/internal/buildinfo.SourceDirty=$source_dirty"

run_make build explicit 'relative output/ply'
assert_single_call "$gobin" "$gopath" build -ldflags "$identity_flags" -o "$make_workspace/relative output/ply" ./cmd/ply

run_make build explicit "$tmp_dir/absolute output/ply"
assert_single_call "$gobin" "$gopath" build -ldflags "$identity_flags" -o "$tmp_dir/absolute output/ply" ./cmd/ply

fail() {
	printf 'make install contract: %s\n' "$*" >&2
	exit 1
}

repo_makefile=$makefile_under_test
build_fixture="$tmp_dir/build-fixture"
mkdir -p "$build_fixture"
cp "$repo_makefile" "$build_fixture/Makefile"
makefile_under_test="$build_fixture/Makefile"
unbound_identity_flags="-X github.com/devdimensionlab/plybuild/internal/buildinfo.SourceCommit= -X github.com/devdimensionlab/plybuild/internal/buildinfo.SourceDirty="

run_make
assert_single_call "$gobin" "$gopath" build -ldflags "$unbound_identity_flags" -o "$build_fixture/.build/bin/ply" ./cmd/ply
[[ -d "$build_fixture/.build/bin" && ! -L "$build_fixture/.build" && ! -L "$build_fixture/.build/bin" ]] ||
	fail 'default build did not create the physical repository-local .build/bin directory'
[[ ! -e "$make_workspace/ply" && ! -e "$make_workspace/.build" ]] ||
	fail 'default build wrote output into the caller working directory'

: >"$calls_file"
(cd "$build_fixture" && FAKE_GO_CALLS="$calls_file" GOBIN="$gobin" GOPATH="$gopath" GO="$fake_go" \
	make --no-print-directory build)
assert_single_call "$gobin" "$gopath" build -ldflags "$unbound_identity_flags" -o "$(cd "$build_fixture" && pwd -P)/.build/bin/ply" ./cmd/ply

rm -rf -- "$build_fixture/.build"
run_make build explicit 'relative output/ply'
assert_single_call "$gobin" "$gopath" build -ldflags "$unbound_identity_flags" -o "$make_workspace/relative output/ply" ./cmd/ply
[[ ! -e "$build_fixture/.build" ]] || fail 'explicit build output created the default .build directory'

symlink_target="$tmp_dir/external build target"
mkdir -p "$symlink_target"
ln -s "$symlink_target" "$build_fixture/.build"
if run_make build explicit 2>"$tmp_dir/symlink-build-error"; then
	fail 'default build accepted a symlinked .build root'
fi
[[ ! -s "$calls_file" ]] || fail 'default build invoked go through a symlinked .build root'
[[ -z "$(find "$symlink_target" -mindepth 1 -print)" ]] || fail 'default build wrote through a symlinked .build root'
rm -- "$build_fixture/.build"

mkdir -p "$build_fixture/.build"
ln -s "$symlink_target" "$build_fixture/.build/bin"
printf 'keep\n' >"$symlink_target/keep"
if run_make build explicit 2>"$tmp_dir/symlink-bin-build-error"; then
	fail 'default build accepted a symlinked .build/bin'
fi
[[ ! -s "$calls_file" ]] || fail 'default build invoked go through a symlinked .build/bin'
[[ "$(find "$symlink_target" -mindepth 1 -print)" == "$symlink_target/keep" && "$(cat "$symlink_target/keep")" == keep ]] ||
	fail 'default build changed the external target of a symlinked .build/bin'
rm -- "$build_fixture/.build/bin"
rmdir -- "$build_fixture/.build"
rm -- "$symlink_target/keep"
makefile_under_test=$repo_makefile

# Without an explicit GO seam, ordinary targets select Go through scripts/with-go; stale PATH go never runs.
selector_calls="$tmp_dir/selector-calls"
stale_bin="$tmp_dir/selector stale/bin"
selected_go="$tmp_dir/selector sdk/bin/go"
selected_version=$(awk '$1 == "toolchain" { print $2 }' "$repo_root/go.mod")
mkdir -p "$stale_bin" "$(dirname "$selected_go")"
printf '#!/usr/bin/env bash\nprintf "stale %%s\\n" "$*" >>"$FAKE_SELECTOR_CALLS"\nexit 97\n' >"$stale_bin/go"
cat >"$selected_go" <<EOF
#!/usr/bin/env bash
set -euo pipefail
if [[ "\$*" == 'env GOVERSION GOROOT' ]]; then
	printf '%s\n%s\n' '$selected_version' '$tmp_dir/selector sdk'
	exit 0
fi
printf 'selected GOTOOLCHAIN=%s %s\n' "\$GOTOOLCHAIN" "\$*" >>"\$FAKE_SELECTOR_CALLS"
EOF
chmod +x "$stale_bin/go" "$selected_go"

run_default() {
	local target=$1
	shift
	: >"$selector_calls"
	(cd "$make_workspace" && env -u GO -u PLY_GO PATH="$stale_bin:$PATH" GOTOOLCHAIN=auto \
		FAKE_SELECTOR_CALLS="$selector_calls" "$@" \
		make --no-print-directory -f "$repo_makefile" "PLY_BUILD_OUTPUT=$tmp_dir/selector output/ply" "$target")
}

if run_default build 2>"$tmp_dir/selector-unset-error"; then
	fail 'default build ran without an explicit PLY_GO selection'
fi
grep -Fq 'PLY_GO must be an absolute path to the selected go binary' "$tmp_dir/selector-unset-error" ||
	fail 'default build did not require the explicit PLY_GO selection'
[[ ! -s "$selector_calls" ]] || fail 'default build invoked a go without an explicit selection'

run_default build PLY_GO="$selected_go" 2>"$tmp_dir/selector-build-error" ||
	{ sed -n '1,40p' "$tmp_dir/selector-build-error" >&2; fail 'default build failed with explicit PLY_GO'; }
[[ $(cat "$selector_calls") == "selected GOTOOLCHAIN=local build -ldflags $identity_flags -o $tmp_dir/selector output/ply ./cmd/ply" ]] ||
	{ sed -n '1,20p' "$selector_calls" >&2; fail 'default build did not run only the selected go'; }

run_default install PLY_GO="$selected_go" 2>"$tmp_dir/selector-install-error" ||
	{ sed -n '1,40p' "$tmp_dir/selector-install-error" >&2; fail 'default install failed with explicit PLY_GO'; }
[[ $(cat "$selector_calls") == 'selected GOTOOLCHAIN=local install ./cmd/ply' ]] ||
	{ sed -n '1,20p' "$selector_calls" >&2; fail 'default install did not run only the selected go'; }

new_clean_fixture() {
	local root="$tmp_dir/clean-fixture-$1"
	mkdir -p "$root"
	cp "$repo_makefile" "$root/Makefile"
	cp "$repo_root/go.mod" "$root/go.mod"
	cp "$repo_root/.gitignore" "$root/.gitignore"
	printf 'source\n' >"$root/source.txt"
	git -C "$root" init -q
	git -C "$root" add .
	git -C "$root" -c user.name='Clean Fixture' -c user.email=clean@example.invalid \
		-c commit.gpgsign=false commit -q -m fixture
	printf '%s\n' "$root"
}

run_clean() {
	local root=$1
	shift
	(cd "$make_workspace" && make --no-print-directory -f "$root/Makefile" "$@" clean)
}

clean_root=$(new_clean_fixture clean)
run_clean "$clean_root" || fail 'clean failed without a .build directory'
[[ ! -e "$clean_root/.build" && ! -L "$clean_root/.build" ]] || fail 'clean created .build'

external_target="$tmp_dir/external symlink target"
external_output="$tmp_dir/explicit external output/ply"
mkdir -p "$external_target" "$(dirname "$external_output")" "$make_workspace/.build"
printf 'keep\n' >"$external_target/keep"
printf 'keep\n' >"$external_output"
printf 'keep\n' >"$make_workspace/.build/keep"
printf 'keep\n' >"$clean_root/.build-neighbor"
mkdir -p "$clean_root/.build/bin" "$clean_root/.build/dist/nested/read-only"
printf 'binary\n' >"$clean_root/.build/bin/ply"
printf 'artifact\n' >"$clean_root/.build/dist/nested/read-only/ply"
ln -s "$external_target" "$clean_root/.build/dist/outside-link"
chmod a-w "$clean_root/.build/dist/nested/read-only" "$clean_root/.build/dist/nested"
run_clean "$clean_root" "PLY_BUILD_OUTPUT=$external_output" || fail 'clean failed for an owned .build tree'
[[ ! -e "$clean_root/.build" && ! -L "$clean_root/.build" ]] || fail 'clean left .build behind'
[[ -f "$external_target/keep" ]] || fail 'clean followed a symlink out of .build'
[[ -f "$external_output" ]] || fail 'clean removed explicit external build output'
[[ -f "$make_workspace/.build/keep" ]] || fail 'clean removed the caller working directory .build'
[[ -f "$clean_root/.build-neighbor" && -f "$clean_root/source.txt" && -f "$clean_root/go.mod" ]] ||
	fail 'clean removed neighbor or source files'
run_clean "$clean_root" || fail 'repeated clean was not idempotent'

mkdir -p "$clean_root/.build/dist/partial"
chmod a-w "$clean_root/.build/dist/partial"
(cd "$clean_root" && make --no-print-directory clean) || fail 'clean could not finish a partially removed .build'
[[ ! -e "$clean_root/.build" ]] || fail 'clean left a partially removed .build behind'

for nested_mode in 000 600; do
	mkdir -p "$clean_root/.build/a/b"
	printf 'generated\n' >"$clean_root/.build/a/b/file"
	chmod "$nested_mode" "$clean_root/.build/a/b"
	chmod "$nested_mode" "$clean_root/.build/a"
	run_clean "$clean_root" >"$tmp_dir/nested-mode-clean" 2>&1 ||
		fail "one clean did not remove nested mode $nested_mode directories: $(cat "$tmp_dir/nested-mode-clean")"
	[[ ! -e "$clean_root/.build" ]] || fail "one clean left nested mode $nested_mode directories behind"
done

expect_clean_refusal() {
	local label=$1
	local root=$2
	local message=$3
	if run_clean "$root" >"$tmp_dir/clean-refusal" 2>&1; then
		fail "clean accepted $label"
	fi
	grep -F "$message" "$tmp_dir/clean-refusal" >/dev/null ||
		fail "clean refusal for $label did not explain itself: $(cat "$tmp_dir/clean-refusal")"
}

symlink_root=$(new_clean_fixture symlink)
ln -s "$external_target" "$symlink_root/.build"
expect_clean_refusal 'a symlinked .build root' "$symlink_root" 'clean: refusing symlinked'
[[ -L "$symlink_root/.build" && -f "$external_target/keep" ]] || fail 'clean changed a symlinked .build root'

file_root=$(new_clean_fixture file)
printf 'not a directory\n' >"$file_root/.build"
expect_clean_refusal 'a non-directory .build' "$file_root" 'clean: refusing non-directory'
[[ -f "$file_root/.build" ]] || fail 'clean removed a non-directory .build'

tracked_root=$(new_clean_fixture tracked)
mkdir -p "$tracked_root/.build/bin"
printf 'tracked\n' >"$tracked_root/.build/tracked.txt"
printf 'generated\n' >"$tracked_root/.build/bin/ply"
git -C "$tracked_root" add -f .build/tracked.txt
git -C "$tracked_root" -c user.name='Clean Fixture' -c user.email=clean@example.invalid \
	-c commit.gpgsign=false commit -q -m tracked
expect_clean_refusal 'tracked files under .build' "$tracked_root" 'clean: refusing tracked files'
[[ -f "$tracked_root/.build/tracked.txt" && -f "$tracked_root/.build/bin/ply" ]] ||
	fail 'clean removed content beside tracked files under .build'

nested_root=$(new_clean_fixture nested)
nested_project="$nested_root/nested-project"
mkdir -p "$nested_project/.build/bin"
cp "$repo_makefile" "$nested_project/Makefile"
cp "$repo_root/go.mod" "$nested_project/go.mod"
printf 'generated\n' >"$nested_project/.build/bin/ply"
expect_clean_refusal 'a Makefile below the Git root' "$nested_project" 'clean: not the Ply repository root'
[[ -f "$nested_project/.build/bin/ply" ]] || fail 'clean removed .build below a foreign Git root'

foreign_root=$(new_clean_fixture foreign)
printf 'module example.invalid/foreign\n' >"$foreign_root/go.mod"
git -C "$foreign_root" -c user.name='Clean Fixture' -c user.email=clean@example.invalid \
	-c commit.gpgsign=false commit -q -am foreign
mkdir -p "$foreign_root/.build/bin"
printf 'generated\n' >"$foreign_root/.build/bin/ply"
expect_clean_refusal 'a foreign module identity' "$foreign_root" 'clean: not the Ply repository root'
[[ -f "$foreign_root/.build/bin/ply" ]] || fail 'clean removed .build in a foreign repository'

repo_recording="$repo_root/test/fixtures/.recording-go-leak-probe"
if FAKE_GO_CALLS="$repo_recording" "$fixture_go" install ./cmd/ply 2>"$tmp_dir/leak-error"; then
	printf 'recording double accepted a repository-local output path\n' >&2
	exit 1
fi
if [[ -e "$repo_recording" ]]; then
	printf 'recording double leaked into repository: %s\n' "$repo_recording" >&2
	exit 1
fi
grep -F 'refusing to write recording inside repository' "$tmp_dir/leak-error" >/dev/null

printf 'make install contract: PASS (install/build invocations, build identity, caller and .build output paths, clean refusals/idempotence and leak guard checked)\n'
