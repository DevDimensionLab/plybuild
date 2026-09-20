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
Each turn receives an isolated scratch directory that is deleted automatically.
Successful supervisor logs are deleted; failed or interrupted logs are retained.
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

cleanup_active_scratch() {
	local scratch=${ACTIVE_SCRATCH_ROOT:-}
	local expected=${ACTIVE_TURN_DIR:-}
	[[ -n "$scratch" ]] || return 0
	if [[ -z "$expected" || "$scratch" != "$expected/scratch" ]]; then
		printf '%s\n' 'codex-dev-start: refusing to clean an unrecognized scratch path' >&2
		return 1
	fi
	if [[ -L "$scratch" || ( -e "$scratch" && ! -d "$scratch" ) ]]; then
		printf '%s\n' 'codex-dev-start: refusing to clean a non-directory scratch path' >&2
		return 1
	fi
	if [[ -d "$scratch" ]]; then
		find "$scratch" -type d -exec chmod u+rwx {} + 2>/dev/null || {
			printf '%s\n' 'codex-dev-start: could not make scratch directories removable' >&2
			return 1
		}
		rm -rf -- "$scratch" || {
			printf '%s\n' 'codex-dev-start: could not clean the active turn scratch directory' >&2
			return 1
		}
	fi
	ACTIVE_SCRATCH_ROOT=''
	ACTIVE_TURN_DIR=''
}

