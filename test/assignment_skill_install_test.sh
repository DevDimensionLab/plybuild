#!/usr/bin/env bash
# Usage: test/assignment_skill_install_test.sh <empty-test-root>
# Exercises scripts/install-assignment-skill.sh and make install-assignment-skill.

set -euo pipefail
export LC_ALL=C LANG=C

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
installer="$repo_root/scripts/install-assignment-skill.sh"
source_file="$repo_root/skills/ply-assignment/SKILL.md"
root=${1-}

fail() {
	printf 'assignment skill install test: %s\n' "$*" >&2
	exit 1
}

[[ -n "$root" && -d "$root" && ! -L "$root" ]] || fail 'test root must be an existing directory'
[[ -z "$(ls -A "$root")" ]] || fail 'test root must be empty'
[[ -f "$source_file" ]] || fail 'skill source is missing'

# run <expected-exit> <stdout-file> <stderr-file> <installer arguments...>
run() {
	local expected=$1 out=$2 err=$3 status=0
	shift 3
	/bin/sh "$installer" "$@" >"$out" 2>"$err" || status=$?
	[[ "$status" -eq "$expected" ]] || fail "exit $status, expected $expected for: $*"
}

tree() {
	(cd "$1" && find . -print0 | sort -z | xargs -0 ls -ld | awk '{print $1, $NF}' && find . -type f -exec cksum {} + | sort)
}

# 14. First installation, through make, then an identical rerun is a no-op.
dest="$root/first/app/.claude/skills/ply-assignment"
make --no-print-directory -C "$repo_root" install-assignment-skill ASSIGNMENT_SKILL_DEST="$dest" >"$root/make.out" 2>&1 ||
	fail "make install-assignment-skill failed: $(cat "$root/make.out")"
cmp -s "$source_file" "$dest/SKILL.md" || fail 'installed bytes differ from source'
[[ "$(stat -f '%Lp' "$dest/SKILL.md" 2>/dev/null || stat -c '%a' "$dest/SKILL.md")" == 644 ]] || fail 'installed mode is not 0644'
before=$(tree "$root/first")
make --no-print-directory -C "$repo_root" install-assignment-skill ASSIGNMENT_SKILL_DEST="$dest" >"$root/make.out" 2>&1 ||
	fail "identical rerun failed: $(cat "$root/make.out")"
[[ "$(tree "$root/first")" == "$before" ]] || fail 'identical rerun changed the destination'

# 15. Conflicts have no effect; a foreign file in the destination is kept.
conflict() {
	local name=$1 setup=$2 case_root="$root/$1"
	local target="$case_root/skill"
	mkdir -p "$case_root"
	eval "$setup"
	local snapshot
	snapshot=$(tree "$case_root")
	run 1 "$case_root.out" "$case_root.err" "$source_file" "$target"
	grep -q '^existing_install_conflict: ' "$case_root.err" || fail "$name: no existing_install_conflict on stderr"
	[[ "$(tree "$case_root")" == "$snapshot" ]] || fail "$name: conflict changed files"
}
conflict different 'mkdir "$target"; printf foreign >"$target/notes.txt"; printf other >"$target/SKILL.md"'
conflict skill-symlink 'mkdir "$target"; printf foreign >"$target/notes.txt"; cp "$source_file" "$case_root/real"; ln -s "$case_root/real" "$target/SKILL.md"'
conflict dest-symlink 'mkdir "$case_root/real"; printf foreign >"$case_root/real/notes.txt"; ln -s "$case_root/real" "$target"'
conflict dest-symlink-slash 'mkdir "$case_root/real"; printf foreign >"$case_root/real/notes.txt"; ln -s "$case_root/real" "$target"; target="$target/"'
conflict dest-file 'printf foreign >"$target"'
conflict skill-directory 'mkdir -p "$target/SKILL.md"; printf foreign >"$target/notes.txt"'

