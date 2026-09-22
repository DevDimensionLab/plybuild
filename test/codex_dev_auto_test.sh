#!/usr/bin/env bash

set -euo pipefail
unset CDPATH

SOURCE_ROOT=$(cd -P "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
TEST_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/codex-dev-auto-test.XXXXXX")
TESTS=0

cleanup() {
	find "$TEST_ROOT" -type d -exec chmod u+rwx {} + 2>/dev/null || true
	rm -rf -- "$TEST_ROOT"
}
trap cleanup EXIT HUP INT TERM

fail() {
	printf 'codex_dev_auto_test: FAIL: %s\n' "$*" >&2
	exit 1
}

assert_file_contains() {
	local file=$1
	local expected=$2
	grep -F -- "$expected" "$file" >/dev/null ||
		fail "$file does not contain: $expected"
	TESTS=$((TESTS + 1))
}

assert_equal() {
	local expected=$1
	local actual=$2
	local label=$3
	[[ "$actual" == "$expected" ]] ||
		fail "$label: expected '$expected', got '$actual'"
	TESTS=$((TESTS + 1))
}

make_fixture() {
	local fixture=$1
	mkdir -p "$fixture/docs/plan" "${fixture}-state" "${fixture}-logs"
	cp "$SOURCE_ROOT/codex-dev-auto.sh" "$fixture/codex-dev-auto.sh"
	chmod 755 "$fixture/codex-dev-auto.sh"
	printf '%s\n' 'initial handover' >"$fixture/docs/plan/quality-handover.md"
	printf '%s\n' '0' >"${fixture}-state/phase"
	printf '%s\n' '0' >"${fixture}-state/launcher-attempts"

	cat >"$fixture/codex-dev-start.sh" <<'LAUNCHER'
#!/usr/bin/env bash
set -euo pipefail
phase=$(sed -n '1p' "$AUTO_TEST_STATE_DIR/phase")
if [[ "${1:-}" == '--check' ]]; then
	if [[ "$phase" == '2' ]]; then
		printf '%s\n' 'codex-dev-start: PASS (session fake, COMPLETE)'
	else
		printf '%s\n' 'codex-dev-start: PASS (session fake, NEXT)'
	fi
	exit 0
fi
supervisor="$TMPDIR/codex-dev-start.fake-$$"
mkdir -p "$supervisor/turn-001"
printf 'final phase %s\n' "$phase" >"$supervisor/turn-001/final-message.txt"
printf '%s\n' '{"type":"thread.started","thread_id":"launcher"}' \
	>"$supervisor/turn-001/events.jsonl"
printf 'launcher stderr phase %s\n' "$phase" \
	>"$supervisor/turn-001/codex.stderr"
attempts=$(sed -n '1p' "$AUTO_TEST_STATE_DIR/launcher-attempts")
attempts=$((attempts + 1))
printf '%s\n' "$attempts" >"$AUTO_TEST_STATE_DIR/launcher-attempts"
if [[ "${AUTO_TEST_LAUNCHER_CAPACITY_ONCE:-0}" == '1' && "$attempts" == '1' ]]; then
	printf '%s\n' \
		'codex-dev-start: Codex terminal failure: Selected model is at capacity. Please try a different model.'
	exit 1
fi
printf 'fake launcher phase %s\n' "$phase"
case $phase in
0) exit 1 ;;
1)
	printf '%s\n' '2' >"$AUTO_TEST_STATE_DIR/phase"
	exit 0
	;;
2) exit 0 ;;
*) exit 70 ;;
esac
LAUNCHER

	cat >"$fixture/fake-codex" <<'CODEX'
#!/usr/bin/env bash
set -euo pipefail
final=''
repo=''
while [[ $# -gt 0 ]]; do
	case $1 in
	exec|--ephemeral|--json|-)
		shift
		;;
	-c|--sandbox|-C|--add-dir|--color|--output-schema|--output-last-message)
		key=$1
		value=${2:-}
		if [[ "$key" == '-C' ]]; then repo=$value; fi
		if [[ "$key" == '--output-last-message' ]]; then final=$value; fi
		shift 2
		;;
	*) shift ;;
	esac
done
[[ -n "$final" && -n "$repo" ]] || exit 64
prompt=$(cat)
[[ "$prompt" == *'conservative automatic controller'* ]] || exit 65
cd "$repo"
phase=$(sed -n '1p' "$AUTO_TEST_STATE_DIR/phase")
mode=${AUTO_TEST_DECISION:-normal}
action='blocked'
reason='synthetic blocker'
decision=''
commit=''
status='NEXT'
if [[ "$mode" == 'blocked' ]]; then
	:
elif [[ "$phase" == '0' ]]; then
	if [[ "$mode" == 'bad-scope' ]]; then
		printf '%s\n' 'not authorized' >product-source.go
		git add product-source.go
	else
		printf '%s\n' 'authorized option 1' >docs/plan/quality-handover.md
		git add docs/plan/quality-handover.md
	fi
	git commit -q --no-gpg-sign -m 'docs: authorize bounded synthetic exception'
	commit=$(git rev-parse HEAD)
	printf '%s\n' '1' >"$AUTO_TEST_STATE_DIR/phase"
	action='restart'
	reason='bounded option 1 was authorized'
	decision='option 1 with exact synthetic guards'
elif [[ "$phase" == '2' ]]; then
	action='complete'
	reason='launcher and roadmap are complete'
	status='COMPLETE'
else
	reason='unexpected synthetic phase'
	status='UNKNOWN'
fi
printf '{"action":"%s","reason":"%s","decision":"%s","commit":"%s","session_status":"%s"}\n' \
	"$action" "$reason" "$decision" "$commit" "$status" >"$final"
printf '%s\n' \
	'{"type":"thread.started","thread_id":"fake"}' \
	'{"type":"turn.started"}' \
	'{"type":"item.completed","item":{"type":"agent_message","text":"synthetic decision"}}' \
	'{"type":"turn.completed","usage":{"input_tokens":1,"output_tokens":1}}'
CODEX

	chmod 755 "$fixture/codex-dev-start.sh" "$fixture/fake-codex"
	(
		cd "$fixture"
		git init -q
		git checkout -q -b codex/upgrade-quality
		git config user.name 'Codex Auto Test'
		git config user.email 'codex-auto-test@example.invalid'
		git config commit.gpgsign false
		git add codex-dev-auto.sh codex-dev-start.sh fake-codex docs
		git commit -q --no-gpg-sign -m 'test fixture'
	)
}

run_wrapper() {
	local fixture=$1
	local output=$2
	shift 2
	(
		cd "$fixture"
		AUTO_TEST_STATE_DIR="${fixture}-state" \
		CODEX_DEV_AUTO_LOG_ROOT="${fixture}-logs" \
		CODEX_DECIDER_BIN="$fixture/fake-codex" \
		CODEX_DEV_AUTO_LOG_KEEP_RUNS=2 \
		CODEX_DEV_AUTO_MAX_CYCLES=4 \
		"$@" /bin/bash ./codex-dev-auto.sh
	) >"$output" 2>&1
}

success_fixture="$TEST_ROOT/success"
success_state="${success_fixture}-state"
success_logs="${success_fixture}-logs"
make_fixture "$success_fixture"
(
	cd "$success_fixture"
	AUTO_TEST_STATE_DIR="$success_state" \
	CODEX_DEV_AUTO_LOG_ROOT="$success_logs" \
	CODEX_DECIDER_BIN="$success_fixture/fake-codex" \
		/bin/bash ./codex-dev-auto.sh --check
) >"$success_state/check.out" 2>&1
assert_file_contains "$success_state/check.out" 'codex-dev-auto: PASS'
assert_file_contains "$success_logs/.codex-dev-auto-log-root" 'codex-dev-auto-v1'

unowned_logs="$TEST_ROOT/unowned-logs"
mkdir "$unowned_logs"
printf '%s\n' 'foreign data' >"$unowned_logs/run-foreign"
set +e
(
	cd "$success_fixture"
	AUTO_TEST_STATE_DIR="$success_state" \
	CODEX_DEV_AUTO_LOG_ROOT="$unowned_logs" \
	CODEX_DECIDER_BIN="$success_fixture/fake-codex" \
		/bin/bash ./codex-dev-auto.sh --check
) >"$success_state/unowned.out" 2>&1
unowned_rc=$?
set -e
assert_equal '1' "$unowned_rc" 'unowned log-root exit status'
assert_file_contains "$success_state/unowned.out" \
	'refusing to adopt a non-empty unowned log root'

mkdir "$success_logs/active.lock"
set +e
run_wrapper "$success_fixture" "$success_state/locked.out" env
locked_rc=$?
set -e
assert_equal '1' "$locked_rc" 'active-lock exit status'
assert_file_contains "$success_state/locked.out" 'another controller may be active'
rmdir "$success_logs/active.lock"

run_wrapper "$success_fixture" "$success_state/run.out" env
assert_equal '2' "$(sed -n '1p' "$success_state/phase")" 'successful phase'
assert_file_contains "$success_state/run.out" 'authorized roadmap COMPLETE'
assert_file_contains "$success_fixture/docs/plan/quality-handover.md" 'authorized option 1'
assert_equal '' "$(cd "$success_fixture" && git status --porcelain=v1 --untracked-files=all)" \
	'successful worktree status'