cleanup_supervisor_log_root() {
	local root=${SUPERVISOR_LOG_ROOT:-}
	local parent=${SUPERVISOR_TMP_ROOT:-}
	local prefix=${SUPERVISOR_LOG_PREFIX:-}
	local basename
	[[ -n "$root" ]] || return 0
	basename=${root##*/}
	if [[ -z "$parent" || -z "$prefix" || "${root%/*}" != "$parent" ||
		"$basename" != "$prefix"* ]]; then
		printf '%s\n' 'codex-dev-start: refusing to clean an unrecognized supervisor log path' >&2
		return 1
	fi
	if [[ -L "$root" || ( -e "$root" && ! -d "$root" ) ]]; then
		printf '%s\n' 'codex-dev-start: refusing to clean a non-directory supervisor log path' >&2
		return 1
	fi
	if [[ -d "$root" ]]; then
		find "$root" -type d -exec chmod u+rwx {} + 2>/dev/null || return 1
		rm -rf -- "$root" || return 1
	fi
	SUPERVISOR_LOG_ROOT=''
}

cleanup_supervisor() {
	cleanup_active_scratch || true
	if [[ -n "${LAUNCHER_SOURCE:-}" ]]; then
		rm -f "$LAUNCHER_SOURCE"
	fi
	if [[ -n "${ACTIVE_EVENT_FIFO:-}" ]]; then
		rm -f "$ACTIVE_EVENT_FIFO"
	fi
	if [[ "${SUPERVISOR_COMPLETED:-no}" == 'yes' ]]; then
		if cleanup_supervisor_log_root; then
			printf '%s\n' 'codex-dev-start: successful supervisor logs cleaned automatically' >&2
		else
			printf '%s\n' 'codex-dev-start: successful supervisor log cleanup failed' >&2
		fi
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
	SUPERVISOR_TMP_ROOT=$resolved_root
	SUPERVISOR_LOG_PREFIX="codex-dev-start.$SESSION_ID."
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
	ACTIVE_TURN_DIR=$turn_dir
	ACTIVE_SCRATCH_ROOT="$turn_dir/scratch"
	mkdir "$ACTIVE_SCRATCH_ROOT" || die 'could not create turn scratch directory'
	raw_log="$turn_dir/events.jsonl"
	child_stderr="$turn_dir/codex.stderr"
	final_message="$turn_dir/final-message.txt"
	ACTIVE_EVENT_FIFO="$turn_dir/events.fifo"
	mkfifo "$ACTIVE_EVENT_FIFO" || die 'could not create event stream pipe'

	printf 'codex-dev-start: starting session %s\n' "$SESSION_ID" >&2
	"$PYTHON_EXECUTABLE" "$EVENT_PARSER_PATH" "$raw_log" <"$ACTIVE_EVENT_FIFO" &
	ACTIVE_PARSER_PID=$!
	TMPDIR="$ACTIVE_SCRATCH_ROOT" \
		CODEX_SESSION_SCRATCH_ROOT="$ACTIVE_SCRATCH_ROOT" \
		"$CODEX_EXECUTABLE" exec --ephemeral -c "service_tier=\"$CODEX_SERVICE_TIER\"" \
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
	if ! cleanup_active_scratch; then
		printf 'codex-dev-start: scratch cleanup failed; inspect %s\n' "$turn_dir" >&2
		return 1
	fi

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
	ACTIVE_TURN_DIR=''
	ACTIVE_SCRATCH_ROOT=''
	SUPERVISOR_LOG_ROOT=''
	SUPERVISOR_TMP_ROOT=''
	SUPERVISOR_LOG_PREFIX=''
	SUPERVISOR_COMPLETED='no'
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
		2) SUPERVISOR_COMPLETED='yes'; exit 0 ;;
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
#|SESSION_ID=2026-09-20T163925+0200-decide-hashicorp-go-uuid-product-direction
#|SESSION_ARCHIVE_REL=docs/plan/agent-sessions/2026-09-20T163925+0200-decide-hashicorp-go-uuid-product-direction.md
#|PREVIOUS_SESSION_ARCHIVE_REL=docs/plan/agent-sessions/2026-09-20T152123+0200-evaluate-hashicorp-go-uuid-dependency.md
# CODEX_MUTABLE_SESSION_HEADER_END

# CODEX_MUTABLE_PROMPT_BEGIN
#|# Mission
#|
#|Continue P7 only by obtaining and recording one bounded product decision for
#|exact-path `github.com/hashicorp/go-uuid`. No exact-path stable release passes
#|every applicable behavior contract, so the preceding evaluation made no
#|dependency change. Choose among the bounded options below; do not silently
#|accept the findings, implement a dependency or source change, combine another
#|dependency group, or begin P8.
#|
#|# Authorized Roadmap
#|
#|P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
#|dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
#|and all target-specific retained-module decisions through exact inherited,
#|unloaded go-syslog v1.0.0. Every earlier outcome and lifecycle ancestor is
#|final. The evaluated go-uuid selection is exact inherited v1.0.1; no prior
#|exception transfers to it. Preserve every guarded selection and do not change
#|a parent, the Go floor, product source, or an unrelated module without the
#|corresponding fresh authorization.
#|
#|Selected `github.com/hashicorp/go-uuid v1.0.1` exists through six exact
#|v1.0.1 requests from Consul API v1.1.0/v1.12.0, Consul SDK v0.1.1/v0.8.0,
#|and Serf v0.8.2/v0.9.6. Go-immutable-radix v1.0.0/v1.3.1 additionally
#|request v1.0.0. There is no direct main-module root or repository import; its
#|why result is negative, production and complete-test loads contain zero target
#|packages, and it is runtime-unreachable. Those facts bound risk but do not
#|constitute qualification or acceptance.
#|
#|# Completed Evaluation
#|
#|Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves Hashicorp's
#|public active unarchived non-fork MPL-2.0 repository. The exact path has four
#|stable lightweight tags, v1.0.0 through proxy-latest v1.0.3, with no
#|prerelease, GitHub Release object, retraction, module deprecation, redirect,
#|alternate path, or `/v2` line. All proxy regular files match Git and sumdb
#|verifies each release. Current master
#|`f405b577e09f44ce7e9b484af2911859f9dab5a6` is 28 commits beyond v1.0.3,
#|declares Go 1.18, and is unreleased.
#|
#|Each stable release has one package, no external module dependency, and a
#|complete minimal production/test closure of target plus standard library.
#|Every release passes verification, listing, native tests, independent repeats,
#|race, vet, and production/test cross-builds under exact Go 1.26.7 and contained
#|Go 1.18.10. V1.0.0/v1.0.1 expose four functions; v1.0.2/v1.0.3 compatibly add
#|two reader-injection functions. The package intentionally formats random bytes
#|as UUID text rather than setting RFC version or variant bits.
#|
#|Every stable release loses underlying randomness-error identity by formatting
#|with `%v`. Under Go 1.18, v1.0.0/v1.0.1 return a non-wrapping error; under Go
#|1.26 their `crypto/rand.Read` path terminates the process when a failing global
#|reader is injected. V1.0.2/v1.0.3 add caller-owned reader APIs but retain the
#|same `%v` identity loss. All stable releases also panic for negative byte
#|counts despite returning an error. Unreleased master changes `%v` to `%w` and
#|passes error-identity fixtures under both SDKs, but retains the negative-size
#|panic. No stable candidate therefore qualifies.
#|
#|Disposable direct v1.0.1 adds only one edge and content checksum. Direct
#|v1.0.2/v1.0.3 alter only the target selection and add one edge, but tidy
#|removes the manufactured root and restores selected v1.0.1. Direct v1.0.0
#|breaks the guarded graph by removing mvn-pom-mutator v0.2.3 and makes project
#|loads fail. No projection was applied. The project remains 234 modules, 3,599
#|edges, 355 production entries, 429 complete-test entries, 197 module-backed
#|entries across 41 loaded modules, 1,067 sum lines, and the recorded 432-line
#|tidy projection. `go.mod`/`go.sum` SHA-256 values remain
#|`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
#|`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
#|
#|Fresh primary vulnerability data has no target record. Exact OSV, repository
#|advisory, and exact global GitHub advisory queries are empty for all four
#|releases. Exact-Go isolated scans are empty; old-Go reachable findings belong
#|only to its standard library. Base and v1.0.1/v1.0.2/v1.0.3 project scans are
#|identical, with zero target assignment or trace. Every earlier guarded
#|selection/request/load/advisory condition remains exact, and no new guarded
#|finding appeared.
#|
#|The exact-Go canonical project gate passes verify, build, count-one, two
#|count-ten repeats, race, vet, pinned lint, API/CLI, empty-HOME, four
#|cross-builds, full preflight, all 15 audit controls, all eight mutation
#|meta-stages, 80/80 kills, and host/snapshot/Docker acceptance. The fresh
#|clean-tree scorecard is 27/27 Q0-Q2 PASS at L2. No dependency implementation
#|commit exists.
#|
#|# Decision Required
#|
#|Choose exactly one direction:
#|
#|1. **Retain exact inherited, unloaded v1.0.1 under a target-specific
#|   exception (recommended).** Make no product-source or dependency-metadata
#|   change. Accept only the completed randomness-error identity, Go 1.26 fatal
#|   injected-reader, negative-size panic, non-RFC generation, API/behavior,
#|   MVS, vulnerability, and related recorded findings. The exception is
#|   go-uuid-specific and non-transferable. It remains valid only while exact
#|   v1.0.1, all six v1.0.1 requests and both v1.0.0 requests remain unchanged,
#|   no direct root or repository import is added, both target loads remain
#|   zero, runtime unreachability holds, every earlier guard remains intact, and
#|   no new advisory or independent defect appears. Any change requires the
#|   owning fresh dependency and product decision.
#|2. **Authorize a separately evaluated Go-1.18-compatible patch, fork, or
#|   replacement.** Scope a new bounded design/evaluation to preserve the
#|   required API, wrap reader errors, define non-panicking negative-size
#|   behavior, and prove provenance, maintenance, MVS, license, vulnerability,
#|   compatibility, and full quality effects before implementation. Do not
#|   implement it in this decision session.
#|3. **Authorize a separately evaluated parent-chain removal.** Scope a fresh
#|   graph/product study across the exact Consul API, Consul SDK, Serf, and
#|   go-immutable-radix requesters. This may affect earlier guarded decisions and
#|   cannot be performed or assumed here.
#|
#|Selecting v1.0.3 alone is not a qualifying fourth option: it retains both
#|behavior blockers, introduces a direct root that tidy removes, and provides no
#|loaded project consumer benefit. An unreleased master, redirect, fork, or
#|alternate path is likewise not an authorized substitute.
#|
#|# Required Reading And Guards
#|
#|At start read this archive, its answered go-uuid evaluation, the answered
#|go-syslog decision and evaluation, the answered go-sockaddr, go-rootcerts,
#|go-retryablehttp, go-multierror, and go-hclog decisions, rolling handover,
#|roadmap, `go.mod`, and `go.sum`. Verify branch, clean ordinary and ignored
#|status, handoff ancestry and exact changed-file set, reciprocal archive chain,
#|exact Go 1.26.7 identity, all 20 guarded selections/requests/why/import/load
#|conditions, module hashes, advisory guards, and
#|`./codex-dev-start.sh --check`. Stop if any guard differs.
#|
#|# Three Moves
#|
#|First, obtain one explicit product choice and do nothing else while it is
#|pending. Second, revalidate the narrow guards and record exactly the authorized
#|choice, limits, expiry triggers, and no-change result; a patch/replacement or
#|parent-removal choice authorizes only its next bounded evaluation, not an
#|implementation here. Third, update the roadmap and rolling handover, answer
#|this archive, prepare exactly one reciprocal NEXT mission for the authorized
#|result, and commit the handoff. Do not execute the successor.
#|
#|# Automatic Handoff
#|
#|After one coherent bounded outcome, make only the required local
#|`docs: prepare next agent session` commit. Do not create a dependency commit,
#|launch a successor, push, merge, publish, release, stash, revert, bypass
#|cleanup, remove the worktree, combine another dependency group, or begin P8.
# CODEX_MUTABLE_PROMPT_END
