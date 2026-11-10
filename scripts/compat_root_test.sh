check_compat_root() {
	local compat="$fake_repo/.build/compatibility" outer="$temp_root/outer-tmp" rc pid child
	export FAKE_LOG="$temp_root/fake.log" TMPDIR="$outer"
	mkdir -p "$compat" "$outer" "$temp_root/elsewhere" "$temp_root/find-bin"
	printf 'keep\n' >"$compat/neighbor"
	printf 'keep\n' >"$outer/sentinel"
	printf '#!/usr/bin/env bash\nexit 42\n' >"$temp_root/find-bin/find"
	chmod +x "$temp_root/find-bin/find"
	expect() { # exit command...
		: >"$FAKE_LOG"
		set +e; (cd / && "${@:2}") >"$temp_root/root.out" 2>&1; rc=$?; set -e
		[[ $rc == "$1" ]] || { cat "$temp_root/root.out" >&2; fail "$2 exited $rc, expected $1"; }
	}
	owned() {
		[[ $(awk '$1 != "hang" { print $2 }' "$FAKE_LOG" | sort -u | wc -l | tr -d ' ') == 1 ]] &&
			! awk -v p="$compat/$prefix." '$1 != "hang" && ($1 != $2 "/" || index($2, p) != 1 || $2 !~ /\/tmp$/)' \
				"$FAKE_LOG" | grep -q . || fail "tool temp outside one owned run: $(cat "$FAKE_LOG")"
		[[ $(ls -A "$compat") == neighbor && $(cat "$compat/neighbor") == keep &&
			$(ls -A "$outer") == sentinel ]] || fail "run data left behind: $(ls -A "$compat" "$outer")"
	}
	expect "${pass_rc:-0}" run_pass; owned
	expect 1 run_fail; owned
	PATH="$temp_root/find-bin:$PATH" expect 1 run_fail; owned

	: >"$FAKE_LOG"
	trap 'kill -TERM -- "-$pid" 2>/dev/null || true; compat_cleanup' EXIT
	set -m; (cd / && FAKE_MODE=hang checker_exec=exec run_pass) >"$temp_root/root.out" 2>&1 & pid=$!; set +m
	for _ in $(seq 100); do child=$(awk '$1 == "hang" { print $2 }' "$FAKE_LOG"); [[ -z $child ]] || break; sleep 0.1; done
	[[ -n $child ]] || fail 'hanging tool never started'
	kill -TERM -- "-$pid"
	set +e; wait "$pid"; rc=$?; set -e
	trap compat_cleanup EXIT
	[[ $rc == 143 ]] || fail "TERM exited $rc, expected 143"
	! kill -0 "$child" 2>/dev/null || fail "child $child outlived the run"
	owned

	for link in ${compat_links:-.build .build/compatibility .build/compatibility/reports}; do
		rm -rf -- "$fake_repo/.build"
		mkdir -p "$fake_repo/$(dirname "$link")"
		ln -s "$temp_root/elsewhere" "$fake_repo/$link"
		REPORT_OUT= expect 2 run_pass
		[[ ! -s $FAKE_LOG && -z $(ls -A "$temp_root/elsewhere") ]] || fail "wrote through symlinked $link"
	done
	rm -rf -- "$fake_repo/.build"
}
