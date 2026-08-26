#!/usr/bin/env bash

set -euo pipefail
unset CDPATH

EXPECTED_BRANCH='codex/upgrade-quality'
ARCHIVE_PREFIX='docs/plan/agent-sessions/'
CODEX_SERVICE_TIER='default'

usage() {
	cat <<EOF
Usage: $(basename "$0") [--check|--print-prompt|--help]

Supervise fresh non-interactive Codex turns in normal service mode.

  --check         Validate the worktree, archive chain, and Codex executable.
  --print-prompt  Validate and print the prompt without starting Codex.
  --help          Show this help text.

CODEX_BIN may name an alternate Codex executable for contract tests.
Raw JSONL and diagnostics are stored in a unique directory outside the worktree.
EOF
}

die() {
	printf 'codex-dev-start: %s\n' "$*" >&2
	exit 1
}

extract_header_value() {
	local key=$1
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
	' "$LAUNCHER_SOURCE"
}

extract_launcher_prompt() {
	local source_file=$1
	awk '
		$0 == "# CODEX_MUTABLE_PROMPT_BEGIN" {
			begin_count++
			if (begin_count != 1 || inside || end_count) invalid = 1
			inside = 1
			next
		}
		$0 == "# CODEX_MUTABLE_PROMPT_END" {
			end_count++
			if (!inside || end_count != 1) invalid = 1
			inside = 0
			next
		}
		inside {
			if (substr($0, 1, 2) != "#|") invalid = 1
			print substr($0, 3)
		}
		END {
			if (begin_count != 1 || end_count != 1 || inside || invalid) exit 2
		}
	' "$source_file"
}

validate_mutable_regions() {
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
			if (prompt_lines == 1 && substr($0, 3) != "# Mission") exit 2
			next
		}
		{ previous = $0 }
		END {
			if (region || stable_count != 1 || header_begin_count != 1 ||
			    header_end_count != 1 || prompt_begin_count != 1 ||
			    prompt_end_count != 1 || header_begin_line <= stable_line ||
			    prompt_begin_line <= header_begin_line) exit 2
		}
	' "$LAUNCHER_SOURCE" || die 'launcher mutable regions are malformed'
}

compact_timestamp_is_valid() {
	printf '%s\n' "$1" | awk '
		/^[0-9][0-9][0-9][0-9]-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])T([01][0-9]|2[0-3])[0-5][0-9][0-5][0-9][+-](0[0-9]|1[0-4])[0-5][0-9]$/ {
			year = substr($0, 1, 4) + 0
			month = substr($0, 6, 2) + 0
			day = substr($0, 9, 2) + 0
			offset_hour = substr($0, 19, 2) + 0
			offset_minute = substr($0, 21, 2) + 0
			max_day = 31
			if (month == 4 || month == 6 || month == 9 || month == 11) max_day = 30
			if (month == 2) {
				max_day = 28
				if (year % 400 == 0 || (year % 4 == 0 && year % 100 != 0)) max_day = 29
			}
			if (day <= max_day && !(offset_hour == 14 && offset_minute != 0)) valid = 1
		}
		END { exit !valid }
	'
}

session_id_is_safe() {
	local id=$1
	local timestamp=${id:0:22}
	local checkpoint=${id:22}
	compact_timestamp_is_valid "$timestamp" || return 1
	printf '%s\n' "$checkpoint" | awk '
		/^-[a-z0-9]+(-[a-z0-9]+)*$/ { valid = 1 }
		END { exit !valid }
	'
}