# A SKILL.md directory or symlink that appears between the checks and the
# link call is a conflict. The injected fake ln/link records what lies under
# the foreign target right after the real call, so even a write that is later
# removed fails the test; afterwards the foreign tree must be unchanged.
race() {
	local name=$1 case_root="$root/$1"
	local target="$case_root/skill" bin="$case_root/bin" status=0
	mkdir -p "$bin" "$case_root/foreign"
	printf foreign >"$case_root/foreign/notes.txt"
	printf '%s\n' '#!/bin/sh' 'for last; do :; done' \
		'case $RACE in' \
		'dir) mkdir "$last" && printf foreign >"$last/notes.txt"; watch=$last ;;' \
		'symlink) /bin/ln -s "$RACE_TARGET" "$last"; watch=$RACE_TARGET ;;' \
		'esac' \
		'status=0; "/bin/${0##*/}" "$@" || status=$?' \
		'ls -A "$watch" >"$RACE_LOG"' \
		'exit $status' >"$bin/ln"
	cp "$bin/ln" "$bin/link"
	chmod 755 "$bin/ln" "$bin/link"
	PATH="$bin:$PATH" RACE=$name RACE_TARGET="$case_root/foreign" RACE_LOG="$case_root.watch" \
		/bin/sh "$installer" "$source_file" "$target" >"$case_root.out" 2>"$case_root.err" || status=$?
	[[ -f "$case_root.watch" ]] || fail "race $name: the link call was not reached"
	[[ "$(cat "$case_root.watch")" == notes.txt ]] || fail "race $name: the link call wrote under the foreign target: $(cat "$case_root.watch")"
	[[ "$status" -eq 1 ]] || fail "race $name: exit $status, expected 1"
	grep -q '^existing_install_conflict: ' "$case_root.err" || fail "race $name: no existing_install_conflict on stderr"
	! grep -q '^installed' "$case_root.out" || fail "race $name: reported installed"
	[[ "$(ls -A "$target")" == SKILL.md ]] || fail "race $name: destination holds $(ls -A "$target")"
	local foreign="$case_root/foreign"
	[[ "$name" == symlink ]] || foreign="$target/SKILL.md"
	[[ "$(ls -A "$foreign")" == notes.txt && "$(cat "$foreign/notes.txt")" == foreign ]] ||
		fail "race $name: foreign tree changed"
}
race dir
race symlink

# A failing utility leaves no file behind. A link failure while SKILL.md is
# still missing is an installation I/O error, not a conflict.
failing() {
	local name=$1 tool=$2 case_root="$root/$1"
	local target="$case_root/skill" bin="$case_root/bin" status=0
	mkdir -p "$bin"
	printf '%s\n' '#!/bin/sh' 'echo "injected failure" >&2' 'exit 1' >"$bin/$tool"
	chmod 755 "$bin/$tool"
	PATH="$bin:$PATH" /bin/sh "$installer" "$source_file" "$target" >"$case_root.out" 2>"$case_root.err" || status=$?
	[[ "$status" -eq 1 ]] || fail "$name: exit $status, expected 1"
	grep -q '^install_io_error: ' "$case_root.err" || fail "$name: no install_io_error on stderr: $(cat "$case_root.err")"
	! grep -q 'existing_install_conflict' "$case_root.err" || fail "$name: reported a conflict"
	[[ -z "$(ls -A "$target")" ]] || fail "$name: left $(ls -A "$target")"
}
failing link-fails link
failing cat-fails cat

# A foreign file beside an identical SKILL.md is kept by a no-op rerun.
printf foreign >"$dest/notes.txt"
run 0 "$root/keep.out" "$root/keep.err" "$source_file" "$dest"
[[ "$(cat "$dest/notes.txt")" == foreign ]] || fail 'foreign file changed'

# 16. An empty or missing destination fails without effect, also through make.
mkdir "$root/empty"
snapshot=$(cd "$root/empty" && ls -A)
(cd "$root/empty" && run 2 "$root/empty.out" "$root/empty.err" "$source_file" "")
(cd "$root/empty" && run 2 "$root/missing.out" "$root/missing.err" "$source_file")
status=0
(cd "$root/empty" && make --no-print-directory -C "$repo_root" install-assignment-skill >"$root/make-empty.out" 2>&1) || status=$?
[[ "$status" -ne 0 ]] || fail 'make without ASSIGNMENT_SKILL_DEST succeeded'
[[ "$(cd "$root/empty" && ls -A)" == "$snapshot" ]] || fail 'empty destination had an effect'
[[ ! -e "$repo_root/SKILL.md" ]] || fail 'empty destination wrote into the repository'

printf '%s\n' 'assignment skill install test: PASS'
