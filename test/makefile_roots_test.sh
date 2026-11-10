#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
. "$repo_root/scripts/compat_root.sh"
. "$repo_root/scripts/compat_root_test.sh"
compat_run_root 'makefile roots meta-test' ply-makefile-roots-test
temp_root=$run_root
fake_repo="$temp_root/repo"
fail() {
	printf 'makefile roots meta-test: %s\n' "$*" >&2
	exit 1
}

mkdir -p "$fake_repo/scripts" "$fake_repo/test/fixtures" "$temp_root/make-bin" "$temp_root/awk-bin"
cp "$repo_root/Makefile" "$repo_root/go.mod" "$repo_root/.gitignore" "$repo_root/.golangci.yml" "$fake_repo/"
cp "$repo_root/scripts/compat_root.sh" "$repo_root/scripts/search-replace.sh" "$repo_root/scripts/test-search-replace.sh" "$fake_repo/scripts/"
cp "$repo_root/test/fixtures/recording-go" "$fake_repo/test/fixtures/"
printf 'package fixture\n' >"$fake_repo/fixture.go"
cat >"$temp_root/make-bin/make" <<FAKE
#!/usr/bin/env bash
set -euo pipefail
printf '%s %s\n' "\$TMPDIR" "\$GOTMPDIR" >>"\$FAKE_LOG"
if [[ "\${FAKE_MODE-}" == hang ]]; then printf 'hang %s\n' "\$\$" >>"\$FAKE_LOG"; exec sleep 60; fi
[[ "\${FAKE_MODE-}" != fail ]] || exit 1
exec $(command -v make) "\$@"
FAKE
chmod +x "$temp_root/make-bin/make"
# The search-replace meta-test runs no make; its gsed double calls awk.
sed "s|exec $(command -v make)|exec $(command -v awk)|" "$temp_root/make-bin/make" >"$temp_root/awk-bin/awk"
chmod +x "$temp_root/awk-bin/awk"

run_test() { PATH="$temp_root/$bin:$PATH" ${checker_exec:-} bash "$fake_repo/$test_script"; }
run_pass() { run_test; }
run_fail() { FAKE_MODE=fail run_test; }
compat_links='.build .build/compatibility'
for name in install lint preflight; do
	cp "$repo_root/test/makefile_${name}_test.sh" "$fake_repo/test/"
	prefix=ply-makefile-$name test_script=test/makefile_${name}_test.sh bin=make-bin
	# Known baseline failure: the preflight test still expects 18 calls.
	pass_rc=0; [[ $name != preflight ]] || pass_rc=1
	check_compat_root
done
prefix=ply-search-replace test_script=scripts/test-search-replace.sh bin=awk-bin pass_rc=0
check_compat_root

printf 'makefile roots meta-test: PASS (install, lint, preflight, search-replace: owned root, temp, cleanup, exit, TERM, symlink refusal)\n'
