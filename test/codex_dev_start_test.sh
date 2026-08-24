#!/usr/bin/env bash

set -euo pipefail
export LC_ALL=C LANG=C
unset CDPATH \
	GIT_ALTERNATE_OBJECT_DIRECTORIES \
	GIT_COMMON_DIR \
	GIT_CONFIG \
	GIT_CONFIG_PARAMETERS \
	GIT_DIR \
	GIT_GRAFT_FILE \
	GIT_IMPLICIT_WORK_TREE \
	GIT_INDEX_FILE \
	GIT_NO_REPLACE_OBJECTS \
	GIT_OBJECT_DIRECTORY \
	GIT_PREFIX \
	GIT_REPLACE_REF_BASE \
	GIT_SHALLOW_FILE \
	GIT_WORK_TREE
export GIT_CONFIG_NOSYSTEM=1
export GIT_CONFIG_SYSTEM=/dev/null
export GIT_CONFIG_GLOBAL=/dev/null
export GIT_CONFIG_COUNT=0

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
launcher="$repo_root/codex-dev-start.sh"
temp_root=$(mktemp -d "${TMPDIR:-/tmp}/codex-dev-start-test.XXXXXX")
trap 'rm -rf "$temp_root"' EXIT

expected_skeleton_sha256='3729ed7b311218543d1139ab69d52ca90ff1f3d332d55b09d273dd881f453cf6'
pass_count=0
nested_mode='no'
environment_probe='no'

fail() {
	printf 'codex dev start contract: FAIL: %s\n' "$*" >&2
	exit 1
}

pass() {
	pass_count=$((pass_count + 1))
	printf 'ok %s - %s\n' "$pass_count" "$1"
}

case $# in
0) ;;
1)
	case $1 in
	--internal-source-archive-child) nested_mode='yes' ;;
	--internal-environment-probe) environment_probe='yes' ;;
	*) fail "unknown argument: $1" ;;
	esac
	;;
*) fail 'expected at most one argument' ;;
esac

[[ -z "${CODEX_DEV_START_NESTED+x}" ]] ||
	fail 'CODEX_DEV_START_NESTED must be unset; nested mode is an internal argument'

if [[ "$environment_probe" == 'yes' ]]; then
	[[ -z "${CDPATH+x}" ]] || fail 'CDPATH was not sanitized'
	[[ "$GIT_CONFIG_NOSYSTEM" == '1' ]] || fail 'Git system config was not disabled'
	[[ "$GIT_CONFIG_SYSTEM" == '/dev/null' ]] || fail 'Git system config path was not sanitized'
	[[ "$GIT_CONFIG_GLOBAL" == '/dev/null' ]] || fail 'Git global config path was not sanitized'
	[[ "$GIT_CONFIG_COUNT" == '0' ]] || fail 'Git command config was not sanitized'
	[[ -z "${GIT_CONFIG_PARAMETERS+x}" ]] || fail 'Git config parameters were not sanitized'
	[[ -z "${GIT_DIR+x}" ]] || fail 'Git repository directory was not sanitized'
	[[ -z "${GIT_WORK_TREE+x}" ]] || fail 'Git worktree path was not sanitized'
	[[ -z "${GIT_INDEX_FILE+x}" ]] || fail 'Git index path was not sanitized'
	if git config --global --get core.fsmonitor >/dev/null 2>&1; then
		fail 'sanitized Git environment still reads a global fsmonitor'
	fi
	printf '%s\n' 'codex dev start environment: PASS'
	exit 0
fi

assert_contains() {
	local file=$1
	local expected=$2
	if ! grep -F -- "$expected" "$file" >/dev/null; then
		printf '%s\n' '--- actual output ---' >&2
		sed -n '1,80p' "$file" >&2
		fail "expected '$expected' in $file"
	fi
}

assert_not_contains() {
	local file=$1
	local unexpected=$2
	if grep -F -- "$unexpected" "$file" >/dev/null; then
		fail "unexpected '$unexpected' in $file"
	fi
}

sha256_file() {
	local file=$1
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$file" | awk '{ print $1 }'
	elif command -v shasum >/dev/null 2>&1; then
		shasum -a 256 "$file" | awk '{ print $1 }'
	elif command -v openssl >/dev/null 2>&1; then
		openssl dgst -sha256 "$file" | awk '{ print $NF }'
	else
		fail 'no SHA-256 implementation is available'
	fi
}

extract_header_value() {
	local source_file=$1
	local key=$2
	awk -v prefix="#|$key=" '
		$0 == "# CODEX_MUTABLE_SESSION_HEADER_BEGIN" { inside = 1; next }
		$0 == "# CODEX_MUTABLE_SESSION_HEADER_END" { inside = 0; next }
		inside && index($0, prefix) == 1 {
			count++
			value = substr($0, length(prefix) + 1)
		}
		END {
			if (count != 1) exit 2
			print value
		}
	' "$source_file"
}

extract_launcher_prompt() {
	local source_file=$1
	awk '
		$0 == "# CODEX_MUTABLE_PROMPT_BEGIN" { region = 1; next }
		$0 == "# CODEX_MUTABLE_PROMPT_END" { exit }
		region {
			if (substr($0, 1, 2) != "#|") exit 2
			print substr($0, 3)
		}
	' "$source_file"
}

extract_archive_prompt() {
	local archive=$1
	awk '
		$0 == "<!-- CODEX_SESSION_PROMPT_BEGIN -->" { inside = 1; next }
		$0 == "<!-- CODEX_SESSION_PROMPT_END -->" { exit }
		inside { print }
	' "$archive"
}

