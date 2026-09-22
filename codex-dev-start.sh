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
#|SESSION_ID=2026-09-22T202511+0200-decide-mitchellh-mapstructure-product-direction
#|SESSION_ARCHIVE_REL=docs/plan/agent-sessions/2026-09-22T202511+0200-decide-mitchellh-mapstructure-product-direction.md
#|PREVIOUS_SESSION_ARCHIVE_REL=docs/plan/agent-sessions/2026-09-22T192114+0200-evaluate-mitchellh-mapstructure-dependency.md
# CODEX_MUTABLE_SESSION_HEADER_END

# CODEX_MUTABLE_PROMPT_BEGIN
#|# Mission
#|
#|Continue P7 by recording exactly one bounded product decision for selected
#|indirect exact-path `github.com/mitchellh/mapstructure v1.5.0`. The completed
#|evaluation found that no canonical exact-path stable qualifies: every one of
#|the 17 Go-1.18-compatible stable releases retains stale destination slice
#|elements during an ordinary documented decode, and v1.2.2 also fails its
#|repeated upstream suite. The exact-path repository is archived and its owner
#|states that it will never tag another release; the blessed active fork uses
#|alternate module path `github.com/go-viper/mapstructure/v2` and is not an
#|eligible exact-path release. Choose only one of the three directions below,
#|apply it exactly, and prepare only its reciprocal successor if the choice
#|requires one. Do not repeat the dependency evaluation, evaluate another
#|dependency group, or begin P8.
#|
#|# Defensive Decision Scope
#|
#|This is an ordinary dependency-quality product decision. Use only the
#|completed public release/repository metadata, static source and graph facts,
#|small bounded ordinary map-decoding results, project projections, and
#|advisory identities recorded here and in the answered evaluation. Do not
#|fuzz, stress, probe resource exhaustion, create oversized, deeply nested,
#|cyclic, malformed, adversarial, or escape-sequence payloads, reproduce a
#|security issue, or perform security or exploitability analysis.
#|
#|Every disposable cache, tool, report, project copy, fixture, or advisory
#|response must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write
#|to `/private/tmp`, `/tmp`, a sibling of the managed root, or another external
#|root. Verify containment and remove task-owned scratch evidence before
#|handoff.
#|
#|# Authorized Roadmap
#|
#|P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
#|dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
#|and target-specific option-1 decisions through exact selected, direct,
#|production-loaded, runtime-relevant go-homedir v1.1.0. P8 remains queued.
#|
#|Go-homedir, Promptui, emoji/v2, kr/text, and kr/pty option-1 decisions are
#|separately exact, unqualified, target-specific, non-transferable, and final.
#|Their selected and historical requests, requester imports, complete routes,
#|why/import/load/runtime facts, release/repository/source identities,
#|qualification results, graph/module/tidy/Go-floor state, projections,
#|advisories, all earlier guards, and compatible-route conditions remain
#|expiry guards. Any change requires the corresponding fresh dependency and
#|product decision before merge. Do not transfer an exception or reopen
#|go-homedir, Promptui, emoji/v2, kr/text, kr/pty, kr/pretty, Cast, Viper, or an
#|earlier decision.
#|
#|# Completed Evaluation Is Final
#|
#|The evaluation began from clean branch `codex/upgrade-quality` at handoff
#|HEAD `bff26c941f30ef51a4f5a35c064a558ef1b45d16`, parent
#|`0f2bf5aa44b1c0f57003cea4f17dad96740e55d2`, tree
#|`6771799cf803b4591bcb56d086aeed8c534e2019`. Exact Google UUID v1.4.0
#|dependency commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an
#|ancestor. The evaluation changed no product source or dependency metadata
#|and prepared this decision-only handoff. Verify the new handoff HEAD, parent,
#|tree, exact changed set, ancestry, reciprocal archive chain, and clean
#|ordinary and ignored status rather than assuming them.
#|
#|Exact go-import metadata and all release module files resolve
#|`github.com/mitchellh/mapstructure` to public, enabled, archived, non-fork
#|MIT repository `mitchellh/mapstructure`, ID 10166531, owned by `mitchellh`
#|and defaulting to `main`. The proxy exposes exactly 17 stable releases from
#|v1.0.0 through v1.5.0; `@latest` is v1.5.0 at
#|2022-04-20T22:31:31Z. Exact `/v2` and `/v3` lines are absent. There are no
#|GitHub Releases, retractions, module deprecations, replacements, or
#|prereleases. All 17 tags are lightweight and form continuous ancestry;
#|v1.0.0 has an unverifiable old key and the later 16 commits are validly
#|signed. V1.0.0-v1.1.2 have no Go directive and v1.2.0-v1.5.0 declare Go
#|1.14. All have no requirements, replacements, or retractions.
#|
#|Selected v1.5.0 is validly signed commit
#|`ab69d8d93410fce4361f4912bb1ff88110a81311`, tree
#|`5a1166013faa55170e1acce958c1a3863faf3f80`, parent
#|`bd687ea300c090473812a1a5730c3a335fbb5b72`. Its proxy source/module sums
#|are `h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=` and
#|`h1:bFUtVrKA4DC2yAKiSyO/QUcy7e+RRV2QTWOzhPopBRo=`. Its proxy ZIP has
#|SHA-256 `118d5b2cb65c50dba967fb6d708f450a9caf93f321f8fc99080675b2ee374199`
#|and byte-matches its exact 13-file Git regular-file manifest at
#|`7c590ec67ae92865b8628f7f8c7d15487f31c7b83e4846941482a051920bcf09`.
#|Sixteen proxy archives byte-match current tag manifests. V1.2.0 is the sole
#|historical tag-mutation exception: its immutable proxy/sumdb archive matches
#|historical commit `047abd31f2839526c057afa2966f1bc4a374d025`, while the current
#|tag points to signed descendant `9e4011917e467353a5cacb6b11f3991b381055c8`
#|and differs only by two CHANGELOG lines. The selected archive is unaffected.
#|
#|The original owner states that archived repositories will accept no further
#|issues or pull requests and will never receive another tag, and names
#|`go-viper/mapstructure` as the blessed fork. That active repository is a fork
#|of the exact-path repository but its current module declaration is
#|`github.com/go-viper/mapstructure/v2` with Go 1.18. It is an alternate path,
#|not an eligible release or owner for the exact module line. No canonical
#|exact-path stable therefore has a genuine supported owner.
#|
#|V1.5.0 is a standard-library-only, single-package 13-file module with no
#|build tags, cgo, generated code, embed, platform-specific implementation,
#|symlink, submodule, package-owned mutable global, external I/O, or cleanup
#|boundary. Its exported surface comprises Decode/WeakDecode variants,
#|NewDecoder and Decoder.Decode, DecoderConfig and Metadata, Error, decode-hook
#|types/execution/composition, and conversion hooks. Decoder/config/output/
#|metadata/hook state is caller-owned; callers must not reuse a config after
#|NewDecoder and must synchronize shared state. Independent decoders are
#|deterministic and race-free in the completed bounded fixtures. Errors are
#|sorted for deterministic rendering. The reflection and conversion boundary
#|is the core behavior under evaluation.
#|
#|Under exact Go 1.26.7 and contained Go 1.18.10, all 17 releases pass module
#|verification, build, count-one tests, race count-one tests, vet, and Darwin
#|amd64/arm64, Linux amd64/arm64/386, Windows amd64/386, FreeBSD amd64, Plan 9
#|amd64, and js/wasm cross-builds. Repeated count-ten tests pass except at v1.2.2,
#|whose upstream metadata-order assertions fail under both SDKs. Selected
#|v1.5.0 has complete production/test closures of 76/133 packages under Go
#|1.26.7 and 57/88 under Go 1.18.10, with three module-backed entries and one
#|loaded module. Small bounded ordinary deterministic, invalid-output, and
#|four-independent-decoder concurrency fixtures pass count-ten and race-count-
#|ten under both SDKs.
#|
#|The decisive ordinary fixture decodes a one-element string slice into a
#|pre-populated two-element destination. Every stable from v1.0.0 through
#|v1.5.0 returns `["new" "stale"]` instead of replacing it with the documented
#|one-element result `["new"]` under both exact SDKs.
#|Upstream merged PR #266 three days after the selected tag and commit
#|`33d262eb9c0f5add618477081144108325cbc55c` adds the missing destination
#|slice truncation; main records it for unreleased v1.5.1. Later unreleased
#|commits also fix decode-hook error wrapping and a named-string
#|TextUnmarshaller panic. Those branch fixes are not eligible releases. No
#|canonical exact-path stable qualifies.
#|
#|The graph contains nine exact target requests. Main and Viper v1.15.0 request
#|v1.5.0. Historical Viper v1.10.1 and crypt v0.4.0 request v1.4.3; historical
#|Kong `v0.2.1-0.20190708041108-0548c6b1afae` and Consul API v1.1.0/v1.12.0
#|request v1.1.2; Serf v0.8.2/v0.9.6 request pseudo-version
#|`v0.0.0-20160808181253-ca63d7c062ee`. Main's explicit request is indirect.
#|Current repository source has no target import. Viper v1.15.0 imports the
#|target in production and test source; project `cmd` imports Viper, so target
#|why is positive through `cmd -> Viper -> mapstructure`, and one target
#|package is production/complete-test loaded. The project calls Viper config
#|setup/read APIs but has no repository-local mapstructure Decode/Unmarshal
#|call. Historical Viper, Consul, and Serf versions genuinely import the
#|target; Kong and crypt are metadata-only requesters. Only current Viper is
#|loaded among target requesters; the historical requester why results are
#|negative.
#|
#|A disposable exact-Go-1.26.7
#|`go get github.com/mitchellh/mapstructure@v1.5.0` is a byte no-op. The
#|unchanged project remains 234 selected modules, 3,599 graph edges, 355
#|production entries, 429 complete-test entries, 197 module-backed entries
#|across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` SHA-256 is
#|`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
#|`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
#|The unchanged 432-line tidy projection retains SHA-256
#|`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
#|the common applied 52/948-line hashes remain
#|`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
#|`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
#|No lower projection was retained because the ordinary failure spans the full
#|line.
#|
#|All 42 earlier guarded selections remain exact. Their 247 sorted incoming
#|edges retain SHA-256
#|`4caf9e803c1b324d2afcd17b9ba664e6d2c8d220cd74b264484107ea3dd93f42`.
#|Thirty-seven why results remain negative; only closed kr/pretty, kr/text,
#|emoji/v2, Promptui, and go-homedir are positive. Promptui and go-homedir are
#|the only guarded repository imports; emoji/v2, Promptui, and go-homedir are
#|the only loaded guarded modules. Every go-homedir, Promptui, emoji/v2,
#|kr/text, and kr/pty expiry guard remains exact.
#|
#|Exact-version OSV and GitHub global queries are empty for all 17 releases,
#|and the repository advisory endpoint is empty. Pinned isolated govulncheck
#|v1.8.0 under exact Go 1.26.7 has no target module/package/symbol/test-symbol
#|result for any stable. Base and no-op-v1.5.0 project populations are
#|identical at 30 module, 22 package, 20 symbol, and 20 test-symbol OSV IDs,
#|with no target trace. Guard OSV retains only Gorilla WebSocket GO-2026-6278/
#|GHSA-w67g-5rqw-f597 and go-retryablehttp GO-2024-2947/
#|GHSA-v6v8-xj6m-xwqh; x/mod v0.14.0 retains GO-2026-6179 and GO-2026-6180.
#|The Go index remains 518,501 bytes/1,402 records at SHA-256
#|`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
#|Last-Modified 2026-09-17T17:29:18Z. The PUBLISHED 2,807-byte memberlist CNA
#|response remains at SHA-256
#|`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
#|Advisory absence does not override the ordinary failure or unsupported owner.
#|
#|Final unchanged-project exact Go 1.26.7 module verification, build,
#|count-one tests, race count-one tests, and vet pass under `umask 022`. The Go
#|1.18, source/API/CLI/help/launcher/Make/quality contracts and every earlier
#|decision remain exact. Accepted quality remains 27/27 Q0-Q2 PASS at L2.
#|Every task-owned scratch artifact was contained beneath the managed session
#|scratch root and removed; only the pre-existing launcher-owned Node compile
#|cache remains. Treat all completed release, source, behavior, closure,
#|route, projection, advisory, and final-gate findings as final.
#|
#|# Choose Exactly One Direction
#|
#|1. Explicitly retain exact selected indirect, production-loaded
#|   `github.com/mitchellh/mapstructure v1.5.0` without product-source or
#|   dependency-metadata changes under a mapstructure-specific,
#|   non-transferable exception. Call it unqualified: no canonical exact-path
#|   stable qualifies and no exact-path stable has a genuine supported owner.
#|   Accept only the completed release/source, ordinary behavior, exact owner/
#|   request, route, why/import/load/runtime, graph/tidy/Go-floor, earlier-
#|   guard, and advisory facts. Define the exact expiry guards below. This
#|   exception must not broaden or expire any closed exception.
#|2. Authorize exactly one later bounded measurement-only genuine owner/request
#|   study named `Ply Viper Mapstructure Ownership Study`. Its sole scope is
#|   the existing main indirect v1.5.0 request, current `cmd -> Viper v1.15.0
#|   -> mapstructure v1.5.0` import/load route, historical requester provenance,
#|   and the owner-declared `go-viper/mapstructure` alternate-path fork. Its
#|   exact question is whether a genuine supported Go-1.18-compatible owner
#|   route can remove the unsupported exact-path request while preserving
#|   public API/CLI behavior, all 42 earlier selections/247 edges, and every
#|   closed exception. Do not run the study, change source or any selection,
#|   reopen Viper or an earlier decision, or introduce a fork/replacement in
#|   this decision-recording move; prepare one reciprocal measurement-only
#|   successor.
#|3. Stop P7 unresolved with no product-source, dependency-metadata, roadmap,
#|   or P8 implementation change and no successor execution.
#|
#|Do not invent a fourth option or combine options. Physical selection,
#|indirectness, positive why, production loading, and advisory absence are not
#|qualification. Do not call v1.5.0 qualified, select a lower release, add or
#|change a root, use the unreleased branch, a fork, replacement, pseudo-version,
#|or alternate path, change Viper or any earlier guard, transfer an exception,
#|or begin P8.
#|
#|# Option 1 Expiry Guards
#|
#|If option 1 is selected, preserve and record at minimum:
#|
#|- exact selected path/version `github.com/mitchellh/mapstructure v1.5.0`,
#|  main's indirect request, the other eight exact current/historical requests,
#|  all current requester selections, and every current/historical requester
#|  import or metadata-only boundary;
#|- zero repository target imports, the current `cmd -> Viper -> target` why
#|  route, Viper production/test imports and project config calls, one
#|  production/complete-test target package, and the exact load/runtime-use
#|  boundary;
#|- the exact 17-release line, absent exact-path `/v2` and `/v3`, archived
#|  repository/owner/status/license/default-branch identity, blessed alternate-
#|  path fork identity, release/tag/commit/tree/signature/archive/sumdb/module
#|  facts, the v1.2.0 tag-mutation record, and no replacement/retraction/
#|  deprecation or eligible exact-path alternate;
#|- all 17 stables remaining unqualified for the completed ordinary stale-slice
#|  result, the v1.2.2 repeated-suite result, no future qualified Go-1.18-
#|  compatible exact-path stable, and no genuine supported exact-path owner;
#|- baseline 234/3,599/355/429/197/41/1,067 state, exact module hashes, common
#|  tidy state, no-op v1.5.0 projection, and no dependency implementation;
#|- exact Go 1.18 floor, exact Go 1.26.7 identity, source/API/CLI/help/launcher/
#|  Make/quality contracts, and accepted 27/27 Q0-Q2 PASS at L2;
#|- all 42 earlier selections/247 incoming edges and their exact hash, plus
#|  every go-homedir, Promptui, emoji/v2, kr/text, and kr/pty expiry guard; and
#|- no new target/closure advisory, independent defect, release, owner,
#|  qualified stable, supported exact-path owner, or compatible genuine route
#|  to a qualified target or closed-exception release.
#|
#|Any target/request/requester-import/owner-route, root/why/import/load/runtime,
#|graph/module/tidy/Go-floor, source/behavior/closure, repository/release/owner,
#|advisory/finding, qualification, supported-owner, earlier guard, closed-
#|exception, or compatible-route change expires option 1 and requires a fresh
#|mapstructure dependency and product decision before merge. Any closed-
#|exception change separately requires its corresponding fresh decision. The
#|exception is not a security claim and cannot transfer.
#|
#|# Decision Recording Contract
#|
#|This successor records a product direction only. Reproduce the handoff,
#|chain, launcher, exact target/request/route, module/graph/tidy, closed-
#|exception, guard, owner/release, and fresh narrow advisory premises without
#|repeating the completed source, behavior, closure, upstream, cross-build,
#|projection, archive, or govulncheck evaluation. Run final exact-Go unchanged-
#|project module verification, build, count-one tests, race count-one tests,
#|and vet.
#|
#|If option 1 is chosen, update the roadmap and rolling handover, answer this
#|archive, and prepare exactly one reciprocal evaluation of the next
#|unanswered selected queue item without executing it. If option 2 is chosen,
#|prepare only the named measurement-only successor. If option 3 is chosen,
#|record the stop without preparing a dependency evaluation or entering P8.
#|Make only the required local handoff documentation commit; do not create a
#|dependency implementation commit.
#|
#|# Required Reading
#|
#|Read this archive, the answered mapstructure evaluation, answered go-homedir,
#|Promptui, emoji/v2, kr/text, and kr/pty decisions/evaluations, kr/pretty and
#|Cast owner records, rolling handover, P7/P8 roadmap, `go.mod`, and `go.sum`.
#|Verify branch, ancestry, clean ordinary/ignored state, reciprocal archive
#|chain, launcher check, exact Go identities, module hashes/counts/tidy
#|projection, all nine target requests, requester imports and routes, why/load/
#|runtime state, all 42 earlier guards and the 247-edge snapshot, all five
#|closed exception boundaries, and fresh advisory identities before recording
#|a choice.
#|
#|# Three Moves
#|
#|First, verify every completed premise without reopening the evaluation.
#|Second, choose and record exactly one option, make only its authorized
#|documentation/handoff changes, and run the unchanged-project final gate.
#|Third, verify scratch cleanup, reciprocal chain, single NEXT state, launcher
#|prompt mirror, exact changed set, and clean status; commit the local handoff
#|without executing its successor.
#|
#|# Automatic Handoff
#|
#|Do not launch a successor, push, merge, publish, release, stash, revert,
#|bypass cleanup, remove the worktree, transfer an exception, reopen
#|go-homedir, Promptui, emoji/v2, kr/text, kr/pty, kr/pretty, Cast, Viper, or an
#|earlier decision, repeat the mapstructure evaluation, evaluate another
#|dependency group, write outside the managed scratch root, or begin P8.
# CODEX_MUTABLE_PROMPT_END
