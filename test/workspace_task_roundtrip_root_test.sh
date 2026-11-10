#!/usr/bin/env bash
set -euo pipefail
repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
base=${PLY_ROUNDTRIP_ROOT:?set PLY_ROUNDTRIP_ROOT to an existing directory for this test}
work=$(mktemp -d "${base%/}/roundtrip-root-test.XXXXXX")
bg=
trap 'set +e; for p in $bg; do kill -TERM -- "-$p" 2>/dev/null; done; wait; rm -rf -- "$work"' EXIT
work=$(cd "$work" && pwd -P)
mkdir "$work/tmp"
export TMPDIR="$work/tmp/" GOTMPDIR="$work/tmp"
fake=$work/bin log=$work/log
mkdir "$fake"
# Fake go/python3 log their temp environment; they and the binaries fake go builds
# pass, fail, hang or wait for a peer.
cat > "$fake/go" <<'EOF'
#!/usr/bin/env bash
printf 'go-%s %s %s\n' "$1" "$TMPDIR" "$GOTMPDIR" >> "$FAKE_LOG"
if [ "$1" = build ]; then while [ "$1" != -o ]; do shift; done; cp "$(dirname "$0")/act" "$2"; exit 0; fi
exec "$(dirname "$0")/act"
EOF
cat > "$fake/python3" <<'EOF'
#!/usr/bin/env bash
printf 'python3 %s %s\n' "$TMPDIR" "$GOTMPDIR" >> "$FAKE_LOG"
exec "$(dirname "$0")/act"
EOF
cat > "$fake/act" <<'EOF'
#!/usr/bin/env bash
case $FAKE_MODE in
pass) exit 0 ;;
fail) exit 7 ;;
hang) echo "hang $$" >> "$FAKE_LOG"; exec sleep 60 ;;
pair) : > "$FAKE_LOG.$$"
	for _ in $(seq 100); do [ "$(ls "$FAKE_LOG".* | wc -l)" -ge 2 ] && exit 0; sleep 0.1; done; exit 9 ;;
esac
EOF
cat > "$fake/find" <<'EOF'
#!/usr/bin/env bash
[ -z "${FAKE_FIND:-}" ] || exit "$FAKE_FIND"
exec /usr/bin/find "$@"
EOF
chmod +x "$fake/go" "$fake/python3" "$fake/act" "$fake/find"

fail() { printf 'FAIL %s: %s\n' "$script" "$1" >&2; cat "$work/out" >&2 2>/dev/null || true; exit 1; }
run() { # mode root; replaces the calling subshell so wait sees the script itself
	exec env FAKE_MODE="$1" FAKE_LOG="$log" GO="$fake/go" PATH="$fake:$PATH" PLY_ROUNDTRIP_ROOT="$2" \
		/bin/bash "$repo_root/test/$script"
}
expect() { # expected-exit mode root
	set +e; (run "$2" "$3") > "$work/out" 2>&1; rc=$?; set -e
	[ "$rc" = "$1" ] || fail "mode $2 root $3 exited $rc, expected $1"
}
start() { # out mode; the script leads its own process group
	set -m; run "$2" "$root" > "$1" 2>&1 & pid=$!; set +m; bg="$bg $pid"
}
hanging_child() {
	for _ in $(seq 100); do child=$(awk '$1 == "hang" {print $2}' "$log"); [ -n "$child" ] && return; sleep 0.1; done
	fail "hang never started"
}
finished() { # pid expected-exit
	set +e; wait "$1"; rc=$?; set -e
	[ "$rc" = "$2" ] || fail "signalled run exited $rc, expected $2"
	! kill -0 "$child" 2>/dev/null || fail "child $child outlived the run"
}
fresh_root() {
	rm -rf -- "$work/root" "$log" "$log".*; mkdir "$work/root"; : > "$log"
	root=$work/root; printf 'keep\n' > "$root/other-data"
}
only_other_data() {
	[ "$(ls -A "$root")" = other-data ] && [ "$(cat "$root/other-data")" = keep ] || fail "run data left in $root"
}
temps_owned() {
	! awk -v r="$root/" '$1 != "hang" && (index($2, r) != 1 || index($3, r) != 1)' "$log" | grep -q . ||
		fail "temp outside owned run dir: $(cat "$log")"
}

for script in workspace_task_run_roundtrip_test.sh workspace_task_queue_roundtrip_test.sh \
	workflow_handoff_roundtrip_test.sh workspace_work_items_roundtrip_test.sh \
	workspace_task_integration_roundtrip_test.sh workspace_goal_delivery_roundtrip_test.sh \
	workspace_task_queue_removed_source_test.sh makefile_quality_test.sh; do
	case $script in workflow_handoff_*|*_integration_*|*_goal_delivery_*|makefile_*) fail_rc=1 ;; *) fail_rc=7 ;; esac
	fresh_root
	mkdir "$work/git-root" "$work/ply-ws" "$work/ply-ws/.ply" "$work/ply-ws/sub"
	git init -q "$work/git-root"
	ln -s "$root" "$work/link-root"
	for bad in '' relative "$work/missing" "$root/other-data" "$work/link-root" "$work/git-root" "$work/ply-ws/sub"; do
		expect 2 pass "$bad"
	done
	[ ! -s "$log" ] || fail "unsuitable root ran commands"
	[ "$(ls -A "$work/git-root")" = .git ] && [ -z "$(ls -A "$work/ply-ws/sub")" ] || fail "unsuitable root was written"
	only_other_data
	rm -rf -- "$work/git-root" "$work/ply-ws" "$work/link-root"

	fresh_root; expect "$fail_rc" fail "$root"; only_other_data; temps_owned
	fresh_root; FAKE_FIND=42 expect "$fail_rc" fail "$root"; only_other_data

	fresh_root; start "$work/out" hang; hanging_child
	kill -TERM -- "-$pid"; finished "$pid" 143; only_other_data; temps_owned

	# HUP to the script alone: it waits for its child, then exits 129 and cleans up.
	fresh_root; start "$work/out" hang; hanging_child
	kill -HUP "$pid"; sleep 0.5
	kill -0 "$pid" 2>/dev/null && [ "$(ls -A "$root" | wc -l | tr -d ' ')" = 2 ] ||
		fail "HUP ended the run while child $child still ran"
	kill -TERM "$child"; finished "$pid" 129; only_other_data; temps_owned

	case $script in workspace_task_run_*|workspace_task_queue_roundtrip_test.sh) ;; *) printf 'ok %s\n' "$script"; continue ;; esac
	fresh_root; expect 0 pass "$root"
	grep -q 'roundtrip: PASS' "$work/out" || fail "missing PASS line"
	only_other_data; temps_owned

	fresh_root
	start "$work/out.1" pair; p1=$pid
	start "$work/out.2" pair; p2=$pid
	wait "$p1" && wait "$p2" || fail "concurrent runs failed: $(cat "$work/out.1" "$work/out.2")"
	[ "$(awk '{print $2}' "$log" | sort -u | wc -l | tr -d ' ')" = 2 ] || fail "runs did not use separate dirs"
	only_other_data; temps_owned
	printf 'ok %s\n' "$script"
done
printf 'roundtrip root and cleanup: PASS\n'