archive_rel_is_safe() {
	local rel=$1
	local leaf
	local id
	case $rel in
	"$ARCHIVE_PREFIX"*.md) ;;
	*) return 1 ;;
	esac
	leaf=${rel#"$ARCHIVE_PREFIX"}
	[[ "$leaf" != */* ]] || return 1
	id=${leaf%.md}
	[[ "$leaf" == "$id.md" ]] || return 1
	session_id_is_safe "$id"
}

archive_field() {
	local archive=$1
	local key=$2
	awk -v prefix="$key: " '
		BEGIN { before_prompt = 1 }
		$0 == "<!-- CODEX_SESSION_PROMPT_BEGIN -->" { before_prompt = 0 }
		before_prompt && index($0, prefix) == 1 {
			count++
			value = substr($0, length(prefix) + 1)
		}
		END {
			if (count != 1) exit 2
			print value
		}
	' "$archive"
}

extract_archive_prompt() {
	local archive=$1
	awk '
		$0 == "<!-- CODEX_SESSION_PROMPT_BEGIN -->" {
			begin_count++
			if (begin_count != 1 || inside || end_count) invalid = 1
			inside = 1
			next
		}
		$0 == "<!-- CODEX_SESSION_PROMPT_END -->" {
			end_count++
			if (!inside || end_count != 1) invalid = 1
			inside = 0
			next
		}
		inside { print }
		END {
			if (begin_count != 1 || end_count != 1 || inside || invalid) exit 2
		}
	' "$archive"
}

archive_link_rel() {
	local value=$1
	local name
	local rel
	[[ "$value" != 'none' ]] || return 1
	case $value in
	\[*\]\(*\)) ;;
	*) return 2 ;;
	esac
	name=${value#\[}
	name=${name%%\]*}
	[[ "$value" == "[$name]($name)" ]] || return 2
	rel="$ARCHIVE_PREFIX$name"
	archive_rel_is_safe "$rel" || return 2
	printf '%s\n' "$rel"
}

created_value_is_valid() {
	local value=$1
	local timestamp
	case $value in
	\`????-??-??T??:??:??[+-]??:??\`) ;;
	*) return 1 ;;
	esac
	timestamp=${value#\`}
	timestamp=${timestamp%\`}
	timestamp=${timestamp//:/}
	compact_timestamp_is_valid "$timestamp"
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
		return 1
	fi
}

validate_archive_file() {
	local rel=$1
	local archive="$REPO_ROOT/$rel"
	local leaf=${rel#"$ARCHIVE_PREFIX"}
	local expected_id=${leaf%.md}
	local status
	local session_id
	local created
	local created_timestamp
	local source
	local prompt_digest
	local prompt_file
	local actual_prompt_digest
	local previous
	local next
	local outcome

	archive_rel_is_safe "$rel" || die 'archive filename is invalid'
	[[ -f "$archive" && ! -L "$archive" ]] ||
		die 'session prompt archive must be a regular, non-symlink file'
	status=$(archive_field "$archive" 'Status') || die 'archive Status metadata is invalid'
	case $status in
	NEXT|'ANSWERED - HISTORY') ;;
	*) die 'archive Status metadata is invalid' ;;
	esac
	session_id=$(archive_field "$archive" 'Session ID') || die 'archive Session ID metadata is invalid'
	[[ "$session_id" == "\`$expected_id\`" ]] || die 'archive ID does not match its filename'
	created=$(archive_field "$archive" 'Created') || die 'archive Created metadata is invalid'
	created_value_is_valid "$created" || die 'archive Created metadata is invalid'
	created_timestamp=${created#\`}
	created_timestamp=${created_timestamp%\`}
	created_timestamp=${created_timestamp//:/}
	[[ "$created_timestamp" == "${expected_id:0:22}" ]] ||
		die 'archive Created metadata does not match its session ID'
	source=$(archive_field "$archive" 'Source') || die 'archive Source metadata is invalid'
	[[ "$source" == '`codex-dev-start.sh`' ]] || die 'archive Source metadata is invalid'
	prompt_digest=$(archive_field "$archive" 'Prompt SHA-256') ||
		die 'archive prompt digest metadata is invalid'
	previous=$(archive_field "$archive" 'Previous') || die 'archive Previous metadata is invalid'
	if [[ "$previous" != 'none' ]]; then
		archive_link_rel "$previous" >/dev/null || die 'archive Previous metadata is invalid'
	fi
	next=$(archive_field "$archive" 'Next') || die 'archive Next metadata is invalid'
	if [[ "$next" != 'none' ]]; then
		archive_link_rel "$next" >/dev/null || die 'archive Next metadata is invalid'
	fi
	outcome=$(archive_field "$archive" 'Outcome') || die 'archive Outcome metadata is invalid'
	[[ -n "$outcome" ]] || die 'archive Outcome metadata is invalid'
	if [[ "$status" == 'NEXT' ]]; then
		[[ "$next" == 'none' && "$outcome" == 'pending' ]] ||
			die 'NEXT archive lifecycle metadata is invalid'
	else
		[[ "$outcome" != 'pending' ]] || die 'answered archive outcome is still pending'
	fi
	prompt_file=$(mktemp "${TMPDIR:-/tmp}/codex-dev-start-archive-prompt.XXXXXX") ||
		die 'could not create archive prompt validation state'
	if ! extract_archive_prompt "$archive" >"$prompt_file"; then
		rm -f "$prompt_file"
		die 'session archive prompt markers are malformed'
	fi
	actual_prompt_digest=$(sha256_file "$prompt_file") || {
		rm -f "$prompt_file"
		die 'no SHA-256 implementation is available'
	}
	rm -f "$prompt_file"
	[[ "$prompt_digest" == "\`$actual_prompt_digest\`" ]] ||
		die 'archive prompt digest does not match content'
}

validate_archive_set() {
	local active="$REPO_ROOT/$SESSION_ARCHIVE_REL"
	local archive
	local rel
	local status
	local next_count=0
	local total_count=0
	local active_status='NEXT'
	local active_previous
	local active_next
	local expected_previous='none'
	local visited
	local current_rel
	local child_rel=''
	local current_previous
	local parent_next
	local visited_count

	[[ -f "$active" && ! -L "$active" ]] ||
		die 'session prompt archive must be a regular, non-symlink file'
	shopt -s nullglob dotglob
	ARCHIVE_FILES=("$ARCHIVE_DIR"/*.md)
	shopt -u nullglob dotglob
	[[ ${#ARCHIVE_FILES[@]} -gt 0 ]] || die 'session archive directory is empty'
	for archive in "${ARCHIVE_FILES[@]}"; do
		rel="$ARCHIVE_PREFIX${archive##*/}"
		validate_archive_file "$rel"
		status=$(archive_field "$archive" 'Status') || die 'archive Status metadata is invalid'
		if [[ "$status" == 'NEXT' ]]; then
			next_count=$((next_count + 1))
		fi
		total_count=$((total_count + 1))
	done

	if [[ "$SESSION_STATUS" == 'COMPLETE' ]]; then
		active_status='ANSWERED - HISTORY'
	fi
	status=$(archive_field "$active" 'Status') || die 'archive Status metadata is invalid'
	[[ "$status" == "$active_status" ]] ||
		die 'active archive status does not match launcher state'
	active_next=$(archive_field "$active" 'Next') || die 'archive Next metadata is invalid'
	[[ "$active_next" == 'none' ]] || die 'active archive tail must have Next: none'
	if [[ "$SESSION_STATUS" == 'NEXT' ]]; then
		[[ "$next_count" -eq 1 ]] ||
			die 'archive set must contain exactly one NEXT session'
	else
		[[ "$next_count" -eq 0 ]] ||
			die 'archive set must contain no NEXT sessions when COMPLETE'
	fi

	active_previous=$(archive_field "$active" 'Previous') ||
		die 'archive Previous metadata is invalid'
	if [[ -n "$PREVIOUS_SESSION_ARCHIVE_REL" ]]; then
		expected_previous="[${PREVIOUS_SESSION_ARCHIVE_REL##*/}](${PREVIOUS_SESSION_ARCHIVE_REL##*/})"
	fi
	[[ "$active_previous" == "$expected_previous" ]] ||
		die 'active archive previous link does not match launcher header'

	visited=$(mktemp "${TMPDIR:-/tmp}/codex-dev-start-visited.XXXXXX") ||
		die 'could not create archive validation state'
	current_rel=$SESSION_ARCHIVE_REL
	while :; do
		if grep -Fx -- "$current_rel" "$visited" >/dev/null; then
			rm -f "$visited"
			die 'archive chain contains a cycle'
		fi
		printf '%s\n' "$current_rel" >>"$visited"
		archive="$REPO_ROOT/$current_rel"
		[[ -f "$archive" && ! -L "$archive" ]] || {
			rm -f "$visited"
			die 'archive link target is missing'
		}
		if [[ -n "$child_rel" ]]; then
			status=$(archive_field "$archive" 'Status') || {
				rm -f "$visited"
				die 'archive Status metadata is invalid'
			}
			[[ "$status" == 'ANSWERED - HISTORY' ]] || {
				rm -f "$visited"
				die 'archive predecessor is not answered history'
			}
			parent_next=$(archive_field "$archive" 'Next') || {
				rm -f "$visited"
				die 'archive Next metadata is invalid'
			}
			[[ "$parent_next" == "[${child_rel##*/}](${child_rel##*/})" ]] || {
				rm -f "$visited"
				die 'archive links are not reciprocal'
			}
		fi
		current_previous=$(archive_field "$archive" 'Previous') || {
			rm -f "$visited"
			die 'archive Previous metadata is invalid'
		}
		[[ "$current_previous" != 'none' ]] || break
		child_rel=$current_rel
		current_rel=$(archive_link_rel "$current_previous") || {
			rm -f "$visited"
			die 'archive Previous metadata is invalid'
		}
	done
	visited_count=$(wc -l <"$visited" | tr -d '[:space:]')
	rm -f "$visited"
	[[ "$visited_count" -eq "$total_count" ]] ||
		die 'archive set contains disconnected history'
	ARCHIVE_TOTAL_COUNT=$total_count
}

validate_active_prompt() {
	local temp_dir
	local prompt_with_sentinel
	temp_dir=$(mktemp -d "${TMPDIR:-/tmp}/codex-dev-start-prompt.XXXXXX") ||
		die 'could not create prompt validation state'
	if ! extract_launcher_prompt "$LAUNCHER_SOURCE" >"$temp_dir/launcher"; then
		rm -rf "$temp_dir"
		die 'launcher prompt is malformed'
	fi
	if ! extract_archive_prompt "$REPO_ROOT/$SESSION_ARCHIVE_REL" >"$temp_dir/archive"; then
		rm -rf "$temp_dir"
		die 'session archive prompt markers are malformed'
	fi
	if ! cmp -s "$temp_dir/launcher" "$temp_dir/archive"; then
		rm -rf "$temp_dir"
		die 'session archive prompt differs from launcher prompt'
	fi
	prompt_with_sentinel=$(cat "$temp_dir/launcher" && printf x) || {
		rm -rf "$temp_dir"
		die 'could not capture the validated launcher prompt'
	}
	SESSION_PROMPT=${prompt_with_sentinel%x}
	rm -rf "$temp_dir"
}

resolve_codex() {
	local candidate
	local candidate_dir

	if [[ "$CODEX_BIN" == */* ]]; then
		candidate=$CODEX_BIN
	else
		candidate=$(type -P -- "$CODEX_BIN" 2>/dev/null) ||
			die 'configured Codex executable is unavailable'
	fi
	[[ -f "$candidate" && -x "$candidate" ]] ||
		die 'configured Codex executable is unavailable'
	candidate_dir=$(cd -P "$(dirname "$candidate")" 2>/dev/null && pwd) ||
		die 'configured Codex executable is unavailable'
	CODEX_EXECUTABLE="$candidate_dir/$(basename "$candidate")"
	[[ -f "$CODEX_EXECUTABLE" && -x "$CODEX_EXECUTABLE" ]] ||
		die 'configured Codex executable is unavailable'
}

resolve_python() {
	local candidate
	local candidate_dir
	candidate=$(type -P -- python3 2>/dev/null) ||
		die 'python3 is required to validate Codex JSONL events'
	[[ -f "$candidate" && -x "$candidate" ]] ||
		die 'python3 is required to validate Codex JSONL events'
	candidate_dir=$(cd -P "$(dirname "$candidate")" 2>/dev/null && pwd) ||
		die 'python3 is required to validate Codex JSONL events'
	PYTHON_EXECUTABLE="$candidate_dir/$(basename "$candidate")"
}

normalized_launcher_digest() {
	local source_file=$1
	local normalized
	local digest
	normalized=$(mktemp "${TMPDIR:-/tmp}/codex-dev-start-skeleton.XXXXXX") ||
		die 'could not create launcher contract validation state'
	if ! awk '
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
	' "$source_file" >"$normalized"; then
		rm -f "$normalized"
		die 'launcher contract normalization failed'
	fi
	digest=$(sha256_file "$normalized") || {
		rm -f "$normalized"
		die 'no SHA-256 implementation is available'
	}
	rm -f "$normalized"
	printf '%s\n' "$digest"
}

validate_authorized_queue() {
	local plan="$REPO_ROOT/docs/plan/quality-upgrade.md"
	local counts
	local active_count
	local queued_count
	counts=$(awk '
		$0 == "<!-- CODEX_AUTHORIZED_CHECKPOINTS_BEGIN -->" {
			begin_count++
			if (inside || begin_count != 1) invalid = 1
			inside = 1
			next
		}
		$0 == "<!-- CODEX_AUTHORIZED_CHECKPOINTS_END -->" {
			end_count++
			if (!inside || end_count != 1) invalid = 1
			inside = 0
			next
		}
		inside {
			entry_count++
			split($0, fields, "|")
			expected[1] = "P2A"
			expected[2] = "P2B"
			expected[3] = "P3"
			expected[4] = "P4"
			expected[5] = "P5"
			expected[6] = "P6"
			expected[7] = "P7"
			expected[8] = "P8"
			if (fields[1] != expected[entry_count] ||
			    fields[2] !~ /^(active|queued|complete)$/ || fields[3] != "") invalid = 1
			if (fields[2] == "active") active++
			if (fields[2] == "queued") queued++
			if (fields[2] == "complete") {
				if (active || queued) invalid = 1
			} else if (fields[2] == "active") {
				if (active != 1 || queued) invalid = 1
			} else if (!active) {
				invalid = 1
			}
			next
		}
		END {
			if (inside || begin_count != 1 || end_count != 1 || entry_count != 8 ||
			    invalid || active > 1) exit 2
			printf "%d %d\n", active, queued
		}
	' "$plan") || die 'authorized checkpoint queue is malformed'
	active_count=${counts%% *}
	queued_count=${counts#* }
	AUTHORIZED_ACTIVE_COUNT=$active_count
	AUTHORIZED_QUEUED_COUNT=$queued_count
	if [[ "$SESSION_STATUS" == 'COMPLETE' ]]; then
		[[ "$active_count" -eq 0 && "$queued_count" -eq 0 ]] ||
			die 'session cannot be COMPLETE while authorized checkpoints remain'
	else
		[[ "$active_count" -eq 1 ]] ||
			die 'NEXT session requires exactly one active authorized checkpoint'
	fi
}

cleanup_supervisor() {
	if [[ -n "${LAUNCHER_SOURCE:-}" ]]; then
		rm -f "$LAUNCHER_SOURCE"
	fi
	if [[ -n "${ACTIVE_EVENT_FIFO:-}" ]]; then
		rm -f "$ACTIVE_EVENT_FIFO"
	fi
}

forward_supervisor_signal() {
	SUPERVISOR_INTERRUPTED='yes'
	printf '%s\n' 'codex-dev-start: interruption requested; stopping the active turn.' >&2
	if [[ -n "${ACTIVE_CODEX_PID:-}" ]]; then
		kill -TERM "$ACTIVE_CODEX_PID" 2>/dev/null || true
	fi
	if [[ -n "${ACTIVE_PARSER_PID:-}" ]]; then
		kill -TERM "$ACTIVE_PARSER_PID" 2>/dev/null || true
	fi
}

snapshot_launcher_source() {
	if [[ -n "${LAUNCHER_SOURCE:-}" ]]; then
		rm -f "$LAUNCHER_SOURCE"
	fi
	[[ -f "$SCRIPT_PATH" && ! -L "$SCRIPT_PATH" ]] ||
		die 'launcher must be a regular, non-symlink file'
	LAUNCHER_SOURCE=$(mktemp "${TMPDIR:-/tmp}/codex-dev-start-source.XXXXXX") ||
		die 'could not create launcher source snapshot'
	if ! cp "$SCRIPT_PATH" "$LAUNCHER_SOURCE"; then
		rm -f "$LAUNCHER_SOURCE"
		LAUNCHER_SOURCE=''
		die 'could not snapshot launcher source'
	fi
}

validate_session_state() {
	local git_root
	local active_branch
	local required_file

	snapshot_launcher_source
	validate_mutable_regions
	SESSION_STATUS=$(extract_header_value 'SESSION_STATUS') || die 'session status is missing'
	SESSION_ID=$(extract_header_value 'SESSION_ID') || die 'session ID is missing'
	SESSION_ARCHIVE_REL=$(extract_header_value 'SESSION_ARCHIVE_REL') ||
		die 'session archive path is missing'
	PREVIOUS_SESSION_ARCHIVE_REL=$(extract_header_value 'PREVIOUS_SESSION_ARCHIVE_REL') ||
		die 'previous session archive path is missing'

	session_id_is_safe "$SESSION_ID" || die 'session ID format is invalid'
	archive_rel_is_safe "$SESSION_ARCHIVE_REL" ||
		die 'session archive path is outside the approved plan directory'
	[[ "$SESSION_ARCHIVE_REL" == "$ARCHIVE_PREFIX$SESSION_ID.md" ]] ||
		die 'session archive path does not match the session ID'
	if [[ -n "$PREVIOUS_SESSION_ARCHIVE_REL" ]]; then
		archive_rel_is_safe "$PREVIOUS_SESSION_ARCHIVE_REL" ||
			die 'previous session archive path is outside the approved plan directory'
		[[ "$PREVIOUS_SESSION_ARCHIVE_REL" != "$SESSION_ARCHIVE_REL" ]] ||
			die 'previous session archive cannot equal the active archive'
	fi

	[[ -d "$REPO_ROOT/.git" || -f "$REPO_ROOT/.git" ]] ||
		die 'launcher directory is not a Git worktree'
	git_root=$(git -C "$REPO_ROOT" rev-parse --show-toplevel 2>/dev/null) ||
		die 'launcher directory is not a Git worktree'
	git_root=$(cd -P "$git_root" && pwd)
	[[ "$git_root" == "$REPO_ROOT" ]] || die 'launcher must run from the repository root'

	active_branch=$(git -C "$REPO_ROOT" symbolic-ref --quiet --short HEAD 2>/dev/null) ||
		die 'worktree must be attached to the expected branch'
	[[ "$active_branch" == "$EXPECTED_BRANCH" ]] || die "expected branch $EXPECTED_BRANCH"

	for required_file in \
		'docs/design/agent-session-continuity.md' \
		'docs/design/quality-lift.md' \
		'docs/plan/quality-upgrade.md' \
		'docs/plan/quality-handover.md'; do
		[[ -f "$REPO_ROOT/$required_file" ]] ||
			die "missing required planning document: $required_file"
		[[ ! -L "$REPO_ROOT/$required_file" ]] ||
			die "planning document must be a regular, non-symlink file: $required_file"
	done
	for required_file in 'docs' 'docs/design' 'docs/plan' 'docs/plan/agent-sessions'; do
		[[ -d "$REPO_ROOT/$required_file" && ! -L "$REPO_ROOT/$required_file" ]] ||
			die "planning directory must not be a symlink: $required_file"
	done
	ARCHIVE_DIR=$(cd -P "$REPO_ROOT/$ARCHIVE_PREFIX" && pwd) ||
		die 'session archive directory is unavailable'
	[[ "$ARCHIVE_DIR" == "$REPO_ROOT/${ARCHIVE_PREFIX%/}" ]] ||
		die 'session archive directory escapes the repository'

	validate_authorized_queue
	validate_archive_set
	validate_active_prompt
	resolve_codex
	resolve_python
	CURRENT_SKELETON_DIGEST=$(normalized_launcher_digest "$LAUNCHER_SOURCE")
	CURRENT_HEAD=$(git -C "$REPO_ROOT" rev-parse HEAD 2>/dev/null) ||
		die 'could not inspect worktree HEAD'
	WORKTREE_STATUS=$(git -C "$REPO_ROOT" status --porcelain=v1 --untracked-files=all) ||
		die 'could not inspect worktree status'
}

create_supervisor_log_root() {
	local requested_root=${TMPDIR:-/tmp}
	local resolved_root
	local created
	resolved_root=$(cd -P "$requested_root" 2>/dev/null && pwd) ||
		die 'supervisor log root is unavailable'
	case "$resolved_root/" in
	"$REPO_ROOT/"*) die 'supervisor log directory must be outside the worktree' ;;
	esac
	created=$(mktemp -d "$resolved_root/codex-dev-start.$SESSION_ID.XXXXXX") ||
		die 'could not create external supervisor log directory'
	SUPERVISOR_LOG_ROOT=$(cd -P "$created" && pwd) ||
		die 'could not resolve external supervisor log directory'
	case "$SUPERVISOR_LOG_ROOT/" in
	"$REPO_ROOT/"*) die 'supervisor log directory must be outside the worktree' ;;
	esac
	printf 'codex-dev-start: logs: %s\n' "$SUPERVISOR_LOG_ROOT" >&2
}

write_event_parser() {
	EVENT_PARSER_PATH="$SUPERVISOR_LOG_ROOT/validate-events.py"
	cat >"$EVENT_PARSER_PATH" <<'PY'
import json
import sys


def concise(value):
    if not isinstance(value, str):
        return ""
    single_line = " ".join(value.split())
    return single_line if len(single_line) <= 180 else single_line[:177] + "..."


raw_path = sys.argv[1]
errors = []
event_count = 0
thread_started = 0
turn_started = 0
completed = 0
failure_terminals = 0
terminal_seen = False

with open(raw_path, "wb") as raw_stream:
    for line_number, raw_line in enumerate(sys.stdin.buffer, 1):
        raw_stream.write(raw_line)
        raw_stream.flush()
        event_count += 1
        if not raw_line.endswith(b"\n"):
            errors.append("line {} is truncated (missing terminal LF)".format(line_number))
        try:
            event = json.loads(raw_line)
        except (UnicodeDecodeError, json.JSONDecodeError) as error:
            errors.append("line {} is not valid JSON: {}".format(line_number, error))
            continue
        if not isinstance(event, dict) or not isinstance(event.get("type"), str):
            errors.append("line {} is not a JSON object with a string type".format(line_number))
            continue
        event_type = event["type"]
        if terminal_seen:
            errors.append("event {} appears after a terminal event".format(event_type))
        if event_type == "thread.started":
            thread_started += 1
            thread_id = concise(event.get("thread_id"))
            print("codex-dev-start: turn thread started{}".format(
                " (" + thread_id + ")" if thread_id else ""), file=sys.stderr, flush=True)
        elif event_type == "turn.started":
            turn_started += 1
        elif event_type in ("item.started", "item.completed"):
            item = event.get("item")
            if isinstance(item, dict):
                item_type = item.get("type")
                if item_type == "agent_message" and event_type == "item.completed":
                    print("codex-dev-start: agent: " + concise(item.get("text")),
                          file=sys.stderr, flush=True)
                elif item_type == "command_execution":
                    command = concise(item.get("command"))
                    if command:
                        state = "completed" if event_type == "item.completed" else "started"
                        print("codex-dev-start: command {}: {}".format(state, command),
                              file=sys.stderr, flush=True)
                elif item_type == "file_change" and event_type == "item.completed":
                    print("codex-dev-start: file change completed", file=sys.stderr, flush=True)
        if event_type == "turn.completed":
            terminal_seen = True
            completed += 1
            print("codex-dev-start: Codex turn completed", file=sys.stderr, flush=True)
        elif event_type in ("turn.failed", "error"):
            terminal_seen = True
            failure_terminals += 1
            detail = concise(event.get("error")) or concise(event.get("message"))
            print("codex-dev-start: Codex terminal failure{}".format(
                ": " + detail if detail else ""), file=sys.stderr, flush=True)

if event_count == 0:
    errors.append("event stream is empty")
if thread_started != 1:
    errors.append("expected exactly one thread.started event, found {}".format(thread_started))
if turn_started != 1:
    errors.append("expected exactly one turn.started event, found {}".format(turn_started))
if completed != 1:
    errors.append("expected exactly one turn.completed event, found {}".format(completed))
if failure_terminals:
    errors.append("stream contains {} failure/error terminal event(s)".format(failure_terminals))

if errors:
    for error in errors:
        print("codex-dev-start: invalid Codex JSONL: " + error, file=sys.stderr)
    raise SystemExit(1)
PY
	chmod 600 "$EVENT_PARSER_PATH"
}

run_codex_turn() {
	local turn_number=$1
	local turn_label
	local turn_dir
	local raw_log
	local child_stderr
	local final_message
	local codex_rc
	local parser_rc

	turn_label=$(printf '%03d' "$turn_number")
	turn_dir="$SUPERVISOR_LOG_ROOT/turn-$turn_label-$SESSION_ID"
	mkdir "$turn_dir" || die 'could not create turn log directory'
	raw_log="$turn_dir/events.jsonl"
	child_stderr="$turn_dir/codex.stderr"
	final_message="$turn_dir/final-message.txt"
	ACTIVE_EVENT_FIFO="$turn_dir/events.fifo"
	mkfifo "$ACTIVE_EVENT_FIFO" || die 'could not create event stream pipe'

	printf 'codex-dev-start: starting session %s\n' "$SESSION_ID" >&2
	"$PYTHON_EXECUTABLE" "$EVENT_PARSER_PATH" "$raw_log" <"$ACTIVE_EVENT_FIFO" &
	ACTIVE_PARSER_PID=$!
	"$CODEX_EXECUTABLE" exec -c "service_tier=\"$CODEX_SERVICE_TIER\"" \
		--sandbox workspace-write -C "$REPO_ROOT" --json \
		--output-last-message "$final_message" "$SESSION_PROMPT" \
		>"$ACTIVE_EVENT_FIFO" 2>"$child_stderr" &
	ACTIVE_CODEX_PID=$!

	set +e
	wait "$ACTIVE_CODEX_PID"
	codex_rc=$?
	if [[ "$SUPERVISOR_INTERRUPTED" == 'yes' ]] && kill -0 "$ACTIVE_CODEX_PID" 2>/dev/null; then
		kill -TERM "$ACTIVE_CODEX_PID" 2>/dev/null || true
		wait "$ACTIVE_CODEX_PID"
		codex_rc=$?
	fi
	ACTIVE_CODEX_PID=''
	wait "$ACTIVE_PARSER_PID"
	parser_rc=$?
	ACTIVE_PARSER_PID=''
	set -e
	rm -f "$ACTIVE_EVENT_FIFO"
	ACTIVE_EVENT_FIFO=''

	if [[ "$SUPERVISOR_INTERRUPTED" == 'yes' ]]; then
		printf 'codex-dev-start: interrupted; inspect %s\n' "$turn_dir" >&2
		return 130
	fi
	if [[ "$codex_rc" -ne 0 ]]; then
		printf 'codex-dev-start: Codex exited %s; inspect %s and %s\n' \
			"$codex_rc" "$raw_log" "$child_stderr" >&2
		return 1
	fi
	if [[ "$parser_rc" -ne 0 ]]; then
		printf 'codex-dev-start: Codex event validation failed; inspect %s\n' "$raw_log" >&2
		return 1
	fi
	return 0
}

added_archives_between() {
	git -C "$REPO_ROOT" diff --name-only --diff-filter=A "$1" "$2" -- "$ARCHIVE_PREFIX" ||
		die 'could not inspect committed archive progression'
}

validate_post_turn_progression() {
	local previous_head=$1
	local previous_id=$2
	local previous_archive_rel=$3
	local previous_archive_count=$4
	local previous_archive="$REPO_ROOT/$previous_archive_rel"
	local previous_status
	local previous_next
	local expected_next
	local added_archives
	local added_count

	validate_session_state
	[[ "$CURRENT_SKELETON_DIGEST" == "$SUPERVISOR_SKELETON_DIGEST" ]] ||
		die 'post-turn launcher contract drifted outside its mutable regions'
	[[ -z "$WORKTREE_STATUS" ]] ||
		die 'post-turn worktree is dirty; refusing to continue'
	[[ "$CURRENT_HEAD" != "$previous_head" ]] ||
		die 'Codex turn made no committed HEAD progress'
	added_archives=$(added_archives_between "$previous_head" "$CURRENT_HEAD")
	if [[ -n "$added_archives" ]]; then
		added_count=$(printf '%s\n' "$added_archives" | wc -l | tr -d '[:space:]')
	else
		added_count=0
	fi

	if [[ "$SESSION_STATUS" == 'COMPLETE' ]]; then
		[[ "$SESSION_ID" == "$previous_id" && "$SESSION_ARCHIVE_REL" == "$previous_archive_rel" ]] ||
			die 'COMPLETE handoff changed the terminal session identity'
		previous_status=$(archive_field "$previous_archive" 'Status') ||
			die 'former archive Status metadata is invalid after the turn'
		[[ "$previous_status" == 'ANSWERED - HISTORY' ]] ||
			die 'former archive was not answered by the committed handoff'
		previous_next=$(archive_field "$previous_archive" 'Next') ||
			die 'former archive Next metadata is invalid after the turn'
		[[ "$ARCHIVE_TOTAL_COUNT" -eq "$previous_archive_count" && "$added_count" -eq 0 ]] ||
			die 'COMPLETE handoff added an unexpected archive'
		[[ "$previous_next" == 'none' ]] || die 'COMPLETE archive is not terminal'
		printf 'codex-dev-start: authorized roadmap COMPLETE at %s\n' "$CURRENT_HEAD" >&2
		return 2
	fi

	[[ "$SESSION_STATUS" == 'NEXT' ]] || die 'post-turn launcher state is invalid'
	[[ "$SESSION_ID" != "$previous_id" ]] ||
		die 'post-turn session identity did not change'
	previous_status=$(archive_field "$previous_archive" 'Status') ||
		die 'former archive Status metadata is invalid after the turn'
	[[ "$previous_status" == 'ANSWERED - HISTORY' ]] ||
		die 'former archive was not answered by the committed handoff'
	previous_next=$(archive_field "$previous_archive" 'Next') ||
		die 'former archive Next metadata is invalid after the turn'
	[[ "$PREVIOUS_SESSION_ARCHIVE_REL" == "$previous_archive_rel" ]] ||
		die 'new session does not identify the former archive as its predecessor'
	[[ "$ARCHIVE_TOTAL_COUNT" -eq $((previous_archive_count + 1)) ]] ||
		die 'post-turn archive graph did not add exactly one session'
	[[ "$added_count" -eq 1 && "$added_archives" == "$SESSION_ARCHIVE_REL" ]] ||
		die 'committed handoff did not add exactly the active NEXT archive'
	expected_next="[${SESSION_ARCHIVE_REL##*/}](${SESSION_ARCHIVE_REL##*/})"
	[[ "$previous_next" == "$expected_next" ]] ||
		die 'former archive does not link to the new NEXT session'
	git -C "$REPO_ROOT" cat-file -e "$CURRENT_HEAD:$SESSION_ARCHIVE_REL" 2>/dev/null ||
		die 'active NEXT archive is not committed at HEAD'
	printf 'codex-dev-start: committed handoff validated: %s -> %s\n' \
		"$previous_id" "$SESSION_ID" >&2
	return 0
}

main() {
	local mode='start'
	local turn_count=0
	local previous_head
	local previous_id
	local previous_archive_rel
	local previous_archive_count
	local turn_rc
	local progress_rc

	case $# in
	0) ;;
	1)
		case $1 in
		--check) mode='check' ;;
		--print-prompt) mode='print' ;;
		--help|-h)
			usage
			exit 0
			;;
		*) die "unknown argument: $1" ;;
		esac
		;;
	*) die 'expected at most one argument' ;;
	esac

	SCRIPT_DIR=$(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd)
	SCRIPT_PATH="$SCRIPT_DIR/$(basename "${BASH_SOURCE[0]}")"
	REPO_ROOT=$SCRIPT_DIR
	CODEX_BIN=${CODEX_BIN:-codex}
	LAUNCHER_SOURCE=''
	ACTIVE_EVENT_FIFO=''
	ACTIVE_CODEX_PID=''
	ACTIVE_PARSER_PID=''
	SUPERVISOR_INTERRUPTED='no'
	trap cleanup_supervisor EXIT

	validate_session_state
	if [[ -n "$WORKTREE_STATUS" ]]; then
		printf '%s\n' \
			'codex-dev-start: worktree has local changes; inspect them before editing.' >&2
	fi

	case $mode in
	check)
		printf 'codex-dev-start: PASS (session %s, %s)\n' "$SESSION_ID" "$SESSION_STATUS"
		exit 0
		;;
	print)
		[[ "$SESSION_STATUS" == 'NEXT' ]] || die 'session has no NEXT task'
		printf '%s' "$SESSION_PROMPT"
		exit 0
		;;
	start) ;;
	*) die 'internal mode error' ;;
	esac

	if [[ "$SESSION_STATUS" == 'COMPLETE' ]]; then
		printf 'codex-dev-start: authorized roadmap already COMPLETE at %s\n' "$CURRENT_HEAD" >&2
		exit 0
	fi
	[[ "$SESSION_STATUS" == 'NEXT' ]] || die 'session has no NEXT task'
	[[ -z "$WORKTREE_STATUS" ]] || die 'worktree must be clean before starting a supervised turn'

	SUPERVISOR_SKELETON_DIGEST=$CURRENT_SKELETON_DIGEST
	create_supervisor_log_root
	write_event_parser
	trap forward_supervisor_signal HUP INT TERM

	while [[ "$SESSION_STATUS" == 'NEXT' ]]; do
		turn_count=$((turn_count + 1))
		previous_head=$CURRENT_HEAD
		previous_id=$SESSION_ID
		previous_archive_rel=$SESSION_ARCHIVE_REL
		previous_archive_count=$ARCHIVE_TOTAL_COUNT

		set +e
		run_codex_turn "$turn_count"
		turn_rc=$?
		set -e
		if [[ "$turn_rc" -eq 130 ]]; then
			exit 130
		fi
		[[ "$turn_rc" -eq 0 ]] || die 'supervised Codex turn failed'

		set +e
		validate_post_turn_progression "$previous_head" "$previous_id" \
			"$previous_archive_rel" "$previous_archive_count"
		progress_rc=$?
		set -e
		case $progress_rc in
		0) ;;
		2) exit 0 ;;
		*) exit "$progress_rc" ;;
		esac
	done

	die 'supervisor loop ended without NEXT or COMPLETE state'
}