validate_mutable_source_structure() {
	local source_file=$1
	awk '
		$0 == "# CODEX_STABLE_EXECUTION_END" {
			stable_count++
			stable_line = NR
			if (previous != "exit 70") exit 2
		}
		$0 == "# CODEX_MUTABLE_SESSION_HEADER_BEGIN" {
			if (region) exit 2
			region = "header"
			header_begin_count++
			header_begin_line = NR
			next
		}
		$0 == "# CODEX_MUTABLE_SESSION_HEADER_END" {
			if (region != "header" || header_lines != 4) exit 2
			region = ""
			header_end_count++
			next
		}
		$0 == "# CODEX_MUTABLE_PROMPT_BEGIN" {
			if (region) exit 2
			region = "prompt"
			prompt_begin_count++
			prompt_begin_line = NR
			next
		}
		$0 == "# CODEX_MUTABLE_PROMPT_END" {
			if (region != "prompt" || prompt_lines < 1) exit 2
			region = ""
			prompt_end_count++
			next
		}
		region == "header" {
			header_lines++
			if (header_lines == 1 && $0 !~ /^#[|]SESSION_STATUS=(NEXT|COMPLETE)$/) exit 2
			if (header_lines == 2 && $0 !~ /^#[|]SESSION_ID=[[:print:]]+$/) exit 2
			if (header_lines == 3 && $0 !~ /^#[|]SESSION_ARCHIVE_REL=[[:print:]]+$/) exit 2
			if (header_lines == 4 && $0 !~ /^#[|]PREVIOUS_SESSION_ARCHIVE_REL=[[:print:]]*$/) exit 2
			if (header_lines > 4) exit 2
			next
		}
		region == "prompt" {
			prompt_lines++
			if (substr($0, 1, 2) != "#|") exit 2
			next
		}
		{ previous = $0 }
		END {
			if (region || stable_count != 1 || header_begin_count != 1 ||
			    header_end_count != 1 || prompt_begin_count != 1 ||
			    prompt_end_count != 1 || header_begin_line <= stable_line ||
			    prompt_begin_line <= header_begin_line) exit 2
		}
	' "$source_file"
}

normalize_launcher_skeleton() {
	local source_file=$1
	awk '
		$0 == "# CODEX_MUTABLE_SESSION_HEADER_BEGIN" {
			if (inside) exit 2
			inside = "header"
			header_begin++
			print
			print "__CODEX_MUTABLE_SESSION_HEADER__"
			next
		}
		$0 == "# CODEX_MUTABLE_SESSION_HEADER_END" {
			if (inside != "header") exit 2
			inside = ""
			header_end++
			print
			next
		}
		$0 == "# CODEX_MUTABLE_PROMPT_BEGIN" {
			if (inside) exit 2
			inside = "prompt"
			prompt_begin++
			print
			print "__CODEX_MUTABLE_PROMPT__"
			next
		}
		$0 == "# CODEX_MUTABLE_PROMPT_END" {
			if (inside != "prompt") exit 2
			inside = ""
			prompt_end++
			print
			next
		}
		!inside { print }
		END {
			if (inside || header_begin != 1 || header_end != 1 ||
			    prompt_begin != 1 || prompt_end != 1) exit 2
		}
	' "$source_file"
}

replace_header_value() {
	local source_file=$1
	local key=$2
	local value=$3
	local replacement="$temp_root/header-replacement"
	awk -v prefix="#|$key=" -v replacement="#|$key=$value" '
		index($0, prefix) == 1 { print replacement; replaced++; next }
		{ print }
		END { if (replaced != 1) exit 2 }
	' "$source_file" >"$replacement" || fail "could not replace $key"
	mv "$replacement" "$source_file"
	chmod +x "$source_file"
}

insert_second_generation_prompt() {
	local source_file=$1
	local replacement="$temp_root/prompt-replacement"
	awk '
		in_prompt && !inserted && $0 == "#|# Mission" {
			print
			print "#|"
			print "#|Resume from the contract-tested second-generation handoff."
			inserted = 1
			next
		}
		$0 == "# CODEX_MUTABLE_PROMPT_BEGIN" { in_prompt = 1 }
		{ print }
		END { if (!inserted) exit 2 }
	' "$source_file" >"$replacement" || fail 'could not mutate second-generation prompt'
	mv "$replacement" "$source_file"
	chmod +x "$source_file"
}

archive_leaf() {
	basename "$1"
}

enumerate_source_archives() {
	local root=$1
	local dotglob_was_set='no'
	local nullglob_was_set='no'

	if shopt -q dotglob; then
		dotglob_was_set='yes'
	fi
	if shopt -q nullglob; then
		nullglob_was_set='yes'
	fi
	shopt -s dotglob nullglob
	SOURCE_ARCHIVES=("$root"/docs/plan/agent-sessions/*.md)
	if [[ "$dotglob_was_set" == 'no' ]]; then
		shopt -u dotglob
	fi
	if [[ "$nullglob_was_set" == 'no' ]]; then
		shopt -u nullglob
	fi
}

source_graph_is_regular() {
	local root=$1
	local path
	local archive

	for path in 'docs' 'docs/design' 'docs/plan' 'docs/plan/agent-sessions'; do
		[[ -d "$root/$path" && ! -L "$root/$path" ]] || return 1
	done
	for path in \
		'codex-dev-start.sh' \
		'docs/design/agent-session-continuity.md' \
		'docs/design/quality-lift.md' \
		'docs/plan/quality-upgrade.md' \
		'docs/plan/quality-handover.md'; do
		[[ -f "$root/$path" && ! -L "$root/$path" ]] || return 1
	done
	enumerate_source_archives "$root"
	[[ ${#SOURCE_ARCHIVES[@]} -gt 0 ]] || return 1
	for archive in "${SOURCE_ARCHIVES[@]}"; do
		[[ -f "$archive" && ! -L "$archive" ]] || return 1
		[[ "${archive##*/}" != .* ]] || return 1
	done
}

write_archive() {
	local fixture=$1
	local archive_rel=$2
	local status=$3
	local previous=$4
	local next=$5
	local outcome=$6
	local prompt_source=$7
	local archive="$fixture/$archive_rel"
	local session_id
	local created
	local prompt_file
	local prompt_digest
	session_id=$(archive_leaf "$archive_rel")
	session_id=${session_id%.md}
	created="${session_id:0:13}:${session_id:13:2}:${session_id:15:2}${session_id:17:3}:${session_id:20:2}"
	prompt_file=$(mktemp "$temp_root/archive-prompt.XXXXXX")
	extract_launcher_prompt "$prompt_source" >"$prompt_file"
	prompt_digest=$(sha256_file "$prompt_file")
	mkdir -p "$(dirname "$archive")"
	{
		printf '# Agent Session: Contract Fixture\n\n'
		printf 'Status: %s\n' "$status"
		printf 'Session ID: `%s`\n' "$session_id"
		printf 'Created: `%s`\n' "$created"
		printf 'Source: `codex-dev-start.sh`\n'
		printf 'Prompt SHA-256: `%s`\n' "$prompt_digest"
		if [[ "$previous" == 'none' ]]; then
			printf 'Previous: none\n'
		else
			printf 'Previous: [%s](%s)\n' "$previous" "$previous"
		fi
		if [[ "$next" == 'none' ]]; then
			printf 'Next: none\n'
		else
			printf 'Next: [%s](%s)\n' "$next" "$next"
		fi
		printf 'Outcome: %s\n\n' "$outcome"
		printf 'The block below is the byte-exact Codex prompt argument, including its terminal LF.\n\n'
		printf '<!-- CODEX_SESSION_PROMPT_BEGIN -->\n'
		cat "$prompt_file"
		printf '<!-- CODEX_SESSION_PROMPT_END -->\n'
	} >"$archive"
	rm -f "$prompt_file"
}

refresh_archive_digest() {
	local archive=$1
	local prompt_file
	local replacement
	local prompt_digest
	prompt_file=$(mktemp "$temp_root/archive-prompt-refresh.XXXXXX")
	replacement=$(mktemp "$temp_root/archive-refresh.XXXXXX")
	extract_archive_prompt "$archive" >"$prompt_file"
	prompt_digest=$(sha256_file "$prompt_file")
	awk -v digest="Prompt SHA-256: \`$prompt_digest\`" '
		/^Prompt SHA-256: / { print digest; replaced++; next }
		{ print }
		END { if (replaced != 1) exit 2 }
	' "$archive" >"$replacement" || fail 'could not refresh archive prompt digest'
	mv "$replacement" "$archive"
	rm -f "$prompt_file"
}

expect_failure() {
	local name=$1
	local expected=$2
	shift 2
	if "$@" >"$temp_root/failure.stdout" 2>"$temp_root/failure.stderr"; then
		fail "$name unexpectedly succeeded"
	fi
	assert_contains "$temp_root/failure.stderr" "$expected"
	pass "$name fails closed"
}

git_at() {
	local root=$1
	shift
	env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null \
		git -C "$root" "$@"
}

git_fixture() {
	git_at "$fixture_root" "$@"
}

commit_fixture() {
	git_fixture add .
	git_fixture -c commit.gpgSign=false -c core.hooksPath="$temp_root/empty-hooks" \
		commit -qm "$1"
}

complete_authorized_queue() {
	local plan=$1
	local replacement="$temp_root/completed-authorized-queue"
	awk '
		$0 == "<!-- CODEX_AUTHORIZED_CHECKPOINTS_BEGIN -->" { inside = 1 }
		inside && $0 ~ /^P(2A|2B|[3-8])[|](active|queued)$/ {
			sub(/[|](active|queued)$/, "|complete")
		}
		{ print }
		$0 == "<!-- CODEX_AUTHORIZED_CHECKPOINTS_END -->" { inside = 0 }
	' "$plan" >"$replacement" || fail 'could not complete authorized queue fixture'
	mv "$replacement" "$plan"
}

[[ -x "$launcher" ]] || fail 'launcher is not executable'
test_shell=${BASH:-/bin/bash}
"$test_shell" -n "$launcher"

unsafe_header="$temp_root/unsafe-header-launcher"
unsafe_header_sentinel="$temp_root/unsafe-header-ran"
awk -v sentinel="$unsafe_header_sentinel" '
	{ print }
	$0 == "# CODEX_MUTABLE_SESSION_HEADER_BEGIN" {
		printf "printf unsafe >\047%s\047\n", sentinel
	}
' "$launcher" >"$unsafe_header"
chmod +x "$unsafe_header"
"$test_shell" "$unsafe_header" --help >/dev/null 2>&1 || true
[[ ! -e "$unsafe_header_sentinel" ]] || fail 'mutable header executed before validation'

unsafe_prompt="$temp_root/unsafe-prompt-launcher"
unsafe_prompt_sentinel="$temp_root/unsafe-prompt-ran"
awk -v sentinel="$unsafe_prompt_sentinel" '
	{ print }
	!injected && $0 == "# CODEX_MUTABLE_PROMPT_BEGIN" {
		printf "printf unsafe >\047%s\047\n", sentinel
		injected = 1
	}
' "$launcher" >"$unsafe_prompt"
chmod +x "$unsafe_prompt"
"$test_shell" "$unsafe_prompt" --help >/dev/null 2>&1 || true
[[ ! -e "$unsafe_prompt_sentinel" ]] || fail 'mutable prompt executed before validation'

validate_mutable_source_structure "$launcher" ||
	fail 'mutable regions are not inert tail data in the required normal form'
pass 'mutable data is unreachable shell input with a strict normal form'

normalized="$temp_root/normalized-launcher"
normalize_launcher_skeleton "$launcher" >"$normalized"
actual_skeleton_sha256=$(sha256_file "$normalized")
[[ "$actual_skeleton_sha256" == "$expected_skeleton_sha256" ]] ||
	fail "stable skeleton hash is $actual_skeleton_sha256, expected $expected_skeleton_sha256"
pass 'normalized stable skeleton hash'

active_status=$(extract_header_value "$launcher" 'SESSION_STATUS') || fail 'missing session status'
active_id=$(extract_header_value "$launcher" 'SESSION_ID') || fail 'missing session ID'
active_archive_rel=$(extract_header_value "$launcher" 'SESSION_ARCHIVE_REL') || fail 'missing session archive'
active_previous_rel=$(extract_header_value "$launcher" 'PREVIOUS_SESSION_ARCHIVE_REL') ||
	fail 'missing previous session archive'
case $active_status in
NEXT|COMPLETE) ;;
*) fail 'checked-in launcher has an invalid session status' ;;
esac
[[ "$active_archive_rel" == "docs/plan/agent-sessions/$active_id.md" ]] ||
	fail 'active archive is not derived from the session ID'
pass 'active metadata is derived from the launcher without test mutation'

source_graph_is_regular "$repo_root" ||
	fail 'checked-in launcher and planning graph must contain only regular files and directories'
repo_archives=("${SOURCE_ARCHIVES[@]}")
source_type_probe="$temp_root/source type probe"
mkdir -p "$source_type_probe/docs/design" \
	"$source_type_probe/docs/plan/agent-sessions"
cp "$launcher" "$source_type_probe/codex-dev-start.sh"
cp "$repo_root/docs/design/quality-lift.md" \
	"$repo_root/docs/design/agent-session-continuity.md" \
	"$source_type_probe/docs/design/"
cp "$repo_root/docs/plan/quality-upgrade.md" \
	"$repo_root/docs/plan/quality-handover.md" \
	"$source_type_probe/docs/plan/"
cp "${repo_archives[@]}" \
	"$source_type_probe/docs/plan/agent-sessions/"
mv "$source_type_probe/docs/plan/quality-handover.md" \
	"$temp_root/source-type-handover-target"
ln -s "$temp_root/source-type-handover-target" \
	"$source_type_probe/docs/plan/quality-handover.md"
if source_graph_is_regular "$source_type_probe"; then
	fail 'source graph type guard accepted a symlinked planning file'
fi
unlink "$source_type_probe/docs/plan/quality-handover.md"
mv "$temp_root/source-type-handover-target" \
	"$source_type_probe/docs/plan/quality-handover.md"
printf '# invalid hidden archive\n' \
	>"$source_type_probe/docs/plan/agent-sessions/.hidden-invalid.md"
if source_graph_is_regular "$source_type_probe"; then
	fail 'source graph type guard accepted a hidden archive'
fi
pass 'source graph guard rejects symlinked authority and hidden archives before copying'

real_fixture_root="$temp_root/checked-in source graph"
mkdir -p "$real_fixture_root/docs/design" \
	"$real_fixture_root/docs/plan/agent-sessions" \
	"$temp_root/empty-hooks" "$temp_root/empty-template"
cp "$launcher" "$real_fixture_root/codex-dev-start.sh"
chmod +x "$real_fixture_root/codex-dev-start.sh"
cp "$repo_root/docs/design/quality-lift.md" \
	"$repo_root/docs/design/agent-session-continuity.md" \
	"$real_fixture_root/docs/design/"
cp "$repo_root/docs/plan/quality-upgrade.md" \
	"$repo_root/docs/plan/quality-handover.md" \
	"$real_fixture_root/docs/plan/"
cp "${repo_archives[@]}" \
	"$real_fixture_root/docs/plan/agent-sessions/"
git_at "$real_fixture_root" init -q --template="$temp_root/empty-template"
git_at "$real_fixture_root" symbolic-ref HEAD refs/heads/codex/upgrade-quality
git_at "$real_fixture_root" config user.name 'Codex Launcher Test'
git_at "$real_fixture_root" config user.email 'codex-launcher-test@example.invalid'
git_at "$real_fixture_root" config commit.gpgSign false
git_at "$real_fixture_root" config core.hooksPath "$temp_root/empty-hooks"
git_at "$real_fixture_root" add .
git_at "$real_fixture_root" -c commit.gpgSign=false \
	-c core.hooksPath="$temp_root/empty-hooks" commit -qm 'checked-in source graph'
CODEX_BIN=true "$real_fixture_root/codex-dev-start.sh" --check \
	>"$temp_root/real-check.stdout" 2>"$temp_root/real-check.stderr"
assert_contains "$temp_root/real-check.stdout" "PASS (session $active_id, $active_status)"
[[ ! -s "$temp_root/real-check.stderr" ]] || fail 'clean checked-in source graph warned'
pass 'checked-in launcher and complete archive graph agree in a synthetic repository'

fixture_root="$temp_root/repository with spaces"
mkdir -p "$fixture_root/docs/design" "$fixture_root/docs/plan" "$fixture_root/test" \
	"$temp_root/empty-hooks" "$temp_root/empty-template"
fixture_root=$(cd -P "$fixture_root" && pwd)
cp "$launcher" "$fixture_root/codex-dev-start.sh"
chmod +x "$fixture_root/codex-dev-start.sh"
cp "$repo_root/test/codex_dev_start_test.sh" "$fixture_root/test/codex_dev_start_test.sh"
chmod +x "$fixture_root/test/codex_dev_start_test.sh"
active_status='NEXT'
active_id='2098-12-31T235959+0000-contract-next'
active_archive_rel="docs/plan/agent-sessions/$active_id.md"
active_previous_rel=''
replace_header_value "$fixture_root/codex-dev-start.sh" 'SESSION_STATUS' "$active_status"
replace_header_value "$fixture_root/codex-dev-start.sh" 'SESSION_ID' "$active_id"
replace_header_value "$fixture_root/codex-dev-start.sh" 'SESSION_ARCHIVE_REL' "$active_archive_rel"
replace_header_value "$fixture_root/codex-dev-start.sh" 'PREVIOUS_SESSION_ARCHIVE_REL' ''
printf '# fixture\n' >"$fixture_root/docs/design/quality-lift.md"
printf '# fixture\n' >"$fixture_root/docs/design/agent-session-continuity.md"
{
	printf '# fixture\n\n'
	printf '<!-- CODEX_AUTHORIZED_CHECKPOINTS_BEGIN -->\n'
	printf '%s\n' 'P2A|active' 'P2B|queued' 'P3|queued' 'P4|queued' \
		'P5|queued' 'P6|queued' 'P7|queued' 'P8|queued'
	printf '<!-- CODEX_AUTHORIZED_CHECKPOINTS_END -->\n'
} >"$fixture_root/docs/plan/quality-upgrade.md"
printf '# fixture\n' >"$fixture_root/docs/plan/quality-handover.md"
active_previous_leaf='none'
if [[ -n "$active_previous_rel" ]]; then
	active_previous_leaf=$(archive_leaf "$active_previous_rel")
fi
write_archive "$fixture_root" "$active_archive_rel" 'NEXT' \
	"$active_previous_leaf" 'none' 'pending' "$fixture_root/codex-dev-start.sh"

git_fixture init -q --template="$temp_root/empty-template"
git_fixture symbolic-ref HEAD refs/heads/codex/upgrade-quality
git_fixture config user.name 'Codex Launcher Test'
git_fixture config user.email 'codex-launcher-test@example.invalid'
git_fixture config commit.gpgSign false
git_fixture config core.hooksPath "$temp_root/empty-hooks"
commit_fixture 'test fixture'

fake_codex="$temp_root/fake codex"
cat >"$fake_codex" <<'FAKE_CODEX'
#!/usr/bin/env bash
set -euo pipefail
: "${FAKE_CODEX_RECORD_DIR:?}"
mkdir -p "$FAKE_CODEX_RECORD_DIR"
printf '%s\n' "$#" >"$FAKE_CODEX_RECORD_DIR/argc"
index=0
for argument in "$@"; do
	index=$((index + 1))
	printf '%s' "$argument" >"$FAKE_CODEX_RECORD_DIR/arg-$index"
done
exit "${FAKE_CODEX_EXIT:-0}"
FAKE_CODEX
chmod +x "$fake_codex"
mkdir -p "$temp_root/tools"
cp "$fake_codex" "$temp_root/tools/fake codex"
chmod +x "$temp_root/tools/fake codex"

"$launcher" --help >"$temp_root/help.stdout" 2>"$temp_root/help.stderr"
assert_contains "$temp_root/help.stdout" 'Usage: codex-dev-start.sh'
assert_contains "$temp_root/help.stdout" '--check'
assert_contains "$temp_root/help.stdout" '--print-prompt'
[[ ! -s "$temp_root/help.stderr" ]] || fail '--help wrote to stderr'
pass 'help does not require a prepared worktree'

expect_failure 'unknown argument' 'unknown argument: --unknown' "$launcher" --unknown
expect_failure 'multiple arguments' 'expected at most one argument' "$launcher" --check extra

if ! (
		cd "$temp_root"
		CDPATH="$temp_root" CODEX_BIN='tools/fake codex' \
			FAKE_CODEX_RECORD_DIR="$temp_root/not-called" \
			"repository with spaces/codex-dev-start.sh" --check
	) >"$temp_root/check.stdout" 2>"$temp_root/check.stderr"; then
	sed -n '1,40p' "$temp_root/check.stderr" >&2
	fail 'relative launcher or Codex path failed with inherited CDPATH'
fi
assert_contains "$temp_root/check.stdout" "PASS (session $active_id, NEXT)"
[[ ! -s "$temp_root/check.stderr" ]] || fail 'clean --check wrote a warning'
[[ ! -e "$temp_root/not-called" ]] || fail '--check invoked Codex'
pass 'check works from another directory without invoking Codex'

CODEX_BIN="$fake_codex" FAKE_CODEX_RECORD_DIR="$temp_root/not-called-print" \
	"$fixture_root/codex-dev-start.sh" --print-prompt \
	>"$temp_root/printed-prompt" 2>"$temp_root/printed-prompt.stderr"
[[ ! -e "$temp_root/not-called-print" ]] || fail '--print-prompt invoked Codex'
[[ ! -s "$temp_root/printed-prompt.stderr" ]] || fail 'clean prompt wrote a warning'
extract_archive_prompt "$fixture_root/$active_archive_rel" >"$temp_root/archive-prompt"
cmp -s "$temp_root/archive-prompt" "$temp_root/printed-prompt" ||
	fail 'printed prompt is not byte-exact archive input'
IFS= read -r first_prompt_line <"$temp_root/printed-prompt"
[[ "$first_prompt_line" == '# Mission' ]] || fail 'mission is not the first prompt content'
previous_line=0
for heading in \
	'# Mission' \
	'# Authorized Roadmap' \
	'# Measurements At Start' \
	'# Role And Boundaries' \
	'# Required Reading' \
	'# Three Moves' \
	'# Automatic Handoff'; do
	[[ $(grep -cFx "$heading" "$temp_root/printed-prompt") -eq 1 ]] ||
		fail "prompt section must occur exactly once: $heading"
	line=$(grep -nFx "$heading" "$temp_root/printed-prompt" | cut -d: -f1)
	[[ "$line" -gt "$previous_line" ]] || fail "prompt section is out of order: $heading"
	previous_line=$line
done
pass 'archived prompt is exact, mission-first, and ordered'

record_dir="$temp_root/codex call"
set +e
CODEX_BIN="$fake_codex" FAKE_CODEX_RECORD_DIR="$record_dir" FAKE_CODEX_EXIT=23 \
	"$fixture_root/codex-dev-start.sh" \
	>"$temp_root/start.stdout" 2>"$temp_root/start.stderr"
start_status=$?
set -e
[[ "$start_status" -eq 23 ]] || fail "Codex exit 23 became $start_status"
[[ $(cat "$record_dir/argc") -eq 5 ]] || fail 'Codex did not receive exactly five arguments'
[[ $(cat "$record_dir/arg-1") == '-c' ]] || fail 'first Codex argument is not -c'
[[ $(cat "$record_dir/arg-2") == 'service_tier="default"' ]] ||
	fail 'Codex service tier is not explicitly normal'
[[ $(cat "$record_dir/arg-3") == '-C' ]] || fail 'third Codex argument is not -C'
[[ $(cat "$record_dir/arg-4") == "$fixture_root" ]] || fail 'Codex worktree argument is wrong'
cmp -s "$record_dir/arg-5" "$temp_root/archive-prompt" ||
	fail 'Codex argument differs from archived bytes, including terminal LF'
[[ ! -s "$temp_root/start.stdout" ]] || fail 'launcher polluted Codex stdout'
[[ ! -s "$temp_root/start.stderr" ]] || fail 'clean start wrote a warning'
pass 'interactive invocation forces normal service tier and preserves exact prompt bytes'

dirty_name='private-untracked-filename-must-not-leak'
printf 'fixture\n' >"$fixture_root/$dirty_name"
CODEX_BIN="$fake_codex" FAKE_CODEX_RECORD_DIR="$temp_root/not-called-dirty" \
	"$fixture_root/codex-dev-start.sh" --print-prompt \
	>"$temp_root/dirty-prompt" 2>"$temp_root/dirty-prompt.stderr"
cmp -s "$temp_root/dirty-prompt" "$temp_root/archive-prompt" ||
	fail 'dirty state changed the archived prompt argument'
assert_contains "$temp_root/dirty-prompt.stderr" \
	'worktree has local changes; inspect them before editing.'
assert_not_contains "$temp_root/dirty-prompt" "$dirty_name"
assert_not_contains "$temp_root/dirty-prompt.stderr" "$dirty_name"
rm "$fixture_root/$dirty_name"
pass 'dirty worktree is allowed without prompt drift or path leakage'

real_git=$(command -v git)
git_wrapper_dir="$temp_root/git-wrapper"
mkdir -p "$git_wrapper_dir"
cat >"$git_wrapper_dir/git" <<'FAKE_GIT'
#!/usr/bin/env bash
set -euo pipefail
is_status='no'
for argument in "$@"; do
	if [[ "$argument" == 'status' ]]; then
		is_status='yes'
	fi
done
if [[ "$is_status" == 'yes' && "${FAKE_GIT_FAIL_STATUS:-0}" == '1' ]]; then
	exit 9
fi
if [[ "$is_status" == 'yes' && -n "${FAKE_GIT_SWAP_SOURCE:-}" ]]; then
	cp "$FAKE_GIT_SWAP_SOURCE" "$FAKE_GIT_SWAP_TARGET"
	chmod +x "$FAKE_GIT_SWAP_TARGET"
fi
exec "$REAL_GIT" "$@"
FAKE_GIT
chmod +x "$git_wrapper_dir/git"
expect_failure 'Git status failure' 'could not inspect worktree status' \
	env PATH="$git_wrapper_dir:$PATH" REAL_GIT="$real_git" FAKE_GIT_FAIL_STATUS=1 \
		CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check

cp "$fixture_root/codex-dev-start.sh" "$temp_root/pre-swap-launcher"
cp "$fixture_root/codex-dev-start.sh" "$temp_root/mutated-during-status"
insert_second_generation_prompt "$temp_root/mutated-during-status"
PATH="$git_wrapper_dir:$PATH" REAL_GIT="$real_git" \
	FAKE_GIT_SWAP_SOURCE="$temp_root/mutated-during-status" \
	FAKE_GIT_SWAP_TARGET="$fixture_root/codex-dev-start.sh" \
	CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --print-prompt \
	>"$temp_root/snapshot-prompt" 2>"$temp_root/snapshot-prompt.stderr"
cmp -s "$temp_root/snapshot-prompt" "$temp_root/archive-prompt" ||
	fail 'launcher did not use the prompt snapshot validated before Git status'
assert_contains "$temp_root/snapshot-prompt.stderr" \
	'worktree has local changes; inspect them before editing.'
cp "$temp_root/pre-swap-launcher" "$fixture_root/codex-dev-start.sh"
chmod +x "$fixture_root/codex-dev-start.sh"
pass 'Git status fails closed and prompt validation uses one source snapshot'

git_fixture switch -qc wrong-branch
expect_failure 'wrong branch' 'expected branch codex/upgrade-quality' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
git_fixture switch -q codex/upgrade-quality

mkdir -p "$fixture_root/nested"
cp "$launcher" "$fixture_root/nested/codex-dev-start.sh"
chmod +x "$fixture_root/nested/codex-dev-start.sh"
expect_failure 'wrong repository root' 'launcher directory is not a Git worktree' \
	env CODEX_BIN="$fake_codex" "$fixture_root/nested/codex-dev-start.sh" --check
rm -rf "$fixture_root/nested"

mv "$fixture_root/codex-dev-start.sh" "$temp_root/symlink-launcher-target"
ln -s "$temp_root/symlink-launcher-target" "$fixture_root/codex-dev-start.sh"
expect_failure 'symlink launcher' 'launcher must be a regular, non-symlink file' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
unlink "$fixture_root/codex-dev-start.sh"
mv "$temp_root/symlink-launcher-target" "$fixture_root/codex-dev-start.sh"
chmod +x "$fixture_root/codex-dev-start.sh"

expect_failure 'missing Codex executable' 'configured Codex executable is unavailable' \
	env CODEX_BIN="$temp_root/absent-codex" "$fixture_root/codex-dev-start.sh" --check
expect_failure 'shell builtin Codex target' 'configured Codex executable is unavailable' \
	env CODEX_BIN=source "$fixture_root/codex-dev-start.sh" --check

handover="$fixture_root/docs/plan/quality-handover.md"
mv "$handover" "$temp_root/quality-handover.md"
expect_failure 'missing planning document' 'missing required planning document' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
mv "$temp_root/quality-handover.md" "$handover"

mv "$handover" "$temp_root/quality-handover-symlink-target.md"
ln -s "$temp_root/quality-handover-symlink-target.md" "$handover"
expect_failure 'symlink planning file' 'planning document must be a regular, non-symlink file' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
unlink "$handover"
mv "$temp_root/quality-handover-symlink-target.md" "$handover"

mv "$fixture_root/docs/design" "$temp_root/design-directory-target"
ln -s "$temp_root/design-directory-target" "$fixture_root/docs/design"
expect_failure 'symlink planning directory' 'planning directory must not be a symlink' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
unlink "$fixture_root/docs/design"
mv "$temp_root/design-directory-target" "$fixture_root/docs/design"

archive="$fixture_root/$active_archive_rel"
cp "$archive" "$temp_root/archive-original"
rm "$archive"
expect_failure 'missing prompt archive' 'session prompt archive must be a regular, non-symlink file' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
cp "$temp_root/archive-original" "$archive"

rm "$archive"
ln -s "$temp_root/archive-original" "$archive"
expect_failure 'symlink prompt archive' 'session prompt archive must be a regular, non-symlink file' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
rm "$archive"
cp "$temp_root/archive-original" "$archive"

traversal_launcher="$fixture_root/codex-dev-start.sh"
cp "$traversal_launcher" "$temp_root/launcher-original"
replace_header_value "$traversal_launcher" 'SESSION_ARCHIVE_REL' \
	'docs/plan/agent-sessions/../quality-handover.md'
expect_failure 'archive path traversal' 'session archive path is outside the approved plan directory' \
	env CODEX_BIN="$fake_codex" "$traversal_launcher" --check
cp "$temp_root/launcher-original" "$traversal_launcher"
chmod +x "$traversal_launcher"
pass 'archive path traversal is rejected before file access'

cp "$traversal_launcher" "$temp_root/launcher-before-invalid-time"
invalid_time_id='2099-13-40T256199+1461-invalid-time'
replace_header_value "$traversal_launcher" 'SESSION_ID' "$invalid_time_id"
replace_header_value "$traversal_launcher" 'SESSION_ARCHIVE_REL' \
	"docs/plan/agent-sessions/$invalid_time_id.md"
expect_failure 'invalid session timestamp' 'session ID format is invalid' \
	env CODEX_BIN="$fake_codex" "$traversal_launcher" --check
cp "$temp_root/launcher-before-invalid-time" "$traversal_launcher"
chmod +x "$traversal_launcher"

sed -e 's/^Status: NEXT$/Status: ANSWERED - HISTORY/' \
	-e 's/^Outcome: pending$/Outcome: unexpectedly answered/' \
	"$temp_root/archive-original" >"$archive"
expect_failure 'active archive status drift' 'active archive status does not match launcher state' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
cp "$temp_root/archive-original" "$archive"

sed 's/^Created: .*$/Created: `2026-99-40T25:61:99+14:61`/' \
	"$temp_root/archive-original" >"$archive"
expect_failure 'invalid created metadata' 'archive Created metadata is invalid' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
cp "$temp_root/archive-original" "$archive"

sed 's/^Created: .*$/Created: `2026-08-25T06:15:32+02:00`/' \
	"$temp_root/archive-original" >"$archive"
expect_failure 'Created and ID mismatch' 'archive Created metadata does not match its session ID' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
cp "$temp_root/archive-original" "$archive"

sed 's/^Session ID: /Session ID: `wrong` # /' \
	"$temp_root/archive-original" >"$archive"
expect_failure 'archive ID drift' 'archive ID does not match its filename' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
cp "$temp_root/archive-original" "$archive"

sed 's/^# Mission$/# Altered Mission/' \
	"$temp_root/archive-original" >"$archive"
refresh_archive_digest "$archive"
expect_failure 'archive prompt drift' 'session archive prompt differs from launcher prompt' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
cp "$temp_root/archive-original" "$archive"

sed '/^<!-- CODEX_SESSION_PROMPT_END -->$/d' \
	"$temp_root/archive-original" >"$archive"
expect_failure 'malformed archive markers' 'session archive prompt markers are malformed' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
cp "$temp_root/archive-original" "$archive"

duplicate_id='2099-01-01T000000+0000-duplicate-next'
duplicate_rel="docs/plan/agent-sessions/$duplicate_id.md"
write_archive "$fixture_root" "$duplicate_rel" 'NEXT' 'none' 'none' 'pending' "$launcher"
expect_failure 'duplicate NEXT archive' 'archive set must contain exactly one NEXT session' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
rm "$fixture_root/$duplicate_rel"

orphan_id='2099-01-03T000000+0000-disconnected-history'
orphan_rel="docs/plan/agent-sessions/$orphan_id.md"
write_archive "$fixture_root" "$orphan_rel" 'ANSWERED - HISTORY' \
	'none' 'none' 'orphaned history' "$fixture_root/codex-dev-start.sh"
expect_failure 'disconnected history' 'archive set contains disconnected history' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
rm "$fixture_root/$orphan_rel"

old_launcher="$temp_root/first-generation-launcher"
cp "$fixture_root/codex-dev-start.sh" "$old_launcher"
new_id='2099-01-02T030405+0000-second-generation'
new_rel="docs/plan/agent-sessions/$new_id.md"
old_leaf=$(archive_leaf "$active_archive_rel")
new_leaf=$(archive_leaf "$new_rel")
replace_header_value "$fixture_root/codex-dev-start.sh" 'SESSION_ID' "$new_id"
replace_header_value "$fixture_root/codex-dev-start.sh" 'SESSION_ARCHIVE_REL' "$new_rel"
replace_header_value "$fixture_root/codex-dev-start.sh" 'PREVIOUS_SESSION_ARCHIVE_REL' "$active_archive_rel"
insert_second_generation_prompt "$fixture_root/codex-dev-start.sh"
write_archive "$fixture_root" "$active_archive_rel" 'ANSWERED - HISTORY' \
	'none' "$new_leaf" 'launcher contract passed' "$old_launcher"
write_archive "$fixture_root" "$new_rel" 'NEXT' "$old_leaf" 'none' 'pending' \
	"$fixture_root/codex-dev-start.sh"
commit_fixture 'second generation'
CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check \
	>"$temp_root/second-check.stdout" 2>"$temp_root/second-check.stderr"
assert_contains "$temp_root/second-check.stdout" "PASS (session $new_id, NEXT)"
[[ ! -s "$temp_root/second-check.stderr" ]] || fail 'clean second generation warned'
pass 'a second restart generation needs no launcher-test mutation'

if [[ "$nested_mode" == 'no' ]]; then
	source_archive_root="$temp_root/source archive"
	mkdir -p "$source_archive_root/test" "$source_archive_root/docs/design" \
		"$source_archive_root/docs/plan/agent-sessions"
	cp "$fixture_root/codex-dev-start.sh" "$source_archive_root/codex-dev-start.sh"
	chmod +x "$source_archive_root/codex-dev-start.sh"
	cp "$fixture_root/test/codex_dev_start_test.sh" \
		"$source_archive_root/test/codex_dev_start_test.sh"
	chmod +x "$source_archive_root/test/codex_dev_start_test.sh"
	cp "$fixture_root/docs/design/quality-lift.md" \
		"$fixture_root/docs/design/agent-session-continuity.md" \
		"$source_archive_root/docs/design/"
	cp "$fixture_root/docs/plan/quality-upgrade.md" \
		"$fixture_root/docs/plan/quality-handover.md" \
		"$source_archive_root/docs/plan/"
	source_graph_is_regular "$fixture_root" ||
		fail 'second-generation source graph contains a non-regular source'
	fixture_archives=("${SOURCE_ARCHIVES[@]}")
	cp "${fixture_archives[@]}" \
		"$source_archive_root/docs/plan/agent-sessions/"
	if ! "$test_shell" "$source_archive_root/test/codex_dev_start_test.sh" \
			--internal-source-archive-child \
			>"$temp_root/nested-test.stdout" 2>"$temp_root/nested-test.stderr"; then
		printf '%s\n' '--- nested stdout ---' >&2
		sed -n '1,120p' "$temp_root/nested-test.stdout" >&2
		printf '%s\n' '--- nested stderr ---' >&2
		sed -n '1,120p' "$temp_root/nested-test.stderr" >&2
		fail 'contract suite failed from a source archive with a second-generation header'
	fi
	assert_contains "$temp_root/nested-test.stdout" 'codex dev start contract: PASS (47 checks)'
	[[ ! -s "$temp_root/nested-test.stderr" ]] || fail 'nested second-generation test warned'
	pass 'contract suite boots from a source archive with a second-generation header'
	if CODEX_DEV_START_NESTED=1 "$test_shell" \
			"$source_archive_root/test/codex_dev_start_test.sh" \
			--internal-source-archive-child \
			>"$temp_root/ambient-nested.stdout" 2>"$temp_root/ambient-nested.stderr"; then
		fail 'ambient nested sentinel unexpectedly disabled contract controls'
	fi
	assert_contains "$temp_root/ambient-nested.stderr" \
		'CODEX_DEV_START_NESTED must be unset; nested mode is an internal argument'
	poisoned_git_config="$temp_root/poisoned-git-config"
	git config -f "$poisoned_git_config" core.fsmonitor /usr/bin/false
	CDPATH="$temp_root" GIT_CONFIG_GLOBAL="$poisoned_git_config" \
		GIT_DIR=/private/tmp/definitely-missing-git-dir \
		GIT_WORK_TREE=/private/tmp/definitely-missing-git-worktree \
		GIT_INDEX_FILE=/private/tmp/definitely-missing-git-index \
		"$test_shell" "$source_archive_root/test/codex_dev_start_test.sh" \
			--internal-environment-probe \
			>"$temp_root/environment-probe.stdout" \
			2>"$temp_root/environment-probe.stderr"
	assert_contains "$temp_root/environment-probe.stdout" \
		'codex dev start environment: PASS'
	[[ ! -s "$temp_root/environment-probe.stderr" ]] ||
		fail 'environment sanitation probe wrote to stderr'
	pass 'ambient state cannot disable or poison the source-archive contract'
fi

cp "$fixture_root/$active_archive_rel" "$temp_root/history-before-cycle"
sed "s|^Previous: none$|Previous: [$new_leaf]($new_leaf)|" \
	"$temp_root/history-before-cycle" >"$fixture_root/$active_archive_rel"
expect_failure 'archive cycle' 'archive chain contains a cycle' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
cp "$temp_root/history-before-cycle" "$fixture_root/$active_archive_rel"

awk '
	!mutated && $0 == "# Mission" {
		print "# Mutated Historical Mission"
		mutated = 1
		next
	}
	{ print }
	END { if (!mutated) exit 2 }
' "$temp_root/history-before-cycle" >"$fixture_root/$active_archive_rel" ||
	fail 'could not mutate historical prompt fixture'
expect_failure 'historical prompt drift' 'archive prompt digest does not match content' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
cp "$temp_root/history-before-cycle" "$fixture_root/$active_archive_rel"
pass 'cycles, disconnected history, and historical prompt drift fail closed'

mv "$fixture_root/$active_archive_rel" "$temp_root/missing-previous"
expect_failure 'missing previous archive' 'archive link target is missing' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
mv "$temp_root/missing-previous" "$fixture_root/$active_archive_rel"

cp "$fixture_root/$active_archive_rel" "$temp_root/previous-original"
sed "s/^Next: .*$/Next: none/" "$temp_root/previous-original" \
	>"$fixture_root/$active_archive_rel"
expect_failure 'broken archive backlink' 'archive links are not reciprocal' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
cp "$temp_root/previous-original" "$fixture_root/$active_archive_rel"
pass 'missing links and broken archive backlinks fail closed'

write_archive "$fixture_root" "$new_rel" 'ANSWERED - HISTORY' "$old_leaf" \
	'none' 'objective complete' "$fixture_root/codex-dev-start.sh"
replace_header_value "$fixture_root/codex-dev-start.sh" 'SESSION_STATUS' 'COMPLETE'
expect_failure 'premature COMPLETE with authorized work' \
	'session cannot be COMPLETE while authorized checkpoints remain' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
complete_authorized_queue "$fixture_root/docs/plan/quality-upgrade.md"
commit_fixture 'complete lifecycle'
CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check \
	>"$temp_root/complete-check.stdout" 2>"$temp_root/complete-check.stderr"
assert_contains "$temp_root/complete-check.stdout" "PASS (session $new_id, COMPLETE)"
[[ ! -s "$temp_root/complete-check.stderr" ]] || fail 'clean COMPLETE check warned'
cp "$fixture_root/$new_rel" "$temp_root/complete-tail-original"
sed "s|^Next: none$|Next: [$old_leaf]($old_leaf)|" \
	"$temp_root/complete-tail-original" >"$fixture_root/$new_rel"
expect_failure 'non-terminal COMPLETE tail' 'active archive tail must have Next: none' \
	env CODEX_BIN="$fake_codex" "$fixture_root/codex-dev-start.sh" --check
cp "$temp_root/complete-tail-original" "$fixture_root/$new_rel"
expect_failure 'completed session start' 'session has no NEXT task' \
	env CODEX_BIN="$fake_codex" FAKE_CODEX_RECORD_DIR="$temp_root/not-called-complete" \
	"$fixture_root/codex-dev-start.sh"
[[ ! -e "$temp_root/not-called-complete" ]] || fail 'COMPLETE session invoked Codex'
pass 'COMPLETE state validates history but cannot replay a task'

[[ -z $(git_fixture status --porcelain=v1 --untracked-files=all) ]] ||
	fail 'launcher contract tests left the fixture worktree dirty'
pass 'hermetic fixture remains clean after all lifecycle probes'

expected_pass_count=49
if [[ "$nested_mode" == 'yes' ]]; then
	expected_pass_count=47
fi
[[ "$pass_count" -eq "$expected_pass_count" ]] ||
	fail "contract control count is $pass_count, expected $expected_pass_count"
printf 'codex dev start contract: PASS (%s checks)\n' "$pass_count"