latest=$(sed -n '1p' "$success_logs/latest-run.txt")
assert_file_contains "$latest/run.log" 'authorization committed; restarting'
assert_file_contains "$latest/run.log" 'authorized roadmap COMPLETE'
assert_equal '2' "$(find "$latest" -type f -name decision-final.json | wc -l | tr -d '[:space:]')" \
	'decision output count'
assert_equal '0' "$(find "$latest" -type d -name 'codex-dev-start.*' | wc -l | tr -d '[:space:]')" \
	'raw launcher directory count'
assert_equal '0' "$(find "$latest" -type d -name 'decision-scratch' | wc -l | tr -d '[:space:]')" \
	'decision scratch count'
assert_equal '0' "$(find "$success_logs" -mindepth 1 -maxdepth 1 -type d -name 'active.lock' | wc -l | tr -d '[:space:]')" \
	'controller lock count'
assert_file_contains "$latest/cycle-001/launcher-final-messages.log" 'final phase 0'
(
	cd "$success_fixture"
	AUTO_TEST_STATE_DIR="$success_state" \
	CODEX_DEV_AUTO_LOG_ROOT="$success_logs" \
	CODEX_DECIDER_BIN="$success_fixture/fake-codex" \
		/bin/bash ./codex-dev-auto.sh --status
) >"$success_state/status.out" 2>&1
assert_file_contains "$success_state/status.out" '"action": "complete"'

capacity_fixture="$TEST_ROOT/capacity"
capacity_state="${capacity_fixture}-state"
capacity_logs="${capacity_fixture}-logs"
make_fixture "$capacity_fixture"
run_wrapper "$capacity_fixture" "$capacity_state/run.out" \
	env AUTO_TEST_LAUNCHER_CAPACITY_ONCE=1 \
	CODEX_DEV_AUTO_CAPACITY_RETRY_DELAY_SECONDS=0
assert_equal '2' "$(sed -n '1p' "$capacity_state/phase")" 'capacity-retry phase'
assert_equal '' "$(cd "$capacity_fixture" && git status --porcelain=v1 --untracked-files=all)" \
	'capacity-retry worktree status'
capacity_latest=$(sed -n '1p' "$capacity_logs/latest-run.txt")
assert_file_contains "$capacity_latest/run.log" \
	'model at capacity; retrying valid NEXT session after 0s'
assert_equal '0' "$(find "$capacity_latest/cycle-001" -type f -name 'decision-final.json' | wc -l | tr -d '[:space:]')" \
	'capacity-cycle decision output count'
assert_equal '2' "$(find "$capacity_latest" -type f -name 'decision-final.json' | wc -l | tr -d '[:space:]')" \
	'capacity-run decision output count'
assert_equal '0' "$(find "$capacity_latest" -type d -name 'codex-dev-start.*' | wc -l | tr -d '[:space:]')" \
	'capacity-run raw launcher directory count'

blocked_fixture="$TEST_ROOT/blocked"
blocked_state="${blocked_fixture}-state"
blocked_logs="${blocked_fixture}-logs"
make_fixture "$blocked_fixture"
set +e
run_wrapper "$blocked_fixture" "$blocked_state/run.out" \
	env AUTO_TEST_DECISION=blocked
blocked_rc=$?
set -e
assert_equal '2' "$blocked_rc" 'blocked exit status'
assert_equal '0' "$(sed -n '1p' "$blocked_state/phase")" 'blocked phase'
assert_equal '' "$(cd "$blocked_fixture" && git status --porcelain=v1 --untracked-files=all)" \
	'blocked worktree status'
blocked_latest=$(sed -n '1p' "$blocked_logs/latest-run.txt")
assert_file_contains "$blocked_latest/run.log" 'real blocker retained for human review'

scope_fixture="$TEST_ROOT/scope"
scope_state="${scope_fixture}-state"
scope_logs="${scope_fixture}-logs"
make_fixture "$scope_fixture"
set +e
run_wrapper "$scope_fixture" "$scope_state/run.out" \
	env AUTO_TEST_DECISION=bad-scope
scope_rc=$?
set -e
assert_equal '2' "$scope_rc" 'bad-scope exit status'
scope_latest=$(sed -n '1p' "$scope_logs/latest-run.txt")
assert_file_contains "$scope_latest/run.log" \
	'commit changed files outside the authorization scope'

printf 'codex_dev_auto_test: PASS (%s assertions)\n' "$TESTS"