main "$@"
exit 70
# CODEX_STABLE_EXECUTION_END

# CODEX_MUTABLE_SESSION_HEADER_BEGIN
#|SESSION_STATUS=NEXT
#|SESSION_ID=2026-08-26T204147+0200-repair-q17-partial-failures
#|SESSION_ARCHIVE_REL=docs/plan/agent-sessions/2026-08-26T204147+0200-repair-q17-partial-failures.md
#|PREVIOUS_SESSION_ARCHIVE_REL=docs/plan/agent-sessions/2026-08-26T195729+0200-repair-q19-empty-populations.md
# CODEX_MUTABLE_SESSION_HEADER_END

# CODEX_MUTABLE_PROMPT_BEGIN
#|# Mission
#|
#|Continue P4 with one Q1.7-only test move. Repair the eight confirmed operations
#|whose per-item failure content, continuation, or ordered partial result lacks
#|an exact executable contract. Reclassify all 17 in-scope operations from
#|executable control flow and do not infer coverage from the 23 existing focused
#|tests or from exit-code assertions alone.
#|
#|# Authorized Roadmap
#|
#|P3 is complete. P4 remains active only for Q1.7 remediation and later combined
#|manual evidence; P5-P8 remain queued in the machine-readable block in
#|`docs/plan/quality-upgrade.md`. Q1.6 is independently supported by its complete
#|57-struct, 74-field review, and Q1.9 is independently supported by its complete
#|169-range classification. Neither support permits a schema-2 receipt in this
#|test-only session. The launcher must remain NEXT until P4 and P5-P8 finish.
#|
#|This mission authorizes read-only reclassification of the 17 Q1.7 operations;
#|test-only characterization of `initCmd`, `Bitbucket.SynchronizeAllRepos`,
#|`spring.DeleteDemoFiles`, `Repository.upgradeDependencies`,
#|`Repository.upgradePluginsOnModel`, `maven.RemoveDeprecated`,
#|`template.MergeTemplates`, and `GitCloudConfig.ValidTemplatesFrom`; focused
#|test helpers needed to exercise their existing behavior; one focused test-only
#|implementation commit; clean automated measurement; and the normal separate
#|continuity commit. It does not authorize production changes, another seam,
#|manual evidence or receipts, Q1.9 changes, audit/parser/scanner changes,
#|`.quality/inventory` or baseline changes, executable P5 harnesses, T1-T10
#|changes, acceptance changes, P6-P8 work, or a Q1.7 PASS claim unless the
#|complete classified population supports it.
#|
#|# Measurements At Start
#|
#|The Q1.9 move produced test-only commit
#|`3a8f0bcd030baf787a29440ee8e4e4a087edb33b`, exact tree
#|`75dee3dcc1cb2fd91f7eb359e96b31cc4be3d986`, on unchanged production commit
#|`c382caa38be167fe17f847370ad8a12270644de3`. Its exhaustive 169-range record
#|classifies 72 repaired, 83 independently guarded, and 14 fixture/support range
#|sites; all 64 local `tests` loops are now directly guarded. Q1.9 is supported,
#|but no manual-evidence document or receipt was created.
#|
#|The clean no-evidence scorecard SHA-256 is
#|`08e5ff9d450a6fa3b8824d82c720d339cf0d44003a490d01d3a412681081a615`.
#|It records 453 tests across all 27 packages, L0 8/8, L1 six PASS and
#|Q1.6/Q1.7/Q1.9 UNMEASURABLE, six improved ratchets, two held, zero regressed,
#|zero non-comparable, and zero dirty paths. The audit exits 1 for 13 documented
#|non-passing criteria, never 2.
#|
#|# Q1.7 Scope And Blockers
#|
#|Include a production collection-item failure path when it reports the failure
#|and continues, and an operation when it can return a populated partial
#|aggregate with an error. Exclude fail-fast errors before partial state,
#|intentionally ignored errors with no reported partial result, non-error
#|business warnings, sequential non-collection phases, and test support. Inspect
#|executable control flow; names and search hits are not classification evidence.
#|
#|The complete scope is 17 operations. Nine have exact content coverage. The
#|eight confirmed gaps are the `initCmd` project-loop failures;
#|`SynchronizeAllRepos` clone/pull warning; `DeleteDemoFiles` discovery/deletion
#|failures; dependency version, maximum-version, and upgrade failures; per-plugin
#|upgrade failures; `RemoveDeprecated` partial-result and warn/continue paths;
#|`MergeTemplates` merge warning; and `ValidTemplatesFrom` ordered partial result
#|with error. A new contract must assert the exact observable failure content and
#|also prove later-item continuation or the exact populated partial result where
#|that behavior exists. Prefer existing doubles and log-capture helpers. Preserve
#|all production behavior and process-state restoration.
#|
#|# Role And Boundaries
#|
#|Work autonomously in this worktree on `codex/upgrade-quality`. Confirm the new
#|continuity HEAD and exact ancestry before editing. Make one focused test-only
#|implementation commit, then the separate `docs: prepare next agent session`
#|continuity commit. Do not push, merge, publish, distribute, remove the
#|worktree, stash inherited changes, revert user work, or run destructive Git
#|commands.
#|
#|If any confirmed blocker cannot be exercised through current seams without a
#|production change, exhaust safe test-only approaches, record the exact
#|controllability gap, and keep Q1.7 blocked. Do not change production or add a
#|seam under this mission. Preserve test names and data already in place, exact
#|failure strings and ordering, public behavior, side effects, callback order,
#|and Go 1.18 compatibility.
#|
#|Keep every generated report outside the worktree. Set
#|`API_COMPAT_REPORT_OUT` and `CLI_COMPAT_REPORT_OUT` to external paths whenever
#|running compatibility. The measured-tree identity counts ignored, untracked,
#|and tracked bytes, so confirm both ordinary and ignored status before every
#|authoritative audit.
#|
#|# Required Reading
#|
#|Before editing, confirm branch, HEAD, clean and ignored status, reciprocal
#|archive links, launcher `--check`, exact ancestry, and the authorized
#|checkpoint block. Read the rolling handover, this archive, the complete P4
#|entry and gate, both design documents, `.quality/README.md`, the Q1.7 criterion
#|and schema-2 validation path, all 17 production operations and enclosing call
#|flows, every existing partial-failure test and helper, the complete audit
#|wrapper, T15 repair, API/CLI and launcher contracts, and Make and acceptance
#|gates. Do not classify an operation by name or search hit alone.
#|
#|# Three Moves
#|
#|1. Regenerate the complete 17-operation Q1.7 record from the clean current
#|   tree. For each operation record path, function, every included failure
#|   branch, existing exact-content/continuation coverage, new contract, or exact
#|   exclusion outside the worktree. Add focused test-only contracts for every
#|   exercisable branch among the eight gaps. Require exact failure content and
#|   later-item continuation or exact ordered partial results as applicable.
#|2. Run the exact affected tests, the complete 17-operation focused population,
#|   complete uncached tests, race, and vet. Inspect every new assertion in its
#|   executable context. If any included failure path can still occur without an
#|   exact content assertion and continuation/partial-result proof, keep Q1.7
#|   blocked. Do not create manual evidence in this session.
#|3. Run API/CLI and subprocess compatibility with external report paths; pinned
#|   lint; launcher and Make contracts; complete preflight; all four host
#|   acceptance flows; the 15-control audit meta-suite; empty-HOME count-2; and a
#|   clean no-evidence audit. Require audit exit never 2 and zero dirty paths.
#|   Record the exact population, implementation commit, scorecard, and next P4
#|   blocker or evidence-only next step in continuity.
#|
#|# Automatic Handoff
#|
#|Before this session ends, finish and commit the coherent Q1.7 test-only move or
#|record an exact resumable blocker. Rewrite the rolling handover, record the P4
#|result in the roadmap, answer this archive, create exactly one reciprocally
#|linked NEXT archive for the next coherent P4 move, replace only the launcher's
#|mutable regions, run launcher and handoff contracts, and make the separate
#|continuity-only commit. Do not launch a real successor, push, merge, publish,
#|distribute, stash, revert, or remove the worktree. COMPLETE remains invalid
#|while P4 or P5-P8 is unfinished.
# CODEX_MUTABLE_PROMPT_END
